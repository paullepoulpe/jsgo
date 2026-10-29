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
	"strconv"
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
	jutil.CheckNotNull(supplier)
	lines, err := supplier.TryReadLines()
	if err != nil {
		panic(jutil.NewRuntimeExceptionWithCause(err))
	}
	return deduplicateLines(lines)
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
	return deduplicateLines(jutil.CharStreams_ReadLines(reader))
}

func deduplicateLines(lines []string) []string {
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
	var sb strings.Builder
	if src := error.SourceName(); src != nil {
		sb.WriteString(normalizeSourceNameStr(*src))
	} else {
		sb.WriteString("null")
	}
	sb.WriteString(":")
	if withMetaData {
		sb.WriteString(strconv.Itoa(error.GetLineNumber()))
	}
	firstLine, rest, hasMore := strings.Cut(error.Description(), "\n")
	sb.WriteString("  ")
	sb.WriteString(firstLine)

	// Add the rest of the message as a comment.
	if withMetaData {
		for hasMore {
			var line string
			line, rest, hasMore = strings.Cut(rest, "\n")
			sb.WriteString("\n# ")
			sb.WriteString(line)
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
		normalized := normalizeSourceNameStr(*sourceName)
		return &normalized
	}
	return sourceName
}

// normalizeSourceNameStr strips the first "blaze-out/[^/]*/(bin|genfiles)/" segment.
// e.g. "blaze-out/k8-fastbuild/genfiles/some/path/foo.js" -> "some/path/foo.js"
func normalizeSourceNameStr(s string) string {
	const prefix = "blaze-out/"
	offset := 0
	for {
		idx := strings.Index(s[offset:], prefix)
		if idx < 0 {
			return s
		}
		matchStart := offset + idx
		rest := s[matchStart+len(prefix):]
		slash := strings.IndexByte(rest, '/')
		if slash >= 0 {
			afterSlash := rest[slash:]
			var tail string
			var ok bool
			if tail, ok = strings.CutPrefix(afterSlash, "/bin/"); ok {
			} else if tail, ok = strings.CutPrefix(afterSlash, "/genfiles/"); ok {
			}
			if ok {
				if matchStart == 0 {
					return tail
				}
				return s[:matchStart] + tail
			}
		}
		offset = matchStart + len(prefix)
	}
}

// Java: AllowlistWarningsGuard.getFirstLine(String)
func AllowlistWarningsGuard_GetFirstLine(warning string) string {
	lineLength := strings.IndexByte(warning, '\n')
	if lineLength > 0 {
		warning = warning[:lineLength]
	}
	return warning
}
