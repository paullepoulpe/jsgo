// Port of com/google/javascript/jscomp/SourceInformationAnnotatorTest.java @ closure-compiler v20260415

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func sourceInformationAnnotatorTestOriginalName(n *rhino.Node) string {
	if n.GetOriginalName() == nil {
		return "<null>"
	}
	return *n.GetOriginalName()
}

func TestSourceInformationAnnotator_PreserveAnnotatedName(t *testing.T) {
	root := rhino.NewNode(rhino.Token_SCRIPT)
	name := rhino.Node_NewString("foo")
	bar := "bar"
	name.SetOriginalName(&bar)
	root.AddChildToBack(name)

	NodeTraversal_Traverse(NewCompiler(), root, SourceInformationAnnotator_Create())
	if got := sourceInformationAnnotatorTestOriginalName(name); got != "bar" {
		t.Errorf("original name: got %s, want bar", got)
	}
}

func TestSourceInformationAnnotator_SetOriginalGetpropNames(t *testing.T) {
	root := rhino.NewNode(rhino.Token_SCRIPT)
	getprop := rhino.IR_Getprop(rhino.IR_Name("x"), "y")
	root.AddChildToBack(rhino.IR_ExprResult(getprop))

	NodeTraversal_Traverse(NewCompiler(), root, SourceInformationAnnotator_Create())

	if got := sourceInformationAnnotatorTestOriginalName(getprop); got != "y" {
		t.Errorf("original name: got %s, want y", got)
	}
}

func TestSourceInformationAnnotator_SetOriginalOptChainGetpropNames(t *testing.T) {
	root := rhino.NewNode(rhino.Token_SCRIPT)
	getprop := rhino.IR_StartOptChainGetprop(rhino.IR_Name("x"), "y")
	root.AddChildToBack(rhino.IR_ExprResult(getprop))

	NodeTraversal_Traverse(NewCompiler(), root, SourceInformationAnnotator_Create())

	if got := sourceInformationAnnotatorTestOriginalName(getprop); got != "y" {
		t.Errorf("original name: got %s, want y", got)
	}
}

// Java: SourceInformationAnnotatorTest.doesNotSetOriginalStringName()
func TestSourceInformationAnnotator_DoesNotSetOriginalStringName(t *testing.T) {
	root := rhino.NewNode(rhino.Token_SCRIPT)
	string_ := rhino.IR_String("x")
	root.AddChildToBack(rhino.IR_ExprResult(string_))

	NodeTraversal_Traverse(NewCompiler(), root, SourceInformationAnnotator_Create())

	// No need for the original name because strings are almost never
	// mangled by JSCompiler and source information mapping "identifier"s
	// don't care about raw strings.
	if string_.GetOriginalName() != nil {
		t.Errorf("original name: got %s, want null", *string_.GetOriginalName())
	}
}
