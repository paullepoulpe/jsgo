// Port of com/google/javascript/jscomp/Es6TemplateLiterals.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: Es6TemplateLiterals.TEMPLATELIT_VAR
const Es6TemplateLiterals_TEMPLATELIT_VAR = "$jscomp$templatelit$"

// Es6TemplateLiterals is a helper class for transpiling ES6 template
// literals.
type Es6TemplateLiterals struct {
	astFactory                       *AstFactory
	compiler                         AbstractCompiler
	namespace                        rhino.StaticScope
	createTemplateTagFirstArg        RuntimeJsLibManager_JsLibField
	createTemplateTagFirstArgWithRaw RuntimeJsLibManager_JsLibField
}

// Java: Es6TemplateLiterals(AbstractCompiler)
func NewEs6TemplateLiterals(compiler AbstractCompiler) *Es6TemplateLiterals {
	e := &Es6TemplateLiterals{}
	e.compiler = compiler
	e.astFactory = compiler.CreateAstFactory()
	e.namespace = e.compiler.GetTranspilationNamespace()

	runtimeJsLibManager := compiler.GetRuntimeJsLibManager()
	e.createTemplateTagFirstArg = runtimeJsLibManager.GetJsLibField("$jscomp.createTemplateTagFirstArg")
	e.createTemplateTagFirstArgWithRaw =
		runtimeJsLibManager.GetJsLibField("$jscomp.createTemplateTagFirstArgWithRaw")
	return e
}

// VisitTemplateLiteral converts `${a} b ${c} d ${e}` to (a + " b " + c + "
// d " + e).
//
// n is a TEMPLATELIT node that is not prefixed with a tag.
//
// Java: Es6TemplateLiterals.visitTemplateLiteral(NodeTraversal, Node)
func (e *Es6TemplateLiterals) VisitTemplateLiteral(t *NodeTraversal, n *rhino.Node) {
	length := n.GetChildCount()
	if length == 0 {
		n.ReplaceWith(e.astFactory.CreateString("\"\""))
	} else {
		first := n.RemoveFirstChild()
		jutil.CheckState(first.IsTemplateLitString() && first.GetCookedString() != nil)
		firstStr := e.astFactory.CreateString(*first.GetCookedString())
		if length == 1 {
			n.ReplaceWith(firstStr)
		} else {
			// Add the first string with the first substitution expression
			add := e.astFactory.CreateAdd(firstStr, n.RemoveFirstChild().RemoveFirstChild())
			// Process the rest of the template literal
			for i := 2; i < length; i++ {
				child := n.RemoveFirstChild()
				if child.IsTemplateLitString() {
					jutil.CheckState(child.GetCookedString() != nil)
					if *child.GetCookedString() == "" {
						continue
					} else if i == 2 && *first.GetCookedString() == "" {
						// So that `${hello} world` gets translated into (hello + " world")
						// instead of ("" + hello + " world").
						add = add.GetSecondChild().Detach()
					}
				}
				var right *rhino.Node
				if child.IsTemplateLitString() {
					right = e.astFactory.CreateString(*child.GetCookedString())
				} else {
					right = child.RemoveFirstChild()
				}
				add = e.astFactory.CreateAdd(add, right)
			}
			n.ReplaceWith(add.SrcrefTreeIfMissing(n))
		}
	}
	t.ReportCodeChange()
}

