package kkok

import (
	"reflect"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
)

func testConfigDefault(t *testing.T) {
	t.Parallel()

	c := NewConfig()
	if c.InitialInterval != defaultInitialInterval {
		t.Error(`c.InitialInterval != defaultInitialInterval`)
	}
	if c.MaxInterval != defaultMaxInterval {
		t.Error(`c.MaxInterval != defaultMaxInterval`)
	}
	if c.Addr != defaultAddr {
		t.Error(`c.Addr != defaultAddr`)
	}
}

func testConfigLoad(t *testing.T) {
	t.Parallel()

	c := NewConfig()
	md, err := toml.DecodeFile("testdata/config.toml", c)
	if err != nil {
		t.Fatal(err)
	}
	if len(md.Undecoded()) > 0 {
		// As md reports values decoded with toml.Unmarshaler as undecoded,
		// this is just not an error.
		t.Log("undecoded:", md.Undecoded())
	}
	if c.InitialInterval != 100 {
		t.Error(`c.InitialInterval != 100`)
	}
	if c.InitialDuration() != 100*time.Second {
		t.Error(`c.InitialDuration() != 100 * time.Second`)
	}
	if c.MaxInterval != 1000 {
		t.Error(`c.MaxInterval != 1000`)
	}
	if c.MaxDuration() != 1000*time.Second {
		t.Error(`c.MaxDuration() != 1000 * time.Second`)
	}
	if c.Addr != "localhost:12229" {
		t.Error(`c.Addr != "localhost:12229"`)
	}
	if len(c.APIToken) != 0 {
		t.Error(`len(c.APIToken) != 0`)
	}

	// log
	if c.Log.Level != "debug" {
		t.Error(`c.Log.Level != "debug"`)
	}

	// sources
	if len(c.Sources) != 1 {
		t.Fatal(`len(c.Sources) != 1`)
	}
	s := c.Sources[0]
	if s.Type != "maildir" {
		t.Error(`s.Type != "maildir"`)
	}
	if s.Params["dir"] != "/var/mail/ymmt2005" {
		t.Error(`s.Params["dir"] != "/var/mail/ymmt2005"`)
	}

	// routes
	if len(c.Routes) != 2 {
		t.Fatal(`len(c.Routes) != 2`)
	}
	rnotify, ok := c.Routes["notify"]
	if !ok {
		t.Fatal(`no "notify" route`)
	}
	if len(rnotify) != 2 {
		t.Fatal(`len(rnotify) != 2`)
	}
	if rnotify[0].Type != "slack" {
		t.Error(`rnotify[0].Type != "slack"`)
	}
	if rnotify[0].Params["url"] != "https://hooks.slack.com/xxxx" {
		t.Error(`rnotify[0].Params["url"] != "https://hooks.slack.com/xxxx"`)
	}
	if rnotify[1].Type != "email" {
		t.Error(`rnotify[1].Type != "email"`)
	}
	if !reflect.DeepEqual(rnotify[1].Params["to"], []interface{}{
		"ymmt2005@example.com", "ymmt@example.org",
	}) {
		t.Error(`rnotify[1].Params["to"] != ["ymmt2005@example.com", "ymmt@example.org"]`)
	}

	remerg, ok := c.Routes["emergency"]
	if !ok {
		t.Fatal(`no "emergency" route`)
	}
	if len(remerg) != 1 {
		t.Fatal(`len(remerg) != 1`)
	}
	if remerg[0].Type != "twilio" {
		t.Error(`remerg.Type != "twilio"`)
	}
	if remerg[0].Params["token"] != "yyyy" {
		t.Error(`remerg.Params["token"] != "yyyy"`)
	}

	// filters
	if len(c.Filters) != 2 {
		t.Fatal(`len(c.Filters) != 2`)
	}
	f1 := c.Filters[0]
	if f1.Type != "route" {
		t.Error(`f1.Type != "route"`)
	}
	if f1.Params["id"] != "default" {
		t.Error(`f1.Params["id"] != "default"`)
	}
	if !reflect.DeepEqual(f1.Params["routes"], []interface{}{"notify"}) {
		t.Error(`f1.Params["routes"] != ["notify"]`)
	}
}

func TestConfig(t *testing.T) {
	t.Run("Default", testConfigDefault)
	t.Run("Load", testConfigLoad)
}
