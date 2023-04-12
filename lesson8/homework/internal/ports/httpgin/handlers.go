package httpgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"homework8/internal/ads"
	"homework8/internal/app"
	"net/http"
)

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		}

		ad, err := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		}

		adID := c.GetInt("ad_id")

		ad, err := a.ChangeAdStatus(c, int64(adID), reqBody.UserID, reqBody.Published)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		}

		adID := c.GetInt("ad_id")

		ad, err := a.UpdateAd(c, int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для получения всех опубликованных объявлений
func listAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody listAdsRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		}

		adverts, err := a.AdsByFilter(c, createAdFilter(c, reqBody))
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, AdsSuccessResponse(adverts))
	}
}

func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		}

		u, err := a.CreateUser(c, reqBody.Nickname, reqBody.Email)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(u))
	}
}

func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateUserRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		}

		userID := c.GetInt("user_id")

		u, err := a.UpdateUser(c, int64(userID), reqBody.Nickname, reqBody.Email)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(u))
	}
}

func createAdFilter(c *gin.Context, r listAdsRequest) *ads.Filter {
	f := ads.NewFilter()

	f.Title = c.GetString("title")
	f.Created = c.GetTime("created")
	if _, ok := c.GetQuery("user_id"); ok {
		f.UserID = r.UserID
	}
	if _, ok := c.GetQuery("published"); ok {
		f.Published = r.Published
	}

	return f
}
