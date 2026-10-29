// Port of com/google/javascript/jscomp/Denormalize.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	denormalizeVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_EXPR_RESULT)
		s.Add(rhino.Token_VAR)
		s.Add(rhino.Token_ASSIGN)
		s.Add(rhino.Token_OR)
		s.Add(rhino.Token_AND)
		s.Add(rhino.Token_COALESCE)
		return s
	}()
	denormalizeVarInlinerVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_NAME)
		return s
	}()
)

// Go-only: this pass uses the typed AST view (internal/ast). visit hands
// each rewrite the kind it works on (ast.Statement, ast.Assign, ast.Binary),
// and the rewrites read their shapes by role (stmt.Next(), forNode.Init(),
// assign.Target(), op.Left()...) instead of child positions, so the port no
// longer follows the Java line by line. Java's isCollapsableAssign and
// isCollapsableLogicalAssign, with the name check their callers make, are
// one comma-ok function each that returns the node to rewrite. Step 3 is
// Go-only code (denormalizeVarInliner, in place of Java's
// ReferenceCollector); its scope bookkeeping uses the rhino API.

// Denormalize: the goal with this pass is to reverse the simplifications
// done in the normalization pass that are not handled by other passes (such
// as CollapseVariableDeclarations) to avoid making the resulting code
// larger.
//
// Currently this pass only does a few things:
//
// 1. Push statements into for-loop initializer. This: var a = 0;
// for(;a<0;a++) {} becomes: for(var a = 0;a<0;a++) {}
//
// 2. Fold assignments like x = x + 1 into x += 1
//
// 3. Inline 'var' keyword. For instance: `var x; if (y) { x = 0; }` becomes
// `if (y) { var x = 0; }`, effectively undoing what HoistVarsOutOfBlocks
// does.
type Denormalize struct {
	compiler         AbstractCompiler
	outputFeatureSet *parser.FeatureSet
	hasLogicalAssign bool
	varInliner       denormalizeVarInliner
}

// Java: Denormalize(AbstractCompiler, FeatureSet)
func NewDenormalize(compiler AbstractCompiler, outputFeatureSet *parser.FeatureSet) *Denormalize {
	return &Denormalize{
		compiler:         compiler,
		outputFeatureSet: outputFeatureSet,
		hasLogicalAssign: outputFeatureSet != nil && outputFeatureSet.Has(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT),
	}
}

