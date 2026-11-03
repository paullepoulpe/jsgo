// Port of com/google/javascript/jscomp/CrossChunkReferenceCollector.java @ closure-compiler v20261028

/*
 * Copyright 2017 The Closure Compiler Authors.
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

// CrossChunkReferenceCollector collects global variable references for use
// by CrossChunkCodeMotion.
type CrossChunkReferenceCollector struct {
	vars         []*Var
	varsByNameID map[uint32]*Var
	maxVarIndex  int

	stmtSlab jutil.Slab[CrossChunkReferenceCollector_TopLevelStatement]

	// List of all top-level statements in the order they appear in the AST.
	topLevelStatements []*CrossChunkReferenceCollector_TopLevelStatement

	// JavaScript compiler to use in traversing.
	compiler AbstractCompiler

	referenceCollector *ReferenceCollector

	statementCounter       int
	topLevelStatementDraft *CrossChunkReferenceCollector_TopLevelStatementDraft
	draftVal               CrossChunkReferenceCollector_TopLevelStatementDraft
	draftRefsBuf           [16]*Reference
	draftVarsBuf           [16]*Var
}

// NewCrossChunkReferenceCollector: constructor initializes reference collector.
//
// Java: CrossChunkReferenceCollector(AbstractCompiler, ScopeCreator)
func NewCrossChunkReferenceCollector(compiler AbstractCompiler, creator ScopeCreator) *CrossChunkReferenceCollector {
	c := &CrossChunkReferenceCollector{
		compiler:    compiler,
		maxVarIndex: -1,
	}
	c.draftVal.nonDeclarationReferences = c.draftRefsBuf[:0]
	c.draftVal.nonDeclarationVars = c.draftVarsBuf[:0]
	c.referenceCollector = NewReferenceCollectorWithHoistedOrder(
		compiler,
		crossChunkTopLevelStatementBehavior{c: c},
		creator,
		c.shouldCollectVar,
		/* collectHoistedFunctionsInReferenceOrder= */ false)
	return c
}

// Java: CrossChunkReferenceCollector.shouldCollectVar(Var)
func (c *CrossChunkReferenceCollector) shouldCollectVar(v *Var) bool {
	// Only global, non-exported names can be moved
	if !v.IsGlobal() {
		return false
	}
	if _, ok := c.varsByNameID[v.nameID]; ok {
		return true
	}
	return !c.compiler.GetCodingConvention().IsExportedWithLocal(v.GetName() /* local= */, false)
}

func (c *CrossChunkReferenceCollector) allocTopLevelStatement(draft *CrossChunkReferenceCollector_TopLevelStatementDraft) *CrossChunkReferenceCollector_TopLevelStatement {
	var refs []*Reference
	var vars []*Var
	if n := len(draft.nonDeclarationReferences); n > 0 {
		refs = append([]*Reference(nil), draft.nonDeclarationReferences...)
		if len(draft.nonDeclarationVars) == n {
			vars = append([]*Var(nil), draft.nonDeclarationVars...)
		}
	}
	s := c.stmtSlab.New()
	*s = CrossChunkReferenceCollector_TopLevelStatement{
		outer:                    c,
		originalOrder:            draft.originalOrder,
		chunk:                    draft.chunk,
		statementNode:            draft.statementNode,
		nonDeclarationReferences: refs,
		nonDeclarationVars:       vars,
		declaredNameReference:    draft.declaredNameReference,
		declaredVar:              draft.declaredVar,
		declaredValueNode:        draft.declaredValueNode,
	}
	return s
}

// Process collects references and top-level statements from the given externs
// and root.
//
// Java: CrossChunkReferenceCollector.process(Node, Node)
func (c *CrossChunkReferenceCollector) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckStateMsg(len(c.topLevelStatements) == 0 && c.statementCounter == 0, "process() called more than once")
	c.referenceCollector.Process(externs, root)
}

// ProcessRoot collects references and top-level statements from the given root.
//
// Go-only: language-adapter — the two-argument overload keeps the plain
// name because it implements CompilerPass.process(Node, Node).
//
// Java: CrossChunkReferenceCollector.process(Node)
func (c *CrossChunkReferenceCollector) ProcessRoot(root rhino.Node) {
	jutil.CheckStateMsg(len(c.topLevelStatements) == 0 && c.statementCounter == 0, "process() called more than once")
	c.referenceCollector.ProcessRoot(root)
}

