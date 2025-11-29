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

type SongRepository struct {
	db *db.Db
}

func NewSongRepository(db *db.Db) *SongRepository {
	return &SongRepository{
		db: db,
	}
}

func (r *SongRepository) Create(ctx context.Context, entity *model.Song) (*model.Song, error) {
	query := `
		INSERT INTO songs (title, artist_id, album_id, duration, file_path, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx,
		query,
		entity.Title,
		entity.ArtistID,
		entity.AlbumID,
		entity.Duration,
		entity.FilePath,
		now,
		now,
	).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)

	return entity, err
}

func (r *SongRepository) Update(ctx context.Context, entity *model.Song) (*model.Song, error) {
	query := `
		UPDATE songs 
		SET title = $1, artist_id = $2, album_id = $3, duration = $4, file_path = $5, updated_at = $6
		WHERE id = $7
		RETURNING updated_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx,
		query,
		entity.Title,
		entity.ArtistID,
		entity.AlbumID,
		entity.Duration,
		entity.FilePath,
		now,
		entity.ID,
	).Scan(&entity.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("song not found")
		}
		return nil, err
	}
	return entity, nil
}

func (r *SongRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM songs WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("song not found")
	}
	return nil
}

func (r *SongRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Song, int64, error) {
	var songs []model.Song
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
	countQuery := "SELECT COUNT(*) FROM songs " + whereClause
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting songs: %w", err)
	}

	// Получение данных
	selectQuery := `
		SELECT id, title, artist_id, album_id, duration, file_path, created_at, updated_at
		FROM songs
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
		var song model.Song
		var albumID sql.NullInt64
		err := rows.Scan(
			&song.ID,
			&song.Title,
			&song.ArtistID,
			&albumID,
			&song.Duration,
			&song.FilePath,
			&song.CreatedAt,
			&song.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		if albumID.Valid {
			albumIDUint := uint(albumID.Int64)
			song.AlbumID = &albumIDUint
		}
		songs = append(songs, song)
	}

	return songs, total, rows.Err()
}

func (r *SongRepository) ExistsInAlbum(ctx context.Context, albumID uint, songName string) bool {
	var count int64
	query := `
		SELECT COUNT(*) 
		FROM songs
		WHERE album_id = $1 AND LOWER(title) = LOWER($2)
		LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, albumID, songName).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (r *SongRepository) GetByID(ctx context.Context, id uint) (*model.Song, error) {
	var song model.Song
	var albumID sql.NullInt64
	query := `
		SELECT id, title, artist_id, album_id, duration, file_path, created_at, updated_at
		FROM songs
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&song.ID,
		&song.Title,
		&song.ArtistID,
		&albumID,
		&song.Duration,
		&song.FilePath,
		&song.CreatedAt,
		&song.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("song not found")
		}
		return nil, err
	}

	if albumID.Valid {
		albumIDUint := uint(albumID.Int64)
		song.AlbumID = &albumIDUint
	}

	return &song, nil
}

func (r *SongRepository) GetByArtistID(ctx context.Context, artistID uint, sort string, limit, offset int) ([]model.Song, int64, error) {
	var songs []model.Song
	var total int64

	orderBy := "title ASC"
	if sort == "desc" {
		orderBy = "title DESC"
	}

	// Подсчет общего количества
	countQuery := `SELECT COUNT(*) FROM songs WHERE artist_id = $1`
	err := r.db.QueryRowContext(ctx, countQuery, artistID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting songs: %w", err)
	}

	// Получение данных
	query := fmt.Sprintf(`
		SELECT id, title, artist_id, album_id, duration, file_path, created_at, updated_at
		FROM songs
		WHERE artist_id = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, artistID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var song model.Song
		var albumID sql.NullInt64
		err := rows.Scan(
			&song.ID,
			&song.Title,
			&song.ArtistID,
			&albumID,
			&song.Duration,
			&song.FilePath,
			&song.CreatedAt,
			&song.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		if albumID.Valid {
			albumIDUint := uint(albumID.Int64)
			song.AlbumID = &albumIDUint
		}
		songs = append(songs, song)
	}

	return songs, total, rows.Err()
}

func (r *SongRepository) GetByAlbumID(ctx context.Context, albumID uint, sort string, limit, offset int) ([]model.Song, int64, error) {
	var songs []model.Song
	var total int64

	orderBy := "created_at DESC"
	if sort == "asc" {
		orderBy = "created_at ASC"
	}

	// Подсчет общего количества
	countQuery := `SELECT COUNT(*) FROM songs WHERE album_id = $1`
	err := r.db.QueryRowContext(ctx, countQuery, albumID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting songs: %w", err)
	}

	// Получение данных
	query := fmt.Sprintf(`
		SELECT id, title, artist_id, album_id, duration, file_path, created_at, updated_at
		FROM songs
		WHERE album_id = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, albumID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

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
			return nil, 0, err
		}
		if albumIDVal.Valid {
			albumIDUint := uint(albumIDVal.Int64)
			song.AlbumID = &albumIDUint
		}
		songs = append(songs, song)
	}

	return songs, total, rows.Err()
}

func (r *SongRepository) GetFullInfo(ctx context.Context, id uint) (*model.Song, *model.Artist, *model.Album, error) {
	song, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}

	// Получение артиста
	var artist model.Artist
	artistQuery := `
		SELECT id, name, description, formation_year, user_id, created_at, updated_at
		FROM artists
		WHERE id = $1
	`
	err = r.db.QueryRowContext(ctx, artistQuery, song.ArtistID).Scan(
		&artist.ID,
		&artist.Name,
		&artist.Description,
		&artist.FormationYear,
		&artist.UserID,
		&artist.CreatedAt,
		&artist.UpdatedAt,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting artist: %w", err)
	}

	// Получение альбома (если есть)
	var album *model.Album
	if song.AlbumID != nil {
		albumQuery := `
			SELECT id, title, artist_id, release_date, cover_art_url, created_at, updated_at
			FROM albums
			WHERE id = $1
		`
		var albumData model.Album
		err = r.db.QueryRowContext(ctx, albumQuery, *song.AlbumID).Scan(
			&albumData.ID,
			&albumData.Title,
			&albumData.ArtistID,
			&albumData.ReleaseDate,
			&albumData.CoverArtURL,
			&albumData.CreatedAt,
			&albumData.UpdatedAt,
		)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, nil, nil, fmt.Errorf("error getting album: %w", err)
			}
		} else {
			album = &albumData
		}
	}

	return song, &artist, album, nil
}
