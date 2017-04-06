package freq

import (
	"testing"
	"time"
)

func testSampleGC(t *testing.T) {
	t.Parallel()

	s := NewSample(10 * time.Millisecond)

	if len(s.samples) != 0 {
		t.Fatal(`len(s.samples) != 0`)
	}

	// empty
	now := time.Now()
	s.gc(now)
	if len(s.samples) != 0 {
		t.Error(`len(s.samples) != 0`)
	}

	// stale only
	stale := now.Add(-20 * time.Millisecond)
	s.Add(stale)
	s.gc(now)
	if len(s.samples) != 0 {
		t.Error(`len(s.samples) != 0`)
	}

	// non-stale only
	fresh1 := now.Add(-5 * time.Millisecond)
	fresh2 := now.Add(-4 * time.Millisecond)
	s.Add(fresh1)
	s.Add(fresh2)
	s.gc(now)
	if len(s.samples) != 2 {
		t.Error(`len(s.samples) != 2`)
	}

	// stale and non-stale
	s = NewSample(10 * time.Millisecond)
	s.Add(stale)
	s.Add(fresh1)
	s.Add(fresh2)
	s.gc(now)
	if len(s.samples) != 2 {
		t.Fatal(`len(s.samples) != 2`)
	}
	if !s.samples[0].Equal(fresh1) {
		t.Error(`!s.samples[0].Equal(fresh1)`)
	}
	if !s.samples[1].Equal(fresh2) {
		t.Error(`!s.samples[1].Equal(fresh2)`)
	}

	s = NewSample(10 * time.Millisecond)
	s.Add(stale)
	s.Add(stale)
	s.Add(fresh1)
	s.gc(now)
	if len(s.samples) != 1 {
		t.Fatal(`len(s.samples) != 1`)
	}
	if !s.samples[0].Equal(fresh1) {
		t.Error(`!s.samples[0].Equal(fresh1)`)
	}
}

func testSampleCount(t *testing.T) {
	t.Parallel()

	s := NewSample(10 * time.Millisecond)
	if s.Count() != 0 {
		t.Error(`s.Count() != 0`)
	}

	s.Add(time.Now())
	if s.Count() != 1 {
		t.Error(`s.Count() != 1`)
	}

	s.Add(time.Now())
	if s.Count() != 2 {
		t.Error(`s.Count() != 2`)
	}

	time.Sleep(11 * time.Millisecond)
	if s.Count() != 0 {
		t.Error(`s.Count() != 0`)
	}
}

func TestSample(t *testing.T) {
	t.Run("GC", testSampleGC)
	t.Run("Count", testSampleCount)
}
