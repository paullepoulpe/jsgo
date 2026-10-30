// Port of com/google/javascript/jscomp/deps/SimpleDependencyInfo.java @ closure-compiler v20261021
// (with the AutoValue-generated deps/AutoValue_SimpleDependencyInfo.java)

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

package jscomp

import (
	"slices"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	simpleDepModuleNone uint8 = iota
	simpleDepModuleGoog
	simpleDepModuleES6
)

// SimpleDependencyInfo holds JS dependency information for a single .js file.
type SimpleDependencyInfo struct {
	name                      string
	pathRelativeToClosureBase string
	provides                  []string
	requires                  []*DependencyInfo_Require
	typeRequires              []string
	loadFlags                 LoadFlags
	moduleKind                uint8
	hasExternsAnnotation      bool
	hasNoCompileAnnotation    bool
}

var _ DependencyInfo = (*SimpleDependencyInfo)(nil)

type canonicalLangFlags struct {
	only [1]LoadFlag
	goog [2]LoadFlag
	es6  [2]LoadFlag
}

func makeCanonicalLangFlags(lang string) canonicalLangFlags {
	return canonicalLangFlags{
		only: [1]LoadFlag{{Key: "lang", Value: lang}},
		goog: [2]LoadFlag{{Key: "lang", Value: lang}, {Key: "module", Value: "goog"}},
		es6:  [2]LoadFlag{{Key: "lang", Value: lang}, {Key: "module", Value: "es6"}},
	}
}

var (
	simpleDependencyInfo_emptyStrings  = []string{}
	simpleDependencyInfo_emptyRequires = []*DependencyInfo_Require{}

	canonFlagsES5     = makeCanonicalLangFlags("es5")
	canonFlagsES6     = makeCanonicalLangFlags("es6")
	canonFlagsES7     = makeCanonicalLangFlags("es7")
	canonFlagsES8     = makeCanonicalLangFlags("es8")
	canonFlagsES9     = makeCanonicalLangFlags("es9")
	canonFlagsES2020  = makeCanonicalLangFlags("es_2020")
	canonFlagsES2021  = makeCanonicalLangFlags("es_2021")
	canonFlagsESNext  = makeCanonicalLangFlags("es_next")
	canonFlagsTSUnsup = makeCanonicalLangFlags("ts_unsupported")
)

func lookupCanonicalLangFlags(lang string) *canonicalLangFlags {
	switch lang {
	case "es5":
		return &canonFlagsES5
	case "es6":
		return &canonFlagsES6
	case "es7":
		return &canonFlagsES7
	case "es8":
		return &canonFlagsES8
	case "es9":
		return &canonFlagsES9
	case "es_2020":
		return &canonFlagsES2020
	case "es_2021":
		return &canonFlagsES2021
	case "es_next":
		return &canonFlagsESNext
	case "ts_unsupported":
		return &canonFlagsTSUnsup
	default:
		return nil
	}
}

func simpleDependencyInfo_moduleKindOf(mod string) uint8 {
	switch mod {
	case "goog":
		return simpleDepModuleGoog
	case "es6":
		return simpleDepModuleES6
	default:
		return simpleDepModuleNone
	}
}

func normalizeLoadFlagsEntries(entries []LoadFlag, cloneCustom bool) (LoadFlags, uint8) {
	switch len(entries) {
	case 0:
		return EmptyLoadFlags, simpleDepModuleNone
	case 1:
		e := entries[0]
		switch e.Key {
		case "module":
			if e.Value == "goog" {
				return GoogModuleLoadFlags, simpleDepModuleGoog
			}
			if e.Value == "es6" {
				return ES6ModuleLoadFlags, simpleDepModuleES6
			}
		case "lang":
			if c := lookupCanonicalLangFlags(e.Value); c != nil {
				return LoadFlags{entries: c.only[:]}, simpleDepModuleNone
			}
		}
	case 2:
		if entries[0].Key == "lang" && entries[1].Key == "module" {
			if c := lookupCanonicalLangFlags(entries[0].Value); c != nil {
				switch entries[1].Value {
				case "goog":
					return LoadFlags{entries: c.goog[:]}, simpleDepModuleGoog
				case "es6":
					return LoadFlags{entries: c.es6[:]}, simpleDepModuleES6
				}
			}
		}
	}
	var mod string
	for _, e := range entries {
		if e.Key == "module" {
			mod = e.Value
			break
		}
	}
	if cloneCustom {
		entries = slices.Clone(entries)
	}
	return LoadFlags{entries: entries}, simpleDependencyInfo_moduleKindOf(mod)
}

