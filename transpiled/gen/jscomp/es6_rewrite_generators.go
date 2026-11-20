/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/Es6RewriteGenerators.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.Es6RewriteGenerators
type Es6RewriteGenerators struct {
	jscompAsyncExecuteFunction RuntimeJsLibManager_JsLibField             `gson:"jscompAsyncExecuteFunction"`
	jscompAsyncExecuteProgram  RuntimeJsLibManager_JsLibField             `gson:"jscompAsyncExecuteProgram"`
	generatorContext           RuntimeJsLibManager_JsLibField             `gson:"generatorContext"`
	jscompCreateGenerator      RuntimeJsLibManager_JsLibField             `gson:"jscompCreateGenerator"`
	compiler                   AbstractCompiler                           `gson:"compiler"`
	namespace                  pkg_rhino.StaticScope                      `gson:"namespace"`
	astFactory                 *AstFactory                                `gson:"astFactory"`
	nullableStringType         *pkg_colors.Color                          `gson:"nullableStringType"`
	generatorContextType       j2gort.Ext_com_google_common_base_Supplier `gson:"generatorContextType"`
	propertyIteratorType       AstFactory_Type                            `gson:"propertyIteratorType"`
	uniqueIdSupplier           *UniqueIdSupplier                          `gson:"uniqueIdSupplier"`
	idHash_                    atomic.Int32
}

func (this_ *Es6RewriteGenerators) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteGenerators) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteGenerators) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteGenerators", "Es6RewriteGenerators"
}
func (this_ *Es6RewriteGenerators) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteGenerators) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteGenerators", "com.google.javascript.jscomp.Es6RewriteGenerators", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteGenerators](v_)
	})
}
func (*Es6RewriteGenerators) Implements_jscomp_CompilerPass() {}

func StaticGet_Es6RewriteGenerators_GENERATOR_FUNCTION() string { return "$jscomp$generator$function" }
func StaticGet_Es6RewriteGenerators_GENERATOR_CONTEXT() string  { return "$jscomp$generator$context" }
func StaticGet_Es6RewriteGenerators_GENERATOR_ARGUMENTS() string {
	return "$jscomp$generator$arguments"
}
func StaticGet_Es6RewriteGenerators_GENERATOR_THIS() string         { return "$jscomp$generator$this" }
func StaticGet_Es6RewriteGenerators_GENERATOR_FORIN_PREFIX() string { return "$jscomp$generator$forin" }

var static_Es6RewriteGenerators_transpiledFeatures *pkg_parser.FeatureSet

func StaticGet_Es6RewriteGenerators_TranspiledFeatures() *pkg_parser.FeatureSet {
	clinit_Es6RewriteGenerators()
	return static_Es6RewriteGenerators_transpiledFeatures
}
func StaticSet_Es6RewriteGenerators_TranspiledFeatures(v *pkg_parser.FeatureSet) {
	clinit_Es6RewriteGenerators()
	static_Es6RewriteGenerators_transpiledFeatures = v
}

var clinitDone_Es6RewriteGenerators uint32
var clinitState_Es6RewriteGenerators j2gort.ClinitState

func clinit_Es6RewriteGenerators() {
	if atomic.LoadUint32(&clinitDone_Es6RewriteGenerators) != 0 {
		return
	}
	clinitSlow_Es6RewriteGenerators()
}

func clinitSlow_Es6RewriteGenerators() {
	if !j2gort.EnterClinit(&clinitState_Es6RewriteGenerators, &clinitDone_Es6RewriteGenerators) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es6RewriteGenerators, &clinitDone_Es6RewriteGenerators)
	static_Es6RewriteGenerators_transpiledFeatures = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS())
}

func (this_ *Es6RewriteGenerators) initFields_Es6RewriteGenerators() {
	_ = this_
}

// Java constructor: Es6RewriteGenerators(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6RewriteGenerators(compiler AbstractCompiler) *Es6RewriteGenerators {
	clinit_Es6RewriteGenerators()
	this_ := &Es6RewriteGenerators{}
	InitEs6RewriteGenerators(this_, compiler)
	return this_
}

func InitEs6RewriteGenerators(this_ *Es6RewriteGenerators, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_Es6RewriteGenerators()
	j2gort.CheckNotNullIface(compiler)
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.namespace = compiler.GetTranspilationNamespace()
	if compiler.HasOptimizationColors() {
		this_.nullableStringType = pkg_colors.Color_CreateUnion2(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID(), pkg_colors.StaticGet_StandardColors_STRING())
	} else {
		this_.nullableStringType = nil
	}
	var runtimeJsLibManager *RuntimeJsLibManager = compiler.GetRuntimeJsLibManager()
	_ = runtimeJsLibManager
	this_.generatorContext = runtimeJsLibManager.GetJsLibField("$jscomp.generator.Context")
	this_.jscompAsyncExecuteFunction = runtimeJsLibManager.GetJsLibField("$jscomp.asyncExecutePromiseGeneratorFunction")
	this_.jscompAsyncExecuteProgram = runtimeJsLibManager.GetJsLibField("$jscomp.asyncExecutePromiseGeneratorProgram")
	this_.jscompCreateGenerator = runtimeJsLibManager.GetJsLibField("$jscomp.generator.createGenerator")
	this_.generatorContextType = j2gort.Ext_com_google_common_base_Suppliers_Memoize(j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
		return AstFactory_S_Type_1_Node(this_.astFactory.CreateNewNode(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.generatorContext), []pkg_rhino.Node{}))
	}))
	this_.propertyIteratorType = AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_TYPE(), pkg_colors.StaticGet_StandardColors_TOP_OBJECT())
	this_.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
}

// Java method: Es6RewriteGenerators.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6RewriteGenerators) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	TranspilationPasses_ProcessTranspile(this_.compiler, root, StaticGet_Es6RewriteGenerators_TranspiledFeatures(), []NodeTraversal_Callback{NewEs6RewriteGenerators_GeneratorFunctionsTranspiler(this_)})
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, StaticGet_Es6RewriteGenerators_TranspiledFeatures())
}

// Java class: com.google.javascript.jscomp.Es6RewriteGenerators.YieldExposer
type Es6RewriteGenerators_YieldExposer struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer_     *Es6RewriteGenerators
	decomposer *ExpressionDecomposer `gson:"decomposer"`
}

func (this_ *Es6RewriteGenerators_YieldExposer) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteGenerators_YieldExposer) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteGenerators_YieldExposer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteGenerators$YieldExposer", "YieldExposer"
}
func (this_ *Es6RewriteGenerators_YieldExposer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*Es6RewriteGenerators_YieldExposer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteGenerators.YieldExposer", "com.google.javascript.jscomp.Es6RewriteGenerators$YieldExposer", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteGenerators_YieldExposer](v_)
	})
}
func (this_ *Es6RewriteGenerators_YieldExposer) SetOuter_(o *Es6RewriteGenerators) { this_.outer_ = o }
func (this_ *Es6RewriteGenerators_YieldExposer) GetOuter_() *Es6RewriteGenerators {
	return this_.outer_
}
func (*Es6RewriteGenerators_YieldExposer) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*Es6RewriteGenerators_YieldExposer) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *Es6RewriteGenerators_YieldExposer) initFields_Es6RewriteGenerators_YieldExposer() {
	_ = this_
}

// Java constructor: YieldExposer()
func NewEs6RewriteGenerators_YieldExposer(outer_ *Es6RewriteGenerators) *Es6RewriteGenerators_YieldExposer {
	this_ := &Es6RewriteGenerators_YieldExposer{}
	this_.outer_ = outer_
	InitEs6RewriteGenerators_YieldExposer(this_)
	return this_
}

func InitEs6RewriteGenerators_YieldExposer(this_ *Es6RewriteGenerators_YieldExposer) {
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_Es6RewriteGenerators_YieldExposer()
	this_.decomposer = this_.outer_.compiler.CreateDefaultExpressionDecomposer()
}

// Java method: YieldExposer.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteGenerators_YieldExposer) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	n.SetGeneratorMarker(false)
	if n.IsFunction() {
		return false
	}
	if n.IsYield() {
		this_.VisitYield(n)
		return false
	}
	return true
}

// Java method: YieldExposer.visitYield(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_YieldExposer) VisitYield(n pkg_rhino.Node) {
	_ = n
	if n.GetParent().IsExprResult() {
		return
	}
	if this_.decomposer.CanExposeExpression(n) != StaticGet_ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE() {
		this_.decomposer.MaybeExposeExpression(n)
	} else {
		var link string = "https://github.com/google/closure-compiler/wiki/FAQ#i-get-an-undecomposable-expression-error-for-my-yield-or-await-expression-what-do-i-do"
		_ = link
		var suggestion string = "Please rewrite the yield or await as a separate statement."
		_ = suggestion
		var message string = ("Undecomposable expression: " + j2gort.StrValString(suggestion) + "\nSee " + j2gort.StrValString(link))
		_ = message
		this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_TranspilationUtil_CANNOT_CONVERT(), []string{message}))
	}
}

// Java class: com.google.javascript.jscomp.Es6RewriteGenerators.GeneratorFunctionsTranspiler
type Es6RewriteGenerators_GeneratorFunctionsTranspiler struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *Es6RewriteGenerators
}

func (this_ *Es6RewriteGenerators_GeneratorFunctionsTranspiler) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteGenerators_GeneratorFunctionsTranspiler) IsJavaNil_() bool {
	return this_ == nil
}
func (*Es6RewriteGenerators_GeneratorFunctionsTranspiler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteGenerators$GeneratorFunctionsTranspiler", "GeneratorFunctionsTranspiler"
}
func (this_ *Es6RewriteGenerators_GeneratorFunctionsTranspiler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6RewriteGenerators_GeneratorFunctionsTranspiler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteGenerators.GeneratorFunctionsTranspiler", "com.google.javascript.jscomp.Es6RewriteGenerators$GeneratorFunctionsTranspiler", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteGenerators_GeneratorFunctionsTranspiler](v_)
	})
}
func (this_ *Es6RewriteGenerators_GeneratorFunctionsTranspiler) SetOuter_(o *Es6RewriteGenerators) {
	this_.outer_ = o
}
func (this_ *Es6RewriteGenerators_GeneratorFunctionsTranspiler) GetOuter_() *Es6RewriteGenerators {
	return this_.outer_
}
func (*Es6RewriteGenerators_GeneratorFunctionsTranspiler) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*Es6RewriteGenerators_GeneratorFunctionsTranspiler) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *Es6RewriteGenerators_GeneratorFunctionsTranspiler) initFields_Es6RewriteGenerators_GeneratorFunctionsTranspiler() {
	_ = this_
}

// Java constructor: GeneratorFunctionsTranspiler()
func NewEs6RewriteGenerators_GeneratorFunctionsTranspiler(outer_ *Es6RewriteGenerators) *Es6RewriteGenerators_GeneratorFunctionsTranspiler {
	this_ := &Es6RewriteGenerators_GeneratorFunctionsTranspiler{}
	this_.outer_ = outer_
	InitEs6RewriteGenerators_GeneratorFunctionsTranspiler(this_)
	return this_
}

func InitEs6RewriteGenerators_GeneratorFunctionsTranspiler(this_ *Es6RewriteGenerators_GeneratorFunctionsTranspiler) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6RewriteGenerators_GeneratorFunctionsTranspiler()
}

// Java method: GeneratorFunctionsTranspiler.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteGenerators_GeneratorFunctionsTranspiler) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsGeneratorFunction() {
		{
			var val_1 Es6RewriteGenerators_SingleGeneratorFunctionTranspiler
			val_1.outer_ = this_.outer_
			InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler(&val_1, n, this_.outer_.uniqueIdSupplier.GetUniqueId(t.GetInput()))
			(&val_1).Transpile()
		}
	}
}

// Java class: com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler struct {
	outer_                 *Es6RewriteGenerators
	uniqueId               string                                                                       `gson:"uniqueId"`
	context                *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext `gson:"context"`
	originalGeneratorBody  pkg_rhino.Node                                                               `gson:"originalGeneratorBody"`
	forInCounter           int32                                                                        `gson:"forInCounter"`
	newGeneratorHoistBlock pkg_rhino.Node                                                               `gson:"newGeneratorHoistBlock"`
	idHash_                atomic.Int32
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) IsJavaNil_() bool {
	return this_ == nil
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler", "SingleGeneratorFunctionTranspiler"
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler", "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler](v_)
	})
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) SetOuter_(o *Es6RewriteGenerators) {
	this_.outer_ = o
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) GetOuter_() *Es6RewriteGenerators {
	return this_.outer_
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler() {
	_ = this_
	this_.uniqueId = j2gort.NullString

}

// Java constructor: SingleGeneratorFunctionTranspiler(com.google.javascript.rhino.Node genFunc, java.lang.String uniqueId)
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler(outer_ *Es6RewriteGenerators, genFunc pkg_rhino.Node, uniqueId string) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler {
	this_ := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler{}
	this_.outer_ = outer_
	InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler(this_, genFunc, uniqueId)
	return this_
}

func InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler(this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler, genFunc pkg_rhino.Node, uniqueId string) {
	_ = genFunc
	_ = uniqueId
	this_.initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler()
	this_.originalGeneratorBody = genFunc.GetLastChild()
	this_.uniqueId = uniqueId
	this_.context = NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext(this_)
}

// Java method: SingleGeneratorFunctionTranspiler.hoistVarNode(com.google.javascript.rhino.Node node)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) priv_SingleGeneratorFunctionTranspiler_hoistVarNode_2661(node pkg_rhino.Node) {
	_ = node
	if !(node.IsVar()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
	}
	node.InsertBefore(this_.newGeneratorHoistBlock.GetLastChild())
}

