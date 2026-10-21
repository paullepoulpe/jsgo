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
	"iter"
	"reflect"
	"slices"

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
	ccArena *crossChunkArena

	compiler AbstractCompiler
	graph    *JSChunkGraph

	// Map from chunk to the node in that chunk that should parent variable
	// declarations that have to be moved into that chunk
	chunkInsertionPointMap *jutil.LinkedHashMap[*JSChunk, *rhino.Node]

	scratchImmovableChunks jutil.BitSet

	parentChunkCanSeeSymbolsDeclaredInChildren bool
}

func (c *CrossChunkCodeMotion) allocGlobalSymbol(var_ *Var) *CrossChunkCodeMotion_GlobalSymbol {
	if ca := c.ccArena; ca != nil {
		if len(ca.symSlab) == cap(ca.symSlab) {
			if ca.symChunkIdx > 0 {
				ca.symChunks[ca.symChunkIdx-1] = ca.symSlab
			}
			if ca.symChunkIdx < len(ca.symChunks) {
				ca.symSlab = ca.symChunks[ca.symChunkIdx][:0]
				ca.symChunkIdx++
			} else {
				newCap := 64
				if cap(ca.symSlab) >= 64 {
					newCap = min(cap(ca.symSlab)*2, 1024)
				}
				ca.symSlab = make([]CrossChunkCodeMotion_GlobalSymbol, 0, newCap)
				ca.symChunks = append(ca.symChunks, ca.symSlab)
				ca.symChunkIdx++
			}
		}
		idx := len(ca.symSlab)
		ca.symSlab = ca.symSlab[:idx+1]
		s := &ca.symSlab[idx]
		s.outer = c
		s.dsgs = s.dsgBuf[:0]
		s.referencingSlice = s.referencingBuf[:0]
		s.preorderNumber = -1
		s.singleImmovableChunk = -1
		s.var_ = var_
		return s
	}
	return NewCrossChunkCodeMotion_GlobalSymbol(c, var_)
}

func (c *CrossChunkCodeMotion) allocDsg(declaredGlobalSymbol *CrossChunkCodeMotion_GlobalSymbol, currentChunk *JSChunk) *CrossChunkCodeMotion_DeclarationStatementGroup {
	if ca := c.ccArena; ca != nil {
		if len(ca.dsgSlab) == cap(ca.dsgSlab) {
			if ca.dsgChunkIdx > 0 {
				ca.dsgChunks[ca.dsgChunkIdx-1] = ca.dsgSlab
			}
			if ca.dsgChunkIdx < len(ca.dsgChunks) {
				ca.dsgSlab = ca.dsgChunks[ca.dsgChunkIdx][:0]
				ca.dsgChunkIdx++
			} else {
				newCap := 64
				if cap(ca.dsgSlab) >= 64 {
					newCap = min(cap(ca.dsgSlab)*2, 1024)
				}
				ca.dsgSlab = make([]CrossChunkCodeMotion_DeclarationStatementGroup, 0, newCap)
				ca.dsgChunks = append(ca.dsgChunks, ca.dsgSlab)
				ca.dsgChunkIdx++
			}
		}
		idx := len(ca.dsgSlab)
		ca.dsgSlab = ca.dsgSlab[:idx+1]
		g := &ca.dsgSlab[idx]
		g.declaredGlobalSymbol = declaredGlobalSymbol
		g.currentChunk = currentChunk
		g.statements = g.stmtBuf[:0]
		g.referencedSlice = g.referencedBuf[:0]
		return g
	}
	return NewCrossChunkCodeMotion_DeclarationStatementGroup(declaredGlobalSymbol, currentChunk)
}

