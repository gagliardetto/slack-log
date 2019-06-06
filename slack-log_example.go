package slacklog

import (
	"github.com/gagliardetto/slack-log/emoji"
	flatcolors "github.com/gagliardetto/slack-log/flat-colors"
)

func example_simplesend_test() {
	conf := &Config{
		Channel:  "#general",
		Username: "system-name",
		HookURL:  "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
	}

	client, err := New(conf)
	if err != nil {
		return
	}

	err = client.SimpleSend("something happened!", emoji.OpenMouth)
	if err != nil {
		return
	}
}

func example_Log15_test() {
	conf := &Config{
		Channel:  "#general",
		Username: "system-name",
		HookURL:  "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
	}

	client, err := New(conf)
	if err != nil {
		return
	}

	err = client.
		Log15(
			"Look at this",
			"key", "value",
		).
		SetColor(flatcolors.DarkerRed).
		SetIcon(emoji.TrollFace).
		Send()

	if err != nil {
		return
	}
}

func example_WithField_test() {
	conf := &Config{
		Channel:  "#general",
		Username: "system-name",
		HookURL:  "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
	}

	client, err := New(conf)
	if err != nil {
		return
	}

	err = client.
		WithField("id", "userID").
		WithField("email", "user@example.com").
		WithField("paying", true).
		SetMessage("New user!").
		SetColor(flatcolors.DarkerRed).
		SetIcon(emoji.TrollFace).
		Send()

	if err != nil {
		return
	}
}

func example_WithFields_test() {
	conf := &Config{
		Channel:  "#general",
		Username: "system-name",
		HookURL:  "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
	}

	client, err := New(conf)
	if err != nil {
		return
	}

	err = client.
		WithFields(ContextFields{
			"id":     "userID",
			"email":  "user@example.com",
			"paying": true,
		}).
		SetMessage("New user!").
		SetColor(flatcolors.DarkerRed).
		SetIcon(emoji.TrollFace).
		Send()

	if err != nil {
		return
	}
}
