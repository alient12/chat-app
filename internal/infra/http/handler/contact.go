package handler

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/contactrepo"
	"chatapp/internal/domain/repository/userrepo"
	"chatapp/internal/infra/http/request"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Contact struct {
	urepo userrepo.Repository
	repo  contactrepo.Repository
}

func NewContact(repo contactrepo.Repository, urepo userrepo.Repository) *Contact {
	return &Contact{
		urepo: urepo,
		repo:  repo,
	}
}

func (cc *Contact) Create(c echo.Context) error {
	var req request.ContactCreate
	var idPtr *uint64

	if err := c.Bind(&req); err != nil {
		log.Print("cannot bind")
		return echo.ErrBadRequest
	}
	if err := req.Validate(); err != nil {
		log.Print("cannot validate")
		return echo.ErrBadRequest
	}

	if id, err := strconv.ParseUint(c.Param("uid"), 10, 64); err == nil {
		idPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	// check auth
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

	// check for existance of user IDs
	for _, userID := range []uint64{req.ID, *idPtr} {
		users := cc.urepo.Get(c.Request().Context(), userrepo.GetCommand{ID: &userID})
		if len(users) == 0 {
			return echo.ErrBadRequest
		}
	}

	if err := cc.repo.Add(c.Request().Context(), model.Contact{
		UserID:      *idPtr,
		ContactID:   req.ID,
		ContactName: req.Name,
	}); err != nil {
		if errors.Is(err, contactrepo.ErrContactDuplicate) {
			return echo.ErrBadRequest
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, req.Name)
}

func (cc *Contact) Get(c echo.Context) error {
	var idPtr *uint64

	if id, err := strconv.ParseUint(c.Param("uid"), 10, 64); err == nil {
		idPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	// check auth
	token := c.QueryParam("token")
	if token != "" {
		// check auth by query params
		if ckID, _, err := CheckJWTLocalStorage(token); err != nil {
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

	contacts := cc.repo.Get(c.Request().Context(), *idPtr)

	return c.JSON(http.StatusOK, contacts)
}

func (cc *Contact) Delete(c echo.Context) error {
	var cIDPtr, idPtr *uint64

	if id, err := strconv.ParseUint(c.Param("cid"), 10, 64); err == nil {
		cIDPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	if id, err := strconv.ParseUint(c.Param("uid"), 10, 64); err == nil {
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

	if err := cc.repo.Delete(c.Request().Context(), *idPtr, *cIDPtr); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, *cIDPtr)
}

func (cc *Contact) Register(g *echo.Group) {
	g.POST("/users/:uid/contacts", cc.Create)
	g.GET("/users/:uid/contacts", cc.Get)
	g.DELETE("/users/:uid/contacts/:cid", cc.Delete)
}
