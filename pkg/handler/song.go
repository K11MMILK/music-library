package handler

import (
	"net/http"
	"strconv"
	musiclibrary "time-tracker"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary CreateSong
// @Tags song
// @Description Create a new song
// @ID create-song
// @Accept  json
// @Produce  json
// @Param input body musiclibrary.CreateSongInput true "Song information"
// @Success 200 {object} map[string]interface{} "Returns song ID"
// @Failure 400 {object} errorResponse "Invalid input"
// @Failure 500 {object} errorResponse "Failed to create song"
// @Router /api/song/ [post]
func (h *Handler) createSong(c *gin.Context) {
	var song musiclibrary.Song
	if err := c.BindJSON(&song); err != nil {
		logrus.WithError(err).Error("Failed to bind JSON for sign up")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id, err := h.services.Song.CreateSong(song)
	if err != nil {
		logrus.WithError(err).Error("Failed to create song")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create song"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id": id,
	}).Info("Song created successfully")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary GetAllSongs
// @Tags song
// @Description Get all songs
// @ID getAllSongs
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllSongsResponse "Returns a list of all songs"
// @Failure 500 {object} errorResponse "Failed to get all songs"
// @Router /api/song/ [get]
func (h *Handler) getAllSongs(c *gin.Context) {
	songList, err := h.services.Song.GetAllSongs()
	if err != nil {
		logrus.WithError(err).Error("Failed to get all songs")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all songs"})
		return
	}

	logrus.Info("Retrieved all songs successfully")

	c.JSON(http.StatusOK, getAllSongsResponse{
		Data: songList,
	})
}

// @Summary UpdateSong
// @Tags song
// @Description Update an existing song
// @ID update-song
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param input body musiclibrary.UpdateSongInput true "Song information"
// @Success 200 {object} statusResponse "Returns status of the operation"
// @Failure 400 {object} errorResponse "Invalid input or ID"
// @Failure 500 {object} errorResponse "Failed to update song"
// @Router /api/song/{id} [put]
func (h *Handler) updateSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var input musiclibrary.UpdateSongInput
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("Failed to bind JSON for update song")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = h.services.Song.UpdateSong(id, input)
	if err != nil {
		logrus.WithError(err).Error("Failed to update song")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id": id,
	}).Info("Song updated successfully")

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary DeleteSong
// @Tags song
// @Description Delete an existing song
// @ID delete-song
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Success 200 {object} statusResponse "Returns status of the operation"
// @Failure 400 {object} errorResponse "Invalid song ID"
// @Failure 500 {object} errorResponse "Failed to delete song"
// @Router /api/song/{id} [delete]
func (h *Handler) deleteSong(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	err = h.services.Song.DeleteSong(id)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete song")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id": id,
	}).Info("Song deleted successfully")

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary GetSongsWithFilter
// @Tags song
// @Description Get songs with filtering
// @ID getSongsWithFilter
// @Accept  json
// @Produce  json
// @Param group query string false "Group filter"
// @Param song query string false "Song filter"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Limit for pagination"
// @Success 200 {object} getAllSongsResponse
// @Failure 500 {object} errorResponse
// @Router /api/song/filter [get]
func (h *Handler) getSongsWithFilter(c *gin.Context) {
	filters := map[string]string{
		"songname":    c.Query("songname"),
		"releasedate": c.Query("releasedate"),
		"link":        c.Query("link"),
		"text":        c.Query("text"),
		"groupname":   c.Query("groupname"),
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	songs, err := h.services.Song.GetSongsWithFilter(filters, page, limit)
	if err != nil {
		logrus.WithError(err).Error("Failed to get songs with filters")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get songs"})
		return
	}

	c.JSON(http.StatusOK, getAllSongsResponse{
		Data: songs,
	})
}

type getAllSongsResponse struct {
	Data []musiclibrary.Song `json:"data"`
}
