package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zayn1510/goarchi/config"
	"github.com/zayn1510/goarchi/database/migrations"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [direction]",
	Short: "Run database migrations (up or down)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		direction := args[0]

		if direction != "up" && direction != "down" {
			fmt.Println("Invalid direction. Use 'up' or 'down'.")
			os.Exit(1)
		}

		db := config.GetDB()

		if len(migrations.AllMigrations) == 0 {
			fmt.Println("No migrations found.")
			return
		}

		for _, migration := range migrations.AllMigrations {
			fmt.Println("Running migration:", migration.Name)

			var err error
			if direction == "up" {
				err = migration.Up(db)
			} else {
				err = migration.Down(db)
			}

			if err != nil {
				fmt.Println("❌ Failed:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Done:", migration.Name)
		}

		fmt.Println("\n✅ All migrations completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
