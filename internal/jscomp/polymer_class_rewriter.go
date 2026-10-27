// Port of com/google/javascript/jscomp/PolymerClassRewriter.java @ closure-compiler v20261021

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

package jscomp

import (
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	polymerClassRewriter_VIRTUAL_FILE                = "<PolymerClassRewriter.java>"
	polymerClassRewriter_POLYMER_ELEMENT_PROP_CONFIG = "PolymerElementProperties"
)

var (
	PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT = DiagnosticType_Error(
		"JSC_POLYMER_IMPLICIT_GLOBAL_CONFLICT",
		"Implicit global name for Polymer element conflicts with existing var {0}. Either"+
			" give the element a lhs or rename {0}. (Or move to class-based Polymer 2"+
			" elements)")

	PolymerClassRewriter_POLYMER_ELEMENT_CONFLICT = DiagnosticType_Error(
		"JSC_POLYMER_ELEMENT_CONFLICT",
		"Cannot generate correct types for Polymer call due to PolymerElement definition at"+
			" {0}:{2}:{1}.\n"+
			"Rename the local PolymerElement to avoid shadowing the PolymerElement externs.")
)

// PolymerClassRewriter rewrites a given call to Polymer({}) to a set of
// declarations and assignments which can be understood by the compiler.
//
// Java: final class PolymerClassRewriter
type PolymerClassRewriter struct {
	compiler                   AbstractCompiler
	externsInsertionRef        rhino.Node
	propertySinkExternInjected bool
}

// Java: PolymerClassRewriter(AbstractCompiler)
func NewPolymerClassRewriter(compiler AbstractCompiler) *PolymerClassRewriter {
	r := &PolymerClassRewriter{}
	r.propertySinkExternInjected = false
	r.compiler = compiler
	r.externsInsertionRef = compiler.GetSynthesizedExternsInput().GetAstRoot(compiler)
	return r
}

// Java: PolymerClassRewriter.isIIFE(Node)
func PolymerClassRewriter_IsIIFE(n rhino.Node) bool {
	call, ok := ast.AsCall(n)
	return ok && call.Callee().Node().IsFunction()
}

// Java: PolymerClassRewriter.isFunctionArgInGoogLoadModule(Node)
func PolymerClassRewriter_IsFunctionArgInGoogLoadModule(n rhino.Node) bool {
	if !n.IsFunction() {
		return false
	}

	parent := n.GetParent()
	return !parent.IsNil() && NodeUtil_IsBundledGoogModuleCall(parent)
}

// insertGeneratedDeclarationCodeToGlobalScope accepts declaration code
// generated for a nonGlobal Polymer call and inserts that into the AST
// depending on the enclosing scope of the Polymer call.
//
// enclosingNode is the enclosing scope of the Polymer call decided by the
// rewritePolymerCall; declarationCode the declaration code generated for
// Polymer call.
//
// Java: PolymerClassRewriter.insertGeneratedDeclarationCodeToGlobalScope(Node, Node)
func (r *PolymerClassRewriter) insertGeneratedDeclarationCodeToGlobalScope(
	enclosingNode rhino.Node, declarationCode rhino.Node) {
	switch enclosingNode.GetToken() {
	case rhino.Token_MODULE_BODY:
		{
			insertionPoint := r.getNodeForInsertion(enclosingNode.GetParent())
			insertionPoint.AddChildToFront(declarationCode)
			r.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(insertionPoint))
		}
	case rhino.Token_SCRIPT:
		{
			enclosingNode.AddChildToFront(declarationCode)
			r.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(enclosingNode))
		}
	case rhino.Token_CALL:
		{
			// This case represents only the Polymer calls which are enclosed inside an IIFE
			jutil.CheckState(PolymerClassRewriter_IsIIFE(enclosingNode))
			enclosingNodeForIIFE := NodeUtil_GetEnclosingNode(
				enclosingNode.GetParent(), func(node rhino.Node) bool { return node.IsScript() || node.IsModuleBody() })
			if enclosingNodeForIIFE.IsScript() {
				enclosingNodeForIIFE.AddChildToFront(declarationCode)
				r.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(enclosingNodeForIIFE))
			} else {
				jutil.CheckState(enclosingNodeForIIFE.IsModuleBody())
				insertionPoint := r.getNodeForInsertion(enclosingNodeForIIFE.GetParent())
				insertionPoint.AddChildToFront(declarationCode)
				r.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(insertionPoint))
			}
		}
	case rhino.Token_FUNCTION:
		{
			// This case represents only the Polymer calls that are inside a function which is an arg
			// to goog.loadModule
			jutil.CheckState(PolymerClassRewriter_IsFunctionArgInGoogLoadModule(enclosingNode))
			enclosingScript := NodeUtil_GetEnclosingScript(enclosingNode)
			insertionPoint := r.getNodeForInsertion(enclosingScript)
			insertionPoint.AddChildToFront(declarationCode)
			r.compiler.ReportChangeToChangeScope(insertionPoint)
		}
	default:
		panic(jutil.NewRuntimeException("Enclosing node for Polymer is incorrect"))
	}
}

// getNodeForInsertion returns a SCRIPT node in which to insert new global
// declarations
//
// New globals are put in the externs if converting a @typeSummary and in
// source code otherwise.
//
// Java: PolymerClassRewriter.getNodeForInsertion(Node)
func (r *PolymerClassRewriter) getNodeForInsertion(enclosingScript rhino.Node) rhino.Node {
	if NodeUtil_IsFromTypeSummary(enclosingScript) {
		return r.compiler.GetSynthesizedTypeSummaryInput().GetAstRoot(r.compiler)
	} else {
		return r.compiler.GetNodeForCodeInsertion(nil)
	}
}

// insertGeneratedPropsAndBehaviorCode accepts code generated for a nonGlobal
// Polymer call and inserts that code into the AST depending on the enclosing
// scope of the Polymer call.
//
// enclosingNode is the enclosing scope of the Polymer call decided by the
// rewritePolymerCall; statements the code generated for Polymer's properties
// and behavior.
//
// Java: PolymerClassRewriter.insertGeneratedPropsAndBehaviorCode(Node, Node)
func (r *PolymerClassRewriter) insertGeneratedPropsAndBehaviorCode(enclosingNode rhino.Node, statements rhino.Node) {
	switch enclosingNode.GetToken() {
	case rhino.Token_MODULE_BODY:
		{
			if enclosingNode.GetParent().GetBooleanProp(rhino.Node_GOOG_MODULE) {
				// The goog.module('ns'); call must remain the first statement in the module.
				insertionPoint := polymerClassRewriter_getInsertionPointForGoogModule(enclosingNode)
				enclosingNode.AddChildrenAfter(statements, insertionPoint)
			} else {
				enclosingNode.AddChildrenToFront(statements)
			}
		}
	case rhino.Token_SCRIPT:
		// If children are added to front, it will cause runtime error in the code because Polymer's
		// properties will be accessed before the Polymer object is defined. The Polymer object gets
		// defined in `insertGeneratedDeclarationCodeToGlobalScope`.
		enclosingNode.AddChildrenToBack(statements)
		r.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(enclosingNode))
	case rhino.Token_CALL:
		{
			// This case represents only the Polymer calls which are enclosed inside an IIFE
			jutil.CheckState(PolymerClassRewriter_IsIIFE(enclosingNode))
			functionNode := ast.MustFunction(ast.MustCall(enclosingNode).Callee().Node())
			functionBlock := functionNode.Body().Node()
			functionBlock.AddChildrenToFront(statements)
		}
	case rhino.Token_FUNCTION:
		// This case represents only the Polymer calls that are inside a function which is an arg
		// to goog.loadModule
		jutil.CheckState(PolymerClassRewriter_IsFunctionArgInGoogLoadModule(enclosingNode))
		functionBlock := ast.MustFunction(enclosingNode).Body().Node()
		insertionPoint := polymerClassRewriter_getInsertionPointForGoogModule(functionBlock)
		// Node insertionPoint will be null here if functionBlock does not contain a goog.module()
		// Missing goog.module inside the loadModule's functionBlock is semantically incorrect
		// That will cause the compiler to crash in closureRewriteModule pass.
		if !insertionPoint.IsNil() {
			functionBlock.AddChildrenAfter(statements, insertionPoint)
		}
	default:
	}
}

