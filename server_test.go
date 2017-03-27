package kkok

import (
	"bytes"
	"encoding/json"
	"mime"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testAlertHandler struct {
	alerts []*Alert
}

func (h *testAlertHandler) Handle(alerts []*Alert) {
	h.alerts = alerts
}

func jsonRequest(method, p string, j string) *http.Request {
	if len(p) == 0 || p[0] != '/' {
		p = "/" + p
	}
	rd := bytes.NewReader([]byte(j))
	r := httptest.NewRequest(method, "http://localhost"+p, rd)
	r.Header.Set("Content-Type", "application/json")
	return r
}

func record(token string, r *http.Request) *httptest.ResponseRecorder {
	k := NewKkok()
	d := NewDispatcher(0, 0, new(testAlertHandler))
	h := &apiHandler{token, k, d}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	return w
}

func recordWithKkok(k *Kkok, r *http.Request) *httptest.ResponseRecorder {
	d := NewDispatcher(0, 0, new(testAlertHandler))
	h := &apiHandler{"", k, d}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	return w
}

func recordWithDispatcher(d *Dispatcher, r *http.Request) *httptest.ResponseRecorder {
	k := NewKkok()
	h := &apiHandler{"", k, d}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	return w
}

func assertContentType(t *testing.T, expected string, rr *httptest.ResponseRecorder) {
	actual, _, err := mime.ParseMediaType(rr.HeaderMap.Get("Content-Type"))
	if err != nil {
		t.Fatal(err)
	}

	if expected != actual {
		t.Error("Content-Type mismatch: expected=" + expected + ", actual=" + actual)
	}
}

func testRecvJSON(t *testing.T, rr *httptest.ResponseRecorder, i interface{}) {
	if rr.Code != http.StatusOK {
		t.Fatalf("not 200: %d", rr.Code)
	}
	assertContentType(t, "application/json", rr)
	err := json.Unmarshal(rr.Body.Bytes(), i)
	if err != nil {
		t.Fatal(err)
	}
}

func testServerVersionGet(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "http://localhost/version", nil)
	w := record("", r)
	if w.Code != http.StatusOK {
		t.Fatal(`w.Code != http.StatusOK`)
	}
	assertContentType(t, "text/plain", w)
	if w.Body.String() != Version {
		t.Error(`w.Body.String() != Version`)
	}
}

func testServerVersionOverrideGet(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("POST", "http://localhost/version", nil)
	r.Header.Set("X-HTTP-Method-Override", "GET")
	w := record("", r)
	if w.Code != http.StatusOK {
		t.Fatal(`w.Code != http.StatusOK`)
	}
	assertContentType(t, "text/plain", w)
	if w.Body.String() != Version {
		t.Error(`w.Body.String() != Version`)
	}
}

func testServerAuthToken(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "http://localhost/notfound", nil)
	w := record("", r)
	if w.Code != http.StatusNotFound {
		t.Error(`w.Code != http.StatusNotFound`)
	}

	r = httptest.NewRequest("GET", "http://localhost/notfound", nil)
	w = record("hoge", r)
	if w.Code != http.StatusForbidden {
		t.Error(`w.Code != http.StatusForbidden`)
	}

	r = httptest.NewRequest("GET", "http://localhost/notfound", nil)
	r.Header.Set("Authorization", "Bearer hoge")
	w = record("hoge", r)
	if w.Code != http.StatusNotFound {
		t.Error(`w.Code != http.StatusNotFound`)
	}

	r = httptest.NewRequest("GET", "http://localhost/notfound", nil)
	r.Header.Set("Authorization", "hoge")
	w = record("hoge", r)
	if w.Code != http.StatusForbidden {
		t.Error(`w.Code != http.StatusForbidden`)
	}
}

func testServerAlertsGet(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "http://localhost/alerts", nil)
	d := NewDispatcher(0, 0, new(testAlertHandler))
	d.put(&Alert{
		From:  "hoge",
		Host:  "localhost",
		Title: "aaa",
	})
	w := recordWithDispatcher(d, r)
	if w.Code != http.StatusOK {
		t.Error(`w.Code != http.StatusOK`)
	}

	var alerts []*Alert
	testRecvJSON(t, w, &alerts)

	if len(alerts) != 1 {
		t.Fatal(`len(alerts) != 1`)
	}

	a := alerts[0]

	if a.From != "hoge" {
		t.Error(`a.From != "hoge"`)
	}
	if a.Host != "localhost" {
		t.Error(`a.Host != "localhost"`)
	}
	if a.Title != "aaa" {
		t.Error(`a.Title != "aaa"`)
	}
}

