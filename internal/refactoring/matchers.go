// Port of com/google/javascript/refactoring/Matchers.java @ closure-compiler v20261021
// Port of com/google/javascript/refactoring/Matcher.java @ closure-compiler v20261021
// Port of com/google/javascript/refactoring/Match.java @ closure-compiler v20261021
// Port of com/google/javascript/refactoring/NodeMetadata.java @ closure-compiler v20261021
// Port of com/google/javascript/refactoring/JsSourceMatcher.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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

package refactoring

import (
	"regexp"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Matcher is a function that knows how to match a Node for a specific
// pattern. For example of Matchers, see Matchers.
//
// Go-only: language-adapter — Java's @FunctionalInterface Matcher is a func
// type so closures convert directly without an adapter wrapper.
type Matcher func(n rhino.Node, metadata *NodeMetadata) bool

// Matches returns true if the specified Node and NodeMetadata match the
// given pattern.
//
// Java: Matcher.matches(Node, NodeMetadata)
func (m Matcher) Matches(n rhino.Node, metadata *NodeMetadata) bool {
	return m(n, metadata)
}

// Match is an object that contains the information for a given match.
type Match struct {
	node     rhino.Node
	metadata *NodeMetadata
}

// Java: Match(Node, NodeMetadata)
func NewMatch(node rhino.Node, metadata *NodeMetadata) *Match {
	return &Match{node: node, metadata: metadata}
}

// GetNode returns the node that matched the given conditions.
//
// Java: Match.getNode()
func (m *Match) GetNode() rhino.Node {
	return m.node
}

// GetMetadata returns the metadata for this match.
//
// Java: Match.getMetadata()
func (m *Match) GetMetadata() *NodeMetadata {
	return m.metadata
}

// NodeMetadata is the class that holds metadata (or meta objects) for use by
// JsFlume that aren't contained within the Node itself.
type NodeMetadata struct {
	compiler  jscomp.AbstractCompiler
	traversal *jscomp.NodeTraversal
}

// Java: NodeMetadata(AbstractCompiler)
func NewNodeMetadata(compiler jscomp.AbstractCompiler) *NodeMetadata {
	return newNodeMetadataWithTraversal(compiler, nil)
}

// Java: NodeMetadata(AbstractCompiler, NodeTraversal)
func newNodeMetadataWithTraversal(compiler jscomp.AbstractCompiler, traversal *jscomp.NodeTraversal) *NodeMetadata {
	return &NodeMetadata{compiler: compiler, traversal: traversal}
}

// Java: NodeMetadata.fromTraversal(NodeTraversal)
func NodeMetadata_FromTraversal(traversal *jscomp.NodeTraversal) *NodeMetadata {
	return newNodeMetadataWithTraversal(traversal.GetCompiler(), traversal)
}

// Java: NodeMetadata.getCompiler()
func (m *NodeMetadata) GetCompiler() jscomp.AbstractCompiler {
	return m.compiler
}

// Java: NodeMetadata.getTraversal()
func (m *NodeMetadata) GetTraversal() *jscomp.NodeTraversal {
	return m.traversal
}

// JsSourceMatcher is a Matcher that can take arbitrary JS source code and use
// it as a template to find matches in other source code.
//
// Java: public final class JsSourceMatcher implements Matcher
type JsSourceMatcher struct {
	matcher *jscomp.TemplateAstMatcher
}

// NewJsSourceMatcher constructs this matcher with a Function node that serves
// as the template to match all other nodes against. The body of the function
// will be used to match against.
//
// Java: JsSourceMatcher(AbstractCompiler, Node, TypeMatchingStrategy)
func NewJsSourceMatcher(
	compiler jscomp.AbstractCompiler, templateNode rhino.Node, typeMatchingStrategy jscomp.TypeMatchingStrategy) *JsSourceMatcher {
	return &JsSourceMatcher{
		matcher: jscomp.NewTemplateAstMatcher(compiler, templateNode, typeMatchingStrategy),
	}
}

// Java: JsSourceMatcher.matches(Node, NodeMetadata)
func (m *JsSourceMatcher) Matches(n rhino.Node, metadata *NodeMetadata) bool {
	return m.matcher.Matches(n)
}

// GetTemplateNodeToMatchMap returns a map from named template node strings to
// Nodes that were the equivalent matches from the last matched template.
//
// Java: JsSourceMatcher.getTemplateNodeToMatchMap()
func (m *JsSourceMatcher) GetTemplateNodeToMatchMap() *jutil.LinkedHashMap[string, rhino.Node] {
	return m.matcher.GetTemplateNodeToMatchMap()
}

// Matchers contains common Matchers useful for JsFlume refactorings.

// Matchers_Anything returns a Matcher that matches every node.
//
// Java: Matchers.anything()
func Matchers_Anything() Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		return true
	}
}

