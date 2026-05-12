package main

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/zayn1510/goarchi/app/middleware"
	"github.com/zayn1510/goarchi/console"
	"github.com/zayn1510/goarchi/routers"
)

func main() {
	if len(os.Args) > 1 {
		console.Execute()
		return
	}

	c := gin.Default()
	middleware.SetCors(c)
	routers.RegisterRoutes(c)
	if err := c.Run(":8080"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
