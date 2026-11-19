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
// Port of com/google/javascript/jscomp/J2clPropertyInlinerPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.J2clPropertyInlinerPass
type J2clPropertyInlinerPass struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *J2clPropertyInlinerPass) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPropertyInlinerPass) IsJavaNil_() bool { return this_ == nil }
func (*J2clPropertyInlinerPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPropertyInlinerPass", "J2clPropertyInlinerPass"
}
func (this_ *J2clPropertyInlinerPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*J2clPropertyInlinerPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPropertyInlinerPass", "com.google.javascript.jscomp.J2clPropertyInlinerPass", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPropertyInlinerPass](v_)
	})
}
func (*J2clPropertyInlinerPass) Implements_jscomp_CompilerPass() {}

func (this_ *J2clPropertyInlinerPass) initFields_J2clPropertyInlinerPass() {
	_ = this_
}

// Java constructor: J2clPropertyInlinerPass(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewJ2clPropertyInlinerPass(compiler AbstractCompiler) *J2clPropertyInlinerPass {
	this_ := &J2clPropertyInlinerPass{}
	InitJ2clPropertyInlinerPass(this_, compiler)
	return this_
}

func InitJ2clPropertyInlinerPass(this_ *J2clPropertyInlinerPass, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_J2clPropertyInlinerPass()
	this_.compiler = compiler
}

// Java method: J2clPropertyInlinerPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *J2clPropertyInlinerPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(J2clSourceFileChecker_ShouldRunJ2clPasses(this_.compiler)) {
		return
	}
	{
		var val_1 J2clPropertyInlinerPass_StaticFieldGetterSetterInliner
		val_1.outer_ = this_
		InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner(&val_1, root)
		(&val_1).priv_StaticFieldGetterSetterInliner_run_3484()
	}
	GatherGetterAndSetterProperties_Update(this_.compiler, externs, root)
}

// Java class: com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner struct {
	outer_  *J2clPropertyInlinerPass
	root    pkg_rhino.Node `gson:"root"`
	idHash_ atomic.Int32
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) IsJavaNil_() bool {
	return this_ == nil
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner", "StaticFieldGetterSetterInliner"
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner", "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner](v_)
	})
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) SetOuter_(o *J2clPropertyInlinerPass) {
	this_.outer_ = o
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) GetOuter_() *J2clPropertyInlinerPass {
	return this_.outer_
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner() {
	_ = this_
}

// Java constructor: StaticFieldGetterSetterInliner(com.google.javascript.rhino.Node root)
func NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner(outer_ *J2clPropertyInlinerPass, root pkg_rhino.Node) *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner {
	this_ := &J2clPropertyInlinerPass_StaticFieldGetterSetterInliner{}
	this_.outer_ = outer_
	InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner(this_, root)
	return this_
}

func InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner(this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner, root pkg_rhino.Node) {
	_ = root
	this_.initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner()
	this_.root = root
}

// Java method: StaticFieldGetterSetterInliner.run()
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) priv_StaticFieldGetterSetterInliner_run_3484() {
	var gatherer *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters = NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters(this_)
	_ = gatherer
	NodeTraversal_Traverse(this_.outer_.compiler, this_.root, gatherer)
	var result j2gort.Ext_java_util_Map = gatherer.priv_GatherJ2CLClassGetterSetters_getResults_3490()
	_ = result
	NodeTraversal_Traverse(this_.outer_.compiler, this_.root, NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties(result))
	NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass(this_, result).priv_InlinePropertiesPass_run_3493()
}

