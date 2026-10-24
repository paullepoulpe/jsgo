// Port of com/google/javascript/jscomp/ReplaceMessagesForChrome.java @ closure-compiler v20261021

/*
 * Copyright 2012 The Closure Compiler Authors.
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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReplaceMessagesForChrome replaces user-visible messages with appropriate
// calls to chrome.i18n.getMessage. The first argument to getMessage is the id
// of the message, as a string. If the message contains placeholders, the
// second argument is an array of the values being used for the placeholders,
// sorted by placeholder name.
//
// Java: class ReplaceMessagesForChrome extends JsMessageVisitor
type ReplaceMessagesForChrome struct {
	JsMessageVisitorBase

	astFactory             *AstFactory
	introducesRegexpSyntax bool
}

// Java: ReplaceMessagesForChrome.DECLARE_ICU_TEMPLATE_NOT_SUPPORTED
var ReplaceMessagesForChrome_DECLARE_ICU_TEMPLATE_NOT_SUPPORTED = DiagnosticType_Error(
	"JSC_DECLARE_ICU_TEMPLATE",
	"goog.i18n.messages.declareIcuTemplate() is not supported for Chrome i18n.")

// Java: ReplaceMessagesForChrome(AbstractCompiler, JsMessage.IdGenerator)
func NewReplaceMessagesForChrome(compiler AbstractCompiler, idGenerator JsMessage_IdGenerator) *ReplaceMessagesForChrome {
	r := &ReplaceMessagesForChrome{}
	r.InitJsMessageVisitorBase(r, compiler, idGenerator)
	r.introducesRegexpSyntax = false
	r.astFactory = compiler.CreateAstFactory()
	return r
}

// Java: ReplaceMessagesForChrome.getChromeI18nGetMessageNode(String)
func (r *ReplaceMessagesForChrome) GetChromeI18nGetMessageNode(messageId string) *rhino.Node {
	getMessage := r.astFactory.CreateQNameWithUnknownType("chrome.i18n.getMessage")
	return r.astFactory.CreateCall(
		getMessage, AstFactory_Type_OfColor(colors.StandardColors_STRING), r.astFactory.CreateString(messageId))
}

// Java: ReplaceMessagesForChrome.processIcuTemplateDefinition(IcuTemplateDefinition)
func (r *ReplaceMessagesForChrome) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	// TODO(bradfordcsmith): Add support for this.
	r.compiler.Report(JSError_MakeWithNode(definition.GetMessageNode(), ReplaceMessagesForChrome_DECLARE_ICU_TEMPLATE_NOT_SUPPORTED))
}

// Java: ReplaceMessagesForChrome.processJsMessageDefinition(JsMessageDefinition)
func (r *ReplaceMessagesForChrome) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	r.introducesRegexpSyntax = false
	newValue := r.GetNewValueNode(definition.GetMessage(), definition)
	definition.GetMessageNode().ReplaceWith(newValue)
	if r.introducesRegexpSyntax {
		NodeUtil_AddFeatureToScript(
			NodeUtil_GetEnclosingScript(newValue), parser.FeatureSet_Feature_REGEXP_SYNTAX, r.compiler)
		r.introducesRegexpSyntax = false
	}
	r.compiler.ReportChangeToEnclosingScope(newValue)
}

// Java: ReplaceMessagesForChrome.getNewValueNode(JsMessage, JsMessageDefinition)
func (r *ReplaceMessagesForChrome) GetNewValueNode(message *JsMessage, definition JsMessageDefinition) *rhino.Node {
	newValueNode := r.GetChromeI18nGetMessageNode(message.GetId())

	isHtml := definition.ShouldEscapeLessThan()
	if !message.JsPlaceholderNames().IsEmpty() || isHtml {
		// Output the placeholders, sorted alphabetically by placeholder name,
		// regardless of what order they appear in the original message.
		placeholderNames := message.JsPlaceholderNames().ToSlice()
		slices.SortStableFunc(placeholderNames, jutil.StringCompareTo) // Ordering.natural().sortedCopy

		placeholderValueArray := r.astFactory.CreateArraylit()
		placeholderValueMap := definition.GetPlaceholderValueMap()
		for _, name := range placeholderNames {
			// JsMessageVisitor ensures that every placeholder name appearing in the message string
			// has a corresponding value node. It will report an error and avoid passing any messages
			// violating this invariant to processMessage().
			value := jutil.CheckNotNull(placeholderValueMap.Get(name)).CloneTree()
			placeholderValueArray.AddChildToBack(value)
		}
		if isHtml {
			args := r.astFactory.CreateArraylit(
				r.astFactory.CreateString(message.GetId()), placeholderValueArray)
			options := r.astFactory.CreateArraylit(
				r.astFactory.CreateObjectLit(
					r.astFactory.CreateStringKey("escapeLt", r.astFactory.CreateBoolean(true))))
			regexp := rhino.IR_Regexp(r.astFactory.CreateString("Chrome\\/(\\d+)")).SetColor(colors.StandardColors_UNKNOWN)
			userAgent := r.astFactory.CreateQNameWithUnknownType("navigator.userAgent")
			version := r.astFactory.CreateGetElem(
				r.astFactory.CreateOr(
					r.astFactory.CreateCall(
						r.astFactory.CreateGetPropWithUnknownType(regexp, "exec"),
						AstFactory_Type_OfColorId(colors.StandardColors_ARRAY_ID),
						userAgent),
					r.astFactory.CreateArraylit()),
				r.astFactory.CreateNumber(1))
			condition := rhino.IR_Ge(version, r.astFactory.CreateNumber(79)).SetColor(colors.StandardColors_BOOLEAN)
			args = r.astFactory.CreateCall(
				r.astFactory.CreateGetPropWithUnknownType(args, "concat"),
				AstFactory_Type_OfColorId(colors.StandardColors_ARRAY_ID),
				r.astFactory.CreateHook(condition, options, r.astFactory.CreateArraylit()))
			newValueNode = r.astFactory.CreateCallWithUnknownType(
				r.astFactory.CreateQNameWithUnknownType("chrome.i18n.getMessage.apply"),
				r.astFactory.CreateNull(),
				args)
			r.introducesRegexpSyntax = true
		} else {
			newValueNode.AddChildToBack(placeholderValueArray)
		}
	}

	newValueNode.SrcrefTreeIfMissing(definition.GetMessageNode())
	return newValueNode
}
