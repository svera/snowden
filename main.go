package main

import (
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	api := slack.New("xoxp-2185617187-93896688804-105315722433-abd3e0ffb6c88dc578fd18df1dc82f5a")

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
