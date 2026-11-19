/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RewriteAsyncFunctions.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RewriteAsyncFunctions
type RewriteAsyncFunctions struct {
	namespace    pkg_rhino.StaticScope                              `gson:"namespace"`
	contextStack j2gort.Deque[RewriteAsyncFunctions_LexicalContext] `gson:"contextStack"`
	compiler     AbstractCompiler                                   `gson:"compiler"`
	astFactory   *AstFactory                                        `gson:"astFactory"`
	idHash_      atomic.Int32
}

func (this_ *RewriteAsyncFunctions) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteAsyncFunctions) IsJavaNil_() bool { return this_ == nil }
func (*RewriteAsyncFunctions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteAsyncFunctions", "RewriteAsyncFunctions"
}
func (this_ *RewriteAsyncFunctions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteAsyncFunctions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteAsyncFunctions", "com.google.javascript.jscomp.RewriteAsyncFunctions", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteAsyncFunctions](v_)
	})
}
func (*RewriteAsyncFunctions) Implements_jscomp_NodeTraversal_Callback() {}
func (*RewriteAsyncFunctions) Implements_jscomp_CompilerPass()           {}

func StaticGet_RewriteAsyncFunctions_ASYNC_ARGUMENTS() string  { return "$jscomp$async$arguments$" }
func StaticGet_RewriteAsyncFunctions_ASYNC_THIS() string       { return "$jscomp$async$this$" }
func StaticGet_RewriteAsyncFunctions_ASYNC_NEW_TARGET() string { return "$jscomp$async$new$target$" }
func StaticGet_RewriteAsyncFunctions_ASYNC_SUPER_PROP_GETTER_PREFIX() string {
	return "$jscomp$async$super$get$"
}

var static_RewriteAsyncFunctions_transpiledFeatures *pkg_parser.FeatureSet

func StaticGet_RewriteAsyncFunctions_TranspiledFeatures() *pkg_parser.FeatureSet {
	clinit_RewriteAsyncFunctions()
	return static_RewriteAsyncFunctions_transpiledFeatures
}
func StaticSet_RewriteAsyncFunctions_TranspiledFeatures(v *pkg_parser.FeatureSet) {
	clinit_RewriteAsyncFunctions()
	static_RewriteAsyncFunctions_transpiledFeatures = v
}

var clinitDone_RewriteAsyncFunctions uint32
var clinitState_RewriteAsyncFunctions j2gort.ClinitState

func clinit_RewriteAsyncFunctions() {
	if atomic.LoadUint32(&clinitDone_RewriteAsyncFunctions) != 0 {
		return
	}
	clinitSlow_RewriteAsyncFunctions()
}

func clinitSlow_RewriteAsyncFunctions() {
	if !j2gort.EnterClinit(&clinitState_RewriteAsyncFunctions, &clinitDone_RewriteAsyncFunctions) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RewriteAsyncFunctions, &clinitDone_RewriteAsyncFunctions)
	static_RewriteAsyncFunctions_transpiledFeatures = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS())
}

func (this_ *RewriteAsyncFunctions) initFields_RewriteAsyncFunctions() {
	_ = this_
}

// Java constructor: RewriteAsyncFunctions(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.AstFactory astFactory, com.google.javascript.rhino.StaticScope namespace)
func NewRewriteAsyncFunctions(compiler AbstractCompiler, astFactory *AstFactory, namespace pkg_rhino.StaticScope) *RewriteAsyncFunctions {
	clinit_RewriteAsyncFunctions()
	this_ := &RewriteAsyncFunctions{}
	InitRewriteAsyncFunctions(this_, compiler, astFactory, namespace)
	return this_
}

func InitRewriteAsyncFunctions(this_ *RewriteAsyncFunctions, compiler AbstractCompiler, astFactory *AstFactory, namespace pkg_rhino.StaticScope) {
	_ = compiler
	_ = astFactory
	_ = namespace
	this_.initFields_RewriteAsyncFunctions()
	this_.compiler = j2gort.CheckNotNullIface(compiler)
	this_.astFactory = j2gort.CheckNotNullPtr(astFactory)
	this_.namespace = j2gort.CheckNotNullIface(namespace)
	this_.contextStack = j2gort.Deque[RewriteAsyncFunctions_LexicalContext]{}
}

// Java static method: RewriteAsyncFunctions.create(com.google.javascript.jscomp.AbstractCompiler compiler)
func RewriteAsyncFunctions_Create(compiler AbstractCompiler) *RewriteAsyncFunctions {
	clinit_RewriteAsyncFunctions()
	_ = compiler
	var astFactory *AstFactory = compiler.CreateAstFactory()
	_ = astFactory
	var namespace pkg_rhino.StaticScope = compiler.GetTranspilationNamespace()
	_ = namespace
	return NewRewriteAsyncFunctions(compiler, astFactory, namespace)
}

// Java method: RewriteAsyncFunctions.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RewriteAsyncFunctions) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	TranspilationPasses_ProcessTranspile(this_.compiler, root, StaticGet_RewriteAsyncFunctions_TranspiledFeatures(), []NodeTraversal_Callback{this_})
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, StaticGet_RewriteAsyncFunctions_TranspiledFeatures())
}

// Java method: RewriteAsyncFunctions.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteAsyncFunctions) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if parent == 0 {
		if !(this_.contextStack.IsEmpty()) {
			j2gort.CheckState(false)
		}
		var uniqueId string = this_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		_ = uniqueId
		this_.contextStack.Push(NewRewriteAsyncFunctions_RootContext(this_, n, uniqueId))
	} else {
		var parentContext RewriteAsyncFunctions_LexicalContext = this_.contextStack.Peek()
		_ = parentContext
		var nodeContext RewriteAsyncFunctions_LexicalContext = parentContext.GetContextForNode(t, n)
		_ = nodeContext
		if !(nodeContext == parentContext || ((nodeContext == nil || nodeContext.IsJavaNil_()) && (parentContext == nil || parentContext.IsJavaNil_()))) {
			this_.contextStack.Push(nodeContext)
		}
	}
	return true
}

// Java method: RewriteAsyncFunctions.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteAsyncFunctions) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var context RewriteAsyncFunctions_LexicalContext = this_.contextStack.Peek()
	_ = context
	context.Visit(t, n)
	if context.GetRewriteAsyncFunctions_LexicalContextBase().GetContextRootNode() == n {
		this_.contextStack.Pop()
	}
}

