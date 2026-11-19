/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/AstFactory.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_ast_factory_1 j2gort.Ext_com_google_common_base_Supplier

func init() {
	lambdaSingleton_ast_factory_1 = j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
		return j2gort.PtrToAny(pkg_colors.Color_CreateUnion_1_Color([]*pkg_colors.Color{pkg_colors.StaticGet_StandardColors_BIGINT(), pkg_colors.StaticGet_StandardColors_STRING(), pkg_colors.StaticGet_StandardColors_NUMBER()}))
	})
}

// Java class: com.google.javascript.jscomp.AstFactory
type AstFactory struct {
	colorRegistry       *pkg_colors.ColorRegistry        `gson:"colorRegistry"`
	registry            *pkg_rhino.JSTypeRegistry        `gson:"registry"`
	unknownType         pkg_rhino.JSType                 `gson:"unknownType"`
	lifeCycleStage      *AbstractCompiler_LifeCycleStage `gson:"lifeCycleStage"`
	runtimeJsLibManager *RuntimeJsLibManager             `gson:"runtimeJsLibManager"`
	idHash_             atomic.Int32
	typeMode            AstFactory_TypeMode `gson:"typeMode"`
}

func (this_ *AstFactory) BindSelf_(s any) {
	_ = s
}
func (this_ *AstFactory) IsJavaNil_() bool { return this_ == nil }
func (*AstFactory) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AstFactory", "AstFactory"
}
func (this_ *AstFactory) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AstFactory) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AstFactory", "com.google.javascript.jscomp.AstFactory", func(v_ any) bool {
		return j2gort.InstanceOf[*AstFactory](v_)
	})
}

var static_AstFactory_DOT_SPLITTER *j2gort.Ext_com_google_common_base_Splitter

func StaticGet_AstFactory_DOT_SPLITTER() *j2gort.Ext_com_google_common_base_Splitter {
	clinit_AstFactory()
	return static_AstFactory_DOT_SPLITTER
}
func StaticSet_AstFactory_DOT_SPLITTER(v *j2gort.Ext_com_google_common_base_Splitter) {
	clinit_AstFactory()
	static_AstFactory_DOT_SPLITTER = v
}

var static_AstFactory_bigintNumberStringColor j2gort.Ext_java_util_function_Supplier

func StaticGet_AstFactory_BigintNumberStringColor() j2gort.Ext_java_util_function_Supplier {
	clinit_AstFactory()
	return static_AstFactory_bigintNumberStringColor
}
func StaticSet_AstFactory_BigintNumberStringColor(v j2gort.Ext_java_util_function_Supplier) {
	clinit_AstFactory()
	static_AstFactory_bigintNumberStringColor = v
}

var clinitDone_AstFactory uint32
var clinitState_AstFactory j2gort.ClinitState

func clinit_AstFactory() {
	if atomic.LoadUint32(&clinitDone_AstFactory) != 0 {
		return
	}
	clinitSlow_AstFactory()
}

func clinitSlow_AstFactory() {
	if !j2gort.EnterClinit(&clinitState_AstFactory, &clinitDone_AstFactory) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AstFactory, &clinitDone_AstFactory)
	static_AstFactory_DOT_SPLITTER = j2gort.Ext_com_google_common_base_Splitter_On_1_string(".")
	static_AstFactory_bigintNumberStringColor = j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](j2gort.Ext_com_google_common_base_Suppliers_Memoize(lambdaSingleton_ast_factory_1))
}

func (this_ *AstFactory) initFields_AstFactory() {
	_ = this_
}

// Java constructor: AstFactory(com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage lifeCycleStage, com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.jscomp.js.RuntimeJsLibManager runtimeJsLibManager)
func NewAstFactory_3_AbstractCompiler_LifeCycleStage_JSTypeRegistry_RuntimeJsLibManager(lifeCycleStage *AbstractCompiler_LifeCycleStage, registry *pkg_rhino.JSTypeRegistry, runtimeJsLibManager *RuntimeJsLibManager) *AstFactory {
	clinit_AstFactory()
	this_ := &AstFactory{}
	InitAstFactory_3_AbstractCompiler_LifeCycleStage_JSTypeRegistry_RuntimeJsLibManager(this_, lifeCycleStage, registry, runtimeJsLibManager)
	return this_
}

func InitAstFactory_3_AbstractCompiler_LifeCycleStage_JSTypeRegistry_RuntimeJsLibManager(this_ *AstFactory, lifeCycleStage *AbstractCompiler_LifeCycleStage, registry *pkg_rhino.JSTypeRegistry, runtimeJsLibManager *RuntimeJsLibManager) {
	_ = lifeCycleStage
	_ = registry
	_ = runtimeJsLibManager
	this_.initFields_AstFactory()
	this_.lifeCycleStage = lifeCycleStage
	this_.registry = registry
	this_.runtimeJsLibManager = runtimeJsLibManager
	this_.colorRegistry = nil
	this_.unknownType = this_.priv_AstFactory_getNativeType_742(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	this_.typeMode = StaticGet_AstFactory_TypeMode_JSTYPE()
}

// Java constructor: AstFactory(com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage lifeCycleStage, com.google.javascript.jscomp.js.RuntimeJsLibManager runtimeJsLibManager)
func NewAstFactory_2(lifeCycleStage *AbstractCompiler_LifeCycleStage, runtimeJsLibManager *RuntimeJsLibManager) *AstFactory {
	clinit_AstFactory()
	this_ := &AstFactory{}
	InitAstFactory_2(this_, lifeCycleStage, runtimeJsLibManager)
	return this_
}

func InitAstFactory_2(this_ *AstFactory, lifeCycleStage *AbstractCompiler_LifeCycleStage, runtimeJsLibManager *RuntimeJsLibManager) {
	_ = lifeCycleStage
	_ = runtimeJsLibManager
	this_.initFields_AstFactory()
	this_.lifeCycleStage = lifeCycleStage
	this_.registry = nil
	this_.runtimeJsLibManager = runtimeJsLibManager
	this_.colorRegistry = nil
	this_.unknownType = nil
	this_.typeMode = StaticGet_AstFactory_TypeMode_NONE()
}

// Java constructor: AstFactory(com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage lifeCycleStage, com.google.javascript.jscomp.colors.ColorRegistry colorRegistry, com.google.javascript.jscomp.js.RuntimeJsLibManager runtimeJsLibManager)
func NewAstFactory_3_AbstractCompiler_LifeCycleStage_ColorRegistry_RuntimeJsLibManager(lifeCycleStage *AbstractCompiler_LifeCycleStage, colorRegistry *pkg_colors.ColorRegistry, runtimeJsLibManager *RuntimeJsLibManager) *AstFactory {
	clinit_AstFactory()
	this_ := &AstFactory{}
	InitAstFactory_3_AbstractCompiler_LifeCycleStage_ColorRegistry_RuntimeJsLibManager(this_, lifeCycleStage, colorRegistry, runtimeJsLibManager)
	return this_
}

func InitAstFactory_3_AbstractCompiler_LifeCycleStage_ColorRegistry_RuntimeJsLibManager(this_ *AstFactory, lifeCycleStage *AbstractCompiler_LifeCycleStage, colorRegistry *pkg_colors.ColorRegistry, runtimeJsLibManager *RuntimeJsLibManager) {
	_ = lifeCycleStage
	_ = colorRegistry
	_ = runtimeJsLibManager
	this_.initFields_AstFactory()
	this_.lifeCycleStage = lifeCycleStage
	this_.registry = nil
	this_.runtimeJsLibManager = runtimeJsLibManager
	this_.colorRegistry = colorRegistry
	this_.unknownType = nil
	this_.typeMode = StaticGet_AstFactory_TypeMode_COLOR()
}

// Java static method: AstFactory.createFactoryWithoutTypes(com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage lifeCycleStage, com.google.javascript.jscomp.js.RuntimeJsLibManager runtimeJsLibManager)
func AstFactory_CreateFactoryWithoutTypes(lifeCycleStage *AbstractCompiler_LifeCycleStage, runtimeJsLibManager *RuntimeJsLibManager) *AstFactory {
	clinit_AstFactory()
	_ = lifeCycleStage
	_ = runtimeJsLibManager
	return NewAstFactory_2(lifeCycleStage, runtimeJsLibManager)
}

// Java static method: AstFactory.createFactoryWithTypes(com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage lifeCycleStage, com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.jscomp.js.RuntimeJsLibManager runtimeJsLibManager)
func AstFactory_CreateFactoryWithTypes(lifeCycleStage *AbstractCompiler_LifeCycleStage, registry *pkg_rhino.JSTypeRegistry, runtimeJsLibManager *RuntimeJsLibManager) *AstFactory {
	clinit_AstFactory()
	_ = lifeCycleStage
	_ = registry
	_ = runtimeJsLibManager
	return NewAstFactory_3_AbstractCompiler_LifeCycleStage_JSTypeRegistry_RuntimeJsLibManager(lifeCycleStage, registry, runtimeJsLibManager)
}

// Java static method: AstFactory.createFactoryWithColors(com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage lifeCycleStage, com.google.javascript.jscomp.colors.ColorRegistry colorRegistry, com.google.javascript.jscomp.js.RuntimeJsLibManager runtimeJsLibManager)
func AstFactory_CreateFactoryWithColors(lifeCycleStage *AbstractCompiler_LifeCycleStage, colorRegistry *pkg_colors.ColorRegistry, runtimeJsLibManager *RuntimeJsLibManager) *AstFactory {
	clinit_AstFactory()
	_ = lifeCycleStage
	_ = colorRegistry
	_ = runtimeJsLibManager
	return NewAstFactory_3_AbstractCompiler_LifeCycleStage_ColorRegistry_RuntimeJsLibManager(lifeCycleStage, colorRegistry, runtimeJsLibManager)
}

// Java method: AstFactory.isAddingTypes()
func (this_ *AstFactory) IsAddingTypes() bool {
	return func() bool {
		re_1 := StaticGet_AstFactory_TypeMode_JSTYPE()
		ae_2 := this_.typeMode
		if re_1 == 0 {
			j2gort.ThrowNPE()
		}
		return re_1 == ae_2
	}()
}

// Java method: AstFactory.isAddingColors()
func (this_ *AstFactory) IsAddingColors() bool {
	return func() bool {
		re_3 := StaticGet_AstFactory_TypeMode_COLOR()
		ae_4 := this_.typeMode
		if re_3 == 0 {
			j2gort.ThrowNPE()
		}
		return re_3 == ae_4
	}()
}

// Java method: AstFactory.assertNotAddingColors()
func (this_ *AstFactory) priv_AstFactory_assertNotAddingColors_620() {
	if !(!(this_.IsAddingColors())) {
		j2gort.CheckStateStr(false, "method not supported for colors")
	}
}

// Java method: AstFactory.exprResult(com.google.javascript.rhino.Node expr)
func (this_ *AstFactory) ExprResult(expr pkg_rhino.Node) pkg_rhino.Node {
	_ = expr
	return pkg_rhino.IR_ExprResult(expr).Srcref(expr)
}

// Java method: AstFactory.createEmpty()
func (this_ *AstFactory) CreateEmpty() pkg_rhino.Node {
	return pkg_rhino.IR_Empty()
}

// Java method: AstFactory.createBlock(com.google.javascript.rhino.Node[] statements)
func (this_ *AstFactory) CreateBlock(statements []pkg_rhino.Node) pkg_rhino.Node {
	_ = statements
	return pkg_rhino.IR_Block_1_Node_2(statements)
}

// Java method: AstFactory.createIf(com.google.javascript.rhino.Node cond, com.google.javascript.rhino.Node then)
func (this_ *AstFactory) CreateIf_2(cond pkg_rhino.Node, then pkg_rhino.Node) pkg_rhino.Node {
	_ = cond
	_ = then
	return pkg_rhino.IR_IfNode_2(cond, then)
}

// Java method: AstFactory.createIf(com.google.javascript.rhino.Node cond, com.google.javascript.rhino.Node then, com.google.javascript.rhino.Node elseNode)
func (this_ *AstFactory) CreateIf_3(cond pkg_rhino.Node, then pkg_rhino.Node, elseNode pkg_rhino.Node) pkg_rhino.Node {
	_ = cond
	_ = then
	_ = elseNode
	return pkg_rhino.IR_IfNode_3(cond, then, elseNode)
}

// Java method: AstFactory.createFor(com.google.javascript.rhino.Node init, com.google.javascript.rhino.Node cond, com.google.javascript.rhino.Node incr, com.google.javascript.rhino.Node body)
func (this_ *AstFactory) CreateFor(init pkg_rhino.Node, cond pkg_rhino.Node, incr pkg_rhino.Node, body pkg_rhino.Node) pkg_rhino.Node {
	_ = init
	_ = cond
	_ = incr
	_ = body
	return pkg_rhino.IR_ForNode(init, cond, incr, body)
}

// Java method: AstFactory.createBreak()
func (this_ *AstFactory) CreateBreak() pkg_rhino.Node {
	return pkg_rhino.IR_BreakNode_0()
}

// Java method: AstFactory.createLabel(com.google.javascript.rhino.Node label, com.google.javascript.rhino.Node stmt)
func (this_ *AstFactory) CreateLabel(label pkg_rhino.Node, stmt pkg_rhino.Node) pkg_rhino.Node {
	_ = label
	_ = stmt
	return pkg_rhino.IR_Label(label, stmt)
}

// Java method: AstFactory.createCatch(com.google.javascript.rhino.Node error, com.google.javascript.rhino.Node block)
func (this_ *AstFactory) CreateCatch(error_ pkg_rhino.Node, block pkg_rhino.Node) pkg_rhino.Node {
	_ = error_
	_ = block
	return pkg_rhino.IR_CatchNode(error_, block)
}

// Java method: AstFactory.createTryFinally(com.google.javascript.rhino.Node tryBlock, com.google.javascript.rhino.Node finallyBlock)
func (this_ *AstFactory) CreateTryFinally(tryBlock pkg_rhino.Node, finallyBlock pkg_rhino.Node) pkg_rhino.Node {
	_ = tryBlock
	_ = finallyBlock
	return pkg_rhino.IR_TryFinally(tryBlock, finallyBlock)
}

// Java method: AstFactory.createTryCatchFinally(com.google.javascript.rhino.Node tryBlock, com.google.javascript.rhino.Node catchNode, com.google.javascript.rhino.Node finallyBlock)
func (this_ *AstFactory) CreateTryCatchFinally(tryBlock pkg_rhino.Node, catchNode pkg_rhino.Node, finallyBlock pkg_rhino.Node) pkg_rhino.Node {
	_ = tryBlock
	_ = catchNode
	_ = finallyBlock
	if !(tryBlock.IsBlock()) {
		j2gort.CheckState(false)
	}
	if !(catchNode.IsCatch()) {
		j2gort.CheckState(false)
	}
	if !(finallyBlock.IsBlock()) {
		j2gort.CheckState(false)
	}
	return pkg_rhino.IR_TryCatchFinally(tryBlock, catchNode, finallyBlock)
}

// Java method: AstFactory.createThrow(com.google.javascript.rhino.Node expr)
func (this_ *AstFactory) CreateThrow(expr pkg_rhino.Node) pkg_rhino.Node {
	_ = expr
	return pkg_rhino.IR_ThrowNode(expr)
}

// Java method: AstFactory.createReturn(com.google.javascript.rhino.Node value)
func (this_ *AstFactory) CreateReturn(value pkg_rhino.Node) pkg_rhino.Node {
	_ = value
	return pkg_rhino.IR_ReturnNode_1(value)
}

// Java method: AstFactory.createYield(com.google.javascript.jscomp.AstFactory.Type type, com.google.javascript.rhino.Node value)
func (this_ *AstFactory) CreateYield(type_ AstFactory_Type, value pkg_rhino.Node) pkg_rhino.Node {
	_ = type_
	_ = value
	var result pkg_rhino.Node = pkg_rhino.IR_YieldNode(value)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	return result
}

// Java method: AstFactory.createAwait(com.google.javascript.jscomp.AstFactory.Type type, com.google.javascript.rhino.Node value)
func (this_ *AstFactory) CreateAwait(type_ AstFactory_Type, value pkg_rhino.Node) pkg_rhino.Node {
	_ = type_
	_ = value
	var result pkg_rhino.Node = pkg_rhino.IR_Await(value)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	return result
}

// Java method: AstFactory.createString(java.lang.String value)
func (this_ *AstFactory) CreateString(value string) pkg_rhino.Node {
	_ = value
	var result pkg_rhino.Node = pkg_rhino.IR_String(value)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE(), pkg_colors.StaticGet_StandardColors_STRING()), result)
	return result
}

