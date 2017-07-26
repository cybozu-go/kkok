package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
)

// Command is the interface for NewCommand
type Command interface {
	SetFlags(f *flag.FlagSet)
	Execute(ctx context.Context, f *flag.FlagSet) error
}

type subcmd struct {
	Command
	name     string
	synopsis string
	usage    string
}

func (c subcmd) Name() string {
	return c.name
}

func (c subcmd) Synopsis() string {
	return c.synopsis
}

func (c subcmd) Usage() string {
	return c.usage
}

func (c subcmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	err := c.Command.Execute(ctx, f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
