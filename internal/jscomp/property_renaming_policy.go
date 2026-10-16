// Port of com/google/javascript/jscomp/PropertyRenamingPolicy.java @ closure-compiler v20260415

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