func simpleDependencyInfo_flagsForLangModule(lang, mod string) LoadFlags {
	if lang == "" {
		switch mod {
		case "":
			return EmptyLoadFlags
		case "goog":
			return GoogModuleLoadFlags
		case "es6":
			return ES6ModuleLoadFlags
		default:
			return LoadFlags{entries: []LoadFlag{{Key: "module", Value: mod}}}
		}
	}
	if c := lookupCanonicalLangFlags(lang); c != nil {
		switch mod {
		case "":
			return LoadFlags{entries: c.only[:]}
		case "goog":
			return LoadFlags{entries: c.goog[:]}
		case "es6":
			return LoadFlags{entries: c.es6[:]}
		}
	}
	if mod == "" {
		return LoadFlags{entries: []LoadFlag{{Key: "lang", Value: lang}}}
	}
	return LoadFlags{entries: []LoadFlag{{Key: "lang", Value: lang}, {Key: "module", Value: mod}}}
}

func newSimpleDependencyInfoDirect(
	srcPathRelativeToClosure string,
	pathOfDefiningFile string,
	provides []string,
	requires []*DependencyInfo_Require,
	typeRequires []string,
	loadFlags LoadFlags,
	hasExternsAnnotation bool,
	hasNoCompileAnnotation bool,
) *SimpleDependencyInfo {
	provCopy, typeReqCopy := simpleDependencyInfo_emptyStrings, simpleDependencyInfo_emptyStrings
	np, nt := len(provides), len(typeRequires)
	if np > 0 || nt > 0 {
		combined := make([]string, np+nt)
		copy(combined[:np], provides)
		copy(combined[np:], typeRequires)
		if np > 0 {
			provCopy = combined[:np:np]
		}
		if nt > 0 {
			typeReqCopy = combined[np:]
		}
	}
	normFlags, modKind := normalizeLoadFlagsEntries(loadFlags.entries, false)
	return &SimpleDependencyInfo{
		name:                      pathOfDefiningFile,
		pathRelativeToClosureBase: srcPathRelativeToClosure,
		provides:                  provCopy,
		requires:                  simpleDependencyInfo_copyRequires(requires),
		typeRequires:              typeReqCopy,
		loadFlags:                 normFlags,
		moduleKind:                modKind,
		hasExternsAnnotation:      hasExternsAnnotation,
		hasNoCompileAnnotation:    hasNoCompileAnnotation,
	}
}

// Java: AutoValue_SimpleDependencyInfo.getName()
func (d *SimpleDependencyInfo) GetName() string { return d.name }

// Java: AutoValue_SimpleDependencyInfo.getPathRelativeToClosureBase()
func (d *SimpleDependencyInfo) GetPathRelativeToClosureBase() string {
	return d.pathRelativeToClosureBase
}

// Java: AutoValue_SimpleDependencyInfo.getProvides()
func (d *SimpleDependencyInfo) GetProvides() []string { return d.provides }

// Java: AutoValue_SimpleDependencyInfo.getRequires()
func (d *SimpleDependencyInfo) GetRequires() []*DependencyInfo_Require { return d.requires }

// Java: AutoValue_SimpleDependencyInfo.getTypeRequires()
func (d *SimpleDependencyInfo) GetTypeRequires() []string { return d.typeRequires }

// Java: AutoValue_SimpleDependencyInfo.getLoadFlags()
func (d *SimpleDependencyInfo) GetLoadFlags() LoadFlags {
	return d.loadFlags
}

// Java: AutoValue_SimpleDependencyInfo.getHasExternsAnnotation()
func (d *SimpleDependencyInfo) GetHasExternsAnnotation() bool { return d.hasExternsAnnotation }

// Java: AutoValue_SimpleDependencyInfo.getHasNoCompileAnnotation()
func (d *SimpleDependencyInfo) GetHasNoCompileAnnotation() bool { return d.hasNoCompileAnnotation }

// Java: DependencyInfo.getRequiredSymbols()
func (d *SimpleDependencyInfo) GetRequiredSymbols() []string {
	return DependencyInfo_GetRequiredSymbols(d)
}

// Java: DependencyInfo.isEs6Module()
func (d *SimpleDependencyInfo) IsEs6Module() bool { return d.moduleKind == simpleDepModuleES6 }

// Java: DependencyInfo.isGoogModule()
func (d *SimpleDependencyInfo) IsGoogModule() bool { return d.moduleKind == simpleDepModuleGoog }

