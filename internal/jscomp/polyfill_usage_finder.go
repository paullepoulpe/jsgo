// Port of com/google/javascript/jscomp/PolyfillUsageFinder.java @ closure-compiler v20261015

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PolyfillUsageFinder detects all potential usages of polyfilled classes or
// methods
type PolyfillUsageFinder struct {
	compiler  AbstractCompiler
	polyfills *PolyfillUsageFinder_Polyfills
}

// PolyfillUsageFinder_Polyfill represents a single polyfill: specifically,
// for a native symbol, a set of native and polyfill versions, and a library
// to ensure is injected if the output version is less than the native
// version.
type PolyfillUsageFinder_Polyfill struct {
	// The full name of the polyfill, e.g `Map` or `String.prototype.includes`
	nativeSymbol string

	// The language version at (or above) which the native symbol is
	// available and sufficient. If the language out flag is at least as high
	// as nativeVersion then the polyfill is not needed. This string should
	// be one of those returned by FeatureSet.version().
	nativeVersion string

	// The required language version for the polyfill to work. This should
	// not be higher than nativeVersion, but may be the same in cases where
	// there is no polyfill provided. This is used to emit a warning if the
	// language out flag is too low. This string should be one of those
	// returned by FeatureSet.version().
	polyfillVersion string

	// Runtime library to inject for the polyfill, e.g. "es6/map".
	library string

	kind PolyfillUsageFinder_Polyfill_Kind
}

// PolyfillUsageFinder_Polyfill_Kind is Polyfill.Kind.
type PolyfillUsageFinder_Polyfill_Kind int32

const (
	PolyfillUsageFinder_Polyfill_Kind_STATIC PolyfillUsageFinder_Polyfill_Kind = iota // Map or Array.of
	PolyfillUsageFinder_Polyfill_Kind_METHOD                                          // String.prototype.includes
)

var polyfillUsageFinder_Polyfill_Kind_names = [...]string{"STATIC", "METHOD"}

// PolyfillUsageFinder_Polyfill_Kind_values is Kind.values().
func PolyfillUsageFinder_Polyfill_Kind_values() []PolyfillUsageFinder_Polyfill_Kind {
	return []PolyfillUsageFinder_Polyfill_Kind{PolyfillUsageFinder_Polyfill_Kind_STATIC, PolyfillUsageFinder_Polyfill_Kind_METHOD}
}

// String is Enum.name().
func (k PolyfillUsageFinder_Polyfill_Kind) String() string {
	return polyfillUsageFinder_Polyfill_Kind_names[k]
}

// Ordinal is Enum.ordinal().
func (k PolyfillUsageFinder_Polyfill_Kind) Ordinal() int { return int(k) }

// Java: PolyfillUsageFinder.Polyfill(String, String, String, String, Kind)
func NewPolyfillUsageFinder_Polyfill(
	nativeSymbol string,
	nativeVersion string,
	polyfillVersion string,
	library string,
	kind PolyfillUsageFinder_Polyfill_Kind) *PolyfillUsageFinder_Polyfill {
	p := &PolyfillUsageFinder_Polyfill{}
	p.nativeSymbol = nativeSymbol
	p.nativeVersion = nativeVersion
	p.polyfillVersion = polyfillVersion
	p.library = library
	p.kind = kind
	return p
}

// PolyfillUsageFinder_Polyfills maps from polyfill names to the actual
// Polyfill object.
type PolyfillUsageFinder_Polyfills struct {
	// Map of method polyfills, keyed by native method name.
	methods *jutil.ListMultimap[string, *PolyfillUsageFinder_Polyfill]
	// Map of static polyfills, keyed by fully-qualified native name.
	statics *jutil.LinkedHashMap[string, *PolyfillUsageFinder_Polyfill]
	// Set of suffixes of qualified names.
	suffixes *jutil.LinkedHashSet[string]
	// Map of all polyfills, keyed by their native version (the first ECMAScript spec version in
	// which they are defined)
	byNativeVersion *jutil.ListMultimap[string, *PolyfillUsageFinder_Polyfill]
}

