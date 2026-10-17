// Port of com/google/javascript/jscomp/CrossChunkMethodMotion.java @ closure-compiler v20260415

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CrossChunkMethodMotion moves prototype methods into later chunks.
type CrossChunkMethodMotion struct {
	compiler        AbstractCompiler
	idGenerator     *IdGenerator
	analyzer        *AnalyzePrototypeProperties
	chunkGraph      *JSChunkGraph
	noStubFunctions bool
	astFactory      *AstFactory
}

// Internal errors
var CrossChunkMethodMotion_NULL_COMMON_CHUNK_ERROR = DiagnosticType_Error("JSC_INTERNAL_ERROR_CHUNK_DEPEND", "null deepest common chunk")

const CrossChunkMethodMotion_STUB_METHOD_NAME = "JSCompiler_stubMethod"
const CrossChunkMethodMotion_UNSTUB_METHOD_NAME = "JSCompiler_unstubMethod"

const CrossChunkMethodMotion_STUB_DECLARATIONS = "var JSCompiler_stubMap = [];" +
	"function JSCompiler_stubMethod(JSCompiler_stubMethod_id) {" +
	"  return function() {" +
	"    return JSCompiler_stubMap[JSCompiler_stubMethod_id].apply(" +
	"        this, arguments);" +
	"  };" +
	"}" +
	"function JSCompiler_unstubMethod(" +
	"    JSCompiler_unstubMethod_id, JSCompiler_unstubMethod_body) {" +
	"  return JSCompiler_stubMap[JSCompiler_unstubMethod_id] = " +
	"      JSCompiler_unstubMethod_body;" +
	"}"

// NewCrossChunkMethodMotion creates a new pass for moving prototype
// properties.
//
// compiler is the compiler; idGenerator an id generator for method stubs;
// canModifyExterns: if true, then we can move prototype properties that are
// declared in the externs file; noStubFunctions: if true, we can move
// methods without stub functions in the parent chunk.
//
// Java: CrossChunkMethodMotion(AbstractCompiler, IdGenerator, boolean, boolean)
func NewCrossChunkMethodMotion(
	compiler AbstractCompiler,
	idGenerator *IdGenerator,
	canModifyExterns bool,
	noStubFunctions bool) *CrossChunkMethodMotion {
	m := &CrossChunkMethodMotion{}
	m.compiler = compiler
	m.idGenerator = idGenerator
	m.chunkGraph = compiler.GetChunkGraph()
	m.analyzer =
		NewAnalyzePrototypeProperties(
			compiler, m.chunkGraph, canModifyExterns /* anchorUnusedVars= */, false, noStubFunctions)
	m.noStubFunctions = noStubFunctions
	m.astFactory = compiler.CreateAstFactory()
	return m
}

// Java: CrossChunkMethodMotion.process(Node, Node)
func (m *CrossChunkMethodMotion) Process(externRoot *rhino.Node, root *rhino.Node) {
	// If there are < 2 chunks, then we will never move anything,
	// so we're done.
	if m.chunkGraph.GetChunkCount() > 1 {
		m.analyzer.Process(externRoot, root)
		m.MoveMethods(m.analyzer.GetAllNameInfo())
	}
}

