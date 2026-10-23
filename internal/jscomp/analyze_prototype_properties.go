// Port of com/google/javascript/jscomp/AnalyzePrototypeProperties.java @ closure-compiler v20261015

/*
 * Copyright 2006 The Closure Compiler Authors.
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
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AnalyzePrototypeProperties analyzes properties on prototypes.
//
// Uses a reference graph to analyze prototype properties. Each unique
// property name is represented by a node in this graph. An edge from
// property A to property B means that there's a GETPROP access of a property
// B on some object inside of a method named A.
//
// Global functions are also represented by nodes in this graph, with
// similar semantics.
type analyzeProtoPropEdge struct {
	chunk *JSChunk
	dest  *AnalyzePrototypeProperties_NameInfo
}

type AnalyzePrototypeProperties struct {
	compiler               AbstractCompiler
	canModifyExterns       bool
	anchorUnusedVars       bool
	rootScopeUsesAreGlobal bool

	chunkGraph *JSChunkGraph
	firstChunk *JSChunk

	// A dummy node for representing global references.
	globalNode *AnalyzePrototypeProperties_NameInfo

	// A dummy node for representing extern references.
	externNode *AnalyzePrototypeProperties_NameInfo

	// A dummy node for representing all anonymous functions with no names.
	anonymousNode *AnalyzePrototypeProperties_NameInfo

	// All the real NameInfo for prototype properties, hashed by the name
	// of the property that they represent.
	propertyNameInfo *jutil.LinkedHashMap[string, *AnalyzePrototypeProperties_NameInfo]

	// All the NameInfo for global functions, hashed by the name of the
	// global variable that it's assigned to.
	varNameInfo *jutil.LinkedHashMap[string, *AnalyzePrototypeProperties_NameInfo]

	seededNodes []*AnalyzePrototypeProperties_NameInfo

	scopeCreator *SyntacticScopeCreator

	nameInfoChunks   [][]AnalyzePrototypeProperties_NameInfo
	nameInfoIdx      int
	globalFuncSlab   []AnalyzePrototypeProperties_GlobalFunction
	assignPropSlab   []AnalyzePrototypeProperties_AssignmentPrototypeProperty
	edgeSlab         []analyzeProtoPropEdge
}

func (a *AnalyzePrototypeProperties) allocGlobalFunction(nameNode *rhino.Node, var_ *Var, chunk *JSChunk) *AnalyzePrototypeProperties_GlobalFunction {
	parent := nameNode.GetParent()
	jutil.CheckState(
		(NodeUtil_IsNameDeclaration(parent) && var_.IsGlobal()) ||
			NodeUtil_IsFunctionDeclaration(parent))
	if len(a.globalFuncSlab) == 0 {
		a.globalFuncSlab = make([]AnalyzePrototypeProperties_GlobalFunction, jscompMaxSmallAlloc/int(unsafe.Sizeof(AnalyzePrototypeProperties_GlobalFunction{})))
	}
	g := &a.globalFuncSlab[0]
	a.globalFuncSlab = a.globalFuncSlab[1:]
	g.var_ = var_
	g.chunk = chunk
	return g
}

func (a *AnalyzePrototypeProperties) allocAssignmentPrototypeProperty(node *rhino.Node, rootVar *Var, chunk *JSChunk) *AnalyzePrototypeProperties_AssignmentPrototypeProperty {
	if len(a.assignPropSlab) == 0 {
		a.assignPropSlab = make([]AnalyzePrototypeProperties_AssignmentPrototypeProperty, jscompMaxSmallAlloc/int(unsafe.Sizeof(AnalyzePrototypeProperties_AssignmentPrototypeProperty{})))
	}
	p := &a.assignPropSlab[0]
	a.assignPropSlab = a.assignPropSlab[1:]
	p.exprNode = node
	p.rootVar = rootVar
	p.chunk = chunk
	return p
}

func (a *AnalyzePrototypeProperties) allocNameInfo(name string) *AnalyzePrototypeProperties_NameInfo {
	if len(a.nameInfoChunks) == 0 || a.nameInfoIdx == len(a.nameInfoChunks[len(a.nameInfoChunks)-1]) {
		nextCap := 64
		if len(a.nameInfoChunks) > 0 {
			nextCap = len(a.nameInfoChunks[len(a.nameInfoChunks)-1]) * 2
			if nextCap > 112 {
				nextCap = 112
			}
		}
		a.nameInfoChunks = append(a.nameInfoChunks, make([]AnalyzePrototypeProperties_NameInfo, nextCap))
		a.nameInfoIdx = 0
	}
	chunk := a.nameInfoChunks[len(a.nameInfoChunks)-1]
	info := &chunk[a.nameInfoIdx]
	a.nameInfoIdx++
	info.outer = a
	info.name = name
	info.declarations.InitWithBuffer(info.declBuf[:])
	info.outEdges = info.outEdgesBuf[:0]
	return info
}

func (a *AnalyzePrototypeProperties) connectSymbolNodes(src *AnalyzePrototypeProperties_NameInfo, chunk *JSChunk, dest *AnalyzePrototypeProperties_NameInfo) {
	if src == a.globalNode {
		if dest.hasGlobalEdge && dest.globalEdgeChunk == chunk {
			return
		}
		dest.hasGlobalEdge = true
		dest.globalEdgeChunk = chunk
		dest.MarkReference(chunk)
		if !dest.seeded {
			dest.seeded = true
			a.seededNodes = append(a.seededNodes, dest)
		}
		return
	}
	if src == a.externNode {
		if dest.hasExternEdge && dest.externEdgeChunk == chunk {
			return
		}
		dest.hasExternEdge = true
		dest.externEdgeChunk = chunk
		dest.MarkReference(chunk)
		if !dest.seeded {
			dest.seeded = true
			a.seededNodes = append(a.seededNodes, dest)
		}
		return
	}
	if dest.lastSrc == src && dest.lastChunk == chunk {
		return
	}
	if n := len(src.outEdges); n > 0 && src.outEdges[n-1].dest == dest && src.outEdges[n-1].chunk == chunk {
		return
	}
	dest.lastSrc = src
	dest.lastChunk = chunk
	if len(src.outEdges) == cap(src.outEdges) {
		newCap := max(cap(src.outEdges)*2, 8)
		if newCap <= 256 {
			if len(a.edgeSlab) < newCap {
				a.edgeSlab = make([]analyzeProtoPropEdge, jscompMaxSmallAlloc/int(unsafe.Sizeof(analyzeProtoPropEdge{})))
			}
			newBuf := a.edgeSlab[:len(src.outEdges):newCap]
			a.edgeSlab = a.edgeSlab[newCap:]
			copy(newBuf, src.outEdges)
			src.outEdges = newBuf
		}
	}
	src.outEdges = append(src.outEdges, analyzeProtoPropEdge{chunk: chunk, dest: dest})
}

// Release drops the scope creator of this AnalyzePrototypeProperties.
func (a *AnalyzePrototypeProperties) Release() {
	if a != nil && a.scopeCreator != nil {
		a.scopeCreator = nil
	}
}

// Constants for symbol types, for easier readability.
const AnalyzePrototypeProperties_PROPERTY = AnalyzePrototypeProperties_SymbolType_PROPERTY
const AnalyzePrototypeProperties_VAR = AnalyzePrototypeProperties_SymbolType_VAR

// Properties that are implicitly used as part of the JS language.
var AnalyzePrototypeProperties_IMPLICITLY_USED_PROPERTIES = []string{"length", "toString", "valueOf"}

// NewAnalyzePrototypeProperties creates a new pass for analyzing prototype
// properties.
//
// compiler is the compiler; chunkGraph is the graph for resolving chunk
// dependencies; if canModifyExterns, then we can move prototype properties
// that are declared in the externs file; if anchorUnusedVars, then we must
// mark all vars as referenced, even if they are never used; if
// rootScopeUsesAreGlobal, all uses in root level scope are treated as
// references from '[global]', even if they are assignments to a property.
//
// Java: AnalyzePrototypeProperties(AbstractCompiler, JSChunkGraph, boolean, boolean, boolean)
func NewAnalyzePrototypeProperties(
	compiler AbstractCompiler,
	chunkGraph *JSChunkGraph,
	canModifyExterns bool,
	anchorUnusedVars bool,
	rootScopeUsesAreGlobal bool) *AnalyzePrototypeProperties {
	a := &AnalyzePrototypeProperties{
		propertyNameInfo: jutil.NewLinkedHashMap[string, *AnalyzePrototypeProperties_NameInfo](),
		varNameInfo:      jutil.NewLinkedHashMap[string, *AnalyzePrototypeProperties_NameInfo](),
	}
	a.globalNode = NewAnalyzePrototypeProperties_NameInfo(a, "[global]")
	a.externNode = NewAnalyzePrototypeProperties_NameInfo(a, "[extern]")
	a.anonymousNode = NewAnalyzePrototypeProperties_NameInfo(a, "[anonymous]")
	a.compiler = compiler
	a.chunkGraph = chunkGraph
	a.canModifyExterns = canModifyExterns
	a.anchorUnusedVars = anchorUnusedVars
	a.rootScopeUsesAreGlobal = rootScopeUsesAreGlobal

	if chunkGraph.GetChunkCount() > 1 {
		a.firstChunk = chunkGraph.GetRootChunk()
	} else {
		a.firstChunk = nil
	}

	a.globalNode.MarkReference(nil)
	a.externNode.MarkReference(nil)

	for _, property := range AnalyzePrototypeProperties_IMPLICITLY_USED_PROPERTIES {
		nameInfo := a.GetNameInfoForName(property, AnalyzePrototypeProperties_PROPERTY)
		if chunkGraph == nil {
			a.connectSymbolNodes(a.externNode, nil, nameInfo)
		} else if chunkGraph.GetChunkCount() > 0 {
			a.connectSymbolNodes(a.externNode, chunkGraph.GetRootChunk(), nameInfo)
		}
	}
	return a
}

// Java: AnalyzePrototypeProperties.process(Node, Node)
func (a *AnalyzePrototypeProperties) Process(externRoot *rhino.Node, root *rhino.Node) {
	jutil.CheckState(a.compiler.GetLifeCycleStage().IsNormalized())
	if !a.canModifyExterns {
		NodeTraversal_Traverse(a.compiler, externRoot, &AnalyzePrototypeProperties_ProcessExternProperties{outer: a})
	}

	if a.scopeCreator == nil {
		a.scopeCreator = NewSyntacticScopeCreator(a.compiler)
	}
	var propProc AnalyzePrototypeProperties_ProcessProperties
	propProc.outer = a
	propProc.symbolStack = propProc.symbolStackBuf[:0]
	NodeTraversal_Builder_().SetCompiler(a.compiler).SetScopeCreator(a.scopeCreator).SetCallback(&propProc).Traverse(root)

	a.propagateReferences()
}

func (a *AnalyzePrototypeProperties) propagateReferences() {
	if len(a.seededNodes) == 0 {
		return
	}
	totalNodes := a.propertyNameInfo.Size() + a.varNameInfo.Size() + 2
	ringCap := 16
	for ringCap < totalNodes+1 {
		ringCap <<= 1
	}
	var ringBuf [64]*AnalyzePrototypeProperties_NameInfo
	var ring []*AnalyzePrototypeProperties_NameInfo
	if ringCap <= len(ringBuf) {
		ring = ringBuf[:ringCap]
	} else {
		ring = make([]*AnalyzePrototypeProperties_NameInfo, ringCap)
	}
	var head, tail uint32
	mask := uint32(ringCap - 1)
	for _, info := range a.seededNodes {
		if len(info.outEdges) > 0 && !info.inWorkSet {
			info.inWorkSet = true
			ring[tail&mask] = info
			tail++
		}
	}
	for head != tail {
		idx := head & mask
		start := ring[idx]
		ring[idx] = nil
		head++
		start.inWorkSet = false
		if !start.referenced {
			continue
		}
		startChunk := start.deepestCommonChunkRef
		for _, e := range start.outEdges {
			targetChunk := e.chunk
			if startChunk != nil && a.chunkGraph.DependsOn(startChunk, targetChunk) {
				targetChunk = startChunk
			}
			dest := e.dest
			if dest.MarkReference(targetChunk) && len(dest.outEdges) > 0 && !dest.inWorkSet {
				dest.inWorkSet = true
				ring[tail&mask] = dest
				tail++
			}
		}
	}
}

// GetAllNameInfo returns information on all prototype properties.
//
// Java: AnalyzePrototypeProperties.getAllNameInfo()
func (a *AnalyzePrototypeProperties) GetAllNameInfo() []*AnalyzePrototypeProperties_NameInfo {
	result := make([]*AnalyzePrototypeProperties_NameInfo, 0, a.propertyNameInfo.Size()+a.varNameInfo.Size())
	for info := range a.propertyNameInfo.Values() {
		result = append(result, info)
	}
	for info := range a.varNameInfo.Values() {
		result = append(result, info)
	}
	return result
}

// GetNameInfoForName gets the name info for the property or variable of a
// given name, and creates a new one if necessary.
//
// name is the name of the symbol; type_ is the type of symbol.
//
// Java: AnalyzePrototypeProperties.getNameInfoForName(String, SymbolType)
func (a *AnalyzePrototypeProperties) GetNameInfoForName(name string, type_ AnalyzePrototypeProperties_SymbolType) *AnalyzePrototypeProperties_NameInfo {
	map_ := a.varNameInfo
	if type_ == AnalyzePrototypeProperties_PROPERTY {
		map_ = a.propertyNameInfo
	}
	if info := map_.Get(name); info != nil {
		return info
	}
	nameInfo := a.allocNameInfo(name)
	map_.PutNew(name, nameInfo)
	return nameInfo
}

// AnalyzePrototypeProperties_ProcessProperties is the inner class
// ProcessProperties.
type AnalyzePrototypeProperties_ProcessProperties struct {
	outer *AnalyzePrototypeProperties

	// There are two types of context information on this stack:
	// 1) Every scope has a NameContext corresponding to its scope.
	//    Variables are given VAR contexts.
	//    Prototype properties are given PROPERTY contexts.
	//    The global scope is given the special [global] context.
	//    And function expressions that we aren't able to give a reasonable
	//    name are given a special [anonymous] context.
	// 2) Every assignment of a prototype property of a non-function is
	//    given a name context. These contexts do not have scopes.
	//
	// Go-only: language-adapter — the ArrayDeque used as a stack (push/pop at
	// the head, iterated from the head) is a slice whose top is its last
	// element, iterated from the end.
	symbolStack    []AnalyzePrototypeProperties_NameContext
	symbolStackBuf [16]AnalyzePrototypeProperties_NameContext
}

// push is Deque.push on symbolStack.
//
// Go-only: language-adapter — see symbolStack.
func (p *AnalyzePrototypeProperties_ProcessProperties) push(c AnalyzePrototypeProperties_NameContext) {
	p.symbolStack = append(p.symbolStack, c)
}

// pop is Deque.pop on symbolStack.
//
// Go-only: language-adapter — see symbolStack.
func (p *AnalyzePrototypeProperties_ProcessProperties) pop() {
	if len(p.symbolStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	p.symbolStack = p.symbolStack[:len(p.symbolStack)-1]
}

// Java: AnalyzePrototypeProperties.ProcessProperties.enterScope(NodeTraversal)
func (p *AnalyzePrototypeProperties_ProcessProperties) EnterScope(t *NodeTraversal) {
	a := p.outer
	root := t.GetScopeRoot()
	if root.IsFunction() {
		n := t.GetCurrentNode()
		if propName, ok := p.getPrototypePropertyNameFromRValueString(n); ok {
			p.push(AnalyzePrototypeProperties_NameContext{
				name:      a.GetNameInfoForName(propName, AnalyzePrototypeProperties_PROPERTY),
				scopeRoot: root,
			})
		} else if p.IsGlobalFunctionDeclaration(t, n) {
			parent := n.GetParent()
			var name string
			if parent.IsName() {
				name = parent.GetString() /* VAR */
			} else {
				name = n.GetFirstChild().GetString() /* named function */
			}
			p.push(AnalyzePrototypeProperties_NameContext{
				name:      a.GetNameInfoForName(name, AnalyzePrototypeProperties_VAR),
				scopeRoot: t.GetClosestHoistScopeRoot(),
			})
		} else {
			// NOTE(nicksantos): We use the same anonymous node for all
			// functions that do not have reasonable names. I can't remember
			// at the moment why we do this. I think it's because anonymous
			// nodes can never have in-edges. They're just there as a placeholder
			// for scope information, and do not matter in the edge propagation.
			p.push(AnalyzePrototypeProperties_NameContext{name: a.anonymousNode, scopeRoot: root})
		}
	} else if t.InGlobalScope() {
		p.push(AnalyzePrototypeProperties_NameContext{name: a.globalNode, scopeRoot: root})
	} else {
		// TODO(moz): It's not yet clear if we need another kind of NameContext for block scopes
		// in ES6, use anonymous node for now and investigate later.
		// TODO(b/189993301): Test effects of pushing anonymousNode for member and computed field
		// defs in CrossChunkMethodMotion (and test class fields there as well)
		jutil.CheckStateMsg1(
			NodeUtil_CreatesBlockScope(root) ||
				root.IsModuleBody() ||
				root.IsComputedFieldDef() ||
				root.IsMemberFieldDef(),
			"%s", root)
		p.push(AnalyzePrototypeProperties_NameContext{name: a.anonymousNode, scopeRoot: root})
	}
}

