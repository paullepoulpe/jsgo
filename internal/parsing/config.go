// Port of com/google/javascript/jscomp/parsing/Config.java @ closure-compiler v20261021

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
	"maps"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
)

var config_emptyStringSet = jutil.NewLinkedHashSet[string]()

// Config is the configuration for the AST factory. Should be shared across
// AST creation for all files of a compilation process.
type Config struct {
	languageMode          Config_LanguageMode
	strictMode            Config_StrictMode
	jsDocParsingMode      Config_JsDocParsing
	runMode               Config_RunMode
	annotations           map[string]Annotation
	extraAnnotations      *jutil.LinkedHashSet[string]
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
	for mode := Config_LanguageMode(0); int(mode) < len(config_LanguageMode_names); mode++ {
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
	for mode := Config_LanguageMode(0); int(mode) < len(config_LanguageMode_names); mode++ {
		if mode.FeatureSet().Contains(featureSet) {
			return mode
		}
	}
	fsStr := "null"
	if featureSet != nil {
		fsStr = featureSet.String()
	}
	panic(jutil.NewIllegalStateExceptionMsg("No input language mode supports feature set: " + fsStr))
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

// LookupAnnotation resolves a JSDoc @annotation name without materializing a
// merged map copy when extraAnnotations are configured.
func (c *Config) LookupAnnotation(name string) (Annotation, bool) {
	if c.annotations != nil {
		ann, ok := c.annotations[name]
		return ann, ok
	}
	if ann, ok := Annotation_recognizedAnnotations[name]; ok {
		return ann, true
	}
	if c.extraAnnotations.Contains(name) {
		return Annotation_NOT_IMPLEMENTED, true
	}
	return 0, false
}

// Annotations are the recognized JSDoc annotations, mapped from their name
// to their internal representation.
//
// Java: Config.annotations()
func (c *Config) Annotations() map[string]Annotation {
	if c.annotations != nil {
		return c.annotations
	}
	if c.extraAnnotations.IsEmpty() {
		return Annotation_recognizedAnnotations
	}
	m := make(map[string]Annotation, len(Annotation_recognizedAnnotations)+c.extraAnnotations.Size())
	for k, v := range Annotation_recognizedAnnotations { // order-independent: populating unordered lookup map
		m[k] = v
	}
	for k := range c.extraAnnotations.All() {
		m[k] = Annotation_NOT_IMPLEMENTED
	}
	return m
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

// AnnotationNames is annotations().keySet() in insertion order.
//
// Java: Config.annotationNames()
func (c *Config) AnnotationNames() []string {
	if c.extraAnnotations.IsEmpty() {
		return slices.Clone(annotation_recognizedOrder)
	}
	out := make([]string, 0, len(annotation_recognizedOrder)+c.extraAnnotations.Size())
	out = append(out, annotation_recognizedOrder...)
	for name := range c.extraAnnotations.All() {
		out = append(out, name)
	}
	return out
}

// Java: Config.toBuilder()
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
	return &Config_Builder{
		cfg: Config{
			languageMode:          Config_LanguageMode_UNSUPPORTED,
			strictMode:            Config_StrictMode_STRICT,
			jsDocParsingMode:      Config_JsDocParsing_TYPES_ONLY,
			runMode:               Config_RunMode_STOP_AFTER_ERROR,
			annotations:           Annotation_recognizedAnnotations,
			suppressionNames:      config_emptyStringSet,
			closurePrimitiveNames: config_emptyStringSet,
			parseInlineSourceMaps: false,
		},
	}
}

// Config_Builder is the builder for a Config.
type Config_Builder struct {
	cfg Config
}

// NewConfig_Builder returns a builder initialized to default Config values.
func NewConfig_Builder() *Config_Builder {
	return Config_Builder_()
}

// NewConfig_BuilderFrom returns a builder initialized from source.
func NewConfig_BuilderFrom(source *Config) *Config_Builder {
	return &Config_Builder{cfg: *source}
}

// Java: Config.Builder.setLanguageMode(LanguageMode)
func (b *Config_Builder) SetLanguageMode(languageMode Config_LanguageMode) *Config_Builder {
	b.cfg.languageMode = languageMode
	return b
}

// Java: Config.Builder.setStrictMode(StrictMode)
func (b *Config_Builder) SetStrictMode(strictMode Config_StrictMode) *Config_Builder {
	b.cfg.strictMode = strictMode
	return b
}

// Java: Config.Builder.setJsDocParsingMode(JsDocParsing)
func (b *Config_Builder) SetJsDocParsingMode(jsDocParsingMode Config_JsDocParsing) *Config_Builder {
	b.cfg.jsDocParsingMode = jsDocParsingMode
	return b
}

// Java: Config.Builder.setRunMode(RunMode)
func (b *Config_Builder) SetRunMode(runMode Config_RunMode) *Config_Builder {
	b.cfg.runMode = runMode
	return b
}

// Java: Config.Builder.setParseInlineSourceMaps(boolean)
func (b *Config_Builder) SetParseInlineSourceMaps(parseInlineSourceMaps bool) *Config_Builder {
	b.cfg.parseInlineSourceMaps = parseInlineSourceMaps
	return b
}

// SetSuppressionNames stores a copy of names.
//
// Java: Config.Builder.setSuppressionNames(Iterable<String>)
func (b *Config_Builder) SetSuppressionNames(names iter.Seq[string]) *Config_Builder {
	b.cfg.suppressionNames = config_stringSetCopyOf(names)
	return b
}

// SetClosurePrimitiveNames stores a copy of names.
//
// Java: Config.Builder.setClosurePrimitiveNames(Iterable<String>)
func (b *Config_Builder) SetClosurePrimitiveNames(names iter.Seq[string]) *Config_Builder {
	b.cfg.closurePrimitiveNames = config_stringSetCopyOf(names)
	return b
}

// Java: Config.Builder.setExtraAnnotationNames(Iterable<String>)
func (b *Config_Builder) SetExtraAnnotationNames(names iter.Seq[string]) *Config_Builder {
	if names == nil {
		panic(jutil.NewNullPointerException())
	}
	var extra *jutil.LinkedHashSet[string]
	for unrecognizedAnnotation := range names {
		if unrecognizedAnnotation == "" {
			continue
		}
		if _, recognized := Annotation_recognizedAnnotations[unrecognizedAnnotation]; recognized {
			continue
		}
		if extra == nil {
			extra = jutil.NewLinkedHashSet[string]()
		}
		if !extra.Add(unrecognizedAnnotation) {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Multiple entries with same key: " + unrecognizedAnnotation + "=" +
					Annotation_NOT_IMPLEMENTED.String() + " and " +
					unrecognizedAnnotation + "=" + Annotation_NOT_IMPLEMENTED.String()))
		}
	}
	if extra.IsEmpty() {
		b.cfg.annotations = Annotation_recognizedAnnotations
		b.cfg.extraAnnotations = nil
	} else {
		b.cfg.annotations = nil
		b.cfg.extraAnnotations = extra
	}
	return b
}

