// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Norris Boyd
 *   Roger Lawrence
 *   Mike McCabe
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

// AST type check methods

// Java: Node.isAdd()
func (n Node) IsAdd() bool {
	return n.rec().token == Token_ADD
}

// Java: Node.isSub()
func (n Node) IsSub() bool {
	return n.rec().token == Token_SUB
}

// Java: Node.isAnd()
func (n Node) IsAnd() bool {
	return n.rec().token == Token_AND
}

// Java: Node.isAssignAnd()
func (n Node) IsAssignAnd() bool {
	return n.rec().token == Token_ASSIGN_AND
}

// Java: Node.isArrayLit()
func (n Node) IsArrayLit() bool {
	return n.rec().token == Token_ARRAYLIT
}

// Java: Node.isArrayPattern()
func (n Node) IsArrayPattern() bool {
	return n.rec().token == Token_ARRAY_PATTERN
}

// Java: Node.isAssign()
func (n Node) IsAssign() bool {
	return n.rec().token == Token_ASSIGN
}

// Java: Node.isAssignAdd()
func (n Node) IsAssignAdd() bool {
	return n.rec().token == Token_ASSIGN_ADD
}

// Java: Node.isNormalBlock()
func (n Node) IsNormalBlock() bool {
	return n.IsBlock()
}

// Java: Node.isBlock()
func (n Node) IsBlock() bool {
	return n.rec().token == Token_BLOCK
}

// Java: Node.isRoot()
func (n Node) IsRoot() bool {
	return n.rec().token == Token_ROOT
}

// Java: Node.isAwait()
func (n Node) IsAwait() bool {
	return n.rec().token == Token_AWAIT
}

// Java: Node.isBigInt()
func (n Node) IsBigInt() bool {
	return n.rec().token == Token_BIGINT
}

// Java: Node.isBitNot()
func (n Node) IsBitNot() bool {
	return n.rec().token == Token_BITNOT
}

// Java: Node.isBreak()
func (n Node) IsBreak() bool {
	return n.rec().token == Token_BREAK
}

// Java: Node.isCall()
func (n Node) IsCall() bool {
	return n.rec().token == Token_CALL
}

// Java: Node.isCase()
func (n Node) IsCase() bool {
	return n.rec().token == Token_CASE
}

// Java: Node.isCast()
func (n Node) IsCast() bool {
	return n.rec().token == Token_CAST
}

// Java: Node.isCatch()
func (n Node) IsCatch() bool {
	return n.rec().token == Token_CATCH
}

// Java: Node.isClass()
func (n Node) IsClass() bool {
	return n.rec().token == Token_CLASS
}

// Java: Node.isClassMembers()
func (n Node) IsClassMembers() bool {
	return n.rec().token == Token_CLASS_MEMBERS
}

// Java: Node.isComma()
func (n Node) IsComma() bool {
	return n.rec().token == Token_COMMA
}

// Java: Node.isComputedProp()
func (n Node) IsComputedProp() bool {
	return n.rec().token == Token_COMPUTED_PROP
}

// Java: Node.isContinue()
func (n Node) IsContinue() bool {
	return n.rec().token == Token_CONTINUE
}

// Java: Node.isConst()
func (n Node) IsConst() bool {
	return n.rec().token == Token_CONST
}

// Java: Node.isDebugger()
func (n Node) IsDebugger() bool {
	return n.rec().token == Token_DEBUGGER
}

// Java: Node.isDec()
func (n Node) IsDec() bool {
	return n.rec().token == Token_DEC
}

// Java: Node.isDefaultCase()
func (n Node) IsDefaultCase() bool {
	return n.rec().token == Token_DEFAULT_CASE
}

// Java: Node.isDefaultValue()
func (n Node) IsDefaultValue() bool {
	return n.rec().token == Token_DEFAULT_VALUE
}

// Java: Node.isDelProp()
func (n Node) IsDelProp() bool {
	return n.rec().token == Token_DELPROP
}

// Java: Node.isDestructuringLhs()
func (n Node) IsDestructuringLhs() bool {
	return n.rec().token == Token_DESTRUCTURING_LHS
}

// Java: Node.isDestructuringPattern()
func (n Node) IsDestructuringPattern() bool {
	return n.IsObjectPattern() || n.IsArrayPattern()
}

