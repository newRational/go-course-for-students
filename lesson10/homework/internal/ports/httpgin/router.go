package httpgin

import (
	"github.com/gin-gonic/gin"

	"homework10/internal/app"
)

func AppRouter(r gin.IRouter, a app.App) {
	g := r.Group("/api/v1")

	users := g.Group("/users")
	{
		users.POST("", createUser(a))
		users.PUT("/:user_id", updateUser(a))
		users.GET("/:user_id", showUser(a))
		users.DELETE("/:user_id", deleteUser(a))
	}

	ads := g.Group("/ads")
	{
		ads.GET("", listAds(a))
		ads.POST("", createAd(a))
		ads.GET("/:ad_id", showAd(a))
		ads.DELETE("/:ad_id", deleteAd(a))
		ads.PUT("/:ad_id", updateAd(a))
		ads.PUT("/:ad_id/status", changeAdStatus(a))
	}
}
