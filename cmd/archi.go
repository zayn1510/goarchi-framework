package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zayn1510/goarchi/config"
	"github.com/zayn1510/goarchi/core/tools"
	"github.com/zayn1510/goarchi/database/migrations"
	"os"
	"strings"
	"time"
)

var makeControllerCmd = &cobra.Command{
	Use:   "controller [name]",
	Short: "Generate a new controller",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := strings.ToLower(args[0])
		parts := strings.Split(path, "/")
		name := parts[len(parts)-1]

		structName := strings.Title(name) + "Controller"
		filePath := fmt.Sprintf("app/controllers/%s_controller.go", path)
		buf, err := tools.GenerateController(structName)
		if err != nil {
			fmt.Println("Failed to execute template:", err)
			return
		}

		dir := "app/controllers/" + strings.Join(parts[:len(parts)-1], "/")
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println("Failed to create folder:", err)
			return
		}
		if err := os.WriteFile(filePath, []byte(buf), 0644); err != nil {
			fmt.Println("Failed to create controller:", err)
			return
		}

		fmt.Printf("%s\n  → %s\n",
			color.HiGreenString("✅ Controller created successfully!"),
			color.YellowString(filePath),
		)
	},
}

var makeServiceCmd = &cobra.Command{
	Use:   "service [name]",
	Short: "Generate a new service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := strings.ToLower(args[0])
		parts := strings.Split(path, "/")
		name := parts[len(parts)-1]

		structName := strings.Title(name) + "Service"
		filePath := fmt.Sprintf("app/services/%s_service.go", path)
		content, err := tools.GenerateServices(structName)
		if err != nil {
			fmt.Println("Failed to execute template:", err)
			return
		}

		dir := "app/services/" + strings.Join(parts[:len(parts)-1], "/")
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println("Failed to create folder:", err)
			return
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Println("Failed to create service:", err)
			return
		}

		fmt.Printf("%s\n  → %s\n",
			color.HiGreenString("✅ Service created successfully!"),
			color.YellowString(filePath),
		)
	},
}

var makeRequestCmd = &cobra.Command{
	Use:   "request [name] [fields]",
	Short: "Generate a new request with optional fields",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := strings.ToLower(args[0])
		parts := strings.Split(path, "/")
		name := parts[len(parts)-1]

		structName := strings.Title(name) + "Request"
		filePath := fmt.Sprintf("app/requests/%s_request.go", path)

		var fieldsBuilder strings.Builder
		for _, fieldArg := range args[1:] {
			parts := strings.Split(fieldArg, ":")
			if len(parts) != 2 {
				fmt.Printf("Invalid field '%s'. Use format name:type\n", fieldArg)
				return
			}
			fieldName := strings.Title(parts[0])
			fieldType := parts[1]
			fieldsBuilder.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\" validate:\"required\"`\n", fieldName, fieldType, parts[0]))
		}

		content, err := tools.GenerateRequest(structName, fieldsBuilder)
		if err != nil {
			fmt.Println("Failed to execute template:", err)
			return
		}

		dir := "app/requests/" + strings.Join(parts[:len(parts)-1], "/")
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println("Failed to create folder:", err)
			return
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Println("Failed to create request:", err)
			return
		}

		fmt.Printf("%s\n  → %s\n",
			color.HiGreenString("✅ Request created successfully!"),
			color.YellowString(filePath),
		)
	},
}

var makeResourceCmd = &cobra.Command{
	Use:   "resource [name] [fields]",
	Short: "Generate a new resource with optional fields",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := strings.ToLower(args[0])
		parts := strings.Split(path, "/")
		name := parts[len(parts)-1]

		structName := strings.Title(name) + "Resource"
		filePath := fmt.Sprintf("app/resources/%s_resource.go", path)

		var fieldsBuilder strings.Builder
		for _, fieldArg := range args[1:] {
			parts := strings.Split(fieldArg, ":")
			if len(parts) != 2 {
				fmt.Printf("Invalid field '%s'. Use format name:type\n", fieldArg)
				return
			}
			fieldName := strings.Title(parts[0])
			fieldType := parts[1]
			fieldsBuilder.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldType, parts[0]))
		}

		content, err := tools.GenerateResource(structName, fieldsBuilder)
		if err != nil {
			fmt.Println("Failed to execute template:", err)
			return
		}

		dir := "app/resources/" + strings.Join(parts[:len(parts)-1], "/")
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println("Failed to create folder:", err)
			return
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Println("Failed to create resource:", err)
			return
		}

		fmt.Printf("%s\n  → %s\n",
			color.HiGreenString("✅ Resource created successfully!"),
			color.YellowString(filePath),
		)
	},
}