// Java: AnalyzePrototypeProperties.ProcessProperties.exitScope(NodeTraversal)
func (p *AnalyzePrototypeProperties_ProcessProperties) ExitScope(t *NodeTraversal) {
	p.pop()
}

// Java: AnalyzePrototypeProperties.ProcessProperties.shouldTraverse(NodeTraversal, Node, Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if !p.outer.rootScopeUsesAreGlobal {
		// Process assignment of a non-function expression to a prototype property.
		if propName, ok := p.processNonFunctionPrototypeAssignString(n, parent); ok {
			p.push(AnalyzePrototypeProperties_NameContext{
				name: p.outer.GetNameInfoForName(propName, AnalyzePrototypeProperties_PROPERTY),
			})
		}
	}
	return true
}

// Java: AnalyzePrototypeProperties.ProcessProperties.visit(NodeTraversal, Node, Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	a := p.outer
	switch n.GetToken() {
	case rhino.Token_SUPER:
		// Example:
		// class X extends Y {
		//   method() {
		//     return () => super.x;
		//   }
		// }
		// Names associated with the arrow function, the method body, and the method itself
		// should be marked as referencing super, but not the class definition or anything
		// containing it.
		for i := len(p.symbolStack) - 1; i >= 0; i-- {
			context := &p.symbolStack[i]
			context.name.referencesSuper = true
			if NodeUtil_IsMethodDeclaration(context.scopeRoot) {
				break
			}
		}

	case rhino.Token_OPTCHAIN_GETPROP:
		p.AddSymbolUse(n.GetString(), t.GetChunk(), AnalyzePrototypeProperties_PROPERTY)

	case rhino.Token_GETPROP:
		propName := n.GetString()

		if n.IsQualifiedName() {
			if propName == "prototype" {
				if p.HandlePossibleAssignmentToPrototype(t, n) {
					// The reference is being assigned to not read from, so don't record this as a
					// reference.
					return
				}
			} else if a.compiler.GetCodingConvention().IsExportedWithLocal(propName, false /* local= */) {
				// TODO(bradfordcsmith): We don't seem to have any tests that cover this case.
				// This class has no unit tests of its own and it is only used by
				// CrossChunkMethodMotion.
				p.AddGlobalUseOfSymbol(propName, t.GetChunk(), AnalyzePrototypeProperties_PROPERTY)
				return
			} else {
				if parent.IsAssign() && n == parent.GetFirstChild() {
					if _, ok := p.getPrototypePropertyNameFromRValueString(n); ok {
						// e.g. `Foo.prototype.bar = something`
						// We record the declaration of `bar` when we look at the `Foo.prototype`
						// Node via the call to handlePossibleAssignmentToPrototype() above.
						// Now that we're looking at the whole `Foo.prototype.bar` node,
						// we just need to make sure we don't record it as a read reference.
						return
					}
				}
			}
		}

		p.AddSymbolUse(propName, t.GetChunk(), AnalyzePrototypeProperties_PROPERTY)

	case rhino.Token_OBJECTLIT:
		// Make sure that we're not handling object literals being
		// assigned to a prototype, as in:
		// Foo.prototype = {bar: 3, baz: 5};
		if bestLValueEndsWithPrototype(NodeUtil_GetBestLValue(n)) {
			return
		}

		// Fall through.
		fallthrough
	case rhino.Token_OBJECT_PATTERN:
		// `var x = {a: 1, b: 2}` and `var {a: x, b: y} = obj;`
		// should count as a use of property a and b.
		for propNode := n.GetFirstChild(); propNode != nil; propNode = propNode.GetNext() {
			switch propNode.GetToken() {
			case rhino.Token_COMPUTED_PROP,
				rhino.Token_ITER_REST,
				rhino.Token_OBJECT_REST,
				rhino.Token_ITER_SPREAD,
				rhino.Token_OBJECT_SPREAD:

			case rhino.Token_STRING_KEY,
				rhino.Token_GETTER_DEF,
				rhino.Token_SETTER_DEF,
				rhino.Token_MEMBER_FUNCTION_DEF:
				if !propNode.IsQuotedStringKey() {
					// May be STRING, GET, or SET, but NUMBER isn't interesting.
					p.AddSymbolUse(propNode.GetString(), t.GetChunk(), AnalyzePrototypeProperties_PROPERTY)
				}

			default:
				panic(jutil.NewIllegalStateExceptionMsg(
					"Unexpected child of " + n.GetToken().String() + ": " + propNode.ToStringTree()))
			}
		}

	case rhino.Token_CLASS:
		classMembers := n.GetLastChild()
		for child := classMembers.GetFirstChild(); child != nil; child = child.GetNext() {
			if child.IsMemberFunctionDef() || child.IsSetterDef() || child.IsGetterDef() {
				p.ProcessMemberDef(t, child)
			}
		}

	case rhino.Token_NAME:
		name := n.GetString()
		if name == "" {
			break
		}

		var_ := t.LookupVar(name)
		if var_ != nil {
			// Only process global functions.
			if var_.IsGlobal() {
				if var_.GetInitialValue() != nil && var_.GetInitialValue().IsFunction() {
					if t.InGlobalHoistScope() {
						if !p.ProcessGlobalFunctionDeclaration(t, n, var_) {
							p.AddGlobalUseOfSymbol(name, t.GetChunk(), AnalyzePrototypeProperties_VAR)
						}
					} else {
						p.AddSymbolUse(name, t.GetChunk(), AnalyzePrototypeProperties_VAR)
					}
				}

				// If it is not a global, it might be accessing a local of the outer
				// scope. If that's the case the functions between the variable's
				// declaring scope and the variable reference scope cannot be moved.
			} else if varScopeRoot := var_.GetScope().(*Scope).GetRootNode(); varScopeRoot != t.GetScopeRoot() {
				for i := len(p.symbolStack) - 1; i >= 0; i-- {
					context := &p.symbolStack[i]
					if context.scopeRoot == varScopeRoot {
						break
					}

					context.name.readClosureVariables = true
				}
			}
		}

	default:
	}
	// Process prototype assignments to non-functions.
	if !a.rootScopeUsesAreGlobal && p.ProcessNonFunctionPrototypeAssign(n, parent) != nil {
		p.pop()
	}
}

