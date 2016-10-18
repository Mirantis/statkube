package main

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/Mirantis/statkube/models"
)

func main() {
	dbstring, exists := os.LookupEnv("STATKUBE_DB")
	if !exists {
		panic("Set the value of STATKUBE_DB")
	}
	db, err := gorm.Open("postgres", dbstring)
	if err != nil {
		panic("Error connecting to db.\n" + err.Error())
	}
	db.AutoMigrate(
		&models.Developer{},
		&models.Account{},
		&models.Email{},
		&models.Company{},
		&models.WorkPeriod{},
		&models.PullRequest{},
	)
	defer db.Close()
}
