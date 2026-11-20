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
// Port of com/google/javascript/jscomp/CrossChunkMethodMotion.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.CrossChunkMethodMotion
type CrossChunkMethodMotion struct {
	compiler        AbstractCompiler            `gson:"compiler"`
	idGenerator     *IdGenerator                `gson:"idGenerator"`
	analyzer        *AnalyzePrototypeProperties `gson:"analyzer"`
	chunkGraph      *JSChunkGraph               `gson:"chunkGraph"`
	astFactory      *AstFactory                 `gson:"astFactory"`
	idHash_         atomic.Int32
	noStubFunctions bool `gson:"noStubFunctions"`
}

func (this_ *CrossChunkMethodMotion) BindSelf_(s any) {
	_ = s
}
func (this_ *CrossChunkMethodMotion) IsJavaNil_() bool { return this_ == nil }
func (*CrossChunkMethodMotion) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CrossChunkMethodMotion", "CrossChunkMethodMotion"
}
func (this_ *CrossChunkMethodMotion) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CrossChunkMethodMotion) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CrossChunkMethodMotion", "com.google.javascript.jscomp.CrossChunkMethodMotion", func(v_ any) bool {
		return j2gort.InstanceOf[*CrossChunkMethodMotion](v_)
	})
}
func (*CrossChunkMethodMotion) Implements_jscomp_CompilerPass() {}

var static_CrossChunkMethodMotion_NULL_COMMON_CHUNK_ERROR *DiagnosticType

func StaticGet_CrossChunkMethodMotion_NULL_COMMON_CHUNK_ERROR() *DiagnosticType {
	clinit_CrossChunkMethodMotion()
	return static_CrossChunkMethodMotion_NULL_COMMON_CHUNK_ERROR
}
func StaticSet_CrossChunkMethodMotion_NULL_COMMON_CHUNK_ERROR(v *DiagnosticType) {
	clinit_CrossChunkMethodMotion()
	static_CrossChunkMethodMotion_NULL_COMMON_CHUNK_ERROR = v
}
func StaticGet_CrossChunkMethodMotion_STUB_METHOD_NAME() string   { return "JSCompiler_stubMethod" }
func StaticGet_CrossChunkMethodMotion_UNSTUB_METHOD_NAME() string { return "JSCompiler_unstubMethod" }
func StaticGet_CrossChunkMethodMotion_STUB_DECLARATIONS() string {
	return "var JSCompiler_stubMap = [];function JSCompiler_stubMethod(JSCompiler_stubMethod_id) {  return function() {    return JSCompiler_stubMap[JSCompiler_stubMethod_id].apply(        this, arguments);  };}function JSCompiler_unstubMethod(    JSCompiler_unstubMethod_id, JSCompiler_unstubMethod_body) {  return JSCompiler_stubMap[JSCompiler_unstubMethod_id] =       JSCompiler_unstubMethod_body;}"
}

var clinitDone_CrossChunkMethodMotion uint32
var clinitState_CrossChunkMethodMotion j2gort.ClinitState

func clinit_CrossChunkMethodMotion() {
	if atomic.LoadUint32(&clinitDone_CrossChunkMethodMotion) != 0 {
		return
	}
	clinitSlow_CrossChunkMethodMotion()
}

func clinitSlow_CrossChunkMethodMotion() {
	if !j2gort.EnterClinit(&clinitState_CrossChunkMethodMotion, &clinitDone_CrossChunkMethodMotion) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CrossChunkMethodMotion, &clinitDone_CrossChunkMethodMotion)
	static_CrossChunkMethodMotion_NULL_COMMON_CHUNK_ERROR = DiagnosticType_Error("JSC_INTERNAL_ERROR_CHUNK_DEPEND", "null deepest common chunk")
}

func (this_ *CrossChunkMethodMotion) initFields_CrossChunkMethodMotion() {
	_ = this_
}

