// internal/domain/short/config.go
package domain

type LimitHandling string

const (
	MoveLastToClosed LimitHandling = "moveLastToClosed"
	PushFront        LimitHandling = "pushFront"
)

type Config struct {
	MaxCount      int           `yaml:"maxCount"`
	LimitHandling LimitHandling `yaml:"limitHandling"`
}

func DefaultConfig() Config {
	return Config{
		MaxCount:      3,
		LimitHandling: MoveLastToClosed,
	}
}
