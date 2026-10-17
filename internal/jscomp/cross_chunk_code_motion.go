// Port of com/google/javascript/jscomp/CrossChunkCodeMotion.java @ closure-compiler v20261015

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
	"reflect"

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CrossChunkCodeMotion is a compiler pass for moving global variable
// declarations and assignments to their properties to a deeper chunk if
// possible.
//
//  1. Collect all global-level statements and the references they contain.
//  2. Statements that do not declare global variables are assumed to exist
//     for their side-effects and are considered immovable.
//  3. Statements that declare global variables may be movable. See
//     CrossChunkReferenceCollector.
//  4. Within each chunk gather all declarations for a single global
//     variable into a DeclarationStatementGroup (DSG). Keep track of the
//     references to other globals that appear in the DSG. A DSG is movable
//     if all of its statements are movable.
//  5. Gather the DSGs for each global variable and note all of the chunks
//     that contain immovable references to it.
//  6. The global variables form a directed graph. Global A has an edge to B
//     if A has a DSG with a reference to B.
//  7. Convert this to a directed-acyclic graph by grouping together all of
//     the strongly-connected global variables into GlobalSymbolCycles. There
//     is an edge from GlobalSymbolCycle A to GlobalSymbolCycle B if A
//     contains a DSG (in one of its GlobalSymbols) that contains a reference
//     to a GlobalSymbol in B.
//  8. Sort the GlobalSymbolCycles into a list c[1], c[2], c[3],..., c[n],
//     such that there are no references to GlobalSymbols in c[i] from DSGs
//     in GlobalSymbolCycles c[i+1]...c[n].
//  9. Traverse the list in that order, so statements for GlobalSymbol X will
//     only be moved after all statements that refer to X have already been
//     moved.
//  10. Within a GlobalSymbolCycle, combine statements from DSGs in the same
//     chunk and keep them in the same order when moving them.
type CrossChunkCodeMotion struct {
	// Allocated & cleaned up by process()
	cccmLog diagnostic.LogFile

	compiler AbstractCompiler
	graph    *JSChunkGraph

	// Map from chunk to the node in that chunk that should parent variable
	// declarations that have to be moved into that chunk
	chunkInsertionPointMap *jutil.LinkedHashMap[*JSChunk, *rhino.Node]

	parentChunkCanSeeSymbolsDeclaredInChildren bool
}

// NewCrossChunkCodeMotion creates an instance.
//
// compiler is the compiler.
//
// Java: CrossChunkCodeMotion(AbstractCompiler, JSChunkGraph, boolean)
func NewCrossChunkCodeMotion(
	compiler AbstractCompiler,
	graph *JSChunkGraph,
	parentChunkCanSeeSymbolsDeclaredInChildren bool) *CrossChunkCodeMotion {
	return &CrossChunkCodeMotion{
		chunkInsertionPointMap: jutil.NewLinkedHashMap[*JSChunk, *rhino.Node](),
		compiler:               compiler,
		graph:                  graph,
		parentChunkCanSeeSymbolsDeclaredInChildren: parentChunkCanSeeSymbolsDeclaredInChildren,
	}
}

// Java: CrossChunkCodeMotion.process(Node, Node)
func (c *CrossChunkCodeMotion) Process(externs *rhino.Node, root *rhino.Node) {
	logFile := c.compiler.CreateOrReopenIndexedLog(reflect.TypeFor[CrossChunkCodeMotion](), "cccm.log")
	defer logFile.Close()
	defer func() { c.cccmLog = nil }()
	c.cccmLog = logFile
	c.cccmLog.LogString("doing cccm")
	// If there are <2 chunks, then we will never move anything, so we're done
	if c.graph.GetChunkCount() > 1 {
		referenceCollector :=
			NewCrossChunkReferenceCollector(c.compiler, NewSyntacticScopeCreator(c.compiler))
		referenceCollector.ProcessRoot(root)
		globalSymbols :=
			NewCrossChunkCodeMotion_GlobalSymbolCollector(c).CollectGlobalSymbols(referenceCollector)
		c.MoveGlobalSymbols(globalSymbols)
		c.AddInstanceofGuards(globalSymbols)
	} else {
		c.cccmLog.LogString("only one chunk exists")
	}
}

// Java: CrossChunkCodeMotion.addInstanceofGuards(Collection<GlobalSymbol>)
func (c *CrossChunkCodeMotion) AddInstanceofGuards(globalSymbols *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol]) {
	for globalSymbol := range globalSymbols.All() {
		for instanceofReference := range globalSymbol.instanceofReferencesToGuard.All() {
			if !globalSymbol.DeclarationsCoverChunk(instanceofReference.GetChunk()) {
				c.AddGuardToInstanceofReference(instanceofReference.GetReference().GetNode())
			}
		}
	}
}

