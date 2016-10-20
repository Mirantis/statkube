package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Mirantis/statkube/db"
	"github.com/Mirantis/statkube/models"
	"github.com/jinzhu/gorm"
)

type company struct {
	CompanyName string `json:"company_name"`
	EndDate     string `json:"end_date"`
}

type user struct {
	Companies []company `json:"companies"`
	UserName  string    `json:"user_name"`
	Emails    []string  `json:"strings"`
	GithubID  string    `json:"github_id"`
}

func handle(u user, db *gorm.DB) {
	if u.GithubID == "" {
		return
	}
	var developer models.Developer
	db.FirstOrCreate(&developer, models.Developer{GithubID: u.GithubID})
	developer.FullName = u.UserName
	for _, email := range u.Emails {
		var emailDB models.Email
		db.FirstOrInit(&email, models.Email{Email: email})
		emailDB.Developer = developer
		db.Save(&emailDB)
	}
	for _, company := range u.Companies {
		var companyDB models.Company
		var workPeriod models.WorkPeriod
        var endDatePt *time.Time
		db.FirstOrCreate(&companyDB, models.Company{Name: company.CompanyName})
        fmt.Printf("EndDate: %s", company.EndDate)
		if company.EndDate != "" {
            endDate, err := time.Parse("2006-Jan-02", company.EndDate)
			if err != nil {
				panic(fmt.Sprintf(
					"Error parsing date %s\n%s",
					company.EndDate,
					err.Error(),
				))
			}
            endDatePt = &endDate
		}
		db.FirstOrCreate(&workPeriod, models.WorkPeriod{
			DeveloperID: developer.ID,
			CompanyID:   companyDB.ID,
			Finished:  endDatePt,
		})
	}
}

func initReader(r io.Reader) *json.Decoder {
	decoder := json.NewDecoder(r)
	// Skip '{', '"users"', '['
	for i := 0; i < 3; i++ {
		_, err := decoder.Token()
		if err != nil {
			panic("Error parsing file:\n" + err.Error())
		}
	}
	return decoder
}

func main() {
	filename, exists := os.LookupEnv("EMPLOYMENT_FILE")
	if !exists {
		filename = "default_data.json"
	}
	f, err := os.Open(filename)
	if err != nil {
		panic("Error opening file:\n" + err.Error())
	}
	db := db.GetDB()
	decoder := initReader(f)
	for decoder.More() {
		var u user
		err = decoder.Decode(&u)
		if err != nil {
			panic("Error parsing list:\n" + err.Error())
		}
		handle(u, db)
	}
}
