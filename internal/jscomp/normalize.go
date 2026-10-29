// Port of com/google/javascript/jscomp/Normalize.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast). NormalizeStatements
// lists its work per node kind in two typed visitors (before and after the
// children are traversed) instead of switching on the token, and the rewrites
// name the children they read (fn.Body(), label.Body(), loop.Target()...)
// instead of using child positions.

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
func (n *Normalize) ReportCodeChange(changeDescription string, node rhino.Node) {
	if n.assertOnChange {
		panic(jutil.NewIllegalStateExceptionMsg("Normalize constraints violated:\n" + changeDescription))
	}
	n.compiler.ReportChangeToEnclosingScope(node)
}

// Normalize_IsFunctionExpressionNameNode: is this a name node of a function
// expression?
//
// Java: Normalize.isFunctionExpressionNameNode(Node)
func Normalize_IsFunctionExpressionNameNode(n rhino.Node) bool {
	name, ok := ast.AsName(n)
	return ok && normalizeIsFunctionExpressionName(name)
}

// normalizeIsFunctionExpressionName reports whether x is the name of a
// function expression, including the empty name of an anonymous one (which
// fn.Name() does not return: hence the first child).
func normalizeIsFunctionExpressionName(x ast.Name) bool {
	fn, ok := ast.AsFunction(x.Node().GetParent())
	return ok && x.Is(fn.Node().GetFirstChild()) && NodeUtil_IsFunctionExpression(fn.Node())
}

// Java: Normalize.process(Node, Node)
func (n *Normalize) Process(externs rhino.Node, root rhino.Node) {
	normStmts := NewNormalize_NormalizeStatements(n.compiler, n.assertOnChange, nil)
	normStmts.initFastContextualRenamer()
	defer normStmts.releaseFastContextualRenamer()
	sc1 := NewSyntacticScopeCreator(n.compiler)
	b1 := NodeTraversal_Builder{
		compiler:     n.compiler,
		callback:     normStmts,
		scopeCreator: sc1,
	}
	b1.TraverseRoots(externs, root)

	sc2 := NewSyntacticScopeCreatorWithRedeclarationHandler(n.compiler, NewNormalize_DuplicateDeclarationHandler(n))
	b2 := NodeTraversal_Builder{
		compiler:     n.compiler,
		callback:     NewNormalize_PropagateConstantPropertyOverVars(n.compiler, n.assertOnChange),
		scopeCreator: sc2,
	}
	b2.TraverseRoots(externs, root)

	if !n.compiler.GetLifeCycleStage().IsNormalized() {
		n.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	}
}

// SplitExportDeclaration splits ES6 export combined with a variable or
// function declaration.
//
// Java: Normalize.NormalizeStatements.splitExportDeclaration(Node)
func (s *Normalize_NormalizeStatements) SplitExportDeclaration(export ast.Export) {
	n := export.Node()
	if n.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
		return
	}
	var destructuringLhsNodes []ast.DestructuringLhs
	c := export.Body().Node()
	if NodeUtil_IsDeclaration(c) {
		c.Detach()

		exportSpecs := rhino.NewNode(rhino.Token_EXPORT_SPECS).Srcref(n)
		n.AddChildToFront(exportSpecs)
		if decl, ok := ast.AsNameDeclaration(c); ok {
			NodeUtil_VisitLhsNodesInNode(c, func(name rhino.Node) { s.AddNameNodeToExportSpecs(exportSpecs, name) })
			// Split up var declarations onto separate lines.
			for d, ok := decl.FirstDeclarator(); ok; d, ok = decl.FirstDeclarator() {
				d.Node().Detach()
				newDeclaration := rhino.NewNodeWithChild(decl.Token(), d.Node()).Srcref(n)
				newDeclaration.InsertBefore(n)
				if lhs, ok := d.AsDestructuringLhs(); ok {
					// `export {a, b} = ...` changed to `var {a,b} = ...; export {a, b};`
					// Hence we must normalize the destructuring declaration.
					destructuringLhsNodes = append(destructuringLhsNodes, lhs)
				}
			}
		} else {
			// A function or class declaration: its first child is its name.
			name := c.GetFirstChild()
			c.InsertBefore(n)
			s.AddNameNodeToExportSpecs(exportSpecs, name)
		}

		// normalize the newly added destructuring var declarations
		for _, destructuringLhsNode := range destructuringLhsNodes {
			s.NormalizeDestructuringLhs(destructuringLhsNode)
		}

		s.ReportCodeChange("combined export and declaration", n.GetParent())
	}
}

