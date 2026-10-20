// Port of com/google/javascript/jscomp/PeepholeOptimizationsPass.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
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
	featureOptimBases         []*AbstractPeepholeOptimizationBase
	activeMaskByToken         [256]uint16
	useActiveMask             bool
	standardOptimizationsOnly bool
	retraverseOnChange        bool
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

func (p *PeepholeOptimizationsPass) addTokensToMask(bit uint16, tokens ...rhino.Token) {
	for _, tok := range tokens {
		p.activeMaskByToken[tok] |= bit
	}
}

func (p *PeepholeOptimizationsPass) initDispatchTable() {
	p.featureOptimBases = p.featureOptimBases[:0]
	if len(p.peepholeOptimizations) > 16 {
		p.useActiveMask = false
		p.standardOptimizationsOnly = false
		for _, optim := range p.peepholeOptimizations {
			p.featureOptimBases = append(p.featureOptimBases, optim.GetAbstractPeepholeOptimizationBase())
		}
		return
	}
	p.useActiveMask = true
	p.standardOptimizationsOnly = true
	clear(p.activeMaskByToken[:])

	for idx, optim := range p.peepholeOptimizations {
		bit := uint16(1) << idx
		switch opt := optim.(type) {
		case *MinimizeExitPoints:
			p.addTokensToMask(bit,
				rhino.Token_LABEL,
				rhino.Token_FOR,
				rhino.Token_FOR_IN,
				rhino.Token_FOR_OF,
				rhino.Token_FOR_AWAIT_OF,
				rhino.Token_WHILE,
				rhino.Token_DO,
				rhino.Token_BLOCK,
				rhino.Token_SWITCH,
				rhino.Token_SCRIPT,
			)
		case *PeepholeMinimizeConditions:
			p.addTokensToMask(bit,
				rhino.Token_THROW,
				rhino.Token_RETURN,
				rhino.Token_NOT,
				rhino.Token_IF,
				rhino.Token_EXPR_RESULT,
				rhino.Token_HOOK,
				rhino.Token_WHILE,
				rhino.Token_DO,
				rhino.Token_FOR,
				rhino.Token_BLOCK,
			)
		case *PeepholeSubstituteAlternateSyntax:
			p.addTokensToMask(bit,
				rhino.Token_ASSIGN_SUB,
				rhino.Token_NEW,
				rhino.Token_CALL,
				rhino.Token_RETURN,
				rhino.Token_NAME,
				rhino.Token_GETPROP,
				rhino.Token_TEMPLATELIT,
			)
			if opt.late {
				p.addTokensToMask(bit,
					rhino.Token_TRUE,
					rhino.Token_FALSE,
					rhino.Token_ARRAYLIT,
					rhino.Token_AND,
					rhino.Token_OR,
					rhino.Token_BITOR,
					rhino.Token_BITXOR,
					rhino.Token_BITAND,
					rhino.Token_COALESCE,
					rhino.Token_MUL,
				)
			} else {
				p.addTokensToMask(bit, rhino.Token_EXPR_RESULT)
			}
		case *PeepholeReplaceKnownMethods:
			p.addTokensToMask(bit, rhino.Token_CALL)
		case *PeepholeRemoveDeadCode:
			p.featureOptimBases = append(p.featureOptimBases, &opt.AbstractPeepholeOptimizationBase)
			p.addTokensToMask(bit,
				rhino.Token_ASSIGN,
				rhino.Token_COMMA,
				rhino.Token_SCRIPT,
				rhino.Token_BLOCK,
				rhino.Token_EXPR_RESULT,
				rhino.Token_HOOK,
				rhino.Token_SWITCH,
				rhino.Token_IF,
				rhino.Token_FOR,
				rhino.Token_DO,
				rhino.Token_TRY,
				rhino.Token_LABEL,
				rhino.Token_ARRAY_PATTERN,
				rhino.Token_OBJECT_PATTERN,
				rhino.Token_VAR,
				rhino.Token_CONST,
				rhino.Token_LET,
				rhino.Token_DEFAULT_VALUE,
				rhino.Token_OPTCHAIN_GETPROP,
				rhino.Token_OPTCHAIN_CALL,
				rhino.Token_OPTCHAIN_GETELEM,
			)
		case *J2clEqualitySameRewriterPass:
			if opt.shouldRunJ2clPasses {
				p.addTokensToMask(bit, rhino.Token_CALL)
			}
		case *J2clStringValueOfRewriterPass:
			if opt.shouldRunJ2clPasses {
				p.addTokensToMask(bit, rhino.Token_CALL)
			}
		case *J2clUndefinedChecksRewriterPass:
			if opt.shouldRunJ2clPasses {
				p.addTokensToMask(bit, rhino.Token_CALL)
			}
		case *PeepholeFoldConstants:
			p.addTokensToMask(bit,
				rhino.Token_OPTCHAIN_CALL,
				rhino.Token_CALL,
				rhino.Token_NEW,
				rhino.Token_TYPEOF,
				rhino.Token_ITER_SPREAD,
				rhino.Token_ARRAYLIT,
				rhino.Token_OBJECTLIT,
				rhino.Token_NOT,
				rhino.Token_POS,
				rhino.Token_NEG,
				rhino.Token_BITNOT,
				rhino.Token_VOID,
				rhino.Token_OPTCHAIN_GETPROP,
				rhino.Token_GETPROP,
				rhino.Token_TEMPLATELIT,
				rhino.Token_GETELEM,
				rhino.Token_OPTCHAIN_GETELEM,
				rhino.Token_INSTANCEOF,
				rhino.Token_AND,
				rhino.Token_OR,
				rhino.Token_COALESCE,
				rhino.Token_LSH,
				rhino.Token_RSH,
				rhino.Token_URSH,
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
				rhino.Token_ADD,
				rhino.Token_SUB,
				rhino.Token_DIV,
				rhino.Token_MOD,
				rhino.Token_MUL,
				rhino.Token_EXPONENT,
				rhino.Token_BITAND,
				rhino.Token_BITOR,
				rhino.Token_BITXOR,
				rhino.Token_LT,
				rhino.Token_GT,
				rhino.Token_LE,
				rhino.Token_GE,
				rhino.Token_EQ,
				rhino.Token_NE,
				rhino.Token_SHEQ,
				rhino.Token_SHNE,
			)
			if opt.late {
				p.addTokensToMask(bit, rhino.Token_ASSIGN)
			}
		case *PeepholeCollectPropertyAssignments:
			p.addTokensToMask(bit,
				rhino.Token_SCRIPT,
				rhino.Token_BLOCK,
			)
		case *ExploitAssigns:
			p.addTokensToMask(bit,
				rhino.Token_SCRIPT,
				rhino.Token_BLOCK,
				rhino.Token_MODULE_BODY,
			)
		case *StatementFusion:
			p.addTokensToMask(bit, rhino.Token_BLOCK)
		case *OptimizeLetAndConstPeephole:
			p.featureOptimBases = append(p.featureOptimBases, &opt.AbstractPeepholeOptimizationBase)
			p.addTokensToMask(bit,
				rhino.Token_LET,
				rhino.Token_CONST,
			)
		default:
			p.standardOptimizationsOnly = false
			p.featureOptimBases = append(p.featureOptimBases, optim.GetAbstractPeepholeOptimizationBase())
			for tok := range p.activeMaskByToken {
				p.activeMaskByToken[tok] |= bit
			}
		}
	}
}