// Java method: SingleGeneratorFunctionTranspiler.hoistFunctionDeclarationNode(com.google.javascript.rhino.Node node)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) priv_SingleGeneratorFunctionTranspiler_hoistFunctionDeclarationNode_2662(node pkg_rhino.Node) {
	_ = node
	if !(node.IsFunction()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
	}
	this_.newGeneratorHoistBlock.AddChildToFront(node)
}

// Java method: SingleGeneratorFunctionTranspiler.isTranspiledAsyncFunction(com.google.javascript.rhino.Node generatorFunction)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) priv_SingleGeneratorFunctionTranspiler_isTranspiledAsyncFunction_2663(generatorFunction pkg_rhino.Node) bool {
	_ = generatorFunction
	if generatorFunction.GetParent().IsCall() && (generatorFunction.GetPrevious() != 0) {
		var callTarget pkg_rhino.Node = generatorFunction.GetParent().GetFirstChild()
		_ = callTarget
		if ((generatorFunction.GetPrevious() == callTarget) && (generatorFunction.GetNext() == 0)) && this_.outer_.jscompAsyncExecuteFunction.Matches(callTarget) {
			if !(generatorFunction.GetGrandparent().IsReturn()) {
				j2gort.CheckState(false)
			}
			if !(generatorFunction.GetGrandparent().GetNext() == 0) {
				j2gort.CheckState(false)
			}
			return true
		}
	}
	return false
}

// Java method: SingleGeneratorFunctionTranspiler.transpile()
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) Transpile() {
	var generatorFunction pkg_rhino.Node = this_.originalGeneratorBody.GetParent()
	_ = generatorFunction
	if !(generatorFunction.IsGeneratorFunction()) {
		j2gort.CheckState(false)
	}
	generatorFunction.PutBooleanProp(pkg_rhino.StaticGet_Node_GENERATOR_FN(), false)
	var program pkg_rhino.Node = pkg_rhino.Node(0)
	_ = program
	var programType *pkg_colors.Color = pkg_colors.StaticGet_StandardColors_TOP_OBJECT()
	_ = programType
	var generatorBody pkg_rhino.Node = this_.outer_.astFactory.CreateBlock([]pkg_rhino.Node{})
	_ = generatorBody
	var changeScopeNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = changeScopeNode
	if this_.priv_SingleGeneratorFunctionTranspiler_isTranspiledAsyncFunction_2663(generatorFunction) {
		var callTarget pkg_rhino.Node = generatorFunction.GetPrevious()
		_ = callTarget
		if !(callTarget.IsGetProp() || callTarget.IsName()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(callTarget))
		}
		this_.newGeneratorHoistBlock = generatorFunction.GetGrandparent().GetParent()
		if !(this_.newGeneratorHoistBlock.IsBlock()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(this_.newGeneratorHoistBlock))
		}
		changeScopeNode = NodeUtil_GetEnclosingFunction(this_.newGeneratorHoistBlock)
		if !(changeScopeNode.IsFunction()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(changeScopeNode))
		}
		callTarget.ReplaceWith(this_.outer_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.outer_.namespace, this_.outer_.jscompAsyncExecuteProgram).SrcrefTree(callTarget))
		program = this_.originalGeneratorBody.GetParent()
		this_.originalGeneratorBody.GetPrevious().AddChildToBack(this_.context.GetJsContextNameNode(this_.originalGeneratorBody))
		this_.originalGeneratorBody.ReplaceWith(generatorBody)
	} else {
		changeScopeNode = generatorFunction
		var genFuncName pkg_rhino.Node = generatorFunction.GetFirstChild()
		_ = genFuncName
		if !(genFuncName.IsName()) {
			j2gort.CheckState(false)
		}
		if j2gort.CheckNotNullString(genFuncName.GetString()) == "" {
			genFuncName.SetString(this_.context.GetScopedName("$jscomp$generator$function"))
			if this_.outer_.astFactory.IsAddingColors() {
				genFuncName.SetColor(generatorFunction.GetColor())
			}
			genFuncName.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
		}
		program = this_.outer_.astFactory.CreateFunction("", pkg_rhino.IR_ParamList([]pkg_rhino.Node{this_.context.GetJsContextNameNode(this_.originalGeneratorBody)}), generatorBody, AstFactory_S_Type_1_Color(programType))
		var createGenerator pkg_rhino.Node = this_.outer_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.outer_.namespace, this_.outer_.jscompCreateGenerator)
		_ = createGenerator
		this_.newGeneratorHoistBlock = this_.outer_.astFactory.CreateBlock([]pkg_rhino.Node{this_.outer_.astFactory.CreateReturn(this_.outer_.astFactory.CreateCallWithUnknownType(createGenerator, []pkg_rhino.Node{genFuncName.CloneNode_0(), program})).SrcrefTree(this_.originalGeneratorBody)})
		this_.originalGeneratorBody.ReplaceWith(this_.newGeneratorHoistBlock)
	}
	this_.outer_.compiler.ReportChangeToChangeScope(program)
	NodeTraversal_Traverse(this_.outer_.compiler, this_.originalGeneratorBody, NewEs6RewriteGenerators_YieldNodeMarker())
	var shouldAddFinalJump bool = !(this_.priv_SingleGeneratorFunctionTranspiler_isEndOfBlockUnreachable_2685(this_.originalGeneratorBody))
	_ = shouldAddFinalJump
	for this_.originalGeneratorBody.HasChildren() {
		this_.TranspileStatement_1(this_.originalGeneratorBody.RemoveFirstChild())
	}
	var finalBlock pkg_rhino.Node = this_.outer_.astFactory.CreateBlock([]pkg_rhino.Node{})
	_ = finalBlock
	if shouldAddFinalJump {
		finalBlock.AddChildToBack(this_.context.CallContextMethodResult(this_.originalGeneratorBody, "jumpToEnd", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), []pkg_rhino.Node{}))
	}
	this_.context.currentCase.JumpTo(this_.context.programEndCase, finalBlock)
	{
		r_2 := this_.context.currentCase
		r_2.mayFallThrough = true
	}
	this_.context.FinalizeTransformation(generatorBody)
	this_.context.CheckStateIsEmpty()
	this_.outer_.compiler.ReportChangeToChangeScope(changeScopeNode)
}

// Java method: SingleGeneratorFunctionTranspiler.transpileStatement(com.google.javascript.rhino.Node statement)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileStatement_1(statement pkg_rhino.Node) {
	_ = statement
	this_.TranspileStatement_3(statement, nil, nil)
}

// Java method: SingleGeneratorFunctionTranspiler.transpileStatement(com.google.javascript.rhino.Node statement, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case breakCase, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case continueCase)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileStatement_3(statement pkg_rhino.Node, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = statement
	_ = breakCase
	_ = continueCase
	if !(pkg_rhino.IR_MayBeStatement(statement)) {
		j2gort.CheckState(false)
	}
	if !(statement.GetParent() == 0) {
		j2gort.CheckState(false)
	}
	if !(statement.IsGeneratorMarker()) {
		this_.TranspileUnmarkedNode(statement)
		return
	}
	switch statement.TokenSwitchOrd_() {
	case int32(92) /* Token.LABEL */ :
		this_.TranspileLabel(statement)
	case int32(90) /* Token.BLOCK */ :
		this_.TranspileBlock(statement)
	case int32(93) /* Token.EXPR_RESULT */ :
		this_.TranspileExpressionResult(statement)
	case int32(84) /* Token.VAR */ :
		this_.TranspileVar(statement)
	case int32(0) /* Token.RETURN */ :
		this_.TranspileReturn(statement)
	case int32(43) /* Token.THROW */ :
		this_.TranspileThrow(statement)
	case int32(74) /* Token.IF */ :
		this_.TranspileIf(statement, breakCase)
	case int32(80) /* Token.FOR */ :
		this_.TranspileFor(statement, breakCase, continueCase)
	case int32(81) /* Token.FOR_IN */ :
		this_.TranspileForIn(statement, breakCase, continueCase)
	case int32(79) /* Token.DO */ :
		this_.TranspileDo(statement, breakCase, continueCase)
	case int32(48) /* Token.TRY */ :
		this_.TranspileTry(statement, breakCase)
	case int32(75) /* Token.SWITCH */ :
		this_.TranspileSwitch(statement, breakCase)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unsupported token: "+j2gort.StrVal(statement.GetToken()))))
	}
}

// Java method: SingleGeneratorFunctionTranspiler.transpileUnmarkedNode(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileUnmarkedNode(n pkg_rhino.Node) {
	_ = n
	if !(!(n.IsGeneratorMarker())) {
		j2gort.CheckState(false)
	}
	if n.IsFunction() {
		var functionName string = n.GetFirstChild().GetString()
		_ = functionName
		if !(!(j2gort.CheckNotNullString(functionName) == "") && !(javalang.StringStartsWith(j2gort.CheckNotNullString(functionName), "$jscomp$generator$function"))) {
			j2gort.CheckState(false)
		}
		this_.priv_SingleGeneratorFunctionTranspiler_hoistFunctionDeclarationNode_2662(n)
		return
	}
	this_.context.TranspileUnmarkedBlock(func() pkg_rhino.Node {
		if n.IsBlock() || n.IsAddedBlock() {
			return n
		}
		return pkg_rhino.IR_Block_1_Node(n)
	}())
}

// Java method: SingleGeneratorFunctionTranspiler.transpileLabel(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileLabel(n pkg_rhino.Node) {
	_ = n
	var labelNames j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = labelNames
	for n.IsLabel() {
		labelNames.Add_1(j2gort.BoxNode(n.RemoveFirstChild()))
		n = n.RemoveFirstChild()
	}
	var continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
	_ = continueCase
	if NodeUtil_IsLoopStructure(n) {
		continueCase = this_.context.CreateCase()
	} else {
		continueCase = nil
	}
	var breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.context.CreateCase()
	_ = breakCase
	this_.context.PushLabels(labelNames, breakCase, continueCase)
	this_.TranspileStatement_3(n, breakCase, continueCase)
	this_.context.PopLabels(labelNames)
	if breakCase != this_.context.currentCase {
		this_.context.SwitchCaseTo(breakCase)
	}
}

// Java method: SingleGeneratorFunctionTranspiler.transpileBlock(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileBlock(n pkg_rhino.Node) {
	_ = n
	for n.HasChildren() {
		this_.TranspileStatement_1(n.RemoveFirstChild())
	}
}

// Java method: SingleGeneratorFunctionTranspiler.transpileExpressionResult(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileExpressionResult(n pkg_rhino.Node) {
	_ = n
	var exposedExpression pkg_rhino.Node = this_.ExposeYieldAndTranspileRest(n.RemoveFirstChild())
	_ = exposedExpression
	var decomposed pkg_rhino.Node = this_.TranspileYields(exposedExpression)
	_ = decomposed
	if !(exposedExpression.IsYield()) {
		n.AddChildToFront(this_.PrepareNodeForWrite(decomposed))
		n.SetGeneratorMarker(false)
		this_.context.WriteGeneratedNode(n)
	}
}

// Java method: SingleGeneratorFunctionTranspiler.transpileVar(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileVar(n pkg_rhino.Node) {
	_ = n
	n.SetGeneratorMarker(false)
	var newVars pkg_rhino.Node = n.CloneNode_0()
	_ = newVars
	for n.HasChildren() {
		var var_ pkg_rhino.Node = pkg_rhino.Node(0)
		_ = var_
		for ((func() pkg_rhino.Node { var asgn_3 pkg_rhino.Node = n.RemoveFirstChild(); var_ = asgn_3; return asgn_3 }()) != 0) && !(var_.IsGeneratorMarker()) {
			newVars.AddChildToBack(var_)
		}
		if newVars.HasChildren() {
			this_.TranspileUnmarkedNode(newVars)
			newVars = n.CloneNode_0()
		}
		if var_ != 0 {
			if !(var_.IsGeneratorMarker()) {
				j2gort.CheckState(false)
			}
			var_.AddChildToFront(this_.MaybeDecomposeExpression(var_.RemoveFirstChild()))
			var_.SetGeneratorMarker(false)
			newVars.AddChildToBack(var_)
		}
	}
	if newVars.HasChildren() {
		this_.TranspileUnmarkedNode(newVars)
	}
}

