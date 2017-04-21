package exec

import (
	"time"

	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/kkok/util"
	"github.com/pkg/errors"
)

func ctor(params map[string]interface{}) (kkok.Transport, error) {
	cl, err := util.GetStringSlice("command", params)
	if err != nil {
		return nil, errors.Wrap(err, "exec: command")
	}
	if len(cl) == 0 {
		return nil, errors.New("exec: empty command")
	}

	tr := newTransport(cl[0], cl[1:]...)

	label, err := util.GetString("label", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "exec: label")
	}
	tr.label = label

	all, err := util.GetBool("all", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "exec: all")
	}
	tr.all = all

	ts, err := util.GetInt("timeout", params)
	switch {
	case err == nil:
		if ts < 0 {
			return nil, errors.New("exec: invalid timeout")
		}
		tr.timeout = time.Duration(ts) * time.Second
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "exec: timeout")
	}

	return tr, nil
}

func init() {
	kkok.RegisterTransport(transportType, ctor)
}
