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
// Port of com/google/javascript/jscomp/AbstractPeepholeOptimization.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java polymorphic class interface: com.google.javascript.jscomp.AbstractPeepholeOptimization
type AbstractPeepholeOptimization interface {
	Implements_jscomp_AbstractPeepholeOptimization()
	IsJavaNil_() bool
	GetAbstractPeepholeOptimizationBase() *AbstractPeepholeOptimizationBase
	OptimizeSubtree(pkg_rhino.Node) pkg_rhino.Node
	Report(*DiagnosticType, pkg_rhino.Node)
	AreNodesEqualForInlining(pkg_rhino.Node, pkg_rhino.Node) bool
	IsASTNormalized() bool
	BeginTraversal(AbstractCompiler)
	EndTraversal()
	MayEffectMutableState(pkg_rhino.Node) bool
	MayHaveSideEffects(pkg_rhino.Node) bool
	GetSideEffectFreeNumberValue(pkg_rhino.Node) *j2gort.BoxDouble
	GetSideEffectFreeNumberValueNoConversion(pkg_rhino.Node) *j2gort.BoxDouble
	GetSideEffectFreeBigIntValue(pkg_rhino.Node) *j2gort.Ext_java_math_BigInteger
	GetSideEffectFreeStringValue(pkg_rhino.Node) string
	GetSideEffectFreeBooleanValue(pkg_rhino.Node) pkg_base.Tri
	NodeTypeMayHaveSideEffects(pkg_rhino.Node) bool
	IsEcmaScript5OrGreater() bool
	GetCodingConvention() CodingConvention
	ReportChangeToEnclosingScope(pkg_rhino.Node)
	DeleteNode(pkg_rhino.Node)
	MarkFunctionsDeleted(pkg_rhino.Node)
	MarkNewScopesChanged(pkg_rhino.Node)
	AddFeatureToEnclosingScript(*pkg_parser.FeatureSet_Feature)
	GetNewFeatures() j2gort.Ext_com_google_common_collect_ImmutableSet
	ClearNewFeatures()
}

func ToIface_AbstractPeepholeOptimization[T interface {
	AbstractPeepholeOptimization
	comparable
}](v T) AbstractPeepholeOptimization {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.AbstractPeepholeOptimization
type AbstractPeepholeOptimizationBase struct {
	self_       AbstractPeepholeOptimization
	compiler    AbstractCompiler                                  `gson:"compiler"`
	astAnalyzer *AstAnalyzer                                      `gson:"astAnalyzer"`
	newFeatures j2gort.OrderedSet[*pkg_parser.FeatureSet_Feature] `gson:"newFeatures"`
	idHash_     atomic.Int32
}

func (this_ *AbstractPeepholeOptimizationBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(AbstractPeepholeOptimization)
}
func (this_ *AbstractPeepholeOptimizationBase) IsJavaNil_() bool { return this_ == nil }
func (*AbstractPeepholeOptimizationBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractPeepholeOptimization", "AbstractPeepholeOptimization"
}
func (this_ *AbstractPeepholeOptimizationBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractPeepholeOptimization", "com.google.javascript.jscomp.AbstractPeepholeOptimization", func(v_ any) bool {
		return j2gort.InstanceOf[AbstractPeepholeOptimization](v_)
	})
}
func (b *AbstractPeepholeOptimizationBase) GetAbstractPeepholeOptimizationBase() *AbstractPeepholeOptimizationBase {
	return b
}
func (*AbstractPeepholeOptimizationBase) Implements_jscomp_AbstractPeepholeOptimization() {}

func (this_ *AbstractPeepholeOptimizationBase) initFields_AbstractPeepholeOptimization() {
	_ = this_
	this_.newFeatures = j2gort.NewOrderedSet[*pkg_parser.FeatureSet_Feature]()
}

func InitAbstractPeepholeOptimization(this_ *AbstractPeepholeOptimizationBase) {
	this_.initFields_AbstractPeepholeOptimization()
}

// Java method: AbstractPeepholeOptimization.optimizeSubtree(com.google.javascript.rhino.Node subtree)
func (this_ *AbstractPeepholeOptimizationBase) OptimizeSubtree(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	panic("abstract or native method invoked")
}

