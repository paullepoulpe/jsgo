/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/InstrumentAsyncContext.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.InstrumentAsyncContext
type InstrumentAsyncContext struct {
	start                 RuntimeJsLibManager_JsLibField                                             `gson:"start"`
	compiler              AbstractCompiler                                                           `gson:"compiler"`
	astFactory            *AstFactory                                                                `gson:"astFactory"`
	uniqueIdSupplier      *UniqueIdSupplier                                                          `gson:"uniqueIdSupplier"`
	tryFunctionStack      j2gort.Deque[pkg_rhino.Node]                                               `gson:"tryFunctionStack"`
	needsInstrumentation  j2gort.OrderedMap[pkg_rhino.Node, *InstrumentAsyncContext_FunctionContext] `gson:"needsInstrumentation"`
	alreadyInstrumented   j2gort.OrderedSet[pkg_rhino.Node]                                          `gson:"alreadyInstrumented"`
	hasSuper              j2gort.OrderedSet[pkg_rhino.Node]                                          `gson:"hasSuper"`
	idHash_               atomic.Int32
	shouldInstrumentAwait bool `gson:"shouldInstrumentAwait"`
}

func (this_ *InstrumentAsyncContext) BindSelf_(s any) {
	_ = s
}
func (this_ *InstrumentAsyncContext) IsJavaNil_() bool { return this_ == nil }
func (*InstrumentAsyncContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InstrumentAsyncContext", "InstrumentAsyncContext"
}
func (this_ *InstrumentAsyncContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InstrumentAsyncContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InstrumentAsyncContext", "com.google.javascript.jscomp.InstrumentAsyncContext", func(v_ any) bool {
		return j2gort.InstanceOf[*InstrumentAsyncContext](v_)
	})
}
func (*InstrumentAsyncContext) Implements_jscomp_CompilerPass()           {}
func (*InstrumentAsyncContext) Implements_jscomp_NodeTraversal_Callback() {}

func StaticGet_InstrumentAsyncContext_FACTORY() string { return "\u1d43\u1d9cfactory" }
func StaticGet_InstrumentAsyncContext_SUSPEND() string { return "\u1d43\u1d9csuspend" }
func StaticGet_InstrumentAsyncContext_RESUME() string  { return "\u1d43\u1d9cresume" }
func (this_ *InstrumentAsyncContext) initFields_InstrumentAsyncContext() {
	_ = this_
	this_.tryFunctionStack = j2gort.Deque[pkg_rhino.Node]{}
	this_.needsInstrumentation = j2gort.NewOrderedMap[pkg_rhino.Node, *InstrumentAsyncContext_FunctionContext]()
	this_.alreadyInstrumented = j2gort.NewOrderedSet[pkg_rhino.Node]()
	this_.hasSuper = j2gort.NewOrderedSet[pkg_rhino.Node]()
}

// Java constructor: InstrumentAsyncContext(com.google.javascript.jscomp.AbstractCompiler compiler, boolean shouldInstrumentAwait)
func NewInstrumentAsyncContext(compiler AbstractCompiler, shouldInstrumentAwait bool) *InstrumentAsyncContext {
	this_ := &InstrumentAsyncContext{}
	InitInstrumentAsyncContext(this_, compiler, shouldInstrumentAwait)
	return this_
}

func InitInstrumentAsyncContext(this_ *InstrumentAsyncContext, compiler AbstractCompiler, shouldInstrumentAwait bool) {
	_ = compiler
	_ = shouldInstrumentAwait
	this_.initFields_InstrumentAsyncContext()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
	this_.shouldInstrumentAwait = shouldInstrumentAwait
	this_.start = compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.asyncContextStart")
}

// Java method: InstrumentAsyncContext.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InstrumentAsyncContext) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(this_.compiler.GetLifeCycleStage()))
	}
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: InstrumentAsyncContext.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InstrumentAsyncContext) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsTry() || n.IsFunction() {
		this_.tryFunctionStack.Push(j2gort.CheckNotNullNode(n))
	}
	return true
}

// Java method: InstrumentAsyncContext.isReentrance(com.google.javascript.rhino.Node n)
func (this_ *InstrumentAsyncContext) priv_InstrumentAsyncContext_isReentrance_3407(n pkg_rhino.Node) bool {
	_ = n
	if n.IsYield() {
		return true
	} else if !(this_.shouldInstrumentAwait) {
		return false
	}
	return (n.IsForAwaitOf() || n.IsAwait())
}

