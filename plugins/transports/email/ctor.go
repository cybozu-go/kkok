package email

import (
	"text/template"

	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/kkok/util"
	"github.com/pkg/errors"
)

func ctor(params map[string]interface{}) (kkok.Transport, error) {
	from, err := util.GetString("from", params)
	if err != nil {
		return nil, errors.Wrap(err, "email: from")
	}

	tr := &transport{
		from: from,
	}

	label, err := util.GetString("label", params)
	switch {
	case err == nil:
		tr.label = label
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "email: label")
	}

	host, err := util.GetString("host", params)
	switch {
	case err == nil:
		tr.host = host
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "email: host")
	}

	port, err := util.GetInt("port", params)
	switch {
	case err == nil:
		tr.port = port
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "email: port")
	}

	user, err := util.GetString("user", params)
	switch {
	case err == nil:
		tr.username = user
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "email: user")
	}

	password, err := util.GetString("password", params)
	switch {
	case err == nil:
		tr.password = password
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "email: password")
	}

	to, err := util.GetStringSlice("to", params)
	switch {
	case err == nil:
		tr.to = to
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "email: to")
	}

	cc, err := util.GetStringSlice("cc", params)
	switch {
	case err == nil:
		tr.cc = cc
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "email: cc")
	}

	bcc, err := util.GetStringSlice("bcc", params)
	switch {
	case err == nil:
		tr.bcc = bcc
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "email: bcc")
	}

	toFile, err := util.GetString("to_file", params)
	switch {
	case err == nil:
		tr.toFile = toFile
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "email: to_file")
	}

	ccFile, err := util.GetString("cc_file", params)
	switch {
	case err == nil:
		tr.ccFile = ccFile
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "email: cc_file")
	}

	bccFile, err := util.GetString("bcc_file", params)
	switch {
	case err == nil:
		tr.bccFile = bccFile
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "email: bcc_file")
	}

	tmplPath, err := util.GetString("template", params)
	switch {
	case err == nil:
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			return nil, errors.Wrap(err, "email: template")
		}
		tr.tmpl = tmpl
		tr.tmplPath = tmplPath
	case util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "email: template")
	}

	return tr, nil
}

func init() {
	kkok.RegisterTransport(transportType, ctor)
}
