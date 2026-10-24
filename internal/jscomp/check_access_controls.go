// Port of com/google/javascript/jscomp/CheckAccessControls.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/AutoBuilder_CheckAccessControls_PropertyReference_Builder.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckAccessControls is a compiler pass that checks that the programmer has
// obeyed all the access control restrictions indicated by JSDoc annotations,
// like @private and @deprecated.
//
// Because access control restrictions are attached to type information, this
// pass must run after TypeInference, and InferJSDocInfo.
type CheckAccessControls struct {
	compiler     AbstractCompiler
	typeRegistry *rhino.JSTypeRegistry

	// State about the current traversal.
	deprecationDepth int
	// NOTE: LinkedList is almost always the wrong choice, but in this case we have at most a small
	// handful of elements, it provides the smoothest API (push, pop, and a peek that doesn't throw
	// on empty), and (unlike ArrayDeque) is null-permissive. No other option meets all these needs.
	//
	// Go-only: language-adapter — the Deque is a slice whose last element is
	// the head (push/pop at the end); iteration runs from the end.
	currentClassStack []rhino.ObjectType
	// Go-only: language-adapter — the Guava Table<JSType, String,
	// ConstantDeclaration> of Tables.newCustomTable(new LinkedHashMap<>(),
	// LinkedHashMap::new): the row keys hash and compare with
	// JSType.hashCode()/equals() (ObjectLinkedHashMap); the table is only
	// looked up, never iterated.
	constPropertyInits *jutil.ObjectLinkedHashMap[rhino.JSType, *jutil.LinkedHashMap[string, *CheckAccessControls_ConstantDeclaration]]

	defaultVisibilityForFiles *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility]

	scratchPropRef           CheckAccessControls_PropertyReference
	lastHoistScopeRoot       *rhino.Node
	lastHoistScopeDeprecated bool
	lastVarNode              *rhino.Node
	lastVar                  *Var
}

var checkAccessControlsActiveToken = func() [256]bool {
	var t [256]bool
	t[rhino.Token_NAME] = true
	t[rhino.Token_GETPROP] = true
	t[rhino.Token_STRING_KEY] = true
	t[rhino.Token_GETTER_DEF] = true
	t[rhino.Token_SETTER_DEF] = true
	t[rhino.Token_MEMBER_FUNCTION_DEF] = true
	t[rhino.Token_MEMBER_FIELD_DEF] = true
	t[rhino.Token_FUNCTION] = true
	t[rhino.Token_CLASS] = true
	return t
}()

// Java: CheckAccessControls.DEPRECATED_NAME
var CheckAccessControls_DEPRECATED_NAME = DiagnosticType_Disabled("JSC_DEPRECATED_VAR", "Variable {0} has been deprecated.")

// Java: CheckAccessControls.DEPRECATED_NAME_REASON
var CheckAccessControls_DEPRECATED_NAME_REASON = DiagnosticType_Disabled("JSC_DEPRECATED_VAR_REASON", "Variable {0} has been deprecated: {1}")

// Java: CheckAccessControls.DEPRECATED_PROP
var CheckAccessControls_DEPRECATED_PROP = DiagnosticType_Disabled(
	"JSC_DEPRECATED_PROP", "Property {0} of type {1} has been deprecated.")

// Java: CheckAccessControls.DEPRECATED_PROP_REASON
var CheckAccessControls_DEPRECATED_PROP_REASON = DiagnosticType_Disabled(
	"JSC_DEPRECATED_PROP_REASON", "Property {0} of type {1} has been deprecated: {2}")

// Java: CheckAccessControls.DEPRECATED_CLASS
var CheckAccessControls_DEPRECATED_CLASS = DiagnosticType_Disabled("JSC_DEPRECATED_CLASS", "Class {0} has been deprecated.")

// Java: CheckAccessControls.DEPRECATED_CLASS_REASON
var CheckAccessControls_DEPRECATED_CLASS_REASON = DiagnosticType_Disabled("JSC_DEPRECATED_CLASS_REASON", "Class {0} has been deprecated: {1}")

// Java: CheckAccessControls.BAD_PACKAGE_PROPERTY_ACCESS
var CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS = DiagnosticType_Error(
	"JSC_BAD_PACKAGE_PROPERTY_ACCESS",
	"Access to package-private property {0} of {1} not allowed here.")

// Java: CheckAccessControls.BAD_PRIVATE_GLOBAL_ACCESS
var CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS = DiagnosticType_Error(
	"JSC_BAD_PRIVATE_GLOBAL_ACCESS",
	"Access to private variable {0} not allowed outside file {1}.")

// Java: CheckAccessControls.BAD_PRIVATE_PROPERTY_ACCESS
var CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS = DiagnosticType_Warning(
	"JSC_BAD_PRIVATE_PROPERTY_ACCESS",
	"Access to private property {0} of {1} not allowed here.")

// Java: CheckAccessControls.BAD_PROTECTED_PROPERTY_ACCESS
var CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS = DiagnosticType_Warning(
	"JSC_BAD_PROTECTED_PROPERTY_ACCESS",
	"Access to protected property {0} of {1} not allowed here.")

// Java: CheckAccessControls.BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY
var CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY = DiagnosticType_Error(
	"JSC_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY",
	"Overridden property {0} in file with fileoverview visibility {1}"+
		" must explicitly redeclare superclass visibility")

// Java: CheckAccessControls.PRIVATE_OVERRIDE
var CheckAccessControls_PRIVATE_OVERRIDE = DiagnosticType_Warning("JSC_PRIVATE_OVERRIDE", "Overriding private property of {0}.")

// Java: CheckAccessControls.EXTEND_FINAL_CLASS
var CheckAccessControls_EXTEND_FINAL_CLASS = DiagnosticType_Error(
	"JSC_EXTEND_FINAL_CLASS", "{0} is not allowed to extend final class {1}.")

// Java: CheckAccessControls.VISIBILITY_MISMATCH
var CheckAccessControls_VISIBILITY_MISMATCH = DiagnosticType_Warning(
	"JSC_VISIBILITY_MISMATCH", "Overriding {0} property of {1} with {2} property.")

// Java: CheckAccessControls.CONST_PROPERTY_REASSIGNED_VALUE
var CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE = DiagnosticType_Warning(
	"JSC_CONSTANT_PROPERTY_REASSIGNED_VALUE",
	base.JSCompStrings_Lines(
		"constant property {0} assigned a value more than once", //
		"Initialized at {1}"))

// Java: CheckAccessControls.FINAL_PROPERTY_OVERRIDDEN
var CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN = DiagnosticType_Warning(
	"JSC_FINAL_PROPERTY_OVERRIDDEN",
	base.JSCompStrings_Lines(
		"@final method or property {0} overridden", //
		"Initialized at {1}"))

// Java: CheckAccessControls.CONST_PROPERTY_DELETED
var CheckAccessControls_CONST_PROPERTY_DELETED = DiagnosticType_Warning(
	"JSC_CONSTANT_PROPERTY_DELETED", "constant property {0} cannot be deleted")

// CheckAccessControls_Constancy distinguishes between different kinds of
// "constant" JSDoc to provide more useful error messages.
//
// Java: CheckAccessControls.Constancy
type CheckAccessControls_Constancy int32

const (
	CheckAccessControls_Constancy_FINAL          CheckAccessControls_Constancy = iota // @final
	CheckAccessControls_Constancy_OTHER_CONSTANT                                      // e.g. @const, @define, or @desc but not @final
	CheckAccessControls_Constancy_MUTABLE
)

var checkAccessControls_Constancy_names = [...]string{"FINAL", "OTHER_CONSTANT", "MUTABLE"}

// CheckAccessControls_Constancy_values is Constancy.values().
func CheckAccessControls_Constancy_values() []CheckAccessControls_Constancy {
	return []CheckAccessControls_Constancy{
		CheckAccessControls_Constancy_FINAL, CheckAccessControls_Constancy_OTHER_CONSTANT, CheckAccessControls_Constancy_MUTABLE,
	}
}

// String is Enum.name().
func (c CheckAccessControls_Constancy) String() string { return checkAccessControls_Constancy_names[c] }

// Ordinal is Enum.ordinal().
func (c CheckAccessControls_Constancy) Ordinal() int { return int(c) }

// Java: CheckAccessControls.ConstantDeclaration
type CheckAccessControls_ConstantDeclaration struct {
	node       *rhino.Node
	annotation CheckAccessControls_Constancy
}

// Java: CheckAccessControls.ConstantDeclaration(Node, Constancy)
func NewCheckAccessControls_ConstantDeclaration(node *rhino.Node, annotation CheckAccessControls_Constancy) *CheckAccessControls_ConstantDeclaration {
	d := &CheckAccessControls_ConstantDeclaration{}
	d.node = node
	d.annotation = annotation
	return d
}

// Java: CheckAccessControls(AbstractCompiler)
func NewCheckAccessControls(compiler AbstractCompiler) *CheckAccessControls {
	c := &CheckAccessControls{}
	c.deprecationDepth = 0
	c.currentClassStack = nil
	c.constPropertyInits = jutil.NewObjectLinkedHashMap[rhino.JSType, *jutil.LinkedHashMap[string, *CheckAccessControls_ConstantDeclaration]]()
	c.compiler = compiler
	c.typeRegistry = compiler.GetTypeRegistry()
	return c
}

// Java: CheckAccessControls.process(Node, Node)
func (c *CheckAccessControls) Process(externs *rhino.Node, root *rhino.Node) {
	collectPass := NewCollectFileOverviewVisibility(c.compiler)
	collectPass.Process(externs, root)
	c.defaultVisibilityForFiles = collectPass.GetFileOverviewVisibilityMap()

	sc := NewSyntacticScopeCreator(c.compiler)
	defer func() {
		c.lastVarNode = nil
		c.lastVar = nil
	}()
	b := NodeTraversal_Builder{
		compiler:     c.compiler,
		callback:     c,
		scopeCreator: sc,
	}

	c.lastHoistScopeRoot = nil
	c.lastVarNode = nil
	c.lastVar = nil
	b.Traverse(externs)
	c.lastHoistScopeRoot = nil
	c.lastVarNode = nil
	c.lastVar = nil
	b.Traverse(root)
}

