# 📦 Goarchi

> Simple Layered Architecture Generator for Golang

Goarchi is a CLI tool that helps you scaffold a clean, layered Go project structure — controllers, services, models, migrations, requests, and resources — in seconds.

---

## 🚀 Getting Started

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

Output:
```
📦 Building Goarchi binary...
✅ Binary generated successfully!
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

## 🔧 Usage

```bash
goarchi make [command] [name] [fields...]
```

### Available Commands

#### 🔧 Controller
```bash
goarchi make controller [name]
```
Generates a controller file at `app/controllers/[name]_controller.go`.

**Example:**
```bash
goarchi make controller user
# → app/controllers/user_controller.go (struct: UserController)
```

---

#### 🛠️ Service
```bash
goarchi make service [name]
```
Generates a service file at `app/services/[name]_service.go`.

**Example:**
```bash
goarchi make service user
# → app/services/user_service.go (struct: UserService)
```

---

#### 📝 Request
```bash
goarchi make request [name] [field:type...]
```
Generates a request struct with validation tags at `app/requests/[name]_request.go`.

**Example:**
```bash
goarchi make request user name:string age:int email:string
# → app/requests/user_request.go (struct: UserRequest)
```

---

#### 📦 Resource
```bash
goarchi make resource [name] [field:type...]
```
Generates a response formatter (DTO) at `app/resources/[name]_resource.go`.

**Example:**
```bash
goarchi make resource user id:int name:string email:string
# → app/resources/user_resource.go (struct: UserResource)
```

---

#### 🧩 Model
```bash
goarchi make model [name] [field:type;gorm_tag...]
```
Generates a GORM model at `app/models/[name].go`. `CreatedAt` and `UpdatedAt` are added automatically.

**Example:**
```bash
goarchi make model users "id:int;primaryKey" "name:string;not null"
# → app/models/users.go (struct: User)
```

For relations (foreign key):
```bash
goarchi make model users "Role:foreignKey:RoleID"
# → adds: Role Role `gorm:"foreignKey:RoleID"`
```

---

#### 🗃️ Migration (Generate file)
```bash
goarchi make migration [name]
```
Generates a timestamped migration file at `database/migrations/[timestamp]_[name].go`.

**Example:**
```bash
goarchi make migration create_users_table
# → database/migrations/20240101120000_create_users_table.go
```

---

#### 🧬 Migrate (Run)
```bash
goarchi make migrate up
goarchi make migrate down
```
Runs all migrations registered in `database/migrations/migrations_list.go`.

- `up` — applies all migrations in `AllMigrations`
- `down` — rolls back all migrations in `AllMigrations`

---

## 🐳 Docker Setup

```bash
goarchi docker:init
```

Interactive wizard that generates all Docker configuration files for your project.

**Prompts:**

1. **Go Version** — `1.22`, `1.23`, or `latest`
2. **Database** — PostgreSQL or MySQL
3. **Air Hot Reload** — yes or no
4. **Nginx** — use as reverse proxy or skip
5. **DB GUI**
   - MySQL → phpMyAdmin / Adminer / Skip
   - PostgreSQL → Adminer / Skip
6. **App name** — used as prefix for all container names

**Generated files:**

| File | Always |
|------|--------|
| `Dockerfile.dev` | ✅ |
| `docker-compose.yml` | ✅ |
| `install.sh` | ✅ |
| `nginx/default.conf` | Only if Nginx selected |

**Example session:**
```
🐳 Goarchi Docker Setup
========================
Choose Go Version:
1. 1.22  2. 1.23  3. latest
Enter choice (1-3): 2

Choose Database:
1. PostgreSQL  2. MySQL
Enter choice (1-2): 2

Use Air Hot Reload? (y/n): y

Use Nginx as reverse proxy? (y/n): y

Choose DB GUI:
1. phpMyAdmin  2. Adminer  3. Skip
Enter choice (1-3): 1

Enter your app name: myapp

✅ Docker configuration generated successfully!
📄 Dockerfile.dev
📄 docker-compose.yml
📄 nginx/default.conf
📄 install.sh

🔖 Container names:
   App   : myapp_api
   DB    : myapp_db
   Nginx : myapp_nginx
   DB GUI (phpmyadmin) : myapp_gui

▶️  To deploy, run: bash install.sh
```

**Ports:**

| Service | Port |
|---------|------|
| App | `9034:8080` |
| Nginx | `9080:80` |
| DB GUI (phpMyAdmin / Adminer) | `9090:80` |
| PostgreSQL | `9432:5432` |
| MySQL | `9339:3306` |

---

## 📁 Project Structure

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
│   └── migrations/        # SQL migration files
├── routers/
│   └── web.go             # Route definitions
├── public/                # Static files
├── main.go                # App entry point
└── config.yaml            # App configuration
```

---

## ⚙️ Configuration

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
DB_HOST=db
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

> ⚠️ When using Docker, `DB_HOST` should match the service name in `docker-compose.yml` (default: `db`).

---

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.