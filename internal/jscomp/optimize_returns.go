// Port of com/google/javascript/jscomp/OptimizeReturns.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// OptimizeReturns is a compiler pass to optimize function return results.
// Currently this pass looks for results that are completely unused and
// rewrites them as: "return x()" --> "x(); return"
//
// Future work: expand this to look for use context to avoid unneeded type
// coercion:
//
// "return x.toString()" --> "return x"
//
// "return !!x" --> "return x"
type OptimizeReturns struct {
	compiler AbstractCompiler

	// Allocated & cleaned up by process()
	decisionsLog diagnostic.LogFile
	fnNodes      []rhino.Node
}

// Java: OptimizeReturns(AbstractCompiler)
func NewOptimizeReturns(compiler AbstractCompiler) *OptimizeReturns {
	return &OptimizeReturns{compiler: compiler}
}

// Process is CompilerPass.process.
//
// Go-only: language-adapter — the two Java overloads process(Node, Node)
// (CompilerPass) and process(Node, Node, ReferenceMap)
// (CallGraphCompilerPass) cannot share a Go name; the second one is
// ProcessWithReferenceMap, and AsCallGraphCompilerPass adapts it.
//
// Java: OptimizeReturns.process(Node, Node)
func (o *OptimizeReturns) Process(externs rhino.Node, root rhino.Node) {
	OptimizeCalls_Builder_().
		SetCompiler(o.compiler).
		SetConsiderExterns(false).
		AddPass(o.AsCallGraphCompilerPass()).
		Build().
		Process(externs, root)
}

// AsCallGraphCompilerPass returns this pass as an
// OptimizeCalls.CallGraphCompilerPass (what Java does by passing `this`).
//
// Go-only: language-adapter — see Process.
func (o *OptimizeReturns) AsCallGraphCompilerPass() OptimizeCalls_CallGraphCompilerPass {
	return OptimizeCalls_CallGraphCompilerPassFunc(o.ProcessWithReferenceMap)
}

// Java: OptimizeReturns.process(Node, Node, ReferenceMap)
func (o *OptimizeReturns) ProcessWithReferenceMap(externs rhino.Node, root rhino.Node, definitions *OptimizeCalls_ReferenceMap) {
	logFile := o.compiler.CreateOrReopenIndexedLog("OptimizeReturns", "decisions.log")
	defer logFile.Close()
	defer func() {
		o.decisionsLog = nil
		clear(o.fnNodes)
		o.fnNodes = o.fnNodes[:0]
	}()
	o.decisionsLog = logFile // avoid passing the log file through a bunch of methods
	// Find all function nodes whose callers ignore the return values.
	var toOptimize [][]rhino.Node

	// Find all the candidates before modifying the AST.
	for key, refs := range definitions.GetNameReferences() {
		if o.isCandidate(key, refs) {
			if o.decisionsLog.IsLogging() {
				o.decisionsLog.LogFormat("name %s\tremoving return value", key)
			}
			toOptimize = append(toOptimize, refs)
		}
	}

	for key, refs := range definitions.GetPropReferences() {
		if o.isCandidate(key, refs) {
			if o.decisionsLog.IsLogging() {
				o.decisionsLog.LogFormat("property %s\tremoving return value", key)
			}
			toOptimize = append(toOptimize, refs)
		}
	}

	// Now modify the AST
	for _, refs := range toOptimize {
		fns := OptimizeCalls_ReferenceMap_AppendFunctionNodes(o.fnNodes[:0], refs)
		o.fnNodes = fns
		for _, fn := range fns {
			o.rewriteReturns(ast.MustFunction(fn))
		}
	}
}

// isCandidate reports whether a reference set is a candidate for return-value
// removal:
//   - all call sites are known (not aliased, not exported)
//   - all call sites ignore the return value
//   - there is at least one known function definition
//   - there is at least one use
//
// NOTE: unknown definitions are allowed, as only known definitions will be
// rewritten.
//
// Java: OptimizeReturns.isCandidate(String, List<Node>)
func (o *OptimizeReturns) isCandidate(name string, refs []rhino.Node) bool {
	if !OptimizeCalls_MayBeOptimizableName(o.compiler, name) {
		if o.decisionsLog.IsLogging() {
			o.decisionsLog.LogFormat("%s\tnot an optimizable name", name)
		}
		return false
	}

	seenCandidateDefinition := false
	seenUse := false
	for _, n := range refs {
		// Assume indirect definitions references use the result
		if OptimizeCalls_ReferenceMap_IsCallTarget(n) || OptimizeCalls_ReferenceMap_IsOptChainCallTarget(n) {
			callNode := OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(n)
			if NodeUtil_IsExpressionResultUsed(callNode) {
				if o.decisionsLog.IsLogging() { // avoid build location string when not logging
					o.decisionsLog.LogFormat("%s\treturn value used: %s", name, callNode.GetLocation())
				}
				return false
			}
			seenUse = true
		} else if o.isCandidateDefinition(n) {
			seenCandidateDefinition = true
		} else if !OptimizeCalls_IsAllowedReference(n) {
			// If this isn't a non-aliasing reference (typeof, instanceof, etc)
			// then there is nothing that can be done.
			if o.decisionsLog.IsLogging() { // avoid build location string when not logging
				o.decisionsLog.LogFormat("%s\tdisallowed reference: %s", name, n.GetLocation())
			}
			return false
		}
	}

	if !seenUse {
		if o.decisionsLog.IsLogging() {
			o.decisionsLog.LogFormat("%s\tno usage seen", name)
		}
		return false
	}
	if !seenCandidateDefinition {
		if o.decisionsLog.IsLogging() {
			o.decisionsLog.LogFormat("%s\tno definition seen", name)
		}
		return false
	}
	return true
}

