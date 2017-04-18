package slack

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/cybozu-go/cmd"
	"github.com/cybozu-go/log"
)

const (
	queueSize = 100

	// See https://api.slack.com/docs/rate-limits
	slackAPIDuration = 1100 * time.Millisecond

	slackAPITimeout = 5 * time.Second
)

var (
	sendCh     = make(chan *slackMessage, queueSize)
	httpClient = &cmd.HTTPClient{
		Client:   &http.Client{},
		Severity: log.LvDebug,
	}
)

func enqueue(m *slackMessage) bool {
	select {
	case sendCh <- m:
		return true
	default:
		return false
	}
}

type slackMessage struct {
	url      *url.URL
	maxRetry int
	payload  []byte
}

func (m *slackMessage) wait(ctx context.Context, duration time.Duration) bool {
	select {
	case <-ctx.Done():
		return false
	case <-time.After(duration):
		return true
	}
}

func (m *slackMessage) do(ctx context.Context) (*http.Response, error) {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	req := &http.Request{
		Method:        "POST",
		URL:           m.url,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
		Body:          ioutil.NopCloser(bytes.NewReader(m.payload)),
		ContentLength: int64(len(m.payload)),
		Host:          m.url.Host,
	}

	// use context.Background to send alerts gracefully.
	ctx, cancel := context.WithTimeout(context.Background(), slackAPITimeout)
	defer cancel()
	return httpClient.Do(req.WithContext(ctx))
}

func (m *slackMessage) send(ctx context.Context) bool {
	var retries int

RETRY:
	resp, err := m.do(ctx)
	if err != nil {
		log.Error("[slack] do", map[string]interface{}{
			log.FnError: err.Error(),
			log.FnURL:   m.url.String(),
		})
		if retries < m.maxRetry {
			retries++
			if m.wait(ctx, slackAPIDuration) {
				goto RETRY
			}
		}
		log.Error("[slack] gave up", nil)
		return false
	}
	defer func(resp *http.Response) {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}(resp)

	sleepDuration := slackAPIDuration

	switch {
	case (200 <= resp.StatusCode) && (resp.StatusCode < 300):
		log.Info("[slack] sent alerts", nil)
		return true

	case resp.StatusCode == 429:
		// rate limit exceeded
		log.Warn("[slack] rate limit exceeds", nil)
		ssec := resp.Header.Get("Retry-After")
		if len(ssec) > 0 {
			sec, err := strconv.Atoi(ssec)
			if err == nil {
				sleepDuration = time.Duration(sec) * time.Second
			}
		}

	case resp.StatusCode >= 500:
		// temporary server failure, hopefully.
		log.Error("[slack] failed to send", map[string]interface{}{
			log.FnURL:            m.url.String(),
			log.FnHTTPStatusCode: resp.StatusCode,
		})

	default:
		// mainly because the request was bad.
		fields := map[string]interface{}{
			log.FnURL:            m.url,
			log.FnHTTPStatusCode: resp.StatusCode,
		}
		data, _ := ioutil.ReadAll(resp.Body)
		if len(data) > 0 {
			fields[log.FnError] = string(data)
		}
		log.Error("[slack] request failed", fields)
		return false
	}

	if retries < m.maxRetry {
		retries++
		if m.wait(ctx, sleepDuration) {
			goto RETRY
		}
	}
	log.Error("[slack] gave up", nil)
	return false
}

// A goroutine to send messages to slack complying its rate limits.
func dequeueAndSend(ctx context.Context) error {
	for {
		select {
		case m := <-sendCh:
			m.send(ctx)
			m.wait(ctx, slackAPIDuration)
		case <-ctx.Done():
			return nil
		}
	}
}
