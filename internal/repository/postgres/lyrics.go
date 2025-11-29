package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"music-lib/internal/model"
	"music-lib/pkg/db"
)

type LyricsRepository struct {
	db *db.Db
}

func NewLyricsRepository(db *db.Db) *LyricsRepository {
	return &LyricsRepository{
		db: db,
	}
}

func (r *LyricsRepository) GetBySongID(ctx context.Context, songID uint) (*model.Lyrics, error) {
	query := `
		SELECT song_id, created_at, updated_at
		FROM lyrics
		WHERE song_id = $1
	`
	
	var lyrics model.Lyrics
	err := r.db.QueryRowContext(ctx, query, songID).Scan(
		&lyrics.SongID,
		&lyrics.CreatedAt,
		&lyrics.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting lyrics by song id: %w", err)
	}
	
	// Затем получаем куплеты текста
	coupletsQuery := `
		SELECT id, lyrics_id, number, text
		FROM couplets
		WHERE lyrics_id = $1
		ORDER BY number ASC
	`
	
	rows, err := r.db.QueryContext(ctx, coupletsQuery, lyrics.SongID)
	if err != nil {
		return nil, fmt.Errorf("error getting couplets for lyrics: %w", err)
	}
	defer rows.Close()
	
	var couplets []model.Couplet
	for rows.Next() {
		var couplet model.Couplet
		err := rows.Scan(
			&couplet.ID,
			&couplet.LyricsID,
			&couplet.Number,
			&couplet.Text,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning couplet: %w", err)
		}
		couplets = append(couplets, couplet)
	}
	
	lyrics.Couplets = couplets
	return &lyrics, nil
}

func (r *LyricsRepository) Upsert(ctx context.Context, entity *model.Lyrics) error {
	// Сначала проверим, существует ли запись
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM lyrics WHERE song_id = $1)`
	err := r.db.QueryRowContext(ctx, checkQuery, entity.SongID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if lyrics exists: %w", err)
	}
	
	if exists {
		// Обновляем существующую запись
		query := `
			UPDATE lyrics
			SET updated_at = $1
			WHERE song_id = $2
		`
		_, err := r.db.ExecContext(ctx, query, entity.UpdatedAt, entity.SongID)
		if err != nil {
			return fmt.Errorf("error updating lyrics: %w", err)
		}
		
		// Удаляем старые куплеты
		deleteCoupletsQuery := `DELETE FROM couplets WHERE lyrics_id = $1`
		_, err = r.db.ExecContext(ctx, deleteCoupletsQuery, entity.SongID)
		if err != nil {
			return fmt.Errorf("error deleting old couplets: %w", err)
		}
	} else {
		// Создаем новую запись
		query := `
			INSERT INTO lyrics (song_id, created_at, updated_at)
			VALUES ($1, $2, $3)
		`
		_, err := r.db.ExecContext(ctx, query, entity.SongID, entity.CreatedAt, entity.UpdatedAt)
		if err != nil {
			return fmt.Errorf("error creating lyrics: %w", err)
		}
	}
	
	// Вставляем куплеты
	for _, couplet := range entity.Couplets {
		coupletQuery := `
			INSERT INTO couplets (lyrics_id, number, text, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err := r.db.ExecContext(ctx, coupletQuery, 
			entity.SongID, 
			couplet.Number, 
			couplet.Text, 
			entity.CreatedAt, 
			entity.UpdatedAt)
		if err != nil {
			return fmt.Errorf("error inserting couplet: %w", err)
		}
	}
	
	return nil
}

func (r *LyricsRepository) DeleteBySongID(ctx context.Context, songID uint) error {
	query := `DELETE FROM lyrics WHERE song_id = $1`

	result, err := r.db.ExecContext(ctx, query, songID)
	if err != nil {
		return fmt.Errorf("error deleting lyrics: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("lyrics not found for song id %d", songID)
	}

	return nil
}