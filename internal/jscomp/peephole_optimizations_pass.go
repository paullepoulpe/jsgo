// Port of com/google/javascript/jscomp/PeepholeOptimizationsPass.java @ closure-compiler v20261021

/*
 * Copyright 2010 The Closure Compiler Authors.
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
	"math/bits"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	peepholeIDUndefined        = uint32(strpool.Intern("undefined"))
	peepholeIDWindow           = uint32(strpool.Intern("window"))
	peepholeIDObject           = uint32(strpool.Intern("Object"))
	peepholeIDArray            = uint32(strpool.Intern("Array"))
	peepholeIDBoolean          = uint32(strpool.Intern("Boolean"))
	peepholeIDString           = uint32(strpool.Intern("String"))
	peepholeIDParseInt         = uint32(strpool.Intern("parseInt"))
	peepholeIDParseFloat       = uint32(strpool.Intern("parseFloat"))
	peepholeIDMath             = uint32(strpool.Intern("Math"))
	peepholeIDNumber           = uint32(strpool.Intern("Number"))
	peepholeIDGoog             = uint32(strpool.Intern("goog"))
	peepholeIDConcat           = uint32(strpool.Intern("concat"))
	peepholeIDSubstr           = uint32(strpool.Intern("substr"))
	peepholeIDSubstring        = uint32(strpool.Intern("substring"))
	peepholeIDSlice            = uint32(strpool.Intern("slice"))
	peepholeIDDefineProperties = uint32(strpool.Intern("defineProperties"))
	peepholeIDBind             = uint32(strpool.Intern("bind"))
	peepholeIDCall             = uint32(strpool.Intern("call"))
	peepholeIDApply            = uint32(strpool.Intern("apply"))
)

// PeepholeOptimizationsPass is a compiler pass to run various peephole
// optimizations (e.g. constant folding, some useless code removal, some
// minimizations).
type PeepholeOptimizationsPass struct {
	compiler AbstractCompiler
	passName string
	// NOTE: Use a native array rather than a List to avoid creating iterators
	// for every node in the AST.
	peepholeOptimizations     []AbstractPeepholeOptimization
	pendingFeatures           *parser.FeatureSet
	visitTokenSet             ast.TokenSet
	activeMaskByToken         [256]uint16
	useActiveMask             bool
	standardOptimizationsOnly bool
	hasJ2clCallPass           bool
	hasOptimizeLetAndConst    bool
	retraverseOnChange        bool
}

var (
	_ NodeTraversal_Callback         = (*PeepholeOptimizationsPass)(nil)
	_ nodeTraversal_noShouldTraverse = (*PeepholeOptimizationsPass)(nil)
	_ nodeTraversal_visitTokens      = (*PeepholeOptimizationsPass)(nil)
)

func (*PeepholeOptimizationsPass) nodeTraversalNoShouldTraverse() {}

func (p *PeepholeOptimizationsPass) visitTokens() *ast.TokenSet {
	if p.useActiveMask {
		return &p.visitTokenSet
	}
	return nil
}

// Java: PeepholeOptimizationsPass.PeepCallback.shouldTraverse(NodeTraversal, Node, Node)
func (p *PeepholeOptimizationsPass) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

// NewPeepholeOptimizationsPass creates a peephole optimization pass that
// runs the given optimizations.
//
// Java: PeepholeOptimizationsPass(AbstractCompiler, String, AbstractPeepholeOptimization...)
func NewPeepholeOptimizationsPass(
	compiler AbstractCompiler, passName string, optimizations ...AbstractPeepholeOptimization) *PeepholeOptimizationsPass {
	return NewPeepholeOptimizationsPassOfList(compiler, passName, optimizations)
}

// Java: PeepholeOptimizationsPass(AbstractCompiler, String, List<AbstractPeepholeOptimization>)
func NewPeepholeOptimizationsPassOfList(
	compiler AbstractCompiler,
	passName string,
	optimizations []AbstractPeepholeOptimization) *PeepholeOptimizationsPass {
	p := &PeepholeOptimizationsPass{}
	p.compiler = compiler
	p.passName = passName
	p.peepholeOptimizations = append([]AbstractPeepholeOptimization(nil), optimizations...)
	p.retraverseOnChange = true
	p.initDispatchTable()
	return p
}

func (p *PeepholeOptimizationsPass) initDispatchTable() {
	if len(p.peepholeOptimizations) > 16 {
		p.useActiveMask = false
		p.standardOptimizationsOnly = false
		return
	}
	p.useActiveMask = true
	p.standardOptimizationsOnly = true
	p.hasJ2clCallPass = false
	p.hasOptimizeLetAndConst = false
	clear(p.activeMaskByToken[:])

	for idx, optim := range p.peepholeOptimizations {
		bit := uint16(1) << idx
		cfg := optim.peepholeConfig()
		if !cfg.isStandard {
			p.standardOptimizationsOnly = false
			for tok := range p.activeMaskByToken {
				p.activeMaskByToken[tok] |= bit
			}
			continue
		}
		if cfg.isJ2clCallPass {
			p.hasJ2clCallPass = true
		}
		if cfg.isOptimizeLetAndConst {
			p.hasOptimizeLetAndConst = true
		}
		for w, word := range cfg.tokens {
			for word != 0 {
				b := bits.TrailingZeros64(word)
				word &^= uint64(1) << b
				p.activeMaskByToken[(w<<6)|b] |= bit
			}
		}
	}

	p.visitTokenSet = ast.TokenSet{}
	p.visitTokenSet.Add(rhino.Token_SCRIPT)
	p.visitTokenSet.Add(rhino.Token_FUNCTION)
	for tok, mask := range p.activeMaskByToken {
		if mask != 0 {
			p.visitTokenSet.Add(rhino.Token(tok))
		}
	}
}

// Java: PeepholeOptimizationsPass.setRetraverseOnChange(boolean)
func (p *PeepholeOptimizationsPass) SetRetraverseOnChange(retraverse bool) {
	p.retraverseOnChange = retraverse
}

// Java: PeepholeOptimizationsPass.process(Node, Node)
func (p *PeepholeOptimizationsPass) Process(externs rhino.Node, root rhino.Node) {
	parallel := p.canCloneForParallel() && parallelThreadCount(p.compiler, 2) > 1
	if parallel {
		prewarmCompilerForParallelReads(p.compiler)
		restoreReg := enableParallelTypeRegistryReads(p.compiler)
		defer restoreReg()
	}

	p.BeginTraversal()

	// Repeat to an internal fixed point.
	for changedScopeNodes := p.compiler.GetChangeTracker().GetChangedScopeNodesForPass(p.passName); changedScopeNodes == nil || len(changedScopeNodes) != 0; changedScopeNodes = p.compiler.GetChangeTracker().GetChangedScopeNodesForPass(p.passName) {
		if changedScopeNodes == nil {
			// changedScopeNodes is null if this is the first run of peepholeOptimizationsPass.
			if !parallel || !p.processFullRootParallel(root) {
				NodeTraversal_Traverse(p.compiler, root, p)
			}
		} else if !parallel || !p.processChangedScopesParallel(changedScopeNodes) {
			NodeTraversal_TraverseScopeRoots(
				p.compiler, changedScopeNodes, p, false /* traverseNested */)
		}

		// Cancel the fixed point if requested.
		if !p.retraverseOnChange {
			break
		}
	}

	p.EndTraversal()
}

