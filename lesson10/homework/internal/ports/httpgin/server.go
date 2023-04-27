package httpgin

import (
	"github.com/gin-contrib/cors"
	"net/http"

	"github.com/gin-gonic/gin"

	"homework10/internal/app"
)

func NewHTTPServer(port string, a app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	handler.Use(gin.Recovery(), logger(), cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type"},
	}))

	AppRouter(handler, a)
	s := &http.Server{Addr: port, Handler: handler}

	return s
}
