// Port of com/google/javascript/jscomp/Normalize.java @ closure-compiler v20260415

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Normalize: the goal with this pass is to simplify the other passes, by
// making less complex statements.
//
// Starting with statements like: var a = 0, b = foo();
//
// Which become: var a = 0; var b = foo();
//
// The key here is only to break down things that help the other passes and
// can be put back together in a form that is at least as small when all is
// said and done.
//
// This pass currently does the following:
//
//  1. Simplifies the AST by splitting var/let/const statements, moving
//     initializers out of for loops, and converting whiles to fors.
//  2. Makes all variable names globally unique (extern or otherwise) so that
//     no value is ever shadowed (note: "arguments" may require special
//     handling).
//  3. Removes duplicate variable declarations.
//  4. Marks constants with the IS_CONSTANT_NAME annotation.
//  5. Rewrite body of arrow function as a block.
//  6. Take var statements out from for-loop initializer. This: for(var a =
//     0;a<0;a++) {} becomes: var a = 0; for(a;a<0;a++) {}
type Normalize struct {
	compiler       AbstractCompiler
	astFactory     *AstFactory
	assertOnChange bool
}

// Java: Normalize(Builder)
func NewNormalize(builder *Normalize_Builder) *Normalize {
	n := &Normalize{}
	n.compiler = builder.compiler
	n.assertOnChange = builder.assertOnChange
	n.astFactory = builder.compiler.CreateAstFactory()
	return n
}

// Java: Normalize.createNormalizeForOptimizations(AbstractCompiler)
func Normalize_CreateNormalizeForOptimizations(compiler AbstractCompiler) *Normalize {
	// The default option values are the right ones for optimizations
	return Normalize_Builder_(compiler).Build()
}

// Go-only: language-adapter — the trailing underscore keeps the mechanical
// name Normalize_Builder for the nested Builder type.
//
// Java: Normalize.builder(AbstractCompiler)
func Normalize_Builder_(compiler AbstractCompiler) *Normalize_Builder {
	return NewNormalize_Builder(compiler)
}

// Normalize_Builder configures and builds a Normalize object.
type Normalize_Builder struct {
	compiler       AbstractCompiler
	assertOnChange bool
}

// Java: Normalize.Builder(AbstractCompiler)
func NewNormalize_Builder(compiler AbstractCompiler) *Normalize_Builder {
	b := &Normalize_Builder{}
	b.assertOnChange = false
	b.compiler = compiler
	return b
}

// AssertOnChange: if the Normalize pass finds work to do, it will throw an
// exception.
//
// This is intended for use in validating that an AST is already normalized.
//
// This option is false by default.
//
// Java: Normalize.Builder.assertOnChange(boolean)
func (b *Normalize_Builder) AssertOnChange(assertOnChange bool) *Normalize_Builder {
	b.assertOnChange = assertOnChange
	return b
}

// Java: Normalize.Builder.build()
func (b *Normalize_Builder) Build() *Normalize {
	return NewNormalize(b)
}

// Java: Normalize.reportCodeChange(String, Node)
func (n *Normalize) ReportCodeChange(changeDescription string, node *rhino.Node) {
	if n.assertOnChange {
		panic(jutil.NewIllegalStateExceptionMsg("Normalize constraints violated:\n" + changeDescription))
	}
	n.compiler.ReportChangeToEnclosingScope(node)
}

// Normalize_IsFunctionExpressionNameNode: is this a name node of a function
// expression?
//
// Java: Normalize.isFunctionExpressionNameNode(Node)
func Normalize_IsFunctionExpressionNameNode(n *rhino.Node) bool {
	if n == nil || !n.IsName() {
		return false
	}
	parent := n.GetParent()
	if parent == nil {
		return false
	}
	return NodeUtil_IsFunctionExpression(parent) && n.IsFirstChildOf(parent)
}

