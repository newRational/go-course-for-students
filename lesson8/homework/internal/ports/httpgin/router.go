package httpgin

import (
	"github.com/gin-gonic/gin"

	"homework8/internal/app"
)

func AppRouter(r gin.IRouter, a app.App) {
	g := r.Group("/api/v1")

	users := g.Group("/users")
	{
		users.POST("", createUser(a))
		users.PUT("/:user_id", updateUser(a))
	}

	ads := g.Group("/ads")
	{
		ads.GET("", listAds(a))
		ads.POST("", createAd(a))
		ads.GET("/:ad_id", showAd(a))
		ads.PUT("/:ad_id", updateAd(a))
		ads.PUT("/:ad_id/status", changeAdStatus(a))
	}
}
