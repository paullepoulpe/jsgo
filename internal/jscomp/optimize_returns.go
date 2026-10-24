// Port of com/google/javascript/jscomp/OptimizeReturns.java @ closure-compiler v20261015

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
	"reflect"

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
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
	logFile := o.compiler.CreateOrReopenIndexedLog(reflect.TypeFor[OptimizeReturns](), "decisions.log")
	defer logFile.Close()
	defer func() {
		o.decisionsLog = nil
		clear(o.fnNodes)
		o.fnNodes = o.fnNodes[:0]
	}()
	o.decisionsLog = logFile // avoid passing the log file through a bunch of methods
	// Find all function nodes whose callers ignore the return values.
	var toOptimize []*jutil.ArrayList[rhino.Node]

	// Find all the candidates before modifying the AST.
	for key, refs := range definitions.GetNameReferences() {
		if o.IsCandidate(key, refs.Slice()) {
			if o.decisionsLog.IsLogging() {
				o.decisionsLog.LogFormat("name %s\tremoving return value", key)
			}
			toOptimize = append(toOptimize, refs)
		}
	}

	for key, refs := range definitions.GetPropReferences() {
		if o.IsCandidate(key, refs.Slice()) {
			if o.decisionsLog.IsLogging() {
				o.decisionsLog.LogFormat("property %s\tremoving return value", key)
			}
			toOptimize = append(toOptimize, refs)
		}
	}

	// Now modify the AST
	for _, refs := range toOptimize {
		fns := OptimizeCalls_ReferenceMap_AppendFunctionNodes(o.fnNodes[:0], refs.Slice())
		o.fnNodes = fns
		for _, fn := range fns {
			o.RewriteReturns(fn)
		}
	}
}

// IsCandidate: this reference set is a candidate for return-value-removal
// if: - if the all call sites are known (not aliased, not exported) - if all
// call sites do not use the return value - if there is at least one known
// function definition - if there is at least one use NOTE: unknown
// definitions are allowed, as only known definitions will be removed.
//
// Java: OptimizeReturns.isCandidate(String, List<Node>)
func (o *OptimizeReturns) IsCandidate(name string, refs []rhino.Node) bool {
	if !OptimizeCalls_MayBeOptimizableName(o.compiler, name) {
		if o.decisionsLog.IsLogging() {
			o.decisionsLog.LogFormat("%s\tnot an optimizable name", name)
		}
		return false
	}

	seenCandidateDefiniton := false
	seenUse := false
	for _, n := range refs {
		// Assume indirect definitions references use the result
		if OptimizeCalls_ReferenceMap_IsCallTarget(n) || OptimizeCalls_ReferenceMap_IsOptChainCallTarget(n) {
			callNode := OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(n)
			if NodeUtil_IsExpressionResultUsed(callNode) {
				// At least one call site uses the return value, this
				// is not a candidate.
				if o.decisionsLog.IsLogging() { // avoid build location string when not logging
					o.decisionsLog.LogFormat("%s\treturn value used: %s", name, callNode.GetLocation())
				}
				return false
			}
			seenUse = true
		} else if o.IsCandidateDefinition(n) {
			// NOTE: While is is possible to optimize calls to functions for which we know
			// only some of the definition are candidates but to keep things simple, only
			// optimize if all of the definitions are known.
			seenCandidateDefiniton = true
		} else {
			// If this isn't an non-aliasing reference (typeof, instanceof, etc)
			// then there is nothing that can be done.
			if !OptimizeCalls_IsAllowedReference(n) {
				if o.decisionsLog.IsLogging() { // avoid build location string when not logging
					o.decisionsLog.LogFormat("%s\tdisallowed reference: %s", name, n.GetLocation())
				}
				return false
			}
		}
	}

	if !seenUse {
		if o.decisionsLog.IsLogging() {
			o.decisionsLog.LogFormat("%s\tno usage seen", name)
		}
		return false
	}
	if !seenCandidateDefiniton {
		if o.decisionsLog.IsLogging() {
			o.decisionsLog.LogFormat("%s\tno definition seen", name)
		}
		return false
	}
	return true
}

// Java: OptimizeReturns.isCandidateDefinition(Node)
func (o *OptimizeReturns) IsCandidateDefinition(n rhino.Node) bool {
	parent := n.GetParent()
	if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
		return true
	} else if OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n) {
		if OptimizeReturns_IsCandidateFunction(parent.GetLastChild()) {
			return true
		}
	} else if n.IsName() {
		if n.HasChildren() && OptimizeReturns_IsCandidateFunction(n.GetFirstChild()) {
			return true
		}
	} else if o.IsClassMemberDefinition(n) {
		return true
	}

	return false
}

// Java: OptimizeReturns.isClassMemberDefinition(Node)
func (o *OptimizeReturns) IsClassMemberDefinition(n rhino.Node) bool {
	return n.IsMemberFunctionDef() && n.GetParent().IsClassMembers()
}

// Java: OptimizeReturns.isCandidateFunction(Node)
func OptimizeReturns_IsCandidateFunction(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		// Named function expression can be recursive, this creates an alias of the name, meaning
		// it might be used in an unexpected way.
		return !NodeUtil_IsNamedFunctionExpression(n)
	case rhino.Token_COMMA, rhino.Token_CAST:
		return OptimizeReturns_IsCandidateFunction(n.GetLastChild())
	case rhino.Token_HOOK:
		return OptimizeReturns_IsCandidateFunction(n.GetSecondChild()) && OptimizeReturns_IsCandidateFunction(n.GetLastChild())
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		return OptimizeReturns_IsCandidateFunction(n.GetFirstChild()) && OptimizeReturns_IsCandidateFunction(n.GetLastChild())
	default:
		return false
	}
}

// RewriteReturns, for the supplied function node, rewrites all the return
// expressions so that: return foo(); becomes: foo(); return; Useless return
// will be removed later by the peephole optimization passes.
//
// Java: OptimizeReturns.rewriteReturns(Node)
func (o *OptimizeReturns) RewriteReturns(fnNode rhino.Node) {
	jutil.CheckState(fnNode.IsFunction())
	body := fnNode.GetLastChild()
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
	if n.IsReturn() && n.HasOneChild() {
		result := n.GetFirstChild()
		keepValue := !o.IsRemovableValue(result)
		result.Detach()
		if keepValue {
			rhino.IR_ExprResult(result).Srcref(result).InsertBefore(n)
		} else {
			NodeUtil_MarkFunctionsDeleted(result, o.compiler)
		}
		o.compiler.ReportChangeToEnclosingScope(body)
	}
}

// IsRemovableValue: just remove objects that don't reference properties
// (object literals) or names (functions) So we don't need to update the
// graph.
//
// Java: OptimizeReturns.isRemovableValue(Node)
func (o *OptimizeReturns) IsRemovableValue(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_TEMPLATELIT, rhino.Token_ARRAYLIT:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if (!child.IsEmpty()) && !o.IsRemovableValue(child) {
				return false
			}
		}
		return true
	case rhino.Token_REGEXP, rhino.Token_STRINGLIT, rhino.Token_NUMBER, rhino.Token_NULL, rhino.Token_TRUE, rhino.Token_FALSE, rhino.Token_TEMPLATELIT_STRING:
		return true
	case rhino.Token_TEMPLATELIT_SUB, rhino.Token_CAST, rhino.Token_NOT, rhino.Token_VOID, rhino.Token_NEG:
		return o.IsRemovableValue(n.GetFirstChild())
	default:
		return false
	}
}
