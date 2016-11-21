package logic

import (
	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
)

type senderInterface interface {
	PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error)
}

type githubInterface interface {
	ListFiles(owner string, repo string, number int, options *github.ListOptions) ([]*github.CommitFile, *github.Response, error)
}
