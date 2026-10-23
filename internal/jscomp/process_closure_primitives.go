// Port of com/google/javascript/jscomp/ProcessClosurePrimitives.java @ closure-compiler v20261015

/*
 * Copyright 2006 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ProcessClosurePrimitives performs some Closure-specific simplifications
// including rewriting goog.base, goog.addDependency.
//
// Records forwardDeclared names in the internal compiler state.
type ProcessClosurePrimitives struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler

	knownClosureSubclasses *jutil.LinkedHashSet[string]

	closureModules *jutil.LinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]
}

var ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR = DiagnosticType_Error(
	"JSC_EXPECTED_OBJECTLIT_ERROR",
	"method \"{0}\" expected an object literal argument")

var ProcessClosurePrimitives_EXPECTED_STRING_ERROR = DiagnosticType_Error(
	"JSC_EXPECTED_STRING_ERROR", "method \"{0}\" expected a string argument")

var ProcessClosurePrimitives_INVALID_STYLE_ERROR = DiagnosticType_Error(
	"JSC_INVALID_CSS_NAME_MAP_STYLE_ERROR",
	"Invalid CSS name map style {0}")

var ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE = DiagnosticType_Warning(
	"JSC_WEAK_NAMESPACE_TYPE",
	"Provided symbol declared with type Object. This is rarely useful. "+
		"For more information see "+
		"https://github.com/google/closure-compiler/wiki/A-word-about-the-type-Object")

var ProcessClosurePrimitives_CLASS_NAMESPACE_ERROR = DiagnosticType_Error(
	"JSC_CLASS_NAMESPACE_ERROR",
	"\"{0}\" cannot be both provided and declared as a class. Try var {0} = class '{'...'}'"+
		" (metadata {1})")

var ProcessClosurePrimitives_FUNCTION_NAMESPACE_ERROR = DiagnosticType_Error(
	"JSC_FUNCTION_NAMESPACE_ERROR",
	"\"{0}\" cannot be both provided and declared as a function. (metadata {1})")

var ProcessClosurePrimitives_INVALID_PROVIDE_ERROR = DiagnosticType_Error(
	"JSC_INVALID_PROVIDE_ERROR",
	"\"{0}\" is not a valid {1} qualified name")

var ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR = DiagnosticType_Error(
	"JSC_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR",
	"goog.setCssNameMapping only takes an object literal with string values")

var ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP = DiagnosticType_Warning(
	"INVALID_CSS_RENAMING_MAP",
	"Invalid entries in css renaming map: {0}")

var ProcessClosurePrimitives_BASE_CLASS_ERROR = DiagnosticType_Error(
	"JSC_BASE_CLASS_ERROR",
	"incorrect use of {0}.base: {1}")

var ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR = DiagnosticType_Error(
	"JSC_POSSIBLE_BASE_CLASS_ERROR",
	"potentially incorrect use of {0}.base: {1}\n"+
		"Note: if this .base method is not on a Closure subclass, this error is a false"+
		" positive. Suppress with /** @suppress '{closureClassChecks}' */")

var ProcessClosurePrimitives_INVALID_FORWARD_DECLARE = DiagnosticType_Error("JSC_INVALID_FORWARD_DECLARE", "Malformed goog.forwardDeclare")

var ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR = DiagnosticType_Error(
	"JSC_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR",
	"Closure primitive method {0} may not be aliased")

var ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR = DiagnosticType_Error(
	"JSC_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR",
	"Closure primitive method {0} may not be aliased  outside a module (ES "+
		"module, CommonJS module, or goog.module)")

var ProcessClosurePrimitives_INVALID_RENAME_FUNCTION = DiagnosticType_Error("JSC_INVALID_RENAME_FUNCTION", "{0} call is invalid: {1}")

var ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL = DiagnosticType_Error("JSC_INVALID_GOOG_WEAK_USAGE", "{0} call is invalid: {1}")

// ProcessClosurePrimitives_GOOG is the root Closure namespace
const ProcessClosurePrimitives_GOOG = "goog"

var processClosurePrimitives_GOOG_INHERITS = rhino.QualifiedName_Of("goog.inherits")

