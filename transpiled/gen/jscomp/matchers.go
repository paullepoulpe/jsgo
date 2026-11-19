/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/refactoring/Matchers.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.refactoring.Matchers
type Matchers struct {
	idHash_ atomic.Int32
}

func (this_ *Matchers) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers) IsJavaNil_() bool { return this_ == nil }
func (*Matchers) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers", "Matchers"
}
func (this_ *Matchers) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.Matchers", "com.google.javascript.refactoring.Matchers", func(v_ any) bool {
		return j2gort.InstanceOf[*Matchers](v_)
	})
}

func (this_ *Matchers) initFields_Matchers() {
	_ = this_
}

// Java constructor: Matchers()
func NewMatchers() *Matchers {
	this_ := &Matchers{}
	InitMatchers(this_)
	return this_
}

func InitMatchers(this_ *Matchers) {
	this_.initFields_Matchers()
}

// Java static method: Matchers.anything()
func Matchers_Anything() Matcher {
	return NewMatchers__Anon_135()
}

// Java static method: Matchers.allOf(com.google.javascript.refactoring.Matcher[] matchers)
func Matchers_AllOf(matchers []Matcher) Matcher {
	_ = matchers
	return NewMatchers__Anon_136(matchers)
}

// Java static method: Matchers.anyOf(com.google.javascript.refactoring.Matcher[] matchers)
func Matchers_AnyOf(matchers []Matcher) Matcher {
	_ = matchers
	return NewMatchers__Anon_137(matchers)
}

// Java static method: Matchers.not(com.google.javascript.refactoring.Matcher matcher)
func Matchers_Not(matcher Matcher) Matcher {
	_ = matcher
	return NewMatchers__Anon_138(matcher)
}

// Java static method: Matchers.constructor()
func Matchers_Constructor_0() Matcher {
	return Matchers_Constructor_1(j2gort.NullString)
}

// Java static method: Matchers.constructor(java.lang.@org.jspecify.annotations.Nullable String name)
func Matchers_Constructor_1(name string) Matcher {
	_ = name
	return NewMatchers__Anon_139(name)
}

// Java static method: Matchers.newClass()
func Matchers_NewClass_0() Matcher {
	return NewMatchers__Anon_140()
}

// Java static method: Matchers.newClass(java.lang.String className)
func Matchers_NewClass_1(className string) Matcher {
	_ = className
	return NewMatchers__Anon_141(className)
}

// Java static method: Matchers.functionCall()
func Matchers_FunctionCall_0() Matcher {
	return Matchers_FunctionCall_1(j2gort.NullString)
}

// Java static method: Matchers.functionCall(java.lang.@org.jspecify.annotations.Nullable String name)
func Matchers_FunctionCall_1(name string) Matcher {
	_ = name
	return NewMatchers__Anon_142(name)
}

// Java static method: Matchers.functionCallWithNumArgs(int numArgs)
func Matchers_FunctionCallWithNumArgs_1(numArgs int32) Matcher {
	_ = numArgs
	return NewMatchers__Anon_143(numArgs)
}

// Java static method: Matchers.functionCallWithNumArgs(java.lang.String name, int numArgs)
func Matchers_FunctionCallWithNumArgs_2(name string, numArgs int32) Matcher {
	_ = name
	_ = numArgs
	return Matchers_AllOf([]Matcher{Matchers_FunctionCallWithNumArgs_1(numArgs), Matchers_FunctionCall_1(name)})
}

// Java static method: Matchers.googRequirelike(java.lang.String namespace)
func Matchers_GoogRequirelike_1(namespace string) Matcher {
	_ = namespace
	return NewLambda_Matcher(func(node pkg_rhino.Node, metadata *NodeMetadata) bool {
		_ = node
		_ = metadata
		return ((Matchers_GoogRequirelike_0().Matches(node, metadata) && node.GetSecondChild().IsStringLit()) && (j2gort.CheckNotNullString(node.GetSecondChild().GetString()) == namespace))
	})
}

