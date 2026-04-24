// Port of com/google/javascript/jscomp/CheckSuper.java @ closure-compiler v20260415

// Check for errors related to the `super` keyword.
//
// NOTE: One might be tempted to make one or both of these refactoring
// improvements to this class.
//
//  1. Use an inner class as the Callback instead of making the whole class
//     implement it.
//  2. Use a ScopedCallback since we need to track function scopes here.
//
// Unfortunately, this class is used with CombinedCompilerPass, which will
// pass instances of it directly to NodeTraversal.traverse().

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)


// CheckSuper is a compiler pass that checks for errors related to the
// `super` keyword. It implements CompilerPass and acts as its own callback
// so it can be used with CombinedCompilerPass.
type CheckSuper struct {
	path         string
	diags        []diag.Diagnostic
	contextStack []superContext
}

// Go-only: constructor (Java uses new CheckSuper(compiler))
// NewCheckSuper creates a new CheckSuper pass for the given source path.
func NewCheckSuper(path string) *CheckSuper {
	return &CheckSuper{path: path}
}

// process mirrors Java's CompilerPass.process(Node externs, Node root) —
// entry point: runs NodeTraversal.traverse(compiler, root, this) and validates context stack.
func (c *CheckSuper) process(script *rhino.Node) {
	c.contextStack = nil
	c.traverseNode(script, nil)
	// Java also calls: checkState(this.contextStack.isEmpty(), ImmutableList.of(this.contextStack))
}

// createStaticBlockContext mirrors Java's CheckSuper.createStaticBlockContext.
// Creates a SuperPropertyAccessAllowedContext for a class static block.
func (c *CheckSuper) createStaticBlockContext(staticBlock *rhino.Node) superContext {
	return &superPropertyAccessAllowedContext{contextNode: staticBlock}
}

// Check walks the AST and returns diagnostics for super-keyword errors.
// Go-only: Go entry point that wraps process() and returns diagnostics.
func (c *CheckSuper) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.process(script)
	return c.diags
}

// report emits a diagnostic at node n with the given code and message.
func (c *CheckSuper) report(n *rhino.Node, code, message string) {
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    code,
		Level:   diag.Error,
		Message: message,
	})
}

// Go-only: inlined NodeTraversal.traverse logic (Java delegates to NodeTraversal)
// traverseNode is the main DFS walk — calls shouldTraverse and visit in the
// correct pre/post order, mirroring Java's NodeTraversal.traverse logic.
func (c *CheckSuper) traverseNode(n *rhino.Node, parent *rhino.Node) {
	if n == nil {
		return
	}
	if !c.shouldTraverse(n, parent) {
		return
	}
	for _, child := range n.Children() {
		c.traverseNode(child, n)
	}
	c.visit(n, parent)
}

// shouldTraverse mirrors Java's CheckSuper.shouldTraverse. Returns true always
// (we never skip subtrees), but manages the context stack on the way in.
func (c *CheckSuper) shouldTraverse(n *rhino.Node, parent *rhino.Node) bool {
	switch n.Token() {
	case rhino.SCRIPT:
		// Mirrors Java's ROOT case — push a SuperNotAllowedContext.
		if len(c.contextStack) == 0 {
			c.contextStack = append(c.contextStack, &superNotAllowedContext{contextNode: n})
		}

	case rhino.FUNCTION:
		currentCtx := c.peekContext()
		newCtx := c.getContextForFunctionNode(currentCtx, n)
		if newCtx != currentCtx {
			c.contextStack = append(c.contextStack, newCtx)
		}

	case rhino.CLASS_FIELD_DEF:
		// Mirrors Java's MEMBER_FIELD_DEF case.
		fieldCtx := c.createMemberFieldDefContext(n)
		c.contextStack = append(c.contextStack, fieldCtx)

	// Java also handles BLOCK for class static blocks via
	// NodeUtil.isClassStaticBlock. Our parser doesn't emit static blocks
	// yet, so we skip that case.
	// TODO (user): For class fields and static blocks when parser supports them.

	case rhino.CLASS:
		// TODO (user): For class fields
		break
	}

	return true
}

