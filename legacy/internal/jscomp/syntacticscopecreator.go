// Port of com/google/javascript/jscomp/SyntacticScopeCreator.java @ closure-compiler v20260415
//
// The syntactic scope creator scans the parse tree to create a Scope object
// containing all variable declarations in that scope. This class adds support
// for block-level scopes introduced in ECMAScript 6.
//
// This implementation is not thread-safe.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// RedeclarationHandler is the injectable interface for handling duplicate
// variable declarations. Mirrors SyntacticScopeCreator.RedeclarationHandler.
type RedeclarationHandler interface {
	// OnRedeclaration is called when a variable is declared more than once
	// in the same scope.
	OnRedeclaration(s *UntypedScope, name string, n *rhino.Node)
}

// DefaultRedeclarationHandler silently ignores duplicate declarations.
// Mirrors SyntacticScopeCreator.DefaultRedeclarationHandler in Java.
type DefaultRedeclarationHandler struct{}

// OnRedeclaration is a no-op — Java's default handler also does nothing.
func (DefaultRedeclarationHandler) OnRedeclaration(_ *UntypedScope, _ string, _ *rhino.Node) {}

// arguments is the implicit variable declared in every function scope.
const argumentsVar = "arguments"

// SyntacticScopeCreator scans the parse tree to create UntypedScope objects
// containing all variable declarations in a scope. It implements
// UntypedScopeCreator and mirrors
// com/google/javascript/jscomp/SyntacticScopeCreator.java.
type SyntacticScopeCreator struct {
	// redeclarationHandler is called when a name is declared multiple times
	// in the same scope.
	redeclarationHandler RedeclarationHandler
	// treatProvidesAsRedeclarations controls whether goog.provide calls that
	// re-declare an existing var are reported to the handler.
	treatProvidesAsRedeclarations bool
}

// NewSyntacticScopeCreator returns a SyntacticScopeCreator with the default
// (no-op) redeclaration handler. Mirrors
// SyntacticScopeCreator(AbstractCompiler) in Java. The AbstractCompiler
// parameter is omitted because it is not yet ported into this package.
func NewSyntacticScopeCreator() *SyntacticScopeCreator {
	return &SyntacticScopeCreator{
		redeclarationHandler: DefaultRedeclarationHandler{},
	}
}

// NewSyntacticScopeCreatorWithHandler returns a SyntacticScopeCreator that
// delegates duplicate declarations to handler.
// Mirrors SyntacticScopeCreator(AbstractCompiler, RedeclarationHandler).
func NewSyntacticScopeCreatorWithHandler(handler RedeclarationHandler) *SyntacticScopeCreator {
	return &SyntacticScopeCreator{
		redeclarationHandler: handler,
	}
}

// CreateUntypedScope implements UntypedScopeCreator.
// It creates a UntypedScope for n, populating it with all variables declared
// directly in n's scope. Mirrors SyntacticScopeCreator.createScope(Node, Scope).
func (c *SyntacticScopeCreator) CreateUntypedScope(n *rhino.Node, parent *UntypedScope) *UntypedScope {
	var scope *UntypedScope
	if parent == nil {
		scope = CreateGlobalUntypedScope(n)
	} else {
		scope = CreateChildUntypedScope(parent, n)
	}
	ss := &scopeScanner{
		creator: c,
		scope:   scope,
	}
	ss.populate()
	return scope
}

// scopeScanner is the inner scanner that populates a UntypedScope.
// Mirrors SyntacticScopeCreator.ScopeScanner.
type scopeScanner struct {
	creator *SyntacticScopeCreator
	scope   *UntypedScope
}

// populate fills the scope with variable declarations. Mirrors
// ScopeScanner.populate() in Java.
func (ss *scopeScanner) populate() {
	n := ss.scope.GetRootNode()
	switch n.Token() {
	case rhino.FUNCTION:
		fnNameNode := n.FirstChild()
		if fnNameNode == nil {
			return
		}
		args := fnNameNode.Next()
		if args == nil {
			return
		}
		// Declare function parameters.
		if args.Token() == rhino.PARAM_LIST {
			ss.declareLHS(ss.scope, args)
		}
		// Bleed the function name into the scope if it is a function expression
		// and has not already been declared via a parameter.
		fnName := fnNameNode.GetString()
		if fnName != "" && isFunctionExpression(n) {
			ss.declareVar(ss.scope, fnNameNode)
		}
		return

	case rhino.CLASS:
		classNameNode := n.FirstChild()
		if classNameNode != nil && !classNameNode.IsEmpty() && isClassExpression(n) {
			ss.declareVar(ss.scope, classNameNode)
		}
		return

	case rhino.ROOT, rhino.SCRIPT:
		ss.scanVars(n, ss.scope, ss.scope)
		return

	case rhino.MODULE_BODY:
		ss.scanVars(n, ss.scope, ss.scope)
		return

	case rhino.FOR, rhino.FOR_OF, rhino.FOR_AWAIT_OF, rhino.FOR_IN, rhino.SWITCH_BODY:
		ss.scanVars(n, nil, ss.scope)
		return

	case rhino.BLOCK:
		if isFunctionBlock(n) || isClassStaticBlock(n) {
			ss.scanVars(n, ss.scope, ss.scope)
		} else {
			ss.scanVars(n, nil, ss.scope)
		}
		return

	case rhino.MEMBER_FIELD_DEF:
		// MEMBER_FIELD_DEF and COMPUTED_FIELD_DEF scopes only created to scope
		// `this` and `super` correctly.
		return
	}
}