// Java method: RewriteAsyncFunctions.convertAsyncFunction(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.RewriteAsyncFunctions.FunctionContext functionContext)
func (this_ *RewriteAsyncFunctions) priv_RewriteAsyncFunctions_convertAsyncFunction_4957(t *NodeTraversal, functionContext *RewriteAsyncFunctions_FunctionContext) {
	_ = t
	_ = functionContext
	var originalFunction pkg_rhino.Node = functionContext.GetContextRootNode()
	_ = originalFunction
	originalFunction.SetIsAsyncFunction(false)
	var originalBody pkg_rhino.Node = originalFunction.GetLastChild()
	_ = originalBody
	if originalFunction.IsFromExterns() {
		if !(NodeUtil_IsEmptyBlock(originalBody)) {
			originalBody.ReplaceWith(this_.astFactory.CreateBlock([]pkg_rhino.Node{}))
			NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(originalBody, this_.compiler)
		}
		return
	}
	var newBody pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{})
	_ = newBody
	originalBody.ReplaceWith(newBody)
	if functionContext.mustAddAsyncThisVariable {
		newBody.AddChildToBack(this_.astFactory.CreateSingleConstNameDeclaration(("$jscomp$async$this$" + j2gort.StrValString(functionContext.asyncThisAndArgumentsContext.GetRewriteAsyncFunctions_LexicalContextBase().uniqueId)), this_.astFactory.CreateThis(functionContext.typeOfThis)))
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), this_.compiler)
	}
	if functionContext.mustAddAsyncArgumentsVariable {
		newBody.AddChildToBack(this_.astFactory.CreateArgumentsAliasDeclaration(("$jscomp$async$arguments$" + j2gort.StrValString(functionContext.asyncThisAndArgumentsContext.GetRewriteAsyncFunctions_LexicalContextBase().uniqueId))))
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), this_.compiler)
	}
	if functionContext.mustAddAsyncNewTargetVariable {
		newBody.AddChildToBack(this_.astFactory.CreateSingleConstNameDeclaration(("$jscomp$async$new$target$" + j2gort.StrValString(functionContext.asyncThisAndArgumentsContext.GetRewriteAsyncFunctions_LexicalContextBase().uniqueId)), this_.astFactory.CreateNewTarget(functionContext.typeOfNewTarget)))
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), this_.compiler)
	}
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(functionContext.superPropertyWrappers.priv_SuperPropertyWrappers_asCollection_4963()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			superPropertyWrapperInfo := j2gort.CastPtr[RewriteAsyncFunctions_SuperPropertyWrapperInfo](rawElem_2)
			_ = superPropertyWrapperInfo
			var arrowFunction pkg_rhino.Node = functionContext.priv_FunctionContext_createWrapperArrowFunction_4974(superPropertyWrapperInfo)
			_ = arrowFunction
			var arrowFunctionDeclarationStatement pkg_rhino.Node = this_.astFactory.CreateSingleConstNameDeclaration(superPropertyWrapperInfo.wrapperFunctionName, arrowFunction)
			_ = arrowFunctionDeclarationStatement
			newBody.AddChildToBack(arrowFunctionDeclarationStatement)
			this_.compiler.ReportChangeToChangeScope(arrowFunction)
			var enclosingScript pkg_rhino.Node = t.GetCurrentScript()
			_ = enclosingScript
			NodeUtil_AddFeatureToScript(enclosingScript, pkg_parser.StaticGet_FeatureSet_Feature_ARROW_FUNCTIONS(), this_.compiler)
			NodeUtil_AddFeatureToScript(enclosingScript, pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), this_.compiler)
		}
	}
	if !(originalBody.IsBlock()) {
		originalBody = this_.astFactory.CreateBlock([]pkg_rhino.Node{this_.astFactory.CreateReturn(originalBody)}).SrcrefTreeIfMissing(originalBody)
	}
	var generatorFunction pkg_rhino.Node = this_.astFactory.CreateZeroArgGeneratorFunction("", originalBody, nil)
	_ = generatorFunction
	this_.compiler.ReportChangeToChangeScope(generatorFunction)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS(), this_.compiler)
	newBody.AddChildToBack(this_.astFactory.CreateReturn(this_.astFactory.CreateJscompAsyncExecutePromiseGeneratorFunctionCall(this_.namespace, generatorFunction)))
	newBody.SrcrefTreeIfMissing(originalBody)
	this_.compiler.ReportChangeToEnclosingScope(newBody)
}

// Java class: com.google.javascript.jscomp.RewriteAsyncFunctions.SuperPropertyWrapperInfo
type RewriteAsyncFunctions_SuperPropertyWrapperInfo struct {
	outer_                          *RewriteAsyncFunctions
	wrapperFunctionName             string          `gson:"wrapperFunctionName"`
	wrapperFunctionReturnType       AstFactory_Type `gson:"wrapperFunctionReturnType"`
	firstInstanceOfSuperDotProperty pkg_rhino.Node  `gson:"firstInstanceOfSuperDotProperty"`
	idHash_                         atomic.Int32
}

func (this_ *RewriteAsyncFunctions_SuperPropertyWrapperInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteAsyncFunctions_SuperPropertyWrapperInfo) IsJavaNil_() bool { return this_ == nil }
func (*RewriteAsyncFunctions_SuperPropertyWrapperInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteAsyncFunctions$SuperPropertyWrapperInfo", "SuperPropertyWrapperInfo"
}
func (this_ *RewriteAsyncFunctions_SuperPropertyWrapperInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteAsyncFunctions_SuperPropertyWrapperInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteAsyncFunctions.SuperPropertyWrapperInfo", "com.google.javascript.jscomp.RewriteAsyncFunctions$SuperPropertyWrapperInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteAsyncFunctions_SuperPropertyWrapperInfo](v_)
	})
}
func (this_ *RewriteAsyncFunctions_SuperPropertyWrapperInfo) SetOuter_(o *RewriteAsyncFunctions) {
	this_.outer_ = o
}
func (this_ *RewriteAsyncFunctions_SuperPropertyWrapperInfo) GetOuter_() *RewriteAsyncFunctions {
	return this_.outer_
}

