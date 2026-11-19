/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/FeatureSet.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.parser.FeatureSet
type FeatureSet struct {
	lazyFeatures j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"-"`
	bits         int64                                             `gson:"bits"`
}

func (this_ *FeatureSet) BindSelf_(s any) {
	_ = s
}
func (this_ *FeatureSet) IsJavaNil_() bool { return this_ == nil }
func (*FeatureSet) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.FeatureSet", "FeatureSet"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.FeatureSet", "com.google.javascript.jscomp.parsing.parser.FeatureSet", func(v_ any) bool {
		return j2gort.InstanceOf[*FeatureSet](v_)
	})
}

var static_FeatureSet_FEATURE_VALUES []*FeatureSet_Feature

func StaticGet_FeatureSet_FEATURE_VALUES() []*FeatureSet_Feature {
	clinit_FeatureSet()
	return static_FeatureSet_FEATURE_VALUES
}
func StaticSet_FeatureSet_FEATURE_VALUES(v []*FeatureSet_Feature) {
	clinit_FeatureSet()
	static_FeatureSet_FEATURE_VALUES = v
}

var static_FeatureSet_TYPESCRIPT_BITS int64

func StaticGet_FeatureSet_TYPESCRIPT_BITS() int64 {
	clinit_FeatureSet()
	return static_FeatureSet_TYPESCRIPT_BITS
}
func StaticSet_FeatureSet_TYPESCRIPT_BITS(v int64) {
	clinit_FeatureSet()
	static_FeatureSet_TYPESCRIPT_BITS = v
}

var static_FeatureSet_BARE_MINIMUM *FeatureSet

func StaticGet_FeatureSet_BARE_MINIMUM() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_BARE_MINIMUM
}
func StaticSet_FeatureSet_BARE_MINIMUM(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_BARE_MINIMUM = v
}

var static_FeatureSet_ES3 *FeatureSet

func StaticGet_FeatureSet_ES3() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES3
}
func StaticSet_FeatureSet_ES3(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES3 = v
}

var static_FeatureSet_ES5 *FeatureSet

func StaticGet_FeatureSet_ES5() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES5
}
func StaticSet_FeatureSet_ES5(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES5 = v
}

var static_FeatureSet_ES2015_MODULES *FeatureSet

func StaticGet_FeatureSet_ES2015_MODULES() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2015_MODULES
}
func StaticSet_FeatureSet_ES2015_MODULES(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2015_MODULES = v
}

var static_FeatureSet_ES2015 *FeatureSet

func StaticGet_FeatureSet_ES2015() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2015
}
func StaticSet_FeatureSet_ES2015(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2015 = v
}

var static_FeatureSet_ES2016_MODULES *FeatureSet

func StaticGet_FeatureSet_ES2016_MODULES() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2016_MODULES
}
func StaticSet_FeatureSet_ES2016_MODULES(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2016_MODULES = v
}

var static_FeatureSet_ES2016 *FeatureSet

func StaticGet_FeatureSet_ES2016() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2016
}
func StaticSet_FeatureSet_ES2016(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2016 = v
}

var static_FeatureSet_ES2017_MODULES *FeatureSet

func StaticGet_FeatureSet_ES2017_MODULES() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2017_MODULES
}
func StaticSet_FeatureSet_ES2017_MODULES(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2017_MODULES = v
}

var static_FeatureSet_ES2017 *FeatureSet

func StaticGet_FeatureSet_ES2017() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2017
}
func StaticSet_FeatureSet_ES2017(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2017 = v
}

var static_FeatureSet_ES2018_MODULES *FeatureSet

func StaticGet_FeatureSet_ES2018_MODULES() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2018_MODULES
}
func StaticSet_FeatureSet_ES2018_MODULES(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2018_MODULES = v
}

var static_FeatureSet_ES2018 *FeatureSet

func StaticGet_FeatureSet_ES2018() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2018
}
func StaticSet_FeatureSet_ES2018(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2018 = v
}

var static_FeatureSet_ES2019_MODULES *FeatureSet

func StaticGet_FeatureSet_ES2019_MODULES() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2019_MODULES
}
func StaticSet_FeatureSet_ES2019_MODULES(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2019_MODULES = v
}

var static_FeatureSet_ES2019 *FeatureSet

func StaticGet_FeatureSet_ES2019() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2019
}
func StaticSet_FeatureSet_ES2019(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2019 = v
}

var static_FeatureSet_ES2020_MODULES *FeatureSet

func StaticGet_FeatureSet_ES2020_MODULES() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2020_MODULES
}
func StaticSet_FeatureSet_ES2020_MODULES(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2020_MODULES = v
}

var static_FeatureSet_ES2020 *FeatureSet

func StaticGet_FeatureSet_ES2020() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2020
}
func StaticSet_FeatureSet_ES2020(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2020 = v
}

var static_FeatureSet_ES2021_MODULES *FeatureSet

func StaticGet_FeatureSet_ES2021_MODULES() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2021_MODULES
}
func StaticSet_FeatureSet_ES2021_MODULES(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2021_MODULES = v
}

var static_FeatureSet_ES2021 *FeatureSet

func StaticGet_FeatureSet_ES2021() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2021
}
func StaticSet_FeatureSet_ES2021(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2021 = v
}

var static_FeatureSet_ES2022_MODULES *FeatureSet

func StaticGet_FeatureSet_ES2022_MODULES() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2022_MODULES
}
func StaticSet_FeatureSet_ES2022_MODULES(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2022_MODULES = v
}

var static_FeatureSet_ES2022 *FeatureSet

func StaticGet_FeatureSet_ES2022() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES2022
}
func StaticSet_FeatureSet_ES2022(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES2022 = v
}

var static_FeatureSet_ES_NEXT *FeatureSet

func StaticGet_FeatureSet_ES_NEXT() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES_NEXT
}
func StaticSet_FeatureSet_ES_NEXT(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES_NEXT = v
}

var static_FeatureSet_ES_UNSTABLE *FeatureSet

func StaticGet_FeatureSet_ES_UNSTABLE() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES_UNSTABLE
}
func StaticSet_FeatureSet_ES_UNSTABLE(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES_UNSTABLE = v
}

var static_FeatureSet_ES_UNSUPPORTED *FeatureSet

func StaticGet_FeatureSet_ES_UNSUPPORTED() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ES_UNSUPPORTED
}
func StaticSet_FeatureSet_ES_UNSUPPORTED(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ES_UNSUPPORTED = v
}

var static_FeatureSet_BROWSER_2019 *FeatureSet

func StaticGet_FeatureSet_BROWSER_2019() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_BROWSER_2019
}
func StaticSet_FeatureSet_BROWSER_2019(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_BROWSER_2019 = v
}

var static_FeatureSet_BROWSER_2020 *FeatureSet

func StaticGet_FeatureSet_BROWSER_2020() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_BROWSER_2020
}
func StaticSet_FeatureSet_BROWSER_2020(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_BROWSER_2020 = v
}

var static_FeatureSet_BROWSER_2021 *FeatureSet

func StaticGet_FeatureSet_BROWSER_2021() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_BROWSER_2021
}
func StaticSet_FeatureSet_BROWSER_2021(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_BROWSER_2021 = v
}

var static_FeatureSet_BROWSER_2022 *FeatureSet

func StaticGet_FeatureSet_BROWSER_2022() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_BROWSER_2022
}
func StaticSet_FeatureSet_BROWSER_2022(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_BROWSER_2022 = v
}

var static_FeatureSet_BROWSER_2023 *FeatureSet

func StaticGet_FeatureSet_BROWSER_2023() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_BROWSER_2023
}
func StaticSet_FeatureSet_BROWSER_2023(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_BROWSER_2023 = v
}

var static_FeatureSet_BROWSER_2024 *FeatureSet

func StaticGet_FeatureSet_BROWSER_2024() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_BROWSER_2024
}
func StaticSet_FeatureSet_BROWSER_2024(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_BROWSER_2024 = v
}

var static_FeatureSet_BROWSER_2025 *FeatureSet

func StaticGet_FeatureSet_BROWSER_2025() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_BROWSER_2025
}
func StaticSet_FeatureSet_BROWSER_2025(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_BROWSER_2025 = v
}

var static_FeatureSet_BROWSER_2026 *FeatureSet

func StaticGet_FeatureSet_BROWSER_2026() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_BROWSER_2026
}
func StaticSet_FeatureSet_BROWSER_2026(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_BROWSER_2026 = v
}

var static_FeatureSet_ALL *FeatureSet

func StaticGet_FeatureSet_ALL() *FeatureSet {
	clinit_FeatureSet()
	return static_FeatureSet_ALL
}
func StaticSet_FeatureSet_ALL(v *FeatureSet) {
	clinit_FeatureSet()
	static_FeatureSet_ALL = v
}

var clinitDone_FeatureSet uint32
var clinitState_FeatureSet j2gort.ClinitState

func clinit_FeatureSet() {
	if atomic.LoadUint32(&clinitDone_FeatureSet) != 0 {
		return
	}
	clinitSlow_FeatureSet()
}

func clinitSlow_FeatureSet() {
	if !j2gort.EnterClinit(&clinitState_FeatureSet, &clinitDone_FeatureSet) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FeatureSet, &clinitDone_FeatureSet)
	static_FeatureSet_FEATURE_VALUES = FeatureSet_Feature_Values()
	if int32(len(StaticGet_FeatureSet_FEATURE_VALUES())) > int32(64) {
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "FeatureSet bitmask overflow: "+j2gort.StrValInt32(int32(len(StaticGet_FeatureSet_FEATURE_VALUES()))))))
	}
	static_FeatureSet_TYPESCRIPT_BITS = StaticGet_FeatureSet_LangVersion_TYPESCRIPT().priv_LangVersion_featureBits_7226()
	static_FeatureSet_BARE_MINIMUM = NewFeatureSet(int64(0))
	static_FeatureSet_ES3 = StaticGet_FeatureSet_BARE_MINIMUM().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES3().priv_LangVersion_featureBits_7226())
	static_FeatureSet_ES5 = StaticGet_FeatureSet_ES3().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES5().priv_LangVersion_featureBits_7226())
	static_FeatureSet_ES2015_MODULES = StaticGet_FeatureSet_ES5().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES2015().priv_LangVersion_featureBits_7226())
	static_FeatureSet_ES2015 = StaticGet_FeatureSet_ES2015_MODULES().Without_2(StaticGet_FeatureSet_Feature_MODULES(), []*FeatureSet_Feature{})
	static_FeatureSet_ES2016_MODULES = StaticGet_FeatureSet_ES2015_MODULES().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES2016().priv_LangVersion_featureBits_7226())
	static_FeatureSet_ES2016 = StaticGet_FeatureSet_ES2016_MODULES().Without_2(StaticGet_FeatureSet_Feature_MODULES(), []*FeatureSet_Feature{})
	static_FeatureSet_ES2017_MODULES = StaticGet_FeatureSet_ES2016_MODULES().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES2017().priv_LangVersion_featureBits_7226())
	static_FeatureSet_ES2017 = StaticGet_FeatureSet_ES2017_MODULES().Without_2(StaticGet_FeatureSet_Feature_MODULES(), []*FeatureSet_Feature{})
	static_FeatureSet_ES2018_MODULES = StaticGet_FeatureSet_ES2017_MODULES().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES2018().priv_LangVersion_featureBits_7226())
	static_FeatureSet_ES2018 = StaticGet_FeatureSet_ES2018_MODULES().Without_2(StaticGet_FeatureSet_Feature_MODULES(), []*FeatureSet_Feature{})
	static_FeatureSet_ES2019_MODULES = StaticGet_FeatureSet_ES2018_MODULES().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES2019().priv_LangVersion_featureBits_7226())
	static_FeatureSet_ES2019 = StaticGet_FeatureSet_ES2019_MODULES().Without_2(StaticGet_FeatureSet_Feature_MODULES(), []*FeatureSet_Feature{})
	static_FeatureSet_ES2020_MODULES = StaticGet_FeatureSet_ES2019_MODULES().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES2020().priv_LangVersion_featureBits_7226())
	static_FeatureSet_ES2020 = StaticGet_FeatureSet_ES2020_MODULES().Without_2(StaticGet_FeatureSet_Feature_MODULES(), []*FeatureSet_Feature{})
	static_FeatureSet_ES2021_MODULES = StaticGet_FeatureSet_ES2020_MODULES().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES2021().priv_LangVersion_featureBits_7226())
	static_FeatureSet_ES2021 = StaticGet_FeatureSet_ES2021_MODULES().Without_2(StaticGet_FeatureSet_Feature_MODULES(), []*FeatureSet_Feature{})
	static_FeatureSet_ES2022_MODULES = StaticGet_FeatureSet_ES2021_MODULES().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES2022().priv_LangVersion_featureBits_7226())
	static_FeatureSet_ES2022 = StaticGet_FeatureSet_ES2022_MODULES().Without_2(StaticGet_FeatureSet_Feature_MODULES(), []*FeatureSet_Feature{})
	static_FeatureSet_ES_NEXT = StaticGet_FeatureSet_ES2022_MODULES().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES_NEXT().priv_LangVersion_featureBits_7226())
	static_FeatureSet_ES_UNSTABLE = StaticGet_FeatureSet_ES_NEXT().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES_UNSTABLE().priv_LangVersion_featureBits_7226())
	static_FeatureSet_ES_UNSUPPORTED = StaticGet_FeatureSet_ES_UNSTABLE().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_LangVersion_ES_UNSUPPORTED().priv_LangVersion_featureBits_7226())
	static_FeatureSet_BROWSER_2019 = StaticGet_FeatureSet_ES2018_MODULES().Without_2(StaticGet_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD(), []*FeatureSet_Feature{StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST(), StaticGet_FeatureSet_Feature_ASYNC_GENERATORS(), StaticGet_FeatureSet_Feature_FOR_AWAIT_OF(), StaticGet_FeatureSet_Feature_REGEXP_FLAG_S(), StaticGet_FeatureSet_Feature_REGEXP_LOOKBEHIND(), StaticGet_FeatureSet_Feature_REGEXP_NAMED_GROUPS(), StaticGet_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE()})
	static_FeatureSet_BROWSER_2020 = StaticGet_FeatureSet_ES2019_MODULES().Without_2(StaticGet_FeatureSet_Feature_REGEXP_FLAG_S(), []*FeatureSet_Feature{StaticGet_FeatureSet_Feature_REGEXP_NAMED_GROUPS(), StaticGet_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE(), StaticGet_FeatureSet_Feature_REGEXP_LOOKBEHIND()})
	static_FeatureSet_BROWSER_2021 = StaticGet_FeatureSet_ES2020_MODULES().Without_2(StaticGet_FeatureSet_Feature_REGEXP_LOOKBEHIND(), []*FeatureSet_Feature{})
	static_FeatureSet_BROWSER_2022 = StaticGet_FeatureSet_ES2021_MODULES().Without_2(StaticGet_FeatureSet_Feature_REGEXP_LOOKBEHIND(), []*FeatureSet_Feature{})
	static_FeatureSet_BROWSER_2023 = StaticGet_FeatureSet_ES2022_MODULES().Without_2(StaticGet_FeatureSet_Feature_REGEXP_LOOKBEHIND(), []*FeatureSet_Feature{StaticGet_FeatureSet_Feature_CLASS_STATIC_BLOCK()})
	static_FeatureSet_BROWSER_2024 = StaticGet_FeatureSet_ES2022_MODULES()
	static_FeatureSet_BROWSER_2025 = StaticGet_FeatureSet_ES2022_MODULES()
	static_FeatureSet_BROWSER_2026 = StaticGet_FeatureSet_ES2022_MODULES()
	static_FeatureSet_ALL = StaticGet_FeatureSet_ES_UNSUPPORTED().priv_FeatureSet_withBits_7207(StaticGet_FeatureSet_TYPESCRIPT_BITS())
}

func (this_ *FeatureSet) initFields_FeatureSet() {
	_ = this_
}

// Java constructor: FeatureSet(long bits)
func NewFeatureSet(bits int64) *FeatureSet {
	clinit_FeatureSet()
	this_ := &FeatureSet{}
	InitFeatureSet(this_, bits)
	return this_
}

func InitFeatureSet(this_ *FeatureSet, bits int64) {
	_ = bits
	this_.initFields_FeatureSet()
	this_.bits = bits
}

// Java static method: FeatureSet.ofBits(long bits)
func FeatureSet_OfBits(bits int64) *FeatureSet {
	clinit_FeatureSet()
	_ = bits
	if bits == int64(0) {
		return StaticGet_FeatureSet_BARE_MINIMUM()
	} else {
		return NewFeatureSet(bits)
	}
}

// Java method: FeatureSet.withBits(long extraBits)
func (this_ *FeatureSet) priv_FeatureSet_withBits_7207(extraBits int64) *FeatureSet {
	_ = extraBits
	var newBits int64 = (this_.bits | extraBits)
	_ = newBits
	if newBits == this_.bits {
		return this_
	} else {
		return FeatureSet_OfBits(newBits)
	}
}

// Java method: FeatureSet.version()
func (this_ *FeatureSet) Version() string {
	if StaticGet_FeatureSet_ES3().Contains_1_FeatureSet(this_) {
		return "es3"
	}
	if StaticGet_FeatureSet_ES5().Contains_1_FeatureSet(this_) {
		return "es5"
	}
	if StaticGet_FeatureSet_ES2015_MODULES().Contains_1_FeatureSet(this_) {
		return "es6"
	}
	if StaticGet_FeatureSet_ES2016_MODULES().Contains_1_FeatureSet(this_) {
		return "es7"
	}
	if StaticGet_FeatureSet_ES2017_MODULES().Contains_1_FeatureSet(this_) {
		return "es8"
	}
	if StaticGet_FeatureSet_ES2018_MODULES().Contains_1_FeatureSet(this_) {
		return "es9"
	}
	if StaticGet_FeatureSet_ES2019_MODULES().Contains_1_FeatureSet(this_) {
		return "es_2019"
	}
	if StaticGet_FeatureSet_ES2020_MODULES().Contains_1_FeatureSet(this_) {
		return "es_2020"
	}
	if StaticGet_FeatureSet_ES2021_MODULES().Contains_1_FeatureSet(this_) {
		return "es_2021"
	}
	if StaticGet_FeatureSet_ES2022_MODULES().Contains_1_FeatureSet(this_) {
		return "es_2022"
	}
	if StaticGet_FeatureSet_ES_NEXT().Contains_1_FeatureSet(this_) {
		return "es_next"
	}
	if StaticGet_FeatureSet_ES_UNSTABLE().Contains_1_FeatureSet(this_) {
		return "es_unstable"
	}
	if StaticGet_FeatureSet_ES_UNSUPPORTED().Contains_1_FeatureSet(this_) {
		return "es_unsupported"
	}
	if StaticGet_FeatureSet_ALL().Contains_1_FeatureSet(this_) {
		return "all"
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(this_.String()))))
}

// Java method: FeatureSet.versionForDebugging()
func (this_ *FeatureSet) VersionForDebugging() string {
	return this_.Version()
}

// Java method: FeatureSet.without(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature featureToRemove, com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature[] moreFeaturesToRemove)
func (this_ *FeatureSet) Without_2(featureToRemove *FeatureSet_Feature, moreFeaturesToRemove []*FeatureSet_Feature) *FeatureSet {
	_ = featureToRemove
	_ = moreFeaturesToRemove
	var mask int64 = javalang.Shl64(int64(1), int(j2gort.EnumOrdinal(featureToRemove)))
	_ = mask
	for _, el_1 := range j2gort.CheckNotNullSlice(moreFeaturesToRemove) {
		f := el_1
		_ = f
		mask = (mask | javalang.Shl64(int64(1), int(j2gort.EnumOrdinal(f))))
	}
	var newBits int64 = (this_.bits & ^(mask))
	_ = newBits
	if newBits == this_.bits {
		return this_
	} else {
		return FeatureSet_OfBits(newBits)
	}
}

// Java method: FeatureSet.without(com.google.javascript.jscomp.parsing.parser.FeatureSet other)
func (this_ *FeatureSet) Without_1(other *FeatureSet) *FeatureSet {
	_ = other
	var newBits int64 = (this_.bits & ^(other.bits))
	_ = newBits
	if newBits == this_.bits {
		return this_
	} else {
		return FeatureSet_OfBits(newBits)
	}
}

// Java method: FeatureSet.withoutTypes()
func (this_ *FeatureSet) WithoutTypes() *FeatureSet {
	var newBits int64 = (this_.bits & ^(StaticGet_FeatureSet_TYPESCRIPT_BITS()))
	_ = newBits
	if newBits == this_.bits {
		return this_
	} else {
		return FeatureSet_OfBits(newBits)
	}
}

// Java method: FeatureSet.union(com.google.javascript.jscomp.parsing.parser.FeatureSet other)
func (this_ *FeatureSet) Union(other *FeatureSet) *FeatureSet {
	_ = other
	var newBits int64 = (this_.bits | other.bits)
	_ = newBits
	if newBits == this_.bits {
		return this_
	}
	if newBits == other.bits {
		return other
	}
	return FeatureSet_OfBits(newBits)
}

// Java method: FeatureSet.containsAtLeastOneOf(com.google.javascript.jscomp.parsing.parser.FeatureSet other)
func (this_ *FeatureSet) ContainsAtLeastOneOf(other *FeatureSet) bool {
	_ = other
	return ((this_.bits & other.bits) != int64(0))
}

// Java method: FeatureSet.contains(com.google.javascript.jscomp.parsing.parser.FeatureSet other)
func (this_ *FeatureSet) Contains_1_FeatureSet(other *FeatureSet) bool {
	_ = other
	return ((this_.bits & other.bits) == other.bits)
}

// Java method: FeatureSet.contains(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func (this_ *FeatureSet) Contains_1_FeatureSet_Feature(feature *FeatureSet_Feature) bool {
	_ = feature
	return ((this_.bits & (javalang.Shl64(int64(1), int(j2gort.EnumOrdinal(feature))))) != int64(0))
}

// Java method: FeatureSet.with(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func (this_ *FeatureSet) With_1_FeatureSet_Feature(feature *FeatureSet_Feature) *FeatureSet {
	_ = feature
	var newBits int64 = (this_.bits | (javalang.Shl64(int64(1), int(j2gort.EnumOrdinal(feature)))))
	_ = newBits
	if newBits == this_.bits {
		return this_
	} else {
		return FeatureSet_OfBits(newBits)
	}
}

// Java method: FeatureSet.with(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature[] newFeatures)
func (this_ *FeatureSet) With_1_FeatureSet_Feature_7218(newFeatures []*FeatureSet_Feature) *FeatureSet {
	_ = newFeatures
	var newBits int64 = this_.bits
	_ = newBits
	for _, el_2 := range j2gort.CheckNotNullSlice(newFeatures) {
		f := el_2
		_ = f
		newBits = (newBits | javalang.Shl64(int64(1), int(j2gort.EnumOrdinal(f))))
	}
	if newBits == this_.bits {
		return this_
	} else {
		return FeatureSet_OfBits(newBits)
	}
}

// Java method: FeatureSet.with(java.util.Set<com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature> newFeatures)
func (this_ *FeatureSet) With_1_Ext_java_util_Set(newFeatures j2gort.Ext_java_util_Set) *FeatureSet {
	_ = newFeatures
	var newBits int64 = this_.bits
	_ = newBits
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(newFeatures))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			f := j2gort.CastPtr[FeatureSet_Feature](rawElem_4)
			_ = f
			newBits = (newBits | javalang.Shl64(int64(1), int(j2gort.EnumOrdinal(f))))
		}
	}
	if newBits == this_.bits {
		return this_
	} else {
		return FeatureSet_OfBits(newBits)
	}
}

// Java method: FeatureSet.with(com.google.javascript.jscomp.parsing.parser.FeatureSet newFeatures)
func (this_ *FeatureSet) With_1_FeatureSet(newFeatures *FeatureSet) *FeatureSet {
	_ = newFeatures
	return this_.Union(newFeatures)
}

// Java method: FeatureSet.has(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func (this_ *FeatureSet) Has(feature *FeatureSet_Feature) bool {
	_ = feature
	return ((feature != nil) && ((this_.bits & (javalang.Shl64(int64(1), int(j2gort.EnumOrdinal(feature))))) != int64(0)))
}

// Java method: FeatureSet.getFeatures()
func (this_ *FeatureSet) GetFeatures() j2gort.Ext_com_google_common_collect_ImmutableSet {
	var f j2gort.Ext_com_google_common_collect_ImmutableSet = this_.lazyFeatures
	_ = f
	if f == nil {
		if this_.bits == int64(0) {
			f = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
		} else {
			var set j2gort.Ext_java_util_EnumSet = j2gort.Ext_java_util_EnumSet_NoneOf(j2gort.ClassLiteral("com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature", "Feature"))
			_ = set
			var rem int64 = this_.bits
			_ = rem
			for rem != int64(0) {
				var idx int32 = j2gort.BoxLong_NumberOfTrailingZeros(rem)
				_ = idx
				set.Add_1(j2gort.PtrToAny(StaticGet_FeatureSet_FEATURE_VALUES()[idx]))
				rem = (rem & (rem - int64(1)))
			}
			f = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](set))
		}
		this_.lazyFeatures = f
	}
	return f
}

// Java method: FeatureSet.equals(java.lang.Object other)
func (this_ *FeatureSet) Equals(arg_5 any) bool {
	other := arg_5
	_ = other
	var featureSet *FeatureSet
	_ = featureSet
	var ok_6 bool
	featureSet, ok_6 = j2gort.TryCast[*FeatureSet](other)
	return ok_6 && (featureSet.bits == this_.bits)
}

// Java method: FeatureSet.hashCode()
func (this_ *FeatureSet) HashCode() int32 {
	return this_.GetFeatures().HashCode()
}

// Java method: FeatureSet.toString()
func (this_ *FeatureSet) String() string {
	return this_.GetFeatures().String()
}

// Java static method: FeatureSet.valueOf(java.lang.String name)
func FeatureSet_ValueOf(name string) *FeatureSet {
	clinit_FeatureSet()
	_ = name
	switch j2gort.CheckNotNullString(name) {
	case "es3":
		return StaticGet_FeatureSet_ES3()
	case "es5":
		return StaticGet_FeatureSet_ES5()
	case "es_2015", "es6":
		return StaticGet_FeatureSet_ES2015()
	case "es_2016", "es7":
		return StaticGet_FeatureSet_ES2016()
	case "es_2017", "es8":
		return StaticGet_FeatureSet_ES2017()
	case "es_2018", "es9":
		return StaticGet_FeatureSet_ES2018()
	case "es_2019":
		return StaticGet_FeatureSet_ES2019()
	case "es_2020":
		return StaticGet_FeatureSet_ES2020()
	case "es_2021":
		return StaticGet_FeatureSet_ES2021()
	case "es_2022":
		return StaticGet_FeatureSet_ES2022()
	case "es_next":
		return StaticGet_FeatureSet_ES_NEXT()
	case "es_unstable":
		return StaticGet_FeatureSet_ES_UNSTABLE()
	case "es_unsupported":
		return StaticGet_FeatureSet_ES_UNSUPPORTED()
	case "all":
		return StaticGet_FeatureSet_ALL()
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "No such FeatureSet: "+j2gort.StrValString(name))))
	}
	panic("unreachable switch expression")
}

// Java static method: FeatureSet.browserFeatureSetValueOf(java.lang.String name)
func FeatureSet_BrowserFeatureSetValueOf(name string) *FeatureSet {
	clinit_FeatureSet()
	_ = name
	var featureSet *FeatureSet = FeatureSet_ValueOf(name)
	_ = featureSet
	if featureSet.Contains_1_FeatureSet(StaticGet_FeatureSet_ES_NEXT()) {
		return featureSet
	}
	if featureSet.Contains_1_FeatureSet(StaticGet_FeatureSet_ES2022()) {
		return StaticGet_FeatureSet_BROWSER_2023()
	}
	if featureSet.Contains_1_FeatureSet(StaticGet_FeatureSet_ES2021()) {
		return StaticGet_FeatureSet_BROWSER_2022()
	}
	if featureSet.Contains_1_FeatureSet(StaticGet_FeatureSet_ES2020()) {
		return StaticGet_FeatureSet_BROWSER_2021()
	}
	if featureSet.Contains_1_FeatureSet(StaticGet_FeatureSet_ES2019()) {
		return StaticGet_FeatureSet_BROWSER_2020()
	}
	if featureSet.Contains_1_FeatureSet(StaticGet_FeatureSet_ES2018()) {
		return StaticGet_FeatureSet_BROWSER_2019()
	}
	return featureSet
}

// Java static method: FeatureSet.all()
func FeatureSet_All() *FeatureSet {
	clinit_FeatureSet()
	return StaticGet_FeatureSet_ALL()
}

// Java static method: FeatureSet.latest()
func FeatureSet_Latest() *FeatureSet {
	clinit_FeatureSet()
	return StaticGet_FeatureSet_ES_UNSUPPORTED()
}

// Java enum: com.google.javascript.jscomp.parsing.parser.FeatureSet.LangVersion
type FeatureSet_LangVersion uint8

const (
	enum_FeatureSet_LangVersion_ES3            FeatureSet_LangVersion = 1
	enum_FeatureSet_LangVersion_ES5            FeatureSet_LangVersion = 2
	enum_FeatureSet_LangVersion_ES2015         FeatureSet_LangVersion = 3
	enum_FeatureSet_LangVersion_ES2016         FeatureSet_LangVersion = 4
	enum_FeatureSet_LangVersion_ES2017         FeatureSet_LangVersion = 5
	enum_FeatureSet_LangVersion_ES2018         FeatureSet_LangVersion = 6
	enum_FeatureSet_LangVersion_ES2019         FeatureSet_LangVersion = 7
	enum_FeatureSet_LangVersion_ES2020         FeatureSet_LangVersion = 8
	enum_FeatureSet_LangVersion_ES2021         FeatureSet_LangVersion = 9
	enum_FeatureSet_LangVersion_ES2022         FeatureSet_LangVersion = 10
	enum_FeatureSet_LangVersion_ES_NEXT        FeatureSet_LangVersion = 11
	enum_FeatureSet_LangVersion_ES_UNSTABLE    FeatureSet_LangVersion = 12
	enum_FeatureSet_LangVersion_ES_UNSUPPORTED FeatureSet_LangVersion = 13
	enum_FeatureSet_LangVersion_TYPESCRIPT     FeatureSet_LangVersion = 14
)

var enumBaseTable_FeatureSet_LangVersion = [...]j2gort.EnumBase{
	{Name_: "ES3", Ordinal_: 0},
	{Name_: "ES5", Ordinal_: 1},
	{Name_: "ES2015", Ordinal_: 2},
	{Name_: "ES2016", Ordinal_: 3},
	{Name_: "ES2017", Ordinal_: 4},
	{Name_: "ES2018", Ordinal_: 5},
	{Name_: "ES2019", Ordinal_: 6},
	{Name_: "ES2020", Ordinal_: 7},
	{Name_: "ES2021", Ordinal_: 8},
	{Name_: "ES2022", Ordinal_: 9},
	{Name_: "ES_NEXT", Ordinal_: 10},
	{Name_: "ES_UNSTABLE", Ordinal_: 11},
	{Name_: "ES_UNSUPPORTED", Ordinal_: 12},
	{Name_: "TYPESCRIPT", Ordinal_: 13},
}
var enumValues_FeatureSet_LangVersion = []FeatureSet_LangVersion{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

func StaticGet_FeatureSet_LangVersion_ES3() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES3
}
func StaticGet_FeatureSet_LangVersion_ES5() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES5
}
func StaticGet_FeatureSet_LangVersion_ES2015() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES2015
}
func StaticGet_FeatureSet_LangVersion_ES2016() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES2016
}
func StaticGet_FeatureSet_LangVersion_ES2017() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES2017
}
func StaticGet_FeatureSet_LangVersion_ES2018() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES2018
}
func StaticGet_FeatureSet_LangVersion_ES2019() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES2019
}
func StaticGet_FeatureSet_LangVersion_ES2020() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES2020
}
func StaticGet_FeatureSet_LangVersion_ES2021() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES2021
}
func StaticGet_FeatureSet_LangVersion_ES2022() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES2022
}
func StaticGet_FeatureSet_LangVersion_ES_NEXT() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES_NEXT
}
func StaticGet_FeatureSet_LangVersion_ES_UNSTABLE() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES_UNSTABLE
}
func StaticGet_FeatureSet_LangVersion_ES_UNSUPPORTED() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_ES_UNSUPPORTED
}
func StaticGet_FeatureSet_LangVersion_TYPESCRIPT() FeatureSet_LangVersion {
	return enum_FeatureSet_LangVersion_TYPESCRIPT
}

func FeatureSet_LangVersion_Values() []FeatureSet_LangVersion {
	return j2gort.CloneSlice(enumValues_FeatureSet_LangVersion)
}

func FeatureSet_LangVersion_ValueOf(name string) FeatureSet_LangVersion {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_FeatureSet_LangVersion {
		if enumBaseTable_FeatureSet_LangVersion[i_].Name_ == name {
			return FeatureSet_LangVersion(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.FeatureSet.LangVersion." + name))
}

// Java method: LangVersion.featureBits()
func (this_ FeatureSet_LangVersion) priv_LangVersion_featureBits_7226() int64 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	var mask int64 = int64(0)
	_ = mask
	for _, el_7 := range j2gort.CheckNotNullSlice(StaticGet_FeatureSet_FEATURE_VALUES()) {
		feature := el_7
		_ = feature
		if feature.version == this_ {
			mask = (mask | javalang.Shl64(int64(1), int(j2gort.EnumOrdinal(feature))))
		}
	}
	return mask
}

func (this_ FeatureSet_LangVersion) IsJavaNil_() bool                              { return this_ == 0 }
func (FeatureSet_LangVersion) HasIdentityEquals_()                                 {}
func (FeatureSet_LangVersion) Implements_j2gort_Ext_java_lang_Enum()               {}
func (FeatureSet_LangVersion) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ FeatureSet_LangVersion) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (FeatureSet_LangVersion) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.FeatureSet$LangVersion", "LangVersion"
}
func (this_ FeatureSet_LangVersion) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_FeatureSet_LangVersion[this_-1]
}
func (this_ FeatureSet_LangVersion) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ FeatureSet_LangVersion) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ FeatureSet_LangVersion) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FeatureSet_LangVersion[this_-1].Name_
}
func (this_ FeatureSet_LangVersion) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_FeatureSet_LangVersion[this_-1].Name_
}
func (this_ FeatureSet_LangVersion) Equals(other_ any) bool {
	o_, ok_ := other_.(FeatureSet_LangVersion)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ FeatureSet_LangVersion) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_FeatureSet_LangVersion[this_-1].HashCode()
}
func (this_ FeatureSet_LangVersion) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FeatureSet_LangVersion[this_-1].CompareTo(other_)
}
func (this_ FeatureSet_LangVersion) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_FeatureSet_LangVersion[this_-1].Name_)
}
func init() {
	vals_ := []any{FeatureSet_LangVersion(1), FeatureSet_LangVersion(2), FeatureSet_LangVersion(3), FeatureSet_LangVersion(4), FeatureSet_LangVersion(5), FeatureSet_LangVersion(6), FeatureSet_LangVersion(7), FeatureSet_LangVersion(8), FeatureSet_LangVersion(9), FeatureSet_LangVersion(10), FeatureSet_LangVersion(11), FeatureSet_LangVersion(12), FeatureSet_LangVersion(13), FeatureSet_LangVersion(14)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.FeatureSet.LangVersion", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.FeatureSet$LangVersion", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.FeatureSet.LangVersion", "com.google.javascript.jscomp.parsing.parser.FeatureSet$LangVersion", func(v_ any) bool {
		return j2gort.InstanceOf[FeatureSet_LangVersion](v_)
	})
}

// Java enum: com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature
type FeatureSet_Feature struct {
	j2gort.EnumBase
	name    string
	version FeatureSet_LangVersion
}

func (this_ *FeatureSet_Feature) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *FeatureSet_Feature) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*FeatureSet_Feature) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*FeatureSet_Feature) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *FeatureSet_Feature) initFields_FeatureSet_Feature() {
	_ = this_
	this_.name = j2gort.NullString
}

func initEnumFeatureSet_Feature(this_ *FeatureSet_Feature, enumName_ string, enumOrd_ int32, name string, version FeatureSet_LangVersion) *FeatureSet_Feature {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = name
	_ = version
	this_.initFields_FeatureSet_Feature()
	this_.name = name
	this_.version = version
	return this_
}

var enum_FeatureSet_Feature_REGEXP_SYNTAX *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_REGEXP_SYNTAX() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_REGEXP_SYNTAX
}

var enum_FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS
}

var enum_FeatureSet_Feature_GETTER *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_GETTER() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_GETTER
}

var enum_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES
}

var enum_FeatureSet_Feature_SETTER *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_SETTER() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_SETTER
}

var enum_FeatureSet_Feature_STRING_CONTINUATION *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_STRING_CONTINUATION() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_STRING_CONTINUATION
}

var enum_FeatureSet_Feature_TRAILING_COMMA *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_TRAILING_COMMA() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_TRAILING_COMMA
}

var enum_FeatureSet_Feature_ARRAY_DESTRUCTURING *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_ARRAY_DESTRUCTURING() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_ARRAY_DESTRUCTURING
}

var enum_FeatureSet_Feature_ARRAY_PATTERN_REST *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_ARRAY_PATTERN_REST() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_ARRAY_PATTERN_REST
}

var enum_FeatureSet_Feature_ARROW_FUNCTIONS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_ARROW_FUNCTIONS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_ARROW_FUNCTIONS
}

var enum_FeatureSet_Feature_BINARY_LITERALS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_BINARY_LITERALS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_BINARY_LITERALS
}

var enum_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION
}

var enum_FeatureSet_Feature_CLASSES *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_CLASSES() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_CLASSES
}

var enum_FeatureSet_Feature_CLASS_GETTER_SETTER *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_CLASS_GETTER_SETTER
}

var enum_FeatureSet_Feature_COMPUTED_PROPERTIES *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_COMPUTED_PROPERTIES
}

var enum_FeatureSet_Feature_CONST_DECLARATIONS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_CONST_DECLARATIONS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_CONST_DECLARATIONS
}

var enum_FeatureSet_Feature_DEFAULT_PARAMETERS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_DEFAULT_PARAMETERS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_DEFAULT_PARAMETERS
}

var enum_FeatureSet_Feature_FOR_OF *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_FOR_OF() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_FOR_OF
}

var enum_FeatureSet_Feature_GENERATORS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_GENERATORS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_GENERATORS
}

var enum_FeatureSet_Feature_LET_DECLARATIONS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_LET_DECLARATIONS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_LET_DECLARATIONS
}

var enum_FeatureSet_Feature_MEMBER_DECLARATIONS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_MEMBER_DECLARATIONS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_MEMBER_DECLARATIONS
}

var enum_FeatureSet_Feature_NEW_TARGET *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_NEW_TARGET() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_NEW_TARGET
}

var enum_FeatureSet_Feature_OBJECT_DESTRUCTURING *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_OBJECT_DESTRUCTURING() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_OBJECT_DESTRUCTURING
}

var enum_FeatureSet_Feature_OCTAL_LITERALS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_OCTAL_LITERALS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_OCTAL_LITERALS
}

var enum_FeatureSet_Feature_REGEXP_FLAG_U *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_REGEXP_FLAG_U() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_REGEXP_FLAG_U
}

var enum_FeatureSet_Feature_REGEXP_FLAG_Y *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_REGEXP_FLAG_Y() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_REGEXP_FLAG_Y
}

var enum_FeatureSet_Feature_REST_PARAMETERS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_REST_PARAMETERS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_REST_PARAMETERS
}

var enum_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES
}

var enum_FeatureSet_Feature_SPREAD_EXPRESSIONS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_SPREAD_EXPRESSIONS
}

var enum_FeatureSet_Feature_SUPER *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_SUPER() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_SUPER
}

var enum_FeatureSet_Feature_TEMPLATE_LITERALS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_TEMPLATE_LITERALS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_TEMPLATE_LITERALS
}

var enum_FeatureSet_Feature_MODULES *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_MODULES() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_MODULES
}

var enum_FeatureSet_Feature_EXPONENT_OP *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_EXPONENT_OP() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_EXPONENT_OP
}

var enum_FeatureSet_Feature_ASYNC_FUNCTIONS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_ASYNC_FUNCTIONS
}

var enum_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD
}

var enum_FeatureSet_Feature_OBJECT_PATTERN_REST *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_OBJECT_PATTERN_REST
}

var enum_FeatureSet_Feature_ASYNC_GENERATORS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_ASYNC_GENERATORS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_ASYNC_GENERATORS
}

var enum_FeatureSet_Feature_FOR_AWAIT_OF *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_FOR_AWAIT_OF() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_FOR_AWAIT_OF
}

var enum_FeatureSet_Feature_REGEXP_FLAG_S *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_REGEXP_FLAG_S() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_REGEXP_FLAG_S
}

var enum_FeatureSet_Feature_REGEXP_NAMED_GROUPS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_REGEXP_NAMED_GROUPS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_REGEXP_NAMED_GROUPS
}

var enum_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE
}

var enum_FeatureSet_Feature_REGEXP_LOOKBEHIND *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_REGEXP_LOOKBEHIND() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_REGEXP_LOOKBEHIND
}

var enum_FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP
}

var enum_FeatureSet_Feature_OPTIONAL_CATCH_BINDING *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_OPTIONAL_CATCH_BINDING() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_OPTIONAL_CATCH_BINDING
}

var enum_FeatureSet_Feature_DYNAMIC_IMPORT *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_DYNAMIC_IMPORT
}

var enum_FeatureSet_Feature_BIGINT *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_BIGINT() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_BIGINT
}

var enum_FeatureSet_Feature_IMPORT_META *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_IMPORT_META() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_IMPORT_META
}

var enum_FeatureSet_Feature_NULL_COALESCE_OP *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_NULL_COALESCE_OP() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_NULL_COALESCE_OP
}

var enum_FeatureSet_Feature_OPTIONAL_CHAINING *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_OPTIONAL_CHAINING
}

var enum_FeatureSet_Feature_NUMERIC_SEPARATOR *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_NUMERIC_SEPARATOR() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_NUMERIC_SEPARATOR
}

var enum_FeatureSet_Feature_LOGICAL_ASSIGNMENT *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_LOGICAL_ASSIGNMENT
}

var enum_FeatureSet_Feature_PUBLIC_CLASS_FIELDS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_PUBLIC_CLASS_FIELDS
}

var enum_FeatureSet_Feature_CLASS_STATIC_BLOCK *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_CLASS_STATIC_BLOCK() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_CLASS_STATIC_BLOCK
}

var enum_FeatureSet_Feature_REGEXP_FLAG_D *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_REGEXP_FLAG_D() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_REGEXP_FLAG_D
}

var enum_FeatureSet_Feature_TOP_LEVEL_AWAIT *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_TOP_LEVEL_AWAIT() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_TOP_LEVEL_AWAIT
}

var enum_FeatureSet_Feature_ES_NEXT_RUNTIME *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_ES_NEXT_RUNTIME() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_ES_NEXT_RUNTIME
}

var enum_FeatureSet_Feature_ES_UNSTABLE_RUNTIME *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_ES_UNSTABLE_RUNTIME() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_ES_UNSTABLE_RUNTIME
}

var enum_FeatureSet_Feature_PRIVATE_ELEMENTS *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_PRIVATE_ELEMENTS
}

var enum_FeatureSet_Feature_TYPE_ANNOTATION *FeatureSet_Feature

func StaticGet_FeatureSet_Feature_TYPE_ANNOTATION() *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return enum_FeatureSet_Feature_TYPE_ANNOTATION
}

var clinitDone_FeatureSet_Feature uint32
var clinitState_FeatureSet_Feature j2gort.ClinitState

func clinit_FeatureSet_Feature() {
	if atomic.LoadUint32(&clinitDone_FeatureSet_Feature) != 0 {
		return
	}
	clinitSlow_FeatureSet_Feature()
}

func clinitSlow_FeatureSet_Feature() {
	if !j2gort.EnterClinit(&clinitState_FeatureSet_Feature, &clinitDone_FeatureSet_Feature) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FeatureSet_Feature, &clinitDone_FeatureSet_Feature)
	enum_FeatureSet_Feature_REGEXP_SYNTAX = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_REGEXP_SYNTAX, "REGEXP_SYNTAX", int32(0), "RegExp syntax", StaticGet_FeatureSet_LangVersion_ES3())
	enum_FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS, "ES3_KEYWORDS_AS_IDENTIFIERS", int32(1), "ES3 keywords as identifiers", StaticGet_FeatureSet_LangVersion_ES5())
	enum_FeatureSet_Feature_GETTER = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_GETTER, "GETTER", int32(2), "getters", StaticGet_FeatureSet_LangVersion_ES5())
	enum_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES, "KEYWORDS_AS_PROPERTIES", int32(3), "reserved words as properties", StaticGet_FeatureSet_LangVersion_ES5())
	enum_FeatureSet_Feature_SETTER = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_SETTER, "SETTER", int32(4), "setters", StaticGet_FeatureSet_LangVersion_ES5())
	enum_FeatureSet_Feature_STRING_CONTINUATION = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_STRING_CONTINUATION, "STRING_CONTINUATION", int32(5), "string continuation", StaticGet_FeatureSet_LangVersion_ES5())
	enum_FeatureSet_Feature_TRAILING_COMMA = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_TRAILING_COMMA, "TRAILING_COMMA", int32(6), "trailing comma", StaticGet_FeatureSet_LangVersion_ES5())
	enum_FeatureSet_Feature_ARRAY_DESTRUCTURING = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_ARRAY_DESTRUCTURING, "ARRAY_DESTRUCTURING", int32(7), "array destructuring", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_ARRAY_PATTERN_REST = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_ARRAY_PATTERN_REST, "ARRAY_PATTERN_REST", int32(8), "array pattern rest", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_ARROW_FUNCTIONS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_ARROW_FUNCTIONS, "ARROW_FUNCTIONS", int32(9), "arrow function", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_BINARY_LITERALS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_BINARY_LITERALS, "BINARY_LITERALS", int32(10), "binary literal", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION, "BLOCK_SCOPED_FUNCTION_DECLARATION", int32(11), "block-scoped function declaration", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_CLASSES = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_CLASSES, "CLASSES", int32(12), "class", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_CLASS_GETTER_SETTER = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_CLASS_GETTER_SETTER, "CLASS_GETTER_SETTER", int32(13), "class getters/setters", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_COMPUTED_PROPERTIES = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_COMPUTED_PROPERTIES, "COMPUTED_PROPERTIES", int32(14), "computed property", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_CONST_DECLARATIONS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_CONST_DECLARATIONS, "CONST_DECLARATIONS", int32(15), "const declaration", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_DEFAULT_PARAMETERS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_DEFAULT_PARAMETERS, "DEFAULT_PARAMETERS", int32(16), "default parameter", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_FOR_OF = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_FOR_OF, "FOR_OF", int32(17), "for-of loop", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_GENERATORS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_GENERATORS, "GENERATORS", int32(18), "generator", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_LET_DECLARATIONS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_LET_DECLARATIONS, "LET_DECLARATIONS", int32(19), "let declaration", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_MEMBER_DECLARATIONS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_MEMBER_DECLARATIONS, "MEMBER_DECLARATIONS", int32(20), "member declaration", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_NEW_TARGET = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_NEW_TARGET, "NEW_TARGET", int32(21), "new.target", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_OBJECT_DESTRUCTURING = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_OBJECT_DESTRUCTURING, "OBJECT_DESTRUCTURING", int32(22), "object destructuring", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_OCTAL_LITERALS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_OCTAL_LITERALS, "OCTAL_LITERALS", int32(23), "octal literal", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_REGEXP_FLAG_U = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_REGEXP_FLAG_U, "REGEXP_FLAG_U", int32(24), "RegExp flag 'u'", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_REGEXP_FLAG_Y = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_REGEXP_FLAG_Y, "REGEXP_FLAG_Y", int32(25), "RegExp flag 'y'", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_REST_PARAMETERS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_REST_PARAMETERS, "REST_PARAMETERS", int32(26), "rest parameter", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES, "SHORTHAND_OBJECT_PROPERTIES", int32(27), "shorthand object property", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_SPREAD_EXPRESSIONS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_SPREAD_EXPRESSIONS, "SPREAD_EXPRESSIONS", int32(28), "spread expression", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_SUPER = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_SUPER, "SUPER", int32(29), "super", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_TEMPLATE_LITERALS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_TEMPLATE_LITERALS, "TEMPLATE_LITERALS", int32(30), "template literal", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_MODULES = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_MODULES, "MODULES", int32(31), "modules", StaticGet_FeatureSet_LangVersion_ES2015())
	enum_FeatureSet_Feature_EXPONENT_OP = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_EXPONENT_OP, "EXPONENT_OP", int32(32), "exponent operator (**)", StaticGet_FeatureSet_LangVersion_ES2016())
	enum_FeatureSet_Feature_ASYNC_FUNCTIONS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_ASYNC_FUNCTIONS, "ASYNC_FUNCTIONS", int32(33), "async function", StaticGet_FeatureSet_LangVersion_ES2017())
	enum_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD, "OBJECT_LITERALS_WITH_SPREAD", int32(34), "object literals with spread", StaticGet_FeatureSet_LangVersion_ES2018())
	enum_FeatureSet_Feature_OBJECT_PATTERN_REST = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_OBJECT_PATTERN_REST, "OBJECT_PATTERN_REST", int32(35), "object pattern rest", StaticGet_FeatureSet_LangVersion_ES2018())
	enum_FeatureSet_Feature_ASYNC_GENERATORS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_ASYNC_GENERATORS, "ASYNC_GENERATORS", int32(36), "async generator functions", StaticGet_FeatureSet_LangVersion_ES2018())
	enum_FeatureSet_Feature_FOR_AWAIT_OF = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_FOR_AWAIT_OF, "FOR_AWAIT_OF", int32(37), "for-await-of loop", StaticGet_FeatureSet_LangVersion_ES2018())
	enum_FeatureSet_Feature_REGEXP_FLAG_S = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_REGEXP_FLAG_S, "REGEXP_FLAG_S", int32(38), "RegExp flag 's'", StaticGet_FeatureSet_LangVersion_ES2018())
	enum_FeatureSet_Feature_REGEXP_NAMED_GROUPS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_REGEXP_NAMED_GROUPS, "REGEXP_NAMED_GROUPS", int32(39), "RegExp named groups", StaticGet_FeatureSet_LangVersion_ES2018())
	enum_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE, "REGEXP_UNICODE_PROPERTY_ESCAPE", int32(40), "RegExp unicode property escape", StaticGet_FeatureSet_LangVersion_ES2018())
	enum_FeatureSet_Feature_REGEXP_LOOKBEHIND = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_REGEXP_LOOKBEHIND, "REGEXP_LOOKBEHIND", int32(41), "RegExp Lookbehind", StaticGet_FeatureSet_LangVersion_ES2018())
	enum_FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP, "UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP", int32(42), "Unescaped unicode line or paragraph separator", StaticGet_FeatureSet_LangVersion_ES2019())
	enum_FeatureSet_Feature_OPTIONAL_CATCH_BINDING = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_OPTIONAL_CATCH_BINDING, "OPTIONAL_CATCH_BINDING", int32(43), "Optional catch binding", StaticGet_FeatureSet_LangVersion_ES2019())
	enum_FeatureSet_Feature_DYNAMIC_IMPORT = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_DYNAMIC_IMPORT, "DYNAMIC_IMPORT", int32(44), "Dynamic module import", StaticGet_FeatureSet_LangVersion_ES2020())
	enum_FeatureSet_Feature_BIGINT = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_BIGINT, "BIGINT", int32(45), "bigint", StaticGet_FeatureSet_LangVersion_ES2020())
	enum_FeatureSet_Feature_IMPORT_META = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_IMPORT_META, "IMPORT_META", int32(46), "import.meta", StaticGet_FeatureSet_LangVersion_ES2020())
	enum_FeatureSet_Feature_NULL_COALESCE_OP = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_NULL_COALESCE_OP, "NULL_COALESCE_OP", int32(47), "Nullish coalescing", StaticGet_FeatureSet_LangVersion_ES2020())
	enum_FeatureSet_Feature_OPTIONAL_CHAINING = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_OPTIONAL_CHAINING, "OPTIONAL_CHAINING", int32(48), "Optional chaining", StaticGet_FeatureSet_LangVersion_ES2020())
	enum_FeatureSet_Feature_NUMERIC_SEPARATOR = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_NUMERIC_SEPARATOR, "NUMERIC_SEPARATOR", int32(49), "numeric separator", StaticGet_FeatureSet_LangVersion_ES2021())
	enum_FeatureSet_Feature_LOGICAL_ASSIGNMENT = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_LOGICAL_ASSIGNMENT, "LOGICAL_ASSIGNMENT", int32(50), "Logical assignments", StaticGet_FeatureSet_LangVersion_ES2021())
	enum_FeatureSet_Feature_PUBLIC_CLASS_FIELDS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_PUBLIC_CLASS_FIELDS, "PUBLIC_CLASS_FIELDS", int32(51), "Public class fields", StaticGet_FeatureSet_LangVersion_ES2022())
	enum_FeatureSet_Feature_CLASS_STATIC_BLOCK = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_CLASS_STATIC_BLOCK, "CLASS_STATIC_BLOCK", int32(52), "Class static block", StaticGet_FeatureSet_LangVersion_ES2022())
	enum_FeatureSet_Feature_REGEXP_FLAG_D = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_REGEXP_FLAG_D, "REGEXP_FLAG_D", int32(53), "RegExp flag 'd'", StaticGet_FeatureSet_LangVersion_ES2022())
	enum_FeatureSet_Feature_TOP_LEVEL_AWAIT = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_TOP_LEVEL_AWAIT, "TOP_LEVEL_AWAIT", int32(54), "Top-level await", StaticGet_FeatureSet_LangVersion_ES_UNSUPPORTED())
	enum_FeatureSet_Feature_ES_NEXT_RUNTIME = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_ES_NEXT_RUNTIME, "ES_NEXT_RUNTIME", int32(55), "es_next runtime", StaticGet_FeatureSet_LangVersion_ES_NEXT())
	enum_FeatureSet_Feature_ES_UNSTABLE_RUNTIME = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_ES_UNSTABLE_RUNTIME, "ES_UNSTABLE_RUNTIME", int32(56), "es_unstable runtime", StaticGet_FeatureSet_LangVersion_ES_UNSTABLE())
	enum_FeatureSet_Feature_PRIVATE_ELEMENTS = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_PRIVATE_ELEMENTS, "PRIVATE_ELEMENTS", int32(57), "Private elements", StaticGet_FeatureSet_LangVersion_ES_UNSUPPORTED())
	enum_FeatureSet_Feature_TYPE_ANNOTATION = &FeatureSet_Feature{}
	initEnumFeatureSet_Feature(enum_FeatureSet_Feature_TYPE_ANNOTATION, "TYPE_ANNOTATION", int32(58), "type annotation", StaticGet_FeatureSet_LangVersion_TYPESCRIPT())
	enumValues_FeatureSet_Feature = []*FeatureSet_Feature{enum_FeatureSet_Feature_REGEXP_SYNTAX, enum_FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS, enum_FeatureSet_Feature_GETTER, enum_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES, enum_FeatureSet_Feature_SETTER, enum_FeatureSet_Feature_STRING_CONTINUATION, enum_FeatureSet_Feature_TRAILING_COMMA, enum_FeatureSet_Feature_ARRAY_DESTRUCTURING, enum_FeatureSet_Feature_ARRAY_PATTERN_REST, enum_FeatureSet_Feature_ARROW_FUNCTIONS, enum_FeatureSet_Feature_BINARY_LITERALS, enum_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION, enum_FeatureSet_Feature_CLASSES, enum_FeatureSet_Feature_CLASS_GETTER_SETTER, enum_FeatureSet_Feature_COMPUTED_PROPERTIES, enum_FeatureSet_Feature_CONST_DECLARATIONS, enum_FeatureSet_Feature_DEFAULT_PARAMETERS, enum_FeatureSet_Feature_FOR_OF, enum_FeatureSet_Feature_GENERATORS, enum_FeatureSet_Feature_LET_DECLARATIONS, enum_FeatureSet_Feature_MEMBER_DECLARATIONS, enum_FeatureSet_Feature_NEW_TARGET, enum_FeatureSet_Feature_OBJECT_DESTRUCTURING, enum_FeatureSet_Feature_OCTAL_LITERALS, enum_FeatureSet_Feature_REGEXP_FLAG_U, enum_FeatureSet_Feature_REGEXP_FLAG_Y, enum_FeatureSet_Feature_REST_PARAMETERS, enum_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES, enum_FeatureSet_Feature_SPREAD_EXPRESSIONS, enum_FeatureSet_Feature_SUPER, enum_FeatureSet_Feature_TEMPLATE_LITERALS, enum_FeatureSet_Feature_MODULES, enum_FeatureSet_Feature_EXPONENT_OP, enum_FeatureSet_Feature_ASYNC_FUNCTIONS, enum_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD, enum_FeatureSet_Feature_OBJECT_PATTERN_REST, enum_FeatureSet_Feature_ASYNC_GENERATORS, enum_FeatureSet_Feature_FOR_AWAIT_OF, enum_FeatureSet_Feature_REGEXP_FLAG_S, enum_FeatureSet_Feature_REGEXP_NAMED_GROUPS, enum_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE, enum_FeatureSet_Feature_REGEXP_LOOKBEHIND, enum_FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP, enum_FeatureSet_Feature_OPTIONAL_CATCH_BINDING, enum_FeatureSet_Feature_DYNAMIC_IMPORT, enum_FeatureSet_Feature_BIGINT, enum_FeatureSet_Feature_IMPORT_META, enum_FeatureSet_Feature_NULL_COALESCE_OP, enum_FeatureSet_Feature_OPTIONAL_CHAINING, enum_FeatureSet_Feature_NUMERIC_SEPARATOR, enum_FeatureSet_Feature_LOGICAL_ASSIGNMENT, enum_FeatureSet_Feature_PUBLIC_CLASS_FIELDS, enum_FeatureSet_Feature_CLASS_STATIC_BLOCK, enum_FeatureSet_Feature_REGEXP_FLAG_D, enum_FeatureSet_Feature_TOP_LEVEL_AWAIT, enum_FeatureSet_Feature_ES_NEXT_RUNTIME, enum_FeatureSet_Feature_ES_UNSTABLE_RUNTIME, enum_FeatureSet_Feature_PRIVATE_ELEMENTS, enum_FeatureSet_Feature_TYPE_ANNOTATION}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature", []any{enum_FeatureSet_Feature_REGEXP_SYNTAX, enum_FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS, enum_FeatureSet_Feature_GETTER, enum_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES, enum_FeatureSet_Feature_SETTER, enum_FeatureSet_Feature_STRING_CONTINUATION, enum_FeatureSet_Feature_TRAILING_COMMA, enum_FeatureSet_Feature_ARRAY_DESTRUCTURING, enum_FeatureSet_Feature_ARRAY_PATTERN_REST, enum_FeatureSet_Feature_ARROW_FUNCTIONS, enum_FeatureSet_Feature_BINARY_LITERALS, enum_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION, enum_FeatureSet_Feature_CLASSES, enum_FeatureSet_Feature_CLASS_GETTER_SETTER, enum_FeatureSet_Feature_COMPUTED_PROPERTIES, enum_FeatureSet_Feature_CONST_DECLARATIONS, enum_FeatureSet_Feature_DEFAULT_PARAMETERS, enum_FeatureSet_Feature_FOR_OF, enum_FeatureSet_Feature_GENERATORS, enum_FeatureSet_Feature_LET_DECLARATIONS, enum_FeatureSet_Feature_MEMBER_DECLARATIONS, enum_FeatureSet_Feature_NEW_TARGET, enum_FeatureSet_Feature_OBJECT_DESTRUCTURING, enum_FeatureSet_Feature_OCTAL_LITERALS, enum_FeatureSet_Feature_REGEXP_FLAG_U, enum_FeatureSet_Feature_REGEXP_FLAG_Y, enum_FeatureSet_Feature_REST_PARAMETERS, enum_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES, enum_FeatureSet_Feature_SPREAD_EXPRESSIONS, enum_FeatureSet_Feature_SUPER, enum_FeatureSet_Feature_TEMPLATE_LITERALS, enum_FeatureSet_Feature_MODULES, enum_FeatureSet_Feature_EXPONENT_OP, enum_FeatureSet_Feature_ASYNC_FUNCTIONS, enum_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD, enum_FeatureSet_Feature_OBJECT_PATTERN_REST, enum_FeatureSet_Feature_ASYNC_GENERATORS, enum_FeatureSet_Feature_FOR_AWAIT_OF, enum_FeatureSet_Feature_REGEXP_FLAG_S, enum_FeatureSet_Feature_REGEXP_NAMED_GROUPS, enum_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE, enum_FeatureSet_Feature_REGEXP_LOOKBEHIND, enum_FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP, enum_FeatureSet_Feature_OPTIONAL_CATCH_BINDING, enum_FeatureSet_Feature_DYNAMIC_IMPORT, enum_FeatureSet_Feature_BIGINT, enum_FeatureSet_Feature_IMPORT_META, enum_FeatureSet_Feature_NULL_COALESCE_OP, enum_FeatureSet_Feature_OPTIONAL_CHAINING, enum_FeatureSet_Feature_NUMERIC_SEPARATOR, enum_FeatureSet_Feature_LOGICAL_ASSIGNMENT, enum_FeatureSet_Feature_PUBLIC_CLASS_FIELDS, enum_FeatureSet_Feature_CLASS_STATIC_BLOCK, enum_FeatureSet_Feature_REGEXP_FLAG_D, enum_FeatureSet_Feature_TOP_LEVEL_AWAIT, enum_FeatureSet_Feature_ES_NEXT_RUNTIME, enum_FeatureSet_Feature_ES_UNSTABLE_RUNTIME, enum_FeatureSet_Feature_PRIVATE_ELEMENTS, enum_FeatureSet_Feature_TYPE_ANNOTATION})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.FeatureSet$Feature", []any{enum_FeatureSet_Feature_REGEXP_SYNTAX, enum_FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS, enum_FeatureSet_Feature_GETTER, enum_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES, enum_FeatureSet_Feature_SETTER, enum_FeatureSet_Feature_STRING_CONTINUATION, enum_FeatureSet_Feature_TRAILING_COMMA, enum_FeatureSet_Feature_ARRAY_DESTRUCTURING, enum_FeatureSet_Feature_ARRAY_PATTERN_REST, enum_FeatureSet_Feature_ARROW_FUNCTIONS, enum_FeatureSet_Feature_BINARY_LITERALS, enum_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION, enum_FeatureSet_Feature_CLASSES, enum_FeatureSet_Feature_CLASS_GETTER_SETTER, enum_FeatureSet_Feature_COMPUTED_PROPERTIES, enum_FeatureSet_Feature_CONST_DECLARATIONS, enum_FeatureSet_Feature_DEFAULT_PARAMETERS, enum_FeatureSet_Feature_FOR_OF, enum_FeatureSet_Feature_GENERATORS, enum_FeatureSet_Feature_LET_DECLARATIONS, enum_FeatureSet_Feature_MEMBER_DECLARATIONS, enum_FeatureSet_Feature_NEW_TARGET, enum_FeatureSet_Feature_OBJECT_DESTRUCTURING, enum_FeatureSet_Feature_OCTAL_LITERALS, enum_FeatureSet_Feature_REGEXP_FLAG_U, enum_FeatureSet_Feature_REGEXP_FLAG_Y, enum_FeatureSet_Feature_REST_PARAMETERS, enum_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES, enum_FeatureSet_Feature_SPREAD_EXPRESSIONS, enum_FeatureSet_Feature_SUPER, enum_FeatureSet_Feature_TEMPLATE_LITERALS, enum_FeatureSet_Feature_MODULES, enum_FeatureSet_Feature_EXPONENT_OP, enum_FeatureSet_Feature_ASYNC_FUNCTIONS, enum_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD, enum_FeatureSet_Feature_OBJECT_PATTERN_REST, enum_FeatureSet_Feature_ASYNC_GENERATORS, enum_FeatureSet_Feature_FOR_AWAIT_OF, enum_FeatureSet_Feature_REGEXP_FLAG_S, enum_FeatureSet_Feature_REGEXP_NAMED_GROUPS, enum_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE, enum_FeatureSet_Feature_REGEXP_LOOKBEHIND, enum_FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP, enum_FeatureSet_Feature_OPTIONAL_CATCH_BINDING, enum_FeatureSet_Feature_DYNAMIC_IMPORT, enum_FeatureSet_Feature_BIGINT, enum_FeatureSet_Feature_IMPORT_META, enum_FeatureSet_Feature_NULL_COALESCE_OP, enum_FeatureSet_Feature_OPTIONAL_CHAINING, enum_FeatureSet_Feature_NUMERIC_SEPARATOR, enum_FeatureSet_Feature_LOGICAL_ASSIGNMENT, enum_FeatureSet_Feature_PUBLIC_CLASS_FIELDS, enum_FeatureSet_Feature_CLASS_STATIC_BLOCK, enum_FeatureSet_Feature_REGEXP_FLAG_D, enum_FeatureSet_Feature_TOP_LEVEL_AWAIT, enum_FeatureSet_Feature_ES_NEXT_RUNTIME, enum_FeatureSet_Feature_ES_UNSTABLE_RUNTIME, enum_FeatureSet_Feature_PRIVATE_ELEMENTS, enum_FeatureSet_Feature_TYPE_ANNOTATION})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature", "com.google.javascript.jscomp.parsing.parser.FeatureSet$Feature", func() {
		clinit_FeatureSet_Feature()
	})
}

var enumValues_FeatureSet_Feature []*FeatureSet_Feature

func FeatureSet_Feature_Values() []*FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	return j2gort.CloneSlice(enumValues_FeatureSet_Feature)
}

func FeatureSet_Feature_ValueOf(name string) *FeatureSet_Feature {
	clinit_FeatureSet_Feature()
	for _, c := range enumValues_FeatureSet_Feature {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature." + name))
}

// Java method: Feature.toString()
func (this_ *FeatureSet_Feature) String() string {
	return this_.name
}

func (this_ *FeatureSet_Feature) IsJavaNil_() bool { return this_ == nil }
func (*FeatureSet_Feature) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.FeatureSet$Feature", "Feature"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature", "com.google.javascript.jscomp.parsing.parser.FeatureSet$Feature", func(v_ any) bool {
		return j2gort.InstanceOf[*FeatureSet_Feature](v_)
	})
}
