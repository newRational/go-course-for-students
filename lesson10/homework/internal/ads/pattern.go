package ads

import "time"

type Pattern struct {
	TitleFits     func(title string) bool
	TextFits      func(text string) bool
	UserIDFits    func(userID int64) bool
	PublishedFits func(published bool) bool
	CreatedFits   func(created time.Time) bool
	UpdatedFits   func(updated time.Time) bool
}

func DefaultPattern() *Pattern {
	return &Pattern{
		TitleFits: func(string) bool {
			return true
		},
		TextFits: func(string) bool {
			return true
		},
		UserIDFits: func(int64) bool {
			return true
		},
		PublishedFits: func(bool) bool {
			return true
		},
		CreatedFits: func(time.Time) bool {
			return true
		},
		UpdatedFits: func(time.Time) bool {
			return true
		},
	}
}

func (p *Pattern) Fits(ad *Ad) bool {
	if !p.TitleFits(ad.Title) {
		return false
	}
	if !p.TextFits(ad.Text) {
		return false
	}
	if !p.UserIDFits(ad.UserID) {
		return false
	}
	if !p.PublishedFits(ad.Published) {
		return false
	}
	if !p.CreatedFits(ad.Created) {
		return false
	}
	if !p.UpdatedFits(ad.Updated) {
		return false
	}
	return true
}

func (p *Pattern) SetTitleFits(f func(string) bool) *Pattern {
	pat := *p
	pat.TitleFits = f
	return &pat
}

func (p *Pattern) SetTextFits(f func(string) bool) *Pattern {
	pat := *p
	pat.TextFits = f
	return &pat
}

func (p *Pattern) SetUserIDFits(f func(int64) bool) *Pattern {
	pat := *p
	pat.UserIDFits = f
	return &pat
}

func (p *Pattern) SetPublishedFits(f func(bool) bool) *Pattern {
	pat := *p
	pat.PublishedFits = f
	return &pat
}

func (p *Pattern) SetCreatedFits(f func(time.Time) bool) *Pattern {
	pat := *p
	pat.CreatedFits = f
	return &pat
}

func (p *Pattern) SetUpdatedFits(f func(time.Time) bool) *Pattern {
	pat := *p
	pat.UpdatedFits = f
	return &pat
}

//type Pattern Ad
//
//func NewPattern() *Pattern {
//	return &Pattern{
//		UserID:    -1,
//		Published: true,
//	}
//}
//
//func (p *Pattern) Match(ad *Ad) bool {
//	if p.Title != "" && p.Title != ad.Title {
//		return false
//	}
//	if p.Text != "" && p.Text != ad.Text {
//		return false
//	}
//	if p.UserID != -1 && p.UserID != ad.UserID {
//		return false
//	}
//	if p.Published != ad.Published {
//		return false
//	}
//	t := time.Time{}
//	if p.Created != t {
//		pY, pM, pD := p.Created.UTC().Date()
//		y, m, d := ad.Created.UTC().Date()
//		if y != pY || m != pM || d != pD {
//			return false
//		}
//	}
//	if p.Updated != t {
//		pY, pM, pD := p.Created.UTC().Date()
//		y, m, d := ad.Created.UTC().Date()
//		if y != pY || m != pM || d != pD {
//			return false
//		}
//	}
//	return true
//}