// Java method: AstFactory.createNumber(double value)
func (this_ *AstFactory) CreateNumber(value float64) pkg_rhino.Node {
	_ = value
	var result pkg_rhino.Node = pkg_rhino.IR_Number(value)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE(), pkg_colors.StaticGet_StandardColors_NUMBER()), result)
	return result
}

// Java method: AstFactory.createBoolean(boolean value)
func (this_ *AstFactory) CreateBoolean(value bool) pkg_rhino.Node {
	_ = value
	var result pkg_rhino.Node
	_ = result
	if value {
		result = pkg_rhino.IR_TrueNode()
	} else {
		result = pkg_rhino.IR_FalseNode()
	}
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE(), pkg_colors.StaticGet_StandardColors_BOOLEAN()), result)
	return result
}

// Java method: AstFactory.createNull()
func (this_ *AstFactory) CreateNull() pkg_rhino.Node {
	var result pkg_rhino.Node = pkg_rhino.IR_NullNode()
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_NULL_TYPE(), pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), result)
	return result
}

// Java method: AstFactory.createVoid(com.google.javascript.rhino.Node child)
func (this_ *AstFactory) CreateVoid(child pkg_rhino.Node) pkg_rhino.Node {
	_ = child
	var result pkg_rhino.Node = pkg_rhino.IR_VoidNode(child)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE(), pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), result)
	return result
}

// Java method: AstFactory.createUndefinedValue()
func (this_ *AstFactory) CreateUndefinedValue() pkg_rhino.Node {
	return this_.CreateVoid(this_.CreateNumber(float64(int32(0))))
}

// Java method: AstFactory.createNot(com.google.javascript.rhino.Node child)
func (this_ *AstFactory) CreateNot(child pkg_rhino.Node) pkg_rhino.Node {
	_ = child
	var result pkg_rhino.Node = pkg_rhino.IR_Not(child)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE(), pkg_colors.StaticGet_StandardColors_BOOLEAN()), result)
	return result
}

// Java method: AstFactory.createThis(com.google.javascript.jscomp.AstFactory.Type thisType)
func (this_ *AstFactory) CreateThis(thisType AstFactory_Type) pkg_rhino.Node {
	_ = thisType
	var result pkg_rhino.Node = pkg_rhino.IR_ThisNode()
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(thisType, result)
	return result
}

// Java method: AstFactory.createSuper(com.google.javascript.jscomp.AstFactory.Type superType)
func (this_ *AstFactory) CreateSuper(superType AstFactory_Type) pkg_rhino.Node {
	_ = superType
	var result pkg_rhino.Node = pkg_rhino.IR_SuperNode()
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(superType, result)
	return result
}

// Java method: AstFactory.createNewTarget(com.google.javascript.jscomp.AstFactory.Type newTargetType)
func (this_ *AstFactory) CreateNewTarget(newTargetType AstFactory_Type) pkg_rhino.Node {
	_ = newTargetType
	var result pkg_rhino.Node = pkg_rhino.IR_NewTarget()
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(newTargetType, result)
	return result
}

// Java method: AstFactory.createThisForEs6Class(com.google.javascript.rhino.Node functionNode)
func (this_ *AstFactory) CreateThisForEs6Class(functionNode pkg_rhino.Node) pkg_rhino.Node {
	_ = functionNode
	if !(functionNode.IsClass()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(functionNode))
	}
	var result pkg_rhino.Node = pkg_rhino.IR_ThisNode()
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(this_.priv_AstFactory_getTypeOfThisForEs6Class_649(functionNode), result)
	return result
}

// Java method: AstFactory.createThisForEs6ClassMember(com.google.javascript.rhino.Node memberNode)
func (this_ *AstFactory) CreateThisForEs6ClassMember(memberNode pkg_rhino.Node) pkg_rhino.Node {
	_ = memberNode
	if !(memberNode.GetParent().IsClassMembers()) {
		j2gort.CheckArgument(false)
	}
	if !((memberNode.IsMemberFunctionDef() || memberNode.IsMemberFieldDef()) || memberNode.IsComputedFieldDef()) {
		j2gort.CheckArgument(false)
	}
	var classNode pkg_rhino.Node = memberNode.GetGrandparent()
	_ = classNode
	if memberNode.IsStaticMember() {
		var result pkg_rhino.Node = pkg_rhino.IR_ThisNode()
		_ = result
		this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_1_Node(classNode), result)
		return result
	} else {
		return this_.CreateThisForEs6Class(classNode)
	}
}

// Java method: AstFactory.getTypeOfThisForFunctionNode(com.google.javascript.rhino.Node functionNode)
func (this_ *AstFactory) priv_AstFactory_getTypeOfThisForFunctionNode_648(functionNode pkg_rhino.Node) pkg_rhino.JSType {
	_ = functionNode
	this_.priv_AstFactory_assertNotAddingColors_620()
	if this_.IsAddingTypes() {
		var functionType *pkg_rhino.FunctionType = this_.priv_AstFactory_getFunctionType_650(functionNode)
		_ = functionType
		return j2gort.CheckNotNullIface(functionType.GetTypeOfThis())
	} else {
		return nil
	}
}

// Java method: AstFactory.getTypeOfThisForEs6Class(com.google.javascript.rhino.Node functionNode)
func (this_ *AstFactory) priv_AstFactory_getTypeOfThisForEs6Class_649(functionNode pkg_rhino.Node) AstFactory_Type {
	_ = functionNode
	if !(functionNode.IsClass()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(functionNode))
	}
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		return AstFactory_S_Type_1_JSType(this_.priv_AstFactory_getTypeOfThisForFunctionNode_648(functionNode))
	case int32(1) /* TypeMode.COLOR */ :
		return AstFactory_S_Type_1_Color(this_.priv_AstFactory_getInstanceOfColor_711(functionNode.GetColor()))
	case int32(2) /* TypeMode.NONE */ :
		return AstFactory_NoTypeInformation()
	}
	panic("unreachable switch expression")
}

// Java method: AstFactory.getFunctionType(com.google.javascript.rhino.Node functionNode)
func (this_ *AstFactory) priv_AstFactory_getFunctionType_650(functionNode pkg_rhino.Node) *pkg_rhino.FunctionType {
	_ = functionNode
	if !(functionNode.IsFunction() || functionNode.IsClass()) {
		j2gort.CheckStateFmt(false, "not a function or class: %s", j2gort.BoxNode(functionNode))
	}
	this_.priv_AstFactory_assertNotAddingColors_620()
	var typeBeforeCast pkg_rhino.JSType = functionNode.GetJSTypeBeforeCast()
	_ = typeBeforeCast
	var functionType *pkg_rhino.FunctionType = nil
	_ = functionType
	if typeBeforeCast != nil && !typeBeforeCast.IsJavaNil_() {
		functionType = typeBeforeCast.AssertFunctionType()
	} else {
		functionType = functionNode.GetJSTypeRequired().AssertFunctionType()
	}
	return functionType
}

// Java method: AstFactory.createThisAliasReferenceForEs6Class(java.lang.String aliasName, com.google.javascript.rhino.Node functionNode)
func (this_ *AstFactory) CreateThisAliasReferenceForEs6Class(aliasName string, functionNode pkg_rhino.Node) pkg_rhino.Node {
	_ = aliasName
	_ = functionNode
	return this_.CreateName_2_string_AstFactory_Type(aliasName, this_.priv_AstFactory_getTypeOfThisForEs6Class_649(functionNode))
}

// Java method: AstFactory.createSingleNameDeclaration(com.google.javascript.rhino.Token tokenType, java.lang.String name, com.google.javascript.rhino.Node value)
func (this_ *AstFactory) CreateSingleNameDeclaration(tokenType pkg_rhino.Token, name string, value pkg_rhino.Node) pkg_rhino.Node {
	_ = tokenType
	_ = name
	_ = value
	switch (tokenType).Ordinal() {
	case int32(111) /* Token.LET */ :
		return this_.CreateSingleLetNameDeclaration_2(name, value)
	case int32(84) /* Token.VAR */ :
		return this_.CreateSingleVarNameDeclaration_2(name, value)
	case int32(97) /* Token.CONST */ :
		return this_.CreateSingleConstNameDeclaration(name, value)
	default:
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Unexpeted token type: "+j2gort.StrVal(tokenType))))
	}
	panic("unreachable switch expression")
}

// Java method: AstFactory.createSingleLetNameDeclaration(java.lang.String variableName)
func (this_ *AstFactory) CreateSingleLetNameDeclaration_1(variableName string) pkg_rhino.Node {
	_ = variableName
	return pkg_rhino.IR_Let_1(this_.CreateName_2_string_AstFactory_Type(variableName, AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE(), pkg_colors.StaticGet_StandardColors_NULL_OR_VOID())))
}

// Java method: AstFactory.createSingleLetNameDeclaration(java.lang.String variableName, com.google.javascript.rhino.Node value)
func (this_ *AstFactory) CreateSingleLetNameDeclaration_2(variableName string, value pkg_rhino.Node) pkg_rhino.Node {
	_ = variableName
	_ = value
	return pkg_rhino.IR_Let_2(this_.CreateName_2_string_AstFactory_Type(variableName, AstFactory_S_Type_1_Node(value)), value)
}

// Java method: AstFactory.createSingleVarNameDeclaration(java.lang.String variableName)
func (this_ *AstFactory) CreateSingleVarNameDeclaration_1(variableName string) pkg_rhino.Node {
	_ = variableName
	return pkg_rhino.IR_Var_1(this_.CreateName_2_string_AstFactory_Type(variableName, AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE(), pkg_colors.StaticGet_StandardColors_NULL_OR_VOID())))
}

// Java method: AstFactory.createSingleVarNameDeclaration(java.lang.String variableName, com.google.javascript.rhino.Node value)
func (this_ *AstFactory) CreateSingleVarNameDeclaration_2(variableName string, value pkg_rhino.Node) pkg_rhino.Node {
	_ = variableName
	_ = value
	return pkg_rhino.IR_Var_2(this_.CreateName_2_string_AstFactory_Type(variableName, AstFactory_S_Type_1_Node(value)), value)
}

// Java method: AstFactory.createSingleConstNameDeclaration(java.lang.String variableName, com.google.javascript.rhino.Node value)
func (this_ *AstFactory) CreateSingleConstNameDeclaration(variableName string, value pkg_rhino.Node) pkg_rhino.Node {
	_ = variableName
	_ = value
	var nameNode pkg_rhino.Node = this_.CreateConstantName(variableName, AstFactory_S_Type_2_JSType_Color(value.GetJSType(), value.GetColor()))
	_ = nameNode
	return pkg_rhino.IR_ConstNode(nameNode, value)
}

// Java method: AstFactory.createSingleConstObjectPatternDeclaration(com.google.javascript.rhino.Node objectPattern, com.google.javascript.rhino.Node value)
func (this_ *AstFactory) CreateSingleConstObjectPatternDeclaration(objectPattern pkg_rhino.Node, value pkg_rhino.Node) pkg_rhino.Node {
	_ = objectPattern
	_ = value
	if !(objectPattern.IsObjectPattern()) {
		j2gort.CheckStateFmt(false, "not an object pattern: %s", j2gort.BoxNode(objectPattern))
	}
	return pkg_rhino.IR_ConstNode(objectPattern, value)
}

