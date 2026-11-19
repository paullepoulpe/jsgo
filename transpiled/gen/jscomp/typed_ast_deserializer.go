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
// Port of com/google/javascript/jscomp/serialization/TypedAstDeserializer.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_typed_ast_deserializer_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_typed_ast_deserializer_2 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_typed_ast_deserializer_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.BoxString(j2gort.CastPtr[SourceFile](mrp0_).GetName()) })
	lambdaSingleton_typed_ast_deserializer_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return TypedAstDeserializer_CreateStubWeakScriptNode(j2gort.CastPtr[SourceFile](mrp0_))
	})
}

// Java class: com.google.javascript.jscomp.serialization.TypedAstDeserializer
type TypedAstDeserializer struct {
	syntheticExterns              *SourceFile                                               `gson:"syntheticExterns"`
	colorPoolBuilder              j2gort.Ext_com_google_common_base_Optional                `gson:"colorPoolBuilder"`
	requiredInputFiles            j2gort.Ext_com_google_common_base_Optional                `gson:"requiredInputFiles"`
	filePoolBuilder               j2gort.OrderedMap[string, *SourceFile]                    `gson:"filePoolBuilder"`
	canonicalFilePool             j2gort.OrderedMap[string, *SourceFile]                    `gson:"canonicalFilePool"`
	typedAstFilesystem            j2gort.Ext_java_util_concurrent_ConcurrentMap             `gson:"typedAstFilesystem"`
	externProperties              j2gort.Ext_com_google_common_collect_ImmutableSet_Builder `gson:"externProperties"`
	runtimeLibraries              j2gort.Ext_com_google_common_collect_ImmutableSet_Builder `gson:"runtimeLibraries"`
	syntheticExternsDeserializers j2gort.Ext_java_util_ArrayList                            `gson:"syntheticExternsDeserializers"`
	idHash_                       atomic.Int32
	mode                          TypedAstDeserializer_Mode `gson:"mode"`
}

func (this_ *TypedAstDeserializer) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedAstDeserializer) IsJavaNil_() bool { return this_ == nil }
func (*TypedAstDeserializer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.TypedAstDeserializer", "TypedAstDeserializer"
}
func (this_ *TypedAstDeserializer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedAstDeserializer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.TypedAstDeserializer", "com.google.javascript.jscomp.serialization.TypedAstDeserializer", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedAstDeserializer](v_)
	})
}

func (this_ *TypedAstDeserializer) initFields_TypedAstDeserializer() {
	_ = this_
	this_.filePoolBuilder = j2gort.NewOrderedMap[string, *SourceFile]()
	this_.canonicalFilePool = j2gort.NewOrderedMap[string, *SourceFile]()
	this_.typedAstFilesystem = j2gort.NewExt_java_util_concurrent_ConcurrentHashMap()
	this_.externProperties = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	this_.runtimeLibraries = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	this_.syntheticExternsDeserializers = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: TypedAstDeserializer(com.google.javascript.jscomp.SourceFile syntheticExterns, com.google.common.base.Optional<com.google.javascript.jscomp.serialization.ColorPool.Builder> existingColorPool, com.google.common.base.Optional<com.google.common.collect.ImmutableSet<com.google.javascript.jscomp.SourceFile>> requiredInputFiles, com.google.javascript.jscomp.serialization.TypedAstDeserializer.Mode mode, boolean includeTypeInformation)
func NewTypedAstDeserializer(syntheticExterns *SourceFile, existingColorPool j2gort.Ext_com_google_common_base_Optional, requiredInputFiles j2gort.Ext_com_google_common_base_Optional, mode TypedAstDeserializer_Mode, includeTypeInformation bool) *TypedAstDeserializer {
	this_ := &TypedAstDeserializer{}
	InitTypedAstDeserializer(this_, syntheticExterns, existingColorPool, requiredInputFiles, mode, includeTypeInformation)
	return this_
}

func InitTypedAstDeserializer(this_ *TypedAstDeserializer, syntheticExterns *SourceFile, existingColorPool j2gort.Ext_com_google_common_base_Optional, requiredInputFiles j2gort.Ext_com_google_common_base_Optional, mode TypedAstDeserializer_Mode, includeTypeInformation bool) {
	_ = syntheticExterns
	_ = existingColorPool
	_ = requiredInputFiles
	_ = mode
	_ = includeTypeInformation
	this_.initFields_TypedAstDeserializer()
	this_.syntheticExterns = syntheticExterns
	this_.mode = mode
	if includeTypeInformation {
		this_.colorPoolBuilder = j2gort.Ext_com_google_common_base_Optional_Of(j2gort.PtrToAny(j2gort.CastPtr[ColorPool_Builder](existingColorPool.Or_1_any(j2gort.PtrToAny(ColorPool_S_Builder())))))
	} else {
		this_.colorPoolBuilder = j2gort.Ext_com_google_common_base_Optional_Absent()
	}
	this_.requiredInputFiles = requiredInputFiles
}

// Java static method: TypedAstDeserializer.deserializeFullAst(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.SourceFile syntheticExterns, com.google.common.collect.ImmutableSet<com.google.javascript.jscomp.SourceFile> requiredInputFiles, java.io.InputStream typedAstsStream, boolean includeTypeInformation, boolean resolveSourceMapAnnotations, boolean parseInlineSourceMaps)
func TypedAstDeserializer_DeserializeFullAst(compiler AbstractCompiler, syntheticExterns *SourceFile, requiredInputFiles j2gort.Ext_com_google_common_collect_ImmutableSet, typedAstsStream j2gort.Ext_java_io_InputStream, includeTypeInformation bool, resolveSourceMapAnnotations bool, parseInlineSourceMaps bool) TypedAstDeserializer_DeserializedAst {
	_ = compiler
	_ = syntheticExterns
	_ = requiredInputFiles
	_ = typedAstsStream
	_ = includeTypeInformation
	_ = resolveSourceMapAnnotations
	_ = parseInlineSourceMaps
	var sourceFilesByName j2gort.Ext_com_google_common_collect_ImmutableMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](requiredInputFiles.Stream().Collect_1(j2gort.Ext_com_google_common_collect_ImmutableMap_ToImmutableMap_2(lambdaSingleton_typed_ast_deserializer_1, j2gort.Ext_java_util_function_Function_Identity())))
	_ = sourceFilesByName
	return TypedAstDeserializer_Deserialize(compiler, syntheticExterns, j2gort.Ext_com_google_common_base_Optional_Of(requiredInputFiles), sourceFilesByName, j2gort.Ext_com_google_common_base_Optional_Absent(), typedAstsStream, StaticGet_TypedAstDeserializer_Mode_FULL_AST(), includeTypeInformation, resolveSourceMapAnnotations, parseInlineSourceMaps)
}

