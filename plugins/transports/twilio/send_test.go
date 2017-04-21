package twilio

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var (
	validData = []byte("valid data")
)

type testHandler struct {
	errors     int
	statusCode int
	username   string
	password   string
}

func newTestHandler() *testHandler {
	return &testHandler{
		statusCode: http.StatusOK,
		username:   "user",
		password:   "password",
	}
}

func newTestSMS(us string) *twilioSMS {
	u, _ := url.Parse(us)
	return &twilioSMS{
		url:      u,
		username: "user",
		password: "password",
		payload:  string(validData),
	}
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.errors > 0 {
		h.errors--
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if h.statusCode != 0 && h.statusCode != http.StatusOK {
		http.Error(w, "failed", h.statusCode)
		return
	}

	u, p, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "no basic auth", http.StatusBadRequest)
		return
	}

	if u != h.username {
		http.Error(w, "invalid user", http.StatusForbidden)
		return
	}

	if p != h.password {
		http.Error(w, "invalid password", http.StatusForbidden)
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

func testSendConnect(t *testing.T) {
	t.Skip("This test is not enough good to run always.")
	t.Parallel()

	serv := httptest.NewUnstartedServer(newTestHandler())
	m := newTestSMS(serv.URL)
	m.maxRetry = 2
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

	// invalid data
	serv := httptest.NewServer(newTestHandler())
	m := newTestSMS(serv.URL)
	m.payload = "invalid"
	m.maxRetry = 3
	if m.send(context.Background()) {
		t.Error(`m.send(context.Background())`)
	}
	serv.Close()

	serv = httptest.NewServer(newTestHandler())
	m = newTestSMS(serv.URL)
	if !m.send(context.Background()) {
		t.Error(`!m.send(context.Background())`)
	}
	serv.Close()

	h := newTestHandler()
	h.errors = 1
	serv = httptest.NewServer(h)
	m = newTestSMS(serv.URL)
	if m.send(context.Background()) {
		t.Error(`m.send(context.Background())`)
	}
	serv.Close()

	h = newTestHandler()
	h.errors = 1
	serv = httptest.NewServer(h)
	m = newTestSMS(serv.URL)
	m.maxRetry = 1
	if !m.send(context.Background()) {
		t.Error(`!m.send(context.Background())`)
	}
	serv.Close()
}

func testSendCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	h := newTestHandler()
	h.errors = 1
	serv := httptest.NewServer(h)
	m := newTestSMS(serv.URL)
	m.maxRetry = 2
	// should give up
	if m.send(ctx) {
		t.Error(`m.send(context.Background())`)
	}
	serv.Close()

	h = newTestHandler()
	serv = httptest.NewServer(h)
	m = newTestSMS(serv.URL)
	// should succeed
	if !m.send(ctx) {
		t.Error(`!m.send(context.Background())`)
	}
	serv.Close()
}

func testSendRate(t *testing.T) {
	t.Parallel()

	now := time.Now()
	h := newTestHandler()
	h.statusCode = http.StatusTooManyRequests
	serv := httptest.NewServer(h)
	m := newTestSMS(serv.URL)
	m.maxRetry = 2
	if m.send(context.Background()) {
		t.Error(`m.send(context.Background())`)
	}
	serv.Close()

	if time.Now().Sub(now) < (2 * twilioSMSInterval) {
		t.Error(`time.Now().Sub(now) < (2 * twilioSMSInterval)`)
	}

	if time.Now().Sub(now) > (3 * twilioSMSInterval) {
		t.Error(`time.Now().Sub(now) > (3 * twilioSMSInterval)`)
	}
}

func testSendBad(t *testing.T) {
	t.Parallel()

	h := newTestHandler()
	h.statusCode = http.StatusBadRequest
	serv := httptest.NewServer(h)
	m := newTestSMS(serv.URL)
	m.maxRetry = 3
	now := time.Now()
	if m.send(context.Background()) {
		t.Error(`m.send(context.Background())`)
	}
	serv.Close()
	if time.Now().Sub(now) > twilioSMSInterval {
		t.Error(`time.Now().Sub(now) > twilioSMSInterval`)
	}

	h = newTestHandler()
	serv = httptest.NewServer(h)
	m = newTestSMS(serv.URL)
	m.username = "baduser"
	m.maxRetry = 3
	now = time.Now()
	if m.send(context.Background()) {
		t.Error(`m.send(context.Background())`)
	}
	serv.Close()
	if time.Now().Sub(now) > twilioSMSInterval {
		t.Error(`time.Now().Sub(now) > twilioSMSInterval`)
	}

	h = newTestHandler()
	serv = httptest.NewServer(h)
	m = newTestSMS(serv.URL)
	m.password = "badpass"
	m.maxRetry = 3
	now = time.Now()
	if m.send(context.Background()) {
		t.Error(`m.send(context.Background())`)
	}
	serv.Close()
	if time.Now().Sub(now) > twilioSMSInterval {
		t.Error(`time.Now().Sub(now) > twilioSMSInterval`)
	}
}
