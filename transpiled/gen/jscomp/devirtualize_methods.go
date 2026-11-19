/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/DevirtualizeMethods.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_devirtualize_methods_1 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_devirtualize_methods_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return pkg_rhino.Node(rawNode_).IsSuper() })
}

// Java class: com.google.javascript.jscomp.DevirtualizeMethods
type DevirtualizeMethods struct {
	compiler                    AbstractCompiler            `gson:"compiler"`
	refMap                      *OptimizeCalls_ReferenceMap `gson:"refMap"`
	idHash_                     atomic.Int32
	slab_OptimizeCalls_NodeList j2gort.PooledSlab[OptimizeCalls_NodeList]
}

func (this_ *DevirtualizeMethods) BindSelf_(s any) {
	_ = s
}
func (this_ *DevirtualizeMethods) IsJavaNil_() bool { return this_ == nil }
func (*DevirtualizeMethods) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DevirtualizeMethods", "DevirtualizeMethods"
}
func (this_ *DevirtualizeMethods) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DevirtualizeMethods) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DevirtualizeMethods", "com.google.javascript.jscomp.DevirtualizeMethods", func(v_ any) bool {
		return j2gort.InstanceOf[*DevirtualizeMethods](v_)
	})
}
func (*DevirtualizeMethods) Implements_jscomp_OptimizeCalls_CallGraphCompilerPass() {}

func (this_ *DevirtualizeMethods) initFields_DevirtualizeMethods() {
	_ = this_
}

// Java constructor: DevirtualizeMethods(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewDevirtualizeMethods(compiler AbstractCompiler) *DevirtualizeMethods {
	this_ := &DevirtualizeMethods{}
	InitDevirtualizeMethods(this_, compiler)
	return this_
}

func InitDevirtualizeMethods(this_ *DevirtualizeMethods, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_DevirtualizeMethods()
	this_.compiler = compiler
}

// Java method: DevirtualizeMethods.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root, com.google.javascript.jscomp.OptimizeCalls.ReferenceMap refMap)
func (this_ *DevirtualizeMethods) Process_3(externs pkg_rhino.Node, root pkg_rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
	_ = externs
	_ = root
	_ = refMap
	defer this_.ReleaseSlabs_()
	if !(this_.refMap == nil) {
		j2gort.CheckStateStr(false, "`process` should only be called once.")
	}
	this_.refMap = refMap
	{
		cur_1 := refMap.props.IterEntries()
		for cur_1.NextEntry() {
			this_.priv_DevirtualizeMethods_processReferenceList_2418(cur_1.Key(), cur_1.Value())
		}
	}
}

