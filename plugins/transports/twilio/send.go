package twilio

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
)

const (
	queueSize = 100

	// see https://www.twilio.com/docs/api/rest/sending-messages#rate-limiting
	twilioSMSInterval = 1100 * time.Millisecond

	twilioTimeout = 5 * time.Second
)

var (
	sendCh     = make(chan *twilioSMS, queueSize)
	httpClient = &well.HTTPClient{
		Client:   &http.Client{},
		Severity: log.LvDebug,
	}
)

func enqueue(m *twilioSMS) bool {
	select {
	case sendCh <- m:
		return true
	default:
		return false
	}
}

type twilioSMS struct {
	url      *url.URL
	username string
	password string
	maxRetry int
	payload  string
}

func (m *twilioSMS) wait(ctx context.Context, duration time.Duration) bool {
	select {
	case <-ctx.Done():
		return false
	case <-time.After(duration):
		return true
	}
}

func (m *twilioSMS) do(ctx context.Context) (*http.Response, []byte, error) {
	header := make(http.Header)
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	req := &http.Request{
		Method:        "POST",
		URL:           m.url,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
		Body:          ioutil.NopCloser(strings.NewReader(m.payload)),
		ContentLength: int64(len(m.payload)),
		Host:          m.url.Host,
	}
	req.SetBasicAuth(m.username, m.password)

	// use context.Background to send SMS gracefully.
	ctx, cancel := context.WithTimeout(context.Background(), twilioTimeout)
	defer cancel()

	resp, err := httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, nil, err
	}
	return resp, data, nil
}

func (m *twilioSMS) send(ctx context.Context) bool {
	var retries int

RETRY:
	resp, body, err := m.do(ctx)
	if err != nil {
		log.Error("[twilio] do", map[string]interface{}{
			log.FnError: err.Error(),
		})
		if retries < m.maxRetry {
			retries++
			if m.wait(ctx, twilioSMSInterval) {
				goto RETRY
			}
		}
		log.Error("[twilio] gave up", nil)
		return false
	}

	switch {
	case (200 <= resp.StatusCode) && (resp.StatusCode < 300):
		log.Info("[twilio] sent SMS", nil)
		return true

	case resp.StatusCode == 429, resp.StatusCode >= 500:
		// temporary server failure, hopefully.
		log.Error("[twilio] failed to send", map[string]interface{}{
			log.FnHTTPStatusCode: resp.StatusCode,
		})

	default:
		// mainly because the request was bad.
		fields := map[string]interface{}{
			log.FnHTTPStatusCode: resp.StatusCode,
		}
		if len(body) > 0 {
			var e map[string]interface{}
			err := json.Unmarshal(body, &e)
			if err == nil {
				fields["exception"] = e
			}
		}
		log.Error("[twilio] request failed", fields)
		return false
	}

	if retries < m.maxRetry {
		retries++
		if m.wait(ctx, twilioSMSInterval) {
			goto RETRY
		}
	}
	log.Error("[twilio] gave up", nil)
	return false
}

// A goroutine to send SMS complying its rate limits.
func dequeueAndSend(ctx context.Context) error {
	for {
		select {
		case m := <-sendCh:
			m.send(ctx)
			m.wait(ctx, twilioSMSInterval)
		case <-ctx.Done():
			return nil
		}
	}
}