// Java: AnalyzePrototypeProperties.ProcessProperties.addSymbolUse(String, JSChunk, SymbolType)
func (p *AnalyzePrototypeProperties_ProcessProperties) AddSymbolUse(name string, chunk *JSChunk, type_ AnalyzePrototypeProperties_SymbolType) {
	a := p.outer
	info := a.GetNameInfoForName(name, type_)
	var def *AnalyzePrototypeProperties_NameInfo
	// Skip all anonymous nodes. We care only about symbols with names.
	for i := len(p.symbolStack) - 1; i >= 0; i-- {
		context := p.symbolStack[i]
		def = context.name
		if def != a.anonymousNode {
			break
		}
	}
	if def == nil {
		panic(jutil.NewNullPointerException()) // def.equals(info) on a null def
	}
	if def != info {
		a.connectSymbolNodes(def, chunk, info)
	}
}

// ProcessNonFunctionPrototypeAssign: if this is a non-function prototype
// assign, return the prop name. Otherwise, return null.
//
// Java: AnalyzePrototypeProperties.ProcessProperties.processNonFunctionPrototypeAssign(Node, Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) ProcessNonFunctionPrototypeAssign(n *rhino.Node, parent *rhino.Node) *string {
	if s, ok := p.processNonFunctionPrototypeAssignString(n, parent); ok {
		return &s
	}
	return nil
}

