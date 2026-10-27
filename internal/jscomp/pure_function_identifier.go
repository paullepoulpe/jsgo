// Port of com/google/javascript/jscomp/PureFunctionIdentifier.java @ closure-compiler v20261021

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
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: callable-leaf, r-value, try-catch, assignment, and invocation
// shapes are read through package ast while keeping PureFunctionIdentifier's
// rhino.Node pass API intact.

// PureFunctionIdentifier is a compiler pass that computes function purity
// and annotates invocation nodes with those purities.
//
// A function is pure if it has no outside visible side effects, and the
// result of the computation does not depend on external factors that are
// beyond the control of the application; repeated calls to the function
// should return the same value as long as global state hasn't changed.
//
// `Date.now` is an example of a function that has no side effects but is
// not pure.
//
// Functions are not tracked individually but rather in aggregate by their
// name. This is because it's impossible to determine exactly which function
// named "foo" is being called at a particular site. Therefore, if any
// function "foo" has a particular side-effect, all invocations "foo" are
// assumed to trigger it.
//
// This pass could be greatly improved by proper tracking of locals within
// function bodies. Every instance of the call to
// NodeUtil#evaluatesToLocalValue(Node) and NodeUtil#allArgsUnescapedLocal(Node)
// do not actually take into account local variables. They only assume
// literals, primitives, and operations on primitives are local.
type PureFunctionIdentifier struct {
	compiler                 AbstractCompiler
	astAnalyzer              *AstAnalyzer
	validateArtificialPurity bool

	// Map of function names to the summary of the functions with that name.
	//
	// Variable names are recorded as-is. Property names are prefixed with
	// PROP_NAME_PREFIX to differentiate them from variable names.
	//
	// See AmbiguatedFunctionSummary.
	summariesByName     map[string]*PureFunctionIdentifier_AmbiguatedFunctionSummary
	propSummariesByName map[string]*PureFunctionIdentifier_AmbiguatedFunctionSummary

	// Mapping from function node to summaries for all names associated with
	// that node.
	//
	// This is a multimap because you can construct situations in which a
	// function node has multiple names, and therefore multiple associated
	// summaries. For example:
	//
	//	// Not enough type information to collapse/disambiguate properties on "staticMethod".
	//	SomeClass.staticMethod = function anotherName() {};
	//	OtherClass.staticMethod = function() { global++; }
	//
	// In this situation we want to keep the side effects for "staticMethod"
	// which are "global" separate from "anotherName". Hence the function node
	// should point to the AmbiguatedFunctionSummary for both "staticMethod"
	// and "anotherName".
	//
	// We could instead store a map of FUNCTION nodes to names, and then join
	// that with the name of names to infos. However, since names are 1:1 with
	// infos, it's more effecient to "pre-join" in this way.
	//
	// Go-only: language-adapter — an ArrayListMultimap keyed by
	// identity-hashed Nodes that is only looked up, so insertion order.
	summariesForAllNamesOfFunctionByNode map[rhino.Node][]*PureFunctionIdentifier_AmbiguatedFunctionSummary

	// Set of FUNCTION nodes which are artificially marked as pure, via a
	// `@nosideeffects` JSDoc tag in source code, which triggers the compiler
	// to ignore any actual side effects.
	//
	// We store this here instead of in the AmbiguatedFunctionSummary because
	// 1) this is only for debugging, and we don't want to accidentally depend
	// on it for side effect analysis and 2) we expect the set of functions
	// with artificial purity annotations to be small and not scale linearly
	// with the size of the application.
	artificiallyPureLiteralsForDebugging *jutil.LinkedHashSet[rhino.Node]

	// List of all function call sites. Storing them here during the function
	// analysis traversal prevents us from doing a second traversal to annotate
	// them with side-effects. We can just iterate the list.
	allFunctionCalls []rhino.Node

	// A graph linking the summary of a function callee to the summaries of
	// its callers.
	//
	// Each node represents an aggregate summary of every function with a
	// particular name. The edge values indicate the details of the invocation
	// necessary to propagate function impurity from callee to caller.
	//
	// Once all the invocation edges are in place, this graph will be
	// traversed to transitively propagate side-effect information across it's
	// entire structure. The resultant side-effects can then be attached to
	// invocation sites.
	reverseCallGraph *graph.LinkedDirectedGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]

	// A summary for a function for which no definition was found.
	//
	// We assume it has all possible side-effects. It's useful for references
	// like function parameters, or inner functions.
	unknownFunctionSummary      *PureFunctionIdentifier_AmbiguatedFunctionSummary
	unknownFunctionSummarySlice []*PureFunctionIdentifier_AmbiguatedFunctionSummary

	summarySlab           jutil.Slab[PureFunctionIdentifier_AmbiguatedFunctionSummary]
	propSlab              jutil.Slab[PureFunctionIdentifier_SideEffectPropagation]
	scratchCallableLeaves []rhino.Node

	assumeGettersArePure bool

	hasProcessed bool
}

// A prefix to differentiate property names from variable names.
// TODO(nickreid): This pass could be made more efficient if props and
// variables were maintained in separate datastructures. We wouldn't allocate
// a bunch of extra strings.
const PureFunctionIdentifier_PROP_NAME_PREFIX = "."

// PureFunctionIdentifier_IMPLICIT_PURE_FN is a function node representing a
// function implicit in the AST that is known to be pure
//
// For example: `class C {}` would get this node.
var PureFunctionIdentifier_IMPLICIT_PURE_FN = rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_Block())

// Enable this error if you are debugging why a `@nosideeffects` annotation
// seems to be ignored.
var PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION = DiagnosticType_Disabled(
	"JSC_UNUSED_ARTIFICIAL_PURE_ANNOTATION",
	"Artificial @nosideeffects annotation cannot be enforced: found ambiguous definitions of"+
		" {0}.\n\n"+
		"Side-effectful definitions:\n  {1}",
)

