// Port of com/google/javascript/jscomp/GuardedCallback.java @ closure-compiler v20261015

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

// GuardedCallback is an AST traverser that keeps track of whether access to
// a generic resource are "guarded" or not. A guarded use is one that will
// not cause runtime errors if the resource does not exist. The resource
// (type T in the signature) is any value type computable from a node, such
// as a qualified name or property name. It is up to the subclass to map the
// currently visited token to the T in order to call IsGuarded, which
// depends entirely on the context from ancestor nodes and previous calls to
// IsGuarded for the same resource.
//
// More precisely, a resource may be guarded either intrinsically or
// conditionally, as follows. A use is intrinsically guarded if it occurs in
// a context where its specific value is immediately discarded, such as
// coercion to boolean or string. A use is conditionally guarded if it
// occurs in a context conditioned on an intrinsically guarded use (such as
// the "then" or "else" block of an "if", the second or third argument of a
// ternary operator, right-hand arguments of logical operators, or later in
// a block in which an unconditional "throw" or "return" was found in a
// guarded context).
//
// For example, the following are all intrinsically guarded uses of x:
//
//	// Coerced to boolean:
//	if (x);
//	x ? y : z;
//	x && y;
//	Boolean(x);
//	!x;
//	x == y; x != y; x === y; x !== y;
//	x instanceof Foo;
//	typeof x === 'string';
//
//	// Coerced to string:
//	String(x);
//	typeof x;
//
//	// Immediately discarded (but doesn't make much sense in my contexts):
//	x = y;
//
// The following uses of x are conditionally guarded:
//
//	if (x) x();
//	!x ? null : x;
//	typeof x == 'function' ? x : () => {};
//	x && x.y;
//	!x || x.y;
//	if (!x) return; x();
//	x ?? x = 3;
//
// Note that there is no logic to determine which branch is guarded, so any
// usages in either branch will pass after such a check. As such, the
// following are also considered guarded, though a human can easily see
// that this is spurious:
//
//	if (!x) x();
//	if (x) { } else x();
//	!x && x();
//	var y = x != null ? null : x;
//
// Note also that the call or property access is not necessary to make a use
// unguarded: the final example immediately above would be unguarded if it
// weren't for the x != null condition, since it allows the value of x to
// leak out in an uncontrolled way.
//
// This class overrides the Callback API methods with final methods of its
// own, and defines the template method VisitGuarded to perform the normal
// work for individual nodes. The only other API is IsGuarded, which allows
// checking if a T in the current node's context is guarded, either
// intrinsically or conditionally. If it is intrinsically guarded, then it
// may be recorded as a condition for the purpose of guarding future
// contexts.
type GuardedCallback[T comparable] interface {
	NodeTraversal_Callback

	// VisitGuarded performs specific traversal behavior. Should call
	// IsGuarded at least once.
	//
	// Java: GuardedCallback.visitGuarded(NodeTraversal, Node, Node)
	VisitGuarded(traversal *NodeTraversal, n *rhino.Node, parent *rhino.Node)

	// Java: GuardedCallback.isGuarded(T)
	IsGuarded(resource T) bool

	GetGuardedCallbackBase() *GuardedCallbackBase[T]
}

