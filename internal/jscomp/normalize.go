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
	"strconv"

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

func (*Normalize_PropagateConstantPropertyOverVars) onlyNameLeaves()                {}
func (*Normalize_PropagateConstantPropertyOverVars) nodeTraversalNoShouldTraverse() {}
func (*ValidityCheck_VerifyConstants) onlyNameLeaves()                              {}
func (*ValidityCheck_VerifyConstants) nodeTraversalNoShouldTraverse()               {}

// Normalize_PropagateConstantPropertyOverVars propagates constant
// annotations and IS_CONSTANT_NAME property over the Var graph.
//
// Also invokes t.getScope() on every scope, for use with the
// DuplicateDeclarationHandler.
type Normalize_PropagateConstantPropertyOverVars struct {
	compiler       AbstractCompiler
	convention     CodingConvention
	assertOnChange bool
}

// Java: Normalize.PropagateConstantPropertyOverVars(AbstractCompiler, boolean)
func NewNormalize_PropagateConstantPropertyOverVars(
	compiler AbstractCompiler, forbidChanges bool) *Normalize_PropagateConstantPropertyOverVars {
	p := &Normalize_PropagateConstantPropertyOverVars{
		compiler:       compiler,
		convention:     compiler.GetCodingConvention(),
		assertOnChange: forbidChanges,
	}
	return p
}

// Java: Normalize.PropagateConstantPropertyOverVars.enterScope(NodeTraversal)
func (p *Normalize_PropagateConstantPropertyOverVars) EnterScope(t *NodeTraversal) {
	// Cause the scope to be created, which will cause duplicate
	// to be found.
	if normalizeScopeMayHaveDeclarations(t.GetScopeRoot()) {
		t.GetScope()
	}
}

// Java: Normalize.PropagateConstantPropertyOverVars.exitScope(NodeTraversal)
func (p *Normalize_PropagateConstantPropertyOverVars) ExitScope(t *NodeTraversal) {
	// Nothing to do.
}

// Java: Normalize.PropagateConstantPropertyOverVars.shouldTraverse(NodeTraversal, Node, Node)
func (p *Normalize_PropagateConstantPropertyOverVars) ShouldTraverse(
	t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

// Java: Normalize.PropagateConstantPropertyOverVars.visit(NodeTraversal, Node, Node)
func (p *Normalize_PropagateConstantPropertyOverVars) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Note: Constant properties annotations are not propagated.
	if !n.IsName() || n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
		return
	}
	id := n.GetStringID()
	if id <= rhino.RhinoStringPool_EmptyID {
		return
	}

	// Find the JSDocInfo for a top-level variable
	var_ := t.LookupVarByID(id)
	var shouldBeConstant bool
	if var_ != nil {
		switch var_.nameHash {
		case 1:
			shouldBeConstant = false
		case 2:
			shouldBeConstant = true
		default:
			if var_.IsConst() || Normalize_IsFunctionExpressionNameNode(var_.GetNode()) {
				shouldBeConstant = true
			} else if info := var_.GetJSDocInfo(); info != nil && info.IsConstant() {
				shouldBeConstant = true
			} else if p.convention.IsConstant(n.GetString()) {
				shouldBeConstant = true
			}
			if shouldBeConstant {
				var_.nameHash = 2
			} else {
				var_.nameHash = 1
			}
		}
	} else {
		shouldBeConstant = p.convention.IsConstant(n.GetString())
	}
	if shouldBeConstant {
		if p.assertOnChange {
			panic(jutil.NewIllegalStateExceptionMsg(
				"Should be const but not marked as const.\n" +
					"  name: " +
					n.GetString() +
					"\n" +
					"  parent:" +
					n.GetParent().ToStringTree()))
		}
		n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
}