// Java static method: Matchers.googRequirelike()
func Matchers_GoogRequirelike_0() Matcher {
	return Matchers_AnyOf([]Matcher{Matchers_GoogRequire(), Matchers_GoogRequireType(), Matchers_GoogForwardDeclare()})
}

// Java static method: Matchers.googRequire()
func Matchers_GoogRequire() Matcher {
	return Matchers_FunctionCall_1("goog.require")
}

// Java static method: Matchers.googRequireType()
func Matchers_GoogRequireType() Matcher {
	return Matchers_FunctionCall_1("goog.requireType")
}

// Java static method: Matchers.googForwardDeclare()
func Matchers_GoogForwardDeclare() Matcher {
	return Matchers_FunctionCall_1("goog.forwardDeclare")
}

// Java static method: Matchers.googModule()
func Matchers_GoogModule() Matcher {
	return Matchers_FunctionCall_1("goog.module")
}

// Java static method: Matchers.googModuleOrProvide()
func Matchers_GoogModuleOrProvide() Matcher {
	return Matchers_AnyOf([]Matcher{Matchers_GoogModule(), Matchers_FunctionCall_1("goog.provide")})
}

// Java static method: Matchers.propertyAccess()
func Matchers_PropertyAccess_0() Matcher {
	return Matchers_PropertyAccess_1(j2gort.NullString)
}

// Java static method: Matchers.propertyAccess(java.lang.@org.jspecify.annotations.Nullable String name)
func Matchers_PropertyAccess_1(name string) Matcher {
	_ = name
	return NewMatchers__Anon_144(name)
}

// Java static method: Matchers.enumDefinition()
func Matchers_EnumDefinition() Matcher {
	return NewMatchers__Anon_145()
}

// Java static method: Matchers.enumDefinitionOfType(java.lang.String type)
func Matchers_EnumDefinitionOfType(type_ string) Matcher {
	_ = type_
	return NewMatchers__Anon_146(type_)
}

// Java static method: Matchers.assignmentWithRhs(com.google.javascript.refactoring.Matcher rhsMatcher)
func Matchers_AssignmentWithRhs(rhsMatcher Matcher) Matcher {
	_ = rhsMatcher
	return NewMatchers__Anon_147(rhsMatcher)
}

// Java static method: Matchers.prototypeVariableDeclaration()
func Matchers_PrototypeVariableDeclaration() Matcher {
	return Matchers_MatcherForPrototypeDeclaration(false)
}

// Java static method: Matchers.prototypeMethodDeclaration()
func Matchers_PrototypeMethodDeclaration() Matcher {
	return Matchers_MatcherForPrototypeDeclaration(true)
}

// Java static method: Matchers.jsDocType(java.lang.String type)
func Matchers_JsDocType(type_ string) Matcher {
	_ = type_
	return NewMatchers__Anon_148(type_)
}

// Java static method: Matchers.constructorPropertyDeclaration()
func Matchers_ConstructorPropertyDeclaration() Matcher {
	return NewMatchers__Anon_149()
}

// Java static method: Matchers.isPrivate()
func Matchers_IsPrivate() Matcher {
	return NewMatchers__Anon_150()
}

// Java static method: Matchers.getJsType(com.google.javascript.refactoring.NodeMetadata metadata, java.lang.String type)
func Matchers_GetJsType_2_NodeMetadata_string(metadata *NodeMetadata, type_ string) pkg_rhino.JSType {
	_ = metadata
	_ = type_
	return metadata.GetCompiler().GetTypeRegistry().GetGlobalType(type_)
}

// Java static method: Matchers.getJsType(com.google.javascript.refactoring.NodeMetadata metadata, com.google.javascript.rhino.jstype.JSTypeNative nativeType)
func Matchers_GetJsType_2_NodeMetadata_JSTypeNative(metadata *NodeMetadata, nativeType pkg_rhino.JSTypeNative) pkg_rhino.JSType {
	_ = metadata
	_ = nativeType
	return metadata.GetCompiler().GetTypeRegistry().GetNativeType(nativeType)
}

