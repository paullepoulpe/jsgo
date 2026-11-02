// Port of com/google/javascript/jscomp/FunctionArgumentInjector.java @ closure-compiler v20261028

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// A string to use to represent "this".  Anything that is not a valid
// identifier can be used, so we use "this".
const FunctionArgumentInjector_THIS_MARKER = "this"

var FunctionArgumentInjector_THIS_MARKER_NODE = rhino.IR_Name(FunctionArgumentInjector_THIS_MARKER)

var (
	calleeArgumentsIdentID = uint32(strpool.Intern("arguments"))
	calleeEvalIdentID      = uint32(strpool.Intern("eval"))
	calleeThisMarkerID     = uint32(strpool.Intern(FunctionArgumentInjector_THIS_MARKER))
)

type calleeParamSummary struct {
	nameID   uint32
	refCount uint32
	modified bool
	isRest   bool
}

func (p calleeParamSummary) nameStr() string {
	return strpool.String(strpool.ID(p.nameID))
}

type calleeSideEffectSummary struct {
	computed             bool
	thisAfterSideEffect  bool
	overflow             bool
	afterSideEffectMask  uint64
	namesBeforeParamMask uint64
}

// calleeSummary caches structural and parameter-usage properties of a candidate
// FUNCTION node from a single AST body walk, avoiding repeated O(Params * |Body|)
// traversals across call-site checks, cost estimation, and block mutation.
//
// All small buffers are value arrays with separate overflow slices (never slices
// aliasing sibling array fields), so copying calleeSummary by value is safe.
type calleeSummary struct {
	fnNode                 rhino.Node
	paramBloom             uint64
	bodyAstSize            uint32
	shallowReturnCount     uint32
	numParams              uint16
	numModified            uint16
	paramsPositional       bool
	isTrivialBody          bool
	hasInnerClosures       bool
	hasShallowDeclarations bool
	hasLabels              bool
	hasInnerFunctionDecls  bool
	hasVarOrLetDecls       bool
	referencesThis         bool
	blocksMinReqs          bool
	hasParamShadowing      bool
	hasSubScopeDecls       bool

	paramsBuf             [4]calleeParamSummary
	paramsOverflow        []calleeParamSummary
	modifiedOrderBuf      [4]uint16
	modifiedOrderOverflow []uint16

	sideEffects calleeSideEffectSummary
}

func calleeIDBloomBit(id uint32) uint64 {
	return uint64(1) << (uint(id) & 63)
}

func (s *calleeSummary) params() []calleeParamSummary {
	if int(s.numParams) <= len(s.paramsBuf) {
		return s.paramsBuf[:s.numParams]
	}
	return s.paramsOverflow[:s.numParams]
}

func (s *calleeSummary) modifiedParamIndices() []uint16 {
	if int(s.numModified) <= len(s.modifiedOrderBuf) {
		return s.modifiedOrderBuf[:s.numModified]
	}
	return s.modifiedOrderOverflow[:s.numModified]
}

func (s *calleeSummary) addParam(p calleeParamSummary) {
	s.paramBloom |= calleeIDBloomBit(p.nameID)
	idx := int(s.numParams)
	if idx < len(s.paramsBuf) {
		s.paramsBuf[idx] = p
	} else if idx == len(s.paramsBuf) {
		s.paramsOverflow = make([]calleeParamSummary, len(s.paramsBuf)+1, 8)
		copy(s.paramsOverflow, s.paramsBuf[:])
		s.paramsOverflow[idx] = p
	} else {
		s.paramsOverflow = append(s.paramsOverflow, p)
	}
	s.numParams++
}

func (s *calleeSummary) addModifiedParam(paramIdx uint16) {
	idx := int(s.numModified)
	if idx < len(s.modifiedOrderBuf) {
		s.modifiedOrderBuf[idx] = paramIdx
	} else if idx == len(s.modifiedOrderBuf) {
		s.modifiedOrderOverflow = make([]uint16, len(s.modifiedOrderBuf)+1, 8)
		copy(s.modifiedOrderOverflow, s.modifiedOrderBuf[:])
		s.modifiedOrderOverflow[idx] = paramIdx
	} else {
		s.modifiedOrderOverflow = append(s.modifiedOrderOverflow, paramIdx)
	}
	s.numModified++
}

