// pkg/short/domain/service/list.go
package service

import (
	"act/pkg/short/domain/model"
	"act/pkg/short/domain/ports/primary"
	"act/pkg/short/domain/ports/secondary"
	"fmt"
)

type listService struct {
	storage secondary.ListStorage
}

func NewListService(storage secondary.ListStorage) primary.ListService {
	return &listService{storage: storage}
}

func (s *listService) GetList(name string) (*model.ShortList, error) {
	return s.storage.Load(name)
}

func (s *listService) CreateList(name string, config model.Config) error {
	if s.storage.Exists(name) {
		return fmt.Errorf("list already exists: %s", name)
	}
	list := model.NewShortList(name, config)
	return s.storage.Save(list)
}

func (s *listService) AddItem(listName, item string) error {
	list, err := s.storage.Load(listName)
	if err != nil {
		return err
	}

	if err := list.AddToOpen(item); err != nil {
		return err
	}

	return s.storage.Save(list)
}

func (s *listService) MoveToOpen(listName string, index int) error {
	list, err := s.storage.Load(listName)
	if err != nil {
		return err
	}

	if err := list.MoveToOpen(index); err != nil {
		return err
	}

	return s.storage.Save(list)
}

func (s *listService) MoveToClosed(listName string, index int) error {
	list, err := s.storage.Load(listName)
	if err != nil {
		return err
	}

	if err := list.MoveToClosed(index); err != nil {
		return err
	}

	return s.storage.Save(list)
}
func (s *listService) UpdateConfig(listName string, config model.Config) error {
	list, err := s.storage.Load(listName)
	if err != nil {
		return err
	}

	list.Config = config
	return s.storage.Save(list)
}
