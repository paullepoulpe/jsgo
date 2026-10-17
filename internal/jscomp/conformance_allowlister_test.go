// Port of test/com/google/javascript/jscomp/ConformanceAllowlisterTest.java @ closure-compiler v20260415

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"slices"
	"strconv"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// conformanceAllowlisterTest_entry is one (source file, line) pair of the
// ImmutableMultimap<String, Integer> testConformanceAllowlister returns.
type conformanceAllowlisterTest_entry struct {
	file string
	line int
}

func conformanceAllowlisterTest_source(path string) []*SourceFile {
	return []*SourceFile{SourceFile_FromCode(
		path,
		"var foo = document.getElementById('name');\n"+
			"foo."+map[bool]string{true: "innerHTML", false: "title"}[path == "/test/entry.js"]+" = 'test';\n")}
}

// conformanceAllowlisterTest_containsExactly is Truth's containsExactly on
// the multimap: the pairs in any order.
func conformanceAllowlisterTest_containsExactly(t *testing.T, got []conformanceAllowlisterTest_entry, want ...conformanceAllowlisterTest_entry) {
	t.Helper()
	key := func(e conformanceAllowlisterTest_entry) string { return e.file + ":" + strconv.Itoa(e.line) }
	var g, w []string
	for _, e := range got {
		g = append(g, key(e))
	}
	for _, e := range want {
		w = append(w, key(e))
	}
	slices.Sort(g)
	slices.Sort(w)
	if !slices.Equal(g, w) {
		t.Errorf("violations = %q, want %q", g, w)
	}
}

func TestConformanceAllowlister_ConformanceAllowlistAddNew(t *testing.T) {
	sources := conformanceAllowlisterTest_source("/entry.js")

	requirement := Requirement_NewBuilder().
		SetType(Requirement_Type_BANNED_PROPERTY).
		SetErrorMessage("Lorem Ipsum").
		AddValue("Object.prototype.title").
		Build()

	conformanceAllowlisterTest_containsExactly(t, conformanceAllowlisterTest_testConformanceAllowlister(t, sources, requirement),
		conformanceAllowlisterTest_entry{"/entry.js", 2})
}

func TestConformanceAllowlister_ConformanceAllowlistRemove(t *testing.T) {
	sources := conformanceAllowlisterTest_source("/entry.js")

	whitelistRequirement := Requirement_NewBuilder().
		SetType(Requirement_Type_BANNED_PROPERTY).
		SetErrorMessage("Lorem Ipsum").
		AddValue("Object.prototype.innerHTML").
		AddWhitelist("/entry.js").
		Build()

	conformanceAllowlisterTest_containsExactly(t, conformanceAllowlisterTest_testConformanceAllowlister(t, sources, whitelistRequirement))

	allowlistRequirement := Requirement_NewBuilder().
		SetType(Requirement_Type_BANNED_PROPERTY).
		SetErrorMessage("Lorem Ipsum").
		AddValue("Object.prototype.innerHTML").
		AddAllowlist("/entry.js").
		Build()

	conformanceAllowlisterTest_containsExactly(t, conformanceAllowlisterTest_testConformanceAllowlister(t, sources, allowlistRequirement))
}

func TestConformanceAllowlister_ConformanceAllowlistPreserve(t *testing.T) {
	sources := conformanceAllowlisterTest_source("/entry.js")

	whitelistRequirement := Requirement_NewBuilder().
		SetType(Requirement_Type_BANNED_PROPERTY).
		SetErrorMessage("Lorem Ipsum").
		AddValue("Object.prototype.title").
		AddWhitelist("/entry.js").
		Build()

	conformanceAllowlisterTest_containsExactly(t, conformanceAllowlisterTest_testConformanceAllowlister(t, sources, whitelistRequirement),
		conformanceAllowlisterTest_entry{"/entry.js", 2})

	allowlistRequirement := Requirement_NewBuilder().
		SetType(Requirement_Type_BANNED_PROPERTY).
		SetErrorMessage("Lorem Ipsum").
		AddValue("Object.prototype.title").
		AddAllowlist("/entry.js").
		Build()

	conformanceAllowlisterTest_containsExactly(t, conformanceAllowlisterTest_testConformanceAllowlister(t, sources, allowlistRequirement),
		conformanceAllowlisterTest_entry{"/entry.js", 2})
}

// TODO(bangert): Evaluate if this is always the best behaviour.
// The current behaviour pushes the behaviour of how to cluster the allowlist
// to the program driving ConformanceAllowlister.
func TestConformanceAllowlister_ConformanceAllowlistBreaksDownFolder(t *testing.T) {
	sources := conformanceAllowlisterTest_source("/test/entry.js")

	whitelistRequirement := Requirement_NewBuilder().
		SetType(Requirement_Type_BANNED_PROPERTY).
		SetErrorMessage("Lorem Ipsum").
		AddValue("Object.prototype.innerHTML").
		AddWhitelist("/test/").
		Build()

	conformanceAllowlisterTest_containsExactly(t, conformanceAllowlisterTest_testConformanceAllowlister(t, sources, whitelistRequirement),
		conformanceAllowlisterTest_entry{"/test/entry.js", 2})

	allowlistRequirement := Requirement_NewBuilder().
		SetType(Requirement_Type_BANNED_PROPERTY).
		SetErrorMessage("Lorem Ipsum").
		AddValue("Object.prototype.innerHTML").
		AddAllowlist("/test/").
		Build()

	conformanceAllowlisterTest_containsExactly(t, conformanceAllowlisterTest_testConformanceAllowlister(t, sources, allowlistRequirement),
		conformanceAllowlisterTest_entry{"/test/entry.js", 2})
}

// Java: ConformanceAllowlisterTest.testConformanceAllowlister(ImmutableList<SourceFile>, Requirement)
func conformanceAllowlisterTest_testConformanceAllowlister(
	t *testing.T, sources []*SourceFile, config *Requirement) []conformanceAllowlisterTest_entry {
	t.Helper()
	options := NewCompilerOptions()
	options.SetCheckTypes(true)
	// TODO(bangert): Support banned property on OBJECT even if types are not
	// checked.
	options.SetChecksOnly(true)
	externs := []*SourceFile{}
	compiler := NewCompiler()
	result := compiler.Compile(externs, sources, options)
	if !result.Success {
		t.Fatalf("compilation failed: %v", result.Errors)
	}

	var errors []conformanceAllowlisterTest_entry
	keys := jutil.NewLinkedHashSet[string]()
	for node := range ConformanceAllowlister_GetViolatingNodes(
		compiler, compiler.GetExternsRoot(), compiler.GetJsRoot(), config).All() {
		file := *node.GetSourceFileName()
		errors = append(errors, conformanceAllowlisterTest_entry{file, node.GetLineno()})
		keys.Add(file)
	}
	paths := ConformanceAllowlister_GetViolatingPaths(
		compiler, compiler.GetExternsRoot(), compiler.GetJsRoot(), config)
	g, w := keys.ToSlice(), paths.ToSlice()
	slices.Sort(g)
	slices.Sort(w)
	if !slices.Equal(g, w) {
		t.Errorf("violating node files %q, violating paths %q", g, w)
	}
	return errors
}
