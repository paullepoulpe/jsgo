// Port of com/google/javascript/jscomp/PolymerClassDefinition.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PolymerClassDefinition is a parsed Polymer class (element) definition.
// Includes convenient fields for rewriting the class.
//
// Java: final class PolymerClassDefinition
type PolymerClassDefinition struct {
	// The declaration style used for the Polymer definition
	defType PolymerClassDefinition_DefinitionType

	// The Polymer call or class node which defines the Element.
	definition rhino.Node

	input *CompilerInput

	// The target node (LHS) for the Polymer element definition.
	target rhino.Node

	// Whether the target of this element is a generated node
	hasGeneratedLhs bool

	// The object literal passed to the call to the Polymer() function.
	descriptor rhino.Node

	// The constructor function for the element.
	constructor *PolymerPass_MemberDefinition

	// The name of the native HTML element which this element extends.
	nativeBaseElement *string

	// Properties declared in the Polymer "properties" block.
	props []*PolymerPass_MemberDefinition

	// Properties declared in the Polymer behavior's "properties" block.
	behaviorProps *jutil.LinkedHashMap[*PolymerPass_MemberDefinition, *PolymerBehaviorExtractor_BehaviorDefinition]

	// Methods on the element.
	methods []*PolymerPass_MemberDefinition

	// Flattened list of behavior definitions used by this element.
	behaviors []*PolymerBehaviorExtractor_BehaviorDefinition

	// Language features that should be carried over to the extraction destination.
	features *parser.FeatureSet

	interfaceName *string
}

// Java: static enum PolymerClassDefinition.DefinitionType
type PolymerClassDefinition_DefinitionType int32

const (
	PolymerClassDefinition_DefinitionType_ObjectLiteral PolymerClassDefinition_DefinitionType = iota
	PolymerClassDefinition_DefinitionType_ES6Class
)

var polymerClassDefinition_DefinitionType_names = [...]string{"ObjectLiteral", "ES6Class"}

// PolymerClassDefinition_DefinitionType_values is DefinitionType.values().
func PolymerClassDefinition_DefinitionType_values() []PolymerClassDefinition_DefinitionType {
	return []PolymerClassDefinition_DefinitionType{
		PolymerClassDefinition_DefinitionType_ObjectLiteral,
		PolymerClassDefinition_DefinitionType_ES6Class,
	}
}

// String is DefinitionType.name().
func (t PolymerClassDefinition_DefinitionType) String() string {
	return polymerClassDefinition_DefinitionType_names[t]
}

