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

// Go-only: the parameters of a function and the arguments of a call are
// read through the typed view (internal/ast), and so are the targets of
// assignments and the declared names; internal parameter/local sets and
// argument replacement maps are keyed by 32-bit strpool IDs instead of
// hashing materialized strings.

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

// FunctionArgumentInjector is a nifty set of functions to deal with the
// issues of replacing function parameters with a set of call argument
// expressions.
type FunctionArgumentInjector struct {
	astAnalyzer *AstAnalyzer

	sentinelUndefinedNode rhino.Node
	templateUndefinedNode rhino.Node
	sentinelThisNode      rhino.Node
	paramBindings         []callParamBinding
	scratchLocals         inlineIDSet
	referencedAfter       FunctionArgumentInjector_ReferencedAfterSideEffect

	// readOnlySummaries is enabled on parallel Phase 2 workers so lazy side-effect
	// summaries are stored in worker-local sideEffectCache instead of mutating
	// shared FunctionState summaries concurrently.
	readOnlySummaries bool
	sideEffectCache   map[rhino.Node]calleeSideEffectSummary
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

// Java: FunctionArgumentInjector(AstAnalyzer)
func NewFunctionArgumentInjector(astAnalyzer *AstAnalyzer) *FunctionArgumentInjector {
	f := &FunctionArgumentInjector{
		astAnalyzer:           astAnalyzer,
		sentinelUndefinedNode: NodeUtil_NewUndefinedNode(rhino.NilNode),
		templateUndefinedNode: NodeUtil_NewUndefinedNode(rhino.NilNode),
		sentinelThisNode:      rhino.IR_ThisNode(),
	}
	f.referencedAfter.outer = f
	return f
}

func (f *FunctionArgumentInjector) getTemplateUndefinedNode(callNode rhino.Node) rhino.Node {
	if !f.templateUndefinedNode.IsNil() && f.templateUndefinedNode.ResetTemplateUndefinedSrcref(callNode) {
		return f.templateUndefinedNode
	}
	return NodeUtil_NewUndefinedNode(callNode)
}

// Inject, with the map provided, replaces the names with expression trees.
//
// node is the root node of the tree within which to perform the
// substitutions; parent is the parent root node; replacements is the map of
// names to template node trees with which to replace the name Nodes.
// Returns the root node or its replacement.
//
// Java: FunctionArgumentInjector.inject(AbstractCompiler, Node, Node, Map<String, Node>)
func (f *FunctionArgumentInjector) Inject(compiler AbstractCompiler, node rhino.Node, parent rhino.Node, replacements *jutil.LinkedHashMap[string, rhino.Node]) rhino.Node {
	return f.InjectWithReplaceThis(compiler, node, parent, replacements, true /* replaceThis= */)
}

// Java: FunctionArgumentInjector.inject(AbstractCompiler, Node, Node, Map<String, Node>, boolean)
func (f *FunctionArgumentInjector) InjectWithReplaceThis(
	compiler AbstractCompiler,
	node rhino.Node,
	parent rhino.Node,
	replacements *jutil.LinkedHashMap[string, rhino.Node],
	replaceThis bool) rhino.Node {
	var idMap inlineIDNodeMap
	idMap.clear()
	hadThis := false
	if replacements != nil {
		for k, v := range replacements.All() {
			id := rhino.RhinoStringPool_InternID(k)
			if id == calleeThisMarkerID {
				hadThis = true
			}
			idMap.put(id, v)
		}
	}
	res := f.injectByID(compiler, node, parent, &idMap, replaceThis)
	if hadThis && idMap.get(calleeThisMarkerID).IsNil() {
		replacements.Remove(FunctionArgumentInjector_THIS_MARKER)
	}
	return res
}

func (f *FunctionArgumentInjector) injectByID(
	compiler AbstractCompiler,
	node rhino.Node,
	parent rhino.Node,
	replacements *inlineIDNodeMap,
	replaceThis bool) rhino.Node {
	if node.IsName() {
		replacementTemplate := replacements.get(node.GetStringID())
		if !replacementTemplate.IsNil() {
			// This should not be replacing declared names.
			jutil.CheckStateMsg1(!(parent.IsFunction() || parent.IsVar() || parent.IsCatch()), "%s", parent)
			if !parent.IsNil() && node.CanCloneLeafInPlace() && replacementTemplate.CanCloneLeafInPlace() {
				node.OverwriteLeafFromClone(replacementTemplate)
				return node
			}
			// The name may need to be replaced more than once,
			// so we need to clone the node.
			replacement := replacementTemplate.CloneTree()
			node.ReplaceWith(replacement)
			return replacement
		}
	} else if replaceThis && node.IsThis() {
		replacementTemplate := replacements.get(calleeThisMarkerID)
		rhino.CheckNotNull(replacementTemplate)
		if !replacementTemplate.IsThis() {
			if !parent.IsNil() && node.CanCloneLeafInPlace() && replacementTemplate.CanCloneLeafInPlace() {
				node.OverwriteLeafFromClone(replacementTemplate)
				if compiler.GetAstAnalyzer().MayHaveSideEffects(replacementTemplate) {
					replacements.remove(calleeThisMarkerID)
				}
				return node
			}
			// The name may need to be replaced more than once,
			// so we need to clone the node.
			replacement := replacementTemplate.CloneTree()
			node.ReplaceWith(replacement)

			// Remove the value.  This isn't required but it ensures that we won't
			// inject side-effects multiple times as it will trigger the null
			// check above if we do.
			if compiler.GetAstAnalyzer().MayHaveSideEffects(replacementTemplate) {
				replacements.remove(calleeThisMarkerID)
			}

			return replacement
		}
	} else if (node.IsFunction() && !node.IsArrowFunction()) ||
		node.IsMemberFieldDef() ||
		NodeUtil_IsClassStaticBlock(node) {
		// Once we enter another non-arrow function, class field initializer, or static block,
		// the "this" value changes. Don't try to replace it within an inner scope.
		replaceThis = false
	}

	if field, ok := ast.AsComputedFieldDef(node); ok {
		// The key is evaluated where the class is, the value later, in the
		// field's own this.
		f.injectByID(compiler, field.Key().Node(), node, replacements, replaceThis)
		if value, ok := field.Value(); ok { // read after the key was replaced
			f.injectByID(compiler, value.Node(), node, replacements /* replaceThis= */, false)
		}
		return node
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		// We have to reassign c in case it was replaced, because the removed c's
		// getNext() would no longer be correct.
		c = f.injectByID(compiler, c, node, replacements, replaceThis)
	}

	return node
}

type callParamWalkMode uint8

const (
	callParamWalkOwned callParamWalkMode = iota
	callParamWalkTemplate
	callParamWalkReadOnly
)

type callParamBinding struct {
	nameID    uint32
	paramNode rhino.Node
	argNode   rhino.Node
}

func (b callParamBinding) nameStr() string {
	return rhino.RhinoStringPool_StringFromID(b.nameID)
}

func hasDuplicateBindingID(bindings []callParamBinding) bool {
	n := len(bindings)
	if n <= 1 {
		return false
	}
	if n <= inlineSmallCap {
		for i := 1; i < n; i++ {
			id := bindings[i].nameID
			for j := 0; j < i; j++ {
				if bindings[j].nameID == id {
					return true
				}
			}
		}
		return false
	}
	seen := make(map[uint32]struct{}, n)
	for _, b := range bindings {
		if _, dup := seen[b.nameID]; dup {
			return true
		}
		seen[b.nameID] = struct{}{}
	}
	return false
}

func (f *FunctionArgumentInjector) resolveUndefinedArg(
	callNode rhino.Node,
	mode callParamWalkMode,
	undefTemplate *rhino.Node,
) rhino.Node {
	switch mode {
	case callParamWalkOwned:
		return NodeUtil_NewUndefinedNode(callNode)
	case callParamWalkReadOnly:
		return f.sentinelUndefinedNode
	default:
		if undefTemplate.IsNil() {
			*undefTemplate = f.getTemplateUndefinedNode(callNode)
		}
		return *undefTemplate
	}
}

func (f *FunctionArgumentInjector) appendCallParamBindings(
	dst []callParamBinding,
	fnNode rhino.Node,
	callNode rhino.Node,
	safeNameIdSupplier func() string,
	mode callParamWalkMode,
) []callParamBinding {
	rhino.CheckNotNull(fnNode)
	var undefTemplate rhino.Node

	// CALL NODE: [ NAME, ARG1, ARG2, ... ]
	cArg, hasArg := ast.MustNormalOrOptChainCall(callNode).FirstArg()
	if hasArg && NodeUtil_IsFunctionObjectCall(callNode) {
		thisArg := cArg.Node()
		if fnNode.IsArrowFunction() {
			if mode == callParamWalkOwned {
				thisArg = rhino.IR_ThisNode()
			} else {
				thisArg = f.sentinelThisNode
			}
		}
		dst = append(dst, callParamBinding{
			nameID:    calleeThisMarkerID,
			paramNode: FunctionArgumentInjector_THIS_MARKER_NODE,
			argNode:   thisArg,
		})
		cArg, hasArg = cArg.Next()
	} else {
		// 'apply' isn't supported yet.
		jutil.CheckStateMsg1(!NodeUtil_IsFunctionObjectApply(callNode), "%s", callNode)
		dst = append(dst, callParamBinding{
			nameID:    calleeThisMarkerID,
			paramNode: FunctionArgumentInjector_THIS_MARKER_NODE,
			argNode:   f.resolveUndefinedArg(callNode, mode, &undefTemplate),
		})
	}

	// (A plain loop: ranging over Params() with a return in the body costs
	// a closure call per parameter in this hot path.)
	for fnParam, more := ast.MustFunction(fnNode).Params().FirstParam(); more; fnParam, more = fnParam.Next() {
		if rest, ok := fnParam.AsIterRest(); ok {
			restName, ok := rest.Target().AsName()
			if hasArg {
				jutil.CheckStateMsg1(ok, "%s", rest.Target().Node())
				// The rest param takes the remaining arguments, as an array.
				array := rhino.IR_Arraylit()
				array.SrcrefTreeIfMissing(cArg.Node())
				for ; hasArg; cArg, hasArg = cArg.Next() {
					array.AddChildToBack(cArg.Node().CloneTree())
				}
				return append(dst, callParamBinding{
					nameID:    restName.Node().GetStringID(),
					paramNode: restName.Node(),
					argNode:   array,
				})
			}
			jutil.CheckStateMsg1(ok, "%s", fnParam.Node())
			dst = append(dst, callParamBinding{
				nameID:    restName.Node().GetStringID(),
				paramNode: restName.Node(),
				argNode:   rhino.IR_Arraylit().Srcref(fnParam.Node()),
			})
			continue
		}
		name, ok := fnParam.AsName()
		jutil.CheckStateMsg1(ok, "%s", fnParam.Node())
		nameNode := name.Node()
		if hasArg {
			dst = append(dst, callParamBinding{
				nameID:    nameNode.GetStringID(),
				paramNode: nameNode,
				argNode:   cArg.Node(),
			})
			cArg, hasArg = cArg.Next()
		} else {
			dst = append(dst, callParamBinding{
				nameID:    nameNode.GetStringID(),
				paramNode: nameNode,
				argNode:   f.resolveUndefinedArg(callNode, mode, &undefTemplate),
			})
		}
	}

	// Add temp names for arguments that don't have named parameters in the
	// called function.
	for ; hasArg; cArg, hasArg = cArg.Next() {
		uniquePlaceholder := FunctionArgumentInjector_GetUniqueAnonymousParameterName(safeNameIdSupplier)
		if mode == callParamWalkReadOnly {
			dst = append(dst, callParamBinding{
				nameID:    rhino.RhinoStringPool_InternID(uniquePlaceholder),
				paramNode: FunctionArgumentInjector_THIS_MARKER_NODE,
				argNode:   cArg.Node(),
			})
		} else {
			paramNode := rhino.IR_Name(uniquePlaceholder)
			dst = append(dst, callParamBinding{
				nameID:    paramNode.GetStringID(),
				paramNode: paramNode,
				argNode:   cArg.Node(),
			})
		}
	}
	return dst
}

// GetFunctionCallParameterMap gets a mapping for function parameter names to
// call arguments.
//
// Java: FunctionArgumentInjector.getFunctionCallParameterMap(Node, Node, Supplier<String>)
func (f *FunctionArgumentInjector) GetFunctionCallParameterMap(
	fnNode rhino.Node, callNode rhino.Node, safeNameIdSupplier func() string) *jutil.LinkedHashMap[string, FunctionArgumentInjector_ParamArgPair] {
	argMap := jutil.NewImmutableMapBuilder[string, FunctionArgumentInjector_ParamArgPair]()
	for _, b := range f.appendCallParamBindings(nil, fnNode, callNode, safeNameIdSupplier, callParamWalkOwned) {
		argMap.Put(b.nameStr(), NewFunctionArgumentInjector_ParamArgPair(b.paramNode, b.argNode))
	}
	return argMap.BuildOrThrow()
}

func (f *FunctionArgumentInjector) fillCallParamNodeMap(
	fnNode rhino.Node,
	callNode rhino.Node,
	safeNameIdSupplier func() string,
	dst *inlineIDNodeMap,
) {
	dst.clear()
	f.paramBindings = f.appendCallParamBindings(f.paramBindings[:0], fnNode, callNode, safeNameIdSupplier, callParamWalkTemplate)
	for _, b := range f.paramBindings {
		if !dst.put(b.nameID, b.argNode).IsNil() {
			f.GetFunctionCallParameterMap(fnNode, callNode, safeNameIdSupplier)
		}
	}
}

func (f *FunctionArgumentInjector) fillCallParamBindings(
	dst []callParamBinding,
	fnNode rhino.Node,
	callNode rhino.Node,
	safeNameIdSupplier func() string,
) []callParamBinding {
	dst = f.appendCallParamBindings(dst[:0], fnNode, callNode, safeNameIdSupplier, callParamWalkTemplate)
	if hasDuplicateBindingID(dst) {
		f.GetFunctionCallParameterMap(fnNode, callNode, safeNameIdSupplier)
	}
	return dst
}

func (f *FunctionArgumentInjector) getCallParamBindingsScratch(
	fnNode rhino.Node,
	callNode rhino.Node,
	safeNameIdSupplier func() string,
) []callParamBinding {
	f.paramBindings = f.appendCallParamBindings(f.paramBindings[:0], fnNode, callNode, safeNameIdSupplier, callParamWalkReadOnly)
	if hasDuplicateBindingID(f.paramBindings) {
		f.GetFunctionCallParameterMap(fnNode, callNode, safeNameIdSupplier)
	}
	return f.paramBindings
}

// FunctionArgumentInjector_GetUniqueAnonymousParameterName: parameter names
// will be name unique when at a later time.
//
// Java: FunctionArgumentInjector.getUniqueAnonymousParameterName(Supplier<String>)
func FunctionArgumentInjector_GetUniqueAnonymousParameterName(safeNameIdSupplier func() string) string {
	return "JSCompiler_inline_anon_param_" + safeNameIdSupplier()
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

// FunctionArgumentInjector_CanNameValueChange is similar to NodeUtil.isLValue
// except that object properties and array member modification aren't
// important ("o" in "o.a = 2" is still "o" after assignment, where in as "o
// = x", "o" is now "x").
//
// This also looks for the redefinition of a name. function (x) {var x;}
//
// n is the NAME node in question.
//
// Java: FunctionArgumentInjector.canNameValueChange(Node)
func FunctionArgumentInjector_CanNameValueChange(n rhino.Node) bool {
	if !NodeUtil_IsLValue(n) {
		return false
	}
	if NodeUtil_IsDeclarationLValue(n) {
		declaring := NodeUtil_GetDeclaringParent(n)
		return !declaring.IsConst() && !declaring.IsLet()
	}
	return true
}

func bindingsFromArgMap(dst []callParamBinding, argMap *jutil.LinkedHashMap[string, FunctionArgumentInjector_ParamArgPair]) []callParamBinding {
	dst = dst[:0]
	if argMap == nil || argMap.IsEmpty() {
		return dst
	}
	for k, pair := range argMap.All() {
		dst = append(dst, callParamBinding{
			nameID:    rhino.RhinoStringPool_InternID(k),
			paramNode: pair.ParamNode(),
			argNode:   pair.Arg(),
		})
	}
	return dst
}

func findThisArgInBindings(bindings []callParamBinding) rhino.Node {
	if len(bindings) > 0 && bindings[0].nameID == calleeThisMarkerID {
		return bindings[0].argNode
	}
	for _, b := range bindings {
		if b.nameID == calleeThisMarkerID {
			return b.argNode
		}
	}
	return rhino.NilNode
}

// findRequiresTempsUpToParameterID returns the strpool ID of the last parameter
// whose call argument requires hoisting into a temporary (or 0 if none do).
//
// If an argument requires a temp (e.g. it is side-effectful, affects mutable
// state or due to other reasons), we need to evaluate it into a temporary at
// the callsite before inlining the function block into that callsite.
// However, that argument's early evaluation could also change the values of
// previous arguments. Hence, in this function we decide to hoist all those
// previous args and evaluate them before this side-effectful argument, even
// if they are not actually affected by that side-effectful argument.
func (f *FunctionArgumentInjector) findRequiresTempsUpToParameterID(
	compiler AbstractCompiler,
	fnNode rhino.Node,
	bindings []callParamBinding,
	modifiedParameters *inlineIDSet,
	convention CodingConvention,
	summary *calleeSummary,
	renamedClone bool,
) uint32 {
	useSummary := summary != nil &&
		(renamedClone || summary.fnNode == fnNode) &&
		summary.canReuseAfterRename(renamedClone)
	if !useSummary && len(bindings) > 1 {
		var localSummary calleeSummary
		if scanCalleeSummary("", fnNode, &localSummary) {
			summary = &localSummary
			renamedClone = false
			useSummary = true
		}
	}

	block := ast.MustFunction(fnNode).Body().Node()
	var requiresTempsUpToID uint32

	argCount := len(bindings)
	var isTrivialBody bool
	if useSummary {
		isTrivialBody = summary.isTrivialBody
	} else if argCount > 1 {
		isTrivialBody = !block.HasChildren() ||
			(block.HasOneChild() && !f.BodyMayHaveConditionalCode(block.GetLastChild()))
	}
	thisArg := findThisArgInBindings(bindings)
	hasMinimalParameters := (NodeUtil_IsUndefined(thisArg) ||
		(fnNode.IsArrowFunction() && thisArg.IsThis())) &&
		argCount <= 2 // this + one parameter

	useSideEffectSummary := useSummary && summary.canReuseSideEffects(renamedClone)
	var seSummary calleeSideEffectSummary
	var seSummaryReady bool
	var sideEffectsCollected bool

	isReferencedAfterSideEffect := func(paramID uint32, paramIdx int) bool {
		if useSideEffectSummary {
			if !seSummaryReady {
				var ok bool
				seSummary, ok = f.ensureCalleeSideEffects(summary)
				if ok {
					seSummaryReady = true
				} else {
					useSideEffectSummary = false
				}
			}
			if seSummaryReady {
				if paramID == calleeThisMarkerID {
					return seSummary.thisAfterSideEffect
				}
				if paramIdx < 0 || paramIdx >= int(summary.numParams) {
					return false
				}
				return (seSummary.afterSideEffectMask & (uint64(1) << uint(paramIdx))) != 0
			}
		}
		if !sideEffectsCollected {
			f.collectReferencedAfterSideEffectForBindings(bindings, block)
			sideEffectsCollected = true
		}
		return f.referencedAfter.parametersReferenced.contains(paramID)
	}

	hasNamesReferencedBefore := func(paramID uint32, paramIdx int) bool {
		if useSideEffectSummary {
			if !seSummaryReady {
				var ok bool
				seSummary, ok = f.ensureCalleeSideEffects(summary)
				if ok {
					seSummaryReady = true
				} else {
					useSideEffectSummary = false
				}
			}
			if seSummaryReady {
				if paramIdx < 0 || paramIdx >= int(summary.numParams) {
					return false
				}
				return (seSummary.namesBeforeParamMask & (uint64(1) << uint(paramIdx))) != 0
			}
		}
		if !sideEffectsCollected {
			f.collectReferencedAfterSideEffectForBindings(bindings, block)
			sideEffectsCollected = true
		}
		return f.referencedAfter.parametersWithNamesReferencedBeforeParameter.contains(paramID)
	}

	astAnalyzer := compiler.GetAstAnalyzer()
	var summaryParams []calleeParamSummary
	if useSummary {
		summaryParams = summary.params()
	}

	// Check for arguments that are evaluated more than once.
	for entryIdx, b := range bindings {
		paramID := b.nameID
		paramIdx := -1
		if useSummary {
			if renamedClone {
				if entryIdx > 0 && entryIdx-1 < len(summaryParams) {
					paramIdx = entryIdx - 1
				}
			} else if paramID != calleeThisMarkerID {
				paramIdx = summary.findParamByID(paramID)
			}
		}

		cArg := b.argNode
		if !modifiedParameters.isEmpty() && modifiedParameters.contains(paramID) {
			requiresTempsUpToID = paramID
			continue
		}
		// Stores whether this arg need to get hoisted using a temporary
		requiresTemporary := false
		var references int
		if useSummary {
			if paramIdx >= 0 {
				references = int(summaryParams[paramIdx].refCount)
			}
		} else if paramID != calleeThisMarkerID {
			references = NodeUtil_GetNameReferenceCount(block, rhino.RhinoStringPool_StringFromID(paramID))
		}

		argSideEffects := astAnalyzer.MayHaveSideEffects(cArg)
		if !argSideEffects && references == 0 {
			requiresTemporary = false
		} else if isTrivialBody &&
			hasMinimalParameters &&
			references == 1 &&
			// The below line is checking: Can this be affected by side-effects in the function body?
			!(NodeUtil_CanBeSideEffected(cArg) && isReferencedAfterSideEffect(paramID, paramIdx)) &&
			// The below line is checking: Can the function body be affected by the side effects of
			// the argument? If the argument has side effects and the function body contains names
			// that are referenced before the argument, then we need to inline using temporaries. If
			// not, then we can enter this condition's else-if block.
			(!argSideEffects || !hasNamesReferencedBefore(paramID, paramIdx)) {
			// For functions with a trivial body, and where the parameter evaluation order
			// can't change, and there aren't any side-effect before the parameter, we can
			// avoid creating a temporary.
			//
			// This is done to help inline common trivial functions
			// TODO: b/407603216 - return `true` when return expressions don't have
			// references that can be side-effected
			requiresTemporary = false
		} else if astAnalyzer.MayEffectMutableState(cArg) && references > 0 {
			// Note: Mutable arguments should be assigned to temps, as the
			// may be within in a loop:
			//   function x(a) {
			//     for(var i=0; i<0; i++) {
			//       foo(a);
			//     }
			//   x( [] );
			// The parameter in the call to foo should not become "[]".
			requiresTemporary = true
		} else if argSideEffects {
			// even if there are no references, we still need to evaluate the
			// expression if it has side-effects.
			requiresTemporary = true
		} else if NodeUtil_CanBeSideEffected(cArg) && isReferencedAfterSideEffect(paramID, paramIdx) {
			requiresTemporary = true
		} else if references > 1 {
			// Safe is a misnomer, this is a check for "large".
			switch cArg.GetToken() {
			case rhino.Token_NAME:
				name := cArg.GetString()
				// Don't worry about whether this is global or local, just check if it is
				// "exported" in either case.
				requiresTemporary =
					(convention.IsExportedWithLocal(name, true) || convention.IsExportedWithLocal(name, false))
			case rhino.Token_THIS:
				requiresTemporary = false
			case rhino.Token_STRINGLIT:
				requiresTemporary = (jutil.StringLength(cArg.GetString()) >= 2)
			default:
				requiresTemporary = !NodeUtil_IsImmutableValue(cArg)
			}
		}

		if requiresTemporary {
			requiresTempsUpToID = paramID
		}
	}

	return requiresTempsUpToID
}

func (f *FunctionArgumentInjector) hasCallArgumentsNeedingTempsForBindings(
	compiler AbstractCompiler,
	fnNode rhino.Node,
	bindings []callParamBinding,
	modifiedParameters *inlineIDSet,
	convention CodingConvention,
	summary *calleeSummary,
) bool {
	jutil.CheckArgumentMsg1(fnNode.IsFunction(), "%s", fnNode)
	if !modifiedParameters.isEmpty() {
		return true
	}
	if len(bindings) == 0 {
		return false
	}
	requiresTempsUpToID := f.findRequiresTempsUpToParameterID(
		compiler, fnNode, bindings, modifiedParameters, convention, summary, false)
	if requiresTempsUpToID == 0 {
		return false
	}
	thisArg := findThisArgInBindings(bindings)
	for _, b := range bindings {
		if b.nameID == calleeThisMarkerID &&
			(NodeUtil_IsUndefined(thisArg) ||
				(fnNode.IsArrowFunction() && thisArg.IsThis())) {
			continue
		}
		return true
	}
	return false
}

// Java: FunctionArgumentInjector.gatherCallArgumentsNeedingTemps(AbstractCompiler, Node, ImmutableMap<String, ParamArgPair>, ImmutableSet<String>, CodingConvention)
func (f *FunctionArgumentInjector) GatherCallArgumentsNeedingTemps(
	compiler AbstractCompiler,
	fnNode rhino.Node,
	argMap *jutil.LinkedHashMap[string, FunctionArgumentInjector_ParamArgPair],
	modifiedParameters *jutil.LinkedHashSet[string],
	convention CodingConvention) *jutil.LinkedHashSet[string] {
	jutil.CheckArgumentMsg1(fnNode.IsFunction(), "%s", fnNode)
	if argMap.IsEmpty() {
		return modifiedParameters
	}
	var bindingsBuf [inlineSmallCap]callParamBinding
	bindings := bindingsFromArgMap(bindingsBuf[:0], argMap)
	var modifiedIDs inlineIDSet
	modifiedIDs.fromLinkedHashSet(modifiedParameters, true)
	requiresTempsUpToID := f.findRequiresTempsUpToParameterID(
		compiler, fnNode, bindings, &modifiedIDs, convention, nil, false)
	if requiresTempsUpToID == 0 {
		return modifiedParameters
	}
	var dstIDs inlineIDSet
	dstIDs.addAll(&modifiedIDs)
	addCallArgumentsNeedingTempsByID(fnNode, bindings, requiresTempsUpToID, &dstIDs)
	return dstIDs.toLinkedHashSet()
}

func (f *FunctionArgumentInjector) gatherCallArgumentsNeedingTempsByID(
	compiler AbstractCompiler,
	fnNode rhino.Node,
	bindings []callParamBinding,
	modifiedParameters *inlineIDSet,
	convention CodingConvention,
	summary *calleeSummary,
	renamedClone bool,
	dst *inlineIDSet,
) {
	jutil.CheckArgumentMsg1(fnNode.IsFunction(), "%s", fnNode)
	dst.clear()
	dst.addAll(modifiedParameters)
	if len(bindings) == 0 {
		return
	}
	requiresTempsUpToID := f.findRequiresTempsUpToParameterID(
		compiler, fnNode, bindings, modifiedParameters, convention, summary, renamedClone)
	if requiresTempsUpToID == 0 {
		return
	}
	addCallArgumentsNeedingTempsByID(fnNode, bindings, requiresTempsUpToID, dst)
}

func addCallArgumentsNeedingTempsByID(
	fnNode rhino.Node,
	bindings []callParamBinding,
	requiresTempsUpToID uint32,
	dst *inlineIDSet,
) {
	thisArg := findThisArgInBindings(bindings)
	// mark all names up to requiresTempsUpToID as namesNeedingTemps
	for _, b := range bindings {
		if b.nameID == calleeThisMarkerID &&
			(NodeUtil_IsUndefined(thisArg) ||
				(fnNode.IsArrowFunction() && thisArg.IsThis())) {
			// Skip generating a temporary for an implicit "this" arg.
			continue
		}
		dst.add(b.nameID)
		if b.nameID == requiresTempsUpToID {
			break
		}
	}
}

// BodyMayHaveConditionalCode: we consider a return or expression trivial if
// it doesn't contain a conditional expression, a function, or a non-static
// class field initializer.
//
// Java: FunctionArgumentInjector.bodyMayHaveConditionalCode(Node)
func (f *FunctionArgumentInjector) BodyMayHaveConditionalCode(n rhino.Node) bool {
	if !n.IsReturn() && !n.IsExprResult() {
		return true
	}
	return f.MayHaveConditionalCode(n)
}

// MayHaveConditionalCode: we consider an expression trivial if it doesn't
// contain a conditional expression, a function, or a non-static class field
// initializer.
//
// Java: FunctionArgumentInjector.mayHaveConditionalCode(Node)
func (f *FunctionArgumentInjector) MayHaveConditionalCode(n rhino.Node) bool {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_FUNCTION,
			rhino.Token_AND,
			rhino.Token_OR,
			rhino.Token_HOOK,
			rhino.Token_COALESCE,
			rhino.Token_OPTCHAIN_CALL,
			rhino.Token_OPTCHAIN_GETELEM,
			rhino.Token_OPTCHAIN_GETPROP:
			return true
		default:
			if NodeUtil_IsNonStaticClassFieldInitializer(c) {
				return true
			}
		}
		if f.MayHaveConditionalCode(c) {
			return true
		}
	}
	return false
}

