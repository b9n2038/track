// internal/ports/primary/rating/service.go
package rating

import (
	"track/internal/track/domain/rating"
	"context"
	"time"
)

type Service interface {
	CreateDayRating(ctx context.Context, date time.Time, rating rating.Rating) (rating.DayRating, error)
	GetWeekRatings(ctx context.Context, year, week int) ([]rating.DayRating, error)
	GetDateRangeRatings(ctx context.Context, start, end time.Time) ([]rating.DayRating, error)
}
