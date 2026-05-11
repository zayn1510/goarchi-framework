package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var dockerInitCmd = &cobra.Command{
	Use:   "docker:init",
	Short: "Generate Dockerfile and docker-compose for Goarchi project",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("🐳 Goarchi Docker Setup")
		fmt.Println("========================")

		// GO VERSION
		fmt.Println("Choose Go Version:")
		fmt.Println("1. 1.22")
		fmt.Println("2. 1.23")
		fmt.Println("3. latest")
		fmt.Print("Enter choice (1-3): ")
		goChoice, _ := reader.ReadString('\n')
		goChoice = strings.TrimSpace(goChoice)
		goVersion := "latest"
		switch goChoice {
		case "1":
			goVersion = "1.22"
		case "2":
			goVersion = "1.23"
		case "3":
			goVersion = "latest"
		}

		// DATABASE
		fmt.Println("\nChoose Database:")
		fmt.Println("1. PostgreSQL")
		fmt.Println("2. MySQL")
		fmt.Print("Enter choice (1-2): ")
		dbChoice, _ := reader.ReadString('\n')
		dbChoice = strings.TrimSpace(dbChoice)
		database := "postgres"
		if dbChoice == "2" {
			database = "mysql"
		}

		// AIR
		fmt.Print("\nUse Air Hot Reload? (y/n): ")
		airChoice, _ := reader.ReadString('\n')
		useAir := strings.TrimSpace(strings.ToLower(airChoice)) == "y"

		// NGINX
		fmt.Print("\nUse Nginx as reverse proxy? (y/n): ")
		nginxChoice, _ := reader.ReadString('\n')
		useNginx := strings.TrimSpace(strings.ToLower(nginxChoice)) == "y"

		// DB GUI
		dbGui := "none"
		if database == "mysql" {
			fmt.Println("\nChoose DB GUI:")
			fmt.Println("1. phpMyAdmin")
			fmt.Println("2. Adminer")
			fmt.Println("3. Skip")
			fmt.Print("Enter choice (1-3): ")
			guiChoice, _ := reader.ReadString('\n')
			switch strings.TrimSpace(guiChoice) {
			case "1":
				dbGui = "phpmyadmin"
			case "2":
				dbGui = "adminer"
			default:
				dbGui = "none"
			}
		} else {
			fmt.Println("\nChoose DB GUI:")
			fmt.Println("1. Adminer")
			fmt.Println("2. Skip")
			fmt.Print("Enter choice (1-2): ")
			guiChoice, _ := reader.ReadString('\n')
			if strings.TrimSpace(guiChoice) == "1" {
				dbGui = "adminer"
			}
		}

		// APP NAME
		fmt.Print("\nEnter your app name (used for container names, e.g. myapp): ")
		appName, _ := reader.ReadString('\n')
		appName = strings.TrimSpace(strings.ToLower(appName))
		if appName == "" {
			appName = "goarchi"
		}
		appName = strings.ReplaceAll(appName, " ", "_")
		appName = strings.ReplaceAll(appName, "-", "_")

		// GENERATE FILES
		if err := os.WriteFile("Dockerfile.dev", []byte(generateDockerfile(goVersion, useAir)), 0644); err != nil {
			fmt.Println("❌ Failed to generate Dockerfile.dev:", err)
			return
		}
		if err := os.WriteFile("docker-compose.yml", []byte(generateCompose(database, appName, useNginx, dbGui)), 0644); err != nil {
			fmt.Println("❌ Failed to generate docker-compose.yml:", err)
			return
		}
		if useNginx {
			if err := os.MkdirAll("nginx", os.ModePerm); err != nil {
				fmt.Println("❌ Failed to create nginx folder:", err)
				return
			}
			if err := os.WriteFile("nginx/default.conf", []byte(generateNginxConf()), 0644); err != nil {
				fmt.Println("❌ Failed to generate nginx/default.conf:", err)
				return
			}
		}
		if err := os.WriteFile("install.sh", []byte(generateInstallScript(database, appName, useNginx, dbGui)), 0755); err != nil {
			fmt.Println("❌ Failed to generate install.sh:", err)
			return
		}

		fmt.Println("\n✅ Docker configuration generated successfully!")
		fmt.Println("📄 Dockerfile.dev")
		fmt.Println("📄 docker-compose.yml")
		if useNginx {
			fmt.Println("📄 nginx/default.conf")
		}
		fmt.Println("📄 install.sh")
		fmt.Println("")
		fmt.Printf("🔖 Container names:\n")
		fmt.Printf("   App : %s_api\n", appName)
		fmt.Printf("   DB  : %s_db\n", appName)
		if useNginx {
			fmt.Printf("   Nginx : %s_nginx\n", appName)
		}
		if dbGui != "none" {
			fmt.Printf("   DB GUI (%s) : %s_gui\n", dbGui, appName)
		}
		fmt.Println("")
		fmt.Println("▶️  To deploy, run: bash install.sh")
	},
}