// Java method: AstFactory.createArgumentsReference()
func (this_ *AstFactory) CreateArgumentsReference() pkg_rhino.Node {
	var result pkg_rhino.Node = pkg_rhino.IR_Name("arguments")
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		result.SetJSType(this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ARGUMENTS_TYPE()))
	case int32(1) /* TypeMode.COLOR */ :
		result.SetColor(this_.colorRegistry.Get(pkg_colors.StaticGet_StandardColors_ARGUMENTS_ID()))
	case int32(2) /* TypeMode.NONE */ :
	}
	return result
}

// Java method: AstFactory.createArgumentsAliasDeclaration(java.lang.String aliasName)
func (this_ *AstFactory) CreateArgumentsAliasDeclaration(aliasName string) pkg_rhino.Node {
	_ = aliasName
	return this_.CreateSingleConstNameDeclaration(aliasName, this_.CreateArgumentsReference())
}

// Java method: AstFactory.createName(java.lang.String name, com.google.javascript.jscomp.AstFactory.Type type)
func (this_ *AstFactory) CreateName_2_string_AstFactory_Type(name string, type_ AstFactory_Type) pkg_rhino.Node {
	_ = name
	_ = type_
	if !(!(j2gort.CheckNotNullString(name) == "$jscomp")) {
		j2gort.CheckArgumentFmt(false, "Use createQName(RuntimeJsLibManager.Field) to reference $jscomp.* methods. Found: %s", j2gort.BoxString(name))
	}
	var result pkg_rhino.Node = pkg_rhino.IR_Name(name)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	if this_.lifeCycleStage.IsNormalized() && javalang.StringStartsWith(j2gort.CheckNotNullString(name), "$jscomp") {
		result.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
	return result
}

// Java method: AstFactory.createConstantName(java.lang.String name, com.google.javascript.jscomp.AstFactory.Type type)
func (this_ *AstFactory) CreateConstantName(name string, type_ AstFactory_Type) pkg_rhino.Node {
	_ = name
	_ = type_
	if !(!(j2gort.CheckNotNullString(name) == "$jscomp")) {
		j2gort.CheckArgumentFmt(false, "Use createQName(RuntimeJsLibManager.Field) to reference $jscomp.* methods. Found: %s", j2gort.BoxString(name))
	}
	var result pkg_rhino.Node = pkg_rhino.IR_Name(name)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	if this_.lifeCycleStage.IsNormalized() {
		result.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
	return result
}

// Java method: AstFactory.createName(com.google.javascript.rhino.@org.jspecify.annotations.Nullable StaticScope scope, java.lang.String name)
func (this_ *AstFactory) CreateName_2_StaticScope_string(scope pkg_rhino.StaticScope, name string) pkg_rhino.Node {
	_ = scope
	_ = name
	if !(!(j2gort.CheckNotNullString(name) == "$jscomp")) {
		j2gort.CheckArgumentFmt(false, "Use createQName(RuntimeJsLibManager.Field) to reference $jscomp.* methods. Found: %s", j2gort.BoxString(name))
	}
	var result pkg_rhino.Node = pkg_rhino.IR_Name(name)
	_ = result
	if this_.lifeCycleStage.IsNormalized() || !(func() bool {
		re_5 := this_.typeMode
		ae_6 := StaticGet_AstFactory_TypeMode_NONE()
		if re_5 == 0 {
			j2gort.ThrowNPE()
		}
		return re_5 == ae_6
	}()) {
		j2gort.CheckNotNullIface(scope)
		var var_ pkg_rhino.StaticSlot = j2gort.AsIface[pkg_rhino.StaticSlot](scope.GetSlot(name))
		_ = var_
		if var_ == nil || var_.IsJavaNil_() {
			{
				var iifeCond_1 bool
				{

					re_7 := this_.typeMode
					ae_8 := StaticGet_AstFactory_TypeMode_NONE()
					if re_7 == 0 {
						j2gort.ThrowNPE()
					}
					iifeCond_1 = re_7 == ae_8
				}
				if !(iifeCond_1) {
					j2gort.CheckStateFmt(false, "Missing var %s in scope %s", j2gort.BoxString(name), scope)
				}
			}
		} else {
			var declaration pkg_rhino.StaticRef = j2gort.CheckNotNullIface(j2gort.AsIface[pkg_rhino.StaticRef](var_.GetDeclaration()))
			_ = declaration
			var varDefinitionNode pkg_rhino.Node = j2gort.CheckNotNullNode(declaration.GetNode())
			_ = varDefinitionNode
			if varDefinitionNode.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME()) {
				result.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
			}
			switch (this_.typeMode).Ordinal() {
			case int32(0) /* TypeMode.JSTYPE */ :
				var definitionType pkg_rhino.JSType = varDefinitionNode.GetJSType()
				_ = definitionType
				result.SetJSType(func() pkg_rhino.JSType {
					if definitionType != nil && !definitionType.IsJavaNil_() {
						return definitionType
					}
					return this_.unknownType
				}())
			case int32(1) /* TypeMode.COLOR */ :
				var definitionColor *pkg_colors.Color = varDefinitionNode.GetColor()
				_ = definitionColor
				result.SetColor(func() *pkg_colors.Color {
					if definitionColor != nil {
						return definitionColor
					}
					return pkg_colors.StaticGet_StandardColors_UNKNOWN()
				}())
			case int32(2) /* TypeMode.NONE */ :
			}
		}
	}
	return result
}

// Java method: AstFactory.createNameWithUnknownType(java.lang.String name)
func (this_ *AstFactory) CreateNameWithUnknownType(name string) pkg_rhino.Node {
	_ = name
	return this_.CreateName_2_string_AstFactory_Type(name, AstFactory_S_Type_2_JSType_Color(this_.unknownType, pkg_colors.StaticGet_StandardColors_UNKNOWN()))
}

// Java method: AstFactory.createQNameUsingJSTypeInfo(com.google.javascript.jscomp.TypedScope globalTypedScope, java.lang.String qname)
func (this_ *AstFactory) CreateQNameUsingJSTypeInfo(globalTypedScope *TypedScope, qname string) pkg_rhino.Node {
	_ = globalTypedScope
	_ = qname
	if !((globalTypedScope == nil) || globalTypedScope.IsGlobal()) {
		j2gort.CheckArgumentMsg(false, j2gort.PtrToAny(globalTypedScope))
	}
	this_.priv_AstFactory_assertNotAddingColors_620()
	var nameParts j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(StaticGet_AstFactory_DOT_SPLITTER(), qname)
	_ = nameParts
	if !(!(nameParts.IsEmpty())) {
		j2gort.CheckState(false)
	}
	var receiverPart string = j2gort.ListGetStr(nameParts, int32(0))
	_ = receiverPart
	var receiver pkg_rhino.Node = pkg_rhino.IR_Name(receiverPart)
	_ = receiver
	if this_.IsAddingTypes() {
		var var_ *TypedVar = j2gort.CheckNotNullPtr(j2gort.CastPtr[TypedVar](globalTypedScope.GetVar(receiverPart)))
		_ = var_
		receiver.SetJSType(j2gort.CheckNotNullIface(var_.GetType()))
	}
	var otherParts j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](nameParts.SubList(int32(1), nameParts.Size()))
	_ = otherParts
	return this_.CreateGetPropsWithoutColors(receiver, j2gort.AsIface[j2gort.Ext_java_lang_Iterable](otherParts))
}

// Java method: AstFactory.createQName(com.google.javascript.rhino.StaticScope scope, java.lang.String qname)
func (this_ *AstFactory) CreateQName_2_StaticScope_string(scope pkg_rhino.StaticScope, qname string) pkg_rhino.Node {
	_ = scope
	_ = qname
	return this_.CreateQName_2_StaticScope_Ext_java_lang_Iterable(scope, j2gort.SplitterSplitStr(StaticGet_AstFactory_DOT_SPLITTER(), qname))
}

// Java method: AstFactory.createQName(com.google.javascript.rhino.StaticScope scope, java.lang.Iterable<java.lang.String> names)
func (this_ *AstFactory) CreateQName_2_StaticScope_Ext_java_lang_Iterable(scope pkg_rhino.StaticScope, names j2gort.Ext_java_lang_Iterable) pkg_rhino.Node {
	_ = scope
	_ = names
	var baseName string = j2gort.CheckNotNullString(j2gort.CastString(j2gort.Ext_com_google_common_collect_Iterables_GetFirst(names, nil)))
	_ = baseName
	var propertyNames j2gort.Ext_java_lang_Iterable = j2gort.Ext_com_google_common_collect_Iterables_Skip(names, int32(1))
	_ = propertyNames
	return this_.CreateQName_3_StaticScope_string_Ext_java_lang_Iterable(scope, baseName, propertyNames)
}

// Java method: AstFactory.createQName(com.google.javascript.rhino.StaticScope scope, java.lang.String baseName, java.lang.String[] propertyNames)
func (this_ *AstFactory) CreateQName_3_StaticScope_string_slice_string(scope pkg_rhino.StaticScope, baseName string, propertyNames []string) pkg_rhino.Node {
	_ = scope
	_ = baseName
	_ = propertyNames
	j2gort.CheckNotNullString(baseName)
	return this_.CreateQName_3_StaticScope_string_Ext_java_lang_Iterable(scope, baseName, j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Arrays_AsList(func() []any {
		arr_9 := propertyNames
		if arr_9 == nil {
			return nil
		}
		out_10 := make([]any, len(arr_9))
		for i_, el_11 := range arr_9 {
			out_10[i_] = j2gort.BoxString(el_11)
		}
		return out_10
	}())))
}

// Java method: AstFactory.createQName(com.google.javascript.rhino.StaticScope scope, java.lang.String baseName, java.lang.Iterable<java.lang.String> propertyNames)
func (this_ *AstFactory) CreateQName_3_StaticScope_string_Ext_java_lang_Iterable(scope pkg_rhino.StaticScope, baseName string, propertyNames j2gort.Ext_java_lang_Iterable) pkg_rhino.Node {
	_ = scope
	_ = baseName
	_ = propertyNames
	var baseNameNode pkg_rhino.Node = this_.CreateName_2_StaticScope_string(scope, baseName)
	_ = baseNameNode
	return this_.priv_AstFactory_createQName_670(scope, baseName, baseNameNode, propertyNames)
}

// Java method: AstFactory.createQName(com.google.javascript.rhino.StaticScope scope, java.lang.String baseName, com.google.javascript.rhino.Node baseNameNode, java.lang.Iterable<java.lang.String> propertyNames)
func (this_ *AstFactory) priv_AstFactory_createQName_670(scope pkg_rhino.StaticScope, baseName string, baseNameNode pkg_rhino.Node, propertyNames j2gort.Ext_java_lang_Iterable) pkg_rhino.Node {
	_ = scope
	_ = baseName
	_ = baseNameNode
	_ = propertyNames
	var qname pkg_rhino.Node = baseNameNode
	_ = qname
	var name string = baseName
	_ = name
	{
		cur_12 := j2gort.NewIterCursor(j2gort.CheckNotNull(propertyNames))
		for cur_12.HasNext() {
			propertyName := cur_12.NextStr()
			_ = propertyName
			name = (j2gort.StrValString(name) + "." + j2gort.StrValString(propertyName))
			var type_ AstFactory_Type = nil
			_ = type_
			if this_.IsAddingTypes() || this_.IsAddingColors() {
				var def pkg_rhino.Node = j2gort.AsIface[pkg_rhino.StaticRef](j2gort.CheckNotNullIface(j2gort.AsIface[pkg_rhino.StaticSlot](scope.GetSlot(name))).GetDeclaration()).GetNode()
				_ = def
				type_ = AstFactory_S_Type_1_Node(def)
			}
			qname = this_.CreateGetProp(qname, propertyName, type_)
		}
	}
	return qname
}

// Java method: AstFactory.createQName(com.google.javascript.rhino.StaticScope scope, com.google.javascript.jscomp.js.RuntimeJsLibManager.JsLibField field)
func (this_ *AstFactory) CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(scope pkg_rhino.StaticScope, field RuntimeJsLibManager_JsLibField) pkg_rhino.Node {
	_ = scope
	_ = field
	var qname string = j2gort.AsIface[RuntimeJsLibManager_InjectedJsLibField](field.AssertInjected()).QualifiedName()
	_ = qname
	var parts j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(StaticGet_AstFactory_DOT_SPLITTER(), qname)
	_ = parts
	var baseName string = j2gort.CheckNotNullString(j2gort.CastString(j2gort.Ext_com_google_common_collect_Iterables_GetFirst(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](parts), nil)))
	_ = baseName
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(baseName), "$jscomp")) {
		j2gort.CheckStateFmt(false, "Unexpected Field name %s", j2gort.BoxString(baseName))
	}
	var baseNameNode pkg_rhino.Node
	_ = baseNameNode
	if j2gort.CheckNotNullString(baseName) == "$jscomp" {
		baseNameNode = this_.priv_AstFactory_createJscomp_676()
	} else {
		baseNameNode = this_.CreateName_2_string_AstFactory_Type(baseName, AstFactory_S_Type_2_JSType_Color(this_.unknownType, pkg_colors.StaticGet_StandardColors_UNKNOWN()))
	}
	if this_.lifeCycleStage.IsNormalized() {
		baseNameNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
	var propertyNames j2gort.Ext_java_lang_Iterable = j2gort.Ext_com_google_common_collect_Iterables_Skip(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](parts), int32(1))
	_ = propertyNames
	return this_.priv_AstFactory_createQName_670(scope, baseName, baseNameNode, propertyNames)
}

