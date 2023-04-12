package httpgin

import (
	"github.com/gin-gonic/gin"

	"homework8/internal/app"
)

func AppRouter(r gin.IRouter, a app.App) {
	r.POST("/api/v1/ads", createAd(a))                    // Метод для создания объявления (ad)
	r.PUT("/api/v1/ads/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("api/v1/ads/:ad_id", updateAd(a))               // Метод для обновления текста(Text) или заголовка(Title) объявления
	r.GET("/api/v1/ads", listAds(a))                      // Метод для получения всех опубликованных объявлений
}
