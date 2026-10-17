// Port of com/google/javascript/jscomp/modules/ResolveExportResult.java @ closure-compiler v20260415

/*
 * Copyright 2019 The Closure Compiler Authors.
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
)

// ResolveExportResult is the result of resolving an export, which can be a
// valid binding, ambiguous, not found, or an error.
type ResolveExportResult struct {
	binding *Binding
	state   resolveExportResult_State
}

// Java: ResolveExportResult.State
type resolveExportResult_State int32

const (
	resolveExportResult_State_RESOLVED resolveExportResult_State = iota
	resolveExportResult_State_AMBIGUOUS
	resolveExportResult_State_NOT_FOUND
	resolveExportResult_State_ERROR
)

// Java: ResolveExportResult(Binding, State)
func newResolveExportResult(binding *Binding, state resolveExportResult_State) *ResolveExportResult {
	return &ResolveExportResult{binding: binding, state: state}
}

// Copy creates a new result that has the given node for the source of the
// binding and given type of binding.
//
// Java: ResolveExportResult.copy(Node, Binding.CreatedBy)
func (r *ResolveExportResult) Copy(sourceNode *rhino.Node, createdBy Binding_CreatedBy) *ResolveExportResult {
	jutil.CheckNotNull(sourceNode)

	if r.binding == nil {
		return r
	}

	return newResolveExportResult(r.binding.Copy(sourceNode, createdBy), r.state)
}

// HadError is true if there was an error resolving the export, false
// otherwise.
//
// Java: ResolveExportResult.hadError()
func (r *ResolveExportResult) HadError() bool {
	return r.state == resolveExportResult_State_ERROR
}

// IsAmbiguous is true if the export is ambiguous, false otherwise.
//
// Java: ResolveExportResult.isAmbiguous()
func (r *ResolveExportResult) IsAmbiguous() bool {
	return r.state == resolveExportResult_State_AMBIGUOUS
}

// Resolved is true if the export was successfully resolved, false
// otherwise.
//
// Java: ResolveExportResult.resolved()
func (r *ResolveExportResult) Resolved() bool {
	return r.state == resolveExportResult_State_RESOLVED
}

// Found is true if the export key exists on the given module, even if it is
// ambiguous or had an error.
//
// Java: ResolveExportResult.found()
func (r *ResolveExportResult) Found() bool {
	return r.state != resolveExportResult_State_NOT_FOUND
}

// GetBinding may return nil.
//
// Java: ResolveExportResult.getBinding()
func (r *ResolveExportResult) GetBinding() *Binding {
	return r.binding
}

// ResolveExportResult_AMBIGUOUS: the result of resolving the export was
// ambiguous.
//
// This happens when there are multiple export * from statements that end up
// causing the same key to be re-exported.
//
// When resolving an import or transitive export, if the result is
// ambiguous, an error should be reported at the import / transitive export
// site, an then ERROR returned so that more ambiguous errors are not
// reported.
var ResolveExportResult_AMBIGUOUS = newResolveExportResult(nil, resolveExportResult_State_AMBIGUOUS)

// ResolveExportResult_NOT_FOUND: the export was not found because the
// module never exported the key.
//
// When resolving an import or transitive export, if the result is not
// found, an error should be reported at the import / transitive export
// site, an then ERROR returned so that more ambiguous errors are not
// reported.
var ResolveExportResult_NOT_FOUND = newResolveExportResult(nil, resolveExportResult_State_NOT_FOUND)

// ResolveExportResult_ERROR: there was an error resolving the export.
//
// This can mean that:
//
//  1. When resolving a transitive export, the transitive export was not
//     found.
//  2. When resolving a transitive export, the transitive export was
//     ambiguous.
//  3. There was a cycle resolving an export.
//  4. The requested module does not exist.
//
// When resolving an import or transitive export, if the result is ERROR,
// then resolving should also return ERROR. No error needs to be reported,
// this is an indication that something has already been reported.
var ResolveExportResult_ERROR = newResolveExportResult(nil, resolveExportResult_State_ERROR)

// Java: ResolveExportResult.of(Binding)
func ResolveExportResult_Of(binding *Binding) *ResolveExportResult {
	jutil.CheckNotNull(binding)
	return newResolveExportResult(binding, resolveExportResult_State_RESOLVED)
}

// String is Object.toString().
func (r *ResolveExportResult) String() string {
	return jutil.ObjectToString("com.google.javascript.jscomp.modules.ResolveExportResult", r)
}
