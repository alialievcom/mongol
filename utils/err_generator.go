package utils

import "github.com/gin-gonic/gin"

type response struct {
	Error string `json:"error" example:"message"`
}

func ErrorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response{msg})
}
