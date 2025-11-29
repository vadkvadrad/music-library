package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"music-lib/internal/model"
	"music-lib/pkg/db"
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

func (r *SongRepository) ExistsInAlbum(ctx context.Context, albumID uint, songName string) bool {
	query := `SELECT COUNT(*) FROM songs WHERE album_id = $1 AND LOWER(title) = LOWER($2)`
	
	var count int64
	err := r.db.QueryRowContext(ctx, query, albumID, songName).Scan(&count)
	if err != nil {
		return false
	}
	
	return count > 0
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
	
	// Загружаем текст песни, если он есть
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
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("error getting lyrics for song: %w", err)
		}
		// Если текста нет, просто продолжаем без него
	} else {
		// Загружаем куплеты текста
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
	}
	
	return &song, nil
}

func (r *SongRepository) GetByArtistID(ctx context.Context, artistID uint, sort string, limit, offset int) ([]model.Song, int64, error){
	var songs []model.Song
	
	countQuery := `SELECT COUNT(*) FROM songs WHERE artist_id = $1`
	
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, artistID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting songs: %w", err)
	}
	
	searchQuery := `
		SELECT id, title, artist_id, album_id, duration, file_path, created_at, updated_at
		FROM songs
		WHERE artist_id = $1
	`
	
	// Добавляем сортировку
	orderClause := "ORDER BY created_at ASC" // по умолчанию
	switch sort {
	case "title":
		orderClause = "ORDER BY title ASC"
	case "title_desc":
		orderClause = "ORDER BY title DESC"
	case "duration":
		orderClause = "ORDER BY duration ASC"
	case "duration_desc":
		orderClause = "ORDER BY duration DESC"
	case "created_at_desc":
		orderClause = "ORDER BY created_at DESC"
	}
	
	searchQuery += " " + orderClause + " LIMIT $2 OFFSET $3"
	
	rows, err := r.db.QueryContext(ctx, searchQuery, artistID, limit, offset)
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

func (r *SongRepository) GetByAlbumID(ctx context.Context, albumID uint, sort string, limit, offset int) ([]model.Song, int64, error){
	var songs []model.Song
	
	countQuery := `SELECT COUNT(*) FROM songs WHERE album_id = $1`
	
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, albumID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting songs: %w", err)
	}
	
	searchQuery := `
		SELECT id, title, artist_id, album_id, duration, file_path, created_at, updated_at
		FROM songs
		WHERE album_id = $1
	`
	
	// Добавляем сортировку
	orderClause := "ORDER BY created_at ASC" // по умолчанию
	switch sort {
	case "title":
		orderClause = "ORDER BY title ASC"
	case "title_desc":
		orderClause = "ORDER BY title DESC"
	case "duration":
		orderClause = "ORDER BY duration ASC"
	case "duration_desc":
		orderClause = "ORDER BY duration DESC"
	case "created_at_desc":
		orderClause = "ORDER BY created_at DESC"
	}
	
	searchQuery += " " + orderClause + " LIMIT $2 OFFSET $3"
	
	rows, err := r.db.QueryContext(ctx, searchQuery, albumID, limit, offset)
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

func (r *SongRepository) GetFullInfo(ctx context.Context, id uint) (*model.Song, *model.Artist, *model.Album, error) {
	// Получаем песню
	song, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting song: %w", err)
	}
	
	// Получаем артиста
	artistQuery := `
		SELECT id, name, description, formation_year, user_id, created_at, updated_at
		FROM artists
		WHERE id = $1
	`
	
	var artist model.Artist
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
		if err == sql.ErrNoRows {
			return nil, nil, nil, fmt.Errorf("artist not found with id %d", song.ArtistID)
		}
		return nil, nil, nil, fmt.Errorf("error getting artist: %w", err)
	}
	
	// Получаем альбом
	albumQuery := `
		SELECT id, title, artist_id, release_date, cover_art_url, created_at, updated_at
		FROM albums
		WHERE id = $1
	`
	
	var album model.Album
	err = r.db.QueryRowContext(ctx, albumQuery, song.AlbumID).Scan(
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
			return nil, nil, nil, fmt.Errorf("album not found with id %d", song.AlbumID)
		}
		return nil, nil, nil, fmt.Errorf("error getting album: %w", err)
	}
	
	return song, &artist, &album, nil
}