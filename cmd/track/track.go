/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package track

import (
	"github.com/spf13/cobra"
)

// trackCmd represents the track command
var TrackCmd = &cobra.Command{
	Use:   "track",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// rootCmd.AddCommand(trackCmd)
	TrackCmd.AddCommand(dayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