// Java method: SingleGeneratorFunctionTranspiler.transpileReturn(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileReturn(n pkg_rhino.Node) {
	_ = n
	n.AddChildToFront(this_.context.ReturnExpression(n, this_.PrepareNodeForWrite(this_.MaybeDecomposeExpression(n.RemoveFirstChild()))))
	this_.context.WriteGeneratedNode(n)
	{
		r_4 := this_.context.currentCase
		r_4.mayFallThrough = false
	}
}

// Java method: SingleGeneratorFunctionTranspiler.transpileThrow(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileThrow(n pkg_rhino.Node) {
	_ = n
	n.AddChildToFront(this_.PrepareNodeForWrite(this_.MaybeDecomposeExpression(n.RemoveFirstChild())))
	this_.context.WriteGeneratedNode(n)
	{
		r_5 := this_.context.currentCase
		r_5.mayFallThrough = false
	}
}

// Java method: SingleGeneratorFunctionTranspiler.exposeYieldAndTranspileRest(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) ExposeYieldAndTranspileRest(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsGeneratorMarker()) {
		j2gort.CheckState(false)
	}
	if n.IsYield() {
		return n
	}
	var isExpression bool = pkg_rhino.IR_MayBeExpression(n)
	_ = isExpression
	var block pkg_rhino.Node = this_.outer_.astFactory.CreateBlock([]pkg_rhino.Node{func() pkg_rhino.Node {
		if isExpression {
			return this_.outer_.astFactory.CreateReturn(n)
		}
		return n
	}()})
	_ = block
	NodeTraversal_Traverse(this_.outer_.compiler, n, NewEs6RewriteGenerators_YieldExposer(this_.outer_))
	NodeTraversal_Traverse(this_.outer_.compiler, block, NewEs6RewriteGenerators_YieldNodeMarker())
	var decomposed pkg_rhino.Node = block.GetLastChild().Detach()
	_ = decomposed
	this_.TranspileStatement_1(block)
	if isExpression {
		return decomposed.RemoveFirstChild()
	} else {
		return decomposed
	}
}

// Java method: SingleGeneratorFunctionTranspiler.maybeDecomposeExpression(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) MaybeDecomposeExpression(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if (n == 0) || !(n.IsGeneratorMarker()) {
		return n
	}
	return this_.TranspileYields(this_.ExposeYieldAndTranspileRest(n))
}

// Java method: SingleGeneratorFunctionTranspiler.prepareNodeForWrite(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) PrepareNodeForWrite(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if n == 0 {
		return pkg_rhino.Node(0)
	}
	var wrapper pkg_rhino.Node
	_ = wrapper
	if pkg_rhino.IR_MayBeStatement(n) {
		wrapper = this_.outer_.astFactory.CreateBlock([]pkg_rhino.Node{n})
	} else {
		wrapper = this_.outer_.astFactory.ExprResult(n)
	}
	NodeTraversal_Traverse(this_.outer_.compiler, wrapper, NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler(this_.context))
	if !(wrapper.HasOneChild()) {
		j2gort.CheckState(false)
	}
	return wrapper.RemoveFirstChild()
}

// Java method: SingleGeneratorFunctionTranspiler.transpileYields(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileYields(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsGeneratorMarker()) {
		return n
	}
	var jumpToSection *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.context.CreateCase()
	_ = jumpToSection
	var yieldNode pkg_rhino.Node = this_.FindYield(n)
	_ = yieldNode
	var yieldExpression pkg_rhino.Node = this_.PrepareNodeForWrite(this_.MaybeDecomposeExpression(yieldNode.RemoveFirstChild()))
	_ = yieldExpression
	if yieldNode.IsYieldAll() {
		this_.context.YieldAll(yieldExpression, jumpToSection, yieldNode)
	} else {
		this_.context.YieldValue(yieldExpression, jumpToSection, yieldNode)
	}
	this_.context.SwitchCaseTo(jumpToSection)
	var yieldResult pkg_rhino.Node = this_.context.YieldResult(yieldNode)
	_ = yieldResult
	if yieldNode == n {
		return yieldResult
	}
	yieldNode.ReplaceWith(yieldResult)
	for yieldResult != n {
		yieldResult = yieldResult.GetParent()
		yieldResult.SetGeneratorMarker(false)
	}
	return n
}

// Java method: SingleGeneratorFunctionTranspiler.transpileIf(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case breakCase)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileIf(n pkg_rhino.Node, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = n
	_ = breakCase
	var condition pkg_rhino.Node = this_.MaybeDecomposeExpression(n.RemoveFirstChild())
	_ = condition
	var ifBlock pkg_rhino.Node = n.GetFirstChild()
	_ = ifBlock
	var elseBlock pkg_rhino.Node = ifBlock.GetNext()
	_ = elseBlock
	if !(ifBlock.IsGeneratorMarker()) && ((elseBlock == 0) || !(elseBlock.IsGeneratorMarker())) {
		n.AddChildToFront(condition)
		n.SetGeneratorMarker(false)
		this_.TranspileUnmarkedNode(n)
		return
	}
	ifBlock.Detach()
	if elseBlock == 0 {
		elseBlock = this_.outer_.astFactory.CreateBlock([]pkg_rhino.Node{}).Srcref(n)
	} else {
		elseBlock.Detach()
	}
	if ifBlock.IsGeneratorMarker() && !(elseBlock.IsGeneratorMarker()) {
		condition = this_.outer_.astFactory.CreateNot(condition).Srcref(condition)
		var tmpNode pkg_rhino.Node = ifBlock
		_ = tmpNode
		ifBlock = elseBlock
		elseBlock = tmpNode
	}
	if !(ifBlock.IsGeneratorMarker()) {
		var endCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.context.MaybeCreateCase(breakCase)
		_ = endCase
		var jumpToBlock pkg_rhino.Node = this_.context.CreateJumpToBlock(endCase, false, ifBlock)
		_ = jumpToBlock
		for jumpToBlock.HasChildren() {
			var jumpToNode pkg_rhino.Node = jumpToBlock.RemoveFirstChild()
			_ = jumpToNode
			jumpToNode.SetGeneratorSafe(true)
			ifBlock.AddChildToBack(jumpToNode)
		}
		this_.TranspileUnmarkedNode(this_.outer_.astFactory.CreateIf_2(condition, ifBlock).Srcref(n))
		this_.TranspileStatement_1(elseBlock)
		this_.context.SwitchCaseTo(endCase)
		return
	}
	var ifCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.context.CreateCase()
	_ = ifCase
	var endCase_2 *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.context.MaybeCreateCase(breakCase)
	_ = endCase_2
	condition = this_.PrepareNodeForWrite(condition)
	var newIfBlock pkg_rhino.Node = this_.context.CreateJumpToBlock(ifCase, true, n)
	_ = newIfBlock
	this_.context.WriteGeneratedNode(this_.outer_.astFactory.CreateIf_2(this_.PrepareNodeForWrite(condition), newIfBlock).Srcref(n))
	this_.TranspileStatement_1(elseBlock)
	this_.context.WriteJumpTo(endCase_2, elseBlock)
	this_.context.SwitchCaseTo(ifCase)
	this_.TranspileStatement_1(ifBlock)
	this_.context.SwitchCaseTo(endCase_2)
}

// Java method: SingleGeneratorFunctionTranspiler.transpileFor(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case breakCase, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case continueCase)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileFor(n pkg_rhino.Node, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = n
	_ = breakCase
	_ = continueCase
	var init pkg_rhino.Node = this_.MaybeDecomposeExpression(n.RemoveFirstChild())
	_ = init
	var condition pkg_rhino.Node = n.GetFirstChild()
	_ = condition
	var increment pkg_rhino.Node = condition.GetNext()
	_ = increment
	var body pkg_rhino.Node = increment.GetNext()
	_ = body
	if (!(condition.IsGeneratorMarker()) && !(increment.IsGeneratorMarker())) && !(body.IsGeneratorMarker()) {
		n.AddChildToFront(init)
		n.SetGeneratorMarker(false)
		this_.TranspileUnmarkedNode(n)
		return
	}
	if !(init.IsEmpty()) {
		if pkg_rhino.IR_MayBeExpression(init) {
			init = this_.outer_.astFactory.ExprResult(init).Srcref(init)
		}
		this_.TranspileUnmarkedNode(init)
	}
	var startCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.context.CreateCase()
	_ = startCase
	var incrementCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.context.MaybeCreateCase(continueCase)
	_ = incrementCase
	var endCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.context.MaybeCreateCase(breakCase)
	_ = endCase
	this_.context.SwitchCaseTo(startCase)
	if !(condition.IsEmpty()) {
		condition = this_.PrepareNodeForWrite(this_.MaybeDecomposeExpression(condition.Detach()))
		this_.context.WriteGeneratedNode(this_.outer_.astFactory.CreateIf_2(this_.outer_.astFactory.CreateNot(condition).Srcref(condition), this_.context.CreateJumpToBlock(endCase, true, n)).Srcref(n))
	}
	this_.context.PushBreakContinueContext(endCase, incrementCase)
	this_.TranspileStatement_1(body.Detach())
	this_.context.PopBreakContinueContext()
	this_.context.SwitchCaseTo(incrementCase)
	if !(increment.IsEmpty()) {
		increment = this_.MaybeDecomposeExpression(increment.Detach())
		this_.TranspileUnmarkedNode(this_.outer_.astFactory.ExprResult(increment).Srcref(increment))
	}
	this_.context.WriteJumpTo(startCase, n)
	this_.context.SwitchCaseTo(endCase)
}

// Java method: SingleGeneratorFunctionTranspiler.transpileForIn(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case breakCase, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case continueCase)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileForIn(n pkg_rhino.Node, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = n
	_ = breakCase
	_ = continueCase
	var detachedExpr pkg_rhino.Node = this_.MaybeDecomposeExpression(n.GetSecondChild().Detach())
	_ = detachedExpr
	var target pkg_rhino.Node = n.GetFirstChild()
	_ = target
	var body pkg_rhino.Node = n.GetSecondChild()
	_ = body
	if !(target.IsGeneratorMarker()) && !(body.IsGeneratorMarker()) {
		detachedExpr.InsertAfter(target)
		n.SetGeneratorMarker(false)
		this_.TranspileUnmarkedNode(n)
		return
	}
	var init pkg_rhino.Node = pkg_rhino.Node(0)
	_ = init
	if target.Detach().IsVar() {
		if !(!(target.IsGeneratorMarker())) {
			j2gort.CheckState(false)
		}
		init = target
		if !(!(init.GetFirstChild().HasChildren())) {
			j2gort.CheckState(false)
		}
		target = init.GetFirstChild().CloneNode_0()
	} else {
		init = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_VAR()).Srcref(target)
	}
	var child pkg_rhino.Node = this_.context.CallContextMethod(target, "forIn", this_.outer_.propertyIteratorType, []pkg_rhino.Node{detachedExpr})
	_ = child
	var forIn pkg_rhino.Node = this_.outer_.astFactory.CreateName_2_string_AstFactory_Type((j2gort.StrValString(this_.context.GetScopedName("$jscomp$generator$forin")) + "$" + j2gort.StrValInt32(func() int32 { old_6 := this_.forInCounter; nw_7 := old_6 + 1; this_.forInCounter = nw_7; return old_6 }())), AstFactory_S_Type_1_Node(child)).Srcref(target)
	_ = forIn
	forIn.AddChildToFront(child)
	init.AddChildToBack(forIn)
	var forInGetNext pkg_rhino.Node = this_.outer_.astFactory.CreateGetProp(forIn.CloneNode_0(), "getNext", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT())).Srcref(detachedExpr)
	_ = forInGetNext
	var forCond pkg_rhino.Node = this_.outer_.astFactory.CreateNe(this_.outer_.astFactory.CreateAssign_2_Node_Node(target.SetColor(this_.outer_.nullableStringType), this_.outer_.astFactory.CreateCall(forInGetNext, AstFactory_S_Type_1_Color(this_.outer_.nullableStringType), []pkg_rhino.Node{}).Srcref(detachedExpr)).Srcref(detachedExpr), this_.outer_.astFactory.CreateNull().Srcref(forIn)).Srcref(detachedExpr)
	_ = forCond
	forCond.SetGeneratorMarker(target.IsGeneratorMarker())
	var forNode pkg_rhino.Node = pkg_rhino.IR_ForNode(init, forCond, pkg_rhino.IR_Empty().Srcref(n), body.Detach()).Srcref(n)
	_ = forNode
	this_.TranspileFor(forNode, breakCase, continueCase)
}

