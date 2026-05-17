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

		// cek binary govulncheck
		path, err := exec.LookPath("govulncheck")
		if err != nil {
			color.Red("govulncheck not installed")

			color.Cyan("[INFO] Installing govulncheck with Go 1.26 toolchain...")

			install := exec.Command(
				"go",
				"install",
				"golang.org/x/vuln/cmd/govulncheck@latest",
			)

			// paksa pakai toolchain lokal (NO auto-switch)
			install.Env = append(os.Environ(),
				"GOTOOLCHAIN=local",
				"CGO_ENABLED=1",
			)

			install.Stdout = os.Stdout
			install.Stderr = os.Stderr

			if err := install.Run(); err != nil {
				color.Red("Failed to install govulncheck: %v", err)
				return
			}

			color.Green("[OK] govulncheck installed.")
			fmt.Println("")
		} else {
			color.Green("[OK] govulncheck found at: %s", path)
		}

		fmt.Println("")

		// jalankan tanpa reinstall, tanpa noise toolchain switch
		check := exec.Command("govulncheck", "./...")

		check.Env = append(os.Environ(),
			"GOTOOLCHAIN=local",
			"CGO_ENABLED=1",
		)

		check.Stdout = os.Stdout
		check.Stderr = os.Stderr

		if err := check.Run(); err != nil {
			fmt.Println("")
			color.Red("govulncheck detected issues (or tool mismatch)")
			fmt.Println("   If this looks like a toolchain error, rebuild govulncheck:")
			fmt.Println("   go install golang.org/x/vuln/cmd/govulncheck@latest")
			fmt.Println("")
			return
		}

		fmt.Println("")
		color.Green("No vulnerabilities found!")
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(govulncheckCmd)
}
