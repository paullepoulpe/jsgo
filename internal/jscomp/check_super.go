// Port of com/google/javascript/jscomp/CheckSuper.java @ closure-compiler v20261021

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckSuper checks for errors related to the `super` keyword.
//
// NOTE: One might be tempted to make one or both of these refactoring
// improvements to this class.
//
//  1. Use an inner class as the Callback instead of making the whole class
//     implement it.
//  2. Use a ScopedCallback since we need to track function scopes here.
//
// Unfortunately, this class is used with CombinedCompilerPass, which will pass
// instances of it directly to NodeTraversal.traverse().
type CheckSuper struct {
	compiler AbstractCompiler

	// Go-only: language-adapter — the Java ArrayDeque<Context> used as a
	// stack; the top of the stack is the last element.
	contextStack []CheckSuper_Context
}

// Java: CheckSuper.MISSING_CALL_TO_SUPER
var CheckSuper_MISSING_CALL_TO_SUPER = DiagnosticType_Error("JSC_MISSING_CALL_TO_SUPER", "constructor is missing a call to super()")

// For a derived class constructor that has a `super()` call or returns a value, but might still
// return without calling `super()`. Unlike MISSING_CALL_TO_SUPER, this can be suppressed, since
// the constructor might intentionally return a different object.
//
// Java: CheckSuper.MISSING_ROOT_LEVEL_CALL_TO_SUPER
var CheckSuper_MISSING_ROOT_LEVEL_CALL_TO_SUPER = DiagnosticType_Error(
	"JSC_MISSING_ROOT_LEVEL_CALL_TO_SUPER",
	"constructor must unconditionally call super() in a root-level statement before any"+
		" return")

// Java: CheckSuper.THIS_BEFORE_SUPER
var CheckSuper_THIS_BEFORE_SUPER = DiagnosticType_Error("JSC_THIS_BEFORE_SUPER", "cannot access this before calling super()")

// Java: CheckSuper.SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR
var CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR = DiagnosticType_Error(
	"JSC_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR",
	"cannot access super properties before calling super()")

// Java: CheckSuper.INVALID_SUPER_CALL
var CheckSuper_INVALID_SUPER_CALL = DiagnosticType_Error(
	"JSC_INVALID_SUPER_CALL", "super() not allowed except in the constructor of a subclass")

// The JS spec allows calls to `super()` in an arrow function within a constructor,
// as long as `super()` executes exactly once before any references to `this` or `super.prop`.
// However, doing that makes it very hard to statically determine whether `super()` is being
// called when it should be.
// There's really no good reason to call `super()` in an arrow function.
// It indicates that your code is overly complicated and you should refactor, so we will not
// allow it.
//
// Java: CheckSuper.SUPER_CALL_IN_ARROW
var CheckSuper_SUPER_CALL_IN_ARROW = DiagnosticType_Error(
	"JSC_SUPER_CALL_IN_ARROW",
	"closure-compiler does not allow calls to `super()` in arrow functions")

// Java: CheckSuper.INVALID_SUPER_USAGE
var CheckSuper_INVALID_SUPER_USAGE = DiagnosticType_Error(
	"JSC_INVALID_SUPER_USAGE", "''super'' may only be used in a call or property access")

// Java: CheckSuper.INVALID_SUPER_ACCESS
var CheckSuper_INVALID_SUPER_ACCESS = DiagnosticType_Error(
	"JSC_INVALID_SUPER_ACCESS", "''super'' may only be accessed within a method")

// Java: CheckSuper.INVALID_SUPER_CALL_WITH_SUGGESTION
var CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION = DiagnosticType_Error(
	"JSC_INVALID_SUPER_CALL_WITH_SUGGESTION",
	"super() not allowed here. Did you mean super.{0}?")

// Java: CheckSuper(AbstractCompiler)
func NewCheckSuper(compiler AbstractCompiler) *CheckSuper {
	c := &CheckSuper{}
	c.contextStack = nil
	c.compiler = compiler
	return c
}

