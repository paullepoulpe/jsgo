// Port of com/google/javascript/jscomp/RemoveUnusedCode.java @ closure-compiler v20261028

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
	"github.com/paullepoulpe/jsgo/internal/strpool"
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
	removables         rucRemovableQueue
	reasonNode         rhino.Node
	flags              uint8
	nonRemovableReason rucPinReason
}

// Java: RemoveUnusedCode.RealVarInfo(String)
func NewRemoveUnusedCode_RealVarInfo(outer *RemoveUnusedCode, var_ *Var) *RemoveUnusedCode_RealVarInfo {
	v := outer.varInfoSlab.Alloc()
	v.outer = outer
	v.var_ = var_
	v.flags = rucVarFlagEntirelyRemovable
	return v
}

// Java: RemoveUnusedCode.RealVarInfo.addRemovable(Removable)
func (v *RemoveUnusedCode_RealVarInfo) AddRemovable(removable *rucRemovable) {
	if v.flags&rucVarFlagCanonical != 0 {
		v.outer.ConsiderForIndependentRemoval(removable)
		return
	}
	if v.outer.recMode != rucRecNone {
		v.outer.onVarMutated(v.var_)
	}
	if removable.IsVariableAssignment() {
		if localValue := removable.GetLocalAssignedValue(); !localValue.IsNil() {
			_, isFn := ast.AsFunction(localValue)
			_, isCls := ast.AsClass(localValue)
			if isFn || isCls {
				v.flags |= rucVarFlagHasFuncOrClassLit
			}
		} else if !removable.isUninitVarDecl() {
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
		v.removables.push(removable)
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
	if v.flags&rucVarFlagCanonical == 0 && v.outer.recMode != rucRecNone &&
		reason != rucPinNotRemovingGlobals && reason != rucPinNotRemovingLocals {
		v.outer.onVarPinned(v.var_)
	}
	if v.flags&rucVarFlagEntirelyRemovable == 0 {
		return
	}
	v.flags &^= rucVarFlagEntirelyRemovable
	v.nonRemovableReason = reason
	v.reasonNode = reasonNode
	if v.outer.unremovableLog.IsLogging() {
		v.outer.unremovableLog.LogString(v.formatUnremovableLog())
	}
	if v.outer.activeGlobalRec != nil && v.var_.IsGlobal() {
		prev := v.outer.activeGlobalRec
		v.outer.activeGlobalRec = nil
		v.removables.drain(v.outer.ConsiderForIndependentRemoval)
		v.outer.activeGlobalRec = prev
	} else {
		v.removables.drain(v.outer.ConsiderForIndependentRemoval)
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
	if v.flags&rucVarFlagCanonical == 0 {
		if v.flags&rucVarFlagEntirelyRemovable != 0 && v.outer.recMode != rucRecNone {
			v.outer.onVarMutated(v.var_)
		}
		v.flags |= rucVarFlagHasNonLocalOrNonLit
	}
}

// Java: RemoveUnusedCode.RealVarInfo.hasFunctionOrClassLiteralValue()
func (v *RemoveUnusedCode_RealVarInfo) HasFunctionOrClassLiteralValue() bool {
	return v.flags&rucVarFlagHasFuncOrClassLit != 0
}

// Java: RemoveUnusedCode.RealVarInfo.removeAllRemovables()
func (v *RemoveUnusedCode_RealVarInfo) RemoveAllRemovables() {
	jutil.CheckState(v.flags&rucVarFlagEntirelyRemovable != 0)
	v.removables.drain(func(cur *rucRemovable) {
		cur.Remove(v.outer)
	})
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
	removable   *rucRemovable
	nextByKey   *rucPolyfillInfo
	owner       string
	key         string
	keyID       strpool.ID
	kind        rucPolyfillKind
	isRemovable bool
}

// createPolyfillInfo makes a new rucPolyfillInfo, including the Polyfill
// Removable, parsing the name to classify global, static, or prototype polyfills.
//
// Java: RemoveUnusedCode.createPolyfillInfo(Node, Scope, String)
func (r *RemoveUnusedCode) createPolyfillInfo(call ast.NormalOrOptChainCall, scope *Scope, name string) *rucPolyfillInfo {
	_, isExprStmt := ast.AsExprStmt(call.Node().GetParent())
	jutil.CheckState(isExprStmt)
	removable := r.newPolyfillRemovable(call, scope)
	info := r.polyfillSlab.Alloc()
	lastDot := strings.LastIndexByte(name, '.')
	if lastDot < 0 {
		*info = rucPolyfillInfo{
			outer:       r,
			removable:   removable,
			key:         name,
			keyID:       strpool.Intern(name),
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
		keyID:       strpool.Intern(prop),
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
	if !p.isRemovable || (!p.removable.isPolyfillPatch() && p.outer.isGuardedUsage(n)) {
		return
	}
	_, isGetProp := ast.AsNormalOrOptChainGetProp(n)
	_, isName := ast.AsName(n)
	if isGetProp || (p.kind == rucPolyfillGlobal && isName) {
		if p.kind == rucPolyfillPrototype && p.outer.disambiguateProperties &&
			p.outer.polyfillUsageFinder.IsDisjointFromIteratorNode(n, p.owner) {
			return
		}
		p.isRemovable = false
		p.removable.ApplyContinuations(p.outer)
	}
}

// Java: RemoveUnusedCode.removeExpressionCompletely(Node)
func (r *RemoveUnusedCode) RemoveExpressionCompletely(expression rhino.Node) {
	jutil.CheckStateMsg1(!NodeUtil_IsExpressionResultUsed(expression), "%s", expression)
	parent := expression.GetParent()
	if _, isExprStmt := ast.AsExprStmt(parent); isExprStmt {
		NodeUtil_DeleteNode(parent, r.compiler)
	} else if commaBin, isComma := ast.AsBinaryOf(parent, rhino.Token_COMMA); isComma {
		// Expression is probably the first child of the comma,
		// but it could be the second if the entire comma expression value is unused.
		var otherChild rhino.Node
		if commaBin.Left().Is(expression) {
			otherChild = commaBin.Right().Detach().Node()
		} else {
			otherChild = commaBin.Left().Detach().Node()
		}
		r.ReplaceNodeWith(parent, otherChild)
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
