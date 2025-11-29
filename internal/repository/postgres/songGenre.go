package postgres

import (
	"context"
	"errors"
	"music-lib/internal/model"
	"music-lib/pkg/db"
)

type SongGenreRepository struct {
	db *db.Db
}

func NewSongGenreRepository(db *db.Db) *SongGenreRepository {
	return &SongGenreRepository{
		db: db,
	}
}

func (r *SongGenreRepository) Create(ctx context.Context, entity *model.SongGenre) (*model.SongGenre, error) {
	query := `
		INSERT INTO song_genres (song_id, genre_id)
		VALUES ($1, $2)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, entity.SongID, entity.GenreID).Scan(&entity.ID)
	return entity, err
}

func (r *SongGenreRepository) Update(ctx context.Context, entity *model.SongGenre) (*model.SongGenre, error) {
	query := `
		UPDATE song_genres 
		SET song_id = $1, genre_id = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, entity.SongID, entity.GenreID, entity.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("song_genre not found")
	}
	return entity, nil
}

func (r *SongGenreRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM song_genres WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("song_genre not found")
	}
	return nil
}
