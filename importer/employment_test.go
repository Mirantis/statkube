package importer

import (
	"strings"
	"testing"

	"github.com/Mirantis/statkube/models"
)

// TestLoadNewDeveloper tests whether we correctly load a user with github_id
func TestLoadNewDeveloper(t *testing.T) {
	var developer models.Developer
	var emails []models.Email
	db := loadData()
	db.Where("github_id = ?", "user1").Find(&developer)
	if developer.GithubID != "user1" {
		t.Error("Didn't find user1")
	}
	db.Model(&developer).Related(&emails)
	if len(emails) != 2 {
		t.Error("Number of emails wrong")
	}
}

// TestLoadNewDeveloperByLaunchpad tests whether we correctly load a user without github_id
func TestLoadNewDeveloperByLaunchpad(t *testing.T) {
	var developer models.Developer
	var emails []models.Email
	db := loadData()
	db.Where("launchpad_id = ?", "user2").Find(&developer)
	if developer.LaunchpadID != "user2" {
		t.Error("Didn't find user2")
	}
	db.Model(&developer).Related(&emails)
	if len(emails) != 2 {
		t.Error("Number of emails wrong")
	}
}

func TestComplementGithubID(t *testing.T) {
	var developer models.Developer
	var count int
	db := loadData()
	LoadAll(strings.NewReader(TEST_DATA_COMPLEMENTED), db)
	db.Table("developers").Where("launchpad_id = ?", "user2").Count(&count)
	if count != 1 {
		t.Error("User duplicated instead of complimenting her data")
	}
	db.Where("launchpad_id = ?", "user2").Find(&developer)
	if developer.GithubID != "user2-gh" {
		t.Error("Github id not supplemented")
	}
}