// GuardedCallbackBase holds GuardedCallback's fields and implementation.
type GuardedCallbackBase[T comparable] struct {
	self GuardedCallback[T]

	// Compiler is needed for coding convention (isPropertyTestFunction).
	compiler         AbstractCompiler
	codingConvention CodingConvention
	// Map from short-circuiting conditional nodes (AND, OR, COALESCE, IF,
	// and HOOK) to the set of resources each node guards. This is saved
	// separately from just `guarded` because the guard must not go into
	// effect until after traversal of the first child is complete. Before
	// traversing the second child any node, its values in this map are
	// moved into `guarded` and `installedGuards` (the latter allowing
	// removal at the correct time).
	//
	// Go-only: language-adapter — MultimapBuilder.hashKeys().hashSetValues():
	// Java iterates a key's values in String hash order, here in insertion
	// order. The order is not observable: the values only flow into
	// `guarded` (a multiset queried for membership) and `installedGuards`
	// (whose lists only flow back into this multimap and into
	// guarded.removeAll), and nothing iterates the keys.
	registeredGuards map[*rhino.Node][]T
	// Set of currently-guarded resources. Elements are added to this set
	// just before traversing the second or later (i.e. "then" or "else")
	// child of a short-circuiting conditional node, and then removed after
	// traversing the last child.
	guarded map[T]struct{}
	// Resources that are currently installed as guarded but will need to be
	// removed from `guarded` after visiting all the key nodes' children.
	installedGuards map[*rhino.Node][]T
	// A stack of `Context` objects describing the current node's context:
	// specifically, whether it is inherently safe, and a link to one or more
	// conditional nodes in the current statement directly above it (for
	// registering safe resources as guards).
	//
	// Go-only: language-adapter — the ArrayDeque used as a stack: the top
	// is the last element of the slice; emptyDepth counts leading EMPTY frames
	// at the bottom of the stack when contextStack is empty.
	contextStack    []*guardedCallback_Context
	emptyDepth      int
	onlyNameLeaves  bool
	polyfills       *PolyfillUsageFinder_Polyfills
	contextStackBuf [128]*guardedCallback_Context
	contextNodes    [128]guardedCallback_Context
}

// InitGuardedCallbackBase is the GuardedCallback(AbstractCompiler)
// constructor; self is the callback being constructed.
//
// Java: GuardedCallback(AbstractCompiler)
func (b *GuardedCallbackBase[T]) InitGuardedCallbackBase(self GuardedCallback[T], compiler AbstractCompiler) {
	b.self = self
	b.registeredGuards = nil
	b.guarded = nil
	b.installedGuards = nil
	b.compiler = compiler
	if compiler != nil {
		b.codingConvention = compiler.GetCodingConvention()
	}
	b.contextStack = b.contextStackBuf[:0]
	b.emptyDepth = 0
}

func (b *GuardedCallbackBase[T]) putRegisteredGuard(node *rhino.Node, resource T) {
	list := b.registeredGuards[node]
	if slices.Contains(list, resource) {
		return
	}
	if b.registeredGuards == nil {
		b.registeredGuards = make(map[*rhino.Node][]T)
	}
	b.registeredGuards[node] = append(list, resource)
}

// GetGuardedCallbackBase gives access to the base through the interface.
func (b *GuardedCallbackBase[T]) GetGuardedCallbackBase() *GuardedCallbackBase[T] { return b }

// Go-only: language-adapter — ArrayDeque.peek(): null when empty.
func (b *GuardedCallbackBase[T]) peekContext() *guardedCallback_Context {
	if len(b.contextStack) > 0 {
		return b.contextStack[len(b.contextStack)-1]
	}
	if b.emptyDepth > 0 {
		return guardedCallback_Context_EMPTY
	}
	return nil
}

func (b *GuardedCallbackBase[T]) canPossiblyBeGuarded() bool {
	if len(b.contextStack) > 0 && b.contextStack[len(b.contextStack)-1].safe {
		return true
	}
	return len(b.guarded) > 0
}

// Go-only: language-adapter — ArrayDeque.push(E) (null-hostile).
func (b *GuardedCallbackBase[T]) pushContext(c *guardedCallback_Context) {
	jutil.CheckNotNull(c)
	if len(b.contextStack) == 0 && c == guardedCallback_Context_EMPTY {
		b.emptyDepth++
		return
	}
	b.contextStack = append(b.contextStack, c)
}

