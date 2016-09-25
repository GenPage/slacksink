# Slack Sink

This is a simple Slack CLI tool for sending data from stdin to a Slack channel,
group, or user.

## Installation

    $ go get github.com/zerok/slacksink

## Usage

    $ export SLACK_TOKEN=your-token
    $ echo "hello" | slacksink --attachment --channel="#team-channel" \
      --message="Some title"


## History

This is basically a complete rewrite of Steve Kaliski's stdslack tool. Our
use-cases seem to be too different from even trying to contribute some of these
changes back so I opted for a complete fork.
