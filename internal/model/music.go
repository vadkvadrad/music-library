package model

import (
	"time"
)

// Артист
type Artist struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	FormationYear time.Time `json:"formation_year"`
	Albums        []Album   `json:"albums"`
	UserID        uint      `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Альбом
type Album struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	ArtistID    uint      `json:"artist_id"`
	Songs       []Song    `json:"songs"`
	ReleaseDate time.Time `json:"release_date"`
	CoverArtURL string    `json:"cover_art_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Песня
type Song struct {
	ID         uint        `json:"id"`
	Title      string      `json:"title"`
	ArtistID   uint        `json:"artist_id"`
	AlbumID    uint        `json:"album_id"`
	SongGenres []SongGenre `json:"song_genres"`
	Duration   int         `json:"duration"`
	FilePath   string      `json:"file_path"`
	Lyrics     Lyrics      `json:"lyrics"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// Текст песни
type Lyrics struct {
	SongID    uint      `json:"song_id"`
	Couplets  []Couplet `json:"couplets"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Куплеты песни
type Couplet struct {
	ID       uint   `json:"id"`
	LyricsID uint   `json:"lyrics_id"`
	Number   uint   `json:"number"`
	Text     string `json:"text"`
}

// Жанр
type Genre struct {
	ID         uint        `json:"id"`
	Name       string      `json:"name"`
	SongGenres []SongGenre `json:"song_genres"` // Связь один-ко-многим
}

// Промежуточная таблица
type SongGenre struct {
	ID      uint  `json:"id"`
	SongID  uint  `json:"song_id"`
	GenreID uint  `json:"genre_id"`
	Song    Song  `json:"song"`
	Genre   Genre `json:"genre"`
}