func (p *AnalyzePrototypeProperties_ProcessProperties) processNonFunctionPrototypeAssignString(n *rhino.Node, parent *rhino.Node) (string, bool) {
	if p.IsAssignRValue(n, parent) && !n.IsFunction() {
		return p.getPrototypePropertyNameFromRValueString(n)
	}
	return "", false
}

// IsGlobalFunctionDeclaration determines whether n is the FUNCTION node in a
// global function declaration.
//
// Java: AnalyzePrototypeProperties.ProcessProperties.isGlobalFunctionDeclaration(NodeTraversal, Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) IsGlobalFunctionDeclaration(t *NodeTraversal, n *rhino.Node) bool {
	if !(NodeUtil_IsFunctionDeclaration(n) || (n.IsFunction() && n.GetParent().IsName())) {
		return false
	}

	// Make sure we're not in a function scope, or if we are then the function we're looking at
	// is defined in the global scope.
	return t.InGlobalHoistScope() ||
		(t.GetScopeRoot() == n &&
			t.GetScope().GetParent().GetClosestHoistScope().IsGlobal())
}

// IsAssignRValue returns true if this is the r-value of an assignment.
//
// Java: AnalyzePrototypeProperties.ProcessProperties.isAssignRValue(Node, Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) IsAssignRValue(n *rhino.Node, parent *rhino.Node) bool {
	return parent != nil && parent.IsAssign() && parent.GetFirstChild() != n
}

