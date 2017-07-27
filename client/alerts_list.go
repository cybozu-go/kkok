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

type alertsListCommand struct {
	oneline bool
}

func (c *alertsListCommand) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.oneline, "oneline", false, "show alerts shortly")
}

func (c *alertsListCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	data, err := Call(ctx, "GET", "/alerts", nil)
	if err != nil {
		return handleError(err)
	}
	var alerts []*kkok.Alert
	err = json.Unmarshal(data, &alerts)
	if err != nil {
		return handleError(err)
	}

	if len(alerts) == 0 {
		fmt.Fprintln(os.Stderr, "no alerts")
	} else {
		for i, a := range alerts {
			c.printAlert(i, a)
		}
	}
	return handleError(nil)
}

func (c *alertsListCommand) printAlert(i int, a *kkok.Alert) {
	dt := a.Date.UTC().Format("2006-01-02T15:04:05.000")
	if c.oneline {
		fmt.Printf("%2d: %s %20s %s\n", i, dt, a.From+"@"+a.Host, a.Title)
		return
	}

	fmt.Printf(`#%d
Date: %s
From: %s
Host: %s
Title: %s
Message:
%s
Info: %+v

`, i, dt, a.From, a.Host, a.Title, a.Message, a.Info)
}

// AlertsListCommand implements "alerts list" subcommand.
func AlertsListCommand() subcommands.Command {
	return subcmd{
		&alertsListCommand{},
		"list",
		"show list of pending alerts",
		`list:
    Show list of pending alerts.
`}
}
