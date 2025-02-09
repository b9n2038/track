/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"act/pkg/track/adapters/primary/cli"
	"act/pkg/track/adapters/secondary/file"
	"act/pkg/track/application/rating"
	"log"
	"os"
	"path/filepath"
)

func Execute() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

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