// Matchers_AllOf returns a Matcher that returns true only if all of the
// provided matchers match.
//
// Java: Matchers.allOf(Matcher...)
func Matchers_AllOf(matchers ...Matcher) Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		for _, m := range matchers {
			if !m.Matches(node, metadata) {
				return false
			}
		}
		return true
	}
}

// Matchers_AnyOf returns a Matcher that returns true if any of the provided
// matchers match.
//
// Java: Matchers.anyOf(Matcher...)
func Matchers_AnyOf(matchers ...Matcher) Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		for _, m := range matchers {
			if m.Matches(node, metadata) {
				return true
			}
		}
		return false
	}
}

// Matchers_Not returns a Matcher that matches the opposite of the provided
// matcher.
//
// Java: Matchers.not(Matcher)
func Matchers_Not(matcher Matcher) Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		return !matcher.Matches(node, metadata)
	}
}

// Matchers_Constructor returns a matcher that matches any constructor
// definitions.
//
// Java: Matchers.constructor()
func Matchers_Constructor() Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		info := node.GetJSDocInfo()
		return info != nil && info.IsConstructor()
	}
}

// Matchers_ConstructorWithName returns a matcher that matches constructor
// definitions of the specified name.
//
// Java: Matchers.constructor(String)
func Matchers_ConstructorWithName(name string) Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		info := node.GetJSDocInfo()
		if info == nil || !info.IsConstructor() {
			return false
		}
		var target rhino.Node
		if fn, ok := ast.AsFunction(node); ok {
			target = fn.NameNode().Node()
		} else if assign, ok := ast.AsAssign(node); ok {
			target = assign.Target().Node()
		} else {
			target = node.GetFirstChild()
		}
		return (target.IsGetProp() || target.IsName()) && target.MatchesQualifiedName(name)
	}
}

// Matchers_NewClass returns a Matcher that matches constructing new objects.
// This will match the NEW node of the JS Compiler AST.
//
// Java: Matchers.newClass()
func Matchers_NewClass() Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		return node.IsNew()
	}
}

// Matchers_NewClassWithName returns a Matcher that matches constructing
// objects of the provided class name.
//
// Java: Matchers.newClass(String)
func Matchers_NewClassWithName(className string) Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		if !node.IsNew() {
			return false
		}
		providedJsType := matchers_getJsType(metadata, className)
		if providedJsType == nil {
			return false
		}
		jsType := node.GetJSType()
		if jsType == nil {
			return false
		}
		jsType = jsType.RestrictByNotNullOrUndefined()
		return matchers_areTypesEquivalentIgnoringGenerics(jsType, providedJsType)
	}
}

// Matchers_FunctionCall returns a Matcher that matches any function call.
//
// Java: Matchers.functionCall()
func Matchers_FunctionCall() Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		call, ok := ast.AsCall(node)
		return ok && Matchers_PropertyAccess().Matches(call.Callee().Node(), metadata)
	}
}

