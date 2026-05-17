package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zayn1510/goarchi/config"
	"github.com/zayn1510/goarchi/database/seeders"
)

var seedCmd = &cobra.Command{
	Use:   "seed [direction] [name?]",
	Short: "Run database seeders (up or down)",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		direction := args[0]

		if direction != "up" && direction != "down" {
			fmt.Println("Invalid direction. Use 'up' or 'down'.")
			os.Exit(1)
		}

		db := config.GetDB()

		if len(seeders.AllSeeders) == 0 {
			fmt.Println("No seeders found.")
			return
		}

		target := ""
		if len(args) == 2 {
			target = args[1]
		}

		found := false
		for _, seeder := range seeders.AllSeeders {
			if target != "" && seeder.Name != target {
				continue
			}

			found = true

			var err error
			if direction == "up" {
				fmt.Println("⬆Seeding:", seeder.Name)
				err = seeder.Up(db)
				if err != nil {
					fmt.Printf("Failed [%s]: %v\n", seeder.Name, err)
					os.Exit(1)
				}
				fmt.Println("Seeded:", seeder.Name)
			} else {
				fmt.Println("⬇Dropping seed:", seeder.Name)
				err = seeder.Down(db)
				if err != nil {
					fmt.Printf("Failed [%s]: %v\n", seeder.Name, err)
					os.Exit(1)
				}
				fmt.Println("Dropped seed:", seeder.Name)
			}
		}

		if target != "" && !found {
			fmt.Printf("Seeder '%s' not found.\n", target)
			fmt.Println("\nAvailable seeders:")
			for _, s := range seeders.AllSeeders {
				fmt.Println(" -", s.Name)
			}
			os.Exit(1)
		}

		if direction == "up" {
			fmt.Println("\nAll seeders applied successfully!")
		} else {
			fmt.Println("\nAll seeds dropped successfully!")
		}
	},
}