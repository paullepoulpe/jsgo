/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/serialization/SerializeTypedAstPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.serialization.SerializeTypedAstPass
type SerializeTypedAstPass struct {
	compiler             AbstractCompiler                       `gson:"compiler"`
	consumer             j2gort.Ext_java_util_function_Consumer `gson:"consumer"`
	serializationOptions *SerializationOptions                  `gson:"serializationOptions"`
	idHash_              atomic.Int32
}

func (this_ *SerializeTypedAstPass) BindSelf_(s any) {
	_ = s
}
func (this_ *SerializeTypedAstPass) IsJavaNil_() bool { return this_ == nil }
func (*SerializeTypedAstPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.SerializeTypedAstPass", "SerializeTypedAstPass"
}
func (this_ *SerializeTypedAstPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SerializeTypedAstPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.SerializeTypedAstPass", "com.google.javascript.jscomp.serialization.SerializeTypedAstPass", func(v_ any) bool {
		return j2gort.InstanceOf[*SerializeTypedAstPass](v_)
	})
}
func (*SerializeTypedAstPass) Implements_jscomp_CompilerPass() {}

func (this_ *SerializeTypedAstPass) initFields_SerializeTypedAstPass() {
	_ = this_
}

// Java constructor: SerializeTypedAstPass(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.function.Consumer<com.google.javascript.jscomp.serialization.TypedAst> astConsumer, com.google.javascript.jscomp.serialization.SerializationOptions serializationOptions)
func NewSerializeTypedAstPass(compiler AbstractCompiler, astConsumer j2gort.Ext_java_util_function_Consumer, serializationOptions *SerializationOptions) *SerializeTypedAstPass {
	this_ := &SerializeTypedAstPass{}
	InitSerializeTypedAstPass(this_, compiler, astConsumer, serializationOptions)
	return this_
}

func InitSerializeTypedAstPass(this_ *SerializeTypedAstPass, compiler AbstractCompiler, astConsumer j2gort.Ext_java_util_function_Consumer, serializationOptions *SerializationOptions) {
	_ = compiler
	_ = astConsumer
	_ = serializationOptions
	this_.initFields_SerializeTypedAstPass()
	this_.compiler = compiler
	this_.consumer = astConsumer
	this_.serializationOptions = serializationOptions
}

// Java static method: SerializeTypedAstPass.createFromOutputStream(com.google.javascript.jscomp.AbstractCompiler c, java.io.OutputStream out, com.google.javascript.jscomp.serialization.SerializationOptions serializationOptions)
func SerializeTypedAstPass_CreateFromOutputStream(c AbstractCompiler, out j2gort.Ext_java_io_OutputStream, serializationOptions *SerializationOptions) *SerializeTypedAstPass {
	_ = c
	_ = out
	_ = serializationOptions
	var toOutputStream j2gort.Ext_java_util_function_Consumer = j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_1 any) {
		ast := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst](lp_1)
		_ = ast
		func() {
			defer func() {
				if rec_2 := recover(); rec_2 != nil {
					jex_3 := javalang.RecoverException(rec_2)
					_ = jex_3
					if j2gort.InstanceOf[javaio.IIOException](jex_3) {
						e := j2gort.AsIface[javaio.IIOException](jex_3)
						_ = e
						panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Cannot write to stream", e)))
					} else {
						panic(rec_2)
					}
				}
			}()
			j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst_List](j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst_List_NewBuilder_0().AddTypedAsts_1_Ext_com_google_javascript_jscomp_serialization_TypedAst(ast).Build()).WriteTo_1_Ext_java_io_OutputStream(out)
		}()
	})
	_ = toOutputStream
	return NewSerializeTypedAstPass(c, toOutputStream, serializationOptions)
}

// Java static method: SerializeTypedAstPass.createFromPath(com.google.javascript.jscomp.AbstractCompiler compiler, java.nio.file.Path outputPath, com.google.javascript.jscomp.serialization.SerializationOptions serializationOptions)
func SerializeTypedAstPass_CreateFromPath(compiler AbstractCompiler, outputPath j2gort.Ext_java_nio_file_Path, serializationOptions *SerializationOptions) *SerializeTypedAstPass {
	_ = compiler
	_ = outputPath
	_ = serializationOptions
	var toPath j2gort.Ext_java_util_function_Consumer = j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_4 any) {
		ast := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst](lp_4)
		_ = ast
		func() {
			defer func() {
				if rec_5 := recover(); rec_5 != nil {
					jex_6 := javalang.RecoverException(rec_5)
					_ = jex_6
					if j2gort.InstanceOf[javaio.IIOException](jex_6) {
						e := j2gort.AsIface[javaio.IIOException](jex_6)
						_ = e
						panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Cannot create TypedAst output file", e)))
					} else {
						panic(rec_5)
					}
				}
			}()
			var out j2gort.Ext_java_io_OutputStream = NewFastGzipOutputStream_1(j2gort.Ext_java_nio_file_Files_NewOutputStream(outputPath, []j2gort.Ext_java_nio_file_OpenOption{}))
			_ = out
			defer j2gort.CloseResource(out)
			j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst_List](j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst_List_NewBuilder_0().AddTypedAsts_1_Ext_com_google_javascript_jscomp_serialization_TypedAst(ast).Build()).WriteTo_1_Ext_java_io_OutputStream(out)
		}()
	})
	_ = toPath
	return NewSerializeTypedAstPass(compiler, toPath, serializationOptions)
}

// Java method: SerializeTypedAstPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *SerializeTypedAstPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		var val_7 RemoveCastNodes
		InitRemoveCastNodes(&val_7, this_.compiler)
		(&val_7).Process_2(externs, root)
	}
	clinit_TypedAstSerializer()
	var val_8 TypedAstSerializer
	InitTypedAstSerializer(&val_8, this_.compiler, this_.serializationOptions)
	var serializer *TypedAstSerializer = &val_8
	_ = serializer
	var ast *j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst = serializer.SerializeRoots(externs, root)
	_ = ast
	this_.consumer.Accept(j2gort.PtrToAny(ast))
}