// Java method: DevirtualizeMethods.processReferenceList(java.lang.String name, com.google.javascript.jscomp.OptimizeCalls.NodeList sites)
func (this_ *DevirtualizeMethods) priv_DevirtualizeMethods_processReferenceList_2418(name string, sites *OptimizeCalls_NodeList) {
	_ = name
	_ = sites
	defer this_.slab_OptimizeCalls_NodeList.Reset()
	var defSites *OptimizeCalls_NodeList = nil
	_ = defSites
	var fns *OptimizeCalls_NodeList = nil
	_ = fns
	{
		var i int32 = int32(0)
		_ = i
		var count int32 = sites.Size()
		_ = count
		for ; i < count; i = i + 1 {
			var s pkg_rhino.Node = sites.Get(i)
			_ = s
			var prevCount int32
			_ = prevCount
			if fns == nil {
				prevCount = int32(0)
			} else {
				prevCount = fns.Size()
			}
			fns = allocPooledVia_OptimizeCalls_ReferenceMap_AppendDefinitionFunctionNodes(&this_.slab_OptimizeCalls_NodeList, s, fns)
			var newCount int32
			_ = newCount
			if fns == nil {
				newCount = int32(0)
			} else {
				newCount = fns.Size()
			}
			if newCount > prevCount {
				if defSites == nil {
					defSites = allocPooled_OptimizeCalls_NodeList(&this_.slab_OptimizeCalls_NodeList)
				}
				if !(defSites.Contains(s)) {
					defSites.Add(s)
				}
			}
		}
	}
	if ((defSites == nil) || (fns == nil)) || fns.IsEmpty() {
		return
	}
	var canonicalDefinitionSite pkg_rhino.Node = defSites.Get(int32(0))
	_ = canonicalDefinitionSite
	var callSites *OptimizeCalls_NodeList = nil
	_ = callSites
	{
		var i_2 int32 = int32(0)
		_ = i_2
		var count_2 int32 = sites.Size()
		_ = count_2
		for ; i_2 < count_2; i_2 = i_2 + 1 {
			var s_2 pkg_rhino.Node = sites.Get(i_2)
			_ = s_2
			if !(defSites.Contains(s_2)) {
				if callSites == nil {
					callSites = allocPooled_OptimizeCalls_NodeList(&this_.slab_OptimizeCalls_NodeList)
				}
				callSites.Add(s_2)
			}
		}
	}
	if (callSites == nil) || callSites.IsEmpty() {
		return
	}
	{
		var i_3 int32 = int32(0)
		_ = i_3
		var count_3 int32 = defSites.Size()
		_ = count_3
		for ; i_3 < count_3; i_3 = i_3 + 1 {
			if !(this_.priv_DevirtualizeMethods_isEligibleDefinitionSite_2421(name, defSites.Get(i_3))) {
				return
			}
		}
	}
	{
		var i_4 int32 = int32(0)
		_ = i_4
		var count_4 int32 = fns.Size()
		_ = count_4
		for ; i_4 < count_4; i_4 = i_4 + 1 {
			if !(this_.priv_DevirtualizeMethods_isEligibleDefinitionFunction_2422(fns.Get(i_4))) {
				return
			}
		}
	}
	{
		var i_5 int32 = int32(0)
		_ = i_5
		var count_5 int32 = callSites.Size()
		_ = count_5
		for ; i_5 < count_5; i_5 = i_5 + 1 {
			if !(this_.priv_DevirtualizeMethods_isEligibleCallSite_2423(callSites.Get(i_5), canonicalDefinitionSite)) {
				return
			}
		}
	}
	if !(this_.priv_DevirtualizeMethods_allDefinitionsEquivalent_2424(fns)) {
		return
	}
	var devirtualizedName string = DevirtualizeMethods_RewrittenMethodNameOf(name)
	_ = devirtualizedName
	var bestJSDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(canonicalDefinitionSite)
	_ = bestJSDoc
	var isConstantName bool = ((bestJSDoc != nil) && bestJSDoc.IsConstant())
	_ = isConstantName
	{
		var i_6 int32 = int32(0)
		_ = i_6
		var count_6 int32 = callSites.Size()
		_ = count_6
		for ; i_6 < count_6; i_6 = i_6 + 1 {
			this_.priv_DevirtualizeMethods_rewriteCall_2425(callSites.Get(i_6), devirtualizedName, isConstantName)
		}
	}
	this_.priv_DevirtualizeMethods_rewriteDefinition_2426(canonicalDefinitionSite, devirtualizedName, isConstantName)
}

