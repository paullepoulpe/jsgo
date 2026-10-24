// Port of com/google/javascript/jscomp/PolymerPassStaticUtils.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const polymerPassStaticUtils_VIRTUAL_FILE = "<PolymerPassStaticUtils.java>"

var polymerPassStaticUtils_POLYMER_DOT_ELEMENT = rhino.QualifiedName_Of("Polymer.Element")

// PolymerPassStaticUtils holds simple static utility functions shared
// between the PolymerPass and its helper classes.
//
// Java: final class PolymerPassStaticUtils
type PolymerPassStaticUtils struct{}

// PolymerPassStaticUtils_IsPolymerCall returns whether the call represents a
// call to the Polymer function.
//
// Java: PolymerPassStaticUtils.isPolymerCall(Node)
func PolymerPassStaticUtils_IsPolymerCall(call rhino.Node) bool {
	if call.IsNil() || !call.IsCall() {
		return false
	}
	name := call.GetFirstChild()
	// When imported from an ES module, the rewriting should set the original name.
	// When imported from an goog module (TS), we'll have a GETPROP like
	// `module$polymer$polymer_legacy.Polymer`.
	if name.MatchesName("Polymer") {
		return true
	}
	if originalQualifiedName := name.GetOriginalQualifiedName(); originalQualifiedName != nil && *originalQualifiedName == "Polymer" {
		return true
	}
	return name.IsGetProp() && name.GetString() == "Polymer"
}

// PolymerPassStaticUtils_IsPolymerClass returns whether the class extends
// PolymerElement.
//
// Java: PolymerPassStaticUtils.isPolymerClass(Node)
func PolymerPassStaticUtils_IsPolymerClass(cls rhino.Node) bool {
	if cls.IsNil() || !cls.IsClass() {
		return false
	}
	// A class with the @polymer annotation is always considered a Polymer element.
	info := NodeUtil_GetBestJSDocInfo(cls)
	if info != nil && info.IsPolymer() {
		return true
	}
	heritage := cls.GetSecondChild()
	// In Polymer 3, the base class was renamed from `Polymer.Element` to `PolymerElement`. When
	// imported from an ES module, the rewriting should set the original name to `PolymerElement`.
	// When imported from an goog module (TS), we'll have a GETPROP like
	// `module$polymer$polymer_element.PolymerElement`.
	// YT Polymer components may also extend the `PolymerElementWithoutHtml`,
	// `PolymerLiteControllerBase`, `PolymerBehavior`, or `PolymerBehaviorWithoutHtml` base classes.
	return !heritage.IsEmpty() &&
		(polymerPassStaticUtils_POLYMER_DOT_ELEMENT.Matches(heritage) ||
			polymerPassStaticUtils_matches(heritage, "PolymerElement") ||
			polymerPassStaticUtils_matches(heritage, "PolymerElementWithoutHtml") ||
			polymerPassStaticUtils_matches(heritage, "PolymerLiteControllerBase") ||
			polymerPassStaticUtils_matches(heritage, "PolymerBehavior") ||
			polymerPassStaticUtils_matches(heritage, "PolymerBehaviorWithoutHtml"))
}

// Java: PolymerPassStaticUtils.matches(Node, String)
func polymerPassStaticUtils_matches(n rhino.Node, name string) bool {
	if n.MatchesName(name) {
		return true
	}
	if originalQualifiedName := n.GetOriginalQualifiedName(); originalQualifiedName != nil && name == *originalQualifiedName {
		return true
	}
	return n.IsGetProp() && n.GetString() == name
}

// PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets: the "$" member in a
// Polymer element is a map of statically created nodes in its local DOM. This
// method is used to rewrite usage of this map from "this.$.foo" to
// "this.$['foo']" to avoid JSC_POSSIBLE_INEXISTENT_PROPERTY errors. Excludes
// function calls like "bar.$.foo()" since some libraries place methods under
// a "$" member.
//
// Java: PolymerPassStaticUtils.switchDollarSignPropsToBrackets(Node, AbstractCompiler)
func PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(def rhino.Node, compiler AbstractCompiler) {
	jutil.CheckState(def.IsObjectLit() || def.IsClassMembers())
	for keyNode := def.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		value := keyNode.GetFirstChild()
		if !value.IsNil() && value.IsFunction() {
			NodeUtil_VisitPostOrder(
				value.GetLastChild(),
				NodeUtil_VisitorFunc(func(n rhino.Node) {
					if !n.IsGetProp() {
						return
					}

					dollarSign := n.GetFirstChild()
					if !dollarSign.IsGetProp() || dollarSign.GetString() != "$" {
						return
					}

					// Some libraries like Mojo JS Bindings and jQuery place methods in a "$" member
					// e.g. "foo.$.closePipe()". Avoid converting to brackets for these cases.
					if n.GetParent().IsCall() && n.GetParent().HasOneChild() {
						return
					}

					key := rhino.IR_String(n.GetString()).ClonePropsFrom(n).Srcref(n)
					getelem := rhino.IR_Getelem(dollarSign.Detach(), key).ClonePropsFrom(n).Srcref(n)
					n.ReplaceWith(getelem)
					compiler.ReportChangeToEnclosingScope(getelem)
				}))
		}
	}
}

// PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys makes sure that
// the keys for listeners and hostAttributes blocks are quoted to avoid
// renaming.
//
// Java: PolymerPassStaticUtils.quoteListenerAndHostAttributeKeys(Node, AbstractCompiler)
func PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(objLit rhino.Node, compiler AbstractCompiler) {
	jutil.CheckState(objLit.IsObjectLit())
	for keyNode := objLit.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		if !keyNode.IsStringKey() {
			// We should only quote string keys. If this is not a string key, then we should skip it.
			continue
		}
		if keyNode.GetString() != "listeners" && keyNode.GetString() != "hostAttributes" {
			continue
		}
		for keyToQuote := keyNode.GetFirstFirstChild(); !keyToQuote.IsNil(); keyToQuote = keyToQuote.GetNext() {
			if !keyToQuote.IsQuotedStringKey() {
				keyToQuote.SetQuotedStringKey()
				compiler.ReportChangeToEnclosingScope(keyToQuote)
			}
		}
	}
}

// PolymerPassStaticUtils_ProtectObserverAndPropertyFunctionKeys ensures that
// methods registered in the static `observers` block, as well as methods
// referenced in the `observer` and `computed` fields on a Polymer property
// block declaration are exported and marked `@noCollapse`.
//
// Java: PolymerPassStaticUtils.protectObserverAndPropertyFunctionKeys(Node)
func PolymerPassStaticUtils_ProtectObserverAndPropertyFunctionKeys(behaviorObjLit rhino.Node) {
	jutil.CheckState(behaviorObjLit.IsObjectLit())

	methodsToProtect := jutil.NewLinkedHashSet[string]()
	methodsToProtect.AddAll(polymerPassStaticUtils_addObserverMethodsToProtect(behaviorObjLit).ToSlice()...)
	methodsToProtect.AddAll(polymerPassStaticUtils_addPropertyMethodsToProtect(behaviorObjLit).ToSlice()...)

	for keyNode := behaviorObjLit.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		// Note: This also is true for object literal methods.
		if !NodeUtil_IsObjectLitKey(keyNode) {
			continue
		}

		keyName := NodeUtil_GetObjectOrClassLitKeyName(keyNode)
		if !methodsToProtect.Contains(keyName) {
			continue
		}

		polymerPassStaticUtils_addExportAndNoCollapseToKey(keyNode)
	}
}

// polymerPassStaticUtils_addObserverMethodsToProtect finds methods
// referenced in the 'observers' property and returns them as a Set.
//
// Java: PolymerPassStaticUtils.addObserverMethodsToProtect(Node)
func polymerPassStaticUtils_addObserverMethodsToProtect(behaviorObjLit rhino.Node) *jutil.LinkedHashSet[string] {
	methodsToProtect := jutil.NewLinkedHashSet[string]()
	observers := NodeUtil_GetFirstPropMatchingKey(behaviorObjLit, "observers")
	if !observers.IsNil() && observers.IsArrayLit() {
		for child := observers.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if !child.IsStringLit() {
				continue
			}

			name := polymerPassStaticUtils_extractMethodName(child.GetString())
			if name != nil {
				methodsToProtect.Add(*name)
			}
		}
	}
	return methodsToProtect
}

