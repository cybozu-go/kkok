package client

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

func (c versionCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	fmt.Println("client version:", clientVersion)
	data, err := Call(ctx, "GET", "/version", nil)
	if err == nil {
		fmt.Println("server version:", string(data))
	}
	return handleError(err)
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
