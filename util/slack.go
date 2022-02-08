package util

import (
	"os"

	"github.com/slack-go/slack"
)

var (
	SLACK_TOKEN        = os.Getenv("SLACK_TOKEN")
	SlackClient        *slack.Client
	BroadCastChannelID string
	SlackURL           string
)
