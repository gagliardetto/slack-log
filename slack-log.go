package slacklog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gagliardetto/slack-log/utils"
)

// curl -X POST \
// --data-urlencode 'payload={"channel": "#general", "username": "origin-system-name", "text": "something happened!", "icon_emoji": ":open_mouth:"}' \
// https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX

// MessageInterface is the interface that implements all the methods of the slack-log message.
type MessageInterface interface {
	Send() error
	Say(text string) error
	Shout(text string) error

	SetTitle(title string) MessageInterface

	Log15(title string, ctx ...interface{}) MessageInterface
	WithField(key string, value interface{}) MessageInterface
	WithFields(fields Fields) MessageInterface
}

// Config is used to config the client for a Slack hook
type Config struct {
	Channel  string
	Username string
	HookURL  string
}

// Client is a client of a Slack hook
type Client struct {
	channel    string
	username   string
	hookURL    *url.URL
	httpClient *http.Client
}

// New returns a new Client that can be used to send slack notifications
func New(conf *Config) (*Client, error) {
	if conf == nil {
		return nil, errors.New("conf is nil")
	}
	c := &Client{}

	if conf.Channel == "" {
		return nil, errors.New("Channel not specified")
	}
	c.channel = conf.Channel

	if conf.Username == "" {
		return nil, errors.New("Username not specified")
	}
	c.username = conf.Username

	u, err := url.Parse(conf.HookURL)
	if err != nil {
		return nil, fmt.Errorf("error while parsing HookURL: %v", err)
	}
	c.hookURL = u

	c.httpClient = utils.NewHTTPClient()

	return c, nil
}

var Debug bool

// Send sends the message to the slack channel
func (c *Client) Send(msg Msg) error {

	channelName := c.channel

	payload := M{
		"blocks":   msg.Blocks(),
		"channel":  channelName,
		"username": c.username,
		"text":     "hello there",
	}

	jsonPayload, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		return err
	}

	if Debug {
		fmt.Println(string(jsonPayload))
	}

	req, err := http.NewRequest("POST", c.hookURL.String(), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(jsonPayload)))

	response, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	var respText bytes.Buffer
	_, err = io.Copy(&respText, response.Body)
	if err != nil {
		return err
	}
	if respText.String() != "ok" {
		return fmt.Errorf("error from server: %v", respText.String())
	}

	return nil
}

// Shout sends a message via a hook
func (c *Client) Shout(text string) error {
	return c.Send(NewLogMessage().SetTitle(text))
}

// Say sends a message via a hook
func (c *Client) Say(text string) error {
	msg := &Message{}

	section := NewSection().
		AddField(
			NewField().
				Type(MARKDOWN).
				Text(text),
		)
	msg.AddBlock(section)
	if err := msg.Validate(); err != nil {
		return err
	}
	return c.Send(msg)
}
