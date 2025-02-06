package cli

import (
	"act/pkg/track/application/rating"
	"act/pkg/track/ports/primary/cli"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

func NewRootCmd(service rating.Service) *cobra.Command {
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

func newAddCmd(service rating.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "add [rating]",
		Short: "Add a rating for today",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.AddDayRating(ctx, time.Date, args[0])
		},
	}
}
