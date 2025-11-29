package postgres

import (
	"context"
	"database/sql"
	"errors"
	"music-lib/internal/model"
	"music-lib/pkg/db"
)

type GenreRepository struct {
	db *db.Db
}

func NewGenreRepository(db *db.Db) *GenreRepository {
	return &GenreRepository{
		db: db,
	}
}

func (r *GenreRepository) Create(ctx context.Context, entity *model.Genre) (*model.Genre, error) {
	query := `
		INSERT INTO genres (name)
		VALUES ($1)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, entity.Name).Scan(&entity.ID)
	return entity, err
}

func (r *GenreRepository) Update(ctx context.Context, entity *model.Genre) (*model.Genre, error) {
	query := `
		UPDATE genres 
		SET name = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, entity.Name, entity.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("genre not found")
	}
	return entity, nil
}

func (r *GenreRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM genres WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("genre not found")
	}
	return nil
}

func (r *GenreRepository) IsExists(ctx context.Context, name string) bool {
	var count int64
	query := `
		SELECT COUNT(*) 
		FROM genres
		WHERE LOWER(name) = LOWER($1)
		LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, name).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (r *GenreRepository) GetById(ctx context.Context, id uint) (*model.Genre, error) {
	var genre model.Genre
	query := `SELECT id, name FROM genres WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(&genre.ID, &genre.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("genre not found")
		}
		return nil, err
	}
	return &genre, nil
}

func (r *GenreRepository) GetByIds(ctx context.Context, ids []uint) ([]model.Genre, error) {
	if len(ids) == 0 {
		return []model.Genre{}, nil
	}

	// Используем ANY для массива
	query := `SELECT id, name FROM genres WHERE id = ANY($1::int[])`

	rows, err := r.db.QueryContext(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []model.Genre
	for rows.Next() {
		var genre model.Genre
		err := rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}

	return genres, rows.Err()
}
