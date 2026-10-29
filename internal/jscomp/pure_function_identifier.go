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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	pfiDotStringID         = rhino.RhinoStringPool_InternID(".")
	pfiConstructorStringID = rhino.RhinoStringPool_InternID("constructor")
	pfiCallStringID        = rhino.RhinoStringPool_InternID("call")
	pfiApplyStringID       = rhino.RhinoStringPool_InternID("apply")
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
type PureFunctionIdentifier struct {
	compiler                 AbstractCompiler
	astAnalyzer              *AstAnalyzer
	accessorSummary          *AccessorSummary
	codingConvention         CodingConvention
	assumeGettersArePure     bool
	validateArtificialPurity bool
	hasProcessed             bool

	summaries     []pfiSummary
	edges         []pfiEdge
	callSites     []pfiCallSite
	multiCallees  []pfiSummaryID
	nameSummaries map[uint32]pfiSummaryID
	propSummaries map[uint32]pfiSummaryID

	// Maps FUNCTION Node -> pfiSummaryID (when high bit clear) or index into
	// fnMultiSummaries (when pfiMultiSummaryBit is set).
	fnSummaries      map[rhino.Node]uint32
	fnMultiSummaries [][]pfiSummaryID

	scratchCallableLeaves []rhino.Node
	scratchCalleeIDs      []pfiSummaryID

	// Debugging structures, populated only when validateArtificialPurity is true.
	artificiallyPureFns []rhino.Node
	artificiallyPureSet map[rhino.Node]struct{}
	debugReasons        map[pfiSummaryID]*jutil.LinkedHashSet[rhino.Node]
}

// A prefix to differentiate property names from variable names.
const PureFunctionIdentifier_PROP_NAME_PREFIX = "."

// PureFunctionIdentifier_IMPLICIT_PURE_FN is a function node representing a
// function implicit in the AST that is known to be pure.
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
	compiler = jutil.CheckNotNull(compiler)
	return &PureFunctionIdentifier{
		compiler:                 compiler,
		astAnalyzer:              compiler.GetAstAnalyzer(),
		assumeGettersArePure:     assumeGettersArePure,
		validateArtificialPurity: validateArtificialPurity,
	}
}

// Java: PureFunctionIdentifier.process(Node, Node, ReferenceMap)
func (p *PureFunctionIdentifier) Process(externs rhino.Node, root rhino.Node, references *OptimizeCalls_ReferenceMap) {
	jutil.CheckState(p.compiler.GetLifeCycleStage().IsNormalized())
	jutil.CheckStateMsg(
		!p.hasProcessed, "PureFunctionIdentifier::process may only be called once per instance.")
	p.hasProcessed = true
	p.codingConvention = p.compiler.GetCodingConvention()
	if !p.assumeGettersArePure {
		p.accessorSummary = p.compiler.GetAccessorSummary()
	}

	p.populateDatastructuresForAnalysisTraversal(references)

	NodeTraversal_Traverse(p.compiler, externs, &pfiExternFunctionAnnotationAnalyzer{outer: p})
	sc := NewSyntacticScopeCreator(p.compiler)
	b := NodeTraversal_Builder{
		compiler:     p.compiler,
		callback:     newPFIFunctionBodyAnalyzer(p),
		scopeCreator: sc,
	}
	b.Traverse(root)

	p.propagateSideEffects()
	p.validateArtificialPurityAnnotations()
	p.markPureFunctionCalls()

	// Release transient arenas immediately once invocation flags are stamped.
	p.summaries = nil
	p.edges = nil
	p.callSites = nil
	p.multiCallees = nil
	p.nameSummaries = nil
	p.propSummaries = nil
	p.fnSummaries = nil
	p.fnMultiSummaries = nil
	p.scratchCallableLeaves = nil
	p.scratchCalleeIDs = nil
	p.artificiallyPureFns = nil
	p.artificiallyPureSet = nil
	p.debugReasons = nil
}

