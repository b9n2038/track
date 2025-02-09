// pkg/domain/short/shortlist.go
package domain

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
