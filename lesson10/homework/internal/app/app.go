package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/ads"
	"homework10/internal/users"

	"github.com/newRational/vld"
)

//go:generate mockery --name App
type App interface {
	CreateAd(ctx context.Context, title, text string, userID int64) (*ads.Ad, error)
	AdByID(ctx context.Context, ID int64) (*ads.Ad, error)
	AdsByPattern(ctx context.Context, p *ads.Pattern) ([]*ads.Ad, error)
	UpdateAd(ctx context.Context, ID, userID int64, title, text string) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, ID, userID int64, published bool) (*ads.Ad, error)
	DeleteAd(ctx context.Context, ID, userID int64) (*ads.Ad, error)

	CreateUser(ctx context.Context, nick, email string) (*users.User, error)
	UserByID(ctx context.Context, ID int64) (*users.User, error)
	UpdateUser(ctx context.Context, ID int64, nick, email string) (*users.User, error)
	DeleteUser(ctx context.Context, ID int64) (*users.User, error)
}

type AdApp struct {
	adRepo   ads.Repository
	userRepo users.Repository
}

var (
	ErrBadRequest            = fmt.Errorf("bad request")
	ErrForbidden             = fmt.Errorf("forbidden")
	ErrInternalAdRepoError   = fmt.Errorf("internal ad repo error")
	ErrInternalUserRepoError = fmt.Errorf("internal user repo error")
)

func NewApp(adRepo ads.Repository, userRepo users.Repository) App {
	return &AdApp{
		adRepo:   adRepo,
		userRepo: userRepo,
	}
}

func (a *AdApp) CreateAd(ctx context.Context, title, text string, userID int64) (*ads.Ad, error) {
	_, err := a.userRepo.UserByID(ctx, userID)
	if errors.Is(err, userrepo.ErrNoUser) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalUserRepoError
	}

	ad := &ads.Ad{
		ID:      -1,
		Title:   title,
		Text:    text,
		UserID:  userID,
		Created: time.Now().UTC(),
		Updated: time.Now().UTC(),
	}
	if err = vld.Validate(*ad); err != nil {
		return nil, ErrBadRequest
	}

	id, err := a.adRepo.AddAd(ctx, ad)
	if errors.Is(err, adrepo.ErrAdAlreadyExists) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalAdRepoError
	}

	ad.ID = id
	return ad, nil
}

func (a *AdApp) UpdateAd(ctx context.Context, ID, userID int64, title, text string) (*ads.Ad, error) {
	_, err := a.userRepo.UserByID(ctx, userID)
	if errors.Is(err, userrepo.ErrNoUser) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalUserRepoError
	}

	ad, err := a.adRepo.AdByID(ctx, ID)
	if errors.Is(err, adrepo.ErrNoAd) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalAdRepoError
	}

	if ad.UserID != userID {
		return nil, ErrForbidden
	}

	if err = vld.Validate(ads.Ad{Title: title, Text: text}); err != nil {
		return nil, ErrBadRequest
	}

	ad.Text = text
	ad.Title = title
	ad.Updated = time.Now().UTC()

	return ad, nil
}

func (a *AdApp) ChangeAdStatus(ctx context.Context, ID, userID int64, published bool) (*ads.Ad, error) {
	_, err := a.userRepo.UserByID(ctx, userID)
	if errors.Is(err, userrepo.ErrNoUser) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalUserRepoError
	}

	ad, err := a.adRepo.AdByID(ctx, ID)
	if errors.Is(err, adrepo.ErrNoAd) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalAdRepoError
	}

	if ad.UserID != userID {
		return nil, ErrForbidden
	}

	ad.Published = published
	ad.Updated = time.Now().UTC()

	return ad, nil
}

func (a *AdApp) DeleteAd(ctx context.Context, ID, userID int64) (*ads.Ad, error) {
	_, err := a.userRepo.UserByID(ctx, userID)
	if errors.Is(err, userrepo.ErrNoUser) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalUserRepoError
	}

	ad, err := a.adRepo.AdByID(ctx, ID)
	if errors.Is(err, adrepo.ErrNoAd) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalAdRepoError
	}

	if ad.UserID != userID {
		return nil, ErrForbidden
	}

	if err = a.adRepo.DeleteAd(ctx, ID); err != nil {
		return nil, ErrInternalAdRepoError
	}

	return ad, nil
}

func (a *AdApp) AdByID(ctx context.Context, ID int64) (*ads.Ad, error) {
	ad, err := a.adRepo.AdByID(ctx, ID)
	if errors.Is(err, adrepo.ErrNoAd) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalAdRepoError
	}

	return ad, nil
}

func (a *AdApp) AdsByPattern(ctx context.Context, p *ads.Pattern) ([]*ads.Ad, error) {
	adverts, err := a.adRepo.AdsByPattern(ctx, p)
	if errors.Is(err, adrepo.ErrNoAd) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalAdRepoError
	}

	return adverts, nil
}

func (a *AdApp) CreateUser(ctx context.Context, nick, email string) (*users.User, error) {
	u := &users.User{
		ID:       -1,
		Nickname: nick,
		Email:    email,
	}

	if err := vld.Validate(*u); err != nil {
		return nil, ErrBadRequest
	}

	id, err := a.userRepo.AddUser(ctx, u)
	if errors.Is(err, userrepo.ErrUserAlreadyExists) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalUserRepoError
	}

	u.ID = id
	return u, nil
}

func (a *AdApp) UpdateUser(ctx context.Context, ID int64, nick, email string) (*users.User, error) {
	u, err := a.userRepo.UserByID(ctx, ID)
	if errors.Is(err, userrepo.ErrNoUser) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalUserRepoError
	}

	if err = vld.Validate(users.User{Nickname: nick, Email: email}); err != nil {
		return nil, ErrBadRequest
	}

	u.Nickname = nick
	u.Email = email

	return u, nil
}

func (a *AdApp) UserByID(ctx context.Context, ID int64) (*users.User, error) {
	u, err := a.userRepo.UserByID(ctx, ID)
	if errors.Is(err, userrepo.ErrNoUser) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalUserRepoError
	}

	return u, nil
}

func (a *AdApp) DeleteUser(ctx context.Context, ID int64) (*users.User, error) {
	u, err := a.userRepo.UserByID(ctx, ID)
	if errors.Is(err, userrepo.ErrNoUser) {
		return nil, ErrBadRequest
	} else if err != nil {
		return nil, ErrInternalUserRepoError
	}

	if err = a.userRepo.DeleteUser(ctx, ID); err != nil {
		return nil, ErrInternalUserRepoError
	}

	return u, nil
}
