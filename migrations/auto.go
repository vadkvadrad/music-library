package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	postgresDb "music-lib/migrations/postgres"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Загрузка переменных окружения (опционально, если файл существует)
	_ = godotenv.Load(".env")

	// Подключение к базе данных
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatalf("DSN environment variable is required")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Выбор действия: migrate или drop
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <migrate|drop>", os.Args[0])
	}

	action := os.Args[1]
	switch action {
	case "migrate":
		fmt.Println("Running migrations...")
		if err := postgresDb.MigrateTables(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations completed successfully.")

	case "drop":
		fmt.Println("Dropping tables...")
		if err := postgresDb.DropTables(db); err != nil {
			log.Fatalf("Failed to drop tables: %v", err)
		}
		fmt.Println("Tables dropped successfully.")

	default:
		log.Fatalf("Unknown action: %s. Use 'migrate' or 'drop'.", action)
	}
}