// Java: PolyfillUsageFinder.Polyfills(ImmutableMultimap<String, Polyfill>, ImmutableMap<String, Polyfill>, ImmutableMultimap<String, Polyfill>)
func newPolyfillUsageFinder_Polyfills(
	methods *jutil.ListMultimap[string, *PolyfillUsageFinder_Polyfill],
	statics *jutil.LinkedHashMap[string, *PolyfillUsageFinder_Polyfill],
	byNativeVersion *jutil.ListMultimap[string, *PolyfillUsageFinder_Polyfill]) *PolyfillUsageFinder_Polyfills {
	p := &PolyfillUsageFinder_Polyfills{}
	p.methods = methods
	p.statics = statics
	p.byNativeVersion = byNativeVersion
	suffixes := jutil.NewLinkedHashSet[string]()
	for arg := range statics.Keys() {
		suffixes.Add(jutil.StringSubstringFrom(arg, jutil.StringLastIndexOf(arg, ".")+1))
	}
	p.suffixes = suffixes
	return p
}

// PolyfillUsageFinder_Polyfills_FromTable builds a Polyfills instance from a
// polyfill table, which is a simple text file with lines containing
// space-separated tokens:
//
//	[NATIVE_SYMBOL] [NATIVE_VERSION] [POLYFILL_VERSION] [LIBRARY]
//
// For example,
//
//	Array.prototype.fill es6 es3 es6/array/fill
//	Map es6 es3 es6/map
//	WeakMap es6 es6
//
// The last line, WeakMap, does not have a polyfill available, so the
// library token is empty.
//
// Java: PolyfillUsageFinder.Polyfills.fromTable(String)
func PolyfillUsageFinder_Polyfills_FromTable(table string) *PolyfillUsageFinder_Polyfills {
	methods := jutil.NewImmutableListMultimapBuilder[string, *PolyfillUsageFinder_Polyfill]()
	statics := jutil.NewImmutableMapBuilder[string, *PolyfillUsageFinder_Polyfill]()
	byNativeVersion := jutil.NewImmutableListMultimapBuilder[string, *PolyfillUsageFinder_Polyfill]()
	for _, line := range jutil.Splitter_OnCharOmitEmptyStrings(table, '\n') {
		tokens := jutil.Splitter_OnCharOmitEmptyStrings(jutil.StringTrim(line), ' ')
		if len(tokens) == 1 && tokens[0] == "" {
			continue
		} else if len(tokens) < 3 {
			panic(jutil.NewIllegalArgumentExceptionMsg("Invalid table: too few tokens on line: " + line))
		}
		symbol := tokens[0]
		isPrototypeMethod := strings.Contains(symbol, ".prototype.")
		nativeVersionStr := tokens[1]
		polyfillVersionStr := tokens[2]
		library := ""
		if len(tokens) > 3 {
			library = tokens[3]
		}
		kind := PolyfillUsageFinder_Polyfill_Kind_STATIC
		if isPrototypeMethod {
			kind = PolyfillUsageFinder_Polyfill_Kind_METHOD
		}
		polyfill :=
			NewPolyfillUsageFinder_Polyfill(
				symbol,
				nativeVersionStr,
				polyfillVersionStr,
				library,
				kind)
		if isPrototypeMethod {
			methods.Put(jutil.StringReplaceAll(symbol, ".*\\.prototype\\.", ""), polyfill)
		} else {
			statics.Put(symbol, polyfill)
		}
		byNativeVersion.Put(nativeVersionStr, polyfill)
	}
	return newPolyfillUsageFinder_Polyfills(methods.Build(), statics.BuildOrThrow(), byNativeVersion.Build())
}

// Java: PolyfillUsageFinder.Polyfills.getPolyfillsNewerThan(LanguageMode)
func (p *PolyfillUsageFinder_Polyfills) GetPolyfillsNewerThan(languageMode CompilerOptions_LanguageMode) []*PolyfillUsageFinder_Polyfill {
	featureSet := languageMode.ToFeatureSet()
	result := []*PolyfillUsageFinder_Polyfill{}

	for nativeVersionStr := range p.byNativeVersion.KeySet() {
		polyfillNativeFeatureSet := PolyfillUsageFinder_GetPolyfillSupportedFeatureSet(nativeVersionStr)
		if !featureSet.Contains(polyfillNativeFeatureSet) {
			result = append(result, p.byNativeVersion.Get(nativeVersionStr)...)
		}
	}
	return result
}

