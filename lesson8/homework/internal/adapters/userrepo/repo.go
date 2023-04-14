package userrepo

import (
	"context"
	"errors"
	"homework8/internal/users"
	"sync"
)

type RepoMap struct {
	storage map[int64]*users.User
	m       sync.RWMutex
}

func New() users.Repository {
	return &RepoMap{
		storage: make(map[int64]*users.User),
		m:       sync.RWMutex{},
	}
}

func (r *RepoMap) UserById(_ context.Context, ID int64) (u *users.User, err error) {
	r.m.RLock()
	u, ok := r.storage[ID]
	r.m.RUnlock()

	if !ok {
		return nil, errors.New("ad doesn't exist")
	}

	return u, nil
}

func (r *RepoMap) AddUser(_ context.Context, u *users.User) (ID int64, err error) {
	r.m.Lock()
	defer r.m.Unlock()

	_, ok := r.storage[u.ID]
	if ok {
		return -1, errors.New("user already exists")
	}

	u.ID = int64(len(r.storage))
	r.storage[u.ID] = u

	return u.ID, nil
}