// Java: ProcessClosurePrimitives(AbstractCompiler)
func NewProcessClosurePrimitives(compiler AbstractCompiler) *ProcessClosurePrimitives {
	p := &ProcessClosurePrimitives{}
	p.knownClosureSubclasses = jutil.NewLinkedHashSet[string]()
	p.compiler = compiler
	p.closureModules =
		jutil.CheckNotNullMsg(compiler.GetModuleMetadataMap(), "Need to run GatherModuleMetadata").
			GetModulesByGoogNamespace()
	return p
}

// Java: ProcessClosurePrimitives.process(Node, Node)
func (p *ProcessClosurePrimitives) Process(externs *rhino.Node, root *rhino.Node) {
	// Replace and validate other Closure primitives
	NodeTraversal_TraverseRoots(p.compiler, p, externs, root)
}

// Java: ProcessClosurePrimitives.visit(NodeTraversal, Node, Node)
func (p *ProcessClosurePrimitives) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CALL:
		{
			p.CheckGoogFunctions(t, n)
			p.MaybeProcessClassBaseCall(n)
			p.CheckPropertyRenameCall(n)
		}
	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		if !t.InGlobalHoistScope() || n.IsFromExterns() {
			return
		}
		if !NodeUtil_IsFunctionDeclaration(n) && !NodeUtil_IsClassDeclaration(n) {
			return
		}
		name := n.GetFirstChild().GetString()
		pn := p.closureModules.Get(name)
		if pn != nil && pn.IsGoogProvide() {
			diagnostic := ProcessClosurePrimitives_FUNCTION_NAMESPACE_ERROR
			if n.IsClass() {
				diagnostic = ProcessClosurePrimitives_CLASS_NAMESPACE_ERROR
			}
			p.compiler.Report(
				JSError_MakeWithNode(
					n,
					diagnostic,
					name,
					pn.String()))
		}
	case rhino.Token_EXPR_RESULT:
		assign := n.GetFirstChild()
		if !assign.IsAssign() || n.IsFromExterns() {
			break
		}
		info := assign.GetJSDocInfo()
		if info == nil || info.GetType() == nil {
			break
		}
		lhs, ok := assign.GetFirstChild().GetQualifiedNameString()
		if !ok {
			break
		}
		p.CheckPossibleGoogProvideInit(lhs, info, n)
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		if !n.GetFirstChild().IsName() {
			break
		}
		p.CheckPossibleGoogProvideInit(n.GetFirstChild().GetString(), n.GetJSDocInfo(), n)
	default:
	}
}

// Java: ProcessClosurePrimitives.checkPossibleGoogProvideInit(String, JSDocInfo, Node)
func (p *ProcessClosurePrimitives) CheckPossibleGoogProvideInit(
	namespace string, info *rhino.JSDocInfo, definition *rhino.Node) {
	if info == nil || definition.IsFromExterns() {
		return
	}

	metadata := p.closureModules.Get(namespace)
	if metadata == nil || !metadata.IsGoogProvide() {
		return
	}

	// Validate that the namespace is not declared as a generic object type.
	expr := info.GetType()
	if expr == nil {
		return
	}
	n := expr.GetRoot()
	if n.GetToken() == rhino.Token_BANG {
		n = n.GetFirstChild()
	}
	if n.IsStringLit() &&
		!n.HasChildren() && // templated object types are ok.
		n.GetString() == "Object" {
		p.compiler.Report(JSError_MakeWithNode(definition, ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE))
	}
}

// Java: ProcessClosurePrimitives.checkGoogFunctions(NodeTraversal, Node)
func (p *ProcessClosurePrimitives) CheckGoogFunctions(t *NodeTraversal, call *rhino.Node) {
	callee := call.GetFirstChild()
	if !callee.IsGetProp() {
		return
	}

	receiver := callee.GetFirstChild()
	if !receiver.IsName() || receiver.GetString() != ProcessClosurePrimitives_GOOG {
		return
	}

	// For the sake of simplicity, we report code changes
	// when we see a provides/requires, and don't worry about
	// reporting the change when we actually do the replacement.
	methodName := callee.GetString()
	switch methodName {
	case "inherits":
		// Note: inherits is allowed in local scope
		p.ProcessInheritsCall(call)
	case "addDependency":
		if p.ValidateUnaliasablePrimitiveCall(t, call, methodName) {
			p.ProcessAddDependency(call)
		}
	case "setCssNameMapping":
		_ = ProcessClosurePrimitives_ProcessSetCssNameMapping(p.compiler, call, call.GetParent())
	case "forwardDeclare":
		if p.ValidatePrimitiveCallWithMessage(
			t,
			call,
			methodName,
			ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR) {
			p.ProcessForwardDeclare(call)
		}
	case "weakUsage":
		p.ValidateWeakUsageCall(call)
	default:
	}
}