// Java: Node.isDo()
func (n Node) IsDo() bool {
	return n.rec().token == Token_DO
}

// Java: Node.isEmpty()
func (n Node) IsEmpty() bool {
	return n.rec().token == Token_EMPTY
}

// Java: Node.isExponent()
func (n Node) IsExponent() bool {
	return n.rec().token == Token_EXPONENT
}

// Java: Node.isAssignExponent()
func (n Node) IsAssignExponent() bool {
	return n.rec().token == Token_ASSIGN_EXPONENT
}

// Java: Node.isExport()
func (n Node) IsExport() bool {
	return n.rec().token == Token_EXPORT
}

// Java: Node.isExportSpec()
func (n Node) IsExportSpec() bool {
	return n.rec().token == Token_EXPORT_SPEC
}

// Java: Node.isExportSpecs()
func (n Node) IsExportSpecs() bool {
	return n.rec().token == Token_EXPORT_SPECS
}

// Java: Node.isExprResult()
func (n Node) IsExprResult() bool {
	return n.rec().token == Token_EXPR_RESULT
}

// Java: Node.isFalse()
func (n Node) IsFalse() bool {
	return n.rec().token == Token_FALSE
}

// Java: Node.isVanillaFor()
func (n Node) IsVanillaFor() bool {
	return n.rec().token == Token_FOR
}

// Java: Node.isForIn()
func (n Node) IsForIn() bool {
	return n.rec().token == Token_FOR_IN
}

// Java: Node.isForOf()
func (n Node) IsForOf() bool {
	return n.rec().token == Token_FOR_OF
}

// Java: Node.isForAwaitOf()
func (n Node) IsForAwaitOf() bool {
	return n.rec().token == Token_FOR_AWAIT_OF
}

// Java: Node.isFunction()
func (n Node) IsFunction() bool {
	return n.rec().token == Token_FUNCTION
}

// Java: Node.isGetterDef()
func (n Node) IsGetterDef() bool {
	return n.rec().token == Token_GETTER_DEF
}

// Java: Node.isGetElem()
func (n Node) IsGetElem() bool {
	return n.rec().token == Token_GETELEM
}

// Java: Node.isGetProp()
func (n Node) IsGetProp() bool {
	return n.rec().token == Token_GETPROP
}

// Java: Node.isHook()
func (n Node) IsHook() bool {
	return n.rec().token == Token_HOOK
}

// Java: Node.isIf()
func (n Node) IsIf() bool {
	return n.rec().token == Token_IF
}

// Java: Node.isImport()
func (n Node) IsImport() bool {
	return n.rec().token == Token_IMPORT
}

// Java: Node.isImportMeta()
func (n Node) IsImportMeta() bool {
	return n.rec().token == Token_IMPORT_META
}

// Java: Node.isImportStar()
func (n Node) IsImportStar() bool {
	return n.rec().token == Token_IMPORT_STAR
}

// Java: Node.isImportSpec()
func (n Node) IsImportSpec() bool {
	return n.rec().token == Token_IMPORT_SPEC
}

// Java: Node.isImportSpecs()
func (n Node) IsImportSpecs() bool {
	return n.rec().token == Token_IMPORT_SPECS
}

// Java: Node.isIn()
func (n Node) IsIn() bool {
	return n.rec().token == Token_IN
}

// Java: Node.isInc()
func (n Node) IsInc() bool {
	return n.rec().token == Token_INC
}

// Java: Node.isInstanceOf()
func (n Node) IsInstanceOf() bool {
	return n.rec().token == Token_INSTANCEOF
}

// Java: Node.isInterface()
func (n Node) IsInterface() bool {
	return n.rec().token == Token_INTERFACE
}

// Java: Node.isInterfaceMembers()
func (n Node) IsInterfaceMembers() bool {
	return n.rec().token == Token_INTERFACE_MEMBERS
}

// Java: Node.isRecordType()
func (n Node) IsRecordType() bool {
	return n.rec().token == Token_RECORD_TYPE
}

// Java: Node.isCallSignature()
func (n Node) IsCallSignature() bool {
	return n.rec().token == Token_CALL_SIGNATURE
}

// Java: Node.isIndexSignature()
func (n Node) IsIndexSignature() bool {
	return n.rec().token == Token_INDEX_SIGNATURE
}

