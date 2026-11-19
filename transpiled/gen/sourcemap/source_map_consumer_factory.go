/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/debugging/sourcemap/SourceMapConsumerFactory.java (generated mechanically by j2go).

package sourcemap

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.debugging.sourcemap.SourceMapConsumerFactory
type SourceMapConsumerFactory struct {
	idHash_ atomic.Int32
}

func (this_ *SourceMapConsumerFactory) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapConsumerFactory) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapConsumerFactory) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapConsumerFactory", "SourceMapConsumerFactory"
}
func (this_ *SourceMapConsumerFactory) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapConsumerFactory) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapConsumerFactory", "com.google.debugging.sourcemap.SourceMapConsumerFactory", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapConsumerFactory](v_)
	})
}

func (this_ *SourceMapConsumerFactory) initFields_SourceMapConsumerFactory() {
	_ = this_
}

// Java constructor: SourceMapConsumerFactory()
func NewSourceMapConsumerFactory() *SourceMapConsumerFactory {
	this_ := &SourceMapConsumerFactory{}
	InitSourceMapConsumerFactory(this_)
	return this_
}

func InitSourceMapConsumerFactory(this_ *SourceMapConsumerFactory) {
	this_.initFields_SourceMapConsumerFactory()
}

// Java static method: SourceMapConsumerFactory.parse(java.lang.String contents)
func SourceMapConsumerFactory_Parse_1(contents string) SourceMapping {
	_ = contents
	return SourceMapConsumerFactory_Parse_2(contents, nil)
}

// Java static method: SourceMapConsumerFactory.parse(java.lang.String contents, com.google.debugging.sourcemap.@org.jspecify.annotations.Nullable SourceMapSupplier supplier)
func SourceMapConsumerFactory_Parse_2(contents string, supplier SourceMapSupplier) SourceMapping {
	_ = contents
	_ = supplier
	if javalang.StringStartsWith(j2gort.CheckNotNullString(contents), "/** Begin line maps. **/") {
		panic(NewSourceMapParseException("This appears to be a V1 SourceMap, which is not supported."))
	} else if javalang.StringStartsWith(j2gort.CheckNotNullString(contents), "{") {
		var sourceMapObject *SourceMapObject = SourceMapObjectParser_Parse(contents)
		_ = sourceMapObject
		switch sourceMapObject.GetVersion() {
		case int32(3):
			{
				var val_1 SourceMapConsumerV3
				InitSourceMapConsumerV3(&val_1)
				var consumer *SourceMapConsumerV3 = &val_1
				_ = consumer
				consumer.Parse_2(sourceMapObject, supplier)
				return consumer
			}
		default:
			panic(NewSourceMapParseException(("Unknown source map version:" + j2gort.StrValInt32(sourceMapObject.GetVersion()))))
		}
	}
	panic(NewSourceMapParseException("unable to detect source map format"))
}

// Java static method: SourceMapConsumerFactory.parseFast(java.lang.String contents)
func SourceMapConsumerFactory_ParseFast(contents string) *SourceMapConsumerV3 {
	_ = contents
	if javalang.StringStartsWith(j2gort.CheckNotNullString(contents), "{") {
		var sourceMapObject *SourceMapObject = SourceMapObjectParser_ParseFast(contents)
		_ = sourceMapObject
		if !(sourceMapObject.GetVersion() == int32(3)) {
			j2gort.CheckArgumentFmt(false, "Unknown source map version:%s", j2gort.ValueOfInt(sourceMapObject.GetVersion()))
		}
		var val_2 SourceMapConsumerV3
		InitSourceMapConsumerV3(&val_2)
		var consumer *SourceMapConsumerV3 = &val_2
		_ = consumer
		consumer.Parse_2(sourceMapObject, nil)
		return consumer
	}
	panic(NewSourceMapParseException("unable to detect source map format"))
}
