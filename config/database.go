package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"gorm.io/gorm"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
)

var dbconn *gorm.DB
var once sync.Once
var dbConnected bool
var DB_PREFIX string

func LoadEnv() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Gagal load .env: %v", err)
		}
	}

	DB_PREFIX = os.Getenv("DB_PREFIX")
}

func GetDBPrefix(tableName string) string {
	return DB_PREFIX + "_" + tableName
}

func ConnectDB() {
	LoadEnv();
	once.Do(func() {

		DB_CONNECTION := os.Getenv("DB_CONNECTION")
		DB_HOST := os.Getenv("DB_HOST")
		DB_NAME := os.Getenv("DB_NAME")
		DB_USER := os.Getenv("DB_USER")
		DB_PASS := os.Getenv("DB_PASS")
		DB_PORT := os.Getenv("DB_PORT")

		if DB_CONNECTION == "" {
			log.Fatal("DB_CONNECTION belum diisi!")
		}

		if DB_HOST == "" || DB_NAME == "" || DB_USER == "" || DB_PORT == "" {
			log.Fatal("Pastikan semua variabel database di .env sudah diisi!")
		}

		if !dbConnected {

			for i := 0; i < 5; i++ {

				var dialector gorm.Dialector

				switch DB_CONNECTION {

				case "mysql":

					dsn := fmt.Sprintf(
						"%s:%s@tcp(%s:%s)/%s?parseTime=True",
						DB_USER,
						DB_PASS,
						DB_HOST,
						DB_PORT,
						DB_NAME,
					)

					dialector = mysql.Open(dsn)

				case "postgres":

					dsn := fmt.Sprintf(
						"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
						DB_HOST,
						DB_USER,
						DB_PASS,
						DB_NAME,
						DB_PORT,
					)

					dialector = postgres.Open(dsn)

				default:
					log.Fatalf("Database driver '%s' tidak didukung!", DB_CONNECTION)
				}

				database, err := gorm.Open(dialector, &gorm.Config{
					SkipDefaultTransaction: true,
					PrepareStmt:            true,
				})

				if err == nil {

					sqlDB, err := database.DB()

					if err == nil && sqlDB.Ping() == nil {

						dbconn = database
						dbConnected = true

						log.Printf("Database connection successful using %s!", DB_CONNECTION)

						return
					}

					log.Println("Connection successful, but ping failed, retrying...")
				}

				log.Println("Connection failed, retrying...", err)

				time.Sleep(2 * time.Second)
			}

			log.Fatal("Failed to connect to the database after 5 attempts.")

		} else {
			log.Println("Database connection was successful previously, not retrying.")
		}
	})
}

func GetDB() *gorm.DB {

	ConnectDB()

	if dbconn == nil {
		log.Fatal("Database is not connected. Make sure ConnectDB() has been called!")
	}

	return dbconn
}