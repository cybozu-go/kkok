package slack

import (
	"bytes"
	"testing"
	"time"

	"github.com/cybozu-go/kkok"
)

func TestEscapeSlack(t *testing.T) {
	t.Parallel()

	if EscapeSlack("hoge<>&") != "hoge&lt;&gt;&amp;" {
		t.Error(`EscapeSlack("hoge<>&") != "hoge&lt;&gt;&amp;"`)
	}

	if EscapeSlack("Hello & <world> 🌊") != "Hello &amp; &lt;world&gt; 🌊" {
		t.Error(`EscapeSlack("Hello & <world> 🌊") != "Hello &amp; &lt;world&gt; 🌊"`)
	}

	// " should not be escaped to &quot;
	if EscapeSlack("\"") != "\"" {
		t.Error(`EscapeSlack("\"") != "\""`)
	}
}

func TestDefaultTemplate(t *testing.T) {
	t.Parallel()

	a := &kkok.Alert{
		From:    "from1",
		Title:   "title1",
		Date:    time.Date(2016, 2, 3, 4, 5, 6, 123456, time.UTC),
		Host:    "localhost",
		Message: "foo\nbar\nHello & world!",
		Info: map[string]interface{}{
			"info1": true,
		},
		Sub: []*kkok.Alert{
			{From: "from2", Message: "msg"},
		},
	}

	buf := new(bytes.Buffer)
	err := defaultTemplate.Execute(buf, a)
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != "foo\nbar\nHello &amp; world!" {
		t.Error(`buf.String() != "foo\nbar\nHello &amp; world!"`)
	}
}
