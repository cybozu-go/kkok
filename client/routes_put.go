package client

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/subcommands"
)

type routesPutCommand struct{}

func (c routesPutCommand) SetFlags(f *flag.FlagSet) {}

func (c routesPutCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	args := f.Args()
	infile := os.Stdin

	switch len(args) {
	case 2:
		g, err := os.Open(args[1])
		if err != nil {
			return handleError(err)
		}
		defer g.Close()
		infile = g
	case 1:
		// ok
	default:
		f.Usage()
		return subcommands.ExitUsageError
	}
	id := args[0]

	data, err := ioutil.ReadAll(infile)
	if err != nil {
		return handleError(err)
	}

	_, err = Call(ctx, "PUT", "/routes/"+id, data)
	if err == nil {
		fmt.Println("success")
	}
	return handleError(err)
}

// RoutesPutCommand implements "routes put" subcommand.
func RoutesPutCommand() subcommands.Command {
	return subcmd{
		routesPutCommand{},
		"put",
		"add or update a route",
		`put ID [FILENAME]:

    Create a new route with ID, or edit the existing route matching ID.
    If FILENAME is given, it should be a JSON file.  If FILENAME is not
    given, JSON will be read from stdin.

The JSON must be an array of objects.  Each object must have "type"
field that speficies the type of transport such as "email" or "slack".

An example JSON may look like:

[
    {
        "type": "email",
        "from": "kkok@example.org",
        "to": ["ymmt2005@example.org"]
    },
    {
        "type": "slack",
        "url": "https://hooks.slack.com/xxxxxxx"
    }
]
`}
}
