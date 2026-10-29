// Port of com/google/javascript/jscomp/modules/ModuleMetadataMap.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/modules/AutoValue_ModuleMetadataMap_ModuleMetadata.java @ closure-compiler v20261021

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
	"iter"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// StringCountEntry pairs a distinct string with its occurrence count.
type StringCountEntry struct {
	Key   string
	Count int
}

const stringCountsMapThreshold = 8

// StringCounts is a compact, insertion-ordered string multiset whose zero
// value is empty and allocates zero heap memory until an element is added.
type StringCounts struct {
	entries []StringCountEntry
	index   map[string]int
	size    int
}

// Add records one occurrence of s, preserving first-seen element order.
func (sc *StringCounts) Add(s string) *StringCounts {
	if sc.index != nil {
		if idx, ok := sc.index[s]; ok {
			sc.entries[idx].Count++
			sc.size++
			return sc
		}
		sc.index[s] = len(sc.entries)
		sc.entries = append(sc.entries, StringCountEntry{Key: s, Count: 1})
		sc.size++
		return sc
	}
	for i := range sc.entries {
		if sc.entries[i].Key == s {
			sc.entries[i].Count++
			sc.size++
			return sc
		}
	}
	if len(sc.entries) >= stringCountsMapThreshold {
		sc.index = make(map[string]int, len(sc.entries)+4)
		for i, e := range sc.entries {
			sc.index[e.Key] = i
		}
		sc.index[s] = len(sc.entries)
	}
	sc.entries = append(sc.entries, StringCountEntry{Key: s, Count: 1})
	sc.size++
	return sc
}

// AddAll adds all elements yielded by es.
func (sc *StringCounts) AddAll(es iter.Seq[string]) *StringCounts {
	for s := range es {
		sc.Add(s)
	}
	return sc
}

// Size returns the total number of occurrences across all distinct elements.
func (sc *StringCounts) Size() int {
	if sc == nil {
		return 0
	}
	return sc.size
}

// DistinctCount returns the number of distinct elements.
func (sc *StringCounts) DistinctCount() int {
	if sc == nil {
		return 0
	}
	return len(sc.entries)
}

// IsEmpty reports whether the multiset has no elements.
func (sc *StringCounts) IsEmpty() bool {
	return sc == nil || sc.size == 0
}

// Count returns the multiplicity of s (0 if absent).
func (sc *StringCounts) Count(s string) int {
	if sc == nil || len(sc.entries) == 0 {
		return 0
	}
	if sc.index != nil {
		if idx, ok := sc.index[s]; ok {
			return sc.entries[idx].Count
		}
		return 0
	}
	for _, e := range sc.entries {
		if e.Key == s {
			return e.Count
		}
	}
	return 0
}

// Contains reports whether s occurs at least once.
func (sc *StringCounts) Contains(s string) bool {
	return sc.Count(s) > 0
}

// All yields each element Count times in first-seen element order.
func (sc *StringCounts) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		if sc == nil {
			return
		}
		for _, e := range sc.entries {
			for range e.Count {
				if !yield(e.Key) {
					return
				}
			}
		}
	}
}

// Entries returns the underlying (Key, Count) slice in first-seen order.
func (sc *StringCounts) Entries() []StringCountEntry {
	if sc == nil {
		return nil
	}
	return sc.entries
}

// ToSlice materializes all occurrences (with multiplicity) into a slice.
func (sc *StringCounts) ToSlice() []string {
	if sc == nil || sc.size == 0 {
		return []string{}
	}
	out := make([]string, 0, sc.size)
	for _, e := range sc.entries {
		for range e.Count {
			out = append(out, e.Key)
		}
	}
	return out
}

// ElementSet returns a zero-allocation view over the distinct keys in first-seen order.
func (sc *StringCounts) ElementSet() StringCountsElements {
	return StringCountsElements{sc: sc}
}

