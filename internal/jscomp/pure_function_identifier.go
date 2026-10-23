// Port of com/google/javascript/jscomp/PureFunctionIdentifier.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	summariesForAllNamesOfFunctionByNode map[*rhino.Node][]*PureFunctionIdentifier_AmbiguatedFunctionSummary

	// Set of FUNCTION nodes which are artificially marked as pure, via a
	// `@nosideeffects` JSDoc tag in source code, which triggers the compiler
	// to ignore any actual side effects.
	//
	// We store this here instead of in the AmbiguatedFunctionSummary because
	// 1) this is only for debugging, and we don't want to accidentally depend
	// on it for side effect analysis and 2) we expect the set of functions
	// with artificial purity annotations to be small and not scale linearly
	// with the size of the application.
	artificiallyPureLiteralsForDebugging *jutil.LinkedHashSet[*rhino.Node]

	// List of all function call sites. Storing them here during the function
	// analysis traversal prevents us from doing a second traversal to annotate
	// them with side-effects. We can just iterate the list.
	allFunctionCalls []*rhino.Node

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
	reverseCallGraph *LinkedDirectedGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]

	// A summary for a function for which no definition was found.
	//
	// We assume it has all possible side-effects. It's useful for references
	// like function parameters, or inner functions.
	unknownFunctionSummary      *PureFunctionIdentifier_AmbiguatedFunctionSummary
	unknownFunctionSummarySlice []*PureFunctionIdentifier_AmbiguatedFunctionSummary

	summarySlab           []PureFunctionIdentifier_AmbiguatedFunctionSummary
	propSlab              []PureFunctionIdentifier_SideEffectPropagation
	scratchCallableLeaves []*rhino.Node

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
	p.artificiallyPureLiteralsForDebugging = jutil.NewLinkedHashSet[*rhino.Node]()
	p.reverseCallGraph = LinkedDirectedGraph_CreateWithoutAnnotations[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]()
	p.unknownFunctionSummary =
		p.createSummaryInGraph("<unknown>").
			SetMutatesGlobalStateAndAllOtherFlags(nil)
	p.unknownFunctionSummarySlice = p.unknownFunctionSummary.selfSlice[:]
	p.hasProcessed = false
	p.compiler = jutil.CheckNotNull(compiler)
	p.assumeGettersArePure = assumeGettersArePure
	p.astAnalyzer = compiler.GetAstAnalyzer()
	p.validateArtificialPurity = validateArtificialPurity
	return p
}

func (p *PureFunctionIdentifier) createSummaryInGraph(name string) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	if len(p.summarySlab) == 0 {
		p.summarySlab = make([]PureFunctionIdentifier_AmbiguatedFunctionSummary, 256)
	}
	s := &p.summarySlab[0]
	p.summarySlab = p.summarySlab[1:]
	s.name = name
	n := p.reverseCallGraph.allocNode(s)
	p.reverseCallGraph.nodes.PutNew(s, n)
	s.graphNode = n
	s.selfSlice[0] = s
	return s
}