// Java method: SingleGeneratorFunctionTranspiler.transpileDo(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case breakCase, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case continueCase)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileDo(n pkg_rhino.Node, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = n
	_ = breakCase
	_ = continueCase
	var startCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.context.CreateCase()
	_ = startCase
	breakCase = this_.context.MaybeCreateCase(breakCase)
	continueCase = this_.context.MaybeCreateCase(continueCase)
	this_.context.SwitchCaseTo(startCase)
	var body pkg_rhino.Node = n.RemoveFirstChild()
	_ = body
	this_.context.PushBreakContinueContext(breakCase, continueCase)
	this_.TranspileStatement_1(body)
	this_.context.PopBreakContinueContext()
	this_.context.SwitchCaseTo(continueCase)
	var condition pkg_rhino.Node = this_.PrepareNodeForWrite(this_.MaybeDecomposeExpression(n.RemoveFirstChild()))
	_ = condition
	this_.context.WriteGeneratedNode(this_.outer_.astFactory.CreateIf_2(condition, this_.context.CreateJumpToBlock(startCase, false, n)).Srcref(n))
	this_.context.SwitchCaseTo(breakCase)
}

// Java method: SingleGeneratorFunctionTranspiler.transpileTry(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case breakCase)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileTry(n pkg_rhino.Node, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = n
	_ = breakCase
	var tryBlock pkg_rhino.Node = n.RemoveFirstChild()
	_ = tryBlock
	var catchBlock pkg_rhino.Node = n.RemoveFirstChild()
	_ = catchBlock
	var finallyBlock pkg_rhino.Node = n.RemoveFirstChild()
	_ = finallyBlock
	var catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
	_ = catchCase
	if catchBlock.HasChildren() {
		catchCase = this_.context.CreateCase()
	} else {
		catchCase = nil
	}
	var finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
	_ = finallyCase
	if finallyBlock == 0 {
		finallyCase = nil
	} else {
		finallyCase = this_.context.CreateCase()
	}
	var endCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.context.MaybeCreateCase(breakCase)
	_ = endCase
	this_.context.EnterTryBlock(catchCase, finallyCase, tryBlock)
	this_.TranspileStatement_1(tryBlock)
	if finallyBlock == 0 {
		this_.context.LeaveTryBlock(catchCase, endCase, tryBlock)
	} else {
		this_.context.SwitchCaseTo(finallyCase)
		this_.context.EnterFinallyBlock(catchCase, finallyCase, finallyBlock)
		this_.TranspileStatement_1(finallyBlock)
		this_.context.LeaveFinallyBlock(endCase, finallyBlock)
	}
	if catchBlock.HasChildren() {
		if !(catchBlock.GetFirstChild().IsCatch()) {
			j2gort.CheckState(false)
		}
		this_.context.SwitchCaseTo(catchCase)
		var exceptionName pkg_rhino.Node = catchBlock.GetFirstFirstChild().Detach()
		_ = exceptionName
		this_.context.EnterCatchBlock(finallyCase, exceptionName)
		var catchBody pkg_rhino.Node = catchBlock.GetFirstFirstChild().Detach()
		_ = catchBody
		if !(catchBody.IsBlock()) {
			j2gort.CheckState(false)
		}
		this_.TranspileStatement_1(catchBody)
		this_.context.LeaveCatchBlock(finallyCase, catchBody)
	}
	this_.context.SwitchCaseTo(endCase)
}

// Java method: SingleGeneratorFunctionTranspiler.transpileSwitch(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case breakCase)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileSwitch(n pkg_rhino.Node, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = n
	_ = breakCase
	n.AddChildToFront(this_.MaybeDecomposeExpression(n.RemoveFirstChild()))
	var hasGeneratorMarker bool = false
	_ = hasGeneratorMarker
	var switchBody pkg_rhino.Node = n.GetSecondChild()
	_ = switchBody
	{
		var caseSection pkg_rhino.Node = switchBody.GetFirstChild()
		_ = caseSection
		for ; caseSection != 0; caseSection = caseSection.GetNext() {
			if caseSection.IsGeneratorMarker() {
				hasGeneratorMarker = true
				break
			}
		}
	}
	if !(hasGeneratorMarker) {
		n.SetGeneratorMarker(false)
		this_.TranspileUnmarkedNode(n)
		return
	}
	var detachedCases []*Es6RewriteGenerators__Local_42_SwitchCase = []*Es6RewriteGenerators__Local_42_SwitchCase{}
	_ = detachedCases
	var canSkipUnmarkedCases bool = true
	_ = canSkipUnmarkedCases
	{
		var caseSection_2 pkg_rhino.Node = switchBody.GetFirstChild()
		_ = caseSection_2
		for ; caseSection_2 != 0; caseSection_2 = caseSection_2.GetNext() {
			if !(caseSection_2.IsDefaultCase()) && caseSection_2.GetFirstChild().IsGeneratorMarker() {
				this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_TranspilationUtil_CANNOT_CONVERT_YET(), []string{"Case statements that contain yields"}))
				return
			}
			var body pkg_rhino.Node = caseSection_2.GetLastChild()
			_ = body
			if !(body.HasChildren()) || (canSkipUnmarkedCases && !(body.IsGeneratorMarker())) {
				continue
			}
			canSkipUnmarkedCases = this_.priv_SingleGeneratorFunctionTranspiler_isEndOfBlockUnreachable_2685(body)
			var generatedCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.context.CreateCase()
			_ = generatedCase
			generatedCase.caseBlock.Srcref(body)
			var newBody pkg_rhino.Node = pkg_rhino.IR_Block_1_Node(this_.context.CreateJumpToNode(generatedCase, body))
			_ = newBody
			newBody.SetIsAddedBlock(true)
			newBody.SetGeneratorSafe(true)
			body.ReplaceWith(newBody)
			detachedCases = append(detachedCases, NewEs6RewriteGenerators__Local_42_SwitchCase(this_, generatedCase, body))
			caseSection_2.SetGeneratorMarker(false)
		}
	}
	var endCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.context.MaybeCreateCase(breakCase)
	_ = endCase
	n.SetGeneratorMarker(false)
	switchBody.SetGeneratorMarker(false)
	this_.TranspileUnmarkedNode(n)
	this_.context.WriteJumpTo(endCase, n)
	this_.context.PushBreakContext(endCase)
	for _, detachedCase := range detachedCases {
		_ = detachedCase
		var generatedCase_2 *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = detachedCase.generatedCase
		_ = generatedCase_2
		this_.context.SwitchCaseTo(generatedCase_2)
		this_.TranspileStatement_1(detachedCase.body)
	}

	this_.context.PopBreakContext()
	this_.context.SwitchCaseTo(endCase)
}

// Java method: SingleGeneratorFunctionTranspiler.findYield(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) FindYield(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var yieldFinder *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder = NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder()
	_ = yieldFinder
	NodeTraversal_Traverse(this_.outer_.compiler, n, yieldFinder)
	return yieldFinder.GetYieldNode()
}

// Java method: SingleGeneratorFunctionTranspiler.isEndOfBlockUnreachable(com.google.javascript.rhino.Node block)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) priv_SingleGeneratorFunctionTranspiler_isEndOfBlockUnreachable_2685(block pkg_rhino.Node) bool {
	_ = block
	if !(block.IsBlock()) {
		j2gort.CheckState(false)
	}
	if !(block.HasChildren()) {
		return false
	}
	switch block.GetLastChild().TokenSwitchOrd_() {
	case int32(90) /* Token.BLOCK */ :
		return this_.priv_SingleGeneratorFunctionTranspiler_isEndOfBlockUnreachable_2685(block.GetLastChild())
	case int32(0) /* Token.RETURN */, int32(43) /* Token.THROW */, int32(83) /* Token.CONTINUE */, int32(82) /* Token.BREAK */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java class: com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.YieldFinder
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	yieldNode pkg_rhino.Node `gson:"yieldNode"`
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) IsJavaNil_() bool {
	return this_ == nil
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$YieldFinder", "YieldFinder"
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.YieldFinder", "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$YieldFinder", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder](v_)
	})
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder() {
	_ = this_
}

// Java constructor: YieldFinder()
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder {
	this_ := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder{}
	InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder(this_)
	return this_
}

func InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder(this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) {
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder()
}

// Java method: YieldFinder.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if n.IsFunction() {
		return false
	}
	if n.IsYield() {
		if !(this_.yieldNode == 0) {
			j2gort.CheckState(false)
		}
		this_.yieldNode = n
		return false
	}
	return true
}

// Java method: YieldFinder.getYieldNode()
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) GetYieldNode() pkg_rhino.Node {
	j2gort.CheckNotNullNode(this_.yieldNode)
	return this_.yieldNode
}

// Java class: com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext struct {
	outer_                  *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler
	namedLabels             j2gort.OrderedMap[string, *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases] `gson:"namedLabels"`
	breakCases              j2gort.Deque[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case]                    `gson:"breakCases"`
	continueCases           j2gort.Deque[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case]                    `gson:"continueCases"`
	catchCases              j2gort.Deque[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase]               `gson:"catchCases"`
	finallyCases            j2gort.Deque[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case]                    `gson:"finallyCases"`
	catchNames              j2gort.OrderedSet[string]                                                                                          `gson:"catchNames"`
	allCases                j2gort.Ext_java_util_ArrayList                                                                                     `gson:"allCases"`
	switchBreaks            []pkg_rhino.Node                                                                                                   `gson:"switchBreaks"`
	programEndCase          *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case                                  `gson:"programEndCase"`
	currentCase             *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case                                  `gson:"currentCase"`
	caseIdCounter           int32                                                                                                              `gson:"caseIdCounter"`
	nestedFinallyBlockCount int32                                                                                                              `gson:"nestedFinallyBlockCount"`
	idHash_                 atomic.Int32
	thisReferenceFound      bool `gson:"thisReferenceFound"`
	argumentsReferenceFound bool `gson:"argumentsReferenceFound"`
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) IsJavaNil_() bool {
	return this_ == nil
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$TranspilationContext", "TranspilationContext"
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext", "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$TranspilationContext", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext](v_)
	})
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) SetOuter_(o *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) {
	this_.outer_ = o
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) GetOuter_() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler {
	return this_.outer_
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext() {
	_ = this_
	this_.namedLabels = j2gort.NewOrderedMap[string, *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases]()
	this_.breakCases = j2gort.Deque[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case]{}
	this_.continueCases = j2gort.Deque[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case]{}
	this_.catchCases = j2gort.Deque[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase]{}
	this_.finallyCases = j2gort.Deque[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case]{}
	this_.catchNames = j2gort.NewOrderedSet[string]()
	this_.allCases = j2gort.NewExt_java_util_ArrayList()
	this_.switchBreaks = []pkg_rhino.Node{}

}

// Java constructor: TranspilationContext()
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext(outer_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext {
	this_ := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext{}
	this_.outer_ = outer_
	InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext(this_)
	return this_
}

func InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext(this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) {
	this_.initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext()
	this_.programEndCase = NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case(this_)
	if !(this_.programEndCase.id == int32(0)) {
		j2gort.CheckState(false)
	}
	this_.currentCase = NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case(this_)
	if !(this_.currentCase.id == int32(1)) {
		j2gort.CheckState(false)
	}
	this_.allCases.Add_1(j2gort.PtrToAny(this_.currentCase))
}

// Java method: TranspilationContext.optimizeCaseIds()
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) OptimizeCaseIds() {
	if !(!(this_.allCases.IsEmpty())) {
		j2gort.CheckStateMsg(false, this_.allCases)
	}
	{
		cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.allCases))
		for cur_10.HasNext() {
			rawElem_11 := cur_10.Next()
			currentCase := j2gort.CastPtr[Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case](rawElem_11)
			_ = currentCase
			if currentCase.jumpTo != nil {
				for currentCase.jumpTo.jumpTo != nil {
					{
						r_12 := currentCase
						r_12.jumpTo = currentCase.jumpTo.jumpTo
					}
				}
				if (currentCase.embedInto != 0) && (currentCase.references.Size() == int32(1)) {
					{
						r_13 := currentCase.jumpTo
						r_13.embedInto = currentCase.embedInto
					}
				}
				{
					r_14 := currentCase
					r_14.embedInto = pkg_rhino.Node(0)
				}
				{
					cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(currentCase.references))
					for cur_15.HasNext() {
						reference := pkg_rhino.Node(cur_15.NextNode())
						_ = reference
						reference.SetDouble(float64(currentCase.jumpTo.id))
					}
				}
				currentCase.jumpTo.references.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](currentCase.references))
				currentCase.references.Clear()
			}
		}
	}
	var it j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](this_.allCases.Iterator())
	_ = it
	var prevCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = j2gort.CastPtr[Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case](it.Next())
	_ = prevCase
	if !(prevCase.id == int32(1)) {
		j2gort.CheckState(false)
	}
	for it.HasNext() {
		var currentCase_2 *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = j2gort.CastPtr[Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case](it.Next())
		_ = currentCase_2
		if currentCase_2.references.IsEmpty() {
			if !(currentCase_2.embedInto == 0) {
				j2gort.CheckState(false)
			}
			if prevCase.mayFallThrough {
				prevCase.caseBlock.AddChildrenToBack(currentCase_2.caseBlock.RemoveChildren())
				{
					r_16 := prevCase
					r_16.mayFallThrough = currentCase_2.mayFallThrough
				}
			}
			it.Remove()
			continue
		}
		if currentCase_2.embedInto != 0 {
			if !(currentCase_2.jumpTo == nil) {
				j2gort.CheckState(false)
			}
			if (currentCase_2.references.Size() == int32(1)) && !(currentCase_2.mayFallThrough) {
				currentCase_2.embedInto.ReplaceWith(currentCase_2.caseBlock)
				it.Remove()
				continue
			}
		}
		if prevCase.jumpTo == currentCase_2 {
			if !(prevCase.mayFallThrough) {
				j2gort.CheckState(false)
			}
			if !(!(prevCase.caseBlock.HasChildren())) {
				j2gort.CheckState(false)
			}
			if !(currentCase_2.jumpTo == nil) {
				j2gort.CheckState(false)
			}
			prevCase.caseBlock.AddChildrenToBack(currentCase_2.caseBlock.RemoveChildren())
			{
				r_17 := prevCase
				r_17.mayFallThrough = currentCase_2.mayFallThrough
			}
			{
				cur_18 := j2gort.NewIterCursor(j2gort.CheckNotNull(currentCase_2.references))
				for cur_18.HasNext() {
					reference_2 := pkg_rhino.Node(cur_18.NextNode())
					_ = reference_2
					reference_2.SetDouble(float64(prevCase.id))
				}
			}
			{
				r_19 := prevCase
				r_19.jumpTo = currentCase_2.jumpTo
			}
			prevCase.references.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](currentCase_2.references))
			it.Remove()
			continue
		}
		prevCase = currentCase_2
	}
}