// visit mirrors Java's CheckSuper.visit — run post-order for each node.
func (c *CheckSuper) visit(n *rhino.Node, parent *rhino.Node) {
	currentCtx := c.peekContext()
	if currentCtx == nil {
		return
	}

	switch n.Token() {
	case rhino.SUPER:
		if c.isSuperConstructorCall(n) {
			// Note: defer recording the call in currentContext.visitSuperConstructorCall
			// until visiting the parent CALL node. This ensures any this/super
			// references in the call arguments are treated as invalid, pre-super()
			// call, references.
		} else if c.isSuperPropertyAccess(n) {
			currentCtx.visitSuperPropertyAccess(c, n)
		} else {
			// super used some way other than super(), super.prop, or super[expr].
			c.report(n, diag.JSC_INVALID_SUPER_USAGE,
				"'super' may only be used in a call or property access")
		}

	case rhino.CALL:
		callee := n.FirstChild()
		if callee != nil && callee.Token() == rhino.SUPER {
			currentCtx.visitSuperConstructorCall(c, callee)
		}

	case rhino.THIS:
		currentCtx.visitThis(c, n)

	case rhino.RETURN:
		currentCtx.visitReturn(c, n)

	case rhino.SCRIPT:
		// Nothing to do for ROOT/SCRIPT visit.
	}

	if n == currentCtx.getContextNode() {
		currentCtx.visitContextNode(c)
		c.contextStack = c.contextStack[:len(c.contextStack)-1]
	}
}

// Go-only: stack peek helper (Java uses ArrayDeque.peek())
// peekContext returns the top of the context stack, or nil if empty.
func (c *CheckSuper) peekContext() superContext {
	if len(c.contextStack) == 0 {
		return nil
	}
	return c.contextStack[len(c.contextStack)-1]
}

// getContextForFunctionNode mirrors Java's CheckSuper.getContextForFunctionNode.
func (c *CheckSuper) getContextForFunctionNode(currentCtx superContext, fn *rhino.Node) superContext {
	if isMethodDeclaration(fn) {
		if isEs6Constructor(fn) {
			return &constructorContext{contextNode: fn, hasParentClass: classHasParent(fn)}
		}
		return newMethodContext(fn)
	}
	// Arrow function context varies depending on the actual context, but
	// super is never allowed in normal functions.
	if fn.GetBooleanProp(rhino.PropArrowFn) {
		return currentCtx.getContextForArrowFunctionNode(fn)
	}
	return &superNotAllowedContext{contextNode: fn}
}

// isSuperConstructorCall reports whether superNode is the callee of a CALL.
func (c *CheckSuper) isSuperConstructorCall(superNode *rhino.Node) bool {
	p := superNode.Parent()
	if p == nil {
		return false
	}
	return p.Token() == rhino.CALL && isFirstChild(superNode, p)
}

// isSuperPropertyAccess reports whether superNode is the object of a GETPROP
// or GETELEM (i.e., super.prop or super[expr]).
func (c *CheckSuper) isSuperPropertyAccess(superNode *rhino.Node) bool {
	p := superNode.Parent()
	if p == nil {
		return false
	}
	// isNormalGet: parent is GETPROP or GETELEM and superNode is first child.
	return (p.Token() == rhino.GETPROP || p.Token() == rhino.GETELEM) &&
		isFirstChild(superNode, p)
}

// createMemberFieldDefContext mirrors Java's createMemberFieldDefContext.
func (c *CheckSuper) createMemberFieldDefContext(fieldNode *rhino.Node) superContext {
	return &superPropertyAccessAllowedContext{contextNode: fieldNode}
}

// ─── helpers ──────────────────────────────────────────────────────────────────

// isFirstChild reports whether n is the first child of parent.
func isFirstChild(n, parent *rhino.Node) bool {
	kids := parent.Children()
	return len(kids) > 0 && kids[0] == n
}

