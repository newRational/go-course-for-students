package adrepo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"homework10/internal/ads"
)

type RepoTestSuite struct {
	suite.Suite
	repo   ads.Repository
	filled time.Time
}

func (s *RepoTestSuite) SetupTest() {
	s.repo = New()
	s.filled = time.Now()
	for i := 0; i < 12; i++ {
		_, _ = s.repo.AddAd(context.Background(), &ads.Ad{
			ID:        -1,
			Title:     fmt.Sprintf("title in group %d", i%3),
			Text:      fmt.Sprintf("%d ad text", i),
			UserID:    int64(i / 4),
			Published: i%2 == 0,
			Created:   s.filled,
			Updated:   s.filled,
		})
	}
}

func (s *RepoTestSuite) TestAddAd() {
	type args struct {
		ctx context.Context
		ad  *ads.Ad
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
		err     error
	}{
		{
			name: "ok add 12th ad",
			args: args{
				ctx: context.Background(),
				ad: &ads.Ad{
					ID:    -1,
					Title: "12th ad title",
					Text:  "12th ad text",
				},
			},
			want:    12,
			wantErr: false,
		},
		{
			name: "ok add 13th ad",
			args: args{
				ctx: context.Background(),
				ad: &ads.Ad{
					ID:    -1,
					Title: "13th ad title",
					Text:  "13th ad text",
				},
			},
			want:    13,
			wantErr: false,
		},
		{
			name: "wrong add 14th ad (given ID=0)",
			args: args{
				ctx: context.Background(),
				ad: &ads.Ad{
					ID:    0,
					Title: "14th ad title",
					Text:  "14th ad text",
				},
			},
			wantErr: true,
			want:    -1,
			err:     ErrAdAlreadyExists,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ID, err := s.repo.AddAd(tt.args.ctx, tt.args.ad)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				assert.Equal(t, int64(-1), ID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, ID)
			}
		})
	}
}

func (s *RepoTestSuite) TestAdByID() {
	type args struct {
		ctx context.Context
		ID  int64
	}
	tests := []struct {
		name    string
		args    args
		want    *ads.Ad
		wantErr bool
		err     error
	}{
		{
			name: "ok get ad by ID=0",
			args: args{
				ctx: context.Background(),
				ID:  0,
			},
			want: &ads.Ad{
				ID:        0,
				Title:     fmt.Sprintf("title in group %d", 0%3),
				Text:      fmt.Sprintf("%d ad text", 0),
				UserID:    0,
				Published: true,
				Created:   s.filled,
				Updated:   s.filled,
			},
			wantErr: false,
		},
		{
			name: "ok get ad by ID=5",
			args: args{
				ctx: context.Background(),
				ID:  5,
			},
			want: &ads.Ad{
				ID:        5,
				Title:     fmt.Sprintf("title in group %d", 5%3),
				Text:      fmt.Sprintf("%d ad text", 5),
				UserID:    1,
				Published: false,
				Created:   s.filled,
				Updated:   s.filled,
			},
			wantErr: false,
		},
		{
			name: "ok get ad by ID=11",
			args: args{
				ctx: context.Background(),
				ID:  11,
			},
			want: &ads.Ad{
				ID:        11,
				Title:     fmt.Sprintf("title in group %d", 11%3),
				Text:      fmt.Sprintf("%d ad text", 11),
				UserID:    2,
				Published: false,
				Created:   s.filled,
				Updated:   s.filled,
			},
			wantErr: false,
		},
		{
			name: "wrong get ad by ID=12",
			args: args{
				ctx: context.Background(),
				ID:  12,
			},
			want:    nil,
			wantErr: true,
			err:     ErrNoAd,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ad, err := s.repo.AdByID(tt.args.ctx, tt.args.ID)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				assert.Nil(t, ad)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, ad.ID)
				assert.Equal(t, tt.want.Title, ad.Title)
				assert.Equal(t, tt.want.Text, ad.Text)
				assert.Equal(t, tt.want.UserID, ad.UserID)
				assert.Equal(t, tt.want.Published, ad.Published)
			}
		})
	}
}