// Java static method: TypedAstDeserializer.deserializeRuntimeLibraries(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.SourceFile syntheticExterns, com.google.common.base.Optional<com.google.javascript.jscomp.serialization.ColorPool.Builder> colorPool, java.io.InputStream typedAstsStream, boolean resolveSourceMapAnnotations, boolean parseInlineSourceMaps)
func TypedAstDeserializer_DeserializeRuntimeLibraries(compiler AbstractCompiler, syntheticExterns *SourceFile, colorPool j2gort.Ext_com_google_common_base_Optional, typedAstsStream j2gort.Ext_java_io_InputStream, resolveSourceMapAnnotations bool, parseInlineSourceMaps bool) TypedAstDeserializer_DeserializedAst {
	_ = compiler
	_ = syntheticExterns
	_ = colorPool
	_ = typedAstsStream
	_ = resolveSourceMapAnnotations
	_ = parseInlineSourceMaps
	return TypedAstDeserializer_Deserialize(compiler, syntheticExterns, j2gort.Ext_com_google_common_base_Optional_Absent(), j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0(), colorPool, typedAstsStream, StaticGet_TypedAstDeserializer_Mode_RUNTIME_LIBRARY_ONLY(), colorPool.IsPresent(), resolveSourceMapAnnotations, parseInlineSourceMaps)
}