func (c *CheckAccessControls) lookupVar(scope *Scope, name string) *Var {
	if name == "" {
		return nil
	}
	return scope.GetVar(name)
}

// Java: CheckAccessControls.enterAccessControlScope(Node)
func (c *CheckAccessControls) EnterAccessControlScope(root *rhino.Node) {
	scopeType := c.BestInstanceTypeForMethodOrCtor(root)

	if CheckAccessControls_IsMarkedDeprecated(root) {
		c.deprecationDepth++
	}
	c.currentClassStack = append(c.currentClassStack, scopeType)
}

// Java: CheckAccessControls.exitAccessControlScope(Node)
func (c *CheckAccessControls) ExitAccessControlScope(root *rhino.Node) {
	if CheckAccessControls_IsMarkedDeprecated(root) {
		c.deprecationDepth--
	}
	// Deque.pop(): NoSuchElementException when empty.
	if len(c.currentClassStack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	c.currentClassStack[len(c.currentClassStack)-1] = nil
	c.currentClassStack = c.currentClassStack[:len(c.currentClassStack)-1]
}

// CheckAccessControls_PrimaryAccessControlScopeRootFor maps node to the
// primary root of an access-control scope if it is some root, or null if it
// is a non-root of the scope.
//
// We define access-control scopes differently from Scopes because of
// mismatches between ECMAScript scoping and our AST structure (e.g. the
// `extends` clause of a CLASS). NodeTraversal is designed to walk the AST in
// ECMAScript scope order, which is not the pre-order traversal that we would
// prefer here. This requires us to treat access-control scopes as a forest
// with a primary root.
//
// Each access-control scope corresponds to some "owner" JavaScript type which
// is used when computing access-controls. At this time, each also corresponds
// to an AST subtree.
//
// Known mismatches:
//
//   - CLASS extends clause: secondary root of the scope defined by the CLASS.
//
// Java: CheckAccessControls.primaryAccessControlScopeRootFor(Node)
func CheckAccessControls_PrimaryAccessControlScopeRootFor(node *rhino.Node) *rhino.Node {
	if CheckAccessControls_IsExtendsTarget(node) {
		return node.GetParent()
	} else if CheckAccessControls_IsFunctionOrClass(node) {
		return node
	} else {
		return nil
	}
}

// BestInstanceTypeForMethodOrCtor returns the instance object type that best
// represents a method or constructor definition, or null if there is no
// representative type.
//
//   - Prototype methods => The instance type having that prototype
//   - Instance methods => The type the method is attached to
//   - Constructors => The type that constructor instantiates
//   - Object-literal members => The object-literal type
//
// Java: CheckAccessControls.bestInstanceTypeForMethodOrCtor(Node)
func (c *CheckAccessControls) BestInstanceTypeForMethodOrCtor(n *rhino.Node) rhino.ObjectType {
	jutil.CheckStateMsg1(CheckAccessControls_IsFunctionOrClass(n), "%s", n)
	parent := n.GetParent()

	// We need to handle declaration syntaxes separately in a way that we can't determine based on
	// the type of just one node.
	// TODO(nickreid): Determine if these can be replaced with FUNCTION and CLASS cases below.
	if NodeUtil_IsFunctionDeclaration(n) ||
		// isClassDeclaration returns false for many instances that we need to handle here.
		n.IsClass() {
		return CheckAccessControls_InstanceTypeFor(n.GetJSType())
	}

	// All the remaining cases can be isolated based on `parent`.
	switch parent.GetToken() {
	case rhino.Token_NAME:
		return CheckAccessControls_InstanceTypeFor(n.GetJSType())
	case rhino.Token_ASSIGN:
		lValue := parent.GetFirstChild()
		if NodeUtil_IsNormalGet(lValue) {
			// We have an assignment of the form `a.b = ...`.
			lValueType := lValue.GetJSType()
			if lValueType != nil && (lValueType.IsConstructor() || lValueType.IsInterface()) {
				// Case `a.B = ...`
				return CheckAccessControls_InstanceTypeFor(lValueType)
			} else if NodeUtil_IsPrototypeProperty(lValue) {
				// Case `a.B.prototype = ...`
				return CheckAccessControls_InstanceTypeFor(NodeUtil_GetPrototypeClassName(lValue).GetJSType())
			} else {
				// Case `a.b = ...`
				return CheckAccessControls_InstanceTypeFor(lValue.GetFirstChild().GetJSType())
			}
		} else {
			// We have an assignment of the form "a = ...", so pull the type off the "a".
			return CheckAccessControls_InstanceTypeFor(lValue.GetJSType())
		}
	case rhino.Token_STRING_KEY,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_COMPUTED_PROP:
		grandparent := parent.GetParent()
		greatGrandparent := grandparent.GetParent()

		if grandparent.IsObjectLit() {
			if grandparent.GetJSType().IsFunctionPrototypeType() {
				// Case: `grandparent` is an object-literal prototype.
				// Example: `Foo.prototype = { a: function() {} };` where `parent` is "a".
				return CheckAccessControls_InstanceTypeFor(grandparent.GetJSType())
			}
			return nil
		} else if greatGrandparent.IsClass() {
			// Case: `n` is a class member definition.
			// Example: `class Foo { a() {} }` where `parent` is "a".
			return CheckAccessControls_InstanceTypeFor(greatGrandparent.GetJSType())
		} else {
			// This would indicate the AST is malformed.
			panic(jutil.NewAssertionError(jutil.ToJavaString(greatGrandparent)))
		}
	default:
		return nil
	}
}

// CheckAccessControls_InstanceTypeFor returns the type that best represents
// the instance type for type.
//
//   - Prototype type => The instance type having that prototype
//   - Instance type => The type
//   - Constructor type => The type that constructor instantiates
//   - Object-literal type => The type
//
// Java: CheckAccessControls.instanceTypeFor(JSType)
func CheckAccessControls_InstanceTypeFor(type_ rhino.JSType) rhino.ObjectType {
	if type_ == nil {
		return nil
	} else if type_.IsUnionType() {
		return nil // A union has no meaningful instance type.
	} else if type_.IsInstanceType() || type_.IsUnknownType() {
		return type_.ToMaybeObjectType()
	} else if type_.IsConstructor() || type_.IsInterface() {
		return type_.ToMaybeFunctionType().GetInstanceType()
	} else if type_.IsFunctionType() {
		return nil // Functions that aren't ctors or interfaces have no instance type.
	} else if type_.IsFunctionPrototypeType() {
		return CheckAccessControls_InstanceTypeFor(type_.ToMaybeObjectType().GetOwnerFunction())
	}

	return type_.ToMaybeObjectType()
}

// Java: CheckAccessControls.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckAccessControls) ShouldTraverse(traversal *NodeTraversal, node *rhino.Node, parent *rhino.Node) bool {
	accessControlRoot := CheckAccessControls_PrimaryAccessControlScopeRootFor(node)
	if accessControlRoot != nil {
		c.EnterAccessControlScope(accessControlRoot)
	}

	return true
}

// Java: CheckAccessControls.visit(NodeTraversal, Node, Node)
func (c *CheckAccessControls) Visit(traversal *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	tok := node.GetToken()
	if !checkAccessControlsActiveToken[tok] {
		if t := node.GetJSType(); t == nil || !t.IsFunctionType() {
			if parent != nil && parent.IsClass() && parent.GetFirstChild().GetNext() == node {
				c.ExitAccessControlScope(parent)
			}
			return
		}
	} else if tok == rhino.Token_NAME {
		if t := node.GetJSType(); t == nil || !t.IsFunctionType() || !t.ToMaybeFunctionType().IsConstructor() {
			if name := node.GetString(); name != "" {
				scope := traversal.GetScope()
				if var_ := c.lookupVar(scope, name); var_ != nil {
					if docInfo := var_.GetJSDocInfo(); docInfo != nil && docInfo.IsDeprecated() {
						if c.ShouldEmitDeprecationWarning(traversal, node) {
							if docInfo.GetDeprecationReason() != nil {
								c.compiler.Report(
									JSError_MakeWithNode(node, CheckAccessControls_DEPRECATED_NAME_REASON, name, *docInfo.GetDeprecationReason()))
							} else {
								c.compiler.Report(JSError_MakeWithNode(node, CheckAccessControls_DEPRECATED_NAME, name))
							}
						}
					}
					c.checkNameVisibilityWithVar(var_, node)
				}
			}
			if parent != nil && parent.IsClass() && parent.GetFirstChild().GetNext() == node {
				c.ExitAccessControlScope(parent)
			}
			return
		}
	}

	identifierBehaviour := CheckAccessControls_IdentifierBehaviour_Select(node)
	propRef := c.CreatePropertyReference(node)

	c.CheckDeprecation(node, propRef, identifierBehaviour, traversal)
	var scope *Scope
	if tok == rhino.Token_NAME && identifierBehaviour != CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE {
		scope = traversal.GetScope()
	}
	c.CheckVisibility(node, propRef, identifierBehaviour, scope)
	c.CheckConstantProperty(propRef, identifierBehaviour)

	c.CheckFinalClassOverrides(node)

	accessControlRoot := CheckAccessControls_PrimaryAccessControlScopeRootFor(node)
	if accessControlRoot != nil {
		c.ExitAccessControlScope(accessControlRoot)
	}
}

// Java: CheckAccessControls.checkDeprecation(Node, PropertyReference, IdentifierBehaviour, NodeTraversal)
func (c *CheckAccessControls) CheckDeprecation(
	node *rhino.Node,
	propRef *CheckAccessControls_PropertyReference,
	identifierBehaviour CheckAccessControls_IdentifierBehaviour,
	traversal *NodeTraversal) {

	switch identifierBehaviour {
	case CheckAccessControls_IdentifierBehaviour_ES5_CLASS_INVOCATION,
		CheckAccessControls_IdentifierBehaviour_ES6_CLASS_INVOCATION,
		CheckAccessControls_IdentifierBehaviour_ES6_CLASS_NAMESPACE:
		// At these usages, treat the deprecation applied to type-declaration as referring to the
		// type, not the identifier (e.g. "the use of class `Foo` is deprecated").
		c.CheckTypeDeprecation(traversal, node)
	case CheckAccessControls_IdentifierBehaviour_NON_CONSTRUCTOR:
		// For all identifiers that are not constructors, deprecation refers to the identifier
		// (e.g.
		// "the use of variable `x` is deprecated").
		c.CheckNameDeprecation(traversal, node)
	default:
	}

	if propRef != nil && identifierBehaviour != CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE {
		c.CheckPropertyDeprecation(traversal, propRef)
	}
}