// isMethodDeclaration reports whether fn is a method inside a class body or
// object literal. Mirrors NodeUtil.isMethodDeclaration in Java.
// In our AST a class method is a FUNCTION whose parent is the class-body BLOCK
// whose parent is a CLASS. An object-literal method would be a FUNCTION directly
// inside an OBJECTLIT (our parser doesn't yet produce these, but we guard for it).
func isMethodDeclaration(fn *rhino.Node) bool {
	if fn.Token() != rhino.FUNCTION {
		return false
	}
	p := fn.Parent()
	if p == nil {
		return false
	}
	// Class body: parent is BLOCK, grandparent is CLASS.
	if p.Token() == rhino.BLOCK {
		gp := p.Parent()
		if gp != nil && gp.Token() == rhino.CLASS {
			return true
		}
	}
	// Object-literal method: parent is OBJECTLIT.
	if p.Token() == rhino.OBJECTLIT {
		return true
	}
	return false
}

// isEs6Constructor reports whether fn is the constructor method of an ES6 class.
// Mirrors NodeUtil.isEs6Constructor in Java.
func isEs6Constructor(fn *rhino.Node) bool {
	if !isMethodDeclaration(fn) {
		return false
	}
	// Our parser places the method name as children[0] of the FUNCTION.
	first := fn.FirstChild()
	if first == nil {
		return false
	}
	return first.String_() == "constructor"
}

// Go-only: helper (Java uses !classNode.getSecondChild().isEmpty() inline in ConstructorContext ctor)
// classHasParent returns true if the class enclosing fn has an extends clause.
// fn must be inside a class body BLOCK.
func classHasParent(fn *rhino.Node) bool {
	cls := getEnclosingClass(fn)
	if cls == nil {
		return false
	}
	kids := cls.Children()
	// CLASS children: [0]=name, [1]=superName (empty string if no extends).
	if len(kids) < 2 {
		return false
	}
	return kids[1].String_() != ""
}

// getEnclosingClass walks up fn's parent chain to find the enclosing CLASS node.
// Mirrors NodeUtil.getEnclosingClass in Java.
func getEnclosingClass(fn *rhino.Node) *rhino.Node {
	for cur := fn.Parent(); cur != nil; cur = cur.Parent() {
		if cur.Token() == rhino.CLASS {
			return cur
		}
	}
	return nil
}

// ─── context interface and implementations ────────────────────────────────────

// superContext tracks lexical context and determines correct traversal behaviour.
// Mirrors Java's abstract Context class.
type superContext interface {
	getContextNode() *rhino.Node
	getContextForArrowFunctionNode(arrowFn *rhino.Node) superContext
	visitSuperConstructorCall(c *CheckSuper, superNode *rhino.Node)
	visitSuperPropertyAccess(c *CheckSuper, superNode *rhino.Node)
	visitThis(c *CheckSuper, thisNode *rhino.Node)
	visitReturn(c *CheckSuper, returnNode *rhino.Node)
	visitContextNode(c *CheckSuper)
}

// ── superPropertyAccessAllowedContext ────────────────────────────────────────

// superPropertyAccessAllowedContext mirrors Java's SuperPropertyAccessAllowedContext.
// Used for public class fields and static blocks where super.prop is allowed but
// super() is not.
type superPropertyAccessAllowedContext struct {
	contextNode *rhino.Node
}

func (ctx *superPropertyAccessAllowedContext) getContextNode() *rhino.Node {
	return ctx.contextNode
}

func (ctx *superPropertyAccessAllowedContext) getContextForArrowFunctionNode(arrowFn *rhino.Node) superContext {
	return ctx
}

func (ctx *superPropertyAccessAllowedContext) visitSuperPropertyAccess(c *CheckSuper, superNode *rhino.Node) {
	// Super property should be allowed on a public field and inside static blocks.
	_ = superNode
}