// Java: AutoValue_SimpleDependencyInfo.toString()
func (d *SimpleDependencyInfo) String() string {
	return "SimpleDependencyInfo{" +
		"name=" + d.name + ", " +
		"pathRelativeToClosureBase=" + d.pathRelativeToClosureBase + ", " +
		"provides=" + jutil.CollectionToString(d.provides) + ", " +
		"requires=" + jutil.CollectionToString(d.requires) + ", " +
		"typeRequires=" + jutil.CollectionToString(d.typeRequires) + ", " +
		"loadFlags=" + d.loadFlags.String() + ", " +
		"hasExternsAnnotation=" + strconv.FormatBool(d.hasExternsAnnotation) + ", " +
		"hasNoCompileAnnotation=" + strconv.FormatBool(d.hasNoCompileAnnotation) +
		"}"
}

// Java: AutoValue_SimpleDependencyInfo.equals(Object)
func (d *SimpleDependencyInfo) Equals(o any) bool {
	that, ok := o.(*SimpleDependencyInfo)
	if !ok || that == nil {
		return false
	}
	if that == d {
		return true
	}
	return d.name == that.name &&
		d.pathRelativeToClosureBase == that.pathRelativeToClosureBase &&
		slices.Equal(d.provides, that.provides) &&
		slices.EqualFunc(d.requires, that.requires, func(a, b *DependencyInfo_Require) bool {
			return a == b || (a != nil && a.Equals(b))
		}) &&
		slices.Equal(d.typeRequires, that.typeRequires) &&
		d.loadFlags.Equals(that.loadFlags) &&
		d.hasExternsAnnotation == that.hasExternsAnnotation &&
		d.hasNoCompileAnnotation == that.hasNoCompileAnnotation
}

func simpleDependencyInfo_stringSliceHash(ss []string) int32 {
	h := int32(1)
	for _, s := range ss {
		h = 31*h + jutil.StringHashCode(s)
	}
	return h
}

func simpleDependencyInfo_requireSliceHash(rs []*DependencyInfo_Require) int32 {
	h := int32(1)
	for _, r := range rs {
		var eh int32
		if r != nil {
			eh = r.HashCode()
		}
		h = 31*h + eh
	}
	return h
}

// Java: AutoValue_SimpleDependencyInfo.hashCode()
func (d *SimpleDependencyInfo) HashCode() int32 {
	h := int32(1)
	h *= 1000003
	h ^= jutil.StringHashCode(d.name)
	h *= 1000003
	h ^= jutil.StringHashCode(d.pathRelativeToClosureBase)
	h *= 1000003
	h ^= simpleDependencyInfo_stringSliceHash(d.provides)
	h *= 1000003
	h ^= simpleDependencyInfo_requireSliceHash(d.requires)
	h *= 1000003
	h ^= simpleDependencyInfo_stringSliceHash(d.typeRequires)
	h *= 1000003
	h ^= d.loadFlags.HashCode()
	h *= 1000003
	h ^= simpleDependencyInfo_boolHash(d.hasExternsAnnotation)
	h *= 1000003
	h ^= simpleDependencyInfo_boolHash(d.hasNoCompileAnnotation)
	return h
}

func simpleDependencyInfo_boolHash(b bool) int32 {
	if b {
		return 1231
	}
	return 1237
}

// SimpleDependencyInfo_Builder constructs instances of SimpleDependencyInfo.
type SimpleDependencyInfo_Builder struct {
	info SimpleDependencyInfo
}

// SimpleDependencyInfo_Builder_ is SimpleDependencyInfo.builder(String, String).
func SimpleDependencyInfo_Builder_(srcPathRelativeToClosure string, pathOfDefiningFile string) *SimpleDependencyInfo_Builder {
	return &SimpleDependencyInfo_Builder{
		info: SimpleDependencyInfo{
			name:                      pathOfDefiningFile,
			pathRelativeToClosureBase: srcPathRelativeToClosure,
			provides:                  simpleDependencyInfo_emptyStrings,
			requires:                  simpleDependencyInfo_emptyRequires,
			typeRequires:              simpleDependencyInfo_emptyStrings,
			loadFlags:                 EmptyLoadFlags,
		},
	}
}

func simpleDependencyInfo_copyStrings(s []string) []string {
	if len(s) == 0 {
		return simpleDependencyInfo_emptyStrings
	}
	return slices.Clone(s)
}

func simpleDependencyInfo_copyRequires(requires []*DependencyInfo_Require) []*DependencyInfo_Require {
	if len(requires) == 0 {
		return simpleDependencyInfo_emptyRequires
	}
	for _, r := range requires {
		jutil.CheckNotNull(r)
	}
	return slices.Clone(requires)
}

