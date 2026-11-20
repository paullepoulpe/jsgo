/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/Es6ForOfConverter.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Es6ForOfConverter
type Es6ForOfConverter struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler   AbstractCompiler      `gson:"compiler"`
	namer      *DefaultNameGenerator `gson:"namer"`
	astFactory *AstFactory           `gson:"astFactory"`
	namespace  pkg_rhino.StaticScope `gson:"namespace"`
}

func (this_ *Es6ForOfConverter) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6ForOfConverter) IsJavaNil_() bool { return this_ == nil }
func (*Es6ForOfConverter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6ForOfConverter", "Es6ForOfConverter"
}
func (this_ *Es6ForOfConverter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6ForOfConverter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6ForOfConverter", "com.google.javascript.jscomp.Es6ForOfConverter", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6ForOfConverter](v_)
	})
}
func (*Es6ForOfConverter) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*Es6ForOfConverter) Implements_jscomp_CompilerPass()                            {}
func (*Es6ForOfConverter) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_Es6ForOfConverter_transpiledFeatures *pkg_parser.FeatureSet

func StaticGet_Es6ForOfConverter_TranspiledFeatures() *pkg_parser.FeatureSet {
	clinit_Es6ForOfConverter()
	return static_Es6ForOfConverter_transpiledFeatures
}
func StaticSet_Es6ForOfConverter_TranspiledFeatures(v *pkg_parser.FeatureSet) {
	clinit_Es6ForOfConverter()
	static_Es6ForOfConverter_transpiledFeatures = v
}
func StaticGet_Es6ForOfConverter_ITER_BASE() string   { return "$jscomp$iter$" }
func StaticGet_Es6ForOfConverter_ITER_RESULT() string { return "$jscomp$key$" }

var clinitDone_Es6ForOfConverter uint32
var clinitState_Es6ForOfConverter j2gort.ClinitState

func clinit_Es6ForOfConverter() {
	if atomic.LoadUint32(&clinitDone_Es6ForOfConverter) != 0 {
		return
	}
	clinitSlow_Es6ForOfConverter()
}

func clinitSlow_Es6ForOfConverter() {
	if !j2gort.EnterClinit(&clinitState_Es6ForOfConverter, &clinitDone_Es6ForOfConverter) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es6ForOfConverter, &clinitDone_Es6ForOfConverter)
	static_Es6ForOfConverter_transpiledFeatures = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_FOR_OF())
}

func (this_ *Es6ForOfConverter) initFields_Es6ForOfConverter() {
	_ = this_
}

// Java constructor: Es6ForOfConverter(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6ForOfConverter(compiler AbstractCompiler) *Es6ForOfConverter {
	clinit_Es6ForOfConverter()
	this_ := &Es6ForOfConverter{}
	InitEs6ForOfConverter(this_, compiler)
	return this_
}

func InitEs6ForOfConverter(this_ *Es6ForOfConverter, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6ForOfConverter()
	this_.compiler = compiler
	this_.namer = NewDefaultNameGenerator()
	this_.astFactory = compiler.CreateAstFactory()
	this_.namespace = compiler.GetTranspilationNamespace()
}

// Java method: Es6ForOfConverter.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6ForOfConverter) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	TranspilationPasses_ProcessTranspile(this_.compiler, root, StaticGet_Es6ForOfConverter_TranspiledFeatures(), []NodeTraversal_Callback{this_})
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, StaticGet_Es6ForOfConverter_TranspiledFeatures())
}

// Java method: Es6ForOfConverter.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6ForOfConverter) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsForOf() {
		this_.priv_Es6ForOfConverter_visitForOf_2503(t, n)
	}
}

