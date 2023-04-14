package ads

import "time"

type Pattern Ad

func NewPattern() *Pattern {
	return &Pattern{
		UserID:    -1,
		Published: true,
	}
}

func (p *Pattern) Match(ad *Ad) bool {
	if p.Title != "" && p.Title != ad.Title {
		return false
	}
	if p.Text != "" && p.Text != ad.Text {
		return false
	}
	if p.UserID != -1 && p.UserID != ad.UserID {
		return false
	}
	if p.Published != ad.Published {
		return false
	}
	t := time.Time{}
	if p.Created != t {
		pY, pM, pD := p.Created.Date()
		y, m, d := ad.Created.Date()
		if y != pY || m != pM || d != pD {
			return false
		}
	}
	if p.Updated != t {
		pY, pM, pD := p.Created.Date()
		y, m, d := ad.Created.Date()
		if y != pY || m != pM || d != pD {
			return false
		}
	}
	return true
}
