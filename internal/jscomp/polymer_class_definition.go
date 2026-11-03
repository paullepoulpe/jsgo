// Port of com/google/javascript/jscomp/PolymerClassDefinition.java @ closure-compiler v20261028

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type polymerBehaviorProp struct {
	prop     *PolymerPass_MemberDefinition
	behavior *PolymerBehaviorExtractor_BehaviorDefinition
}

// PolymerClassDefinition is a parsed Polymer class (element) definition.
// Includes convenient fields for rewriting the class.
//
// Java: final class PolymerClassDefinition
type PolymerClassDefinition struct {
	defType           PolymerClassDefinition_DefinitionType
	definition        rhino.Node
	input             *CompilerInput
	target            rhino.Node
	hasGeneratedLhs   bool
	descriptor        rhino.Node
	constructor       *PolymerPass_MemberDefinition
	nativeBaseElement *string
	props             []*PolymerPass_MemberDefinition
	behaviorProps     []polymerBehaviorProp
	methods           []*PolymerPass_MemberDefinition
	behaviors         []*PolymerBehaviorExtractor_BehaviorDefinition
	features          parser.FeatureSet
	interfaceName     string
}

// Java: static enum PolymerClassDefinition.DefinitionType
type PolymerClassDefinition_DefinitionType int32

const (
	PolymerClassDefinition_DefinitionType_ObjectLiteral PolymerClassDefinition_DefinitionType = iota
	PolymerClassDefinition_DefinitionType_ES6Class
)

var polymerClassDefinition_DefinitionType_names = [...]string{"ObjectLiteral", "ES6Class"}

// String is DefinitionType.name().
func (t PolymerClassDefinition_DefinitionType) String() string {
	return polymerClassDefinition_DefinitionType_names[t]
}

// NewPolymerClassDefinition constructs a PolymerClassDefinition.
func NewPolymerClassDefinition(
	defType PolymerClassDefinition_DefinitionType,
	definition rhino.Node,
	target rhino.Node,
	hasGeneratedLhs bool,
	descriptor rhino.Node,
	classInfo *rhino.JSDocInfo,
	constructor *PolymerPass_MemberDefinition,
	nativeBaseElement *string,
	props []*PolymerPass_MemberDefinition,
	behaviorProps []polymerBehaviorProp,
	methods []*PolymerPass_MemberDefinition,
	behaviors []*PolymerBehaviorExtractor_BehaviorDefinition,
	features parser.FeatureSet,
	input *CompilerInput) *PolymerClassDefinition {
	jutil.CheckState(descriptor.IsNil() || descriptor.IsObjectLit())
	return &PolymerClassDefinition{
		defType:           defType,
		definition:        definition,
		target:            target,
		hasGeneratedLhs:   hasGeneratedLhs,
		descriptor:        descriptor,
		constructor:       constructor,
		nativeBaseElement: nativeBaseElement,
		props:             props,
		behaviorProps:     behaviorProps,
		methods:           methods,
		behaviors:         behaviors,
		features:          features,
		input:             input,
	}
}

func (d *PolymerClassDefinition) targetQualifiedName() string {
	qname, ok := d.target.GetQualifiedNameString()
	jutil.CheckState(ok)
	return qname
}

