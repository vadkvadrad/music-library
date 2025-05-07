package response

import "time"

type PaginatedResponse struct {
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}

type SearchResult map[string]PaginatedResponse

// Для ответов с деталями артиста
type ArtistDTO struct {
	ID            uint       `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	FormationYear time.Time  `json:"formation_year"`
	Albums        []AlbumDTO `json:"albums,omitempty"`
}

// Для ответов с деталями альбома
type AlbumDTO struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"` // Формат: "2006-01-02"
	CoverArtURL string    `json:"cover_art_url"`
	Songs       []SongDTO `json:"songs,omitempty"`
}

// Для ответов с деталями песни
type SongDTO struct {
	ID       uint      `json:"id"`
	Title    string    `json:"title"`
	AlbumID  uint      `json:"album_id"`
	Duration int       `json:"duration"`
	FilePath string    `json:"file_path"` // или URL для скачивания
	Lyrics   LyricsDTO `json:"lyrics,omitempty"`
}

// Для ответов с текстом песни
type LyricsDTO struct {
	Couplets []CoupletDTO `json:"text"`
}

type CoupletDTO struct {
	Number  uint   `json:"number"`
	Couplet string `json:"couplet"`
}

type LoginResponse struct {
	Token string `json:"jwt_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type RegisterResponse struct {
	SessionId string `json:"session_id" example:"UexEJzPJ3M"`
}

type VerifyResponse struct {
	Token string `json:"jwt_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type SearchErrorResponse struct {
	Error error  `json:"error"`
	Tip   string `json:"tip"`
}

type AddSongResponse struct {
	AlbumID   uint   `json:"added_to_album"`
	AlbumName string `json:"album_name"`
	ArtistID  uint   `json:"artist_id"`
}
