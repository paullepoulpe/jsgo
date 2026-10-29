// Port of test/com/google/javascript/jscomp/serialization/JSDocSerializerTest.java @ closure-compiler v20261021

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestJSDocSerializer_PreserveIsTypeSummary(t *testing.T) {
	builder := rhino.JSDocInfo_Builder_()
	builder.RecordTypeSummary()
	originalJsDocInfo := builder.Build()
	convertedJsDocInfo := JSDocSerializer_ConvertJSDocInfoForOptimizations(originalJsDocInfo)
	jsaTrue(t, convertedJsDocInfo.IsTypeSummary())
}

func TestJSDocSerializer_PreserveIsSassGeneratedCssTs(t *testing.T) {
	builder := rhino.JSDocInfo_Builder_()
	builder.RecordSassGeneratedCssTs()
	originalJsDocInfo := builder.Build()
	convertedJsDocInfo := JSDocSerializer_ConvertJSDocInfoForOptimizations(originalJsDocInfo)
	jsaTrue(t, convertedJsDocInfo.IsSassGeneratedCssTs())
}

func TestJSDocSerializer_PreserveSupportsClosureUnaware_simple(t *testing.T) {
	builder := rhino.JSDocInfo_Builder_()
	builder.RecordTypeSummary()
	jsaTrue(t, builder.RecordClosureUnawareCodeWithMode(rhino.JSDocInfo_PerFileClosureUnawareMode_SIMPLE))
	originalJsDocInfo := builder.Build()

	convertedJsDocInfo := JSDocSerializer_ConvertJSDocInfoForOptimizations(originalJsDocInfo)
	jsaEq(t, *convertedJsDocInfo.GetPerFileClosureUnawareMode(), rhino.JSDocInfo_PerFileClosureUnawareMode_SIMPLE)
}

// Go-only: extra coverage — the serialized form of a JSDoc with several
// fields, and its deserialization.
func TestJSDocSerializer_SerializedFields(t *testing.T) {
	builder := rhino.JSDocInfo_Builder_()
	builder.RecordConstancy()
	builder.RecordModifies(jutil.NewLinkedHashSet("this", "x"))
	builder.RecordDescription("desc")
	builder.RecordMeaning("meaning")
	builder.RecordSuppressions(jutil.NewLinkedHashSet("messageConventions", "partialAlias"))
	builder.RecordStableIdGenerator()
	pool := StringPool_Builder_()
	ser := JSDocSerializer_SerializeJsdoc(builder.Build(), pool)
	jsaEq(t, ser.String(), "kind: JSDOC_CONST\nkind: JSDOC_ID_GENERATOR_STABLE\n"+
		"kind: JSDOC_MODIFIES_THIS\nkind: JSDOC_MODIFIES_ARGUMENTS\nkind: JSDOC_SUPPRESS_PARTIAL_ALIAS\n"+
		"kind: JSDOC_SUPPRESS_MESSAGE_CONVENTION\ndescription_pointer: 1\nmeaning_pointer: 2\n")
	de := JSDocSerializer_DeserializeJsdoc(ser, pool.Build())
	jsaTrue(t, de.IsConstant())
	jsaTrue(t, de.IsStableIdGenerator())
	jsaDeepEq(t, de.GetModifies().ToSlice(), []string{"arguments", "this"})
	jsaDeepEq(t, de.GetSuppressions().ToSlice(), []string{"messageConventions", "partialAlias"})
	jsaEq(t, *de.GetDescription(), "desc")
	jsaTrue(t, JSDocSerializer_SerializeJsdoc(rhino.JSDocInfo_Builder_().Build(), pool) == nil)

	bad := &OptimizationJsdoc{kind_: []int32{500}}
	jsaThrows[IMalformedTypedAstException](t, func() { JSDocSerializer_DeserializeJsdoc(bad, pool.Build()) })
}
