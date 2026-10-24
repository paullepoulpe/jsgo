// Port of com/google/javascript/jscomp/ClosurePrimitiveErrors.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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

// ClosurePrimitiveErrors holds common errors for Closure primitives that are
// reported in multiple passses.
type ClosurePrimitiveErrors struct{}

func newClosurePrimitiveErrors() *ClosurePrimitiveErrors { return &ClosurePrimitiveErrors{} }

var (
	ClosurePrimitiveErrors_DUPLICATE_MODULE = DiagnosticType_Error(
		"JSC_DUPLICATE_MODULE",
		"Duplicate module: {0}. Module already declared in other file {1}")

	ClosurePrimitiveErrors_DUPLICATE_NAMESPACE = DiagnosticType_Error(
		"JSC_DUPLICATE_NAMESPACE",
		"Duplicate namespace: {0}. Namespace already provided in other file {1}")

	ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE = DiagnosticType_Error(
		"JSC_DUPLICATE_NAMESPACE_AND_MODULE",
		"Duplicate module: {0}. Module already goog.provided in other file {1}")

	ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE = DiagnosticType_Error(
		"JSC_INVALID_DESTRUCTURING_FORWARD_DECLARE",
		"Cannot destructure a forward-declared type")

	ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET = DiagnosticType_Error(
		"JSC_MODULE_USES_GOOG_MODULE_GET",
		"It's illegal to use a 'goog.module.get' at the module top-level."+
			" Did you mean to use goog.require instead?")

	ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE = DiagnosticType_Error(
		"JSC_GOOG_MODULE_INVALID_FORWARD_DECLARE_NAMESPACE",
		"goog.forwardDeclare parameter must be a string literal.")

	ClosurePrimitiveErrors_INVALID_GET_NAMESPACE = DiagnosticType_Error(
		"JSC_GOOG_MODULE_INVALID_GET_NAMESPACE",
		"goog.module.get parameter must be a string literal.")

	ClosurePrimitiveErrors_INVALID_REQUIRE_DYNAMIC = DiagnosticType_Error(
		"JSC_GOOG_MODULE_INVALID_REQUIRE_DYNAMIC_NAMESPACE",
		"goog.requireDynamic parameter must be a string literal.")

	ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE = DiagnosticType_Error(
		"JSC_GOOG_MODULE_INVALID_REQUIRE_NAMESPACE",
		"goog.require parameter must be a string literal.")

	ClosurePrimitiveErrors_INVALID_REQUIRE_TYPE_NAMESPACE = DiagnosticType_Error(
		"JSC_GOOG_MODULE_INVALID_REQUIRE_TYPE_NAMESPACE",
		"goog.requireType parameter must be a string literal.")

	ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE = DiagnosticType_Error(
		"JSC_MISSING_MODULE_OR_PROVIDE", "Required namespace \"{0}\" never defined.")

	ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE = DiagnosticType_Error(
		"JSC_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE",
		"Required namespace \"{0}\" never defined.")

	ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE = DiagnosticType_Error(
		"JSC_GOOG_MODULE_INVALID_GET_CALL_SCOPE",
		"goog.module.get values should not be assigned to global variables.")

	ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR = DiagnosticType_Error(
		"JSC_INVALID_CLOSURE_CALL_ERROR",
		"Closure primitive methods (goog.provide, goog.require, goog.define, etc) must be called "+
			"at file scope.")

	ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR = DiagnosticType_Error(
		"JSC_TOO_MANY_ARGUMENTS_ERROR", "method \"{0}\" called with more than one argument")

	ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR = DiagnosticType_Error("JSC_NULL_ARGUMENT_ERROR", "method \"{0}\" called without an argument")

	ClosurePrimitiveErrors_INVALID_ARGUMENT_ERROR = DiagnosticType_Error(
		"JSC_INVALID_ARGUMENT_ERROR", "method \"{0}\" called with invalid argument")

	ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE = DiagnosticType_Error(
		"JSC_GOOG_MODULE_GET_OF_WEAK_MODULE",
		"goog.module.get of \"{0}\" references a namespace that is only imported by"+
			" goog.requireType.")
)
