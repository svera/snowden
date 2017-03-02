package logic

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
)

type senderInterface interface {
	PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error)
}

type githubInterface interface {
	ListFiles(ctx context.Context, owner string, repo string, number int, options *github.ListOptions) ([]*github.CommitFile, *github.Response, error)
	Get(ctx context.Context, owner string, repo string, number int) (*github.PullRequest, *github.Response, error)
}
