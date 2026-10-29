// Port of com/google/javascript/jscomp/parsing/parser/PredefinedName.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/Identifiers.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/ClassOrObjectElementInfo.java @ closure-compiler v20261021

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

// PredefinedName: the set of all non-keyword, non-reserved words used in
// javascript.
//
// Java: final class PredefinedName (static constants only; the private
// constructor has no Go counterpart).
const (
	PredefinedName_AS   = "as"
	PredefinedName_FROM = "from"
	PredefinedName_GET  = "get"
	PredefinedName_OF   = "of"
	PredefinedName_SET  = "set"
)

// Identifiers: JS identifier parsing utilities.
//
// Java: final class Identifiers (static methods only; the private
// constructor has no Go counterpart).

const (
	asciiFlagIdentStart uint8 = 1 << iota
	asciiFlagIdentPart
)

var asciiIdentFlags = func() [128]uint8 {
	var t [128]uint8
	for c := uint8('A'); c <= 'Z'; c++ {
		t[c] = asciiFlagIdentStart | asciiFlagIdentPart
	}
	for c := uint8('a'); c <= 'z'; c++ {
		t[c] = asciiFlagIdentStart | asciiFlagIdentPart
	}
	for c := uint8('0'); c <= '9'; c++ {
		t[c] = asciiFlagIdentPart
	}
	t['_'] = asciiFlagIdentStart | asciiFlagIdentPart
	t['$'] = asciiFlagIdentStart | asciiFlagIdentPart
	return t
}()

// Identifiers_IsIdentifierStart reports whether ch can start an identifier.
//
// Java: Identifiers.isIdentifierStart(char)
func Identifiers_IsIdentifierStart(ch uint16) bool {
	// Most code is written in pure ASCII, so create a fast path here.
	if ch <= 127 {
		return asciiIdentFlags[ch]&asciiFlagIdentStart != 0
	}

	// Handle non-ASCII characters.
	// TODO(tjgq): This should include all characters with the ID_Start property.
	return jutil.CharacterIsLetter(ch)
}

// Identifiers_IsIdentifierPart reports whether ch can be part of an
// identifier.
//
// Java: Identifiers.isIdentifierPart(char)
func Identifiers_IsIdentifierPart(ch uint16) bool {
	// Most code is written in pure ASCII, so create a fast path here.
	if ch <= 127 {
		return asciiIdentFlags[ch]&asciiFlagIdentPart != 0
	}

	// Handle non-ASCII characters.
	// TODO(tjgq): This should include all characters with the ID_Continue property, plus
	// Zero Width Non-Joiner and Zero Width Joiner.
	return Identifiers_IsIdentifierStart(ch) || jutil.CharacterIsDigit(ch)
}

// ClassOrObjectElementInfo gathers information for ES6 class / Object
// literal methods, getters, setters and fields.
//
// Gathers: the start position, whether it is a class element (false for
// object literal), whether it is static (only for class elements), and,
// after construction, the name IdentifierToken or ParseTree.
//
// Used by Parser but separated out to enforce visibility.
type ClassOrObjectElementInfo struct {
	// The start position of the element.
	start SourcePosition

	// Whether this is a class member. False for object literal elements.
	isClassMember bool

	// Whether this class element is static. Can only be true if
	// isClassMember is true.
	isStatic bool

	name     Token
	nameExpr ParseTree
}

// NewClassOrObjectElementInfo constructs a ClassOrObjectElementInfo.
//
// Java: ClassOrObjectElementInfo(SourcePosition, boolean, boolean)
func NewClassOrObjectElementInfo(start SourcePosition, isClassMember, isStatic bool) *ClassOrObjectElementInfo {
	return &ClassOrObjectElementInfo{
		start:         start,
		isClassMember: isClassMember,
		isStatic:      isStatic,
	}
}

// ClassOrObjectElementInfo_CreateClassMemberInfo creates the info for a
// class member.
//
// Java: ClassOrObjectElementInfo.createClassMemberInfo(SourcePosition, boolean)
func ClassOrObjectElementInfo_CreateClassMemberInfo(start SourcePosition, isStatic bool) *ClassOrObjectElementInfo {
	return NewClassOrObjectElementInfo(start /* isClassMember= */, true, isStatic)
}

// ClassOrObjectElementInfo_CreateObjectLiteralElementInfo creates the info
// for an object literal element.
//
// Java: ClassOrObjectElementInfo.createObjectLiteralElementInfo(SourcePosition)
func ClassOrObjectElementInfo_CreateObjectLiteralElementInfo(start SourcePosition) *ClassOrObjectElementInfo {
	return NewClassOrObjectElementInfo(start /* isClassMember= */, false /* isStatic= */, false)
}

// SetName sets the name token.
//
// Java: ClassOrObjectElementInfo.setName(IdentifierToken)
func (c *ClassOrObjectElementInfo) SetName(name Token) {
	jutil.CheckState(c.nameExpr == nil)
	c.name = name
}

// HasName reports whether a name token was set.
//
// Java: ClassOrObjectElementInfo.hasName()
func (c *ClassOrObjectElementInfo) HasName() bool {
	return c.name.IsPresent()
}

// GetName returns the name token.
//
// Java: ClassOrObjectElementInfo.getName()
func (c *ClassOrObjectElementInfo) GetName() Token {
	if c.name.IsNull() {
		panic(jutil.NewNullPointerException())
	}
	return c.name
}

// SetNameExpr sets the name expression.
//
// Java: ClassOrObjectElementInfo.setNameExpr(ParseTree)
func (c *ClassOrObjectElementInfo) SetNameExpr(nameExpr ParseTree) {
	jutil.CheckState(c.name.IsNull())
	c.nameExpr = nameExpr
}

// HasNameExpr reports whether a name expression was set.
//
// Java: ClassOrObjectElementInfo.hasNameExpr()
func (c *ClassOrObjectElementInfo) HasNameExpr() bool {
	return c.nameExpr != nil
}

// GetNameExpr returns the name expression.
//
// Java: ClassOrObjectElementInfo.getNameExpr()
func (c *ClassOrObjectElementInfo) GetNameExpr() ParseTree {
	if c.nameExpr == nil {
		panic(jutil.NewNullPointerException())
	}
	return c.nameExpr
}