// Go-only: language-adapter — ArrayDeque.pop(): throws
// NoSuchElementException when empty.
func (b *GuardedCallbackBase[T]) popContext() *guardedCallback_Context {
	if len(b.contextStack) > 0 {
		top := b.contextStack[len(b.contextStack)-1]
		b.contextStack = b.contextStack[:len(b.contextStack)-1]
		return top
	}
	if b.emptyDepth > 0 {
		b.emptyDepth--
		return guardedCallback_Context_EMPTY
	}
	panic(jutil.NewNoSuchElementException())
}

// Java: GuardedCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *GuardedCallbackBase[T]) ShouldTraverse(traversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	// Note that shouldTraverse() operates primarily on `parent`, while
	// visit() uses `n`. This is intentional. To see why, consider traversing
	// the following tree:
	//
	//    if (x) y; else z;
	//
	// 1. shouldTraverse(`if`):
	//    a. parent is null, so pushes an EMPTY onto the context stack.
	// 2. shouldTraverse(`x`):
	//    a. parent is `if`, so pushes Context(`if`, true); guards is empty.
	// 3. visit(`x`)
	//    a. guarded and installedGuards are both empty, so nothing is
	//       removed.
	//    b. visitGuarded(`x`) will call isGuarded("x"), which looks at the
	//       top of the stack and sees that the context is safe (true) and
	//       that there is a linked conditional node (the `if`); adds
	//       {`if`: "x"} to registeredGuards.
	//    b. Context(`if`, true) is popped off the stack.
	// 4. shouldTraverse(`y`):
	//    a. parent is still `if`, but since `y` is the second child it is no
	//       longer safe, so another EMPTY is pushed.
	//    b. the {`if`: "x"} guard is moved from registered to installed.
	// 5. visit(`y`):
	//    a. nothing is installed on `y` so no guards are removed.
	//    b. visitGuarded(`y`) will call isGuarded("y"), which will return
	//       false since "y" is neither intrinsically or conditionally
	//       guarded; if we'd called isResourceRequired("x"), it would return
	//       false because "x" is currently an element of guarded.
	//    c. one empty context is popped.
	// 6. shouldTraverse(`z`), visit(`z`)
	//    a. see steps 4-5, nothing really changes here.
	// 7. visit(`if`)
	//    a. the installed {`if`: "x"} guard is removed.
	//    c. pop the final empty context from the stack.

	if parent == nil {
		// The first node gets an empty context.
		if len(b.contextStack) == 0 {
			b.emptyDepth++
		} else {
			b.contextStack = append(b.contextStack, guardedCallback_Context_EMPTY)
		}
		return true
	}
	if b.onlyNameLeaves && n.GetFirstChild() == nil && len(b.registeredGuards) == 0 {
		nt := n.GetToken()
		if nt != rhino.Token_NAME {
			if len(b.installedGuards) == 0 || (nt != rhino.Token_RETURN && nt != rhino.Token_THROW) {
				return false
			}
		} else if b.polyfills != nil && !b.polyfills.mightHaveStaticNode(n) {
			return false
		}
	}

	// Before traversing any children, we update the stack
	ptok := parent.GetToken()
	if len(b.contextStack) == 0 {
		if uint(ptok) < 256 && guardedCallback_tokenFlags[ptok]&guardedCallback_flagEmptyDescend == 0 {
			b.emptyDepth++
		} else {
			next := guardedCallback_Context_EMPTY.descendIntoWithToken(b.compiler, b.codingConvention, parent, ptok, n, &b.contextNodes[0])
			if next == guardedCallback_Context_EMPTY {
				b.emptyDepth++
			} else {
				b.contextStack = append(b.contextStack, next)
			}
		}
	} else {
		var slot *guardedCallback_Context
		if idx := len(b.contextStack); idx < len(b.contextNodes) {
			slot = &b.contextNodes[idx]
		}
		top := b.contextStack[len(b.contextStack)-1]
		b.contextStack = append(b.contextStack, top.descendIntoWithToken(b.compiler, b.codingConvention, parent, ptok, n, slot))
	}

	// If the parent has any guards registered on it, then add them to
	// both `guarded` and `installedGuards`.
	if len(b.registeredGuards) > 0 &&
		uint(ptok) < 256 && guardedCallback_tokenFlags[ptok]&guardedCallback_flagCanHaveGuards != 0 {
		if resources, ok := b.registeredGuards[parent]; ok {
			delete(b.registeredGuards, parent)
			if b.guarded == nil {
				b.guarded = make(map[T]struct{})
			}
			if b.installedGuards == nil {
				b.installedGuards = make(map[*rhino.Node][]T)
			}
			for _, resource := range resources {
				b.guarded[resource] = struct{}{}
			}
			b.installedGuards[parent] = append(b.installedGuards[parent], resources...)
		}
	}
	return true
}