// Java: CheckAccessControls.checkVisibility(Node, PropertyReference, IdentifierBehaviour, Scope)
func (c *CheckAccessControls) CheckVisibility(
	node *rhino.Node,
	propRef *CheckAccessControls_PropertyReference,
	identifierBehaviour CheckAccessControls_IdentifierBehaviour,
	scope *Scope) {
	if identifierBehaviour == CheckAccessControls_IdentifierBehaviour_ES6_CLASS_INVOCATION {
		c.CheckEs6ConstructorInvocationVisibility(node)
	}

	if identifierBehaviour != CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE {
		c.CheckNameVisibility(scope, node)
	}

	if propRef != nil && identifierBehaviour != CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE {
		c.CheckPropertyVisibility(propRef)
	}
}

// CheckTypeDeprecation reports deprecation issue with regard to a type usage.
//
// Precondition: n has a constructor JSType.
//
// Java: CheckAccessControls.checkTypeDeprecation(NodeTraversal, Node)
func (c *CheckAccessControls) CheckTypeDeprecation(t *NodeTraversal, n *rhino.Node) {
	if !c.ShouldEmitDeprecationWarning(t, n) {
		return
	}

	instanceType := n.GetJSType().ToMaybeFunctionType().GetInstanceType()

	deprecationInfo := CheckAccessControls_GetTypeDeprecationInfo(instanceType)
	if deprecationInfo == nil {
		return
	}

	message := CheckAccessControls_DEPRECATED_CLASS_REASON
	if *deprecationInfo == "" {
		message = CheckAccessControls_DEPRECATED_CLASS
	}
	c.compiler.Report(JSError_MakeWithNode(n, message, instanceType.String(), *deprecationInfo))
}

// CheckNameDeprecation checks the given NAME node to ensure that access
// restrictions are obeyed.
//
// Java: CheckAccessControls.checkNameDeprecation(NodeTraversal, Node)
func (c *CheckAccessControls) CheckNameDeprecation(t *NodeTraversal, n *rhino.Node) {
	if !n.IsName() {
		return
	}

	var_ := c.lookupVar(t.GetScope(), n.GetString())
	c.lastVarNode = n
	c.lastVar = var_
	if var_ == nil {
		return
	}
	docInfo := var_.GetJSDocInfo()
	if docInfo == nil || !docInfo.IsDeprecated() {
		return
	}

	if !c.ShouldEmitDeprecationWarning(t, n) {
		return
	}

	if docInfo.GetDeprecationReason() != nil {
		c.compiler.Report(
			JSError_MakeWithNode(n, CheckAccessControls_DEPRECATED_NAME_REASON, n.GetString(), *docInfo.GetDeprecationReason()))
	} else {
		c.compiler.Report(JSError_MakeWithNode(n, CheckAccessControls_DEPRECATED_NAME, n.GetString()))
	}
}

// CheckPropertyDeprecation checks the given GETPROP node to ensure that
// access restrictions are obeyed.
//
// Java: CheckAccessControls.checkPropertyDeprecation(NodeTraversal, PropertyReference)
func (c *CheckAccessControls) CheckPropertyDeprecation(t *NodeTraversal, propRef *CheckAccessControls_PropertyReference) {
	if !c.ShouldEmitDeprecationWarningForPropRef(t, propRef) {
		return
	}

	// Don't bother checking constructors.
	if propRef.SourceNode().GetParent().IsNew() {
		return
	}

	objectType := CheckAccessControls_CastToObject(CheckAccessControls_Dereference(propRef.ReceiverType()))
	propertyName := propRef.Name()

	if objectType != nil {
		deprecationInfo := CheckAccessControls_GetPropertyDeprecationInfo(objectType, propertyName)
		if deprecationInfo != nil {
			if *deprecationInfo != "" {
				c.compiler.Report(
					JSError_MakeWithNode(
						propRef.SourceNode(),
						CheckAccessControls_DEPRECATED_PROP_REASON,
						propertyName,
						propRef.GetReadableTypeNameOrDefault(),
						*deprecationInfo))
			} else {
				c.compiler.Report(
					JSError_MakeWithNode(
						propRef.SourceNode(),
						CheckAccessControls_DEPRECATED_PROP,
						propertyName,
						propRef.GetReadableTypeNameOrDefault()))
			}
		}
	}
}

// CheckNameVisibility reports an error if the given name is not visible in
// the current context.
//
// scope is the current scope; name the name node.
//
// Java: CheckAccessControls.checkNameVisibility(Scope, Node)
func (c *CheckAccessControls) CheckNameVisibility(scope *Scope, name *rhino.Node) {
	if !name.IsName() {
		return
	}

	var var_ *Var
	if c.lastVarNode == name {
		var_ = c.lastVar
	} else {
		var_ = c.lookupVar(scope, name.GetString())
	}
	if var_ == nil {
		return
	}
	c.checkNameVisibilityWithVar(var_, name)
}

func (c *CheckAccessControls) checkNameVisibilityWithVar(var_ *Var, name *rhino.Node) {
	var v rhino.JSDocInfo_Visibility
	if c.defaultVisibilityForFiles.IsEmpty() {
		if jsDocInfo := var_.GetJSDocInfo(); jsDocInfo != nil {
			v = jsDocInfo.GetVisibility()
		} else {
			return
		}
	} else {
		v = AccessControlUtils_GetEffectiveNameVisibility(name, var_, c.defaultVisibilityForFiles)
	}

	switch v {
	case rhino.JSDocInfo_Visibility_PACKAGE:
		if !c.IsPackageAccessAllowed(var_, name) {
			c.compiler.Report(
				JSError_MakeWithNode(
					name,
					CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS,
					name.GetString(),
					var_.GetSourceFile().GetName()))
		}
	case rhino.JSDocInfo_Visibility_PRIVATE:
		if !CheckAccessControls_IsPrivateAccessAllowed(var_, name) {
			c.compiler.Report(
				JSError_MakeWithNode(
					name,
					CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS,
					name.GetString(),
					var_.GetSourceFile().GetName()))
		}
	default:
		// Nothing to do for PUBLIC and PROTECTED
		// (which is irrelevant for names).
	}
}

// Java: CheckAccessControls.isPrivateAccessAllowed(Var, Node)
func CheckAccessControls_IsPrivateAccessAllowed(var_ *Var, name *rhino.Node) bool {
	varSrc := var_.GetSourceFile()
	refSrc := name.GetStaticSourceFile()

	return varSrc == nil || refSrc == nil || varSrc.GetName() == refSrc.GetName()
}

// Java: CheckAccessControls.isPackageAccessAllowed(Var, Node)
func (c *CheckAccessControls) IsPackageAccessAllowed(var_ *Var, name *rhino.Node) bool {
	varSrc := var_.GetSourceFile()
	refSrc := name.GetStaticSourceFile()
	if varSrc == nil && refSrc == nil {
		// If the source file of either var or name is unavailable, conservatively assume they belong
		// to different packages.
		//
		// TODO(brndn): by contrast, isPrivateAccessAllowed does allow private access when a source
		// file is unknown. I didn't change it in order not to break existing code.
		return false
	}

	codingConvention := c.compiler.GetCodingConvention()
	srcPackage := codingConvention.GetPackageName(varSrc)
	refPackage := codingConvention.GetPackageName(refSrc)
	return srcPackage != nil && refPackage != nil && *srcPackage == *refPackage
}

// Java: CheckAccessControls.checkPropertyOverrideVisibilityIsSame(Visibility, Visibility, Visibility, PropertyReference)
func (c *CheckAccessControls) CheckPropertyOverrideVisibilityIsSame(
	overriding rhino.JSDocInfo_Visibility,
	overridden rhino.JSDocInfo_Visibility,
	fileOverview *rhino.JSDocInfo_Visibility,
	propRef *CheckAccessControls_PropertyReference) {
	if overriding == rhino.JSDocInfo_Visibility_INHERITED &&
		overriding != overridden &&
		fileOverview != nil &&
		*fileOverview != rhino.JSDocInfo_Visibility_INHERITED {
		c.compiler.Report(
			JSError_MakeWithNode(
				propRef.SourceNode(),
				CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY,
				propRef.Name(),
				fileOverview.String()))
	}
}

// Java: CheckAccessControls.getOverridingPropertyVisibility(PropertyReference)
func CheckAccessControls_GetOverridingPropertyVisibility(propRef *CheckAccessControls_PropertyReference) *rhino.JSDocInfo_Visibility {
	overridingInfo := propRef.GetJSDocInfo()
	if overridingInfo == nil || !overridingInfo.IsOverride() {
		return nil
	}
	visibility := overridingInfo.GetVisibility()
	return &visibility
}

// CheckFinalClassOverrides checks if a constructor is trying to override a
// final class.
//
// Java: CheckAccessControls.checkFinalClassOverrides(Node)
func (c *CheckAccessControls) CheckFinalClassOverrides(ctor *rhino.Node) {
	if !CheckAccessControls_IsFunctionOrClass(ctor) ||
		// checking class constructors is redundant because we already check the same thing on
		// the CLASS node
		NodeUtil_IsEs6ConstructorMemberFunctionDef(ctor.GetParent()) {
		return
	}

	ctorType := ctor.GetJSType().ToMaybeFunctionType()
	if ctorType == nil || !ctorType.IsConstructor() {
		return
	}
	finalParentClass := CheckAccessControls_GetSuperClassInstanceIfFinal(ctorType)
	if finalParentClass != nil {
		c.compiler.Report(
			JSError_MakeWithNode(
				ctor,
				CheckAccessControls_EXTEND_FINAL_CLASS,
				jutil.ToJavaString(ctorType.GetDisplayName()),
				jutil.ToJavaString(finalParentClass.GetDisplayName())))
	}
}

