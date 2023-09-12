package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/goasali/toolkit/http/controllers"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if re := recover(); re != nil {
				log.Error(re)
				res := controllers.Response{}
				if err, ok := re.(error); ok {
					res.HandleError(c, err)
				} else {
					res.SendMessage(c, http.StatusInternalServerError, "errors.internal_server")
				}
			}
		}()
		c.Next()
	}
}
