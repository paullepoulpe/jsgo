// Port of com/google/javascript/jscomp/parsing/Config.java @ closure-compiler v20260415
// (with its generated subclass AutoValue_Config.java folded in)

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

// Package parsing is the port of com.google.javascript.jscomp.parsing: the
// IRFactory that turns parse trees into rhino Nodes, the JSDoc parser, and
// their configuration.
package parsing

import (
	"iter"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
)

// Config is the configuration for the AST factory. Should be shared across
// AST creation for all files of a compilation process.
//
// Go-only: language-adapter — Config is an @AutoValue abstract class; its
// generated implementation AutoValue_Config (fields, accessors, toString,
// equals, hashCode, toBuilder and the Builder subclass) is folded into this
// struct and Config_Builder.
type Config struct {
	languageMode          Config_LanguageMode
	strictMode            Config_StrictMode
	jsDocParsingMode      Config_JsDocParsing
	runMode               Config_RunMode
	annotations           *jutil.LinkedHashMap[string, Annotation]
	suppressionNames      *jutil.LinkedHashSet[string]
	closurePrimitiveNames *jutil.LinkedHashSet[string]
	parseInlineSourceMaps bool
}

// Config_StrictMode is the level of language strictness required for the
// input source code.
type Config_StrictMode int32

const (
	Config_StrictMode_STRICT Config_StrictMode = iota
	Config_StrictMode_SLOPPY
)

var config_StrictMode_names = [...]string{"STRICT", "SLOPPY"}

// Config_StrictMode_values is StrictMode.values().
//
// Java: Config.StrictMode.values()
func Config_StrictMode_values() []Config_StrictMode {
	return []Config_StrictMode{Config_StrictMode_STRICT, Config_StrictMode_SLOPPY}
}

// Name is Enum.name().
//
// Java: Config.StrictMode.name()
func (m Config_StrictMode) Name() string { return config_StrictMode_names[m] }

// String is Enum.toString().
//
// Java: Config.StrictMode.toString()
func (m Config_StrictMode) String() string { return m.Name() }

// Java: Config.StrictMode.isStrict()
func (m Config_StrictMode) IsStrict() bool {
	return m == Config_StrictMode_STRICT
}

// Config_LanguageMode is the JavaScript mode.
type Config_LanguageMode int32

// Note that minimumRequiredFor() relies on these being defined in order from
// fewest features to most features, and _STRICT versions should be supplied
// after unspecified strictness.
const (
	Config_LanguageMode_ECMASCRIPT3 Config_LanguageMode = iota
	Config_LanguageMode_ECMASCRIPT5
	Config_LanguageMode_ECMASCRIPT_2015
	Config_LanguageMode_ECMASCRIPT_2016
	Config_LanguageMode_ECMASCRIPT_2017
	Config_LanguageMode_ECMASCRIPT_2018
	Config_LanguageMode_ECMASCRIPT_2019
	Config_LanguageMode_ECMASCRIPT_2020
	Config_LanguageMode_ECMASCRIPT_2021
	Config_LanguageMode_ECMASCRIPT_2022
	// NOTE: When adding a new language level here, also update
	// latestEcmaScript() below to return it.
	Config_LanguageMode_ES_NEXT
	Config_LanguageMode_UNSTABLE
	Config_LanguageMode_UNSUPPORTED
)

var config_LanguageMode_names = [...]string{
	"ECMASCRIPT3",
	"ECMASCRIPT5",
	"ECMASCRIPT_2015",
	"ECMASCRIPT_2016",
	"ECMASCRIPT_2017",
	"ECMASCRIPT_2018",
	"ECMASCRIPT_2019",
	"ECMASCRIPT_2020",
	"ECMASCRIPT_2021",
	"ECMASCRIPT_2022",
	"ES_NEXT",
	"UNSTABLE",
	"UNSUPPORTED",
}

