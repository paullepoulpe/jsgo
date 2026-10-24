// Port of com/google/javascript/jscomp/lint/CheckArrayWithGoogObject.java @ closure-compiler v20261015

/*
 * Copyright 2015 The Closure Compiler Authors.
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

// CheckArrayWithGoogObject lints against passing arrays to goog.object
// methods with the intention of iterating over them as though with a for-in
// loop, which is discouraged with arrays.
type CheckArrayWithGoogObject struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

var checkArrayWithGoogObject_GOOG_OBJECT_METHODS = jutil.NewLinkedHashSet(
	"goog.object.forEach",
	"goog.object.filter",
	"goog.object.map",
	"goog.object.some",
	"goog.object.every",
	"goog.object.getCount",
	"goog.object.getAnyKey",
	"goog.object.getAnyValue",
	"goog.object.contains",
	"goog.object.getValues",
	"goog.object.getKeys",
	"goog.object.findKey",
	"goog.object.findValue",
	"goog.object.isEmpty",
	"goog.object.clear",
	"goog.object.remove",
	"goog.object.equals",
	"goog.object.clone",
	"goog.object.transpose")

var CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT = DiagnosticType_Warning(
	"JSC_ARRAY_PASSED_TO_GOOG_OBJECT",
	"{0} expects an object, not an array. Did you mean to use goog.array?")

// NewCheckArrayWithGoogObject is new CheckArrayWithGoogObject(compiler).
//
// Java: CheckArrayWithGoogObject(AbstractCompiler)
func NewCheckArrayWithGoogObject(compiler AbstractCompiler) *CheckArrayWithGoogObject {
	return &CheckArrayWithGoogObject{compiler: compiler}
}

// Java: CheckArrayWithGoogObject.isGoogObjectIterationOverArray(Node)
func (c *CheckArrayWithGoogObject) IsGoogObjectIterationOverArray(n rhino.Node) bool {
	if !n.IsCall() {
		return false
	}
	if !n.GetFirstChild().IsQualifiedName() {
		return false
	}

	name := *n.GetFirstChild().GetQualifiedName()
	if !checkArrayWithGoogObject_GOOG_OBJECT_METHODS.Contains(name) {
		return false
	}

	firstArg := n.GetSecondChild()
	if firstArg.IsNil() {
		return false
	}
	type_ := firstArg.GetJSType()
	return type_ != nil && c.ContainsArray(type_)
}

// Java: CheckArrayWithGoogObject.containsArray(JSType)
func (c *CheckArrayWithGoogObject) ContainsArray(type_ rhino.JSType) bool {
	// Check if type is itself an array
	if type_.IsArrayType() {
		return true
	}
	templatizedType := type_.ToMaybeTemplatizedType()
	if templatizedType != nil && templatizedType.GetReferencedType().IsArrayType() {
		return true
	}

	// Check if this is a union that contains an array
	if type_.IsUnionType() {
		arrayType := c.compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)
		for _, alternate := range type_.ToMaybeUnionType().GetAlternates() {
			if alternate.IsSubtypeOf(arrayType) {
				return true
			}
		}
	}
	return false
}

// Java: CheckArrayWithGoogObject.visit(NodeTraversal, Node, Node)
func (c *CheckArrayWithGoogObject) Visit(unused *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if c.IsGoogObjectIterationOverArray(n) {
		c.compiler.Report(
			JSError_MakeWithNode(n, CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT, *n.GetFirstChild().GetQualifiedName()))
	}
}

// Java: CheckArrayWithGoogObject.process(Node, Node)
func (c *CheckArrayWithGoogObject) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}