// validateArtificialPurityAnnotations checks whether any function literals
// marked as artificially pure still have an overall summary with side-effects.
//
// Java: PureFunctionIdentifier.validateArtificialPurity()
func (p *PureFunctionIdentifier) validateArtificialPurityAnnotations() {
	if !p.validateArtificialPurity {
		return
	}
	for _, function := range p.artificiallyPureFns {
		single, multi := p.getFunctionSummaries(function)
		if single != 0 {
			p.validateArtificialPurityForSummary(function, single)
		} else {
			for _, sid := range multi {
				p.validateArtificialPurityForSummary(function, sid)
			}
		}
	}
}

func (p *PureFunctionIdentifier) validateArtificialPurityForSummary(function rhino.Node, sid pfiSummaryID) {
	s := &p.summaries[sid]
	if s.bitmask == 0 {
		return
	}
	if s.nameKind() == pfiNameProp && s.nameID == pfiConstructorStringID {
		// .constructor is a special case - constructors are almost always invoked via
		// `new Foo();` rather than `new something.constructor();`.
		return
	}
	var set *jutil.LinkedHashSet[rhino.Node]
	if p.debugReasons != nil {
		set = p.debugReasons[sid]
	}
	if set == nil {
		// stream() of a null set.
		panic(jutil.NewNullPointerException())
	}
	reasons := set.ToSlice()
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
			s.displayName(),
			sideEffectfulDefinitions))
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
// definitely a reference reading a value.
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
		return false
	}
}

func (p *PureFunctionIdentifier) appendGoogCacheCallableLeaves(cacheCall *CodingConvention_Cache) []rhino.Node {
	jutil.CheckNotNull(cacheCall)
	leaves := p.scratchCallableLeaves[:0]
	if !pureFunctionIdentifier_collectCallableLeavesSlice(cacheCall.valueFn, &leaves, p.astAnalyzer) {
		// ImmutableList.Builder.addAll(null)
		panic(jutil.NewNullPointerException())
	}
	if !cacheCall.keyFn.IsNil() {
		if !pureFunctionIdentifier_collectCallableLeavesSlice(cacheCall.keyFn, &leaves, p.astAnalyzer) {
			// ImmutableList.Builder.addAll(null)
			panic(jutil.NewNullPointerException())
		}
	}
	p.scratchCallableLeaves = leaves
	return leaves
}