// CheckConstantProperty determines whether the given constant property got
// reassigned.
//
// Java: CheckAccessControls.checkConstantProperty(PropertyReference, IdentifierBehaviour)
func (c *CheckAccessControls) CheckConstantProperty(
	propRef *CheckAccessControls_PropertyReference, identifierBehaviour CheckAccessControls_IdentifierBehaviour) {
	if propRef == nil || identifierBehaviour == CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE {
		return
	}

	sourceNode := propRef.SourceNode()
	if !propRef.Mutation() && !(propRef.Declaration() && sourceNode.IsFromExterns()) {
		return
	}

	objectType := CheckAccessControls_Dereference(propRef.ReceiverType())
	propertyName := propRef.Name()

	constness := c.IsPropertyDeclaredConstant(objectType, propertyName)
	if constness == CheckAccessControls_Constancy_MUTABLE {
		return
	}

	if sourceNode.IsFromExterns() && propRef.Declaration() {
		// Treat stub declarations in externs as inits, but never warn on them.
		c.RecordConstPropertyInit(propRef, objectType, constness)
		return
	}

	if !propRef.Mutation() {
		return
	}

	if propRef.IsDeletion() {
		c.compiler.Report(JSError_MakeWithNode(sourceNode, CheckAccessControls_CONST_PROPERTY_DELETED, propertyName))
		return
	}

	// Can't check for constant properties on generic function types.
	// TODO(johnlenz): I'm not 100% certain this is necessary, or if
	// the type is being inspected incorrectly.
	if objectType.IsFunctionType() && !objectType.ToMaybeFunctionType().IsConstructor() {
		return
	}

	if objectType.IsStructuralType() && !propRef.Declaration() {
		// We don't know the claess this structural type matches, so assume all assignments are bad.
		c.compiler.Report(
			JSError_MakeWithNode(
				sourceNode,
				CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE,
				propertyName,
				"unknown location due to structural typing"))
		return
	}

	init := c.GetConstPropertyInit(propRef, objectType)
	if init != nil {
		diagnostic := CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE
		if init.annotation == CheckAccessControls_Constancy_FINAL {
			diagnostic = CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN
		}
		c.compiler.Report(JSError_MakeWithNode(sourceNode, diagnostic, propertyName, init.node.GetLocation()))
	}

	c.RecordConstPropertyInit(propRef, objectType, constness)
}

// checkAccessControls_tableGet is Table.get(rowKey, columnKey).
//
// Go-only: language-adapter — see the constPropertyInits field.
func (c *CheckAccessControls) checkAccessControls_tableGet(rowKey rhino.JSType, columnKey string) *CheckAccessControls_ConstantDeclaration {
	row, ok := c.constPropertyInits.Lookup(rowKey)
	if !ok {
		return nil
	}
	return row.Get(columnKey)
}

// checkAccessControls_tableRowPutIfAbsent is
// Table.row(rowKey).putIfAbsent(columnKey, value).
//
// Go-only: language-adapter — see the constPropertyInits field.
func (c *CheckAccessControls) checkAccessControls_tableRowPutIfAbsent(
	rowKey rhino.JSType, columnKey string, value *CheckAccessControls_ConstantDeclaration) {
	row, ok := c.constPropertyInits.Lookup(rowKey)
	if ok {
		if row.Get(columnKey) != nil {
			return
		}
	} else {
		row = jutil.NewLinkedHashMap[string, *CheckAccessControls_ConstantDeclaration]()
		c.constPropertyInits.Put(rowKey, row)
	}
	row.Put(columnKey, value)
}

// Java: CheckAccessControls.getConstPropertyInit(PropertyReference, ObjectType)
func (c *CheckAccessControls) GetConstPropertyInit(
	ref *CheckAccessControls_PropertyReference, type_ rhino.ObjectType) *CheckAccessControls_ConstantDeclaration {
	name := ref.Name()
	for type_ != nil {
		init := c.checkAccessControls_tableGet(type_, name)
		if init != nil {
			return init
		}
		canonicalInit := c.checkAccessControls_tableGet(CheckAccessControls_GetCanonicalInstance(type_), name)
		if canonicalInit != nil {
			return canonicalInit
		}

		type_ = type_.GetImplicitPrototype()
	}

	return nil
}

// Java: CheckAccessControls.recordConstPropertyInit(PropertyReference, ObjectType, Constancy)
func (c *CheckAccessControls) RecordConstPropertyInit(
	ref *CheckAccessControls_PropertyReference, type_ rhino.ObjectType, annotation CheckAccessControls_Constancy) {
	c.checkAccessControls_tableRowPutIfAbsent(
		type_, ref.Name(), NewCheckAccessControls_ConstantDeclaration(ref.SourceNode(), annotation))

	// Add the prototype when we're looking at an instance object
	if type_.IsInstanceType() {
		prototype := type_.GetImplicitPrototype()
		if prototype != nil && prototype.HasProperty(ref.Name()) {
			c.checkAccessControls_tableRowPutIfAbsent(
				prototype, ref.Name(), NewCheckAccessControls_ConstantDeclaration(ref.SourceNode(), annotation))
		}
	}
}

// CheckAccessControls_GetCanonicalInstance returns an object with the same
// nominal type as obj, but without any possible extra properties that exist
// on obj.
//
// Java: CheckAccessControls.getCanonicalInstance(ObjectType)
func CheckAccessControls_GetCanonicalInstance(obj rhino.ObjectType) rhino.ObjectType {
	ctor := obj.GetConstructor()
	if ctor == nil {
		return obj
	}
	return ctor.GetInstanceType()
}

// CheckAccessControls_Dereference dereferences a type, autoboxing it and
// filtering out null.
//
// Java: CheckAccessControls.dereference(JSType)
func CheckAccessControls_Dereference(type_ rhino.JSType) rhino.ObjectType {
	if type_ == nil {
		return nil
	}
	return type_.Dereference()
}

// Java: CheckAccessControls.typeOrUnknown(JSType)
func (c *CheckAccessControls) TypeOrUnknown(type_ rhino.JSType) rhino.JSType {
	if type_ == nil {
		return c.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	}
	return type_
}

// Java: CheckAccessControls.typeOrUnknown(ObjectType)
func (c *CheckAccessControls) TypeOrUnknownObjectType(type_ rhino.ObjectType) rhino.ObjectType {
	return c.TypeOrUnknown(type_).(rhino.ObjectType)
}

// Java: CheckAccessControls.boxedOrUnknown(JSType)
func (c *CheckAccessControls) BoxedOrUnknown(type_ rhino.JSType) rhino.ObjectType {
	return c.TypeOrUnknownObjectType(CheckAccessControls_Dereference(type_))
}

// CheckPropertyVisibility reports an error if the given property is not
// visible in the current context.
//
// This method covers both:
//
//   - accesses to properties during execution
//   - overrides of properties during declaration
//
// TODO(nickreid): Things would probably be a lot simpler, though a bit
// duplicated, if these two concepts were separated. Much of the underlying
// logic could stop checking various inconsistent definitions of "is this an
// override".
//
// Java: CheckAccessControls.checkPropertyVisibility(PropertyReference)
func (c *CheckAccessControls) CheckPropertyVisibility(propRef *CheckAccessControls_PropertyReference) {
	if NodeUtil_IsEs6ConstructorMemberFunctionDef(propRef.SourceNode()) {
		// Class ctor *declarations* can never violate visibility restrictions. They are not
		// accesses and we don't consider them overrides.
		//
		// TODO(nickreid): It would be a lot cleaner if we could model this using `PropertyReference`
		// rather than defining a special case here. I think the problem is that the current
		// implementation of this method conflates "override" with "declaration". But that only works
		// because it ignores cases where there's no overridden definition.
		return
	}

	rawReferenceType := c.TypeOrUnknownObjectType(propRef.ReceiverType()).Autobox()
	referenceType := CheckAccessControls_CastToObject(rawReferenceType)

	propertyName := propRef.Name()

	// Is this a normal property access, or are we trying to override
	// an existing property?
	isOverride := propRef.IsDocumentedDeclaration() || propRef.Override()

	objectType := AccessControlUtils_GetObjectType(referenceType, isOverride, propertyName)

	var definingSource rhino.StaticSourceFile
	var fileOverviewVisibility *rhino.JSDocInfo_Visibility
	if !c.defaultVisibilityForFiles.IsEmpty() {
		definingSource = AccessControlUtils_GetDefiningSource(propRef.SourceNode(), referenceType, propertyName)
		fileOverviewVisibility = checkAccessControls_visibilityFor(c.defaultVisibilityForFiles, definingSource)
	}

	var defSite *rhino.Node
	if objectType != nil {
		defSite = objectType.GetOwnPropertyDefSite(propertyName)
		if defSite == nil && fileOverviewVisibility == nil {
			return
		}
	} else if fileOverviewVisibility == nil {
		return
	}

	var visibility rhino.JSDocInfo_Visibility
	if propRef.Override() {
		effObjType := objectType
		if !isOverride {
			effObjType = AccessControlUtils_GetObjectType(referenceType, true, propertyName)
		}
		overridden := AccessControlUtils_GetOverriddenPropertyVisibility(effObjType, propertyName)
		visibility = AccessControlUtils_GetEffectiveVisibilityForOverriddenProperty(
			overridden, fileOverviewVisibility, propertyName)
	} else {
		effObjType := objectType
		if isOverride {
			effObjType = AccessControlUtils_GetObjectType(referenceType, false, propertyName)
		}
		visibility = CheckAccessControls_GetEffectiveVisibilityForNonOverriddenProperty(
			propRef, effObjType, fileOverviewVisibility)
	}

	if isOverride {
		overriding := CheckAccessControls_GetOverridingPropertyVisibility(propRef)
		if overriding != nil {
			c.CheckPropertyOverrideVisibilityIsSame(
				*overriding, visibility, fileOverviewVisibility, propRef)
		}
	}

	reportType := rawReferenceType
	if objectType != nil {
		if defSite == nil {
			// Assume the property is public.
			return
		}
		reportType = objectType
		definingSource = defSite.GetStaticSourceFile()
	} else if fileOverviewVisibility == nil {
		// We can only check visibility references if we know what file
		// it was defined in.
		// Otherwise just assume the property is public.
		return
	}

	referenceSource := propRef.SourceNode().GetStaticSourceFile()

	if isOverride {
		sameInput := referenceSource != nil && referenceSource.GetName() == definingSource.GetName()
		c.CheckPropertyOverrideVisibility(
			propRef, visibility, fileOverviewVisibility, reportType, sameInput)
	} else {
		c.CheckPropertyAccessVisibility(
			propRef, visibility, reportType, referenceSource, definingSource)
	}
}

