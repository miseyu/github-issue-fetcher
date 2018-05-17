package lib

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Settings has configure
type Settings struct {
	Format struct {
		Subject string
		Line    string
	}
	URL string
}

func getUser() (string, error) {
	if os.Getenv("GITHUB_USER") != "" {
		return os.Getenv("GITHUB_USER"), nil
	}

	output, _ := exec.Command("git", "config", "github-issue-fetcher.user").Output()

	if len(output) >= 1 {
		return strings.TrimRight(string(output), "\n"), nil
	}

	errText := `!!!! GitHub User required. Please execute the following command. !!!!

    $ github-issue-fetcher init`

	return "", errors.New(errText)
}

func getAccessToken() (string, error) {
	if os.Getenv("GITHUB_ACCESS_TOKEN") != "" {
		return os.Getenv("GITHUB_ACCESS_TOKEN"), nil
	}

	output, _ := exec.Command("git", "config", "github-issue-fetcher.token").Output()

	if len(output) >= 1 {
		return strings.TrimRight(string(output), "\n"), nil
	}

	errText := `!!!! GitHub Personal access token required. Please execute the following command. !!!!

    $ github-issue-fetcher init`

	return "", errors.New(errText)
}

func getParallelNum() (int, error) {
	return 5, nil
}

func getClient(ctx context.Context, accessToken string) *github.Client {
	sts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return github.NewClient(oauth2.NewClient(ctx, sts))
}

func getClientScopes(ctx context.Context, client *github.Client) ([]string, error) {
	_, response, err := client.Users.Get(ctx, "")
	return strings.Split(response.Header.Get("X-OAuth-Scopes"), ", "), err
}