// Java: CheckSuper.process(Node, Node)
func (c *CheckSuper) Process(externs rhino.Node, root rhino.Node) {
	if c.processParallel(root) {
		return
	}
	NodeTraversal_Traverse(c.compiler, root, c)
	// Java: checkState(this.contextStack.isEmpty(), ImmutableList.of(this.contextStack));
	// the message is only built on failure.
	if len(c.contextStack) != 0 {
		jutil.CheckStateMsg1(false, "%s", "["+c.contextStackString()+"]")
	}
}

func (c *CheckSuper) processParallel(root rhino.Node) bool {
	scripts, ok := collectRootScripts(root)
	if !ok {
		return false
	}
	numWorkers := parallelThreadCount(c.compiler, len(scripts))
	if numWorkers <= 1 {
		return false
	}
	prewarmCompilerForParallelReads(c.compiler)
	runParallelWork(c.compiler, numWorkers, len(scripts), func(_ int) (func(int, *recordingCompiler), func()) {
		var workerRC *recordingCompiler
		var workerPass *CheckSuper
		runItem := func(idx int, rc *recordingCompiler) {
			if workerPass == nil || workerRC != rc {
				workerRC = rc
				workerPass = NewCheckSuper(rc)
			}
			b := NodeTraversal_Builder{
				compiler: rc,
				callback: workerPass,
			}
			b.traverseScriptUnderRoot(root, root, scripts[idx], nil)
			if len(workerPass.contextStack) != 0 {
				jutil.CheckStateMsg1(false, "%s", "["+workerPass.contextStackString()+"]")
			}
		}
		return runItem, nil
	}, false)
	return true
}

// contextStackString renders the context stack the way ArrayDeque.toString()
// does (top of the stack first).
//
// Go-only: language-adapter — ArrayDeque.toString() for the stack slice.
func (c *CheckSuper) contextStackString() string {
	elems := slices.Clone(c.contextStack)
	slices.Reverse(elems)
	return jutil.CollectionToString(elems)
}

// Go-only: language-adapter — ArrayDeque.push(Context).
func (c *CheckSuper) pushContext(context CheckSuper_Context) {
	if context == nil {
		panic(jutil.NewNullPointerException())
	}
	c.contextStack = append(c.contextStack, context)
}

// Go-only: language-adapter — ArrayDeque.peek(): null when empty.
func (c *CheckSuper) peekContext() CheckSuper_Context {
	if len(c.contextStack) == 0 {
		return nil
	}
	return c.contextStack[len(c.contextStack)-1]
}

// Go-only: language-adapter — ArrayDeque.pop(): NoSuchElementException when
// empty.
func (c *CheckSuper) popContext() CheckSuper_Context {
	if len(c.contextStack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	top := c.contextStack[len(c.contextStack)-1]
	c.contextStack[len(c.contextStack)-1] = nil
	c.contextStack = c.contextStack[:len(c.contextStack)-1]
	return top
}

// Java: CheckSuper.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckSuper) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT:
		jutil.CheckState(len(c.contextStack) == 0)
		c.pushContext(NewCheckSuper_SuperNotAllowedContext(n))

	case rhino.Token_FUNCTION:
		{
			currentContext := c.peekContext()
			newContext := c.GetContextForFunctionNode(currentContext, n)
			if newContext != currentContext {
				c.pushContext(newContext)
			}
		}
	case rhino.Token_MEMBER_FIELD_DEF:
		fieldContext := c.CreateMemberFieldDefContext(n)
		c.pushContext(fieldContext)
	case rhino.Token_BLOCK: // For class static blocks
		if NodeUtil_IsClassStaticBlock(n) {
			newContext := c.CreateStaticBlockContext(n)
			c.pushContext(newContext)
		}

	case rhino.Token_CLASS: // TODO (user): For class fields
	default:
	}

	return true
}

