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
	"strings"

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

// PolymerClassRewriter rewrites a given call to Polymer({}) or class extending
// Polymer.Element to declarations and assignments understood by the compiler.
//
// Java: final class PolymerClassRewriter
type PolymerClassRewriter struct {
	compiler                   AbstractCompiler
	externsInsertionRef        rhino.Node
	propertySinkExternInjected bool
}

// Java: PolymerClassRewriter(AbstractCompiler)
func NewPolymerClassRewriter(compiler AbstractCompiler) *PolymerClassRewriter {
	return &PolymerClassRewriter{
		compiler:            compiler,
		externsInsertionRef: compiler.GetSynthesizedExternsInput().GetAstRoot(compiler),
	}
}

func polymerIsIIFE(n rhino.Node) bool {
	call, ok := ast.AsCall(n)
	return ok && call.Callee().Node().IsFunction()
}

func polymerIsFunctionArgInGoogLoadModule(n rhino.Node) bool {
	if !n.IsFunction() {
		return false
	}
	parent := n.GetParent()
	return !parent.IsNil() && NodeUtil_IsBundledGoogModuleCall(parent)
}

func (r *PolymerClassRewriter) insertGeneratedDeclarationCodeToGlobalScope(
	enclosingNode rhino.Node, declarationCode rhino.Node) {
	switch enclosingNode.GetToken() {
	case rhino.Token_MODULE_BODY:
		insertionPoint := r.getNodeForInsertion(enclosingNode.GetParent())
		insertionPoint.AddChildToFront(declarationCode)
		r.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(insertionPoint))
	case rhino.Token_SCRIPT:
		enclosingNode.AddChildToFront(declarationCode)
		r.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(enclosingNode))
	case rhino.Token_CALL:
		jutil.CheckState(polymerIsIIFE(enclosingNode))
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
	case rhino.Token_FUNCTION:
		jutil.CheckState(polymerIsFunctionArgInGoogLoadModule(enclosingNode))
		enclosingScript := NodeUtil_GetEnclosingScript(enclosingNode)
		insertionPoint := r.getNodeForInsertion(enclosingScript)
		insertionPoint.AddChildToFront(declarationCode)
		r.compiler.ReportChangeToChangeScope(insertionPoint)
	default:
		panic(jutil.NewRuntimeException("Enclosing node for Polymer is incorrect"))
	}
}

func (r *PolymerClassRewriter) getNodeForInsertion(enclosingScript rhino.Node) rhino.Node {
	if NodeUtil_IsFromTypeSummary(enclosingScript) {
		return r.compiler.GetSynthesizedTypeSummaryInput().GetAstRoot(r.compiler)
	}
	return r.compiler.GetNodeForCodeInsertion(nil)
}

func (r *PolymerClassRewriter) insertGeneratedPropsAndBehaviorCode(enclosingNode rhino.Node, statements rhino.Node) {
	switch enclosingNode.GetToken() {
	case rhino.Token_MODULE_BODY:
		if enclosingNode.GetParent().GetBooleanProp(rhino.Node_GOOG_MODULE) {
			// The goog.module('ns'); call must remain the first statement in the module.
			insertionPoint := getInsertionPointForGoogModule(enclosingNode)
			enclosingNode.AddChildrenAfter(statements, insertionPoint)
		} else {
			enclosingNode.AddChildrenToFront(statements)
		}
	case rhino.Token_SCRIPT:
		enclosingNode.AddChildrenToBack(statements)
		r.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(enclosingNode))
	case rhino.Token_CALL:
		jutil.CheckState(polymerIsIIFE(enclosingNode))
		functionNode := ast.MustFunction(ast.MustCall(enclosingNode).Callee().Node())
		functionNode.Body().Node().AddChildrenToFront(statements)
	case rhino.Token_FUNCTION:
		jutil.CheckState(polymerIsFunctionArgInGoogLoadModule(enclosingNode))
		functionBlock := ast.MustFunction(enclosingNode).Body().Node()
		insertionPoint := getInsertionPointForGoogModule(functionBlock)
		if !insertionPoint.IsNil() {
			functionBlock.AddChildrenAfter(statements, insertionPoint)
		}
	}
}

