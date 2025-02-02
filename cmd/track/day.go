/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package track

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var score int
var url string

// dayCmd represents the day command
var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		number, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("argument must be between 1 and 5")
		}
		// Run the custom validation logic
		if number >= 1 && number <= 5 {
			score = number
			return nil
		}
		return fmt.Errorf("invalid score (must be between 1 and 5): %s", args[0])
	}, Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("day called with score: %d\n", score)
	},
}

func init() {
	// dayCmd.Flags().StringArrayVarP(&url, "url", "u", "", "The url")
	// dayCmd.ValidArgs(*[]"score")
	// dayCmd.Args(int(&score))
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
