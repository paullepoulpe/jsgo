// Port of com/google/javascript/jscomp/MapDiagnosticToSuppression.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/args4j"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
)

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
	parsed := false
	jutil.Try(func() {
		parser.ParseArgument(args...)
		parsed = true
	}, func(ex args4j.ICmdLineException) {
		jutil.System_Err().Println("Could not parse arguments: " + ex.Error())
		MapDiagnosticToSuppression_PrintHelp()
	})
	if !parsed {
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