// RewritePolymerCall rewrites a given call to Polymer({}) to declarations and
// assignments understood by the compiler.
//
// Java: PolymerClassRewriter.rewritePolymerCall(PolymerClassDefinition, NodeTraversal)
func (r *PolymerClassRewriter) RewritePolymerCall(cls *PolymerClassDefinition, traversal *NodeTraversal) {
	callParent := cls.definition.GetParent()
	exprRoot := callParent.GetParent()
	if callParent.IsExprResult() {
		exprRoot = callParent
	}
	jutil.CheckStateMsg1(NodeUtil_IsStatementParent(exprRoot.GetParent()), "%s", exprRoot.GetParent())
	objLit := rhino.CheckNotNull(cls.descriptor)
	targetQName := cls.targetQualifiedName()

	// Add @lends to the object literal.
	objLitDoc := rhino.JSDocInfo_Builder_().ParseDocumentation()
	jsTypeExpression := rhino.NewJSTypeExpression(
		rhino.IR_String(targetQName+".prototype").Srcref(exprRoot),
		polymerSourceFileName(exprRoot))
	objLitDoc.RecordLends(jsTypeExpression)
	objLit.SetJSDocInfo(objLitDoc.Build())

	r.addTypesToFunctions(objLit, targetQName, cls.defType)
	PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(objLit, r.compiler)
	PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(objLit, r.compiler)
	PolymerPassStaticUtils_ProtectObserverAndPropertyFunctionKeys(objLit)

	for _, prop := range cls.props {
		if prop.value.IsObjectLit() {
			PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(prop.value, r.compiler)
		}
	}

	propsAndBehaviorBlock := rhino.IR_Block()
	constructorDoc := r.getConstructorDoc(cls)

	// Remove the original constructor JSDoc from the objlit.
	if ctorKey := cls.constructor.value.GetParent(); !ctorKey.IsNil() {
		ctorKey.SetJSDocInfo(nil)
	}

	// Check for a conflicting definition of PolymerElement.
	if !traversal.InGlobalScope() {
		polymerElement := traversal.GetScope().GetVar("PolymerElement")
		if polymerElement != nil && !polymerElement.GetScope().(*Scope).IsGlobal() {
			nameNode := polymerElement.GetNameNode()
			r.compiler.Report(
				JSError_MakeWithNode(
					cls.constructor.value,
					PolymerClassRewriter_POLYMER_ELEMENT_CONFLICT,
					polymerSourceFileName(nameNode),
					strconv.Itoa(nameNode.GetLineno()),
					strconv.Itoa(nameNode.GetCharno())))
		}
	}

	declarationCode := r.generateDeclarationCode(exprRoot, cls, constructorDoc, traversal)
	basePath := targetQName + ".prototype."
	r.appendBehaviorPropertiesToBlock(cls, propsAndBehaviorBlock, basePath, false)
	r.appendPropertiesToBlock(cls.props, propsAndBehaviorBlock, basePath, false)
	r.appendBehaviorMembersToBlock(cls, propsAndBehaviorBlock)
	readOnlyPropsAll := r.parseReadOnlyProperties(cls, propsAndBehaviorBlock)
	r.createExportsAndExterns(cls, readOnlyPropsAll)
	r.removePropertyDocs(objLit, cls.defType)

	propsAndBehaviorCode := propsAndBehaviorBlock.RemoveChildren()
	parent := exprRoot.GetParent()

	if !traversal.InGlobalScope() && cls.hasGeneratedLhs && !cls.target.IsGetProp() {
		enclosingNode := NodeUtil_GetEnclosingNode(
			parent,
			func(node rhino.Node) bool {
				return node.IsScript() ||
					node.IsModuleBody() ||
					polymerIsIIFE(node) ||
					polymerIsFunctionArgInGoogLoadModule(node)
			})
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
		assignExpr := polymerVarToAssign(exprRoot)
		exprRoot.ReplaceWith(assignExpr)
		r.compiler.ReportChangeToEnclosingScope(assignExpr)
	}

	if !cls.descriptor.IsNil() {
		props := NodeUtil_GetFirstPropMatchingKey(cls.descriptor, "properties")
		if !props.IsNil() && props.IsObjectLit() {
			r.addPropertiesConfigObjectReflection(cls, props)
		}
	}
}

