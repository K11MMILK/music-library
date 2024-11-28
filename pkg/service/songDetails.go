package service

import (
	musiclibrary "time-tracker"
	"time-tracker/pkg/repository"
)

type SongDetailsService struct {
	repo repository.SongDetails
}

func NewSongDetailsService(repo repository.SongDetails) *SongDetailsService {
	return &SongDetailsService{repo: repo}
}

func (s *SongDetailsService) GetSongDetailsById(songId int) ([]musiclibrary.SongDetails, error) {
	return s.repo.GetSongDetailsById(songId)
}

func (s *SongDetailsService) UpdateSongDetails(id int, input musiclibrary.UpdateSongDetailsInput) error {
	return s.repo.UpdateSongDetails(id, input)
}
func (s *SongDetailsService) GetSongText(songId int, page int, limit int) ([]string, error) {
	return s.repo.GetSongText(songId, page, limit)
}
