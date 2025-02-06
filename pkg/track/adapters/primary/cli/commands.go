package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github/b9n2038/act/pkg/track/ports/primary/cli"
	"time"
)

func NewRootCmd(service primary.RatingService) *cobra.Command {
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

func newAddCmd(service primary.RatingService) *cobra.Command {
	return &cobra.Command{
		Use:   "add [rating]",
		Short: "Add a rating for today",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.Add()
		},
	}
}
