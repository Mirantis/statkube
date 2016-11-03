package importer

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/Mirantis/statkube/models"
)

const TEST_REPOS = `[
	"kubernetes"
]`

const TEST_DATA = `{
	"users": [{
		"github_id": "user1",
		"emails": ["user1@example.com", "user1@example.pl"],
		"companies": [{
				"company_name": "Intel",
				"end_date": "2015-May-01"
			}, {
				"company_name": "Mirantis",
				"end_date": null
			}

		]
	}, {
		"launchpad_id": "user2",
		"emails": ["user2@example.com", "user2@example.pl"],
		"companies": [{
				"company_name": "Intel",
				"end_date": "2015-May-01"
			}, {
				"company_name": "Mirantis",
				"end_date": null
			}

		]
	}],
	"companies": [
		{
			"company_name": "Mirantis",
			"domains": ["mirantis.com"]
		}
	]
} 
`

const TEST_DATA_COMPLEMENTED = `{
	"users": [{
		"github_id": "user1",
		"emails": ["user1@example.com", "user1@example.pl"],
		"companies": [{
				"company_name": "Intel",
				"end_date": "2015-May-01"
			}, {
				"company_name": "Mirantis",
				"end_date": null
			}

		]
	}, {
		"launchpad_id": "user2",
		"github_id": "user2-gh",
		"emails": ["user2@example.com", "user2@example.pl"],
		"companies": [{
				"company_name": "Intel",
				"end_date": "2015-May-01"
			}, {
				"company_name": "Mirantis",
				"end_date": null
			}

		]
	}],
	"companies": [
	]
} 
`

func getDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(fmt.Sprintf("Failed to create in-memory db: %v", err.Error()))
	}
	models.Migrate(db)
	return db
}

func getTestData(s string) *strings.Reader {
	return strings.NewReader(s)
}

func loadData() *gorm.DB {
	db := getDB()
	data := getTestData(TEST_DATA)
	LoadAll(data, db)
	data = getTestData(TEST_REPOS)
	LoadRepos(data, db)
	return db
}

// getK8sRepo returns the kubernetes repo object from db
func getK8sRepo(db *gorm.DB) *models.Repository {
	var repo models.Repository
	db.First(&repo, "repo = ?", "kubernetes")
	return &repo
}
