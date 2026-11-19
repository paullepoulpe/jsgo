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
// Port of com/google/javascript/jscomp/Es6TemplateLiterals.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Es6TemplateLiterals
type Es6TemplateLiterals struct {
	astFactory *AstFactory `gson:"astFactory"`
	compiler AbstractCompiler `gson:"compiler"`
	namespace pkg_rhino.StaticScope `gson:"namespace"`
	createTemplateTagFirstArg RuntimeJsLibManager_JsLibField `gson:"createTemplateTagFirstArg"`
	createTemplateTagFirstArgWithRaw RuntimeJsLibManager_JsLibField `gson:"createTemplateTagFirstArgWithRaw"`
	idHash_ atomic.Int32
}

func (this_ *Es6TemplateLiterals) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6TemplateLiterals) IsJavaNil_() bool { return this_ == nil }
func (*Es6TemplateLiterals) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.Es6TemplateLiterals", "Es6TemplateLiterals" }
func (this_ *Es6TemplateLiterals) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6TemplateLiterals) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6TemplateLiterals", "com.google.javascript.jscomp.Es6TemplateLiterals", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6TemplateLiterals](v_)
	})
}

func StaticGet_Es6TemplateLiterals_TEMPLATELIT_VAR() string { return "$jscomp$templatelit$" }
func (this_ *Es6TemplateLiterals) initFields_Es6TemplateLiterals() {
	_ = this_
}

// Java constructor: Es6TemplateLiterals(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6TemplateLiterals(compiler AbstractCompiler) *Es6TemplateLiterals {
	this_ := &Es6TemplateLiterals{}
	InitEs6TemplateLiterals(this_, compiler)
	return this_
}

func InitEs6TemplateLiterals(this_ *Es6TemplateLiterals, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_Es6TemplateLiterals()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.namespace = this_.compiler.GetTranspilationNamespace()
	var runtimeJsLibManager *RuntimeJsLibManager = compiler.GetRuntimeJsLibManager()
	_ = runtimeJsLibManager
	this_.createTemplateTagFirstArg = runtimeJsLibManager.GetJsLibField("$jscomp.createTemplateTagFirstArg")
	this_.createTemplateTagFirstArgWithRaw = runtimeJsLibManager.GetJsLibField("$jscomp.createTemplateTagFirstArgWithRaw")
}

// Java method: Es6TemplateLiterals.visitTemplateLiteral(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *Es6TemplateLiterals) VisitTemplateLiteral(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var length int32 = n.GetChildCount()
	_ = length
	if ((length == int32(0))) {
		n.ReplaceWith(this_.astFactory.CreateString("\"\""))
	} else {
		var first pkg_rhino.Node = n.RemoveFirstChild()
		_ = first
		if !((first.IsTemplateLitString() && (first.GetCookedString() != j2gort.NullString))) { j2gort.CheckState(false) }
		var firstStr pkg_rhino.Node = this_.astFactory.CreateString(first.GetCookedString())
		_ = firstStr
		if ((length == int32(1))) {
			n.ReplaceWith(firstStr)
		} else {
			var add pkg_rhino.Node = this_.astFactory.CreateAdd(firstStr, n.RemoveFirstChild().RemoveFirstChild())
			_ = add
			{
				var i int32 = int32(2)
				_ = i
				for ; (i < length); i = i + 1 {
					var child pkg_rhino.Node = n.RemoveFirstChild()
					_ = child
					if (child.IsTemplateLitString()) {
						if !((child.GetCookedString() != j2gort.NullString)) { j2gort.CheckState(false) }
						if ((j2gort.CheckNotNullString(child.GetCookedString()) == "")) {
							continue
						} else if (((i == int32(2)) && (j2gort.CheckNotNullString(first.GetCookedString()) == ""))) {
							add = add.GetSecondChild().Detach()
						}
					}
					add = this_.astFactory.CreateAdd(add, func() pkg_rhino.Node { if child.IsTemplateLitString() { return this_.astFactory.CreateString(child.GetCookedString()) }; return child.RemoveFirstChild() }())
				}
			}
			n.ReplaceWith(add.SrcrefTreeIfMissing(n))
		}
	}
	t.ReportCodeChange_0()
}

