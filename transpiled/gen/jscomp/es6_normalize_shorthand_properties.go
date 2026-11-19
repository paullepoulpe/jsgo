/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/Es6NormalizeShorthandProperties.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Es6NormalizeShorthandProperties
type Es6NormalizeShorthandProperties struct {
	AbstractPeepholeTranspilationBase
}

func (this_ *Es6NormalizeShorthandProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6NormalizeShorthandProperties) IsJavaNil_() bool { return this_ == nil }
func (*Es6NormalizeShorthandProperties) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.Es6NormalizeShorthandProperties", "Es6NormalizeShorthandProperties" }
func (this_ *Es6NormalizeShorthandProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeTranspilationBase.IdentityHash_()
}
func (*Es6NormalizeShorthandProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6NormalizeShorthandProperties", "com.google.javascript.jscomp.Es6NormalizeShorthandProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6NormalizeShorthandProperties](v_)
	})
}
func (*Es6NormalizeShorthandProperties) Implements_jscomp_AbstractPeepholeTranspilation() {}

func (this_ *Es6NormalizeShorthandProperties) initFields_Es6NormalizeShorthandProperties() {
	_ = this_
}

// Java constructor: Es6NormalizeShorthandProperties(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6NormalizeShorthandProperties(compiler AbstractCompiler) *Es6NormalizeShorthandProperties {
	this_ := &Es6NormalizeShorthandProperties{}
	InitEs6NormalizeShorthandProperties(this_, compiler)
	return this_
}

func InitEs6NormalizeShorthandProperties(this_ *Es6NormalizeShorthandProperties, compiler AbstractCompiler) {
	_ = compiler
	InitAbstractPeepholeTranspilation(&this_.AbstractPeepholeTranspilationBase)
	this_.initFields_Es6NormalizeShorthandProperties()
}

// Java method: Es6NormalizeShorthandProperties.transpileSubtree(com.google.javascript.rhino.Node subtree)
func (this_ *Es6NormalizeShorthandProperties) TranspileSubtree(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	if (subtree.IsStringKey()) {
		subtree.SetShorthandProperty(false)
	}
	return subtree
}

// Java method: Es6NormalizeShorthandProperties.getTranspiledAwayFeatures()
func (this_ *Es6NormalizeShorthandProperties) GetTranspiledAwayFeatures() *pkg_parser.FeatureSet {
	return pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES())
}

