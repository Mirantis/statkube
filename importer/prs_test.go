package importer

import (
	"fmt"
	"testing"
	"time"

	"github.com/Mirantis/statkube/models"
	"github.com/google/go-github/github"
)

type staticPRScanner struct {
	prs []github.PullRequest
	i   int
}

func (s *staticPRScanner) More() bool {
	fmt.Printf("%d", s.i)
	return s.i < len(s.prs)
}

func (s *staticPRScanner) Scan() *github.PullRequest {
	s.i++
	return &s.prs[s.i-1]
}

type staticGithubProvider struct {
	prs     []github.PullRequest
	commits []*github.RepositoryCommit
}

func (s *staticGithubProvider) ListCommits(_, _ string, _ int) ([]*github.RepositoryCommit, error) {
	return s.commits, nil
}

func (s *staticGithubProvider) ListPRs(_, _ string, _ time.Time) PRScanner {
	return &staticPRScanner{prs: s.prs}
}

// TestByWorkPeriod tests finding user that has github_id and work_period found
func TestGetByWorkPeriod(t *testing.T) {
	db := loadData()
	var pr models.PullRequest
	var company models.Company
	repository := getK8sRepo(db)
	provider := &staticGithubProvider{
		prs:     []github.PullRequest{},
		commits: []*github.RepositoryCommit{{}},
	}
	one := 1
	closed := "closed"
	user1 := "user1"
	url := "http://github.com/kubernetes/kubernetes/pulls/1"
	created := time.Date(2016, time.November, 8, 12, 0, 0, 0, time.UTC)
	merged := time.Date(2016, time.November, 8, 13, 0, 0, 0, time.UTC)
	HandlePR(&github.PullRequest{
		ID:     &one,
		Number: &one,
		State:  &closed,
		User: &github.User{
			Login: &user1,
		},
		CreatedAt: &created,
		MergedAt:  &merged,
		URL:       &url,
	}, provider, repository, db)
	db.First(&pr)
	db.Model(&pr).Related(&company)
	if company.Name != "Mirantis" {
		t.Error(fmt.Sprintf("PR didn't get company by work period. Instead got: %s", company.Name))
	}
}

// TestByEmail tests finding user that has no github_id but a matching e-mail
func TestGetByEmail(t *testing.T) {
	db := loadData()
	var pr models.PullRequest
	var company models.Company
	repository := getK8sRepo(db)
	email := "user2@example.com"
	provider := &staticGithubProvider{
		prs: []github.PullRequest{},
		commits: []*github.RepositoryCommit{{
			Commit: &github.Commit{
				Author: &github.CommitAuthor{
					Email: &email,
				},
			},
		}},
	}
	one := 1
	closed := "closed"
	user2 := "user2"
	url := "http://github.com/kubernetes/kubernetes/pulls/1"
	created := time.Date(2015, time.April, 8, 12, 0, 0, 0, time.UTC)
	merged := time.Date(2015, time.April, 8, 13, 0, 0, 0, time.UTC)
	HandlePR(&github.PullRequest{
		ID:     &one,
		Number: &one,
		State:  &closed,
		User: &github.User{
			Login: &user2,
		},
		CreatedAt: &created,
		MergedAt:  &merged,
		URL:       &url,
	}, provider, repository, db)
	db.First(&pr)
	db.Model(&pr).Related(&company)
	if company.Name != "Intel" {
		t.Error(fmt.Sprintf("PR didn't get company by work period. Instead got: %s", company.Name))
	}
}

// TestByDomain tests finding user that is not in the DB, but his domain is matched
func TestGetByDomain(t *testing.T) {
	db := loadData()
	var pr models.PullRequest
	var company models.Company
	repository := getK8sRepo(db)
	email := "user3@mirantis.com"
	provider := &staticGithubProvider{
		prs: []github.PullRequest{},
		commits: []*github.RepositoryCommit{{
			Commit: &github.Commit{
				Author: &github.CommitAuthor{
					Email: &email,
				},
			},
		}},
	}
	one := 1
	closed := "closed"
	user3 := "user3"
	url := "http://github.com/kubernetes/kubernetes/pulls/1"
	created := time.Date(2015, time.April, 8, 12, 0, 0, 0, time.UTC)
	merged := time.Date(2015, time.April, 8, 13, 0, 0, 0, time.UTC)
	HandlePR(&github.PullRequest{
		ID:     &one,
		Number: &one,
		State:  &closed,
		User: &github.User{
			Login: &user3,
		},
		CreatedAt: &created,
		MergedAt:  &merged,
		URL:       &url,
	}, provider, repository, db)
	db.First(&pr)
	db.Model(&pr).Related(&company)
	if company.Name != "Mirantis" {
		t.Error(fmt.Sprintf("PR didn't get company by domain. Instead got: %s", company.Name))
	}
}

// TestGetIndependent tests *independent user
func TestGetIndependent(t *testing.T) {
	db := loadData()
	var pr models.PullRequest
	var company models.Company
	repository := getK8sRepo(db)
	email := "user4@example.com"
	provider := &staticGithubProvider{
		prs: []github.PullRequest{},
		commits: []*github.RepositoryCommit{{
			Commit: &github.Commit{
				Author: &github.CommitAuthor{
					Email: &email,
				},
			},
		}},
	}
	one := 1
	closed := "closed"
	user3 := "user3"
	url := "http://github.com/kubernetes/kubernetes/pulls/1"
	created := time.Date(2015, time.April, 8, 12, 0, 0, 0, time.UTC)
	merged := time.Date(2015, time.April, 8, 13, 0, 0, 0, time.UTC)
	HandlePR(&github.PullRequest{
		ID:     &one,
		Number: &one,
		State:  &closed,
		User: &github.User{
			Login: &user3,
		},
		CreatedAt: &created,
		MergedAt:  &merged,
		URL:       &url,
	}, provider, repository, db)
	db.First(&pr)
	db.Model(&pr).Related(&company)
	if company.Name != "*independent" {
		t.Error(fmt.Sprintf("The user is not *independent. Instead: %s", company.Name))
	}
}