// checkAccessControls_visibilityFor is
// ImmutableMap<StaticSourceFile, Visibility>.get(file): null when absent.
//
// Go-only: language-adapter — a nullable Visibility is a pointer.
func checkAccessControls_visibilityFor(
	m *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility],
	file rhino.StaticSourceFile) *rhino.JSDocInfo_Visibility {
	if file == nil || m.IsEmpty() {
		return nil
	}
	if v, ok := m.Lookup(file); ok {
		return &v
	}
	return nil
}

// CheckEs6ConstructorInvocationVisibility reports visibility violations on
// ES6 class constructor invocations.
//
// Precondition: target has an ES6 class JSType.
//
// Java: CheckAccessControls.checkEs6ConstructorInvocationVisibility(Node)
func (c *CheckAccessControls) CheckEs6ConstructorInvocationVisibility(target *rhino.Node) {
	ctorType := target.GetJSType().ToMaybeFunctionType()
	prototypeType := ctorType.GetPrototype()

	// We use the class definition site because classes automatically get a implicit constructor,
	// so there may not be a definition node.
	classDefinition := ctorType.GetSource()

	var definingSource rhino.StaticSourceFile
	if classDefinition != nil {
		definingSource = AccessControlUtils_GetDefiningSource(classDefinition, prototypeType, "constructor")
	}

	// Synthesize a `PropertyReference` for this constructor call as if we're accessing
	// `Foo.prototype.constructor`. This object allows us to reuse the
	// `checkPropertyAccessVisibility` method which actually reports violations.
	fauxCtorRef := CheckAccessControls_PropertyReference_Builder_().
		SetSourceNode(target).
		SetName("constructor").
		SetReceiverType(prototypeType).
		SetMutation(false).    // This shouldn't matter.
		SetDeclaration(false). // This shouldn't matter.
		SetOverride(false).    // This shouldn't matter.
		SetReadableTypeName(func() string { return ctorType.GetInstanceType().String() }).
		Build()

	annotatedCtorVisibility :=
		// This function defaults to `INHERITED` which isn't what we want here, but it does handle
		// combining inline and `@fileoverview` visibilities.
		CheckAccessControls_GetEffectiveVisibilityForNonOverriddenProperty(
			fauxCtorRef, prototypeType, checkAccessControls_visibilityFor(c.defaultVisibilityForFiles, definingSource))
	effectiveCtorVisibility := annotatedCtorVisibility
	if annotatedCtorVisibility == rhino.JSDocInfo_Visibility_INHERITED {
		effectiveCtorVisibility = rhino.JSDocInfo_Visibility_PUBLIC
	}

	c.CheckPropertyAccessVisibility(
		fauxCtorRef,
		effectiveCtorVisibility,
		ctorType,
		target.GetStaticSourceFile(),
		definingSource)
}

// Java: CheckAccessControls.checkPropertyOverrideVisibility(PropertyReference, Visibility, Visibility, JSType, boolean)
func (c *CheckAccessControls) CheckPropertyOverrideVisibility(
	propRef *CheckAccessControls_PropertyReference,
	visibility rhino.JSDocInfo_Visibility,
	fileOverviewVisibility *rhino.JSDocInfo_Visibility,
	objectType rhino.JSType,
	sameInput bool) {
	overridingVisibility := rhino.JSDocInfo_Visibility_INHERITED
	if propRef.Override() {
		overridingVisibility = propRef.GetJSDocInfo().GetVisibility()
	}

	// Check that:
	// (a) the property *can* be overridden,
	// (b) the visibility of the override is the same as (or broader than) the
	//     visibility of the original property,
	// (c) the visibility is explicitly redeclared if the override is in
	//     a file with default visibility in the @fileoverview block.
	if visibility == rhino.JSDocInfo_Visibility_PRIVATE && !sameInput {
		c.compiler.Report(JSError_MakeWithNode(propRef.SourceNode(), CheckAccessControls_PRIVATE_OVERRIDE, objectType.String()))
	} else if !CheckAccessControls_CanOverrideVisibility(visibility, overridingVisibility) &&
		fileOverviewVisibility == nil {
		c.compiler.Report(
			JSError_MakeWithNode(
				propRef.SourceNode(),
				CheckAccessControls_VISIBILITY_MISMATCH,
				visibility.String(),
				objectType.String(),
				overridingVisibility.String()))
	}
}

// Java: CheckAccessControls.canOverrideVisibility(Visibility, Visibility)
func CheckAccessControls_CanOverrideVisibility(
	superclassVisibility rhino.JSDocInfo_Visibility, subclassVisibility rhino.JSDocInfo_Visibility) bool {
	// This allows INHERITED to override anything, PUBLIC to override anything (except INHERITED),
	// and PROTECTED to override anything (except PUBLIC or INHERITED). PRIVATE was already handled
	// in a previous check, leaving PACKAGE as the lowest visibility.
	return superclassVisibility.Ordinal()-subclassVisibility.Ordinal() <= 0
}

// Java: CheckAccessControls.checkPropertyAccessVisibility(PropertyReference, Visibility, JSType, StaticSourceFile, StaticSourceFile)
func (c *CheckAccessControls) CheckPropertyAccessVisibility(
	propRef *CheckAccessControls_PropertyReference,
	visibility rhino.JSDocInfo_Visibility,
	objectType rhino.JSType,
	referenceSource rhino.StaticSourceFile,
	definingSource rhino.StaticSourceFile) {
	// private access is always allowed in the same file.
	if referenceSource != nil &&
		definingSource != nil &&
		referenceSource.GetName() == definingSource.GetName() {
		return
	}

	ownerType := CheckAccessControls_InstanceTypeFor(objectType)

	switch visibility {
	case rhino.JSDocInfo_Visibility_PACKAGE:
		c.CheckPackagePropertyVisibility(propRef, referenceSource, definingSource)
	case rhino.JSDocInfo_Visibility_PRIVATE:
		c.CheckPrivatePropertyVisibility(propRef, ownerType)
	case rhino.JSDocInfo_Visibility_PROTECTED:
		c.CheckProtectedPropertyVisibility(propRef, ownerType)
	default:
	}
}

// Java: CheckAccessControls.checkPackagePropertyVisibility(PropertyReference, StaticSourceFile, StaticSourceFile)
func (c *CheckAccessControls) CheckPackagePropertyVisibility(
	propRef *CheckAccessControls_PropertyReference,
	referenceSource rhino.StaticSourceFile,
	definingSource rhino.StaticSourceFile) {
	codingConvention := c.compiler.GetCodingConvention()
	refPackage := codingConvention.GetPackageName(referenceSource)
	defPackage := codingConvention.GetPackageName(definingSource)
	if refPackage == nil || defPackage == nil || *refPackage != *defPackage {
		c.compiler.Report(
			JSError_MakeWithNode(
				propRef.SourceNode(),
				CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS,
				propRef.Name(),
				propRef.GetReadableTypeNameOrDefault()))
	}
}

// Java: CheckAccessControls.checkPrivatePropertyVisibility(PropertyReference, ObjectType)
func (c *CheckAccessControls) CheckPrivatePropertyVisibility(
	propRef *CheckAccessControls_PropertyReference, ownerType rhino.ObjectType) {

	// private access is not allowed outside the file from a different
	// enclosing class.
	// TODO(tbreisacher): Should we also include the filename where ownerType is defined?
	var readableTypeName string
	if ownerType == nil || ownerType.Equals(propRef.ReceiverType()) {
		readableTypeName = propRef.GetReadableTypeNameOrDefault()
	} else {
		readableTypeName = ownerType.String()
	}
	c.compiler.Report(
		JSError_MakeWithNode(
			propRef.SourceNode(), CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS, propRef.Name(), readableTypeName))
}

// Java: CheckAccessControls.checkProtectedPropertyVisibility(PropertyReference, ObjectType)
func (c *CheckAccessControls) CheckProtectedPropertyVisibility(
	propRef *CheckAccessControls_PropertyReference, ownerType rhino.ObjectType) {
	// There are 3 types of legal accesses of a protected property:
	// 1) Accesses in the same file
	// 2) Overriding the property in a subclass
	// 3) Accessing the property from inside a subclass
	// The first two have already been checked for.
	if ownerType != nil {
		for i := len(c.currentClassStack) - 1; i >= 0; i-- {
			scopeType := c.currentClassStack[i]
			if scopeType == nil {
				continue
			} else if scopeType.IsSubtypeOf(ownerType) {
				return
			}
		}
	}

	c.compiler.Report(
		JSError_MakeWithNode(
			propRef.SourceNode(),
			CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS,
			propRef.Name(),
			propRef.GetReadableTypeNameOrDefault()))
}

// ShouldEmitDeprecationWarning determines whether a deprecation warning
// should be emitted.
//
// t is the current traversal; n the node which we are checking.
//
// Java: CheckAccessControls.shouldEmitDeprecationWarning(NodeTraversal, Node)
func (c *CheckAccessControls) ShouldEmitDeprecationWarning(t *NodeTraversal, n *rhino.Node) bool {
	// In the global scope, there are only two kinds of accesses that should
	// be flagged for warnings:
	// 1) Calls of deprecated functions and methods.
	// 2) Instantiations of deprecated classes.
	// For now, we just let everything else by.
	if t.InGlobalScope() {
		if !NodeUtil_IsInvocationTarget(n) && !n.IsNew() {
			return false
		}
	}

	return !c.CanAccessDeprecatedTypes(t)
}