func testServerAlertsPost(t *testing.T) {
	t.Parallel()

	j := `{
    "From": "hoge",
    "Host": "localhost",
    "Title": "aaa",
    "Info": {"option1": true, "option2": 1},
    "Stats": {"stat1": 1.23}
}`
	r := jsonRequest("POST", "/alerts", j)
	d := NewDispatcher(0, 0, new(testAlertHandler))
	w := recordWithDispatcher(d, r)
	if w.Code != http.StatusOK {
		t.Fatal(`w.Code != http.StatusOK`)
	}

	alerts := d.take()
	if len(alerts) != 1 {
		t.Fatal(`len(alerts) != 1`)
	}

	a := alerts[0]
	if a.From != "hoge" {
		t.Error(`a.From != "hoge"`)
	}
	if a.Host != "localhost" {
		t.Error(`a.Host != "localhost"`)
	}
	if a.Title != "aaa" {
		t.Error(`a.Title != "aaa"`)
	}
	if opt1, ok := a.Info["option1"]; !ok {
		t.Error(`opt1, ok := a.Info["option1"]; !ok`)
	} else if !opt1.(bool) {
		t.Error(`!opt1.(bool)`)
	}
	if opt2, ok := a.Info["option2"]; !ok {
		t.Error(`opt2, ok := a.Info["option2"]; !ok`)
	} else if int(opt2.(float64)) != 1 {
		t.Error(`int(opt2.(float64)) != 1`)
	}
	if a.Stats != nil {
		t.Error(`a.Stats != nil`)
	}

	j2 := `{
    "From": "fuga",
    "Host": "localhost",
    "Title": "aaa"
}`
	r = jsonRequest("POST", "/alerts", j2)
	w = recordWithDispatcher(d, r)
	if w.Code != http.StatusOK {
		t.Fatal(`w.Code != http.StatusOK`)
	}

	alerts = d.take()
	if len(alerts) != 1 {
		t.Fatal(`len(alerts) != 1`)
	}
	a = alerts[0]
	if a.From != "fuga" {
		t.Error(`a.From != "fuga"`)
	}
	if a.Info != nil {
		t.Error(`a.Info != nil`)
	}
}

func testServerAlertsBad(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("PUT", "http://localhost/alerts", nil)
	w := record("", r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Error(`w.Code != http.StatusMethodNotAllowed`)
	}
}

func hasString(s string, ss []string) bool {
	for _, t := range ss {
		if t == s {
			return true
		}
	}
	return false
}

func testServerFiltersGet(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "http://localhost/filters", nil)
	k := NewKkok()
	f1 := &dupFilter{}
	f1.Init("f1", nil)
	err := k.AddStaticFilter(f1)
	if err != nil {
		t.Fatal(err)
	}

	f2 := &dupFilter{}
	f2.Init("f2", nil)
	k.PutFilter(f2)

	w := recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Error(`w.Code != http.StatusOK`)
	}

	var filterIDs []string
	testRecvJSON(t, w, &filterIDs)

	if len(filterIDs) != 2 {
		t.Error(`len(filterIDs) != 2`)
	}
	if !hasString("f1", filterIDs) {
		t.Error(`!hasString("f1", filterIDs)`)
	}
	if !hasString("f2", filterIDs) {
		t.Error(`!hasString("f2", filterIDs)`)
	}
}

