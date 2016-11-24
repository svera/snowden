# Snowden
This little tool allows users to be notified when one or more of their watched files or folders are included in a new or updated 
Github pull request, so they can review the changes made on them.

## Dependencies

Go > 1.5

This program is designed to work in conjunction with https://github.com/adnanh/webhook. Take a look at its documentation to know
more about how install and configure it.

You'll also need to register Snowden as a Slack application to obtain the credentials required to send messages.
Go to https://api.slack.com/apps?new_app=1 for more information.

## Usage

* Download an install Snowden: `go get github.com/svera/snowden`

* Create a configuration file `etc/webhook/snowden.yml`. You can follow the provided `snowden.sample.yml` as an example.

* Add the following to webhooks' `hooks.json` configuration file and launch it:
    ```
    {
        "id": "snowden",
        "execute-command": "/path/to/snowden/executable",
        "command-working-directory": "/path/of/working/directory",
        "response-message": "I got the payload!",
        "response-headers":
        [
            {
                "name": "Access-Control-Allow-Origin",
                "value": "*"
            }
        ],
        "pass-arguments-to-command":
        [
            {
                "source": "payload",
                "name": "action"
            },
            {
                "source": "payload",
                "name": "pull_request.head.repo.owner.login"
            },
            {
                "source": "payload",
                "name": "pull_request.head.repo.name"
            },      
            {
                "source": "payload",
                "name": "number"
            },
            {
                "source": "payload",
                "name": "pull_request.title"
            },
            {
                "source": "payload",
                "name": "pull_request.body"
            }            
        ]
    }
    ```

* Head to your Github repository settings and add a new webhook that triggers with pull request events, setting its URL 
to `<webhooks server domain>/hooks/snowden`.

You can also run it manually (without relying on webhooks server) for testing purposes. Snowden expects 6 parameters, in this order:

* The event that thrown the webhook. Snowden only responds to `opened` and `reopened`.
Read https://developer.github.com/v3/activity/events/types/ for more information.
* The pull request owner.
* The repository name where the pull request is aimed.
* The pull request number.
* The pull request title.
* The pull request description.
