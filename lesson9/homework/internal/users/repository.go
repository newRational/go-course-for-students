package users

import "context"

type Repository interface {
	UserById(ctx context.Context, ID int64) (*User, error)
	AddUser(ctx context.Context, ad *User) (int64, error)
	DeleteUser(ctx context.Context, ID int64) error
}
