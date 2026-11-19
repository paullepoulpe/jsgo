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
// Port of com/google/javascript/jscomp/serialization/TypedAstSerializer.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_typed_ast_serializer_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_typed_ast_serializer_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.PtrToAny(j2gort.CastPtr[SourceFile](mrp0_).GetProto()) })
}

// Java class: com.google.javascript.jscomp.serialization.TypedAstSerializer
type TypedAstSerializer struct {
	compiler           AbstractCompiler                      `gson:"compiler"`
	serializationMode  *SerializationOptions                 `gson:"serializationMode"`
	stringPool         *StringPool_Builder                   `gson:"stringPool"`
	subtreeSourceFiles j2gort.Deque[*SourceFile]             `gson:"subtreeSourceFiles"`
	sourceFilePointers j2gort.OrderedMap[*SourceFile, int32] `gson:"sourceFilePointers"`
	typeSerializer     TypedAstSerializer_TypeSerializer     `gson:"typeSerializer"`
	previousLine       int32                                 `gson:"previousLine"`
	previousColumn     int32                                 `gson:"previousColumn"`
	idHash_            atomic.Int32
}

func (this_ *TypedAstSerializer) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedAstSerializer) IsJavaNil_() bool { return this_ == nil }
func (*TypedAstSerializer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.TypedAstSerializer", "TypedAstSerializer"
}
func (this_ *TypedAstSerializer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedAstSerializer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.TypedAstSerializer", "com.google.javascript.jscomp.serialization.TypedAstSerializer", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedAstSerializer](v_)
	})
}

func (this_ *TypedAstSerializer) initFields_TypedAstSerializer() {
	_ = this_
	this_.stringPool = StringPool_S_Builder()
	this_.subtreeSourceFiles = j2gort.Deque[*SourceFile]{}
	this_.sourceFilePointers = j2gort.NewOrderedMapCap[*SourceFile, int32](int32(256))

}

// Java constructor: TypedAstSerializer(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.serialization.SerializationOptions serializationMode)
func NewTypedAstSerializer(compiler AbstractCompiler, serializationMode *SerializationOptions) *TypedAstSerializer {
	this_ := &TypedAstSerializer{}
	InitTypedAstSerializer(this_, compiler, serializationMode)
	return this_
}

func InitTypedAstSerializer(this_ *TypedAstSerializer, compiler AbstractCompiler, serializationMode *SerializationOptions) {
	_ = compiler
	_ = serializationMode
	this_.initFields_TypedAstSerializer()
	this_.compiler = compiler
	this_.serializationMode = serializationMode
}

// Java method: TypedAstSerializer.serializeRoots(com.google.javascript.rhino.Node externsRoot, com.google.javascript.rhino.Node jsRoot)
func (this_ *TypedAstSerializer) SerializeRoots(externsRoot pkg_rhino.Node, jsRoot pkg_rhino.Node) *j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst {
	_ = externsRoot
	_ = jsRoot
	j2gort.PurgeAllRuntimePools()

	if !(externsRoot.IsRoot()) {
		j2gort.CheckArgument(false)
	}
	if !(jsRoot.IsRoot()) {
		j2gort.CheckArgument(false)
	}
	if this_.compiler.HasOptimizationColors() {
		this_.typeSerializer = ToIface_TypedAstSerializer_TypeSerializer(this_.priv_TypedAstSerializer_createColorTypeSerializer_7775(this_.compiler, this_.serializationMode, this_.stringPool))
	} else if this_.compiler.HasTypeCheckingRun() {
		this_.typeSerializer = ToIface_TypedAstSerializer_TypeSerializer(TypedAstSerializer_CreateJSTypeSerializer(this_.compiler, this_.serializationMode, this_.stringPool))
	} else {
		this_.typeSerializer = NewTypedAstSerializer_NoOpTypeSerializer()
	}
	var builder *j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst_Builder = j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst_NewBuilder_0()
	_ = builder
	{
		var script pkg_rhino.Node = externsRoot.GetFirstChild()
		_ = script
		for ; script != 0; script = script.GetNext() {
			if NodeUtil_IsFromTypeSummary(script) {
				continue
			}
			builder.AddExternAst_1_Ext_com_google_javascript_jscomp_serialization_LazyAst(this_.priv_TypedAstSerializer_serializeScriptNode_7767(script))
		}
	}
	{
		var script_2 pkg_rhino.Node = jsRoot.GetFirstChild()
		_ = script_2
		for ; script_2 != 0; script_2 = script_2.GetNext() {
			if NodeUtil_IsFromTypeSummary(script_2) {
				continue
			}
			builder.AddCodeAst_1_Ext_com_google_javascript_jscomp_serialization_LazyAst(this_.priv_TypedAstSerializer_serializeScriptNode_7767(script_2))
		}
	}
	var sourceFiles *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFilePool = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFilePool](j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFilePool_NewBuilder_0().AddAllSourceFile(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](j2gort.AsIface[j2gort.Ext_java_util_Set](this_.sourceFilePointers.LiveKeySet()).Stream().Map(lambdaSingleton_typed_ast_serializer_1).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList()))).Build())
	_ = sourceFiles
	if this_.compiler.GetExternProperties() != nil {
		var externsSummary *j2gort.Ext_com_google_javascript_jscomp_serialization_ExternsSummary_Builder = j2gort.Ext_com_google_javascript_jscomp_serialization_ExternsSummary_NewBuilder_0()
		_ = externsSummary
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.compiler.GetExternProperties()))
			for cur_1.HasNext() {
				prop := cur_1.NextStr()
				_ = prop
				externsSummary.AddPropNamePtr(this_.stringPool.Put(prop))
			}
		}
		builder.SetExternsSummary_1_Ext_com_google_javascript_jscomp_serialization_ExternsSummary_Builder(externsSummary)
	}
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypedAst](builder.AddAllRuntimeLibraryToInject(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.serializationMode.RuntimeLibraries())).SetTypePool_1_Ext_com_google_javascript_jscomp_serialization_TypePool(this_.typeSerializer.GenerateTypePool()).SetStringPool_1_Ext_com_google_javascript_jscomp_serialization_StringPoolProto(this_.stringPool.Build().ToProto()).SetSourceFilePool_1_Ext_com_google_javascript_jscomp_serialization_SourceFilePool(sourceFiles).Build())
}

