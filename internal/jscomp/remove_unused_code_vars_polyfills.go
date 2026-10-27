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

// RemoveUnusedCode_RealVarInfo tracks the removable code and other state
// related to variables we may be able to remove.
type RemoveUnusedCode_RealVarInfo struct {
	outer   *RemoveUnusedCode
	var_    *Var
	remHead *RemoveUnusedCode_RemovableBase
	remTail *RemoveUnusedCode_RemovableBase
	flags   uint8
}

// Java: RemoveUnusedCode.RealVarInfo(String)
func NewRemoveUnusedCode_RealVarInfo(outer *RemoveUnusedCode, var_ *Var) *RemoveUnusedCode_RealVarInfo {
	v := outer.varInfoSlab.New()
	*v = RemoveUnusedCode_RealVarInfo{outer: outer, var_: var_, flags: rucVarFlagEntirelyRemovable}
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
		v.SetIsExplicitlyNotRemovable(
			func() string { return "hasNonLocalOrNonLiteralValue && requiresLocalLiteralValueForRemoval" })
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

// Java: RemoveUnusedCode.RealVarInfo.setIsExplicitlyNotRemovable(Supplier<String>)
func (v *RemoveUnusedCode_RealVarInfo) SetIsExplicitlyNotRemovable(reasonSupplier func() string) {
	if v.flags&rucVarFlagEntirelyRemovable != 0 {
		v.flags &^= rucVarFlagEntirelyRemovable
		if v.outer.unremovableLog.IsLogging() {
			v.outer.unremovableLog.LogString(v.var_.GetName() + ": " + reasonSupplier())
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

const RemoveUnusedCode_DOT_PROTOTYPE = ".prototype"

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
	owner       string
	key         string
	kind        rucPolyfillKind
	isRemovable bool
}

// createPolyfillInfo makes a new rucPolyfillInfo, including the Polyfill
// Removable, parsing the name to classify global, static, or prototype polyfills.
//
// Java: RemoveUnusedCode.createPolyfillInfo(Node, Scope, String)
func (r *RemoveUnusedCode) createPolyfillInfo(call rhino.Node, scope *Scope, name string) *rucPolyfillInfo {
	jutil.CheckState(call.GetParent().IsExprResult())
	if !r.guardedUsagesComputed {
		r.guardedUsagesComputed = true
		NewPolyfillUsageFinder(r.compiler, r.polyfillsFromTable).
			TraverseOnlyGuarded(r.root, r.StorePolyfill)
	}
	// Make the removable and polyfill info. Add continuations for all arguments.
	builder := NewRemoveUnusedCode_RemovableBuilder(r)
	for arg := range ast.MustNormalOrOptChainCall(call).Args() {
		builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, arg.Node(), scope))
	}
	removable := builder.BuildPolyfill(call.GetParent())
	lastDot := jutil.StringLastIndexOf(name, ".")
	if lastDot < 0 {
		return &rucPolyfillInfo{outer: r, removable: removable, key: name, kind: rucPolyfillGlobal, isRemovable: true}
	}
	owner := jutil.StringSubstring(name, 0, lastDot)
	prop := jutil.StringSubstringFrom(name, lastDot+1)
	if strings.HasSuffix(owner, RemoveUnusedCode_DOT_PROTOTYPE) {
		owner = jutil.StringSubstring(owner, 0, jutil.StringLength(owner)-jutil.StringLength(RemoveUnusedCode_DOT_PROTOTYPE))
		return &rucPolyfillInfo{outer: r, removable: removable, owner: owner, key: prop, kind: rucPolyfillPrototype, isRemovable: true}
	}
	return &rucPolyfillInfo{outer: r, removable: removable, owner: owner, key: prop, kind: rucPolyfillStatic, isRemovable: true}
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

// Java: RemoveUnusedCode.VanillaForNameDeclaration(RemovableBuilder, Node)
func NewRemoveUnusedCode_VanillaForNameDeclaration(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, nameNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	d := outer.removables.New()
	d.initFromBuilder(remKindVanillaForNameDecl, nameNode, builder)
	if _, ok := ast.MustName(nameNode).Init(); !ok {
		d.flags |= remFlagUninitVarDecl
	}
	return d
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
