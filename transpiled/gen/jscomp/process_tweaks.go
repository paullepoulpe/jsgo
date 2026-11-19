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
// Port of com/google/javascript/jscomp/ProcessTweaks.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ProcessTweaks
type ProcessTweaks struct {
	compiler    AbstractCompiler `gson:"compiler"`
	idHash_     atomic.Int32
	stripTweaks bool `gson:"stripTweaks"`
}

func (this_ *ProcessTweaks) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessTweaks) IsJavaNil_() bool { return this_ == nil }
func (*ProcessTweaks) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessTweaks", "ProcessTweaks"
}
func (this_ *ProcessTweaks) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessTweaks) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessTweaks", "com.google.javascript.jscomp.ProcessTweaks", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessTweaks](v_)
	})
}
func (*ProcessTweaks) Implements_jscomp_CompilerPass() {}

var static_ProcessTweaks_ID_MATCHER j2gort.Ext_com_google_common_base_CharMatcher

func StaticGet_ProcessTweaks_ID_MATCHER() j2gort.Ext_com_google_common_base_CharMatcher {
	clinit_ProcessTweaks()
	return static_ProcessTweaks_ID_MATCHER
}
func StaticSet_ProcessTweaks_ID_MATCHER(v j2gort.Ext_com_google_common_base_CharMatcher) {
	clinit_ProcessTweaks()
	static_ProcessTweaks_ID_MATCHER = v
}

var static_ProcessTweaks_TWEAK_MULTIPLY_REGISTERED_ERROR *DiagnosticType

func StaticGet_ProcessTweaks_TWEAK_MULTIPLY_REGISTERED_ERROR() *DiagnosticType {
	clinit_ProcessTweaks()
	return static_ProcessTweaks_TWEAK_MULTIPLY_REGISTERED_ERROR
}
func StaticSet_ProcessTweaks_TWEAK_MULTIPLY_REGISTERED_ERROR(v *DiagnosticType) {
	clinit_ProcessTweaks()
	static_ProcessTweaks_TWEAK_MULTIPLY_REGISTERED_ERROR = v
}

var static_ProcessTweaks_NON_LITERAL_TWEAK_ID_ERROR *DiagnosticType

func StaticGet_ProcessTweaks_NON_LITERAL_TWEAK_ID_ERROR() *DiagnosticType {
	clinit_ProcessTweaks()
	return static_ProcessTweaks_NON_LITERAL_TWEAK_ID_ERROR
}
func StaticSet_ProcessTweaks_NON_LITERAL_TWEAK_ID_ERROR(v *DiagnosticType) {
	clinit_ProcessTweaks()
	static_ProcessTweaks_NON_LITERAL_TWEAK_ID_ERROR = v
}

var static_ProcessTweaks_INVALID_TWEAK_DEFAULT_VALUE_WARNING *DiagnosticType

func StaticGet_ProcessTweaks_INVALID_TWEAK_DEFAULT_VALUE_WARNING() *DiagnosticType {
	clinit_ProcessTweaks()
	return static_ProcessTweaks_INVALID_TWEAK_DEFAULT_VALUE_WARNING
}
func StaticSet_ProcessTweaks_INVALID_TWEAK_DEFAULT_VALUE_WARNING(v *DiagnosticType) {
	clinit_ProcessTweaks()
	static_ProcessTweaks_INVALID_TWEAK_DEFAULT_VALUE_WARNING = v
}

var static_ProcessTweaks_NON_GLOBAL_TWEAK_INIT_ERROR *DiagnosticType

func StaticGet_ProcessTweaks_NON_GLOBAL_TWEAK_INIT_ERROR() *DiagnosticType {
	clinit_ProcessTweaks()
	return static_ProcessTweaks_NON_GLOBAL_TWEAK_INIT_ERROR
}
func StaticSet_ProcessTweaks_NON_GLOBAL_TWEAK_INIT_ERROR(v *DiagnosticType) {
	clinit_ProcessTweaks()
	static_ProcessTweaks_NON_GLOBAL_TWEAK_INIT_ERROR = v
}

var static_ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING *DiagnosticType

func StaticGet_ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING() *DiagnosticType {
	clinit_ProcessTweaks()
	return static_ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING
}
func StaticSet_ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING(v *DiagnosticType) {
	clinit_ProcessTweaks()
	static_ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING = v
}

var static_ProcessTweaks_INVALID_TWEAK_ID_ERROR *DiagnosticType

func StaticGet_ProcessTweaks_INVALID_TWEAK_ID_ERROR() *DiagnosticType {
	clinit_ProcessTweaks()
	return static_ProcessTweaks_INVALID_TWEAK_ID_ERROR
}
func StaticSet_ProcessTweaks_INVALID_TWEAK_ID_ERROR(v *DiagnosticType) {
	clinit_ProcessTweaks()
	static_ProcessTweaks_INVALID_TWEAK_ID_ERROR = v
}

var static_ProcessTweaks_TWEAK_FUNCTIONS_MAP j2gort.Ext_java_util_Map