func (this_ *RewriteAsyncFunctions_SuperPropertyWrapperInfo) initFields_RewriteAsyncFunctions_SuperPropertyWrapperInfo() {
	_ = this_
	this_.wrapperFunctionName = j2gort.NullString
}

// Java constructor: SuperPropertyWrapperInfo(com.google.javascript.rhino.Node firstSuperDotPropertyNode, java.lang.String wrapperFunctionName, com.google.javascript.jscomp.AstFactory.Type wrapperFunctionReturnType)
func NewRewriteAsyncFunctions_SuperPropertyWrapperInfo(outer_ *RewriteAsyncFunctions, firstSuperDotPropertyNode pkg_rhino.Node, wrapperFunctionName string, wrapperFunctionReturnType AstFactory_Type) *RewriteAsyncFunctions_SuperPropertyWrapperInfo {
	this_ := &RewriteAsyncFunctions_SuperPropertyWrapperInfo{}
	this_.outer_ = outer_
	InitRewriteAsyncFunctions_SuperPropertyWrapperInfo(this_, firstSuperDotPropertyNode, wrapperFunctionName, wrapperFunctionReturnType)
	return this_
}

func InitRewriteAsyncFunctions_SuperPropertyWrapperInfo(this_ *RewriteAsyncFunctions_SuperPropertyWrapperInfo, firstSuperDotPropertyNode pkg_rhino.Node, wrapperFunctionName string, wrapperFunctionReturnType AstFactory_Type) {
	_ = firstSuperDotPropertyNode
	_ = wrapperFunctionName
	_ = wrapperFunctionReturnType
	this_.initFields_RewriteAsyncFunctions_SuperPropertyWrapperInfo()
	this_.firstInstanceOfSuperDotProperty = firstSuperDotPropertyNode
	this_.wrapperFunctionName = wrapperFunctionName
	this_.wrapperFunctionReturnType = wrapperFunctionReturnType
}

// Java method: SuperPropertyWrapperInfo.getPropertyType()
func (this_ *RewriteAsyncFunctions_SuperPropertyWrapperInfo) priv_SuperPropertyWrapperInfo_getPropertyType_4958() *pkg_colors.Color {
	return this_.firstInstanceOfSuperDotProperty.GetColor()
}

// Java method: SuperPropertyWrapperInfo.createWrapperFunctionNameNode()
func (this_ *RewriteAsyncFunctions_SuperPropertyWrapperInfo) priv_SuperPropertyWrapperInfo_createWrapperFunctionNameNode_4959() pkg_rhino.Node {
	return this_.outer_.astFactory.CreateName_2_string_AstFactory_Type(this_.wrapperFunctionName, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()))
}

// Java method: SuperPropertyWrapperInfo.createWrapperFunctionCallNode()
func (this_ *RewriteAsyncFunctions_SuperPropertyWrapperInfo) priv_SuperPropertyWrapperInfo_createWrapperFunctionCallNode_4960() pkg_rhino.Node {
	return this_.outer_.astFactory.CreateCall(this_.priv_SuperPropertyWrapperInfo_createWrapperFunctionNameNode_4959(), this_.wrapperFunctionReturnType, []pkg_rhino.Node{})
}

// Java class: com.google.javascript.jscomp.RewriteAsyncFunctions.SuperPropertyWrappers
type RewriteAsyncFunctions_SuperPropertyWrappers struct {
	outer_                *RewriteAsyncFunctions
	propertyNameToTypeMap j2gort.OrderedMap[string, *RewriteAsyncFunctions_SuperPropertyWrapperInfo] `gson:"propertyNameToTypeMap"`
	idHash_               atomic.Int32
}

func (this_ *RewriteAsyncFunctions_SuperPropertyWrappers) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteAsyncFunctions_SuperPropertyWrappers) IsJavaNil_() bool { return this_ == nil }
func (*RewriteAsyncFunctions_SuperPropertyWrappers) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteAsyncFunctions$SuperPropertyWrappers", "SuperPropertyWrappers"
}
func (this_ *RewriteAsyncFunctions_SuperPropertyWrappers) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteAsyncFunctions_SuperPropertyWrappers) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteAsyncFunctions.SuperPropertyWrappers", "com.google.javascript.jscomp.RewriteAsyncFunctions$SuperPropertyWrappers", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteAsyncFunctions_SuperPropertyWrappers](v_)
	})
}
func (this_ *RewriteAsyncFunctions_SuperPropertyWrappers) SetOuter_(o *RewriteAsyncFunctions) {
	this_.outer_ = o
}
func (this_ *RewriteAsyncFunctions_SuperPropertyWrappers) GetOuter_() *RewriteAsyncFunctions {
	return this_.outer_
}

func (this_ *RewriteAsyncFunctions_SuperPropertyWrappers) initFields_RewriteAsyncFunctions_SuperPropertyWrappers() {
	_ = this_
	this_.propertyNameToTypeMap = j2gort.NewOrderedMap[string, *RewriteAsyncFunctions_SuperPropertyWrapperInfo]()
}

// Java constructor: SuperPropertyWrappers()
func NewRewriteAsyncFunctions_SuperPropertyWrappers(outer_ *RewriteAsyncFunctions) *RewriteAsyncFunctions_SuperPropertyWrappers {
	this_ := &RewriteAsyncFunctions_SuperPropertyWrappers{}
	this_.outer_ = outer_
	InitRewriteAsyncFunctions_SuperPropertyWrappers(this_)
	return this_
}

func InitRewriteAsyncFunctions_SuperPropertyWrappers(this_ *RewriteAsyncFunctions_SuperPropertyWrappers) {
	this_.initFields_RewriteAsyncFunctions_SuperPropertyWrappers()
}

