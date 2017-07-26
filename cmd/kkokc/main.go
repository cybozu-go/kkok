package main

import (
	"context"
	"flag"
	"net/url"
	"os"

	"github.com/cybozu-go/cmd"
	"github.com/cybozu-go/log"
	sub "github.com/google/subcommands"
)

// TokenEnv is the environment variable name to specify kkok auth token.
const TokenEnv = "KKOK_TOKEN"

var (
	flgURL   = flag.String("url", "http://localhost:19898/", "URL of kkok server")
	flgToken = flag.String("token", os.Getenv(TokenEnv), "authentication token")
	kkokURL  *url.URL
)

func main() {
	sub.ImportantFlag("url")
	sub.ImportantFlag("token")
	sub.Register(sub.HelpCommand(), "misc")
	sub.Register(sub.FlagsCommand(), "misc")
	sub.Register(sub.CommandsCommand(), "misc")
	sub.Register(VersionCommand(), "server")
	flag.Parse()
	err := cmd.LogConfig{}.Apply()
	if err != nil {
		log.ErrorExit(err)
	}
	kkokURL, err = url.Parse(*flgURL)
	if err != nil {
		log.ErrorExit(err)
	}

	exitStatus := sub.ExitSuccess
	cmd.Go(func(ctx context.Context) error {
		exitStatus = sub.Execute(ctx)
		return nil
	})
	cmd.Stop()
	cmd.Wait()
	os.Exit(int(exitStatus))
}