// Java: ProcessClosurePrimitives.validateWeakUsageCall(Node)
func (p *ProcessClosurePrimitives) ValidateWeakUsageCall(call *rhino.Node) {
	// goog.weakUsage() should have exactly one argument, and it should be a name (possibly
	// qualified).
	childCount := call.GetChildCount()
	arg := call.GetSecondChild()
	if childCount != 2 {
		calleeName, ok := call.GetFirstChild().GetQualifiedNameString()
		if !ok {
			calleeName = "null"
		}
		p.compiler.Report(
			JSError_MakeWithNode(
				call,
				ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL,
				calleeName,
				"should have exactly one argument, not "+strconv.Itoa(childCount-1)))
	}
	if childCount >= 2 && !call.GetSecondChild().IsQualifiedName() {
		calleeName, ok := call.GetFirstChild().GetQualifiedNameString()
		if !ok {
			calleeName = "null"
		}
		p.compiler.Report(
			JSError_MakeWithNode(
				call,
				ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL,
				calleeName,
				"argument should be a name or qualified name, not "+arg.String()))
	}
}

// ValidateUnaliasablePrimitiveCall verifies that a) the call is in the top
// level of a file and b) the return value is unused
//
// This method is for primitives that never return a value.
//
// Java: ProcessClosurePrimitives.validateUnaliasablePrimitiveCall(NodeTraversal, Node, String)
func (p *ProcessClosurePrimitives) ValidateUnaliasablePrimitiveCall(t *NodeTraversal, n *rhino.Node, methodName string) bool {
	return p.ValidatePrimitiveCallWithMessage(t, n, methodName, ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR)
}

// ValidatePrimitiveCallWithMessage: methodName is the list of primitive
// types classed together with this one; invalidAliasingError is which
// DiagnosticType to emit if this call is aliased. this depends on whether
// the primitive is sometimes aliasiable in a module or never aliasable.
//
// Java: ProcessClosurePrimitives.validatePrimitiveCallWithMessage(NodeTraversal, Node, String, DiagnosticType)
func (p *ProcessClosurePrimitives) ValidatePrimitiveCallWithMessage(
	t *NodeTraversal, n *rhino.Node, methodName string, invalidAliasingError *DiagnosticType) bool {
	// Ignore invalid primitives if we didn't strip module sugar.
	if p.compiler.GetOptions().ShouldPreserveGoogModule() {
		return true
	}

	if !t.InGlobalHoistScope() && !t.InModuleScope() {
		p.compiler.Report(JSError_MakeWithNode(n, ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR))
		return false
	} else if !n.GetParent().IsExprResult() && !t.InModuleScope() {
		// If the call is in the global hoist scope, but the result is used
		p.compiler.Report(JSError_MakeWithNode(n, invalidAliasingError, ProcessClosurePrimitives_GOOG+"."+methodName))
		return false
	}
	return true
}

