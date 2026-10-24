// Port of com/google/javascript/jscomp/type/FlowScope.java @ closure-compiler v20261021

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

// FlowScope is a symbol table for inferring types during data flow
// analysis.
//
// Each flow scope represents the types of all variables in the scope at a
// particular point in the flow analysis.
//
// Go-only: language-adapter — getParentScope() and getSlot(String) keep
// StaticScope's Go signatures (see rhino.StaticTypedScope); FlowScope
// implementations return StaticTypedScope / StaticTypedSlot values from
// them.
type FlowScope interface {
	rhino.StaticTypedScope
	LatticeElement

	// WithSyntacticScope returns a flow scope with the given syntactic
	// scope, which may be required to be a specific subclass, such as
	// TypedScope.
	//
	// Java: FlowScope.withSyntacticScope(StaticTypedScope)
	WithSyntacticScope(scope rhino.StaticTypedScope) FlowScope

	// InferSlotType returns a flow scope with the type of the given symbol
	// updated to type. Throws IllegalArgumentException if no slot for this
	// symbol exists.
	//
	// Java: FlowScope.inferSlotType(String, JSType)
	InferSlotType(symbol string, type_ rhino.JSType) FlowScope

	// InferQualifiedSlot returns a flow scope with the type of the given
	// symbol updated to inferredType. Updates are not performed in-place.
	//
	// When traversing the control flow of a function, simple names are
	// declared at the bottom of the flow lattice. But there are far too many
	// qualified names to be able to do this and be performant. So the
	// bottoms of qualified names are declared lazily.
	//
	// Therefore, when inferring a qualified slot, we need both the "bottom"
	// type of the slot when we enter the scope, and the current type being
	// inferred.
	//
	// Java: FlowScope.inferQualifiedSlot(Node, String, JSType, JSType, boolean)
	InferQualifiedSlot(
		node *rhino.Node, symbol string, bottomType rhino.JSType, inferredType rhino.JSType, declare bool) FlowScope

	// GetDeclarationScope returns the underlying TypedScope.
	//
	// Java: FlowScope.getDeclarationScope()
	GetDeclarationScope() rhino.StaticTypedScope
}