// MoveMethods moves methods deeper in the chunk graph when possible.
//
// Java: CrossChunkMethodMotion.moveMethods(Collection<NameInfo>)
func (m *CrossChunkMethodMotion) MoveMethods(allNameInfo []*AnalyzePrototypeProperties_NameInfo) {
	hasStubDeclaration := m.idGenerator.HasGeneratedAnyIds()
	for _, nameInfo := range allNameInfo {
		if !nameInfo.IsReferenced() {
			// The code below can't do anything with unreferenced name
			// infos.  They should be skipped to avoid NPE since their
			// deepestCommonChunkRef is null.
			continue
		}

		if nameInfo.ReadsClosureVariables() {
			continue
		}

		deepestCommonChunkRef := nameInfo.GetDeepestCommonChunkRef()
		if deepestCommonChunkRef == nil {
			m.compiler.Report(JSError_Make(CrossChunkMethodMotion_NULL_COMMON_CHUNK_ERROR))
			continue
		}

		// nameInfo.getDeclarations().descendingIterator()
		declarations := nameInfo.GetDeclarations().Slice()
		for i := len(declarations) - 1; i >= 0; i-- {
			symbol := declarations[i]
			if prototypeProperty, ok := symbol.(AnalyzePrototypeProperties_PrototypeProperty); ok {
				m.TryToMovePrototypeMethod(nameInfo, deepestCommonChunkRef, prototypeProperty)
			} else if classMemberFunction, ok := symbol.(*AnalyzePrototypeProperties_ClassMemberFunction); ok {
				m.TryToMoveMemberFunction(nameInfo, deepestCommonChunkRef, classMemberFunction)
			} // else it's a variable definition, and we don't move those.
		}
	}

	if !m.noStubFunctions && !hasStubDeclaration && m.idGenerator.HasGeneratedAnyIds() {
		// Declare stub functions in the top-most chunk.
		declarations := m.compiler.ParseSyntheticCode("chunk_method_stubbing", CrossChunkMethodMotion_STUB_DECLARATIONS)
		NodeUtil_MarkNewScopesChanged(declarations, m.compiler)
		firstScript := m.compiler.GetNodeForCodeInsertion(nil)
		firstScript.AddChildrenToFront(declarations.RemoveChildren())
		m.compiler.ReportChangeToEnclosingScope(firstScript)
	}
}

// Java: CrossChunkMethodMotion.tryToMovePrototypeMethod(NameInfo, JSChunk, PrototypeProperty)
func (m *CrossChunkMethodMotion) TryToMovePrototypeMethod(
	nameInfo *AnalyzePrototypeProperties_NameInfo, deepestCommonChunkRef *JSChunk, prop AnalyzePrototypeProperties_PrototypeProperty) {

	// We should only move a property across chunks if:
	// 1) We can move it deeper in the chunk graph, and
	// 2) it's a function, and
	// 3) it is not a GETTER_DEF or a SETTER_DEF, and
	// 4) it does not refer to `super`
	// 5) the class is available in the global scope.
	//
	// #1 should be obvious. #2 is more subtle. It's possible
	// to copy off of a prototype, as in the code:
	// for (var k in Foo.prototype) {
	//   doSomethingWith(Foo.prototype[k]);
	// }
	// This is a common way to implement pseudo-multiple inheritance in JS.
	//
	// So if we move a prototype method into a deeper chunk, we must
	// replace it with a stub function so that it preserves its original
	// behavior.
	if prop.GetRootVar() == nil || !prop.GetRootVar().IsGlobal() {
		return
	}

	if nameInfo.ReferencesSuper() {
		// It is illegal to move `super` outside of a member function def.
		return
	}

	value := prop.GetValue()
	valueParent := value.GetParent()
	// Only attempt to move normal functions.
	if !value.IsFunction() ||
		// A GET or SET can't be deferred like a normal
		// FUNCTION property definition as a mix-in would get the result
		// of a GET instead of the function itself.
		valueParent.IsGetterDef() ||
		valueParent.IsSetterDef() {
		return
	}

	if m.chunkGraph.DependsOn(deepestCommonChunkRef, prop.GetChunk()) {
		if CrossChunkMethodMotion_HasUnmovableRedeclaration(nameInfo, prop) {
			// If it has been redeclared on the same object, skip it.
			return
		}
		destParent := m.compiler.GetNodeForCodeInsertion(deepestCommonChunkRef)
		if valueParent.IsMemberFunctionDef() {
			m.MovePrototypeObjectLiteralMethodShorthand(nameInfo.name, destParent, value)
		} else if valueParent.IsStringKey() {
			m.MovePrototypeObjectLiteralProperty(nameInfo.name, destParent, value)
		} else {
			// Note that computed properties should have been filtered out by
			// AnalyzePrototypeProperties, because they don't have a recognizable property name.
			// Getters and setters are filtered out by the code above
			jutil.CheckStateMsg1(valueParent.IsAssign(), "%s", valueParent)
			m.MovePrototypeDotMethodAssignment(destParent, value)
		}
	}
}