func generateDockerfile(goVersion string, useAir bool) string {
	cmd := `CMD ["go", "run", "main.go"]`
	if useAir {
		cmd = `CMD ["air"]`
	}
	airInstall := ""
	if useAir {
		airInstall = `RUN go install github.com/air-verse/air@latest`
	}
	return fmt.Sprintf(`FROM golang:%s

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

%s

EXPOSE 8080

%s
`, goVersion, airInstall, cmd)
}

func generateNginxConf() string {
	return `server {
    listen 80;

    location / {
        proxy_pass http://api:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
`
}

func generateCompose(database string, appName string, useNginx bool, dbGui string) string {
	var b strings.Builder

	b.WriteString("services:\n")

	// =========================
	// APP
	// =========================
	b.WriteString(fmt.Sprintf(`
  # =========================
  # GO APP
  # =========================
  api:
    build:
      context: .
      dockerfile: Dockerfile.dev
    container_name: %s_api
    ports:
      - "9034:8080"
    volumes:
      - .:/app
    restart: always
    depends_on:
      - db
    networks:
      - %s_network
    environment:
      - APP_ENV=docker
      - DB_CONNECTION=%s
      - DB_HOST=db
      - DB_NAME=${DB_NAME}
      - DB_USER=${DB_USER}
      - DB_PASS=${DB_PASS}
      - DB_PORT=${DB_PORT}
      - DB_PREFIX=${DB_PREFIX}
      - JWT_SECRET_KEY=${JWT_SECRET_KEY}
      - JWT_EXPIRED_TOKEN=${JWT_EXPIRED_TOKEN}
`, appName, appName, database))

	// =========================
	// NGINX
	// =========================
	if useNginx {
		b.WriteString(fmt.Sprintf(`
  # =========================
  # NGINX
  # =========================
  nginx:
    image: nginx:latest
    container_name: %s_nginx
    ports:
      - "9080:80"
    volumes:
      - ./nginx/default.conf:/etc/nginx/conf.d/default.conf:ro
    depends_on:
      - api
    networks:
      - %s_network
    restart: always
`, appName, appName))
	}

	// =========================
	// DATABASE
	// =========================
	if database == "mysql" {
		b.WriteString(fmt.Sprintf(`
  # =========================
  # MYSQL
  # =========================
  db:
    image: mysql:8
    container_name: %s_db
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME}
      MYSQL_USER: ${DB_USER}
      MYSQL_PASSWORD: ${DB_PASS}
    ports:
      - "9339:3306"
    volumes:
      - %s_db_data:/var/lib/mysql
    networks:
      - %s_network
`, appName, appName, appName))
	} else {
		b.WriteString(fmt.Sprintf(`
  # =========================
  # POSTGRESQL
  # =========================
  db:
    image: postgres:16
    container_name: %s_db
    restart: always
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASS}
      POSTGRES_DB: ${DB_NAME}
    ports:
      - "9432:5432"
    volumes:
      - %s_db_data:/var/lib/postgresql/data
    networks:
      - %s_network
`, appName, appName, appName))
	}

	// =========================
	// DB GUI
	// =========================
	if dbGui == "phpmyadmin" {
		b.WriteString(fmt.Sprintf(`
  # =========================
  # PHPMYADMIN
  # =========================
  gui:
    image: phpmyadmin/phpmyadmin
    container_name: %s_gui
    environment:
      - PMA_HOST=db
      - PMA_USER=${DB_USER}
      - PMA_PASSWORD=${DB_PASS}
    ports:
      - "9090:80"
    depends_on:
      - db
    networks:
      - %s_network
    restart: always
`, appName, appName))
	} else if dbGui == "adminer" {
		b.WriteString(fmt.Sprintf(`
  # =========================
  # ADMINER
  # =========================
  gui:
    image: adminer
    container_name: %s_gui
    ports:
      - "9090:8080"
    depends_on:
      - db
    networks:
      - %s_network
    restart: always
`, appName, appName))
	}

	// =========================
	// VOLUMES & NETWORKS
	// =========================
	b.WriteString(fmt.Sprintf(`
# =========================
# VOLUMES
# =========================
volumes:
  %s_db_data:

# =========================
# NETWORKS
# =========================
networks:
  %s_network:
    driver: bridge
`, appName, appName))

	return b.String()
}

