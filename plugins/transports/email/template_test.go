package email

import (
	"bytes"
	"testing"
	"time"

	"github.com/cybozu-go/kkok"
)

type templateTest struct {
	alert  *kkok.Alert
	result string
}

var (
	testDate = time.Date(2016, 12, 24, 1, 2, 34, 56789, time.UTC)
	testTZ   = time.FixedZone("-1", -3600)

	templateTestData = map[string]templateTest{
		"nomsg": {&kkok.Alert{
			From:  "from1",
			Date:  testDate,
			Host:  "localhost",
			Title: "title1",
		}, `From: from1
Date: 2016-12-24T01:02:34.000056789Z
Host: localhost
Title: title1
`},
		"msgnonl": {&kkok.Alert{
			From:    "from1",
			Date:    testDate,
			Host:    "localhost",
			Title:   "title1",
			Message: "foo\nbar",
		}, `From: from1
Date: 2016-12-24T01:02:34.000056789Z
Host: localhost
Title: title1

foo
bar`},
		"msg": {&kkok.Alert{
			From:    "from1",
			Date:    testDate,
			Host:    "localhost",
			Title:   "title1",
			Message: "foo\nbar\n",
		}, `From: from1
Date: 2016-12-24T01:02:34.000056789Z
Host: localhost
Title: title1

foo
bar
`},
		"info": {&kkok.Alert{
			From:  "from1",
			Date:  testDate,
			Host:  "localhost",
			Title: "title1",
			Info:  map[string]interface{}{"foo": 1, "bar": true},
		}, `From: from1
Date: 2016-12-24T01:02:34.000056789Z
Host: localhost
Title: title1

Info:
  bar=true
  foo=1
`},
		"msginfo": {&kkok.Alert{
			From:    "from1",
			Date:    testDate,
			Host:    "localhost",
			Title:   "title1",
			Message: "foo\nbar\n",
			Info:    map[string]interface{}{"foo": 1, "bar": true},
		}, `From: from1
Date: 2016-12-24T01:02:34.000056789Z
Host: localhost
Title: title1

foo
bar

Info:
  bar=true
  foo=1
`},
		"infosub": {&kkok.Alert{
			From:  "from1",
			Date:  testDate,
			Host:  "localhost",
			Title: "title1",
			Info:  map[string]interface{}{"foo": 1, "bar": true},
			Sub: []*kkok.Alert{
				{
					From:  "from2",
					Date:  time.Date(2011, 2, 3, 4, 5, 6, 0, time.UTC),
					Host:  "host2",
					Title: "title long",
				},
			},
		}, `From: from1
Date: 2016-12-24T01:02:34.000056789Z
Host: localhost
Title: title1

Info:
  bar=true
  foo=1

Sub alerts:
-------------------------------------------------------
From: from2
Date: 2011-02-03T04:05:06Z
Host: host2
Title: title long
-------------------------------------------------------
`},
		"sub2": {&kkok.Alert{
			From:  "from1",
			Date:  testDate,
			Host:  "localhost",
			Title: "title1",
			Info:  map[string]interface{}{"foo": 1, "bar": true},
			Sub: []*kkok.Alert{
				{
					From:  "from2",
					Date:  time.Date(2011, 2, 3, 4, 5, 6, 0, time.UTC),
					Host:  "host2",
					Title: "title2",
				},
				{
					From:    "from3",
					Date:    time.Date(2011, 2, 3, 4, 5, 6, 0, testTZ),
					Host:    "host3",
					Title:   "title3",
					Message: "hello\nworld!\n",
				},
			},
		}, `From: from1
Date: 2016-12-24T01:02:34.000056789Z
Host: localhost
Title: title1

Info:
  bar=true
  foo=1

Sub alerts:
-------------------------------------------------------
From: from2
Date: 2011-02-03T04:05:06Z
Host: host2
Title: title2
-------------------------------------------------------
From: from3
Date: 2011-02-03T05:05:06Z
Host: host3
Title: title3

hello
world!
-------------------------------------------------------
`},
	}
)

func testTemplateOne(t *testing.T, alert *kkok.Alert, expected string) {
	t.Parallel()

	buf := new(bytes.Buffer)
	err := defaultTemplate.Execute(buf, alert)
	if err != nil {
		t.Fatal(err)
	}
	result := buf.String()
	if result != expected {
		t.Error(`result != expected`)
		t.Logf("%#v", result)
		t.Logf("%#v", expected)
	}
}

func TestTemplate(t *testing.T) {
	for key, data := range templateTestData {
		t.Run(key, func(t *testing.T) {
			testTemplateOne(t, data.alert, data.result)
		})
	}
}