func StaticGet_ProcessTweaks_TWEAK_FUNCTIONS_MAP() j2gort.Ext_java_util_Map {
	clinit_ProcessTweaks()
	return static_ProcessTweaks_TWEAK_FUNCTIONS_MAP
}
func StaticSet_ProcessTweaks_TWEAK_FUNCTIONS_MAP(v j2gort.Ext_java_util_Map) {
	clinit_ProcessTweaks()
	static_ProcessTweaks_TWEAK_FUNCTIONS_MAP = v
}

var clinitDone_ProcessTweaks uint32
var clinitState_ProcessTweaks j2gort.ClinitState

func clinit_ProcessTweaks() {
	if atomic.LoadUint32(&clinitDone_ProcessTweaks) != 0 {
		return
	}
	clinitSlow_ProcessTweaks()
}

func clinitSlow_ProcessTweaks() {
	if !j2gort.EnterClinit(&clinitState_ProcessTweaks, &clinitDone_ProcessTweaks) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ProcessTweaks, &clinitDone_ProcessTweaks)
	static_ProcessTweaks_ID_MATCHER = j2gort.Ext_com_google_common_base_CharMatcher_S_InRange(uint16(97), uint16(122)).Or_1_Ext_com_google_common_base_CharMatcher(j2gort.Ext_com_google_common_base_CharMatcher_S_InRange(uint16(65), uint16(90))).Or_1_Ext_com_google_common_base_CharMatcher(j2gort.Ext_com_google_common_base_CharMatcher_S_AnyOf("0123456789_."))
	static_ProcessTweaks_TWEAK_MULTIPLY_REGISTERED_ERROR = DiagnosticType_Error("JSC_TWEAK_MULTIPLY_REGISTERED_ERROR", "Tweak {0} has already been registered.")
	static_ProcessTweaks_NON_LITERAL_TWEAK_ID_ERROR = DiagnosticType_Error("JSC_NON_LITERAL_TWEAK_ID_ERROR", "tweak ID must be a string literal")
	static_ProcessTweaks_INVALID_TWEAK_DEFAULT_VALUE_WARNING = DiagnosticType_Warning("JSC_INVALID_TWEAK_DEFAULT_VALUE_WARNING", "tweak {0} registered with {1} must have a default value that is a literal of type {2}")
	static_ProcessTweaks_NON_GLOBAL_TWEAK_INIT_ERROR = DiagnosticType_Error("JSC_NON_GLOBAL_TWEAK_INIT_ERROR", "tweak declaration {0} must occur in the global scope")
	static_ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING = DiagnosticType_Warning("JSC_TWEAK_WRONG_GETTER_TYPE_WARNING", "tweak getter function {0} used for tweak registered using {1}")
	static_ProcessTweaks_INVALID_TWEAK_ID_ERROR = DiagnosticType_Error("JSC_INVALID_TWEAK_ID_ERROR", "tweak ID contains illegal characters. Only letters, numbers, _ and . are allowed")
	StaticSet_ProcessTweaks_TWEAK_FUNCTIONS_MAP(j2gort.NewExt_java_util_LinkedHashMap())
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_java_util_EnumSet_AllOf(j2gort.ClassLiteral("com.google.javascript.jscomp.ProcessTweaks.TweakFunction", "TweakFunction"))))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			func_ := j2gort.CastPtr[ProcessTweaks_TweakFunction](rawElem_2)
			_ = func_
			j2gort.MapPutStr(StaticGet_ProcessTweaks_TWEAK_FUNCTIONS_MAP(), func_.GetName(), j2gort.PtrToAny(func_))
		}
	}
}

func (this_ *ProcessTweaks) initFields_ProcessTweaks() {
	_ = this_
}

// Java constructor: ProcessTweaks(com.google.javascript.jscomp.AbstractCompiler compiler, boolean stripTweaks)
func NewProcessTweaks(compiler AbstractCompiler, stripTweaks bool) *ProcessTweaks {
	clinit_ProcessTweaks()
	this_ := &ProcessTweaks{}
	InitProcessTweaks(this_, compiler, stripTweaks)
	return this_
}

func InitProcessTweaks(this_ *ProcessTweaks, compiler AbstractCompiler, stripTweaks bool) {
	_ = compiler
	_ = stripTweaks
	this_.initFields_ProcessTweaks()
	this_.compiler = compiler
	this_.stripTweaks = stripTweaks
}

// Java method: ProcessTweaks.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ProcessTweaks) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var result *ProcessTweaks_CollectTweaksResult = this_.priv_ProcessTweaks_collectTweaks_4547(root)
	_ = result
	if this_.stripTweaks {
		this_.priv_ProcessTweaks_stripAllCalls_4546(result.tweakInfos)
	}
}

