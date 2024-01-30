package handler

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/userrepo"
	"chatapp/internal/infra/http/request"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/argon2"
)

type User struct {
	repo userrepo.Repository
}

var (
	user_count uint64
	mu         sync.Mutex
)

func NewUser(repo userrepo.Repository) *User {
	return &User{
		repo: repo,
	}
}

func GenerateHash(p string, s string) string {
	salt := []byte(s)
	pepper := []byte(os.Getenv("PEPPER"))
	pass := append([]byte(p), salt...)
	pass = append(pass, pepper...)
	hash := argon2.IDKey(pass, salt, 1, 64*1024, 4, 32)
	hashString := fmt.Sprintf("%x", hash)
	return hashString
}

func GenerateID(t model.IDType) uint64 {
	// mask: (3 bits 0) + (61 bits 1) = 000111111111...
	// id: (3 bits t) + (32 bits user_count) + (29 bits random number)
	mask := uint64(1<<61 - 1)
	id := (user_count<<29)&mask | (uint64(rand.Uint32() >> 3))
	id = id | (uint64(t) << 61)
	return id
}

func GetIDType(id uint64) model.IDType {
	return model.IDType(id >> 61)
}

func (u *User) Create(c echo.Context) error {
	var req request.UserCreate

	if err := c.Bind(&req); err != nil {
		log.Print("cannot bind")
		return echo.ErrBadRequest
	}
	if err := req.Validate(); err != nil {
		log.Print("cannot validate")
		return echo.ErrBadRequest
	}

	id := GenerateID(model.UserIDType)
	salt := fmt.Sprintf("%06d ", id)[2:5]
	hash := GenerateHash(req.Password, salt)

	phPtr := &req.Phone
	unPtr := &req.Username

	users := u.repo.Get(c.Request().Context(), userrepo.GetCommand{
		ID:       nil,
		Username: nil,
		Phone:    phPtr,
	})
	if len(users) != 0 {
		return userrepo.ErrPhoneDuplicate
	}

	users = u.repo.Get(c.Request().Context(), userrepo.GetCommand{
		ID:       nil,
		Username: unPtr,
		Phone:    nil,
	})
	if len(users) != 0 {
		return userrepo.ErrUsernameDuplicate
	}

	if err := u.repo.Add(c.Request().Context(), model.User{
		ID:        id,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Phone:     req.Phone,
		Username:  req.Username,
		Password:  hash,
		Image:     req.Image,
		Bio:       req.Bio,
	}); err != nil {
		if errors.Is(err, userrepo.ErrUserIDDuplicate) {
			log.Print("duplicated id")
			return echo.ErrInternalServerError
		} else if errors.Is(err, userrepo.ErrImageSrcDuplicate) {
			log.Print("duplicte image source")
			return echo.ErrBadRequest
		} else {
			return err
		}
	}
	mu.Lock()
	user_count++
	mu.Unlock()

	// cookie based
	// if err := GenJWT(c, id, req.Username); err != nil {
	// 	return echo.ErrInternalServerError
	// }

	token, err := GenJWTLocalStorage(id, req.Username)
	if err != nil {
		return echo.ErrInternalServerError
	}

	id_username := struct {
		ID       uint64
		Username string
		Token    string
	}{ID: id, Username: req.Username, Token: token}

	return c.JSON(http.StatusOK, id_username)
}

func (u *User) Get(c echo.Context) error {
	var idPtr *uint64
	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
		idPtr = &id
	}

	var unPtr *string
	if un := c.QueryParam("username"); un != "" {
		unPtr = &un
	}

	var phPtr *string
	if ph := c.QueryParam("phone"); ph != "" {
		phPtr = &ph
	}

	users := u.repo.Get(c.Request().Context(), userrepo.GetCommand{
		ID:       idPtr,
		Username: unPtr,
		Phone:    phPtr,
	})
	if len(users) == 0 {
		return echo.ErrNotFound
	}

	for i := range users {
		users[i].Password = ""
	}

	return c.JSON(http.StatusOK, users[0])
}

func (u *User) GetByKeyword(c echo.Context) error {
	var keyPtr *string
	if key := c.QueryParam("keyword"); key != "" {
		keyPtr = &key
	}

	users := u.repo.Get(c.Request().Context(), userrepo.GetCommand{
		ID:       nil,
		Username: nil,
		Phone:    nil,
		Keyword:  keyPtr,
	})
	if len(users) == 0 {
		return echo.ErrNotFound
	}

	for i := range users {
		users[i].Password = ""
	}

	return c.JSON(http.StatusOK, users)
}

func (u *User) Delete(c echo.Context) error {
	var idPtr *uint64
	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
		idPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	// check auth
	req := struct {
		Token string `json:"token,omitempty"`
	}{}
	if err := c.Bind(&req); err == nil {
		// check auth by headers
		if ckID, _, err := CheckJWTLocalStorage(req.Token); err != nil {
			return err
		} else {
			if ckID != *idPtr {
				return echo.ErrUnauthorized
			}
		}
	} else {
		// check auth by cookies
		if ckID, _, err := CheckJWT(c); err != nil {
			return err
		} else {
			if ckID != *idPtr {
				return echo.ErrUnauthorized
			}
		}
	}

	if err := u.repo.Delete(c.Request().Context(), *idPtr); err != nil {
		return err
	}

	Logout(c)

	return c.JSON(http.StatusOK, *idPtr)
}