// Java method: SuperPropertyWrappers.getOrCreateSuperPropertyWrapperInfo(com.google.javascript.rhino.Node superDotPropertyNode)
func (this_ *RewriteAsyncFunctions_SuperPropertyWrappers) priv_SuperPropertyWrappers_getOrCreateSuperPropertyWrapperInfo_4961(superDotPropertyNode pkg_rhino.Node) *RewriteAsyncFunctions_SuperPropertyWrapperInfo {
	_ = superDotPropertyNode
	if !(superDotPropertyNode.IsGetProp()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(superDotPropertyNode))
	}
	var superNode pkg_rhino.Node = superDotPropertyNode.GetFirstChild()
	_ = superNode
	if !(superNode.IsSuper()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(superNode))
	}
	var propertyName string = superDotPropertyNode.GetString()
	_ = propertyName
	var propertyType *pkg_colors.Color = superDotPropertyNode.GetColor()
	_ = propertyType
	var superPropertyWrapperInfo *RewriteAsyncFunctions_SuperPropertyWrapperInfo = nil
	_ = superPropertyWrapperInfo
	if this_.propertyNameToTypeMap.ContainsKey(propertyName) {
		superPropertyWrapperInfo = this_.propertyNameToTypeMap.Get(propertyName)
		var existingColor *pkg_colors.Color = superPropertyWrapperInfo.priv_SuperPropertyWrapperInfo_getPropertyType_4958()
		_ = existingColor
		if !(j2gort.ObjEqCustomPtr(existingColor, propertyType)) {
			j2gort.CheckStateFmt(false, "Previous reference type: %s differs from current reference type: %s", j2gort.PtrToAny(existingColor), j2gort.PtrToAny(propertyType))
		}
	} else {
		superPropertyWrapperInfo = this_.priv_SuperPropertyWrappers_createNewInfo_4962(superDotPropertyNode)
		this_.propertyNameToTypeMap.Put(propertyName, superPropertyWrapperInfo)
	}
	return superPropertyWrapperInfo
}

// Java method: SuperPropertyWrappers.createNewInfo(com.google.javascript.rhino.Node firstSuperDotPropertyNode)
func (this_ *RewriteAsyncFunctions_SuperPropertyWrappers) priv_SuperPropertyWrappers_createNewInfo_4962(firstSuperDotPropertyNode pkg_rhino.Node) *RewriteAsyncFunctions_SuperPropertyWrapperInfo {
	_ = firstSuperDotPropertyNode
	if !(firstSuperDotPropertyNode.IsGetProp()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(firstSuperDotPropertyNode))
	}
	var propertyName string = firstSuperDotPropertyNode.GetString()
	_ = propertyName
	var wrapperFunctionName string = ("$jscomp$async$super$get$" + j2gort.StrValString(this_.outer_.compiler.GetUniqueIdSupplier().GetUniqueId(this_.outer_.compiler.GetInput(NodeUtil_GetInputId(firstSuperDotPropertyNode)))) + "$" + j2gort.StrValString(propertyName))
	_ = wrapperFunctionName
	return NewRewriteAsyncFunctions_SuperPropertyWrapperInfo(this_.outer_, firstSuperDotPropertyNode, wrapperFunctionName, AstFactory_S_Type_1_Node(firstSuperDotPropertyNode))
}

// Java method: SuperPropertyWrappers.asCollection()
func (this_ *RewriteAsyncFunctions_SuperPropertyWrappers) priv_SuperPropertyWrappers_asCollection_4963() j2gort.Ext_java_util_Collection {
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.propertyNameToTypeMap.LiveValues())
}

// Java polymorphic class interface: com.google.javascript.jscomp.RewriteAsyncFunctions.LexicalContext
type RewriteAsyncFunctions_LexicalContext interface {
	Implements_jscomp_RewriteAsyncFunctions_LexicalContext()
	IsJavaNil_() bool
	GetRewriteAsyncFunctions_LexicalContextBase() *RewriteAsyncFunctions_LexicalContextBase
	GetContextRootNode() pkg_rhino.Node
	GetContextForNode(*NodeTraversal, pkg_rhino.Node) RewriteAsyncFunctions_LexicalContext
	Visit(*NodeTraversal, pkg_rhino.Node)
}

func ToIface_RewriteAsyncFunctions_LexicalContext[T interface {
	RewriteAsyncFunctions_LexicalContext
	comparable
}](v T) RewriteAsyncFunctions_LexicalContext {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_RewriteAsyncFunctions_LexicalContext[T RewriteAsyncFunctions_LexicalContext](v T) RewriteAsyncFunctions_LexicalContext {
	var vi RewriteAsyncFunctions_LexicalContext = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.RewriteAsyncFunctions.LexicalContext
type RewriteAsyncFunctions_LexicalContextBase struct {
	uniqueId        string         `gson:"uniqueId"`
	contextRootNode pkg_rhino.Node `gson:"contextRootNode"`
	idHash_         atomic.Int32
}

func (this_ *RewriteAsyncFunctions_LexicalContextBase) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteAsyncFunctions_LexicalContextBase) IsJavaNil_() bool { return this_ == nil }
func (*RewriteAsyncFunctions_LexicalContextBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteAsyncFunctions$LexicalContext", "LexicalContext"
}
func (this_ *RewriteAsyncFunctions_LexicalContextBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteAsyncFunctions.LexicalContext", "com.google.javascript.jscomp.RewriteAsyncFunctions$LexicalContext", func(v_ any) bool {
		return j2gort.InstanceOf[RewriteAsyncFunctions_LexicalContext](v_)
	})
}
func (b *RewriteAsyncFunctions_LexicalContextBase) GetRewriteAsyncFunctions_LexicalContextBase() *RewriteAsyncFunctions_LexicalContextBase {
	return b
}
func (*RewriteAsyncFunctions_LexicalContextBase) Implements_jscomp_RewriteAsyncFunctions_LexicalContext() {
}

func (this_ *RewriteAsyncFunctions_LexicalContextBase) initFields_RewriteAsyncFunctions_LexicalContext() {
	_ = this_
	this_.uniqueId = j2gort.NullString
}

func InitRewriteAsyncFunctions_LexicalContext(this_ *RewriteAsyncFunctions_LexicalContextBase, contextRootNode pkg_rhino.Node, uniqueId string) {
	_ = contextRootNode
	_ = uniqueId
	this_.initFields_RewriteAsyncFunctions_LexicalContext()
	this_.contextRootNode = j2gort.CheckNotNullNode(contextRootNode)
	this_.uniqueId = uniqueId
}

// Java method: LexicalContext.getContextRootNode()
func (this_ *RewriteAsyncFunctions_LexicalContextBase) GetContextRootNode() pkg_rhino.Node {
	return this_.contextRootNode
}