// GetReferences gets the reference collection for the given variable.
//
// Java: CrossChunkReferenceCollector.getReferences(Var)
func (c *CrossChunkReferenceCollector) GetReferences(v *Var) *ReferenceCollection {
	return c.referenceCollector.GetReferences(v)
}

// Java: CrossChunkReferenceCollector.getGlobalVariableNamesMap()
func (c *CrossChunkReferenceCollector) GetGlobalVariableNamesMap() *OrderedStringMap[*Var] {
	out := NewOrderedStringMapWithCapacity[*Var](len(c.vars))
	for _, v := range c.vars {
		out.Put(v.GetName(), v)
	}
	return out
}

type crossChunkTopLevelStatementBehavior struct {
	c *CrossChunkReferenceCollector
}

var _ ReferenceCollector_NodeBehavior = crossChunkTopLevelStatementBehavior{}

// Java: CrossChunkReferenceCollector.TopLevelStatementCollectingBehavior.beforeTraverseNode(NodeTraversal, Node, Node)
func (b crossChunkTopLevelStatementBehavior) BeforeTraverseNode(t *NodeTraversal, n, parent rhino.Node) bool {
	c := b.c
	if !parent.IsNil() && NodeUtil_IsTopLevel(parent) {
		jutil.CheckStateMsg1(c.topLevelStatementDraft == nil, "%s", n)
		c.topLevelStatementDraft = c.InitializeDraftStatement(t.GetChunk(), n)
	}
	return true
}

// Java: CrossChunkReferenceCollector.TopLevelStatementCollectingBehavior.afterVisitNode(NodeTraversal, Node, Node, Reference)
func (b crossChunkTopLevelStatementBehavior) AfterVisitNode(t *NodeTraversal, n, parent rhino.Node, addedReference *Reference) {
	c := b.c
	draft := c.topLevelStatementDraft
	if draft == nil {
		return
	}
	if n == draft.statementNode {
		c.topLevelStatements = append(c.topLevelStatements, c.allocTopLevelStatement(draft))
		c.topLevelStatementDraft = nil
	} else if addedReference != nil && n.IsName() {
		id := n.GetStringID()
		v := addedReference.getSymbolV()
		if existing, ok := c.varsByNameID[id]; ok {
			jutil.CheckState(existing == v || (existing != nil && existing.Equals(v)))
		} else {
			if c.varsByNameID == nil {
				c.varsByNameID = make(map[uint32]*Var, 64)
			}
			c.varsByNameID[id] = v
			c.vars = append(c.vars, v)
			if idx := v.GetIndex(); idx > c.maxVarIndex {
				c.maxVarIndex = idx
			}
		}
		if addedReference.GetNode() == draft.declaredNameNode {
			draft.declaredNameReference = addedReference
			draft.declaredVar = v
		} else {
			draft.nonDeclarationReferences = append(draft.nonDeclarationReferences, addedReference)
			draft.nonDeclarationVars = append(draft.nonDeclarationVars, v)
		}
	}
}

// Java: CrossChunkReferenceCollector.TopLevelStatementCollectingBehavior.afterExitScope(NodeTraversal, ReferenceMap)
func (b crossChunkTopLevelStatementBehavior) AfterExitScope(t *NodeTraversal, refs *ReferenceCollector) {
}