// RewritePolymerClassDeclaration rewrites a class extending Polymer.Element to
// declarations and assignments understood by the compiler.
//
// Java: PolymerClassRewriter.rewritePolymerClassDeclaration(Node, NodeTraversal, PolymerClassDefinition)
func (r *PolymerClassRewriter) RewritePolymerClassDeclaration(
	clazz rhino.Node, traversal *NodeTraversal, cls *PolymerClassDefinition) {
	targetQName := cls.targetQualifiedName()

	if !cls.descriptor.IsNil() {
		r.addTypesToFunctions(cls.descriptor, targetQName, cls.defType)
	}
	PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(
		NodeUtil_GetClassMembers(clazz), r.compiler)

	for _, prop := range cls.props {
		if prop.value.IsObjectLit() {
			PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(prop.value, r.compiler)
		}
	}

	block := rhino.IR_Block()
	basePath := targetQName + ".prototype."
	r.appendBehaviorPropertiesToBlock(cls, block, basePath, false)
	r.appendPropertiesToBlock(cls.props, block, basePath, false)

	readOnlyProps := r.parseReadOnlyProperties(cls, block)
	r.createExportsAndExterns(cls, readOnlyProps)

	jsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(clazz)
	classInfo := rhino.JSDocInfo_Builder_MaybeCopyFrom(jsDocInfoNode.GetJSDocInfo())
	interfaceName := cls.GetInterfaceName(r.compiler.GetUniqueIdSupplier())
	interfaceType := rhino.NewJSTypeExpression(
		rhino.NewNodeWithChild(rhino.Token_BANG, rhino.IR_String(interfaceName)).SrcrefTree(jsDocInfoNode),
		polymerClassRewriter_VIRTUAL_FILE)
	classInfo.RecordImplementedInterface(interfaceType)
	jsDocInfoNode.SetJSDocInfo(classInfo.Build())

	insertAfterReference := NodeUtil_GetEnclosingStatement(clazz)
	if block.HasChildren() {
		r.removePropertyDocs(cls.descriptor, cls.defType)
		insertAfterReference.GetParent().AddChildrenAfter(block.RemoveChildren(), insertAfterReference)
		r.compiler.ReportChangeToEnclosingScope(insertAfterReference)
	}

	addReturnTypeIfMissing(cls, "is", rhino.NewJSTypeExpression(rhino.IR_String("string"), polymerClassRewriter_VIRTUAL_FILE))

	type_ := rhino.NewNode(rhino.Token_BANG)
	array := rhino.IR_String("Array")
	type_.AddChildToBack(array)
	array.AddChildToBack(rhino.NewNodeWithChild(rhino.Token_BLOCK, rhino.IR_String("string")))
	addReturnTypeIfMissing(cls, "observers", rhino.NewJSTypeExpression(type_, polymerClassRewriter_VIRTUAL_FILE))
	addReturnTypeIfMissing(
		cls,
		"properties",
		rhino.NewJSTypeExpression(rhino.IR_String(polymerClassRewriter_POLYMER_ELEMENT_PROP_CONFIG), polymerClassRewriter_VIRTUAL_FILE))

	if !cls.descriptor.IsNil() {
		r.convertSimpleObserverStringsToReferences(cls)
		r.addPropertiesConfigObjectReflection(cls, cls.descriptor)
		r.validateComputedPropertiesReflectionCalls(cls)
		r.validateComplexObserverReflectionCalls(cls)
	}
}

func addReturnTypeIfMissing(
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
					rhino.NewNodeWithChild(rhino.Token_BANG, rhino.IR_String(thisType)).SrcrefTree(keyNode),
					polymerClassRewriter_VIRTUAL_FILE))
			keyNode.SetJSDocInfo(fnDoc.Build())
		}
	}

	for _, property := range PolymerPassStaticUtils_ExtractProperties(
		objLit, defType, r.compiler, rhino.NilNode) {
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
				rhino.NewNodeWithChild(rhino.Token_BANG, rhino.IR_String(thisType)).SrcrefTree(defaultValueKey),
				polymerClassRewriter_VIRTUAL_FILE))
		fnDoc.RecordReturnType(PolymerPassStaticUtils_GetTypeFromProperty(property, r.compiler))
		defaultValueKey.SetJSDocInfo(fnDoc.Build())
	}
}

