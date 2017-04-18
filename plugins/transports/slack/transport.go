package slack

import (
	"bytes"
	"encoding/json"
	"net/url"
	"text/template"

	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/log"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

const (
	transportType = "slack"

	// Slack's recommendation is 20.  Hard maximum is 100.
	maxAttachments = 20

	defaultRetry = 3

	rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"
)

type transport struct {
	url       *url.URL
	label     string
	maxRetry  int
	name      string
	icon      string
	channel   string
	color     *otto.Script
	origColor string
	tmplPath  string
	tmpl      *template.Template
	enqueue   func(*slackMessage) bool
}

func (t *transport) String() string {
	if len(t.label) > 0 {
		return t.label
	}

	return transportType
}

func (t *transport) Params() kkok.PluginParams {
	m := map[string]interface{}{
		"url":       t.url.String(),
		"max_retry": t.maxRetry,
	}

	if len(t.label) > 0 {
		m["label"] = t.label
	}
	if len(t.name) > 0 {
		m["name"] = t.name
	}
	if len(t.icon) > 0 {
		m["icon"] = t.icon
	}
	if len(t.channel) > 0 {
		m["channel"] = t.channel
	}
	if len(t.origColor) > 0 {
		m["color"] = t.origColor
	}
	if len(t.tmplPath) > 0 {
		m["template"] = t.tmplPath
	}

	return kkok.PluginParams{
		Type:   transportType,
		Params: m,
	}
}

func (t *transport) format(a *kkok.Alert) (*attachment, error) {
	at := &attachment{
		Fallback: a.Title,
		Title:    EscapeSlack(a.Title),
	}

	if t.color != nil {
		v, err := a.Eval(t.color)
		if err != nil {
			return nil, errors.Wrap(err, t.String())
		}
		switch {
		case v.IsString():
			at.Color = v.String()
		case v.IsNull(), v.IsUndefined():
		default:
			return nil, errors.New("slack: non-string color")
		}
	}

	tmpl := t.tmpl
	if tmpl == nil {
		tmpl = defaultTemplate
	}
	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, a)
	if err != nil {
		return nil, errors.Wrap(err, t.String())
	}
	at.Text = buf.String()

	at.addField("From", a.From)
	if !a.Date.IsZero() {
		at.addField("Date", a.Date.Format(rfc3339Milli))
	}
	if len(a.Host) > 0 {
		at.addField("Host", a.Host)
	}
	for k, v := range a.Info {
		at.addField("Info."+k, v)
	}

	return at, nil
}

func (t *transport) send(alerts []*kkok.Alert) error {
	m := &message{
		Name:    t.name,
		Icon:    t.icon,
		Channel: t.channel,
	}

	for _, a := range alerts {
		at, err := t.format(a)
		if err != nil {
			return err
		}
		m.Attachments = append(m.Attachments, at)
	}

	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	eq := t.enqueue
	if eq == nil {
		eq = enqueue
	}

	ok := eq(&slackMessage{
		url:      t.url,
		maxRetry: t.maxRetry,
		payload:  data,
	})
	if !ok {
		return errors.New("slack queue is full")
	}

	return nil
}

func (t *transport) Deliver(alerts []*kkok.Alert) error {
	for i := 0; i < len(alerts); i += maxAttachments {
		pos := i + maxAttachments
		if pos > len(alerts) {
			pos = len(alerts)
		}

		err := t.send(alerts[i:pos])
		if err == nil {
			continue
		}

		fields := map[string]interface{}{
			log.FnError: err.Error(),
		}
		if len(t.label) > 0 {
			fields["label"] = t.label
		}
		log.Error("[slack] failed to enqueue", fields)
	}

	return nil
}