// MovePrototypeObjectLiteralProperty moves a property defined by object
// literal assigned to `.prototype`.
//
//	Foo.prototype = { propName: function() {}};
//
// Java: CrossChunkMethodMotion.movePrototypeObjectLiteralProperty(String, Node, Node)
func (m *CrossChunkMethodMotion) MovePrototypeObjectLiteralProperty(
	propName string, destParent *rhino.Node, functionNode *rhino.Node) {
	jutil.CheckStateMsg1(functionNode.IsFunction(), "%s", functionNode)
	stringKey := functionNode.GetParent()
	jutil.CheckStateMsg1(stringKey.IsStringKey(), "%s", stringKey)
	prototypeObjectLiteral := stringKey.GetParent()
	jutil.CheckStateMsg1(prototypeObjectLiteral.IsObjectLit(), "%s", prototypeObjectLiteral)
	assignNode := prototypeObjectLiteral.GetParent()
	jutil.CheckStateMsg1(
		assignNode.IsAssign() && prototypeObjectLiteral.IsSecondChildOf(assignNode), "%s", assignNode)
	ownerDotPrototypeNode := assignNode.GetFirstChild()
	jutil.CheckStateMsg1(
		ownerDotPrototypeNode.IsQualifiedName() &&
			ownerDotPrototypeNode.GetString() == "prototype",
		"%s", ownerDotPrototypeNode)

	if m.noStubFunctions {
		// Remove the definition from the object literal
		stringKey.Detach()
		m.compiler.ReportChangeToEnclosingScope(prototypeObjectLiteral)

		// Prepend definition to new chunk
		// Foo.prototype.propName = function() {};
		ownerDotPrototypeDotPropName :=
			m.astFactory.CreateGetProp(ownerDotPrototypeNode.CloneTree(), propName, AstFactory_Type_(functionNode))
		functionNode.Detach()
		definitionStatement :=
			m.astFactory.
				CreateAssignStatement(ownerDotPrototypeDotPropName, functionNode).
				SrcrefTreeIfMissing(stringKey)
		destParent.AddChildToFront(definitionStatement)
		m.compiler.ReportChangeToEnclosingScope(destParent)
	} else {
		stubId := m.idGenerator.NewId()
		// { propName: function() {} } => { propName: JSCompiler_stubMethod(0) }
		stubCall := m.CreateStubCall(functionNode, stubId)
		functionNode.ReplaceWith(stubCall)
		m.compiler.ReportChangeToEnclosingScope(prototypeObjectLiteral)

		// Prepend definition to new chunk
		// Foo.prototype.propName = function() {};
		ownerDotPrototypeDotPropName :=
			m.astFactory.CreateGetProp(ownerDotPrototypeNode.CloneTree(), propName, AstFactory_Type_(functionNode))
		unstubCall := m.CreateUnstubCall(functionNode, stubId)
		definitionStatement :=
			m.astFactory.
				CreateAssignStatement(ownerDotPrototypeDotPropName, unstubCall).
				SrcrefTreeIfMissing(stringKey)
		destParent.AddChildToFront(definitionStatement)
		m.compiler.ReportChangeToEnclosingScope(destParent)
	}
}

