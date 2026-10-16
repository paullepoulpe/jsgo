package jutil

// Go-only: language-adapter — java.lang.Thread identity. The port runs
// sequentially (PORTING.md §5): there is exactly one thread, so
// Thread.currentThread() always returns the same object and ThreadLocals are
// plain variables.

// Thread is java.lang.Thread (identity only).
type Thread struct {
	name string
}

var thread_main = &Thread{name: "main"}

// Thread_CurrentThread is Thread.currentThread().
func Thread_CurrentThread() *Thread { return thread_main }

// GetName is Thread.getName().
func (t *Thread) GetName() string { return t.name }

// String is Thread.toString().
func (t *Thread) String() string { return "Thread[#1," + t.name + ",5,main]" }
