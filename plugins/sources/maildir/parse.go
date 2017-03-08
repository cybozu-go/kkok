package maildir

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"mime"
	"mime/quotedprintable"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/ianaindex"

	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/log"
	"github.com/pkg/errors"
)

const (
	maxMailSize = 10 * 1024 * 1024 // 10 MiB
)

var (
	mimeDecoder = mime.WordDecoder{CharsetReader: charsetReader}

	// for parseDate
	dateLayouts []string

	rePseudoHeader = regexp.MustCompile(`^([a-zA-Z0-9-]+):[ \t]*(.*)$`)
)

func charsetReader(cs string, r io.Reader) (io.Reader, error) {
	enc, err := ianaindex.MIME.Encoding(cs)

	// ianaindex returns nil for character sets like ks_c_5601-1987.
	// We fallback to htmlindex in such cases.
	if err == nil && enc == nil {
		enc, err = htmlindex.Get(cs)
	}

	if err != nil {
		return nil, err
	}

	// enc may be nil for ASCII or UTF-8 character sets.
	if enc != nil {
		r = enc.NewDecoder().Reader(r)
	}
	return r, nil
}

func init() {
	// Generate layouts based on RFC 5322, section 3.3.
	dows := [...]string{"", "Mon, "}   // day-of-week
	days := [...]string{"2", "02"}     // day = 1*2DIGIT
	years := [...]string{"2006", "06"} // year = 4*DIGIT / 2*DIGIT
	seconds := [...]string{":05", ""}  // second
	// "-0700 (MST)" is not in RFC 5322, but is common.
	zones := [...]string{"-0700", "MST", "-0700 (MST)"} // zone = (("+" / "-") 4DIGIT) / "GMT" / ...

	for _, dow := range dows {
		for _, day := range days {
			for _, year := range years {
				for _, second := range seconds {
					for _, zone := range zones {
						s := dow + day + " Jan " + year + " 15:04" + second + " " + zone
						dateLayouts = append(dateLayouts, s)
					}
				}
			}
		}
	}
}

// getFrom parses mail address line such as "Foo Bar <foobar@example.com>"
// and returns "Foo Bar".  If no name is given, this returns the mail address.
func getName(address string) string {
	a, err := mail.ParseAddress(address)
	if err != nil {
		return address
	}
	if len(a.Name) > 0 {
		return a.Name
	}
	return a.Address
}

// equivalent function for mail.ParseDate in Go 1.8+
func parseDate(s string) (time.Time, error) {
	for _, layout := range dateLayouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("mail: header could not be parsed")
}

func decodeBody(m *mail.Message) ([]byte, error) {
	r := m.Body
	switch t := strings.ToLower(m.Header.Get("Content-Transfer-Encoding")); t {
	case "base64":
		r = base64.NewDecoder(base64.StdEncoding, r)
	case "quoted-printable":
		r = quotedprintable.NewReader(r)
	case "", "7bit", "8bit":
	default:
		return nil, errors.New("unsupported transfer encoding: " + t)
	}

	ct := m.Header.Get("Content-Type")
	if len(ct) > 0 {
		mt, params, err := mime.ParseMediaType(ct)
		if err != nil {
			return nil, errors.Wrap(err, "bad content-type")
		}
		if !strings.HasPrefix(mt, "text/") {
			return nil, errors.New("non text content: " + mt)
		}
		if cs, ok := params["charset"]; ok {
			r, err = charsetReader(cs, r)
			if err != nil {
				return nil, errors.New("unsupported character set: " + cs)
			}
		}
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return bytes.Replace(data, []byte{'\r'}, nil, -1), nil
}

func parseBody(data []byte, a *kkok.Alert) {
	headers := make(map[string]string)

	buf := data[:]
	for {
		idx := bytes.IndexByte(buf, '\n')

		// no newlines are left
		if idx == -1 {
			if len(buf) > 0 {
				// incomplete pseudo headers should be ignored.
				headers = nil
				buf = data[:]
			}
			break
		}

		if idx == 0 {
			if len(headers) > 0 {
				// empty line separates body from (pseudo) headers.
				buf = buf[1:]
			}
			break
		}

		matches := rePseudoHeader.FindSubmatch(buf[:idx])
		if len(matches) != 3 {
			// discard headers as no pseudo headers shall be included
			headers = nil
			buf = data[:]
			break
		}

		headers[string(matches[1])] = string(matches[2])
		buf = buf[(idx + 1):]
	}

	for k, v := range headers {
		switch k {
		case "From":
			a.From = v
		case "Date":
			if dt, err := parseDate(v); err == nil {
				a.Date = dt
				break
			}
			if dt, err := time.Parse(time.RFC3339Nano, v); err == nil {
				a.Date = dt
				break
			}
			log.Warn("ignored ill-formatted date", map[string]interface{}{
				"value": v,
			})
		case "Host":
			a.Host = v
		case "Title":
			a.Title = v
		default:
			a.Info[k] = v
		}
	}
	a.Message = string(buf)
}

// parse read a mail source from r and generate an alert.
func parse(r io.Reader) (*kkok.Alert, error) {
	r = &io.LimitedReader{
		R: r,
		N: maxMailSize,
	}

	m, err := mail.ReadMessage(r)
	if err != nil {
		return nil, err
	}

	var dt time.Time
	if hDate := m.Header.Get("Date"); len(hDate) > 0 {
		dt2, err := parseDate(hDate)
		if err == nil {
			dt = dt2
		}
	}

	title, _ := mimeDecoder.DecodeHeader(m.Header.Get("Subject"))

	a := &kkok.Alert{
		From:  getName(m.Header.Get("From")),
		Date:  dt,
		Title: title,
		Info:  make(map[string]interface{}),
	}

	body, err := decodeBody(m)
	if err != nil {
		log.Error("failed to decode mail body", map[string]interface{}{
			log.FnError:    err.Error(),
			"mail_from":    a.From,
			"mail_subject": a.Title,
		})
		return nil, err
	}

	parseBody(body, a)

	err = a.Validate()
	if err != nil {
		return nil, err
	}
	return a, nil
}
