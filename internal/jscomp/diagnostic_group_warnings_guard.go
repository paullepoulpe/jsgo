// Port of com/google/javascript/jscomp/DiagnosticGroupWarningsGuard.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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

import "github.com/paullepoulpe/jsgo/internal/base"

// DiagnosticGroupWarningsGuard sets the level for a particular
// DiagnosticGroup.
type DiagnosticGroupWarningsGuard struct {
	WarningsGuardBase
	group *DiagnosticGroup
	level CheckLevel
}

var _ WarningsGuard = (*DiagnosticGroupWarningsGuard)(nil)

const diagnosticGroupWarningsGuard_serialVersionUID int64 = 1

// Java: DiagnosticGroupWarningsGuard(DiagnosticGroup, CheckLevel)
func NewDiagnosticGroupWarningsGuard(
	group *DiagnosticGroup, level CheckLevel) *DiagnosticGroupWarningsGuard {
	return &DiagnosticGroupWarningsGuard{group: group, level: level}
}

// Level may return nil.
//
// Java: DiagnosticGroupWarningsGuard.level(JSError)
func (g *DiagnosticGroupWarningsGuard) Level(error *JSError) *CheckLevel {
	if g.group.Matches(error) {
		return g.level.Ptr()
	}
	return nil
}

// Java: DiagnosticGroupWarningsGuard.mustRunChecks(DiagnosticGroup)
func (g *DiagnosticGroupWarningsGuard) MustRunChecks(otherGroup *DiagnosticGroup) base.Tri {
	if g.level.IsOn() {
		for _, t := range otherGroup.GetTypes() {
			if g.group.MatchesType(t) {
				return base.Tri_TRUE
			}
		}
		return base.Tri_UNKNOWN
	}
	for _, t := range otherGroup.GetTypes() {
		if !g.group.MatchesType(t) {
			return base.Tri_UNKNOWN
		}
	}
	return base.Tri_FALSE
}

// Java: DiagnosticGroupWarningsGuard.toString()
func (g *DiagnosticGroupWarningsGuard) String() string {
	return g.group.String() + "(" + g.level.String() + ")"
}
