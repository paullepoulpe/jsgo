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
// Port of com/google/debugging/sourcemap/SourceMapGeneratorFactory.java (generated mechanically by j2go).

package sourcemap

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.debugging.sourcemap.SourceMapGeneratorFactory
type SourceMapGeneratorFactory struct {
	idHash_ atomic.Int32
}

func (this_ *SourceMapGeneratorFactory) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapGeneratorFactory) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapGeneratorFactory) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapGeneratorFactory", "SourceMapGeneratorFactory"
}
func (this_ *SourceMapGeneratorFactory) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapGeneratorFactory) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapGeneratorFactory", "com.google.debugging.sourcemap.SourceMapGeneratorFactory", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapGeneratorFactory](v_)
	})
}

func (this_ *SourceMapGeneratorFactory) initFields_SourceMapGeneratorFactory() {
	_ = this_
}

// Java constructor: SourceMapGeneratorFactory()
func NewSourceMapGeneratorFactory() *SourceMapGeneratorFactory {
	this_ := &SourceMapGeneratorFactory{}
	InitSourceMapGeneratorFactory(this_)
	return this_
}

func InitSourceMapGeneratorFactory(this_ *SourceMapGeneratorFactory) {
	this_.initFields_SourceMapGeneratorFactory()
}

// Java static method: SourceMapGeneratorFactory.getInstance(com.google.debugging.sourcemap.SourceMapFormat format)
func SourceMapGeneratorFactory_GetInstance(format SourceMapFormat) SourceMapGenerator {
	_ = format
	switch (format).Ordinal() {
	case int32(0) /* SourceMapFormat.DEFAULT */, int32(1) /* SourceMapFormat.V3 */ :
		return NewSourceMapGeneratorV3()
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Unexpected format: "+j2gort.StrVal(format))))
	}
	panic("unreachable")
}
