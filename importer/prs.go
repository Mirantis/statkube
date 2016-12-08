package importer

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"

	"github.com/Mirantis/statkube/models"
)

// SleepTillFactory returns a SleepTill function that performs sleep to a given timestamp
func SleepTillFactory(now func() time.Time, sleep func(time.Duration)) func(time.Time) {
	return func(till time.Time) {
		sleep(till.Sub(now()))
	}

}

type PRScanner interface {
	More() bool
	Scan() *github.PullRequest
}

type GithubPRScanner struct {
	client   GHClientWrapper
	data     []*github.PullRequest
	i        int
	opt      *github.PullRequestListOptions
	nextPage int
	user     string
	repo     string
	limit    time.Time
	limitMet bool
}

// Implementation test...
var _ PRScanner = &GithubPRScanner{}

func CheckLimits(resp github.Response, client GithubProvider, sleepTill func(time.Time)) {
	if resp.Rate.Remaining > 0 {
		return
	}
	fmt.Printf("Sleeping till %s\n", resp.Rate.Reset)
	sleepTill(resp.Rate.Reset.Time)

	for {
		limit := client.GetLimits()
		if limit.Remaining > 0 {
			return
		}
		fmt.Printf("Sleeping more till %s\n", limit.Reset.Time)
		sleepTill(limit.Reset.Time)
	}
}

func (s *GithubPRScanner) More() bool {
	sleepTill := SleepTillFactory(time.Now, time.Sleep)
	if s.i == len(s.data) {
		if s.nextPage == 0 {
			return false
		}
		s.opt.ListOptions.Page = s.nextPage
		data, resp, err := s.client.client.PullRequests.List(
			s.user, s.repo, s.opt,
		)
		CheckLimits(*resp, s.client, sleepTill)
		if err != nil {
			fmt.Printf("%v", err.Error())
			return false
		}
		s.data = data
		s.i = 0
		s.nextPage = resp.NextPage
	}
	if s.data[s.i].UpdatedAt.Before(s.limit) {
		return false
	}
	return true
}

func (s *GithubPRScanner) Scan() *github.PullRequest {
	s.i++
	return s.data[s.i-1]
}

type GithubProvider interface {
	ListCommits(user string, repo string, pr int) ([]*github.RepositoryCommit, error)
	ListPRs(user string, repo string, limit time.Time) PRScanner
	GetLimits() *github.Rate
}

type GHClientWrapper struct {
	client *github.Client
}

// Implementation test...
var _ GithubProvider = &GHClientWrapper{}

func NewClient(token string) *GHClientWrapper {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return &GHClientWrapper{client: github.NewClient(tc)}
}

func (client GHClientWrapper) ListCommits(user string, repo string, prNo int) ([]*github.RepositoryCommit, error) {
	sleepTill := SleepTillFactory(time.Now, time.Sleep)
	commits, resp, err := client.client.PullRequests.ListCommits(
		user, repo, prNo, nil,
	)
	CheckLimits(*resp, client, sleepTill)
	if err != nil {
		return nil, err
	}
	return commits, nil
}

func (client GHClientWrapper) ListPRs(user, repo string, limit time.Time) PRScanner {
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 1000},
		State:       "closed",
		Sort:        "updated",
		Direction:   "desc",
	}
	return &GithubPRScanner{
		client: client, user: user, repo: repo, limit: limit, nextPage: 1, opt: opt,
	}
}

func (client GHClientWrapper) GetLimits() *github.Rate {
	limits, _, err := client.client.RateLimits()
	if err != nil {
		panic(err.Error())
	}
	return limits.Core
}

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
func getEmails(pr *github.PullRequest, client GithubProvider, user, repository string) map[string]struct{} {
	commits, _ := client.ListCommits(
		user, repository, *pr.Number,
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
	search := db.Joins("INNER JOIN developers ON developers.id = work_periods.developer_id").
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
		domain = curDomain
	}
	for {
		search := db.Joins("INNER JOIN domains ON domains.company_id = companies.id").
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

func deduceCompanyAndDev(pr *github.PullRequest, client GithubProvider, db *gorm.DB, user, repository string) (*models.Company, *models.Developer) {
	company, developer := deduceFromWorkPeriod(pr, db)
	if company != nil {
		return company, developer
	}
	emails := getEmails(pr, client, user, repository)
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

// HandlePR handles single pull request
func HandlePR(pr *github.PullRequest, client GithubProvider, repository *models.Repository, db *gorm.DB) {
	var prDB models.PullRequest
	if pr.MergedAt == nil {
		return
	}

	db.FirstOrInit(&prDB, models.PullRequest{Url: *pr.URL, RepositoryID: repository.ID})

	company, developer := deduceCompanyAndDev(pr, client, db, repository.User, repository.Repo)

	prDB.Company = *company
	prDB.Developer = *developer
	prDB.Created = *pr.CreatedAt
	prDB.Merged = pr.MergedAt
	db.Save(&prDB)
}
