package kkok

import (
	"regexp"
	"sync"

	"github.com/cybozu-go/log"
	"github.com/pkg/errors"
)

var (
	reRouteID = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// Kkok is the struct to compose kkok.
//
// Internal APIs to work on generators/routes/filters are provided by this.
type Kkok struct {
	// lock for routes
	lkr sync.Mutex

	// routes maps route ID to a route (= list of transports).
	routes map[string][]Transport

	// lock for filters
	lkf sync.Mutex

	// filters are ordered as defined.
	filters []Filter
}

// NewKkok constructs a new empty Kkok.
func NewKkok() *Kkok {
	return &Kkok{
		routes:  make(map[string][]Transport),
		filters: make([]Filter, 0, 10),
	}
}

// RouteIDs return a slice of route IDs.
func (k *Kkok) RouteIDs() []string {
	k.lkr.Lock()
	defer k.lkr.Unlock()

	ids := make([]string, 0, len(k.routes))
	for id := range k.routes {
		ids = append(ids, id)
	}
	return ids
}

// Filters return a snapshot of the current filters.
func (k *Kkok) Filters() []Filter {
	k.lkf.Lock()
	defer k.lkf.Unlock()

	k.gc()
	filters := make([]Filter, len(k.filters))
	for i, f := range k.filters {
		filters[i] = f
	}
	return filters
}

// Handle implements AlertHandler interface.
func (k *Kkok) Handle(alerts []*Alert) {
	if len(alerts) == 0 {
		return
	}

	var err error
	for _, f := range k.Filters() {
		if f.Disabled() {
			continue
		}

		err = f.Reload()
		if err != nil {
			log.Error("[kkok] failed to reload", map[string]interface{}{
				log.FnError: err.Error(),
				"filter":    f.ID(),
			})
			return
		}

		alerts, err = f.Process(alerts)
		if err != nil {
			log.Error("[kkok] failed to filter alerts", map[string]interface{}{
				log.FnError: err.Error(),
				"filter":    f.ID(),
				"nalerts":   len(alerts),
			})
			return
		}

		if len(alerts) == 0 {
			log.Info("[kkok] filters reduced all alerts", map[string]interface{}{
				"filter": f.ID(),
			})
			return
		}
	}

	k.sendAlerts(alerts)
}

// AddRoute adds or replaces a route with id.
func (k *Kkok) AddRoute(id string, route []Transport) error {
	if !reRouteID.MatchString(id) {
		return errors.New("invalid route id: " + id)
	}

	k.lkr.Lock()
	k.routes[id] = route
	k.lkr.Unlock()
	return nil
}

// AddStaticFilter adds a filter statically.
func (k *Kkok) AddStaticFilter(filter Filter) error {
	k.lkf.Lock()
	defer k.lkf.Unlock()

	k.gc()

	id := filter.ID()

	for _, f := range k.filters {
		if f.ID() == id {
			return errors.New("duplicate filter id: " + id)
		}
	}

	k.filters = append(k.filters, filter)
	return nil
}

// PutFilter adds or replaces a filter with filter.ID().
func (k *Kkok) PutFilter(filter Filter) {
	k.lkf.Lock()
	defer k.lkf.Unlock()

	k.gc()

	id := filter.ID()

	for i, f := range k.filters {
		if f.ID() != id {
			continue
		}
		if f.Dynamic() {
			filter.SetDynamic()
		}
		k.filters[i] = filter
		return
	}

	filter.SetDynamic()
	k.filters = append(k.filters, filter)
}

// Internal APIs

func (k *Kkok) getRoute(id string) []Transport {
	k.lkr.Lock()
	defer k.lkr.Unlock()

	return k.routes[id]
}

func (k *Kkok) gc() {
	n := 0
	for _, f := range k.filters {
		if !f.Dynamic() || !f.Expired() {
			k.filters[n] = f
			n++
		}
	}
	k.filters = k.filters[0:n]
}

func (k *Kkok) removeFilter(id string) error {
	k.lkf.Lock()
	defer k.lkf.Unlock()

	n := 0
	for _, f := range k.filters {
		if f.ID() != id {
			if !f.Dynamic() || !f.Expired() {
				k.filters[n] = f
				n++
			}
			continue
		}

		if !f.Dynamic() {
			return errors.New("static filters cannot be removed")
		}
	}

	k.filters = k.filters[0:n]
	return nil
}

func (k *Kkok) getFilter(id string) Filter {
	k.lkf.Lock()
	defer k.lkf.Unlock()

	for _, f := range k.filters {
		if f.ID() == id {
			return f
		}
	}
	return nil
}

func (k *Kkok) sendAlerts(alerts []*Alert) {
	k.lkr.Lock()
	defer k.lkr.Unlock()

	routedAlerts := make(map[string][]*Alert)
	for id := range k.routes {
		routedAlerts[id] = make([]*Alert, 0)
	}

	for _, a := range alerts {
		for _, id := range a.Routes {
			_, ok := routedAlerts[id]
			if !ok {
				log.Warn("[kkok] unknown route", map[string]interface{}{
					"route": id,
					"title": a.Title,
				})
				continue
			}
			routedAlerts[id] = append(routedAlerts[id], a)
		}
	}

	for id, r := range k.routes {
		alerts = routedAlerts[id]
		if len(alerts) == 0 {
			continue
		}

		log.Info("[kkok] sending alerts", map[string]interface{}{
			"route":   id,
			"nalerts": len(alerts),
		})

		for _, t := range r {
			err := t.Deliver(alerts)
			if err != nil {
				log.Error("[kkok] failed to send alerts", map[string]interface{}{
					log.FnError: err.Error(),
					"route":     id,
					"transport": t.String(),
				})
			}
		}
	}
}
