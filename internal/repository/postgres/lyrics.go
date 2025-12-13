package postgres

import (
	"context"
	"database/sql"
	"errors"
	"music-lib/internal/model"
	"music-lib/pkg/db"
	"time"
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
	var lyrics model.Lyrics
	query := `
		SELECT song_id, created_at, updated_at
		FROM lyrics
		WHERE song_id = $1
	`

	// Конвертируем uint в int для PostgreSQL
	songIDInt := int(songID)
	err := r.db.QueryRowContext(ctx, query, songIDInt).Scan(
		&lyrics.SongID,
		&lyrics.CreatedAt,
		&lyrics.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Загружаем куплеты
	coupletsQuery := `
		SELECT id, lyrics_id, number, text
		FROM couplets
		WHERE lyrics_id = $1
		ORDER BY number ASC
	`

	rows, err := r.db.QueryContext(ctx, coupletsQuery, songIDInt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var couplet model.Couplet
		err := rows.Scan(
			&couplet.ID,
			&couplet.LyricsID,
			&couplet.Number,
			&couplet.Text,
		)
		if err != nil {
			return nil, err
		}
		lyrics.Couplets = append(lyrics.Couplets, couplet)
	}

	return &lyrics, rows.Err()
}

func (r *LyricsRepository) Upsert(ctx context.Context, lyrics *model.Lyrics) error {
	// Сначала вставляем или обновляем lyrics
	query := `
		INSERT INTO lyrics (song_id, created_at, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (song_id) 
		DO UPDATE SET updated_at = $3
		RETURNING created_at, updated_at
	`

	now := time.Now()
	songIDInt := int(lyrics.SongID)
	err := r.db.QueryRowContext(ctx, query, songIDInt, now, now).Scan(
		&lyrics.CreatedAt,
		&lyrics.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Удаляем старые куплеты
	deleteQuery := `DELETE FROM couplets WHERE lyrics_id = $1`
	_, err = r.db.ExecContext(ctx, deleteQuery, songIDInt)
	if err != nil {
		return err
	}

	// Вставляем новые куплеты
	if len(lyrics.Couplets) > 0 {
		insertCoupletQuery := `
			INSERT INTO couplets (lyrics_id, number, text)
			VALUES ($1, $2, $3)
			RETURNING id
		`

		for i := range lyrics.Couplets {
			lyrics.Couplets[i].LyricsID = lyrics.SongID
			lyricsIDInt := int(lyrics.Couplets[i].LyricsID)
			err := r.db.QueryRowContext(ctx, insertCoupletQuery,
				lyricsIDInt,
				lyrics.Couplets[i].Number,
				lyrics.Couplets[i].Text,
			).Scan(&lyrics.Couplets[i].ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *LyricsRepository) DeleteBySongID(ctx context.Context, songID uint) error {
	// Конвертируем uint в int для PostgreSQL
	songIDInt := int(songID)
	
	// Удаляем куплеты (каскадное удаление должно сработать, но для надежности удаляем явно)
	deleteCoupletsQuery := `DELETE FROM couplets WHERE lyrics_id = $1`
	_, err := r.db.ExecContext(ctx, deleteCoupletsQuery, songIDInt)
	if err != nil {
		return err
	}

	// Удаляем lyrics
	deleteQuery := `DELETE FROM lyrics WHERE song_id = $1`
	result, err := r.db.ExecContext(ctx, deleteQuery, songIDInt)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("lyrics not found")
	}
	return nil
}