// Java method: TranspilationContext.eliminateSwitchBreaks()
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) EliminateSwitchBreaks() {
	for _, breakNode := range this_.switchBreaks {
		_ = breakNode
		var prevStatement pkg_rhino.Node = breakNode.GetPrevious()
		_ = prevStatement
		if !(prevStatement != 0) {
			j2gort.CheckState(false)
		}
		if !(prevStatement.IsExprResult()) {
			j2gort.CheckState(false)
		}
		prevStatement.ReplaceWith(pkg_rhino.IR_ReturnNode_1(prevStatement.RemoveFirstChild()))
		breakNode.Detach()
	}

	this_.switchBreaks = this_.switchBreaks[:0]
}

// Java method: TranspilationContext.finalizeTransformation(com.google.javascript.rhino.Node generatorBody)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) FinalizeTransformation(generatorBody pkg_rhino.Node) {
	_ = generatorBody
	this_.OptimizeCaseIds()
	if (this_.allCases.Size() == int32(2)) || (this_.allCases.Size() == int32(3)) {
		generatorBody.AddChildToBack(this_.outer_.outer_.astFactory.CreateIf_2(this_.outer_.outer_.astFactory.CreateEq(this_.CallContextMethod(generatorBody, "getNextAddressJsc", AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE(), pkg_colors.StaticGet_StandardColors_NUMBER()), []pkg_rhino.Node{}), this_.outer_.outer_.astFactory.CreateNumber(float64(int32(1)))), j2gort.CastPtr[Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case](this_.allCases.Remove_1_int32(int32(0))).caseBlock).SrcrefTreeIfMissing(generatorBody))
	}
	if this_.allCases.Size() == int32(2) {
		generatorBody.AddChildToBack(this_.outer_.outer_.astFactory.CreateIf_2(this_.outer_.outer_.astFactory.CreateNe(this_.CallContextMethod(generatorBody, "getNextAddressJsc", AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE(), pkg_colors.StaticGet_StandardColors_NUMBER()), []pkg_rhino.Node{}), this_.outer_.outer_.astFactory.CreateNumber(float64(j2gort.CastPtr[Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case](this_.allCases.Get(int32(1))).id))), j2gort.CastPtr[Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case](this_.allCases.Remove_1_int32(int32(0))).caseBlock).SrcrefTreeIfMissing(generatorBody))
	}
	if this_.allCases.Size() == int32(1) {
		generatorBody.AddChildrenToBack(j2gort.CastPtr[Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case](this_.allCases.Remove_1_int32(int32(0))).caseBlock.RemoveChildren())
		this_.EliminateSwitchBreaks()
		return
	}
	var switchNode pkg_rhino.Node = pkg_rhino.IR_SwitchNode(this_.CallContextMethod(generatorBody, "getNextAddressJsc", AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE(), pkg_colors.StaticGet_StandardColors_NUMBER()), []pkg_rhino.Node{}), []pkg_rhino.Node{}).Srcref(generatorBody)
	_ = switchNode
	generatorBody.AddChildToBack(switchNode)
	var switchBody pkg_rhino.Node = switchNode.GetSecondChild().Srcref(generatorBody)
	_ = switchBody
	{
		cur_21 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.allCases))
		for cur_21.HasNext() {
			rawElem_22 := cur_21.Next()
			currentCase := j2gort.CastPtr[Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case](rawElem_22)
			_ = currentCase
			switchBody.AddChildToBack(currentCase.CreateCaseNode())
		}
	}
	this_.allCases.Clear()
}

// Java method: TranspilationContext.checkStateIsEmpty()
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CheckStateIsEmpty() {
	if !(this_.namedLabels.IsEmpty()) {
		j2gort.CheckState(false)
	}
	if !(this_.breakCases.IsEmpty()) {
		j2gort.CheckState(false)
	}
	if !(this_.continueCases.IsEmpty()) {
		j2gort.CheckState(false)
	}
	if !(this_.catchCases.IsEmpty()) {
		j2gort.CheckState(false)
	}
	if !(this_.finallyCases.IsEmpty()) {
		j2gort.CheckState(false)
	}
	if !(this_.nestedFinallyBlockCount == int32(0)) {
		j2gort.CheckState(false)
	}
	if !(this_.allCases.IsEmpty()) {
		j2gort.CheckState(false)
	}
}

// Java method: TranspilationContext.transpileUnmarkedBlock(com.google.javascript.rhino.Node block)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) TranspileUnmarkedBlock(block pkg_rhino.Node) {
	_ = block
	if block.HasChildren() {
		NodeTraversal_Traverse(this_.outer_.outer_.compiler, block, NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler(this_))
		for block.HasChildren() {
			this_.WriteGeneratedNode(block.RemoveFirstChild())
		}
	}
}

// Java method: TranspilationContext.writeGeneratedNode(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) WriteGeneratedNode(n pkg_rhino.Node) {
	_ = n
	this_.currentCase.AddNode(n)
}

// Java method: TranspilationContext.writeGeneratedNodeAndBreak(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) WriteGeneratedNodeAndBreak(n pkg_rhino.Node) {
	_ = n
	this_.WriteGeneratedNode(n)
	this_.WriteGeneratedNode(this_.CreateBreakNodeFor(n))
	{
		r_23 := this_.currentCase
		r_23.mayFallThrough = false
	}
}

// Java method: TranspilationContext.createCase()
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CreateCase() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	return NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case(this_)
}

// Java method: TranspilationContext.maybeCreateCase(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case other)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) MaybeCreateCase(other *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	_ = other
	if other != nil {
		return other
	}
	return this_.CreateCase()
}

// Java method: TranspilationContext.getJsContextNameNode(com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) GetJsContextNameNode(sourceNode pkg_rhino.Node) pkg_rhino.Node {
	_ = sourceNode
	return this_.outer_.outer_.astFactory.CreateName_2_string_AstFactory_Type(this_.GetScopedName("$jscomp$generator$context"), j2gort.AsIface[AstFactory_Type](this_.outer_.outer_.generatorContextType.Get())).Srcref(sourceNode)
}

// Java method: TranspilationContext.getScopedName(java.lang.String name)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) GetScopedName(name string) string {
	_ = name
	return (j2gort.StrValString(name) + "$" + j2gort.StrValString(this_.outer_.uniqueId))
}

// Java method: TranspilationContext.getContextField(com.google.javascript.rhino.Node sourceNode, java.lang.String fieldName)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) GetContextField(sourceNode pkg_rhino.Node, fieldName string) pkg_rhino.Node {
	_ = sourceNode
	_ = fieldName
	return this_.outer_.outer_.astFactory.CreateGetPropWithUnknownType(this_.GetJsContextNameNode(sourceNode), fieldName).Srcref(sourceNode)
}

// Java method: TranspilationContext.callContextMethod(com.google.javascript.rhino.Node sourceNode, java.lang.String methodName, com.google.javascript.jscomp.AstFactory.Type type, com.google.javascript.rhino.Node[] args)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CallContextMethod(sourceNode pkg_rhino.Node, methodName string, type_ AstFactory_Type, args []pkg_rhino.Node) pkg_rhino.Node {
	_ = sourceNode
	_ = methodName
	_ = type_
	_ = args
	var contextField pkg_rhino.Node = this_.GetContextField(sourceNode, methodName)
	_ = contextField
	return this_.outer_.outer_.astFactory.CreateCall(contextField, type_, args).Srcref(sourceNode)
}

// Java method: TranspilationContext.callContextMethodResult(com.google.javascript.rhino.Node sourceNode, java.lang.String methodName, com.google.javascript.jscomp.AstFactory.Type type, com.google.javascript.rhino.Node[] args)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CallContextMethodResult(sourceNode pkg_rhino.Node, methodName string, type_ AstFactory_Type, args []pkg_rhino.Node) pkg_rhino.Node {
	_ = sourceNode
	_ = methodName
	_ = type_
	_ = args
	return this_.outer_.outer_.astFactory.ExprResult(this_.CallContextMethod(sourceNode, methodName, type_, args)).Srcref(sourceNode)
}

// Java method: TranspilationContext.returnContextMethod(com.google.javascript.rhino.Node sourceNode, java.lang.String methodName, com.google.javascript.jscomp.AstFactory.Type type, com.google.javascript.rhino.Node[] args)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) ReturnContextMethod(sourceNode pkg_rhino.Node, methodName string, type_ AstFactory_Type, args []pkg_rhino.Node) pkg_rhino.Node {
	_ = sourceNode
	_ = methodName
	_ = type_
	_ = args
	return this_.outer_.outer_.astFactory.CreateReturn(this_.CallContextMethod(sourceNode, methodName, type_, args)).Srcref(sourceNode)
}

// Java method: TranspilationContext.createBreakNodeFor(com.google.javascript.rhino.Node preBreak)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CreateBreakNodeFor(preBreak pkg_rhino.Node) pkg_rhino.Node {
	_ = preBreak
	var breakNode pkg_rhino.Node = pkg_rhino.IR_BreakNode_0().Srcref(preBreak)
	_ = breakNode
	this_.switchBreaks = append(this_.switchBreaks, breakNode)
	return breakNode
}

// Java method: TranspilationContext.createJumpToNode(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case section, com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CreateJumpToNode(section *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode pkg_rhino.Node) pkg_rhino.Node {
	_ = section
	_ = sourceNode
	return this_.ReturnContextMethod(sourceNode, "jumpTo", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), []pkg_rhino.Node{section.GetNumber(sourceNode)})
}

// Java method: TranspilationContext.writeJumpTo(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case section, com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) WriteJumpTo(section *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode pkg_rhino.Node) {
	_ = section
	_ = sourceNode
	this_.currentCase.JumpTo(section, this_.CreateJumpToBlock(section, false, sourceNode))
}

// Java method: TranspilationContext.createJumpToBlock(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case section, boolean allowEmbedding, com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CreateJumpToBlock(section *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, allowEmbedding bool, sourceNode pkg_rhino.Node) pkg_rhino.Node {
	_ = section
	_ = allowEmbedding
	_ = sourceNode
	if !(section.embedInto == 0) {
		j2gort.CheckState(false)
	}
	var jumpBlock pkg_rhino.Node = this_.outer_.outer_.astFactory.CreateBlock([]pkg_rhino.Node{this_.CallContextMethodResult(sourceNode, "jumpTo", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), []pkg_rhino.Node{section.GetNumber(sourceNode)}), this_.CreateBreakNodeFor(sourceNode)}).Srcref(sourceNode)
	_ = jumpBlock
	if allowEmbedding {
		{
			r_24 := section
			r_24.embedInto = jumpBlock
		}
	}
	return jumpBlock
}

