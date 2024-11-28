package service

import (
	timetracker "time-tracker"
	"time-tracker/pkg/repository"
)

type AuthServise struct {
	repo repository.Authorisation
}

func NewAuthService(repo repository.Authorisation) *AuthServise {
	return &AuthServise{repo: repo}
}

func (s *AuthServise) CreateSong(song timetracker.Song) (int, error) {
	return s.repo.CreateSong(song)
}

func (s *AuthServise) GetAllSongs() ([]timetracker.Song, error) {
	return s.repo.GetAllSongs()
}

func (s *AuthServise) DeleteSong(id int) error {
	return s.repo.DeleteSong(id)
}

func (s *AuthServise) UpdateSong(id int, input timetracker.UpdateSongInput) error {
	return s.repo.UpdateSong(id, input)
}
func (s *AuthServise) GetSongsWithFilter(filters map[string]string, page int, limit int) ([]timetracker.Song, error) {
	return s.repo.GetSongsWithFilter(filters, page, limit)
}