// polymerPassStaticUtils_addPropertyMethodsToProtect finds methods
// referenced in the 'properties' property and returns them as a Set.
//
// Java: PolymerPassStaticUtils.addPropertyMethodsToProtect(Node)
func polymerPassStaticUtils_addPropertyMethodsToProtect(behaviorObjLit rhino.Node) *jutil.LinkedHashSet[string] {
	methodsToProtect := jutil.NewLinkedHashSet[string]()
	properties := NodeUtil_GetFirstPropMatchingKey(behaviorObjLit, "properties")
	if !properties.IsNil() && properties.IsObjectLit() {
		for prop := properties.GetFirstChild(); !prop.IsNil(); prop = prop.GetNext() {
			if !prop.GetLastChild().IsObjectLit() {
				continue
			}

			observer := NodeUtil_GetFirstPropMatchingKey(prop.GetLastChild(), "observer")
			if !observer.IsNil() && observer.IsStringLit() {
				methodsToProtect.Add(observer.GetString())
			}

			computed := NodeUtil_GetFirstPropMatchingKey(prop.GetLastChild(), "computed")
			if !computed.IsNil() && computed.IsStringLit() {
				name := polymerPassStaticUtils_extractMethodName(computed.GetString())
				if name != nil {
					methodsToProtect.Add(*name)
				}
			}
		}
	}
	return methodsToProtect
}

// Java: PolymerPassStaticUtils.extractMethodName(String)
func polymerPassStaticUtils_extractMethodName(signature string) *string {
	// '(' is ASCII: its byte index is positive exactly when its UTF-16 index is.
	openParenIndex := strings.IndexByte(signature, '(')
	if openParenIndex > 0 {
		name := jutil.StringTrim(signature[:openParenIndex])
		return &name
	}
	return nil
}

// polymerPassStaticUtils_addExportAndNoCollapseToKey ensures the given key
// node (or methodFunction) is exported and adds `@nocollapse` to its JSDoc.
//
// keyNode is the key node from an object literal.
//
// Java: PolymerPassStaticUtils.addExportAndNoCollapseToKey(Node)
func polymerPassStaticUtils_addExportAndNoCollapseToKey(keyNode rhino.Node) {
	newDocs := rhino.JSDocInfo_Builder_MaybeCopyFrom(keyNode.GetJSDocInfo())
	newDocs.RecordNoCollapse()
	newDocs.RecordExport()
	keyNode.SetJSDocInfo(newDocs.Build())
}

// PolymerPassStaticUtils_ExtractProperties finds a list of MemberDefinitions
// for the `properties` block of the given descriptor Object literal.
//
// descriptor is the Polymer properties configuration object literal node.
// constructor: if we are finding the properties of an ES6 class, the
// constructor function node for that class, otherwise null. We'll prefer
// JSDoc from property initialization statements in this constructor over the
// JSDoc within the Polymer properties configuration object.
//
// Java: PolymerPassStaticUtils.extractProperties(Node, PolymerClassDefinition.DefinitionType, AbstractCompiler, Node)
func PolymerPassStaticUtils_ExtractProperties(
	descriptor rhino.Node,
	defType PolymerClassDefinition_DefinitionType,
	compiler AbstractCompiler,
	constructor rhino.Node) []*PolymerPass_MemberDefinition {
	properties := descriptor
	if defType == PolymerClassDefinition_DefinitionType_ObjectLiteral {
		properties = NodeUtil_GetFirstPropMatchingKey(descriptor, "properties")
	}
	if properties.IsNil() {
		return nil
	}

	constructorPropertyJsDoc := jutil.NewLinkedHashMap[string, *rhino.JSDocInfo]()
	if !constructor.IsNil() {
		polymerPassStaticUtils_collectConstructorPropertyJsDoc(constructor, constructorPropertyJsDoc)
	}
	enclosingModule := NodeUtil_GetEnclosingModuleIfPresent(descriptor)

	var members []*PolymerPass_MemberDefinition
	for keyNode := properties.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		// The JSDoc for a Polymer property in the constructor should win over the JSDoc in the
		// Polymer properties configuration object.
		constructorJsDoc := constructorPropertyJsDoc.Get(keyNode.GetString())
		propertiesConfigJsDoc := NodeUtil_GetBestJSDocInfo(keyNode)
		var bestJsDoc *rhino.JSDocInfo
		if constructorJsDoc != nil {
			bestJsDoc = constructorJsDoc
			if propertiesConfigJsDoc != nil {
				// Warn if the user put JSDoc in both places.
				compiler.Report(JSError_MakeWithNode(keyNode, PolymerPassErrors_POLYMER_MISPLACED_PROPERTY_JSDOC))
			}
		} else {
			bestJsDoc = propertiesConfigJsDoc
		}
		members = append(members,
			NewPolymerPass_MemberDefinition(bestJsDoc, keyNode, keyNode.GetFirstChild(), enclosingModule))
	}
	return members
}