// Java method: InstrumentAsyncContext.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent)
func (this_ *InstrumentAsyncContext) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() && javalang.StringStartsWith(j2gort.CheckNotNullString(n.GetString()), "\u1d43\u1d9cfactory") {
		this_.alreadyInstrumented.Add(t.GetEnclosingFunction())
	} else if n.IsSuper() {
		this_.hasSuper.Add(t.GetEnclosingFunction())
	} else if this_.priv_InstrumentAsyncContext_isReentrance_3407(n) {
		this_.priv_InstrumentAsyncContext_instrumentReentrance_3408(t, n, parent)
	} else {
		var iifeCond_1 bool
		{

			rn_1 := n
			an_2 := this_.tryFunctionStack.PeekFirst()
			iifeCond_1 = j2gort.CheckNotNullNode(rn_1) == an_2
		}
		if iifeCond_1 {
			this_.tryFunctionStack.Pop()
			if this_.alreadyInstrumented.Contains(n) {
				return
			}
			var context *InstrumentAsyncContext_FunctionContext
			_ = context
			if n.IsGeneratorFunction() {
				context = j2gort.OrderedMapComputeIfAbsent(&this_.needsInstrumentation, n, nil, func(lp_7 pkg_rhino.Node) *InstrumentAsyncContext_FunctionContext {
					unused := lp_7
					_ = unused
					return this_.priv_InstrumentAsyncContext_createNewFunctionContext_3420(t.GetInput())
				})

			} else {
				context = this_.needsInstrumentation.Get(n)
			}
			if context == nil {
				return
			}
			if n.IsTry() {
				this_.InstrumentTry(n, context)
			} else if n.IsGeneratorFunction() {
				this_.InstrumentGeneratorFunction(t, n, context)
			} else {
				if !(n.IsAsyncFunction()) {
					j2gort.CheckState(false)
				}
				this_.InstrumentAsyncFunction(n, t.GetInput(), context)
			}
		}
	}
}

// Java method: InstrumentAsyncContext.instrumentReentrance(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InstrumentAsyncContext) priv_InstrumentAsyncContext_instrumentReentrance_3408(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var enclosingTry pkg_rhino.Node = this_.tryFunctionStack.PeekFirst()
	_ = enclosingTry
	if (enclosingTry != 0) && !(enclosingTry.IsTry()) {
		enclosingTry = pkg_rhino.Node(0)
	}
	var enclosingFunction pkg_rhino.Node = t.GetEnclosingFunction()
	_ = enclosingFunction
	if enclosingFunction == 0 {
		return
	} else if this_.alreadyInstrumented.Contains(enclosingFunction) {
		return
	} else if parent.IsFunction() {
		return
	} else if parent.IsReturn() && (enclosingTry == 0) {
		return
	}
	var functionContext *InstrumentAsyncContext_FunctionContext = j2gort.OrderedMapComputeIfAbsent(&this_.needsInstrumentation, enclosingFunction, nil, func(lp_12 pkg_rhino.Node) *InstrumentAsyncContext_FunctionContext {
		unused := lp_12
		_ = unused
		return this_.priv_InstrumentAsyncContext_createNewFunctionContext_3420(t.GetInput())
	})

	_ = functionContext
	if enclosingTry != 0 {
		this_.needsInstrumentation.Put(enclosingTry, functionContext)
	}
	if n.IsForAwaitOf() {
		this_.InstrumentForAwaitOf(t, n, parent, functionContext)
		return
	}
	var placeholder pkg_rhino.Node = pkg_rhino.IR_Empty()
	_ = placeholder
	var arg pkg_rhino.Node = n.GetFirstChild()
	_ = arg
	if arg == 0 {
		n.AddChildToBack(functionContext.Suspend_0().SrcrefTreeIfMissing(n))
	} else {
		arg.ReplaceWith(placeholder)
		placeholder.ReplaceWith(functionContext.Suspend_1(arg).SrcrefTreeIfMissing(arg))
	}
	n.ReplaceWith(placeholder)
	placeholder.ReplaceWith(functionContext.Resume_1(n).SrcrefTreeIfMissing(n))
}

// Java method: InstrumentAsyncContext.instrumentForAwaitOf(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent, com.google.javascript.jscomp.InstrumentAsyncContext.FunctionContext functionContext)
func (this_ *InstrumentAsyncContext) InstrumentForAwaitOf(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) {
	_ = t
	_ = n
	_ = parent
	_ = functionContext
	var placeholder pkg_rhino.Node = pkg_rhino.IR_Empty()
	_ = placeholder
	var arg pkg_rhino.Node = n.GetSecondChild()
	_ = arg
	arg.ReplaceWith(placeholder)
	placeholder.ReplaceWith(functionContext.Suspend_1(arg).SrcrefTreeIfMissing(arg))
	var body pkg_rhino.Node = n.GetLastChild()
	_ = body
	if !(body.IsBlock()) {
		j2gort.CheckState(false)
	}
	body.ReplaceWith(placeholder)
	placeholder.ReplaceWith(pkg_rhino.IR_Block_1_Node_2([]pkg_rhino.Node{pkg_rhino.IR_ExprResult(functionContext.Resume_0().SrcrefTreeIfMissing(arg)), pkg_rhino.IR_TryFinally(body, pkg_rhino.IR_Block_1_Node(pkg_rhino.IR_ExprResult(functionContext.Suspend_0().SrcrefTreeIfMissing(arg))))}))
	if !(parent.IsBlock()) {
		var block pkg_rhino.Node = pkg_rhino.IR_Block_0()
		_ = block
		n.ReplaceWith(block)
		block.AddChildToFront(n)
		n = block
	}
	pkg_rhino.IR_ExprResult(functionContext.Resume_0().SrcrefTreeIfMissing(arg)).InsertAfter(n)
}

