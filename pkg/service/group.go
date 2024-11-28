package service

import (
	timetracker "time-tracker"
	"time-tracker/pkg/repository"
)

type GroupServise struct {
	repo repository.Group
}

func NewGroupService(repo repository.Group) *GroupServise {
	return &GroupServise{repo: repo}
}

func (s *GroupServise) CreateGroup(Group timetracker.Group) (int, error) {
	return s.repo.CreateGroup(Group)
}

func (s *GroupServise) GetAllGroups() ([]timetracker.Group, error) {
	return s.repo.GetAllGroups()
}

func (s *GroupServise) DeleteGroup(id int) error {
	return s.repo.DeleteGroup(id)
}

func (s *GroupServise) UpdateGroup(id int, input timetracker.UpdateGroupInput) error {
	return s.repo.UpdateGroup(id, input)
}
func (s *GroupServise) GetGroupsWithFilter(filters map[string]string, page int, limit int) ([]timetracker.Group, error) {
	return s.repo.GetGroupsWithFilter(filters, page, limit)
}
