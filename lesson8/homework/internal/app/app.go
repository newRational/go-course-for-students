package app

import (
	"context"
	"fmt"

	"github.com/newRational/vld"

	"homework8/internal/ads"
)

type App interface {
	CreateAd(ctx context.Context, title, text string, authorId int64) (*ads.Ad, error)
	UpdateAd(ctx context.Context, ID, authorId int64, title, text string) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, ID, authorId int64, published bool) (*ads.Ad, error)
	ListAds(ctx context.Context) ([]*ads.Ad, error)
}

type AdApp struct {
	repo Repository
}

type Repository interface {
	AdById(ctx context.Context, id int64) (*ads.Ad, error)
	AddAd(ctx context.Context, ad *ads.Ad) (int64, error)
	Ads(ctx context.Context) (ads []*ads.Ad, err error)
}

var ErrBadRequest = fmt.Errorf("bad request")
var ErrForbidden = fmt.Errorf("forbidden")

func NewApp(repo Repository) App {
	return &AdApp{repo: repo}
}

func (a *AdApp) CreateAd(ctx context.Context, title, text string, authorID int64) (*ads.Ad, error) {
	ad := &ads.Ad{
		ID:       -1,
		Title:    title,
		Text:     text,
		AuthorID: authorID,
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

func (a *AdApp) UpdateAd(ctx context.Context, ID, authorID int64, title, text string) (*ads.Ad, error) {
	ad, err := a.repo.AdById(ctx, ID)
	if err != nil {
		return nil, err
	}

	if ad.AuthorID != authorID {
		return nil, ErrForbidden
	}

	if err = vld.Validate(ads.Ad{Title: title, Text: text}); err != nil {
		return nil, ErrBadRequest
	}

	ad.Text = text
	ad.Title = title

	return ad, nil
}

func (a *AdApp) ChangeAdStatus(ctx context.Context, ID, authorID int64, published bool) (*ads.Ad, error) {
	ad, err := a.repo.AdById(ctx, ID)
	if err != nil {
		return nil, err
	}

	if ad.AuthorID != authorID {
		return nil, ErrForbidden
	}

	ad.Published = published
	return ad, nil
}

func (a *AdApp) ListAds(ctx context.Context) ([]*ads.Ad, error) {
	adverts, err := a.repo.Ads(ctx)
	if err != nil {
		return nil, err
	}

	var pubAds []*ads.Ad

	for _, ad := range adverts {
		if ad.Published == true {
			pubAds = append(pubAds, ad)
		}
	}

	return pubAds, nil
}