// Java static method: TypedAstDeserializer.deserialize(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.SourceFile syntheticExterns, com.google.common.base.Optional<com.google.common.collect.ImmutableSet<com.google.javascript.jscomp.SourceFile>> requiredInputFiles, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.SourceFile> scriptSourceFiles, com.google.common.base.Optional<com.google.javascript.jscomp.serialization.ColorPool.Builder> colorPool, java.io.InputStream typedAstStream, com.google.javascript.jscomp.serialization.TypedAstDeserializer.Mode mode, boolean includeTypeInformation, boolean resolveSourceMapAnnotations, boolean parseInlineSourceMaps)
func TypedAstDeserializer_Deserialize(compiler AbstractCompiler, syntheticExterns *SourceFile, requiredInputFiles j2gort.Ext_com_google_common_base_Optional, scriptSourceFiles j2gort.Ext_com_google_common_collect_ImmutableMap, colorPool j2gort.Ext_com_google_common_base_Optional, typedAstStream j2gort.Ext_java_io_InputStream, mode TypedAstDeserializer_Mode, includeTypeInformation bool, resolveSourceMapAnnotations bool, parseInlineSourceMaps bool) TypedAstDeserializer_DeserializedAst {
	_ = compiler
	_ = syntheticExterns
	_ = requiredInputFiles
	_ = scriptSourceFiles
	_ = colorPool
	_ = typedAstStream
	_ = mode
	_ = includeTypeInformation
	_ = resolveSourceMapAnnotations
	_ = parseInlineSourceMaps
	if !(colorPool.IsPresent() == (func() bool {
		re_1 := mode
		ae_2 := StaticGet_TypedAstDeserializer_Mode_RUNTIME_LIBRARY_ONLY()
		if re_1 == 0 {
			j2gort.ThrowNPE()
		}
		return re_1 == ae_2
	}() && includeTypeInformation)) {
		j2gort.CheckArgumentStr(false, "ColorPool.Builder required iff deserializing runtime libraries & including types")
	}
	var val_3 TypedAstDeserializer
	InitTypedAstDeserializer(&val_3, syntheticExterns, colorPool, requiredInputFiles, mode, includeTypeInformation)
	var deserializer *TypedAstDeserializer = &val_3
	_ = deserializer
	deserializer.priv_TypedAstDeserializer_addToFilePool_7774(syntheticExterns.GetName(), syntheticExterns)
	scriptSourceFiles.ForEach_1(func() j2gort.Ext_java_util_function_BiConsumer {
		bound_4 := deserializer
		return j2gort.NewLambda_Ext_java_util_function_BiConsumer(func(mrp0_ any, mrp1_ any) {
			bound_4.priv_TypedAstDeserializer_addToFilePool_7774(j2gort.CastString(mrp0_), j2gort.CastPtr[SourceFile](mrp1_))
		})
	}())
	{
		var iifeCond_1 bool
		{

			re_5 := mode
			ae_6 := StaticGet_TypedAstDeserializer_Mode_RUNTIME_LIBRARY_ONLY()
			if re_5 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_5 == ae_6
		}
		if !(iifeCond_1) {
			compiler.InitRuntimeLibraryTypedAsts(deserializer.colorPoolBuilder)
		}
	}
	TypedAstDeserializer_DeserializeTypedAsts(typedAstStream, deserializer, compiler, resolveSourceMapAnnotations, parseInlineSourceMaps)
	deserializer.typedAstFilesystem.Put_2(j2gort.PtrToAny(syntheticExterns), func() j2gort.Ext_java_util_function_Supplier {
		bound_7 := NewTypedAstDeserializer_SyntheticExternsDeserializer(syntheticExterns, deserializer.syntheticExternsDeserializers)
		return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.BoxNode(bound_7.Deserialize()) })
	}())
	return deserializer.priv_TypedAstDeserializer_toDeserializedAst_7772()
}

// Java method: TypedAstDeserializer.toDeserializedAst()
func (this_ *TypedAstDeserializer) priv_TypedAstDeserializer_toDeserializedAst_7772() TypedAstDeserializer_DeserializedAst {
	var registry j2gort.Ext_com_google_common_base_Optional
	_ = registry
	{
		var iifeCond_1 bool
		{

			re_8 := this_.mode
			ae_9 := StaticGet_TypedAstDeserializer_Mode_RUNTIME_LIBRARY_ONLY()
			if re_8 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_8 == ae_9
		}
		if iifeCond_1 || !(this_.colorPoolBuilder.IsPresent()) {
			registry = j2gort.Ext_com_google_common_base_Optional_Absent()
		} else {
			registry = j2gort.Ext_com_google_common_base_Optional_Of(j2gort.PtrToAny(j2gort.CastPtr[ColorPool_Builder](this_.colorPoolBuilder.Get()).Build().GetRegistry()))
		}
	}
	return TypedAstDeserializer_DeserializedAst_Create(this_.typedAstFilesystem, registry, j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.externProperties.Build()), j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.runtimeLibraries.Build()))
}