// Java method: TypedAstSerializer.serializeScriptNode(com.google.javascript.rhino.Node script)
func (this_ *TypedAstSerializer) priv_TypedAstSerializer_serializeScriptNode_7767(script pkg_rhino.Node) *j2gort.Ext_com_google_javascript_jscomp_serialization_LazyAst {
	_ = script
	if !(script.IsScript()) {
		j2gort.CheckState(false)
	}
	this_.previousLine = func() int32 { var asgn_2 int32 = int32(0); this_.previousColumn = asgn_2; return asgn_2 }()
	var sourceFile int32 = this_.priv_TypedAstSerializer_getSourceFilePointer_7770(script)
	_ = sourceFile
	var scriptProto *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode = this_.priv_TypedAstSerializer_visit_7769(script)
	_ = scriptProto
	this_.subtreeSourceFiles.Clear()
	var encodedSourceMap string = this_.compiler.GetBase64SourceMapContents(script.GetSourceFileName())
	_ = encodedSourceMap
	var lazyAstBuilder *j2gort.Ext_com_google_javascript_jscomp_serialization_LazyAst_Builder = j2gort.Ext_com_google_javascript_jscomp_serialization_LazyAst_NewBuilder_0().SetScript(scriptProto.ToByteString()).SetSourceFile(sourceFile)
	_ = lazyAstBuilder
	if encodedSourceMap != j2gort.NullString {
		lazyAstBuilder.SetSourceMappingUrl(encodedSourceMap)
	}
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_LazyAst](lazyAstBuilder.Build())
}

// Java method: TypedAstSerializer.createWithPositionInfo(com.google.javascript.rhino.Node n)
func (this_ *TypedAstSerializer) priv_TypedAstSerializer_createWithPositionInfo_7768(n pkg_rhino.Node) *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_Builder {
	_ = n
	if !(n.GetLength() >= int32(0)) {
		j2gort.CheckState(false)
	}
	var currentLine int32 = n.GetLineno()
	_ = currentLine
	var currentColumn int32 = n.GetCharno()
	_ = currentColumn
	var builder *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_Builder = j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_NewBuilder_0().SetRelativeLine((currentLine - this_.previousLine)).SetRelativeColumn((currentColumn - this_.previousColumn))
	_ = builder
	this_.previousLine = currentLine
	this_.previousColumn = currentColumn
	return builder
}

// Java method: TypedAstSerializer.visit(com.google.javascript.rhino.Node n)
func (this_ *TypedAstSerializer) priv_TypedAstSerializer_visit_7769(n pkg_rhino.Node) *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode {
	_ = n
	var builder *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_Builder = this_.priv_TypedAstSerializer_createWithPositionInfo_7768(n)
	_ = builder
	this_.priv_TypedAstSerializer_addType_7771(n, builder)
	var serializedJsdoc *j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc = JSDocSerializer_SerializeJsdoc(n.GetJSDocInfo(), this_.stringPool)
	_ = serializedJsdoc
	if serializedJsdoc != nil {
		builder.SetJsdoc_1_Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc(serializedJsdoc)
	}
	builder.SetKind(this_.priv_TypedAstSerializer_kindTranslator_7774(n))
	this_.priv_TypedAstSerializer_valueTranslator_7772(builder, n)
	builder.SetBooleanProperties(n.SerializeProperties())
	var sourceFile int32 = this_.priv_TypedAstSerializer_getSourceFilePointer_7770(n)
	_ = sourceFile
	builder.SetSourceFile(sourceFile)
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			builder.AddChild_1_Ext_com_google_javascript_jscomp_serialization_AstNode(this_.priv_TypedAstSerializer_visit_7769(child))
		}
	}
	var shadowedCode pkg_rhino.Node = n.GetClosureUnawareShadow()
	_ = shadowedCode
	if shadowedCode != 0 {
		builder.AddChild_1_Ext_com_google_javascript_jscomp_serialization_AstNode(this_.priv_TypedAstSerializer_visit_7769(shadowedCode.GetFirstFirstChild().GetFirstChild()))
	}
	if sourceFile != int32(0) {
		this_.subtreeSourceFiles.RemoveLast()
	}
	this_.priv_TypedAstSerializer_setOriginalName_7773(builder, n)
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode](builder.Build())
}