func bestLValueEndsWithPrototype(lv *rhino.Node) bool {
	if lv == nil || lv.GetParent() == nil {
		return false
	}
	if lv.GetParent().IsClassMembers() && !lv.IsComputedProp() {
		if lv.GetString() != "prototype" {
			return false
		}
		nameNode := NodeUtil_GetNameNode(lv.GetGrandparent())
		return nameNode != nil && nameNode.IsQualifiedName()
	}
	if NodeUtil_MayBeObjectLitKey(lv) {
		if NodeUtil_GetObjectOrClassLitKeyName(lv) != "prototype" {
			return false
		}
		owner := NodeUtil_GetBestLValue(lv.GetParent())
		if owner == nil {
			return false
		}
		_, ok := nodeUtil_getBestLValueNameString(owner)
		return ok
	}
	if lv.IsGetProp() {
		return lv.GetString() == "prototype" && lv.IsQualifiedName()
	}
	return false
}

// GetPrototypePropertyNameFromRValue returns the name of a prototype
// property being assigned to this r-value.
//
// Returns null if this is not the R-value of a prototype property, or if
// the R-value is used in multiple expressions (i.e., if there's a prototype
// property assignment in a more complex expression).
//
// Java: AnalyzePrototypeProperties.ProcessProperties.getPrototypePropertyNameFromRValue(Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) GetPrototypePropertyNameFromRValue(rValue *rhino.Node) *string {
	if s, ok := p.getPrototypePropertyNameFromRValueString(rValue); ok {
		return &s
	}
	return nil
}

func (p *AnalyzePrototypeProperties_ProcessProperties) getPrototypePropertyNameFromRValueString(rValue *rhino.Node) (string, bool) {
	lValue := NodeUtil_GetBestLValue(rValue)
	if lValue == nil ||
		lValue.GetParent() == nil ||
		!((NodeUtil_MayBeObjectLitKey(lValue) && !lValue.IsQuotedStringKey()) ||
			NodeUtil_IsExprAssign(lValue.GetGrandparent())) {
		return "", false
	}

	if lValue.GetParent().IsClassMembers() && !lValue.IsComputedProp() {
		nameNode := NodeUtil_GetNameNode(lValue.GetGrandparent())
		if nameNode == nil {
			return "", false
		}
		if !lValue.IsStaticMember() {
			if !nameNode.IsQualifiedName() {
				return "", false
			}
		} else {
			if !bestLValueEndsWithPrototype(nameNode) {
				return "", false
			}
		}
		return lValue.GetString(), true
	}
	if NodeUtil_MayBeObjectLitKey(lValue) {
		owner := NodeUtil_GetBestLValue(lValue.GetParent())
		if !bestLValueEndsWithPrototype(owner) {
			return "", false
		}
		key := NodeUtil_GetObjectOrClassLitKeyName(lValue)
		if !rhino.TokenStream_IsJSIdentifier(key) {
			return "", false
		}
		return key, true
	}
	if lValue.IsGetProp() {
		if !bestLValueEndsWithPrototype(lValue.GetFirstChild()) {
			return "", false
		}
		return lValue.GetString(), true
	}
	return "", false
}

// ProcessGlobalFunctionDeclaration processes a NAME node to see if it's a
// global function declaration. If it is, record it and return true.
// Otherwise, return false.
//
// Java: AnalyzePrototypeProperties.ProcessProperties.processGlobalFunctionDeclaration(NodeTraversal, Node, Var)
func (p *AnalyzePrototypeProperties_ProcessProperties) ProcessGlobalFunctionDeclaration(t *NodeTraversal, nameNode *rhino.Node, v *Var) bool {
	a := p.outer
	firstChild := nameNode.GetFirstChild()
	parent := nameNode.GetParent()

	if // Check for a named FUNCTION.
	p.IsGlobalFunctionDeclaration(t, parent) ||
		// Check for a VAR declaration.
		(firstChild != nil && p.IsGlobalFunctionDeclaration(t, firstChild)) {
		name := nameNode.GetString()
		a.GetNameInfoForName(name, AnalyzePrototypeProperties_VAR).
			GetDeclarations().
			Add(a.allocGlobalFunction(nameNode, v, t.GetChunk()))

		// If the function name is exported, we should create an edge here
		// so that it's never removed.
		if a.compiler.GetCodingConvention().IsExportedWithLocal(name, false /* local= */) ||
			a.anchorUnusedVars {
			p.AddGlobalUseOfSymbol(name, t.GetChunk(), AnalyzePrototypeProperties_VAR)
		}

		return true
	}
	return false
}

// HandlePossibleAssignmentToPrototype examines a qualified name ending in
// `.prototype`.
//
// If it is part of an assignment like `foo.prototype = {}` or
// `foo.prototype.bar = x`, record this reference as the definition of one or
// more prototype properties and return `true`.
//
// ref is a reference to some qualified name that ends with `.prototype`.
// Returns true if a declaration was added.
//
// Java: AnalyzePrototypeProperties.ProcessProperties.handlePossibleAssignmentToPrototype(NodeTraversal, Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) HandlePossibleAssignmentToPrototype(t *NodeTraversal, ref *rhino.Node) bool {
	a := p.outer
	root := NodeUtil_GetRootOfQualifiedName(ref)

	n := ref.GetParent()
	switch n.GetToken() {
	case rhino.Token_GETPROP:
		// Foo.prototype.getBar = function() { ... }
		parent := n.GetParent()
		grandParent := parent.GetParent()

		if NodeUtil_IsExprAssign(grandParent) &&
			NodeUtil_IsNameDeclOrSimpleAssignLhs(n, parent) {
			var prop AnalyzePrototypeProperties_PrototypeProperty = a.allocAssignmentPrototypeProperty(grandParent, p.MaybeGetVar(t, root), t.GetChunk())
			a.GetNameInfoForName(n.GetString(), AnalyzePrototypeProperties_PROPERTY).GetDeclarations().Add(prop)
			return true
		}
	case rhino.Token_ASSIGN:
		// Foo.prototype = { "getBar" : function() { ... } }
		map_ := n.GetSecondChild()
		if map_.IsObjectLit() {
			for key := map_.GetFirstChild(); key != nil; key = key.GetNext() {
				if !key.IsQuotedStringKey() && !key.IsComputedProp() {
					// We won't consider quoted or computed properties for any kind of modification,
					// so key may be STRING_KEY, GETTER_DEF, SETTER_DEF, or MEMBER_FUNCTION_DEF
					name := key.GetString()
					var prop AnalyzePrototypeProperties_PrototypeProperty = NewAnalyzePrototypeProperties_LiteralPrototypeProperty(
						key.GetFirstChild(), n, p.MaybeGetVar(t, root), t.GetChunk())
					a.GetNameInfoForName(name, AnalyzePrototypeProperties_PROPERTY).GetDeclarations().Add(prop)
				}
			}
			return true
		}
	default:
	}
	return false
}

