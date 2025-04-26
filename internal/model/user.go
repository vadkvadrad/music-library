package model

import (
	"math/rand"

	"gorm.io/gorm"
)

type role string

const (
	RoleUser  role = "user"
	RoleAdmin role = "admin"
)

// Пользователь
type User struct {
	gorm.Model
	Name       string  `gorm:"uniqueIndex;not null" json:"name"`
	Email      string  `gorm:"uniqueIndex;not null" json:"email"`
	Password   string  `gorm:"not null" json:"-"`
	Role       role    `gorm:"type:varchar(20);default:'user'" json:"role"`
	SessionId  string  `gorm:"index" json:"session_id"`
	Code       string  `json:"code"`
	IsVerified bool    `gorm:"default:false" json:"is_verified"`
	Profile    Profile `gorm:"foreignKey:UserID"`
	Artist     Artist  `gorm:"foreignKey:UserID"`
}

func (u *User) Generate() {
	u.SessionId = randLettersRunes(10)
	u.Code = randNumbersRunes(4)
}

var lettersRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numbersRunes = []rune("0123456789")

func randLettersRunes(n int) string {
	b := make([]rune, n)
	for i := range n {
		b[i] = lettersRunes[rand.Intn(len(lettersRunes))]
	}
	return string(b)
}

func randNumbersRunes(n int) string {
	b := make([]rune, n)
	for i := range n {
		b[i] = numbersRunes[rand.Intn(len(numbersRunes))]
	}
	return string(b)
}
