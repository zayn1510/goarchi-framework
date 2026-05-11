package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func RunInstall() {

	binaryName := "goarchi"

	// Windows binary extension
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	fmt.Println("📦 Building Goarchi binary...")

	// Build binary
	cmd := exec.Command("go", "build", "-o", binaryName, "cli/main.go")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Failed to build binary:", err)
		return
	}

	fmt.Println("\n✅ Binary generated successfully!")

	// =========================
	// INSTALLATION GUIDE
	// =========================

	switch runtime.GOOS {

	case "windows":

		fmt.Println("\n📌 Windows Installation")
		fmt.Println("Move goarchi.exe to a folder inside your PATH.")
		fmt.Println("\nExample:")
		fmt.Println("C:\\Users\\YourUser\\go\\bin")

		fmt.Println("\nThen verify:")
		fmt.Println("goarchi version")

	case "darwin":

		fmt.Println("\n📌 macOS Installation")
		fmt.Println("Run the following command:")

		fmt.Printf("\nsudo mv %s /usr/local/bin/goarchi\n", binaryName)

		fmt.Println("\nThen verify:")
		fmt.Println("goarchi version")

	case "linux":

		fmt.Println("\n📌 Linux Installation")
		fmt.Println("Run the following command:")

		fmt.Printf("\nsudo mv %s /usr/local/bin/goarchi\n", binaryName)

		fmt.Println("\nThen verify:")
		fmt.Println("goarchi version")

	default:

		fmt.Println("\n⚠️ Unsupported OS detected.")
		fmt.Println("Please move the binary manually to your system PATH.")
	}
}