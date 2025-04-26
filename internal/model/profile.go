package model

import "time"

// Профиль
type Profile struct {
	UserID      uint         `gorm:"primaryKey" json:"user_id"`
	Bio         string       `json:"bio"`
	AvatarURL   string       `json:"avatar_url"`
	Favorites   []Favorite   `gorm:"foreignKey:ProfileID"`
	Collections []Collection `gorm:"foreignKey:ProfileID"`
	History     []History    `gorm:"foreignKey:ProfileID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Избранное
type Favorite struct {
	ID         uint   `gorm:"primaryKey"`
	ProfileID  uint   `gorm:"index;not null"`
	ObjectType string `gorm:"not null"` // song, artist, album
	ObjectID   uint   `gorm:"not null"`
	CreatedAt  time.Time
}

// Коллекция
type Collection struct {
	ID          uint   `gorm:"primaryKey"`
	ProfileID   uint   `gorm:"index;not null"`
	Name        string `gorm:"not null"`
	Description string
	Items       []CollectionItem `gorm:"foreignKey:CollectionID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Элемент коллекции
type CollectionItem struct {
	CollectionID uint `gorm:"primaryKey"`
	SongID       uint `gorm:"primaryKey"`
	Position     int  `gorm:"not null"`
}

// История прослушиваний
type History struct {
	ID        uint      `gorm:"primaryKey"`
	ProfileID uint      `gorm:"index;not null"`
	SongID    uint      `gorm:"index;not null"`
	PlayedAt  time.Time `gorm:"index"`
}