// Java: PolymerClassDefinition(DefinitionType, Node, Node, boolean, Node, JSDocInfo, MemberDefinition, String, List<MemberDefinition>, Map<MemberDefinition, BehaviorDefinition>, List<MemberDefinition>, ImmutableList<BehaviorDefinition>, FeatureSet, CompilerInput)
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
	behaviorProps *jutil.LinkedHashMap[*PolymerPass_MemberDefinition, *PolymerBehaviorExtractor_BehaviorDefinition],
	methods []*PolymerPass_MemberDefinition,
	behaviors []*PolymerBehaviorExtractor_BehaviorDefinition,
	features *parser.FeatureSet,
	input *CompilerInput) *PolymerClassDefinition {
	d := &PolymerClassDefinition{}
	d.interfaceName = nil
	d.defType = defType
	d.definition = definition
	d.target = target
	d.hasGeneratedLhs = hasGeneratedLhs
	jutil.CheckState(descriptor.IsNil() || descriptor.IsObjectLit())
	d.descriptor = descriptor
	d.constructor = constructor
	d.nativeBaseElement = nativeBaseElement
	d.props = props
	d.behaviorProps = behaviorProps
	d.methods = methods
	d.behaviors = behaviors
	d.features = features
	d.input = input
	return d
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
	descriptor := NodeUtil_GetArgumentForCallOrNew(callNode, 0)
	if descriptor.IsNil() || !descriptor.IsObjectLit() {
		// report bad class definition
		compiler.Report(JSError_MakeWithNode(callNode, PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID))
		return nil
	}

	paramCount := callNode.GetChildCount() - 1
	if paramCount != 1 {
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
	if NodeUtil_IsNameDeclaration(callNode.GetGrandparent()) {
		target = rhino.IR_Name(callNode.GetParent().GetString())
	} else if callNode.GetParent().IsAssign() {
		if polymerClassDefinition_isGoogModuleExports(callNode.GetParent()) {
			// `exports = Polymer({` in a goog.module requires special handling, as adding a
			// duplicate assignment to exports just confuses the compiler. Create a dummy declaration
			// var exportsForPolymer$jscomp0 = Polymer({ // ...
			target = polymerClassDefinition_createDummyGoogModuleExportsTarget(compiler, callNode)
		} else {
			target = callNode.GetParent().GetFirstChild().CloneTree()
		}
	} else {
		var elNameStringBase string
		if elName.IsQualifiedName() {
			elNameStringBase = strings.ReplaceAll(*elName.GetQualifiedName(), ".", "$")
		} else {
			elNameStringBase = elName.GetString()
		}
		elNameString := jutil.CaseFormat_LowerHyphenToUpperCamel(elNameStringBase)
		elNameString += "Element"
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
	properties := []*PolymerPass_MemberDefinition{}
	behaviorProps := jutil.NewLinkedHashMap[*PolymerPass_MemberDefinition, *PolymerBehaviorExtractor_BehaviorDefinition]()
	for _, behavior := range behaviors {
		for _, prop := range behavior.props {
			behaviorProps.Put(prop, behavior)
		}
	}
	properties = polymerClassDefinition_overwriteMembersIfPresent(
		properties,
		PolymerPassStaticUtils_ExtractProperties(
			descriptor, PolymerClassDefinition_DefinitionType_ObjectLiteral, compiler /* constructor= */, rhino.NilNode))

	// Behaviors might get included multiple times for the same element. See test case
	// testDuplicatedBehaviorsAreCopiedOnce
	// In those cases, behaviorProps will have repeated names. We must remove those duplicates.
	polymerClassDefinition_removeDuplicateBehaviorProps(behaviorProps)

	// Remove behavior properties which are already present in element properties
	polymerClassDefinition_removeBehaviorPropsOverlappingWithElementProps(behaviorProps, properties)

	var newFeatures *parser.FeatureSet
	if len(behaviors) != 0 {
		newFeatures = behaviors[0].features
		for i := 1; i < len(behaviors); i++ {
			newFeatures = newFeatures.Union(behaviors[i].features)
		}
	}

	methods := []*PolymerPass_MemberDefinition{}
	for keyNode := descriptor.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		isFunctionDefinition := keyNode.IsMemberFunctionDef() ||
			(keyNode.IsStringKey() && keyNode.GetFirstChild().IsFunction())
		if isFunctionDefinition {
			methods = append(methods,
				NewPolymerPass_MemberDefinition(
					NodeUtil_GetBestJSDocInfo(keyNode),
					keyNode,
					keyNode.GetFirstChild(),
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

// Java: PolymerClassDefinition.isGoogModuleExports(Node)
func polymerClassDefinition_isGoogModuleExports(assign rhino.Node) bool {
	// Verify this is an assignment to a name `exports`
	if !assign.GetParent().IsExprResult() || !assign.GetFirstChild().MatchesName("exports") {
		return false
	}
	// Verify the assignment is within either a goog.module or goog.loadModule
	containingBlock := assign.GetGrandparent()
	return (containingBlock.IsModuleBody() &&
		containingBlock.GetParent().GetBooleanProp(rhino.Node_GOOG_MODULE)) ||
		NodeUtil_IsBundledGoogModuleScopeRoot(containingBlock)
}

// polymerClassDefinition_createDummyGoogModuleExportsTarget: adding our usual
// multiple assignments
//
//	/** @constructor @extends {PolymerElement} */
//	var FooElement = function() {};
//	FooElement = Polymer() {
//
//	}
//
// confuses goog.module rewriting when the Polymer call is assigned to
// `exports. Instead create a new module local, and export that name.
//
// callNode is a Polymer({}) call. Returns the new target node for the
// Polymer call.
//
// Java: PolymerClassDefinition.createDummyGoogModuleExportsTarget(AbstractCompiler, Node)
func polymerClassDefinition_createDummyGoogModuleExportsTarget(compiler AbstractCompiler, callNode rhino.Node) rhino.Node {
	madeUpName := "exportsForPolymer$jscomp" + compiler.GetUniqueNameIdSupplier()()
	assignExpr := callNode.GetGrandparent()

	exportName := callNode.GetParent().GetFirstChild()
	target := rhino.IR_Name(madeUpName).ClonePropsFrom(exportName).Srcref(exportName)
	callNode.ReplaceWith(target)
	newDecl := rhino.IR_VarWithValue(target.CloneNode(), callNode).Srcref(assignExpr)
	newDecl.InsertBefore(assignExpr)
	newDecl.SetJSDocInfo(assignExpr.GetJSDocInfo())
	assignExpr.SetJSDocInfo(nil)
	return target
}

// PolymerClassDefinition_ExtractFromClassNode validates the class definition
// and if valid, extracts the class definition from the AST. As opposed to
// the Polymer 1 extraction, this operation is non-destructive.
//
// Java: PolymerClassDefinition.extractFromClassNode(Node, AbstractCompiler)
func PolymerClassDefinition_ExtractFromClassNode(classNode rhino.Node, compiler AbstractCompiler) *PolymerClassDefinition {
	jutil.CheckState(!classNode.IsNil() && classNode.IsClass())

	// The supported case is for the config getter to return an object literal descriptor.
	var propertiesDescriptor rhino.Node
	propertiesGetter := NodeUtil_GetFirstGetterMatchingKey(NodeUtil_GetClassMembers(classNode), "properties")
	if !propertiesGetter.IsNil() {
		if !propertiesGetter.IsStaticMember() {
			// report bad class definition
			compiler.Report(
				JSError_MakeWithNode(classNode, PolymerPassErrors_POLYMER_CLASS_PROPERTIES_NOT_STATIC))
		} else {
			for child := NodeUtil_GetFunctionBody(propertiesGetter.GetFirstChild()).GetFirstChild(); !child.IsNil(); child = child.GetNext() {
				if child.IsReturn() {
					if child.HasChildren() && child.GetFirstChild().IsObjectLit() {
						propertiesDescriptor = child.GetFirstChild()
						break
					} else {
						compiler.Report(
							JSError_MakeWithNode(
								propertiesGetter, PolymerPassErrors_POLYMER_CLASS_PROPERTIES_INVALID))
					}
				}
			}
		}
	}

	var target rhino.Node
	if NodeUtil_IsNameDeclaration(classNode.GetGrandparent()) {
		target = rhino.IR_Name(classNode.GetParent().GetString())
	} else if classNode.GetParent().IsAssign() &&
		classNode.GetParent().GetFirstChild().IsQualifiedName() {
		target = classNode.GetParent().GetFirstChild()
	} else if !classNode.GetFirstChild().IsEmpty() {
		target = classNode.GetFirstChild()
	} else {
		// issue error - no name found
		compiler.Report(JSError_MakeWithNode(classNode, PolymerPassErrors_POLYMER_CLASS_UNNAMED))
		return nil
	}
	enclosingModule := NodeUtil_GetEnclosingModuleIfPresent(classNode)

	classInfo := NodeUtil_GetBestJSDocInfo(classNode)

	constructor := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(classNode)

	properties := PolymerPassStaticUtils_ExtractProperties(
		propertiesDescriptor, PolymerClassDefinition_DefinitionType_ES6Class, compiler, constructor)

	methods := polymerClassDefinition_extractClassMethods(classNode, enclosingModule)

	var ctorInfo *rhino.JSDocInfo
	if !constructor.IsNil() {
		ctorInfo = NodeUtil_GetBestJSDocInfo(constructor)

		// Add properties assigned to functions in the constructor (i.e. `this.prop = function`).
		methods = polymerClassDefinition_addFunctionAssignmentsFromConstructor(constructor, enclosingModule, methods)
	}
	input := compiler.GetInput(NodeUtil_GetEnclosingScript(classNode).GetInputId())

	return NewPolymerClassDefinition(
		PolymerClassDefinition_DefinitionType_ES6Class,
		classNode,
		target,
		/* hasGeneratedLhs= */ false,
		propertiesDescriptor,
		classInfo,
		NewPolymerPass_MemberDefinition(ctorInfo, rhino.NilNode, constructor, enclosingModule),
		nil,
		properties,
		nil,
		methods,
		nil,
		nil,
		input)
}

// polymerClassDefinition_extractClassMethods extracts true class methods and
// fields assigned to functions from an ES6 class node.
//
// Java: PolymerClassDefinition.extractClassMethods(Node, Node)
func polymerClassDefinition_extractClassMethods(classNode rhino.Node, enclosingModule rhino.Node) []*PolymerPass_MemberDefinition {
	methods := []*PolymerPass_MemberDefinition{}
	for keyNode := NodeUtil_GetClassMembers(classNode).GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		if keyNode.IsMemberFunctionDef() || polymerClassDefinition_isFunctionField(keyNode) {
			methods = append(methods,
				NewPolymerPass_MemberDefinition(
					NodeUtil_GetBestJSDocInfo(keyNode),
					keyNode,
					keyNode.GetFirstChild(),
					enclosingModule))
		}
	}
	return methods
}

// polymerClassDefinition_isFunctionField detects patterns like
// `field = function() {}` and `field = () => {})`.
//
// Java: PolymerClassDefinition.isFunctionField(Node)
func polymerClassDefinition_isFunctionField(field rhino.Node) bool {
	if !field.IsMemberFieldDef() {
		return false
	}
	value := field.GetFirstChild()
	return !value.IsNil() && value.IsFunction()
}

// polymerClassDefinition_addFunctionAssignmentsFromConstructor adds
// properties assigned to functions in the constructor body (e.g.,
// `this.prop = function()`) to the given list of methods.
//
// Go-only: language-adapter — the Java method appends to the given list in
// place; this returns the appended slice.
//
// Java: PolymerClassDefinition.addFunctionAssignmentsFromConstructor(Node, Node, List<MemberDefinition>)
func polymerClassDefinition_addFunctionAssignmentsFromConstructor(
	constructor rhino.Node, enclosingModule rhino.Node, methods []*PolymerPass_MemberDefinition) []*PolymerPass_MemberDefinition {
	body := NodeUtil_GetFunctionBody(constructor.GetFirstChild())
	for stmt := body.GetFirstChild(); !stmt.IsNil(); stmt = stmt.GetNext() {
		if !stmt.IsExprResult() || !stmt.GetFirstChild().IsAssign() {
			continue
		}

		assignment := stmt.GetFirstChild()
		propName := assignment.GetFirstChild()
		propValue := assignment.GetLastChild()
		if propName.IsGetProp() && propName.GetFirstChild().IsThis() && propValue.IsFunction() {
			methods = append(methods,
				NewPolymerPass_MemberDefinition(
					NodeUtil_GetBestJSDocInfo(propName), propName, propValue, enclosingModule))
		}
	}
	return methods
}

// polymerClassDefinition_overwriteMembersIfPresent appends a list of new
// MemberDefinitions to the end of a list and removes any previous
// MemberDefinition in the list which has the same name as the new member.
//
// Go-only: language-adapter — the Java method mutates the given list in
// place; this returns the updated slice.
//
// Java: PolymerClassDefinition.overwriteMembersIfPresent(List<MemberDefinition>, List<MemberDefinition>)
func polymerClassDefinition_overwriteMembersIfPresent(
	list []*PolymerPass_MemberDefinition, newMembers []*PolymerPass_MemberDefinition) []*PolymerPass_MemberDefinition {
	for _, newMember := range newMembers {
		for i, member := range list {
			if member.name.GetString() == newMember.name.GetString() {
				list = append(list[:i:i], list[i+1:]...)
				break
			}
		}
		list = append(list, newMember)
	}
	return list
}

// polymerClassDefinition_removeDuplicateBehaviorProps removes duplicate
// properties from the given map, keeping only the first property
//
// Duplicates occur when either the same behavior is transitively added
// multiple times to Polymer element, or when two unique added behaviors share
// a property with the same name.
//
// Java: PolymerClassDefinition.removeDuplicateBehaviorProps(Map<MemberDefinition, BehaviorDefinition>)
func polymerClassDefinition_removeDuplicateBehaviorProps(
	behaviorProps *jutil.LinkedHashMap[*PolymerPass_MemberDefinition, *PolymerBehaviorExtractor_BehaviorDefinition]) {
	if behaviorProps == nil {
		return
	}
	seen := jutil.NewLinkedHashSet[string]()
	for memberDefinition := range behaviorProps.Keys() {
		propertyName := memberDefinition.name.GetString()
		if !seen.Add(propertyName) {
			behaviorProps.Remove(memberDefinition) // behaviorsItr.remove()
		}
	}
}

// polymerClassDefinition_removeBehaviorPropsOverlappingWithElementProps
// removes any behavior properties from the given map that have the same name
// as a property in the given list.
//
// For example, if a Polymer element with a property "name" depends on a
// behavior with a property "name", then this method removes the behavior
// property in favor of the element property.
//
// Java: PolymerClassDefinition.removeBehaviorPropsOverlappingWithElementProps(Map<MemberDefinition, BehaviorDefinition>, List<MemberDefinition>)
func polymerClassDefinition_removeBehaviorPropsOverlappingWithElementProps(
	behaviorProps *jutil.LinkedHashMap[*PolymerPass_MemberDefinition, *PolymerBehaviorExtractor_BehaviorDefinition],
	polymerElementProps []*PolymerPass_MemberDefinition) {
	if behaviorProps == nil {
		return
	}
	elementPropNames := jutil.NewLinkedHashSet[string]()
	for _, x := range polymerElementProps {
		elementPropNames.Add(x.name.GetString())
	}
	for memberDefinition := range behaviorProps.Keys() {
		if elementPropNames.Contains(memberDefinition.name.GetString()) {
			behaviorProps.Remove(memberDefinition) // behaviorsItr.remove()
		}
	}
}

// Java: PolymerClassDefinition.toString()
func (d *PolymerClassDefinition) String() string {
	return jutil.MoreObjects_ToStringHelper("PolymerClassDefinition").
		Add("defType", d.defType).
		Add("definition", d.definition).
		Add("target", d.target).
		Add("nativeBaseElement", d.nativeBaseElement).
		OmitNullValues().
		String()
}

// Java: PolymerClassDefinition.getInterfaceName(UniqueIdSupplier)
func (d *PolymerClassDefinition) GetInterfaceName(uniqueIdSupplier *UniqueIdSupplier) string {
	if d.interfaceName == nil {
		name := "Polymer" +
			strings.ReplaceAll(*d.target.GetQualifiedName(), ".", "_") +
			"Interface$" +
			uniqueIdSupplier.GetUniqueId(d.input)
		d.interfaceName = &name
	}
	return *d.interfaceName
}