// Java method: StaticFieldGetterSetterInliner.matchesJ2clGetKeySignature(java.lang.String className, com.google.javascript.rhino.Node getKey)
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) priv_StaticFieldGetterSetterInliner_matchesJ2clGetKeySignature_3485(className string, getKey pkg_rhino.Node) bool {
	_ = className
	_ = getKey
	if !(getKey.GetFirstChild().IsFunction()) {
		return false
	}
	var getFunction pkg_rhino.Node = getKey.GetFirstChild()
	_ = getFunction
	if !(getFunction.HasChildren()) || !(getFunction.GetLastChild().IsBlock()) {
		return false
	}
	var getBlock pkg_rhino.Node = getFunction.GetLastChild()
	_ = getBlock
	if !(getBlock.HasOneChild()) {
		return false
	}
	var returnStatement pkg_rhino.Node = getBlock.GetFirstChild()
	_ = returnStatement
	if !(returnStatement.IsReturn()) {
		return false
	}
	var commaExpr pkg_rhino.Node = returnStatement.GetOnlyChild()
	_ = commaExpr
	if !(commaExpr.IsComma()) {
		return false
	}
	var callNode pkg_rhino.Node = commaExpr.GetFirstChild()
	_ = callNode
	if !(callNode.IsCall() && callNode.HasOneChild()) {
		return false
	}
	var clinitQname pkg_rhino.Node = callNode.GetOnlyChild()
	_ = clinitQname
	if !(clinitQname.IsQualifiedName()) {
		return false
	}
	var clinitQnameString string = clinitQname.GetQualifiedName()
	_ = clinitQnameString
	if !(this_.LooksLikeAJ2clGeneratedClassClinitName(className, clinitQnameString)) {
		return false
	}
	var propReference pkg_rhino.Node = commaExpr.GetSecondChild()
	_ = propReference
	if !(propReference.IsQualifiedName()) {
		return false
	}
	var propReferenceString string = propReference.GetQualifiedName()
	_ = propReferenceString
	if !(this_.LooksLikeAJ2clGeneratedStaticPropertyName(className, propReferenceString)) {
		return false
	}
	return true
}

// Java method: StaticFieldGetterSetterInliner.matchesJ2clSetKeySignature(java.lang.String className, com.google.javascript.rhino.Node setKey)
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) priv_StaticFieldGetterSetterInliner_matchesJ2clSetKeySignature_3486(className string, setKey pkg_rhino.Node) bool {
	_ = className
	_ = setKey
	if !(setKey.GetFirstChild().IsFunction()) {
		return false
	}
	var setFunction pkg_rhino.Node = setKey.GetFirstChild()
	_ = setFunction
	if (!(setFunction.HasChildren()) || !(setFunction.GetLastChild().IsBlock())) || !(setFunction.GetSecondChild().IsParamList()) {
		return false
	}
	if !(setFunction.GetSecondChild().HasOneChild()) {
		return false
	}
	var setBlock pkg_rhino.Node = setFunction.GetLastChild()
	_ = setBlock
	if !(setBlock.HasOneChild()) {
		return false
	}
	var exprResult pkg_rhino.Node = setBlock.GetOnlyChild()
	_ = exprResult
	if !(exprResult.IsExprResult()) {
		return false
	}
	var commaExpr pkg_rhino.Node = exprResult.GetOnlyChild()
	_ = commaExpr
	if !(commaExpr.IsComma()) {
		return false
	}
	var clinitCall pkg_rhino.Node = commaExpr.GetFirstChild()
	_ = clinitCall
	if !(clinitCall.IsCall() && clinitCall.HasOneChild()) {
		return false
	}
	var clinitQname pkg_rhino.Node = clinitCall.GetOnlyChild()
	_ = clinitQname
	if !(clinitQname.IsQualifiedName()) {
		return false
	}
	var clinitQnameString string = clinitQname.GetQualifiedName()
	_ = clinitQnameString
	if !(this_.LooksLikeAJ2clGeneratedClassClinitName(className, clinitQnameString)) {
		return false
	}
	var assign pkg_rhino.Node = commaExpr.GetSecondChild()
	_ = assign
	if !(assign.IsAssign()) {
		return false
	}
	var lhsQname pkg_rhino.Node = assign.GetFirstChild()
	_ = lhsQname
	if !(lhsQname.IsQualifiedName()) {
		return false
	}
	var lhsQnameString string = lhsQname.GetQualifiedName()
	_ = lhsQnameString
	if !(this_.LooksLikeAJ2clGeneratedStaticPropertyName(className, lhsQnameString)) {
		return false
	}
	return true
}

