// Port of com/google/javascript/jscomp/ProcessClosurePrimitives.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast). The checks
// take the ast.Call they look at and read it by role: the callee (a GetProp
// whose object is the name goog), the method name (its property), the
// arguments (FirstArg, Next, NumArgs) instead of GetNext chains from the
// callee and child counts that include it. goog.base calls read `this` and
// the method name as the first two arguments. Kept on rhino: the checks of
// function and class declarations and of the enclosing function or class
// (they read parents and siblings of any kind), findGoogInheritsCall (a
// walk over the following statements), and the nodes built for the
// rewrites.

// ProcessClosurePrimitives performs some Closure-specific simplifications
// including rewriting goog.base, goog.addDependency.
//
// Records forwardDeclared names in the internal compiler state.
type ProcessClosurePrimitives struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler

	knownClosureSubclasses map[string]struct{}

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
	p.compiler = compiler
	p.closureModules =
		jutil.CheckNotNullMsg(compiler.GetModuleMetadataMap(), "Need to run GatherModuleMetadata").
			GetModulesByGoogNamespace()
	return p
}

// Java: ProcessClosurePrimitives.process(Node, Node)
func (p *ProcessClosurePrimitives) Process(externs rhino.Node, root rhino.Node) {
	// Replace and validate other Closure primitives
	NodeTraversal_TraverseRoots(p.compiler, p, externs, root)
}

// Java: ProcessClosurePrimitives.visit(NodeTraversal, Node, Node)
func (p *ProcessClosurePrimitives) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		p.CheckGoogFunctions(t, call)
		p.MaybeProcessClassBaseCall(call)
		p.CheckPropertyRenameCall(call)
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
		assign, ok := ast.MustExprStmt(n).Expr().AsAssign()
		if !ok || n.IsFromExterns() {
			break
		}
		info := assign.Node().GetJSDocInfo()
		if info == nil || info.GetType() == nil {
			break
		}
		lhs, ok := assign.Target().Node().GetQualifiedNameString()
		if !ok {
			break
		}
		p.CheckPossibleGoogProvideInit(lhs, info, n)
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		first, _ := ast.MustNameDeclaration(n).FirstDeclarator()
		name, ok := first.AsName() // false for a destructuring declaration
		if !ok {
			break
		}
		p.CheckPossibleGoogProvideInit(name.Ident(), n.GetJSDocInfo(), n)
	default:
	}
}