// Java: CrossChunkReferenceCollector.initializeDraftStatement(JSChunk, Node)
func (c *CrossChunkReferenceCollector) InitializeDraftStatement(chunk *JSChunk, statementNode rhino.Node) *CrossChunkReferenceCollector_TopLevelStatementDraft {
	draft := &c.draftVal
	refs := draft.nonDeclarationReferences[:0]
	vars := draft.nonDeclarationVars[:0]
	*draft = CrossChunkReferenceCollector_TopLevelStatementDraft{
		originalOrder:            c.statementCounter,
		chunk:                    chunk,
		statementNode:            statementNode,
		nonDeclarationReferences: refs,
		nonDeclarationVars:       vars,
	}
	c.statementCounter++
	// Determine whether this statement declares a name or not.
	// If so, save its name node and value node, if any.
	if decl, ok := ast.AsNameDeclaration(statementNode); ok {
		// variable declaration
		if first, ok := decl.FirstDeclarator(); ok {
			if name, ok := first.AsName(); ok {
				draft.declaredNameNode = name.Node()
				if initExpr, ok := name.Init(); ok {
					draft.declaredValueNode = initExpr.Node()
				}
			}
		}
	} else if cls, ok := ast.AsClass(statementNode); ok {
		draft.declaredNameNode = cls.NameNode()
		draft.declaredValueNode = statementNode
	} else if fn, ok := ast.AsFunction(statementNode); ok {
		// function declaration
		draft.declaredNameNode = fn.NameNode().Node()
		draft.declaredValueNode = statementNode
	} else if exprStmt, ok := ast.AsExprStmt(statementNode); ok {
		expr := exprStmt.Expr()
		if assign, ok := expr.AsAssign(); ok {
			lhs := assign.Target()
			rhs := assign.Value()
			if name, ok := lhs.AsName(); ok {
				// `varName = value;`
				draft.declaredNameNode = name.Node()
				draft.declaredValueNode = rhs.Node()
			} else if normalGet, ok := lhs.AsNormalGet(); ok &&
				(normalGet.Token() == rhino.Token_GETPROP || CrossChunkReferenceCollector_IsPotentiallyMoveableGetElem(lhs.Node())) {
				obj := normalGet.Object()
				for prop, ok := obj.AsGetProp(); ok; prop, ok = obj.AsGetProp() {
					obj = prop.Object()
				}
				if name, ok := obj.AsName(); ok {
					// `varName.some.property = value;`
					// OR
					// `varName[someProperty] = value;`
					draft.declaredNameNode = name.Node()
					draft.declaredValueNode = rhs.Node()
				}
			}
		} else if call, ok := expr.AsCall(); ok {
			var nameNode rhino.Node
			var valueNode rhino.Node
			relationship :=
				c.compiler.GetCodingConvention().GetClassesDefinedByCall(call.Node())
			if relationship != nil {
				// Check for $jscomp.inherits(SubC, SuperC), goog.inherits(Sub, SuperC), etc.
				declaredName := relationship.subclassName
				for callArg := range call.Args() {
					// We're assuming that the child class must be an argument to the function that
					// establishes its inheritance, which is true for `goog.inherits()` and
					// `$jscomp.inherits()`
					// TODO(bradfordcsmith): handle cases like `goog.inherits(x.ChildClass, SuperClass)`
					if argName, ok := callArg.AsName(); ok && declaredName == argName.Ident() {
						nameNode = argName.Node()
						break
					}
				}
			} else if NodeUtil_IsObjectDefinePropertiesDefinition(call.Node()) {
				// Check for $jscomp$global.Object.defineProperties.
				targetArg := call.MustFirstArg()

				// Get the global var being referenced in the first parameter of Object.defineProperties.
				// Can be 'Foo' or 'Foo.prototype'.
				if targetName, ok := targetArg.AsName(); ok {
					nameNode = targetName.Node()
				} else {
					nameNode = targetArg.Node().GetFirstChild()
				}

				// Second parameter of defineProperties is the value node.
				secondArg, _ := targetArg.Next()
				valueNode = secondArg.Node()
			}

			if !nameNode.IsNil() {
				draft.declaredNameNode = nameNode
				draft.declaredValueNode = valueNode
			}
		}
	}
	return draft
}

// Java: CrossChunkReferenceCollector.isPotentiallyMoveableGetElem(Node)
func CrossChunkReferenceCollector_IsPotentiallyMoveableGetElem(maybeGetElem rhino.Node) bool {
	getElem, ok := ast.AsGetElem(maybeGetElem)
	if !ok {
		return false
	}

	// We'll only consider a GETELEM as potentially moveable if the property being indexed is on
	// `Symbol` (e.g. `Symbol.iterator`). We don't generally want to consider GETELEM nodes as
	// moveable as registries (ex. `GlobalRegistry['Foo'] = Foo`) will create a reference from `Foo`
	// to `GlobalRegistry`. As other pieces of code reference `GlobalRegistry` this will ultimately
	// cause a large cluster of symbols that would have to be considered as moving together, but
	// likely will no to be able to.
	// Therefore, we we'll only consider a subset of cases, like `Symbol.iterator`, as moveable as
	// these large reference clusters are unlikely in this case. This should allow classes that
	// define these properties to potentialy be moved to later chunks.
	property, ok := getElem.Index().AsGetProp()
	if !ok {
		return false
	}
	name, ok := property.Object().AsName()
	return ok && name.Ident() == "Symbol"
}

