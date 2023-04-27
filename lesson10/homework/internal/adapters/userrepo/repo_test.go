package userrepo_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/suite"

	"homework10/internal/adapters/userrepo"
	"homework10/internal/users"
)

type EmptyRepoTestSuite struct {
	suite.Suite
	repo users.Repository
}

type FilledRepoTestSuite struct {
	suite.Suite
	repo users.Repository
}

type RepoTest struct {
	name    string
	in      any
	wantVal any
	wantErr error
}

func (s *EmptyRepoTestSuite) SetupTest() {
	s.repo = userrepo.New()
}

func (s *FilledRepoTestSuite) SetupTest() {
	s.repo = userrepo.New()
	for i := 0; i < 5; i++ {
		_, _ = s.repo.AddUser(context.Background(), &users.User{
			ID:       -1,
			Nickname: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@gmail.com", i),
		})
	}
}

func (s *EmptyRepoTestSuite) TestAddUser() {
	tests := []RepoTest{
		{
			name: "ok add 0th user",
			in: &users.User{
				ID:       -1,
				Nickname: "user0",
				Email:    "user0@gmail.com",
			},
			wantVal: int64(0),
			wantErr: nil,
		},
		{
			name: "ok add 1st user",
			in: &users.User{
				ID:       -1,
				Nickname: "user1",
				Email:    "user1@gmail.com",
			},
			wantVal: int64(1),
			wantErr: nil,
		},
		{
			name: "wrong add 2nd user (given ID=0)",
			in: &users.User{
				ID:       0,
				Nickname: "user0",
				Email:    "user0@gmail.com",
			},
			wantVal: int64(-1),
			wantErr: userrepo.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ID, err := s.repo.AddUser(context.Background(), tt.in.(*users.User))
			if err != nil {
				assert.ErrorIs(t, err, userrepo.ErrUserAlreadyExists)
				assert.Equal(t, int64(-1), ID)
			} else {
				assert.Equal(t, tt.wantVal.(int64), ID)
			}
		})
	}
}

func (s *FilledRepoTestSuite) TestUserByID() {
	tests := []RepoTest{
		{
			name: "ok get user by ID=0",
			in:   int64(0),
			wantVal: &users.User{
				ID:       0,
				Nickname: "user0",
				Email:    "user0@gmail.com",
			},
			wantErr: nil,
		},
		{
			name: "ok get user by ID=3",
			in:   int64(3),
			wantVal: &users.User{
				ID:       3,
				Nickname: "user3",
				Email:    "user3@gmail.com",
			},
			wantErr: nil,
		},
		{
			name: "ok get user by ID=4",
			in:   int64(4),
			wantVal: &users.User{
				ID:       4,
				Nickname: "user4",
				Email:    "user4@gmail.com",
			},
			wantErr: nil,
		},
		{
			name:    "wrong get user by ID=5",
			in:      int64(5),
			wantVal: nil,
			wantErr: userrepo.ErrNoUser,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			u, err := s.repo.UserByID(context.Background(), tt.in.(int64))
			if err != nil {
				assert.ErrorIs(t, err, userrepo.ErrNoUser)
				assert.Nil(t, u)
			} else {
				assert.Equal(t, tt.wantVal.(*users.User).ID, u.ID)
				assert.Equal(t, tt.wantVal.(*users.User).Nickname, u.Nickname)
				assert.Equal(t, tt.wantVal.(*users.User).Email, u.Email)
			}
		})
	}
}

func (s *FilledRepoTestSuite) TestDeleteUser() {
	tests := []RepoTest{
		{
			name:    "ok delete user by ID=0",
			in:      int64(0),
			wantErr: nil,
		},
		{
			name:    "ok delete user by ID=3",
			in:      int64(3),
			wantErr: nil,
		},
		{
			name:    "wrong delete user by ID=4",
			in:      int64(4),
			wantErr: userrepo.ErrNoUser,
		},
		{
			name:    "wrong delete user ID=10",
			in:      int64(10),
			wantErr: userrepo.ErrNoUser,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.repo.DeleteUser(context.Background(), tt.in.(int64))
			if err != nil {
				assert.ErrorIs(t, err, userrepo.ErrNoUser)
			} else {
				_, err = s.repo.UserByID(context.Background(), tt.in.(int64))
				assert.ErrorIs(t, err, userrepo.ErrNoUser)
			}
		})
	}
}

func TestEmptyRepoTestSuite(t *testing.T) {
	suite.Run(t, new(EmptyRepoTestSuite))
}

func TestFilledRepoTestSuite(t *testing.T) {
	suite.Run(t, new(FilledRepoTestSuite))
}
