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
// Port of com/google/javascript/jscomp/parsing/parser/ClassOrObjectElementInfo.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.parser.ClassOrObjectElementInfo
type ClassOrObjectElementInfo struct {
	name          Token          `gson:"name"`
	nameExpr      ParseTree      `gson:"nameExpr"`
	start         SourcePosition `gson:"start"`
	idHash_       atomic.Int32
	isClassMember bool `gson:"isClassMember"`
	isStatic      bool `gson:"isStatic"`
}

func (this_ *ClassOrObjectElementInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *ClassOrObjectElementInfo) IsJavaNil_() bool { return this_ == nil }
func (*ClassOrObjectElementInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.ClassOrObjectElementInfo", "ClassOrObjectElementInfo"
}
func (this_ *ClassOrObjectElementInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClassOrObjectElementInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.ClassOrObjectElementInfo", "com.google.javascript.jscomp.parsing.parser.ClassOrObjectElementInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*ClassOrObjectElementInfo](v_)
	})
}

func (this_ *ClassOrObjectElementInfo) initFields_ClassOrObjectElementInfo() {
	_ = this_
}

// Java constructor: ClassOrObjectElementInfo(com.google.javascript.jscomp.parsing.parser.util.SourcePosition start, boolean isClassMember, boolean isStatic)
func NewClassOrObjectElementInfo(start SourcePosition, isClassMember bool, isStatic bool) *ClassOrObjectElementInfo {
	this_ := &ClassOrObjectElementInfo{}
	InitClassOrObjectElementInfo(this_, start, isClassMember, isStatic)
	return this_
}

func InitClassOrObjectElementInfo(this_ *ClassOrObjectElementInfo, start SourcePosition, isClassMember bool, isStatic bool) {
	_ = start
	_ = isClassMember
	_ = isStatic
	this_.initFields_ClassOrObjectElementInfo()
	this_.start = start
	this_.isClassMember = isClassMember
	this_.isStatic = isStatic
}

// Java static method: ClassOrObjectElementInfo.createClassMemberInfo(com.google.javascript.jscomp.parsing.parser.util.SourcePosition start, boolean isStatic)
func ClassOrObjectElementInfo_CreateClassMemberInfo(start SourcePosition, isStatic bool) *ClassOrObjectElementInfo {
	_ = start
	_ = isStatic
	return NewClassOrObjectElementInfo(start, true, isStatic)
}

// Java static method: ClassOrObjectElementInfo.createObjectLiteralElementInfo(com.google.javascript.jscomp.parsing.parser.util.SourcePosition start)
func ClassOrObjectElementInfo_CreateObjectLiteralElementInfo(start SourcePosition) *ClassOrObjectElementInfo {
	_ = start
	return NewClassOrObjectElementInfo(start, false, false)
}

// Java method: ClassOrObjectElementInfo.setName(com.google.javascript.jscomp.parsing.parser.Token name)
func (this_ *ClassOrObjectElementInfo) SetName(name Token) {
	_ = name
	if this_.nameExpr != nil && !this_.nameExpr.IsJavaNil_() {
		j2gort.CheckState(false)
	}
	this_.name = name
}

// Java method: ClassOrObjectElementInfo.hasName()
func (this_ *ClassOrObjectElementInfo) HasName() bool {
	return !(this_.name).IsJavaNil_()
}

// Java method: ClassOrObjectElementInfo.getName()
func (this_ *ClassOrObjectElementInfo) GetName() Token {
	if (this_.name).IsJavaNil_() {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException")))
	}
	return this_.name
}

// Java method: ClassOrObjectElementInfo.setNameExpr(com.google.javascript.jscomp.parsing.parser.trees.ParseTree nameExpr)
func (this_ *ClassOrObjectElementInfo) SetNameExpr(nameExpr ParseTree) {
	_ = nameExpr
	if !((this_.name).IsJavaNil_()) {
		j2gort.CheckState(false)
	}
	this_.nameExpr = nameExpr
}

// Java method: ClassOrObjectElementInfo.hasNameExpr()
func (this_ *ClassOrObjectElementInfo) HasNameExpr() bool {
	return this_.nameExpr != nil && !this_.nameExpr.IsJavaNil_()
}

// Java method: ClassOrObjectElementInfo.getNameExpr()
func (this_ *ClassOrObjectElementInfo) GetNameExpr() ParseTree {
	return j2gort.CheckNotNullIface(this_.nameExpr)
}
