package route

import (
	"math"
	"time"

	"github.com/cybozu-go/kkok"
)

type filter struct {
	kkok.BaseFilter

	// constants
	routes       []string
	replace      bool
	autoMute     bool
	muteDuration time.Duration
	muteRoutes   []string

	// states
	muteUntil time.Time
}

func (f *filter) Params() kkok.PluginParams {
	m := map[string]interface{}{
		"routes":       f.routes,
		"replace":      f.replace,
		"auto_mute":    f.autoMute,
		"mute_seconds": int(math.Trunc(f.muteDuration.Seconds())),
		"mute_routes":  f.muteRoutes,
	}

	f.BaseFilter.AddParams(m)

	return kkok.PluginParams{
		Type:   filterType,
		Params: m,
	}
}

func (f *filter) route(a *kkok.Alert) {
	routes := f.routes

	if f.autoMute {
		now := time.Now()
		if now.Before(f.muteUntil) {
			routes = f.muteRoutes
		} else {
			f.muteUntil = now.Add(f.muteDuration)
		}
	}

	if f.replace || len(a.Routes) == 0 {
		a.Routes = routes
	} else {
		m := make(map[string]struct{}, len(a.Routes))
		for _, r := range a.Routes {
			m[r] = struct{}{}
		}
		for _, r := range routes {
			if _, ok := m[r]; !ok {
				a.Routes = append(a.Routes, r)
			}
		}
	}
}

func (f *filter) Process(alerts []*kkok.Alert) ([]*kkok.Alert, error) {
	if f.BaseFilter.All() {
		ok, err := f.BaseFilter.EvalAllAlerts(alerts)
		if err != nil {
			return nil, err
		}

		if ok {
			for _, a := range alerts {
				f.route(a)
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
			f.route(a)
		}
	}
	return alerts, nil
}