// Java: Node.isLabel()
func (n Node) IsLabel() bool {
	return n.rec().token == Token_LABEL
}

// Java: Node.isLabelName()
func (n Node) IsLabelName() bool {
	return n.rec().token == Token_LABEL_NAME
}

// Java: Node.isLet()
func (n Node) IsLet() bool {
	return n.rec().token == Token_LET
}

// Java: Node.isMemberFunctionDef()
func (n Node) IsMemberFunctionDef() bool {
	return n.rec().token == Token_MEMBER_FUNCTION_DEF
}

// Java: Node.isMemberVariableDef()
func (n Node) IsMemberVariableDef() bool {
	return n.rec().token == Token_MEMBER_VARIABLE_DEF
}

// Java: Node.isMemberFieldDef()
func (n Node) IsMemberFieldDef() bool {
	return n.rec().token == Token_MEMBER_FIELD_DEF
}

// Java: Node.isComputedFieldDef()
func (n Node) IsComputedFieldDef() bool {
	return n.rec().token == Token_COMPUTED_FIELD_DEF
}

// Java: Node.isModuleBody()
func (n Node) IsModuleBody() bool {
	return n.rec().token == Token_MODULE_BODY
}

// Java: Node.isName()
func (n Node) IsName() bool {
	return n.rec().token == Token_NAME
}

// Java: Node.isNE()
func (n Node) IsNE() bool {
	return n.rec().token == Token_NE
}

// Java: Node.isSHNE()
func (n Node) IsSHNE() bool {
	return n.rec().token == Token_SHNE
}

// Java: Node.isEQ()
func (n Node) IsEQ() bool {
	return n.rec().token == Token_EQ
}

// Java: Node.isSHEQ()
func (n Node) IsSHEQ() bool {
	return n.rec().token == Token_SHEQ
}

// Java: Node.isNeg()
func (n Node) IsNeg() bool {
	return n.rec().token == Token_NEG
}

// Java: Node.isNew()
func (n Node) IsNew() bool {
	return n.rec().token == Token_NEW
}

// Java: Node.isNot()
func (n Node) IsNot() bool {
	return n.rec().token == Token_NOT
}

// Java: Node.isNull()
func (n Node) IsNull() bool {
	return n.rec().token == Token_NULL
}

// Java: Node.isNullishCoalesce()
func (n Node) IsNullishCoalesce() bool {
	return n.rec().token == Token_COALESCE
}

// Java: Node.isAssignNullishCoalesce()
func (n Node) IsAssignNullishCoalesce() bool {
	return n.rec().token == Token_ASSIGN_COALESCE
}

// Java: Node.isNumber()
func (n Node) IsNumber() bool {
	return n.rec().token == Token_NUMBER
}

// Java: Node.isObjectLit()
func (n Node) IsObjectLit() bool {
	return n.rec().token == Token_OBJECTLIT
}

// Java: Node.isObjectPattern()
func (n Node) IsObjectPattern() bool {
	return n.rec().token == Token_OBJECT_PATTERN
}

// Java: Node.isOptChainCall()
func (n Node) IsOptChainCall() bool {
	return n.rec().token == Token_OPTCHAIN_CALL
}

// Java: Node.isOptChainGetElem()
func (n Node) IsOptChainGetElem() bool {
	return n.rec().token == Token_OPTCHAIN_GETELEM
}

// Java: Node.isOptChainGetProp()
func (n Node) IsOptChainGetProp() bool {
	return n.rec().token == Token_OPTCHAIN_GETPROP
}

// Java: Node.isOr()
func (n Node) IsOr() bool {
	return n.rec().token == Token_OR
}

// Java: Node.isAssignOr()
func (n Node) IsAssignOr() bool {
	return n.rec().token == Token_ASSIGN_OR
}

// Java: Node.isParamList()
func (n Node) IsParamList() bool {
	return n.rec().token == Token_PARAM_LIST
}

// Java: Node.isPos()
func (n Node) IsPos() bool {
	return n.rec().token == Token_POS
}

// Java: Node.isRegExp()
func (n Node) IsRegExp() bool {
	return n.rec().token == Token_REGEXP
}

// Java: Node.isRest()
func (n Node) IsRest() bool {
	return n.rec().token == Token_ITER_REST || n.rec().token == Token_OBJECT_REST
}

