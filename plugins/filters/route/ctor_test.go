package route

import (
	"reflect"
	"testing"
	"time"
)

type ctorTest struct {
	params map[string]interface{}
	filter *filter
}

var ctorTestData = map[string]ctorTest{
	"f1": {nil, &filter{}},
	"f2": {map[string]interface{}{"routes": 1}, nil},
	"f3": {
		map[string]interface{}{"routes": []string{"r1", "r2"}},
		&filter{
			routes: []string{"r1", "r2"},
		},
	},
	"f4": {
		map[string]interface{}{"mute_routes": []string{"r1", "r2"}},
		&filter{
			muteRoutes: []string{"r1", "r2"},
		},
	},
	"f5":  {map[string]interface{}{"replace": nil}, &filter{}},
	"f6":  {map[string]interface{}{"replace": "true"}, nil},
	"f7":  {map[string]interface{}{"replace": true}, &filter{replace: true}},
	"f8":  {map[string]interface{}{"auto_mute": nil}, &filter{}},
	"f9":  {map[string]interface{}{"auto_mute": "true"}, nil},
	"f10": {map[string]interface{}{"auto_mute": true}, &filter{autoMute: true}},
	"f11": {map[string]interface{}{"mute_seconds": 3.1}, nil},
	"f12": {
		map[string]interface{}{"mute_seconds": 3.0},
		&filter{muteDuration: 3 * time.Second},
	},
}

func testCtorOne(t *testing.T, id string, data ctorTest) {
	t.Parallel()

	f, err := ctor(id, data.params)
	if err != nil {
		if data.filter != nil {
			// unexpected
			t.Error(err)
		}
		return
	}

	if !reflect.DeepEqual(f, data.filter) {
		t.Error(`!reflect.DeepEqual(f, data.filter)`, f, data.filter)
	}
}

func TestCtor(t *testing.T) {
	for id, data := range ctorTestData {
		if data.filter != nil {
			data.filter.Init(id, data.params)
			if data.filter.muteDuration == 0 {
				data.filter.muteDuration = time.Second * defaultMuteSeconds
			}
		}
		t.Run(id, func(t *testing.T) { testCtorOne(t, id, data) })
	}
}
