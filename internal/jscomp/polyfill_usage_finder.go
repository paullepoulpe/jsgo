// Port of com/google/javascript/jscomp/PolyfillUsageFinder.java @ closure-compiler v20261021

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
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/resources"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	pufEmptyID      = uint32(strpool.Intern(""))
	pufGoogID       = uint32(strpool.Intern("goog"))
	pufGlobalID     = uint32(strpool.Intern("global"))
	pufWindowID     = uint32(strpool.Intern("window"))
	pufGoogGlobalID = uint32(strpool.Intern("goog$global"))
	pufGlobalThisID = uint32(strpool.Intern("globalThis"))
)

type pufIDSet struct {
	bits []uint64
	high map[uint32]struct{}
}

func (s *pufIDSet) add(id uint32) {
	if id < 131072 {
		word := int(id >> 6)
		if word >= len(s.bits) {
			grown := make([]uint64, max(word+1, len(s.bits)*2))
			copy(grown, s.bits)
			s.bits = grown
		}
		s.bits[word] |= uint64(1) << (id & 63)
		return
	}
	if s.high == nil {
		s.high = make(map[uint32]struct{})
	}
	s.high[id] = struct{}{}
}

func (s *pufIDSet) has(id uint32) bool {
	word := int(id >> 6)
	if word < len(s.bits) {
		return (s.bits[word] & (uint64(1) << (id & 63))) != 0
	}
	if s.high == nil {
		return false
	}
	_, ok := s.high[id]
	return ok
}

type pufMethodEntry struct {
	methodName   string
	dotGuardName string
	polyfills    []*PolyfillUsageFinder_Polyfill
}

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
	// Map of all polyfills, keyed by their native version (the first ECMAScript spec version in
	// which they are defined)
	nativeVersions  []string
	byNativeVersion map[string][]*PolyfillUsageFinder_Polyfill

	staticsMap     map[string]*PolyfillUsageFinder_Polyfill
	staticsOneByID map[uint32]*PolyfillUsageFinder_Polyfill
	staticsPairMap map[[2]uint32]*PolyfillUsageFinder_Polyfill
	staticNameIDs  pufIDSet
	staticRootIDs  pufIDSet
	suffixIDs      pufIDSet
	anyPropIDs     pufIDSet
	methodsByID    map[uint32]*pufMethodEntry
}

func (p *PolyfillUsageFinder_Polyfills) getStatic(name string) *PolyfillUsageFinder_Polyfill {
	return p.staticsMap[name]
}

func (p *PolyfillUsageFinder_Polyfills) getMethods(name string) []*PolyfillUsageFinder_Polyfill {
	id := strpool.Lookup(name)
	if id == 0 {
		return nil
	}
	if entry := p.methodsByID[uint32(id)]; entry != nil {
		return entry.polyfills
	}
	return nil
}