// String renders in Guava AbstractMultiset.toString() format: "[a, b x 2]".
func (sc *StringCounts) String() string {
	if sc == nil || len(sc.entries) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, e := range sc.entries {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(e.Key)
		if e.Count > 1 {
			sb.WriteString(" x ")
			sb.WriteString(strconv.Itoa(e.Count))
		}
	}
	sb.WriteByte(']')
	return sb.String()
}

// StringCountsElements is a zero-allocation view over the distinct strings in a StringCounts.
type StringCountsElements struct {
	sc *StringCounts
}

// Size returns the number of distinct elements.
func (e StringCountsElements) Size() int { return e.sc.DistinctCount() }

// IsEmpty reports whether there are no distinct elements.
func (e StringCountsElements) IsEmpty() bool { return e.sc.IsEmpty() }

// Contains reports whether s is present.
func (e StringCountsElements) Contains(s string) bool { return e.sc.Contains(s) }

// All yields each distinct element once in first-seen order.
func (e StringCountsElements) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		if e.sc == nil {
			return
		}
		for _, entry := range e.sc.entries {
			if !yield(entry.Key) {
				return
			}
		}
	}
}

// ToSlice returns the distinct elements as a slice in first-seen order.
func (e StringCountsElements) ToSlice() []string {
	if e.sc == nil || len(e.sc.entries) == 0 {
		return nil
	}
	out := make([]string, len(e.sc.entries))
	for i, entry := range e.sc.entries {
		out[i] = entry.Key
	}
	return out
}

// String renders like AbstractCollection.toString(): "[a, b]".
func (e StringCountsElements) String() string {
	if e.sc == nil || len(e.sc.entries) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, entry := range e.sc.entries {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(entry.Key)
	}
	sb.WriteByte(']')
	return sb.String()
}

// ModuleMetadataMap contains metadata around modules (or scripts) that is
// useful for checking imports / requires.
//
// TODO(johnplaisted): There's an opportunity for reuse here in
// ClosureRewriteModules, which would involve putting this in some common
// location. Currently this is only used as a helper class for
// Es6RewriteModules. CompilerInput already has some (not all) of this
// information but it is not always populated. It may also be ideal to
// include CommonJS here too as ES6 modules can import them. That would allow
// decoupling of how these modules are written; right now Es6RewriteModule
// only checks this for goog.requires and goog: imports, not for ES6 path
// imports.
type ModuleMetadataMap struct {
	// Map from module path to module. These modules represent files and thus will contain all goog
	// namespaces that are in the file. These are not the same modules in modulesByGoogNamespace.
	modulesByPath *jutil.LinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]

	// Map from Closure namespace to module. These modules represent just the single namespace and
	// thus each module has only one goog namespace in its ModuleMetadata#googNamespaces().
	// These are not the same modules in modulesByPath.
	modulesByGoogNamespace *jutil.LinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]

	moduleMetadata *jutil.LinkedHashSet[*ModuleMetadataMap_ModuleMetadata]
}

// ModuleMetadataMap_ModuleType lists the various types of Javascript
// modules and scripts that can be found in the JS Compiler.
//
// Java: ModuleMetadataMap.ModuleType
type ModuleMetadataMap_ModuleType int32

const (
	ModuleMetadataMap_ModuleType_ES6_MODULE ModuleMetadataMap_ModuleType = iota
	// A goog.module that does not declare a legacy namespace.
	ModuleMetadataMap_ModuleType_GOOG_MODULE
	// A goog.module that declares a legacy namespace with goog.module.declareLegacyNamespace.
	ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE
	ModuleMetadataMap_ModuleType_COMMON_JS
	// The following two cases are not actually modules, but are useful to include in the map.
	ModuleMetadataMap_ModuleType_GOOG_PROVIDE
	ModuleMetadataMap_ModuleType_SCRIPT
)

var moduleMetadataMap_ModuleType_names = [...]string{
	"ES6_MODULE", "GOOG_MODULE", "LEGACY_GOOG_MODULE", "COMMON_JS", "GOOG_PROVIDE", "SCRIPT",
}

