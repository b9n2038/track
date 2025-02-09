// cmd/shortlist/main.go
package main

import (
	"log"
	"os"

	"act/pkg/short/domain/service"
	"act/pkg/short/infrastructure/cli"
	"act/pkg/short/infrastructure/storage"
)

func main() {
	baseDir := os.Getenv("HOME")
	if baseDir == "" {
		baseDir = "."
	}

	store := storage.NewMarkdownStorage(baseDir)
	listService := service.NewListService(store)
	app := cli.NewApp(listService)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