func (u *User) Update(c echo.Context) error {
	var idPtr *uint64
	var unPtr, phPtr *string

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err == nil {
		idPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	var req request.UserCreate

	if err := c.Bind(&req); err != nil {
		log.Print("cannot bind")
		return echo.ErrBadRequest
	}
	if err := req.Validate(); err != nil {
		log.Print("cannot validate")
		return echo.ErrBadRequest
	}

	users := u.repo.Get(c.Request().Context(), userrepo.GetCommand{
		ID:       idPtr,
		Username: unPtr,
		Phone:    phPtr,
	})
	if len(users) == 0 {
		return echo.ErrNotFound
	}
	user := users[0]

	if req.Token != "" {
		// check auth by headers
		if ckID, _, err := CheckJWTLocalStorage(req.Token); err != nil {
			return err
		} else {
			if ckID != *idPtr {
				return echo.ErrUnauthorized
			}
		}
	} else {
		// check auth by cookies
		if ckID, _, err := CheckJWT(c); err != nil {
			return err
		} else {
			if ckID != *idPtr {
				return echo.ErrUnauthorized
			}
		}
	}

	var fn, ln, ph, un, pass, im, bio string
	if req.Firstname != "" {
		fn = req.Firstname
	} else {
		fn = user.Firstname
	}
	if req.Lastname != "" {
		ln = req.Lastname
	} else {
		ln = user.Lastname
	}
	if req.Phone != "" {
		if req.Phone != user.Phone {
			phPtr = &req.Phone
			users := u.repo.Get(c.Request().Context(), userrepo.GetCommand{
				ID:       nil,
				Username: nil,
				Phone:    phPtr,
			})
			if len(users) > 1 {
				return userrepo.ErrPhoneDuplicate
			}
		}
		ph = req.Phone
	} else {
		ph = user.Phone
	}
	if req.Username != "" {
		if req.Username != user.Username {
			unPtr = &req.Username
			users := u.repo.Get(c.Request().Context(), userrepo.GetCommand{
				ID:       nil,
				Username: unPtr,
				Phone:    nil,
			})
			if len(users) > 1 {
				return userrepo.ErrUsernameDuplicate
			}
		}
		un = req.Username
	} else {
		un = user.Username
	}
	if req.Image != "" {
		im = req.Image
	} else {
		im = user.Image
	}
	if req.Bio != "" {
		bio = req.Bio
	} else {
		bio = user.Bio
	}
	if req.Password != "" {
		salt := fmt.Sprintf("%06d ", *idPtr)[2:5]
		pass = GenerateHash(req.Password, salt)
	} else {
		pass = user.Password
	}

	if err := u.repo.Update(c.Request().Context(), model.User{
		ID:        id,
		Firstname: fn,
		Lastname:  ln,
		Phone:     ph,
		Username:  un,
		Password:  pass,
		Image:     im,
		Bio:       bio,
	}); err != nil {
		return err
	}

	// if err := RefJWT(c); err != nil {
	// 	if !errors.Is(err, echo.ErrTooEarly) {
	// 		return err
	// 	}
	// }

	return c.JSON(http.StatusOK, *idPtr)
}

func (u *User) Login(c echo.Context) error {
	var req request.UserLogin

	if err := c.Bind(&req); err != nil {
		log.Print("cannot bind")
		return echo.ErrBadRequest
	}
	if err := req.Validate(); err != nil {
		log.Print("cannot validate")
		return echo.ErrBadRequest
	}

	// check if user is loged in by cookie
	if _, ckUn, err := CheckJWT(c); err != nil {
		if !errors.Is(err, echo.ErrUnauthorized) {
			return err
		}
	} else {
		if ckUn == req.Username {
			if err := RefJWT(c); err != nil {
				return err
			}
			return nil
		}
	}
	unPtr := &req.Username

	users := u.repo.Get(c.Request().Context(), userrepo.GetCommand{
		ID:       nil,
		Username: unPtr,
		Phone:    nil,
	})
	if len(users) == 0 {
		return echo.ErrUnauthorized
	}

	user := users[0]
	salt := fmt.Sprintf("%06d ", user.ID)[2:5]
	hash := GenerateHash(req.Password, salt)
	if user.Password != hash {
		return echo.ErrUnauthorized
	}

	// cookie based
	// if err := GenJWT(c, user.ID, user.Username); err != nil {
	// 	return echo.ErrInternalServerError
	// }
	token, err := GenJWTLocalStorage(user.ID, user.Username)
	if err != nil {
		return echo.ErrInternalServerError
	}

	id_username := struct {
		ID       uint64
		Username string
		Token    string
	}{ID: user.ID, Username: user.Username, Token: token}

	return c.JSON(http.StatusOK, id_username)
}

func (u *User) Register(g *echo.Group) {
	g.POST("/register", u.Create)
	g.POST("/login", u.Login)
	g.GET("/users/:id", u.Get)
	g.GET("/users", u.GetByKeyword)
	g.PATCH("/users/:id", u.Update)
	g.DELETE("/users/:id", u.Delete)
}
