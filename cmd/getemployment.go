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
	Emails      []string     `json:"emails"`
	GithubID    string       `json:"github_id"`
	LaunchpadID string       `json:"launchpad_id"`
}

func loadRepos(f io.Reader, db *gorm.DB) {
	var repo string
	decoder := json.NewDecoder(f)
	skipToken(decoder, 1) // "["
	for decoder.More() {
		decoder.Decode(&repo)
		var repoDB models.Repository
		db.FirstOrCreate(&repoDB, models.Repository{User: "kubernetes", Repo: repo})
	}

}

// findDeveloper finds the user in the database or initialize it. If necessary, fill the 'github_id'
func findDeveloper(u *user, db *gorm.DB) *models.Developer {
	var developer models.Developer
	// Try to find the user
	if u.GithubID != "" {
		search := db.Where("developers.github_id = ?", u.GithubID).First(&developer)
		if !search.RecordNotFound() {
			return &developer
		}
	}
	if u.LaunchpadID != "" {
		search := db.Where("developers.launchpad_id = ?", u.LaunchpadID).First(&developer)
		if !search.RecordNotFound() {
			if developer.GithubID != u.GithubID {
				developer.GithubID = u.GithubID
				db.Save(&developer)
			}
			return &developer
		}
	}
	if u.LaunchpadID == "" && u.GithubID == "" {
		return nil
	}
	developer.LaunchpadID = u.LaunchpadID
	developer.GithubID = u.GithubID
	db.Create(&developer)
	return &developer

}

// handle handles single entry in the json
func handle(u user, db *gorm.DB) {
	MIN_TIME := time.Unix(0, 0)
	MAX_TIME := time.Unix(1<<40-1, 0)
	developer := findDeveloper(&u, db)
	if developer == nil {
		return
	}
	developer.FullName = u.UserName
	for _, email := range u.Emails {
		var emailDB models.Email
		db.FirstOrInit(&emailDB, models.Email{Email: email})
		emailDB.DeveloperID = developer.ID
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

// skipToken skips given number of tokens in parsed json
func skipToken(decoder *json.Decoder, n int) {
	for i := 0; i < n; i++ {
		_, err := decoder.Token()
		if err != nil {
			panic("Error parsing file:\n" + err.Error())
		}
	}
}

// initReader prepares the reader to iterate users
func initReader(r io.Reader) *json.Decoder {
	decoder := json.NewDecoder(r)
	// Skip '{', '"users"', '['
	skipToken(decoder, 3)
	return decoder
}

// loadCompanies loads data about companies and their domains
func loadCompanies(decoder *json.Decoder, db *gorm.DB) {
	for decoder.More() {
		var company company
		var companyDB models.Company
		decoder.Decode(&company)
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

// loadEmployment loads data about developers, their emails and work periods
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
		filename = "repos.json"
	}
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(fmt.Sprintln("Error opening repos file: %v ", err.Error()))
	}
	loadRepos(f, db)
	filename, exists = os.LookupEnv("EMPLOYMENT_FILE")
	if !exists {
		filename = "default_data.json"
	}
	f, err = os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(fmt.Sprintln("Error opening file %s: %v ", filename, err.Error()))
	}
	decoder := initReader(f)
	loadEmployment(decoder, db)
	// skip ']', '"companies", '['
	skipToken(decoder, 3)
	loadCompanies(decoder, db)
}
