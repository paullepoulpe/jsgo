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
// Port of com/google/javascript/jscomp/Es6ConvertSuperConstructorCalls.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls
type Es6ConvertSuperConstructorCalls struct {
	jscompInherits           RuntimeJsLibManager_JsLibField                                 `gson:"jscompInherits"`
	jscompConstruct          RuntimeJsLibManager_JsLibField                                 `gson:"jscompConstruct"`
	compiler                 AbstractCompiler                                               `gson:"compiler"`
	constructorDataStack     j2gort.Deque[*Es6ConvertSuperConstructorCalls_ConstructorData] `gson:"constructorDataStack"`
	astFactory               *AstFactory                                                    `gson:"astFactory"`
	globalNamespace          *GlobalNamespace                                               `gson:"globalNamespace"`
	transpilationNamespace   pkg_rhino.StaticScope                                          `gson:"transpilationNamespace"`
	uniqueIdSupplier         *UniqueIdSupplier                                              `gson:"uniqueIdSupplier"`
	idHash_                  atomic.Int32
	es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation `gson:"es6SubclassTranspilation"`
}

func (this_ *Es6ConvertSuperConstructorCalls) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6ConvertSuperConstructorCalls) IsJavaNil_() bool { return this_ == nil }
func (*Es6ConvertSuperConstructorCalls) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls", "Es6ConvertSuperConstructorCalls"
}
func (this_ *Es6ConvertSuperConstructorCalls) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6ConvertSuperConstructorCalls) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls", "com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6ConvertSuperConstructorCalls](v_)
	})
}
func (*Es6ConvertSuperConstructorCalls) Implements_jscomp_NodeTraversal_Callback() {}

func StaticGet_Es6ConvertSuperConstructorCalls_TMP_ERROR() string  { return "$jscomp$tmp$error" }
func StaticGet_Es6ConvertSuperConstructorCalls_SUPER_THIS() string { return "$jscomp$super$this" }
func (this_ *Es6ConvertSuperConstructorCalls) initFields_Es6ConvertSuperConstructorCalls() {
	_ = this_
}

// Java constructor: Es6ConvertSuperConstructorCalls(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.CompilerOptions.Es6SubclassTranspilation es6SubclassTranspilation)
func NewEs6ConvertSuperConstructorCalls(compiler AbstractCompiler, es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation) *Es6ConvertSuperConstructorCalls {
	this_ := &Es6ConvertSuperConstructorCalls{}
	InitEs6ConvertSuperConstructorCalls(this_, compiler, es6SubclassTranspilation)
	return this_
}

func InitEs6ConvertSuperConstructorCalls(this_ *Es6ConvertSuperConstructorCalls, compiler AbstractCompiler, es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation) {
	_ = compiler
	_ = es6SubclassTranspilation
	this_.initFields_Es6ConvertSuperConstructorCalls()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.transpilationNamespace = compiler.GetTranspilationNamespace()
	this_.constructorDataStack = j2gort.Deque[*Es6ConvertSuperConstructorCalls_ConstructorData]{}
	this_.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
	var runtimeJsLibManager *RuntimeJsLibManager = compiler.GetRuntimeJsLibManager()
	_ = runtimeJsLibManager
	this_.jscompInherits = runtimeJsLibManager.GetJsLibField("$jscomp.inherits")
	this_.jscompConstruct = runtimeJsLibManager.GetJsLibField("$jscomp.construct")
	this_.es6SubclassTranspilation = es6SubclassTranspilation
}

// Java method: Es6ConvertSuperConstructorCalls.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6ConvertSuperConstructorCalls) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() && !(NodeUtil_GetFeatureSetOfScript(n).Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_SUPER())) {
		return false
	}
	if n.IsFunction() {
		this_.constructorDataStack.Push(NewEs6ConvertSuperConstructorCalls_ConstructorData(n))
	} else if n.IsSuper() {
		if !(n.IsFirstChildOf(parent)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
		}
		if parent.IsGetProp() {
			t.Report_3(n, StaticGet_TranspilationUtil_CANNOT_CONVERT_YET(), []string{"super access with no extends clause"})
			return false
		}
		if !(parent.IsCall()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
		}
		var constructorData *Es6ConvertSuperConstructorCalls_ConstructorData = j2gort.CheckNotNullPtr(this_.constructorDataStack.Peek())
		_ = constructorData
		constructorData.superCalls.Add_1(j2gort.BoxNode(parent))
	}
	return true
}

// Java method: Es6ConvertSuperConstructorCalls.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6ConvertSuperConstructorCalls) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var constructorData *Es6ConvertSuperConstructorCalls_ConstructorData = this_.constructorDataStack.Peek()
	_ = constructorData
	if (constructorData != nil) && (n == constructorData.constructor) {
		this_.constructorDataStack.Pop()
		this_.priv_Es6ConvertSuperConstructorCalls_visitSuper_2482(t, constructorData)
	}
}