// Java method: TypedAstDeserializer.deserializeTypedAst(com.google.javascript.jscomp.serialization.TypedAst typedAstProto, com.google.javascript.jscomp.AbstractCompiler compiler, boolean resolveSourceMapAnnotations, boolean parseInlineSourceMaps)
func (this_ *TypedAstDeserializer) priv_TypedAstDeserializer_deserializeTypedAst_7773(typedAstProto *j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst, compiler AbstractCompiler, resolveSourceMapAnnotations bool, parseInlineSourceMaps bool) {
	_ = typedAstProto
	_ = compiler
	_ = resolveSourceMapAnnotations
	_ = parseInlineSourceMaps
	var fileShard []*SourceFile = j2gort.ImmListToSliceNullable[*SourceFile](this_.priv_TypedAstDeserializer_toFileShard_7776(typedAstProto))
	_ = fileShard
	{
		var iifeCond_1 bool
		{

			re_10 := this_.mode
			ae_11 := StaticGet_TypedAstDeserializer_Mode_FULL_AST()
			if re_10 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_10 == ae_11
		}
		if iifeCond_1 {
			var containsRequiredInput bool = false
			_ = containsRequiredInput
			{
				cur_12 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_Iterables_Concat_2(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](typedAstProto.GetExternAstList()), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](typedAstProto.GetCodeAstList()))))
				for cur_12.HasNext() {
					rawElem_13 := cur_12.Next()
					lazyAst := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_LazyAst](rawElem_13)
					_ = lazyAst
					var file *SourceFile = fileShard[(lazyAst.GetSourceFile() - int32(1))]
					_ = file
					if j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.requiredInputFiles.Get()).Contains(j2gort.PtrToAny(file)) {
						containsRequiredInput = true
						break
					}
				}
			}
			if !(containsRequiredInput) {
				return
			}
		}
	}
	var runtimeLibraryList j2gort.Ext_com_google_protobuf_ProtocolStringList = j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](typedAstProto.GetRuntimeLibraryToInjectList())
	_ = runtimeLibraryList
	{
		var i int32 = int32(0)
		_ = i
		for ; i < runtimeLibraryList.Size(); i = i + 1 {
			j2gort.ImmSetBuilderAddStrVoid(this_.runtimeLibraries, j2gort.ListGetStr(runtimeLibraryList, i))
		}
	}
	var stringShard *StringPool = StringPool_FromProto(typedAstProto.GetStringPool())
	_ = stringShard
	var colorShard j2gort.Ext_com_google_common_base_Optional = this_.colorPoolBuilder.Transform(j2gort.NewLambda_Ext_com_google_common_base_Function(func(lp_14 any) any {
		builder := j2gort.CastPtr[ColorPool_Builder](lp_14)
		_ = builder
		return j2gort.PtrToAny(builder.AddShard(typedAstProto.GetTypePool(), stringShard))
	}))
	_ = colorShard
	{
		cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(typedAstProto.GetExternsSummary().GetPropNamePtrList()))
		for cur_15.HasNext() {
			rawElem_16 := cur_15.Next()
			x := javalang.D2I(j2gort.UnboxToDouble(rawElem_16))
			_ = x
			j2gort.ImmSetBuilderAddStrVoid(this_.externProperties, stringShard.Get(x))
		}
	}
	{
		cur_17 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_Iterables_Concat_2(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](typedAstProto.GetExternAstList()), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](typedAstProto.GetCodeAstList()))))
		for cur_17.HasNext() {
			rawElem_18 := cur_17.Next()
			lazyAst_2 := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_LazyAst](rawElem_18)
			_ = lazyAst_2
			this_.priv_TypedAstDeserializer_initLazyAstDeserializer_7777(lazyAst_2, j2gort.ImmListFromTypedSlice(fileShard), colorShard, stringShard, compiler, resolveSourceMapAnnotations, parseInlineSourceMaps)
		}
	}
}

// Java method: TypedAstDeserializer.addToFilePool(java.lang.String filename, com.google.javascript.jscomp.SourceFile file)
func (this_ *TypedAstDeserializer) priv_TypedAstDeserializer_addToFilePool_7774(filename string, file *SourceFile) {
	_ = filename
	_ = file
	this_.filePoolBuilder.Put(filename, file)
	var canonicalPath string = SourceFile_CanonicalizePathForMatching(filename)
	_ = canonicalPath
	if !(j2gort.CheckNotNullString(canonicalPath) == filename) {
		this_.canonicalFilePool.PutIfAbsent(canonicalPath, file)
	}
}

// Java method: TypedAstDeserializer.findExistingFile(java.lang.String filename)
func (this_ *TypedAstDeserializer) priv_TypedAstDeserializer_findExistingFile_7775(filename string) *SourceFile {
	_ = filename
	var exact *SourceFile = this_.filePoolBuilder.Get(filename)
	_ = exact
	if exact != nil {
		return exact
	}
	var canonicalPath string = SourceFile_CanonicalizePathForMatching(filename)
	_ = canonicalPath
	var canonical *SourceFile = this_.canonicalFilePool.Get(canonicalPath)
	_ = canonical
	if canonical != nil {
		return canonical
	}
	return this_.filePoolBuilder.Get(canonicalPath)
}

// Java method: TypedAstDeserializer.toFileShard(com.google.javascript.jscomp.serialization.TypedAst typedAstProto)
func (this_ *TypedAstDeserializer) priv_TypedAstDeserializer_toFileShard_7776(typedAstProto *j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = typedAstProto
	var fileShardBuilder []*SourceFile = []*SourceFile{}
	_ = fileShardBuilder
	var protos j2gort.Ext_java_util_List = typedAstProto.GetSourceFilePool().GetSourceFileList()
	_ = protos
	{
		var i int32 = int32(0)
		_ = i
		for ; i < protos.Size(); i = i + 1 {
			var p *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto](protos.Get(i))
			_ = p
			var existingFile *SourceFile = this_.priv_TypedAstDeserializer_findExistingFile_7775(p.GetFilename())
			_ = existingFile
			if existingFile != nil {
				existingFile.RestoreCachedStateFrom(p)
				fileShardBuilder = append(fileShardBuilder, j2gort.CheckNotNullElem(existingFile))
			} else {
				var protoFile *SourceFile = SourceFile_FromProto_1(p)
				_ = protoFile
				this_.priv_TypedAstDeserializer_addToFilePool_7774(p.GetFilename(), protoFile)
				fileShardBuilder = append(fileShardBuilder, j2gort.CheckNotNullElem(protoFile))
			}
		}
	}
	return j2gort.ImmListFromTypedSlice(fileShardBuilder)
}

