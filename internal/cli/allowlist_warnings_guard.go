// Port of com/google/javascript/jscomp/AllowlistWarningsGuard.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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

package cli

import (
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// AllowlistWarningsGuard is an extension of WarningsGuard that provides
// functionality to maintain a list of warnings (allowlist). It is
// subclasses' responsibility to decide what to do with the allowlist by
// implementing the level function. Warnings are defined by the name of the JS
// file and the first line of warnings description.
type AllowlistWarningsGuard struct {
	jscomp.WarningsGuardBase

	// The set of allowlisted warnings, same format as formatWarning.
	allowlist map[string]struct{}
}

// AllowlistWarningsGuard_LINE_SPLITTER is Splitter.on('\n'): splitToList.
//
// Go-only: language-adapter — Guava's Splitter.on(char) without options
// splits like strings.Split (empty pieces are kept).
var AllowlistWarningsGuard_LINE_SPLITTER = func(s string) []string { return strings.Split(s, "\n") }

// Pattern to match line number in error descriptions.
var AllowlistWarningsGuard_LINE_NUMBER = jutil.Pattern_Compile(":-?\\d+")

// Java: AllowlistWarningsGuard()
func NewAllowlistWarningsGuard() *AllowlistWarningsGuard {
	return NewAllowlistWarningsGuardWithAllowlist([]string{})
}

// NewAllowlistWarningsGuardWithAllowlist: this class depends on an input set
// that contains the allowlist. The format of each allowlist string is:
// <file-name>:<line-number>? <warning-description> # <optional-comment>
//
// allowlist is the set of JS-warnings that are allowlisted. This is expected
// to have similar format as formatWarning(JSError).
//
// Java: AllowlistWarningsGuard(Set<String>)
func NewAllowlistWarningsGuardWithAllowlist(allowlist []string) *AllowlistWarningsGuard {
	g := &AllowlistWarningsGuard{}
	jutil.CheckNotNull(allowlist)
	g.allowlist = AllowlistWarningsGuard_NormalizeAllowlist(allowlist)
	return g
}

// AllowlistWarningsGuard_NormalizeAllowlist loads legacy warnings list from
// the set of strings. During development line numbers are changed very often
// - we just cut them and compare without ones.
//
// Also remove lines starting with "#" or are blank lines.
//
// Returns known legacy warnings without line numbers.
//
// Java: AllowlistWarningsGuard.normalizeAllowlist(Set<String>)
func AllowlistWarningsGuard_NormalizeAllowlist(allowlist []string) map[string]struct{} {
	result := make(map[string]struct{}, len(allowlist))
	for _, line := range allowlist {
		trimmed := jutil.StringTrim(line)
		if trimmed == "" || trimmed[0] == '#' {
			// strip out empty lines and comments.
			continue
		}

		// Strip line number for matching.
		result[AllowlistWarningsGuard_LINE_NUMBER.Matcher(trimmed).ReplaceFirst(":")] = struct{}{}
	}
	return result
}

// Java: AllowlistWarningsGuard.level(JSError)
func (g *AllowlistWarningsGuard) Level(error *jscomp.JSError) *jscomp.CheckLevel {
	if error.DefaultLevel() == jscomp.CheckLevel_ERROR {
		return nil
	}
	if len(g.allowlist) > 0 && g.ContainWarning(g.FormatWarning(error)) {
		// If the message matches the guard we use WARNING, so that it
		// - Shows up on stderr, and
		// - Gets caught by the AllowlistBuilder downstream in the pipeline
		return jscomp.CheckLevel_WARNING.Ptr()
	}
	return nil
}

// ContainWarning determines whether a given warning is included in the
// allowlist.
//
// formattedWarning is the warning formatted by formattedWarning. Returns
// whether the given warning is allowlisted or not.
//
// Java: AllowlistWarningsGuard.containWarning(String)
func (g *AllowlistWarningsGuard) ContainWarning(formattedWarning string) bool {
	_, ok := g.allowlist[formattedWarning]
	return ok
}

// Java: AllowlistWarningsGuard.getPriority()
func (g *AllowlistWarningsGuard) GetPriority() int {
	return jscomp.WarningsGuard_Priority_SUPPRESS_BY_ALLOWLIST.GetValue()
}

// AllowlistWarningsGuard_FromFile creates a warnings guard from a file.
//
// Java: AllowlistWarningsGuard.fromFile(File)
func AllowlistWarningsGuard_FromFile(file *jutil.File) *AllowlistWarningsGuard {
	return NewAllowlistWarningsGuardWithAllowlist(AllowlistWarningsGuard_LoadAllowlistedJsWarnings(file))
}

// AllowlistWarningsGuard_LoadAllowlistedJsWarnings loads legacy warnings
// list from the file.
//
// Returns the lines of the file.
//
// Java: AllowlistWarningsGuard.loadAllowlistedJsWarnings(File)
func AllowlistWarningsGuard_LoadAllowlistedJsWarnings(file *jutil.File) []string {
	return AllowlistWarningsGuard_LoadAllowlistedJsWarningsFromCharSource(
		jutil.Files_AsCharSource(file, jutil.StandardCharsets_UTF_8))
}

// AllowlistWarningsGuard_LoadAllowlistedJsWarningsFromCharSource loads
// legacy warnings list from the file.
//
// Returns the lines of the file.
//
// Java: AllowlistWarningsGuard.loadAllowlistedJsWarnings(CharSource)
func AllowlistWarningsGuard_LoadAllowlistedJsWarningsFromCharSource(supplier *jutil.CharSource) []string {
	var result []string
	jutil.Try(func() {
		result = AllowlistWarningsGuard_LoadAllowlistedJsWarningsFromReader(supplier.OpenStream())
	}, func(e jutil.IIOException) {
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})
	return result
}

// AllowlistWarningsGuard_LoadAllowlistedJsWarningsFromReader loads legacy
// warnings list from the file. It panics with an IOException when reading
// fails.
//
// Returns the lines of the file.
//
// Java: AllowlistWarningsGuard.loadAllowlistedJsWarnings(Reader)
func AllowlistWarningsGuard_LoadAllowlistedJsWarningsFromReader(reader jutil.JavaReader) []string {
	// TODO(nicksantos): This is a weird API.
	jutil.CheckNotNull(reader)
	lines := jutil.CharStreams_ReadLines(reader)
	seen := make(map[string]struct{}, len(lines))
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if _, ok := seen[line]; !ok {
			seen[line] = struct{}{}
			result = append(result, line)
		}
	}
	return result
}

