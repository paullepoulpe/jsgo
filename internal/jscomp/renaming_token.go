// Port of com/google/javascript/jscomp/RenamingToken.java @ closure-compiler v20261021

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
