// Port of com/google/javascript/jscomp/LinterMain.java @ closure-compiler v20260415

/*
 * Copyright 2015 The Closure Compiler Authors.
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

import "github.com/paullepoulpe/jsgo/internal/args4j"

// LinterMain is a minimal binary that just runs the "lint" checks which can
// be run on a single file at a time. This means some checks in the
// lintChecks DiagnosticGroup are skipped, since they depend on type
// information.
//
// Java: final class LinterMain
type LinterMain struct {
	fix bool

	resolutionMode ModuleLoader_ResolutionMode

	files []string
}

// linterMain_ResolutionModeClass is the Class of the enum
// ModuleLoader.ResolutionMode, the type of the --resolution_mode option.
//
// Go-only: language-adapter — reflection (see package args4j).
var linterMain_ResolutionModeClass = args4j.EnumClassOf(
	"com.google.javascript.jscomp.deps.ModuleLoader$ResolutionMode",
	ModuleLoader_ResolutionMode_values())

// Java: LinterMain()
func newLinterMain() *LinterMain {
	m := &LinterMain{}
	m.fix = false
	m.resolutionMode = ModuleLoader_ResolutionMode_BROWSER
	m.files = []string{}
	return m
}

// Args4jFields lists the fields that carry args4j's @Option and @Argument
// annotations, in declaration order.
//
// Go-only: language-adapter — annotations (see package args4j).
func (m *LinterMain) Args4jFields() []args4j.AnnotatedField {
	const c = "com.google.javascript.jscomp.LinterMain"
	return []args4j.AnnotatedField{
		{
			Option: &args4j.Option{Name: "--fix", Usage: "Fix lint warnings automatically"},
			Field:  args4j.BooleanField(c, "fix", &m.fix),
		},
		{
			Option: &args4j.Option{
				Name:  "--resolution_mode",
				Usage: "Expected path resolution mode for resolving module specifiers."},
			Field: args4j.EnumField(c, "resolutionMode", linterMain_ResolutionModeClass, &m.resolutionMode),
		},
		{
			Argument: &args4j.Argument{},
			Field:    args4j.ListField(c, "files", args4j.Class_String, &m.files),
		},
	}
}

// LinterMain_Main throws IOException and CmdLineException.
//
// Java: LinterMain.main(String[])
func LinterMain_Main(args []string) {
	newLinterMain().run(args)
}

// Java: LinterMain.run(String[])
func (m *LinterMain) run(args []string) {
	parser := args4j.NewCmdLineParser(m)
	parser.ParseArgument(args...)

	linter := Linter_Builder_().WithModuleResolutionMode(m.resolutionMode).Build()

	for _, filename := range m.files {
		if m.fix {
			linter.FixRepeatedly(filename)
		} else {
			linter.Lint(filename)
		}
	}
}
