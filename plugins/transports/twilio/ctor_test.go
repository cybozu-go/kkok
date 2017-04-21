package twilio

import (
	"reflect"
	"testing"
)

type ctorTest struct {
	params map[string]interface{}
	tr     *transport
}

var ctorTestData = map[string]ctorTest{
	"tr1": {nil, nil},
	"tr2": {map[string]interface{}{
		"account": "account",
		"token":   "token",
	}, nil},
	"tr3": {map[string]interface{}{
		"account": "account",
		"from":    "+123456789",
	}, nil},
	"tr4": {map[string]interface{}{
		"token": "token",
		"from":  "+123456789",
	}, nil},
	"tr5": {map[string]interface{}{
		"account": "account",
		"token":   "token",
		"from":    "+123456789",
	}, &transport{
		account:   "account",
		sid:       "account",
		token:     "token",
		from:      "+123456789",
		maxLength: defaultLength,
		maxRetry:  defaultRetry,
	}},
	"tr6": {map[string]interface{}{
		"account": "account",
		"token":   "token",
		"from":    "12345",
	}, &transport{
		account:   "account",
		sid:       "account",
		token:     "token",
		from:      "12345",
		maxLength: defaultLength,
		maxRetry:  defaultRetry,
	}},
	"tr7": {map[string]interface{}{
		"account": "account",
		"token":   "token",
		"from":    "123456",
	}, &transport{
		account:   "account",
		sid:       "account",
		token:     "token",
		from:      "123456",
		maxLength: defaultLength,
		maxRetry:  defaultRetry,
	}},
	"tr8": {map[string]interface{}{
		"account": "account",
		"token":   "token",
		"from":    "1234",
	}, nil},
	"tr9": {map[string]interface{}{
		"account": "account",
		"token":   "token",
		"from":    "12345678",
	}, nil},
	"tr10": {map[string]interface{}{
		"account": "account",
		"token":   "token",
		"from":    "abc+12345678",
	}, nil},
	"tr11": {map[string]interface{}{
		"account": "account",
		"token":   "token",
		"from":    "+12345678abc",
	}, nil},
	"tr12": {map[string]interface{}{
		"account": "account",
		"key_sid": "sid",
		"token":   "token",
		"from":    "+123456789",
	}, &transport{
		account:   "account",
		sid:       "sid",
		token:     "token",
		from:      "+123456789",
		maxLength: defaultLength,
		maxRetry:  defaultRetry,
	}},
	"tr13": {map[string]interface{}{
		"account": "account",
		"key_sid": true,
		"token":   "token",
		"from":    "+123456789",
	}, nil},
	"tr14": {map[string]interface{}{
		"account": "account",
		"token":   "token",
		"from":    "+123456789",
		"to":      []string{"+123456789", "123456"},
	}, &transport{
		account:   "account",
		sid:       "account",
		token:     "token",
		from:      "+123456789",
		to:        []string{"+123456789", "123456"},
		maxLength: defaultLength,
		maxRetry:  defaultRetry,
	}},
	"tr15": {map[string]interface{}{
		"account": "account",
		"token":   "token",
		"from":    "+123456789",
		"to":      []string{"abc"},
	}, nil},
	"tr16": {map[string]interface{}{
		"account": "account",
		"token":   "token",
		"from":    "+123456789",
		"to_file": "/path/to/file",
	}, &transport{
		account:   "account",
		sid:       "account",
		token:     "token",
		from:      "+123456789",
		toFile:    "/path/to/file",
		maxLength: defaultLength,
		maxRetry:  defaultRetry,
	}},
	"tr17": {map[string]interface{}{
		"account":    "account",
		"token":      "token",
		"from":       "+123456789",
		"max_length": 20,
	}, &transport{
		account:   "account",
		sid:       "account",
		token:     "token",
		from:      "+123456789",
		maxLength: 20,
		maxRetry:  defaultRetry,
	}},
	"tr18": {map[string]interface{}{
		"account":    "account",
		"token":      "token",
		"from":       "+123456789",
		"max_length": 0,
	}, nil},
	"tr19": {map[string]interface{}{
		"account":    "account",
		"token":      "token",
		"from":       "+123456789",
		"max_length": -3,
	}, nil},
	"tr20": {map[string]interface{}{
		"account":    "account",
		"token":      "token",
		"from":       "+123456789",
		"max_length": 2000,
	}, nil},
	"tr21": {map[string]interface{}{
		"account":   "account",
		"token":     "token",
		"from":      "+123456789",
		"max_retry": 0,
	}, &transport{
		account:   "account",
		sid:       "account",
		token:     "token",
		from:      "+123456789",
		maxLength: defaultLength,
		maxRetry:  0,
	}},
	"tr22": {map[string]interface{}{
		"account":  "account",
		"token":    "token",
		"from":     "+123456789",
		"template": "testdata/test.tmpl",
	}, &transport{
		account:   "account",
		sid:       "account",
		token:     "token",
		from:      "+123456789",
		maxLength: defaultLength,
		maxRetry:  defaultRetry,
		tmplPath:  "testdata/test.tmpl",
	}},
	"tr23": {map[string]interface{}{
		"account":  "account",
		"token":    "token",
		"from":     "+123456789",
		"template": "testdata/invalid.tmpl",
	}, nil},
	"tr24": {map[string]interface{}{
		"account":    "account",
		"token":      "token",
		"from":       "+123456789",
		"count_only": true,
	}, &transport{
		account:   "account",
		sid:       "account",
		token:     "token",
		from:      "+123456789",
		maxLength: defaultLength,
		maxRetry:  defaultRetry,
		countOnly: true,
	}},
	"tr25": {map[string]interface{}{
		"label":   "abc",
		"account": "account",
		"token":   "token",
		"from":    "+123456789",
	}, &transport{
		label:     "abc",
		account:   "account",
		sid:       "account",
		token:     "token",
		from:      "+123456789",
		maxLength: defaultLength,
		maxRetry:  defaultRetry,
	}},
}

func testCtorOne(t *testing.T, data ctorTest) {
	t.Parallel()

	tr, err := ctor(data.params)
	if err != nil {
		if data.tr != nil {
			// unexpected
			t.Error(err)
		}
		return
	}

	tr2, ok := tr.(*transport)
	if !ok {
		t.Fatal("not a proper transport")
	}
	if tr2.url == nil {
		t.Error(`tr2.url == nil`)
	} else {
		tr2.url = nil
	}
	tr2.tmpl = nil
	tr2.enqueue = nil
	if !reflect.DeepEqual(tr2, data.tr) {
		t.Error(`!reflect.DeepEqual(tr2, data.tr)`)
		t.Logf("%#v, %#v", tr2, data.tr)
	}
}

func TestCtor(t *testing.T) {
	for id, data := range ctorTestData {
		t.Run(id, func(t *testing.T) { testCtorOne(t, data) })
	}
}
