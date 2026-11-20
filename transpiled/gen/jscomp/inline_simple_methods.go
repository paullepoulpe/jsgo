/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/InlineSimpleMethods.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.InlineSimpleMethods
type InlineSimpleMethods struct {
	nonInlineableProperties j2gort.OrderedSet[string]                        `gson:"nonInlineableProperties"`
	methodDefinitions       j2gort.Ext_com_google_common_collect_SetMultimap `gson:"methodDefinitions"`
	compiler                AbstractCompiler                                 `gson:"compiler"`
	astAnalyzer             *AstAnalyzer                                     `gson:"astAnalyzer"`
	idHash_                 atomic.Int32
}

func (this_ *InlineSimpleMethods) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineSimpleMethods) IsJavaNil_() bool { return this_ == nil }
func (*InlineSimpleMethods) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineSimpleMethods", "InlineSimpleMethods"
}
func (this_ *InlineSimpleMethods) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineSimpleMethods) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineSimpleMethods", "com.google.javascript.jscomp.InlineSimpleMethods", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineSimpleMethods](v_)
	})
}
func (*InlineSimpleMethods) Implements_jscomp_CompilerPass() {}

func (this_ *InlineSimpleMethods) initFields_InlineSimpleMethods() {
	_ = this_
	this_.nonInlineableProperties = j2gort.NewOrderedSet[string]()
	this_.methodDefinitions = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_LinkedHashMultimap_Create_0())
}

// Java constructor: InlineSimpleMethods(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewInlineSimpleMethods(compiler AbstractCompiler) *InlineSimpleMethods {
	this_ := &InlineSimpleMethods{}
	InitInlineSimpleMethods(this_, compiler)
	return this_
}

func InitInlineSimpleMethods(this_ *InlineSimpleMethods, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_InlineSimpleMethods()
	this_.compiler = compiler
	this_.astAnalyzer = compiler.GetAstAnalyzer()
	j2gort.OrderedSetAddAllStr(&this_.nonInlineableProperties, j2gort.AsIface[j2gort.Ext_java_util_Collection](compiler.GetExternProperties()))
}

// Java method: InlineSimpleMethods.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineSimpleMethods) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(this_.compiler.GetLifeCycleStage()))
	}
	if !(this_.methodDefinitions.IsEmpty()) {
		j2gort.CheckState(false)
	}
	if !(externs != 0) {
		j2gort.CheckState(false)
	}
	NodeTraversal_TraverseRoots(this_.compiler, NewInlineSimpleMethods_GatherSignatures(this_), externs, root)
	NodeTraversal_Traverse(this_.compiler, root, NewInlineSimpleMethods_InlineTrivialAccessors(this_))
}

// Java static method: InlineSimpleMethods.isPropertyTree(com.google.javascript.rhino.Node expectedGetprop)
func InlineSimpleMethods_IsPropertyTree(expectedGetprop pkg_rhino.Node) bool {
	_ = expectedGetprop
	if !(expectedGetprop.IsGetProp()) {
		return false
	}
	var getpropLhs pkg_rhino.Node = expectedGetprop.GetFirstChild()
	_ = getpropLhs
	return (getpropLhs.IsThis() || InlineSimpleMethods_IsPropertyTree(getpropLhs))
}

// Java static method: InlineSimpleMethods.replaceThis(com.google.javascript.rhino.Node expectedGetprop, com.google.javascript.rhino.Node replacement)
func InlineSimpleMethods_ReplaceThis(expectedGetprop pkg_rhino.Node, replacement pkg_rhino.Node) {
	_ = expectedGetprop
	_ = replacement
	var leftChild pkg_rhino.Node = expectedGetprop.GetFirstChild()
	_ = leftChild
	if leftChild.IsThis() {
		leftChild.ReplaceWith(replacement)
	} else {
		InlineSimpleMethods_ReplaceThis(leftChild, replacement)
	}
}

