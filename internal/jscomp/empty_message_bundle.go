// Port of com/google/javascript/jscomp/EmptyMessageBundle.java @ closure-compiler v20260415

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
	"iter"
)

// EmptyMessageBundle is an implementation of MessageBundle that has no
// translations.
type EmptyMessageBundle struct{}

// NewEmptyMessageBundle is new EmptyMessageBundle().
func NewEmptyMessageBundle() *EmptyMessageBundle { return &EmptyMessageBundle{} }

// IdGenerator gets a dummy message ID generator.
//
// Java: EmptyMessageBundle.idGenerator()
func (b *EmptyMessageBundle) IdGenerator() JsMessage_IdGenerator {
	return nil
}

// GetMessage returns null, to indicate it has no message replacements.
//
// Java: EmptyMessageBundle.getMessage(String)
func (b *EmptyMessageBundle) GetMessage(id string) *JsMessage {
	return nil
}

// GetAllMessages returns an empty list of messages.
//
// Java: EmptyMessageBundle.getAllMessages()
func (b *EmptyMessageBundle) GetAllMessages() iter.Seq[*JsMessage] {
	return func(yield func(*JsMessage) bool) {}
}
