// Port of com/google/javascript/jscomp/IcuTemplateDefinition.java @ closure-compiler v20261015

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// IcuTemplateDefinition contains a JsMessage representing an ICU template
// message and information about the
// `goog.i18n.messages.declareIcuTemplate()` call from which it was
// extracted.
type IcuTemplateDefinition interface {
	// GetMessage is the JsMessage object built from the
	// `declareIcuTemplate()` call.
	//
	// Java: IcuTemplateDefinition.getMessage()
	GetMessage() *JsMessage

	// GetMessageNode is the RHS node of the message assignment statement.
	//
	// e.g. The declareIcuTemplate() call in:
	//
	//	const MSG_HELLO =
	//	    declareIcuTemplate(
	//	        'Hello, {NAME}.',
	//	        {
	//	          description: 'A greeting',
	//	          example: {
	//	            'NAME': 'Jane'
	//	          });
	//
	// Java: IcuTemplateDefinition.getMessageNode()
	GetMessageNode() *rhino.Node

	// GetTemplateTextNode is the Node representing the message text
	// template.
	//
	// This node may be a literal string or a concatenation of literal
	// strings.
	//
	// Java: IcuTemplateDefinition.getTemplateTextNode()
	GetTemplateTextNode() *rhino.Node
}
