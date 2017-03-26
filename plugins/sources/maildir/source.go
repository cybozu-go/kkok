package maildir

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/kkok/util"
	"github.com/cybozu-go/log"
	"github.com/pkg/errors"
)

const (
	defaultInterval = 10 // seconds
)

// source implements kkok.Source.
type source struct {
	dir      string
	interval time.Duration
}

func (s *source) Run(ctx context.Context, post func(*kkok.Alert)) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(s.interval):
		}

		for _, a := range scan(s.dir) {
			post(a)
		}
	}
}

func ctor(params map[string]interface{}) (kkok.Source, error) {
	dir, err := util.GetString("dir", params)
	if err != nil {
		return nil, errors.Wrap(err, "maildir: dir")
	}

	if !filepath.IsAbs(dir) {
		return nil, errors.New(`maildir: dir is not an absolute path`)
	}

	fi, err := os.Stat(dir)
	if err != nil {
		log.Warn("directory does not exist", map[string]interface{}{
			"source": "maildir",
			"dir":    dir,
		})
	} else {
		if !fi.IsDir() {
			return nil, errors.New(`maildir: dir is not a directory`)
		}
	}

	interval := time.Second * defaultInterval

	switch i, err := util.GetInt("interval", params); {
	case util.IsNotFound(err):
	case err != nil:
		return nil, errors.Wrap(err, "maildir: interval")
	default:
		interval = time.Duration(i) * time.Second
	}

	if interval <= 0 {
		return nil, errors.New(`maildir: invalid interval value`)
	}

	return &source{dir, interval}, nil
}

func init() {
	kkok.RegisterSource("maildir", ctor)
}
