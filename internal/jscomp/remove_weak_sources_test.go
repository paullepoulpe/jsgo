// Port of com/google/javascript/jscomp/RemoveWeakSourcesTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/RemoveWeakSourcesTest.java)

/*
 * Copyright 2018 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type RemoveWeakSourcesTest struct {
	CompilerTestCaseBase
}

// Java: RemoveWeakSourcesTest()
func NewRemoveWeakSourcesTest() *RemoveWeakSourcesTest {
	c := &RemoveWeakSourcesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRemoveWeakSourcesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRemoveWeakSourcesTest() })
}

// Java: RemoveWeakSourcesTest.getProcessor(Compiler)
func (c *RemoveWeakSourcesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewRemoveWeakSources(compiler)
}

// Java: RemoveWeakSourcesTest.testSingleChunk()
func (c *RemoveWeakSourcesTest) TestSingleChunk() {
	strongSrc1 := SourceFile_FromCodeWithKind("a.js", "var a = 1;", rhino.StaticSourceFile_SourceKind_STRONG)
	strongSrc2 := SourceFile_FromCodeWithKind("b.js", "var b = 2;", rhino.StaticSourceFile_SourceKind_STRONG)
	weakSrc1 := SourceFile_FromCodeWithKind("c.js", "var c = 3;", rhino.StaticSourceFile_SourceKind_WEAK)
	weakSrc2 := SourceFile_FromCodeWithKind("d.js", "var d = 4;", rhino.StaticSourceFile_SourceKind_WEAK)
	emptyWeakSrc1 := SourceFile_FromCode("c.js", "")
	emptyWeakSrc2 := SourceFile_FromCode("d.js", "")

	// Expect the weak sources to be emptied and moved to the end.
	c.Test(
		c.Srcs(strongSrc1, weakSrc1, strongSrc2, weakSrc2),
		c.Expected(strongSrc1, strongSrc2, emptyWeakSrc1, emptyWeakSrc2))
}

// Java: RemoveWeakSourcesTest.testMultipleChunks()
func (c *RemoveWeakSourcesTest) TestMultipleChunks() {
	strongSrc1 := SourceFile_FromCodeWithKind("a.js", "var a = 1;", rhino.StaticSourceFile_SourceKind_STRONG)
	strongSrc2 := SourceFile_FromCodeWithKind("b.js", "var b = 2;", rhino.StaticSourceFile_SourceKind_STRONG)
	weakSrc1 := SourceFile_FromCodeWithKind("c.js", "var c = 3;", rhino.StaticSourceFile_SourceKind_WEAK)
	weakSrc2 := SourceFile_FromCodeWithKind("d.js", "var d = 4;", rhino.StaticSourceFile_SourceKind_WEAK)
	weakSrc3 := SourceFile_FromCodeWithKind("e.js", "var e = 5;", rhino.StaticSourceFile_SourceKind_WEAK)
	emptyWeakSrc1 := SourceFile_FromCode("c.js", "")
	emptyWeakSrc2 := SourceFile_FromCode("d.js", "")
	emptyWeakSrc3 := SourceFile_FromCode("e.js", "")
	fillFileSrc := SourceFile_FromCode("$fillFile", "")

	chunkBefore1 := NewJSChunk("m1")
	chunkBefore1.AddSourceFile(strongSrc1)
	chunkBefore1.AddSourceFile(weakSrc1)
	chunkAfter1 := NewJSChunk("m1")
	chunkAfter1.AddSourceFile(strongSrc1)

	chunkBefore2 := NewJSChunk("m2")
	chunkBefore2.AddSourceFile(weakSrc2)
	chunkBefore2.AddSourceFile(strongSrc2)
	chunkAfter2 := NewJSChunk("m2")
	chunkAfter2.AddSourceFile(strongSrc2)

	chunkBefore3 := NewJSChunk("m3")
	chunkBefore3.AddSourceFile(weakSrc3)
	chunkAfter3 := NewJSChunk("m3")
	chunkAfter3.AddSourceFile(fillFileSrc)

	weakChunk := NewJSChunk("$weak$")
	weakChunk.AddSourceFile(emptyWeakSrc1)
	weakChunk.AddSourceFile(emptyWeakSrc2)
	weakChunk.AddSourceFile(emptyWeakSrc3)

	// Expect the weak sources to be emptied and moved to a separate final chunk.
	c.Test(
		c.Srcs([]*JSChunk{chunkBefore1, chunkBefore2, chunkBefore3}),
		c.Expected([]*JSChunk{chunkAfter1, chunkAfter2, chunkAfter3, weakChunk}))
}
