package handler

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/chatrepo"
	"chatapp/internal/domain/repository/userrepo"
	"chatapp/internal/infra/http/request"
	"chatapp/internal/util"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Chat struct {
	urepo userrepo.Repository
	repo  chatrepo.Repository
}

func NewChat(repo chatrepo.Repository, urepo userrepo.Repository) *Chat {
	return &Chat{
		urepo: urepo,
		repo:  repo,
	}
}

var (
	chat_count uint64
)

func GenerateChatID(t model.ChatIDType) uint64 {
	id := (chat_count << 29) | (uint64(rand.Uint32() >> 3))
	id = id | (uint64(t) << 61)
	return id
}

func GetChatIDType(id uint64) model.ChatIDType {
	return model.ChatIDType(id >> 61)
}

func (ch *Chat) Create(c echo.Context) error {
	var req request.ChatCreate

	if err := c.Bind(&req); err != nil {
		log.Print("cannot bind")
		return echo.ErrBadRequest
	}
	if err := req.Validate(); err != nil {
		log.Print("cannot validate")
		return echo.ErrBadRequest
	}

	var people []uint64 = req.People

	// check if user is logged in
	if ckID, _, err := CheckJWT(c); err != nil {
		return err
	} else {
		if !util.InSlice(people, ckID) {
			people = append(people, ckID)
		}
	}

	sort.Slice(people, func(i, j int) bool { return people[i] < people[j] })

	// check for existance of user IDs
	for _, userID := range people {
		users := ch.urepo.Get(c.Request().Context(), userrepo.GetCommand{ID: &userID})
		if len(users) == 0 {
			return echo.ErrBadRequest
		}
	}

	if len(people) == 2 {
		chats := ch.repo.Get(c.Request().Context(), chatrepo.GetCommand{
			ID:     nil,
			UserID: &people,
		})
		for _, chat := range chats {
			// check if dual chat exists
			if len(chat.People) == 2 {
				return echo.ErrBadRequest
			}
		}
	} else {
		// Cannot create group with this api
		return echo.ErrBadRequest
	}

	id := GenerateChatID(1)
	if err := ch.repo.Add(c.Request().Context(), model.Chat{
		ID:        id,
		People:    people,
		CreatedAt: time.Now(),
	}); err != nil {
		if errors.Is(err, chatrepo.ErrChatIDDuplicate) {
			return echo.ErrInternalServerError
		} else if errors.Is(err, chatrepo.ErrDualChatDuplicate) {
			return echo.ErrBadRequest
		} else {
			return err
		}
	}
	mu.Lock()
	chat_count++
	mu.Unlock()

	return c.JSON(http.StatusOK, id)
}

func (ch *Chat) GetByID(c echo.Context) error {
	var chIDPtr, idPtr *uint64
	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
		chIDPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	// check auth
	if ckID, _, err := CheckJWT(c); err != nil {
		return err
	} else {
		idPtr = &ckID
	}

	chat := ch.repo.Get(c.Request().Context(), chatrepo.GetCommand{
		ID:     chIDPtr,
		UserID: nil,
	})[0]

	// check if user has access to the chat
	if !util.InSlice(chat.People, *idPtr) {
		return echo.ErrForbidden
	}

	return c.JSON(http.StatusOK, chat)
}

func (ch *Chat) Get(c echo.Context) error {
	var idPtr *uint64

	// check auth
	if ckID, _, err := CheckJWT(c); err != nil {
		return err
	} else {
		idPtr = &ckID
	}

	people := make([]uint64, 0)
	people = append(people, *idPtr)

	chats := ch.repo.Get(c.Request().Context(), chatrepo.GetCommand{
		ID:     nil,
		UserID: &people,
	})

	return c.JSON(http.StatusOK, chats)
}

func (ch *Chat) Delete(c echo.Context) error {
	var chIDPtr, idPtr *uint64
	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
		chIDPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	// check auth
	if ckID, _, err := CheckJWT(c); err != nil {
		return err
	} else {
		idPtr = &ckID
	}

	chat := ch.repo.Get(c.Request().Context(), chatrepo.GetCommand{
		ID:     chIDPtr,
		UserID: nil,
	})[0]

	// check if user has access to the chat
	if !util.InSlice(chat.People, *idPtr) {
		return echo.ErrForbidden
	}

	if err := ch.repo.Delete(c.Request().Context(), *chIDPtr); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, *chIDPtr)
}

func (ch *Chat) Register(g *echo.Group) {
	g.POST("/chats", ch.Create)
	g.GET("/chats", ch.Get)
	g.GET("/chats/:id", ch.GetByID)
	g.DELETE("/chats/:id", ch.Delete)
}
