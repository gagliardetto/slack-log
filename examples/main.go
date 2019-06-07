package main

import (
	"os"

	slacklog "github.com/gagliardetto/slack-log"
	"github.com/gagliardetto/slack-log/emoji"
	flatcolors "github.com/gagliardetto/slack-log/flat-colors"
)

func main() {
	conf := &slacklog.Config{
		Channel:  "#general",
		Username: "system-name",
		//HookURL:  "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
		HookURL: os.Getenv("SLACK_WEBHOOK"),
	}

	client, err := slacklog.New(conf)
	if err != nil {
		panic(err)
	}

	err = client.
		WithFields(slacklog.ContextFields{
			{Key: "foo", Val: "bar"},
			{Key: "bool", Val: true},
			{Key: "int", Val: 12345},
			{Key: "float", Val: 4.5678},
		}).
		SetMessage("test message!").
		SetColor(flatcolors.LightBlue2).
		SetIcon(emoji.Smile).
		Send()

	if err != nil {
		panic(err)
	}
}