// Java method: InstrumentAsyncContext.instrumentTry(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.InstrumentAsyncContext.FunctionContext functionContext)
func (this_ *InstrumentAsyncContext) InstrumentTry(n pkg_rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) {
	_ = n
	_ = functionContext
	var block pkg_rhino.Node = n.GetSecondChild()
	_ = block
	if block.HasChildren() {
		block = block.GetFirstChild().GetSecondChild()
	} else {
		block = block.GetNext()
	}
	if !(block.IsBlock()) {
		j2gort.CheckState(false)
	}
	block.AddChildToFront(pkg_rhino.IR_ExprResult(functionContext.Resume_0().SrcrefTreeIfMissing(block)))
}

// Java method: InstrumentAsyncContext.instrumentGeneratorFunction(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node f, com.google.javascript.jscomp.InstrumentAsyncContext.FunctionContext functionContext)
func (this_ *InstrumentAsyncContext) InstrumentGeneratorFunction(t *NodeTraversal, f pkg_rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) {
	_ = t
	_ = f
	_ = functionContext
	if !(f.GetLastChild().HasChildren()) {
		return
	}
	if this_.hasSuper.Contains(f) {
		this_.InstrumentGeneratorMethod(t, f, functionContext)
		return
	}
	f.SetIsGeneratorFunction(false)
	var generatorBody pkg_rhino.Node = f.GetLastChild().Detach()
	_ = generatorBody
	var inner pkg_rhino.Node = this_.astFactory.CreateEmptyGeneratorFunction(AstFactory_S_Type_1_Node(f))
	_ = inner
	if f.IsAsyncFunction() {
		f.SetIsAsyncFunction(false)
		inner.SetIsAsyncFunction(true)
	}
	inner.GetLastChild().ReplaceWith(generatorBody)
	generatorBody = this_.AddFinallySuspend(generatorBody, functionContext)
	generatorBody.AddChildToFront(pkg_rhino.IR_ExprResult(functionContext.Resume_0()))
	var newOuterBody pkg_rhino.Node = pkg_rhino.IR_Block_1_Node_2([]pkg_rhino.Node{this_.priv_InstrumentAsyncContext_createFactorySuspendCall_3419(functionContext.FactoryName()), pkg_rhino.IR_ReturnNode_1(this_.astFactory.CreateCallWithUnknownType(inner, []pkg_rhino.Node{}))})
	_ = newOuterBody
	this_.priv_InstrumentAsyncContext_addSuspendResumeVarsAfter_3416(newOuterBody.GetFirstChild(), functionContext)
	f.AddChildToBack(newOuterBody.SrcrefTreeIfMissing(generatorBody))
	var val_13 ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext
	InitThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(&val_13, newOuterBody, false, this_.uniqueIdSupplier.GetUniqueId(t.GetInput()))
	var thisContext *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext = &val_13
	_ = thisContext
	var val_14 ThisAndArgumentsReferenceUpdater
	InitThisAndArgumentsReferenceUpdater(&val_14, this_.compiler, thisContext, this_.astFactory)
	var updater *ThisAndArgumentsReferenceUpdater = &val_14
	_ = updater
	NodeTraversal_Traverse(this_.compiler, generatorBody, updater)
	thisContext.AddVarDeclarations(this_.compiler, this_.astFactory, t.GetCurrentScript())
	this_.compiler.ReportChangeToChangeScope(f)
	this_.compiler.ReportChangeToChangeScope(inner)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION(), this_.compiler)
}