// Java constructor: CrossChunkMethodMotion(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.IdGenerator idGenerator, boolean canModifyExterns, boolean noStubFunctions)
func NewCrossChunkMethodMotion(compiler AbstractCompiler, idGenerator *IdGenerator, canModifyExterns bool, noStubFunctions bool) *CrossChunkMethodMotion {
	clinit_CrossChunkMethodMotion()
	this_ := &CrossChunkMethodMotion{}
	InitCrossChunkMethodMotion(this_, compiler, idGenerator, canModifyExterns, noStubFunctions)
	return this_
}

func InitCrossChunkMethodMotion(this_ *CrossChunkMethodMotion, compiler AbstractCompiler, idGenerator *IdGenerator, canModifyExterns bool, noStubFunctions bool) {
	_ = compiler
	_ = idGenerator
	_ = canModifyExterns
	_ = noStubFunctions
	this_.initFields_CrossChunkMethodMotion()
	this_.compiler = compiler
	this_.idGenerator = idGenerator
	this_.chunkGraph = compiler.GetChunkGraph()
	this_.analyzer = NewAnalyzePrototypeProperties(compiler, this_.chunkGraph, canModifyExterns, false, noStubFunctions)
	this_.noStubFunctions = noStubFunctions
	this_.astFactory = compiler.CreateAstFactory()
}

// Java method: CrossChunkMethodMotion.process(com.google.javascript.rhino.Node externRoot, com.google.javascript.rhino.Node root)
func (this_ *CrossChunkMethodMotion) Process_2(externRoot pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externRoot
	_ = root
	defer func() {
		if this_.analyzer != nil {
			this_.analyzer.ReleaseSlabs_()
			this_.analyzer = nil
		}
	}()

	if this_.chunkGraph.GetChunkCount() > int32(1) {
		this_.analyzer.Process_2(externRoot, root)
		this_.priv_CrossChunkMethodMotion_moveMethods_2280(this_.analyzer.GetAllNameInfo())
	}
}

// Java method: CrossChunkMethodMotion.moveMethods(java.util.Collection<com.google.javascript.jscomp.AnalyzePrototypeProperties.NameInfo> allNameInfo)
func (this_ *CrossChunkMethodMotion) priv_CrossChunkMethodMotion_moveMethods_2280(allNameInfo j2gort.Ext_java_util_Collection) {
	_ = allNameInfo
	var hasStubDeclaration bool = this_.idGenerator.HasGeneratedAnyIds()
	_ = hasStubDeclaration
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(allNameInfo))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			nameInfo := j2gort.CastPtr[AnalyzePrototypeProperties_NameInfo](rawElem_2)
			_ = nameInfo
			if !(nameInfo.IsReferenced()) {
				continue
			}
			if nameInfo.ReadsClosureVariables() {
				continue
			}
			var deepestCommonChunkRef *JSChunk = nameInfo.GetDeepestCommonChunkRef()
			_ = deepestCommonChunkRef
			if deepestCommonChunkRef == nil {
				this_.compiler.Report_1(JSError_Make_2(StaticGet_CrossChunkMethodMotion_NULL_COMMON_CHUNK_ERROR(), []string{}))
				continue
			}
			var declarations j2gort.Ext_java_util_Iterator = nameInfo.GetDeclarations().DescendingIterator()
			_ = declarations
			for declarations.HasNext() {
				var symbol AnalyzePrototypeProperties_Symbol = j2gort.AsIface[AnalyzePrototypeProperties_Symbol](declarations.Next())
				_ = symbol
				var prototypeProperty AnalyzePrototypeProperties_PrototypeProperty
				_ = prototypeProperty
				var classMemberFunction *AnalyzePrototypeProperties_ClassMemberFunction
				_ = classMemberFunction
				var ok_3 bool
				prototypeProperty, ok_3 = j2gort.TryCast[AnalyzePrototypeProperties_PrototypeProperty](symbol)
				if ok_3 {
					this_.priv_CrossChunkMethodMotion_tryToMovePrototypeMethod_2281(nameInfo, deepestCommonChunkRef, prototypeProperty)
				} else {
					var ok_4 bool
					classMemberFunction, ok_4 = j2gort.TryCast[*AnalyzePrototypeProperties_ClassMemberFunction](symbol)
					if ok_4 {
						this_.priv_CrossChunkMethodMotion_tryToMoveMemberFunction_2287(nameInfo, deepestCommonChunkRef, classMemberFunction)
					}
				}
			}
		}
	}
	if (!(this_.noStubFunctions) && !(hasStubDeclaration)) && this_.idGenerator.HasGeneratedAnyIds() {
		var declarations_2 pkg_rhino.Node = this_.compiler.ParseSyntheticCode("chunk_method_stubbing", "var JSCompiler_stubMap = [];function JSCompiler_stubMethod(JSCompiler_stubMethod_id) {  return function() {    return JSCompiler_stubMap[JSCompiler_stubMethod_id].apply(        this, arguments);  };}function JSCompiler_unstubMethod(    JSCompiler_unstubMethod_id, JSCompiler_unstubMethod_body) {  return JSCompiler_stubMap[JSCompiler_unstubMethod_id] =       JSCompiler_unstubMethod_body;}")
		_ = declarations_2
		NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(declarations_2, this_.compiler)
		var firstScript pkg_rhino.Node = this_.compiler.GetNodeForCodeInsertion(nil)
		_ = firstScript
		firstScript.AddChildrenToFront(declarations_2.RemoveChildren())
		this_.compiler.ReportChangeToEnclosingScope(firstScript)
	}
}