// Java method: TranspilationContext.replaceBreakContinueWithJump(com.google.javascript.rhino.Node sourceNode, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case section, int breakSuppressors)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) ReplaceBreakContinueWithJump(sourceNode pkg_rhino.Node, section *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, breakSuppressors int32) {
	_ = sourceNode
	_ = section
	_ = breakSuppressors
	var jumpMethod string = j2gort.NullString
	_ = jumpMethod
	if this_.nestedFinallyBlockCount > int32(0) {
		jumpMethod = "jumpThroughFinallyBlocks"
	} else if this_.finallyCases.IsEmpty() || (this_.finallyCases.GetFirst().id < section.id) {
		jumpMethod = "jumpTo"
	} else {
		if !(this_.finallyCases.GetFirst().id != section.id) {
			j2gort.CheckState(false)
		}
		jumpMethod = "jumpThroughFinallyBlocks"
	}
	if breakSuppressors == int32(0) {
		this_.CallContextMethodResult(sourceNode, jumpMethod, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), []pkg_rhino.Node{section.GetNumber(sourceNode)}).InsertBefore(sourceNode)
		sourceNode.ReplaceWith(this_.CreateBreakNodeFor(sourceNode))
	} else {
		sourceNode.ReplaceWith(this_.ReturnContextMethod(sourceNode, jumpMethod, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), []pkg_rhino.Node{section.GetNumber(sourceNode)}))
	}
}

// Java method: TranspilationContext.yieldValue(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node expression, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case jumpToSection, com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) YieldValue(expression pkg_rhino.Node, jumpToSection *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode pkg_rhino.Node) {
	_ = expression
	_ = jumpToSection
	_ = sourceNode
	var args j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = args
	args.Add_1(j2gort.BoxNode(func() pkg_rhino.Node {
		if expression == 0 {
			return this_.outer_.outer_.astFactory.CreateUndefinedValue().SrcrefTree(sourceNode)
		}
		return expression
	}()))
	args.Add_1(j2gort.BoxNode(jumpToSection.GetNumber(sourceNode)))
	this_.outer_.context.WriteGeneratedNode(this_.ReturnContextMethod(sourceNode, "yield", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN()), func() []pkg_rhino.Node {
		arr_28 := args.ToArray_1_slice_any(func() []any {
			arr_25 := make([]pkg_rhino.Node, int32(0))
			if arr_25 == nil {
				return nil
			}
			out_26 := make([]any, len(arr_25))
			for i_, el_27 := range arr_25 {
				out_26[i_] = j2gort.BoxNode(el_27)
			}
			return out_26
		}())
		if arr_28 == nil {
			return nil
		}
		out_29 := make([]pkg_rhino.Node, len(arr_28))
		for i_, el_30 := range arr_28 {
			out_29[i_] = j2gort.CastNode[pkg_rhino.Node](el_30)
		}
		return out_29
	}()))
	{
		r_31 := this_.outer_.context.currentCase
		r_31.mayFallThrough = false
	}
}

// Java method: TranspilationContext.yieldAll(com.google.javascript.rhino.Node expression, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case jumpToSection, com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) YieldAll(expression pkg_rhino.Node, jumpToSection *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode pkg_rhino.Node) {
	_ = expression
	_ = jumpToSection
	_ = sourceNode
	this_.WriteGeneratedNode(this_.ReturnContextMethod(sourceNode, "yieldAll", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN()), []pkg_rhino.Node{expression, jumpToSection.GetNumber(sourceNode)}))
	{
		r_32 := this_.outer_.context.currentCase
		r_32.mayFallThrough = false
	}
}

// Java method: TranspilationContext.returnExpression(com.google.javascript.rhino.Node sourceNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node expression)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) ReturnExpression(sourceNode pkg_rhino.Node, expression pkg_rhino.Node) pkg_rhino.Node {
	_ = sourceNode
	_ = expression
	if expression == 0 {
		return this_.CallContextMethod(sourceNode, "return", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), []pkg_rhino.Node{})
	}
	return this_.CallContextMethod(sourceNode, "return", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), []pkg_rhino.Node{expression})
}

// Java method: TranspilationContext.yieldResult(com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) YieldResult(sourceNode pkg_rhino.Node) pkg_rhino.Node {
	_ = sourceNode
	return this_.CallContextMethod(sourceNode, "getYieldResultJsc", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN()), []pkg_rhino.Node{})
}

// Java method: TranspilationContext.addCatchFinallyCases(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case catchCase, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case finallyCase)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) priv_TranspilationContext_addCatchFinallyCases_2711(catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = catchCase
	_ = finallyCase
	if finallyCase != nil {
		if !(this_.catchCases.IsEmpty()) {
			r_33 := this_.catchCases.GetFirst()
			r_33.finallyBlocks = r_33.finallyBlocks + 1
		}
		this_.finallyCases.AddFirst(finallyCase)
	}
	if catchCase != nil {
		this_.catchCases.AddFirst(NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase(catchCase))
	}
}

// Java method: TranspilationContext.getNextCatchCase()
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) priv_TranspilationContext_getNextCatchCase_2712() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	{
		cur_34 := this_.catchCases.Iter()
		for cur_34.HasNext() {

			catchCase := cur_34.Next()
			_ = catchCase
			if catchCase.finallyBlocks == int32(0) {
				return catchCase.catchCase
			}
			break
		}
	}
	return nil
}

// Java method: TranspilationContext.getNextFinallyCase()
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) priv_TranspilationContext_getNextFinallyCase_2713() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	if this_.finallyCases.IsEmpty() {
		return nil
	} else {
		return this_.finallyCases.GetFirst()
	}
}

// Java method: TranspilationContext.removeCatchFinallyCases(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case catchCase, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case finallyCase)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) priv_TranspilationContext_removeCatchFinallyCases_2714(catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = catchCase
	_ = finallyCase
	if catchCase != nil {
		var lastCatch *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase = this_.catchCases.RemoveFirst()
		_ = lastCatch
		if !(lastCatch.finallyBlocks == int32(0)) {
			j2gort.CheckState(false)
		}
		if !(lastCatch.catchCase == catchCase) {
			j2gort.CheckState(false)
		}
	}
	if finallyCase != nil {
		if !(this_.catchCases.IsEmpty()) {
			var finallyBlocks int32 = func() int32 {
				r_36 := this_.catchCases.GetFirst()
				old_37 := r_36.finallyBlocks
				nw_38 := old_37 - 1
				r_36.finallyBlocks = nw_38
				return nw_38
			}()
			_ = finallyBlocks
			if !(finallyBlocks >= int32(0)) {
				j2gort.CheckState(false)
			}
		}
		var lastFinally *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.finallyCases.RemoveFirst()
		_ = lastFinally
		if !(lastFinally == finallyCase) {
			j2gort.CheckState(false)
		}
	}
}

// Java method: TranspilationContext.enterTryBlock(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case catchCase, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case finallyCase, com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) EnterTryBlock(catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode pkg_rhino.Node) {
	_ = catchCase
	_ = finallyCase
	_ = sourceNode
	this_.priv_TranspilationContext_addCatchFinallyCases_2711(catchCase, finallyCase)
	var methodName string = j2gort.NullString
	_ = methodName
	var args j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = args
	if catchCase == nil {
		methodName = "setFinallyBlock"
		args.Add_1(j2gort.BoxNode(finallyCase.GetNumber(sourceNode)))
	} else {
		methodName = "setCatchFinallyBlocks"
		args.Add_1(j2gort.BoxNode(catchCase.GetNumber(sourceNode)))
		if finallyCase != nil {
			args.Add_1(j2gort.BoxNode(finallyCase.GetNumber(sourceNode)))
		}
	}
	this_.WriteGeneratedNode(this_.CallContextMethodResult(sourceNode, methodName, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), func() []pkg_rhino.Node {
		arr_42 := args.ToArray_1_slice_any(func() []any {
			arr_39 := make([]pkg_rhino.Node, int32(0))
			if arr_39 == nil {
				return nil
			}
			out_40 := make([]any, len(arr_39))
			for i_, el_41 := range arr_39 {
				out_40[i_] = j2gort.BoxNode(el_41)
			}
			return out_40
		}())
		if arr_42 == nil {
			return nil
		}
		out_43 := make([]pkg_rhino.Node, len(arr_42))
		for i_, el_44 := range arr_42 {
			out_43[i_] = j2gort.CastNode[pkg_rhino.Node](el_44)
		}
		return out_43
	}()))
}

// Java method: TranspilationContext.leaveTryBlock(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case catchCase, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case endCase, com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) LeaveTryBlock(catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, endCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode pkg_rhino.Node) {
	_ = catchCase
	_ = endCase
	_ = sourceNode
	this_.priv_TranspilationContext_removeCatchFinallyCases_2714(catchCase, nil)
	var args j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = args
	args.Add_1(j2gort.BoxNode(endCase.GetNumber(sourceNode)))
	var nextCatchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.priv_TranspilationContext_getNextCatchCase_2712()
	_ = nextCatchCase
	if nextCatchCase != nil {
		args.Add_1(j2gort.BoxNode(nextCatchCase.GetNumber(sourceNode)))
	}
	this_.WriteGeneratedNodeAndBreak(this_.CallContextMethodResult(sourceNode, "leaveTryBlock", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), func() []pkg_rhino.Node {
		arr_48 := args.ToArray_1_slice_any(func() []any {
			arr_45 := make([]pkg_rhino.Node, int32(0))
			if arr_45 == nil {
				return nil
			}
			out_46 := make([]any, len(arr_45))
			for i_, el_47 := range arr_45 {
				out_46[i_] = j2gort.BoxNode(el_47)
			}
			return out_46
		}())
		if arr_48 == nil {
			return nil
		}
		out_49 := make([]pkg_rhino.Node, len(arr_48))
		for i_, el_50 := range arr_48 {
			out_49[i_] = j2gort.CastNode[pkg_rhino.Node](el_50)
		}
		return out_49
	}()))
}

// Java method: TranspilationContext.enterCatchBlock(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case finallyCase, com.google.javascript.rhino.Node exceptionName)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) EnterCatchBlock(finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, exceptionName pkg_rhino.Node) {
	_ = finallyCase
	_ = exceptionName
	if !(exceptionName.IsName()) {
		j2gort.CheckState(false)
	}
	this_.priv_TranspilationContext_addCatchFinallyCases_2711(nil, finallyCase)
	var nextCatchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.priv_TranspilationContext_getNextCatchCase_2712()
	_ = nextCatchCase
	if this_.catchNames.Add(exceptionName.GetString()) {
		this_.outer_.priv_SingleGeneratorFunctionTranspiler_hoistVarNode_2661(pkg_rhino.IR_Var_1(exceptionName.CloneNode_0()).Srcref(exceptionName))
	}
	var args j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = args
	if nextCatchCase != nil {
		args.Add_1(j2gort.BoxNode(nextCatchCase.GetNumber(exceptionName)))
	}
	var enterCatchBlockCall pkg_rhino.Node = this_.CallContextMethod(exceptionName, "enterCatchBlock", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN()), func() []pkg_rhino.Node {
		arr_54 := args.ToArray_1_slice_any(func() []any {
			arr_51 := make([]pkg_rhino.Node, int32(0))
			if arr_51 == nil {
				return nil
			}
			out_52 := make([]any, len(arr_51))
			for i_, el_53 := range arr_51 {
				out_52[i_] = j2gort.BoxNode(el_53)
			}
			return out_52
		}())
		if arr_54 == nil {
			return nil
		}
		out_55 := make([]pkg_rhino.Node, len(arr_54))
		for i_, el_56 := range arr_54 {
			out_55[i_] = j2gort.CastNode[pkg_rhino.Node](el_56)
		}
		return out_55
	}())
	_ = enterCatchBlockCall
	exceptionName.SetColor(enterCatchBlockCall.GetColor())
	this_.WriteGeneratedNode(this_.outer_.outer_.astFactory.ExprResult(this_.outer_.outer_.astFactory.CreateAssign_2_Node_Node(exceptionName, enterCatchBlockCall).Srcref(exceptionName)).Srcref(exceptionName))
}

// Java method: TranspilationContext.leaveCatchBlock(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case finallyCase, com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) LeaveCatchBlock(finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode pkg_rhino.Node) {
	_ = finallyCase
	_ = sourceNode
	if finallyCase != nil {
		this_.priv_TranspilationContext_removeCatchFinallyCases_2714(nil, finallyCase)
		this_.WriteJumpTo(finallyCase, sourceNode)
	}
}

