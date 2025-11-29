package sql

import (
	"context"
	"database/sql"
	"fmt"
	"music-lib/internal/model"
)

type AlbumRepository struct {
	db *DB
}

func NewAlbumRepository(db *DB) *AlbumRepository {
	return &AlbumRepository{
		db: db,
	}
}

func (r *AlbumRepository) Create(ctx context.Context, entity *model.Album) (*model.Album, error) {
	query := `
		INSERT INTO albums (title, artist_id, release_date, cover_art_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id uint
	err := r.db.QueryRowContext(ctx, query,
		entity.Title,
		entity.ArtistID,
		entity.ReleaseDate,
		entity.CoverArtURL,
		entity.CreatedAt,
		entity.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to create album: %w", err)
	}

	entity.ID = id
	return entity, nil
}

func (r *AlbumRepository) Update(ctx context.Context, entity *model.Album) (*model.Album, error) {
	query := `
		UPDATE albums
		SET title = $1, artist_id = $2, release_date = $3, cover_art_url = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		entity.Title,
		entity.ArtistID,
		entity.ReleaseDate,
		entity.CoverArtURL,
		entity.UpdatedAt,
		entity.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update album: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no album found with id %d", entity.ID)
	}

	return entity, nil
}

func (r *AlbumRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM albums WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete album: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no album found with id %d", id)
	}

	return nil
}

func (r *AlbumRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Album, int64, error) {
	var albums []model.Album
	
	countQuery := `SELECT COUNT(*) FROM albums`
	if query != "" {
		countQuery += ` WHERE LOWER(title) LIKE LOWER($1)`
	}
	
	var total int64
	countArgs := []interface{}{}
	if query != "" {
		countArgs = []interface{}{"%" + query + "%"}
	}
	
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting albums: %w", err)
	}
	
	searchQuery := `
		SELECT id, title, artist_id, release_date, cover_art_url, created_at, updated_at
		FROM albums
	`
	
	searchArgs := []interface{}{}
	argIndex := 1
	
	if query != "" {
		searchQuery += ` WHERE LOWER(title) LIKE LOWER($` + fmt.Sprintf("%d", argIndex) + `) `
		searchArgs = append(searchArgs, "%"+query+"%")
		argIndex++
	}
	
	searchQuery += ` ORDER BY title ASC LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)
	searchArgs = append(searchArgs, limit, offset)
	
	rows, err := r.db.QueryContext(ctx, searchQuery, searchArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("error searching albums: %w", err)
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
			return nil, 0, fmt.Errorf("error scanning album: %w", err)
		}
		albums = append(albums, album)
	}
	
	return albums, total, nil
}

func (r *AlbumRepository) GetByID(ctx context.Context, id uint) (*model.Album, error) {
	query := `
		SELECT id, title, artist_id, release_date, cover_art_url, created_at, updated_at
		FROM albums
		WHERE id = $1
	`
	
	var album model.Album
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
			return nil, fmt.Errorf("album not found with id %d", id)
		}
		return nil, fmt.Errorf("error getting album by id: %w", err)
	}
	
	return &album, nil
}

func (r *AlbumRepository) GetWithSongs(ctx context.Context, id uint) (*model.Album, error) {
	// Сначала получаем альбом
	album, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Затем получаем песни альбома
	songsQuery := `
		SELECT id, title, artist_id, album_id, duration, file_path, created_at, updated_at
		FROM songs
		WHERE album_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, songsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error getting songs for album: %w", err)
	}
	defer rows.Close()
	
	var songs []model.Song
	for rows.Next() {
		var song model.Song
		err := rows.Scan(
			&song.ID,
			&song.Title,
			&song.ArtistID,
			&song.AlbumID,
			&song.Duration,
			&song.FilePath,
			&song.CreatedAt,
			&song.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning song: %w", err)
		}
		songs = append(songs, song)
	}
	
	album.Songs = songs
	return album, nil
}