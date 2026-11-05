// Port of com/google/javascript/jscomp/ConstCheck.java @ closure-compiler v20261028

/*
 * Copyright 2004 The Closure Compiler Authors.
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

// ConstCheck verifies that constants are only assigned a value once.
//
// e.g.
//
//	XX = 3;    // error!
//	XX++;      // error!
//
// TODO(tbreisacher): Consider merging this with CheckAccessControls so that
// all const-related checks are in the same place.
type ConstCheck struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler             AbstractCompiler
	initializedConstants map[rhino.Node]struct{}
	onGlobalOp           func(kind constCheckGlobalOpKind, node, nameNode rhino.Node, var_ *Var)
}

type constCheckGlobalOpKind uint8

const (
	constCheckOpExternInit constCheckGlobalOpKind = iota + 1
	constCheckOpInitOrAssign
	constCheckOpIncDec
)

type constCheckGlobalOp struct {
	kind     constCheckGlobalOpKind
	eventIdx int
	node     rhino.Node
	nameNode rhino.Node
	var_     *Var
}

type constCheckWorkerState struct {
	rc     recordingCompiler
	sc     *SyntacticScopeCreator
	cc     *ConstCheck
	curIdx int
}

// Java: ConstCheck.CONST_REASSIGNED_VALUE_ERROR
var ConstCheck_CONST_REASSIGNED_VALUE_ERROR = DiagnosticType_Warning(
	"JSC_CONSTANT_REASSIGNED_VALUE_ERROR",
	"constant {0} assigned a value more than once.\nOriginal definition at {1}")

// NewConstCheck creates an instance.
//
// Java: ConstCheck(AbstractCompiler)
func NewConstCheck(compiler AbstractCompiler) *ConstCheck {
	c := &ConstCheck{}
	c.compiler = compiler
	c.initializedConstants = make(map[rhino.Node]struct{})
	return c
}

func (c *ConstCheck) addInitializedConstant(var_ *Var) bool {
	key := var_.nameNode
	if key.IsNil() {
		key = var_.scope.rootNode
	}
	if _, exists := c.initializedConstants[key]; exists {
		return false
	}
	c.initializedConstants[key] = struct{}{}
	return true
}

// Java: ConstCheck.process(Node, Node)
func (c *ConstCheck) Process(externs rhino.Node, root rhino.Node) {
	defer func() {
		clear(c.initializedConstants)
	}()
	if c.processParallel(externs, root) {
		return
	}
	sc := NewSyntacticScopeCreator(c.compiler)
	b := NodeTraversal_Builder{
		compiler:     c.compiler,
		callback:     c,
		scopeCreator: sc,
	}
	b.TraverseRoots(externs, root)
}

func (c *ConstCheck) processParallel(externs rhino.Node, root rhino.Node) bool {
	compiler := c.compiler
	if compiler == nil || compiler.IsDebugLoggingEnabled() {
		return false
	}
	scopeRoot, items, ok := collectExternAndRootScripts(externs, root)
	if !ok {
		return false
	}
	if len(items) == 0 {
		return true
	}
	numWorkers := parallelThreadCount(compiler, len(items))
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(compiler)
	lazyGlobal := newLazyGlobalScope(compiler, scopeRoot)
	eventsByItem := make([][]recordedCompilerEvent, len(items))
	globalOpsByItem := make([][]constCheckGlobalOp, len(items))

	panicsByItem, workerPanic := collectParallelWorkersDeferredPins(
		len(items),
		numWorkers,
		func(_ int, w *constCheckWorkerState) {
			initRecordingCompiler(&w.rc, compiler)
			w.sc = NewSyntacticScopeCreator(&w.rc)
			w.cc = &ConstCheck{
				compiler:             &w.rc,
				initializedConstants: make(map[rhino.Node]struct{}, 32),
			}
			w.cc.onGlobalOp = func(kind constCheckGlobalOpKind, node, nameNode rhino.Node, var_ *Var) {
				globalOpsByItem[w.curIdx] = append(globalOpsByItem[w.curIdx], constCheckGlobalOp{
					kind:     kind,
					eventIdx: len(w.rc.parallelRecorder.events),
					node:     node,
					nameNode: nameNode,
					var_:     var_,
				})
			}
		},
		func(w *constCheckWorkerState, idx int) {
			resetRecordingCompiler(&w.rc, compiler)
			w.curIdx = idx
			clear(w.cc.initializedConstants)
			item := items[idx]
			wt := acquireNodeTraversalDirect(&w.rc, w.cc, w.sc, false)
			defer releaseNodeTraversal(wt)
			func() {
				defer func() {
					if r := recover(); r != nil {
						wt.ThrowUnexpectedException(jutil.AsThrowable(r))
					}
				}()
				wt.lazyGlobalScope = lazyGlobal
				wt.InitTraversal(scopeRoot)
				wt.currentNode = scopeRoot
				wt.PushScope(scopeRoot)
				wt.TraverseBranch(item.script, item.parent)
				wt.PopScope()
			}()
		},
		func(w *constCheckWorkerState, idx int) {
			eventsByItem[idx] = w.rc.takeEvents()
		},
		nil,
		false,
	)

	for i := range items {
		evs := eventsByItem[i]
		evPos := 0
		for _, op := range globalOpsByItem[i] {
			if op.eventIdx > evPos && evPos < len(evs) {
				end := min(op.eventIdx, len(evs))
				replayCompilerEvents(compiler, evs[evPos:end])
				evPos = end
			}
			switch op.kind {
			case constCheckOpExternInit:
				c.addInitializedConstant(op.var_)
			case constCheckOpInitOrAssign:
				if !c.addInitializedConstant(op.var_) {
					c.ReportError(op.node, op.var_, op.nameNode.GetString())
				}
			case constCheckOpIncDec:
				c.ReportError(op.node, op.var_, op.nameNode.GetString())
			}
		}
		if evPos < len(evs) {
			replayCompilerEvents(compiler, evs[evPos:])
		}
		if p := panicsByItem[i]; p != nil {
			rethrowParallelPanic(compiler, p)
		}
	}
	if workerPanic != nil {
		rethrowParallelPanic(compiler, workerPanic)
	}
	return true
}

var constCheckVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_NAME)
	s.Add(rhino.Token_ASSIGN)
	s.Add(rhino.Token_ASSIGN_BITOR)
	s.Add(rhino.Token_ASSIGN_BITXOR)
	s.Add(rhino.Token_ASSIGN_BITAND)
	s.Add(rhino.Token_ASSIGN_LSH)
	s.Add(rhino.Token_ASSIGN_RSH)
	s.Add(rhino.Token_ASSIGN_URSH)
	s.Add(rhino.Token_ASSIGN_ADD)
	s.Add(rhino.Token_ASSIGN_SUB)
	s.Add(rhino.Token_ASSIGN_MUL)
	s.Add(rhino.Token_ASSIGN_DIV)
	s.Add(rhino.Token_ASSIGN_MOD)
	s.Add(rhino.Token_ASSIGN_EXPONENT)
	s.Add(rhino.Token_INC)
	s.Add(rhino.Token_DEC)
	return s
}()

func (*ConstCheck) onlyNameLeaves()            {}
func (*ConstCheck) visitTokens() *ast.TokenSet { return &constCheckVisitTokens }

var _ nodeTraversal_visitTokens = (*ConstCheck)(nil)

func lookupConstVar(t *NodeTraversal, nameNode rhino.Node) *Var {
	if id := nameNode.GetStringID(); id != 0 {
		return t.GetScope().GetVarByID(id)
	}
	return t.LookupVar(nameNode.GetString())
}

// Java: ConstCheck.visit(NodeTraversal, Node, Node)
func (c *ConstCheck) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		if NodeUtil_IsNameDeclaration(parent) {
			var_ := lookupConstVar(t, n)
			if c.IsConstant(var_, n) {
				// If a constant is declared in externs, add it to initializedConstants to indicate
				// that it is initialized externally.
				if n.IsFromExterns() {
					if c.onGlobalOp != nil && var_.IsGlobal() {
						c.onGlobalOp(constCheckOpExternInit, n, n, var_)
					} else {
						c.addInitializedConstant(var_)
					}
				} else if n.HasChildren() {
					if c.onGlobalOp != nil && var_.IsGlobal() {
						c.onGlobalOp(constCheckOpInitOrAssign, n, n, var_)
					} else if !c.addInitializedConstant(var_) {
						c.ReportError(n, var_, n.GetString())
					}
				}
			}
		}
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
		rhino.Token_ASSIGN_EXPONENT:
		lhs := n.GetFirstChild()
		if lhs.IsName() {
			var_ := lookupConstVar(t, lhs)
			if c.IsConstant(var_, lhs) {
				if c.onGlobalOp != nil && var_.IsGlobal() {
					c.onGlobalOp(constCheckOpInitOrAssign, n, lhs, var_)
				} else if !c.addInitializedConstant(var_) {
					c.ReportError(n, var_, lhs.GetString())
				}
			} else if var_ != nil && var_.IsGoogModuleExports() && !c.addInitializedConstant(var_) {
				c.compiler.Report(
					JSError_MakeWithNode(n, ConstCheck_CONST_REASSIGNED_VALUE_ERROR, "exports",
						jutil.FormatOptionalString(n.SourceFileName())))
			}
		}
	case rhino.Token_INC, rhino.Token_DEC:
		lhs := n.GetFirstChild()
		if lhs.IsName() {
			var_ := lookupConstVar(t, lhs)
			if c.IsConstant(var_, lhs) {
				if c.onGlobalOp != nil && var_.IsGlobal() {
					c.onGlobalOp(constCheckOpIncDec, n, lhs, var_)
				} else {
					c.ReportError(n, var_, lhs.GetString())
				}
			}
		}
	default:
	}
}

// IsConstant gets whether a variable is a constant initialized to a literal
// value at the point where it is declared.
//
// Java: ConstCheck.isConstant(Var, Node)
func (c *ConstCheck) IsConstant(var_ *Var, nameNode rhino.Node) bool {
	jutil.CheckNotNullMsg1(var_, "Found unexpected undeclared name %s", nameNode)
	return var_.IsConst() || var_.IsDeclaredOrInferredConst()
}

// ReportError reports a reassigned constant error.
//
// Java: ConstCheck.reportError(Node, Var, String)
func (c *ConstCheck) ReportError(n rhino.Node, var_ *Var, name string) {
	info := NodeUtil_GetBestJSDocInfo(n)
	if info == nil || !info.HasSuppression("const") {
		declNode := var_.GetNode()
		declaredPosition := jutil.FormatOptionalString(declNode.SourceFileName()) + ":" +
			strconv.Itoa(declNode.GetLineno())
		c.compiler.Report(JSError_MakeWithNode(n, ConstCheck_CONST_REASSIGNED_VALUE_ERROR, name, declaredPosition))
	}
}
