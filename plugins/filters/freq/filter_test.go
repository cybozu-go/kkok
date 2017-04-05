package freq

import (
	"reflect"
	"testing"
	"time"

	"github.com/cybozu-go/kkok"
)

func testCalc(t *testing.T) {
	t.Parallel()

	f1 := newFilter()
	f1.Init("f1", nil)
	f1.duration = 10 * time.Millisecond
	a := &kkok.Alert{}

	f1.calc(a, time.Now())
	freq, ok := a.Stats["f1"]
	if !ok {
		t.Fatal("no stats")
	}
	if freq != float64(1)/10 {
		t.Error(`freq != float64(1)/10`)
	}

	f1.calc(a, time.Now())
	freq, ok = a.Stats["f1"]
	if !ok {
		t.Fatal("no stats")
	}
	if freq != float64(2)/10 {
		t.Error(`freq != float64(2)/10`)
	}

	time.Sleep(20 * time.Millisecond)

	f1.calc(a, time.Now())
	freq, ok = a.Stats["f1"]
	if !ok {
		t.Fatal("no stats")
	}
	if freq != float64(1)/10 {
		t.Error(`freq != float64(1)/10`)
	}
}

func testProcess(t *testing.T) {
	t.Parallel()

	f1 := newFilter()
	f1.Init("f1", map[string]interface{}{"if": "alert.Title == 'foo'"})
	f1.key = "test"

	alerts := []*kkok.Alert{
		{Title: "foo"},
		{Title: "bar"},
		{Title: "zot"},
	}

	alerts, err := f1.Process(alerts)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := alerts[0].Stats["test"]; !ok {
		t.Error(`_, ok := alerts[0].Stats["test"]; !ok`)
	}

	f2 := newFilter()
	f2.Init("f2", map[string]interface{}{"all": true, "if": "alerts.length>10"})

	// mismatch
	alerts = []*kkok.Alert{
		{Title: "foo"},
		{Title: "bar"},
		{Title: "zot"},
	}
	alerts, err = f2.Process(alerts)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range alerts {
		if len(a.Stats) != 0 {
			t.Error(`len(a.Stats) != 0`)
			t.Log(a)
		}
	}

	f3 := newFilter()
	f3.Init("f3", map[string]interface{}{"all": true, "if": "alerts.length>2"})
	f3.key = "test"

	// match
	alerts, err = f3.Process(alerts)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range alerts {
		if _, ok := a.Stats["test"]; !ok {
			t.Error(`_, ok := a.Stats["test"]; !ok`)
			t.Log(a)
		}
	}
}

func testParamsDefault(t *testing.T) {
	t.Parallel()

	f1 := newFilter()
	pp := f1.Params()

	if pp.Type != filterType {
		t.Error(`pp.Type != filterType`)
	}
	if pp.Params["duration"].(int) != int(defaultDuration.Seconds()) {
		t.Error(`pp.Params["duration"].(int) != defaultDuration`)
	}
	if pp.Params["divisor"].(float64) != defaultDivisor {
		t.Error(`pp.Params["divisor"].(float64) != defaultDivisor`)
	}
	if _, ok := pp.Params["classify"]; ok {
		t.Error(`_, ok :- pp.Params["classify"]; ok`)
	}
	if _, ok := pp.Params["key"]; ok {
		t.Error(`_, ok :- pp.Params["key"]; ok`)
	}
}

func testParamsExplicit(t *testing.T) {
	t.Parallel()

	f1 := newFilter()
	f1.duration = 2 * time.Second
	f1.divisor = 1.0
	f1.cl = clTitle
	f1.key = "hoge"

	pp := f1.Params()

	if pp.Type != "freq" {
		t.Error(`pp.Type != "freq"`)
	}

	params := map[string]interface{}{
		"duration": 2,
		"divisor":  1.0,
		"classify": "Title",
		"key":      "hoge",
	}
	if !reflect.DeepEqual(pp.Params, params) {
		t.Error(`!reflect.DeepEqual(pp.Params, params)`)
		t.Logf("%#v", pp.Params)
		t.Logf("%#v", params)
	}
}

func testParams(t *testing.T) {
	t.Run("Default", testParamsDefault)
	t.Run("Explicit", testParamsExplicit)
}

func TestFilter(t *testing.T) {
	t.Run("Calc", testCalc)
	t.Run("Process", testProcess)
	t.Run("Params", testParams)
}