// Java method: AstFactory.createQNameWithUnknownType(java.lang.String qname)
func (this_ *AstFactory) CreateQNameWithUnknownType_1_string(qname string) pkg_rhino.Node {
	_ = qname
	return this_.priv_AstFactory_createQNameWithUnknownType_674(j2gort.SplitterSplitStr(StaticGet_AstFactory_DOT_SPLITTER(), qname))
}

// Java method: AstFactory.createQNameWithUnknownType(com.google.javascript.jscomp.js.RuntimeJsLibManager.JsLibField field)
func (this_ *AstFactory) CreateQNameWithUnknownType_1_RuntimeJsLibManager_JsLibField(field RuntimeJsLibManager_JsLibField) pkg_rhino.Node {
	_ = field
	var parts j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(StaticGet_AstFactory_DOT_SPLITTER(), j2gort.AsIface[RuntimeJsLibManager_InjectedJsLibField](field.AssertInjected()).QualifiedName())
	_ = parts
	var baseName string = j2gort.CheckNotNullString(j2gort.CastString(j2gort.Ext_com_google_common_collect_Iterables_GetFirst(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](parts), nil)))
	_ = baseName
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(baseName), "$jscomp")) {
		j2gort.CheckStateFmt(false, "Unexpected Field name %s", j2gort.BoxString(baseName))
	}
	var baseNameNode pkg_rhino.Node
	_ = baseNameNode
	if j2gort.CheckNotNullString(baseName) == "$jscomp" {
		baseNameNode = this_.priv_AstFactory_createJscomp_676()
	} else {
		baseNameNode = this_.CreateName_2_string_AstFactory_Type(baseName, AstFactory_S_Type_2_JSType_Color(this_.unknownType, pkg_colors.StaticGet_StandardColors_UNKNOWN()))
	}
	if this_.lifeCycleStage.IsNormalized() {
		baseNameNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
	var propertyNames j2gort.Ext_java_lang_Iterable = j2gort.Ext_com_google_common_collect_Iterables_Skip(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](parts), int32(1))
	_ = propertyNames
	return this_.CreateGetPropsWithUnknownType(baseNameNode, propertyNames)
}

// Java method: AstFactory.createQNameWithUnknownType(java.lang.Iterable<java.lang.String> names)
func (this_ *AstFactory) priv_AstFactory_createQNameWithUnknownType_674(names j2gort.Ext_java_lang_Iterable) pkg_rhino.Node {
	_ = names
	var baseName string = j2gort.CheckNotNullString(j2gort.CastString(j2gort.Ext_com_google_common_collect_Iterables_GetFirst(names, nil)))
	_ = baseName
	var propertyNames j2gort.Ext_java_lang_Iterable = j2gort.Ext_com_google_common_collect_Iterables_Skip(names, int32(1))
	_ = propertyNames
	return this_.CreateQNameWithUnknownType_2(baseName, propertyNames)
}

// Java method: AstFactory.createQNameWithUnknownType(java.lang.String baseName, java.lang.Iterable<java.lang.String> propertyNames)
func (this_ *AstFactory) CreateQNameWithUnknownType_2(baseName string, propertyNames j2gort.Ext_java_lang_Iterable) pkg_rhino.Node {
	_ = baseName
	_ = propertyNames
	var baseNameNode pkg_rhino.Node = this_.CreateNameWithUnknownType(baseName)
	_ = baseNameNode
	return this_.CreateGetPropsWithUnknownType(baseNameNode, propertyNames)
}

// Java method: AstFactory.createJscomp()
func (this_ *AstFactory) priv_AstFactory_createJscomp_676() pkg_rhino.Node {
	var jscomp pkg_rhino.Node = pkg_rhino.IR_Name("$jscomp")
	_ = jscomp
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSType_Color(this_.unknownType, pkg_colors.StaticGet_StandardColors_UNKNOWN()), jscomp)
	if this_.lifeCycleStage.IsNormalized() {
		jscomp.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
	return jscomp
}

// Java method: AstFactory.createPrototypeAccess(com.google.javascript.rhino.Node receiver)
func (this_ *AstFactory) CreatePrototypeAccess(receiver pkg_rhino.Node) pkg_rhino.Node {
	_ = receiver
	var result pkg_rhino.Node = pkg_rhino.IR_Getprop_2(receiver, "prototype")
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		result.SetJSType(this_.priv_AstFactory_getJsTypeForProperty_743(receiver, "prototype"))
	case int32(1) /* TypeMode.COLOR */ :
		j2gort.CheckNotNullPtr(receiver.GetColor())
		var possiblePrototypes []*pkg_colors.Color = receiver.GetColor().GetPrototypesArray()
		_ = possiblePrototypes
		result.SetColor(func() *pkg_colors.Color {
			if int32(len(possiblePrototypes)) == int32(0) {
				return pkg_colors.StaticGet_StandardColors_UNKNOWN()
			}
			return pkg_colors.Color_CreateUnion_1_Color(possiblePrototypes)
		}())
	case int32(2) /* TypeMode.NONE */ :
	}
	return result
}

// Java method: AstFactory.createJSCompDotGlobalAccess(com.google.javascript.rhino.StaticScope scope, java.lang.String qname)
func (this_ *AstFactory) CreateJSCompDotGlobalAccess(scope pkg_rhino.StaticScope, qname string) pkg_rhino.Node {
	_ = scope
	_ = qname
	var global RuntimeJsLibManager_JsLibField = this_.runtimeJsLibManager.GetJsLibField("$jscomp.global")
	_ = global
	var jscompDotGlobal pkg_rhino.Node = this_.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(scope, global)
	_ = jscompDotGlobal
	var result pkg_rhino.Node = this_.CreateQName_2_StaticScope_string(scope, qname)
	_ = result
	var qnameRoot pkg_rhino.Node = NodeUtil_GetRootOfQualifiedName_1_Node(result)
	_ = qnameRoot
	qnameRoot.ReplaceWith(this_.CreateGetProp(jscompDotGlobal, qnameRoot.GetString(), AstFactory_S_Type_1_Node(qnameRoot)))
	return result
}

// Java method: AstFactory.createGetPropWithoutColor(com.google.javascript.rhino.Node receiver, java.lang.String propertyName)
func (this_ *AstFactory) CreateGetPropWithoutColor(receiver pkg_rhino.Node, propertyName string) pkg_rhino.Node {
	_ = receiver
	_ = propertyName
	this_.priv_AstFactory_assertNotAddingColors_620()
	var result pkg_rhino.Node = pkg_rhino.IR_Getprop_2(receiver, propertyName)
	_ = result
	if this_.IsAddingTypes() {
		result.SetJSType(this_.priv_AstFactory_getJsTypeForProperty_743(receiver, propertyName))
	}
	return result
}

// Java method: AstFactory.createGetProp(com.google.javascript.rhino.Node receiver, java.lang.String propertyName, com.google.javascript.jscomp.AstFactory.Type type)
func (this_ *AstFactory) CreateGetProp(receiver pkg_rhino.Node, propertyName string, type_ AstFactory_Type) pkg_rhino.Node {
	_ = receiver
	_ = propertyName
	_ = type_
	var result pkg_rhino.Node = pkg_rhino.IR_Getprop_2(receiver, propertyName)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	return result
}

// Java method: AstFactory.createGetPropsWithoutColors(com.google.javascript.rhino.Node receiver, java.lang.Iterable<java.lang.String> propertyNames)
func (this_ *AstFactory) CreateGetPropsWithoutColors(receiver pkg_rhino.Node, propertyNames j2gort.Ext_java_lang_Iterable) pkg_rhino.Node {
	_ = receiver
	_ = propertyNames
	this_.priv_AstFactory_assertNotAddingColors_620()
	var result pkg_rhino.Node = receiver
	_ = result
	{
		cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNull(propertyNames))
		for cur_13.HasNext() {
			propertyName := cur_13.NextStr()
			_ = propertyName
			result = this_.CreateGetPropWithoutColor(result, propertyName)
		}
	}
	return result
}

// Java method: AstFactory.createGetPropsWithUnknownType(com.google.javascript.rhino.Node receiver, java.lang.Iterable<java.lang.String> propertyNames)
func (this_ *AstFactory) CreateGetPropsWithUnknownType(receiver pkg_rhino.Node, propertyNames j2gort.Ext_java_lang_Iterable) pkg_rhino.Node {
	_ = receiver
	_ = propertyNames
	var result pkg_rhino.Node = receiver
	_ = result
	{
		cur_14 := j2gort.NewIterCursor(j2gort.CheckNotNull(propertyNames))
		for cur_14.HasNext() {
			propertyName := cur_14.NextStr()
			_ = propertyName
			result = this_.CreateGetPropWithUnknownType(result, propertyName)
		}
	}
	return result
}

// Java method: AstFactory.createGetPropWithUnknownType(com.google.javascript.rhino.Node receiver, java.lang.String propertyName)
func (this_ *AstFactory) CreateGetPropWithUnknownType(receiver pkg_rhino.Node, propertyName string) pkg_rhino.Node {
	_ = receiver
	_ = propertyName
	var result pkg_rhino.Node = pkg_rhino.IR_Getprop_2(receiver, propertyName)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSType_Color(this_.unknownType, pkg_colors.StaticGet_StandardColors_UNKNOWN()), result)
	return result
}

// Java method: AstFactory.createStartOptChainGetprop(com.google.javascript.rhino.Node receiver, java.lang.String propertyName, com.google.javascript.jscomp.AstFactory.Type type)
func (this_ *AstFactory) CreateStartOptChainGetprop(receiver pkg_rhino.Node, propertyName string, type_ AstFactory_Type) pkg_rhino.Node {
	_ = receiver
	_ = propertyName
	_ = type_
	var result pkg_rhino.Node = pkg_rhino.IR_StartOptChainGetprop(receiver, propertyName)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	return result
}

// Java method: AstFactory.createContinueOptChainGetprop(com.google.javascript.rhino.Node receiver, java.lang.String propertyName, com.google.javascript.jscomp.AstFactory.Type type)
func (this_ *AstFactory) CreateContinueOptChainGetprop(receiver pkg_rhino.Node, propertyName string, type_ AstFactory_Type) pkg_rhino.Node {
	_ = receiver
	_ = propertyName
	_ = type_
	var result pkg_rhino.Node = pkg_rhino.IR_ContinueOptChainGetprop(receiver, propertyName)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	return result
}

// Java method: AstFactory.createStartOptChainGetelem(com.google.javascript.rhino.Node receiver, com.google.javascript.rhino.Node elem, com.google.javascript.jscomp.AstFactory.Type type)
func (this_ *AstFactory) CreateStartOptChainGetelem(receiver pkg_rhino.Node, elem pkg_rhino.Node, type_ AstFactory_Type) pkg_rhino.Node {
	_ = receiver
	_ = elem
	_ = type_
	var result pkg_rhino.Node = pkg_rhino.IR_StartOptChainGetelem(receiver, elem)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	return result
}

// Java method: AstFactory.createContinueOptChainGetelem(com.google.javascript.rhino.Node receiver, com.google.javascript.rhino.Node elem, com.google.javascript.jscomp.AstFactory.Type type)
func (this_ *AstFactory) CreateContinueOptChainGetelem(receiver pkg_rhino.Node, elem pkg_rhino.Node, type_ AstFactory_Type) pkg_rhino.Node {
	_ = receiver
	_ = elem
	_ = type_
	var result pkg_rhino.Node = pkg_rhino.IR_ContinueOptChainGetelem(receiver, elem)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	return result
}

// Java method: AstFactory.createStartOptChainCall(com.google.javascript.rhino.Node receiver, com.google.javascript.jscomp.AstFactory.Type type, com.google.javascript.rhino.Node[] args)
func (this_ *AstFactory) CreateStartOptChainCall(receiver pkg_rhino.Node, type_ AstFactory_Type, args []pkg_rhino.Node) pkg_rhino.Node {
	_ = receiver
	_ = type_
	_ = args
	var result pkg_rhino.Node = pkg_rhino.IR_StartOptChainCall(receiver, args)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	return result
}

// Java method: AstFactory.createContinueOptChainCall(com.google.javascript.rhino.Node receiver, com.google.javascript.jscomp.AstFactory.Type type, com.google.javascript.rhino.Node[] args)
func (this_ *AstFactory) CreateContinueOptChainCall(receiver pkg_rhino.Node, type_ AstFactory_Type, args []pkg_rhino.Node) pkg_rhino.Node {
	_ = receiver
	_ = type_
	_ = args
	var result pkg_rhino.Node = pkg_rhino.IR_ContinueOptChainCall(receiver, args)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	return result
}

// Java method: AstFactory.createGetElem(com.google.javascript.rhino.Node receiver, com.google.javascript.rhino.Node key)
func (this_ *AstFactory) CreateGetElem(receiver pkg_rhino.Node, key pkg_rhino.Node) pkg_rhino.Node {
	_ = receiver
	_ = key
	var result pkg_rhino.Node = pkg_rhino.IR_Getelem(receiver, key)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSType_Color(this_.unknownType, pkg_colors.StaticGet_StandardColors_UNKNOWN()), result)
	return result
}

// Java method: AstFactory.createDelProp(com.google.javascript.rhino.Node target)
func (this_ *AstFactory) CreateDelProp(target pkg_rhino.Node) pkg_rhino.Node {
	_ = target
	var result pkg_rhino.Node = pkg_rhino.IR_Delprop(target)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE(), pkg_colors.StaticGet_StandardColors_BOOLEAN()), result)
	return result
}

// Java method: AstFactory.createStringKey(java.lang.String key, com.google.javascript.rhino.Node value)
func (this_ *AstFactory) CreateStringKey(key string, value pkg_rhino.Node) pkg_rhino.Node {
	_ = key
	_ = value
	var result pkg_rhino.Node = pkg_rhino.IR_StringKey_2(key, value)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSType_Color(value.GetJSType(), value.GetColor()), result)
	return result
}

