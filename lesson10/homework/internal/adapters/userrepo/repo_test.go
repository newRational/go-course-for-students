package userrepo

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"homework10/internal/users"
)

type RepoTestSuite struct {
	suite.Suite
	repo users.Repository
}

func (s *RepoTestSuite) SetupTest() {
	s.repo = New()
	for i := 0; i < 5; i++ {
		_, _ = s.repo.AddUser(context.Background(), &users.User{
			ID:       -1,
			Nickname: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@gmail.com", i),
		})
	}
}

func (s *RepoTestSuite) TestAddUser() {
	type args struct {
		ctx  context.Context
		user *users.User
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
		err     error
	}{
		{
			name: "ok add 5th user",
			args: args{
				ctx: context.Background(),
				user: &users.User{
					ID:       -1,
					Nickname: "user5",
					Email:    "user5@gmail.com",
				},
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "ok add 6th user",
			args: args{
				ctx: context.Background(),
				user: &users.User{
					ID:       -1,
					Nickname: "user6",
					Email:    "user6@gmail.com",
				},
			},
			want:    6,
			wantErr: false,
		},
		{
			name: "wrong add 7th user (given ID=0)",
			args: args{
				ctx: context.Background(),
				user: &users.User{
					ID:       0,
					Nickname: "user7",
					Email:    "user7@gmail.com",
				},
			},
			want:    -1,
			wantErr: true,
			err:     ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ID, err := s.repo.AddUser(tt.args.ctx, tt.args.user)
			if tt.wantErr {
				assert.ErrorIs(t, err, ErrUserAlreadyExists)
				assert.Equal(t, int64(-1), ID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, ID)
			}
		})
	}
}

func (s *RepoTestSuite) TestUserByID() {
	type args struct {
		ctx context.Context
		ID  int64
	}
	tests := []struct {
		name    string
		args    args
		want    *users.User
		wantErr bool
		err     error
	}{
		{
			name: "ok get user by ID=0",
			args: args{
				ctx: context.Background(),
				ID:  0,
			},
			want: &users.User{
				ID:       0,
				Nickname: "user0",
				Email:    "user0@gmail.com",
			},
			wantErr: false,
		},
		{
			name: "ok get user by ID=3",
			args: args{
				ctx: context.Background(),
				ID:  3,
			},
			want: &users.User{
				ID:       3,
				Nickname: "user3",
				Email:    "user3@gmail.com",
			},
			wantErr: false,
		},
		{
			name: "ok get user by ID=4",
			args: args{
				ctx: context.Background(),
				ID:  4,
			},
			want: &users.User{
				ID:       4,
				Nickname: "user4",
				Email:    "user4@gmail.com",
			},
			wantErr: false,
		},
		{
			name: "wrong get user by ID=5",
			args: args{
				ctx: context.Background(),
				ID:  5,
			},
			want:    nil,
			wantErr: true,
			err:     ErrNoUser,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			u, err := s.repo.UserByID(tt.args.ctx, tt.args.ID)
			if tt.wantErr {
				assert.ErrorIs(t, err, ErrNoUser)
				assert.Nil(t, u)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, u.ID)
				assert.Equal(t, tt.want.Nickname, u.Nickname)
				assert.Equal(t, tt.want.Email, u.Email)
			}
		})
	}
}

func (s *RepoTestSuite) TestDeleteUser() {
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
			name: "ok delete user by ID=0",
			args: args{
				ctx: context.Background(),
				ID:  0,
			},
			wantErr: false,
		},
		{
			name: "ok delete user by ID=3",
			args: args{
				ctx: context.Background(),
				ID:  3,
			},
			wantErr: false,
		},
		{
			name: "wrong delete user by ID=5",
			args: args{
				ctx: context.Background(),
				ID:  5,
			},
			wantErr: true,
			err:     ErrNoUser,
		},
		{
			name: "wrong delete user ID=10",
			args: args{
				ctx: context.Background(),
				ID:  10,
			},
			wantErr: true,
			err:     ErrNoUser,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.repo.DeleteUser(tt.args.ctx, tt.args.ID)
			if tt.wantErr {
				assert.ErrorIs(t, err, ErrNoUser)
			} else {
				assert.NoError(t, err)
				_, err = s.repo.UserByID(context.Background(), tt.args.ID)
				assert.ErrorIs(t, err, ErrNoUser)
			}
		})
	}
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}