// Java: ProcessClosurePrimitives.maybeProcessClassBaseCall(Node)
func (p *ProcessClosurePrimitives) MaybeProcessClassBaseCall(n *rhino.Node) {
	// Two things must hold for every base call:
	// 1) We must be calling it on "this".
	// 2) We must be calling it on a prototype method of the same name as
	//    the one we're in, OR we must be calling it from a constructor.
	// If both of those things are true, then we can rewrite:
	// <pre>
	// function Foo() {
	//   Foo.base(this);
	// }
	// goog.inherits(Foo, BaseFoo);
	// Foo.prototype.bar = function() {
	//   Foo.base(this, 'bar', 1);
	// };
	// </pre>
	// as the easy-to-optimize:
	// <pre>
	// function Foo() {
	//   BaseFoo.call(this);
	// }
	// goog.inherits(Foo, BaseFoo);
	// Foo.prototype.bar = function() {
	//   Foo.superClass_.bar.call(this, 1);
	// };
	//
	// Most of the logic here is just to make sure the AST's
	// structure is what we expect it to be.

	callTarget := n.GetFirstChild()
	if !callTarget.IsGetProp() || callTarget.GetString() != "base" {
		return
	}

	baseContainerNode := callTarget.GetFirstChild()
	if !baseContainerNode.IsUnscopedQualifiedName() {
		// Some unknown "base" method.
		return
	}
	baseContainer, _ := baseContainerNode.GetQualifiedNameString()

	enclosingFnNameNode := ProcessClosurePrimitives_GetEnclosingDeclNameNode(n)
	if enclosingFnNameNode == nil || !enclosingFnNameNode.IsUnscopedQualifiedName() {
		// some unknown container method or a MEMBER_FUNCTION_DEF.
		if p.knownClosureSubclasses.Contains(baseContainer) {
			p.ReportBadBaseMethodUse(n, baseContainer, "Could not find enclosing method.")
		} else if ProcessClosurePrimitives_BaseUsedInClass(n) {
			clazz := NodeUtil_GetEnclosingClass(n)
			// TODO(lharker): this check ignores class expressions like "foo.Bar = class extends X {}"
			if (clazz.GetFirstChild().IsName() &&
				clazz.GetFirstChild().GetString() == baseContainer) ||
				(clazz.GetSecondChild().IsName() &&
					clazz.GetSecondChild().GetString() == baseContainer) {
				p.ReportBadBaseMethodUse(
					n,
					clazz.GetFirstChild().GetString(),
					"base method is not allowed in ES6 class. Use super instead.")
			}
		}
		return
	}

	if ProcessClosurePrimitives_BaseUsedInClass(n) {
		clazz := NodeUtil_GetEnclosingClass(n)
		name := NodeUtil_GetBestLValueName(clazz)
		p.ReportBadBaseMethodUse(
			n, jutil.ToJavaString(name), "base method is not allowed in ES6 class. Use super instead.")
		return
	}

	enclosingQname, _ := enclosingFnNameNode.GetQualifiedNameString()
	if !(jutil.StringIndexOf(enclosingQname, ".prototype.") >= 0) {
		p.RewriteBaseCallInConstructor(enclosingQname, baseContainer, n, enclosingFnNameNode)
	} else {
		p.RewriteBaseCallInMethod(enclosingQname, baseContainer, n, enclosingFnNameNode)
	}
}

// Java: ProcessClosurePrimitives.rewriteBaseCallInConstructor(String, String, Node, Node)
func (p *ProcessClosurePrimitives) RewriteBaseCallInConstructor(
	enclosingQname string, baseContainer string, n *rhino.Node, enclosingFnNameNode *rhino.Node) {
	// Check if this is some other "base" method.
	if enclosingQname != baseContainer {
		// Report misuse of "base" methods from other known classes.
		if p.knownClosureSubclasses.Contains(baseContainer) {
			p.ReportBadBaseMethodUse(
				n, baseContainer, "Must be used within "+baseContainer+" methods")
		}
		return
	}

	// Determine if this is a class with a "base" method created by
	// goog.inherits.
	enclosingParent := enclosingFnNameNode.GetParent()
	ctorDeclarationNode := enclosingParent
	if enclosingParent.IsAssign() {
		ctorDeclarationNode = enclosingParent.GetParent()
	}
	baseClassNode := ProcessClosurePrimitives_FindGoogInheritsCall(ctorDeclarationNode)
	if baseClassNode == nil {
		// If there is no "goog.inherits", this might be some other "base" method.
		return
	}

	// This is the expected method, validate its parameters.
	callee := n.GetFirstChild()
	thisArg := callee.GetNext()
	if thisArg == nil || !thisArg.IsThis() {
		p.ReportBadBaseMethodUse(n, baseContainer, "First argument must be 'this'.")
		return
	}

	// Handle methods.
	methodNameNode := thisArg.GetNext()
	if methodNameNode == nil ||
		!methodNameNode.IsStringLit() ||
		methodNameNode.GetString() != "constructor" {
		p.ReportBadBaseMethodUse(n, baseContainer, "Second argument must be 'constructor'.")
		return
	}

	// We're good to go.
	baseClassQname, ok := baseClassNode.GetQualifiedNameString()
	if !ok {
		baseClassQname = "null"
	}
	callee.ReplaceWith(
		NodeUtil_NewQNameWithOriginalName(
			p.compiler,
			baseClassQname+".call",
			callee,
			enclosingQname+".base"))
	methodNameNode.Detach()
	p.compiler.ReportChangeToEnclosingScope(n)
}

