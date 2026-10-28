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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsMessageExtractor extracts messages and message comments from JS code.
type JsMessageExtractor struct {
	idGenerator             JsMessage_IdGenerator
	options                 *CompilerOptions
	extractExternalMessages bool
}

// Java: JsMessageExtractor(JsMessage.IdGenerator)
func NewJsMessageExtractor(idGenerator JsMessage_IdGenerator) *JsMessageExtractor {
	return NewJsMessageExtractorWithOptions(idGenerator, NewCompilerOptions(), false)
}

// Java: JsMessageExtractor(JsMessage.IdGenerator, CompilerOptions, boolean)
func NewJsMessageExtractorWithOptions(
	idGenerator JsMessage_IdGenerator, options *CompilerOptions, extractExternalMessages bool,
) *JsMessageExtractor {
	return &JsMessageExtractor{
		idGenerator:             idGenerator,
		options:                 options,
		extractExternalMessages: extractExternalMessages,
	}
}

// ExtractMessages extracts JS messages from JavaScript code.
//
// Java: JsMessageExtractor.extractMessages(SourceFile...)
func (e *JsMessageExtractor) ExtractMessages(inputs ...*SourceFile) []*JsMessage {
	return e.ExtractMessagesFromIterable(inputs)
}

// ExtractMessagesFromIterable extracts JS messages from JavaScript code.
//
// Java: JsMessageExtractor.extractMessages(Iterable<SourceFile>)
func (e *JsMessageExtractor) ExtractMessagesFromIterable(inputs []*SourceFile) []*JsMessage {
	compiler := NewCompiler()
	compiler.Init([]*SourceFile{}, slices.Clone(inputs), e.options)
	Compiler_RunInCompilerThread(compiler,
		func() any {
			compiler.ParseInputs()
			return nil
		})

	var messages []*JsMessage
	recordMessage := func(jsMessage *JsMessage) {
		if e.extractExternalMessages || !jsMessage.IsExternal() {
			messages = append(messages, jsMessage)
		}
	}
	if len(compiler.GetErrors()) == 0 {
		visitor := NewJsMessageVisitor(
			compiler,
			e.idGenerator,
			func(def JsMessageDefinition) { recordMessage(def.GetMessage()) },
			func(def IcuTemplateDefinition) { recordMessage(def.GetMessage()) },
			nil,
		)
		visitor.Process(rhino.NilNode, compiler.GetRoot())
	}

	errors := compiler.GetErrors()
	if len(errors) != 0 {
		var msg strings.Builder
		msg.WriteString("JSCompiler errors\n")
		formatter := NewLightweightMessageFormatter(compiler)
		for _, err := range errors {
			msg.WriteString(formatter.FormatError(err))
		}
		panic(jutil.NewRuntimeException(msg.String()))
	}

	return messages
}
