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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// pfiSummaryID indexes into PureFunctionIdentifier.summaries.
// Index 0 is reserved for the conservative <unknown> summary.
type pfiSummaryID uint32

const (
	pfiUnknownSummaryID pfiSummaryID = 0
	pfiMultiSummaryBit  uint32       = 1 << 31
)

// Side-effect bitmask flags on pfiSummary.bitmask.
const (
	pfiThrows             uint8 = 1 << 0
	pfiMutatesGlobalState uint8 = 1 << 1
	pfiMutatesThis        uint8 = 1 << 2
	pfiMutatesArguments   uint8 = 1 << 3
	pfiAllFlags           uint8 = pfiThrows | pfiMutatesGlobalState | pfiMutatesThis | pfiMutatesArguments
)

// Auxiliary flags and name kind on pfiSummary.flags.
const (
	pfiSummaryCollectDebug uint8 = 1 << 0
	pfiSummaryInQueue      uint8 = 1 << 1
	pfiNameKindShift       uint8 = 2
)

type pfiNameKind uint8

const (
	pfiNameVar pfiNameKind = iota
	pfiNameProp
	pfiNameUnknown
	pfiNameImplicitPureCtor
	pfiNameAnonymous
)

// pfiSummary is a pointer-free summary for the set of functions sharing a name.
type pfiSummary struct {
	mutatedParamMask uint64
	nameID           uint32
	firstCallerEdge  int32
	bitmask          uint8
	flags            uint8
}

func (s *pfiSummary) nameKind() pfiNameKind {
	return pfiNameKind(s.flags >> pfiNameKindShift)
}

func (s *pfiSummary) displayName() string {
	switch s.nameKind() {
	case pfiNameVar:
		return strpool.String(strpool.ID(s.nameID))
	case pfiNameProp:
		return "." + strpool.String(strpool.ID(s.nameID))
	case pfiNameUnknown:
		return "<unknown>"
	case pfiNameImplicitPureCtor:
		return "<implicit pure class ctor"
	case pfiNameAnonymous:
		return "<anonymous>"
	default:
		return ""
	}
}

// Flags on pfiEdge.flags.
const (
	pfiEdgeAlias              uint8 = 1 << 0
	pfiEdgePropagateThrows    uint8 = 1 << 1
	pfiEdgeAlwaysNonLocalArgs uint8 = 1 << 2
	pfiEdgeThisCaller         uint8 = 1 << 3
	pfiEdgeThisGlobal         uint8 = 1 << 4
)

// pfiEdge is a pointer-free reverse-call-graph edge from callee summary to caller summary.
type pfiEdge struct {
	nonLocalArgMask uint64
	caller          pfiSummaryID
	nextCallerEdge  int32
	callerFn        rhino.Node
	flags           uint8
}

// Flags on pfiCallSite.flags.
const (
	pfiCallIsCallOrTagged   uint8 = 1 << 0
	pfiCallViaCallOrApply   uint8 = 1 << 1
	pfiCallIsSuper          uint8 = 1 << 2
	pfiCallKnownPureForMark uint8 = 1 << 3
	pfiCallSkipAnalyzer     uint8 = 1 << 4
	pfiCallIsCacheCall      uint8 = 1 << 5
)

// pfiCallSite records an invocation node and its resolved callee summary ID(s)
// during the single AST traversal so MarkPureFunctionCalls needs no AST/name re-lookup.
type pfiCallSite struct {
	callNode    rhino.Node
	firstCallee uint32 // pfiSummaryID when calleeCount == 1; offset into multiCallees when calleeCount > 1
	calleeCount uint32
	flags       uint8
}

func (p *PureFunctionIdentifier) allocSummary(nameID uint32, kind pfiNameKind) pfiSummaryID {
	id := pfiSummaryID(len(p.summaries))
	p.summaries = append(p.summaries, pfiSummary{
		nameID:          nameID,
		firstCallerEdge: -1,
		flags:           uint8(kind) << pfiNameKindShift,
	})
	return id
}

func (p *PureFunctionIdentifier) putSummaryForFunctionNode(fn rhino.Node, summary pfiSummaryID) {
	cur := p.fnSummaries[fn]
	if cur == 0 {
		p.fnSummaries[fn] = uint32(summary)
		return
	}
	if cur&pfiMultiSummaryBit == 0 {
		idx := uint32(len(p.fnMultiSummaries))
		p.fnMultiSummaries = append(p.fnMultiSummaries, []pfiSummaryID{pfiSummaryID(cur), summary})
		p.fnSummaries[fn] = pfiMultiSummaryBit | idx
		return
	}
	idx := cur &^ pfiMultiSummaryBit
	p.fnMultiSummaries[idx] = append(p.fnMultiSummaries[idx], summary)
}

func (p *PureFunctionIdentifier) getFunctionSummaries(fn rhino.Node) (single pfiSummaryID, multi []pfiSummaryID) {
	cur := p.fnSummaries[fn]
	if cur == 0 {
		return 0, nil
	}
	if cur&pfiMultiSummaryBit == 0 {
		return pfiSummaryID(cur), nil
	}
	return 0, p.fnMultiSummaries[cur&^pfiMultiSummaryBit]
}