func testServerFiltersIDGet(t *testing.T) {
	t.Parallel()

	k := NewKkok()
	f1 := &dupFilter{}
	f1.Init("f1", map[string]interface{}{
		"all": true,
	})
	err := k.AddStaticFilter(f1)
	if err != nil {
		t.Fatal(err)
	}

	f2 := &dupFilter{}
	f2.Init("f2", nil)
	k.PutFilter(f2)

	r := httptest.NewRequest("GET", "http://localhost/filters/f1", nil)
	w := recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Error(`w.Code != http.StatusOK`)
	}

	var pp PluginParams
	testRecvJSON(t, w, &pp)
	if pp.Type != "dup" {
		t.Error(`pp.Type != "dup"`)
	}
	if !pp.Params["all"].(bool) {
		t.Error(`!pp.Params["all"].(bool)`)
	}

	r = httptest.NewRequest("GET", "http://localhost/filters/nosuchfilter", nil)
	w = recordWithKkok(k, r)
	if w.Code != http.StatusNotFound {
		t.Error(`w.Code != http.StatusNotFound`)
	}
}

func testServerFiltersIDPut(t *testing.T) {
	t.Parallel()

	j := `{
    "type": "dup",
    "disabled": true
}`
	r := jsonRequest("PUT", "/filters/f1", j)
	k := NewKkok()
	w := recordWithKkok(k, r)
	if w.Code != http.StatusInternalServerError {
		t.Error(`w.Code != http.StatusInternalServerError`)
	}

	RegisterFilter("dup", func(id string, params map[string]interface{}) (Filter, error) {
		f := &dupFilter{}
		err := f.Init(id, params)
		if err != nil {
			return nil, err
		}
		return f, nil
	})

	r = jsonRequest("PUT", "/filters/f1", j)
	w = recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Fatal(`w.Code != http.StatusOK`)
	}

	if len(k.Filters()) != 1 {
		t.Error(`len(k.Filters()) != 1`)
	}

	// overwrite
	r = jsonRequest("PUT", "/filters/f1", j)
	w = recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Fatal(`w.Code != http.StatusOK`)
	}

	if len(k.Filters()) != 1 {
		t.Error(`len(k.Filters()) != 1`)
	}

	r = jsonRequest("PUT", "/filters/f2", j)
	w = recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Fatal(`w.Code != http.StatusOK`)
	}

	if len(k.Filters()) != 2 {
		t.Error(`len(k.Filters()) != 2`)
	}

	// overwrite
	r = jsonRequest("PUT", "/filters/f2", j)
	w = recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Fatal(`w.Code != http.StatusOK`)
	}

	if len(k.Filters()) != 2 {
		t.Error(`len(k.Filters()) != 2`)
	}
}

func testServerFiltersIDDelete(t *testing.T) {
	t.Parallel()

	k := NewKkok()
	f1 := &dupFilter{}
	f1.Init("f1", map[string]interface{}{
		"all": true,
	})
	err := k.AddStaticFilter(f1)
	if err != nil {
		t.Fatal(err)
	}

	f2 := &dupFilter{}
	f2.Init("f2", nil)
	k.PutFilter(f2)

	r := httptest.NewRequest("DELETE", "http://localhost/filters/nosuchfilter", nil)
	w := recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Error(`w.Code != http.StatusOK`)
	}

	if len(k.Filters()) != 2 {
		t.Error(`len(k.Filters()) != 2`)
	}

	// static filter f1 cannot be removed
	r = httptest.NewRequest("DELETE", "http://localhost/filters/f1", nil)
	w = recordWithKkok(k, r)
	if w.Code != http.StatusInternalServerError {
		t.Error(`w.Code != http.StatusInternalServerError`)
	}

	if len(k.Filters()) != 2 {
		t.Error(`len(k.Filters()) != 2`)
	}

	r = httptest.NewRequest("DELETE", "http://localhost/filters/f2", nil)
	w = recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Error(`w.Code != http.StatusOK`)
	}

	if len(k.Filters()) != 1 {
		t.Error(`len(k.Filters()) != 1`)
	}

	if k.Filters()[0].ID() != "f1" {
		t.Error(`k.Filters()[0].ID() != "f1"`)
	}
}

func testServerRoutesGet(t *testing.T) {
	t.Parallel()

	tr1 := &testTransport{}
	tr2 := &testTransport{}
	k := NewKkok()
	err := k.AddRoute("r1", []Transport{tr1})
	if err != nil {
		t.Error(err)
	}
	err = k.AddRoute("r2", []Transport{tr2})
	if err != nil {
		t.Error(err)
	}

	r := httptest.NewRequest("GET", "http://localhost/routes", nil)
	w := recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Error(`w.Code != http.StatusOK`)
	}

	var ids []string
	testRecvJSON(t, w, &ids)

	if len(ids) != 2 {
		t.Error(`len(ids) != 2`)
	}
	if !hasString("r1", ids) {
		t.Error(`!hasString("r1", ids)`)
	}
	if !hasString("r2", ids) {
		t.Error(`!hasString("r2", ids)`)
	}
}