// Java: CrossChunkReferenceCollector.getTopLevelStatements()
func (c *CrossChunkReferenceCollector) GetTopLevelStatements() []*CrossChunkReferenceCollector_TopLevelStatement {
	return c.topLevelStatements
}

// CanMoveValue determines whether the given value is eligible to be moved
// across chunks.
//
// Java: CrossChunkReferenceCollector.canMoveValue(Scope, Node)
func (c *CrossChunkReferenceCollector) CanMoveValue(scope *Scope, valueNode rhino.Node) bool {
	// the value is only movable if it's
	// - nothing,
	// - a constant literal,
	// - a function
	// - class without computed properties, or whose computed properties expressions are moveable
	// - an array/object literal of movable values
	// - a template literal with movable substitutions
	// - a function stub generated by CrossChunkMethodMotion.
	if valueNode.IsNil() || NodeUtil_IsLiteralValue(valueNode, true) || valueNode.IsFunction() {
		return true
	}

	// Consider the node eligible for move if the current statement or expression is annotated
	// with @pureOrBreakMyCode.
	jsDocInfo := NodeUtil_GetBestJSDocInfo(valueNode)
	if jsDocInfo != nil && jsDocInfo.IsPureOrBreakMyCode() {
		return true
	}

	switch valueNode.GetToken() {
	case rhino.Token_CLASS:
		// NOTE: Why aren't class extends expressions checked here?
		cls := ast.MustClass(valueNode)
		for member := range cls.Body().Members() {
			if cp, ok := member.AsComputedProp(); ok {
				method, hasMethod := cp.Value()
				jutil.CheckStateMsg1(hasMethod && method.Node().IsFunction(), "%s", method.Node())
				if !c.CanMoveValue(scope, cp.Key().Node()) {
					return false
				}
			} else if cf, ok := member.AsComputedFieldDef(); ok {
				if !c.CanMoveValue(scope, cf.Key().Node()) {
					return false
				}
			} else if mf, ok := member.AsMemberFieldDef(); ok {
				if val, ok := mf.Value(); ok && mf.IsStatic() && !c.CanMoveValue(scope, val.Node()) {
					// Unlike non-static fields,the RHS of a static field is executed at
					// class definition time. So, we must check canMoveValue for static fields.
					return false
				}
			} else if _, ok := member.AsBlock(); ok {
				// TODO(bradfordcsmith): Ideally could move these in some cases, fix later
				return false
			} else {
				_, isMethod := member.AsMethodDef()
				jutil.CheckStateMsg1(isMethod, "%s", member.Node())
			}
		}
		return true

	case rhino.Token_CALL:
		// In general it is not safe to move function calls, but we carve out an exception
		// for the special stub method calls used for CrossChunkMethodMotion.
		// Case: `JSCompiler_stubMethod(x)`
		call := ast.MustCall(valueNode)
		calleeName, ok := call.Callee().AsName()
		return ok && calleeName.Ident() == CrossChunkMethodMotion_STUB_METHOD_NAME

	case rhino.Token_ARRAYLIT:
		// Movable if all of the array values are movable.
		for elem := range ast.MustArrayLit(valueNode).Elements() {
			if !c.CanMoveValue(scope, elem.Node()) {
				return false
			}
		}
		return true

	case rhino.Token_OBJECTLIT:
		// Movable if all of the keys and values are movable.
		for member := range ast.MustObjectLit(valueNode).Members() {
			if _, ok := member.AsMethodDef(); ok {
				continue
			} else if cp, ok := member.AsComputedProp(); ok {
				val, _ := cp.Value()
				if !c.CanMoveValue(scope, cp.Key().Node()) ||
					!c.CanMoveValue(scope, val.Node()) {
					return false
				}
			} else if sk, ok := member.AsStringKey(); ok {
				if !c.CanMoveValue(scope, sk.Value().Node()) {
					return false
				}
			} else if sp, ok := member.AsObjectSpread(); ok {
				if !c.CanMoveValue(scope, sp.Expr().Node()) {
					return false
				}
			} else {
				panic(jutil.NewIllegalStateExceptionMsg(
					"Unexpected child of OBJECTLIT: " + member.Node().ToStringTree()))
			}
		}
		return true

	case rhino.Token_NAME:
		// If the value is guaranteed to never be changed after
		// this reference, then we can move it.
		v := scope.GetVarByID(valueNode.GetStringID())
		if v != nil && v.IsGlobal() {
			refCollection := c.GetReferences(v)
			if refCollection != nil &&
				refCollection.IsWellDefined() &&
				refCollection.IsAssignedOnceInLifetime() {
				return true
			}
		}

	case rhino.Token_TEMPLATELIT:
		// A template literal is movable if all of the substitutions it contains are movable.
		for part := range ast.MustTemplateLit(valueNode).Parts() {
			if sub, ok := part.AsTemplateLitSub(); ok {
				if !c.CanMoveValue(scope, sub.Expr().Node()) {
					return false
				}
			} else {
				_, isStr := part.AsTemplateLitString()
				jutil.CheckStateMsg1(isStr, "%s", part.Node())
			}
		}
		return true

	case rhino.Token_GETPROP:
		// `Symbol` is a built-in JavaScript object that is guaranteed to be unique. Therefore it
		// is movable. Any properties on `Symbol` are also movable (e.g `Symbol.iterator`).
		name, ok := ast.MustGetProp(valueNode).Object().AsName()
		return ok && name.Ident() == "Symbol"

	default:
	}

	return false
}

