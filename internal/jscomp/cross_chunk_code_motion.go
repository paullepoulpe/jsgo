// Port of com/google/javascript/jscomp/CrossChunkCodeMotion.java @ closure-compiler v20261028

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/ast"
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
	chunkInsertionPointMap map[*JSChunk]rhino.Node

	scratchImmovableChunks smallBitSet
	scratchDsgCycle        cccmDsgCycle

	symbolSlab jutil.Slab[cccmGlobalSymbol]
	dsgSlab    jutil.Slab[cccmDsg]
	cycleSlab  jutil.Slab[cccmSymbolCycle]

	parentChunkCanSeeSymbolsDeclaredInChildren bool
}

func (c *CrossChunkCodeMotion) allocGlobalSymbol(var_ *Var) *cccmGlobalSymbol {
	s := c.symbolSlab.New()
	s.outer = c
	s.preorderNumber = -1
	s.var_ = var_
	s.dsgs = s.dsgBuf[:0]
	s.referencingSlice = s.referencingBuf[:0]
	return s
}

func (c *CrossChunkCodeMotion) allocDsg(declaredGlobalSymbol *cccmGlobalSymbol, currentChunk *JSChunk) *cccmDsg {
	g := c.dsgSlab.New()
	g.declaredGlobalSymbol = declaredGlobalSymbol
	g.currentChunk = currentChunk
	g.statements = g.stmtBuf[:0]
	g.referenced.clear()
	return g
}

func (c *CrossChunkCodeMotion) allocCycle() *cccmSymbolCycle {
	y := c.cycleSlab.New()
	y.outer = c
	y.symbolsSlice = y.symbolsBuf[:0]
	return y
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
		chunkInsertionPointMap: make(map[*JSChunk]rhino.Node),
		compiler:               compiler,
		graph:                  graph,
		parentChunkCanSeeSymbolsDeclaredInChildren: parentChunkCanSeeSymbolsDeclaredInChildren,
	}
}

// Java: CrossChunkCodeMotion.process(Node, Node)
func (c *CrossChunkCodeMotion) Process(externs rhino.Node, root rhino.Node) {
	logFile := c.compiler.CreateOrReopenIndexedLog("CrossChunkCodeMotion", "cccm.log")
	defer logFile.Close()
	defer func() {
		c.cccmLog = diagnostic.LogFile{}
		c.symbolSlab = jutil.Slab[cccmGlobalSymbol]{}
		c.dsgSlab = jutil.Slab[cccmDsg]{}
		c.cycleSlab = jutil.Slab[cccmSymbolCycle]{}
		c.scratchImmovableChunks = smallBitSet{}
		clear(c.chunkInsertionPointMap)
	}()
	c.cccmLog = logFile
	c.cccmLog.LogString("doing cccm")
	// If there are <2 chunks, then we will never move anything, so we're done
	if c.graph.GetChunkCount() > 1 {
		referenceCollector :=
			NewCrossChunkReferenceCollector(c.compiler, NewSyntacticScopeCreator(c.compiler))
		referenceCollector.ProcessRoot(root)
		collector := cccmSymbolCollector{outer: c}
		globalSymbols := collector.collectGlobalSymbols(referenceCollector)
		c.moveGlobalSymbols(globalSymbols)
		c.addInstanceofGuards(globalSymbols)
	} else {
		c.cccmLog.LogString("only one chunk exists")
	}
}

// Java: CrossChunkCodeMotion.addInstanceofGuards(Collection<GlobalSymbol>)
func (c *CrossChunkCodeMotion) addInstanceofGuards(globalSymbols []*cccmGlobalSymbol) {
	for _, globalSymbol := range globalSymbols {
		refs := globalSymbol.instanceofReferencesToGuard
		for i := len(refs) - 1; i >= 0; i-- {
			ref := refs[i]
			if !globalSymbol.declsCoverChunk(ref.targetChunk()) {
				c.addGuardToInstanceofRef(ref.node)
			}
		}
	}
}