// polymerPassStaticUtils_collectConstructorPropertyJsDoc finds the
// properties that are initialized in the given constructor, and return a map
// from each property name to its JSDoc.
//
// node is the constructor function node to traverse; m the map from property
// name to JSDoc.
//
// Java: PolymerPassStaticUtils.collectConstructorPropertyJsDoc(Node, Map<String, JSDocInfo>)
func polymerPassStaticUtils_collectConstructorPropertyJsDoc(node rhino.Node, m *jutil.LinkedHashMap[string, *rhino.JSDocInfo]) {
	rhino.CheckNotNull(node)
	for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if child.IsGetProp() && child.GetFirstChild().IsThis() {
			// We found a "this.foo" expression. Map "foo" to its JSDoc.
			m.Put(child.GetString(), NodeUtil_GetBestJSDocInfo(child))
		} else {
			// Recurse through every other kind of node, because properties are not necessarily declared
			// at the top level of the constructor body; e.g. they could be declared as part of an
			// assignment, or within an if statement. We're being overly loose here, e.g. we shouldn't
			// traverse into a nested function where "this" doesn't refer to our prototype, but
			// hopefully this is good enough for our purposes.
			polymerPassStaticUtils_collectConstructorPropertyJsDoc(child, m)
		}
	}
}

// PolymerPassStaticUtils_GetTypeFromProperty gets the JSTypeExpression for a
// given property using its "type" key.
//
// See https://github.com/Polymer/polymer/blob/0.8-preview/PRIMER.md#configuring-properties
//
// Java: PolymerPassStaticUtils.getTypeFromProperty(MemberDefinition, AbstractCompiler)
func PolymerPassStaticUtils_GetTypeFromProperty(
	property *PolymerPass_MemberDefinition, compiler AbstractCompiler) *rhino.JSTypeExpression {
	if property.info != nil && property.info.HasType() {
		return property.info.GetType().Copy()
	}

	var typeValue rhino.Node
	if property.value.IsObjectLit() {
		typeValue = NodeUtil_GetFirstPropMatchingKey(property.value, "type")
		if typeValue.IsNil() || !typeValue.IsName() {
			compiler.Report(JSError_MakeWithNode(property.name, PolymerPassErrors_POLYMER_INVALID_PROPERTY))
			return nil
		}
	} else if property.value.IsName() {
		typeValue = property.value
	} else {
		typeValue = rhino.NilNode
	}

	if typeValue.IsNil() {
		compiler.Report(JSError_MakeWithNode(property.value, PolymerPassErrors_POLYMER_INVALID_PROPERTY))
		return nil
	}

	typeString := typeValue.GetString()
	var typeNode rhino.Node
	switch typeString {
	case "Boolean", "String", "Number":
		typeNode = rhino.IR_String(jutil.StringToLowerCaseASCII(typeString))
	case "Array", "Function", "Object", "Date":
		typeNode = rhino.NewNodeWithChild(rhino.Token_BANG, rhino.IR_String(typeString))
	default:
		compiler.Report(JSError_MakeWithNode(property.value, PolymerPassErrors_POLYMER_INVALID_PROPERTY))
		typeNode = rhino.NewNode(rhino.Token_QMARK)
	}

	typeNode.SrcrefTree(typeValue)
	return rhino.NewJSTypeExpression(typeNode, polymerPassStaticUtils_VIRTUAL_FILE)
}

// PolymerPassStaticUtils_GetPolymerElementType returns the PolymerElement
// type string for a class definition.
//
// Java: PolymerPassStaticUtils.getPolymerElementType(PolymerClassDefinition)
func PolymerPassStaticUtils_GetPolymerElementType(cls *PolymerClassDefinition) string {
	nativeElementName := ""
	if cls.nativeBaseElement != nil {
		nativeElementName = jutil.CaseFormat_LowerHyphenToUpperCamel(*cls.nativeBaseElement)
	}
	return "Polymer" + nativeElementName + "Element" // String.format("Polymer%sElement", ...)
}

// Java: PolymerPassStaticUtils()
func newPolymerPassStaticUtils() *PolymerPassStaticUtils { return &PolymerPassStaticUtils{} }
