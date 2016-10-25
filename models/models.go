package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Developer struct {
	gorm.Model
	FullName    string `gorm:"not null"`
	Emails      []Email
	WorkPeriods []WorkPeriod
	GithubID    string `gorm:"not null"`
}

type Email struct {
	gorm.Model
	Developer   Developer
	DeveloperID uint
	Email       string `gorm:"not null, unique"`
}

type Company struct {
	gorm.Model
	Name        string `gorm:"not null"`
	WorkPeriods []WorkPeriod
}

type Domain struct {
	gorm.Model
	Domain    string `gorm:"not null"`
	Company   Company
	CompanyID uint `gorm:"not null"`
}

type WorkPeriod struct {
	gorm.Model
	Company     Company
	CompanyID   uint
	Developer   Developer
	DeveloperID uint
	Position    string
	Started     time.Time
	Finished    time.Time
}

type PullRequest struct {
	gorm.Model
	Developer   Developer
	DeveloperID uint
	Company     Company
	CompanyID   uint `gorm:"not null"`
	Url         string
	Created     time.Time
	Merged      *time.Time
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&Developer{},
		&Email{},
		&Company{},
		&Domain{},
		&WorkPeriod{},
		&PullRequest{},
	)
}

type DevStats struct {
	FullName string
	PRCount  uint
}

func GetDevStats(db *gorm.DB) ([]DevStats, error) {
	//TODO: add time logic from GetCompanyStats
	var developers []DevStats
	rows, err := db.Table("developers").
		Select("developers.full_name, COUNT(prs.developer_id)").
		Joins("left join pull_requests prs on prs.developer_id = developers.id").
		Group("prs.developer_id").
		Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		var count uint
		err := rows.Scan(&name, &count)
		if err != nil {
			return nil, err
		}
		developers = append(developers, DevStats{name, count})
	}

	return developers, nil
}

func GetCompanyStats(db *gorm.DB, start, end time.Time) ([]DevStats, error) {
	var stats []DevStats

	rows, err := db.Table("companies").
		Select("companies.name, COUNT(pull_requests.id)").
		Joins("left join pull_requests on pull_requests.company_id=companies.id").
		Where("pull_requests.merged BETWEEN ? AND ?", start, end). //filter by start and end date when it's available
		Group("companies.id").
		Rows()

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		var count uint
		err := rows.Scan(&name, &count)
		if err != nil {
			return nil, err
		}
		stats = append(stats, DevStats{name, count})
	}

	return stats, nil
}
