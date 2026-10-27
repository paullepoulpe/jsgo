// Port of com/google/javascript/jscomp/ConformanceAllowlister.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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

// ConformanceAllowlister creates or updates conformance allowlist/whitelist
// entries.
//
// Java: class ConformanceAllowlister
type ConformanceAllowlister struct {
	// Go-only: language-adapter — distinct identity.
	_ byte
}

// ConformanceAllowlister_GetViolatingPaths returns, in encounter order, the
// source names of the conformance violations of requirement. Like
// toImmutableSet, it throws NullPointerException for a violation without a
// source name.
//
// Java: ConformanceAllowlister.getViolatingPaths(Compiler, Node, Node, Requirement)
func ConformanceAllowlister_GetViolatingPaths(
	compiler *Compiler, externs rhino.Node, ast rhino.Node, requirement *Requirement) *jutil.LinkedHashSet[string] {
	result := jutil.NewLinkedHashSet[string]()
	for _, e := range ConformanceAllowlister_GetConformanceErrors(compiler, externs, ast, requirement) {
		result.Add(*jutil.CheckNotNull(e.SourceName()))
	}
	return result
}

// ConformanceAllowlister_GetViolatingNodes returns, in encounter order, the
// nodes of the conformance violations of requirement. Like toImmutableSet,
// it throws NullPointerException for a violation without a node.
//
// Java: ConformanceAllowlister.getViolatingNodes(Compiler, Node, Node, Requirement)
func ConformanceAllowlister_GetViolatingNodes(
	compiler *Compiler, externs rhino.Node, ast rhino.Node, requirement *Requirement) *jutil.LinkedHashSet[rhino.Node] {
	result := jutil.NewLinkedHashSet[rhino.Node]()
	for _, e := range ConformanceAllowlister_GetConformanceErrors(compiler, externs, ast, requirement) {
		result.Add(rhino.CheckNotNull(e.Node()))
	}
	return result
}

// Java: ConformanceAllowlister.getConformanceErrors(Compiler, Node, Node, Requirement)
func ConformanceAllowlister_GetConformanceErrors(
	compiler *Compiler, externs rhino.Node, ast rhino.Node, requirement *Requirement) []*JSError {
	cleanedRequirement :=
		requirement.ToBuilder().
			ClearWhitelist().
			ClearAllowlist().
			ClearWhitelistRegexp().
			ClearAllowlistRegexp().
			ClearWhitelistEntry().
			ClearAllowlistEntry().
			SetSeverity(Requirement_Severity_ERROR).
			Build() // So we only have one type of error.
	cleanedConfig :=
		ConformanceConfig_NewBuilder().AddRequirement(cleanedRequirement).Build()

	oldErrorManager := compiler.GetErrorManager()
	var errors []*JSError
	func() {
		defer compiler.SetErrorManager(oldErrorManager)
		// TODO(bangert): handle invalid conformance requirements
		compiler.SetErrorManager(
			&conformanceAllowlister_ErrorManager{
				ThreadSafeDelegatingErrorManager: NewThreadSafeDelegatingErrorManager(oldErrorManager),
				errors:                           &errors,
			})
		check :=
			NewCheckConformanceWithConfigs(
				compiler, []*ConformanceConfig{cleanedConfig} /* reportingMode= */, nil)
		check.Process(externs, ast)
	}()
	return errors
}

// conformanceAllowlister_ErrorManager is the anonymous
// ThreadSafeDelegatingErrorManager of getConformanceErrors: it collects the
// conformance violations instead of reporting them. Its synchronized
// method holds the object's monitor, the embedded manager's.
type conformanceAllowlister_ErrorManager struct {
	*ThreadSafeDelegatingErrorManager
	errors *[]*JSError
}

// Java: ThreadSafeDelegatingErrorManager.shouldReportConformanceViolation(Requirement, Optional<RequirementScopeEntry>, JSError, LibraryLevelNonAllowlistedConformanceViolationsBehavior, boolean)
func (m *conformanceAllowlister_ErrorManager) ShouldReportConformanceViolation(
	requirement *Requirement,
	whitelistEntry *RequirementScopeEntry,
	diagnostic *JSError,
	behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior,
	isAllowlisted bool) bool {
	m.monitor.Lock()
	defer m.monitor.Unlock()
	*m.errors = append(*m.errors, diagnostic)
	return false
}