// Java method: CrossChunkMethodMotion.tryToMovePrototypeMethod(com.google.javascript.jscomp.AnalyzePrototypeProperties.NameInfo nameInfo, com.google.javascript.jscomp.JSChunk deepestCommonChunkRef, com.google.javascript.jscomp.AnalyzePrototypeProperties.PrototypeProperty prop)
func (this_ *CrossChunkMethodMotion) priv_CrossChunkMethodMotion_tryToMovePrototypeMethod_2281(nameInfo *AnalyzePrototypeProperties_NameInfo, deepestCommonChunkRef *JSChunk, prop AnalyzePrototypeProperties_PrototypeProperty) {
	_ = nameInfo
	_ = deepestCommonChunkRef
	_ = prop
	if (prop.GetRootVar() == nil) || !(prop.GetRootVar().IsGlobal()) {
		return
	}
	if nameInfo.ReferencesSuper() {
		return
	}
	var value pkg_rhino.Node = prop.GetValue()
	_ = value
	var valueParent pkg_rhino.Node = value.GetParent()
	_ = valueParent
	if (!(value.IsFunction()) || valueParent.IsGetterDef()) || valueParent.IsSetterDef() {
		return
	}
	if this_.chunkGraph.DependsOn(deepestCommonChunkRef, prop.GetChunk()) {
		if CrossChunkMethodMotion_HasUnmovableRedeclaration(nameInfo, Upcast_AnalyzePrototypeProperties_Property(prop)) {
			return
		}
		var destParent pkg_rhino.Node = this_.compiler.GetNodeForCodeInsertion(deepestCommonChunkRef)
		_ = destParent
		if valueParent.IsMemberFunctionDef() {
			this_.priv_CrossChunkMethodMotion_movePrototypeObjectLiteralMethodShorthand_2284(nameInfo.name, destParent, value)
		} else if valueParent.IsStringKey() {
			this_.priv_CrossChunkMethodMotion_movePrototypeObjectLiteralProperty_2282(nameInfo.name, destParent, value)
		} else {
			if !(valueParent.IsAssign()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(valueParent))
			}
			this_.priv_CrossChunkMethodMotion_movePrototypeDotMethodAssignment_2283(destParent, value)
		}
	}
}

