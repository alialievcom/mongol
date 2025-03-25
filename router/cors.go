package router

import (
	"github.com/AliAlievMos/mongol/models"
	"github.com/gin-gonic/gin"
)

func corsMiddleware(cfg models.Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.Origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", cfg.Headers)
		c.Writer.Header().Set("Access-Control-Allow-Methods", cfg.Methods)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