// Java static method: InlineSimpleMethods.returnedExpression(com.google.javascript.rhino.Node fn)
func InlineSimpleMethods_ReturnedExpression(fn pkg_rhino.Node) pkg_rhino.Node {
	_ = fn
	var expectedBlock pkg_rhino.Node = NodeUtil_GetFunctionBody(fn)
	_ = expectedBlock
	if !(expectedBlock.HasOneChild()) {
		return pkg_rhino.Node(0)
	}
	var expectedReturn pkg_rhino.Node = expectedBlock.GetFirstChild()
	_ = expectedReturn
	if !(expectedReturn.IsReturn()) {
		return pkg_rhino.Node(0)
	}
	if !(expectedReturn.HasOneChild()) {
		return pkg_rhino.Node(0)
	}
	return expectedReturn.GetOnlyChild()
}

// Java static method: InlineSimpleMethods.isEmptyMethod(com.google.javascript.rhino.Node fn)
func InlineSimpleMethods_IsEmptyMethod(fn pkg_rhino.Node) bool {
	_ = fn
	return NodeUtil_IsEmptyBlock(NodeUtil_GetFunctionBody(fn))
}

// Java method: InlineSimpleMethods.allDefinitionsEquivalent(java.util.Set<com.google.javascript.rhino.Node> definitions)
func (this_ *InlineSimpleMethods) priv_InlineSimpleMethods_allDefinitionsEquivalent_3368(definitions j2gort.Ext_java_util_Set) bool {
	_ = definitions
	var first pkg_rhino.Node = pkg_rhino.Node(0)
	_ = first
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(definitions))
		for cur_1.HasNext() {
			n := pkg_rhino.Node(cur_1.NextNode())
			_ = n
			if first == 0 {
				first = n
			} else if !(this_.compiler.AreNodesEqualForInlining(first, n)) {
				return false
			}
		}
	}
	return true
}

// Java method: InlineSimpleMethods.anyDefinitionsNotInlineable(java.util.Set<com.google.javascript.rhino.Node> definitions)
func (this_ *InlineSimpleMethods) priv_InlineSimpleMethods_anyDefinitionsNotInlineable_3369(definitions j2gort.Ext_java_util_Set) bool {
	_ = definitions
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(definitions))
		for cur_2.HasNext() {
			n := pkg_rhino.Node(cur_2.NextNode())
			_ = n
			if n.IsAsyncFunction() || n.IsGeneratorFunction() {
				return true
			}
			if NodeUtil_FunctionParametersMayHaveSideEffects(n, this_.astAnalyzer) {
				return true
			}
			var jsDocInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n.GetParent())
			_ = jsDocInfo
			if (jsDocInfo != nil) && jsDocInfo.IsNoInline() {
				return true
			}
		}
	}
	return false
}

// Java method: InlineSimpleMethods.inlinePropertyReturn(com.google.javascript.rhino.Node call, com.google.javascript.rhino.Node returnedValue)
func (this_ *InlineSimpleMethods) priv_InlineSimpleMethods_inlinePropertyReturn_3370(call pkg_rhino.Node, returnedValue pkg_rhino.Node) {
	_ = call
	_ = returnedValue
	var getProp pkg_rhino.Node = returnedValue.CloneTree_0()
	_ = getProp
	InlineSimpleMethods_ReplaceThis(getProp, call.GetFirstChild().RemoveFirstChild())
	call.ReplaceWith(getProp)
	this_.compiler.ReportChangeToEnclosingScope(getProp)
}

// Java method: InlineSimpleMethods.inlineConstReturn(com.google.javascript.rhino.Node call, com.google.javascript.rhino.Node returnedValue)
func (this_ *InlineSimpleMethods) priv_InlineSimpleMethods_inlineConstReturn_3371(call pkg_rhino.Node, returnedValue pkg_rhino.Node) {
	_ = call
	_ = returnedValue
	var retValue pkg_rhino.Node = returnedValue.CloneTree_0()
	_ = retValue
	call.ReplaceWith(retValue)
	this_.compiler.ReportChangeToEnclosingScope(retValue)
}

