package github

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Client struct {
	Github *github.Client
}

func NewAuthenticatedClient(accessToken string) *Client {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &Client{
		client,
	}

}

func (client *Client) GetBranch(ctx context.Context, owner, repository, branchName string) (*github.Branch, bool) {

	branch, res, err := client.Github.Repositories.GetBranch(ctx, owner, repository, branchName)

	if err != nil {

		if hit, hourWait := blockIfRateLimitIsHit(err); hit {

			<-hourWait

			client.GetBranch(ctx, owner, repository, branchName)

		}

		if res.StatusCode == 404 {
			return nil, false
		}

		log.WithFields(log.Fields{
			"error": true,
		}).Fatal(err.Error())

	}

	return branch, true

}

func blockIfRateLimitIsHit(err error) (bool, <-chan time.Time) {

	if _, ok := err.(*github.RateLimitError); ok {

		log.Warn("Rate limt hit, waiting 1 hour")

		return true, time.After(time.Hour)

	}

	return false, nil

}
func getNamesFromBranches(branches []*github.Branch) []string {

	branchNames := make([]string, len(branches))

	for i, branch := range branches {

		branchNames[i] = branch.GetName()

	}
	return branchNames
}