// config_LanguageMode_featureSet is the enum field LanguageMode.featureSet,
// indexed by ordinal.
var config_LanguageMode_featureSet = [...]*parser.FeatureSet{
	parser.FeatureSet_ES3,
	parser.FeatureSet_ES5,
	parser.FeatureSet_ES2015_MODULES,
	parser.FeatureSet_ES2016_MODULES,
	parser.FeatureSet_ES2017_MODULES,
	parser.FeatureSet_ES2018_MODULES,
	parser.FeatureSet_ES2019_MODULES,
	parser.FeatureSet_ES2020_MODULES,
	parser.FeatureSet_ES2021_MODULES,
	parser.FeatureSet_ES2022_MODULES,
	parser.FeatureSet_ES_NEXT,
	parser.FeatureSet_ES_UNSTABLE,
	parser.FeatureSet_ES_UNSUPPORTED,
}

// Config_LanguageMode_values is LanguageMode.values().
//
// Java: Config.LanguageMode.values()
func Config_LanguageMode_values() []Config_LanguageMode {
	values := make([]Config_LanguageMode, len(config_LanguageMode_names))
	for i := range values {
		values[i] = Config_LanguageMode(i)
	}
	return values
}

// Name is Enum.name().
//
// Java: Config.LanguageMode.name()
func (m Config_LanguageMode) Name() string { return config_LanguageMode_names[m] }

// String is Enum.toString().
//
// Java: Config.LanguageMode.toString()
func (m Config_LanguageMode) String() string { return m.Name() }

// FeatureSet is the public final field LanguageMode.featureSet.
//
// Java: Config.LanguageMode.featureSet
func (m Config_LanguageMode) FeatureSet() *parser.FeatureSet {
	return config_LanguageMode_featureSet[m]
}

// Config_LanguageMode_MinimumRequiredFor returns the lowest LanguageMode that
// supports the specified feature.
//
// Java: Config.LanguageMode.minimumRequiredFor(FeatureSet.Feature)
func Config_LanguageMode_MinimumRequiredFor(feature parser.FeatureSet_Feature) Config_LanguageMode {
	// relies on the LanguageMode enums being in the right order
	for _, mode := range Config_LanguageMode_values() {
		if mode.FeatureSet().Has(feature) {
			return mode
		}
	}
	panic(jutil.NewIllegalStateExceptionMsg("No input language mode supports feature: " + feature.String()))
}

// Config_LanguageMode_MinimumRequiredForSet returns the lowest LanguageMode
// that supports the specified feature set.
//
// Java: Config.LanguageMode.minimumRequiredForSet(FeatureSet)
func Config_LanguageMode_MinimumRequiredForSet(featureSet *parser.FeatureSet) Config_LanguageMode {
	for _, mode := range Config_LanguageMode_values() {
		if mode.FeatureSet().Contains(featureSet) {
			return mode
		}
	}
	panic(jutil.NewIllegalStateExceptionMsg("No input language mode supports feature set: " + jutil.ToJavaString(featureSet)))
}

// Java: Config.LanguageMode.latestEcmaScript()
func Config_LanguageMode_LatestEcmaScript() Config_LanguageMode {
	return Config_LanguageMode_ECMASCRIPT_2021
}

// Config_JsDocParsing says whether to parse the descriptions of JsDoc
// comments.
type Config_JsDocParsing int32

const (
	Config_JsDocParsing_TYPES_ONLY Config_JsDocParsing = iota
	Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE
	Config_JsDocParsing_INCLUDE_DESCRIPTIONS_WITH_WHITESPACE
	Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	Config_JsDocParsing_LICENSE_COMMENTS_ONLY
)

var config_JsDocParsing_names = [...]string{
	"TYPES_ONLY",
	"INCLUDE_DESCRIPTIONS_NO_WHITESPACE",
	"INCLUDE_DESCRIPTIONS_WITH_WHITESPACE",
	"INCLUDE_ALL_COMMENTS",
	"LICENSE_COMMENTS_ONLY",
}

// Config_JsDocParsing_values is JsDocParsing.values().
//
// Java: Config.JsDocParsing.values()
func Config_JsDocParsing_values() []Config_JsDocParsing {
	values := make([]Config_JsDocParsing, len(config_JsDocParsing_names))
	for i := range values {
		values[i] = Config_JsDocParsing(i)
	}
	return values
}

