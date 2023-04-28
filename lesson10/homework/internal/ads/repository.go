package ads

import (
	"context"
)

//go:generate mockery --name Repository
type Repository interface {
	AdByID(ctx context.Context, ID int64) (*Ad, error)
	AddAd(ctx context.Context, ad *Ad) (int64, error)
	AdsByPattern(ctx context.Context, p *Pattern) ([]*Ad, error)
	DeleteAd(ctx context.Context, ID int64) error
}