var makeModelCmd = &cobra.Command{
	Use:   "model [path] [fields]",
	Short: "Generate a new model with fields and GORM tags",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := strings.ToLower(args[0])
		parts := strings.Split(path, "/")
		name := parts[len(parts)-1]
		structName := strings.Title(strings.TrimSuffix(name, "s"))

		dir := "app/models/" + strings.Join(parts[:len(parts)-1], "/")
		filePath := fmt.Sprintf("%s/%s.go", dir, name)

		var fieldsBuilder strings.Builder
		for _, fieldArg := range args[1:] {
			parts := strings.Split(fieldArg, ":")
			if len(parts) < 2 {
				fmt.Printf("Invalid field '%s'. Use format name:type;tag1;tag2 or Struct:foreignKey:Field\n", fieldArg)
				return
			}

			if strings.ToUpper(parts[0][:1]) == parts[0][:1] && parts[1] == "foreignKey" {
				structName := parts[0]
				foreignKey := parts[2]
				fieldsBuilder.WriteString(fmt.Sprintf("\t%s %s `gorm:\"foreignKey:%s\"`\n", structName, structName, foreignKey))
				continue
			}

			fieldName := strings.Title(parts[0])
			tagParts := strings.Split(parts[1], ";")
			fieldType := tagParts[0]
			gormTags := strings.Join(tagParts[1:], ";")
			jsonTag := strings.ToLower(parts[0])

			var tagBuilder strings.Builder
			if gormTags != "" {
				tagBuilder.WriteString(fmt.Sprintf("gorm:\"%s\" ", gormTags))
			}
			tagBuilder.WriteString(fmt.Sprintf("json:\"%s\"", jsonTag))

			fieldsBuilder.WriteString(fmt.Sprintf("\t%s %s `%s`\n", fieldName, fieldType, tagBuilder.String()))
		}
		fieldsBuilder.WriteString("\tCreatedAt time.Time `json:\"created_at\"`\n")
		fieldsBuilder.WriteString("\tUpdatedAt time.Time `json:\"updated_at\"`\n")

		content, err := tools.GenerateModel(structName, fieldsBuilder)
		if err != nil {
			fmt.Println("Failed to create model:", err)
			return
		}
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println("Failed to create folder:", err)
			return
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Println("Failed to create model:", err)
			return
		}

		fmt.Printf("%s\n  → %s\n",
			color.HiGreenString("✅ Model created successfully!"),
			color.YellowString(filePath),
		)
	},
}

var makeMigrationCmd = &cobra.Command{
	Use:   "migration [name]",
	Short: "Generate a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		timestamp := time.Now().Format("20060102150405")
		fileName := fmt.Sprintf("database/migrations/%s_%s.go", timestamp, name)

		content, err := tools.GenerateMigration()
		if err != nil {
			fmt.Println("Failed to execute template:", err)
			return
		}
		if err := os.MkdirAll("database/migrations", os.ModePerm); err != nil {
			fmt.Println("Failed to create folder:", err)
			return
		}
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			fmt.Println("Failed to create migration:", err)
			return
		}

		fmt.Printf("%s\n  → %s\n",
			color.HiGreenString("✅ Migration created successfully!"),
			color.YellowString(fileName),
		)
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate [direction]",
	Short: "Run migrations (up or down)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		direction := args[0]

		if direction != "up" && direction != "down" {
			fmt.Println("Invalid direction. Use 'up' or 'down'.")
			return
		}

		db := config.GetDB()

		for _, migration := range migrations.AllMigrations {
			fmt.Println("Running migration:", migration.Name)

			var migErr error
			if direction == "up" {
				migErr = migration.Up(db)
			} else {
				migErr = migration.Down(db)
			}

			if migErr != nil {
				fmt.Println("Failed:", migErr)
				return
			}
			fmt.Println("Done:", migration.Name)
		}
	},
}

func init() {
	makeCmd.AddCommand(makeControllerCmd)
	makeCmd.AddCommand(makeServiceCmd)
	makeCmd.AddCommand(makeRequestCmd)
	makeCmd.AddCommand(makeResourceCmd)
	makeCmd.AddCommand(makeModelCmd)
	makeCmd.AddCommand(makeMigrationCmd)
	makeCmd.AddCommand(migrateCmd)
}