// Java method: LexicalContext.getContextForNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *RewriteAsyncFunctions_LexicalContextBase) GetContextForNode(t *NodeTraversal, n pkg_rhino.Node) RewriteAsyncFunctions_LexicalContext {
	_ = t
	_ = n
	panic("abstract or native method invoked")
}

// Java method: LexicalContext.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *RewriteAsyncFunctions_LexicalContextBase) Visit(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	panic("abstract or native method invoked")
}

// Java class: com.google.javascript.jscomp.RewriteAsyncFunctions.RootContext
type RewriteAsyncFunctions_RootContext struct {
	RewriteAsyncFunctions_LexicalContextBase
	outer_ *RewriteAsyncFunctions
}

func (this_ *RewriteAsyncFunctions_RootContext) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteAsyncFunctions_RootContext) IsJavaNil_() bool { return this_ == nil }
func (*RewriteAsyncFunctions_RootContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteAsyncFunctions$RootContext", "RootContext"
}
func (this_ *RewriteAsyncFunctions_RootContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RewriteAsyncFunctions_LexicalContextBase.IdentityHash_()
}
func (*RewriteAsyncFunctions_RootContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteAsyncFunctions.RootContext", "com.google.javascript.jscomp.RewriteAsyncFunctions$RootContext", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteAsyncFunctions_RootContext](v_)
	})
}
func (this_ *RewriteAsyncFunctions_RootContext) SetOuter_(o *RewriteAsyncFunctions) { this_.outer_ = o }
func (this_ *RewriteAsyncFunctions_RootContext) GetOuter_() *RewriteAsyncFunctions {
	return this_.outer_
}
func (*RewriteAsyncFunctions_RootContext) Implements_jscomp_RewriteAsyncFunctions_LexicalContext() {}

func (this_ *RewriteAsyncFunctions_RootContext) initFields_RewriteAsyncFunctions_RootContext() {
	_ = this_
}

// Java constructor: RootContext(com.google.javascript.rhino.Node contextRootNode, java.lang.String uniqueId)
func NewRewriteAsyncFunctions_RootContext(outer_ *RewriteAsyncFunctions, contextRootNode pkg_rhino.Node, uniqueId string) *RewriteAsyncFunctions_RootContext {
	this_ := &RewriteAsyncFunctions_RootContext{}
	this_.outer_ = outer_
	InitRewriteAsyncFunctions_RootContext(this_, contextRootNode, uniqueId)
	return this_
}

func InitRewriteAsyncFunctions_RootContext(this_ *RewriteAsyncFunctions_RootContext, contextRootNode pkg_rhino.Node, uniqueId string) {
	_ = contextRootNode
	_ = uniqueId
	InitRewriteAsyncFunctions_LexicalContext(&this_.RewriteAsyncFunctions_LexicalContextBase, contextRootNode, uniqueId)
	this_.initFields_RewriteAsyncFunctions_RootContext()
}

// Java method: RootContext.getContextForNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *RewriteAsyncFunctions_RootContext) GetContextForNode(t *NodeTraversal, n pkg_rhino.Node) RewriteAsyncFunctions_LexicalContext {
	_ = t
	_ = n
	if n.IsFunction() {
		var uniqueId string = this_.outer_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		_ = uniqueId
		return NewRewriteAsyncFunctions_FunctionContext_2(this_.outer_, n, uniqueId)
	} else {
		return this_
	}
}

// Java method: RootContext.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *RewriteAsyncFunctions_RootContext) Visit(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
}

// Java class: com.google.javascript.jscomp.RewriteAsyncFunctions.ParameterListContext
type RewriteAsyncFunctions_ParameterListContext struct {
	RewriteAsyncFunctions_LexicalContextBase
	outer_          *RewriteAsyncFunctions
	functionContext *RewriteAsyncFunctions_FunctionContext `gson:"functionContext"`
}

func (this_ *RewriteAsyncFunctions_ParameterListContext) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteAsyncFunctions_ParameterListContext) IsJavaNil_() bool { return this_ == nil }
func (*RewriteAsyncFunctions_ParameterListContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteAsyncFunctions$ParameterListContext", "ParameterListContext"
}
func (this_ *RewriteAsyncFunctions_ParameterListContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RewriteAsyncFunctions_LexicalContextBase.IdentityHash_()
}
func (*RewriteAsyncFunctions_ParameterListContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteAsyncFunctions.ParameterListContext", "com.google.javascript.jscomp.RewriteAsyncFunctions$ParameterListContext", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteAsyncFunctions_ParameterListContext](v_)
	})
}
func (this_ *RewriteAsyncFunctions_ParameterListContext) SetOuter_(o *RewriteAsyncFunctions) {
	this_.outer_ = o
}
func (this_ *RewriteAsyncFunctions_ParameterListContext) GetOuter_() *RewriteAsyncFunctions {
	return this_.outer_
}
func (*RewriteAsyncFunctions_ParameterListContext) Implements_jscomp_RewriteAsyncFunctions_LexicalContext() {
}

func (this_ *RewriteAsyncFunctions_ParameterListContext) initFields_RewriteAsyncFunctions_ParameterListContext() {
	_ = this_
}

// Java constructor: ParameterListContext(com.google.javascript.jscomp.RewriteAsyncFunctions.FunctionContext functionContext, com.google.javascript.rhino.Node contextRootNode, java.lang.String uniqueId)
func NewRewriteAsyncFunctions_ParameterListContext(outer_ *RewriteAsyncFunctions, functionContext *RewriteAsyncFunctions_FunctionContext, contextRootNode pkg_rhino.Node, uniqueId string) *RewriteAsyncFunctions_ParameterListContext {
	this_ := &RewriteAsyncFunctions_ParameterListContext{}
	this_.outer_ = outer_
	InitRewriteAsyncFunctions_ParameterListContext(this_, functionContext, contextRootNode, uniqueId)
	return this_
}