// ShouldEmitDeprecationWarningForPropRef determines whether a deprecation
// warning should be emitted.
//
// Java: CheckAccessControls.shouldEmitDeprecationWarning(NodeTraversal, PropertyReference)
func (c *CheckAccessControls) ShouldEmitDeprecationWarningForPropRef(t *NodeTraversal, propRef *CheckAccessControls_PropertyReference) bool {
	// In the global scope, there are only two kinds of accesses that should
	// be flagged for warnings:
	// 1) Calls of deprecated functions and methods.
	// 2) Instantiations of deprecated classes.
	// For now, we just let everything else by.
	if t.InGlobalScope() {
		if !NodeUtil_IsInvocationTarget(propRef.SourceNode()) {
			return false
		}
	}

	// We can always assign to a deprecated property, to keep it up to date.
	if propRef.Mutation() {
		return false
	}

	// Don't warn if the node is just declaring the property, not reading it.
	if propRef.Declaration() {
		jsdoc := propRef.GetJSDocInfo()
		if jsdoc != nil && jsdoc.IsDeprecated() {
			return false
		}
	}

	return !c.CanAccessDeprecatedTypes(t)
}

// CanAccessDeprecatedTypes returns whether it's currently OK to access
// deprecated names and properties.
//
// There are 3 exceptions when we're allowed to use a deprecated type or
// property:
//  1. When we're in a deprecated function.
//  2. When we're in a deprecated class.
//  3. When we're in a static method of a deprecated class.
//
// Java: CheckAccessControls.canAccessDeprecatedTypes(NodeTraversal)
func (c *CheckAccessControls) CanAccessDeprecatedTypes(t *NodeTraversal) bool {
	// Case #1
	if c.deprecationDepth > 0 {
		return true
	}
	hoistRoot := t.GetClosestHoistScopeRoot()
	if hoistRoot == c.lastHoistScopeRoot && hoistRoot != nil {
		return c.lastHoistScopeDeprecated
	}
	scopeRoot := hoistRoot
	if NodeUtil_IsFunctionBlock(scopeRoot) {
		scopeRoot = scopeRoot.GetParent()
	}
	scopeRootParent := scopeRoot.GetParent()

	// Cases 2 and 3 are required to handle ES5-style class methods since they aren't nested inside
	// their class. This is tested in the CheckAccessControlsOldSyntaxTest class.
	res :=
		// Case #2
		(CheckAccessControls_GetTypeDeprecationInfo(c.GetTypeOfThis(scopeRoot)) != nil) ||
			// Case #3
			(scopeRootParent != nil &&
				scopeRootParent.IsAssign() &&
				CheckAccessControls_GetTypeDeprecationInfo(c.BestInstanceTypeForMethodOrCtor(scopeRoot)) != nil)
	c.lastHoistScopeRoot = hoistRoot
	c.lastHoistScopeDeprecated = res
	return res
}

// CheckAccessControls_IsMarkedDeprecated returns whether this node roots a
// subtree under which references to deprecated constructs are allowed.
//
// Java: CheckAccessControls.isMarkedDeprecated(Node)
func CheckAccessControls_IsMarkedDeprecated(n *rhino.Node) bool {
	return CheckAccessControls_GetDeprecationReason(NodeUtil_GetBestJSDocInfo(n)) != nil
}

// CheckAccessControls_GetTypeDeprecationInfo returns the deprecation reason
// for the type if it is marked as being deprecated. Returns empty string if
// the type is deprecated but no reason was given. Returns null if the type is
// not deprecated.
//
// Java: CheckAccessControls.getTypeDeprecationInfo(JSType)
func CheckAccessControls_GetTypeDeprecationInfo(type_ rhino.JSType) *string {
	if type_ == nil {
		return nil
	}

	return CheckAccessControls_GetDeprecationReason(type_.GetJSDocInfo())
}

var checkAccessControls_emptyDeprecationReason = ""

// Java: CheckAccessControls.getDeprecationReason(JSDocInfo)
func CheckAccessControls_GetDeprecationReason(info *rhino.JSDocInfo) *string {
	if info != nil && info.IsDeprecated() {
		if info.GetDeprecationReason() != nil {
			return info.GetDeprecationReason()
		}
		return &checkAccessControls_emptyDeprecationReason
	}
	return nil
}

// IsPropertyDeclaredConstant returns if a property is declared constant.
//
// Java: CheckAccessControls.isPropertyDeclaredConstant(ObjectType, String)
func (c *CheckAccessControls) IsPropertyDeclaredConstant(objectType rhino.ObjectType, prop string) CheckAccessControls_Constancy {
	for ; objectType != nil; objectType = objectType.GetImplicitPrototype() {
		docInfo := objectType.GetOwnPropertyJSDocInfo(prop)
		if docInfo == nil {
			continue
		} else if docInfo.IsFinal() {
			return CheckAccessControls_Constancy_FINAL
		} else if docInfo.IsConstant() {
			return CheckAccessControls_Constancy_OTHER_CONSTANT
		}
	}
	return CheckAccessControls_Constancy_MUTABLE
}

// CheckAccessControls_GetPropertyDeprecationInfo returns the deprecation
// reason for the property if it is marked as being deprecated. Returns empty
// string if the property is deprecated but no reason was given. Returns null
// if the property is not deprecated.
//
// Java: CheckAccessControls.getPropertyDeprecationInfo(ObjectType, String)
func CheckAccessControls_GetPropertyDeprecationInfo(type_ rhino.ObjectType, prop string) *string {
	for ; type_ != nil; type_ = type_.GetImplicitPrototype() {
		if info := type_.GetOwnPropertyJSDocInfo(prop); info != nil {
			if depReason := CheckAccessControls_GetDeprecationReason(info); depReason != nil {
				return depReason
			}
		}
	}
	return nil
}

// CheckAccessControls_GetSuperClassInstanceIfFinal: if the superclass is
// final, this method returns an instance of the superclass.
//
// Java: CheckAccessControls.getSuperClassInstanceIfFinal(FunctionType)
func CheckAccessControls_GetSuperClassInstanceIfFinal(subCtor rhino.FunctionType) rhino.ObjectType {
	ctor := subCtor.GetSuperClassConstructor()
	var doc *rhino.JSDocInfo
	if ctor != nil {
		doc = ctor.GetJSDocInfo()
	}
	if doc != nil && doc.IsFinal() {
		return ctor.GetInstanceType()
	}

	return nil
}

// Java: CheckAccessControls.castToObject(JSType)
func CheckAccessControls_CastToObject(type_ rhino.JSType) rhino.ObjectType {
	if type_ == nil {
		return nil
	}
	return type_.ToMaybeObjectType()
}

// Java: CheckAccessControls.isFunctionOrClass(Node)
func CheckAccessControls_IsFunctionOrClass(n *rhino.Node) bool {
	return n.IsFunction() || n.IsClass()
}

// Java: CheckAccessControls.isExtendsTarget(Node)
func CheckAccessControls_IsExtendsTarget(node *rhino.Node) bool {
	parent := node.GetParent()
	return parent.IsClass() && node.IsSecondChildOf(parent)
}

// Java: CheckAccessControls.getTypeOfThis(Node)
func (c *CheckAccessControls) GetTypeOfThis(scopeRoot *rhino.Node) rhino.JSType {
	if scopeRoot.IsRoot() || scopeRoot.IsScript() {
		return CheckAccessControls_CastToObject(scopeRoot.GetJSType())
	} else if scopeRoot.IsModuleBody() {
		return nil
	} else if NodeUtil_IsClassStaticBlock(scopeRoot) {
		classNode := NodeUtil_GetEnclosingClass(scopeRoot)
		return classNode.GetJSType()
	}

	jutil.CheckArgumentMsg1(scopeRoot.IsFunction(), "%s", scopeRoot)

	nodeType := scopeRoot.GetJSType()
	if nodeType != nil && nodeType.IsFunctionType() {
		return nodeType.ToMaybeFunctionType().GetTypeOfThis()
	} else {
		// Executed when the current scope has not been typechecked.
		return nil
	}
}

// CheckAccessControls_IdentifierBehaviour is the set of ways in which JSDoc
// and identifier usage can interact.
//
// Before ES6, for a given type, there was no way to differentate the
// declaration of the constructor function, namespace object, and type; the
// declaration of all three constructs was the same node. This lead to some
// assumptions were made about how access-control modifiers applied to each.
//
// At ES6, class-syntax separated the constructor function from the namespace
// object and type declaration. This allowed finer grained control over
// access-control modifiers; however it broke some of the eariler assumptions.
//
// This type exists to simplify maintaining both sets of assumptions. It
// allows other code to branch on behaviour in a more obvious way.
//
// TODO(b/113127707): Make this unnecessary by better modeling or decomposing
// this pass.
//
// Java: CheckAccessControls.IdentifierBehaviour
type CheckAccessControls_IdentifierBehaviour int32

const (
	CheckAccessControls_IdentifierBehaviour_NON_CONSTRUCTOR CheckAccessControls_IdentifierBehaviour = iota
	CheckAccessControls_IdentifierBehaviour_ES5_CLASS_INVOCATION
	CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE
	CheckAccessControls_IdentifierBehaviour_ES6_CLASS_INVOCATION
	CheckAccessControls_IdentifierBehaviour_ES6_CLASS_NAMESPACE
)

var checkAccessControls_IdentifierBehaviour_names = [...]string{
	"NON_CONSTRUCTOR", "ES5_CLASS_INVOCATION", "ES5_CLASS_NAMESPACE", "ES6_CLASS_INVOCATION", "ES6_CLASS_NAMESPACE",
}

// CheckAccessControls_IdentifierBehaviour_values is IdentifierBehaviour.values().
func CheckAccessControls_IdentifierBehaviour_values() []CheckAccessControls_IdentifierBehaviour {
	return []CheckAccessControls_IdentifierBehaviour{
		CheckAccessControls_IdentifierBehaviour_NON_CONSTRUCTOR,
		CheckAccessControls_IdentifierBehaviour_ES5_CLASS_INVOCATION,
		CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE,
		CheckAccessControls_IdentifierBehaviour_ES6_CLASS_INVOCATION,
		CheckAccessControls_IdentifierBehaviour_ES6_CLASS_NAMESPACE,
	}
}