// declareLHS declares all LHS names in n into scope s.
// Mirrors ScopeScanner.declareLHS(Scope, Node).
func (ss *scopeScanner) declareLHS(s *UntypedScope, n *rhino.Node) {
	// Fast path: single NAME child.
	first := n.FirstChild()
	if first != nil && first.Next() == nil && first.Token() == rhino.NAME {
		ss.declareVar(s, first)
		return
	}
	// General: visit all LHS names.
	visitLhsNodesInNode(n, func(lhs *rhino.Node) {
		ss.declareVar(s, lhs)
	})
}

// scanVars traverses n and declares variables into hoistScope (for var/import)
// and blockScope (for let/const/function/class).
// Mirrors ScopeScanner.scanVars(Node, Scope, Scope).
func (ss *scopeScanner) scanVars(n *rhino.Node, hoistScope, blockScope *UntypedScope) {
	switch n.Token() {
	case rhino.VAR:
		if hoistScope != nil {
			ss.declareLHS(hoistScope, n)
		}
		return

	case rhino.LET, rhino.CONST:
		if blockScope != nil {
			ss.declareLHS(blockScope, n)
		}
		return

	case rhino.IMPORT:
		if hoistScope != nil {
			ss.declareLHS(hoistScope, n)
		}
		return

	case rhino.EXPORT:
		// The first child of EXPORT can be a declaration.
		first := n.FirstChild()
		if first != nil {
			ss.scanVars(first, hoistScope, blockScope)
		}
		return

	case rhino.FUNCTION:
		if isFunctionExpression(n) || blockScope == nil {
			return
		}
		fnNameNode := n.FirstChild()
		if fnNameNode == nil {
			return
		}
		fnName := fnNameNode.GetString()
		if fnName == "" {
			return
		}
		ss.declareVar(blockScope, fnNameNode)
		return // do not descend into function children

	case rhino.CLASS:
		if isClassExpression(n) || blockScope == nil {
			return
		}
		classNameNode := n.FirstChild()
		if classNameNode == nil {
			return
		}
		className := classNameNode.GetString()
		if className == "" {
			return
		}
		ss.declareVar(blockScope, classNameNode)
		return // do not descend into class children

	case rhino.CATCH:
		// The first child is the catch var; the second is the block body.
		if blockScope != nil {
			ss.declareLHS(blockScope, n)
		}
		block := n.GetSecondChild()
		if block != nil {
			ss.scanVars(block, hoistScope, blockScope)
		}
		return
	}

	// For statement-level nodes, recurse into children that may contain decls.
	if isShallowStatementTree(n) {
		isBlockStart := blockScope != nil && n == blockScope.GetRootNode()
		enteringNewBlock := !isBlockStart && createsBlockScope(n)
		if enteringNewBlock && hoistScope == nil {
			return
		}
		childBlockScope := blockScope
		if enteringNewBlock {
			childBlockScope = nil
		}
		for child := n.FirstChild(); child != nil; child = child.Next() {
			ss.scanVars(child, hoistScope, childBlockScope)
		}
	}
}

// declareVar declares name n in scope s, or calls the redeclaration handler
// if the name is already declared. Mirrors ScopeScanner.declareVar(Scope, Node).
func (ss *scopeScanner) declareVar(s *UntypedScope, n *rhino.Node) {
	if n.Token() != rhino.NAME && n.Token() != rhino.IMPORT_STAR {
		return
	}
	name := n.GetString()
	existing := s.GetOwnSlot(name)
	if existing != nil {
		if existing.GetNode() == n {
			return
		}
		// Handle implicit goog namespace replacement.
		if existing.IsImplicitGoogNamespace() {
			s.RemoveVar(name)
			existing = nil
		}
	}
	isFnScope := s.IsFunctionScope()
	isFnBlockScope := s.IsFunctionBlockScope()
	if existing != nil ||
		!isShadowingAllowed(name, s) ||
		((isFnScope || isFnBlockScope) && name == argumentsVar) {
		ss.creator.redeclarationHandler.OnRedeclaration(s, name, n)
	} else {
		s.Declare(name, n)
	}
}