func (s *RepoTestSuite) TestAdsByPattern() {
	type args struct {
		ctx context.Context
		pat *ads.Pattern
	}
	tests := []struct {
		name    string
		args    args
		want    []*ads.Ad
		wantErr bool
		err     error
	}{
		{
			name: "ok get ads by pattern with title=\"title in group 1\" and userID=1",
			args: args{
				ctx: context.Background(),
				pat: ads.DefaultPattern().
					SetTitleFits(func(title string) bool {
						return title == "title in group 1"
					}).
					SetUserIDFits(func(userID int64) bool {
						return userID == 1
					}),
			},
			want: []*ads.Ad{
				{
					ID:        4,
					Title:     fmt.Sprintf("title in group %d", 4%3),
					Text:      fmt.Sprintf("%d ad text", 4),
					UserID:    1,
					Published: true,
					Created:   s.filled,
					Updated:   s.filled,
				},
				{
					ID:        7,
					Title:     fmt.Sprintf("title in group %d", 7%3),
					Text:      fmt.Sprintf("%d ad text", 7),
					UserID:    1,
					Published: false,
					Created:   s.filled,
					Updated:   s.filled,
				},
			},
			wantErr: false,
		},
		{
			name: "ok get ads by pattern with published=true and created=s.filled",
			args: args{
				ctx: context.Background(),
				pat: ads.DefaultPattern().
					SetPublishedFits(func(published bool) bool {
						return published
					}).
					SetCreatedFits(func(created time.Time) bool {
						pY, pM, pD := s.filled.UTC().Date()
						y, m, d := created.UTC().Date()
						return y == pY && m == pM && d == pD
					}),
			},
			want: []*ads.Ad{
				{
					ID:        0,
					Title:     fmt.Sprintf("title in group %d", 0%3),
					Text:      fmt.Sprintf("%d ad text", 0),
					UserID:    0,
					Published: true,
					Created:   s.filled,
					Updated:   s.filled,
				},
				{
					ID:        2,
					Title:     fmt.Sprintf("title in group %d", 2%3),
					Text:      fmt.Sprintf("%d ad text", 2),
					UserID:    0,
					Published: true,
					Created:   s.filled,
					Updated:   s.filled,
				},
				{
					ID:        4,
					Title:     fmt.Sprintf("title in group %d", 4%3),
					Text:      fmt.Sprintf("%d ad text", 4),
					UserID:    1,
					Published: true,
					Created:   s.filled,
					Updated:   s.filled,
				},
				{
					ID:        6,
					Title:     fmt.Sprintf("title in group %d", 6%3),
					Text:      fmt.Sprintf("%d ad text", 6),
					UserID:    1,
					Published: true,
					Created:   s.filled,
					Updated:   s.filled,
				},
				{
					ID:        8,
					Title:     fmt.Sprintf("title in group %d", 8%3),
					Text:      fmt.Sprintf("%d ad text", 8),
					UserID:    2,
					Published: true,
					Created:   s.filled,
					Updated:   s.filled,
				},
				{
					ID:        10,
					Title:     fmt.Sprintf("title in group %d", 10%3),
					Text:      fmt.Sprintf("%d ad text", 10),
					UserID:    2,
					Published: true,
					Created:   s.filled,
					Updated:   s.filled,
				},
			},
			wantErr: false,
		},
		{
			name: "ok get ads by pattern with userID=2 and updated=s.filled",
			args: args{
				ctx: context.Background(),
				pat: ads.DefaultPattern().
					SetUserIDFits(func(userID int64) bool {
						return userID == 2
					}).
					SetUpdatedFits(func(updated time.Time) bool {
						pY, pM, pD := s.filled.UTC().Date()
						y, m, d := updated.UTC().Date()
						return y == pY && m == pM && d == pD
					}),
			},
			want: []*ads.Ad{
				{
					ID:        8,
					Title:     fmt.Sprintf("title in group %d", 8%3),
					Text:      fmt.Sprintf("%d ad text", 8),
					UserID:    2,
					Published: true,
					Created:   s.filled,
					Updated:   s.filled,
				},
				{
					ID:        9,
					Title:     fmt.Sprintf("title in group %d", 9%3),
					Text:      fmt.Sprintf("%d ad text", 9),
					UserID:    2,
					Published: false,
					Created:   s.filled,
					Updated:   s.filled,
				},
				{
					ID:        10,
					Title:     fmt.Sprintf("title in group %d", 10%3),
					Text:      fmt.Sprintf("%d ad text", 10),
					UserID:    2,
					Published: true,
					Created:   s.filled,
					Updated:   s.filled,
				},
				{
					ID:        11,
					Title:     fmt.Sprintf("title in group %d", 11%3),
					Text:      fmt.Sprintf("%d ad text", 11),
					UserID:    2,
					Published: false,
					Created:   s.filled,
					Updated:   s.filled,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			adverts, err := s.repo.AdsByPattern(tt.args.ctx, tt.args.pat)
			assert.Len(t, adverts, len(tt.want))
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				for i := range adverts {
					assert.Contains(t, adverts, tt.want[i])
				}
			}
		})
	}
}

func (s *RepoTestSuite) TestDeleteAd() {
	type args struct {
		ctx context.Context
		ID  int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "ok delete ad by ID=3",
			args: args{
				ctx: context.Background(),
				ID:  3,
			},
			wantErr: false,
		},
		{
			name: "ok delete ad by ID=11",
			args: args{
				ctx: context.Background(),
				ID:  11,
			},
			wantErr: false,
		},
		{
			name: "wrong delete ad by ID=12",
			args: args{
				ctx: context.Background(),
				ID:  12,
			},
			wantErr: true,
			err:     ErrNoAd,
		},
		{
			name: "wrong delete ad by ID=100",
			args: args{
				ctx: context.Background(),
				ID:  100,
			},
			wantErr: true,
			err:     ErrNoAd,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.repo.DeleteAd(tt.args.ctx, tt.args.ID)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				_, err = s.repo.AdByID(tt.args.ctx, tt.args.ID)
				assert.ErrorIs(t, err, ErrNoAd)
			}
		})
	}
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}