// Java: CheckSuper.getContextForFunctionNode(Context, Node)
func (c *CheckSuper) GetContextForFunctionNode(currentContext CheckSuper_Context, fn rhino.Node) CheckSuper_Context {
	if NodeUtil_IsMethodDeclaration(fn) {
		if NodeUtil_IsEs6Constructor(fn) {
			return NewCheckSuper_ConstructorContext(fn)
		} else {
			return NewCheckSuper_MethodContext(fn)
		}
	} else {
		// Arrow function context varies depending on the actual context, but
		// super is never allowed in normal functions.
		if fn.IsArrowFunction() {
			return currentContext.GetContextForArrowFunctionNode(fn)
		}
		return NewCheckSuper_SuperNotAllowedContext(fn)
	}
}

// Java: CheckSuper.visit(NodeTraversal, Node, Node)
func (c *CheckSuper) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	currentContext := jutil.CheckNotNull(c.peekContext())

	switch n.GetToken() {
	case rhino.Token_SUPER:
		if c.IsSuperConstructorCall(n) {
			// Note: defer recording the call in `currentContext.visitSuperConstructorCall` until
			// visitng the parent CALL node. This ensures any this/super references in the
			// call arguments are treated as invalid, pre-super() call, references.
		} else if c.IsSuperPropertyAccess(n) {
			currentContext.VisitSuperPropertyAccess(t, n)
		} else {
			// super used some way other than `super()`, `super.prop`, or `super[expr]`.
			t.Report(n, CheckSuper_INVALID_SUPER_USAGE)
		}
	case rhino.Token_CALL:
		callee := n.GetFirstChild()
		if callee.IsSuper() {
			currentContext.VisitSuperConstructorCall(t, callee)
		}
	case rhino.Token_THIS:
		currentContext.VisitThis(t, n)
	case rhino.Token_RETURN:
		currentContext.VisitReturn(t, n)
	case rhino.Token_ROOT:
		jutil.CheckState(len(c.contextStack) == 1)
	default:
	}

	if n == currentContext.GetContextNode() {
		currentContext.VisitContextNode(t)
		c.popContext()
	}
}

// Java: CheckSuper.isSuperConstructorCall(Node)
func (c *CheckSuper) IsSuperConstructorCall(superNode rhino.Node) bool {
	jutil.CheckStateMsg1(superNode.IsSuper(), "%s", superNode)
	parent := superNode.GetParent()
	return parent.IsCall() && superNode.IsFirstChildOf(parent)
}

// Java: CheckSuper.isSuperPropertyAccess(Node)
func (c *CheckSuper) IsSuperPropertyAccess(superNode rhino.Node) bool {
	jutil.CheckStateMsg1(superNode.IsSuper(), "%s", superNode)
	parent := superNode.GetParent()
	return NodeUtil_IsNormalGet(parent) && superNode.IsFirstChildOf(parent)
}

// CheckSuper_Context tracks lexical context and determines correct traversal
// behavior based on it.
//
// Java: CheckSuper.Context
type CheckSuper_Context interface {
	// Java: CheckSuper.Context.getContextNode()
	GetContextNode() rhino.Node

	// GetContextForArrowFunctionNode: the correct context for an arrow
	// function depends on the enclosing context.
	//
	// Java: CheckSuper.Context.getContextForArrowFunctionNode(Node)
	GetContextForArrowFunctionNode(arrowFn rhino.Node) CheckSuper_Context

	// VisitSuperConstructorCall handles a super constructor call.
	//
	// Java: CheckSuper.Context.visitSuperConstructorCall(NodeTraversal, Node)
	VisitSuperConstructorCall(t *NodeTraversal, superNode rhino.Node)

	// VisitSuperPropertyAccess handles a super property reference.
	// (`super.prop` or `super[expr]`)
	//
	// Java: CheckSuper.Context.visitSuperPropertyAccess(NodeTraversal, Node)
	VisitSuperPropertyAccess(t *NodeTraversal, superNode rhino.Node)

	// Java: CheckSuper.Context.visitThis(NodeTraversal, Node)
	VisitThis(t *NodeTraversal, thisNode rhino.Node)

	// Java: CheckSuper.Context.visitReturn(NodeTraversal, Node)
	VisitReturn(t *NodeTraversal, returnNode rhino.Node)

	// Java: CheckSuper.Context.visitContextNode(NodeTraversal)
	VisitContextNode(t *NodeTraversal)

	GetCheckSuper_ContextBase() *CheckSuper_ContextBase
}