// Java method: Es6TemplateLiterals.visitTaggedTemplateLiteral(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node insertBefore)
func (this_ *Es6TemplateLiterals) VisitTaggedTemplateLiteral(t *NodeTraversal, n pkg_rhino.Node, insertBefore pkg_rhino.Node) {
	_ = t
	_ = n
	_ = insertBefore
	var templateLit pkg_rhino.Node = n.GetLastChild()
	_ = templateLit
	var siteObject pkg_rhino.Node = this_.priv_Es6TemplateLiterals_createCookedStringArray_2779(templateLit)
	_ = siteObject
	var callTemplateTagArgCreator pkg_rhino.Node = pkg_rhino.Node(0)
	_ = callTemplateTagArgCreator
	if (Es6TemplateLiterals_CookedAndRawStringsSame(templateLit)) {
		callTemplateTagArgCreator = this_.astFactory.CreateCall(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.createTemplateTagFirstArg), AstFactory_S_Type_1_Node(siteObject), []pkg_rhino.Node{siteObject.CloneTree_0()})
	} else {
		var raw pkg_rhino.Node = this_.priv_Es6TemplateLiterals_createRawStringArray_2778(templateLit)
		_ = raw
		callTemplateTagArgCreator = this_.astFactory.CreateCall(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.createTemplateTagFirstArgWithRaw), AstFactory_S_Type_1_Node(siteObject), []pkg_rhino.Node{siteObject.CloneTree_0(), raw})
	}
	var jsDocInfoBuilder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
	_ = jsDocInfoBuilder
	jsDocInfoBuilder.RecordNoInline()
	var info *pkg_rhino.JSDocInfo = jsDocInfoBuilder.Build_0()
	_ = info
	var input *CompilerInput = t.GetInput()
	_ = input
	var uniqueId string = this_.compiler.GetUniqueIdSupplier().GetUniqueId(input)
	_ = uniqueId
	var tagFnFirstArgDeclaration pkg_rhino.Node = this_.astFactory.CreateSingleVarNameDeclaration_2(("$jscomp$templatelit$" + j2gort.StrValString(uniqueId)), callTemplateTagArgCreator).SetJSDocInfo(info).SrcrefTreeIfMissing(n)
	_ = tagFnFirstArgDeclaration
	tagFnFirstArgDeclaration.InsertBefore(insertBefore)
	t.ReportCodeChange_1(tagFnFirstArgDeclaration)
	var tagFnFirstArg pkg_rhino.Node = tagFnFirstArgDeclaration.GetFirstChild().CloneNode_0()
	_ = tagFnFirstArg
	var call pkg_rhino.Node = this_.astFactory.CreateCall(n.RemoveFirstChild(), AstFactory_S_Type_1_Node(n), []pkg_rhino.Node{tagFnFirstArg})
	_ = call
	{
		var child pkg_rhino.Node = templateLit.GetFirstChild()
		_ = child
		for ; (child != 0); child = child.GetNext() {
			if (!(child.IsTemplateLitString())) {
				call.AddChildToBack(child.RemoveFirstChild())
			}
		}
	}
	call.SrcrefTreeIfMissing(templateLit)
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), !(NodeUtil_IsNormalOrOptChainGet(call.GetFirstChild())))
	n.ReplaceWith(call)
	t.ReportCodeChange_0()
}

// Java method: Es6TemplateLiterals.createRawStringArray(com.google.javascript.rhino.Node n)
func (this_ *Es6TemplateLiterals) priv_Es6TemplateLiterals_createRawStringArray_2778(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var array pkg_rhino.Node = this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{})
	_ = array
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; (child != 0); child = child.GetNext() {
			if (child.IsTemplateLitString()) {
				array.AddChildToBack(this_.astFactory.CreateString(child.GetRawString()))
			}
		}
	}
	return array
}

// Java method: Es6TemplateLiterals.createCookedStringArray(com.google.javascript.rhino.Node n)
func (this_ *Es6TemplateLiterals) priv_Es6TemplateLiterals_createCookedStringArray_2779(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var templateArrayType *pkg_colors.Color
	_ = templateArrayType
	if this_.astFactory.IsAddingColors() {
		templateArrayType = this_.compiler.GetColorRegistry().Get(pkg_colors.StaticGet_StandardColors_I_TEMPLATE_ARRAY_ID())
	} else {
		templateArrayType = nil
	}
	var array pkg_rhino.Node = this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{}).SetColor(templateArrayType)
	_ = array
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; (child != 0); child = child.GetNext() {
			if (child.IsTemplateLitString()) {
				if ((child.GetCookedString() != j2gort.NullString)) {
					array.AddChildToBack(this_.astFactory.CreateString(child.GetCookedString()))
				} else {
					array.AddChildToBack(this_.astFactory.CreateVoid(this_.astFactory.CreateNumber(float64(int32(0)))))
				}
			}
		}
	}
	return array
}

// Java static method: Es6TemplateLiterals.cookedAndRawStringsSame(com.google.javascript.rhino.Node n)
func Es6TemplateLiterals_CookedAndRawStringsSame(n pkg_rhino.Node) bool {
	_ = n
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; (child != 0); child = child.GetNext() {
			if (!(child.IsTemplateLitString())) {
				continue
			}
			if (((child.GetCookedString() == j2gort.NullString) || !(func() bool { sr_1 := child.GetCookedString(); se_2 := child.GetRawString(); return j2gort.CheckNotNullString(sr_1) == se_2 }()))) {
				return false
			}
		}
	}
	return true
}

