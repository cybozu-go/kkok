package client

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type filtersDeleteCommand struct{}

func (c filtersDeleteCommand) SetFlags(f *flag.FlagSet) {}

func (c filtersDeleteCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	args := f.Args()
	if len(args) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	id := args[0]

	_, err := Call(ctx, "DELETE", "/filters/"+id, nil)
	if err == nil {
		fmt.Println("success")
	}
	return handleError(err)
}

// FiltersDeleteCommand implements "filters delete" subcommand.
func FiltersDeleteCommand() subcommands.Command {
	return subcmd{
		filtersDeleteCommand{},
		"delete",
		"delete a filter",
		`delete ID:
    Delete a filter.  ID is the filter ID.
`}
}
