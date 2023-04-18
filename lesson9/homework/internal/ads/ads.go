package ads

import "time"

type Ad struct {
	ID        int64
	Title     string `validate:"min:1;max:99"`
	Text      string `validate:"min:1;max:499"`
	UserID    int64
	Published bool
	Created   time.Time
	Updated   time.Time
}