// Java method: Es6ConvertSuperConstructorCalls.visitSuper(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls.ConstructorData constructorData)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_visitSuper_2482(t *NodeTraversal, constructorData *Es6ConvertSuperConstructorCalls_ConstructorData) {
	_ = t
	_ = constructorData
	var constructor pkg_rhino.Node = constructorData.constructor
	_ = constructor
	var superCalls j2gort.Ext_java_util_List = constructorData.superCalls
	_ = superCalls
	if superCalls.IsEmpty() {
		return
	}
	if constructor.IsFromExterns() {
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(superCalls))
			for cur_1.HasNext() {
				superCall := pkg_rhino.Node(cur_1.NextNode())
				_ = superCall
				var enclosingStatement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(superCall)
				_ = enclosingStatement
				var enclosingScope pkg_rhino.Node = enclosingStatement.GetParent()
				_ = enclosingScope
				enclosingStatement.Detach()
				this_.compiler.ReportChangeToEnclosingScope(enclosingScope)
			}
		}
		return
	}
	var uniqueSuperThisName string = ("$jscomp$super$this" + "$" + j2gort.StrValString(this_.uniqueIdSupplier.GetUniqueId(this_.compiler.GetInput(NodeUtil_GetInputId(constructor)))))
	_ = uniqueSuperThisName
	var superClassNameNode pkg_rhino.Node = this_.priv_Es6ConvertSuperConstructorCalls_getSuperClassQNameNode_2499(constructor)
	_ = superClassNameNode
	var superClassQName string = superClassNameNode.GetQualifiedName()
	_ = superClassQName
	var thisType AstFactory_Type = this_.priv_Es6ConvertSuperConstructorCalls_getTypeOfThisForConstructor_2498(constructorData.constructor)
	_ = thisType
	if this_.priv_Es6ConvertSuperConstructorCalls_isNativeObjectClass_2493(t, superClassQName) {
		this_.priv_Es6ConvertSuperConstructorCalls_replaceSuperCallsWithThis_2483(superCalls, thisType)
	} else if this_.priv_Es6ConvertSuperConstructorCalls_isKnownNativeClass_2495(t, superClassQName) {
		this_.priv_Es6ConvertSuperConstructorCalls_rewriteSuperCallsToJsCompConstructCalls_2484(constructor, superCalls, superClassNameNode, thisType, uniqueSuperThisName)
	} else if this_.priv_Es6ConvertSuperConstructorCalls_isNativeErrorClass_2494(t, superClassQName) {
		this_.priv_Es6ConvertSuperConstructorCalls_rewriteNativeErrorSuperCalls_2485(t, superCalls, superClassNameNode, thisType)
	} else if this_.priv_Es6ConvertSuperConstructorCalls_isKnownToReturnOnlyUndefined_2488(superClassQName) {
		this_.priv_Es6ConvertSuperConstructorCalls_rewriteSuperCalls_2486(superCalls, superClassNameNode, thisType)
	} else {
		this_.priv_Es6ConvertSuperConstructorCalls_rewriteSuperAndPreserveReturnedThis_2487(constructor, superCalls, superClassNameNode, thisType, uniqueSuperThisName)
	}
}

// Java method: Es6ConvertSuperConstructorCalls.replaceSuperCallsWithThis(java.util.List<com.google.javascript.rhino.Node> superCalls, com.google.javascript.jscomp.AstFactory.Type thisType)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_replaceSuperCallsWithThis_2483(superCalls j2gort.Ext_java_util_List, thisType AstFactory_Type) {
	_ = superCalls
	_ = thisType
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(superCalls))
		for cur_2.HasNext() {
			superCall := pkg_rhino.Node(cur_2.NextNode())
			_ = superCall
			var sideEffectArgs []pkg_rhino.Node = []pkg_rhino.Node{}
			_ = sideEffectArgs
			var callee pkg_rhino.Node = superCall.RemoveFirstChild()
			_ = callee
			if !(callee.IsSuper()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(callee))
			}
			for superCall.HasChildren() {
				var arg pkg_rhino.Node = superCall.RemoveFirstChild()
				_ = arg
				if arg.IsSpread() {
					arg = arg.RemoveFirstChild()
				}
				if this_.compiler.GetAstAnalyzer().MayHaveSideEffects(arg) {
					sideEffectArgs = append(sideEffectArgs, arg)
				}
			}
			var replacement pkg_rhino.Node = this_.astFactory.CreateThis(thisType).Srcref(superCall)
			_ = replacement
			if len(sideEffectArgs) != 0 {
				var exprs pkg_rhino.Node = sideEffectArgs[int32(0)]
				_ = exprs
				{
					var i int32 = int32(1)
					_ = i
					for ; i < int32(len(sideEffectArgs)); i = i + 1 {
						exprs = this_.astFactory.CreateComma(exprs, sideEffectArgs[i]).SrcrefTreeIfMissing(superCall)
					}
				}
				replacement = this_.astFactory.CreateComma(exprs, replacement).SrcrefTreeIfMissing(superCall)
			}
			superCall.ReplaceWith(replacement)
			this_.compiler.ReportChangeToEnclosingScope(replacement)
		}
	}
}

// Java method: Es6ConvertSuperConstructorCalls.rewriteSuperCallsToJsCompConstructCalls(com.google.javascript.rhino.Node constructor, java.util.List<com.google.javascript.rhino.Node> superCalls, com.google.javascript.rhino.Node superClassNameNode, com.google.javascript.jscomp.AstFactory.Type thisType, java.lang.String uniqueSuperThisName)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_rewriteSuperCallsToJsCompConstructCalls_2484(constructor pkg_rhino.Node, superCalls j2gort.Ext_java_util_List, superClassNameNode pkg_rhino.Node, thisType AstFactory_Type, uniqueSuperThisName string) {
	_ = constructor
	_ = superCalls
	_ = superClassNameNode
	_ = thisType
	_ = uniqueSuperThisName
	var constructorBody pkg_rhino.Node = j2gort.CheckNotNullNode(constructor.GetChildAtIndex(int32(2)))
	_ = constructorBody
	var firstStatement pkg_rhino.Node = constructorBody.GetFirstChild()
	_ = firstStatement
	j2gort.CheckNotNullNode(firstStatement)
	var firstSuperCall pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](superCalls.Get(int32(0)))
	_ = firstSuperCall
	if ((constructorBody.HasOneChild() && firstStatement.IsExprResult()) && firstStatement.HasOneChild()) && (firstStatement.GetFirstChild() == firstSuperCall) {
		if !(superCalls.Size() == int32(1)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(constructor))
		}
		firstStatement.ReplaceWith(this_.astFactory.CreateReturn(this_.priv_Es6ConvertSuperConstructorCalls_createJSCompConstructorCall_2491(superClassNameNode, firstSuperCall, thisType)))
	} else {
		var typeOfThis AstFactory_Type = this_.priv_Es6ConvertSuperConstructorCalls_getTypeOfThisForConstructor_2498(constructor)
		_ = typeOfThis
		this_.priv_Es6ConvertSuperConstructorCalls_updateThisToSuperThis_2497(typeOfThis, constructorBody, superCalls, uniqueSuperThisName)
		var declaration pkg_rhino.Node = this_.astFactory.CreateSingleVarNameDeclaration_1(uniqueSuperThisName).SrcrefTree(constructorBody)
		_ = declaration
		var insertBeforePoint pkg_rhino.Node = NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(constructorBody)
		_ = insertBeforePoint
		if insertBeforePoint != 0 {
			declaration.InsertBefore(insertBeforePoint)
		} else {
			constructorBody.AddChildToBack(declaration)
		}
		constructorBody.AddChildToBack(this_.astFactory.CreateReturn(this_.astFactory.CreateName_2_string_AstFactory_Type(uniqueSuperThisName, typeOfThis)).SrcrefTree(constructorBody))
		{
			cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(superCalls))
			for cur_3.HasNext() {
				superCall := pkg_rhino.Node(cur_3.NextNode())
				_ = superCall
				superCall.ReplaceWith(this_.astFactory.CreateAssign_2_Node_Node(this_.astFactory.CreateName_2_string_AstFactory_Type(uniqueSuperThisName, typeOfThis).Srcref(superCall), this_.priv_Es6ConvertSuperConstructorCalls_createJSCompConstructorCall_2491(superClassNameNode, superCall, thisType)).Srcref(superCall))
			}
		}
	}
	this_.compiler.ReportChangeToEnclosingScope(constructorBody)
}

