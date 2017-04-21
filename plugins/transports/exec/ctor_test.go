package exec

import (
	"reflect"
	"testing"
	"time"
)

type ctorTest struct {
	params map[string]interface{}
	tr     *transport
}

var ctorTestData = map[string]ctorTest{
	"tr1": {nil, nil},
	"tr2": {map[string]interface{}{
		"command": "abc",
	}, nil},
	"tr3": {map[string]interface{}{
		"command": []string{},
	}, nil},
	"tr4": {map[string]interface{}{
		"command": []string{"curl"},
	}, &transport{
		command: "curl",
		timeout: defaultTimeout,
	}},
	"tr5": {map[string]interface{}{
		"command": []string{"curl", "-d", "@-", "http://example.com"},
	}, &transport{
		command: "curl",
		args:    []string{"-d", "@-", "http://example.com"},
		timeout: defaultTimeout,
	}},
	"tr6": {map[string]interface{}{
		"label":   "abc",
		"command": []string{"curl"},
	}, &transport{
		label:   "abc",
		command: "curl",
		timeout: defaultTimeout,
	}},
	"tr7": {map[string]interface{}{
		"command": []string{"curl"},
		"all":     true,
	}, &transport{
		command: "curl",
		timeout: defaultTimeout,
		all:     true,
	}},
	"tr8": {map[string]interface{}{
		"command": []string{"curl"},
		"timeout": -1,
	}, nil},
	"tr9": {map[string]interface{}{
		"command": []string{"curl"},
		"timeout": 0,
	}, &transport{
		command: "curl",
		timeout: 0,
	}},
	"tr10": {map[string]interface{}{
		"command": []string{"curl"},
		"timeout": 9,
	}, &transport{
		command: "curl",
		timeout: 9 * time.Second,
	}},
}

func testCtorOne(t *testing.T, data ctorTest) {
	t.Parallel()

	tr, err := ctor(data.params)
	if err != nil {
		if data.tr != nil {
			// unexpected
			t.Error(err)
		}
		return
	}

	tr2, ok := tr.(*transport)
	if !ok {
		t.Fatal("not a proper transport")
	}
	if len(tr2.args) == 0 {
		tr2.args = nil
	}
	if !reflect.DeepEqual(tr2, data.tr) {
		t.Error(`!reflect.DeepEqual(tr2, data.tr)`)
		t.Logf("%#v, %#v", tr2, data.tr)
	}
}

func TestCtor(t *testing.T) {
	for id, data := range ctorTestData {
		t.Run(id, func(t *testing.T) { testCtorOne(t, data) })
	}
}
