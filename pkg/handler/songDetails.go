package handler

import (
	"net/http"
	"strconv"
	musiclibrary "time-tracker"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary GetSongDetailsById
// @Tags songDetails
// @Description Get song details by song ID
// @ID get-songDetails-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Success 200 {object} songDetailsByIdResponse "Song details data"
// @Failure 400 {object} errorResponse "Invalid song ID"
// @Failure 404 {object} errorResponse "Song not found"
// @Failure 500 {object} errorResponse "Failed to get song details"
// @Router /api/songDetails/{id} [get]
func (h *Handler) getSongDetailsById(c *gin.Context) {
	logrus.Debug("getSongDetailsById handler called")

	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	songDetails, err := h.services.SongDetails.GetSongDetailsById(songId)
	if err != nil {
		logrus.WithError(err).Error("Failed to get songDetails by ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get songDetails by ID"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id": songId,
		"count":   len(songDetails),
	}).Info("Retrieved songDetails by ID successfully")

	var songDetailsDL []musiclibrary.SongDetailsDL
	for _, element := range songDetails {
		songDetailsDL = append(songDetailsDL, musiclibrary.SongDetailsDL{Id: element.Id, ReleaseDate: element.ReleaseDate, Link: element.Link, SongId: element.SongId})
	}
	c.JSON(http.StatusOK, songDetailsByIdResponse{
		Data: songDetailsDL,
	})
}

// @Summary UpdateSongDetails
// @Tags songDetails
// @Description Update song details by songDetails ID
// @ID update-songDetails
// @Accept  json
// @Produce  json
// @Param id path int true "SongDetails ID"
// @Param input body musiclibrary.UpdateSongDetailsInput true "SongDetails info"
// @Success 200 {object} statusResponse "Status of the operation"
// @Failure 400 {object} errorResponse "Invalid input or ID"
// @Failure 404 {object} errorResponse "SongDetails not found"
// @Failure 500 {object} errorResponse "Failed to update song details"
// @Router /api/songDetails/{id} [put]
func (h *Handler) updateSongDetails(c *gin.Context) {
	logrus.Debug("updateSongDetails handler called")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Invalid songDetails ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid songDetails ID"})
		return
	}

	var input musiclibrary.UpdateSongDetailsInput
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("Failed to bind JSON for update songDetails")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = h.services.SongDetails.UpdateSongDetails(id, input)
	if err != nil {
		logrus.WithError(err).Error("Failed to update songDetails")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update songDetails"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"songDetails_id": id,
	}).Info("SongDetails updated successfully")

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary GetSongText
// @Tags songDetails
// @Description Get song text with pagination by song ID
// @ID get-song-text
// @Accept  json
// @Produce  json
// @Param id path int true "SongDetails ID"
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Limit of verses per page" default(10)
// @Success 200 {object} songTextResponse "Song text with pagination"
// @Failure 400 {object} errorResponse "Invalid songDetails ID or pagination parameters"
// @Failure 500 {object} errorResponse "Failed to get songDetails"
// @Router /api/songText/{id}/filter [get]
func (h *Handler) getSongText(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Invalid songDetails ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid songDetails ID"})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	songText, err := h.services.SongDetails.GetSongText(id, page, limit)

	if err != nil {
		logrus.WithError(err).Error("Failed to get songDetails by ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get songDetails by ID"})
		return
	}
	c.JSON(http.StatusOK, songTextResponse{
		Data: songText,
	})
}

type songDetailsByIdResponse struct {
	Data []musiclibrary.SongDetailsDL `json:"data"`
}
type songTextResponse struct {
	Data []string `json:"data"`
}