// Java: Node.isObjectRest()
func (n Node) IsObjectRest() bool {
	return n.rec().token == Token_OBJECT_REST
}

// Java: Node.isReturn()
func (n Node) IsReturn() bool {
	return n.rec().token == Token_RETURN
}

// Java: Node.isScript()
func (n Node) IsScript() bool {
	return n.rec().token == Token_SCRIPT
}

// Java: Node.isSetterDef()
func (n Node) IsSetterDef() bool {
	return n.rec().token == Token_SETTER_DEF
}

// Java: Node.isSpread()
func (n Node) IsSpread() bool {
	return n.rec().token == Token_ITER_SPREAD || n.rec().token == Token_OBJECT_SPREAD
}

// Java: Node.isString()
func (n Node) IsString() bool {
	return n.rec().token == Token_STRINGLIT
}

// Java: Node.isStringKey()
func (n Node) IsStringKey() bool {
	return n.rec().token == Token_STRING_KEY
}

// Java: Node.isStringLit()
func (n Node) IsStringLit() bool {
	return n.rec().token == Token_STRINGLIT
}

// Java: Node.isSuper()
func (n Node) IsSuper() bool {
	return n.rec().token == Token_SUPER
}

// Java: Node.isSwitch()
func (n Node) IsSwitch() bool {
	return n.rec().token == Token_SWITCH
}

// Java: Node.isSwitchBody()
func (n Node) IsSwitchBody() bool {
	return n.rec().token == Token_SWITCH_BODY
}

// Java: Node.isTaggedTemplateLit()
func (n Node) IsTaggedTemplateLit() bool {
	return n.rec().token == Token_TAGGED_TEMPLATELIT
}

// Java: Node.isTemplateLit()
func (n Node) IsTemplateLit() bool {
	return n.rec().token == Token_TEMPLATELIT
}

// Java: Node.isTemplateLitString()
func (n Node) IsTemplateLitString() bool {
	return n.rec().token == Token_TEMPLATELIT_STRING
}

// Java: Node.isTemplateLitSub()
func (n Node) IsTemplateLitSub() bool {
	return n.rec().token == Token_TEMPLATELIT_SUB
}

// Java: Node.isThis()
func (n Node) IsThis() bool {
	return n.rec().token == Token_THIS
}

// Java: Node.isThrow()
func (n Node) IsThrow() bool {
	return n.rec().token == Token_THROW
}

// Java: Node.isTrue()
func (n Node) IsTrue() bool {
	return n.rec().token == Token_TRUE
}

// Java: Node.isTry()
func (n Node) IsTry() bool {
	return n.rec().token == Token_TRY
}

// Java: Node.isTypeOf()
func (n Node) IsTypeOf() bool {
	return n.rec().token == Token_TYPEOF
}

// Java: Node.isVar()
func (n Node) IsVar() bool {
	return n.rec().token == Token_VAR
}

// Java: Node.isVoid()
func (n Node) IsVoid() bool {
	return n.rec().token == Token_VOID
}

// Java: Node.isWhile()
func (n Node) IsWhile() bool {
	return n.rec().token == Token_WHILE
}

// Java: Node.isWith()
func (n Node) IsWith() bool {
	return n.rec().token == Token_WITH
}

// Java: Node.isYield()
func (n Node) IsYield() bool {
	return n.rec().token == Token_YIELD
}

// Java: Node.isDeclare()
func (n Node) IsDeclare() bool {
	return n.rec().token == Token_DECLARE
}

// Go-only: language-adapter — recover the Java Object field only for printing
// and equality; hot getters use the tagged pointer directly.
func (n Node) typeOrColorValue() any {
	if n.tc() == 0 {
		return nil
	}
	if (n.rec().kindFlags & kindFlagTypeIsColor) != 0 {
		if c := n.store().colorAt(n.tc()); c != nil {
			return c
		}
		return nil
	}
	return jsTypeOf(n.store().typeAt(n.tc()))
}

func (n Node) nodeStringValue() string {
	return rhinoStringFromID(rhinoStringID(n.rec().payload))
}

// CanResetInPlace reports whether n has no cold extBox, extra side-table entry
// or diagnostic pin, making in-place token/payload mutation exact.
func (n Node) CanResetInPlace() bool {
	return (n.rec().kindFlags & (kindFlagHasExtBox | kindFlagHasExtra)) == 0
}

