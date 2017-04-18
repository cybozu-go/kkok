package slack

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"testing"
	"text/template"
	"time"

	"github.com/cybozu-go/kkok"
)

func TestTransport(t *testing.T) {
	t.Run("String", testString)
	t.Run("Params", testParams)
	t.Run("Format", testFormat)
	t.Run("Deliver", testDeliver)
}

func testString(t *testing.T) {
	t.Parallel()

	tr := &transport{}
	if tr.String() != transportType {
		t.Error(`tr.String() != transportType`)
	}

	tr = &transport{
		label: "label",
	}
	if tr.String() != "label" {
		t.Error(`tr.String() != "label"`)
	}
}

func testParams(t *testing.T) {
	t.Parallel()

	u, _ := url.Parse("https://slack.com/hoge/fuga")
	tr := &transport{
		url:       u,
		label:     "label",
		maxRetry:  0,
		name:      "name",
		icon:      ":sushi:",
		channel:   "#channel",
		origColor: `"danger"`,
		tmplPath:  "/path/to/template",
	}

	pp := tr.Params()
	if pp.Type != transportType {
		t.Error(`pp.Type != transportType`)
	}

	if pp.Params["url"] != "https://slack.com/hoge/fuga" {
		t.Error(`pp.Params["url"] != "https://slack.com/hoge/fuga"`)
	}
	if pp.Params["label"] != "label" {
		t.Error(`pp.Params["label"] != "label"`)
	}
	if pp.Params["max_retry"] != 0 {
		t.Error(`pp.Params["max_retry"] != 0`)
	}
	if pp.Params["name"] != "name" {
		t.Error(`pp.Params["name"] != "name"`)
	}
	if pp.Params["icon"] != ":sushi:" {
		t.Error(`pp.Params["icon"] != ":sushi:"`)
	}
	if pp.Params["channel"] != "#channel" {
		t.Error(`pp.Params["channel"] != "#channel"`)
	}
	if pp.Params["color"] != "\"danger\"" {
		t.Error(`pp.Params["color"] != "\"danger\""`)
	}
	if pp.Params["template"] != "/path/to/template" {
		t.Error(`pp.Params["template"] != "/path/to/template"`)
	}

	tr = &transport{
		url: u,
	}
	pp = tr.Params()
	if len(pp.Params) != 2 {
		t.Error(`len(pp.Params) != 2`)
	}
	if _, ok := pp.Params["max_retry"]; !ok {
		t.Error(`_, ok := pp.Params["max_retry"]; !ok`)
	}
}

func testFormat(t *testing.T) {
	t.Run("Color", testFormatColor)
	t.Run("Template", testFormatTemplate)
	t.Run("Default", testFormatDefault)
	t.Run("Custom", testFormatCustom)
}

func testFormatColor(t *testing.T) {
	t.Parallel()

	s, err := kkok.CompileJS(`1`)
	if err != nil {
		t.Fatal(err)
	}

	u, _ := url.Parse("https://slack.com/hoge/fuga")
	tr := &transport{
		url:   u,
		color: s,
	}

	_, err = tr.format(&kkok.Alert{})
	if err == nil {
		t.Error(`err == nil`)
	}

	s, err = kkok.CompileJS(`a.Info.severity`)
	if err != nil {
		t.Fatal(err)
	}
	tr.color = s
	_, err = tr.format(&kkok.Alert{})
	if err == nil {
		t.Error(`err == nil`)
	}

	s, err = kkok.CompileJS(`alert.Info.severity`)
	if err != nil {
		t.Fatal(err)
	}
	tr.color = s
	at, err := tr.format(&kkok.Alert{})
	if err != nil {
		t.Fatal(err)
	}
	if len(at.Color) != 0 {
		t.Error(`len(at.Color) != 0`, at.Color)
	}

	at, err = tr.format(&kkok.Alert{Info: map[string]interface{}{
		"severity": "good",
	}})
	if err != nil {
		t.Fatal(err)
	}
	if at.Color != "good" {
		t.Error(`at.Color != "good"`)
	}
}

func testFormatTemplate(t *testing.T) {
	t.Parallel()

	u, _ := url.Parse("https://slack.com/hoge/fuga")
	tr := &transport{
		url: u,
	}

	at, err := tr.format(&kkok.Alert{
		Message: "msg & msg",
	})
	if err != nil {
		t.Fatal(err)
	}
	if at.Text != "msg &amp; msg" {
		t.Error(`at.Text != "msg &amp; msg"`)
	}

	tmpl := template.Must(newTemplate().Parse(`hoge & fuga`))
	tr.tmpl = tmpl
	at, err = tr.format(&kkok.Alert{
		Message: "msg & msg",
	})
	if err != nil {
		t.Fatal(err)
	}
	if at.Text != "hoge & fuga" {
		t.Error(`at.Text != "hoge & fuga"`)
	}
}