func (ctx *superPropertyAccessAllowedContext) visitSuperConstructorCall(c *CheckSuper, superNode *rhino.Node) {
	// Super constructor calls are only allowed in constructor() methods.
	c.report(superNode, diag.JSC_INVALID_SUPER_CALL,
		"super() not allowed except in the constructor of a subclass")
}

func (ctx *superPropertyAccessAllowedContext) visitThis(c *CheckSuper, thisNode *rhino.Node) {
	// Ignored except in constructor methods. Mirrors Java Context base class default.
	_ = thisNode
}
func (ctx *superPropertyAccessAllowedContext) visitReturn(c *CheckSuper, returnNode *rhino.Node) {
	// Ignored except in constructor methods. Mirrors Java Context base class default.
	_ = returnNode
}
func (ctx *superPropertyAccessAllowedContext) visitContextNode(c *CheckSuper) {
	// Called when visiting the root node of this context — no-op here. Mirrors Java Context base class default.
	_ = c
}

// ── superNotAllowedContext ────────────────────────────────────────────────────

// superNotAllowedContext mirrors Java's SuperNotAllowedContext. Lexical context
// when not within a method, class, or object literal.
type superNotAllowedContext struct {
	contextNode *rhino.Node
}

func (ctx *superNotAllowedContext) getContextNode() *rhino.Node {
	return ctx.contextNode
}

func (ctx *superNotAllowedContext) getContextForArrowFunctionNode(arrowFn *rhino.Node) superContext {
	// Since super isn't allowed in this context, it isn't allowed in any of
	// its arrow functions either.
	return ctx
}

func (ctx *superNotAllowedContext) visitSuperConstructorCall(c *CheckSuper, superNode *rhino.Node) {
	// We're not within a constructor method.
	c.report(superNode, diag.JSC_INVALID_SUPER_CALL,
		"super() not allowed except in the constructor of a subclass")
}

func (ctx *superNotAllowedContext) visitSuperPropertyAccess(c *CheckSuper, superNode *rhino.Node) {
	// We're not within a method.
	c.report(superNode, diag.JSC_INVALID_SUPER_ACCESS,
		"'super' may only be accessed within a method")
}

func (ctx *superNotAllowedContext) visitThis(c *CheckSuper, thisNode *rhino.Node) {
	// Ignored except in constructor methods. Mirrors Java Context base class default.
	_ = thisNode
}
func (ctx *superNotAllowedContext) visitReturn(c *CheckSuper, returnNode *rhino.Node) {
	// Ignored except in constructor methods. Mirrors Java Context base class default.
	_ = returnNode
}
func (ctx *superNotAllowedContext) visitContextNode(c *CheckSuper) {
	// Called when visiting the root node of this context — no-op here. Mirrors Java Context base class default.
	_ = c
}

// ── methodContext ─────────────────────────────────────────────────────────────

// methodContext mirrors Java's MethodContext. Lexical context within a
// non-constructor method of either a class or an object literal.
type methodContext struct {
	contextNode  *rhino.Node
	isClassMethod bool
}

// Go-only: factory function (Java uses new MethodContext(fn) directly)
func newMethodContext(fn *rhino.Node) *methodContext {
	isClass := false
	if p := fn.Parent(); p != nil {
		if p.Token() == rhino.BLOCK {
			if gp := p.Parent(); gp != nil && gp.Token() == rhino.CLASS {
				isClass = true
			}
		}
	}
	return &methodContext{contextNode: fn, isClassMethod: isClass}
}

func (ctx *methodContext) getContextNode() *rhino.Node { return ctx.contextNode }

func (ctx *methodContext) getContextForArrowFunctionNode(arrowFn *rhino.Node) superContext {
	// In a method an arrow function keeps the same context as the method.
	return ctx
}

func (ctx *methodContext) visitSuperPropertyAccess(c *CheckSuper, superNode *rhino.Node) {
	// super property access in a method is perfectly valid, so do not report an error.
	_ = superNode
}