// Java: Normalize.process(Node, Node)
func (n *Normalize) Process(externs *rhino.Node, root *rhino.Node) {
	renamer := MakeDeclaredNamesUnique_Builder_().WithAssertOnChange(n.assertOnChange).Build()
	NodeTraversal_TraverseRoots(
		n.compiler, NewNormalize_NormalizeStatements(n.compiler, n.assertOnChange, renamer), externs, root)

	NodeTraversal_Builder_().
		SetCompiler(n.compiler).
		SetCallback(NewNormalize_PropagateConstantPropertyOverVars(n.compiler, n.assertOnChange)).
		SetScopeCreator(NewSyntacticScopeCreatorWithRedeclarationHandler(n.compiler, NewNormalize_DuplicateDeclarationHandler(n))).
		TraverseRoots(externs, root)

	if !n.compiler.GetLifeCycleStage().IsNormalized() {
		n.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	}
}

// Normalize_PropagateConstantPropertyOverVars propagates constant
// annotations and IS_CONSTANT_NAME property over the Var graph.
//
// Also invokes t.getScope() on every scope, for use with the
// DuplicateDeclarationHandler.
type Normalize_PropagateConstantPropertyOverVars struct {
	compiler       AbstractCompiler
	assertOnChange bool
}

// Java: Normalize.PropagateConstantPropertyOverVars(AbstractCompiler, boolean)
func NewNormalize_PropagateConstantPropertyOverVars(
	compiler AbstractCompiler, forbidChanges bool) *Normalize_PropagateConstantPropertyOverVars {
	p := &Normalize_PropagateConstantPropertyOverVars{}
	p.compiler = compiler
	p.assertOnChange = forbidChanges
	return p
}

// Java: Normalize.PropagateConstantPropertyOverVars.enterScope(NodeTraversal)
func (p *Normalize_PropagateConstantPropertyOverVars) EnterScope(t *NodeTraversal) {
	// Cause the scope to be created, which will cause duplicate
	// to be found.
	t.GetScope()
}

// Java: Normalize.PropagateConstantPropertyOverVars.exitScope(NodeTraversal)
func (p *Normalize_PropagateConstantPropertyOverVars) ExitScope(t *NodeTraversal) {
	// Nothing to do.
}

// Java: Normalize.PropagateConstantPropertyOverVars.shouldTraverse(NodeTraversal, Node, Node)
func (p *Normalize_PropagateConstantPropertyOverVars) ShouldTraverse(
	t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return true
}

// Java: Normalize.PropagateConstantPropertyOverVars.visit(NodeTraversal, Node, Node)
func (p *Normalize_PropagateConstantPropertyOverVars) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// Note: Constant properties annotations are not propagated.
	if !n.IsName() || n.GetString() == "" {
		return
	}

	// Find the JSDocInfo for a top-level variable
	var_ := t.GetScope().GetVar(n.GetString())
	var info *rhino.JSDocInfo
	if var_ != nil {
		info = var_.GetJSDocInfo()
	}

	shouldBeConstant :=
		(var_ != nil && (var_.IsConst() || Normalize_IsFunctionExpressionNameNode(var_.GetNode()))) ||
			(info != nil && info.IsConstant()) ||
			NodeUtil_IsConstantByConvention(p.compiler.GetCodingConvention(), n)
	isMarkedConstant := n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)
	if shouldBeConstant && !isMarkedConstant {
		if p.assertOnChange {
			name := n.GetString()
			panic(jutil.NewIllegalStateExceptionMsg(
				"Should be const but not marked as const.\n" +
					"  name: " +
					name +
					"\n" +
					"  parent:" +
					n.GetParent().ToStringTree()))
		}
		n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
}

// Normalize_NormalizeStatements simplifies the AST:
//   - VAR declarations split, so they represent exactly one child
//     declaration.
//   - WHILEs are converted to FORs.
//   - FOR loop are initializers are moved out of the FOR structure.
//   - LABEL node of children other than LABEL, BLOCK, WHILE, FOR, or DO are
//     moved into a block.
//   - Add constant annotations based on coding convention.
type Normalize_NormalizeStatements struct {
	compiler                                AbstractCompiler
	astFactory                              *AstFactory
	assertOnChange                          bool
	makeDeclaredNamesUnique                 *MakeDeclaredNamesUnique
	rewriteLogicalAssignmentOperatorsHelper *RewriteLogicalAssignmentOperatorsHelper
}