// Java: SimpleDependencyInfo.Builder.setProvides(Collection<String>)
func (b *SimpleDependencyInfo_Builder) SetProvides(provides []string) *SimpleDependencyInfo_Builder {
	b.info.provides = simpleDependencyInfo_copyStrings(provides)
	return b
}

// Java: SimpleDependencyInfo.Builder.setProvides(String...)
func (b *SimpleDependencyInfo_Builder) SetProvidesArray(provides ...string) *SimpleDependencyInfo_Builder {
	b.info.provides = simpleDependencyInfo_copyStrings(provides)
	return b
}

// Java: SimpleDependencyInfo.Builder.setRequires(Collection<Require>)
func (b *SimpleDependencyInfo_Builder) SetRequires(requires []*DependencyInfo_Require) *SimpleDependencyInfo_Builder {
	b.info.requires = simpleDependencyInfo_copyRequires(requires)
	return b
}

// Java: SimpleDependencyInfo.Builder.setRequires(Require...)
func (b *SimpleDependencyInfo_Builder) SetRequiresArray(requires ...*DependencyInfo_Require) *SimpleDependencyInfo_Builder {
	b.info.requires = simpleDependencyInfo_copyRequires(requires)
	return b
}

// Java: SimpleDependencyInfo.Builder.setTypeRequires(Collection<String>)
func (b *SimpleDependencyInfo_Builder) SetTypeRequires(typeRequires []string) *SimpleDependencyInfo_Builder {
	b.info.typeRequires = simpleDependencyInfo_copyStrings(typeRequires)
	return b
}

// Java: SimpleDependencyInfo.Builder.setTypeRequires(String...)
func (b *SimpleDependencyInfo_Builder) SetTypeRequiresArray(typeRequires ...string) *SimpleDependencyInfo_Builder {
	b.info.typeRequires = simpleDependencyInfo_copyStrings(typeRequires)
	return b
}

// Java: SimpleDependencyInfo.Builder.setLoadFlags(Map<String, String>)
func (b *SimpleDependencyInfo_Builder) SetLoadFlags(loadFlags jutil.MapSource[string, string]) *SimpleDependencyInfo_Builder {
	if jutil.IsNil(loadFlags) {
		panic(jutil.NewNullPointerException())
	}
	if lf, ok := loadFlags.(LoadFlags); ok {
		return b.SetLoadFlagsValue(lf)
	}
	var buf [4]LoadFlag
	entries := buf[:0]
	for k, v := range loadFlags.All() {
		putLoadFlag(&entries, k, v)
	}
	b.info.loadFlags, b.info.moduleKind = normalizeLoadFlagsEntries(entries, true)
	return b
}

// SetLoadFlagsValue sets LoadFlags directly without boxing into MapSource.
func (b *SimpleDependencyInfo_Builder) SetLoadFlagsValue(loadFlags LoadFlags) *SimpleDependencyInfo_Builder {
	b.info.loadFlags, b.info.moduleKind = normalizeLoadFlagsEntries(loadFlags.entries, false)
	return b
}

// Java: SimpleDependencyInfo.Builder.setHasExternsAnnotation(boolean)
func (b *SimpleDependencyInfo_Builder) SetHasExternsAnnotation(hasExternsAnnotation bool) *SimpleDependencyInfo_Builder {
	b.info.hasExternsAnnotation = hasExternsAnnotation
	return b
}

// Java: SimpleDependencyInfo.Builder.setHasNoCompileAnnotation(boolean)
func (b *SimpleDependencyInfo_Builder) SetHasNoCompileAnnotation(hasNoCompileAnnotation bool) *SimpleDependencyInfo_Builder {
	b.info.hasNoCompileAnnotation = hasNoCompileAnnotation
	return b
}

// Java: SimpleDependencyInfo.Builder.setGoogModule(boolean)
func (b *SimpleDependencyInfo_Builder) SetGoogModule(isModule bool) *SimpleDependencyInfo_Builder {
	if isModule {
		b.info.loadFlags = GoogModuleLoadFlags
		b.info.moduleKind = simpleDepModuleGoog
	} else {
		b.info.loadFlags = EmptyLoadFlags
		b.info.moduleKind = simpleDepModuleNone
	}
	return b
}

// Java: SimpleDependencyInfo.Builder.build()
func (b *SimpleDependencyInfo_Builder) Build() *SimpleDependencyInfo {
	out := b.info
	return &out
}

// Java: SimpleDependencyInfo.EMPTY
var SimpleDependencyInfo_EMPTY = SimpleDependencyInfo_Builder_("", "").Build()