// Java: AnalyzePrototypeProperties.ProcessProperties.processMemberDef(NodeTraversal, Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) ProcessMemberDef(t *NodeTraversal, n *rhino.Node) {
	jutil.CheckState(n.IsMemberFunctionDef() || n.IsGetterDef() || n.IsSetterDef())
	name := n.GetString()
	// Don't want to add a declaration for constructors and static members
	// so they aren't removed
	if NodeUtil_IsEs6ConstructorMemberFunctionDef(n) || n.IsStaticMember() {
		return
	}

	classNameNode := NodeUtil_GetNameNode(n.GetGrandparent())
	var var_ *Var
	if classNameNode != nil && classNameNode.IsName() {
		var_ = t.LookupVar(classNameNode.GetString())
	}
	p.outer.GetNameInfoForName(name, AnalyzePrototypeProperties_PROPERTY).
		GetDeclarations().
		Add(NewAnalyzePrototypeProperties_ClassMemberFunction(n, var_, t.GetChunk()))
}

// Java: AnalyzePrototypeProperties.ProcessProperties.maybeGetVar(NodeTraversal, Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) MaybeGetVar(t *NodeTraversal, maybeName *rhino.Node) *Var {
	if maybeName.IsName() {
		return t.LookupVar(maybeName.GetString())
	}
	return nil
}

// Java: AnalyzePrototypeProperties.ProcessProperties.addGlobalUseOfSymbol(String, JSChunk, SymbolType)
func (p *AnalyzePrototypeProperties_ProcessProperties) AddGlobalUseOfSymbol(name string, chunk *JSChunk, type_ AnalyzePrototypeProperties_SymbolType) {
	p.outer.connectSymbolNodes(p.outer.globalNode, chunk, p.outer.GetNameInfoForName(name, type_))
}

// AnalyzePrototypeProperties_ProcessExternProperties is the inner class
// ProcessExternProperties.
type AnalyzePrototypeProperties_ProcessExternProperties struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *AnalyzePrototypeProperties
}

// Java: AnalyzePrototypeProperties.ProcessExternProperties.visit(NodeTraversal, Node, Node)
func (p *AnalyzePrototypeProperties_ProcessExternProperties) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	a := p.outer
	if n.IsGetProp() {
		a.connectSymbolNodes(a.externNode, a.firstChunk, a.GetNameInfoForName(n.GetString(), AnalyzePrototypeProperties_PROPERTY))
	} else if n.IsMemberFunctionDef() || n.IsGetterDef() || n.IsSetterDef() {
		// As of 2019-08-29 the only user of this class is CrossChunkMethodMotion, which never
		// moves static methods, but that could change. So, we're intentionally including static
		// methods, static getters, and static setters here, because there are cases where they
		// could act like prototype properties.
		//
		// e.g.
		// // externs.js
		// class Foo {
		//   static foo() {}
		// }
		//
		// // src.js
		// /** @record */
		// class ObjWithFooMethod {
		//   foo() {}
		// }
		// /** @type {!ObjWithFooMethod} */
		// let objWithFooMethod = Foo; // yes, this is valid
		//
		a.connectSymbolNodes(a.externNode, a.firstChunk, a.GetNameInfoForName(n.GetString(), AnalyzePrototypeProperties_PROPERTY))
	}
}

// AnalyzePrototypeProperties_PropagateReferences is the inner class
// PropagateReferences.
type AnalyzePrototypeProperties_PropagateReferences struct {
	outer *AnalyzePrototypeProperties
}

// Java: AnalyzePrototypeProperties.PropagateReferences.traverseEdge(NameInfo, JSChunk, NameInfo)
func (p *AnalyzePrototypeProperties_PropagateReferences) TraverseEdge(start *AnalyzePrototypeProperties_NameInfo, edge *JSChunk, dest *AnalyzePrototypeProperties_NameInfo) bool {
	if start.IsReferenced() {
		startChunk := start.GetDeepestCommonChunkRef()
		if startChunk != nil && p.outer.chunkGraph.DependsOn(startChunk, edge) {
			return dest.MarkReference(startChunk)
		} else {
			return dest.MarkReference(edge)
		}
	}
	return false
}

// AnalyzePrototypeProperties_Symbol is the declaration of an abstract
// symbol.
type AnalyzePrototypeProperties_Symbol interface {
	// GetRootVar is the variable for the root of this symbol.
	//
	// Java: AnalyzePrototypeProperties.Symbol.getRootVar()
	GetRootVar() *Var

	// GetChunk returns the chunk where this appears.
	//
	// Java: AnalyzePrototypeProperties.Symbol.getChunk()
	GetChunk() *JSChunk
}

// AnalyzePrototypeProperties_SymbolType is the nested enum SymbolType.
type AnalyzePrototypeProperties_SymbolType int32

const (
	AnalyzePrototypeProperties_SymbolType_PROPERTY AnalyzePrototypeProperties_SymbolType = iota
	AnalyzePrototypeProperties_SymbolType_VAR
)

var analyzePrototypeProperties_SymbolType_names = [...]string{"PROPERTY", "VAR"}

// AnalyzePrototypeProperties_SymbolType_values is SymbolType.values().
func AnalyzePrototypeProperties_SymbolType_values() []AnalyzePrototypeProperties_SymbolType {
	return []AnalyzePrototypeProperties_SymbolType{AnalyzePrototypeProperties_SymbolType_PROPERTY, AnalyzePrototypeProperties_SymbolType_VAR}
}

// String is Enum.name().
func (t AnalyzePrototypeProperties_SymbolType) String() string {
	return analyzePrototypeProperties_SymbolType_names[t]
}

// AnalyzePrototypeProperties_GlobalFunction is a function initialized as a
// VAR statement or LET AND CONST and global or a function declaration.
type AnalyzePrototypeProperties_GlobalFunction struct {
	var_  *Var
	chunk *JSChunk
}

// Java: AnalyzePrototypeProperties.GlobalFunction(Node, Var, JSChunk)
func NewAnalyzePrototypeProperties_GlobalFunction(nameNode *rhino.Node, var_ *Var, chunk *JSChunk) *AnalyzePrototypeProperties_GlobalFunction {
	parent := nameNode.GetParent()
	jutil.CheckState(
		(NodeUtil_IsNameDeclaration(parent) && var_.IsGlobal()) ||
			NodeUtil_IsFunctionDeclaration(parent))
	return &AnalyzePrototypeProperties_GlobalFunction{var_: var_, chunk: chunk}
}

// Java: AnalyzePrototypeProperties.GlobalFunction.getRootVar()
func (g *AnalyzePrototypeProperties_GlobalFunction) GetRootVar() *Var {
	return g.var_
}