// NewNormalize_NormalizeStatements: prefer using Normalize as a compiler
// pass instead of using this class directly. makeDeclaredNamesUnique may be
// nil.
//
// Java: Normalize.NormalizeStatements(AbstractCompiler, boolean, MakeDeclaredNamesUnique)
func NewNormalize_NormalizeStatements(
	compiler AbstractCompiler,
	assertOnChange bool,
	makeDeclaredNamesUnique *MakeDeclaredNamesUnique) *Normalize_NormalizeStatements {
	s := &Normalize_NormalizeStatements{}
	s.compiler = compiler
	s.assertOnChange = assertOnChange
	s.makeDeclaredNamesUnique = makeDeclaredNamesUnique
	s.astFactory = compiler.CreateAstFactory()
	s.rewriteLogicalAssignmentOperatorsHelper =
		NewRewriteLogicalAssignmentOperatorsHelper(
			compiler, s.astFactory, compiler.GetUniqueIdSupplier())
	return s
}

// Java: Normalize.NormalizeStatements.reportCodeChange(String, Node)
func (s *Normalize_NormalizeStatements) ReportCodeChange(changeDescription string, n *rhino.Node) {
	if s.assertOnChange {
		panic(jutil.NewIllegalStateExceptionMsg("Normalize constraints violated:\n" + changeDescription))
	}
	s.compiler.ReportChangeToEnclosingScope(n)
}

// Java: Normalize.NormalizeStatements.enterScope(NodeTraversal)
func (s *Normalize_NormalizeStatements) EnterScope(t *NodeTraversal) {
	if s.makeDeclaredNamesUnique != nil {
		s.makeDeclaredNamesUnique.EnterScope(t)
	}
}

// Java: Normalize.NormalizeStatements.exitScope(NodeTraversal)
func (s *Normalize_NormalizeStatements) ExitScope(t *NodeTraversal) {
	if s.makeDeclaredNamesUnique != nil {
		s.makeDeclaredNamesUnique.ExitScope(t)
	}
}

// Java: Normalize.NormalizeStatements.shouldTraverse(NodeTraversal, Node, Node)
func (s *Normalize_NormalizeStatements) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if s.makeDeclaredNamesUnique != nil {
		_ = s.makeDeclaredNamesUnique.ShouldTraverse(t, n, parent)
	}

	oldParent := n.GetParent()
	s.DoStatementNormalizations(n)
	jutil.CheckStateMsg(
		n.GetParent() == oldParent,
		"n's parent changed during normalization. It needs to stay the same so that we properly"+
			" visit any added nodes.")
	return true
}

// Java: Normalize.NormalizeStatements.visit(NodeTraversal, Node, Node)
func (s *Normalize_NormalizeStatements) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if s.makeDeclaredNamesUnique != nil {
		s.makeDeclaredNamesUnique.Visit(t, n, parent)
	}

	switch n.GetToken() {
	case rhino.Token_WHILE:
		expr := n.GetFirstChild()
		n.SetToken(rhino.Token_FOR)
		empty := rhino.IR_Empty()
		empty.SrcrefIfMissing(n)
		empty.InsertBefore(expr)
		empty.CloneNode().InsertAfter(expr)
		s.ReportCodeChange("WHILE node", n)
	case rhino.Token_FUNCTION:
		if s.VisitFunction(n) {
			s.ReportCodeChange("Function declaration", n)
		}
	case rhino.Token_ARRAYLIT, rhino.Token_CALL, rhino.Token_PARAM_LIST, rhino.Token_NEW,
		rhino.Token_OBJECTLIT, rhino.Token_OPTCHAIN_CALL:
		n.SetTrailingComma(false)
	case rhino.Token_NAME:
		s.AnnotateConstantsByConvention(n)
		s.AnnotateFunctionExpressionNameAsConstant(n)
	case rhino.Token_DESTRUCTURING_LHS:
		s.NormalizeDestructuringLhs(n)
	case rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND, rhino.Token_ASSIGN_COALESCE:
		s.rewriteLogicalAssignmentOperatorsHelper.VisitLogicalAssignmentOperator(t, n)
	default:
	}
}

// AnnotateConstantsByConvention marks names and properties that are
// constants by convention.
//
// Java: Normalize.NormalizeStatements.annotateConstantsByConvention(Node)
func (s *Normalize_NormalizeStatements) AnnotateConstantsByConvention(n *rhino.Node) {
	jutil.CheckState(n.IsName())

	// Need to check that variables have not been renamed, to determine whether
	// coding conventions still apply.
	if s.compiler.GetLifeCycleStage().IsNormalizedObfuscated() {
		return
	}

	if !n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) &&
		NodeUtil_IsConstantByConvention(s.compiler.GetCodingConvention(), n) {
		jutil.CheckStateMsg1(!s.assertOnChange, "Not marked as constant when it should be: %s", n)
		n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
}

