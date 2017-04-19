package freq

import (
	"math"
	"time"

	"github.com/cybozu-go/kkok"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

const (
	filterType = "freq"

	defaultDuration = 600 * time.Second

	defaultDivisor = 10
)

type filter struct {
	kkok.BaseFilter

	// constants
	duration    time.Duration
	divisor     float64
	foreach     *otto.Script
	origForeach string
	key         string

	// states
	samples map[interface{}]*Sample
}

func newFilter() *filter {
	return &filter{
		duration: defaultDuration,
		divisor:  defaultDivisor,
		samples:  make(map[interface{}]*Sample),
	}
}

func (f *filter) Params() kkok.PluginParams {
	m := map[string]interface{}{
		"duration": int(math.Trunc(f.duration.Seconds())),
		"divisor":  f.divisor,
	}
	if len(f.origForeach) > 0 {
		m["foreach"] = f.origForeach
	}
	if len(f.key) > 0 {
		m["key"] = f.key
	}

	f.BaseFilter.AddParams(m)

	return kkok.PluginParams{
		Type:   filterType,
		Params: m,
	}
}

func (f *filter) calc(a *kkok.Alert, now time.Time) error {
	var v interface{}
	if f.foreach != nil {
		vv, err := f.BaseFilter.EvalAlert(a, f.foreach)
		if err != nil {
			return err
		}
		v, err = vv.Export()
		if err != nil {
			return err
		}
	}

	s, ok := f.samples[v]
	if !ok {
		s = NewSample(f.duration)
		f.samples[v] = s
	}

	s.Add(now)
	freq := float64(s.Count()) / f.divisor

	k := f.key
	if len(k) == 0 {
		k = f.ID()
	}
	a.SetStat(k, freq)
	return nil
}

func (f *filter) Process(alerts []*kkok.Alert) ([]*kkok.Alert, error) {
	now := time.Now()

	if f.BaseFilter.All() {
		ok, err := f.BaseFilter.IfAll(alerts)
		if err != nil {
			return nil, errors.Wrap(err, "freq:"+f.ID())
		}

		if ok {
			for _, a := range alerts {
				err = f.calc(a, now)
				if err != nil {
					return nil, errors.Wrap(err, "freq:"+f.ID())
				}
			}
		}
		return alerts, nil
	}

	for _, a := range alerts {
		ok, err := f.BaseFilter.If(a)
		if err != nil {
			return nil, errors.Wrap(err, "freq:"+f.ID())
		}

		if ok {
			err = f.calc(a, now)
			if err != nil {
				return nil, errors.Wrap(err, "freq:"+f.ID())
			}
		}
	}
	return alerts, nil
}