// Java: Denormalize.process(Node, Node)
func (d *Denormalize) Process(externs rhino.Node, root rhino.Node) {
	d.hasLogicalAssign = d.outputFeatureSet != nil && d.outputFeatureSet.Has(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	NodeTraversal_Traverse(d.compiler, root, d)
	// Don't inline the VAR declaration if this compilation involves old-style ctemplates.
	if d.compiler.GetOptions().GetSyntheticBlockStartMarker() == nil {
		d.varInliner.reset()
		b := NodeTraversal_Builder{
			compiler:     d.compiler,
			callback:     &d.varInliner,
			scopeCreator: NewSyntacticScopeCreator(d.compiler),
			obeyDestructuringAndDefaultValueExecutionOrder: true,
		}
		b.Traverse(root)
		d.varInliner.reset()
	}
}

type denormalizeVarState struct {
	declNode   rhino.Node
	assignNode rhino.Node
}

type denormalizeVarInliner struct {
	states          map[*Var]denormalizeVarState
	scopesWithDecls map[rhino.Node]*Scope
}

func (inliner *denormalizeVarInliner) reset() {
	if len(inliner.states) > 1024 {
		inliner.states = nil
	} else if len(inliner.states) > 0 {
		clear(inliner.states)
	}
	if len(inliner.scopesWithDecls) > 256 {
		inliner.scopesWithDecls = nil
	} else if len(inliner.scopesWithDecls) > 0 {
		clear(inliner.scopesWithDecls)
	}
}

func (*denormalizeVarInliner) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

func (*denormalizeVarInliner) nodeTraversalNoShouldTraverse() {}
func (*denormalizeVarInliner) onlyNameLeaves()                {}
func (*denormalizeVarInliner) visitTokens() *ast.TokenSet     { return &denormalizeVarInlinerVisitTokens }
func (*denormalizeVarInliner) EnterScope(t *NodeTraversal)    {}

func (inliner *denormalizeVarInliner) ExitScope(t *NodeTraversal) {
	if len(inliner.scopesWithDecls) == 0 {
		return
	}
	scopeRoot := t.GetScopeRoot()
	scope := inliner.scopesWithDecls[scopeRoot]
	if scope == nil {
		return
	}
	delete(inliner.scopesWithDecls, scopeRoot)
	changed := false
	for _, v := range scope.varList {
		if v == nil {
			continue
		}
		st, ok := inliner.states[v]
		if !ok {
			continue
		}
		delete(inliner.states, v)
		if !st.declNode.IsNil() && !st.assignNode.IsNil() && denormalize_inlineVar(st.declNode, st.assignNode) {
			changed = true
		}
	}
	if changed {
		t.ReportCodeChange()
	}
}

// Visit records, for each variable declared in a function body, its last
// declaration without value (var x;) and its first assignment (x = value)
// in the same hoist scope.
func (inliner *denormalizeVarInliner) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !n.IsName() || parent.IsNil() {
		return
	}
	switch parent.GetToken() {
	case rhino.Token_VAR:
		if n.GetFirstChild().IsNil() && NodeUtil_IsStatement(parent) {
			id := n.GetStringID()
			if id == 0 || id == rhino.RhinoStringPool_EmptyID {
				return
			}
			v := t.LookupVarByID(id)
			if v != nil {
				vRoot := v.scope.GetRootNode()
				if vRoot.IsBlock() && vRoot.GetParent().IsFunction() {
					if inliner.states == nil {
						inliner.states = make(map[*Var]denormalizeVarState, 32)
					}
					st := inliner.states[v]
					st.declNode = n
					inliner.states[v] = st
					if inliner.scopesWithDecls == nil {
						inliner.scopesWithDecls = make(map[rhino.Node]*Scope, 8)
					}
					inliner.scopesWithDecls[vRoot] = v.scope
				}
			}
		}
	case rhino.Token_ASSIGN:
		if parent.GetFirstChild() == n {
			id := n.GetStringID()
			if id == 0 || id == rhino.RhinoStringPool_EmptyID {
				return
			}
			scope := t.resolveLookupScope(len(t.scopes) - 1)
			v := scope.GetVarByID(id)
			if v != nil && scope.GetClosestHoistScope() == v.scope {
				vRoot := v.scope.GetRootNode()
				if vRoot.IsBlock() && vRoot.GetParent().IsFunction() {
					if inliner.states == nil {
						inliner.states = make(map[*Var]denormalizeVarState, 32)
					}
					if st := inliner.states[v]; st.assignNode.IsNil() {
						st.assignNode = n
						inliner.states[v] = st
					}
				}
			}
		}
	}
}

// denormalize_inlineVar moves the var keyword of `var x;` to an assignment
// `x = value;` (a statement) in the same function: decl is the x of the var,
// target the x assigned. It returns false, and changes nothing, when the
// assignment is not a statement.
//
// (Go-only: the shared rewrite of afterExitScope and denormalizeVarInliner.)
func denormalize_inlineVar(decl rhino.Node, target rhino.Node) bool {
	assign, ok := ast.AsAssign(target.GetParent())
	jutil.CheckState(ok)
	stmt, ok := ast.AsExprStmt(assign.Node().GetParent())
	if !ok {
		return false
	}
	lhs, rhs := assign.Target(), assign.Value()
	stmt.Node().ReplaceWith(rhino.IR_VarWithValue(lhs.Detach().Node(), rhs.Detach().Node()))

	var_ := decl.GetParent()
	jutil.CheckStateMsg1(var_.IsVar(), "%s", var_)
	NodeUtil_RemoveChild(var_, decl)
	return true
}

// Java: Denormalize.shouldTraverse(NodeTraversal, Node, Node)
func (*Denormalize) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

func (*Denormalize) nodeTraversalNoShouldTraverse() {}
func (*Denormalize) noVisitLeaves()                 {}
func (*Denormalize) visitTokens() *ast.TokenSet     { return &denormalizeVisitTokens }

// (Go: a token switch, not an ast.Visitor: the pass does little per node, and
// the visitor's dispatch made this traversal about 10% slower.)
//
// Java: Denormalize.visit(NodeTraversal, Node, Node)
func (d *Denormalize) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_EXPR_RESULT, rhino.Token_VAR:
		stmt := ast.MustStatement(n)
		d.MaybeCollapseIntoForStatements(stmt)
	case rhino.Token_ASSIGN:
		assign := ast.MustAssign(n)
		d.MaybeCollapseAssignShorthand(assign)
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		if d.hasLogicalAssign {
			op := ast.MustBinary(n)
			d.MaybeCollapseLogicalAssignShorthand(t, op)
		}
	}
}