// AnnotateFunctionExpressionNameAsConstant annotates function names on
// function expressions to be unconditionally constant
//
// Java: Normalize.NormalizeStatements.annotateFunctionExpressionNameAsConstant(Node)
func (s *Normalize_NormalizeStatements) AnnotateFunctionExpressionNameAsConstant(n *rhino.Node) {
	jutil.CheckStateMsg1(n.IsName(), "Expected NAME node but got %s", n.GetToken().String())
	if Normalize_IsFunctionExpressionNameNode(n) {
		n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
}

// SplitExportDeclaration splits ES6 export combined with a variable or
// function declaration.
//
// Java: Normalize.NormalizeStatements.splitExportDeclaration(Node)
func (s *Normalize_NormalizeStatements) SplitExportDeclaration(n *rhino.Node) {
	if n.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
		return
	}
	var destructuringLhsNodes []*rhino.Node
	c := n.GetFirstChild()
	if NodeUtil_IsDeclaration(c) {
		c.Detach()

		exportSpecs := rhino.NewNode(rhino.Token_EXPORT_SPECS).Srcref(n)
		n.AddChildToFront(exportSpecs)
		if c.IsClass() || c.IsFunction() {
			name := c.GetFirstChild()
			c.InsertBefore(n)
			s.AddNameNodeToExportSpecs(exportSpecs, name)
		} else {
			NodeUtil_VisitLhsNodesInNode(c, func(name *rhino.Node) { s.AddNameNodeToExportSpecs(exportSpecs, name) })
			// Split up var declarations onto separate lines.
			for child := c.GetFirstChild(); child != nil; {
				next := child.GetNext()
				child.Detach()
				newDeclaration := rhino.NewNodeWithChild(c.GetToken(), child).Srcref(n)
				newDeclaration.InsertBefore(n)
				if child.IsDestructuringLhs() {
					// `export {a, b} = ...` changed to `var {a,b} = ...; export {a, b};`
					// Hence we must normalize the destructuring declaration.
					destructuringLhsNodes = append(destructuringLhsNodes, child)
				}
				child = next
			}
		}

		// normalize the newly added destructuring var declarations
		for _, destructuringLhsNode := range destructuringLhsNodes {
			s.NormalizeDestructuringLhs(destructuringLhsNode)
		}

		s.ReportCodeChange("combined export and declaration", n.GetParent())
	}
}

// Java: Normalize.NormalizeStatements.addNameNodeToExportSpecs(Node, Node)
func (s *Normalize_NormalizeStatements) AddNameNodeToExportSpecs(exportSpecs *rhino.Node, name *rhino.Node) {
	exportSpec := rhino.NewNode(rhino.Token_EXPORT_SPEC).Srcref(name)
	exportSpec.AddChildToFront(name.CloneNode())
	exportSpec.AddChildToFront(name.CloneNode())
	exportSpecs.AddChildToBack(exportSpec)
}

// VisitFunction rewrites blockless arrow functions to have a block with a
// single return statement
//
// For example: (x) => x becomes (x) => { return x; }.
//
// This simplifies optimizations as they can now assume all functions have a
// BLOCK.
//
// Java: Normalize.NormalizeStatements.visitFunction(Node)
func (s *Normalize_NormalizeStatements) VisitFunction(n *rhino.Node) bool {
	jutil.CheckStateMsg1(n.IsFunction(), "%s", n)
	if n.IsFunction() && !NodeUtil_GetFunctionBody(n).IsBlock() {
		returnValue := NodeUtil_GetFunctionBody(n)
		body := rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(returnValue.Detach()))
		body.SrcrefTreeIfMissing(returnValue)
		n.AddChildToBack(body)
		s.ReportCodeChange("blockless arrow function", body)
	}
	return false
}

