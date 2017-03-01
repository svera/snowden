package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"

	"strconv"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"github.com/svera/snowden/config"
	"github.com/svera/snowden/logic"
	"github.com/urfave/cli"
)

var cfg *config.Config

func loadConfig() (*config.Config, error) {
	var data []byte
	var err error
	if data, err = config.Load("/etc/webhook/snowden.yml"); err != nil {
		return nil, err
	}
	return config.Parse(data)
}

func main() {
	var err error
	if cfg, err = loadConfig(); err != nil {
		fmt.Println(err.Error())
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	githubClient := github.NewClient(tc)
	slackClient := slack.New(cfg.SlackToken)

	app := cli.NewApp()
	lgc := logic.New(slackClient, githubClient.PullRequests, cfg)

	app.Name = "Snowden"
	app.Usage = "Notify users when their watched files or folders are included in a Github pull request"
	app.Action = func(c *cli.Context) error {
		action := c.Args().Get(0)
		owner := c.Args().Get(1)
		repo := c.Args().Get(2)
		number, _ := strconv.Atoi(c.Args().Get(3))
		title := c.Args().Get(4)
		description := c.Args().Get(5)

		log.Printf("Calling Snowden with params: %s %s %s %d %s %s\n", action, owner, repo, number, title, description)
		if err := lgc.Process(action, owner, repo, number, title, description); err != nil {
			if cfg.Debug {
				log.Printf("Error: %s\n", err.Error())
			}
		}
		return err
	}

	app.Run(os.Args)
}