// CheckSuper_ContextBase holds the fields and base implementations of
// CheckSuper.Context.
type CheckSuper_ContextBase struct {
	contextNode rhino.Node
}

// Java: CheckSuper.Context(Node)
func (b *CheckSuper_ContextBase) InitCheckSuper_ContextBase(contextNode rhino.Node) {
	b.contextNode = rhino.CheckNotNull(contextNode)
}

func (b *CheckSuper_ContextBase) GetCheckSuper_ContextBase() *CheckSuper_ContextBase { return b }

// Java: CheckSuper.Context.getContextNode()
func (b *CheckSuper_ContextBase) GetContextNode() rhino.Node {
	return b.contextNode
}

// Java: CheckSuper.Context.visitThis(NodeTraversal, Node)
func (b *CheckSuper_ContextBase) VisitThis(t *NodeTraversal, thisNode rhino.Node) {
	// Ignored except in constructor methods.
}

// Java: CheckSuper.Context.visitReturn(NodeTraversal, Node)
func (b *CheckSuper_ContextBase) VisitReturn(t *NodeTraversal, returnNode rhino.Node) {
	// Ignored except in constructor methods.
}

// Java: CheckSuper.Context.visitContextNode(NodeTraversal)
func (b *CheckSuper_ContextBase) VisitContextNode(t *NodeTraversal) {
	// Called when visiting the root node of this context after all of its children have been
	// visited.
}

// Java: CheckSuper.SuperPropertyAccessAllowedContext
type CheckSuper_SuperPropertyAccessAllowedContext struct {
	CheckSuper_ContextBase
}

// Java: CheckSuper.SuperPropertyAccessAllowedContext(Node)
func NewCheckSuper_SuperPropertyAccessAllowedContext(contextNode rhino.Node) *CheckSuper_SuperPropertyAccessAllowedContext {
	c := &CheckSuper_SuperPropertyAccessAllowedContext{}
	c.InitCheckSuper_ContextBase(contextNode)
	return c
}

// Java: CheckSuper.SuperPropertyAccessAllowedContext.visitSuperPropertyAccess(NodeTraversal, Node)
func (c *CheckSuper_SuperPropertyAccessAllowedContext) VisitSuperPropertyAccess(t *NodeTraversal, superNode rhino.Node) {
	// Super property should be allowed on a public field and inside static blocks.
}

// Java: CheckSuper.SuperPropertyAccessAllowedContext.visitSuperConstructorCall(NodeTraversal, Node)
func (c *CheckSuper_SuperPropertyAccessAllowedContext) VisitSuperConstructorCall(t *NodeTraversal, superNode rhino.Node) {
	// Super constructor calls are only allowed in constructor() methods.
	t.Report(superNode, CheckSuper_INVALID_SUPER_CALL)
}

// Java: CheckSuper.SuperPropertyAccessAllowedContext.getContextForArrowFunctionNode(Node)
func (c *CheckSuper_SuperPropertyAccessAllowedContext) GetContextForArrowFunctionNode(arrowFn rhino.Node) CheckSuper_Context {
	return c
}

// Java: CheckSuper.createMemberFieldDefContext(Node)
func (c *CheckSuper) CreateMemberFieldDefContext(fieldNode rhino.Node) CheckSuper_Context {
	return NewCheckSuper_SuperPropertyAccessAllowedContext(fieldNode)
}

// Java: CheckSuper.createStaticBlockContext(Node)
func (c *CheckSuper) CreateStaticBlockContext(staticBlock rhino.Node) CheckSuper_Context {
	return NewCheckSuper_SuperPropertyAccessAllowedContext(staticBlock)
}

