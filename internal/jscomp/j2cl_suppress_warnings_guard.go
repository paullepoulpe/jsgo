// Port of com/google/javascript/jscomp/J2clSuppressWarningsGuard.java @ closure-compiler v20261015

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

import "strings"

// J2clSuppressWarningsGuard is a warnings guard that suppresses some
// warnings incompatible with J2CL.
type J2clSuppressWarningsGuard struct {
	WarningsGuardBase
}

// TODO(b/128554878): Cleanup and document all file level suppressions for J2CL generated code.
var J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS *DiagnosticGroup

// Go-only: language-adapter — the static field reads DiagnosticGroups'
// static fields, which are assigned from DiagnosticGroups' init(); this
// file's init() runs after it (Go runs init functions in file name order).
func init() {
	j2clSuppressWarningsGuard_clinit()
}

// Java: J2clSuppressWarningsGuard static initializers
func j2clSuppressWarningsGuard_clinit() {
	name := "j2clIncomaptible"
	J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS =
		NewDiagnosticGroupWithNameFromGroups(
			&name,
			// Do not warn when static overrides do not match. This is safe, because J2CL generated
			// code directly points to declaration and this is also required with collapse properties
			// pass, which does not support dynamic dispatch for static methods.
			DiagnosticGroups_CHECK_STATIC_OVERRIDES,
			// Do not warn on valid Java constructs like "if(false) {...}" and other situations that
			// may arise from transormation of complex Kotlin constructs.
			DiagnosticGroups_CHECK_USELESS_CODE,
			DiagnosticGroups_CONST,
			DiagnosticGroups_EXTRA_REQUIRE,
			// Kotlin allows provably invalid casts so long as it's via a safe cast. This causes J2CL
			// to generate code of the form:
			//   Foo.$isInstance(value) ? /**@type {!Foo}*/ (value) : null
			// However, if value is obviously never a instance of Foo then this will cause invalid
			// cast error. At runtime it's guarded so it would be safe regardless.
			// This also suppresses casts from non-lambda JsFunction implementations to functions.
			// This particular feature is deprecated and will be removed (b/159954752).
			DiagnosticGroups_INVALID_CASTS,
			DiagnosticGroups_LATE_PROVIDE,
			DiagnosticGroups_MISSING_OVERRIDE,
			DiagnosticGroups_MISSING_REQUIRE,
			DiagnosticGroups_STRICT_MODULE_DEP_CHECK,
			DiagnosticGroups_SUSPICIOUS_CODE,
			DiagnosticGroups_UNUSED_LOCAL_VARIABLE,
			// TODO(b/78521031): J2CL targets are not strict missing property compatible.
			DiagnosticGroups_STRICT_MISSING_PROPERTIES,
			DiagnosticGroups_ForName("transitionalSuspiciousCodeWarnings"))
}

// NewJ2clSuppressWarningsGuard is the implicit default constructor.
//
// Java: J2clSuppressWarningsGuard()
func NewJ2clSuppressWarningsGuard() *J2clSuppressWarningsGuard {
	return &J2clSuppressWarningsGuard{}
}

// Java: J2clSuppressWarningsGuard.level(JSError)
func (g *J2clSuppressWarningsGuard) Level(error *JSError) *CheckLevel {
	if !J2clSourceUtils_IsJ2clSource(error.SourceName()) {
		return nil
	}

	// Do not warn about deprecated forwardDeclare usages in J2CL generated code as this particular
	// use-case is still valid, but not elsewhere.
	if DiagnosticGroups_DEPRECATED.Matches(error) &&
		strings.Contains(error.Description(), "forwardDeclare") {
		return CheckLevel_OFF.Ptr()
	}

	if J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS.Matches(error) {
		return CheckLevel_OFF.Ptr()
	}
	return nil
}

// Java: J2clSuppressWarningsGuard.getPriority()
func (g *J2clSuppressWarningsGuard) GetPriority() int {
	return WarningsGuard_Priority_MAX.GetValue()
}
