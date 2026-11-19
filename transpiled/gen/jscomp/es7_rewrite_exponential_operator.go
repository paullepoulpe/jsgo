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
// Port of com/google/javascript/jscomp/Es7RewriteExponentialOperator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Es7RewriteExponentialOperator
type Es7RewriteExponentialOperator struct {
	compiler         AbstractCompiler                       `gson:"compiler"`
	astFactory       *AstFactory                            `gson:"astFactory"`
	uniqueIdSupplier *UniqueIdSupplier                      `gson:"uniqueIdSupplier"`
	mathPowCall      j2gort.Ext_java_util_function_Supplier `gson:"mathPowCall"`
	idHash_          atomic.Int32
}

func (this_ *Es7RewriteExponentialOperator) BindSelf_(s any) {
	_ = s
}
func (this_ *Es7RewriteExponentialOperator) IsJavaNil_() bool { return this_ == nil }
func (*Es7RewriteExponentialOperator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es7RewriteExponentialOperator", "Es7RewriteExponentialOperator"
}
func (this_ *Es7RewriteExponentialOperator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es7RewriteExponentialOperator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es7RewriteExponentialOperator", "com.google.javascript.jscomp.Es7RewriteExponentialOperator", func(v_ any) bool {
		return j2gort.InstanceOf[*Es7RewriteExponentialOperator](v_)
	})
}
func (*Es7RewriteExponentialOperator) Implements_jscomp_NodeTraversal_Callback() {}
func (*Es7RewriteExponentialOperator) Implements_jscomp_CompilerPass()           {}

var static_Es7RewriteExponentialOperator_TRANSPILE_EXPONENT_USING_BIGINT *DiagnosticType

func StaticGet_Es7RewriteExponentialOperator_TRANSPILE_EXPONENT_USING_BIGINT() *DiagnosticType {
	clinit_Es7RewriteExponentialOperator()
	return static_Es7RewriteExponentialOperator_TRANSPILE_EXPONENT_USING_BIGINT
}
func StaticSet_Es7RewriteExponentialOperator_TRANSPILE_EXPONENT_USING_BIGINT(v *DiagnosticType) {
	clinit_Es7RewriteExponentialOperator()
	static_Es7RewriteExponentialOperator_TRANSPILE_EXPONENT_USING_BIGINT = v
}
func StaticGet_Es7RewriteExponentialOperator_TEMP_VAR_NAME_PREFIX() string {
	return "$jscomp$exp$assign$tmp"
}
func StaticGet_Es7RewriteExponentialOperator_TEMP_INDEX_VAR_NAME_PREFIX() string {
	return "$jscomp$exp$assign$tmpindex"
}

var clinitDone_Es7RewriteExponentialOperator uint32
var clinitState_Es7RewriteExponentialOperator j2gort.ClinitState

func clinit_Es7RewriteExponentialOperator() {
	if atomic.LoadUint32(&clinitDone_Es7RewriteExponentialOperator) != 0 {
		return
	}
	clinitSlow_Es7RewriteExponentialOperator()
}

func clinitSlow_Es7RewriteExponentialOperator() {
	if !j2gort.EnterClinit(&clinitState_Es7RewriteExponentialOperator, &clinitDone_Es7RewriteExponentialOperator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es7RewriteExponentialOperator, &clinitDone_Es7RewriteExponentialOperator)
	static_Es7RewriteExponentialOperator_TRANSPILE_EXPONENT_USING_BIGINT = DiagnosticType_Error("JSC_TRANSPILE_EXPONENT_USING_BIGINT", "Cannot transpile `**` operator applied to BigInt operands.")
}

func (this_ *Es7RewriteExponentialOperator) initFields_Es7RewriteExponentialOperator() {
	_ = this_
	this_.mathPowCall = j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](j2gort.Ext_com_google_common_base_Suppliers_Memoize(func() j2gort.Ext_com_google_common_base_Supplier {
		bound_1 := this_
		return j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any { return j2gort.BoxNode(bound_1.priv_Es7RewriteExponentialOperator_createMathPowCall_2784()) })
	}()))
}

// Java constructor: Es7RewriteExponentialOperator(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs7RewriteExponentialOperator(compiler AbstractCompiler) *Es7RewriteExponentialOperator {
	clinit_Es7RewriteExponentialOperator()
	this_ := &Es7RewriteExponentialOperator{}
	InitEs7RewriteExponentialOperator(this_, compiler)
	return this_
}

