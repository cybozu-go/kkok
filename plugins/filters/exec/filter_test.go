package exec

import (
	"os/exec"
	"reflect"
	"testing"
	"time"

	"github.com/cybozu-go/kkok"
)

var testAlertsData = []*kkok.Alert{
	{From: "from1", Title: "title1"},
	{From: "from2", Title: "title2"},
	{From: "from3", Title: "title3"},
	{From: "from1", Title: "title4"},
}

func TestFilter(t *testing.T) {
	t.Run("Params", testParams)
	t.Run("Exec", testExec)
	t.Run("Process", testProcess)
}

func testParams(t *testing.T) {
	t.Run("Default", testParamsDefault)
	t.Run("Explicit", testParamsExplicit)
}

func testParamsDefault(t *testing.T) {
	t.Parallel()

	f := newFilter()
	f.command = []string{"foo", "bar"}
	pp := f.Params()

	if pp.Type != filterType {
		t.Error(`pp.Type != filterType`)
	}

	if !reflect.DeepEqual(pp.Params["command"], []string{"foo", "bar"}) {
		t.Error(`!reflect.DeepEqual(pp.Params["command"], []string{"foo", "bar"})`)
	}
	if pp.Params["timeout"].(int) != int(defaultTimeout.Seconds()) {
		t.Error(`pp.Params["timeout"].(int) != int(defaultTimeout.Seconds())`)
	}
}

func testParamsExplicit(t *testing.T) {
	t.Parallel()

	f := newFilter()
	f.command = []string{"foo", "bar"}
	f.timeout = 10 * time.Second
	pp := f.Params()

	if pp.Type != filterType {
		t.Error(`pp.Type != filterType`)
	}

	if !reflect.DeepEqual(pp.Params["command"], []string{"foo", "bar"}) {
		t.Error(`!reflect.DeepEqual(pp.Params["command"], []string{"foo", "bar"})`)
	}
	if pp.Params["timeout"].(int) != 10 {
		t.Error(`pp.Params["timeout"].(int) != 10`)
	}
}

func testExec(t *testing.T) {
	t.Run("Success", testExecSuccess)
	t.Run("Error", testExecError)
	t.Run("Timeout", testExecTimeout)
}

func testExecSuccess(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("cat"); err != nil {
		t.Skip("cat is not found")
	}

	f := newFilter()
	f.command = []string{"cat"}

	j, err := f.exec([]byte("abc"))
	if err != nil {
		t.Fatal(err)
	}
	if string(j) != "abc" {
		t.Error(`string(j) != "abc"`)
	}
}

func testExecError(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("/bin/sh"); err != nil {
		t.Skip("/bin/sh is not found")
	}

	f := newFilter()
	f.command = []string{"/bin/sh", "-c", `
echo foo 1>&2
exit 3
`}

	_, err := f.exec([]byte("abc"))
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testExecTimeout(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("sleep"); err != nil {
		t.Skip("sleep is not found")
	}

	f := newFilter()
	f.timeout = 10 * time.Millisecond
	f.command = []string{"sleep", "1"}

	_, err := f.exec([]byte("abc"))
	if err == nil {
		t.Error(`err == nil`)
	}
	t.Log(err)
}

func testProcess(t *testing.T) {
	t.Run("One", testProcessOne)
	t.Run("All", testProcessAll)
}

func testProcessOne(t *testing.T) {
	t.Run("Success", testProcessSuccess)
	t.Run("Failure", testProcessFailure)
	t.Run("InvalidJSON", testProcessInvalidJSON)
	t.Run("If", testProcessIf)
}

func testProcessSuccess(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq is not found")
	}

	f := newFilter()
	f.command = []string{"jq", `. + {"Message": "emergency"}`}

	alerts, err := f.Process(testAlertsData)
	if err != nil {
		t.Fatal(err)
	}

	if len(alerts) != len(testAlertsData) {
		t.Fatal(`len(alerts) != len(testAlertsData)`)
	}

	for i, a := range alerts {
		if a.From != testAlertsData[i].From {
			t.Error(`a.From != testAlertsData[i].From; i=`, i)
		}
		if a.Title != testAlertsData[i].Title {
			t.Error(`a.Title != testAlertsData[i].Title; i=`, i)
		}
		if a.Message != "emergency" {
			t.Error(`a.Message != "emergency"; i=`, i)
		}
	}
}

