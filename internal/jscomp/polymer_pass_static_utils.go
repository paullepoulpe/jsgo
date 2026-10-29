// Port of com/google/javascript/jscomp/PolymerPassStaticUtils.java @ closure-compiler v20261021

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

// Port of com/google/common/base/CaseFormat.java (Guava 33.4.6-jre)

/*
 * Copyright (C) 2006 The Guava Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const polymerPassStaticUtils_VIRTUAL_FILE = "<PolymerPassStaticUtils.java>"

var polymerPassStaticUtils_POLYMER_DOT_ELEMENT = rhino.QualifiedName_Of("Polymer.Element")

// PolymerPassStaticUtils_IsPolymerCall returns whether the call represents a
// call to the Polymer function.
//
// Java: PolymerPassStaticUtils.isPolymerCall(Node)
func PolymerPassStaticUtils_IsPolymerCall(call rhino.Node) bool {
	callNode, ok := ast.AsCall(call)
	if !ok {
		return false
	}
	return polymerMatchesNameOrProp(callNode.Callee().Node(), "Polymer")
}

// PolymerPassStaticUtils_IsPolymerClass returns whether the class extends
// PolymerElement.
//
// Java: PolymerPassStaticUtils.isPolymerClass(Node)
func PolymerPassStaticUtils_IsPolymerClass(cls rhino.Node) bool {
	classNode, ok := ast.AsClass(cls)
	if !ok {
		return false
	}
	// A class with the @polymer annotation is always considered a Polymer element.
	info := NodeUtil_GetBestJSDocInfo(cls)
	if info != nil && info.IsPolymer() {
		return true
	}
	ext, ok := classNode.Extends()
	if !ok {
		return false
	}
	heritage := ext.Node()
	return polymerPassStaticUtils_POLYMER_DOT_ELEMENT.Matches(heritage) ||
		polymerMatchesNameOrProp(heritage, "PolymerElement") ||
		polymerMatchesNameOrProp(heritage, "PolymerElementWithoutHtml") ||
		polymerMatchesNameOrProp(heritage, "PolymerLiteControllerBase") ||
		polymerMatchesNameOrProp(heritage, "PolymerBehavior") ||
		polymerMatchesNameOrProp(heritage, "PolymerBehaviorWithoutHtml")
}

func polymerMatchesNameOrProp(n rhino.Node, name string) bool {
	if n.MatchesName(name) {
		return true
	}
	if originalQualifiedName := n.GetOriginalQualifiedName(); originalQualifiedName != nil && name == *originalQualifiedName {
		return true
	}
	gp, ok := ast.AsGetProp(n)
	return ok && gp.Property() == name
}

// PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets rewrites usage of a
// Polymer element's "$" local-DOM map from "this.$.foo" to "this.$['foo']" to
// avoid JSC_POSSIBLE_INEXISTENT_PROPERTY errors, excluding zero-arg function
// calls like "bar.$.foo()".
//
// Java: PolymerPassStaticUtils.switchDollarSignPropsToBrackets(Node, AbstractCompiler)
func PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(def rhino.Node, compiler AbstractCompiler) {
	jutil.CheckState(def.IsObjectLit() || def.IsClassMembers())
	for keyNode := def.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		fn, ok := ast.AsFunction(keyNode.GetFirstChild())
		if !ok {
			continue
		}
		NodeUtil_VisitPostOrder(
			fn.Body().Node(),
			NodeUtil_VisitorFunc(func(n rhino.Node) {
				gp, ok := ast.AsGetProp(n)
				if !ok {
					return
				}
				dollarSign, ok := gp.Object().AsGetProp()
				if !ok || dollarSign.Property() != "$" {
					return
				}
				// Some libraries like Mojo JS Bindings and jQuery place methods in a "$" member
				// e.g. "foo.$.closePipe()". Avoid converting to brackets for these cases.
				if call, ok := ast.AsCall(n.GetParent()); ok && !call.HasArgs() {
					return
				}
				key := rhino.IR_String(gp.Property()).ClonePropsFrom(n).Srcref(n)
				getelem := rhino.IR_Getelem(dollarSign.Detach().Node(), key).ClonePropsFrom(n).Srcref(n)
				n.ReplaceWith(getelem)
				compiler.ReportChangeToEnclosingScope(getelem)
			}))
	}
}

// PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys ensures keys for
// listeners and hostAttributes blocks are quoted to avoid renaming.
//
// Java: PolymerPassStaticUtils.quoteListenerAndHostAttributeKeys(Node, AbstractCompiler)
func PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(objLit rhino.Node, compiler AbstractCompiler) {
	for member := range ast.MustObjectLit(objLit).Members() {
		sk, ok := member.AsStringKey()
		if !ok {
			continue
		}
		if sk.Key() != "listeners" && sk.Key() != "hostAttributes" {
			continue
		}
		for keyToQuote := sk.Value().Node().GetFirstChild(); !keyToQuote.IsNil(); keyToQuote = keyToQuote.GetNext() {
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
// block declaration, are exported and marked `@nocollapse`.
//
// Java: PolymerPassStaticUtils.protectObserverAndPropertyFunctionKeys(Node)
func PolymerPassStaticUtils_ProtectObserverAndPropertyFunctionKeys(behaviorObjLit rhino.Node) {
	jutil.CheckState(behaviorObjLit.IsObjectLit())

	var methodsToProtect map[string]struct{}
	addMethod := func(name string) {
		if methodsToProtect == nil {
			methodsToProtect = make(map[string]struct{})
		}
		methodsToProtect[name] = struct{}{}
	}

	observers := NodeUtil_GetFirstPropMatchingKey(behaviorObjLit, "observers")
	if !observers.IsNil() && observers.IsArrayLit() {
		for child := observers.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if child.IsStringLit() {
				if name, ok := polymerExtractMethodName(child.GetString()); ok {
					addMethod(name)
				}
			}
		}
	}

	properties := NodeUtil_GetFirstPropMatchingKey(behaviorObjLit, "properties")
	if propsObj, ok := ast.AsObjectLit(properties); ok {
		for member := range propsObj.Members() {
			sk, ok := member.AsStringKey()
			if !ok {
				continue
			}
			propVal := sk.Value().Node()
			if !propVal.IsObjectLit() {
				continue
			}

			observer := NodeUtil_GetFirstPropMatchingKey(propVal, "observer")
			if !observer.IsNil() && observer.IsStringLit() {
				addMethod(observer.GetString())
			}

			computed := NodeUtil_GetFirstPropMatchingKey(propVal, "computed")
			if !computed.IsNil() && computed.IsStringLit() {
				if name, ok := polymerExtractMethodName(computed.GetString()); ok {
					addMethod(name)
				}
			}
		}
	}

	if len(methodsToProtect) == 0 {
		return
	}

	for keyNode := behaviorObjLit.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		if !NodeUtil_IsObjectLitKey(keyNode) {
			continue
		}
		if _, protect := methodsToProtect[NodeUtil_GetObjectOrClassLitKeyName(keyNode)]; !protect {
			continue
		}
		newDocs := rhino.JSDocInfo_Builder_MaybeCopyFrom(keyNode.GetJSDocInfo())
		newDocs.RecordNoCollapse()
		newDocs.RecordExport()
		keyNode.SetJSDocInfo(newDocs.Build())
	}
}

func polymerExtractMethodName(signature string) (string, bool) {
	// '(' is ASCII: its byte index is positive iff its UTF-16 index is.
	if openParenIndex := strings.IndexByte(signature, '('); openParenIndex > 0 {
		return jutil.StringTrim(signature[:openParenIndex]), true
	}
	return "", false
}

// PolymerPassStaticUtils_ExtractProperties extracts MemberDefinitions for the
// `properties` block of the given descriptor Object literal.
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

	var constructorPropertyJsDoc map[string]*rhino.JSDocInfo
	if !constructor.IsNil() {
		constructorPropertyJsDoc = make(map[string]*rhino.JSDocInfo)
		collectConstructorPropertyJsDoc(constructor, constructorPropertyJsDoc)
	}
	enclosingModule := NodeUtil_GetEnclosingModuleIfPresent(descriptor)

	var members []*PolymerPass_MemberDefinition
	for keyNode := properties.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		// The JSDoc for a Polymer property in the constructor wins over the JSDoc in the
		// Polymer properties configuration object.
		constructorJsDoc := constructorPropertyJsDoc[keyNode.GetString()]
		propertiesConfigJsDoc := NodeUtil_GetBestJSDocInfo(keyNode)
		bestJsDoc := propertiesConfigJsDoc
		if constructorJsDoc != nil {
			bestJsDoc = constructorJsDoc
			if propertiesConfigJsDoc != nil {
				compiler.Report(JSError_MakeWithNode(keyNode, PolymerPassErrors_POLYMER_MISPLACED_PROPERTY_JSDOC))
			}
		}
		members = append(members,
			NewPolymerPass_MemberDefinition(bestJsDoc, keyNode, keyNode.GetFirstChild(), enclosingModule))
	}
	return members
}

func collectConstructorPropertyJsDoc(node rhino.Node, m map[string]*rhino.JSDocInfo) {
	rhino.CheckNotNull(node)
	for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if gp, ok := ast.AsGetProp(child); ok && gp.Object().Node().IsThis() {
			m[gp.Property()] = NodeUtil_GetBestJSDocInfo(child)
		} else {
			collectConstructorPropertyJsDoc(child, m)
		}
	}
}

// PolymerPassStaticUtils_GetTypeFromProperty gets the JSTypeExpression for a
// given property using its "type" key.
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
	}

	if typeValue.IsNil() {
		compiler.Report(JSError_MakeWithNode(property.value, PolymerPassErrors_POLYMER_INVALID_PROPERTY))
		return nil
	}

	var typeNode rhino.Node
	switch typeString := typeValue.GetString(); typeString {
	case "Boolean":
		typeNode = rhino.IR_String("boolean")
	case "String":
		typeNode = rhino.IR_String("string")
	case "Number":
		typeNode = rhino.IR_String("number")
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
	if cls.nativeBaseElement == nil {
		return "PolymerElement"
	}
	return "Polymer" + caseFormatLowerHyphenToUpperCamel(*cls.nativeBaseElement) + "Element"
}

// caseFormatLowerCamelToUpperUnderscore is Guava's
// CaseFormat.LOWER_CAMEL.to(CaseFormat.UPPER_UNDERSCORE, s): words split at
// each ASCII uppercase letter after index 0, are upper-cased
// (Ascii.toUpperCase), and are joined with '_'.
func caseFormatLowerCamelToUpperUnderscore(s string) string {
	extra := 0
	needsChange := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			needsChange = true
		} else if i > 0 && c >= 'A' && c <= 'Z' {
			extra++
			needsChange = true
		}
	}
	if !needsChange {
		return s
	}
	var sb strings.Builder
	sb.Grow(len(s) + extra)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			sb.WriteByte(c - ('a' - 'A'))
		} else if i > 0 && c >= 'A' && c <= 'Z' {
			sb.WriteByte('_')
			sb.WriteByte(c)
		} else {
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

// caseFormatLowerHyphenToUpperCamel is Guava's
// CaseFormat.LOWER_HYPHEN.to(CaseFormat.UPPER_CAMEL, s): words end at each
// '-', each is normalized by UPPER_CAMEL (first char Ascii.toUpperCase, the
// rest Ascii.toLowerCase), and they are joined without a separator.
func caseFormatLowerHyphenToUpperCamel(s string) string {
	hyphens := 0
	needsChange := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' {
			hyphens++
			needsChange = true
		} else if i == 0 && c >= 'a' && c <= 'z' {
			needsChange = true
		} else if i > 0 && c >= 'A' && c <= 'Z' {
			needsChange = true
		}
	}
	if !needsChange {
		return s
	}
	var sb strings.Builder
	sb.Grow(len(s) - hyphens)
	atWordStart := true
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' {
			atWordStart = true
			continue
		}
		if atWordStart {
			atWordStart = false
			if c >= 'a' && c <= 'z' {
				c -= 'a' - 'A'
			}
		} else if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		sb.WriteByte(c)
	}
	return sb.String()
}
