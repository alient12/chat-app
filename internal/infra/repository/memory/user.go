package memory

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/userrepo"
	"context"
	"sync"

	"github.com/labstack/echo/v4"
)

type Repository struct {
	users map[uint64]model.User
	lock  sync.RWMutex
}

func New() *Repository {
	return &Repository{
		users: make(map[uint64]model.User),
		lock:  sync.RWMutex{},
	}
}

func (r *Repository) Add(_ context.Context, m model.User) error {
	r.lock.RLock()
	if _, ok := r.users[m.ID]; ok {
		return userrepo.ErrUserIDDuplicate
	}
	r.lock.RUnlock()

	r.lock.Lock()
	r.users[m.ID] = m
	r.lock.Unlock()

	return nil
}

func (r *Repository) Get(_ context.Context, cmd userrepo.GetCommand) []model.User {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var users []model.User

	if cmd.ID != nil {
		user, ok := r.users[*cmd.ID]
		if !ok {
			return nil
		}

		users = []model.User{user}
	} else {
		for _, user := range r.users {
			users = append(users, user)
		}
	}

	for i := 0; i < len(users); i++ {
		if cmd.Username != nil {
			if users[i].Username != *cmd.Username {
				users = append(users[:i], users[i+1:]...)
				i--
				continue
			}
		}

		if cmd.Phone != nil {
			if users[i].Phone != *cmd.Phone {
				users = append(users[:i], users[i+1:]...)
				i--
				continue
			}
		}
	}

	return users
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	r.lock.RLock()
	if _, ok := r.users[id]; !ok {
		return echo.ErrBadRequest
	}
	r.lock.RUnlock()

	r.lock.Lock()
	delete(r.users, id)
	r.lock.Unlock()

	return nil
}

func (r *Repository) Update(_ context.Context, m model.User) error {
	r.lock.RLock()
	if _, ok := r.users[m.ID]; !ok {
		return echo.ErrNotFound
	}
	r.lock.RUnlock()

	r.lock.Lock()
	r.users[m.ID] = m
	r.lock.Unlock()

	return nil
}
