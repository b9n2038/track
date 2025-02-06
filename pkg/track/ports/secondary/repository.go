// internal/ports/secondary/repository.go
package secondary

import (
	"context"
	"act/pkg/track/domain"
	"time"
)

type RatingRepository interface {
	Save(ctx context.Context, r rating.DayRating) error
	GetByID(ctx context.Context, id string) (rating.DayRating, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]rating.DayRating, error)
	GetByWeek(ctx context.Context, year, week int) ([]rating.DayRating, error)
}