// RewritePolymerCall rewrites a given call to Polymer({}) to a set of
// declarations and assignments which can be understood by the compiler.
//
// cls is the extracted PolymerClassDefinition for the Polymer element created
// by this call; traversal the NodeTraversal used here to identify the scope
// in which Polymer exists.
//
// Java: PolymerClassRewriter.rewritePolymerCall(PolymerClassDefinition, NodeTraversal)
func (r *PolymerClassRewriter) RewritePolymerCall(cls *PolymerClassDefinition, traversal *NodeTraversal) {
	callParent := cls.definition.GetParent()
	// Determine whether we are in a Polymer({}) call at the top level versus in an assignment.
	exprRoot := callParent.GetParent()
	if callParent.IsExprResult() {
		exprRoot = callParent
	}
	jutil.CheckStateMsg1(NodeUtil_IsStatementParent(exprRoot.GetParent()), "%s", exprRoot.GetParent())
	objLit := rhino.CheckNotNull(cls.descriptor)

	// Add {@code @lends} to the object literal.
	objLitDoc := rhino.JSDocInfo_Builder_().ParseDocumentation()
	jsTypeExpression := rhino.NewJSTypeExpression(
		rhino.IR_String(*cls.target.GetQualifiedName()+".prototype").Srcref(exprRoot),
		jutil.ToJavaString(exprRoot.GetSourceFileName()))
	objLitDoc.RecordLends(jsTypeExpression)
	objLit.SetJSDocInfo(objLitDoc.Build())

	r.addTypesToFunctions(objLit, *cls.target.GetQualifiedName(), cls.defType)
	PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(objLit, r.compiler)
	PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(objLit, r.compiler)
	PolymerPassStaticUtils_ProtectObserverAndPropertyFunctionKeys(objLit)

	for _, prop := range cls.props {
		if prop.value.IsObjectLit() {
			PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(prop.value, r.compiler)
		}
	}

	// The propsAndBehaviorBlock holds code generated for the  Polymer's properties and behaviors
	propsAndBehaviorBlock := rhino.IR_Block()

	constructorDoc := r.getConstructorDoc(cls)

	// Remove the original constructor JS docs from the objlit.
	ctorKey := cls.constructor.value.GetParent()
	if !ctorKey.IsNil() {
		ctorKey.SetJSDocInfo(nil)
	}

	// Check for a conflicting definition of PolymerElement
	if !traversal.InGlobalScope() {
		polymerElement := traversal.GetScope().GetVar("PolymerElement")
		if polymerElement != nil && !polymerElement.GetScope().(*Scope).IsGlobal() {
			nameNode := polymerElement.GetNameNode()
			r.compiler.Report(
				JSError_MakeWithNode(
					cls.constructor.value,
					PolymerClassRewriter_POLYMER_ELEMENT_CONFLICT,
					jutil.ToJavaString(nameNode.GetSourceFileName()),
					strconv.Itoa(nameNode.GetLineno()),
					strconv.Itoa(nameNode.GetCharno())))
		}
	}

	declarationCode := r.generateDeclarationCode(exprRoot, cls, constructorDoc, traversal)
	basePath := *cls.target.GetQualifiedName() + ".prototype."
	r.appendBehaviorPropertiesToBlock(
		cls, propsAndBehaviorBlock, basePath /* isExternsBlock= */, false)
	r.appendPropertiesToBlock(
		cls.props, propsAndBehaviorBlock, basePath /* isExternsBlock= */, false)
	r.appendBehaviorMembersToBlock(cls, propsAndBehaviorBlock)
	readOnlyPropsAll := r.parseReadOnlyProperties(cls, propsAndBehaviorBlock)
	r.createExportsAndExterns(cls, readOnlyPropsAll)
	r.removePropertyDocs(objLit, cls.defType)

	propsAndBehaviorCode := propsAndBehaviorBlock.RemoveChildren()
	parent := exprRoot.GetParent()

	// Put the type declaration in to either the enclosing module scope, if in a module, or the
	// enclosing script node. Compiler support for local scopes like IIFEs is sometimes lacking but
	// module scopes are well-supported. If this is not in a module or the global scope it is likely
	// exported.
	if !traversal.InGlobalScope() && cls.hasGeneratedLhs && !cls.target.IsGetProp() {
		enclosingNode := NodeUtil_GetEnclosingNode(
			parent,
			func(node rhino.Node) bool {
				return node.IsScript() ||
					node.IsModuleBody() ||
					PolymerClassRewriter_IsIIFE(node) ||
					PolymerClassRewriter_IsFunctionArgInGoogLoadModule(node)
			})

		// For module, IIFE and goog.LoadModule enclosed Polymer calls, the declaration code and the
		// code generated from properties and behavior have to be hoisted in different places within
		// the AST. We want to insert the generated declarations to global scope, and insert the
		// propsAndbehaviorCode in the same scope. Hence, dealing with them separately.
		r.insertGeneratedDeclarationCodeToGlobalScope(enclosingNode, declarationCode)
		if !propsAndBehaviorCode.IsNil() {
			r.insertGeneratedPropsAndBehaviorCode(enclosingNode, propsAndBehaviorCode)
		}
	} else {
		beforeRoot := exprRoot.GetPrevious()
		if beforeRoot.IsNil() {
			if !propsAndBehaviorCode.IsNil() {
				parent.AddChildrenToFront(propsAndBehaviorCode)
			}
			parent.AddChildToFront(declarationCode)
		} else {
			if !propsAndBehaviorCode.IsNil() {
				parent.AddChildrenAfter(propsAndBehaviorCode, beforeRoot)
			}
			declarationCode.InsertAfter(beforeRoot)
		}
		r.compiler.ReportChangeToEnclosingScope(parent)
	}
	if !propsAndBehaviorCode.IsNil() {
		r.compiler.ReportChangeToEnclosingScope(propsAndBehaviorCode)
	}

	// Since behavior files might contain language features that aren't present in the class file,
	// we might need to update the FeatureSet.
	if cls.features != nil {
		scriptNode := NodeUtil_GetEnclosingScript(parent)
		oldFeatures, _ := scriptNode.GetProp(rhino.Node_FEATURE_SET).(*parser.FeatureSet)
		newFeatures := oldFeatures.Union(cls.features)
		if !newFeatures.Equals(oldFeatures) {
			scriptNode.PutProp(rhino.Node_FEATURE_SET, newFeatures)
			r.compiler.ReportChangeToChangeScope(scriptNode)
		}
	}

	if NodeUtil_IsNameDeclaration(exprRoot) {
		assignExpr := polymerClassRewriter_varToAssign(exprRoot)
		exprRoot.ReplaceWith(assignExpr)
		r.compiler.ReportChangeToEnclosingScope(assignExpr)
	}

	// If property renaming is enabled, wrap the properties object literal
	// in a reflection call so that the properties are renamed consistently
	// with the class members.
	if !cls.descriptor.IsNil() {
		props := NodeUtil_GetFirstPropMatchingKey(cls.descriptor, "properties")
		if !props.IsNil() && props.IsObjectLit() {
			r.addPropertiesConfigObjectReflection(cls, props)
		}
	}
}