// Java method: TypedAstSerializer.getSourceFilePointer(com.google.javascript.rhino.Node n)
func (this_ *TypedAstSerializer) priv_TypedAstSerializer_getSourceFilePointer_7770(n pkg_rhino.Node) int32 {
	_ = n
	var sourceFile *SourceFile = j2gort.CastPtr[SourceFile](n.GetStaticSourceFile())
	_ = sourceFile
	if sourceFile == nil {
		j2gort.CheckStateFmt(this_.subtreeSourceFiles.PeekLast().IsExtern(), "Unexpected null SourceFile for node %s with parent %s", j2gort.BoxString(n.ToStringTree()), j2gort.BoxNode(n.GetParent()))
		return int32(0)
	}
	{
		var iifeCond_1 bool
		{

			rp_3 := sourceFile
			ap_4 := this_.subtreeSourceFiles.PeekLast()
			iifeCond_1 = j2gort.CheckNotNullPtr(rp_3) == ap_4
		}
		if iifeCond_1 {
			return int32(0)
		}
	}
	this_.subtreeSourceFiles.AddLast(sourceFile)
	ptr, ptrOk_ := this_.sourceFilePointers.GetOk(sourceFile)
	_ = ptrOk_
	_ = ptr
	if !ptrOk_ {
		ptr, ptrOk_ = (int32(1) + this_.sourceFilePointers.Size()), true
		this_.sourceFilePointers.Put(sourceFile, ptr)
	}
	return func() int32 {
		if !ptrOk_ {
			j2gort.ThrowNPE(

			// Java method: TypedAstSerializer.addType(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.serialization.AstNode.Builder builder)
			)
		}
		return ptr
	}()
}

func (this_ *TypedAstSerializer) priv_TypedAstSerializer_addType_7771(n pkg_rhino.Node, builder *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_Builder) {
	_ = n
	_ = builder
	this_.typeSerializer.AddTypeForNode(n, builder)
}

// Java method: TypedAstSerializer.valueTranslator(com.google.javascript.jscomp.serialization.AstNode.Builder builder, com.google.javascript.rhino.Node n)
func (this_ *TypedAstSerializer) priv_TypedAstSerializer_valueTranslator_7772(builder *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_Builder, n pkg_rhino.Node) {
	_ = builder
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(32) /* Token.NAME */, int32(35) /* Token.STRINGLIT */, int32(100) /* Token.STRING_KEY */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(99) /* Token.LABEL_NAME */, int32(119) /* Token.IMPORT_STAR */ :
		builder.SetStringValuePointer(this_.stringPool.Put(n.GetString()))
		return
	case int32(133) /* Token.TEMPLATELIT_STRING */ :
		builder.SetTemplateStringValue_1_Ext_com_google_javascript_jscomp_serialization_TemplateStringValue(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TemplateStringValue](j2gort.Ext_com_google_javascript_jscomp_serialization_TemplateStringValue_NewBuilder_0().SetRawStringPointer(this_.stringPool.Put(n.GetRawString())).SetCookedStringPointer(func() int32 {
			if n.GetCookedString() == j2gort.NullString {
				return int32(-1)
			}
			return this_.stringPool.Put(n.GetCookedString())
		}()).Build()))
		return
	case int32(33) /* Token.NUMBER */ :
		builder.SetDoubleValue(n.GetDouble())
		return
	case int32(34) /* Token.BIGINT */ :
		builder.SetStringValuePointer(this_.stringPool.Put(n.GetBigInt().String()))
		return
	default:
		return
	}
}

// Java method: TypedAstSerializer.setOriginalName(com.google.javascript.jscomp.serialization.AstNode.Builder builder, com.google.javascript.rhino.Node n)
func (this_ *TypedAstSerializer) priv_TypedAstSerializer_setOriginalName_7773(builder *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_Builder, n pkg_rhino.Node) {
	_ = builder
	_ = n
	var originalName string = n.GetOriginalName()
	_ = originalName
	if originalName == j2gort.NullString {
		builder.SetOriginalNamePointer(int32(0))
		return
	}
	builder.SetOriginalNamePointer(this_.stringPool.Put(originalName))
}

