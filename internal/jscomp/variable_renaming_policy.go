// Port of com/google/javascript/jscomp/VariableRenamingPolicy.java @ closure-compiler v20261015

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

// VariableRenamingPolicy holds the policies to determine which variables
// should be renamed.
type VariableRenamingPolicy int32

const (
	// VariableRenamingPolicy_OFF renames no variables.
	VariableRenamingPolicy_OFF VariableRenamingPolicy = iota

	// VariableRenamingPolicy_LOCAL renames local variables only.
	VariableRenamingPolicy_LOCAL

	// VariableRenamingPolicy_ALL renames all variables and functions unless
	// they are exported or externed.
	VariableRenamingPolicy_ALL
)

var variableRenamingPolicy_names = [...]string{"OFF", "LOCAL", "ALL"}

// VariableRenamingPolicy_values is VariableRenamingPolicy.values().
func VariableRenamingPolicy_values() []VariableRenamingPolicy {
	return []VariableRenamingPolicy{VariableRenamingPolicy_OFF, VariableRenamingPolicy_LOCAL, VariableRenamingPolicy_ALL}
}

// Name is Enum.name().
func (p VariableRenamingPolicy) Name() string { return variableRenamingPolicy_names[p] }

// String is Enum.toString().
func (p VariableRenamingPolicy) String() string { return p.Name() }

// Ordinal is Enum.ordinal().
func (p VariableRenamingPolicy) Ordinal() int { return int(p) }
