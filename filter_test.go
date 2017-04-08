package kkok

import (
	"os/exec"
	"testing"
	"time"
)

func newBaseFilter(id string, params map[string]interface{}) (*BaseFilter, error) {
	b := new(BaseFilter)
	err := b.Init(id, params)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func testBaseFilterAll(t *testing.T) {
	t.Parallel()
	params := map[string]interface{}{
		"label":    "テスト",
		"disabled": true,
		"all":      true,
		"if":       "alerts.length > 1",
	}

	b, err := newBaseFilter("base", params)
	if err != nil {
		t.Fatal(err)
	}

	if b.ID() != "base" {
		t.Error(`b.ID() != "base"`)
	}
	if b.Label() != "テスト" {
		t.Error(`b.Label() != "テスト"`)
	}
	if b.Dynamic() {
		t.Error(`b.Dynamic()`)
	}
	if !b.Disabled() {
		t.Error(`!b.Disabled()`)
	}
	b.Enable(true)
	if b.Disabled() {
		t.Error(`b.Disabled()`)
	}
	if !b.All() {
		t.Error(`!b.All()`)
	}

	ok, err := b.EvalAllAlerts(nil)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("condition should not be met")
	}

	ok, err = b.EvalAllAlerts([]*Alert{{}, {}})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("condition should be met")
	}
}

func testBaseFilterOne(t *testing.T) {
	t.Parallel()
	params := map[string]interface{}{
		"all": false,
		"if":  "alert.From == 'hoge'",
	}

	b, err := newBaseFilter("base", params)
	if err != nil {
		t.Fatal(err)
	}

	if b.Disabled() {
		t.Error(`b.Disabled()`)
	}
	if b.All() {
		t.Error(`b.All()`)
	}

	ok, err := b.EvalAlert(&Alert{})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("condition should not be met")
	}

	ok, err = b.EvalAlert(&Alert{From: "hoge"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("condition should be met")
	}

	params = map[string]interface{}{
		"all": false,
		"if":  "alert.Info.Hoge == 'fuga'",
	}
	b2, err := newBaseFilter("base", params)
	if err != nil {
		t.Fatal(err)
	}
	ok, err = b2.EvalAlert(&Alert{})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("condition should not be met")
	}

	params = map[string]interface{}{
		"all": false,
		"if":  "!alert.Info.Hoge",
	}
	b3, err := newBaseFilter("base", params)
	if err != nil {
		t.Fatal(err)
	}
	ok, err = b3.EvalAlert(&Alert{})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("condition should be met")
	}
}

func testBaseFilterParseError(t *testing.T) {
	t.Parallel()
	params := map[string]interface{}{
		"all": false,
		"if":  "alert.From =",
	}

	_, err := newBaseFilter("id", params)
	if err == nil {
		t.Fatal("if must cause a parse error")
	}
	t.Log(err)
}

func testBaseFilterCommand(t *testing.T) {
	jq, err := exec.LookPath("jq")
	if err != nil {
		t.Skip(err)
	}

	t.Parallel()
	params := map[string]interface{}{
		"all": false,
		"if":  []interface{}{jq, "-e", `.From == "hoge"`},
	}

	b, err := newBaseFilter("base", params)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := b.EvalAlert(&Alert{})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("condition should not be met")
	}

	ok, err = b.EvalAlert(&Alert{From: "hoge"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("condition should be met")
	}
}

func testBaseFilterExpire(t *testing.T) {
	t.Parallel()

	params := map[string]interface{}{
		"expire": "hoge",
	}

	_, err := newBaseFilter("id", params)
	if err == nil {
		t.Fatal("expire must cause a parse error")
	}
	t.Log(err)

	now := time.Now().UTC()

	params["expire"] = now.Add(-1 * time.Hour).Format(time.RFC3339)
	f, err := newBaseFilter("id", params)
	if err != nil {
		t.Fatal(err)
	}

	if f.expire.IsZero() {
		t.Error(`f.expire.IsZero()`)
	}

	// static filters never expires.
	if f.Expired() {
		t.Error(`f.Expired()`)
	}
	// set it dynamic, then it gets expired.
	f.SetDynamic()
	if !f.Expired() {
		t.Error(`!f.Expired()`)
	}

	params["expire"] = now.Add(1 * time.Hour).Format(time.RFC3339)
	f, err = newBaseFilter("id", params)
	if err != nil {
		t.Fatal(err)
	}

	f.SetDynamic()
	if f.expire.IsZero() {
		t.Error(`f.expire.IsZero()`)
	}
	if f.Expired() {
		t.Error(`f.Expired()`)
	}
}

func testBaseFilterAddParams(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	params := map[string]interface{}{
		"label":    "label1",
		"disabled": true,
		"all":      true,
		"if":       "alerts.length > 1",
		"expire":   now.Format(time.RFC3339Nano),
	}

	f, err := newBaseFilter("id", params)
	if err != nil {
		t.Fatal(err)
	}
	f.SetDynamic()

	m := make(map[string]interface{})
	f.AddParams(m)

	if m["label"].(string) != "label1" {
		t.Error(`m["label"].(string) != "label1"`)
	}

	if !m["disabled"].(bool) {
		t.Error(`!m["disabled"].(bool)`)
	}

	if !m["all"].(bool) {
		t.Error(`!m["all"].(bool)`)
	}

	if m["if"].(string) != "alerts.length > 1" {
		t.Error(`m["if"].(string) != "alerts.length > 1"`)
	}

	if !now.Equal(m["expire"].(time.Time)) {
		t.Error(`!now.Equal(m["expire"].(time.Time))`)
	}
}

func TestBaseFilter(t *testing.T) {
	t.Run("All", testBaseFilterAll)
	t.Run("One", testBaseFilterOne)
	t.Run("ParseError", testBaseFilterParseError)
	t.Run("Command", testBaseFilterCommand)
	t.Run("Expire", testBaseFilterExpire)
	t.Run("AddParams", testBaseFilterAddParams)
}
