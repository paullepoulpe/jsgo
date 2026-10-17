// Port of com/google/javascript/jscomp/SourceInformationAnnotatorTest.java @ closure-compiler v20261015

/*
 * Copyright 2011 The Closure Compiler Authors.
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

// SourceInformationAnnotatorTest is the JUnit 4 class
// SourceInformationAnnotatorTest (run with RunJUnit4, which reports a
// panic as the failure of one test method).
type SourceInformationAnnotatorTest struct {
	t *testing.T
}

func (c *SourceInformationAnnotatorTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *SourceInformationAnnotatorTest) JUnitBefore()           {}
func (c *SourceInformationAnnotatorTest) JUnitAfter()            {}

func TestSourceInformationAnnotatorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &SourceInformationAnnotatorTest{} })
}

func sourceInformationAnnotatorTestOriginalName(n *rhino.Node) string {
	if n.GetOriginalName() == nil {
		return "<null>"
	}
	return *n.GetOriginalName()
}

func (c *SourceInformationAnnotatorTest) TestPreserveAnnotatedName() {
	root := rhino.NewNode(rhino.Token_SCRIPT)
	name := rhino.Node_NewString("foo")
	bar := "bar"
	name.SetOriginalName(&bar)
	root.AddChildToBack(name)

	NodeTraversal_Traverse(NewCompiler(), root, SourceInformationAnnotator_Create())
	if got := sourceInformationAnnotatorTestOriginalName(name); got != "bar" {
		c.t.Errorf("original name: got %s, want bar", got)
	}
}

func (c *SourceInformationAnnotatorTest) TestSetOriginalGetpropNames() {
	root := rhino.NewNode(rhino.Token_SCRIPT)
	getprop := rhino.IR_Getprop(rhino.IR_Name("x"), "y")
	root.AddChildToBack(rhino.IR_ExprResult(getprop))

	NodeTraversal_Traverse(NewCompiler(), root, SourceInformationAnnotator_Create())

	if got := sourceInformationAnnotatorTestOriginalName(getprop); got != "y" {
		c.t.Errorf("original name: got %s, want y", got)
	}
}

func (c *SourceInformationAnnotatorTest) TestSetOriginalOptChainGetpropNames() {
	root := rhino.NewNode(rhino.Token_SCRIPT)
	getprop := rhino.IR_StartOptChainGetprop(rhino.IR_Name("x"), "y")
	root.AddChildToBack(rhino.IR_ExprResult(getprop))

	NodeTraversal_Traverse(NewCompiler(), root, SourceInformationAnnotator_Create())

	if got := sourceInformationAnnotatorTestOriginalName(getprop); got != "y" {
		c.t.Errorf("original name: got %s, want y", got)
	}
}

// Java: SourceInformationAnnotatorTest.doesNotSetOriginalStringName()
func (c *SourceInformationAnnotatorTest) TestDoesNotSetOriginalStringName() {
	root := rhino.NewNode(rhino.Token_SCRIPT)
	string_ := rhino.IR_String("x")
	root.AddChildToBack(rhino.IR_ExprResult(string_))

	NodeTraversal_Traverse(NewCompiler(), root, SourceInformationAnnotator_Create())

	// No need for the original name because strings are almost never
	// mangled by JSCompiler and source information mapping "identifier"s
	// don't care about raw strings.
	if string_.GetOriginalName() != nil {
		c.t.Errorf("original name: got %s, want null", *string_.GetOriginalName())
	}
}