func (r *PolymerClassRewriter) parseReadOnlyProperties(
	cls *PolymerClassDefinition, block rhino.Node) []*PolymerPass_MemberDefinition {
	qualifiedPath := cls.targetQualifiedName() + ".prototype."
	var readOnlyProps []*PolymerPass_MemberDefinition

	addReadOnlyIfApplicable := func(prop *PolymerPass_MemberDefinition) {
		if !prop.value.IsObjectLit() {
			return
		}
		readOnlyValue := NodeUtil_GetFirstPropMatchingKey(prop.value, "readOnly")
		if !readOnlyValue.IsNil() && readOnlyValue.IsTrue() {
			setter := r.makeReadOnlySetter(prop, qualifiedPath)
			setter.SrcrefTreeIfMissing(prop.name)
			block.AddChildToBack(setter)
			readOnlyProps = append(readOnlyProps, prop)
		}
	}

	for _, prop := range cls.props {
		addReadOnlyIfApplicable(prop)
	}
	for _, bp := range cls.behaviorProps {
		addReadOnlyIfApplicable(bp.prop)
	}
	return readOnlyProps
}

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

func (r *PolymerClassRewriter) generateDeclarationCode(
	exprRoot rhino.Node,
	cls *PolymerClassDefinition,
	constructorDoc *rhino.JSDocInfo_Builder,
	traversal *NodeTraversal) rhino.Node {
	if cls.target.IsGetProp() {
		assign := rhino.IR_Assign(cls.target.CloneTree(), cls.constructor.value.CloneTree())
		NodeUtil_MarkNewScopesChanged(assign, r.compiler)
		assign.SetJSDocInfo(constructorDoc.Build())
		exprResult := rhino.IR_ExprResult(assign)
		exprResult.SrcrefTreeIfMissing(cls.target)
		return exprResult
	}
	var_ := rhino.IR_VarWithValue(cls.target.CloneTree(), cls.constructor.value.CloneTree())
	NodeUtil_MarkNewScopesChanged(var_, r.compiler)
	var_.SrcrefTreeIfMissing(exprRoot)
	var_.SetJSDocInfo(constructorDoc.Build())
	name := cls.target.GetString()
	if existingVar := traversal.GetScope().GetSlot(name); existingVar != nil && cls.hasGeneratedLhs {
		r.compiler.Report(JSError_MakeWithNode(cls.constructor.value, PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT, name))
	}
	return var_
}

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

func (r *PolymerClassRewriter) getPropertyNode(prop *PolymerPass_MemberDefinition, basePath string) rhino.Node {
	if prop.name.IsQuotedStringKey() {
		return rhino.NilNode
	}
	propertyNode := rhino.IR_ExprResult(NodeUtil_NewQName(r.compiler, basePath+prop.name.GetString()))
	propertyNode.SrcrefTreeIfMissing(prop.name)
	return propertyNode
}

func (r *PolymerClassRewriter) appendBehaviorPropertiesToBlock(
	cls *PolymerClassDefinition, block rhino.Node, basePath string, isExternsBlock bool) {
	if len(cls.behaviors) == 0 || len(cls.behaviorProps) == 0 {
		return
	}
	for _, bp := range cls.behaviorProps {
		prop := bp.prop
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
			infoBuilder := r.getJSDocInfoBuilderForBehavior(bp.behavior, prop)
			infoBuilder.RecordType(propType)
			info = infoBuilder.Build()
		}
		ast.MustExprStmt(propertyNode).Expr().Node().SetJSDocInfo(info)
		block.AddChildToBack(propertyNode)
	}
}

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

func (r *PolymerClassRewriter) createVarsInExternsBlock(
	block rhino.Node,
	propertyNames *jutil.LinkedHashSet[string],
	propType *rhino.JSTypeExpression,
	prop *PolymerPass_MemberDefinition) {
	for propName := range propertyNames.All() {
		varName := "PolymerDummyVar" + r.compiler.GetUniqueNameIdSupplier()()
		n := rhino.Node_NewStringWithToken(rhino.Token_NAME, varName)
		var_ := rhino.NewNodeWithChild(rhino.Token_VAR, n)

		newType := createNewTypeExpressionForExtern(propName, propType.GetSourceName(), prop)
		oldInfo := rhino.JSDocInfo_Builder_MaybeCopyFrom(prop.info).Build()
		newInfo := rhino.JSDocInfo_Builder_CopyFromWithNewType(oldInfo, newType)
		var_.SetJSDocInfo(newInfo.Build())
		block.AddChildToBack(var_)
	}
}

