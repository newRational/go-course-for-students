package users

import "context"

type Repository interface {
	UserById(ctx context.Context, id int64) (*User, error)
	AddUser(ctx context.Context, ad *User) (int64, error)
}