// Java: Normalize.NormalizeStatements.addNameNodeToExportSpecs(Node, Node)
func (s *Normalize_NormalizeStatements) AddNameNodeToExportSpecs(exportSpecs rhino.Node, name rhino.Node) {
	exportSpec := rhino.NewNode(rhino.Token_EXPORT_SPEC).Srcref(name)
	exportSpec.AddChildToFront(name.CloneNode())
	exportSpec.AddChildToFront(name.CloneNode())
	exportSpecs.AddChildToBack(exportSpec)
}

// NormalizeBlocklessFunction rewrites blockless arrow functions to have a
// block with a single return statement
//
// For example: (x) => x becomes (x) => { return x; }.
//
// This simplifies optimizations as they can now assume all functions have a
// BLOCK.
//
// Java: Normalize.NormalizeStatements.normalizeBlocklessFunction(Node)
func (s *Normalize_NormalizeStatements) NormalizeBlocklessFunction(fn ast.Function) {
	if _, ok := fn.Body().AsBlock(); ok {
		return
	}
	returnValue := fn.Body().Node()
	body := rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(returnValue.Detach()))
	body.SrcrefTreeIfMissing(returnValue)
	fn.Node().AddChildToBack(body) // in the place of the body just detached
	s.ReportCodeChange("blockless arrow function", body)
}

// DoStatementNormalizations does normalizations that introduce new siblings
// or parents: the handlers of normalizePreOrderVisitor.
//
// Java: Normalize.NormalizeStatements.doStatementNormalizations(Node)
func (s *Normalize_NormalizeStatements) DoStatementNormalizations(_ *NodeTraversal, n rhino.Node) {
	normalizePreOrderVisitor.Visit(s, n)
}

// normalizeStatements normalizes the statements of n: a root, a script, a
// module body or a block.
//
// Java: Normalize.NormalizeStatements.doStatementNormalizations(Node) (ROOT, SCRIPT, MODULE_BODY, BLOCK)
func (s *Normalize_NormalizeStatements) normalizeStatements(n rhino.Node) {
	s.ExtractForInitializer(n, rhino.NilNode, rhino.NilNode)
	s.SplitVarDeclarations(n)
}

// Java: Normalize.NormalizeStatements.doStatementNormalizations(Node) (BLOCK)
func (s *Normalize_NormalizeStatements) normalizeBlock(x ast.Block) {
	s.normalizeStatements(x.Node())
	if NodeUtil_IsClassStaticBlock(x.Node()) {
		s.MoveNamedFunctions(x)
	}
}

// Java: Normalize.NormalizeStatements.doStatementNormalizations(Node) (LABEL)
func (s *Normalize_NormalizeStatements) normalizeLabel(x ast.Label) {
	s.NormalizeLabels(x)
	s.ExtractForInitializer(x.Node(), rhino.NilNode, rhino.NilNode)
}

// Java: Normalize.NormalizeStatements.doStatementNormalizations(Node) (FUNCTION)
func (s *Normalize_NormalizeStatements) normalizeFunction(x ast.Function) {
	s.NormalizeBlocklessFunction(x)
	body, ok := x.Body().AsBlock()
	jutil.CheckState(ok) // normalizeBlocklessFunction gave it one
	s.MoveNamedFunctions(body)
}

