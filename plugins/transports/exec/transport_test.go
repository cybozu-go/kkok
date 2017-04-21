package exec

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"reflect"
	"testing"
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

	tr := newTransport("curl", "-d", "@-", "-f", "http://some.host/")
	pp := tr.Params()
	if pp.Type != transportType {
		t.Error(`pp.Type != transportType`)
	}

	if !reflect.DeepEqual(pp.Params, map[string]interface{}{
		"command": []string{"curl", "-d", "@-", "-f", "http://some.host/"},
		"timeout": int(defaultTimeout.Seconds()),
	}) {
		t.Error(`unexpected params`, pp.Params)
	}

	tr.label = "label"
	tr.timeout = 3 * time.Second
	tr.all = true
	pp = tr.Params()

	if !reflect.DeepEqual(pp.Params, map[string]interface{}{
		"label":   "label",
		"command": []string{"curl", "-d", "@-", "-f", "http://some.host/"},
		"timeout": 3,
		"all":     true,
	}) {
		t.Error(`unexpected params`, pp.Params)
	}
}

func testDeliver(t *testing.T) {
	_, err := exec.LookPath("curl")
	if err != nil {
		t.Skip("curl is not available")
	}
	t.Run("One", testDeliverOne)
	t.Run("All", testDeliverAll)
	t.Run("Timeout", testDeliverTimeout)
	t.Run("Error", testDeliverError)
}

type testHandler struct {
	expectAlert  *kkok.Alert
	expectAlerts []*kkok.Alert
	wait         time.Duration
	statusCode   int
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.statusCode != http.StatusOK && h.statusCode != 0 {
		http.Error(w, "failed", h.statusCode)
		return
	}

	if h.wait != 0 {
		time.Sleep(h.wait)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.expectAlert != nil {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "bad content type", http.StatusBadRequest)
			return
		}

		a := new(kkok.Alert)
		err := json.Unmarshal(data, a)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !reflect.DeepEqual(a, h.expectAlert) {
			http.Error(w, "unexpected alert", http.StatusBadRequest)
			return
		}
	}

	if h.expectAlerts != nil {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "bad content type", http.StatusBadRequest)
			return
		}

		var alerts []*kkok.Alert
		err := json.Unmarshal(data, &alerts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !reflect.DeepEqual(alerts, h.expectAlerts) {
			http.Error(w, "unexpected alerts", http.StatusBadRequest)
			return
		}
	}
}

func testSetup(h *testHandler) (*transport, *httptest.Server) {
	serv := httptest.NewServer(h)
	tr := newTransport("curl", "--data-binary", "@-",
		"-f", "-s", "-o", os.DevNull,
		"-H", "Content-Type: application/json", serv.URL)
	return tr, serv
}

func testDeliverOne(t *testing.T) {
	t.Parallel()

	a := &kkok.Alert{
		From:    "from",
		Title:   "title",
		Date:    time.Date(2011, 2, 3, 4, 5, 6, 123456000, time.UTC),
		Host:    "localhost",
		Message: "msg",
		Info: map[string]interface{}{
			"info1": true,
		},
	}
	h := &testHandler{
		expectAlert: a,
	}
	tr, serv := testSetup(h)
	defer serv.Close()

	err := tr.Deliver([]*kkok.Alert{a, a})
	if err != nil {
		t.Error(err)
	}
}

func testDeliverAll(t *testing.T) {
	t.Parallel()

	a := &kkok.Alert{
		From:    "from",
		Title:   "title",
		Date:    time.Date(2011, 2, 3, 4, 5, 6, 123456000, time.UTC),
		Host:    "localhost",
		Message: "msg",
		Info: map[string]interface{}{
			"info1": true,
		},
	}
	h := &testHandler{
		expectAlerts: []*kkok.Alert{a, a},
	}
	tr, serv := testSetup(h)
	defer serv.Close()
	tr.all = true

	err := tr.Deliver([]*kkok.Alert{a, a})
	if err != nil {
		t.Error(err)
	}
}

func testDeliverTimeout(t *testing.T) {
	t.Parallel()

	a := &kkok.Alert{
		From:    "from",
		Title:   "title",
		Date:    time.Date(2011, 2, 3, 4, 5, 6, 123456000, time.UTC),
		Host:    "localhost",
		Message: "msg",
		Info: map[string]interface{}{
			"info1": true,
		},
	}
	h := &testHandler{
		expectAlert: a,
		wait:        200 * time.Millisecond,
	}
	tr, serv := testSetup(h)
	defer serv.Close()
	tr.timeout = 100 * time.Millisecond

	err := tr.Deliver([]*kkok.Alert{a, a})
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testDeliverError(t *testing.T) {
	t.Parallel()

	a := &kkok.Alert{
		From:    "from",
		Title:   "title",
		Date:    time.Date(2011, 2, 3, 4, 5, 6, 123456000, time.UTC),
		Host:    "localhost",
		Message: "msg",
		Info: map[string]interface{}{
			"info1": true,
		},
	}
	h := &testHandler{
		expectAlert: a,
		statusCode:  500,
	}
	tr, serv := testSetup(h)
	defer serv.Close()

	err := tr.Deliver([]*kkok.Alert{a, a})
	if err == nil {
		t.Error(`err == nil`)
	}
}
