package httpgin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"homework9/internal/ads"
	"homework9/internal/app"
)

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ad, err := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
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
			return
		}

		v := c.Param("ad_id")
		adID, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ad, err := a.ChangeAdStatus(c, int64(adID), reqBody.UserID, reqBody.Published)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
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
			return
		}

		adID := c.GetInt64("ad_id")

		ad, err := a.UpdateAd(c, adID, reqBody.UserID, reqBody.Title, reqBody.Text)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для получения объявления по его ID
func showAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		v := c.Param("ad_id")
		adID, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ad, err := a.AdByID(c, int64(adID))
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
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
		var reqParams listAdsRequest

		if err := c.Bind(&reqParams); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		p, err := createAdPattern(c, reqParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		adverts, err := a.AdsByPattern(c, p)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, AdsSuccessResponse(adverts))
	}
}

// Метод для удаления объявления по его ID
func deleteAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody deleteAdRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		v := c.Param("ad_id")
		adID, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		err = a.DeleteAd(c, int64(adID), reqBody.UserID)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.Status(http.StatusOK)
	}
}

// Метод для создания пользователя
func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		u, err := a.CreateUser(c, reqBody.Nickname, reqBody.Email)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(u))
	}
}

// Метода для обновления пользователя
func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateUserRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		userID := c.GetInt64("user_id")

		u, err := a.UpdateUser(c, userID, reqBody.Nickname, reqBody.Email)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(u))
	}
}

func showUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		v := c.Param("user_id")
		userID, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		u, err := a.UserByID(c, int64(userID))
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(u))
	}
}

// Метод для удаления пользователя по его ID
func deleteUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		v := c.Param("user_id")
		userID, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		err = a.DeleteUser(c, int64(userID))
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			}
			return
		}

		c.Status(http.StatusOK)
	}
}

// Метод для генерации шаблона для выборки объявлений
func createAdPattern(c *gin.Context, params listAdsRequest) (*ads.Pattern, error) {
	f := ads.NewPattern()

	f.Title = params.Title
	f.Created = params.Created
	if v, ok := c.GetQuery("user_id"); ok {
		id, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		f.UserID = int64(id)
	}
	if v, ok := c.GetQuery("published"); ok {
		p, err := strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}
		f.Published = p
	}

	return f, nil
}