// Java method: Es6ConvertSuperConstructorCalls.rewriteNativeErrorSuperCalls(com.google.javascript.jscomp.NodeTraversal t, java.util.List<com.google.javascript.rhino.Node> superCalls, com.google.javascript.rhino.Node superClassNameNode, com.google.javascript.jscomp.AstFactory.Type thisType)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_rewriteNativeErrorSuperCalls_2485(t *NodeTraversal, superCalls j2gort.Ext_java_util_List, superClassNameNode pkg_rhino.Node, thisType AstFactory_Type) {
	_ = t
	_ = superCalls
	_ = superClassNameNode
	_ = thisType
	{
		cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(superCalls))
		for cur_4.HasNext() {
			superCall := pkg_rhino.Node(cur_4.NextNode())
			_ = superCall
			var newSuperCall pkg_rhino.Node = this_.priv_Es6ConvertSuperConstructorCalls_createNewSuperCall_2490(superClassNameNode, superCall, thisType, AstFactory_S_Type_1_Node(superCall))
			_ = newSuperCall
			this_.priv_Es6ConvertSuperConstructorCalls_replaceNativeErrorSuperCall_2492(superCall, newSuperCall, t.GetInput())
		}
	}
}

// Java method: Es6ConvertSuperConstructorCalls.rewriteSuperCalls(java.util.List<com.google.javascript.rhino.Node> superCalls, com.google.javascript.rhino.Node superClassNameNode, com.google.javascript.jscomp.AstFactory.Type thisType)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_rewriteSuperCalls_2486(superCalls j2gort.Ext_java_util_List, superClassNameNode pkg_rhino.Node, thisType AstFactory_Type) {
	_ = superCalls
	_ = superClassNameNode
	_ = thisType
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(superCalls))
		for cur_5.HasNext() {
			superCall := pkg_rhino.Node(cur_5.NextNode())
			_ = superCall
			var newSuperCall pkg_rhino.Node = this_.priv_Es6ConvertSuperConstructorCalls_createNewSuperCall_2490(superClassNameNode, superCall, thisType, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()))
			_ = newSuperCall
			var superCallParent pkg_rhino.Node = superCall.GetParent()
			_ = superCallParent
			if superCallParent.HasOneChild() && NodeUtil_IsStatement(superCallParent) {
				superCall.ReplaceWith(newSuperCall)
			} else {
				superCall.ReplaceWith(this_.astFactory.CreateComma(newSuperCall, this_.astFactory.CreateThis(thisType)).SrcrefTreeIfMissing(superCall))
			}
			this_.compiler.ReportChangeToEnclosingScope(superCallParent)
		}
	}
}

// Java method: Es6ConvertSuperConstructorCalls.rewriteSuperAndPreserveReturnedThis(com.google.javascript.rhino.Node constructor, java.util.List<com.google.javascript.rhino.Node> superCalls, com.google.javascript.rhino.Node superClassNameNode, com.google.javascript.jscomp.AstFactory.Type thisType, java.lang.String uniqueSuperThisName)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_rewriteSuperAndPreserveReturnedThis_2487(constructor pkg_rhino.Node, superCalls j2gort.Ext_java_util_List, superClassNameNode pkg_rhino.Node, thisType AstFactory_Type, uniqueSuperThisName string) {
	_ = constructor
	_ = superCalls
	_ = superClassNameNode
	_ = thisType
	_ = uniqueSuperThisName
	if this_.es6SubclassTranspilation == StaticGet_CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT() {
		this_.priv_Es6ConvertSuperConstructorCalls_rewriteSuperCallsToJsCompConstructCalls_2484(constructor, superCalls, superClassNameNode, thisType, uniqueSuperThisName)
		return
	}
	var constructorBody pkg_rhino.Node = j2gort.CheckNotNullNode(constructor.GetChildAtIndex(int32(2)))
	_ = constructorBody
	var firstStatement pkg_rhino.Node = constructorBody.GetFirstChild()
	_ = firstStatement
	var firstSuperCall pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](superCalls.Get(int32(0)))
	_ = firstSuperCall
	if ((constructorBody.HasOneChild() && firstStatement.IsExprResult()) && firstStatement.HasOneChild()) && (firstStatement.GetFirstChild() == firstSuperCall) {
		if !(superCalls.Size() == int32(1)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(constructor))
		}
		var newReturn pkg_rhino.Node = this_.astFactory.CreateOr(this_.priv_Es6ConvertSuperConstructorCalls_createNewSuperCall_2490(superClassNameNode, j2gort.CastNode[pkg_rhino.Node](superCalls.Get(int32(0))), AstFactory_S_Type_1_Node(j2gort.CastNode[pkg_rhino.Node](superCalls.Get(int32(0)))), AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), this_.astFactory.CreateThis(thisType))
		_ = newReturn
		firstStatement.ReplaceWith(pkg_rhino.IR_ReturnNode_1(newReturn).SrcrefTreeIfMissing(firstStatement))
	} else {
		var typeOfThis AstFactory_Type = this_.priv_Es6ConvertSuperConstructorCalls_getTypeOfThisForConstructor_2498(constructor)
		_ = typeOfThis
		var declaration pkg_rhino.Node = pkg_rhino.IR_Var_1(this_.astFactory.CreateName_2_string_AstFactory_Type(uniqueSuperThisName, typeOfThis)).SrcrefTree(constructorBody)
		_ = declaration
		var insertBeforePoint pkg_rhino.Node = NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(constructorBody)
		_ = insertBeforePoint
		if insertBeforePoint != 0 {
			declaration.InsertBefore(insertBeforePoint)
		} else {
			constructorBody.AddChildToBack(declaration)
		}
		constructorBody.AddChildToBack(pkg_rhino.IR_ReturnNode_1(this_.astFactory.CreateName_2_string_AstFactory_Type(uniqueSuperThisName, typeOfThis)).SrcrefTree(constructorBody))
		this_.priv_Es6ConvertSuperConstructorCalls_updateThisToSuperThis_2497(typeOfThis, constructorBody, superCalls, uniqueSuperThisName)
		{
			cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(superCalls))
			for cur_6.HasNext() {
				superCall := pkg_rhino.Node(cur_6.NextNode())
				_ = superCall
				var newSuperCall pkg_rhino.Node = this_.priv_Es6ConvertSuperConstructorCalls_createNewSuperCall_2490(superClassNameNode, superCall, typeOfThis, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN()))
				_ = newSuperCall
				superCall.ReplaceWith(this_.astFactory.CreateAssign_2_Node_Node(this_.astFactory.CreateName_2_string_AstFactory_Type(uniqueSuperThisName, typeOfThis), this_.astFactory.CreateOr(newSuperCall, this_.astFactory.CreateThis(typeOfThis))).SrcrefTreeIfMissing(superCall))
			}
		}
	}
	this_.compiler.ReportChangeToEnclosingScope(constructorBody)
}

