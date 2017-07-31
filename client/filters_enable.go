package client

import (
	"context"
	"flag"
	"fmt"
	"path"

	"github.com/google/subcommands"
)

type filtersEnableCommand struct{}

func (c filtersEnableCommand) SetFlags(f *flag.FlagSet) {}

func (c filtersEnableCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	args := f.Args()
	if len(args) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	id := args[0]

	_, err := Call(ctx, "PUT", path.Join("/filters", id, "enable"), nil)
	if err == nil {
		fmt.Println("success")
	}
	return handleError(err)
}

// FiltersEnableCommand implements "filters enable" subcommand.
func FiltersEnableCommand() subcommands.Command {
	return subcmd{
		filtersEnableCommand{},
		"enable",
		"enable a filter",
		`enable ID:
    Enable a filter.
`}
}
