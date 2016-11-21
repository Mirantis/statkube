package main

import (
	"os"
	"time"

	"github.com/Mirantis/statkube/db"
	"github.com/Mirantis/statkube/importer"
	"github.com/Mirantis/statkube/models"
)

func main() {
	var repositories []models.Repository
	limitStr := os.Args[1]
	limit, err := time.Parse("2006-01-02", limitStr)
	if err != nil {
		panic(err.Error())
	}
	db := db.GetDB()
	token, exists := os.LookupEnv("GITHUB_TOKEN")
	if !exists {
		panic("Set GITHUB_TOKEN")
	}
	client := importer.NewClient(token)

	db.Find(&repositories)

	for _, repository := range repositories {
		prs := client.ListPRs(repository.User, repository.Repo, limit)
		for prs.More() {
			importer.HandlePR(prs.Scan(), client, &repository, db)
		}
	}
}
