// pkg/short/domain/model/shortlist.go
package model

import "errors"

var (
	ErrListFull     = errors.New("open list is full")
	ErrInvalidIndex = errors.New("invalid index")
)

type ShortList struct {
	Name   string
	Config Config
	Open   []string
	Closed []string
}

func NewShortList(name string, config Config) *ShortList {
	return &ShortList{
		Name:   name,
		Config: config,
		Open:   make([]string, 0),
		Closed: make([]string, 0),
	}
}

func (s *ShortList) AddToOpen(item string) error {
	if len(s.Open) >= s.Config.MaxCount {
		switch s.Config.LimitHandling {
		case MoveLastToClosed:
			s.Closed = append(s.Closed, s.Open[len(s.Open)-1])
			s.Open = s.Open[:len(s.Open)-1]
		case PushFront:
			s.Closed = append(s.Closed, s.Open[s.Config.MaxCount-1])
			copy(s.Open[1:], s.Open[:s.Config.MaxCount-1])
			s.Open[0] = item
			return nil
		}
	}
	s.Open = append(s.Open, item)
	return nil
}
func (s *ShortList) MoveToOpen(index int) error {
	if index < 0 || index >= len(s.Closed) {
		return ErrInvalidIndex
	}

	// Get the item to move
	item := s.Closed[index]

	// Check if open list is full before moving
	if len(s.Open) >= s.Config.MaxCount {
		switch s.Config.LimitHandling {
		case MoveLastToClosed:
			// Move last open item to closed
			s.Closed = append(s.Closed, s.Open[len(s.Open)-1])
			s.Open = s.Open[:len(s.Open)-1]
		case PushFront:
			// Move last open item to closed
			s.Closed = append(s.Closed, s.Open[s.Config.MaxCount-1])
			// Shift items right
			copy(s.Open[1:], s.Open[:s.Config.MaxCount-1])
			s.Open[0] = item
			// Remove the moved item from closed
			s.Closed = append(s.Closed[:index], s.Closed[index+1:]...)
			return nil
		}
	}

	// Remove from closed list
	s.Closed = append(s.Closed[:index], s.Closed[index+1:]...)
	// Add to open list
	s.Open = append(s.Open, item)

	return nil
}

func (s *ShortList) MoveToClosed(index int) error {
	if index < 0 || index >= len(s.Open) {
		return ErrInvalidIndex
	}

	// Get the item to move
	item := s.Open[index]

	// Remove from open list
	s.Open = append(s.Open[:index], s.Open[index+1:]...)

	// Add to closed list
	s.Closed = append(s.Closed, item)

	return nil
}
