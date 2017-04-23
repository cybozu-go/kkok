package edit

import (
	"reflect"
	"testing"
	"time"

	"github.com/cybozu-go/kkok"
)

func testConvert(t *testing.T, js string, alert, expect *kkok.Alert) {
	s, err := kkok.CompileJS(js)
	if err != nil {
		t.Fatal(err)
	}

	vm := kkok.NewVM()

	obj, err := toObject(alert)
	if err != nil {
		t.Fatal(err)
	}
	vm.Set("alert", obj.Value())
	_, err = vm.Run(s)
	if err != nil {
		t.Fatal(err)
	}
	a, err := fromObject(obj)
	if err != nil {
		if expect == nil {
			return
		}
		t.Fatal(err)
	}

	if len(a.Info) == 0 {
		a.Info = nil
	}
	if len(a.Stats) == 0 {
		a.Stats = nil
	}

	if !reflect.DeepEqual(a, expect) {
		t.Error(`!reflect.DeepEqual(a, expect)`)
		t.Logf("%#v", a)
	}
}

func TestConvert(t *testing.T) {
	newAlert := func() *kkok.Alert {
		return &kkok.Alert{
			From:    "from",
			Date:    time.Date(2011, 2, 3, 4, 5, 6, 123000000, time.UTC),
			Host:    "localhost",
			Title:   "title",
			Message: "こんにちは",
		}
	}
	newAlertFrom := func(from string) *kkok.Alert {
		a := newAlert()
		a.From = from
		return a
	}
	newAlertHost := func(host string) *kkok.Alert {
		a := newAlert()
		a.Host = host
		return a
	}
	newAlertTitle := func(title string) *kkok.Alert {
		a := newAlert()
		a.Title = title
		return a
	}
	newAlertMessage := func(msg string) *kkok.Alert {
		a := newAlert()
		a.Message = msg
		return a
	}
	newAlertAt := func(at time.Time) *kkok.Alert {
		a := newAlert()
		a.Date = at
		return a
	}
	newAlertRoutes := func(routes []string) *kkok.Alert {
		a := newAlert()
		a.Routes = routes
		return a
	}
	newAlertInfo := func(info map[string]interface{}) *kkok.Alert {
		a := newAlert()
		a.Info = info
		return a
	}
	newAlertStats := func(stats map[string]float64) *kkok.Alert {
		a := newAlert()
		a.Stats = stats
		return a
	}
	newAlertSub := func(sub []*kkok.Alert) *kkok.Alert {
		a := newAlert()
		a.Sub = sub
		return a
	}

	data := map[string]struct {
		alert  *kkok.Alert
		js     string
		expect *kkok.Alert
	}{
		"none": {
			newAlert(),
			``,
			newAlert(),
		},
		"from": {
			newAlert(),
			`alert.From = "テスト";`,
			newAlertFrom("テスト"),
		},
		"invalid_from": {
			newAlert(),
			`alert.From = true;`,
			nil,
		},
		"empty_from": {
			newAlert(),
			`alert.From = "";`,
			nil,
		},
		"host": {
			newAlert(),
			`alert.Host = "テスト";`,
			newAlertHost("テスト"),
		},
		"invalid_host": {
			newAlert(),
			`alert.Host = true;`,
			nil,
		},
		"title": {
			newAlert(),
			`alert.Title = "テスト";`,
			newAlertTitle("テスト"),
		},
		"invalid_title": {
			newAlert(),
			`alert.Title = true;`,
			nil,
		},
		"empty_title": {
			newAlert(),
			`alert.Title = "";`,
			nil,
		},
		"message": {
			newAlert(),
			`alert.Message = "テスト";`,
			newAlertMessage("テスト"),
		},
		"invalid_message": {
			newAlert(),
			`alert.Message = true;`,
			nil,
		},
		"date_modify": {
			newAlert(),
			`alert.Date.setYear(2017);`,
			newAlertAt(time.Date(2017, 2, 3, 4, 5, 6, 123000000, time.UTC)),
		},
		"date_replace": {
			newAlert(),
			`alert.Date = new Date("2016-12-31T11:22:33Z");`,
			newAlertAt(time.Date(2016, 12, 31, 11, 22, 33, 0, time.UTC)),
		},
		"invalid_date": {
			newAlert(),
			`alert.Date = "2016-12-31T11:22:33Z";`,
			nil,
		},
		"route_add1": {
			newAlert(),
			`alert.Routes.push("r1");`,
			newAlertRoutes([]string{"r1"}),
		},
		"route_add2": {
			newAlertRoutes([]string{"r1"}),
			`alert.Routes.push("r2");`,
			newAlertRoutes([]string{"r1", "r2"}),
		},
		"route_modify": {
			newAlertRoutes([]string{"r1"}),
			`alert.Routes[0] = "r2";`,
			newAlertRoutes([]string{"r2"}),
		},
		"route_delete": {
			newAlertRoutes([]string{"r1", "r2"}),
			`alert.Routes.splice(0, alert.Routes.length);`,
			newAlert(),
		},
		"route_replace": {
			newAlertRoutes([]string{"r1", "r2"}),
			`alert.Routes = ["r3"];`,
			newAlertRoutes([]string{"r3"}),
		},
		"invalid_routes": {
			newAlert(),
			`alert.Routes.push(true);`,
			nil,
		},
		"info_add1": {
			newAlert(),
			`alert.Info.test = 10;`,
			newAlertInfo(map[string]interface{}{"test": 10}),
		},
		"info_add2": {
			newAlertInfo(map[string]interface{}{"test": 10}),
			`alert.Info.test2 = "hoge";`,
			newAlertInfo(map[string]interface{}{"test": 10, "test2": "hoge"}),
		},
		"info_modify": {
			newAlertInfo(map[string]interface{}{"test": 10}),
			`alert.Info.test = "hoge";`,
			newAlertInfo(map[string]interface{}{"test": "hoge"}),
		},
		"info_delete": {
			newAlertInfo(map[string]interface{}{"test": 10}),
			`delete alert.Info.test;`,
			newAlert(),
		},
		"info_replace": {
			newAlertInfo(map[string]interface{}{"test": 10}),
			`alert.Info = {"test2": 3.14};`,
			newAlertInfo(map[string]interface{}{"test2": 3.14}),
		},
		"invalid_info": {
			newAlert(),
			`alert.Info = 3;`,
			nil,
		},
		"stats_add1": {
			newAlert(),
			`alert.Stats.test = 1;`,
			newAlertStats(map[string]float64{"test": 1}),
		},
		"stats_add2": {
			newAlertStats(map[string]float64{"test": 1}),
			`alert.Stats.test2 = 3.14;`,
			newAlertStats(map[string]float64{"test": 1, "test2": 3.14}),
		},
		"stats_modify": {
			newAlertStats(map[string]float64{"test": 3.14}),
			`alert.Stats.test = 1;`,
			newAlertStats(map[string]float64{"test": 1}),
		},
		"stats_delete": {
			newAlertStats(map[string]float64{"test": 1, "test2": 3.14}),
			`delete alert.Stats["test2"];`,
			newAlertStats(map[string]float64{"test": 1}),
		},
		"stats_replace": {
			newAlertStats(map[string]float64{"test": 1}),
			`alert.Stats = {"test2": 3.14};`,
			newAlertStats(map[string]float64{"test2": 3.14}),
		},
		"invalid_stats1": {
			newAlertStats(map[string]float64{"test": 1}),
			`alert.Stats.test = {};`,
			nil,
		},
		"invalid_stats2": {
			newAlert(),
			`alert.Stats = true;`,
			nil,
		},
		"sub": {
			newAlertSub([]*kkok.Alert{{From: "from1"}, {From: "from2"}}),
			`alert.Sub.length == 2;`,
			newAlertSub([]*kkok.Alert{{From: "from1"}, {From: "from2"}}),
		},
	}

	for k, v := range data {
		v := v // to avoid modification at later loop
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			testConvert(t, v.js, v.alert, v.expect)
		})
	}
}
