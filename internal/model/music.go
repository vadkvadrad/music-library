package model

import (
	"time"
)

// Артист
type Artist struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"index;not null"`
	Description   string
	FormationYear time.Time
	Albums        []Album `gorm:"foreignKey:ArtistID"`
	UserID        uint    `gorm:"index;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Альбом
type Album struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"index"`
	ArtistID    uint   `gorm:"index"`
	Songs       []Song `gorm:"foreignKey:AlbumID"`
	ReleaseDate time.Time
	CoverArtURL string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Песня
type Song struct {
	ID         uint        `gorm:"primaryKey"`
	Title      string      `gorm:"index"`
	ArtistID   uint        `gorm:"index"`
	AlbumID    uint        `gorm:"index"`
	SongGenres []SongGenre `gorm:"foreignKey:SongID"`
	Duration   int
	FilePath   string
	Lyrics     Lyrics `gorm:"foreignKey:SongID;references:ID;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Текст песни
type Lyrics struct {
	SongID    uint      `gorm:"primaryKey;autoIncrement:false"`
	Couplets  []Couplet `gorm:"foreignKey:LyricsID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Куплеты песни
type Couplet struct {
	ID       uint `gorm:"primaryKey"`
	LyricsID uint `gorm:"index"`
	Number   uint
	Text     string
}

// Жанр
type Genre struct {
	ID         uint        `gorm:"primaryKey"`
	Name       string      `gorm:"uniqueIndex"`
	SongGenres []SongGenre `gorm:"foreignKey:GenreID"` // Связь один-ко-многим
}

// Промежуточная таблица
type SongGenre struct {
	ID      uint  `gorm:"primaryKey"`
	SongID  uint  `gorm:"index"`
	GenreID uint  `gorm:"index"`
	Song    Song  `gorm:"foreignKey:SongID"`
	Genre   Genre `gorm:"foreignKey:GenreID"`
}
