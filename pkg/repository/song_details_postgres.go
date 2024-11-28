package repository

import (
	"fmt"
	"strings"
	musiclibrary "time-tracker"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SongDetailPostgres struct {
	db *sqlx.DB
}

func NewSongDetailsPostgres(db *sqlx.DB) *SongDetailPostgres {
	return &SongDetailPostgres{db: db}
}

func (r *SongDetailPostgres) GetSongDetailsById(songId int) ([]musiclibrary.SongDetails, error) {
	logrus.WithField("songId", songId).Debug("Fetching song details by song ID")
	var details []musiclibrary.SongDetails
	query := fmt.Sprintf("SELECT id, songId  AS \"songid\", releaseDate AS \"releasedate\", text, link FROM %s WHERE songid = $1", songDetailsTable)
	err := r.db.Select(&details, query, songId)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch song details by song ID")
		return nil, err
	}
	logrus.WithField("count", len(details)).Info("Fetched song details by song ID successfully")
	return details, err
}

func (r *SongDetailPostgres) UpdateSongDetails(id int, input musiclibrary.UpdateSongDetailsInput) error {
	logrus.WithField("id", id).Debug("Updating song detail")
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.ReleaseDate != "" {
		setValues = append(setValues, fmt.Sprintf("releasedate=$%d", argId))
		args = append(args, input.ReleaseDate)
		argId++
	}

	if input.Text != "" {
		setValues = append(setValues, fmt.Sprintf("text=$%d", argId))
		args = append(args, input.Text)
		argId++
	}

	if input.Link != "" {
		setValues = append(setValues, fmt.Sprintf("link=$%d", argId))
		args = append(args, input.Link)
		argId++
	}

	if argId > 1 {
		setQuery := strings.Join(setValues, ", ")
		query := fmt.Sprintf("UPDATE %s SET %s WHERE songid = $%d", songDetailsTable, setQuery, argId)
		args = append(args, id)
		_, err := r.db.Exec(query, args...)
		if err != nil {
			logrus.WithError(err).Error("Failed to update song detail")
			return err
		}
		logrus.WithField("id", id).Info("Song detail updated successfully")
	}
	return nil
}
func (r *SongDetailPostgres) GetSongText(songId int, page int, limit int) ([]string, error) {
	logrus.WithField("songId", songId).Debug("Fetching song text by song ID")
	var details musiclibrary.SongDetailsT
	query := fmt.Sprintf("SELECT id, songId AS \"songid\", text FROM %s WHERE songid = $1", songDetailsTable)

	err := r.db.Get(&details, query, songId)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch song text by song ID")
		return nil, err
	}

	verses := strings.Split(details.Text, "\n\n")

	start := (page - 1) * limit
	end := start + limit

	if start >= len(verses) {
		return nil, nil
	}

	if end > len(verses) {
		end = len(verses)
	}

	paginatedVerses := verses[start:end]

	logrus.WithFields(logrus.Fields{
		"songId": songId,
		"page":   page,
		"limit":  limit,
	}).Info("Fetched song details with pagination successfully")

	return paginatedVerses, nil
}
