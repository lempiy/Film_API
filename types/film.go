package types

type Film struct {
	ID          int    `json:"id"`
	Name string `json:"name"`
	Year int    `json:"year"`
	AddedAt string `json:"added_at"`
	Genres []Genre `json:"genres, omitempty"`
}

type PostFilm struct {
	*Film
	Genres []int `json:"genres"`
}

type GetFilmParams struct {
	Limit int `query:"limit"`
	Offset int `query:"offset"`
}