// RewriteBaseCallInMethod: enclosingQname is a Java String that is never
// null here (a Go string is never null).
//
// Java: ProcessClosurePrimitives.rewriteBaseCallInMethod(String, String, Node, Node)
func (p *ProcessClosurePrimitives) RewriteBaseCallInMethod(
	enclosingQname string, baseContainer string, n *rhino.Node, enclosingFnNameNode *rhino.Node) {
	if !p.knownClosureSubclasses.Contains(baseContainer) {
		// Can't determine if this is a known "class" that has a known "base" method. Don't rewrite
		// the "base" method call because if it's not a Closure subclass, the rewriting would be
		// wrong. Instead emit an error the user can either suppress or fix.
		p.ReportPossibleBadBaseMethodUse(
			n,
			baseContainer,
			"base used outside a known Closure subclass, in: "+
				baseContainer+
				"\nIf "+
				baseContainer+
				" is actually a Closure subclass, the compiler is unable to statically determine"+
				" this. This can happen when calling .base off of some import or other alias of the"+
				" original Closure subclass, such as from a goog.require in another file. Instead"+
				" of using .base, rewrite this call as ParentClass.prototype.method.call(this)")
		return
	}

	misuseOfBase :=
		!enclosingFnNameNode.GetFirstFirstChild().MatchesQualifiedName(baseContainer)
	if misuseOfBase {
		// Report misuse of "base" methods from other known classes.
		p.ReportBadBaseMethodUse(n, baseContainer, "Must be used within "+baseContainer+" methods")
		return
	}

	// The super class is known.
	callee := n.GetFirstChild()
	thisArg := callee.GetNext()
	if thisArg == nil || !thisArg.IsThis() {
		p.ReportBadBaseMethodUse(n, baseContainer, "First argument must be 'this'.")
		return
	}

	// Handle methods.
	methodNameNode := thisArg.GetNext()
	if methodNameNode == nil || !methodNameNode.IsStringLit() {
		p.ReportBadBaseMethodUse(n, baseContainer, "Second argument must name a method.")
		return
	}

	methodName := methodNameNode.GetString()
	ending := ".prototype." + methodName
	if !jutil.StringEndsWith(enclosingQname, ending) {
		p.ReportBadBaseMethodUse(n, baseContainer, "Enclosing method does not match "+methodName)
		return
	}

	// We're good to go.
	className := enclosingFnNameNode.GetFirstFirstChild()
	classQname, ok := className.GetQualifiedNameString()
	if !ok {
		classQname = "null"
	}
	callee.ReplaceWith(
		NodeUtil_NewQNameWithOriginalName(
			p.compiler,
			classQname+".superClass_."+methodName+".call",
			callee,
			enclosingQname+".base"))
	methodNameNode.Detach()
	p.compiler.ReportChangeToEnclosingScope(n)
}

// Java: ProcessClosurePrimitives.findGoogInheritsCall(Node)
func ProcessClosurePrimitives_FindGoogInheritsCall(ctorDeclarationNode *rhino.Node) *rhino.Node {
	maybeInheritsExpr := ctorDeclarationNode.GetNext()

	for maybeInheritsExpr != nil &&
		(maybeInheritsExpr.IsEmpty() || NodeUtil_IsExprAssign(maybeInheritsExpr)) {
		// Skip empty nodes (e.g. those created by an unnecessary semicolon) and potential aliases,
		// e.g. "exports = Ctor;" in a goog.module.
		maybeInheritsExpr = maybeInheritsExpr.GetNext()
	}
	var baseClassNode *rhino.Node = nil
	if maybeInheritsExpr != nil && NodeUtil_IsExprCall(maybeInheritsExpr) {
		callNode := maybeInheritsExpr.GetFirstChild()
		if processClosurePrimitives_GOOG_INHERITS.Matches(callNode.GetFirstChild()) &&
			callNode.GetLastChild().IsQualifiedName() {
			baseClassNode = callNode.GetLastChild()
		}
	}
	return baseClassNode
}

