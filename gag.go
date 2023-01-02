// Package gaglog augments the standard log package with duration based filtering.
package gaglog

import (
	"io"
	"regexp"
	"sync"
	"time"
)

type Gags map[*regexp.Regexp]time.Duration

// GagFilter is an io.Writer that can be used with a logger that
// will gag out log messages that match a regex for a duration.
type GagFilter struct {
	// The map with keys of regex.Regexp and values of time.Durations.
	// Example: Gags{regexp.MustCompile("bar"): time.Millisecond * 2}
	Gags

	// The underlying io.Writer where log messages that pass the filter
	// will be set.
	Writer io.Writer

	// Once the filter is in use, it is not safe to modify
	// the structure.
	once  sync.Once
	clock clock

	// Internally package works with time.Time instead of time.Duration.
	// On first call of Write, gags is initialized (see gaglog.init).
	// gags does not have to be protected by a sync.Mutex as log.Logger
	// serializes access to the Writer.
	gags map[*regexp.Regexp]time.Time
}

func (f *GagFilter) Write(p []byte) (n int, err error) {
	// Note in general that io.Writer can receive any byte sequence
	// to write, but the "log" package always guarantees that we only
	// get a single line.

	if f.check(p) {
		return len(p), nil
	}

	return f.Writer.Write(p)
}

func (f *GagFilter) check(line []byte) bool {
	f.once.Do(f.init)

	for r, d := range f.Gags {
		if r.FindIndex(line) != nil {
			now := f.clock.Now()
			if now.Before(f.gags[r].Add(d)) {
				return true
			}
			f.gags[r] = now
			break
		}
	}

	return false
}

func (f *GagFilter) init() {
	if f.clock == nil {
		f.clock = realClock{}
	}

	f.gags = make(map[*regexp.Regexp]time.Time, len(f.Gags))
	for r := range f.Gags {
		f.gags[r] = time.Time{}
	}
}

// clock interface is used in favor of easier testing of duration/time.
// source: https://stackoverflow.com/a/18970352
type clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}
type realClock struct{}

// Now is wrapper for stdlib time.Now
func (realClock) Now() time.Time {
	return time.Now()
}

// After is wrapper for stdlib time.After
func (realClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}
