/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/b9n2038/act/pkg/track/adapters/primary/cli"
	"github.com/b9n2038/act/pkg/track/adapters/secondary/file"
	"github.com/b9n2038/act/pkg/track/application/rating"
	"log"
	"os"
	"path/filepath"
)

// func Execute() {
// 	err := rootCmd.Execute()
// 	if err != nil {
// 		os.Exit(1)
// 	}
// }
//

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	repoPath := filepath.Join(homeDir, ".track.rating.json")
	repo, err := file.NewFileRepository(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	// Setup service
	service := rating.NewService(repo)

	rootCmd := cli.NewRootCmd(service)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
