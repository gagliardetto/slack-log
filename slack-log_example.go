package slacklog

import (
	"os"

	"github.com/gagliardetto/slack-log/emoji"
)

func example_Say_test() {
	conf := &Config{
		Channel:  "#general",
		Username: "system-name",
		HookURL:  "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
	}

	client, err := New(conf)
	if err != nil {
		return
	}

	err = client.Say("something happened!")
	if err != nil {
		return
	}
}

func example_Shout_test() {
	conf := &Config{
		Channel:  "#general",
		Username: "system-name",
		HookURL:  "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
	}

	client, err := New(conf)
	if err != nil {
		return
	}

	err = client.Shout("something happened!")
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

	msg := NewLogMessage().
		Log15(
			"Look at this",
			"key", "value",
		)

	if err := client.Send(msg); err != nil {
		panic(err)
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

	msg := NewLogMessage().
		WithField("id", "userID").
		WithField("email", "user@example.com").
		WithField("paying", true).
		SetTitle("New user!")

	if err := client.Send(msg); err != nil {
		panic(err)
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

	msg := NewLogMessage().
		WithFields(Fields{
			{Key: "id", Val: "userID"},
			{Key: "email", Val: "user@example.com"},
			{Key: "paying", Val: true},
		}).
		SetTitle("New user!")

	if err := client.Send(msg); err != nil {
		panic(err)
	}
}
func example_all() {
	conf := &Config{
		Channel:  "#random",
		Username: "system-name",
		HookURL:  os.ExpandEnv("$SLACK_WEBHOOK"),
	}

	client, err := New(conf)
	if err != nil {
		return
	}

	{
		err = client.Say("Say: Hello world! " + emoji.Blush)
		if err != nil {
			panic(err)
		}

		err = client.Shout("Shout: Hello world! " + emoji.Blush)
		if err != nil {
			panic(err)
		}
	}
	{
		msg := &Message{}

		header := NewHeader().
			Text(">>> Header: Hello world").
			Emoji(true)

		msg.AddBlock(header)

		section1 := NewSection().
			AddField(
				NewField().
					Type(PLAIN_TEXT).
					Text("Foo"),
			)

		msg.AddBlock(section1)

		section2 := NewSection().
			AddField(
				NewField().
					Type(MARKDOWN).
					Text(
						"*Hello* **world**",
						"This is another line.",
					),
			)

		msg.AddBlock(section2)
		msg.AddBlock(NewDivider())

		img := NewImage().
			URL("https://i1.wp.com/thetempest.co/wp-content/uploads/2017/08/The-wise-words-of-Michael-Scott-Imgur-2.jpg?w=1024&ssl=1").
			Alt("This is a quote").
			Title(
				NewHeader().
					Emoji(true).
					Text("An inspiring quote for *you*"),
			)
		msg.AddBlock(img)
		msg.AddBlock(NewDivider())

		if err := msg.Validate(); err != nil {
			panic(err)
		}

		err = client.Send(msg)
		if err != nil {
			panic(err)
		}
	}
	{
		msg := NewLogMessage().
			Log15(
				"Log15: Look at this",
				"key", "value",
			)

		if err := client.Send(msg); err != nil {
			panic(err)
		}
	}
	{
		msg := NewLogMessage().
			SetTitle("Message: Hello world!").
			WithField("id", "userID").
			WithField("email", "user@example.com").
			WithField("paying", true)

		if err := client.Send(msg); err != nil {
			panic(err)
		}
	}
	{
		msg := NewLogMessage().
			SetTitle("WithFields: New user!").
			WithFields(Fields{
				{Key: "id", Val: "userID"},
				{Key: "email", Val: "user@example.com"},
				{Key: "paying", Val: true},
			})

		if err := client.Send(msg); err != nil {
			panic(err)
		}
	}
}
