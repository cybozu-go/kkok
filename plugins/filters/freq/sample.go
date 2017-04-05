package freq

import "time"

// Sample keeps the latest events for a given period.
type Sample struct {
	duration time.Duration
	samples  []time.Time
}

// NewSample creates a new Sample object.
func NewSample(duration time.Duration) *Sample {
	return &Sample{duration, make([]time.Time, 0, 10)}
}

// gc expunges stale samples by overwriting non-stale samples
// onto the stale ones.  No (re-)allocation happens.
func (s *Sample) gc(now time.Time) {
	expire := now.Add(-s.duration)

	idx := -1
	for i, t := range s.samples {
		if t.After(expire) {
			break
		}
		idx = i
	}

	if idx == -1 {
		return
	}

	remaining := len(s.samples) - idx - 1

	if remaining == 0 {
		s.samples = s.samples[:0]
		return
	}

	// move non-stale samples, if any.
	liveIdx := idx + 1
	for i := 0; i < remaining; i++ {
		s.samples[i] = s.samples[liveIdx+i]
	}

	s.samples = s.samples[0:remaining]
}

// Add adds a sample.
func (s *Sample) Add(t time.Time) {
	s.samples = append(s.samples, t)
}

// Count returns the number of samples for the latest duration.
func (s *Sample) Count() int {
	s.gc(time.Now())
	return len(s.samples)
}
