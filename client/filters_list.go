package client

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type filtersListCommand struct{}

func (c filtersListCommand) SetFlags(f *flag.FlagSet) {}

func (c filtersListCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	data, err := Call(ctx, "GET", "/filters", nil)
	if err != nil {
		return handleError(err)
	}

	var filterIDs []string
	err = json.Unmarshal(data, &filterIDs)
	if err == nil {
		for _, fid := range filterIDs {
			fmt.Println(fid)
		}
	}
	return handleError(err)
}

// FiltersListCommand implements "filters list" subcommand.
func FiltersListCommand() subcommands.Command {
	return subcmd{
		filtersListCommand{},
		"list",
		"show filter IDs",
		`list:
    Show filter IDs.
`}
}