// Matchers_FunctionCallWithName returns a Matcher that matches all nodes
// that are function calls that match the provided name.
//
// Java: Matchers.functionCall(String)
func Matchers_FunctionCallWithName(name string) Matcher {
	propMatcher := Matchers_PropertyAccessWithName(&name)
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		call, ok := ast.AsCall(node)
		return ok && propMatcher.Matches(call.Callee().Node(), metadata)
	}
}

// Matchers_FunctionCallWithNumArgs returns a Matcher that matches any
// function call that has the given number of arguments.
//
// Java: Matchers.functionCallWithNumArgs(int)
func Matchers_FunctionCallWithNumArgs(numArgs int) Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		call, ok := ast.AsCall(node)
		return ok && call.NumArgs() == numArgs
	}
}

// Matchers_FunctionCallWithNumArgsAndName returns a Matcher that matches any
// function call that has the given number of arguments and the given name.
//
// Java: Matchers.functionCallWithNumArgs(String, int)
func Matchers_FunctionCallWithNumArgsAndName(name string, numArgs int) Matcher {
	return Matchers_AllOf(Matchers_FunctionCallWithNumArgs(numArgs), Matchers_FunctionCallWithName(name))
}

// Java: Matchers.googRequirelike(String)
func Matchers_GoogRequirelikeWithNamespace(namespace string) Matcher {
	reqMatcher := Matchers_GoogRequirelike()
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		if !reqMatcher.Matches(node, metadata) {
			return false
		}
		arg, ok := ast.MustCall(node).FirstArg()
		if !ok {
			return false
		}
		strLit, ok := arg.AsStringLit()
		return ok && strLit.Value() == namespace
	}
}

// Java: Matchers.googRequirelike()
func Matchers_GoogRequirelike() Matcher {
	return Matchers_AnyOf(Matchers_GoogRequire(), Matchers_GoogRequireType(), Matchers_GoogForwardDeclare())
}

// Java: Matchers.googRequire()
func Matchers_GoogRequire() Matcher {
	return Matchers_FunctionCallWithName("goog.require")
}

// Java: Matchers.googRequireType()
func Matchers_GoogRequireType() Matcher {
	return Matchers_FunctionCallWithName("goog.requireType")
}

// Java: Matchers.googForwardDeclare()
func Matchers_GoogForwardDeclare() Matcher {
	return Matchers_FunctionCallWithName("goog.forwardDeclare")
}

// Java: Matchers.googModule()
func Matchers_GoogModule() Matcher {
	return Matchers_FunctionCallWithName("goog.module")
}

// Java: Matchers.googModuleOrProvide()
func Matchers_GoogModuleOrProvide() Matcher {
	return Matchers_AnyOf(Matchers_GoogModule(), Matchers_FunctionCallWithName("goog.provide"))
}

// Matchers_PropertyAccess returns a Matcher that matches any property
// access.
//
// Java: Matchers.propertyAccess()
func Matchers_PropertyAccess() Matcher {
	return Matchers_PropertyAccessWithName(nil)
}

// Matchers_PropertyAccessWithName returns a Matcher that matches nodes
// representing a GETPROP access of an object property.
//
// Java: Matchers.propertyAccess(String)
func Matchers_PropertyAccessWithName(name *string) Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		if _, ok := ast.AsGetProp(node); !ok {
			return false
		}
		if name == nil || node.MatchesQualifiedName(*name) {
			return true
		}
		if strings.Contains(*name, ".prototype.") {
			return matchers_matchesPrototypeInstanceVar(node, metadata, *name)
		}
		return false
	}
}

// Matchers_EnumDefinition returns a Matcher that matches definitions of any
// enum.
//
// Java: Matchers.enumDefinition()
func Matchers_EnumDefinition() Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		jsType := node.GetJSType()
		return jsType != nil && jsType.IsEnumType()
	}
}

