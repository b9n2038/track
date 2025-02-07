package cli

import (
	ratingService "act/pkg/track/application/rating" // aliased this import
	"act/pkg/track/domain/rating"
	"strconv"

	//"act/pkg/track/ports/primary/rating"
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func NewRootCmd(service ratingService.Service) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "rating",
		Short: "Day rating management tool",
	}

	rootCmd.AddCommand(
		newAddCmd(service),
		// newListCmd(service),
		// newWeekCmd(service),
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

func newAddCmd(service ratingService.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "add [rating]",
		Short: "Add a rating for today",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(cmd.OutOrStderr, "call add day rating")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			value, err := rating.NewRating(parseInt(args[0])) // using domain package
			if err != nil {
				return err
			}

			dayRating, err := service.AddDayRating(ctx, time.Now(), value)
			if err != nil {
				return err
			}

			fmt.Printf("Added rating: %s\n", dayRating)
			return nil
		},
	}
}

// func newListCmd(service rating.Service) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "list",
// 		Short: "List""
// 		// Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Print(cmd.OutOrStderr, "list")
// 			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 			defer cancel()
// 			service.List(ctx, time.Now(), args[0])
// 			return nil
// 		},
// 	}
// }