// PolymerClassDefinition_ExtractFromCallNode validates the class definition
// and if valid, destructively extracts the class definition from the AST.
//
// Java: PolymerClassDefinition.extractFromCallNode(Node, AbstractCompiler, ModuleMetadata, PolymerBehaviorExtractor)
func PolymerClassDefinition_ExtractFromCallNode(
	callNode rhino.Node,
	compiler AbstractCompiler,
	moduleMetadata *ModuleMetadataMap_ModuleMetadata,
	behaviorExtractor *PolymerBehaviorExtractor) *PolymerClassDefinition {
	call, ok := ast.AsCall(callNode)
	var descriptor rhino.Node
	if ok {
		if firstArg, hasArg := call.FirstArg(); hasArg {
			descriptor = firstArg.Node()
		}
	}
	if !descriptor.IsObjectLit() {
		compiler.Report(JSError_MakeWithNode(callNode, PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID))
		return nil
	}

	if call.NumArgs() != 1 {
		compiler.Report(JSError_MakeWithNode(callNode, PolymerPassErrors_POLYMER_UNEXPECTED_PARAMS))
		return nil
	}

	elName := NodeUtil_GetFirstPropMatchingKey(descriptor, "is")
	if elName.IsNil() {
		compiler.Report(JSError_MakeWithNode(callNode, PolymerPassErrors_POLYMER_MISSING_IS))
		return nil
	}
	enclosingModule := NodeUtil_GetEnclosingModuleIfPresent(callNode)

	hasGeneratedLhs := false
	var target rhino.Node
	callParent := callNode.GetParent()
	if NodeUtil_IsNameDeclaration(callParent.GetParent()) {
		target = rhino.IR_Name(callParent.GetString())
	} else if assign, ok := ast.AsAssign(callParent); ok {
		if isPolymerGoogModuleExports(assign) {
			target = createDummyGoogModuleExportsTarget(compiler, callNode, assign)
		} else {
			target = assign.Target().Node().CloneTree()
		}
	} else {
		var elNameStringBase string
		if qname, ok := elName.GetQualifiedNameString(); ok {
			elNameStringBase = strings.ReplaceAll(qname, ".", "$")
		} else {
			elNameStringBase = elName.GetString()
		}
		elNameString := caseFormatLowerHyphenToUpperCamel(elNameStringBase) + "Element"
		target = rhino.IR_Name(elNameString)
		hasGeneratedLhs = true
	}

	classInfo := NodeUtil_GetBestJSDocInfo(target)

	var ctorInfo *rhino.JSDocInfo
	constructor := NodeUtil_GetFirstPropMatchingKey(descriptor, "factoryImpl")
	if constructor.IsNil() {
		constructor = NodeUtil_EmptyFunction()
		compiler.ReportChangeToChangeScope(constructor)
		constructor.SrcrefTree(callNode)
	} else {
		ctorInfo = NodeUtil_GetBestJSDocInfo(constructor)
	}

	baseClass := NodeUtil_GetFirstPropMatchingKey(descriptor, "extends")
	var nativeBaseElement *string
	if !baseClass.IsNil() {
		s := baseClass.GetString()
		nativeBaseElement = &s
	}

	behaviorArray := NodeUtil_GetFirstPropMatchingKey(descriptor, "behaviors")
	behaviors := behaviorExtractor.ExtractBehaviors(behaviorArray, moduleMetadata)
	properties := overwriteMembersIfPresent(
		[]*PolymerPass_MemberDefinition{},
		PolymerPassStaticUtils_ExtractProperties(
			descriptor, PolymerClassDefinition_DefinitionType_ObjectLiteral, compiler, rhino.NilNode))

	// Collect behavior properties in first-seen order, skipping duplicates and
	// any property already declared directly on the element.
	var behaviorProps []polymerBehaviorProp
	if len(behaviors) > 0 {
		seenPropNames := make(map[string]struct{}, len(properties))
		for _, p := range properties {
			seenPropNames[p.name.GetString()] = struct{}{}
		}
		for _, behavior := range behaviors {
			for _, prop := range behavior.props {
				propName := prop.name.GetString()
				if _, exists := seenPropNames[propName]; !exists {
					seenPropNames[propName] = struct{}{}
					behaviorProps = append(behaviorProps, polymerBehaviorProp{
						prop:     prop,
						behavior: behavior,
					})
				}
			}
		}
	}

	var newFeatures parser.FeatureSet
	for _, behavior := range behaviors {
		if !behavior.features.IsZero() {
			newFeatures = newFeatures.Union(behavior.features)
		}
	}

	methods := []*PolymerPass_MemberDefinition{}
	for member := range ast.MustObjectLit(descriptor).Members() {
		if fnDef, ok := member.AsMemberFunctionDef(); ok {
			methods = append(methods,
				NewPolymerPass_MemberDefinition(
					NodeUtil_GetBestJSDocInfo(fnDef.Node()),
					fnDef.Node(),
					fnDef.Function().Node(),
					enclosingModule))
		} else if sk, ok := member.AsStringKey(); ok && sk.Value().Node().IsFunction() {
			methods = append(methods,
				NewPolymerPass_MemberDefinition(
					NodeUtil_GetBestJSDocInfo(sk.Node()),
					sk.Node(),
					sk.Value().Node(),
					enclosingModule))
		}
	}
	input := compiler.GetInput(NodeUtil_GetEnclosingScript(callNode).GetInputId())

	return NewPolymerClassDefinition(
		PolymerClassDefinition_DefinitionType_ObjectLiteral,
		callNode,
		target,
		hasGeneratedLhs,
		descriptor,
		classInfo,
		NewPolymerPass_MemberDefinition(ctorInfo, rhino.NilNode, constructor, enclosingModule),
		nativeBaseElement,
		properties,
		behaviorProps,
		methods,
		behaviors,
		newFeatures,
		input)
}