// Java: GuardedCallback.visit(NodeTraversal, Node, Node)
func (b *GuardedCallbackBase[T]) Visit(traversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// Remove any guards registered on this node by its children, which are
	// no longer relevant. This happens first because these were registered
	// on a "parent", but now this is that parent (i.e. `n` here vs `parent`
	// in isGuarded).
	ntok := n.GetToken()
	if len(b.installedGuards) > 0 {
		if parent != nil && uint(ntok) < 256 && guardedCallback_tokenFlags[ntok]&guardedCallback_flagCanHaveGuards != 0 {
			if resources, ok := b.installedGuards[n]; ok {
				delete(b.installedGuards, n)
				for _, r := range resources {
					delete(b.guarded, r)
				}
			}
		}

		// Check for abrupt returns (`return` and `throw`).
		if ntok == rhino.Token_RETURN || ntok == rhino.Token_THROW {
			// If found, any guards installed on a parent IF should be promoted to
			// the grandparent. This allows a small amount of flow-sensitivity, in
			// that
			//   if (!x) return; x();
			// has the guard for `x` promoted from the `if` to the outer block, so
			// that it guards the next statement.
			b.PromoteAbruptReturns(parent)
		}
	}

	// Finally continue on to whatever the traversal would normally do.
	if b.polyfills == nil {
		b.self.VisitGuarded(traversal, n, parent)
	} else if ntok == rhino.Token_GETPROP || ntok == rhino.Token_OPTCHAIN_GETPROP {
		if b.polyfills.mightHavePropNode(n) {
			b.self.VisitGuarded(traversal, n, parent)
		}
	} else if ntok == rhino.Token_NAME {
		if b.polyfills.mightHaveStaticNode(n) {
			b.self.VisitGuarded(traversal, n, parent)
		}
	}

	// After children have been traversed, pop the top of the conditional
	// stack.
	if len(b.contextStack) > 0 {
		b.contextStack = b.contextStack[:len(b.contextStack)-1]
	} else if b.emptyDepth > 0 {
		b.emptyDepth--
	} else {
		panic(jutil.NewNoSuchElementException())
	}
}

// Java: GuardedCallback.promoteAbruptReturns(Node)
func (b *GuardedCallbackBase[T]) PromoteAbruptReturns(parent *rhino.Node) {
	// If the parent is a BLOCK (e.g. `if (x) { return; }`) then go up one
	// level.
	if parent.IsBlock() {
		parent = parent.GetParent()
	}
	// If there were any guards registered the parent IF, then promote them
	// up one level.
	if parent.IsIf() {
		if resources := b.installedGuards[parent]; len(resources) > 0 {
			grandparent := parent.GetParent()
			if grandparent.IsBlock() || grandparent.IsScript() {
				for _, r := range resources {
					b.putRegisteredGuard(grandparent, r)
				}
			}
		}
	}
}

