package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Mirantis/statkube/db"
	"github.com/Mirantis/statkube/models"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
)

// getDeveloper gets or creates a Developer object basing on github_id from pr
func getDeveloper(pr *github.PullRequest, db *gorm.DB) *models.Developer {
	var developer models.Developer
	db.FirstOrCreate(&developer, models.Developer{GithubID: *pr.User.Login})
	return &developer

}

// assumeIndependent returns a pair of company/developer with *independent
// company and username from pr
func assumeIndependent(pr *github.PullRequest, db *gorm.DB) (*models.Company, *models.Developer) {
	var company models.Company
	db.FirstOrCreate(&company, models.Company{Name: "*independent"})
	return &company, getDeveloper(pr, db)
}

// deduceFromWorkPeriod finds the pair of company/developer by github_id
// and the date of  pull request
func deduceFromWorkPeriod(pr *github.PullRequest, db *gorm.DB) (*models.Company, *models.Developer) {
	var workPeriod models.WorkPeriod
	var company models.Company
	var developer models.Developer
	search := db.Joins("JOIN developers ON developers.id = work_periods.developer_id").
		Where("developers.github_id = ?", pr.User.Login).
		Where("? BETWEEN work_periods.started AND work_periods.finished", pr.CreatedAt).
		First(&workPeriod)
	if search.RecordNotFound() {
		return nil, nil
	}
	db.Model(&workPeriod).Related(&company)
	db.Model(&workPeriod).Related(&developer)
	return &company, &developer
}

// getEmails returns the set of all emails in prs
func getEmails(pr *github.PullRequest, client *github.Client) map[string]struct{} {
	commits, _, _ := client.PullRequests.ListCommits(
		"kubernetes", "kubernetes", *pr.Number, nil,
	)
	if len(commits) == 0 {
		fmt.Printf(
			"PR empty %s\n", *pr.URL,
		)
		return nil
	}
	emailSet := make(map[string]struct{})
	for _, commit := range commits {
		emailSet[(*commit.Commit.Author.Email)] = struct{}{}
	}
	return emailSet
}

// deduceFromEmail will find the user by his e-mail to see if she is in the DB
// This approach will succeed if the user is not registered with his github_id
func deduceFromEmail(pr *github.PullRequest, emails map[string]struct{}, db *gorm.DB) (*models.Company, *models.Developer) {
	var email string
	var workPeriod models.WorkPeriod
	var company models.Company
	var developer models.Developer
	var emailsObjs []models.Email
	// Hack for getting one value from dict
	for curEmail, _ := range emails {
		email = curEmail
		break
	}
	search := db.Joins("RIGHT JOIN developers ON developers.id = work_periods.developer_id").
		Joins("LEFT JOIN emails ON emails.developer_id = developers.id").
		Where("emails.email = ?", email).
		Where("? BETWEEN work_periods.started AND work_periods.finished", pr.CreatedAt).
		First(&workPeriod)
	if search.RecordNotFound() {
		return nil, nil
	}
	db.Model(&workPeriod).Related(&company)
	db.Model(&workPeriod).Related(&developer)
	if len(emails) > 1 { // Check consistency
		db.Model(&developer).Related(&emailsObjs)
		for _, e := range emailsObjs {
			_, exists := emails[e.Email]
			if !exists {
				fmt.Printf("Inconsistent emails %v\n", emails)
			}
		}
	}
	return &company, &developer
}

// deduceFromDomain will try to find the company by domain of emails
// it will be used if we can't find the user affiliation in our db
func deduceFromDomain(pr *github.PullRequest, emails map[string]struct{}, db *gorm.DB) (*models.Company, *models.Developer) {
	var company models.Company
	domain := ""
	for email, _ := range emails {
		bits := strings.Split(email, "@")
		if len(bits) != 2 {
			fmt.Printf("Invalid email %s\n", email)
			return nil, nil
		}
		curDomain := bits[1]
		if domain != "" && curDomain != domain {
			fmt.Printf("Inconsistent domains %v\n", emails)
			return nil, nil
		}
	}
	for {
		search := db.Joins("RIGHT JOIN domains ON domains.company_id = companies.id").
			Where("domains.domain = ?", domain).
			First(&company)
		if search.RecordNotFound() {
			bits := strings.SplitN(domain, ".", 2)
			if len(bits) == 1 {
				return nil, nil
			}
			domain = bits[1]
		} else {
			break
		}

	}
	developer := getDeveloper(pr, db)
	return &company, developer
}

// deduceCompanyAndDev tries get the developer and company of a pr. It tries
// to give the best result using the strategies in this order
// * Find the developer and his affiliation by github_id
// * Find the developer and his affiliation by e-mails in commits
// * Find the affiliation by e-mails in commits and create/use a dev by github_id
// * Declare developer "independent" and use his github_id

func deduceCompanyAndDev(pr *github.PullRequest, client *github.Client, db *gorm.DB) (*models.Company, *models.Developer) {
	company, developer := deduceFromWorkPeriod(pr, db)
	if company != nil {
		return company, developer
	}
	emails := getEmails(pr, client)
	company, developer = deduceFromEmail(pr, emails, db)
	if company != nil {
		return company, developer
	}
	company, developer = deduceFromDomain(pr, emails, db)
	if company != nil {
		return company, developer
	}
	return assumeIndependent(pr, db)

}

// handlePR handles single pull request
func handlePR(pr *github.PullRequest, client *github.Client, repository *models.Repository, db *gorm.DB) {
	var prDB models.PullRequest
	if pr.MergedAt == nil {
		return
	}

	db.FirstOrInit(&prDB, models.PullRequest{Url: *pr.URL, RepositoryID: repository.ID})

	company, developer := deduceCompanyAndDev(pr, client, db)

	prDB.Company = *company
	prDB.Developer = *developer
	prDB.Created = *pr.CreatedAt
	prDB.Merged = pr.MergedAt
	db.Save(&prDB)
}

func main() {
	var repositories []models.Repository
	limitStr := os.Args[1]
	limit, err := time.Parse("2006-01-02", limitStr)
	if err != nil {
		panic(err.Error())
	}
	db := db.GetDB()
	token, exists := os.LookupEnv("GITHUB_TOKEN")
	if !exists {
		panic("Set GITHUB_TOKEN")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	db.Find(&repositories)
	for _, repository := range repositories {
		opt := &github.PullRequestListOptions{
			ListOptions: github.ListOptions{PerPage: 1000},
			State:       "closed",
			Sort:        "updated",
			Direction:   "desc",
		}
		for {
			limitMet := false
			prs, resp, err := client.PullRequests.List(
				repository.User, repository.Repo, opt,
			)
			fmt.Printf("repo: %v found: %v prs\n", repository.Repo, len(prs))
			if err != nil {
				panic(err.Error())
			}
			for _, pr := range prs {
				//if pr is updated before limit, break as prs are sorted by updatedAt
				if pr.UpdatedAt.Before(limit) {
					limitMet = true
					break
				}
				handlePR(pr, client, &repository, db)
			}
			if resp.NextPage == 0 || limitMet {
				break
			}
			opt.ListOptions.Page = resp.NextPage
		}
	}
}
