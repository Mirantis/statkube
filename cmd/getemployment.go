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
	CompanyName string   `json:"company_name"`
	Domains     []string `json:"domains"`
}

type employment struct {
	CompanyName string `json:"company_name"`
	EndDate     string `json:"end_date"`
}

type user struct {
	Employments []employment `json:"companies"`
	UserName    string       `json:"user_name"`
	Emails      []string     `json:"strings"`
	GithubID    string       `json:"github_id"`
}

func handle(u user, db *gorm.DB) {
	MIN_TIME := time.Unix(0, 0)
	MAX_TIME := time.Unix(1<<40-1, 0)
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
	startDate := MIN_TIME
	for _, employment := range u.Employments {
		var companyDB models.Company
		var workPeriod models.WorkPeriod
		var endDate time.Time
		var err error
		db.FirstOrCreate(&companyDB, models.Company{Name: employment.CompanyName})
		if employment.EndDate != "" {
			endDate, err = time.Parse("2006-Jan-02", employment.EndDate)
			if err != nil {
				panic(fmt.Sprintf(
					"Error parsing date %s\n%s",
					employment.EndDate,
					err.Error(),
				))
			}
		} else {
			endDate = MAX_TIME
		}
		db.FirstOrCreate(&workPeriod, models.WorkPeriod{
			DeveloperID: developer.ID,
			CompanyID:   companyDB.ID,
			Started:     startDate,
		})
		workPeriod.Finished = endDate
		db.Save(&workPeriod)
		startDate = endDate
	}
}

func skipToken(decoder *json.Decoder, n int) {
	for i := 0; i < n; i++ {
		_, err := decoder.Token()
		if err != nil {
			panic("Error parsing file:\n" + err.Error())
		}
	}
}

func initReader(r io.Reader) *json.Decoder {
	decoder := json.NewDecoder(r)
	// Skip '{', '"users"', '['
	skipToken(decoder, 3)
	return decoder
}

func loadCompanies(decoder *json.Decoder, db *gorm.DB) {
	for decoder.More() {
		var company company
		var companyDB models.Company
		decoder.Decode(&company)
		fmt.Printf("Company: %s\n", company)
		db.FirstOrCreate(&companyDB, models.Company{Name: company.CompanyName})
		for _, domain := range company.Domains {
			var domainDB models.Domain
			db.FirstOrCreate(&domainDB, models.Domain{
				Domain:    domain,
				CompanyID: companyDB.ID,
			})
		}
	}
}

func loadEmployment(decoder *json.Decoder, db *gorm.DB) {
	for decoder.More() {
		var u user
		err := decoder.Decode(&u)
		if err != nil {
			panic("Error parsing list:\n" + err.Error())
		}
		handle(u, db)
	}
}

func main() {
	db := db.GetDB()
	filename, exists := os.LookupEnv("EMPLOYMENT_FILE")
	if !exists {
		filename = "default_data.json"
	}
	f, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintln("Error opening file %s: %v ", filename, err.Error()))
	}
	decoder := initReader(f)
	loadEmployment(decoder, db)
	// skip ']', '"companies", '['
	skipToken(decoder, 3)
	loadCompanies(decoder, db)
}
