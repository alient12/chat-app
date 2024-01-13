package filemem

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/filerepo"
	"chatapp/internal/util"
	"context"
	"strings"
	"sync"
)

type Repository struct {
	files map[uint64]model.File
	lock  sync.RWMutex
}

func New() *Repository {
	return &Repository{
		files: make(map[uint64]model.File),
		lock:  sync.RWMutex{},
	}
}

func (r *Repository) Add(_ context.Context, m model.File) error {
	r.lock.RLock()
	if _, ok := r.files[m.ID]; ok {
		return filerepo.ErrFileIDDuplicate
	}
	r.lock.RUnlock()

	r.lock.Lock()
	r.files[m.ID] = m
	r.lock.Unlock()

	return nil
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	r.lock.RLock()
	if _, ok := r.files[id]; !ok {
		return filerepo.ErrIDNotFound
	}
	r.lock.RUnlock()

	r.lock.Lock()
	delete(r.files, id)
	r.lock.Unlock()

	return nil
}

func (r *Repository) Update(_ context.Context, m model.File) error {
	r.lock.RLock()
	if _, ok := r.files[m.ID]; !ok {
		return filerepo.ErrIDNotFound
	}
	r.lock.RUnlock()

	r.lock.Lock()
	r.files[m.ID] = m
	r.lock.Unlock()

	return nil
}

func (r *Repository) Get(_ context.Context, cmd filerepo.GetCommand) []model.File {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var files []model.File

	if cmd.ID != nil {
		file, ok := r.files[*cmd.ID]
		if !ok {
			return nil
		}

		files = []model.File{file}
	} else {
		for _, file := range r.files {
			files = append(files, file)
		}
	}

	for i := 0; i < len(files); i++ {
		if cmd.UserID != nil {
			if *cmd.UserID != files[i].UserID {
				files = append(files[:i], files[i+1:]...)
				i--
				continue
			}
		}
		if cmd.UserID != nil {
			if *cmd.UserID != files[i].UserID {
				files = append(files[:i], files[i+1:]...)
				i--
				continue
			}
		}

		if cmd.FileName != nil {
			if *cmd.FileName != files[i].FileName {
				files = append(files[:i], files[i+1:]...)
				i--
				continue
			}
		}

		if cmd.ContentType != nil {
			if *cmd.ContentType != files[i].ContentType {
				files = append(files[:i], files[i+1:]...)
				i--
				continue
			}
		}

		if cmd.ChatID != nil {
			if !util.InSlice(files[i].ChatIDs, *cmd.ChatID) {
				files = append(files[:i], files[i+1:]...)
				i--
				continue
			}
		}

		if cmd.Keyword != nil {
			if !strings.Contains(files[i].FileName, *cmd.Keyword) {
				files = append(files[:i], files[i+1:]...)
				i--
				continue
			}
		}
	}

	return files
}
