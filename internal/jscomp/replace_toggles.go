// Port of com/google/javascript/jscomp/ReplaceToggles.java @ closure-compiler v20261021

/*
 * Copyright 2023 The Closure Compiler Authors.
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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReplaceToggles replaces calls to id generators with ids.
//
// Use this to get unique and short ids.
type ReplaceToggles struct {
	compiler   AbstractCompiler
	astFactory *AstFactory

	// nil is Java null.
	ordinalMapping map[string]int
}

// NOTE: These diagnostics are only checked in Stage 2 (optimization), since none of this
// code should ever be hand-written: the `CLOSURE_TOGGLE_ORDINALS` bootstrap and the calls to
// `goog.readToggleInternalDoNotCallDirectly` are all generated automatically by the build
// system, so it's unexpected that anyone should ever run into these diagnostics when doing
// ordinary development.

// Java: ReplaceToggles.INVALID_TOGGLE_PARAMETER
var ReplaceToggles_INVALID_TOGGLE_PARAMETER = DiagnosticType_Error(
	"JSC_INVALID_TOGGLE_PARAMETER",
	"goog.readToggleInternalDoNotCallDirectly must be called with a string literal.")

// Java: ReplaceToggles.INVALID_ORDINAL_MAPPING
var ReplaceToggles_INVALID_ORDINAL_MAPPING = DiagnosticType_Error(
	"JSC_INVALID_ORDINAL_MAPPING",
	"CLOSURE_TOGGLE_ORDINALS must be initialized with an object literal mapping strings to"+
		" booleans or unique whole numbers: {0}")

// NOTE: These values are chosen as negative integers because actual toggle ordinals must always
// be non-negative (at least zero).  Any negative integers would do to distinguish them from real
// toggle ordinals, but -1 and -2 are the simplest.
const (
	// Java: ReplaceToggles.TRUE_VALUE
	ReplaceToggles_TRUE_VALUE = -2
	// Java: ReplaceToggles.FALSE_VALUE
	ReplaceToggles_FALSE_VALUE = -1
)

// Java: ReplaceToggles(AbstractCompiler)
func NewReplaceToggles(compiler AbstractCompiler) *ReplaceToggles {
	r := &ReplaceToggles{ordinalMapping: nil}
	r.compiler = compiler
	r.astFactory = compiler.CreateAstFactory()
	return r
}

// Java: ReplaceToggles.process(Node, Node)
func (r *ReplaceToggles) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, NewReplaceToggles_Traversal(r))
}

// Java: ReplaceToggles.ORDINAL_VAR_NAME
const ReplaceToggles_ORDINAL_VAR_NAME = "CLOSURE_TOGGLE_ORDINALS"

// Java: ReplaceToggles.readToggleFunctionName
var ReplaceToggles_readToggleFunctionName = rhino.QualifiedName_Of("goog.readToggleInternalDoNotCallDirectly")

// ReplaceToggles_Traversal is the inner class Traversal.
type ReplaceToggles_Traversal struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	this_ *ReplaceToggles // the enclosing instance
}

// Java: Traversal()
func NewReplaceToggles_Traversal(outer *ReplaceToggles) *ReplaceToggles_Traversal {
	return &ReplaceToggles_Traversal{this_: outer}
}

// Java: Traversal.visit(NodeTraversal, Node, Node)
func (v *ReplaceToggles_Traversal) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	r := v.this_
	// Look for `var CLOSURE_TOGGLE_ORDINALS = {...};`, record the mapping, and then delete the
	// declaration from the AST since we should no longer need it (and the optimizer won't
	// delete it in case the global assignment was an intended side effect).
	if NodeUtil_IsNameDeclaration(n) && n.GetFirstChild().MatchesName(ReplaceToggles_ORDINAL_VAR_NAME) {
		rhs := n.GetFirstFirstChild()

		if rhs.IsNil() && n.IsVar() {
			// An empty var is fine; it should get deleted later.
			return
		} else if !rhs.IsObjectLit() {
			r.compiler.Report(JSError_MakeWithNode(n, ReplaceToggles_INVALID_ORDINAL_MAPPING, "not an object literal"))
			return
		} else if r.ordinalMapping != nil {
			r.compiler.Report(JSError_MakeWithNode(n, ReplaceToggles_INVALID_ORDINAL_MAPPING, "multiple initialized copies"))
			return
		}

		mapping := make(map[string]int)
		ordinals := make(map[int]bool)
		for c := rhs.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if !c.IsStringKey() && !c.IsStringLit() {
				r.compiler.Report(JSError_MakeWithNode(c, ReplaceToggles_INVALID_ORDINAL_MAPPING, "non-string key"))
				return
			}
			key := c.GetString()
			if _, dup := mapping[key]; dup {
				r.compiler.Report(JSError_MakeWithNode(c, ReplaceToggles_INVALID_ORDINAL_MAPPING, "duplicate key: "+key))
				return
			}
			child := c.GetFirstChild()
			doubleValue := NodeUtil_GetNumberValue(child)
			intValue := -1
			if doubleValue != nil {
				intValue = int(jutil.D2I(*doubleValue))
			}
			if child.IsTrue() || child.IsFalse() {
				if child.IsTrue() {
					intValue = ReplaceToggles_TRUE_VALUE
				} else {
					intValue = ReplaceToggles_FALSE_VALUE
				}
			} else if !child.IsNumber() || intValue < 0 || float64(intValue) != *doubleValue {
				r.compiler.Report(
					JSError_MakeWithNode(
						c, ReplaceToggles_INVALID_ORDINAL_MAPPING, "value not a boolean or whole number literal"))
				return
			} else if ordinals[intValue] {
				r.compiler.Report(
					JSError_MakeWithNode(c, ReplaceToggles_INVALID_ORDINAL_MAPPING, "duplicate ordinal: "+strconv.Itoa(intValue)))
				return
			}
			mapping[key] = intValue
			ordinals[intValue] = true
		}
		r.ordinalMapping = mapping

		// NOTE: We do not support a simple assignment without `var` since reassignment or later
		// augmentation (i.e. `CLOSURE_TOGGLE_ORDINALS['foo'] = true`) is not allowed.
		return
	}

	if !n.IsCall() {
		return
	}
	qname := NodeUtil_GetCallTargetResolvingIndirectCalls(n)
	if !ReplaceToggles_readToggleFunctionName.Matches(qname) {
		return
	}

	arg := n.GetSecondChild()
	if arg.IsNil() || !arg.IsStringLit() || !n.HasTwoChildren() {
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceToggles_INVALID_TOGGLE_PARAMETER))
		return
	}

	ordinal, hasOrdinal := r.ordinalMapping[arg.GetString()]
	if !hasOrdinal || ordinal < 0 {
		// No ordinals given: hard-code `true` if explicitly set as true, or `false` otherwise.
		n.ReplaceWith(
			r.astFactory.
				CreateBoolean(hasOrdinal && ordinal == ReplaceToggles_TRUE_VALUE).
				SrcrefTreeIfMissing(n))
		t.ReportCodeChange()
		return
	}

	// Replace with a lookup into the data structure
	// Note: the choice of 30 here is per spec, since the bootstrap must be written in agreement
	// with this convention.  30 was chosen over 32 to ensure all numbers are SMI and will fit
	// neatly in 4 bytes, whereas larger ints require more space in the VM.
	index := ordinal / 30
	bit := ordinal % 30
	getElem := r.astFactory.CreateGetElem(
		r.astFactory.CreateQNameWithUnknownType("goog.TOGGLES_"),
		r.astFactory.CreateNumber(float64(index)))
	// There are two different ways to write the bitwiseAnd:
	//  1. `getElem & mask`
	//  2. `getElem >> bit & 1` (note: sign extension is irrelevant here, and `>>` is shorter)
	// where x is the `goog.TOGGLES_[index]` lookup, and `mask` is `1 << bit`.
	// When `bit < 14`, option 1 is shorter.  When `bit > 16`, option 2 is shorter.
	// We arbitrarily prefer option 2 in the break-even range 14..16.
	var bitAnd rhino.Node
	if bit < 14 {
		bitAnd = r.astFactory.CreateBitwiseAnd(getElem, r.astFactory.CreateNumber(float64(int32(1)<<bit)))
	} else {
		bitAnd = r.astFactory.CreateBitwiseAnd(
			r.astFactory.CreateRightShift(getElem, r.astFactory.CreateNumber(float64(bit))),
			r.astFactory.CreateNumber(1))
	}
	n.ReplaceWith(r.astFactory.CreateNot(r.astFactory.CreateNot(bitAnd)).SrcrefTreeIfMissing(n))
	t.ReportCodeChange()
}