// RewritePolymerClassDeclaration rewrites a class which extends
// Polymer.Element to a set of declarations and assignments which can be
// understood by the compiler.
//
// clazz is the class node; cls the extracted PolymerClassDefinition for the
// Polymer element created by this call.
//
// Java: PolymerClassRewriter.rewritePolymerClassDeclaration(Node, NodeTraversal, PolymerClassDefinition)
func (r *PolymerClassRewriter) RewritePolymerClassDeclaration(
	clazz rhino.Node, traversal *NodeTraversal, cls *PolymerClassDefinition) {

	if !cls.descriptor.IsNil() {
		r.addTypesToFunctions(cls.descriptor, *cls.target.GetQualifiedName(), cls.defType)
	}
	PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(
		NodeUtil_GetClassMembers(clazz), r.compiler)

	for _, prop := range cls.props {
		if prop.value.IsObjectLit() {
			PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(prop.value, r.compiler)
		}
	}

	// For simplicity add everything into a block, before adding it to the AST.
	block := rhino.IR_Block()

	r.appendBehaviorPropertiesToBlock(
		cls, block, *cls.target.GetQualifiedName()+".prototype." /* isExternsBlock= */, false)
	// For each Polymer property we found in the "properties" configuration object, append a
	// property declaration to the prototype (e.g. "/** @type {string} */ MyElement.prototype.foo").
	r.appendPropertiesToBlock(
		cls.props,
		block,
		*cls.target.GetQualifiedName()+".prototype." /* isExternsBlock= */, false)

	readOnlyProps := r.parseReadOnlyProperties(cls, block)
	r.createExportsAndExterns(cls, readOnlyProps)

	// If an external interface is required, mark the class as implementing it
	jsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(clazz)
	classInfo := rhino.JSDocInfo_Builder_MaybeCopyFrom(jsDocInfoNode.GetJSDocInfo())
	interfaceName := cls.GetInterfaceName(r.compiler.GetUniqueIdSupplier())
	interfaceType := rhino.NewJSTypeExpression(
		rhino.NewNodeWithChild(rhino.Token_BANG, rhino.IR_String(interfaceName)).SrcrefTree(jsDocInfoNode), polymerClassRewriter_VIRTUAL_FILE)
	classInfo.RecordImplementedInterface(interfaceType)
	jsDocInfoNode.SetJSDocInfo(classInfo.Build())

	insertAfterReference := NodeUtil_GetEnclosingStatement(clazz)
	if block.HasChildren() {
		r.removePropertyDocs(cls.descriptor, cls.defType)
		newInsertAfterReference := block.GetLastChild()
		insertAfterReference.
			GetParent().
			AddChildrenAfter(block.RemoveChildren(), insertAfterReference)
		r.compiler.ReportChangeToEnclosingScope(insertAfterReference)
		insertAfterReference = newInsertAfterReference
	}

	polymerClassRewriter_addReturnTypeIfMissing(cls, "is", rhino.NewJSTypeExpression(rhino.IR_String("string"), polymerClassRewriter_VIRTUAL_FILE))

	type_ := rhino.NewNode(rhino.Token_BANG)
	array := rhino.IR_String("Array")
	type_.AddChildToBack(array)
	arrayTemplateType := rhino.NewNodeWithChild(rhino.Token_BLOCK, rhino.IR_String("string"))
	array.AddChildToBack(arrayTemplateType)
	polymerClassRewriter_addReturnTypeIfMissing(cls, "observers", rhino.NewJSTypeExpression(type_, polymerClassRewriter_VIRTUAL_FILE))
	polymerClassRewriter_addReturnTypeIfMissing(
		cls,
		"properties",
		rhino.NewJSTypeExpression(rhino.IR_String(polymerClassRewriter_POLYMER_ELEMENT_PROP_CONFIG), polymerClassRewriter_VIRTUAL_FILE))

	// If property renaming is enabled, wrap the properties object literal
	// in a reflection call so that the properties are renamed consistently
	// with the class members.
	//
	// Also add reflection and sinks for computed properties and complex observers
	// and switch simple observers to direct function references.
	if !cls.descriptor.IsNil() {
		r.convertSimpleObserverStringsToReferences(cls)
		r.addPropertiesConfigObjectReflection(cls, cls.descriptor)
	}

	if !cls.descriptor.IsNil() {
		r.validateComputedPropertiesReflectionCalls(cls)
		r.validateComplexObserverReflectionCalls(cls)
	}
}

// polymerClassRewriter_addReturnTypeIfMissing adds return type information
// to class getters
//
// Java: PolymerClassRewriter.addReturnTypeIfMissing(PolymerClassDefinition, String, JSTypeExpression)
func polymerClassRewriter_addReturnTypeIfMissing(
	cls *PolymerClassDefinition, getterPropName string, jsType *rhino.JSTypeExpression) {
	classMembers := NodeUtil_GetClassMembers(cls.definition)
	getter := NodeUtil_GetFirstGetterMatchingKey(classMembers, getterPropName)
	if !getter.IsNil() {
		info := NodeUtil_GetBestJSDocInfo(getter)
		if info == nil || !info.HasReturnType() {
			builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(info)
			builder.RecordReturnType(jsType)
			jsType.GetRoot().SrcrefTreeIfMissing(getter)
			getter.SetJSDocInfo(builder.Build())
		}
	}
}

// addPropertiesConfigObjectReflection wraps the properties config object in
// an objectReflect call
//
// Java: PolymerClassRewriter.addPropertiesConfigObjectReflection(PolymerClassDefinition, Node)
func (r *PolymerClassRewriter) addPropertiesConfigObjectReflection(
	cls *PolymerClassDefinition, propertiesLiteral rhino.Node) {
	rhino.CheckNotNull(propertiesLiteral)
	jutil.CheckState(propertiesLiteral.IsObjectLit())

	parent := propertiesLiteral.GetParent()

	objReflectCall := rhino.IR_Call(
		NodeUtil_NewQName(r.compiler, "$jscomp.reflectObject"),
		cls.target.CloneTree(),
		propertiesLiteral.Detach()).
		SrcrefTreeIfMissing(propertiesLiteral)
	parent.AddChildToFront(objReflectCall)
	r.compiler.ReportChangeToEnclosingScope(parent)
	r.compiler.GetRuntimeJsLibManager().EnsureLibraryInjected("util/reflectobject", false)
}

