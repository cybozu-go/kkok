package group

import (
	"reflect"
	"testing"
	"time"

	"github.com/cybozu-go/kkok"
)

func testParamsDefault(t *testing.T) {
	t.Parallel()

	f := &filter{}
	pp := f.Params()

	if pp.Type != filterType {
		t.Error(`pp.Type != filterType`)
	}

	if len(pp.Params) != 0 {
		t.Error(`len(pp.Params) != 0`, pp.Params)
	}
}

func testParamsExplicit(t *testing.T) {
	t.Parallel()

	f := &filter{
		origBy:  `alert.Host`,
		from:    "from1",
		title:   "title1",
		message: "msg",
		routes:  []string{"r1", "r2"},
	}

	pp := f.Params()
	if pp.Type != filterType {
		t.Error(`pp.Type != filterType`)
	}

	params := map[string]interface{}{
		"by":      "alert.Host",
		"from":    "from1",
		"title":   "title1",
		"message": "msg",
		"routes":  []string{"r1", "r2"},
	}
	if !reflect.DeepEqual(pp.Params, params) {
		t.Error(`!reflect.DeepEqual(pp.Params, params)`)
		t.Log(pp.Params)
		t.Log(params)
	}
}

func testParams(t *testing.T) {
	t.Run("Default", testParamsDefault)
	t.Run("Explicit", testParamsExplicit)
}

func testMergeZero(t *testing.T) {
	t.Parallel()

	f := &filter{}
	a := f.mergeAlerts(nil)
	if a != nil {
		t.Error(`a != nil`)
	}
}

func testMergeOne(t *testing.T) {
	t.Parallel()

	f := &filter{}
	a := &kkok.Alert{
		From:  "from1",
		Title: "title1",
	}

	b := f.mergeAlerts([]*kkok.Alert{a})
	if a != b {
		t.Error(`a != b`)
	}
}

func testMergeSame(t *testing.T) {
	t.Parallel()

	f := &filter{
		from:    "newfrom",
		title:   "newtitle",
		message: "newmsg",
		routes:  []string{"newr"},
	}

	now := time.Now()
	old := now.Add(-1 * time.Minute)
	a := &kkok.Alert{
		From:    "from1",
		Date:    old,
		Host:    "host1",
		Title:   "title1",
		Message: "msg",
		Routes:  []string{"r1"},
		Info:    map[string]interface{}{"foo": "bar"},
	}
	b := &kkok.Alert{
		From:    "from1",
		Date:    old,
		Host:    "host1",
		Title:   "title1",
		Message: "msg",
		Routes:  []string{"r1"},
		Info:    map[string]interface{}{"foo": "bar"},
	}

	c := f.mergeAlerts([]*kkok.Alert{a, b})
	if c == a {
		t.Error(`c == a`)
	}
	if c == b {
		t.Error(`c == b`)
	}

	if c.From != "from1" {
		t.Error(`c.From != "from1"`)
	}
	if c.Date.Equal(old) {
		t.Error(`c.Date.Equal(old)`)
	}
	if c.Host != "host1" {
		t.Error(`c.Host != "host1"`)
	}
	if c.Title != "title1" {
		t.Error(`c.Title != "title1"`)
	}
	if c.Message != "msg" {
		t.Error(`c.Message != "msg"`)
	}
	if !reflect.DeepEqual(c.Routes, []string{"newr"}) {
		t.Error(`!reflect.DeepEqual(c.Routes, []string{"newr"})`)
	}
	if c.Info != nil {
		t.Error(`c.Info != nil`)
	}

	if len(c.Sub) != 2 {
		t.Fatal(`len(c.Sub) != 2`)
	}
	if c.Sub[0] != a {
		t.Error(`c.Sub[0] != a`)
	}
	if c.Sub[1] != b {
		t.Error(`c.Sub[1] != b`)
	}
}

