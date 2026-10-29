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
	"fmt"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
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
	compiler     AbstractCompiler
	contextStack []superContext
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

type superContextKind uint8

const (
	superCtxNotAllowed superContextKind = iota
	superCtxPropertyAccessAllowed
	superCtxMethod
	superCtxConstructor
	superCtxConstructorArrow
)

// superContext tracks lexical context by value on CheckSuper.contextStack.
type superContext struct {
	contextNode                          rhino.Node
	firstSuperCall                       rhino.Node
	thisAccessedBeforeSuper              rhino.Node
	superPropertyAccessedBeforeSuperCall rhino.Node
	constructorIdx                       int
	kind                                 superContextKind
	isClassMethod                        bool
	hasParentClass                       bool
	returnsAValue                        bool
}

func (ctx superContext) String() string {
	return fmt.Sprintf("superContext(%d, %v)", ctx.kind, ctx.contextNode)
}

// Java: CheckSuper(AbstractCompiler)
func NewCheckSuper(compiler AbstractCompiler) *CheckSuper {
	return &CheckSuper{compiler: compiler}
}

// Java: CheckSuper.process(Node, Node)
func (c *CheckSuper) Process(externs rhino.Node, root rhino.Node) {
	if c.processParallel(root) {
		return
	}
	NodeTraversal_Traverse(c.compiler, root, c)
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

func (c *CheckSuper) contextStackString() string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i := len(c.contextStack) - 1; i >= 0; i-- {
		if i < len(c.contextStack)-1 {
			sb.WriteString(", ")
		}
		sb.WriteString(c.contextStack[i].String())
	}
	sb.WriteByte(']')
	return sb.String()
}

func (c *CheckSuper) pushContext(ctx superContext) {
	rhino.CheckNotNull(ctx.contextNode)
	c.contextStack = append(c.contextStack, ctx)
}

func (c *CheckSuper) popContext() {
	last := len(c.contextStack) - 1
	if last < 0 {
		panic(jutil.NewNoSuchElementException())
	}
	c.contextStack[last] = superContext{}
	c.contextStack = c.contextStack[:last]
}

// Java: CheckSuper.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckSuper) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT:
		jutil.CheckState(len(c.contextStack) == 0)
		c.pushContext(superContext{kind: superCtxNotAllowed, contextNode: n})

	case rhino.Token_FUNCTION:
		c.pushFunctionContextIfNeeded(n)

	case rhino.Token_MEMBER_FIELD_DEF:
		c.pushContext(superContext{kind: superCtxPropertyAccessAllowed, contextNode: n})

	case rhino.Token_BLOCK: // For class static blocks
		if NodeUtil_IsClassStaticBlock(n) {
			c.pushContext(superContext{kind: superCtxPropertyAccessAllowed, contextNode: n})
		}

	case rhino.Token_CLASS: // TODO (user): For class fields
	default:
	}

	return true
}

func (c *CheckSuper) pushFunctionContextIfNeeded(fn rhino.Node) {
	if NodeUtil_IsMethodDeclaration(fn) {
		if NodeUtil_IsEs6Constructor(fn) {
			classNode := NodeUtil_GetEnclosingClass(fn)
			_, hasParentClass := ast.MustClass(classNode).Extends()
			c.pushContext(superContext{
				kind:           superCtxConstructor,
				contextNode:    fn,
				hasParentClass: hasParentClass,
			})
		} else {
			objLitOrClassMembers := rhino.CheckNotNull(fn.GetGrandparent())
			isClassMethod := !objLitOrClassMembers.IsObjectLit()
			if isClassMethod {
				jutil.CheckStateMsg1(objLitOrClassMembers.IsClassMembers(), "%s", objLitOrClassMembers)
			}
			c.pushContext(superContext{
				kind:          superCtxMethod,
				contextNode:   fn,
				isClassMethod: isClassMethod,
			})
		}
		return
	}

	// Arrow function context varies depending on the enclosing context, while
	// super is never allowed in normal non-method functions.
	if ast.MustFunction(fn).IsArrow() {
		topIdx := len(c.contextStack) - 1
		if topIdx < 0 {
			panic(jutil.NewNullPointerException())
		}
		// Only constructors create a distinct arrow context to distinguish return
		// statements and `super()` calls directly in the constructor from those in
		// arrow functions; all other contexts share their enclosing context.
		if c.contextStack[topIdx].kind == superCtxConstructor {
			c.pushContext(superContext{
				kind:           superCtxConstructorArrow,
				contextNode:    fn,
				constructorIdx: topIdx,
			})
		}
		return
	}

	c.pushContext(superContext{kind: superCtxNotAllowed, contextNode: fn})
}

var checkSuperVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	for _, t := range [...]rhino.Token{
		rhino.Token_SUPER,
		rhino.Token_CALL,
		rhino.Token_THIS,
		rhino.Token_RETURN,
		rhino.Token_ROOT,
		rhino.Token_FUNCTION,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_BLOCK,
	} {
		s.Add(t)
	}
	return s
}()

func (*CheckSuper) visitTokens() *ast.TokenSet        { return &checkSuperVisitTokens }
func (*CheckSuper) combinedVisitTokens() ast.TokenSet { return checkSuperVisitTokens }