func testFormatDefault(t *testing.T) {
	t.Parallel()

	u, _ := url.Parse("https://slack.com/hoge/fuga")
	tr := &transport{
		url: u,
	}

	testTitle := "title & title"
	at, err := tr.format(&kkok.Alert{
		From:  "from",
		Title: testTitle,
	})
	if err != nil {
		t.Fatal(err)
	}

	if at.Fallback != testTitle {
		t.Error(`at.Fallback != testTitle`)
	}
	if at.Title != EscapeSlack(testTitle) {
		t.Error(`at.Title != EscapeSlack(testTitle)`)
	}

	if len(at.Color) != 0 {
		t.Error(`len(at.Color) != 0`)
	}
	if len(at.Fields) != 1 {
		t.Error(`len(at.Fields) != 1`)
	}
	f0 := at.Fields[0]
	if f0.Title != "From" {
		t.Error(`f0.Title != "From"`)
	}
	if f0.Value != "from" {
		t.Error(`f0.Value != "from"`)
	}
}

func testFormatCustom(t *testing.T) {
	t.Parallel()

	u, _ := url.Parse("https://slack.com/hoge/fuga")
	tr := &transport{
		url: u,
	}

	at, err := tr.format(&kkok.Alert{
		From:  "kkok",
		Date:  time.Date(2011, 2, 3, 4, 5, 6, 123456000, time.UTC),
		Host:  "localhost",
		Title: "title1",
		Info: map[string]interface{}{
			"info1": true,
			"info2": "msg & msg",
			"info3": "long\nmessage\n",
			"info4": []int{1, 2, 3},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range at.Fields {
		switch f.Title {
		case "From":
			if f.Value != "kkok" {
				t.Error(`f.Value != "kkok"`)
			}

		case "Date":
			if f.Value != "2011-02-03T04:05:06.123Z" {
				t.Error(`f.Value != "2011-02-03T04:05:06.123Z"`)
			}

		case "Host":
			if f.Value != "localhost" {
				t.Error(`f.Value != "localhost"`)
			}

		case "Info.info1":
			if f.Value != "true" {
				t.Error(`f.Value != "true"`)
			}

		case "Info.info2":
			if f.Value != "msg &amp; msg" {
				t.Error(`f.Value != "msg &amp; msg"`)
			}

		case "Info.info3":
			if f.Value != "long\nmessage\n" {
				t.Error(`f.Value != "long\nmessage\n"`)
			}

		case "Info.info4":

		default:
			t.Error("unexpected field: ", f.Title)
		}
	}
}

func testDeliver(t *testing.T) {
	t.Parallel()

	ch := make(chan *slackMessage, 10)
	eq := func(m *slackMessage) bool {
		ch <- m
		return true
	}

	alerts1 := make([]*kkok.Alert, 1)
	alerts2 := make([]*kkok.Alert, maxAttachments*2)
	alerts3 := make([]*kkok.Alert, maxAttachments*2+1)
	for i := range alerts1 {
		alerts1[i] = &kkok.Alert{
			From:  "kkok",
			Title: strconv.Itoa(i),
		}
	}
	for i := range alerts2 {
		alerts2[i] = &kkok.Alert{
			From:  "kkok",
			Title: strconv.Itoa(i),
		}
	}
	for i := range alerts3 {
		alerts3[i] = &kkok.Alert{
			From:  "kkok",
			Title: strconv.Itoa(i),
		}
	}

	u, _ := url.Parse("https://slack.com/hoge/fuga")
	tr := &transport{
		url:     u,
		enqueue: eq,
	}

	dequeue := func() int {
		var n int
		for {
			select {
			case <-ch:
				n++
			default:
				return n
			}
		}
	}

	err := tr.Deliver(alerts1)
	if err != nil {
		t.Fatal(err)
	}
	if dequeue() != 1 {
		t.Error(`dequeue() != 1`)
	}

	err = tr.Deliver(alerts2)
	if err != nil {
		t.Fatal(err)
	}
	if dequeue() != 2 {
		t.Error(`dequeue() != 2`)
	}

	err = tr.Deliver(alerts3)
	if err != nil {
		t.Fatal(err)
	}
	if dequeue() != 3 {
		t.Error(`dequeue() != 3`)
	}
}

func TestPost(t *testing.T) {
	us := os.Getenv("TEST_SLACK_URL")
	if len(us) == 0 {
		t.Skip("No TEST_SLACK_URL envvar")
	}

	u, err := url.Parse(us)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan *slackMessage, 10)
	eq := func(m *slackMessage) bool {
		ch <- m
		return true
	}

	color, err := kkok.CompileJS(`"danger"`)
	if err != nil {
		t.Fatal(err)
	}

	tr := &transport{
		url:     u,
		name:    "kkok",
		icon:    ":sushi:",
		channel: "#random",
		color:   color,
		enqueue: eq,
	}

	err = tr.Deliver([]*kkok.Alert{
		{
			From:  "from1",
			Title: "test & test",
			Date:  time.Date(2011, 2, 3, 4, 5, 6, 120000000, time.UTC),
			Host:  "localhost",
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