func (p *PureFunctionIdentifier) getSummaryForReference(nameRef rhino.Node) (pfiSummaryID, bool) {
	switch nameRef.GetToken() {
	case rhino.Token_NAME:
		s := p.nameSummaries[nameRef.GetStringID()]
		return s, s != 0
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		s := p.propSummaries[nameRef.GetStringID()]
		return s, s != 0
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected name reference: " + nameRef.String()))
	}
}

func (p *PureFunctionIdentifier) getSummaryForReferenceOrDefault(nameRef rhino.Node) pfiSummaryID {
	if s, ok := p.getSummaryForReference(nameRef); ok {
		return s
	}
	return pfiUnknownSummaryID
}

func (p *PureFunctionIdentifier) addArtificiallyPureFunction(fn rhino.Node) {
	if p.artificiallyPureSet == nil {
		p.artificiallyPureSet = make(map[rhino.Node]struct{}, 8)
	}
	if _, exists := p.artificiallyPureSet[fn]; exists {
		return
	}
	p.artificiallyPureSet[fn] = struct{}{}
	p.artificiallyPureFns = append(p.artificiallyPureFns, fn)
}

func (p *PureFunctionIdentifier) maybeCollectDebugReason(id pfiSummaryID, s *pfiSummary, reason rhino.Node) {
	if s.flags&pfiSummaryCollectDebug == 0 || reason.IsNil() {
		return
	}
	fn := reason
	if !fn.IsFunction() {
		fn = NodeUtil_GetEnclosingFunction(reason)
	}
	if p.debugReasons == nil {
		p.debugReasons = make(map[pfiSummaryID]*jutil.LinkedHashSet[rhino.Node], 8)
	}
	set := p.debugReasons[id]
	if set == nil {
		set = jutil.NewLinkedHashSet[rhino.Node]()
		p.debugReasons[id] = set
	}
	set.Add(fn)
}

func (p *PureFunctionIdentifier) setSummaryThrows(id pfiSummaryID, s *pfiSummary, reason rhino.Node) {
	p.maybeCollectDebugReason(id, s, reason)
	s.bitmask |= pfiThrows
}

func (p *PureFunctionIdentifier) setSummaryMutatesThis(id pfiSummaryID, s *pfiSummary, reason rhino.Node) {
	p.maybeCollectDebugReason(id, s, reason)
	s.bitmask |= pfiMutatesThis
}

func (p *PureFunctionIdentifier) setSummaryMutatesArguments(id pfiSummaryID, s *pfiSummary, reason rhino.Node) {
	p.maybeCollectDebugReason(id, s, reason)
	s.mutatedParamMask = ^uint64(0)
	s.bitmask |= pfiMutatesArguments
}

func (p *PureFunctionIdentifier) setSummaryMutatesArgumentAt(id pfiSummaryID, s *pfiSummary, paramIndex int, reason rhino.Node) {
	jutil.CheckArgumentMsg1(paramIndex >= 0, "%s", paramIndex)
	p.maybeCollectDebugReason(id, s, reason)
	if paramIndex < 64 {
		s.mutatedParamMask |= uint64(1) << paramIndex
	} else {
		s.mutatedParamMask = ^uint64(0)
	}
	s.bitmask |= pfiMutatesArguments
}

func (p *PureFunctionIdentifier) setSummaryMutatesGlobalState(id pfiSummaryID, s *pfiSummary, reason rhino.Node) {
	p.maybeCollectDebugReason(id, s, reason)
	s.mutatedParamMask = ^uint64(0)
	s.bitmask |= pfiAllFlags
}

func (p *PureFunctionIdentifier) addEdge(callee pfiSummaryID, caller pfiSummaryID, callerFn rhino.Node, flags uint8, nonLocalArgMask uint64) int32 {
	idx := int32(len(p.edges))
	p.edges = append(p.edges, pfiEdge{
		nonLocalArgMask: nonLocalArgMask,
		caller:          caller,
		nextCallerEdge:  p.summaries[callee].firstCallerEdge,
		callerFn:        callerFn,
		flags:           flags,
	})
	p.summaries[callee].firstCallerEdge = idx
	return idx
}

func (p *PureFunctionIdentifier) addAliasEdge(callee pfiSummaryID, caller pfiSummaryID) {
	p.addEdge(callee, caller, rhino.NilNode, pfiEdgeAlias|pfiEdgePropagateThrows, 0)
}