// PinForDiagnostic marks n so later optimization passes do not recycle n
// in-place, preserving n's original length/token for deferred JSError formatting.
func (n Node) PinForDiagnostic() {
	n.rec().kindFlags |= kindFlagHasExtra
}

// MutateGetPropToNameInPlace rewrites an attached GETPROP node into a NAME node
// in-place (preserving its sibling links, JSType/Color, and Srcref(n)) when n
// has no cold extBox/extra metadata. Interns name before originalName to match
// NodeUtil_NewNameWithOriginalName's exact RhinoStringPool interning order.
func (n Node) MutateGetPropToNameInPlace(name string, originalName string, isConstant bool) bool {
	if !n.CanResetInPlace() {
		return false
	}
	n.DetachChildren()
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
	nameID := rhinoStringPoolIntern(name)
	origID := rhinoStringPoolIntern(originalName)
	if (n.props() & propFlagHasJSDoc) != 0 {
		n.setAux(nodeAux{})
	}
	n.rec().token = Token_NAME
	n.rec().kindFlags = uint8(node_kindString) | (n.rec().kindFlags & kindFlagTypeIsColor)
	flags := n.props() & propFlag_IS_CONSTANT_NAME
	if isConstant {
		flags |= propFlag_IS_CONSTANT_NAME
	}
	n.setProps(flags)
	n.rec().payload = uint32(nameID)
	n.node_setOriginalNameID(origID)
	return true
}

// MutateCallTargetGetPropToNameInPlace rewrites a CALL's first-child GETPROP
// `o.foo` into `NAME(newMethodName)` in-place and inserts detached `o` right
// after it (`foo(o, ...)`), matching DevirtualizeMethods.RewriteCall.
func (n Node) MutateCallTargetGetPropToNameInPlace(call Node, receiver Node, newMethodName string, isConstantName bool) bool {
	if !n.CanResetInPlace() {
		return false
	}
	n.DetachChildren()
	receiver.SrcrefIfMissing(n)
	nameID := rhinoStringPoolIntern(newMethodName)
	origID := n.node_originalNameID()
	if (n.props() & propFlagHasJSDoc) != 0 {
		n.setAux(nodeAux{})
	}
	n.rec().token = Token_NAME
	n.rec().kindFlags = uint8(node_kindString) | (n.rec().kindFlags & kindFlagTypeIsColor)
	flags := uint64(0)
	if isConstantName {
		flags |= propFlag_IS_CONSTANT_NAME
	}
	n.setProps(flags)
	n.rec().payload = uint32(nameID)
	n.node_setOriginalNameID(origID)
	receiver.InsertAfter(n)
	return true
}

// MutateThisToNameIDInPlace rewrites an attached leaf THIS node into
// NAME(nameID) in-place, preserving Srcref(n) and CopyTypeFrom(n).
func (n Node) MutateThisToNameIDInPlace(nameID uint32) bool {
	if !n.CanResetInPlace() || n.rec().first != 0 {
		return false
	}
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
	if (n.props() & propFlagHasJSDoc) != 0 {
		n.setAux(nodeAux{})
	}
	n.rec().token = Token_NAME
	n.rec().kindFlags = uint8(node_kindString) | (n.rec().kindFlags & kindFlagTypeIsColor)
	n.setProps(0)
	n.rec().payload = uint32(nameID)
	return true
}

// ResetDetachedGetPropToNameID resets a detached GETPROP node (CanResetInPlace)
// into IR_Name(nameID).SrcrefIfMissing(self), optionally setting IS_CONSTANT_NAME.
func (n Node) ResetDetachedGetPropToNameID(nameID uint32, isConstantName bool) {
	origID := n.node_originalNameID()
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = Token_NAME
	n.rec().kindFlags = uint8(node_kindString)
	flags := uint64(0)
	if isConstantName {
		flags |= propFlag_IS_CONSTANT_NAME
	}
	n.setProps(flags)
	n.rec().payload = uint32(nameID)
	n.node_setOriginalNameID(origID)
}