// FindParametersReferencedAfterSideEffect bootstraps a traversal to look for
// parameters referenced after a non-local side-effect, and parameters with
// names referenced before the parameter is referenced.
//
// NOTE: This assumes no-inner functions.
//
// parameters is the set of parameter names; root is the function code
// block. Returns two sets: Set #1 is the subset of parameters referenced
// after the first seen non-local side-effect. Set #2 is the subset of
// parameters with names referenced before the parameter is referenced.
//
// Java: FunctionArgumentInjector.findParametersReferencedAfterSideEffect(ImmutableSet<String>, Node)
func (f *FunctionArgumentInjector) FindParametersReferencedAfterSideEffect(
	parameters *jutil.LinkedHashSet[string], root rhino.Node) *FunctionArgumentInjector_SetContainer {
	f.scratchLocals.clear()
	gatherLocalNameIDs(root, &f.scratchLocals)

	collector := &f.referencedAfter
	collector.reset(f, nil, &f.scratchLocals)
	collector.parameterIDs.fromLinkedHashSet(parameters, true)
	NodeUtil_VisitPostOrderWithPred(root, collector, collector.Apply)
	return NewFunctionArgumentInjector_SetContainer(
		collector.parametersReferenced.toLinkedHashSet(),
		collector.parametersWithNamesReferencedBeforeParameter.toLinkedHashSet(),
	)
}