// Java method: TypedAstDeserializer.initLazyAstDeserializer(com.google.javascript.jscomp.serialization.LazyAst lazyAst, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.SourceFile> fileShard, com.google.common.base.Optional<com.google.javascript.jscomp.serialization.ColorPool.ShardView> colorShard, com.google.javascript.jscomp.serialization.StringPool stringShard, com.google.javascript.jscomp.AbstractCompiler compiler, boolean resolveSourceMapAnnotations, boolean parseInlineSourceMaps)
func (this_ *TypedAstDeserializer) priv_TypedAstDeserializer_initLazyAstDeserializer_7777(lazyAst *j2gort.Ext_com_google_javascript_jscomp_serialization_LazyAst, fileShard j2gort.Ext_com_google_common_collect_ImmutableList, colorShard j2gort.Ext_com_google_common_base_Optional, stringShard *StringPool, compiler AbstractCompiler, resolveSourceMapAnnotations bool, parseInlineSourceMaps bool) {
	_ = lazyAst
	_ = fileShard
	_ = colorShard
	_ = stringShard
	_ = compiler
	_ = resolveSourceMapAnnotations
	_ = parseInlineSourceMaps
	var file *SourceFile = j2gort.CastPtr[SourceFile](fileShard.Get((lazyAst.GetSourceFile() - int32(1))))
	_ = file
	if (this_.requiredInputFiles.IsPresent() && !(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.requiredInputFiles.Get()).Contains(j2gort.PtrToAny(file)))) && !(this_.syntheticExterns == file) {
		return
	}
	if file.IsWeak() {
		this_.typedAstFilesystem.ComputeIfAbsent(j2gort.PtrToAny(file), lambdaSingleton_typed_ast_deserializer_2)
		return
	}
	var val_19 ScriptNodeDeserializer
	InitScriptNodeDeserializer(&val_19, lazyAst, stringShard, colorShard, fileShard)
	var deserializer *ScriptNodeDeserializer = &val_19
	_ = deserializer
	if this_.syntheticExterns == file {
		this_.syntheticExternsDeserializers.Add_1(deserializer)
	} else {
		func() {
			m_20 := this_.typedAstFilesystem
			k_21 := j2gort.PtrToAny(file)
			if _, ok_23 := j2gort.MapTryGetComputeAny(m_20, k_21); !ok_23 {
				m_20.ComputeIfAbsent(k_21, j2gort.NewLambda_Ext_java_util_function_Function(func(lp_24 any) any {
					f := j2gort.CastPtr[SourceFile](lp_24)
					_ = f
					return func() j2gort.Ext_java_util_function_Supplier {
						bound_25 := deserializer
						return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.BoxNode(bound_25.DeserializeNew()) })
					}()
				}))
			}
		}()
	}
	TypedAstDeserializer_AddInputSourceMap(deserializer, compiler, resolveSourceMapAnnotations, parseInlineSourceMaps)
}

// Java static method: TypedAstDeserializer.addInputSourceMap(com.google.javascript.jscomp.serialization.ScriptNodeDeserializer deserializer, com.google.javascript.jscomp.AbstractCompiler compiler, boolean resolveSourceMapAnnotations, boolean parseInlineSourceMaps)
func TypedAstDeserializer_AddInputSourceMap(deserializer *ScriptNodeDeserializer, compiler AbstractCompiler, resolveSourceMapAnnotations bool, parseInlineSourceMaps bool) {
	_ = deserializer
	_ = compiler
	_ = resolveSourceMapAnnotations
	_ = parseInlineSourceMaps
	var sourceFile *SourceFile = deserializer.GetSourceFile()
	_ = sourceFile
	var sourceMappingURL string = deserializer.GetSourceMappingURL()
	_ = sourceMappingURL
	if ((sourceMappingURL != j2gort.NullString) && (j2gort.JStrLen(sourceMappingURL) > int32(0))) && resolveSourceMapAnnotations {
		var base64EncodedSourceMap string = SourceMapResolver_AddBase64PrefixToEncodedSourceMap(sourceMappingURL)
		_ = base64EncodedSourceMap
		var sourceMapSourceFile *SourceFile = SourceMapResolver_ExtractSourceMap(sourceFile, base64EncodedSourceMap, parseInlineSourceMaps)
		_ = sourceMapSourceFile
		if sourceMapSourceFile != nil {
			compiler.AddInputSourceMap(sourceFile.GetName(), NewSourceMapInput(sourceMapSourceFile))
		}
	}
}

// Java static method: TypedAstDeserializer.createStubWeakScriptNode(com.google.javascript.jscomp.SourceFile file)
func TypedAstDeserializer_CreateStubWeakScriptNode(file *SourceFile) j2gort.Ext_java_util_function_Supplier {
	_ = file
	return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		var stubScript pkg_rhino.Node = pkg_rhino.IR_Script_0().SetStaticSourceFile(pkg_rhino.ToIface_StaticSourceFile(file))
		_ = stubScript
		stubScript.SetFeatureSet(j2gort.PtrToAny(pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM()))
		return j2gort.BoxNode(stubScript)
	})
}

