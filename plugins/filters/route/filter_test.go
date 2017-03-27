package route

import (
	"reflect"
	"testing"
	"time"

	"github.com/cybozu-go/kkok"
)

func testRouteAdd(t *testing.T) {
	t.Parallel()

	f1 := &filter{
		routes:       []string{"r1"},
		muteDuration: 10 * time.Millisecond,
		muteRoutes:   []string{"r2"},
	}
	a := &kkok.Alert{}

	f1.route(a)
	if !reflect.DeepEqual(a.Routes, []string{"r1"}) {
		t.Error(`!reflect.DeepEqual(a.Routes, []string{"r1"})`)
	}

	// no duplicate, no muteRoutes
	f1.route(a)
	if !reflect.DeepEqual(a.Routes, []string{"r1"}) {
		t.Error(`!reflect.DeepEqual(a.Routes, []string{"r1"})`)
	}

	f2 := &filter{
		routes:       []string{"r2", "r1"},
		muteDuration: 10 * time.Millisecond,
	}

	f2.route(a)
	if !reflect.DeepEqual(a.Routes, []string{"r1", "r2"}) {
		t.Error(`!reflect.DeepEqual(a.Routes, []string{"r2", "r1"})`)
	}

	// no side effects
	if !reflect.DeepEqual(f1.routes, []string{"r1"}) {
		t.Error(`!reflect.DeepEqual(f1.routes, []string{"r1"})`)
	}
}

func testRouteReplace(t *testing.T) {
	t.Parallel()

	f1 := &filter{
		routes: []string{"r1"},
	}
	a := &kkok.Alert{}

	f1.route(a)
	if !reflect.DeepEqual(a.Routes, []string{"r1"}) {
		t.Fatal(`!reflect.DeepEqual(a.Routes, []string{"r1"})`)
	}

	f2 := &filter{
		routes:  []string{"r4"},
		replace: true,
	}
	f2.route(a)

	if !reflect.DeepEqual(a.Routes, []string{"r4"}) {
		t.Error(`!reflect.DeepEqual(a.Routes, []string{"r4"})`)
	}
}

func testRouteMute(t *testing.T) {
	f1 := &filter{
		routes:       []string{"r1"},
		autoMute:     true,
		muteDuration: 10 * time.Millisecond,
		muteRoutes:   []string{"r2", "r3"},
	}
	a1 := &kkok.Alert{}
	a2 := &kkok.Alert{}
	a3 := &kkok.Alert{}
	f1.route(a1)
	f1.route(a2)
	time.Sleep(11 * time.Millisecond)
	f1.route(a3)

	if !reflect.DeepEqual(a1.Routes, []string{"r1"}) {
		t.Error(`!reflect.DeepEqual(a1.Routes, []string{"r1"})`)
	}
	if !reflect.DeepEqual(a2.Routes, []string{"r2", "r3"}) {
		t.Error(`!reflect.DeepEqual(a1.Routes, []string{"r2", "r3"})`)
	}
	if !reflect.DeepEqual(a3.Routes, []string{"r1"}) {
		t.Error(`!reflect.DeepEqual(a2.Routes, []string{"r1"})`)
	}
}

func testRoute(t *testing.T) {
	t.Run("Add", testRouteAdd)
	t.Run("Replace", testRouteReplace)
	t.Run("Mute", testRouteMute)
}

func testProcess(t *testing.T) {
	t.Parallel()

	f1 := &filter{
		routes: []string{"r1", "r2"},
	}
	f1.Init("f1", map[string]interface{}{"if": "alert.Title == 'foo'"})

	alerts := []*kkok.Alert{
		{Title: "foo"},
		{Title: "bar"},
		{Title: "zot"},
	}

	alerts, err := f1.Process(alerts)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(alerts[0].Routes, []string{"r1", "r2"}) {
		t.Error(`!reflect.DeepEqual(alerts[0].Routes, []string{"r1", "r2"})`)
	}
	if len(alerts[1].Routes) != 0 {
		t.Error(`len(alerts[1].Routes) != 0`)
	}
	if len(alerts[2].Routes) != 0 {
		t.Error(`len(alerts[2].Routes) != 0`)
	}

	f2 := &filter{
		routes:  []string{"r1", "r2"},
		replace: true,
	}
	f2.Init("f2", map[string]interface{}{"all": true, "if": "alerts.length>10"})

	// mismatch
	alerts, err = f2.Process(alerts)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(alerts[0].Routes, []string{"r1", "r2"}) {
		t.Error(`!reflect.DeepEqual(alerts[0].Routes, []string{"r1", "r2"})`)
	}
	if len(alerts[1].Routes) != 0 {
		t.Error(`len(alerts[1].Routes) != 0`)
	}
	if len(alerts[2].Routes) != 0 {
		t.Error(`len(alerts[2].Routes) != 0`)
	}

	f3 := &filter{
		routes:  []string{"r3"},
		replace: true,
	}
	f3.Init("f3", map[string]interface{}{"all": true, "if": "alerts.length>2"})

	// match
	alerts, err = f3.Process(alerts)
	if err != nil {
		t.Fatal(err)
	}
	for i, a := range alerts {
		if !reflect.DeepEqual(a.Routes, []string{"r3"}) {
			t.Error(`!reflect.DeepEqual(a.Routes, []string{"r3"})`, i)
		}
	}
}

func testParams(t *testing.T) {
	t.Parallel()

	f1 := &filter{
		routes:       []string{"r1"},
		replace:      true,
		autoMute:     true,
		muteDuration: 2 * time.Second,
		muteRoutes:   []string{"r2"},
	}
	pp := f1.Params()

	if pp.Type != "route" {
		t.Error(`pp.Type != "route"`)
	}

	params := map[string]interface{}{
		"routes":       []string{"r1"},
		"replace":      true,
		"auto_mute":    true,
		"mute_seconds": 2,
		"mute_routes":  []string{"r2"},
	}
	if !reflect.DeepEqual(pp.Params, params) {
		t.Error(`!reflect.DeepEqual(pp.Params, params)`)
		t.Logf("%#v", pp.Params)
		t.Logf("%#v", params)
	}
}

func TestFilter(t *testing.T) {
	t.Run("Route", testRoute)
	t.Run("Process", testProcess)
	t.Run("Params", testParams)
}
