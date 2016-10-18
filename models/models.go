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
