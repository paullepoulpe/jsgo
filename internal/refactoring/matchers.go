// Port of com/google/javascript/refactoring/Matchers.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Matchers is the class that contains common Matchers that are useful to
// everyone.
//
// TODO(mknichel): Make sure all this code works with goog.scope.
//
// Go-only: language-adapter — the Java class has only a private
// constructor and static methods; it has no Go type. Its anonymous Matcher
// classes are matchers_matcherFunc closures.

// matchers_matcherFunc is a Matcher implemented by a closure (the anonymous
// `new Matcher() {...}` classes and Matcher lambdas of Matchers).
type matchers_matcherFunc func(node rhino.Node, metadata *NodeMetadata) bool

// Matches calls f.
func (f matchers_matcherFunc) Matches(node rhino.Node, metadata *NodeMetadata) bool {
	return f(node, metadata)
}

// Matchers_Anything returns a Matcher that matches every node.
//
// Java: Matchers.anything()
func Matchers_Anything() Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		return true
	})
}

// Matchers_AllOf returns a Matcher that returns true only if all of the
// provided matchers match.
//
// Java: Matchers.allOf(Matcher...)
func Matchers_AllOf(matchers ...Matcher) Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		for _, m := range matchers {
			if !m.Matches(node, metadata) {
				return false
			}
		}
		return true
	})
}

// Matchers_AnyOf returns a Matcher that returns true if any of the provided
// matchers match.
//
// Java: Matchers.anyOf(Matcher...)
func Matchers_AnyOf(matchers ...Matcher) Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		for _, m := range matchers {
			if m.Matches(node, metadata) {
				return true
			}
		}
		return false
	})
}

// Matchers_Not returns a Matcher that matches the opposite of the provided
// matcher.
//
// Java: Matchers.not(Matcher)
func Matchers_Not(matcher Matcher) Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		return !matcher.Matches(node, metadata)
	})
}

// Matchers_Constructor returns a matcher that matches any constructor
// definitions.
//
// Java: Matchers.constructor()
func Matchers_Constructor() Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		info := node.GetJSDocInfo()
		return info != nil && info.IsConstructor()
	})
}

// Matchers_ConstructorWithName returns a matcher that matches constructor
// definitions of the specified name.
//
// name is the name of the class constructor to match.
//
// Java: Matchers.constructor(String)
func Matchers_ConstructorWithName(name string) Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		info := node.GetJSDocInfo()
		if info != nil && info.IsConstructor() {
			firstChild := node.GetFirstChild()
			// TODO(mknichel): Make sure this works with the following cases:
			// ns = {
			//   /** @constructor */
			//   name: function() {}
			// }
			if (firstChild.IsGetProp() || firstChild.IsName()) &&
				firstChild.MatchesQualifiedName(name) {
				return true
			}
		}
		return false
	})
}

// Matchers_NewClass returns a Matcher that matches constructing new objects.
// This will match the NEW node of the JS Compiler AST.
//
// Java: Matchers.newClass()
func Matchers_NewClass() Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		return node.IsNew()
	})
}

// Matchers_NewClassWithName returns a Matcher that matches constructing
// objects of the provided class name. This will match the NEW node of the
// JS Compiler AST.
//
// className is the name of the class to return matching NEW nodes.
//
// Java: Matchers.newClass(String)
func Matchers_NewClassWithName(className string) Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
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
	})
}

// Matchers_FunctionCall returns a Matcher that matches any function call.
//
// Java: Matchers.functionCall()
func Matchers_FunctionCall() Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		return node.IsCall() && Matchers_PropertyAccess().Matches(node.GetFirstChild(), metadata)
	})
}

// Matchers_FunctionCallWithName returns a Matcher that matches all nodes
// that are function calls that match the provided name.
//
// name is the name of the function to match. For non-static functions, this
// must be the fully qualified name that includes the type of the object.
// For instance: ns.AppContext.prototype.get will match appContext.get and
// this.get when called from the AppContext class.
//
// Java: Matchers.functionCall(String)
func Matchers_FunctionCallWithName(name string) Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		// TODO(mknichel): Handle the case when functions are applied through .call or .apply.
		return node.IsCall() && Matchers_PropertyAccessWithName(&name).Matches(node.GetFirstChild(), metadata)
	})
}

