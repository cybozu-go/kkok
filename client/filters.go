package client

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

type filtersCommand struct{}

func (c filtersCommand) SetFlags(f *flag.FlagSet) {}

func (c filtersCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	newc := NewCommander(f, "kkokc filters")
	newc.Register(FiltersListCommand(), "")
	newc.Register(FiltersShowCommand(), "")
	return newc.Execute(ctx)
}

// FiltersCommand implements "filters" subcommand.
func FiltersCommand() subcommands.Command {
	return subcmd{
		filtersCommand{},
		"filters",
		"call /filters/... API",
		"filters ACTION ...",
	}
}
