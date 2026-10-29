// Port of com/google/javascript/jscomp/AnalyzePrototypeProperties.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
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

	// All the real NameInfo for prototype properties, keyed by property name ID
	// and recorded in insertion order in propertyInfos.
	propertyNameInfo map[uint32]*AnalyzePrototypeProperties_NameInfo
	propertyInfos    []*AnalyzePrototypeProperties_NameInfo

	// All the NameInfo for global functions, keyed by global variable name ID
	// and recorded in insertion order in varInfos.
	varNameInfo map[uint32]*AnalyzePrototypeProperties_NameInfo
	varInfos    []*AnalyzePrototypeProperties_NameInfo

	seededNodes []*AnalyzePrototypeProperties_NameInfo
	nameSlab    jutil.Slab[AnalyzePrototypeProperties_NameInfo]
}

var analyzeProtoPrototypeID = rhino.RhinoStringPool_InternID("prototype")

func (a *AnalyzePrototypeProperties) allocNameInfo(name string) *AnalyzePrototypeProperties_NameInfo {
	info := a.nameSlab.New()
	info.outer = a
	info.name = name
	info.declarations = info.declBuf[:0]
	info.outEdges = info.outEdgesBuf[:0]
	return info
}

func (a *AnalyzePrototypeProperties) allocNameInfoByID(id uint32) *AnalyzePrototypeProperties_NameInfo {
	info := a.nameSlab.New()
	info.outer = a
	info.nameID = id
	info.declarations = info.declBuf[:0]
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
	src.outEdges = append(src.outEdges, analyzeProtoPropEdge{chunk: chunk, dest: dest})
}

// Constants for symbol types, for easier readability.
const AnalyzePrototypeProperties_PROPERTY = AnalyzePrototypeProperties_SymbolType_PROPERTY
const AnalyzePrototypeProperties_VAR = AnalyzePrototypeProperties_SymbolType_VAR

// Properties that are implicitly used as part of the JS language.
var AnalyzePrototypeProperties_IMPLICITLY_USED_PROPERTIES = []string{"length", "toString", "valueOf"}

var analyzeExternPropTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_GETPROP)
	s.Add(rhino.Token_MEMBER_FUNCTION_DEF)
	s.Add(rhino.Token_GETTER_DEF)
	s.Add(rhino.Token_SETTER_DEF)
	return s
}()

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
		propertyNameInfo:       make(map[uint32]*AnalyzePrototypeProperties_NameInfo, 16),
		varNameInfo:            make(map[uint32]*AnalyzePrototypeProperties_NameInfo, 16),
		compiler:               compiler,
		chunkGraph:             chunkGraph,
		canModifyExterns:       canModifyExterns,
		anchorUnusedVars:       anchorUnusedVars,
		rootScopeUsesAreGlobal: rootScopeUsesAreGlobal,
	}
	a.globalNode = a.allocNameInfo("[global]")
	a.externNode = a.allocNameInfo("[extern]")
	a.anonymousNode = a.allocNameInfo("[anonymous]")

	if chunkGraph != nil && chunkGraph.GetChunkCount() > 1 {
		a.firstChunk = chunkGraph.GetRootChunk()
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
func (a *AnalyzePrototypeProperties) Process(externRoot rhino.Node, root rhino.Node) {
	jutil.CheckState(a.compiler.GetLifeCycleStage().IsNormalized())
	if !a.canModifyExterns {
		NodeTraversal_Traverse(a.compiler, externRoot, &analyzeProcessExternProperties{outer: a})
	}

	var propProc AnalyzePrototypeProperties_ProcessProperties
	propProc.outer = a
	propProc.symbolStack = propProc.symbolStackBuf[:0]
	NodeTraversal_Builder_().SetCompiler(a.compiler).SetCallback(&propProc).Traverse(root)

	a.propagateReferences()
}

func (a *AnalyzePrototypeProperties) propagateReferences() {
	if len(a.seededNodes) == 0 {
		return
	}
	totalNodes := len(a.propertyInfos) + len(a.varInfos) + 2
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

// GetAllNameInfo returns information on all prototype properties and global functions.
//
// Java: AnalyzePrototypeProperties.getAllNameInfo()
func (a *AnalyzePrototypeProperties) GetAllNameInfo() []*AnalyzePrototypeProperties_NameInfo {
	result := make([]*AnalyzePrototypeProperties_NameInfo, 0, len(a.propertyInfos)+len(a.varInfos))
	result = append(result, a.propertyInfos...)
	result = append(result, a.varInfos...)
	return result
}

func (a *AnalyzePrototypeProperties) getNameInfoByID(id uint32, type_ AnalyzePrototypeProperties_SymbolType) *AnalyzePrototypeProperties_NameInfo {
	if type_ == AnalyzePrototypeProperties_PROPERTY {
		if info := a.propertyNameInfo[id]; info != nil {
			return info
		}
		info := a.allocNameInfoByID(id)
		a.propertyNameInfo[id] = info
		a.propertyInfos = append(a.propertyInfos, info)
		return info
	}
	if info := a.varNameInfo[id]; info != nil {
		return info
	}
	info := a.allocNameInfoByID(id)
	a.varNameInfo[id] = info
	a.varInfos = append(a.varInfos, info)
	return info
}

// GetNameInfoForName gets the name info for the property or variable of a
// given name, and creates a new one if necessary.
//
// Java: AnalyzePrototypeProperties.getNameInfoForName(String, SymbolType)
func (a *AnalyzePrototypeProperties) GetNameInfoForName(name string, type_ AnalyzePrototypeProperties_SymbolType) *AnalyzePrototypeProperties_NameInfo {
	info := a.getNameInfoByID(rhino.RhinoStringPool_InternID(name), type_)
	if info.name == "" && name != "" {
		info.name = name
	}
	return info
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
	symbolStack    []AnalyzePrototypeProperties_NameContext
	symbolStackBuf [16]AnalyzePrototypeProperties_NameContext
}

func (p *AnalyzePrototypeProperties_ProcessProperties) push(c AnalyzePrototypeProperties_NameContext) {
	p.symbolStack = append(p.symbolStack, c)
}

func (p *AnalyzePrototypeProperties_ProcessProperties) pop() {
	if len(p.symbolStack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	p.symbolStack = p.symbolStack[:len(p.symbolStack)-1]
}

// Java: AnalyzePrototypeProperties.ProcessProperties.enterScope(NodeTraversal)
func (p *AnalyzePrototypeProperties_ProcessProperties) EnterScope(t *NodeTraversal) {
	a := p.outer
	root := t.GetScopeRoot()
	if root.IsFunction() {
		n := t.GetCurrentNode()
		if propID, ok := p.getPrototypePropertyNameFromRValueID(n); ok {
			p.push(AnalyzePrototypeProperties_NameContext{
				name:      a.getNameInfoByID(propID, AnalyzePrototypeProperties_PROPERTY),
				scopeRoot: root,
			})
		} else if p.IsGlobalFunctionDeclaration(t, n) {
			parent := n.GetParent()
			var nameID uint32
			if parent.IsName() {
				nameID = parent.GetStringID() /* VAR */
			} else {
				nameID = n.GetFirstChild().GetStringID() /* named function */
			}
			p.push(AnalyzePrototypeProperties_NameContext{
				name:      a.getNameInfoByID(nameID, AnalyzePrototypeProperties_VAR),
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
func (p *AnalyzePrototypeProperties_ProcessProperties) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if !p.outer.rootScopeUsesAreGlobal {
		// Process assignment of a non-function expression to a prototype property.
		if propID, ok := p.processNonFunctionPrototypeAssignID(n, parent); ok {
			p.push(AnalyzePrototypeProperties_NameContext{
				name: p.outer.getNameInfoByID(propID, AnalyzePrototypeProperties_PROPERTY),
			})
		}
	}
	return true
}

// Java: AnalyzePrototypeProperties.ProcessProperties.visit(NodeTraversal, Node, Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
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
		p.addSymbolUseByID(n.GetStringID(), t.GetChunk(), AnalyzePrototypeProperties_PROPERTY)

	case rhino.Token_GETPROP:
		propID := n.GetStringID()

		if n.IsQualifiedName() {
			if propID == analyzeProtoPrototypeID {
				if p.HandlePossibleAssignmentToPrototype(t, n) {
					// The reference is being assigned to not read from, so don't record this as a
					// reference.
					return
				}
			} else if a.compiler.GetCodingConvention().IsExportedWithLocal(rhino.RhinoStringPool_StringFromID(propID), false /* local= */) {
				p.addGlobalUseOfSymbolByID(propID, t.GetChunk(), AnalyzePrototypeProperties_PROPERTY)
				return
			} else if assign, ok := ast.AsAssign(parent); ok && assign.Target().Is(n) {
				if _, ok := p.getPrototypePropertyNameFromRValueID(n); ok {
					// e.g. `Foo.prototype.bar = something`
					// We record the declaration of `bar` when we look at the `Foo.prototype`
					// Node via the call to handlePossibleAssignmentToPrototype() above.
					// Now that we're looking at the whole `Foo.prototype.bar` node,
					// we just need to make sure we don't record it as a read reference.
					return
				}
			}
		}

		p.addSymbolUseByID(propID, t.GetChunk(), AnalyzePrototypeProperties_PROPERTY)

	case rhino.Token_OBJECTLIT:
		// Make sure that we're not handling object literals being
		// assigned to a prototype, as in:
		// Foo.prototype = {bar: 3, baz: 5};
		if bestLValueEndsWithPrototype(NodeUtil_GetBestLValue(n)) {
			return
		}
		fallthrough

	case rhino.Token_OBJECT_PATTERN:
		// `var x = {a: 1, b: 2}` and `var {a: x, b: y} = obj;`
		// should count as a use of property a and b.
		for propNode := n.GetFirstChild(); !propNode.IsNil(); propNode = propNode.GetNext() {
			if pkm, ok := ast.AsPlainKeyMember(propNode); ok {
				if !pkm.IsQuoted() {
					// May be STRING, GET, or SET, but NUMBER isn't interesting.
					p.addSymbolUseByID(propNode.GetStringID(), t.GetChunk(), AnalyzePrototypeProperties_PROPERTY)
				}
				continue
			}
			switch propNode.GetToken() {
			case rhino.Token_COMPUTED_PROP,
				rhino.Token_ITER_REST,
				rhino.Token_OBJECT_REST,
				rhino.Token_ITER_SPREAD,
				rhino.Token_OBJECT_SPREAD:
			default:
				panic(jutil.NewIllegalStateExceptionMsg(
					"Unexpected child of " + n.GetToken().String() + ": " + propNode.ToStringTree()))
			}
		}

	case rhino.Token_CLASS:
		for member := range ast.MustClass(n).Body().Members() {
			if methodDef, ok := member.AsMethodDef(); ok {
				p.processMemberDef(t, methodDef)
			}
		}

	case rhino.Token_NAME:
		nameID := n.GetStringID()
		if nameID == 0 {
			break
		}

		var_ := t.LookupVarByID(nameID)
		if var_ != nil {
			// Only process global functions.
			if var_.IsGlobal() {
				if !var_.GetInitialValue().IsNil() && var_.GetInitialValue().IsFunction() {
					if t.InGlobalHoistScope() {
						if !p.ProcessGlobalFunctionDeclaration(t, n, var_) {
							p.addGlobalUseOfSymbolByID(nameID, t.GetChunk(), AnalyzePrototypeProperties_VAR)
						}
					} else {
						p.addSymbolUseByID(nameID, t.GetChunk(), AnalyzePrototypeProperties_VAR)
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
	if !a.rootScopeUsesAreGlobal {
		if _, ok := p.processNonFunctionPrototypeAssignID(n, parent); ok {
			p.pop()
		}
	}
}

func (p *AnalyzePrototypeProperties_ProcessProperties) addSymbolUseByID(id uint32, chunk *JSChunk, type_ AnalyzePrototypeProperties_SymbolType) {
	a := p.outer
	info := a.getNameInfoByID(id, type_)
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
		panic(jutil.NewNullPointerException())
	}
	if def != info {
		a.connectSymbolNodes(def, chunk, info)
	}
}

// Java: AnalyzePrototypeProperties.ProcessProperties.addSymbolUse(String, JSChunk, SymbolType)
func (p *AnalyzePrototypeProperties_ProcessProperties) AddSymbolUse(name string, chunk *JSChunk, type_ AnalyzePrototypeProperties_SymbolType) {
	p.addSymbolUseByID(rhino.RhinoStringPool_InternID(name), chunk, type_)
}

func (p *AnalyzePrototypeProperties_ProcessProperties) processNonFunctionPrototypeAssignID(n rhino.Node, parent rhino.Node) (uint32, bool) {
	if isAssignRValue(n, parent) && !n.IsFunction() {
		return p.getPrototypePropertyNameFromRValueID(n)
	}
	return 0, false
}

// IsGlobalFunctionDeclaration determines whether n is the FUNCTION node in a
// global function declaration.
//
// Java: AnalyzePrototypeProperties.ProcessProperties.isGlobalFunctionDeclaration(NodeTraversal, Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) IsGlobalFunctionDeclaration(t *NodeTraversal, n rhino.Node) bool {
	if !(NodeUtil_IsFunctionDeclaration(n) || (n.IsFunction() && n.GetParent().IsName())) {
		return false
	}

	// Make sure we're not in a function scope, or if we are then the function we're looking at
	// is defined in the global scope.
	return t.InGlobalHoistScope() ||
		(t.GetScopeRoot() == n &&
			t.GetScope().GetParent().GetClosestHoistScope().IsGlobal())
}

func isAssignRValue(n rhino.Node, parent rhino.Node) bool {
	assign, ok := ast.AsAssign(parent)
	return ok && !assign.Target().Is(n)
}

func bestLValueEndsWithPrototype(lv rhino.Node) bool {
	if lv.IsNil() || lv.GetParent().IsNil() {
		return false
	}
	if lv.GetParent().IsClassMembers() && !lv.IsComputedProp() {
		if lv.GetStringID() != analyzeProtoPrototypeID {
			return false
		}
		nameNode := NodeUtil_GetNameNode(lv.GetGrandparent())
		return !nameNode.IsNil() && nameNode.IsQualifiedName()
	}
	if NodeUtil_MayBeObjectLitKey(lv) {
		if !lv.IsComputedProp() {
			if lv.GetStringID() != analyzeProtoPrototypeID {
				return false
			}
		} else if NodeUtil_GetObjectOrClassLitKeyName(lv) != "prototype" {
			return false
		}
		owner := NodeUtil_GetBestLValue(lv.GetParent())
		if owner.IsNil() {
			return false
		}
		_, ok := nodeUtil_getBestLValueNameString(owner)
		return ok
	}
	if lv.IsGetProp() {
		return lv.GetStringID() == analyzeProtoPrototypeID && lv.IsQualifiedName()
	}
	return false
}

func (p *AnalyzePrototypeProperties_ProcessProperties) getPrototypePropertyNameFromRValueID(rValue rhino.Node) (uint32, bool) {
	lValue := NodeUtil_GetBestLValue(rValue)
	if lValue.IsNil() ||
		lValue.GetParent().IsNil() ||
		!((NodeUtil_MayBeObjectLitKey(lValue) && !lValue.IsQuotedStringKey()) ||
			NodeUtil_IsExprAssign(lValue.GetGrandparent())) {
		return 0, false
	}

	if lValue.GetParent().IsClassMembers() && !lValue.IsComputedProp() {
		nameNode := NodeUtil_GetNameNode(lValue.GetGrandparent())
		if nameNode.IsNil() {
			return 0, false
		}
		if !lValue.IsStaticMember() {
			if !nameNode.IsQualifiedName() {
				return 0, false
			}
		} else {
			if !bestLValueEndsWithPrototype(nameNode) {
				return 0, false
			}
		}
		return lValue.GetStringID(), true
	}
	if NodeUtil_MayBeObjectLitKey(lValue) {
		owner := NodeUtil_GetBestLValue(lValue.GetParent())
		if !bestLValueEndsWithPrototype(owner) {
			return 0, false
		}
		key := NodeUtil_GetObjectOrClassLitKeyName(lValue)
		if !rhino.TokenStream_IsJSIdentifier(key) {
			return 0, false
		}
		if !lValue.IsComputedProp() {
			return lValue.GetStringID(), true
		}
		return rhino.RhinoStringPool_InternID(key), true
	}
	if gp, ok := ast.AsGetProp(lValue); ok {
		if !bestLValueEndsWithPrototype(gp.Object().Node()) {
			return 0, false
		}
		return lValue.GetStringID(), true
	}
	return 0, false
}

// ProcessGlobalFunctionDeclaration processes a NAME node to see if it's a
// global function declaration. If it is, record it and return true.
// Otherwise, return false.
//
// Java: AnalyzePrototypeProperties.ProcessProperties.processGlobalFunctionDeclaration(NodeTraversal, Node, Var)
func (p *AnalyzePrototypeProperties_ProcessProperties) ProcessGlobalFunctionDeclaration(t *NodeTraversal, nameNode rhino.Node, v *Var) bool {
	a := p.outer
	nameView := ast.MustName(nameNode)
	parent := nameNode.GetParent()
	initExpr, hasInit := nameView.Init()

	if // Check for a named FUNCTION.
	p.IsGlobalFunctionDeclaration(t, parent) ||
		// Check for a VAR declaration.
		(hasInit && p.IsGlobalFunctionDeclaration(t, initExpr.Node())) {
		jutil.CheckState(
			(NodeUtil_IsNameDeclaration(parent) && v.IsGlobal()) ||
				NodeUtil_IsFunctionDeclaration(parent))
		nameID := nameNode.GetStringID()
		a.getNameInfoByID(nameID, AnalyzePrototypeProperties_VAR).AddDeclaration(
			AnalyzePrototypeProperties_Symbol{
				kind:    AnalyzePrototypeProperties_SymbolGlobalFunction,
				rootVar: v,
				chunk:   t.GetChunk(),
			})

		// If the function name is exported, we should create an edge here
		// so that it's never removed.
		if a.anchorUnusedVars ||
			a.compiler.GetCodingConvention().IsExportedWithLocal(rhino.RhinoStringPool_StringFromID(nameID), false /* local= */) {
			p.addGlobalUseOfSymbolByID(nameID, t.GetChunk(), AnalyzePrototypeProperties_VAR)
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
// Java: AnalyzePrototypeProperties.ProcessProperties.handlePossibleAssignmentToPrototype(NodeTraversal, Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) HandlePossibleAssignmentToPrototype(t *NodeTraversal, ref rhino.Node) bool {
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
			a.getNameInfoByID(n.GetStringID(), AnalyzePrototypeProperties_PROPERTY).AddDeclaration(
				AnalyzePrototypeProperties_Symbol{
					kind:    AnalyzePrototypeProperties_SymbolAssignmentProperty,
					node:    grandParent,
					rootVar: p.MaybeGetVar(t, root),
					chunk:   t.GetChunk(),
				})
			return true
		}
	case rhino.Token_ASSIGN:
		// Foo.prototype = { "getBar" : function() { ... } }
		assign := ast.MustAssign(n)
		if objLit, ok := assign.Value().AsObjectLit(); ok {
			rootVar := p.MaybeGetVar(t, root)
			chunk := t.GetChunk()
			for member := range objLit.Members() {
				if pkm, ok := member.AsPlainKeyMember(); ok && !pkm.IsQuoted() {
					// We won't consider quoted or computed properties for any kind of modification,
					// so member may be STRING_KEY, GETTER_DEF, SETTER_DEF, or MEMBER_FUNCTION_DEF.
					a.getNameInfoByID(pkm.Node().GetStringID(), AnalyzePrototypeProperties_PROPERTY).AddDeclaration(
						AnalyzePrototypeProperties_Symbol{
							kind:       AnalyzePrototypeProperties_SymbolLiteralProperty,
							node:       pkm.Value().Node(),
							assignNode: n,
							rootVar:    rootVar,
							chunk:      chunk,
						})
				}
			}
			return true
		}
	default:
	}
	return false
}

func (p *AnalyzePrototypeProperties_ProcessProperties) processMemberDef(t *NodeTraversal, methodDef ast.MethodDef) {
	n := methodDef.Node()
	// Don't want to add a declaration for constructors and static members
	// so they aren't removed.
	if NodeUtil_IsEs6ConstructorMemberFunctionDef(n) || methodDef.IsStatic() {
		return
	}

	classNameNode := NodeUtil_GetNameNode(n.GetGrandparent())
	var var_ *Var
	if classNameNode.IsName() {
		var_ = t.LookupVarByID(classNameNode.GetStringID())
	}
	p.outer.getNameInfoByID(n.GetStringID(), AnalyzePrototypeProperties_PROPERTY).AddDeclaration(
		AnalyzePrototypeProperties_Symbol{
			kind:    AnalyzePrototypeProperties_SymbolClassMember,
			node:    n,
			rootVar: var_,
			chunk:   t.GetChunk(),
		})
}

// Java: AnalyzePrototypeProperties.ProcessProperties.maybeGetVar(NodeTraversal, Node)
func (p *AnalyzePrototypeProperties_ProcessProperties) MaybeGetVar(t *NodeTraversal, maybeName rhino.Node) *Var {
	if maybeName.IsName() {
		return t.LookupVarByID(maybeName.GetStringID())
	}
	return nil
}

func (p *AnalyzePrototypeProperties_ProcessProperties) addGlobalUseOfSymbolByID(id uint32, chunk *JSChunk, type_ AnalyzePrototypeProperties_SymbolType) {
	p.outer.connectSymbolNodes(p.outer.globalNode, chunk, p.outer.getNameInfoByID(id, type_))
}

// Java: AnalyzePrototypeProperties.ProcessProperties.addGlobalUseOfSymbol(String, JSChunk, SymbolType)
func (p *AnalyzePrototypeProperties_ProcessProperties) AddGlobalUseOfSymbol(name string, chunk *JSChunk, type_ AnalyzePrototypeProperties_SymbolType) {
	p.addGlobalUseOfSymbolByID(rhino.RhinoStringPool_InternID(name), chunk, type_)
}

type analyzeProcessExternProperties struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *AnalyzePrototypeProperties
}

func (p *analyzeProcessExternProperties) visitTokens() *ast.TokenSet {
	return &analyzeExternPropTokens
}

// Java: AnalyzePrototypeProperties.ProcessExternProperties.visit(NodeTraversal, Node, Node)
func (p *analyzeProcessExternProperties) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETPROP,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF:
		a := p.outer
		// Includes GETPROP as well as static/instance methods, getters, and setters
		// in extern classes, which could act like prototype properties via @record.
		a.connectSymbolNodes(a.externNode, a.firstChunk, a.getNameInfoByID(n.GetStringID(), AnalyzePrototypeProperties_PROPERTY))
	}
}

// AnalyzePrototypeProperties_SymbolType is the nested enum SymbolType.
type AnalyzePrototypeProperties_SymbolType int32

const (
	AnalyzePrototypeProperties_SymbolType_PROPERTY AnalyzePrototypeProperties_SymbolType = iota
	AnalyzePrototypeProperties_SymbolType_VAR
)

var analyzePrototypeProperties_SymbolType_names = [...]string{"PROPERTY", "VAR"}

// String is Enum.name().
func (t AnalyzePrototypeProperties_SymbolType) String() string {
	return analyzePrototypeProperties_SymbolType_names[t]
}

// AnalyzePrototypeProperties_SymbolKind identifies the kind of symbol declaration.
type AnalyzePrototypeProperties_SymbolKind uint8

const (
	AnalyzePrototypeProperties_SymbolGlobalFunction AnalyzePrototypeProperties_SymbolKind = iota
	AnalyzePrototypeProperties_SymbolClassMember
	AnalyzePrototypeProperties_SymbolAssignmentProperty
	AnalyzePrototypeProperties_SymbolLiteralProperty
)

// AnalyzePrototypeProperties_Symbol is a declaration of a global function,
// class member method/accessor, or `.prototype` property.
type AnalyzePrototypeProperties_Symbol struct {
	rootVar    *Var
	chunk      *JSChunk
	node       rhino.Node
	assignNode rhino.Node
	kind       AnalyzePrototypeProperties_SymbolKind
}

func (s AnalyzePrototypeProperties_Symbol) Kind() AnalyzePrototypeProperties_SymbolKind {
	return s.kind
}

func (s AnalyzePrototypeProperties_Symbol) IsProperty() bool {
	return s.kind != AnalyzePrototypeProperties_SymbolGlobalFunction
}

// Java: AnalyzePrototypeProperties.Symbol.getRootVar()
func (s AnalyzePrototypeProperties_Symbol) GetRootVar() *Var {
	return s.rootVar
}

// Java: AnalyzePrototypeProperties.Symbol.getChunk()
func (s AnalyzePrototypeProperties_Symbol) GetChunk() *JSChunk {
	return s.chunk
}

// GetDefinitionNode returns the MEMBER_FUNCTION_DEF, GETTER_DEF, or
// SETTER_DEF node that defines a ClassMember symbol.
func (s AnalyzePrototypeProperties_Symbol) GetDefinitionNode() rhino.Node {
	return s.node
}

// GetPrototype returns the GETPROP node that refers to `.prototype` for an
// AssignmentProperty or LiteralProperty symbol.
func (s AnalyzePrototypeProperties_Symbol) GetPrototype() rhino.Node {
	if s.kind == AnalyzePrototypeProperties_SymbolAssignmentProperty {
		return ast.MustExprStmt(s.node).Expr().MustAssign().Target().MustGetProp().Object().Node()
	}
	return ast.MustAssign(s.assignNode).Target().Node()
}

// GetValue returns the value node of an AssignmentProperty or LiteralProperty symbol.
func (s AnalyzePrototypeProperties_Symbol) GetValue() rhino.Node {
	if s.kind == AnalyzePrototypeProperties_SymbolAssignmentProperty {
		return ast.MustExprStmt(s.node).Expr().MustAssign().Value().Node()
	}
	return s.node
}

// AnalyzePrototypeProperties_NameContext is the context of the current
// name. This includes the NameInfo and the scope root if it is a scope-defining
// name (function).
type AnalyzePrototypeProperties_NameContext struct {
	name      *AnalyzePrototypeProperties_NameInfo
	scopeRoot rhino.Node
}

// AnalyzePrototypeProperties_NameInfo is information on all properties or
// global variables of a given name.
type AnalyzePrototypeProperties_NameInfo struct {
	outer *AnalyzePrototypeProperties

	name   string
	nameID uint32

	referenced            bool
	declarations          []AnalyzePrototypeProperties_Symbol
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

// Java: AnalyzePrototypeProperties.NameInfo.toString()
func (i *AnalyzePrototypeProperties_NameInfo) String() string {
	if i.name == "" && i.nameID != 0 {
		i.name = rhino.RhinoStringPool_StringFromID(i.nameID)
	}
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
// and records the deepest common chunk reference.
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

// GetDeclarations returns all the prototype property declarations of this
// property name.
//
// Java: AnalyzePrototypeProperties.NameInfo.getDeclarations()
func (i *AnalyzePrototypeProperties_NameInfo) GetDeclarations() []AnalyzePrototypeProperties_Symbol {
	return i.declarations
}

// AddDeclaration appends a declaration symbol to this NameInfo.
func (i *AnalyzePrototypeProperties_NameInfo) AddDeclaration(sym AnalyzePrototypeProperties_Symbol) {
	i.declarations = append(i.declarations, sym)
}