func isPolymerGoogModuleExports(assign ast.Assign) bool {
	assignNode := assign.Node()
	if !assignNode.GetParent().IsExprResult() || !assign.Target().Node().MatchesName("exports") {
		return false
	}
	containingBlock := assignNode.GetGrandparent()
	return (containingBlock.IsModuleBody() &&
		containingBlock.GetParent().GetBooleanProp(rhino.Node_GOOG_MODULE)) ||
		NodeUtil_IsBundledGoogModuleScopeRoot(containingBlock)
}

func createDummyGoogModuleExportsTarget(compiler AbstractCompiler, callNode rhino.Node, assign ast.Assign) rhino.Node {
	madeUpName := "exportsForPolymer$jscomp" + compiler.GetUniqueNameIdSupplier()()
	assignExpr := assign.Node().GetParent()

	exportName := assign.Target().Node()
	target := rhino.IR_Name(madeUpName).ClonePropsFrom(exportName).Srcref(exportName)
	callNode.ReplaceWith(target)
	newDecl := rhino.IR_VarWithValue(target.CloneNode(), callNode).Srcref(assignExpr)
	newDecl.InsertBefore(assignExpr)
	newDecl.SetJSDocInfo(assignExpr.GetJSDocInfo())
	assignExpr.SetJSDocInfo(nil)
	return target
}

// PolymerClassDefinition_ExtractFromClassNode validates the class definition
// and if valid, non-destructively extracts the class definition from the AST.
//
// Java: PolymerClassDefinition.extractFromClassNode(Node, AbstractCompiler)
func PolymerClassDefinition_ExtractFromClassNode(classNode rhino.Node, compiler AbstractCompiler) *PolymerClassDefinition {
	cls := ast.MustClass(classNode)

	var propertiesDescriptor rhino.Node
	propertiesGetter := NodeUtil_GetFirstGetterMatchingKey(cls.Body().Node(), "properties")
	if !propertiesGetter.IsNil() {
		getterDef := ast.MustGetterDef(propertiesGetter)
		if !getterDef.IsStatic() {
			compiler.Report(
				JSError_MakeWithNode(classNode, PolymerPassErrors_POLYMER_CLASS_PROPERTIES_NOT_STATIC))
		} else {
			for stmt := range getterDef.Function().Body().MustBlock().Statements() {
				if ret, ok := stmt.AsReturn(); ok {
					if val, hasVal := ret.Value(); hasVal && val.Node().IsObjectLit() {
						propertiesDescriptor = val.Node()
						break
					}
					compiler.Report(
						JSError_MakeWithNode(
							propertiesGetter, PolymerPassErrors_POLYMER_CLASS_PROPERTIES_INVALID))
				}
			}
		}
	}

	var target rhino.Node
	if NodeUtil_IsNameDeclaration(classNode.GetGrandparent()) {
		target = rhino.IR_Name(classNode.GetParent().GetString())
	} else if assign, ok := ast.AsAssign(classNode.GetParent()); ok &&
		assign.Target().Node().IsQualifiedName() {
		target = assign.Target().Node()
	} else if name, ok := cls.Name(); ok {
		target = name.Node()
	} else {
		compiler.Report(JSError_MakeWithNode(classNode, PolymerPassErrors_POLYMER_CLASS_UNNAMED))
		return nil
	}
	enclosingModule := NodeUtil_GetEnclosingModuleIfPresent(classNode)

	classInfo := NodeUtil_GetBestJSDocInfo(classNode)
	constructor := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(classNode)
	properties := PolymerPassStaticUtils_ExtractProperties(
		propertiesDescriptor, PolymerClassDefinition_DefinitionType_ES6Class, compiler, constructor)
	methods := extractPolymerClassMethods(cls, enclosingModule)

	var ctorInfo *rhino.JSDocInfo
	if !constructor.IsNil() {
		ctorInfo = NodeUtil_GetBestJSDocInfo(constructor)
		methods = addFunctionAssignmentsFromConstructor(constructor, enclosingModule, methods)
	}
	input := compiler.GetInput(NodeUtil_GetEnclosingScript(classNode).GetInputId())

	return NewPolymerClassDefinition(
		PolymerClassDefinition_DefinitionType_ES6Class,
		classNode,
		target,
		false,
		propertiesDescriptor,
		classInfo,
		NewPolymerPass_MemberDefinition(ctorInfo, rhino.NilNode, constructor, enclosingModule),
		nil,
		properties,
		nil,
		methods,
		nil,
		parser.FeatureSet{},
		input)
}

