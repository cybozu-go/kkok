package email

import (
	"net/mail"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/cybozu-go/kkok"
)

var (
	testHost string
	testPort int
)

func init() {
	testHost = os.Getenv("TEST_MAILHOST")
	sport := os.Getenv("TEST_MAILPORT")
	if len(sport) == 0 {
		testPort = 8025
		return
	}

	i, err := strconv.Atoi(sport)
	if err != nil {
		panic(err)
	}
	testPort = i
}

func TestTransport(t *testing.T) {
	t.Run("Params", testParams)
	t.Run("String", testString)
	t.Run("AddressList", testAddressList)
	t.Run("Compose", testCompose)
	t.Run("Deliver", testDeliver)
}

func testParams(t *testing.T) {
	t.Parallel()

	tr := &transport{
		from: "kkok@example.com",
	}
	pp := tr.Params()

	if pp.Type != transportType {
		t.Error(`tr.Type != transportType`)
	}
	if len(pp.Params) != 1 {
		t.Error(`len(pp.Params) != 1`)
	}
	if pp.Params["from"] != "kkok@example.com" {
		t.Error(`pp.Params["from"] != "kkok@example.com"`)
	}

	tr = &transport{
		label:    "foo",
		host:     "h",
		port:     1025,
		username: "test",
		password: "secret",
		from:     "kkok@example.com",
		to:       []string{"to@example.org"},
		cc:       []string{"cc@example.org"},
		bcc:      []string{"bcc@example.org"},
		toFile:   "/path/to/to_file",
		ccFile:   "/path/to/cc_file",
		bccFile:  "/path/to/bcc_file",
		tmplPath: "/path/to/template_fille",
	}
	pp = tr.Params()
	m := pp.Params

	if m["label"].(string) != "foo" {
		t.Error(`m["label"].(string) != "foo"`)
	}
	if m["host"].(string) != "h" {
		t.Error(`m["host"].(string) != "h"`)
	}
	if m["port"].(int) != 1025 {
		t.Error(`m["port"].(int) != 1025`)
	}
	if m["user"] != "test" {
		t.Error(`m["user"] != "test"`)
	}
	if m["password"] != "secret" {
		t.Error(`m["password"] != "secret"`)
	}
	if m["from"] != "kkok@example.com" {
		t.Error(`m["from"] != "kkok@example.com"`)
	}
	if !reflect.DeepEqual(m["to"], []string{"to@example.org"}) {
		t.Error(`!reflect.DeepEqual(m["to"], []string{"to@example.org"})`)
	}
	if !reflect.DeepEqual(m["cc"], []string{"cc@example.org"}) {
		t.Error(`!reflect.DeepEqual(m["cc"], []string{"cc@example.org"})`)
	}
	if !reflect.DeepEqual(m["bcc"], []string{"bcc@example.org"}) {
		t.Error(`!reflect.DeepEqual(m["bcc"], []string{"bcc@example.org"})`)
	}
	if m["to_file"].(string) != "/path/to/to_file" {
		t.Error(`m["to_file"] != "/path/to/to_file"`)
	}
	if m["cc_file"].(string) != "/path/to/cc_file" {
		t.Error(`m["cc_file"].(string) != "/path/to/cc_file"`)
	}
	if m["bcc_file"].(string) != "/path/to/bcc_file" {
		t.Error(`m["bcc_file"].(string) != "/path/to/bcc_file"`)
	}
	if m["template"].(string) != "/path/to/template_fille" {
		t.Error(`m["template"].(string) != "/path/to/template_fille"`)
	}
}

func testString(t *testing.T) {
	t.Parallel()

	tr := &transport{}
	if tr.String() != "email" {
		t.Error(`tr.String() != "email"`)
	}

	tr = &transport{
		label: "test",
	}
	if tr.String() != "test" {
		t.Error(`tr.String() != "test"`)
	}
}

func testAddressList(t *testing.T) {
	t.Run("NoFile", testAddressListNoFile)
	t.Run("NotFound", testAddressListNotFound)
	t.Run("File", testAddressListFile)
}

func testAddressListNoFile(t *testing.T) {
	al := []string{"abc", "def"}
	result, err := getAddressList(al, "")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, []string{"abc", "def"}) {
		t.Error(`!reflect.DeepEqual(result, []string{"abc", "def"})`)
	}
}