// MutateAttachedStmtToVarFromSrcref resets an attached childless statement node
// (CanResetInPlace) into IR_Var(nameChild).SrcrefIfMissing(srcref) in-place.
func (n Node) MutateAttachedStmtToVarFromSrcref(nameChild Node, srcref Node, jsdoc *JSDocInfo) {
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
	origID := srcref.node_originalNameID() // before n changes: srcref may be n
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = Token_VAR
	n.rec().kindFlags = 0
	n.setProps(0)
	n.rec().sourceFile = srcref.rec().sourceFile
	n.rec().payload = 0
	n.node_setOriginalNameID(origID)
	n.rec().linenoCharno = srcref.rec().linenoCharno
	n.rec().length16 = srcref.rec().length16
	n.AddChildToFront(nameChild)
	if jsdoc != nil {
		n.SetJSDocInfo(jsdoc)
	}
}

// ResetDetachedToFreshNameID resets a detached childless node (CanResetInPlace)
// into a fresh IR_Name(nameID) with unset source location.
func (n Node) ResetDetachedToFreshNameID(nameID uint32) Node {
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = Token_NAME
	n.rec().kindFlags = uint8(node_kindString)
	n.setProps(0)
	n.rec().sourceFile = 0
	n.rec().payload = uint32(nameID)
	n.node_setOriginalNameID(0)
	n.rec().linenoCharno = -1
	n.rec().length16 = 0
	return n
}

// CanCloneLeafInPlace reports whether src is a childless leaf node with no cold
// extBox/extra metadata or overflow source file, allowing zero-allocation
// in-place leaf cloning via OverwriteLeafFromClone.
func (n Node) CanCloneLeafInPlace() bool {
	return n.rec().first == 0 && n.CanResetInPlace()
}

// OverwriteLeafFromClone overwrites an attached childless leaf node (which must
// satisfy CanCloneLeafInPlace) with a clone of src (which must also satisfy
// CanCloneLeafInPlace), matching src.CloneTree() + n.ReplaceWith(clone) with
// zero allocations.
func (n Node) OverwriteLeafFromClone(src Node) {
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
	oldSF := n.rec().sourceFile
	oldOrigID := n.node_originalNameID()
	srcOrigID := src.node_originalNameID()
	n.rec().token = src.rec().token
	n.rec().kindFlags = src.rec().kindFlags & (kindMask | kindFlagTypeIsColor)
	n.setTC(n.store().translateTC(src.store(), src.tc(), src.rec().kindFlags&kindFlagTypeIsColor != 0))
	n.setAux(src.auxValue())
	n.setProps(src.props())
	n.rec().payload = src.rec().payload
	if src.rec().sourceFile == 0 {
		n.rec().sourceFile = oldSF
	} else {
		n.rec().sourceFile = src.rec().sourceFile
		n.rec().linenoCharno = src.rec().linenoCharno
		n.rec().length16 = src.rec().length16
	}
	if srcOrigID == 0 {
		srcOrigID = oldOrigID
	}
	n.node_setOriginalNameID(srcOrigID)
}

// ResetTemplateUndefinedSrcref updates a detached 2-node `void 0` template in
// place to match NodeUtil_NewUndefinedNode(src). Returns false if either node in
// the template or src has cold/overflow metadata so the caller can allocate a
// fresh template.
func (n Node) ResetTemplateUndefinedSrcref(src Node) bool {
	num := n.firstNode()
	if num.IsNil() || !n.CanResetInPlace() || !num.CanResetInPlace() {
		return false
	}
	if src.IsNil() {
		n.setTC(0)
		n.setAux(nodeAux{})
		n.rec().kindFlags = 0
		n.setProps(0)
		n.rec().sourceFile = 0
		n.rec().payload = 0
		n.rec().linenoCharno = -1
		n.rec().length16 = 0
		num.setTC(0)
		num.setAux(nodeAux{})
		num.rec().kindFlags = uint8(node_kindNumber)
		num.setProps(0)
		num.rec().sourceFile = 0
		num.rec().payload = 0
		num.rec().linenoCharno = -1
		num.rec().length16 = 0
		n.node_setOriginalNameID(0)
		num.node_setOriginalNameID(0)
		return true
	}
	if !src.CanResetInPlace() || src.node_originalNameID() != 0 {
		return false
	}
	sf := src.rec().sourceFile
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().kindFlags = 0
	n.setProps(0)
	n.rec().sourceFile = sf
	n.rec().payload = 0
	n.rec().linenoCharno = src.rec().linenoCharno
	n.rec().length16 = src.rec().length16
	num.setTC(0)
	num.setAux(nodeAux{})
	num.rec().kindFlags = uint8(node_kindNumber)
	num.setProps(0)
	num.rec().sourceFile = sf
	num.rec().payload = 0
	num.rec().linenoCharno = src.rec().linenoCharno
	num.rec().length16 = src.rec().length16
	n.node_setOriginalNameID(0)
	num.node_setOriginalNameID(0)
	return true
}