// Java method: CrossChunkMethodMotion.movePrototypeObjectLiteralProperty(java.lang.String propName, com.google.javascript.rhino.Node destParent, com.google.javascript.rhino.Node functionNode)
func (this_ *CrossChunkMethodMotion) priv_CrossChunkMethodMotion_movePrototypeObjectLiteralProperty_2282(propName string, destParent pkg_rhino.Node, functionNode pkg_rhino.Node) {
	_ = propName
	_ = destParent
	_ = functionNode
	if !(functionNode.IsFunction()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(functionNode))
	}
	var stringKey pkg_rhino.Node = functionNode.GetParent()
	_ = stringKey
	if !(stringKey.IsStringKey()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(stringKey))
	}
	var prototypeObjectLiteral pkg_rhino.Node = stringKey.GetParent()
	_ = prototypeObjectLiteral
	if !(prototypeObjectLiteral.IsObjectLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(prototypeObjectLiteral))
	}
	var assignNode pkg_rhino.Node = prototypeObjectLiteral.GetParent()
	_ = assignNode
	if !(assignNode.IsAssign() && prototypeObjectLiteral.IsSecondChildOf(assignNode)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(assignNode))
	}
	var ownerDotPrototypeNode pkg_rhino.Node = assignNode.GetFirstChild()
	_ = ownerDotPrototypeNode
	if !(ownerDotPrototypeNode.IsQualifiedName() && (j2gort.CheckNotNullString(ownerDotPrototypeNode.GetString()) == "prototype")) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(ownerDotPrototypeNode))
	}
	if this_.noStubFunctions {
		stringKey.Detach()
		this_.compiler.ReportChangeToEnclosingScope(prototypeObjectLiteral)
		var ownerDotPrototypeDotPropName pkg_rhino.Node = this_.astFactory.CreateGetProp(ownerDotPrototypeNode.CloneTree_0(), propName, AstFactory_S_Type_1_Node(functionNode))
		_ = ownerDotPrototypeDotPropName
		functionNode.Detach()
		var definitionStatement pkg_rhino.Node = this_.astFactory.CreateAssignStatement(ownerDotPrototypeDotPropName, functionNode).SrcrefTreeIfMissing(stringKey)
		_ = definitionStatement
		destParent.AddChildToFront(definitionStatement)
		this_.compiler.ReportChangeToEnclosingScope(destParent)
	} else {
		var stubId int32 = this_.idGenerator.NewId()
		_ = stubId
		var stubCall pkg_rhino.Node = this_.priv_CrossChunkMethodMotion_createStubCall_2285(functionNode, stubId)
		_ = stubCall
		functionNode.ReplaceWith(stubCall)
		this_.compiler.ReportChangeToEnclosingScope(prototypeObjectLiteral)
		var ownerDotPrototypeDotPropName_2 pkg_rhino.Node = this_.astFactory.CreateGetProp(ownerDotPrototypeNode.CloneTree_0(), propName, AstFactory_S_Type_1_Node(functionNode))
		_ = ownerDotPrototypeDotPropName_2
		var unstubCall pkg_rhino.Node = this_.priv_CrossChunkMethodMotion_createUnstubCall_2286(functionNode, stubId)
		_ = unstubCall
		var definitionStatement_2 pkg_rhino.Node = this_.astFactory.CreateAssignStatement(ownerDotPrototypeDotPropName_2, unstubCall).SrcrefTreeIfMissing(stringKey)
		_ = definitionStatement_2
		destParent.AddChildToFront(definitionStatement_2)
		this_.compiler.ReportChangeToEnclosingScope(destParent)
	}
}

