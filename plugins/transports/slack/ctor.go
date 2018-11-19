package slack

import (
	"net/url"

	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/kkok/util"
	"github.com/cybozu-go/well"
	"github.com/pkg/errors"
)

func ctor(params map[string]interface{}) (kkok.Transport, error) {
	us, err := util.GetString("url", params)
	if err != nil {
		return nil, errors.Wrap(err, "slack: url")
	}

	u, err := url.ParseRequestURI(us)
	if err != nil {
		return nil, errors.Wrap(err, "slack: url")
	}

	tr := &transport{
		url:      u,
		maxRetry: defaultRetry,
	}

	label, err := util.GetString("label", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "slack: label")
	}
	tr.label = label

	maxRetry, err := util.GetInt("max_retry", params)
	switch {
	case err == nil:
		tr.maxRetry = maxRetry
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "slack: max_retry")
	}

	name, err := util.GetString("name", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "slack: name")
	}
	tr.name = name

	icon, err := util.GetString("icon", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "slack: icon")
	}
	tr.icon = icon

	channel, err := util.GetString("channel", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "slack: channel")
	}
	tr.channel = channel

	color, err := util.GetString("color", params)
	switch {
	case err == nil:
		cs, err := kkok.CompileJS(color)
		if err != nil {
			return nil, errors.Wrap(err, "slack: color")
		}
		tr.color = cs
		tr.origColor = color
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "slack: color")
	}

	tmplPath, err := util.GetString("template", params)
	switch {
	case err == nil:
		tmpl, err := newTemplate().ParseFiles(tmplPath)
		if err != nil {
			return nil, errors.Wrap(err, "slack: template")
		}
		tr.tmpl = tmpl
		tr.tmplPath = tmplPath
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "slack: template")
	}

	return tr, nil
}

func init() {
	kkok.RegisterTransport(transportType, ctor)
	well.Go(dequeueAndSend)
}
