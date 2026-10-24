// Port of com/google/javascript/jscomp/parsing/parser/trees/ParseTree.java @ closure-compiler v20261021

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

package parser

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ParseTree is an abstract syntax tree for JavaScript parse trees.
// Immutable.
// A plain old data structure. Should include data members and simple
// accessors only.
//
// Derived classes should have a 'Tree' suffix. Each concrete derived class
// should have a ParseTreeType whose name matches the derived class name.
//
// A parse tree derived from source should have a non-null location. A parse
// tree that is synthesized by the compiler may have a null location.
//
// When adding a new subclass of ParseTree you must also do the following:
//   - add a new entry to ParseTreeType
//   - add ParseTree.asXTree()
//
// Java: class ParseTree. Every tree embeds ParseTreeBase; fields are read
// through GetParseTreeBase(), and the asX() casts are Go type assertions
// (tree.(*XTree)).
type ParseTree interface {
	GetParseTreeBase() *ParseTreeBase
	GetStart() SourcePosition
	GetEnd() SourcePosition
	IsPattern() bool
	IsValidNonVanillaAssignmentTarget() bool
	IsValidAssignmentTarget() bool
	IsRestParameter() bool
	String() string
}

// ParseTreeBase holds ParseTree's fields and methods.
type ParseTreeBase struct {
	Type     ParseTreeType
	Location SourceRange
}

// InitParseTreeBase is the Go form of the (protected) ParseTree
// constructor.
//
// Java: ParseTree(ParseTreeType, SourceRange)
func (t *ParseTreeBase) InitParseTreeBase(type_ ParseTreeType, location SourceRange) {
	t.Type = type_
	t.Location = location
}

// GetParseTreeBase returns the base struct (field access through the
// interface).
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (t *ParseTreeBase) GetParseTreeBase() *ParseTreeBase {
	return t
}

// GetStart returns the start position.
//
// Java: ParseTree.getStart()
func (t *ParseTreeBase) GetStart() SourcePosition {
	return t.Location.Start
}

// GetEnd returns the end position.
//
// Java: ParseTree.getEnd()
func (t *ParseTreeBase) GetEnd() SourcePosition {
	return t.Location.End
}

// IsPattern reports whether the tree (inside any parentheses) is an array or
// object pattern.
//
// Java: ParseTree.isPattern() — Go-only: language-adapter — the loop over
// parentheses is ParenExpressionTree.IsPattern.
func (t *ParseTreeBase) IsPattern() bool {
	switch t.Type {
	case ParseTreeType_ARRAY_PATTERN, ParseTreeType_OBJECT_PATTERN:
		return true
	default:
		return false
	}
}

// IsValidNonVanillaAssignmentTarget: is valid assignment target for
// non-vanilla assignment operators like `+=`, `-+`, `**=`, etc
//
// Java: ParseTree.isValidNonVanillaAssignmentTarget() — Go-only:
// language-adapter — the loop over parentheses is
// ParenExpressionTree.IsValidNonVanillaAssignmentTarget.
func (t *ParseTreeBase) IsValidNonVanillaAssignmentTarget() bool {
	switch t.Type {
	case ParseTreeType_IDENTIFIER_EXPRESSION, ParseTreeType_MEMBER_EXPRESSION, ParseTreeType_MEMBER_LOOKUP_EXPRESSION, ParseTreeType_DEFAULT_PARAMETER:
		return true
	case ParseTreeType_ARRAY_PATTERN, ParseTreeType_OBJECT_PATTERN:
		return false
	default:
		return false
	}
}

// IsValidAssignmentTarget: is valid assignment target for any assignment
// operator like `=`, `+=`, `-+`, `**=`, etc
//
// Java: ParseTree.isValidAssignmentTarget() — Go-only: language-adapter —
// the loop over parentheses is ParenExpressionTree.IsValidAssignmentTarget.
func (t *ParseTreeBase) IsValidAssignmentTarget() bool {
	switch t.Type {
	case ParseTreeType_IDENTIFIER_EXPRESSION,
		ParseTreeType_MEMBER_EXPRESSION,
		ParseTreeType_MEMBER_LOOKUP_EXPRESSION,
		ParseTreeType_ARRAY_PATTERN,
		ParseTreeType_OBJECT_PATTERN,
		ParseTreeType_DEFAULT_PARAMETER:
		return true
	default:
		return false
	}
}

// IsRestParameter reports whether this is an ITER_REST tree.
//
// Java: ParseTree.isRestParameter()
func (t *ParseTreeBase) IsRestParameter() bool {
	return t.Type == ParseTreeType_ITER_REST
}

// String renders "<type>@<location>".
//
// Java: ParseTree.toString()
func (t *ParseTreeBase) String() string {
	return t.Type.String() + "@" + jutil.ToJavaString(t.Location)
}