// Java: PeepholeOptimizationsPass.setRetraverseOnChange(boolean)
func (p *PeepholeOptimizationsPass) SetRetraverseOnChange(retraverse bool) {
	p.retraverseOnChange = retraverse
}

// Java: PeepholeOptimizationsPass.process(Node, Node)
func (p *PeepholeOptimizationsPass) Process(externs *rhino.Node, root *rhino.Node) {
	p.BeginTraversal()
	cb := newPeepholeOptimizationsPass_PeepCallback(p)

	// Repeat to an internal fixed point.
	for changedScopeNodes := p.compiler.GetChangeTracker().GetChangedScopeNodesForPass(p.passName); changedScopeNodes == nil || len(changedScopeNodes) != 0; changedScopeNodes = p.compiler.GetChangeTracker().GetChangedScopeNodesForPass(p.passName) {

		if changedScopeNodes == nil {
			// changedScopeNodes is null if this is the first run of peepholeOptimizationsPass.
			NodeTraversal_Traverse(p.compiler, root, cb)
		} else {
			NodeTraversal_TraverseScopeRoots(
				p.compiler, changedScopeNodes, cb, false /* traverseNested */)
		}

		// Cancel the fixed point if requested.
		if !p.retraverseOnChange {
			break
		}
	}

	p.EndTraversal()
}