func (f *FunctionArgumentInjector) collectReferencedAfterSideEffectForBindings(
	bindings []callParamBinding, root rhino.Node,
) {
	f.scratchLocals.clear()
	gatherLocalNameIDs(root, &f.scratchLocals)

	collector := &f.referencedAfter
	collector.reset(f, nil, &f.scratchLocals)
	for _, b := range bindings {
		collector.parameterIDs.add(b.nameID)
	}
	NodeUtil_VisitPostOrderWithPred(root, collector, collector.Apply)
}

func (f *FunctionArgumentInjector) ensureCalleeSideEffects(summary *calleeSummary) (calleeSideEffectSummary, bool) {
	if summary == nil || int(summary.numParams) > 64 {
		return calleeSideEffectSummary{overflow: true}, false
	}
	if summary.sideEffects.computed {
		return summary.sideEffects, true
	}
	if cached, ok := f.sideEffectCache[summary.fnNode]; ok {
		if !f.readOnlySummaries {
			summary.sideEffects = cached
		}
		return cached, true
	}

	root := ast.MustFunction(summary.fnNode).Body().Node()
	f.scratchLocals.clear()
	gatherLocalNameIDs(root, &f.scratchLocals)

	collector := &f.referencedAfter
	collector.reset(f, summary, &f.scratchLocals)
	NodeUtil_VisitPostOrderWithPred(root, collector, collector.Apply)
	collector.summary = nil

	se := calleeSideEffectSummary{
		computed:            true,
		thisAfterSideEffect: collector.parametersReferenced.contains(calleeThisMarkerID),
	}
	for i, p := range summary.params() {
		bit := uint64(1) << uint(i)
		if collector.parametersReferenced.contains(p.nameID) {
			se.afterSideEffectMask |= bit
		}
		if collector.parametersWithNamesReferencedBeforeParameter.contains(p.nameID) {
			se.namesBeforeParamMask |= bit
		}
	}
	if f.readOnlySummaries {
		if f.sideEffectCache == nil {
			f.sideEffectCache = make(map[rhino.Node]calleeSideEffectSummary, 16)
		}
		f.sideEffectCache[summary.fnNode] = se
	} else {
		summary.sideEffects = se
	}
	return se, true
}

