// Port of com/google/javascript/jscomp/RemoveUnusedCode.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	rucVarFlagEntirelyRemovable uint8 = 1 << iota
	rucVarFlagHasNonLocalOrNonLit
	rucVarFlagHasFuncOrClassLit
	rucVarFlagRequiresLocalLit
	rucVarFlagCanonical
)

type rucPinReason uint8

const (
	rucPinNone rucPinReason = iota
	rucPinNonLocalValueWithPropRef
	rucPinReferenceFound
	rucPinExportedClass
	rucPinCatchVariable
	rucPinEnhancedForVar
	rucPinForLoopInitSideEffects
	rucPinBaseClassSideEffects
	rucPinClassBodySideEffects
	rucPinParamInFnUsingArguments
	rucPinNotRemovingGlobals
	rucPinNotRemovingLocals
)

// RemoveUnusedCode_RealVarInfo tracks the removable code and other state
// related to variables we may be able to remove.
type RemoveUnusedCode_RealVarInfo struct {
	outer              *RemoveUnusedCode
	var_               *Var
	remHead            *RemoveUnusedCode_RemovableBase
	remTail            *RemoveUnusedCode_RemovableBase
	reasonNode         rhino.Node
	flags              uint8
	nonRemovableReason rucPinReason
}

// Java: RemoveUnusedCode.RealVarInfo(String)
func NewRemoveUnusedCode_RealVarInfo(outer *RemoveUnusedCode, var_ *Var) *RemoveUnusedCode_RealVarInfo {
	v := outer.varInfoSlab.New()
	v.outer = outer
	v.var_ = var_
	v.flags = rucVarFlagEntirelyRemovable
	return v
}

// Java: RemoveUnusedCode.RealVarInfo.addRemovable(Removable)
func (v *RemoveUnusedCode_RealVarInfo) AddRemovable(removable *RemoveUnusedCode_RemovableBase) {
	if v.flags&rucVarFlagCanonical != 0 {
		v.outer.ConsiderForIndependentRemoval(removable)
		return
	}
	if removable.IsVariableAssignment() {
		if localValue := removable.GetLocalAssignedValue(); !localValue.IsNil() {
			if localValue.IsFunction() || localValue.IsClass() {
				v.flags |= rucVarFlagHasFuncOrClassLit
			}
		} else if removable.flags&remFlagUninitVarDecl == 0 {
			v.flags |= rucVarFlagHasNonLocalOrNonLit
		}
	} else if removable.IsPrototypeAssignment() && !removable.IsAssignedValueLocal() {
		v.flags |= rucVarFlagHasNonLocalOrNonLit
	}
	if removable.PreventsRemovalOfVariableWithNonLocalValueOrPrototype() {
		v.flags |= rucVarFlagRequiresLocalLit
	}
	if v.flags&(rucVarFlagHasNonLocalOrNonLit|rucVarFlagRequiresLocalLit) == (rucVarFlagHasNonLocalOrNonLit | rucVarFlagRequiresLocalLit) {
		v.pin(rucPinNonLocalValueWithPropRef)
	}

	if v.flags&rucVarFlagEntirelyRemovable != 0 {
		removable.propNext = nil
		if v.remHead == nil {
			v.remHead = removable
			v.remTail = removable
		} else {
			v.remTail.propNext = removable
			v.remTail = removable
		}
	} else {
		v.outer.ConsiderForIndependentRemoval(removable)
	}
}

// Java: RemoveUnusedCode.RealVarInfo.isRemovable()
func (v *RemoveUnusedCode_RealVarInfo) IsRemovable() bool {
	return v.flags&rucVarFlagEntirelyRemovable != 0
}

func (v *RemoveUnusedCode_RealVarInfo) pin(reason rucPinReason) {
	v.pinAt(reason, rhino.NilNode)
}

func (v *RemoveUnusedCode_RealVarInfo) pinAt(reason rucPinReason, reasonNode rhino.Node) {
	if v.flags&rucVarFlagEntirelyRemovable == 0 {
		return
	}
	v.flags &^= rucVarFlagEntirelyRemovable
	v.nonRemovableReason = reason
	v.reasonNode = reasonNode
	if v.outer.unremovableLog.IsLogging() {
		v.outer.unremovableLog.LogString(v.formatUnremovableLog())
	}
	cur := v.remHead
	v.remHead = nil
	v.remTail = nil
	for cur != nil {
		next := cur.propNext
		cur.propNext = nil
		v.outer.ConsiderForIndependentRemoval(cur)
		cur = next
	}
}

func (v *RemoveUnusedCode_RealVarInfo) formatUnremovableLog() string {
	var reason string
	switch v.nonRemovableReason {
	case rucPinNonLocalValueWithPropRef:
		reason = "hasNonLocalOrNonLiteralValue && requiresLocalLiteralValueForRemoval"
	case rucPinReferenceFound:
		reason = "reference found: " + v.reasonNode.GetLocation()
	case rucPinExportedClass:
		reason = "exported class"
	case rucPinCatchVariable:
		reason = "catch variable"
	case rucPinEnhancedForVar:
		reason = "for-of or for-in loop variable"
	case rucPinForLoopInitSideEffects:
		reason = "for-loop variable initialization has side-effects"
	case rucPinBaseClassSideEffects:
		reason = "base class expression has side-effects"
	case rucPinClassBodySideEffects:
		reason = "class body has side-effects"
	case rucPinParamInFnUsingArguments:
		reason = "parameter in function using arguments"
	case rucPinNotRemovingGlobals:
		reason = "not removing globals"
	case rucPinNotRemovingLocals:
		reason = "not removing locals"
	}
	return v.var_.GetName() + ": " + reason
}