func (s *calleeSummary) findParamByID(id uint32) int {
	if s.numParams == 0 || (s.paramBloom&calleeIDBloomBit(id)) == 0 {
		return -1
	}
	params := s.params()
	for i := range params {
		if params[i].nameID == id {
			return i
		}
	}
	return -1
}

func (s *calleeSummary) canReuseAfterRename(renamedClone bool) bool {
	return s.paramsPositional && (!renamedClone || (!s.hasParamShadowing && !s.hasInnerFunctionDecls))
}

func (s *calleeSummary) canReuseSideEffects(renamedClone bool) bool {
	return s.paramsPositional && int(s.numParams) <= 64 &&
		(!renamedClone || (!s.hasSubScopeDecls && !s.hasParamShadowing && !s.hasInnerFunctionDecls))
}

type calleeScanner struct {
	dst               *calleeSummary
	body              rhino.Node
	fnNameID          uint32
	fnRecursionNameID uint32
	hasFnName         bool
	hasFnRecursion    bool
	hasCondNode       bool
}

func (cs *calleeScanner) walk(n, parent rhino.Node, inInnerClosure, inInnerFn, inReceiverScope, inShallowStmt bool) {
	dst := cs.dst
	dst.bodyAstSize++
	tok := n.GetToken()

	isNonStaticFieldInit := !parent.IsStaticMember() &&
		(parent.IsMemberFieldDef() || (parent.IsComputedFieldDef() && n != parent.GetFirstChild()))
	if tok == rhino.Token_FUNCTION || isNonStaticFieldInit {
		inInnerClosure = true
		dst.hasInnerClosures = true
		cs.hasCondNode = true
	}

	if inShallowStmt {
		if NodeUtil_IsDeclaration(n) {
			dst.hasShallowDeclarations = true
		}
		if tok == rhino.Token_RETURN {
			dst.shallowReturnCount++
		}
	}

	switch tok {
	case rhino.Token_FUNCTION:
		if !inInnerFn && NodeUtil_IsFunctionDeclaration(n) {
			dst.hasInnerFunctionDecls = true
		}
		dst.hasSubScopeDecls = true
		inInnerFn = true
		if !n.IsArrowFunction() {
			inReceiverScope = false
		}
	case rhino.Token_THIS:
		if inReceiverScope {
			dst.referencesThis = true
		}
	case rhino.Token_SUPER:
		if inReceiverScope {
			dst.referencesThis = true
		}
		dst.blocksMinReqs = true
	case rhino.Token_NEW_TARGET:
		dst.blocksMinReqs = true
	case rhino.Token_LABEL:
		dst.hasLabels = true
	case rhino.Token_VAR:
		if !inInnerFn {
			dst.hasVarOrLetDecls = true
		}
	case rhino.Token_LET:
		if !inInnerFn {
			dst.hasVarOrLetDecls = true
		}
		if parent != cs.body {
			dst.hasSubScopeDecls = true
		}
	case rhino.Token_CONST, rhino.Token_CLASS:
		if parent != cs.body {
			dst.hasSubScopeDecls = true
		}
	case rhino.Token_CATCH:
		dst.hasSubScopeDecls = true
	case rhino.Token_AND,
		rhino.Token_OR,
		rhino.Token_HOOK,
		rhino.Token_COALESCE,
		rhino.Token_OPTCHAIN_CALL,
		rhino.Token_OPTCHAIN_GETELEM,
		rhino.Token_OPTCHAIN_GETPROP:
		cs.hasCondNode = true
	case rhino.Token_NAME:
		nameID := n.GetStringID()
		if (inReceiverScope && nameID == calleeArgumentsIdentID) ||
			nameID == calleeEvalIdentID ||
			(cs.hasFnName && nameID == cs.fnNameID) ||
			(cs.hasFnRecursion && nameID == cs.fnRecursionNameID) {
			dst.blocksMinReqs = true
		}
		if idx := dst.findParamByID(nameID); idx >= 0 {
			params := dst.params()
			params[idx].refCount++
			if !params[idx].isRest && !params[idx].modified &&
				(inInnerClosure || FunctionArgumentInjector_CanNameValueChange(n)) {
				params[idx].modified = true
				dst.addModifiedParam(uint16(idx))
			}
			role := ast.RoleIn(n, parent)
			if inInnerFn ||
				role == ast.RoleNameDeclarationDeclarators ||
				role == ast.RoleCatchParam ||
				role == ast.RoleParamListParams ||
				role == ast.RoleFunctionName ||
				role == ast.RoleClassName ||
				NodeUtil_IsDeclarationLValue(n) {
				dst.hasParamShadowing = true
			}
		}
	}

	childrenShallow := inShallowStmt &&
		(tok == rhino.Token_ROOT || tok == rhino.Token_BLOCK || (tok != rhino.Token_FUNCTION && NodeUtil_IsShallowStatementTree(parent)))
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		cs.walk(c, n, inInnerClosure, inInnerFn, inReceiverScope, childrenShallow)
	}
}

