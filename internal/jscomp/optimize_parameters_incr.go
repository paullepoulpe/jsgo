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

import "github.com/paullepoulpe/jsgo/internal/rhino"

type optimizeParametersCandidate struct {
	id     uint32
	isProp bool
	dirty  bool
	refs   []rhino.Node
}

func (o *OptimizeParameters) hasGlobalVarByID(id uint32) bool {
	if o.refMap != nil {
		return o.refMap.HasGlobalVarByID(id)
	}
	return o.globalScope != nil && o.globalScope.GetVarByID(id) != nil
}

func (o *OptimizeParameters) doesFnRefOwnArguments(fn rhino.Node) bool {
	if o.refMap != nil && o.refMap.index != nil {
		return o.refMap.index.doesFunctionRefOwnArguments(fn)
	}
	return NodeUtil_DoesFunctionReferenceOwnArgumentsObject(fn)
}

func (o *OptimizeParameters) changeStamp() int32 {
	if o.compiler != nil && o.compiler.changeTracker != nil {
		return o.compiler.changeTracker.GetChangeStamp()
	}
	return 0
}

// optimizeParameters_hasExtraOrTrailingUndefinedArg reports whether any normal,
// optional-chain, or `new` call target in refs passes more than maxFormalsCount
// arguments or has a trailing `undefined` / `void 0` argument.
func optimizeParameters_hasExtraOrTrailingUndefinedArg(refs []rhino.Node, maxFormalsCount int) bool {
	for _, n := range refs {
		if !OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
			continue
		}
		arg := OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
		if arg.IsNil() {
			continue
		}
		count := 0
		var last rhino.Node
		for ; !arg.IsNil(); arg = arg.GetNext() {
			count++
			if count > maxFormalsCount {
				return true
			}
			last = arg
		}
		if NodeUtil_IsUndefined(last) {
			return true
		}
	}
	return false
}

func (o *OptimizeParameters) runWithReferenceMap(refMap *OptimizeCalls_ReferenceMap) {
	ix := refMap.index
	logging := o.decisionsLog.IsLogging()
	useClean := !logging && ix != nil

	var toOptimize []optimizeParametersCandidate

	refMap.forEachNameByID(func(id uint32, refs []rhino.Node, slotDirty bool) {
		if useClean && ix.oc.canSkipClean(false, slotDirty, id, refs, ocCleanParams) {
			return
		}
		o.convertedTTL = false
		if o.analyzeCandidateByID("name", id, refs) {
			if o.convertedTTL && useClean {
				ix.oc.syncLiveMutatedScopes(o.compiler)
			}
			toOptimize = append(toOptimize, optimizeParametersCandidate{
				id:     id,
				isProp: false,
				dirty:  o.convertedTTL,
				refs:   refs,
			})
		} else if useClean {
			ix.oc.markClean(false, id, refs, ocCleanParams)
		}
	})

	refMap.forEachPropByID(func(id uint32, refs []rhino.Node, slotDirty bool) {
		if useClean && ix.oc.canSkipClean(true, slotDirty, id, refs, ocCleanParams) {
			return
		}
		o.convertedTTL = false
		if o.analyzeCandidateByID("property", id, refs) {
			if o.convertedTTL && useClean {
				ix.oc.syncLiveMutatedScopes(o.compiler)
			}
			toOptimize = append(toOptimize, optimizeParametersCandidate{
				id:     id,
				isProp: true,
				dirty:  o.convertedTTL,
				refs:   refs,
			})
		} else if useClean {
			ix.oc.markClean(true, id, refs, ocCleanParams)
		}
	})

	// NOTE: Optimizations must keep ReferenceMap consistent:
	// 1. Optional-arg elimination only moves parameter names/defaults into bodies.
	// 2. Constant-arg elimination only runs on single-definition functions, moving
	//    arguments from call sites into the body.
	// 3. Unused-arg elimination may delete references and therefore runs last.
	for i := range toOptimize {
		s0 := o.changeStamp()
		o.tryEliminateOptionalArgs(toOptimize[i].refs)
		if o.changeStamp() != s0 {
			toOptimize[i].dirty = true
		}
	}

	for i := range toOptimize {
		s0 := o.changeStamp()
		o.tryEliminateConstantArgs(toOptimize[i].refs)
		if o.changeStamp() != s0 {
			toOptimize[i].dirty = true
		}
	}

	optimizer := optimizeParametersUnusedOptimizer{outer: o}
	for i := range toOptimize {
		s0 := o.changeStamp()
		rem0 := len(optimizer.toRemove)
		rep0 := len(optimizer.toReplaceWithZero)
		optimizer.tryEliminateUnusedArgs(toOptimize[i].refs)
		if o.changeStamp() != s0 || len(optimizer.toRemove) != rem0 || len(optimizer.toReplaceWithZero) != rep0 {
			toOptimize[i].dirty = true
		}
	}
	optimizer.applyChanges()

	if useClean && len(toOptimize) > 0 {
		ix.oc.syncLiveMutatedScopes(o.compiler)
		for i := range toOptimize {
			cand := &toOptimize[i]
			if !cand.dirty && !ix.oc.refsTouchLiveMutatedScopes(cand.refs) {
				ix.oc.markClean(cand.isProp, cand.id, cand.refs, ocCleanParams)
			}
		}
	}
}
