// Port of com/google/javascript/jscomp/testing/JSCompCorrespondences.java @ closure-compiler v20261021

/*
 * Copyright 2010 The Closure Compiler Authors.
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
	"sync"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Well known Correspondence instances for use in tests.

// Java: JSCompCorrespondences.DIAGNOSTIC_EQUALITY
var JSCompCorrespondences_DIAGNOSTIC_EQUALITY = Correspondence_Transforming((*JSError).Type, "has diagnostic type equal to")

// Java: JSCompCorrespondences.OWNING_DIAGNOSTIC_GROUP
var JSCompCorrespondences_OWNING_DIAGNOSTIC_GROUP = Correspondence_From(
	func(actual *JSError, expected *DiagnosticGroup) bool { return expected.Matches(actual) }, "is part of diagnostic group")

// Java: JSCompCorrespondences.DESCRIPTION_EQUALITY
var JSCompCorrespondences_DESCRIPTION_EQUALITY = Correspondence_Transforming((*JSError).Description, "has description equal to")

// Java: JSCompCorrespondences.INPUT_NAME_EQUALITY
var JSCompCorrespondences_INPUT_NAME_EQUALITY = Correspondence_Transforming((*CompilerInput).GetName, "has name equal to")

// Java: JSCompCorrespondences.EQUALITY_WHEN_PARSED_AS_EXPRESSION
var JSCompCorrespondences_EQUALITY_WHEN_PARSED_AS_EXPRESSION = Correspondence_From(
	func(actual rhino.Node, expected string) bool {
		return jsCompCorrespondences_parseExpr(expected).IsEquivalentTo(actual)
	},
	"matches nodes parsed from")

// Java: JSCompCorrespondences.referenceEquality()
func JSCompCorrespondences_ReferenceEquality() *Correspondence {
	return jsCompCorrespondences_REFERENCE_EQUALITY
}

// Java: JSCompCorrespondences.REFERENCE_EQUALITY
var jsCompCorrespondences_REFERENCE_EQUALITY = Correspondence_From(func(a, b any) bool { return a == b }, "is same instance as")

// jsCompCorrespondences_COMPILER_FOR_PARSING is a compiler for parsing
// snippets of code into AST as leniently as possible.
//
// Go-only: language-adapter — created on first use, as Java initializes the
// class's static fields on first use; a package-level Go var would build a
// Compiler before any test runs.
var jsCompCorrespondences_COMPILER_FOR_PARSING = sync.OnceValue(func() *Compiler { return NewCompiler() })

// jsCompCorrespondences_parseExpr parses expr into an expression AST as
// leniently as possible.
//
// Java: JSCompCorrespondences.parseExpr(String)
func jsCompCorrespondences_parseExpr(expr string) rhino.Node {
	exprRoot := jsCompCorrespondences_COMPILER_FOR_PARSING().
		ParseOfSourceFile(SourceFile_FromCode("expr", "("+expr+")")). // SCRIPT
		GetFirstFirstChild()                                          // EXPR_RESULT > expr
	return rhino.CheckNotNullMsg(exprRoot, "Failed to parse expression")
}