// IsGuarded determines if the given resource is guarded, either
// intrinsically or conditionally. If the former, any ancestor conditional
// nodes are registered as feature-testing the resource.
//
// Java: GuardedCallback.isGuarded(T)
func (b *GuardedCallbackBase[T]) IsGuarded(resource T) bool {
	// Check if this polyfill is already guarded. If so, return true right
	// away.
	if len(b.guarded) > 0 {
		if _, ok := b.guarded[resource]; ok {
			return true
		}
	}

	// If not, see if this is itself a feature check guard. This is defined
	// as a usage of the polyfill in such a way that throws away the actual
	// value and only cares about its truthiness or typeof. We walk up the
	// ancestor tree through a small set of node types and if this is
	// detected to be a guard, then the conditional node is marked as a guard
	// for this polyfill.
	if len(b.contextStack) == 0 {
		return false
	}
	context := b.contextStack[len(b.contextStack)-1]
	if !context.safe {
		return false
	}

	// Loop over all the linked conditionals and register this as a guard.
	for context != nil && context.conditional != nil {
		b.putRegisteredGuard(context.conditional, resource)
		context = context.linked
	}
	return true
}

// guardedCallback_Context is the context of a node, keeping track of
// whether it is safe for possibly-undefined values, and whether there are
// any conditionals upstream in the tree.
type guardedCallback_Context struct {
	// The most recent conditional.
	conditional *rhino.Node
	// Whether this position is safe for an undefined type.
	safe bool
	// A very naive linked list for storing additional conditional nodes.
	linked *guardedCallback_Context
}

// guardedCallback_Context_EMPTY: an empty instance: unsafe and with no
// linked conditional nodes.
//
// Java: GuardedCallback.Context.EMPTY
var guardedCallback_Context_EMPTY = newGuardedCallback_Context(nil, false, nil)

// Java: GuardedCallback.Context(Node, boolean, Context)
func newGuardedCallback_Context(conditional *rhino.Node, safe bool, linked *guardedCallback_Context) *guardedCallback_Context {
	return &guardedCallback_Context{conditional: conditional, safe: safe, linked: linked}
}

// Link returns a new Context with a new conditional node and safety
// status. If the current context already has a conditional, then it is
// linked so that both can be marked when necessary.
//
// Java: GuardedCallback.Context.link(Node, boolean)
func (c *guardedCallback_Context) Link(newConditional *rhino.Node, newSafe bool) *guardedCallback_Context {
	return c.linkInto(newConditional, newSafe, nil)
}

func (c *guardedCallback_Context) linkInto(newConditional *rhino.Node, newSafe bool, slot *guardedCallback_Context) *guardedCallback_Context {
	var linked *guardedCallback_Context
	if c.conditional != nil {
		linked = c
	}
	if slot != nil {
		*slot = guardedCallback_Context{
			conditional: newConditional,
			safe:        newSafe,
			linked:      linked,
		}
		return slot
	}
	return newGuardedCallback_Context(newConditional, newSafe, linked)
}

var guardedCallback_Context_SAFE_EMPTY = newGuardedCallback_Context(nil, true, nil)

// Propagate returns a new Context with a different safety bit, but doesn't
// change anything else.
//
// Java: GuardedCallback.Context.propagate(boolean)
func (c *guardedCallback_Context) Propagate(newSafe bool) *guardedCallback_Context {
	return c.propagateInto(newSafe, nil)
}

func (c *guardedCallback_Context) propagateInto(newSafe bool, slot *guardedCallback_Context) *guardedCallback_Context {
	if newSafe == c.safe {
		return c
	}
	if c.conditional == nil && c.linked == nil {
		if newSafe {
			return guardedCallback_Context_SAFE_EMPTY
		}
		return guardedCallback_Context_EMPTY
	}
	if slot != nil {
		*slot = guardedCallback_Context{
			conditional: c.conditional,
			safe:        newSafe,
			linked:      c.linked,
		}
		return slot
	}
	return newGuardedCallback_Context(c.conditional, newSafe, c.linked)
}