// Java static method: Matchers.areTypesEquivalentIgnoringGenerics(com.google.javascript.rhino.jstype.JSType a, com.google.javascript.rhino.jstype.JSType b)
func Matchers_AreTypesEquivalentIgnoringGenerics(a pkg_rhino.JSType, b pkg_rhino.JSType) bool {
	_ = a
	_ = b
	var equivalent bool = a.Equals(b)
	_ = equivalent
	if equivalent {
		return true
	}
	if a.IsTemplatizedType() {
		return a.ToMaybeTemplatizedType().GetReferencedType().Equals(b)
	}
	return false
}

// Java static method: Matchers.matchesPrototypeInstanceVar(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata, java.lang.String name)
func Matchers_MatchesPrototypeInstanceVar(node pkg_rhino.Node, metadata *NodeMetadata, name string) bool {
	_ = node
	_ = metadata
	_ = name
	var parts []string = j2gort.JStrSplit(name, ".prototype.")
	_ = parts
	var className string = parts[int32(0)]
	_ = className
	var propertyName string = parts[int32(1)]
	_ = propertyName
	var providedJsType pkg_rhino.JSType = Matchers_GetJsType_2_NodeMetadata_string(metadata, className)
	_ = providedJsType
	if providedJsType == nil || providedJsType.IsJavaNil_() {
		return false
	}
	var jsType pkg_rhino.JSType = nil
	_ = jsType
	if node.HasChildren() {
		jsType = node.GetFirstChild().GetJSType()
	}
	if jsType == nil || jsType.IsJavaNil_() {
		return false
	}
	jsType = jsType.RestrictByNotNullOrUndefined()
	if (!(jsType.IsUnknownType()) && !(jsType.IsAllType())) && jsType.IsSubtypeOf_1(providedJsType) {
		if node.IsName() && func() bool {
			sr_1 := propertyName
			se_2 := node.GetString()
			return j2gort.CheckNotNullString(sr_1) == se_2
		}() {
			return true
		} else if node.IsGetProp() && func() bool {
			sr_3 := propertyName
			se_4 := node.GetString()
			return j2gort.CheckNotNullString(sr_3) == se_4
		}() {
			return true
		}
	}
	return false
}

// Java static method: Matchers.matcherForPrototypeDeclaration(boolean requireFunctionType)
func Matchers_MatcherForPrototypeDeclaration(requireFunctionType bool) Matcher {
	_ = requireFunctionType
	return NewMatchers__Anon_151(requireFunctionType)
}

// Java class:
type Matchers__Anon_135 struct {
	idHash_ atomic.Int32
}

func (this_ *Matchers__Anon_135) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_135) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_135) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$1", ""
}
func (this_ *Matchers__Anon_135) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_135) HasIdentityEquals_()        {}
func (*Matchers__Anon_135) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_135) initFields_Matchers__Anon_135() {
	_ = this_
}

func NewMatchers__Anon_135() *Matchers__Anon_135 {
	this_ := &Matchers__Anon_135{}
	this_.initFields_Matchers__Anon_135()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_135) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	return true
}

// Java class:
type Matchers__Anon_136 struct {
	cap_0_matchers []Matcher
	idHash_        atomic.Int32
}

func (this_ *Matchers__Anon_136) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_136) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_136) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$2", ""
}
func (this_ *Matchers__Anon_136) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_136) HasIdentityEquals_()        {}
func (*Matchers__Anon_136) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_136) initFields_Matchers__Anon_136() {
	_ = this_
}

func NewMatchers__Anon_136(cap_p_matchers []Matcher) *Matchers__Anon_136 {
	this_ := &Matchers__Anon_136{}
	this_.cap_0_matchers = cap_p_matchers
	this_.initFields_Matchers__Anon_136()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_136) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	for _, el_5 := range j2gort.CheckNotNullSlice(this_.cap_0_matchers) {
		m := el_5
		_ = m
		if !(m.Matches(node, metadata)) {
			return false
		}
	}
	return true
}

