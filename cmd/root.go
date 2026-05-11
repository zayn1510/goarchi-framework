package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "goarchi",
	Short: color.HiCyanString("Goarchi CLI for generating Golang boilerplate code"),
	Long: color.HiWhiteString(`
%s

%s

%s

%s

%s

%s

%s

%s

%s

%s
`,
		color.New(color.FgHiBlue, color.Bold).Sprint("📦 Goarchi - Simple Layered Architecture Generator for Golang"),

		color.HiGreenString("🔧 Controller:\n  goarchi make controller [name]")+
			"\n    → Generate a controller (e.g. UserController)",

		color.HiGreenString("🛠️  Service:\n  goarchi make service [name]")+
			"\n    → Generate a service layer (e.g. UserService)",

		color.HiGreenString("📝 Request:\n  goarchi make request [name] [fields...]")+
			"\n    → Generate a request struct with validation (e.g. name:string age:int)",

		color.HiGreenString("📦 Resource:\n  goarchi make resource [name]")+
			"\n    → Generate a response formatter (DTO/transformer)",

		color.HiGreenString("🧩 Model:\n  goarchi make model [name] [fields...]")+
			"\n    → Generate a GORM model with tags\n    → Example: goarchi make model users \"id:int;primaryKey\" \"name:string;not null\"",

		color.HiGreenString("🗃️  Migration:\n  goarchi make migration [name]")+
			"\n    → Generate a migration file in 'database/migrations'",

		color.HiGreenString("🧬 Migrate:\n  goarchi make migrate [up|down]")+
			"\n    → 'up' applies migrations, 'down' rolls them back",

		color.HiCyanString("🐳 Docker Setup:\n  goarchi docker:init")+
			"\n    → Interactive wizard to generate Dockerfile.dev, docker-compose.yml, nginx config, and install.sh",

		color.HiYellowString("📌 Installation (Linux/macOS/Windows):")+
			"\n  go run cli/main.go build"+
			"\n  sudo mv goarchi /usr/local/bin/goarchi"+
			"\n  → Then use 'goarchi' globally from any folder",
	),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}