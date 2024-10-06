package musiclibrary

type Song struct {
	Id    int    `json:"id" db:"id"`
	Group string `json:"group" db:"group" binding:"required"`
	Song  string `json:"song" db:"song" binding:"required"`
}

type UpdateSongInput struct {
	Group *string `json:"group" example:"Metallica"`
	Song  *string `json:"song" example:"Enter Sandman"`
}

type CreateSongInput struct {
	Group string `json:"group" binding:"required" example:"Metallica"`
	Song  string `json:"song" binding:"required" example:"Enter Sandman"`
}

type UpdateSongDetailsInput struct {
	ReleaseDate string `json:"releaseDate" example:"2024-01-01"`
	Text        string `json:"text" example:"Song lyrics here"`
	Link        string `json:"link" example:"https://example.com/song"`
}
type SongDetails struct {
	Id          int    `json:"id" db:"id"`
	SongId      int    `json:"songId" db:"songId"`
	ReleaseDate string `json:"releaseDate" db:"releaseDate"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
}
type SongDetailsDL struct {
	Id          int    `json:"id" db:"id"`
	SongId      int    `json:"songId" db:"songId"`
	ReleaseDate string `json:"releaseDate" db:"releaseDate"`
	Link        string `json:"link" db:"link"`
}
type SongDetailsT struct {
	Id     int    `json:"id" db:"id"`
	SongId int    `json:"songId" db:"songId"`
	Text   string `json:"text" db:"text"`
}