// Name is Enum.name().
//
// Java: Config.JsDocParsing.name()
func (p Config_JsDocParsing) Name() string { return config_JsDocParsing_names[p] }

// String is Enum.toString().
//
// Java: Config.JsDocParsing.toString()
func (p Config_JsDocParsing) String() string { return p.Name() }

// Java: Config.JsDocParsing.shouldParseDescriptions()
func (p Config_JsDocParsing) ShouldParseDescriptions() bool {
	return p != Config_JsDocParsing_TYPES_ONLY
}

// Java: Config.JsDocParsing.shouldPreserveWhitespace()
func (p Config_JsDocParsing) ShouldPreserveWhitespace() bool {
	return p == Config_JsDocParsing_INCLUDE_DESCRIPTIONS_WITH_WHITESPACE ||
		p == Config_JsDocParsing_INCLUDE_ALL_COMMENTS ||
		p == Config_JsDocParsing_LICENSE_COMMENTS_ONLY
}

// Config_RunMode says whether to keep going after encountering a parse
// error.
type Config_RunMode int32

const (
	Config_RunMode_STOP_AFTER_ERROR Config_RunMode = iota
	Config_RunMode_KEEP_GOING
)

var config_RunMode_names = [...]string{"STOP_AFTER_ERROR", "KEEP_GOING"}

// Config_RunMode_values is RunMode.values().
//
// Java: Config.RunMode.values()
func Config_RunMode_values() []Config_RunMode {
	return []Config_RunMode{Config_RunMode_STOP_AFTER_ERROR, Config_RunMode_KEEP_GOING}
}

// Name is Enum.name().
//
// Java: Config.RunMode.name()
func (m Config_RunMode) Name() string { return config_RunMode_names[m] }

// String is Enum.toString().
//
// Java: Config.RunMode.toString()
func (m Config_RunMode) String() string { return m.Name() }

// NewConfig is the private AutoValue_Config constructor.
//
// Java: AutoValue_Config(Config.LanguageMode, Config.StrictMode, Config.JsDocParsing,
// Config.RunMode, ImmutableMap<String, Annotation>, ImmutableSet<String>, ImmutableSet<String>,
// boolean)
func NewConfig(
	languageMode Config_LanguageMode,
	strictMode Config_StrictMode,
	jsDocParsingMode Config_JsDocParsing,
	runMode Config_RunMode,
	annotations *jutil.LinkedHashMap[string, Annotation],
	suppressionNames *jutil.LinkedHashSet[string],
	closurePrimitiveNames *jutil.LinkedHashSet[string],
	parseInlineSourceMaps bool,
) *Config {
	c := &Config{}
	c.languageMode = languageMode
	c.strictMode = strictMode
	c.jsDocParsingMode = jsDocParsingMode
	c.runMode = runMode
	c.annotations = annotations
	c.suppressionNames = suppressionNames
	c.closurePrimitiveNames = closurePrimitiveNames
	c.parseInlineSourceMaps = parseInlineSourceMaps
	return c
}

// LanguageMode is the language level to accept.
//
// Java: Config.languageMode()
func (c *Config) LanguageMode() Config_LanguageMode {
	return c.languageMode
}

// StrictMode says whether to assume input is strict mode compliant.
//
// Java: Config.strictMode()
func (c *Config) StrictMode() Config_StrictMode {
	return c.strictMode
}

// JsDocParsingMode says how to parse the descriptions of JsDoc comments.
//
// Java: Config.jsDocParsingMode()
func (c *Config) JsDocParsingMode() Config_JsDocParsing {
	return c.jsDocParsingMode
}

// RunMode says whether to keep going after encountering a parse error.
//
// Java: Config.runMode()
func (c *Config) RunMode() Config_RunMode {
	return c.runMode
}

// Annotations are the recognized JSDoc annotations, mapped from their name
// to their internal representation.
//
// Java: Config.annotations()
func (c *Config) Annotations() *jutil.LinkedHashMap[string, Annotation] {
	return c.annotations
}

// SuppressionNames is the set of recognized names in a `@suppress` tag.
//
// Java: Config.suppressionNames()
func (c *Config) SuppressionNames() *jutil.LinkedHashSet[string] {
	return c.suppressionNames
}