var moduleMetadataMap_ModuleType_description = [...]string{
	"an ES6 module",
	"a goog.module",
	"a goog.module",
	"a CommonJS module",
	"a script file that contains at least one goog.provide",
	"a script file that does not contain a goog.provide",
}

// String is Enum.name().
func (t ModuleMetadataMap_ModuleType) String() string { return moduleMetadataMap_ModuleType_names[t] }

// Ordinal is Enum.ordinal().
func (t ModuleMetadataMap_ModuleType) Ordinal() int { return int(t) }

// Description is the field ModuleType.description.
func (t ModuleMetadataMap_ModuleType) Description() string {
	return moduleMetadataMap_ModuleType_description[t]
}

// Java: ModuleMetadataMap(Map<String, ModuleMetadata>, Map<String, ModuleMetadata>)
func NewModuleMetadataMap(
	modulesByPath *jutil.LinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata],
	modulesByGoogNamespace *jutil.LinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]) *ModuleMetadataMap {
	m := &ModuleMetadataMap{
		modulesByPath:          moduleMetadataMap_immutableMapCopyOf(modulesByPath),
		modulesByGoogNamespace: moduleMetadataMap_immutableMapCopyOf(modulesByGoogNamespace),
	}
	m.moduleMetadata = jutil.NewLinkedHashSetWithCapacity[*ModuleMetadataMap_ModuleMetadata](
		modulesByPath.Size() + modulesByGoogNamespace.Size())
	for v := range modulesByPath.Values() {
		m.moduleMetadata.Add(jutil.CheckNotNull(v))
	}
	for v := range modulesByGoogNamespace.Values() {
		m.moduleMetadata.Add(jutil.CheckNotNull(v))
	}
	return m
}

// moduleMetadataMap_immutableMapCopyOf is ImmutableMap.copyOf(Map), which
// rejects null values.
func moduleMetadataMap_immutableMapCopyOf(
	m *jutil.LinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]) *jutil.LinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata] {
	out := jutil.NewLinkedHashMapWithCapacity[string, *ModuleMetadataMap_ModuleMetadata](m.Size())
	for k, v := range m.All() {
		out.Put(k, jutil.CheckNotNull(v))
	}
	return out
}

// ModuleMetadataMap_ModuleMetadata is the struct containing basic
// information about a module/script including its type and goog
// namespaces.
//
// It uses reference equality (Java overrides equals and hashCode to call
// Object's), so the Go pointer identity is its equality.
//
// Go-only: language-adapter — the AutoValue-generated
// AutoValue_ModuleMetadataMap_ModuleMetadata is folded into this struct,
// and the seven string multisets are inline zero-alloc StringCounts values.
type ModuleMetadataMap_ModuleMetadata struct {
	rootNode                          rhino.Node
	path                              *ModuleLoader_ModulePath
	googNamespaces                    StringCounts
	stronglyRequiredGoogNamespaces    StringCounts
	dynamicallyRequiredGoogNamespaces StringCounts
	maybeRequiredGoogNamespaces       StringCounts
	weaklyRequiredGoogNamespaces      StringCounts
	es6ImportSpecifiers               StringCounts
	nestedModules                     []*ModuleMetadataMap_ModuleMetadata
	readToggles                       StringCounts
	moduleType                        ModuleMetadataMap_ModuleType
	usesClosure                       bool
	isTestOnly                        bool
}

// Java: ModuleMetadataMap.ModuleMetadata.moduleType()
func (m *ModuleMetadataMap_ModuleMetadata) ModuleType() ModuleMetadataMap_ModuleType {
	return m.moduleType
}

// Java: ModuleMetadataMap.ModuleMetadata.isEs6Module()
func (m *ModuleMetadataMap_ModuleMetadata) IsEs6Module() bool {
	return m.ModuleType() == ModuleMetadataMap_ModuleType_ES6_MODULE
}

// Java: ModuleMetadataMap.ModuleMetadata.isGoogModule()
func (m *ModuleMetadataMap_ModuleMetadata) IsGoogModule() bool {
	return m.IsNonLegacyGoogModule() || m.IsLegacyGoogModule()
}