// Java class:
type Matchers__Anon_137 struct {
	cap_0_matchers []Matcher
	idHash_        atomic.Int32
}

func (this_ *Matchers__Anon_137) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_137) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_137) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$3", ""
}
func (this_ *Matchers__Anon_137) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_137) HasIdentityEquals_()        {}
func (*Matchers__Anon_137) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_137) initFields_Matchers__Anon_137() {
	_ = this_
}

func NewMatchers__Anon_137(cap_p_matchers []Matcher) *Matchers__Anon_137 {
	this_ := &Matchers__Anon_137{}
	this_.cap_0_matchers = cap_p_matchers
	this_.initFields_Matchers__Anon_137()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_137) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	for _, el_6 := range j2gort.CheckNotNullSlice(this_.cap_0_matchers) {
		m := el_6
		_ = m
		if m.Matches(node, metadata) {
			return true
		}
	}
	return false
}

// Java class:
type Matchers__Anon_138 struct {
	cap_0_matcher Matcher
	idHash_       atomic.Int32
}

func (this_ *Matchers__Anon_138) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_138) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_138) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$4", ""
}
func (this_ *Matchers__Anon_138) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_138) HasIdentityEquals_()        {}
func (*Matchers__Anon_138) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_138) initFields_Matchers__Anon_138() {
	_ = this_
}

func NewMatchers__Anon_138(cap_p_matcher Matcher) *Matchers__Anon_138 {
	this_ := &Matchers__Anon_138{}
	this_.cap_0_matcher = cap_p_matcher
	this_.initFields_Matchers__Anon_138()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_138) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	return !(this_.cap_0_matcher.Matches(node, metadata))
}

// Java class:
type Matchers__Anon_139 struct {
	cap_0_name string
	idHash_    atomic.Int32
}

func (this_ *Matchers__Anon_139) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_139) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_139) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$5", ""
}
func (this_ *Matchers__Anon_139) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_139) HasIdentityEquals_()        {}
func (*Matchers__Anon_139) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_139) initFields_Matchers__Anon_139() {
	_ = this_
}

func NewMatchers__Anon_139(cap_p_name string) *Matchers__Anon_139 {
	this_ := &Matchers__Anon_139{}
	this_.cap_0_name = cap_p_name
	this_.initFields_Matchers__Anon_139()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_139) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	var info *pkg_rhino.JSDocInfo = node.GetJSDocInfo()
	_ = info
	if (info != nil) && info.IsConstructor() {
		var firstChild pkg_rhino.Node = node.GetFirstChild()
		_ = firstChild
		if this_.cap_0_name == j2gort.NullString {
			return true
		}
		if (firstChild.IsGetProp() || firstChild.IsName()) && firstChild.MatchesQualifiedName_1_string(this_.cap_0_name) {
			return true
		}
	}
	return false
}

// Java class:
type Matchers__Anon_140 struct {
	idHash_ atomic.Int32
}

func (this_ *Matchers__Anon_140) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_140) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_140) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$6", ""
}
func (this_ *Matchers__Anon_140) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_140) HasIdentityEquals_()        {}
func (*Matchers__Anon_140) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_140) initFields_Matchers__Anon_140() {
	_ = this_
}

func NewMatchers__Anon_140() *Matchers__Anon_140 {
	this_ := &Matchers__Anon_140{}
	this_.initFields_Matchers__Anon_140()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_140) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	return node.IsNew()
}

// Java class:
type Matchers__Anon_141 struct {
	cap_0_className string
	idHash_         atomic.Int32
}

func (this_ *Matchers__Anon_141) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_141) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_141) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$7", ""
}
func (this_ *Matchers__Anon_141) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_141) HasIdentityEquals_()        {}
func (*Matchers__Anon_141) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_141) initFields_Matchers__Anon_141() {
	_ = this_
}

