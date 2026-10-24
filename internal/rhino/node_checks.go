package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Preconditions.checkNotNull for a Node. The generic jutil versions take an
// any, which would put the 16-byte handle on the heap on every call; these
// do the same check without allocating.

// CheckNotNull is Preconditions.checkNotNull(T) for a Node.
func CheckNotNull(n Node) Node {
	if n == 0 {
		panic(jutil.NewNullPointerException())
	}
	return n
}

// CheckNotNullMsg is Preconditions.checkNotNull(T, String, Object...) for a Node.
func CheckNotNullMsg(n Node, template string, args ...any) Node {
	if n == 0 {
		failNotNull(template, args...)
	}
	return n
}

// CheckNotNullMsg1 is Preconditions.checkNotNull(T, String, Object) for a Node.
func CheckNotNullMsg1[A any](n Node, template string, p1 A) Node {
	if n == 0 {
		failNotNull(template, p1)
	}
	return n
}

// CheckNotNullMsg2 is Preconditions.checkNotNull(T, String, Object, Object) for a Node.
func CheckNotNullMsg2[A, B any](n Node, template string, p1 A, p2 B) Node {
	if n == 0 {
		failNotNull(template, p1, p2)
	}
	return n
}

//go:noinline
func failNotNull(template string, args ...any) {
	panic(jutil.NewNullPointerExceptionMsg(jutil.LenientFormat(template, args...)))
}
