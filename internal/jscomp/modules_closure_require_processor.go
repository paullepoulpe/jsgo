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

// Equals is the record's component-wise equals.
//
// Java: ClosureRequireProcessor.Require.equals(Object)
func (r *ClosureRequireProcessor_Require) Equals(o any) bool {
	that, ok := o.(*ClosureRequireProcessor_Require)
	if !ok || that == nil {
		return false
	}
	return r.localName == that.localName &&
		r.importRecord.Equals(that.importRecord) &&
		r.createdBy == that.createdBy
}

// String is the record's toString.
//
// Java: ClosureRequireProcessor.Require.toString()
func (r *ClosureRequireProcessor_Require) String() string {
	return "Require[localName=" + r.localName +
		", importRecord=" + r.importRecord.String() +
		", createdBy=" + r.createdBy.String() + "]"
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
	if nameDeclaration.GetFirstChild().IsDestructuringLhs() {
		rhs = nameDeclaration.GetFirstChild().GetSecondChild()
	} else {
		rhs = nameDeclaration.GetFirstFirstChild()
	}
	// This may be a require, requireType, or forwardDeclare.
	requireKind := ClosureRequireProcessor_GetModuleDependencyTypeFromRhs(rhs)
	if requireKind == nil {
		return []*ClosureRequireProcessor_Require{}
	}

	return newClosureRequireProcessor(nameDeclaration, *requireKind).GetAllRequiresInDeclaration()
}

var closureRequireProcessor_GOOG_DEPENDENCY_CALLS = func() *jutil.LinkedHashMap[string, Binding_CreatedBy] {
	m := jutil.NewLinkedHashMap[string, Binding_CreatedBy]()
	m.Put("require", Binding_CreatedBy_GOOG_REQUIRE)
	m.Put("requireType", Binding_CreatedBy_GOOG_REQUIRE_TYPE)
	m.Put("forwardDeclare", Binding_CreatedBy_GOOG_FORWARD_DECLARE)
	return m
}()

// ClosureRequireProcessor_GetModuleDependencyTypeFromRhs checks if the given
// rvalue is a goog.require(Type) or goog.forwardDeclare call, and if so
// returns which one.
//
// Returns a Closure require (where CreatedBy#isClosureImport() is true) or
// null.
//
// Java: ClosureRequireProcessor.getModuleDependencyTypeFromRhs(Node)
func ClosureRequireProcessor_GetModuleDependencyTypeFromRhs(value rhino.Node) *Binding_CreatedBy {
	if value.IsNil() ||
		!value.IsCall() ||
		!value.HasTwoChildren() ||
		!value.GetSecondChild().IsStringLit() {
		return nil
	}
	callee := value.GetFirstChild()
	if !callee.IsGetProp() {
		return nil
	}
	owner := callee.GetFirstChild()
	if !owner.IsName() || owner.GetString() != "goog" {
		return nil
	}
	if createdBy, ok := closureRequireProcessor_GOOG_DEPENDENCY_CALLS.Lookup(callee.GetString()); ok {
		return &createdBy
	}
	return nil
}

// GetAllRequiresInDeclaration returns a new list of all required names in
// nameDeclaration
//
// Java: ClosureRequireProcessor.getAllRequiresInDeclaration()
func (p *ClosureRequireProcessor) GetAllRequiresInDeclaration() []*ClosureRequireProcessor_Require {
	var rhs rhino.Node
	if p.nameDeclaration.GetFirstChild().IsDestructuringLhs() {
		rhs = p.nameDeclaration.GetFirstChild().GetSecondChild()
	} else {
		rhs = p.nameDeclaration.GetFirstFirstChild()
	}

	namespace := rhs.GetSecondChild().GetString()

	if p.nameDeclaration.GetFirstChild().IsName() {
		// const modA = goog.require('modA');
		lhs := p.nameDeclaration.GetFirstChild()
		return []*ClosureRequireProcessor_Require{
			ClosureRequireProcessor_Require_Create(
				lhs.GetString(),
				Import_Builder_().
					ModuleRequest(namespace).
					LocalName(lhs.GetString()).
					ImportName(Export_NAMESPACE).
					ImportNode(p.nameDeclaration).
					NameNode(lhs).
					Build(),
				p.requireKind),
		}
	} else {
		// const {x, y} = goog.require('modA');
		objectPattern := p.nameDeclaration.GetFirstFirstChild()
		if !objectPattern.IsObjectPattern() {
			// bad JS, ignore
			return []*ClosureRequireProcessor_Require{}
		}
		return p.GetAllRequiresFromDestructuring(objectPattern, namespace)
	}
}

// GetAllRequiresFromDestructuring returns all requires from destructruring,
// like `const {x, y, z} = goog.require('a');`
//
// Java: ClosureRequireProcessor.getAllRequiresFromDestructuring(Node, String)
func (p *ClosureRequireProcessor) GetAllRequiresFromDestructuring(
	objectPattern rhino.Node, namespace string) []*ClosureRequireProcessor_Require {
	requireBuilder := jutil.ImmutableList_Builder_[*ClosureRequireProcessor_Require]()
	for key := objectPattern.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
		if !key.IsStringKey() {
			// Bad code, just ignore. We warn elsewhere.
			continue
		}
		lhs := key.GetOnlyChild()
		if !lhs.IsName() {
			// Bad code ( e.g. `const {a = 0} = goog.require(...)`). We warn elsewhere.
			continue
		}

		requireBuilder.Add(
			ClosureRequireProcessor_Require_Create(
				lhs.GetString(),
				Import_Builder_().
					ModuleRequest(namespace).
					LocalName(lhs.GetString()).
					ImportName(key.GetString()).
					ImportNode(p.nameDeclaration).
					NameNode(lhs).
					Build(),
				p.requireKind))
	}
	return requireBuilder.Build()
}
