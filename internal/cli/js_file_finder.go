// Port of com/google/javascript/jscomp/CommandLineRunner.java @ closure-compiler v20261021

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

package cli

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CommandLineRunner_FindJsFiles returns all the JavaScript files from the
// set of patterns. The patterns support globs, such as '*.js' for all JS
// files in a directory and '**.js' for all JS files within the directory and
// sub-directories.
//
// Java: CommandLineRunner.findJsFiles(Collection<String>)
func CommandLineRunner_FindJsFiles(patterns []string) []string {
	return commandLineRunner_findJsFilesSorted(patterns, false)
}

type jsInputEntry struct {
	orig     string
	orderIdx int
}

// jsInputsCollector tracks normalized-absolute-path -> original-path entries
// while preserving insertion order (or sorting keys alphabetically on output).
type jsInputsCollector struct {
	byAbs map[string]jsInputEntry
	order []string
}

func (c *jsInputsCollector) put(absPath, origPath string) {
	if c.byAbs == nil {
		c.byAbs = make(map[string]jsInputEntry)
	}
	if e, exists := c.byAbs[absPath]; exists {
		e.orig = origPath
		c.byAbs[absPath] = e
		return
	}
	idx := len(c.order)
	c.order = append(c.order, absPath)
	c.byAbs[absPath] = jsInputEntry{orig: origPath, orderIdx: idx}
}

func (c *jsInputsCollector) remove(absPath string) {
	e, exists := c.byAbs[absPath]
	if !exists {
		return
	}
	delete(c.byAbs, absPath)
	c.order[e.orderIdx] = ""
}

// commandLineRunner_findJsFilesSorted returns all the JavaScript files from
// the set of patterns; sortAlphabetically is whether the output filenames
// should be in alphabetical order. It returns the list of JS filenames found
// by expanding the patterns.
//
// Java: CommandLineRunner.findJsFiles(Collection<String>, boolean)
func commandLineRunner_findJsFilesSorted(patterns []string, sortAlphabetically bool) []string {
	if len(patterns) == 0 {
		return nil
	}
	// A map from normalized absolute paths to original paths. We need to return original paths to
	// support allowlist files that depend on them.
	var allJsInputs jsInputsCollector
	var excludes map[string]struct{}
	var wd *jutil.Path
	for _, pattern := range patterns {
		if !strings.Contains(pattern, "*") && !strings.HasPrefix(pattern, "!") {
			if jutil.File_IsDirectoryPath(pattern) {
				commandLineRunner_matchPaths(jutil.NewFileWithParent(jutil.NewFile(pattern), "**.js").String(), &allJsInputs, &excludes, &wd)
			} else {
				original := jutil.Path_Of(pattern)
				norm := original.Normalize()
				if !norm.IsAbsolute() && wd == nil {
					wd = jutil.Path_WorkingDir()
				}
				pathStringAbsolute := norm.ToAbsolutePathWithWd(wd).String()
				if _, excluded := excludes[pathStringAbsolute]; !excluded {
					allJsInputs.put(pathStringAbsolute, original.String())
				}
			}
		} else {
			commandLineRunner_matchPaths(pattern, &allJsInputs, &excludes, &wd)
		}
	}

	if len(allJsInputs.byAbs) == 0 {
		return nil
	}
	if len(allJsInputs.order) != len(allJsInputs.byAbs) {
		allJsInputs.order = slices.DeleteFunc(allJsInputs.order, func(s string) bool { return s == "" })
	}
	if sortAlphabetically {
		slices.SortFunc(allJsInputs.order, jutil.StringCompareTo)
	}
	result := make([]string, len(allJsInputs.order))
	for i, absPath := range allJsInputs.order {
		result[i] = allJsInputs.byAbs[absPath].orig
	}
	return result
}

