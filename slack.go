package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

// SlackBot abstracts interactions with a slack chat
type SlackBot struct {
	api     *slack.Client
	channel string
}

// NewSlackBot constructs a new slack bot.
func NewSlackBot(token, channel string) *SlackBot {
	return &SlackBot{
		api:     slack.New(token),
		channel: channel,
	}
}

// Post posts a message to the slack chat.
func (b *SlackBot) Post(msg string) error {
	_, _, err := b.api.PostMessage(b.channel, msg, slack.PostMessageParameters{})
	if err != nil {
		return fmt.Errorf("can't post message to slack: %s", err)
	}

	return nil
}