func generateInstallScript(database string, appName string, useNginx bool, dbGui string) string {
	dbContainer := appName + "_db"
	apiContainer := appName + "_api"

	var dbReadyCheck string
	var dbInfo string
	var extraInfo string

	if database == "mysql" {
		dbReadyCheck = fmt.Sprintf(`# --- Wait for MySQL ---
info "Waiting for MySQL to be ready..."
RETRIES=20
until docker exec %s mysqladmin ping -u root -p"${MYSQL_ROOT_PASSWORD}" --silent 2>/dev/null; do
  RETRIES=$((RETRIES - 1))
  if [ $RETRIES -le 0 ]; then
    error "MySQL is not ready. Check logs: docker logs %s"
  fi
  echo -n "."
  sleep 3
done
echo ""
log "MySQL is ready."`, dbContainer, dbContainer)
		dbInfo = `echo -e "  MySQL       : ${CYAN}localhost:9339${NC}"`
	} else {
		dbReadyCheck = fmt.Sprintf(`# --- Wait for PostgreSQL ---
info "Waiting for PostgreSQL to be ready..."
RETRIES=20
until docker exec %s pg_isready -U "${DB_USER}" --silent 2>/dev/null; do
  RETRIES=$((RETRIES - 1))
  if [ $RETRIES -le 0 ]; then
    error "PostgreSQL is not ready. Check logs: docker logs %s"
  fi
  echo -n "."
  sleep 3
done
echo ""
log "PostgreSQL is ready."`, dbContainer, dbContainer)
		dbInfo = `echo -e "  PostgreSQL  : ${CYAN}localhost:9432${NC}"`
	}

	if useNginx {
		extraInfo += `echo -e "  Nginx       : ${CYAN}http://localhost:9080${NC}"` + "\n"
	}
	if dbGui != "none" {
		extraInfo += `echo -e "  DB GUI      : ${CYAN}http://localhost:9090${NC}"` + "\n"
	}

	containerGrep := apiContainer + "|" + dbContainer
	if useNginx {
		containerGrep += "|" + appName + "_nginx"
	}
	if dbGui != "none" {
		containerGrep += "|" + appName + "_gui"
	}

	return fmt.Sprintf(`#!/bin/bash

# ================================================
#   %s - Deploy Script
# ================================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

log()   { echo -e "${GREEN}[OK]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }
info()  { echo -e "${CYAN}[INFO]${NC} $1"; }

echo ""
echo "================================================"
echo "   %s - Deploy Script"
echo "================================================"
echo ""

command -v docker &>/dev/null || error "Docker not found. Install: https://docs.docker.com/engine/install/"

if docker compose version &>/dev/null; then
  DC="docker compose"
elif command -v docker-compose &>/dev/null; then
  DC="docker-compose"
else
  error "Docker Compose not found."
fi

if [ ! -f ".env" ]; then
  error ".env file not found! Copy from .env.example and fill in the values."
fi

set -o allexport
source .env
set +o allexport

info "Checking for existing containers..."
if docker ps -a --format '{{.Names}}' | grep -qE "%s"; then
  warn "Old containers found, stopping..."
  $DC down
  log "Old containers stopped."
fi

info "Building Docker image..."
$DC build --no-cache || error "Build failed!"
log "Build complete."

info "Starting all containers..."
$DC up -d || error "Failed to start containers!"
log "All containers running."

%s

echo ""
echo "================================================"
echo -e "${GREEN}  Deploy successful!${NC}"
echo "================================================"
echo ""
echo -e "  App  : ${CYAN}http://localhost:9034${NC}"
%s
%s
echo ""
echo -e "  DB Name : ${DB_NAME}"
echo -e "  DB User : ${DB_USER}"
echo -e "  DB Pass : (hidden)"
echo ""
echo "  Logs : docker logs -f %s"
echo "  Stop : $DC down"
echo "================================================"
echo ""
`, appName, appName, containerGrep, dbReadyCheck, dbInfo, extraInfo, apiContainer)
}

func init() {
	rootCmd.AddCommand(dockerInitCmd)
}