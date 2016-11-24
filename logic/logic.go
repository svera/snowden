package logic

import (
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"github.com/svera/snowden/config"
)

// Logic is a struct that contains all the code and properties needed to send a message to subscribers of files and folders
// of a specific repository.
type Logic struct {
	sender senderInterface
	gh     githubInterface
	cfg    *config.Config
}

// New returns a new instance of Logic.
func New(sender senderInterface, gh githubInterface, cfg *config.Config) *Logic {
	return &Logic{
		sender: sender,
		gh:     gh,
		cfg:    cfg,
	}
}

// Process will, given an owner of a repository, its name and a PR number, alert subscriber about a new PR being opened if
// it affects any of their watched files or folders.
func (l *Logic) Process(action string, owner string, repo string, number int, title string, description string) error {
	var err error

	if _, ok := l.cfg.Watched[repo]; !ok {
		return nil
	}

	if action != "opened" && action != "reopened" {
		return nil
	}

	var subs []string
	if files, _, err := l.gh.ListFiles(owner, repo, number, &github.ListOptions{}); err == nil {
		for _, file := range files {
			l.appendSubscribers(&subs, file.Filename, l.cfg.Watched[repo], owner)
		}
		if len(subs) > 0 {
			if err = l.notify(subs, owner, repo, number, title, description); err != nil {
				return err
			}
		}
		return nil
	}
	return err
}

func (l *Logic) appendSubscribers(subs *[]string, fileName *string, rules []config.Rule, owner string) {
	for _, rule := range rules {
		for _, exception := range rule.Exceptions {
			if exception == owner {
				return
			}
		}
		for _, name := range rule.Names {
			if strings.HasPrefix(*fileName, name) {
				for _, new := range rule.Subscribers {
					l.appendIfNotIn(subs, new)
				}
			}
		}
	}
}

func (l *Logic) appendIfNotIn(subs *[]string, subscriber string) {
	for _, v := range *subs {
		if v == subscriber {
			return
		}
	}
	*subs = append(*subs, subscriber)
}

func (l *Logic) notify(subs []string, owner string, repo string, number int, title string, description string) error {
	url := fmt.Sprintf("https://github.com/%s/%s/pull/%d", owner, repo, number)
	for _, subscriber := range subs {
		params := slack.PostMessageParameters{
			Markdown: true,
		}
		attachment := slack.Attachment{
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: title,
					Value: description,
				},
			},
		}
		params.Attachments = []slack.Attachment{attachment}
		_, _, err := l.sender.PostMessage(
			subscriber,
			fmt.Sprintf("Psssst... *%s* has opened a PR %s that affects files watched by you!", owner, url),
			params,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