// NormalizeDestructuringLhs splits a var destructuring declaration into stub
// declarations of individual LHS name nodes following by the destructuring
// pattern assignment.
//
// Java: Normalize.NormalizeStatements.normalizeDestructuringLhs(Node)
func (s *Normalize_NormalizeStatements) NormalizeDestructuringLhs(destructuringLhs ast.DestructuringLhs) {
	if s.IsVarDestructuringDeclaration(destructuringLhs) {
		s.RewriteVarDestructuringDeclaration(destructuringLhs)
	}
}

// IsVarDestructuringDeclaration: is this a var destructuring LHS like
// `[a,b]` in `var [a,b] = ...`
//
// (Its pattern is an array or object pattern: the grammar says so.)
//
// Java: Normalize.NormalizeStatements.isVarDestructuringDeclaration(Node)
func (s *Normalize_NormalizeStatements) IsVarDestructuringDeclaration(destructuringLhs ast.DestructuringLhs) bool {
	return destructuringLhs.Node().GetParent().IsVar() // `var [a,b] = [5,6];`
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
func (s *Normalize_NormalizeStatements) RewriteVarDestructuringDeclaration(destructuringLhs ast.DestructuringLhs) {
	n := destructuringLhs.Node()
	jutil.CheckStateMsg(
		!n.GetGrandparent().IsExport(),
		"Export destructuring declarations should be already handled in `splitExportDeclaration`"+
			" during pre-order traversal.")
	var_ := n.GetParent()

	// create a stub declaration for each name in the destructuring pattern
	NodeUtil_VisitLhsNodesInNode(
		n,
		func(name rhino.Node) {
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
	destructuringPattern := destructuringLhs.Pattern().Node()
	rhs, ok := destructuringLhs.Init() // read both before detaching either
	jutil.CheckStateMsg(ok, "Expected a value.")
	destructuringPattern.Detach()
	rhs.Node().Detach()
	assign := s.astFactory.CreateAssign(destructuringPattern, rhs.Node())
	assign.Srcref(var_)
	expr := s.astFactory.ExprResult(assign)
	expr.Srcref(var_)

	var_.ReplaceWith(expr)
	n.Detach()
	s.ReportCodeChange("Var destructuring declaration rewritten", expr)
}

// NormalizeLabels limits the number of special cases where LABELs need to be
// handled. Only BLOCK and loops are allowed to be labeled. Loop labels must
// remain in place as the named continues are not allowed for labeled blocks.
//
// Java: Normalize.NormalizeStatements.normalizeLabels(Node)
func (s *Normalize_NormalizeStatements) NormalizeLabels(label ast.Label) {
	body := label.Body()
	// TODO(moz): Avoid adding blocks for cases like "label: let x;"
	_, isLabel := body.AsLabel()
	_, isBlock := body.AsBlock()
	_, isLoop := body.AsLoop()
	if isLabel || isBlock || isLoop {
		return
	}
	block := ast.NewBlock()
	block.Node().SrcrefIfMissing(body.Node())
	label.SetBody(block.AsStatement())
	block.PrependStatement(body)
	s.ReportCodeChange("LABEL normalization", label.Node())
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
func (s *Normalize_NormalizeStatements) ExtractForInitializer(n rhino.Node, before rhino.Node, beforeParent rhino.Node) {
	var next rhino.Node
	for c := n.GetFirstChild(); !c.IsNil(); c = next {
		next = c.GetNext()
		insertBefore := before
		insertBeforeParent := beforeParent
		if insertBefore.IsNil() {
			insertBefore = c
			insertBeforeParent = n
		}
		if _, ok := ast.AsLabel(c); ok {
			s.ExtractForInitializer(c, insertBefore, insertBeforeParent)
		} else if loop, ok := ast.AsEnhancedFor(c); ok {
			decl, ok := loop.Target().AsNameDeclaration()
			if !ok || decl.Token() != rhino.Token_VAR {
				continue
			}
			lhs, ok := decl.FirstDeclarator()
			jutil.CheckState(ok)
			if destructuringLhs, ok := lhs.AsDestructuringLhs(); ok {
				// Transform:
				//    for (var [a, b = 3] in c) {}
				// to:
				//    var a; var b; for ([a, b = 3] in c) {}
				NodeUtil_VisitLhsNodesInNode(
					destructuringLhs.Node(),
					func(name rhino.Node) {
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
				destructuringPattern := destructuringLhs.Pattern().Node()
				destructuringPattern.Detach()
				decl.Node().ReplaceWith(destructuringPattern)
			} else {
				// Transform:
				//    for (var a in b) {}
				// to:
				//    var a; for (a in b) {};
				name := lhs.Node().CloneNode()
				decl.Node().ReplaceWith(name)
				decl.Node().InsertBefore(insertBefore)
			}
			s.ReportCodeChange("FOR-IN var declaration", n)
		} else if forLoop, ok := ast.AsFor(c); ok {
			forInit, ok := forLoop.Init()
			if !ok {
				continue
			}
			init := forInit.Node()
			if init.IsLet() || init.IsConst() || init.IsClass() || init.IsFunction() {
				continue
			}

			empty := rhino.IR_Empty()
			empty.SrcrefIfMissing(c)
			init.ReplaceWith(empty)

			var newStatement rhino.Node
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
func (s *Normalize_NormalizeStatements) SplitVarDeclarations(n rhino.Node) {
	var next rhino.Node
	for c := n.GetFirstChild(); !c.IsNil(); c = next {
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
func (s *Normalize_NormalizeStatements) MoveNamedFunctions(body ast.Block) {
	n := body.Node()
	jutil.CheckState(n.GetParent().IsFunction() || NodeUtil_IsClassStaticBlock(n))
	var insertAfter rhino.Node
	current := n.GetFirstChild()
	// Skip any declarations at the beginning of the function body, they
	// are already in the right place.
	for !current.IsNil() && NodeUtil_IsFunctionDeclaration(current) {
		insertAfter = current
		current = current.GetNext()
	}

	// Find any remaining declarations and move them.
	for !current.IsNil() {
		// Save off the next node as the current node maybe removed.
		next := current.GetNext()
		if NodeUtil_IsFunctionDeclaration(current) {
			// Remove the declaration from the body.
			current.Detach()

			// Read the function at the top of the function body (after any
			// previous declarations).
			insertAfter = Normalize_NormalizeStatements_AddToFront(n, current, insertAfter)
			s.ReportCodeChange(
				"Move function declaration not at top of function or class static block", n)
		}
		current = next
	}
}

// Java: Normalize.NormalizeStatements.normalizeAssignShorthand(Node)
func (s *Normalize_NormalizeStatements) NormalizeAssignShorthand(shorthand ast.AssignOp) {
	target, ok := shorthand.Target().AsName()
	if !ok {
		return
	}

	n := shorthand.Node()
	lhs := target.Node()
	rhs := shorthand.Value().Node() // read before detaching the target
	lhs.Detach()
	rhs.Detach()

	op := NodeUtil_GetOpFromAssignmentOp(n)
	opNode :=
		rhino.NewNodeWith2Children(op, lhs.CloneNode().Srcref(lhs), rhs).Srcref(n).CopyTypeFrom(n)

	// Change the original ASSIGN_X to an ASSIGN.
	n.SetToken(rhino.Token_ASSIGN)
	lhs.CopyTypeFrom(n)
	n.AddChildToBack(lhs)
	n.AddChildToBack(opNode)

	s.ReportCodeChange("assign shorthand", n)
}

// Normalize_NormalizeStatements_AddToFront: after is the child node to
// insert the newChild after, or nil if newChild should be added to the front
// of parent's child list. Returns the inserted child node.
//
// Java: Normalize.NormalizeStatements.addToFront(Node, Node, Node)
func Normalize_NormalizeStatements_AddToFront(parent rhino.Node, newChild rhino.Node, after rhino.Node) rhino.Node {
	if after.IsNil() {
		parent.AddChildToFront(newChild)
	} else {
		newChild.InsertAfter(after)
	}
	return newChild
}
