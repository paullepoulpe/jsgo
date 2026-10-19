package jutil

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// Go-only: language-adapter — Java's object monitors (`synchronized`
// methods and blocks) and java.util.concurrent.locks.ReentrantLock.
//
// Both are re-entrant: the thread that holds the lock may take it again,
// and it is released when every acquisition has been matched by a release.
// Go's sync.Mutex is not: a goroutine that locks a mutex it already holds
// blocks forever. Translated code re-enters its monitors in ways that are
// easy to miss, e.g. ThreadSafeDelegatingErrorManager.generateReport
// formats a diagnostic, which loads an input source map, which reports a
// warning through the same synchronized error manager. So every
// translation of `synchronized` or ReentrantLock uses a Monitor, never a
// sync.Mutex (PORTING.md §5).

// Monitor is a re-entrant mutual-exclusion lock owned by a goroutine: the
// monitor of a Java object or class, or a ReentrantLock.
//
// The owner is identified by GoroutineID: a goroutine plays the part of a
// Java thread. The zero value is an unlocked monitor. A Monitor must not be
// copied after first use.
//
// A `synchronized` method or block translates to
//
//	m.Lock()
//	defer m.Unlock()
//
// at its start (or around the block, in a closure), and
// `lock.lock(); try { ... } finally { lock.unlock(); }` the same way.
type Monitor struct {
	// mu is held while some goroutine owns the monitor.
	mu sync.Mutex
	// owner is the GoroutineID of the owner, 0 when the monitor is free.
	// Other goroutines read it (to see that they are not the owner), so it
	// is atomic.
	owner atomic.Int64
	// holds is the owner's hold count; only the owner reads or writes it.
	holds int
}

// Lock acquires the monitor, blocking while another goroutine holds it. If
// the current goroutine already holds it, Lock only increments the hold
// count (monitorenter on an owned monitor; ReentrantLock.lock()).
func (m *Monitor) Lock() {
	g := GoroutineID()
	if m.owner.Load() == g {
		m.holds++
		return
	}
	m.mu.Lock()
	m.owner.Store(g)
	m.holds = 1
}

// Unlock releases one hold of the monitor; the monitor is free once every
// Lock has been matched (monitorexit; ReentrantLock.unlock()). Like
// ReentrantLock.unlock, it panics with IllegalMonitorStateException if the
// current goroutine does not hold the monitor.
func (m *Monitor) Unlock() {
	if m.owner.Load() != GoroutineID() {
		panic(NewIllegalMonitorStateException())
	}
	m.holds--
	if m.holds == 0 {
		m.owner.Store(0)
		m.mu.Unlock()
	}
}

// LockFast acquires the monitor for a leaf synchronized method whose body
// cannot re-enter the same monitor, but which may be called re-entrantly while
// another method on the same object holds Lock(). When the monitor has no
// re-entrant owner (owner == 0), it avoids the GoroutineID() stack walk.
func (m *Monitor) LockFast() {
	if owner := m.owner.Load(); owner != 0 && owner == GoroutineID() {
		m.holds++
		return
	}
	m.mu.Lock()
	m.holds = 1
}

// UnlockFast releases one hold acquired by LockFast.
func (m *Monitor) UnlockFast() {
	m.holds--
	if m.holds == 0 {
		m.mu.Unlock()
	}
}

// IsHeldByCurrentGoroutine reports whether the current goroutine holds the
// monitor (Thread.holdsLock(obj); ReentrantLock.isHeldByCurrentThread()).
func (m *Monitor) IsHeldByCurrentGoroutine() bool {
	return m.owner.Load() == GoroutineID()
}

// GetHoldCount is the number of holds the current goroutine has on the
// monitor, 0 if it does not hold it (ReentrantLock.getHoldCount()).
func (m *Monitor) GetHoldCount() int {
	if m.owner.Load() != GoroutineID() {
		return 0
	}
	return m.holds
}

// GoroutineID returns the identity of the current goroutine: the number
// the runtime gives it, as in the "goroutine N [...]" header of its stack
// trace, which never changes during the goroutine's life and is never 0.
// It is the owner identity Monitor records.
//
// Go has no API for this on purpose; the only supported way to obtain it
// is runtime.Stack, whose first line is that header. It costs a stack walk
// (bounded by the runtime's traceback limit), so it belongs on the paths
// that take a monitor, not in per-node code.
func GoroutineID() int64 {
	var buf [64]byte
	b := buf[:runtime.Stack(buf[:], false)]
	// "goroutine 123 [running]:\n..."
	const prefix = "goroutine "
	if len(b) < len(prefix) || string(b[:len(prefix)]) != prefix {
		panic(NewIllegalStateExceptionMsg("unexpected goroutine header: " + string(b)))
	}
	var id int64
	for _, c := range b[len(prefix):] {
		if c < '0' || c > '9' {
			break
		}
		id = id*10 + int64(c-'0')
	}
	if id == 0 {
		panic(NewIllegalStateExceptionMsg("unexpected goroutine header: " + string(b)))
	}
	return id
}