// Java method: DevirtualizeMethods.isPrototypeOrStaticMethodDefinition(com.google.javascript.rhino.Node node)
func (this_ *DevirtualizeMethods) priv_DevirtualizeMethods_isPrototypeOrStaticMethodDefinition_2419(node pkg_rhino.Node) bool {
	_ = node
	var parent pkg_rhino.Node = node.GetParent()
	_ = parent
	var grandparent pkg_rhino.Node = node.GetGrandparent()
	_ = grandparent
	if (parent == 0) || (grandparent == 0) {
		return false
	}
	switch node.TokenSwitchOrd_() {
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		if NodeUtil_IsEs6ConstructorMemberFunctionDef(node) {
			return false
		}
		if parent.IsObjectLit() {
			return DevirtualizeMethods_IsPrototypeObjectLit(parent)
		}
		return true
	case int32(26) /* Token.GETPROP */ :
		if (!(node.IsFirstChildOf(parent)) || !(NodeUtil_IsExprAssign(grandparent))) || !(parent.GetLastChild().IsFunction()) {
			return false
		}
		if NodeUtil_IsPrototypeProperty(node) {
			return true
		}
		if this_.priv_DevirtualizeMethods_isDefinitelyCtorOrInterface_2420(node.GetFirstChild()) {
			return true
		}
		return false
	case int32(100) /* Token.STRING_KEY */ :
		if !(parent.IsObjectLit()) {
			j2gort.CheckArgumentMsg(false, j2gort.BoxNode(parent))
		}
		return (DevirtualizeMethods_IsPrototypeObjectLit(parent) && node.GetFirstChild().IsFunction())
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: DevirtualizeMethods.isPrototypeObjectLit(com.google.javascript.rhino.Node objectLit)
func DevirtualizeMethods_IsPrototypeObjectLit(objectLit pkg_rhino.Node) bool {
	_ = objectLit
	if !(objectLit.IsObjectLit()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(objectLit))
	}
	var parent pkg_rhino.Node = objectLit.GetParent()
	_ = parent
	return (objectLit.IsSecondChildOf(parent) && NodeUtil_IsPrototypeAssignment(parent.GetFirstChild()))
}

// Java method: DevirtualizeMethods.isDefinitelyCtorOrInterface(com.google.javascript.rhino.Node receiver)
func (this_ *DevirtualizeMethods) priv_DevirtualizeMethods_isDefinitelyCtorOrInterface_2420(receiver pkg_rhino.Node) bool {
	_ = receiver
	var qname string = receiver.GetQualifiedName()
	_ = qname
	if qname == j2gort.NullString {
		return false
	}
	var var_ *Var = j2gort.CastPtr[Var](this_.refMap.GetGlobalScope().GetVar(qname))
	_ = var_
	if var_ == nil {
		return false
	}
	if var_.IsClass() {
		return true
	}
	var jsdoc *pkg_rhino.JSDocInfo = var_.GetJSDocInfo()
	_ = jsdoc
	if jsdoc == nil {
		return false
	} else if jsdoc.IsConstructorOrInterface() {
		return true
	}
	return false
}

// Java method: DevirtualizeMethods.isEligibleDefinitionSite(java.lang.String name, com.google.javascript.rhino.Node definitionSite)
func (this_ *DevirtualizeMethods) priv_DevirtualizeMethods_isEligibleDefinitionSite_2421(name string, definitionSite pkg_rhino.Node) bool {
	_ = name
	_ = definitionSite
	switch definitionSite.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(100) /* Token.STRING_KEY */, int32(108) /* Token.MEMBER_FIELD_DEF */ :
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.BoxString(definitionSite.String()))))
	}
	var codingConvention CodingConvention = this_.compiler.GetCodingConvention()
	_ = codingConvention
	if codingConvention.IsExported_2(name, false) {
		return false
	}
	if !(this_.priv_DevirtualizeMethods_isPrototypeOrStaticMethodDefinition_2419(definitionSite)) {
		return false
	}
	return true
}

// Java method: DevirtualizeMethods.isEligibleDefinitionFunction(com.google.javascript.rhino.Node definitionFunction)
func (this_ *DevirtualizeMethods) priv_DevirtualizeMethods_isEligibleDefinitionFunction_2422(definitionFunction pkg_rhino.Node) bool {
	_ = definitionFunction
	if !(definitionFunction.IsFunction()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(definitionFunction))
	}
	if definitionFunction.IsArrowFunction() {
		return false
	}
	{
		var ancestor pkg_rhino.Node = definitionFunction.GetParent()
		_ = ancestor
		for ; ancestor != 0; ancestor = ancestor.GetParent() {
			if DevirtualizeMethods_IsScopingOrBranchingConstruct(ancestor) {
				return false
			}
			if ancestor.IsClass() && DevirtualizeMethods_LocalNameIsDeclaredByClass(ancestor) {
				return false
			}
		}
	}
	if NodeUtil_Has(definitionFunction, lambdaSingleton_devirtualize_methods_1, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue()) {
		return false
	}
	if NodeUtil_DoesFunctionReferenceOwnArgumentsObject(definitionFunction) {
		return false
	}
	return true
}

// Java method: DevirtualizeMethods.isEligibleCallSite(com.google.javascript.rhino.Node access, com.google.javascript.rhino.Node definitionSite)
func (this_ *DevirtualizeMethods) priv_DevirtualizeMethods_isEligibleCallSite_2423(access pkg_rhino.Node, definitionSite pkg_rhino.Node) bool {
	_ = access
	_ = definitionSite
	var invocation pkg_rhino.Node = access.GetParent()
	_ = invocation
	if !(NodeUtil_IsInvocationTarget(access)) || !(invocation.IsCall()) {
		return false
	}
	var chunkGraph *JSChunkGraph = this_.compiler.GetChunkGraph()
	_ = chunkGraph
	var definitionChunk *JSChunk = this_.priv_DevirtualizeMethods_chunkForNode_2428(definitionSite)
	_ = definitionChunk
	var callChunk *JSChunk = this_.priv_DevirtualizeMethods_chunkForNode_2428(access)
	_ = callChunk
	if definitionChunk == callChunk {
	} else if callChunk == nil {
		return false
	} else if !(chunkGraph.DependsOn(callChunk, definitionChunk)) {
		return false
	}
	return true
}

