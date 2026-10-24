// Port of com/google/javascript/jscomp/ReplaceMessagesConstants.java @ closure-compiler v20261021

/*
 * Copyright 2021 The Closure Compiler Authors.
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

// ReplaceMessagesConstants holds constants that need to be shared between
// `ReplaceMessages` and other passes.
type ReplaceMessagesConstants struct{}

const (
	// ReplaceMessagesConstants_DEFINE_MSG_CALLEE: `goog.getMsg()` calls will
	// be converted into a call to this method which is defined in synthetic
	// externs.
	ReplaceMessagesConstants_DEFINE_MSG_CALLEE = "__jscomp_define_msg__"

	// ReplaceMessagesConstants_FALLBACK_MSG_CALLEE:
	// `goog.getMsgWithFallback(MSG_NEW, MSG_OLD)` will be converted into a
	// call to this method which is defined in * synthetic externs.
	ReplaceMessagesConstants_FALLBACK_MSG_CALLEE = "__jscomp_msg_fallback__"
)

// Java: ReplaceMessagesConstants.PROTECTED_FUNCTION_NAMES
var replaceMessagesConstants_PROTECTED_FUNCTION_NAMES = jutil.NewLinkedHashSet(
	ReplaceMessagesConstants_DEFINE_MSG_CALLEE, ReplaceMessagesConstants_FALLBACK_MSG_CALLEE)

// Java: ReplaceMessagesConstants.isProtectedMessage(Node)
func ReplaceMessagesConstants_IsProtectedMessage(n rhino.Node) bool {
	if !n.IsCall() {
		return false
	}
	callee := n.GetFirstChild()
	if !callee.IsName() {
		return false
	}
	return replaceMessagesConstants_PROTECTED_FUNCTION_NAMES.Contains(callee.GetString())
}