// Java: ModuleMetadataMap.ModuleMetadata.isNonLegacyGoogModule()
func (m *ModuleMetadataMap_ModuleMetadata) IsNonLegacyGoogModule() bool {
	return m.ModuleType() == ModuleMetadataMap_ModuleType_GOOG_MODULE
}

// Java: ModuleMetadataMap.ModuleMetadata.isLegacyGoogModule()
func (m *ModuleMetadataMap_ModuleMetadata) IsLegacyGoogModule() bool {
	return m.ModuleType() == ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE
}

// Java: ModuleMetadataMap.ModuleMetadata.isGoogProvide()
func (m *ModuleMetadataMap_ModuleMetadata) IsGoogProvide() bool {
	return m.ModuleType() == ModuleMetadataMap_ModuleType_GOOG_PROVIDE
}

// Java: ModuleMetadataMap.ModuleMetadata.hasLegacyGoogNamespaces()
func (m *ModuleMetadataMap_ModuleMetadata) HasLegacyGoogNamespaces() bool {
	return m.IsGoogProvide() || m.IsLegacyGoogModule()
}

// Java: ModuleMetadataMap.ModuleMetadata.isCommonJs()
func (m *ModuleMetadataMap_ModuleMetadata) IsCommonJs() bool {
	return m.ModuleType() == ModuleMetadataMap_ModuleType_COMMON_JS
}

// Java: ModuleMetadataMap.ModuleMetadata.isNonProvideScript()
func (m *ModuleMetadataMap_ModuleMetadata) IsNonProvideScript() bool {
	return m.ModuleType() == ModuleMetadataMap_ModuleType_SCRIPT
}

// IsModule returns whether this is a module (with it's own local scope).
//
// Java: ModuleMetadataMap.ModuleMetadata.isModule()
func (m *ModuleMetadataMap_ModuleMetadata) IsModule() bool {
	switch m.ModuleType() {
	case ModuleMetadataMap_ModuleType_GOOG_PROVIDE, ModuleMetadataMap_ModuleType_SCRIPT:
		return false
	default: // COMMON_JS, ES6_MODULE, GOOG_MODULE, LEGACY_GOOG_MODULE
		return true
	}
}

// RootNode is the AST node that represents the root of this module.
//
// May be nil if this is a synthetic piece of metadata - e.g. in a test, or
// something used as a fallback.
//
// Java: ModuleMetadataMap.ModuleMetadata.rootNode()
func (m *ModuleMetadataMap_ModuleMetadata) RootNode() rhino.Node {
	return m.rootNode
}

// UsesClosure is whether this file uses Closure Library at all. Note that a
// file could use Closure Library even without calling
// goog.provide/module/require - there are some primitives in base.js that
// can be used without being required like goog.isArray.
//
// If this is true this indicates the base.js is needed and is not part of
// this script - it is an EXTERNAL dependencym otherwise false. If this is
// also false if Closure Library is part of this script - e.g. a bundle with
// base.js. So something could be a "goog.provide'd file", but not use
// Closure if the bundle already contains Closure.
//
// Java: ModuleMetadataMap.ModuleMetadata.usesClosure()
func (m *ModuleMetadataMap_ModuleMetadata) UsesClosure() bool {
	return m.usesClosure
}

// IsTestOnly is whether goog.setTestOnly was called.
//
// Java: ModuleMetadataMap.ModuleMetadata.isTestOnly()
func (m *ModuleMetadataMap_ModuleMetadata) IsTestOnly() bool {
	return m.isTestOnly
}

// GoogNamespaces are the Closure namespaces that this file is associated
// with. Created by goog.provide, goog.module, and goog.declareModuleId.
//
// This is a multiset as it does not warn on duplicate namespaces, but will
// still encapsulate that information with this multiset.
//
// Java: ModuleMetadataMap.ModuleMetadata.googNamespaces()
func (m *ModuleMetadataMap_ModuleMetadata) GoogNamespaces() *StringCounts {
	return &m.googNamespaces
}

