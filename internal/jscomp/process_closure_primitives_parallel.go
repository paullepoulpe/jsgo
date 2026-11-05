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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var processClosurePrimitivesExternVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_CALL)
	return s
}()

type closurePrimitivesActionKind uint8

const (
	cpActionInherits closurePrimitivesActionKind = iota + 1
	cpActionAddDep
	cpActionForwardDeclare
	cpActionBaseCall
)

type closurePrimitivesAction struct {
	kind     closurePrimitivesActionKind
	eventIdx int
	node     rhino.Node
	qname    string
}

type closurePrimitivesWorker struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	p        *ProcessClosurePrimitives
	actions  []closurePrimitivesAction
	rc       recordingCompiler
	curIndex int
}

func (*closurePrimitivesWorker) noVisitLeaves() {}
func (*closurePrimitivesWorker) visitTokens() *ast.TokenSet {
	return &processClosurePrimitivesVisitTokens
}

func (w *closurePrimitivesWorker) Visit(t *NodeTraversal, n rhino.Node, _ rhino.Node) {
	p := w.p
	switch n.GetToken() {
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		w.checkGoogFunctions(t, call)
		if callTarget, ok := call.Callee().AsGetProp(); ok && callTarget.Property() == "base" {
			w.actions = append(w.actions, closurePrimitivesAction{
				kind:     cpActionBaseCall,
				eventIdx: len(w.rc.parallelRecorder.events),
				node:     n,
			})
		}
		p.CheckPropertyRenameCall(call)
	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		if !t.InGlobalHoistScope() || n.IsFromExterns() {
			return
		}
		if !NodeUtil_IsFunctionDeclaration(n) && !NodeUtil_IsClassDeclaration(n) {
			return
		}
		var name string
		if fn, ok := ast.AsFunction(n); ok {
			fnName, _ := fn.Name()
			name = fnName.Ident()
		} else {
			clsName, _ := ast.MustClass(n).Name()
			name = clsName.Ident()
		}
		pn := p.closureModules.Get(name)
		if pn != nil && pn.IsGoogProvide() {
			diagnostic := ProcessClosurePrimitives_FUNCTION_NAMESPACE_ERROR
			if n.IsClass() {
				diagnostic = ProcessClosurePrimitives_CLASS_NAMESPACE_ERROR
			}
			p.compiler.Report(JSError_MakeWithNode(n, diagnostic, name, pn.String()))
		}
	case rhino.Token_EXPR_RESULT:
		assign, ok := ast.MustExprStmt(n).Expr().AsAssign()
		if !ok || n.IsFromExterns() {
			break
		}
		info := assign.Node().GetJSDocInfo()
		if info == nil || info.GetType() == nil {
			break
		}
		lhs, ok := assign.Target().Node().GetQualifiedNameString()
		if !ok {
			break
		}
		p.CheckPossibleGoogProvideInit(lhs, info, n)
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		first, _ := ast.MustNameDeclaration(n).FirstDeclarator()
		name, ok := first.AsName()
		if !ok {
			break
		}
		p.CheckPossibleGoogProvideInit(name.Ident(), n.GetJSDocInfo(), n)
	}
}

