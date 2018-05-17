package lib

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/google/go-github/github"
)

// Format is Formatter
type Format struct {
	ctx    context.Context
	client *github.Client
	debug  bool
}

// NewFormat is an initializer
func NewFormat(ctx context.Context, client *github.Client, debug bool) *Format {
	return &Format{ctx: ctx, client: client, debug: debug}
}

// Line is line infomation
type Line struct {
	title    string
	repoName string
	url      string
	user     string
}

// NewLineByIssue is an initializer by Issue
func NewLineByIssue(repoName string, issue github.Issue) Line {
	return Line{
		title:    *issue.Title,
		repoName: repoName,
		url:      *issue.HTMLURL,
		user:     *issue.User.Login,
	}
}

// NewLineByPullRequest is an initializer by PR
func NewLineByPullRequest(repoName string, pr github.PullRequest) Line {
	return Line{
		title:    *pr.Title,
		repoName: repoName,
		url:      *pr.HTMLURL,
		user:     *pr.User.Login,
	}
}

// Line returns Issue/PR info retrieving from GitHub
func (f *Format) Line(event *github.Event, i int) Line {
	payload := event.Payload()
	var line Line

	switch *event.Type {
	case "IssuesEvent":
		e := payload.(*github.IssuesEvent)
		issue := getIssue(f.ctx, f.client, *event.Repo.Name, *e.Issue.Number)

		if issue != nil {
			line = NewLineByIssue(*event.Repo.Name, *issue)
		} else {
			line = NewLineByIssue(*event.Repo.Name, *e.Issue)
		}
	case "IssueCommentEvent":
		e := payload.(*github.IssueCommentEvent)
		issue := getIssue(f.ctx, f.client, *event.Repo.Name, *e.Issue.Number)

		if issue != nil {
			if issue.PullRequestLinks == nil {
				line = NewLineByIssue(*event.Repo.Name, *issue)
			} else {
				pr := getPullRequest(f.ctx, f.client, *event.Repo.Name, *e.Issue.Number)
				line = NewLineByPullRequest(*event.Repo.Name, *pr)
			}
		} else {
			line = NewLineByIssue(*event.Repo.Name, *e.Issue)
		}
	case "PullRequestEvent":
		e := payload.(*github.PullRequestEvent)
		pr := getPullRequest(f.ctx, f.client, *event.Repo.Name, e.GetNumber())
		if pr != nil {
			line = NewLineByPullRequest(*event.Repo.Name, *pr)
		} else {
			line = NewLineByPullRequest(*event.Repo.Name, *e.PullRequest)
		}
	case "PullRequestReviewCommentEvent":
		e := payload.(*github.PullRequestReviewCommentEvent)
		pr := getPullRequest(f.ctx, f.client, *event.Repo.Name, *e.PullRequest.Number)
		if pr != nil {
			line = NewLineByPullRequest(*event.Repo.Name, *pr)
		} else {
			line = NewLineByPullRequest(*event.Repo.Name, *e.PullRequest)
		}
	}

	return line
}

func getIssue(ctx context.Context, client *github.Client, repoFullName string, number int) *github.Issue {
	owner, repo := getOwnerRepo(repoFullName)
	issue, _, _ := client.Issues.Get(ctx, owner, repo, number)
	return issue
}

func getPullRequest(ctx context.Context, client *github.Client, repoFullName string, number int) *github.PullRequest {
	owner, repo := getOwnerRepo(repoFullName)
	pr, _, _ := client.PullRequests.Get(ctx, owner, repo, number)
	return pr
}

func getOwnerRepo(repoFullName string) (string, string) {
	s := strings.Split(repoFullName, "/")
	owner := s[0]
	repo := s[1]
	return owner, repo
}

// All returns all lines which are formatted and sorted
func (f *Format) All(lines Lines) (string, error) {
	var result, prevRepoName, currentRepoName string

	sort.Sort(lines)

	for _, line := range lines {
		currentRepoName = line.repoName

		if currentRepoName != prevRepoName {
			prevRepoName = currentRepoName
			result += fmt.Sprintf("\n%s\n\n", currentRepoName)
		}

		result += fmt.Sprintf("%s\n%s\n\n", line.title, line.url)
	}

	return result, nil
}

// Lines has sort.Interface
type Lines []Line

func (l Lines) Len() int {
	return len(l)
}

func (l Lines) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l Lines) Less(i, j int) bool {
	return l[i].url < l[j].url
}
