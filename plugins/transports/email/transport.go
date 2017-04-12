package email

import (
	"bufio"
	"bytes"
	"os"
	"text/template"

	gomail "gopkg.in/gomail.v2"

	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/log"
	"github.com/pkg/errors"
)

const (
	transportType = "email"

	defaultHost = "localhost"
	defaultPort = 25

	mailer = "kkok " + kkok.Version
)

type transport struct {
	label    string
	host     string
	port     int
	username string
	password string
	from     string
	to       []string
	cc       []string
	bcc      []string
	toFile   string
	ccFile   string
	bccFile  string
	tmplPath string
	tmpl     *template.Template
}

func (t *transport) String() string {
	if len(t.label) == 0 {
		return "email"
	}
	return t.label
}

func (t *transport) Params() kkok.PluginParams {
	m := map[string]interface{}{
		"from": t.from,
	}
	if len(t.label) != 0 {
		m["label"] = t.label
	}
	if len(t.host) != 0 {
		m["host"] = t.host
	}
	if t.port != 0 {
		m["port"] = t.port
	}
	if len(t.username) != 0 {
		m["user"] = t.username
	}
	if len(t.password) != 0 {
		m["password"] = t.password
	}
	if len(t.to) > 0 {
		m["to"] = t.to
	}
	if len(t.cc) > 0 {
		m["cc"] = t.cc
	}
	if len(t.bcc) > 0 {
		m["bcc"] = t.bcc
	}
	if len(t.toFile) != 0 {
		m["to_file"] = t.toFile
	}
	if len(t.ccFile) != 0 {
		m["cc_file"] = t.ccFile
	}
	if len(t.bccFile) != 0 {
		m["bcc_file"] = t.bccFile
	}
	if len(t.tmplPath) != 0 {
		m["template"] = t.tmplPath
	}

	return kkok.PluginParams{
		Type:   transportType,
		Params: m,
	}
}

func getAddressList(s []string, filename string) ([]string, error) {
	if len(filename) == 0 {
		return s, nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 {
			continue
		}
		s = append(s, t)
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (t *transport) compose(alert *kkok.Alert, to, cc, bcc []string) (*gomail.Message, error) {
	tmpl := t.tmpl
	if tmpl == nil {
		tmpl = defaultTemplate
	}

	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, alert)
	if err != nil {
		return nil, errors.Wrap(err, transportType)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", t.from, alert.From)
	if len(to) > 0 {
		m.SetHeader("To", to...)
	}
	if len(cc) > 0 {
		m.SetHeader("Cc", cc...)
	}
	if len(bcc) > 0 {
		m.SetHeader("Bcc", bcc...)
	}
	m.SetHeader("Subject", alert.Title)
	m.SetDateHeader("Date", alert.Date)
	m.SetHeader("X-Mailer", mailer)
	m.SetBody("text/plain", buf.String())
	return m, nil
}

func (t *transport) Deliver(alerts []*kkok.Alert) error {
	to, err := getAddressList(t.to, t.toFile)
	if err != nil {
		return errors.Wrap(err, transportType)
	}
	cc, err := getAddressList(t.cc, t.ccFile)
	if err != nil {
		return errors.Wrap(err, transportType)
	}
	bcc, err := getAddressList(t.bcc, t.bccFile)
	if err != nil {
		return errors.Wrap(err, transportType)
	}

	if len(to)+len(cc)+len(bcc) == 0 {
		return nil
	}

	host := t.host
	if len(host) == 0 {
		host = defaultHost
	}
	port := t.port
	if port == 0 {
		port = defaultPort
	}
	s, err := gomail.NewDialer(host, port, t.username, t.password).Dial()
	if err != nil {
		return errors.Wrap(err, transportType)
	}
	defer s.Close()

	for _, a := range alerts {
		m, err := t.compose(a, to, cc, bcc)
		if err != nil {
			return err
		}

		err = gomail.Send(s, m)
		fields := map[string]interface{}{
			"transport": transportType,
			"from":      a.From,
			"title":     a.Title,
			"host":      a.Host,
		}
		if len(to) > 0 {
			fields["to"] = to
		}
		if len(cc) > 0 {
			fields["cc"] = cc
		}
		if len(bcc) > 0 {
			fields["bcc"] = bcc
		}
		if err == nil {
			log.Info("sent a mail", fields)
		} else {
			fields[log.FnError] = err.Error()
			log.Error("failed to send a mail", fields)
		}
	}

	return nil
}
