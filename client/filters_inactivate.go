package client

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"path"
	"time"

	"github.com/google/subcommands"
)

type filtersInactivateCommand struct{}

func (c filtersInactivateCommand) SetFlags(f *flag.FlagSet) {}

type inactivateData struct {
	Until time.Time `json:"until"`
}

func (c filtersInactivateCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	args := f.Args()
	if len(args) != 2 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	id := args[0]

	dt, err := time.Parse(time.RFC3339, args[1])
	if err != nil {
		duration, err := time.ParseDuration(args[1])
		if err != nil {
			return handleError(errors.New("wrong date or duration"))
		}
		dt = time.Now().Add(duration)
	}

	data := inactivateData{dt}
	_, err = Call(ctx, "PUT", path.Join("/filters", id, "inactivate"), data)
	if err == nil {
		fmt.Println("success")
	}
	return handleError(err)
}

// FiltersInactivateCommand implements "filters inactivate" subcommand.
func FiltersInactivateCommand() subcommands.Command {
	return subcmd{
		filtersInactivateCommand{},
		"inactivate",
		"inactivate a filter",
		`inactivate ID DATE_OR_DURATION:
    Inactivate a filter.  DATE_OR_DURATION must be either RFC3339
    date string or a duration that can be parsed by time.Duration.

    Ref: https://golang.org/pkg/time/#ParseDuration

Example:
    kkokc filters inactivate ID 30m
        Inactivate filter ID for 30 minutes.
    kkokc filters inactivate ID 2017-12-24T11:22:33Z
        Inactivate filter ID until Dec. 24th 2017 11:22:33 UTC.
`}
}
