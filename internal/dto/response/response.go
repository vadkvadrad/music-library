package response

type GetSongResponse struct {
	ReleaseDate string `json:"release_date" example:"06.08.1965"`
	Text        string `json:"text" example:"Yesterday, all my troubles seemed so far away..."`
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
