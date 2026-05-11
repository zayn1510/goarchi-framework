package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type AppConfig struct {
	Version string `yaml:"version"`
}

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Generate code components like controller, model, etc.",
}
var airCmd = &cobra.Command{
	Use:   "air",
	Short: "Run Goarchi in development mode with Air (hot reload)",
	Run: func(cmd *cobra.Command, args []string) {
		RunDev()
	},
}
var installCmd = &cobra.Command{
	Use:   "build",
	Short: "Install Goarchi CLI globally",
	Run: func(cmd *cobra.Command, args []string) {
		RunInstall()
	},
}
var installDBCommand = &cobra.Command{
	Use:   "install-db",
	Short: "Install Database Clients (GORM, sqlx, PostgreSQL, MongoDB)",
	Long: `Install various database clients including:
- GORM + selected driver (MySQL, PostgreSQL, etc.)
- sqlx
- PostgreSQL driver
- MongoDB driver`,
	Run: func(cmd *cobra.Command, args []string) {
		options := []string{
			"GORM (gorm.io/gorm)",
			"sqlx (github.com/jmoiron/sqlx)",
			"PostgreSQL Driver (github.com/lib/pq)",
			"MongoDB Driver (go.mongodb.org/mongo-driver)",
		}

		packages := map[int]string{
			1: "gorm.io/gorm",
			2: "github.com/jmoiron/sqlx",
			3: "github.com/lib/pq",
			4: "go.mongodb.org/mongo-driver",
		}

		reader := bufio.NewReader(os.Stdin)

		// Select client
		var choice int
		for {
			fmt.Println("Choose a database client to install:")
			for i, opt := range options {
				fmt.Printf("  %d. %s\n", i+1, opt)
			}
			fmt.Print("\nEnter your choice (1-4): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			num, err := strconv.Atoi(input)
			if err == nil && num >= 1 && num <= len(options) {
				choice = num
				break
			}
			fmt.Println("❌ Invalid choice. Please enter a number between 1 and 4.\n")
		}

		// Handle GORM special case
		if choice == 1 {
			fmt.Println("📦 Installing GORM core (gorm.io/gorm)...")
			cmdGorm := exec.Command("go", "get", "gorm.io/gorm")
			cmdGorm.Stdout = os.Stdout
			cmdGorm.Stderr = os.Stderr
			if err := cmdGorm.Run(); err != nil {
				fmt.Println("❌ Failed to install GORM:", err)
				return
			}

			// Driver selection for GORM
			fmt.Println("Choose GORM driver to install:")
			drivers := map[int]string{
				1: "gorm.io/driver/mysql",
				2: "gorm.io/driver/postgres",
				3: "gorm.io/driver/sqlite",
				4: "gorm.io/driver/sqlserver",
			}
			driverOptions := []string{
				"MySQL (gorm.io/driver/mysql)",
				"PostgreSQL (gorm.io/driver/postgres)",
				"SQLite (gorm.io/driver/sqlite)",
				"SQL Server (gorm.io/driver/sqlserver)",
			}

			var driverChoice int
			for {
				for i, opt := range driverOptions {
					fmt.Printf("  %d. %s\n", i+1, opt)
				}
				fmt.Print("\nEnter your choice (1-4): ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				num, err := strconv.Atoi(input)
				if err == nil && num >= 1 && num <= len(drivers) {
					driverChoice = num
					break
				}
				fmt.Println("❌ Invalid driver. Please enter a valid number.\n")
			}

			driver := drivers[driverChoice]
			fmt.Printf("📦 Installing GORM driver: %s\n", driver)
			cmdDriver := exec.Command("go", "get", driver)
			cmdDriver.Stdout = os.Stdout
			cmdDriver.Stderr = os.Stderr
			if err := cmdDriver.Run(); err != nil {
				fmt.Printf("❌ Failed to install %s: %s\n", driver, err)
			} else {
				fmt.Printf("✅ Successfully installed %s\n", driver)
			}
			return
		}

		// Install non-GORM package
		packagePath := packages[choice]
		fmt.Printf("📦 Installing %s ...\n", packagePath)
		cmdInstall := exec.Command("go", "get", packagePath)
		cmdInstall.Stdout = os.Stdout
		cmdInstall.Stderr = os.Stderr

		if err := cmdInstall.Run(); err != nil {
			fmt.Printf("❌ Failed to install %s: %s\n", packagePath, err)
		} else {
			fmt.Printf("✅ Successfully installed %s\n", packagePath)
		}
	},
}

var installRouter = &cobra.Command{
	Use:   "router",
	Short: "Install a Go HTTP router framework",
	Long: `Install a popular Go web framework for building RESTful APIs or web applications.

You can choose one of the available router frameworks:
  1. Gin   - A high-performance HTTP web framework
  2. Fiber - An Express.js inspired web framework built on Fasthttp
  3. Echo  - A minimalist and extensible web framework

Simply select a framework by entering the corresponding number.`,
	Run: func(cmd *cobra.Command, args []string) {
		options := []string{
			"Gin (github.com/gin-gonic/gin)",
			"Fiber (github.com/gofiber/fiber/v2)",
			"Echo (github.com/labstack/echo/v4)",
		}

		packages := map[int]string{
			1: "github.com/gin-gonic/gin",
			2: "github.com/gofiber/fiber/v2",
			3: "github.com/labstack/echo/v4",
		}

		fmt.Println("Choose a framework to install:")
		for i, opt := range options {
			fmt.Printf("  %d. %s\n", i+1, opt)
		}
		fmt.Print("\nEnter your choice (1-3): ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		var choice int
		for {
			fmt.Println("Choose a framework to install:")
			for i, opt := range options {
				fmt.Printf("  %d. %s\n", i+1, opt)
			}
			fmt.Print("\nEnter your choice (1-3): ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			num, err := strconv.Atoi(input)
			if err == nil && num >= 1 && num <= len(options) {
				choice = num
				break
			}
			fmt.Println("❌ Invalid choice. Please enter a number between 1 and 3.\n")
		}
		packagePath := packages[choice]
		fmt.Printf("📦 Installing %s ...\n", packagePath)

		cmdInstall := exec.Command("go", "get", packagePath)
		cmdInstall.Stdout = os.Stdout
		cmdInstall.Stderr = os.Stderr

		if err := cmdInstall.Run(); err != nil {
			fmt.Printf("❌ Failed to install %s: %s\n", packagePath, err)
		} else {
			fmt.Printf("✅ Successfully installed %s\n", packagePath)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version of the Goarchi Framework CLI tool",
	Long:  `Displays the current version of the Goarchi Framework CLI tool as defined in the configuration file. This command helps you keep track of the version you are using and check for any updates or new releases.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := readConfig("config.yaml")
		if err != nil {
			fmt.Println("Error reading config:", err)
			return
		}
		// Displaying version
		fmt.Printf("Goarchi Framework CLI Tool - Version: %s\n", config.Version)
	},
}
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean unused Go modules",
	Run: func(cmd *cobra.Command, args []string) {
		tidy := exec.Command("go", "mod", "tidy")
		tidy.Stdout = os.Stdout
		tidy.Stderr = os.Stderr
		if err := tidy.Run(); err != nil {
			fmt.Println("❌ Failed to tidy modules:", err)
			return
		}
		fmt.Println("✅ Unused packages removed!")
	},
}

func readConfig(filename string) (*AppConfig, error) {
	// Membaca isi file YAML
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parsing YAML
	var config AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(makeCmd)
	rootCmd.AddCommand(airCmd)
	rootCmd.AddCommand(installRouter)
	rootCmd.AddCommand(installDBCommand)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(cleanCmd)
}
