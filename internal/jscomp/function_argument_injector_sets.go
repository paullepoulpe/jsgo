// Port of com/google/javascript/jscomp/FunctionArgumentInjector.java @ closure-compiler v20261021

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

const inlineSmallCap = 16

// inlineIDSet is a reusable insertion-ordered set of 32-bit strpool IDs.
// For <= inlineSmallCap elements it uses a linear scan over items (initially
// backed by buf); larger sets lazily populate index.
type inlineIDSet struct {
	items []uint32
	index map[uint32]struct{}
	buf   [inlineSmallCap]uint32
}

func (s *inlineIDSet) clear() {
	if s.items == nil {
		s.items = s.buf[:0]
	} else {
		s.items = s.items[:0]
	}
	if len(s.index) > 0 {
		clear(s.index)
	}
}

func (s *inlineIDSet) isEmpty() bool {
	return s == nil || len(s.items) == 0
}

func (s *inlineIDSet) size() int {
	if s == nil {
		return 0
	}
	return len(s.items)
}

func (s *inlineIDSet) contains(id uint32) bool {
	if s == nil || len(s.items) == 0 || id == 0 {
		return false
	}
	if len(s.items) <= inlineSmallCap {
		for _, existing := range s.items {
			if existing == id {
				return true
			}
		}
		return false
	}
	_, ok := s.index[id]
	return ok
}

func (s *inlineIDSet) add(id uint32) bool {
	if id == 0 {
		return false
	}
	if len(s.items) <= inlineSmallCap {
		for _, existing := range s.items {
			if existing == id {
				return false
			}
		}
		if s.items == nil {
			s.items = s.buf[:0]
		}
		if len(s.items) == inlineSmallCap {
			if s.index == nil {
				s.index = make(map[uint32]struct{}, inlineSmallCap*2)
			} else {
				clear(s.index)
			}
			for _, existing := range s.items {
				s.index[existing] = struct{}{}
			}
			s.index[id] = struct{}{}
		}
		s.items = append(s.items, id)
		return true
	}
	if _, exists := s.index[id]; exists {
		return false
	}
	s.index[id] = struct{}{}
	s.items = append(s.items, id)
	return true
}

func (s *inlineIDSet) addAll(other *inlineIDSet) {
	if other == nil {
		return
	}
	for _, id := range other.items {
		s.add(id)
	}
}

func (s *inlineIDSet) copyFrom(other *inlineIDSet) {
	s.clear()
	if other != nil {
		for _, id := range other.items {
			s.add(id)
		}
	}
}

func (s *inlineIDSet) fromLinkedHashSet(src *jutil.LinkedHashSet[string], intern bool) {
	s.clear()
	if src == nil || src.IsEmpty() {
		return
	}
	for name := range src.All() {
		var id uint32
		if intern {
			id = rhino.RhinoStringPool_InternID(name)
		} else {
			id = rhino.RhinoStringPool_LookupID(name)
		}
		if id != 0 {
			s.add(id)
		}
	}
}

func (s *inlineIDSet) toLinkedHashSet() *jutil.LinkedHashSet[string] {
	out := jutil.NewLinkedHashSet[string]()
	if s != nil {
		for _, id := range s.items {
			out.Add(rhino.RhinoStringPool_StringFromID(id))
		}
	}
	return out
}

type inlineIDNodeEntry struct {
	id   uint32
	node rhino.Node
}

// inlineIDNodeMap is a reusable map of 32-bit strpool IDs to replacement Nodes.
// For <= inlineSmallCap entries it scans entries linearly; larger maps lazily
// populate index.
type inlineIDNodeMap struct {
	entries []inlineIDNodeEntry
	index   map[uint32]int
	buf     [inlineSmallCap]inlineIDNodeEntry
}

func (m *inlineIDNodeMap) clear() {
	if m.entries == nil {
		m.entries = m.buf[:0]
	} else {
		clear(m.entries)
		m.entries = m.entries[:0]
	}
	if len(m.index) > 0 {
		clear(m.index)
	}
}

func (m *inlineIDNodeMap) get(id uint32) rhino.Node {
	if m == nil || len(m.entries) == 0 || id == 0 {
		return rhino.NilNode
	}
	if len(m.entries) <= inlineSmallCap {
		for i := range m.entries {
			if m.entries[i].id == id {
				return m.entries[i].node
			}
		}
		return rhino.NilNode
	}
	if idx, ok := m.index[id]; ok {
		return m.entries[idx].node
	}
	return rhino.NilNode
}