// recordCallSite resolves invocation's callee summaries and invocation flags once
// during the analyzer traversal and appends a compact pfiCallSite record.
func (p *PureFunctionIdentifier) recordCallSite(invocation rhino.Node) pfiCallSite {
	var flags uint8
	tok := invocation.GetToken()
	isCallOrTagged := tok == rhino.Token_CALL || tok == rhino.Token_OPTCHAIN_CALL || tok == rhino.Token_TAGGED_TEMPLATELIT
	if isCallOrTagged {
		flags |= pfiCallIsCallOrTagged
		if !p.astAnalyzer.FunctionCallHasSideEffects(invocation) {
			flags |= pfiCallKnownPureForMark
			if tok == rhino.Token_CALL {
				flags |= pfiCallSkipAnalyzer
			}
		}
	} else if tok == rhino.Token_NEW {
		if !p.astAnalyzer.ConstructorCallHasSideEffects(invocation) {
			flags |= pfiCallKnownPureForMark | pfiCallSkipAnalyzer
		}
	}

	callee := pureFunctionIdentifier_invocationCallee(invocation)
	calleeNode := callee.Node()
	if calleeNode.IsSuper() {
		flags |= pfiCallIsSuper
	}

	// Fast exit when neither VisitCall nor MarkPureFunctionCalls needs callee summaries.
	if flags&pfiCallSkipAnalyzer != 0 {
		site := pfiCallSite{callNode: invocation, flags: flags}
		p.callSites = append(p.callSites, site)
		return site
	}

	isViaCallOrApply := PureFunctionIdentifier_IsInvocationViaCallOrApply(invocation)
	if isViaCallOrApply {
		flags |= pfiCallViaCallOrApply
	}

	var cacheCall *CodingConvention_Cache
	if tok == rhino.Token_CALL {
		cacheCall = p.codingConvention.DescribeCachingCall(invocation)
		if cacheCall != nil {
			flags |= pfiCallIsCacheCall
		}
	}

	if cacheCall == nil {
		var singleCallee rhino.Node
		if isViaCallOrApply {
			singleCallee = calleeNode.GetFirstChild()
		} else {
			switch calleeNode.GetToken() {
			case rhino.Token_FUNCTION, rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_NAME:
				singleCallee = calleeNode
			}
		}
		if !singleCallee.IsNil() {
			if singleCallee.IsFunction() {
				single, multi := p.getFunctionSummaries(singleCallee)
				jutil.CheckStateMsg1(single != 0 || len(multi) != 0, "Function missed during analysis: %s", singleCallee)
				if single != 0 {
					site := pfiCallSite{
						callNode:    invocation,
						firstCallee: uint32(single),
						calleeCount: 1,
						flags:       flags,
					}
					p.callSites = append(p.callSites, site)
					return site
				}
				offset := uint32(len(p.multiCallees))
				p.multiCallees = append(p.multiCallees, multi...)
				site := pfiCallSite{
					callNode:    invocation,
					firstCallee: offset,
					calleeCount: uint32(len(multi)),
					flags:       flags,
				}
				p.callSites = append(p.callSites, site)
				return site
			}
			sid := p.getSummaryForReferenceOrDefault(singleCallee)
			site := pfiCallSite{
				callNode:    invocation,
				firstCallee: uint32(sid),
				calleeCount: 1,
				flags:       flags,
			}
			p.callSites = append(p.callSites, site)
			return site
		}
	}

	var callees []rhino.Node
	if cacheCall != nil {
		callees = p.appendGoogCacheCallableLeaves(cacheCall)
	} else {
		p.scratchCallableLeaves = p.scratchCallableLeaves[:0]
		if pureFunctionIdentifier_collectCallableLeavesSlice(calleeNode, &p.scratchCallableLeaves, p.astAnalyzer) {
			callees = p.scratchCallableLeaves
		}
	}

	if callees == nil {
		site := pfiCallSite{
			callNode:    invocation,
			firstCallee: uint32(pfiUnknownSummaryID),
			calleeCount: 1,
			flags:       flags,
		}
		p.callSites = append(p.callSites, site)
		return site
	}

	jutil.CheckStateMsg(len(callees) != 0, "Unexpected empty callees for valid result")
	ids := p.scratchCalleeIDs[:0]
	for _, c := range callees {
		if c.IsFunction() {
			single, multi := p.getFunctionSummaries(c)
			jutil.CheckStateMsg1(single != 0 || len(multi) != 0, "Function missed during analysis: %s", c)
			if single != 0 {
				ids = append(ids, single)
			} else {
				ids = append(ids, multi...)
			}
		} else {
			ids = append(ids, p.getSummaryForReferenceOrDefault(c))
		}
	}
	p.scratchCalleeIDs = ids[:0]

	if len(ids) == 1 {
		site := pfiCallSite{
			callNode:    invocation,
			firstCallee: uint32(ids[0]),
			calleeCount: 1,
			flags:       flags,
		}
		p.callSites = append(p.callSites, site)
		return site
	}
	offset := uint32(len(p.multiCallees))
	p.multiCallees = append(p.multiCallees, ids...)
	site := pfiCallSite{
		callNode:    invocation,
		firstCallee: offset,
		calleeCount: uint32(len(ids)),
		flags:       flags,
	}
	p.callSites = append(p.callSites, site)
	return site
}

