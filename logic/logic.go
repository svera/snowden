package logic

import (
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"github.com/svera/snowden/config"
)

type Logic struct {
	sender senderInterface
	gh     githubInterface
	cfg    *config.Config
}

func New(sender senderInterface, gh githubInterface, cfg *config.Config) *Logic {
	return &Logic{
		sender: sender,
		gh:     gh,
		cfg:    cfg,
	}
}

func (l *Logic) Process(owner string, repo string, number int) {
	if _, ok := l.cfg.Watched[repo]; !ok {
		return
	}

	if f, _, err := l.gh.ListFiles(owner, repo, number, &github.ListOptions{}); err == nil {
		for _, file := range f {
			subs := l.subscribers(file.Filename, l.cfg.Watched[repo])
			if len(subs) > 0 {
				l.notify(subs, owner)
				return
			}
		}
	}
}

func (l *Logic) subscribers(fileName *string, watched []config.Watch) []string {
	for _, file := range watched {
		for _, name := range file.Name {
			if strings.HasPrefix(name, *fileName) {
				return file.Subscribers
			}
		}
	}
	return []string{}
}

func (l *Logic) notify(subs []string, owner string) {
	for _, subscriber := range subs {
		params := slack.PostMessageParameters{}
		_, _, err := l.sender.PostMessage(
			subscriber,
			fmt.Sprintf("User %s has opened a new PR that affects one or more files watched by you", owner),
			params,
		)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}

}
