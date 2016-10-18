package models

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Developer struct {
	gorm.Model
	FullName    string `gorm:"not null"`
	Accounts    []Account
	Emails      []Email
	WorkPeriods []WorkPeriod
}

type Account struct {
	gorm.Model
	Developer   Developer
	DeveloperID int
	Username    string
}

type Email struct {
	gorm.Model
	Developer   Developer
	DeveloperID int
	Email       string `gorm:"not null, unique"`
}

type Company struct {
	gorm.Model
	Name        string `gorm:"not null"`
	WorkPeriods []WorkPeriod
}

type WorkPeriod struct {
	gorm.Model
	Company     Company
	CompanyID   int
	Developer   Developer
	DeveloperID int
	Position    string
	Started     time.Time `gorm:"null"`
	Finished    time.Time `gorm:"not null"`
}

type PullRequest struct {
	gorm.Model
	WorkPeriod   WorkPeriod
	WorkPeriodId sql.NullInt64
	Developer    Developer
	DeveloperId  int
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&Developer{},
		&Account{},
		&Email{},
		&Company{},
		&WorkPeriod{},
		&PullRequest{},
	)
}

type DevStats struct {
	FullName string
	PRCount  int
}

func GetDevStats(db *gorm.DB) ([]DevStats, error) {
	var developers []DevStats
	rows, err := db.Table("developers").Select("developers.full_name, COUNT(prs.developer_id)").Joins("left join pull_requests prs on prs.developer_id = developers.id").Group("prs.developer_id").Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		var count int
		err := rows.Scan(&name, &count)
		if err != nil {
			return nil, err
		}
		developers = append(developers, DevStats{name, count})
	}

	return developers, nil
}