type normalizeRenameUndo struct {
	name    string
	oldRepl string
	hadOld  bool
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
	convention                              CodingConvention
	astFactory                              *AstFactory
	assertOnChange                          bool
	isNormalizedObfuscated                  bool
	makeDeclaredNamesUnique                 *MakeDeclaredNamesUnique
	rewriteLogicalAssignmentOperatorsHelper *RewriteLogicalAssignmentOperatorsHelper

	useFastContextualRenamer bool
	hasJscompVar             bool
	nameUsage                map[string]int
	activeReplacements       map[string]string
	undoStack                []normalizeRenameUndo
	scopeUndoOffsets         []int

	// Go-only: the work per node kind, before its children are traversed
	// (doStatementNormalizations) and after (visit); see initVisitors.
	preOrder  ast.Visitor[*NodeTraversal]
	postOrder ast.Visitor[*NodeTraversal]
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
	s := &Normalize_NormalizeStatements{
		compiler:                compiler,
		convention:              compiler.GetCodingConvention(),
		assertOnChange:          assertOnChange,
		isNormalizedObfuscated:  compiler.GetLifeCycleStage().IsNormalizedObfuscated(),
		makeDeclaredNamesUnique: makeDeclaredNamesUnique,
		astFactory:              compiler.CreateAstFactory(),
	}
	s.rewriteLogicalAssignmentOperatorsHelper =
		NewRewriteLogicalAssignmentOperatorsHelper(
			compiler, s.astFactory, compiler.GetUniqueIdSupplier())
	s.initVisitors()
	return s
}

// initVisitors lists what the pass does to each kind of node: before its
// children are traversed, to add siblings or parents that the traversal then
// visits (doStatementNormalizations, called from shouldTraverse), and after
// (visit).
func (s *Normalize_NormalizeStatements) initVisitors() {
	// Java: Normalize.NormalizeStatements.doStatementNormalizations(Node)
	s.preOrder = ast.Visitor[*NodeTraversal]{
		Root:       func(_ *NodeTraversal, x ast.Root) { s.normalizeStatements(x.Node()) },
		Script:     func(_ *NodeTraversal, x ast.Script) { s.normalizeStatements(x.Node()) },
		ModuleBody: func(_ *NodeTraversal, x ast.ModuleBody) { s.normalizeStatements(x.Node()) },
		Block:      s.normalizeBlock,
		Label:      s.normalizeLabel,
		Function:   s.normalizeFunction,
		// Perform the split in pre-order traversal to accurately normalize destructuring
		// declarations.
		Export: func(_ *NodeTraversal, x ast.Export) { s.SplitExportDeclaration(x) },
		AssignOp: func(_ *NodeTraversal, x ast.AssignOp) {
			// Logical assignments should be handled in visit(), not here
			if !NodeUtil_IsLogicalAssignmentOp(x.Node()) {
				s.NormalizeAssignShorthand(x)
			}
		},
	}

	// Java: Normalize.NormalizeStatements.visit(NodeTraversal, Node, Node)
	s.postOrder = ast.Visitor[*NodeTraversal]{
		Name:             s.visitName,
		ImportStar:       s.visitImportStar,
		While:            s.visitWhile,
		DestructuringLhs: func(_ *NodeTraversal, x ast.DestructuringLhs) { s.NormalizeDestructuringLhs(x) },
		AssignOp: func(t *NodeTraversal, x ast.AssignOp) {
			if NodeUtil_IsLogicalAssignmentOp(x.Node()) {
				s.rewriteLogicalAssignmentOperatorsHelper.VisitLogicalAssignmentOperator(t, x.Node())
			}
		},
		// The kinds that can have a trailing comma.
		ArrayLit:     func(_ *NodeTraversal, x ast.ArrayLit) { x.Node().SetTrailingComma(false) },
		Call:         func(_ *NodeTraversal, x ast.Call) { x.Node().SetTrailingComma(false) },
		New:          func(_ *NodeTraversal, x ast.New) { x.Node().SetTrailingComma(false) },
		ObjectLit:    func(_ *NodeTraversal, x ast.ObjectLit) { x.Node().SetTrailingComma(false) },
		OptChainCall: func(_ *NodeTraversal, x ast.OptChainCall) { x.Node().SetTrailingComma(false) },
		ParamList:    func(_ *NodeTraversal, x ast.ParamList) { x.Node().SetTrailingComma(false) },
	}
}

func (s *Normalize_NormalizeStatements) initFastContextualRenamer() {
	s.useFastContextualRenamer = true
	s.nameUsage = make(map[string]int)
	s.activeReplacements = make(map[string]string)
}

