package adrepo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"homework10/internal/adapters/adrepo"
	"homework10/internal/ads"
)

type EmptyRepoTestSuite struct {
	suite.Suite
	repo ads.Repository
}

type FilledRepoTestSuite struct {
	suite.Suite
	repo   ads.Repository
	filled time.Time
}

type RepoTest struct {
	name    string
	in      any
	wantVal any
	wantErr error
}

func (s *EmptyRepoTestSuite) SetupTest() {
	s.repo = adrepo.New()
}

func (s *FilledRepoTestSuite) SetupTest() {
	s.repo = adrepo.New()
	s.filled = time.Now().UTC()
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

func (s *EmptyRepoTestSuite) TestAddAd() {
	tests := []RepoTest{
		{
			name: "ok add 0th ad",
			in: &ads.Ad{
				ID:    -1,
				Title: "0th ad title",
				Text:  "0th ad text",
			},
			wantVal: int64(0),
			wantErr: nil,
		},
		{
			name: "ok add 1th ad",
			in: &ads.Ad{
				ID:    -1,
				Title: "1st ad title",
				Text:  "1st ad text",
			},
			wantVal: int64(1),
			wantErr: nil,
		},
		{
			name: "wrong add 2nd ad (given ID is 0)",
			in: &ads.Ad{
				ID:    0,
				Title: "again 0th ad title",
				Text:  "again 0th ad text",
			},
			wantVal: int64(-1),
			wantErr: adrepo.ErrAdAlreadyExists,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ID, err := s.repo.AddAd(context.Background(), tt.in.(*ads.Ad))
			if err != nil {
				assert.ErrorIs(t, err, adrepo.ErrAdAlreadyExists)
				assert.Equal(t, int64(-1), ID)
			} else {
				assert.Equal(t, tt.wantVal.(int64), ID)
			}
		})
	}
}

func (s *FilledRepoTestSuite) TestAdByID() {
	tests := []RepoTest{
		{
			name: "ok get ad by ID 0",
			in:   int64(0),
			wantVal: &ads.Ad{
				ID:        0,
				Title:     fmt.Sprintf("title in group %d", 0%3),
				Text:      fmt.Sprintf("%d ad text", 0),
				UserID:    int64(0 / 4),
				Published: true,
				Created:   s.filled,
				Updated:   s.filled,
			},
			wantErr: nil,
		},
		{
			name: "ok get ad by ID 5",
			in:   int64(5),
			wantVal: &ads.Ad{
				ID:        5,
				Title:     fmt.Sprintf("title in group %d", 5%3),
				Text:      fmt.Sprintf("%d ad text", 5),
				UserID:    int64(5 / 4),
				Published: false,
				Created:   s.filled,
				Updated:   s.filled,
			},
		},
		{
			name: "ok get ad by ID 11",
			in:   int64(11),
			wantVal: &ads.Ad{
				ID:        11,
				Title:     fmt.Sprintf("title in group %d", 11%3),
				Text:      fmt.Sprintf("%d ad text", 11),
				UserID:    int64(11 / 4),
				Published: false,
				Created:   s.filled,
				Updated:   s.filled,
			},
		},
		{
			name:    "wrong get ad by ID 12",
			in:      int64(12),
			wantVal: nil,
			wantErr: adrepo.ErrNoAd,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ad, err := s.repo.AdByID(context.Background(), tt.in.(int64))
			if err != nil {
				assert.ErrorIs(t, err, adrepo.ErrNoAd)
				assert.Nil(t, ad)
			} else {
				assert.Equal(t, tt.wantVal.(*ads.Ad).ID, ad.ID)
				assert.Equal(t, tt.wantVal.(*ads.Ad).Title, ad.Title)
				assert.Equal(t, tt.wantVal.(*ads.Ad).Text, ad.Text)
				assert.Equal(t, tt.wantVal.(*ads.Ad).UserID, ad.UserID)
				assert.Equal(t, tt.wantVal.(*ads.Ad).Published, ad.Published)

				assert.Equal(t, tt.wantVal.(*ads.Ad).Created.Year(), ad.Created.Year())
				assert.Equal(t, tt.wantVal.(*ads.Ad).Created.Month(), ad.Created.Month())
				assert.Equal(t, tt.wantVal.(*ads.Ad).Created.Day(), ad.Created.Day())

				assert.Equal(t, tt.wantVal.(*ads.Ad).Updated.Year(), ad.Updated.Year())
				assert.Equal(t, tt.wantVal.(*ads.Ad).Updated.Month(), ad.Updated.Month())
				assert.Equal(t, tt.wantVal.(*ads.Ad).Updated.Day(), ad.Updated.Day())
			}
		})
	}
}

func (s *FilledRepoTestSuite) TestAdsByPattern() {
	tests := []RepoTest{
		{
			name: "ok get ads by title \"title in group 1\"",
			in:   &ads.Pattern{},
			wantVal: []*ads.Ad{
				{
					ID:        -1,
					Title:     fmt.Sprintf("title in group %d", 1%3),
					Text:      fmt.Sprintf("%d ad text", 1),
					UserID:    int64(1 / 4),
					Published: false,
					Created:   s.filled,
					Updated:   s.filled,
				},
				{
					ID:        -1,
					Title:     fmt.Sprintf("title in group %d", 4%3),
					Text:      fmt.Sprintf("%d ad text", 4),
					UserID:    int64(4 / 4),
					Published: true,
					Created:   s.filled,
					Updated:   s.filled,
				},
				{
					ID:        -1,
					Title:     fmt.Sprintf("title in group %d", 7%3),
					Text:      fmt.Sprintf("%d ad text", 7),
					UserID:    int64(7 / 4),
					Published: false,
					Created:   s.filled,
					Updated:   s.filled,
				},
				{
					ID:        -1,
					Title:     fmt.Sprintf("title in group %d", 10%3),
					Text:      fmt.Sprintf("%d ad text", 10),
					UserID:    int64(10 / 4),
					Published: true,
					Created:   s.filled,
					Updated:   s.filled,
				},
			},
			wantErr: nil,
		},
	}
}

func TestEmptyRepoTestSuite(t *testing.T) {
	suite.Run(t, new(EmptyRepoTestSuite))
}

func TestFilledRepoTestSuite(t *testing.T) {
	suite.Run(t, new(FilledRepoTestSuite))
}
