package exec

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
	"f1": {nil, nil},
	"f2": {map[string]interface{}{"command": []string{}}, nil},
	"f3": {map[string]interface{}{"command": "cat"}, nil},
	"f4": {map[string]interface{}{"command": []string{"cat"}}, &filter{
		command: []string{"cat"},
		timeout: defaultTimeout,
	}},
	"f5": {map[string]interface{}{
		"command": []string{"cat"},
		"timeout": "-1",
	}, nil},
	"f6": {map[string]interface{}{
		"command": []string{"cat"},
		"timeout": -1,
	}, nil},
	"f7": {map[string]interface{}{
		"command": []string{"cat"},
		"timeout": 9,
	}, &filter{
		command: []string{"cat"},
		timeout: 9 * time.Second,
	}},
	"f8": {map[string]interface{}{
		"command": []string{"cat"},
		"timeout": float64(9),
	}, &filter{
		command: []string{"cat"},
		timeout: 9 * time.Second,
	}},
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

	ff, ok := f.(*filter)
	if !ok {
		t.Fatal("not a proper filter")
	}
	if !reflect.DeepEqual(ff, data.filter) {
		t.Error(`!reflect.DeepEqual(ff, data.filter)`, ff, data.filter)
	}
}

func TestCtor(t *testing.T) {
	for id, data := range ctorTestData {
		if data.filter != nil {
			data.filter.Init(id, data.params)
		}
		t.Run(id, func(t *testing.T) { testCtorOne(t, id, data) })
	}
}
