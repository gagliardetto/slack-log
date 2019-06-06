package main

import (
	slacklog "github.com/gagliardetto/slack-log"
	"github.com/gagliardetto/slack-log/emoji"
	flatcolors "github.com/gagliardetto/slack-log/flat-colors"
)

func main() {
	conf := &slacklog.Config{
		Channel:  "#general",
		Username: "system-name",
		HookURL:  "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
	}

	client, err := slacklog.New(conf)
	if err != nil {
		panic(err)
	}

	err = client.
		WithFields(slacklog.ContextFields{
			"id":     "userID",
			"email":  "user@example.com",
			"paying": true,
		}).
		SetMessage("New user!").
		SetColor(flatcolors.LightBlue2).
		SetIcon(emoji.Smile).
		Send()

	if err != nil {
		panic(err)
	}
}
