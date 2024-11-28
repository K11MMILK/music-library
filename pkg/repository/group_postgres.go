package repository

import (
	"fmt"
	"strings"
	musiclibrary "time-tracker"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type GroupPostgres struct {
	db *sqlx.DB
}

func NewGroupPostgres(db *sqlx.DB) *GroupPostgres {
	return &GroupPostgres{db: db}
}
func (r *GroupPostgres) CreateGroup(group musiclibrary.Group) (int, error) {
	logrus.Debug("Creating group")
	var id int
	query := fmt.Sprintf("INSERT INTO %s (groupName) VALUES ($1) RETURNING id", groupsTable)
	row := r.db.QueryRow(query, group.GroupName)
	if err := row.Scan(&id); err != nil {
		logrus.WithError(err).Error("Failed to create group")
		return 0, err
	}
	logrus.WithField("id", id).Info("Group created successfully")
	return id, nil
}

func (r *GroupPostgres) GetAllGroups() ([]musiclibrary.Group, error) {
	logrus.Debug("Fetching all groups")
	var groupList []musiclibrary.Group
	query := fmt.Sprintf("SELECT * FROM %s", groupsTable)
	err := r.db.Select(&groupList, query)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch all groups")
		return nil, err
	}
	logrus.WithField("count", len(groupList)).Info("Fetched all groups successfully")
	return groupList, err
}

func (r *GroupPostgres) DeleteGroup(id int) error {
	logrus.WithField("id", id).Debug("Deleting group")
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", groupsTable)
	_, err := r.db.Exec(query, id)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete group")
		return err
	}
	logrus.WithField("id", id).Info("Group deleted successfully")
	return nil
}

func (r *GroupPostgres) UpdateGroup(id int, input musiclibrary.UpdateGroupInput) error {
	logrus.WithField("id", id).Debug("Updating group")
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.GroupName != nil {
		setValues = append(setValues, fmt.Sprintf("groupName=$%d", argId))
		args = append(args, *input.GroupName)
		argId++
	}

	if argId > 1 {
		setQuery := strings.Join(setValues, ", ")
		query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", groupsTable, setQuery, argId)
		args = append(args, id)
		_, err := r.db.Exec(query, args...)
		if err != nil {
			logrus.WithError(err).Error("Failed to update group")
			return err
		}
		logrus.WithField("id", id).Info("Group updated successfully")
	}

	return nil
}

func (r *GroupPostgres) GetGroupsWithFilter(filters map[string]string, page, limit int) ([]musiclibrary.Group, error) {
	var groups []musiclibrary.Group
	var conditions []string
	var args []interface{}
	argId := 1

	query := `SELECT * FROM groupss WHERE 1=1`

	if group, ok := filters["groupname"]; ok && group != "" {
		conditions = append(conditions, fmt.Sprintf("groupname ILIKE $%d", argId))
		args = append(args, "%"+group+"%")
		argId++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argId, argId+1)
	args = append(args, limit, (page-1)*limit)

	err := r.db.Select(&groups, query, args...)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
