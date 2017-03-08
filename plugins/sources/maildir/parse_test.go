package maildir

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/cybozu-go/kkok"
)

var (
	tzP0900 = time.FixedZone("0900", 9*3600)
	tzM0100 = time.FixedZone("-0100", -1*3600)

	expects = map[string]*kkok.Alert{
		"simple": {
			From:    "ymmt2005",
			Date:    time.Date(2017, 2, 28, 13, 56, 18, 0, tzP0900),
			Title:   "Alert from example monitor",
			Message: "t\nt\n\nt\n",
		},

		"simple_dos": {
			From:    "ymmt2005",
			Date:    time.Date(2017, 2, 28, 13, 56, 18, 0, tzP0900),
			Title:   "Alert from example monitor",
			Message: "t\nt\n\nt\n",
		},

		"confusing": {
			From:    "ymmt2005@example.com",
			Date:    time.Date(2017, 2, 28, 13, 56, 29, 0, time.UTC),
			Title:   "Alert from example monitor",
			Message: "From: example monitor\nconfusing text\n",
		},

		"pseudo": {
			From:    "example monitor",
			Date:    time.Date(2017, 2, 28, 13, 56, 25, 0, time.UTC),
			Title:   "Alert from example monitor",
			Host:    "host-1",
			Message: "test\n",
			Info: map[string]interface{}{
				"Option1": "123",
				"Option2": "aaa",
			},
		},

		"pseudo_dos": {
			From:    "example monitor",
			Date:    time.Date(2017, 2, 28, 13, 56, 25, 0, time.UTC),
			Title:   "Alert from example monitor",
			Host:    "host-1",
			Message: "test\n",
			Info: map[string]interface{}{
				"Option1": "123",
				"Option2": "aaa",
			},
		},

		"pseudo_partial_dos": {
			From:    "example monitor",
			Date:    time.Date(2017, 2, 28, 13, 56, 25, 0, time.UTC),
			Title:   "Alert from example monitor",
			Host:    "host-1",
			Message: "test\n",
			Info: map[string]interface{}{
				"Option1": "123",
				"Option2": "aaa",
			},
		},

		"pseudo_only": {
			From:  "example monitor",
			Date:  time.Date(2011, 2, 3, 11, 22, 33, 0, time.UTC),
			Title: "hoge",
			Host:  "host-1",
			Info: map[string]interface{}{
				"Option1": "123",
				"Option2": "aaa",
			},
		},

		"pseudo_base64": {
			From:    "monitor",
			Date:    time.Date(2017, 3, 3, 14, 21, 57, 0, time.UTC),
			Title:   "♡テスト",
			Host:    "example.com",
			Message: "hoge\n",
		},

		"qprintable": {
			From:    "monitor",
			Date:    time.Date(2017, 2, 28, 5, 55, 2, 0, time.UTC),
			Title:   "テスト",
			Message: "テストテスト\ntest\ntest\n",
		},

		"base64": {
			From:    "monitor@example.com",
			Date:    time.Date(2017, 2, 28, 6, 10, 15, 0, tzM0100),
			Title:   "test2",
			Host:    "test-host",
			Message: "$ € 円\n",
		},
	}
)

func alertEqual(t *testing.T, actual, expected *kkok.Alert) {
	if actual.From != expected.From {
		t.Errorf(`actual.From != expected.From (%q, %q)`, actual.From, expected.From)
	}
	if !actual.Date.Equal(expected.Date) {
		t.Error(`!actual.Date.Equal(expected.Date)`,
			actual.Date.String(), expected.Date.String())
	}
	if actual.Host != expected.Host {
		t.Errorf(`actual.Host != expected.Host (%q, %q)`, actual.Host, expected.Host)
	}
	if actual.Title != expected.Title {
		t.Errorf(`actual.Title != expected.Title (%q, %q)`, actual.Title, expected.Title)
	}
	if actual.Message != expected.Message {
		t.Errorf(`actual.Message != expected.Message (%q, %q)`, actual.Message, expected.Message)
	}
	if (len(expected.Info) != 0 || len(actual.Info) != 0) &&
		!reflect.DeepEqual(actual.Info, expected.Info) {
		t.Error(`!reflect.DeepEqual(actual.Info, expected.Info)`,
			actual.Info, expected.Info)
	}
}

func testParse(t *testing.T, src string) {
	t.Parallel()

	f, err := os.Open("testdata/new/" + src + ".txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	a, err := parse(f)
	if err != nil {
		t.Fatal(err)
	}

	alertEqual(t, a, expects[src])
}

func TestParse(t *testing.T) {
	t.Run("Simple/Unix", func(t *testing.T) { testParse(t, "simple") })
	t.Run("Simple/Dos", func(t *testing.T) { testParse(t, "simple_dos") })
	t.Run("Confusing", func(t *testing.T) { testParse(t, "confusing") })
	t.Run("Pseudo/Unix", func(t *testing.T) { testParse(t, "pseudo") })
	t.Run("Pseudo/Dos", func(t *testing.T) { testParse(t, "pseudo_dos") })
	t.Run("Pseudo/DosUnix", func(t *testing.T) { testParse(t, "pseudo_partial_dos") })
	t.Run("Pseudo/Only", func(t *testing.T) { testParse(t, "pseudo_only") })
	t.Run("Pseudo/Base64", func(t *testing.T) { testParse(t, "pseudo_base64") })
	t.Run("QuotedPrintable", func(t *testing.T) { testParse(t, "qprintable") })
	t.Run("Base64", func(t *testing.T) { testParse(t, "base64") })
}
