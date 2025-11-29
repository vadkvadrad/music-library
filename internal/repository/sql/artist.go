package sql

import (
	"context"
	"database/sql"
	"fmt"
	"music-lib/internal/model"
)

type ArtistRepository struct {
	db *DB
}

func NewArtistRepository(db *DB) *ArtistRepository {
	return &ArtistRepository{
		db: db,
	}
}

func (r *ArtistRepository) Create(ctx context.Context, entity *model.Artist) (*model.Artist, error) {
	query := `
		INSERT INTO artists (name, description, formation_year, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id uint
	err := r.db.QueryRowContext(ctx, query,
		entity.Name,
		entity.Description,
		entity.FormationYear,
		entity.UserID,
		entity.CreatedAt,
		entity.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to create artist: %w", err)
	}

	entity.ID = id
	return entity, nil
}

func (r *ArtistRepository) Update(ctx context.Context, entity *model.Artist) (*model.Artist, error) {
	query := `
		UPDATE artists
		SET name = $1, description = $2, formation_year = $3, user_id = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		entity.Name,
		entity.Description,
		entity.FormationYear,
		entity.UserID,
		entity.UpdatedAt,
		entity.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update artist: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no artist found with id %d", entity.ID)
	}

	return entity, nil
}

func (r *ArtistRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM artists WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete artist: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no artist found with id %d", id)
	}

	return nil
}

func (r *ArtistRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Artist, int64, error) {
	var artists []model.Artist
	
	countQuery := `SELECT COUNT(*) FROM artists`
	if query != "" {
		countQuery += ` WHERE LOWER(name) LIKE LOWER($1)`
	}
	
	var total int64
	countArgs := []interface{}{}
	if query != "" {
		countArgs = []interface{}{"%" + query + "%"}
	}
	
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting artists: %w", err)
	}
	
	searchQuery := `
		SELECT id, name, description, formation_year, user_id, created_at, updated_at
		FROM artists
	`
	
	searchArgs := []interface{}{}
	argIndex := 1
	
	if query != "" {
		searchQuery += ` WHERE LOWER(name) LIKE LOWER($` + fmt.Sprintf("%d", argIndex) + `) `
		searchArgs = append(searchArgs, "%"+query+"%")
		argIndex++
	}
	
	searchQuery += ` ORDER BY name ASC LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)
	searchArgs = append(searchArgs, limit, offset)
	
	rows, err := r.db.QueryContext(ctx, searchQuery, searchArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("error searching artists: %w", err)
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
			return nil, 0, fmt.Errorf("error scanning artist: %w", err)
		}
		artists = append(artists, artist)
	}
	
	return artists, total, nil
}

func (r *ArtistRepository) GetByID(ctx context.Context, id uint) (*model.Artist, error) {
	query := `
		SELECT id, name, description, formation_year, user_id, created_at, updated_at
		FROM artists
		WHERE id = $1
	`
	
	var artist model.Artist
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
			return nil, fmt.Errorf("artist not found with id %d", id)
		}
		return nil, fmt.Errorf("error getting artist by id: %w", err)
	}
	
	return &artist, nil
}

func (r *ArtistRepository) GetWithAlbums(ctx context.Context, id uint) (*model.Artist, error) {
	// Сначала получаем артиста
	artist, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Затем получаем альбомы артиста
	albumsQuery := `
		SELECT id, title, artist_id, release_date, cover_art_url, created_at, updated_at
		FROM albums
		WHERE artist_id = $1
		ORDER BY release_date DESC
	`
	
	rows, err := r.db.QueryContext(ctx, albumsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error getting albums for artist: %w", err)
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
			return nil, fmt.Errorf("error scanning album: %w", err)
		}
		albums = append(albums, album)
	}
	
	artist.Albums = albums
	return artist, nil
}