func createNewTypeExpressionForExtern(
	propName string, sourceName string, prop *PolymerPass_MemberDefinition) *rhino.JSTypeExpression {
	propertyName := rhino.Node_NewStringWithToken(rhino.Token_STRING_KEY, propName)
	colon := rhino.NewNodeWithChild(rhino.Token_COLON, propertyName)
	colon.AddChildToBack(rhino.NewNode(rhino.Token_QMARK))
	leftBracket := rhino.NewNodeWithChild(rhino.Token_LB, colon)
	leftCurly := rhino.NewNodeWithChild(rhino.Token_LC, leftBracket).SrcrefTreeIfMissing(prop.name)
	return rhino.NewJSTypeExpression(leftCurly, sourceName)
}

func (r *PolymerClassRewriter) removePropertyDocs(
	objLit rhino.Node, defType PolymerClassDefinition_DefinitionType) {
	for _, prop := range PolymerPassStaticUtils_ExtractProperties(
		objLit, defType, r.compiler, rhino.NilNode) {
		prop.name.SetJSDocInfo(nil)
	}
}

func (r *PolymerClassRewriter) getJSDocInfoBuilderForBehavior(
	behavior *PolymerBehaviorExtractor_BehaviorDefinition, behaviorFunctionOrProp *PolymerPass_MemberDefinition) *rhino.JSDocInfo_Builder {
	if !behavior.isGlobalDeclaration &&
		behaviorFunctionOrProp.info != nil &&
		behaviorFunctionOrProp.info.ContainsTypeDeclaration() &&
		!behavior.behaviorModule.IsNil() {
		return rhino.JSDocInfo_Builder_MaybeCopyFromAndReplaceNames(
			behaviorFunctionOrProp.info, behavior.GetModuleLocalNames(r.compiler))
	}
	return rhino.JSDocInfo_Builder_MaybeCopyFrom(behaviorFunctionOrProp.info)
}