// Java static method: TypedAstDeserializer.deserializeTypedAsts(java.io.InputStream typedAstsStream, com.google.javascript.jscomp.serialization.TypedAstDeserializer deserializer, com.google.javascript.jscomp.AbstractCompiler compiler, boolean resolveSourceMapAnnotations, boolean parseInlineSourceMaps)
func TypedAstDeserializer_DeserializeTypedAsts(typedAstsStream j2gort.Ext_java_io_InputStream, deserializer *TypedAstDeserializer, compiler AbstractCompiler, resolveSourceMapAnnotations bool, parseInlineSourceMaps bool) {
	_ = typedAstsStream
	_ = deserializer
	_ = compiler
	_ = resolveSourceMapAnnotations
	_ = parseInlineSourceMaps
	func() {
		defer func() {
			if rec_26 := recover(); rec_26 != nil {
				jex_27 := javalang.RecoverException(rec_26)
				_ = jex_27
				if j2gort.InstanceOf[javaio.IIOException](jex_27) {
					ex := j2gort.AsIface[javaio.IIOException](jex_27)
					_ = ex
					panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Cannot read from TypedAST input stream", ex)))
				} else {
					panic(rec_26)
				}
			}
		}()
		var codedInput j2gort.Ext_com_google_protobuf_CodedInputStream = j2gort.Ext_com_google_protobuf_CodedInputStream_NewInstance_1_Ext_java_io_InputStream(typedAstsStream)
		_ = codedInput
		var typedAstBuilder *j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst_Builder = j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst_NewBuilder_0()
		_ = typedAstBuilder
		for !(codedInput.IsAtEnd()) {
			var tag int32 = codedInput.ReadTag()
			_ = tag
			if (j2gort.Ext_com_google_protobuf_WireFormat_GetTagFieldNumber(tag) != int32(1)) || (j2gort.Ext_com_google_protobuf_WireFormat_GetTagWireType(tag) != int32(2)) {
				panic(j2gort.NewExt_com_google_protobuf_InvalidProtocolBufferException_1_string(("Unexpected field number " + j2gort.StrValInt32(j2gort.Ext_com_google_protobuf_WireFormat_GetTagFieldNumber(tag)) + " or wire type " + j2gort.StrValInt32(j2gort.Ext_com_google_protobuf_WireFormat_GetTagWireType(tag)))))
			}
			codedInput.ReadMessage_2_Ext_com_google_protobuf_MessageLite_Builder_Ext_com_google_protobuf_ExtensionRegistryLite(j2gort.AsIface[j2gort.Ext_com_google_protobuf_MessageLite_Builder](typedAstBuilder), j2gort.AsIface[j2gort.Ext_com_google_protobuf_ExtensionRegistryLite](j2gort.Ext_com_google_protobuf_ExtensionRegistry_GetEmptyRegistry()))
			var typedAst *j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst](typedAstBuilder.Build())
			_ = typedAst
			typedAstBuilder.Clear()
			codedInput.ResetSizeCounter()
			deserializer.priv_TypedAstDeserializer_deserializeTypedAst_7773(typedAst, compiler, resolveSourceMapAnnotations, parseInlineSourceMaps)
		}
	}()
}

// Java enum: com.google.javascript.jscomp.serialization.TypedAstDeserializer.Mode
type TypedAstDeserializer_Mode uint8

const (
	enum_TypedAstDeserializer_Mode_RUNTIME_LIBRARY_ONLY TypedAstDeserializer_Mode = 1
	enum_TypedAstDeserializer_Mode_FULL_AST             TypedAstDeserializer_Mode = 2
)

var enumBaseTable_TypedAstDeserializer_Mode = [...]j2gort.EnumBase{
	{Name_: "RUNTIME_LIBRARY_ONLY", Ordinal_: 0},
	{Name_: "FULL_AST", Ordinal_: 1},
}
var enumValues_TypedAstDeserializer_Mode = []TypedAstDeserializer_Mode{1, 2}

func StaticGet_TypedAstDeserializer_Mode_RUNTIME_LIBRARY_ONLY() TypedAstDeserializer_Mode {
	return enum_TypedAstDeserializer_Mode_RUNTIME_LIBRARY_ONLY
}
func StaticGet_TypedAstDeserializer_Mode_FULL_AST() TypedAstDeserializer_Mode {
	return enum_TypedAstDeserializer_Mode_FULL_AST
}

func TypedAstDeserializer_Mode_Values() []TypedAstDeserializer_Mode {
	return j2gort.CloneSlice(enumValues_TypedAstDeserializer_Mode)
}

