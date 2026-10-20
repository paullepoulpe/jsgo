package jutil

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type outputWriterReentrySink struct {
	write      func([]byte) (int, error)
	closeCount int
}

func (s *outputWriterReentrySink) Write(p []byte) (int, error) { return s.write(p) }
func (s *outputWriterReentrySink) Close() error {
	s.closeCount++
	return nil
}

func outputWriterReentryPanic(f func()) (caught any) {
	defer func() { caught = recover() }()
	f()
	return nil
}

// A stream may synchronously call its writer before it copies the borrowed p.
// The nested write must not change the outer bytes still owned by that call.
func TestOutputStreamWriterReentrantWritesKeepBorrowedBytes(t *testing.T) {
	for _, mode := range []string{"utf8_outer_larger", "utf8_nested_larger", "buffered_utf16"} {
		t.Run(mode, func(t *testing.T) {
			outer := strings.Repeat("O", 16384) + "é😀"
			nested := strings.Repeat("N", 8192) + "Ω😀"
			if mode == "utf8_nested_larger" {
				outer, nested = nested, outer
			}
			var captured bytes.Buffer
			var writer *OutputStreamWriter
			var buffered *BufferedWriter
			fired := false
			sink := &outputWriterReentrySink{}
			sink.write = func(p []byte) (int, error) {
				if !fired {
					fired = true
					if buffered != nil {
						buffered.Write(nested)
						buffered.Flush()
					} else {
						writer.Write(nested)
					}
				}
				return captured.Write(p)
			}
			writer = NewOutputStreamWriter(sink)
			if mode == "buffered_utf16" {
				buffered = NewBufferedWriter(writer)
				buffered.Write(outer)
				buffered.Flush()
			} else {
				writer.Write(outer)
			}
			if got := captured.String(); got != nested+outer {
				t.Fatalf("nested write changed in-flight outer bytes: got %d bytes, want %d", len(got), len(nested)+len(outer))
			}
			// Further ordinary writes and empty writes remain usable after the
			// reentrant operation, regardless of which call needed more scratch.
			writer.Write("")
			writer.Flush()
			writer.Write("after")
			if captured.String() != nested+outer+"after" {
				t.Fatal("writer failed to resume after reentrant output")
			}
			if buffered != nil {
				buffered.Close()
			} else {
				writer.Close()
			}
			if sink.closeCount != 1 {
				t.Fatalf("underlying close count = %d, want 1", sink.closeCount)
			}
		})
	}

	t.Run("split_surrogate_state_is_not_rolled_back", func(t *testing.T) {
		var captured bytes.Buffer
		var writer *OutputStreamWriter
		nested := strings.Repeat("N", 8192)
		fired := false
		sink := &outputWriterReentrySink{write: func(p []byte) (int, error) {
			if !fired {
				fired = true
				writer.Write(nested)
			}
			return captured.Write(p)
		}}
		writer = NewOutputStreamWriter(sink)
		buffered := NewBufferedWriter(writer)
		// The first flush ends in a high surrogate. The nested write emits
		// its replacement; the remaining outer low surrogate is then lone.
		buffered.Write(strings.Repeat("a", 8191) + "😀Z")
		buffered.Close()
		want := "?" + nested + strings.Repeat("a", 8191) + "?Z"
		if got := captured.String(); got != want {
			t.Fatalf("reentrant write changed split-surrogate encoding: got %d bytes, want %d", len(got), len(want))
		}
	})
}

func TestOutputStreamWriterReentrantFailuresKeepBytesAndErrors(t *testing.T) {
	for _, mode := range []string{"partial_error", "sink_panic", "buffered_write_string_error"} {
		t.Run(mode, func(t *testing.T) {
			outer := strings.Repeat("O", 16384) + "é😀"
			nested := strings.Repeat("N", 8192) + "Ω"
			var captured bytes.Buffer
			var writer *OutputStreamWriter
			fired, calls := false, 0
			panicValue := &struct{ message string }{"sink panic"}
			sink := &outputWriterReentrySink{write: func(p []byte) (int, error) {
				calls++
				if !fired {
					fired = true
					writer.Write(nested)
					captured.Write(p[:7])
					if mode == "sink_panic" {
						panic(panicValue)
					}
					return 7, errors.New("partial sink failure")
				}
				return captured.Write(p)
			}}
			writer = NewOutputStreamWriter(sink)
			var buffered *BufferedWriter
			if mode == "buffered_write_string_error" {
				buffered = NewBufferedWriter(writer)
				n, err := buffered.WriteString(outer)
				if n != 0 || err == nil || err.Error() != "partial sink failure" {
					t.Fatalf("WriteString failure = (%d, %v), want (0, partial sink failure)", n, err)
				}
			} else {
				caught := outputWriterReentryPanic(func() { writer.Write(outer) })
				if mode == "sink_panic" {
					if caught != panicValue {
						t.Fatalf("sink panic was replaced: %v", caught)
					}
				} else if err, ok := caught.(*IOException); !ok || err.GetMessage() != "partial sink failure" {
					t.Fatalf("partial error panic = %v, want IOException(partial sink failure)", caught)
				}
			}
			if got := captured.String(); got != nested+outer[:7] || calls != 2 {
				t.Fatalf("failure changed borrowed prefix or retried output: got %d bytes in %d calls", len(got), calls)
			}
			writer.Write("after")
			if captured.String() != nested+outer[:7]+"after" {
				t.Fatal("writer could not resume after downstream failure")
			}
			if buffered != nil {
				buffered.Close()
			} else {
				writer.Close()
			}
			if sink.closeCount != 1 {
				t.Fatalf("underlying close count = %d, want 1", sink.closeCount)
			}
		})
	}
}

func TestBufferedWriterReentrantCloseFailureKeepsCloseSemantics(t *testing.T) {
	var captured bytes.Buffer
	var writer *OutputStreamWriter
	fired := false
	nested := strings.Repeat("N", 8192)
	sink := &outputWriterReentrySink{write: func(p []byte) (int, error) {
		if !fired {
			fired = true
			writer.Write(nested)
			captured.Write(p[:3])
			return 3, errors.New("close flush failure")
		}
		return captured.Write(p)
	}}
	writer = NewOutputStreamWriter(sink)
	buffered := NewBufferedWriter(writer)
	buffered.Write("outer tail")
	caught := outputWriterReentryPanic(buffered.Close)
	if err, ok := caught.(*IOException); !ok || err.GetMessage() != "close flush failure" {
		t.Fatalf("close failure = %v, want IOException(close flush failure)", caught)
	}
	if got := captured.String(); got != nested+"out" {
		t.Fatalf("close flush changed borrowed bytes: got %d bytes", len(got))
	}
	if sink.closeCount != 1 {
		t.Fatalf("failed close reached stream %d times, want 1", sink.closeCount)
	}
	buffered.Close()
	writer.Close()
	if sink.closeCount != 1 {
		t.Fatal("second close was not a no-op")
	}
	if got := throwableString(func() { buffered.Write("") }); got != "java.io.IOException: Stream closed" {
		t.Fatalf("write after failed close: %s", got)
	}
}
