// Port of com/google/javascript/jscomp/CheckAccessControls.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/AutoBuilder_CheckAccessControls_PropertyReference_Builder.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/CollectFileOverviewVisibility.java @ closure-compiler v20261021

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/ast"
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

	constPropertyInits      map[checkAccessControlsConstPropKey]CheckAccessControls_ConstantDeclaration
	equivConstPropertyInits []checkAccessControlsEquivConstPropEntry

	defaultVisibilityForFiles map[rhino.StaticSourceFile]rhino.JSDocInfo_Visibility

	scratchPropRef           CheckAccessControls_PropertyReference
	lastHoistScopeRoot       rhino.Node
	lastHoistScopeDeprecated bool
	lastVarNode              rhino.Node
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

var checkAccessControlsVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	for tok := rhino.Token(0); tok < rhino.Token_STRING_TYPE; tok++ {
		switch tok {
		case rhino.Token_ROOT,
			rhino.Token_SCRIPT,
			rhino.Token_MODULE_BODY,
			rhino.Token_BLOCK,
			rhino.Token_EXPR_RESULT,
			rhino.Token_LABEL,
			rhino.Token_LABEL_NAME,
			rhino.Token_VAR,
			rhino.Token_LET,
			rhino.Token_CONST,
			rhino.Token_DESTRUCTURING_LHS,
			rhino.Token_ARRAY_PATTERN,
			rhino.Token_OBJECT_PATTERN,
			rhino.Token_IF,
			rhino.Token_SWITCH,
			rhino.Token_SWITCH_BODY,
			rhino.Token_CASE,
			rhino.Token_DEFAULT_CASE,
			rhino.Token_WHILE,
			rhino.Token_DO,
			rhino.Token_FOR,
			rhino.Token_FOR_IN,
			rhino.Token_FOR_OF,
			rhino.Token_FOR_AWAIT_OF,
			rhino.Token_WITH,
			rhino.Token_RETURN,
			rhino.Token_THROW,
			rhino.Token_BREAK,
			rhino.Token_CONTINUE,
			rhino.Token_DEBUGGER,
			rhino.Token_TRY,
			rhino.Token_CATCH,
			rhino.Token_PARAM_LIST,
			rhino.Token_CLASS_MEMBERS,
			rhino.Token_IMPORT,
			rhino.Token_IMPORT_SPECS,
			rhino.Token_IMPORT_SPEC,
			rhino.Token_IMPORT_STAR,
			rhino.Token_EXPORT,
			rhino.Token_EXPORT_SPECS,
			rhino.Token_EXPORT_SPEC,
			rhino.Token_TEMPLATELIT_SUB,
			rhino.Token_TEMPLATELIT_STRING:
			continue
		}
		s.Add(tok)
	}
	return s
}()

var _ nodeTraversal_visitTokens = (*CheckAccessControls)(nil)

func (*CheckAccessControls) visitTokens() *ast.TokenSet {
	return &checkAccessControlsVisitTokens
}

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
	"constant property {0} assigned a value more than once\nInitialized at {1}")

// Java: CheckAccessControls.FINAL_PROPERTY_OVERRIDDEN
var CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN = DiagnosticType_Warning(
	"JSC_FINAL_PROPERTY_OVERRIDDEN",
	"@final method or property {0} overridden\nInitialized at {1}")

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

// String is Enum.name().
func (c CheckAccessControls_Constancy) String() string { return checkAccessControls_Constancy_names[c] }

// Ordinal is Enum.ordinal().
func (c CheckAccessControls_Constancy) Ordinal() int { return int(c) }

// Java: CheckAccessControls.ConstantDeclaration
type CheckAccessControls_ConstantDeclaration struct {
	node       rhino.Node
	annotation CheckAccessControls_Constancy
}

type checkAccessControlsConstPropKey struct {
	owner rhino.JSType
	prop  string
}

type checkAccessControlsEquivConstPropEntry struct {
	owner rhino.JSType
	prop  string
	decl  CheckAccessControls_ConstantDeclaration
}

// Java: CheckAccessControls(AbstractCompiler)
func NewCheckAccessControls(compiler AbstractCompiler) *CheckAccessControls {
	return &CheckAccessControls{
		compiler:     compiler,
		typeRegistry: compiler.GetTypeRegistry(),
	}
}

