package httpgin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	green  = "\033[97;42m"
	yellow = "\033[90;43m"
	red    = "\033[97;41m"
	cyan   = "\033[97;46m"
	reset  = "\033[0m"
)

func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log.SetPrefix("[ADAPP] - ")
		code := fmt.Sprintf("<"+StatusCodeColor(c.Writer.Status())+"%d"+reset+">", c.Writer.Status())
		path := c.Request.URL.Path
		log.Print("- ", code, " - "+c.Request.Method+":\t", path)
	}
}

func StatusCodeColor(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return cyan
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}
