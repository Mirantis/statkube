package main

import (
	"github.com/Mirantis/statkube/db"
	"github.com/Mirantis/statkube/models"
)

func main() {
	db := db.GetDB()
	models.Migrate(db)
	defer db.Close()
}