// Java method: CrossChunkMethodMotion.movePrototypeDotMethodAssignment(com.google.javascript.rhino.Node destParent, com.google.javascript.rhino.Node functionNode)
func (this_ *CrossChunkMethodMotion) priv_CrossChunkMethodMotion_movePrototypeDotMethodAssignment_2283(destParent pkg_rhino.Node, functionNode pkg_rhino.Node) {
	_ = destParent
	_ = functionNode
	if !(functionNode.IsFunction()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(functionNode))
	}
	var assignNode pkg_rhino.Node = functionNode.GetParent()
	_ = assignNode
	if !(assignNode.IsAssign() && functionNode.IsSecondChildOf(assignNode)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(assignNode))
	}
	var definitionStatement pkg_rhino.Node = assignNode.GetParent()
	_ = definitionStatement
	if !(definitionStatement.IsExprResult()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(assignNode))
	}
	if this_.noStubFunctions {
		var assignStatementParent pkg_rhino.Node = definitionStatement.GetParent()
		_ = assignStatementParent
		definitionStatement.Detach()
		this_.compiler.ReportChangeToEnclosingScope(assignStatementParent)
		destParent.AddChildToFront(definitionStatement)
		this_.compiler.ReportChangeToEnclosingScope(destParent)
	} else {
		var stubId int32 = this_.idGenerator.NewId()
		_ = stubId
		var originalDefinitionPlaceholder pkg_rhino.Node = this_.astFactory.CreateEmpty()
		_ = originalDefinitionPlaceholder
		functionNode.ReplaceWith(originalDefinitionPlaceholder)
		var newDefinitionStatement pkg_rhino.Node = definitionStatement.CloneTree_0()
		_ = newDefinitionStatement
		var newDefinitionPlaceholder pkg_rhino.Node = newDefinitionStatement.GetOnlyChild().GetLastChild()
		_ = newDefinitionPlaceholder
		var stubCall pkg_rhino.Node = this_.priv_CrossChunkMethodMotion_createStubCall_2285(functionNode, stubId)
		_ = stubCall
		originalDefinitionPlaceholder.ReplaceWith(stubCall)
		this_.compiler.ReportChangeToEnclosingScope(definitionStatement)
		var unstubCall pkg_rhino.Node = this_.priv_CrossChunkMethodMotion_createUnstubCall_2286(functionNode, stubId)
		_ = unstubCall
		newDefinitionPlaceholder.ReplaceWith(unstubCall)
		destParent.AddChildToFront(newDefinitionStatement)
		this_.compiler.ReportChangeToEnclosingScope(destParent)
	}
}

// Java method: CrossChunkMethodMotion.movePrototypeObjectLiteralMethodShorthand(java.lang.String propName, com.google.javascript.rhino.Node destParent, com.google.javascript.rhino.Node functionNode)
func (this_ *CrossChunkMethodMotion) priv_CrossChunkMethodMotion_movePrototypeObjectLiteralMethodShorthand_2284(propName string, destParent pkg_rhino.Node, functionNode pkg_rhino.Node) {
	_ = propName
	_ = destParent
	_ = functionNode
	if !(functionNode.IsFunction()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(functionNode))
	}
	var memberFunctionDef pkg_rhino.Node = functionNode.GetParent()
	_ = memberFunctionDef
	if !(memberFunctionDef.IsMemberFunctionDef()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(memberFunctionDef))
	}
	var prototypeObjectLiteral pkg_rhino.Node = memberFunctionDef.GetParent()
	_ = prototypeObjectLiteral
	if !(prototypeObjectLiteral.IsObjectLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(prototypeObjectLiteral))
	}
	var assignNode pkg_rhino.Node = prototypeObjectLiteral.GetParent()
	_ = assignNode
	if !(assignNode.IsAssign() && prototypeObjectLiteral.IsSecondChildOf(assignNode)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(assignNode))
	}
	var ownerDotPrototypeNode pkg_rhino.Node = assignNode.GetFirstChild()
	_ = ownerDotPrototypeNode
	if !(ownerDotPrototypeNode.IsQualifiedName() && (j2gort.CheckNotNullString(ownerDotPrototypeNode.GetString()) == "prototype")) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(ownerDotPrototypeNode))
	}
	var originalScript pkg_rhino.Node = NodeUtil_GetEnclosingScript(memberFunctionDef)
	_ = originalScript
	if this_.noStubFunctions {
		memberFunctionDef.Detach()
		this_.compiler.ReportChangeToEnclosingScope(prototypeObjectLiteral)
		var ownerDotPrototypeDotPropName pkg_rhino.Node = this_.astFactory.CreateGetProp(ownerDotPrototypeNode.CloneTree_0(), propName, AstFactory_S_Type_1_Node(functionNode))
		_ = ownerDotPrototypeDotPropName
		var definitionStatement pkg_rhino.Node = this_.astFactory.CreateAssignStatement(ownerDotPrototypeDotPropName, functionNode.Detach()).SrcrefTreeIfMissing(memberFunctionDef)
		_ = definitionStatement
		destParent.AddChildToFront(definitionStatement)
		NodeUtil_AddFeaturesToScript(destParent, NodeUtil_GetFeatureSetOfScript(originalScript), this_.compiler)
		this_.compiler.ReportChangeToEnclosingScope(destParent)
	} else {
		var stubId int32 = this_.idGenerator.NewId()
		_ = stubId
		var stubCall pkg_rhino.Node = this_.priv_CrossChunkMethodMotion_createStubCall_2285(functionNode, stubId)
		_ = stubCall
		memberFunctionDef.ReplaceWith(this_.astFactory.CreateStringKey(propName, stubCall))
		this_.compiler.ReportChangeToEnclosingScope(prototypeObjectLiteral)
		var ownerDotPrototypeDotPropName_2 pkg_rhino.Node = this_.astFactory.CreateGetProp(ownerDotPrototypeNode.CloneTree_0(), propName, AstFactory_S_Type_1_Node(functionNode))
		_ = ownerDotPrototypeDotPropName_2
		var unstubCall pkg_rhino.Node = this_.priv_CrossChunkMethodMotion_createUnstubCall_2286(functionNode.Detach(), stubId)
		_ = unstubCall
		var definitionStatement_2 pkg_rhino.Node = this_.astFactory.CreateAssignStatement(ownerDotPrototypeDotPropName_2, unstubCall).SrcrefTreeIfMissing(memberFunctionDef)
		_ = definitionStatement_2
		destParent.AddChildToFront(definitionStatement_2)
		NodeUtil_AddFeaturesToScript(destParent, NodeUtil_GetFeatureSetOfScript(originalScript), this_.compiler)
		this_.compiler.ReportChangeToEnclosingScope(destParent)
	}
}