// Java: CheckAccessControls.process(Node, Node)
func (c *CheckAccessControls) Process(externs rhino.Node, root rhino.Node) {
	collectPass := NewCollectFileOverviewVisibility(c.compiler)
	collectPass.Process(externs, root)
	if collectPass.hasDuplicate {
		panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key"))
	}
	c.defaultVisibilityForFiles = collectPass.visibilities

	sc := NewSyntacticScopeCreator(c.compiler)
	defer func() {
		c.lastVarNode = rhino.NilNode
		c.lastVar = nil
	}()
	b := NodeTraversal_Builder{
		compiler:     c.compiler,
		callback:     c,
		scopeCreator: sc,
	}

	c.lastHoistScopeRoot = rhino.NilNode
	c.lastVarNode = rhino.NilNode
	c.lastVar = nil
	b.Traverse(externs)
	c.lastHoistScopeRoot = rhino.NilNode
	c.lastVarNode = rhino.NilNode
	c.lastVar = nil
	b.Traverse(root)
}

func (c *CheckAccessControls) lookupVarNode(scope *Scope, n rhino.Node) *Var {
	id := n.GetStringID()
	if id == 0 || id == rhino.RhinoStringPool_EmptyID {
		return nil
	}
	return scope.GetVarByID(id)
}

// Java: CheckAccessControls.enterAccessControlScope(Node)
func (c *CheckAccessControls) EnterAccessControlScope(root rhino.Node) {
	scopeType := c.BestInstanceTypeForMethodOrCtor(root)

	if CheckAccessControls_IsMarkedDeprecated(root) {
		c.deprecationDepth++
	}
	c.currentClassStack = append(c.currentClassStack, scopeType)
}