// FunctionArgumentInjector_SetContainer is the record SetContainer.
type FunctionArgumentInjector_SetContainer struct {
	parametersReferencedAfterSideEffect          *jutil.LinkedHashSet[string]
	parametersWithNamesReferencedBeforeParameter *jutil.LinkedHashSet[string]
}

// Java: FunctionArgumentInjector.SetContainer(ImmutableSet<String>, ImmutableSet<String>)
func NewFunctionArgumentInjector_SetContainer(
	parametersReferencedAfterSideEffect *jutil.LinkedHashSet[string],
	parametersWithNamesReferencedBeforeParameter *jutil.LinkedHashSet[string]) *FunctionArgumentInjector_SetContainer {
	return &FunctionArgumentInjector_SetContainer{
		parametersReferencedAfterSideEffect:          parametersReferencedAfterSideEffect,
		parametersWithNamesReferencedBeforeParameter: parametersWithNamesReferencedBeforeParameter,
	}
}

// Java: FunctionArgumentInjector.SetContainer.parametersReferencedAfterSideEffect()
func (s *FunctionArgumentInjector_SetContainer) ParametersReferencedAfterSideEffect() *jutil.LinkedHashSet[string] {
	return s.parametersReferencedAfterSideEffect
}

// Java: FunctionArgumentInjector.SetContainer.parametersWithNamesReferencedBeforeParameter()
func (s *FunctionArgumentInjector_SetContainer) ParametersWithNamesReferencedBeforeParameter() *jutil.LinkedHashSet[string] {
	return s.parametersWithNamesReferencedBeforeParameter
}