// Java method: TypedAstSerializer.kindTranslator(com.google.javascript.rhino.Node n)
func (this_ *TypedAstSerializer) priv_TypedAstSerializer_kindTranslator_7774(n pkg_rhino.Node) *j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_SOURCE_FILE
	case int32(33) /* Token.NUMBER */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_NUMBER_LITERAL
	case int32(35) /* Token.STRINGLIT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_STRING_LITERAL
	case int32(34) /* Token.BIGINT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_BIGINT_LITERAL
	case int32(42) /* Token.REGEXP */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_REGEX_LITERAL
	case int32(38) /* Token.FALSE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_FALSE
	case int32(39) /* Token.TRUE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_TRUE
	case int32(36) /* Token.NULL */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_NULL
	case int32(37) /* Token.THIS */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_THIS
	case int32(87) /* Token.VOID */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_VOID
	case int32(46) /* Token.ARRAYLIT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ARRAY_LITERAL
	case int32(47) /* Token.OBJECTLIT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_OBJECT_LITERAL
	case int32(32) /* Token.NAME */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_IDENTIFIER
	case int32(51) /* Token.ASSIGN */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGNMENT
	case int32(28) /* Token.CALL */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_CALL
	case int32(23) /* Token.NEW */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_NEW
	case int32(26) /* Token.GETPROP */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_PROPERTY_ACCESS
	case int32(27) /* Token.GETELEM */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ELEMENT_ACCESS
	case int32(50) /* Token.COMMA */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_COMMA
	case int32(68) /* Token.OR */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_BOOLEAN_OR
	case int32(69) /* Token.AND */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_BOOLEAN_AND
	case int32(67) /* Token.HOOK */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_HOOK
	case int32(4) /* Token.EQ */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_EQUAL
	case int32(5) /* Token.NE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_NOT_EQUAL
	case int32(6) /* Token.LT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_LESS_THAN
	case int32(7) /* Token.LE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_LESS_THAN_EQUAL
	case int32(8) /* Token.GT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_GREATER_THAN
	case int32(9) /* Token.GE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_GREATER_THAN_EQUAL
	case int32(40) /* Token.SHEQ */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_TRIPLE_EQUAL
	case int32(41) /* Token.SHNE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_NOT_TRIPLE_EQUAL
	case int32(19) /* Token.NOT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_NOT
	case int32(21) /* Token.POS */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_POSITIVE
	case int32(22) /* Token.NEG */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_NEGATIVE
	case int32(25) /* Token.TYPEOF */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_TYPEOF
	case int32(45) /* Token.INSTANCEOF */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_INSTANCEOF
	case int32(44) /* Token.IN */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_IN
	case int32(13) /* Token.ADD */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ADD
	case int32(14) /* Token.SUB */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_SUBTRACT
	case int32(15) /* Token.MUL */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_MULTIPLY
	case int32(16) /* Token.DIV */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_DIVIDE
	case int32(17) /* Token.MOD */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_MODULO
	case int32(18) /* Token.EXPONENT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_EXPONENT
	case int32(20) /* Token.BITNOT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_BITWISE_NOT
	case int32(1) /* Token.BITOR */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_BITWISE_OR
	case int32(3) /* Token.BITAND */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_BITWISE_AND
	case int32(2) /* Token.BITXOR */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_BITWISE_XOR
	case int32(10) /* Token.LSH */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_LEFT_SHIFT
	case int32(11) /* Token.RSH */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_RIGHT_SHIFT
	case int32(12) /* Token.URSH */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_UNSIGNED_RIGHT_SHIFT
	case int32(71) /* Token.INC */ :
		if n.GetBooleanProp(pkg_rhino.StaticGet_Node_INCRDECR_PROP()) {
			return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_POST_INCREMENT
		} else {
			return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_PRE_INCREMENT
		}
	case int32(72) /* Token.DEC */ :
		if n.GetBooleanProp(pkg_rhino.StaticGet_Node_INCRDECR_PROP()) {
			return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_POST_DECREMENT
		} else {
			return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_PRE_DECREMENT
		}
	case int32(58) /* Token.ASSIGN_ADD */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_ADD
	case int32(59) /* Token.ASSIGN_SUB */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_SUBTRACT
	case int32(60) /* Token.ASSIGN_MUL */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_MULTIPLY
	case int32(61) /* Token.ASSIGN_DIV */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_DIVIDE
	case int32(62) /* Token.ASSIGN_MOD */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_MODULO
	case int32(63) /* Token.ASSIGN_EXPONENT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_EXPONENT
	case int32(52) /* Token.ASSIGN_BITOR */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_BITWISE_OR
	case int32(54) /* Token.ASSIGN_BITAND */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_BITWISE_AND
	case int32(53) /* Token.ASSIGN_BITXOR */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_BITWISE_XOR
	case int32(55) /* Token.ASSIGN_LSH */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_LEFT_SHIFT
	case int32(56) /* Token.ASSIGN_RSH */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_RIGHT_SHIFT
	case int32(57) /* Token.ASSIGN_URSH */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_UNSIGNED_RIGHT_SHIFT
	case int32(114) /* Token.YIELD */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_YIELD
	case int32(115) /* Token.AWAIT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_AWAIT
	case int32(24) /* Token.DELPROP */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_DELETE
	case int32(130) /* Token.TAGGED_TEMPLATELIT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_TAGGED_TEMPLATELIT
	case int32(131) /* Token.TEMPLATELIT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_TEMPLATELIT
	case int32(132) /* Token.TEMPLATELIT_SUB */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_TEMPLATELIT_SUB
	case int32(133) /* Token.TEMPLATELIT_STRING */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_TEMPLATELIT_STRING
	case int32(135) /* Token.NEW_TARGET */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_NEW_TARGET
	case int32(129) /* Token.COMPUTED_PROP */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_COMPUTED_PROP
	case int32(136) /* Token.IMPORT_META */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_IMPORT_META
	case int32(29) /* Token.OPTCHAIN_GETPROP */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_OPTCHAIN_PROPERTY_ACCESS
	case int32(31) /* Token.OPTCHAIN_CALL */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_OPTCHAIN_CALL
	case int32(30) /* Token.OPTCHAIN_GETELEM */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_OPTCHAIN_ELEMENT_ACCESS
	case int32(70) /* Token.COALESCE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_COALESCE
	case int32(124) /* Token.DYNAMIC_IMPORT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_DYNAMIC_IMPORT
	case int32(64) /* Token.ASSIGN_OR */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_OR
	case int32(65) /* Token.ASSIGN_AND */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_AND
	case int32(66) /* Token.ASSIGN_COALESCE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ASSIGN_COALESCE
	case int32(93) /* Token.EXPR_RESULT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_EXPRESSION_STATEMENT
	case int32(82) /* Token.BREAK */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_BREAK_STATEMENT
	case int32(83) /* Token.CONTINUE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_CONTINUE_STATEMENT
	case int32(98) /* Token.DEBUGGER */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_DEBUGGER_STATEMENT
	case int32(79) /* Token.DO */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_DO_STATEMENT
	case int32(80) /* Token.FOR */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_FOR_STATEMENT
	case int32(81) /* Token.FOR_IN */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_FOR_IN_STATEMENT
	case int32(112) /* Token.FOR_OF */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_FOR_OF_STATEMENT
	case int32(113) /* Token.FOR_AWAIT_OF */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_FOR_AWAIT_OF_STATEMENT
	case int32(74) /* Token.IF */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_IF_STATEMENT
	case int32(0) /* Token.RETURN */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_RETURN_STATEMENT
	case int32(75) /* Token.SWITCH */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_SWITCH_STATEMENT
	case int32(43) /* Token.THROW */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_THROW_STATEMENT
	case int32(48) /* Token.TRY */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_TRY_STATEMENT
	case int32(78) /* Token.WHILE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_WHILE_STATEMENT
	case int32(88) /* Token.EMPTY */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_EMPTY
	case int32(85) /* Token.WITH */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_WITH
	case int32(116) /* Token.IMPORT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_IMPORT
	case int32(120) /* Token.EXPORT */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_EXPORT
	case int32(84) /* Token.VAR */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_VAR_DECLARATION
	case int32(97) /* Token.CONST */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_CONST_DECLARATION
	case int32(111) /* Token.LET */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_LET_DECLARATION
	case int32(73) /* Token.FUNCTION */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_FUNCTION_LITERAL
	case int32(105) /* Token.CLASS */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_CLASS_LITERAL
	case int32(90) /* Token.BLOCK */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_BLOCK
	case int32(92) /* Token.LABEL */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_LABELED_STATEMENT
	case int32(99) /* Token.LABEL_NAME */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_LABELED_NAME
	case int32(106) /* Token.CLASS_MEMBERS */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_CLASS_MEMBERS
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_METHOD_DECLARATION
	case int32(108) /* Token.MEMBER_FIELD_DEF */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_FIELD_DECLARATION
	case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_COMPUTED_PROP_FIELD
	case int32(49) /* Token.PARAM_LIST */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_PARAMETER_LIST
	case int32(100) /* Token.STRING_KEY */ :
		if n.IsQuotedStringKey() {
			return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_QUOTED_STRING_KEY
		} else {
			return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_RENAMABLE_STRING_KEY
		}
	case int32(76) /* Token.CASE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_CASE
	case int32(77) /* Token.DEFAULT_CASE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_DEFAULT_CASE
	case int32(86) /* Token.CATCH */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_CATCH
	case int32(110) /* Token.SUPER */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_SUPER
	case int32(102) /* Token.ARRAY_PATTERN */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ARRAY_PATTERN
	case int32(103) /* Token.OBJECT_PATTERN */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_OBJECT_PATTERN
	case int32(104) /* Token.DESTRUCTURING_LHS */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_DESTRUCTURING_LHS
	case int32(134) /* Token.DEFAULT_VALUE */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_DEFAULT_VALUE
	case int32(95) /* Token.GETTER_DEF */ :
		if n.IsQuotedStringKey() {
			return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_QUOTED_GETTER_DEF
		} else {
			return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_RENAMABLE_GETTER_DEF
		}
	case int32(96) /* Token.SETTER_DEF */ :
		if n.IsQuotedStringKey() {
			return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_QUOTED_SETTER_DEF
		} else {
			return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_RENAMABLE_SETTER_DEF
		}
	case int32(117) /* Token.IMPORT_SPECS */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_IMPORT_SPECS
	case int32(118) /* Token.IMPORT_SPEC */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_IMPORT_SPEC
	case int32(119) /* Token.IMPORT_STAR */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_IMPORT_STAR
	case int32(121) /* Token.EXPORT_SPECS */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_EXPORT_SPECS
	case int32(122) /* Token.EXPORT_SPEC */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_EXPORT_SPEC
	case int32(123) /* Token.MODULE_BODY */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_MODULE_BODY
	case int32(125) /* Token.ITER_REST */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ITER_REST
	case int32(127) /* Token.ITER_SPREAD */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_ITER_SPREAD
	case int32(126) /* Token.OBJECT_REST */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_OBJECT_REST
	case int32(128) /* Token.OBJECT_SPREAD */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_OBJECT_SPREAD
	case int32(91) /* Token.SWITCH_BODY */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_SWITCH_BODY
	case int32(89) /* Token.ROOT */, int32(137) /* Token.STRING_TYPE */, int32(138) /* Token.BOOLEAN_TYPE */, int32(139) /* Token.NUMBER_TYPE */, int32(140) /* Token.FUNCTION_TYPE */, int32(141) /* Token.PARAMETERIZED_TYPE */, int32(142) /* Token.UNION_TYPE */, int32(143) /* Token.ANY_TYPE */, int32(144) /* Token.NULLABLE_TYPE */, int32(145) /* Token.VOID_TYPE */, int32(146) /* Token.REST_PARAMETER_TYPE */, int32(147) /* Token.NAMED_TYPE */, int32(149) /* Token.RECORD_TYPE */, int32(150) /* Token.UNDEFINED_TYPE */, int32(148) /* Token.OPTIONAL_PARAMETER */, int32(151) /* Token.ARRAY_TYPE */, int32(152) /* Token.GENERIC_TYPE */, int32(153) /* Token.GENERIC_TYPE_LIST */, int32(164) /* Token.INTERFACE */, int32(165) /* Token.INTERFACE_EXTENDS */, int32(166) /* Token.INTERFACE_MEMBERS */, int32(172) /* Token.MEMBER_VARIABLE_DEF */, int32(167) /* Token.ENUM */, int32(168) /* Token.ENUM_MEMBERS */, int32(169) /* Token.IMPLEMENTS */, int32(170) /* Token.TYPE_ALIAS */, int32(171) /* Token.DECLARE */, int32(173) /* Token.INDEX_SIGNATURE */, int32(174) /* Token.CALL_SIGNATURE */, int32(175) /* Token.NAMESPACE */, int32(176) /* Token.NAMESPACE_ELEMENTS */, int32(154) /* Token.ANNOTATION */, int32(155) /* Token.PIPE */, int32(156) /* Token.STAR */, int32(157) /* Token.EOC */, int32(158) /* Token.QMARK */, int32(159) /* Token.BANG */, int32(160) /* Token.EQUALS */, int32(161) /* Token.LB */, int32(162) /* Token.LC */, int32(163) /* Token.COLON */, int32(101) /* Token.CAST */, int32(177) /* Token.PLACEHOLDER1 */, int32(178) /* Token.PLACEHOLDER2 */, int32(179) /* Token.PLACEHOLDER3 */ :
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unserializable token for node: "+j2gort.StrVal(n))))
}

