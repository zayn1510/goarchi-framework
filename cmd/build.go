package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func RunInstall() {

	if !isGoCompatible() {
		fmt.Println("Go 1.26+ required to build Goarchi")
		fmt.Println("Current:", runtime.Version())
		return
	}

	binaryName := "goarchi"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	fmt.Println("Building Goarchi binary...")

	cmd := exec.Command("go", "build", "-ldflags",
		"-X 'main.Version=dev'", "-o", binaryName, "cli/main.go")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to build binary:", err)
		return
	}

	fmt.Println("\nBinary generated successfully!")
	fmt.Println("Go version used:", runtime.Version())

	printInstallGuide(binaryName)
}

func isGoCompatible() bool {
	version := runtime.Version()
	if !strings.HasPrefix(version, "go") {
		return false
	}

	v := strings.TrimPrefix(version, "go")

	parts := strings.Split(v, ".")
	if len(parts) < 2 {
		return false
	}

	majorMinor := parts[0] + "." + parts[1]

	f, err := strconv.ParseFloat(majorMinor, 64)
	if err != nil {
		return false
	}

	// minimum Go 1.26
	return f >= 1.26
}
func printInstallGuide(binaryName string) {
	switch runtime.GOOS {

	case "windows":
		fmt.Println("\nWindows Installation")
		fmt.Println("Move binary to PATH (e.g. C:\\Users\\<user>\\go\\bin)")

	case "darwin", "linux":
		fmt.Println("\nUnix Installation")
		fmt.Printf("sudo mv %s /usr/local/bin/goarchi\n", binaryName)

	default:
		fmt.Println("\nManual install required")
	}
}