// Java method: Es6ConvertSuperConstructorCalls.isKnownToReturnOnlyUndefined(java.lang.String functionQName)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_isKnownToReturnOnlyUndefined_2488(functionQName string) bool {
	_ = functionQName
	if this_.globalNamespace == nil {
		return false
	}
	var globalName *GlobalNamespace_Name = j2gort.CastPtr[GlobalNamespace_Name](this_.globalNamespace.GetSlot(functionQName))
	_ = globalName
	if globalName == nil {
		return false
	}
	var declarationRef *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](globalName.GetDeclaration())
	_ = declarationRef
	if declarationRef == nil {
		{
			cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(globalName.GetRefs()))
			for cur_7.HasNext() {
				rawElem_8 := cur_7.Next()
				ref := j2gort.CastPtr[GlobalNamespace_Ref](rawElem_8)
				_ = ref
				if ref.IsSet() {
					declarationRef = ref
				}
			}
		}
	}
	if declarationRef == nil {
		return false
	}
	var declaredVarOrProp pkg_rhino.Node = declarationRef.GetNode()
	_ = declaredVarOrProp
	if declaredVarOrProp.IsFromExterns() {
		return false
	}
	var declaration pkg_rhino.Node = declaredVarOrProp.GetParent()
	_ = declaration
	var declaredValue pkg_rhino.Node = pkg_rhino.Node(0)
	_ = declaredValue
	if declaration.IsFunction() {
		declaredValue = declaration
	} else if NodeUtil_IsNameDeclaration(declaration) && declaredVarOrProp.IsName() {
		if declaredVarOrProp.HasChildren() {
			declaredValue = j2gort.CheckNotNullNode(declaredVarOrProp.GetFirstChild())
		} else {
			return false
		}
	} else if declaration.IsAssign() && (declaration.GetFirstChild() == declaredVarOrProp) {
		declaredValue = j2gort.CheckNotNullNode(declaration.GetSecondChild())
	} else if declaration.IsObjectLit() && declaredVarOrProp.HasOneChild() {
		declaredValue = j2gort.CheckNotNullNode(declaredVarOrProp.GetFirstChild())
	} else {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected declaration format:\n"+j2gort.StrValString(declaration.ToStringTree()))))
	}
	return this_.priv_Es6ConvertSuperConstructorCalls_isNodeKnownToOnlyReturnUndefined_2489(declaredValue)
}

// Java method: Es6ConvertSuperConstructorCalls.isNodeKnownToOnlyReturnUndefined(com.google.javascript.rhino.Node node)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_isNodeKnownToOnlyReturnUndefined_2489(node pkg_rhino.Node) bool {
	_ = node
	if node.IsFunction() {
		var functionBody pkg_rhino.Node = j2gort.CheckNotNullNode(node.GetChildAtIndex(int32(2)))
		_ = functionBody
		return !(NewEs6ConvertSuperConstructorCalls_UndefinedReturnValueCheck(this_).MayReturnDefinedValue(functionBody))
	} else if node.IsQualifiedName() {
		return this_.priv_Es6ConvertSuperConstructorCalls_isKnownToReturnOnlyUndefined_2488(node.GetQualifiedName())
	} else if node.IsHook() {
		var left pkg_rhino.Node = node.GetSecondChild()
		_ = left
		var right pkg_rhino.Node = node.GetLastChild()
		_ = right
		return (this_.priv_Es6ConvertSuperConstructorCalls_isNodeKnownToOnlyReturnUndefined_2489(left) && this_.priv_Es6ConvertSuperConstructorCalls_isNodeKnownToOnlyReturnUndefined_2489(right))
	} else {
		return false
	}
}

