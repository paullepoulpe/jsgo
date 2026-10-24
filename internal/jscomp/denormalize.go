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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
		inliner := newDenormalizeVarInliner()
		sc := NewSyntacticScopeCreator(d.compiler)
		b := NodeTraversal_Builder{
			compiler:     d.compiler,
			callback:     inliner,
			scopeCreator: sc,
			obeyDestructuringAndDefaultValueExecutionOrder: true,
		}
		b.Traverse(root)
	}
}

type denormalizeVarState struct {
	declNode   rhino.Node
	assignNode rhino.Node
}

type denormalizeVarInliner struct {
	states          []denormalizeVarState
	scopesWithDecls map[rhino.Node]*Scope
	statesBuf       [32]denormalizeVarState
}

func newDenormalizeVarInliner() *denormalizeVarInliner {
	inliner := &denormalizeVarInliner{}
	inliner.states = inliner.statesBuf[:0]
	return inliner
}

func (inliner *denormalizeVarInliner) getOrAllocState(v *Var) *denormalizeVarState {
	idx := v.nameHash
	if idx == 0 {
		idx = int32(len(inliner.states) + 1)
		v.nameHash = idx
		inliner.states = append(inliner.states, denormalizeVarState{})
	}
	return &inliner.states[idx-1]
}

func (inliner *denormalizeVarInliner) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

func (inliner *denormalizeVarInliner) EnterScope(t *NodeTraversal) {}

func (inliner *denormalizeVarInliner) ExitScope(t *NodeTraversal) {
	if len(inliner.scopesWithDecls) == 0 {
		return
	}
	scopeRoot := t.GetScopeRoot()
	scope := inliner.scopesWithDecls[scopeRoot]
	if scope == nil {
		return
	}
	changed := false
	for _, v := range scope.varList {
		if v == nil {
			continue
		}
		idx := v.nameHash
		if idx == 0 {
			continue
		}
		st := &inliner.states[idx-1]
		if !st.declNode.IsNil() && !st.assignNode.IsNil() {
			lhs := st.assignNode
			assignNode := lhs.GetParent()
			if assignNode.GetParent().IsExprResult() {
				rhs := lhs.GetNext()
				assignNode.GetParent().ReplaceWith(rhino.IR_VarWithValue(lhs.Detach(), rhs.Detach()))

				var_ := st.declNode.GetParent()
				jutil.CheckStateMsg1(var_.IsVar(), "%s", var_)
				NodeUtil_RemoveChild(var_, st.declNode)

				changed = true
			}
		}
	}
	if changed {
		t.ReportCodeChange()
	}
}

func (inliner *denormalizeVarInliner) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !n.IsName() || parent.IsNil() {
		return
	}
	switch parent.GetToken() {
	case rhino.Token_VAR:
		if !n.HasChildren() && NodeUtil_IsStatement(parent) {
			name := n.GetString()
			if name == "" {
				return
			}
			v := t.LookupVar(name)
			if v != nil {
				vRoot := v.scope.GetRootNode()
				if vRoot.IsBlock() && vRoot.GetParent().IsFunction() {
					st := inliner.getOrAllocState(v)
					st.declNode = n
					if inliner.scopesWithDecls == nil {
						inliner.scopesWithDecls = make(map[rhino.Node]*Scope, 8)
					}
					inliner.scopesWithDecls[vRoot] = v.scope
				}
			}
		}
	case rhino.Token_ASSIGN:
		if parent.GetFirstChild() == n {
			name := n.GetString()
			if name == "" {
				return
			}
			scope := t.resolveLookupScope(len(t.scopes) - 1)
			v := scope.GetVar(name)
			if v != nil && scope.GetClosestHoistScope() == v.scope {
				vRoot := v.scope.GetRootNode()
				if vRoot.IsBlock() && vRoot.GetParent().IsFunction() {
					st := inliner.getOrAllocState(v)
					if st.assignNode.IsNil() {
						st.assignNode = n
					}
				}
			}
		}
	}
}

// Java: Denormalize.shouldTraverse(NodeTraversal, Node, Node)
func (d *Denormalize) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

// AfterExitScope implements step 3 (inlining the var keyword).
//
// Java: Denormalize.afterExitScope(NodeTraversal, ReferenceMap)
func (d *Denormalize) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	scopeRoot := t.GetScopeRoot()
	if scopeRoot.IsBlock() && scopeRoot.GetParent().IsFunction() {
		changed := false
		for v := range t.GetScope().GetVarIterable() {
			references := referenceMap.GetReferences(v)
			var declaration *Reference
			var assign *Reference
			for r := range references.All() {
				if r.IsVarDeclaration() &&
					NodeUtil_IsStatement(r.GetNode().GetParent()) &&
					!r.IsInitializingDeclaration() {
					declaration = r
				} else if assign == nil &&
					r.IsSimpleAssignmentToName() &&
					r.GetScope().GetClosestHoistScope() == t.GetScope() {
					assign = r
				}
			}
			if declaration != nil && assign != nil {
				lhs := assign.GetNode()
				assignNode := lhs.GetParent()
				if assignNode.GetParent().IsExprResult() {
					rhs := lhs.GetNext()
					assignNode.GetParent().ReplaceWith(rhino.IR_VarWithValue(lhs.Detach(), rhs.Detach()))

					var_ := declaration.GetNode().GetParent()
					jutil.CheckStateMsg1(var_.IsVar(), "%s", var_)
					NodeUtil_RemoveChild(var_, declaration.GetNode())

					changed = true
				}
			}
		}

		if changed {
			t.ReportCodeChange()
		}
	}
}