// Java method: ProcessTweaks.stripAllCalls(java.util.Map<java.lang.String,com.google.javascript.jscomp.ProcessTweaks.TweakInfo> tweakInfos)
func (this_ *ProcessTweaks) priv_ProcessTweaks_stripAllCalls_4546(tweakInfos *j2gort.OrderedMap[string, *ProcessTweaks_TweakInfo]) {
	_ = tweakInfos
	{
		cur_3 := tweakInfos.IterValues()
		for cur_3.HasNext() {

			tweakInfo := cur_3.
				Next()
			_ = tweakInfo
			var isRegistered bool = tweakInfo.IsRegistered()
			_ = isRegistered
			for _, functionCall := range tweakInfo.functionCalls {
				_ = functionCall
				var callNode pkg_rhino.Node = functionCall.callNode
				_ = callNode
				var parent pkg_rhino.Node = callNode.GetParent()
				_ = parent
				if functionCall.tweakFunc.IsGetterFunction() {
					var newValue pkg_rhino.Node = pkg_rhino.Node(0)
					_ = newValue
					if isRegistered {
						newValue = tweakInfo.GetDefaultValueNode().CloneNode_0()
					} else {
						var registerFunction *ProcessTweaks_TweakFunction = functionCall.tweakFunc.registerFunction
						_ = registerFunction
						newValue = registerFunction.CreateDefaultValueNode()
					}
					callNode.ReplaceWith(newValue)
					this_.compiler.ReportChangeToEnclosingScope(parent)
				} else {
					var voidZeroNode pkg_rhino.Node = pkg_rhino.IR_VoidNode(pkg_rhino.IR_Number(float64(int32(0))).Srcref(callNode)).Srcref(callNode)
					_ = voidZeroNode
					callNode.ReplaceWith(voidZeroNode)
					this_.compiler.ReportChangeToEnclosingScope(parent)
				}
			}

		}
	}
}

// Java method: ProcessTweaks.collectTweaks(com.google.javascript.rhino.Node root)
func (this_ *ProcessTweaks) priv_ProcessTweaks_collectTweaks_4547(root pkg_rhino.Node) *ProcessTweaks_CollectTweaksResult {
	_ = root
	var pass *ProcessTweaks_CollectTweaks = NewProcessTweaks_CollectTweaks(this_)
	_ = pass
	NodeTraversal_Traverse(this_.compiler, root, pass)
	var tweakInfos *j2gort.OrderedMap[string, *ProcessTweaks_TweakInfo] = pass.allTweaks
	_ = tweakInfos
	{
		cur_7 := tweakInfos.IterValues()
		for cur_7.HasNext() {

			tweakInfo := cur_7.
				Next()
			_ = tweakInfo
			tweakInfo.EmitAllWarnings()
		}
	}
	return NewProcessTweaks_CollectTweaksResult(tweakInfos)
}

// Java enum: com.google.javascript.jscomp.ProcessTweaks.TweakFunction
type ProcessTweaks_TweakFunction struct {
	j2gort.EnumBase
	name             string
	expectedTypeName string
	validNodeTypeA   pkg_rhino.Token
	validNodeTypeB   pkg_rhino.Token
	registerFunction *ProcessTweaks_TweakFunction
}

func (this_ *ProcessTweaks_TweakFunction) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *ProcessTweaks_TweakFunction) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*ProcessTweaks_TweakFunction) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*ProcessTweaks_TweakFunction) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *ProcessTweaks_TweakFunction) initFields_ProcessTweaks_TweakFunction() {
	_ = this_
	this_.name = j2gort.NullString
	this_.expectedTypeName = j2gort.NullString
}

func initEnumProcessTweaks_TweakFunction_3(this_ *ProcessTweaks_TweakFunction, enumName_ string, enumOrd_ int32, name string, expectedTypeName string, validNodeTypeA pkg_rhino.Token) *ProcessTweaks_TweakFunction {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = name
	_ = expectedTypeName
	_ = validNodeTypeA
	initEnumProcessTweaks_TweakFunction_5(this_, enumName_, enumOrd_, name, expectedTypeName, validNodeTypeA, pkg_rhino.StaticGet_Token_EMPTY(), nil)
	return this_
}

func initEnumProcessTweaks_TweakFunction_4(this_ *ProcessTweaks_TweakFunction, enumName_ string, enumOrd_ int32, name string, expectedTypeName string, validNodeTypeA pkg_rhino.Token, validNodeTypeB pkg_rhino.Token) *ProcessTweaks_TweakFunction {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = name
	_ = expectedTypeName
	_ = validNodeTypeA
	_ = validNodeTypeB
	initEnumProcessTweaks_TweakFunction_5(this_, enumName_, enumOrd_, name, expectedTypeName, validNodeTypeA, validNodeTypeB, nil)
	return this_
}

func initEnumProcessTweaks_TweakFunction_2(this_ *ProcessTweaks_TweakFunction, enumName_ string, enumOrd_ int32, name string, registerFunction *ProcessTweaks_TweakFunction) *ProcessTweaks_TweakFunction {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = name
	_ = registerFunction
	initEnumProcessTweaks_TweakFunction_5(this_, enumName_, enumOrd_, name, j2gort.NullString, pkg_rhino.StaticGet_Token_EMPTY(), pkg_rhino.StaticGet_Token_EMPTY(), registerFunction)
	return this_
}

func initEnumProcessTweaks_TweakFunction_5(this_ *ProcessTweaks_TweakFunction, enumName_ string, enumOrd_ int32, name string, expectedTypeName string, validNodeTypeA pkg_rhino.Token, validNodeTypeB pkg_rhino.Token, registerFunction *ProcessTweaks_TweakFunction) *ProcessTweaks_TweakFunction {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = name
	_ = expectedTypeName
	_ = validNodeTypeA
	_ = validNodeTypeB
	_ = registerFunction
	this_.initFields_ProcessTweaks_TweakFunction()
	this_.name = name
	this_.expectedTypeName = expectedTypeName
	this_.validNodeTypeA = validNodeTypeA
	this_.validNodeTypeB = validNodeTypeB
	this_.registerFunction = registerFunction
	return this_
}

