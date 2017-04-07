package freq

import (
	"time"

	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/kkok/util"
	"github.com/pkg/errors"
)

const filterType = "freq"

func ctor(id string, params map[string]interface{}) (kkok.Filter, error) {
	f := newFilter()
	durationSeconds, err := util.GetInt("duration", params)
	switch {
	case err == nil:
		if durationSeconds <= 0 {
			return nil, errors.New("freq: invalid duration")
		}
		f.duration = time.Duration(durationSeconds) * time.Second
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "freq: duration")
	}

	divisor, err := util.GetFloat64("divisor", params)
	switch {
	case err == nil:
		if divisor <= 0 {
			return nil, errors.New("freq: invalid divisor")
		}
		f.divisor = divisor
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "freq: divisor")
	}

	s, err := util.GetString("foreach", params)
	switch {
	case err == nil:
		if len(s) == 0 {
			break
		}
		foreach, err := kkok.CompileJS(s)
		if err != nil {
			return nil, errors.Wrap(err, "freq: foreach")
		}
		f.foreach = foreach
		f.origForeach = s
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "freq: classify")
	}

	key, err := util.GetString("key", params)
	switch {
	case err == nil:
		f.key = key
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "freq: key")
	}

	err = f.Init(id, params)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func init() {
	kkok.RegisterFilter(filterType, ctor)
}