// Java: AnalyzePrototypeProperties.GlobalFunction.getChunk()
func (g *AnalyzePrototypeProperties_GlobalFunction) GetChunk() *JSChunk {
	return g.chunk
}

// AnalyzePrototypeProperties_Property represents property definitions done
// either directly on a '.prototype' object or as a member method, getter, or
// setter.
type AnalyzePrototypeProperties_Property interface {
	AnalyzePrototypeProperties_Symbol

	// Go-only: language-adapter — Property adds no method to Symbol in
	// Java; the marker keeps `symbol instanceof Property` false for a
	// GlobalFunction, since Go interfaces are satisfied structurally.
	isAnalyzePrototypePropertiesProperty()
}

// AnalyzePrototypeProperties_ClassMemberFunction is the nested class
// ClassMemberFunction.
type AnalyzePrototypeProperties_ClassMemberFunction struct {
	node  *rhino.Node
	var_  *Var
	chunk *JSChunk
}

// Java: AnalyzePrototypeProperties.ClassMemberFunction(Node, Var, JSChunk)
func NewAnalyzePrototypeProperties_ClassMemberFunction(node *rhino.Node, var_ *Var, chunk *JSChunk) *AnalyzePrototypeProperties_ClassMemberFunction {
	jutil.CheckState(node.GetParent().IsClassMembers())
	jutil.CheckState(node.IsMemberFunctionDef() || node.IsSetterDef() || node.IsGetterDef())
	return &AnalyzePrototypeProperties_ClassMemberFunction{node: node, var_: var_, chunk: chunk}
}

func (c *AnalyzePrototypeProperties_ClassMemberFunction) isAnalyzePrototypePropertiesProperty() {}

// Java: AnalyzePrototypeProperties.ClassMemberFunction.getRootVar()
func (c *AnalyzePrototypeProperties_ClassMemberFunction) GetRootVar() *Var {
	return c.var_
}

// Java: AnalyzePrototypeProperties.ClassMemberFunction.getChunk()
func (c *AnalyzePrototypeProperties_ClassMemberFunction) GetChunk() *JSChunk {
	return c.chunk
}

// GetFunctionNode returns the function node within the definition.
//
// Java: AnalyzePrototypeProperties.ClassMemberFunction.getFunctionNode()
func (c *AnalyzePrototypeProperties_ClassMemberFunction) GetFunctionNode() *rhino.Node {
	return c.node.GetOnlyChild()
}

// GetDefinitionNode returns the MEMBER_FUNCTION_DEF, GETTER_DEF, or
// SETTER_DEF node that defines the property.
//
// Java: AnalyzePrototypeProperties.ClassMemberFunction.getDefinitionNode()
func (c *AnalyzePrototypeProperties_ClassMemberFunction) GetDefinitionNode() *rhino.Node {
	return c.node
}

// AnalyzePrototypeProperties_PrototypeProperty represents a property that is
// defined on a '.prototype' object.
//
// This includes both of these cases.
//
//	Foo.prototype = {
//	      a: 1, // LiteralPrototypeProperty
//	      b: 2  // LiteralPrototypeProperty
//	    };
//	Bar.prototype.method = function() {}; // AssignmentPrototypeProperty.
//
// Since there are two ways of assigning properties to prototypes, we hide
// them behind this interface so they can both be removed regardless of
// type.
type AnalyzePrototypeProperties_PrototypeProperty interface {
	AnalyzePrototypeProperties_Property

	// GetPrototype returns the GETPROP node that refers to the prototype.
	//
	// Java: AnalyzePrototypeProperties.PrototypeProperty.getPrototype()
	GetPrototype() *rhino.Node

	// GetValue returns the value of this property.
	//
	// Java: AnalyzePrototypeProperties.PrototypeProperty.getValue()
	GetValue() *rhino.Node
}

// AnalyzePrototypeProperties_AssignmentPrototypeProperty: properties created
// via EXPR assignment:
//
//	function Foo() { ... };
//	Foo.prototype.bar = function() { ... };
type AnalyzePrototypeProperties_AssignmentPrototypeProperty struct {
	exprNode *rhino.Node
	rootVar  *Var
	chunk    *JSChunk
}

// NewAnalyzePrototypeProperties_AssignmentPrototypeProperty: node is an
// EXPR node.
//
// Java: AnalyzePrototypeProperties.AssignmentPrototypeProperty(Node, Var, JSChunk)
func NewAnalyzePrototypeProperties_AssignmentPrototypeProperty(node *rhino.Node, rootVar *Var, chunk *JSChunk) *AnalyzePrototypeProperties_AssignmentPrototypeProperty {
	return &AnalyzePrototypeProperties_AssignmentPrototypeProperty{exprNode: node, rootVar: rootVar, chunk: chunk}
}

func (p *AnalyzePrototypeProperties_AssignmentPrototypeProperty) isAnalyzePrototypePropertiesProperty() {
}

// Java: AnalyzePrototypeProperties.AssignmentPrototypeProperty.getRootVar()
func (p *AnalyzePrototypeProperties_AssignmentPrototypeProperty) GetRootVar() *Var {
	return p.rootVar
}

// Java: AnalyzePrototypeProperties.AssignmentPrototypeProperty.getPrototype()
func (p *AnalyzePrototypeProperties_AssignmentPrototypeProperty) GetPrototype() *rhino.Node {
	return p.GetAssignNode().GetFirstFirstChild()
}

// Java: AnalyzePrototypeProperties.AssignmentPrototypeProperty.getValue()
func (p *AnalyzePrototypeProperties_AssignmentPrototypeProperty) GetValue() *rhino.Node {
	return p.GetAssignNode().GetLastChild()
}

// Java: AnalyzePrototypeProperties.AssignmentPrototypeProperty.getAssignNode()
func (p *AnalyzePrototypeProperties_AssignmentPrototypeProperty) GetAssignNode() *rhino.Node {
	return p.exprNode.GetFirstChild()
}

// Java: AnalyzePrototypeProperties.AssignmentPrototypeProperty.getChunk()
func (p *AnalyzePrototypeProperties_AssignmentPrototypeProperty) GetChunk() *JSChunk {
	return p.chunk
}

// AnalyzePrototypeProperties_LiteralPrototypeProperty: properties created
// via object literals:
//
//	function Foo() { ... };
//	Foo.prototype = {bar: function() { ... };
type AnalyzePrototypeProperties_LiteralPrototypeProperty struct {
	value   *rhino.Node
	assign  *rhino.Node
	rootVar *Var
	chunk   *JSChunk
}