var enum_ProcessTweaks_TweakFunction_REGISTER_BOOLEAN *ProcessTweaks_TweakFunction

func StaticGet_ProcessTweaks_TweakFunction_REGISTER_BOOLEAN() *ProcessTweaks_TweakFunction {
	clinit_ProcessTweaks_TweakFunction()
	return enum_ProcessTweaks_TweakFunction_REGISTER_BOOLEAN
}

var enum_ProcessTweaks_TweakFunction_REGISTER_NUMBER *ProcessTweaks_TweakFunction

func StaticGet_ProcessTweaks_TweakFunction_REGISTER_NUMBER() *ProcessTweaks_TweakFunction {
	clinit_ProcessTweaks_TweakFunction()
	return enum_ProcessTweaks_TweakFunction_REGISTER_NUMBER
}

var enum_ProcessTweaks_TweakFunction_REGISTER_STRING *ProcessTweaks_TweakFunction

func StaticGet_ProcessTweaks_TweakFunction_REGISTER_STRING() *ProcessTweaks_TweakFunction {
	clinit_ProcessTweaks_TweakFunction()
	return enum_ProcessTweaks_TweakFunction_REGISTER_STRING
}

var enum_ProcessTweaks_TweakFunction_GET_BOOLEAN *ProcessTweaks_TweakFunction

func StaticGet_ProcessTweaks_TweakFunction_GET_BOOLEAN() *ProcessTweaks_TweakFunction {
	clinit_ProcessTweaks_TweakFunction()
	return enum_ProcessTweaks_TweakFunction_GET_BOOLEAN
}

var enum_ProcessTweaks_TweakFunction_GET_NUMBER *ProcessTweaks_TweakFunction

func StaticGet_ProcessTweaks_TweakFunction_GET_NUMBER() *ProcessTweaks_TweakFunction {
	clinit_ProcessTweaks_TweakFunction()
	return enum_ProcessTweaks_TweakFunction_GET_NUMBER
}

var enum_ProcessTweaks_TweakFunction_GET_STRING *ProcessTweaks_TweakFunction

func StaticGet_ProcessTweaks_TweakFunction_GET_STRING() *ProcessTweaks_TweakFunction {
	clinit_ProcessTweaks_TweakFunction()
	return enum_ProcessTweaks_TweakFunction_GET_STRING
}

var clinitDone_ProcessTweaks_TweakFunction uint32
var clinitState_ProcessTweaks_TweakFunction j2gort.ClinitState

func clinit_ProcessTweaks_TweakFunction() {
	if atomic.LoadUint32(&clinitDone_ProcessTweaks_TweakFunction) != 0 {
		return
	}
	clinitSlow_ProcessTweaks_TweakFunction()
}

func clinitSlow_ProcessTweaks_TweakFunction() {
	if !j2gort.EnterClinit(&clinitState_ProcessTweaks_TweakFunction, &clinitDone_ProcessTweaks_TweakFunction) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ProcessTweaks_TweakFunction, &clinitDone_ProcessTweaks_TweakFunction)
	enum_ProcessTweaks_TweakFunction_REGISTER_BOOLEAN = &ProcessTweaks_TweakFunction{}
	initEnumProcessTweaks_TweakFunction_4(enum_ProcessTweaks_TweakFunction_REGISTER_BOOLEAN, "REGISTER_BOOLEAN", int32(0), "goog.tweak.registerBoolean", "boolean", pkg_rhino.StaticGet_Token_TRUE(), pkg_rhino.StaticGet_Token_FALSE())
	enum_ProcessTweaks_TweakFunction_REGISTER_NUMBER = &ProcessTweaks_TweakFunction{}
	initEnumProcessTweaks_TweakFunction_3(enum_ProcessTweaks_TweakFunction_REGISTER_NUMBER, "REGISTER_NUMBER", int32(1), "goog.tweak.registerNumber", "number", pkg_rhino.StaticGet_Token_NUMBER())
	enum_ProcessTweaks_TweakFunction_REGISTER_STRING = &ProcessTweaks_TweakFunction{}
	initEnumProcessTweaks_TweakFunction_3(enum_ProcessTweaks_TweakFunction_REGISTER_STRING, "REGISTER_STRING", int32(2), "goog.tweak.registerString", "string", pkg_rhino.StaticGet_Token_STRINGLIT())
	enum_ProcessTweaks_TweakFunction_GET_BOOLEAN = &ProcessTweaks_TweakFunction{}
	initEnumProcessTweaks_TweakFunction_2(enum_ProcessTweaks_TweakFunction_GET_BOOLEAN, "GET_BOOLEAN", int32(3), "goog.tweak.getBoolean", StaticGet_ProcessTweaks_TweakFunction_REGISTER_BOOLEAN())
	enum_ProcessTweaks_TweakFunction_GET_NUMBER = &ProcessTweaks_TweakFunction{}
	initEnumProcessTweaks_TweakFunction_2(enum_ProcessTweaks_TweakFunction_GET_NUMBER, "GET_NUMBER", int32(4), "goog.tweak.getNumber", StaticGet_ProcessTweaks_TweakFunction_REGISTER_NUMBER())
	enum_ProcessTweaks_TweakFunction_GET_STRING = &ProcessTweaks_TweakFunction{}
	initEnumProcessTweaks_TweakFunction_2(enum_ProcessTweaks_TweakFunction_GET_STRING, "GET_STRING", int32(5), "goog.tweak.getString", StaticGet_ProcessTweaks_TweakFunction_REGISTER_STRING())
	enumValues_ProcessTweaks_TweakFunction = []*ProcessTweaks_TweakFunction{enum_ProcessTweaks_TweakFunction_REGISTER_BOOLEAN, enum_ProcessTweaks_TweakFunction_REGISTER_NUMBER, enum_ProcessTweaks_TweakFunction_REGISTER_STRING, enum_ProcessTweaks_TweakFunction_GET_BOOLEAN, enum_ProcessTweaks_TweakFunction_GET_NUMBER, enum_ProcessTweaks_TweakFunction_GET_STRING}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ProcessTweaks.TweakFunction", []any{enum_ProcessTweaks_TweakFunction_REGISTER_BOOLEAN, enum_ProcessTweaks_TweakFunction_REGISTER_NUMBER, enum_ProcessTweaks_TweakFunction_REGISTER_STRING, enum_ProcessTweaks_TweakFunction_GET_BOOLEAN, enum_ProcessTweaks_TweakFunction_GET_NUMBER, enum_ProcessTweaks_TweakFunction_GET_STRING})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ProcessTweaks$TweakFunction", []any{enum_ProcessTweaks_TweakFunction_REGISTER_BOOLEAN, enum_ProcessTweaks_TweakFunction_REGISTER_NUMBER, enum_ProcessTweaks_TweakFunction_REGISTER_STRING, enum_ProcessTweaks_TweakFunction_GET_BOOLEAN, enum_ProcessTweaks_TweakFunction_GET_NUMBER, enum_ProcessTweaks_TweakFunction_GET_STRING})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.ProcessTweaks.TweakFunction", "com.google.javascript.jscomp.ProcessTweaks$TweakFunction", func() {
		clinit_ProcessTweaks_TweakFunction()
	})
}

