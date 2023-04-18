package userrepo

import (
	"context"
	"fmt"
	"sync"

	"homework9/internal/users"
)

var (
	ErrNoUser            = fmt.Errorf("user does not exist")
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
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

func (r *RepoMap) UserByID(_ context.Context, ID int64) (u *users.User, err error) {
	r.m.RLock()
	defer r.m.RUnlock()
	u, ok := r.storage[ID]

	if !ok {
		return nil, ErrNoUser
	}

	return u, nil
}

func (r *RepoMap) AddUser(_ context.Context, u *users.User) (ID int64, err error) {
	r.m.Lock()
	defer r.m.Unlock()

	_, ok := r.storage[u.ID]
	if ok {
		return -1, ErrUserAlreadyExists
	}

	u.ID = int64(len(r.storage))
	r.storage[u.ID] = u

	return u.ID, nil
}

func (r *RepoMap) DeleteUser(_ context.Context, ID int64) error {
	r.m.Lock()
	defer r.m.Unlock()

	_, ok := r.storage[ID]
	if ok {
		return ErrNoUser
	}

	delete(r.storage, ID)

	return nil
}
