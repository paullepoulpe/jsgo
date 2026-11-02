// Port of com/google/javascript/refactoring/RequireAliasGenerator.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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

package refactoring

import (
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// RequireAliasGenerator implements go/js-style#file-goog-require for
// selecting import aliases.
//
// Instances act as lazy generators for an infinite sequence of names, in
// order of preferability. Callers can pull from this sequence until they
// find a suitable name. This design allows higher-level concerns, such as
// uniqueness, to be managed by higher-level code.
//
// This impelementation is a heuristic for selecting likely good aliases. In
// addition to conforming to the style guide, the ideal heurisitc is simple
// engough to explain in a few bullets, in priority order:
//
//   - Alias is not skiplisted
//   - Alias is always only composed of segments of the full namespace
//   - Shorter aliases are preferred over longer ones
type RequireAliasGenerator struct {
	parts      []string
	capitalize bool
}

var requireAliasGenerator_SKIPLISTED_ALIASES = requireAliasGenerator_createSkiplistedAliases()

// Java: RequireAliasGenerator.createSkiplistedAliases()
func requireAliasGenerator_createSkiplistedAliases() map[string]struct{} {
	seeds := []string{
		"", //
		"array",
		"event",
		"map",
		"math",
		"object",
		"promise",
		"set",
		"string",
		"thenable",
	}

	builder := make(map[string]struct{}, len(seeds)*2)
	for _, seed := range seeds {
		builder[requireAliasGenerator_captializeFirstChar(seed)] = struct{}{}
		builder[jutil.StringToLowerCaseASCII(seed)] = struct{}{}
	}
	return builder
}

// RequireAliasGenerator_Over returns the Iterable<String> of aliases for
// namespace.
//
// Java: RequireAliasGenerator.over(String)
func RequireAliasGenerator_Over(namespace string) *RequireAliasGenerator {
	return newRequireAliasGenerator(namespace)
}

// Java: RequireAliasGenerator(String)
func newRequireAliasGenerator(namespace string) *RequireAliasGenerator {
	g := &RequireAliasGenerator{}
	// DOT_SPLITTER = Splitter.on('.')
	g.parts = strings.Split(namespace, ".")
	slices.Reverse(g.parts)
	g.capitalize = requireAliasGenerator_isUpperCase(jutil.StringCharAt(g.parts[0], 0))
	return g
}

// All iterates the (infinite) sequence: Java's for-each over the Iterable.
//
// Java: RequireAliasGenerator.iterator()
func (g *RequireAliasGenerator) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		c := requireAliasCursor{outer: g}
		for {
			alias, ok := c.next()
			if !ok || !yield(alias) {
				return
			}
		}
	}
}

type requireAliasCursor struct {
	outer  *RequireAliasGenerator
	concat string
	index  int
	suffix int
}

// Java: RequireAliasGenerator.NameIterator.computeNext()
func (c *requireAliasCursor) next() (string, bool) {
	for {
		candidate := c.nextUnfiltered()
		if _, skipped := requireAliasGenerator_SKIPLISTED_ALIASES[candidate]; !skipped {
			return candidate, true
		}
	}
}

// Java: RequireAliasGenerator.NameIterator.computeNextWithoutFiltering()
func (c *requireAliasCursor) nextUnfiltered() string {
	if c.index >= len(c.outer.parts) {
		s := c.concat + strconv.Itoa(c.suffix)
		c.suffix++
		return s
	}

	part := c.outer.parts[c.index]
	c.index++
	if c.outer.capitalize {
		part = requireAliasGenerator_captializeFirstChar(part)
	}

	c.concat = part + requireAliasGenerator_captializeFirstChar(c.concat)
	return c.concat
}

// Java: RequireAliasGenerator.captializeFirstChar(String)
func requireAliasGenerator_captializeFirstChar(w string) string {
	if w == "" || requireAliasGenerator_isUpperCase(jutil.StringCharAt(w, 0)) {
		return w
	}
	c := jutil.StringCharAt(w, 0)
	if c >= 'a' && c <= 'z' {
		c -= 'a' - 'A'
	}
	return jutil.StringFromChar(c) + jutil.StringSubstringFrom(w, 1)
}

// requireAliasGenerator_isUpperCase is Guava's Ascii.isUpperCase(char).
func requireAliasGenerator_isUpperCase(c uint16) bool {
	return c >= 'A' && c <= 'Z'
}