// Java: PureFunctionIdentifier.populateDatastructuresForAnalysisTraversal(ReferenceMap)
func (p *PureFunctionIdentifier) populateDatastructuresForAnalysisTraversal(referenceMap *OptimizeCalls_ReferenceMap) {
	nameCount := referenceMap.names.Size()
	propCount := referenceMap.props.Size()
	totalRefs := nameCount + propCount

	p.summaries = make([]pfiSummary, 1, totalRefs+4)
	p.summaries[pfiUnknownSummaryID] = pfiSummary{
		mutatedParamMask: ^uint64(0),
		firstCallerEdge:  -1,
		bitmask:          pfiAllFlags,
		flags:            uint8(pfiNameUnknown) << pfiNameKindShift,
	}
	p.nameSummaries = make(map[uint32]pfiSummaryID, nameCount)
	p.propSummaries = make(map[uint32]pfiSummaryID, propCount)
	p.fnSummaries = make(map[rhino.Node]uint32, max(totalRefs/2, 16))
	p.edges = make([]pfiEdge, 0, max(totalRefs/2, 16))
	p.callSites = make([]pfiCallSite, 0, max(totalRefs/2, 16))

	var dotPrefixedNames map[uint32]pfiSummaryID
	for id, refs := range referenceMap.GetNameReferencesByID() {
		if len(refs) == 0 {
			continue
		}
		jutil.CheckState(id > rhino.RhinoStringPool_EmptyID && id != pfiDotStringID)
		sid := p.allocSummary(id, pfiNameVar)
		p.nameSummaries[id] = sid
		if s := strpool.String(strpool.ID(id)); len(s) > 1 && s[0] == '.' {
			if propID := strpool.Lookup(s[1:]); propID != 0 {
				if dotPrefixedNames == nil {
					dotPrefixedNames = make(map[uint32]pfiSummaryID, 4)
				}
				dotPrefixedNames[uint32(propID)] = sid
			}
		}
	}

	hasDotNameCollision := false
	for id, refs := range referenceMap.GetPropReferencesByID() {
		if len(refs) == 0 {
			continue
		}
		jutil.CheckState(id > rhino.RhinoStringPool_EmptyID)
		if dotPrefixedNames != nil {
			if existing := dotPrefixedNames[id]; existing != 0 {
				p.propSummaries[id] = existing
				hasDotNameCollision = true
				continue
			}
		}
		p.propSummaries[id] = p.allocSummary(id, pfiNameProp)
	}

	if !hasDotNameCollision {
		for id, refs := range referenceMap.GetNameReferencesByID() {
			if len(refs) == 0 {
				continue
			}
			p.populateFunctionDefinitionsWithSummary(p.nameSummaries[id], refs)
		}
		for id, refs := range referenceMap.GetPropReferencesByID() {
			if len(refs) == 0 {
				continue
			}
			p.populateFunctionDefinitionsWithSummary(p.propSummaries[id], refs)
		}
	} else {
		type summaryRefs struct {
			sid  pfiSummaryID
			refs []rhino.Node
		}
		bySummary := make(map[pfiSummaryID]int, totalRefs)
		ordered := make([]summaryRefs, 0, totalRefs)
		for id, refs := range referenceMap.GetNameReferencesByID() {
			if len(refs) == 0 {
				continue
			}
			sid := p.nameSummaries[id]
			bySummary[sid] = len(ordered)
			ordered = append(ordered, summaryRefs{sid: sid, refs: refs})
		}
		for id, refs := range referenceMap.GetPropReferencesByID() {
			if len(refs) == 0 {
				continue
			}
			sid := p.propSummaries[id]
			if idx, ok := bySummary[sid]; ok {
				combined := make([]rhino.Node, 0, len(ordered[idx].refs)+len(refs))
				combined = append(combined, ordered[idx].refs...)
				combined = append(combined, refs...)
				ordered[idx].refs = combined
			} else {
				bySummary[sid] = len(ordered)
				ordered = append(ordered, summaryRefs{sid: sid, refs: refs})
			}
		}
		for _, item := range ordered {
			p.populateFunctionDefinitionsWithSummary(item.sid, item.refs)
		}
	}

	p.putSummaryForFunctionNode(
		PureFunctionIdentifier_IMPLICIT_PURE_FN,
		p.allocSummary(0, pfiNameImplicitPureCtor))
}