// Java method: StaticFieldGetterSetterInliner.looksLikeAJ2clGeneratedClassClinitName(java.lang.String className, java.lang.String candidate)
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) LooksLikeAJ2clGeneratedClassClinitName(className string, candidate string) bool {
	_ = className
	_ = candidate
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(candidate), j2gort.CheckNotNullString(className))) {
		return false
	}
	var candidateAfterClassName string = j2gort.JStrSubFrom(candidate, j2gort.JStrLen(className))
	_ = candidateAfterClassName
	return ((j2gort.CheckNotNullString(candidateAfterClassName) == "$$0clinit") || (j2gort.CheckNotNullString(candidateAfterClassName) == ".$clinit"))
}

// Java method: StaticFieldGetterSetterInliner.looksLikeAJ2clGeneratedStaticPropertyName(java.lang.String className, java.lang.String candidate)
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) LooksLikeAJ2clGeneratedStaticPropertyName(className string, candidate string) bool {
	_ = className
	_ = candidate
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(candidate), j2gort.CheckNotNullString(className))) {
		return false
	}
	var candidateAfterClassName string = j2gort.JStrSubFrom(candidate, j2gort.JStrLen(className))
	_ = candidateAfterClassName
	return (javalang.StringStartsWith(j2gort.CheckNotNullString(candidateAfterClassName), "$$0") || javalang.StringStartsWith(j2gort.CheckNotNullString(candidateAfterClassName), ".$"))
}

// Java polymorphic class interface: com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clProperty
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty interface {
	Implements_jscomp_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty()
	IsJavaNil_() bool
	GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase() *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase
	Remove()
}

func ToIface_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty[T interface {
	J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty
	comparable
}](v T) J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty[T J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty](v T) J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty {
	var vi J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clProperty
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase struct {
	getKey         pkg_rhino.Node `gson:"getKey"`
	setKey         pkg_rhino.Node `gson:"setKey"`
	idHash_        atomic.Int32
	isSafeToInline bool `gson:"isSafeToInline"`
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase) IsJavaNil_() bool {
	return this_ == nil
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner$J2clProperty", "J2clProperty"
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clProperty", "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner$J2clProperty", func(v_ any) bool {
		return j2gort.InstanceOf[J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty](v_)
	})
}
func (b *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase) GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase() *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase {
	return b
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase) Implements_jscomp_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty() {
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase) initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty() {
	_ = this_
}

func InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty(this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase, getKey pkg_rhino.Node, setKey pkg_rhino.Node) {
	_ = getKey
	_ = setKey
	this_.initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty()
	this_.getKey = getKey
	this_.setKey = setKey
	this_.isSafeToInline = true
}

// Java method: J2clProperty.remove()
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase) Remove() {
	panic("abstract or native method invoked")
}

// Java class: com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clPropertyEs5
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5 struct {
	J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase
	outer_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5) IsJavaNil_() bool {
	return this_ == nil
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner$J2clPropertyEs5", "J2clPropertyEs5"
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase.IdentityHash_()
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clPropertyEs5", "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner$J2clPropertyEs5", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5](v_)
	})
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5) SetOuter_(o *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) {
	this_.outer_ = o
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5) GetOuter_() *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner {
	return this_.outer_
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5) Implements_jscomp_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty() {
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5) initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5() {
	_ = this_
}

// Java constructor: J2clPropertyEs5(com.google.javascript.rhino.Node getKey, com.google.javascript.rhino.Node setKey)
func NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5(outer_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner, getKey pkg_rhino.Node, setKey pkg_rhino.Node) *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5 {
	this_ := &J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5{}
	this_.outer_ = outer_
	InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5(this_, getKey, setKey)
	return this_
}

func InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5(this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5, getKey pkg_rhino.Node, setKey pkg_rhino.Node) {
	_ = getKey
	_ = setKey
	InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty(&this_.J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase, getKey, setKey)
	this_.initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5()
	if !(getKey.IsStringKey() && (j2gort.CheckNotNullString(getKey.GetString()) == "get")) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(getKey))
	}
	if !((setKey == 0) || (setKey.IsStringKey() && (j2gort.CheckNotNullString(setKey.GetString()) == "set"))) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(setKey))
	}
}

// Java method: J2clPropertyEs5.remove()
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5) Remove() {
	var nodeToDetach pkg_rhino.Node = this_.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().getKey.GetGrandparent()
	_ = nodeToDetach
	var objectLit pkg_rhino.Node = nodeToDetach.GetParent()
	_ = objectLit
	if !(objectLit.IsObjectLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(objectLit))
	}
	nodeToDetach.Detach()
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(nodeToDetach, this_.outer_.outer_.compiler)
	this_.outer_.outer_.compiler.ReportChangeToEnclosingScope(objectLit)
	if !(objectLit.HasChildren()) {
		objectLit.GetGrandparent().Detach()
	}
}

// Java class: com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clPropertyEs6
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6 struct {
	J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase
	outer_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6) IsJavaNil_() bool {
	return this_ == nil
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner$J2clPropertyEs6", "J2clPropertyEs6"
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase.IdentityHash_()
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clPropertyEs6", "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner$J2clPropertyEs6", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6](v_)
	})
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6) SetOuter_(o *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) {
	this_.outer_ = o
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6) GetOuter_() *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner {
	return this_.outer_
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6) Implements_jscomp_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty() {
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6) initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6() {
	_ = this_
}

// Java constructor: J2clPropertyEs6(com.google.javascript.rhino.Node getKey, com.google.javascript.rhino.Node setKey)
func NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6(outer_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner, getKey pkg_rhino.Node, setKey pkg_rhino.Node) *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6 {
	this_ := &J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6{}
	this_.outer_ = outer_
	InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6(this_, getKey, setKey)
	return this_
}

func InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6(this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6, getKey pkg_rhino.Node, setKey pkg_rhino.Node) {
	_ = getKey
	_ = setKey
	InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty(&this_.J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase, getKey, setKey)
	this_.initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6()
	if !(getKey.IsGetterDef()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(getKey))
	}
	if !((setKey == 0) || setKey.IsSetterDef()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(setKey))
	}
}

// Java method: J2clPropertyEs6.remove()
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6) Remove() {
	var classMembers pkg_rhino.Node = this_.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().getKey.GetParent()
	_ = classMembers
	if !(classMembers.IsClassMembers()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(classMembers))
	}
	this_.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().getKey.Detach()
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(this_.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().getKey, this_.outer_.outer_.compiler)
	if this_.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().setKey != 0 {
		this_.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().setKey.Detach()
		NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(this_.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().setKey, this_.outer_.outer_.compiler)
	}
	this_.outer_.outer_.compiler.ReportChangeToEnclosingScope(classMembers)
}

// Java class: com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.GatherJ2CLClassGetterSetters
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_               *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner
	j2clPropertiesByName j2gort.Ext_java_util_Map `gson:"j2clPropertiesByName"`
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) IsJavaNil_() bool {
	return this_ == nil
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner$GatherJ2CLClassGetterSetters", "GatherJ2CLClassGetterSetters"
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.GatherJ2CLClassGetterSetters", "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner$GatherJ2CLClassGetterSetters", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters](v_)
	})
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) SetOuter_(o *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) {
	this_.outer_ = o
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) GetOuter_() *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner {
	return this_.outer_
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters() {
	_ = this_
	this_.j2clPropertiesByName = j2gort.NewExt_java_util_LinkedHashMap()
}

// Java constructor: GatherJ2CLClassGetterSetters()
func NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters(outer_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters {
	this_ := &J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters{}
	this_.outer_ = outer_
	InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters(this_)
	return this_
}

func InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters(this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters()
}

// Java method: GatherJ2CLClassGetterSetters.getResults()
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) priv_GatherJ2CLClassGetterSetters_getResults_3490() j2gort.Ext_java_util_Map {
	return this_.j2clPropertiesByName
}

// Java method: GatherJ2CLClassGetterSetters.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsClass() {
		this_.VisitClass(n)
	} else if NodeUtil_IsObjectDefinePropertiesDefinition(n) {
		this_.VisitObjectDefineProperties(n)
	}
}

