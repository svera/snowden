package main

import (
	"errors"
	"fmt"
	"os"

	"strconv"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"github.com/svera/snowden/config"
	"github.com/svera/snowden/logic"
	"github.com/urfave/cli"
)

var cfg *config.Config

func loadConfig() (*config.Config, error) {
	f, err := os.Open("/etc/webhook/snowden.yml")
	if err != nil {
		return nil, errors.New("Couldn't load configuration file. Check that snowden.yml exists and that it can be read. Exiting...")
	}
	if cfg, err = config.Load(f); err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	var err error
	if cfg, err = loadConfig(); err != nil {
		fmt.Println(err.Error())
		return
	}

	githubClient := github.NewClient(nil)
	slackClient := slack.New(cfg.SlackToken)

	app := cli.NewApp()
	lgc := logic.New(slackClient, githubClient.PullRequests, cfg)

	app.Name = "Snowden"
	app.Usage = "Pass the ping parameters to see them"
	app.Action = func(c *cli.Context) error {
		owner := c.Args().Get(0)
		repo := c.Args().Get(1)
		number, _ := strconv.Atoi(c.Args().Get(2))

		lgc.Process(owner, repo, number)
		return nil
	}

	app.Run(os.Args)
}
