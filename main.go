package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"github.com/svera/github-ping/config"
	"github.com/urfave/cli"
)

func loadConfig() (*config.Config, error) {
	var cfg *config.Config

	f, err := os.Open("/etc/webhook/ping.yml")
	if err != nil {
		return nil, errors.New("Couldn't load configuration file. Check that ping.yml exists and that it can be read. Exiting...")
	}
	if cfg, err = config.Load(f); err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	var cfg *config.Config
	var err error

	if cfg, err = loadConfig(); err != nil {
		fmt.Println(err.Error())
		return
	}

	app := cli.NewApp()
	api := slack.New(cfg.SlackToken)

	app.Name = "Github ping"
	app.Usage = "Pass the ping parameters to see them"
	app.Action = func(c *cli.Context) error {
		client := github.NewClient(nil)
		userName := c.Args().Get(0)
		if _, _, err := client.Users.Get(userName); err == nil {
			params := slack.PostMessageParameters{}
			_, _, err := api.PostMessage("@svera", "Hola mundo", params)
			if err != nil {
				fmt.Printf("%s\n", err)
			}
		}
		return nil
	}

	app.Run(os.Args)
}
