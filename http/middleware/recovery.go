package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Recovery() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if re := recover(); re != nil {
				log.Error(re)
				context.JSON(500, gin.H{"message": re, "status": false})
			}
		}()
		context.Next()
	}
}
