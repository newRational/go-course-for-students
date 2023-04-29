package app

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/ads"
	adrepoMock "homework10/internal/ads/mocks"
	"homework10/internal/users"
	userrepoMock "homework10/internal/users/mocks"
)

type AppTestSuite struct {
	suite.Suite
	adRepo   *adrepoMock.Repository
	userRepo *userrepoMock.Repository
	app      App
}

func (s *AppTestSuite) SetupSuite() {
	s.adRepo = adrepoMock.NewRepository(s.T())
	s.userRepo = userrepoMock.NewRepository(s.T())
	s.app = NewApp(s.adRepo, s.userRepo)
}

func (s *AppTestSuite) TestAdApp_CreateAd() {
	type args struct {
		ctx    context.Context
		title  string
		text   string
		userID int64
	}
	tests := []struct {
		name    string
		args    args
		setMock func()
		want    *ads.Ad
		wantErr bool
		err     error
	}{
		{
			name: "wrong userID",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, userrepo.ErrNoUser).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from userRepo.UserByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalUserRepoError,
		},
		{
			name: "validate error",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "err ad already exists",
			args: args{
				ctx:   context.Background(),
				title: "title",
				text:  "text",
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AddAd", mock.Anything, mock.Anything).
					Return(int64(-1), adrepo.ErrAdAlreadyExists).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from adRepo.AddAd func",
			args: args{
				ctx:   context.Background(),
				title: "title",
				text:  "text",
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AddAd", mock.Anything, mock.Anything).
					Return(int64(-1), fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalAdRepoError,
		},
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				title: "title",
				text:  "text",
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AddAd", mock.Anything, mock.Anything).
					Return(int64(0), nil).
					Once()
			},
			want: &ads.Ad{
				ID:      0,
				Title:   "title",
				Text:    "text",
				UserID:  0,
				Created: time.Now().UTC(),
				Updated: time.Now().UTC(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setMock()
			ad, err := s.app.CreateAd(tt.args.ctx, tt.args.title, tt.args.text, tt.args.userID)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, ad.ID)
				assert.Equal(t, tt.want.Title, ad.Title)
				assert.Equal(t, tt.want.Text, ad.Text)
				assert.Equal(t, tt.want.UserID, ad.UserID)
				assert.InDelta(t, tt.want.Created.Unix(), ad.Created.Unix(), 1)
				assert.InDelta(t, tt.want.Updated.Unix(), ad.Updated.Unix(), 1)
			}
		})
	}
}

func (s *AppTestSuite) TestAdApp_UpdateAd() {
	type args struct {
		ctx    context.Context
		adId   int64
		title  string
		text   string
		userID int64
	}
	tests := []struct {
		name    string
		args    args
		setMock func()
		want    *ads.Ad
		wantErr bool
		err     error
	}{
		{
			name: "wrong userID",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, userrepo.ErrNoUser).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from userRepo.UserByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalUserRepoError,
		},
		{
			name: "err no ad",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(nil, adrepo.ErrNoAd).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from adRepo.AdByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalAdRepoError,
		},
		{
			name: "forbidden request from user",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(&ads.Ad{}, nil).
					Once()
			},
			wantErr: true,
			err:     ErrForbidden,
		},
		{
			name: "validate error",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(&ads.Ad{}, nil).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				title: "new title",
				text:  "new text",
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(&ads.Ad{}, nil).
					Once()
			},
			want: &ads.Ad{
				Title:   "new title",
				Text:    "new text",
				Updated: time.Now().UTC(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setMock()
			ad, err := s.app.UpdateAd(tt.args.ctx, tt.args.adId, tt.args.userID, tt.args.title, tt.args.text)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, ad.ID)
				assert.Equal(t, tt.want.Title, ad.Title)
				assert.Equal(t, tt.want.Text, ad.Text)
				assert.Equal(t, tt.want.UserID, ad.UserID)
				assert.InDelta(t, tt.want.Updated.Unix(), ad.Updated.Unix(), 1)
			}
		})
	}
}