// Java: PureFunctionIdentifier(AbstractCompiler, boolean, boolean)
func NewPureFunctionIdentifier(
	compiler AbstractCompiler, assumeGettersArePure bool, validateArtificialPurity bool) *PureFunctionIdentifier {
	p := &PureFunctionIdentifier{}
	p.artificiallyPureLiteralsForDebugging = jutil.NewLinkedHashSet[rhino.Node]()
	p.reverseCallGraph = graph.LinkedDirectedGraph_CreateWithoutAnnotations[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]()
	p.unknownFunctionSummary =
		p.createSummaryInGraph("<unknown>").
			SetMutatesGlobalStateAndAllOtherFlags(rhino.NilNode)
	p.unknownFunctionSummarySlice = p.unknownFunctionSummary.selfSlice[:]
	p.hasProcessed = false
	p.compiler = jutil.CheckNotNull(compiler)
	p.assumeGettersArePure = assumeGettersArePure
	p.astAnalyzer = compiler.GetAstAnalyzer()
	p.validateArtificialPurity = validateArtificialPurity
	return p
}

func (p *PureFunctionIdentifier) createSummaryInGraph(name string) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	s := p.summarySlab.New()
	s.name = name
	n := p.reverseCallGraph.AllocNode(s)
	p.reverseCallGraph.Nodes.PutNew(s, n)
	s.graphNode = n
	s.selfSlice[0] = s
	return s
}

func (p *PureFunctionIdentifier) putSummaryForFunctionNode(fn rhino.Node, summary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	existing := p.summariesForAllNamesOfFunctionByNode[fn]
	if len(existing) == 0 {
		p.summariesForAllNamesOfFunctionByNode[fn] = summary.selfSlice[:1:1]
	} else {
		p.summariesForAllNamesOfFunctionByNode[fn] = append(existing, summary)
	}
}

func (p *PureFunctionIdentifier) getSummaryOrDefault(name string) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	if s := p.summariesByName[name]; s != nil {
		return s
	}
	return p.unknownFunctionSummary
}

