package model

import "time"

// Профиль
type Profile struct {
	UserID      uint         `json:"user_id"`
	Bio         string       `json:"bio"`
	AvatarURL   string       `json:"avatar_url"`
	Favorites   []Favorite   `json:"favorites"`
	Collections []Collection `json:"collections"`
	History     []History    `json:"history"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Избранное
type Favorite struct {
	ID         uint      `json:"id"`
	ProfileID  uint      `json:"profile_id"`
	ObjectType string    `json:"object_type"` // song, artist, album
	ObjectID   uint      `json:"object_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// Коллекция
type Collection struct {
	ID          uint             `json:"id"`
	ProfileID   uint             `json:"profile_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Items       []CollectionItem `json:"items"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// Элемент коллекции
type CollectionItem struct {
	CollectionID uint `json:"collection_id"`
	SongID       uint `json:"song_id"`
	Position     int  `json:"position"`
}

// История прослушиваний
type History struct {
	ID        uint      `json:"id"`
	ProfileID uint      `json:"profile_id"`
	SongID    uint      `json:"song_id"`
	PlayedAt  time.Time `json:"played_at"`
}
