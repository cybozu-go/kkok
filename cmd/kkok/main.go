package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/cybozu-go/cmd"
	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/log"
)

const (
	defaultConfigPath = "/etc/kkok.toml"
)

var (
	configPath = flag.String("f", defaultConfigPath, "Configuration file name")
	flgVersion = flag.Bool("v", false, "Print version and exit")
)

func main() {
	flag.Parse()

	if *flgVersion {
		fmt.Println(kkok.Version)
		return
	}

	cfg := kkok.NewConfig()
	_, err := toml.DecodeFile(*configPath, cfg)
	if err != nil {
		log.ErrorExit(err)
	}

	err = cfg.Log.Apply()
	if err != nil {
		log.ErrorExit(err)
	}

	k := kkok.NewKkok()

	// register routes
	for id, pl := range cfg.Routes {
		r := make([]kkok.Transport, len(pl))
		for i, p := range pl {
			t, err := kkok.NewTransport(p.Type, p.Params)
			if err != nil {
				log.ErrorExit(err)
			}
			r[i] = t
		}
		err = k.AddRoute(id, r)
		if err != nil {
			log.ErrorExit(err)
		}
	}

	// register filters
	idMap := make(map[string]struct{})
	for _, p := range cfg.Filters {
		f, err := kkok.NewFilter(p.Type, p.Params)
		if err != nil {
			log.ErrorExit(err)
		}
		id := f.ID()
		if _, ok := idMap[id]; ok {
			fmt.Fprintln(os.Stderr, "duplicate filter ID: "+id)
			os.Exit(1)
		}
		idMap[id] = struct{}{}
		k.AddFilter(f)
	}

	// start dispatcher
	d := kkok.NewDispatcher(cfg.InitialDuration(), cfg.MaxDuration(), k)
	cmd.Go(d.Run)

	for _, p := range cfg.Sources {
		src, err := kkok.NewSource(p.Type, p.Params)
		if err != nil {
			log.ErrorExit(err)
		}
		cmd.Go(func(ctx context.Context) error {
			return src.Run(ctx, d.Post)
		})
	}

	// start API server
	s, err := kkok.NewHTTPServer(cfg.Addr, cfg.APIToken, k, d)
	if err != nil {
		log.ErrorExit(err)
	}
	err = s.ListenAndServe()
	if err != nil {
		log.ErrorExit(err)
	}

	// all done.  wait for completion.
	err = cmd.Wait()
	if err != nil && !cmd.IsSignaled(err) {
		log.ErrorExit(err)
	}
}