// PolyfillUsageFinder_GetPolyfillSupportedFeatureSet converts a polyfill
// native version string as passed to $jscomp.polyfill, e.g. "es6", to a
// FeatureSet.
//
// Java: PolyfillUsageFinder.getPolyfillSupportedFeatureSet(String)
func PolyfillUsageFinder_GetPolyfillSupportedFeatureSet(nativeVersionStr string) *parser.FeatureSet {
	// When calculating whether a polyfill should be included for a given language_out associated
	// with a BFSY, we want to also compare using a BFSY FeatureSet. Needed because BFSY FeatureSets
	// often exclude features but the ES20XX FeatureSets do not meaning we could unnecessarily
	// include polyfills that, in reality, could be excluded.
	return parser.FeatureSet_BrowserFeatureSetValueOf(nativeVersionStr).WithoutFeatures(parser.FeatureSet_Feature_MODULES)
}

// PolyfillUsageFinder_PolyfillUsage is the record PolyfillUsage(polyfill,
// node, name, isExplicitGlobal).
type PolyfillUsageFinder_PolyfillUsage struct {
	polyfill         *PolyfillUsageFinder_Polyfill
	node             *rhino.Node
	name             string
	isExplicitGlobal bool
}

// Java: PolyfillUsageFinder.PolyfillUsage(Polyfill, Node, String, boolean)
func NewPolyfillUsageFinder_PolyfillUsage(
	polyfill *PolyfillUsageFinder_Polyfill, node *rhino.Node, name string, isExplicitGlobal bool) *PolyfillUsageFinder_PolyfillUsage {
	jutil.CheckNotNullMsg(polyfill, "polyfill")
	jutil.CheckNotNullMsg(node, "node")
	// Go-only: language-adapter — requireNonNull(name): a Go string cannot
	// be null.
	return &PolyfillUsageFinder_PolyfillUsage{
		polyfill:         polyfill,
		node:             node,
		name:             name,
		isExplicitGlobal: isExplicitGlobal,
	}
}

// Java: PolyfillUsageFinder.PolyfillUsage.polyfill()
func (u *PolyfillUsageFinder_PolyfillUsage) Polyfill() *PolyfillUsageFinder_Polyfill {
	return u.polyfill
}

// Java: PolyfillUsageFinder.PolyfillUsage.node()
func (u *PolyfillUsageFinder_PolyfillUsage) Node() *rhino.Node { return u.node }

// Java: PolyfillUsageFinder.PolyfillUsage.name()
func (u *PolyfillUsageFinder_PolyfillUsage) Name() string { return u.name }

// Java: PolyfillUsageFinder.PolyfillUsage.isExplicitGlobal()
func (u *PolyfillUsageFinder_PolyfillUsage) IsExplicitGlobal() bool { return u.isExplicitGlobal }

// Java: PolyfillUsageFinder.PolyfillUsage.createExplicit(Polyfill, Node, String)
func PolyfillUsageFinder_PolyfillUsage_CreateExplicit(polyfill *PolyfillUsageFinder_Polyfill, node *rhino.Node, name string) *PolyfillUsageFinder_PolyfillUsage {
	return NewPolyfillUsageFinder_PolyfillUsage(polyfill, node, name /* isExplicitGlobal= */, true)
}

// Java: PolyfillUsageFinder.PolyfillUsage.createNonExplicit(Polyfill, Node, String)
func PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill *PolyfillUsageFinder_Polyfill, node *rhino.Node, name string) *PolyfillUsageFinder_PolyfillUsage {
	return NewPolyfillUsageFinder_PolyfillUsage(polyfill, node, name /* isExplicitGlobal= */, false)
}

// Java: PolyfillUsageFinder(AbstractCompiler, Polyfills)
func NewPolyfillUsageFinder(compiler AbstractCompiler, polyfills *PolyfillUsageFinder_Polyfills) *PolyfillUsageFinder {
	f := &PolyfillUsageFinder{}
	f.polyfills = polyfills
	f.compiler = compiler
	return f
}

