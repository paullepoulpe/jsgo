// Port of com/google/javascript/jscomp/base/JSCompStrings.java @ closure-compiler v20260415

package base

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSCompStrings holds basic String functions used by JSComp.
type JSCompStrings struct{}

// JSCompStrings_Lines joins the lines with "\n".
//
// Java: JSCompStrings.lines(String...)
func JSCompStrings_Lines(lines ...string) string {
	return strings.Join(lines, "\n")
}

// JSCompStrings_LinesOfIterable joins the lines with "\n".
//
// Java: JSCompStrings.lines(Iterable<String>)
func JSCompStrings_LinesOfIterable(lines []string) string {
	return strings.Join(lines, "\n")
}

// NewJSCompStrings is the private constructor, which always throws.
//
// Java: JSCompStrings()
func NewJSCompStrings() *JSCompStrings {
	panic(jutil.NewAssertionError())
}
