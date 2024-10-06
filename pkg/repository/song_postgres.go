package repository

import (
	"fmt"
	"strings"
	musiclibrary "time-tracker"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SongPostgres struct {
	db *sqlx.DB
}

func NewSongPostgres(db *sqlx.DB) *SongPostgres {
	return &SongPostgres{db: db}
}

func (r *SongPostgres) CreateSong(song musiclibrary.Song) (int, error) {
	logrus.Debug("Creating song")
	var id int
	query := fmt.Sprintf("INSERT INTO %s (\"group\", song) VALUES ($1, $2) RETURNING id", songsTable)
	row := r.db.QueryRow(query, song.Group, song.Song)
	if err := row.Scan(&id); err != nil {
		logrus.WithError(err).Error("Failed to create song")
		return 0, err
	}
	query = fmt.Sprintf("INSERT INTO %s (songId) VALUES ($1)", songDetailsTable)
	r.db.QueryRow(query, id)
	logrus.WithField("id", id).Info("Song created successfully")
	return id, nil
}

func (r *SongPostgres) GetAllSongs() ([]musiclibrary.Song, error) {
	logrus.Debug("Fetching all songs")
	var songList []musiclibrary.Song
	query := fmt.Sprintf("SELECT * FROM %s", songsTable)
	err := r.db.Select(&songList, query)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch all songs")
		return nil, err
	}
	logrus.WithField("count", len(songList)).Info("Fetched all songs successfully")
	return songList, err
}

func (r *SongPostgres) DeleteSong(id int) error {
	logrus.WithField("id", id).Debug("Deleting song")
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", songsTable)
	_, err := r.db.Exec(query, id)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete song")
		return err
	}
	logrus.WithField("id", id).Info("Song deleted successfully")
	return nil
}

func (r *SongPostgres) UpdateSong(id int, input musiclibrary.UpdateSongInput) error {
	logrus.WithField("id", id).Debug("Updating song")
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Group != nil {
		setValues = append(setValues, fmt.Sprintf("\"group\"=$%d", argId))
		args = append(args, *input.Group)
		argId++
	}

	if input.Song != nil {
		setValues = append(setValues, fmt.Sprintf("song=$%d", argId))
		args = append(args, *input.Song)
		argId++
	}

	if argId > 1 {
		setQuery := strings.Join(setValues, ", ")
		query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", songsTable, setQuery, argId)
		args = append(args, id)
		_, err := r.db.Exec(query, args...)
		if err != nil {
			logrus.WithError(err).Error("Failed to update song")
			return err
		}
		logrus.WithField("id", id).Info("Song updated successfully")
	}

	return nil
}

func (r *SongPostgres) GetSongsWithFilter(filters map[string]string, page, limit int) ([]musiclibrary.Song, error) {
	var songs []musiclibrary.Song
	var conditions []string
	var args []interface{}
	argId := 1

	query := `SELECT * FROM songs WHERE 1=1`

	if group, ok := filters["group"]; ok && group != "" {
		conditions = append(conditions, fmt.Sprintf("group=$%d", argId))
		args = append(args, group)
		argId++
	}

	if song, ok := filters["song"]; ok && song != "" {
		conditions = append(conditions, fmt.Sprintf("song=$%d", argId))
		args = append(args, song)
		argId++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argId, argId+1)
	args = append(args, limit, (page-1)*limit)

	err := r.db.Select(&songs, query, args...)
	if err != nil {
		return nil, err
	}

	return songs, nil
}
