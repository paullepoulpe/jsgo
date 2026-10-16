// Port of com/google/javascript/jscomp/ReferenceCollection.java @ closure-compiler v20260415

/*
 * Copyright 2017 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ReferenceCollection is a collection of references. Can be subclassed to
// apply checks or store additional state when adding.
type ReferenceCollection struct {
	references *jutil.ArrayList[*Reference]
}

// Java: ReferenceCollection()
func NewReferenceCollection() *ReferenceCollection {
	return &ReferenceCollection{references: jutil.NewArrayList[*Reference]()}
}

// All is the Iterable's iterator().
//
// Java: ReferenceCollection.iterator()
func (c *ReferenceCollection) All() iter.Seq[*Reference] {
	return c.references.All()
}

// Java: ReferenceCollection.add(Reference)
func (c *ReferenceCollection) Add(reference *Reference) {
	c.references.Add(reference)
}

// IsWellDefined determines if the variable for this reference collection is
// "well-defined." A variable is well-defined if we can prove at compile-time
// that it's assigned a value before it's used.
//
// Notice that if this function returns false, this doesn't imply that the
// variable is used before it's assigned. It just means that we don't have
// enough information to make a definitive judgment.
//
// Java: ReferenceCollection.isWellDefined()
func (c *ReferenceCollection) IsWellDefined() bool {
	size := c.references.Size()
	if size == 0 {
		return false
	}

	// If this is a declaration that does not instantiate the variable,
	// it's not well-defined.
	init := c.GetInitializingReference()
	if init == nil {
		return false
	}

	jutil.CheckState(c.references.Get(0).IsDeclaration())
	initBlock := init.GetBasicBlock()
	for i := 1; i < size; i++ {
		if !initBlock.ProvablyExecutesBefore(c.references.Get(i).GetBasicBlock()) {
			return false
		}
	}

	return true
}

// IsEscaped reports whether the variable is escaped into an inner function.
//
// Java: ReferenceCollection.isEscaped()
func (c *ReferenceCollection) IsEscaped() bool {
	var hoistScope *Scope
	for ref := range c.references.All() {
		if hoistScope == nil {
			hoistScope = ref.GetScope().GetClosestHoistScope()
		} else if hoistScope != ref.GetScope().GetClosestHoistScope() {
			return true
		}
	}
	return false
}

// IsInitializingDeclarationAt: index is the index into the references array
// to look for an assigning declaration.
//
// This is either the declaration if a value is assigned (such as "var a =
// 2", "function a()...", "... catch (a)...").
//
// Java: ReferenceCollection.isInitializingDeclarationAt(int)
func (c *ReferenceCollection) IsInitializingDeclarationAt(index int) bool {
	maybeInit := c.references.Get(index)
	if maybeInit.IsInitializingDeclaration() {
		// This is a declaration that represents the initial value.
		// Specifically, var declarations without assignments such as "var a;"
		// are not.
		return true
	}
	return false
}

// IsInitializingAssignmentAt: index is the index into the references array
// to look for an initialized assignment reference. That is, an assignment
// immediately follow a variable declaration that itself does not initialize
// the variable.
//
// Java: ReferenceCollection.isInitializingAssignmentAt(int)
func (c *ReferenceCollection) IsInitializingAssignmentAt(index int) bool {
	if index < c.references.Size() && index > 0 {
		maybeDecl := c.references.Get(index - 1)
		if maybeDecl.IsVarDeclaration() || maybeDecl.IsLetDeclaration() {
			jutil.CheckState(!maybeDecl.IsInitializingDeclaration())
			maybeInit := c.references.Get(index)
			if maybeInit.IsSimpleAssignmentToName() {
				return true
			}
		}
	}
	return false
}

// GetInitializingReference returns the reference that provides the value
// for the variable at the time of the first read, if known, otherwise null.
//
// This is either the variable declaration ("var a = ...") or first reference
// following the declaration if it is an assignment in the same JSChunk.
//
// Java: ReferenceCollection.getInitializingReference()
func (c *ReferenceCollection) GetInitializingReference() *Reference {
	if c.IsInitializingDeclarationAt(0) {
		return c.references.Get(0)
	} else if c.IsInitializingAssignmentAt(1) &&
		// Check that these references are in the same file as a way of ensuring they are in the
		// same chunk: checking for the same chunk would also be reasonable, but is slightly more
		// complicated and as of time of writing, didn't seem to help code size in practice.
		c.AreInSameFile(0, 1) {
		return c.references.Get(1)
	}
	return nil
}

// Java: ReferenceCollection.areInSameFile(int, int)
func (c *ReferenceCollection) AreInSameFile(ref0 int, ref1 int) bool {
	return jutil.ObjectsEquals(c.references.Get(ref0).GetInputId(), c.references.Get(ref1).GetInputId())
}

// GetInitializingReferenceForConstants: constants are allowed to be defined
// after their first use.
//
// Java: ReferenceCollection.getInitializingReferenceForConstants()
func (c *ReferenceCollection) GetInitializingReferenceForConstants() *Reference {
	size := c.references.Size()
	for i := 0; i < size; i++ {
		if c.IsInitializingDeclarationAt(i) || c.IsInitializingAssignmentAt(i) {
			return c.references.Get(i)
		}
	}
	return nil
}

// IsAssignedOnceInLifetime returns whether the variable is only assigned a
// value once for its lifetime.
//
// Java: ReferenceCollection.isAssignedOnceInLifetime()
func (c *ReferenceCollection) IsAssignedOnceInLifetime() bool {
	ref := c.GetOneAndOnlyAssignment()
	if ref == nil {
		return false
	}

	// Make sure this assignment is not in a loop or an enclosing function.
	for block := ref.GetBasicBlock(); block != nil; block = block.GetParent() {
		if block.IsFunction() {
			if ref.getSymbolV().getScopeS().GetClosestHoistScope() !=
				ref.GetScope().GetClosestHoistScope() {
				return false
			}
			break
		}

		if block.GetRoot() == ref.getSymbolV().getScopeS().GetRootNode() {
			break
		}

		if block.IsLoop() {
			return false
		}
	}

	return true
}

// IsReferencedWeakly returns whether the variable is ever referenced weakly.
//
// Java: ReferenceCollection.isReferencedWeakly()
func (c *ReferenceCollection) IsReferencedWeakly() bool {
	for ref := range c.references.All() {
		if NodeUtil_IsGoogWeakUsageCall(ref.GetParent()) {
			return true
		}
	}
	return false
}

// GetOneAndOnlyAssignment returns the one and only assignment. Returns null
// if the number of assignments is not exactly one.
//
// Java: ReferenceCollection.getOneAndOnlyAssignment()
func (c *ReferenceCollection) GetOneAndOnlyAssignment() *Reference {
	var assignment *Reference
	size := c.references.Size()
	for i := 0; i < size; i++ {
		ref := c.references.Get(i)
		if ref.IsLvalue() || ref.IsInitializingDeclaration() {
			if assignment == nil {
				assignment = ref
			} else {
				return nil
			}
		}
	}
	return assignment
}

// IsNeverAssigned returns whether the variable is never assigned a value.
//
// Java: ReferenceCollection.isNeverAssigned()
func (c *ReferenceCollection) IsNeverAssigned() bool {
	size := c.references.Size()
	for i := 0; i < size; i++ {
		ref := c.references.Get(i)
		if ref.IsLvalue() || ref.IsInitializingDeclaration() {
			return false
		}
	}
	return true
}

// Java: ReferenceCollection.firstReferenceIsAssigningDeclaration()
func (c *ReferenceCollection) FirstReferenceIsAssigningDeclaration() bool {
	size := c.references.Size()
	return size > 0 && c.references.Get(0).IsInitializingDeclaration()
}

// Java: ReferenceCollection.toString()
func (c *ReferenceCollection) String() string {
	return jutil.MoreObjects_ToStringHelper("ReferenceCollection").
		Add("initRef", c.GetInitializingReference()).
		Add("references", c.references).
		AddBool("wellDefined", c.IsWellDefined()).
		AddBool("assignedOnce", c.IsAssignedOnceInLifetime()).
		String()
}
