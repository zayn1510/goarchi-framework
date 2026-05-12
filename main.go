package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zayn1510/goarchi/app/middleware"
	"github.com/zayn1510/goarchi/routers"
)

func main() {
	c := gin.Default()
	middleware.SetCors(c)
	routers.RegisterRoutes(c)
	if err := c.Run(":8080"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