// MovePrototypeDotMethodAssignment moves a property defined by assignment to
// `.prototype` or `.prototype.propName`.
//
//	Foo.prototype.propName = function() {};
//
// Java: CrossChunkMethodMotion.movePrototypeDotMethodAssignment(Node, Node)
func (m *CrossChunkMethodMotion) MovePrototypeDotMethodAssignment(destParent *rhino.Node, functionNode *rhino.Node) {
	jutil.CheckStateMsg1(functionNode.IsFunction(), "%s", functionNode)
	assignNode := functionNode.GetParent()
	jutil.CheckStateMsg1(assignNode.IsAssign() && functionNode.IsSecondChildOf(assignNode), "%s", assignNode)
	definitionStatement := assignNode.GetParent()
	jutil.CheckStateMsg1(definitionStatement.IsExprResult(), "%s", assignNode)

	if m.noStubFunctions {
		// Remove the definition statement from its current location
		assignStatementParent := definitionStatement.GetParent()
		definitionStatement.Detach()
		m.compiler.ReportChangeToEnclosingScope(assignStatementParent)

		// Prepend definition to new chunk
		// Foo.prototype.propName = function() {};
		destParent.AddChildToFront(definitionStatement)
		m.compiler.ReportChangeToEnclosingScope(destParent)
	} else {
		stubId := m.idGenerator.NewId()

		// replace function definition with temporary placeholder so we can clone the whole
		// assignment statement without cloning the function definition itself.
		originalDefinitionPlaceholder := m.astFactory.CreateEmpty()
		functionNode.ReplaceWith(originalDefinitionPlaceholder)
		newDefinitionStatement := definitionStatement.CloneTree()
		newDefinitionPlaceholder :=
			newDefinitionStatement. // EXPR_RESULT
						GetOnlyChild(). // ASSIGN
						GetLastChild()  // EMPTY RHS node

		// convert original assignment statement to
		// owner.prototype.propName = JSCompiler_stubMethod(0);
		stubCall := m.CreateStubCall(functionNode, stubId)
		originalDefinitionPlaceholder.ReplaceWith(stubCall)
		m.compiler.ReportChangeToEnclosingScope(definitionStatement)

		// Prepend new definition to new chunk
		// Foo.prototype.propName = JSCompiler_unstubMethod(0, function() {});
		unstubCall := m.CreateUnstubCall(functionNode, stubId)
		newDefinitionPlaceholder.ReplaceWith(unstubCall)
		destParent.AddChildToFront(newDefinitionStatement)
		m.compiler.ReportChangeToEnclosingScope(destParent)
	}
}