// Java: CommandLineRunner.matchPaths(String, Map<String, String>, Set<String>)
func commandLineRunner_matchPaths(
	pattern string,
	allJsInputs *jsInputsCollector,
	excludes *map[string]struct{},
	wd **jutil.Path,
) {
	remove := strings.HasPrefix(pattern, "!")
	if remove {
		pattern = pattern[1:]
	}

	separator := jutil.File_separator

	// Split the pattern into two pieces: the globbing part
	// and the non-globbing prefix.
	// Go-only: language-adapter — Guava Splitter.on / Joiner.on with the
	// ASCII separator are strings.Split / strings.Join.
	patternParts := strings.Split(pattern, separator)
	prefix := "."
	for i := 0; i < len(patternParts); i++ {
		if strings.Contains(patternParts[i], "*") {
			if i > 0 {
				prefix = strings.Join(patternParts[:i], separator)
				pattern = strings.Join(patternParts[i:], separator)
			}
			break
		}
	}

	matcher := compileGlobMatcher(prefix + separator + pattern)
	visitFile := func(p *jutil.Path) {
		var norm *jutil.Path
		matched := matcher.MatchString(p.String())
		if !matched {
			norm = p.Normalize()
			matched = matcher.MatchString(norm.String())
		}
		if !matched {
			return
		}
		if norm == nil {
			norm = p.Normalize()
		}
		if !norm.IsAbsolute() && *wd == nil {
			*wd = jutil.Path_WorkingDir()
		}
		pathStringAbsolute := norm.ToAbsolutePathWithWd(*wd).String()
		if remove {
			if *excludes == nil {
				*excludes = make(map[string]struct{})
			}
			(*excludes)[pathStringAbsolute] = struct{}{}
			allJsInputs.remove(pathStringAbsolute)
		} else if _, excluded := (*excludes)[pathStringAbsolute]; !excluded {
			allJsInputs.put(pathStringAbsolute, p.String())
		}
	}

	startPath := jutil.Path_Of(prefix)
	startStr := startPath.String()
	st, err := os.Stat(startStr)
	if err != nil {
		// Broken symlink: visit the link itself as a file; nonexistent start: skip.
		if lst, lerr := os.Lstat(startStr); lerr == nil && !lst.IsDir() {
			visitFile(startPath)
		}
		return
	}
	if !st.IsDir() {
		visitFile(startPath)
		return
	}
	evalRoot, err := filepath.EvalSymlinks(startStr)
	if err != nil {
		return
	}
	commandLineRunner_walkDir(evalRoot, startPath, nil, visitFile)
}

func commandLineRunner_walkDir(
	evalRoot string,
	logicalRoot *jutil.Path,
	ancestors []os.FileInfo,
	visitFile func(*jutil.Path),
) {
	_ = filepath.WalkDir(evalRoot, func(walkPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || walkPath == evalRoot {
			return nil
		}
		var rel string
		switch evalRoot {
		case ".":
			rel = walkPath
		case "/":
			rel = walkPath[1:]
		default:
			rel = walkPath[len(evalRoot)+1:]
		}
		if d.IsDir() {
			return nil
		}
		logicalPath := logicalRoot.ResolveString(rel)
		if d.Type()&fs.ModeSymlink != 0 {
			targetInfo, err := os.Stat(walkPath)
			if err == nil && targetInfo.IsDir() {
				for _, anc := range ancestors {
					if os.SameFile(anc, targetInfo) {
						return nil
					}
				}
				nextAncestors := slices.Clone(ancestors)
				for cur := filepath.Dir(walkPath); ; {
					dirInfo, derr := os.Stat(cur)
					if derr != nil || os.SameFile(dirInfo, targetInfo) {
						return nil
					}
					nextAncestors = append(nextAncestors, dirInfo)
					if cur == evalRoot {
						break
					}
					parent := filepath.Dir(cur)
					if parent == cur {
						break
					}
					cur = parent
				}
				evalTarget, err := filepath.EvalSymlinks(walkPath)
				if err != nil {
					return nil
				}
				commandLineRunner_walkDir(evalTarget, logicalPath, nextAncestors, visitFile)
				return nil
			}
		}
		visitFile(logicalPath)
		return nil
	})
}