// StronglyRequiredGoogNamespaces are the Closure namespaces this file
// strongly requires, i.e., arguments to goog.require calls.
//
// This is a multiset as it does not warn on duplicate namespaces, but will
// still encapsulate that information with this multiset.
//
// Java: ModuleMetadataMap.ModuleMetadata.stronglyRequiredGoogNamespaces()
func (m *ModuleMetadataMap_ModuleMetadata) StronglyRequiredGoogNamespaces() *StringCounts {
	return &m.stronglyRequiredGoogNamespaces
}

// DynamicallyRequiredGoogNamespaces are the Closure namespaces this file
// dynamically require, i.e., arguments to goog.requireDynamic() calls.
//
// This is a multiset as it does not warn on duplicate namespaces, but will
// still encapsulate that information with this multiset.
//
// Java: ModuleMetadataMap.ModuleMetadata.dynamicallyRequiredGoogNamespaces()
func (m *ModuleMetadataMap_ModuleMetadata) DynamicallyRequiredGoogNamespaces() *StringCounts {
	return &m.dynamicallyRequiredGoogNamespaces
}

// MaybeRequiredGoogNamespaces are the Closure namespaces this file "maybe"
// require, i.e., arguments to
// goog.maybeRequireFrameworkInternalOnlyDoNotCallOrElse() calls.
//
// This is a multiset as it does not warn on duplicate namespaces, but will
// still encapsulate that information with this multiset.
//
// Java: ModuleMetadataMap.ModuleMetadata.maybeRequiredGoogNamespaces()
func (m *ModuleMetadataMap_ModuleMetadata) MaybeRequiredGoogNamespaces() *StringCounts {
	return &m.maybeRequiredGoogNamespaces
}

// WeaklyRequiredGoogNamespaces are the Closure namespaces this file weakly
// requires, i.e., arguments to goog.requireType calls.
//
// This is a multiset as it does not warn on duplicate namespaces, but will
// still encapsulate that information with this multiset.
//
// Java: ModuleMetadataMap.ModuleMetadata.weaklyRequiredGoogNamespaces()
func (m *ModuleMetadataMap_ModuleMetadata) WeaklyRequiredGoogNamespaces() *StringCounts {
	return &m.weaklyRequiredGoogNamespaces
}

// Es6ImportSpecifiers is the raw text of all ES6 import specifiers
// (includes "export from" as well).
//
// Java: ModuleMetadataMap.ModuleMetadata.es6ImportSpecifiers()
func (m *ModuleMetadataMap_ModuleMetadata) Es6ImportSpecifiers() *StringCounts {
	return &m.es6ImportSpecifiers
}

// Java: ModuleMetadataMap.ModuleMetadata.nestedModules()
func (m *ModuleMetadataMap_ModuleMetadata) NestedModules() []*ModuleMetadataMap_ModuleMetadata {
	return m.nestedModules
}

// ReadToggles are the arguments to
// goog.readToggleInternalDoNotCallDirectly() calls.
//
// This is a multiset as it does not warn on duplicate toggles, but will
// still encapsulate that information with this multiset.
//
// Java: ModuleMetadataMap.ModuleMetadata.readToggles()
func (m *ModuleMetadataMap_ModuleMetadata) ReadToggles() *StringCounts {
	return &m.readToggles
}

