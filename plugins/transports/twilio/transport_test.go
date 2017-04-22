package twilio

import (
	"context"
	"net/url"
	"os"
	"reflect"
	"testing"
	"text/template"
	"time"

	"github.com/cybozu-go/kkok"
)

func TestTransport(t *testing.T) {
	t.Run("String", testString)
	t.Run("Params", testParams)
	t.Run("Deliver", testDeliver)
}

func testString(t *testing.T) {
	t.Parallel()

	tr, err := newTransport("abc")
	if err != nil {
		t.Fatal(err)
	}
	if tr.String() != transportType {
		t.Error(`tr.String() != transportType`)
	}

	tr.label = "label"
	if tr.String() != "label" {
		t.Error(`tr.String() != "label"`)
	}
}

func testParams(t *testing.T) {
	t.Parallel()

	tr, err := newTransport("abc")
	if err != nil {
		t.Fatal(err)
	}
	tr.token = "token"
	tr.from = "+123456789"

	pp := tr.Params()
	if pp.Type != transportType {
		t.Error(`pp.Type != transportType`)
	}
	if !reflect.DeepEqual(pp.Params, map[string]interface{}{
		"account":    "abc",
		"token":      "token",
		"from":       "+123456789",
		"max_length": defaultLength,
		"max_retry":  defaultRetry,
	}) {
		t.Error(`pp.Params is not expected`, pp.Params)
	}

	tr.label = "label"
	tr.sid = "key1"
	tr.to = []string{"+987654321"}
	tr.toFile = "/path/to/file"
	tr.maxLength = 1
	tr.maxRetry = 0
	tr.countOnly = true
	tr.tmplPath = "/path/to/template"
	pp = tr.Params()
	if pp.Type != transportType {
		t.Error(`pp.Type != transportType`)
	}
	if !reflect.DeepEqual(pp.Params, map[string]interface{}{
		"label":      "label",
		"account":    "abc",
		"key_sid":    "key1",
		"token":      "token",
		"from":       "+123456789",
		"to":         []string{"+987654321"},
		"to_file":    "/path/to/file",
		"max_length": 1,
		"max_retry":  0,
		"template":   "/path/to/template",
		"count_only": true,
	}) {
		t.Error(`pp.Params is not expected`)
		t.Logf("%#v", pp.Params)
	}
}

func testDeliver(t *testing.T) {
	t.Run("SMS", testDeliverSMS)
	t.Run("Queue", testDeliverQueue)
	t.Run("Template", testDeliverTemplate)
	t.Run("CountOnly", testDeliverCountOnly)
	t.Run("ToFile", testDeliverToFile)
	t.Run("MaxLength", testDeliverMaxLength)
}

func testDequeue(ch <-chan *twilioSMS) []*twilioSMS {
	var r []*twilioSMS
	for {
		select {
		case m := <-ch:
			r = append(r, m)
		default:
			return r
		}
	}
}

func testDeliverSMS(t *testing.T) {
	t.Parallel()

	tr, err := newTransport("abc")
	if err != nil {
		t.Fatal(err)
	}
	tr.token = "token"
	tr.from = "+123456789"

	ch := make(chan *twilioSMS, 10)
	tr.enqueue = func(m *twilioSMS) bool {
		ch <- m
		return true
	}

	a := &kkok.Alert{
		From:    "from1",
		Title:   "title1",
		Host:    "host1",
		Message: "msg",
	}
	err = tr.Deliver([]*kkok.Alert{a})
	if err != nil {
		t.Fatal(err)
	}
	msgs := testDequeue(ch)
	if len(msgs) != 0 {
		t.Error(`len(msgs) != 0`)
	}

	tr.to = []string{"+123456789", "+100000000"}
	err = tr.Deliver([]*kkok.Alert{a})
	if err != nil {
		t.Fatal(err)
	}
	msgs = testDequeue(ch)
	if len(msgs) != 2 {
		t.Error(`len(msgs) != 2`)
	}

	for i, m := range msgs {
		if m.url.String() != twilioEndPoint+"abc/Messages.json" {
			t.Error(`m.url.String() != twilioEndPoint+"abc/Messages.json"`)
		}
		if m.username != tr.sid {
			t.Error(`m.username != tr.sid`)
		}
		if m.password != tr.token {
			t.Error(`m.password != tr.token`)
		}
		if m.maxRetry != tr.maxRetry {
			t.Error(`m.maxRetry != tr.maxRetry`)
		}

		v, err := url.ParseQuery(m.payload)
		if err != nil {
			t.Fatal(err)
		}
		if v.Get("To") != tr.to[i] {
			t.Error(`v.Get("To") != tr.to[i]; i=`, i)
		}
		if v.Get("From") != tr.from {
			t.Error(`v.Get("From") != tr.from`)
		}
		t.Log(v.Get("Body"))
	}
}