func NewMatchers__Anon_141(cap_p_className string) *Matchers__Anon_141 {
	this_ := &Matchers__Anon_141{}
	this_.cap_0_className = cap_p_className
	this_.initFields_Matchers__Anon_141()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_141) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	if !(node.IsNew()) {
		return false
	}
	var providedJsType pkg_rhino.JSType = Matchers_GetJsType_2_NodeMetadata_string(metadata, this_.cap_0_className)
	_ = providedJsType
	if providedJsType == nil || providedJsType.IsJavaNil_() {
		return false
	}
	var jsType pkg_rhino.JSType = node.GetJSType()
	_ = jsType
	if jsType == nil || jsType.IsJavaNil_() {
		return false
	}
	jsType = jsType.RestrictByNotNullOrUndefined()
	return Matchers_AreTypesEquivalentIgnoringGenerics(jsType, providedJsType)
}

// Java class:
type Matchers__Anon_142 struct {
	cap_0_name string
	idHash_    atomic.Int32
}

func (this_ *Matchers__Anon_142) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_142) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_142) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$8", ""
}
func (this_ *Matchers__Anon_142) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_142) HasIdentityEquals_()        {}
func (*Matchers__Anon_142) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_142) initFields_Matchers__Anon_142() {
	_ = this_
}

func NewMatchers__Anon_142(cap_p_name string) *Matchers__Anon_142 {
	this_ := &Matchers__Anon_142{}
	this_.cap_0_name = cap_p_name
	this_.initFields_Matchers__Anon_142()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_142) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	return (node.IsCall() && Matchers_PropertyAccess_1(this_.cap_0_name).Matches(node.GetFirstChild(), metadata))
}

// Java class:
type Matchers__Anon_143 struct {
	cap_0_numArgs int32
	idHash_       atomic.Int32
}

func (this_ *Matchers__Anon_143) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_143) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_143) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$9", ""
}
func (this_ *Matchers__Anon_143) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_143) HasIdentityEquals_()        {}
func (*Matchers__Anon_143) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_143) initFields_Matchers__Anon_143() {
	_ = this_
}

func NewMatchers__Anon_143(cap_p_numArgs int32) *Matchers__Anon_143 {
	this_ := &Matchers__Anon_143{}
	this_.cap_0_numArgs = cap_p_numArgs
	this_.initFields_Matchers__Anon_143()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_143) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	return (node.IsCall() && node.HasXChildren((this_.cap_0_numArgs + int32(1))))
}

// Java class:
type Matchers__Anon_144 struct {
	cap_0_name string
	idHash_    atomic.Int32
}

func (this_ *Matchers__Anon_144) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_144) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_144) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$10", ""
}
func (this_ *Matchers__Anon_144) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_144) HasIdentityEquals_()        {}
func (*Matchers__Anon_144) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_144) initFields_Matchers__Anon_144() {
	_ = this_
}

func NewMatchers__Anon_144(cap_p_name string) *Matchers__Anon_144 {
	this_ := &Matchers__Anon_144{}
	this_.cap_0_name = cap_p_name
	this_.initFields_Matchers__Anon_144()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_144) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	if node.IsGetProp() {
		if this_.cap_0_name == j2gort.NullString {
			return true
		}
		if node.MatchesQualifiedName_1_string(this_.cap_0_name) {
			return true
		} else if j2gort.JStrIndexOf(this_.cap_0_name, ".prototype.") >= 0 {
			return Matchers_MatchesPrototypeInstanceVar(node, metadata, this_.cap_0_name)
		}
	}
	return false
}

// Java class:
type Matchers__Anon_145 struct {
	idHash_ atomic.Int32
}

func (this_ *Matchers__Anon_145) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_145) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_145) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$11", ""
}
func (this_ *Matchers__Anon_145) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_145) HasIdentityEquals_()        {}
func (*Matchers__Anon_145) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_145) initFields_Matchers__Anon_145() {
	_ = this_
}

