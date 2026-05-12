package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var govulncheckCmd = &cobra.Command{
	Use:   "govulncheck",
	Short: "Check for known vulnerabilities in dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("")
		color.Cyan("[INFO] Checking for vulnerabilities in dependencies...")
		fmt.Println("")

		// cek apakah govulncheck sudah terinstall
		if _, err := exec.LookPath("govulncheck"); err != nil {
			color.Yellow("[WARN] govulncheck is not installed.")
			fmt.Println("Installing govulncheck...")
			fmt.Println("")

			install := exec.Command("go", "install", "golang.org/x/vuln/cmd/govulncheck@latest")
			install.Stdout = os.Stdout
			install.Stderr = os.Stderr
			if err := install.Run(); err != nil {
				color.Red("❌ Failed to install govulncheck: %v", err)
				return
			}
			color.Green("[OK] govulncheck installed successfully!")
			fmt.Println("")
		}

		// jalankan govulncheck
		check := exec.Command("govulncheck", "./...")
		check.Stdout = os.Stdout
		check.Stderr = os.Stderr

		if err := check.Run(); err != nil {
			fmt.Println("")
			color.Red("❌ Vulnerabilities found! Fix them by running:")
			fmt.Println("   go get -u ./...")
			fmt.Println("   go mod tidy")
			fmt.Println("")
			return
		}

		fmt.Println("")
		color.Green("✅ No vulnerabilities found!")
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(govulncheckCmd)
}