// Java method: CrossChunkMethodMotion.createStubCall(com.google.javascript.rhino.Node originalDefinition, int stubId)
func (this_ *CrossChunkMethodMotion) priv_CrossChunkMethodMotion_createStubCall_2285(originalDefinition pkg_rhino.Node, stubId int32) pkg_rhino.Node {
	_ = originalDefinition
	_ = stubId
	return this_.astFactory.CreateCall(this_.astFactory.CreateNameWithUnknownType("JSCompiler_stubMethod"), AstFactory_S_Type_1_Node(originalDefinition), []pkg_rhino.Node{this_.astFactory.CreateNumber(float64(stubId))}).SrcrefTreeIfMissing(originalDefinition)
}

// Java method: CrossChunkMethodMotion.createUnstubCall(com.google.javascript.rhino.Node functionNode, int stubId)
func (this_ *CrossChunkMethodMotion) priv_CrossChunkMethodMotion_createUnstubCall_2286(functionNode pkg_rhino.Node, stubId int32) pkg_rhino.Node {
	_ = functionNode
	_ = stubId
	return this_.astFactory.CreateCall(this_.astFactory.CreateNameWithUnknownType("JSCompiler_unstubMethod"), AstFactory_S_Type_1_Node(functionNode), []pkg_rhino.Node{this_.astFactory.CreateNumber(float64(stubId)), functionNode}).SrcrefTreeIfMissing(functionNode)
}

// Java method: CrossChunkMethodMotion.tryToMoveMemberFunction(com.google.javascript.jscomp.AnalyzePrototypeProperties.NameInfo nameInfo, com.google.javascript.jscomp.JSChunk deepestCommonChunkRef, com.google.javascript.jscomp.AnalyzePrototypeProperties.ClassMemberFunction classMemberFunction)
func (this_ *CrossChunkMethodMotion) priv_CrossChunkMethodMotion_tryToMoveMemberFunction_2287(nameInfo *AnalyzePrototypeProperties_NameInfo, deepestCommonChunkRef *JSChunk, classMemberFunction *AnalyzePrototypeProperties_ClassMemberFunction) {
	_ = nameInfo
	_ = deepestCommonChunkRef
	_ = classMemberFunction
	var rootVar *Var = classMemberFunction.GetRootVar()
	_ = rootVar
	if (rootVar == nil) || !(rootVar.IsGlobal()) {
		return
	}
	var definitionNode pkg_rhino.Node = classMemberFunction.GetDefinitionNode()
	_ = definitionNode
	if !(definitionNode.IsMemberFunctionDef()) {
		return
	}
	if NodeUtil_IsEs6ConstructorMemberFunctionDef(definitionNode) {
		return
	}
	if nameInfo.ReferencesSuper() {
		return
	}
	if this_.chunkGraph.DependsOn(deepestCommonChunkRef, classMemberFunction.GetChunk()) {
		if CrossChunkMethodMotion_HasUnmovableRedeclaration(nameInfo, ToIface_AnalyzePrototypeProperties_Property(classMemberFunction)) {
			return
		}
		var destinationParent pkg_rhino.Node = this_.compiler.GetNodeForCodeInsertion(deepestCommonChunkRef)
		_ = destinationParent
		var classNameNode pkg_rhino.Node = rootVar.GetNameNode()
		_ = classNameNode
		if this_.noStubFunctions {
			this_.priv_CrossChunkMethodMotion_moveClassInstanceMethodWithoutStub_2288(classNameNode, definitionNode, destinationParent)
		} else {
			this_.priv_CrossChunkMethodMotion_moveClassInstanceMethodWithStub_2289(classNameNode, definitionNode, destinationParent)
		}
	}
}

