// Port of com/google/javascript/jscomp/PolymerPassStaticUtilsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/PolymerPassStaticUtilsTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PolymerPassStaticUtilsTest is the unit test for PolymerPassStaticUtils.
type PolymerPassStaticUtilsTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T
}

// Java: PolymerPassStaticUtilsTest()
func NewPolymerPassStaticUtilsTest() *PolymerPassStaticUtilsTest {
	return &PolymerPassStaticUtilsTest{}
}

func TestPolymerPassStaticUtilsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPolymerPassStaticUtilsTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *PolymerPassStaticUtilsTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore: the class has no @Before method.
func (c *PolymerPassStaticUtilsTest) JUnitBefore() {}

// JUnitAfter: the class has no @After method.
func (c *PolymerPassStaticUtilsTest) JUnitAfter() {}

// Java: PolymerPassStaticUtilsTest.testGetPolymerElementType_noNativeElement()
func (c *PolymerPassStaticUtilsTest) TestGetPolymerElementType_noNativeElement() {
	def := NewPolymerClassDefinition(
		PolymerClassDefinition_DefinitionType_ObjectLiteral,
		rhino.Node{},
		rhino.Node{},
		false,
		rhino.IR_Objectlit(),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil)
	AssertThat(PolymerPassStaticUtils_GetPolymerElementType(def)).IsEqualTo("PolymerElement")
}

// Java: PolymerPassStaticUtilsTest.testGetPolymerElementType_inputBaseElement()
func (c *PolymerPassStaticUtilsTest) TestGetPolymerElementType_inputBaseElement() {
	input := "input"
	def := NewPolymerClassDefinition(
		PolymerClassDefinition_DefinitionType_ES6Class,
		rhino.Node{},
		rhino.Node{},
		false,
		rhino.IR_Objectlit(),
		nil,
		nil,
		&input,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil)
	AssertThat(PolymerPassStaticUtils_GetPolymerElementType(def)).IsEqualTo("PolymerInputElement")
}

// TODO(jlklein): Add unit tests for remaining util functions.