func (p *PolyfillUsageFinder_Polyfills) lookupStaticChain(reversedIDs []uint32) *PolyfillUsageFinder_Polyfill {
	switch len(reversedIDs) {
	case 1:
		return p.staticsOneByID[reversedIDs[0]]
	case 2:
		return p.staticsPairMap[[2]uint32{reversedIDs[1], reversedIDs[0]}]
	default:
		var sb strings.Builder
		for i := len(reversedIDs) - 1; i >= 0; i-- {
			if i < len(reversedIDs)-1 {
				sb.WriteByte('.')
			}
			sb.WriteString(strpool.String(strpool.ID(reversedIDs[i])))
		}
		return p.staticsMap[sb.String()]
	}
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
var defaultPolyfillsText = sync.OnceValue(func() string {
	return resources.ResourceLoader_LoadTextResource("com.google.javascript.jscomp.PolyfillUsageFinder", "js/polyfills.txt")
})

var defaultPolyfillsTable = sync.OnceValue(func() *PolyfillUsageFinder_Polyfills {
	return polyfillUsageFinder_Polyfills_parseTable(defaultPolyfillsText())
})

// Java: PolyfillUsageFinder.Polyfills.fromTable(String)
func PolyfillUsageFinder_Polyfills_FromTable(table string) *PolyfillUsageFinder_Polyfills {
	if table == defaultPolyfillsText() {
		return defaultPolyfillsTable()
	}
	return polyfillUsageFinder_Polyfills_parseTable(table)
}

func polyfillUsageFinder_Polyfills_parseTable(table string) *PolyfillUsageFinder_Polyfills {
	p := &PolyfillUsageFinder_Polyfills{
		byNativeVersion: make(map[string][]*PolyfillUsageFinder_Polyfill),
		staticsMap:      make(map[string]*PolyfillUsageFinder_Polyfill),
		staticsOneByID:  make(map[uint32]*PolyfillUsageFinder_Polyfill),
		staticsPairMap:  make(map[[2]uint32]*PolyfillUsageFinder_Polyfill),
		methodsByID:     make(map[uint32]*pufMethodEntry),
	}
	for _, line := range strings.Split(table, "\n") {
		if line == "" {
			continue
		}
		trimmed := strings.TrimFunc(line, func(r rune) bool { return r <= ' ' })
		if trimmed == "" {
			continue
		}
		var tokens []string
		for _, tok := range strings.Split(trimmed, " ") {
			if tok != "" {
				tokens = append(tokens, tok)
			}
		}
		if len(tokens) < 3 {
			panic(jutil.NewIllegalArgumentExceptionMsg("Invalid table: too few tokens on line: " + line))
		}
		symbol := tokens[0]
		protoIdx := strings.LastIndex(symbol, ".prototype.")
		isPrototypeMethod := protoIdx >= 0
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
			methodName := symbol[protoIdx+len(".prototype."):]
			methodID := uint32(strpool.Intern(methodName))
			entry := p.methodsByID[methodID]
			if entry == nil {
				entry = &pufMethodEntry{
					methodName:   methodName,
					dotGuardName: "." + methodName,
				}
				p.methodsByID[methodID] = entry
			}
			entry.polyfills = append(entry.polyfills, polyfill)
			p.anyPropIDs.add(methodID)
		} else {
			if _, exists := p.staticsMap[symbol]; exists {
				panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " + symbol))
			}
			p.staticsMap[symbol] = polyfill
			dot := strings.IndexByte(symbol, '.')
			if dot < 0 {
				symID := uint32(strpool.Intern(symbol))
				p.staticsOneByID[symID] = polyfill
				p.staticNameIDs.add(symID)
				p.suffixIDs.add(symID)
				p.anyPropIDs.add(symID)
			} else {
				var rootID uint32
				if dot > 0 {
					rootID = uint32(strpool.Intern(symbol[:dot]))
					p.staticRootIDs.add(rootID)
				}
				lastDot := strings.LastIndexByte(symbol, '.')
				suffixID := uint32(strpool.Intern(symbol[lastDot+1:]))
				p.suffixIDs.add(suffixID)
				p.anyPropIDs.add(suffixID)
				if dot > 0 && dot == lastDot {
					p.staticsPairMap[[2]uint32{rootID, suffixID}] = polyfill
				}
			}
		}
		if _, exists := p.byNativeVersion[nativeVersionStr]; !exists {
			p.nativeVersions = append(p.nativeVersions, nativeVersionStr)
		}
		p.byNativeVersion[nativeVersionStr] = append(p.byNativeVersion[nativeVersionStr], polyfill)
	}
	return p
}