// String is Enum.name().
func (b CheckAccessControls_IdentifierBehaviour) String() string {
	return checkAccessControls_IdentifierBehaviour_names[b]
}

// Ordinal is Enum.ordinal().
func (b CheckAccessControls_IdentifierBehaviour) Ordinal() int { return int(b) }

// Java: CheckAccessControls.IdentifierBehaviour.select(Node)
func CheckAccessControls_IdentifierBehaviour_Select(target *rhino.Node) CheckAccessControls_IdentifierBehaviour {
	type_ := target.GetJSType()
	if type_ == nil || !type_.IsFunctionType() {
		// If we aren't sure what we're dealing with be more strict.
		return CheckAccessControls_IdentifierBehaviour_NON_CONSTRUCTOR
	}

	ctorType := type_.ToMaybeFunctionType()
	if !ctorType.IsConstructor() {
		return CheckAccessControls_IdentifierBehaviour_NON_CONSTRUCTOR
	}

	isInvocation := NodeUtil_IsInvocationTarget(target) || CheckAccessControls_IsExtendsTarget(target)
	isEs6 := (ctorType.GetSource() != nil) && ctorType.GetSource().IsClass()

	if !isEs6 {
		if isInvocation {
			return CheckAccessControls_IdentifierBehaviour_ES5_CLASS_INVOCATION
		}
		return CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE
	} else {
		if isInvocation {
			return CheckAccessControls_IdentifierBehaviour_ES6_CLASS_INVOCATION
		}
		return CheckAccessControls_IdentifierBehaviour_ES6_CLASS_NAMESPACE
	}
}

// CheckAccessControls_PropertyReference is a representation of an object
// property reference in JS code.
//
// This is an abstraction to smooth over the various AST structures that can
// act on properties. It is not useful for names (variables) or anonymous JS
// constructs.
//
// This class should only be used within CheckAccessControls. Having
// package-private visibility is a quirk of AutoValue.
//
// sourceNode: the Node that spawned this reference. receiverType: the type
// from which the property is referenced, not necessarily the one that
// declared it. readableTypeName: a lazy source for a human-readable type name
// to use when generating messages. Most users probably want
// GetReadableTypeNameOrDefault().
//
// Java: CheckAccessControls.PropertyReference (a record)
type CheckAccessControls_PropertyReference struct {
	sourceNode       *rhino.Node
	name             string
	receiverType     rhino.ObjectType
	mutation         bool
	declaration      bool
	override         bool
	readableTypeName func() string
	typeRegistry     *rhino.JSTypeRegistry
	readableTypeNode *rhino.Node
	jsdoc            *rhino.JSDocInfo
	hasJSDoc         bool
}

// checkAccessControls_requireNonNull is Objects.requireNonNull(obj, message).
func checkAccessControls_requireNonNull(isNull bool, message string) {
	if isNull {
		panic(jutil.NewNullPointerExceptionMsg(message))
	}
}

// NewCheckAccessControls_PropertyReference is the record's canonical
// constructor.
//
// Java: CheckAccessControls.PropertyReference(Node, String, ObjectType, boolean, boolean, boolean, Supplier<String>)
func NewCheckAccessControls_PropertyReference(
	sourceNode *rhino.Node,
	name string,
	receiverType rhino.ObjectType,
	mutation bool,
	declaration bool,
	override bool,
	readableTypeName func() string) *CheckAccessControls_PropertyReference {
	checkAccessControls_requireNonNull(sourceNode == nil, "sourceNode")
	// Go-only: language-adapter — a Go string cannot be null:
	// requireNonNull(name, "name") cannot fail.
	checkAccessControls_requireNonNull(receiverType == nil, "receiverType")
	checkAccessControls_requireNonNull(readableTypeName == nil, "readableTypeName")
	r := &CheckAccessControls_PropertyReference{}
	r.sourceNode = sourceNode
	r.name = name
	r.receiverType = receiverType
	r.mutation = mutation
	r.declaration = declaration
	r.override = override
	r.readableTypeName = readableTypeName
	return r
}

// Java: CheckAccessControls.PropertyReference.sourceNode()
func (r *CheckAccessControls_PropertyReference) SourceNode() *rhino.Node { return r.sourceNode }

// Java: CheckAccessControls.PropertyReference.name()
func (r *CheckAccessControls_PropertyReference) Name() string { return r.name }

// Java: CheckAccessControls.PropertyReference.receiverType()
func (r *CheckAccessControls_PropertyReference) ReceiverType() rhino.ObjectType {
	return r.receiverType
}

// Java: CheckAccessControls.PropertyReference.mutation()
func (r *CheckAccessControls_PropertyReference) Mutation() bool { return r.mutation }

// Java: CheckAccessControls.PropertyReference.declaration()
func (r *CheckAccessControls_PropertyReference) Declaration() bool { return r.declaration }

// Java: CheckAccessControls.PropertyReference.override()
func (r *CheckAccessControls_PropertyReference) Override() bool { return r.override }

// Java: CheckAccessControls.PropertyReference.readableTypeName()
func (r *CheckAccessControls_PropertyReference) ReadableTypeName() func() string {
	if r.readableTypeName != nil {
		return r.readableTypeName
	}
	return func() string {
		if r.readableTypeNode != nil && r.typeRegistry != nil {
			return r.typeRegistry.GetReadableTypeName(r.readableTypeNode)
		}
		return ""
	}
}

// Java: CheckAccessControls.PropertyReference.builder()
func CheckAccessControls_PropertyReference_Builder_() *CheckAccessControls_PropertyReference_Builder {
	return NewAutoBuilder_CheckAccessControls_PropertyReference_Builder()
}

// CheckAccessControls_PropertyReference_Builder is the @AutoBuilder interface
// CheckAccessControls.PropertyReference.Builder.
//
// Go-only: language-adapter — the AutoBuilder-generated
// AutoBuilder_CheckAccessControls_PropertyReference_Builder (its only
// implementation) is folded into this struct.
type CheckAccessControls_PropertyReference_Builder struct {
	sourceNode       *rhino.Node
	name             *string
	receiverType     rhino.ObjectType
	mutation         bool
	declaration      bool
	override         bool
	readableTypeName func() string
	set_0            int8
}

// Java: AutoBuilder_CheckAccessControls_PropertyReference_Builder()
func NewAutoBuilder_CheckAccessControls_PropertyReference_Builder() *CheckAccessControls_PropertyReference_Builder {
	return &CheckAccessControls_PropertyReference_Builder{}
}

// Java: AutoBuilder_CheckAccessControls_PropertyReference_Builder(CheckAccessControls.PropertyReference)
func NewAutoBuilder_CheckAccessControls_PropertyReference_BuilderFromSource(
	source *CheckAccessControls_PropertyReference) *CheckAccessControls_PropertyReference_Builder {
	b := &CheckAccessControls_PropertyReference_Builder{}
	b.sourceNode = source.SourceNode()
	name := source.Name()
	b.name = &name
	b.receiverType = source.ReceiverType()
	b.mutation = source.Mutation()
	b.declaration = source.Declaration()
	b.override = source.Override()
	b.readableTypeName = source.ReadableTypeName()
	b.set_0 = int8(7)
	return b
}

// Java: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setSourceNode(Node)
func (b *CheckAccessControls_PropertyReference_Builder) SetSourceNode(sourceNode *rhino.Node) *CheckAccessControls_PropertyReference_Builder {
	if sourceNode == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null sourceNode"))
	}
	b.sourceNode = sourceNode
	return b
}

// Java: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setName(String)
func (b *CheckAccessControls_PropertyReference_Builder) SetName(name string) *CheckAccessControls_PropertyReference_Builder {
	// Go-only: language-adapter — a Go string cannot be null, so the
	// "Null name" check cannot fail.
	b.name = &name
	return b
}

// Java: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setReceiverType(ObjectType)
func (b *CheckAccessControls_PropertyReference_Builder) SetReceiverType(receiverType rhino.ObjectType) *CheckAccessControls_PropertyReference_Builder {
	if receiverType == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null receiverType"))
	}
	b.receiverType = receiverType
	return b
}

// Java: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setMutation(boolean)
func (b *CheckAccessControls_PropertyReference_Builder) SetMutation(mutation bool) *CheckAccessControls_PropertyReference_Builder {
	b.mutation = mutation
	b.set_0 |= int8(1)
	return b
}

// Java: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setDeclaration(boolean)
func (b *CheckAccessControls_PropertyReference_Builder) SetDeclaration(declaration bool) *CheckAccessControls_PropertyReference_Builder {
	b.declaration = declaration
	b.set_0 |= int8(2)
	return b
}

// Java: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setOverride(boolean)
func (b *CheckAccessControls_PropertyReference_Builder) SetOverride(override bool) *CheckAccessControls_PropertyReference_Builder {
	b.override = override
	b.set_0 |= int8(4)
	return b
}

// Java: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setReadableTypeName(Supplier<String>)
func (b *CheckAccessControls_PropertyReference_Builder) SetReadableTypeName(readableTypeName func() string) *CheckAccessControls_PropertyReference_Builder {
	if readableTypeName == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null readableTypeName"))
	}
	b.readableTypeName = readableTypeName
	return b
}

// Java: AutoBuilder_CheckAccessControls_PropertyReference_Builder.build()
func (b *CheckAccessControls_PropertyReference_Builder) Build() *CheckAccessControls_PropertyReference {
	if b.set_0 != 7 ||
		b.sourceNode == nil ||
		b.name == nil ||
		b.receiverType == nil ||
		b.readableTypeName == nil {
		var missing strings.Builder
		if b.sourceNode == nil {
			missing.WriteString(" sourceNode")
		}
		if b.name == nil {
			missing.WriteString(" name")
		}
		if b.receiverType == nil {
			missing.WriteString(" receiverType")
		}
		if (b.set_0 & 1) == 0 {
			missing.WriteString(" mutation")
		}
		if (b.set_0 & 2) == 0 {
			missing.WriteString(" declaration")
		}
		if (b.set_0 & 4) == 0 {
			missing.WriteString(" override")
		}
		if b.readableTypeName == nil {
			missing.WriteString(" readableTypeName")
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing.String()))
	}
	return NewCheckAccessControls_PropertyReference(
		b.sourceNode,
		*b.name,
		b.receiverType,
		b.mutation,
		b.declaration,
		b.override,
		b.readableTypeName)
}