// Java method: GatherJ2CLClassGetterSetters.visitObjectDefineProperties(com.google.javascript.rhino.Node n)
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) VisitObjectDefineProperties(n pkg_rhino.Node) {
	_ = n
	var className pkg_rhino.Node = n.GetSecondChild()
	_ = className
	if !(className.IsName()) {
		return
	}
	var classNameString string = className.GetString()
	_ = classNameString
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(NodeUtil_GetObjectDefinedPropertiesKeys(n)))
		for cur_2.HasNext() {
			p := pkg_rhino.Node(cur_2.NextNode())
			_ = p
			var name string = p.GetString()
			_ = name
			var propertyLiteral pkg_rhino.Node = p.GetFirstChild()
			_ = propertyLiteral
			var getKey pkg_rhino.Node = pkg_rhino.Node(0)
			_ = getKey
			var setKey pkg_rhino.Node = pkg_rhino.Node(0)
			_ = setKey
			var hasSetter bool = false
			_ = hasSetter
			{
				var innerKey pkg_rhino.Node = propertyLiteral.GetFirstChild()
				_ = innerKey
				for ; innerKey != 0; innerKey = innerKey.GetNext() {
					if !(innerKey.IsStringKey()) {
						continue
					}
					switch j2gort.CheckNotNullString(innerKey.GetString()) {
					case "get":
						if this_.outer_.priv_StaticFieldGetterSetterInliner_matchesJ2clGetKeySignature_3485(classNameString, innerKey) {
							getKey = innerKey
						}
					case "set":
						hasSetter = true
						if this_.outer_.priv_StaticFieldGetterSetterInliner_matchesJ2clSetKeySignature_3486(classNameString, innerKey) {
							setKey = innerKey
						}
					default:
					}
				}
			}
			if (getKey != 0) && (!(hasSetter) || (setKey != 0)) {
				j2gort.MapPutStr(this_.j2clPropertiesByName, (j2gort.StrValString(classNameString) + "." + j2gort.StrValString(name)), NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs5(this_.outer_, getKey, setKey))
			}
		}
	}
}

// Java method: GatherJ2CLClassGetterSetters.visitClass(com.google.javascript.rhino.Node classNode)
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) VisitClass(classNode pkg_rhino.Node) {
	_ = classNode
	var className string = NodeUtil_GetName(classNode)
	_ = className
	var classMembers pkg_rhino.Node = NodeUtil_GetClassMembers(classNode)
	_ = classMembers
	var setterDefByName j2gort.OrderedMap[string, pkg_rhino.Node] = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
	_ = setterDefByName
	var getterDefByName j2gort.OrderedMap[string, pkg_rhino.Node] = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
	_ = getterDefByName
	{
		var memberFunction pkg_rhino.Node = classMembers.GetFirstChild()
		_ = memberFunction
		for ; memberFunction != 0; memberFunction = memberFunction.GetNext() {
			if !(memberFunction.IsStaticMember()) {
				continue
			}
			switch memberFunction.TokenSwitchOrd_() {
			case int32(95) /* Token.GETTER_DEF */ :
				getterDefByName.Put(memberFunction.GetString(), memberFunction)
			case int32(96) /* Token.SETTER_DEF */ :
				setterDefByName.Put(memberFunction.GetString(), memberFunction)
			default:
			}
		}
	}
	{
		cur_3 := getterDefByName.IterKeys()
		for cur_3.HasNext() {
			propertyName := cur_3.
				Next()
			_ = propertyName
			var getterDef pkg_rhino.Node = getterDefByName.Get(propertyName)
			_ = getterDef
			var setterDef pkg_rhino.Node = setterDefByName.Get(propertyName)
			_ = setterDef
			if this_.outer_.priv_StaticFieldGetterSetterInliner_matchesJ2clGetKeySignature_3485(className, getterDef) && ((setterDef == 0) || this_.outer_.priv_StaticFieldGetterSetterInliner_matchesJ2clSetKeySignature_3486(className, setterDef)) {
				j2gort.MapPutStr(this_.j2clPropertiesByName, (j2gort.StrValString(className) + "." + j2gort.StrValString(propertyName)), NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyEs6(this_.outer_, getterDef, setterDef))
			}
		}
	}
}