// DoStatementNormalizations does normalizations that introduce new siblings
// or parents.
//
// Java: Normalize.NormalizeStatements.doStatementNormalizations(Node)
func (s *Normalize_NormalizeStatements) DoStatementNormalizations(n *rhino.Node) {
	if n.IsLabel() {
		s.NormalizeLabels(n)
	}

	// Only inspect the children of SCRIPTs, BLOCKs and LABELs, as all these
	// are the only legal place for VARs and FOR statements.
	if NodeUtil_IsStatementBlock(n) || n.IsLabel() {
		// makes sure that var [a,b] destructuring from for-loop initializer is also extracted
		s.ExtractForInitializer(n, nil, nil)
	}

	// Only inspect the children of SCRIPTs, BLOCKs, as all these
	// are the only legal place for VARs.
	if NodeUtil_IsStatementBlock(n) {
		s.SplitVarDeclarations(n)
	}

	if n.IsFunction() {
		s.MoveNamedFunctions(n.GetLastChild())
	}

	if NodeUtil_IsClassStaticBlock(n) {
		s.MoveNamedFunctions(n)
	}

	if n.IsExport() {
		// Perform the split in pre-order traversal to accurately normalize destructuring
		// declarations.
		s.SplitExportDeclaration(n)
	}

	if NodeUtil_IsCompoundAssignmentOp(n) && !NodeUtil_IsLogicalAssignmentOp(n) {
		// Logical assignments should be handled in visit(), not here
		s.NormalizeAssignShorthand(n)
	}
}

// NormalizeDestructuringLhs splits a var destructuring declaration into stub
// declarations of individual LHS name nodes following by the destructuring
// pattern assignment.
//
// Java: Normalize.NormalizeStatements.normalizeDestructuringLhs(Node)
func (s *Normalize_NormalizeStatements) NormalizeDestructuringLhs(n *rhino.Node) {
	if s.IsVarDestructuringDeclaration(n) {
		s.RewriteVarDestructuringDeclaration(n)
	}
}

// IsVarDestructuringDeclaration: is this a var destructuring LHS like
// `[a,b]` in `var [a,b] = ...`
//
// Java: Normalize.NormalizeStatements.isVarDestructuringDeclaration(Node)
func (s *Normalize_NormalizeStatements) IsVarDestructuringDeclaration(n *rhino.Node) bool {
	return n.IsDestructuringLhs() && // `[a,b]`
		(n.GetFirstChild().IsArrayPattern() ||
			n.GetFirstChild().
				IsObjectPattern()) && // `var [a,b] = [5,6];` or `var {a,b} = {a: 5, b: 6};`
		n.GetParent().IsVar() // `var [a,b] = [5,6];`
}

// RewriteVarDestructuringDeclaration transforms a var destructuring array
// declarations into stub declarations of the individual lhs names within the
// array. For example:
//
//	var [a, b = 3] = ....
//
// to:
//
//	var a; var b; [a, b = 3] = ...
//
// Same for object destructuring declarations. For example:
//
//	var {a, b = 3} = ....
//
// to:
//
//	var a; var b; ({a, b = 3} = ...);
//
// Before this function is used, extractForInitializer has already run during
// pre-order traversal, so we need not worry about destructuring declarations
// in the for loop initializers (e.g. `for (var [a,b] ...)` is already
// rewritten before this). Also, export declarations are already handled in
// `splitExportDeclaration` during pre-order traversal.
//
// Java: Normalize.NormalizeStatements.rewriteVarDestructuringDeclaration(Node)
func (s *Normalize_NormalizeStatements) RewriteVarDestructuringDeclaration(destructuringLhs *rhino.Node) {
	jutil.CheckStateMsg(
		!destructuringLhs.GetGrandparent().IsExport(),
		"Export destructuring declarations should be already handled in `splitExportDeclaration`"+
			" during pre-order traversal.")
	var_ := destructuringLhs.GetParent()

	// create a stub declaration for each name in the destructuring pattern
	NodeUtil_VisitLhsNodesInNode(
		destructuringLhs,
		func(name *rhino.Node) {
			// Add a declaration outside the destructuring pattern for the given name.
			jutil.CheckStateMsg1(
				name.IsName(),
				"lhs in destructuring declaration should be a simple name. (%s)",
				name)
			newName := rhino.IR_Name(name.GetString()).Srcref(name)
			if name.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
				// if old name was a const, new name should be too
				// e.g. when rewriting `{VALUE} = ...` the `VALUE` is const by coding convention
				newName.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
			}
			newVar := rhino.IR_Var(newName).Srcref(name)
			newVar.InsertBefore(var_)
		})

	// Transform destructuring var declaration to assignment. That is, `var [a, b] = ...` to `[a,
	// b] = ...` and `var {a, b} = ...` to `({a, b} = ...);`
	destructuringPattern := destructuringLhs.RemoveFirstChild()
	jutil.CheckStateMsg(destructuringPattern.IsDestructuringPattern(), "Expected destructuring pattern.")

	rhs := destructuringLhs.RemoveFirstChild()
	assign := s.astFactory.CreateAssign(destructuringPattern, rhs)
	assign.Srcref(var_)
	expr := s.astFactory.ExprResult(assign)
	expr.Srcref(var_)

	var_.ReplaceWith(expr)
	destructuringLhs.Detach()
	s.ReportCodeChange("Var destructuring declaration rewritten", expr)
}