func (p *PeepholeOptimizationsPass) canCloneForParallel() bool {
	for _, optim := range p.peepholeOptimizations {
		if !optim.peepholeConfig().isStandard {
			return false
		}
	}
	return true
}

func (p *PeepholeOptimizationsPass) cloneForWorker(rc *recordingCompiler) *PeepholeOptimizationsPass {
	clonedOpts := make([]AbstractPeepholeOptimization, len(p.peepholeOptimizations))
	for i, optim := range p.peepholeOptimizations {
		clonedOpts[i] = optim.cloneForWorker()
	}
	cp := NewPeepholeOptimizationsPassOfList(rc, p.passName, clonedOpts)
	cp.retraverseOnChange = p.retraverseOnChange
	return cp
}

func enclosingScriptForPeephole(n rhino.Node) rhino.Node {
	for cur := n; !cur.IsNil(); cur = cur.GetParent() {
		if cur.IsScript() {
			return cur
		}
	}
	return rhino.NilNode
}

type peepholeScriptBucket struct {
	origIndices []int
}

func (p *PeepholeOptimizationsPass) processFullRootParallel(root rhino.Node) bool {
	scripts, ok := collectRootScripts(root)
	if !ok {
		return false
	}
	numWorkers := parallelThreadCount(p.compiler, len(scripts))
	if numWorkers <= 1 {
		return false
	}
	runParallelWork(p.compiler, numWorkers, len(scripts), func(_ int) (func(int, *recordingCompiler), func()) {
		var workerRC *recordingCompiler
		var workerPass *PeepholeOptimizationsPass

		cleanup := func() {
			if workerPass != nil {
				workerPass.EndTraversal()
			}
		}
		runItem := func(idx int, rc *recordingCompiler) {
			if workerPass == nil || workerRC != rc {
				workerRC = rc
				workerPass = p.cloneForWorker(rc)
				workerPass.BeginTraversal()
			}
			b := NodeTraversal_Builder{
				compiler: rc,
				callback: workerPass,
			}
			b.traverseScriptUnderRoot(root, root, scripts[idx], nil)
		}
		return runItem, cleanup
	}, false)
	return true
}