// Java method: InstrumentAsyncContext.instrumentGeneratorMethod(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node f, com.google.javascript.jscomp.InstrumentAsyncContext.FunctionContext functionContext)
func (this_ *InstrumentAsyncContext) InstrumentGeneratorMethod(t *NodeTraversal, f pkg_rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) {
	_ = t
	_ = f
	_ = functionContext
	var innerMethodName string = (j2gort.StrValString(func() string {
		if f.GetParent().IsMemberFunctionDef() {
			return f.GetParent().GetString()
		}
		return ""
	}()) + "$jscomp$" + j2gort.StrValString(this_.uniqueIdSupplier.GetUniqueId(t.GetInput())))
	_ = innerMethodName
	var outerParams pkg_rhino.Node = f.GetSecondChild()
	_ = outerParams
	var innerParams pkg_rhino.Node = outerParams.CloneTree_0()
	_ = innerParams
	var innerFactoryName pkg_rhino.Node = functionContext.FactoryName()
	_ = innerFactoryName
	var outerFactoryName pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(("\u1d43\u1d9cfactory" + j2gort.StrValString(this_.uniqueIdSupplier.GetUniqueId(t.GetInput()))))
	_ = outerFactoryName
	innerParams.AddChildToFront(innerFactoryName)
	var unknownType AstFactory_Type = AstFactory_S_Type_1_Node(innerFactoryName)
	_ = unknownType
	var generatorBody pkg_rhino.Node = this_.AddFinallySuspend(f.GetLastChild(), functionContext)
	_ = generatorBody
	generatorBody.AddChildToFront(pkg_rhino.IR_ExprResult(functionContext.Resume_0()))
	var argumentsParam pkg_rhino.Node = pkg_rhino.Node(0)
	_ = argumentsParam
	var argumentsReference pkg_rhino.Node = pkg_rhino.Node(0)
	_ = argumentsReference
	var val_15 InstrumentAsyncContext_ArgumentsRenamer
	val_15.outer_ = this_
	InitInstrumentAsyncContext_ArgumentsRenamer(&val_15)
	var renamer *InstrumentAsyncContext_ArgumentsRenamer = &val_15
	_ = renamer
	NodeTraversal_Traverse(this_.compiler, generatorBody, renamer)
	if renamer.argumentsName != j2gort.NullString {
		argumentsReference = this_.astFactory.CreateArgumentsReference()
		argumentsParam = this_.astFactory.CreateName_2_string_AstFactory_Type(renamer.argumentsName, AstFactory_S_Type_1_Node(argumentsReference))
		argumentsParam.InsertAfter(innerFactoryName)
	}
	var outerBody pkg_rhino.Node = pkg_rhino.IR_Block_0()
	_ = outerBody
	generatorBody.ReplaceWith(outerBody)
	var innerCall pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(this_.astFactory.CreateGetProp(this_.astFactory.CreateThis(unknownType), innerMethodName, AstFactory_S_Type_1_Node(f)), []pkg_rhino.Node{})
	_ = innerCall
	var innerParam pkg_rhino.Node = innerParams.GetFirstChild()
	_ = innerParam
	innerCall.AddChildToBack(outerFactoryName.CloneTree_0())
	innerParam = innerParam.GetNext()
	if (innerParam != 0) && (innerParam == argumentsParam) {
		innerCall.AddChildToBack(argumentsReference)
	}
	{
		var outerParam pkg_rhino.Node = outerParams.GetFirstChild()
		_ = outerParam
		for ; outerParam != 0; outerParam = outerParam.GetNext() {
			var argParamPair *InstrumentAsyncContext_ArgParamPair = this_.priv_InstrumentAsyncContext_simplifyParameter_3413(outerParam.CloneTree_0(), t)
			_ = argParamPair
			outerParam.ReplaceWith(func() pkg_rhino.Node {
				var asgn_16 pkg_rhino.Node = argParamPair.parameter
				outerParam = asgn_16
				return asgn_16
			}())
			innerCall.AddChildToBack(argParamPair.argument)
		}
	}
	var newInnerMethod pkg_rhino.Node = this_.astFactory.CreateFunction("", innerParams, generatorBody, AstFactory_S_Type_1_Node(f))
	_ = newInnerMethod
	pkg_rhino.IR_MemberFunctionDef(innerMethodName, newInnerMethod).InsertAfter(f.GetParent())
	newInnerMethod.GetFirstChild().MakeNonIndexable()
	newInnerMethod.SetIsAsyncFunction(f.IsAsyncFunction())
	newInnerMethod.SetIsGeneratorFunction(true)
	f.SetIsGeneratorFunction(false)
	f.SetIsAsyncFunction(false)
	var fp pkg_rhino.Node = f.GetParent()
	_ = fp
	var np pkg_rhino.Node = newInnerMethod.GetParent()
	_ = np
	np.SetStaticMember(fp.IsStaticMember())
	var jsdoc *pkg_rhino.JSDocInfo = fp.GetJSDocInfo()
	_ = jsdoc
	if jsdoc != nil {
		var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
		_ = builder
		if jsdoc.IsNoCollapse() {
			builder.RecordNoCollapse()
		}
		if jsdoc.IsNoInline() {
			builder.RecordNoInline()
		}
		np.SetJSDocInfo(builder.Build_0())
	}
	outerBody.AddChildToBack(this_.priv_InstrumentAsyncContext_createFactorySuspendCall_3419(outerFactoryName))
	this_.priv_InstrumentAsyncContext_addSuspendResumeVarsAtStartOf_3417(generatorBody, innerFactoryName, functionContext)
	outerBody.AddChildToBack(pkg_rhino.IR_ReturnNode_1(innerCall))
	f.SrcrefTreeIfMissing(f)
	newInnerMethod.GetParent().SrcrefTreeIfMissing(f)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_MEMBER_DECLARATIONS(), this_.compiler)
	this_.compiler.ReportChangeToChangeScope(f)
	this_.compiler.ReportChangeToChangeScope(newInnerMethod)
	this_.compiler.ReportChangeToEnclosingScope(f.GetParent())
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION(), this_.compiler)
}