// Java method: AbstractPeepholeOptimization.report(com.google.javascript.jscomp.DiagnosticType diagnostic, com.google.javascript.rhino.Node n)
func (this_ *AbstractPeepholeOptimizationBase) Report(diagnostic *DiagnosticType, n pkg_rhino.Node) {
	_ = diagnostic
	_ = n
	var error_ *JSError = JSError_Make_3(n, diagnostic, []string{n.String()})
	_ = error_
	this_.compiler.Report_1(error_)
}

// Java method: AbstractPeepholeOptimization.areNodesEqualForInlining(com.google.javascript.rhino.Node n1, com.google.javascript.rhino.Node n2)
func (this_ *AbstractPeepholeOptimizationBase) AreNodesEqualForInlining(n1 pkg_rhino.Node, n2 pkg_rhino.Node) bool {
	_ = n1
	_ = n2
	j2gort.CheckNotNullIface(this_.compiler)
	return this_.compiler.AreNodesEqualForInlining(n1, n2)
}

// Java method: AbstractPeepholeOptimization.isASTNormalized()
func (this_ *AbstractPeepholeOptimizationBase) IsASTNormalized() bool {
	j2gort.CheckNotNullIface(this_.compiler)
	return this_.compiler.GetLifeCycleStage().IsNormalized()
}

// Java method: AbstractPeepholeOptimization.beginTraversal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *AbstractPeepholeOptimizationBase) BeginTraversal(compiler AbstractCompiler) {
	_ = compiler
	this_.compiler = j2gort.CheckNotNullIface(compiler)
	this_.astAnalyzer = compiler.GetAstAnalyzer()
}

// Java method: AbstractPeepholeOptimization.endTraversal()
func (this_ *AbstractPeepholeOptimizationBase) EndTraversal() {
	j2gort.CheckStateFmt(this_.newFeatures.IsEmpty(), "Expected getNewFeatures() to be empty in endTraversal() but found %s for %s", &this_.newFeatures, j2gort.BoxString(j2gort.CheckNotNullPtr(j2gort.GetObjectClass(this_.self_)).GetName()))
	this_.compiler = nil
	this_.astAnalyzer = nil
}

// Java method: AbstractPeepholeOptimization.mayEffectMutableState(com.google.javascript.rhino.Node n)
func (this_ *AbstractPeepholeOptimizationBase) MayEffectMutableState(n pkg_rhino.Node) bool {
	_ = n
	return this_.astAnalyzer.MayEffectMutableState(n)
}

// Java method: AbstractPeepholeOptimization.mayHaveSideEffects(com.google.javascript.rhino.Node n)
func (this_ *AbstractPeepholeOptimizationBase) MayHaveSideEffects(n pkg_rhino.Node) bool {
	_ = n
	return this_.astAnalyzer.MayHaveSideEffects(n)
}

// Java method: AbstractPeepholeOptimization.getSideEffectFreeNumberValue(com.google.javascript.rhino.Node n)
func (this_ *AbstractPeepholeOptimizationBase) GetSideEffectFreeNumberValue(n pkg_rhino.Node) *j2gort.BoxDouble {
	_ = n
	var value *j2gort.BoxDouble = NodeUtil_GetNumberValue(n)
	_ = value
	if (value != nil) && this_.astAnalyzer.MayHaveSideEffects(n) {
		value = nil
	}
	return value
}

// Java method: AbstractPeepholeOptimization.getSideEffectFreeNumberValueNoConversion(com.google.javascript.rhino.Node n)
func (this_ *AbstractPeepholeOptimizationBase) GetSideEffectFreeNumberValueNoConversion(n pkg_rhino.Node) *j2gort.BoxDouble {
	_ = n
	var value *j2gort.BoxDouble = NodeUtil_GetNumberValueNoConversions(n)
	_ = value
	if (value != nil) && this_.astAnalyzer.MayHaveSideEffects(n) {
		value = nil
	}
	return value
}

// Java method: AbstractPeepholeOptimization.getSideEffectFreeBigIntValue(com.google.javascript.rhino.Node n)
func (this_ *AbstractPeepholeOptimizationBase) GetSideEffectFreeBigIntValue(n pkg_rhino.Node) *j2gort.Ext_java_math_BigInteger {
	_ = n
	var value *j2gort.Ext_java_math_BigInteger = NodeUtil_GetBigIntValue(n)
	_ = value
	if (value != nil) && this_.astAnalyzer.MayHaveSideEffects(n) {
		value = nil
	}
	return value
}

