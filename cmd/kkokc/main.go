package main

import (
	"context"
	"flag"
	"net/url"
	"os"

	"github.com/cybozu-go/kkok/client"
	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
	sub "github.com/google/subcommands"
)

// TokenEnv is the environment variable name to specify kkok auth token.
const TokenEnv = "KKOK_TOKEN"

var (
	flgURL   = flag.String("url", "http://localhost:19898/", "URL of kkok server")
	flgToken = flag.String("token", os.Getenv(TokenEnv), "authentication token")
)

func main() {
	sub.ImportantFlag("url")
	sub.ImportantFlag("token")
	sub.Register(sub.HelpCommand(), "misc")
	sub.Register(sub.FlagsCommand(), "misc")
	sub.Register(sub.CommandsCommand(), "misc")
	sub.Register(client.VersionCommand(), "")
	sub.Register(client.AlertsCommand(), "")
	sub.Register(client.FiltersCommand(), "")
	sub.Register(client.RoutesCommand(), "")
	flag.Parse()
	err := well.LogConfig{}.Apply()
	if err != nil {
		log.ErrorExit(err)
	}
	u, err := url.Parse(*flgURL)
	if err != nil {
		log.ErrorExit(err)
	}
	client.Setup(u, *flgToken)

	exitStatus := sub.ExitSuccess
	well.Go(func(ctx context.Context) error {
		exitStatus = sub.Execute(ctx)
		return nil
	})
	well.Stop()
	well.Wait()
	os.Exit(int(exitStatus))
}
