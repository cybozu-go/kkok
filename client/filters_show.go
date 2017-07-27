package client

import (
	"context"
	"encoding/json"
	"flag"

	"github.com/cybozu-go/kkok"
	"github.com/google/subcommands"
)

type filtersShowCommand struct {
	showJSON bool
}

func (c *filtersShowCommand) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.showJSON, "json", false, "output JSON")
}

func (c *filtersShowCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	args := f.Args()
	if len(args) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	id := args[0]

	data, err := Call(ctx, "GET", "/filters/"+id, nil)
	if err != nil {
		return handleError(err)
	}
	var params kkok.PluginParams
	err = json.Unmarshal(data, &params)
	if err == nil {
		printParams(params, c.showJSON)
	}
	return handleError(err)
}

// FiltersShowCommand implements "filters show" subcommand.
func FiltersShowCommand() subcommands.Command {
	return subcmd{
		&filtersShowCommand{},
		"show",
		"show filter parameters",
		`show [-json] ID:
    Show parameters of a filter.  ID is the filter ID.
    If -json is specified, parameters are shown in JSON.
`}
}
