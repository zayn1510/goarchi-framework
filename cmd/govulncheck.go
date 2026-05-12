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

		// selalu reinstall govulncheck biar sesuai versi Go yang aktif
		color.Cyan("[INFO] Installing/updating govulncheck...")
		install := exec.Command("go", "install", "golang.org/x/vuln/cmd/govulncheck@latest")
		install.Stdout = os.Stdout
		install.Stderr = os.Stderr
		if err := install.Run(); err != nil {
			color.Red("❌ Failed to install govulncheck: %v", err)
			return
		}
		color.Green("[OK] govulncheck ready.")
		fmt.Println("")

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