func TypedAstDeserializer_Mode_ValueOf(name string) TypedAstDeserializer_Mode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_TypedAstDeserializer_Mode {
		if enumBaseTable_TypedAstDeserializer_Mode[i_].Name_ == name {
			return TypedAstDeserializer_Mode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.serialization.TypedAstDeserializer.Mode." + name))
}

func (this_ TypedAstDeserializer_Mode) IsJavaNil_() bool                              { return this_ == 0 }
func (TypedAstDeserializer_Mode) HasIdentityEquals_()                                 {}
func (TypedAstDeserializer_Mode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (TypedAstDeserializer_Mode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ TypedAstDeserializer_Mode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (TypedAstDeserializer_Mode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.TypedAstDeserializer$Mode", "Mode"
}
func (this_ TypedAstDeserializer_Mode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_TypedAstDeserializer_Mode[this_-1]
}
func (this_ TypedAstDeserializer_Mode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ TypedAstDeserializer_Mode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ TypedAstDeserializer_Mode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypedAstDeserializer_Mode[this_-1].Name_
}
func (this_ TypedAstDeserializer_Mode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_TypedAstDeserializer_Mode[this_-1].Name_
}
func (this_ TypedAstDeserializer_Mode) Equals(other_ any) bool {
	o_, ok_ := other_.(TypedAstDeserializer_Mode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ TypedAstDeserializer_Mode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_TypedAstDeserializer_Mode[this_-1].HashCode()
}
func (this_ TypedAstDeserializer_Mode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypedAstDeserializer_Mode[this_-1].CompareTo(other_)
}
func (this_ TypedAstDeserializer_Mode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_TypedAstDeserializer_Mode[this_-1].Name_)
}
func init() {
	vals_ := []any{TypedAstDeserializer_Mode(1), TypedAstDeserializer_Mode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.serialization.TypedAstDeserializer.Mode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.serialization.TypedAstDeserializer$Mode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.TypedAstDeserializer.Mode", "com.google.javascript.jscomp.serialization.TypedAstDeserializer$Mode", func(v_ any) bool {
		return j2gort.InstanceOf[TypedAstDeserializer_Mode](v_)
	})
}

// Java class: com.google.javascript.jscomp.serialization.TypedAstDeserializer.SyntheticExternsDeserializer
type TypedAstDeserializer_SyntheticExternsDeserializer struct {
	syntheticExterns              *SourceFile                    `gson:"syntheticExterns"`
	syntheticExternsDeserializers j2gort.Ext_java_util_ArrayList `gson:"syntheticExternsDeserializers"`
	idHash_                       atomic.Int32
}

func (this_ *TypedAstDeserializer_SyntheticExternsDeserializer) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedAstDeserializer_SyntheticExternsDeserializer) IsJavaNil_() bool {
	return this_ == nil
}
func (*TypedAstDeserializer_SyntheticExternsDeserializer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.TypedAstDeserializer$SyntheticExternsDeserializer", "SyntheticExternsDeserializer"
}
func (this_ *TypedAstDeserializer_SyntheticExternsDeserializer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedAstDeserializer_SyntheticExternsDeserializer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.TypedAstDeserializer.SyntheticExternsDeserializer", "com.google.javascript.jscomp.serialization.TypedAstDeserializer$SyntheticExternsDeserializer", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedAstDeserializer_SyntheticExternsDeserializer](v_)
	})
}

func (this_ *TypedAstDeserializer_SyntheticExternsDeserializer) initFields_TypedAstDeserializer_SyntheticExternsDeserializer() {
	_ = this_
}

// Java constructor: SyntheticExternsDeserializer(com.google.javascript.jscomp.SourceFile syntheticExterns, java.util.ArrayList<com.google.javascript.jscomp.serialization.ScriptNodeDeserializer> syntheticExternsDeserializers)
func NewTypedAstDeserializer_SyntheticExternsDeserializer(syntheticExterns *SourceFile, syntheticExternsDeserializers j2gort.Ext_java_util_ArrayList) *TypedAstDeserializer_SyntheticExternsDeserializer {
	this_ := &TypedAstDeserializer_SyntheticExternsDeserializer{}
	InitTypedAstDeserializer_SyntheticExternsDeserializer(this_, syntheticExterns, syntheticExternsDeserializers)
	return this_
}

func InitTypedAstDeserializer_SyntheticExternsDeserializer(this_ *TypedAstDeserializer_SyntheticExternsDeserializer, syntheticExterns *SourceFile, syntheticExternsDeserializers j2gort.Ext_java_util_ArrayList) {
	_ = syntheticExterns
	_ = syntheticExternsDeserializers
	this_.initFields_TypedAstDeserializer_SyntheticExternsDeserializer()
	this_.syntheticExterns = syntheticExterns
	this_.syntheticExternsDeserializers = syntheticExternsDeserializers
}

// Java method: SyntheticExternsDeserializer.deserialize()
func (this_ *TypedAstDeserializer_SyntheticExternsDeserializer) Deserialize() pkg_rhino.Node {
	var script pkg_rhino.Node = pkg_rhino.IR_Script_0()
	_ = script
	script.SetStaticSourceFile(pkg_rhino.ToIface_StaticSourceFile(this_.syntheticExterns))
	{
		cur_28 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.syntheticExternsDeserializers))
		for cur_28.HasNext() {
			rawElem_29 := cur_28.Next()
			d := j2gort.CastPtr[ScriptNodeDeserializer](rawElem_29)
			_ = d
			script.AddChildrenToBack(d.DeserializeNew().RemoveChildren())
		}
	}
	return script
}

