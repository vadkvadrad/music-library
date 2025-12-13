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

type AlbumRepository struct {
	db *db.Db
}

func NewAlbumRepository(db *db.Db) *AlbumRepository {
	return &AlbumRepository{
		db: db,
	}
}

func (r *AlbumRepository) Create(ctx context.Context, entity *model.Album) (*model.Album, error) {
	query := `
		INSERT INTO albums (title, artist_id, release_date, cover_art_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx,
		query,
		entity.Title,
		entity.ArtistID,
		entity.ReleaseDate,
		entity.CoverArtURL,
		now,
		now,
	).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)

	return entity, err
}

func (r *AlbumRepository) Update(ctx context.Context, entity *model.Album) (*model.Album, error) {
	query := `
		UPDATE albums 
		SET title = $1, artist_id = $2, release_date = $3, cover_art_url = $4, updated_at = $5
		WHERE id = $6
		RETURNING updated_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx,
		query,
		entity.Title,
		entity.ArtistID,
		entity.ReleaseDate,
		entity.CoverArtURL,
		now,
		entity.ID,
	).Scan(&entity.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("album not found")
		}
		return nil, err
	}
	return entity, nil
}

func (r *AlbumRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM albums WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("album not found")
	}
	return nil
}

func (r *AlbumRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Album, int64, error) {
	var albums []model.Album
	var total int64

	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if query != "" {
		whereClause = "WHERE LOWER(title) LIKE LOWER($" + fmt.Sprintf("%d", argIndex) + ")"
		args = append(args, "%"+query+"%")
		argIndex++
	}

	// Подсчет общего количества
	countQuery := "SELECT COUNT(*) FROM albums " + whereClause
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting albums: %w", err)
	}

	// Получение данных
	selectQuery := `
		SELECT id, title, artist_id, release_date, cover_art_url, created_at, updated_at
		FROM albums
		` + whereClause + `
		ORDER BY title ASC
		LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

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
			return nil, 0, err
		}
		albums = append(albums, album)
	}

	return albums, total, rows.Err()
}

func (r *AlbumRepository) GetByID(ctx context.Context, id uint) (*model.Album, error) {
	var album model.Album
	query := `
		SELECT id, title, artist_id, release_date, cover_art_url, created_at, updated_at
		FROM albums
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
			return nil, errors.New("album not found")
		}
		return nil, err
	}
	return &album, nil
}

func (r *AlbumRepository) GetWithSongs(ctx context.Context, id uint) (*model.Album, error) {
	album, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (r *AlbumRepository) GetSongsByAlbumID(ctx context.Context, albumID uint) ([]model.Song, error) {
	query := `
		SELECT id, title, artist_id, album_id, duration, file_path, created_at, updated_at
		FROM songs
		WHERE album_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []model.Song
	for rows.Next() {
		var song model.Song
		var albumIDVal sql.NullInt64
		err := rows.Scan(
			&song.ID,
			&song.Title,
			&song.ArtistID,
			&albumIDVal,
			&song.Duration,
			&song.FilePath,
			&song.CreatedAt,
			&song.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if albumIDVal.Valid {
			albumIDUint := uint(albumIDVal.Int64)
			song.AlbumID = &albumIDUint
		}
		songs = append(songs, song)
	}

	return songs, rows.Err()
}