func InitEs7RewriteExponentialOperator(this_ *Es7RewriteExponentialOperator, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_Es7RewriteExponentialOperator()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
}

// Java method: Es7RewriteExponentialOperator.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es7RewriteExponentialOperator) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(this_.compiler, root, pkg_parser.StaticGet_FeatureSet_Feature_EXPONENT_OP())
}

// Java method: Es7RewriteExponentialOperator.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es7RewriteExponentialOperator) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		var scriptFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(n)
		_ = scriptFeatures
		return ((scriptFeatures == nil) || scriptFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_EXPONENT_OP()))
	}
	return true
}

// Java method: Es7RewriteExponentialOperator.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es7RewriteExponentialOperator) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(18) /* Token.EXPONENT */ :
		if this_.priv_Es7RewriteExponentialOperator_checkOperatorType_2785(n) {
			this_.priv_Es7RewriteExponentialOperator_visitExponentiationOperator_2780(n)
		}
	case int32(63) /* Token.ASSIGN_EXPONENT */ :
		if this_.priv_Es7RewriteExponentialOperator_checkOperatorType_2785(n) {
			this_.priv_Es7RewriteExponentialOperator_visitExponentiationAssignmentOperator_2781(t, n)
		}
	default:
	}
}

// Java method: Es7RewriteExponentialOperator.visitExponentiationOperator(com.google.javascript.rhino.Node operator)
func (this_ *Es7RewriteExponentialOperator) priv_Es7RewriteExponentialOperator_visitExponentiationOperator_2780(operator pkg_rhino.Node) {
	_ = operator
	var callClone pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](this_.mathPowCall.Get()).CloneTree_0()
	_ = callClone
	callClone.AddChildToBack(operator.RemoveFirstChild())
	callClone.AddChildToBack(operator.RemoveFirstChild())
	callClone.SrcrefTreeIfMissing(operator)
	operator.ReplaceWith(callClone)
	this_.compiler.ReportChangeToEnclosingScope(callClone)
}

// Java method: Es7RewriteExponentialOperator.visitExponentiationAssignmentOperator(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node operator)
func (this_ *Es7RewriteExponentialOperator) priv_Es7RewriteExponentialOperator_visitExponentiationAssignmentOperator_2781(t *NodeTraversal, operator pkg_rhino.Node) {
	_ = t
	_ = operator
	var enclosingStatement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(operator)
	_ = enclosingStatement
	var left pkg_rhino.Node = operator.RemoveFirstChild()
	_ = left
	var replacement pkg_rhino.Node = pkg_rhino.Node(0)
	_ = replacement
	if left.IsName() {
		replacement = this_.priv_Es7RewriteExponentialOperator_handleLHSName_2782(operator, left)
	} else {
		if !(left.IsGetProp() || left.IsGetElem()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(left))
		}
		replacement = this_.priv_Es7RewriteExponentialOperator_handleLHSPropertyReference_2783(t, operator, left, enclosingStatement)
	}
	operator.ReplaceWith(replacement)
	this_.compiler.ReportChangeToEnclosingScope(enclosingStatement)
}

// Java method: Es7RewriteExponentialOperator.handleLHSName(com.google.javascript.rhino.Node operator, com.google.javascript.rhino.Node left)
func (this_ *Es7RewriteExponentialOperator) priv_Es7RewriteExponentialOperator_handleLHSName_2782(operator pkg_rhino.Node, left pkg_rhino.Node) pkg_rhino.Node {
	_ = operator
	_ = left
	var callClone pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](this_.mathPowCall.Get()).CloneTree_0()
	_ = callClone
	callClone.AddChildToBack(left.CloneTree_0())
	callClone.AddChildToBack(operator.RemoveFirstChild())
	return this_.astFactory.CreateAssign_2_Node_Node(left, callClone).SrcrefTreeIfMissing(operator)
}

