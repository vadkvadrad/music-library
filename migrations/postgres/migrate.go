package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func MigrateTables(db *sql.DB) error {
	migrationDir := "migrations/postgres"

	// Получаем список всех SQL файлов миграций
	files, err := filepath.Glob(filepath.Join(migrationDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// Сортируем файлы по имени (они должны быть в порядке 001_, 002_, и т.д.)
	sort.Strings(files)

	// Выполняем каждую миграцию
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Удаляем комментарии (строки, начинающиеся с --)
		lines := strings.Split(string(content), "\n")
		var cleanLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			// Пропускаем пустые строки и комментарии
			if trimmed != "" && !strings.HasPrefix(trimmed, "--") {
				cleanLines = append(cleanLines, line)
			}
		}

		// Разделяем по ";", но только если после точки с запятой идет перенос строки или конец файла
		// Это более безопасный способ разделения команд
		statements := make([]string, 0)
		current := ""

		for _, line := range cleanLines {
			current += line + "\n"
			// Если строка заканчивается на ";", это конец команды
			if strings.HasSuffix(strings.TrimSpace(line), ";") {
				stmt := strings.TrimSpace(current)
				if stmt != "" && stmt != ";" {
					statements = append(statements, stmt)
				}
				current = ""
			}
		}

		// Добавляем последнюю команду, если она есть
		if strings.TrimSpace(current) != "" {
			statements = append(statements, strings.TrimSpace(current))
		}

		// Выполняем каждую команду отдельно
		for _, stmt := range statements {
			if stmt == "" {
				continue
			}

			_, err = db.Exec(stmt)
			if err != nil {
				return fmt.Errorf("failed to execute migration from %s: %w\nStatement: %s", file, err, stmt)
			}
		}
	}

	return nil
}
