// Port of com/google/javascript/jscomp/ClosureUnawareCodeWarningsGuard.java @ closure-compiler v20261015

/*
 * Copyright 2024 The Closure Compiler Authors.
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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ClosureUnawareCodeWarningsGuard is a warnings guard that suppresses
// warnings that are spurious for code that is unaware of Closure compiler's
// requirements.
type ClosureUnawareCodeWarningsGuard struct {
	WarningsGuardBase
}

var _ WarningsGuard = (*ClosureUnawareCodeWarningsGuard)(nil)

// Java: ClosureUnawareCodeWarningsGuard()
func NewClosureUnawareCodeWarningsGuard() *ClosureUnawareCodeWarningsGuard {
	return &ClosureUnawareCodeWarningsGuard{}
}

var closureUnawareCodeWarningsGuard_DEFAULT_CLOSURE_UNAWARE_CODE_SUPPRESSIONS = NewDiagnosticGroupWithName(
	closureUnawareCodeWarningsGuard_name("closureUnawareCodeJSDocIncompatible"),
	RhinoErrorReporter_TYPE_PARSE_ERROR,
	RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR,
	RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING,
	RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING,
	RhinoErrorReporter_UNNECESSARY_ESCAPE,
	RhinoErrorReporter_BAD_JSDOC_ANNOTATION,
	RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES,
	RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR)

func closureUnawareCodeWarningsGuard_name(s string) *string { return &s }

// Level may return nil.
//
// Java: ClosureUnawareCodeWarningsGuard.level(JSError)
func (g *ClosureUnawareCodeWarningsGuard) Level(error *JSError) *CheckLevel {
	if !closureUnawareCodeWarningsGuard_DEFAULT_CLOSURE_UNAWARE_CODE_SUPPRESSIONS.Matches(error) {
		return nil
	}
	node := error.Node()
	if node.IsNil() {
		return nil
	}
	if node.IsClosureUnawareCode() {
		return CheckLevel_OFF.Ptr()
	}
	return nil
}

// Java: ClosureUnawareCodeWarningsGuard.getPriority()
func (g *ClosureUnawareCodeWarningsGuard) GetPriority() int {
	return WarningsGuard_Priority_MAX.GetValue()
}

// String is Object.toString().
func (g *ClosureUnawareCodeWarningsGuard) String() string {
	return jutil.ObjectToString("com.google.javascript.jscomp.ClosureUnawareCodeWarningsGuard", g)
}