func (p *PeepholeOptimizationsPass) processChangedScopesParallel(changedScopeNodes []rhino.Node) bool {
	if len(changedScopeNodes) < 2 {
		return false
	}
	bucketByScript := make(map[rhino.Node]int, len(changedScopeNodes))
	var buckets []peepholeScriptBucket
	for origIdx, scopeNode := range changedScopeNodes {
		script := enclosingScriptForPeephole(scopeNode)
		if script.IsNil() {
			return false
		}
		bIdx, exists := bucketByScript[script]
		if !exists {
			bIdx = len(buckets)
			bucketByScript[script] = bIdx
			buckets = append(buckets, peepholeScriptBucket{origIndices: []int{origIdx}})
		} else {
			buckets[bIdx].origIndices = append(buckets[bIdx].origIndices, origIdx)
		}
	}
	numWorkers := parallelThreadCount(p.compiler, len(buckets))
	if numWorkers <= 1 {
		return false
	}

	eventsByOrigIdx := make([][]recordedCompilerEvent, len(changedScopeNodes))
	panicsByOrigIdx := make([]any, len(changedScopeNodes))

	runParallelWork(p.compiler, numWorkers, len(buckets), func(_ int) (func(int, *recordingCompiler), func()) {
		var workerRC *recordingCompiler
		var workerPass *PeepholeOptimizationsPass

		cleanup := func() {
			if workerPass != nil {
				workerPass.EndTraversal()
			}
		}
		runItem := func(bIdx int, rc *recordingCompiler) {
			if workerPass == nil || workerRC != rc {
				workerRC = rc
				workerPass = p.cloneForWorker(rc)
				workerPass.BeginTraversal()
			}
			for _, origIdx := range buckets[bIdx].origIndices {
				resetRecordingCompiler(rc, p.compiler)
				panicked := runParallelItemIndexed(rc, origIdx, panicsByOrigIdx, func(rc *recordingCompiler, itemIdx int) {
					b := NodeTraversal_Builder{
						compiler: rc,
						callback: workerPass,
					}
					b.traverseScopedRoot(changedScopeNodes[itemIdx], false)
				})
				eventsByOrigIdx[origIdx] = rc.takeEvents()
				if panicked {
					workerPass.pendingFeatures = nil
					return
				}
			}
		}
		return runItem, cleanup
	}, false)

	for i := range changedScopeNodes {
		if len(eventsByOrigIdx[i]) != 0 {
			replayCompilerEvents(p.compiler, eventsByOrigIdx[i])
		}
		if r := panicsByOrigIdx[i]; r != nil {
			if ie, ok := r.(parallelInternalError); ok {
				p.compiler.ThrowInternalError(ie.msg, ie.cause)
			}
			panic(r)
		}
	}
	return true
}

// Java: PeepholeOptimizationsPass.PeepCallback.visit(NodeTraversal, Node, Node)
func (p *PeepholeOptimizationsPass) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	p.visitOptimize(n)
	if p.pendingFeatures != nil && (n.IsScript() || n == t.scopeTraversalRoot) {
		NodeUtil_AddFeaturesToScript(t.GetCurrentScript(), p.pendingFeatures, p.compiler)
		p.pendingFeatures = nil
	}
}