// Java method: DevirtualizeMethods.allDefinitionsEquivalent(com.google.javascript.jscomp.OptimizeCalls.NodeList definitions)
func (this_ *DevirtualizeMethods) priv_DevirtualizeMethods_allDefinitionsEquivalent_2424(definitions *OptimizeCalls_NodeList) bool {
	_ = definitions
	var count int32 = definitions.Size()
	_ = count
	if count == int32(0) {
		return true
	}
	var definition pkg_rhino.Node = definitions.Get(int32(0))
	_ = definition
	if !(definition.IsFunction()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(definition))
	}
	{
		var i int32 = int32(1)
		_ = i
		for ; i < count; i = i + 1 {
			if !(this_.compiler.AreNodesEqualForInlining(definitions.Get(i), definition)) {
				return false
			}
		}
	}
	return true
}

// Java method: DevirtualizeMethods.rewriteCall(com.google.javascript.rhino.Node getprop, java.lang.String newMethodName, boolean isConstantName)
func (this_ *DevirtualizeMethods) priv_DevirtualizeMethods_rewriteCall_2425(getprop pkg_rhino.Node, newMethodName string, isConstantName bool) {
	_ = getprop
	_ = newMethodName
	_ = isConstantName
	if !(getprop.IsGetProp()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(getprop))
	}
	var call pkg_rhino.Node = getprop.GetParent()
	_ = call
	if !(call.IsCall()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(call))
	}
	var receiver pkg_rhino.Node = getprop.GetFirstChild()
	_ = receiver
	receiver.Detach()
	getprop.ReplaceWith(receiver)
	var newReceiver pkg_rhino.Node = pkg_rhino.IR_Name(newMethodName).CopyTypeFrom(getprop).Srcref(getprop)
	_ = newReceiver
	if isConstantName {
		newReceiver.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
	call.AddChildToFront(newReceiver)
	if receiver.IsSuper() {
		receiver.SetToken(pkg_rhino.StaticGet_Token_THIS())
	}
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	this_.compiler.ReportChangeToEnclosingScope(call)
}

// Java method: DevirtualizeMethods.rewriteDefinition(com.google.javascript.rhino.Node definitionSite, java.lang.String newMethodName, boolean isConstantName)
func (this_ *DevirtualizeMethods) priv_DevirtualizeMethods_rewriteDefinition_2426(definitionSite pkg_rhino.Node, newMethodName string, isConstantName bool) {
	_ = definitionSite
	_ = newMethodName
	_ = isConstantName
	var function pkg_rhino.Node = pkg_rhino.Node(0)
	_ = function
	var subtreeToRemove pkg_rhino.Node = pkg_rhino.Node(0)
	_ = subtreeToRemove
	var nameSource pkg_rhino.Node = pkg_rhino.Node(0)
	_ = nameSource
	switch definitionSite.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */ :
		function = definitionSite.GetParent().GetLastChild()
		nameSource = definitionSite
		subtreeToRemove = NodeUtil_GetEnclosingStatement(definitionSite)
	case int32(100) /* Token.STRING_KEY */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		function = definitionSite.GetLastChild()
		nameSource = definitionSite
		subtreeToRemove = definitionSite
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.BoxString(definitionSite.String()))))
	}
	var statement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(definitionSite)
	_ = statement
	var newNameNode pkg_rhino.Node = pkg_rhino.IR_Name(newMethodName).SrcrefIfMissing(nameSource)
	_ = newNameNode
	var newVarNode pkg_rhino.Node = pkg_rhino.IR_Var_1(newNameNode).SrcrefIfMissing(nameSource)
	_ = newVarNode
	newVarNode.InsertBefore(statement)
	if isConstantName {
		newNameNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
	var originalJSDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(definitionSite)
	_ = originalJSDoc
	newVarNode.SetJSDocInfo(originalJSDoc)
	function.Detach()
	newNameNode.AddChildToFront(function)
	var selfName string = (j2gort.StrValString(newMethodName) + "$self")
	_ = selfName
	var paramList pkg_rhino.Node = function.GetSecondChild()
	_ = paramList
	paramList.AddChildToFront(pkg_rhino.IR_Name(selfName).SrcrefIfMissing(function))
	this_.compiler.ReportChangeToEnclosingScope(paramList)
	this_.priv_DevirtualizeMethods_replaceReferencesToThis_2427(function.GetSecondChild(), selfName)
	this_.priv_DevirtualizeMethods_replaceReferencesToThis_2427(function.GetLastChild(), selfName)
	NodeUtil_DeleteNode(subtreeToRemove, this_.compiler)
	this_.compiler.ReportChangeToEnclosingScope(newVarNode)
}