// Java method: AstFactory.createComputedProperty(com.google.javascript.rhino.Node key, com.google.javascript.rhino.Node value)
func (this_ *AstFactory) CreateComputedProperty(key pkg_rhino.Node, value pkg_rhino.Node) pkg_rhino.Node {
	_ = key
	_ = value
	var result pkg_rhino.Node = pkg_rhino.IR_ComputedProp(key, value)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSType_Color(value.GetJSType(), value.GetColor()), result)
	return result
}

// Java method: AstFactory.createGetterDef(java.lang.String name, com.google.javascript.rhino.Node value)
func (this_ *AstFactory) CreateGetterDef(name string, value pkg_rhino.Node) pkg_rhino.Node {
	_ = name
	_ = value
	var returnType pkg_rhino.JSType = value.GetJSType()
	_ = returnType
	var functionNode pkg_rhino.Node = this_.CreateZeroArgFunction("", pkg_rhino.IR_Block_1_Node(this_.CreateReturn(value)), returnType)
	_ = functionNode
	var getterNode pkg_rhino.Node = pkg_rhino.Node_NewString_2(pkg_rhino.StaticGet_Token_GETTER_DEF(), name)
	_ = getterNode
	getterNode.AddChildToFront(functionNode)
	return getterNode
}

// Java method: AstFactory.createIn(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *AstFactory) CreateIn(left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = left
	_ = right
	var result pkg_rhino.Node = pkg_rhino.IR_In(left, right)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE(), pkg_colors.StaticGet_StandardColors_BOOLEAN()), result)
	return result
}

// Java method: AstFactory.createComma(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *AstFactory) CreateComma(left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = left
	_ = right
	var result pkg_rhino.Node = pkg_rhino.IR_Comma(left, right)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSType_Color(right.GetJSType(), right.GetColor()), result)
	return result
}

// Java method: AstFactory.createCommas(com.google.javascript.rhino.Node first, com.google.javascript.rhino.Node second, com.google.javascript.rhino.Node[] rest)
func (this_ *AstFactory) CreateCommas(first pkg_rhino.Node, second pkg_rhino.Node, rest []pkg_rhino.Node) pkg_rhino.Node {
	_ = first
	_ = second
	_ = rest
	var result pkg_rhino.Node = this_.CreateComma(first, second)
	_ = result
	for _, el_15 := range j2gort.CheckNotNullSlice(rest) {
		next := el_15
		_ = next
		result = this_.CreateComma(result, next)
	}
	return result
}

// Java method: AstFactory.createAnd(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *AstFactory) CreateAnd(left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = left
	_ = right
	var result pkg_rhino.Node = pkg_rhino.IR_And(left, right)
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		var leftType pkg_rhino.JSType = j2gort.CheckNotNullIface(left.GetJSType())
		_ = leftType
		var rightType pkg_rhino.JSType = j2gort.CheckNotNullIface(right.GetJSType())
		_ = rightType
		result.SetJSType(this_.registry.CreateUnionType_2(leftType, rightType))
	case int32(1) /* TypeMode.COLOR */ :
		var leftColor *pkg_colors.Color = j2gort.CheckNotNullPtr(left.GetColor())
		_ = leftColor
		var rightColor *pkg_colors.Color = j2gort.CheckNotNullPtr(right.GetColor())
		_ = rightColor
		result.SetColor(pkg_colors.Color_CreateUnion2(leftColor, rightColor))
	case int32(2) /* TypeMode.NONE */ :
	}
	return result
}

// Java method: AstFactory.createOr(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *AstFactory) CreateOr(left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = left
	_ = right
	var result pkg_rhino.Node = pkg_rhino.IR_Or(left, right)
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		var leftType pkg_rhino.JSType = j2gort.CheckNotNullIface(left.GetJSType())
		_ = leftType
		var rightType pkg_rhino.JSType = j2gort.CheckNotNullIface(right.GetJSType())
		_ = rightType
		result.SetJSType(this_.registry.CreateUnionType_2(leftType, rightType))
	case int32(1) /* TypeMode.COLOR */ :
		var leftColor *pkg_colors.Color = j2gort.CheckNotNullPtr(left.GetColor())
		_ = leftColor
		var rightColor *pkg_colors.Color = j2gort.CheckNotNullPtr(right.GetColor())
		_ = rightColor
		result.SetColor(pkg_colors.Color_CreateUnion2(leftColor, rightColor))
	case int32(2) /* TypeMode.NONE */ :
	}
	return result
}

// Java method: AstFactory.createAdd(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *AstFactory) CreateAdd(left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = left
	_ = right
	var result pkg_rhino.Node = pkg_rhino.IR_Add(left, right)
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		result.SetJSType(this_.priv_AstFactory_getNativeType_742(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER_STRING()))
	case int32(1) /* TypeMode.COLOR */ :
		result.SetColor(j2gort.CastPtr[pkg_colors.Color](StaticGet_AstFactory_BigintNumberStringColor().Get()))
	case int32(2) /* TypeMode.NONE */ :
	}
	return result
}

// Java method: AstFactory.createSub(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *AstFactory) CreateSub(left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = left
	_ = right
	var result pkg_rhino.Node = pkg_rhino.IR_Sub(left, right)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE(), pkg_colors.StaticGet_StandardColors_NUMBER()), result)
	return result
}

// Java method: AstFactory.createInc(com.google.javascript.rhino.Node operand, boolean isPost)
func (this_ *AstFactory) CreateInc(operand pkg_rhino.Node, isPost bool) pkg_rhino.Node {
	_ = operand
	_ = isPost
	var result pkg_rhino.Node = pkg_rhino.IR_Inc(operand, isPost)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE(), pkg_colors.StaticGet_StandardColors_NUMBER()), result)
	return result
}

// Java method: AstFactory.createLessThan(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *AstFactory) CreateLessThan(left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = left
	_ = right
	var result pkg_rhino.Node = pkg_rhino.IR_Lt(left, right)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE(), pkg_colors.StaticGet_StandardColors_BOOLEAN()), result)
	return result
}

// Java method: AstFactory.createBitwiseAnd(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *AstFactory) CreateBitwiseAnd(left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = left
	_ = right
	var result pkg_rhino.Node = pkg_rhino.IR_BitwiseAnd(left, right)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE(), pkg_colors.StaticGet_StandardColors_NUMBER()), result)
	return result
}

// Java method: AstFactory.createRightShift(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *AstFactory) CreateRightShift(left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = left
	_ = right
	var result pkg_rhino.Node = pkg_rhino.IR_RightShift(left, right)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE(), pkg_colors.StaticGet_StandardColors_NUMBER()), result)
	return result
}

// Java method: AstFactory.createCall(com.google.javascript.rhino.Node callee, com.google.javascript.jscomp.AstFactory.Type resultType, com.google.javascript.rhino.Node[] args)
func (this_ *AstFactory) CreateCall(callee pkg_rhino.Node, resultType AstFactory_Type, args []pkg_rhino.Node) pkg_rhino.Node {
	_ = callee
	_ = resultType
	_ = args
	var result pkg_rhino.Node = NodeUtil_NewCallNode(callee, args)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(resultType, result)
	return result
}

// Java method: AstFactory.createCallWithUnknownType(com.google.javascript.rhino.Node callee, com.google.javascript.rhino.Node[] args)
func (this_ *AstFactory) CreateCallWithUnknownType(callee pkg_rhino.Node, args []pkg_rhino.Node) pkg_rhino.Node {
	_ = callee
	_ = args
	return this_.CreateCall(callee, AstFactory_S_Type_2_JSType_Color(this_.unknownType, pkg_colors.StaticGet_StandardColors_UNKNOWN()), args)
}

// Java method: AstFactory.createObjectDotAssignCall(com.google.javascript.rhino.StaticScope scope, com.google.javascript.jscomp.AstFactory.Type returnType, com.google.javascript.rhino.Node[] args)
func (this_ *AstFactory) CreateObjectDotAssignCall(scope pkg_rhino.StaticScope, returnType AstFactory_Type, args []pkg_rhino.Node) pkg_rhino.Node {
	_ = scope
	_ = returnType
	_ = args
	var objAssign pkg_rhino.Node = this_.CreateQName_3_StaticScope_string_slice_string(scope, "Object", []string{"assign"})
	_ = objAssign
	var result pkg_rhino.Node = this_.CreateCall(objAssign, returnType, args)
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		var returnJSType pkg_rhino.JSType = returnType.GetJSType(this_.registry)
		_ = returnJSType
		var objAssignType pkg_rhino.JSType = pkg_rhino.ToIface_JSType(this_.registry.CreateFunctionTypeWithVarArgs(returnJSType, []pkg_rhino.JSType{this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE()), this_.registry.CreateUnionType_1_JSTypeNative([]pkg_rhino.JSTypeNative{pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE(), pkg_rhino.StaticGet_JSTypeNative_NULL_TYPE()})}))
		_ = objAssignType
		objAssign.SetJSType(objAssignType)
	case int32(1) /* TypeMode.COLOR */, int32(2) /* TypeMode.NONE */ :
	}
	return result
}

// Java method: AstFactory.createJscompAssignCall(com.google.javascript.rhino.StaticScope scope, com.google.javascript.jscomp.AstFactory.Type returnType, com.google.javascript.rhino.Node[] args)
func (this_ *AstFactory) CreateJscompAssignCall(scope pkg_rhino.StaticScope, returnType AstFactory_Type, args []pkg_rhino.Node) pkg_rhino.Node {
	_ = scope
	_ = returnType
	_ = args
	var jscompAssign RuntimeJsLibManager_JsLibField = this_.runtimeJsLibManager.GetJsLibField("$jscomp.assign")
	_ = jscompAssign
	var call pkg_rhino.Node = this_.CreateCall(this_.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(scope, jscompAssign), returnType, args)
	_ = call
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	return call
}

// Java method: AstFactory.createNewNode(com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node[] args)
func (this_ *AstFactory) CreateNewNode(target pkg_rhino.Node, args []pkg_rhino.Node) pkg_rhino.Node {
	_ = target
	_ = args
	var result pkg_rhino.Node = pkg_rhino.IR_NewNode(target, args)
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		var instanceType pkg_rhino.JSType = target.GetJSType()
		_ = instanceType
		if instanceType.IsFunctionType() {
			instanceType = j2gort.AsIface[pkg_rhino.JSType](instanceType.ToMaybeFunctionType().GetInstanceType())
		} else {
			instanceType = this_.priv_AstFactory_getNativeType_742(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
		}
		result.SetJSType(instanceType)
	case int32(1) /* TypeMode.COLOR */ :
		result.SetColor(this_.priv_AstFactory_getInstanceOfColor_711(target.GetColor()))
	case int32(2) /* TypeMode.NONE */ :
	}
	return result
}

// Java method: AstFactory.getInstanceOfColor(com.google.javascript.jscomp.colors.Color color)
func (this_ *AstFactory) priv_AstFactory_getInstanceOfColor_711(color *pkg_colors.Color) *pkg_colors.Color {
	_ = color
	var possibleInstanceColors []*pkg_colors.Color = color.GetInstanceColorsArray()
	_ = possibleInstanceColors
	if int32(len(possibleInstanceColors)) == int32(0) {
		return pkg_colors.StaticGet_StandardColors_UNKNOWN()
	} else {
		return pkg_colors.Color_CreateUnion_1_Color(possibleInstanceColors)
	}
}

// Java method: AstFactory.createConstructorCall(com.google.javascript.jscomp.AstFactory.Type classType, com.google.javascript.rhino.Node callee, com.google.javascript.rhino.Node[] args)
func (this_ *AstFactory) CreateConstructorCall(classType AstFactory_Type, callee pkg_rhino.Node, args []pkg_rhino.Node) pkg_rhino.Node {
	_ = classType
	_ = callee
	_ = args
	var result pkg_rhino.Node = NodeUtil_NewCallNode(callee, args)
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		var classJSType pkg_rhino.JSType = classType.GetJSType(this_.registry)
		_ = classJSType
		var constructorType *pkg_rhino.FunctionType = j2gort.CheckNotNullPtr(classJSType.ToMaybeFunctionType())
		_ = constructorType
		var instanceType pkg_rhino.ObjectType = j2gort.CheckNotNullIface(constructorType.GetInstanceType())
		_ = instanceType
		result.SetJSType(instanceType)
	case int32(1) /* TypeMode.COLOR */ :
		result.SetColor(this_.priv_AstFactory_getInstanceOfColor_711(classType.GetColor(this_.colorRegistry)))
	case int32(2) /* TypeMode.NONE */ :
	}
	return result
}

// Java method: AstFactory.createAssignStatement(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node rhs)
func (this_ *AstFactory) CreateAssignStatement(lhs pkg_rhino.Node, rhs pkg_rhino.Node) pkg_rhino.Node {
	_ = lhs
	_ = rhs
	return this_.ExprResult(this_.CreateAssign_2_Node_Node(lhs, rhs))
}

// Java method: AstFactory.createAssign(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node rhs)
func (this_ *AstFactory) CreateAssign_2_Node_Node(lhs pkg_rhino.Node, rhs pkg_rhino.Node) pkg_rhino.Node {
	_ = lhs
	_ = rhs
	var result pkg_rhino.Node = pkg_rhino.IR_Assign(lhs, rhs)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_1_Node(rhs), result)
	return result
}

// Java method: AstFactory.createAssign(java.lang.String lhsName, com.google.javascript.rhino.Node rhs)
func (this_ *AstFactory) CreateAssign_2_string_Node(lhsName string, rhs pkg_rhino.Node) pkg_rhino.Node {
	_ = lhsName
	_ = rhs
	var name pkg_rhino.Node = this_.CreateName_2_string_AstFactory_Type(lhsName, AstFactory_S_Type_1_Node(rhs))
	_ = name
	return this_.CreateAssign_2_Node_Node(name, rhs)
}