// Java method: InstrumentAsyncContext.simplifyParameter(com.google.javascript.rhino.Node param, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *InstrumentAsyncContext) priv_InstrumentAsyncContext_simplifyParameter_3413(param pkg_rhino.Node, t *NodeTraversal) *InstrumentAsyncContext_ArgParamPair {
	_ = param
	_ = t
	if param.IsRest() {
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS(), this_.compiler)
		var paramElement pkg_rhino.Node = param.RemoveFirstChild()
		_ = paramElement
		var simplified *InstrumentAsyncContext_ArgParamPair = this_.priv_InstrumentAsyncContext_simplifyParameter_3413(paramElement, t)
		_ = simplified
		return NewInstrumentAsyncContext_ArgParamPair(pkg_rhino.IR_IterSpread(simplified.argument).SrcrefTreeIfMissing(param), pkg_rhino.IR_IterRest(simplified.parameter).SrcrefTreeIfMissing(param))
	} else if param.IsDefaultValue() {
		return this_.priv_InstrumentAsyncContext_simplifyParameter_3413(param.RemoveFirstChild(), t)
	} else if param.IsArrayPattern() || param.IsObjectPattern() {
		var newParam pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(("$jscomp$param$" + j2gort.StrValString(this_.uniqueIdSupplier.GetUniqueId(t.GetInput()))), AstFactory_S_Type_1_Node(param)).SrcrefTreeIfMissing(param)
		_ = newParam
		return NewInstrumentAsyncContext_ArgParamPair(newParam.CloneTree_0(), newParam)
	} else if param.IsName() {
		var newParam_2 pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type((j2gort.StrValString(param.GetString()) + "$jscomp$param$" + j2gort.StrValString(this_.uniqueIdSupplier.GetUniqueId(t.GetInput()))), AstFactory_S_Type_1_Node(param)).SrcrefTreeIfMissing(param)
		_ = newParam_2
		return NewInstrumentAsyncContext_ArgParamPair(newParam_2.CloneTree_0(), newParam_2)
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected parameter: "+j2gort.StrVal(param))))
}

// Java method: InstrumentAsyncContext.instrumentAsyncFunction(com.google.javascript.rhino.Node f, com.google.javascript.jscomp.CompilerInput input, com.google.javascript.jscomp.InstrumentAsyncContext.FunctionContext functionContext)
func (this_ *InstrumentAsyncContext) InstrumentAsyncFunction(f pkg_rhino.Node, input *CompilerInput, functionContext *InstrumentAsyncContext_FunctionContext) {
	_ = f
	_ = input
	_ = functionContext
	var block pkg_rhino.Node = f.GetLastChild()
	_ = block
	var newBlock pkg_rhino.Node = this_.AddFinallySuspend(block, functionContext)
	_ = newBlock
	newBlock.AddChildToFront(this_.priv_InstrumentAsyncContext_createFactoryCall_3418(functionContext.factoryName, []pkg_rhino.Node{}).SrcrefTreeIfMissing(newBlock))
	this_.priv_InstrumentAsyncContext_addSuspendResumeVarsAfter_3416(newBlock.GetFirstChild(), functionContext)
	this_.compiler.ReportChangeToChangeScope(f)
}

// Java method: InstrumentAsyncContext.addFinallySuspend(com.google.javascript.rhino.Node block, com.google.javascript.jscomp.InstrumentAsyncContext.FunctionContext functionContext)
func (this_ *InstrumentAsyncContext) AddFinallySuspend(block pkg_rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) pkg_rhino.Node {
	_ = block
	_ = functionContext
	var placeholder pkg_rhino.Node = pkg_rhino.IR_Empty()
	_ = placeholder
	block.ReplaceWith(placeholder)
	var finallyBlock pkg_rhino.Node = pkg_rhino.IR_Block_1_Node(pkg_rhino.IR_ExprResult(functionContext.Suspend_0()))
	_ = finallyBlock
	var newBlock pkg_rhino.Node = pkg_rhino.IR_Block_1_Node(pkg_rhino.IR_TryFinally(block, finallyBlock)).SrcrefTreeIfMissing(block)
	_ = newBlock
	placeholder.ReplaceWith(newBlock)
	return newBlock
}

// Java method: InstrumentAsyncContext.addSuspendResumeVarsAfter(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.InstrumentAsyncContext.FunctionContext functionContext)
func (this_ *InstrumentAsyncContext) priv_InstrumentAsyncContext_addSuspendResumeVarsAfter_3416(node pkg_rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) {
	_ = node
	_ = functionContext
	var resumeCall pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(functionContext.factoryName.CloneNode_0(), []pkg_rhino.Node{this_.astFactory.CreateNumber(float64(int32(1)))})
	_ = resumeCall
	pkg_rhino.IR_Var_2(functionContext.resumeName, resumeCall).SrcrefTreeIfMissing(node).InsertAfter(node)
	var suspendCall pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(functionContext.factoryName.CloneNode_0(), []pkg_rhino.Node{})
	_ = suspendCall
	pkg_rhino.IR_Var_2(functionContext.suspendName, suspendCall).SrcrefTreeIfMissing(node).InsertAfter(node)
}

// Java method: InstrumentAsyncContext.addSuspendResumeVarsAtStartOf(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node factoryName, com.google.javascript.jscomp.InstrumentAsyncContext.FunctionContext functionContext)
func (this_ *InstrumentAsyncContext) priv_InstrumentAsyncContext_addSuspendResumeVarsAtStartOf_3417(node pkg_rhino.Node, factoryName pkg_rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) {
	_ = node
	_ = factoryName
	_ = functionContext
	var resumeCall pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(factoryName.CloneNode_0(), []pkg_rhino.Node{this_.astFactory.CreateNumber(float64(int32(1)))})
	_ = resumeCall
	node.AddChildToFront(pkg_rhino.IR_Var_2(functionContext.resumeName, resumeCall).SrcrefTreeIfMissing(node))
	var suspendCall pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(factoryName.CloneNode_0(), []pkg_rhino.Node{})
	_ = suspendCall
	node.AddChildToFront(pkg_rhino.IR_Var_2(functionContext.suspendName, suspendCall).SrcrefTreeIfMissing(node))
}

