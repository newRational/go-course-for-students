package adrepo

import (
	"context"
	"errors"
	"homework6/internal/ads"
	"homework6/internal/app"
)

type RepoMap struct {
	storage map[int64]*ads.Ad
}

func New() app.Repository {
	return &RepoMap{
		storage: make(map[int64]*ads.Ad),
	}
}

func (r *RepoMap) AdById(ctx context.Context, ID int64) (ad *ads.Ad, err error) {
	defer func() {
		if ctxErr := ctx.Err(); ctxErr != nil {
			err = ctxErr
		}
	}()

	ad, ok := r.storage[ID]
	if !ok {
		return nil, errors.New("ad doesn't exist")
	}

	return ad, nil
}

func (r *RepoMap) AddAd(ctx context.Context, ad *ads.Ad) (ID int64, err error) {
	defer func() {
		if ctxErr := ctx.Err(); ctxErr != nil {
			err = ctxErr
		}
	}()

	_, ok := r.storage[ad.ID]
	if ok {
		return -1, errors.New("ad already exists")
	}

	ad.ID = int64(len(r.storage))
	r.storage[ad.ID] = ad

	return ad.ID, nil
}
