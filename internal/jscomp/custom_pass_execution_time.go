// Port of com/google/javascript/jscomp/CustomPassExecutionTime.java @ closure-compiler v20261021

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

// CustomPassExecutionTime is the custom pass type. Controls the compilation
// stage at which each custom pass executes.
type CustomPassExecutionTime int32

const (
	CustomPassExecutionTime_BEFORE_CHECKS CustomPassExecutionTime = iota
	CustomPassExecutionTime_BEFORE_OPTIMIZATIONS
	CustomPassExecutionTime_BEFORE_OPTIMIZATION_LOOP
	CustomPassExecutionTime_AFTER_OPTIMIZATION_LOOP
)

var customPassExecutionTime_names = [...]string{
	"BEFORE_CHECKS", "BEFORE_OPTIMIZATIONS", "BEFORE_OPTIMIZATION_LOOP", "AFTER_OPTIMIZATION_LOOP",
}

// Name is Enum.name().
func (t CustomPassExecutionTime) Name() string { return customPassExecutionTime_names[t] }

// String is Enum.toString().
func (t CustomPassExecutionTime) String() string { return t.Name() }

// Ordinal is Enum.ordinal().
func (t CustomPassExecutionTime) Ordinal() int { return int(t) }
