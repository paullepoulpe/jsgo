// Port of com/google/javascript/jscomp/OptimizeConstructors.java @ closure-compiler v20261021

/*
 * Copyright 2021 The Closure Compiler Authors.
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

// OptimizeConstructors optimizes class declarations by removing explicit
// constructor declarations if the implicit constructor is sufficient.
//
// # The constructor can be removed, if
//
// (0) the class in question defined using ES class syntax
//
// (2) the constructor is empty save for the call to super
//
// (3) either:
//
// - all constructor arguments are forward to the super constructor (via
// rest or arguments)
//
// - super constructor parameters match the number, order, and default
// values of constructor to be removed and the constructor doesn't look for
// additional parameters via references to 'arguments'
//
// When inspecting super class constructors, we don't need to worry about
// (1) reassignment of subclasses or (2) escapes of the superclass. We only
// need to make sure the superclass isn't reassigned. For that we can simply
// check for direct assignments and whether the name is defined in the
// externs as references. The two other reassignment cases that are
// hard/impossible to detect: "globalThis" properties and assignment through
// eval are only only supported with externs.
//
// Note that an alternative exists where the super class could be ignored, if
// all constructor references to the class were known (including direct
// calls, super calls, implicit constructor calls, etc) to validate the
// number of parameters passed. This would require backing off on various
// escapes and is believed it would be both more complicated and less
// effective.
type OptimizeConstructors struct {
	compiler    AbstractCompiler
	astAnalyzer *AstAnalyzer

	// All constructor definition nodes that are to be removed.
	removableConstructors []ast.MemberFunctionDef
}

// Java: OptimizeConstructors(AbstractCompiler)
func NewOptimizeConstructors(compiler AbstractCompiler) *OptimizeConstructors {
	compiler = jutil.CheckNotNull(compiler)
	return &OptimizeConstructors{
		compiler:    compiler,
		astAnalyzer: compiler.GetAstAnalyzer(),
	}
}

// Process is CompilerPass.process.
//
// Go-only: language-adapter — the two Java overloads process(Node, Node)
// (CompilerPass) and process(Node, Node, ReferenceMap)
// (CallGraphCompilerPass) cannot share a Go name; the second one is
// ProcessWithReferenceMap, and AsCallGraphCompilerPass adapts it.
//
// Java: OptimizeConstructors.process(Node, Node)
func (o *OptimizeConstructors) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(o.compiler.GetLifeCycleStage() == AbstractCompiler_LifeCycleStage_NORMALIZED)

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
func (o *OptimizeConstructors) AsCallGraphCompilerPass() OptimizeCalls_CallGraphCompilerPass {
	return OptimizeCalls_CallGraphCompilerPassFunc(o.ProcessWithReferenceMap)
}

// Java: OptimizeConstructors.process(Node, Node, ReferenceMap)
func (o *OptimizeConstructors) ProcessWithReferenceMap(externs rhino.Node, root rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
	defer func() {
		clear(o.removableConstructors)
		o.removableConstructors = o.removableConstructors[:0]
	}()

	for key, value := range refMap.GetNameReferences() {
		o.addConstructorsToBeRemoved(key, value)
	}

	// NOTE: Normally, the ReferenceMap must be kept in a consistent state (removing references
	// as they are removed from the AST) so the next pass can reuse the ReferenceMap. However,
	// currently this pass runs by itself so that work is avoided.
	for _, ctor := range o.removableConstructors {
		o.removeConstructorMethod(ctor)
	}
}

// Java: OptimizeConstructors.addConstructorsToBeRemoved(String, ArrayList<Node>)
func (o *OptimizeConstructors) addConstructorsToBeRemoved(name string, refs []rhino.Node) {
	if !OptimizeCalls_MayBeOptimizableName(o.compiler, name) {
		return
	}

	var candidateClassDefinition rhino.Node
	var candidateConstructor ast.MemberFunctionDef
	hasCandidateConstructor := false

	// Lazily init as most symbols aren't classes, and most classes don't have subclasses.
	var subclassConstructors []ast.MemberFunctionDef

	for _, n := range refs {
		definition := optimizeConstructors_getClassDefinitionOrFunction(n)
		if !definition.IsNil() {
			if !candidateClassDefinition.IsNil() {
				// As a simplification only allow one definition.
				return
			}

			// Be lazy about constructor analysis as it is expected that most classes
			// don't have subclasses.
			candidateClassDefinition = definition

			if cls, ok := ast.AsClass(candidateClassDefinition); ok {
				// Handle no super class:
				//   - no extends clause
				//   - constructor body is empty
				//   - constructor parameters are side-effect free (no destructuring or side-effect defaults)
				if _, hasExtends := cls.Extends(); !hasExtends {
					if ctor, hasCtor := cls.Constructor(); hasCtor {
						if !ctor.Function().Body().Node().HasChildren() && o.hasRemovableParameterList(ctor) {
							candidateConstructor = ctor
							hasCandidateConstructor = true
						}
					}
				}
			}
		} else if ast.RoleIn(n, n.GetParent()) == ast.RoleClassExtends {
			// If this is a reference to the class in an extend clause, then this is a subclass
			// definition whose constructor should be considered for removal.
			subclass := ast.MustClass(n.GetParent())
			if subCtor, hasCtor := subclass.Constructor(); hasCtor &&
				optimizeConstructors_constructorHasRemovableDefinition(subCtor) {
				subclassConstructors = append(subclassConstructors, subCtor)
			}
		} else if optimizeConstructors_isAssigningReference(n) {
			// Any assignment that isn't a class definition covered above is an invalidating assignment.
			return
		}
	}

	if candidateClassDefinition.IsNil() {
		return
	}

	if hasCandidateConstructor {
		o.removableConstructors = append(o.removableConstructors, candidateConstructor)
	}

	// Nothing invalidated the superclass definition; check whether candidate
	// subclasses are equivalent to the superclass.
	if len(subclassConstructors) > 0 {
		summary := optimizeConstructors_buildSummary(candidateClassDefinition)
		for _, subCtor := range subclassConstructors {
			if summary.isEquivalentConstructorDefinition(subCtor) {
				o.removableConstructors = append(o.removableConstructors, subCtor)
			}
		}
	}
}

// Java: OptimizeConstructors.isAssigningReference(Node)
func optimizeConstructors_isAssigningReference(n rhino.Node) bool {
	parent := n.GetParent()
	switch parent.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		return n.HasChildren() // value assigned
	case rhino.Token_STRING_KEY:
		return parent.GetParent().IsObjectPattern()
	case rhino.Token_COMPUTED_PROP:
		return parent.GetLastChild() == n && parent.GetParent().IsObjectPattern()
	case rhino.Token_ARRAY_PATTERN,
		rhino.Token_DEFAULT_VALUE, // object or array or function parameter
		rhino.Token_PARAM_LIST,
		rhino.Token_OBJECT_REST,
		rhino.Token_ITER_REST,
		rhino.Token_INC,
		rhino.Token_DEC:
		return true
	case rhino.Token_FUNCTION, rhino.Token_CLASS, rhino.Token_CATCH:
		return parent.GetFirstChild() == n
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
		// If n is the target of the assignment (the first child/LHS of the assignment node),
		// then it is an assigning reference.
		return parent.GetFirstChild() == n
	default:
		return false
	}
}

// Java: OptimizeConstructors.getClassDefinitionOrFunction(Node)
func optimizeConstructors_getClassDefinitionOrFunction(n rhino.Node) rhino.Node {
	parent := n.GetParent()

	var expr rhino.Node
	if OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n) {
		expr = ast.MustAssign(parent).Value().Node()
	} else if name, ok := ast.AsName(n); ok && n.HasChildren() {
		init, _ := name.Init()
		expr = init.Node()
	} else if (parent.IsFunction() || parent.IsClass()) && n.IsFirstChildOf(parent) {
		expr = parent
	} else {
		return rhino.NilNode
	}

	for {
		if cast, ok := ast.AsCast(expr); ok {
			expr = cast.Expr().Node()
		} else if comma, ok := ast.AsBinaryOf(expr, rhino.Token_COMMA); ok {
			expr = comma.Right().Node()
		} else {
			break
		}
	}

	switch expr.GetToken() {
	case rhino.Token_FUNCTION:
		// TODO(b/176208718): ideally this only returns true for normal functions, but it is
		// harmless to include other function types.
		return expr
	case rhino.Token_CLASS:
		if _, hasCtor := ast.MustClass(expr).Constructor(); hasCtor {
			return expr
		}
	}
	return rhino.NilNode
}

// optimizeConstructorsSummary summarizes a superclass constructor definition
// to avoid repeated analysis when checking subclass constructors.
type optimizeConstructorsSummary struct {
	formalParameterCount int
	isVarArgs            bool
}

// Java: OptimizeConstructors.ClassConstructorSummary.build(Node)
func optimizeConstructors_buildSummary(classDefinition rhino.Node) optimizeConstructorsSummary {
	var fn ast.Function
	if cls, ok := ast.AsClass(classDefinition); ok {
		ctor, hasCtor := cls.Constructor()
		jutil.CheckState(hasCtor)
		fn = ctor.Function()
	} else {
		fn = ast.MustFunction(classDefinition)
	}

	hasVarArgs := NodeUtil_DoesFunctionReferenceOwnArgumentsObject(fn.Node()) ||
		optimizeConstructors_functionHasRest(fn)
	return optimizeConstructorsSummary{
		formalParameterCount: fn.Params().NumParams(),
		isVarArgs:            hasVarArgs,
	}
}

// Java: OptimizeConstructors.ClassConstructorSummary.isEquivalentConstructorDefinition(Node)
func (s optimizeConstructorsSummary) isEquivalentConstructorDefinition(ctor ast.MemberFunctionDef) bool {
	fn := ctor.Function()
	if optimizeConstructors_functionHasRest(fn) {
		return true
	}
	if s.isVarArgs {
		// The count of parameters may matter, and the function may not pass them all on.
		return false
	}
	return s.formalParameterCount == fn.Params().NumParams()
}

// Java: OptimizeConstructors.functionHasRest(Node)
func optimizeConstructors_functionHasRest(fn ast.Function) bool {
	lastParam, ok := fn.Params().LastParam()
	return ok && lastParam.Node().IsRest()
}

// Java: OptimizeConstructors.constructorHasRemovableDefinition(Node)
func optimizeConstructors_constructorHasRemovableDefinition(ctor ast.MemberFunctionDef) bool {
	fn := ctor.Function()
	superCall, ok := optimizeConstructors_getOnlySuperCall(fn)
	if !ok {
		return false
	}

	params := fn.Params()
	if params.NumParams() != superCall.NumArgs() {
		return false
	}

	// TODO(johnlenz): broaden the recognized patterns: default parameters
	param, hasParam := params.FirstParam()
	arg, _ := superCall.FirstArg()
	for hasParam {
		if rest, isRest := param.AsIterRest(); isRest {
			spread, isSpread := arg.AsIterSpread()
			if !isSpread || !rest.Target().Node().MatchesNameNode(spread.Expr().Node()) {
				return false
			}
		} else if !param.Node().MatchesNameNode(arg.Node()) {
			return false
		}
		param, hasParam = param.Next()
		arg, _ = arg.Next()
	}

	return true
}

// Java: OptimizeConstructors.hasRemovableParameterList(Node)
func (o *OptimizeConstructors) hasRemovableParameterList(ctor ast.MemberFunctionDef) bool {
	for param := range ctor.Function().Params().Params() {
		if _, ok := param.AsName(); ok {
			continue
		}
		if rest, ok := param.AsIterRest(); ok && rest.Target().Node().IsName() {
			continue
		}
		if dv, ok := param.AsDefaultValue(); ok &&
			dv.Target().Node().IsName() &&
			!o.astAnalyzer.MayHaveSideEffects(dv.Value().Node()) {
			continue
		}
		// not a simple parameter (destructuring can throw, etc)
		return false
	}
	return true
}

// Java: OptimizeConstructors.getOnlySuperCall(Node)
func optimizeConstructors_getOnlySuperCall(fn ast.Function) (ast.Call, bool) {
	block, ok := fn.Body().AsBlock()
	if !ok {
		return ast.Call{}, false
	}
	stmt, ok := block.OnlyStatement()
	if !ok {
		return ast.Call{}, false
	}
	exprStmt, ok := stmt.AsExprStmt()
	if !ok {
		return ast.Call{}, false
	}
	call, ok := exprStmt.Expr().AsCall()
	if !ok || !call.Callee().Node().IsSuper() {
		return ast.Call{}, false
	}
	return call, true
}

// Java: OptimizeConstructors.removeConstructorMethod(Node)
func (o *OptimizeConstructors) removeConstructorMethod(ctor ast.MemberFunctionDef) {
	o.compiler.ReportFunctionDeleted(ctor.Function().Node())
	o.compiler.ReportChangeToEnclosingScope(ctor.Node())
	ctor.Detach()
}
