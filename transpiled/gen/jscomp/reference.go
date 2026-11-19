/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/Reference.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Reference
type Reference struct {
	basicBlock *BasicBlock `gson:"basicBlock"`
	scope *Scope `gson:"scope"`
	inputId *pkg_rhino.InputId `gson:"inputId"`
	nameNode pkg_rhino.Node `gson:"nameNode"`
	idHash_ atomic.Int32
}

func (this_ *Reference) BindSelf_(s any) {
	_ = s
}
func (this_ *Reference) IsJavaNil_() bool { return this_ == nil }
func (*Reference) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.Reference", "Reference" }
func (this_ *Reference) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Reference) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Reference", "com.google.javascript.jscomp.Reference", func(v_ any) bool {
		return j2gort.InstanceOf[*Reference](v_)
	})
}
func (*Reference) Implements_rhino_StaticRef() {}

var static_Reference_DECLARATION_PARENTS j2gort.Ext_com_google_common_collect_ImmutableSet
func StaticGet_Reference_DECLARATION_PARENTS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_Reference()
	return static_Reference_DECLARATION_PARENTS
}
func StaticSet_Reference_DECLARATION_PARENTS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_Reference()
	static_Reference_DECLARATION_PARENTS = v
}
var clinitDone_Reference uint32
var clinitState_Reference j2gort.ClinitState
func clinit_Reference() {
	if atomic.LoadUint32(&clinitDone_Reference) != 0 {
		return
	}
	clinitSlow_Reference()
}

func clinitSlow_Reference() {
	if !j2gort.EnterClinit(&clinitState_Reference, &clinitDone_Reference) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Reference, &clinitDone_Reference)
	static_Reference_DECLARATION_PARENTS = j2gort.Ext_com_google_common_collect_Sets_ImmutableEnumSet_2(j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_VAR()), []j2gort.Ext_java_lang_Enum{j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_LET()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_CONST()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_PARAM_LIST()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_FUNCTION()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_CLASS()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_CATCH())})
}

func (this_ *Reference) initFields_Reference() {
	_ = this_
}

// Java constructor: Reference(com.google.javascript.rhino.Node nameNode, com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.BasicBlock basicBlock)
func NewReference_3(nameNode pkg_rhino.Node, t *NodeTraversal, basicBlock *BasicBlock) *Reference {
	this_ := &Reference{}
	InitReference_3(this_, nameNode, t, basicBlock)
	return this_
}

func InitReference_3(this_ *Reference, nameNode pkg_rhino.Node, t *NodeTraversal, basicBlock *BasicBlock) {
	_ = nameNode
	_ = t
	_ = basicBlock
	InitReference_4(this_, nameNode, basicBlock, t.GetScope(), t.GetInput().GetInputId())
}

// Java constructor: Reference(com.google.javascript.rhino.Node nameNode, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable BasicBlock basicBlock, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Scope scope, com.google.javascript.rhino.InputId inputId)
func NewReference_4(nameNode pkg_rhino.Node, basicBlock *BasicBlock, scope *Scope, inputId *pkg_rhino.InputId) *Reference {
	this_ := &Reference{}
	InitReference_4(this_, nameNode, basicBlock, scope, inputId)
	return this_
}

func InitReference_4(this_ *Reference, nameNode pkg_rhino.Node, basicBlock *BasicBlock, scope *Scope, inputId *pkg_rhino.InputId) {
	_ = nameNode
	_ = basicBlock
	_ = scope
	_ = inputId
	this_.initFields_Reference()
	this_.nameNode = nameNode
	this_.basicBlock = basicBlock
	this_.scope = scope
	this_.inputId = inputId
}

// Java method: Reference.toString()
func (this_ *Reference) String() string {
	return this_.nameNode.String()
}

// Java static method: Reference.createRefForTest(com.google.javascript.jscomp.CompilerInput input)
func Reference_CreateRefForTest(input *CompilerInput) *Reference {
	_ = input
	return NewReference_4(pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_NAME()), nil, nil, input.GetInputId())
}

// Java method: Reference.cloneWithNewScope(com.google.javascript.jscomp.Scope newScope)
func (this_ *Reference) CloneWithNewScope(newScope *Scope) *Reference {
	_ = newScope
	return NewReference_4(this_.nameNode, this_.basicBlock, newScope, this_.inputId)
}

// Java method: Reference.getSymbol()
func (this_ *Reference) GetSymbol() any {
	return j2gort.PtrToAny(j2gort.CastPtr[Var](this_.scope.GetVar(this_.nameNode.GetString())))
}

// Java method: Reference.getNode()
func (this_ *Reference) GetNode() pkg_rhino.Node {
	return this_.nameNode
}

// Java method: Reference.getInputId()
func (this_ *Reference) GetInputId() *pkg_rhino.InputId {
	return this_.inputId
}

// Java method: Reference.getSourceFile()
func (this_ *Reference) GetSourceFile() pkg_rhino.StaticSourceFile {
	return this_.nameNode.GetStaticSourceFile()
}

// Java method: Reference.isDeclaration()
func (this_ *Reference) IsDeclaration() bool {
	return Reference_IsDeclarationHelper(this_.nameNode)
}