// TraverseExcludingGuarded passes all polyfill usages found, in postorder,
// to the given polyfillConsumer
//
// Excludes polyfill usages behind a guard, like
// `if (Promise) return Promise.resolve('ok');`
//
// Java: PolyfillUsageFinder.traverseExcludingGuarded(Node, Consumer<PolyfillUsage>)
func (f *PolyfillUsageFinder) TraverseExcludingGuarded(root *rhino.Node, polyfillConsumer func(*PolyfillUsageFinder_PolyfillUsage)) {
	NodeTraversal_Traverse(
		f.compiler, root, newPolyfillUsageFinder_Traverser(f, f.compiler, polyfillConsumer, polyfillUsageFinder_Guard_ONLY_UNGUARDED))
}

// TraverseIncludingGuarded passes all polyfill usages found, in postorder,
// to the given polyfillConsumer
//
// Includes polyfill usages that are behind a guard, like
// `if (Promise) return Promise.resolve('ok');`
//
// Java: PolyfillUsageFinder.traverseIncludingGuarded(Node, Consumer<PolyfillUsage>)
func (f *PolyfillUsageFinder) TraverseIncludingGuarded(root *rhino.Node, polyfillConsumer func(*PolyfillUsageFinder_PolyfillUsage)) {
	NodeTraversal_Traverse(
		f.compiler, root, newPolyfillUsageFinder_Traverser(f, f.compiler, polyfillConsumer, polyfillUsageFinder_Guard_ALL))
}

// TraverseOnlyGuarded passes all polyfill usages found, in postorder, to
// the given polyfillConsumer
//
// Only includes polyfill usages that are behind a guard, like
// `if (Promise) return Promise.resolve('ok');`
//
// Java: PolyfillUsageFinder.traverseOnlyGuarded(Node, Consumer<PolyfillUsage>)
func (f *PolyfillUsageFinder) TraverseOnlyGuarded(root *rhino.Node, polyfillConsumer func(*PolyfillUsageFinder_PolyfillUsage)) {
	NodeTraversal_Traverse(
		f.compiler, root, newPolyfillUsageFinder_Traverser(f, f.compiler, polyfillConsumer, polyfillUsageFinder_Guard_ONLY_GUARDED))
}

// polyfillUsageFinder_Guard is the private enum PolyfillUsageFinder.Guard.
type polyfillUsageFinder_Guard int32

const (
	polyfillUsageFinder_Guard_ONLY_GUARDED polyfillUsageFinder_Guard = iota
	polyfillUsageFinder_Guard_ONLY_UNGUARDED
	polyfillUsageFinder_Guard_ALL
)

var polyfillUsageFinder_Guard_names = [...]string{"ONLY_GUARDED", "ONLY_UNGUARDED", "ALL"}

// String is Enum.name().
func (g polyfillUsageFinder_Guard) String() string { return polyfillUsageFinder_Guard_names[g] }

// Java: PolyfillUsageFinder.Guard.shouldInclude(boolean)
func (g polyfillUsageFinder_Guard) ShouldInclude(isGuarded bool) bool {
	switch g {
	case polyfillUsageFinder_Guard_ALL:
		return true
	case polyfillUsageFinder_Guard_ONLY_GUARDED:
		return isGuarded
	case polyfillUsageFinder_Guard_ONLY_UNGUARDED:
		return !isGuarded
	}
	panic(jutil.NewIllegalStateExceptionMsg(g.String()))
}

// Java: private class PolyfillUsageFinder.Traverser extends GuardedCallback<String>
type polyfillUsageFinder_Traverser struct {
	GuardedCallbackBase[string]
	outer *PolyfillUsageFinder

	polyfillConsumer func(*PolyfillUsageFinder_PolyfillUsage)
	// Whether to emit usages like Promise in `if (Promise) return Promise.resolve('ok');}`
	includeGuardedUsages polyfillUsageFinder_Guard
}

