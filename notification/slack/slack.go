package slack

import (
	"github.com/slack-go/slack"

	"github.com/agflow/tools/agerr"
)

const (
	defaultChannel = "test-notifications"
	// ColorGood is slack's color "good"
	ColorGood = "good"
	// ColorWarning is slack's color "warning"
	ColorWarning = "warning"
	// ColorDanger is slack's color "danger"
	ColorDanger = "danger"
)

// Client is wrapper of a slack.Client
type Client struct {
	slackCli *slack.Client
	enabled  bool
}

// New return a new notifications/slack.Client
func New(token string, enabled bool) *Client {
	return &Client{slackCli: slack.New(token), enabled: enabled}
}

func getChannel(channel string) string {
	if channel == "" {
		return defaultChannel
	}
	return channel
}

// Send sends a notification message to slack
func (c *Client) Send(channel, msg string) error {
	if !c.enabled {
		return nil
	}

	channel = getChannel(channel)

	_, _, err := c.slackCli.PostMessage(
		channel,
		slack.MsgOptionAsUser(true),
		slack.MsgOptionText(msg, false))
	return agerr.Wrap("can't send slack notification: %w", err)
}

// SendWithColor sends a notification message to slack as an attachment with color
func (c *Client) SendWithColor(channel, msg, color string) error {
	if !c.enabled {
		return nil
	}

	channel = getChannel(channel)
	attachment := slack.Attachment{
		Text:  msg,
		Color: color,
	}

	_, _, err := c.slackCli.PostMessage(channel,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	return agerr.Wrap("can't send slack notification with color: %w", err)
}
