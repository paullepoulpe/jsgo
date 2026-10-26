// Port of com/google/javascript/jscomp/JsMessageExtractor.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsMessageExtractor extracts messages and message comments from JS code.
//
// Uses a special prefix (e.g. MSG_) to determine which variables are
// messages. Here are the recognized formats:
//
//	var MSG_FOO = "foo";
//	var MSG_FOO_HELP = "this message is used for foo";
//
//	var MSG_BAR = function(a, b) {
//	  return a + " bar " + b;
//	}
//	var MSG_BAR_HELP = "the bar message";
//
// This class enforces the policy that message variable names must be unique
// across all JS files.
//
// Java: final class JsMessageExtractor
type JsMessageExtractor struct {
	idGenerator             JsMessage_IdGenerator
	options                 *CompilerOptions
	extractExternalMessages bool
}

// Java: JsMessageExtractor(JsMessage.IdGenerator)
func NewJsMessageExtractor(idGenerator JsMessage_IdGenerator) *JsMessageExtractor {
	return NewJsMessageExtractorWithOptions(idGenerator, NewCompilerOptions() /* extractExternalMessages= */, false)
}

// Java: JsMessageExtractor(JsMessage.IdGenerator, CompilerOptions, boolean)
func NewJsMessageExtractorWithOptions(
	idGenerator JsMessage_IdGenerator, options *CompilerOptions, extractExternalMessages bool) *JsMessageExtractor {
	e := &JsMessageExtractor{}
	e.idGenerator = idGenerator
	e.options = options
	e.extractExternalMessages = extractExternalMessages
	return e
}

// JsMessageExtractor_ExtractMessagesVisitor is the visitor that collects
// messages.
//
// Java: class JsMessageExtractor.ExtractMessagesVisitor extends JsMessageVisitor
type JsMessageExtractor_ExtractMessagesVisitor struct {
	JsMessageVisitorBase
	this *JsMessageExtractor // JsMessageExtractor.this

	// We use List here as we want to preserve insertion-order for found
	// messages.
	// Take into account that messages with the same id could be present in
	// the result list. Message could have the same id only in case if they
	// are unnamed and have the same text but located in different source
	// files.
	messages []*JsMessage
}

// Java: ExtractMessagesVisitor(AbstractCompiler)
func newJsMessageExtractor_ExtractMessagesVisitor(this *JsMessageExtractor, compiler AbstractCompiler) *JsMessageExtractor_ExtractMessagesVisitor {
	v := &JsMessageExtractor_ExtractMessagesVisitor{this: this}
	v.InitJsMessageVisitorBase(v, compiler, this.idGenerator)
	return v
}

// Java: ExtractMessagesVisitor.processJsMessageDefinition(JsMessageDefinition)
func (v *JsMessageExtractor_ExtractMessagesVisitor) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	v.processJsMessage(definition.GetMessage())
}

// Java: ExtractMessagesVisitor.processIcuTemplateDefinition(IcuTemplateDefinition)
func (v *JsMessageExtractor_ExtractMessagesVisitor) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	v.processJsMessage(definition.GetMessage())
}

// Java: ExtractMessagesVisitor.processJsMessage(JsMessage)
func (v *JsMessageExtractor_ExtractMessagesVisitor) processJsMessage(jsMessage *JsMessage) {
	if v.this.extractExternalMessages || !jsMessage.IsExternal() {
		v.messages = append(v.messages, jsMessage)
	}
}

// GetMessages returns extracted messages: the collection of JsMessage
// objects that was found in js sources.
//
// Java: ExtractMessagesVisitor.getMessages()
func (v *JsMessageExtractor_ExtractMessagesVisitor) GetMessages() []*JsMessage {
	return v.messages
}

// ExtractMessages extracts JS messages from JavaScript code.
//
// Java: JsMessageExtractor.extractMessages(SourceFile...)
func (e *JsMessageExtractor) ExtractMessages(inputs ...*SourceFile) []*JsMessage {
	return e.ExtractMessagesFromIterable(append([]*SourceFile{}, inputs...))
}

// ExtractMessagesFromIterable extracts JS messages from JavaScript code.
//
// inputs are the JavaScript source code inputs. It returns the extracted
// messages collection. It throws RuntimeException if there are problems
// parsing the JS code or the JS messages, or if two messages have the same
// key.
//
// Java: JsMessageExtractor.extractMessages(Iterable<SourceFile>)
func (e *JsMessageExtractor) ExtractMessagesFromIterable(inputs []*SourceFile) []*JsMessage {
	compiler := NewCompiler()
	compiler.Init([]*SourceFile{}, append([]*SourceFile{}, inputs...), e.options)
	Compiler_RunInCompilerThread(compiler,
		func() any {
			compiler.ParseInputs()
			return nil
		})

	extractCompilerPass := newJsMessageExtractor_ExtractMessagesVisitor(e, compiler)
	if len(compiler.GetErrors()) == 0 {
		extractCompilerPass.Process(rhino.NilNode, compiler.GetRoot())
	}

	errors := compiler.GetErrors()
	// Check for errors.
	if len(errors) != 0 {
		var msg strings.Builder
		msg.WriteString("JSCompiler errors\n")
		var formatter MessageFormatter = NewLightweightMessageFormatter(compiler)
		for _, err := range errors {
			msg.WriteString(formatter.FormatError(err))
		}
		panic(jutil.NewRuntimeException(msg.String()))
	}

	return extractCompilerPass.GetMessages()
}
