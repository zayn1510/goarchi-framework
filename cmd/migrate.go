package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zayn1510/goarchi/config"
	"github.com/zayn1510/goarchi/database/migrations"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [direction] [name?]",
	Short: "Run database migrations (up or down)",
	Args:  cobra.RangeArgs(1, 2),
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

		target := ""
		if len(args) == 2 {
			target = args[1]
		}

		found := false
		for _, migration := range migrations.AllMigrations {
			if target != "" && migration.Name != target {
				continue
			}

			found = true

			var err error
			if direction == "up" {
				fmt.Println("⬆Migrating:", migration.Name)
				err = migration.Up(db)
				if err != nil {
					fmt.Printf("Failed [%s]: %v\n", migration.Name, err)
					os.Exit(1)
				}
				fmt.Println("Migrated:", migration.Name)
			} else {
				fmt.Println("⬇Dropping:", migration.Name)
				err = migration.Down(db)
				if err != nil {
					fmt.Printf("Failed [%s]: %v\n", migration.Name, err)
					os.Exit(1)
				}
				fmt.Println("Dropped:", migration.Name)
			}
		}

		if target != "" && !found {
			fmt.Printf("Migration '%s' not found.\n", target)
			fmt.Println("\nAvailable migrations:")
			for _, m := range migrations.AllMigrations {
				fmt.Println(" -", m.Name)
			}
			os.Exit(1)
		}

		if direction == "up" {
			fmt.Println("\nAll migrations applied successfully!")
		} else {
			fmt.Println("\nAll tables dropped successfully!")
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
