package kkok

import (
	"context"
	"sync"
	"time"
)

// alertPool pools Alert objects for some duration.
type alertPool struct {
	mu     sync.Mutex
	alerts []*Alert
}

// put puts an alert into the pool.
func (p *alertPool) put(a *Alert) {
	p.mu.Lock()
	p.alerts = append(p.alerts, a)
	p.mu.Unlock()
}

// empty returns true if the pool is empty.
func (p *alertPool) empty() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.alerts) == 0
}

// peek returns a (deep) copy of currently pooled alerts.
func (p *alertPool) peek() []*Alert {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.alerts) == 0 {
		return nil
	}

	c := make([]*Alert, len(p.alerts))
	for i, a := range p.alerts {
		c[i] = a.Clone()
	}
	return c
}

// take returns pooled alerts and clears the pool.
func (p *alertPool) take() []*Alert {
	p.mu.Lock()
	defer p.mu.Unlock()

	c := p.alerts
	p.alerts = nil
	return c
}

// AlertHandler is an interface for NewDispatcher.
type AlertHandler interface {
	Handle([]*Alert)
}

// Dispatcher accepts and pools alerts then dispatches them periodically.
type Dispatcher struct {
	alertPool
	initInterval time.Duration
	maxInterval  time.Duration
	handler      AlertHandler
}

// NewDispatcher creates Dispatcher.
// init and max is the initial and maximum duration between dispatches.
// handler handles pooled alerts.  To start dispatching, invoke Run.
func NewDispatcher(init, max time.Duration, handler AlertHandler) *Dispatcher {
	if init <= 0 {
		init = 1 * time.Second
	}
	if max < init {
		max = init
	}
	return &Dispatcher{
		initInterval: init,
		maxInterval:  max,
		handler:      handler,
	}
}

// Post puts an alert into the pool.
func (d *Dispatcher) Post(a *Alert) {
	d.alertPool.put(a)
}

// Run starts dispatching alerts until ctx is canceled.
// This method will always return nil.
func (d *Dispatcher) Run(ctx context.Context) error {
	cur := d.initInterval

	for {
		select {
		case <-ctx.Done():
			// process pooled alerts before quit, if any.
			if d.alertPool.empty() {
				return nil
			}
		case <-time.After(cur):
		}

		alerts := d.alertPool.take()
		if len(alerts) == 0 {
			cur = d.initInterval
			continue
		}

		d.handler.Handle(alerts)

		cur = cur * 2
		if cur > d.maxInterval {
			cur = d.maxInterval
		}
	}
}
