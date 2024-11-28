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
	query := fmt.Sprintf("INSERT INTO %s (songName, groupId) VALUES ($1, $2) RETURNING id", songsTable)
	row := r.db.QueryRow(query, song.SongName, song.GroupId)
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

	if input.SongName != nil {
		setValues = append(setValues, fmt.Sprintf("songname=$%d", argId))
		args = append(args, *input.SongName)
		argId++
	}
	if input.GroupId != nil {
		setValues = append(setValues, fmt.Sprintf("groupid=$%d", argId))
		args = append(args, *input.GroupId)
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

	query := `
		SELECT s.* 
		FROM songs s 
		JOIN songDetails sd ON s.id = sd.songId 
		JOIN groupss g ON s.groupId = g.id 
		WHERE 1=1
	`

	if song, ok := filters["songname"]; ok && song != "" {
		conditions = append(conditions, fmt.Sprintf("s.songname ILIKE $%d", argId))
		args = append(args, "%"+song+"%")
		argId++
	}

	if releaseDate, ok := filters["releasedate"]; ok && releaseDate != "" {
		conditions = append(conditions, fmt.Sprintf("TO_CHAR(sd.releaseDate, 'YYYY-MM-DD') ILIKE $%d", argId))
		args = append(args, "%"+releaseDate+"%")
		argId++
	}

	if link, ok := filters["link"]; ok && link != "" {
		conditions = append(conditions, fmt.Sprintf("sd.link ILIKE $%d", argId))
		args = append(args, "%"+link+"%")
		argId++
	}

	if text, ok := filters["text"]; ok && text != "" {
		conditions = append(conditions, fmt.Sprintf("sd.text ILIKE $%d", argId))
		args = append(args, "%"+text+"%")
		argId++
	}

	if groupName, ok := filters["groupname"]; ok && groupName != "" {
		conditions = append(conditions, fmt.Sprintf("g.groupname ILIKE $%d", argId))
		args = append(args, "%"+groupName+"%")
		argId++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY s.id LIMIT $%d OFFSET $%d", argId, argId+1)
	args = append(args, limit, (page-1)*limit)

	err := r.db.Select(&songs, query, args...)
	if err != nil {
		return nil, err
	}

	return songs, nil
}