// ProcessInheritsCall processes the goog.inherits call.
//
// Java: ProcessClosurePrimitives.processInheritsCall(Node)
func (p *ProcessClosurePrimitives) ProcessInheritsCall(n *rhino.Node) {
	if n.HasXChildren(3) {
		subClass := n.GetSecondChild()
		superClass := subClass.GetNext()
		if subClass.IsUnscopedQualifiedName() && superClass.IsUnscopedQualifiedName() {
			if subQname, ok := subClass.GetQualifiedNameString(); ok {
				p.knownClosureSubclasses.Add(subQname)
			}
		}
	}
}

// ProcessClosurePrimitives_GetEnclosingDeclNameNode returns the qualified
// name node of the function whose scope we're in, or null if it cannot be
// found.
//
// Java: ProcessClosurePrimitives.getEnclosingDeclNameNode(Node)
func ProcessClosurePrimitives_GetEnclosingDeclNameNode(n *rhino.Node) *rhino.Node {
	fn := NodeUtil_GetEnclosingFunction(n)
	if fn == nil {
		return nil
	}
	return NodeUtil_GetNameNode(fn)
}

// ProcessClosurePrimitives_BaseUsedInClass verifies if goog.base call is
// used in a class
//
// Java: ProcessClosurePrimitives.baseUsedInClass(Node)
func ProcessClosurePrimitives_BaseUsedInClass(n *rhino.Node) bool {
	for curr := n; curr != nil; curr = curr.GetParent() {
		if curr.IsClassMembers() {
			return true
		}
	}
	return false
}

// ReportBadBaseMethodUse reports an incorrect use of super-method calling.
//
// Java: ProcessClosurePrimitives.reportBadBaseMethodUse(Node, String, String)
func (p *ProcessClosurePrimitives) ReportBadBaseMethodUse(n *rhino.Node, className string, extraMessage string) {
	p.compiler.Report(JSError_MakeWithNode(n, ProcessClosurePrimitives_BASE_CLASS_ERROR, className, extraMessage))
}

// ReportPossibleBadBaseMethodUse reports a potential incorrect use of
// super-method calling.
//
// Use this instead of reportBadBaseMethodUse when reporting a potential
// false positive that should be suppressible in JS code, instead of a
// coding pattern that is definitely bad and is not suppressible.
//
// Java: ProcessClosurePrimitives.reportPossibleBadBaseMethodUse(Node, String, String)
func (p *ProcessClosurePrimitives) ReportPossibleBadBaseMethodUse(n *rhino.Node, className string, extraMessage string) {
	p.compiler.Report(JSError_MakeWithNode(n, ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR, className, extraMessage))
}

// processClosurePrimitives_cssRenamingMap is the anonymous CssRenamingMap of
// processSetCssNameMapping.
type processClosurePrimitives_cssRenamingMap struct {
	cssNames *jutil.LinkedHashMap[string, string]
	style    CssRenamingMap_Style
}

// Java: ProcessClosurePrimitives.processSetCssNameMapping(...).new CssRenamingMap() {...}.get(String)
func (m *processClosurePrimitives_cssRenamingMap) Get(value string) *string {
	if m.cssNames.ContainsKey(value) {
		s := m.cssNames.Get(value)
		return &s
	} else {
		return &value
	}
}

// Java: ProcessClosurePrimitives.processSetCssNameMapping(...).new CssRenamingMap() {...}.getStyle()
func (m *processClosurePrimitives_cssRenamingMap) GetStyle() CssRenamingMap_Style {
	return m.style
}

