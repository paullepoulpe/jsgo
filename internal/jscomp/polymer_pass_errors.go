// Port of com/google/javascript/jscomp/PolymerPassErrors.java @ closure-compiler v20260415

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

// PolymerPassErrors holds the static error constants related to the
// PolymerPass.
type PolymerPassErrors struct{}

func newPolymerPassErrors() *PolymerPassErrors { return &PolymerPassErrors{} }

var (
	// TODO(jlklein): Switch back to an error when everyone is upgraded to
	// Polymer 1.0
	PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID = DiagnosticType_Warning(
		"JSC_POLYMER_DESCRIPTOR_NOT_VALID",
		"The argument to Polymer() is not an obj lit or the Polymer 2 class does not have a"+
			" static getter named 'config'. Ignoring this definition.")

	PolymerPassErrors_POLYMER_INVALID_BEHAVIOR = DiagnosticType_Error(
		"JSC_POLYMER_INVALID_BEHAVIOR", "A Polymer behavior may not include an ''is'' property.")

	PolymerPassErrors_POLYMER_MISSING_IS = DiagnosticType_Error("JSC_POLYMER_MISSING_IS",
		"The class descriptor must include an ''is'' property.")

	PolymerPassErrors_POLYMER_UNEXPECTED_PARAMS = DiagnosticType_Error(
		"JSC_POLYMER_UNEXPECTED_PARAMS", "The class definition has too many arguments.")

	PolymerPassErrors_POLYMER_MISSING_EXTERNS = DiagnosticType_Error(
		"JSC_POLYMER_MISSING_EXTERNS", "Missing Polymer externs.")

	PolymerPassErrors_POLYMER_INVALID_PROPERTY = DiagnosticType_Error(
		"JSC_POLYMER_INVALID_PROPERTY", "Polymer property has an invalid or missing type.")

	PolymerPassErrors_POLYMER_INVALID_EXTENDS = DiagnosticType_Error(
		"JSC_POLYMER_INVALID_EXTENDS",
		"Cannot extend HTML element ''{0}''. The element is probably either misspelled,"+
			" or needs to be added to the list of known elements.")

	PolymerPassErrors_POLYMER_INVALID_BEHAVIOR_ARRAY = DiagnosticType_Error(
		"JSC_POLYMER_INVALID_BEHAVIOR_ARRAY", "The behaviors property must be an array literal.")

	PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR = DiagnosticType_Error(
		"JSC_POLYMER_UNQUALIFIED_BEHAVIOR",
		"Behaviors must be global names or qualified names that are declared as object literals or "+
			"array literals of other valid Behaviors.")

	PolymerPassErrors_POLYMER_UNANNOTATED_BEHAVIOR = DiagnosticType_Error(
		"JSC_POLYMER_UNANNOTATED_BEHAVIOR",
		"Behavior declarations must be annotated with @polymerBehavior.")

	PolymerPassErrors_POLYMER_CLASS_PROPERTIES_INVALID = DiagnosticType_Error(
		"JSC_POLYMER_CLASS_PROPERTIES_INVALID",
		"The Polymer element class 'properties' getter does not return an object literal. "+
			"Ignoring this definition.")

	PolymerPassErrors_POLYMER_CLASS_PROPERTIES_NOT_STATIC = DiagnosticType_Error(
		"JSC_POLYMER_CLASS_PROPERTIES_NOT_STATIC",
		"The Polymer element class 'properties' getter is not declared static. "+
			"Ignoring this definition.")

	PolymerPassErrors_POLYMER_CLASS_UNNAMED = DiagnosticType_Warning(
		"JSC_POLYMER2_UNNAMED",
		"Unable to locate a valid name for the Polymer element class."+
			"Ignoring this definition.")

	PolymerPassErrors_POLYMER_MISPLACED_PROPERTY_JSDOC = DiagnosticType_Warning(
		"JSC_POLYMER_MISPLACED_PROPERTY_JSDOC",
		"When a Polymer property is declared in the constructor, its JSDoc "+
			"should only be in the constructor, not on the Polymer properties configuration.")

	PolymerPassErrors_POLYMER_UNPARSABLE_STRING = DiagnosticType_Error(
		"JSC_POLYMER_UNPARSABLE_STRING",
		"The Polymer computed property or complex observer string could not be parsed.")
)
