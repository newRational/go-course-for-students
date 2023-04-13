package users

type User struct {
	ID       int64
	Nickname string `validate:"min:1;max:50"`
	Email    string `validate:"min:1;max:50;email"`
}
