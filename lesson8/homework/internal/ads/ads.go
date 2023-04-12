package ads

import "time"

type Ad struct {
	ID        int64
	Title     string `validate:"min:1;max:99"`
	Text      string `validate:"min:1;max:499"`
	UserID    int64
	Published bool
	Created   time.Time
	Changed   time.Time
}

type Filter Ad

func NewFilter() *Filter {
	return &Filter{
		UserID:    -1,
		Published: true,
	}
}

func (o *Filter) Suits(ad *Ad) bool {
	if o.Title != "" && o.Title != ad.Title {
		return false
	}
	if o.Text != "" && o.Text != ad.Text {
		return false
	}
	if o.UserID != -1 && o.UserID != ad.UserID {
		return false
	}
	if o.Published != ad.Published {
		return false
	}
	t := time.Time{}
	if o.Created != t && o.Created != ad.Created {
		return false
	}
	if o.Changed != t && o.Changed != ad.Changed {
		return false
	}
	return true
}
