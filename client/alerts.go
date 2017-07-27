package client

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

type alertsCommand struct{}

func (c alertsCommand) SetFlags(f *flag.FlagSet) {}

func (c alertsCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	newc := NewCommander(f, "kkokc alerts")
	newc.Register(AlertsListCommand(), "")
	newc.Register(AlertsPostCommand(), "")
	newc.Register(AlertsPostJSONCommand(), "")
	return newc.Execute(ctx)
}

// AlertsCommand implements "alerts" subcommand.
func AlertsCommand() subcommands.Command {
	return subcmd{
		alertsCommand{},
		"alerts",
		"call /alerts API",
		"alerts ACTION ...",
	}
}
