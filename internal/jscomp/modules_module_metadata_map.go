// Port of com/google/javascript/jscomp/modules/ModuleMetadataMap.java @ closure-compiler v20260415
// Port of com/google/javascript/jscomp/modules/AutoValue_ModuleMetadataMap_ModuleMetadata.java @ closure-compiler v20260415

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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

// Java: ModuleMetadataMap.ModuleType.values()
func ModuleMetadataMap_ModuleType_values() []ModuleMetadataMap_ModuleType {
	return []ModuleMetadataMap_ModuleType{
		ModuleMetadataMap_ModuleType_ES6_MODULE, ModuleMetadataMap_ModuleType_GOOG_MODULE,
		ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE, ModuleMetadataMap_ModuleType_COMMON_JS,
		ModuleMetadataMap_ModuleType_GOOG_PROVIDE, ModuleMetadataMap_ModuleType_SCRIPT,
	}
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
	m.moduleMetadata = jutil.NewLinkedHashSet[*ModuleMetadataMap_ModuleMetadata]()
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
	out := jutil.NewLinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]()
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
// AutoValue_ModuleMetadataMap_ModuleMetadata is folded into this struct.
type ModuleMetadataMap_ModuleMetadata struct {
	moduleType                        ModuleMetadataMap_ModuleType
	rootNode                          *rhino.Node
	usesClosure                       bool
	isTestOnly                        bool
	googNamespaces                    *jutil.ImmutableMultiset[string]
	stronglyRequiredGoogNamespaces    *jutil.ImmutableMultiset[string]
	dynamicallyRequiredGoogNamespaces *jutil.ImmutableMultiset[string]
	maybeRequiredGoogNamespaces       *jutil.ImmutableMultiset[string]
	weaklyRequiredGoogNamespaces      *jutil.ImmutableMultiset[string]
	es6ImportSpecifiers               *jutil.ImmutableMultiset[string]
	nestedModules                     []*ModuleMetadataMap_ModuleMetadata
	readToggles                       *jutil.ImmutableMultiset[string]
	path                              *ModuleLoader_ModulePath
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
func (m *ModuleMetadataMap_ModuleMetadata) RootNode() *rhino.Node {
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
func (m *ModuleMetadataMap_ModuleMetadata) GoogNamespaces() *jutil.ImmutableMultiset[string] {
	return m.googNamespaces
}

// StronglyRequiredGoogNamespaces are the Closure namespaces this file
// strongly requires, i.e., arguments to goog.require calls.
//
// This is a multiset as it does not warn on duplicate namespaces, but will
// still encapsulate that information with this multiset.
//
// Java: ModuleMetadataMap.ModuleMetadata.stronglyRequiredGoogNamespaces()
func (m *ModuleMetadataMap_ModuleMetadata) StronglyRequiredGoogNamespaces() *jutil.ImmutableMultiset[string] {
	return m.stronglyRequiredGoogNamespaces
}

// DynamicallyRequiredGoogNamespaces are the Closure namespaces this file
// dynamically require, i.e., arguments to goog.requireDynamic() calls.
//
// This is a multiset as it does not warn on duplicate namespaces, but will
// still encapsulate that information with this multiset.
//
// Java: ModuleMetadataMap.ModuleMetadata.dynamicallyRequiredGoogNamespaces()
func (m *ModuleMetadataMap_ModuleMetadata) DynamicallyRequiredGoogNamespaces() *jutil.ImmutableMultiset[string] {
	return m.dynamicallyRequiredGoogNamespaces
}

// MaybeRequiredGoogNamespaces are the Closure namespaces this file "maybe"
// require, i.e., arguments to
// goog.maybeRequireFrameworkInternalOnlyDoNotCallOrElse() calls.
//
// This is a multiset as it does not warn on duplicate namespaces, but will
// still encapsulate that information with this multiset.
//
// Java: ModuleMetadataMap.ModuleMetadata.maybeRequiredGoogNamespaces()
func (m *ModuleMetadataMap_ModuleMetadata) MaybeRequiredGoogNamespaces() *jutil.ImmutableMultiset[string] {
	return m.maybeRequiredGoogNamespaces
}

// WeaklyRequiredGoogNamespaces are the Closure namespaces this file weakly
// requires, i.e., arguments to goog.requireType calls.
//
// This is a multiset as it does not warn on duplicate namespaces, but will
// still encapsulate that information with this multiset.
//
// Java: ModuleMetadataMap.ModuleMetadata.weaklyRequiredGoogNamespaces()
func (m *ModuleMetadataMap_ModuleMetadata) WeaklyRequiredGoogNamespaces() *jutil.ImmutableMultiset[string] {
	return m.weaklyRequiredGoogNamespaces
}

// Es6ImportSpecifiers is the raw text of all ES6 import specifiers
// (includes "export from" as well).
//
// Java: ModuleMetadataMap.ModuleMetadata.es6ImportSpecifiers()
func (m *ModuleMetadataMap_ModuleMetadata) Es6ImportSpecifiers() *jutil.ImmutableMultiset[string] {
	return m.es6ImportSpecifiers
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
func (m *ModuleMetadataMap_ModuleMetadata) ReadToggles() *jutil.ImmutableMultiset[string] {
	return m.readToggles
}

// Path may return nil.
//
// Java: ModuleMetadataMap.ModuleMetadata.path()
func (m *ModuleMetadataMap_ModuleMetadata) Path() *ModuleLoader_ModulePath {
	return m.path
}

// Java: ModuleMetadataMap.ModuleMetadata.builder()
func ModuleMetadataMap_ModuleMetadata_Builder_() *ModuleMetadataMap_ModuleMetadata_Builder {
	return newAutoValue_ModuleMetadataMap_ModuleMetadata_Builder()
}

// Equals is reference equality.
//
// Java: ModuleMetadataMap.ModuleMetadata.equals(Object)
func (m *ModuleMetadataMap_ModuleMetadata) Equals(other any) bool {
	o, ok := other.(*ModuleMetadataMap_ModuleMetadata)
	return ok && o == m
}

// HashCode is the identity hash code.
//
// Java: ModuleMetadataMap.ModuleMetadata.hashCode()
func (m *ModuleMetadataMap_ModuleMetadata) HashCode() int32 {
	return jutil.IdentityHashCode(m)
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.toString()
func (m *ModuleMetadataMap_ModuleMetadata) String() string {
	rootNode := "null"
	if m.rootNode != nil {
		rootNode = m.rootNode.String()
	}
	path := "null"
	if m.path != nil {
		path = m.path.String()
	}
	return "ModuleMetadata{" +
		"moduleType=" + m.moduleType.String() + ", " +
		"rootNode=" + rootNode + ", " +
		"usesClosure=" + jutil.ToJavaString(m.usesClosure) + ", " +
		"isTestOnly=" + jutil.ToJavaString(m.isTestOnly) + ", " +
		"googNamespaces=" + m.googNamespaces.String() + ", " +
		"stronglyRequiredGoogNamespaces=" + m.stronglyRequiredGoogNamespaces.String() + ", " +
		"dynamicallyRequiredGoogNamespaces=" + m.dynamicallyRequiredGoogNamespaces.String() + ", " +
		"maybeRequiredGoogNamespaces=" + m.maybeRequiredGoogNamespaces.String() + ", " +
		"weaklyRequiredGoogNamespaces=" + m.weaklyRequiredGoogNamespaces.String() + ", " +
		"es6ImportSpecifiers=" + m.es6ImportSpecifiers.String() + ", " +
		"nestedModules=" + jutil.CollectionToString(m.nestedModules) + ", " +
		"readToggles=" + m.readToggles.String() + ", " +
		"path=" + path +
		"}"
}

// ModuleMetadataMap_ModuleMetadata_Builder is the builder for
// ModuleMetadata (AutoValue_ModuleMetadataMap_ModuleMetadata.Builder).
type ModuleMetadataMap_ModuleMetadata_Builder struct {
	moduleType                                *ModuleMetadataMap_ModuleType
	rootNode                                  *rhino.Node
	usesClosure                               bool
	isTestOnly                                bool
	googNamespacesBuilder_                    *jutil.ImmutableMultiset_Builder[string]
	googNamespaces                            *jutil.ImmutableMultiset[string]
	stronglyRequiredGoogNamespacesBuilder_    *jutil.ImmutableMultiset_Builder[string]
	stronglyRequiredGoogNamespaces            *jutil.ImmutableMultiset[string]
	dynamicallyRequiredGoogNamespacesBuilder_ *jutil.ImmutableMultiset_Builder[string]
	dynamicallyRequiredGoogNamespaces         *jutil.ImmutableMultiset[string]
	maybeRequiredGoogNamespacesBuilder_       *jutil.ImmutableMultiset_Builder[string]
	maybeRequiredGoogNamespaces               *jutil.ImmutableMultiset[string]
	weaklyRequiredGoogNamespacesBuilder_      *jutil.ImmutableMultiset_Builder[string]
	weaklyRequiredGoogNamespaces              *jutil.ImmutableMultiset[string]
	es6ImportSpecifiersBuilder_               *jutil.ImmutableMultiset_Builder[string]
	es6ImportSpecifiers                       *jutil.ImmutableMultiset[string]
	nestedModulesBuilder_                     *jutil.ImmutableList_Builder[*ModuleMetadataMap_ModuleMetadata]
	nestedModules                             []*ModuleMetadataMap_ModuleMetadata
	readTogglesBuilder_                       *jutil.ImmutableMultiset_Builder[string]
	readToggles                               *jutil.ImmutableMultiset[string]
	path                                      *ModuleLoader_ModulePath
	set0                                      byte
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder()
func newAutoValue_ModuleMetadataMap_ModuleMetadata_Builder() *ModuleMetadataMap_ModuleMetadata_Builder {
	return &ModuleMetadataMap_ModuleMetadata_Builder{}
}

// Java: ModuleMetadataMap.ModuleMetadata.Builder.addGoogNamespace(String)
func (b *ModuleMetadataMap_ModuleMetadata_Builder) AddGoogNamespace(namespace string) *ModuleMetadataMap_ModuleMetadata_Builder {
	b.GoogNamespacesBuilder().Add(namespace)
	return b
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.moduleType(ModuleType)
func (b *ModuleMetadataMap_ModuleMetadata_Builder) ModuleType(moduleType ModuleMetadataMap_ModuleType) *ModuleMetadataMap_ModuleMetadata_Builder {
	b.moduleType = &moduleType
	return b
}

// GetModuleType is the getter overload moduleType().
//
// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.moduleType()
func (b *ModuleMetadataMap_ModuleMetadata_Builder) GetModuleType() ModuleMetadataMap_ModuleType {
	if b.moduleType == nil {
		panic(jutil.NewIllegalStateExceptionMsg("Property \"moduleType\" has not been set"))
	}
	return *b.moduleType
}

// RootNode takes a possibly-nil root.
//
// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.rootNode(Node)
func (b *ModuleMetadataMap_ModuleMetadata_Builder) RootNode(rootNode *rhino.Node) *ModuleMetadataMap_ModuleMetadata_Builder {
	b.rootNode = rootNode
	return b
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.usesClosure(boolean)
func (b *ModuleMetadataMap_ModuleMetadata_Builder) UsesClosure(usesClosure bool) *ModuleMetadataMap_ModuleMetadata_Builder {
	b.usesClosure = usesClosure
	b.set0 |= 1
	return b
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.isTestOnly(boolean)
func (b *ModuleMetadataMap_ModuleMetadata_Builder) IsTestOnly(isTestOnly bool) *ModuleMetadataMap_ModuleMetadata_Builder {
	b.isTestOnly = isTestOnly
	b.set0 |= 2
	return b
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.googNamespacesBuilder()
func (b *ModuleMetadataMap_ModuleMetadata_Builder) GoogNamespacesBuilder() *jutil.ImmutableMultiset_Builder[string] {
	if b.googNamespacesBuilder_ == nil {
		b.googNamespacesBuilder_ = jutil.NewImmutableMultiset_Builder[string]()
	}
	return b.googNamespacesBuilder_
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.stronglyRequiredGoogNamespacesBuilder()
func (b *ModuleMetadataMap_ModuleMetadata_Builder) StronglyRequiredGoogNamespacesBuilder() *jutil.ImmutableMultiset_Builder[string] {
	if b.stronglyRequiredGoogNamespacesBuilder_ == nil {
		b.stronglyRequiredGoogNamespacesBuilder_ = jutil.NewImmutableMultiset_Builder[string]()
	}
	return b.stronglyRequiredGoogNamespacesBuilder_
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.dynamicallyRequiredGoogNamespacesBuilder()
func (b *ModuleMetadataMap_ModuleMetadata_Builder) DynamicallyRequiredGoogNamespacesBuilder() *jutil.ImmutableMultiset_Builder[string] {
	if b.dynamicallyRequiredGoogNamespacesBuilder_ == nil {
		b.dynamicallyRequiredGoogNamespacesBuilder_ = jutil.NewImmutableMultiset_Builder[string]()
	}
	return b.dynamicallyRequiredGoogNamespacesBuilder_
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.maybeRequiredGoogNamespacesBuilder()
func (b *ModuleMetadataMap_ModuleMetadata_Builder) MaybeRequiredGoogNamespacesBuilder() *jutil.ImmutableMultiset_Builder[string] {
	if b.maybeRequiredGoogNamespacesBuilder_ == nil {
		b.maybeRequiredGoogNamespacesBuilder_ = jutil.NewImmutableMultiset_Builder[string]()
	}
	return b.maybeRequiredGoogNamespacesBuilder_
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.weaklyRequiredGoogNamespacesBuilder()
func (b *ModuleMetadataMap_ModuleMetadata_Builder) WeaklyRequiredGoogNamespacesBuilder() *jutil.ImmutableMultiset_Builder[string] {
	if b.weaklyRequiredGoogNamespacesBuilder_ == nil {
		b.weaklyRequiredGoogNamespacesBuilder_ = jutil.NewImmutableMultiset_Builder[string]()
	}
	return b.weaklyRequiredGoogNamespacesBuilder_
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.es6ImportSpecifiersBuilder()
func (b *ModuleMetadataMap_ModuleMetadata_Builder) Es6ImportSpecifiersBuilder() *jutil.ImmutableMultiset_Builder[string] {
	if b.es6ImportSpecifiersBuilder_ == nil {
		b.es6ImportSpecifiersBuilder_ = jutil.NewImmutableMultiset_Builder[string]()
	}
	return b.es6ImportSpecifiersBuilder_
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.nestedModulesBuilder()
func (b *ModuleMetadataMap_ModuleMetadata_Builder) NestedModulesBuilder() *jutil.ImmutableList_Builder[*ModuleMetadataMap_ModuleMetadata] {
	if b.nestedModulesBuilder_ == nil {
		b.nestedModulesBuilder_ = jutil.ImmutableList_Builder_[*ModuleMetadataMap_ModuleMetadata]()
	}
	return b.nestedModulesBuilder_
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.readTogglesBuilder()
func (b *ModuleMetadataMap_ModuleMetadata_Builder) ReadTogglesBuilder() *jutil.ImmutableMultiset_Builder[string] {
	if b.readTogglesBuilder_ == nil {
		b.readTogglesBuilder_ = jutil.NewImmutableMultiset_Builder[string]()
	}
	return b.readTogglesBuilder_
}

// Path takes a possibly-nil path.
//
// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.path(ModulePath)
func (b *ModuleMetadataMap_ModuleMetadata_Builder) Path(path *ModuleLoader_ModulePath) *ModuleMetadataMap_ModuleMetadata_Builder {
	b.path = path
	return b
}

// Java: AutoValue_ModuleMetadataMap_ModuleMetadata.Builder.build()
func (b *ModuleMetadataMap_ModuleMetadata_Builder) Build() *ModuleMetadataMap_ModuleMetadata {
	if b.googNamespacesBuilder_ != nil {
		b.googNamespaces = b.googNamespacesBuilder_.Build()
	} else if b.googNamespaces == nil {
		b.googNamespaces = jutil.ImmutableMultiset_Of[string]()
	}
	if b.stronglyRequiredGoogNamespacesBuilder_ != nil {
		b.stronglyRequiredGoogNamespaces = b.stronglyRequiredGoogNamespacesBuilder_.Build()
	} else if b.stronglyRequiredGoogNamespaces == nil {
		b.stronglyRequiredGoogNamespaces = jutil.ImmutableMultiset_Of[string]()
	}
	if b.dynamicallyRequiredGoogNamespacesBuilder_ != nil {
		b.dynamicallyRequiredGoogNamespaces = b.dynamicallyRequiredGoogNamespacesBuilder_.Build()
	} else if b.dynamicallyRequiredGoogNamespaces == nil {
		b.dynamicallyRequiredGoogNamespaces = jutil.ImmutableMultiset_Of[string]()
	}
	if b.maybeRequiredGoogNamespacesBuilder_ != nil {
		b.maybeRequiredGoogNamespaces = b.maybeRequiredGoogNamespacesBuilder_.Build()
	} else if b.maybeRequiredGoogNamespaces == nil {
		b.maybeRequiredGoogNamespaces = jutil.ImmutableMultiset_Of[string]()
	}
	if b.weaklyRequiredGoogNamespacesBuilder_ != nil {
		b.weaklyRequiredGoogNamespaces = b.weaklyRequiredGoogNamespacesBuilder_.Build()
	} else if b.weaklyRequiredGoogNamespaces == nil {
		b.weaklyRequiredGoogNamespaces = jutil.ImmutableMultiset_Of[string]()
	}
	if b.es6ImportSpecifiersBuilder_ != nil {
		b.es6ImportSpecifiers = b.es6ImportSpecifiersBuilder_.Build()
	} else if b.es6ImportSpecifiers == nil {
		b.es6ImportSpecifiers = jutil.ImmutableMultiset_Of[string]()
	}
	if b.nestedModulesBuilder_ != nil {
		b.nestedModules = b.nestedModulesBuilder_.Build()
	} else if b.nestedModules == nil {
		b.nestedModules = []*ModuleMetadataMap_ModuleMetadata{}
	}
	if b.readTogglesBuilder_ != nil {
		b.readToggles = b.readTogglesBuilder_.Build()
	} else if b.readToggles == nil {
		b.readToggles = jutil.ImmutableMultiset_Of[string]()
	}
	if b.set0 != 3 ||
		b.moduleType == nil {
		missing := ""
		if b.moduleType == nil {
			missing += " moduleType"
		}
		if (b.set0 & 1) == 0 {
			missing += " usesClosure"
		}
		if (b.set0 & 2) == 0 {
			missing += " isTestOnly"
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing))
	}
	return &ModuleMetadataMap_ModuleMetadata{
		moduleType:                        *b.moduleType,
		rootNode:                          b.rootNode,
		usesClosure:                       b.usesClosure,
		isTestOnly:                        b.isTestOnly,
		googNamespaces:                    b.googNamespaces,
		stronglyRequiredGoogNamespaces:    b.stronglyRequiredGoogNamespaces,
		dynamicallyRequiredGoogNamespaces: b.dynamicallyRequiredGoogNamespaces,
		maybeRequiredGoogNamespaces:       b.maybeRequiredGoogNamespaces,
		weaklyRequiredGoogNamespaces:      b.weaklyRequiredGoogNamespaces,
		es6ImportSpecifiers:               b.es6ImportSpecifiers,
		nestedModules:                     b.nestedModules,
		readToggles:                       b.readToggles,
		path:                              b.path,
	}
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
