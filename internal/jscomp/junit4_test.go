package jscomp

// Go-only: test-helper — the JUnit 4 runner for ported Java test classes.
//
// A ported test class is a struct (usually embedding CompilerTestCaseBase or
// CompilerTypeTestCaseBase) whose Java @Test methods become exported Go
// methods with the same name and no parameters or results:
//
//	@Test public void testFoo() { ... }   →   func (c *FooTest) TestFoo() { ... }
//
// and one Go test function runs them all, each on a fresh instance (JUnit 4
// creates one instance per test method), between the @Before and @After
// methods:
//
//	func TestFooTest(t *testing.T) { RunJUnit4(t, func() JUnitTestCase { return NewFooTest() }) }
//
// Each test method becomes a subtest named after it, so
// `go test -run 'TestFooTest/TestFoo$'` runs a single Java test.
//
// Assertions (AssertThat, AssertWithMessage, Fail, the Truth subjects) panic
// with *jutil.AssertionError, as Java's throw AssertionError; the runner reports
// that panic (and any other panic, e.g. a Java exception) as a failure of
// that test method only.

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"sort"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JUnitTestCase is one instance of a ported JUnit 4 test class.
type JUnitTestCase interface {
	// JUnitSetT gives the instance the Go test running its test method.
	JUnitSetT(t *testing.T)
	// JUnitBefore runs the class's @Before methods (superclass ones first).
	JUnitBefore()
	// JUnitAfter runs the class's @After methods (subclass ones first).
	JUnitAfter()
}

// RunJUnit4 runs every test method of the class built by newTest: every
// exported method whose name starts with "Test" and that takes no
// parameters and returns nothing. Methods are run in name order.
func RunJUnit4(t *testing.T, newTest func() JUnitTestCase) {
	t.Helper()
	proto := newTest()
	typ := reflect.TypeOf(proto)
	var names []string
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if !strings.HasPrefix(m.Name, "Test") || m.Type.NumIn() != 1 || m.Type.NumOut() != 0 {
			continue
		}
		names = append(names, m.Name)
	}
	sort.Strings(names)
	if len(names) == 0 {
		t.Fatalf("%s has no test methods", typ)
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			runJUnit4Method(t, newTest(), name)
		})
	}
}

// runJUnit4Method runs one test method on a fresh instance: @Before, the
// test, then @After (also when the test failed), like JUnit's
// RunBefores/RunAfters statements.
func runJUnit4Method(t *testing.T, inst JUnitTestCase, name string) {
	t.Helper()
	inst.JUnitSetT(t)
	method := reflect.ValueOf(inst).MethodByName(name)
	defer junit4Guard(t, "@After", inst.JUnitAfter)
	junit4Guard(t, "", func() {
		inst.JUnitBefore()
		method.Call(nil)
	})
}

// junit4Guard runs body, reporting a panic as a test failure. A test that
// already failed through t.Fatal (runtime.Goexit) is not reported again.
func junit4Guard(t *testing.T, phase string, body func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		prefix := ""
		if phase != "" {
			prefix = phase + ": "
		}
		if ae, ok := r.(*jutil.AssertionError); ok {
			t.Errorf("%s%s\n%s", prefix, ae.GetMessage(), junit4ShortStack())
			return
		}
		t.Errorf("%s%s\n%s", prefix, junit4Describe(r), debug.Stack())
	}()
	body()
}

// junit4Describe renders a recovered panic like a Java stack trace head.
func junit4Describe(r any) string {
	if th, ok := r.(jutil.Throwable); ok {
		return th.Error()
	}
	if err, ok := r.(error); ok {
		return err.Error()
	}
	return fmt.Sprint(r)
}

// junit4ShortStack is the stack of a failed assertion without the runtime,
// reflection and harness frames, so that the failing test line shows first.
func junit4ShortStack() string {
	lines := strings.Split(string(debug.Stack()), "\n")
	var out []string
	for i := 1; i+1 < len(lines); i += 2 {
		fn, loc := lines[i], lines[i+1]
		if strings.Contains(fn, "runtime/") || strings.Contains(fn, "reflect.") ||
			strings.Contains(fn, "testing.") || strings.Contains(fn, "jscomp.junit4") ||
			strings.Contains(fn, "jscomp.RunJUnit4") || strings.Contains(loc, "truth_test.go") ||
			strings.Contains(fn, "panic(") {
			continue
		}
		out = append(out, fn, loc)
	}
	return strings.Join(out, "\n")
}
