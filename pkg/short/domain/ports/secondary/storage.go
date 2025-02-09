// pkg/short/domain/ports/storage.go
package secondary

import "act/pkg/short/domain/model"

type ListStorage interface {
	Save(list *model.ShortList) error
	Load(name string) (*model.ShortList, error)
	Exists(name string) bool
}
