package sql

import (
	"context"
	"database/sql"
	"fmt"
	"music-lib/internal/model"
)

type LyricsRepository struct {
	db *DB
}

func NewLyricsRepository(db *DB) *LyricsRepository {
	return &LyricsRepository{
		db: db,
	}
}

func (r *LyricsRepository) Create(ctx context.Context, entity *model.Lyrics) (*model.Lyrics, error) {
	query := `
		INSERT INTO lyrics (song_id, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING song_id
	`

	var songID uint
	err := r.db.QueryRowContext(ctx, query,
		entity.SongID,
		entity.CreatedAt,
		entity.UpdatedAt,
	).Scan(&songID)

	if err != nil {
		return nil, fmt.Errorf("failed to create lyrics: %w", err)
	}

	entity.SongID = songID
	return entity, nil
}

func (r *LyricsRepository) Update(ctx context.Context, entity *model.Lyrics) (*model.Lyrics, error) {
	query := `
		UPDATE lyrics
		SET updated_at = $1
		WHERE song_id = $2
	`

	result, err := r.db.ExecContext(ctx, query,
		entity.UpdatedAt,
		entity.SongID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update lyrics: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no lyrics found for song id %d", entity.SongID)
	}

	return entity, nil
}

func (r *LyricsRepository) Delete(ctx context.Context, songID uint) error {
	query := `DELETE FROM lyrics WHERE song_id = $1`

	result, err := r.db.ExecContext(ctx, query, songID)
	if err != nil {
		return fmt.Errorf("failed to delete lyrics: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no lyrics found for song id %d", songID)
	}

	return nil
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
			return nil, fmt.Errorf("lyrics not found for song id %d", songID)
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