func testProcessFailure(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq is not found")
	}

	f := newFilter()
	f.command = []string{"jq", `[`}

	_, err := f.Process(testAlertsData)
	if err == nil {
		t.Error(`err == nil`)
	}
	t.Log(err)
}

func testProcessInvalidJSON(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("echo"); err != nil {
		t.Skip("echo is not found")
	}

	f := newFilter()
	f.command = []string{"echo", "aaa"}

	_, err := f.Process(testAlertsData)
	if err == nil {
		t.Error(`err == nil`)
	}
	t.Log(err)
}

func testProcessIf(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq is not found")
	}

	f := newFilter()
	f.command = []string{"jq", `. + {"Message": "emergency"}`}
	err := f.Init("f", map[string]interface{}{
		"if": `alert.From=="from2"`,
	})
	if err != nil {
		t.Fatal(err)
	}

	alerts, err := f.Process(testAlertsData)
	if err != nil {
		t.Fatal(err)
	}

	if len(alerts) != len(testAlertsData) {
		t.Fatal(`len(alerts) != len(testAlertsData)`)
	}

	for i, a := range alerts {
		if a.From != testAlertsData[i].From {
			t.Error(`a.From != testAlertsData[i].From; i=`, i)
		}
		if a.Title != testAlertsData[i].Title {
			t.Error(`a.Title != testAlertsData[i].Title; i=`, i)
		}
		if a.From == "from2" {
			if a.Message != "emergency" {
				t.Error(`a.Message != "emergency"; i=`, i)
			}
		} else {
			if len(a.Message) != 0 {
				t.Error(`len(a.Message) != 0; i=`, i)
			}
		}
	}
}

func testProcessAll(t *testing.T) {
	t.Run("Success", testProcessAllSuccess)
	t.Run("Failure", testProcessAllFailure)
	t.Run("InvalidJSON", testProcessAllInvalidJSON)
	t.Run("If", testProcessAllIf)
}

func testProcessAllSuccess(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq is not found")
	}

	f := newFilter()
	f.command = []string{"jq", `[.[1]]`}
	err := f.Init("f", map[string]interface{}{
		"all": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	alerts, err := f.Process(testAlertsData)
	if err != nil {
		t.Fatal(err)
	}

	if len(alerts) != 1 {
		t.Fatal(`len(alerts) != 1`)
	}

	if alerts[0].From != "from2" {
		t.Error(`alerts[0].From != "from2"`)
	}
	if alerts[0].Title != "title2" {
		t.Error(`alerts[0].Title != "title2"`)
	}
}

func testProcessAllFailure(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq is not found")
	}

	f := newFilter()
	f.command = []string{"jq", `[`}
	err := f.Init("f", map[string]interface{}{
		"all": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.Process(testAlertsData)
	if err == nil {
		t.Error(`err == nil`)
	}
	t.Log(err)
}

func testProcessAllInvalidJSON(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("echo"); err != nil {
		t.Skip("echo is not found")
	}

	f := newFilter()
	f.command = []string{"echo", "aaa"}
	err := f.Init("f", map[string]interface{}{
		"all": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.Process(testAlertsData)
	if err == nil {
		t.Error(`err == nil`)
	}
	t.Log(err)
}

func testProcessAllIf(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq is not found")
	}

	f := newFilter()
	f.command = []string{"jq", `[.[1]]`}
	err := f.Init("f", map[string]interface{}{
		"all": true,
		"if":  "alerts.length > 10",
	})
	if err != nil {
		t.Fatal(err)
	}

	alerts, err := f.Process(testAlertsData)
	if err != nil {
		t.Fatal(err)
	}

	if len(alerts) != len(testAlertsData) {
		t.Fatal(`len(alerts) != len(testAlertsData`)
	}
}