// Derived properties.

// Java: CheckAccessControls.PropertyReference.getParentNode()
func (r *CheckAccessControls_PropertyReference) GetParentNode() *rhino.Node {
	return r.SourceNode().GetParent()
}

// Java: CheckAccessControls.PropertyReference.getJSType()
func (r *CheckAccessControls_PropertyReference) GetJSType() rhino.JSType {
	return r.SourceNode().GetJSType()
}

// Java: CheckAccessControls.PropertyReference.getJSDocInfo()
func (r *CheckAccessControls_PropertyReference) GetJSDocInfo() *rhino.JSDocInfo {
	if r.hasJSDoc {
		return r.jsdoc
	}
	r.jsdoc = NodeUtil_GetBestJSDocInfo(r.SourceNode())
	r.hasJSDoc = true
	return r.jsdoc
}

// Java: CheckAccessControls.PropertyReference.isDocumentedDeclaration()
func (r *CheckAccessControls_PropertyReference) IsDocumentedDeclaration() bool {
	return r.Declaration() && (r.GetJSDocInfo() != nil)
}

// Java: CheckAccessControls.PropertyReference.isDeletion()
func (r *CheckAccessControls_PropertyReference) IsDeletion() bool {
	return r.SourceNode().GetParent().IsDelProp()
}

// Java: CheckAccessControls.PropertyReference.getReadableTypeNameOrDefault()
func (r *CheckAccessControls_PropertyReference) GetReadableTypeNameOrDefault() string {
	var preferred string
	if r.readableTypeName != nil {
		preferred = r.readableTypeName()
	} else if r.readableTypeNode != nil && r.typeRegistry != nil {
		preferred = r.typeRegistry.GetReadableTypeName(r.readableTypeNode)
	}
	if preferred == "" {
		return r.ReceiverType().String()
	}
	return preferred
}

// Java: CheckAccessControls.createPropertyReference(Node)
func (c *CheckAccessControls) CreatePropertyReference(sourceNode *rhino.Node) *CheckAccessControls_PropertyReference {
	switch sourceNode.GetToken() {
	case rhino.Token_GETPROP:
		parent := sourceNode.GetParent()
		isLValue := NodeUtil_IsLValue(sourceNode)
		firstChild := sourceNode.GetFirstChild()
		ref := &c.scratchPropRef
		ref.sourceNode = sourceNode
		ref.name = sourceNode.GetString()
		ref.receiverType = c.BoxedOrUnknown(firstChild.GetJSType())
		// Props are always mutated as L-values, even when assigned `undefined`.
		ref.mutation = isLValue || parent.IsDelProp()
		if isLValue {
			jsdoc := NodeUtil_GetBestJSDocInfo(sourceNode)
			ref.declaration = parent.IsExprResult() || (jsdoc != nil && jsdoc.IsConstant())
			// TODO(b/113704668): This definition is way too loose. It was used to prevent
			// breakages during refactoring and should be tightened.
			ref.override = jsdoc != nil
			ref.jsdoc = jsdoc
			ref.hasJSDoc = true
		} else {
			ref.declaration = parent.IsExprResult()
			ref.override = false
			ref.jsdoc = nil
			ref.hasJSDoc = false
		}
		ref.readableTypeName = nil
		ref.typeRegistry = c.typeRegistry
		ref.readableTypeNode = firstChild
		return ref

	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF:
		parent := sourceNode.GetParent()
		switch parent.GetToken() {
		case rhino.Token_OBJECTLIT:
			// TODO(b/80580110): Eventually object-literal members should be covered by
			// `PropertyReference`s. However, doing so initially would have caused too many errors
			// in existing code and delayed support for class syntax.
			parentType := parent.GetJSType()
			if !parentType.IsLiteralObject() {
				// Only add a mutation if the object type is actually a literal object (e.g. a
				// global namespace).  OBJECTLIT tokens are often used to fulfill structural types,
				// which is fine for writing constant properties the first time.
				return nil
			}
			ref := &c.scratchPropRef
			ref.sourceNode = sourceNode
			ref.name = sourceNode.GetString()
			ref.receiverType = c.TypeOrUnknownObjectType(rhino.ObjectType_Cast(parentType))
			ref.mutation = true
			ref.declaration = true
			ref.override = false
			ref.readableTypeName = nil
			ref.typeRegistry = c.typeRegistry
			ref.readableTypeNode = parent
			ref.jsdoc = nil
			ref.hasJSDoc = false
			return ref

		case rhino.Token_OBJECT_PATTERN:
			ref := &c.scratchPropRef
			ref.sourceNode = sourceNode
			ref.name = sourceNode.GetString()
			ref.receiverType = c.TypeOrUnknownObjectType(rhino.ObjectType_Cast(parent.GetJSType()))
			ref.mutation = false
			ref.declaration = false
			ref.override = false
			ref.readableTypeName = nil
			ref.typeRegistry = c.typeRegistry
			ref.readableTypeNode = parent
			ref.jsdoc = nil
			ref.hasJSDoc = false
			return ref

		case rhino.Token_CLASS_MEMBERS:
			jsdoc := NodeUtil_GetBestJSDocInfo(sourceNode)
			receiverType := c.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_UNKNOWN_TYPE)
			ctorType := parent.GetParent().GetJSType()
			if ctorType != nil && ctorType.IsFunctionType() {
				ctorFunctionType := ctorType.ToMaybeFunctionType()
				if sourceNode.IsStaticMember() {
					receiverType = ctorFunctionType
				} else if sourceNode.IsMemberFieldDef() {
					receiverType = ctorFunctionType.GetInstanceType()
				} else {
					receiverType = ctorFunctionType.GetPrototype()
				}
			}
			ref := &c.scratchPropRef
			ref.sourceNode = sourceNode
			ref.name = sourceNode.GetString()
			ref.receiverType = receiverType
			ref.mutation = !(sourceNode.IsMemberFieldDef() && !sourceNode.HasChildren())
			ref.declaration = true
			// TODO(b/113704668): This definition is way too loose. It was used to prevent
			// breakages during refactoring and should be tightened.
			ref.override = jsdoc != nil
			ref.readableTypeName = nil
			ref.typeRegistry = nil
			ref.readableTypeNode = nil
			ref.jsdoc = jsdoc
			ref.hasJSDoc = true
			return ref

		default:
			panic(jutil.NewAssertionError("null"))
		}
	default:
		return nil
	}
}

// CheckAccessControls_GetEffectivePropertyVisibility returns the effective
// visibility of the given property. This can differ from the property's
// declared visibility if the property is inherited from a superclass, or if
// the file's @fileoverview JsDoc specifies a default visibility.
//
// TODO(b/111789692): The following methods are forked from
// `AccessControlUtils`. Consider consolidating them.
//
// referenceType is the JavaScript type of the property; fileVisibilityMap a
// map of @fileoverview visibility annotations, used to compute the
// property's default visibility.
//
// Java: CheckAccessControls.getEffectivePropertyVisibility(PropertyReference, ObjectType, ImmutableMap<StaticSourceFile, Visibility>)
func CheckAccessControls_GetEffectivePropertyVisibility(
	propRef *CheckAccessControls_PropertyReference,
	referenceType rhino.ObjectType,
	fileVisibilityMap *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility]) rhino.JSDocInfo_Visibility {
	propertyName := propRef.Name()
	isOverride := propRef.Override()

	definingSource := AccessControlUtils_GetDefiningSource(propRef.SourceNode(), referenceType, propertyName)
	fileOverviewVisibility := checkAccessControls_visibilityFor(fileVisibilityMap, definingSource)
	objectType := AccessControlUtils_GetObjectType(referenceType, isOverride, propertyName)

	if isOverride {
		overridden := AccessControlUtils_GetOverriddenPropertyVisibility(objectType, propertyName)
		return AccessControlUtils_GetEffectiveVisibilityForOverriddenProperty(
			overridden, fileOverviewVisibility, propertyName)
	} else {
		return CheckAccessControls_GetEffectiveVisibilityForNonOverriddenProperty(
			propRef, objectType, fileOverviewVisibility)
	}
}

// CheckAccessControls_GetEffectiveVisibilityForNonOverriddenProperty returns
// the effective visibility of the given non-overridden property.
// Non-overridden properties without an explicit visibility annotation receive
// the default visibility declared in the file's @fileoverview block, if one
// exists.
//
// TODO(b/111789692): The following methods are forked from
// `AccessControlUtils`. Consider consolidating them.
//
// Java: CheckAccessControls.getEffectiveVisibilityForNonOverriddenProperty(PropertyReference, ObjectType, Visibility)
func CheckAccessControls_GetEffectiveVisibilityForNonOverriddenProperty(
	propRef *CheckAccessControls_PropertyReference,
	objectType rhino.ObjectType,
	fileOverviewVisibility *rhino.JSDocInfo_Visibility) rhino.JSDocInfo_Visibility {
	propertyName := propRef.Name()
	raw := rhino.JSDocInfo_Visibility_INHERITED
	if objectType != nil {
		jsdoc := objectType.GetOwnPropertyJSDocInfo(propertyName)
		if jsdoc != nil {
			raw = jsdoc.GetVisibility()
		}
	}
	type_ := propRef.GetJSType()
	createdFromGoogProvide := (type_ != nil && type_.IsLiteralObject())
	// Ignore @fileoverview visibility when computing the effective visibility
	// for properties created by goog.provide.
	//
	// ProcessClosurePrimitives rewrites goog.provide()s as object literal
	// declarations, but the exact form depends on the ordering of the
	// input files. If goog.provide('a.b.c') occurs in the inputs before
	// goog.provide('a'), it is rewritten like
	//
	// var a={};a.b={}a.b.c={};
	//
	// If the file containing goog.provide('a.b.c') also declares
	// a @fileoverview visibility, it must not apply to b, as this would make
	// every a.b.* namespace effectively package-private.
	if raw != rhino.JSDocInfo_Visibility_INHERITED || fileOverviewVisibility == nil || createdFromGoogProvide {
		return raw
	}
	return *fileOverviewVisibility
}