func (p *PureFunctionIdentifier) getSummaryForReference(nameRef rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	if name, ok := ast.AsName(nameRef); ok {
		return p.summariesByName[name.Ident()]
	}
	if gp, ok := ast.AsNormalOrOptChainGetProp(nameRef); ok {
		return p.propSummariesByName[gp.Property()]
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected name reference: " + nameRef.String()))
}

func (p *PureFunctionIdentifier) getSummaryForReferenceOrDefault(nameRef rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	if s := p.getSummaryForReference(nameRef); s != nil {
		return s
	}
	return p.unknownFunctionSummary
}

// allocPropagationForInvocation is SideEffectPropagation.forInvocation,
// allocated from a slab.
func (p *PureFunctionIdentifier) allocPropagationForInvocation(
	invocation rhino.Node,
	callScope *Scope,
	callerEntry *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry,
	argumentsMatchParameters bool,
	propagateThrows bool) *PureFunctionIdentifier_SideEffectPropagation {
	jutil.CheckArgumentMsg1(NodeUtil_IsInvocation(invocation), "%s", invocation)
	prop := p.propSlab.New()
	// `this` in an arrow function is the `this` of the enclosing function, which is not
	// tracked.
	prop.calleeThisEqualsCallerThis = PureFunctionIdentifier_SideEffectPropagation_CalleeAndCallerShareThis(invocation) && !callerEntry.root.IsArrowFunction()
	prop.propagateThrows = propagateThrows
	prop.invocation = invocation
	prop.callScope = callScope
	prop.callerEntry = callerEntry
	prop.argumentsMatchParameters = argumentsMatchParameters
	return prop
}

// Java: PureFunctionIdentifier.process(Node, Node, ReferenceMap)
func (p *PureFunctionIdentifier) Process(externs rhino.Node, root rhino.Node, references *OptimizeCalls_ReferenceMap) {
	jutil.CheckState(p.compiler.GetLifeCycleStage().IsNormalized())
	jutil.CheckStateMsg(
		!p.hasProcessed, "PureFunctionIdentifier::process may only be called once per instance.")
	p.hasProcessed = true

	p.PopulateDatastructuresForAnalysisTraversal(references)

	NodeTraversal_Traverse(p.compiler, externs, NewPureFunctionIdentifier_ExternFunctionAnnotationAnalyzer(p))
	sc := NewSyntacticScopeCreator(p.compiler)
	b := NodeTraversal_Builder{
		compiler:     p.compiler,
		callback:     NewPureFunctionIdentifier_FunctionBodyAnalyzer(p),
		scopeCreator: sc,
	}
	b.Traverse(root)

	p.PropagateSideEffects()

	p.ValidateArtificialPurity()

	p.MarkPureFunctionCalls()
}

// ValidateArtificialPurity checks whether any function literals marked as
// artificially pure still have an overall AmbiguatedFunctionSummary with
// some side-effects, for debugging purposes.
//
// This can happen when there are multiple definitions of the same function
// name, and one is marked as artificially pure but another is
// side-effect-full.
//
// Java: PureFunctionIdentifier.validateArtificialPurity()
func (p *PureFunctionIdentifier) ValidateArtificialPurity() {
	if !p.validateArtificialPurity {
		// Avoid unnecessary computation.
		return
	}
	for function := range p.artificiallyPureLiteralsForDebugging.All() {
		for _, summary := range p.summariesForAllNamesOfFunctionByNode[function] {
			if summary.HasNoFlagsSet() {
				// no flags set == the function is successfully treated as pure.
				continue
			}
			if summary.name == ".constructor" {
				// .constructor is a special case - almost all the time, constructors get invoked by
				// `new Foo();` instead of `new something.constructor();` are a special case - most of the
				// time they are invoked not by the
				// this.constructor property & JSCompiler doesn't support this.constructor + property
				// renaming well. So don't report this error, as it's not particularly helpful.
				continue
			}
			var reasons []rhino.Node
			if summary.impureFunctionReasonsForDebugging != nil {
				reasons = summary.impureFunctionReasonsForDebugging.ToSlice()
			} else {
				// stream() of a null set.
				panic(jutil.NewNullPointerException())
			}
			slices.SortStableFunc(reasons, PureFunctionIdentifier_SOURCE_LOCATION_COMPARATOR)
			formatted := make([]string, 0, len(reasons))
			for _, reason := range reasons {
				formatted = append(formatted, PureFunctionIdentifier_FormatSourceLocation(reason))
			}
			sideEffectfulDefinitions := strings.Join(formatted, "\n  ")
			p.compiler.Report(
				JSError_MakeWithNode(
					function,
					PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION,
					summary.name,
					sideEffectfulDefinitions))
		}
	}
}

// PureFunctionIdentifier_SOURCE_LOCATION_COMPARATOR is
// Comparator.comparing(Node::getSourceFileName).thenComparingInt(Node::getLineno).thenComparingInt(Node::getCharno).
func PureFunctionIdentifier_SOURCE_LOCATION_COMPARATOR(a, b rhino.Node) int {
	aName, bName := a.GetSourceFileName(), b.GetSourceFileName()
	if aName == nil || bName == nil {
		// String.compareTo on a null source file name.
		panic(jutil.NewNullPointerException())
	}
	if c := jutil.StringCompareTo(*aName, *bName); c != 0 {
		return c
	}
	if c := a.GetLineno() - b.GetLineno(); c != 0 {
		if c < 0 {
			return -1
		}
		return 1
	}
	if a.GetCharno() < b.GetCharno() {
		return -1
	} else if a.GetCharno() > b.GetCharno() {
		return 1
	}
	return 0
}

// Java: PureFunctionIdentifier.formatSourceLocation(Node)
func PureFunctionIdentifier_FormatSourceLocation(node rhino.Node) string {
	name := "null"
	if n := node.GetSourceFileName(); n != nil {
		name = *n
	}
	return name + ":" + strconv.Itoa(node.GetLineno()) + ":" + strconv.Itoa(node.GetCharno())
}

// PureFunctionIdentifier_CollectCallableLeaves traverses an expr to collect
// nodes representing potential callables that it may resolve to well known
// callables.
//
// See collectCallableLeavesInternal. Returns the disovered callables, or
// null if an unexpected possible value was found.
//
// Java: PureFunctionIdentifier.collectCallableLeaves(Node, AstAnalyzer)
func PureFunctionIdentifier_CollectCallableLeaves(
	expr rhino.Node, astAnalyzer *AstAnalyzer) []rhino.Node {
	var callables []rhino.Node
	if pureFunctionIdentifier_collectCallableLeavesSlice(expr, &callables, astAnalyzer) {
		return callables
	}
	return nil
}

func pureFunctionIdentifier_invocationCallee(invocation rhino.Node) ast.Expr {
	if callOrNew, ok := ast.AsCallOrNew(invocation); ok {
		return callOrNew.Callee()
	}
	return ast.MustTaggedTemplateLit(invocation).Tag()
}

func pureFunctionIdentifier_collectCallableLeavesSlice(
	expr rhino.Node, results *[]rhino.Node, astAnalyzer *AstAnalyzer) bool {
	switch expr.GetToken() {
	case rhino.Token_FUNCTION, rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_NAME:
		*results = append(*results, expr)
		return true

	case rhino.Token_SUPER:
		clazz := ast.MustClass(rhino.CheckNotNull(NodeUtil_GetEnclosingClass(expr)))
		return pureFunctionIdentifier_collectCallableLeavesSlice(clazz.ExtendsNode(), results, astAnalyzer)

	case rhino.Token_CLASS:
		cls := ast.MustClass(expr)
		for member := range cls.Body().Members() {
			if mf, ok := member.AsMemberFieldDef(); ok && !mf.IsStatic() {
				if init, hasInit := mf.Value(); hasInit && astAnalyzer.MayHaveSideEffects(init.Node()) {
					return false
				}
			} else if cf, ok := member.AsComputedFieldDef(); ok && !cf.IsStatic() {
				if init, hasInit := cf.Value(); hasInit && astAnalyzer.MayHaveSideEffects(init.Node()) {
					return false
				}
			}
		}

		ctorDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(expr)
		if !ctorDef.IsNil() {
			return pureFunctionIdentifier_collectCallableLeavesSlice(
				ast.MustMemberFunctionDef(ctorDef).Function().Node(), results, astAnalyzer)
		} else if ext, hasExt := cls.Extends(); !hasExt {
			*results = append(*results, PureFunctionIdentifier_IMPLICIT_PURE_FN)
			return true
		} else {
			return pureFunctionIdentifier_collectCallableLeavesSlice(ext.Node(), results, astAnalyzer)
		}

	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		bin := ast.MustBinary(expr)
		return pureFunctionIdentifier_collectCallableLeavesSlice(bin.Left().Node(), results, astAnalyzer) &&
			pureFunctionIdentifier_collectCallableLeavesSlice(bin.Right().Node(), results, astAnalyzer)

	case rhino.Token_COMMA:
		return pureFunctionIdentifier_collectCallableLeavesSlice(ast.MustBinary(expr).Right().Node(), results, astAnalyzer)

	case rhino.Token_ASSIGN:
		return pureFunctionIdentifier_collectCallableLeavesSlice(ast.MustAssign(expr).Value().Node(), results, astAnalyzer)

	case rhino.Token_HOOK:
		hook := ast.MustHook(expr)
		return pureFunctionIdentifier_collectCallableLeavesSlice(hook.Then().Node(), results, astAnalyzer) &&
			pureFunctionIdentifier_collectCallableLeavesSlice(hook.Else().Node(), results, astAnalyzer)

	case rhino.Token_NEW_TARGET, rhino.Token_THIS:
		return false
	default:
		return false
	}
}

// PureFunctionIdentifier_IsDefinitelyRValue returns true only if rvalue is
// defintely a reference reading a value.
//
// For the most part it's sufficient to cover cases where a nominal function
// reference might reasonably be expected, since those are the values that
// matter to analysis.
//
// It's very important that this never returns true for an L-value,
// including when new syntax is added to the language. That would cause some
// impure functions to be considered pure. Therefore, this method is a very
// explict allowlist. Anything that's unrecognized is considered not an
// R-value. This is insurance against new syntax.
//
// New cases can be added as needed to increase the accuracy of the analysis.
// They just have to be verified as always R-values.
//
// Java: PureFunctionIdentifier.isDefinitelyRValue(Node)
func PureFunctionIdentifier_IsDefinitelyRValue(rvalue rhino.Node) bool {
	parent := rvalue.GetParent()

	switch parent.GetToken() {
	case rhino.Token_AND,
		rhino.Token_ARRAYLIT,
		rhino.Token_CALL,
		rhino.Token_COALESCE,
		rhino.Token_COMMA,
		rhino.Token_EQ,
		rhino.Token_GETELEM,
		rhino.Token_GETPROP,
		rhino.Token_HOOK,
		rhino.Token_INSTANCEOF,
		rhino.Token_NEW,
		rhino.Token_NOT,
		rhino.Token_NAME,
		rhino.Token_OPTCHAIN_CALL,
		rhino.Token_OPTCHAIN_GETELEM,
		rhino.Token_OPTCHAIN_GETPROP,
		rhino.Token_OR,
		rhino.Token_RETURN,
		rhino.Token_SHEQ,
		rhino.Token_TAGGED_TEMPLATELIT,
		rhino.Token_TYPEOF,
		rhino.Token_YIELD:
		return true
	case rhino.Token_CASE:
		return ast.MustCase(parent).Value().Is(rvalue)
	case rhino.Token_IF:
		return ast.MustIf(parent).Cond().Is(rvalue)
	case rhino.Token_SWITCH:
		return ast.MustSwitch(parent).Discriminant().Is(rvalue)
	case rhino.Token_WHILE:
		return ast.MustWhile(parent).Cond().Is(rvalue)
	case rhino.Token_EXPR_RESULT:
		// Extern declarations are sometimes stubs. These must be considered L-values with no
		// associated R-values.
		return !rvalue.IsFromExterns()
	case rhino.Token_ASSIGN:
		return ast.MustAssign(parent).Value().Is(rvalue)
	case rhino.Token_CLASS: // `extends` clause.
		ext, hasExt := ast.MustClass(parent).Extends()
		return hasExt && ext.Is(rvalue)
	case rhino.Token_STRING_KEY: // Assignment to an object literal property. Excludes object destructuring.
		_, isObjLit := ast.AsObjectLit(parent.GetParent())
		return isObjLit
	default:
		// Anything not explicitly listed may not be an R-value. We only worry about the likely
		// cases for nominal function values since those are what interest us and its safe to miss
		// some R-values. It's more important that we correctly identify L-values.
		return false
	}
}

// Java: PureFunctionIdentifier.getGoogCacheCallableExpression(Cache)
func (p *PureFunctionIdentifier) GetGoogCacheCallableExpression(cacheCall *CodingConvention_Cache) []rhino.Node {
	jutil.CheckNotNull(cacheCall)

	valueFnCallables := PureFunctionIdentifier_CollectCallableLeaves(cacheCall.valueFn, p.astAnalyzer)
	if valueFnCallables == nil {
		// ImmutableList.Builder.addAll(null)
		panic(jutil.NewNullPointerException())
	}
	builder := append([]rhino.Node(nil), valueFnCallables...)
	if !cacheCall.keyFn.IsNil() {
		keyFnCallables := PureFunctionIdentifier_CollectCallableLeaves(cacheCall.keyFn, p.astAnalyzer)
		if keyFnCallables == nil {
			// ImmutableList.Builder.addAll(null)
			panic(jutil.NewNullPointerException())
		}
		builder = append(builder, keyFnCallables...)
	}
	return builder
}

// Java: PureFunctionIdentifier.getSummariesForCallee(Node)
func (p *PureFunctionIdentifier) GetSummariesForCallee(invocation rhino.Node) []*PureFunctionIdentifier_AmbiguatedFunctionSummary {
	jutil.CheckArgumentMsg1(NodeUtil_IsInvocation(invocation), "%s", invocation)

	cacheCall := p.compiler.GetCodingConvention().DescribeCachingCall(invocation)

	if cacheCall == nil {
		var singleCallee rhino.Node
		callee := pureFunctionIdentifier_invocationCallee(invocation)
		if PureFunctionIdentifier_IsInvocationViaCallOrApply(invocation) {
			singleCallee = callee.MustNormalOrOptChainGet().Object().Node()
		} else {
			first := callee.Node()
			switch first.GetToken() {
			case rhino.Token_FUNCTION, rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_NAME:
				singleCallee = first
			}
		}
		if !singleCallee.IsNil() {
			if singleCallee.IsFunction() {
				summariesForFunction := p.summariesForAllNamesOfFunctionByNode[singleCallee]
				jutil.CheckStateMsg1(len(summariesForFunction) != 0, "Function missed during analysis: %s", singleCallee)
				return summariesForFunction
			}
			if s := p.getSummaryForReference(singleCallee); s != nil {
				return s.selfSlice[:]
			}
			return p.unknownFunctionSummarySlice
		}
	}

	var callees []rhino.Node
	if cacheCall != nil {
		callees = p.GetGoogCacheCallableExpression(cacheCall)
	} else {
		p.scratchCallableLeaves = p.scratchCallableLeaves[:0]
		if pureFunctionIdentifier_collectCallableLeavesSlice(pureFunctionIdentifier_invocationCallee(invocation).Node(), &p.scratchCallableLeaves, p.astAnalyzer) {
			callees = p.scratchCallableLeaves
		}
	}

	if callees == nil {
		return p.unknownFunctionSummarySlice
	}

	jutil.CheckStateMsg(len(callees) != 0, "Unexpected empty callees for valid result")
	results := make([]*PureFunctionIdentifier_AmbiguatedFunctionSummary, 0, len(callees))
	for _, callee := range callees {
		if callee.IsFunction() {
			jutil.CheckStateMsg1(callee.IsFunction(), "%s", callee)

			summariesForFunction :=
				p.summariesForAllNamesOfFunctionByNode[callee]
			jutil.CheckStateMsg1(len(summariesForFunction) != 0, "Function missed during analysis: %s", callee)

			results = append(results, summariesForFunction...)
		} else {
			results = append(results, p.getSummaryForReferenceOrDefault(callee))
		}
	}
	return results
}

func (p *PureFunctionIdentifier) ensureReverseCallGraphCapacity(totalNodes int) {
	if totalNodes <= 8 {
		return
	}
	p.reverseCallGraph.Nodes.EnsureCapacity(totalNodes)
}

// PopulateDatastructuresForAnalysisTraversal fills all of the auxiliary
// data-structures used by this pass based on the results in referenceMap.
//
// This is the first step of analysis. These structures will be used by a
// traversal that analyzes the bodies of located functions for side-effects.
// That traversal is separate because it needs access to scopes and also
// depends on global knowledge of functions.
//
// Java: PureFunctionIdentifier.populateDatastructuresForAnalysisTraversal(ReferenceMap)
func (p *PureFunctionIdentifier) PopulateDatastructuresForAnalysisTraversal(referenceMap *OptimizeCalls_ReferenceMap) {
	totalRefs := referenceMap.names.Size() + referenceMap.props.Size()
	p.summariesByName = make(map[string]*PureFunctionIdentifier_AmbiguatedFunctionSummary, totalRefs)
	p.propSummariesByName = make(map[string]*PureFunctionIdentifier_AmbiguatedFunctionSummary, referenceMap.props.Size())
	p.summariesForAllNamesOfFunctionByNode = make(map[rhino.Node][]*PureFunctionIdentifier_AmbiguatedFunctionSummary, max(totalRefs/2, 16))
	p.ensureReverseCallGraphCapacity(totalRefs + 2)

	hasDotNameCollision := false
	for key, refs := range referenceMap.GetNameReferences() {
		if len(refs) == 0 {
			continue
		}
		jutil.CheckState(key != "" && key != PureFunctionIdentifier_PROP_NAME_PREFIX)
		p.summariesByName[key] = p.createSummaryInGraph(key)
	}
	for key, refs := range referenceMap.GetPropReferences() {
		if len(refs) == 0 {
			continue
		}
		jutil.CheckState(key != "")
		propKey := p.allocPropKey(key)
		if existing := p.summariesByName[propKey]; existing != nil {
			p.propSummariesByName[key] = existing
			hasDotNameCollision = true
		} else {
			s := p.createSummaryInGraph(propKey)
			p.summariesByName[propKey] = s
			p.propSummariesByName[key] = s
		}
	}

	if !hasDotNameCollision {
		for key, refs := range referenceMap.GetNameReferences() {
			if len(refs) == 0 {
				continue
			}
			p.populateFunctionDefinitionsWithSummary(p.summariesByName[key], refs)
		}
		for key, refs := range referenceMap.GetPropReferences() {
			if len(refs) == 0 {
				continue
			}
			p.populateFunctionDefinitionsWithSummary(p.propSummariesByName[key], refs)
		}
	} else {
		referencesByName := jutil.Maps_NewHashMapWithExpectedSize[string, []rhino.Node](max(totalRefs, 12))
		for key, refs := range referenceMap.GetNameReferences() {
			if len(refs) == 0 {
				continue
			}
			referencesByName.Put(key, refs)
		}
		for key, refs := range referenceMap.GetPropReferences() {
			if len(refs) == 0 {
				continue
			}
			propKey := p.allocPropKey(key)
			if existing, ok := referencesByName.Lookup(propKey); ok {
				combined := make([]rhino.Node, 0, len(existing)+len(refs))
				combined = append(combined, existing...)
				combined = append(combined, refs...)
				referencesByName.Put(propKey, combined)
			} else {
				referencesByName.Put(propKey, refs)
			}
		}
		for name, references := range referencesByName.All() {
			p.PopulateFunctionDefinitions(name, references)
		}
	}

	p.putSummaryForFunctionNode(
		PureFunctionIdentifier_IMPLICIT_PURE_FN,
		p.createSummaryInGraph("<implicit pure class ctor"))
}

func (p *PureFunctionIdentifier) allocPropKey(key string) string {
	return PureFunctionIdentifier_PROP_NAME_PREFIX + key
}

// PopulateFunctionDefinitions, for a name and its set of references, records
// the set of functions that may define that name or skiplists the name if
// there are unclear definitions.
//
// name is a variable or property name; references the set of all nodes
// representing R- and L-value references to name.
//
// Java: PureFunctionIdentifier.populateFunctionDefinitions(String, List<Node>)
func (p *PureFunctionIdentifier) PopulateFunctionDefinitions(name string, references []rhino.Node) {
	summaryForName := jutil.CheckNotNull(p.summariesByName[name])
	p.populateFunctionDefinitionsWithSummary(summaryForName, references)
}

func (p *PureFunctionIdentifier) populateFunctionDefinitionsWithSummary(summaryForName *PureFunctionIdentifier_AmbiguatedFunctionSummary, references []rhino.Node) {
	// Make sure we get absolutely every R-value assigned to `name` or at the very least detect
	// there are some we're missing. Overlooking a single R-value would invalidate the analysis.

	invalid := false
	rvaluesAssignedToName := p.scratchCallableLeaves[:0]
	for _, reference := range references {
		// Eliminate any references that we're sure are R-values themselves. Otherwise
		// there's a high probability we'll inspect an R-value for futher R-values. We wouldn't
		// find any, and then we'd have to consider `name` impure.
		if !PureFunctionIdentifier_IsDefinitelyRValue(reference) {
			// For anything that might be an L-reference, get the expression being assigned to it.
			rvalue := NodeUtil_GetRValueOfLValue(reference)
			if rvalue.IsNil() {
				invalid = true
				break
			} else {
				// If the assigned R-value is an analyzable expression, collect all the possible
				// FUNCTIONs that could result from that expression. If the expression isn't analyzable,
				// represent that with `null` so we can skiplist `name`.
				prevLen := len(rvaluesAssignedToName)
				if !pureFunctionIdentifier_collectCallableLeavesSlice(rvalue, &rvaluesAssignedToName, p.astAnalyzer) {
					invalid = true
					break
				}
				if p.validateArtificialPurity && !reference.IsFromExterns() {
					for _, callable := range rvaluesAssignedToName[prevLen:] {
						jsdoc := NodeUtil_GetBestJSDocInfo(callable)
						isArtificiallyPure := jsdoc != nil && jsdoc.IsNoSideEffects()
						if isArtificiallyPure {
							p.artificiallyPureLiteralsForDebugging.Add(callable)
							summaryForName.SetCollectImpureDebuggingReason()
						}
					}
				}
			}
		}
	}
	p.scratchCallableLeaves = rvaluesAssignedToName[:0]

	if len(rvaluesAssignedToName) == 0 || invalid {
		// Any of:
		// - There are no L-values with this name.
		// - There's a an L-value and we can't find the associated R-values.
		// - There's a an L-value with R-values are not all known to be callable.
		summaryForName.SetMutatesGlobalStateAndAllOtherFlags(rhino.NilNode)
	} else {
		for _, rvalue := range rvaluesAssignedToName {
			if rvalue.IsFunction() {
				p.putSummaryForFunctionNode(rvalue, summaryForName)
			} else if NodeUtil_IsUndefined(rvalue) {
			} else {
				rvalueSummary := p.getSummaryForReferenceOrDefault(rvalue)

				p.reverseCallGraph.ConnectNodes(
					rvalueSummary.graphNode,
					PureFunctionIdentifier_SideEffectPropagation_ForAlias(),
					summaryForName.graphNode)
			}
		}
	}
}

// PropagateSideEffects propagates side effect information in
// reverseCallGraph from callees to callers.
//
// This is an iterative process executed until a fixed point, where no
// caller summary would be given new side-effects from from any callee
// summary, is reached.
//
// Java: PureFunctionIdentifier.propagateSideEffects()
func (p *PureFunctionIdentifier) PropagateSideEffects() {
	graph.FixedPointGraphTraversal_NewTraversal[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation](
		graph.FixedPointGraphTraversal_EdgeCallbackFunc[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation](
			func(source *PureFunctionIdentifier_AmbiguatedFunctionSummary,
				edge *PureFunctionIdentifier_SideEffectPropagation,
				destination *PureFunctionIdentifier_AmbiguatedFunctionSummary) bool {
				return edge.Propagate(source, destination)
			})).
		ComputeFixedPoint(p.reverseCallGraph)
}

// MarkPureFunctionCalls sets no side effect property at pure-function call
// sites.
//
// Java: PureFunctionIdentifier.markPureFunctionCalls()
func (p *PureFunctionIdentifier) MarkPureFunctionCalls() {
	for _, callNode := range p.allFunctionCalls {
		calleeSummaries := p.GetSummariesForCallee(callNode)

		// Default to side effects, non-local results
		flags := rhino.NewNode_SideEffectFlags()
		if len(calleeSummaries) == 0 {
			flags.SetAllFlags()
		} else {
			flags.ClearAllFlags()
			for _, calleeSummary := range calleeSummaries {
				jutil.CheckNotNull(calleeSummary)
				if calleeSummary.MutatesGlobalState() {
					flags.SetMutatesGlobalState()
				}

				if calleeSummary.MutatesArguments() {
					flags.SetMutatesArguments()
				}

				if calleeSummary.FunctionThrows() {
					flags.SetThrows()
				}

				if PureFunctionIdentifier_IsCallOrTaggedTemplateLit(callNode) {
					if calleeSummary.MutatesThis() {
						// A summary for "f" maps to both "f()" and "f.call()" nodes.
						if PureFunctionIdentifier_IsInvocationViaCallOrApply(callNode) {
							flags.SetMutatesArguments() // `this` is actually an argument.
						} else {
							flags.SetMutatesThis()
						}
					}
				}
			}
		}

		if _, isSuper := pureFunctionIdentifier_invocationCallee(callNode).AsSuper(); isSuper {
			// All `super()` calls (i.e. from subclass constructors) implicitly mutate `this`; they
			// determine its value in the caller scope. Concretely, `super()` calls must not be removed
			// or reordered. Marking them this way ensures that without pinning the enclosing function.
			flags.SetMutatesThis()
		}

		// Handle special cases (Math, RegExp)
		if PureFunctionIdentifier_IsCallOrTaggedTemplateLit(callNode) {
			if !p.astAnalyzer.FunctionCallHasSideEffects(callNode) {
				flags.ClearAllFlags()
			}
		} else if callNode.IsNew() {
			// Handle known cases now (Object, Date, RegExp, etc)
			if !p.astAnalyzer.ConstructorCallHasSideEffects(callNode) {
				flags.ClearAllFlags()
			}
		}

		if callNode.GetSideEffectFlags() != flags.ValueOf() {
			callNode.SetSideEffectFlagsObject(flags)
			p.compiler.ReportChangeToEnclosingScope(callNode)
		}
	}
}

// PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer inspects function
// JSDoc for side effects and applies them to the associated
// AmbiguatedFunctionSummary.
//
// This callback is only meant for use on externs.
type PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer struct {
	outer *PureFunctionIdentifier
}

// Java: PureFunctionIdentifier.ExternFunctionAnnotationAnalyzer()
func NewPureFunctionIdentifier_ExternFunctionAnnotationAnalyzer(outer *PureFunctionIdentifier) *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer {
	return &PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer{outer: outer}
}

// Java: PureFunctionIdentifier.ExternFunctionAnnotationAnalyzer.shouldTraverse(NodeTraversal, Node, Node)
func (a *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) ShouldTraverse(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) bool {
	return true
}

// Java: PureFunctionIdentifier.ExternFunctionAnnotationAnalyzer.visit(NodeTraversal, Node, Node)
func (a *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) Visit(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) {
	if !node.IsFunction() {
		return
	}

	for _, definitionSummary := range a.outer.summariesForAllNamesOfFunctionByNode[node] {
		a.UpdateSideEffectsForExternFunction(node, definitionSummary)
	}
}

// Java: PureFunctionIdentifier.ExternFunctionAnnotationAnalyzer.updateSideEffectsForExternFunction(Node, AmbiguatedFunctionSummary)
func (a *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) UpdateSideEffectsForExternFunction(
	externFunction rhino.Node, summary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	jutil.CheckArgument(externFunction.IsFunction())
	jutil.CheckArgument(externFunction.IsFromExterns())

	info := NodeUtil_GetBestJSDocInfo(externFunction)
	if info == nil {
		// We don't know anything about this function so we assume it has side effects.
		summary.SetMutatesGlobalStateAndAllOtherFlags(externFunction)
		return
	}

	if info.ModifiesThis() {
		summary.SetMutatesThis(externFunction)
	}
	if info.HasSideEffectsArgumentsAnnotation() {
		summary.SetMutatesArguments(externFunction)
	}
	if len(info.GetThrowsAnnotations()) != 0 {
		summary.SetThrows(externFunction)
	}

	if !info.IsNoSideEffects() && summary.HasNoFlagsSet() {
		// We don't know anything about this function so we assume it has side effects.
		summary.SetMutatesGlobalStateAndAllOtherFlags(externFunction)
	}
}

// Java: PureFunctionIdentifier.RHS_IS_ALWAYS_LOCAL
func PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL(lhs rhino.Node) bool { return true }

// Java: PureFunctionIdentifier.RHS_IS_NEVER_LOCAL
func PureFunctionIdentifier_RHS_IS_NEVER_LOCAL(lhs rhino.Node) bool { return false }

// Java: PureFunctionIdentifier.FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE
func PureFunctionIdentifier_FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE(lhs rhino.Node) bool {
	rhs := NodeUtil_GetRValueOfLValue(lhs)
	return rhs.IsNil() || NodeUtil_EvaluatesToLocalValue(rhs)
}

// PureFunctionIdentifier_AmbiguatedFunctionSummary is a summary for the set
// of functions that share a particular name.
//
// Side-effects of the functions are the most significant aspect of this
// summary. Because the functions are "ambiguated", the recorded side-effects
// are the union of all side effects detected in any member of the set.
//
// Name in this context refers to a short name, not a qualified name; only
// the last segment of a qualified name is used.
type PureFunctionIdentifier_AmbiguatedFunctionSummary struct {
	// The name shared by the set of functions that defined this summary.
	name string
	// The node holding this summary in the reverse call graph.
	graphNode graph.DiGraph_DiGraphNode[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]
	// The side effect flags for this set of functions.
	// TODO(nickreid): Replace this with a `Node.SideEffectFlags`.
	bitmask int
	// The indices of the parameters (0..63) mutated by this function, if MUTATES_ARGUMENTS is set.
	mutatedParamMask             int64
	collectImpureDebuggingReason bool
	// FUNCTION nodes with a name corresponding to this summary that are
	// impure. Only ever initialized if `collectImpureDebuggingReason` is
	// true, and if PureFunctionIdentifier finds that this summary has a
	// non-extern function that's "artificially pure", i.e. has a
	// @nosideeffects annotation.
	impureFunctionReasonsForDebugging *jutil.LinkedHashSet[rhino.Node] // lazily initialized
	selfSlice                         [1]*PureFunctionIdentifier_AmbiguatedFunctionSummary
}

// Side effect types:
const (
	PureFunctionIdentifier_AmbiguatedFunctionSummary_THROWS               = 1 << 0
	PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_GLOBAL_STATE = 1 << 1
	PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_THIS         = 1 << 2
	PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_ARGUMENTS    = 1 << 3
)

// PureFunctionIdentifier_AmbiguatedFunctionSummary_CreateInGraph adds a new
// summary node to graph, storing the node and returning the summary.
//
// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.createInGraph(DiGraph<AmbiguatedFunctionSummary, SideEffectPropagation>, String)
func PureFunctionIdentifier_AmbiguatedFunctionSummary_CreateInGraph(
	g graph.DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation], name string) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	return NewPureFunctionIdentifier_AmbiguatedFunctionSummary(g, name)
}

// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary(DiGraph<AmbiguatedFunctionSummary, SideEffectPropagation>, String)
func NewPureFunctionIdentifier_AmbiguatedFunctionSummary(
	g graph.DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation], name string) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	s := &PureFunctionIdentifier_AmbiguatedFunctionSummary{bitmask: 0}
	s.name = name // Java: checkNotNull(name); a Go string is never null.
	s.graphNode = g.CreateDiGraphNode(s)
	s.selfSlice[0] = s
	return s
}

// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.setMask(int)
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetMask(mask int) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	s.bitmask |= mask
	return s
}

// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.getMask(int)
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) GetMask(mask int) bool {
	return (s.bitmask & mask) != 0
}

// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.maybeCollectImpureDebuggingReason(Node)
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) MaybeCollectImpureDebuggingReason(debuggingReason rhino.Node) {
	if !s.collectImpureDebuggingReason || debuggingReason.IsNil() {
		return
	}
	if s.impureFunctionReasonsForDebugging == nil {
		s.impureFunctionReasonsForDebugging = jutil.NewLinkedHashSet[rhino.Node]()
	}
	s.impureFunctionReasonsForDebugging.Add(NodeUtil_GetEnclosingFunction(debuggingReason))
}

// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.mutatesThis()
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) MutatesThis() bool {
	// MUTATES_GLOBAL_STATE implies MUTATES_THIS
	return s.GetMask(PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_THIS)
}