// Java method: TranspilationContext.enterFinallyBlock(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case catchCase, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case finallyCase, com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) EnterFinallyBlock(catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode pkg_rhino.Node) {
	_ = catchCase
	_ = finallyCase
	_ = sourceNode
	this_.priv_TranspilationContext_removeCatchFinallyCases_2714(catchCase, finallyCase)
	var nextCatchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.priv_TranspilationContext_getNextCatchCase_2712()
	_ = nextCatchCase
	var nextFinallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = this_.priv_TranspilationContext_getNextFinallyCase_2713()
	_ = nextFinallyCase
	var args j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = args
	if this_.nestedFinallyBlockCount == int32(0) {
		if (nextCatchCase != nil) || (nextFinallyCase != nil) {
			args.Add_1(j2gort.BoxNode(func() pkg_rhino.Node {
				if nextCatchCase == nil {
					return this_.outer_.outer_.astFactory.CreateNumber(float64(int32(0))).Srcref(sourceNode)
				}
				return nextCatchCase.GetNumber(sourceNode)
			}()))
			if nextFinallyCase != nil {
				args.Add_1(j2gort.BoxNode(nextFinallyCase.GetNumber(sourceNode)))
			}
		}
	} else {
		args.Add_1(j2gort.BoxNode(func() pkg_rhino.Node {
			if nextCatchCase == nil {
				return this_.outer_.outer_.astFactory.CreateNumber(float64(int32(0))).Srcref(sourceNode)
			}
			return nextCatchCase.GetNumber(sourceNode)
		}()))
		args.Add_1(j2gort.BoxNode(func() pkg_rhino.Node {
			if nextFinallyCase == nil {
				return this_.outer_.outer_.astFactory.CreateNumber(float64(int32(0))).Srcref(sourceNode)
			}
			return nextFinallyCase.GetNumber(sourceNode)
		}()))
		args.Add_1(j2gort.BoxNode(this_.outer_.outer_.astFactory.CreateNumber(float64(this_.nestedFinallyBlockCount)).Srcref(sourceNode)))
	}
	this_.WriteGeneratedNode(this_.CallContextMethodResult(sourceNode, "enterFinallyBlock", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), func() []pkg_rhino.Node {
		arr_60 := args.ToArray_1_slice_any(func() []any {
			arr_57 := make([]pkg_rhino.Node, int32(0))
			if arr_57 == nil {
				return nil
			}
			out_58 := make([]any, len(arr_57))
			for i_, el_59 := range arr_57 {
				out_58[i_] = j2gort.BoxNode(el_59)
			}
			return out_58
		}())
		if arr_60 == nil {
			return nil
		}
		out_61 := make([]pkg_rhino.Node, len(arr_60))
		for i_, el_62 := range arr_60 {
			out_61[i_] = j2gort.CastNode[pkg_rhino.Node](el_62)
		}
		return out_61
	}()))
	this_.nestedFinallyBlockCount = this_.nestedFinallyBlockCount + 1
}

// Java method: TranspilationContext.leaveFinallyBlock(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case endCase, com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) LeaveFinallyBlock(endCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode pkg_rhino.Node) {
	_ = endCase
	_ = sourceNode
	var args j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = args
	args.Add_1(j2gort.BoxNode(endCase.GetNumber(sourceNode)))
	{
		var iifeCond_1 int32
		{

			old_63 := this_.nestedFinallyBlockCount
			nw_64 := old_63 - 1
			this_.nestedFinallyBlockCount = nw_64
			iifeCond_1 = nw_64
		}
		if iifeCond_1 != int32(0) {
			args.Add_1(j2gort.BoxNode(this_.outer_.outer_.astFactory.CreateNumber(float64(this_.nestedFinallyBlockCount)).Srcref(sourceNode)))
		}
	}
	this_.WriteGeneratedNodeAndBreak(this_.CallContextMethodResult(sourceNode, "leaveFinallyBlock", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), func() []pkg_rhino.Node {
		arr_68 := args.ToArray_1_slice_any(func() []any {
			arr_65 := make([]pkg_rhino.Node, int32(0))
			if arr_65 == nil {
				return nil
			}
			out_66 := make([]any, len(arr_65))
			for i_, el_67 := range arr_65 {
				out_66[i_] = j2gort.BoxNode(el_67)
			}
			return out_66
		}())
		if arr_68 == nil {
			return nil
		}
		out_69 := make([]pkg_rhino.Node, len(arr_68))
		for i_, el_70 := range arr_68 {
			out_69[i_] = j2gort.CastNode[pkg_rhino.Node](el_70)
		}
		return out_69
	}()))
}

// Java method: TranspilationContext.switchCaseTo(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case caseSection)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) SwitchCaseTo(caseSection *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = caseSection
	this_.currentCase.WillFollowBy(caseSection)
	this_.allCases.Add_1(j2gort.PtrToAny(caseSection))
	this_.currentCase = caseSection
}

// Java method: TranspilationContext.pushLabels(java.util.ArrayList<com.google.javascript.rhino.Node> labelNames, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case breakCase, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case continueCase)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) PushLabels(labelNames j2gort.Ext_java_util_ArrayList, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = labelNames
	_ = breakCase
	_ = continueCase
	{
		cur_71 := j2gort.NewIterCursor(j2gort.CheckNotNull(labelNames))
		for cur_71.HasNext() {
			labelName := pkg_rhino.Node(cur_71.NextNode())
			_ = labelName
			if !(labelName.IsLabelName()) {
				j2gort.CheckState(false)
			}
			this_.namedLabels.Put(labelName.GetString(), NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases(breakCase, continueCase))
		}
	}
}

// Java method: TranspilationContext.popLabels(java.util.ArrayList<com.google.javascript.rhino.Node> labelNames)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) PopLabels(labelNames j2gort.Ext_java_util_ArrayList) {
	_ = labelNames
	{
		cur_72 := j2gort.NewIterCursor(j2gort.CheckNotNull(labelNames))
		for cur_72.HasNext() {
			labelName := pkg_rhino.Node(cur_72.NextNode())
			_ = labelName
			if !(labelName.IsLabelName()) {
				j2gort.CheckState(false)
			}
			this_.namedLabels.Remove(labelName.GetString())
		}
	}
}

// Java method: TranspilationContext.pushBreakContext(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case breakCase)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) PushBreakContext(breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = breakCase
	this_.breakCases.Push(breakCase)
}

// Java method: TranspilationContext.pushBreakContinueContext(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case breakCase, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case continueCase)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) PushBreakContinueContext(breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = breakCase
	_ = continueCase
	this_.PushBreakContext(breakCase)
	this_.continueCases.Push(continueCase)
}

// Java method: TranspilationContext.popBreakContext()
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) PopBreakContext() {
	this_.breakCases.Pop()
}

// Java method: TranspilationContext.popBreakContinueContext()
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) PopBreakContinueContext() {
	this_.PopBreakContext()
	this_.continueCases.Pop()
}

// Java class: com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case struct {
	outer_         *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext
	references     j2gort.Ext_java_util_ArrayList                                                    `gson:"references"`
	jumpTo         *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case `gson:"jumpTo"`
	id             int32                                                                             `gson:"id"`
	caseBlock      pkg_rhino.Node                                                                    `gson:"caseBlock"`
	embedInto      pkg_rhino.Node                                                                    `gson:"embedInto"`
	idHash_        atomic.Int32
	mayFallThrough bool `gson:"mayFallThrough"`
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) IsJavaNil_() bool {
	return this_ == nil
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$TranspilationContext$Case", "Case"
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case", "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$TranspilationContext$Case", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case](v_)
	})
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) SetOuter_(o *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) {
	this_.outer_ = o
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) GetOuter_() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext {
	return this_.outer_
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case() {
	_ = this_
	this_.references = j2gort.NewExt_java_util_ArrayList()
	this_.mayFallThrough = true
}

// Java constructor: Case()
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case(outer_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	this_ := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case{}
	this_.outer_ = outer_
	InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case(this_)
	return this_
}

func InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case(this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	this_.initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case()
	this_.id = func() int32 {
		old_73 := this_.outer_.caseIdCounter
		nw_74 := old_73 + 1
		this_.outer_.caseIdCounter = nw_74
		return old_73
	}()
	this_.caseBlock = this_.outer_.outer_.outer_.astFactory.CreateBlock([]pkg_rhino.Node{}).Srcref(this_.outer_.outer_.originalGeneratorBody)
}

// Java method: Case.createCaseNode()
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) CreateCaseNode() pkg_rhino.Node {
	return pkg_rhino.IR_CaseNode(this_.outer_.outer_.outer_.astFactory.CreateNumber(float64(this_.id)).Srcref(this_.caseBlock), this_.caseBlock).Srcref(this_.caseBlock)
}

// Java method: Case.getNumber(com.google.javascript.rhino.Node sourceNode)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) GetNumber(sourceNode pkg_rhino.Node) pkg_rhino.Node {
	_ = sourceNode
	if this_.jumpTo != nil {
		return this_.jumpTo.GetNumber(sourceNode)
	}
	var node pkg_rhino.Node = this_.outer_.outer_.outer_.astFactory.CreateNumber(float64(this_.id)).Srcref(sourceNode)
	_ = node
	this_.references.Add_1(j2gort.BoxNode(node))
	return node
}

// Java method: Case.jumpTo(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case other, com.google.javascript.rhino.Node jumpBlock)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) JumpTo(other *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, jumpBlock pkg_rhino.Node) {
	_ = other
	_ = jumpBlock
	if !(jumpBlock.IsBlock()) {
		j2gort.CheckState(false)
	}
	if !(this_.jumpTo == nil) {
		j2gort.CheckState(false)
	}
	this_.WillFollowBy(other)
	this_.caseBlock.AddChildrenToBack(jumpBlock.RemoveChildren())
	this_.mayFallThrough = false
}

// Java method: Case.willFollowBy(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case other)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) WillFollowBy(other *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = other
	if (this_.jumpTo == nil) && !(this_.caseBlock.HasChildren()) {
		if !(other.jumpTo == nil) {
			j2gort.CheckState(false)
		}
		this_.jumpTo = other
	}
}

// Java method: Case.addNode(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) AddNode(n pkg_rhino.Node) {
	_ = n
	if !(this_.jumpTo == nil) {
		j2gort.CheckState(false)
	}
	if !(pkg_rhino.IR_MayBeStatement(n)) {
		j2gort.CheckState(false)
	}
	this_.caseBlock.AddChildToBack(n)
}

// Java class: com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.UnmarkedNodeTranspiler
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler struct {
	outer_              *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext
	breakSuppressors    int32 `gson:"breakSuppressors"`
	continueSuppressors int32 `gson:"continueSuppressors"`
	idHash_             atomic.Int32
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) IsJavaNil_() bool {
	return this_ == nil
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$TranspilationContext$UnmarkedNodeTranspiler", "UnmarkedNodeTranspiler"
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.UnmarkedNodeTranspiler", "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$TranspilationContext$UnmarkedNodeTranspiler", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler](v_)
	})
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) SetOuter_(o *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) {
	this_.outer_ = o
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) GetOuter_() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext {
	return this_.outer_
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler() {
	_ = this_
}

// Java constructor: UnmarkedNodeTranspiler()
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler(outer_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler {
	this_ := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler{}
	this_.outer_ = outer_
	InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler(this_)
	return this_
}

func InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler(this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) {
	this_.initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler()
}

// Java method: UnmarkedNodeTranspiler.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if n.IsGeneratorSafe() {
		n.SetGeneratorSafe(false)
		return false
	}
	if !(!(n.IsGeneratorMarker())) {
		j2gort.CheckState(false)
	}
	if !(!(n.IsSuper())) {
		j2gort.CheckStateStr(false, "Reference to SUPER is not supported")
	}
	if NodeUtil_IsLoopStructure(n) {
		this_.continueSuppressors = this_.continueSuppressors + 1
		this_.breakSuppressors = this_.breakSuppressors + 1
	} else if n.IsSwitch() {
		this_.breakSuppressors = this_.breakSuppressors + 1
	}
	if n.IsBreak() || n.IsContinue() {
		if n.HasChildren() {
			this_.VisitNamedBreakContinue(n)
		} else {
			this_.VisitBreakContinue(n)
		}
		return false
	}
	return !(n.IsFunction())
}

// Java method: UnmarkedNodeTranspiler.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if NodeUtil_IsLoopStructure(n) {
		this_.continueSuppressors = this_.continueSuppressors - 1
		this_.breakSuppressors = this_.breakSuppressors - 1
	} else if n.IsSwitch() {
		this_.breakSuppressors = this_.breakSuppressors - 1
	} else if n.IsThis() {
		this_.VisitThis(n)
	} else if n.IsReturn() {
		this_.VisitReturn(n)
	} else if n.IsName() && (j2gort.CheckNotNullString(n.GetString()) == "arguments") {
		this_.VisitArguments(n)
	} else if n.IsVar() {
		if parent.IsVanillaFor() {
			this_.priv_UnmarkedNodeTranspiler_visitVanillaForLoopVar_2739(n)
		} else if parent.IsForIn() {
			this_.priv_UnmarkedNodeTranspiler_visitForInLoopVar_2740(n)
		} else {
			this_.VisitVarStatement(n)
		}
	}
}

// Java method: UnmarkedNodeTranspiler.visitReturn(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitReturn(n pkg_rhino.Node) {
	_ = n
	n.AddChildToFront(this_.outer_.ReturnExpression(n, n.RemoveFirstChild()))
}