func NewMatchers__Anon_145() *Matchers__Anon_145 {
	this_ := &Matchers__Anon_145{}
	this_.initFields_Matchers__Anon_145()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_145) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	var jsType pkg_rhino.JSType = node.GetJSType()
	_ = jsType
	return ((jsType != nil && !jsType.IsJavaNil_()) && jsType.IsEnumType())
}

// Java class:
type Matchers__Anon_146 struct {
	cap_0_type string
	idHash_    atomic.Int32
}

func (this_ *Matchers__Anon_146) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_146) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_146) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$12", ""
}
func (this_ *Matchers__Anon_146) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_146) HasIdentityEquals_()        {}
func (*Matchers__Anon_146) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_146) initFields_Matchers__Anon_146() {
	_ = this_
}

func NewMatchers__Anon_146(cap_p_type string) *Matchers__Anon_146 {
	this_ := &Matchers__Anon_146{}
	this_.cap_0_type = cap_p_type
	this_.initFields_Matchers__Anon_146()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_146) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	var providedJsType pkg_rhino.JSType = Matchers_GetJsType_2_NodeMetadata_string(metadata, this_.cap_0_type)
	_ = providedJsType
	if providedJsType == nil || providedJsType.IsJavaNil_() {
		return false
	}
	providedJsType = providedJsType.RestrictByNotNullOrUndefined()
	var jsType pkg_rhino.JSType = node.GetJSType()
	_ = jsType
	return (((jsType != nil && !jsType.IsJavaNil_()) && jsType.IsEnumType()) && providedJsType.Equals(jsType.ToMaybeEnumType().GetElementsType().GetPrimitiveType()))
}

// Java class:
type Matchers__Anon_147 struct {
	cap_0_rhsMatcher Matcher
	idHash_          atomic.Int32
}

func (this_ *Matchers__Anon_147) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_147) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_147) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$13", ""
}
func (this_ *Matchers__Anon_147) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_147) HasIdentityEquals_()        {}
func (*Matchers__Anon_147) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_147) initFields_Matchers__Anon_147() {
	_ = this_
}

func NewMatchers__Anon_147(cap_p_rhsMatcher Matcher) *Matchers__Anon_147 {
	this_ := &Matchers__Anon_147{}
	this_.cap_0_rhsMatcher = cap_p_rhsMatcher
	this_.initFields_Matchers__Anon_147()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_147) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	return (node.IsAssign() && this_.cap_0_rhsMatcher.Matches(node.GetLastChild(), metadata))
}

// Java class:
type Matchers__Anon_148 struct {
	cap_0_type string
	idHash_    atomic.Int32
}

func (this_ *Matchers__Anon_148) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_148) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_148) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$14", ""
}
func (this_ *Matchers__Anon_148) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_148) HasIdentityEquals_()        {}
func (*Matchers__Anon_148) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_148) initFields_Matchers__Anon_148() {
	_ = this_
}

func NewMatchers__Anon_148(cap_p_type string) *Matchers__Anon_148 {
	this_ := &Matchers__Anon_148{}
	this_.cap_0_type = cap_p_type
	this_.initFields_Matchers__Anon_148()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_148) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	var providedJsType pkg_rhino.JSType = Matchers_GetJsType_2_NodeMetadata_string(metadata, this_.cap_0_type)
	_ = providedJsType
	if providedJsType == nil || providedJsType.IsJavaNil_() {
		return false
	}
	providedJsType = providedJsType.RestrictByNotNullOrUndefined()
	var jsDoc *pkg_rhino.JSDocInfo
	_ = jsDoc
	if NodeUtil_IsNameDeclaration(node.GetParent()) {
		jsDoc = node.GetParent().GetJSDocInfo()
	} else {
		jsDoc = node.GetJSDocInfo()
	}
	var jsType pkg_rhino.JSType = node.GetJSType()
	_ = jsType
	return ((((jsDoc != nil) && jsDoc.HasType()) && (jsType != nil && !jsType.IsJavaNil_())) && providedJsType.Equals(jsType.RestrictByNotNullOrUndefined()))
}

