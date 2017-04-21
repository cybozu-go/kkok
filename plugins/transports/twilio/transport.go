package twilio

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"os"
	"text/template"

	"github.com/cybozu-go/kkok"
	"github.com/pkg/errors"
)

const (
	transportType = "twilio"

	defaultLength = 160
	defaultRetry  = 3

	countOnlyMessageOne = "There is an alert."
	countOnlyMessage    = "There are %d alerts."

	twilioEndPoint = "https://api.twilio.com/2010-04-01/Accounts/"
)

type transport struct {
	url       *url.URL
	label     string
	account   string
	sid       string
	token     string
	from      string
	to        []string
	toFile    string
	maxLength int
	maxRetry  int
	countOnly bool
	tmplPath  string
	tmpl      *template.Template
	enqueue   func(*twilioSMS) bool
}

func newTransport(account string) (*transport, error) {
	u, err := url.ParseRequestURI(twilioEndPoint + account + "/Messages.json")
	if err != nil {
		return nil, err
	}

	tr := &transport{
		url:       u,
		account:   account,
		sid:       account,
		maxLength: defaultLength,
		maxRetry:  defaultRetry,
		tmpl:      defaultTemplate,
		enqueue:   enqueue,
	}
	return tr, nil
}

func (t *transport) String() string {
	if len(t.label) > 0 {
		return t.label
	}

	return transportType
}

func (t *transport) Params() kkok.PluginParams {
	m := map[string]interface{}{
		"account":    t.account,
		"token":      t.token,
		"from":       t.from,
		"max_length": t.maxLength,
		"max_retry":  t.maxRetry,
	}

	if len(t.label) > 0 {
		m["label"] = t.label
	}
	if t.account != t.sid {
		m["key_sid"] = t.sid
	}
	if len(t.to) > 0 {
		m["to"] = t.to
	}
	if len(t.toFile) > 0 {
		m["to_file"] = t.toFile
	}
	if len(t.tmplPath) > 0 {
		m["template"] = t.tmplPath
	}
	if t.countOnly {
		m["count_only"] = t.countOnly
	}

	return kkok.PluginParams{
		Type:   transportType,
		Params: m,
	}
}

func (t *transport) recipients() ([]string, error) {
	if len(t.toFile) == 0 {
		return t.to, nil
	}

	to := t.to
	f, err := os.Open(t.toFile)
	if err != nil {
		if os.IsNotExist(err) {
			return to, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}
		to = append(to, text)
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}
	return to, nil
}

func (t *transport) send(to []string, msg string) error {
	if len(msg) > t.maxLength {
		umsg := []rune(msg)
		if len(umsg) > t.maxLength {
			msg = string(umsg[:t.maxLength])
		}
	}

	v := url.Values{}
	v.Set("From", t.from)
	v.Set("Body", msg)
	for _, r := range to {
		v.Set("To", r)
		m := &twilioSMS{
			url:      t.url,
			username: t.sid,
			password: t.token,
			maxRetry: t.maxRetry,
			payload:  v.Encode(),
		}

		if !t.enqueue(m) {
			return errors.New("twilio: queue is full")
		}
	}
	return nil
}

func (t *transport) Deliver(alerts []*kkok.Alert) error {
	to, err := t.recipients()
	if err != nil {
		return errors.Wrap(err, transportType)
	}

	if len(to) == 0 {
		return nil
	}

	if t.countOnly {
		if len(alerts) <= 1 {
			return t.send(to, countOnlyMessageOne)
		}
		return t.send(to, fmt.Sprintf(countOnlyMessage, len(alerts)))
	}

	buf := new(bytes.Buffer)
	for _, a := range alerts {
		err := t.tmpl.Execute(buf, a)
		if err != nil {
			return errors.Wrap(err, transportType)
		}
		err = t.send(to, buf.String())
		if err != nil {
			return errors.Wrap(err, transportType)
		}
		buf.Reset()
	}

	return nil
}