func (m *inlineIDNodeMap) put(id uint32, node rhino.Node) rhino.Node {
	if id == 0 {
		return rhino.NilNode
	}
	if len(m.entries) <= inlineSmallCap {
		for i := range m.entries {
			if m.entries[i].id == id {
				prev := m.entries[i].node
				m.entries[i].node = node
				return prev
			}
		}
		if m.entries == nil {
			m.entries = m.buf[:0]
		}
		if len(m.entries) == inlineSmallCap {
			if m.index == nil {
				m.index = make(map[uint32]int, inlineSmallCap*2)
			} else {
				clear(m.index)
			}
			for i, e := range m.entries {
				if e.id != 0 {
					m.index[e.id] = i
				}
			}
			m.index[id] = inlineSmallCap
		}
		m.entries = append(m.entries, inlineIDNodeEntry{id: id, node: node})
		return rhino.NilNode
	}
	if idx, ok := m.index[id]; ok {
		prev := m.entries[idx].node
		m.entries[idx].node = node
		return prev
	}
	m.index[id] = len(m.entries)
	m.entries = append(m.entries, inlineIDNodeEntry{id: id, node: node})
	return rhino.NilNode
}

func (m *inlineIDNodeMap) remove(id uint32) {
	if m == nil || len(m.entries) == 0 || id == 0 {
		return
	}
	if len(m.entries) <= inlineSmallCap {
		for i := range m.entries {
			if m.entries[i].id == id {
				m.entries[i] = inlineIDNodeEntry{}
				return
			}
		}
		return
	}
	if idx, ok := m.index[id]; ok {
		m.entries[idx] = inlineIDNodeEntry{}
		delete(m.index, id)
	}
}

// A string to use to represent "this".  Anything that is not a valid
// identifier can be used, so we use "this".
const FunctionArgumentInjector_THIS_MARKER = "this"

var FunctionArgumentInjector_THIS_MARKER_NODE = rhino.IR_Name(FunctionArgumentInjector_THIS_MARKER)

const FunctionArgumentInjector_REST_MARKER = "rest param"

var (
	calleeArgumentsIdentID = uint32(strpool.Intern("arguments"))
	calleeEvalIdentID      = uint32(strpool.Intern("eval"))
	calleeThisMarkerID     = uint32(strpool.Intern(FunctionArgumentInjector_THIS_MARKER))
	calleeRestMarkerID     = uint32(strpool.Intern(FunctionArgumentInjector_REST_MARKER))
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
	return !renamedClone || (!s.hasParamShadowing && !s.hasInnerFunctionDecls)
}