func InitRewriteAsyncFunctions_ParameterListContext(this_ *RewriteAsyncFunctions_ParameterListContext, functionContext *RewriteAsyncFunctions_FunctionContext, contextRootNode pkg_rhino.Node, uniqueId string) {
	_ = functionContext
	_ = contextRootNode
	_ = uniqueId
	InitRewriteAsyncFunctions_LexicalContext(&this_.RewriteAsyncFunctions_LexicalContextBase, contextRootNode, uniqueId)
	this_.initFields_RewriteAsyncFunctions_ParameterListContext()
	this_.functionContext = j2gort.CheckNotNullPtr(functionContext)
}

// Java method: ParameterListContext.getContextForNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *RewriteAsyncFunctions_ParameterListContext) GetContextForNode(t *NodeTraversal, n pkg_rhino.Node) RewriteAsyncFunctions_LexicalContext {
	_ = t
	_ = n
	if n.IsFunction() {
		var uniqueId string = this_.outer_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		_ = uniqueId
		return NewRewriteAsyncFunctions_FunctionContext_3(this_.outer_, this_.priv_ParameterListContext_getAsyncContextForParameterList_4967(), n, uniqueId)
	} else {
		return this_
	}
}

// Java method: ParameterListContext.getAsyncContextForParameterList()
func (this_ *RewriteAsyncFunctions_ParameterListContext) priv_ParameterListContext_getAsyncContextForParameterList_4967() *RewriteAsyncFunctions_FunctionContext {
	var asyncContext *RewriteAsyncFunctions_FunctionContext = this_.functionContext.asyncThisAndArgumentsContext
	_ = asyncContext
	if asyncContext == this_.functionContext {
		return nil
	} else {
		return asyncContext
	}
}

// Java method: ParameterListContext.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *RewriteAsyncFunctions_ParameterListContext) Visit(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if this_.priv_ParameterListContext_getAsyncContextForParameterList_4967() != nil {
		this_.functionContext.Visit(t, n)
	}
}

// Java class: com.google.javascript.jscomp.RewriteAsyncFunctions.FunctionContext
type RewriteAsyncFunctions_FunctionContext struct {
	RewriteAsyncFunctions_LexicalContextBase
	outer_                        *RewriteAsyncFunctions
	asyncThisAndArgumentsContext  *RewriteAsyncFunctions_FunctionContext       `gson:"asyncThisAndArgumentsContext"`
	superPropertyWrappers         *RewriteAsyncFunctions_SuperPropertyWrappers `gson:"superPropertyWrappers"`
	typeOfThis                    AstFactory_Type                              `gson:"typeOfThis"`
	typeOfNewTarget               AstFactory_Type                              `gson:"typeOfNewTarget"`
	mustAddAsyncThisVariable      bool                                         `gson:"mustAddAsyncThisVariable"`
	mustAddAsyncArgumentsVariable bool                                         `gson:"mustAddAsyncArgumentsVariable"`
	mustAddAsyncNewTargetVariable bool                                         `gson:"mustAddAsyncNewTargetVariable"`
}

func (this_ *RewriteAsyncFunctions_FunctionContext) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteAsyncFunctions_FunctionContext) IsJavaNil_() bool { return this_ == nil }
func (*RewriteAsyncFunctions_FunctionContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteAsyncFunctions$FunctionContext", "FunctionContext"
}
func (this_ *RewriteAsyncFunctions_FunctionContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RewriteAsyncFunctions_LexicalContextBase.IdentityHash_()
}
func (*RewriteAsyncFunctions_FunctionContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteAsyncFunctions.FunctionContext", "com.google.javascript.jscomp.RewriteAsyncFunctions$FunctionContext", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteAsyncFunctions_FunctionContext](v_)
	})
}
func (this_ *RewriteAsyncFunctions_FunctionContext) SetOuter_(o *RewriteAsyncFunctions) {
	this_.outer_ = o
}
func (this_ *RewriteAsyncFunctions_FunctionContext) GetOuter_() *RewriteAsyncFunctions {
	return this_.outer_
}
func (*RewriteAsyncFunctions_FunctionContext) Implements_jscomp_RewriteAsyncFunctions_LexicalContext() {
}

func (this_ *RewriteAsyncFunctions_FunctionContext) initFields_RewriteAsyncFunctions_FunctionContext() {
	_ = this_
	this_.superPropertyWrappers = NewRewriteAsyncFunctions_SuperPropertyWrappers(this_.outer_)

}

// Java constructor: FunctionContext(com.google.javascript.rhino.Node contextRootNode, java.lang.String uniqueId)
func NewRewriteAsyncFunctions_FunctionContext_2(outer_ *RewriteAsyncFunctions, contextRootNode pkg_rhino.Node, uniqueId string) *RewriteAsyncFunctions_FunctionContext {
	this_ := &RewriteAsyncFunctions_FunctionContext{}
	this_.outer_ = outer_
	InitRewriteAsyncFunctions_FunctionContext_2(this_, contextRootNode, uniqueId)
	return this_
}

func InitRewriteAsyncFunctions_FunctionContext_2(this_ *RewriteAsyncFunctions_FunctionContext, contextRootNode pkg_rhino.Node, uniqueId string) {
	_ = contextRootNode
	_ = uniqueId
	InitRewriteAsyncFunctions_LexicalContext(&this_.RewriteAsyncFunctions_LexicalContextBase, contextRootNode, uniqueId)
	this_.initFields_RewriteAsyncFunctions_FunctionContext()
	if contextRootNode.IsAsyncFunction() {
		this_.asyncThisAndArgumentsContext = this_
	} else {
		this_.asyncThisAndArgumentsContext = nil
	}
}

// Java constructor: FunctionContext(com.google.javascript.jscomp.RewriteAsyncFunctions.@org.jspecify.annotations.Nullable FunctionContext outerAsyncContext, com.google.javascript.rhino.Node contextRootNode, java.lang.String uniqueId)
func NewRewriteAsyncFunctions_FunctionContext_3(outer_ *RewriteAsyncFunctions, outerAsyncContext *RewriteAsyncFunctions_FunctionContext, contextRootNode pkg_rhino.Node, uniqueId string) *RewriteAsyncFunctions_FunctionContext {
	this_ := &RewriteAsyncFunctions_FunctionContext{}
	this_.outer_ = outer_
	InitRewriteAsyncFunctions_FunctionContext_3(this_, outerAsyncContext, contextRootNode, uniqueId)
	return this_
}