// FunctionArgumentInjector_ReferencedAfterSideEffect collects parameter
// names referenced after a non-local side-effect.
//
// Assumptions:
//
//   - We assume parameters are not modified in the function body (that is
//     checked separately).
//   - There are no inner functions (also checked separately).
//
// As we are trying to replace parameters with there passed in values we are
// interested in anything that may affect those value. So, ignoring changes
// to local variables, we look for things that may affect anything outside
// the local-state. Once such a side-effect is seen any following reference
// to the function parameters are collected. These will need to be assigned
// to temporaries to prevent changes to their value as would have happened
// during the function call.
//
// To properly handle loop structures all references to the function
// parameters are recorded and the decision to keep or throw away those
// references is deferred until exiting the loop structure.
type FunctionArgumentInjector_ReferencedAfterSideEffect struct {
	outer                *FunctionArgumentInjector
	parameterIDs         inlineIDSet
	summary              *calleeSummary
	locals               *inlineIDSet
	ownedLocals          inlineIDSet
	sideEffectSeen       bool
	parametersReferenced inlineIDSet

	nameNodeHasBeenSeen                          bool
	parametersWithNamesReferencedBeforeParameter inlineIDSet

	loopsEntered int
}

func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) reset(
	outer *FunctionArgumentInjector,
	summary *calleeSummary,
	locals *inlineIDSet,
) {
	r.outer = outer
	r.parameterIDs.clear()
	r.summary = summary
	r.locals = locals
	r.sideEffectSeen = false
	r.nameNodeHasBeenSeen = false
	r.loopsEntered = 0
	r.parametersReferenced.clear()
	r.parametersWithNamesReferencedBeforeParameter.clear()
}