// Java: CheckSuper.visit(NodeTraversal, Node, Node)
func (c *CheckSuper) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if len(c.contextStack) == 0 {
		panic(jutil.NewNullPointerException())
	}
	ctx := &c.contextStack[len(c.contextStack)-1]

	switch n.GetToken() {
	case rhino.Token_SUPER:
		if c.IsSuperConstructorCall(n) {
			// Defer recording the call until visiting the parent CALL node so any
			// this/super references in call arguments are treated as pre-super() references.
		} else if c.IsSuperPropertyAccess(n) {
			c.visitSuperPropertyAccess(ctx, t, n)
		} else {
			// super used some way other than `super()`, `super.prop`, or `super[expr]`.
			t.Report(n, CheckSuper_INVALID_SUPER_USAGE)
		}
	case rhino.Token_CALL:
		callee := ast.MustCall(n).Callee().Node()
		if callee.IsSuper() {
			c.visitSuperConstructorCall(ctx, t, callee)
		}
	case rhino.Token_THIS:
		if ctx.kind == superCtxConstructor && ctx.firstSuperCall.IsNil() {
			ctx.thisAccessedBeforeSuper = n
		}
	case rhino.Token_RETURN:
		if ctx.kind == superCtxConstructor {
			if _, hasVal := ast.MustReturn(n).Value(); hasVal {
				ctx.returnsAValue = true
			}
		}
	case rhino.Token_ROOT:
		jutil.CheckState(len(c.contextStack) == 1)
	default:
	}

	if n == ctx.contextNode {
		if ctx.kind == superCtxConstructor {
			c.finishConstructorContext(ctx, t)
		}
		c.popContext()
	}
}

func (c *CheckSuper) visitSuperConstructorCall(ctx *superContext, t *NodeTraversal, superNode rhino.Node) {
	switch ctx.kind {
	case superCtxNotAllowed, superCtxPropertyAccessAllowed:
		t.Report(superNode, CheckSuper_INVALID_SUPER_CALL)
	case superCtxMethod:
		if ctx.isClassMethod {
			if methodDef, ok := ast.AsMethodDef(rhino.CheckNotNull(ctx.contextNode.GetParent())); ok {
				t.Report(superNode, CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION, methodDef.Key())
				return
			}
		}
		t.Report(superNode, CheckSuper_INVALID_SUPER_CALL)
	case superCtxConstructor:
		if ctx.firstSuperCall.IsNil() {
			ctx.firstSuperCall = superNode
		}
	case superCtxConstructorArrow:
		t.Report(superNode, CheckSuper_SUPER_CALL_IN_ARROW)
		ctor := &c.contextStack[ctx.constructorIdx]
		if ctor.firstSuperCall.IsNil() {
			ctor.firstSuperCall = superNode
		}
	}
}

func (c *CheckSuper) visitSuperPropertyAccess(ctx *superContext, t *NodeTraversal, superNode rhino.Node) {
	switch ctx.kind {
	case superCtxNotAllowed:
		t.Report(superNode, CheckSuper_INVALID_SUPER_ACCESS)
	case superCtxPropertyAccessAllowed, superCtxMethod:
		// Allowed in methods, class fields, and static blocks.
	case superCtxConstructor:
		if ctx.firstSuperCall.IsNil() {
			ctx.superPropertyAccessedBeforeSuperCall = superNode
		}
	case superCtxConstructorArrow:
		ctor := &c.contextStack[ctx.constructorIdx]
		if ctor.firstSuperCall.IsNil() {
			ctor.superPropertyAccessedBeforeSuperCall = superNode
		}
	}
}

func (c *CheckSuper) finishConstructorContext(ctx *superContext, t *NodeTraversal) {
	if !ctx.hasParentClass {
		if !ctx.firstSuperCall.IsNil() {
			t.Report(ctx.firstSuperCall, CheckSuper_INVALID_SUPER_CALL)
		}
		return
	}
	if ctx.firstSuperCall.IsNil() && !ctx.returnsAValue {
		// ErrorToFixMapper expects the error on the parent MEMBER_FUNCTION_DEF.
		t.Report(ctx.contextNode.GetParent(), CheckSuper_MISSING_CALL_TO_SUPER)
	}
	if !ctx.thisAccessedBeforeSuper.IsNil() {
		t.Report(ctx.thisAccessedBeforeSuper, CheckSuper_THIS_BEFORE_SUPER)
	}
	if !ctx.superPropertyAccessedBeforeSuperCall.IsNil() {
		t.Report(ctx.superPropertyAccessedBeforeSuperCall, CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR)
	}
}

// Java: CheckSuper.isSuperConstructorCall(Node)
func (c *CheckSuper) IsSuperConstructorCall(superNode rhino.Node) bool {
	jutil.CheckStateMsg1(superNode.IsSuper(), "%s", superNode)
	call, ok := ast.AsCall(superNode.GetParent())
	return ok && call.Callee().Is(superNode)
}

// Java: CheckSuper.isSuperPropertyAccess(Node)
func (c *CheckSuper) IsSuperPropertyAccess(superNode rhino.Node) bool {
	jutil.CheckStateMsg1(superNode.IsSuper(), "%s", superNode)
	get, ok := ast.AsNormalGet(superNode.GetParent())
	return ok && get.Object().Is(superNode)
}
