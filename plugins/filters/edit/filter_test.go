package edit

import (
	"reflect"
	"testing"
	"time"

	"github.com/cybozu-go/kkok"
)

func TestFilter(t *testing.T) {
	t.Run("Params", testParams)
	t.Run("Edit", testEdit)
	t.Run("Process", testProcess)
}

func testParams(t *testing.T) {
	t.Parallel()

	f := &filter{origCode: "alert.From = 'foo';"}
	f.Init("id", map[string]interface{}{"label": "label1"})
	pp := f.Params()

	if pp.Type != filterType {
		t.Error(`pp.Type != filterType`)
	}

	if pp.Params["code"] != "alert.From = 'foo';" {
		t.Error(`pp.Params["code"] != "alert.From = 'foo';"`)
	}
	if pp.Params["label"] != "label1" {
		t.Error(`pp.Params["label"] != "label1"`)
	}
}

func testEdit(t *testing.T) {
	t.Run("Success", testEditSuccess)
	t.Run("ScriptErropr", testEditScriptError)
	t.Run("BadValue", testEditBadValue)
}

func testCreateFilter(t *testing.T, code string) *filter {
	f := &filter{}
	s, err := kkok.CompileJS(code)
	if err != nil {
		t.Fatal(err)
	}
	f.code = s
	f.BaseFilter.Init("id", nil)
	return f
}

func testEditSuccess(t *testing.T) {
	t.Parallel()

	f := testCreateFilter(t, `
alert.From = "[foo] " + alert.From;
alert.Date.setYear(2017);
`)
	a := &kkok.Alert{
		From:    "from",
		Date:    time.Date(2011, 2, 3, 4, 5, 6, 123000000, time.UTC),
		Host:    "localhost",
		Title:   "title",
		Message: "msg",
		Info:    map[string]interface{}{"info1": true},
	}
	aa, err := f.edit(a)
	if err != nil {
		t.Fatal(err)
	}
	a.From = "[foo] from"
	a.Date = time.Date(2017, 2, 3, 4, 5, 6, 123000000, time.UTC)
	aa.Stats = nil

	if !reflect.DeepEqual(aa, a) {
		t.Error(`!reflect.DeepEqual(aa, a)`)
		t.Logf("%#v", aa)
	}
}

func testEditScriptError(t *testing.T) {
	t.Parallel()

	f := testCreateFilter(t, `a.hoge = 3;`)
	_, err := f.edit(&kkok.Alert{})
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testEditBadValue(t *testing.T) {
	t.Parallel()

	f := testCreateFilter(t, `alert.Date = 3;`)
	a := &kkok.Alert{
		From:    "from",
		Date:    time.Date(2011, 2, 3, 4, 5, 6, 123000000, time.UTC),
		Host:    "localhost",
		Title:   "title",
		Message: "msg",
		Info:    map[string]interface{}{"info1": true},
	}
	_, err := f.edit(a)
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testProcess(t *testing.T) {
	f := &filter{}
	s, err := kkok.CompileJS(`alert.Title = "a";`)
	if err != nil {
		t.Fatal(err)
	}
	f.code = s
	f.BaseFilter.Init("id", map[string]interface{}{
		"if": `alert.From == "foo"`,
	})

	alerts := []*kkok.Alert{
		{From: "foo", Title: "title"},
		{From: "bar", Title: "title"},
		{From: "foo", Title: "title"},
		{From: "zot", Title: "title"},
	}

	alerts, err = f.Process(alerts)
	if err != nil {
		t.Fatal(err)
	}

	if len(alerts) != 4 {
		t.Error(`len(alerts) != 4`)
	}

	for _, a := range alerts {
		if a.From == "foo" && a.Title != "a" {
			t.Error(`a.From == "foo" && a.Title != "a"`)
			t.Logf("%#v", a)
		}
		if a.From != "foo" && a.Title != "title" {
			t.Error(`a.From != "foo" && a.Title != "title"`)
			t.Logf("%#v", a)
		}
	}
}