func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) paramCount() int {
	if r.summary != nil {
		return int(r.summary.numParams) + 1
	}
	return r.parameterIDs.size()
}

func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) hasParameterID(nameID uint32) bool {
	if r.summary != nil {
		return r.summary.findParamByID(nameID) >= 0
	}
	return r.parameterIDs.contains(nameID)
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect(ImmutableSet<String>, ImmutableSet<String>)
func NewFunctionArgumentInjector_ReferencedAfterSideEffect(
	outer *FunctionArgumentInjector, parameters *jutil.LinkedHashSet[string], locals *jutil.LinkedHashSet[string]) *FunctionArgumentInjector_ReferencedAfterSideEffect {
	r := &FunctionArgumentInjector_ReferencedAfterSideEffect{
		outer: outer,
	}
	r.parameterIDs.fromLinkedHashSet(parameters, true)
	r.ownedLocals.fromLinkedHashSet(locals, true)
	r.locals = &r.ownedLocals
	return r
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.getParametersReferencedAfterSideEffect()
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) GetParametersReferencedAfterSideEffect() *jutil.LinkedHashSet[string] {
	return r.parametersReferenced.toLinkedHashSet()
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.getParametersWithNamesReferencedBeforeParameter()
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) GetParametersWithNamesReferencedBeforeParameter() *jutil.LinkedHashSet[string] {
	return r.parametersWithNamesReferencedBeforeParameter.toLinkedHashSet()
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.apply(Node)
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) Apply(node rhino.Node) bool {
	// Keep track of any loop structures entered.
	if NodeUtil_IsLoopStructure(node) {
		r.loopsEntered++
	}

	// If we have found all the parameters, don't bother looking
	// at the children.
	return !(r.sideEffectSeen && r.paramCount() == r.parametersReferenced.size())
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.inLoop()
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) InLoop() bool {
	return r.loopsEntered != 0
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.visit(Node)
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) Visit(n rhino.Node) {
	if n.IsName() {
		nameID := n.GetStringID()
		if r.hasParameterID(nameID) && r.nameNodeHasBeenSeen {
			// We have seen a name node before this parameter. If this parameter has side effects, it
			// is not safe to inline.
			r.parametersWithNamesReferencedBeforeParameter.add(nameID)
		}
		r.nameNodeHasBeenSeen = true
	}

	// If we are exiting a loop.
	if NodeUtil_IsLoopStructure(n) {
		r.loopsEntered--
		if !r.InLoop() && !r.sideEffectSeen {
			// Now that the loops has been fully traversed and
			// no side-effects have been seen, throw away
			// the references seen in them.
			r.parametersReferenced.clear()
		}
	}

	if !r.sideEffectSeen {
		// Look for side-effects.
		if r.HasNonLocalSideEffect(n) {
			r.sideEffectSeen = true
		}
	}

	// If traversing the nodes of a loop save any references
	// that are seen.
	if r.InLoop() || r.sideEffectSeen {
		// Record references to parameters.
		if n.IsName() {
			nameID := n.GetStringID()
			if r.hasParameterID(nameID) {
				r.parametersReferenced.add(nameID)
			}
		} else if n.IsThis() {
			r.parametersReferenced.add(calleeThisMarkerID)
		}
	}
}

// HasNonLocalSideEffect returns whether the node may have non-local
// side-effects.
//
// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.hasNonLocalSideEffect(Node)
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) HasNonLocalSideEffect(n rhino.Node) bool {
	sideEffect := false
	type_ := n.GetToken()
	// Note: Only care about changes to non-local names, specifically
	// ignore VAR declaration assignments.
	if assign, ok := ast.AsAssignment(n); ok {
		// Ignore changes to local names.
		if !r.IsLocalName(assign.Target().Node()) {
			sideEffect = true
		}
	} else if incDec, ok := ast.AsIncDec(n); ok {
		if !r.IsLocalName(incDec.Target().Node()) {
			sideEffect = true
		}
	} else if type_ == rhino.Token_CALL {
		sideEffect = r.outer.astAnalyzer.FunctionCallHasSideEffects(n)
	} else if type_ == rhino.Token_NEW {
		sideEffect = r.outer.astAnalyzer.ConstructorCallHasSideEffects(n)
	} else if type_ == rhino.Token_DELPROP {
		sideEffect = true
	}

	return sideEffect
}