// Java class:
type Matchers__Anon_149 struct {
	idHash_ atomic.Int32
}

func (this_ *Matchers__Anon_149) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_149) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_149) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$15", ""
}
func (this_ *Matchers__Anon_149) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_149) HasIdentityEquals_()        {}
func (*Matchers__Anon_149) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_149) initFields_Matchers__Anon_149() {
	_ = this_
}

func NewMatchers__Anon_149() *Matchers__Anon_149 {
	this_ := &Matchers__Anon_149{}
	this_.initFields_Matchers__Anon_149()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_149) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	if (!(node.IsAssign()) || !(node.GetFirstChild().IsGetProp())) || !(node.GetFirstFirstChild().IsThis()) {
		return false
	}
	for (node != 0) && !(node.IsFunction()) {
		node = node.GetParent()
	}
	if (node != 0) && node.IsFunction() {
		var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(node)
		_ = jsDoc
		if jsDoc != nil {
			return jsDoc.IsConstructor()
		}
	}
	return false
}

// Java class:
type Matchers__Anon_150 struct {
	idHash_ atomic.Int32
}

func (this_ *Matchers__Anon_150) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_150) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_150) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$16", ""
}
func (this_ *Matchers__Anon_150) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_150) HasIdentityEquals_()        {}
func (*Matchers__Anon_150) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_150) initFields_Matchers__Anon_150() {
	_ = this_
}

func NewMatchers__Anon_150() *Matchers__Anon_150 {
	this_ := &Matchers__Anon_150{}
	this_.initFields_Matchers__Anon_150()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_150) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(node)
	_ = jsDoc
	if jsDoc != nil {
		return (jsDoc.GetVisibility() == pkg_rhino.StaticGet_JSDocInfo_Visibility_PRIVATE())
	}
	return false
}

// Java class:
type Matchers__Anon_151 struct {
	idHash_                   atomic.Int32
	cap_0_requireFunctionType bool
}

func (this_ *Matchers__Anon_151) BindSelf_(s any) {
	_ = s
}
func (this_ *Matchers__Anon_151) IsJavaNil_() bool { return this_ == nil }
func (*Matchers__Anon_151) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.Matchers$17", ""
}
func (this_ *Matchers__Anon_151) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Matchers__Anon_151) HasIdentityEquals_()        {}
func (*Matchers__Anon_151) Implements_jscomp_Matcher() {}

func (this_ *Matchers__Anon_151) initFields_Matchers__Anon_151() {
	_ = this_
}

func NewMatchers__Anon_151(cap_p_requireFunctionType bool) *Matchers__Anon_151 {
	this_ := &Matchers__Anon_151{}
	this_.cap_0_requireFunctionType = cap_p_requireFunctionType
	this_.initFields_Matchers__Anon_151()
	return this_
}

// Java method: .matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *Matchers__Anon_151) Matches(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	var firstChild pkg_rhino.Node = node.GetFirstChild()
	_ = firstChild
	if (node.IsGetProp() && firstChild.IsGetProp()) && func() bool {
		sr_7 := "prototype"
		se_8 := firstChild.GetString()
		return j2gort.CheckNotNullString(sr_7) == se_8
	}() {
		var fnJsType pkg_rhino.JSType = Matchers_GetJsType_2_NodeMetadata_JSTypeNative(metadata, pkg_rhino.StaticGet_JSTypeNative_FUNCTION_FUNCTION_TYPE())
		_ = fnJsType
		var jsType pkg_rhino.JSType = node.GetJSType()
		_ = jsType
		if jsType == nil || jsType.IsJavaNil_() {
			return false
		} else if this_.cap_0_requireFunctionType {
			return jsType.CanCastTo(fnJsType)
		} else {
			return !(jsType.CanCastTo(fnJsType))
		}
	}
	return false
}