func (s *AppTestSuite) TestAdApp_ChangeAdStatus() {
	type args struct {
		ctx       context.Context
		adId      int64
		published bool
		userID    int64
	}
	tests := []struct {
		name    string
		args    args
		setMock func()
		want    *ads.Ad
		wantErr bool
		err     error
	}{
		{
			name: "wrong userID",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, userrepo.ErrNoUser).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from userRepo.UserByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalUserRepoError,
		},
		{
			name: "err no ad",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(nil, adrepo.ErrNoAd).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from adRepo.AdByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalAdRepoError,
		},
		{
			name: "forbidden request from user",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(&ads.Ad{}, nil).
					Once()
			},
			wantErr: true,
			err:     ErrForbidden,
		},
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(&ads.Ad{}, nil).
					Once()
			},
			want: &ads.Ad{
				Published: true,
				Updated:   time.Now().UTC(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setMock()
			ad, err := s.app.ChangeAdStatus(tt.args.ctx, tt.args.adId, tt.args.userID, tt.args.published)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, ad.ID)
				assert.Equal(t, tt.want.Title, ad.Title)
				assert.Equal(t, tt.want.Text, ad.Text)
				assert.Equal(t, tt.want.UserID, ad.UserID)
				assert.InDelta(t, tt.want.Updated.Unix(), ad.Updated.Unix(), 1)
			}
		})
	}
}

func (s *AppTestSuite) TestAdApp_DeleteAd() {
	type args struct {
		ctx    context.Context
		adId   int64
		userID int64
	}
	tests := []struct {
		name    string
		args    args
		setMock func()
		want    *ads.Ad
		wantErr bool
		err     error
	}{
		{
			name: "wrong userID",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, userrepo.ErrNoUser).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from userRepo.UserByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalUserRepoError,
		},
		{
			name: "err no ad",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(nil, adrepo.ErrNoAd).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from adRepo.AdByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalAdRepoError,
		},
		{
			name: "forbidden request from user",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(&ads.Ad{}, nil).
					Once()
			},
			wantErr: true,
			err:     ErrForbidden,
		},
		{
			name: "unknown error from adRepo.AdByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(&ads.Ad{}, nil).
					Once()

				s.adRepo.
					On("DeleteAd", mock.Anything, mock.Anything).
					Return(fmt.Errorf("unknown error from adRepo.AdByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalAdRepoError,
		},
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(&ads.Ad{}, nil).
					Once()

				s.adRepo.
					On("DeleteAd", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			want:    &ads.Ad{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setMock()
			ad, err := s.app.DeleteAd(tt.args.ctx, tt.args.adId, tt.args.userID)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, ad.ID)
				assert.Equal(t, tt.want.Title, ad.Title)
				assert.Equal(t, tt.want.Text, ad.Text)
				assert.Equal(t, tt.want.UserID, ad.UserID)
				assert.InDelta(t, tt.want.Updated.Unix(), ad.Updated.Unix(), 1)
			}
		})
	}
}

func (s *AppTestSuite) TestAdApp_AdByID() {
	type args struct {
		ctx  context.Context
		adId int64
	}
	tests := []struct {
		name    string
		args    args
		setMock func()
		want    *ads.Ad
		wantErr bool
		err     error
	}{
		{
			name: "err no ad",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(nil, adrepo.ErrNoAd).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from adRepo.AdByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalAdRepoError,
		},
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.adRepo.
					On("AdByID", mock.Anything, mock.Anything).
					Return(&ads.Ad{}, nil).
					Once()
			},
			want:    &ads.Ad{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setMock()
			ad, err := s.app.AdByID(tt.args.ctx, tt.args.adId)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, ad.ID)
				assert.Equal(t, tt.want.Title, ad.Title)
				assert.Equal(t, tt.want.Text, ad.Text)
				assert.Equal(t, tt.want.UserID, ad.UserID)
				assert.InDelta(t, tt.want.Updated.Unix(), ad.Updated.Unix(), 1)
			}
		})
	}
}

func (s *AppTestSuite) TestAdApp_AdsByPattern() {
	type args struct {
		ctx context.Context
		p   *ads.Pattern
	}
	tests := []struct {
		name    string
		args    args
		setMock func()
		want    []*ads.Ad
		wantErr bool
		err     error
	}{
		{
			name: "err no ad",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.adRepo.
					On("AdsByPattern", mock.Anything, mock.Anything).
					Return(nil, adrepo.ErrNoAd).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from adRepo.AdByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.adRepo.
					On("AdsByPattern", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalAdRepoError,
		},
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				p:   ads.DefaultPattern(),
			},
			setMock: func() {
				s.adRepo.
					On("AdsByPattern", mock.Anything, mock.Anything).
					Return([]*ads.Ad{{ID: 0}, {ID: 1}, {ID: 2}}, nil).
					Once()
			},
			want:    []*ads.Ad{{ID: 0}, {ID: 1}, {ID: 2}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setMock()
			adverts, err := s.app.AdsByPattern(tt.args.ctx, tt.args.p)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				for i := range adverts {
					assert.Contains(t, adverts, tt.want[i])
				}
			}
		})
	}
}