var enumValues_ProcessTweaks_TweakFunction []*ProcessTweaks_TweakFunction

func ProcessTweaks_TweakFunction_Values() []*ProcessTweaks_TweakFunction {
	clinit_ProcessTweaks_TweakFunction()
	return j2gort.CloneSlice(enumValues_ProcessTweaks_TweakFunction)
}

func ProcessTweaks_TweakFunction_ValueOf(name string) *ProcessTweaks_TweakFunction {
	clinit_ProcessTweaks_TweakFunction()
	for _, c := range enumValues_ProcessTweaks_TweakFunction {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ProcessTweaks.TweakFunction." + name))
}

// Java method: TweakFunction.isValidNodeType(com.google.javascript.rhino.Token type)
func (this_ *ProcessTweaks_TweakFunction) IsValidNodeType(type_ pkg_rhino.Token) bool {
	_ = type_
	return ((type_ == this_.validNodeTypeA) || (type_ == this_.validNodeTypeB))
}

// Java method: TweakFunction.isCorrectRegisterFunction(com.google.javascript.jscomp.ProcessTweaks.TweakFunction registerFunction)
func (this_ *ProcessTweaks_TweakFunction) IsCorrectRegisterFunction(registerFunction *ProcessTweaks_TweakFunction) bool {
	_ = registerFunction
	j2gort.CheckNotNullPtr(registerFunction)
	return (this_.registerFunction == registerFunction)
}

// Java method: TweakFunction.isGetterFunction()
func (this_ *ProcessTweaks_TweakFunction) IsGetterFunction() bool {
	return (this_.registerFunction != nil)
}

// Java method: TweakFunction.getName()
func (this_ *ProcessTweaks_TweakFunction) GetName() string {
	return this_.name
}

// Java method: TweakFunction.getExpectedTypeName()
func (this_ *ProcessTweaks_TweakFunction) GetExpectedTypeName() string {
	return this_.expectedTypeName
}

// Java method: TweakFunction.createDefaultValueNode()
func (this_ *ProcessTweaks_TweakFunction) CreateDefaultValueNode() pkg_rhino.Node {
	switch this_.Ordinal() {
	case int32(0) /* TweakFunction.REGISTER_BOOLEAN */ :
		return pkg_rhino.IR_FalseNode()
	case int32(1) /* TweakFunction.REGISTER_NUMBER */ :
		return pkg_rhino.IR_Number(float64(int32(0)))
	case int32(2) /* TweakFunction.REGISTER_STRING */ :
		return pkg_rhino.IR_String("")
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
	}
	panic("unreachable switch expression")
}

func (this_ *ProcessTweaks_TweakFunction) IsJavaNil_() bool { return this_ == nil }
func (*ProcessTweaks_TweakFunction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessTweaks$TweakFunction", "TweakFunction"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessTweaks.TweakFunction", "com.google.javascript.jscomp.ProcessTweaks$TweakFunction", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessTweaks_TweakFunction](v_)
	})
}

