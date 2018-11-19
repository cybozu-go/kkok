package kkok

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cybozu-go/well"
)

func TestAlertPool(t *testing.T) {
	t.Parallel()

	var p alertPool

	a1 := &Alert{
		From:  "from1",
		Title: "title1",
	}
	a2 := &Alert{
		From:  "from2",
		Title: "title2",
		Info:  map[string]interface{}{"info2": 2},
	}
	p.put(a1)
	p.put(a2)

	alerts := p.peek()
	if len(alerts) != 2 {
		t.Error(`len(alerts) != 2`)
	}

	alerts = p.take()
	if len(alerts) != 2 {
		t.Error(`len(alerts) != 2`)
	}

	alerts = p.take()
	if len(alerts) > 0 {
		t.Error(`len(alerts) > 0`)
	}
}

var nchan = make(chan int, 10)

type testHandler struct{}

func (t testHandler) Handle(alerts []*Alert) {
	nchan <- len(alerts)
	return
}

func TestDispatcher(t *testing.T) {
	if len(os.Getenv("TEST_ALL")) == 0 {
		t.Skip("TEST_ALL envvar is not set")
	}

	t.Parallel()

	d := NewDispatcher(3*time.Millisecond, 10*time.Millisecond, testHandler{})

	env := well.NewEnvironment(context.Background())
	env.Go(d.Run)

	d.put(&Alert{})
	d.put(&Alert{})

	// this will NOT split alerts
	time.Sleep(2 * time.Millisecond)

	d.put(&Alert{})

	// this will split alerts.
	time.Sleep(6 * time.Millisecond)
	d.put(&Alert{})

	// wait long enough
	time.Sleep(10 * time.Millisecond)

	d.put(&Alert{})
	d.put(&Alert{})

	time.Sleep(4 * time.Millisecond)

	d.put(&Alert{})
	d.put(&Alert{})

	// this will NOT split alerts as the interval is extended to 6 ms.
	time.Sleep(4 * time.Millisecond)

	d.put(&Alert{})
	d.put(&Alert{})

	n := <-nchan
	if n != 3 {
		t.Error(`n != 3`)
	}
	n = <-nchan
	if n != 1 {
		t.Error(`n != 1`)
	}
	n = <-nchan
	if n != 2 {
		t.Error(`n != 2`)
	}
	n = <-nchan
	if n != 4 {
		t.Error(`n != 4`)
	}

	env.Cancel(nil)
	err := env.Wait()
	if err != nil {
		t.Error(err)
	}
}
