package app

import (
	"context"
	"fmt"
	"github.com/newRational/vld"
	"homework6/internal/ads"
)

type App interface {
	CreateAd(ctx context.Context, title, text string, authorId int64) (*ads.Ad, error)
	UpdateAd(ctx context.Context, ID, authorId int64, title, text string) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, ID, authorId int64, published bool) (*ads.Ad, error)
}

type AdApp struct {
	repo Repository
}

type Repository interface {
	AdById(ctx context.Context, id int64) (*ads.Ad, error)
	AddAd(ctx context.Context, ad *ads.Ad) (int64, error)
}

var ErrBadRequest = fmt.Errorf("bad request")
var ErrForbidden = fmt.Errorf("forbidden")

func NewApp(repo Repository) App {
	return &AdApp{repo: repo}
}

func (a *AdApp) CreateAd(ctx context.Context, title, text string, authorId int64) (*ads.Ad, error) {
	ad := &ads.Ad{
		ID:       -1,
		Title:    title,
		Text:     text,
		AuthorID: authorId,
	}
	if err := vld.Validate(*ad); err != nil {
		return nil, ErrBadRequest
	}

	id, err := a.repo.AddAd(ctx, ad)
	if err != nil {
		return nil, err
	}

	ad.ID = id
	return ad, nil
}

func (a *AdApp) UpdateAd(ctx context.Context, ID, authorId int64, title, text string) (*ads.Ad, error) {
	ad, err := a.repo.AdById(ctx, ID)
	if err != nil {
		return nil, err
	}

	if ad.AuthorID != authorId {
		return nil, ErrForbidden
	}

	if err = vld.Validate(ads.Ad{Title: title, Text: text}); err != nil {
		return nil, ErrBadRequest
	}

	ad.Text = text
	ad.Title = title

	return ad, nil
}

func (a *AdApp) ChangeAdStatus(ctx context.Context, ID, authorId int64, published bool) (*ads.Ad, error) {
	ad, err := a.repo.AdById(ctx, ID)
	if err != nil {
		return nil, err
	}

	if ad.AuthorID != authorId {
		return nil, ErrForbidden
	}

	ad.Published = published
	return ad, nil
}