func InitRewriteAsyncFunctions_FunctionContext_3(this_ *RewriteAsyncFunctions_FunctionContext, outerAsyncContext *RewriteAsyncFunctions_FunctionContext, contextRootNode pkg_rhino.Node, uniqueId string) {
	_ = outerAsyncContext
	_ = contextRootNode
	_ = uniqueId
	InitRewriteAsyncFunctions_LexicalContext(&this_.RewriteAsyncFunctions_LexicalContextBase, contextRootNode, uniqueId)
	this_.initFields_RewriteAsyncFunctions_FunctionContext()
	if !(contextRootNode.IsFunction()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(contextRootNode))
	}
	if contextRootNode.IsAsyncFunction() {
		if contextRootNode.IsArrowFunction() {
			if outerAsyncContext == nil {
				this_.asyncThisAndArgumentsContext = this_
			} else {
				this_.asyncThisAndArgumentsContext = outerAsyncContext
			}
		} else {
			this_.asyncThisAndArgumentsContext = this_
		}
	} else {
		if contextRootNode.IsArrowFunction() {
			this_.asyncThisAndArgumentsContext = outerAsyncContext
		} else {
			this_.asyncThisAndArgumentsContext = nil
		}
	}
}

// Java method: FunctionContext.getContextForNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *RewriteAsyncFunctions_FunctionContext) GetContextForNode(t *NodeTraversal, n pkg_rhino.Node) RewriteAsyncFunctions_LexicalContext {
	_ = t
	_ = n
	if n == this_.GetRewriteAsyncFunctions_LexicalContextBase().contextRootNode {
		return this_
	} else if n.IsFunction() {
		var uniqueId string = this_.outer_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		_ = uniqueId
		return NewRewriteAsyncFunctions_FunctionContext_3(this_.outer_, this_.asyncThisAndArgumentsContext, n, uniqueId)
	} else if n.IsParamList() {
		var uniqueId_2 string = this_.outer_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		_ = uniqueId_2
		return NewRewriteAsyncFunctions_ParameterListContext(this_.outer_, this_, n, uniqueId_2)
	} else {
		return this_
	}
}

// Java method: FunctionContext.recordAsyncThisReplacementWasDone(com.google.javascript.jscomp.AstFactory.Type typeOfThis)
func (this_ *RewriteAsyncFunctions_FunctionContext) priv_FunctionContext_recordAsyncThisReplacementWasDone_4968(typeOfThis AstFactory_Type) {
	_ = typeOfThis
	{
		r_3 := this_.asyncThisAndArgumentsContext
		r_3.mustAddAsyncThisVariable = true
	}
	{
		r_4 := this_.asyncThisAndArgumentsContext
		r_4.typeOfThis = typeOfThis
	}
}

// Java method: FunctionContext.getOrCreateSuperPropertyWrapperInfo(com.google.javascript.rhino.Node superDotPropertyNode)
func (this_ *RewriteAsyncFunctions_FunctionContext) priv_FunctionContext_getOrCreateSuperPropertyWrapperInfo_4969(superDotPropertyNode pkg_rhino.Node) *RewriteAsyncFunctions_SuperPropertyWrapperInfo {
	_ = superDotPropertyNode
	return this_.asyncThisAndArgumentsContext.superPropertyWrappers.priv_SuperPropertyWrappers_getOrCreateSuperPropertyWrapperInfo_4961(superDotPropertyNode)
}

// Java method: FunctionContext.recordAsyncArgumentsReplacementWasDone()
func (this_ *RewriteAsyncFunctions_FunctionContext) priv_FunctionContext_recordAsyncArgumentsReplacementWasDone_4970() {
	{
		r_5 := this_.asyncThisAndArgumentsContext
		r_5.mustAddAsyncArgumentsVariable = true
	}
}

// Java method: FunctionContext.recordAsyncNewTargetReplacementWasDone(com.google.javascript.jscomp.AstFactory.Type typeOfNewTarget)
func (this_ *RewriteAsyncFunctions_FunctionContext) priv_FunctionContext_recordAsyncNewTargetReplacementWasDone_4971(typeOfNewTarget AstFactory_Type) {
	_ = typeOfNewTarget
	{
		r_6 := this_.asyncThisAndArgumentsContext
		r_6.mustAddAsyncNewTargetVariable = true
	}
	{
		r_7 := this_.asyncThisAndArgumentsContext
		r_7.typeOfNewTarget = typeOfNewTarget
	}
}

// Java method: FunctionContext.createThisVariableReference(com.google.javascript.jscomp.AstFactory.Type typeOfThis)
func (this_ *RewriteAsyncFunctions_FunctionContext) priv_FunctionContext_createThisVariableReference_4972(typeOfThis AstFactory_Type) pkg_rhino.Node {
	_ = typeOfThis
	this_.priv_FunctionContext_recordAsyncThisReplacementWasDone_4968(typeOfThis)
	return this_.outer_.astFactory.CreateName_2_string_AstFactory_Type(("$jscomp$async$this$" + j2gort.StrValString(this_.asyncThisAndArgumentsContext.GetRewriteAsyncFunctions_LexicalContextBase().uniqueId)), typeOfThis)
}

// Java method: FunctionContext.createNewTargetVariableReference(com.google.javascript.jscomp.AstFactory.Type typeOfNewTarget)
func (this_ *RewriteAsyncFunctions_FunctionContext) priv_FunctionContext_createNewTargetVariableReference_4973(typeOfNewTarget AstFactory_Type) pkg_rhino.Node {
	_ = typeOfNewTarget
	this_.priv_FunctionContext_recordAsyncNewTargetReplacementWasDone_4971(typeOfNewTarget)
	return this_.outer_.astFactory.CreateName_2_string_AstFactory_Type(("$jscomp$async$new$target$" + j2gort.StrValString(this_.asyncThisAndArgumentsContext.GetRewriteAsyncFunctions_LexicalContextBase().uniqueId)), typeOfNewTarget)
}