// Matchers_EnumDefinitionOfType returns a Matcher that matches definitions
// of an enum of the given type.
//
// Java: Matchers.enumDefinitionOfType(String)
func Matchers_EnumDefinitionOfType(type_ string) Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		providedJsType := matchers_getJsType(metadata, type_)
		if providedJsType == nil {
			return false
		}
		providedJsType = providedJsType.RestrictByNotNullOrUndefined()

		jsType := node.GetJSType()
		return jsType != nil &&
			jsType.IsEnumType() &&
			providedJsType.Equals(jsType.ToMaybeEnumType().GetElementsType().GetPrimitiveType())
	}
}

// Matchers_AssignmentWithRhs returns a Matcher that matches an ASSIGN node
// where the RHS of the assignment matches the given rhsMatcher.
//
// Java: Matchers.assignmentWithRhs(Matcher)
func Matchers_AssignmentWithRhs(rhsMatcher Matcher) Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		assign, ok := ast.AsAssign(node)
		return ok && rhsMatcher.Matches(assign.Value().Node(), metadata)
	}
}

// Matchers_PrototypeVariableDeclaration returns a Matcher that matches a
// declaration of a variable on the prototype of a class.
//
// Java: Matchers.prototypeVariableDeclaration()
func Matchers_PrototypeVariableDeclaration() Matcher {
	return matchers_matcherForPrototypeDeclaration( /* requireFunctionType= */ false)
}

// Matchers_PrototypeMethodDeclaration returns a Matcher that matches a
// declaration of a method on the prototype of a class.
//
// Java: Matchers.prototypeMethodDeclaration()
func Matchers_PrototypeMethodDeclaration() Matcher {
	return matchers_matcherForPrototypeDeclaration( /* requireFunctionType= */ true)
}

// Matchers_JsDocType returns a Matcher that matches nodes that contain JS
// Doc that specify the `@type` annotation equivalent to the provided type.
//
// Java: Matchers.jsDocType(String)
func Matchers_JsDocType(type_ string) Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		providedJsType := matchers_getJsType(metadata, type_)
		if providedJsType == nil {
			return false
		}
		providedJsType = providedJsType.RestrictByNotNullOrUndefined()
		// The JSDoc for a variable declaration is on the VAR/CONST/LET node, but the type only
		// exists on the NAME node.
		var jsDoc *rhino.JSDocInfo
		if decl, ok := ast.AsNameDeclaration(node.GetParent()); ok {
			jsDoc = decl.Node().GetJSDocInfo()
		} else {
			jsDoc = node.GetJSDocInfo()
		}
		jsType := node.GetJSType()
		return jsDoc != nil &&
			jsDoc.HasType() &&
			jsType != nil &&
			providedJsType.Equals(jsType.RestrictByNotNullOrUndefined())
	}
}

// Matchers_ConstructorPropertyDeclaration returns a Matcher that matches
// against properties that are declared in the constructor.
//
// Java: Matchers.constructorPropertyDeclaration()
func Matchers_ConstructorPropertyDeclaration() Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		assign, ok := ast.AsAssign(node)
		if !ok {
			return false
		}
		gp, ok := ast.AsGetProp(assign.Target().Node())
		if !ok || !gp.Object().Node().IsThis() {
			return false
		}
		for !node.IsNil() && !node.IsFunction() {
			node = node.GetParent()
		}
		if _, ok := ast.AsFunction(node); ok {
			if jsDoc := jscomp.NodeUtil_GetBestJSDocInfo(node); jsDoc != nil {
				return jsDoc.IsConstructor()
			}
		}
		return false
	}
}

// Matchers_IsPrivate returns a Matcher that matches against nodes that are
// declared `@private`.
//
// Java: Matchers.isPrivate()
func Matchers_IsPrivate() Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		jsDoc := jscomp.NodeUtil_GetBestJSDocInfo(node)
		return jsDoc != nil && jsDoc.GetVisibility() == rhino.JSDocInfo_Visibility_PRIVATE
	}
}

