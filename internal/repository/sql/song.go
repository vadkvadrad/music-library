package sql

import (
	"context"
	"database/sql"
	"fmt"
	"music-lib/internal/model"
)

type SongRepository struct {
	db *DB
}

func NewSongRepository(db *DB) *SongRepository {
	return &SongRepository{
		db: db,
	}
}

func (r *SongRepository) Create(ctx context.Context, entity *model.Song) (*model.Song, error) {
	query := `
		INSERT INTO songs (title, artist_id, album_id, duration, file_path, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id uint
	err := r.db.QueryRowContext(ctx, query,
		entity.Title,
		entity.ArtistID,
		entity.AlbumID,
		entity.Duration,
		entity.FilePath,
		entity.CreatedAt,
		entity.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to create song: %w", err)
	}

	entity.ID = id
	return entity, nil
}

func (r *SongRepository) Update(ctx context.Context, entity *model.Song) (*model.Song, error) {
	query := `
		UPDATE songs
		SET title = $1, artist_id = $2, album_id = $3, duration = $4, file_path = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		entity.Title,
		entity.ArtistID,
		entity.AlbumID,
		entity.Duration,
		entity.FilePath,
		entity.UpdatedAt,
		entity.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update song: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no song found with id %d", entity.ID)
	}

	return entity, nil
}

func (r *SongRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM songs WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no song found with id %d", id)
	}

	return nil
}

func (r *SongRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Song, int64, error) {
	var songs []model.Song
	
	countQuery := `SELECT COUNT(*) FROM songs`
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
		return nil, 0, fmt.Errorf("error counting songs: %w", err)
	}
	
	searchQuery := `
		SELECT id, title, artist_id, album_id, duration, file_path, created_at, updated_at
		FROM songs
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
		return nil, 0, fmt.Errorf("error searching songs: %w", err)
	}
	defer rows.Close()
	
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
			return nil, 0, fmt.Errorf("error scanning song: %w", err)
		}
		songs = append(songs, song)
	}
	
	return songs, total, nil
}

func (r *SongRepository) GetByID(ctx context.Context, id uint) (*model.Song, error) {
	query := `
		SELECT id, title, artist_id, album_id, duration, file_path, created_at, updated_at
		FROM songs
		WHERE id = $1
	`
	
	var song model.Song
	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("song not found with id %d", id)
		}
		return nil, fmt.Errorf("error getting song by id: %w", err)
	}
	
	return &song, nil
}

func (r *SongRepository) GetWithLyrics(ctx context.Context, id uint) (*model.Song, error) {
	// Сначала получаем песню
	song, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Затем получаем текст песни
	lyricsQuery := `
		SELECT song_id, created_at, updated_at
		FROM lyrics
		WHERE song_id = $1
	`
	
	var lyrics model.Lyrics
	err = r.db.QueryRowContext(ctx, lyricsQuery, id).Scan(
		&lyrics.SongID,
		&lyrics.CreatedAt,
		&lyrics.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// Если текста нет, возвращаем песню без текста
			return song, nil
		}
		return nil, fmt.Errorf("error getting lyrics for song: %w", err)
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
	song.Lyrics = lyrics
	return song, nil
}