func (ctx *methodContext) visitSuperConstructorCall(c *CheckSuper, superNode *rhino.Node) {
	if ctx.isClassMethod {
		// super() — not allowed in non-constructor methods.
		propName := ctx.getPropertyName()
		if propName != "" {
			// Maybe the user was confused by other languages where super() invokes
			// the same method in the parent class?
			c.report(superNode, diag.JSC_INVALID_SUPER_CALL_WITH_SUGGESTION,
				fmt.Sprintf("super() not allowed here. Did you mean super.%s?", propName))
		} else {
			c.report(superNode, diag.JSC_INVALID_SUPER_CALL,
				"super() not allowed except in the constructor of a subclass")
		}
	} else {
		// object literal methods cannot contain super() calls
		c.report(superNode, diag.JSC_INVALID_SUPER_CALL,
			"super() not allowed except in the constructor of a subclass")
	}
}

func (ctx *methodContext) visitThis(c *CheckSuper, thisNode *rhino.Node) {
	// Ignored except in constructor methods. Mirrors Java Context base class default.
	_ = thisNode
}
func (ctx *methodContext) visitReturn(c *CheckSuper, returnNode *rhino.Node) {
	// Ignored except in constructor methods. Mirrors Java Context base class default.
	_ = returnNode
}
func (ctx *methodContext) visitContextNode(c *CheckSuper) {
	// Called when visiting the root node of this context — no-op here. Mirrors Java Context base class default.
	_ = c
}

// getPropertyName returns the method name for the diagnostic suggestion, or "".
// In Java, the MEMBER_FUNCTION_DEF / GETTER_DEF / SETTER_DEF parent node carries
// the name. In our AST the name is children[0] of the FUNCTION itself.
func (ctx *methodContext) getPropertyName() string {
	first := ctx.contextNode.FirstChild()
	if first == nil {
		return ""
	}
	return first.String_()
}

// ── constructorContext ────────────────────────────────────────────────────────

// constructorContext mirrors Java's ConstructorContext. Lexical context within
// a class constructor method.
type constructorContext struct {
	contextNode    *rhino.Node
	hasParentClass bool

	// Will be set to the first super() call that appears lexically, if any.
	firstSuperCall *rhino.Node
	// Call to super() isn't required if the constructor returns a value.
	returnsAValue bool
	thisAccessedBeforeSuper         *rhino.Node
	superPropertyAccessedBeforeSuperCall *rhino.Node
}

func (ctx *constructorContext) getContextNode() *rhino.Node { return ctx.contextNode }

func (ctx *constructorContext) getContextForArrowFunctionNode(arrowFn *rhino.Node) superContext {
	// Unlike normal methods, we need a separate context for arrow functions in
	// constructors. We need to distinguish between return statements and super()
	// calls that appear directly in the constructor vs. in arrow functions.
	return &constructorArrowContext{constructorCtx: ctx, contextNode: arrowFn}
}

func (ctx *constructorContext) visitSuperConstructorCall(c *CheckSuper, superNode *rhino.Node) {
	if ctx.firstSuperCall == nil {
		ctx.firstSuperCall = superNode
	}
}

func (ctx *constructorContext) visitSuperPropertyAccess(c *CheckSuper, superNode *rhino.Node) {
	if ctx.firstSuperCall == nil {
		ctx.superPropertyAccessedBeforeSuperCall = superNode
	}
}

func (ctx *constructorContext) visitThis(c *CheckSuper, thisNode *rhino.Node) {
	if ctx.firstSuperCall == nil {
		ctx.thisAccessedBeforeSuper = thisNode
	}
}

func (ctx *constructorContext) visitReturn(c *CheckSuper, returnNode *rhino.Node) {
	if returnNode.FirstChild() != nil {
		ctx.returnsAValue = true
	}
}

