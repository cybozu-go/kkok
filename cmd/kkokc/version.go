package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/cybozu-go/kkok"
	"github.com/google/subcommands"
)

const clientVersion = kkok.Version

type versionCommand struct{}

func (c versionCommand) SetFlags(f *flag.FlagSet) {}

func (c versionCommand) Execute(ctx context.Context, f *flag.FlagSet) error {
	fmt.Println("client version:", clientVersion)
	data, err := Call(ctx, "GET", "/version", nil)
	if err != nil {
		return err
	}
	fmt.Println("server version:", string(data))
	return nil
}

// VersionCommand implements "version" subcommand.
func VersionCommand() subcommands.Command {
	return subcmd{
		versionCommand{},
		"version",
		"print client/server versions",
		`version:
    Print client/server versions.
`,
	}
}