// Java: PolyfillUsageFinder.Polyfills.getPolyfillsNewerThan(LanguageMode)
func (p *PolyfillUsageFinder_Polyfills) GetPolyfillsNewerThan(languageMode CompilerOptions_LanguageMode) []*PolyfillUsageFinder_Polyfill {
	featureSet := languageMode.ToFeatureSet()
	result := []*PolyfillUsageFinder_Polyfill{}

	for _, nativeVersionStr := range p.nativeVersions {
		polyfillNativeFeatureSet := PolyfillUsageFinder_GetPolyfillSupportedFeatureSet(nativeVersionStr)
		if !featureSet.Contains(polyfillNativeFeatureSet) {
			result = append(result, p.byNativeVersion[nativeVersionStr]...)
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
	node             rhino.Node
	name             string
	isExplicitGlobal bool
}

// Java: PolyfillUsageFinder.PolyfillUsage(Polyfill, Node, String, boolean)
func NewPolyfillUsageFinder_PolyfillUsage(
	polyfill *PolyfillUsageFinder_Polyfill, node rhino.Node, name string, isExplicitGlobal bool) *PolyfillUsageFinder_PolyfillUsage {
	jutil.CheckNotNullMsg(polyfill, "polyfill")
	rhino.CheckNotNullMsg(node, "node")
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
func (u *PolyfillUsageFinder_PolyfillUsage) Node() rhino.Node { return u.node }

// Java: PolyfillUsageFinder.PolyfillUsage.name()
func (u *PolyfillUsageFinder_PolyfillUsage) Name() string { return u.name }

// Java: PolyfillUsageFinder.PolyfillUsage.isExplicitGlobal()
func (u *PolyfillUsageFinder_PolyfillUsage) IsExplicitGlobal() bool { return u.isExplicitGlobal }

// Java: PolyfillUsageFinder.PolyfillUsage.createExplicit(Polyfill, Node, String)
func PolyfillUsageFinder_PolyfillUsage_CreateExplicit(polyfill *PolyfillUsageFinder_Polyfill, node rhino.Node, name string) *PolyfillUsageFinder_PolyfillUsage {
	return NewPolyfillUsageFinder_PolyfillUsage(polyfill, node, name /* isExplicitGlobal= */, true)
}

// Java: PolyfillUsageFinder.PolyfillUsage.createNonExplicit(Polyfill, Node, String)
func PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill *PolyfillUsageFinder_Polyfill, node rhino.Node, name string) *PolyfillUsageFinder_PolyfillUsage {
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
func (f *PolyfillUsageFinder) TraverseExcludingGuarded(root rhino.Node, polyfillConsumer func(*PolyfillUsageFinder_PolyfillUsage)) {
	sc := NewSyntacticScopeCreator(f.compiler)
	NodeTraversal_Builder_().
		SetCompiler(f.compiler).
		SetCallback(newPolyfillUsageFinder_Traverser(f, f.compiler, polyfillConsumer, polyfillUsageFinder_Guard_ONLY_UNGUARDED)).
		SetScopeCreator(sc).
		Traverse(root)
}

// TraverseIncludingGuarded passes all polyfill usages found, in postorder,
// to the given polyfillConsumer
//
// Includes polyfill usages that are behind a guard, like
// `if (Promise) return Promise.resolve('ok');`
//
// Java: PolyfillUsageFinder.traverseIncludingGuarded(Node, Consumer<PolyfillUsage>)
func (f *PolyfillUsageFinder) TraverseIncludingGuarded(root rhino.Node, polyfillConsumer func(*PolyfillUsageFinder_PolyfillUsage)) {
	sc := NewSyntacticScopeCreator(f.compiler)
	NodeTraversal_Builder_().
		SetCompiler(f.compiler).
		SetCallback(newPolyfillUsageFinder_Traverser(f, f.compiler, polyfillConsumer, polyfillUsageFinder_Guard_ALL)).
		SetScopeCreator(sc).
		Traverse(root)
}

// TraverseOnlyGuarded passes all polyfill usages found, in postorder, to
// the given polyfillConsumer
//
// Only includes polyfill usages that are behind a guard, like
// `if (Promise) return Promise.resolve('ok');`
//
// Java: PolyfillUsageFinder.traverseOnlyGuarded(Node, Consumer<PolyfillUsage>)
func (f *PolyfillUsageFinder) TraverseOnlyGuarded(root rhino.Node, polyfillConsumer func(*PolyfillUsageFinder_PolyfillUsage)) {
	sc := NewSyntacticScopeCreator(f.compiler)
	NodeTraversal_Builder_().
		SetCompiler(f.compiler).
		SetCallback(newPolyfillUsageFinder_Traverser(f, f.compiler, polyfillConsumer, polyfillUsageFinder_Guard_ONLY_GUARDED)).
		SetScopeCreator(sc).
		Traverse(root)
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

func (tr *polyfillUsageFinder_Traverser) ShouldTraverse(traversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if !parent.IsNil() && n.GetFirstChild().IsNil() && len(tr.registeredGuards) == 0 {
		switch n.GetToken() {
		case rhino.Token_NAME:
			if !tr.outer.polyfills.staticNameIDs.has(n.GetStringID()) {
				return false
			}
		case rhino.Token_RETURN, rhino.Token_THROW:
			if len(tr.installedGuards) == 0 {
				return false
			}
		default:
			return false
		}
	}
	return tr.GuardedCallbackBase.ShouldTraverse(traversal, n, parent)
}

// Java: PolyfillUsageFinder.Traverser.visitGuarded(NodeTraversal, Node, Node)
func (tr *polyfillUsageFinder_Traverser) VisitGuarded(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) {
	if tr.includeGuardedUsages == polyfillUsageFinder_Guard_ONLY_GUARDED && !tr.canPossiblyBeGuarded() {
		return
	}
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
func (tr *polyfillUsageFinder_Traverser) VisitName(traversal *NodeTraversal, nameNode rhino.Node) {
	polyfill := tr.outer.polyfills.staticsOneByID[nameNode.GetStringID()]
	if polyfill == nil {
		// no polyfill exists for this name
		return
	}
	name := polyfill.nativeSymbol
	if tr.outer.IsPolyfillVisibleInScope(traversal, name) &&
		tr.includeGuardedUsages.ShouldInclude(tr.canPossiblyBeGuarded() && tr.IsGuarded(name)) {
		tr.polyfillConsumer(PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill, nameNode, name))
	}
}

// Java: PolyfillUsageFinder.Traverser.visitGetPropChain(NodeTraversal, Node)
func (tr *polyfillUsageFinder_Traverser) VisitGetPropChain(traversal *NodeTraversal, getPropNode rhino.Node) {
	propID := getPropNode.GetStringID()
	if !tr.outer.polyfills.anyPropIDs.has(propID) {
		return
	}
	// First see if we have a usage that matches a full static polyfill name.
	// e.g. `Array.from` or `globalThis.Promise.allSettled`
	staticPolyfillUsage :=
		tr.outer.MaybeCreateStaticPolyfillUsageForGetPropChain(traversal, getPropNode)
	if staticPolyfillUsage != nil {
		if tr.includeGuardedUsages.ShouldInclude(tr.canPossiblyBeGuarded() && tr.IsGuarded(staticPolyfillUsage.Name())) {
			tr.polyfillConsumer(staticPolyfillUsage)
		}
	} else {
		// We don't have a static polyfill usage, but this could still be a reference to one of
		// several possible method polyfills.
		// e.g. `obj.includes(x)` could be a usage of `Array.prototype.includes` or
		// `String.prototype.includes`.
		methodEntry := tr.outer.polyfills.methodsByID[propID]
		// Note that we use ".foo" as the guard check for methods to keep them distinct in case
		// there is also a static "foo" polyfill.
		if methodEntry != nil &&
			tr.includeGuardedUsages.ShouldInclude(tr.canPossiblyBeGuarded() && tr.IsGuarded(methodEntry.dotGuardName)) {
			for _, polyfill := range methodEntry.polyfills {
				tr.polyfillConsumer(
					PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill, getPropNode, methodEntry.methodName))
			}
		}
	}
}

// Java: PolyfillUsageFinder.maybeCreateStaticPolyfillUsageForGetPropChain(NodeTraversal, Node)
func (f *PolyfillUsageFinder) MaybeCreateStaticPolyfillUsageForGetPropChain(
	traversal *NodeTraversal, getPropNode rhino.Node) *PolyfillUsageFinder_PolyfillUsage {
	jutil.CheckArgumentMsg1(getPropNode.IsGetProp() || getPropNode.IsOptChainGetProp(), "%s", getPropNode)
	lastID := getPropNode.GetStringID()
	if !f.polyfills.suffixIDs.has(lastID) {
		// Save execution time by bailing out early if the property name at the end of the chain
		// doesn't match any of the known polyfills.
		return nil
	}
	firstChild := getPropNode.GetFirstChild()
	if firstChild.IsName() {
		rootID := firstChild.GetStringID()
		if rootID == 0 || rootID == pufEmptyID {
			return nil
		}
		if rootID == pufWindowID || rootID == pufGoogGlobalID || rootID == pufGlobalThisID {
			if polyfill := f.polyfills.staticsOneByID[lastID]; polyfill != nil {
				return PolyfillUsageFinder_PolyfillUsage_CreateExplicit(polyfill, getPropNode, polyfill.nativeSymbol)
			}
			return nil
		}
		if !f.polyfills.staticRootIDs.has(rootID) {
			return nil
		}
		if polyfill := f.polyfills.staticsPairMap[[2]uint32{rootID, lastID}]; polyfill != nil &&
			f.isPolyfillVisibleInScopeNode(traversal, firstChild) {
			return PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill, getPropNode, polyfill.nativeSymbol)
		}
		return nil
	}
	if !firstChild.IsGetProp() && !firstChild.IsOptChainGetProp() {
		return nil
	}
	// NOTE: We are not using isQualifiedName() and getQualifiedName() here, because we want to
	// locate the owner node and also have this code work for optional chains.
	var reversedIDsBuf [8]uint32
	reversedIDs := reversedIDsBuf[:0]
	reversedIDs = append(reversedIDs, lastID)
	var ownerNode rhino.Node
	for ownerNode = firstChild; ownerNode.IsGetProp() || ownerNode.IsOptChainGetProp(); ownerNode = ownerNode.GetFirstChild() {
		reversedIDs = append(reversedIDs, ownerNode.GetStringID())
	}
	if !ownerNode.IsName() {
		// Static polyfills are always fully qualified names beginning with a NAME node.
		// e.g. `Array.from` or `globalThis.Promise`
		return nil
	}
	rootID := ownerNode.GetStringID()
	if rootID == 0 || rootID == pufEmptyID {
		return nil
	}
	// If the full name starts with a known global value, like `goog.global.` or `globalThis.`,
	// we must strip that off before matching with the known polyfill names. Also, the presence of
	// the explicit global value name means we don't have to check the scope for a shadowing
	// variable.
	if rootID == pufWindowID || rootID == pufGoogGlobalID || rootID == pufGlobalThisID {
		if polyfill := f.polyfills.lookupStaticChain(reversedIDs); polyfill != nil {
			return PolyfillUsageFinder_PolyfillUsage_CreateExplicit(polyfill, getPropNode, polyfill.nativeSymbol)
		}
		return nil
	}
	if rootID == pufGoogID && reversedIDs[len(reversedIDs)-1] == pufGlobalID {
		if polyfill := f.polyfills.lookupStaticChain(reversedIDs[:len(reversedIDs)-1]); polyfill != nil {
			return PolyfillUsageFinder_PolyfillUsage_CreateExplicit(polyfill, getPropNode, polyfill.nativeSymbol)
		}
		return nil
	}
	if !f.polyfills.staticRootIDs.has(rootID) {
		return nil
	}
	reversedIDs = append(reversedIDs, rootID)
	if polyfill := f.polyfills.lookupStaticChain(reversedIDs); polyfill != nil {
		// If we see a declaration of the name, then it is defined by the source code,
		// and we won't count it as a reference to our polyfill.
		// Checking the scope is relatively expensive, so we don't want to do it until
		// we've confirmed that this node looks like it could be a polyfill reference.
		if f.isPolyfillVisibleInScopeNode(traversal, ownerNode) {
			return PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill, getPropNode, polyfill.nativeSymbol)
		}
	}
	return nil
}

func (f *PolyfillUsageFinder) isPolyfillVisibleInScopeNode(t *NodeTraversal, rootNameNode rhino.Node) bool {
	if f.compiler.GetLifeCycleStage() == AbstractCompiler_LifeCycleStage_NORMALIZED {
		return true
	}
	return t.GetScope().GetVar(rootNameNode.GetString()) == nil
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

func pufFindGlobalPrefix(qualifiedName string) (string, bool) {
	for _, global := range polyfillUsageFinder_GLOBAL_NAMES {
		if strings.HasPrefix(qualifiedName, global) {
			return global, true
		}
	}
	return "", false
}

// Java: PolyfillUsageFinder.findGlobalPrefix(String)
func PolyfillUsageFinder_FindGlobalPrefix(qualifiedName string) *string {
	if prefix, ok := pufFindGlobalPrefix(qualifiedName); ok {
		return &prefix
	}
	return nil
}

// Java: PolyfillUsageFinder.GLOBAL_NAMES
var polyfillUsageFinder_GLOBAL_NAMES = []string{"goog.global.", "window.", "goog$global.", "globalThis."}