// Java method: Es6ConvertSuperConstructorCalls.createNewSuperCall(com.google.javascript.rhino.Node superClassQNameNode, com.google.javascript.rhino.Node superCall, com.google.javascript.jscomp.AstFactory.Type thisType, com.google.javascript.jscomp.AstFactory.Type resultType)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_createNewSuperCall_2490(superClassQNameNode pkg_rhino.Node, superCall pkg_rhino.Node, thisType AstFactory_Type, resultType AstFactory_Type) pkg_rhino.Node {
	_ = superClassQNameNode
	_ = superCall
	_ = thisType
	_ = resultType
	if !(superClassQNameNode.IsQualifiedName()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(superClassQNameNode))
	}
	if !(superCall.IsCall()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(superCall))
	}
	var callee pkg_rhino.Node = superCall.RemoveFirstChild()
	_ = callee
	if !(callee.IsSuper()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(callee))
	}
	var args j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = args
	var hasSpreadArg bool = false
	_ = hasSpreadArg
	for superCall.HasChildren() {
		var arg pkg_rhino.Node = superCall.RemoveFirstChild()
		_ = arg
		hasSpreadArg = (hasSpreadArg || arg.IsSpread())
		args.Add_1(j2gort.BoxNode(arg))
	}
	if hasSpreadArg {
		var superClassDotApply pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(superClassQNameNode.CloneTree_0(), "apply").SrcrefTree(callee)
		_ = superClassDotApply
		var newSuperCall pkg_rhino.Node = this_.astFactory.CreateCall(superClassDotApply, resultType, []pkg_rhino.Node{}).Srcref(superCall)
		_ = newSuperCall
		newSuperCall.AddChildToBack(this_.astFactory.CreateThis(thisType).Srcref(callee))
		newSuperCall.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), false)
		if Es6ConvertSuperConstructorCalls_IsSingleSpreadOfArguments(args) {
			newSuperCall.AddChildToBack(j2gort.CastNode[pkg_rhino.Node](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(args)).GetOnlyChild().Detach())
		} else {
			newSuperCall.AddChildToBack(this_.astFactory.CreateArraylit_1_Ext_java_lang_Iterable(args).Srcref(superCall))
		}
		return newSuperCall
	} else {
		var superClassDotCall pkg_rhino.Node = this_.astFactory.CreateGetProp(superClassQNameNode.CloneTree_0(), "call", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT())).SrcrefTree(callee)
		_ = superClassDotCall
		var newSuperCall_2 pkg_rhino.Node = this_.astFactory.CreateCall(superClassDotCall, resultType, []pkg_rhino.Node{}).Srcref(superCall)
		_ = newSuperCall_2
		newSuperCall_2.AddChildToBack(this_.astFactory.CreateThis(thisType).Srcref(callee))
		newSuperCall_2.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), false)
		{
			cur_9 := j2gort.NewIterCursor(args)
			for cur_9.HasNext() {
				arg_2 := pkg_rhino.Node(cur_9.NextNode())
				_ = arg_2
				newSuperCall_2.AddChildToBack(arg_2)
			}
		}
		return newSuperCall_2
	}
}

// Java method: Es6ConvertSuperConstructorCalls.createJSCompConstructorCall(com.google.javascript.rhino.Node superClassQNameNode, com.google.javascript.rhino.Node superCall, com.google.javascript.jscomp.AstFactory.Type thisType)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_createJSCompConstructorCall_2491(superClassQNameNode pkg_rhino.Node, superCall pkg_rhino.Node, thisType AstFactory_Type) pkg_rhino.Node {
	_ = superClassQNameNode
	_ = superCall
	_ = thisType
	if !(superClassQNameNode.IsQualifiedName()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(superClassQNameNode))
	}
	if !(superCall.IsCall()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(superCall))
	}
	var callee pkg_rhino.Node = j2gort.CheckNotNullNode(superCall.RemoveFirstChild())
	_ = callee
	if !(callee.IsSuper()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(callee))
	}
	var jscompDotConstruct pkg_rhino.Node = this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.transpilationNamespace, this_.jscompConstruct).SrcrefTree(callee)
	_ = jscompDotConstruct
	var superClassQName pkg_rhino.Node = superClassQNameNode.CloneTree_0()
	_ = superClassQName
	var superCallArgList j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = superCallArgList
	for superCall.HasChildren() {
		superCallArgList.Add_1(j2gort.BoxNode(superCall.RemoveFirstChild()))
	}
	var superArgs pkg_rhino.Node
	_ = superArgs
	if Es6ConvertSuperConstructorCalls_IsSingleSpreadOfArguments(superCallArgList) {
		superArgs = j2gort.CastNode[pkg_rhino.Node](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(superCallArgList)).GetOnlyChild().Detach()
	} else {
		superArgs = this_.astFactory.CreateArraylit_1_Ext_java_lang_Iterable(superCallArgList).Srcref(superCall)
	}
	var thisDotConstructor pkg_rhino.Node = this_.astFactory.CreateGetProp(this_.astFactory.CreateThis(thisType), "constructor", AstFactory_S_Type_1_Node(superClassQNameNode)).SrcrefTree(superCall)
	_ = thisDotConstructor
	return this_.astFactory.CreateCall(jscompDotConstruct, AstFactory_S_Type_1_Node(superCall), []pkg_rhino.Node{superClassQName, superArgs, thisDotConstructor}).Srcref(superCall)
}

// Java static method: Es6ConvertSuperConstructorCalls.isSingleSpreadOfArguments(java.util.List<com.google.javascript.rhino.Node> nodeList)
func Es6ConvertSuperConstructorCalls_IsSingleSpreadOfArguments(nodeList j2gort.Ext_java_util_List) bool {
	_ = nodeList
	return ((nodeList.Size() == int32(1)) && Es6ConvertSuperConstructorCalls_IsSpreadOfArguments(j2gort.CastNode[pkg_rhino.Node](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](nodeList)))))
}

// Java static method: Es6ConvertSuperConstructorCalls.isSpreadOfArguments(com.google.javascript.rhino.Node node)
func Es6ConvertSuperConstructorCalls_IsSpreadOfArguments(node pkg_rhino.Node) bool {
	_ = node
	return (node.IsSpread() && node.GetOnlyChild().MatchesName_1_string("arguments"))
}

