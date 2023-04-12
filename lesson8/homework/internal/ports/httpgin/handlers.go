package httpgin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"homework8/internal/app"
)

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
		}

		ad, err := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, AdErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, AdErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
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
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
		}

		adID := c.GetInt("ad_id")

		ad, err := a.ChangeAdStatus(c, int64(adID), reqBody.UserID, reqBody.Published)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, AdErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, AdErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
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
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
		}

		adID := c.GetInt("ad_id")

		ad, err := a.UpdateAd(c, int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, AdErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, AdErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func listAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ads, err := a.ListAds(c)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, AdErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, AdErrorResponse(err))
			} else {
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}

		c.JSON(http.StatusOK, AdsSuccessResponse(ads))
	}
}
