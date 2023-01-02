package gaglog

import (
	"bytes"
	"io"
	"log"
	"regexp"
	"sync"
	"testing"
	"time"
)

type fakeClock struct {
	mu  sync.Mutex
	old time.Time
}

// Now increments a clock by one millisecond each time it is called.
func (f *fakeClock) Now() time.Time {
	f.mu.Lock()
	f.old = f.old.Add(time.Millisecond)
	now := f.old
	f.mu.Unlock()

	return now
}

// After is wrapper for stdlib time.After.
func (*fakeClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func TestGagFilter_impl(t *testing.T) {
	var _ io.Writer = new(GagFilter)
}

func TestGagFilter(t *testing.T) {
	buf := new(bytes.Buffer)
	filter := &GagFilter{
		Writer: buf,
		Gags: Gags{
			regexp.MustCompile("bar"): time.Millisecond * 2,
		},

		clock: &fakeClock{old: time.Now()},
	}

	logger := log.New(filter, "", 0)
	for i := 0; i <= 4; i++ {
		logger.Printf("foo-%d\n", i) // all foo-0 - foo-4 should be printed
		logger.Printf("bar-%d\n", i) // bar-0, bar-2 and bar-4 should be printed
	}

	result := buf.String()
	expected := "foo-0\nbar-0\nfoo-1\nfoo-2\nbar-2\nfoo-3\nfoo-4\nbar-4\n"
	if result != expected {
		t.Fatalf("bad: %#v", result)
	}
}