// Path may return nil.
//
// Java: ModuleMetadataMap.ModuleMetadata.path()
func (m *ModuleMetadataMap_ModuleMetadata) Path() *ModuleLoader_ModulePath {
	return m.path
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.toString()
func (m *ModuleMetadataMap_ModuleMetadata) String() string {
	rootNode := "null"
	if !m.rootNode.IsNil() {
		rootNode = m.rootNode.String()
	}
	path := "null"
	if m.path != nil {
		path = m.path.String()
	}
	nested := "[]"
	if len(m.nestedModules) > 0 {
		var sb strings.Builder
		sb.WriteByte('[')
		for i, nm := range m.nestedModules {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(nm.String())
		}
		sb.WriteByte(']')
		nested = sb.String()
	}
	return "ModuleMetadata{" +
		"moduleType=" + m.moduleType.String() + ", " +
		"rootNode=" + rootNode + ", " +
		"usesClosure=" + strconv.FormatBool(m.usesClosure) + ", " +
		"isTestOnly=" + strconv.FormatBool(m.isTestOnly) + ", " +
		"googNamespaces=" + m.googNamespaces.String() + ", " +
		"stronglyRequiredGoogNamespaces=" + m.stronglyRequiredGoogNamespaces.String() + ", " +
		"dynamicallyRequiredGoogNamespaces=" + m.dynamicallyRequiredGoogNamespaces.String() + ", " +
		"maybeRequiredGoogNamespaces=" + m.maybeRequiredGoogNamespaces.String() + ", " +
		"weaklyRequiredGoogNamespaces=" + m.weaklyRequiredGoogNamespaces.String() + ", " +
		"es6ImportSpecifiers=" + m.es6ImportSpecifiers.String() + ", " +
		"nestedModules=" + nested + ", " +
		"readToggles=" + m.readToggles.String() + ", " +
		"path=" + path +
		"}"
}

// ModuleMetadataMap_ModuleMetadata_Builder is a 1-word value builder for
// ModuleMetadata (AutoValue_ModuleMetadataMap_ModuleMetadata.Builder).
type ModuleMetadataMap_ModuleMetadata_Builder struct {
	m *ModuleMetadataMap_ModuleMetadata
}

// Java: ModuleMetadataMap.ModuleMetadata.builder()
func ModuleMetadataMap_ModuleMetadata_Builder_() ModuleMetadataMap_ModuleMetadata_Builder {
	return ModuleMetadataMap_ModuleMetadata_Builder{m: &ModuleMetadataMap_ModuleMetadata{}}
}

// Java: ModuleMetadataMap.ModuleMetadata.Builder.addGoogNamespace(String)
func (b ModuleMetadataMap_ModuleMetadata_Builder) AddGoogNamespace(namespace string) ModuleMetadataMap_ModuleMetadata_Builder {
	b.m.googNamespaces.Add(namespace)
	return b
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.moduleType(ModuleType)
func (b ModuleMetadataMap_ModuleMetadata_Builder) ModuleType(moduleType ModuleMetadataMap_ModuleType) ModuleMetadataMap_ModuleMetadata_Builder {
	b.m.moduleType = moduleType
	return b
}

// GetModuleType is the getter overload moduleType().
//
// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.moduleType()
func (b ModuleMetadataMap_ModuleMetadata_Builder) GetModuleType() ModuleMetadataMap_ModuleType {
	return b.m.moduleType
}

// RootNode takes a possibly-nil root.
//
// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.rootNode(Node)
func (b ModuleMetadataMap_ModuleMetadata_Builder) RootNode(rootNode rhino.Node) ModuleMetadataMap_ModuleMetadata_Builder {
	b.m.rootNode = rootNode
	return b
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.usesClosure(boolean)
func (b ModuleMetadataMap_ModuleMetadata_Builder) UsesClosure(usesClosure bool) ModuleMetadataMap_ModuleMetadata_Builder {
	b.m.usesClosure = usesClosure
	return b
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.isTestOnly(boolean)
func (b ModuleMetadataMap_ModuleMetadata_Builder) IsTestOnly(isTestOnly bool) ModuleMetadataMap_ModuleMetadata_Builder {
	b.m.isTestOnly = isTestOnly
	return b
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.googNamespacesBuilder()
func (b ModuleMetadataMap_ModuleMetadata_Builder) GoogNamespacesBuilder() *StringCounts {
	return &b.m.googNamespaces
}

// HasGoogNamespace reports whether namespace has been added to googNamespacesBuilder.
func (b ModuleMetadataMap_ModuleMetadata_Builder) HasGoogNamespace(namespace string) bool {
	return b.m.googNamespaces.Contains(namespace)
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.stronglyRequiredGoogNamespacesBuilder()
func (b ModuleMetadataMap_ModuleMetadata_Builder) StronglyRequiredGoogNamespacesBuilder() *StringCounts {
	return &b.m.stronglyRequiredGoogNamespaces
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.dynamicallyRequiredGoogNamespacesBuilder()
func (b ModuleMetadataMap_ModuleMetadata_Builder) DynamicallyRequiredGoogNamespacesBuilder() *StringCounts {
	return &b.m.dynamicallyRequiredGoogNamespaces
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.maybeRequiredGoogNamespacesBuilder()
func (b ModuleMetadataMap_ModuleMetadata_Builder) MaybeRequiredGoogNamespacesBuilder() *StringCounts {
	return &b.m.maybeRequiredGoogNamespaces
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.weaklyRequiredGoogNamespacesBuilder()
func (b ModuleMetadataMap_ModuleMetadata_Builder) WeaklyRequiredGoogNamespacesBuilder() *StringCounts {
	return &b.m.weaklyRequiredGoogNamespaces
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.es6ImportSpecifiersBuilder()
func (b ModuleMetadataMap_ModuleMetadata_Builder) Es6ImportSpecifiersBuilder() *StringCounts {
	return &b.m.es6ImportSpecifiers
}

// AddNestedModule appends module to the nestedModules builder list.
func (b ModuleMetadataMap_ModuleMetadata_Builder) AddNestedModule(module *ModuleMetadataMap_ModuleMetadata) {
	b.m.nestedModules = append(b.m.nestedModules, jutil.CheckNotNull(module))
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.readTogglesBuilder()
func (b ModuleMetadataMap_ModuleMetadata_Builder) ReadTogglesBuilder() *StringCounts {
	return &b.m.readToggles
}

// Path takes a possibly-nil path.
//
// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.path(ModulePath)
func (b ModuleMetadataMap_ModuleMetadata_Builder) Path(path *ModuleLoader_ModulePath) ModuleMetadataMap_ModuleMetadata_Builder {
	b.m.path = path
	return b
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.build()
func (b ModuleMetadataMap_ModuleMetadata_Builder) Build() *ModuleMetadataMap_ModuleMetadata {
	return b.m
}

// GetModulesByPath returns the map from module path to module. These
// modules represent files and thus ModuleMetadata#googNamespaces() contains
// all Closure namespaces in the file. These are not the same modules from
// ModuleMetadataMap#getModulesByGoogNamespace(). It is not valid to call
// ModuleRenaming#getGlobalName(ModuleMetadata, String) on
// ModuleType#GOOG_PROVIDE modules from this map that have more than one
// Closure namespace as it is ambiguous.
//
// Java: ModuleMetadataMap.getModulesByPath()
func (m *ModuleMetadataMap) GetModulesByPath() *jutil.LinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata] {
	return m.modulesByPath
}

// GetModulesByGoogNamespace returns the map from Closure namespace to
// module. These modules represent the Closure namespace and thus
// ModuleMetadata#googNamespaces() will have size 1. As a result, it is
// valid to call ModuleRenaming#getGlobalName(ModuleMetadata, String) on
// these modules. These are not the same modules from
// ModuleMetadataMap#getModulesByPath().
//
// Java: ModuleMetadataMap.getModulesByGoogNamespace()
func (m *ModuleMetadataMap) GetModulesByGoogNamespace() *jutil.LinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata] {
	return m.modulesByGoogNamespace
}

// GetAllModuleMetadata is the set of all modules across both maps.
//
// goog.loadModule calls have no associated path, and non-Closure modules
// have no namespaces.
//
// Java: ModuleMetadataMap.getAllModuleMetadata()
func (m *ModuleMetadataMap) GetAllModuleMetadata() *jutil.LinkedHashSet[*ModuleMetadataMap_ModuleMetadata] {
	return m.moduleMetadata
}

// Java: ModuleMetadataMap.emptyForTesting()
func ModuleMetadataMap_EmptyForTesting() *ModuleMetadataMap {
	return NewModuleMetadataMap(
		jutil.NewLinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata](),
		jutil.NewLinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]())
}
