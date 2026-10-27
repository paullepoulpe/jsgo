// Port of com/google/javascript/jscomp/RenamingMap.java @ closure-compiler v20261021

/*
 * Copyright 2013 The Closure Compiler Authors.
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

// RenamingMap is the interface used by to substitute names.
type RenamingMap interface {
	// Get returns the replacement for value, or nil (Java null).
	//
	// Java: RenamingMap.get(String)
	Get(value string) *string
}

// RenamingToken is a marker object for the command line api to pass
// special values to the ReplaceIdGenerators pass.
type RenamingToken int32

const (
	RenamingToken_INCONSISTENT RenamingToken = iota
	RenamingToken_STABLE
	RenamingToken_DISABLE
)

var renamingToken_names = [...]string{"INCONSISTENT", "STABLE", "DISABLE"}

// Name is Enum.name().
func (r RenamingToken) Name() string { return renamingToken_names[r] }

// String is Enum.toString().
func (r RenamingToken) String() string { return r.Name() }

// Ordinal is Enum.ordinal().
func (r RenamingToken) Ordinal() int { return int(r) }

// Java: RenamingToken.get(String)
func (r RenamingToken) Get(value string) *string {
	return nil
}

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

// PropertyRenamingPolicy holds the policies to determine how properties
// should be renamed.
type PropertyRenamingPolicy int32

const (
	// PropertyRenamingPolicy_OFF renames no properties.
	PropertyRenamingPolicy_OFF PropertyRenamingPolicy = iota

	// PropertyRenamingPolicy_ALL_UNQUOTED renames all properties that
	// aren't explicitly quoted and aren't externally defined (i.e. declared
	// in an externs file). This policy achieves better compaction than the
	// others.
	//
	// See RenameProperties.
	PropertyRenamingPolicy_ALL_UNQUOTED
)

var propertyRenamingPolicy_names = [...]string{"OFF", "ALL_UNQUOTED"}

// PropertyRenamingPolicy_values is PropertyRenamingPolicy.values().
func PropertyRenamingPolicy_values() []PropertyRenamingPolicy {
	return []PropertyRenamingPolicy{PropertyRenamingPolicy_OFF, PropertyRenamingPolicy_ALL_UNQUOTED}
}

// Name is Enum.name().
func (p PropertyRenamingPolicy) Name() string { return propertyRenamingPolicy_names[p] }

// String is Enum.toString().
func (p PropertyRenamingPolicy) String() string { return p.Name() }

// Ordinal is Enum.ordinal().
func (p PropertyRenamingPolicy) Ordinal() int { return int(p) }