// Java method: UnmarkedNodeTranspiler.visitNamedBreakContinue(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitNamedBreakContinue(n pkg_rhino.Node) {
	_ = n
	if !(n.GetFirstChild().IsLabelName()) {
		j2gort.CheckState(false)
	}
	var cases *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases = this_.outer_.namedLabels.Get(n.GetFirstChild().GetString())
	_ = cases
	if cases != nil {
		var caseSection *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
		_ = caseSection
		if n.IsBreak() {
			caseSection = cases.breakCase
		} else {
			caseSection = cases.continueCase
		}
		this_.outer_.outer_.context.ReplaceBreakContinueWithJump(n, caseSection, this_.breakSuppressors)
	}
}

// Java method: UnmarkedNodeTranspiler.visitBreakContinue(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitBreakContinue(n pkg_rhino.Node) {
	_ = n
	var caseSection *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = nil
	_ = caseSection
	if n.IsBreak() && (this_.breakSuppressors == int32(0)) {
		caseSection = this_.outer_.breakCases.GetFirst()
	}
	if n.IsContinue() && (this_.continueSuppressors == int32(0)) {
		caseSection = this_.outer_.continueCases.GetFirst()
	}
	if caseSection != nil {
		this_.outer_.outer_.context.ReplaceBreakContinueWithJump(n, caseSection, this_.breakSuppressors)
	}
}

// Java method: UnmarkedNodeTranspiler.visitThis(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitThis(n pkg_rhino.Node) {
	_ = n
	var newThis pkg_rhino.Node = this_.outer_.outer_.outer_.astFactory.CreateName_2_string_AstFactory_Type(this_.outer_.outer_.context.GetScopedName("$jscomp$generator$this"), AstFactory_S_Type_1_Node(n))
	_ = newThis
	n.ReplaceWith(newThis)
	if !(this_.outer_.thisReferenceFound) {
		var var_ pkg_rhino.Node = pkg_rhino.IR_Var_2(newThis.CloneNode_0().Srcref(n), n).Srcref(this_.outer_.outer_.newGeneratorHoistBlock)
		_ = var_
		this_.outer_.outer_.priv_SingleGeneratorFunctionTranspiler_hoistVarNode_2661(var_)
		this_.outer_.thisReferenceFound = true
	}
}

// Java method: UnmarkedNodeTranspiler.visitArguments(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitArguments(n pkg_rhino.Node) {
	_ = n
	var newArguments pkg_rhino.Node = this_.outer_.outer_.outer_.astFactory.CreateName_2_string_AstFactory_Type(this_.outer_.outer_.context.GetScopedName("$jscomp$generator$arguments"), AstFactory_S_Type_1_Node(n)).Srcref(n)
	_ = newArguments
	n.ReplaceWith(newArguments)
	if !(this_.outer_.argumentsReferenceFound) {
		var var_ pkg_rhino.Node = pkg_rhino.IR_Var_2(newArguments.CloneNode_0(), n).Srcref(this_.outer_.outer_.newGeneratorHoistBlock)
		_ = var_
		this_.outer_.outer_.priv_SingleGeneratorFunctionTranspiler_hoistVarNode_2661(var_)
		this_.outer_.argumentsReferenceFound = true
	}
}

// Java method: UnmarkedNodeTranspiler.visitVarStatement(com.google.javascript.rhino.Node varStatement)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitVarStatement(varStatement pkg_rhino.Node) {
	_ = varStatement
	var commaExpression pkg_rhino.Node = this_.priv_UnmarkedNodeTranspiler_extractAssignmentsToCommaExpression_2741(varStatement)
	_ = commaExpression
	if commaExpression == 0 {
		varStatement.Detach()
	} else {
		varStatement.ReplaceWith(this_.outer_.outer_.outer_.astFactory.ExprResult(commaExpression))
	}
	this_.outer_.outer_.priv_SingleGeneratorFunctionTranspiler_hoistVarNode_2661(varStatement)
}

// Java method: UnmarkedNodeTranspiler.visitVanillaForLoopVar(com.google.javascript.rhino.Node varDeclaration)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) priv_UnmarkedNodeTranspiler_visitVanillaForLoopVar_2739(varDeclaration pkg_rhino.Node) {
	_ = varDeclaration
	var commaExpression pkg_rhino.Node = this_.priv_UnmarkedNodeTranspiler_extractAssignmentsToCommaExpression_2741(varDeclaration)
	_ = commaExpression
	if commaExpression == 0 {
		varDeclaration.ReplaceWith(pkg_rhino.IR_Empty())
	} else {
		varDeclaration.ReplaceWith(commaExpression)
	}
	this_.outer_.outer_.priv_SingleGeneratorFunctionTranspiler_hoistVarNode_2661(varDeclaration)
}

// Java method: UnmarkedNodeTranspiler.visitForInLoopVar(com.google.javascript.rhino.Node varDeclaration)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) priv_UnmarkedNodeTranspiler_visitForInLoopVar_2740(varDeclaration pkg_rhino.Node) {
	_ = varDeclaration
	var varName pkg_rhino.Node = varDeclaration.GetOnlyChild()
	_ = varName
	if !(!(varName.HasChildren())) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(varName))
	}
	var clonedVarName pkg_rhino.Node = varName.CloneNode_0().SetJSDocInfo(nil)
	_ = clonedVarName
	varDeclaration.ReplaceWith(clonedVarName)
	this_.outer_.outer_.priv_SingleGeneratorFunctionTranspiler_hoistVarNode_2661(varDeclaration)
}

// Java method: UnmarkedNodeTranspiler.extractAssignmentsToCommaExpression(com.google.javascript.rhino.Node varDeclaration)
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) priv_UnmarkedNodeTranspiler_extractAssignmentsToCommaExpression_2741(varDeclaration pkg_rhino.Node) pkg_rhino.Node {
	_ = varDeclaration
	var assignments []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = assignments
	{
		var varName pkg_rhino.Node = varDeclaration.GetFirstChild()
		_ = varName
		for ; varName != 0; varName = varName.GetNext() {
			if varName.HasChildren() {
				var copiedVarName pkg_rhino.Node = varName.CloneNode_0().SetJSDocInfo(nil)
				_ = copiedVarName
				var assign pkg_rhino.Node = this_.outer_.outer_.outer_.astFactory.CreateAssign_2_Node_Node(copiedVarName, varName.RemoveFirstChild()).Srcref(varName)
				_ = assign
				assignments = append(assignments, assign)
			}
		}
	}
	var commaExpression pkg_rhino.Node = pkg_rhino.Node(0)
	_ = commaExpression
	for _, assignment := range assignments {
		_ = assignment
		if commaExpression == 0 {
			commaExpression = assignment
		} else {
			commaExpression = this_.outer_.outer_.outer_.astFactory.CreateComma(commaExpression, assignment).Srcref(assignment)
		}
	}

	return commaExpression
}

// Java class: com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.CatchCase
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase struct {
	catchCase     *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case `gson:"catchCase"`
	finallyBlocks int32                                                                             `gson:"finallyBlocks"`
	idHash_       atomic.Int32
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase) IsJavaNil_() bool {
	return this_ == nil
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$TranspilationContext$CatchCase", "CatchCase"
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.CatchCase", "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$TranspilationContext$CatchCase", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase](v_)
	})
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase) initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase() {
	_ = this_
}

// Java constructor: CatchCase(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case catchCase)
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase(catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase {
	this_ := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase{}
	InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase(this_, catchCase)
	return this_
}

func InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase(this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase, catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = catchCase
	this_.initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase()
	this_.catchCase = catchCase
}

// Java class: com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.LabelCases
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases struct {
	breakCase    *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case `gson:"breakCase"`
	continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case `gson:"continueCase"`
	idHash_      atomic.Int32
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases) IsJavaNil_() bool {
	return this_ == nil
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$TranspilationContext$LabelCases", "LabelCases"
}
func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.LabelCases", "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$TranspilationContext$LabelCases", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases](v_)
	})
}

func (this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases) initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases() {
	_ = this_
}

// Java constructor: LabelCases(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case breakCase, com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.@org.jspecify.annotations.Nullable Case continueCase)
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases(breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases {
	this_ := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases{}
	InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases(this_, breakCase, continueCase)
	return this_
}

func InitEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases(this_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	_ = breakCase
	_ = continueCase
	this_.initFields_Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases()
	this_.breakCase = breakCase
	this_.continueCase = continueCase
}

// Java class: com.google.javascript.jscomp.Es6RewriteGenerators.YieldNodeMarker
type Es6RewriteGenerators_YieldNodeMarker struct {
	idHash_ atomic.Int32
}

func (this_ *Es6RewriteGenerators_YieldNodeMarker) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteGenerators_YieldNodeMarker) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteGenerators_YieldNodeMarker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteGenerators$YieldNodeMarker", "YieldNodeMarker"
}
func (this_ *Es6RewriteGenerators_YieldNodeMarker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteGenerators_YieldNodeMarker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteGenerators.YieldNodeMarker", "com.google.javascript.jscomp.Es6RewriteGenerators$YieldNodeMarker", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteGenerators_YieldNodeMarker](v_)
	})
}
func (*Es6RewriteGenerators_YieldNodeMarker) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *Es6RewriteGenerators_YieldNodeMarker) initFields_Es6RewriteGenerators_YieldNodeMarker() {
	_ = this_
}

// Java constructor: YieldNodeMarker()
func NewEs6RewriteGenerators_YieldNodeMarker() *Es6RewriteGenerators_YieldNodeMarker {
	this_ := &Es6RewriteGenerators_YieldNodeMarker{}
	InitEs6RewriteGenerators_YieldNodeMarker(this_)
	return this_
}

func InitEs6RewriteGenerators_YieldNodeMarker(this_ *Es6RewriteGenerators_YieldNodeMarker) {
	this_.initFields_Es6RewriteGenerators_YieldNodeMarker()
}

// Java method: YieldNodeMarker.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteGenerators_YieldNodeMarker) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	return !(n.IsFunction())
}

// Java method: YieldNodeMarker.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteGenerators_YieldNodeMarker) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsYield() {
		n.SetGeneratorMarker(true)
	}
	if (parent != 0) && n.IsGeneratorMarker() {
		parent.SetGeneratorMarker(true)
	}
}

// Java class: SwitchCase
type Es6RewriteGenerators__Local_42_SwitchCase struct {
	outer_        *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler
	generatedCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case `gson:"generatedCase"`
	body          pkg_rhino.Node                                                                    `gson:"body"`
	idHash_       atomic.Int32
}

func (this_ *Es6RewriteGenerators__Local_42_SwitchCase) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteGenerators__Local_42_SwitchCase) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteGenerators__Local_42_SwitchCase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteGenerators$SingleGeneratorFunctionTranspiler$1SwitchCase", "SwitchCase"
}
func (this_ *Es6RewriteGenerators__Local_42_SwitchCase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteGenerators__Local_42_SwitchCase) HasIdentityEquals_() {}
func (this_ *Es6RewriteGenerators__Local_42_SwitchCase) SetOuter_(o *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) {
	this_.outer_ = o
}
func (this_ *Es6RewriteGenerators__Local_42_SwitchCase) GetOuter_() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler {
	return this_.outer_
}

func (this_ *Es6RewriteGenerators__Local_42_SwitchCase) initFields_Es6RewriteGenerators__Local_42_SwitchCase() {
	_ = this_
}

// Java constructor: SwitchCase(com.google.javascript.jscomp.Es6RewriteGenerators.SingleGeneratorFunctionTranspiler.TranspilationContext.Case generatedCase, com.google.javascript.rhino.Node caseNode)
func NewEs6RewriteGenerators__Local_42_SwitchCase(outer_ *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler, generatedCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, caseNode pkg_rhino.Node) *Es6RewriteGenerators__Local_42_SwitchCase {
	this_ := &Es6RewriteGenerators__Local_42_SwitchCase{}
	this_.outer_ = outer_
	InitEs6RewriteGenerators__Local_42_SwitchCase(this_, generatedCase, caseNode)
	return this_
}

func InitEs6RewriteGenerators__Local_42_SwitchCase(this_ *Es6RewriteGenerators__Local_42_SwitchCase, generatedCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, caseNode pkg_rhino.Node) {
	_ = generatedCase
	_ = caseNode
	this_.initFields_Es6RewriteGenerators__Local_42_SwitchCase()
	this_.generatedCase = generatedCase
	this_.body = caseNode
}
func (*Es6RewriteGenerators_YieldExposer) NoOpVisit_()                                  {}
func (*Es6RewriteGenerators_YieldExposer) CanReleaseSyntacticScopes_()                  {}
func (*Es6RewriteGenerators_GeneratorFunctionsTranspiler) AlwaysShouldTraverse_()       {}
func (*Es6RewriteGenerators_GeneratorFunctionsTranspiler) CanReleaseSyntacticScopes_()  {}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) NoOpVisit_() {}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) CanReleaseSyntacticScopes_() {
}
func (*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) CanReleaseSyntacticScopes_() {
}
func (*Es6RewriteGenerators_YieldNodeMarker) CanReleaseSyntacticScopes_() {}
