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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// SimpleDependencyInfo is a class to hold JS dependency information for a
// single .js file.
//
// Go-only: language-adapter — SimpleDependencyInfo is an @AutoValue: the
// abstract class and the generated AutoValue_SimpleDependencyInfo (fields,
// equals, hashCode, toString, Builder) are one struct here (PORTING.md §4).
type SimpleDependencyInfo struct {
	name                      string
	pathRelativeToClosureBase string
	provides                  []string
	requires                  []*DependencyInfo_Require
	typeRequires              []string
	loadFlags                 *jutil.LinkedHashMap[string, string]
	hasExternsAnnotation      bool
	hasNoCompileAnnotation    bool
}

var _ DependencyInfo = (*SimpleDependencyInfo)(nil)

var (
	simpleDependencyInfo_emptyStrings  = []string{}
	simpleDependencyInfo_emptyRequires = []*DependencyInfo_Require{}
	simpleDependencyInfo_EMPTY_FLAGS   = jutil.NewLinkedHashMap[string, string]()
	simpleDependencyInfo_ES6_FLAGS     = func() *jutil.LinkedHashMap[string, string] {
		m := jutil.NewLinkedHashMap[string, string]()
		m.Put("module", "es6")
		return m
	}()
	simpleDependencyInfo_langModuleFlags = func() map[[2]string]*jutil.LinkedHashMap[string, string] {
		langs := []string{"", "es5", "es6", "es7", "es8", "es9", "es_2020", "es_2021", "es_next", "ts_unsupported"}
		mods := []string{"", "goog", "es6"}
		out := make(map[[2]string]*jutil.LinkedHashMap[string, string], len(langs)*len(mods))
		for _, l := range langs {
			for _, mod := range mods {
				if l == "" && mod == "" {
					out[[2]string{l, mod}] = simpleDependencyInfo_EMPTY_FLAGS
					continue
				}
				if l == "" && mod == "goog" {
					out[[2]string{l, mod}] = simpleDependencyInfo_Builder_GOOG_MODULE_FLAGS
					continue
				}
				if l == "" && mod == "es6" {
					out[[2]string{l, mod}] = simpleDependencyInfo_ES6_FLAGS
					continue
				}
				m := jutil.NewLinkedHashMap[string, string]()
				if l != "" {
					m.Put("lang", l)
				}
				if mod != "" {
					m.Put("module", mod)
				}
				out[[2]string{l, mod}] = m
			}
		}
		return out
	}()
)

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
	var provCopy, typeReqCopy []string
	np, nt := len(provides), len(typeRequires)
	if np == 0 && nt == 0 {
		provCopy = simpleDependencyInfo_emptyStrings
		typeReqCopy = simpleDependencyInfo_emptyStrings
	} else {
		combined := make([]string, np+nt)
		copy(combined[:np], provides)
		copy(combined[np:], typeRequires)
		if np == 0 {
			provCopy = simpleDependencyInfo_emptyStrings
		} else {
			provCopy = combined[:np:np]
		}
		if nt == 0 {
			typeReqCopy = simpleDependencyInfo_emptyStrings
		} else {
			typeReqCopy = combined[np:]
		}
	}
	reqCopy := simpleDependencyInfo_copyRequires(requires)
	if loadFlags == nil || loadFlags.IsEmpty() {
		loadFlags = simpleDependencyInfo_EMPTY_FLAGS
	}
	return &SimpleDependencyInfo{
		name:                      pathOfDefiningFile,
		pathRelativeToClosureBase: srcPathRelativeToClosure,
		provides:                  provCopy,
		requires:                  reqCopy,
		typeRequires:              typeReqCopy,
		loadFlags:                 loadFlags,
		hasExternsAnnotation:      hasExternsAnnotation,
		hasNoCompileAnnotation:    hasNoCompileAnnotation,
	}
}

