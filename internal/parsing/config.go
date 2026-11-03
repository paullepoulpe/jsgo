// Port of com/google/javascript/jscomp/parsing/Config.java @ closure-compiler v20261028

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
)

// Config is the configuration for the AST factory. Should be shared across
// AST creation for all files of a compilation process.
type Config struct {
	languageMode          Config_LanguageMode
	strictMode            Config_StrictMode
	jsDocParsingMode      Config_JsDocParsing
	runMode               Config_RunMode
	annotations           map[string]Annotation
	extraAnnotations      []string
	suppressionNames      map[string]struct{}
	closurePrimitiveNames map[string]struct{}
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
var config_LanguageMode_featureSet = [...]parser.FeatureSet{
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
func (m Config_LanguageMode) FeatureSet() parser.FeatureSet {
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
func Config_LanguageMode_MinimumRequiredForSet(featureSet parser.FeatureSet) Config_LanguageMode {
	for mode := Config_LanguageMode(0); int(mode) < len(config_LanguageMode_names); mode++ {
		if mode.FeatureSet().Contains(featureSet) {
			return mode
		}
	}
	panic(jutil.NewIllegalStateExceptionMsg("No input language mode supports feature set: " + featureSet.String()))
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

var config_JsDocParsing_all = [...]Config_JsDocParsing{
	Config_JsDocParsing_TYPES_ONLY,
	Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE,
	Config_JsDocParsing_INCLUDE_DESCRIPTIONS_WITH_WHITESPACE,
	Config_JsDocParsing_INCLUDE_ALL_COMMENTS,
	Config_JsDocParsing_LICENSE_COMMENTS_ONLY,
}

// Config_JsDocParsing_values is JsDocParsing.values().
//
// Java: Config.JsDocParsing.values()
func Config_JsDocParsing_values() []Config_JsDocParsing {
	return config_JsDocParsing_all[:]
}

// Name is Enum.name().
//
// Java: Config.JsDocParsing.name()
func (p Config_JsDocParsing) Name() string { return config_JsDocParsing_names[p] }

// String is Enum.toString().
//
// Java: Config.JsDocParsing.toString()
func (p Config_JsDocParsing) String() string { return config_JsDocParsing_names[p] }

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
	if slices.Contains(c.extraAnnotations, name) {
		return Annotation_NOT_IMPLEMENTED, true
	}
	return 0, false
}

// SuppressionNames is the set of recognized names in a `@suppress` tag.
//
// Java: Config.suppressionNames()
func (c *Config) SuppressionNames() map[string]struct{} {
	return c.suppressionNames
}

// ClosurePrimitiveNames is the set of recognized names in a
// `@closurePrimitive` tag.
//
// Java: Config.closurePrimitiveNames()
func (c *Config) ClosurePrimitiveNames() map[string]struct{} {
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
	if len(c.extraAnnotations) == 0 {
		return slices.Clone(annotation_recognizedOrder)
	}
	out := make([]string, 0, len(annotation_recognizedOrder)+len(c.extraAnnotations))
	out = append(out, annotation_recognizedOrder...)
	out = append(out, c.extraAnnotations...)
	return out
}

// Java: Config.toBuilder()
func (c *Config) ToBuilder() *Config_Builder {
	return &Config_Builder{cfg: *c}
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
			parseInlineSourceMaps: false,
		},
	}
}

// Config_Builder is the builder for a Config.
type Config_Builder struct {
	cfg Config
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
	var extra []string
	for unrecognizedAnnotation := range names {
		if unrecognizedAnnotation == "" {
			continue
		}
		if _, recognized := Annotation_recognizedAnnotations[unrecognizedAnnotation]; recognized {
			continue
		}
		if slices.Contains(extra, unrecognizedAnnotation) {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Multiple entries with same key: " + unrecognizedAnnotation + "=" +
					Annotation_NOT_IMPLEMENTED.String() + " and " +
					unrecognizedAnnotation + "=" + Annotation_NOT_IMPLEMENTED.String()))
		}
		extra = append(extra, unrecognizedAnnotation)
	}
	b.cfg.extraAnnotations = extra
	if len(extra) == 0 {
		b.cfg.annotations = Annotation_recognizedAnnotations
	} else {
		b.cfg.annotations = nil
	}
	return b
}

// Java: Config.Builder.build()
func (b *Config_Builder) Build() *Config {
	c := b.cfg
	if c.annotations == nil && len(c.extraAnnotations) == 0 {
		c.annotations = Annotation_recognizedAnnotations
	}
	return &c
}

func config_stringSetCopyOf(names iter.Seq[string]) map[string]struct{} {
	if names == nil {
		panic(jutil.NewNullPointerException())
	}
	var set map[string]struct{}
	for name := range names {
		if set == nil {
			set = make(map[string]struct{})
		}
		set[name] = struct{}{}
	}
	return set
}