// Java: RemoveUnusedCode.RealVarInfo.setHasNonLocalOrNonLiteralValue()
func (v *RemoveUnusedCode_RealVarInfo) SetHasNonLocalOrNonLiteralValue() {
	v.flags |= rucVarFlagHasNonLocalOrNonLit
}

// Java: RemoveUnusedCode.RealVarInfo.hasFunctionOrClassLiteralValue()
func (v *RemoveUnusedCode_RealVarInfo) HasFunctionOrClassLiteralValue() bool {
	return v.flags&rucVarFlagHasFuncOrClassLit != 0
}

// Java: RemoveUnusedCode.RealVarInfo.removeAllRemovables()
func (v *RemoveUnusedCode_RealVarInfo) RemoveAllRemovables() {
	jutil.CheckState(v.flags&rucVarFlagEntirelyRemovable != 0)
	cur := v.remHead
	v.remHead = nil
	v.remTail = nil
	for cur != nil {
		next := cur.propNext
		cur.propNext = nil
		cur.Remove(v.outer)
		cur = next
	}
}

type rucPolyfillKind uint8

const (
	rucPolyfillGlobal rucPolyfillKind = iota
	rucPolyfillStatic
	rucPolyfillPrototype
)

// rucPolyfillInfo stores information about definitions and usages of polyfills.
//
// Java: RemoveUnusedCode.PolyfillInfo / GlobalPolyfillInfo / StaticPropertyPolyfillInfo / PrototypePropertyPolyfillInfo
type rucPolyfillInfo struct {
	outer       *RemoveUnusedCode
	removable   *RemoveUnusedCode_RemovableBase
	nextByKey   *rucPolyfillInfo
	owner       string
	key         string
	keyID       uint32
	kind        rucPolyfillKind
	isRemovable bool
}

// createPolyfillInfo makes a new rucPolyfillInfo, including the Polyfill
// Removable, parsing the name to classify global, static, or prototype polyfills.
//
// Java: RemoveUnusedCode.createPolyfillInfo(Node, Scope, String)
func (r *RemoveUnusedCode) createPolyfillInfo(call ast.NormalOrOptChainCall, scope *Scope, name string) *rucPolyfillInfo {
	jutil.CheckState(call.Node().GetParent().IsExprResult())
	if !r.guardedUsagesComputed {
		r.guardedUsagesComputed = true
		NewPolyfillUsageFinder(r.compiler, r.polyfillsFromTable).
			TraverseOnlyGuarded(r.root, r.StorePolyfill)
	}
	removable := r.newPolyfillRemovable(call, scope)
	info := r.polyfillSlab.New()
	lastDot := strings.LastIndexByte(name, '.')
	if lastDot < 0 {
		*info = rucPolyfillInfo{
			outer:       r,
			removable:   removable,
			key:         name,
			keyID:       rhino.RhinoStringPool_InternID(name),
			kind:        rucPolyfillGlobal,
			isRemovable: true,
		}
		return info
	}
	owner, prop := name[:lastDot], name[lastDot+1:]
	kind := rucPolyfillStatic
	if trimmed, ok := strings.CutSuffix(owner, ".prototype"); ok {
		owner = trimmed
		kind = rucPolyfillPrototype
	}
	*info = rucPolyfillInfo{
		outer:       r,
		removable:   removable,
		owner:       owner,
		key:         prop,
		keyID:       rhino.RhinoStringPool_InternID(prop),
		kind:        kind,
		isRemovable: true,
	}
	return info
}

func (p *rucPolyfillInfo) getName() string {
	switch p.kind {
	case rucPolyfillStatic:
		return p.owner + "." + p.key
	case rucPolyfillPrototype:
		return p.owner + ".prototype." + p.key
	default:
		return p.key
	}
}

// considerPossibleReference accepts a NAME or GETPROP node whose (property)
// string matches key and marks the polyfill non-removable if referenced.
//
// Java: RemoveUnusedCode.PolyfillInfo.considerPossibleReference(Node)
func (p *rucPolyfillInfo) considerPossibleReference(n rhino.Node) {
	if !p.isRemovable || (p.removable.flags&remFlagSubBit == 0 && p.outer.guardedUsages != nil && p.outer.guardedUsages.Contains(n)) {
		return
	}
	if NodeUtil_IsNormalOrOptChainGetProp(n) || (p.kind == rucPolyfillGlobal && n.IsName()) {
		p.isRemovable = false
		p.removable.ApplyContinuations(p.outer)
	}
}

// Java: RemoveUnusedCode.removeExpressionCompletely(Node)
func (r *RemoveUnusedCode) RemoveExpressionCompletely(expression rhino.Node) {
	jutil.CheckStateMsg1(!NodeUtil_IsExpressionResultUsed(expression), "%s", expression)
	parent := expression.GetParent()
	if parent.IsExprResult() {
		NodeUtil_DeleteNode(parent, r.compiler)
	} else if parent.IsComma() {
		// Expression is probably the first child of the comma,
		// but it could be the second if the entire comma expression value is unused.
		otherChild := expression.GetNext()
		if otherChild.IsNil() {
			otherChild = expression.GetPrevious()
		}
		r.ReplaceNodeWith(parent, otherChild.Detach())
	} else {
		// value isn't needed, but we need to keep the AST valid.
		r.ReplaceNodeWith(expression, rhino.IR_Number(0).Srcref(expression))
	}
}

// Java: RemoveUnusedCode.replaceNodeWith(Node, Node)
func (r *RemoveUnusedCode) ReplaceNodeWith(n rhino.Node, replacement rhino.Node) {
	r.compiler.ReportChangeToEnclosingScope(n)
	n.ReplaceWith(replacement)
	NodeUtil_MarkFunctionsDeleted(n, r.compiler)
}
