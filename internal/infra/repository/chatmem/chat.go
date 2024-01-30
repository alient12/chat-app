package chatmem

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/chatrepo"
	"chatapp/internal/util"
	"context"
	"sync"

	"github.com/labstack/echo/v4"
)

type Repository struct {
	chats map[uint64]model.Chat
	lock  sync.RWMutex
}

func New() *Repository {
	return &Repository{
		chats: make(map[uint64]model.Chat),
		lock:  sync.RWMutex{},
	}
}

func (r *Repository) Add(_ context.Context, m model.Chat) error {
	r.lock.RLock()
	if _, ok := r.chats[m.ID]; ok {
		return chatrepo.ErrChatIDDuplicate
	}
	r.lock.RUnlock()

	r.lock.Lock()
	r.chats[m.ID] = m
	r.lock.Unlock()

	return nil
}

func (r *Repository) Get(_ context.Context, cmd chatrepo.GetCommand) []model.Chat {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var chats []model.Chat

	if cmd.ID != nil {
		chat, ok := r.chats[*cmd.ID]
		if !ok {
			return nil
		}

		chats = []model.Chat{chat}
	} else {
		for _, chat := range r.chats {
			chats = append(chats, chat)
		}
	}

	for i := 0; i < len(chats); i++ {
		if cmd.UserID != nil {
			if !util.IsSubset(*cmd.UserID, chats[i].People) {
				chats = append(chats[:i], chats[i+1:]...)
				i--
				continue
			}
		}
	}

	return chats
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	r.lock.RLock()
	if _, ok := r.chats[id]; !ok {
		return echo.ErrBadRequest
	}
	r.lock.RUnlock()

	r.lock.Lock()
	delete(r.chats, id)
	r.lock.Unlock()

	return nil
}

func (r *Repository) Update(_ context.Context, m model.Chat) error {
	r.lock.RLock()
	if _, ok := r.chats[m.ID]; !ok {
		return echo.ErrNotFound
	}
	r.lock.RUnlock()

	r.lock.Lock()
	r.chats[m.ID] = m
	r.lock.Unlock()

	return nil
}
