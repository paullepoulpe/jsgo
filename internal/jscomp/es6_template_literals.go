// Port of com/google/javascript/jscomp/Es6TemplateLiterals.java @ closure-compiler v20261021

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: Es6TemplateLiterals.TEMPLATELIT_VAR
const Es6TemplateLiterals_TEMPLATELIT_VAR = "$jscomp$templatelit$"

// Es6TemplateLiterals is a helper for transpiling ES6 template literals.
type Es6TemplateLiterals struct {
	astFactory                       *AstFactory
	compiler                         AbstractCompiler
	namespace                        rhino.StaticScope
	createTemplateTagFirstArg        RuntimeJsLibManager_JsLibField
	createTemplateTagFirstArgWithRaw RuntimeJsLibManager_JsLibField
}

// Java: Es6TemplateLiterals(AbstractCompiler)
func NewEs6TemplateLiterals(compiler AbstractCompiler) *Es6TemplateLiterals {
	runtimeJsLibManager := compiler.GetRuntimeJsLibManager()
	return &Es6TemplateLiterals{
		compiler:                         compiler,
		astFactory:                       compiler.CreateAstFactory(),
		namespace:                        compiler.GetTranspilationNamespace(),
		createTemplateTagFirstArg:        runtimeJsLibManager.GetJsLibField("$jscomp.createTemplateTagFirstArg"),
		createTemplateTagFirstArgWithRaw: runtimeJsLibManager.GetJsLibField("$jscomp.createTemplateTagFirstArgWithRaw"),
	}
}

// VisitTemplateLiteral converts `${a} b ${c} d ${e}` to `(a + " b " + c + " d " + e)`.
//
// Java: Es6TemplateLiterals.visitTemplateLiteral(NodeTraversal, Node)
func (e *Es6TemplateLiterals) VisitTemplateLiteral(t *NodeTraversal, n rhino.Node) {
	length := n.GetChildCount()
	if length == 0 {
		n.ReplaceWith(e.astFactory.CreateString("\"\""))
	} else {
		firstTls, ok := ast.AsTemplateLitString(n.RemoveFirstChild())
		jutil.CheckState(ok)
		firstCooked, hasFirstCooked := firstTls.Cooked()
		jutil.CheckState(hasFirstCooked)
		firstStr := e.astFactory.CreateString(firstCooked)
		if length == 1 {
			n.ReplaceWith(firstStr)
		} else {
			firstSub := ast.MustTemplateLitSub(n.RemoveFirstChild())
			add := e.astFactory.CreateAdd(firstStr, firstSub.Expr().Node().Detach())
			for i := 2; i < length; i++ {
				child := n.RemoveFirstChild()
				var right rhino.Node
				if tls, isStr := ast.AsTemplateLitString(child); isStr {
					cooked, hasCooked := tls.Cooked()
					jutil.CheckState(hasCooked)
					if cooked == "" {
						continue
					}
					if i == 2 && firstCooked == "" {
						// Translate `${hello} world` into `(hello + " world")` instead of `("" + hello + " world")`.
						add = add.GetSecondChild().Detach()
					}
					right = e.astFactory.CreateString(cooked)
				} else {
					sub := ast.MustTemplateLitSub(child)
					right = sub.Expr().Node().Detach()
				}
				add = e.astFactory.CreateAdd(add, right)
			}
			n.ReplaceWith(add.SrcrefTreeIfMissing(n))
		}
	}
	t.ReportCodeChange()
}

// VisitTaggedTemplateLiteral converts a tagged template into a call to the tag.
//
// Java: Es6TemplateLiterals.visitTaggedTemplateLiteral(NodeTraversal, Node, Node)
func (e *Es6TemplateLiterals) VisitTaggedTemplateLiteral(t *NodeTraversal, n rhino.Node, insertBefore rhino.Node) {
	tagged := ast.MustTaggedTemplateLit(n)
	templateLit := tagged.Template()
	siteObject := e.createCookedStringArray(templateLit)

	var callTemplateTagArgCreator rhino.Node
	if cookedAndRawStringsSame(templateLit) {
		callTemplateTagArgCreator = e.astFactory.CreateCall(
			e.astFactory.CreateQNameOfField(e.namespace, e.createTemplateTagFirstArg),
			AstFactory_Type_(siteObject),
			siteObject.CloneTree())
	} else {
		raw := e.createRawStringArray(templateLit)
		callTemplateTagArgCreator = e.astFactory.CreateCall(
			e.astFactory.CreateQNameOfField(e.namespace, e.createTemplateTagFirstArgWithRaw),
			AstFactory_Type_(siteObject),
			siteObject.CloneTree(),
			raw)
	}
	jsDocInfoBuilder := rhino.JSDocInfo_Builder_()
	jsDocInfoBuilder.RecordNoInline()
	info := jsDocInfoBuilder.Build()

	uniqueId := e.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
	tagFnFirstArgDeclaration := e.astFactory.
		CreateSingleVarNameDeclarationWithValue(Es6TemplateLiterals_TEMPLATELIT_VAR+uniqueId, callTemplateTagArgCreator).
		SetJSDocInfo(info).
		SrcrefTreeIfMissing(n)

	tagFnFirstArgDeclaration.InsertBefore(insertBefore)
	t.ReportCodeChangeWithNode(tagFnFirstArgDeclaration)

	tagFnFirstArg := tagFnFirstArgDeclaration.GetFirstChild().CloneNode()
	call := e.astFactory.CreateCall(n.RemoveFirstChild(), AstFactory_Type_(n), tagFnFirstArg)
	for part := range templateLit.Parts() {
		if sub, isSub := part.AsTemplateLitSub(); isSub {
			call.AddChildToBack(sub.Expr().Node().Detach())
		}
	}
	call.SrcrefTreeIfMissing(templateLit.Node())
	callTyped := ast.MustCall(call)
	_, isPropOrElemGet := ast.AsNormalOrOptChainGet(callTyped.Callee().Node())
	callTyped.SetFreeCall(!isPropOrElemGet)
	n.ReplaceWith(call)
	t.ReportCodeChange()
}

func (e *Es6TemplateLiterals) createRawStringArray(tpl ast.TemplateLit) rhino.Node {
	array := e.astFactory.CreateArraylit()
	for part := range tpl.Parts() {
		if tls, ok := part.AsTemplateLitString(); ok {
			array.AddChildToBack(e.astFactory.CreateString(tls.Raw()))
		}
	}
	return array
}

func (e *Es6TemplateLiterals) createCookedStringArray(tpl ast.TemplateLit) rhino.Node {
	var templateArrayType *colors.Color
	if e.astFactory.IsAddingColors() {
		templateArrayType = e.compiler.GetColorRegistry().Get(colors.StandardColors_I_TEMPLATE_ARRAY_ID)
	}
	array := e.astFactory.CreateArraylit().SetColor(templateArrayType)
	for part := range tpl.Parts() {
		if tls, ok := part.AsTemplateLitString(); ok {
			if cooked, hasCooked := tls.Cooked(); hasCooked {
				array.AddChildToBack(e.astFactory.CreateString(cooked))
			} else {
				array.AddChildToBack(e.astFactory.CreateVoid(e.astFactory.CreateNumber(0)))
			}
		}
	}
	return array
}

func cookedAndRawStringsSame(tpl ast.TemplateLit) bool {
	for part := range tpl.Parts() {
		if tls, ok := part.AsTemplateLitString(); ok {
			cooked, hasCooked := tls.Cooked()
			if !hasCooked || cooked != tls.Raw() {
				return false
			}
		}
	}
	return true
}
