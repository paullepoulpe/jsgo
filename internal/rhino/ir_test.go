// Port of com/google/javascript/rhino/IRTest.java @ closure-compiler v20261015

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   John Lenz
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func testIR(t *testing.T, node Node, expectedStructure string) {
	t.Helper()
	if got := node.ToStringTree(); got != expectedStructure {
		t.Errorf("toStringTree:\n%s\nwant:\n%s", got, expectedStructure)
	}
}

func TestIR_Empty(t *testing.T) {
	testIR(t, IR_Empty(), "EMPTY\n")
}

func TestIR_NewTarget(t *testing.T) {
	testIR(t, IR_NewTarget(), "NEW_TARGET\n")
}

func TestIR_Function(t *testing.T) {
	testIR(t,
		IR_Function(IR_Name("hi"), IR_ParamList(), IR_Block()),
		"FUNCTION hi\n    NAME hi\n    PARAM_LIST\n    BLOCK\n")
}

func TestIR_ArrowFunction(t *testing.T) {
	testIR(t,
		IR_ArrowFunction(IR_Name("hi"), IR_ParamList(), IR_Block()),
		"FUNCTION hi [arrow_fn: 1]\n    NAME hi\n    PARAM_LIST\n    BLOCK\n")
}

func TestIR_ParamList(t *testing.T) {
	testIR(t, IR_ParamList(), "PARAM_LIST\n")

	testIR(t,
		IR_ParamList(IR_Name("a"), IR_Name("b")),
		"PARAM_LIST\n    NAME a\n    NAME b\n")
}

func TestIR_Block(t *testing.T) {
	testIR(t, IR_Block(), "BLOCK\n")

	testIR(t,
		IR_BlockOfNodes(IR_Empty(), IR_Empty()),
		"BLOCK\n    EMPTY\n    EMPTY\n")

	testIR(t,
		IR_BlockOfList([]Node{IR_Empty(), IR_Empty()}),
		"BLOCK\n    EMPTY\n    EMPTY\n")
}

func TestIR_Script(t *testing.T) {
	testIR(t, IR_Script(), "SCRIPT\n")

	testIR(t,
		IR_ScriptOfNodes(IR_Empty(), IR_Empty()),
		"SCRIPT\n    EMPTY\n    EMPTY\n")

	testIR(t,
		IR_ScriptOfList([]Node{IR_Empty(), IR_Empty()}),
		"SCRIPT\n    EMPTY\n    EMPTY\n")
}

func TestIR_ScriptThrows(t *testing.T) {
	caught := false
	jutil.Try(func() {
		IR_ScriptOfNodes(IR_ReturnNode())
	}, func(e jutil.IIllegalStateException) {
		caught = true
	})
	if !caught {
		t.Fatal("expected exception was not seen")
	}
}

func TestIR_Var(t *testing.T) {
	testIR(t,
		IR_Var(IR_Name("a")),
		"VAR\n    NAME a\n")

	testIR(t,
		IR_VarWithValue(IR_Name("a"), IR_TrueNode()),
		"VAR\n    NAME a\n        TRUE\n")
}

func TestIR_Return(t *testing.T) {
	testIR(t, IR_ReturnNode(), "RETURN\n")

	testIR(t,
		IR_ReturnNodeWithExpr(IR_Name("a")),
		"RETURN\n    NAME a\n")
}

func TestIR_Throw(t *testing.T) {
	testIR(t,
		IR_ThrowNode(IR_Name("a")),
		"THROW\n    NAME a\n")
}

func TestIR_ExprResult(t *testing.T) {
	testIR(t,
		IR_ExprResult(IR_Name("a")),
		"EXPR_RESULT\n    NAME a\n")
}

func TestIR_If(t *testing.T) {
	testIR(t,
		IR_IfNode(IR_Name("a"), IR_Block()),
		"IF\n    NAME a\n    BLOCK\n")

	testIR(t,
		IR_IfNodeWithElse(IR_Name("a"), IR_Block(), IR_Block()),
		"IF\n    NAME a\n    BLOCK\n    BLOCK\n")
}

func TestIR_Issue727_1(t *testing.T) {
	testIR(t,
		IR_TryFinally(IR_Block(), IR_Block()),
		"TRY\n    BLOCK\n    BLOCK\n    BLOCK\n")
}

func TestIR_Issue727_2(t *testing.T) {
	testIR(t,
		IR_TryCatch(IR_Block(), IR_CatchNode(IR_Name("e"), IR_Block())),
		"TRY\n    BLOCK\n    BLOCK\n        CATCH\n            NAME e\n            BLOCK\n")
}

func TestIR_Issue727_3(t *testing.T) {
	testIR(t,
		IR_TryCatchFinally(IR_Block(), IR_CatchNode(IR_Name("e"), IR_Block()), IR_Block()),
		"TRY\n    BLOCK\n    BLOCK\n        CATCH\n            NAME e\n            BLOCK\n    BLOCK\n")
}

func TestIR_Add(t *testing.T) {
	testIR(t,
		IR_Add(IR_Cast(IR_Number(1), nil), IR_Number(2)),
		"ADD\n    CAST\n        NUMBER 1.0\n    NUMBER 2.0\n")
}

func TestIR_VarWithTemplateLitOnRHS(t *testing.T) {
	testIR(t,
		IR_VarWithValue(IR_Name("x"), NewNodeWithChild(Token_TEMPLATELIT, IR_String(""))),
		strings.Join([]string{
			"VAR", //
			"    NAME x",
			"        TEMPLATELIT",
			"            STRINGLIT ",
			""}, "\n"))

	testIR(t,
		IR_VarWithValue(IR_Name("x"), NewNodeWith2Children(Token_TAGGED_TEMPLATELIT, IR_Name("y"), IR_String(""))),
		strings.Join([]string{
			"VAR", //
			"    NAME x",
			"        TAGGED_TEMPLATELIT",
			"            NAME y",
			"            STRINGLIT ",
			""}, "\n"))
}
