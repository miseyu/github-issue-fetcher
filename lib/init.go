package lib

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// Init initializes github-issue-fetcher settings
func Init() error {
	fmt.Print("** github-issue-fetcher Initialization **\n")

	if err := setUser(); err != nil {
		return err
	}

	time.Sleep(500 * time.Millisecond)
	ctx := context.Background()

	if err := setAccessToken(ctx); err != nil {
		return err
	}

	time.Sleep(500 * time.Millisecond)

	return nil
}

func setUser() error {
	fmt.Print(`
== [Step: 1/3] GitHub user ==

`)

	var msg string

	if _, err := getUser(); err == nil {
		msg = "Already initialized."
	} else {
		var user, answer string

		fmt.Print("What's your GitHub account? ")
		fmt.Scanln(&user)

		if len(user) >= 1 {
			fmt.Printf(`
The following command will be executed.

    $ git config --global github-issue-fetcher.user %s

`, user)

			fmt.Print("Are you sure? [y/n] ")
			fmt.Scanln(&answer)

			if string(answer[0]) != "y" {
				return errors.New("Canceled")
			}

			cmd := exec.Command("git", "config", "--global", "github-issue-fetcher.user", user)
			if err := cmd.Run(); err != nil {
				return err
			}

			msg = "Thanks!"
		}
	}

	fmt.Printf(`%s You can get it with the following command.

    $ git config --global github-issue-fetcher.user

`, msg)

	return nil
}

func setAccessToken(ctx context.Context) error {
	fmt.Print(`
== [Step: 2/3] GitHub personal access token ==

To get new token with ` + "`repo`" + ` scope, visit
https://github.com/settings/tokens/new

`)

	var msg string

	accessToken, err := getAccessToken()

	if err == nil {
		msg = "Already initialized."
	} else {
		var answer string

		fmt.Print("What's your GitHub personal access token? ")
		fmt.Scanln(&accessToken)

		if len(accessToken) >= 1 {
			fmt.Printf(`
The following command will be executed.

    $ git config --global github-issue-fetcher.token %s

`, accessToken)

			fmt.Print("Are you sure? [y/n] ")
			fmt.Scanln(&answer)

			if string(answer[0]) != "y" {
				return errors.New("Canceled")
			}

			cmd := exec.Command("git", "config", "--global", "github-issue-fetcher.token", accessToken)
			if err := cmd.Run(); err != nil {
				return err
			}

			msg = "Thanks!"
		}
	}

	fmt.Printf(`%s You can get it with the following command.

    $ git config --global github-issue-fetcher.token

`, msg)

	scopes, err := getClientScopes(ctx, getClient(ctx, accessToken))
	if err != nil {
		return err
	}

	if !isValidScopes(scopes) {
		return errors.New(`!!!! ` + "`repo`" + ` scopes are required. !!!!

You need personal access token which has ` + "`repo`" + `
scopes. Please add these scopes to your personal access
token, visit https://github.com/settings/tokens

`)

	}

	return nil
}

func isValidScopes(scopes []string) bool {
	var found1 bool

	for _, v := range scopes {
		if v == "repo" {
			found1 = true
		}
	}

	return found1
}