func (p *PeepholeOptimizationsPass) visitOptimize(n rhino.Node) {
	if !p.useActiveMask {
		currentNode := n
		for _, optim := range p.peepholeOptimizations {
			currentNode = optim.OptimizeSubtree(currentNode)
			if currentNode.IsNil() {
				return
			}
		}
		return
	}

	currentNode := n
	tok := currentNode.GetToken()
	mask := p.activeMaskByToken[tok]
	if mask == 0 {
		return
	}
	if p.standardOptimizationsOnly {
		switch tok {
		case rhino.Token_NAME:
			if currentNode.GetStringID() != peepholeIDUndefined {
				return
			}
		case rhino.Token_GETPROP:
			switch first := currentNode.GetFirstChild(); first.GetToken() {
			case rhino.Token_NAME:
				if first.GetStringID() != peepholeIDWindow {
					return
				}
			case rhino.Token_OBJECTLIT, rhino.Token_ARRAYLIT, rhino.Token_STRINGLIT:
			default:
				return
			}
		case rhino.Token_GETELEM:
			switch currentNode.GetFirstChild().GetToken() {
			case rhino.Token_OBJECTLIT, rhino.Token_ARRAYLIT, rhino.Token_STRINGLIT:
			default:
				return
			}
		case rhino.Token_CALL:
			if !p.hasJ2clCallPass {
				callee := currentNode.GetFirstChild()
				switch callee.GetToken() {
				case rhino.Token_NAME:
					switch callee.GetStringID() {
					case peepholeIDObject, peepholeIDArray, peepholeIDBoolean, peepholeIDString, peepholeIDParseInt, peepholeIDParseFloat:
					default:
						return
					}
				case rhino.Token_GETPROP:
					recv := callee.GetFirstChild()
					switch recv.GetToken() {
					case rhino.Token_STRINGLIT, rhino.Token_ARRAYLIT:
					case rhino.Token_NAME:
						switch recv.GetStringID() {
						case peepholeIDObject, peepholeIDArray, peepholeIDMath, peepholeIDNumber, peepholeIDGoog:
						default:
							switch callee.GetStringID() {
							case peepholeIDConcat, peepholeIDSubstr, peepholeIDSubstring, peepholeIDSlice:
							default:
								return
							}
						}
					default:
						switch callee.GetStringID() {
						case peepholeIDConcat, peepholeIDSubstr, peepholeIDSubstring, peepholeIDSlice,
							peepholeIDDefineProperties, peepholeIDBind, peepholeIDCall, peepholeIDApply:
						default:
							return
						}
					}
				case rhino.Token_CALL:
				default:
					return
				}
			}
		case rhino.Token_EXPR_RESULT:
			expr := currentNode.GetFirstChild()
			switch expr.GetToken() {
			case rhino.Token_ASSIGN,
				rhino.Token_ASSIGN_BITOR,
				rhino.Token_ASSIGN_BITXOR,
				rhino.Token_ASSIGN_BITAND,
				rhino.Token_ASSIGN_LSH,
				rhino.Token_ASSIGN_RSH,
				rhino.Token_ASSIGN_URSH,
				rhino.Token_ASSIGN_ADD,
				rhino.Token_ASSIGN_SUB,
				rhino.Token_ASSIGN_MUL,
				rhino.Token_ASSIGN_DIV,
				rhino.Token_ASSIGN_MOD,
				rhino.Token_ASSIGN_EXPONENT,
				rhino.Token_ASSIGN_OR,
				rhino.Token_ASSIGN_AND,
				rhino.Token_ASSIGN_COALESCE:
				// An assignment is a statement the standard optimizations
				// leave alone only if it certainly changes state: dead-code
				// removal drops a side-effect-free one (`[1].a = 2`), as
				// AstAnalyzer.checkAssignmentForChangeState decides.
				if peepholeAssignmentCertainlyChangesState(expr) {
					return
				}
			case rhino.Token_INC,
				rhino.Token_DEC,
				rhino.Token_DELPROP,
				rhino.Token_YIELD,
				rhino.Token_AWAIT,
				rhino.Token_DYNAMIC_IMPORT:
				return
			}
		case rhino.Token_ASSIGN:
			left := currentNode.GetFirstChild()
			right := currentNode.GetLastChild()
			if !(left.IsName() && right.IsName()) && !(left.IsDestructuringPattern() && !left.HasChildren()) {
				switch right.GetToken() {
				case rhino.Token_ADD,
					rhino.Token_BITAND,
					rhino.Token_BITOR,
					rhino.Token_BITXOR,
					rhino.Token_DIV,
					rhino.Token_EXPONENT,
					rhino.Token_LSH,
					rhino.Token_MOD,
					rhino.Token_MUL,
					rhino.Token_RSH,
					rhino.Token_SUB,
					rhino.Token_URSH:
				default:
					return
				}
			}
		case rhino.Token_VAR:
			if first := currentNode.GetFirstChild(); first.IsNil() || !first.IsDestructuringLhs() {
				return
			}
		case rhino.Token_LET, rhino.Token_CONST:
			if !p.hasOptimizeLetAndConst {
				if first := currentNode.GetFirstChild(); first.IsNil() || !first.IsDestructuringLhs() {
					return
				}
			}
		case rhino.Token_NEW:
			if !currentNode.GetFirstChild().IsName() {
				return
			}
		case rhino.Token_OBJECTLIT, rhino.Token_ARRAYLIT:
			if currentNode.GetFirstChild().IsNil() {
				return
			}
		case rhino.Token_VOID:
			if first := currentNode.GetFirstChild(); !first.IsNil() && first.IsNumber() && first.GetDouble() == 0 {
				return
			}
		case rhino.Token_NEG:
			if first := currentNode.GetFirstChild(); !first.IsNil() && first.IsNumber() {
				return
			}
		}
	}

	for mask != 0 {
		idx := bits.TrailingZeros16(mask)
		bit := uint16(1) << idx
		mask &^= bit
		currentNode = p.peepholeOptimizations[idx].OptimizeSubtree(currentNode)
		if currentNode.IsNil() {
			return
		}
		if newTok := currentNode.GetToken(); newTok != tok {
			tok = newTok
			higherMask := ^((bit << 1) - 1)
			mask = p.activeMaskByToken[tok] & higherMask
		}
	}
}