// SimpleDependencyInfo_Builder_ is SimpleDependencyInfo.builder(String, String).
//
// Go-only: language-adapter — the trailing underscore keeps the mechanical
// name SimpleDependencyInfo_Builder for the nested Builder type.
func SimpleDependencyInfo_Builder_(srcPathRelativeToClosure string, pathOfDefiningFile string) *SimpleDependencyInfo_Builder {
	return &SimpleDependencyInfo_Builder{
		name:                      pathOfDefiningFile,
		pathRelativeToClosureBase: srcPathRelativeToClosure,
		provides:                  simpleDependencyInfo_emptyStrings,
		requires:                  simpleDependencyInfo_emptyRequires,
		typeRequires:              simpleDependencyInfo_emptyStrings,
		loadFlags:                 simpleDependencyInfo_EMPTY_FLAGS,
		set0:                      1 | 2 | 4 | 8,
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
func (d *SimpleDependencyInfo) IsEs6Module() bool { return DependencyInfo_IsEs6Module(d) }

// Java: DependencyInfo.isGoogModule()
func (d *SimpleDependencyInfo) IsGoogModule() bool { return DependencyInfo_IsGoogModule(d) }

// Java: AutoValue_SimpleDependencyInfo.toString()
func (d *SimpleDependencyInfo) String() string {
	return "SimpleDependencyInfo{" +
		"name=" + d.name + ", " +
		"pathRelativeToClosureBase=" + d.pathRelativeToClosureBase + ", " +
		"provides=" + jutil.CollectionToString(d.provides) + ", " +
		"requires=" + jutil.CollectionToString(d.requires) + ", " +
		"typeRequires=" + jutil.CollectionToString(d.typeRequires) + ", " +
		"loadFlags=" + d.loadFlags.String() + ", " +
		"hasExternsAnnotation=" + jutil.ToJavaString(d.hasExternsAnnotation) + ", " +
		"hasNoCompileAnnotation=" + jutil.ToJavaString(d.hasNoCompileAnnotation) +
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
	return d.name == that.GetName() &&
		d.pathRelativeToClosureBase == that.GetPathRelativeToClosureBase() &&
		jutil.ObjectsEquals(jutil.NewArrayListFrom(d.provides), jutil.NewArrayListFrom(that.GetProvides())) &&
		jutil.ObjectsEquals(jutil.NewArrayListFrom(d.requires), jutil.NewArrayListFrom(that.GetRequires())) &&
		jutil.ObjectsEquals(jutil.NewArrayListFrom(d.typeRequires), jutil.NewArrayListFrom(that.GetTypeRequires())) &&
		d.loadFlags.Equals(that.GetLoadFlags()) &&
		d.hasExternsAnnotation == that.GetHasExternsAnnotation() &&
		d.hasNoCompileAnnotation == that.GetHasNoCompileAnnotation()
}

// Java: AutoValue_SimpleDependencyInfo.hashCode()
func (d *SimpleDependencyInfo) HashCode() int32 {
	h := int32(1)
	h *= 1000003
	h ^= jutil.StringHashCode(d.name)
	h *= 1000003
	h ^= jutil.StringHashCode(d.pathRelativeToClosureBase)
	h *= 1000003
	h ^= jutil.NewArrayListFrom(d.provides).HashCode()
	h *= 1000003
	h ^= jutil.NewArrayListFrom(d.requires).HashCode()
	h *= 1000003
	h ^= jutil.NewArrayListFrom(d.typeRequires).HashCode()
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

// SimpleDependencyInfo_Builder is the builder for constructing instances of
// SimpleDependencyInfo. Use SimpleDependencyInfo_Builder_ to create an
// instance.
type SimpleDependencyInfo_Builder struct {
	name                      string
	pathRelativeToClosureBase string
	provides                  []string
	requires                  []*DependencyInfo_Require
	typeRequires              []string
	loadFlags                 *jutil.LinkedHashMap[string, string]
	hasExternsAnnotation      bool
	hasNoCompileAnnotation    bool
	set0                      byte
}

// Java: AutoValue_SimpleDependencyInfo.Builder()
func NewSimpleDependencyInfo_Builder() *SimpleDependencyInfo_Builder {
	return &SimpleDependencyInfo_Builder{}
}

// Java: SimpleDependencyInfo.Builder.from(DependencyInfo)
func SimpleDependencyInfo_Builder_From(copy DependencyInfo) *SimpleDependencyInfo_Builder {
	return NewSimpleDependencyInfo_Builder().
		SetName(copy.GetName()).
		SetPathRelativeToClosureBase(copy.GetPathRelativeToClosureBase()).
		SetProvides(copy.GetProvides()).
		SetRequires(copy.GetRequires()).
		SetTypeRequires(copy.GetTypeRequires()).
		SetLoadFlags(copy.GetLoadFlags()).
		SetHasExternsAnnotation(copy.GetHasExternsAnnotation()).
		SetHasNoCompileAnnotation(copy.GetHasNoCompileAnnotation())
}

// Java: SimpleDependencyInfo.Builder.setName(String)
func (b *SimpleDependencyInfo_Builder) SetName(name string) *SimpleDependencyInfo_Builder {
	b.name = name
	b.set0 |= 4
	return b
}

// Java: SimpleDependencyInfo.Builder.setPathRelativeToClosureBase(String)
func (b *SimpleDependencyInfo_Builder) SetPathRelativeToClosureBase(srcPathRelativeToClosure string) *SimpleDependencyInfo_Builder {
	b.pathRelativeToClosureBase = srcPathRelativeToClosure
	b.set0 |= 8
	return b
}

func simpleDependencyInfo_copyStrings(s []string) []string {
	if len(s) == 0 {
		return simpleDependencyInfo_emptyStrings
	}
	out := make([]string, len(s))
	copy(out, s)
	return out
}

// Java: SimpleDependencyInfo.Builder.setProvides(Collection<String>)
func (b *SimpleDependencyInfo_Builder) SetProvides(provides []string) *SimpleDependencyInfo_Builder {
	b.provides = simpleDependencyInfo_copyStrings(provides)
	return b
}

// Java: SimpleDependencyInfo.Builder.setProvides(String...)
func (b *SimpleDependencyInfo_Builder) SetProvidesArray(provides ...string) *SimpleDependencyInfo_Builder {
	b.provides = simpleDependencyInfo_copyStrings(provides)
	return b
}

// Java: SimpleDependencyInfo.Builder.setRequires(Collection<Require>)
func (b *SimpleDependencyInfo_Builder) SetRequires(requires []*DependencyInfo_Require) *SimpleDependencyInfo_Builder {
	b.requires = simpleDependencyInfo_copyRequires(requires)
	return b
}

// Java: SimpleDependencyInfo.Builder.setRequires(Require...)
func (b *SimpleDependencyInfo_Builder) SetRequiresArray(requires ...*DependencyInfo_Require) *SimpleDependencyInfo_Builder {
	b.requires = simpleDependencyInfo_copyRequires(requires)
	return b
}

// simpleDependencyInfo_copyRequires is ImmutableList.copyOf, which rejects
// null elements.
func simpleDependencyInfo_copyRequires(requires []*DependencyInfo_Require) []*DependencyInfo_Require {
	if len(requires) == 0 {
		return simpleDependencyInfo_emptyRequires
	}
	for _, r := range requires {
		jutil.CheckNotNull(r)
	}
	out := make([]*DependencyInfo_Require, len(requires))
	copy(out, requires)
	return out
}

// Java: SimpleDependencyInfo.Builder.setTypeRequires(Collection<String>)
func (b *SimpleDependencyInfo_Builder) SetTypeRequires(typeRequires []string) *SimpleDependencyInfo_Builder {
	b.typeRequires = simpleDependencyInfo_copyStrings(typeRequires)
	return b
}

// Java: SimpleDependencyInfo.Builder.setTypeRequires(String...)
func (b *SimpleDependencyInfo_Builder) SetTypeRequiresArray(typeRequires ...string) *SimpleDependencyInfo_Builder {
	b.typeRequires = simpleDependencyInfo_copyStrings(typeRequires)
	return b
}

// Java: SimpleDependencyInfo.Builder.setLoadFlags(Map<String, String>)
func (b *SimpleDependencyInfo_Builder) SetLoadFlags(loadFlags *jutil.LinkedHashMap[string, string]) *SimpleDependencyInfo_Builder {
	jutil.CheckNotNull(loadFlags)
	if loadFlags.IsEmpty() {
		b.loadFlags = simpleDependencyInfo_EMPTY_FLAGS
	} else if loadFlags == simpleDependencyInfo_Builder_GOOG_MODULE_FLAGS || loadFlags == simpleDependencyInfo_ES6_FLAGS {
		b.loadFlags = loadFlags
	} else {
		b.loadFlags = jutil.NewLinkedHashMapFrom(loadFlags)
	}
	return b
}

// Java: SimpleDependencyInfo.Builder.setHasExternsAnnotation(boolean)
func (b *SimpleDependencyInfo_Builder) SetHasExternsAnnotation(hasExternsAnnotation bool) *SimpleDependencyInfo_Builder {
	b.hasExternsAnnotation = hasExternsAnnotation
	b.set0 |= 1
	return b
}

// Java: SimpleDependencyInfo.Builder.setHasNoCompileAnnotation(boolean)
func (b *SimpleDependencyInfo_Builder) SetHasNoCompileAnnotation(hasNoCompileAnnotation bool) *SimpleDependencyInfo_Builder {
	b.hasNoCompileAnnotation = hasNoCompileAnnotation
	b.set0 |= 2
	return b
}

var simpleDependencyInfo_Builder_GOOG_MODULE_FLAGS = func() *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	m.Put("module", "goog")
	return m
}()

// Java: SimpleDependencyInfo.Builder.setGoogModule(boolean)
func (b *SimpleDependencyInfo_Builder) SetGoogModule(isModule bool) *SimpleDependencyInfo_Builder {
	if isModule {
		b.loadFlags = simpleDependencyInfo_Builder_GOOG_MODULE_FLAGS
		return b
	}
	b.loadFlags = simpleDependencyInfo_EMPTY_FLAGS
	return b
}

// Java: SimpleDependencyInfo.Builder.build()
func (b *SimpleDependencyInfo_Builder) Build() *SimpleDependencyInfo {
	if b.set0 != (1|2|4|8) || b.provides == nil ||
		b.requires == nil || b.typeRequires == nil || b.loadFlags == nil {
		missing := jutil.NewStringBuilder()
		if b.set0&4 == 0 {
			missing.Append(" name")
		}
		if b.set0&8 == 0 {
			missing.Append(" pathRelativeToClosureBase")
		}
		if b.provides == nil {
			missing.Append(" provides")
		}
		if b.requires == nil {
			missing.Append(" requires")
		}
		if b.typeRequires == nil {
			missing.Append(" typeRequires")
		}
		if b.loadFlags == nil {
			missing.Append(" loadFlags")
		}
		if b.set0&1 == 0 {
			missing.Append(" hasExternsAnnotation")
		}
		if b.set0&2 == 0 {
			missing.Append(" hasNoCompileAnnotation")
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing.String()))
	}
	return &SimpleDependencyInfo{
		name:                      b.name,
		pathRelativeToClosureBase: b.pathRelativeToClosureBase,
		provides:                  b.provides,
		requires:                  b.requires,
		typeRequires:              b.typeRequires,
		loadFlags:                 b.loadFlags,
		hasExternsAnnotation:      b.hasExternsAnnotation,
		hasNoCompileAnnotation:    b.hasNoCompileAnnotation,
	}
}

// Java: SimpleDependencyInfo.EMPTY
var SimpleDependencyInfo_EMPTY = SimpleDependencyInfo_Builder_("", "").Build()
