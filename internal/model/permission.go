package model

import "time"

type Resource string
type Permission string

const (
	SongResource   Resource = "song"
	AlbumResource  Resource = "album"
	ArtistResource Resource = "artist"

	EditPermission Permission = "edit"
	ViewPermission Permission = "view"
)

type ResourcePermission struct {
	ID           uint       `gorm:"primaryKey"`
	UserID       uint       `gorm:"index"`
	ResourceID   uint       `gorm:"index"`
	ResourceType Resource   `gorm:"text"`
	Permission   Permission `gorm:"text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