// Java: Config.Builder.build()
func (b *Config_Builder) Build() *Config {
	c := b.cfg
	if c.suppressionNames == nil {
		c.suppressionNames = config_emptyStringSet
	}
	if c.closurePrimitiveNames == nil {
		c.closurePrimitiveNames = config_emptyStringSet
	}
	if c.annotations == nil && c.extraAnnotations.IsEmpty() {
		c.annotations = Annotation_recognizedAnnotations
	}
	return &c
}

// SetAnnotations sets the recognized annotations map directly.
//
// Java: Config.Builder.setAnnotations(ImmutableMap<String, Annotation>)
func (b *Config_Builder) SetAnnotations(names map[string]Annotation) *Config_Builder {
	if names == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null annotations"))
	}
	b.cfg.annotations = names
	b.cfg.extraAnnotations = nil
	return b
}

// Config_BuildAnnotations creates the annotation names from the
// user-specified annotation allowlist.
//
// Java: Config.buildAnnotations(Iterable<String>)
func Config_BuildAnnotations(allowlist iter.Seq[string]) map[string]Annotation {
	var extra []string
	for unrecognizedAnnotation := range allowlist {
		if unrecognizedAnnotation == "" {
			continue
		}
		if _, recognized := Annotation_recognizedAnnotations[unrecognizedAnnotation]; !recognized {
			extra = append(extra, unrecognizedAnnotation)
		}
	}
	if len(extra) == 0 {
		return Annotation_recognizedAnnotations
	}
	m := maps.Clone(Annotation_recognizedAnnotations)
	for _, unrecognizedAnnotation := range extra {
		if _, dup := m[unrecognizedAnnotation]; dup {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Multiple entries with same key: " + unrecognizedAnnotation + "=" +
					Annotation_NOT_IMPLEMENTED.String() + " and " +
					unrecognizedAnnotation + "=" + Annotation_NOT_IMPLEMENTED.String()))
		}
		m[unrecognizedAnnotation] = Annotation_NOT_IMPLEMENTED
	}
	return m
}

func config_stringSetCopyOf(names iter.Seq[string]) *jutil.LinkedHashSet[string] {
	if names == nil {
		panic(jutil.NewNullPointerException())
	}
	var set *jutil.LinkedHashSet[string]
	for name := range names {
		if set == nil {
			set = jutil.NewLinkedHashSet[string]()
		}
		set.Add(name)
	}
	if set == nil {
		return config_emptyStringSet
	}
	return set
}
