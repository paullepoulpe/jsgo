// Port of com/google/javascript/jscomp/parsing/parser/ClassOrObjectElementInfo.java @ closure-compiler v20261015

/*
 * Copyright 2025 The Closure Compiler Authors.
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

	name     *IdentifierToken
	nameExpr ParseTree
}

// NewClassOrObjectElementInfo constructs a ClassOrObjectElementInfo.
//
// Java: ClassOrObjectElementInfo(SourcePosition, boolean, boolean)
func NewClassOrObjectElementInfo(start SourcePosition, isClassMember, isStatic bool) *ClassOrObjectElementInfo {
	c := &ClassOrObjectElementInfo{}
	c.start = start
	c.isClassMember = isClassMember
	c.isStatic = isStatic
	return c
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
func (c *ClassOrObjectElementInfo) SetName(name *IdentifierToken) {
	jutil.CheckState(c.nameExpr == nil)
	c.name = name
}

// HasName reports whether a name token was set.
//
// Java: ClassOrObjectElementInfo.hasName()
func (c *ClassOrObjectElementInfo) HasName() bool {
	return c.name != nil
}

// GetName returns the name token.
//
// Java: ClassOrObjectElementInfo.getName()
func (c *ClassOrObjectElementInfo) GetName() *IdentifierToken {
	return jutil.CheckNotNull(c.name)
}

// SetNameExpr sets the name expression.
//
// Java: ClassOrObjectElementInfo.setNameExpr(ParseTree)
func (c *ClassOrObjectElementInfo) SetNameExpr(nameExpr ParseTree) {
	jutil.CheckState(c.name == nil)
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
	return jutil.CheckNotNull(c.nameExpr)
}