// Java class: com.google.javascript.jscomp.ProcessTweaks.CollectTweaksResult
type ProcessTweaks_CollectTweaksResult struct {
	tweakInfos *j2gort.OrderedMap[string, *ProcessTweaks_TweakInfo] `gson:"tweakInfos"`
	idHash_    atomic.Int32
}

func (this_ *ProcessTweaks_CollectTweaksResult) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessTweaks_CollectTweaksResult) IsJavaNil_() bool { return this_ == nil }
func (*ProcessTweaks_CollectTweaksResult) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessTweaks$CollectTweaksResult", "CollectTweaksResult"
}
func (this_ *ProcessTweaks_CollectTweaksResult) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessTweaks_CollectTweaksResult) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessTweaks.CollectTweaksResult", "com.google.javascript.jscomp.ProcessTweaks$CollectTweaksResult", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessTweaks_CollectTweaksResult](v_)
	})
}

func (this_ *ProcessTweaks_CollectTweaksResult) initFields_ProcessTweaks_CollectTweaksResult() {
	_ = this_
}

// Java constructor: CollectTweaksResult(java.util.Map<java.lang.String,com.google.javascript.jscomp.ProcessTweaks.TweakInfo> tweakInfos)
func NewProcessTweaks_CollectTweaksResult(tweakInfos *j2gort.OrderedMap[string, *ProcessTweaks_TweakInfo]) *ProcessTweaks_CollectTweaksResult {
	this_ := &ProcessTweaks_CollectTweaksResult{}
	InitProcessTweaks_CollectTweaksResult(this_, tweakInfos)
	return this_
}

func InitProcessTweaks_CollectTweaksResult(this_ *ProcessTweaks_CollectTweaksResult, tweakInfos *j2gort.OrderedMap[string, *ProcessTweaks_TweakInfo]) {
	_ = tweakInfos
	this_.initFields_ProcessTweaks_CollectTweaksResult()
	this_.tweakInfos = tweakInfos
}

// Java class: com.google.javascript.jscomp.ProcessTweaks.CollectTweaks
type ProcessTweaks_CollectTweaks struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_    *ProcessTweaks
	allTweaks *j2gort.OrderedMap[string, *ProcessTweaks_TweakInfo] `gson:"allTweaks"`
}

func (this_ *ProcessTweaks_CollectTweaks) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessTweaks_CollectTweaks) IsJavaNil_() bool { return this_ == nil }
func (*ProcessTweaks_CollectTweaks) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessTweaks$CollectTweaks", "CollectTweaks"
}
func (this_ *ProcessTweaks_CollectTweaks) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ProcessTweaks_CollectTweaks) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessTweaks.CollectTweaks", "com.google.javascript.jscomp.ProcessTweaks$CollectTweaks", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessTweaks_CollectTweaks](v_)
	})
}
func (this_ *ProcessTweaks_CollectTweaks) SetOuter_(o *ProcessTweaks)                           { this_.outer_ = o }
func (this_ *ProcessTweaks_CollectTweaks) GetOuter_() *ProcessTweaks                            { return this_.outer_ }
func (*ProcessTweaks_CollectTweaks) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ProcessTweaks_CollectTweaks) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *ProcessTweaks_CollectTweaks) initFields_ProcessTweaks_CollectTweaks() {
	_ = this_
	this_.allTweaks = &j2gort.OrderedMap[string, *ProcessTweaks_TweakInfo]{

		// Java constructor: CollectTweaks()
	}
}

func NewProcessTweaks_CollectTweaks(outer_ *ProcessTweaks) *ProcessTweaks_CollectTweaks {
	this_ := &ProcessTweaks_CollectTweaks{}
	this_.outer_ = outer_
	InitProcessTweaks_CollectTweaks(this_)
	return this_
}

func InitProcessTweaks_CollectTweaks(this_ *ProcessTweaks_CollectTweaks) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ProcessTweaks_CollectTweaks()
}

// Java method: CollectTweaks.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessTweaks_CollectTweaks) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsCall()) {
		return
	}
	var callName string = n.GetFirstChild().GetQualifiedName()
	_ = callName
	var tweakFunc *ProcessTweaks_TweakFunction = j2gort.CastPtr[ProcessTweaks_TweakFunction](j2gort.MapGetStr(StaticGet_ProcessTweaks_TWEAK_FUNCTIONS_MAP(), callName))
	_ = tweakFunc
	if tweakFunc == nil {
		return
	}
	var tweakIdNode pkg_rhino.Node = n.GetSecondChild()
	_ = tweakIdNode
	if !(tweakIdNode.IsStringLit()) {
		this_.outer_.compiler.Report_1(JSError_Make_3(tweakIdNode, StaticGet_ProcessTweaks_NON_LITERAL_TWEAK_ID_ERROR(), []string{}))
		return
	}
	var tweakId string = tweakIdNode.GetString()
	_ = tweakId
	var tweakInfo *ProcessTweaks_TweakInfo = j2gort.OrderedMapComputeIfAbsent(this_.allTweaks, tweakId, nil, func(mrp0_ string) *ProcessTweaks_TweakInfo {
		return NewProcessTweaks_TweakInfo(this_.outer_, mrp0_)
	})
	_ = tweakInfo
	switch j2gort.CheckNotNullPtr(tweakFunc).Ordinal() {
	case int32(0) /* TweakFunction.REGISTER_BOOLEAN */, int32(1) /* TweakFunction.REGISTER_NUMBER */, int32(2) /* TweakFunction.REGISTER_STRING */ :
		if !(StaticGet_ProcessTweaks_ID_MATCHER().MatchesAllOf(j2gort.BoxString(tweakId))) {
			this_.outer_.compiler.Report_1(JSError_Make_3(tweakIdNode, StaticGet_ProcessTweaks_INVALID_TWEAK_ID_ERROR(), []string{}))
		}
		if !(t.InGlobalHoistScope()) {
			this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ProcessTweaks_NON_GLOBAL_TWEAK_INIT_ERROR(), []string{tweakId}))
			break
		}
		if tweakInfo.IsRegistered() {
			this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ProcessTweaks_TWEAK_MULTIPLY_REGISTERED_ERROR(), []string{tweakId}))
			break
		}
		var tweakDefaultValueNode pkg_rhino.Node = tweakIdNode.GetNext().GetNext()
		_ = tweakDefaultValueNode
		tweakInfo.AddRegisterCall(tweakFunc, n, tweakDefaultValueNode)
	case int32(3) /* TweakFunction.GET_BOOLEAN */, int32(4) /* TweakFunction.GET_NUMBER */, int32(5) /* TweakFunction.GET_STRING */ :
		tweakInfo.AddGetterCall(tweakFunc, n)
	}
}