// Java polymorphic class interface: com.google.javascript.jscomp.serialization.TypedAstDeserializer.DeserializedAst
type TypedAstDeserializer_DeserializedAst interface {
	Implements_jscomp_TypedAstDeserializer_DeserializedAst()
	IsJavaNil_() bool
	GetTypedAstDeserializer_DeserializedAstBase() *TypedAstDeserializer_DeserializedAstBase
	GetFilesystem() j2gort.Ext_java_util_concurrent_ConcurrentMap
	GetColorRegistry() j2gort.Ext_com_google_common_base_Optional
	GetExternProperties() j2gort.Ext_com_google_common_collect_ImmutableSet
	GetRuntimeLibraries() j2gort.Ext_com_google_common_collect_ImmutableSet
}

func ToIface_TypedAstDeserializer_DeserializedAst[T interface {
	TypedAstDeserializer_DeserializedAst
	comparable
}](v T) TypedAstDeserializer_DeserializedAst {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_TypedAstDeserializer_DeserializedAst[T TypedAstDeserializer_DeserializedAst](v T) TypedAstDeserializer_DeserializedAst {
	var vi TypedAstDeserializer_DeserializedAst = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.serialization.TypedAstDeserializer.DeserializedAst
type TypedAstDeserializer_DeserializedAstBase struct {
	idHash_ atomic.Int32
}

func (this_ *TypedAstDeserializer_DeserializedAstBase) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedAstDeserializer_DeserializedAstBase) IsJavaNil_() bool { return this_ == nil }
func (*TypedAstDeserializer_DeserializedAstBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.TypedAstDeserializer$DeserializedAst", "DeserializedAst"
}
func (this_ *TypedAstDeserializer_DeserializedAstBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.TypedAstDeserializer.DeserializedAst", "com.google.javascript.jscomp.serialization.TypedAstDeserializer$DeserializedAst", func(v_ any) bool {
		return j2gort.InstanceOf[TypedAstDeserializer_DeserializedAst](v_)
	})
}
func (b *TypedAstDeserializer_DeserializedAstBase) GetTypedAstDeserializer_DeserializedAstBase() *TypedAstDeserializer_DeserializedAstBase {
	return b
}
func (*TypedAstDeserializer_DeserializedAstBase) Implements_jscomp_TypedAstDeserializer_DeserializedAst() {
}

func (this_ *TypedAstDeserializer_DeserializedAstBase) initFields_TypedAstDeserializer_DeserializedAst() {
	_ = this_
}

func InitTypedAstDeserializer_DeserializedAst(this_ *TypedAstDeserializer_DeserializedAstBase) {
	this_.initFields_TypedAstDeserializer_DeserializedAst()
}

// Java method: DeserializedAst.getFilesystem()
func (this_ *TypedAstDeserializer_DeserializedAstBase) GetFilesystem() j2gort.Ext_java_util_concurrent_ConcurrentMap {
	panic("abstract or native method invoked")
}

// Java method: DeserializedAst.getColorRegistry()
func (this_ *TypedAstDeserializer_DeserializedAstBase) GetColorRegistry() j2gort.Ext_com_google_common_base_Optional {
	panic("abstract or native method invoked")
}

// Java method: DeserializedAst.getExternProperties()
func (this_ *TypedAstDeserializer_DeserializedAstBase) GetExternProperties() j2gort.Ext_com_google_common_collect_ImmutableSet {
	panic("abstract or native method invoked")
}

// Java method: DeserializedAst.getRuntimeLibraries()
func (this_ *TypedAstDeserializer_DeserializedAstBase) GetRuntimeLibraries() j2gort.Ext_com_google_common_collect_ImmutableSet {
	panic("abstract or native method invoked")
}

// Java static method: DeserializedAst.create(java.util.concurrent.ConcurrentMap<com.google.javascript.jscomp.SourceFile,java.util.function.Supplier<com.google.javascript.rhino.Node>> filesystem, com.google.common.base.Optional<com.google.javascript.jscomp.colors.ColorRegistry> colorRegistry, com.google.common.collect.ImmutableSet<java.lang.String> externProperties, com.google.common.collect.ImmutableSet<java.lang.String> runtimeLibraries)
func TypedAstDeserializer_DeserializedAst_Create(filesystem j2gort.Ext_java_util_concurrent_ConcurrentMap, colorRegistry j2gort.Ext_com_google_common_base_Optional, externProperties j2gort.Ext_com_google_common_collect_ImmutableSet, runtimeLibraries j2gort.Ext_com_google_common_collect_ImmutableSet) TypedAstDeserializer_DeserializedAst {
	_ = filesystem
	_ = colorRegistry
	_ = externProperties
	_ = runtimeLibraries
	return NewAutoValue_TypedAstDeserializer_DeserializedAst(filesystem, colorRegistry, externProperties, runtimeLibraries)
}