// addTypesToFunctions adds an @this annotation to all functions in the
// objLit.
//
// Java: PolymerClassRewriter.addTypesToFunctions(Node, String, PolymerClassDefinition.DefinitionType)
func (r *PolymerClassRewriter) addTypesToFunctions(
	objLit rhino.Node, thisType string, defType PolymerClassDefinition_DefinitionType) {
	for member := range ast.MustObjectLit(objLit).Members() {
		keyNode := member.Node()
		var value rhino.Node
		if pkm, ok := member.AsPlainKeyMember(); ok {
			value = pkm.Value().Node()
		} else if cp, ok := member.AsComputedProp(); ok {
			if val, hasVal := cp.Value(); hasVal {
				value = val.Node()
			}
		}
		if value.IsFunction() {
			fnDoc := rhino.JSDocInfo_Builder_MaybeCopyFrom(keyNode.GetJSDocInfo())
			fnDoc.RecordThisType(
				rhino.NewJSTypeExpression(
					rhino.NewNodeWithChild(rhino.Token_BANG, rhino.IR_String(thisType)).SrcrefTree(keyNode), polymerClassRewriter_VIRTUAL_FILE))
			keyNode.SetJSDocInfo(fnDoc.Build())
		}
	}

	// Add @this and @return to default property values.
	for _, property := range PolymerPassStaticUtils_ExtractProperties(
		objLit, defType, r.compiler /* constructor= */, rhino.NilNode) {
		if !property.value.IsObjectLit() {
			continue
		}

		defaultValue := NodeUtil_GetFirstPropMatchingKey(property.value, "value")
		if defaultValue.IsNil() || !defaultValue.IsFunction() {
			continue
		}
		defaultValueKey := defaultValue.GetParent()
		fnDoc := rhino.JSDocInfo_Builder_MaybeCopyFrom(defaultValueKey.GetJSDocInfo())
		fnDoc.RecordThisType(
			rhino.NewJSTypeExpression(
				rhino.NewNodeWithChild(rhino.Token_BANG, rhino.IR_String(thisType)).SrcrefTree(defaultValueKey), polymerClassRewriter_VIRTUAL_FILE))
		fnDoc.RecordReturnType(PolymerPassStaticUtils_GetTypeFromProperty(property, r.compiler))
		defaultValueKey.SetJSDocInfo(fnDoc.Build())
	}
}

// parseReadOnlyProperties generates the _set* setters for readonly
// properties and appends them to the given block.
//
// Returns a List of all readonly properties.
//
// Java: PolymerClassRewriter.parseReadOnlyProperties(PolymerClassDefinition, Node)
func (r *PolymerClassRewriter) parseReadOnlyProperties(
	cls *PolymerClassDefinition, block rhino.Node) []*PolymerPass_MemberDefinition {
	qualifiedPath := *cls.target.GetQualifiedName() + ".prototype."
	var readOnlyProps []*PolymerPass_MemberDefinition

	for _, prop := range cls.props {
		// Generate the setter for readOnly properties.
		if prop.value.IsObjectLit() {
			readOnlyValue := NodeUtil_GetFirstPropMatchingKey(prop.value, "readOnly")
			if !readOnlyValue.IsNil() && readOnlyValue.IsTrue() {
				setter := r.makeReadOnlySetter(prop, qualifiedPath)
				setter.SrcrefTreeIfMissing(prop.name)
				block.AddChildToBack(setter)
				readOnlyProps = append(readOnlyProps, prop)
			}
		}
	}

	if cls.behaviorProps != nil {
		for prop := range cls.behaviorProps.Keys() {
			// Generate the setter for readOnly properties.
			if prop.value.IsObjectLit() {
				readOnlyValue := NodeUtil_GetFirstPropMatchingKey(prop.value, "readOnly")
				if !readOnlyValue.IsNil() && readOnlyValue.IsTrue() {
					setter := r.makeReadOnlySetter(prop, qualifiedPath)
					setter.SrcrefTreeIfMissing(prop.name)
					block.AddChildToBack(setter)
					readOnlyProps = append(readOnlyProps, prop)
				}
			}
		}
	}
	return readOnlyProps
}

// getConstructorDoc returns the proper constructor doc for the Polymer call.
//
// Java: PolymerClassRewriter.getConstructorDoc(PolymerClassDefinition)
func (r *PolymerClassRewriter) getConstructorDoc(cls *PolymerClassDefinition) *rhino.JSDocInfo_Builder {
	constructorDoc := rhino.JSDocInfo_Builder_MaybeCopyFrom(cls.constructor.info)
	constructorDoc.RecordConstructor()

	baseType := rhino.NewJSTypeExpression(
		rhino.NewNodeWithChild(rhino.Token_BANG, rhino.IR_String(PolymerPassStaticUtils_GetPolymerElementType(cls))).
			SrcrefTree(cls.definition),
		polymerClassRewriter_VIRTUAL_FILE)
	constructorDoc.RecordBaseType(baseType)

	interfaceName := cls.GetInterfaceName(r.compiler.GetUniqueIdSupplier())
	interfaceType := rhino.NewJSTypeExpression(
		rhino.NewNodeWithChild(rhino.Token_BANG, rhino.IR_String(interfaceName)).SrcrefTree(cls.definition),
		polymerClassRewriter_VIRTUAL_FILE)
	constructorDoc.RecordImplementedInterface(interfaceType)

	return constructorDoc
}

// generateDeclarationCode appends var declaration code created from the
// Polymer call to the given block
//
// Java: PolymerClassRewriter.generateDeclarationCode(Node, PolymerClassDefinition, JSDocInfo.Builder, NodeTraversal)
func (r *PolymerClassRewriter) generateDeclarationCode(
	exprRoot rhino.Node,
	cls *PolymerClassDefinition,
	constructorDoc *rhino.JSDocInfo_Builder,
	traversal *NodeTraversal) rhino.Node {
	if cls.target.IsGetProp() {
		// foo.bar = Polymer({...});
		assign := rhino.IR_Assign(cls.target.CloneTree(), cls.constructor.value.CloneTree())
		NodeUtil_MarkNewScopesChanged(assign, r.compiler)
		assign.SetJSDocInfo(constructorDoc.Build())
		exprResult := rhino.IR_ExprResult(assign)
		exprResult.SrcrefTreeIfMissing(cls.target)
		return exprResult
	} else {
		// var foo = Polymer({...}); OR Polymer({...});
		var_ := rhino.IR_VarWithValue(cls.target.CloneTree(), cls.constructor.value.CloneTree())
		NodeUtil_MarkNewScopesChanged(var_, r.compiler)
		var_.SrcrefTreeIfMissing(exprRoot)
		var_.SetJSDocInfo(constructorDoc.Build())
		name := cls.target.GetString()
		existingVar := traversal.GetScope().GetSlot(name)
		if existingVar != nil && cls.hasGeneratedLhs {
			r.compiler.Report(JSError_MakeWithNode(cls.constructor.value, PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT, name))
		}
		return var_
	}
}

// replaceJSDocAndAddNewVars replaces JSDoc types in externs with unknown
// type. For a JSDoc like @type {{ propertyName : string }}, collects all such
// "propertyName"s, and generates extern vars with an attached
// {{propertyName: ?}} JsDoc. This is to prevent renaming of vars in source
// code with the same names as propertyNames.
//
// Java: PolymerClassRewriter.replaceJSDocAndAddNewVars(MemberDefinition, JSTypeExpression, Node)
func (r *PolymerClassRewriter) replaceJSDocAndAddNewVars(
	prop *PolymerPass_MemberDefinition, propType *rhino.JSTypeExpression, block rhino.Node) *rhino.JSDocInfo {
	infoBuilder := rhino.JSDocInfo_Builder_MaybeCopyFrom(prop.info)
	infoBuilder.RecordType(propType)
	origInfo := infoBuilder.Build()
	propertyNames := propType.GetRecordPropertyNames()
	r.createVarsInExternsBlock(block, propertyNames, propType, prop)
	unknown := rhino.NewJSTypeExpression(rhino.NewNode(rhino.Token_QMARK), propType.GetSourceName())
	newInfoBuilder := rhino.JSDocInfo_Builder_MaybeCopyFromWithNewType(origInfo, unknown)
	return newInfoBuilder.Build()
}