// ProcessClosurePrimitives_ProcessSetCssNameMapping processes a call to
// goog.setCssNameMapping(). Either the argument to goog.setCssNameMapping()
// is valid, in which case it will be used to create a CssRenamingMap for the
// compiler of this CompilerPass, or it is invalid and a JSCompiler error
// will be reported.
//
// Java: ProcessClosurePrimitives.processSetCssNameMapping(AbstractCompiler, Node, Node)
func ProcessClosurePrimitives_ProcessSetCssNameMapping(
	compiler AbstractCompiler, n *rhino.Node, parent *rhino.Node) CssRenamingMap {
	left := n.GetFirstChild()
	arg := left.GetNext()
	if !ProcessClosurePrimitives_VerifySetCssNameMapping(compiler, left, arg) {
		return nil
	}
	// Translate OBJECTLIT into SubstitutionMap. All keys and
	// values must be strings, or an error will be thrown.
	cssNames := jutil.NewLinkedHashMap[string, string]()

	for key := arg.GetFirstChild(); key != nil; key = key.GetNext() {
		value := key.GetFirstChild()
		if !key.IsStringKey() || value == nil || !value.IsStringLit() {
			compiler.Report(JSError_MakeWithNode(n, ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR))
			return nil
		}
		cssNames.Put(key.GetString(), value.GetString())
	}

	styleStr := "BY_PART"
	if arg.GetNext() != nil {
		styleStr = arg.GetNext().GetString()
	}

	var style CssRenamingMap_Style
	invalidStyle := false
	jutil.Try(func() {
		style = CssRenamingMap_Style_ValueOf(styleStr)
	}, func(e jutil.IIllegalArgumentException) {
		compiler.Report(JSError_MakeWithNode(n, ProcessClosurePrimitives_INVALID_STYLE_ERROR, styleStr))
		invalidStyle = true
	})
	if invalidStyle {
		return nil
	}

	if style == CssRenamingMap_Style_BY_PART {
		// Make sure that no class keys contain -'s
		// Variable keys can contain -'s since they're always named BY_WHOLE.
		errors := jutil.NewArrayList[string]()
		for key := range cssNames.Keys() {
			if jutil.StringIndexOf(key, "-") >= 0 && !jutil.StringStartsWith(key, "--") {
				errors.Add(key)
			}
		}
		if !errors.IsEmpty() {
			compiler.Report(JSError_MakeWithNode(n, ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP, errors.String()))
		}
	} else if style == CssRenamingMap_Style_BY_WHOLE {
		// Verifying things is a lot trickier here. We just do a quick
		// n^2 check over the map which makes sure that if "a-b" in
		// the map, then map(a-b) = map(a)-map(b).
		// To speed things up, only consider cases where len(b) <= 10
		errors := jutil.NewArrayList[string]()
		for bKey, bValue := range cssNames.All() {
			if jutil.StringLength(bKey) > 10 {
				continue
			}
			for aKey, aValue := range cssNames.All() {
				combined, ok := cssNames.Lookup(aKey + "-" + bKey)
				if ok && combined != aValue+"-"+bValue {
					errors.Add(
						"map(" +
							aKey +
							"-" +
							bKey +
							") != map(" +
							aKey +
							")-map(" +
							bKey +
							")")
				}
			}
		}
		if !errors.IsEmpty() {
			compiler.Report(JSError_MakeWithNode(n, ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP, errors.String()))
		}
	}

	return &processClosurePrimitives_cssRenamingMap{cssNames: cssNames, style: style}
}

// ProcessAddDependency processes a goog.addDependency() call and record any
// forward declarations.
//
// Java: ProcessClosurePrimitives.processAddDependency(Node)
func (p *ProcessClosurePrimitives) ProcessAddDependency(n *rhino.Node) {
	// We can't modify parent, so just create a node that will
	// get compiled out.
	emptyNode := rhino.IR_Number(0)
	n.ReplaceWith(emptyNode)
	p.compiler.ReportChangeToEnclosingScope(emptyNode)
}

// ProcessForwardDeclare processes a goog.forwardDeclare() call and record
// the specified forward declaration.
//
// Java: ProcessClosurePrimitives.processForwardDeclare(Node)
func (p *ProcessClosurePrimitives) ProcessForwardDeclare(n *rhino.Node) {
	if !n.GetParent().IsExprResult() {
		//  Ignore "const Foo = goog.forwardDeclare('my.Foo');". It's legal but does not actually
		// forward declare the type 'my.Foo'.
		return
	}
	convention := p.compiler.GetCodingConvention()

	var typeDeclaration *string = nil
	jutil.Try(func() {
		typeDeclaration = processClosurePrimitives_getOnlyElement(
			convention.IdentifyTypeDeclarationCall(n))
	}, func(e jutil.Throwable) {
		switch e.(type) {
		case jutil.INullPointerException, jutil.INoSuchElementException, jutil.IIllegalArgumentException:
			p.compiler.Report(
				JSError_MakeWithNode(
					n,
					ProcessClosurePrimitives_INVALID_FORWARD_DECLARE,
					"A single type could not identified for the goog.forwardDeclare statement"))
		default:
			panic(e)
		}
	})

	if typeDeclaration != nil {
		p.compiler.ForwardDeclareType(*typeDeclaration)
	}
}

