package cli

import (
	ratingService "track/internal/track/application/rating" // aliased this import
	"track/internal/track/domain/rating"

	// "track/internal/track/domain/short"
	"strconv"

	//"track/internal/track/ports/primary/rating"
	"context"
	"fmt"
	"time"

	"github.com/snabb/isoweek"
	"github.com/spf13/cobra"
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
// this is broken
func parseDayID(id string) (time.Time, error) {
	// Parse format: YYwWW-D
	if len(id) != 7 || id[2] != 'w' || id[5] != '-' {
		return time.Time{}, fmt.Errorf("invalid format, expected YYwWW-D")
	}

	year, _ := strconv.Atoi("20" + id[0:2])
	week, _ := strconv.Atoi(id[3:5])
	day, _ := strconv.Atoi(id[6:])

	day = day % 7
	if day < 1 || day > 7 {
		return time.Time{}, fmt.Errorf("invalid format, isoWeekDay range from 0-6, with Mon as start of the week")
	}

	// Get the date for the start of the week
	// date1 := isoWeekStart(year, week)
	// fmt.Println("date for start of week", date1.Local())

	//get start of date
	startYr, startMth, startDay := isoweek.StartDate(year, week)

	date := time.Date(startYr, startMth, startDay, 0, 0, 0, 0, time.UTC)
	// Add days
	date = date.AddDate(0, 0, day-1)

	return date, nil
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
		dayID    string
		weekday  string
		target   time.Time
		fillGaps bool
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
				// fmt.Printf("handle weekday %s set\n", weekday)
				iWeekday, iWeekdayErr := strconv.Atoi(weekday)
				if iWeekdayErr != nil {
					fmt.Errorf("invalid weekday format: %w", iWeekdayErr)
					return iWeekdayErr
				}
				weekday := iWeekday
				_, week := target.ISOWeek()

				start_yr, start_mth, start_day := isoweek.StartDate(target.Year(), week)
				// fmt.Printf("start of isoweek %d, %d, %s, %d\n", week, start_yr, start_mth, start_day)
				target = time.Date(start_yr, start_mth, start_day, 0, 0, 0, 0, time.UTC)
				target = target.AddDate(0, 0, weekday-1)
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
	//may combine these
	cmd.Flags().StringVarP(&dayID, "long", "l", "", "Day ID in format YYwWW-D. 25w05-3")
	cmd.Flags().StringVarP(&weekday, "weekday", "d", "", "Week Day 1-7 (e.g. 1 = Monday")
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
