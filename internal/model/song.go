package model

import (
	"gorm.io/gorm"
)

type Song struct {
	gorm.Model  `json:"-" swaggerignore:"true"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	Text        string `json:"text"`
	ReleaseDate string `json:"release_date"`
	Link        string `json:"link"`
	Owner       string `json:"owner"`
}
