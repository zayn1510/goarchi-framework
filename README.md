# Goarchi

> Simple Layered Architecture Generator for Golang

Goarchi is a CLI tool that helps you scaffold a clean, layered Go project structure — controllers, services, models, migrations, requests, resources, and seeders — in seconds.

---

## Requirements

- Go `1.26.3` or higher

---

## Getting Started

### 1. Create a new project

Clone Goarchi as your project base, then point the remote to your own repo:

```bash
git clone https://github.com/zayn1510/goarchi.git myapp
cd myapp

# Remove Goarchi's remote and point to your own repository
git remote remove origin
git remote add origin https://github.com/username/myapp.git
git branch -M main
git push -u origin main
```

---

### 2. Install the CLI

Build the binary from inside your project:

```bash
go run cli/main.go build
```

Move binary to PATH:

**Linux / macOS:**
```bash
sudo mv goarchi /usr/local/bin/goarchi
```

**Windows:**

Move `goarchi.exe` to a folder that's in your `PATH`, or add its directory to the environment variable manually.

Verify:
```bash
goarchi version
```

---

### 3. Setup environment

```bash
cp .env.example .env
```

Edit `.env` with your database credentials and app settings.

---

### 4. Setup Docker (optional)

```bash
goarchi docker:init
```

Follow the interactive wizard, then deploy with:

```bash
bash install.sh
```

---

### 5. Run the app

```bash
go run main.go
```

---

## Usage

```bash
goarchi make [command] [name] [fields...]
```

### Available Commands

#### Controller
```bash
goarchi make controller [name]
```
Generates a controller file at `app/controllers/[name]_controller.go`.

**Example:**
```bash
goarchi make controller user
# -> app/controllers/user_controller.go (struct: UserController)
```

---

#### Service
```bash
goarchi make service [name]
```
Generates a service file at `app/services/[name]_service.go`.

**Example:**
```bash
goarchi make service user
# -> app/services/user_service.go (struct: UserService)
```

---

#### Request
```bash
goarchi make request [name] [field:type...]
```
Generates a request struct with validation tags at `app/requests/[name]_request.go`.

**Example:**
```bash
goarchi make request user name:string age:int email:string
# -> app/requests/user_request.go (struct: UserRequest)
```

---

#### Resource
```bash
goarchi make resource [name] [field:type...]
```
Generates a response formatter (DTO) at `app/resources/[name]_resource.go`.

**Example:**
```bash
goarchi make resource user id:int name:string email:string
# -> app/resources/user_resource.go (struct: UserResource)
```

---

#### Model
```bash
goarchi make model [name] [field:type;gorm_tag...]
```
Generates a GORM model at `app/models/[name].go`. `CreatedAt` and `UpdatedAt` are added automatically.

**Example:**
```bash
goarchi make model users "id:int;primaryKey" "name:string;not null"
# -> app/models/users.go (struct: User)
```

For relations (foreign key):
```bash
goarchi make model users "Role:foreignKey:RoleID"
# -> adds: Role *Role `gorm:"foreignKey:RoleID"`
```

---

#### Migration (Generate file)
```bash
goarchi make migration [name]
```
Generates a timestamped migration file at `database/migrations/[timestamp]_[name].go`.

**Example:**
```bash
goarchi make migration create_users_table
# -> database/migrations/20240101120000_create_users_table.go
```

Don't forget to register it in `database/migrations/migrations_list.go`:
```go
{
    Name: "create_users_table",
    Up:   CreateUsersTable,
    Down: DropUsersTable,
}
```

---

#### Seeder (Generate file)
```bash
goarchi make seeder [name]
```
Generates a timestamped seeder file at `database/seeders/[timestamp]_[name].go`.

**Example:**
```bash
goarchi make seeder role_seeder
# -> database/seeders/20240101120000_role_seeder.go
```

Don't forget to register it in `database/seeders/seeder_list.go`:
```go
{
    Name: "role_seeder",
    Up:   SeedRoleSeeder,
    Down: DropSeedRoleSeeder,
}
```