// Java method: AbstractPeepholeOptimization.getSideEffectFreeStringValue(com.google.javascript.rhino.Node n)
func (this_ *AbstractPeepholeOptimizationBase) GetSideEffectFreeStringValue(n pkg_rhino.Node) string {
	_ = n
	var value string = NodeUtil_GetStringValue(n)
	_ = value
	if (value != j2gort.NullString) && this_.astAnalyzer.MayHaveSideEffects(n) {
		value = j2gort.NullString
	}
	return value
}

// Java method: AbstractPeepholeOptimization.getSideEffectFreeBooleanValue(com.google.javascript.rhino.Node n)
func (this_ *AbstractPeepholeOptimizationBase) GetSideEffectFreeBooleanValue(n pkg_rhino.Node) pkg_base.Tri {
	_ = n
	var value pkg_base.Tri = NodeUtil_GetBooleanValue(n)
	_ = value
	if (value != pkg_base.StaticGet_Tri_UNKNOWN()) && this_.astAnalyzer.MayHaveSideEffects(n) {
		value = pkg_base.StaticGet_Tri_UNKNOWN()
	}
	return value
}

// Java method: AbstractPeepholeOptimization.nodeTypeMayHaveSideEffects(com.google.javascript.rhino.Node n)
func (this_ *AbstractPeepholeOptimizationBase) NodeTypeMayHaveSideEffects(n pkg_rhino.Node) bool {
	_ = n
	return this_.astAnalyzer.NodeTypeMayHaveSideEffects(n)
}

// Java method: AbstractPeepholeOptimization.isEcmaScript5OrGreater()
func (this_ *AbstractPeepholeOptimizationBase) IsEcmaScript5OrGreater() bool {
	return ((this_.compiler != nil && !this_.compiler.IsJavaNil_()) && this_.compiler.GetOptions().GetOutputFeatureSet().Contains_1_FeatureSet(pkg_parser.StaticGet_FeatureSet_ES5()))
}

// Java method: AbstractPeepholeOptimization.getCodingConvention()
func (this_ *AbstractPeepholeOptimizationBase) GetCodingConvention() CodingConvention {
	return this_.compiler.GetCodingConvention()
}

// Java method: AbstractPeepholeOptimization.reportChangeToEnclosingScope(com.google.javascript.rhino.Node n)
func (this_ *AbstractPeepholeOptimizationBase) ReportChangeToEnclosingScope(n pkg_rhino.Node) {
	_ = n
	this_.compiler.ReportChangeToEnclosingScope(n)
}

// Java method: AbstractPeepholeOptimization.deleteNode(com.google.javascript.rhino.Node property)
func (this_ *AbstractPeepholeOptimizationBase) DeleteNode(property pkg_rhino.Node) {
	_ = property
	j2gort.CheckNotNullIface(this_.compiler)
	NodeUtil_DeleteNode(property, this_.compiler)
}

// Java method: AbstractPeepholeOptimization.markFunctionsDeleted(com.google.javascript.rhino.Node function)
func (this_ *AbstractPeepholeOptimizationBase) MarkFunctionsDeleted(function pkg_rhino.Node) {
	_ = function
	j2gort.CheckNotNullIface(this_.compiler)
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(function, this_.compiler)
}

// Java method: AbstractPeepholeOptimization.markNewScopesChanged(com.google.javascript.rhino.Node n)
func (this_ *AbstractPeepholeOptimizationBase) MarkNewScopesChanged(n pkg_rhino.Node) {
	_ = n
	j2gort.CheckNotNullIface(this_.compiler)
	NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(n, this_.compiler)
}

// Java method: AbstractPeepholeOptimization.addFeatureToEnclosingScript(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func (this_ *AbstractPeepholeOptimizationBase) AddFeatureToEnclosingScript(feature *pkg_parser.FeatureSet_Feature) {
	_ = feature
	this_.newFeatures.Add(feature)
}

// Java method: AbstractPeepholeOptimization.getNewFeatures()
func (this_ *AbstractPeepholeOptimizationBase) GetNewFeatures() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.newFeatures.ToImmutableSet()
}

// Java method: AbstractPeepholeOptimization.clearNewFeatures()
func (this_ *AbstractPeepholeOptimizationBase) ClearNewFeatures() {
	this_.newFeatures.Clear()
}