// Java method: DevirtualizeMethods.replaceReferencesToThis(com.google.javascript.rhino.Node node, java.lang.String name)
func (this_ *DevirtualizeMethods) priv_DevirtualizeMethods_replaceReferencesToThis_2427(node pkg_rhino.Node, name string) {
	_ = node
	_ = name
	if ((node.IsFunction() && !(node.IsArrowFunction())) || node.IsMemberFieldDef()) || NodeUtil_IsClassStaticBlock(node) {
		return
	}
	if node.IsComputedFieldDef() {
		var key pkg_rhino.Node = node.GetFirstChild()
		_ = key
		if key != 0 {
			if key.IsThis() {
				var newName pkg_rhino.Node = pkg_rhino.IR_Name(name).Srcref(key).CopyTypeFrom(key)
				_ = newName
				key.ReplaceWith(newName)
				this_.compiler.ReportChangeToEnclosingScope(newName)
			} else {
				this_.priv_DevirtualizeMethods_replaceReferencesToThis_2427(key, name)
			}
		}
		return
	}
	{
		var child pkg_rhino.Node = node.GetFirstChild()
		_ = child
		for child != 0 {
			var next pkg_rhino.Node = child.GetNext()
			_ = next
			if child.IsThis() {
				var newName_2 pkg_rhino.Node = pkg_rhino.IR_Name(name).Srcref(child).CopyTypeFrom(child)
				_ = newName_2
				child.ReplaceWith(newName_2)
				this_.compiler.ReportChangeToEnclosingScope(newName_2)
			} else {
				this_.priv_DevirtualizeMethods_replaceReferencesToThis_2427(child, name)
			}
			child = next
		}
	}
}

// Java method: DevirtualizeMethods.chunkForNode(com.google.javascript.rhino.Node node)
func (this_ *DevirtualizeMethods) priv_DevirtualizeMethods_chunkForNode_2428(node pkg_rhino.Node) *JSChunk {
	_ = node
	var script pkg_rhino.Node = NodeUtil_GetEnclosingScript(node)
	_ = script
	var input *CompilerInput = this_.compiler.GetInput(script.GetInputId())
	_ = input
	return input.GetChunk()
}

// Java static method: DevirtualizeMethods.rewrittenMethodNameOf(java.lang.String originalMethodName)
func DevirtualizeMethods_RewrittenMethodNameOf(originalMethodName string) string {
	_ = originalMethodName
	return ("JSCompiler_StaticMethods_" + j2gort.StrValString(originalMethodName))
}

// Java static method: DevirtualizeMethods.isScopingOrBranchingConstruct(com.google.javascript.rhino.Node node)
func DevirtualizeMethods_IsScopingOrBranchingConstruct(node pkg_rhino.Node) bool {
	_ = node
	return ((((NodeUtil_IsControlStructure(node) || node.IsAnd()) || node.IsOr()) || node.IsFunction()) || node.IsBlock())
}

// Java static method: DevirtualizeMethods.localNameIsDeclaredByClass(com.google.javascript.rhino.Node clazz)
func DevirtualizeMethods_LocalNameIsDeclaredByClass(clazz pkg_rhino.Node) bool {
	_ = clazz
	if !(clazz.IsClass()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(clazz))
	}
	if clazz.GetFirstChild().IsEmpty() {
		return false
	} else if NodeUtil_IsStatement(clazz) {
		return false
	}
	return true
}
func (this_ *DevirtualizeMethods) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.refMap = nil
	this_.slab_OptimizeCalls_NodeList.Release()
}
