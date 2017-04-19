package group

import (
	"time"

	"github.com/cybozu-go/kkok"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

const (
	filterType = "group"

	defaultFrom  = "filter:"
	defaultTitle = "merged alert"
)

type filter struct {
	kkok.BaseFilter

	// constants
	by      *otto.Script
	origBy  string
	from    string
	title   string
	message string
	routes  []string
}

func (f *filter) Params() kkok.PluginParams {
	m := make(map[string]interface{})
	if len(f.origBy) > 0 {
		m["by"] = f.origBy
	}
	if len(f.from) > 0 {
		m["from"] = f.from
	}
	if len(f.title) > 0 {
		m["title"] = f.title
	}
	if len(f.message) > 0 {
		m["message"] = f.message
	}
	if len(f.routes) > 0 {
		m["routes"] = f.routes
	}

	f.BaseFilter.AddParams(m)

	return kkok.PluginParams{
		Type:   filterType,
		Params: m,
	}
}

func (f *filter) mergeAlerts(alerts []*kkok.Alert) *kkok.Alert {
	if len(alerts) == 0 {
		return nil
	}

	a0 := alerts[0]
	if len(alerts) == 1 {
		return a0
	}

	newSub := make([]*kkok.Alert, len(alerts))
	copy(newSub, alerts)
	newAlert := &kkok.Alert{
		From:    a0.From,
		Date:    time.Now().UTC(),
		Host:    a0.Host,
		Title:   a0.Title,
		Message: a0.Message,
		Routes:  f.routes,
		Sub:     newSub,
	}

	for i, a := range alerts {
		if i == 0 {
			continue
		}

		if a.From != a0.From {
			if len(f.from) == 0 {
				newAlert.From = defaultFrom + f.ID()
			} else {
				newAlert.From = f.from
			}
		}

		if a.Host != a0.Host {
			newAlert.Host = "localhost"
		}

		if a.Title != a0.Title {
			if len(f.title) == 0 {
				newAlert.Title = defaultTitle
			} else {
				newAlert.Title = f.title
			}
		}

		if a.Message != a0.Message {
			newAlert.Message = f.message
		}
	}

	return newAlert
}

func (f *filter) Process(alerts []*kkok.Alert) ([]*kkok.Alert, error) {
	var newAlerts []*kkok.Alert
	var groups map[interface{}][]*kkok.Alert

	addGroup := func(a *kkok.Alert) error {
		var key interface{}
		if f.by != nil {
			v, err := f.BaseFilter.EvalAlert(a, f.by)
			if err != nil {
				return err
			}
			key, err = v.Export()
			if err != nil {
				return err
			}
		}
		groups[key] = append(groups[key], a)
		return nil
	}

	if f.BaseFilter.All() {
		ok, err := f.BaseFilter.IfAll(alerts)
		if err != nil {
			return nil, errors.Wrap(err, "group:"+f.ID())
		}
		if !ok {
			return alerts, nil
		}

		groups = make(map[interface{}][]*kkok.Alert)
		for _, a := range alerts {
			err = addGroup(a)
			if err != nil {
				return nil, errors.Wrap(err, "group:"+f.ID())
			}
		}
	} else {
		for _, a := range alerts {
			ok, err := f.BaseFilter.If(a)
			if err != nil {
				return nil, errors.Wrap(err, "group:"+f.ID())
			}

			if !ok {
				newAlerts = append(newAlerts, a)
				continue
			}

			if groups == nil {
				groups = make(map[interface{}][]*kkok.Alert)
			}
			err = addGroup(a)
			if err != nil {
				return nil, errors.Wrap(err, "group:"+f.ID())
			}
		}
	}

	for _, g := range groups {
		newAlerts = append(newAlerts, f.mergeAlerts(g))
	}
	return newAlerts, nil
}