// Java method: CrossChunkMethodMotion.moveClassInstanceMethodWithoutStub(com.google.javascript.rhino.Node classNameNode, com.google.javascript.rhino.Node methodDefinition, com.google.javascript.rhino.Node destinationParent)
func (this_ *CrossChunkMethodMotion) priv_CrossChunkMethodMotion_moveClassInstanceMethodWithoutStub_2288(classNameNode pkg_rhino.Node, methodDefinition pkg_rhino.Node, destinationParent pkg_rhino.Node) {
	_ = classNameNode
	_ = methodDefinition
	_ = destinationParent
	if !(methodDefinition.IsMemberFunctionDef()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(methodDefinition))
	}
	var classMembers pkg_rhino.Node = j2gort.CheckNotNullNode(methodDefinition.GetParent())
	_ = classMembers
	if !(classMembers.IsClassMembers()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(classMembers))
	}
	var classNode pkg_rhino.Node = classMembers.GetParent()
	_ = classNode
	if !(classNode.IsClass()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(classNode))
	}
	var originalScript pkg_rhino.Node = NodeUtil_GetEnclosingScript(methodDefinition)
	_ = originalScript
	methodDefinition.Detach()
	this_.compiler.ReportChangeToEnclosingScope(classMembers)
	var functionNode pkg_rhino.Node = j2gort.CheckNotNullNode(methodDefinition.GetOnlyChild())
	_ = functionNode
	var clonedNameNode pkg_rhino.Node = classNameNode.CloneNode_0().CopyTypeFrom(classNode)
	_ = clonedNameNode
	var classNameDotPrototypeDotPropName pkg_rhino.Node = this_.astFactory.CreateGetProp(this_.astFactory.CreatePrototypeAccess(clonedNameNode), methodDefinition.GetString(), AstFactory_S_Type_1_Node(functionNode))
	_ = classNameDotPrototypeDotPropName
	functionNode.Detach()
	var definitionStatementNode pkg_rhino.Node = this_.astFactory.CreateAssignStatement(classNameDotPrototypeDotPropName, functionNode).SrcrefTreeIfMissing(methodDefinition)
	_ = definitionStatementNode
	destinationParent.AddChildToFront(definitionStatementNode)
	NodeUtil_AddFeaturesToScript(destinationParent, NodeUtil_GetFeatureSetOfScript(originalScript), this_.compiler)
	this_.compiler.ReportChangeToEnclosingScope(destinationParent)
}

