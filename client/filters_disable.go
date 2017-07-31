package client

import (
	"context"
	"flag"
	"fmt"
	"path"

	"github.com/google/subcommands"
)

type filtersDisableCommand struct{}

func (c filtersDisableCommand) SetFlags(f *flag.FlagSet) {}

func (c filtersDisableCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	args := f.Args()
	if len(args) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	id := args[0]

	_, err := Call(ctx, "PUT", path.Join("/filters", id, "disable"), nil)
	if err == nil {
		fmt.Println("success")
	}
	return handleError(err)
}

// FiltersDisableCommand implements "filters disable" subcommand.
func FiltersDisableCommand() subcommands.Command {
	return subcmd{
		filtersDisableCommand{},
		"disable",
		"disable a filter",
		`disable ID:
    Disable a filter.
`}
}