// scanCalleeSummary populates dst from a single pass over fnNode's parameters
// and body. It returns false if fnNode has non-inlinable parameter patterns or
// duplicate formal parameter names (though body properties are always populated
// when fnNode is a FUNCTION node).
func scanCalleeSummary(fnName string, fnNode rhino.Node, dst *calleeSummary) bool {
	*dst = calleeSummary{fnNode: fnNode}
	if !fnNode.IsFunction() {
		return false
	}
	fn := ast.MustFunction(fnNode)

	paramsPositional := true
	for p, more := fn.Params().FirstParam(); more; p, more = p.Next() {
		var pNameID uint32
		var isRest bool
		switch p.Token() {
		case rhino.Token_NAME:
			pNameID = p.Node().GetStringID()
		case rhino.Token_ITER_REST:
			restName, ok := p.MustIterRest().Target().AsName()
			if !ok {
				paramsPositional = false
				continue
			}
			pNameID = restName.Node().GetStringID()
			isRest = true
		default:
			paramsPositional = false
			continue
		}
		if dst.findParamByID(pNameID) >= 0 {
			paramsPositional = false
			continue
		}
		dst.addParam(calleeParamSummary{
			nameID: pNameID,
			isRest: isRest,
		})
	}
	dst.paramsPositional = paramsPositional

	body := fn.Body().Node()
	cs := calleeScanner{
		dst:  dst,
		body: body,
	}
	if fnName != "" {
		if id := uint32(strpool.Lookup(fnName)); id != 0 {
			cs.fnNameID = id
			cs.hasFnName = true
		}
	}
	if name, ok := fn.Name(); ok {
		if recID := name.Node().GetStringID(); recID != 0 && (!cs.hasFnName || recID != cs.fnNameID) {
			cs.fnRecursionNameID = recID
			cs.hasFnRecursion = true
		}
	}

	cs.walk(body, fnNode, false, false, true, true)

	if !body.HasChildren() {
		dst.isTrivialBody = true
	} else if body.HasOneChild() {
		onlyChild := body.GetFirstChild()
		dst.isTrivialBody = (onlyChild.IsReturn() || onlyChild.IsExprResult()) && !cs.hasCondNode
	}

	return paramsPositional
}

func resolveCalleeSummary(fnName string, fnNode rhino.Node, summary *calleeSummary) *calleeSummary {
	if summary != nil && summary.fnNode == fnNode {
		return summary
	}
	local := new(calleeSummary)
	scanCalleeSummary(fnName, fnNode, local)
	return local
}

func (f *FunctionArgumentInjector) findModifiedParametersByIDInto(
	fnNode rhino.Node, _ *inlineIDSet, dstModified *inlineIDSet,
) {
	dstModified.clear()
	fn := ast.MustFunction(fnNode)
	if !fn.Params().HasParams() {
		return
	}
	for p, more := fn.Params().FirstParam(); more; p, more = p.Next() {
		if p.Token() != rhino.Token_NAME && p.Token() != rhino.Token_ITER_REST {
			panic(jutil.NewIllegalStateExceptionMsg("Not supported: " + p.Node().String()))
		}
	}
	var sum calleeSummary
	scanCalleeSummary("", fnNode, &sum)
	params := sum.params()
	for _, idx := range sum.modifiedParamIndices() {
		dstModified.add(params[idx].nameID)
	}
}

func gatherLocalNameIDs(n rhino.Node, names *inlineIDSet) {
	if n.IsFunction() {
		if NodeUtil_IsFunctionDeclaration(n) {
			name, _ := ast.MustFunction(n).Name()
			names.add(name.Node().GetStringID())
		}
		// Don't traverse into inner function scopes;
		return
	} else if n.IsName() {
		switch ast.RoleOf(n) {
		case ast.RoleNameDeclarationDeclarators, ast.RoleCatchParam:
			names.add(n.GetStringID())
		}
	}

	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		gatherLocalNameIDs(c, names)
	}
}
