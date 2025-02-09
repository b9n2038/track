package cli

import (
	ratingService "act/pkg/track/application/rating" // aliased this import
	"act/pkg/track/domain/rating"
	// "act/pkg/track/domain/short"
	"strconv"

	//"act/pkg/track/ports/primary/rating"
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func NewRootCmd(service *ratingService.Service) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "rating",
		Short: "Day rating tool",
	}

	rootCmd.AddCommand(
		newAddCmd(service),
		newListCmd(service),
		newWeekCmd(service),
	)

	return rootCmd
}

// Helper function to parse string to int
func parseInt(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return value
}

// Helper function to parse day ID format YYwWW-D
func parseDayID(id string) (time.Time, error) {
	// Parse format: YYwWW-D
	if len(id) != 7 || id[2] != 'w' || id[5] != '-' {
		return time.Time{}, fmt.Errorf("invalid format, expected YYwWW-D")
	}

	year, _ := strconv.Atoi("20" + id[0:2])
	week, _ := strconv.Atoi(id[3:5])
	day, _ := strconv.Atoi(id[6:])

	// Get the date for the start of the week
	date := isoWeekStart(year, week)
	// Add days
	date = date.AddDate(0, 0, day)

	return date, nil
}

// Helper function to get the start of an ISO week
func isoWeekStart(year, week int) time.Time {
	// Find a day in the week
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.Local)
	_, w := jan4.ISOWeek()
	// Move to the desired week
	daysSinceJan4 := (week - w) * 7
	return jan4.AddDate(0, 0, daysSinceJan4-int(jan4.Weekday())+1)
}

func newAddCmd(service *ratingService.Service) *cobra.Command {
	var (
		dayID    string
		fillGaps bool
	)

	cmd := &cobra.Command{
		Use:   "add [rating]",
		Short: "Add a rating between 1 and 5, for today",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			value, err := rating.NewRating(parseInt(args[0])) // using domain package
			if err != nil {
				return err
			}

			if dayID == "" {
				// Add rating for today
				dayRating, err := service.AddDayRating(ctx, time.Now(), value)
				if err != nil {
					return err
				}
				fmt.Printf("Added rating: %s\n", dayRating)
			} else {
				// Add rating for specified day
				date, err := parseDayID(dayID)
				if err != nil {
					return fmt.Errorf("invalid day ID format: %w", err)
				}

				dayRating, err := service.AddDayRating(ctx, date, value)
				if err != nil {
					return err
				}
				fmt.Printf("Added rating for %s: %s\n", dayID, dayRating)

				// Fill gaps if requested
				if fillGaps {
					lastRating, err := service.GetLastRatingBefore(ctx, date)
					if err != nil {
						return fmt.Errorf("getting last rating: %w", err)
					}

					if !lastRating.Date.IsZero() {
						filled, err := service.FillMissingRatings(ctx, lastRating.Date, date, value)
						if err != nil {
							return fmt.Errorf("filling gaps: %w", err)
						}
						if len(filled) > 0 {
							fmt.Printf("Filled %d missing days with rating %s\n", len(filled), value.String())
						}
					}
				}
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&dayID, "day", "d", "", "Day ID in format YYwWW-D (e.g., 25w05-3)")
	cmd.Flags().BoolVarP(&fillGaps, "fill", "f", false, "Fill missing days from last entry")
	return cmd
}

func newListCmd(service *ratingService.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List ratings for current week",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			ratings, err := service.GetCurrentWeekRatings(ctx)
			if err != nil {
				return err
			}

			// Display ratings
			for _, r := range ratings {
				fmt.Printf("%s: %s %s\n", r.Label(), r.Rating.String(), r.Rating.Emoji())
			}
			return nil
		},
	}
}

func newWeekCmd(service *ratingService.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "week",
		Short: "Show ratings for current week with optional trend analysis",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Current week summary and details (existing code)
			year, week := time.Now().ISOWeek()
			summary, err := service.GetWeekSummary(ctx, year, week)
			if err != nil {
				return fmt.Errorf("getting week summary: %w", err)
			}

			// Print summary header
			fmt.Printf("Week %d, %d Summary:\n", summary.Week, summary.Year)
			fmt.Printf("─────────────────────\n")

			if summary.DayCount == 0 {
				fmt.Println("No ratings recorded this week")
				return nil
			}

			// Print stats
			fmt.Printf("Days Rated: %d\n", summary.DayCount)
			fmt.Printf("Average:    %.1f\n", summary.Average)
			fmt.Printf("Best Day:   %s %s\n", summary.Best.Label(), summary.Best.Rating.Emoji())
			fmt.Printf("Worst Day:  %s %s\n", summary.Worst.Label(), summary.Worst.Rating.Emoji())

			// Get detailed ratings for the week
			ratings, err := service.GetWeekRatings(ctx, year, week)
			if err != nil {
				return fmt.Errorf("getting week ratings: %w", err)
			}

			// Print daily breakdown
			fmt.Printf("\nDaily Breakdown:\n")
			fmt.Printf("───────────────\n")
			for _, r := range ratings {
				fmt.Printf("%s: %s %s\n",
					r.Date.Format("Monday"),
					r.Rating.String(),
					r.Rating.Emoji(),
				)
			}
			// Add 13-week trend analysis
			fmt.Printf("\n13-Week Trend:\n")
			fmt.Printf("────���────────\n")

			// Calculate start date (13 weeks ago)
			currentDate := time.Now()
			startDate := currentDate.AddDate(0, 0, -13*7)

			trends, err := service.GetDateRangeRatings(ctx, startDate, currentDate)
			if err != nil {
				return fmt.Errorf("getting trend data: %w", err)
			}

			// Group by week and calculate averages
			weeklyAverages := make(map[int]float64)
			weeklyCount := make(map[int]int)

			for _, t := range trends {
				_, w := t.Date.ISOWeek()
				weeklyAverages[w] += float64(t.Rating)
				weeklyCount[w]++
			}

			// Print trend with arrow indicators
			var lastAvg float64
			for i := 0; i < 13; i++ {
				weekDate := currentDate.AddDate(0, 0, -i*7)
				_, w := weekDate.ISOWeek()

				if count := weeklyCount[w]; count > 0 {
					avg := weeklyAverages[w] / float64(count)
					trend := "→"
					if i > 0 {
						if avg > lastAvg {
							trend = "↑"
						} else if avg < lastAvg {
							trend = "↓"
						}
					}
					fmt.Printf("Week %02d: %.1f %s (%d days)\n", w, avg, trend, count)
					lastAvg = avg
				} else {
					fmt.Printf("Week %02d: No data\n", w)
				}
			}

			return nil
		},
	}

	return cmd
}
