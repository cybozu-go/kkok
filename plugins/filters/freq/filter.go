package freq

import (
	"math"
	"time"

	"github.com/cybozu-go/kkok"
)

const (
	defaultDuration = 600 * time.Second

	defaultDivisor = 10
)

type filter struct {
	kkok.BaseFilter

	// constants
	duration time.Duration
	divisor  float64
	cl       clType
	key      string

	// states
	samples map[string]*Sample
}

func newFilter() *filter {
	return &filter{
		duration: defaultDuration,
		divisor:  defaultDivisor,
		cl:       clNone,
		samples:  make(map[string]*Sample),
	}
}

func (f *filter) Params() kkok.PluginParams {
	m := map[string]interface{}{
		"duration": int(math.Trunc(f.duration.Seconds())),
		"divisor":  f.divisor,
	}
	cls := f.cl.String()
	if len(cls) > 0 {
		m["classify"] = cls
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

func (f *filter) calc(a *kkok.Alert, now time.Time) {
	v := f.cl.Value(a)
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
}

func (f *filter) Process(alerts []*kkok.Alert) ([]*kkok.Alert, error) {
	now := time.Now()

	if f.BaseFilter.All() {
		ok, err := f.BaseFilter.EvalAllAlerts(alerts)
		if err != nil {
			return nil, err
		}

		if ok {
			for _, a := range alerts {
				f.calc(a, now)
			}
		}
		return alerts, nil
	}

	for _, a := range alerts {
		ok, err := f.BaseFilter.EvalAlert(a)
		if err != nil {
			return nil, err
		}

		if ok {
			f.calc(a, now)
		}
	}
	return alerts, nil
}