// Matchers_FunctionCallWithNumArgs returns a Matcher that matches any
// function call that has the given number of arguments.
//
// Java: Matchers.functionCallWithNumArgs(int)
func Matchers_FunctionCallWithNumArgs(numArgs int) Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		return node.IsCall() && node.HasXChildren(numArgs+1)
	})
}

// Matchers_FunctionCallWithNumArgsAndName returns a Matcher that matches any
// function call that has the given number of arguments and the given name.
//
// name is the name of the function to match. For non-static functions, this
// must be the fully qualified name that includes the type of the object.
// For instance: ns.AppContext.prototype.get will match appContext.get and
// this.get when called from the AppContext class.
//
// Java: Matchers.functionCallWithNumArgs(String, int)
func Matchers_FunctionCallWithNumArgsAndName(name string, numArgs int) Matcher {
	return Matchers_AllOf(Matchers_FunctionCallWithNumArgs(numArgs), Matchers_FunctionCallWithName(name))
}

// Java: Matchers.googRequirelike(String)
func Matchers_GoogRequirelikeWithNamespace(namespace string) Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		return Matchers_GoogRequirelike().Matches(node, metadata) &&
			node.GetSecondChild().IsStringLit() &&
			node.GetSecondChild().GetString() == namespace
	})
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
// name is the name of the property to match. For non-static properties,
// this must be the fully qualified name that includes the type of the
// object. For instance: ns.AppContext.prototype.root will match
// appContext.root and this.root when accessed from the AppContext.
//
// Java: Matchers.propertyAccess(String)
func Matchers_PropertyAccessWithName(name *string) Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		if node.IsGetProp() {
			if name == nil {
				return true
			}
			if node.MatchesQualifiedName(*name) {
				return true
			} else if strings.Contains(*name, ".prototype.") {
				return matchers_matchesPrototypeInstanceVar(node, metadata, *name)
			}
		}
		return false
	})
}

// Matchers_EnumDefinition returns a Matcher that matches definitions of any
// enum.
//
// Java: Matchers.enumDefinition()
func Matchers_EnumDefinition() Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		jsType := node.GetJSType()
		return jsType != nil && jsType.IsEnumType()
	})
}

// Matchers_EnumDefinitionOfType returns a Matcher that matches definitions
// of an enum of the given type.
//
// Java: Matchers.enumDefinitionOfType(String)
func Matchers_EnumDefinitionOfType(type_ string) Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		providedJsType := matchers_getJsType(metadata, type_)
		if providedJsType == nil {
			return false
		}
		providedJsType = providedJsType.RestrictByNotNullOrUndefined()

		jsType := node.GetJSType()
		return jsType != nil &&
			jsType.IsEnumType() &&
			providedJsType.Equals(jsType.ToMaybeEnumType().GetElementsType().GetPrimitiveType())
	})
}

// Matchers_AssignmentWithRhs returns a Matcher that matches an ASSIGN node
// where the RHS of the assignment matches the given rhsMatcher.
//
// Java: Matchers.assignmentWithRhs(Matcher)
func Matchers_AssignmentWithRhs(rhsMatcher Matcher) Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		return node.IsAssign() && rhsMatcher.Matches(node.GetLastChild(), metadata)
	})
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
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		providedJsType := matchers_getJsType(metadata, type_)
		if providedJsType == nil {
			return false
		}
		providedJsType = providedJsType.RestrictByNotNullOrUndefined()
		// The JSDoc for a variable declaration is on the VAR/CONST/LET node, but the type only
		// exists on the NAME node.
		// TODO(mknichel): Make NodeUtil.getBestJSDoc public and use that.
		var jsDoc *rhino.JSDocInfo
		if jscomp.NodeUtil_IsNameDeclaration(node.GetParent()) {
			jsDoc = node.GetParent().GetJSDocInfo()
		} else {
			jsDoc = node.GetJSDocInfo()
		}
		jsType := node.GetJSType()
		return jsDoc != nil &&
			jsDoc.HasType() &&
			jsType != nil &&
			providedJsType.Equals(jsType.RestrictByNotNullOrUndefined())
	})
}

