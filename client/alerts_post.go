package client

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cybozu-go/kkok"
	"github.com/google/subcommands"
)

type alertsPostCommand struct {
	from string
	host string
}

func (c *alertsPostCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.from, "from", "kkokc", "value of From field")
	hname, _ := os.Hostname()
	if len(hname) == 0 {
		hname = "localhost"
	}
	f.StringVar(&c.host, "host", hname, "value of Host field")
}

func (c *alertsPostCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	a := &kkok.Alert{
		From: c.from,
		Date: time.Now().UTC(),
		Host: c.host,
	}

	args := f.Args()
	switch len(args) {
	case 2:
		a.Message = args[1]
	case 1:
		fmt.Fprintln(os.Stderr, "reading message from stdin...")
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return handleError(err)
		}
		a.Message = string(data)
	default:
		f.Usage()
		return subcommands.ExitUsageError
	}
	a.Title = args[0]

	err := a.Validate()
	if err != nil {
		return handleError(err)
	}

	_, err = Call(ctx, "POST", "/alerts", a)
	if err == nil {
		fmt.Println("success")
	}
	return handleError(err)
}

// AlertsPostCommand implements "alerts post" subcommand.
func AlertsPostCommand() subcommands.Command {
	return subcmd{
		&alertsPostCommand{},
		"post",
		"post a new alert",
		`post [options] TITLE [MESSAGE]:
    Post a new alert.  If MESSAGE is not specified, it will be read
    from stdin.
`}
}