// Java static method: TypedAstSerializer.createJSTypeSerializer(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.serialization.SerializationOptions serializationMode, com.google.javascript.jscomp.serialization.StringPool.Builder stringPoolBuilder)
func TypedAstSerializer_CreateJSTypeSerializer(compiler AbstractCompiler, serializationMode *SerializationOptions, stringPoolBuilder *StringPool_Builder) *TypedAstSerializer_JSTypeSerializer {
	_ = compiler
	_ = serializationMode
	_ = stringPoolBuilder
	var serializeTypesToPointers *SerializeTypesToPointers = SerializeTypesToPointers_Create(compiler, stringPoolBuilder, serializationMode)
	_ = serializeTypesToPointers
	serializeTypesToPointers.GatherTypesOnAst(compiler.GetRoot())
	return NewTypedAstSerializer_JSTypeSerializer(serializeTypesToPointers.GetTypePointersByJstype(), serializeTypesToPointers.GetTypePool())
}

// Java method: TypedAstSerializer.createColorTypeSerializer(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.serialization.SerializationOptions serializationMode, com.google.javascript.jscomp.serialization.StringPool.Builder stringPoolBuilder)
func (this_ *TypedAstSerializer) priv_TypedAstSerializer_createColorTypeSerializer_7775(compiler AbstractCompiler, serializationMode *SerializationOptions, stringPoolBuilder *StringPool_Builder) *TypedAstSerializer_ColorTypeSerializer {
	_ = compiler
	_ = serializationMode
	_ = stringPoolBuilder
	var usedPropertyNames *j2gort.Ext_java_util_LinkedHashSet = TypedAstSerializer_CollectUsedPropertyNames(compiler)
	_ = usedPropertyNames
	var val_5 ColorSerializer
	InitColorSerializer(&val_5, serializationMode, func() ColorSerializer_StringPoolIndexFn {
		bound_6 := stringPoolBuilder
		return NewLambda_ColorSerializer_StringPoolIndexFn(func(mrp0_ string) int32 { return bound_6.Put(mrp0_) })
	}(), func() ColorSerializer_PropertyFilter {
		bound_7 := usedPropertyNames
		return NewLambda_ColorSerializer_PropertyFilter(func(mrp0_ string) bool { return bound_7.Contains(j2gort.BoxString(mrp0_)) })
	}())
	var colorSerializer *ColorSerializer = &val_5
	_ = colorSerializer
	return NewTypedAstSerializer_ColorTypeSerializer(colorSerializer, compiler.GetColorRegistry())
}

