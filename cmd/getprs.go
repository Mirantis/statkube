package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Mirantis/statkube/db"
	"github.com/Mirantis/statkube/models"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
)

func getDeveloper(pr *github.PullRequest, db *gorm.DB) *models.Developer {
	var developer models.Developer
	db.FirstOrCreate(&developer, models.Developer{GithubID: *pr.User.Login})
	return &developer

}

func assumeIndependent(pr *github.PullRequest, db *gorm.DB) (*models.Company, *models.Developer) {
	var company models.Company
	db.FirstOrCreate(&company, models.Company{Name: "*independent"})
	return &company, getDeveloper(pr, db)
}

func deduceFromEmail(pr *github.PullRequest, client *github.Client, db *gorm.DB) (*models.Company, *models.Developer) {
	var company models.Company
	commits, _, _ := client.PullRequests.ListCommits(
		"kubernetes", "kubernetes", *pr.Number, nil,
	)
	if len(commits) == 0 {
		fmt.Printf(
			"PR empty %s\n", *pr.URL,
		)
		return nil, nil
	}
	email := *commits[0].Commit.Author.Email
	for _, commit := range commits {
		if *commit.Commit.Author.Email != email {
			fmt.Printf(
				"Inconsistent emails in PR: %s, %s != %s\n",
				*pr.URL, email, *commit.Commit.Author.Email,
			)
			return nil, nil
		}
	}
	domain := strings.Split(email, "@")[1]
	search := db.Joins("RIGHT JOIN domains ON domains.company_id = companies.id").
		Where("domains.domain = ?", domain).
		First(&company)
	if search.RecordNotFound() {
		fmt.Printf("No company for domain %s\n", domain)
		return nil, nil
	}
	developer := getDeveloper(pr, db)
	return &company, developer
}

func deduceCompanyAndDev(pr *github.PullRequest, client *github.Client, db *gorm.DB) (*models.Company, *models.Developer) {
	var workPeriod models.WorkPeriod
	search := db.Joins("JOIN developers ON developers.id = work_periods.developer_id").
		Where("developers.github_id = ?", pr.User.Login).
		Where("? BETWEEN work_periods.started AND work_periods.finished", pr.CreatedAt).
		First(&workPeriod)
	if !search.RecordNotFound() {
		return &workPeriod.Company, &workPeriod.Developer
	}
	company, developer := deduceFromEmail(pr, client, db)
	if company != nil {
		return company, developer
	}
	return assumeIndependent(pr, db)

}

func handlePR(pr *github.PullRequest, client *github.Client, db *gorm.DB) {
	var prDB models.PullRequest
	if pr.MergedAt == nil {
		return
	}

	db.FirstOrInit(&prDB, models.PullRequest{Url: *pr.URL})

	company, developer := deduceCompanyAndDev(pr, client, db)

	prDB.Company = *company
	prDB.Developer = *developer
	prDB.Created = *pr.CreatedAt
	prDB.Merged = pr.MergedAt
	db.Save(&prDB)
}

func main() {
	db := db.GetDB()
	token, exists := os.LookupEnv("GITHUB_TOKEN")
	if !exists {
		panic("Set GITHUB_TOKEN")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 1000},
		State:       "closed",
	}
	for {
		prs, resp, err := client.PullRequests.List(
			"kubernetes", "kubernetes", opt,
		)
		if err != nil {
			panic(err.Error())
		}
		for _, pr := range prs {
			handlePR(pr, client, db)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage

	}
}
