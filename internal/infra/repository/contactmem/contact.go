package contactmem

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/contactrepo"
	"context"
	"sync"

	"github.com/labstack/echo/v4"
)

type Key struct {
	UserID    uint64
	ContactID uint64
}

type Repository struct {
	contacts map[Key]model.Contact
	lock     sync.RWMutex
}

func New() *Repository {
	return &Repository{
		contacts: make(map[Key]model.Contact),
		lock:     sync.RWMutex{},
	}
}

func (r *Repository) Add(_ context.Context, m model.Contact) error {
	key := Key{m.UserID, m.ContactID}
	r.lock.RLock()
	if _, ok := r.contacts[key]; ok {
		return contactrepo.ErrContactDuplicate
	}
	r.lock.RUnlock()

	r.lock.Lock()
	r.contacts[key] = m
	r.lock.Unlock()

	return nil
}

func (r *Repository) Get(_ context.Context, uid uint64) []model.Contact {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var contacts []model.Contact

	for _, contact := range r.contacts {
		contacts = append(contacts, contact)
	}

	for i := 0; i < len(contacts); i++ {
		if contacts[i].UserID != uid {
			contacts = append(contacts[:i], contacts[i+1:]...)
			i--
			continue
		}
	}

	return contacts
}

func (r *Repository) Delete(_ context.Context, uid uint64, cid uint64) error {
	key := Key{uid, cid}
	r.lock.RLock()
	if _, ok := r.contacts[key]; !ok {
		return echo.ErrBadRequest
	}
	r.lock.RUnlock()

	r.lock.Lock()
	delete(r.contacts, key)
	r.lock.Unlock()

	return nil
}