// getPropertyNode returns a node from a property's definition in the Polymer
// element or behavior
//
// Java: PolymerClassRewriter.getPropertyNode(MemberDefinition, String)
func (r *PolymerClassRewriter) getPropertyNode(prop *PolymerPass_MemberDefinition, basePath string) rhino.Node {
	// If a property string is quoted, make sure the added prototype properties are also quoted
	if prop.name.IsQuotedStringKey() {
		return rhino.NilNode
	}
	propertyNode := rhino.IR_ExprResult(NodeUtil_NewQName(r.compiler, basePath+prop.name.GetString()))
	propertyNode.SrcrefTreeIfMissing(prop.name)

	return propertyNode
}

// appendBehaviorPropertiesToBlock iterates through all the behaviors of this
// polymer call, and appends the properties of each behavior to the given
// block.
//
// Java: PolymerClassRewriter.appendBehaviorPropertiesToBlock(PolymerClassDefinition, Node, String, boolean)
func (r *PolymerClassRewriter) appendBehaviorPropertiesToBlock(
	cls *PolymerClassDefinition, block rhino.Node, basePath string, isExternsBlock bool) {
	if len(cls.behaviors) == 0 || cls.behaviorProps == nil {
		return
	}

	for prop, behavior := range cls.behaviorProps.All() {
		propertyNode := r.getPropertyNode(prop, basePath)
		if propertyNode.IsNil() {
			continue
		}

		propType := PolymerPassStaticUtils_GetTypeFromProperty(prop, r.compiler)
		if propType == nil {
			continue
		}

		var info *rhino.JSDocInfo
		if isExternsBlock {
			info = r.replaceJSDocAndAddNewVars(prop, propType, block)
		} else {
			infoBuilder := r.getJSDocInfoBuilderForBehavior(behavior, prop)
			infoBuilder.RecordType(propType)
			info = infoBuilder.Build()
		}
		ast.MustExprStmt(propertyNode).Expr().Node().SetJSDocInfo(info)
		block.AddChildToBack(propertyNode)
	}
}

// appendPropertiesToBlock appends all of the given properties to the given
// block.
//
// Java: PolymerClassRewriter.appendPropertiesToBlock(List<MemberDefinition>, Node, String, boolean)
func (r *PolymerClassRewriter) appendPropertiesToBlock(
	props []*PolymerPass_MemberDefinition, block rhino.Node, basePath string, isExternsBlock bool) {
	for _, prop := range props {

		propertyNode := r.getPropertyNode(prop, basePath)
		if propertyNode.IsNil() {
			continue
		}

		propType := PolymerPassStaticUtils_GetTypeFromProperty(prop, r.compiler)
		if propType == nil {
			continue
		}

		var info *rhino.JSDocInfo
		if isExternsBlock {
			info = r.replaceJSDocAndAddNewVars(prop, propType, block)
		} else {
			infoBuilder := rhino.JSDocInfo_Builder_MaybeCopyFrom(prop.info)
			infoBuilder.RecordType(propType)
			info = infoBuilder.Build()
		}
		ast.MustExprStmt(propertyNode).Expr().Node().SetJSDocInfo(info)
		block.AddChildToBack(propertyNode)
	}
}

// createVarsInExternsBlock: for a JSDoc like @type {{ propertyName : string }},
// collects all such "propertyName"s, and generates extern vars with an
// attached {{propertyName: ?}} JsDoc. This is to prevent renaming of vars in
// source code with the same names as propertyNames.
//
// If we do not preserve the property names, then 540 targets get broken on
// TGP testing. This indicates that those targets probably accidentally relied
// on properties not being renamed, and we did not find it important to clean
// up all those targets' JS source.
//
// Java: PolymerClassRewriter.createVarsInExternsBlock(Node, ImmutableSet<String>, JSTypeExpression, MemberDefinition)
func (r *PolymerClassRewriter) createVarsInExternsBlock(
	block rhino.Node,
	propertyNames *jutil.LinkedHashSet[string],
	propType *rhino.JSTypeExpression,
	prop *PolymerPass_MemberDefinition) {

	for propName := range propertyNames.All() {
		varName := "PolymerDummyVar" + r.compiler.GetUniqueNameIdSupplier()()
		n := rhino.Node_NewStringWithToken(rhino.Token_NAME, varName)
		var_ := rhino.NewNode(rhino.Token_VAR)
		var_.AddChildToBack(n)

		// Forming @type {{ propertyName : ? }}
		newType := polymerClassRewriter_createNewTypeExpressionForExtern(propName, propType.GetSourceName(), prop)

		oldInfoBuilder := rhino.JSDocInfo_Builder_MaybeCopyFrom(prop.info)
		info := oldInfoBuilder.Build()

		newInfo := rhino.JSDocInfo_Builder_CopyFromWithNewType(info, newType)
		var_.SetJSDocInfo(newInfo.Build())
		block.AddChildToBack(var_)
	}
}

// polymerClassRewriter_createNewTypeExpressionForExtern creates a new type
// expression for JSDoc like @type {{ propertyName : ? }}
//
// Java: PolymerClassRewriter.createNewTypeExpressionForExtern(String, String, MemberDefinition)
func polymerClassRewriter_createNewTypeExpressionForExtern(
	propName string, sourceName string, prop *PolymerPass_MemberDefinition) *rhino.JSTypeExpression {
	leftCurly := rhino.NewNode(rhino.Token_LC)
	leftBracket := rhino.NewNode(rhino.Token_LB)
	colon := rhino.NewNode(rhino.Token_COLON)
	propertyName := rhino.Node_NewString(propName)
	propertyName.SetToken(rhino.Token_STRING_KEY)
	unknown := rhino.NewNode(rhino.Token_QMARK)
	colon.AddChildToBack(propertyName)
	colon.AddChildToBack(unknown)
	leftBracket.AddChildToBack(colon)
	leftCurly.AddChildToBack(leftBracket)
	leftCurly.SrcrefTreeIfMissing(prop.name)
	return rhino.NewJSTypeExpression(leftCurly, sourceName)
}

// removePropertyDocs removes all JSDocs from properties of a class
// definition
//
// Java: PolymerClassRewriter.removePropertyDocs(Node, PolymerClassDefinition.DefinitionType)
func (r *PolymerClassRewriter) removePropertyDocs(
	objLit rhino.Node, defType PolymerClassDefinition_DefinitionType) {
	for _, prop := range PolymerPassStaticUtils_ExtractProperties(
		objLit, defType, r.compiler /* constructor= */, rhino.NilNode) {
		prop.name.SetJSDocInfo(nil)
	}
}

