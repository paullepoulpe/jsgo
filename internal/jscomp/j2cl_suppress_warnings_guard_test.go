// Port of com/google/javascript/jscomp/J2clSuppressWarningsGuardTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/J2clSuppressWarningsGuardTest.java)

/*
 * Copyright 2019 The Closure Compiler Authors.
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

// J2clSuppressWarningsGuardTest holds the tests for
// J2clSuppressWarningsGuard.
type J2clSuppressWarningsGuardTest struct {
	t *testing.T
}

func (c *J2clSuppressWarningsGuardTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *J2clSuppressWarningsGuardTest) JUnitBefore()           {}
func (c *J2clSuppressWarningsGuardTest) JUnitAfter()            {}

func TestJ2clSuppressWarningsGuardTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &J2clSuppressWarningsGuardTest{} })
}

// Java: J2clSuppressWarningsGuardTest.testSuppress_j2cl()
func (c *J2clSuppressWarningsGuardTest) TestSuppress_j2cl() {
	var guard WarningsGuard = NewJ2clSuppressWarningsGuard()

	AssertThat(guard.Level(J2clSuppressWarningsGuardTest_MakeError("hello.java.js"))).IsNull()
	AssertThat(guard.Level(J2clSuppressWarningsGuardTest_MakeJ2clSuppressedError("hello.java.js"))).IsEqualTo(CheckLevel_OFF.Ptr())
	AssertThat(guard.Level(J2clSuppressWarningsGuardTest_MakeJ2clSuppressedError("hello.impl.java.js"))).IsEqualTo(CheckLevel_OFF.Ptr())
	AssertThat(guard.Level(J2clSuppressWarningsGuardTest_MakeJ2clSuppressedError("foo/hello.impl.java.js"))).IsEqualTo(CheckLevel_OFF.Ptr())
}

// Java: J2clSuppressWarningsGuardTest.testSuppress_nonJ2cl()
func (c *J2clSuppressWarningsGuardTest) TestSuppress_nonJ2cl() {
	var guard WarningsGuard = NewJ2clSuppressWarningsGuard()

	AssertThat(guard.Level(J2clSuppressWarningsGuardTest_MakeError("hello.js"))).IsNull()
	AssertThat(guard.Level(J2clSuppressWarningsGuardTest_MakeJ2clSuppressedError("hello.js"))).IsNull()
	AssertThat(guard.Level(J2clSuppressWarningsGuardTest_MakeJ2clSuppressedError("kajava.js"))).IsNull()
}

// Java: J2clSuppressWarningsGuardTest.makeError(String)
func J2clSuppressWarningsGuardTest_MakeError(sourcePath string) *JSError {
	n := rhino.NewNode(rhino.Token_EMPTY)
	n.SetSourceFileForTesting(sourcePath)
	return JSError_MakeWithNode(n, TypeCheck_INEXISTENT_PROPERTY)
}

// Java: J2clSuppressWarningsGuardTest.makeJ2clSuppressedError(String)
func J2clSuppressWarningsGuardTest_MakeJ2clSuppressedError(sourcePath string) *JSError {
	n := rhino.NewNode(rhino.Token_EMPTY)
	n.SetSourceFileForTesting(sourcePath)
	return JSError_MakeWithNode(n, CheckSideEffects_USELESS_CODE_ERROR)
}
