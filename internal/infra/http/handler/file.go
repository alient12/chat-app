package handler

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/chatrepo"
	"chatapp/internal/domain/repository/filerepo"
	"chatapp/internal/domain/repository/userrepo"
	"chatapp/internal/util"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type File struct {
	urepo  userrepo.Repository
	chrepo chatrepo.Repository
	repo   filerepo.Repository
}

func NewFile(repo filerepo.Repository, urepo userrepo.Repository, chrepo chatrepo.Repository) *File {
	return &File{
		urepo:  urepo,
		chrepo: chrepo,
		repo:   repo,
	}
}

var (
	File_count uint64
)

func GenerateFileID() uint64 {
	id := (File_count << 32) | (uint64(rand.Uint32()))
	return id
}

func (f *File) Create(c echo.Context) error {
	var chatIDPtr, idPtr *uint64
	var chatIDs []uint64
	var isProfileContent bool = true
	var dir string

	// check auth
	if ckID, _, err := CheckJWT(c); err != nil {
		return err
	} else {
		idPtr = &ckID
	}

	// check if user is uploading file for a chat
	chIDParam := c.Param("chatid")
	if chIDParam != "" {
		id, err := strconv.ParseUint(chIDParam, 10, 64)
		if err != nil {
			return echo.ErrBadRequest
		}
		chatIDPtr = &id

		chat := f.chrepo.Get(c.Request().Context(), chatrepo.GetCommand{
			ID:     chatIDPtr,
			UserID: nil,
		})[0]

		// check if user has access to the chat
		if !util.InSlice(chat.People, *idPtr) {
			return echo.ErrForbidden
		} else {
			chatIDs = append(chatIDs, *chatIDPtr)
		}

		isProfileContent = false
	}

	ufile, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := ufile.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Read the first 512 bytes of the file
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return err
	}

	// Seek to the beginning of the file
	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		// handle error
	}

	// Use these bytes to detect the content type
	contentType := http.DetectContentType(buffer)

	// Destination
	if isProfileContent {
		dir = filepath.Join("files", "profiles")
	} else {
		dir = filepath.Join("files", "chats")
	}
	dir = filepath.Join(dir, fmt.Sprint(*idPtr))
	dstPath := filepath.Join(dir, ufile.Filename)

	// make sure the path exists
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	id := GenerateFileID()

	file := model.File{
		ID:               id,
		UserID:           *idPtr,
		FileName:         ufile.Filename,
		Size:             ufile.Size,
		ContentType:      contentType,
		FilePath:         dstPath,
		ChatIDs:          chatIDs,
		IsProfileContent: isProfileContent,
		CreatedAt:        time.Now(),
	}

	if err := util.AddMetadata(&file); err != nil {
		log.Print(err)
	}

	if err := f.repo.Add(c.Request().Context(), file); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, id)
}

func (f *File) Get(c echo.Context) error {
	var fileIDPtr, idPtr *uint64
	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
		fileIDPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	// check auth
	if ckID, _, err := CheckJWT(c); err != nil {
		return err
	} else {
		idPtr = &ckID
	}

	file := f.repo.Get(c.Request().Context(), filerepo.GetCommand{
		ID:          fileIDPtr,
		UserID:      nil,
		FileName:    nil,
		ContentType: nil,
		ChatID:      nil,
		Keyword:     nil,
	})[0]

	alowedToDownlaod := false
	if file.IsProfileContent {
		alowedToDownlaod = true
	} else {
		// get list of chats that have access to file
		for _, chatid := range file.ChatIDs {
			chat := f.chrepo.Get(c.Request().Context(), chatrepo.GetCommand{
				ID:     &chatid,
				UserID: nil,
			})[0]

			// check if user has access to the chat to download the file
			if util.InSlice(chat.People, *idPtr) {
				alowedToDownlaod = true
				break
			}
		}
	}

	if !alowedToDownlaod {
		return echo.ErrUnauthorized
	}
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=\""+file.FilePath+"\"")
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	return c.File(file.FilePath)
}

func (f *File) Delete(c echo.Context) error {
	var fileIDPtr, idPtr, chatIDPtr *uint64

	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
		fileIDPtr = &id
	} else {
		return echo.ErrBadRequest
	}
	if id, err := strconv.ParseUint(c.Param("chatid"), 10, 64); err == nil {
		chatIDPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	// check auth
	if ckID, _, err := CheckJWT(c); err != nil {
		return err
	} else {
		idPtr = &ckID
	}

	// check if user has access to the chat to delete the file
	chat := f.chrepo.Get(c.Request().Context(), chatrepo.GetCommand{
		ID:     chatIDPtr,
		UserID: nil,
	})[0]

	if !util.InSlice(chat.People, *idPtr) {
		return echo.ErrUnauthorized
	}

	file := f.repo.Get(c.Request().Context(), filerepo.GetCommand{
		ID:          fileIDPtr,
		UserID:      nil,
		FileName:    nil,
		ContentType: nil,
		ChatID:      nil,
	})[0]

	// remove the chatid from file.chatids
	newChatIDs := []uint64{}
	for _, id := range file.ChatIDs {
		if id != *chatIDPtr {
			newChatIDs = append(newChatIDs, id)
		}
	}
	file.ChatIDs = newChatIDs

	// if no chat uses the file, delete the file from server
	if len(file.ChatIDs) == 0 {
		// Delete file from server
		err := os.Remove(file.FilePath)
		if err != nil {
			log.Fatal(err)
		}
		// Delete File model database
		if err := f.repo.Delete(c.Request().Context(), *fileIDPtr); err != nil {
			return err
		}
	} else {
		// Update database with noew File
		if err := f.repo.Update(c.Request().Context(), file); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, *idPtr)
}

func (f *File) DeleteProfileContent(c echo.Context) error {
	var fileIDPtr, idPtr *uint64

	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
		fileIDPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	// check auth
	if ckID, _, err := CheckJWT(c); err != nil {
		return err
	} else {
		idPtr = &ckID
	}

	file := f.repo.Get(c.Request().Context(), filerepo.GetCommand{
		ID:          fileIDPtr,
		UserID:      nil,
		FileName:    nil,
		ContentType: nil,
		ChatID:      nil,
	})[0]

	// check if user has access
	if file.UserID == *idPtr && file.IsProfileContent {
		// Delete file from server
		err := os.Remove(file.FilePath)
		if err != nil {
			log.Fatal(err)
		}
		// Delete File model database
		if err := f.repo.Delete(c.Request().Context(), *fileIDPtr); err != nil {
			return err
		}
	} else {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, *idPtr)
}

func (f *File) Register(g *echo.Group) {
	g.POST("/files/upload/:chatid", f.Create)
	g.GET("/files/download/:id", f.Get)
	g.DELETE("/files/chats/:chatid/:id", f.Delete)
	g.DELETE("/files/:id", f.DeleteProfileContent)
}
