package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/ads"
	"homework8/internal/users"
	"time"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	AuthorID  int64  `json:"author_id"`
	Published bool   `json:"published"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type listAdsRequest struct {
	Title     string    `form:"title"`
	UserID    int64     `form:"user_id"`
	Published bool      `form:"published"`
	Created   time.Time `form:"created" time_format:"2006-01-02"`
}

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email" `
}

type userResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"email"`
}

type updateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"email"`
}

func AdSuccessResponse(ad *ads.Ad) gin.H {
	return gin.H{
		"data": adResponse{
			ID:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorID:  ad.UserID,
			Published: ad.Published,
		},
		"error": nil,
	}
}

func ErrorResponse(err error) gin.H {
	return gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

func AdsSuccessResponse(a []*ads.Ad) gin.H {
	var response []adResponse
	for i := range a {
		response = append(response, adResponse{
			ID:        a[i].ID,
			Title:     a[i].Title,
			Text:      a[i].Text,
			AuthorID:  a[i].UserID,
			Published: a[i].Published,
		})
	}

	return gin.H{
		"data":  response,
		"error": nil,
	}
}

func UserSuccessResponse(u *users.User) gin.H {
	return gin.H{
		"data": userResponse{
			ID:       u.ID,
			Nickname: u.Nickname,
			Email:    u.Email,
		},
		"error": nil,
	}
}

/*
func AdSuccessResponse(ad *ads.Ad) gin.H {
	return gin.H{
		"data": adResponse{
			ID:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorID:  ad.AuthorID,
			Published: ad.Published,
		},
		"error": nil,
	}
}
*/
