// Port of com/google/javascript/jscomp/base/JSCompObjects.java @ closure-compiler v20260415

package base

import "github.com/paullepoulpe/jsgo/internal/jutil"

// JSCompObjects holds basic Object functions used by JSComp.
type JSCompObjects struct{}

// JSCompObjects_Identical is reference equality (`a == b`).
//
// Go-only: language-adapter — Java `==` on references. Go's == on interface
// values compares pointers for reference types and values for Go value types
// (which stand in for Java enums and interned strings, where Java code
// relies on identity being value equality).
//
// Java: JSCompObjects.identical(Object, Object)
func JSCompObjects_Identical(a, b any) bool {
	return a == b
}

// NewJSCompObjects is the private constructor, which always throws.
//
// Java: JSCompObjects()
func NewJSCompObjects() *JSCompObjects {
	panic(jutil.NewAssertionError())
}
