package main

import (
    "fmt"
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/joho/godotenv"
    postgresDb "music-lib/migrations/postgres"
)

func main() {
    // Загрузка переменных окружения
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Подключение к базе данных
    dsn := os.Getenv("DSN")
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Выбор действия: migrate или drop
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
