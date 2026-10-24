// Port of com/google/javascript/jscomp/DiagnosticToSuppressionMapper.java @ closure-compiler v20261021

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// DiagnosticToSuppressionMapper is used to print a map from diagnostic id to
// suppression.
type DiagnosticToSuppressionMapper struct {
	validSuppressions *jutil.LinkedHashSet[string]
	diagnosticGroups  *jutil.LinkedHashMap[string, *DiagnosticGroup]
}

// DiagnosticToSuppressionMapper_OutputFormat
//
// Java: DiagnosticToSuppressionMapper.OutputFormat
type DiagnosticToSuppressionMapper_OutputFormat int32

const (
	DiagnosticToSuppressionMapper_OutputFormat_MD DiagnosticToSuppressionMapper_OutputFormat = iota
	DiagnosticToSuppressionMapper_OutputFormat_JSON
)

var diagnosticToSuppressionMapper_OutputFormat_names = [...]string{"MD", "JSON"}

// Java: DiagnosticToSuppressionMapper.OutputFormat.values()
func DiagnosticToSuppressionMapper_OutputFormat_values() []DiagnosticToSuppressionMapper_OutputFormat {
	return []DiagnosticToSuppressionMapper_OutputFormat{
		DiagnosticToSuppressionMapper_OutputFormat_MD, DiagnosticToSuppressionMapper_OutputFormat_JSON,
	}
}

// String is Enum.name().
func (f DiagnosticToSuppressionMapper_OutputFormat) String() string {
	return diagnosticToSuppressionMapper_OutputFormat_names[f]
}

// NewDiagnosticToSuppressionMapper takes validSuppressions, the set of
// diagnostic group names valid to @suppress in code, and diagnosticGroups,
// a map from suppression to DiagnosticGroup.
//
// Java: DiagnosticToSuppressionMapper(Set<String>, ImmutableMap<String, DiagnosticGroup>)
func NewDiagnosticToSuppressionMapper(
	validSuppressions *jutil.LinkedHashSet[string], diagnosticGroups *jutil.LinkedHashMap[string, *DiagnosticGroup]) *DiagnosticToSuppressionMapper {
	return &DiagnosticToSuppressionMapper{validSuppressions: validSuppressions, diagnosticGroups: diagnosticGroups}
}

// PrintSuppressions prints a sorted map of all valid suppressions to
// System.out.
//
// Java: DiagnosticToSuppressionMapper.printSuppressions(OutputFormat)
func (m *DiagnosticToSuppressionMapper) PrintSuppressions(output DiagnosticToSuppressionMapper_OutputFormat) {
	switch output {
	case DiagnosticToSuppressionMapper_OutputFormat_JSON:
		diagnosticToSuppressionMapper_printAsJson(m.CreateSuppressionMap())
	case DiagnosticToSuppressionMapper_OutputFormat_MD:
		diagnosticToSuppressionMapper_printAsMarkdown(m.CreateSuppressionMap())
	}
}

// CreateSuppressionMap returns the map from diagnostic key to suppression,
// sorted by key (ImmutableSortedMap).
//
// Java: DiagnosticToSuppressionMapper.createSuppressionMap()
func (m *DiagnosticToSuppressionMapper) CreateSuppressionMap() *jutil.TreeMap[string, string] {
	validSuppressions := jutil.NewLinkedHashSet(m.validSuppressions.ToSlice()...)
	// deprecated in favor of "visibility"
	validSuppressions.Remove("accessControls")
	validSuppressions.Remove("missingSourcesWarnings")

	diagnosticToSuppression := jutil.NewObjectLinkedHashMap[*DiagnosticType, string]()

	for suppression, group := range m.diagnosticGroups.All() {
		if !validSuppressions.Contains(suppression) {
			continue
		}

		for diagnosticType := range group.GetTypes().All() {
			// Map.merge(key, value, remappingFunction)
			if old, ok := diagnosticToSuppression.Lookup(diagnosticType); ok {
				diagnosticToSuppression.Put(diagnosticType, m.PreferMoreTargetedGroup(old, suppression))
			} else {
				diagnosticToSuppression.Put(diagnosticType, suppression)
			}
		}
	}

	// toImmutableSortedMap(naturalOrder(), entry -> entry.getKey().key, Map.Entry::getValue): the
	// keys are distinct because DiagnosticType equality is key equality.
	result := jutil.NewTreeMap[string, string](jutil.StringCompareTo)
	for diagnosticType, suppression := range diagnosticToSuppression.All() {
		result.Put(diagnosticType.key, suppression)
	}
	return result
}

// PreferMoreTargetedGroup uses a more targeted suppression if one exists.
//
// Java: DiagnosticToSuppressionMapper.preferMoreTargetedGroup(String, String)
func (m *DiagnosticToSuppressionMapper) PreferMoreTargetedGroup(next string, current string) string {
	nextGroup := m.diagnosticGroups.Get(next)
	currentGroup := m.diagnosticGroups.Get(current)
	if nextGroup.GetTypes().Size() < currentGroup.GetTypes().Size() {
		return next
	}
	return current
}

// Java: DiagnosticToSuppressionMapper.printAsJson(ImmutableSortedMap<String, String>)
func diagnosticToSuppressionMapper_printAsJson(diagnosticToSuppression *jutil.TreeMap[string, string]) {
	jutil.System_Out().Println(jutil.NewGson().ToJson(diagnosticToSuppression))
}

// Java: DiagnosticToSuppressionMapper.printAsMarkdown(ImmutableSortedMap<String, String>)
func diagnosticToSuppressionMapper_printAsMarkdown(diagnosticToSuppression *jutil.TreeMap[string, string]) {
	jutil.System_Out().Println("| Error | Suppression tag |")
	jutil.System_Out().Println("|---|---|")

	for k, v := range diagnosticToSuppression.All() {
		jutil.System_Out().Println(jutil.StringFormat("|%s|%s|", k, v))
	}
}