// Java: PolyfillUsageFinder.Traverser(AbstractCompiler, Consumer<PolyfillUsage>, Guard)
func newPolyfillUsageFinder_Traverser(
	outer *PolyfillUsageFinder,
	compiler AbstractCompiler,
	polyfillConsumer func(*PolyfillUsageFinder_PolyfillUsage),
	includeGuardedUsages polyfillUsageFinder_Guard) *polyfillUsageFinder_Traverser {
	tr := &polyfillUsageFinder_Traverser{outer: outer}
	tr.InitGuardedCallbackBase(tr, compiler)
	tr.polyfillConsumer = polyfillConsumer
	tr.includeGuardedUsages = includeGuardedUsages
	return tr
}

// Java: PolyfillUsageFinder.Traverser.visitGuarded(NodeTraversal, Node, Node)
func (tr *polyfillUsageFinder_Traverser) VisitGuarded(traversal *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	switch node.GetToken() {
	case rhino.Token_NAME:
		tr.VisitName(traversal, node)
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		tr.VisitGetPropChain(traversal, node)
	default:
		// nothing to do
	}
}

// Java: PolyfillUsageFinder.Traverser.visitName(NodeTraversal, Node)
func (tr *polyfillUsageFinder_Traverser) VisitName(traversal *NodeTraversal, nameNode *rhino.Node) {
	name := nameNode.GetString()
	polyfill := tr.outer.polyfills.statics.Get(name)
	if polyfill == nil {
		// no polyfill exists for this name
		return
	}

	if tr.outer.IsPolyfillVisibleInScope(traversal, name) &&
		tr.includeGuardedUsages.ShouldInclude(tr.IsGuarded(name)) {
		tr.polyfillConsumer(PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill, nameNode, name))
	}
}

// Java: PolyfillUsageFinder.Traverser.visitGetPropChain(NodeTraversal, Node)
func (tr *polyfillUsageFinder_Traverser) VisitGetPropChain(traversal *NodeTraversal, getPropNode *rhino.Node) {
	// First see if we have a usage that matches a full static polyfill name.
	// e.g. `Array.from` or `globalThis.Promise.allSettled`
	staticPolyfillUsage :=
		tr.outer.MaybeCreateStaticPolyfillUsageForGetPropChain(traversal, getPropNode)
	if staticPolyfillUsage != nil {
		if tr.includeGuardedUsages.ShouldInclude(tr.IsGuarded(staticPolyfillUsage.Name())) {
			tr.polyfillConsumer(staticPolyfillUsage)
		}
	} else {
		// We don't have a static polyfill usage, but this could still be a reference to one of
		// several possible method polyfills.
		// e.g. `obj.includes(x)` could be a usage of `Array.prototype.includes` or
		// `String.prototype.includes`.
		propertyName := getPropNode.GetString()
		methodPolyfills := tr.outer.polyfills.methods.Get(propertyName)
		// Note that we use ".foo" as the guard check for methods to keep them distinct in case
		// there is also a static "foo" polyfill.
		if len(methodPolyfills) != 0 &&
			tr.includeGuardedUsages.ShouldInclude(tr.IsGuarded("."+propertyName)) {
			for _, polyfill := range methodPolyfills {
				tr.polyfillConsumer(
					PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill, getPropNode, propertyName))
			}
		}
	}
}