// Java class: com.google.javascript.jscomp.ProcessTweaks.TweakFunctionCall
type ProcessTweaks_TweakFunctionCall struct {
	tweakFunc *ProcessTweaks_TweakFunction `gson:"tweakFunc"`
	callNode  pkg_rhino.Node               `gson:"callNode"`
	valueNode pkg_rhino.Node               `gson:"valueNode"`
	idHash_   atomic.Int32
}

func (this_ *ProcessTweaks_TweakFunctionCall) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessTweaks_TweakFunctionCall) IsJavaNil_() bool { return this_ == nil }
func (*ProcessTweaks_TweakFunctionCall) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessTweaks$TweakFunctionCall", "TweakFunctionCall"
}
func (this_ *ProcessTweaks_TweakFunctionCall) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessTweaks_TweakFunctionCall) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessTweaks.TweakFunctionCall", "com.google.javascript.jscomp.ProcessTweaks$TweakFunctionCall", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessTweaks_TweakFunctionCall](v_)
	})
}

func (this_ *ProcessTweaks_TweakFunctionCall) initFields_ProcessTweaks_TweakFunctionCall() {
	_ = this_
}

// Java constructor: TweakFunctionCall(com.google.javascript.jscomp.ProcessTweaks.TweakFunction tweakFunc, com.google.javascript.rhino.Node callNode)
func NewProcessTweaks_TweakFunctionCall_2(tweakFunc *ProcessTweaks_TweakFunction, callNode pkg_rhino.Node) *ProcessTweaks_TweakFunctionCall {
	this_ := &ProcessTweaks_TweakFunctionCall{}
	InitProcessTweaks_TweakFunctionCall_2(this_, tweakFunc, callNode)
	return this_
}

func InitProcessTweaks_TweakFunctionCall_2(this_ *ProcessTweaks_TweakFunctionCall, tweakFunc *ProcessTweaks_TweakFunction, callNode pkg_rhino.Node) {
	_ = tweakFunc
	_ = callNode
	InitProcessTweaks_TweakFunctionCall_3(this_, tweakFunc, callNode, pkg_rhino.Node(0))
}

// Java constructor: TweakFunctionCall(com.google.javascript.jscomp.ProcessTweaks.TweakFunction tweakFunc, com.google.javascript.rhino.Node callNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node valueNode)
func NewProcessTweaks_TweakFunctionCall_3(tweakFunc *ProcessTweaks_TweakFunction, callNode pkg_rhino.Node, valueNode pkg_rhino.Node) *ProcessTweaks_TweakFunctionCall {
	this_ := &ProcessTweaks_TweakFunctionCall{}
	InitProcessTweaks_TweakFunctionCall_3(this_, tweakFunc, callNode, valueNode)
	return this_
}

func InitProcessTweaks_TweakFunctionCall_3(this_ *ProcessTweaks_TweakFunctionCall, tweakFunc *ProcessTweaks_TweakFunction, callNode pkg_rhino.Node, valueNode pkg_rhino.Node) {
	_ = tweakFunc
	_ = callNode
	_ = valueNode
	this_.initFields_ProcessTweaks_TweakFunctionCall()
	this_.callNode = callNode
	this_.tweakFunc = tweakFunc
	this_.valueNode = valueNode
}

// Java class: com.google.javascript.jscomp.ProcessTweaks.TweakInfo
type ProcessTweaks_TweakInfo struct {
	outer_           *ProcessTweaks
	tweakId          string                             `gson:"tweakId"`
	functionCalls    []*ProcessTweaks_TweakFunctionCall `gson:"functionCalls"`
	registerCall     *ProcessTweaks_TweakFunctionCall   `gson:"registerCall"`
	defaultValueNode pkg_rhino.Node                     `gson:"defaultValueNode"`
	idHash_          atomic.Int32
}