// ClosurePrimitiveNames is the set of recognized names in a
// `@closurePrimitive` tag.
//
// Java: Config.closurePrimitiveNames()
func (c *Config) ClosurePrimitiveNames() *jutil.LinkedHashSet[string] {
	return c.closurePrimitiveNames
}

// ParseInlineSourceMaps says whether to parse inline source maps
// (//# sourceMappingURL=data:...).
//
// Java: Config.parseInlineSourceMaps()
func (c *Config) ParseInlineSourceMaps() bool {
	return c.parseInlineSourceMaps
}

// AnnotationNames is annotations().keySet().
//
// Java: Config.annotationNames()
func (c *Config) AnnotationNames() *jutil.LinkedHashSet[string] {
	keys := jutil.NewLinkedHashSet[string]()
	for k := range c.Annotations().Keys() {
		keys.Add(k)
	}
	return keys
}

// Java: AutoValue_Config.toString()
func (c *Config) String() string {
	return "Config{" +
		"languageMode=" + c.languageMode.String() + ", " +
		"strictMode=" + c.strictMode.String() + ", " +
		"jsDocParsingMode=" + c.jsDocParsingMode.String() + ", " +
		"runMode=" + c.runMode.String() + ", " +
		"annotations=" + c.annotations.String() + ", " +
		"suppressionNames=" + c.suppressionNames.String() + ", " +
		"closurePrimitiveNames=" + c.closurePrimitiveNames.String() + ", " +
		"parseInlineSourceMaps=" + jutil.ToJavaString(c.parseInlineSourceMaps) +
		"}"
}

// Java: AutoValue_Config.equals(Object)
func (c *Config) Equals(o any) bool {
	if that, ok := o.(*Config); ok && that == c {
		return true
	}
	if that, ok := o.(*Config); ok && that != nil {
		return c.languageMode == that.LanguageMode() &&
			c.strictMode == that.StrictMode() &&
			c.jsDocParsingMode == that.JsDocParsingMode() &&
			c.runMode == that.RunMode() &&
			c.annotations.Equals(that.Annotations()) &&
			c.suppressionNames.Equals(that.SuppressionNames()) &&
			c.closurePrimitiveNames.Equals(that.ClosurePrimitiveNames()) &&
			c.parseInlineSourceMaps == that.ParseInlineSourceMaps()
	}
	return false
}

// HashCode mirrors the generated hashCode. Java hashes the enum fields by
// identity (Enum.hashCode), which differs from run to run; here they hash by
// value (jutil.IdentityHashCode of a value type).
//
// Java: AutoValue_Config.hashCode()
func (c *Config) HashCode() int32 {
	var h int32 = 1
	h *= 1000003
	h ^= jutil.IdentityHashCode(c.languageMode)
	h *= 1000003
	h ^= jutil.IdentityHashCode(c.strictMode)
	h *= 1000003
	h ^= jutil.IdentityHashCode(c.jsDocParsingMode)
	h *= 1000003
	h ^= jutil.IdentityHashCode(c.runMode)
	h *= 1000003
	h ^= c.annotations.HashCode()
	h *= 1000003
	h ^= c.suppressionNames.HashCode()
	h *= 1000003
	h ^= c.closurePrimitiveNames.HashCode()
	h *= 1000003
	if c.parseInlineSourceMaps {
		h ^= 1231
	} else {
		h ^= 1237
	}
	return h
}

// Java: Config.toBuilder() / AutoValue_Config.toBuilder()
func (c *Config) ToBuilder() *Config_Builder {
	return NewConfig_BuilderFrom(c)
}

// Config_Builder_ returns a builder with every property set to its default.
//
// Go-only: language-adapter — named with a trailing underscore because the
// nested type Config_Builder has the same Go name.
//
// Java: Config.builder()
func Config_Builder_() *Config_Builder {
	return NewConfig_Builder().
		SetLanguageMode(Config_LanguageMode_UNSUPPORTED).
		SetStrictMode(Config_StrictMode_STRICT).
		SetJsDocParsingMode(Config_JsDocParsing_TYPES_ONLY).
		SetRunMode(Config_RunMode_STOP_AFTER_ERROR).
		SetExtraAnnotationNames(jutil.NewLinkedHashSet[string]().All()).
		SetSuppressionNames(jutil.NewLinkedHashSet[string]().All()).
		SetClosurePrimitiveNames(jutil.NewLinkedHashSet[string]().All()).
		SetParseInlineSourceMaps(false)
}