// Descend returns a new context given the current context and the next
// parent node. Child is only used to determine whether we're looking at the
// first child or not.
//
// Java: GuardedCallback.Context.descend(AbstractCompiler, Node, Node)
func (c *guardedCallback_Context) Descend(compiler AbstractCompiler, parent *rhino.Node, child *rhino.Node) *guardedCallback_Context {
	return c.descendInto(compiler, parent, child, nil)
}

func (c *guardedCallback_Context) descendInto(compiler AbstractCompiler, parent *rhino.Node, child *rhino.Node, slot *guardedCallback_Context) *guardedCallback_Context {
	var cc CodingConvention
	if compiler != nil {
		cc = compiler.GetCodingConvention()
	}
	return c.descendIntoWithToken(compiler, cc, parent, parent.GetToken(), child, slot)
}

func (c *guardedCallback_Context) descendIntoWithToken(compiler AbstractCompiler, cc CodingConvention, parent *rhino.Node, ptok rhino.Token, child *rhino.Node, slot *guardedCallback_Context) *guardedCallback_Context {
	switch ptok {
	case rhino.Token_CAST:
		// Casts are irrelevant.
		return c
	case rhino.Token_COMMA:
		// `Promise, whatever` is safe.
		// `whatever, Promise` is same as outer context.
		if child == parent.GetLastChild() {
			return c
		}
		return c.propagateInto(true, slot)
	case rhino.Token_AND, rhino.Token_HOOK:
		// `Promise && whatever` never returns Promise itself, so it is safe.
		// `whatever && Promise` may return Promise, so return outer context.
		// `Promise ? whatever : whatever` is a safe conditional.
		// `whatever ? Promise : whatever` (etc) is same as outer context.
		if child == parent.GetFirstChild() {
			return c.linkInto(parent, true, slot)
		}
		return c
	case rhino.Token_OR, rhino.Token_COALESCE:
		// `Promise || whatever` and `Promise ?? whatever`
		// may return Promise (unsafe), but is itself a conditional.
		// `whatever || Promise` and `whatever ?? Promise`
		// is same as outer context.
		if child == parent.GetFirstChild() {
			return c.linkInto(parent, false, slot)
		}
		return c
	case rhino.Token_IF:
		// `if (Promise) whatever` is a safe conditional.
		// `if (whatever) { ... }` is nothing.
		// TODO(sdh): Handle do/while/for/for-of/for-in?
		if child == parent.GetFirstChild() {
			return c.linkInto(parent, true, slot)
		}
		return guardedCallback_Context_EMPTY
	case rhino.Token_INSTANCEOF, rhino.Token_ASSIGN:
		// `Promise instanceof whatever` is safe, `whatever instanceof Promise`
		// is not. `Promise = whatever` is a bad idea, but it's safe w.r.t.
		// polyfills.
		return c.propagateInto(child == parent.GetFirstChild(), slot)
	case rhino.Token_TYPEOF, rhino.Token_NOT, rhino.Token_EQ, rhino.Token_NE, rhino.Token_SHEQ, rhino.Token_SHNE:
		// `typeof Promise` is always safe, as is `Promise == whatever`, etc.
		return c.propagateInto(true, slot)
	case rhino.Token_CALL:
		// `String(Promise)` is safe, `Promise(whatever)` or
		// `whatever(Promise)` is not.
		return c.propagateInto(child != parent.GetFirstChild() && guardedCallback_isPropertyTestFunctionWithCC(compiler, cc, parent), slot)
	case rhino.Token_ROOT:
		// This case causes problems for isStatement() so handle it
		// separately.
		return guardedCallback_Context_EMPTY
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_GETPROP:
		if child == parent.GetFirstChild() {
			// thisNode?.rest.of.chain
			// OR firstChild?.thisNode.rest.of.chain
			// For the first case `thisNode` should be considered
			// intrinsically guarded.
			return c.linkInto(parent, parent.IsOptionalChainStart(), slot)
		} else {
			// `first?.(thisNode)`
			// or `first?.[thisNode]`
			// or `first?.thisNode`
			return c.propagateInto(false, slot)
		}
	default:
		// Expressions propagate linked conditionals; statements do not.
		if c.conditional == nil || NodeUtil_IsStatement(parent) {
			return guardedCallback_Context_EMPTY
		}
		return c.propagateInto(false, slot)
	}
}

