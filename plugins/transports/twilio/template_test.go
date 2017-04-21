package twilio

import (
	"bytes"
	"testing"
	"time"

	"github.com/cybozu-go/kkok"
)

func TestTemplate(t *testing.T) {
	t.Parallel()

	a := &kkok.Alert{
		From:    "from",
		Date:    time.Date(2011, 2, 3, 4, 5, 6, 123456000, time.UTC),
		Host:    "localhost",
		Title:   "title",
		Message: "msg",
		Routes:  []string{"r1", "r2"},
		Info: map[string]interface{}{
			"info1": true,
			"info2": "foo",
		},
		Sub: []*kkok.Alert{{}, {}},
	}

	buf := new(bytes.Buffer)
	err := defaultTemplate.Execute(buf, a)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(buf.String())
}