// Java class: com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.DetermineInlinableProperties
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	propertiesByName j2gort.Ext_java_util_Map `gson:"propertiesByName"`
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties) IsJavaNil_() bool {
	return this_ == nil
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner$DetermineInlinableProperties", "DetermineInlinableProperties"
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.DetermineInlinableProperties", "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner$DetermineInlinableProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties](v_)
	})
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties) initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties() {
	_ = this_
}

// Java constructor: DetermineInlinableProperties(java.util.Map<java.lang.String,com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clProperty> allGetterSetters)
func NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties(allGetterSetters j2gort.Ext_java_util_Map) *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties {
	this_ := &J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties{}
	InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties(this_, allGetterSetters)
	return this_
}

func InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties(this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties, allGetterSetters j2gort.Ext_java_util_Map) {
	_ = allGetterSetters
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties()
	this_.propertiesByName = allGetterSetters
}

// Java method: DetermineInlinableProperties.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if (NodeUtil_IsCompoundAssignmentOp(n) || n.IsInc()) || n.IsDec() {
		var assignmentTarget pkg_rhino.Node = n.GetFirstChild()
		_ = assignmentTarget
		if assignmentTarget.IsGetProp() {
			var accessName string = assignmentTarget.GetQualifiedName()
			_ = accessName
			var prop J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty = j2gort.AsIface[J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty](j2gort.MapGetStr(this_.propertiesByName, accessName))
			_ = prop
			if prop != nil && !prop.IsJavaNil_() {
				{
					r_4 := prop
					r_4.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().isSafeToInline = false
				}
			}
		}
	}
}

// Java class: com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.InlinePropertiesPass
type J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_           *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner
	propertiesByName j2gort.Ext_java_util_Map `gson:"propertiesByName"`
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) IsJavaNil_() bool {
	return this_ == nil
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner$InlinePropertiesPass", "InlinePropertiesPass"
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.InlinePropertiesPass", "com.google.javascript.jscomp.J2clPropertyInlinerPass$StaticFieldGetterSetterInliner$InlinePropertiesPass", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass](v_)
	})
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) SetOuter_(o *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner) {
	this_.outer_ = o
}
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) GetOuter_() *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner {
	return this_.outer_
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass() {
	_ = this_
}

// Java constructor: InlinePropertiesPass(java.util.Map<java.lang.String,com.google.javascript.jscomp.J2clPropertyInlinerPass.StaticFieldGetterSetterInliner.J2clProperty> allGetterSetters)
func NewJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass(outer_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner, allGetterSetters j2gort.Ext_java_util_Map) *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass {
	this_ := &J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass{}
	this_.outer_ = outer_
	InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass(this_, allGetterSetters)
	return this_
}

func InitJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass(this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass, allGetterSetters j2gort.Ext_java_util_Map) {
	_ = allGetterSetters
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass()
	this_.propertiesByName = allGetterSetters
}

// Java method: InlinePropertiesPass.run()
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) priv_InlinePropertiesPass_run_3493() {
	NodeTraversal_Traverse(this_.outer_.outer_.compiler, this_.outer_.root, this_)
	{
		cur_5 := j2gort.NewMapValCursor(this_.propertiesByName)
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			prop := j2gort.AsIface[J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty](rawElem_6)
			_ = prop
			if prop.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().isSafeToInline {
				prop.Remove()
			}
		}
	}
}