func testMergeDifferent(t *testing.T) {
	t.Parallel()

	f := &filter{
		from:    "newfrom",
		title:   "newtitle",
		message: "newmsg",
		routes:  []string{"newr"},
	}

	now := time.Now()
	old1 := now.Add(-1 * time.Minute)
	old2 := now.Add(-2 * time.Minute)
	a := &kkok.Alert{
		From:    "from1",
		Date:    old1,
		Host:    "host1",
		Title:   "title1",
		Message: "msg1",
		Routes:  []string{"r1"},
		Info:    map[string]interface{}{"foo": "bar1"},
	}
	b := &kkok.Alert{
		From:    "from2",
		Date:    old2,
		Host:    "host2",
		Title:   "title2",
		Message: "msg2",
		Routes:  []string{"r2"},
		Info:    map[string]interface{}{"foo": "bar2"},
	}

	c := f.mergeAlerts([]*kkok.Alert{a, b})
	if c == a {
		t.Error(`c == a`)
	}
	if c == b {
		t.Error(`c == b`)
	}

	if c.From != "newfrom" {
		t.Error(`c.From != "newfrom"`)
	}
	if c.Date.Equal(old1) || c.Date.Equal(old2) {
		t.Error(`c.Date.Equal(old1) || c.Date.Equal(old2)`)
	}
	if c.Host != "localhost" {
		t.Error(`c.Host != "localhost"`)
	}
	if c.Title != "newtitle" {
		t.Error(`c.Title != "newtitle"`)
	}
	if c.Message != "newmsg" {
		t.Error(`c.Message != "newmsg"`)
	}
	if !reflect.DeepEqual(c.Routes, []string{"newr"}) {
		t.Error(`!reflect.DeepEqual(c.Routes, []string{"newr"})`)
	}
	if c.Info != nil {
		t.Error(`c.Info != nil`)
	}

	if len(c.Sub) != 2 {
		t.Fatal(`len(c.Sub) != 2`)
	}
	if c.Sub[0] != a {
		t.Error(`c.Sub[0] != a`)
	}
	if c.Sub[1] != b {
		t.Error(`c.Sub[1] != b`)
	}
}

func testMergeDefault(t *testing.T) {
	t.Parallel()

	f := &filter{}
	err := f.Init("f", nil)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	old1 := now.Add(-1 * time.Minute)
	old2 := now.Add(-2 * time.Minute)
	a := &kkok.Alert{
		From:    "from1",
		Date:    old1,
		Host:    "host1",
		Title:   "title1",
		Message: "msg1",
		Routes:  []string{"r1"},
		Info:    map[string]interface{}{"foo": "bar1"},
	}
	b := &kkok.Alert{
		From:    "from2",
		Date:    old2,
		Host:    "host2",
		Title:   "title2",
		Message: "msg2",
		Routes:  []string{"r2"},
		Info:    map[string]interface{}{"foo": "bar2"},
	}

	c := f.mergeAlerts([]*kkok.Alert{a, b})
	if c.From != defaultFrom+"f" {
		t.Error(`c.From != defaultFrom+"f"`)
	}
	if c.Title != defaultTitle {
		t.Error(`c.Title != defaultTitle`)
	}
	if len(c.Message) != 0 {
		t.Error(`len(c.Message) != 0`)
	}
	if len(c.Routes) != 0 {
		t.Error(`len(c.Routes) != 0`)
	}
}

func testMerge(t *testing.T) {
	t.Run("Zero", testMergeZero)
	t.Run("One", testMergeOne)
	t.Run("Same", testMergeSame)
	t.Run("Different", testMergeDifferent)
	t.Run("Default", testMergeDefault)
}

var testAlertsData = []*kkok.Alert{
	{From: "from1", Info: map[string]interface{}{"info1": 3}},
	{From: "from2", Info: map[string]interface{}{"info1": 3}},
	{From: "from1", Info: map[string]interface{}{"info1": 3}},
	{From: "from2", Info: map[string]interface{}{"info1": 3}},
	{From: "from3", Title: "t1"},
	{From: "from3", Title: "t2"},
	{From: "from3", Title: "t1"},
}

func testProcessAll(t *testing.T) {
	t.Parallel()

	f := &filter{}

	alerts, err := f.Process(testAlertsData)
	if err != nil {
		t.Fatal(err)
	}
	if len(alerts) != 1 {
		t.Error(`len(alerts) != 1`)
	}

	f.Init("f", map[string]interface{}{"all": true})
	alerts, err = f.Process(testAlertsData)
	if err != nil {
		t.Fatal(err)
	}
	if len(alerts) != 1 {
		t.Error(`len(alerts) != 1`)
	}
}