// CrossChunkCodeMotion_GlobalSymbolCollector collects all global symbols,
// their declaration statements and references.
type CrossChunkCodeMotion_GlobalSymbolCollector struct {
	outer *CrossChunkCodeMotion

	globalSymbolforVar *jutil.ObjectLinkedHashMap[*Var, *CrossChunkCodeMotion_GlobalSymbol]

	// Returning the symbols in the reverse order in which they are defined
	// helps to minimize unnecessary reordering of declaration statements.
	symbolStack *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol]
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector()
func NewCrossChunkCodeMotion_GlobalSymbolCollector(outer *CrossChunkCodeMotion) *CrossChunkCodeMotion_GlobalSymbolCollector {
	return &CrossChunkCodeMotion_GlobalSymbolCollector{
		outer:              outer,
		globalSymbolforVar: jutil.NewObjectLinkedHashMap[*Var, *CrossChunkCodeMotion_GlobalSymbol](),
		symbolStack:        jutil.NewArrayDeque[*CrossChunkCodeMotion_GlobalSymbol](),
	}
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.collectGlobalSymbols(CrossChunkReferenceCollector)
func (g *CrossChunkCodeMotion_GlobalSymbolCollector) CollectGlobalSymbols(referenceCollector *CrossChunkReferenceCollector) *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol] {

	for _, statement := range referenceCollector.GetTopLevelStatements() {
		if statement.IsDeclarationStatement() {
			g.ProcessDeclarationStatement(statement)
		} else {
			g.ProcessImmovableStatement(statement)
		}
	}
	return g.symbolStack
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.processImmovableStatement(TopLevelStatement)
func (g *CrossChunkCodeMotion_GlobalSymbolCollector) ProcessImmovableStatement(statement *CrossChunkReferenceCollector_TopLevelStatement) {
	for _, ref := range statement.GetNonDeclarationReferences() {
		g.ProcessImmovableReference(ref, statement.GetChunk())
	}
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.processImmovableReference(Reference, JSChunk)
func (g *CrossChunkCodeMotion_GlobalSymbolCollector) ProcessImmovableReference(ref *Reference, chunk *JSChunk) {
	globalSymbol := g.GetGlobalSymbol(ref.getSymbolV())
	if g.outer.parentChunkCanSeeSymbolsDeclaredInChildren {
		// It is possible to move the declaration of `Foo` after
		// `'undefined' != typeof Foo && x instanceof Foo`.
		// We'll add the undefined check, if necessary.
		n := ref.GetNode()
		if g.outer.IsGuardedInstanceofReference(n) || g.outer.IsUndefinedTypeofGuardReference(n) {
			return
		} else if g.outer.IsUnguardedInstanceofReference(n) {
			instanceofReference :=
				NewCrossChunkCodeMotion_ImmovableInstanceofReference(chunk, ref)

			globalSymbol.instanceofReferencesToGuard.Push(instanceofReference)
			return
		}
	}
	globalSymbol.AddImmovableReference(chunk)
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.processDeclarationStatement(TopLevelStatement)
func (g *CrossChunkCodeMotion_GlobalSymbolCollector) ProcessDeclarationStatement(statement *CrossChunkReferenceCollector_TopLevelStatement) {
	declaredSymbol :=
		g.GetGlobalSymbol(statement.GetDeclaredNameReference().getSymbolV())
	dsg := declaredSymbol.AddDeclarationStatement(statement)
	g.ProcessDeclarationStatementContainedReferences(statement, declaredSymbol, dsg)
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.processDeclarationStatementContainedReferences(TopLevelStatement, GlobalSymbol, DeclarationStatementGroup)
func (g *CrossChunkCodeMotion_GlobalSymbolCollector) ProcessDeclarationStatementContainedReferences(
	statement *CrossChunkReferenceCollector_TopLevelStatement, declaredSymbol *CrossChunkCodeMotion_GlobalSymbol, dsg *CrossChunkCodeMotion_DeclarationStatementGroup) {
	for _, ref := range statement.GetNonDeclarationReferences() {
		refSymbol := g.GetGlobalSymbol(ref.getSymbolV())
		if refSymbol == declaredSymbol {
			continue // ignore circular reference
		}
		if g.outer.parentChunkCanSeeSymbolsDeclaredInChildren {
			// It is possible to move the declaration of `Foo` after
			// `'undefined' != typeof Foo && x instanceof Foo`.
			// We'll add the undefined check, if necessary.
			n := ref.GetNode()
			if g.outer.IsGuardedInstanceofReference(n) || g.outer.IsUndefinedTypeofGuardReference(n) {
				continue
			} else if g.outer.IsUnguardedInstanceofReference(n) {
				instanceofReference :=
					NewCrossChunkCodeMotion_MovableInstanceofReference(dsg, ref)

				refSymbol.instanceofReferencesToGuard.Push(instanceofReference)
				continue
			}
		}
		dsg.AddReferenceToGlobalSymbol(refSymbol)
		refSymbol.AddReferringGlobalSymbol(declaredSymbol)
	}
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.getGlobalSymbol(Var)
func (g *CrossChunkCodeMotion_GlobalSymbolCollector) GetGlobalSymbol(var_ *Var) *CrossChunkCodeMotion_GlobalSymbol {
	globalSymbol := g.globalSymbolforVar.Get(var_)
	if globalSymbol == nil {
		globalSymbol = NewCrossChunkCodeMotion_GlobalSymbol(g.outer, var_)
		g.globalSymbolforVar.Put(var_, globalSymbol)
		g.symbolStack.Push(globalSymbol)
	}
	return globalSymbol
}

// MoveGlobalSymbols moves all of the declaration statements that can move to
// their best possible chunk location.
//
// Java: CrossChunkCodeMotion.moveGlobalSymbols(Collection<GlobalSymbol>)
func (c *CrossChunkCodeMotion) MoveGlobalSymbols(globalSymbols *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol]) {
	for globalSymbolCycle := range NewCrossChunkCodeMotion_OrderAndCombineGlobalSymbols(c, globalSymbols).OrderAndCombine().All() {
		// Symbols whose declarations refer to each other must be grouped together and their
		// declaration statements moved together.
		globalSymbolCycle.MoveDeclarationStatements()
	}
}

// CrossChunkCodeMotion_GlobalSymbol represents a global symbol whose
// declaration statements may be moved.
type CrossChunkCodeMotion_GlobalSymbol struct {
	outer *CrossChunkCodeMotion

	var_ *Var
	// As we traverse the statements in execution order the top of the stack
	// represents the most recently seen DSG for the variable.
	dsgStack *jutil.ArrayDeque[*CrossChunkCodeMotion_DeclarationStatementGroup]

	chunksWithImmovableReferences *jutil.BitSet

	// Symbols whose declaration statements refer to this symbol.
	//
	// This is a LinkedHashSet in order to enforce a consistent ordering when
	// we iterate over these to identify cycles. This guarantees that the order
	// in which statements are moved won't depend on the arbitrary ordering of
	// LinkedHashSet.
	referencingGlobalSymbols *jutil.LinkedHashSet[*CrossChunkCodeMotion_GlobalSymbol]

	// Instanceof references we may need to update with a guard after moving
	// declarations.
	instanceofReferencesToGuard *jutil.ArrayDeque[CrossChunkCodeMotion_InstanceofReference]
	// Used by OrderAndCombineGlobalSymbols to find reference cycles.
	preorderNumber int
	// Used by OrderAndCombineGlobalSymbols to find reference cycles.
	hasBeenAssignedToAStronglyConnectedComponent bool
	// Used to confirm all symbols get moved in the correct order.
	isMoveDeclarationStatementsDone bool
}

// Java: CrossChunkCodeMotion.GlobalSymbol(Var)
func NewCrossChunkCodeMotion_GlobalSymbol(outer *CrossChunkCodeMotion, var_ *Var) *CrossChunkCodeMotion_GlobalSymbol {
	return &CrossChunkCodeMotion_GlobalSymbol{
		outer:                         outer,
		dsgStack:                      jutil.NewArrayDeque[*CrossChunkCodeMotion_DeclarationStatementGroup](),
		chunksWithImmovableReferences: jutil.NewBitSetWithSize(outer.graph.GetChunkCount()),
		referencingGlobalSymbols:      jutil.NewLinkedHashSet[*CrossChunkCodeMotion_GlobalSymbol](),
		instanceofReferencesToGuard:   jutil.NewArrayDeque[CrossChunkCodeMotion_InstanceofReference](),
		preorderNumber:                -1,
		hasBeenAssignedToAStronglyConnectedComponent: false,
		isMoveDeclarationStatementsDone:              false,
		var_:                                         var_,
	}
}

// Java: CrossChunkCodeMotion.GlobalSymbol.toString()
func (s *CrossChunkCodeMotion_GlobalSymbol) String() string {
	return s.var_.GetName()
}

// Java: CrossChunkCodeMotion.GlobalSymbol.addImmovableReference(JSChunk)
func (s *CrossChunkCodeMotion_GlobalSymbol) AddImmovableReference(chunk *JSChunk) {
	s.chunksWithImmovableReferences.Set(chunk.GetIndex())
}

// AddDeclarationStatement adds the statement to the appropriate
// DeclarationStatementGroup and returns it.
//
// Java: CrossChunkCodeMotion.GlobalSymbol.addDeclarationStatement(TopLevelStatement)
func (s *CrossChunkCodeMotion_GlobalSymbol) AddDeclarationStatement(statement *CrossChunkReferenceCollector_TopLevelStatement) *CrossChunkCodeMotion_DeclarationStatementGroup {
	chunk := statement.GetChunk()
	lastDsg := s.dsgStack.Peek()
	if lastDsg == nil {
		lastDsg = NewCrossChunkCodeMotion_DeclarationStatementGroup(s, chunk)
		s.dsgStack.Push(lastDsg)
	}
	var statementDsg *CrossChunkCodeMotion_DeclarationStatementGroup
	if chunk == lastDsg.currentChunk {
		statementDsg = lastDsg
	} else {
		// new chunk requires a new DSG
		statementDsg = NewCrossChunkCodeMotion_DeclarationStatementGroup(s, chunk)
		s.dsgStack.Push(statementDsg)
	}
	statementDsg.statementStack.Push(statement)
	return statementDsg
}

// Java: CrossChunkCodeMotion.GlobalSymbol.addReferringGlobalSymbol(GlobalSymbol)
func (s *CrossChunkCodeMotion_GlobalSymbol) AddReferringGlobalSymbol(declaredSymbol *CrossChunkCodeMotion_GlobalSymbol) {
	s.referencingGlobalSymbols.Add(declaredSymbol)
}

// DeclarationsCoverChunk: does the chunk depend on at least one of the
// chunks containing declaration statements for this symbol?
//
// Java: CrossChunkCodeMotion.GlobalSymbol.declarationsCoverChunk(JSChunk)
func (s *CrossChunkCodeMotion_GlobalSymbol) DeclarationsCoverChunk(chunk *JSChunk) bool {
	for dsg := range s.dsgStack.All() {
		if chunk == dsg.currentChunk || s.outer.graph.DependsOn(chunk, dsg.currentChunk) {
			return true
		}
	}
	return false
}

// CrossChunkCodeMotion_GlobalSymbolCycle represents a set of global symbols
// that all refer to each other, and so must be considered for movement as a
// group.
type CrossChunkCodeMotion_GlobalSymbolCycle struct {
	outer   *CrossChunkCodeMotion
	symbols *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol]
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle()
func NewCrossChunkCodeMotion_GlobalSymbolCycle(outer *CrossChunkCodeMotion) *CrossChunkCodeMotion_GlobalSymbolCycle {
	return &CrossChunkCodeMotion_GlobalSymbolCycle{
		outer:   outer,
		symbols: jutil.NewArrayDeque[*CrossChunkCodeMotion_GlobalSymbol](),
	}
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle.addSymbol(GlobalSymbol)
func (y *CrossChunkCodeMotion_GlobalSymbolCycle) AddSymbol(symbol *CrossChunkCodeMotion_GlobalSymbol) {
	y.symbols.Add(symbol)
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle.moveDeclarationStatements()
func (y *CrossChunkCodeMotion_GlobalSymbolCycle) MoveDeclarationStatements() {
	for symbol := range y.symbols.All() {
		jutil.CheckStateMsg1(!symbol.isMoveDeclarationStatementsDone, "duplicate attempt to move %s", symbol)
	}
	chunksWithImmovableReferences := jutil.NewBitSetWithSize(y.outer.graph.GetChunkCount())
	cyclesLatestFirst := y.GetDsgCyclesLatestFirst()
	for _, dsgCycle := range cyclesLatestFirst {
		// Each move may change chunksWithImmovableReferences
		for symbol := range y.symbols.All() {
			chunksWithImmovableReferences.Or(symbol.chunksWithImmovableReferences)
		}
		dsgCycle.MoveToPreferredChunk(chunksWithImmovableReferences)
	}
	for symbol := range y.symbols.All() {
		symbol.isMoveDeclarationStatementsDone = true
	}
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle.getDsgCyclesLatestFirst()
func (y *CrossChunkCodeMotion_GlobalSymbolCycle) GetDsgCyclesLatestFirst() []*CrossChunkCodeMotion_DeclarationStatementGroupCycle {
	var cyclesLatestFirst []*CrossChunkCodeMotion_DeclarationStatementGroupCycle
	dsgsLatestFirst := y.GetDsgsLatestFirst()
	var cycle *CrossChunkCodeMotion_DeclarationStatementGroupCycle
	for dsg := range dsgsLatestFirst.All() {
		if cycle == nil || cycle.currentChunk != dsg.currentChunk {
			cycle = NewCrossChunkCodeMotion_DeclarationStatementGroupCycle(y.outer, dsg.currentChunk)
			cyclesLatestFirst = append(cyclesLatestFirst, cycle)
		}
		cycle.dsgs.Add(dsg)
	}
	return cyclesLatestFirst
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle.getDsgsLatestFirst()
func (y *CrossChunkCodeMotion_GlobalSymbolCycle) GetDsgsLatestFirst() *jutil.ArrayDeque[*CrossChunkCodeMotion_DeclarationStatementGroup] {
	resultStack := jutil.NewArrayDeque[*CrossChunkCodeMotion_DeclarationStatementGroup]()
	for symbol := range y.symbols.All() {
		stack1 := resultStack
		stack2 := jutil.NewArrayDequeFrom(symbol.dsgStack.All())
		resultStack = jutil.NewArrayDequeWithCapacity[*CrossChunkCodeMotion_DeclarationStatementGroup](stack1.Size() + stack2.Size())
		for {
			if stack1.IsEmpty() {
				resultStack.AddAll(stack2.All())
				break
			} else if stack2.IsEmpty() {
				resultStack.AddAll(stack1.All())
				break
			} else {
				dsg1 := stack1.Peek()
				dsg2 := stack2.Peek()
				if dsg1.currentChunk.GetIndex() > dsg2.currentChunk.GetIndex() {
					resultStack.Add(stack1.Pop())
					jutil.CheckStateMsg(
						stack1.IsEmpty() ||
							stack1.Peek().currentChunk.GetIndex() <= dsg1.currentChunk.GetIndex(),
						"DSG stacks are out of order.")
				} else {
					resultStack.Add(stack2.Pop())
					jutil.CheckStateMsg(
						stack2.IsEmpty() ||
							stack2.Peek().currentChunk.GetIndex() <= dsg2.currentChunk.GetIndex(),
						"DSG stacks are out of order.")
				}
			}
		}
	}
	return resultStack
}

// Java: CrossChunkCodeMotion.addGuardToInstanceofReference(Node)
func (c *CrossChunkCodeMotion) AddGuardToInstanceofReference(referenceNode *rhino.Node) {
	jutil.CheckStateMsg1(
		c.IsUnguardedInstanceofReference(referenceNode),
		"instanceof Reference is already guarded: %s",
		referenceNode)
	instanceofNode := jutil.CheckNotNull(referenceNode.GetParent())
	referenceForTypeOf := referenceNode.CloneNode()
	tmp := rhino.IR_Block()
	// Wrap "foo instanceof Bar" in
	// "('function' == typeof Bar && foo instanceof Bar)"
	instanceofNode.ReplaceWith(tmp)
	and :=
		rhino.IR_And(
			rhino.NewNodeWith2Children(rhino.Token_EQ, rhino.IR_String("function"), rhino.NewNodeWithChild(rhino.Token_TYPEOF, referenceForTypeOf)),
			instanceofNode)
	and.SrcrefTreeIfMissing(instanceofNode)
	tmp.ReplaceWith(and)
	c.compiler.ReportChangeToEnclosingScope(and)
}

// CrossChunkCodeMotion_DeclarationStatementGroup is a group of declaration
// statements that must be moved (or not) as a group.
//
// All of the statements must be in the same chunk initially. If there are
// declarations for the same variable in different chunks, they will be
// grouped separately.
type CrossChunkCodeMotion_DeclarationStatementGroup struct {
	declaredGlobalSymbol    *CrossChunkCodeMotion_GlobalSymbol
	referencedGlobalSymbols *jutil.LinkedHashSet[*CrossChunkCodeMotion_GlobalSymbol]

	// chunk containing the statements
	currentChunk *JSChunk

	// statements in the group, latest first
	statementStack *jutil.ArrayDeque[*CrossChunkReferenceCollector_TopLevelStatement]
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroup(GlobalSymbol, JSChunk)
func NewCrossChunkCodeMotion_DeclarationStatementGroup(declaredGlobalSymbol *CrossChunkCodeMotion_GlobalSymbol, currentChunk *JSChunk) *CrossChunkCodeMotion_DeclarationStatementGroup {
	return &CrossChunkCodeMotion_DeclarationStatementGroup{
		referencedGlobalSymbols: jutil.NewLinkedHashSet[*CrossChunkCodeMotion_GlobalSymbol](),
		statementStack:          jutil.NewArrayDeque[*CrossChunkReferenceCollector_TopLevelStatement](),
		declaredGlobalSymbol:    declaredGlobalSymbol,
		currentChunk:            currentChunk,
	}
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroup.allStatementsCanMove()
func (g *CrossChunkCodeMotion_DeclarationStatementGroup) AllStatementsCanMove() bool {
	for s := range g.statementStack.All() {
		if !s.IsMovableDeclaration() {
			return false
		}
	}
	return true
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroup.addReferenceToGlobalSymbol(GlobalSymbol)
func (g *CrossChunkCodeMotion_DeclarationStatementGroup) AddReferenceToGlobalSymbol(refSymbol *CrossChunkCodeMotion_GlobalSymbol) {
	g.referencedGlobalSymbols.Add(refSymbol)
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroup.makeReferencesImmovable()
func (g *CrossChunkCodeMotion_DeclarationStatementGroup) MakeReferencesImmovable() {
	g.declaredGlobalSymbol.AddImmovableReference(g.currentChunk)
	for symbol := range g.referencedGlobalSymbols.All() {
		jutil.CheckStateMsg2(
			!symbol.isMoveDeclarationStatementsDone,
			"symbol %s moved before referring symbol %s",
			symbol,
			g.declaredGlobalSymbol)
		symbol.AddImmovableReference(g.currentChunk)
	}
}

// CrossChunkCodeMotion_OrderAndCombineGlobalSymbols orders
// DeclarationStatementGroups so that each DSG appears in the list only after
// all of the DSGs that contain references to it.
//
// DSGs that form cycles are combined into a single DSG. This happens when
// declarations within a chunk form cycles by referring to each other.
//
// This is an implementation of the path-based strong component algorithm as
// it is described in the Wikipedia article
// https://en.wikipedia.org/wiki/Path-based_strong_component_algorithm.
type CrossChunkCodeMotion_OrderAndCombineGlobalSymbols struct {
	outer        *CrossChunkCodeMotion
	inputSymbols *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol]
	// Tracks DSGs that may be part of a strongly connected component
	// (reference cycle).
	componentContents *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol]
	// Tracks DSGs that may be the root of a strongly connected component
	// (reference cycle).
	componentRoots *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol]
	// Filled with lists of strongly connected DSGs as they are discovered.
	stronglyConnectedSymbols *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbolCycle]

	preorderCounter int
}

// Java: CrossChunkCodeMotion.OrderAndCombineGlobalSymbols(Collection<GlobalSymbol>)
func NewCrossChunkCodeMotion_OrderAndCombineGlobalSymbols(outer *CrossChunkCodeMotion, dsgs *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol]) *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols {
	return &CrossChunkCodeMotion_OrderAndCombineGlobalSymbols{
		outer:                    outer,
		componentContents:        jutil.NewArrayDeque[*CrossChunkCodeMotion_GlobalSymbol](),
		componentRoots:           jutil.NewArrayDeque[*CrossChunkCodeMotion_GlobalSymbol](),
		preorderCounter:          0,
		inputSymbols:             dsgs,
		stronglyConnectedSymbols: jutil.NewArrayDequeWithCapacity[*CrossChunkCodeMotion_GlobalSymbolCycle](dsgs.Size()),
	}
}

// Java: CrossChunkCodeMotion.OrderAndCombineGlobalSymbols.orderAndCombine()
func (o *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) OrderAndCombine() *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbolCycle] {
	for globalSymbol := range o.inputSymbols.All() {
		if globalSymbol.preorderNumber < 0 {
			o.ProcessGlobalSymbol(globalSymbol)
		} // else already processed
	}
	// At this point stronglyConnectedSymbols has been filled.
	return o.stronglyConnectedSymbols
}

// ProcessGlobalSymbol determines to which strongly connected component this
// GlobalSymbol belongs.
//
// Called exactly once for each GlobalSymbol. When this method returns we
// will have determined which strongly connected component contains
// globalSymbol. Either:
//
//   - globalSymbol was the first member of the strongly connected component
//     we encountered, and the entire component has now been added to
//     stronglyConnectedSymbols.
//   - OR, we encountered the first member of the strongly connected
//     component in an earlier call to this method on the call stack. In
//     this case the top of componentRoots will be the first member, and
//     we'll add the component to stronglyConnectedSymbols once execution has
//     fallen back to the call made with that GlobalSymbol.
//
// Java: CrossChunkCodeMotion.OrderAndCombineGlobalSymbols.processGlobalSymbol(GlobalSymbol)
func (o *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) ProcessGlobalSymbol(symbol *CrossChunkCodeMotion_GlobalSymbol) {
	// preorderNumber is used to track whether we've already processed a DSG and the order in
	// which they have been processed. It is initially -1.
	jutil.CheckStateMsg1(symbol.preorderNumber < 0, "already processed: %s", symbol)
	symbol.preorderNumber = o.preorderCounter
	o.preorderCounter++
	o.componentRoots.Push(symbol)    // could be the start of a new strongly connected component
	o.componentContents.Push(symbol) // could be part of an existing strongly connected component
	for referringSymbol := range symbol.referencingGlobalSymbols.All() {
		if referringSymbol.preorderNumber < 0 {
			o.ProcessGlobalSymbol(referringSymbol)
		} else {
			if !referringSymbol.hasBeenAssignedToAStronglyConnectedComponent {
				// This GlobalSymbol is part of a not-yet-completed strongly connected component.
				// Back off the potential roots stack to the earliest symbol that is part of the
				// component.
				for o.componentRoots.Peek().preorderNumber > referringSymbol.preorderNumber {
					o.componentRoots.Pop()
				}
			}
		}
	}
	if o.componentRoots.Peek() == symbol {
		// After exploring all paths from here, this symbol is still at the top of the potential
		// component roots stack, so this symbol is the root of a strongly connected component.
		o.componentRoots.Pop()
		cycle := NewCrossChunkCodeMotion_GlobalSymbolCycle(o.outer)
		// this symbol and all those after it on the componentContents stack are part of a single
		// cycle.
		var connectedSymbol *CrossChunkCodeMotion_GlobalSymbol
		for {
			connectedSymbol = o.componentContents.Pop()
			cycle.AddSymbol(connectedSymbol)
			connectedSymbol.hasBeenAssignedToAStronglyConnectedComponent = true
			if connectedSymbol == symbol {
				break
			}
		}
		o.stronglyConnectedSymbols.Add(cycle)
	}
}

// CrossChunkCodeMotion_DeclarationStatementGroupCycle is one or more
// DeclarationStatementGroups that that share the same chunk and whose
// declared global symbols refer to each other.
type CrossChunkCodeMotion_DeclarationStatementGroupCycle struct {
	outer        *CrossChunkCodeMotion
	currentChunk *JSChunk
	dsgs         *jutil.ArrayDeque[*CrossChunkCodeMotion_DeclarationStatementGroup]
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle(JSChunk)
func NewCrossChunkCodeMotion_DeclarationStatementGroupCycle(outer *CrossChunkCodeMotion, currentChunk *JSChunk) *CrossChunkCodeMotion_DeclarationStatementGroupCycle {
	return &CrossChunkCodeMotion_DeclarationStatementGroupCycle{
		outer:        outer,
		currentChunk: currentChunk,
		dsgs:         jutil.NewArrayDeque[*CrossChunkCodeMotion_DeclarationStatementGroup](),
	}
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.moveToPreferredChunk(BitSet)
func (y *CrossChunkCodeMotion_DeclarationStatementGroupCycle) MoveToPreferredChunk(chunksWithImmovableReferences *jutil.BitSet) {
	preferredChunk := y.GetPreferredChunk(chunksWithImmovableReferences)
	if preferredChunk != y.currentChunk {
		if y.outer.cccmLog.IsLogging() {
			// Write a separate log for each global symbol, because this is easier to work
			// with when analyzing the log.
			// Include a header, so we can tell which ones were treated as a cycle
			y.outer.cccmLog.LogString("start DSG Cycle move")
			for _, globalSymbolName := range y.GetGlobalSymbolNames() {
				y.outer.cccmLog.LogSupplier(
					func() string {
						return jutil.StringFormat(
							"Moving DSG for %s from chunk %s to chunk %s",
							globalSymbolName, y.currentChunk, preferredChunk)
					})
			}
		}
		y.MoveStatementsToChunk(preferredChunk)
	}
	// Now that all the statements have been moved, update the current chunk for all the DSGs
	// and treat all the references they contain as now immovable.
	for dsg := range y.dsgs.All() {
		dsg.currentChunk = preferredChunk
		dsg.MakeReferencesImmovable()
	}
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.getGlobalSymbolNames()
func (y *CrossChunkCodeMotion_DeclarationStatementGroupCycle) GetGlobalSymbolNames() []string {
	var builder []string
	for dsg := range y.dsgs.All() {
		builder = append(builder, dsg.declaredGlobalSymbol.var_.GetName())
	}
	return builder
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.moveStatementsToChunk(JSChunk)
func (y *CrossChunkCodeMotion_DeclarationStatementGroupCycle) MoveStatementsToChunk(preferredChunk *JSChunk) {
	destParent :=
		y.outer.chunkInsertionPointMap.ComputeIfAbsent(preferredChunk, y.outer.compiler.GetNodeForCodeInsertion)
	statementsLastFirst := y.GetStatementsLastFirst()
	for statement := range statementsLastFirst.All() {
		statementNode := statement.GetStatementNode()
		// Remove it
		y.outer.compiler.ReportChangeToEnclosingScope(statementNode)
		originalScript := NodeUtil_GetEnclosingScript(statementNode)

		statementNode.Detach()

		// Add it to the new spot
		destParent.AddChildToFront(statementNode)
		NodeUtil_AddFeaturesToScript(
			destParent, NodeUtil_GetFeatureSetOfScript(originalScript), y.outer.compiler)
		y.outer.compiler.ReportChangeToEnclosingScope(statementNode)
	}
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.getStatementsLastFirst()
func (y *CrossChunkCodeMotion_DeclarationStatementGroupCycle) GetStatementsLastFirst() *jutil.ArrayDeque[*CrossChunkReferenceCollector_TopLevelStatement] {
	result := jutil.NewArrayDeque[*CrossChunkReferenceCollector_TopLevelStatement]()

	for dsg := range y.dsgs.All() {
		// combine previous result with statements for current DSG
		stack1 := jutil.NewArrayDequeFrom(dsg.statementStack.All())
		stack2 := result
		result = jutil.NewArrayDequeWithCapacity[*CrossChunkReferenceCollector_TopLevelStatement](stack1.Size() + stack2.Size())
		for {
			if stack1.IsEmpty() {
				result.AddAll(stack2.All())
				break
			} else if stack2.IsEmpty() {
				result.AddAll(stack1.All())
				break
			} else {
				s1 := stack1.Peek()
				s2 := stack2.Peek()
				if s1.GetOriginalOrder() > s2.GetOriginalOrder() {
					result.Add(stack1.Pop())
					jutil.CheckStateMsg(
						stack1.IsEmpty() || stack1.Peek().GetOriginalOrder() < s1.GetOriginalOrder(),
						"Statements are recorded in the wrong order.")
				} else {
					result.Add(stack2.Pop())
					jutil.CheckStateMsg(
						stack2.IsEmpty() || stack2.Peek().GetOriginalOrder() < s2.GetOriginalOrder(),
						"Statements are recorded in the wrong order.")
				}
			}
		}
	}
	return result
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.getPreferredChunk(BitSet)
func (y *CrossChunkCodeMotion_DeclarationStatementGroupCycle) GetPreferredChunk(chunksWithImmovableReferences *jutil.BitSet) *JSChunk {
	if chunksWithImmovableReferences.IsEmpty() {
		return y.currentChunk
	} else if !y.AllStatementsCanMove() {
		return y.currentChunk
	} else {
		return y.outer.graph.GetSmallestCoveringSubtree(y.currentChunk, chunksWithImmovableReferences)
	}
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.allStatementsCanMove()
func (y *CrossChunkCodeMotion_DeclarationStatementGroupCycle) AllStatementsCanMove() bool {
	for dsg := range y.dsgs.All() {
		if !dsg.AllStatementsCanMove() {
			return false
		}
	}
	return true
}

// CrossChunkCodeMotion_InstanceofReference is the nested interface
// InstanceofReference.
type CrossChunkCodeMotion_InstanceofReference interface {
	// Java: CrossChunkCodeMotion.InstanceofReference.getChunk()
	GetChunk() *JSChunk

	// Java: CrossChunkCodeMotion.InstanceofReference.getReference()
	GetReference() *Reference
}

// CrossChunkCodeMotion_ImmovableInstanceofReference is the nested class
// ImmovableInstanceofReference.
type CrossChunkCodeMotion_ImmovableInstanceofReference struct {
	chunk     *JSChunk
	reference *Reference
}

// Java: CrossChunkCodeMotion.ImmovableInstanceofReference(JSChunk, Reference)
func NewCrossChunkCodeMotion_ImmovableInstanceofReference(chunk *JSChunk, reference *Reference) *CrossChunkCodeMotion_ImmovableInstanceofReference {
	return &CrossChunkCodeMotion_ImmovableInstanceofReference{chunk: chunk, reference: reference}
}

// Java: CrossChunkCodeMotion.ImmovableInstanceofReference.getChunk()
func (r *CrossChunkCodeMotion_ImmovableInstanceofReference) GetChunk() *JSChunk {
	return r.chunk
}

// Java: CrossChunkCodeMotion.ImmovableInstanceofReference.getReference()
func (r *CrossChunkCodeMotion_ImmovableInstanceofReference) GetReference() *Reference {
	return r.reference
}

// CrossChunkCodeMotion_MovableInstanceofReference is the nested class
// MovableInstanceofReference.
type CrossChunkCodeMotion_MovableInstanceofReference struct {
	containingDsg *CrossChunkCodeMotion_DeclarationStatementGroup
	reference     *Reference
}

// Java: CrossChunkCodeMotion.MovableInstanceofReference(DeclarationStatementGroup, Reference)
func NewCrossChunkCodeMotion_MovableInstanceofReference(containingDsg *CrossChunkCodeMotion_DeclarationStatementGroup, reference *Reference) *CrossChunkCodeMotion_MovableInstanceofReference {
	return &CrossChunkCodeMotion_MovableInstanceofReference{containingDsg: containingDsg, reference: reference}
}

// Java: CrossChunkCodeMotion.MovableInstanceofReference.getChunk()
func (r *CrossChunkCodeMotion_MovableInstanceofReference) GetChunk() *JSChunk {
	return r.containingDsg.currentChunk
}

// Java: CrossChunkCodeMotion.MovableInstanceofReference.getReference()
func (r *CrossChunkCodeMotion_MovableInstanceofReference) GetReference() *Reference {
	return r.reference
}

// IsUndefinedTypeofGuardReference: is the reference node the first Ref in
// an expression like 'undefined' != typeof Ref && x instanceof Ref?
//
// It's safe to ignore this kind of reference when moving the definition of
// Ref.
//
// Java: CrossChunkCodeMotion.isUndefinedTypeofGuardReference(Node)
func (c *CrossChunkCodeMotion) IsUndefinedTypeofGuardReference(reference *rhino.Node) bool {
	// reference => typeof => `!=`
	maybeTypeofGuard := reference.GetGrandparent()
	if maybeTypeofGuard != nil && c.IsExistenceTypeofGuardFor(maybeTypeofGuard, reference) {
		andNode := maybeTypeofGuard.GetParent()
		return andNode != nil &&
			andNode.IsAnd() &&
			c.IsInstanceofFor(andNode.GetLastChild(), reference)
	} else {
		return false
	}
}

// IsExistenceTypeofGuardFor: is the expression of the form 'undefined' !=
// typeof Ref or 'function' == typeof Ref?
//
// expression is the expression being checked; reference the Ref node must
// be equivalent to this node.
//
// Java: CrossChunkCodeMotion.isExistenceTypeofGuardFor(Node, Node)
func (c *CrossChunkCodeMotion) IsExistenceTypeofGuardFor(expression *rhino.Node, reference *rhino.Node) bool {
	if expression.IsNE() || expression.IsSHNE() {
		undefinedString := expression.GetFirstChild()
		typeofNode := expression.GetLastChild()
		return undefinedString.IsStringLit() &&
			undefinedString.GetString() == "undefined" &&
			typeofNode.IsTypeOf() &&
			typeofNode.GetFirstChild().IsEquivalentTo(reference)
	} else if expression.IsEQ() || expression.IsSHEQ() {
		functionString := expression.GetFirstChild()
		typeofNode := expression.GetLastChild()
		return functionString.IsStringLit() &&
			functionString.GetString() == "function" &&
			typeofNode.IsTypeOf() &&
			typeofNode.GetFirstChild().IsEquivalentTo(reference)
	} else {
		return false
	}
}

// IsGuardedInstanceofReference: is the reference node the second Ref in an
// expression like 'undefined' != typeof Ref && x instanceof Ref?
//
// It's safe to ignore this kind of reference when moving the definition of
// Ref.
//
// Java: CrossChunkCodeMotion.isGuardedInstanceofReference(Node)
func (c *CrossChunkCodeMotion) IsGuardedInstanceofReference(reference *rhino.Node) bool {
	instanceofNode := reference.GetParent()
	if c.IsInstanceofFor(instanceofNode, reference) {
		andNode := instanceofNode.GetParent()
		return andNode != nil &&
			andNode.IsAnd() &&
			c.IsExistenceTypeofGuardFor(andNode.GetFirstChild(), reference)
	} else {
		return false
	}
}

// IsUnguardedInstanceofReference: is the reference the right hand side of
// an instanceof and not guarded?
//
// Java: CrossChunkCodeMotion.isUnguardedInstanceofReference(Node)
func (c *CrossChunkCodeMotion) IsUnguardedInstanceofReference(reference *rhino.Node) bool {
	instanceofNode := reference.GetParent()
	if c.IsInstanceofFor(instanceofNode, reference) {
		if c.compiler.GetAstAnalyzer().MayHaveSideEffects(instanceofNode.GetFirstChild()) {
			return false
		}
		return !c.IsGuardedInstanceofReference(reference)
	} else {
		return false
	}
}

// IsInstanceofFor: is the expression of the form x instanceof Ref?
//
// reference is the Ref node; it must be equivalent to this node.
//
// Java: CrossChunkCodeMotion.isInstanceofFor(Node, Node)
func (c *CrossChunkCodeMotion) IsInstanceofFor(expression *rhino.Node, reference *rhino.Node) bool {
	return expression.IsInstanceOf() && expression.GetLastChild().IsEquivalentTo(reference)
}
