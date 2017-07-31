package client

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

type routesCommand struct{}

func (c routesCommand) SetFlags(f *flag.FlagSet) {}

func (c routesCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	newc := NewCommander(f, "routes")
	newc.Register(RoutesListCommand(), "")
	newc.Register(RoutesShowCommand(), "")
	newc.Register(RoutesPutCommand(), "")
	return newc.Execute(ctx)
}

// RoutesCommand implements "routes" subcommand.
func RoutesCommand() subcommands.Command {
	return subcmd{
		routesCommand{},
		"routes",
		"call /routes/... API",
		"routes ACTION ...",
	}
}