// Java static method: TypedAstSerializer.collectUsedPropertyNames(com.google.javascript.jscomp.AbstractCompiler compiler)
func TypedAstSerializer_CollectUsedPropertyNames(compiler AbstractCompiler) *j2gort.Ext_java_util_LinkedHashSet {
	_ = compiler
	var propertyNames *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	_ = propertyNames
	NodeTraversal_Traverse(compiler, compiler.GetRoot(), NewTypedAstSerializer__Anon_126(propertyNames))
	return propertyNames
}

// Java interface: com.google.javascript.jscomp.serialization.TypedAstSerializer.TypeSerializer
type TypedAstSerializer_TypeSerializer interface {
	Implements_jscomp_TypedAstSerializer_TypeSerializer()
	IsJavaNil_() bool
	AddTypeForNode(pkg_rhino.Node, *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_Builder)
	GenerateTypePool() *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool
}

func ToIface_TypedAstSerializer_TypeSerializer[T interface {
	TypedAstSerializer_TypeSerializer
	comparable
}](v T) TypedAstSerializer_TypeSerializer {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.TypedAstSerializer.TypeSerializer", "com.google.javascript.jscomp.serialization.TypedAstSerializer$TypeSerializer", func(v_ any) bool {
		return j2gort.InstanceOf[TypedAstSerializer_TypeSerializer](v_)
	})
}

// Java class: com.google.javascript.jscomp.serialization.TypedAstSerializer.NoOpTypeSerializer
type TypedAstSerializer_NoOpTypeSerializer struct {
	idHash_ atomic.Int32
}