func compileGlobMatcher(glob string) *regexp.Regexp {
	var sb strings.Builder
	sb.WriteString(`(?s)^`)
	groupOpen := false
	n := len(glob)
	for i := 0; i < n; i++ {
		switch glob[i] {
		case '\\':
			if i+1 == n {
				panic(jutil.NewPatternSyntaxException("No character to escape", glob, i))
			}
			i++
			_, w := utf8.DecodeRuneInString(glob[i:])
			sb.WriteString(regexp.QuoteMeta(glob[i : i+w]))
			i += w - 1
		case '*':
			if i+1 < n && glob[i+1] == '*' {
				sb.WriteString(`.*`)
				i++
			} else {
				sb.WriteString(`[^/]*`)
			}
		case '?':
			sb.WriteString(`[^/]`)
		case '[':
			i = globBracket(glob, i, &sb)
		case '{':
			if groupOpen {
				panic(jutil.NewPatternSyntaxException("Cannot nest groups", glob, i))
			}
			groupOpen = true
			sb.WriteString(`(?:`)
		case '}':
			if groupOpen {
				groupOpen = false
				sb.WriteString(`)`)
			} else {
				sb.WriteString(`\}`)
			}
		case ',':
			if groupOpen {
				sb.WriteString(`|`)
			} else {
				sb.WriteString(`,`)
			}
		default:
			_, w := utf8.DecodeRuneInString(glob[i:])
			sb.WriteString(regexp.QuoteMeta(glob[i : i+w]))
			i += w - 1
		}
	}
	if groupOpen {
		panic(jutil.NewPatternSyntaxException("Missing '}", glob, n-1))
	}
	sb.WriteString(`$`)
	return regexp.MustCompile(sb.String())
}

func globBracket(glob string, start int, sb *strings.Builder) int {
	n := len(glob)
	i := start + 1
	negate := false
	if i < n && glob[i] == '!' {
		negate = true
		i++
	}
	type rng struct{ lo, hi rune }
	var ranges []rng
	for ; i < n && glob[i] != ']'; i++ {
		loIndex := i
		lo, w := utf8.DecodeRuneInString(glob[i:])
		i += w - 1
		if lo == '/' {
			panic(jutil.NewPatternSyntaxException("Explicit 'name separator' in class", glob, i))
		}
		hi := lo
		if i+2 < n && glob[i+1] == '-' && glob[i+2] != ']' {
			i += 2
			var hw int
			hi, hw = utf8.DecodeRuneInString(glob[i:])
			i += hw - 1
			if hi < lo {
				panic(jutil.NewPatternSyntaxException("Invalid range", glob, loIndex))
			}
		}
		ranges = append(ranges, rng{lo, hi})
	}
	if i >= n {
		panic(jutil.NewPatternSyntaxException("Missing ']", glob, n-1))
	}
	if len(ranges) == 0 {
		panic(jutil.NewPatternSyntaxException("Unclosed character class", glob, i))
	}
	sb.WriteByte('[')
	if negate {
		sb.WriteString(`^/`)
	}
	writeClassRune := func(r rune) {
		switch r {
		case '\\', ']', '[', '^', '-':
			sb.WriteByte('\\')
		}
		sb.WriteRune(r)
	}
	for _, r := range ranges {
		for _, part := range [2]rng{{r.lo, min(r.hi, '/'-1)}, {max(r.lo, '/'+1), r.hi}} {
			if part.lo > part.hi {
				continue
			}
			writeClassRune(part.lo)
			if part.hi != part.lo {
				sb.WriteByte('-')
				writeClassRune(part.hi)
			}
		}
	}
	sb.WriteByte(']')
	return i
}