// Java method: AstFactory.createObjectLit(com.google.javascript.rhino.Node[] elements)
func (this_ *AstFactory) CreateObjectLit_1(elements []pkg_rhino.Node) pkg_rhino.Node {
	_ = elements
	var result pkg_rhino.Node = pkg_rhino.IR_Objectlit(elements)
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		result.SetJSType(j2gort.AsIface[pkg_rhino.JSType](this_.registry.CreateAnonymousObjectType(nil)))
	case int32(1) /* TypeMode.COLOR */ :
		result.SetColor(pkg_colors.StaticGet_StandardColors_TOP_OBJECT())
	case int32(2) /* TypeMode.NONE */ :
	}
	return result
}

// Java method: AstFactory.createObjectLit(com.google.javascript.jscomp.AstFactory.Type type, com.google.javascript.rhino.Node[] elements)
func (this_ *AstFactory) CreateObjectLit_2(type_ AstFactory_Type, elements []pkg_rhino.Node) pkg_rhino.Node {
	_ = type_
	_ = elements
	var result pkg_rhino.Node = pkg_rhino.IR_Objectlit(elements)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	return result
}

// Java method: AstFactory.createQuotedStringKey(java.lang.String key, com.google.javascript.rhino.Node value)
func (this_ *AstFactory) CreateQuotedStringKey(key string, value pkg_rhino.Node) pkg_rhino.Node {
	_ = key
	_ = value
	var result pkg_rhino.Node = pkg_rhino.IR_StringKey_2(key, value)
	_ = result
	result.SetQuotedStringKey()
	return result
}

// Java method: AstFactory.createEmptyFunction(com.google.javascript.jscomp.AstFactory.Type type)
func (this_ *AstFactory) CreateEmptyFunction(type_ AstFactory_Type) pkg_rhino.Node {
	_ = type_
	var result pkg_rhino.Node = NodeUtil_EmptyFunction()
	_ = result
	if this_.IsAddingTypes() {
		if !(type_.GetJSType(this_.registry).IsFunctionType()) {
			j2gort.CheckArgumentMsg(false, type_)
		}
	}
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	return result
}

// Java method: AstFactory.createEmptyGeneratorFunction(com.google.javascript.jscomp.AstFactory.Type type)
func (this_ *AstFactory) CreateEmptyGeneratorFunction(type_ AstFactory_Type) pkg_rhino.Node {
	_ = type_
	var result pkg_rhino.Node = this_.CreateEmptyFunction(type_)
	_ = result
	result.SetIsGeneratorFunction(true)
	return result
}

// Java method: AstFactory.createFunction(java.lang.String name, com.google.javascript.rhino.Node paramList, com.google.javascript.rhino.Node body, com.google.javascript.jscomp.AstFactory.Type type)
func (this_ *AstFactory) CreateFunction(name string, paramList pkg_rhino.Node, body pkg_rhino.Node, type_ AstFactory_Type) pkg_rhino.Node {
	_ = name
	_ = paramList
	_ = body
	_ = type_
	var nameNode pkg_rhino.Node = this_.CreateName_2_string_AstFactory_Type(name, type_)
	_ = nameNode
	var result pkg_rhino.Node = pkg_rhino.IR_Function(nameNode, paramList, body)
	_ = result
	if this_.IsAddingTypes() {
		if !(type_.GetJSType(this_.registry).IsFunctionType()) {
			j2gort.CheckArgumentMsg(false, type_)
		}
	}
	this_.priv_AstFactory_setJSTypeOrColor_744(type_, result)
	return result
}

// Java method: AstFactory.createParamList(java.lang.String[] parameterNames)
func (this_ *AstFactory) CreateParamList(parameterNames []string) pkg_rhino.Node {
	_ = parameterNames
	var paramList pkg_rhino.Node = pkg_rhino.IR_ParamList([]pkg_rhino.Node{})
	_ = paramList
	for _, el_16 := range j2gort.CheckNotNullSlice(parameterNames) {
		parameterName := el_16
		_ = parameterName
		paramList.AddChildToBack(this_.CreateNameWithUnknownType(parameterName))
	}
	return paramList
}

// Java method: AstFactory.createZeroArgFunction(java.lang.String name, com.google.javascript.rhino.Node body, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType returnType)
func (this_ *AstFactory) CreateZeroArgFunction(name string, body pkg_rhino.Node, returnType pkg_rhino.JSType) pkg_rhino.Node {
	_ = name
	_ = body
	_ = returnType
	var functionType *pkg_rhino.FunctionType
	_ = functionType
	if this_.IsAddingTypes() {
		functionType = this_.registry.CreateFunctionType_2_JSType_JSType(returnType, []pkg_rhino.JSType{}).ToMaybeFunctionType()
	} else {
		functionType = nil
	}
	return this_.CreateFunction(name, pkg_rhino.IR_ParamList([]pkg_rhino.Node{}), body, AstFactory_S_Type_2_JSType_Color(pkg_rhino.ToIface_JSType(functionType), pkg_colors.StaticGet_StandardColors_TOP_OBJECT()))
}

// Java method: AstFactory.createZeroArgGeneratorFunction(java.lang.String name, com.google.javascript.rhino.Node body, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType returnType)
func (this_ *AstFactory) CreateZeroArgGeneratorFunction(name string, body pkg_rhino.Node, returnType pkg_rhino.JSType) pkg_rhino.Node {
	_ = name
	_ = body
	_ = returnType
	var result pkg_rhino.Node = this_.CreateZeroArgFunction(name, body, returnType)
	_ = result
	result.SetIsGeneratorFunction(true)
	return result
}

// Java method: AstFactory.createZeroArgArrowFunctionForExpression(com.google.javascript.rhino.Node expression)
func (this_ *AstFactory) CreateZeroArgArrowFunctionForExpression(expression pkg_rhino.Node) pkg_rhino.Node {
	_ = expression
	var result pkg_rhino.Node = pkg_rhino.IR_ArrowFunction(pkg_rhino.IR_Name(""), pkg_rhino.IR_ParamList([]pkg_rhino.Node{}), expression)
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		var fb_ pkg_rhino.FunctionType_Builder
		pkg_rhino.InitFunctionType_Builder(&fb_, this_.registry)
		(&fb_).WithReturnType_1(expression.GetJSTypeRequired())
		(&fb_).WithParameters_0()
		var functionType *pkg_rhino.FunctionType = (&fb_).BuildAndResolve()
		_ = functionType
		result.SetJSType(pkg_rhino.ToIface_JSType(functionType))
	case int32(1) /* TypeMode.COLOR */ :
		result.SetColor(pkg_colors.StaticGet_StandardColors_TOP_OBJECT())
	case int32(2) /* TypeMode.NONE */ :
	}
	return result
}

// Java method: AstFactory.createMemberFunctionDef(java.lang.String name, com.google.javascript.rhino.Node function)
func (this_ *AstFactory) CreateMemberFunctionDef(name string, function pkg_rhino.Node) pkg_rhino.Node {
	_ = name
	_ = function
	if !(j2gort.CheckNotNullString(function.GetFirstChild().GetString()) == "") {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(function))
	}
	var result pkg_rhino.Node = pkg_rhino.IR_MemberFunctionDef(name, function)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_1_Node(function), result)
	return result
}

// Java method: AstFactory.createSheq(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func (this_ *AstFactory) CreateSheq(expr1 pkg_rhino.Node, expr2 pkg_rhino.Node) pkg_rhino.Node {
	_ = expr1
	_ = expr2
	var result pkg_rhino.Node = pkg_rhino.IR_Sheq(expr1, expr2)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE(), pkg_colors.StaticGet_StandardColors_BOOLEAN()), result)
	return result
}

// Java method: AstFactory.createEq(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func (this_ *AstFactory) CreateEq(expr1 pkg_rhino.Node, expr2 pkg_rhino.Node) pkg_rhino.Node {
	_ = expr1
	_ = expr2
	var result pkg_rhino.Node = pkg_rhino.IR_Eq(expr1, expr2)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE(), pkg_colors.StaticGet_StandardColors_BOOLEAN()), result)
	return result
}

// Java method: AstFactory.createNe(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func (this_ *AstFactory) CreateNe(expr1 pkg_rhino.Node, expr2 pkg_rhino.Node) pkg_rhino.Node {
	_ = expr1
	_ = expr2
	var result pkg_rhino.Node = pkg_rhino.IR_Ne(expr1, expr2)
	_ = result
	this_.priv_AstFactory_setJSTypeOrColor_744(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE(), pkg_colors.StaticGet_StandardColors_BOOLEAN()), result)
	return result
}

// Java method: AstFactory.createHook(com.google.javascript.rhino.Node condition, com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func (this_ *AstFactory) CreateHook(condition pkg_rhino.Node, expr1 pkg_rhino.Node, expr2 pkg_rhino.Node) pkg_rhino.Node {
	_ = condition
	_ = expr1
	_ = expr2
	var result pkg_rhino.Node = pkg_rhino.IR_Hook(condition, expr1, expr2)
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		result.SetJSType(this_.registry.CreateUnionType_2(expr1.GetJSType(), expr2.GetJSType()))
	case int32(1) /* TypeMode.COLOR */ :
		result.SetColor(pkg_colors.Color_CreateUnion2(expr1.GetColor(), expr2.GetColor()))
	case int32(2) /* TypeMode.NONE */ :
	}
	return result
}

// Java method: AstFactory.createOptionalChainShortCircuit(com.google.javascript.rhino.Node receiver, com.google.javascript.rhino.Node resultExpr)
func (this_ *AstFactory) CreateOptionalChainShortCircuit(receiver pkg_rhino.Node, resultExpr pkg_rhino.Node) pkg_rhino.Node {
	_ = receiver
	_ = resultExpr
	var cond pkg_rhino.Node = this_.CreateEq(receiver, this_.CreateNull())
	_ = cond
	return this_.CreateHook(cond, this_.CreateUndefinedValue(), resultExpr)
}

// Java method: AstFactory.createArraylit(com.google.javascript.rhino.Node[] elements)
func (this_ *AstFactory) CreateArraylit_1_Node(elements []pkg_rhino.Node) pkg_rhino.Node {
	_ = elements
	return this_.CreateArraylit_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Arrays_AsList(func() []any {
		arr_17 := elements
		if arr_17 == nil {
			return nil
		}
		out_18 := make([]any, len(arr_17))
		for i_, el_19 := range arr_17 {
			out_18[i_] = j2gort.BoxNode(el_19)
		}
		return out_18
	}())))
}

// Java method: AstFactory.createArraylit(java.lang.Iterable<com.google.javascript.rhino.Node> elements)
func (this_ *AstFactory) CreateArraylit_1_Ext_java_lang_Iterable(elements j2gort.Ext_java_lang_Iterable) pkg_rhino.Node {
	_ = elements
	var result pkg_rhino.Node = pkg_rhino.IR_Arraylit_1_Ext_java_lang_Iterable(elements)
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		result.SetJSType(pkg_rhino.ToIface_JSType(this_.registry.CreateTemplatizedType_2_ObjectType_JSType(this_.registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE()), []pkg_rhino.JSType{this_.priv_AstFactory_getNativeType_742(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())})))
	case int32(1) /* TypeMode.COLOR */ :
		result.SetColor(this_.colorRegistry.Get(pkg_colors.StaticGet_StandardColors_ARRAY_ID()))
	case int32(2) /* TypeMode.NONE */ :
	}
	return result
}

// Java method: AstFactory.createJSCompMakeIteratorCall(com.google.javascript.rhino.Node iterable, com.google.javascript.rhino.StaticScope scope)
func (this_ *AstFactory) CreateJSCompMakeIteratorCall(iterable pkg_rhino.Node, scope pkg_rhino.StaticScope) pkg_rhino.Node {
	_ = iterable
	_ = scope
	var makeIterator RuntimeJsLibManager_JsLibField = this_.runtimeJsLibManager.GetJsLibField("$jscomp.makeIterator")
	_ = makeIterator
	var makeIteratorName pkg_rhino.Node = this_.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(scope, makeIterator)
	_ = makeIteratorName
	var type_ AstFactory_Type = nil
	_ = type_
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		var iterableType pkg_rhino.JSType = iterable.GetJSType().GetTemplateTypeMap().GetResolvedTemplateType(this_.registry.GetIterableValueTemplate())
		_ = iterableType
		var makeIteratorType pkg_rhino.JSType = makeIteratorName.GetJSType()
		_ = makeIteratorType
		makeIteratorName.SetJSType(this_.priv_AstFactory_replaceTemplate_739(makeIteratorType, j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(iterableType)))
		type_ = AstFactory_S_Type_1_JSType(makeIteratorName.GetJSType().AssertFunctionType().GetReturnType())
	case int32(1) /* TypeMode.COLOR */ :
		type_ = AstFactory_S_Type_1_Color(this_.colorRegistry.Get(pkg_colors.StaticGet_StandardColors_ITERATOR_ID()))
	case int32(2) /* TypeMode.NONE */ :
		type_ = AstFactory_NoTypeInformation()
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
	var call pkg_rhino.Node = this_.CreateCall(makeIteratorName, type_, []pkg_rhino.Node{iterable})
	_ = call
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	return call
}

// Java method: AstFactory.createJscompIteratorCloseCall(com.google.javascript.rhino.Node iterator, com.google.javascript.rhino.Node iterResult, com.google.javascript.rhino.StaticScope scope)
func (this_ *AstFactory) CreateJscompIteratorCloseCall(iterator pkg_rhino.Node, iterResult pkg_rhino.Node, scope pkg_rhino.StaticScope) pkg_rhino.Node {
	_ = iterator
	_ = iterResult
	_ = scope
	var iteratorClose RuntimeJsLibManager_JsLibField = this_.runtimeJsLibManager.GetJsLibField("$jscomp.iteratorClose")
	_ = iteratorClose
	var iteratorCloseName pkg_rhino.Node = this_.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(scope, iteratorClose)
	_ = iteratorCloseName
	return this_.CreateCall(iteratorCloseName, AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE(), pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), []pkg_rhino.Node{iterator, iterResult})
}

