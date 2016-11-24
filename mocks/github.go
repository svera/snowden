package mocks

import "github.com/google/go-github/github"

type Github struct {
	FakeCommitFiles []*github.CommitFile
	FakeResponse    *github.Response
	FakeError       error
}

func (g *Github) ListFiles(owner string, repo string, number int, options *github.ListOptions) ([]*github.CommitFile, *github.Response, error) {
	return g.FakeCommitFiles, g.FakeResponse, g.FakeError
}