// MovePrototypeObjectLiteralMethodShorthand moves a property defined by
// assignment to `.prototype` or `.prototype.propName`.
//
//	Foo.prototype = { propName() {}};
//
// Java: CrossChunkMethodMotion.movePrototypeObjectLiteralMethodShorthand(String, Node, Node)
func (m *CrossChunkMethodMotion) MovePrototypeObjectLiteralMethodShorthand(
	propName string, destParent *rhino.Node, functionNode *rhino.Node) {
	jutil.CheckStateMsg1(functionNode.IsFunction(), "%s", functionNode)
	memberFunctionDef := functionNode.GetParent()
	jutil.CheckStateMsg1(memberFunctionDef.IsMemberFunctionDef(), "%s", memberFunctionDef)
	prototypeObjectLiteral := memberFunctionDef.GetParent()
	jutil.CheckStateMsg1(prototypeObjectLiteral.IsObjectLit(), "%s", prototypeObjectLiteral)
	assignNode := prototypeObjectLiteral.GetParent()
	jutil.CheckStateMsg1(
		assignNode.IsAssign() && prototypeObjectLiteral.IsSecondChildOf(assignNode), "%s", assignNode)
	ownerDotPrototypeNode := assignNode.GetFirstChild()
	jutil.CheckStateMsg1(
		ownerDotPrototypeNode.IsQualifiedName() &&
			ownerDotPrototypeNode.GetString() == "prototype",
		"%s", ownerDotPrototypeNode)
	originalScript := NodeUtil_GetEnclosingScript(memberFunctionDef)

	if m.noStubFunctions {
		// Remove the definition from the object literal
		memberFunctionDef.Detach()
		m.compiler.ReportChangeToEnclosingScope(prototypeObjectLiteral)

		// Prepend definition to new chunk
		// Foo.prototype.propName = function() {};
		ownerDotPrototypeDotPropName :=
			m.astFactory.CreateGetProp(ownerDotPrototypeNode.CloneTree(), propName, AstFactory_Type_(functionNode))
		definitionStatement :=
			m.astFactory.
				CreateAssignStatement(ownerDotPrototypeDotPropName, functionNode.Detach()).
				SrcrefTreeIfMissing(memberFunctionDef)
		destParent.AddChildToFront(definitionStatement)
		NodeUtil_AddFeaturesToScript(
			destParent, NodeUtil_GetFeatureSetOfScript(originalScript), m.compiler)
		m.compiler.ReportChangeToEnclosingScope(destParent)
	} else {
		stubId := m.idGenerator.NewId()
		// { propName() {} } => { propName: JSCompiler_stubMethod(0) }
		stubCall := m.CreateStubCall(functionNode, stubId)
		memberFunctionDef.ReplaceWith(m.astFactory.CreateStringKey(propName, stubCall))
		m.compiler.ReportChangeToEnclosingScope(prototypeObjectLiteral)

		// Prepend definition to new chunk
		// Foo.prototype.propName = function() {};
		ownerDotPrototypeDotPropName :=
			m.astFactory.CreateGetProp(ownerDotPrototypeNode.CloneTree(), propName, AstFactory_Type_(functionNode))
		unstubCall := m.CreateUnstubCall(functionNode.Detach(), stubId)
		definitionStatement :=
			m.astFactory.
				CreateAssignStatement(ownerDotPrototypeDotPropName, unstubCall).
				SrcrefTreeIfMissing(memberFunctionDef)
		destParent.AddChildToFront(definitionStatement)
		NodeUtil_AddFeaturesToScript(
			destParent, NodeUtil_GetFeatureSetOfScript(originalScript), m.compiler)

		m.compiler.ReportChangeToEnclosingScope(destParent)
	}
}

// CreateStubCall returns a new Node to be used as the stub definition for a
// method.
//
// originalDefinition is the function Node whose definition is being
// stubbed; stubId the ID to use for stubbing and unstubbing. Returns a Node
// that looks like JSCompiler_stubMethod(0).
//
// Java: CrossChunkMethodMotion.createStubCall(Node, int)
func (m *CrossChunkMethodMotion) CreateStubCall(originalDefinition *rhino.Node, stubId int) *rhino.Node {
	return m.astFactory.
		CreateCall(
			// We can't look up the type of the stub creating method, because we add its
			// definition after type checking.
			m.astFactory.CreateNameWithUnknownType(CrossChunkMethodMotion_STUB_METHOD_NAME),
			AstFactory_Type_(originalDefinition),
			m.astFactory.CreateNumber(float64(stubId))).
		SrcrefTreeIfMissing(originalDefinition)
}

// CreateUnstubCall returns a new Node to be used as the stub definition for
// a method.
//
// functionNode is the actual function definition to be attached. Must be
// detached now. stubId is the ID to use for stubbing and unstubbing. Returns
// a Node that looks like JSCompiler_unstubMethod(0, function() {}).
//
// Java: CrossChunkMethodMotion.createUnstubCall(Node, int)
func (m *CrossChunkMethodMotion) CreateUnstubCall(functionNode *rhino.Node, stubId int) *rhino.Node {
	return m.astFactory.
		CreateCall(
			// We can't look up the type of the stub creating method, because we add its
			// definition after type checking.
			m.astFactory.CreateNameWithUnknownType(CrossChunkMethodMotion_UNSTUB_METHOD_NAME),
			AstFactory_Type_(functionNode),
			m.astFactory.CreateNumber(float64(stubId)),
			functionNode).
		SrcrefTreeIfMissing(functionNode)
}