// Java method: AstFactory.createJscompArrayFromIteratorCall(com.google.javascript.rhino.Node iterator, com.google.javascript.rhino.StaticScope scope)
func (this_ *AstFactory) CreateJscompArrayFromIteratorCall(iterator pkg_rhino.Node, scope pkg_rhino.StaticScope) pkg_rhino.Node {
	_ = iterator
	_ = scope
	var arrayFromIterator RuntimeJsLibManager_JsLibField = this_.runtimeJsLibManager.GetJsLibField("$jscomp.arrayFromIterator")
	_ = arrayFromIterator
	var makeIteratorName pkg_rhino.Node = this_.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(scope, arrayFromIterator)
	_ = makeIteratorName
	var resultType AstFactory_Type = func() AstFactory_Type {
		switch (this_.typeMode).Ordinal() {
		case int32(0) /* TypeMode.JSTYPE */ :
			var iterableType pkg_rhino.JSType = iterator.GetJSType().GetTemplateTypeMap().GetResolvedTemplateType(this_.registry.GetIteratorValueTemplate())
			_ = iterableType
			var makeIteratorType pkg_rhino.JSType = makeIteratorName.GetJSType()
			_ = makeIteratorType
			makeIteratorName.SetJSType(this_.priv_AstFactory_replaceTemplate_739(makeIteratorType, j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(iterableType)))
			return AstFactory_S_Type_1_JSType(makeIteratorName.GetJSType().AssertFunctionType().GetReturnType())
		case int32(1) /* TypeMode.COLOR */ :
			return AstFactory_S_Type_1_Color(this_.colorRegistry.Get(pkg_colors.StaticGet_StandardColors_ARRAY_ID()))
		case int32(2) /* TypeMode.NONE */ :
			return AstFactory_NoTypeInformation()
		}
		panic("unreachable switch expression")
	}()
	_ = resultType
	var call pkg_rhino.Node = this_.CreateCall(makeIteratorName, resultType, []pkg_rhino.Node{iterator})
	_ = call
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	return call
}

// Java method: AstFactory.createJscompArrayFromIterableCall(com.google.javascript.rhino.Node iterable, com.google.javascript.rhino.StaticScope scope)
func (this_ *AstFactory) CreateJscompArrayFromIterableCall(iterable pkg_rhino.Node, scope pkg_rhino.StaticScope) pkg_rhino.Node {
	_ = iterable
	_ = scope
	var arrayFromIterable RuntimeJsLibManager_JsLibField = this_.runtimeJsLibManager.GetJsLibField("$jscomp.arrayFromIterable")
	_ = arrayFromIterable
	var makeIterableName pkg_rhino.Node = this_.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(scope, arrayFromIterable)
	_ = makeIterableName
	var resultType AstFactory_Type = nil
	_ = resultType
	resultType = func() AstFactory_Type {
		switch (this_.typeMode).Ordinal() {
		case int32(0) /* TypeMode.JSTYPE */ :
			var iterableType pkg_rhino.JSType = iterable.GetJSType().GetTemplateTypeMap().GetResolvedTemplateType(this_.registry.GetIterableValueTemplate())
			_ = iterableType
			var makeIterableType pkg_rhino.JSType = makeIterableName.GetJSType()
			_ = makeIterableType
			makeIterableName.SetJSType(this_.priv_AstFactory_replaceTemplate_739(makeIterableType, j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(iterableType)))
			return AstFactory_S_Type_1_JSType(makeIterableName.GetJSType().AssertFunctionType().GetReturnType())
		case int32(1) /* TypeMode.COLOR */ :
			return func() AstFactory_Type {
				var asgn_20 AstFactory_Type = AstFactory_S_Type_1_Color(this_.colorRegistry.Get(pkg_colors.StaticGet_StandardColors_ARRAY_ID()))
				resultType = asgn_20
				return asgn_20
			}()
		case int32(2) /* TypeMode.NONE */ :
			return AstFactory_NoTypeInformation()
		}
		panic("unreachable switch expression")
	}()
	var call pkg_rhino.Node = this_.CreateCall(makeIterableName, resultType, []pkg_rhino.Node{iterable})
	_ = call
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	return call
}

// Java method: AstFactory.createJSCompMakeAsyncIteratorCall(com.google.javascript.rhino.Node iterable, com.google.javascript.rhino.StaticScope scope)
func (this_ *AstFactory) CreateJSCompMakeAsyncIteratorCall(iterable pkg_rhino.Node, scope pkg_rhino.StaticScope) pkg_rhino.Node {
	_ = iterable
	_ = scope
	var makeAsyncIterator RuntimeJsLibManager_JsLibField = this_.runtimeJsLibManager.GetJsLibField("$jscomp.makeAsyncIterator")
	_ = makeAsyncIterator
	var makeIteratorAsyncName pkg_rhino.Node = this_.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(scope, makeAsyncIterator)
	_ = makeIteratorAsyncName
	var resultType AstFactory_Type = nil
	_ = resultType
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		var asyncIterableType pkg_rhino.JSType = JsIterables_MaybeBoxIterableOrAsyncIterable(iterable.GetJSType(), this_.registry).OrElse(this_.unknownType)
		_ = asyncIterableType
		var makeAsyncIteratorType pkg_rhino.JSType = makeIteratorAsyncName.GetJSType()
		_ = makeAsyncIteratorType
		makeIteratorAsyncName.SetJSType(this_.priv_AstFactory_replaceTemplate_739(makeAsyncIteratorType, j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(asyncIterableType)))
		resultType = AstFactory_S_Type_1_JSType(makeIteratorAsyncName.GetJSType().AssertFunctionType().GetReturnType())
	case int32(1) /* TypeMode.COLOR */ :
		resultType = AstFactory_S_Type_1_Color(this_.colorRegistry.Get(pkg_colors.StaticGet_StandardColors_ASYNC_ITERATOR_ITERABLE_ID()))
	case int32(2) /* TypeMode.NONE */ :
		resultType = AstFactory_NoTypeInformation()
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
	var call pkg_rhino.Node = this_.CreateCall(makeIteratorAsyncName, resultType, []pkg_rhino.Node{iterable})
	_ = call
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	return call
}

// Java method: AstFactory.replaceTemplate(com.google.javascript.rhino.jstype.JSType templatedType, com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.JSType> templateTypes)
func (this_ *AstFactory) priv_AstFactory_replaceTemplate_739(templatedType pkg_rhino.JSType, templateTypes j2gort.Ext_com_google_common_collect_ImmutableList) pkg_rhino.JSType {
	_ = templatedType
	_ = templateTypes
	var typeMap *pkg_rhino.TemplateTypeMap = this_.registry.GetEmptyTemplateTypeMap().CopyWithKeysFromAndValues(templatedType.GetTemplateTypeMap(), templateTypes)
	_ = typeMap
	return this_.registry.ReplaceTemplateTypesPartial(templatedType, typeMap)
}

// Java method: AstFactory.createEmptyAsyncGeneratorWrapperArgument(com.google.javascript.rhino.jstype.JSType asyncGeneratorWrapperType)
func (this_ *AstFactory) CreateEmptyAsyncGeneratorWrapperArgument(asyncGeneratorWrapperType pkg_rhino.JSType) pkg_rhino.Node {
	_ = asyncGeneratorWrapperType
	var generatorType AstFactory_Type = AstFactory_NoTypeInformation()
	_ = generatorType
	if this_.IsAddingTypes() {
		if asyncGeneratorWrapperType.IsUnknownType() {
			generatorType = AstFactory_S_Type_1_JSType(pkg_rhino.ToIface_JSType(this_.registry.CreateFunctionType_2_JSType_JSType(this_.priv_AstFactory_replaceTemplate_739(this_.priv_AstFactory_getNativeType_742(pkg_rhino.StaticGet_JSTypeNative_GENERATOR_TYPE()), j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(this_.unknownType)), []pkg_rhino.JSType{})))
		} else {
			var wrapperFnType *pkg_rhino.FunctionType = asyncGeneratorWrapperType.ToMaybeFunctionType()
			_ = wrapperFnType
			var innerFunctionReturnType pkg_rhino.JSType
			_ = innerFunctionReturnType
			if wrapperFnType.GetParameterCount() == int32(1) {
				innerFunctionReturnType = wrapperFnType.GetParameter(int32(0)).GetJSType()
			} else {
				innerFunctionReturnType = j2gort.Cast[pkg_rhino.FunctionType_Parameter](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](wrapperFnType.GetParameters()))).GetJSType()
			}
			generatorType = AstFactory_S_Type_1_JSType(pkg_rhino.ToIface_JSType(this_.registry.CreateFunctionType_2_JSType_JSType(innerFunctionReturnType, []pkg_rhino.JSType{})))
		}
	} else if this_.IsAddingColors() {
		generatorType = AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT())
	}
	return this_.CreateEmptyGeneratorFunction(generatorType)
}

// Java method: AstFactory.createJscompAsyncExecutePromiseGeneratorFunctionCall(com.google.javascript.rhino.StaticScope scope, com.google.javascript.rhino.Node generatorFunction)
func (this_ *AstFactory) CreateJscompAsyncExecutePromiseGeneratorFunctionCall(scope pkg_rhino.StaticScope, generatorFunction pkg_rhino.Node) pkg_rhino.Node {
	_ = scope
	_ = generatorFunction
	var method RuntimeJsLibManager_JsLibField = this_.runtimeJsLibManager.GetJsLibField("$jscomp.asyncExecutePromiseGeneratorFunction")
	_ = method
	var jscompDotAsyncExecutePromiseGeneratorFunction pkg_rhino.Node = this_.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(scope, method)
	_ = jscompDotAsyncExecutePromiseGeneratorFunction
	var resultType AstFactory_Type = nil
	_ = resultType
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		resultType = AstFactory_S_Type_1_JSType(jscompDotAsyncExecutePromiseGeneratorFunction.GetJSType().AssertFunctionType().GetReturnType())
	case int32(1) /* TypeMode.COLOR */ :
		resultType = AstFactory_S_Type_1_Color(this_.colorRegistry.Get(pkg_colors.StaticGet_StandardColors_PROMISE_ID()))
	case int32(2) /* TypeMode.NONE */ :
		resultType = AstFactory_S_Type_1_JSType(this_.unknownType)
	}
	var call pkg_rhino.Node = this_.CreateCall(jscompDotAsyncExecutePromiseGeneratorFunction, resultType, []pkg_rhino.Node{generatorFunction})
	_ = call
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	return call
}

// Java method: AstFactory.getNativeType(com.google.javascript.rhino.jstype.JSTypeNative nativeType)
func (this_ *AstFactory) priv_AstFactory_getNativeType_742(nativeType pkg_rhino.JSTypeNative) pkg_rhino.JSType {
	_ = nativeType
	j2gort.CheckNotNullPtr(this_.registry)
	return j2gort.CheckNotNullIface(this_.registry.GetNativeType(nativeType))
}

// Java method: AstFactory.getJsTypeForProperty(com.google.javascript.rhino.Node receiver, java.lang.String propertyName)
func (this_ *AstFactory) priv_AstFactory_getJsTypeForProperty_743(receiver pkg_rhino.Node, propertyName string) pkg_rhino.JSType {
	_ = receiver
	_ = propertyName
	var getpropType pkg_rhino.JSType = nil
	_ = getpropType
	var receiverJSType pkg_rhino.JSType = receiver.GetJSType()
	_ = receiverJSType
	if receiverJSType != nil && !receiverJSType.IsJavaNil_() {
		getpropType = receiverJSType.FindPropertyType_1_string(propertyName)
		if getpropType == nil || getpropType.IsJavaNil_() {
			var receiverObjectType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(receiverJSType.Autobox())
			_ = receiverObjectType
			if receiverObjectType == nil || receiverObjectType.IsJavaNil_() {
				getpropType = this_.unknownType
			} else {
				getpropType = receiverObjectType.GetPropertyType_1_string(propertyName)
			}
		}
	}
	if getpropType == nil || getpropType.IsJavaNil_() {
		getpropType = this_.unknownType
	}
	if (getpropType.IsUnknownType() && (j2gort.CheckNotNullString(propertyName) == "global")) && receiver.MatchesName_1_string("$jscomp") {
		getpropType = this_.priv_AstFactory_getNativeType_742(pkg_rhino.StaticGet_JSTypeNative_GLOBAL_THIS())
	}
	return getpropType
}

// Java method: AstFactory.setJSTypeOrColor(com.google.javascript.jscomp.AstFactory.Type type, com.google.javascript.rhino.Node result)
func (this_ *AstFactory) priv_AstFactory_setJSTypeOrColor_744(type_ AstFactory_Type, result pkg_rhino.Node) {
	_ = type_
	_ = result
	switch (this_.typeMode).Ordinal() {
	case int32(0) /* TypeMode.JSTYPE */ :
		result.SetJSType(type_.GetJSType(this_.registry))
	case int32(1) /* TypeMode.COLOR */ :
		result.SetColor(type_.GetColor(this_.colorRegistry))
	case int32(2) /* TypeMode.NONE */ :
	}
}

// Java static method: AstFactory.type(com.google.javascript.rhino.Node node)
func AstFactory_S_Type_1_Node(node pkg_rhino.Node) AstFactory_Type {
	clinit_AstFactory()
	_ = node
	return NewAstFactory_TypeOnNode(node)
}

// Java static method: AstFactory.type(com.google.javascript.rhino.jstype.JSType type)
func AstFactory_S_Type_1_JSType(type_ pkg_rhino.JSType) AstFactory_Type {
	clinit_AstFactory()
	_ = type_
	return NewAstFactory_JSTypeOrColor_2_JSType_Color(type_, nil)
}

