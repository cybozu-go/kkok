package kkok

import (
	"reflect"
	"testing"
	"time"

	"github.com/robertkrimen/otto"
)

func testAlertOtto(t *testing.T) {
	t.Parallel()

	a := &Alert{
		From:    "NTP monitor",
		Date:    time.Now().UTC(),
		Host:    "localhost",
		Title:   "test",
		Message: "test\ntest\n",
		Routes:  []string{"notify"},
		Info:    map[string]interface{}{"domain": "example.org"},
	}

	vm := otto.New()
	err := vm.Set("alert", a)
	if err != nil {
		t.Fatal(err)
	}
	err = vm.Set("alerts", []*Alert{a})
	if err != nil {
		t.Fatal(err)
	}

	script, err := vm.Compile("", "alert.From == 'NTP monitor'")
	if err != nil {
		t.Fatal(err)
	}
	value, err := vm.Run(script)
	if err != nil {
		t.Error(err)
	}
	bvalue, _ := value.ToBoolean()
	if !bvalue {
		t.Error(`!bvalue`)
	}

	script, err = vm.Compile("", "alerts.length == 1")
	if err != nil {
		t.Fatal(err)
	}
	value, err = vm.Run(script)
	if err != nil {
		t.Error(err)
	}
	bvalue, _ = value.ToBoolean()
	if !bvalue {
		t.Error(`!bvalue`)
	}
}

func testAlertClone(t *testing.T) {
	t.Parallel()

	a := &Alert{
		From:    "NTP monitor",
		Date:    time.Now().UTC(),
		Host:    "localhost",
		Title:   "test",
		Message: "test\ntest\n",
		Routes:  []string{"notify"},
		Info:    map[string]interface{}{"domain": "example.org"},
		Sub: []*Alert{
			&Alert{From: "sub1"}, &Alert{From: "sub2"},
		},
	}

	a2 := a.Clone()
	if a2.From != a.From {
		t.Error(`a2.From != a.From`)
	}
	if !a2.Date.Equal(a.Date) {
		t.Error(`!a2.Date.Equal(a.Date)`)
	}
	if a2.Host != a.Host {
		t.Error(`a2.Host != a.Host`)
	}
	if a2.Title != a.Title {
		t.Error(`a2.Title != a.Title`)
	}
	if a2.Message != a.Message {
		t.Error(`a2.Message != a.Message`)
	}
	if !reflect.DeepEqual(a2.Routes, a.Routes) {
		t.Error(`!reflect.DeepEqual(a2.Routes, a.Routes)`)
	}
	if !reflect.DeepEqual(a2.Info, a.Info) {
		t.Error(`!reflect.DeepEqual(a2.Info, a.Info)`)
	}
	if !reflect.DeepEqual(a2.Sub, a.Sub) {
		t.Error(`!reflect.DeepEqual(a2.Sub, a.Sub)`)
	}
}

func TestAlert(t *testing.T) {
	t.Run("Otto", testAlertOtto)
	t.Run("Clone", testAlertClone)
}
