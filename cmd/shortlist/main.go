// cmd/shortlist/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"act/pkg/short/domain/ports/secondary"
	"act/pkg/short/domain/service"
	"act/pkg/short/infrastructure/cli"
	"act/pkg/short/infrastructure/storage"

	"github.com/spf13/cobra"
)

func processPath(path string) string {
	if path == "" {
		return ""
	}

	// Expand ~ to home directory
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}

	// Convert to absolute path
	if !filepath.IsAbs(path) {
		currentDir, err := os.Getwd()
		if err == nil {
			path = filepath.Join(currentDir, path)
		}
	}

	// Clean the path to remove any .. or . elements
	return filepath.Clean(path)
}

func getStoragePath() (string, error) {
	// Create a temporary root command just for parsing the path flag
	tmpRoot := &cobra.Command{
		Use:           "short",
		SilenceErrors: true,
		SilenceUsage:  true,
		Run:           func(cmd *cobra.Command, args []string) {},
	}

	tmpRoot.PersistentFlags().String("path", "", "Override default storage location")

	// Parse flags without executing the command
	args := os.Args[1:]
	if err := tmpRoot.ParseFlags(args); err != nil {
		return "", fmt.Errorf("failed to parse flags: %w", err)
	}

	path, err := tmpRoot.Flags().GetString("path")
	if err != nil {
		return "", fmt.Errorf("failed to get path flag: %w", err)
	}

	return processPath(path), nil
}

func main() {

	// Get storage path from flags
	storagePath, err := getStoragePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing path flag: %v\n", err)
		os.Exit(1)
	}

	// Set up storage
	var store secondary.ListStorage
	if storagePath != "" {
		store = storage.NewMarkdownStorage(storagePath)
	} else {
		defaultPath := filepath.Join(os.Getenv("HOME"), ".shortlist")
		store = storage.NewMarkdownStorage(defaultPath)
	}

	listService := service.NewListService(store)

	// Create and run app
	app := cli.NewApp(listService)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

}
