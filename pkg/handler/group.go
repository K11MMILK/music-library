package handler

import (
	"net/http"
	"strconv"
	musiclibrary "time-tracker"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary CreateGroup
// @Tags group
// @Description Create a new group
// @ID create-group
// @Accept  json
// @Produce  json
// @Param input body musiclibrary.Group true "Group information"
// @Success 200 {object} map[string]interface{} "Returns group ID"
// @Failure 400 {object} errorResponse "Invalid input"
// @Failure 500 {object} errorResponse "Failed to create group"
// @Router /api/group/ [post]
func (h *Handler) createGroup(c *gin.Context) {
	var group musiclibrary.Group
	if err := c.BindJSON(&group); err != nil {
		logrus.WithError(err).Error("Failed to bind JSON for sign up")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id, err := h.services.Group.CreateGroup(group)
	if err != nil {
		logrus.WithError(err).Error("Failed to create group")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"group_id": id,
	}).Info("group created successfully")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary GetAllGroups
// @Tags group
// @Description Get a list of all groups
// @ID get-all-groups
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllGroupsResponse "Returns a list of all groups"
// @Failure 500 {object} errorResponse "Failed to get all groups"
// @Router /api/group/ [get]
func (h *Handler) getAllGroups(c *gin.Context) {
	groupList, err := h.services.Group.GetAllGroups()
	if err != nil {
		logrus.WithError(err).Error("Failed to get all groups")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all groups"})
		return
	}

	logrus.Info("Retrieved all groups successfully")

	c.JSON(http.StatusOK, getAllGroupsResponse{
		Data: groupList,
	})
}

// @Summary UpdateGroup
// @Tags group
// @Description Update an existing group
// @ID update-group
// @Accept  json
// @Produce  json
// @Param id path int true "Group ID"
// @Param input body musiclibrary.UpdateGroupInput true "Group information"
// @Success 200 {object} statusResponse "Returns status of the operation"
// @Failure 400 {object} errorResponse "Invalid input or ID"
// @Failure 500 {object} errorResponse "Failed to update group"
// @Router /api/group/{id} [put]
func (h *Handler) updateGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Invalid group ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var input musiclibrary.UpdateGroupInput
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("Failed to bind JSON for update group")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = h.services.Group.UpdateGroup(id, input)
	if err != nil {
		logrus.WithError(err).Error("Failed to update group")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"group_id": id,
	}).Info("group updated successfully")

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary DeleteGroup
// @Tags group
// @Description Delete an existing group
// @ID delete-group
// @Accept  json
// @Produce  json
// @Param id path int true "Group ID"
// @Success 200 {object} statusResponse "Returns status of the operation"
// @Failure 400 {object} errorResponse "Invalid group ID"
// @Failure 500 {object} errorResponse "Failed to delete group"
// @Router /api/group/{id} [delete]
func (h *Handler) deleteGroup(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Invalid group ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	err = h.services.Group.DeleteGroup(id)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete group")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"group_id": id,
	}).Info("group deleted successfully")

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary GetGroupsWithFilter
// @Tags group
// @Description Get groups with optional filtering and pagination
// @ID get-groups-with-filter
// @Accept  json
// @Produce  json
// @Param groupname query string false "Group name filter"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of groups per page"
// @Success 200 {object} getAllGroupsResponse "Returns a filtered list of groups"
// @Failure 500 {object} errorResponse "Failed to get groups"
// @Router /api/group/filter [get]
func (h *Handler) getGroupsWithFilter(c *gin.Context) {
	filters := map[string]string{
		"groupname": c.Query("groupname"),
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	groups, err := h.services.Group.GetGroupsWithFilter(filters, page, limit)
	if err != nil {
		logrus.WithError(err).Error("Failed to get groups with filters")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get groups"})
		return
	}

	c.JSON(http.StatusOK, getAllGroupsResponse{
		Data: groups,
	})
}

type getAllGroupsResponse struct {
	Data []musiclibrary.Group `json:"data"`
}
