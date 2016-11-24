package main

import (
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

	githubClient := github.NewClient(nil)
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

		return lgc.Process(action, owner, repo, number, title, description)
	}

	app.Run(os.Args)
}