// cccmSymbolCollector collects all global symbols, their declaration statements
// and references.
type cccmSymbolCollector struct {
	outer *CrossChunkCodeMotion

	globalSymbolForVar map[*Var]*cccmGlobalSymbol
	symbolsByVarIndex  []*cccmGlobalSymbol

	// Returning the symbols in the reverse order in which they are defined
	// helps to minimize unnecessary reordering of declaration statements.
	symbolStack []*cccmGlobalSymbol
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.collectGlobalSymbols(CrossChunkReferenceCollector)
func (g *cccmSymbolCollector) collectGlobalSymbols(referenceCollector *CrossChunkReferenceCollector) []*cccmGlobalSymbol {
	if n := len(referenceCollector.vars); n > 0 && len(g.symbolStack) == 0 {
		g.symbolStack = make([]*cccmGlobalSymbol, 0, n)
	}
	if maxIdx := referenceCollector.maxVarIndex; maxIdx >= 0 && maxIdx < 1_000_000 {
		g.symbolsByVarIndex = make([]*cccmGlobalSymbol, maxIdx+1)
	}
	for _, statement := range referenceCollector.GetTopLevelStatements() {
		if statement.IsDeclarationStatement() {
			g.processDeclStmt(statement)
		} else {
			g.processImmovableStmt(statement)
		}
	}
	slices.Reverse(g.symbolStack)
	return g.symbolStack
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.processImmovableStatement(TopLevelStatement)
func (g *cccmSymbolCollector) processImmovableStmt(statement *CrossChunkReferenceCollector_TopLevelStatement) {
	refs := statement.GetNonDeclarationReferences()
	vars := statement.nonDeclarationVars
	chunk := statement.GetChunk()
	if len(vars) == len(refs) {
		for i, ref := range refs {
			g.processImmovableRefWithVar(ref, vars[i], chunk)
		}
	} else {
		for _, ref := range refs {
			g.processImmovableRefWithVar(ref, ref.getSymbolV(), chunk)
		}
	}
}

func (g *cccmSymbolCollector) processImmovableRefWithVar(ref *Reference, v *Var, chunk *JSChunk) {
	globalSymbol := g.globalSymbol(v)
	if g.outer.parentChunkCanSeeSymbolsDeclaredInChildren {
		// It is possible to move the declaration of `Foo` after
		// `'undefined' != typeof Foo && x instanceof Foo`.
		// We'll add the undefined check, if necessary.
		n := ref.GetNode()
		if g.outer.isGuardedInstanceofRef(n) || g.outer.isUndefinedTypeofGuardRef(n) {
			return
		} else if g.outer.isUnguardedInstanceofRef(n) {
			globalSymbol.instanceofReferencesToGuard = append(
				globalSymbol.instanceofReferencesToGuard,
				instanceofRef{node: n, chunk: chunk})
			return
		}
	}
	globalSymbol.addImmovableRef(chunk)
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.processDeclarationStatement(TopLevelStatement)
func (g *cccmSymbolCollector) processDeclStmt(statement *CrossChunkReferenceCollector_TopLevelStatement) {
	v := statement.declaredVar
	if v == nil {
		v = statement.GetDeclaredNameReference().getSymbolV()
	}
	declaredSymbol := g.globalSymbol(v)
	dsg := declaredSymbol.addDeclStmt(statement)
	g.processDeclStmtContainedRefs(statement, declaredSymbol, dsg)
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.processDeclarationStatementContainedReferences(TopLevelStatement, GlobalSymbol, DeclarationStatementGroup)
func (g *cccmSymbolCollector) processDeclStmtContainedRefs(
	statement *CrossChunkReferenceCollector_TopLevelStatement, declaredSymbol *cccmGlobalSymbol, dsg *cccmDsg) {
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
		refSymbol := g.globalSymbol(v)
		if refSymbol == declaredSymbol {
			continue // ignore circular reference
		}
		if g.outer.parentChunkCanSeeSymbolsDeclaredInChildren {
			// It is possible to move the declaration of `Foo` after
			// `'undefined' != typeof Foo && x instanceof Foo`.
			// We'll add the undefined check, if necessary.
			n := ref.GetNode()
			if g.outer.isGuardedInstanceofRef(n) || g.outer.isUndefinedTypeofGuardRef(n) {
				continue
			} else if g.outer.isUnguardedInstanceofRef(n) {
				refSymbol.instanceofReferencesToGuard = append(
					refSymbol.instanceofReferencesToGuard,
					instanceofRef{node: n, dsg: dsg, movable: true})
				continue
			}
		}
		if dsg.addRefToGlobalSymbol(refSymbol) {
			if len(declaredSymbol.dsgs) == 1 {
				refSymbol.referencingSlice = append(refSymbol.referencingSlice, declaredSymbol)
			} else {
				alreadyReferred := false
				for _, prevDsg := range declaredSymbol.dsgs[:len(declaredSymbol.dsgs)-1] {
					if prevDsg.hasRefToGlobalSymbol(refSymbol) {
						alreadyReferred = true
						break
					}
				}
				if !alreadyReferred {
					refSymbol.referencingSlice = append(refSymbol.referencingSlice, declaredSymbol)
				}
			}
		}
	}
}

// Java: CrossChunkCodeMotion.GlobalSymbolCollector.getGlobalSymbol(Var)
func (g *cccmSymbolCollector) globalSymbol(var_ *Var) *cccmGlobalSymbol {
	idx := var_.GetIndex()
	if idx >= 0 && idx < len(g.symbolsByVarIndex) {
		if sym := g.symbolsByVarIndex[idx]; sym != nil {
			if sym.var_ == var_ {
				return sym
			}
		} else {
			globalSymbol := g.outer.allocGlobalSymbol(var_)
			g.symbolsByVarIndex[idx] = globalSymbol
			g.symbolStack = append(g.symbolStack, globalSymbol)
			return globalSymbol
		}
	}
	if sym := g.globalSymbolForVar[var_]; sym != nil {
		return sym
	}
	globalSymbol := g.outer.allocGlobalSymbol(var_)
	if g.globalSymbolForVar == nil {
		g.globalSymbolForVar = make(map[*Var]*cccmGlobalSymbol)
	}
	g.globalSymbolForVar[var_] = globalSymbol
	g.symbolStack = append(g.symbolStack, globalSymbol)
	return globalSymbol
}

// moveGlobalSymbols moves all of the declaration statements that can move to
// their best possible chunk location.
//
// Java: CrossChunkCodeMotion.moveGlobalSymbols(Collection<GlobalSymbol>)
func (c *CrossChunkCodeMotion) moveGlobalSymbols(globalSymbols []*cccmGlobalSymbol) {
	if len(globalSymbols) == 0 {
		return
	}
	o := cccmSymbolOrderer{
		outer:        c,
		inputSymbols: globalSymbols,
		cycles:       make([]*cccmSymbolCycle, 0, len(globalSymbols)),
	}
	o.orderAndCombine()
	for _, globalSymbolCycle := range o.cycles {
		// Symbols whose declarations refer to each other must be grouped together and their
		// declaration statements moved together.
		globalSymbolCycle.moveDeclStmts()
	}
}

// cccmGlobalSymbol represents a global symbol whose declaration statements may
// be moved.
type cccmGlobalSymbol struct {
	outer *CrossChunkCodeMotion

	var_ *Var
	// As we traverse the statements in execution order the last element of
	// dsgs represents the most recently seen DSG for the variable.
	dsgs   []*cccmDsg
	dsgBuf [2]*cccmDsg

	chunksWithImmovableReferences smallBitSet

	// Symbols whose declaration statements refer to this symbol.
	//
	// Insertion order is preserved to enforce a consistent ordering when
	// we iterate over these to identify cycles.
	referencingSlice []*cccmGlobalSymbol
	referencingBuf   [4]*cccmGlobalSymbol

	// Instanceof references we may need to update with a guard after moving
	// declarations (in push order; iterated latest-first).
	instanceofReferencesToGuard []instanceofRef
	// Used by cccmSymbolOrderer to find reference cycles.
	preorderNumber int
	// Used by cccmSymbolOrderer to find reference cycles.
	hasBeenAssignedToAStronglyConnectedComponent bool
	// Used to confirm all symbols get moved in the correct order.
	isMoveDeclarationStatementsDone bool
}

// Java: CrossChunkCodeMotion.GlobalSymbol.toString()
func (s *cccmGlobalSymbol) String() string {
	return s.var_.GetName()
}

// Java: CrossChunkCodeMotion.GlobalSymbol.addImmovableReference(JSChunk)
func (s *cccmGlobalSymbol) addImmovableRef(chunk *JSChunk) {
	s.chunksWithImmovableReferences.Set(chunk.GetIndex())
}

func (s *cccmGlobalSymbol) orImmovableChunksInto(dst *smallBitSet) {
	dst.Or(&s.chunksWithImmovableReferences)
}

// addDeclStmt adds the statement to the appropriate cccmDsg and returns it.
//
// Java: CrossChunkCodeMotion.GlobalSymbol.addDeclarationStatement(TopLevelStatement)
func (s *cccmGlobalSymbol) addDeclStmt(statement *CrossChunkReferenceCollector_TopLevelStatement) *cccmDsg {
	chunk := statement.GetChunk()
	var lastDsg *cccmDsg
	if len(s.dsgs) > 0 {
		lastDsg = s.dsgs[len(s.dsgs)-1]
	} else {
		lastDsg = s.outer.allocDsg(s, chunk)
		s.dsgs = append(s.dsgs, lastDsg)
	}
	var statementDsg *cccmDsg
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

// declsCoverChunk reports whether the chunk depends on at least one of the
// chunks containing declaration statements for this symbol.
//
// Java: CrossChunkCodeMotion.GlobalSymbol.declarationsCoverChunk(JSChunk)
func (s *cccmGlobalSymbol) declsCoverChunk(chunk *JSChunk) bool {
	for i := len(s.dsgs) - 1; i >= 0; i-- {
		dsg := s.dsgs[i]
		if chunk == dsg.currentChunk || s.outer.graph.DependsOn(chunk, dsg.currentChunk) {
			return true
		}
	}
	return false
}

// cccmSymbolCycle represents a set of global symbols that all refer to each
// other, and so must be considered for movement as a group.
type cccmSymbolCycle struct {
	outer        *CrossChunkCodeMotion
	symbolsSlice []*cccmGlobalSymbol
	symbolsBuf   [1]*cccmGlobalSymbol
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle.addSymbol(GlobalSymbol)
func (y *cccmSymbolCycle) addSymbol(symbol *cccmGlobalSymbol) {
	y.symbolsSlice = append(y.symbolsSlice, symbol)
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle.moveDeclarationStatements()
func (y *cccmSymbolCycle) moveDeclStmts() {
	if len(y.symbolsSlice) == 1 {
		symbol := y.symbolsSlice[0]
		jutil.CheckStateMsg1(!symbol.isMoveDeclarationStatementsDone, "duplicate attempt to move %s", symbol)
		if len(symbol.dsgs) == 1 {
			dsg := symbol.dsgs[0]
			dsgCycle := &y.outer.scratchDsgCycle
			dsgCycle.outer = y.outer
			dsgCycle.currentChunk = dsg.currentChunk
			dsgCycle.dsgsSlice = dsgCycle.dsgsBuf[:1]
			dsgCycle.dsgsSlice[0] = dsg
			if symbol.chunksWithImmovableReferences.IsEmpty() {
				dsg.makeRefsImmovable()
			} else if symbol.chunksWithImmovableReferences.Cardinality() == 1 {
				var preferredChunk *JSChunk
				if !dsg.allStmtsCanMove() {
					preferredChunk = dsg.currentChunk
				} else {
					chunkIdx := symbol.chunksWithImmovableReferences.NextSetBit(0)
					targetChunk := y.outer.graph.chunks[chunkIdx]
					if targetChunk == dsg.currentChunk || y.outer.graph.DependsOn(targetChunk, dsg.currentChunk) {
						preferredChunk = targetChunk
					} else {
						preferredChunk = dsg.currentChunk
					}
				}
				dsgCycle.moveToPreferredChunkDirect(preferredChunk)
			} else {
				dsgCycle.moveToPreferredChunk(&symbol.chunksWithImmovableReferences)
			}
			dsgCycle.dsgsSlice[0] = nil
			symbol.isMoveDeclarationStatementsDone = true
			return
		}
		if len(symbol.dsgs) > 1 {
			chunksWithImmovableReferences := &y.outer.scratchImmovableChunks
			chunksWithImmovableReferences.ClearAll()
			dsgCycle := &y.outer.scratchDsgCycle
			dsgCycle.outer = y.outer
			for i := len(symbol.dsgs) - 1; i >= 0; i-- {
				dsg := symbol.dsgs[i]
				dsgCycle.currentChunk = dsg.currentChunk
				dsgCycle.dsgsSlice = dsgCycle.dsgsBuf[:1]
				dsgCycle.dsgsSlice[0] = dsg
				symbol.orImmovableChunksInto(chunksWithImmovableReferences)
				dsgCycle.moveToPreferredChunk(chunksWithImmovableReferences)
			}
			dsgCycle.dsgsSlice[0] = nil
		}
		symbol.isMoveDeclarationStatementsDone = true
		return
	}
	for _, symbol := range y.symbolsSlice {
		jutil.CheckStateMsg1(!symbol.isMoveDeclarationStatementsDone, "duplicate attempt to move %s", symbol)
	}
	chunksWithImmovableReferences := &y.outer.scratchImmovableChunks
	chunksWithImmovableReferences.ClearAll()
	cyclesLatestFirst := y.dsgCyclesLatestFirst()
	for _, dsgCycle := range cyclesLatestFirst {
		// Each move may change chunksWithImmovableReferences
		for _, symbol := range y.symbolsSlice {
			symbol.orImmovableChunksInto(chunksWithImmovableReferences)
		}
		dsgCycle.moveToPreferredChunk(chunksWithImmovableReferences)
	}
	for _, symbol := range y.symbolsSlice {
		symbol.isMoveDeclarationStatementsDone = true
	}
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle.getDsgCyclesLatestFirst()
func (y *cccmSymbolCycle) dsgCyclesLatestFirst() []*cccmDsgCycle {
	var cyclesLatestFirst []*cccmDsgCycle
	dsgsLatestFirst := y.dsgsLatestFirst()
	var cycle *cccmDsgCycle
	for _, dsg := range dsgsLatestFirst {
		if cycle == nil || cycle.currentChunk != dsg.currentChunk {
			cycle = newCccmDsgCycle(y.outer, dsg.currentChunk)
			cyclesLatestFirst = append(cyclesLatestFirst, cycle)
		}
		cycle.dsgsSlice = append(cycle.dsgsSlice, dsg)
	}
	return cyclesLatestFirst
}

// mergeDescendingBy merges desc (already in descending order) with asc (in
// ascending order, iterated in reverse) into a new descending slice, checking
// descending monotonicity on the consumed side at each comparison step.
func mergeDescendingBy[T comparable](desc, asc []T, order func(T) int, strict bool, errMsg string) []T {
	result := make([]T, 0, len(desc)+len(asc))
	i := 0
	j := len(asc) - 1
	for {
		if i == len(desc) {
			for ; j >= 0; j-- {
				result = append(result, asc[j])
			}
			return result
		}
		if j < 0 {
			return append(result, desc[i:]...)
		}
		d := desc[i]
		a := asc[j]
		od := order(d)
		oa := order(a)
		if od > oa {
			result = append(result, d)
			i++
			jutil.CheckStateMsg(
				i == len(desc) || order(desc[i]) < od || (!strict && order(desc[i]) == od),
				errMsg)
		} else {
			result = append(result, a)
			j--
			jutil.CheckStateMsg(
				j < 0 || order(asc[j]) < oa || (!strict && order(asc[j]) == oa),
				errMsg)
		}
	}
}

// Java: CrossChunkCodeMotion.GlobalSymbolCycle.getDsgsLatestFirst()
func (y *cccmSymbolCycle) dsgsLatestFirst() []*cccmDsg {
	var result []*cccmDsg
	for _, symbol := range y.symbolsSlice {
		result = mergeDescendingBy(
			result,
			symbol.dsgs,
			func(dsg *cccmDsg) int { return dsg.currentChunk.GetIndex() },
			false,
			"DSG stacks are out of order.")
	}
	return result
}

// Java: CrossChunkCodeMotion.addGuardToInstanceofReference(Node)
func (c *CrossChunkCodeMotion) addGuardToInstanceofRef(referenceNode rhino.Node) {
	jutil.CheckStateMsg1(
		c.isUnguardedInstanceofRef(referenceNode),
		"instanceof Reference is already guarded: %s",
		referenceNode)
	instanceofBin := ast.MustBinary(rhino.CheckNotNull(referenceNode.GetParent()))
	instanceofNode := instanceofBin.Node()
	referenceForTypeOf := ast.MustExpr(referenceNode.CloneNode())
	tmp := rhino.IR_Block()
	// Wrap "foo instanceof Bar" in
	// "('function' == typeof Bar && foo instanceof Bar)"
	instanceofNode.ReplaceWith(tmp)
	typeofGuard := ast.NewBinary(
		rhino.Token_EQ,
		ast.NewStringLit("function").AsExpr(),
		ast.NewUnary(rhino.Token_TYPEOF, referenceForTypeOf).AsExpr())
	and := ast.NewBinary(rhino.Token_AND, typeofGuard.AsExpr(), instanceofBin.AsExpr()).Node()
	and.SrcrefTreeIfMissing(instanceofNode)
	tmp.ReplaceWith(and)
	c.compiler.ReportChangeToEnclosingScope(and)
}

// cccmDsg is a group of declaration statements that must be moved (or not) as
// a group.
//
// All of the statements must be in the same chunk initially. If there are
// declarations for the same variable in different chunks, they will be
// grouped separately.
type cccmDsg struct {
	declaredGlobalSymbol *cccmGlobalSymbol
	referenced           smallSet[*cccmGlobalSymbol]

	// chunk containing the statements
	currentChunk *JSChunk

	// statements in the group in chronological push order (last element is latest)
	statements []*CrossChunkReferenceCollector_TopLevelStatement
	stmtBuf    [4]*CrossChunkReferenceCollector_TopLevelStatement
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroup.allStatementsCanMove()
func (g *cccmDsg) allStmtsCanMove() bool {
	for i := len(g.statements) - 1; i >= 0; i-- {
		if !g.statements[i].IsMovableDeclaration() {
			return false
		}
	}
	return true
}

func (g *cccmDsg) hasRefToGlobalSymbol(refSymbol *cccmGlobalSymbol) bool {
	return g.referenced.contains(refSymbol)
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroup.addReferenceToGlobalSymbol(GlobalSymbol)
func (g *cccmDsg) addRefToGlobalSymbol(refSymbol *cccmGlobalSymbol) bool {
	return g.referenced.add(refSymbol)
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroup.makeReferencesImmovable()
func (g *cccmDsg) makeRefsImmovable() {
	g.declaredGlobalSymbol.addImmovableRef(g.currentChunk)
	for _, symbol := range g.referenced.items {
		jutil.CheckStateMsg2(
			!symbol.isMoveDeclarationStatementsDone,
			"symbol %s moved before referring symbol %s",
			symbol,
			g.declaredGlobalSymbol)
		symbol.addImmovableRef(g.currentChunk)
	}
}

// cccmSymbolOrderer orders DeclarationStatementGroups so that each DSG appears
// in the list only after all of the DSGs that contain references to it.
//
// DSGs that form cycles are combined into a single DSG. This happens when
// declarations within a chunk form cycles by referring to each other.
//
// This is an implementation of the path-based strong component algorithm as
// it is described in the Wikipedia article
// https://en.wikipedia.org/wiki/Path-based_strong_component_algorithm.
type cccmSymbolOrderer struct {
	outer        *CrossChunkCodeMotion
	inputSymbols []*cccmGlobalSymbol
	// Tracks DSGs that may be part of a strongly connected component
	// (reference cycle).
	componentContents []*cccmGlobalSymbol
	// Tracks DSGs that may be the root of a strongly connected component
	// (reference cycle).
	componentRoots []*cccmGlobalSymbol
	// Filled with lists of strongly connected DSGs as they are discovered.
	cycles []*cccmSymbolCycle

	preorderCounter int
}

func (o *cccmSymbolOrderer) orderAndCombine() {
	for _, globalSymbol := range o.inputSymbols {
		if globalSymbol.preorderNumber < 0 {
			o.processGlobalSymbol(globalSymbol)
		} // else already processed
	}
}

// processGlobalSymbol determines to which strongly connected component this
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
func (o *cccmSymbolOrderer) processGlobalSymbol(symbol *cccmGlobalSymbol) {
	// preorderNumber is used to track whether we've already processed a DSG and the order in
	// which they have been processed. It is initially -1.
	jutil.CheckStateMsg1(symbol.preorderNumber < 0, "already processed: %s", symbol)
	symbol.preorderNumber = o.preorderCounter
	o.preorderCounter++
	o.componentRoots = append(o.componentRoots, symbol)       // could be the start of a new strongly connected component
	o.componentContents = append(o.componentContents, symbol) // could be part of an existing strongly connected component
	for _, referringSymbol := range symbol.referencingSlice {
		if referringSymbol.preorderNumber < 0 {
			o.processGlobalSymbol(referringSymbol)
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
			cycle.addSymbol(connectedSymbol)
			connectedSymbol.hasBeenAssignedToAStronglyConnectedComponent = true
			if connectedSymbol == symbol {
				break
			}
		}
		o.cycles = append(o.cycles, cycle)
	}
}

// cccmDsgCycle is one or more DeclarationStatementGroups that share the same
// chunk and whose declared global symbols refer to each other.
type cccmDsgCycle struct {
	outer        *CrossChunkCodeMotion
	currentChunk *JSChunk
	dsgsSlice    []*cccmDsg
	dsgsBuf      [2]*cccmDsg
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle(JSChunk)
func newCccmDsgCycle(outer *CrossChunkCodeMotion, currentChunk *JSChunk) *cccmDsgCycle {
	y := &cccmDsgCycle{
		outer:        outer,
		currentChunk: currentChunk,
	}
	y.dsgsSlice = y.dsgsBuf[:0]
	return y
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.moveToPreferredChunk(BitSet)
func (y *cccmDsgCycle) moveToPreferredChunk(chunksWithImmovableReferences *smallBitSet) {
	y.moveToPreferredChunkDirect(y.preferredChunk(chunksWithImmovableReferences))
}

func (y *cccmDsgCycle) moveToPreferredChunkDirect(preferredChunk *JSChunk) {
	if preferredChunk != y.currentChunk {
		if y.outer.cccmLog.IsLogging() {
			// Write a separate log for each global symbol, because this is easier to work
			// with when analyzing the log.
			// Include a header, so we can tell which ones were treated as a cycle
			y.outer.cccmLog.LogString("start DSG Cycle move")
			for _, globalSymbolName := range y.globalSymbolNames() {
				y.outer.cccmLog.LogSupplier(
					func() string {
						return "Moving DSG for " + globalSymbolName +
							" from chunk " + cccmChunkName(y.currentChunk) +
							" to chunk " + cccmChunkName(preferredChunk)
					})
			}
		}
		y.moveStmtsToChunk(preferredChunk)
	}
	// Now that all the statements have been moved, update the current chunk for all the DSGs
	// and treat all the references they contain as now immovable.
	for _, dsg := range y.dsgsSlice {
		dsg.currentChunk = preferredChunk
		dsg.makeRefsImmovable()
	}
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.getGlobalSymbolNames()
func (y *cccmDsgCycle) globalSymbolNames() []string {
	builder := make([]string, 0, len(y.dsgsSlice))
	for _, dsg := range y.dsgsSlice {
		builder = append(builder, dsg.declaredGlobalSymbol.var_.GetName())
	}
	return builder
}

func (y *cccmDsgCycle) moveSingleStmtToChunk(statement *CrossChunkReferenceCollector_TopLevelStatement, destParent rhino.Node) {
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
func (y *cccmDsgCycle) moveStmtsToChunk(preferredChunk *JSChunk) {
	destParent, ok := y.outer.chunkInsertionPointMap[preferredChunk]
	if !ok {
		destParent = y.outer.compiler.GetNodeForCodeInsertion(preferredChunk)
		y.outer.chunkInsertionPointMap[preferredChunk] = destParent
	}
	if len(y.dsgsSlice) == 1 {
		stmts := y.dsgsSlice[0].statements
		for i := len(stmts) - 1; i >= 0; i-- {
			y.moveSingleStmtToChunk(stmts[i], destParent)
		}
		return
	}
	for _, statement := range y.stmtsLastFirst() {
		y.moveSingleStmtToChunk(statement, destParent)
	}
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.getStatementsLastFirst()
func (y *cccmDsgCycle) stmtsLastFirst() []*CrossChunkReferenceCollector_TopLevelStatement {
	var result []*CrossChunkReferenceCollector_TopLevelStatement
	for _, dsg := range y.dsgsSlice {
		result = mergeDescendingBy(
			result,
			dsg.statements,
			(*CrossChunkReferenceCollector_TopLevelStatement).GetOriginalOrder,
			true,
			"Statements are recorded in the wrong order.")
	}
	return result
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.getPreferredChunk(BitSet)
func (y *cccmDsgCycle) preferredChunk(chunksWithImmovableReferences *smallBitSet) *JSChunk {
	if chunksWithImmovableReferences.IsEmpty() {
		return y.currentChunk
	} else if !y.allStmtsCanMove() {
		return y.currentChunk
	} else {
		return y.outer.graph.getSmallestCoveringSubtree(y.currentChunk, chunksWithImmovableReferences)
	}
}

// Java: CrossChunkCodeMotion.DeclarationStatementGroupCycle.allStatementsCanMove()
func (y *cccmDsgCycle) allStmtsCanMove() bool {
	for _, dsg := range y.dsgsSlice {
		if !dsg.allStmtsCanMove() {
			return false
		}
	}
	return true
}

type instanceofRef struct {
	node    rhino.Node
	chunk   *JSChunk
	dsg     *cccmDsg
	movable bool
}

func (r instanceofRef) targetChunk() *JSChunk {
	if r.movable {
		return r.dsg.currentChunk
	}
	return r.chunk
}

func cccmChunkName(c *JSChunk) string {
	if c == nil {
		return "null"
	}
	return c.String()
}

// isUndefinedTypeofGuardRef reports whether the reference node is the first Ref
// in an expression like 'undefined' != typeof Ref && x instanceof Ref.
//
// It's safe to ignore this kind of reference when moving the definition of
// Ref.
//
// Java: CrossChunkCodeMotion.isUndefinedTypeofGuardReference(Node)
func (c *CrossChunkCodeMotion) isUndefinedTypeofGuardRef(reference rhino.Node) bool {
	// reference => typeof => `!=`
	typeofOp, ok := ast.AsUnaryOf(reference.GetParent(), rhino.Token_TYPEOF)
	if !ok {
		return false
	}
	maybeTypeofGuard := typeofOp.Node().GetParent()
	if !maybeTypeofGuard.IsNil() && c.isExistenceTypeofGuardFor(maybeTypeofGuard, reference) {
		if andBin, ok := ast.AsBinaryOf(maybeTypeofGuard.GetParent(), rhino.Token_AND); ok {
			return c.isInstanceofFor(andBin.Right().Node(), reference)
		}
	}
	return false
}

// isExistenceTypeofGuardFor reports whether the expression is of the form
// 'undefined' != typeof Ref or 'function' == typeof Ref.
//
// expression is the expression being checked; reference the Ref node must
// be equivalent to this node.
//
// Java: CrossChunkCodeMotion.isExistenceTypeofGuardFor(Node, Node)
func (c *CrossChunkCodeMotion) isExistenceTypeofGuardFor(expression rhino.Node, reference rhino.Node) bool {
	bin, ok := ast.AsBinary(expression)
	if !ok {
		return false
	}
	var expected string
	switch bin.Token() {
	case rhino.Token_NE, rhino.Token_SHNE:
		expected = "undefined"
	case rhino.Token_EQ, rhino.Token_SHEQ:
		expected = "function"
	default:
		return false
	}
	strLit, ok := bin.Left().AsStringLit()
	if !ok || strLit.Value() != expected {
		return false
	}
	typeofOp, ok := ast.AsUnaryOf(bin.Right().Node(), rhino.Token_TYPEOF)
	return ok && typeofOp.Operand().Node().IsEquivalentTo(reference)
}

// isGuardedInstanceofRef reports whether the reference node is the second Ref
// in an expression like 'undefined' != typeof Ref && x instanceof Ref.
//
// It's safe to ignore this kind of reference when moving the definition of
// Ref.
//
// Java: CrossChunkCodeMotion.isGuardedInstanceofReference(Node)
func (c *CrossChunkCodeMotion) isGuardedInstanceofRef(reference rhino.Node) bool {
	instanceofNode := reference.GetParent()
	if c.isInstanceofFor(instanceofNode, reference) {
		if andBin, ok := ast.AsBinaryOf(instanceofNode.GetParent(), rhino.Token_AND); ok {
			return c.isExistenceTypeofGuardFor(andBin.Left().Node(), reference)
		}
	}
	return false
}

// isUnguardedInstanceofRef reports whether the reference is the right hand side
// of an instanceof and not guarded.
//
// Java: CrossChunkCodeMotion.isUnguardedInstanceofReference(Node)
func (c *CrossChunkCodeMotion) isUnguardedInstanceofRef(reference rhino.Node) bool {
	instanceofBin, ok := ast.AsBinaryOf(reference.GetParent(), rhino.Token_INSTANCEOF)
	if !ok || !instanceofBin.Right().Node().IsEquivalentTo(reference) {
		return false
	}
	if c.compiler.GetAstAnalyzer().MayHaveSideEffects(instanceofBin.Left().Node()) {
		return false
	}
	return !c.isGuardedInstanceofRef(reference)
}

// isInstanceofFor reports whether the expression is of the form x instanceof Ref.
//
// reference is the Ref node; it must be equivalent to this node.
//
// Java: CrossChunkCodeMotion.isInstanceofFor(Node, Node)
func (c *CrossChunkCodeMotion) isInstanceofFor(expression rhino.Node, reference rhino.Node) bool {
	bin, ok := ast.AsBinaryOf(expression, rhino.Token_INSTANCEOF)
	return ok && bin.Right().Node().IsEquivalentTo(reference)
}
