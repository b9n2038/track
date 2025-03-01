/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"track/internal/track/adapters/primary/cli"
	"track/internal/track/adapters/secondary/file"
	"track/internal/track/application/rating"
)

func Execute() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	//todo: make configurable
	repoPath := filepath.Join(homeDir, ".track.rating.json")
	repo, err := file.NewFileRepository(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	ratingService := rating.NewService(repo)

	rootCmd := cli.NewRootCmd(ratingService)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