func (s *calleeSummary) canReuseSideEffects(renamedClone bool) bool {
	return int(s.numParams) <= 64 &&
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
// duplicate formal parameter names.
func scanCalleeSummary(fnName string, fnNode rhino.Node, dst *calleeSummary) bool {
	*dst = calleeSummary{fnNode: fnNode}
	if !fnNode.IsFunction() {
		return false
	}
	fn := ast.MustFunction(fnNode)

	for p, more := fn.Params().FirstParam(); more; p, more = p.Next() {
		var pNameID uint32
		var isRest bool
		switch p.Token() {
		case rhino.Token_NAME:
			pNameID = p.Node().GetStringID()
		case rhino.Token_ITER_REST:
			restName, ok := p.MustIterRest().Target().AsName()
			if !ok {
				return false
			}
			pNameID = restName.Node().GetStringID()
			isRest = true
		default:
			return false
		}
		if dst.findParamByID(pNameID) >= 0 {
			return false
		}
		dst.addParam(calleeParamSummary{
			nameID: pNameID,
			isRest: isRest,
		})
	}

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
		recID := name.Node().GetStringID()
		if !cs.hasFnName || recID != cs.fnNameID {
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

	return true
}

// FindModifiedParameters retrieves a set of names that can not be safely
// substituted in place.
//
// Example:
//
//	function(a) {
//	  a = 0;
//	}
//
// Inlining this without taking precautions would cause the call site value
// to be modified (bad).
//
// Java: FunctionArgumentInjector.findModifiedParameters(Node)
func (f *FunctionArgumentInjector) FindModifiedParameters(fnNode rhino.Node) *jutil.LinkedHashSet[string] {
	var paramIDs, modifiedIDs inlineIDSet
	f.findModifiedParametersByIDInto(fnNode, &paramIDs, &modifiedIDs)
	return modifiedIDs.toLinkedHashSet()
}

func (f *FunctionArgumentInjector) findModifiedParametersByIDInto(
	fnNode rhino.Node, dstNames *inlineIDSet, dstModified *inlineIDSet,
) {
	dstModified.clear()
	fn := ast.MustFunction(fnNode)
	params := fn.Params()
	if !params.HasParams() {
		return
	}
	collectFunctionParameterIDs(fnNode, dstNames)
	if !dstNames.isEmpty() {
		collectModifiedParameterIDs(fn.Body().Node(), dstNames, false, dstModified)
	}
}

// FunctionArgumentInjector_FindModifiedParametersIn checks for uses of the
// named value that imply a pass-by-value parameter is expected. This is used
// to prevent cases like:
//
// function (x) { x=2; return x; }
//
// We don't want "undefined" to be substituted for "x", and get undefined=2
//
// n is the node in question; names is the set of names to check;
// inInnerClosure is whether the inspection is occurring inside an inner
// function or non-static class field initializer.
//
// Java: FunctionArgumentInjector.findModifiedParameters(Node, ImmutableSet<String>, boolean)
func FunctionArgumentInjector_FindModifiedParametersIn(
	n rhino.Node, names *jutil.LinkedHashSet[string], inInnerClosure bool) *jutil.LinkedHashSet[string] {
	var nameIDs, modifiedIDs inlineIDSet
	nameIDs.fromLinkedHashSet(names, true)
	if !nameIDs.isEmpty() {
		collectModifiedParameterIDs(n, &nameIDs, inInnerClosure, &modifiedIDs)
	}
	return modifiedIDs.toLinkedHashSet()
}

func collectModifiedParameterIDs(
	n rhino.Node, names *inlineIDSet, inInnerClosure bool, dstModified *inlineIDSet,
) {
	if NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(n) {
		// A function parameter can not be replaced with a direct inlined value
		// if it is referred to by an inner function or a non-static class field
		// initializer. The inner function or field initializer can out live the
		// call we are replacing, so it must capture a unique name.  This approach
		// does not work within loop bodies so those are forbidden elsewhere.
		inInnerClosure = true
	}

	if n.IsName() {
		nameID := n.GetStringID()
		if names.contains(nameID) && (inInnerClosure || FunctionArgumentInjector_CanNameValueChange(n)) {
			dstModified.add(nameID)
		}
	}

	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		collectModifiedParameterIDs(c, names, inInnerClosure, dstModified)
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

// FunctionArgumentInjector_GatherLocalNames gathers any names declared in
// the local scope.
//
// Java: FunctionArgumentInjector.gatherLocalNames(Node, Set<String>)
func FunctionArgumentInjector_GatherLocalNames(n rhino.Node, names *jutil.LinkedHashSet[string]) {
	var localIDs inlineIDSet
	gatherLocalNameIDs(n, &localIDs)
	for _, id := range localIDs.items {
		names.Add(rhino.RhinoStringPool_StringFromID(id))
	}
}

func collectFunctionParameterIDs(fnNode rhino.Node, dst *inlineIDSet) {
	dst.clear()
	for p, more := ast.MustFunction(fnNode).Params().FirstParam(); more; p, more = p.Next() {
		dst.add(functionArgumentInjector_paramID(p))
	}
}

// FunctionArgumentInjector_GetFunctionParameterSet gets a set of function
// parameter names.
//
// Java: FunctionArgumentInjector.getFunctionParameterSet(Node)
func FunctionArgumentInjector_GetFunctionParameterSet(fnNode rhino.Node) *jutil.LinkedHashSet[string] {
	var paramIDs inlineIDSet
	collectFunctionParameterIDs(fnNode, &paramIDs)
	return paramIDs.toLinkedHashSet()
}

// functionArgumentInjector_paramID returns the strpool ID of a parameter name,
// or calleeRestMarkerID for a rest parameter; a default value or a pattern is
// not supported (InlineFunctions does not inline such functions).
func functionArgumentInjector_paramID(p ast.Param) uint32 {
	switch p.Token() {
	case rhino.Token_ITER_REST:
		return calleeRestMarkerID
	case rhino.Token_NAME:
		return p.MustName().Node().GetStringID()
	}
	panic(jutil.NewIllegalStateExceptionMsg("Not supported: " + p.Node().String()))
}
