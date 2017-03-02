package mocks

import (
	"context"

	"github.com/google/go-github/github"
)

// Github mocked
type Github struct {
	FakeCommitFiles []*github.CommitFile
	FakeResponse    *github.Response
	FakePullRequest *github.PullRequest
	FakeError       error
}

// ListFiles mocked
func (g *Github) ListFiles(ctx context.Context, owner string, repo string, number int, options *github.ListOptions) ([]*github.CommitFile, *github.Response, error) {
	return g.FakeCommitFiles, g.FakeResponse, g.FakeError
}

// Get mocked
func (g *Github) Get(ctx context.Context, owner string, repo string, number int) (*github.PullRequest, *github.Response, error) {
	return g.FakePullRequest, g.FakeResponse, g.FakeError
}