func (ctx *constructorContext) visitContextNode(c *CheckSuper) {
	// We've now visited the entire constructor, so we can decide whether to
	// report errors.
	if !ctx.hasParentClass {
		if ctx.firstSuperCall != nil {
			// calling super() only makes sense when there is a super class.
			c.report(ctx.firstSuperCall, diag.JSC_INVALID_SUPER_CALL,
				"super() not allowed except in the constructor of a subclass")
		}
	} else {
		// There is a parent class, so a call to super() is required unless the
		// constructor returns a value.
		if ctx.firstSuperCall == nil && !ctx.returnsAValue {
			// The context node is the function itself, but ErrorToFixMapper
			// expects the error to be reported on the MEMBER_FUNCTION_DEF that
			// is its parent. In our AST the FUNCTION is directly in the class
			// body BLOCK, so we report at the FUNCTION's parent-level name.
			// Java reports at getContextNode().getParent() — which is the
			// MEMBER_FUNCTION_DEF. We report at the FUNCTION node itself since
			// we have no MEMBER_FUNCTION_DEF wrapper.
			reportNode := ctx.contextNode
			c.report(reportNode, diag.JSC_MISSING_CALL_TO_SUPER,
				"constructor is missing a call to super()")
		}
		if ctx.thisAccessedBeforeSuper != nil {
			c.report(ctx.thisAccessedBeforeSuper, diag.JSC_THIS_BEFORE_SUPER,
				"cannot access this before calling super()")
		}
		if ctx.superPropertyAccessedBeforeSuperCall != nil {
			c.report(ctx.superPropertyAccessedBeforeSuperCall,
				diag.JSC_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR,
				"cannot access super properties before calling super()")
		}
	}
}

// ── constructorArrowContext ───────────────────────────────────────────────────

// constructorArrowContext mirrors Java's ConstructorArrowContext. Lexical context
// within an arrow function enclosed by a class constructor method.
type constructorArrowContext struct {
	constructorCtx *constructorContext
	contextNode    *rhino.Node
}

func (ctx *constructorArrowContext) getContextNode() *rhino.Node { return ctx.contextNode }

func (ctx *constructorArrowContext) getContextForArrowFunctionNode(arrowFn *rhino.Node) superContext {
	// Arrow functions within arrow functions share the same context behaviour.
	return ctx
}

func (ctx *constructorArrowContext) visitSuperConstructorCall(c *CheckSuper, superNode *rhino.Node) {
	c.report(superNode, diag.JSC_SUPER_CALL_IN_ARROW,
		"closure-compiler does not allow calls to `super()` in arrow functions")
	// Tell the constructor context about this super() call in order to avoid
	// confusing and likely redundant error messages such as "missing super
	// call" or "`this` accessed before super()".
	ctx.constructorCtx.visitSuperConstructorCall(c, superNode)
}

func (ctx *constructorArrowContext) visitSuperPropertyAccess(c *CheckSuper, superNode *rhino.Node) {
	// We will pretend the arrow function is immediately called, so it's as if
	// the super property reference appeared directly in the constructor.
	// The result is that you get an error for code like the following, even
	// though it isn't technically a JS error.
	//
	//   class Sub extends Base {
	//     constructor() {
	//       let arrow = () => super.prop; // ERROR: super.prop comes before super()
	//       super();
	//       arrow(); // not really executed until here, though
	//     }
	//   }
	//
	// This behaviour is consistent with TypeScript.
	// In general there's no good reason to declare such arrow functions before
	// calling super().
	ctx.constructorCtx.visitSuperPropertyAccess(c, superNode)
}

func (ctx *constructorArrowContext) visitThis(c *CheckSuper, thisNode *rhino.Node) {
	// Arrow functions in constructors don't independently track this-access. Mirrors Java ConstructorArrowContext.
	_ = thisNode
}
func (ctx *constructorArrowContext) visitReturn(c *CheckSuper, returnNode *rhino.Node) {
	// Return statements in arrow functions don't affect constructor return-a-value tracking. Mirrors Java.
	_ = returnNode
}
func (ctx *constructorArrowContext) visitContextNode(c *CheckSuper) {
	// Called when visiting the root node of this context — no-op here. Mirrors Java Context base class default.
	_ = c
}