func testDeliverQueue(t *testing.T) {
	t.Parallel()

	tr, err := newTransport("abc")
	if err != nil {
		t.Fatal(err)
	}
	tr.token = "token"
	tr.from = "+123456789"
	tr.to = []string{"+100000000"}

	tr.enqueue = func(m *twilioSMS) bool {
		return false
	}

	err = tr.Deliver([]*kkok.Alert{{}})
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testDeliverTemplate(t *testing.T) {
	t.Parallel()

	tr, err := newTransport("abc")
	if err != nil {
		t.Fatal(err)
	}
	tr.token = "token"
	tr.from = "+123456789"
	tr.to = []string{"+100000000"}

	tmpl, err := template.New("").Parse(`{{.Message}}`)
	if err != nil {
		t.Fatal(err)
	}
	tr.tmpl = tmpl

	ch := make(chan *twilioSMS, 10)
	tr.enqueue = func(m *twilioSMS) bool {
		ch <- m
		return true
	}

	msgs := []string{"msg1", "msg2"}
	err = tr.Deliver([]*kkok.Alert{{Message: msgs[0]}, {Message: msgs[1]}})
	if err != nil {
		t.Fatal(err)
	}
	ms := testDequeue(ch)
	if len(ms) != 2 {
		t.Error(`len(ms) != 2`)
	}

	for i, m := range ms {
		v, err := url.ParseQuery(m.payload)
		if err != nil {
			t.Fatal(err)
		}
		if v.Get("Body") != msgs[i] {
			t.Error(`v.Get("Body") != msgs[i]; i=`, i)
		}
	}
}

func testDeliverCountOnly(t *testing.T) {
	t.Parallel()

	tr, err := newTransport("abc")
	if err != nil {
		t.Fatal(err)
	}
	tr.token = "token"
	tr.from = "+123456789"
	tr.to = []string{"+100000000"}
	tr.countOnly = true

	ch := make(chan *twilioSMS, 10)
	tr.enqueue = func(m *twilioSMS) bool {
		ch <- m
		return true
	}

	err = tr.Deliver([]*kkok.Alert{{Message: "aaa"}, {Message: "bbb"}})
	if err != nil {
		t.Fatal(err)
	}
	ms := testDequeue(ch)
	if len(ms) != 1 {
		t.Fatal(`len(ms) != 1`)
	}

	v, _ := url.ParseQuery(ms[0].payload)
	t.Log(v.Get("Body"))
}

func testDeliverToFile(t *testing.T) {
	t.Parallel()

	tr, err := newTransport("abc")
	if err != nil {
		t.Fatal(err)
	}
	tr.token = "token"
	tr.from = "+123456789"
	tr.toFile = "testdata/to.txt"

	ch := make(chan *twilioSMS, 10)
	tr.enqueue = func(m *twilioSMS) bool {
		ch <- m
		return true
	}

	err = tr.Deliver([]*kkok.Alert{{Message: "aaa"}})
	if err != nil {
		t.Fatal(err)
	}
	ms := testDequeue(ch)
	if len(ms) != 3 {
		t.Error(`len(ms) != 3`)
	}
}

func testDeliverMaxLength(t *testing.T) {
	tmpl, err := template.New("").Parse(`{{.Message}}`)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]struct {
		msg    string
		max    int
		expect string
	}{
		"Short":     {"short", 6, "short"},
		"UCS2Short": {"ショート", 6, "ショート"},
		"Long":      {"longlonglong", 6, "longlo"},
		"UCS2Long":  {"ロングロングロング", 6, "ロングロング"},
	}

	subtest := func(t *testing.T, msg, expect string, max int) {
		t.Parallel()
		tr, err := newTransport("abc")
		if err != nil {
			t.Fatal(err)
		}
		tr.token = "token"
		tr.from = "+123456789"
		tr.to = []string{"+100000000"}
		tr.tmpl = tmpl
		tr.maxLength = max

		ch := make(chan *twilioSMS, 1)
		tr.enqueue = func(m *twilioSMS) bool {
			ch <- m
			return true
		}

		err = tr.Deliver([]*kkok.Alert{{Message: msg}})
		if err != nil {
			t.Fatal(err)
		}
		ms := testDequeue(ch)
		if len(ms) != 1 {
			t.Fatal(`len(ms) != 1`)
		}

		v, _ := url.ParseQuery(ms[0].payload)
		if v.Get("Body") != expect {
			t.Error(`v.Get("Body") != expect`)
		}
	}

	for k, v := range data {
		t.Run(k, func(t *testing.T) {
			subtest(t, v.msg, v.expect, v.max)
		})
	}
}

func TestPost(t *testing.T) {
	account := os.Getenv("TEST_TWILIO_ACCOUNT")
	if len(account) == 0 {
		t.Skip("No TEST_TWILIO_ACCOUNT envvar")
	}
	sid := os.Getenv("TEST_TWILIO_SID")
	token := os.Getenv("TEST_TWILIO_TOKEN")
	if len(token) == 0 {
		t.Skip("No TEST_TWILIO_TOKEN envvar")
	}
	from := os.Getenv("TEST_TWILIO_FROM")
	if len(from) == 0 {
		t.Skip("No TEST_TWILIO_FROM envvar")
	}
	to := os.Getenv("TEST_TWILIO_TO")
	if len(to) == 0 {
		t.Skip("No TEST_TWILIO_TO envvar")
	}

	tr, err := newTransport(account)
	if err != nil {
		t.Fatal(err)
	}
	if len(sid) > 0 {
		tr.sid = sid
	}
	tr.token = token
	tr.from = from
	tr.to = []string{to}

	ch := make(chan *twilioSMS, 1)
	tr.enqueue = func(m *twilioSMS) bool {
		ch <- m
		return true
	}

	err = tr.Deliver([]*kkok.Alert{
		{
			From:    "from1",
			Title:   "test & test",
			Date:    time.Date(2011, 2, 3, 4, 5, 6, 120000000, time.UTC),
			Host:    "localhost",
			Message: "こんにちは",
			Info: map[string]interface{}{
				"info1": true,
				"info2": "long message\naaa <bbb>\n",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	m := <-ch
	if !m.send(context.Background()) {
		t.Error(`!m.send(context.Background())`)
	}
}
