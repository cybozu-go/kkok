package twilio

import (
	"regexp"
	"text/template"

	"github.com/cybozu-go/cmd"
	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/kkok/util"
	"github.com/pkg/errors"
)

var validPhoneNumber = regexp.MustCompile(`^[0-9]{5,6}$|^\+[0-9]+$`)

func ctor(params map[string]interface{}) (kkok.Transport, error) {
	account, err := util.GetString("account", params)
	if err != nil {
		return nil, errors.Wrap(err, "twilio: account")
	}
	token, err := util.GetString("token", params)
	if err != nil {
		return nil, errors.Wrap(err, "twilio: token")
	}
	from, err := util.GetString("from", params)
	if err != nil {
		return nil, errors.Wrap(err, "twilio: from")
	}
	if !validPhoneNumber.MatchString(from) {
		return nil, errors.New("twilio: invalid phone number: " + from)
	}

	tr, err := newTransport(account)
	if err != nil {
		return nil, err
	}
	tr.token = token
	tr.from = from

	label, err := util.GetString("label", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "twilio: label")
	}
	tr.label = label

	sid, err := util.GetString("key_sid", params)
	switch {
	case err == nil:
		tr.sid = sid
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "twilio: key_sid")
	}

	to, err := util.GetStringSlice("to", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "twilio: to")
	}
	for _, n := range to {
		if !validPhoneNumber.MatchString(n) {
			return nil, errors.New("twilio: invalid phone number: " + n)
		}
	}
	tr.to = to

	toFile, err := util.GetString("to_file", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "twilio: to_file")
	}
	tr.toFile = toFile

	maxLength, err := util.GetInt("max_length", params)
	switch {
	case err == nil:
		if maxLength <= 0 || maxLength > 1600 {
			return nil, errors.New("twilio: invalid max_length value")
		}
		tr.maxLength = maxLength
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "twilio: max_length")
	}

	maxRetry, err := util.GetInt("max_retry", params)
	switch {
	case err == nil:
		tr.maxRetry = maxRetry
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "twilio: max_retry")
	}

	tmplPath, err := util.GetString("template", params)
	switch {
	case err == nil:
		tmpl, err := template.New("").ParseFiles(tmplPath)
		if err != nil {
			return nil, errors.Wrap(err, "twilio: template")
		}
		tr.tmpl = tmpl
		tr.tmplPath = tmplPath
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "twilio: template")
	}

	countOnly, err := util.GetBool("count_only", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "twilio: count_only")
	}
	tr.countOnly = countOnly

	return tr, nil
}

func init() {
	kkok.RegisterTransport(transportType, ctor)
	cmd.Go(dequeueAndSend)
}
