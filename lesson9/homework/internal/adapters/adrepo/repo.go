package adrepo

import (
	"context"
	"fmt"
	"sync"

	"homework9/internal/ads"
)

var (
	ErrNoAd            = fmt.Errorf("ad does not exist")
	ErrAdAlreadyExists = fmt.Errorf("ad already exists")
)

type RepoMap struct {
	storage map[int64]*ads.Ad
	m       sync.RWMutex
}

func New() ads.Repository {
	return &RepoMap{
		storage: make(map[int64]*ads.Ad),
		m:       sync.RWMutex{},
	}
}

func (r *RepoMap) AdByID(_ context.Context, ID int64) (ad *ads.Ad, err error) {
	r.m.RLock()
	defer r.m.RUnlock()
	ad, ok := r.storage[ID]

	if !ok {
		return nil, ErrNoAd
	}

	return ad, nil
}

func (r *RepoMap) AddAd(_ context.Context, ad *ads.Ad) (ID int64, err error) {
	r.m.Lock()
	defer r.m.Unlock()

	_, ok := r.storage[ad.ID]
	if ok {
		return -1, ErrAdAlreadyExists
	}

	ad.ID = int64(len(r.storage))
	r.storage[ad.ID] = ad

	return ad.ID, nil
}

func (r *RepoMap) AdsByPattern(_ context.Context, p *ads.Pattern) (adverts []*ads.Ad, err error) {
	r.m.RLock()
	for _, a := range r.storage {
		if p.Match(a) {
			adverts = append(adverts, a)
		}
	}
	r.m.RUnlock()

	return adverts, nil
}

func (r *RepoMap) DeleteAd(_ context.Context, ID int64) error {
	r.m.Lock()
	defer r.m.Unlock()

	_, ok := r.storage[ID]
	if !ok {
		return ErrNoAd
	}

	delete(r.storage, ID)

	return nil
}
