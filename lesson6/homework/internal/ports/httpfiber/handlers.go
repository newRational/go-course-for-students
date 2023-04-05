package httpfiber

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"homework6/internal/app"
)

// Метод для создания объявления (ad)
func createAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody createAdRequest
		err := c.BodyParser(&reqBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, err := a.CreateAd(c.Context(), reqBody.Title, reqBody.Text, reqBody.UserID)
		if err != nil {
			code, e := strconv.Atoi(err.Error())
			if e != nil {
				c.Status(http.StatusInternalServerError)
			} else {
				c.Status(code)
			}
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody changeAdStatusRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, err := a.ChangeAdStatus(c.Context(), int64(adID), reqBody.UserID, reqBody.Published)

		if err != nil {
			code, e := strconv.Atoi(err.Error())
			if e != nil {
				c.Status(http.StatusInternalServerError)
			} else {
				c.Status(code)
			}
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody updateAdRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, err := a.UpdateAd(c.Context(), int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)

		if err != nil {
			code, e := strconv.Atoi(err.Error())
			if e != nil {
				c.Status(http.StatusInternalServerError)
			} else {
				c.Status(code)
			}
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}
