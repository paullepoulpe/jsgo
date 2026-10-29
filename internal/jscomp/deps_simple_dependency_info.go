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
	loadFlags                 *jutil.LinkedHashMap[string, string]
	moduleKind                uint8
	hasExternsAnnotation      bool
	hasNoCompileAnnotation    bool
}

var _ DependencyInfo = (*SimpleDependencyInfo)(nil)

var (
	simpleDependencyInfo_emptyStrings              = []string{}
	simpleDependencyInfo_emptyRequires             = []*DependencyInfo_Require{}
	simpleDependencyInfo_EMPTY_FLAGS               = jutil.NewLinkedHashMap[string, string]()
	simpleDependencyInfo_Builder_GOOG_MODULE_FLAGS = func() *jutil.LinkedHashMap[string, string] {
		m := jutil.NewLinkedHashMap[string, string]()
		m.Put("module", "goog")
		return m
	}()
	simpleDependencyInfo_ES6_FLAGS = func() *jutil.LinkedHashMap[string, string] {
		m := jutil.NewLinkedHashMap[string, string]()
		m.Put("module", "es6")
		return m
	}()
	simpleDependencyInfo_singletonModuleKind = map[*jutil.LinkedHashMap[string, string]]uint8{}
	simpleDependencyInfo_langModuleFlags     = func() map[[2]string]*jutil.LinkedHashMap[string, string] {
		langs := []string{"", "es5", "es6", "es7", "es8", "es9", "es_2020", "es_2021", "es_next", "ts_unsupported"}
		mods := []string{"", "goog", "es6"}
		out := make(map[[2]string]*jutil.LinkedHashMap[string, string], len(langs)*len(mods))
		for _, l := range langs {
			for _, mod := range mods {
				var m *jutil.LinkedHashMap[string, string]
				switch {
				case l == "" && mod == "":
					m = simpleDependencyInfo_EMPTY_FLAGS
				case l == "" && mod == "goog":
					m = simpleDependencyInfo_Builder_GOOG_MODULE_FLAGS
				case l == "" && mod == "es6":
					m = simpleDependencyInfo_ES6_FLAGS
				default:
					m = jutil.NewLinkedHashMap[string, string]()
					if l != "" {
						m.Put("lang", l)
					}
					if mod != "" {
						m.Put("module", mod)
					}
				}
				out[[2]string{l, mod}] = m
				var kind uint8
				switch mod {
				case "goog":
					kind = simpleDepModuleGoog
				case "es6":
					kind = simpleDepModuleES6
				}
				simpleDependencyInfo_singletonModuleKind[m] = kind
			}
		}
		return out
	}()
)

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

func simpleDependencyInfo_normalizeFlags(
	loadFlags *jutil.LinkedHashMap[string, string],
	cloneCustom bool,
) (*jutil.LinkedHashMap[string, string], uint8) {
	if loadFlags == nil || loadFlags.IsEmpty() {
		return simpleDependencyInfo_EMPTY_FLAGS, simpleDepModuleNone
	}
	if kind, ok := simpleDependencyInfo_singletonModuleKind[loadFlags]; ok {
		return loadFlags, kind
	}
	if loadFlags.Size() <= 2 {
		var lang, mod string
		validOrder := true
		idx := 0
		for k, v := range loadFlags.All() {
			switch k {
			case "lang":
				if idx != 0 || v == "" {
					validOrder = false
				}
				lang = v
			case "module":
				if v == "" {
					validOrder = false
				}
				mod = v
			default:
				validOrder = false
			}
			idx++
		}
		if validOrder {
			if singleton, ok := simpleDependencyInfo_langModuleFlags[[2]string{lang, mod}]; ok {
				return singleton, simpleDependencyInfo_moduleKindOf(mod)
			}
		}
	}
	mod, _ := loadFlags.Lookup("module")
	if cloneCustom {
		loadFlags = jutil.NewLinkedHashMapFrom(loadFlags)
	}
	return loadFlags, simpleDependencyInfo_moduleKindOf(mod)
}

func simpleDependencyInfo_flagsForLangModule(lang, mod string) *jutil.LinkedHashMap[string, string] {
	if m, ok := simpleDependencyInfo_langModuleFlags[[2]string{lang, mod}]; ok {
		return m
	}
	m := jutil.NewLinkedHashMap[string, string]()
	if lang != "" {
		m.Put("lang", lang)
	}
	if mod != "" {
		m.Put("module", mod)
	}
	return m
}

func newSimpleDependencyInfoDirect(
	srcPathRelativeToClosure string,
	pathOfDefiningFile string,
	provides []string,
	requires []*DependencyInfo_Require,
	typeRequires []string,
	loadFlags *jutil.LinkedHashMap[string, string],
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
	normFlags, modKind := simpleDependencyInfo_normalizeFlags(loadFlags, false)
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
func (d *SimpleDependencyInfo) GetLoadFlags() *jutil.LinkedHashMap[string, string] {
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
			loadFlags:                 simpleDependencyInfo_EMPTY_FLAGS,
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
func (b *SimpleDependencyInfo_Builder) SetLoadFlags(loadFlags *jutil.LinkedHashMap[string, string]) *SimpleDependencyInfo_Builder {
	jutil.CheckNotNull(loadFlags)
	b.info.loadFlags, b.info.moduleKind = simpleDependencyInfo_normalizeFlags(loadFlags, true)
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
		b.info.loadFlags = simpleDependencyInfo_Builder_GOOG_MODULE_FLAGS
		b.info.moduleKind = simpleDepModuleGoog
	} else {
		b.info.loadFlags = simpleDependencyInfo_EMPTY_FLAGS
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