---

#### Migrate (Run)
```bash
# Run all migrations
go run cli/main.go migrate up
go run cli/main.go migrate down

# Run a specific migration
go run cli/main.go migrate up create_users_table
go run cli/main.go migrate down create_users_table
```

- `up` — applies migrations
- `down` — drops tables
- Optionally pass a migration name to run only that one

> Migrate is run via `go run cli/main.go` not from the `goarchi` binary. This is because Go compiles to a static binary and cannot access migration files added after build time.

---

#### Seed (Run)
```bash
# Run all seeders
go run cli/main.go seed up
go run cli/main.go seed down

# Run a specific seeder
go run cli/main.go seed up role_seeder
go run cli/main.go seed down role_seeder
```

- `up` — inserts seed data
- `down` — removes seed data
- Optionally pass a seeder name to run only that one

---

#### Check Vulnerabilities
```bash
goarchi govulncheck
```

Automatically installs/updates `govulncheck` and scans all dependencies for known vulnerabilities.

If vulnerabilities are found, fix them by running:
```bash
go get -u ./...
go mod tidy
```

---

## Docker Setup

```bash
goarchi docker:init
```

Interactive wizard that generates all Docker configuration files for your project.

**Prompts:**

1. **Go Version** — `1.24`, `1.25`, `1.26`, or `latest`
2. **Database** — PostgreSQL or MySQL
3. **Air Hot Reload** — yes or no
4. **Nginx** — use as reverse proxy or skip
5. **DB GUI**
   - MySQL -> phpMyAdmin / Adminer / Skip
   - PostgreSQL -> Adminer / Skip
6. **App name** — used as prefix for all container names

**Generated files:**

| File | Always |
|------|--------|
| `Dockerfile.dev` | yes |
| `docker-compose.yml` | yes |
| `install.sh` | yes |
| `nginx/default.conf` | Only if Nginx selected |

**Ports:**

| Service | Port |
|---------|------|
| App | `9034:8080` |
| Nginx | `9080:80` |
| DB GUI (phpMyAdmin / Adminer) | `9090:80` |
| PostgreSQL | `9432:5432` |
| MySQL | `9339:3306` |

---

## Project Structure

```
.
├── app/
│   ├── controllers/       # HTTP handlers
│   ├── middleware/        # JWT, CORS, rate limiting
│   ├── migrate/           # Migration runner
│   ├── models/            # GORM models
│   ├── requests/          # Request structs & validation
│   ├── resources/         # Response formatters (DTO)
│   └── services/          # Business logic
├── cli/
│   └── main.go            # CLI entry point (for building binary)
├── cmd/                   # Cobra CLI commands
├── config/                # App & database configuration
├── core/
│   ├── generate/          # Code generator
│   │   └── templates/     # .tmpl files for each layer
│   └── tools/             # Utility helpers
├── database/
│   ├── migrations/        # Migration files & registry
│   └── seeders/           # Seeder files & registry
├── routers/
│   └── web.go             # Route definitions
├── main.go                # App entry point
└── .env.example           # Environment variable template
```

---

## Configuration

Copy `.env.example` to `.env` and fill in your values:

```bash
cp .env.example .env
```

```env
# ======================
# APP
# ======================
APP_NAME=my-app
APP_ENV=development
APP_PORT=8080

# ======================
# DATABASE
# ======================
DB_CONNECTION=postgres   # postgres or mysql
DB_HOST=127.0.0.1
DB_PORT=5432             # 5432 for postgres, 3306 for mysql
DB_NAME=dbname
DB_USER=dbuser
DB_PASS=dbpass
DB_PREFIX=tbl_

# ======================
# JWT
# ======================
JWT_SECRET_KEY=your-secret
JWT_EXPIRED_TOKEN=5h
```

> When using Docker, `DB_HOST` should match the service name in `docker-compose.yml` (default: `db`).

---

## License

MIT License. See [LICENSE](LICENSE) for details.