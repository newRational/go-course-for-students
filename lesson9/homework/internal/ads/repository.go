package ads

import (
	"context"
)

type Repository interface {
	AdById(ctx context.Context, ID int64) (*Ad, error)
	AddAd(ctx context.Context, ad *Ad) (int64, error)
	AdsByPattern(ctx context.Context, p *Pattern) ([]*Ad, error)
	DeleteAd(ctx context.Context, ID int64) error
}
