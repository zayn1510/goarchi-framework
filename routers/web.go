package routers

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	// api.Use(middleware.RateLimit())
	setUpRouterPing(api)
	
}

func setUpRouterPing(router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping",
		})
	})
	router.GET("welcome", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "welcome to api",
		})
	})
	router.GET("bye", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "bye bye",
		})
	})
}