// IsLocalName returns whether node is a reference to locally declared name.
//
// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.isLocalName(Node)
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) IsLocalName(node rhino.Node) bool {
	if node.IsName() {
		nameID := node.GetStringID()
		return r.hasParameterID(nameID) || r.locals.contains(nameID)
	}
	return false
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

// FunctionArgumentInjector_ParamArgPair stores a pair of function formal
// parameter + argument value at a specific call site.
//
// paramNode is the original formal parameter name, such as x in function
// foo(x) { ... } or, when injecting a value for this,
// FunctionArgumentInjector#THIS_MARKER_NODE; arg is the argument passed as
// `paramNode` at some call site, such as 500 in f(500);.
type FunctionArgumentInjector_ParamArgPair struct {
	paramNode rhino.Node
	arg       rhino.Node
}

// Java: FunctionArgumentInjector.ParamArgPair(Node, Node)
func NewFunctionArgumentInjector_ParamArgPair(paramNode rhino.Node, arg rhino.Node) FunctionArgumentInjector_ParamArgPair {
	rhino.CheckNotNull(paramNode)
	rhino.CheckNotNull(arg)
	return FunctionArgumentInjector_ParamArgPair{paramNode: paramNode, arg: arg}
}

// Java: FunctionArgumentInjector.ParamArgPair.paramNode()
func (p FunctionArgumentInjector_ParamArgPair) ParamNode() rhino.Node { return p.paramNode }

// Java: FunctionArgumentInjector.ParamArgPair.arg()
func (p FunctionArgumentInjector_ParamArgPair) Arg() rhino.Node { return p.arg }