// Java static method: Reference.isDeclarationHelper(com.google.javascript.rhino.Node node)
func Reference_IsDeclarationHelper(node pkg_rhino.Node) bool {
	_ = node
	var parent pkg_rhino.Node = node.GetParent()
	_ = parent
	if ((parent.IsClass() && !(node.IsFirstChildOf(parent)))) {
		return false
	}
	if ((parent.GetParent() == 0)) {
		return false
	}
	if ((NodeUtil_IsNameDeclaration(parent.GetParent()) && (node == parent.GetSecondChild()))) {
		return false
	}
	if ((((((parent.IsDestructuringLhs() || parent.IsDestructuringPattern()) || parent.IsRest()) || ((parent.IsStringKey() && parent.GetParent().IsObjectPattern()))) || (((parent.IsComputedProp() && parent.GetParent().IsObjectPattern()) && (node == parent.GetLastChild())))) || ((parent.IsDefaultValue() && (node == parent.GetFirstChild()))))) {
		return Reference_IsDeclarationHelper(parent)
	}
	if (parent.IsImport()) {
		return true
	}
	if ((parent.IsImportSpec() && (node == parent.GetLastChild()))) {
		return true
	}
	if (parent.IsArrowFunction()) {
		return (node == parent.GetFirstChild())
	}
	return StaticGet_Reference_DECLARATION_PARENTS().Contains(j2gort.BoxNullableVal(parent.GetToken()))
}

// Java method: Reference.isVarDeclaration()
func (this_ *Reference) IsVarDeclaration() bool {
	return this_.GetParent().IsVar()
}

// Java method: Reference.isLetDeclaration()
func (this_ *Reference) IsLetDeclaration() bool {
	return this_.GetParent().IsLet()
}

// Java method: Reference.isDotPropertyAccess()
func (this_ *Reference) IsDotPropertyAccess() bool {
	return this_.GetParent().IsGetProp()
}

// Java method: Reference.isAssignedToObjectDestructuringPattern()
func (this_ *Reference) IsAssignedToObjectDestructuringPattern() bool {
	var parent pkg_rhino.Node = this_.GetParent()
	_ = parent
	return (((parent.IsDestructuringLhs() || parent.IsAssign())) && parent.GetFirstChild().IsObjectPattern())
}

// Java method: Reference.isConstDeclaration()
func (this_ *Reference) IsConstDeclaration() bool {
	return this_.GetParent().IsConst()
}

// Java method: Reference.isHoistedFunction()
func (this_ *Reference) IsHoistedFunction() bool {
	return NodeUtil_IsHoistedFunctionDeclaration(this_.GetParent())
}

// Java method: Reference.isInitializingDeclaration()
func (this_ *Reference) IsInitializingDeclaration() bool {
	return ((((this_.IsDeclaration() && !(this_.GetParent().IsVar())) && !(this_.GetParent().IsLet()))) || this_.nameNode.HasChildren())
}

// Java method: Reference.getAssignedValue()
func (this_ *Reference) GetAssignedValue() pkg_rhino.Node {
	return NodeUtil_GetRValueOfLValue(this_.nameNode)
}

// Java method: Reference.getBasicBlock()
func (this_ *Reference) GetBasicBlock() *BasicBlock {
	return this_.basicBlock
}

// Java method: Reference.getParent()
func (this_ *Reference) GetParent() pkg_rhino.Node {
	return this_.GetNode().GetParent()
}

// Java method: Reference.getGrandparent()
func (this_ *Reference) GetGrandparent() pkg_rhino.Node {
	return this_.GetNode().GetGrandparent()
}

// Java static method: Reference.isLhsOfEnhancedForExpression(com.google.javascript.rhino.Node n)
func Reference_IsLhsOfEnhancedForExpression(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (NodeUtil_IsEnhancedFor(parent) && (parent.GetFirstChild() == n))
}

// Java method: Reference.isSimpleAssignmentToName()
func (this_ *Reference) IsSimpleAssignmentToName() bool {
	var parent pkg_rhino.Node = this_.GetParent()
	_ = parent
	return (parent.IsAssign() && (parent.GetFirstChild() == this_.nameNode))
}

// Java method: Reference.isLvalue()
func (this_ *Reference) IsLvalue() bool {
	var parent pkg_rhino.Node = this_.GetParent()
	_ = parent
	var parentType pkg_rhino.Token = parent.GetToken()
	_ = parentType
	switch (parentType).Ordinal() {
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */:
		return ((this_.nameNode.HasChildren() || Reference_IsLhsOfEnhancedForExpression(this_.nameNode)))
	case int32(134) /* Token.DEFAULT_VALUE */:
		return (parent.GetFirstChild() == this_.nameNode)
	case int32(71) /* Token.INC */, int32(72) /* Token.DEC */, int32(86) /* Token.CATCH */, int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */, int32(49) /* Token.PARAM_LIST */:
		return true
	case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */:
		return (NodeUtil_IsEnhancedFor(parent) && (parent.GetFirstChild() == this_.nameNode))
	case int32(102) /* Token.ARRAY_PATTERN */, int32(100) /* Token.STRING_KEY */, int32(129) /* Token.COMPUTED_PROP */:
		return NodeUtil_IsLhsByDestructuring(this_.nameNode)
	default:
		return ((NodeUtil_IsAssignmentOp(parent) && (parent.GetFirstChild() == this_.nameNode)))
	}
	panic("unreachable switch expression")
}

// Java method: Reference.getScope()
func (this_ *Reference) GetScope() *Scope {
	return this_.scope
}

