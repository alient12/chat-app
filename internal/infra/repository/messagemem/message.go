package messagemem

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/chatrepo"
	"chatapp/internal/domain/repository/messagerepo"
	"context"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
)

type Repository struct {
	messages map[uint64]model.Message
	lock     sync.RWMutex
}

func New() *Repository {
	return &Repository{
		messages: make(map[uint64]model.Message),
		lock:     sync.RWMutex{},
	}
}

func (r *Repository) Add(_ context.Context, m model.Message) error {
	r.lock.RLock()
	if _, ok := r.messages[m.ID]; ok {
		return chatrepo.ErrChatIDDuplicate
	}
	r.lock.RUnlock()

	r.lock.Lock()
	r.messages[m.ID] = m
	r.lock.Unlock()

	return nil
}

func (r *Repository) Get(_ context.Context, cmd messagerepo.GetCommand) []model.Message {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var messages []model.Message

	// not allowed to search just by keyword and content type
	if cmd.ID == nil && cmd.ChatID == nil && cmd.Sender == nil {
		return nil
	}

	if cmd.ID != nil {
		message, ok := r.messages[*cmd.ID]
		if !ok {
			return nil
		}

		messages = []model.Message{message}
	} else {
		for _, message := range r.messages {
			messages = append(messages, message)
		}
	}

	for i := 0; i < len(messages); i++ {
		if cmd.ChatID != nil {
			if messages[i].ChatID != *cmd.ChatID {
				messages = append(messages[:i], messages[i+1:]...)
				i--
				continue
			}
		}

		if cmd.Sender != nil {
			if messages[i].Sender != *cmd.Sender {
				messages = append(messages[:i], messages[i+1:]...)
				i--
				continue
			}
		}

		if cmd.Keyword != nil {
			if !strings.Contains(messages[i].Content, *cmd.Keyword) {
				messages = append(messages[:i], messages[i+1:]...)
				i--
				continue
			}
		}

		if cmd.ContentType != nil {
			if messages[i].ContentType != *cmd.ContentType {
				messages = append(messages[:i], messages[i+1:]...)
				i--
				continue
			}
		}
	}

	return messages
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	r.lock.RLock()
	if _, ok := r.messages[id]; !ok {
		return echo.ErrBadRequest
	}
	r.lock.RUnlock()

	r.lock.Lock()
	delete(r.messages, id)
	r.lock.Unlock()

	return nil
}
