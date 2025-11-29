package db

import (
	"database/sql"
	"fmt"
)

// CreateTables создает все таблицы в базе данных
func CreateTables(db *sql.DB) error {
	queries := []string{
		// Таблица пользователей
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			avatar_url TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Таблица артистов
		`CREATE TABLE IF NOT EXISTS artists (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			formation_year TIMESTAMP WITH TIME ZONE,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Таблица альбомов
		`CREATE TABLE IF NOT EXISTS albums (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			artist_id INTEGER REFERENCES artists(id) ON DELETE CASCADE,
			release_date TIMESTAMP WITH TIME ZONE,
			cover_art_url TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Таблица жанров
		`CREATE TABLE IF NOT EXISTS genres (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Таблица песен
		`CREATE TABLE IF NOT EXISTS songs (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			artist_id INTEGER REFERENCES artists(id) ON DELETE CASCADE,
			album_id INTEGER REFERENCES albums(id) ON DELETE CASCADE,
			duration INTEGER, -- длительность в секундах
			file_path TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Таблица текстов песен
		`CREATE TABLE IF NOT EXISTS lyrics (
			song_id INTEGER PRIMARY KEY REFERENCES songs(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Таблица куплетов
		`CREATE TABLE IF NOT EXISTS couplets (
			id SERIAL PRIMARY KEY,
			lyrics_id INTEGER REFERENCES lyrics(song_id) ON DELETE CASCADE,
			number INTEGER NOT NULL,
			text TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Промежуточная таблица для связи песен и жанров
		`CREATE TABLE IF NOT EXISTS song_genres (
			id SERIAL PRIMARY KEY,
			song_id INTEGER REFERENCES songs(id) ON DELETE CASCADE,
			genre_id INTEGER REFERENCES genres(id) ON DELETE CASCADE,
			UNIQUE(song_id, genre_id)
		);`,

		// Таблица профилей
		`CREATE TABLE IF NOT EXISTS profiles (
			user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
			bio TEXT,
			birth_date DATE,
			location VARCHAR(255),
			website VARCHAR(255),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Таблица разрешений
		`CREATE TABLE IF NOT EXISTS permissions (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for i, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration query %d: %w", i, err)
		}
	}

	return nil
}

// DropTables удаляет все таблицы из базы данных
func DropTables(db *sql.DB) error {
	queries := []string{
		`DROP TABLE IF EXISTS song_genres CASCADE;`,
		`DROP TABLE IF EXISTS couplets CASCADE;`,
		`DROP TABLE IF EXISTS lyrics CASCADE;`,
		`DROP TABLE IF EXISTS songs CASCADE;`,
		`DROP TABLE IF EXISTS albums CASCADE;`,
		`DROP TABLE IF EXISTS genres CASCADE;`,
		`DROP TABLE IF EXISTS artists CASCADE;`,
		`DROP TABLE IF EXISTS profiles CASCADE;`,
		`DROP TABLE IF EXISTS permissions CASCADE;`,
		`DROP TABLE IF EXISTS users CASCADE;`,
	}

	for i, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute drop query %d: %w", i, err)
		}
	}

	return nil
}