// Java: AnalyzePrototypeProperties.LiteralPrototypeProperty(Node, Node, Var, JSChunk)
func NewAnalyzePrototypeProperties_LiteralPrototypeProperty(value *rhino.Node, assign *rhino.Node, rootVar *Var, chunk *JSChunk) *AnalyzePrototypeProperties_LiteralPrototypeProperty {
	return &AnalyzePrototypeProperties_LiteralPrototypeProperty{value: value, assign: assign, rootVar: rootVar, chunk: chunk}
}

func (p *AnalyzePrototypeProperties_LiteralPrototypeProperty) isAnalyzePrototypePropertiesProperty() {
}

// Java: AnalyzePrototypeProperties.LiteralPrototypeProperty.getRootVar()
func (p *AnalyzePrototypeProperties_LiteralPrototypeProperty) GetRootVar() *Var {
	return p.rootVar
}

// Java: AnalyzePrototypeProperties.LiteralPrototypeProperty.getPrototype()
func (p *AnalyzePrototypeProperties_LiteralPrototypeProperty) GetPrototype() *rhino.Node {
	return p.assign.GetFirstChild()
}

// Java: AnalyzePrototypeProperties.LiteralPrototypeProperty.getValue()
func (p *AnalyzePrototypeProperties_LiteralPrototypeProperty) GetValue() *rhino.Node {
	return p.value
}

// Java: AnalyzePrototypeProperties.LiteralPrototypeProperty.getChunk()
func (p *AnalyzePrototypeProperties_LiteralPrototypeProperty) GetChunk() *JSChunk {
	return p.chunk
}

// AnalyzePrototypeProperties_NameContext is the context of the current
// name. This includes the NameInfo and the scope if it is a scope defining
// name (function).
type AnalyzePrototypeProperties_NameContext struct {
	name *AnalyzePrototypeProperties_NameInfo

	// If this is a function context, then scope will be the scope of the
	// corresponding function. Otherwise, it will be null.
	scope     *Scope
	scopeRoot *rhino.Node
}

// Java: AnalyzePrototypeProperties.NameContext(NameInfo, Scope)
func NewAnalyzePrototypeProperties_NameContext(name *AnalyzePrototypeProperties_NameInfo, scope *Scope) *AnalyzePrototypeProperties_NameContext {
	var root *rhino.Node
	if scope != nil {
		root = scope.GetRootNode()
	}
	return &AnalyzePrototypeProperties_NameContext{name: name, scope: scope, scopeRoot: root}
}

// AnalyzePrototypeProperties_NameInfo is information on all properties or
// global variables of a given name.
type AnalyzePrototypeProperties_NameInfo struct {
	outer *AnalyzePrototypeProperties

	name string

	referenced bool
	// Go-only: language-adapter — the ArrayDeque is only appended to (add)
	// and iterated in either direction, so it is an ArrayList.
	declarations          jutil.ArrayList[AnalyzePrototypeProperties_Symbol]
	declBuf               [2]AnalyzePrototypeProperties_Symbol
	deepestCommonChunkRef *JSChunk

	outEdges        []analyzeProtoPropEdge
	outEdgesBuf     [4]analyzeProtoPropEdge
	lastSrc         *AnalyzePrototypeProperties_NameInfo
	lastChunk       *JSChunk
	globalEdgeChunk *JSChunk
	externEdgeChunk *JSChunk
	hasGlobalEdge   bool
	hasExternEdge   bool
	seeded          bool
	inWorkSet       bool

	// True if this property is a function that reads a variable from an
	// outer scope which isn't the global scope.
	readClosureVariables bool

	// does the definition refer to `super`?
	// We cannot move references to `super` outside of the class body.
	referencesSuper bool
}

// NewAnalyzePrototypeProperties_NameInfo constructs a new NameInfo.
//
// name is the name of the property that this represents. May be null to
// signify dummy nodes in the property graph.
//
// Java: AnalyzePrototypeProperties.NameInfo(String)
func NewAnalyzePrototypeProperties_NameInfo(outer *AnalyzePrototypeProperties, name string) *AnalyzePrototypeProperties_NameInfo {
	if outer != nil {
		return outer.allocNameInfo(name)
	}
	info := &AnalyzePrototypeProperties_NameInfo{
		outer: outer,
		name:  name,
	}
	info.declarations.InitWithBuffer(info.declBuf[:])
	return info
}

// Java: AnalyzePrototypeProperties.NameInfo.toString()
func (i *AnalyzePrototypeProperties_NameInfo) String() string {
	return i.name
}

// IsReferenced determines whether we've marked a reference to this property
// name.
//
// Java: AnalyzePrototypeProperties.NameInfo.isReferenced()
func (i *AnalyzePrototypeProperties_NameInfo) IsReferenced() bool {
	return i.referenced
}

// ReadsClosureVariables determines whether it reads a closure variable.
//
// Java: AnalyzePrototypeProperties.NameInfo.readsClosureVariables()
func (i *AnalyzePrototypeProperties_NameInfo) ReadsClosureVariables() bool {
	return i.readClosureVariables
}

// ReferencesSuper: does the definition refer to `super`?
//
// Java: AnalyzePrototypeProperties.NameInfo.referencesSuper()
func (i *AnalyzePrototypeProperties_NameInfo) ReferencesSuper() bool {
	return i.referencesSuper
}

// MarkReference marks a reference in a given chunk to this property name,
// and record the deepest common chunk reference.
//
// chunk is the chunk where it was referenced. Returns whether the name info
// has changed.
//
// Java: AnalyzePrototypeProperties.NameInfo.markReference(JSChunk)
func (i *AnalyzePrototypeProperties_NameInfo) MarkReference(chunk *JSChunk) bool {
	hasChanged := false
	if !i.referenced {
		i.referenced = true
		hasChanged = true
	}

	originalDeepestCommon := i.deepestCommonChunkRef

	if i.deepestCommonChunkRef == nil {
		i.deepestCommonChunkRef = chunk
	} else {
		i.deepestCommonChunkRef =
			i.outer.chunkGraph.GetDeepestCommonDependencyInclusive(i.deepestCommonChunkRef, chunk)
	}

	if originalDeepestCommon != i.deepestCommonChunkRef {
		hasChanged = true
	}

	return hasChanged
}

// GetDeepestCommonChunkRef returns the deepest common chunk of all the
// references to this property.
//
// Java: AnalyzePrototypeProperties.NameInfo.getDeepestCommonChunkRef()
func (i *AnalyzePrototypeProperties_NameInfo) GetDeepestCommonChunkRef() *JSChunk {
	return i.deepestCommonChunkRef
}

// GetDeclarations returns a mutable collection of all the prototype property
// declarations of this property name.
//
// Java: AnalyzePrototypeProperties.NameInfo.getDeclarations()
func (i *AnalyzePrototypeProperties_NameInfo) GetDeclarations() *jutil.ArrayList[AnalyzePrototypeProperties_Symbol] {
	return &i.declarations
}