// Java static method: InstrumentAsyncContext.createVar(com.google.javascript.rhino.Node name, com.google.javascript.rhino.Node value)
func InstrumentAsyncContext_CreateVar(name pkg_rhino.Node, value pkg_rhino.Node) pkg_rhino.Node {
	_ = name
	_ = value
	return pkg_rhino.IR_Var_2(name, value)
}

// Java method: InstrumentAsyncContext.createFactoryCall(com.google.javascript.rhino.Node factoryName, com.google.javascript.rhino.Node[] arg)
func (this_ *InstrumentAsyncContext) priv_InstrumentAsyncContext_createFactoryCall_3418(factoryName pkg_rhino.Node, arg []pkg_rhino.Node) pkg_rhino.Node {
	_ = factoryName
	_ = arg
	var call pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(this_.astFactory.CreateQNameWithUnknownType_1_RuntimeJsLibManager_JsLibField(this_.start), arg)
	_ = call
	return InstrumentAsyncContext_CreateVar(factoryName, call)
}

// Java method: InstrumentAsyncContext.createFactorySuspendCall(com.google.javascript.rhino.Node factoryName)
func (this_ *InstrumentAsyncContext) priv_InstrumentAsyncContext_createFactorySuspendCall_3419(factoryName pkg_rhino.Node) pkg_rhino.Node {
	_ = factoryName
	return this_.priv_InstrumentAsyncContext_createFactoryCall_3418(factoryName, []pkg_rhino.Node{this_.astFactory.CreateNumber(float64(int32(1)))})
}

// Java method: InstrumentAsyncContext.createNewFunctionContext(com.google.javascript.jscomp.CompilerInput input)
func (this_ *InstrumentAsyncContext) priv_InstrumentAsyncContext_createNewFunctionContext_3420(input *CompilerInput) *InstrumentAsyncContext_FunctionContext {
	_ = input
	var uniqueSuffix string = this_.uniqueIdSupplier.GetUniqueId(input)
	_ = uniqueSuffix
	var suspendNode pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(("\u1d43\u1d9csuspend" + j2gort.StrValString(uniqueSuffix)))
	_ = suspendNode
	var resumeNode pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(("\u1d43\u1d9cresume" + j2gort.StrValString(uniqueSuffix)))
	_ = resumeNode
	var factoryName pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(("\u1d43\u1d9cfactory" + j2gort.StrValString(uniqueSuffix)))
	_ = factoryName
	return NewInstrumentAsyncContext_FunctionContext(uniqueSuffix, suspendNode, resumeNode, factoryName, this_.astFactory)
}

// Java class: com.google.javascript.jscomp.InstrumentAsyncContext.ArgParamPair
type InstrumentAsyncContext_ArgParamPair struct {
	argument  pkg_rhino.Node `gson:"argument"`
	parameter pkg_rhino.Node `gson:"parameter"`
	idHash_   atomic.Int32
}

func (this_ *InstrumentAsyncContext_ArgParamPair) BindSelf_(s any) {
	_ = s
}
func (this_ *InstrumentAsyncContext_ArgParamPair) IsJavaNil_() bool { return this_ == nil }
func (*InstrumentAsyncContext_ArgParamPair) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InstrumentAsyncContext$ArgParamPair", "ArgParamPair"
}
func (this_ *InstrumentAsyncContext_ArgParamPair) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InstrumentAsyncContext_ArgParamPair) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InstrumentAsyncContext.ArgParamPair", "com.google.javascript.jscomp.InstrumentAsyncContext$ArgParamPair", func(v_ any) bool {
		return j2gort.InstanceOf[*InstrumentAsyncContext_ArgParamPair](v_)
	})
}

func (this_ *InstrumentAsyncContext_ArgParamPair) initFields_InstrumentAsyncContext_ArgParamPair() {
	_ = this_
}

// Java constructor: ArgParamPair(com.google.javascript.rhino.Node argument, com.google.javascript.rhino.Node parameter)
func NewInstrumentAsyncContext_ArgParamPair(argument pkg_rhino.Node, parameter pkg_rhino.Node) *InstrumentAsyncContext_ArgParamPair {
	this_ := &InstrumentAsyncContext_ArgParamPair{}
	InitInstrumentAsyncContext_ArgParamPair(this_, argument, parameter)
	return this_
}

func InitInstrumentAsyncContext_ArgParamPair(this_ *InstrumentAsyncContext_ArgParamPair, argument pkg_rhino.Node, parameter pkg_rhino.Node) {
	_ = argument
	_ = parameter
	this_.initFields_InstrumentAsyncContext_ArgParamPair()
	this_.argument = argument
	this_.parameter = parameter
}