// CheckSuper_SuperNotAllowedContext is the lexical context when not within a
// method, class, or object literal.
//
// Java: CheckSuper.SuperNotAllowedContext
type CheckSuper_SuperNotAllowedContext struct {
	CheckSuper_ContextBase
}

// Java: CheckSuper.SuperNotAllowedContext(Node)
func NewCheckSuper_SuperNotAllowedContext(contextNode rhino.Node) *CheckSuper_SuperNotAllowedContext {
	c := &CheckSuper_SuperNotAllowedContext{}
	c.InitCheckSuper_ContextBase(contextNode)
	return c
}

// Java: CheckSuper.SuperNotAllowedContext.getContextForArrowFunctionNode(Node)
func (c *CheckSuper_SuperNotAllowedContext) GetContextForArrowFunctionNode(arrowFn rhino.Node) CheckSuper_Context {
	// Since super isn't allowed in this context, it isn't allowed in any of its arrow functions
	// either.
	return c
}

// Java: CheckSuper.SuperNotAllowedContext.visitSuperConstructorCall(NodeTraversal, Node)
func (c *CheckSuper_SuperNotAllowedContext) VisitSuperConstructorCall(t *NodeTraversal, superNode rhino.Node) {
	// We're not within a constructor method.
	t.Report(superNode, CheckSuper_INVALID_SUPER_CALL)
}

// Java: CheckSuper.SuperNotAllowedContext.visitSuperPropertyAccess(NodeTraversal, Node)
func (c *CheckSuper_SuperNotAllowedContext) VisitSuperPropertyAccess(t *NodeTraversal, superNode rhino.Node) {
	// We're not within a method.
	t.Report(superNode, CheckSuper_INVALID_SUPER_ACCESS)
}

// CheckSuper_MethodContext is the lexicial context within a non-constructor
// method of either a class or an object literal.
//
// Java: CheckSuper.MethodContext
type CheckSuper_MethodContext struct {
	CheckSuper_ContextBase

	isClassMethod bool
}

// Java: CheckSuper.MethodContext(Node)
func NewCheckSuper_MethodContext(functionNode rhino.Node) *CheckSuper_MethodContext {
	c := &CheckSuper_MethodContext{}
	c.InitCheckSuper_ContextBase(functionNode)
	jutil.CheckArgumentMsg1(
		functionNode.IsFunction() && NodeUtil_IsMethodDeclaration(functionNode), "%s", functionNode)
	objLitOrClassMembers := rhino.CheckNotNull(functionNode.GetGrandparent())
	if objLitOrClassMembers.IsObjectLit() {
		c.isClassMethod = false
	} else {
		jutil.CheckStateMsg1(objLitOrClassMembers.IsClassMembers(), "%s", objLitOrClassMembers)
		c.isClassMethod = true
	}
	return c
}

// Java: CheckSuper.MethodContext.getContextForArrowFunctionNode(Node)
func (c *CheckSuper_MethodContext) GetContextForArrowFunctionNode(arrowFn rhino.Node) CheckSuper_Context {
	// In a method an arrow function keeps the same context as the method, but a non-arrow
	// function creates a new non-method context.
	return c
}

// Java: CheckSuper.MethodContext.visitSuperPropertyAccess(NodeTraversal, Node)
func (c *CheckSuper_MethodContext) VisitSuperPropertyAccess(t *NodeTraversal, superNode rhino.Node) {
	// super property access in a method is perfectly valid, so do not report an error.
}

// Java: CheckSuper.MethodContext.visitSuperConstructorCall(NodeTraversal, Node)
func (c *CheckSuper_MethodContext) VisitSuperConstructorCall(t *NodeTraversal, superNode rhino.Node) {
	if c.isClassMethod {
		// super() - not allowed in non-constructor methods.
		propName := c.GetPropertyName()
		if propName != nil {
			// Maybe the user was confused by other languages where super() invokes the
			// same method in the parent class?
			t.Report(superNode, CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION, *propName)
		} else {
			t.Report(superNode, CheckSuper_INVALID_SUPER_CALL)
		}
	} else {
		// object literal methods cannot contain super() calls
		t.Report(superNode, CheckSuper_INVALID_SUPER_CALL)
	}
}