func (this_ *TypedAstSerializer_NoOpTypeSerializer) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedAstSerializer_NoOpTypeSerializer) IsJavaNil_() bool { return this_ == nil }
func (*TypedAstSerializer_NoOpTypeSerializer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.TypedAstSerializer$NoOpTypeSerializer", "NoOpTypeSerializer"
}
func (this_ *TypedAstSerializer_NoOpTypeSerializer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedAstSerializer_NoOpTypeSerializer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.TypedAstSerializer.NoOpTypeSerializer", "com.google.javascript.jscomp.serialization.TypedAstSerializer$NoOpTypeSerializer", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedAstSerializer_NoOpTypeSerializer](v_)
	})
}
func (*TypedAstSerializer_NoOpTypeSerializer) Implements_jscomp_TypedAstSerializer_TypeSerializer() {}

func (this_ *TypedAstSerializer_NoOpTypeSerializer) initFields_TypedAstSerializer_NoOpTypeSerializer() {
	_ = this_
}

// Java constructor: NoOpTypeSerializer()
func NewTypedAstSerializer_NoOpTypeSerializer() *TypedAstSerializer_NoOpTypeSerializer {
	this_ := &TypedAstSerializer_NoOpTypeSerializer{}
	InitTypedAstSerializer_NoOpTypeSerializer(this_)
	return this_
}

func InitTypedAstSerializer_NoOpTypeSerializer(this_ *TypedAstSerializer_NoOpTypeSerializer) {
	this_.initFields_TypedAstSerializer_NoOpTypeSerializer()
}

// Java method: NoOpTypeSerializer.addTypeForNode(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.serialization.AstNode.Builder astNodeBuilder)
func (this_ *TypedAstSerializer_NoOpTypeSerializer) AddTypeForNode(node pkg_rhino.Node, astNodeBuilder *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_Builder) {
	_ = node
	_ = astNodeBuilder
}

// Java method: NoOpTypeSerializer.generateTypePool()
func (this_ *TypedAstSerializer_NoOpTypeSerializer) GenerateTypePool() *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool {
	return j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool_GetDefaultInstance()
}

// Java class: com.google.javascript.jscomp.serialization.TypedAstSerializer.JSTypeSerializer
type TypedAstSerializer_JSTypeSerializer struct {
	typesToPointers *pkg_base.LinkedIdentityHashMap[pkg_rhino.JSType, *j2gort.BoxInt] `gson:"typesToPointers"`
	typePool        *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool   `gson:"typePool"`
	idHash_         atomic.Int32
}

func (this_ *TypedAstSerializer_JSTypeSerializer) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedAstSerializer_JSTypeSerializer) IsJavaNil_() bool { return this_ == nil }
func (*TypedAstSerializer_JSTypeSerializer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.TypedAstSerializer$JSTypeSerializer", "JSTypeSerializer"
}
func (this_ *TypedAstSerializer_JSTypeSerializer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedAstSerializer_JSTypeSerializer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.TypedAstSerializer.JSTypeSerializer", "com.google.javascript.jscomp.serialization.TypedAstSerializer$JSTypeSerializer", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedAstSerializer_JSTypeSerializer](v_)
	})
}
func (*TypedAstSerializer_JSTypeSerializer) Implements_jscomp_TypedAstSerializer_TypeSerializer() {}

func (this_ *TypedAstSerializer_JSTypeSerializer) initFields_TypedAstSerializer_JSTypeSerializer() {
	_ = this_
}

// Java constructor: JSTypeSerializer(com.google.javascript.jscomp.base.LinkedIdentityHashMap<com.google.javascript.rhino.jstype.JSType,java.lang.Integer> typesToPointers, com.google.javascript.jscomp.serialization.TypePool typePool)
func NewTypedAstSerializer_JSTypeSerializer(typesToPointers *pkg_base.LinkedIdentityHashMap[pkg_rhino.JSType, *j2gort.BoxInt], typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool) *TypedAstSerializer_JSTypeSerializer {
	this_ := &TypedAstSerializer_JSTypeSerializer{}
	InitTypedAstSerializer_JSTypeSerializer(this_, typesToPointers, typePool)
	return this_
}

func InitTypedAstSerializer_JSTypeSerializer(this_ *TypedAstSerializer_JSTypeSerializer, typesToPointers *pkg_base.LinkedIdentityHashMap[pkg_rhino.JSType, *j2gort.BoxInt], typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool) {
	_ = typesToPointers
	_ = typePool
	this_.initFields_TypedAstSerializer_JSTypeSerializer()
	this_.typesToPointers = typesToPointers
	this_.typePool = typePool
}

// Java method: JSTypeSerializer.addTypeForNode(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.serialization.AstNode.Builder astNodeBuilder)
func (this_ *TypedAstSerializer_JSTypeSerializer) AddTypeForNode(node pkg_rhino.Node, astNodeBuilder *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_Builder) {
	_ = node
	_ = astNodeBuilder
	var type_ pkg_rhino.JSType = node.GetJSType()
	_ = type_
	if type_ != nil && !type_.IsJavaNil_() {
		astNodeBuilder.SetType(j2gort.UnboxInt(j2gort.CheckNotNullPtr(this_.typesToPointers.Get(type_))))
	}
}

// Java method: JSTypeSerializer.generateTypePool()
func (this_ *TypedAstSerializer_JSTypeSerializer) GenerateTypePool() *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool {
	return this_.typePool
}

// Java class: com.google.javascript.jscomp.serialization.TypedAstSerializer.ColorTypeSerializer
type TypedAstSerializer_ColorTypeSerializer struct {
	colorSerializer *ColorSerializer          `gson:"colorSerializer"`
	colorRegistry   *pkg_colors.ColorRegistry `gson:"colorRegistry"`
	idHash_         atomic.Int32
}