// propagateSideEffects propagates side-effect flags from impure callees to
// callers using an impure-seeded FIFO worklist over the flat edge arena.
//
// Java: PureFunctionIdentifier.propagateSideEffects()
func (p *PureFunctionIdentifier) propagateSideEffects() {
	// Reverse caller edge lists so traversal order matches insertion order.
	for i := range p.summaries {
		prev := int32(-1)
		for cur := p.summaries[i].firstCallerEdge; cur >= 0; {
			next := p.edges[cur].nextCallerEdge
			p.edges[cur].nextCallerEdge = prev
			prev = cur
			cur = next
		}
		p.summaries[i].firstCallerEdge = prev
	}

	queue := make([]pfiSummaryID, 0, 64)
	for i := range p.summaries {
		if p.summaries[i].bitmask != 0 && p.summaries[i].firstCallerEdge >= 0 {
			p.summaries[i].flags |= pfiSummaryInQueue
			queue = append(queue, pfiSummaryID(i))
		}
	}

	for head := 0; head < len(queue); head++ {
		if head > 1024 && head*2 >= len(queue) {
			n := copy(queue, queue[head:])
			queue = queue[:n]
			head = 0
		}
		calleeID := queue[head]
		p.summaries[calleeID].flags &^= pfiSummaryInQueue
		calleeMask := p.summaries[calleeID].bitmask
		calleeParams := p.summaries[calleeID].mutatedParamMask

		for eIdx := p.summaries[calleeID].firstCallerEdge; eIdx >= 0; eIdx = p.edges[eIdx].nextCallerEdge {
			e := &p.edges[eIdx]
			if p.propagateEdge(calleeMask, calleeParams, e) {
				caller := &p.summaries[e.caller]
				if caller.firstCallerEdge >= 0 && caller.flags&pfiSummaryInQueue == 0 {
					caller.flags |= pfiSummaryInQueue
					queue = append(queue, e.caller)
				}
			}
		}
	}
}

func (p *PureFunctionIdentifier) propagateEdge(calleeMask uint8, calleeParams uint64, e *pfiEdge) bool {
	caller := &p.summaries[e.caller]
	initialMask := caller.bitmask
	initialParams := caller.mutatedParamMask

	if e.flags&pfiEdgeAlias != 0 {
		caller.bitmask |= calleeMask
		caller.mutatedParamMask |= calleeParams
	} else if calleeMask&pfiMutatesGlobalState != 0 {
		p.setSummaryMutatesGlobalState(e.caller, caller, e.callerFn)
	} else {
		if e.flags&pfiEdgePropagateThrows != 0 && calleeMask&pfiThrows != 0 {
			p.setSummaryThrows(e.caller, caller, e.callerFn)
		}
		if calleeMask&pfiMutatesArguments != 0 &&
			(e.flags&pfiEdgeAlwaysNonLocalArgs != 0 || (calleeParams&e.nonLocalArgMask) != 0) {
			p.setSummaryMutatesGlobalState(e.caller, caller, e.callerFn)
		}
		if calleeMask&pfiMutatesThis != 0 {
			if e.flags&pfiEdgeThisCaller != 0 {
				p.setSummaryMutatesThis(e.caller, caller, e.callerFn)
			} else if e.flags&pfiEdgeThisGlobal != 0 {
				p.setSummaryMutatesGlobalState(e.caller, caller, e.callerFn)
			}
		}
	}

	return caller.bitmask != initialMask || caller.mutatedParamMask != initialParams
}

// markPureFunctionCalls stamps side-effect flags onto recorded invocation nodes.
//
// Java: PureFunctionIdentifier.markPureFunctionCalls()
func (p *PureFunctionIdentifier) markPureFunctionCalls() {
	for i := range p.callSites {
		site := p.callSites[i]
		flags := rhino.NewNode_SideEffectFlags()

		if site.flags&pfiCallKnownPureForMark != 0 {
			flags.ClearAllFlags()
		} else {
			if site.calleeCount == 0 {
				flags.SetAllFlags()
			} else {
				var unionMask uint8
				if site.calleeCount == 1 {
					unionMask = p.summaries[site.firstCallee].bitmask
				} else {
					for _, cid := range p.multiCallees[site.firstCallee : site.firstCallee+site.calleeCount] {
						unionMask |= p.summaries[cid].bitmask
					}
				}
				flags.ClearAllFlags()
				if unionMask&pfiMutatesGlobalState != 0 {
					flags.SetMutatesGlobalState()
				}
				if unionMask&pfiMutatesArguments != 0 {
					flags.SetMutatesArguments()
				}
				if unionMask&pfiThrows != 0 {
					flags.SetThrows()
				}
				if site.flags&pfiCallIsCallOrTagged != 0 && unionMask&pfiMutatesThis != 0 {
					if site.flags&pfiCallViaCallOrApply != 0 {
						flags.SetMutatesArguments() // `this` is passed as an argument.
					} else {
						flags.SetMutatesThis()
					}
				}
			}

			if site.flags&pfiCallIsSuper != 0 {
				// All `super()` calls implicitly mutate `this` in the caller scope.
				flags.SetMutatesThis()
			}
		}

		if site.callNode.GetSideEffectFlags() != flags.ValueOf() {
			site.callNode.SetSideEffectFlagsObject(flags)
			p.compiler.ReportChangeToEnclosingScope(site.callNode)
		}
	}
}