// TryToMoveMemberFunction, if possible, moves a class instance member
// function definition to the deepest chunk common to all uses of the
// method.
//
// nameInfo is information about all definitions of the given property name;
// deepestCommonChunkRef: all uses of the method are either in this chunk or
// in chunks that depend on it; classMemberFunction the definition of the
// method within its class body.
//
// Java: CrossChunkMethodMotion.tryToMoveMemberFunction(NameInfo, JSChunk, ClassMemberFunction)
func (m *CrossChunkMethodMotion) TryToMoveMemberFunction(
	nameInfo *AnalyzePrototypeProperties_NameInfo, deepestCommonChunkRef *JSChunk, classMemberFunction *AnalyzePrototypeProperties_ClassMemberFunction) {

	// We should only move a property across chunks if:
	// 1) We can move it deeper in the chunk graph,
	// 2) and it's a normal member function, and not a GETTER_DEF or a SETTER_DEF, or
	//    or the class constructor.
	// 3) and it does not refer to `super`, which would be invalid outside of the class.
	// 4) and the class is available in the global scope.
	rootVar := classMemberFunction.GetRootVar()
	if rootVar == nil || !rootVar.IsGlobal() {
		return
	}

	definitionNode := classMemberFunction.GetDefinitionNode()

	// Only attempt to move normal member functions.
	// A getter or setter cannot be as easily defined outside of the class to which it belongs.
	if !definitionNode.IsMemberFunctionDef() {
		return
	}

	if NodeUtil_IsEs6ConstructorMemberFunctionDef(definitionNode) {
		// Constructor functions cannot be moved.
		return
	}

	if nameInfo.ReferencesSuper() {
		// Do not move methods containing `super`, because it doesn't work outside of a
		// class method or object literal method.
		return
	}

	if m.chunkGraph.DependsOn(deepestCommonChunkRef, classMemberFunction.GetChunk()) {
		if CrossChunkMethodMotion_HasUnmovableRedeclaration(nameInfo, classMemberFunction) {
			// If it has been redeclared on the same object, skip it.
			return
		}

		destinationParent := m.compiler.GetNodeForCodeInsertion(deepestCommonChunkRef)
		classNameNode := rootVar.GetNameNode()
		if m.noStubFunctions {
			m.MoveClassInstanceMethodWithoutStub(classNameNode, definitionNode, destinationParent)
		} else {
			m.MoveClassInstanceMethodWithStub(classNameNode, definitionNode, destinationParent)
		}
	}
}

// Java: CrossChunkMethodMotion.moveClassInstanceMethodWithoutStub(Node, Node, Node)
func (m *CrossChunkMethodMotion) MoveClassInstanceMethodWithoutStub(
	classNameNode *rhino.Node, methodDefinition *rhino.Node, destinationParent *rhino.Node) {
	jutil.CheckArgumentMsg1(methodDefinition.IsMemberFunctionDef(), "%s", methodDefinition)
	classMembers := jutil.CheckNotNull(methodDefinition.GetParent())
	jutil.CheckStateMsg1(classMembers.IsClassMembers(), "%s", classMembers)
	classNode := classMembers.GetParent()
	jutil.CheckStateMsg1(classNode.IsClass(), "%s", classNode)
	originalScript := NodeUtil_GetEnclosingScript(methodDefinition)

	methodDefinition.Detach()
	m.compiler.ReportChangeToEnclosingScope(classMembers)

	// ClassName.prototype.propertyName = function() {};
	functionNode := jutil.CheckNotNull(methodDefinition.GetOnlyChild())
	clonedNameNode := classNameNode.CloneNode().CopyTypeFrom(classNode)
	classNameDotPrototypeDotPropName :=
		m.astFactory.CreateGetProp(
			m.astFactory.CreatePrototypeAccess(clonedNameNode),
			methodDefinition.GetString(),
			AstFactory_Type_(functionNode))
	functionNode.Detach()
	definitionStatementNode :=
		m.astFactory.
			CreateAssignStatement(classNameDotPrototypeDotPropName, functionNode).
			SrcrefTreeIfMissing(methodDefinition)
	destinationParent.AddChildToFront(definitionStatementNode)
	NodeUtil_AddFeaturesToScript(
		destinationParent, NodeUtil_GetFeatureSetOfScript(originalScript), m.compiler)
	m.compiler.ReportChangeToEnclosingScope(destinationParent)
}