// Java method: InlineSimpleMethods.inlineEmptyMethod(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node call)
func (this_ *InlineSimpleMethods) priv_InlineSimpleMethods_inlineEmptyMethod_3372(t *NodeTraversal, parent pkg_rhino.Node, call pkg_rhino.Node) {
	_ = t
	_ = parent
	_ = call
	if NodeUtil_IsExprCall(parent) {
		parent.ReplaceWith(pkg_rhino.IR_Empty())
		NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(parent, this_.compiler)
	} else {
		var srcLocation pkg_rhino.Node = call
		_ = srcLocation
		call.ReplaceWith(NodeUtil_NewUndefinedNode(srcLocation))
		NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(call, this_.compiler)
	}
	t.ReportCodeChange_0()
}

// Java method: InlineSimpleMethods.argsMayHaveSideEffects(com.google.javascript.rhino.Node call)
func (this_ *InlineSimpleMethods) priv_InlineSimpleMethods_argsMayHaveSideEffects_3373(call pkg_rhino.Node) bool {
	_ = call
	{
		var currentChild pkg_rhino.Node = call.GetSecondChild()
		_ = currentChild
		for ; currentChild != 0; currentChild = currentChild.GetNext() {
			if this_.astAnalyzer.MayHaveSideEffects(currentChild) {
				return true
			}
		}
	}
	return false
}

// Java method: InlineSimpleMethods.addPossibleSignature(java.lang.String name, com.google.javascript.rhino.Node node)
func (this_ *InlineSimpleMethods) priv_InlineSimpleMethods_addPossibleSignature_3374(name string, node pkg_rhino.Node) {
	_ = name
	_ = node
	if (node != 0) && node.IsFunction() {
		this_.priv_InlineSimpleMethods_addSignature_3375(name, node)
	} else {
		this_.nonInlineableProperties.Add(name)
	}
}

// Java method: InlineSimpleMethods.addSignature(java.lang.String name, com.google.javascript.rhino.Node function)
func (this_ *InlineSimpleMethods) priv_InlineSimpleMethods_addSignature_3375(name string, function pkg_rhino.Node) {
	_ = name
	_ = function
	if this_.nonInlineableProperties.Contains(name) {
		return
	}
	j2gort.MultimapPutStr(this_.methodDefinitions, name, j2gort.BoxNode(function))
}

// Java static method: InlineSimpleMethods.getPropName(com.google.javascript.rhino.Node getPropElem)
func InlineSimpleMethods_GetPropName(getPropElem pkg_rhino.Node) string {
	_ = getPropElem
	if getPropElem.IsGetProp() {
		return getPropElem.GetString()
	} else if getPropElem.GetSecondChild().IsStringLit() {
		return getPropElem.GetSecondChild().GetString()
	} else {
		return j2gort.NullString
	}
}

// Java class: com.google.javascript.jscomp.InlineSimpleMethods.InlineTrivialAccessors
type InlineSimpleMethods_InlineTrivialAccessors struct {
	InvocationsCallbackBase
	outer_ *InlineSimpleMethods
}

func (this_ *InlineSimpleMethods_InlineTrivialAccessors) BindSelf_(s any) {
	_ = s
	this_.InvocationsCallbackBase.BindSelf_(s)
}
func (this_ *InlineSimpleMethods_InlineTrivialAccessors) IsJavaNil_() bool { return this_ == nil }
func (*InlineSimpleMethods_InlineTrivialAccessors) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineSimpleMethods$InlineTrivialAccessors", "InlineTrivialAccessors"
}
func (this_ *InlineSimpleMethods_InlineTrivialAccessors) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.InvocationsCallbackBase.IdentityHash_()
}
func (*InlineSimpleMethods_InlineTrivialAccessors) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineSimpleMethods.InlineTrivialAccessors", "com.google.javascript.jscomp.InlineSimpleMethods$InlineTrivialAccessors", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineSimpleMethods_InlineTrivialAccessors](v_)
	})
}
func (this_ *InlineSimpleMethods_InlineTrivialAccessors) SetOuter_(o *InlineSimpleMethods) {
	this_.outer_ = o
}
func (this_ *InlineSimpleMethods_InlineTrivialAccessors) GetOuter_() *InlineSimpleMethods {
	return this_.outer_
}
func (*InlineSimpleMethods_InlineTrivialAccessors) Implements_jscomp_InvocationsCallback() {}
func (*InlineSimpleMethods_InlineTrivialAccessors) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*InlineSimpleMethods_InlineTrivialAccessors) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *InlineSimpleMethods_InlineTrivialAccessors) initFields_InlineSimpleMethods_InlineTrivialAccessors() {
	_ = this_
}

