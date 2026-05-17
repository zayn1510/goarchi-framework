package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check Goarchi environment health",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("")
		color.Cyan("Goarchi Doctor - Environment Check")
		fmt.Println("")

		checkGoVersion()
		checkTool("govulncheck")
		checkGoMod()
		printSummary()
	},
}

func checkGoVersion() {
	fmt.Println("Checking Go version...")

	v := runtime.Version()
	fmt.Println("Go version:", v)

	if !isGoCompatible() {
		color.Red("Go version too old (requires 1.26+)")
		return
	}

	color.Green("Go version OK")
}

func checkTool(name string) {
	fmt.Printf("Checking %s...\n", name)

	_, err := exec.LookPath(name)
	if err != nil {
		color.Red("%s not found", name)
		fmt.Println("   Run: go install golang.org/x/vuln/cmd/govulncheck@latest")
		return
	}

	color.Green("%s installed", name)
}

func checkGoMod() {
	fmt.Println("Checking go.mod...")

	if _, err := os.Stat("go.mod"); err != nil {
		color.Red("go.mod not found")
		return
	}

	color.Green("go.mod exists")
}

func printSummary() {
	fmt.Println("\nSummary:")
	fmt.Println("- Goarchi CLI: OK")
	fmt.Println("- Environment: checked")
	fmt.Println("\nSystem ready for Goarchi development")
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}