// VisitTaggedTemplateLiteral converts a tagged template into a call to the
// tag.
//
// If the cooked and raw strings of the template literal are same, this will
// create a call to createtemplatetagfirstarg that simply calls slice() on the
// cooked array at runtime to make the raw array a copy of the cooked array,
// and returns the cooked array. For example, tag`ab${bar}` will change to:
//
//	// A call to the tagging function:
//	tag($jscomp$templatelit$createTemplateTagFirstArg(['ab'], bar);
//
// If the cooked and raw strings of the template literal are not same, this
// will construct the raw strings array, and call
// createtemplatetagfirstargwithraw with it, which assigns the raw strings
// array to the property 'raw' of the cooked array at runtime, and returns
// the cooked array. For example, tag`a\tb${bar}` will change to:
//
//	// A call to the tagging function:
//	tag($jscomp$templatelit$createTemplateTagFirstArgWithRaw(['a\tb'] /, ['a\\tb']), bar);
//
// See template_literal_test.js for more examples.
//
// n is a TAGGED_TEMPLATELIT node.
//
// Java: Es6TemplateLiterals.visitTaggedTemplateLiteral(NodeTraversal, Node, Node)
func (e *Es6TemplateLiterals) VisitTaggedTemplateLiteral(t *NodeTraversal, n *rhino.Node, insertBefore *rhino.Node) {
	templateLit := n.GetLastChild()
	siteObject := e.CreateCookedStringArray(templateLit)

	// Node holding the function call to the runtime injected function
	var callTemplateTagArgCreator *rhino.Node
	if Es6TemplateLiterals_CookedAndRawStringsSame(templateLit) {
		// The cooked and raw versions of the array are the same, so just call slice() on the
		// cooked array at runtime to make the raw array a copy of the cooked array.
		callTemplateTagArgCreator = e.astFactory.CreateCall(
			e.astFactory.CreateQNameOfField(e.namespace, e.createTemplateTagFirstArg),
			AstFactory_Type_(siteObject),
			siteObject.CloneTree())
	} else {
		// The raw string array is different, so we need to construct it.
		raw := e.CreateRawStringArray(templateLit)
		callTemplateTagArgCreator = e.astFactory.CreateCall(
			e.astFactory.CreateQNameOfField(e.namespace, e.createTemplateTagFirstArgWithRaw),
			AstFactory_Type_(siteObject),
			siteObject.CloneTree(),
			raw)
	}
	jsDocInfoBuilder := rhino.JSDocInfo_Builder_()
	jsDocInfoBuilder.RecordNoInline()
	info := jsDocInfoBuilder.Build()

	input := t.GetInput()
	uniqueId := e.compiler.GetUniqueIdSupplier().GetUniqueId(input)
	// var tagFnFirstArg = $jscomp.createTemplateTagFirstArg...
	tagFnFirstArgDeclaration := e.astFactory.
		CreateSingleVarNameDeclarationWithValue(Es6TemplateLiterals_TEMPLATELIT_VAR+uniqueId, callTemplateTagArgCreator).
		SetJSDocInfo(info).
		SrcrefTreeIfMissing(n)

	// For the first script, insertion point is right after the runtime injected function definition
	// as injecting to the top of the script causes runtime errors
	// https://github.com/google/closure-compiler/issues/3589. For any subsequent script(s), the
	// call is injected to the top of that script.
	tagFnFirstArgDeclaration.InsertBefore(insertBefore)
	t.ReportCodeChangeWithNode(tagFnFirstArgDeclaration)

	// Generate the call expression.
	tagFnFirstArg := tagFnFirstArgDeclaration.GetFirstChild().CloneNode()
	call := e.astFactory.CreateCall(n.RemoveFirstChild(), AstFactory_Type_(n), tagFnFirstArg)
	for child := templateLit.GetFirstChild(); child != nil; child = child.GetNext() {
		if !child.IsTemplateLitString() {
			call.AddChildToBack(child.RemoveFirstChild())
		}
	}
	call.SrcrefTreeIfMissing(templateLit)
	call.PutBooleanProp(rhino.Node_FREE_CALL, !NodeUtil_IsNormalOrOptChainGet(call.GetFirstChild()))
	n.ReplaceWith(call)
	t.ReportCodeChange()
}

// Java: Es6TemplateLiterals.CreateRawStringArray(Node)
func (e *Es6TemplateLiterals) CreateRawStringArray(n *rhino.Node) *rhino.Node {
	array := e.astFactory.CreateArraylit()
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		if child.IsTemplateLitString() {
			array.AddChildToBack(e.astFactory.CreateString(child.GetRawString()))
		}
	}
	return array
}

// Java: Es6TemplateLiterals.CreateCookedStringArray(Node)
func (e *Es6TemplateLiterals) CreateCookedStringArray(n *rhino.Node) *rhino.Node {
	var templateArrayType *colors.Color
	if e.astFactory.IsAddingColors() {
		templateArrayType = e.compiler.GetColorRegistry().Get(colors.StandardColors_I_TEMPLATE_ARRAY_ID)
	}
	array := e.astFactory.
		CreateArraylit().
		SetColor(templateArrayType) // tighten the type from Array to ITemplateArray
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		if child.IsTemplateLitString() {
			if child.GetCookedString() != nil {
				array.AddChildToBack(e.astFactory.CreateString(*child.GetCookedString()))
			} else {
				// undefined cooked string due to exception in template escapes
				array.AddChildToBack(e.astFactory.CreateVoid(e.astFactory.CreateNumber(0)))
			}
		}
	}
	return array
}

// Java: Es6TemplateLiterals.cookedAndRawStringsSame(Node)
func Es6TemplateLiterals_CookedAndRawStringsSame(n *rhino.Node) bool {
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		if !child.IsTemplateLitString() {
			continue
		}
		// getCookedString() returns null when the template literal has an illegal escape sequence.
		if child.GetCookedString() == nil || *child.GetCookedString() != child.GetRawString() {
			return false
		}
	}
	return true
}