func (s *Normalize_NormalizeStatements) releaseFastContextualRenamer() {
	s.nameUsage = nil
	s.activeReplacements = nil
	s.undoStack = nil
	s.scopeUndoOffsets = nil
}

func normalizeHasJscompSeparator(s string) bool {
	const sep = MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR
	if len(s) < len(sep)+2 {
		return false
	}
	last := s[len(s)-1]
	if last < '0' || last > '9' {
		return false
	}
	i := len(s) - 1
	for i > 0 && s[i-1] >= '0' && s[i-1] <= '9' {
		i--
	}
	if s[i] == '0' || i < len(sep)+1 {
		return false
	}
	return s[i-len(sep):i] == sep
}

func (s *Normalize_NormalizeStatements) formatUniqueName(name string, id int) string {
	return name + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR + strconv.Itoa(id)
}

func (s *Normalize_NormalizeStatements) enableJscompVarMode() {
	if s.hasJscompVar {
		return
	}
	s.hasJscompVar = true
	if len(s.nameUsage) == 0 {
		return
	}
	type pending struct {
		name  string
		count int
	}
	var toBackfill []pending
	for origName, count := range s.nameUsage {
		if count > 1 {
			toBackfill = append(toBackfill, pending{name: origName, count: count})
		}
	}
	for _, p := range toBackfill {
		for i := 1; i < p.count; i++ {
			u := s.formatUniqueName(p.name, i)
			if _, exists := s.nameUsage[u]; !exists {
				s.nameUsage[u] = 1
			}
		}
	}
}

// normalizeScopeMayHaveDeclarations reports whether the scope rooted at n
// may declare names; fastEnterScope does not create the others.
func normalizeScopeMayHaveDeclarations(n rhino.Node) bool {
	if fn, ok := ast.AsFunction(n); ok {
		// Its parameters, and the name of a function expression.
		if fn.Params().HasParams() {
			return true
		}
		_, named := fn.Name()
		return named && NodeUtil_IsFunctionExpression(n)
	}
	if _, ok := ast.AsBlock(n); ok {
		if NodeUtil_IsFunctionBlock(n) || NodeUtil_IsClassStaticBlock(n) {
			return !syntacticFunctionBlockHasNoDecls(n)
		}
		// (The children are statements, or the catch of a try.)
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			switch c.GetToken() {
			case rhino.Token_LET, rhino.Token_CONST, rhino.Token_FUNCTION, rhino.Token_CLASS,
				rhino.Token_CATCH, rhino.Token_EXPORT, rhino.Token_LABEL, rhino.Token_BLOCK:
				return true
			}
		}
		return false
	}
	if cls, ok := ast.AsClass(n); ok {
		// The name of a class expression.
		_, named := cls.Name()
		return named && NodeUtil_IsClassExpression(n)
	}
	if loop, ok := ast.AsEnhancedFor(n); ok {
		return normalizeMayDeclareInLoopScope(loop.Target().Node())
	}
	if f, ok := ast.AsFor(n); ok {
		init, ok := f.Init()
		return ok && normalizeMayDeclareInLoopScope(init.Node())
	}
	switch n.GetToken() {
	case rhino.Token_COMPUTED_FIELD_DEF, rhino.Token_MEMBER_FIELD_DEF:
		return false
	}
	return true
}

// normalizeMayDeclareInLoopScope reports whether n, the init of a for or the
// target of a for-in/of, may declare names in the scope of the loop.
func normalizeMayDeclareInLoopScope(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CLASS, rhino.Token_FUNCTION:
		return true
	}
	return false
}

