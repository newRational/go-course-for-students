package adrepo

import (
	"context"
	"errors"
	"homework8/internal/ads"
)

type RepoMap struct {
	storage map[int64]*ads.Ad
}

func New() ads.Repository {
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

//
//func (r *RepoMap) AdsByPublished(ctx context.Context, published bool) (adverts []*ads.Ad, err error) {
//	defer func() {
//		if ctxErr := ctx.Err(); ctxErr != nil {
//			err = ctxErr
//		}
//	}()
//
//	for _, a := range r.storage {
//		if a.Published == published {
//			adverts = append(adverts, a)
//		}
//	}
//
//	return adverts, nil
//}
//
//func (r *RepoMap) AdsByUserID(ctx context.Context, userID int64) (adverts []*ads.Ad, err error) {
//	defer func() {
//		if ctxErr := ctx.Err(); ctxErr != nil {
//			err = ctxErr
//		}
//	}()
//
//	for _, a := range r.storage {
//		if a.UserID == userID {
//			adverts = append(adverts, a)
//		}
//	}
//
//	return adverts, nil
//}
//
//func (r *RepoMap) AdsByCreated(ctx context.Context, created time.Time) (adverts []*ads.Ad, err error) {
//	defer func() {
//		if ctxErr := ctx.Err(); ctxErr != nil {
//			err = ctxErr
//		}
//	}()
//
//	for _, a := range r.storage {
//		if a.Created.Date() == created.Date() {
//			adverts = append(adverts, a)
//		}
//	}
//
//	return adverts, nil
//}
//
//func (r *RepoMap) AdsByTitle(ctx context.Context, title string) (adverts []*ads.Ad, err error) {
//	defer func() {
//		if ctxErr := ctx.Err(); ctxErr != nil {
//			err = ctxErr
//		}
//	}()
//
//	reg, err := regexp.Compile(title + ".*")
//	if err != nil {
//		return nil, err
//	}
//
//	for _, a := range r.storage {
//		if reg.MatchString(a.Title) {
//			adverts = append(adverts, a)
//		}
//	}
//
//	return adverts, nil
//}

func (r *RepoMap) AdsByFilters(ctx context.Context, f *ads.Filter) (adverts []*ads.Ad, err error) {
	defer func() {
		if ctxErr := ctx.Err(); ctxErr != nil {
			err = ctxErr
		}
	}()

	for _, a := range r.storage {
		if f.Suits(a) {
			adverts = append(adverts, a)
		}
	}

	return adverts, nil
}