// Java method: Es6ConvertSuperConstructorCalls.replaceNativeErrorSuperCall(com.google.javascript.rhino.Node superCall, com.google.javascript.rhino.Node newSuperCall, com.google.javascript.jscomp.CompilerInput compilerInput)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_replaceNativeErrorSuperCall_2492(superCall pkg_rhino.Node, newSuperCall pkg_rhino.Node, compilerInput *CompilerInput) {
	_ = superCall
	_ = newSuperCall
	_ = compilerInput
	var superStatement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(superCall)
	_ = superStatement
	var body pkg_rhino.Node = superStatement.GetParent()
	_ = body
	if !(body.IsBlock()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(body))
	}
	var thisType AstFactory_Type = AstFactory_S_Type_1_Node(newSuperCall)
	_ = thisType
	var tmpErrorName string = ("$jscomp$tmp$error" + "$" + j2gort.StrValString(this_.uniqueIdSupplier.GetUniqueId(compilerInput)))
	_ = tmpErrorName
	var getError pkg_rhino.Node = pkg_rhino.IR_Var_1(this_.astFactory.CreateName_2_string_AstFactory_Type(tmpErrorName, thisType)).SrcrefTreeIfMissing(superCall)
	_ = getError
	getError.InsertBefore(superStatement)
	var getTmpError pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(this_.astFactory.CreateName_2_string_AstFactory_Type(tmpErrorName, thisType), newSuperCall)
	_ = getTmpError
	var copyMessage pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(this_.astFactory.CreateGetProp(this_.astFactory.CreateThis(thisType), "message", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_STRING())), this_.astFactory.CreateGetProp(this_.astFactory.CreateName_2_string_AstFactory_Type(tmpErrorName, thisType), "message", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_STRING())))
	_ = copyMessage
	var setStack pkg_rhino.Node = this_.astFactory.CreateAnd(this_.astFactory.CreateIn(this_.astFactory.CreateString("stack"), this_.astFactory.CreateName_2_string_AstFactory_Type(tmpErrorName, thisType)), this_.astFactory.CreateAssign_2_Node_Node(this_.astFactory.CreateGetProp(this_.astFactory.CreateThis(thisType), "stack", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_STRING())), this_.astFactory.CreateGetProp(this_.astFactory.CreateName_2_string_AstFactory_Type(tmpErrorName, thisType), "stack", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_STRING()))))
	_ = setStack
	var superErrorExpr pkg_rhino.Node = this_.astFactory.CreateCommas(getTmpError, copyMessage, []pkg_rhino.Node{setStack, this_.astFactory.CreateThis(thisType)}).SrcrefTreeIfMissing(superCall)
	_ = superErrorExpr
	superCall.ReplaceWith(superErrorExpr)
	this_.compiler.ReportChangeToEnclosingScope(superErrorExpr)
}

// Java method: Es6ConvertSuperConstructorCalls.isNativeObjectClass(com.google.javascript.jscomp.NodeTraversal t, java.lang.String className)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_isNativeObjectClass_2493(t *NodeTraversal, className string) bool {
	_ = t
	_ = className
	return ((j2gort.CheckNotNullString(className) == "Object") && !(this_.priv_Es6ConvertSuperConstructorCalls_isDefinedInSources_2496(t, className)))
}

// Java method: Es6ConvertSuperConstructorCalls.isNativeErrorClass(com.google.javascript.jscomp.NodeTraversal t, java.lang.String superClassName)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_isNativeErrorClass_2494(t *NodeTraversal, superClassName string) bool {
	_ = t
	_ = superClassName
	switch j2gort.CheckNotNullString(superClassName) {
	case "AggregateError", "Error", "EvalError", "RangeError", "ReferenceError", "SyntaxError", "TypeError", "URIError":
		return !(this_.priv_Es6ConvertSuperConstructorCalls_isDefinedInSources_2496(t, superClassName))
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Es6ConvertSuperConstructorCalls.isKnownNativeClass(com.google.javascript.jscomp.NodeTraversal t, java.lang.String className)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_isKnownNativeClass_2495(t *NodeTraversal, className string) bool {
	_ = t
	_ = className
	switch j2gort.CheckNotNullString(className) {
	case "Array", "ArrayBuffer", "Boolean", "DataView", "Date", "Float32Array", "Function", "Generator", "GeneratorFunction", "Int16Array", "Int32Array", "Int8Array", "InternalError", "Iterator", "Map", "Number", "Object", "Promise", "Proxy", "RegExp", "Set", "String", "Symbol", "TypedArray", "Uint16Array", "Uint32Array", "Uint8Array", "Uint8ClampedArray", "WeakMap", "WeakSet":
		return !(this_.priv_Es6ConvertSuperConstructorCalls_isDefinedInSources_2496(t, className))
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Es6ConvertSuperConstructorCalls.isDefinedInSources(com.google.javascript.jscomp.NodeTraversal t, java.lang.String varName)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_isDefinedInSources_2496(t *NodeTraversal, varName string) bool {
	_ = t
	_ = varName
	var objectVar *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(varName))
	_ = objectVar
	return ((objectVar != nil) && !(objectVar.IsExtern()))
}

// Java method: Es6ConvertSuperConstructorCalls.updateThisToSuperThis(com.google.javascript.jscomp.AstFactory.Type typeOfThis, com.google.javascript.rhino.Node constructorBody, java.util.List<com.google.javascript.rhino.Node> superCalls, java.lang.String uniqueSuperThisName)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_updateThisToSuperThis_2497(typeOfThis AstFactory_Type, constructorBody pkg_rhino.Node, superCalls j2gort.Ext_java_util_List, uniqueSuperThisName string) {
	_ = typeOfThis
	_ = constructorBody
	_ = superCalls
	_ = uniqueSuperThisName
	var replaceThisWithSuperThis NodeTraversal_Callback = NewEs6ConvertSuperConstructorCalls__Anon_41(this_, superCalls, uniqueSuperThisName, typeOfThis)
	_ = replaceThisWithSuperThis
	NodeTraversal_Traverse(this_.compiler, constructorBody, replaceThisWithSuperThis)
}

// Java method: Es6ConvertSuperConstructorCalls.getTypeOfThisForConstructor(com.google.javascript.rhino.Node constructor)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_getTypeOfThisForConstructor_2498(constructor pkg_rhino.Node) AstFactory_Type {
	_ = constructor
	if !(constructor.IsFunction()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(constructor))
	}
	var constructorType *pkg_colors.Color = constructor.GetColor()
	_ = constructorType
	if (constructorType != nil) && constructorType.HasInstanceColors() {
		return AstFactory_S_Type_1_Color(pkg_colors.Color_CreateUnion_1_Color(constructorType.GetInstanceColorsArray()))
	} else {
		return AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())
	}
}