// Java method: InlinePropertiesPass.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsGetProp() {
		if parent.IsExprResult() {
			return
		}
		if NodeUtil_IsAssignmentOp(parent) && (parent.GetFirstChild() == n) {
			return
		}
		var accessName string = n.GetQualifiedName()
		_ = accessName
		var prop J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty = j2gort.AsIface[J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty](j2gort.MapGetStr(this_.propertiesByName, accessName))
		_ = prop
		if (prop != nil && !prop.IsJavaNil_()) && prop.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().isSafeToInline {
			var injector *FunctionInjector = NewFunctionInjector_Builder(this_.outer_.outer_.compiler).AssumeStrictThis(true).Build()
			_ = injector
			var functionCall pkg_rhino.Node = pkg_rhino.IR_Call(pkg_rhino.IR_Name("inlined_j2cl_getter"), []pkg_rhino.Node{})
			_ = functionCall
			n.ReplaceWith(functionCall)
			var reference FunctionInjector_Reference = NewFunctionInjector_Reference(functionCall, t.GetScope(), t.GetChunk(), StaticGet_FunctionInjector_InliningMode_DIRECT())
			_ = reference
			var inlinedCall pkg_rhino.Node = injector.Inline(reference, j2gort.NullString, prop.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().getKey.GetFirstChild())
			_ = inlinedCall
			t.GetCompiler().ReportChangeToEnclosingScope(inlinedCall)
		}
	}
	if n.IsAssign() {
		var assignmentTarget pkg_rhino.Node = n.GetFirstChild()
		_ = assignmentTarget
		var assignmentValue pkg_rhino.Node = n.GetLastChild()
		_ = assignmentValue
		if assignmentTarget.IsGetProp() {
			var accessName_2 string = assignmentTarget.GetQualifiedName()
			_ = accessName_2
			var prop_2 J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty = j2gort.AsIface[J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clProperty](j2gort.MapGetStr(this_.propertiesByName, accessName_2))
			_ = prop_2
			if ((prop_2 != nil && !prop_2.IsJavaNil_()) && (prop_2.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().setKey != 0)) && prop_2.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().isSafeToInline {
				var injector_2 *FunctionInjector = NewFunctionInjector_Builder(this_.outer_.outer_.compiler).AssumeStrictThis(true).Build()
				_ = injector_2
				assignmentValue.Detach()
				var functionCall_2 pkg_rhino.Node = pkg_rhino.IR_Call(pkg_rhino.IR_Name("inlined_j2cl_setter"), []pkg_rhino.Node{assignmentValue})
				_ = functionCall_2
				n.ReplaceWith(functionCall_2)
				var reference_2 FunctionInjector_Reference = NewFunctionInjector_Reference(functionCall_2, t.GetScope(), t.GetChunk(), StaticGet_FunctionInjector_InliningMode_BLOCK())
				_ = reference_2
				injector_2.MaybePrepareCall(reference_2)
				var inlinedCall_2 pkg_rhino.Node = injector_2.Inline(reference_2, j2gort.NullString, prop_2.GetJ2clPropertyInlinerPass_StaticFieldGetterSetterInliner_J2clPropertyBase().setKey.GetFirstChild())
				_ = inlinedCall_2
				t.GetCompiler().ReportChangeToEnclosingScope(inlinedCall_2)
			}
		}
	}
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) AlwaysShouldTraverse_() {
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_GatherJ2CLClassGetterSetters) CanReleaseSyntacticScopes_() {
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties) AlwaysShouldTraverse_() {
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_DetermineInlinableProperties) CanReleaseSyntacticScopes_() {
}
func (*J2clPropertyInlinerPass_StaticFieldGetterSetterInliner_InlinePropertiesPass) AlwaysShouldTraverse_() {
}