// Matchers_ConstructorPropertyDeclaration returns a Matcher that matches
// against properties that are declared in the constructor.
//
// Java: Matchers.constructorPropertyDeclaration()
func Matchers_ConstructorPropertyDeclaration() Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		// This will match against code that looks like:
		// /** @constructor */
		// function constructor() {
		//   this.variable = 3;
		// }
		if !node.IsAssign() ||
			!node.GetFirstChild().IsGetProp() ||
			!node.GetFirstFirstChild().IsThis() {
			return false
		}
		for !node.IsNil() && !node.IsFunction() {
			node = node.GetParent()
		}
		if !node.IsNil() && node.IsFunction() {
			jsDoc := jscomp.NodeUtil_GetBestJSDocInfo(node)
			if jsDoc != nil {
				return jsDoc.IsConstructor()
			}
		}
		return false
	})
}

// Matchers_IsPrivate returns a Matcher that matches against nodes that are
// declared `@private`.
//
// Java: Matchers.isPrivate()
func Matchers_IsPrivate() Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		jsDoc := jscomp.NodeUtil_GetBestJSDocInfo(node)
		if jsDoc != nil {
			return jsDoc.GetVisibility() == rhino.JSDocInfo_Visibility_PRIVATE
		}
		return false
	})
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
	equivalent := a.Equals(b)
	if equivalent {
		return true
	}
	if a.IsTemplatizedType() {
		return a.ToMaybeTemplatizedType().GetReferencedType().Equals(b)
	}
	return false
}

// matchers_PROTOTYPE_SPLIT_RE is the regex of name.split(".prototype.") in
// matchesPrototypeInstanceVar: each `.` is Java's any-char-but-a-line-
// terminator.
//
// Go-only: language-adapter — java.util.regex `.prototype.` written for
// RE2 (both `.` matching what Java's does), used with jutil.StringSplitRegexp.
var matchers_PROTOTYPE_SPLIT_RE = regexp.MustCompile(`[^\n\r\x{85}\x{2028}\x{2029}]prototype[^\n\r\x{85}\x{2028}\x{2029}]`)

// matchers_matchesPrototypeInstanceVar checks to see if the node represents
// an access of an instance variable on an object given a prototype
// declaration of an object. For instance, ns.AppContext.prototype.get will
// match appContext.get or this.get when accessed from within the AppContext
// object.
//
// Java: Matchers.matchesPrototypeInstanceVar(Node, NodeMetadata, String)
func matchers_matchesPrototypeInstanceVar(node rhino.Node, metadata *NodeMetadata,
	name string) bool {
	parts := jutil.StringSplitRegexp(name, matchers_PROTOTYPE_SPLIT_RE)
	className := parts[0]
	propertyName := parts[1]
	providedJsType := matchers_getJsType(metadata, className)
	if providedJsType == nil {
		return false
	}
	var jsType rhino.JSType
	if node.HasChildren() {
		jsType = node.GetFirstChild().GetJSType()
	}
	if jsType == nil {
		return false
	}
	jsType = jsType.RestrictByNotNullOrUndefined()
	if !jsType.IsUnknownType() &&
		!jsType.IsAllType() &&
		jsType.IsSubtypeOf(providedJsType) {
		if node.IsName() && propertyName == node.GetString() {
			return true
		} else if node.IsGetProp() && propertyName == node.GetString() {
			return true
		}
	}
	return false
}

// Java: Matchers.matcherForPrototypeDeclaration(boolean)
func matchers_matcherForPrototypeDeclaration(requireFunctionType bool) Matcher {
	return matchers_matcherFunc(func(node rhino.Node, metadata *NodeMetadata) bool {
		// TODO(mknichel): Figure out which node is the best to return for this
		// function: the GETPROP node, or the ASSIGN node when the property is
		// being assigned to.
		// TODO(mknichel): Support matching:
		// foo.prototype = {
		//   bar: 1
		// };
		firstChild := node.GetFirstChild()
		if node.IsGetProp() &&
			firstChild.IsGetProp() &&
			"prototype" == firstChild.GetString() {
			fnJsType := matchers_getJsTypeOfNative(metadata, rhino.JSTypeNative_FUNCTION_FUNCTION_TYPE)
			jsType := node.GetJSType()
			if jsType == nil {
				return false
			} else if requireFunctionType {
				return jsType.CanCastTo(fnJsType)
			} else {
				return !jsType.CanCastTo(fnJsType)
			}
		}
		return false
	})
}

// TODO(mknichel): Add matchers for:
// - Constructor with argument types
// - Function call with argument types
// - Function definitions.
// - Property definitions, references
// - IsStatic
// - JsDocMatcher
