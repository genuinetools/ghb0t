package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/oauth2"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = "ghb0t - %s\n"
	// VERSION is the binary version.
	VERSION = "v0.1.0"
)

var (
	token   string
	seconds int64

	lastChecked time.Time

	debug   bool
	version bool
)

func init() {
	// parse flags
	flag.StringVar(&token, "token", "", "GitHub API token")
	flag.Int64Var(&seconds, "seconds", 30, "seconds to wait before checking for new events")

	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if version {
		fmt.Printf("%s", VERSION)
		os.Exit(0)
	}

	// set log level
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if token == "" {
		usageAndExit("GitHub token cannot be empty.", 1)
	}
}

func main() {
	// On ^C, or SIGTERM handle exit.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for sig := range c {
			logrus.Infof("Received %s, exiting.", sig.String())
			os.Exit(0)
		}
	}()

	// Create the http client.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	// Create the github client.
	client := github.NewClient(tc)

	// Get the authenticated user, the empty string being passed let's the GitHub
	// API know we want ourself.
	user, _, err := client.Users.Get("")
	if err != nil {
		logrus.Fatal(err)
	}
	username := *user.Login

	logrus.Infof("Bot started for user %s.", username)
	for {
		page := 1
		perPage := 20
		if err := getNotifications(client, username, page, perPage); err != nil {
			logrus.Fatal(err)
		}

		// Timeout for the interval.
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}

// getNotifications iterates over all the notifications received by a user.
func getNotifications(client *github.Client, username string, page, perPage int) error {
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

	notifications, resp, err := client.Activity.ListNotifications(opt)
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		// handle event
		if err := handleNotification(client, notification, username); err != nil {
			return err
		}
	}

	// Return early if we are on the last page.
	if page == resp.LastPage || resp.NextPage == 0 {
		return nil
	}

	page = resp.NextPage
	return getNotifications(client, username, page, perPage)
}

func handleNotification(client *github.Client, notification *github.Notification, username string) error {
	// Check if the type is a pull request.
	if *notification.Subject.Type == "PullRequest" {
		// Let's get some information about the pull request.
		parts := strings.Split(*notification.Subject.URL, "/")
		last := parts[len(parts)-1]
		id, err := strconv.Atoi(last)
		if err != nil {
			return err
		}

		pr, _, err := client.PullRequests.Get(*notification.Repository.Owner.Login, *notification.Repository.Name, int(id))
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
			// Never delete the master branch or a branch we do not own.
			if owner == username && branch != "master" {
				_, err := client.Git.DeleteRef(username, *pr.Head.Repo.Name, "heads/"+*pr.Head.Ref)
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

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
