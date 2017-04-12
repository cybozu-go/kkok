package email

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
		"from": "test@example.com",
	}, &transport{
		from: "test@example.com",
	}},
	"tr3": {map[string]interface{}{"from": 3.14}, nil},
	"tr4": {map[string]interface{}{
		"from":  "test@example.com",
		"label": "label1",
	}, &transport{
		from:  "test@example.com",
		label: "label1",
	}},
	"tr5": {map[string]interface{}{
		"from": "test@example.com",
		"host": "smtp.localdomain",
	}, &transport{
		from: "test@example.com",
		host: "smtp.localdomain",
	}},
	"tr6": {map[string]interface{}{
		"from": "test@example.com",
		"port": 1025,
	}, &transport{
		from: "test@example.com",
		port: 1025,
	}},
	"tr7": {map[string]interface{}{
		"from": "test@example.com",
		"user": "user1",
	}, &transport{
		from:     "test@example.com",
		username: "user1",
	}},
	"tr8": {map[string]interface{}{
		"from":     "test@example.com",
		"password": "secret",
	}, &transport{
		from:     "test@example.com",
		password: "secret",
	}},
	"tr9": {map[string]interface{}{
		"from": "test@example.com",
		"to":   "hoge@example.org",
	}, nil},
	"tr10": {map[string]interface{}{
		"from": "test@example.com",
		"to":   []string{"hoge@example.org"},
	}, &transport{
		from: "test@example.com",
		to:   []string{"hoge@example.org"},
	}},
	"tr11": {map[string]interface{}{
		"from": "test@example.com",
		"cc":   []string{"hoge@example.org"},
	}, &transport{
		from: "test@example.com",
		cc:   []string{"hoge@example.org"},
	}},
	"tr12": {map[string]interface{}{
		"from": "test@example.com",
		"bcc":  []string{"hoge@example.org"},
	}, &transport{
		from: "test@example.com",
		bcc:  []string{"hoge@example.org"},
	}},
	"tr13": {map[string]interface{}{
		"from":    "test@example.com",
		"to_file": "/path/to/file",
	}, &transport{
		from:   "test@example.com",
		toFile: "/path/to/file",
	}},
	"tr14": {map[string]interface{}{
		"from":    "test@example.com",
		"cc_file": "/path/to/file",
	}, &transport{
		from:   "test@example.com",
		ccFile: "/path/to/file",
	}},
	"tr15": {map[string]interface{}{
		"from":     "test@example.com",
		"bcc_file": "/path/to/file",
	}, &transport{
		from:    "test@example.com",
		bccFile: "/path/to/file",
	}},
	"tr16": {map[string]interface{}{
		"from":     "test@example.com",
		"template": "/template/not/exist",
	}, nil},
	"tr17": {map[string]interface{}{
		"from":     "test@example.com",
		"template": "testdata/1.txt",
	}, &transport{
		from:     "test@example.com",
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
	tr2.tmpl = nil
	if !reflect.DeepEqual(tr2, data.tr) {
		t.Error(`!reflect.DeepEqual(tr2, data.tr)`, tr2, data.tr)
	}
}

func TestCtor(t *testing.T) {
	for id, data := range ctorTestData {
		t.Run(id, func(t *testing.T) { testCtorOne(t, data) })
	}
}