func (p *PureFunctionIdentifier) putSummaryForFunctionNode(fn *rhino.Node, summary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
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

func (p *PureFunctionIdentifier) getSummaryForReference(nameRef *rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	switch nameRef.GetToken() {
	case rhino.Token_NAME:
		return p.summariesByName[nameRef.GetString()]
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		return p.propSummariesByName[nameRef.GetString()]
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected name reference: " + nameRef.String()))
	}
}

func (p *PureFunctionIdentifier) getSummaryForReferenceOrDefault(nameRef *rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	if s := p.getSummaryForReference(nameRef); s != nil {
		return s
	}
	return p.unknownFunctionSummary
}

func (p *PureFunctionIdentifier) allocPropagationForInvocation(
	invocation *rhino.Node, callerIsArrowFunction bool, propagateThrows bool) *PureFunctionIdentifier_SideEffectPropagation {
	jutil.CheckArgumentMsg1(NodeUtil_IsInvocation(invocation), "%s", invocation)
	if len(p.propSlab) == 0 {
		p.propSlab = make([]PureFunctionIdentifier_SideEffectPropagation, 256)
	}
	prop := &p.propSlab[0]
	p.propSlab = p.propSlab[1:]
	prop.allArgsUnescapedLocal = NodeUtil_AllArgsUnescapedLocal(invocation)
	prop.calleeThisEqualsCallerThis = PureFunctionIdentifier_SideEffectPropagation_CalleeAndCallerShareThis(invocation) && !callerIsArrowFunction
	prop.propagateThrows = propagateThrows
	prop.invocation = invocation
	return prop
}

// Java: PureFunctionIdentifier.process(Node, Node, ReferenceMap)
func (p *PureFunctionIdentifier) Process(externs *rhino.Node, root *rhino.Node, references *OptimizeCalls_ReferenceMap) {
	jutil.CheckState(p.compiler.GetLifeCycleStage().IsNormalized())
	jutil.CheckStateMsg(
		!p.hasProcessed, "PureFunctionIdentifier::process may only be called once per instance.")
	p.hasProcessed = true

	p.PopulateDatastructuresForAnalysisTraversal(references)

	NodeTraversal_Traverse(p.compiler, externs, NewPureFunctionIdentifier_ExternFunctionAnnotationAnalyzer(p))
	sc := NewSyntacticScopeCreator(p.compiler)
	sc.EnableArenaPool()
	b := NodeTraversal_Builder{
		compiler:     p.compiler,
		callback:     NewPureFunctionIdentifier_FunctionBodyAnalyzer(p),
		scopeCreator: sc,
	}
	b.Traverse(root)
	sc.ReleaseArena()

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
			var reasons []*rhino.Node
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
func PureFunctionIdentifier_SOURCE_LOCATION_COMPARATOR(a, b *rhino.Node) int {
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
func PureFunctionIdentifier_FormatSourceLocation(node *rhino.Node) string {
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
	expr *rhino.Node, astAnalyzer *AstAnalyzer) []*rhino.Node {
	var callables []*rhino.Node
	if pureFunctionIdentifier_collectCallableLeavesSlice(expr, &callables, astAnalyzer) {
		return callables
	}
	return nil
}

func pureFunctionIdentifier_collectCallableLeavesSlice(
	expr *rhino.Node, results *[]*rhino.Node, astAnalyzer *AstAnalyzer) bool {
	switch expr.GetToken() {
	case rhino.Token_FUNCTION, rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_NAME:
		*results = append(*results, expr)
		return true

	case rhino.Token_SUPER:
		clazz := jutil.CheckNotNull(NodeUtil_GetEnclosingClass(expr))
		return pureFunctionIdentifier_collectCallableLeavesSlice(clazz.GetSecondChild(), results, astAnalyzer)

	case rhino.Token_CLASS:
		members := expr.GetLastChild()
		if members != nil && members.IsClassMembers() {
			for member := members.GetFirstChild(); member != nil; member = member.GetNext() {
				initializer := member.GetLastChild()
				if initializer != nil &&
					NodeUtil_IsNonStaticClassFieldInitializer(initializer) &&
					astAnalyzer.MayHaveSideEffects(initializer) {
					return false
				}
			}
		}

		ctorDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(expr)
		if ctorDef != nil {
			return pureFunctionIdentifier_collectCallableLeavesSlice(ctorDef.GetOnlyChild(), results, astAnalyzer)
		} else if expr.GetSecondChild().IsEmpty() {
			*results = append(*results, PureFunctionIdentifier_IMPLICIT_PURE_FN)
			return true
		} else {
			return pureFunctionIdentifier_collectCallableLeavesSlice(expr.GetSecondChild(), results, astAnalyzer)
		}

	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		return pureFunctionIdentifier_collectCallableLeavesSlice(expr.GetFirstChild(), results, astAnalyzer) &&
			pureFunctionIdentifier_collectCallableLeavesSlice(expr.GetSecondChild(), results, astAnalyzer)

	case rhino.Token_COMMA, rhino.Token_ASSIGN:
		return pureFunctionIdentifier_collectCallableLeavesSlice(expr.GetSecondChild(), results, astAnalyzer)

	case rhino.Token_HOOK:
		return pureFunctionIdentifier_collectCallableLeavesSlice(expr.GetSecondChild(), results, astAnalyzer) &&
			pureFunctionIdentifier_collectCallableLeavesSlice(expr.GetChildAtIndex(2), results, astAnalyzer)

	case rhino.Token_NEW_TARGET, rhino.Token_THIS:
		return false
	default:
		return false
	}
}

// PureFunctionIdentifier_CollectCallableLeavesInternal traverses an expr to
// collect nodes representing potential callables that it may resolve to well
// known callables.
//
// For example:
//
//	`a.c || b` => [a.c, b]
//	`x ? a.c : b` => [a.c, b]
//	`(function f() { }) && x || class Foo { constructor() { } }` => [function, x, constructor]`
//
// This function is applicable to finding both assignment aliases and call
// targets. That is, one way to find the possible callees of an invocation is
// to pass the complex expression representing the final callee to this
// method.
//
// If a node that isn't understood is detected, false is returned and the
// caller is expected to invalidate the entire collection. If true is
// returned, this method is guaranteed to have added at least one result to
// the collection.
//
// See collectCallableLeaves. expr is a possibly complicated expression;
// results the collection of qualified names and functions. Returns true iff
// only understood results were discovered.
//
// Java: PureFunctionIdentifier.collectCallableLeavesInternal(Node, ArrayList<Node>, AstAnalyzer)
func PureFunctionIdentifier_CollectCallableLeavesInternal(
	expr *rhino.Node, results *jutil.ArrayList[*rhino.Node], astAnalyzer *AstAnalyzer) bool {
	switch expr.GetToken() {
	case rhino.Token_FUNCTION, rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_NAME:
		results.Add(expr)
		return true

	case rhino.Token_SUPER:
		{
			// Pretend that `super` is an alias for the superclass reference.
			clazz := jutil.CheckNotNull(NodeUtil_GetEnclosingClass(expr))
			return PureFunctionIdentifier_CollectCallableLeavesInternal(clazz.GetSecondChild(), results, astAnalyzer)
		}

	case rhino.Token_CLASS:
		{
			// Instance field initializers run when the constructor is invoked (during class
			// instantiation), so any side-effectful instance field initializers make the constructor
			// invocation impure. Static field initializers run when the class definition itself is
			// evaluated, not upon constructor invocation.
			members := expr.GetLastChild()
			if members != nil && members.IsClassMembers() {
				for member := members.GetFirstChild(); member != nil; member = member.GetNext() {
					initializer := member.GetLastChild()
					if initializer != nil &&
						NodeUtil_IsNonStaticClassFieldInitializer(initializer) &&
						astAnalyzer.MayHaveSideEffects(initializer) {
						return false
					}
				}
			}

			// Collect the constructor function, or failing that, the superclass reference.
			ctorDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(expr)
			if ctorDef != nil {
				return PureFunctionIdentifier_CollectCallableLeavesInternal(ctorDef.GetOnlyChild(), results, astAnalyzer)
			} else if expr.GetSecondChild().IsEmpty() {
				results.Add(PureFunctionIdentifier_IMPLICIT_PURE_FN)
				return true // A class an implicit ctor is pure when there is no superclass.
			} else {
				return PureFunctionIdentifier_CollectCallableLeavesInternal(expr.GetSecondChild(), results, astAnalyzer)
			}
		}

	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		return PureFunctionIdentifier_CollectCallableLeavesInternal(expr.GetFirstChild(), results, astAnalyzer) &&
			PureFunctionIdentifier_CollectCallableLeavesInternal(expr.GetSecondChild(), results, astAnalyzer)

	case rhino.Token_COMMA, rhino.Token_ASSIGN:
		return PureFunctionIdentifier_CollectCallableLeavesInternal(expr.GetSecondChild(), results, astAnalyzer)

	case rhino.Token_HOOK:
		return PureFunctionIdentifier_CollectCallableLeavesInternal(expr.GetSecondChild(), results, astAnalyzer) &&
			PureFunctionIdentifier_CollectCallableLeavesInternal(expr.GetChildAtIndex(2), results, astAnalyzer)

	case rhino.Token_NEW_TARGET, rhino.Token_THIS:
		// These could be an alias to any function. Treat them as an unknown callable.
		return false
	default:
		return false // Unsupported call type.
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
func PureFunctionIdentifier_IsDefinitelyRValue(rvalue *rhino.Node) bool {
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
	case rhino.Token_CASE, rhino.Token_IF, rhino.Token_SWITCH, rhino.Token_WHILE:
		return rvalue.IsFirstChildOf(parent) // the condition is always an r-value
	case rhino.Token_EXPR_RESULT:
		// Extern declarations are sometimes stubs. These must be considered L-values with no
		// associated R-values.
		return !rvalue.IsFromExterns()
	case rhino.Token_ASSIGN, rhino.Token_CLASS: // `extends` clause.
		return rvalue.IsSecondChildOf(parent)
	case rhino.Token_STRING_KEY: // Assignment to an object literal property. Excludes object destructuring.
		return parent.GetParent().IsObjectLit()
	default:
		// Anything not explicitly listed may not be an R-value. We only worry about the likely
		// cases for nominal function values since those are what interest us and its safe to miss
		// some R-values. It's more important that we correctly identify L-values.
		return false
	}
}

// Java: PureFunctionIdentifier.getGoogCacheCallableExpression(Cache)
func (p *PureFunctionIdentifier) GetGoogCacheCallableExpression(cacheCall *CodingConvention_Cache) []*rhino.Node {
	jutil.CheckNotNull(cacheCall)

	valueFnCallables := PureFunctionIdentifier_CollectCallableLeaves(cacheCall.valueFn, p.astAnalyzer)
	if valueFnCallables == nil {
		// ImmutableList.Builder.addAll(null)
		panic(jutil.NewNullPointerException())
	}
	builder := append([]*rhino.Node(nil), valueFnCallables...)
	if cacheCall.keyFn != nil {
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
func (p *PureFunctionIdentifier) GetSummariesForCallee(invocation *rhino.Node) []*PureFunctionIdentifier_AmbiguatedFunctionSummary {
	jutil.CheckArgumentMsg1(NodeUtil_IsInvocation(invocation), "%s", invocation)

	cacheCall := p.compiler.GetCodingConvention().DescribeCachingCall(invocation)

	if cacheCall == nil {
		var singleCallee *rhino.Node
		if PureFunctionIdentifier_IsInvocationViaCallOrApply(invocation) {
			singleCallee = invocation.GetFirstFirstChild()
		} else {
			first := invocation.GetFirstChild()
			switch first.GetToken() {
			case rhino.Token_FUNCTION, rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_NAME:
				singleCallee = first
			}
		}
		if singleCallee != nil {
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

	var callees []*rhino.Node
	if cacheCall != nil {
		callees = p.GetGoogCacheCallableExpression(cacheCall)
	} else {
		p.scratchCallableLeaves = p.scratchCallableLeaves[:0]
		if pureFunctionIdentifier_collectCallableLeavesSlice(invocation.GetFirstChild(), &p.scratchCallableLeaves, p.astAnalyzer) {
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
	newNodes := jutil.NewLinkedHashMapWithCapacity[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *LinkedDirectedGraph_LinkedDiGraphNode[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]](totalNodes)
	for k, v := range p.reverseCallGraph.nodes.All() {
		newNodes.PutNew(k, v)
	}
	p.reverseCallGraph.nodes = *newNodes
	if len(p.summarySlab) < totalNodes {
		p.summarySlab = make([]PureFunctionIdentifier_AmbiguatedFunctionSummary, totalNodes)
	}
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
	p.summariesForAllNamesOfFunctionByNode = make(map[*rhino.Node][]*PureFunctionIdentifier_AmbiguatedFunctionSummary, max(totalRefs/2, 16))
	p.ensureReverseCallGraphCapacity(totalRefs + 2)

	hasDotNameCollision := false
	for key, value := range referenceMap.GetNameReferences() {
		if value.IsEmpty() {
			continue
		}
		jutil.CheckState(key != "" && key != PureFunctionIdentifier_PROP_NAME_PREFIX)
		p.summariesByName[key] = p.createSummaryInGraph(key)
	}
	for key, value := range referenceMap.GetPropReferences() {
		if value.IsEmpty() {
			continue
		}
		jutil.CheckState(key != "")
		propKey := PureFunctionIdentifier_PROP_NAME_PREFIX + key
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
		for key, value := range referenceMap.GetNameReferences() {
			refs := value.Slice()
			if len(refs) == 0 {
				continue
			}
			p.populateFunctionDefinitionsWithSummary(p.summariesByName[key], refs)
		}
		for key, value := range referenceMap.GetPropReferences() {
			refs := value.Slice()
			if len(refs) == 0 {
				continue
			}
			p.populateFunctionDefinitionsWithSummary(p.propSummariesByName[key], refs)
		}
	} else {
		referencesByName := jutil.Maps_NewHashMapWithExpectedSize[string, *jutil.ArrayList[*rhino.Node]](max(totalRefs, 12))
		for key, value := range referenceMap.GetNameReferences() {
			if value.IsEmpty() {
				continue
			}
			referencesByName.Put(key, value)
		}
		for key, value := range referenceMap.GetPropReferences() {
			if value.IsEmpty() {
				continue
			}
			propKey := PureFunctionIdentifier_PROP_NAME_PREFIX + key
			if existing := referencesByName.Get(propKey); existing != nil {
				combined := jutil.NewArrayListWithCapacity[*rhino.Node](existing.Size() + value.Size())
				combined.AddAll(existing.Slice()...)
				combined.AddAll(value.Slice()...)
				referencesByName.Put(propKey, combined)
			} else {
				referencesByName.Put(propKey, value)
			}
		}
		for name, references := range referencesByName.All() {
			p.PopulateFunctionDefinitions(name, references.Slice())
		}
	}

	p.putSummaryForFunctionNode(
		PureFunctionIdentifier_IMPLICIT_PURE_FN,
		p.createSummaryInGraph("<implicit pure class ctor"))
}

// PopulateFunctionDefinitions, for a name and its set of references, records
// the set of functions that may define that name or skiplists the name if
// there are unclear definitions.
//
// name is a variable or property name; references the set of all nodes
// representing R- and L-value references to name.
//
// Java: PureFunctionIdentifier.populateFunctionDefinitions(String, List<Node>)
func (p *PureFunctionIdentifier) PopulateFunctionDefinitions(name string, references []*rhino.Node) {
	summaryForName := jutil.CheckNotNull(p.summariesByName[name])
	p.populateFunctionDefinitionsWithSummary(summaryForName, references)
}

func (p *PureFunctionIdentifier) populateFunctionDefinitionsWithSummary(summaryForName *PureFunctionIdentifier_AmbiguatedFunctionSummary, references []*rhino.Node) {
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
			if rvalue == nil {
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
		summaryForName.SetMutatesGlobalStateAndAllOtherFlags(nil)
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
	FixedPointGraphTraversal_NewTraversal[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation](
		FixedPointGraphTraversal_EdgeCallbackFunc[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation](
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

		if callNode.GetFirstChild().IsSuper() {
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
func (a *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) ShouldTraverse(traversal *NodeTraversal, node *rhino.Node, parent *rhino.Node) bool {
	return true
}

// Java: PureFunctionIdentifier.ExternFunctionAnnotationAnalyzer.visit(NodeTraversal, Node, Node)
func (a *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) Visit(traversal *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	if !node.IsFunction() {
		return
	}

	for _, definitionSummary := range a.outer.summariesForAllNamesOfFunctionByNode[node] {
		a.UpdateSideEffectsForExternFunction(node, definitionSummary)
	}
}

// Java: PureFunctionIdentifier.ExternFunctionAnnotationAnalyzer.updateSideEffectsForExternFunction(Node, AmbiguatedFunctionSummary)
func (a *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) UpdateSideEffectsForExternFunction(
	externFunction *rhino.Node, summary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
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
func PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL(lhs *rhino.Node) bool { return true }

// Java: PureFunctionIdentifier.RHS_IS_NEVER_LOCAL
func PureFunctionIdentifier_RHS_IS_NEVER_LOCAL(lhs *rhino.Node) bool { return false }

// Java: PureFunctionIdentifier.FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE
func PureFunctionIdentifier_FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE(lhs *rhino.Node) bool {
	rhs := NodeUtil_GetRValueOfLValue(lhs)
	return rhs == nil || NodeUtil_EvaluatesToLocalValue(rhs)
}

// PureFunctionIdentifier_FunctionBodyAnalyzer inspects function bodies for
// side effects and applies them to the associated AmbiguatedFunctionSummary.
//
// This callback also fills allFunctionCalls
type PureFunctionIdentifier_FunctionBodyAnalyzer struct {
	outer *PureFunctionIdentifier

	// Preloaded with an entry to represent the global scope.
	//
	// Go-only: language-adapter — the ArrayDeque is only used at its end
	// (getLast/addLast/removeLast), so it is a slice.
	functionScopeStack []PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer()
func NewPureFunctionIdentifier_FunctionBodyAnalyzer(outer *PureFunctionIdentifier) *PureFunctionIdentifier_FunctionBodyAnalyzer {
	b := &PureFunctionIdentifier_FunctionBodyAnalyzer{
		outer:              outer,
		functionScopeStack: make([]PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry, 1, 8),
	}
	return b
}

// functionScopeStackGetLast is ArrayDeque.getLast().
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) functionScopeStackGetLast() *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry {
	if len(b.functionScopeStack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	return &b.functionScopeStack[len(b.functionScopeStack)-1]
}

// PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry is a single
// function literal definition and associated side-effect-related
// information.
//
// We define separate objects for function literal to simplify handling of
// nested function literal definitions - we push a new entry onto the stack
// when entering a new function literal, and pop the entry when leaving the
// function literal.
type PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry struct {
	root                *rhino.Node
	skiplistedVars      []*Var
	skiplistedNameNodes []*rhino.Node
	taintedVars         []*Var
	// Whether this function was marked as artificially pure, i.e. had a
	// `@nosideeffects` annotation. This tells the compiler to ignore side
	// effects from within the function body when deciding whether the given
	// function as a whole is pure. Note that within the function body, we
	// still record side effects on individual variables accesses; so e.g.
	// given
	//   /** @nosideeffects */
	//   function foo() { console.log(1); }
	//   foo();
	// `foo()` will be marked as pure, but `console.log` will still be marked
	// as impure. Also note that this doesn't extend to nested function
	// definitions: a nested function does not inherit artificial purity from
	// its enclosing function, which is why we need to track this on the
	// function stack.
	isArtificiallyPure bool
	catchDepth         int // The number of try-catch blocks around the current node.
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.FunctionStackEntry(Node, boolean)
func NewPureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry(root *rhino.Node, isArtificiallyPure bool) *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry {
	return &PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry{
		catchDepth:         0,
		root:               root,
		isArtificiallyPure: isArtificiallyPure,
	}
}

func (e *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) addSkiplistedVar(v *Var) {
	if slices.Contains(e.skiplistedVars, v) {
		return
	}
	e.skiplistedVars = append(e.skiplistedVars, v)
}

func (e *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) isSkiplistedVar(v *Var) bool {
	if slices.Contains(e.skiplistedVars, v) {
		return true
	}
	for _, node := range e.skiplistedNameNodes {
		if node == v.GetNameNode() {
			return true
		}
		if node.GetStringID() == v.nameID && node.GetParent().IsVar() && v.scope.IsHoistScope() {
			return true
		}
	}
	return false
}

func (e *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) addTaintedVar(v *Var) {
	if slices.Contains(e.taintedVars, v) {
		return
	}
	e.taintedVars = append(e.taintedVars, v)
}

func (e *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) isTaintedVar(v *Var) bool {
	return slices.Contains(e.taintedVars, v)
}

// IsLocalVar reports whether v, as seen from scope, is a variable of this
// function that only holds values local to the current invocation, so that
// mutating them is invisible to the callers.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.FunctionStackEntry.isLocalVar(Var, Scope)
func (e *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) IsLocalVar(v *Var, scope *Scope) bool {
	return PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(v, scope) &&
		!e.isSkiplistedVar(v) &&
		// Parameters and catch values can come from other scopes.
		!v.IsParam() &&
		!v.IsCatch()
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.shouldTraverse(NodeTraversal, Node, Node)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) ShouldTraverse(traversal *NodeTraversal, node *rhino.Node, parent *rhino.Node) bool {
	b.AddToCatchDepthIfTryBlock(node, 1)
	return true
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.visit(NodeTraversal, Node, Node)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) Visit(traversal *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	b.AddToCatchDepthIfTryBlock(node, -1)

	if NodeUtil_IsInvocation(node) {
		// We collect these after filtering for side-effects because there's no point re-processing
		// a known pure call. This analysis is run multiple times, but no optimization will make a
		// pure function impure.
		b.outer.allFunctionCalls = append(b.outer.allFunctionCalls, node)
	}

	enclosingFunction := b.functionScopeStackGetLast()
	root := enclosingFunction.root
	if root == nil || enclosingFunction.isArtificiallyPure {
		// Within artificially pure function literal bodies, we do not ever need to call
		// updateSideEffectsForNode because we intentionally ignore any side effects, and still
		// consider the function as a whole to be pure.
		// Note: this doesn't mean that individual invocations within the function body are always
		// considered pure - those will still be marked in a later stage of PureFunctionIdentifier.
		// We just don't propagate any impurity to the enclosing function literal.
		return
	}

	for _, summary := range b.outer.summariesForAllNamesOfFunctionByNode[root] {
		b.UpdateSideEffectsForNode(jutil.CheckNotNull(summary), traversal, node)
	}
}

// UpdateSideEffectsForNode updates the side effects of summary based on a
// given node.
//
// This node should be known to (possibly have) side effects. This method
// does not check if the node (possibly) has side effects.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.updateSideEffectsForNode(AmbiguatedFunctionSummary, NodeTraversal, Node)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) UpdateSideEffectsForNode(
	encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary, traversal *NodeTraversal, node *rhino.Node) {
	if encloserSummary.MutatesGlobalState() {
		return // Functions with MUTATES_GLOBAL_STATE already have all side-effects set.
	}

	switch node.GetToken() {
	case rhino.Token_ASSIGN:
		{
			// e.g.
			// lhs = rhs;
			// ({x, y} = object);
			lhs := node.GetFirstChild()
			// Consider destructured properties or values to be nonlocal.
			if lhs.IsDestructuringPattern() {
				NodeUtil_VisitLhsNodesInNode(
					node,
					func(lhsNode *rhino.Node) {
						b.visitLhsNode(encloserSummary, traversal, lhsNode, PureFunctionIdentifier_RHS_IS_NEVER_LOCAL)
					})
			} else {
				b.visitLhsNode(encloserSummary, traversal, lhs, PureFunctionIdentifier_FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE)
			}
		}

	case rhino.Token_INC, // e.g. x++;
		rhino.Token_DEC,
		rhino.Token_DELPROP:
		b.visitLhsNode(
			encloserSummary,
			traversal,
			node.GetOnlyChild(),
			// The value assigned by a unary op is always local.
			PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL)

	case rhino.Token_FOR_AWAIT_OF:
		b.DeprecatedSetSideEffectsForControlLoss(encloserSummary) // Control is lost during await.
		// Fall through.
		fallthrough
	case rhino.Token_FOR_OF:
		// e.g.
		// for (const {prop1, prop2} of iterable) {...}
		// for ({prop1: x.p1, prop2: x.p2} of iterable) {...}
		NodeUtil_VisitLhsNodesInNode(
			node,
			func(lhsNode *rhino.Node) {
				b.visitLhsNode(
					encloserSummary,
					traversal,
					lhsNode,
					// The RHS of a for-of must always be an iterable, making it a container, so
					// we can't
					// consider its contents to be local
					PureFunctionIdentifier_RHS_IS_NEVER_LOCAL)
			})
		b.CheckIteratesImpureIterable(node, encloserSummary)

	case rhino.Token_FOR_IN:
		// e.g.
		// for (prop in obj) {...}
		// Also this, though not very useful or readable.
		// for ([char1, char2, ...x.rest] in obj) {...}
		NodeUtil_VisitLhsNodesInNode(
			node,
			func(lhsNode *rhino.Node) {
				b.visitLhsNode(
					encloserSummary,
					traversal,
					lhsNode,
					// A for-in always assigns a string, which is a local value by definition.
					PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL)
			})

	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL, rhino.Token_NEW, rhino.Token_TAGGED_TEMPLATELIT:
		b.VisitCall(encloserSummary, node)

	case rhino.Token_DESTRUCTURING_LHS:
		if NodeUtil_IsAnyFor(node.GetParent()) {
			// This case is handled when visiting the enclosing for loop.
			break
		}
		// Assume the value assigned to each item is potentially global state. This is overly
		// conservative but necessary because in the common case the rhs is not a literal.
		NodeUtil_VisitLhsNodesInNode(
			node.GetParent(),
			func(lhsNode *rhino.Node) {
				b.visitLhsNode(encloserSummary, traversal, lhsNode, PureFunctionIdentifier_RHS_IS_NEVER_LOCAL)
			})

	case rhino.Token_NAME:
		// Local variable declarations are not a side-effect, but we do want to track them.
		if NodeUtil_IsNameDeclaration(node.GetParent()) {
			value := node.GetFirstChild()
			// Assignment to local, if the value isn't a safe local value,
			// new object creation or literal or known primitive result
			// value, add it to the local skiplist.
			if value != nil && !NodeUtil_EvaluatesToLocalValue(value) {
				e := b.functionScopeStackGetLast()
				e.skiplistedNameNodes = append(e.skiplistedNameNodes, node)
			}
		}

	case rhino.Token_THROW:
		b.RecordThrowsBasedOnContext(encloserSummary)

	case rhino.Token_YIELD:
		b.CheckIteratesImpureIterable(node, encloserSummary) // `yield*` triggers iteration.
		// 'yield' throws if the caller calls `.throw` on the generator object.
		b.DeprecatedSetSideEffectsForControlLoss(encloserSummary)

	case rhino.Token_AWAIT:
		// 'await' throws if the promise it's waiting on is rejected.
		b.DeprecatedSetSideEffectsForControlLoss(encloserSummary)

	case rhino.Token_OBJECT_REST, rhino.Token_OBJECT_SPREAD:
		if !b.outer.assumeGettersArePure {
			// May trigger a getter.
			encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
		}

	case rhino.Token_ITER_REST, rhino.Token_ITER_SPREAD:
		b.CheckIteratesImpureIterable(node, encloserSummary)

	case rhino.Token_STRING_KEY:
		if node.GetParent().IsObjectPattern() {
			// This is an l-value STRING_KEY.
			// Assumption: GETELEM (via a COMPUTED_PROP) is never side-effectful.
			if b.outer.GetPropertyKind(node.GetString()).HasGetter() {
				encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
			}
		}

	case rhino.Token_OPTCHAIN_GETPROP, rhino.Token_GETPROP:
		// Assumption: GETELEM and OPTCHAIN_GETELEM are never side-effectful.
		if b.outer.GetPropertyKind(node.GetString()).HasGetterOrSetter() {
			encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
		}

	case rhino.Token_DYNAMIC_IMPORT:
		// Modules may be imported for side-effects only. This is frequently
		// a pattern used to load polyfills.
		encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)

	default:
		if NodeUtil_IsCompoundAssignmentOp(node) {
			// e.g.
			// x += 3;
			b.visitLhsNode(
				encloserSummary,
				traversal,
				node.GetFirstChild(),
				// The update assignments (e.g. `+=) always assign primitive, and therefore local,
				// values.
				PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL)
			break
		}

		if b.outer.astAnalyzer.NodeTypeMayHaveSideEffects(node) {
			panic(jutil.NewIllegalArgumentExceptionMsg("Unhandled side effect node type " + node.String()))
		}
	}
}

// CheckIteratesImpureIterable inspects node for impure iteration and
// assigns the appropriate side-effects to encloserSummary if so.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.checkIteratesImpureIterable(Node, AmbiguatedFunctionSummary)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) CheckIteratesImpureIterable(node *rhino.Node, encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	if !NodeUtil_IteratesImpureIterable(node) {
		return
	}
	encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
}

// DeprecatedSetSideEffectsForControlLoss assigns the set of side-effects
// associated with an arbitrary loss of control flow to encloserSummary.
//
// This function is kept to retain behaviour but marks places where the
// analysis is inaccurate.
//
// See b/135475880.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.deprecatedSetSideEffectsForControlLoss(AmbiguatedFunctionSummary)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) DeprecatedSetSideEffectsForControlLoss(encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	b.RecordThrowsBasedOnContext(encloserSummary)
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.recordThrowsBasedOnContext(AmbiguatedFunctionSummary)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) RecordThrowsBasedOnContext(encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	if b.functionScopeStackGetLast().catchDepth == 0 {
		encloserSummary.SetThrows(b.functionScopeStackGetLast().root)
	}
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.enterScope(NodeTraversal)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) EnterScope(t *NodeTraversal) {
	function := t.GetScopeRoot()
	if !function.IsFunction() {
		if b.functionScopeStackGetLast().root == nil {
			// Variables declared outside all functions can never share a container
			// scope with any scope inside a function. Pre-populating an empty
			// Scope avoids scanning the global/module AST in SyntacticScopeCreator.
			if len(t.scopes) == 1 {
				t.scopes[0].untypedScope = Scope_CreateGlobalScope(function)
			} else if parentScope := t.scopes[len(t.scopes)-2].untypedScope; parentScope != nil {
				t.scopes[len(t.scopes)-1].untypedScope = Scope_CreateChildScope(parentScope, function)
			}
		}
		return
	}
	jsdoc := NodeUtil_GetBestJSDocInfo(function)
	isArtificiallyPure := jsdoc != nil && jsdoc.IsNoSideEffects()

	n := len(b.functionScopeStack)
	if n < cap(b.functionScopeStack) {
		b.functionScopeStack = b.functionScopeStack[:n+1]
		e := &b.functionScopeStack[n]
		e.root = function
		e.isArtificiallyPure = isArtificiallyPure
		e.catchDepth = 0
		e.skiplistedVars = e.skiplistedVars[:0]
		e.skiplistedNameNodes = e.skiplistedNameNodes[:0]
		e.taintedVars = e.taintedVars[:0]
	} else {
		b.functionScopeStack = append(b.functionScopeStack, PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry{
			root:               function,
			isArtificiallyPure: isArtificiallyPure,
		})
	}
	if len(b.outer.summariesForAllNamesOfFunctionByNode[function]) == 0 {
		// This function was not part of a definition which is why it was not created by
		// {@link populateDatastructuresForAnalysisTraversal}. For example, an anonymous
		// function.
		summary := b.outer.createSummaryInGraph("<anonymous>")
		b.outer.putSummaryForFunctionNode(function, summary)
	}
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.exitScope(NodeTraversal)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) ExitScope(t *NodeTraversal) {
	idx := len(b.functionScopeStack) - 1
	functionEntry := &b.functionScopeStack[idx]
	if functionEntry.root == nil {
		return
	}
	if t.GetScopeRoot() == functionEntry.root {
		// ArrayDeque.removeLast()
		b.functionScopeStack = b.functionScopeStack[:idx]
	}

	if len(functionEntry.taintedVars) == 0 {
		return
	}

	// If the current scope was never materialized, none of its variables could
	// have been looked up or added to taintedVars.
	scope := t.scopes[len(t.scopes)-1].untypedScope
	if scope == nil {
		return
	}
	functionScope := scope.GetClosestContainerScope()
	if !functionScope.IsFunctionScope() {
		return
	}
	jutil.CheckStateMsg1(
		functionScope.GetRootNode() == functionEntry.root, "%s", functionScope.GetRootNode())

	// Handle deferred local variable modifications:
	for _, sideEffectInfo := range b.outer.summariesForAllNamesOfFunctionByNode[functionEntry.root] {
		jutil.CheckNotNullMsg1(sideEffectInfo, "%s has no side effect info.", functionEntry.root)

		if sideEffectInfo.MutatesGlobalState() {
			continue
		}

		for _, v := range scope.varList {
			if v == nil || !functionEntry.isTaintedVar(v) {
				continue
			}
			isFromDestructuring := NodeUtil_IsLhsByDestructuring(v.GetNameNode())
			if v.IsParam() &&
				// Ignore destructuring parameters because they don't directly correspond to an
				// argument passed to the function for the purposes of "setMutatesArguments"
				!isFromDestructuring &&
				!functionEntry.isSkiplistedVar(v) {
				sideEffectInfo.SetMutatesArguments(v.GetNode())
				continue
			}

			// Take care of locals that might have been tainted.
			if !functionEntry.IsLocalVar(v, scope) {
				// If the function has global side-effects
				// don't bother with the local side-effects.
				sideEffectInfo.SetMutatesGlobalStateAndAllOtherFlags(v.GetNode())
				break
			}
		}
	}
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.isVarDeclaredInSameContainerScope(Var, Scope)
func PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(v *Var, scope *Scope) bool {
	if v == nil {
		return false
	}
	varScope, _ := v.GetScope().(*Scope)
	return varScope.HasSameContainerScope(scope)
}

func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) visitLhsNode(
	encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary,
	traversal *NodeTraversal,
	lhs *rhino.Node,
	hasLocalRhs func(*rhino.Node) bool) {
	if NodeUtil_IsNormalOrOptChainGet(lhs) {
		objectNode := lhs.GetFirstChild()
		if objectNode.IsThis() {
			if b.functionScopeStackGetLast().root.IsArrowFunction() {
				encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
			} else {
				encloserSummary.SetMutatesThis(lhs)
			}
		} else if objectNode.IsName() {
			scope := traversal.GetScope()
			var_ := scope.GetVar(objectNode.GetString())
			if PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(var_, scope) {
				b.functionScopeStackGetLast().addTaintedVar(var_)
			} else {
				encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
			}
		} else {
			encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
		}
	} else {
		jutil.CheckStateMsg1(lhs.IsName(), "%s", lhs)
		scope := traversal.GetScope()
		var_ := scope.GetVar(lhs.GetString())
		if PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(var_, scope) {
			if !hasLocalRhs(lhs) {
				b.functionScopeStackGetLast().addSkiplistedVar(var_)
			}
		} else {
			encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
		}
	}
}

// VisitLhsNode records information about the side effects caused by
// assigning a value to a given LHS.
//
// If the operation modifies this or taints global state, mark the enclosing
// function as having those side effects.
//
// encloserSummary is the function side effect record to be updated; scope
// the variable scope in which the variable assignment occurs; lhs the LHS
// node assigned a value by a given parent node; hasLocalRhs a predicate
// indicating whether a given LHS is being assigned a local value.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.visitLhsNode(AmbiguatedFunctionSummary, Scope, Node, Predicate<Node>)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) VisitLhsNode(
	encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary,
	scope *Scope,
	lhs *rhino.Node,
	hasLocalRhs func(*rhino.Node) bool) {
	if NodeUtil_IsNormalOrOptChainGet(lhs) {
		// Although OPTCHAIN_GETPROP can not be an LHS of an assign, it can be a child to DELPROP.
		// e.g. `delete obj?.prop` <==> `obj == null ?  true : delete obj.prop;`
		// Hence the enclosing function's side effects must be recorded.
		if lhs.GetFirstChild().IsThis() {
			if b.functionScopeStackGetLast().root.IsArrowFunction() {
				// `this` in an arrow function is the `this` of the enclosing function, which is not
				// the receiver the arrow function is called with.
				encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
			} else {
				encloserSummary.SetMutatesThis(lhs)
			}
		} else {
			objectNode := lhs.GetFirstChild()
			if objectNode.IsName() {
				var_ := scope.GetVar(objectNode.GetString())
				if PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(var_, scope) {
					// Maybe a local object modification.  We won't know for sure until
					// we exit the scope and can validate the value of the local.
					b.functionScopeStackGetLast().addTaintedVar(var_)
				} else {
					encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
				}
			} else {
				// Don't track multi level locals: local.prop.prop2++;
				encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
			}
		}
	} else {
		jutil.CheckStateMsg1(lhs.IsName(), "%s", lhs)
		var_ := scope.GetVar(lhs.GetString())
		if PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(var_, scope) {
			if !hasLocalRhs(lhs) {
				// Assigned value is not guaranteed to be a local value,
				// so if we see any property assignments on this variable,
				// they could be tainting a non-local value.
				b.functionScopeStackGetLast().addSkiplistedVar(var_)
			}
		} else {
			encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
		}
	}
}

// VisitCall records information about a call site.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.visitCall(AmbiguatedFunctionSummary, Node)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) VisitCall(callerInfo *PureFunctionIdentifier_AmbiguatedFunctionSummary, invocation *rhino.Node) {
	// Handle special cases (Math, RegExp)
	// TODO: This logic can probably be replaced with @nosideeffects annotations in externs.
	if invocation.IsCall() && !b.outer.astAnalyzer.FunctionCallHasSideEffects(invocation) {
		return
	}

	// Handle known cases now (Object, Date, RegExp, etc)
	if invocation.IsNew() && !b.outer.astAnalyzer.ConstructorCallHasSideEffects(invocation) {
		return
	}

	calleeSummaries := b.outer.GetSummariesForCallee(invocation)
	if len(calleeSummaries) == 0 {
		callerInfo.SetMutatesGlobalStateAndAllOtherFlags(nil)
		return
	}

	enclosingFunction := b.functionScopeStackGetLast()
	propatesThrows := enclosingFunction.catchDepth == 0
	callerIsArrowFunction := enclosingFunction.root.IsArrowFunction()
	for _, calleeInfo := range calleeSummaries {
		edge := b.outer.allocPropagationForInvocation(invocation, callerIsArrowFunction, propatesThrows)
		b.outer.reverseCallGraph.ConnectNodes(calleeInfo.graphNode, edge, callerInfo.graphNode)
	}
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.addToCatchDepthIfTryBlock(Node, int)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) AddToCatchDepthIfTryBlock(n *rhino.Node, delta int) {
	parent := n.GetParent()
	if !n.IsBlock() || !parent.IsTry() || !n.IsFirstChildOf(parent) {
		return
	}

	jsCatch := n.GetNext().GetFirstChild()
	if jsCatch == nil {
		return
	}

	b.functionScopeStackGetLast().catchDepth += delta
}

// Java: PureFunctionIdentifier.isInvocationViaCallOrApply(Node)
func PureFunctionIdentifier_IsInvocationViaCallOrApply(callSite *rhino.Node) bool {
	receiver := callSite.GetFirstFirstChild()
	if receiver == nil ||
		!(receiver.IsName() || receiver.IsGetProp() || receiver.IsOptChainGetProp()) {
		return false
	}

	return NodeUtil_IsFunctionObjectCall(callSite) || NodeUtil_IsFunctionObjectApply(callSite)
}

// Java: PureFunctionIdentifier.isCallOrTaggedTemplateLit(Node)
func PureFunctionIdentifier_IsCallOrTaggedTemplateLit(invocation *rhino.Node) bool {
	return invocation.IsCall() || invocation.IsOptChainCall() || invocation.IsTaggedTemplateLit()
}

// PureFunctionIdentifier_NameForReference returns the unqualified name
// associated with an R-value.
//
// For NAMEs this is the name. For GETPROPs this is the last segment
// including a leading dot.
//
// Java: PureFunctionIdentifier.nameForReference(Node)
func PureFunctionIdentifier_NameForReference(nameRef *rhino.Node) string {
	switch nameRef.GetToken() {
	case rhino.Token_NAME:
		return nameRef.GetString()
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		return PureFunctionIdentifier_PROP_NAME_PREFIX + nameRef.GetString()
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected name reference: " + nameRef.String()))
	}
}

// Java: PureFunctionIdentifier.getPropertyKind(String)
func (p *PureFunctionIdentifier) GetPropertyKind(name string) AccessorSummary_PropertyAccessKind {
	if p.assumeGettersArePure {
		return AccessorSummary_PropertyAccessKind_NORMAL
	}
	return p.compiler.GetAccessorSummary().GetKind(name)
}

// PureFunctionIdentifier_SideEffectPropagation stores all the information
// about a connection between functions needed to propagate side effects from
// one instance of AmbiguatedFunctionSummary to another.
type PureFunctionIdentifier_SideEffectPropagation struct {
	// Whether this propagation represents an aliasing of one name by another.
	// In that case, all side effects of the "callee" just need to be copied
	// onto the "caller".
	callerIsAlias bool

	// If all the arguments passed to the callee are local to the caller.
	allArgsUnescapedLocal bool

	// If you call a function with apply or call, one of the arguments at the
	// call site will be used as 'this' inside the implementation. If this is
	// pass into apply like so: function.apply(this, ...) then 'this' in the
	// caller is tainted.
	calleeThisEqualsCallerThis bool

	// Whether this propagation includes the "throws" bit.
	//
	// In some contexts, such as an invocation inside a "try", the caller is
	// uneffected by the callee throwing.
	propagateThrows bool

	// The token used to invoke the callee by the caller.
	invocation *rhino.Node
}

// Java: PureFunctionIdentifier.SideEffectPropagation(boolean, boolean, boolean, boolean, Node)
func NewPureFunctionIdentifier_SideEffectPropagation(
	callerIsAlias bool,
	allArgsUnescapedLocal bool,
	calleeThisEqualsCallerThis bool,
	propagateThrows bool,
	invocation *rhino.Node) *PureFunctionIdentifier_SideEffectPropagation {
	jutil.CheckArgumentMsg1(invocation == nil || NodeUtil_IsInvocation(invocation), "%s", invocation)

	return &PureFunctionIdentifier_SideEffectPropagation{
		callerIsAlias:              callerIsAlias,
		allArgsUnescapedLocal:      allArgsUnescapedLocal,
		calleeThisEqualsCallerThis: calleeThisEqualsCallerThis,
		propagateThrows:            propagateThrows,
		invocation:                 invocation,
	}
}

var pureFunctionIdentifier_forAliasSingleton = &PureFunctionIdentifier_SideEffectPropagation{
	callerIsAlias:   true,
	propagateThrows: true,
}

// Java: PureFunctionIdentifier.SideEffectPropagation.forAlias()
func PureFunctionIdentifier_SideEffectPropagation_ForAlias() *PureFunctionIdentifier_SideEffectPropagation {
	return pureFunctionIdentifier_forAliasSingleton
}

// Java: PureFunctionIdentifier.SideEffectPropagation.forInvocation(Node, boolean)
func PureFunctionIdentifier_SideEffectPropagation_ForInvocation(
	invocation *rhino.Node, callerIsArrowFunction bool, propagateThrows bool) *PureFunctionIdentifier_SideEffectPropagation {
	jutil.CheckArgumentMsg1(NodeUtil_IsInvocation(invocation), "%s", invocation)

	return NewPureFunctionIdentifier_SideEffectPropagation(
		false,
		NodeUtil_AllArgsUnescapedLocal(invocation),
		// `this` in an arrow function is the `this` of the enclosing function, which is not
		// tracked.
		PureFunctionIdentifier_SideEffectPropagation_CalleeAndCallerShareThis(invocation) && !callerIsArrowFunction,
		propagateThrows,
		invocation)
}

// Java: PureFunctionIdentifier.SideEffectPropagation.calleeAndCallerShareThis(Node)
func PureFunctionIdentifier_SideEffectPropagation_CalleeAndCallerShareThis(invocation *rhino.Node) bool {
	if !PureFunctionIdentifier_IsCallOrTaggedTemplateLit(invocation) {
		return false // Calling a constructor creates a new object bound to `this`.
	}

	callee := invocation.GetFirstChild()
	if callee.IsSuper() {
		return true
	}

	var thisArg *rhino.Node
	if PureFunctionIdentifier_IsInvocationViaCallOrApply(invocation) {
		// If the call site is actually a `.call` or `.apply`, then `this` will be an argument.
		thisArg = invocation.GetSecondChild()
	} else if callee.IsGetProp() || callee.IsOptChainGetProp() {
		thisArg = callee.GetFirstChild()
	} else {
		thisArg = nil
	}

	if thisArg == nil {
		return false // No `this` is being passed.
	} else if thisArg.IsThis() || thisArg.IsSuper() {
		return true
	}

	// TODO(nickreid): If `thisArg` is a known local or known arg we could say something more
	// specific about the effect of the callee mutating `this`.

	// We're not sure what `this` is being passed, so make a conservative choice.
	return false
}

// Propagate propagates the side effects from the callee to the caller.
//
// callee: propagate from; caller: propagate to. Returns true if the
// propagation changed the side effects on the caller.
//
// Java: PureFunctionIdentifier.SideEffectPropagation.propagate(AmbiguatedFunctionSummary, AmbiguatedFunctionSummary)
func (e *PureFunctionIdentifier_SideEffectPropagation) Propagate(callee *PureFunctionIdentifier_AmbiguatedFunctionSummary, caller *PureFunctionIdentifier_AmbiguatedFunctionSummary) bool {
	initialCallerFlags := caller.bitmask

	if e.callerIsAlias {
		caller.SetMask(callee.bitmask)
	} else if callee.MutatesGlobalState() {
		// If the callee modifies global state then so does that caller.
		caller.SetMutatesGlobalStateAndAllOtherFlags(e.invocation)
	} else {
		if e.propagateThrows && callee.FunctionThrows() {
			// If the callee throws an exception then so does the caller.
			caller.SetThrows(e.invocation)
		}
		if callee.MutatesArguments() && !e.allArgsUnescapedLocal {
			// If the callee mutates its input arguments and the arguments escape the caller then it
			// has unbounded side effects.
			caller.SetMutatesGlobalStateAndAllOtherFlags(e.invocation)
		}
		if callee.MutatesThis() {
			if e.invocation.IsNew() {
				// NEWing a constructor provide a unescaped "this" making side-effects impossible.
			} else if e.calleeThisEqualsCallerThis {
				caller.SetMutatesThis(e.invocation)
			} else {
				caller.SetMutatesGlobalStateAndAllOtherFlags(e.invocation)
			}
		}
	}

	return caller.bitmask != initialCallerFlags
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
	graphNode DiGraph_DiGraphNode[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]
	// The side effect flags for this set of functions.
	// TODO(nickreid): Replace this with a `Node.SideEffectFlags`.
	bitmask                      int
	collectImpureDebuggingReason bool
	// FUNCTION nodes with a name corresponding to this summary that are
	// impure. Only ever initialized if `collectImpureDebuggingReason` is
	// true, and if PureFunctionIdentifier finds that this summary has a
	// non-extern function that's "artificially pure", i.e. has a
	// @nosideeffects annotation.
	impureFunctionReasonsForDebugging *jutil.LinkedHashSet[*rhino.Node] // lazily initialized
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
	graph DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation], name string) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	return NewPureFunctionIdentifier_AmbiguatedFunctionSummary(graph, name)
}

// Java: PureFunctionIdentifier.AmbiguatedFunctionSummary(DiGraph<AmbiguatedFunctionSummary, SideEffectPropagation>, String)
func NewPureFunctionIdentifier_AmbiguatedFunctionSummary(
	graph DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation], name string) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	s := &PureFunctionIdentifier_AmbiguatedFunctionSummary{bitmask: 0}
	s.name = name // Java: checkNotNull(name); a Go string is never null.
	s.graphNode = graph.CreateDiGraphNode(s)
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
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) MaybeCollectImpureDebuggingReason(debuggingReason *rhino.Node) {
	if !s.collectImpureDebuggingReason || debuggingReason == nil {
		return
	}
	if s.impureFunctionReasonsForDebugging == nil {
		s.impureFunctionReasonsForDebugging = jutil.NewLinkedHashSet[*rhino.Node]()
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
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetMutatesThis(debuggingReason *rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
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
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetThrows(debuggingReason *rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
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
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetMutatesGlobalStateAndAllOtherFlags(debuggingReason *rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	s.MaybeCollectImpureDebuggingReason(debuggingReason)
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
func (s *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetMutatesArguments(debuggingReason *rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	s.MaybeCollectImpureDebuggingReason(debuggingReason)
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
func (d *PureFunctionIdentifier_Driver) Process(externs *rhino.Node, root *rhino.Node) {
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