func (s *Normalize_NormalizeStatements) fastEnterScope(t *NodeTraversal) {
	declarationRoot := t.GetScopeRoot()
	if len(s.scopeUndoOffsets) == 0 {
		jutil.CheckState(!declarationRoot.IsFunction())
		s.scopeUndoOffsets = append(s.scopeUndoOffsets, 0)
		scope := t.GetScope()
		if s.nameUsage == nil {
			s.nameUsage = make(map[string]int, max(64, len(scope.varList)+len(scope.varList)/2))
		}
		for _, v := range scope.varList {
			if v == nil {
				continue
			}
			name := v.GetName()
			if name != MakeDeclaredNamesUnique_ARGUMENTS {
				if !s.hasJscompVar && normalizeHasJscompSeparator(name) {
					s.enableJscompVarMode()
				}
				s.nameUsage[name] = 1
			}
		}
		return
	}

	s.scopeUndoOffsets = append(s.scopeUndoOffsets, len(s.undoStack))
	if !normalizeScopeMayHaveDeclarations(declarationRoot) {
		return
	}
	for _, v := range t.GetScope().varList {
		if v == nil {
			continue
		}
		name := v.GetName()
		if name == MakeDeclaredNamesUnique_ARGUMENTS {
			continue
		}
		id := s.nameUsage[name]
		if id == 0 && !s.hasJscompVar && normalizeHasJscompSeparator(name) {
			s.enableJscompVarMode()
			id = s.nameUsage[name]
		}
		s.nameUsage[name] = id + 1
		if id != 0 {
			unique := s.formatUniqueName(name, id)
			if s.hasJscompVar {
				for {
					uid, exists := s.nameUsage[unique]
					if !exists || uid == 0 {
						break
					}
					s.nameUsage[unique] = uid + 1
					unique = s.formatUniqueName(name, uid)
				}
				if _, exists := s.nameUsage[unique]; !exists {
					s.nameUsage[unique] = 1
				}
			}
			oldRepl, hadOld := s.activeReplacements[name]
			s.undoStack = append(s.undoStack, normalizeRenameUndo{
				name:    name,
				oldRepl: oldRepl,
				hadOld:  hadOld,
			})
			s.activeReplacements[name] = unique
		}
	}
}

func (s *Normalize_NormalizeStatements) fastExitScope(t *NodeTraversal) {
	if !t.InGlobalScope() {
		offset := s.scopeUndoOffsets[len(s.scopeUndoOffsets)-1]
		s.scopeUndoOffsets = s.scopeUndoOffsets[:len(s.scopeUndoOffsets)-1]
		for i := len(s.undoStack) - 1; i >= offset; i-- {
			u := s.undoStack[i]
			if u.hadOld {
				s.activeReplacements[u.name] = u.oldRepl
			} else {
				delete(s.activeReplacements, u.name)
			}
		}
		s.undoStack = s.undoStack[:offset]
	}
}

func (s *Normalize_NormalizeStatements) fastVisitNameOrImportStar(t *NodeTraversal, n rhino.Node) {
	if len(s.activeReplacements) == 0 {
		return
	}
	if n.IsName() && NodeUtil_IsNonlocalModuleExportName(n) {
		return
	}
	newName, ok := s.activeReplacements[n.GetString()]
	if !ok {
		return
	}
	if s.assertOnChange {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected renaming in MakeDeclaredNamesUnique. new name: " + newName + " for " + n.String()))
	}
	n.SetString(newName)
	t.ReportCodeChange()
	if parent := n.GetParent(); parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
		s.compiler.ReportChangeToEnclosingScope(parent)
	}
}

// Java: Normalize.NormalizeStatements.reportCodeChange(String, Node)
func (s *Normalize_NormalizeStatements) ReportCodeChange(changeDescription string, n rhino.Node) {
	if s.assertOnChange {
		panic(jutil.NewIllegalStateExceptionMsg("Normalize constraints violated:\n" + changeDescription))
	}
	s.compiler.ReportChangeToEnclosingScope(n)
}

// Java: Normalize.NormalizeStatements.enterScope(NodeTraversal)
func (s *Normalize_NormalizeStatements) EnterScope(t *NodeTraversal) {
	if s.useFastContextualRenamer {
		s.fastEnterScope(t)
	} else if s.makeDeclaredNamesUnique != nil {
		s.makeDeclaredNamesUnique.EnterScope(t)
	}
}

// Java: Normalize.NormalizeStatements.exitScope(NodeTraversal)
func (s *Normalize_NormalizeStatements) ExitScope(t *NodeTraversal) {
	if s.useFastContextualRenamer {
		s.fastExitScope(t)
	} else if s.makeDeclaredNamesUnique != nil {
		s.makeDeclaredNamesUnique.ExitScope(t)
	}
}

