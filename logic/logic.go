package logic

import (
	"context"
	"fmt"
	"strings"

	"log"

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
		if l.cfg.Debug == true {
			log.Printf("%s repository is not being watched.\n", repo)
		}
		return nil
	}

	if action != "opened" && action != "reopened" {
		if l.cfg.Debug {
			log.Printf("PR %d has not been opened or reopened.\n", number)
		}
		return nil
	}

	var subs []string
	if files, _, err := l.gh.ListFiles(context.Background(), owner, repo, number, &github.ListOptions{}); err == nil {
		for _, file := range files {
			l.appendSubscribers(&subs, file.Filename, l.cfg.Watched[repo], owner)
		}
		if len(subs) == 0 {
			if l.cfg.Debug {
				log.Printf("No one is subscribed to files in PR %d \n", number)
			}
		}
		if err = l.notify(subs, owner, repo, number, title, description); err != nil {
			return err
		}

		return nil
	} else {
		if l.cfg.Debug {
			log.Printf("Error when trying to list Github repository files: %s", err.Error())
		}
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
			Title:     title,
			TitleLink: url,
			Text:      description,
		}
		params.Attachments = []slack.Attachment{attachment}
		_, _, err := l.sender.PostMessage(
			subscriber,
			fmt.Sprintf("Psssst... *%s* has opened a PR that affects files watched by you!", owner),
			params,
		)
		if err != nil {
			if l.cfg.Debug {
				log.Printf("Error notifying to subscriber %s: %s.\n", err.Error(), subscriber)
			}
		} else {
			if l.cfg.Debug {
				log.Printf("Notification sent to %s.\n", subscriber)
			}
		}
	}
	return nil
}
