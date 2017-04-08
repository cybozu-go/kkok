package group

import (
	"reflect"
	"testing"
)

type ctorTest struct {
	params map[string]interface{}
	filter *filter
}

var ctorTestData = map[string]ctorTest{
	"f1": {nil, &filter{}},
	"f2": {map[string]interface{}{"by": "alert.Host"}, &filter{
		origBy: "alert.Host",
	}},
	"f3": {map[string]interface{}{"by": "}"}, nil},
	"f4": {map[string]interface{}{"from": "from1"}, &filter{
		from: "from1",
	}},
	"f5": {map[string]interface{}{"title": "title1"}, &filter{
		title: "title1",
	}},
	"f6": {map[string]interface{}{"message": "msg"}, &filter{
		message: "msg",
	}},
	"f7": {map[string]interface{}{"routes": []string{"r1", "r2"}}, &filter{
		routes: []string{"r1", "r2"},
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
	ff.by = nil
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
