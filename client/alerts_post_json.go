package client

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/cybozu-go/kkok"
	"github.com/google/subcommands"
)

type alertsPostJSONCommand struct{}

func (c alertsPostJSONCommand) SetFlags(f *flag.FlagSet) {}

func (c alertsPostJSONCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	args := f.Args()
	infile := os.Stdin
	switch len(args) {
	case 0:
	case 1:
		g, err := os.Open(args[0])
		if err != nil {
			return handleError(err)
		}
		defer g.Close()
		infile = g
	default:
		f.Usage()
		return subcommands.ExitUsageError
	}

	dec := json.NewDecoder(infile)
	a := new(kkok.Alert)
	err := dec.Decode(a)
	if err != nil {
		return handleError(err)
	}
	err = a.Validate()
	if err != nil {
		return handleError(err)
	}

	_, err = Call(ctx, "POST", "/alerts", a)
	if err == nil {
		fmt.Println("success")
	}
	return handleError(err)
}

// AlertsPostJSONCommand implements "alerts postJSON" subcommand.
func AlertsPostJSONCommand() subcommands.Command {
	return subcmd{
		alertsPostJSONCommand{},
		"postJSON",
		"post JSON directly as a new alert",
		`postJSON [FILENAME]:
    Post a new alert.  If FILENAME is given, it should be a JSON file
    of the new alert.  If FILENAME is not given, JSON is read from stdin.

The JSON should be an object with these fields:

Name    | Required | Type   | Description
------- | -------- | ------ | ------------------------------------
From    | Yes      | string | Who sent this alert.
Title   | Yes      | string | One-line description of the alert.
Date    | No       | string | RFC3339 format date string.
Host    | No       | string | Where this alert was generated.
Message | No       | string | Multi-line description of the alert.
Info    | No       | object | Additional fields.
`}
}
