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
	ID           uint       `json:"id"`
	UserID       uint       `json:"user_id"`
	ResourceID   uint       `json:"resource_id"`
	ResourceType Resource   `json:"resource_type"`
	Permission   Permission `json:"permission"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
