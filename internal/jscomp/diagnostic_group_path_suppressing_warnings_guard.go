// Port of com/google/javascript/jscomp/DiagnosticGroupPathSuppressingWarningsGuard.java @ closure-compiler v20261021

/*
 * Copyright 2016 The Closure Compiler Authors.
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

import "strings"

// DiagnosticGroupPathSuppressingWarningsGuard is a warnings guard that
// suppresses warnings for a particular diagnostic group for a file that
// contains the specified substring.
type DiagnosticGroupPathSuppressingWarningsGuard struct {
	WarningsGuardBase
	group *DiagnosticGroup
	part  string
}

var _ WarningsGuard = (*DiagnosticGroupPathSuppressingWarningsGuard)(nil)

// Java: DiagnosticGroupPathSuppressingWarningsGuard(DiagnosticGroup, String)
func NewDiagnosticGroupPathSuppressingWarningsGuard(group *DiagnosticGroup, part string) *DiagnosticGroupPathSuppressingWarningsGuard {
	return &DiagnosticGroupPathSuppressingWarningsGuard{group: group, part: part}
}

// Level does not touch warnings in other paths. It may return nil.
//
// Java: DiagnosticGroupPathSuppressingWarningsGuard.level(JSError)
func (g *DiagnosticGroupPathSuppressingWarningsGuard) Level(error *JSError) *CheckLevel {
	if error.SourceName() == nil || !strings.Contains(*error.SourceName(), g.part) {
		return nil
	}
	if g.group.Matches(error) {
		return CheckLevel_OFF.Ptr()
	}
	return nil
}

// Java: DiagnosticGroupPathSuppressingWarningsGuard.toString()
func (g *DiagnosticGroupPathSuppressingWarningsGuard) String() string {
	return g.group.String() + "(" + g.part + ")"
}