// isShadowingAllowed returns false when name conflicts with a parameter in an
// enclosing function scope. Mirrors ScopeScanner.isShadowingAllowed.
func isShadowingAllowed(name string, s *UntypedScope) bool {
	if s.IsFunctionBlockScope() {
		parent := s.GetParent()
		if parent != nil {
			maybeParam := parent.GetOwnSlot(name)
			if maybeParam != nil && maybeParam.IsParam() {
				return false
			}
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// Helpers — mirrors of Java NodeUtil predicates used by ScopeScanner.
// These are minimal implementations; full NodeUtil mirrors live in nodeutil.go.
// ---------------------------------------------------------------------------

// isFunctionExpression reports whether n is a function expression (i.e., a
// FUNCTION node that is not a statement). Mirrors NodeUtil.isFunctionExpression.
func isFunctionExpression(n *rhino.Node) bool {
	if n.Token() != rhino.FUNCTION {
		return false
	}
	parent := n.Parent()
	if parent == nil {
		return false
	}
	switch parent.Token() {
	case rhino.EXPR_RESULT, rhino.CALL, rhino.ASSIGN,
		rhino.RETURN, rhino.VAR, rhino.LET, rhino.CONST,
		rhino.COMMA, rhino.AND, rhino.OR:
		return true
	case rhino.STRING_KEY, rhino.MEMBER_FUNCTION_DEF,
		rhino.GETTER_DEF, rhino.SETTER_DEF:
		return true
	}
	return false
}

// isClassExpression reports whether n is a class expression.
// Mirrors NodeUtil.isClassExpression.
func isClassExpression(n *rhino.Node) bool {
	if n.Token() != rhino.CLASS {
		return false
	}
	parent := n.Parent()
	if parent == nil {
		return false
	}
	switch parent.Token() {
	case rhino.EXPR_RESULT, rhino.ASSIGN,
		rhino.RETURN, rhino.VAR, rhino.LET, rhino.CONST:
		return true
	}
	return false
}

// isFunctionBlock reports whether n is the body block of a FUNCTION node.
// Mirrors NodeUtil.isFunctionBlock.
func isFunctionBlock(n *rhino.Node) bool {
	if n.Token() != rhino.BLOCK {
		return false
	}
	parent := n.Parent()
	return parent != nil && parent.Token() == rhino.FUNCTION
}

// isClassStaticBlock reports whether n is a static block of a CLASS node.
// Mirrors NodeUtil.isClassStaticBlock — simplified: always false since
// BLOCK children of CLASS have no Go equivalent constant yet.
func isClassStaticBlock(_ *rhino.Node) bool {
	return false
}

// createsBlockScope reports whether n introduces a new block-level scope in ES6.
// Mirrors NodeUtil.createsBlockScope.
func createsBlockScope(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.BLOCK, rhino.FOR, rhino.FOR_IN, rhino.FOR_OF,
		rhino.FOR_AWAIT_OF, rhino.SWITCH, rhino.SWITCH_BODY:
		return true
	}
	return false
}

// isShallowStatementTree reports whether we should recurse into n's children
// to look for declarations. Mirrors NodeUtil.isShallowStatementTree.
func isShallowStatementTree(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.BLOCK, rhino.ROOT, rhino.SCRIPT, rhino.MODULE_BODY,
		rhino.IF, rhino.WHILE, rhino.DO_WHILE, rhino.FOR,
		rhino.FOR_IN, rhino.FOR_OF, rhino.FOR_AWAIT_OF,
		rhino.LABEL, rhino.SWITCH, rhino.SWITCH_BODY,
		rhino.CASE, rhino.DEFAULT_CASE, rhino.TRY, rhino.CATCH:
		return true
	}
	return false
}

// visitLhsNodesInNode calls f on every NAME or IMPORT_STAR node that is an
// LHS target within n. Mirrors NodeUtil.visitLhsNodesInNode.
func visitLhsNodesInNode(n *rhino.Node, f func(*rhino.Node)) {
	for child := n.FirstChild(); child != nil; child = child.Next() {
		visitLhsNode(child, f)
	}
}

// visitLhsNode recursively visits all LHS name nodes in pattern node n.
func visitLhsNode(n *rhino.Node, f func(*rhino.Node)) {
	switch n.Token() {
	case rhino.NAME, rhino.IMPORT_STAR:
		f(n)
	case rhino.ARRAY_PATTERN, rhino.OBJECT_PATTERN, rhino.PARAM_LIST,
		rhino.DESTRUCTURING_LHS:
		for child := n.FirstChild(); child != nil; child = child.Next() {
			visitLhsNode(child, f)
		}
	case rhino.DEFAULT_VALUE:
		// The LHS is the first child (before the default value).
		first := n.FirstChild()
		if first != nil {
			visitLhsNode(first, f)
		}
	case rhino.ITER_REST, rhino.OBJECT_REST:
		first := n.FirstChild()
		if first != nil {
			visitLhsNode(first, f)
		}
	case rhino.STRING_KEY, rhino.COMPUTED_PROP:
		// In object destructuring, the value (second child) is the target.
		last := n.LastChild()
		if last != nil {
			visitLhsNode(last, f)
		}
	}
}
