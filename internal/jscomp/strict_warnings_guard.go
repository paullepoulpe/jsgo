// Port of com/google/javascript/jscomp/StrictWarningsGuard.java @ closure-compiler v20260415

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// StrictWarningsGuard: all warnings should be reported as errors.
type StrictWarningsGuard struct {
	WarningsGuardBase
}

var _ WarningsGuard = (*StrictWarningsGuard)(nil)

const strictWarningsGuard_serialVersionUID int64 = 1

var StrictWarningsGuard_UNRAISABLE_WARNING = DiagnosticType_Warning("JSC_UNRAISABLE_WARNING", "{0}")

// Java: StrictWarningsGuard()
func NewStrictWarningsGuard() *StrictWarningsGuard {
	return &StrictWarningsGuard{}
}

// Level may return nil.
//
// Java: StrictWarningsGuard.level(JSError)
func (g *StrictWarningsGuard) Level(error *JSError) *CheckLevel {
	if error.Type() == StrictWarningsGuard_UNRAISABLE_WARNING {
		return nil
	}
	if error.DefaultLevel().IsOn() {
		return CheckLevel_ERROR.Ptr()
	}
	return nil
}

// Java: StrictWarningsGuard.getPriority()
func (g *StrictWarningsGuard) GetPriority() int {
	return WarningsGuard_Priority_STRICT.value() // applied last
}

// String is Object.toString().
func (g *StrictWarningsGuard) String() string {
	return jutil.ObjectToString("com.google.javascript.jscomp.StrictWarningsGuard", g)
}