func testAddressListNotFound(t *testing.T) {
	_, err := getAddressList([]string{}, "/not/found/file")
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testAddressListFile(t *testing.T) {
	result, err := getAddressList([]string{"abc"}, "testdata/1.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, []string{"abc", "def", "ghi"}) {
		t.Error(`!reflect.DeepEqual(result, []string{"abc", "def", "ghi"})`)
	}

	result, err = getAddressList(nil, "testdata/1.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, []string{"def", "ghi"}) {
		t.Error(`!reflect.DeepEqual(result, []string{"def", "ghi"})`)
	}

	result, err = getAddressList(nil, "testdata/2.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, []string{"111", "222", "333"}) {
		t.Error(`!reflect.DeepEqual(result, []string{"111", "222", "333"})`)
	}
}

func testComposeOne(t *testing.T, to, cc, bcc []string) {
	t.Parallel()

	tr := &transport{
		from: "foo@example.com",
	}
	a := &kkok.Alert{
		From:    "test monitor",
		Date:    time.Date(2011, 2, 3, 4, 5, 6, 0, time.UTC),
		Host:    "host1",
		Title:   "test test",
		Message: "こんにちは\n",
	}

	m, err := tr.compose(a, to, cc, bcc)
	if err != nil {
		t.Fatal(err)
	}

	hFrom := m.GetHeader("From")
	if len(hFrom) != 1 {
		t.Fatal(`len(hFrom) != 1`)
	}
	addr, err := mail.ParseAddress(hFrom[0])
	if err != nil {
		t.Fatal(err)
	}
	if addr.Name != a.From {
		t.Error(`addr.Name != a.From`)
	}
	if addr.Address != tr.from {
		t.Error(`addr.Address != tr.from`)
	}

	if len(m.GetHeader("To")) != len(to) {
		t.Error(`len(m.GetHeader("To")) != len(to)`)
	}
	if len(m.GetHeader("Cc")) != len(cc) {
		t.Error(`len(m.GetHeader("Cc")) != len(cc)`)
	}
	if len(m.GetHeader("Bcc")) != len(bcc) {
		t.Error(`len(m.GetHeader("Bcc")) != len(bcc)`)
	}
	hSubject := m.GetHeader("Subject")
	if len(hSubject) != 1 {
		t.Fatal(`len(hSubject) != 1`)
	}
	if hSubject[0] != a.Title {
		t.Error(`hSubject[0] != a.Title`)
	}
	hDate := m.GetHeader("Date")
	if len(hDate) != 1 {
		t.Fatal(`len(hDate) != 1`)
	}
	if hDate[0] != "Thu, 03 Feb 2011 04:05:06 +0000" {
		t.Error(`hDate[0] != "Thu, 03 Feb 2011 04:05:06 +0000"`)
	}
	hMailer := m.GetHeader("X-Mailer")
	if len(hMailer) != 1 {
		t.Fatal(`len(hMailer) != 1`)
	}
}

func testCompose(t *testing.T) {
	t.Parallel()

	t.Run("To", func(t *testing.T) {
		testComposeOne(t, []string{"foo"}, nil, nil)
	})
	t.Run("Cc", func(t *testing.T) {
		testComposeOne(t, nil, []string{"foo", "bar"}, nil)
	})
	t.Run("Bcc", func(t *testing.T) {
		testComposeOne(t, nil, nil, []string{"foo", "bar", "zot"})
	})
}

func testDeliver(t *testing.T) {
	if len(testHost) == 0 {
		t.Skip("No TEST_MAILHOST envvar")
	}
	t.Parallel()

	tr := &transport{
		from: "foo@example.com",
		bcc:  []string{"bar@example.org", "zot@example.org"},
		host: testHost,
		port: testPort,
	}
	err := tr.Deliver([]*kkok.Alert{
		{
			From:    "from1",
			Date:    time.Date(2014, 03, 02, 11, 22, 33, 0, time.UTC),
			Title:   "タイトル",
			Host:    "host1",
			Message: "こんにちは",
		},
		{
			From:    "from2",
			Date:    time.Date(2014, 03, 02, 11, 22, 33, 123456789, time.UTC),
			Title:   "title2",
			Host:    "host2",
			Message: "世界",
		},
	})
	if err != nil {
		t.Error(err)
	}
}