// Java constructor: InlineTrivialAccessors()
func NewInlineSimpleMethods_InlineTrivialAccessors(outer_ *InlineSimpleMethods) *InlineSimpleMethods_InlineTrivialAccessors {
	this_ := &InlineSimpleMethods_InlineTrivialAccessors{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitInlineSimpleMethods_InlineTrivialAccessors(this_)
	return this_
}

func InitInlineSimpleMethods_InlineTrivialAccessors(this_ *InlineSimpleMethods_InlineTrivialAccessors) {
	InitInvocationsCallback(&this_.InvocationsCallbackBase)
	this_.initFields_InlineSimpleMethods_InlineTrivialAccessors()
}

// Java method: InlineTrivialAccessors.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node callNode, com.google.javascript.rhino.Node parent, java.lang.String callName)
func (this_ *InlineSimpleMethods_InlineTrivialAccessors) Visit_4(t *NodeTraversal, callNode pkg_rhino.Node, parent pkg_rhino.Node, callName string) {
	_ = t
	_ = callNode
	_ = parent
	_ = callName
	if this_.outer_.nonInlineableProperties.Contains(callName) {
		return
	}
	var definitions j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.MultimapGetStr(this_.outer_.methodDefinitions, callName))
	_ = definitions
	if (definitions == nil) || definitions.IsEmpty() {
		return
	}
	if (definitions.Size() > int32(1)) && !(this_.outer_.priv_InlineSimpleMethods_allDefinitionsEquivalent_3368(definitions)) {
		this_.outer_.nonInlineableProperties.Add(callName)
		return
	}
	if this_.outer_.priv_InlineSimpleMethods_anyDefinitionsNotInlineable_3369(definitions) {
		this_.outer_.nonInlineableProperties.Add(callName)
		return
	}
	var firstDefinition pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](j2gort.AsIface[j2gort.Ext_java_util_Iterator](definitions.Iterator()).Next())
	_ = firstDefinition
	if !(this_.outer_.priv_InlineSimpleMethods_argsMayHaveSideEffects_3373(callNode)) && !(NodeUtil_ReferencesSuper(callNode)) {
		var returned pkg_rhino.Node = InlineSimpleMethods_ReturnedExpression(firstDefinition)
		_ = returned
		if returned != 0 {
			if InlineSimpleMethods_IsPropertyTree(returned) && !(firstDefinition.IsArrowFunction()) {
				this_.outer_.priv_InlineSimpleMethods_inlinePropertyReturn_3370(callNode, returned)
			} else if NodeUtil_IsLiteralValue(returned, false) && !(this_.outer_.astAnalyzer.MayHaveSideEffects(callNode.GetFirstChild())) {
				this_.outer_.priv_InlineSimpleMethods_inlineConstReturn_3371(callNode, returned)
			}
		} else if InlineSimpleMethods_IsEmptyMethod(firstDefinition) && !(this_.outer_.astAnalyzer.MayHaveSideEffects(callNode.GetFirstChild())) {
			this_.outer_.priv_InlineSimpleMethods_inlineEmptyMethod_3372(t, parent, callNode)
		}
	}
}

// Java class: com.google.javascript.jscomp.InlineSimpleMethods.GatherSignatures
type InlineSimpleMethods_GatherSignatures struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *InlineSimpleMethods
}

func (this_ *InlineSimpleMethods_GatherSignatures) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineSimpleMethods_GatherSignatures) IsJavaNil_() bool { return this_ == nil }
func (*InlineSimpleMethods_GatherSignatures) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineSimpleMethods$GatherSignatures", "GatherSignatures"
}
func (this_ *InlineSimpleMethods_GatherSignatures) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*InlineSimpleMethods_GatherSignatures) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineSimpleMethods.GatherSignatures", "com.google.javascript.jscomp.InlineSimpleMethods$GatherSignatures", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineSimpleMethods_GatherSignatures](v_)
	})
}
func (this_ *InlineSimpleMethods_GatherSignatures) SetOuter_(o *InlineSimpleMethods) {
	this_.outer_ = o
}
func (this_ *InlineSimpleMethods_GatherSignatures) GetOuter_() *InlineSimpleMethods {
	return this_.outer_
}
func (*InlineSimpleMethods_GatherSignatures) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*InlineSimpleMethods_GatherSignatures) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *InlineSimpleMethods_GatherSignatures) initFields_InlineSimpleMethods_GatherSignatures() {
	_ = this_
}

