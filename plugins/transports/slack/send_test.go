package slack

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
)

var (
	validData = []byte("valid data")
)

type testHandler struct {
	errors       int
	rateLimitSec int
	badRequest   bool
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.errors > 0 {
		h.errors--
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if h.rateLimitSec > 0 {
		w.Header().Set("Retry-After", strconv.Itoa(h.rateLimitSec))
		h.rateLimitSec = 0
		http.Error(w, "rate limit exceeds", http.StatusTooManyRequests)
		return
	}

	if h.badRequest {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !bytes.Equal(data, validData) {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}
}

func TestSend(t *testing.T) {
	t.Run("Connect", testSendConnect)
	t.Run("Retry", testSendRetry)
	t.Run("Cancel", testSendCancel)
	t.Run("Rate", testSendRate)
	t.Run("Bad", testSendBad)
}

func makeServ(errors, rateLimit int, badRequest bool) (*httptest.Server, *url.URL) {
	serv := httptest.NewServer(&testHandler{errors, rateLimit, badRequest})
	u, _ := url.Parse(serv.URL)
	return serv, u
}

func testSendConnect(t *testing.T) {
	t.Skip("This test is not enough good to run always.")
	t.Parallel()

	serv := httptest.NewUnstartedServer(&testHandler{0, 0, false})
	m := &slackMessage{&url.URL{}, 2, validData}
	go func() {
		time.Sleep(1500 * time.Millisecond)
		serv.Start()
		u, _ := url.Parse(serv.URL)
		m.url = u
	}()
	if !m.send(context.Background()) {
		t.Error(`!m.send(context.Background())`)
	}
	serv.Close()
}

func testSendRetry(t *testing.T) {
	t.Parallel()

	serv, u := makeServ(0, 0, false)
	// invalid data
	m := &slackMessage{u, 0, nil}
	if m.send(context.Background()) {
		t.Error(`m.send(context.Background())`)
	}
	serv.Close()

	serv, u = makeServ(0, 0, false)
	m = &slackMessage{u, 0, validData}
	if !m.send(context.Background()) {
		t.Error(`!m.send(context.Background())`)
	}
	serv.Close()

	serv, u = makeServ(1, 0, false)
	m = &slackMessage{u, 0, validData}
	if m.send(context.Background()) {
		t.Error(`m.send(context.Background())`)
	}
	serv.Close()

	serv, u = makeServ(1, 0, false)
	m = &slackMessage{u, 1, validData}
	if !m.send(context.Background()) {
		t.Error(`!m.send(context.Background())`)
	}
	serv.Close()
}

func testSendCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	serv, u := makeServ(1, 0, false)
	m := &slackMessage{u, 2, validData}
	// should give up
	if m.send(ctx) {
		t.Error(`m.send(context.Background())`)
	}
	serv.Close()

	serv, u = makeServ(0, 0, false)
	m = &slackMessage{u, 0, validData}
	// should succeed
	if !m.send(ctx) {
		t.Error(`!m.send(context.Background())`)
	}
	serv.Close()
}

func testSendRate(t *testing.T) {
	t.Parallel()

	serv, u := makeServ(0, 2, false)
	m := &slackMessage{u, 0, validData}
	if m.send(context.Background()) {
		t.Error(`m.send(context.Background())`)
	}
	serv.Close()

	now := time.Now()
	serv, u = makeServ(0, 2, false)
	m = &slackMessage{u, 1, validData}
	if !m.send(context.Background()) {
		t.Error(`!m.send(context.Background())`)
	}
	serv.Close()

	if time.Now().Sub(now) < (2 * time.Second) {
		t.Error(`time.Now().Sub(now) < (2 * time.Second)`)
	}
}

func testSendBad(t *testing.T) {
	t.Parallel()

	serv, u := makeServ(0, 0, true)
	m := &slackMessage{u, 3, validData}
	if m.send(context.Background()) {
		t.Error(`m.send(context.Background())`)
	}
	serv.Close()
}