// SetMutatesThis marks the function as having "modifies this" side effects.
//
// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.setMutatesThis(Node)
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetMutatesThis(debuggingReason rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	s.MaybeCollectImpureDebuggingReason(debuggingReason)
	return s.SetMask(PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_THIS)
}

// FunctionThrows returns true if function has an explicit "throw".
//
// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.functionThrows()
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) FunctionThrows() bool {
	// MUTATES_GLOBAL_STATE implies THROWS
	return s.GetMask(PureFunctionIdentifier_AmbiguatedFunctionSummary_THROWS)
}

// SetThrows marks the function as having "throw" side effects.
//
// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.setThrows(Node)
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetThrows(debuggingReason rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	s.MaybeCollectImpureDebuggingReason(debuggingReason)
	return s.SetMask(PureFunctionIdentifier_AmbiguatedFunctionSummary_THROWS)
}

// MutatesGlobalState returns true if function mutates global state.
//
// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.mutatesGlobalState()
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) MutatesGlobalState() bool {
	return s.GetMask(PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_GLOBAL_STATE)
}

// SetMutatesGlobalStateAndAllOtherFlags marks the function as having
// "modifies globals" side effects.
//
// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.setMutatesGlobalStateAndAllOtherFlags(Node)
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetMutatesGlobalStateAndAllOtherFlags(debuggingReason rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	s.MaybeCollectImpureDebuggingReason(debuggingReason)
	s.mutatedParamMask = ^int64(0)
	return s.SetMask(PureFunctionIdentifier_AmbiguatedFunctionSummary_THROWS | PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_THIS | PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_ARGUMENTS | PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_GLOBAL_STATE)
}