// Java static method: AstFactory.type(com.google.javascript.rhino.jstype.JSTypeNative type)
func AstFactory_S_Type_1_JSTypeNative(type_ pkg_rhino.JSTypeNative) AstFactory_Type {
	clinit_AstFactory()
	_ = type_
	return NewAstFactory_JSTypeOrColor_2_JSTypeNative_Color(type_, nil)
}

// Java static method: AstFactory.type(com.google.javascript.jscomp.colors.Color type)
func AstFactory_S_Type_1_Color(type_ *pkg_colors.Color) AstFactory_Type {
	clinit_AstFactory()
	_ = type_
	return NewAstFactory_JSTypeOrColor_2_JSType_Color(nil, type_)
}

// Java static method: AstFactory.type(com.google.javascript.jscomp.colors.ColorId type)
func AstFactory_S_Type_1_ColorId(type_ pkg_colors.ColorId) AstFactory_Type {
	clinit_AstFactory()
	_ = type_
	return NewAstFactory_JSTypeOrColor_2_JSTypeNative_ColorId(pkg_rhino.JSTypeNative(0), type_)
}

// Java static method: AstFactory.type(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.jscomp.colors.Color color)
func AstFactory_S_Type_2_JSType_Color(type_ pkg_rhino.JSType, color *pkg_colors.Color) AstFactory_Type {
	clinit_AstFactory()
	_ = type_
	_ = color
	return NewAstFactory_JSTypeOrColor_2_JSType_Color(type_, color)
}

// Java static method: AstFactory.type(com.google.javascript.rhino.jstype.JSTypeNative type, com.google.javascript.jscomp.colors.Color color)
func AstFactory_S_Type_2_JSTypeNative_Color(type_ pkg_rhino.JSTypeNative, color *pkg_colors.Color) AstFactory_Type {
	clinit_AstFactory()
	_ = type_
	_ = color
	return NewAstFactory_JSTypeOrColor_2_JSTypeNative_Color(type_, color)
}

// Java static method: AstFactory.noTypeInformation()
func AstFactory_NoTypeInformation() AstFactory_Type {
	clinit_AstFactory()
	return NewAstFactory_JSTypeOrColor_2_JSType_Color(nil, nil)
}

// Java enum: com.google.javascript.jscomp.AstFactory.TypeMode
type AstFactory_TypeMode uint8

const (
	enum_AstFactory_TypeMode_JSTYPE AstFactory_TypeMode = 1
	enum_AstFactory_TypeMode_COLOR  AstFactory_TypeMode = 2
	enum_AstFactory_TypeMode_NONE   AstFactory_TypeMode = 3
)

var enumBaseTable_AstFactory_TypeMode = [...]j2gort.EnumBase{
	{Name_: "JSTYPE", Ordinal_: 0},
	{Name_: "COLOR", Ordinal_: 1},
	{Name_: "NONE", Ordinal_: 2},
}
var enumValues_AstFactory_TypeMode = []AstFactory_TypeMode{1, 2, 3}

func StaticGet_AstFactory_TypeMode_JSTYPE() AstFactory_TypeMode {
	return enum_AstFactory_TypeMode_JSTYPE
}
func StaticGet_AstFactory_TypeMode_COLOR() AstFactory_TypeMode {
	return enum_AstFactory_TypeMode_COLOR
}
func StaticGet_AstFactory_TypeMode_NONE() AstFactory_TypeMode {
	return enum_AstFactory_TypeMode_NONE
}

func AstFactory_TypeMode_Values() []AstFactory_TypeMode {
	return j2gort.CloneSlice(enumValues_AstFactory_TypeMode)
}

func AstFactory_TypeMode_ValueOf(name string) AstFactory_TypeMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_AstFactory_TypeMode {
		if enumBaseTable_AstFactory_TypeMode[i_].Name_ == name {
			return AstFactory_TypeMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.AstFactory.TypeMode." + name))
}

func (this_ AstFactory_TypeMode) IsJavaNil_() bool                              { return this_ == 0 }
func (AstFactory_TypeMode) HasIdentityEquals_()                                 {}
func (AstFactory_TypeMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (AstFactory_TypeMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ AstFactory_TypeMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (AstFactory_TypeMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AstFactory$TypeMode", "TypeMode"
}
func (this_ AstFactory_TypeMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_AstFactory_TypeMode[this_-1]
}
func (this_ AstFactory_TypeMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ AstFactory_TypeMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ AstFactory_TypeMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AstFactory_TypeMode[this_-1].Name_
}
func (this_ AstFactory_TypeMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_AstFactory_TypeMode[this_-1].Name_
}
func (this_ AstFactory_TypeMode) Equals(other_ any) bool {
	o_, ok_ := other_.(AstFactory_TypeMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ AstFactory_TypeMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_AstFactory_TypeMode[this_-1].HashCode()
}
func (this_ AstFactory_TypeMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AstFactory_TypeMode[this_-1].CompareTo(other_)
}
func (this_ AstFactory_TypeMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_AstFactory_TypeMode[this_-1].Name_)
}
func init() {
	vals_ := []any{AstFactory_TypeMode(1), AstFactory_TypeMode(2), AstFactory_TypeMode(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AstFactory.TypeMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AstFactory$TypeMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AstFactory.TypeMode", "com.google.javascript.jscomp.AstFactory$TypeMode", func(v_ any) bool {
		return j2gort.InstanceOf[AstFactory_TypeMode](v_)
	})
}

// Java interface: com.google.javascript.jscomp.AstFactory.Type
type AstFactory_Type interface {
	Implements_jscomp_AstFactory_Type()
	IsJavaNil_() bool
	GetJSType(*pkg_rhino.JSTypeRegistry) pkg_rhino.JSType
	GetColor(*pkg_colors.ColorRegistry) *pkg_colors.Color
}

func ToIface_AstFactory_Type[T interface {
	AstFactory_Type
	comparable
}](v T) AstFactory_Type {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AstFactory.Type", "com.google.javascript.jscomp.AstFactory$Type", func(v_ any) bool {
		return j2gort.InstanceOf[AstFactory_Type](v_)
	})
}

// Java class: com.google.javascript.jscomp.AstFactory.TypeOnNode
type AstFactory_TypeOnNode struct {
	n       pkg_rhino.Node `gson:"n"`
	idHash_ atomic.Int32
}

func (this_ *AstFactory_TypeOnNode) BindSelf_(s any) {
	_ = s
}
func (this_ *AstFactory_TypeOnNode) IsJavaNil_() bool { return this_ == nil }
func (*AstFactory_TypeOnNode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AstFactory$TypeOnNode", "TypeOnNode"
}
func (this_ *AstFactory_TypeOnNode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AstFactory_TypeOnNode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AstFactory.TypeOnNode", "com.google.javascript.jscomp.AstFactory$TypeOnNode", func(v_ any) bool {
		return j2gort.InstanceOf[*AstFactory_TypeOnNode](v_)
	})
}
func (*AstFactory_TypeOnNode) Implements_jscomp_AstFactory_Type() {}

func (this_ *AstFactory_TypeOnNode) initFields_AstFactory_TypeOnNode() {
	_ = this_
}

// Java constructor: TypeOnNode(com.google.javascript.rhino.Node n)
func NewAstFactory_TypeOnNode(n pkg_rhino.Node) *AstFactory_TypeOnNode {
	this_ := &AstFactory_TypeOnNode{}
	InitAstFactory_TypeOnNode(this_, n)
	return this_
}

func InitAstFactory_TypeOnNode(this_ *AstFactory_TypeOnNode, n pkg_rhino.Node) {
	_ = n
	this_.initFields_AstFactory_TypeOnNode()
	this_.n = n
}

// Java method: TypeOnNode.getJSType(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func (this_ *AstFactory_TypeOnNode) GetJSType(registry *pkg_rhino.JSTypeRegistry) pkg_rhino.JSType {
	_ = registry
	var jstype pkg_rhino.JSType = this_.n.GetJSType()
	_ = jstype
	if jstype != nil && !jstype.IsJavaNil_() {
		return jstype
	} else {
		return registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	}
}

// Java method: TypeOnNode.getColor(com.google.javascript.jscomp.colors.ColorRegistry registry)
func (this_ *AstFactory_TypeOnNode) GetColor(registry *pkg_colors.ColorRegistry) *pkg_colors.Color {
	_ = registry
	var color *pkg_colors.Color = this_.n.GetColor()
	_ = color
	if color != nil {
		return color
	} else {
		return pkg_colors.StaticGet_StandardColors_UNKNOWN()
	}
}

// Java class: com.google.javascript.jscomp.AstFactory.JSTypeOrColor
type AstFactory_JSTypeOrColor struct {
	jstype       pkg_rhino.JSType   `gson:"jstype"`
	color        *pkg_colors.Color  `gson:"color"`
	colorId      pkg_colors.ColorId `gson:"colorId"`
	idHash_      atomic.Int32
	jstypeNative pkg_rhino.JSTypeNative `gson:"jstypeNative"`
}

func (this_ *AstFactory_JSTypeOrColor) BindSelf_(s any) {
	_ = s
}
func (this_ *AstFactory_JSTypeOrColor) IsJavaNil_() bool { return this_ == nil }
func (*AstFactory_JSTypeOrColor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AstFactory$JSTypeOrColor", "JSTypeOrColor"
}
func (this_ *AstFactory_JSTypeOrColor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AstFactory_JSTypeOrColor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AstFactory.JSTypeOrColor", "com.google.javascript.jscomp.AstFactory$JSTypeOrColor", func(v_ any) bool {
		return j2gort.InstanceOf[*AstFactory_JSTypeOrColor](v_)
	})
}
func (*AstFactory_JSTypeOrColor) Implements_jscomp_AstFactory_Type() {}

func (this_ *AstFactory_JSTypeOrColor) initFields_AstFactory_JSTypeOrColor() {
	_ = this_
}

// Java constructor: JSTypeOrColor(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSTypeNative jstypeNative, com.google.javascript.jscomp.colors.ColorId colorId)
func NewAstFactory_JSTypeOrColor_2_JSTypeNative_ColorId(jstypeNative pkg_rhino.JSTypeNative, colorId pkg_colors.ColorId) *AstFactory_JSTypeOrColor {
	this_ := &AstFactory_JSTypeOrColor{}
	InitAstFactory_JSTypeOrColor_2_JSTypeNative_ColorId(this_, jstypeNative, colorId)
	return this_
}

func InitAstFactory_JSTypeOrColor_2_JSTypeNative_ColorId(this_ *AstFactory_JSTypeOrColor, jstypeNative pkg_rhino.JSTypeNative, colorId pkg_colors.ColorId) {
	_ = jstypeNative
	_ = colorId
	this_.initFields_AstFactory_JSTypeOrColor()
	this_.jstypeNative = jstypeNative
	this_.jstype = nil
	this_.color = nil
	this_.colorId = colorId
}

// Java constructor: JSTypeOrColor(com.google.javascript.rhino.jstype.JSTypeNative jstypeNative, com.google.javascript.jscomp.colors.@org.jspecify.annotations.Nullable Color color)
func NewAstFactory_JSTypeOrColor_2_JSTypeNative_Color(jstypeNative pkg_rhino.JSTypeNative, color *pkg_colors.Color) *AstFactory_JSTypeOrColor {
	this_ := &AstFactory_JSTypeOrColor{}
	InitAstFactory_JSTypeOrColor_2_JSTypeNative_Color(this_, jstypeNative, color)
	return this_
}

func InitAstFactory_JSTypeOrColor_2_JSTypeNative_Color(this_ *AstFactory_JSTypeOrColor, jstypeNative pkg_rhino.JSTypeNative, color *pkg_colors.Color) {
	_ = jstypeNative
	_ = color
	this_.initFields_AstFactory_JSTypeOrColor()
	this_.jstypeNative = jstypeNative
	this_.jstype = nil
	this_.color = color
	this_.colorId = (pkg_colors.ColorId{})
}

// Java constructor: JSTypeOrColor(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType jstype, com.google.javascript.jscomp.colors.@org.jspecify.annotations.Nullable Color color)
func NewAstFactory_JSTypeOrColor_2_JSType_Color(jstype pkg_rhino.JSType, color *pkg_colors.Color) *AstFactory_JSTypeOrColor {
	this_ := &AstFactory_JSTypeOrColor{}
	InitAstFactory_JSTypeOrColor_2_JSType_Color(this_, jstype, color)
	return this_
}

func InitAstFactory_JSTypeOrColor_2_JSType_Color(this_ *AstFactory_JSTypeOrColor, jstype pkg_rhino.JSType, color *pkg_colors.Color) {
	_ = jstype
	_ = color
	this_.initFields_AstFactory_JSTypeOrColor()
	this_.jstype = jstype
	this_.jstypeNative = pkg_rhino.JSTypeNative(0)
	this_.color = color
	this_.colorId = (pkg_colors.ColorId{})
}

// Java method: JSTypeOrColor.getJSType(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func (this_ *AstFactory_JSTypeOrColor) GetJSType(registry *pkg_rhino.JSTypeRegistry) pkg_rhino.JSType {
	_ = registry
	if this_.jstype != nil && !this_.jstype.IsJavaNil_() {
		return this_.jstype
	} else {
		return registry.GetNativeType(this_.jstypeNative)
	}
}

// Java method: JSTypeOrColor.getColor(com.google.javascript.jscomp.colors.ColorRegistry registry)
func (this_ *AstFactory_JSTypeOrColor) GetColor(registry *pkg_colors.ColorRegistry) *pkg_colors.Color {
	_ = registry
	if this_.color != nil {
		return this_.color
	} else {
		return registry.Get(j2gort.CheckNotNullVal(this_.colorId))
	}
}
