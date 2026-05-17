package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zayn1510/goarchi/config"
)

var dbcheckCmd = &cobra.Command{
	Use:   "dbcheck",
	Short: "Check database connection",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking database connection...")

		db := config.GetDB()

		sqlDB, err := db.DB()
		if err != nil {
			fmt.Println("Failed to get database instance:", err)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			fmt.Println("Database connection failed:", err)
			return
		}

		fmt.Println("Database connection successful!")
	},
}

func init() {
	rootCmd.AddCommand(dbcheckCmd)
}