func (c *CrossChunkCodeMotion) allocCycle() *CrossChunkCodeMotion_GlobalSymbolCycle {
	if ca := c.ccArena; ca != nil {
		if len(ca.cycleSlab) == cap(ca.cycleSlab) {
			if ca.cycleChunkIdx > 0 {
				ca.cycleChunks[ca.cycleChunkIdx-1] = ca.cycleSlab
			}
			if ca.cycleChunkIdx < len(ca.cycleChunks) {
				ca.cycleSlab = ca.cycleChunks[ca.cycleChunkIdx][:0]
				ca.cycleChunkIdx++
			} else {
				newCap := 64
				if cap(ca.cycleSlab) >= 64 {
					newCap = min(cap(ca.cycleSlab)*2, 1024)
				}
				ca.cycleSlab = make([]CrossChunkCodeMotion_GlobalSymbolCycle, 0, newCap)
				ca.cycleChunks = append(ca.cycleChunks, ca.cycleSlab)
				ca.cycleChunkIdx++
			}
		}
		idx := len(ca.cycleSlab)
		ca.cycleSlab = ca.cycleSlab[:idx+1]
		y := &ca.cycleSlab[idx]
		y.outer = c
		y.symbolsSlice = y.symbolsBuf[:0]
		return y
	}
	return NewCrossChunkCodeMotion_GlobalSymbolCycle(c)
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
	defer func() {
		c.cccmLog = nil
		c.ccArena = nil
	}()
	c.cccmLog = logFile
	c.cccmLog.LogString("doing cccm")
	// If there are <2 chunks, then we will never move anything, so we're done
	if c.graph.GetChunkCount() > 1 {
		referenceCollector :=
			NewCrossChunkReferenceCollector(c.compiler, NewSyntacticScopeCreator(c.compiler))
		referenceCollector.enableArenaPool()
		defer referenceCollector.release()
		c.ccArena = referenceCollector.ccArena
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
		if globalSymbol.instanceofReferencesToGuard == nil {
			continue
		}
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
	symbolsByVarIndex  []*CrossChunkCodeMotion_GlobalSymbol

	// Returning the symbols in the reverse order in which they are defined
	// helps to minimize unnecessary reordering of declaration statements.
	symbolStack *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol]
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector()
func NewCrossChunkCodeMotion_GlobalSymbolCollector(outer *CrossChunkCodeMotion) *CrossChunkCodeMotion_GlobalSymbolCollector {
	return &CrossChunkCodeMotion_GlobalSymbolCollector{
		outer:       outer,
		symbolStack: jutil.NewArrayDeque[*CrossChunkCodeMotion_GlobalSymbol](),
	}
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.collectGlobalSymbols(CrossChunkReferenceCollector)
func (g *CrossChunkCodeMotion_GlobalSymbolCollector) CollectGlobalSymbols(referenceCollector *CrossChunkReferenceCollector) *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol] {
	if n := len(referenceCollector.vars); n > 0 && g.symbolStack.IsEmpty() {
		g.symbolStack = jutil.NewArrayDequeWithCapacity[*CrossChunkCodeMotion_GlobalSymbol](n)
	}
	if maxIdx := referenceCollector.maxVarIndex; maxIdx >= 0 && maxIdx < 1_000_000 {
		need := maxIdx + 1
		if ca := g.outer.ccArena; ca != nil {
			if cap(ca.symbolsByVarIndex) >= need {
				ca.symbolsByVarIndex = ca.symbolsByVarIndex[:need]
			} else {
				ca.symbolsByVarIndex = make([]*CrossChunkCodeMotion_GlobalSymbol, need)
			}
			g.symbolsByVarIndex = ca.symbolsByVarIndex
		} else {
			g.symbolsByVarIndex = make([]*CrossChunkCodeMotion_GlobalSymbol, need)
		}
	}
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
	refs := statement.GetNonDeclarationReferences()
	vars := statement.nonDeclarationVars
	chunk := statement.GetChunk()
	if len(vars) == len(refs) {
		for i, ref := range refs {
			g.processImmovableReferenceWithVar(ref, vars[i], chunk)
		}
	} else {
		for _, ref := range refs {
			g.ProcessImmovableReference(ref, chunk)
		}
	}
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.processImmovableReference(Reference, JSChunk)
func (g *CrossChunkCodeMotion_GlobalSymbolCollector) ProcessImmovableReference(ref *Reference, chunk *JSChunk) {
	g.processImmovableReferenceWithVar(ref, ref.getSymbolV(), chunk)
}

func (g *CrossChunkCodeMotion_GlobalSymbolCollector) processImmovableReferenceWithVar(ref *Reference, v *Var, chunk *JSChunk) {
	globalSymbol := g.GetGlobalSymbol(v)
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

			if globalSymbol.instanceofReferencesToGuard == nil {
				globalSymbol.instanceofReferencesToGuard = jutil.NewArrayDeque[CrossChunkCodeMotion_InstanceofReference]()
			}
			globalSymbol.instanceofReferencesToGuard.Push(instanceofReference)
			return
		}
	}
	globalSymbol.AddImmovableReference(chunk)
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.processDeclarationStatement(TopLevelStatement)
func (g *CrossChunkCodeMotion_GlobalSymbolCollector) ProcessDeclarationStatement(statement *CrossChunkReferenceCollector_TopLevelStatement) {
	v := statement.declaredVar
	if v == nil {
		v = statement.GetDeclaredNameReference().getSymbolV()
	}
	declaredSymbol := g.GetGlobalSymbol(v)
	dsg := declaredSymbol.AddDeclarationStatement(statement)
	g.ProcessDeclarationStatementContainedReferences(statement, declaredSymbol, dsg)
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.processDeclarationStatementContainedReferences(TopLevelStatement, GlobalSymbol, DeclarationStatementGroup)
func (g *CrossChunkCodeMotion_GlobalSymbolCollector) ProcessDeclarationStatementContainedReferences(
	statement *CrossChunkReferenceCollector_TopLevelStatement, declaredSymbol *CrossChunkCodeMotion_GlobalSymbol, dsg *CrossChunkCodeMotion_DeclarationStatementGroup) {
	refs := statement.GetNonDeclarationReferences()
	vars := statement.nonDeclarationVars
	hasVars := len(vars) == len(refs)
	for i, ref := range refs {
		var v *Var
		if hasVars {
			v = vars[i]
		} else {
			v = ref.getSymbolV()
		}
		refSymbol := g.GetGlobalSymbol(v)
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

				if refSymbol.instanceofReferencesToGuard == nil {
					refSymbol.instanceofReferencesToGuard = jutil.NewArrayDeque[CrossChunkCodeMotion_InstanceofReference]()
				}
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
	idx := var_.GetIndex()
	if idx >= 0 && idx < len(g.symbolsByVarIndex) {
		if sym := g.symbolsByVarIndex[idx]; sym != nil {
			if sym.var_ == var_ {
				return sym
			}
		} else {
			globalSymbol := g.outer.allocGlobalSymbol(var_)
			g.symbolsByVarIndex[idx] = globalSymbol
			g.symbolStack.Push(globalSymbol)
			return globalSymbol
		}
	}
	if g.globalSymbolforVar != nil {
		if sym := g.globalSymbolforVar.Get(var_); sym != nil {
			return sym
		}
	}
	globalSymbol := g.outer.allocGlobalSymbol(var_)
	if g.globalSymbolforVar == nil {
		g.globalSymbolforVar = jutil.NewObjectLinkedHashMap[*Var, *CrossChunkCodeMotion_GlobalSymbol]()
	}
	g.globalSymbolforVar.Put(var_, globalSymbol)
	g.symbolStack.Push(globalSymbol)
	return globalSymbol
}

// MoveGlobalSymbols moves all of the declaration statements that can move to
// their best possible chunk location.
//
// Java: CrossChunkCodeMotion.moveGlobalSymbols(Collection<GlobalSymbol>)
func (c *CrossChunkCodeMotion) MoveGlobalSymbols(globalSymbols *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol]) {
	if globalSymbols.IsEmpty() {
		return
	}
	o := NewCrossChunkCodeMotion_OrderAndCombineGlobalSymbols(c, globalSymbols)
	o.orderAndCombineSlice()
	for _, globalSymbolCycle := range o.cycles {
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
	// As we traverse the statements in execution order the last element of
	// dsgs represents the most recently seen DSG for the variable.
	dsgs   []*CrossChunkCodeMotion_DeclarationStatementGroup
	dsgBuf [1]*CrossChunkCodeMotion_DeclarationStatementGroup

	chunksWithImmovableReferences jutil.BitSet
	singleImmovableChunk          int32

	// Symbols whose declaration statements refer to this symbol.
	//
	// Insertion order is preserved to enforce a consistent ordering when
	// we iterate over these to identify cycles.
	referencingSlice []*CrossChunkCodeMotion_GlobalSymbol
	referencingBuf   [4]*CrossChunkCodeMotion_GlobalSymbol
	referencingMap   map[*CrossChunkCodeMotion_GlobalSymbol]struct{}

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
	s := &CrossChunkCodeMotion_GlobalSymbol{
		outer:                outer,
		preorderNumber:       -1,
		singleImmovableChunk: -1,
		var_:                 var_,
	}
	s.dsgs = s.dsgBuf[:0]
	s.referencingSlice = s.referencingBuf[:0]
	return s
}

func (s *CrossChunkCodeMotion_GlobalSymbol) allDsgsLatestFirst() iter.Seq[*CrossChunkCodeMotion_DeclarationStatementGroup] {
	return func(yield func(*CrossChunkCodeMotion_DeclarationStatementGroup) bool) {
		for i := len(s.dsgs) - 1; i >= 0; i-- {
			if !yield(s.dsgs[i]) {
				return
			}
		}
	}
}

// Java: CrossChunkCodeMotion.GlobalSymbol.toString()
func (s *CrossChunkCodeMotion_GlobalSymbol) String() string {
	return s.var_.GetName()
}

// Java: CrossChunkCodeMotion.GlobalSymbol.addImmovableReference(JSChunk)
func (s *CrossChunkCodeMotion_GlobalSymbol) AddImmovableReference(chunk *JSChunk) {
	idx := int32(chunk.GetIndex())
	if s.singleImmovableChunk == -1 {
		s.singleImmovableChunk = idx
		return
	}
	if s.singleImmovableChunk == idx {
		return
	}
	if s.singleImmovableChunk >= 0 {
		s.chunksWithImmovableReferences.Set(int(s.singleImmovableChunk))
		s.singleImmovableChunk = -2
	}
	s.chunksWithImmovableReferences.Set(int(idx))
}

func (s *CrossChunkCodeMotion_GlobalSymbol) orImmovableChunksInto(dst *jutil.BitSet) {
	if s.singleImmovableChunk >= 0 {
		dst.Set(int(s.singleImmovableChunk))
	} else if s.singleImmovableChunk == -2 {
		dst.Or(&s.chunksWithImmovableReferences)
	}
}

// AddDeclarationStatement adds the statement to the appropriate
// DeclarationStatementGroup and returns it.
//
// Java: CrossChunkCodeMotion.GlobalSymbol.addDeclarationStatement(TopLevelStatement)
func (s *CrossChunkCodeMotion_GlobalSymbol) AddDeclarationStatement(statement *CrossChunkReferenceCollector_TopLevelStatement) *CrossChunkCodeMotion_DeclarationStatementGroup {
	chunk := statement.GetChunk()
	var lastDsg *CrossChunkCodeMotion_DeclarationStatementGroup
	if len(s.dsgs) > 0 {
		lastDsg = s.dsgs[len(s.dsgs)-1]
	} else {
		lastDsg = s.outer.allocDsg(s, chunk)
		s.dsgs = append(s.dsgs, lastDsg)
	}
	var statementDsg *CrossChunkCodeMotion_DeclarationStatementGroup
	if chunk == lastDsg.currentChunk {
		statementDsg = lastDsg
	} else {
		// new chunk requires a new DSG
		statementDsg = s.outer.allocDsg(s, chunk)
		s.dsgs = append(s.dsgs, statementDsg)
	}
	statementDsg.statements = append(statementDsg.statements, statement)
	return statementDsg
}

// Java: CrossChunkCodeMotion.GlobalSymbol.addReferringGlobalSymbol(GlobalSymbol)
func (s *CrossChunkCodeMotion_GlobalSymbol) AddReferringGlobalSymbol(declaredSymbol *CrossChunkCodeMotion_GlobalSymbol) {
	if s.referencingMap != nil {
		if _, exists := s.referencingMap[declaredSymbol]; exists {
			return
		}
		s.referencingMap[declaredSymbol] = struct{}{}
		s.referencingSlice = append(s.referencingSlice, declaredSymbol)
		return
	}
	if slices.Contains(s.referencingSlice, declaredSymbol) {
		return
	}
	if len(s.referencingSlice) == 8 {
		s.referencingMap = make(map[*CrossChunkCodeMotion_GlobalSymbol]struct{}, 16)
		for _, existing := range s.referencingSlice {
			s.referencingMap[existing] = struct{}{}
		}
		s.referencingMap[declaredSymbol] = struct{}{}
	}
	s.referencingSlice = append(s.referencingSlice, declaredSymbol)
}

// DeclarationsCoverChunk: does the chunk depend on at least one of the
// chunks containing declaration statements for this symbol?
//
// Java: CrossChunkCodeMotion.GlobalSymbol.declarationsCoverChunk(JSChunk)
func (s *CrossChunkCodeMotion_GlobalSymbol) DeclarationsCoverChunk(chunk *JSChunk) bool {
	for i := len(s.dsgs) - 1; i >= 0; i-- {
		dsg := s.dsgs[i]
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
	outer        *CrossChunkCodeMotion
	symbolsSlice []*CrossChunkCodeMotion_GlobalSymbol
	symbolsBuf   [1]*CrossChunkCodeMotion_GlobalSymbol
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle()
func NewCrossChunkCodeMotion_GlobalSymbolCycle(outer *CrossChunkCodeMotion) *CrossChunkCodeMotion_GlobalSymbolCycle {
	y := &CrossChunkCodeMotion_GlobalSymbolCycle{
		outer: outer,
	}
	y.symbolsSlice = y.symbolsBuf[:0]
	return y
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle.addSymbol(GlobalSymbol)
func (y *CrossChunkCodeMotion_GlobalSymbolCycle) AddSymbol(symbol *CrossChunkCodeMotion_GlobalSymbol) {
	y.symbolsSlice = append(y.symbolsSlice, symbol)
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle.moveDeclarationStatements()
func (y *CrossChunkCodeMotion_GlobalSymbolCycle) MoveDeclarationStatements() {
	if len(y.symbolsSlice) == 1 {
		symbol := y.symbolsSlice[0]
		jutil.CheckStateMsg1(!symbol.isMoveDeclarationStatementsDone, "duplicate attempt to move %s", symbol)
		if len(symbol.dsgs) == 1 {
			dsg := symbol.dsgs[0]
			var dsgCycle CrossChunkCodeMotion_DeclarationStatementGroupCycle
			dsgCycle.outer = y.outer
			dsgCycle.currentChunk = dsg.currentChunk
			dsgCycle.dsgsSlice = dsgCycle.dsgsBuf[:1]
			dsgCycle.dsgsSlice[0] = dsg
			if symbol.singleImmovableChunk == -1 {
				dsg.MakeReferencesImmovable()
			} else if symbol.singleImmovableChunk >= 0 {
				var preferredChunk *JSChunk
				if !dsg.AllStatementsCanMove() {
					preferredChunk = dsg.currentChunk
				} else {
					targetChunk := y.outer.graph.chunks[symbol.singleImmovableChunk]
					if targetChunk == dsg.currentChunk || y.outer.graph.DependsOn(targetChunk, dsg.currentChunk) {
						preferredChunk = targetChunk
					} else {
						preferredChunk = dsg.currentChunk
					}
				}
				dsgCycle.moveToPreferredChunkDirect(preferredChunk)
			} else {
				dsgCycle.MoveToPreferredChunk(&symbol.chunksWithImmovableReferences)
			}
			symbol.isMoveDeclarationStatementsDone = true
			return
		}
		if len(symbol.dsgs) > 1 {
			chunksWithImmovableReferences := &y.outer.scratchImmovableChunks
			chunksWithImmovableReferences.ClearAll()
			var dsgCycle CrossChunkCodeMotion_DeclarationStatementGroupCycle
			dsgCycle.outer = y.outer
			for i := len(symbol.dsgs) - 1; i >= 0; i-- {
				dsg := symbol.dsgs[i]
				dsgCycle.currentChunk = dsg.currentChunk
				dsgCycle.dsgsSlice = dsgCycle.dsgsBuf[:1]
				dsgCycle.dsgsSlice[0] = dsg
				symbol.orImmovableChunksInto(chunksWithImmovableReferences)
				dsgCycle.MoveToPreferredChunk(chunksWithImmovableReferences)
			}
		}
		symbol.isMoveDeclarationStatementsDone = true
		return
	}
	for _, symbol := range y.symbolsSlice {
		jutil.CheckStateMsg1(!symbol.isMoveDeclarationStatementsDone, "duplicate attempt to move %s", symbol)
	}
	chunksWithImmovableReferences := &y.outer.scratchImmovableChunks
	chunksWithImmovableReferences.ClearAll()
	cyclesLatestFirst := y.GetDsgCyclesLatestFirst()
	for _, dsgCycle := range cyclesLatestFirst {
		// Each move may change chunksWithImmovableReferences
		for _, symbol := range y.symbolsSlice {
			symbol.orImmovableChunksInto(chunksWithImmovableReferences)
		}
		dsgCycle.MoveToPreferredChunk(chunksWithImmovableReferences)
	}
	for _, symbol := range y.symbolsSlice {
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
		cycle.dsgsSlice = append(cycle.dsgsSlice, dsg)
	}
	return cyclesLatestFirst
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle.getDsgsLatestFirst()
func (y *CrossChunkCodeMotion_GlobalSymbolCycle) GetDsgsLatestFirst() *jutil.ArrayDeque[*CrossChunkCodeMotion_DeclarationStatementGroup] {
	resultStack := jutil.NewArrayDeque[*CrossChunkCodeMotion_DeclarationStatementGroup]()
	for _, symbol := range y.symbolsSlice {
		stack1 := resultStack
		stack2 := jutil.NewArrayDequeFrom(symbol.allDsgsLatestFirst())
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
	declaredGlobalSymbol *CrossChunkCodeMotion_GlobalSymbol
	referencedSlice      []*CrossChunkCodeMotion_GlobalSymbol
	referencedBuf        [4]*CrossChunkCodeMotion_GlobalSymbol
	referencedMap        map[*CrossChunkCodeMotion_GlobalSymbol]struct{}

	// chunk containing the statements
	currentChunk *JSChunk

	// statements in the group in chronological push order (last element is latest)
	statements []*CrossChunkReferenceCollector_TopLevelStatement
	stmtBuf    [4]*CrossChunkReferenceCollector_TopLevelStatement
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroup(GlobalSymbol, JSChunk)
func NewCrossChunkCodeMotion_DeclarationStatementGroup(declaredGlobalSymbol *CrossChunkCodeMotion_GlobalSymbol, currentChunk *JSChunk) *CrossChunkCodeMotion_DeclarationStatementGroup {
	g := &CrossChunkCodeMotion_DeclarationStatementGroup{
		declaredGlobalSymbol: declaredGlobalSymbol,
		currentChunk:         currentChunk,
	}
	g.statements = g.stmtBuf[:0]
	g.referencedSlice = g.referencedBuf[:0]
	return g
}

func (g *CrossChunkCodeMotion_DeclarationStatementGroup) allStatementsLatestFirst() iter.Seq[*CrossChunkReferenceCollector_TopLevelStatement] {
	return func(yield func(*CrossChunkReferenceCollector_TopLevelStatement) bool) {
		for i := len(g.statements) - 1; i >= 0; i-- {
			if !yield(g.statements[i]) {
				return
			}
		}
	}
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroup.allStatementsCanMove()
func (g *CrossChunkCodeMotion_DeclarationStatementGroup) AllStatementsCanMove() bool {
	for i := len(g.statements) - 1; i >= 0; i-- {
		if !g.statements[i].IsMovableDeclaration() {
			return false
		}
	}
	return true
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroup.addReferenceToGlobalSymbol(GlobalSymbol)
func (g *CrossChunkCodeMotion_DeclarationStatementGroup) AddReferenceToGlobalSymbol(refSymbol *CrossChunkCodeMotion_GlobalSymbol) {
	if g.referencedMap != nil {
		if _, exists := g.referencedMap[refSymbol]; exists {
			return
		}
		g.referencedMap[refSymbol] = struct{}{}
		g.referencedSlice = append(g.referencedSlice, refSymbol)
		return
	}
	if slices.Contains(g.referencedSlice, refSymbol) {
		return
	}
	if len(g.referencedSlice) == 8 {
		g.referencedMap = make(map[*CrossChunkCodeMotion_GlobalSymbol]struct{}, 16)
		for _, existing := range g.referencedSlice {
			g.referencedMap[existing] = struct{}{}
		}
		g.referencedMap[refSymbol] = struct{}{}
	}
	g.referencedSlice = append(g.referencedSlice, refSymbol)
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroup.makeReferencesImmovable()
func (g *CrossChunkCodeMotion_DeclarationStatementGroup) MakeReferencesImmovable() {
	g.declaredGlobalSymbol.AddImmovableReference(g.currentChunk)
	for _, symbol := range g.referencedSlice {
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
	componentContents []*CrossChunkCodeMotion_GlobalSymbol
	// Tracks DSGs that may be the root of a strongly connected component
	// (reference cycle).
	componentRoots []*CrossChunkCodeMotion_GlobalSymbol
	// Filled with lists of strongly connected DSGs as they are discovered.
	stronglyConnectedSymbols *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbolCycle]
	cycles                   []*CrossChunkCodeMotion_GlobalSymbolCycle

	preorderCounter int
}

// Java: CrossChunkCodeMotion.OrderAndCombineGlobalSymbols(Collection<GlobalSymbol>)
func NewCrossChunkCodeMotion_OrderAndCombineGlobalSymbols(outer *CrossChunkCodeMotion, dsgs *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbol]) *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols {
	o := &CrossChunkCodeMotion_OrderAndCombineGlobalSymbols{
		outer:           outer,
		preorderCounter: 0,
		inputSymbols:    dsgs,
	}
	if ca := outer.ccArena; ca != nil {
		o.componentContents = ca.componentContents[:0]
		o.componentRoots = ca.componentRoots[:0]
		o.cycles = ca.cycles[:0]
	} else if n := dsgs.Size(); n > 0 {
		o.cycles = make([]*CrossChunkCodeMotion_GlobalSymbolCycle, 0, n)
	}
	return o
}

func (o *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) orderAndCombineSlice() {
	for globalSymbol := range o.inputSymbols.All() {
		if globalSymbol.preorderNumber < 0 {
			o.ProcessGlobalSymbol(globalSymbol)
		} // else already processed
	}
	if ca := o.outer.ccArena; ca != nil {
		ca.componentContents = o.componentContents[:0]
		ca.componentRoots = o.componentRoots[:0]
		ca.cycles = o.cycles
	}
}

// Java: CrossChunkCodeMotion.OrderAndCombineGlobalSymbols.orderAndCombine()
func (o *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) OrderAndCombine() *jutil.ArrayDeque[*CrossChunkCodeMotion_GlobalSymbolCycle] {
	o.orderAndCombineSlice()
	o.stronglyConnectedSymbols = jutil.NewArrayDequeWithCapacity[*CrossChunkCodeMotion_GlobalSymbolCycle](len(o.cycles))
	for _, c := range o.cycles {
		o.stronglyConnectedSymbols.Add(c)
	}
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
	o.componentRoots = append(o.componentRoots, symbol)       // could be the start of a new strongly connected component
	o.componentContents = append(o.componentContents, symbol) // could be part of an existing strongly connected component
	for _, referringSymbol := range symbol.referencingSlice {
		if referringSymbol.preorderNumber < 0 {
			o.ProcessGlobalSymbol(referringSymbol)
		} else if !referringSymbol.hasBeenAssignedToAStronglyConnectedComponent {
			// This GlobalSymbol is part of a not-yet-completed strongly connected component.
			// Back off the potential roots stack to the earliest symbol that is part of the
			// component.
			for o.componentRoots[len(o.componentRoots)-1].preorderNumber > referringSymbol.preorderNumber {
				last := len(o.componentRoots) - 1
				o.componentRoots[last] = nil
				o.componentRoots = o.componentRoots[:last]
			}
		}
	}
	if o.componentRoots[len(o.componentRoots)-1] == symbol {
		// After exploring all paths from here, this symbol is still at the top of the potential
		// component roots stack, so this symbol is the root of a strongly connected component.
		lastRoot := len(o.componentRoots) - 1
		o.componentRoots[lastRoot] = nil
		o.componentRoots = o.componentRoots[:lastRoot]
		cycle := o.outer.allocCycle()
		// this symbol and all those after it on the componentContents stack are part of a single
		// cycle.
		for {
			lastContent := len(o.componentContents) - 1
			connectedSymbol := o.componentContents[lastContent]
			o.componentContents[lastContent] = nil
			o.componentContents = o.componentContents[:lastContent]
			cycle.AddSymbol(connectedSymbol)
			connectedSymbol.hasBeenAssignedToAStronglyConnectedComponent = true
			if connectedSymbol == symbol {
				break
			}
		}
		o.cycles = append(o.cycles, cycle)
	}
}

// CrossChunkCodeMotion_DeclarationStatementGroupCycle is one or more
// DeclarationStatementGroups that that share the same chunk and whose
// declared global symbols refer to each other.
type CrossChunkCodeMotion_DeclarationStatementGroupCycle struct {
	outer        *CrossChunkCodeMotion
	currentChunk *JSChunk
	dsgsSlice    []*CrossChunkCodeMotion_DeclarationStatementGroup
	dsgsBuf      [2]*CrossChunkCodeMotion_DeclarationStatementGroup
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle(JSChunk)
func NewCrossChunkCodeMotion_DeclarationStatementGroupCycle(outer *CrossChunkCodeMotion, currentChunk *JSChunk) *CrossChunkCodeMotion_DeclarationStatementGroupCycle {
	y := &CrossChunkCodeMotion_DeclarationStatementGroupCycle{
		outer:        outer,
		currentChunk: currentChunk,
	}
	y.dsgsSlice = y.dsgsBuf[:0]
	return y
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.moveToPreferredChunk(BitSet)
func (y *CrossChunkCodeMotion_DeclarationStatementGroupCycle) MoveToPreferredChunk(chunksWithImmovableReferences *jutil.BitSet) {
	y.moveToPreferredChunkDirect(y.GetPreferredChunk(chunksWithImmovableReferences))
}

func (y *CrossChunkCodeMotion_DeclarationStatementGroupCycle) moveToPreferredChunkDirect(preferredChunk *JSChunk) {
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
	for _, dsg := range y.dsgsSlice {
		dsg.currentChunk = preferredChunk
		dsg.MakeReferencesImmovable()
	}
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.getGlobalSymbolNames()
func (y *CrossChunkCodeMotion_DeclarationStatementGroupCycle) GetGlobalSymbolNames() []string {
	builder := make([]string, 0, len(y.dsgsSlice))
	for _, dsg := range y.dsgsSlice {
		builder = append(builder, dsg.declaredGlobalSymbol.var_.GetName())
	}
	return builder
}

func (y *CrossChunkCodeMotion_DeclarationStatementGroupCycle) moveSingleStatementToChunk(statement *CrossChunkReferenceCollector_TopLevelStatement, destParent *rhino.Node) {
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

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.moveStatementsToChunk(JSChunk)
func (y *CrossChunkCodeMotion_DeclarationStatementGroupCycle) MoveStatementsToChunk(preferredChunk *JSChunk) {
	destParent :=
		y.outer.chunkInsertionPointMap.ComputeIfAbsent(preferredChunk, y.outer.compiler.GetNodeForCodeInsertion)
	if len(y.dsgsSlice) == 1 {
		stmts := y.dsgsSlice[0].statements
		for i := len(stmts) - 1; i >= 0; i-- {
			y.moveSingleStatementToChunk(stmts[i], destParent)
		}
		return
	}
	statementsLastFirst := y.GetStatementsLastFirst()
	for statement := range statementsLastFirst.All() {
		y.moveSingleStatementToChunk(statement, destParent)
	}
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.getStatementsLastFirst()
func (y *CrossChunkCodeMotion_DeclarationStatementGroupCycle) GetStatementsLastFirst() *jutil.ArrayDeque[*CrossChunkReferenceCollector_TopLevelStatement] {
	result := jutil.NewArrayDeque[*CrossChunkReferenceCollector_TopLevelStatement]()

	for _, dsg := range y.dsgsSlice {
		// combine previous result with statements for current DSG
		stack1 := jutil.NewArrayDequeFrom(dsg.allStatementsLatestFirst())
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
	for _, dsg := range y.dsgsSlice {
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