func (p *PureFunctionIdentifier) populateFunctionDefinitionsWithSummary(summaryID pfiSummaryID, references []rhino.Node) {
	invalid := false
	rvaluesAssignedToName := p.scratchCallableLeaves[:0]
	for _, reference := range references {
		if !PureFunctionIdentifier_IsDefinitelyRValue(reference) {
			rvalue := NodeUtil_GetRValueOfLValue(reference)
			if rvalue.IsNil() {
				invalid = true
				break
			}
			prevLen := len(rvaluesAssignedToName)
			if !pureFunctionIdentifier_collectCallableLeavesSlice(rvalue, &rvaluesAssignedToName, p.astAnalyzer) {
				invalid = true
				break
			}
			if p.validateArtificialPurity && !reference.IsFromExterns() {
				for _, callable := range rvaluesAssignedToName[prevLen:] {
					jsdoc := NodeUtil_GetBestJSDocInfo(callable)
					if jsdoc != nil && jsdoc.IsNoSideEffects() {
						p.addArtificiallyPureFunction(callable)
						p.summaries[summaryID].flags |= pfiSummaryCollectDebug
					}
				}
			}
		}
	}
	p.scratchCallableLeaves = rvaluesAssignedToName[:0]

	if len(rvaluesAssignedToName) == 0 || invalid {
		p.setSummaryMutatesGlobalState(summaryID, &p.summaries[summaryID], rhino.NilNode)
		return
	}
	for _, rvalue := range rvaluesAssignedToName {
		if rvalue.IsFunction() {
			p.putSummaryForFunctionNode(rvalue, summaryID)
		} else if NodeUtil_IsUndefined(rvalue) {
			// Assigning `undefined` leaves purity unchanged.
		} else {
			p.addAliasEdge(p.getSummaryForReferenceOrDefault(rvalue), summaryID)
		}
	}
}

var pfiExternFnAnnotationTokens = makeVisitTokenSet(rhino.Token_FUNCTION)

// pfiExternFunctionAnnotationAnalyzer inspects extern function JSDoc annotations
// and applies side effects to the associated summaries.
type pfiExternFunctionAnnotationAnalyzer struct {
	outer *PureFunctionIdentifier
}

var (
	_ nodeTraversal_visitTokens      = (*pfiExternFunctionAnnotationAnalyzer)(nil)
	_ nodeTraversal_noVisitLeaves    = (*pfiExternFunctionAnnotationAnalyzer)(nil)
	_ nodeTraversal_noShouldTraverse = (*pfiExternFunctionAnnotationAnalyzer)(nil)
)

func (*pfiExternFunctionAnnotationAnalyzer) visitTokens() *ast.TokenSet {
	return &pfiExternFnAnnotationTokens
}

func (*pfiExternFunctionAnnotationAnalyzer) noVisitLeaves() {}

func (*pfiExternFunctionAnnotationAnalyzer) nodeTraversalNoShouldTraverse() {}

func (a *pfiExternFunctionAnnotationAnalyzer) ShouldTraverse(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) bool {
	return true
}

func (a *pfiExternFunctionAnnotationAnalyzer) Visit(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) {
	if !node.IsFunction() {
		return
	}
	single, multi := a.outer.getFunctionSummaries(node)
	if single != 0 {
		a.updateSideEffectsForExternFunction(node, single)
	} else {
		for _, sid := range multi {
			a.updateSideEffectsForExternFunction(node, sid)
		}
	}
}

func (a *pfiExternFunctionAnnotationAnalyzer) updateSideEffectsForExternFunction(
	externFunction rhino.Node, sid pfiSummaryID) {
	jutil.CheckArgument(externFunction.IsFunction())
	jutil.CheckArgument(externFunction.IsFromExterns())

	p := a.outer
	s := &p.summaries[sid]
	info := NodeUtil_GetBestJSDocInfo(externFunction)
	if info == nil {
		p.setSummaryMutatesGlobalState(sid, s, externFunction)
		return
	}

	if info.ModifiesThis() {
		p.setSummaryMutatesThis(sid, s, externFunction)
	}
	if info.HasSideEffectsArgumentsAnnotation() {
		p.setSummaryMutatesArguments(sid, s, externFunction)
	}
	if len(info.GetThrowsAnnotations()) != 0 {
		p.setSummaryThrows(sid, s, externFunction)
	}

	if !info.IsNoSideEffects() && s.bitmask == 0 {
		p.setSummaryMutatesGlobalState(sid, s, externFunction)
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
