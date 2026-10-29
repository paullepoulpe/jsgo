// Port of com/google/javascript/jscomp/DiagnosticToSuppressionMapper.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/MapDiagnosticToSuppression.java @ closure-compiler v20261021

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

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/args4j"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
)

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

// CreateSuppressionMap returns the map from diagnostic key to suppression.
//
// Java: DiagnosticToSuppressionMapper.createSuppressionMap()
func (m *DiagnosticToSuppressionMapper) CreateSuppressionMap() map[string]string {
	result := make(map[string]string)

	for suppression, group := range m.diagnosticGroups.All() {
		// "accessControls" and "missingSourcesWarnings" are deprecated in favor of "visibility"
		if suppression == "accessControls" ||
			suppression == "missingSourcesWarnings" ||
			!m.validSuppressions.Contains(suppression) {
			continue
		}

		for _, diagnosticType := range group.GetTypes() {
			key := diagnosticType.key
			// Map.merge(key, value, remappingFunction); DiagnosticType equality is key equality.
			if old, ok := result[key]; ok {
				result[key] = m.PreferMoreTargetedGroup(old, suppression)
			} else {
				result[key] = suppression
			}
		}
	}

	return result
}

// PreferMoreTargetedGroup uses a more targeted suppression if one exists.
//
// Java: DiagnosticToSuppressionMapper.preferMoreTargetedGroup(String, String)
func (m *DiagnosticToSuppressionMapper) PreferMoreTargetedGroup(next string, current string) string {
	nextGroup := m.diagnosticGroups.Get(next)
	currentGroup := m.diagnosticGroups.Get(current)
	if len(nextGroup.GetTypes()) < len(currentGroup.GetTypes()) {
		return next
	}
	return current
}

func diagnosticToSuppressionMapper_sortedKeys(diagnosticToSuppression map[string]string) []string {
	keys := make([]string, 0, len(diagnosticToSuppression))
	for k := range diagnosticToSuppression { // order-independent: sorted immediately below
		keys = append(keys, k)
	}
	slices.SortFunc(keys, jutil.StringCompareTo)
	return keys
}

// Java: DiagnosticToSuppressionMapper.printAsJson(ImmutableSortedMap<String, String>)
func diagnosticToSuppressionMapper_printAsJson(diagnosticToSuppression map[string]string) {
	keys := diagnosticToSuppressionMapper_sortedKeys(diagnosticToSuppression)
	var sb strings.Builder
	sb.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			sb.WriteByte(',')
		}
		kb, _ := json.Marshal(k)
		vb, _ := json.Marshal(diagnosticToSuppression[k])
		sb.Write(kb)
		sb.WriteByte(':')
		sb.Write(vb)
	}
	sb.WriteByte('}')
	jutil.System_Out().Println(sb.String())
}

// Java: DiagnosticToSuppressionMapper.printAsMarkdown(ImmutableSortedMap<String, String>)
func diagnosticToSuppressionMapper_printAsMarkdown(diagnosticToSuppression map[string]string) {
	out := jutil.System_Out()
	out.Println("| Error | Suppression tag |")
	out.Println("|---|---|")

	for _, k := range diagnosticToSuppressionMapper_sortedKeys(diagnosticToSuppression) {
		out.Println("|" + k + "|" + diagnosticToSuppression[k] + "|")
	}
}

// MapDiagnosticToSuppression prints out a map from diagnostic id to
// suppression.
//
// This can be used by tooling to automatically suppress errors given a list
// of diagnostics or to generate documentation.
//
// A few potentially confusing aspects of this class:
//
//   - Some diagnostic ids are not suppressible, and so are not included in
//     the map.
//   - Some diagnostics have multiple possible suppressions. We arbitrarily
//     choose one suppression, with two exceptions: use visibility over
//     accessControls, and use anything over missingSourcesWarnings.
//     TODO(lharker): we could also make this configurable from the command
//     line.
//   - Some suppression groups are not usable in at-suppress annotations, as
//     configured in the Parser, and are not included in the map.
type MapDiagnosticToSuppression struct {
	output DiagnosticToSuppressionMapper_OutputFormat

	help bool
}

// diagnosticToSuppressionMapper_OutputFormatClass is the Class of the enum
// DiagnosticToSuppressionMapper.OutputFormat, the type of the --output
// option.
//
// Go-only: language-adapter — reflection (see package args4j).
var diagnosticToSuppressionMapper_OutputFormatClass = args4j.EnumClassOf(
	"com.google.javascript.jscomp.DiagnosticToSuppressionMapper$OutputFormat",
	DiagnosticToSuppressionMapper_OutputFormat_values())

// Args4jFields lists the fields that carry args4j's @Option annotations, in
// declaration order.
//
// Go-only: language-adapter — annotations (see package args4j).
func (m *MapDiagnosticToSuppression) Args4jFields() []args4j.AnnotatedField {
	const c = "com.google.javascript.jscomp.MapDiagnosticToSuppression"
	return []args4j.AnnotatedField{
		{
			Option: &args4j.Option{
				Name:  "--output",
				Usage: "Sets the desired output format.\n" + "Options: md, json" + "Default: json"},
			Field: args4j.EnumField(c, "output", diagnosticToSuppressionMapper_OutputFormatClass, &m.output),
		},
		{
			Option: &args4j.Option{Name: "--help", Usage: "Usage"},
			Field:  args4j.BooleanField(c, "help", &m.help),
		},
	}
}

// Java: MapDiagnosticToSuppression()
func NewMapDiagnosticToSuppression() *MapDiagnosticToSuppression {
	m := &MapDiagnosticToSuppression{}
	m.output = DiagnosticToSuppressionMapper_OutputFormat_JSON
	m.help = false
	return m
}

// Java: MapDiagnosticToSuppression.main(String[])
func MapDiagnosticToSuppression_Main(args []string) {
	NewMapDiagnosticToSuppression().PrintSuppressions(args)
}

// Java: MapDiagnosticToSuppression.printHelp()
func MapDiagnosticToSuppression_PrintHelp() {
	jutil.System_Out().Println("Usage: MapDiagnosticToSuppression --output={json,md}")
}

// Java: MapDiagnosticToSuppression.printSuppressions(String[])
func (m *MapDiagnosticToSuppression) PrintSuppressions(args []string) {
	parser := args4j.NewCmdLineParser(m)
	if ex := parser.TryParseArgument(args...); ex != nil {
		jutil.System_Err().Println("Could not parse arguments: " + ex.Error())
		MapDiagnosticToSuppression_PrintHelp()
		return
	}

	if m.help {
		MapDiagnosticToSuppression_PrintHelp()
		return
	}

	NewDiagnosticToSuppressionMapper(
		parsing.ParserRunner_GetSuppressionNames(), DiagnosticGroups_GetRegisteredGroups()).
		PrintSuppressions(m.output)
}