// Java: CrossChunkMethodMotion.moveClassInstanceMethodWithStub(Node, Node, Node)
func (m *CrossChunkMethodMotion) MoveClassInstanceMethodWithStub(
	classNameNode *rhino.Node, methodDefinition *rhino.Node, destinationParent *rhino.Node) {
	jutil.CheckArgumentMsg1(methodDefinition.IsMemberFunctionDef(), "%s", methodDefinition)
	classMembers := jutil.CheckNotNull(methodDefinition.GetParent())
	jutil.CheckStateMsg1(classMembers.IsClassMembers(), "%s", classMembers)
	classNode := classMembers.GetParent()
	jutil.CheckStateMsg1(classNode.IsClass(), "%s", classNode)

	stubId := m.idGenerator.NewId()

	// Put a stub definition after the class
	// ClassName.prototype.propertyName = JSCompiler_stubMethod(id);
	clonedNameNode := classNameNode.CloneNode().CopyTypeFrom(classNode)
	classNameDotPrototypeDotPropName :=
		m.astFactory.CreateGetProp(
			m.astFactory.CreatePrototypeAccess(clonedNameNode),
			methodDefinition.GetString(),
			AstFactory_Type_(methodDefinition))
	stubCall := m.CreateStubCall(methodDefinition, stubId)
	stubDefinitionStatement :=
		m.astFactory.
			CreateAssignStatement(classNameDotPrototypeDotPropName, stubCall).
			SrcrefTreeIfMissing(methodDefinition)
	classDefiningStatement := NodeUtil_GetEnclosingStatement(classMembers)
	stubDefinitionStatement.InsertAfter(classDefiningStatement)

	originalScript := NodeUtil_GetEnclosingScript(methodDefinition)

	// remove the definition from the class
	methodDefinition.Detach()

	m.compiler.ReportChangeToEnclosingScope(classMembers)

	// Prepend unstub definition to the new location.
	// ClassName.prototype.propertyName = JSCompiler_unstubMethod(id, function(...) {...});
	classNameDotPrototypeDotPropName2 := classNameDotPrototypeDotPropName.CloneTree()
	functionNode := jutil.CheckNotNull(methodDefinition.GetOnlyChild())
	functionNode.Detach()
	unstubCall := m.CreateUnstubCall(functionNode, stubId)
	statementNode :=
		m.astFactory.
			CreateAssignStatement(classNameDotPrototypeDotPropName2, unstubCall).
			SrcrefTreeIfMissing(methodDefinition)
	destinationParent.AddChildToFront(statementNode)
	NodeUtil_AddFeaturesToScript(
		destinationParent, NodeUtil_GetFeatureSetOfScript(originalScript), m.compiler)
	m.compiler.ReportChangeToEnclosingScope(destinationParent)
}

// Java: CrossChunkMethodMotion.hasUnmovableRedeclaration(NameInfo, Property)
func CrossChunkMethodMotion_HasUnmovableRedeclaration(nameInfo *AnalyzePrototypeProperties_NameInfo, prop AnalyzePrototypeProperties_Property) bool {
	for _, symbol := range nameInfo.GetDeclarations().Slice() {
		if otherProp, ok := symbol.(AnalyzePrototypeProperties_Property); ok {
			// It is possible to do better here if the dependencies are well defined
			// but redefinitions are usually in optional chunks so it isn't likely
			// worth the effort to check.
			if prop != otherProp &&
				prop.GetRootVar() == otherProp.GetRootVar() &&
				prop.GetChunk() != otherProp.GetChunk() {
				return true
			}
		}
	}
	return false
}