func (w *closurePrimitivesWorker) checkGoogFunctions(t *NodeTraversal, call ast.Call) {
	callee, ok := call.Callee().AsGetProp()
	if !ok {
		return
	}
	receiver, ok := callee.Object().AsName()
	if !ok || receiver.Ident() != ProcessClosurePrimitives_GOOG {
		return
	}
	p := w.p
	methodName := callee.Property()
	switch methodName {
	case "inherits":
		if call.NumArgs() == 2 {
			subClass, _ := call.FirstArg()
			superClass, _ := subClass.Next()
			if subClass.Node().IsUnscopedQualifiedName() && superClass.Node().IsUnscopedQualifiedName() {
				if subQname, ok := subClass.Node().GetQualifiedNameString(); ok {
					w.actions = append(w.actions, closurePrimitivesAction{
						kind:     cpActionInherits,
						eventIdx: len(w.rc.parallelRecorder.events),
						qname:    subQname,
					})
				}
			}
		}
	case "addDependency":
		if p.ValidateUnaliasablePrimitiveCall(t, call.Node(), methodName) {
			w.actions = append(w.actions, closurePrimitivesAction{
				kind:     cpActionAddDep,
				eventIdx: len(w.rc.parallelRecorder.events),
				node:     call.Node(),
			})
		}
	case "setCssNameMapping":
		_ = ProcessClosurePrimitives_ProcessSetCssNameMapping(p.compiler, call.Node(), call.Node().GetParent())
	case "forwardDeclare":
		if p.ValidatePrimitiveCallWithMessage(
			t,
			call.Node(),
			methodName,
			ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR) {
			n := call.Node()
			if !n.GetParent().IsExprResult() {
				return
			}
			typeDeclaration, ok := p.compiler.GetCodingConvention().IdentifyTypeDeclarationCall(n)
			if !ok {
				p.compiler.Report(
					JSError_MakeWithNode(
						n,
						ProcessClosurePrimitives_INVALID_FORWARD_DECLARE,
						"A single type could not identified for the goog.forwardDeclare statement"))
				return
			}
			w.actions = append(w.actions, closurePrimitivesAction{
				kind:     cpActionForwardDeclare,
				eventIdx: len(w.rc.parallelRecorder.events),
				qname:    typeDeclaration,
			})
		}
	case "weakUsage":
		p.ValidateWeakUsageCall(call)
	}
}

func (p *ProcessClosurePrimitives) processParallel(externs rhino.Node, root rhino.Node) bool {
	comp := p.compiler
	if comp == nil || comp.IsDebugLoggingEnabled() {
		return false
	}
	scopeRoot, items, ok := collectExternAndRootScripts(externs, root)
	if !ok || len(items) == 0 {
		return false
	}
	numWorkers := parallelThreadCount(comp, len(items))
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(comp)
	eventsByItem := make([][]recordedCompilerEvent, len(items))
	actionsByItem := make([][]closurePrimitivesAction, len(items))

	panicsByItem, workerPanic := collectParallelWorkersDeferredPins(
		len(items),
		numWorkers,
		func(_ int, w *closurePrimitivesWorker) {
			initRecordingCompiler(&w.rc, comp)
			w.p = &ProcessClosurePrimitives{
				compiler:       &w.rc,
				closureModules: p.closureModules,
			}
		},
		func(w *closurePrimitivesWorker, idx int) {
			resetRecordingCompiler(&w.rc, comp)
			w.actions = w.actions[:0]
			w.curIndex = idx
			item := items[idx]
			wt := acquireNodeTraversalDirect(&w.rc, w, ScopeCreator_ASSERT_NO_SCOPES_CREATED, false)
			if item.parent == externs {
				wt.visitTokens = &processClosurePrimitivesExternVisitTokens
			}
			defer releaseNodeTraversal(wt)
			traverseScriptUnderRoot(wt, scopeRoot, item.parent, item.script, nil)
		},
		func(w *closurePrimitivesWorker, idx int) {
			eventsByItem[idx] = w.rc.takeEvents()
			if len(w.actions) > 0 {
				acts := make([]closurePrimitivesAction, len(w.actions))
				copy(acts, w.actions)
				actionsByItem[idx] = acts
			}
		},
		nil,
		false,
	)

	for i := range items {
		events := eventsByItem[i]
		acts := actionsByItem[i]
		evPos := 0
		for j := range acts {
			act := &acts[j]
			if evPos < act.eventIdx {
				replayCompilerEvents(comp, events[evPos:act.eventIdx])
				evPos = act.eventIdx
			}
			switch act.kind {
			case cpActionInherits:
				if p.knownClosureSubclasses == nil {
					p.knownClosureSubclasses = make(map[string]struct{})
				}
				p.knownClosureSubclasses[act.qname] = struct{}{}
			case cpActionAddDep:
				p.ProcessAddDependency(ast.MustCall(act.node))
			case cpActionForwardDeclare:
				p.compiler.ForwardDeclareType(act.qname)
			case cpActionBaseCall:
				p.MaybeProcessClassBaseCall(ast.MustCall(act.node))
			}
		}
		if evPos < len(events) {
			replayCompilerEvents(comp, events[evPos:])
		}
		if pn := panicsByItem[i]; pn != nil {
			rethrowParallelPanic(comp, pn)
		}
	}
	if workerPanic != nil {
		rethrowParallelPanic(comp, workerPanic)
	}
	return true
}