func testProcessAllIf(t *testing.T) {
	t.Parallel()

	f := &filter{}
	err := f.Init("f", map[string]interface{}{
		"all": true,
		"if":  "alerts.length > 10",
	})
	if err != nil {
		t.Fatal(err)
	}

	alerts, err := f.Process(testAlertsData)
	if err != nil {
		t.Fatal(err)
	}
	if len(alerts) != len(testAlertsData) {
		t.Error(`len(alerts) != len(testAlertsData)`)
	}

	err = f.Init("f", map[string]interface{}{
		"all": true,
		"if":  "alerts.length > 3",
	})
	if err != nil {
		t.Fatal(err)
	}

	alerts, err = f.Process(testAlertsData)
	if err != nil {
		t.Fatal(err)
	}
	if len(alerts) != 1 {
		t.Error(`len(alerts) != 1`)
	}
}

func testProcessIf(t *testing.T) {
	t.Parallel()

	f := &filter{}
	err := f.Init("f", map[string]interface{}{
		"if": "alert.From=='from3'",
	})
	if err != nil {
		t.Fatal(err)
	}

	alerts, err := f.Process(testAlertsData)
	if err != nil {
		t.Fatal(err)
	}
	if len(alerts) != 5 {
		t.Fatal(`len(alerts) != 5`)
	}

	results := []string{"from1", "from2", "from1", "from2", "from3"}
	for i, a := range alerts {
		if a.From != results[i] {
			t.Error(`a.From != results[i]; i=`, i)
		}
	}

	f.from = "test"
	err = f.Init("f", map[string]interface{}{
		"if": "alert.From=='from1' || alert.From=='from3'",
	})
	if err != nil {
		t.Fatal(err)
	}

	alerts, err = f.Process(testAlertsData)
	if err != nil {
		t.Fatal(err)
	}
	if len(alerts) != 3 {
		t.Fatal(`len(alerts) != 3`, alerts)
	}

	if alerts[0].From != "from2" {
		t.Error(`alerts[0].From != "from2"`)
	}
	if alerts[1].From != "from2" {
		t.Error(`alerts[1].From != "from2"`)
	}
	if alerts[2].From != "test" {
		t.Error(`alerts[2].From != "test"`)
	}
}

func testProcessError(t *testing.T) {
	t.Parallel()

	f := &filter{}
	s, err := kkok.CompileJS(`alerts.length`)
	if err != nil {
		t.Fatal(err)
	}
	f.by = s

	_, err = f.Process(testAlertsData)
	if err == nil {
		t.Error(`err == nil`)
	}
}

func testProcessBy(t *testing.T) {
	t.Parallel()

	f := &filter{}
	s, err := kkok.CompileJS(`alert.From`)
	if err != nil {
		t.Fatal(err)
	}
	f.by = s

	alerts, err := f.Process(testAlertsData)
	if err != nil {
		t.Fatal(err)
	}
	if len(alerts) != 3 {
		t.Fatal(`len(alerts) != 3`)
	}

	s, err = kkok.CompileJS(`
if( alert.Info.info1 ) {
    alert.Info.info1;
} else {
    alert.Title;
}`)
	if err != nil {
		t.Fatal(err)
	}
	f.from = "test"
	f.by = s

	alerts, err = f.Process(testAlertsData)
	if err != nil {
		t.Fatal(err)
	}
	if len(alerts) != 3 {
		t.Fatal(`len(alerts) != 3`)
	}
	for _, a := range alerts {
		switch a.From {
		case "test":
			if len(a.Title) != 0 {
				t.Error(`len(a.Title) != 0`)
			}
		case "from3":
			if a.Title != "t1" && a.Title != "t2" {
				t.Error(`a.Title != "t1" && a.Title != "t2"`)
			}
		default:
			t.Error("unexpected From:", a.From)
		}
	}
}

func testProcess(t *testing.T) {
	t.Run("All", testProcessAll)
	t.Run("AllIf", testProcessAllIf)
	t.Run("If", testProcessIf)
	t.Run("Error", testProcessError)
	t.Run("By", testProcessBy)
}

func TestFilter(t *testing.T) {
	t.Run("Params", testParams)
	t.Run("Merge", testMerge)
	t.Run("Process", testProcess)
}
