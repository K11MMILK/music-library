package service

import (
	musiclibrary "time-tracker"
	"time-tracker/pkg/repository"
)

type Group interface {
	CreateGroup(group musiclibrary.Group) (int, error)
	GetAllGroups() ([]musiclibrary.Group, error)
	DeleteGroup(id int) error
	UpdateGroup(id int, input musiclibrary.UpdateGroupInput) error
	GetGroupsWithFilter(map[string]string, int, int) ([]musiclibrary.Group, error)
}

type Song interface {
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

type Service struct {
	Group
	Song
	SongDetails
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Group:       NewGroupService(repos.Group),
		Song:        NewAuthService(repos.Authorisation),
		SongDetails: NewSongDetailsService(repos.SongDetails),
	}
}