// CrossChunkReferenceCollector_TopLevelStatement represents a top-level
// statement and the references to global names it contains.
type CrossChunkReferenceCollector_TopLevelStatement struct {
	outer *CrossChunkReferenceCollector

	// 0-based index indicating original order of this statement in the source.
	originalOrder int

	chunk                    *JSChunk
	statementNode            rhino.Node
	nonDeclarationReferences []*Reference
	nonDeclarationVars       []*Var
	declaredNameReference    *Reference
	declaredVar              *Var
	declaredValueNode        rhino.Node
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.getOriginalOrder()
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetOriginalOrder() int {
	return s.originalOrder
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.getChunk()
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetChunk() *JSChunk {
	return s.chunk
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.getStatementNode()
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetStatementNode() rhino.Node {
	return s.statementNode
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.getNonDeclarationReferences()
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetNonDeclarationReferences() []*Reference {
	return s.nonDeclarationReferences
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.isDeclarationStatement()
func (s *CrossChunkReferenceCollector_TopLevelStatement) IsDeclarationStatement() bool {
	return s.declaredNameReference != nil
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.getDeclaredNameReference()
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetDeclaredNameReference() *Reference {
	return jutil.CheckNotNull(s.declaredNameReference)
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.getDeclaredValueNode()
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetDeclaredValueNode() rhino.Node {
	return s.declaredValueNode
}

// Java: CrossChunkReferenceCollector.TopLevelStatement.isMovableDeclaration()
func (s *CrossChunkReferenceCollector_TopLevelStatement) IsMovableDeclaration() bool {
	return s.IsDeclarationStatement() &&
		s.outer.CanMoveValue(s.declaredNameReference.GetScope(), s.declaredValueNode)
}

// CrossChunkReferenceCollector_TopLevelStatementDraft holds statement info
// temporarily while the statement is being traversed.
type CrossChunkReferenceCollector_TopLevelStatementDraft struct {
	// 0-based index indicating original order of this statement in the source.
	originalOrder int

	chunk                    *JSChunk
	statementNode            rhino.Node
	nonDeclarationReferences []*Reference
	nonDeclarationVars       []*Var
	declaredValueNode        rhino.Node
	declaredNameNode         rhino.Node
	declaredNameReference    *Reference
	declaredVar              *Var
}
