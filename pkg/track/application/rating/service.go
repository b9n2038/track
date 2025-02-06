// internal/application/rating/service.go
package rating

import (
   "context"
   "fmt"
   "time"

   "github/b9n2038/act/pkg/track/domain/rating"
   "github/b9n2038/act/pkg/track/ports/secondary"
)

type Service struct {
   repo secondary.RatingRepository
}

func NewService(repo secondary.RatingRepository) *Service {
   return &Service{
       repo: repo,
   }
}

// AddDayRating creates a new rating for a specific day
func (s *Service) AddDayRating(ctx context.Context, date time.Time, r rating.Rating) (rating.DayRating, error) {
   if !r.IsValid() {
       return rating.DayRating{}, rating.ErrInvalidRating
   }

   // Create new DayRating
   dayRating := rating.DayRating{
       ID:     fmt.Sprintf("%sw%02d-%d", date.Format("06"), date.ISOWeek()), date.Weekday()),
       Date:   date,
       Rating: r,
   }

   // Save to repository
   if err := s.repo.Save(ctx, dayRating); err != nil {
       return rating.DayRating{}, fmt.Errorf("saving day rating: %w", err)
   }

   return dayRating, nil
}

// GetTodayRating gets the rating for the current day
func (s *Service) GetTodayRating(ctx context.Context) (rating.DayRating, error) {
   today := time.Now()
   id := fmt.Sprintf("%sw%02d-%d", today.Format("06"), today.ISOWeek(), today.Weekday())
   
   return s.repo.GetByID(ctx, id)
}

// GetWeekRatings gets all ratings for a specific week
func (s *Service) GetWeekRatings(ctx context.Context, year, week int) ([]rating.DayRating, error) {
   return s.repo.GetByWeek(ctx, year, week)
}

// GetCurrentWeekRatings gets all ratings for the current week
func (s *Service) GetCurrentWeekRatings(ctx context.Context) ([]rating.DayRating, error) {
   now := time.Now()
   year, week := now.ISOWeek()
   return s.GetWeekRatings(ctx, year, week)
}

// GetWeekSummary provides a summary of ratings for a specific week
type WeekSummary struct {
   Year      int
   Week      int
   Average   float64
   Best      rating.DayRating
   Worst     rating.DayRating
   DayCount  int
}

func (s *Service) GetWeekSummary(ctx context.Context, year, week int) (WeekSummary, error) {
   ratings, err := s.repo.GetByWeek(ctx, year, week)
   if err != nil {
       return WeekSummary{}, fmt.Errorf("getting week ratings: %w", err)
   }

   if len(ratings) == 0 {
       return WeekSummary{
           Year: year,
           Week: week,
       }, nil
   }

   var sum int
   best := ratings[0]
   worst := ratings[0]

   for _, dr := range ratings {
       sum += int(dr.Rating)
       
       if dr.Rating > best.Rating {
           best = dr
       }
       if dr.Rating < worst.Rating {
           worst = dr
       }
   }

   return WeekSummary{
       Year:     year,
       Week:     week,
       Average:  float64(sum) / float64(len(ratings)),
       Best:     best,
       Worst:    worst,
       DayCount: len(ratings),
   }, nil
}

// GetDateRangeRatings gets all ratings between start and end dates
func (s *Service) GetDateRangeRatings(ctx context.Context, start, end time.Time) ([]rating.DayRating, error) {
   if end.Before(start) {
       return nil, fmt.Errorf("end date before start date")
   }

   return s.repo.GetByDateRange(ctx, start, end)
}

// UpdateTodayRating updates the rating for the current day
func (s *Service) UpdateTodayRating(ctx context.Context, r rating.Rating) (rating.DayRating, error) {
   if !r.IsValid() {
       return rating.DayRating{}, rating.ErrInvalidRating
   }

   today := time.Now()
   dayRating := rating.DayRating{
       ID:     fmt.Sprintf("%sw%02d-%d", today.Format("06"), today.ISOWeek(), today.Weekday()),
       Date:   today,
       Rating: r,
   }

   if err := s.repo.Save(ctx, dayRating); err != nil {
       return rating.DayRating{}, fmt.Errorf("updating today's rating: %w", err)
   }

   return dayRating, nil
}