// Java method: Es6ConvertSuperConstructorCalls.getSuperClassQNameNode(com.google.javascript.rhino.Node constructor)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_getSuperClassQNameNode_2499(constructor pkg_rhino.Node) pkg_rhino.Node {
	_ = constructor
	var className string = NodeUtil_GetNameNode(constructor).GetQualifiedName()
	_ = className
	var constructorStatement pkg_rhino.Node = j2gort.CheckNotNullNode(NodeUtil_GetEnclosingStatement(constructor))
	_ = constructorStatement
	var superClassNameNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = superClassNameNode
	{
		var statement pkg_rhino.Node = constructorStatement.GetNext()
		_ = statement
		for ; statement != 0; statement = statement.GetNext() {
			superClassNameNode = this_.priv_Es6ConvertSuperConstructorCalls_getSuperClassNameNodeIfIsInheritsStatement_2500(statement, className)
			if superClassNameNode != 0 {
				break
			}
		}
	}
	return j2gort.CheckNotNullNode(superClassNameNode)
}

// Java method: Es6ConvertSuperConstructorCalls.getSuperClassNameNodeIfIsInheritsStatement(com.google.javascript.rhino.Node statement, java.lang.String className)
func (this_ *Es6ConvertSuperConstructorCalls) priv_Es6ConvertSuperConstructorCalls_getSuperClassNameNodeIfIsInheritsStatement_2500(statement pkg_rhino.Node, className string) pkg_rhino.Node {
	_ = statement
	_ = className
	if !(statement.IsExprResult()) {
		return pkg_rhino.Node(0)
	}
	var callNode pkg_rhino.Node = statement.GetFirstChild()
	_ = callNode
	if !(callNode.IsCall()) {
		return pkg_rhino.Node(0)
	}
	var jscompDotInherits pkg_rhino.Node = callNode.GetFirstChild()
	_ = jscompDotInherits
	if !(this_.jscompInherits.Matches(jscompDotInherits)) {
		return pkg_rhino.Node(0)
	}
	var classNameNode pkg_rhino.Node = j2gort.CheckNotNullNode(jscompDotInherits.GetNext())
	_ = classNameNode
	if classNameNode.MatchesQualifiedName_1_string(className) {
		return j2gort.CheckNotNullNode(classNameNode.GetNext())
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java method: Es6ConvertSuperConstructorCalls.setGlobalNamespace(com.google.javascript.jscomp.GlobalNamespace globalNamespace)
func (this_ *Es6ConvertSuperConstructorCalls) SetGlobalNamespace(globalNamespace *GlobalNamespace) {
	_ = globalNamespace
	this_.globalNamespace = globalNamespace
}

// Java class: com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls.ConstructorData
type Es6ConvertSuperConstructorCalls_ConstructorData struct {
	superCalls  j2gort.Ext_java_util_List `gson:"superCalls"`
	constructor pkg_rhino.Node            `gson:"constructor"`
	idHash_     atomic.Int32
}

func (this_ *Es6ConvertSuperConstructorCalls_ConstructorData) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6ConvertSuperConstructorCalls_ConstructorData) IsJavaNil_() bool { return this_ == nil }
func (*Es6ConvertSuperConstructorCalls_ConstructorData) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls$ConstructorData", "ConstructorData"
}
func (this_ *Es6ConvertSuperConstructorCalls_ConstructorData) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6ConvertSuperConstructorCalls_ConstructorData) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls.ConstructorData", "com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls$ConstructorData", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6ConvertSuperConstructorCalls_ConstructorData](v_)
	})
}

func (this_ *Es6ConvertSuperConstructorCalls_ConstructorData) initFields_Es6ConvertSuperConstructorCalls_ConstructorData() {
	_ = this_
}

// Java constructor: ConstructorData(com.google.javascript.rhino.Node constructor)
func NewEs6ConvertSuperConstructorCalls_ConstructorData(constructor pkg_rhino.Node) *Es6ConvertSuperConstructorCalls_ConstructorData {
	this_ := &Es6ConvertSuperConstructorCalls_ConstructorData{}
	InitEs6ConvertSuperConstructorCalls_ConstructorData(this_, constructor)
	return this_
}

func InitEs6ConvertSuperConstructorCalls_ConstructorData(this_ *Es6ConvertSuperConstructorCalls_ConstructorData, constructor pkg_rhino.Node) {
	_ = constructor
	this_.initFields_Es6ConvertSuperConstructorCalls_ConstructorData()
	this_.constructor = constructor
	this_.superCalls = j2gort.NewExt_java_util_ArrayList()
}

// Java class: com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls.UndefinedReturnValueCheck
type Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck struct {
	outer_              *Es6ConvertSuperConstructorCalls
	idHash_             atomic.Int32
	foundNonEmptyReturn bool `gson:"foundNonEmptyReturn"`
}

func (this_ *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) IsJavaNil_() bool {
	return this_ == nil
}
func (*Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls$UndefinedReturnValueCheck", "UndefinedReturnValueCheck"
}
func (this_ *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls.UndefinedReturnValueCheck", "com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls$UndefinedReturnValueCheck", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck](v_)
	})
}
func (this_ *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) SetOuter_(o *Es6ConvertSuperConstructorCalls) {
	this_.outer_ = o
}
func (this_ *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) GetOuter_() *Es6ConvertSuperConstructorCalls {
	return this_.outer_
}

func (this_ *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) initFields_Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck() {
	_ = this_
}

// Java constructor: UndefinedReturnValueCheck()
func NewEs6ConvertSuperConstructorCalls_UndefinedReturnValueCheck(outer_ *Es6ConvertSuperConstructorCalls) *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck {
	this_ := &Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck{}
	this_.outer_ = outer_
	InitEs6ConvertSuperConstructorCalls_UndefinedReturnValueCheck(this_)
	return this_
}

func InitEs6ConvertSuperConstructorCalls_UndefinedReturnValueCheck(this_ *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) {
	this_.initFields_Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck()
}