// Java: OptimizeReturns.isCandidateDefinition(Node)
func (o *OptimizeReturns) isCandidateDefinition(n rhino.Node) bool {
	parent := n.GetParent()
	if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
		return true
	}
	if OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n) {
		return optimizeReturns_isCandidateFunction(ast.MustAssign(parent).Value())
	}
	if name, ok := ast.AsName(n); ok {
		init, hasInit := name.Init()
		return hasInit && optimizeReturns_isCandidateFunction(init)
	}
	return n.IsMemberFunctionDef() && parent.IsClassMembers()
}

// Java: OptimizeReturns.isCandidateFunction(Node)
func optimizeReturns_isCandidateFunction(e ast.Expr) bool {
	switch e.Token() {
	case rhino.Token_FUNCTION:
		// Named function expression can be recursive, this creates an alias of the name, meaning
		// it might be used in an unexpected way.
		return !NodeUtil_IsNamedFunctionExpression(e.Node())
	case rhino.Token_COMMA:
		return optimizeReturns_isCandidateFunction(e.MustBinary().Right())
	case rhino.Token_CAST:
		return optimizeReturns_isCandidateFunction(e.MustCast().Expr())
	case rhino.Token_HOOK:
		h := e.MustHook()
		return optimizeReturns_isCandidateFunction(h.Then()) && optimizeReturns_isCandidateFunction(h.Else())
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		b := e.MustBinary()
		return optimizeReturns_isCandidateFunction(b.Left()) && optimizeReturns_isCandidateFunction(b.Right())
	default:
		return false
	}
}

// rewriteReturns rewrites all return expressions in fn so `return foo();`
// becomes `foo(); return;` (peephole passes remove redundant trailing returns).
//
// Java: OptimizeReturns.rewriteReturns(Node)
func (o *OptimizeReturns) rewriteReturns(fn ast.Function) {
	body := fn.Body().Node()
	o.rewriteReturnsInSubtree(body, body)
}

func (o *OptimizeReturns) rewriteReturnsInSubtree(n rhino.Node, body rhino.Node) {
	parent := n.GetParent()
	if n.IsRoot() || n.IsBlock() || (!n.IsFunction() && NodeUtil_IsShallowStatementTree(parent)) {
		for c := n.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			o.rewriteReturnsInSubtree(c, body)
			c = next
		}
	}
	if ret, ok := ast.AsReturn(n); ok {
		if result, hasVal := ret.Value(); hasVal {
			keepValue := !o.isRemovableValue(result.Node())
			result.Detach()
			if keepValue {
				ast.NewExprStmt(result).Node().Srcref(result.Node()).InsertBefore(n)
			} else {
				NodeUtil_MarkFunctionsDeleted(result.Node(), o.compiler)
			}
			o.compiler.ReportChangeToEnclosingScope(body)
		}
	}
}

// isRemovableValue removes only values that don't reference properties (object
// literals) or names (functions), so the call graph does not need updating.
//
// Java: OptimizeReturns.isRemovableValue(Node)
func (o *OptimizeReturns) isRemovableValue(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_TEMPLATELIT, rhino.Token_ARRAYLIT:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if !child.IsEmpty() && !o.isRemovableValue(child) {
				return false
			}
		}
		return true
	case rhino.Token_REGEXP, rhino.Token_STRINGLIT, rhino.Token_NUMBER, rhino.Token_NULL, rhino.Token_TRUE, rhino.Token_FALSE, rhino.Token_TEMPLATELIT_STRING:
		return true
	case rhino.Token_TEMPLATELIT_SUB:
		return o.isRemovableValue(ast.MustTemplateLitSub(n).Expr().Node())
	case rhino.Token_CAST:
		return o.isRemovableValue(ast.MustCast(n).Expr().Node())
	case rhino.Token_NOT, rhino.Token_VOID, rhino.Token_NEG:
		return o.isRemovableValue(ast.MustUnary(n).Operand().Node())
	default:
		return false
	}
}
