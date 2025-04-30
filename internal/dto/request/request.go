package request

type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"qwerty123"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"strong-password"`
}

type VerifyRequest struct {
	SessionId string `json:"session_id" binding:"required" example:"UexEJzPJ3M"`
	Code      string `json:"code" binding:"required" example:"1234"`
}

type NewProfileRequest struct {
	Bio       string `json:"bio" binding:"required" example:"Hello, i am new artist, gonna make songs fo u"`
	AvatarURL string `json:"avatar_url" binding:"required" example:"http://url/image/avatar.png"`
}

type NewArtistRequest struct {
	ArtistName    string `json:"artist_name" example:"Imagine Dragons"`
	Description   string `json:"description" example:"We are team Imagine Dragons, writing songs for you"`
	FormationYear string `json:"formation_year" example:"2001-02-16"` // Формат: "YYYY-MM-DD"
}

type NewAlbumRequest struct {
	Title       string `json:"title" binding:"required"`
	ReleaseDate string `json:"release_date" binding:"required"`
	CoverArtURL string `json:"cover_art_url" binding:"required"`
}

type NewSongRequest struct {
	Title    string    `json:"title" binding:"required"`
	Genres   []Genres  `json:"genres" binding:"required"`
	Duration int       `json:"duration_sec" binding:"required"`
	FilePath string    `json:"file_path" binding:"required"`
	Lyrics   AddLyrics `json:"lyrics" binding:"required"`
}

type Genres struct {
	GenreID uint `json:"genre_id" binding:"required"`
}

type AddLyrics struct {
	Text []Couplet `json:"text" binding:"required"`
}

type Couplet struct {
	Text string `json:"couplet" binding:"required"`
}

type NewGenreRequest struct {
	Name string `json:"genre_name" binding:"required"`
}
