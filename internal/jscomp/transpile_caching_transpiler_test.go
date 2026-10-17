// Port of test/com/google/javascript/jscomp/transpile/CachingTranspilerTest.java @ closure-compiler v20261015

/*
 * Copyright 2016 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

var (
	jsaCtFOO_JS  = jutil.NewURI("foo.js")
	jsaCtBAR_JS  = jutil.NewURI("bar.js")
	jsaCtQUX_JS  = jutil.NewURI("qux.js")
	jsaCtRESULT1 = NewTranspileResult(jsaCtFOO_JS, "bar", "baz", "")
	jsaCtRESULT2 = NewTranspileResult(jsaCtQUX_JS, "qux", "corge", "")
	jsaCtRESULT3 = NewTranspileResult(jsaCtBAR_JS, "baz", "xyzzy", "")
)

// jsaMockTranspiler stands in for the Mockito mock of Transpiler.
type jsaMockTranspiler struct {
	results       map[string]*TranspileResult
	transpileCall map[string]int
	runtime       string
	runtimeCalls  int
}

func newJsaMockTranspiler() *jsaMockTranspiler {
	return &jsaMockTranspiler{results: map[string]*TranspileResult{}, transpileCall: map[string]int{}}
}

func (m *jsaMockTranspiler) when(path *jutil.URI, code string, r *TranspileResult) {
	m.results[path.String()+"\x00"+code] = r
}

func (m *jsaMockTranspiler) Transpile(path *jutil.URI, code string) *TranspileResult {
	k := path.String() + "\x00" + code
	m.transpileCall[k]++
	return m.results[k]
}

func (m *jsaMockTranspiler) Runtime() string {
	m.runtimeCalls++
	return m.runtime
}

func jsaCtSetUp() (*jsaMockTranspiler, Transpile_Transpiler) {
	delegate := newJsaMockTranspiler()
	return delegate, NewCachingTranspiler(delegate, jutil.CacheBuilder_NewBuilder())
}

func TestCachingTranspiler_TranspileDelegates(t *testing.T) {
	delegate, transpiler := jsaCtSetUp()
	delegate.when(jsaCtFOO_JS, "bar", jsaCtRESULT1)
	jsaTrue(t, transpiler.Transpile(jsaCtFOO_JS, "bar") == jsaCtRESULT1)
}

func TestCachingTranspiler_TranspileCaches(t *testing.T) {
	delegate, transpiler := jsaCtSetUp()
	delegate.when(jsaCtFOO_JS, "bar", jsaCtRESULT1)
	jsaTrue(t, transpiler.Transpile(jsaCtFOO_JS, "bar") == jsaCtRESULT1)
	jsaTrue(t, transpiler.Transpile(jsaCtFOO_JS, "bar") == jsaCtRESULT1)
	jsaEq(t, delegate.transpileCall["foo.js\x00bar"], 1)
}

func TestCachingTranspiler_TranspileDependsOnBothPathAndCode(t *testing.T) {
	delegate, transpiler := jsaCtSetUp()
	delegate.when(jsaCtFOO_JS, "bar", jsaCtRESULT1)
	delegate.when(jsaCtBAR_JS, "bar", jsaCtRESULT2)
	delegate.when(jsaCtFOO_JS, "bard", jsaCtRESULT3)
	jsaTrue(t, transpiler.Transpile(jsaCtFOO_JS, "bar") == jsaCtRESULT1)
	jsaTrue(t, transpiler.Transpile(jsaCtBAR_JS, "bar") == jsaCtRESULT2)
	jsaTrue(t, transpiler.Transpile(jsaCtFOO_JS, "bard") == jsaCtRESULT3)
}

func TestCachingTranspiler_RuntimeDelegates(t *testing.T) {
	delegate, transpiler := jsaCtSetUp()
	delegate.runtime = "xyzzy"
	jsaEq(t, transpiler.Runtime(), "xyzzy")
}

func TestCachingTranspiler_RuntimeCaches(t *testing.T) {
	delegate, transpiler := jsaCtSetUp()
	delegate.runtime = "xyzzy"
	jsaEq(t, transpiler.Runtime(), "xyzzy")
	jsaEq(t, transpiler.Runtime(), "xyzzy")
	jsaEq(t, delegate.runtimeCalls, 1)
}

// jsaFailingTranspiler throws a TranspilationException.
type jsaFailingTranspiler struct{}

func (jsaFailingTranspiler) Transpile(path *jutil.URI, code string) *TranspileResult {
	panic(NewTranspilationException(nil, []*JSError{JSError_Make(jsaSemFOO_TYPE)}, nil))
}

func (jsaFailingTranspiler) Runtime() string { return "" }

// Go-only: extra coverage — a TranspilationException from the delegate is
// rethrown (wrapping the UncheckedExecutionException) with the same report.
func TestCachingTranspiler_TranspilationException(t *testing.T) {
	transpiler := NewCachingTranspiler(jsaFailingTranspiler{}, jutil.CacheBuilder_NewBuilder())
	e := jsaThrows[ITranspilationException](t, func() { transpiler.Transpile(jsaCtFOO_JS, "x") })
	if e != nil {
		te := e.(*TranspilationException)
		jsaEq(t, len(te.Errors()), 1)
		jsaEq(t, te.GetMessage(), "Transpilation failed:\n\nERROR - [TEST_FOO] Foo\n")
		_, ok := te.GetCause().(jutil.IUncheckedExecutionException)
		jsaTrue(t, ok)
	}
}