// Java: Denormalize.visit(NodeTraversal, Node, Node)
func (d *Denormalize) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_EXPR_RESULT, rhino.Token_VAR:
		d.MaybeCollapseIntoForStatements(n, parent)
	case rhino.Token_ASSIGN:
		d.MaybeCollapseAssignShorthand(n, parent)
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		if d.hasLogicalAssign {
			d.MaybeCollapseLogicalAssignShorthand(t, n, parent)
		}
	}
}

// MaybeCollapseIntoForStatements collapses VARs and EXPR_RESULT node into
// FOR loop initializers where possible.
//
// Java: Denormalize.maybeCollapseIntoForStatements(Node, Node)
func (d *Denormalize) MaybeCollapseIntoForStatements(n rhino.Node, parent rhino.Node) {
	// Is the next statement a valid FOR?
	nextSibling := n.GetNext()
	if nextSibling.IsNil() {
		return
	}

	// Only SCRIPT, BLOCK, and LABELs can have FORs that can be collapsed into.
	// LABELs are not supported here.
	if parent.IsNil() || !NodeUtil_IsStatementBlock(parent) {
		return
	}

	// Is the current node something that can be in a for loop initializer?
	if !n.IsExprResult() && !n.IsVar() {
		return
	}

	if nextSibling.IsForIn() || nextSibling.IsForOf() {
		forNode := nextSibling
		forVar := forNode.GetFirstChild()
		if forVar.IsName() &&
			n.IsVar() && n.HasOneChild() {
			name := n.GetFirstChild()
			if !name.HasChildren() &&
				forVar.GetString() == name.GetString() {
				// OK, the names match, and the var declaration does not have an
				// initializer. Move it into the loop.
				n.Detach()
				forVar.ReplaceWith(n)
				d.compiler.ReportChangeToEnclosingScope(parent)
			}
		}
	} else if nextSibling.IsVanillaFor() && nextSibling.GetFirstChild().IsEmpty() {

		// Does the current node contain an in operator?  If so, embedding
		// the expression in a for loop can cause some JavaScript parsers (such
		// as the PlayStation 3's browser based on Access's NetFront
		// browser) to fail to parse the code.
		// See bug 1778863 for details.
		if NodeUtil_Has(n, (rhino.Node).IsIn, denormalize_alwaysTrue) {
			return
		}

		// Move the current node into the FOR loop initializer.
		forNode := nextSibling
		oldInitializer := forNode.GetFirstChild()
		n.Detach()

		var newInitializer rhino.Node
		if n.IsVar() {
			newInitializer = n
		} else {
			// Extract the expression from EXPR_RESULT node.
			jutil.CheckStateMsg1(n.HasOneChild(), "%s", n)
			newInitializer = n.GetFirstChild()
			newInitializer.Detach()
		}

		oldInitializer.ReplaceWith(newInitializer)

		d.compiler.ReportChangeToEnclosingScope(forNode)
	}
}

// denormalize_alwaysTrue is Guava's Predicates.alwaysTrue().
func denormalize_alwaysTrue(rhino.Node) bool { return true }

// Java: Denormalize.maybeCollapseAssignShorthand(Node, Node)
func (d *Denormalize) MaybeCollapseAssignShorthand(n rhino.Node, parent rhino.Node) {
	if !d.IsCollapsableAssign(n) {
		return
	}
	op := n.GetLastChild()
	assignOp := d.GetAssignOpFromOp(op)
	if n.GetFirstChild().GetString() == op.GetFirstChild().GetString() {
		op.SetToken(assignOp)
		opDetached := op.Detach()
		opDetached.SetJSDocInfo(n.GetJSDocInfo())
		n.ReplaceWith(opDetached)
		d.compiler.ReportChangeToEnclosingScope(parent)
	}
}

// Java: Denormalize.maybeCollapseLogicalAssignShorthand(NodeTraversal, Node, Node)
func (d *Denormalize) MaybeCollapseLogicalAssignShorthand(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !d.IsCollapsableLogicalAssign(n) {
		return
	}
	op := n.GetLastChild()
	assignOp := d.GetAssignOpFromOp(n)
	if n.GetFirstChild().GetString() == op.GetFirstChild().GetString() {
		op.SetToken(assignOp)
		opDetached := op.Detach()
		opDetached.SetJSDocInfo(n.GetJSDocInfo())
		n.ReplaceWith(opDetached)
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT, d.compiler)
		d.compiler.ReportChangeToEnclosingScope(parent)
	}
}

// Java: Denormalize.isCollapsableAssign(Node)
func (d *Denormalize) IsCollapsableAssign(n rhino.Node) bool {
	return n.IsAssign() &&
		n.GetFirstChild().IsName() &&
		d.HasCorrespondingAssignmentOp(n.GetLastChild()) &&
		n.GetLastChild().GetFirstChild().IsName()
}

// Java: Denormalize.isCollapsableLogicalAssign(Node)
func (d *Denormalize) IsCollapsableLogicalAssign(n rhino.Node) bool {
	return d.outputFeatureSet.Has(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT) &&
		(n.IsOr() || n.IsAnd() || n.IsNullishCoalesce()) &&
		n.GetFirstChild().IsName() &&
		n.GetLastChild().IsAssign() &&
		n.GetLastChild().GetFirstChild().IsName()
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