// FormatWarning: if subclasses want to modify the formatting, they should
// override #formatWarning(JSError, boolean), not this method.
//
// Java: AllowlistWarningsGuard.formatWarning(JSError)
func (g *AllowlistWarningsGuard) FormatWarning(error *jscomp.JSError) string {
	return g.FormatWarningWithMetaData(error, false)
}

// FormatWarningWithMetaData: if withMetaData is true, include metadata
// that's useful to humans This metadata won't be used for matching the
// warning.
//
// Java: AllowlistWarningsGuard.formatWarning(JSError, boolean)
func (g *AllowlistWarningsGuard) FormatWarningWithMetaData(error *jscomp.JSError, withMetaData bool) string {
	sb := &strings.Builder{}
	sb.WriteString(jutil.ToJavaString(g.NormalizeSourceName(error.SourceName())))
	sb.WriteString(":")
	if withMetaData {
		sb.WriteString(jutil.ToJavaString(error.GetLineNumber()))
	}
	lines := AllowlistWarningsGuard_LINE_SPLITTER(error.Description())
	sb.WriteString("  ")
	sb.WriteString(lines[0])

	// Add the rest of the message as a comment.
	if withMetaData {
		for i := 1; i < len(lines); i++ {
			sb.WriteString("\n# ")
			sb.WriteString(lines[i])
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// NormalizeSourceName: sourceName may be nil (and so is the result then).
//
// Java: AllowlistWarningsGuard.normalizeSourceName(String)
func (g *AllowlistWarningsGuard) NormalizeSourceName(sourceName *string) *string {
	if sourceName != nil {
		// e.g.
		// "blaze-out/k8-fastbuild/genfiles/some/path/foo.js" -> "some/path/foo.js"
		normalized := jutil.StringReplaceFirst(*sourceName, "blaze-out/[^/]*/(bin|genfiles)/", "")
		return &normalized
	}
	return sourceName
}

// Java: AllowlistWarningsGuard.getFirstLine(String)
func AllowlistWarningsGuard_GetFirstLine(warning string) string {
	lineLength := strings.IndexByte(warning, '\n')
	if lineLength > 0 {
		warning = warning[:lineLength]
	}
	return warning
}

// AllowlistWarningsGuard_AllowlistBuilder is the allowlist builder.
type AllowlistWarningsGuard_AllowlistBuilder struct {
	outer *AllowlistWarningsGuard

	warnings []*jscomp.JSError
	// May be nil.
	productName *string
	// May be nil.
	generatorTarget *string
	// May be nil.
	headerNote *string
}

// Java: AllowlistWarningsGuard.AllowlistBuilder()
func NewAllowlistWarningsGuard_AllowlistBuilder(outer *AllowlistWarningsGuard) *AllowlistWarningsGuard_AllowlistBuilder {
	return &AllowlistWarningsGuard_AllowlistBuilder{
		outer: outer,
	}
}

// SetProductName: fill in your product name to get a fun message!
//
// Java: AllowlistWarningsGuard.AllowlistBuilder.setProductName(String)
func (b *AllowlistWarningsGuard_AllowlistBuilder) SetProductName(name *string) *AllowlistWarningsGuard_AllowlistBuilder {
	b.productName = name
	return b
}

// SetGeneratorTarget: fill in instructions on how to generate this
// allowlist.
//
// Java: AllowlistWarningsGuard.AllowlistBuilder.setGeneratorTarget(String)
func (b *AllowlistWarningsGuard_AllowlistBuilder) SetGeneratorTarget(name *string) *AllowlistWarningsGuard_AllowlistBuilder {
	b.generatorTarget = name
	return b
}

// SetNote: a note to include at the top of the allowlist file.
//
// Java: AllowlistWarningsGuard.AllowlistBuilder.setNote(String)
func (b *AllowlistWarningsGuard_AllowlistBuilder) SetNote(note *string) *AllowlistWarningsGuard_AllowlistBuilder {
	b.headerNote = note
	return b
}

// Java: AllowlistWarningsGuard.AllowlistBuilder.report(CheckLevel, JSError)
func (b *AllowlistWarningsGuard_AllowlistBuilder) Report(level *jscomp.CheckLevel, error *jscomp.JSError) {
	if error.DefaultLevel() == jscomp.CheckLevel_ERROR {
		// ERROR-level diagnostics are ignored by AllowlistWarningsGuard (c.f. above getLevel).
		return
	}
	b.warnings = append(b.warnings, error)
}

// WriteAllowlist writes the warnings collected in a format that the
// AllowlistWarningsGuard can read back later. It panics with an IOException
// (FileNotFoundException) when the file cannot be written.
//
// Java: AllowlistWarningsGuard.AllowlistBuilder.writeAllowlist(File)
func (b *AllowlistWarningsGuard_AllowlistBuilder) WriteAllowlist(out *jutil.File) {
	stream := jutil.NewPrintStreamFromFile(out)
	defer stream.Close()
	b.AppendAllowlist(stream)
}

// AppendAllowlist writes the warnings collected in a format that the
// AllowlistWarningsGuard can read back later.
//
// Java: AllowlistWarningsGuard.AllowlistBuilder.appendAllowlist(PrintStream)
func (b *AllowlistWarningsGuard_AllowlistBuilder) AppendAllowlist(out *jutil.PrintStream) {
	out.Append("# This is a list of legacy warnings that have yet to be fixed.\n")

	if b.productName != nil && *b.productName != "" && len(b.warnings) > 0 {
		out.Append(
			"# Please find some time and fix at least one of them " +
				"and it will be the happiest day for " +
				*b.productName +
				".\n")
	}

	if b.generatorTarget != nil && *b.generatorTarget != "" {
		out.Append("# When you fix any of these warnings, run " + *b.generatorTarget + " task.\n")
	}

	if b.headerNote != nil {
		out.Append("#" + strings.Join(strings.Split(*b.headerNote, "\n"), "\n# ") + "\n")
	}

	type allowlistTypeGroup struct {
		diagType *jscomp.DiagnosticType
		seen     map[string]struct{}
		warnings []string
	}
	groupsByKey := make(map[string]*allowlistTypeGroup)
	var groups []*allowlistTypeGroup
	for _, warning := range b.warnings {
		t := warning.Type()
		key := t.Key()
		grp, ok := groupsByKey[key]
		if !ok {
			grp = &allowlistTypeGroup{
				diagType: t,
				seen:     make(map[string]struct{}),
			}
			groupsByKey[key] = grp
			groups = append(groups, grp)
		}
		formatted := b.outer.FormatWarningWithMetaData(warning, true /* withLineNumber */)
		if _, exists := grp.seen[formatted]; !exists {
			grp.seen[formatted] = struct{}{}
			grp.warnings = append(grp.warnings, formatted)
		}
	}

	slices.SortFunc(groups, func(a, c *allowlistTypeGroup) int {
		return a.diagType.CompareTo(c.diagType)
	})

	for _, grp := range groups {
		type_ := grp.diagType
		if jscomp.DiagnosticGroups_DEPRECATED.MatchesType(type_) {
			// Deprecation warnings are not raisable to error, so we don't need them in allowlists.
			continue
		}
		out.Append("\n# Warning ").
			Append(type_.Key()).
			Append(": ").
			Println(AllowlistWarningsGuard_LINE_SPLITTER(type_.FormatPattern())[0])

		slices.SortFunc(grp.warnings, jutil.StringCompareTo)
		for _, warning := range grp.warnings {
			out.Println(warning)
		}
	}
	out.Flush()
}