// NormalizeLabels limits the number of special cases where LABELs need to be
// handled. Only BLOCK and loops are allowed to be labeled. Loop labels must
// remain in place as the named continues are not allowed for labeled blocks.
//
// Java: Normalize.NormalizeStatements.normalizeLabels(Node)
func (s *Normalize_NormalizeStatements) NormalizeLabels(n *rhino.Node) {
	jutil.CheckArgument(n.IsLabel())

	last := n.GetLastChild()
	// TODO(moz): Avoid adding blocks for cases like "label: let x;"
	switch last.GetToken() {
	case rhino.Token_LABEL, rhino.Token_BLOCK, rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF,
		rhino.Token_FOR_AWAIT_OF, rhino.Token_WHILE, rhino.Token_DO:
		return
	default:
		block := rhino.IR_Block()
		block.SrcrefIfMissing(last)
		last.ReplaceWith(block)
		block.AddChildToFront(last)
		s.ReportCodeChange("LABEL normalization", n)
		return
	}
}

// ExtractForInitializer brings the initializers out of FOR loops. These need
// to be placed before any associated LABEL nodes. This needs to be done from
// the top level label first so this is called as a pre-order callback (from
// shouldTraverse).
//
// n is the node to inspect; before is the node to insert the initializer
// before (may be nil); beforeParent is the parent of the node before which
// the initializer will be inserted (may be nil).
//
// Java: Normalize.NormalizeStatements.extractForInitializer(Node, Node, Node)
func (s *Normalize_NormalizeStatements) ExtractForInitializer(n *rhino.Node, before *rhino.Node, beforeParent *rhino.Node) {
	var next *rhino.Node
	for c := n.GetFirstChild(); c != nil; c = next {
		next = c.GetNext()
		var insertBefore *rhino.Node
		if before == nil {
			insertBefore = c
		} else {
			insertBefore = before
		}
		var insertBeforeParent *rhino.Node
		if before == nil {
			insertBeforeParent = n
		} else {
			insertBeforeParent = beforeParent
		}
		switch c.GetToken() {
		case rhino.Token_LABEL:
			s.ExtractForInitializer(c, insertBefore, insertBeforeParent)
		case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
			first := c.GetFirstChild()
			if first.IsVar() {
				lhs := first.GetFirstChild()
				if lhs.IsDestructuringLhs() {
					// Transform:
					//    for (var [a, b = 3] in c) {}
					// to:
					//    var a; var b; for ([a, b = 3] in c) {}
					NodeUtil_VisitLhsNodesInNode(
						lhs,
						func(name *rhino.Node) {
							// Add a declaration outside the for loop for the given name.
							jutil.CheckStateMsg1(
								name.IsName(),
								"lhs in destructuring declaration should be a simple name. (%s)",
								name)
							newName := rhino.IR_Name(name.GetString()).Srcref(name)
							newVar := rhino.IR_Var(newName).Srcref(name)
							newVar.InsertBefore(insertBefore)
						})

					// Transform for (var [a, b]... ) to for ([a, b]...
					destructuringPattern := lhs.RemoveFirstChild()
					first.ReplaceWith(destructuringPattern)
				} else {
					// Transform:
					//    for (var a in b) {}
					// to:
					//    var a; for (a in b) {};
					newStatement := first
					name := newStatement.GetFirstChild().CloneNode()
					first.ReplaceWith(name)
					newStatement.InsertBefore(insertBefore)
				}
				s.ReportCodeChange("FOR-IN var declaration", n)
			}
		case rhino.Token_FOR:
			if !c.GetFirstChild().IsEmpty() {
				init := c.GetFirstChild()

				if init.IsLet() || init.IsConst() || init.IsClass() || init.IsFunction() {
					continue
				}

				empty := rhino.IR_Empty()
				empty.SrcrefIfMissing(c)
				init.ReplaceWith(empty)

				var newStatement *rhino.Node
				// Only VAR statements, and expressions are allowed,
				// but are handled differently.
				if init.IsVar() {
					newStatement = init
				} else {
					newStatement = NodeUtil_NewExpr(init)
				}

				newStatement.InsertBefore(insertBefore)
				s.ReportCodeChange("FOR initializer", n)
			}
		default:
		}
	}
}

