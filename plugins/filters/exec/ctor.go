package exec

import (
	"time"

	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/kkok/util"
	"github.com/pkg/errors"
)

func ctor(id string, params map[string]interface{}) (kkok.Filter, error) {
	f := newFilter()

	command, err := util.GetStringSlice("command", params)
	if err != nil {
		return nil, errors.Wrap(err, "filter:"+id)
	}
	if len(command) == 0 {
		return nil, errors.New("filter:" + id + " empty command")
	}
	f.command = command

	ts, err := util.GetInt("timeout", params)
	switch {
	case err == nil:
		if ts <= 0 {
			return nil, errors.New("filter:" + id + " wrong timeout")
		}
		f.timeout = time.Duration(ts) * time.Second
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "filter:"+id)
	}

	err = f.Init(id, params)
	if err != nil {
		return nil, errors.Wrap(err, "filter:"+id)
	}

	return f, nil
}

func init() {
	kkok.RegisterFilter(filterType, ctor)
}
