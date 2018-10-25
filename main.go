package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/genuinetools/ghb0t/version"
	"github.com/genuinetools/pkg/cli"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var (
	token    string
	interval time.Duration
	enturl   string

	lastChecked time.Time

	debug bool
)

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "ghb0t"
	p.Description = "A GitHub Bot to automatically delete your fork's branches after a pull request has been merged"

	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.StringVar(&token, "token", os.Getenv("GITHUB_TOKEN"), "GitHub API token (or env var GITHUB_TOKEN)")
	p.FlagSet.DurationVar(&interval, "interval", 30*time.Second, "check interval (ex. 5ms, 10s, 1m, 3h)")
	p.FlagSet.StringVar(&enturl, "url", "", "Connect to a specific GitHub server, provide full API URL (ex. https://github.example.com/api/v3/)")

	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		if token == "" {
			return fmt.Errorf("GitHub token cannot be empty")
		}

		return nil
	}

	// Set the main program action.
	p.Action = func(ctx context.Context, args []string) error {
		ticker := time.NewTicker(interval)

		// On ^C, or SIGTERM handle exit.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		go func() {
			for sig := range c {
				cancel()
				ticker.Stop()
				logrus.Infof("Received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

		// Create the http client.
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		// Create the github client.
		client := github.NewClient(tc)
		if enturl != "" {
			var err error
			client.BaseURL, err = url.Parse(enturl)
			if err != nil {
				logrus.Fatalf("failed to parse provided url: %v", err)
			}
		}

		// Get the authenticated user, the empty string being passed let's the GitHub
		// API know we want ourself.
		user, _, err := client.Users.Get(ctx, "")
		if err != nil {
			logrus.Fatal(err)
		}
		username := *user.Login

		logrus.Infof("Bot started for user %s.", username)

		for range ticker.C {
			page := 1
			perPage := 20
			if err := getNotifications(ctx, client, username, page, perPage); err != nil {
				logrus.Warn(err)
			}
		}
		return nil
	}

	// Run our program.
	p.Run()
}

// getNotifications iterates over all the notifications received by a user.
func getNotifications(ctx context.Context, client *github.Client, username string, page, perPage int) error {
	opt := &github.NotificationListOptions{
		All:   true,
		Since: lastChecked,
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}
	if lastChecked.IsZero() {
		lastChecked = time.Now()
	}

	notifications, resp, err := client.Activity.ListNotifications(ctx, opt)
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		// handle event
		if err := handleNotification(ctx, client, notification, username); err != nil {
			return err
		}
	}

	// Return early if we are on the last page.
	if page == resp.LastPage || resp.NextPage == 0 {
		return nil
	}

	page = resp.NextPage
	return getNotifications(ctx, client, username, page, perPage)
}

func handleNotification(ctx context.Context, client *github.Client, notification *github.Notification, username string) error {
	// Check if the type is a pull request.
	if *notification.Subject.Type == "PullRequest" {
		// Let's get some information about the pull request.
		parts := strings.Split(*notification.Subject.URL, "/")
		last := parts[len(parts)-1]
		id, err := strconv.Atoi(last)
		if err != nil {
			return err
		}

		pr, _, err := client.PullRequests.Get(ctx, *notification.Repository.Owner.Login, *notification.Repository.Name, id)
		if err != nil {
			return err
		}

		if *pr.State == "closed" && *pr.Merged {
			// If the PR was made from a repository owned by the current user,
			// let's delete it.
			branch := *pr.Head.Ref
			if pr.Head.Repo == nil {
				return nil
			}
			if pr.Head.Repo.Owner == nil {
				return nil
			}
			owner := *pr.Head.Repo.Owner.Login
			defaultBranch := *notification.Repository.DefaultBranch
			// Never delete the default branch or a branch we do not own.
			if owner == username && branch != defaultBranch {
				_, err := client.Git.DeleteRef(ctx, username, *pr.Head.Repo.Name, strings.Replace("heads/"+*pr.Head.Ref, "#", "%23", -1))
				// 422 is the error code for when the branch does not exist.
				if err != nil && !strings.Contains(err.Error(), " 422 ") {
					return err
				}
				logrus.Infof("Branch %s on %s/%s no longer exists.", branch, owner, *pr.Head.Repo.Name)
			}
		}
	}

	return nil
}
