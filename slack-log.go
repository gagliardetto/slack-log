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
	colorful "github.com/lucasb-eyer/go-colorful"
)

// curl -X POST \
// --data-urlencode 'payload={"channel": "#general", "username": "origin-system-name", "text": "something happened!", "icon_emoji": ":open_mouth:"}' \
// https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX

// ClientInterface is the interface that implements all the methods of the slack-log client
type ClientInterface interface {
	Send() error
	SimpleSend(text, iconEmoji string) error

	SetColor(hex string) ClientInterface
	SetMessage(msg string) ClientInterface
	SetIcon(iconEmoji string) ClientInterface
	OverrideChannel(channel string) ClientInterface

	Log15(msg string, ctx ...interface{}) ClientInterface
	WithField(key string, value interface{}) ClientInterface
	WithFields(fields ContextFields) ClientInterface
}

// Config is used to config the client for a Slack hook
type Config struct {
	Channel       string
	Username      string
	HookURL       string
	DefaultFields map[string]interface{}
}

// Client is a client of a Slack hook
type Client struct {
	channel       string
	username      string
	hookURL       *url.URL
	httpClient    *http.Client
	defaultFields map[string]interface{}

	tmp *temporaryData // temporaryData is reset on Send()
}

type temporaryData struct {
	ctx             *contextFields
	color           string
	message         string
	iconEmoji       string
	channelOverride string
}

func newTmp() *temporaryData {
	return &temporaryData{
		ctx:   newCtx(),
		color: "",
	}
}
func (td *temporaryData) reset() {
	td.ctx = newCtx()
	td.color = ""
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

	if conf.DefaultFields != nil {
		c.defaultFields = conf.DefaultFields
	}

	c.httpClient = utils.NewHTTPClient()

	return c, nil
}

// SetColor sets the color of the lateral bar of the message
func (c *Client) SetColor(hex string) ClientInterface {
	if c.tmp == nil {
		c.tmp = newTmp()
	}

	c.tmp.color = hex
	return c
}

// OverrideChannel overrides the default channel used to send the messages to.
func (c *Client) OverrideChannel(channel string) ClientInterface {
	if c.tmp == nil {
		c.tmp = newTmp()
	}

	c.tmp.channelOverride = channel
	return c
}

// SetMessage sets the main message text
func (c *Client) SetMessage(msg string) ClientInterface {
	if c.tmp == nil {
		c.tmp = newTmp()
	}

	c.tmp.message = msg
	return c
}

// SetIcon sets the icon for the slack message; you can find some in the slack-log/emoji package
func (c *Client) SetIcon(iconEmoji string) ClientInterface {
	if c.tmp == nil {
		c.tmp = newTmp()
	}

	c.tmp.iconEmoji = iconEmoji
	return c
}

// Send sends the message to the slack channel
func (c *Client) Send() error {
	var message string

	if c.tmp == nil {
		c.tmp = newTmp()
	}

	if c.defaultFields != nil {
		for k, v := range c.defaultFields {
			// TODO: overwrite if already set ???
			c.tmp.ctx.fields[k] = v
		}
	}

	if c.tmp.ctx == nil || c.tmp.ctx.len() == 0 {
		message = c.tmp.message
	} else {
		message = c.tmp.ctx.fields.String()
	}

	attachment := map[string]interface{}{
		"text":      message,
		"mrkdwn_in": []interface{}{"text"},
	}
	if c.tmp.color != "" {
		// Try parsing the color
		barColor, err := colorful.Hex(c.tmp.color)
		// If it is a valid color, use it.
		if err == nil {
			attachment["color"] = barColor.Hex()
		}
	}

	var channelName string
	if c.tmp.channelOverride != "" {
		channelName = c.tmp.channelOverride
	} else {
		channelName = c.channel
	}

	payload := map[string]interface{}{
		"attachments": []map[string]interface{}{
			attachment,
		},
		"text":       c.tmp.message,
		"channel":    channelName,
		"username":   c.username,
		"icon_emoji": c.tmp.iconEmoji,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
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

	// reset temporary data
	c.tmp = newTmp()

	return nil
}

// SimpleSend sends a message via a hook
func (c *Client) SimpleSend(text, iconEmoji string) error {
	return c.SetMessage(text).SetIcon(iconEmoji).Send()
}