// Java method: Es6ForOfConverter.visitForOf(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node)
func (this_ *Es6ForOfConverter) priv_Es6ForOfConverter_visitForOf_2503(t *NodeTraversal, node pkg_rhino.Node) {
	_ = t
	_ = node
	var variable pkg_rhino.Node = node.RemoveFirstChild()
	_ = variable
	var iterable pkg_rhino.Node = node.RemoveFirstChild()
	_ = iterable
	var body pkg_rhino.Node = node.RemoveFirstChild()
	_ = body
	var varJSDocInfo *pkg_rhino.JSDocInfo = variable.GetJSDocInfo()
	_ = varJSDocInfo
	var iterName pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(("$jscomp$iter$" + j2gort.StrValString(this_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput()))), AstFactory_S_Type_1_ColorId(pkg_colors.StaticGet_StandardColors_ITERATOR_ID()))
	_ = iterName
	iterName.MakeNonIndexable()
	var getNext pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(this_.astFactory.CreateGetPropWithUnknownType(iterName.CloneTree_0(), "next"), []pkg_rhino.Node{})
	_ = getNext
	var iteratorResultName string = ("$jscomp$key$" + j2gort.StrValString(this_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())) + "$")
	_ = iteratorResultName
	if NodeUtil_IsNameDeclaration(variable) {
		iteratorResultName = (j2gort.StrValString(iteratorResultName) + j2gort.StrValString(variable.GetFirstChild().GetString()))
	} else if variable.IsName() {
		iteratorResultName = (j2gort.StrValString(iteratorResultName) + j2gort.StrValString(variable.GetString()))
	} else {
		iteratorResultName = (j2gort.StrValString(iteratorResultName) + j2gort.StrValString(this_.namer.GenerateNextName()))
	}
	var iterResult pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(iteratorResultName)
	_ = iterResult
	iterResult.MakeNonIndexable()
	var callMakeIterator pkg_rhino.Node = this_.astFactory.CreateJSCompMakeIteratorCall(iterable, this_.namespace).SrcrefTreeIfMissing(iterable)
	_ = callMakeIterator
	var initIter pkg_rhino.Node = pkg_rhino.IR_Var_2(iterName.CloneTree_0(), callMakeIterator).SrcrefTreeIfMissing(iterable)
	_ = initIter
	var initIterResult pkg_rhino.Node = pkg_rhino.IR_Var_2(iterResult.CloneTree_0(), getNext.CloneTree_0()).SrcrefTreeIfMissing(iterable)
	_ = initIterResult
	var cond pkg_rhino.Node = this_.astFactory.CreateNot(this_.astFactory.CreateGetProp(iterResult.CloneTree_0(), "done", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_BOOLEAN())))
	_ = cond
	var incr pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(iterResult.CloneTree_0(), getNext.CloneTree_0())
	_ = incr
	var declarationOrAssign pkg_rhino.Node = pkg_rhino.Node(0)
	_ = declarationOrAssign
	if !(NodeUtil_IsNameDeclaration(variable)) {
		declarationOrAssign = this_.astFactory.CreateAssign_2_Node_Node(variable.CloneTree_0().SetJSDocInfo(nil), this_.astFactory.CreateGetProp(iterResult.CloneTree_0(), "value", AstFactory_S_Type_1_Node(variable)))
		declarationOrAssign.SetJSDocInfo(varJSDocInfo)
		declarationOrAssign = pkg_rhino.IR_ExprResult(declarationOrAssign)
	} else {
		var type_ AstFactory_Type = AstFactory_S_Type_1_Node(variable.GetFirstChild())
		_ = type_
		var declarationType pkg_rhino.Token = variable.GetToken()
		_ = declarationType
		{
			var iifeCond_1 bool
			{

				re_1 := declarationType
				ae_2 := pkg_rhino.StaticGet_Token_VAR()
				if re_1 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_1 == ae_2
			}
			if !(!(iifeCond_1)) {
				j2gort.CheckStateStr(false, "var initializers must've gotten moved out of the loop during normalize")
			}
		}
		declarationOrAssign = this_.astFactory.CreateSingleNameDeclaration(declarationType, variable.GetFirstChild().GetString(), this_.astFactory.CreateGetProp(iterResult.CloneTree_0(), "value", type_))
		if variable.GetFirstChild().GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME()) {
			declarationOrAssign.GetFirstChild().PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
		}
		declarationOrAssign.SetJSDocInfo(varJSDocInfo)
	}
	var newBody pkg_rhino.Node = pkg_rhino.IR_Block_1_Node_2([]pkg_rhino.Node{declarationOrAssign, body}).Srcref(body)
	_ = newBody
	var empty pkg_rhino.Node = this_.astFactory.CreateEmpty()
	_ = empty
	var newFor pkg_rhino.Node = pkg_rhino.IR_ForNode(empty, cond, incr, newBody).SrcrefTreeIfMissing(node)
	_ = newFor
	var callIteratorClose pkg_rhino.Node = this_.astFactory.CreateJscompIteratorCloseCall(iterName.CloneTree_0(), iterResult.CloneTree_0(), this_.namespace).SrcrefTreeIfMissing(node)
	_ = callIteratorClose
	var finallyBlock pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{this_.astFactory.ExprResult(callIteratorClose)})
	_ = finallyBlock
	var labelNames []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = labelNames
	var insertionPoint pkg_rhino.Node = node
	_ = insertionPoint
	for (insertionPoint.GetParent() != 0) && insertionPoint.GetParent().IsLabel() {
		insertionPoint = insertionPoint.GetParent()
		labelNames = append(labelNames, insertionPoint.GetFirstChild().CloneNode_0())
	}
	var innerLoop pkg_rhino.Node = newFor
	_ = innerLoop
	for _, labelName := range labelNames {
		_ = labelName
		innerLoop = this_.astFactory.CreateLabel(labelName, innerLoop)
	}

	var tryBlock pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{innerLoop})
	_ = tryBlock
	var tryFinally pkg_rhino.Node = this_.astFactory.CreateTryFinally(tryBlock, finallyBlock).SrcrefTreeIfMissing(node)
	_ = tryFinally
	insertionPoint.ReplaceWith(tryFinally)
	initIter.InsertBefore(tryFinally)
	initIterResult.InsertAfter(initIter)
	this_.compiler.ReportChangeToEnclosingScope(tryFinally)
}
func (*Es6ForOfConverter) AlwaysShouldTraverse_()      {}
func (*Es6ForOfConverter) CanReleaseSyntacticScopes_() {}