// Java: PolymerClassRewriter.getJSDocInfoBuilderForBehavior(BehaviorDefinition, MemberDefinition)
func (r *PolymerClassRewriter) getJSDocInfoBuilderForBehavior(
	behavior *PolymerBehaviorExtractor_BehaviorDefinition, behaviorFunctionOrProp *PolymerPass_MemberDefinition) *rhino.JSDocInfo_Builder {
	var info *rhino.JSDocInfo_Builder
	if !behavior.isGlobalDeclaration &&
		behaviorFunctionOrProp.info != nil &&
		behaviorFunctionOrProp.info.ContainsTypeDeclaration() {

		if !behavior.behaviorModule.IsNil() {
			// Replace module local names in @type, @return and @param with unknown type
			info = rhino.JSDocInfo_Builder_MaybeCopyFromAndReplaceNames(
				behaviorFunctionOrProp.info, behavior.GetModuleLocalNames(r.compiler))
		} else {
			info = rhino.JSDocInfo_Builder_MaybeCopyFrom(behaviorFunctionOrProp.info)
		}
	} else {
		info = rhino.JSDocInfo_Builder_MaybeCopyFrom(behaviorFunctionOrProp.info)
	}
	return info
}

// appendBehaviorMembersToBlock appends all required behavior functions and
// non-property members to the given block.
//
// Java: PolymerClassRewriter.appendBehaviorMembersToBlock(PolymerClassDefinition, Node)
func (r *PolymerClassRewriter) appendBehaviorMembersToBlock(cls *PolymerClassDefinition, block rhino.Node) {
	qualifiedPath := *cls.target.GetQualifiedName() + ".prototype."
	nameToExprResult := jutil.NewLinkedHashMap[string, rhino.Node]()
	for _, behavior := range cls.behaviors {
		for _, behaviorFunction := range behavior.functionsToCopy {
			fnName := behaviorFunction.name.GetString()
			// Don't copy functions already defined by the element itself.
			if !NodeUtil_GetFirstPropMatchingKey(cls.descriptor, fnName).IsNil() {
				continue
			}

			// Avoid copying over the same function twice. The last definition always wins.
			if nameToExprResult.ContainsKey(fnName) {
				nameToExprResult.Get(fnName).Detach()
			}

			fnValue := behaviorFunction.value.CloneTree()
			NodeUtil_MarkNewScopesChanged(fnValue, r.compiler)
			exprResult := rhino.IR_ExprResult(rhino.IR_Assign(NodeUtil_NewQName(r.compiler, qualifiedPath+fnName), fnValue))
			exprResult.SrcrefTreeIfMissing(behaviorFunction.name)

			info := r.getJSDocInfoBuilderForBehavior(behavior, behaviorFunction)

			// If the function in the behavior is @protected, switch it to @public so that
			// we don't get a visibility warning. This is a bit of a hack but easier than
			// making the type system understand that methods are "inherited" from behaviors.
			if behaviorFunction.info != nil &&
				behaviorFunction.info.GetVisibility() == rhino.JSDocInfo_Visibility_PROTECTED {
				info.OverwriteVisibility(rhino.JSDocInfo_Visibility_PUBLIC)
			}

			// Behaviors whose declarations are not in the global scope may contain references to
			// symbols which do not exist in the element's scope. Only copy a function stub.
			if !behavior.isGlobalDeclaration {
				fn := ast.MustFunction(fnValue)
				body := fn.Body().Node()
				if fn.IsArrow() && !body.IsBlock() {
					// replace `() => <someExpr>` with `() => undefined`
					body.ReplaceWith(NodeUtil_NewUndefinedNode(body))
				} else {
					body.RemoveChildren()
				}
				// Remove any non-named parameters, which may reference locals.
				paramIndex := 0
				for param := range ast.ReadAhead(fn.Params().FirstParam()) {
					polymerClassRewriter_makeParamSafe(param.Node(), paramIndex)
					paramIndex++
				}
				if behaviorFunction.info == nil || behaviorFunction.info.GetReturnType() == nil {
					// Record a return type of ? to avoid the type inferencer assuming this function returns
					// 'undefined' just because it's a stub.
					info.RecordReturnType(
						rhino.NewJSTypeExpression(
							rhino.NewNode(rhino.Token_QMARK).Srcref(fnValue), jutil.ToJavaString(fnValue.GetSourceFileName())))
				}
			}

			exprResult.GetFirstChild().SetJSDocInfo(info.Build())
			block.AddChildToBack(exprResult)
			nameToExprResult.Put(fnName, exprResult)
		}

		// Copy other members.
		for _, behaviorProp := range behavior.nonPropertyMembersToCopy {
			propName := behaviorProp.name.GetString()
			if nameToExprResult.ContainsKey(propName) {
				nameToExprResult.Get(propName).Detach()
			}

			exprResult := rhino.IR_ExprResult(NodeUtil_NewQName(r.compiler, qualifiedPath+propName))
			exprResult.SrcrefTree(behaviorProp.name)

			info := r.getJSDocInfoBuilderForBehavior(behavior, behaviorProp)

			if behaviorProp.name.IsGetterDef() {
				info = rhino.JSDocInfo_Builder_().ParseDocumentation()
				if behaviorProp.info != nil && behaviorProp.info.GetReturnType() != nil {
					info.RecordType(behaviorProp.info.GetReturnType())
					if !behaviorProp.enclosingModule.IsNil() {
						info = rhino.JSDocInfo_Builder_MaybeCopyFromAndReplaceNames(
							info.Build(), behavior.GetModuleLocalNames(r.compiler))
					}
				}
			}

			exprResult.GetFirstChild().SetJSDocInfo(info.Build())
			block.AddChildToBack(exprResult)
			nameToExprResult.Put(propName, exprResult)
		}
	}
}

// polymerClassRewriter_makeParamSafe removes any potential local names
// referenced within a formal parameter
//
// Java: PolymerClassRewriter.makeParamSafe(Node, int)
func polymerClassRewriter_makeParamSafe(param rhino.Node, index int) {
	if rest, ok := ast.AsIterRest(param); ok {
		// The lhs may be a destructuring pattern.
		param = rest.Target().Node()
	} else if defVal, ok := ast.AsDefaultValue(param); ok {
		// Replace default value with void 0, then look at the lhs
		target := defVal.Target().Node()
		defVal.Value().Node().ReplaceWith(NodeUtil_NewUndefinedNode(param))
		param = target
	}

	if param.IsDestructuringPattern() {
		param.ReplaceWith(rhino.IR_Name("param$polymer$" + strconv.Itoa(index)).Srcref(param))
	}
}

// makeReadOnlySetter adds the generated setter for a readonly property.
//
// See https://www.polymer-project.org/0.8/docs/devguide/properties.html#read-only
//
// Java: PolymerClassRewriter.makeReadOnlySetter(MemberDefinition, String)
func (r *PolymerClassRewriter) makeReadOnlySetter(prop *PolymerPass_MemberDefinition, qualifiedPath string) rhino.Node {
	propName := prop.name.GetString()
	setterName :=
		"_set" + jutil.StringToUpperCaseRoot(jutil.StringSubstring(propName, 0, 1)) + jutil.StringSubstringFrom(propName, 1)
	fnNode := rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(rhino.IR_Name(propName)), rhino.IR_Block())
	r.compiler.ReportChangeToChangeScope(fnNode)
	exprResNode := rhino.IR_ExprResult(rhino.IR_Assign(NodeUtil_NewQName(r.compiler, qualifiedPath+setterName), fnNode))

	info := rhino.JSDocInfo_Builder_().ParseDocumentation()
	// This is overriding a generated function which was added to the interface in
	// {@code createExportsAndExterns}.
	info.RecordOverride()
	propType := PolymerPassStaticUtils_GetTypeFromProperty(prop, r.compiler)
	info.RecordParameter(propName, propType)
	exprResNode.GetFirstChild().SetJSDocInfo(info.Build())
	return exprResNode
}