func (this_ *TypedAstSerializer_ColorTypeSerializer) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedAstSerializer_ColorTypeSerializer) IsJavaNil_() bool { return this_ == nil }
func (*TypedAstSerializer_ColorTypeSerializer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.TypedAstSerializer$ColorTypeSerializer", "ColorTypeSerializer"
}
func (this_ *TypedAstSerializer_ColorTypeSerializer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedAstSerializer_ColorTypeSerializer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.TypedAstSerializer.ColorTypeSerializer", "com.google.javascript.jscomp.serialization.TypedAstSerializer$ColorTypeSerializer", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedAstSerializer_ColorTypeSerializer](v_)
	})
}
func (*TypedAstSerializer_ColorTypeSerializer) Implements_jscomp_TypedAstSerializer_TypeSerializer() {
}

func (this_ *TypedAstSerializer_ColorTypeSerializer) initFields_TypedAstSerializer_ColorTypeSerializer() {
	_ = this_
}

// Java constructor: ColorTypeSerializer(com.google.javascript.jscomp.serialization.ColorSerializer colorSerializer, com.google.javascript.jscomp.colors.ColorRegistry colorRegistry)
func NewTypedAstSerializer_ColorTypeSerializer(colorSerializer *ColorSerializer, colorRegistry *pkg_colors.ColorRegistry) *TypedAstSerializer_ColorTypeSerializer {
	this_ := &TypedAstSerializer_ColorTypeSerializer{}
	InitTypedAstSerializer_ColorTypeSerializer(this_, colorSerializer, colorRegistry)
	return this_
}

func InitTypedAstSerializer_ColorTypeSerializer(this_ *TypedAstSerializer_ColorTypeSerializer, colorSerializer *ColorSerializer, colorRegistry *pkg_colors.ColorRegistry) {
	_ = colorSerializer
	_ = colorRegistry
	this_.initFields_TypedAstSerializer_ColorTypeSerializer()
	this_.colorSerializer = colorSerializer
	this_.colorRegistry = colorRegistry
}

// Java method: ColorTypeSerializer.addTypeForNode(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.serialization.AstNode.Builder astNodeBuilder)
func (this_ *TypedAstSerializer_ColorTypeSerializer) AddTypeForNode(node pkg_rhino.Node, astNodeBuilder *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_Builder) {
	_ = node
	_ = astNodeBuilder
	var color *pkg_colors.Color = node.GetColor()
	_ = color
	if color != nil {
		astNodeBuilder.SetType(this_.colorSerializer.AddColor(color))
	}
}

// Java method: ColorTypeSerializer.generateTypePool()
func (this_ *TypedAstSerializer_ColorTypeSerializer) GenerateTypePool() *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool {
	var mismatchLocationsForDebugging *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap = this_.colorRegistry.GetMismatchLocationsForDebugging()
	_ = mismatchLocationsForDebugging
	return this_.colorSerializer.GenerateTypePool(func() ColorSerializer_DisambiguationSupertypesFn {
		bound_8 := this_.colorRegistry
		return NewLambda_ColorSerializer_DisambiguationSupertypesFn(func(mrp0_ *pkg_colors.Color) []*pkg_colors.Color {
			return bound_8.GetDisambiguationSupertypesArray(mrp0_)
		})
	}(), j2gort.NewLambda_Ext_java_util_function_Function(func(lp_9 any) any {
		color := j2gort.CastPtr[pkg_colors.Color](lp_9)
		_ = color
		return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](mismatchLocationsForDebugging.Get(j2gort.BoxNullableVal(color.GetId())))
	}))
}

// Java class:
type TypedAstSerializer__Anon_126 struct {
	cap_0_propertyNames *j2gort.Ext_java_util_LinkedHashSet
	idHash_             atomic.Int32
}

func (this_ *TypedAstSerializer__Anon_126) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedAstSerializer__Anon_126) IsJavaNil_() bool { return this_ == nil }
func (*TypedAstSerializer__Anon_126) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.TypedAstSerializer$1", ""
}
func (this_ *TypedAstSerializer__Anon_126) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedAstSerializer__Anon_126) HasIdentityEquals_()                       {}
func (*TypedAstSerializer__Anon_126) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *TypedAstSerializer__Anon_126) initFields_TypedAstSerializer__Anon_126() {
	_ = this_
}

func NewTypedAstSerializer__Anon_126(cap_p_propertyNames *j2gort.Ext_java_util_LinkedHashSet) *TypedAstSerializer__Anon_126 {
	this_ := &TypedAstSerializer__Anon_126{}
	this_.cap_0_propertyNames = cap_p_propertyNames
	this_.initFields_TypedAstSerializer__Anon_126()
	return this_
}

// Java method: .shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypedAstSerializer__Anon_126) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(n.IsScript()) || !(NodeUtil_IsFromTypeSummary(n)))
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypedAstSerializer__Anon_126) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		j2gort.SetAddStr(this_.cap_0_propertyNames, n.GetString())
	case int32(100) /* Token.STRING_KEY */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
		if !(n.IsQuotedStringKey()) {
			j2gort.SetAddStr(this_.cap_0_propertyNames, n.GetString())
		}
	default:
	}
}
func (*TypedAstSerializer__Anon_126) CanReleaseSyntacticScopes_() {}
