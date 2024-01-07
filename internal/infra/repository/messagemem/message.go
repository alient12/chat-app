package messagemem

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/chatrepo"
	"context"
	"sync"
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