// GetPropertyName returns the property name associated with the method, if
// any, otherwise NULL.
//
// Java: CheckSuper.MethodContext.getPropertyName()
func (c *CheckSuper_MethodContext) GetPropertyName() *string {
	parent := rhino.CheckNotNull(c.GetContextNode().GetParent())
	// ```
	// class X {
	//   propertyName() {}
	//   get propertyName() {}
	//   set propertyName(value) {}
	//   [expression]() {} // no name for this one
	// }
	// ```
	if parent.IsMemberFunctionDef() || parent.IsGetterDef() || parent.IsSetterDef() {
		s := parent.GetString()
		return &s
	} else {
		return nil
	}
}

// CheckSuper_ConstructorContext is the lexical context within a class
// constructor method.
//
// Java: CheckSuper.ConstructorContext
type CheckSuper_ConstructorContext struct {
	CheckSuper_ContextBase

	hasParentClass bool

	// Will be set to the first `super()` call that appears lexically, if any.
	firstSuperCall rhino.Node
	// Call to super() isn't required if the constructor returns a value.
	returnsAValue                        bool
	thisAccessedBeforeSuper              rhino.Node
	superPropertyAccessedBeforeSuperCall rhino.Node
}

// Java: CheckSuper.ConstructorContext(Node)
func NewCheckSuper_ConstructorContext(contextNode rhino.Node) *CheckSuper_ConstructorContext {
	c := &CheckSuper_ConstructorContext{}
	c.InitCheckSuper_ContextBase(contextNode)
	c.firstSuperCall = rhino.NilNode
	c.returnsAValue = false
	c.thisAccessedBeforeSuper = rhino.NilNode
	c.superPropertyAccessedBeforeSuperCall = rhino.NilNode
	jutil.CheckArgumentMsg1(NodeUtil_IsEs6Constructor(contextNode), "%s", contextNode)
	classNode := NodeUtil_GetEnclosingClass(contextNode)
	c.hasParentClass = !classNode.GetSecondChild().IsEmpty()
	return c
}

// Java: CheckSuper.ConstructorContext.getContextForArrowFunctionNode(Node)
func (c *CheckSuper_ConstructorContext) GetContextForArrowFunctionNode(arrowFn rhino.Node) CheckSuper_Context {
	// Unlike normal methods, we need to create a separate context for arrow functions in
	// constructors. We need to distinguish between return statements and `super()` calls
	// that appear directly in the constructor and those that appear in arrow functions.
	return NewCheckSuper_ConstructorArrowContext(c, arrowFn)
}

// Java: CheckSuper.ConstructorContext.visitSuperConstructorCall(NodeTraversal, Node)
func (c *CheckSuper_ConstructorContext) VisitSuperConstructorCall(t *NodeTraversal, superNode rhino.Node) {
	if c.firstSuperCall.IsNil() {
		c.firstSuperCall = superNode
	}
}

// Java: CheckSuper.ConstructorContext.visitSuperPropertyAccess(NodeTraversal, Node)
func (c *CheckSuper_ConstructorContext) VisitSuperPropertyAccess(t *NodeTraversal, superNode rhino.Node) {
	if c.firstSuperCall.IsNil() {
		c.superPropertyAccessedBeforeSuperCall = superNode
	}
}

// Java: CheckSuper.ConstructorContext.visitThis(NodeTraversal, Node)
func (c *CheckSuper_ConstructorContext) VisitThis(t *NodeTraversal, thisNode rhino.Node) {
	if c.firstSuperCall.IsNil() {
		c.thisAccessedBeforeSuper = thisNode
	}
}

// Java: CheckSuper.ConstructorContext.visitReturn(NodeTraversal, Node)
func (c *CheckSuper_ConstructorContext) VisitReturn(t *NodeTraversal, returnNode rhino.Node) {
	if returnNode.HasChildren() {
		c.returnsAValue = true
	}
}