func (this_ *ProcessTweaks_TweakInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessTweaks_TweakInfo) IsJavaNil_() bool { return this_ == nil }
func (*ProcessTweaks_TweakInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessTweaks$TweakInfo", "TweakInfo"
}
func (this_ *ProcessTweaks_TweakInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessTweaks_TweakInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessTweaks.TweakInfo", "com.google.javascript.jscomp.ProcessTweaks$TweakInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessTweaks_TweakInfo](v_)
	})
}
func (this_ *ProcessTweaks_TweakInfo) SetOuter_(o *ProcessTweaks) { this_.outer_ = o }
func (this_ *ProcessTweaks_TweakInfo) GetOuter_() *ProcessTweaks  { return this_.outer_ }

func (this_ *ProcessTweaks_TweakInfo) initFields_ProcessTweaks_TweakInfo() {
	_ = this_
	this_.tweakId = j2gort.NullString
}

// Java constructor: TweakInfo(java.lang.String tweakId)
func NewProcessTweaks_TweakInfo(outer_ *ProcessTweaks, tweakId string) *ProcessTweaks_TweakInfo {
	this_ := &ProcessTweaks_TweakInfo{}
	this_.outer_ = outer_
	InitProcessTweaks_TweakInfo(this_, tweakId)
	return this_
}

func InitProcessTweaks_TweakInfo(this_ *ProcessTweaks_TweakInfo, tweakId string) {
	_ = tweakId
	this_.initFields_ProcessTweaks_TweakInfo()
	this_.tweakId = tweakId
	this_.functionCalls = []*ProcessTweaks_TweakFunctionCall{}
}

// Java method: TweakInfo.emitAllWarnings()
func (this_ *ProcessTweaks_TweakInfo) EmitAllWarnings() {
	if this_.IsRegistered() {
		this_.EmitAllTypeWarnings()
	}
}

// Java method: TweakInfo.emitAllTypeWarnings()
func (this_ *ProcessTweaks_TweakInfo) EmitAllTypeWarnings() {
	for _, call := range this_.functionCalls {
		_ = call
		var valueNode pkg_rhino.Node = call.valueNode
		_ = valueNode
		var tweakFunc *ProcessTweaks_TweakFunction = call.tweakFunc
		_ = tweakFunc
		var registerFunc *ProcessTweaks_TweakFunction = this_.registerCall.tweakFunc
		_ = registerFunc
		if valueNode != 0 {
			if !(registerFunc.IsValidNodeType(valueNode.GetToken())) {
				this_.outer_.compiler.Report_1(JSError_Make_3(valueNode, StaticGet_ProcessTweaks_INVALID_TWEAK_DEFAULT_VALUE_WARNING(), []string{this_.tweakId, registerFunc.GetName(), registerFunc.GetExpectedTypeName()}))
			}
		} else if tweakFunc.IsGetterFunction() {
			if !(tweakFunc.IsCorrectRegisterFunction(registerFunc)) {
				this_.outer_.compiler.Report_1(JSError_Make_3(call.callNode, StaticGet_ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING(), []string{tweakFunc.GetName(), registerFunc.GetName()}))
			}
		}
	}

}

// Java method: TweakInfo.addRegisterCall(com.google.javascript.jscomp.ProcessTweaks.TweakFunction tweakFunc, com.google.javascript.rhino.Node callNode, com.google.javascript.rhino.Node defaultValueNode)
func (this_ *ProcessTweaks_TweakInfo) AddRegisterCall(tweakFunc *ProcessTweaks_TweakFunction, callNode pkg_rhino.Node, defaultValueNode pkg_rhino.Node) {
	_ = tweakFunc
	_ = callNode
	_ = defaultValueNode
	this_.registerCall = NewProcessTweaks_TweakFunctionCall_3(tweakFunc, callNode, defaultValueNode)
	this_.functionCalls = append(this_.functionCalls, this_.registerCall)
}

// Java method: TweakInfo.addGetterCall(com.google.javascript.jscomp.ProcessTweaks.TweakFunction tweakFunc, com.google.javascript.rhino.Node callNode)
func (this_ *ProcessTweaks_TweakInfo) AddGetterCall(tweakFunc *ProcessTweaks_TweakFunction, callNode pkg_rhino.Node) {
	_ = tweakFunc
	_ = callNode
	this_.functionCalls = append(this_.functionCalls, NewProcessTweaks_TweakFunctionCall_2(tweakFunc, callNode))
}

// Java method: TweakInfo.isRegistered()
func (this_ *ProcessTweaks_TweakInfo) IsRegistered() bool {
	return (this_.registerCall != nil)
}

// Java method: TweakInfo.getDefaultValueNode()
func (this_ *ProcessTweaks_TweakInfo) GetDefaultValueNode() pkg_rhino.Node {
	if !(this_.IsRegistered()) {
		j2gort.CheckState(false)
	}
	if this_.defaultValueNode != 0 {
		return this_.defaultValueNode
	}
	if this_.registerCall.valueNode != 0 {
		return this_.registerCall.valueNode
	}
	return this_.registerCall.tweakFunc.CreateDefaultValueNode()
}
func (*ProcessTweaks_CollectTweaks) AlwaysShouldTraverse_()      {}
func (*ProcessTweaks_CollectTweaks) CanReleaseSyntacticScopes_() {}
