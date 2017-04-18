package slack

import (
	"reflect"
	"testing"
)

type ctorTest struct {
	params map[string]interface{}
	tr     *transport
}

var ctorTestData = map[string]ctorTest{
	"tr1": {nil, nil},
	"tr2": {map[string]interface{}{
		"url": "https://slack.com/foo/bar",
	}, &transport{
		maxRetry: defaultRetry,
	}},
	"tr3": {map[string]interface{}{"url": 3.14}, nil},
	"tr4": {map[string]interface{}{"url": "hoge"}, nil},
	"tr5": {map[string]interface{}{
		"url":   "https://slack.com/foo/bar",
		"label": "label1",
	}, &transport{
		label:    "label1",
		maxRetry: defaultRetry,
	}},
	"tr6": {map[string]interface{}{
		"url":       "https://slack.com/foo/bar",
		"max_retry": 0,
	}, &transport{
		maxRetry: 0,
	}},
	"tr7": {map[string]interface{}{
		"url":       "https://slack.com/foo/bar",
		"max_retry": 100,
	}, &transport{
		maxRetry: 100,
	}},
	"tr8": {map[string]interface{}{
		"url":       "https://slack.com/foo/bar",
		"max_retry": "100",
	}, nil},
	"tr9": {map[string]interface{}{
		"url":  "https://slack.com/foo/bar",
		"name": "test",
	}, &transport{
		maxRetry: defaultRetry,
		name:     "test",
	}},
	"tr10": {map[string]interface{}{
		"url":  "https://slack.com/foo/bar",
		"icon": ":sushi:",
	}, &transport{
		maxRetry: defaultRetry,
		icon:     ":sushi:",
	}},
	"tr11": {map[string]interface{}{
		"url":     "https://slack.com/foo/bar",
		"channel": "#random",
	}, &transport{
		maxRetry: defaultRetry,
		channel:  "#random",
	}},
	"tr12": {map[string]interface{}{
		"url":   "https://slack.com/foo/bar",
		"color": "'",
	}, nil},
	"tr13": {map[string]interface{}{
		"url":   "https://slack.com/foo/bar",
		"color": "alert.Info.severity",
	}, &transport{
		maxRetry:  defaultRetry,
		origColor: "alert.Info.severity",
	}},
	"tr14": {map[string]interface{}{
		"url":      "https://slack.com/foo/bar",
		"template": "/template/not/exist",
	}, nil},
	"tr15": {map[string]interface{}{
		"url":      "https://slack.com/foo/bar",
		"template": "testdata/1.txt",
	}, &transport{
		maxRetry: defaultRetry,
		tmplPath: "testdata/1.txt",
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
	if tr2.url == nil {
		t.Error(`tr2.url == nil`)
	} else {
		tr2.url = nil
	}
	tr2.color = nil
	tr2.tmpl = nil
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
