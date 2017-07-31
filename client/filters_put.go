package client

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/subcommands"
)

type filtersPutCommand struct{}

func (c filtersPutCommand) SetFlags(f *flag.FlagSet) {}

func (c filtersPutCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
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

	_, err = Call(ctx, "PUT", "/filters/"+id, data)
	if err == nil {
		fmt.Println("success")
	}
	return handleError(err)
}

// FiltersPutCommand implements "filters put" subcommand.
func FiltersPutCommand() subcommands.Command {
	return subcmd{
		filtersPutCommand{},
		"put",
		"add or update a filter",
		`put ID [FILENAME]:

    Create a new filter with ID, or edit the existing filter matching ID.
    If FILENAME is given, it should be a JSON file.  If FILENAME is not
    given, JSON will be read from stdin.

The JSON must be an object with these fields:

Name     | Required | Type      | Description
-------- | -------- | --------- | ---------------------------------
type     | Yes      | string    | Filter type such as discard or group.
disabled | No       | bool      | If true, the filter will not be used.
all      | No       | bool      | If true, the filter works for all alerts.
if       | No       | see below | Filter condition.
expire   | No       | string    | RFC3339 date string.

Other fields may be used depending on the filter type.

The default values of "disabled" and "all" are false.

"if" may be either a string of JavaScript boolean expression to
test an alert or an array of alerts should be filtered, or an array
of strings to invoke an external command.

For JavaScript expressions, when "all" is false, the filter will
assign each alert as "alert" variable and evaluate the JavaScript
expression.  When "all" is true, the filter will assign an array of
all alerts as "alerts".

For external commands, the filter executes the command by passing
the array of strings to os/exec.Command.  If the command exits
successfully, the filter work for the alerts.  When "all" is false,
the filter feeds a JSON object representing an alert via stdin.

Not all filters can be configured by "all".
For example, group filter always works as if "all" is true.

If "expire" is given, the filter will automatically be removed
at the given date.
`}
}