// Java: CheckAccessControls.exitAccessControlScope(Node)
func (c *CheckAccessControls) ExitAccessControlScope(root rhino.Node) {
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
func CheckAccessControls_PrimaryAccessControlScopeRootFor(node rhino.Node) rhino.Node {
	if CheckAccessControls_IsExtendsTarget(node) {
		return node.GetParent()
	} else if CheckAccessControls_IsFunctionOrClass(node) {
		return node
	} else {
		return rhino.NilNode
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
func (c *CheckAccessControls) BestInstanceTypeForMethodOrCtor(n rhino.Node) rhino.ObjectType {
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
			panic(jutil.NewAssertionError(nodeStringOrNull(greatGrandparent)))
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
func (c *CheckAccessControls) ShouldTraverse(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) bool {
	accessControlRoot := CheckAccessControls_PrimaryAccessControlScopeRootFor(node)
	if !accessControlRoot.IsNil() {
		c.EnterAccessControlScope(accessControlRoot)
	}

	return true
}

// Java: CheckAccessControls.visit(NodeTraversal, Node, Node)
func (c *CheckAccessControls) Visit(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) {
	tok := node.GetToken()
	if !checkAccessControlsActiveToken[tok] {
		if tb := node.GetJSTypeBase(); tb == nil || !tb.Self().IsFunctionType() {
			if !parent.IsNil() && parent.IsClass() && parent.GetFirstChild().GetNext() == node {
				c.ExitAccessControlScope(parent)
			}
			return
		}
	} else if tok == rhino.Token_NAME {
		if tb := node.GetJSTypeBase(); tb == nil || !tb.Self().IsFunctionType() || !tb.Self().ToMaybeFunctionType().IsConstructor() {
			if id := node.GetStringID(); id != 0 && id != rhino.RhinoStringPool_EmptyID {
				scope := traversal.GetScope()
				if var_ := scope.GetVarByID(id); var_ != nil {
					if docInfo := var_.GetJSDocInfo(); docInfo != nil && docInfo.IsDeprecated() {
						if c.ShouldEmitDeprecationWarning(traversal, node) {
							name := node.GetString()
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
			if !parent.IsNil() && parent.IsClass() && parent.GetFirstChild().GetNext() == node {
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
	if !accessControlRoot.IsNil() {
		c.ExitAccessControlScope(accessControlRoot)
	}
}

// Java: CheckAccessControls.checkDeprecation(Node, PropertyReference, IdentifierBehaviour, NodeTraversal)
func (c *CheckAccessControls) CheckDeprecation(
	node rhino.Node,
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
	node rhino.Node,
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
func (c *CheckAccessControls) CheckTypeDeprecation(t *NodeTraversal, n rhino.Node) {
	if !c.ShouldEmitDeprecationWarning(t, n) {
		return
	}

	instanceType := n.GetJSType().ToMaybeFunctionType().GetInstanceType()

	deprecationInfo, ok := CheckAccessControls_GetTypeDeprecationInfo(instanceType)
	if !ok {
		return
	}

	message := CheckAccessControls_DEPRECATED_CLASS_REASON
	if deprecationInfo == "" {
		message = CheckAccessControls_DEPRECATED_CLASS
	}
	c.compiler.Report(JSError_MakeWithNode(n, message, instanceType.String(), deprecationInfo))
}

// CheckNameDeprecation checks the given NAME node to ensure that access
// restrictions are obeyed.
//
// Java: CheckAccessControls.checkNameDeprecation(NodeTraversal, Node)
func (c *CheckAccessControls) CheckNameDeprecation(t *NodeTraversal, n rhino.Node) {
	if !n.IsName() {
		return
	}

	var_ := c.lookupVarNode(t.GetScope(), n)
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

// CheckNameVisibility reports an error if the given name is not visible in
// the current context.
//
// scope is the current scope; name the name node.
//
// Java: CheckAccessControls.checkNameVisibility(Scope, Node)
func (c *CheckAccessControls) CheckNameVisibility(scope *Scope, name rhino.Node) {
	if !name.IsName() {
		return
	}

	var var_ *Var
	if c.lastVarNode == name {
		var_ = c.lastVar
	} else {
		var_ = c.lookupVarNode(scope, name)
	}
	if var_ == nil {
		return
	}
	c.checkNameVisibilityWithVar(var_, name)
}

func (c *CheckAccessControls) checkNameVisibilityWithVar(var_ *Var, name rhino.Node) {
	var v rhino.JSDocInfo_Visibility
	if len(c.defaultVisibilityForFiles) == 0 {
		jsDocInfo := var_.GetJSDocInfo()
		if jsDocInfo == nil {
			return
		}
		v = jsDocInfo.GetVisibility()
	} else {
		v = c.effectiveNameVisibility(name, var_)
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

func (c *CheckAccessControls) effectiveNameVisibility(name rhino.Node, var_ *Var) rhino.JSDocInfo_Visibility {
	raw := rhino.JSDocInfo_Visibility_INHERITED
	if jsDocInfo := var_.GetJSDocInfo(); jsDocInfo != nil {
		raw = jsDocInfo.GetVisibility()
	}
	if raw != rhino.JSDocInfo_Visibility_INHERITED {
		return raw
	}
	defaultVisibilityForFile, hasDefault := c.defaultVisibilityForFiles[var_.GetSourceFile()]
	if !hasDefault {
		return raw
	}
	if type_ := name.GetJSType(); type_ != nil && type_.IsLiteralObject() {
		return raw
	}
	return defaultVisibilityForFile
}

// Java: CheckAccessControls.isPrivateAccessAllowed(Var, Node)
func CheckAccessControls_IsPrivateAccessAllowed(var_ *Var, name rhino.Node) bool {
	varSrc := var_.GetSourceFile()
	refSrc := name.GetStaticSourceFile()

	return varSrc == nil || refSrc == nil || varSrc.GetName() == refSrc.GetName()
}

// Java: CheckAccessControls.isPackageAccessAllowed(Var, Node)
func (c *CheckAccessControls) IsPackageAccessAllowed(var_ *Var, name rhino.Node) bool {
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

// CheckFinalClassOverrides checks if a constructor is trying to override a
// final class.
//
// Java: CheckAccessControls.checkFinalClassOverrides(Node)
func (c *CheckAccessControls) CheckFinalClassOverrides(ctor rhino.Node) {
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
				stringPtrOrNull(ctorType.GetDisplayName()),
				stringPtrOrNull(finalParentClass.GetDisplayName())))
	}
}

// ShouldEmitDeprecationWarning determines whether a deprecation warning
// should be emitted.
//
// t is the current traversal; n the node which we are checking.
//
// Java: CheckAccessControls.shouldEmitDeprecationWarning(NodeTraversal, Node)
func (c *CheckAccessControls) ShouldEmitDeprecationWarning(t *NodeTraversal, n rhino.Node) bool {
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
	if hoistRoot == c.lastHoistScopeRoot && !hoistRoot.IsNil() {
		return c.lastHoistScopeDeprecated
	}
	scopeRoot := hoistRoot
	if NodeUtil_IsFunctionBlock(scopeRoot) {
		scopeRoot = scopeRoot.GetParent()
	}
	scopeRootParent := scopeRoot.GetParent()

	// Cases 2 and 3 are required to handle ES5-style class methods since they aren't nested inside
	// their class. This is tested in the CheckAccessControlsOldSyntaxTest class.
	_, thisDeprecated := CheckAccessControls_GetTypeDeprecationInfo(c.GetTypeOfThis(scopeRoot))
	res := thisDeprecated
	if !res && !scopeRootParent.IsNil() && scopeRootParent.IsAssign() {
		_, res = CheckAccessControls_GetTypeDeprecationInfo(c.BestInstanceTypeForMethodOrCtor(scopeRoot))
	}
	c.lastHoistScopeRoot = hoistRoot
	c.lastHoistScopeDeprecated = res
	return res
}

// CheckAccessControls_IsMarkedDeprecated returns whether this node roots a
// subtree under which references to deprecated constructs are allowed.
//
// Java: CheckAccessControls.isMarkedDeprecated(Node)
func CheckAccessControls_IsMarkedDeprecated(n rhino.Node) bool {
	_, ok := CheckAccessControls_GetDeprecationReason(NodeUtil_GetBestJSDocInfo(n))
	return ok
}

// CheckAccessControls_GetTypeDeprecationInfo returns the deprecation reason
// for the type if it is marked as being deprecated. Returns ("", true) if
// the type is deprecated without a reason, or ("", false) if not deprecated.
//
// Java: CheckAccessControls.getTypeDeprecationInfo(JSType)
func CheckAccessControls_GetTypeDeprecationInfo(type_ rhino.JSType) (string, bool) {
	if type_ == nil {
		return "", false
	}
	return CheckAccessControls_GetDeprecationReason(type_.GetJSDocInfo())
}

// Java: CheckAccessControls.getDeprecationReason(JSDocInfo)
func CheckAccessControls_GetDeprecationReason(info *rhino.JSDocInfo) (string, bool) {
	if info != nil && info.IsDeprecated() {
		if reason := info.GetDeprecationReason(); reason != nil {
			return *reason, true
		}
		return "", true
	}
	return "", false
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
func CheckAccessControls_IsFunctionOrClass(n rhino.Node) bool {
	return n.IsFunction() || n.IsClass()
}

// Java: CheckAccessControls.isExtendsTarget(Node)
func CheckAccessControls_IsExtendsTarget(node rhino.Node) bool {
	return ast.RoleOf(node) == ast.RoleClassExtends
}

// Java: CheckAccessControls.getTypeOfThis(Node)
func (c *CheckAccessControls) GetTypeOfThis(scopeRoot rhino.Node) rhino.JSType {
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

// String is Enum.name().
func (b CheckAccessControls_IdentifierBehaviour) String() string {
	return checkAccessControls_IdentifierBehaviour_names[b]
}

// Ordinal is Enum.ordinal().
func (b CheckAccessControls_IdentifierBehaviour) Ordinal() int { return int(b) }

// Java: CheckAccessControls.IdentifierBehaviour.select(Node)
func CheckAccessControls_IdentifierBehaviour_Select(target rhino.Node) CheckAccessControls_IdentifierBehaviour {
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
	isEs6 := (!ctorType.GetSource().IsNil()) && ctorType.GetSource().IsClass()

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

// CollectFileOverviewVisibility is the compiler pass that collects
// visibility annotations in @fileoverview blocks. Used by
// CheckAccessControls.
type CollectFileOverviewVisibility struct {
	visibilities map[rhino.StaticSourceFile]rhino.JSDocInfo_Visibility
	order        []rhino.StaticSourceFile
	hasDuplicate bool
}

// Java: CollectFileOverviewVisibility(AbstractCompiler)
func NewCollectFileOverviewVisibility(compiler AbstractCompiler) *CollectFileOverviewVisibility {
	return &CollectFileOverviewVisibility{}
}

// Java: CollectFileOverviewVisibility.process(Node, Node)
func (c *CollectFileOverviewVisibility) Process(externs rhino.Node, root rhino.Node) {
	for script := root.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		jutil.CheckState(script.IsScript())
		c.Visit(script)
	}
}

// Java: CollectFileOverviewVisibility.visit(Node)
func (c *CollectFileOverviewVisibility) Visit(scriptNode rhino.Node) {
	jsDocInfo := scriptNode.GetJSDocInfo()
	if jsDocInfo == nil {
		return
	}
	// Go-only: language-adapter — JSDocInfo.getVisibility() never returns
	// null in Java either (it defaults to INHERITED), so Java's
	// `if (v == null) return;` never fires and has no Go counterpart.
	v := jsDocInfo.GetVisibility()
	sf := jutil.CheckNotNull(scriptNode.GetStaticSourceFile())
	if c.visibilities == nil {
		c.visibilities = make(map[rhino.StaticSourceFile]rhino.JSDocInfo_Visibility)
		c.order = append(c.order, sf)
	} else if _, exists := c.visibilities[sf]; exists {
		c.hasDuplicate = true
	} else {
		c.order = append(c.order, sf)
	}
	c.visibilities[sf] = v
}

// Java: CollectFileOverviewVisibility.getFileOverviewVisibilityMap()
func (c *CollectFileOverviewVisibility) GetFileOverviewVisibilityMap() *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility] {
	if c.hasDuplicate {
		panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key"))
	}
	m := jutil.NewLinkedHashMapWithCapacity[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility](len(c.order))
	for _, sf := range c.order {
		m.Put(sf, c.visibilities[sf])
	}
	return m
}