// Java: GuardedCallback.isAbrupt(Node)
func guardedCallback_isAbrupt(n *rhino.Node) bool {
	return n.IsReturn() || n.IsThrow()
}

// guardedCallback_isPropertyTestFunction extends the coding convention's
// idea of property test functions to also include String() and Boolean().
//
// Java: GuardedCallback.isPropertyTestFunction(AbstractCompiler, Node)
func guardedCallback_isPropertyTestFunction(compiler AbstractCompiler, n *rhino.Node) bool {
	var cc CodingConvention
	if compiler != nil {
		cc = compiler.GetCodingConvention()
	}
	return guardedCallback_isPropertyTestFunctionWithCC(compiler, cc, n)
}

func guardedCallback_isPropertyTestFunctionWithCC(compiler AbstractCompiler, cc CodingConvention, n *rhino.Node) bool {
	if cc == nil && compiler != nil {
		cc = compiler.GetCodingConvention()
	}
	if cc != nil && cc.IsPropertyTestFunction(n) {
		return true
	}
	target := n.GetFirstChild()
	if target.IsName() {
		name := target.GetString()
		return name == "String" || name == "Boolean"
	}
	return false
}

const (
	guardedCallback_flagEmptyDescend  uint8 = 1 << 0
	guardedCallback_flagCanHaveGuards uint8 = 1 << 1
)

var guardedCallback_tokenFlags = func() [256]uint8 {
	var flags [256]uint8
	for _, t := range []rhino.Token{
		rhino.Token_COMMA,
		rhino.Token_AND,
		rhino.Token_OR,
		rhino.Token_COALESCE,
		rhino.Token_HOOK,
		rhino.Token_IF,
		rhino.Token_INSTANCEOF,
		rhino.Token_ASSIGN,
		rhino.Token_TYPEOF,
		rhino.Token_NOT,
		rhino.Token_EQ,
		rhino.Token_NE,
		rhino.Token_SHEQ,
		rhino.Token_SHNE,
		rhino.Token_CALL,
		rhino.Token_OPTCHAIN_CALL,
		rhino.Token_OPTCHAIN_GETELEM,
		rhino.Token_OPTCHAIN_GETPROP,
	} {
		flags[t] |= guardedCallback_flagEmptyDescend
	}
	for _, t := range []rhino.Token{
		rhino.Token_AND,
		rhino.Token_OR,
		rhino.Token_COALESCE,
		rhino.Token_HOOK,
		rhino.Token_IF,
		rhino.Token_BLOCK,
		rhino.Token_SCRIPT,
		rhino.Token_OPTCHAIN_CALL,
		rhino.Token_OPTCHAIN_GETELEM,
		rhino.Token_OPTCHAIN_GETPROP,
	} {
		flags[t] |= guardedCallback_flagCanHaveGuards
	}
	return flags
}()

// guardedCallback_CAN_HAVE_GUARDS: tokens that are allowed to have guards
// on them (no point doing a hash lookup on any other type of node).
//
// Java: GuardedCallback.CAN_HAVE_GUARDS
var guardedCallback_CAN_HAVE_GUARDS = jutil.NewEnumSet(
	rhino.Token_AND,
	rhino.Token_OR,
	rhino.Token_COALESCE,
	rhino.Token_HOOK,
	rhino.Token_IF,
	rhino.Token_BLOCK,
	rhino.Token_SCRIPT,
	rhino.Token_OPTCHAIN_CALL,
	rhino.Token_OPTCHAIN_GETELEM,
	rhino.Token_OPTCHAIN_GETPROP)