// SplitVarDeclarations splits a var (or let or const) node such as:
//
// var a, b;
//
// into individual statements:
//
// var a; var b;
//
// n is the node whose children we should inspect.
//
// Java: Normalize.NormalizeStatements.splitVarDeclarations(Node)
func (s *Normalize_NormalizeStatements) SplitVarDeclarations(n *rhino.Node) {
	var next *rhino.Node
	for c := n.GetFirstChild(); c != nil; c = next {
		next = c.GetNext()
		if NodeUtil_IsNameDeclaration(c) {
			if s.assertOnChange && !c.HasChildren() {
				panic(jutil.NewIllegalStateExceptionMsg("Empty VAR node."))
			}

			for c.GetFirstChild() != c.GetLastChild() {
				name := c.GetFirstChild()
				name.Detach()
				newVar := rhino.NewNodeWithChild(c.GetToken(), name).Srcref(n)
				newVar.InsertBefore(c)
				s.ReportCodeChange("VAR with multiple children", n)
			}
		}
	}
}

// MoveNamedFunctions moves all the functions that are valid at the execution
// of the first statement of the function or static block to the beginning of
// the function definition or class static block.
//
// Java: Normalize.NormalizeStatements.moveNamedFunctions(Node)
func (s *Normalize_NormalizeStatements) MoveNamedFunctions(body *rhino.Node) {
	jutil.CheckState(body.GetParent().IsFunction() || NodeUtil_IsClassStaticBlock(body))
	var insertAfter *rhino.Node
	current := body.GetFirstChild()
	// Skip any declarations at the beginning of the function body, they
	// are already in the right place.
	for current != nil && NodeUtil_IsFunctionDeclaration(current) {
		insertAfter = current
		current = current.GetNext()
	}

	// Find any remaining declarations and move them.
	for current != nil {
		// Save off the next node as the current node maybe removed.
		next := current.GetNext()
		if NodeUtil_IsFunctionDeclaration(current) {
			// Remove the declaration from the body.
			current.Detach()

			// Read the function at the top of the function body (after any
			// previous declarations).
			insertAfter = Normalize_NormalizeStatements_AddToFront(body, current, insertAfter)
			s.ReportCodeChange(
				"Move function declaration not at top of function or class static block", body)
		}
		current = next
	}
}

// Java: Normalize.NormalizeStatements.normalizeAssignShorthand(Node)
func (s *Normalize_NormalizeStatements) NormalizeAssignShorthand(shorthand *rhino.Node) {
	if !shorthand.GetFirstChild().IsName() {
		return
	}

	jutil.CheckArgumentMsg(shorthand.HasTwoChildren(), "Expected two children for shorthand assignment.")
	lhs := shorthand.GetFirstChild()
	rhs := shorthand.GetSecondChild()
	lhs.Detach()
	rhs.Detach()

	op := NodeUtil_GetOpFromAssignmentOp(shorthand)
	opNode :=
		rhino.NewNodeWith2Children(op, lhs.CloneNode().Srcref(lhs), rhs).Srcref(shorthand).CopyTypeFrom(shorthand)

	// Change the original ASSIGN_X to an ASSIGN.
	shorthand.SetToken(rhino.Token_ASSIGN)
	shorthand.AddChildToBack(lhs)
	shorthand.AddChildToBack(opNode)

	s.ReportCodeChange("assign shorthand", shorthand)
}

