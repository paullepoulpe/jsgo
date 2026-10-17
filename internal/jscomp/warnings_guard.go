// Port of com/google/javascript/jscomp/WarningsGuard.java @ closure-compiler v20260415

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

// WarningsGuard allows to flexibly manage what to do with a reported
// warning/error.
//
// Guard has several choices: - return OFF - suppress the warning/error -
// return WARNING - return ERROR report it with high severity - return null.
// Does not know what to do with it. Lets the other guard decide what to do
// with it.
//
// Although the interface is very simple, it allows you easily customize what
// warnings you are interested in.
//
// For example there are could be several implementations: StrictGuard -
// {return ERROR}. All warnings should be treat as errors. SilentGuard - {if
// (WARNING) return OFF}. Suppress all warnings but still fail if JS has
// errors. AllowlistGuard (if !allowlistErrors.contains(error) return ERROR)
// return error if it does not present in the allowlist.
type WarningsGuard interface {
	// Level returns a new check level for a given error.
	//
	// nil means that this guard does not know what to do with the error. nil
	// can be used it chain multiple guards; if current guard returns nil,
	// then the next in the chain should process it.
	//
	// Java: WarningsGuard.level(JSError)
	Level(error *JSError) *CheckLevel

	// MustRunChecks: do checks for group still need to be run if this guard
	// is installed?
	//
	//  1. TRUE: Enables one or more types in the group, so it must be checked.
	//  2. FALSE: Disables all types in the group, so it need not be checked.
	//  3. UNKNOWN: Does not affect or only partially disables the group, so
	//     checking is undecided.
	//
	// Java: WarningsGuard.mustRunChecks(DiagnosticGroup)
	MustRunChecks(group *DiagnosticGroup) base.Tri

	// GetPriority is the priority in which warnings guards are applied.
	// Lower means the guard will be applied sooner. Expressed on a scale of
	// 1 to 100.
	//
	// Java: WarningsGuard.getPriority()
	GetPriority() int
}

// WarningsGuardBase holds the base implementations of WarningsGuard;
// subclasses embed it.
type WarningsGuardBase struct {
	// Go-only: language-adapter — Java compares guards by identity
	// (ComposeWarningsGuard keys them in a map and orders them in a
	// TreeSet), but pointers to distinct zero-size Go values may be equal,
	// so the base keeps every guard (e.g. StrictWarningsGuard) non-zero-size.
	_ byte
}

const warningsGuard_serialVersionUID int64 = 1

// Java: WarningsGuard.mustRunChecks(DiagnosticGroup)
func (WarningsGuardBase) MustRunChecks(group *DiagnosticGroup) base.Tri {
	return base.Tri_UNKNOWN
}

// Java: WarningsGuard.getPriority()
func (WarningsGuardBase) GetPriority() int {
	return WarningsGuard_Priority_DEFAULT.value()
}

// WarningsGuard_Priority is the priority of a guard.
//
// Java: WarningsGuard.Priority
type WarningsGuard_Priority int32

const (
	WarningsGuard_Priority_MAX WarningsGuard_Priority = iota
	WarningsGuard_Priority_MIN
	WarningsGuard_Priority_STRICT
	WarningsGuard_Priority_DEFAULT
	WarningsGuard_Priority_SUPPRESS_BY_ALLOWLIST
	WarningsGuard_Priority_SUPPRESS_DOC
	WarningsGuard_Priority_FILTER_BY_PATH
)

var warningsGuard_Priority_names = [...]string{
	"MAX", "MIN", "STRICT", "DEFAULT", "SUPPRESS_BY_ALLOWLIST", "SUPPRESS_DOC", "FILTER_BY_PATH",
}

var warningsGuard_Priority_values = [...]int{1, 100, 100, 50, 40, 20, 1}

// Java: WarningsGuard.Priority.values()
func WarningsGuard_Priority_values() []WarningsGuard_Priority {
	return []WarningsGuard_Priority{
		WarningsGuard_Priority_MAX, WarningsGuard_Priority_MIN, WarningsGuard_Priority_STRICT,
		WarningsGuard_Priority_DEFAULT, WarningsGuard_Priority_SUPPRESS_BY_ALLOWLIST,
		WarningsGuard_Priority_SUPPRESS_DOC, WarningsGuard_Priority_FILTER_BY_PATH,
	}
}

// String is Enum.name().
func (p WarningsGuard_Priority) String() string { return warningsGuard_Priority_names[p] }

// Ordinal is Enum.ordinal().
func (p WarningsGuard_Priority) Ordinal() int { return int(p) }

// value is the field Priority.value.
func (p WarningsGuard_Priority) value() int { return warningsGuard_Priority_values[p] }

// Java: WarningsGuard.Priority.getValue()
func (p WarningsGuard_Priority) GetValue() int { return p.value() }