// MaybeCollapseIntoForStatements collapses VARs and EXPR_RESULT node into
// FOR loop initializers where possible.
//
// Java: Denormalize.maybeCollapseIntoForStatements(Node, Node)
func (d *Denormalize) MaybeCollapseIntoForStatements(stmt ast.Statement) {
	n := stmt.Node()
	// Is the next statement a valid FOR?
	nextSibling, ok := stmt.Next()
	if !ok {
		return
	}

	// Only SCRIPT, BLOCK, and LABELs can have FORs that can be collapsed into.
	// LABELs are not supported here.
	parent := n.GetParent()
	if parent.IsNil() || !NodeUtil_IsStatementBlock(parent) {
		return
	}

	// Is the current node something that can be in a for loop initializer?
	exprStmt, isExprStmt := stmt.AsExprStmt()
	varNode, isVar := denormalize_asVar(stmt)
	if !isExprStmt && !isVar {
		return
	}

	if forNode, ok := nextSibling.AsEnhancedFor(); ok && forNode.Token() != rhino.Token_FOR_AWAIT_OF {
		// (A for-in or for-of loop: Java leaves for-await-of alone.)
		forVar, ok := forNode.Target().AsName()
		if !ok || !isVar {
			return
		}
		if name, ok := denormalize_onlyNameWithoutValue(varNode); ok && forVar.Node().GetStringID() == name.Node().GetStringID() {
			// OK, the names match, and the var declaration does not have an
			// initializer. Move it into the loop.
			varNode.Detach()
			forNode.Target().ReplaceWith(varNode.AsForTarget())
			d.compiler.ReportChangeToEnclosingScope(parent)
		}
	} else if forNode, ok := nextSibling.AsFor(); ok {
		if _, hasInit := forNode.Init(); hasInit {
			return
		}

		// Does the current node contain an in operator?  If so, embedding
		// the expression in a for loop can cause some JavaScript parsers (such
		// as the PlayStation 3's browser based on Access's NetFront
		// browser) to fail to parse the code.
		// See bug 1778863 for details.
		if NodeUtil_Has(n, (rhino.Node).IsIn, denormalize_alwaysTrue) {
			return
		}

		// Move the current node into the FOR loop initializer.
		n.Detach()
		var newInitializer ast.ForInit
		if isVar {
			newInitializer = varNode.AsForInit()
		} else {
			// Extract the expression from EXPR_RESULT node.
			newInitializer = exprStmt.Expr().Detach().AsForInit()
		}
		forNode.SetInit(newInitializer)

		d.compiler.ReportChangeToEnclosingScope(forNode.Node())
	}
}

// denormalize_asVar returns stmt as a var declaration; ok is false for
// anything else, let and const included.
func denormalize_asVar(stmt ast.Statement) (ast.NameDeclaration, bool) {
	decl, ok := stmt.AsNameDeclaration()
	return decl, ok && decl.Token() == rhino.Token_VAR
}

// denormalize_onlyNameWithoutValue returns the name declared by decl if it
// declares a single name without a value (var x;); ok is false otherwise.
func denormalize_onlyNameWithoutValue(decl ast.NameDeclaration) (ast.Name, bool) {
	declarator, ok := decl.OnlyDeclarator()
	if !ok {
		return ast.Name{}, false
	}
	name, ok := declarator.AsName()
	if !ok {
		return ast.Name{}, false
	}
	_, hasValue := name.Init()
	return name, !hasValue
}

// denormalize_alwaysTrue is Guava's Predicates.alwaysTrue().
func denormalize_alwaysTrue(rhino.Node) bool { return true }

// MaybeCollapseAssignShorthand rewrites x = x op y to x op= y.
//
// Java: Denormalize.maybeCollapseAssignShorthand(Node, Node)
func (d *Denormalize) MaybeCollapseAssignShorthand(assign ast.Assign) {
	op, ok := d.IsCollapsableAssign(assign)
	if !ok {
		return
	}
	n := assign.Node()
	parent := n.GetParent()
	// x op y becomes x op= y in place (same children, new token), and takes
	// the place of the assignment.
	opNode := op.Node()
	opNode.SetToken(d.GetAssignOpFromOp(opNode))
	opNode.Detach()
	opNode.SetJSDocInfo(n.GetJSDocInfo())
	n.ReplaceWith(opNode)
	d.compiler.ReportChangeToEnclosingScope(parent)
}

