package kkok

import "testing"

type dupFilter struct {
	BaseFilter
}

func (f *dupFilter) Params() PluginParams {
	p := PluginParams{
		Type:   "dup",
		Params: make(map[string]interface{}),
	}
	f.BaseFilter.AddParams(p.Params)
	return p
}

func (f *dupFilter) Process(alerts []*Alert) ([]*Alert, error) {
	r := make([]*Alert, len(alerts)*2)
	for i, a := range alerts {
		r[i*2] = a
		r[i*2+1] = a
	}
	return r, nil
}

type discardFilter struct {
	BaseFilter
}

func (f *discardFilter) Params() PluginParams {
	p := PluginParams{
		Type:   "discard",
		Params: make(map[string]interface{}),
	}
	f.BaseFilter.AddParams(p.Params)
	return p
}

func (f *discardFilter) Process(alerts []*Alert) ([]*Alert, error) {
	return nil, nil
}

type routeFilter struct {
	BaseFilter

	routes []string
}

func (f *routeFilter) Params() PluginParams {
	p := PluginParams{
		Type: "route",
		Params: map[string]interface{}{
			"routes": f.routes,
		},
	}
	f.BaseFilter.AddParams(p.Params)
	return p
}

func (f *routeFilter) Process(alerts []*Alert) ([]*Alert, error) {
	for _, a := range alerts {
		a.Routes = f.routes
	}
	return alerts, nil
}

type testTransport struct {
	alerts []*Alert
}

func (t *testTransport) Params() PluginParams {
	return PluginParams{
		Type:   "test",
		Params: make(map[string]interface{}),
	}
}

func (t *testTransport) String() string {
	return "test"
}

func (t *testTransport) Deliver(alerts []*Alert) error {
	t.alerts = alerts
	return nil
}

func testFilters(t *testing.T) {
	t.Parallel()

	f1 := &dupFilter{}
	f1.Init("f1", nil)
	f2 := &dupFilter{}
	f2.Init("f2", nil)

	k := NewKkok()
	if len(k.Filters()) != 0 {
		t.Error(`len(k.Filters()) != 0`)
	}

	err := k.AddStaticFilter(f1)
	if err != nil {
		t.Fatal(err)
	}
	if len(k.Filters()) != 1 {
		t.Error(`len(k.Filters()) != 1`)
	}

	if k.getFilter("f1") != Filter(f1) {
		t.Error(`k.getFilter("f1") != Filter(f1)`)
	}

	// replace
	k.PutFilter(f1)
	if len(k.Filters()) != 1 {
		t.Error(`len(k.Filters()) != 1`)
	}

	k.PutFilter(f2)
	if len(k.Filters()) != 2 {
		t.Error(`len(k.Filters()) != 2`)
	}

	if k.getFilter("f2") != Filter(f2) {
		t.Error(`k.getFilter("f2") != Filter(f2)`)
	}

	if err := k.removeFilter("f1"); err == nil {
		t.Fatal("static filter f1 should not be removed")
	}

	if k.getFilter("f1") != Filter(f1) {
		t.Error(`k.getFilter("f1") != Filter(f1)`)
	}

	if err := k.removeFilter("f2"); err != nil {
		t.Error(err)
	}

	if k.getFilter("f2") != nil {
		t.Error(`k.getFilter("f2") != nil`)
	}

	if len(k.Filters()) != 1 {
		t.Error(`len(k.Filters()) != 1`)
	}
}

func testRoutes(t *testing.T) {
	t.Parallel()

	tr1 := &testTransport{}
	tr2 := &testTransport{}

	k := NewKkok()

	if len(k.routes) != 0 {
		t.Error(`len(k.routes) != 0`)
	}

	k.AddRoute("r1", []Transport{tr1})
	if len(k.routes) != 1 {
		t.Error(`len(k.routes) != 1`)
	}

	// in fact, replace
	k.AddRoute("r1", []Transport{tr1})
	if len(k.routes) != 1 {
		t.Error(`len(k.routes) != 1`)
	}

	k.AddRoute("r2", []Transport{tr2})
	if len(k.routes) != 2 {
		t.Error(`len(k.routes) != 2`)
	}

	ids := k.RouteIDs()
	if len(ids) != 2 {
		t.Error(`len(ids) != 2`)
	}
}

func testHandleMultiFilters(t *testing.T) {
	t.Parallel()

	k := NewKkok()

	f1 := &dupFilter{}
	f1.Init("f1", nil)
	err := k.AddStaticFilter(f1)
	if err != nil {
		t.Fatal(err)
	}

	f2 := &discardFilter{}
	f2.Init("f2", nil)
	k.PutFilter(f2)

	tr1 := &testTransport{}
	k.AddRoute("r1", []Transport{tr1})

	a1 := &Alert{Routes: []string{"r1"}}
	k.Handle([]*Alert{a1})
	if len(tr1.alerts) != 0 {
		t.Error(len(tr1.alerts) != 0)
	}

	a2 := &Alert{Routes: []string{"r1"}}
	f2.Enable(false)
	k.Handle([]*Alert{a2})
	if len(tr1.alerts) != 2 {
		t.Error(len(tr1.alerts) != 2)
	}
}

func testHandleMultiRoutes(t *testing.T) {
	t.Parallel()

	k := NewKkok()

	// constructs a route filter which sets 'r1' to `routes` attribute of all alerts
	f := &routeFilter{}
	f.Init("f", nil)
	f.routes = []string{"r1"}
	err := k.AddStaticFilter(f)
	if err != nil {
		t.Fatal(err)
	}

	tr1 := &testTransport{}
	k.AddRoute("r1", []Transport{tr1})

	tr2 := &testTransport{}
	k.AddRoute("r2", []Transport{tr2})

	a1 := &Alert{Routes: []string{}}
	k.Handle([]*Alert{a1})
	if len(tr1.alerts) != 1 {
		t.Error("len(tr1.alerts) != 1")
	}
	if len(tr2.alerts) != 0 {
		t.Error("len(tr2.alerts) != 0")
	}
}

func TestKkok(t *testing.T) {
	t.Run("Filters", testFilters)
	t.Run("Routes", testRoutes)
	t.Run("HandleMultiFilters", testHandleMultiFilters)
	t.Run("HandleMultiRoutes", testHandleMultiRoutes)
}
