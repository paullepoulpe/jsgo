// Port of com/google/javascript/jscomp/modules/ClosureRequireProcessor.java @ closure-compiler v20261021

/*
 * Copyright 2019 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ClosureRequireProcessor handles creating an Import from
// goog.require(Type) or goog.forwardDeclare.
//
// This logic can be used by both goog.modules and ES modules
type ClosureRequireProcessor struct {
	nameDeclaration rhino.Node
	requireKind     Binding_CreatedBy
}

// ClosureRequireProcessor_Require represents a goog.require(Type) or
// goog.forwardDeclare
//
// localName is the name local to the module with the require; e.g. `b` in
// `const b = goog.require('a');`. importRecord is an Import containing all
// metadata about this require. createdBy is whether this is a goog.require,
// goog.requireType, or goog.forwardDeclare.
type ClosureRequireProcessor_Require struct {
	localName    string
	importRecord *Import
	createdBy    Binding_CreatedBy
}

// Java: ClosureRequireProcessor.Require(String, Import, CreatedBy)
func NewClosureRequireProcessor_Require(localName string, importRecord *Import, createdBy Binding_CreatedBy) *ClosureRequireProcessor_Require {
	jutil.CheckNotNullMsg(importRecord, "importRecord")
	jutil.CheckArgument(createdBy.IsClosureImport())
	r := &ClosureRequireProcessor_Require{}
	r.localName = localName
	r.importRecord = importRecord
	r.createdBy = createdBy
	return r
}

// Java: ClosureRequireProcessor.Require.localName()
func (r *ClosureRequireProcessor_Require) LocalName() string { return r.localName }

// Java: ClosureRequireProcessor.Require.importRecord()
func (r *ClosureRequireProcessor_Require) ImportRecord() *Import { return r.importRecord }

// Java: ClosureRequireProcessor.Require.createdBy()
func (r *ClosureRequireProcessor_Require) CreatedBy() Binding_CreatedBy { return r.createdBy }

// Java: ClosureRequireProcessor.Require.create(String, Import, CreatedBy)
func ClosureRequireProcessor_Require_Create(localName string, importRecord *Import, createdBy Binding_CreatedBy) *ClosureRequireProcessor_Require {
	return NewClosureRequireProcessor_Require(localName, importRecord, createdBy)
}

// Java: ClosureRequireProcessor(Node, CreatedBy)
func newClosureRequireProcessor(nameDeclaration rhino.Node, requireKind Binding_CreatedBy) *ClosureRequireProcessor {
	p := &ClosureRequireProcessor{}
	jutil.CheckArgument(NodeUtil_IsNameDeclaration(nameDeclaration))
	p.nameDeclaration = nameDeclaration
	p.requireKind = requireKind
	return p
}

// ClosureRequireProcessor_GetAllRequires returns all Require built from the
// given statement, or null if it is not a require
//
// nameDeclaration is a VAR, LET, or CONST; returns all Requires contained in
// this declaration
//
// Java: ClosureRequireProcessor.getAllRequires(Node)
func ClosureRequireProcessor_GetAllRequires(nameDeclaration rhino.Node) []*ClosureRequireProcessor_Require {
	var rhs rhino.Node
	if decl, ok := ast.AsNameDeclaration(nameDeclaration); ok {
		if first, hasFirst := decl.FirstDeclarator(); hasFirst {
			if dlhs, isDlhs := first.AsDestructuringLhs(); isDlhs {
				if init, hasInit := dlhs.Init(); hasInit {
					rhs = init.Node()
				}
			} else if name, isName := first.AsName(); isName {
				if init, hasInit := name.Init(); hasInit {
					rhs = init.Node()
				}
			}
		}
	}
	// This may be a require, requireType, or forwardDeclare.
	requireKind, ok := ClosureRequireProcessor_GetModuleDependencyTypeFromRhs(rhs)
	if !ok {
		return []*ClosureRequireProcessor_Require{}
	}

	return newClosureRequireProcessor(nameDeclaration, requireKind).GetAllRequiresInDeclaration()
}

// ClosureRequireProcessor_GetModuleDependencyTypeFromRhs checks if the given
// rvalue is a goog.require(Type) or goog.forwardDeclare call, and if so
// returns which one.
//
// Java: ClosureRequireProcessor.getModuleDependencyTypeFromRhs(Node)
func ClosureRequireProcessor_GetModuleDependencyTypeFromRhs(value rhino.Node) (Binding_CreatedBy, bool) {
	call, ok := ast.AsCall(value)
	if !ok {
		return 0, false
	}
	arg, ok := call.OnlyArg()
	if !ok || !arg.Node().IsStringLit() {
		return 0, false
	}
	callee, ok := call.Callee().AsGetProp()
	if !ok {
		return 0, false
	}
	owner, ok := callee.Object().AsName()
	if !ok || owner.Ident() != "goog" {
		return 0, false
	}
	switch callee.Property() {
	case "require":
		return Binding_CreatedBy_GOOG_REQUIRE, true
	case "requireType":
		return Binding_CreatedBy_GOOG_REQUIRE_TYPE, true
	case "forwardDeclare":
		return Binding_CreatedBy_GOOG_FORWARD_DECLARE, true
	default:
		return 0, false
	}
}

// GetAllRequiresInDeclaration returns a new list of all required names in
// nameDeclaration
//
// Java: ClosureRequireProcessor.getAllRequiresInDeclaration()
func (p *ClosureRequireProcessor) GetAllRequiresInDeclaration() []*ClosureRequireProcessor_Require {
	first := ast.MustNameDeclaration(p.nameDeclaration).MustFirstDeclarator()
	var rhs rhino.Node
	if dlhs, ok := first.AsDestructuringLhs(); ok {
		init, _ := dlhs.Init()
		rhs = init.Node()
	} else {
		init, _ := first.MustName().Init()
		rhs = init.Node()
	}

	namespace := ast.MustCall(rhs).MustOnlyArg().MustStringLit().Value()

	if lhs, ok := first.AsName(); ok {
		// const modA = goog.require('modA');
		return []*ClosureRequireProcessor_Require{
			ClosureRequireProcessor_Require_Create(
				lhs.Ident(),
				NewImport(namespace, Export_NAMESPACE, lhs.Ident(), nil, p.nameDeclaration, lhs.Node()),
				p.requireKind),
		}
	}
	// const {x, y} = goog.require('modA');
	objectPattern := first.MustDestructuringLhs().Pattern().Node()
	if !objectPattern.IsObjectPattern() {
		// bad JS, ignore
		return []*ClosureRequireProcessor_Require{}
	}
	return p.GetAllRequiresFromDestructuring(objectPattern, namespace)
}

// GetAllRequiresFromDestructuring returns all requires from destructruring,
// like `const {x, y, z} = goog.require('a');`
//
// Java: ClosureRequireProcessor.getAllRequiresFromDestructuring(Node, String)
func (p *ClosureRequireProcessor) GetAllRequiresFromDestructuring(
	objectPattern rhino.Node, namespace string) []*ClosureRequireProcessor_Require {
	requires := []*ClosureRequireProcessor_Require{}
	for prop := range ast.MustObjectPattern(objectPattern).Props() {
		key, ok := prop.AsStringKey()
		if !ok {
			// Bad code, just ignore. We warn elsewhere.
			continue
		}
		lhs, ok := key.Value().AsName()
		if !ok {
			// Bad code ( e.g. `const {a = 0} = goog.require(...)`). We warn elsewhere.
			continue
		}

		requires = append(requires,
			ClosureRequireProcessor_Require_Create(
				lhs.Ident(),
				NewImport(namespace, key.Key(), lhs.Ident(), nil, p.nameDeclaration, lhs.Node()),
				p.requireKind))
	}
	return requires
}
