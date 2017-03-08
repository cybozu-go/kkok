package kkok

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/cybozu-go/cmd"
	"github.com/cybozu-go/log"
)

const (
	maxJSONLength = 1024 * 1024 * 10 // 10 MiB
)

func getMethod(r *http.Request) string {
	if r.Method == "POST" {
		m := r.Header.Get("X-HTTP-Method-Override")
		if len(m) > 0 {
			return m
		}
	}
	return r.Method
}

type apiHandler struct {
	token string
	k     *Kkok
	d     *Dispatcher
}

func sendJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		et := err.Error()
		fields := cmd.FieldsFromContext(r.Context())
		fields[log.FnError] = et
		log.Error("json.Marshal", fields)
		http.Error(w, et, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}

func recvJSON(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	data, err := ioutil.ReadAll(http.MaxBytesReader(w, r.Body, maxJSONLength))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return false
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		et := err.Error()
		fields := cmd.FieldsFromContext(r.Context())
		fields[log.FnError] = et
		log.Error("json.Unmarshal", fields)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return false
	}

	return true
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	if getMethod(r) != "GET" {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(Version))
}

func (a *apiHandler) authenticate(w http.ResponseWriter, r *http.Request) bool {
	if len(a.token) == 0 {
		return true
	}

	fields := strings.Fields(r.Header.Get("Authorization"))
	if len(fields) != 2 {
		http.Error(w, "auth token is required", http.StatusForbidden)
		return false
	}

	if fields[0] != "Bearer" {
		http.Error(w, "not a Bearer token", http.StatusForbidden)
		return false
	}

	if fields[1] != a.token {
		http.Error(w, "token mismatch", http.StatusForbidden)
		return false
	}

	return true
}

// ServeHTTP implements http.Handler.
func (a *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if p == "/version" {
		handleVersion(w, r)
		return
	}

	// End-points other than /version require authentication.
	if !a.authenticate(w, r) {
		return
	}

	if p == "/alerts" {
		a.handleAlerts(w, r)
		return
	}

	if p == "/filters" {
		a.getFilters(w, r)
		return
	}

	if strings.HasPrefix(p, "/filters/") {
		id := p[9:]
		if !reFilterID.MatchString(id) {
			http.Error(w, "invalid filter id: "+id, http.StatusBadRequest)
		} else {
			a.handleFilter(w, r, id)
		}
		return
	}

	if p == "/routes" {
		a.getRoutes(w, r)
		return
	}

	if strings.HasPrefix(p, "/routes/") {
		id := p[8:]
		if !reRouteID.MatchString(id) {
			http.Error(w, "invalid route id: "+id, http.StatusBadRequest)
		} else {
			a.handleRoute(w, r, id)
		}
		return
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func (a *apiHandler) handleAlerts(w http.ResponseWriter, r *http.Request) {
	switch getMethod(r) {
	case "GET":
		a.getAlerts(w, r)
	case "POST":
		a.postAlert(w, r)
	default:
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
	}
}

func (a *apiHandler) getAlerts(w http.ResponseWriter, r *http.Request) {
	alerts := a.d.peek()
	sendJSON(w, r, alerts)
}

func (a *apiHandler) postAlert(w http.ResponseWriter, r *http.Request) {
	alert := new(Alert)
	if !recvJSON(w, r, alert) {
		return
	}

	err := alert.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if alert.Date.IsZero() {
		alert.Date = time.Now().UTC()
	}

	if len(alert.Host) == 0 {
		h, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			h = r.RemoteAddr
		}
		alert.Host = h
	}

	alert.Routes = nil
	alert.Sub = nil

	a.d.put(alert)

	fields := cmd.FieldsFromContext(r.Context())
	fields["from"] = alert.From
	fields["title"] = alert.Title
	log.Info("new alert", fields)
}

func (a *apiHandler) getFilters(w http.ResponseWriter, r *http.Request) {
	if getMethod(r) != "GET" {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	filters := a.k.Filters()
	ids := make([]string, len(filters))
	for i, f := range filters {
		ids[i] = f.ID()
	}

	sendJSON(w, r, ids)
}

func (a *apiHandler) handleFilter(w http.ResponseWriter, r *http.Request, id string) {
	switch getMethod(r) {
	case "GET":
		a.getFilter(w, r, id)
	case "PUT":
		a.putFilter(w, r, id)
	case "DELETE":
		a.deleteFilter(w, r, id)
	default:
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
	}
}

func (a *apiHandler) getFilter(w http.ResponseWriter, r *http.Request, id string) {
	f := a.k.getFilter(id)
	if f == nil {
		http.NotFound(w, r)
		return
	}

	sendJSON(w, r, f.Params())
}

func (a *apiHandler) putFilter(w http.ResponseWriter, r *http.Request, id string) {
	if !reFilterID.MatchString(id) {
		http.Error(w, "invalid filter id: "+id, http.StatusBadRequest)
		return
	}

	var params PluginParams
	if !recvJSON(w, r, &params) {
		return
	}

	params.Params["id"] = id

	f, err := NewFilter(params.Type, params.Params)
	if err != nil {
		et := err.Error()
		fields := cmd.FieldsFromContext(r.Context())
		fields["filter_id"] = id
		fields[log.FnError] = et
		log.Error("failed to create a new filter", fields)
		http.Error(w, et, http.StatusInternalServerError)
		return
	}

	a.k.AddFilter(f)
}

func (a *apiHandler) deleteFilter(w http.ResponseWriter, r *http.Request, id string) {
	err := a.k.removeFilter(id)
	if err == nil {
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (a *apiHandler) getRoutes(w http.ResponseWriter, r *http.Request) {
	if getMethod(r) != "GET" {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	sendJSON(w, r, a.k.RouteIDs())
}

func (a *apiHandler) handleRoute(w http.ResponseWriter, r *http.Request, id string) {
	switch getMethod(r) {
	case "GET":
		a.getRoute(w, r, id)
	case "PUT":
		a.putRoute(w, r, id)
	default:
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
	}
}

func (a *apiHandler) getRoute(w http.ResponseWriter, r *http.Request, id string) {
	route := a.k.getRoute(id)
	if route == nil {
		http.NotFound(w, r)
		return
	}

	pl := make([]PluginParams, len(route))
	for i, t := range route {
		pl[i] = t.Params()
	}

	sendJSON(w, r, pl)
}

func (a *apiHandler) putRoute(w http.ResponseWriter, r *http.Request, id string) {
	if !reFilterID.MatchString(id) {
		http.Error(w, "invalid filter id: "+id, http.StatusBadRequest)
		return
	}

	var pl []PluginParams
	if !recvJSON(w, r, &pl) {
		return
	}

	route := make([]Transport, len(pl))
	for i, pp := range pl {
		tr, err := NewTransport(pp.Type, pp.Params)
		if err != nil {
			et := err.Error()
			fields := cmd.FieldsFromContext(r.Context())
			fields["route_id"] = id
			fields[log.FnError] = et
			log.Error("failed to create a new transport", fields)
			http.Error(w, et, http.StatusInternalServerError)
			return
		}
		route[i] = tr
	}

	a.k.AddRoute(id, route)
}

// NewHTTPServer returns *cmd.HTTPServer for REST API.
func NewHTTPServer(addr, apiToken string, k *Kkok, d *Dispatcher) (*cmd.HTTPServer, error) {
	s := &cmd.HTTPServer{
		Server: &http.Server{
			Addr:    addr,
			Handler: &apiHandler{apiToken, k, d},
		},
		ShutdownTimeout: 10 * time.Second,
	}
	return s, nil
}
