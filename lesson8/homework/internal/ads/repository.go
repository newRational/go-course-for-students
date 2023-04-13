package ads

import (
	"context"
)

type Repository interface {
	AdById(ctx context.Context, id int64) (*Ad, error)
	AddAd(ctx context.Context, ad *Ad) (int64, error)
	AdsByPattern(ctx context.Context, p *Pattern) (adverts []*Ad, err error)
}