// Java: ProcessClosurePrimitives.checkPossibleGoogProvideInit(String, JSDocInfo, Node)
func (p *ProcessClosurePrimitives) CheckPossibleGoogProvideInit(
	namespace string, info *rhino.JSDocInfo, definition rhino.Node) {
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
func (p *ProcessClosurePrimitives) CheckGoogFunctions(t *NodeTraversal, call ast.Call) {
	callee, ok := call.Callee().AsGetProp()
	if !ok {
		return
	}

	receiver, ok := callee.Object().AsName()
	if !ok || receiver.Ident() != ProcessClosurePrimitives_GOOG {
		return
	}

	// For the sake of simplicity, we report code changes
	// when we see a provides/requires, and don't worry about
	// reporting the change when we actually do the replacement.
	methodName := callee.Property()
	switch methodName {
	case "inherits":
		// Note: inherits is allowed in local scope
		p.ProcessInheritsCall(call)
	case "addDependency":
		if p.ValidateUnaliasablePrimitiveCall(t, call.Node(), methodName) {
			p.ProcessAddDependency(call)
		}
	case "setCssNameMapping":
		_ = ProcessClosurePrimitives_ProcessSetCssNameMapping(p.compiler, call.Node(), call.Node().GetParent())
	case "forwardDeclare":
		if p.ValidatePrimitiveCallWithMessage(
			t,
			call.Node(),
			methodName,
			ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR) {
			p.ProcessForwardDeclare(call.Node())
		}
	case "weakUsage":
		p.ValidateWeakUsageCall(call)
	default:
	}
}

// Java: ProcessClosurePrimitives.validateWeakUsageCall(Node)
func (p *ProcessClosurePrimitives) ValidateWeakUsageCall(call ast.Call) {
	// goog.weakUsage() should have exactly one argument, and it should be a name (possibly
	// qualified).
	calleeName := processClosurePrimitives_qualifiedNameOrNull(call.Callee().Node())
	if numArgs := call.NumArgs(); numArgs != 1 {
		p.compiler.Report(
			JSError_MakeWithNode(
				call.Node(),
				ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL,
				calleeName,
				"should have exactly one argument, not "+strconv.Itoa(numArgs)))
	}
	if arg, ok := call.FirstArg(); ok && !arg.Node().IsQualifiedName() {
		p.compiler.Report(
			JSError_MakeWithNode(
				call.Node(),
				ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL,
				calleeName,
				"argument should be a name or qualified name, not "+arg.Node().String()))
	}
}

// processClosurePrimitives_qualifiedNameOrNull returns the qualified name
// of n, or "null" (what Java prints for a null string) if it has none.
//
// Go-only: the getQualifiedName() results this pass puts in messages.
func processClosurePrimitives_qualifiedNameOrNull(n rhino.Node) string {
	if name, ok := n.GetQualifiedNameString(); ok {
		return name
	}
	return "null"
}

// ValidateUnaliasablePrimitiveCall verifies that a) the call is in the top
// level of a file and b) the return value is unused
//
// This method is for primitives that never return a value.
//
// Java: ProcessClosurePrimitives.validateUnaliasablePrimitiveCall(NodeTraversal, Node, String)
func (p *ProcessClosurePrimitives) ValidateUnaliasablePrimitiveCall(t *NodeTraversal, n rhino.Node, methodName string) bool {
	return p.ValidatePrimitiveCallWithMessage(t, n, methodName, ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR)
}

// ValidatePrimitiveCallWithMessage: methodName is the list of primitive
// types classed together with this one; invalidAliasingError is which
// DiagnosticType to emit if this call is aliased. this depends on whether
// the primitive is sometimes aliasiable in a module or never aliasable.
//
// Java: ProcessClosurePrimitives.validatePrimitiveCallWithMessage(NodeTraversal, Node, String, DiagnosticType)
func (p *ProcessClosurePrimitives) ValidatePrimitiveCallWithMessage(
	t *NodeTraversal, n rhino.Node, methodName string, invalidAliasingError *DiagnosticType) bool {
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
func (p *ProcessClosurePrimitives) MaybeProcessClassBaseCall(call ast.Call) {
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

	n := call.Node()
	callTarget, ok := call.Callee().AsGetProp()
	if !ok || callTarget.Property() != "base" {
		return
	}

	baseContainerNode := callTarget.Object().Node()
	if !baseContainerNode.IsUnscopedQualifiedName() {
		// Some unknown "base" method.
		return
	}
	baseContainer, _ := baseContainerNode.GetQualifiedNameString()

	enclosingFnNameNode := ProcessClosurePrimitives_GetEnclosingDeclNameNode(n)
	if enclosingFnNameNode.IsNil() || !enclosingFnNameNode.IsUnscopedQualifiedName() {
		// some unknown container method or a MEMBER_FUNCTION_DEF.
		if _, ok := p.knownClosureSubclasses[baseContainer]; ok {
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
		p.RewriteBaseCallInConstructor(enclosingQname, baseContainer, call, enclosingFnNameNode)
	} else {
		p.RewriteBaseCallInMethod(enclosingQname, baseContainer, call, enclosingFnNameNode)
	}
}

// Java: ProcessClosurePrimitives.rewriteBaseCallInConstructor(String, String, Node, Node)
func (p *ProcessClosurePrimitives) RewriteBaseCallInConstructor(
	enclosingQname string, baseContainer string, call ast.Call, enclosingFnNameNode rhino.Node) {
	n := call.Node()
	// Check if this is some other "base" method.
	if enclosingQname != baseContainer {
		// Report misuse of "base" methods from other known classes.
		if _, ok := p.knownClosureSubclasses[baseContainer]; ok {
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
	if baseClassNode.IsNil() {
		// If there is no "goog.inherits", this might be some other "base" method.
		return
	}

	// This is the expected method, validate its parameters.
	thisArg, ok := call.FirstArg()
	if !ok || thisArg.Token() != rhino.Token_THIS {
		p.ReportBadBaseMethodUse(n, baseContainer, "First argument must be 'this'.")
		return
	}

	// Handle methods.
	methodName, ok := processClosurePrimitives_stringArgAfter(thisArg)
	if !ok || methodName.Value() != "constructor" {
		p.ReportBadBaseMethodUse(n, baseContainer, "Second argument must be 'constructor'.")
		return
	}

	// We're good to go.
	baseClassQname := processClosurePrimitives_qualifiedNameOrNull(baseClassNode)
	callee := call.Callee().Node()
	callee.ReplaceWith(
		NodeUtil_NewQNameWithOriginalName(
			p.compiler,
			baseClassQname+".call",
			callee,
			enclosingQname+".base"))
	methodName.Detach()
	p.compiler.ReportChangeToEnclosingScope(n)
}

// processClosurePrimitives_stringArgAfter returns the argument after arg;
// ok is false if there is none or it is not a string literal.
//
// Go-only: the `methodNameNode == null || !methodNameNode.isStringLit()`
// checks of the goog.base rewrites.
func processClosurePrimitives_stringArgAfter(arg ast.Arg) (ast.StringLit, bool) {
	next, ok := arg.Next()
	if !ok {
		return ast.StringLit{}, false
	}
	return next.AsStringLit()
}

// RewriteBaseCallInMethod: enclosingQname is a Java String that is never
// null here (a Go string is never null).
//
// Java: ProcessClosurePrimitives.rewriteBaseCallInMethod(String, String, Node, Node)
func (p *ProcessClosurePrimitives) RewriteBaseCallInMethod(
	enclosingQname string, baseContainer string, call ast.Call, enclosingFnNameNode rhino.Node) {
	n := call.Node()
	if _, ok := p.knownClosureSubclasses[baseContainer]; !ok {
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
	thisArg, ok := call.FirstArg()
	if !ok || thisArg.Token() != rhino.Token_THIS {
		p.ReportBadBaseMethodUse(n, baseContainer, "First argument must be 'this'.")
		return
	}

	// Handle methods.
	methodNameNode, ok := processClosurePrimitives_stringArgAfter(thisArg)
	if !ok {
		p.ReportBadBaseMethodUse(n, baseContainer, "Second argument must name a method.")
		return
	}

	methodName := methodNameNode.Value()
	ending := ".prototype." + methodName
	if !jutil.StringEndsWith(enclosingQname, ending) {
		p.ReportBadBaseMethodUse(n, baseContainer, "Enclosing method does not match "+methodName)
		return
	}

	// We're good to go.
	className := enclosingFnNameNode.GetFirstFirstChild()
	classQname := processClosurePrimitives_qualifiedNameOrNull(className)
	callee := call.Callee().Node()
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
func ProcessClosurePrimitives_FindGoogInheritsCall(ctorDeclarationNode rhino.Node) rhino.Node {
	maybeInheritsExpr := ctorDeclarationNode.GetNext()

	for !maybeInheritsExpr.IsNil() &&
		(maybeInheritsExpr.IsEmpty() || NodeUtil_IsExprAssign(maybeInheritsExpr)) {
		// Skip empty nodes (e.g. those created by an unnecessary semicolon) and potential aliases,
		// e.g. "exports = Ctor;" in a goog.module.
		maybeInheritsExpr = maybeInheritsExpr.GetNext()
	}
	var baseClassNode rhino.Node = rhino.NilNode
	if !maybeInheritsExpr.IsNil() && NodeUtil_IsExprCall(maybeInheritsExpr) {
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
func (p *ProcessClosurePrimitives) ProcessInheritsCall(call ast.Call) {
	if call.NumArgs() == 2 {
		// The count proves both arguments are there.
		subClass, _ := call.FirstArg()
		superClass, _ := subClass.Next()
		if subClass.Node().IsUnscopedQualifiedName() && superClass.Node().IsUnscopedQualifiedName() {
			if subQname, ok := subClass.Node().GetQualifiedNameString(); ok {
				if p.knownClosureSubclasses == nil {
					p.knownClosureSubclasses = make(map[string]struct{})
				}
				p.knownClosureSubclasses[subQname] = struct{}{}
			}
		}
	}
}

// ProcessClosurePrimitives_GetEnclosingDeclNameNode returns the qualified
// name node of the function whose scope we're in, or null if it cannot be
// found.
//
// Java: ProcessClosurePrimitives.getEnclosingDeclNameNode(Node)
func ProcessClosurePrimitives_GetEnclosingDeclNameNode(n rhino.Node) rhino.Node {
	fn := NodeUtil_GetEnclosingFunction(n)
	if fn.IsNil() {
		return rhino.NilNode
	}
	return NodeUtil_GetNameNode(fn)
}

// ProcessClosurePrimitives_BaseUsedInClass verifies if goog.base call is
// used in a class
//
// Java: ProcessClosurePrimitives.baseUsedInClass(Node)
func ProcessClosurePrimitives_BaseUsedInClass(n rhino.Node) bool {
	for curr := n; !curr.IsNil(); curr = curr.GetParent() {
		if curr.IsClassMembers() {
			return true
		}
	}
	return false
}

// ReportBadBaseMethodUse reports an incorrect use of super-method calling.
//
// Java: ProcessClosurePrimitives.reportBadBaseMethodUse(Node, String, String)
func (p *ProcessClosurePrimitives) ReportBadBaseMethodUse(n rhino.Node, className string, extraMessage string) {
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
func (p *ProcessClosurePrimitives) ReportPossibleBadBaseMethodUse(n rhino.Node, className string, extraMessage string) {
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
	compiler AbstractCompiler, n rhino.Node, parent rhino.Node) CssRenamingMap {
	call := ast.MustCall(n)
	if !ProcessClosurePrimitives_VerifySetCssNameMapping(compiler, call) {
		return nil
	}
	// VerifySetCssNameMapping checked the arguments: an object literal, then
	// an optional string.
	arg, _ := call.FirstArg()
	obj := ast.MustObjectLit(arg.Node())

	// Translate OBJECTLIT into SubstitutionMap. All keys and
	// values must be strings, or an error will be thrown.
	cssNames := jutil.NewLinkedHashMap[string, string]()

	for member := range obj.Members() {
		key, isKey := member.AsStringKey()
		var value ast.StringLit
		isString := false
		if isKey {
			value, isString = key.Value().AsStringLit()
		}
		if !isString {
			compiler.Report(JSError_MakeWithNode(n, ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR))
			return nil
		}
		cssNames.Put(key.Key(), value.Value())
	}

	styleStr := "BY_PART"
	if style, ok := arg.Next(); ok {
		styleStr = ast.MustStringLit(style.Node()).Value()
	}

	var style CssRenamingMap_Style
	switch styleStr {
	case "BY_WHOLE":
		style = CssRenamingMap_Style_BY_WHOLE
	case "BY_PART":
		style = CssRenamingMap_Style_BY_PART
	default:
		compiler.Report(JSError_MakeWithNode(n, ProcessClosurePrimitives_INVALID_STYLE_ERROR, styleStr))
		return nil
	}

	if style == CssRenamingMap_Style_BY_PART {
		// Make sure that no class keys contain -'s
		// Variable keys can contain -'s since they're always named BY_WHOLE.
		var errors []string
		for key := range cssNames.Keys() {
			if jutil.StringIndexOf(key, "-") >= 0 && !jutil.StringStartsWith(key, "--") {
				errors = append(errors, key)
			}
		}
		if len(errors) > 0 {
			compiler.Report(JSError_MakeWithNode(n, ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP, jutil.CollectionToString(errors)))
		}
	} else if style == CssRenamingMap_Style_BY_WHOLE {
		// Verifying things is a lot trickier here. We just do a quick
		// n^2 check over the map which makes sure that if "a-b" in
		// the map, then map(a-b) = map(a)-map(b).
		// To speed things up, only consider cases where len(b) <= 10
		var errors []string
		for bKey, bValue := range cssNames.All() {
			if jutil.StringLength(bKey) > 10 {
				continue
			}
			for aKey, aValue := range cssNames.All() {
				combined, ok := cssNames.Lookup(aKey + "-" + bKey)
				if ok && combined != aValue+"-"+bValue {
					errors = append(errors,
						"map("+
							aKey+
							"-"+
							bKey+
							") != map("+
							aKey+
							")-map("+
							bKey+
							")")
				}
			}
		}
		if len(errors) > 0 {
			compiler.Report(JSError_MakeWithNode(n, ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP, jutil.CollectionToString(errors)))
		}
	}

	return &processClosurePrimitives_cssRenamingMap{cssNames: cssNames, style: style}
}

// ProcessAddDependency processes a goog.addDependency() call and record any
// forward declarations.
//
// Java: ProcessClosurePrimitives.processAddDependency(Node)
func (p *ProcessClosurePrimitives) ProcessAddDependency(call ast.Call) {
	// We can't modify parent, so just create a node that will
	// get compiled out.
	emptyNode := rhino.IR_Number(0)
	call.Node().ReplaceWith(emptyNode)
	p.compiler.ReportChangeToEnclosingScope(emptyNode)
}

// ProcessForwardDeclare processes a goog.forwardDeclare() call and record
// the specified forward declaration.
//
// Java: ProcessClosurePrimitives.processForwardDeclare(Node)
func (p *ProcessClosurePrimitives) ProcessForwardDeclare(n rhino.Node) {
	if !n.GetParent().IsExprResult() {
		//  Ignore "const Foo = goog.forwardDeclare('my.Foo');". It's legal but does not actually
		// forward declare the type 'my.Foo'.
		return
	}
	convention := p.compiler.GetCodingConvention()

	typeDeclarations := convention.IdentifyTypeDeclarationCall(n)
	if len(typeDeclarations) != 1 {
		p.compiler.Report(
			JSError_MakeWithNode(
				n,
				ProcessClosurePrimitives_INVALID_FORWARD_DECLARE,
				"A single type could not identified for the goog.forwardDeclare statement"))
		return
	}

	p.compiler.ForwardDeclareType(typeDeclarations[0])
}

// ProcessClosurePrimitives_VerifySetCssNameMapping verifies that
// setCssNameMapping is called with the correct methods.
//
// Returns whether the arguments checked out okay.
//
// Go-only: it takes the call; Java takes its callee and first argument.
//
// Java: ProcessClosurePrimitives.verifySetCssNameMapping(AbstractCompiler, Node, Node)
func ProcessClosurePrimitives_VerifySetCssNameMapping(compiler AbstractCompiler, call ast.Call) bool {
	var diagnostic *DiagnosticType = nil
	if firstArg, ok := call.FirstArg(); !ok {
		diagnostic = ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR
	} else if firstArg.Token() != rhino.Token_OBJECTLIT {
		diagnostic = ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR
	} else if secondArg, ok := firstArg.Next(); ok {
		if secondArg.Token() != rhino.Token_STRINGLIT {
			diagnostic = ProcessClosurePrimitives_EXPECTED_STRING_ERROR
		} else if !secondArg.IsLast() {
			diagnostic = ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR
		}
	}
	if diagnostic != nil {
		methodName := call.Callee().Node()
		compiler.Report(JSError_MakeWithNode(methodName, diagnostic, processClosurePrimitives_qualifiedNameOrNull(methodName)))
		return false
	}
	return true
}

// Java: ProcessClosurePrimitives.checkPropertyRenameCall(Node)
func (p *ProcessClosurePrimitives) CheckPropertyRenameCall(call ast.Call) {
	callee := call.Callee().Node()
	// TODO(b/193038601): make this work when the callee is a module import
	if !p.compiler.GetCodingConvention().IsPropertyRenameFunction(callee) {
		return
	}
	calleeName := processClosurePrimitives_qualifiedNameOrNull(callee)

	switch call.NumArgs() {
	case 1, 2:
	default:
		p.compiler.Report(
			JSError_MakeWithNode(
				call.Node(),
				ProcessClosurePrimitives_INVALID_RENAME_FUNCTION,
				calleeName,
				"Must be called with 1 or 2 arguments."))
	}

	var propName ast.StringLit
	isString := false
	if arg, ok := call.FirstArg(); ok {
		propName, isString = arg.AsStringLit()
	}
	if !isString {
		p.compiler.Report(
			JSError_MakeWithNode(
				call.Node(),
				ProcessClosurePrimitives_INVALID_RENAME_FUNCTION,
				calleeName,
				"The first argument must be a string literal."))
	} else if jutil.StringIndexOf(propName.Value(), ".") >= 0 {
		p.compiler.Report(
			JSError_MakeWithNode(
				call.Node(),
				ProcessClosurePrimitives_INVALID_RENAME_FUNCTION,
				calleeName,
				"The first argument must not be a property path."))
	}
}