// Java method: CrossChunkMethodMotion.moveClassInstanceMethodWithStub(com.google.javascript.rhino.Node classNameNode, com.google.javascript.rhino.Node methodDefinition, com.google.javascript.rhino.Node destinationParent)
func (this_ *CrossChunkMethodMotion) priv_CrossChunkMethodMotion_moveClassInstanceMethodWithStub_2289(classNameNode pkg_rhino.Node, methodDefinition pkg_rhino.Node, destinationParent pkg_rhino.Node) {
	_ = classNameNode
	_ = methodDefinition
	_ = destinationParent
	if !(methodDefinition.IsMemberFunctionDef()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(methodDefinition))
	}
	var classMembers pkg_rhino.Node = j2gort.CheckNotNullNode(methodDefinition.GetParent())
	_ = classMembers
	if !(classMembers.IsClassMembers()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(classMembers))
	}
	var classNode pkg_rhino.Node = classMembers.GetParent()
	_ = classNode
	if !(classNode.IsClass()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(classNode))
	}
	var stubId int32 = this_.idGenerator.NewId()
	_ = stubId
	var clonedNameNode pkg_rhino.Node = classNameNode.CloneNode_0().CopyTypeFrom(classNode)
	_ = clonedNameNode
	var classNameDotPrototypeDotPropName pkg_rhino.Node = this_.astFactory.CreateGetProp(this_.astFactory.CreatePrototypeAccess(clonedNameNode), methodDefinition.GetString(), AstFactory_S_Type_1_Node(methodDefinition))
	_ = classNameDotPrototypeDotPropName
	var stubCall pkg_rhino.Node = this_.priv_CrossChunkMethodMotion_createStubCall_2285(methodDefinition, stubId)
	_ = stubCall
	var stubDefinitionStatement pkg_rhino.Node = this_.astFactory.CreateAssignStatement(classNameDotPrototypeDotPropName, stubCall).SrcrefTreeIfMissing(methodDefinition)
	_ = stubDefinitionStatement
	var classDefiningStatement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(classMembers)
	_ = classDefiningStatement
	stubDefinitionStatement.InsertAfter(classDefiningStatement)
	var originalScript pkg_rhino.Node = NodeUtil_GetEnclosingScript(methodDefinition)
	_ = originalScript
	methodDefinition.Detach()
	this_.compiler.ReportChangeToEnclosingScope(classMembers)
	var classNameDotPrototypeDotPropName2 pkg_rhino.Node = classNameDotPrototypeDotPropName.CloneTree_0()
	_ = classNameDotPrototypeDotPropName2
	var functionNode pkg_rhino.Node = j2gort.CheckNotNullNode(methodDefinition.GetOnlyChild())
	_ = functionNode
	functionNode.Detach()
	var unstubCall pkg_rhino.Node = this_.priv_CrossChunkMethodMotion_createUnstubCall_2286(functionNode, stubId)
	_ = unstubCall
	var statementNode pkg_rhino.Node = this_.astFactory.CreateAssignStatement(classNameDotPrototypeDotPropName2, unstubCall).SrcrefTreeIfMissing(methodDefinition)
	_ = statementNode
	destinationParent.AddChildToFront(statementNode)
	NodeUtil_AddFeaturesToScript(destinationParent, NodeUtil_GetFeatureSetOfScript(originalScript), this_.compiler)
	this_.compiler.ReportChangeToEnclosingScope(destinationParent)
}

// Java static method: CrossChunkMethodMotion.hasUnmovableRedeclaration(com.google.javascript.jscomp.AnalyzePrototypeProperties.NameInfo nameInfo, com.google.javascript.jscomp.AnalyzePrototypeProperties.Property prop)
func CrossChunkMethodMotion_HasUnmovableRedeclaration(nameInfo *AnalyzePrototypeProperties_NameInfo, prop AnalyzePrototypeProperties_Property) bool {
	clinit_CrossChunkMethodMotion()
	_ = nameInfo
	_ = prop
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(nameInfo.GetDeclarations()))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			symbol := j2gort.AsIface[AnalyzePrototypeProperties_Symbol](rawElem_6)
			_ = symbol
			var otherProp AnalyzePrototypeProperties_Property
			_ = otherProp
			var ok_7 bool
			otherProp, ok_7 = j2gort.TryCast[AnalyzePrototypeProperties_Property](symbol)
			if ok_7 {
				if (!(prop == otherProp || ((prop == nil || prop.IsJavaNil_()) && (otherProp == nil || otherProp.IsJavaNil_()))) && (prop.GetRootVar() == otherProp.GetRootVar())) && (prop.GetChunk() != otherProp.GetChunk()) {
					return true
				}
			}
		}
	}
	return false
}
