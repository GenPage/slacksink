package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	flag "github.com/ogier/pflag"

	"github.com/nlopes/slack"
)

const usage = `Usage: slacksink [options]

slacksink reads data from stdin and posts it either directly
or as an attachment to Slack.
`

const (
	defaultUsername = "slacksink"
	envToken        = "SLACK_TOKEN"
	envUsername     = "SLACK_USERNAME"
	envFieldPrefix  = "SLACK_FIELD_"
)

var (
	attach   bool
	channel  string
	color    string
	icon     string
	message  string
	token    string
	username string
)

func main() {
	var inputData bytes.Buffer
	var body string

	parseArgs()

	client := slack.New(token)

	if _, err := io.Copy(&inputData, os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read from stdin: %s\n", err.Error())
		os.Exit(1)
	}

	params := slack.NewPostMessageParameters()
	if username != "" {
		params.Username = username
	}
	if icon != "" {
		params.IconURL = icon
	}

	if attach {
		attachment := slack.Attachment{
			Text: inputData.String(),
		}
		if color != "" {
			attachment.Color = color
		}
		attachment.Fields = getFields(os.Environ())
		params.Attachments = []slack.Attachment{attachment}
		body = message
	} else {
		body = inputData.String()
	}

	if _, _, err := client.PostMessage(channel, body, params); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to post message to Slack: %s\n", err.Error())
		os.Exit(1)
	}
}

func getFields(envs []string) []slack.AttachmentField {
	result := make([]slack.AttachmentField, 0, 2)
	for _, env := range envs {
		els := strings.SplitN(env, "=", 2)
		if strings.HasPrefix(els[0], envFieldPrefix) {
			f := slack.AttachmentField{
				Title: strings.TrimPrefix(els[0], envFieldPrefix),
				Value: els[1],
			}
			result = append(result, f)
		}
	}
	return result
}

func parseArgs() {
	flag.StringVar(&token, "token", "", "If no global token has been set, this is used. Uses $SLACK_TOKEN as fallback.")
	flag.StringVar(&channel, "channel", "", "The target channel or group or user")
	flag.StringVar(&message, "message", "", "Used as a header message while stdin becomes an attachment")
	flag.StringVar(&username, "username", "", "The username that is rendered as author of the message")
	flag.BoolVar(&attach, "attachment", false, "Send as attachment")
	flag.StringVar(&color, "color", "", "The color of the border (good, warning, danger, ...)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	if username == "" {
		username = os.Getenv(envUsername)
	}
	if username == "" {
		username = defaultUsername
	}

	if channel == "" {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "Please specify a channel.")
		os.Exit(1)
	}

	if token == "" {
		token = os.Getenv(envToken)
	}
	if token == "" {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "Please specify a token.")
		os.Exit(1)
	}
}