// Config_Builder is the builder for a Config.
//
// Go-only: language-adapter — Config.Builder is an abstract @AutoValue.Builder;
// its generated implementation AutoValue_Config.Builder is folded in. Java
// tracks unset object properties as null; the enum-typed ones are Go value
// types, so a bit in `set` records whether each was set (bit 0 is the
// generated set$0 bit for parseInlineSourceMaps).
type Config_Builder struct {
	languageMode          Config_LanguageMode
	strictMode            Config_StrictMode
	jsDocParsingMode      Config_JsDocParsing
	runMode               Config_RunMode
	annotations           *jutil.LinkedHashMap[string, Annotation]
	suppressionNames      *jutil.LinkedHashSet[string]
	closurePrimitiveNames *jutil.LinkedHashSet[string]
	parseInlineSourceMaps bool
	set                   uint8
}

// Go-only: language-adapter — the bits of Config_Builder.set.
const (
	config_Builder_setParseInlineSourceMaps uint8 = 1 << iota
	config_Builder_setLanguageMode
	config_Builder_setStrictMode
	config_Builder_setJsDocParsingMode
	config_Builder_setRunMode
)

// Java: AutoValue_Config.Builder()
func NewConfig_Builder() *Config_Builder {
	return &Config_Builder{}
}

// Java: AutoValue_Config.Builder(Config)
func NewConfig_BuilderFrom(source *Config) *Config_Builder {
	b := &Config_Builder{}
	b.languageMode = source.LanguageMode()
	b.strictMode = source.StrictMode()
	b.jsDocParsingMode = source.JsDocParsingMode()
	b.runMode = source.RunMode()
	b.annotations = source.Annotations()
	b.suppressionNames = source.SuppressionNames()
	b.closurePrimitiveNames = source.ClosurePrimitiveNames()
	b.parseInlineSourceMaps = source.ParseInlineSourceMaps()
	b.set = config_Builder_setParseInlineSourceMaps | config_Builder_setLanguageMode |
		config_Builder_setStrictMode | config_Builder_setJsDocParsingMode | config_Builder_setRunMode
	return b
}

// Java: Config.Builder.setLanguageMode(LanguageMode)
func (b *Config_Builder) SetLanguageMode(languageMode Config_LanguageMode) *Config_Builder {
	b.languageMode = languageMode
	b.set |= config_Builder_setLanguageMode
	return b
}

// Java: Config.Builder.setStrictMode(StrictMode)
func (b *Config_Builder) SetStrictMode(strictMode Config_StrictMode) *Config_Builder {
	b.strictMode = strictMode
	b.set |= config_Builder_setStrictMode
	return b
}

// Java: Config.Builder.setJsDocParsingMode(JsDocParsing)
func (b *Config_Builder) SetJsDocParsingMode(jsDocParsingMode Config_JsDocParsing) *Config_Builder {
	b.jsDocParsingMode = jsDocParsingMode
	b.set |= config_Builder_setJsDocParsingMode
	return b
}

// Java: Config.Builder.setRunMode(RunMode)
func (b *Config_Builder) SetRunMode(runMode Config_RunMode) *Config_Builder {
	b.runMode = runMode
	b.set |= config_Builder_setRunMode
	return b
}

// Java: Config.Builder.setParseInlineSourceMaps(boolean)
func (b *Config_Builder) SetParseInlineSourceMaps(parseInlineSourceMaps bool) *Config_Builder {
	b.parseInlineSourceMaps = parseInlineSourceMaps
	b.set |= config_Builder_setParseInlineSourceMaps
	return b
}

// SetSuppressionNames stores ImmutableSet.copyOf(names).
//
// Java: Config.Builder.setSuppressionNames(Iterable<String>)
func (b *Config_Builder) SetSuppressionNames(names iter.Seq[string]) *Config_Builder {
	b.suppressionNames = config_immutableSetCopyOf(names)
	return b
}