// BeginTraversal makes sure that all the optimizations have the current
// compiler so they can report errors.
//
// Java: PeepholeOptimizationsPass.beginTraversal()
func (p *PeepholeOptimizationsPass) BeginTraversal() {
	p.pendingFeatures = nil
	for _, optimization := range p.peepholeOptimizations {
		optimization.bindPendingFeatures(&p.pendingFeatures)
		optimization.BeginTraversal(p.compiler)
	}
	p.initDispatchTable()
}

// EndTraversal ends the traversal.
//
// Java: PeepholeOptimizationsPass.endTraversal()
func (p *PeepholeOptimizationsPass) EndTraversal() {
	jutil.CheckStateMsg1(
		p.pendingFeatures == nil || p.pendingFeatures.IsEmpty(),
		"Expected getNewFeatures() to be empty in endTraversal() but found %s",
		p.pendingFeatures)
	for _, optimization := range p.peepholeOptimizations {
		optimization.bindPendingFeatures(nil)
		optimization.EndTraversal()
	}
}

// peepholeAssignmentCertainlyChangesState reports whether the assignment
// assign has side effects whatever its operands are: the conditions under
// which AstAnalyzer.checkAssignmentForChangeState returns true before
// looking at the operands' own side effects (a name target, or a property
// of a receiver that is neither a local value nor rooted in a literal).
// When it returns false, the assignment may be side-effect free and the
// standard optimizations must see it.
//
// Go-only: performance — the PeepCallback prefilter.
func peepholeAssignmentCertainlyChangesState(assign rhino.Node) bool {
	target := assign.GetFirstChild()
	if target.IsName() {
		return true
	}
	if !NodeUtil_IsNormalGet(target) {
		return false
	}
	current := target.GetFirstChild()
	if NodeUtil_EvaluatesToLocalValue(current) {
		return false
	}
	for NodeUtil_IsNormalGet(current) {
		current = current.GetFirstChild()
	}
	return !NodeUtil_IsLiteralValue(current, true)
}
