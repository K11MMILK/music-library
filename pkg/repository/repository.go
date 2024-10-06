package repository

import (
	musiclibrary "time-tracker"

	"github.com/jmoiron/sqlx"
)

type Authorisation interface {
	CreateSong(song musiclibrary.Song) (int, error)
	GetAllSongs() ([]musiclibrary.Song, error)
	DeleteSong(id int) error
	UpdateSong(id int, input musiclibrary.UpdateSongInput) error
	GetSongsWithFilter(map[string]string, int, int) ([]musiclibrary.Song, error)
}
type SongDetails interface {
	GetSongDetailsById(songId int) ([]musiclibrary.SongDetails, error)
	UpdateSongDetails(id int, input musiclibrary.UpdateSongDetailsInput) error
	GetSongText(songId int, page int, limit int) ([]string, error)
}

type Repository struct {
	Authorisation
	SongDetails
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorisation: NewSongPostgres(db),
		SongDetails:   NewSongDetailsPostgres(db),
	}
}