func (s *AppTestSuite) TestAdApp_CreateUser() {
	type args struct {
		ctx   context.Context
		nick  string
		email string
	}
	tests := []struct {
		name    string
		args    args
		setMock func()
		want    *users.User
		wantErr bool
		err     error
	}{
		{
			name: "validate error",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "err user already exists",
			args: args{
				ctx:   context.Background(),
				nick:  "user",
				email: "user@gmail.com",
			},
			setMock: func() {
				s.userRepo.
					On("AddUser", mock.Anything, mock.Anything).
					Return(int64(-1), userrepo.ErrUserAlreadyExists).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from adRepo.AddUser func",
			args: args{
				ctx:   context.Background(),
				nick:  "user",
				email: "user@gmail.com",
			},
			setMock: func() {
				s.userRepo.
					On("AddUser", mock.Anything, mock.Anything).
					Return(int64(-1), fmt.Errorf("unknown error from adRepo.AddUser func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalUserRepoError,
		},
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				nick:  "user",
				email: "user@gmail.com",
			},
			setMock: func() {
				s.userRepo.
					On("AddUser", mock.Anything, mock.Anything).
					Return(int64(0), nil).
					Once()
			},
			want: &users.User{
				Nickname: "user",
				Email:    "user@gmail.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setMock()
			u, err := s.app.CreateUser(tt.args.ctx, tt.args.nick, tt.args.email)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, u.ID)
				assert.Equal(t, tt.want.Nickname, u.Nickname)
				assert.Equal(t, tt.want.Email, u.Email)
			}
		})
	}
}

func (s *AppTestSuite) TestAdApp_UpdateUser() {
	type args struct {
		ctx   context.Context
		id    int64
		nick  string
		email string
	}
	tests := []struct {
		name    string
		args    args
		setMock func()
		want    *users.User
		wantErr bool
		err     error
	}{
		{
			name: "wrong userID",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, userrepo.ErrNoUser).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from userRepo.UserByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalUserRepoError,
		},
		{
			name: "validate error",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(&users.User{}, nil).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				nick:  "new.user",
				email: "new.user@gmail.com",
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(&users.User{}, nil).
					Once()
			},
			want: &users.User{
				Nickname: "new.user",
				Email:    "new.user@gmail.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setMock()
			u, err := s.app.UpdateUser(tt.args.ctx, tt.args.id, tt.args.nick, tt.args.email)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, u.ID)
				assert.Equal(t, tt.want.Nickname, u.Nickname)
				assert.Equal(t, tt.want.Email, u.Email)
			}
		})
	}
}

func (s *AppTestSuite) TestAdApp_UserByID() {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		args    args
		setMock func()
		want    *users.User
		wantErr bool
		err     error
	}{
		{
			name: "wrong userID",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, userrepo.ErrNoUser).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from userRepo.UserByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalUserRepoError,
		},
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(&users.User{}, nil).
					Once()
			},
			want:    &users.User{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setMock()
			u, err := s.app.UserByID(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, u.ID)
				assert.Equal(t, tt.want.Nickname, u.Nickname)
				assert.Equal(t, tt.want.Email, u.Email)
			}
		})
	}
}

func (s *AppTestSuite) TestAdApp_DeleteUser() {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		args    args
		setMock func()
		want    *users.User
		wantErr bool
		err     error
	}{
		{
			name: "wrong userID",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, userrepo.ErrNoUser).
					Once()
			},
			wantErr: true,
			err:     ErrBadRequest,
		},
		{
			name: "unknown error from userRepo.UserByID func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("unknown error from userRepo.UserByID func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalUserRepoError,
		},
		{
			name: "unknown error from userRepo.DeleteUser func",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(&users.User{}, nil).
					Once()

				s.userRepo.
					On("DeleteUser", mock.Anything, mock.Anything).
					Return(fmt.Errorf("unknown error from userRepo.DeleteUser func")).
					Once()
			},
			wantErr: true,
			err:     ErrInternalUserRepoError,
		},
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
			},
			setMock: func() {
				s.userRepo.
					On("UserByID", mock.Anything, mock.Anything).
					Return(&users.User{}, nil).
					Once()

				s.userRepo.
					On("DeleteUser", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			want:    &users.User{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setMock()
			u, err := s.app.DeleteUser(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, u.ID)
				assert.Equal(t, tt.want.Nickname, u.Nickname)
				assert.Equal(t, tt.want.Email, u.Email)
			}
		})
	}
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