// SetClosurePrimitiveNames stores ImmutableSet.copyOf(names).
//
// Java: Config.Builder.setClosurePrimitiveNames(Iterable<String>)
func (b *Config_Builder) SetClosurePrimitiveNames(names iter.Seq[string]) *Config_Builder {
	b.closurePrimitiveNames = config_immutableSetCopyOf(names)
	return b
}

// Java: Config.Builder.setExtraAnnotationNames(Iterable<String>)
func (b *Config_Builder) SetExtraAnnotationNames(names iter.Seq[string]) *Config_Builder {
	return b.SetAnnotations(Config_BuildAnnotations(names))
}

// Java: Config.Builder.build() / AutoValue_Config.Builder.build()
func (b *Config_Builder) Build() *Config {
	const allSet = config_Builder_setParseInlineSourceMaps | config_Builder_setLanguageMode |
		config_Builder_setStrictMode | config_Builder_setJsDocParsingMode | config_Builder_setRunMode
	if b.set != allSet ||
		b.annotations == nil ||
		b.suppressionNames == nil ||
		b.closurePrimitiveNames == nil {
		var missing strings.Builder
		if b.set&config_Builder_setLanguageMode == 0 {
			missing.WriteString(" languageMode")
		}
		if b.set&config_Builder_setStrictMode == 0 {
			missing.WriteString(" strictMode")
		}
		if b.set&config_Builder_setJsDocParsingMode == 0 {
			missing.WriteString(" jsDocParsingMode")
		}
		if b.set&config_Builder_setRunMode == 0 {
			missing.WriteString(" runMode")
		}
		if b.annotations == nil {
			missing.WriteString(" annotations")
		}
		if b.suppressionNames == nil {
			missing.WriteString(" suppressionNames")
		}
		if b.closurePrimitiveNames == nil {
			missing.WriteString(" closurePrimitiveNames")
		}
		if b.set&config_Builder_setParseInlineSourceMaps == 0 {
			missing.WriteString(" parseInlineSourceMaps")
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing.String()))
	}
	return NewConfig(
		b.languageMode,
		b.strictMode,
		b.jsDocParsingMode,
		b.runMode,
		b.annotations,
		b.suppressionNames,
		b.closurePrimitiveNames,
		b.parseInlineSourceMaps)
}

// SetAnnotations is intended to be used internally only (but isn't private
// due to AutoValue).
//
// Java: Config.Builder.setAnnotations(ImmutableMap<String, Annotation>)
func (b *Config_Builder) SetAnnotations(names *jutil.LinkedHashMap[string, Annotation]) *Config_Builder {
	if names == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null annotations"))
	}
	b.annotations = names
	return b
}

// Config_BuildAnnotations creates the annotation names from the
// user-specified annotation allowlist.
//
// Java: Config.buildAnnotations(Iterable<String>)
func Config_BuildAnnotations(allowlist iter.Seq[string]) *jutil.LinkedHashMap[string, Annotation] {
	annotationsBuilder := jutil.NewImmutableMapBuilder[string, Annotation]()
	for k, v := range Annotation_recognizedAnnotations.All() {
		annotationsBuilder.Put(k, v)
	}
	for unrecognizedAnnotation := range allowlist {
		if unrecognizedAnnotation != "" &&
			!Annotation_recognizedAnnotations.ContainsKey(unrecognizedAnnotation) {
			annotationsBuilder.Put(unrecognizedAnnotation, Annotation_NOT_IMPLEMENTED)
		}
	}
	return annotationsBuilder.BuildOrThrow()
}

// config_immutableSetCopyOf is ImmutableSet.copyOf(Iterable<String>).
//
// Go-only: language-adapter — Guava ImmutableSet.copyOf.
func config_immutableSetCopyOf(names iter.Seq[string]) *jutil.LinkedHashSet[string] {
	if names == nil {
		panic(jutil.NewNullPointerException())
	}
	set := jutil.NewLinkedHashSet[string]()
	for name := range names {
		set.Add(name)
	}
	return set
}
