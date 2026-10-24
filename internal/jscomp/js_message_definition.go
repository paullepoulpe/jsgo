// Port of com/google/javascript/jscomp/JsMessageDefinition.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsMessageDefinition is the container class that holds a JsMessage and
// information about the `goog.getMsg()` call it was built from.
type JsMessageDefinition interface {
	// GetMessage is the JsMessage object built from the `goog.getMsg()`
	// call.
	//
	// Java: JsMessageDefinition.getMessage()
	GetMessage() *JsMessage

	// GetMessageNode is the RHS node of the message assignment statement.
	//
	// e.g. The goog.getMsg() call in:
	//
	//	const MSG_HELLO =
	//	    goog.getMsg(
	//	        'Hello, {$name}.',
	//	        { 'name', getName() },
	//	        { unescapeHtmlEntities: true });
	//
	// Java: JsMessageDefinition.getMessageNode()
	GetMessageNode() rhino.Node

	// GetTemplateTextNode is the Node representing the message text
	// template.
	//
	// This node may be a literal string or a concatenation of literal
	// strings.
	//
	// For a goog.getMsg() call this is the first argument.
	//
	// Java: JsMessageDefinition.getTemplateTextNode()
	GetTemplateTextNode() rhino.Node

	// GetPlaceholderValuesNode is the object literal Node that maps
	// placeholder names to expressions providing their values.
	//
	// This value will be null if the message definition didn't specify
	// placeholder values.
	//
	// Java: JsMessageDefinition.getPlaceholderValuesNode()
	GetPlaceholderValuesNode() rhino.Node

	// GetPlaceholderValueMap is a map from placehlolder name to the Node
	// assigned to it in the values map argument of `goog.getMsg()`.
	//
	// This will be an empty map if there was no object or it was an empty
	// object literal.
	//
	// Java: JsMessageDefinition.getPlaceholderValueMap()
	GetPlaceholderValueMap() *jutil.LinkedHashMap[string, rhino.Node]

	// ShouldEscapeLessThan is the value of the 'html' options key in the
	// options bag argument.
	//
	// This value will be false if there was no options bag argument, or if
	// it didn't contain an 'html' property.
	//
	// Java: JsMessageDefinition.shouldEscapeLessThan()
	ShouldEscapeLessThan() bool

	// ShouldUnescapeHtmlEntities is the value of the 'unescapeHtmlEntities'
	// options key in the options bag argument.
	//
	// This value will be false if there was no options bag argument, or if
	// it didn't contain an 'unescapeHtmlEntities' property.
	//
	// Java: JsMessageDefinition.shouldUnescapeHtmlEntities()
	ShouldUnescapeHtmlEntities() bool
}