// MutateAttachedToFreshVarInPlace resets an attached statement node
// (CanResetInPlace) into a Token_VAR node matching n.ReplaceWith(NewNode(Token_VAR))
// (preserving n's Srcref), with optional single child nameChild and optional JSDocInfo.
func (n Node) MutateAttachedToFreshVarInPlace(nameChild Node, jsdoc *JSDocInfo) {
	n.DetachChildren()
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
	origID := n.node_originalNameID()
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = Token_VAR
	n.rec().kindFlags = 0
	n.setProps(0)
	n.rec().payload = 0
	n.node_setOriginalNameID(origID)
	if !nameChild.IsNil() {
		n.AddChildToFront(nameChild)
	}
	if jsdoc != nil {
		n.SetJSDocInfo(jsdoc)
	}
}

// ResetDetachedToNameWithOriginalNameInPlace resets a detached node
// (CanResetInPlace) into a NAME(name) node with Srcref(srcref), originalName,
// and IS_CONSTANT_NAME = isConstant (srcref must satisfy CanResetInPlace).
func (n Node) ResetDetachedToNameWithOriginalNameInPlace(name string, originalName string, srcref Node, isConstant bool) Node {
	n.DetachChildren()
	nameID := rhinoStringPoolIntern(name)
	origID := rhinoStringPoolIntern(originalName)
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = Token_NAME
	n.rec().kindFlags = uint8(node_kindString)
	flags := uint64(0)
	n.rec().sourceFile = srcref.rec().sourceFile
	if isConstant {
		flags |= propFlag_IS_CONSTANT_NAME
	}
	n.setProps(flags)
	n.rec().payload = uint32(nameID)
	n.node_setOriginalNameID(origID)
	n.rec().linenoCharno = srcref.rec().linenoCharno
	n.rec().length16 = srcref.rec().length16
	return n
}

// MutateAttachedToExprResultInPlace resets an attached statement node
// (CanResetInPlace) into IR_ExprResult(exprChild).Srcref(exprChild) + ReplaceWith
// (preserving n's SrcrefIfMissing) in-place (exprChild must be detached and
// satisfy CanResetInPlace).
func (n Node) MutateAttachedToExprResultInPlace(exprChild Node) Node {
	n.DetachChildren()
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
	origID := exprChild.node_originalNameID()
	if origID == 0 {
		origID = n.node_originalNameID()
	}
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = Token_EXPR_RESULT
	n.rec().kindFlags = 0
	n.setProps(0)
	if exprChild.rec().sourceFile != 0 {
		n.rec().sourceFile = exprChild.rec().sourceFile
		n.rec().linenoCharno = exprChild.rec().linenoCharno
		n.rec().length16 = exprChild.rec().length16
	}
	n.rec().payload = 0
	n.node_setOriginalNameID(origID)
	n.AddChildToFront(exprChild)
	return n
}

// ResetDetachedToBinaryOpFromSrcref resets a detached node (CanResetInPlace)
// into NewNodeWith2Children(op, left, right).Srcref(srcref) in-place (srcref
// must satisfy CanResetInPlace; left and right must be detached).
func (n Node) ResetDetachedToBinaryOpFromSrcref(op Token, left Node, right Node, srcref Node) Node {
	n.DetachChildren()
	origID := srcref.node_originalNameID() // before n changes: srcref may be n
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = op
	n.rec().kindFlags = 0
	n.setProps(0)
	n.rec().sourceFile = srcref.rec().sourceFile
	n.rec().payload = 0
	n.node_setOriginalNameID(origID)
	n.rec().linenoCharno = srcref.rec().linenoCharno
	n.rec().length16 = srcref.rec().length16
	n.AddChildToFront(left)
	n.AddChildToBack(right)
	return n
}
