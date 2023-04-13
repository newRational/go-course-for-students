package app

import (
	"context"
	"fmt"
	"time"

	"github.com/newRational/vld"

	"homework8/internal/ads"
	"homework8/internal/users"
)

type App interface {
	CreateAd(ctx context.Context, title, text string, userID int64) (*ads.Ad, error)
	UpdateAd(ctx context.Context, ID, userID int64, title, text string) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, ID, userID int64, published bool) (*ads.Ad, error)
	AdsByFilter(ctx context.Context, o *ads.Filter) ([]*ads.Ad, error)

	CreateUser(ctx context.Context, nick, email string) (*users.User, error)
	UpdateUser(ctx context.Context, ID int64, nick, email string) (*users.User, error)
}

type AdApp struct {
	adRepo   ads.Repository
	userRepo users.Repository
}

var ErrBadRequest = fmt.Errorf("bad request")
var ErrForbidden = fmt.Errorf("forbidden")

func NewApp(adRepo ads.Repository, userRepo users.Repository) App {
	return &AdApp{
		adRepo:   adRepo,
		userRepo: userRepo,
	}
}

func (a *AdApp) CreateAd(ctx context.Context, title, text string, userID int64) (*ads.Ad, error) {
	_, err := a.userRepo.UserById(ctx, userID)
	if err != nil {
		return nil, ErrBadRequest
	}

	ad := &ads.Ad{
		ID:      -1,
		Title:   title,
		Text:    text,
		UserID:  userID,
		Created: time.Now(),
		Changed: time.Now(),
	}
	if err := vld.Validate(*ad); err != nil {
		return nil, ErrBadRequest
	}

	id, err := a.adRepo.AddAd(ctx, ad)
	if err != nil {
		return nil, err
	}

	ad.ID = id
	return ad, nil
}

func (a *AdApp) UpdateAd(ctx context.Context, ID, userID int64, title, text string) (*ads.Ad, error) {
	ad, err := a.adRepo.AdById(ctx, ID)
	if err != nil {
		return nil, err
	}

	_, err = a.userRepo.UserById(ctx, userID)
	if err != nil {
		return nil, ErrBadRequest
	}

	if ad.UserID != userID {
		return nil, ErrForbidden
	}

	if err = vld.Validate(ads.Ad{Title: title, Text: text}); err != nil {
		return nil, ErrBadRequest
	}

	ad.Text = text
	ad.Title = title
	ad.Changed = time.Now()

	return ad, nil
}

func (a *AdApp) ChangeAdStatus(ctx context.Context, ID, userID int64, published bool) (*ads.Ad, error) {
	ad, err := a.adRepo.AdById(ctx, ID)
	if err != nil {
		return nil, err
	}

	_, err = a.userRepo.UserById(ctx, userID)
	if err != nil {
		return nil, ErrBadRequest
	}

	if ad.UserID != userID {
		return nil, ErrForbidden
	}

	ad.Published = published
	ad.Changed = time.Now()

	return ad, nil
}

func (a *AdApp) AdsByFilter(ctx context.Context, f *ads.Filter) ([]*ads.Ad, error) {
	adverts, err := a.adRepo.AdsByFilters(ctx, f)
	if err != nil {
		return nil, err
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
	if err != nil {
		return nil, err
	}

	u.ID = id
	return u, nil
}

func (a *AdApp) UpdateUser(ctx context.Context, ID int64, nick, email string) (*users.User, error) {
	u, err := a.userRepo.UserById(ctx, ID)
	if err != nil {
		return nil, ErrBadRequest
	}

	if err = vld.Validate(users.User{Nickname: nick, Email: email}); err != nil {
		return nil, ErrBadRequest
	}

	u.Nickname = nick
	u.Email = email

	return u, nil
}
