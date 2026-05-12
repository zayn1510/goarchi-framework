package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zayn1510/goarchi/app/middleware"
	"github.com/zayn1510/goarchi/config"
	"github.com/zayn1510/goarchi/database/migrations"
	"github.com/zayn1510/goarchi/routers"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			if len(os.Args) < 3 {
				fmt.Println("Usage: go run main.go migrate [up|down]")
				os.Exit(1)
			}
			runMigrate(os.Args[2])
			return
		}
	}

	c := gin.Default()
	middleware.SetCors(c)
	routers.RegisterRoutes(c)
	err := c.Run(":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func runMigrate(direction string) {
	if direction != "up" && direction != "down" {
		fmt.Println("Invalid direction. Use 'up' or 'down'.")
		os.Exit(1)
	}

	db := config.GetDB()

	for _, migration := range migrations.AllMigrations {
		fmt.Println("Running migration:", migration.Name)

		var err error
		if direction == "up" {
			err = migration.Up(db)
		} else {
			err = migration.Down(db)
		}

		if err != nil {
			fmt.Println("Failed:", err)
			os.Exit(1)
		}
		fmt.Println("Done:", migration.Name)
	}

	fmt.Println("✅ Migration completed successfully!")
}