// Java class: com.google.javascript.jscomp.InstrumentAsyncContext.ArgumentsRenamer
type InstrumentAsyncContext_ArgumentsRenamer struct {
	outer_        *InstrumentAsyncContext
	argumentsName string `gson:"argumentsName"`
	idHash_       atomic.Int32
}

func (this_ *InstrumentAsyncContext_ArgumentsRenamer) BindSelf_(s any) {
	_ = s
}
func (this_ *InstrumentAsyncContext_ArgumentsRenamer) IsJavaNil_() bool { return this_ == nil }
func (*InstrumentAsyncContext_ArgumentsRenamer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InstrumentAsyncContext$ArgumentsRenamer", "ArgumentsRenamer"
}
func (this_ *InstrumentAsyncContext_ArgumentsRenamer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InstrumentAsyncContext_ArgumentsRenamer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InstrumentAsyncContext.ArgumentsRenamer", "com.google.javascript.jscomp.InstrumentAsyncContext$ArgumentsRenamer", func(v_ any) bool {
		return j2gort.InstanceOf[*InstrumentAsyncContext_ArgumentsRenamer](v_)
	})
}
func (this_ *InstrumentAsyncContext_ArgumentsRenamer) SetOuter_(o *InstrumentAsyncContext) {
	this_.outer_ = o
}
func (this_ *InstrumentAsyncContext_ArgumentsRenamer) GetOuter_() *InstrumentAsyncContext {
	return this_.outer_
}
func (*InstrumentAsyncContext_ArgumentsRenamer) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *InstrumentAsyncContext_ArgumentsRenamer) initFields_InstrumentAsyncContext_ArgumentsRenamer() {
	_ = this_
	this_.argumentsName = j2gort.NullString
}

// Java constructor: ArgumentsRenamer()
func NewInstrumentAsyncContext_ArgumentsRenamer(outer_ *InstrumentAsyncContext) *InstrumentAsyncContext_ArgumentsRenamer {
	this_ := &InstrumentAsyncContext_ArgumentsRenamer{}
	this_.outer_ = outer_
	InitInstrumentAsyncContext_ArgumentsRenamer(this_)
	return this_
}

func InitInstrumentAsyncContext_ArgumentsRenamer(this_ *InstrumentAsyncContext_ArgumentsRenamer) {
	this_.initFields_InstrumentAsyncContext_ArgumentsRenamer()
}

// Java method: ArgumentsRenamer.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent)
func (this_ *InstrumentAsyncContext_ArgumentsRenamer) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(n.IsFunction()) || n.IsArrowFunction())
}

// Java method: ArgumentsRenamer.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent)
func (this_ *InstrumentAsyncContext_ArgumentsRenamer) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() && (j2gort.CheckNotNullString(n.GetString()) == "arguments") {
		if this_.argumentsName == j2gort.NullString {
			this_.argumentsName = ("$jscomp$arguments$" + j2gort.StrValString(this_.outer_.uniqueIdSupplier.GetUniqueId(t.GetInput())))
		}
		n.SetString(this_.argumentsName)
		if this_.outer_.compiler.GetOptions().PreservesDetailedSourceInfo() {
			n.SetOriginalName("arguments")
		}
		n.MakeNonIndexable()
		n.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
}

// Java record: com.google.javascript.jscomp.InstrumentAsyncContext.FunctionContext
type InstrumentAsyncContext_FunctionContext struct {
	uniqueSuffix string         `gson:"uniqueSuffix"`
	astFactory   *AstFactory    `gson:"astFactory"`
	suspendName  pkg_rhino.Node `gson:"suspendName"`
	resumeName   pkg_rhino.Node `gson:"resumeName"`
	factoryName  pkg_rhino.Node `gson:"factoryName"`
}

func (this_ *InstrumentAsyncContext_FunctionContext) BindSelf_(s any) {
	_ = s
}
func (this_ *InstrumentAsyncContext_FunctionContext) IsJavaNil_() bool { return this_ == nil }
func (*InstrumentAsyncContext_FunctionContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InstrumentAsyncContext$FunctionContext", "FunctionContext"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InstrumentAsyncContext.FunctionContext", "com.google.javascript.jscomp.InstrumentAsyncContext$FunctionContext", func(v_ any) bool {
		return j2gort.InstanceOf[*InstrumentAsyncContext_FunctionContext](v_)
	})
}
func (*InstrumentAsyncContext_FunctionContext) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *InstrumentAsyncContext_FunctionContext) initFields_InstrumentAsyncContext_FunctionContext() {
	_ = this_
}

// Java constructor: FunctionContext(java.lang.String uniqueSuffix, com.google.javascript.rhino.Node suspendName, com.google.javascript.rhino.Node resumeName, com.google.javascript.rhino.Node factoryName, com.google.javascript.jscomp.AstFactory astFactory)
func NewInstrumentAsyncContext_FunctionContext(uniqueSuffix string, suspendName pkg_rhino.Node, resumeName pkg_rhino.Node, factoryName pkg_rhino.Node, astFactory *AstFactory) *InstrumentAsyncContext_FunctionContext {
	this_ := &InstrumentAsyncContext_FunctionContext{}
	InitInstrumentAsyncContext_FunctionContext(this_, uniqueSuffix, suspendName, resumeName, factoryName, astFactory)
	return this_
}

