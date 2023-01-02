package gaglog

import (
	"io"
	"regexp"
	"testing"
	"time"
)

var messages [][]byte

func init() {
	messages = [][]byte{
		[]byte("foo"),
		[]byte("bar"),
		[]byte("baz"),
		[]byte("qux"),
	}
}

func BenchmarkDiscard(b *testing.B) {
	for i := 0; i < b.N; i++ {
		io.Discard.Write(messages[i%len(messages)])
	}
}

func BenchmarkGagFilterAnything(b *testing.B) {
	filter := &GagFilter{
		Gags: Gags{
			regexp.MustCompile(".*"): time.Millisecond * 100,
		},
		Writer: io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Write(messages[i%len(messages)])
	}
}

func BenchmarkGagFilterFoo(b *testing.B) {
	filter := &GagFilter{
		Gags: Gags{
			regexp.MustCompile("foo"): time.Millisecond * 100,
		},
		Writer: io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Write(messages[i%len(messages)])
	}
}

func BenchmarkGagFilterFooBar(b *testing.B) {
	filter := &GagFilter{
		Gags: Gags{
			regexp.MustCompile("foo"): time.Millisecond * 100,
			regexp.MustCompile("bar"): time.Millisecond * 1000,
		},
		Writer: io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Write(messages[i%len(messages)])
	}
}