// Java method: Es7RewriteExponentialOperator.handleLHSPropertyReference(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node operator, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node enclosingStatement)
func (this_ *Es7RewriteExponentialOperator) priv_Es7RewriteExponentialOperator_handleLHSPropertyReference_2783(t *NodeTraversal, operator pkg_rhino.Node, left pkg_rhino.Node, enclosingStatement pkg_rhino.Node) pkg_rhino.Node {
	_ = t
	_ = operator
	_ = left
	_ = enclosingStatement
	var input *CompilerInput = t.GetInput()
	_ = input
	var uniqueId string = this_.uniqueIdSupplier.GetUniqueId(input)
	_ = uniqueId
	var tempVarName string = ("$jscomp$exp$assign$tmp" + j2gort.StrValString(uniqueId))
	_ = tempVarName
	var newLHS pkg_rhino.Node = pkg_rhino.Node(0)
	_ = newLHS
	var tempPropOrElem pkg_rhino.Node = pkg_rhino.Node(0)
	_ = tempPropOrElem
	var objectNode pkg_rhino.Node = left.RemoveFirstChild()
	_ = objectNode
	var let pkg_rhino.Node = this_.astFactory.CreateSingleLetNameDeclaration_1(tempVarName).SrcrefTree(operator)
	_ = let
	let.InsertBefore(enclosingStatement)
	var tempName pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(tempVarName, AstFactory_S_Type_1_Node(objectNode)).Srcref(objectNode)
	_ = tempName
	var assignTemp pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(tempName, objectNode).Srcref(objectNode)
	_ = assignTemp
	if left.IsGetProp() {
		var propertyName string = left.GetString()
		_ = propertyName
		tempPropOrElem = this_.astFactory.CreateGetProp(tempName.CloneNode_0(), propertyName, AstFactory_S_Type_1_Node(left)).Srcref(left)
		newLHS = this_.astFactory.CreateGetProp(assignTemp, propertyName, AstFactory_S_Type_1_Node(left)).Srcref(left)
	} else {
		if !(left.IsGetElem()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(left))
		}
		var tempIndexVarName string = ("$jscomp$exp$assign$tmpindex" + j2gort.StrValString(uniqueId))
		_ = tempIndexVarName
		var indexExprNode pkg_rhino.Node = left.GetLastChild().Detach()
		_ = indexExprNode
		var letIndex pkg_rhino.Node = this_.astFactory.CreateSingleLetNameDeclaration_1(tempIndexVarName).SrcrefTree(operator)
		_ = letIndex
		letIndex.InsertBefore(enclosingStatement)
		var tempIndexName pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(tempIndexVarName, AstFactory_S_Type_1_Node(indexExprNode)).Srcref(indexExprNode)
		_ = tempIndexName
		var assignTempIndex pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(tempIndexName, indexExprNode).Srcref(indexExprNode)
		_ = assignTempIndex
		tempPropOrElem = this_.astFactory.CreateGetElem(tempName.CloneNode_0(), tempIndexName.CloneNode_0()).CopyTypeFrom(left).Srcref(left)
		newLHS = this_.astFactory.CreateGetElem(assignTemp, assignTempIndex).CopyTypeFrom(left).Srcref(left)
	}
	var callClone pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](this_.mathPowCall.Get()).CloneTree_0()
	_ = callClone
	callClone.AddChildToBack(tempPropOrElem.CloneTree_0())
	callClone.AddChildToBack(operator.RemoveFirstChild())
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
	return this_.astFactory.CreateAssign_2_Node_Node(newLHS, callClone).SrcrefTreeIfMissing(operator)
}

// Java method: Es7RewriteExponentialOperator.createMathPowCall()
func (this_ *Es7RewriteExponentialOperator) priv_Es7RewriteExponentialOperator_createMathPowCall_2784() pkg_rhino.Node {
	return this_.astFactory.CreateCall(this_.astFactory.CreateQName_2_StaticScope_string(this_.compiler.GetTranspilationNamespace(), "Math.pow"), AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NUMBER()), []pkg_rhino.Node{})
}

// Java method: Es7RewriteExponentialOperator.checkOperatorType(com.google.javascript.rhino.Node operator)
func (this_ *Es7RewriteExponentialOperator) priv_Es7RewriteExponentialOperator_checkOperatorType_2785(operator pkg_rhino.Node) bool {
	_ = operator
	if !(operator.IsExponent() || operator.IsAssignExponent()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(operator))
	}
	if this_.compiler.HasTypeCheckingRun() {
		if operator.GetColor().Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_BIGINT())) {
			this_.compiler.Report_1(JSError_Make_3(operator, StaticGet_Es7RewriteExponentialOperator_TRANSPILE_EXPONENT_USING_BIGINT(), []string{}))
			return false
		}
	}
	return true
}
func (*Es7RewriteExponentialOperator) CanReleaseSyntacticScopes_() {}
