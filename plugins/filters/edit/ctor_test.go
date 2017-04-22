package edit

import (
	"reflect"
	"testing"
)

type ctorTest struct {
	params map[string]interface{}
	filter *filter
}

var ctorTestData = map[string]ctorTest{
	"f1": {nil, nil},
	"f2": {map[string]interface{}{
		"code": 1,
	}, nil},
	"f3": {map[string]interface{}{
		"code": `alert.From = "a";`,
	}, &filter{
		origCode: `alert.From = "a";`,
	}},
	"f4": {map[string]interface{}{
		"label": "label",
		"code":  `alert.From = "a";`,
	}, &filter{
		origCode: `alert.From = "a";`,
	}},
	"f5": {map[string]interface{}{
		"all":  true,
		"code": `alert.From = "a";`,
	}, nil},
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
	ff.code = nil
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
