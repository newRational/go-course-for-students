package adrepo

import (
	"context"
	"errors"
	"sync"

	"homework8/internal/ads"
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

func (r *RepoMap) AdById(_ context.Context, ID int64) (ad *ads.Ad, err error) {
	r.m.RLock()
	ad, ok := r.storage[ID]
	r.m.RUnlock()

	if !ok {
		return nil, errors.New("ad doesn't exist")
	}

	return ad, nil
}

func (r *RepoMap) AddAd(_ context.Context, ad *ads.Ad) (ID int64, err error) {
	r.m.Lock()
	defer r.m.Unlock()

	_, ok := r.storage[ad.ID]
	if ok {
		return -1, errors.New("ad already exists")
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