func (r *PolymerClassRewriter) appendBehaviorMembersToBlock(cls *PolymerClassDefinition, block rhino.Node) {
	if len(cls.behaviors) == 0 {
		return
	}
	qualifiedPath := cls.targetQualifiedName() + ".prototype."
	nameToExprResult := make(map[string]rhino.Node)

	for _, behavior := range cls.behaviors {
		for _, behaviorFunction := range behavior.functionsToCopy {
			fnName := behaviorFunction.name.GetString()
			if !NodeUtil_GetFirstPropMatchingKey(cls.descriptor, fnName).IsNil() {
				continue
			}

			// Avoid copying over the same function twice; the last definition wins.
			if prev, exists := nameToExprResult[fnName]; exists {
				prev.Detach()
			}

			fnValue := behaviorFunction.value.CloneTree()
			NodeUtil_MarkNewScopesChanged(fnValue, r.compiler)
			exprResult := rhino.IR_ExprResult(rhino.IR_Assign(NodeUtil_NewQName(r.compiler, qualifiedPath+fnName), fnValue))
			exprResult.SrcrefTreeIfMissing(behaviorFunction.name)

			info := r.getJSDocInfoBuilderForBehavior(behavior, behaviorFunction)
			if behaviorFunction.info != nil &&
				behaviorFunction.info.GetVisibility() == rhino.JSDocInfo_Visibility_PROTECTED {
				info.OverwriteVisibility(rhino.JSDocInfo_Visibility_PUBLIC)
			}

			if !behavior.isGlobalDeclaration {
				fn := ast.MustFunction(fnValue)
				body := fn.Body().Node()
				if fn.IsArrow() && !body.IsBlock() {
					body.ReplaceWith(NodeUtil_NewUndefinedNode(body))
				} else {
					body.RemoveChildren()
				}
				paramIndex := 0
				for param := range ast.ReadAhead(fn.Params().FirstParam()) {
					polymerMakeParamSafe(param.Node(), paramIndex)
					paramIndex++
				}
				if behaviorFunction.info == nil || behaviorFunction.info.GetReturnType() == nil {
					info.RecordReturnType(
						rhino.NewJSTypeExpression(
							rhino.NewNode(rhino.Token_QMARK).Srcref(fnValue),
							polymerSourceFileName(fnValue)))
				}
			}

			exprResult.GetFirstChild().SetJSDocInfo(info.Build())
			block.AddChildToBack(exprResult)
			nameToExprResult[fnName] = exprResult
		}

		for _, behaviorProp := range behavior.nonPropertyMembersToCopy {
			propName := behaviorProp.name.GetString()
			if prev, exists := nameToExprResult[propName]; exists {
				prev.Detach()
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
			nameToExprResult[propName] = exprResult
		}
	}
}

func polymerMakeParamSafe(param rhino.Node, index int) {
	if rest, ok := ast.AsIterRest(param); ok {
		param = rest.Target().Node()
	} else if defVal, ok := ast.AsDefaultValue(param); ok {
		target := defVal.Target().Node()
		defVal.Value().Node().ReplaceWith(NodeUtil_NewUndefinedNode(param))
		param = target
	}
	if param.IsDestructuringPattern() {
		param.ReplaceWith(rhino.IR_Name("param$polymer$" + strconv.Itoa(index)).Srcref(param))
	}
}

func polymerSourceFileName(n rhino.Node) string {
	if s := n.GetSourceFileName(); s != nil {
		return *s
	}
	return "null"
}

func polymerReadOnlySetterName(propName string) string {
	if len(propName) > 0 && propName[0] < 0x80 {
		ch := propName[0]
		if ch >= 'a' && ch <= 'z' {
			ch -= 'a' - 'A'
		}
		return "_set" + string(ch) + propName[1:]
	}
	return "_set" + stringToUpperCaseRoot(jutil.StringSubstring(propName, 0, 1)) + jutil.StringSubstringFrom(propName, 1)
}

func (r *PolymerClassRewriter) makeReadOnlySetter(prop *PolymerPass_MemberDefinition, qualifiedPath string) rhino.Node {
	propName := prop.name.GetString()
	setterName := polymerReadOnlySetterName(propName)
	fnNode := rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(rhino.IR_Name(propName)), rhino.IR_Block())
	r.compiler.ReportChangeToChangeScope(fnNode)
	exprResNode := rhino.IR_ExprResult(rhino.IR_Assign(NodeUtil_NewQName(r.compiler, qualifiedPath+setterName), fnNode))

	info := rhino.JSDocInfo_Builder_().ParseDocumentation()
	info.RecordOverride()
	propType := PolymerPassStaticUtils_GetTypeFromProperty(prop, r.compiler)
	info.RecordParameter(propName, propType)
	exprResNode.GetFirstChild().SetJSDocInfo(info.Build())
	return exprResNode
}

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

	r.appendBehaviorPropertiesToBlock(cls, block, interfaceBasePath, true)
	r.appendPropertiesToBlock(cls.props, block, interfaceBasePath, true)

	// Deduplicate behavior and element methods by name in first-seen key order
	// with last-seen definition winning (matching Java LinkedHashMap.put).
	var methodOrder []string
	var uniqueMethods map[string]*PolymerPass_MemberDefinition
	putMethod := func(method *PolymerPass_MemberDefinition) {
		if uniqueMethods == nil {
			uniqueMethods = make(map[string]*PolymerPass_MemberDefinition)
		}
		name := method.name.GetString()
		if _, exists := uniqueMethods[name]; !exists {
			methodOrder = append(methodOrder, name)
		}
		uniqueMethods[name] = method
	}
	for _, behavior := range cls.behaviors {
		for _, method := range behavior.functionsToCopy {
			putMethod(method)
		}
	}
	for _, method := range cls.methods {
		putMethod(method)
	}
	for _, name := range methodOrder {
		r.addMethodToObjectExternsUsingExportAnnotation(cls, uniqueMethods[name])
	}

	for _, prop := range readOnlyProps {
		propName := prop.name.GetString()
		setterName := polymerReadOnlySetterName(propName)
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

func (r *PolymerClassRewriter) addMethodToObjectExternsUsingExportAnnotation(
	cls *PolymerClassDefinition, method *PolymerPass_MemberDefinition) {
	getprop := NodeUtil_NewQName(
		r.compiler, cls.targetQualifiedName()+".prototype."+method.name.GetString())
	info := rhino.JSDocInfo_Builder_().ParseDocumentation()
	if method.info != nil {
		info.RecordVisibility(method.info.GetVisibility())
	}
	info.RecordExport()
	getprop.SetJSDocInfo(info.Build())
	expression := rhino.IR_ExprResult(getprop).SrcrefTreeIfMissing(method.name)
	insertAfter := cls.definition
	for !NodeUtil_IsStatementBlock(insertAfter.GetParent()) {
		insertAfter = insertAfter.GetParent()
	}
	expression.InsertAfter(insertAfter)
	r.compiler.ReportChangeToEnclosingScope(expression)
}

func polymerVarToAssign(var_ rhino.Node) rhino.Node {
	lhs := ast.MustNameDeclaration(var_).MustFirstDeclarator().Node()
	assign := rhino.IR_Assign(lhs.CloneNode(), lhs.RemoveFirstChild())
	return rhino.IR_ExprResult(assign).SrcrefTreeIfMissing(var_)
}

func (r *PolymerClassRewriter) convertSimpleObserverStringsToReferences(cls *PolymerClassDefinition) {
	for _, prop := range cls.props {
		if prop.value.IsObjectLit() {
			observer := NodeUtil_GetFirstPropMatchingKey(prop.value, "observer")
			if !observer.IsNil() && observer.IsStringLit() {
				observerDirectReference := rhino.IR_GetpropWithMoreProps(
					cls.target.CloneTree(), "prototype", observer.GetString()).Srcref(observer)
				observer.ReplaceWith(observerDirectReference)
				r.compiler.ReportChangeToEnclosingScope(observerDirectReference)
			}
		}
	}
}

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

func (r *PolymerClassRewriter) validateComplexObserverReflectionCalls(cls *PolymerClassDefinition) {
	classMembers := NodeUtil_GetClassMembers(cls.definition)
	getter := NodeUtil_GetFirstGetterMatchingKey(classMembers, "observers")
	if getter.IsNil() {
		return
	}
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

func (r *PolymerClassRewriter) validateReflectedString(methodSignature rhino.Node) {
	jutil.CheckArgument(methodSignature.IsStringLit())
	sig := jutil.StringTrim(methodSignature.GetString())
	openParenIndex := strings.IndexByte(sig, '(')
	if !strings.HasSuffix(sig, ")") || openParenIndex < 1 {
		r.compiler.Report(JSError_MakeWithNode(methodSignature, PolymerPassErrors_POLYMER_UNPARSABLE_STRING))
		return
	}
	if openParenIndex+1 < len(sig)-1 {
		params := jutil.StringTrim(sig[openParenIndex+1 : len(sig)-1])
		r.parseMethodParams(params, methodSignature)
	}
}

func (r *PolymerClassRewriter) parseMethodParams(methodParameters string, methodSignature rhino.Node) {
	nextDelimiter := byte(',')
	for i := 0; i < len(methodParameters); i++ {
		ch := methodParameters[i]
		if ch == nextDelimiter {
			nextDelimiter = ','
		} else if ch == '"' || ch == '\'' {
			nextDelimiter = ch
		}
	}
	if nextDelimiter != ',' {
		r.compiler.Report(JSError_MakeWithNode(methodSignature, PolymerPassErrors_POLYMER_UNPARSABLE_STRING))
	}
}

func getInsertionPointForGoogModule(moduleBody rhino.Node) rhino.Node {
	insertionPoint := moduleBody.GetFirstChild() // goog.module('ns');
	next := insertionPoint.GetNext()
	for isGoogRequireExpr(next) ||
		NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(next) ||
		NodeUtil_IsGoogSetTestOnlyCall(next) {
		insertionPoint = next
		next = next.GetNext()
	}
	return insertionPoint
}

func isGoogRequireExpr(statement rhino.Node) bool {
	if NodeUtil_IsExprCall(statement) &&
		ModuleImportResolver_IsGoogModuleDependencyCall(statement.GetOnlyChild()) {
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
		if init, hasInit := name.Init(); hasInit {
			rhs = init.Node()
		}
	} else if dlhs, ok := first.AsDestructuringLhs(); ok {
		if init, hasInit := dlhs.Init(); hasInit {
			rhs = init.Node()
		}
	}
	return ModuleImportResolver_IsGoogModuleDependencyCall(rhs)
}