func testServerRoutesIDGet(t *testing.T) {
	t.Parallel()

	tr1 := &testTransport{}
	tr2 := &testTransport{}
	k := NewKkok()
	err := k.AddRoute("r1", []Transport{tr1})
	if err != nil {
		t.Error(err)
	}
	err = k.AddRoute("r2", []Transport{tr2})
	if err != nil {
		t.Error(err)
	}

	r := httptest.NewRequest("GET", "http://localhost/routes/r1", nil)
	w := recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Error(`w.Code != http.StatusOK`)
	}

	var pl []PluginParams
	testRecvJSON(t, w, &pl)

	if len(pl) != 1 {
		t.Fatal(`len(pl) != 1`)
	}
	if pl[0].Type != "test" {
		t.Error(`pl[0].Type != "test"`)
	}
	if len(pl[0].Params) != 0 {
		t.Error(`len(pl[0].Params) != 0`)
	}

	r = httptest.NewRequest("GET", "http://localhost/routes/nosuchfilter", nil)
	w = recordWithKkok(k, r)
	if w.Code != http.StatusNotFound {
		t.Error(`w.Code != http.StatusNotFound`)
	}
}

func testServerRoutesIDPut(t *testing.T) {
	t.Parallel()

	j := `[
    {
        "type": "test"
    },
    {
        "type": "test2"
    }
]`
	r := jsonRequest("PUT", "/routes/r1", j)
	k := NewKkok()
	w := recordWithKkok(k, r)
	if w.Code != http.StatusInternalServerError {
		t.Error(`w.Code != http.StatusInternalServerError`)
	}

	RegisterTransport("test", func(params map[string]interface{}) (Transport, error) {
		return &testTransport{}, nil
	})
	RegisterTransport("test2", func(params map[string]interface{}) (Transport, error) {
		return &testTransport{}, nil
	})

	r = jsonRequest("PUT", "/routes/r1", j)
	w = recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Fatal(`w.Code != http.StatusOK`)
	}

	if len(k.RouteIDs()) != 1 {
		t.Error(`len(k.RouteIDs()) != 1`)
	}

	// overwrite
	r = jsonRequest("PUT", "/routes/r1", j)
	w = recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Fatal(`w.Code != http.StatusOK`)
	}

	if len(k.RouteIDs()) != 1 {
		t.Error(`len(k.RouteIDs()) != 1`)
	}

	r = jsonRequest("PUT", "/routes/r2", j)
	w = recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Fatal(`w.Code != http.StatusOK`)
	}

	if len(k.RouteIDs()) != 2 {
		t.Error(`len(k.RouteIDs()) != 2`)
	}

	// overwrite
	r = jsonRequest("PUT", "/routes/r2", j)
	w = recordWithKkok(k, r)
	if w.Code != http.StatusOK {
		t.Fatal(`w.Code != http.StatusOK`)
	}

	if len(k.RouteIDs()) != 2 {
		t.Error(`len(k.RouteIDs()) != 2`)
	}
}

func TestServer(t *testing.T) {
	t.Run("Version/Get", testServerVersionGet)
	t.Run("Version/OverrideGet", testServerVersionOverrideGet)
	t.Run("AuthToken", testServerAuthToken)
	t.Run("Alerts/Get", testServerAlertsGet)
	t.Run("Alerts/Post", testServerAlertsPost)
	t.Run("Alerts/Bad", testServerAlertsBad)
	t.Run("Filters/Get", testServerFiltersGet)
	t.Run("Filters/ID/Get", testServerFiltersIDGet)
	t.Run("Filters/ID/Put", testServerFiltersIDPut)
	t.Run("Filters/ID/Delete", testServerFiltersIDDelete)
	t.Run("Routes/Get", testServerRoutesGet)
	t.Run("Routes/ID/Get", testServerRoutesIDGet)
	t.Run("Routes/ID/Put", testServerRoutesIDPut)
}