// Java constructor: GatherSignatures()
func NewInlineSimpleMethods_GatherSignatures(outer_ *InlineSimpleMethods) *InlineSimpleMethods_GatherSignatures {
	this_ := &InlineSimpleMethods_GatherSignatures{}
	this_.outer_ = outer_
	InitInlineSimpleMethods_GatherSignatures(this_)
	return this_
}

func InitInlineSimpleMethods_GatherSignatures(this_ *InlineSimpleMethods_GatherSignatures) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_InlineSimpleMethods_GatherSignatures()
}

// Java method: GatherSignatures.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineSimpleMethods_GatherSignatures) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
		var name string = InlineSimpleMethods_GetPropName(n)
		_ = name
		if name == j2gort.NullString {
			return
		}
		if j2gort.CheckNotNullString(name) == "prototype" {
			this_.priv_GatherSignatures_processPrototypeParent_3377(parent)
		} else {
			if parent.IsAssign() && n.IsFirstChildOf(parent) {
				this_.outer_.priv_InlineSimpleMethods_addPossibleSignature_3374(name, n.GetNext())
			}
		}
	case int32(47) /* Token.OBJECTLIT */, int32(106) /* Token.CLASS_MEMBERS */ :
		{
			var key pkg_rhino.Node = n.GetFirstChild()
			_ = key
			for ; key != 0; key = key.GetNext() {
				switch key.TokenSwitchOrd_() {
				case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(100) /* Token.STRING_KEY */ :
					this_.outer_.priv_InlineSimpleMethods_addPossibleSignature_3374(key.GetString(), key.GetFirstChild())
				case int32(96) /* Token.SETTER_DEF */, int32(95) /* Token.GETTER_DEF */ :
					this_.outer_.nonInlineableProperties.Add(key.GetString())
				case int32(129) /* Token.COMPUTED_PROP */, int32(128) /* Token.OBJECT_SPREAD */, int32(109) /* Token.COMPUTED_FIELD_DEF */, int32(90) /* Token.BLOCK */ :
				default:
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Unexpected " + j2gort.StrVal(n.GetToken()) + " key: " + j2gort.StrVal(key))))))
				}
			}
		}
	case int32(28) /* Token.CALL */ :
		if this_.outer_.compiler.GetCodingConvention().IsPropertyRenameFunction(n.GetFirstChild()) {
			this_.outer_.nonInlineableProperties.Add(n.GetSecondChild().GetString())
		}
	default:
	}
}

// Java method: GatherSignatures.processPrototypeParent(com.google.javascript.rhino.Node n)
func (this_ *InlineSimpleMethods_GatherSignatures) priv_GatherSignatures_processPrototypeParent_3377(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
		var grandparent pkg_rhino.Node = n.GetGrandparent()
		_ = grandparent
		var name string = InlineSimpleMethods_GetPropName(n)
		_ = name
		if (name != j2gort.NullString) && grandparent.IsAssign() {
			var assignee pkg_rhino.Node = grandparent.GetSecondChild()
			_ = assignee
			this_.outer_.priv_InlineSimpleMethods_addPossibleSignature_3374(name, assignee)
		}
	default:
	}
}
func (*InlineSimpleMethods_InlineTrivialAccessors) AlwaysShouldTraverse_()      {}
func (*InlineSimpleMethods_InlineTrivialAccessors) CanReleaseSyntacticScopes_() {}
func (*InlineSimpleMethods_GatherSignatures) AlwaysShouldTraverse_()            {}
func (*InlineSimpleMethods_GatherSignatures) CanReleaseSyntacticScopes_()       {}