// Java: Matchers.getJsType(NodeMetadata, String)
func matchers_getJsType(metadata *NodeMetadata, type_ string) rhino.JSType {
	return metadata.GetCompiler().GetTypeRegistry().GetGlobalType(type_)
}

// Java: Matchers.getJsType(NodeMetadata, JSTypeNative)
func matchers_getJsTypeOfNative(metadata *NodeMetadata, nativeType rhino.JSTypeNative) rhino.JSType {
	return metadata.GetCompiler().GetTypeRegistry().GetNativeType(nativeType)
}

// Java: Matchers.areTypesEquivalentIgnoringGenerics(JSType, JSType)
func matchers_areTypesEquivalentIgnoringGenerics(a rhino.JSType, b rhino.JSType) bool {
	if a.Equals(b) {
		return true
	}
	if a.IsTemplatizedType() {
		return a.ToMaybeTemplatizedType().GetReferencedType().Equals(b)
	}
	return false
}

// matchers_PROTOTYPE_SPLIT_RE is the regex of name.split(".prototype.") in
// matchesPrototypeInstanceVar: each `.` is Java's any-char-but-a-line-terminator.
var matchers_PROTOTYPE_SPLIT_RE = regexp.MustCompile(`[^\n\r\x{85}\x{2028}\x{2029}]prototype[^\n\r\x{85}\x{2028}\x{2029}]`)

func matchers_splitPrototype(s string) []string {
	matches := matchers_PROTOTYPE_SPLIT_RE.FindAllStringIndex(s, -1)
	var parts []string
	index := 0
	for _, m := range matches {
		if (m[0] == 0 || m[0] == len(s) || s[m[0]]&0xC0 != 0x80) &&
			(m[1] == 0 || m[1] == len(s) || s[m[1]]&0xC0 != 0x80) {
			parts = append(parts, s[index:m[0]])
			index = m[1]
		}
	}
	if index == 0 {
		return []string{s}
	}
	parts = append(parts, s[index:])
	for len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

// Java: Matchers.matchesPrototypeInstanceVar(Node, NodeMetadata, String)
func matchers_matchesPrototypeInstanceVar(node rhino.Node, metadata *NodeMetadata, name string) bool {
	parts := matchers_splitPrototype(name)
	className := parts[0]
	propertyName := parts[1]
	providedJsType := matchers_getJsType(metadata, className)
	if providedJsType == nil {
		return false
	}
	var jsType rhino.JSType
	if gp, ok := ast.AsGetProp(node); ok {
		jsType = gp.Object().Node().GetJSType()
	} else if node.HasChildren() {
		jsType = node.GetFirstChild().GetJSType()
	}
	if jsType == nil {
		return false
	}
	jsType = jsType.RestrictByNotNullOrUndefined()
	if !jsType.IsUnknownType() && !jsType.IsAllType() && jsType.IsSubtypeOf(providedJsType) {
		if nameNode, ok := ast.AsName(node); ok && propertyName == nameNode.Ident() {
			return true
		}
		if gp, ok := ast.AsGetProp(node); ok && propertyName == gp.Property() {
			return true
		}
	}
	return false
}

// Java: Matchers.matcherForPrototypeDeclaration(boolean)
func matchers_matcherForPrototypeDeclaration(requireFunctionType bool) Matcher {
	return func(node rhino.Node, metadata *NodeMetadata) bool {
		gp, ok := ast.AsGetProp(node)
		if !ok {
			return false
		}
		recvProp, ok := ast.AsGetProp(gp.Object().Node())
		if !ok || recvProp.Property() != "prototype" {
			return false
		}
		fnJsType := matchers_getJsTypeOfNative(metadata, rhino.JSTypeNative_FUNCTION_FUNCTION_TYPE)
		jsType := node.GetJSType()
		if jsType == nil {
			return false
		}
		if requireFunctionType {
			return jsType.CanCastTo(fnJsType)
		}
		return !jsType.CanCastTo(fnJsType)
	}
}
