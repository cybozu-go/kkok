package maildir

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/cybozu-go/kkok"
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
	i, ok := params["dir"]
	if !ok {
		return nil, errors.New(`maildir: missing mandatory parmeter "dir"`)
	}

	dir, ok := i.(string)
	if !ok {
		return nil, errors.New(`maildir: dir must be a string`)
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

	if i, ok = params["interval"]; ok {
		switch ii := i.(type) {
		case int:
			interval = time.Duration(ii) * time.Second
		case float64:
			interval = time.Duration(int(ii)) * time.Second
		default:
			return nil, errors.New(`maildir: invalid interval type`)
		}
	}

	if interval <= 0 {
		return nil, errors.New(`maildir: invalid interval value`)
	}

	return &source{dir, interval}, nil
}

func init() {
	kkok.RegisterSource("maildir", ctor)
}