func extractPolymerClassMethods(cls ast.Class, enclosingModule rhino.Node) []*PolymerPass_MemberDefinition {
	methods := []*PolymerPass_MemberDefinition{}
	for member := range cls.Body().Members() {
		if fnDef, ok := member.AsMemberFunctionDef(); ok {
			keyNode := fnDef.Node()
			methods = append(methods,
				NewPolymerPass_MemberDefinition(
					NodeUtil_GetBestJSDocInfo(keyNode),
					keyNode,
					fnDef.Function().Node(),
					enclosingModule))
		} else if fieldDef, ok := member.AsMemberFieldDef(); ok {
			if val, hasVal := fieldDef.Value(); hasVal && val.Node().IsFunction() {
				keyNode := fieldDef.Node()
				methods = append(methods,
					NewPolymerPass_MemberDefinition(
						NodeUtil_GetBestJSDocInfo(keyNode),
						keyNode,
						val.Node(),
						enclosingModule))
			}
		}
	}
	return methods
}

func addFunctionAssignmentsFromConstructor(
	constructor rhino.Node, enclosingModule rhino.Node, methods []*PolymerPass_MemberDefinition) []*PolymerPass_MemberDefinition {
	for stmt := range ast.MustMemberFunctionDef(constructor).Function().Body().MustBlock().Statements() {
		exprStmt, ok := stmt.AsExprStmt()
		if !ok {
			continue
		}
		assign, ok := exprStmt.Expr().AsAssign()
		if !ok {
			continue
		}
		propValue := assign.Value().Node()
		if prop, ok := assign.Target().AsGetProp(); ok && prop.Object().Node().IsThis() && propValue.IsFunction() {
			propName := prop.Node()
			methods = append(methods,
				NewPolymerPass_MemberDefinition(
					NodeUtil_GetBestJSDocInfo(propName), propName, propValue, enclosingModule))
		}
	}
	return methods
}

func overwriteMembersIfPresent(
	list []*PolymerPass_MemberDefinition, newMembers []*PolymerPass_MemberDefinition) []*PolymerPass_MemberDefinition {
	for _, newMember := range newMembers {
		newID := newMember.name.GetStringID()
		for i, member := range list {
			if member.name.GetStringID() == newID {
				list = append(list[:i:i], list[i+1:]...)
				break
			}
		}
		list = append(list, newMember)
	}
	return list
}

// Java: PolymerClassDefinition.toString()
func (d *PolymerClassDefinition) String() string {
	var b strings.Builder
	b.WriteString("PolymerClassDefinition{defType=")
	b.WriteString(d.defType.String())
	if !d.definition.IsNil() {
		b.WriteString(", definition=")
		b.WriteString(d.definition.String())
	}
	if !d.target.IsNil() {
		b.WriteString(", target=")
		b.WriteString(d.target.String())
	}
	if d.nativeBaseElement != nil {
		b.WriteString(", nativeBaseElement=")
		b.WriteString(*d.nativeBaseElement)
	}
	b.WriteByte('}')
	return b.String()
}

// Java: PolymerClassDefinition.getInterfaceName(UniqueIdSupplier)
func (d *PolymerClassDefinition) GetInterfaceName(uniqueIdSupplier *UniqueIdSupplier) string {
	if d.interfaceName == "" {
		d.interfaceName = "Polymer" +
			strings.ReplaceAll(d.targetQualifiedName(), ".", "_") +
			"Interface$" +
			uniqueIdSupplier.GetUniqueId(d.input)
	}
	return d.interfaceName
}