// Java method: UndefinedReturnValueCheck.mayReturnDefinedValue(com.google.javascript.rhino.Node functionBody)
func (this_ *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) MayReturnDefinedValue(functionBody pkg_rhino.Node) bool {
	_ = functionBody
	this_.foundNonEmptyReturn = false
	var checkForDefinedReturnValue NodeTraversal_Callback = NewEs6ConvertSuperConstructorCalls__Anon_40(this_)
	_ = checkForDefinedReturnValue
	NodeTraversal_Traverse(this_.outer_.compiler, functionBody, checkForDefinedReturnValue)
	return this_.foundNonEmptyReturn
}

// Java class:
type Es6ConvertSuperConstructorCalls__Anon_41 struct {
	outer_                    *Es6ConvertSuperConstructorCalls
	cap_0_superCalls          j2gort.Ext_java_util_List
	cap_1_uniqueSuperThisName string
	cap_2_typeOfThis          AstFactory_Type
	idHash_                   atomic.Int32
}

func (this_ *Es6ConvertSuperConstructorCalls__Anon_41) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6ConvertSuperConstructorCalls__Anon_41) IsJavaNil_() bool { return this_ == nil }
func (*Es6ConvertSuperConstructorCalls__Anon_41) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls$1", ""
}
func (this_ *Es6ConvertSuperConstructorCalls__Anon_41) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6ConvertSuperConstructorCalls__Anon_41) HasIdentityEquals_() {}
func (this_ *Es6ConvertSuperConstructorCalls__Anon_41) SetOuter_(o *Es6ConvertSuperConstructorCalls) {
	this_.outer_ = o
}
func (this_ *Es6ConvertSuperConstructorCalls__Anon_41) GetOuter_() *Es6ConvertSuperConstructorCalls {
	return this_.outer_
}
func (*Es6ConvertSuperConstructorCalls__Anon_41) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *Es6ConvertSuperConstructorCalls__Anon_41) initFields_Es6ConvertSuperConstructorCalls__Anon_41() {
	_ = this_
}

func NewEs6ConvertSuperConstructorCalls__Anon_41(outer_ *Es6ConvertSuperConstructorCalls, cap_p_superCalls j2gort.Ext_java_util_List, cap_p_uniqueSuperThisName string, cap_p_typeOfThis AstFactory_Type) *Es6ConvertSuperConstructorCalls__Anon_41 {
	this_ := &Es6ConvertSuperConstructorCalls__Anon_41{}
	this_.outer_ = outer_
	this_.cap_0_superCalls = cap_p_superCalls
	this_.cap_1_uniqueSuperThisName = cap_p_uniqueSuperThisName
	this_.cap_2_typeOfThis = cap_p_typeOfThis
	this_.initFields_Es6ConvertSuperConstructorCalls__Anon_41()
	return this_
}

// Java method: .shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6ConvertSuperConstructorCalls__Anon_41) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if this_.cap_0_superCalls.Contains(j2gort.BoxNode(n)) {
		return false
	} else if n.IsFunction() && !(n.IsArrowFunction()) {
		return false
	} else {
		return true
	}
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6ConvertSuperConstructorCalls__Anon_41) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsThis() {
		var superThis pkg_rhino.Node = this_.outer_.astFactory.CreateName_2_string_AstFactory_Type(this_.cap_1_uniqueSuperThisName, AstFactory_S_Type_1_Node(n)).Srcref(n)
		_ = superThis
		n.ReplaceWith(superThis)
	} else if n.IsReturn() && !(n.HasChildren()) {
		n.AddChildToFront(this_.outer_.astFactory.CreateName_2_string_AstFactory_Type(this_.cap_1_uniqueSuperThisName, this_.cap_2_typeOfThis).Srcref(n))
	}
}

// Java class:
type Es6ConvertSuperConstructorCalls__Anon_40 struct {
	NodeTraversal_AbstractShallowCallbackBase
	outer_ *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck
}

func (this_ *Es6ConvertSuperConstructorCalls__Anon_40) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6ConvertSuperConstructorCalls__Anon_40) IsJavaNil_() bool { return this_ == nil }
func (*Es6ConvertSuperConstructorCalls__Anon_40) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6ConvertSuperConstructorCalls$UndefinedReturnValueCheck$1", ""
}
func (this_ *Es6ConvertSuperConstructorCalls__Anon_40) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractShallowCallbackBase.IdentityHash_()
}
func (*Es6ConvertSuperConstructorCalls__Anon_40) HasIdentityEquals_() {}
func (this_ *Es6ConvertSuperConstructorCalls__Anon_40) SetOuter_(o *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) {
	this_.outer_ = o
}
func (this_ *Es6ConvertSuperConstructorCalls__Anon_40) GetOuter_() *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck {
	return this_.outer_
}
func (*Es6ConvertSuperConstructorCalls__Anon_40) Implements_jscomp_NodeTraversal_AbstractShallowCallback() {
}
func (*Es6ConvertSuperConstructorCalls__Anon_40) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *Es6ConvertSuperConstructorCalls__Anon_40) initFields_Es6ConvertSuperConstructorCalls__Anon_40() {
	_ = this_
}

func NewEs6ConvertSuperConstructorCalls__Anon_40(outer_ *Es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) *Es6ConvertSuperConstructorCalls__Anon_40 {
	this_ := &Es6ConvertSuperConstructorCalls__Anon_40{}
	this_.outer_ = outer_
	InitNodeTraversal_AbstractShallowCallback(&this_.NodeTraversal_AbstractShallowCallbackBase)
	this_.initFields_Es6ConvertSuperConstructorCalls__Anon_40()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6ConvertSuperConstructorCalls__Anon_40) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(this_.outer_.foundNonEmptyReturn) {
		if (n.IsReturn() && n.HasChildren()) && !(n.GetFirstChild().MatchesName_1_string("undefined")) {
			this_.outer_.foundNonEmptyReturn = true
		}
	}
}
func (*Es6ConvertSuperConstructorCalls) CanReleaseSyntacticScopes_()          {}
func (*Es6ConvertSuperConstructorCalls__Anon_41) CanReleaseSyntacticScopes_() {}
func (*Es6ConvertSuperConstructorCalls__Anon_40) CanReleaseSyntacticScopes_() {}