// Java: CheckSuper.ConstructorContext.visitContextNode(NodeTraversal)
func (c *CheckSuper_ConstructorContext) VisitContextNode(t *NodeTraversal) {
	// We've now visited the entire constructor, so we can decide whether to report errors.
	if !c.hasParentClass {
		if !c.firstSuperCall.IsNil() {
			// calling `super()` only makes sense when there is a super class.
			t.Report(c.firstSuperCall, CheckSuper_INVALID_SUPER_CALL)
		}
	} else {
		// There is a parent class, so a call to `super()` is required unless the constructor
		// returns a value.
		if c.firstSuperCall.IsNil() && !c.returnsAValue {
			// The context node is the function itself , but ErrorToFixMapper expects the error to
			// be reported on the MEMBER_FUNCTION_DEF that is its parent.
			t.Report(c.GetContextNode().GetParent(), CheckSuper_MISSING_CALL_TO_SUPER)
		}
		if !c.thisAccessedBeforeSuper.IsNil() {
			t.Report(c.thisAccessedBeforeSuper, CheckSuper_THIS_BEFORE_SUPER)
		}
		if !c.superPropertyAccessedBeforeSuperCall.IsNil() {
			t.Report(c.superPropertyAccessedBeforeSuperCall, CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR)
		}
	}
}

// CheckSuper_ConstructorArrowContext is the lexical context within an arrow
// function enclosed by a class constructor method.
//
// Java: CheckSuper.ConstructorArrowContext
type CheckSuper_ConstructorArrowContext struct {
	CheckSuper_ContextBase

	constructorContext *CheckSuper_ConstructorContext
}

// Java: CheckSuper.ConstructorArrowContext(ConstructorContext, Node)
func NewCheckSuper_ConstructorArrowContext(
	constructorContext *CheckSuper_ConstructorContext, arrowFn rhino.Node) *CheckSuper_ConstructorArrowContext {
	c := &CheckSuper_ConstructorArrowContext{}
	c.InitCheckSuper_ContextBase(arrowFn)
	c.constructorContext = jutil.CheckNotNull(constructorContext)
	return c
}

// Java: CheckSuper.ConstructorArrowContext.getContextForArrowFunctionNode(Node)
func (c *CheckSuper_ConstructorArrowContext) GetContextForArrowFunctionNode(arrowFn rhino.Node) CheckSuper_Context {
	// Arrow functions within arrow functions share the same context behavior.
	return c
}

// Java: CheckSuper.ConstructorArrowContext.visitSuperConstructorCall(NodeTraversal, Node)
func (c *CheckSuper_ConstructorArrowContext) VisitSuperConstructorCall(t *NodeTraversal, superNode rhino.Node) {
	t.Report(superNode, CheckSuper_SUPER_CALL_IN_ARROW)
	// Tell the constructor context about this `super()` call in order to avoid confusing and
	// likely redundant error messages such as "missing super call" or "`this` accessed before
	// super()".
	c.constructorContext.VisitSuperConstructorCall(t, superNode)
}

// Java: CheckSuper.ConstructorArrowContext.visitSuperPropertyAccess(NodeTraversal, Node)
func (c *CheckSuper_ConstructorArrowContext) VisitSuperPropertyAccess(t *NodeTraversal, superNode rhino.Node) {
	// We will pretend the arrow function is immediately called, so it's as if the super
	// property reference appeared directly in the constructor.
	// The result is that you get an error for code like the following, even though it isn't
	// technically a JS error.
	// ```
	// class Sub extends Base {
	//   constructor() {
	//     let arrow = () => super.prop; // ERROR: super.prop comes before super()
	//     super();
	//     arrow(); // not really executed until here, though
	//   }
	// }
	// ```
	// This behavior is consistent with TypeScript.
	// In general there's no good reason to declare such arrow functions before calling `super()`.
	c.constructorContext.VisitSuperPropertyAccess(t, superNode)
}
