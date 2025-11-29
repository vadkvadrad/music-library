package model

import (
	"math/rand"
	"time"
)

type role string

const (
	RoleUser  role = "user"
	RoleAdmin role = "admin"
)

// Пользователь
type User struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Password   string     `json:"-"`
	Role       role       `json:"role"`
	SessionId  string     `json:"session_id"`
	Code       string     `json:"code"`
	IsVerified bool       `json:"is_verified"`
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