// createExportsAndExterns creates exports and externs to protect element
// properties and methods from renaming and dead code removal.
//
// Since Polymer templates, observers, and computed properties rely on string
// references to element properties and methods, and because we don't yet
// have a way to update those references reliably, we instead export or
// extern them.
//
// For properties, we create a new interface called
// `Polymer<ElementName>Interface`, add all element properties to it, mark
// that the element class `@implements` this interface, and add the interface
// to the Closure externs.
//
// For methods, we instead append to `Object.prototype` in the externs using
// `@export` annotations. This approach is a compromise, with the following
// alternatives considered:
//
// Alternative 1: Add methods to our generated
// `Polymer<ElementName>Interface` in the externs. Pro: More optimal than
// `Object.prototype` when type-aware optimizations are enabled. Con 1: When a
// class `@implements` an interface, and when `report_missing_override` is
// enabled, any method on the class that is also in the interface must have an
// `@override` annotation, which means we generate a spurious warning for all
// methods. Con 2: An unresolved bug was encountered (b/115942961) relating to
// a mismatch between the signatures of the class and the generated interface.
//
// Alternative 2: Generate goog.exportProperty calls, which causes aliases on
// the prototype from original to optimized names to be set. Pro: Compiled
// code can still use the optimized name. Con: In practice, for Polymer
// applications, we see a net increase in bundle size due to the high number
// of new `Foo.prototype.originalName = Foo.prototype.z` expressions.
//
// Alternative 3: Append directly to the `Object.prototype` externs, instead
// of using `@export` annotations for the GenerateExports pass. Pro: Doesn't
// depend on the `generate_exports` and `export_local_property_definitions`
// flags. Con: The PolymerPass runs in the type checking phase, so modifying
// `Object.prototype` here causes unwanted type checking effects, such as
// allowing the method to be called on any object, and generating incorrect
// warnings when `report_missing_override` is enabled.
//
// Java: PolymerClassRewriter.createExportsAndExterns(PolymerClassDefinition, List<MemberDefinition>)
func (r *PolymerClassRewriter) createExportsAndExterns(
	cls *PolymerClassDefinition, readOnlyProps []*PolymerPass_MemberDefinition) {
	block := rhino.IR_Block()

	interfaceName := cls.GetInterfaceName(r.compiler.GetUniqueIdSupplier())
	fnNode := NodeUtil_EmptyFunction()
	r.compiler.ReportChangeToChangeScope(fnNode)
	varNode := rhino.IR_VarWithValue(NodeUtil_NewQName(r.compiler, interfaceName), fnNode)

	info := rhino.JSDocInfo_Builder_().ParseDocumentation()
	info.RecordInterface()
	varNode.SetJSDocInfo(info.Build())
	block.AddChildToBack(varNode)
	interfaceBasePath := interfaceName + ".prototype."

	r.appendBehaviorPropertiesToBlock(cls, block, interfaceBasePath /* isExternsBlock= */, true)
	r.appendPropertiesToBlock(cls.props, block, interfaceBasePath /* isExternsBlock= */, true)
	// Methods from behaviors were not already added to our element definition, so we need to
	// export those in addition to methods defined directly on the element. Note it's possible
	// and valid for two behaviors, or a behavior and an element, to implement the same method,
	// so we de-dupe by name. We're not checking that the signatures are compatible in the way
	// that normal class inheritance would, but that's not easy to do since these aren't classes.
	// Class mixins replace Polymer behaviors and are supported directly by Closure, so new code
	// should use those instead.
	uniqueMethods := jutil.NewLinkedHashMap[string, *PolymerPass_MemberDefinition]()
	if cls.behaviors != nil {
		for _, behavior := range cls.behaviors {
			for _, method := range behavior.functionsToCopy {
				uniqueMethods.Put(method.name.GetString(), method)
			}
		}
	}
	for _, method := range cls.methods {
		uniqueMethods.Put(method.name.GetString(), method)
	}
	for method := range uniqueMethods.Values() {
		r.addMethodToObjectExternsUsingExportAnnotation(cls, method)
	}

	for _, prop := range readOnlyProps {
		// Add all _set* functions to avoid renaming.
		propName := prop.name.GetString()
		setterName :=
			"_set" + jutil.StringToUpperCaseRoot(jutil.StringSubstring(propName, 0, 1)) + jutil.StringSubstringFrom(propName, 1)
		setterExprNode := rhino.IR_ExprResult(NodeUtil_NewQName(r.compiler, interfaceBasePath+setterName))

		setterInfo := rhino.JSDocInfo_Builder_().ParseDocumentation()
		propType := PolymerPassStaticUtils_GetTypeFromProperty(prop, r.compiler)
		unknown := rhino.NewJSTypeExpression(rhino.NewNode(rhino.Token_QMARK), propType.GetSourceName())
		setterInfo.RecordParameter(propName, unknown)
		setterExprNode.GetFirstChild().SetJSDocInfo(setterInfo.Build())

		block.AddChildToBack(setterExprNode)
	}

	var stmts rhino.Node
	if NodeUtil_IsFromTypeSummary(cls.input.GetAstRoot(r.compiler)) {
		typeSummaryInsertionRef := r.compiler.GetSynthesizedTypeSummaryInput().GetAstRoot(r.compiler)
		block.SrcrefTreeIfMissing(typeSummaryInsertionRef)
		stmts = block.RemoveChildren()
		typeSummaryInsertionRef.AddChildrenToBack(stmts)
	} else {
		block.SrcrefTreeIfMissing(r.externsInsertionRef)
		stmts = block.RemoveChildren()
		r.externsInsertionRef.AddChildrenToBack(stmts)
	}

	r.compiler.ReportChangeToEnclosingScope(stmts)
}

// addMethodToObjectExternsUsingExportAnnotation adds a method to
// `Object.prototype` in the externs by inserting a `GETPROP` expression with
// an `@export` annotation into the program.
//
// This relies on the `--generate_exports` and
// `export_local_property_definitions` flags to enable the GenerateExports
// pass, which will add properties exported in this way to `Object.prototype`
// in the externs, thus preventing renaming and dead code removal. Note that
// GenerateExports runs after type checking, so extending `Object.prototype`
// does not cause unwanted type checking effects.
//
// Java: PolymerClassRewriter.addMethodToObjectExternsUsingExportAnnotation(PolymerClassDefinition, MemberDefinition)
func (r *PolymerClassRewriter) addMethodToObjectExternsUsingExportAnnotation(
	cls *PolymerClassDefinition, method *PolymerPass_MemberDefinition) {
	getprop := NodeUtil_NewQName(
		r.compiler, *cls.target.GetQualifiedName()+".prototype."+method.name.GetString())
	info := rhino.JSDocInfo_Builder_().ParseDocumentation()
	if method.info != nil {
		// We need to preserve visibility, but other JSDoc doesn't matter (and can cause
		// re-declaration errors).
		info.RecordVisibility(method.info.GetVisibility())
	}
	info.RecordExport()
	getprop.SetJSDocInfo(info.Build())
	expression := rhino.IR_ExprResult(getprop).SrcrefTreeIfMissing(method.name)
	// Walk up until we find a statement we can insert after.
	insertAfter := cls.definition
	for !NodeUtil_IsStatementBlock(insertAfter.GetParent()) {
		insertAfter = insertAfter.GetParent()
	}
	expression.InsertAfter(insertAfter)
	r.compiler.ReportChangeToEnclosingScope(expression)
}

