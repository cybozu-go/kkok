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

type routesShowCommand struct {
	showJSON bool
}

func (c *routesShowCommand) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.showJSON, "json", false, "output JSON")
}

func (c *routesShowCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	args := f.Args()
	if len(args) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	id := args[0]

	data, err := Call(ctx, "GET", "/routes/"+id, nil)
	if err != nil {
		return handleError(err)
	}
	var ppl []kkok.PluginParams
	err = json.Unmarshal(data, &ppl)
	if err != nil {
		return handleError(err)
	}

	if c.showJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")
		enc.Encode(ppl)
		return handleError(nil)
	}

	for _, pp := range ppl {
		printParams(pp, false)
		fmt.Println()
	}
	return handleError(nil)
}

// RoutesShowCommand implements "routes show" subcommand.
func RoutesShowCommand() subcommands.Command {
	return subcmd{
		&routesShowCommand{},
		"show",
		"show transports in a route",
		`show [-json] ID:
    Show transports in a route.  ID is the route ID.
    If -json is specified, information are shown in JSON.
`}
}
