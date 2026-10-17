package jutil

import (
	"sync"
	"testing"
	"time"
)

func TestMonitorReentrant(t *testing.T) {
	var m Monitor
	if m.IsHeldByCurrentGoroutine() || m.GetHoldCount() != 0 {
		t.Fatal("a new monitor is held")
	}
	m.Lock()
	m.Lock() // would deadlock with a sync.Mutex
	func() {
		m.Lock()
		defer m.Unlock()
		if got := m.GetHoldCount(); got != 3 {
			t.Fatalf("hold count %d, want 3", got)
		}
	}()
	m.Unlock()
	if !m.IsHeldByCurrentGoroutine() || m.GetHoldCount() != 1 {
		t.Fatalf("after 3 locks and 2 unlocks: held=%v count=%d", m.IsHeldByCurrentGoroutine(), m.GetHoldCount())
	}
	m.Unlock()
	if m.IsHeldByCurrentGoroutine() || m.GetHoldCount() != 0 {
		t.Fatal("still held after the last unlock")
	}
}

func TestMonitorUnlockNotHeld(t *testing.T) {
	var m Monitor
	caught := false
	Try(func() { m.Unlock() }, func(e IIllegalMonitorStateException) { caught = true })
	if !caught {
		t.Fatal("unlocking a free monitor did not throw IllegalMonitorStateException")
	}

	// Another goroutine's hold cannot be released.
	m.Lock()
	done := make(chan bool)
	go func() {
		defer close(done)
		Try(func() { m.Unlock() }, func(e IIllegalMonitorStateException) { done <- true })
	}()
	if !<-done {
		t.Fatal("unlocking another goroutine's monitor did not throw")
	}
	m.Unlock()
}

func TestMonitorExcludesOtherGoroutines(t *testing.T) {
	var m Monitor
	m.Lock()
	m.Lock()
	acquired := make(chan struct{})
	go func() {
		m.Lock()
		if m.GetHoldCount() != 1 {
			t.Errorf("second goroutine hold count %d, want 1", m.GetHoldCount())
		}
		m.Unlock()
		close(acquired)
	}()
	m.Unlock()
	select {
	case <-acquired:
		t.Fatal("another goroutine acquired a monitor that is still held")
	case <-time.After(50 * time.Millisecond):
	}
	m.Unlock()
	<-acquired

	// A counter updated under the monitor by many goroutines.
	var wg sync.WaitGroup
	n := 0
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				m.Lock()
				m.Lock()
				n++
				m.Unlock()
				m.Unlock()
			}
		}()
	}
	wg.Wait()
	if n != 8*200 {
		t.Fatalf("counter %d, want %d", n, 8*200)
	}
}

func TestGoroutineID(t *testing.T) {
	id := GoroutineID()
	if id <= 0 || GoroutineID() != id {
		t.Fatalf("GoroutineID %d is not stable and positive", id)
	}
	other := make(chan int64)
	go func() { other <- GoroutineID() }()
	if o := <-other; o == id || o <= 0 {
		t.Fatalf("another goroutine has id %d (this one %d)", o, id)
	}
}

func BenchmarkMonitorLockUnlock(b *testing.B) {
	var m Monitor
	for i := 0; i < b.N; i++ {
		m.Lock()
		m.Unlock()
	}
}