// MutatesArguments returns true if function mutates its arguments.
//
// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.mutatesArguments()
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) MutatesArguments() bool {
	// MUTATES_GLOBAL_STATE implies MUTATES_ARGUMENTS
	return s.GetMask(PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_ARGUMENTS)
}

// SetMutatesArguments marks the function as having "modifies arguments"
// side effects.
//
// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.setMutatesArguments(Node)
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetMutatesArguments(debuggingReason rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	s.MaybeCollectImpureDebuggingReason(debuggingReason)
	s.mutatedParamMask = ^int64(0)
	return s.SetMask(PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_ARGUMENTS)
}

// SetMutatesArgumentAt marks the function as mutating the parameter at
// paramIndex.
//
// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.setMutatesArgumentAt(int, Node)
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetMutatesArgumentAt(paramIndex int, debuggingReason rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	jutil.CheckArgumentMsg1(paramIndex >= 0, "%s", paramIndex)
	s.MaybeCollectImpureDebuggingReason(debuggingReason)
	// A parameter that doesn't fit in the mask is approximated by all the parameters.
	var paramBit int64
	if paramIndex < 64 /* Long.SIZE */ {
		paramBit = int64(1) << paramIndex
	} else {
		paramBit = ^int64(0)
	}
	s.mutatedParamMask |= paramBit
	return s.SetMask(PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_ARGUMENTS)
}

// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.hasNoFlagsSet()
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) HasNoFlagsSet() bool {
	return s.bitmask == 0
}

// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.setCollectImpureDebuggingReason()
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetCollectImpureDebuggingReason() {
	s.collectImpureDebuggingReason = true
}

// String is for debugging only (@DoNotCall in Java).
//
// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.toString()
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) String() string {
	return jutil.MoreObjects_ToStringHelper("AmbiguatedFunctionSummary").
		Add("name", s.name).
		// Passing `graphNode` directly causes recursion as its `toString` calls `toString` on the
		// summary it contains.
		Add("graphNode", jutil.IdentityHashCode(s.graphNode)).
		Add("sideEffects", s.SideEffectsToString()).
		String()
}

// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary.sideEffectsToString()
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) SideEffectsToString() string {
	var status []string
	if s.MutatesThis() {
		status = append(status, "this")
	}

	if s.MutatesGlobalState() {
		status = append(status, "global")
	}

	if s.MutatesArguments() {
		status = append(status, "args")
	}

	if s.FunctionThrows() {
		status = append(status, "throw")
	}

	return jutil.CollectionToString(status)
}

// PureFunctionIdentifier_Driver is a compiler pass that constructs a
// reference graph and drives the PureFunctionIdentifier across it.
type PureFunctionIdentifier_Driver struct {
	compiler AbstractCompiler
}

// Java: PureFunctionIdentifier.Driver(AbstractCompiler)
func NewPureFunctionIdentifier_Driver(compiler AbstractCompiler) *PureFunctionIdentifier_Driver {
	return &PureFunctionIdentifier_Driver{compiler: compiler}
}

// Java: PureFunctionIdentifier.Driver.process(Node, Node)
func (d *PureFunctionIdentifier_Driver) Process(externs rhino.Node, root rhino.Node) {
	shouldValidate :=
		d.compiler.
			GetOptions().
			GetWarningsGuard().
			Level(JSError_Make(PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION, "", ""))
	OptimizeCalls_Builder_().
		SetCompiler(d.compiler).
		SetConsiderExterns(true).
		AddPass(
			NewPureFunctionIdentifier(
				d.compiler,
				d.compiler.GetOptions().GetAssumeGettersArePure(),
				shouldValidate == nil || *shouldValidate != CheckLevel_OFF)).
		Build().
		Process(externs, root)
}
