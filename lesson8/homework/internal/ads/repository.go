package ads

import (
	"context"
)

type Repository interface {
	AdById(ctx context.Context, id int64) (*Ad, error)
	AddAd(ctx context.Context, ad *Ad) (int64, error)
	AdsByFilters(ctx context.Context, f *Filter) (adverts []*Ad, err error)
}

/*AdsByPublished(ctx context.Context, published bool) (ads []*Ad, err error)
AdsByUserID(ctx context.Context, userID int64) (adverts []*Ad, err error)
AdsByCreated(ctx context.Context, created time.Time) (adverts []*Ad, err error)
AdsByTitle(ctx context.Context, title string) (adverts []*Ad, err error)*/
