package freq

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
	"f1": {nil, &filter{
		duration: defaultDuration,
		divisor:  defaultDivisor,
	}},
	"f2": {map[string]interface{}{"duration": 3.0}, &filter{
		duration: 3 * time.Second,
		divisor:  defaultDivisor,
	}},
	"f3": {map[string]interface{}{"duration": -1}, nil},
	"f4": {map[string]interface{}{"duration": "1"}, nil},
	"f5": {map[string]interface{}{"divisor": 5}, &filter{
		duration: defaultDuration,
		divisor:  5,
	}},
	"f6": {map[string]interface{}{"divisor": 0}, nil},
	"f7": {map[string]interface{}{"divisor": "5"}, nil},
	"f8": {map[string]interface{}{"classify": "From"}, &filter{
		duration: defaultDuration,
		divisor:  defaultDivisor,
		cl:       clFrom,
	}},
	"f9": {map[string]interface{}{"classify": "from"}, &filter{
		duration: defaultDuration,
		divisor:  defaultDivisor,
		cl:       clFrom,
	}},
	"f10": {map[string]interface{}{"classify": "Title"}, &filter{
		duration: defaultDuration,
		divisor:  defaultDivisor,
		cl:       clTitle,
	}},
	"f11": {map[string]interface{}{"classify": "title"}, &filter{
		duration: defaultDuration,
		divisor:  defaultDivisor,
		cl:       clTitle,
	}},
	"f12": {map[string]interface{}{"classify": "Host"}, &filter{
		duration: defaultDuration,
		divisor:  defaultDivisor,
		cl:       clHost,
	}},
	"f13": {map[string]interface{}{"classify": "host"}, &filter{
		duration: defaultDuration,
		divisor:  defaultDivisor,
		cl:       clHost,
	}},
	"f14": {map[string]interface{}{"classify": "foo"}, nil},
	"f15": {map[string]interface{}{"key": "test"}, &filter{
		duration: defaultDuration,
		divisor:  defaultDivisor,
		key:      "test",
	}},
	"f16": {map[string]interface{}{"key": true}, nil},
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
	ff.samples = nil
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