// Java: Normalize.NormalizeStatements.shouldTraverse(NodeTraversal, Node, Node)
func (s *Normalize_NormalizeStatements) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if !s.useFastContextualRenamer && s.makeDeclaredNamesUnique != nil {
		_ = s.makeDeclaredNamesUnique.ShouldTraverse(t, n, parent)
	}
	if n.GetFirstChild().IsNil() {
		return true
	}

	oldParent := n.GetParent()
	s.DoStatementNormalizations(t, n)
	jutil.CheckStateMsg(
		n.GetParent() == oldParent,
		"n's parent changed during normalization. It needs to stay the same so that we properly"+
			" visit any added nodes.")
	return true
}

// Visit runs the handler of s.postOrder for n's kind.
//
// Java: Normalize.NormalizeStatements.visit(NodeTraversal, Node, Node)
func (s *Normalize_NormalizeStatements) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !s.useFastContextualRenamer && s.makeDeclaredNamesUnique != nil {
		s.makeDeclaredNamesUnique.Visit(t, n, parent)
	}
	s.postOrder.Visit(t, n)
}

// visitName renames x if it is shadowed (with the fast contextual renamer;
// otherwise makeDeclaredNamesUnique did) and marks it constant when it is.
//
// Java: Normalize.NormalizeStatements.visit(NodeTraversal, Node, Node) (NAME)
func (s *Normalize_NormalizeStatements) visitName(t *NodeTraversal, x ast.Name) {
	if !s.useFastContextualRenamer {
		s.AnnotateConstantsByConvention(x)
		s.AnnotateFunctionExpressionNameAsConstant(x)
		return
	}
	n := x.Node()
	s.fastVisitNameOrImportStar(t, n)
	if n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
		return
	}
	if !s.isNormalizedObfuscated && s.convention.IsConstant(x.Ident()) {
		jutil.CheckStateMsg1(!s.assertOnChange, "Not marked as constant when it should be: %s", n)
		n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	} else if normalizeIsFunctionExpressionName(x) {
		n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
}

// visitImportStar renames x if it is shadowed (with the fast contextual
// renamer; otherwise makeDeclaredNamesUnique did).
func (s *Normalize_NormalizeStatements) visitImportStar(t *NodeTraversal, x ast.ImportStar) {
	if s.useFastContextualRenamer {
		s.fastVisitNameOrImportStar(t, x.Node())
	}
}

// visitWhile turns while (cond) body into for (; cond; ) body.
//
// Java: Normalize.NormalizeStatements.visit(NodeTraversal, Node, Node) (WHILE)
func (s *Normalize_NormalizeStatements) visitWhile(_ *NodeTraversal, x ast.While) {
	n := x.Node()
	cond := x.Cond().Node()
	n.SetToken(rhino.Token_FOR)
	empty := rhino.IR_Empty()
	empty.SrcrefIfMissing(n)
	empty.InsertBefore(cond)
	empty.CloneNode().InsertAfter(cond)
	s.ReportCodeChange("WHILE node", n)
}