// Java method: FunctionContext.createWrapperArrowFunction(com.google.javascript.jscomp.RewriteAsyncFunctions.SuperPropertyWrapperInfo wrapperInfo)
func (this_ *RewriteAsyncFunctions_FunctionContext) priv_FunctionContext_createWrapperArrowFunction_4974(wrapperInfo *RewriteAsyncFunctions_SuperPropertyWrapperInfo) pkg_rhino.Node {
	_ = wrapperInfo
	var superDotProperty pkg_rhino.Node = wrapperInfo.firstInstanceOfSuperDotProperty.CloneTree_0()
	_ = superDotProperty
	return this_.outer_.astFactory.CreateZeroArgArrowFunctionForExpression(this_.outer_.astFactory.CreateBlock([]pkg_rhino.Node{this_.outer_.astFactory.CreateReturn(superDotProperty)}))
}

// Java method: FunctionContext.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *RewriteAsyncFunctions_FunctionContext) Visit(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if (this_.GetRewriteAsyncFunctions_LexicalContextBase().contextRootNode == n) && this_.GetRewriteAsyncFunctions_LexicalContextBase().contextRootNode.IsAsyncFunction() {
		this_.outer_.priv_RewriteAsyncFunctions_convertAsyncFunction_4957(t, this_)
	} else if this_.asyncThisAndArgumentsContext != nil {
		switch n.TokenSwitchOrd_() {
		case int32(32) /* Token.NAME */ :
			if n.MatchesName_1_string("arguments") {
				n.SetString(("$jscomp$async$arguments$" + j2gort.StrValString(this_.asyncThisAndArgumentsContext.GetRewriteAsyncFunctions_LexicalContextBase().uniqueId)))
				if this_.outer_.compiler.GetLifeCycleStage().IsNormalized() {
					n.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
				}
				this_.asyncThisAndArgumentsContext.priv_FunctionContext_recordAsyncArgumentsReplacementWasDone_4970()
				this_.outer_.compiler.ReportChangeToChangeScope(this_.GetRewriteAsyncFunctions_LexicalContextBase().contextRootNode)
			}
		case int32(37) /* Token.THIS */ :
			n.ReplaceWith(this_.asyncThisAndArgumentsContext.priv_FunctionContext_createThisVariableReference_4972(AstFactory_S_Type_1_Node(n)))
			this_.outer_.compiler.ReportChangeToChangeScope(this_.GetRewriteAsyncFunctions_LexicalContextBase().contextRootNode)
		case int32(135) /* Token.NEW_TARGET */ :
			n.ReplaceWith(this_.asyncThisAndArgumentsContext.priv_FunctionContext_createNewTargetVariableReference_4973(AstFactory_S_Type_1_Node(n)))
			this_.outer_.compiler.ReportChangeToChangeScope(this_.GetRewriteAsyncFunctions_LexicalContextBase().contextRootNode)
		case int32(110) /* Token.SUPER */ :
			{
				var parent pkg_rhino.Node = n.GetParent()
				_ = parent
				if !(parent.IsGetProp()) {
					this_.outer_.compiler.Report_1(JSError_Make_3(parent, StaticGet_TranspilationUtil_CANNOT_CONVERT_YET(), []string{"super expression"}))
				} else if NodeUtil_IsLValue(parent) {
					this_.outer_.compiler.Report_1(JSError_Make_3(parent.GetParent(), StaticGet_TranspilationUtil_CANNOT_CONVERT_YET(), []string{"assignment to super property"}))
				}
				var superDotProperty pkg_rhino.Node = parent
				_ = superDotProperty
				var superPropertyWrapperInfo *RewriteAsyncFunctions_SuperPropertyWrapperInfo = this_.asyncThisAndArgumentsContext.priv_FunctionContext_getOrCreateSuperPropertyWrapperInfo_4969(superDotProperty)
				_ = superPropertyWrapperInfo
				var getPropReplacement pkg_rhino.Node = superPropertyWrapperInfo.priv_SuperPropertyWrapperInfo_createWrapperFunctionCallNode_4960()
				_ = getPropReplacement
				var grandparent pkg_rhino.Node = superDotProperty.GetParent()
				_ = grandparent
				if grandparent.IsCall() && (grandparent.GetFirstChild() == superDotProperty) {
					getPropReplacement = this_.outer_.astFactory.CreateGetPropWithUnknownType(getPropReplacement, "call")
					var thisAlias pkg_rhino.Node = this_.outer_.astFactory.CreateThisAliasReferenceForEs6Class(("$jscomp$async$this$" + j2gort.StrValString(this_.asyncThisAndArgumentsContext.GetRewriteAsyncFunctions_LexicalContextBase().uniqueId)), NodeUtil_GetEnclosingClass(this_.asyncThisAndArgumentsContext.GetContextRootNode())).Srcref(superDotProperty)
					_ = thisAlias
					thisAlias.InsertAfter(superDotProperty)
					this_.asyncThisAndArgumentsContext.priv_FunctionContext_recordAsyncThisReplacementWasDone_4968(AstFactory_S_Type_1_Node(thisAlias))
				} else if grandparent.IsTaggedTemplateLit() && (grandparent.GetFirstChild() == superDotProperty) {
					var thisAlias_2 pkg_rhino.Node = this_.outer_.astFactory.CreateThisAliasReferenceForEs6Class(("$jscomp$async$this$" + j2gort.StrValString(this_.asyncThisAndArgumentsContext.GetRewriteAsyncFunctions_LexicalContextBase().uniqueId)), NodeUtil_GetEnclosingClass(this_.asyncThisAndArgumentsContext.GetContextRootNode())).Srcref(superDotProperty)
					_ = thisAlias_2
					this_.asyncThisAndArgumentsContext.priv_FunctionContext_recordAsyncThisReplacementWasDone_4968(AstFactory_S_Type_1_Node(thisAlias_2))
					getPropReplacement = this_.outer_.astFactory.CreateCallWithUnknownType(this_.outer_.astFactory.CreateGetPropWithUnknownType(getPropReplacement, "bind"), []pkg_rhino.Node{thisAlias_2})
					grandparent.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
				}
				getPropReplacement.SrcrefTree(superDotProperty)
				superDotProperty.ReplaceWith(getPropReplacement)
				this_.outer_.compiler.ReportChangeToChangeScope(this_.GetRewriteAsyncFunctions_LexicalContextBase().contextRootNode)
			}
		case int32(115) /* Token.AWAIT */ :
			n.ReplaceWith(this_.outer_.astFactory.CreateYield(AstFactory_S_Type_1_Node(n), n.RemoveFirstChild()))
		default:
		}
	}
}
