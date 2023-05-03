package users

import "context"

//go:generate mockery --name Repository
type Repository interface {
	UserByID(ctx context.Context, ID int64) (*User, error)
	AddUser(ctx context.Context, ad *User) (int64, error)
	DeleteUser(ctx context.Context, ID int64) error
}
