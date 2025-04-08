package cli

import (
	ratingService "track/internal/track/application/rating" // aliased this import
	"track/internal/track/domain/rating"

	// "track/internal/track/domain/short"
	"strconv"

	//"track/internal/track/ports/primary/rating"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

func NewRootCmd(ratingService *ratingService.Service) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "track",
		Short: "Track the important stuff",
	}

	rootCmd.AddCommand(
		newDayCmd(ratingService),
	)
	return rootCmd
}

func newDayCmd(service *ratingService.Service) *cobra.Command {
	dayCmd := &cobra.Command{
		Use:   "day",
		Short: "Day rating",
	}

	dayCmd.AddCommand(
		newSetCmd(service),
		newListCmd(service),
		newWeekCmd(service),
	)

	return dayCmd
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

	day = day % 7
	if day < 0 || day > 7 {
		return time.Time{}, fmt.Errorf("invalid format, isoWeekDay range from 1-7, with Mon as start of the week")
	}

	//d,w,y => d, m, y
	//startDate, convert to ordinal days, add weeks, plus day

	// get the absDay of start of year and add week days + weekday
	start := time.Date(year, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	days := ((week - 1) * 7) + day
	date := start.Add(time.Duration(time.Second.Nanoseconds() * int64(86400*days))) //int64(days * time.Second)))

	return date, nil
}

// GetWeekdayInISOWeek returns the date for the specified weekday (1-7)
// in the ISO week of the provided reference date
func GetWeekdayInISOWeek(referenceDate time.Time, weekday int) (time.Time, error) {
	if weekday < 1 || weekday > 7 {
		return time.Time{}, fmt.Errorf("invalid weekday: %d, must be between 1 and 7", weekday)
	}

	// Get the ISO year and week of the reference date
	year, week := referenceDate.ISOWeek()

	// Calculate which day of the week the referenceDate is (in ISO terms, 1=Monday, 7=Sunday)
	isoWeekday := int(referenceDate.Weekday())
	if isoWeekday == 0 { // Sunday in Go is 0
		isoWeekday = 7
	}

	// Calculate the offset to the target weekday
	daysToAdd := weekday - isoWeekday

	// Get the requested weekday in the same ISO week
	targetDate := referenceDate.AddDate(0, 0, daysToAdd)

	// Verify we're still in the same ISO week
	targetYear, targetWeek := targetDate.ISOWeek()
	if targetYear != year || targetWeek != week {
		// This can only happen if we crossed a week boundary
		// Adjust by a week
		if daysToAdd < 0 {
			// If we went back to the previous week, go forward a week
			targetDate = targetDate.AddDate(0, 0, 7)
		} else {
			// If we went to the next week, go back a week
			targetDate = targetDate.AddDate(0, 0, -7)
		}
	}

	return targetDate, nil
}

// Helper function to convert ISO weekday (1-7) to Go's time.Weekday
func convertToWeekday(isoWeekday int) time.Weekday {
	// In ISO, 1=Monday, 7=Sunday
	// In Go, 0=Sunday, 1=Monday, ..., 6=Saturday
	if isoWeekday == 7 {
		return time.Sunday
	}
	return time.Weekday(isoWeekday)
}

//	func isoWeekStart(year, week int) time.Time {
//		// Find a day in the week
//		jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.Local)
//		_, w := jan4.ISOWeek()
//		// Move to the desired week
//		daysSinceJan4 := (week - w) * 7
//		return jan4.AddDate(0, 0, daysSinceJan4-int(jan4.Weekday())+1)
//	}
func newSetCmd(service *ratingService.Service) *cobra.Command {
	var (
		dayID   string
		weekday string
		target  time.Time
	)

	cmd := &cobra.Command{
		Use:   "set [rating]",
		Short: "Set a day rating between 1 and 5, for today.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			value, err := rating.NewRating(parseInt(args[0])) // using domain package
			if err != nil {
				return err
			}

			target = time.Now()

			if dayID != "" {
				// Add rating for specified day
				parsedDate, err := parseDayID(dayID)
				if err != nil {
					fmt.Errorf("invalid day ID format: %w", err)
					return err
				}
				target = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
			}

			if weekday != "" {
				day, iWeekdayErr := strconv.Atoi(weekday)
				if iWeekdayErr != nil {
					fmt.Errorf("invalid weekday format: %w", iWeekdayErr)
					return iWeekdayErr
				}
				target, _ = GetWeekdayInISOWeek(target, day)
				// fmt.Printf("target %d, %d, %d\n", target.Year(), target.Month(), target.Day())
			}
			_, err = service.SetDayRating(ctx, target, value)
			if err != nil {
				return err
			}

			//fmt.Printf("rating: %s\n", dayRating)
			// } else {
			//
			//
			// 	// fmt.Println("parseDayID", date.Local())
			// 	_, err = service.SetDayRating(ctx, date, value)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	//fmt.Printf("Added rating for %s: %s\n", dayID, dayRating)
			//
			// 	// Fill gaps if requested
			// 	if fillGaps {
			// 		lastRating, err := service.GetLastRatingBefore(ctx, date)
			// 		if err != nil {
			// 			return fmt.Errorf("getting last rating: %w", err)
			// 		}
			//
			// 		if !lastRating.Date.IsZero() {
			// 			filled, err := service.FillMissingRatings(ctx, lastRating.Date, date, value)
			// 			if err != nil {
			// 				return fmt.Errorf("filling gaps: %w", err)
			// 			}
			// 			if len(filled) > 0 {
			// 				fmt.Printf("Filled %d missing days with rating %s\n", len(filled), value.String())
			// 			}
			// 		}
			// }
			// },
			//
			return nil
		},
	}

	cmd.Flags().StringVarP(&dayID, "long", "l", "", "Day ID in format YYwWW-D. 25w05-3")
	cmd.Flags().StringVarP(&weekday, "weekday", "d", "", "Week Day 1-7 (e.g. 1 = Monday")
	// cmd.Flags().BoolVarP(&fillGaps, "fill", "f", false, "Fill missing days from last entry")
	return cmd
}

func newListCmd(service *ratingService.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List ratings for a week, default current.",
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
	// cmd.Flags().StringVarP(&dayID, "long", "l", "", "Day ID in format YYwWW-D. 25w05-3")
	return cmd
}

func newWeekCmd(service *ratingService.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
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

			// Print daily list
			fmt.Printf("\nDaily List:\n")
			fmt.Printf("───────────────\n")
			for _, r := range ratings {
				fmt.Printf("%s: %s %s\n",
					r.Date.Format("Mon"),
					r.Rating.String(),
					r.Rating.Emoji(),
				)
			}

			//todo; Print daily grid

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
