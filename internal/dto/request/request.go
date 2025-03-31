package request

// AddSongRequest представляет запрос на добавление песни
// swagger:model
type AddSongRequest struct {
    // Название группы (обязательное)
    Group string `json:"group" binding:"required" example:"The Beatles"`
    
    // Название песни (обязательное)
    Song  string `json:"song" binding:"required" example:"Yesterday"`
    
    // Текст песни (обязательное)
    Text  string `json:"text" binding:"required" example:"Yesterday, all my troubles seemed so far away..."`
    
    // Дата релиза в формате dd.mm.yyyy (обязательное)
    ReleaseDate string `json:"release_date" binding:"required,datetime=02.01.2006" example:"06.08.1965"`
    
    // Ссылка на песню (обязательное)
    Link string `json:"link" binding:"required" example:"https://example.com/yesterday"`
}

// UpdateSongRequest представляет запрос на обновление песни
// swagger:model
type UpdateSongRequest struct {
    // Новое название группы
    Group string `json:"group" example:"The Beatles (Remastered)"`
    
    // Новое название песни
    Song  string `json:"song" example:"Yesterday (Remastered)"`
    
    // Обновленный текст песни
    Text  string `json:"text" example:"Updated lyrics text..."`
    
    // Новая дата релиза в формате dd.mm.yyyy
    ReleaseDate string `json:"release_date" example:"01.01.2023"`
    
    // Новая ссылка на песню
    Link string `json:"link" example:"https://example.com/yesterday-remastered"`
}

// GetSongRequest представляет запрос на получение песни
// swagger:model
type GetSongRequest struct {
    // Название группы (обязательное)
    Group string `json:"group" binding:"required" example:"Queen"`
    
    // Название песни (обязательное)
    Song  string `json:"song" binding:"required" example:"Bohemian Rhapsody"`
}

// GetSongsRequest представляет запрос на получение песен группы
// swagger:model
type GetSongsRequest struct {
    // Название группы (обязательное)
    Group string `json:"group" binding:"required" example:"Queen"`
}

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
