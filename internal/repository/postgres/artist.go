package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"music-lib/internal/model"
	"music-lib/pkg/db"
	"time"
)

type ArtistRepository struct {
	db *db.Db
}

func NewArtistRepository(db *db.Db) *ArtistRepository {
	return &ArtistRepository{
		db: db,
	}
}

func (r *ArtistRepository) Create(ctx context.Context, entity *model.Artist) (*model.Artist, error) {
	query := `
		INSERT INTO artists (name, description, formation_year, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx,
		query,
		entity.Name,
		entity.Description,
		entity.FormationYear,
		entity.UserID,
		now,
		now,
	).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)

	return entity, err
}

func (r *ArtistRepository) Update(ctx context.Context, entity *model.Artist) (*model.Artist, error) {
	query := `
		UPDATE artists 
		SET name = $1, description = $2, formation_year = $3, user_id = $4, updated_at = $5
		WHERE id = $6
		RETURNING updated_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx,
		query,
		entity.Name,
		entity.Description,
		entity.FormationYear,
		entity.UserID,
		now,
		entity.ID,
	).Scan(&entity.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("artist not found")
		}
		return nil, err
	}
	return entity, nil
}

func (r *ArtistRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM artists WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("artist not found")
	}
	return nil
}

func (r *ArtistRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Artist, int64, error) {
	var artists []model.Artist
	var total int64

	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if query != "" {
		whereClause = "WHERE LOWER(name) LIKE LOWER($" + fmt.Sprintf("%d", argIndex) + ")"
		args = append(args, "%"+query+"%")
		argIndex++
	}

	// Подсчет общего количества
	countQuery := "SELECT COUNT(*) FROM artists " + whereClause
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting artists: %w", err)
	}

	// Получение данных
	selectQuery := `
		SELECT id, name, description, formation_year, user_id, created_at, updated_at
		FROM artists
		` + whereClause + `
		ORDER BY name ASC
		LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var artist model.Artist
		err := rows.Scan(
			&artist.ID,
			&artist.Name,
			&artist.Description,
			&artist.FormationYear,
			&artist.UserID,
			&artist.CreatedAt,
			&artist.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		artists = append(artists, artist)
	}

	return artists, total, rows.Err()
}

func (r *ArtistRepository) GetByID(ctx context.Context, id uint) (*model.Artist, error) {
	var artist model.Artist
	query := `
		SELECT id, name, description, formation_year, user_id, created_at, updated_at
		FROM artists
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&artist.ID,
		&artist.Name,
		&artist.Description,
		&artist.FormationYear,
		&artist.UserID,
		&artist.CreatedAt,
		&artist.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("artist not found")
		}
		return nil, err
	}
	return &artist, nil
}

func (r *ArtistRepository) GetByUserID(ctx context.Context, userID uint) (*model.Artist, error) {
	var artist model.Artist
	query := `
		SELECT id, name, description, formation_year, user_id, created_at, updated_at
		FROM artists
		WHERE user_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&artist.ID,
		&artist.Name,
		&artist.Description,
		&artist.FormationYear,
		&artist.UserID,
		&artist.CreatedAt,
		&artist.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("artist not found")
		}
		return nil, err
	}
	return &artist, nil
}

func (r *ArtistRepository) GetWithAlbums(ctx context.Context, id uint) (*model.Artist, error) {
	artist, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return artist, nil
}

func (r *ArtistRepository) GetAlbumsByArtistID(ctx context.Context, artistID uint) ([]model.Album, error) {
	query := `
		SELECT id, title, artist_id, release_date, cover_art_url, created_at, updated_at
		FROM albums
		WHERE artist_id = $1
		ORDER BY release_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []model.Album
	for rows.Next() {
		var album model.Album
		err := rows.Scan(
			&album.ID,
			&album.Title,
			&album.ArtistID,
			&album.ReleaseDate,
			&album.CoverArtURL,
			&album.CreatedAt,
			&album.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, rows.Err()
}

func (r *ArtistRepository) IsExists(ctx context.Context, name string) bool {
	var count int64
	query := `
		SELECT COUNT(*) 
		FROM artists
		WHERE LOWER(name) = LOWER($1)
		LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, name).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (r *ArtistRepository) GetArtistAlbumByUserID(ctx context.Context, userID uint, albumID uint) (*model.Album, int, error) {
	artist, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	var album model.Album
	query := `
		SELECT id, title, artist_id, release_date, cover_art_url, created_at, updated_at
		FROM albums
		WHERE id = $1 AND artist_id = $2
	`

	err = r.db.QueryRowContext(ctx, query, albumID, artist.ID).Scan(
		&album.ID,
		&album.Title,
		&album.ArtistID,
		&album.ReleaseDate,
		&album.CoverArtURL,
		&album.CreatedAt,
		&album.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	// Подсчет количества альбомов
	var count int
	countQuery := `SELECT COUNT(*) FROM albums WHERE artist_id = $1`
	err = r.db.QueryRowContext(ctx, countQuery, artist.ID).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return &album, count, nil
}
