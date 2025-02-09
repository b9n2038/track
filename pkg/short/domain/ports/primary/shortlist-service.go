// pkg/short/domain/ports/service.go
package primary

import "act/pkg/short/domain/model"

type ListService interface {
	CreateList(name string, config model.Config) error
	AddItem(listName, item string) error
	MoveToOpen(listName string, index int) error
	MoveToClosed(listName string, index int) error
	GetList(listName string) (*model.ShortList, error)
	UpdateConfig(listName string, config model.Config) error
}
