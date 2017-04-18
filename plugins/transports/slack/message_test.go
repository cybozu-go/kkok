package slack

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/cybozu-go/kkok"
)

func TestMessage(t *testing.T) {
	t.Parallel()

	a1 := &attachment{
		Fallback: "fallback1",
		Color:    "#D00000",
		Title:    "title1",
		Text:     "text1",
	}

	a1.addField("short string", "short <string>")
	a1.addField("long string", "long\nstring")
	a1.addField("bool", true)
	a1.addField("int", 3)
	a1.addField("float64", 3.14159)
	a1.addField("struct", &kkok.Alert{Host: "localhost", Title: "hoge"})

	a2 := &attachment{
		Fallback: "fallback2",
	}

	m := &message{
		Text:        "text",
		Name:        "name",
		Icon:        ":sushi:",
		Channel:     "#test",
		Attachments: []*attachment{a1, a2},
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	var m2 map[string]interface{}
	err = json.Unmarshal(data, &m2)
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := m2["text"]; !ok {
		t.Error(`v, ok := m2["text"]; !ok`)
	} else if v != "text" {
		t.Error(`v != "text"`)
	}

	if v, ok := m2["username"]; !ok {
		t.Error(`v, ok := m2["username"]; !ok`)
	} else if v != "name" {
		t.Error(`v != "name"`)
	}

	if v, ok := m2["icon_emoji"]; !ok {
		t.Error(`v, ok := m2["icon_emoji"]; !ok`)
	} else if v != ":sushi:" {
		t.Error(`v != ":sushi:"`)
	}

	if v, ok := m2["channel"]; !ok {
		t.Error(`v, ok := m2["channel"]; !ok`)
	} else if v != "#test" {
		t.Error(`v != "#test"`)
	}

	if v, ok := m2["attachments"]; !ok {
		t.Error(`v, ok := m2["attachments"]; !ok`)
	} else if attachments, ok := v.([]interface{}); !ok {
		t.Error(`attachments, ok := v.([]interface{}); !ok`)
	} else if len(attachments) != 2 {
		t.Error(`len(attachments) != 2`)
	} else {
		if _, ok := attachments[0].(map[string]interface{}); !ok {
			t.Error(`ra1, ok := attachments[0].(map[string]interface{})`)
		}
		if _, ok := attachments[0].(map[string]interface{}); !ok {
			t.Error(`ra2, ok := attachments[0].(map[string]interface{})`)
		}
	}

	m3 := &message{}
	data, err = json.Marshal(m3)
	if err != nil {
		t.Fatal(err)
	}
	var m4 map[string]interface{}
	err = json.Unmarshal(data, &m4)
	if err != nil {
		t.Fatal(err)
	}
	if len(m4) != 0 {
		t.Error(`len(m4) != 0`)
	}
}

func TestAttachment(t *testing.T) {
	t.Parallel()

	a1 := &attachment{
		Fallback: "fallback1",
		Color:    "#D00000",
		Title:    "title1",
		Text:     "text1",
	}

	a1.addField("short string", "short <string>")
	a1.addField("long string", "long\nstring")
	a1.addField("bool", true)
	a1.addField("int", 3)
	a1.addField("float64", 3.14159)
	a1.addField("struct", &kkok.Alert{Host: "localhost", Title: "hoge"})

	data, err := json.Marshal(a1)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := m["fallback"]; !ok {
		t.Error(`v, ok := m["fallback"]; !ok`)
	} else if v != "fallback1" {
		t.Error(`v != "fallback1"`)
	}

	if v, ok := m["color"]; !ok {
		t.Error(`v, ok := m["color"]; !ok`)
	} else if v != "#D00000" {
		t.Error(`v != "#D00000"`)
	}

	if v, ok := m["title"]; !ok {
		t.Error(`v, ok := m["title"]; !ok`)
	} else if v != "title1" {
		t.Error(`v != "title1"`)
	}

	if v, ok := m["text"]; !ok {
		t.Error(`v, ok := m["text"]; !ok`)
	} else if v != "text1" {
		t.Error(`v != "text1"`)
	}

	v, ok := m["fields"]
	if !ok {
		t.Fatal("no fields")
	}
	fields, ok := v.([]interface{})
	if !ok {
		t.Fatal("fields is not a slice")
	}
	if len(fields) != 6 {
		t.Fatal(`len(fields) != 6`)
	}
	f5, ok := fields[5].(map[string]interface{})
	if !ok {
		t.Fatal("fields[5] is not an object")
	}
	delete(f5, "value")

	ans := []map[string]interface{}{
		{
			"title": "short string",
			"value": "short &lt;string&gt;",
			"short": true,
		},
		{
			"title": "long string",
			"value": "long\nstring",
			"short": false,
		},
		{
			"title": "bool",
			"value": "true",
			"short": true,
		},
		{
			"title": "int",
			"value": "3",
			"short": true,
		},
		{
			"title": "float64",
			"value": "3.14159",
			"short": true,
		},
		{
			"title": "struct",
			"short": false,
		},
	}
	for i, f := range fields {
		if !reflect.DeepEqual(f, ans[i]) {
			t.Error(`!reflect.DeepEqual(f, ans[i]); i=`, i)
		}
	}

	a2 := &attachment{
		Fallback: "fallback2",
	}
	data, err = json.Marshal(a2)
	if err != nil {
		t.Fatal(err)
	}
	var m2 map[string]interface{}
	err = json.Unmarshal(data, &m2)
	if err != nil {
		t.Fatal(err)
	}
	if len(m2) != 1 {
		t.Error(`len(m2) != 1`)
	}
	if _, ok := m2["fallback"]; !ok {
		t.Error(`_, ok := m2["fallback"]; !ok`)
	}
}
