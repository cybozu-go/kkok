package client

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type routesListCommand struct{}

func (c routesListCommand) SetFlags(f *flag.FlagSet) {}

func (c routesListCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	data, err := Call(ctx, "GET", "/routes", nil)
	if err != nil {
		return handleError(err)
	}

	var ids []string
	err = json.Unmarshal(data, &ids)
	if err == nil {
		for _, rid := range ids {
			fmt.Println(rid)
		}
	}
	return handleError(err)
}

// RoutesListCommand implements "routes list" subcommand.
func RoutesListCommand() subcommands.Command {
	return subcmd{
		routesListCommand{},
		"list",
		"show route IDs",
		`list:
    Show route IDs.
`}
}