// peepholeOptimizationsPass_PeepCallback is the inner class PeepCallback.
type peepholeOptimizationsPass_PeepCallback struct {
	NodeTraversal_AbstractScopedCallbackBase
	this *PeepholeOptimizationsPass
}

func newPeepholeOptimizationsPass_PeepCallback(outer *PeepholeOptimizationsPass) *peepholeOptimizationsPass_PeepCallback {
	return &peepholeOptimizationsPass_PeepCallback{this: outer}
}

// Java: PeepholeOptimizationsPass.PeepCallback.visit(NodeTraversal, Node, Node)
func (c *peepholeOptimizationsPass_PeepCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	p := c.this
	if !p.useActiveMask {
		currentNode := n
		for _, optim := range p.peepholeOptimizations {
			currentNode = optim.OptimizeSubtree(currentNode)
			if currentNode == nil {
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
		if tok == rhino.Token_NAME {
			if currentNode.GetString() != "undefined" {
				return
			}
		} else if tok == rhino.Token_GETPROP {
			switch first := currentNode.GetFirstChild(); first.GetToken() {
			case rhino.Token_NAME:
				if first.GetString() != "window" {
					return
				}
			case rhino.Token_OBJECTLIT, rhino.Token_ARRAYLIT, rhino.Token_STRINGLIT:
			default:
				return
			}
		}
	}

	for mask != 0 {
		idx := bits.TrailingZeros16(mask)
		bit := uint16(1) << idx
		mask &^= bit
		currentNode = p.peepholeOptimizations[idx].OptimizeSubtree(currentNode)
		if currentNode == nil {
			return
		}
		if newTok := currentNode.GetToken(); newTok != tok {
			tok = newTok
			higherMask := ^((bit << 1) - 1)
			mask = p.activeMaskByToken[tok] & higherMask
		}
	}
}

// Java: PeepholeOptimizationsPass.PeepCallback.enterScope(NodeTraversal)
func (c *peepholeOptimizationsPass_PeepCallback) EnterScope(t *NodeTraversal) {}

// Java: PeepholeOptimizationsPass.PeepCallback.exitScope(NodeTraversal)
func (c *peepholeOptimizationsPass_PeepCallback) ExitScope(t *NodeTraversal) {
	// Call updateFeatures() here, instead of immediately after
	// `optim.optimizeSubtree(currentNode)` in visit, to avoid repeating work for
	// every AST node x every peephole optimization. In practice, profiling shows a small
	// but measurable improvement for large projects. See comments on cl/856752797.
	c.UpdateFeatures(t)
}

// Java: PeepholeOptimizationsPass.PeepCallback.updateFeatures(NodeTraversal)
func (c *peepholeOptimizationsPass_PeepCallback) UpdateFeatures(t *NodeTraversal) {
	for _, base := range c.this.featureOptimBases {
		newFeatures := base.newFeatures
		if !newFeatures.IsEmpty() {
			// Go-only: language-adapter — FeatureSet.with(Set<Feature>) takes the Go EnumSet form.
			NodeUtil_AddFeaturesToScript(
				t.GetCurrentScript(),
				parser.FeatureSet_BARE_MINIMUM.WithEnumSet(jutil.NewEnumSet(newFeatures.ToSlice()...)),
				c.this.compiler)
			base.ClearNewFeatures()
		}
	}
}

// BeginTraversal makes sure that all the optimizations have the current
// compiler so they can report errors.
//
// Java: PeepholeOptimizationsPass.beginTraversal()
func (p *PeepholeOptimizationsPass) BeginTraversal() {
	for _, optimization := range p.peepholeOptimizations {
		optimization.BeginTraversal(p.compiler)
	}
	p.initDispatchTable()
}

// EndTraversal ends the traversal.
//
// Java: PeepholeOptimizationsPass.endTraversal()
func (p *PeepholeOptimizationsPass) EndTraversal() {
	for _, optimization := range p.peepholeOptimizations {
		optimization.EndTraversal()
	}
}
