// Port of test/com/google/javascript/jscomp/JSChunkTest.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type jsChunkTest struct {
	chunk1 *JSChunk
	chunk2 *JSChunk // depends on chunk1
	chunk3 *JSChunk // depends on chunk1
	chunk4 *JSChunk // depends on chunk2, chunk3
	chunk5 *JSChunk // depends on chunk1
}

func newJSChunkTest() *jsChunkTest {
	x := &jsChunkTest{}
	x.chunk1 = NewJSChunk("chunk1")

	x.chunk2 = NewJSChunk("chunk2")
	x.chunk2.AddDependency(x.chunk1)

	x.chunk3 = NewJSChunk("chunk3")
	x.chunk3.AddDependency(x.chunk1)

	x.chunk4 = NewJSChunk("chunk4")
	x.chunk4.AddDependency(x.chunk2)
	x.chunk4.AddDependency(x.chunk3)

	x.chunk5 = NewJSChunk("chunk5")
	x.chunk5.AddDependency(x.chunk1)
	return x
}

func jsChunkTest_assertSetEquals(t *testing.T, what string, got *jutil.LinkedHashSet[*JSChunk], want ...*JSChunk) {
	t.Helper()
	if !got.Equals(jutil.NewLinkedHashSet(want...)) {
		t.Errorf("%s: got %v, want %v", what, got, want)
	}
}

func TestJSChunk_Dependencies(t *testing.T) {
	x := newJSChunkTest()
	if !x.chunk1.GetAllDependencies().IsEmpty() {
		t.Errorf("chunk1 deps not empty")
	}
	jsChunkTest_assertSetEquals(t, "chunk2", x.chunk2.GetAllDependencies(), x.chunk1)
	jsChunkTest_assertSetEquals(t, "chunk3", x.chunk3.GetAllDependencies(), x.chunk1)
	jsChunkTest_assertSetEquals(t, "chunk4", x.chunk4.GetAllDependencies(), x.chunk1, x.chunk2, x.chunk3)

	jsChunkTest_assertSetEquals(t, "chunk1+", x.chunk1.GetThisAndAllDependencies(), x.chunk1)
	jsChunkTest_assertSetEquals(t, "chunk2+", x.chunk2.GetThisAndAllDependencies(), x.chunk1, x.chunk2)
	jsChunkTest_assertSetEquals(t, "chunk3+", x.chunk3.GetThisAndAllDependencies(), x.chunk1, x.chunk3)
	jsChunkTest_assertSetEquals(t, "chunk4+", x.chunk4.GetThisAndAllDependencies(), x.chunk1, x.chunk2, x.chunk3, x.chunk4)
}

func TestJSChunk_SortInputs(t *testing.T) {
	a := NewCompilerInput(SourceFile_FromCode("a.js", "goog.require('b');goog.require('c')"))
	b := NewCompilerInput(SourceFile_FromCode("b.js", "goog.provide('b');goog.require('d')"))
	c := NewCompilerInput(SourceFile_FromCode("c.js", "goog.provide('c');goog.require('d')"))
	d := NewCompilerInput(SourceFile_FromCode("d.js", "goog.provide('d')"))

	// Independent chunks.
	e := NewCompilerInput(SourceFile_FromCode("e.js", "goog.provide('e')"))
	f := NewCompilerInput(SourceFile_FromCode("f.js", "goog.provide('f')"))

	type list = []*CompilerInput
	jsChunkTest_assertSortedInputs(t, list{d, b, c, a}, list{a, b, c, d})
	jsChunkTest_assertSortedInputs(t, list{d, b, c, a}, list{d, b, c, a})
	jsChunkTest_assertSortedInputs(t, list{d, c, b, a}, list{d, c, b, a})
	jsChunkTest_assertSortedInputs(t, list{d, b, c, a}, list{d, a, b, c})

	jsChunkTest_assertSortedInputs(t, list{d, b, c, a, e, f}, list{a, b, c, d, e, f})
	jsChunkTest_assertSortedInputs(t, list{e, f, d, b, c, a}, list{e, f, a, b, c, d})
	jsChunkTest_assertSortedInputs(t, list{d, b, c, a, e, f}, list{a, b, c, e, d, f})
	jsChunkTest_assertSortedInputs(t, list{e, d, b, c, a, f}, list{e, a, f, b, c, d})
}

func jsChunkTest_assertSortedInputs(t *testing.T, expected []*CompilerInput, shuffled []*CompilerInput) {
	t.Helper()
	chunk := NewJSChunk("chunk")
	for _, input := range shuffled {
		input.SetChunk(nil)
		chunk.Add(input)
	}
	// new Compiler(System.err) + initCompilerOptionsIfTesting()
	options := NewCompilerOptions()
	compiler := newTestCompiler(options)
	options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	chunk.SortInputsByDeps(compiler)

	if got := chunk.GetInputs(); !slices.Equal(got, expected) {
		t.Errorf("sortInputsByDeps(%v): got %v, want %v", shuffled, got, expected)
	}
}