// Java: PolyfillUsageFinder.maybeCreateStaticPolyfillUsageForGetPropChain(NodeTraversal, Node)
func (f *PolyfillUsageFinder) MaybeCreateStaticPolyfillUsageForGetPropChain(
	traversal *NodeTraversal, getPropNode *rhino.Node) *PolyfillUsageFinder_PolyfillUsage {
	jutil.CheckArgumentMsg1(getPropNode.IsGetProp() || getPropNode.IsOptChainGetProp(), "%s", getPropNode)
	lastComponent := getPropNode.GetString()
	if !f.polyfills.suffixes.Contains(lastComponent) {
		// Save execution time by bailing out early if the property name at the end of the chain
		// doesn't match any of the known polyfills.
		return nil
	}
	// NOTE: We are not using isQualifiedName() and getQualifiedName() here, because we want to
	// locate the owner node and also have this code work for optional chains.
	//
	// Go-only: language-adapter — the ArrayDeque filled with addFirst is a
	// slice filled from the end and reversed when joined.
	reversedComponents := []string{}
	reversedComponents = append(reversedComponents, lastComponent)
	var ownerNode *rhino.Node
	for ownerNode = getPropNode.GetFirstChild(); ownerNode.IsGetProp() || ownerNode.IsOptChainGetProp(); ownerNode = ownerNode.GetFirstChild() {
		reversedComponents = append(reversedComponents, ownerNode.GetString())
	}
	if !ownerNode.IsName() {
		// Static polyfills are always fully qualified names beginning with a NAME node.
		// e.g. `Array.from` or `globalThis.Promise`
		return nil
	}
	rootName := ownerNode.GetString()
	reversedComponents = append(reversedComponents, rootName)
	var sb strings.Builder
	for i := len(reversedComponents) - 1; i >= 0; i-- {
		sb.WriteString(reversedComponents[i])
		if i > 0 {
			sb.WriteByte('.')
		}
	}
	fullName := sb.String()
	globalPrefix := PolyfillUsageFinder_FindGlobalPrefix(fullName)
	if globalPrefix != nil {
		// The full name starts with a known global value, like `goog.global.` or `globalThis.`.
		// We must strip that off before matching with the known polyfill names.
		// (Note that the connecting '.' is included and will also be stripped.)
		// Also, the presence of the explicit global value name means we don't have to check the
		// scope for a shadowing variable as we do below.
		polyfill := f.polyfills.statics.Get(jutil.StringSubstringFrom(fullName, jutil.StringLength(*globalPrefix)))
		if polyfill != nil {
			return PolyfillUsageFinder_PolyfillUsage_CreateExplicit(polyfill, getPropNode, polyfill.nativeSymbol)
		}
	} else {
		polyfill := f.polyfills.statics.Get(fullName)
		if polyfill != nil {
			// If we see a declaration of the name, then it is defined by the source code,
			// and we won't count it as a reference to our polyfill.
			// Checking the scope is relatively expensive, so we don't want to do it until
			// we've confirmed that this node looks like it could be a polyfill reference.
			if f.IsPolyfillVisibleInScope(traversal, rootName) {
				return PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill, getPropNode, polyfill.nativeSymbol)
			}
		}
	}
	return nil
}

// IsPolyfillVisibleInScope returns whether the Polyfill name is visible in
// the current scope.
//
// Java: PolyfillUsageFinder.isPolyfillVisibleInScope(NodeTraversal, String)
func (f *PolyfillUsageFinder) IsPolyfillVisibleInScope(t *NodeTraversal, globalName string) bool {
	// This class is only supposed to traverse over actual sources, not externs,
	// so it shouldn't see the declarations of the things being polyfilled.

	// If the AST is normalized we can avoid building the scope (which is relatively expensive)
	// and rely on the normalization guarantee that every name in AST is unique.

	// NOTE: there is an edge case that we are purposely ignoring here.  This pass expects
	// the traverse the AST without externs, this creates a distinction between undeclared variables
	// (externs) and symbols in the global scope (internal globals).  So then against the
	// un-normalized AST for detecting whether to inject Polyfills, it won't inject,  if there is a
	// global declaration of the same name.  But when removing Polyfills (which happens when the
	// AST is normalized), it won't try to remove these same Polyfills if they were present,
	// because the scope distinction isn't used.  This doesn't change behavior because they
	// wouldn't have been injected in the first place.

	return f.compiler.GetLifeCycleStage() == AbstractCompiler_LifeCycleStage_NORMALIZED ||
		t.GetScope().GetVar(globalName) == nil
}

// Java: PolyfillUsageFinder.findGlobalPrefix(String)
func PolyfillUsageFinder_FindGlobalPrefix(qualifiedName string) *string {
	for _, global := range polyfillUsageFinder_GLOBAL_NAMES {
		if strings.HasPrefix(qualifiedName, global) {
			return &global
		}
	}
	return nil
}

// Java: PolyfillUsageFinder.GLOBAL_NAMES
var polyfillUsageFinder_GLOBAL_NAMES = []string{"goog.global.", "window.", "goog$global.", "globalThis."}