func InitInstrumentAsyncContext_FunctionContext(this_ *InstrumentAsyncContext_FunctionContext, uniqueSuffix string, suspendName pkg_rhino.Node, resumeName pkg_rhino.Node, factoryName pkg_rhino.Node, astFactory *AstFactory) {
	_ = uniqueSuffix
	_ = suspendName
	_ = resumeName
	_ = factoryName
	_ = astFactory
	this_.initFields_InstrumentAsyncContext_FunctionContext()
	this_.uniqueSuffix = uniqueSuffix
	this_.suspendName = suspendName
	this_.resumeName = resumeName
	this_.factoryName = factoryName
	this_.astFactory = astFactory
}

// Java method: FunctionContext.suspend()
func (this_ *InstrumentAsyncContext_FunctionContext) Suspend_0() pkg_rhino.Node {
	return this_.astFactory.CreateCall(this_.suspendName.CloneNode_0(), AstFactory_S_Type_1_Node(this_.suspendName), []pkg_rhino.Node{})
}

// Java method: FunctionContext.suspend(com.google.javascript.rhino.Node inner)
func (this_ *InstrumentAsyncContext_FunctionContext) Suspend_1(inner pkg_rhino.Node) pkg_rhino.Node {
	_ = inner
	return this_.astFactory.CreateCall(this_.suspendName.CloneNode_0(), AstFactory_S_Type_1_Node(inner), []pkg_rhino.Node{inner})
}

// Java method: FunctionContext.resume()
func (this_ *InstrumentAsyncContext_FunctionContext) Resume_0() pkg_rhino.Node {
	return this_.astFactory.CreateCall(this_.resumeName.CloneNode_0(), AstFactory_S_Type_1_Node(this_.resumeName), []pkg_rhino.Node{})
}

// Java method: FunctionContext.resume(com.google.javascript.rhino.Node inner)
func (this_ *InstrumentAsyncContext_FunctionContext) Resume_1(inner pkg_rhino.Node) pkg_rhino.Node {
	_ = inner
	return this_.astFactory.CreateCall(this_.resumeName.CloneNode_0(), AstFactory_S_Type_1_Node(inner), []pkg_rhino.Node{inner})
}

func (this_ *InstrumentAsyncContext_FunctionContext) UniqueSuffix() string { return this_.uniqueSuffix }
func (this_ *InstrumentAsyncContext_FunctionContext) SuspendName() pkg_rhino.Node {
	return this_.suspendName
}
func (this_ *InstrumentAsyncContext_FunctionContext) ResumeName() pkg_rhino.Node {
	return this_.resumeName
}
func (this_ *InstrumentAsyncContext_FunctionContext) FactoryName() pkg_rhino.Node {
	return this_.factoryName
}
func (this_ *InstrumentAsyncContext_FunctionContext) AstFactory() *AstFactory {
	return this_.astFactory
}
func (this_ *InstrumentAsyncContext_FunctionContext) Equals(other_ any) bool {
	o_, ok_ := other_.(*InstrumentAsyncContext_FunctionContext)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.uniqueSuffix != o_.uniqueSuffix {
		return false
	}
	if this_.suspendName != o_.suspendName {
		return false
	}
	if this_.resumeName != o_.resumeName {
		return false
	}
	if this_.factoryName != o_.factoryName {
		return false
	}
	if this_.astFactory != o_.astFactory {
		return false
	}
	return true
}
func (this_ *InstrumentAsyncContext_FunctionContext) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_17 := this_.uniqueSuffix
		h_ = 31 * h_
		if hs_17 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_17)
		}
	}
	h_ = 31*h_ + int32(this_.suspendName)
	h_ = 31*h_ + int32(this_.resumeName)
	h_ = 31*h_ + int32(this_.factoryName)
	{
		hp_18 := this_.astFactory
		h_ = 31 * h_
		if hp_18 != nil {
			h_ += j2gort.IdentityHash(hp_18)
		}
	}
	return h_
}
func (this_ *InstrumentAsyncContext_FunctionContext) String() string {
	_ = fmt.Sprint
	return "FunctionContext[" + "uniqueSuffix=" + j2gort.StrVal(j2gort.BoxString(this_.uniqueSuffix)) + ", " + "suspendName=" + j2gort.StrVal(j2gort.BoxNode(this_.suspendName)) + ", " + "resumeName=" + j2gort.StrVal(j2gort.BoxNode(this_.resumeName)) + ", " + "factoryName=" + j2gort.StrVal(j2gort.BoxNode(this_.factoryName)) + ", " + "astFactory=" + j2gort.StrVal(j2gort.PtrToAny(this_.astFactory)) + "]"
}
func (*InstrumentAsyncContext) CanReleaseSyntacticScopes_()                  {}
func (*InstrumentAsyncContext_ArgumentsRenamer) CanReleaseSyntacticScopes_() {}
