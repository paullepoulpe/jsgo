// Port of test/com/google/javascript/jscomp/transpile/CachingTranspilerTest.java @ closure-compiler v20261021

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

package transpile

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/bundle"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

const (
	jsaCtFOO_JS = "foo.js"
	jsaCtBAR_JS = "bar.js"
	jsaCtQUX_JS = "qux.js"
)

var (
	jsaCtRESULT1 = NewTranspileResult(jsaCtFOO_JS, "bar", "baz", "")
	jsaCtRESULT2 = NewTranspileResult(jsaCtQUX_JS, "qux", "corge", "")
	jsaCtRESULT3 = NewTranspileResult(jsaCtBAR_JS, "baz", "xyzzy", "")
)

// jsaMockTranspiler stands in for the Mockito mock of Transpiler.
type jsaMockTranspiler struct {
	results       map[cacheKey]*TranspileResult
	transpileCall map[cacheKey]int
	runtime       string
	runtimeCalls  int
}

func newJsaMockTranspiler() *jsaMockTranspiler {
	return &jsaMockTranspiler{results: map[cacheKey]*TranspileResult{}, transpileCall: map[cacheKey]int{}}
}

func (m *jsaMockTranspiler) when(path string, code string, r *TranspileResult) {
	m.results[cacheKey{path: path, code: code}] = r
}

func (m *jsaMockTranspiler) Transpile(path string, code string) *TranspileResult {
	k := cacheKey{path: path, code: code}
	m.transpileCall[k]++
	return m.results[k]
}

func (m *jsaMockTranspiler) Runtime() string {
	m.runtimeCalls++
	return m.runtime
}

func jsaCtSetUp() (*jsaMockTranspiler, Transpiler) {
	delegate := newJsaMockTranspiler()
	return delegate, NewCachingTranspiler(delegate, -1)
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
	jsaEq(t, delegate.transpileCall[cacheKey{path: "foo.js", code: "bar"}], 1)
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

func (jsaFailingTranspiler) Transpile(path string, code string) *TranspileResult {
	panic(bundle.NewTranspilationException(nil, []*jscomp.JSError{jscomp.JSError_Make(jscomp.DiagnosticType_Error("TEST_FOO", "Foo"))}, nil))
}

func (jsaFailingTranspiler) Runtime() string { return "" }

// Go-only: extra coverage — a TranspilationException from the delegate is
// propagated with the same report.
func TestCachingTranspiler_TranspilationException(t *testing.T) {
	transpiler := NewCachingTranspiler(jsaFailingTranspiler{}, -1)
	e := jsaThrows[bundle.ITranspilationException](t, func() { transpiler.Transpile(jsaCtFOO_JS, "x") })
	if e != nil {
		te := e.(*bundle.TranspilationException)
		jsaEq(t, len(te.Errors()), 1)
		jsaEq(t, te.GetMessage(), "Transpilation failed:\n\nERROR - [TEST_FOO] Foo\n")
	}
}

type jsaAtomicRuntimeTranspiler struct {
	runtimeCalls atomic.Int32
}

func (*jsaAtomicRuntimeTranspiler) Transpile(path string, code string) *TranspileResult {
	return jsaCtRESULT1
}

func (m *jsaAtomicRuntimeTranspiler) Runtime() string {
	m.runtimeCalls.Add(1)
	return "runtime_lib"
}

func TestCachingTranspiler_ConcurrentRuntimeAndBoundedEviction(t *testing.T) {
	delegate := &jsaAtomicRuntimeTranspiler{}
	transpiler := NewCachingTranspiler(delegate, -1)
	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			jsaEq(t, transpiler.Runtime(), "runtime_lib")
		}()
	}
	wg.Wait()
	jsaEq(t, int(delegate.runtimeCalls.Load()), 1)

	mock := newJsaMockTranspiler()
	mock.when(jsaCtFOO_JS, "1", jsaCtRESULT1)
	mock.when(jsaCtBAR_JS, "2", jsaCtRESULT2)
	mock.when(jsaCtQUX_JS, "3", jsaCtRESULT3)
	bounded := NewCachingTranspiler(mock, 2)
	jsaTrue(t, bounded.Transpile(jsaCtFOO_JS, "1") == jsaCtRESULT1)
	jsaTrue(t, bounded.Transpile(jsaCtBAR_JS, "2") == jsaCtRESULT2)
	jsaTrue(t, bounded.Transpile(jsaCtFOO_JS, "1") == jsaCtRESULT1)
	jsaTrue(t, bounded.Transpile(jsaCtQUX_JS, "3") == jsaCtRESULT3)
	jsaTrue(t, bounded.Transpile(jsaCtFOO_JS, "1") == jsaCtRESULT1)
	jsaTrue(t, bounded.Transpile(jsaCtBAR_JS, "2") == jsaCtRESULT2)
	jsaEq(t, mock.transpileCall[cacheKey{path: jsaCtFOO_JS, code: "1"}], 1)
	jsaEq(t, mock.transpileCall[cacheKey{path: jsaCtBAR_JS, code: "2"}], 2)
}