// AnnotateConstantsByConvention marks names and properties that are
// constants by convention.
//
// Java: Normalize.NormalizeStatements.annotateConstantsByConvention(Node)
func (s *Normalize_NormalizeStatements) AnnotateConstantsByConvention(x ast.Name) {
	// Need to check that variables have not been renamed, to determine whether
	// coding conventions still apply.
	if s.compiler.GetLifeCycleStage().IsNormalizedObfuscated() {
		return
	}

	n := x.Node()
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
func (s *Normalize_NormalizeStatements) AnnotateFunctionExpressionNameAsConstant(x ast.Name) {
	if normalizeIsFunctionExpressionName(x) {
		x.Node().PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
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
// or parents: the handlers of s.preOrder.
//
// Java: Normalize.NormalizeStatements.doStatementNormalizations(Node)
func (s *Normalize_NormalizeStatements) DoStatementNormalizations(t *NodeTraversal, n rhino.Node) {
	s.preOrder.Visit(t, n)
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
func (s *Normalize_NormalizeStatements) normalizeBlock(_ *NodeTraversal, x ast.Block) {
	s.normalizeStatements(x.Node())
	if NodeUtil_IsClassStaticBlock(x.Node()) {
		s.MoveNamedFunctions(x)
	}
}

// Java: Normalize.NormalizeStatements.doStatementNormalizations(Node) (LABEL)
func (s *Normalize_NormalizeStatements) normalizeLabel(_ *NodeTraversal, x ast.Label) {
	s.NormalizeLabels(x)
	s.ExtractForInitializer(x.Node(), rhino.NilNode, rhino.NilNode)
}

// Java: Normalize.NormalizeStatements.doStatementNormalizations(Node) (FUNCTION)
func (s *Normalize_NormalizeStatements) normalizeFunction(_ *NodeTraversal, x ast.Function) {
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
		var insertBefore rhino.Node
		if before.IsNil() {
			insertBefore = c
		} else {
			insertBefore = before
		}
		var insertBeforeParent rhino.Node
		if before.IsNil() {
			insertBeforeParent = n
		} else {
			insertBeforeParent = beforeParent
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

// Normalize_DuplicateDeclarationHandler is the ScopeCreator duplicate
// declaration handler.
type Normalize_DuplicateDeclarationHandler struct {
	outer *Normalize

	hasOkDuplicateDeclaration map[*Var]struct{}
	okDupBuf                  [8]*Var
	okDupLen                  int
}

// Java: Normalize.DuplicateDeclarationHandler()
func NewNormalize_DuplicateDeclarationHandler(outer *Normalize) *Normalize_DuplicateDeclarationHandler {
	return &Normalize_DuplicateDeclarationHandler{
		outer: outer,
	}
}

// OnRedeclaration removes duplicate VAR declarations discovered during scope
// creation.
//
// Java: Normalize.DuplicateDeclarationHandler.onRedeclaration(Scope, String, Node, CompilerInput)
func (h *Normalize_DuplicateDeclarationHandler) OnRedeclaration(
	s *Scope, name string, n rhino.Node, input *CompilerInput) {
	jutil.CheckState(n.IsName())
	parent := n.GetParent()
	v := s.GetVar(name)

	if s.IsGlobal() {
		// We allow variables to be duplicate declared if one
		// declaration appears in source and the other in externs.
		// This deals with issues where a browser built-in is declared
		// in one browser but not in another.
		if v.IsExtern() && !input.IsExtern() {
			for i := 0; i < h.okDupLen; i++ {
				if h.okDupBuf[i] == v {
					goto duplicateNotOk
				}
			}
			if h.okDupLen < len(h.okDupBuf) {
				h.okDupBuf[h.okDupLen] = v
				h.okDupLen++
				return
			}
			if h.hasOkDuplicateDeclaration == nil {
				h.hasOkDuplicateDeclaration = make(map[*Var]struct{}, 32)
			}
			if _, exists := h.hasOkDuplicateDeclaration[v]; !exists {
				h.hasOkDuplicateDeclaration[v] = struct{}{}
				return
			}
		}
	}
duplicateNotOk:

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
	n rhino.Node, parent rhino.Node, grandparent rhino.Node) {
	name, ok := ast.AsName(n)
	jutil.CheckState(ok)
	if value, ok := name.Init(); ok {
		// The  *  is being initialize, preserve the new value.
		n.Detach()
		// Convert "var name = value" to "name = value"
		value.Node().Detach()
		replacement := h.outer.astFactory.CreateAssign(n, value.Node())
		replacement.SetJSDocInfo(parent.GetJSDocInfo())
		replacement.SrcrefIfMissing(parent)
		var statement rhino.Node
		if parent.CanResetInPlace() && replacement.CanResetInPlace() {
			statement = parent.MutateAttachedToExprResultInPlace(replacement)
		} else {
			statement = NodeUtil_NewExpr(replacement)
			parent.ReplaceWith(statement)
		}
		h.outer.ReportCodeChange("Duplicate VAR declaration", statement)
	} else {
		// It is an empty reference. Remove it.
		jutil.CheckStateMsg1(NodeUtil_IsStatementBlock(grandparent), "%s", grandparent)
		parent.Detach()
		h.outer.ReportCodeChange("Duplicate VAR declaration", grandparent)
	}
}