// processClosurePrimitives_getOnlyElement is Guava's
// Iterables.getOnlyElement(Iterable): it throws NullPointerException for a
// null list, NoSuchElementException for an empty one and
// IllegalArgumentException for one with several elements.
//
// Go-only: language-adapter — Guava Iterables.getOnlyElement; a nil slice
// stands for a null List.
func processClosurePrimitives_getOnlyElement(list []string) *string {
	if list == nil {
		panic(jutil.NewNullPointerException())
	}
	if len(list) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	if len(list) > 1 {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"expected one element but was: <" + processClosurePrimitives_joinElements(list) + ">"))
	}
	s := list[0]
	return &s
}

// processClosurePrimitives_joinElements renders the first (at most five)
// elements like Guava's getOnlyElement message.
//
// Go-only: language-adapter — Guava Iterators.getOnlyElement message.
func processClosurePrimitives_joinElements(list []string) string {
	s := ""
	for i := 0; i < 5 && i < len(list); i++ {
		if i > 0 {
			s += ", "
		}
		s += list[i]
	}
	if len(list) > 5 {
		s += ", ..."
	}
	return s
}

// ProcessClosurePrimitives_VerifySetCssNameMapping verifies that
// setCssNameMapping is called with the correct methods.
//
// Returns whether the arguments checked out okay.
//
// Java: ProcessClosurePrimitives.verifySetCssNameMapping(AbstractCompiler, Node, Node)
func ProcessClosurePrimitives_VerifySetCssNameMapping(
	compiler AbstractCompiler, methodName *rhino.Node, firstArg *rhino.Node) bool {
	var diagnostic *DiagnosticType = nil
	if firstArg == nil {
		diagnostic = ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR
	} else if !firstArg.IsObjectLit() {
		diagnostic = ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR
	} else if firstArg.GetNext() != nil {
		secondArg := firstArg.GetNext()
		if !secondArg.IsStringLit() {
			diagnostic = ProcessClosurePrimitives_EXPECTED_STRING_ERROR
		} else if secondArg.GetNext() != nil {
			diagnostic = ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR
		}
	}
	if diagnostic != nil {
		methodQname, ok := methodName.GetQualifiedNameString()
		if !ok {
			methodQname = "null"
		}
		compiler.Report(JSError_MakeWithNode(methodName, diagnostic, methodQname))
		return false
	}
	return true
}

// Java: ProcessClosurePrimitives.checkPropertyRenameCall(Node)
func (p *ProcessClosurePrimitives) CheckPropertyRenameCall(call *rhino.Node) {
	callee := call.GetFirstChild()
	// TODO(b/193038601): make this work when the callee is a module import
	if !p.compiler.GetCodingConvention().IsPropertyRenameFunction(callee) {
		return
	}
	calleeName, ok := callee.GetQualifiedNameString()
	if !ok {
		calleeName = "null"
	}

	switch call.GetChildCount() - 1 {
	case 1, 2:
	default:
		p.compiler.Report(
			JSError_MakeWithNode(
				call,
				ProcessClosurePrimitives_INVALID_RENAME_FUNCTION,
				calleeName,
				"Must be called with 1 or 2 arguments."))
	}

	propName := callee.GetNext()
	if propName == nil || !propName.IsStringLit() {
		p.compiler.Report(
			JSError_MakeWithNode(
				call,
				ProcessClosurePrimitives_INVALID_RENAME_FUNCTION,
				calleeName,
				"The first argument must be a string literal."))
	} else if jutil.StringIndexOf(propName.GetString(), ".") >= 0 {
		p.compiler.Report(
			JSError_MakeWithNode(
				call,
				ProcessClosurePrimitives_INVALID_RENAME_FUNCTION,
				calleeName,
				"The first argument must not be a property path."))
	}
}