// MaybeCollapseLogicalAssignShorthand rewrites x || (x = y) to x ||= y, and
// the same for && and ??.
//
// Java: Denormalize.maybeCollapseLogicalAssignShorthand(NodeTraversal, Node, Node)
func (d *Denormalize) MaybeCollapseLogicalAssignShorthand(t *NodeTraversal, op ast.Binary) {
	assign, ok := d.IsCollapsableLogicalAssign(op)
	if !ok {
		return
	}
	n := op.Node()
	parent := n.GetParent()
	// x = y becomes x ||= y in place (same children, new token), and takes
	// the place of the ||.
	assignNode := assign.Node()
	assignNode.SetToken(d.GetAssignOpFromOp(n))
	assignNode.Detach()
	assignNode.SetJSDocInfo(n.GetJSDocInfo())
	n.ReplaceWith(assignNode)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT, d.compiler)
	d.compiler.ReportChangeToEnclosingScope(parent)
}

// IsCollapsableAssign returns the operation of assign if assign is
// x = x op y, with an op that has an assignment form (x op= y); ok is false
// otherwise.
//
// (Go: Java's isCollapsableAssign returns a boolean, and its caller checks
// that the two names are the same.)
//
// Java: Denormalize.isCollapsableAssign(Node)
func (d *Denormalize) IsCollapsableAssign(assign ast.Assign) (ast.Binary, bool) {
	target, ok := assign.Target().AsName()
	if !ok {
		return ast.Binary{}, false
	}
	op, ok := assign.Value().AsBinary()
	if !ok || !d.HasCorrespondingAssignmentOp(op.Node()) {
		return ast.Binary{}, false
	}
	left, ok := op.Left().AsName()
	return op, ok && target.Node().GetStringID() == left.Node().GetStringID()
}

// IsCollapsableLogicalAssign returns the assignment of op if op is
// x || (x = y), x && (x = y) or x ?? (x = y), and the output can have
// logical assignments; ok is false otherwise.
//
// (Go: Java's isCollapsableLogicalAssign returns a boolean, and its caller
// checks that the two names are the same.)
//
// Java: Denormalize.isCollapsableLogicalAssign(Node)
func (d *Denormalize) IsCollapsableLogicalAssign(op ast.Binary) (ast.Assign, bool) {
	if !d.outputFeatureSet.Has(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT) {
		return ast.Assign{}, false
	}
	switch op.Token() {
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
	default:
		return ast.Assign{}, false
	}
	left, ok := op.Left().AsName()
	if !ok {
		return ast.Assign{}, false
	}
	assign, ok := op.Right().AsAssign()
	if !ok {
		return ast.Assign{}, false
	}
	target, ok := assign.Target().AsName()
	return assign, ok && left.Node().GetStringID() == target.Node().GetStringID()
}

// Java: Denormalize.getAssignOpFromOp(Node)
func (d *Denormalize) GetAssignOpFromOp(n rhino.Node) rhino.Token {
	switch n.GetToken() {
	case rhino.Token_BITOR:
		return rhino.Token_ASSIGN_BITOR
	case rhino.Token_BITXOR:
		return rhino.Token_ASSIGN_BITXOR
	case rhino.Token_BITAND:
		return rhino.Token_ASSIGN_BITAND
	case rhino.Token_LSH:
		return rhino.Token_ASSIGN_LSH
	case rhino.Token_RSH:
		return rhino.Token_ASSIGN_RSH
	case rhino.Token_URSH:
		return rhino.Token_ASSIGN_URSH
	case rhino.Token_ADD:
		return rhino.Token_ASSIGN_ADD
	case rhino.Token_SUB:
		return rhino.Token_ASSIGN_SUB
	case rhino.Token_MUL:
		return rhino.Token_ASSIGN_MUL
	case rhino.Token_EXPONENT:
		return rhino.Token_ASSIGN_EXPONENT
	case rhino.Token_DIV:
		return rhino.Token_ASSIGN_DIV
	case rhino.Token_MOD:
		return rhino.Token_ASSIGN_MOD
	case rhino.Token_OR:
		return rhino.Token_ASSIGN_OR
	case rhino.Token_AND:
		return rhino.Token_ASSIGN_AND
	case rhino.Token_COALESCE:
		return rhino.Token_ASSIGN_COALESCE
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected operator: " + n.String()))
	}
}

// Java: Denormalize.hasCorrespondingAssignmentOp(Node)
func (d *Denormalize) HasCorrespondingAssignmentOp(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_BITOR, rhino.Token_BITXOR, rhino.Token_BITAND, rhino.Token_LSH, rhino.Token_RSH,
		rhino.Token_URSH, rhino.Token_ADD, rhino.Token_SUB, rhino.Token_MUL, rhino.Token_DIV,
		rhino.Token_MOD:
		return true
	default:
		return false
	}
}