// Normalize_NormalizeStatements_AddToFront: after is the child node to
// insert the newChild after, or nil if newChild should be added to the front
// of parent's child list. Returns the inserted child node.
//
// Java: Normalize.NormalizeStatements.addToFront(Node, Node, Node)
func Normalize_NormalizeStatements_AddToFront(parent *rhino.Node, newChild *rhino.Node, after *rhino.Node) *rhino.Node {
	if after == nil {
		parent.AddChildToFront(newChild)
	} else {
		newChild.InsertAfter(after)
	}
	return newChild
}

// Normalize_DuplicateDeclarationHandler is the ScopeCreator duplicate
// declaration handler.
type Normalize_DuplicateDeclarationHandler struct {
	outer *Normalize

	hasOkDuplicateDeclaration *jutil.LinkedHashSet[*Var]
}

// Java: Normalize.DuplicateDeclarationHandler()
func NewNormalize_DuplicateDeclarationHandler(outer *Normalize) *Normalize_DuplicateDeclarationHandler {
	h := &Normalize_DuplicateDeclarationHandler{}
	h.outer = outer
	h.hasOkDuplicateDeclaration = jutil.NewLinkedHashSet[*Var]()
	return h
}

// OnRedeclaration removes duplicate VAR declarations discovered during scope
// creation.
//
// Java: Normalize.DuplicateDeclarationHandler.onRedeclaration(Scope, String, Node, CompilerInput)
func (h *Normalize_DuplicateDeclarationHandler) OnRedeclaration(
	s *Scope, name string, n *rhino.Node, input *CompilerInput) {
	jutil.CheckState(n.IsName())
	parent := n.GetParent()
	v := s.GetVar(name)

	if s.IsGlobal() {
		// We allow variables to be duplicate declared if one
		// declaration appears in source and the other in externs.
		// This deals with issues where a browser built-in is declared
		// in one browser but not in another.
		if v.IsExtern() && !input.IsExtern() {
			if h.hasOkDuplicateDeclaration.Add(v) {
				return
			}
		}
	}

	if parent.IsFunction() {
		if v.GetParentNode().IsVar() {
			s.Undeclare(v)
			s.Declare(name, n, v.GetInput())
			h.ReplaceVarWithAssignment(
				v.GetNameNode(), v.GetParentNode(), v.GetParentNode().GetParent())
		}
	} else if parent.IsVar() {
		jutil.CheckState(parent.HasOneChild())
		h.ReplaceVarWithAssignment(n, parent, parent.GetParent())
	}
}

// ReplaceVarWithAssignment removes the parent VAR. There is only one case
// that need to be handled: "var a = b;" which is replaced with "a = b"
//
// Cases we don't need to handle are VARs with multiple children, which have
// already been split into separate declarations, so there is no need to
// handle that here; "for (var a;;);"/"for (var a of b)"/"for (var a in b)",
// which have been moved out of the loop; and "LABEL: var x;" which has been
// put in a BLOCK
//
// The result of this is that in each case the parent node is replaced which
// is generally dangerous in a traversal but is fine here with the scope
// creator, as the next node of interest is the parent's next sibling.
//
// Java: Normalize.DuplicateDeclarationHandler.replaceVarWithAssignment(Node, Node, Node)
func (h *Normalize_DuplicateDeclarationHandler) ReplaceVarWithAssignment(
	n *rhino.Node, parent *rhino.Node, grandparent *rhino.Node) {
	if n.HasChildren() {
		// The  *  is being initialize, preserve the new value.
		n.Detach()
		// Convert "var name = value" to "name = value"
		value := n.GetFirstChild()
		value.Detach()
		replacement := h.outer.astFactory.CreateAssign(n, value)
		replacement.SetJSDocInfo(parent.GetJSDocInfo())
		replacement.SrcrefIfMissing(parent)
		statement := NodeUtil_NewExpr(replacement)
		parent.ReplaceWith(statement)
		h.outer.ReportCodeChange("Duplicate VAR declaration", statement)
	} else {
		// It is an empty reference. Remove it.
		jutil.CheckStateMsg1(NodeUtil_IsStatementBlock(grandparent), "%s", grandparent)
		parent.Detach()
		h.outer.ReportCodeChange("Duplicate VAR declaration", grandparent)
	}
}
