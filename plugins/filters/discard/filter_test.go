package discard

import (
	"testing"

	"github.com/cybozu-go/kkok"
)

func testFilterAll(t *testing.T) {
	t.Parallel()

	f := &filter{}
	err := f.Init("f", map[string]interface{}{
		"all": true,
		"if":  "alerts.length > 2",
	})
	if err != nil {
		t.Fatal(err)
	}

	alerts := []*kkok.Alert{{}, {}}
	alerts, err = f.Process(alerts)
	if err != nil {
		t.Fatal(err)
	}
	if len(alerts) != 2 {
		t.Error(`len(alerts) != 2`)
	}

	alerts = []*kkok.Alert{{}, {}, {}}
	alerts, err = f.Process(alerts)
	if err != nil {
		t.Fatal(err)
	}
	if len(alerts) != 0 {
		t.Error(`len(alerts) != 0`)
	}
}

func testFilterOne(t *testing.T) {
	t.Parallel()

	f := &filter{}
	err := f.Init("f", map[string]interface{}{
		"if": "alert.From == 'from1'",
	})
	if err != nil {
		t.Fatal(err)
	}

	alerts := []*kkok.Alert{
		{From: "from1"},
		{From: "from2"},
		{From: "from3"},
		{From: "from1"},
		{From: "from4"},
	}

	alerts, err = f.Process(alerts)
	if err != nil {
		t.Fatal(err)
	}
	if len(alerts) != 3 {
		t.Error(`len(alerts) != 3`)
	}
	for i, a := range alerts {
		if a.From == "from1" {
			t.Error(`a.From == "from1"`, i)
		}
	}
}

func TestFilter(t *testing.T) {
	t.Run("All", testFilterAll)
	t.Run("One", testFilterOne)
}

func TestCtor(t *testing.T) {
	t.Parallel()

	f, err := ctor("test", nil)
	if err != nil {
		t.Fatal(err)
	}
	if f.ID() != "test" {
		t.Error(`f.ID() != "test"`)
	}

	pp := f.Params()
	if pp.Type != filterType {
		t.Error(`pp.Type != filterType`)
	}
}