// polymerClassRewriter_varToAssign returns an assign replacing the
// equivalent var or let declaration.
//
// Java: PolymerClassRewriter.varToAssign(Node)
func polymerClassRewriter_varToAssign(var_ rhino.Node) rhino.Node {
	lhs := ast.MustNameDeclaration(var_).MustFirstDeclarator().Node()
	assign := rhino.IR_Assign(lhs.CloneNode(), lhs.RemoveFirstChild())
	return rhino.IR_ExprResult(assign).SrcrefTreeIfMissing(var_)
}

// convertSimpleObserverStringsToReferences converts property observer
// strings to direct function references.
//
// From: `observer: '_observerName'` To:
// `observer: ClassName.prototype._observerName`
//
// Java: PolymerClassRewriter.convertSimpleObserverStringsToReferences(PolymerClassDefinition)
func (r *PolymerClassRewriter) convertSimpleObserverStringsToReferences(cls *PolymerClassDefinition) {
	for _, prop := range cls.props {
		if prop.value.IsObjectLit() {
			observer := NodeUtil_GetFirstPropMatchingKey(prop.value, "observer")
			if !observer.IsNil() && observer.IsStringLit() {
				observerDirectReference := rhino.IR_GetpropWithMoreProps(cls.target.CloneTree(), "prototype", observer.GetString()).
					Srcref(observer)

				observer.ReplaceWith(observerDirectReference)
				r.compiler.ReportChangeToEnclosingScope(observerDirectReference)
			}
		}
	}
}

// validateComputedPropertiesReflectionCalls: for any property in the Polymer
// property configuration object with a `computed` key, parses the method call
// and path arguments and reports errors on invalid syntax.
//
// Java: PolymerClassRewriter.validateComputedPropertiesReflectionCalls(PolymerClassDefinition)
func (r *PolymerClassRewriter) validateComputedPropertiesReflectionCalls(cls *PolymerClassDefinition) {
	for _, prop := range cls.props {
		if prop.value.IsObjectLit() {
			computed := NodeUtil_GetFirstPropMatchingKey(prop.value, "computed")
			if !computed.IsNil() && computed.IsStringLit() {
				r.validateReflectedString(computed)
			}
		}
	}
}

// validateComplexObserverReflectionCalls: for any strings returned in the
// array from the Polymer static observers property, parses the method call
// and path arguments and reports errors on invalid syntax.
//
// Java: PolymerClassRewriter.validateComplexObserverReflectionCalls(PolymerClassDefinition)
func (r *PolymerClassRewriter) validateComplexObserverReflectionCalls(cls *PolymerClassDefinition) {
	classMembers := NodeUtil_GetClassMembers(cls.definition)
	getter := NodeUtil_GetFirstGetterMatchingKey(classMembers, "observers")
	if !getter.IsNil() {
		var complexObservers rhino.Node
		for stmt := range ast.MustGetterDef(getter).Function().Body().MustBlock().Statements() {
			if ret, ok := stmt.AsReturn(); ok {
				if val, hasVal := ret.Value(); hasVal && val.Node().IsArrayLit() {
					complexObservers = val.Node()
					break
				}
			}
		}
		if !complexObservers.IsNil() {
			for complexObserver := complexObservers.GetFirstChild(); !complexObserver.IsNil(); {
				next := complexObserver.GetNext()
				if complexObserver.IsStringLit() {
					r.validateReflectedString(complexObserver)
				}
				complexObserver = next
			}
		}
	}
}

// validateReflectedString: given a Polymer method call string such as:
// `"methodName(path.item, other.property.path)"` parses the string into a
// method name and arguments and reports errors on invalid syntax.
//
// Java: PolymerClassRewriter.validateReflectedString(Node)
func (r *PolymerClassRewriter) validateReflectedString(methodSignature rhino.Node) {
	jutil.CheckArgument(methodSignature.IsStringLit())
	methodSignatureString := jutil.StringTrim(methodSignature.GetString())
	openParenIndex := jutil.StringIndexOfChar(methodSignatureString, '(')
	length := jutil.StringLength(methodSignatureString)
	if jutil.StringCharAt(methodSignatureString, length-1) != ')' ||
		openParenIndex < 1 {
		r.compiler.Report(JSError_MakeWithNode(methodSignature, PolymerPassErrors_POLYMER_UNPARSABLE_STRING))
		return
	}

	// Process any parameters in the method call
	if openParenIndex < length-2 {
		methodParamsString := jutil.StringTrim(
			jutil.StringSubstring(methodSignatureString, openParenIndex+1, length-1))

		r.parseMethodParams(methodParamsString, methodSignature)
	}
}

// parseMethodParams parses the parameters string from a complex observer or
// computed property into distinct parameters and reports errors on invalid
// syntax.
//
// Java: PolymerClassRewriter.parseMethodParams(String, Node)
func (r *PolymerClassRewriter) parseMethodParams(methodParameters string, methodSignature rhino.Node) {
	// The delimiters are ASCII, which never occurs inside a multi-byte
	// sequence: scanning the bytes sees the same delimiters as the chars.
	nextDelimeter := byte(',')
	for i := 0; i < len(methodParameters); i++ {
		if methodParameters[i] == nextDelimeter {
			nextDelimeter = ','
		} else {
			if methodParameters[i] == '"' || methodParameters[i] == '\'' {
				nextDelimeter = methodParameters[i]
			}
		}
	}
	if nextDelimeter != ',' {
		r.compiler.Report(JSError_MakeWithNode(methodSignature, PolymerPassErrors_POLYMER_UNPARSABLE_STRING))
	}
}

// Java: PolymerClassRewriter.getInsertionPointForGoogModule(Node)
func polymerClassRewriter_getInsertionPointForGoogModule(moduleBody rhino.Node) rhino.Node {
	insertionPoint := moduleBody.GetFirstChild() // goog.module('ns');
	next := insertionPoint.GetNext()
	for polymerClassRewriter_isGoogRequireExpr(next) ||
		NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(next) ||
		NodeUtil_IsGoogSetTestOnlyCall(next) {
		insertionPoint = next
		next = next.GetNext()
	}
	return insertionPoint
}

// Java: PolymerClassRewriter.isGoogRequireExpr(Node)
func polymerClassRewriter_isGoogRequireExpr(statement rhino.Node) bool {
	if NodeUtil_IsExprCall(statement) &&
		ModuleImportResolver_IsGoogModuleDependencyCall(statement.GetOnlyChild()) {
		// `goog.require('a.b.c');`
		return true
	}
	decl, ok := ast.AsNameDeclaration(statement)
	if !ok {
		return false
	}
	first, ok := decl.FirstDeclarator()
	if !ok {
		return false
	}
	var rhs rhino.Node
	if name, ok := first.AsName(); ok {
		// `const c = goog.require('a.b.c');`
		if init, hasInit := name.Init(); hasInit {
			rhs = init.Node()
		}
	} else if dlhs, ok := first.AsDestructuringLhs(); ok {
		// `const {D} = goog.require('a.b.c');`
		if init, hasInit := dlhs.Init(); hasInit {
			rhs = init.Node()
		}
	}
	return ModuleImportResolver_IsGoogModuleDependencyCall(rhs)
}
