package main

import (
	"fmt"
	"os"

	"github.com/Mirantis/statkube/db"
	"github.com/Mirantis/statkube/importer"
)

func main() {
	db := db.GetDB()
	filename, exists := os.LookupEnv("EMPLOYMENT_FILE")
	if !exists {
		filename = "repos.json"
	}
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(fmt.Sprintln("Error opening repos file: %v ", err.Error()))
	}
	importer.LoadRepos(f, db)
	filename, exists = os.LookupEnv("EMPLOYMENT_FILE")
	if !exists {
		filename = "default_data.json"
	}
	f, err = os.Open(filename)
	defer f.Close()
	importer.LoadAll(f, db)
	if err != nil {
		panic(fmt.Sprintln("Error opening file %s: %v ", filename, err.Error()))
	}
}
