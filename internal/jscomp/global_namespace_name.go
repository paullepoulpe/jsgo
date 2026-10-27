// Port of com/google/javascript/jscomp/GlobalNamespace.java @ closure-compiler v20261021

/*
 * Copyright 2006 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// GlobalNamespace_Name is a name defined in global scope (e.g. "a" or
// "a.b.c.d").
//
// Instances form a tree describing the "Closure namespaces" in the program.
// As the parse tree traversal proceeds, we'll discover that some names
// correspond to JavaScript objects whose properties we should consider
// collapsing.
type GlobalNamespace_Name struct {
	outer *GlobalNamespace

	// Flat uint32 handles (0 = nil). Keeping outer as the sole pointer at
	// offset 0 shrinks Name from 160B (160 ptrBytes) to 96B (8 ptrBytes).
	selfID           uint32
	baseNameID       uint32
	parentID         uint32
	firstPropID      uint32
	lastPropID       uint32
	nextPropID       uint32
	declarationID    uint32
	initializationID uint32

	// Inline ref storage (paired with gnNamePropRefsMapState in propertyBitSet):
	//   ref0 == 0: 0 refs
	//   ref0 & gnMultiRefsBit == 0 && !mapState: singleton ref0
	//   ref0 & gnMultiRefsBit == 0 && mapState: 1..2 inline map refs (ref0, optional ref1)
	//   ref0 & gnMultiRefsBit != 0: outer.multiRefs[ref0 &^ gnMultiRefsBit]
	ref0 uint32
	ref1 uint32

	firstDeclDocID  uint32
	firstQnameDocID uint32

	globalSets              int32
	localSets               int32
	localSetsWithNoCollapse int32
	aliasingGets            int32
	totalGets               int32
	callGets                int32
	deleteProps             int32
	subclassingGets         int32

	// Bitset containing NameProps
	//
	// Using a bit set over boolean fields is a memory optimization. Large
	// projects can have millions of Name objects, so saving a few bytes per
	// Name is useful.
	propertyBitSet int32
}

// Java: GlobalNamespace.Name(String, Name, SourceKind)
func NewGlobalNamespace_Name(outer *GlobalNamespace, name string, parent *GlobalNamespace_Name, sourceKind GlobalNamespace_SourceKind) *GlobalNamespace_Name {
	if outer == nil {
		if parent != nil && parent.outer != nil {
			outer = parent.outer
		} else {
			outer = &GlobalNamespace{
				globalRoot: globalNamespace_globalRoot,
			}
		}
	}
	var parentID uint32
	if parent != nil {
		parentID = parent.selfID
	}
	return outer.allocName(uint32(strpool.Intern(name)), parentID, sourceKind)
}

// Java: GlobalNamespace.Name.addProperty(String, SourceKind)
func (n *GlobalNamespace_Name) AddProperty(name string, sourceKind GlobalNamespace_SourceKind) *GlobalNamespace_Name {
	baseID := uint32(strpool.Intern(name))
	node := n.addPropertyWithBaseNameID(baseID, sourceKind)
	n.outer.putChildID(n.selfID, baseID, node.selfID)
	return node
}

// addPropertyWithBaseNameID appends a child property Name in O(1) with zero
// heap slice or wrapper allocations.
func (n *GlobalNamespace_Name) addPropertyWithBaseNameID(baseNameID uint32, sourceKind GlobalNamespace_SourceKind) *GlobalNamespace_Name {
	g := n.outer
	node := g.allocName(baseNameID, n.selfID, sourceKind)
	if n.firstPropID == 0 {
		n.firstPropID = node.selfID
		n.lastPropID = node.selfID
	} else {
		g.nameByID(n.lastPropID).nextPropID = node.selfID
		n.lastPropID = node.selfID
	}
	return node
}

// LookupChild returns the child property Name with baseName propName in O(1),
// or nil if absent.
func (n *GlobalNamespace_Name) LookupChild(propName string) *GlobalNamespace_Name {
	if n.firstPropID == 0 {
		return nil
	}
	propID := uint32(strpool.Lookup(propName))
	if propID == 0 {
		return nil
	}
	return n.outer.nameByID(n.outer.lookupChildID(n.selfID, propID))
}

// RootBaseName returns the base name of the top-level ancestor of n with zero
// allocations.
func (n *GlobalNamespace_Name) RootBaseName() string {
	cur := n
	g := n.outer
	for cur.parentID != 0 {
		cur = g.nameByID(cur.parentID)
	}
	return strpool.String(strpool.ID(cur.baseNameID))
}

// MatchesNode reports whether node is a simple or qualified name matching n's
// full name without allocating n.GetFullName().
func (n *GlobalNamespace_Name) MatchesNode(node rhino.Node) bool {
	if n.parentID == 0 && strings.IndexByte(strpool.String(strpool.ID(n.baseNameID)), '.') >= 0 {
		return node.MatchesQualifiedName(n.GetFullName())
	}
	cur := n
	g := n.outer
	for {
		switch node.GetToken() {
		case rhino.Token_NAME, rhino.Token_IMPORT_STAR:
			if cur.parentID != 0 {
				return false
			}
			id := node.GetStringID()
			return id > 1 && id == cur.baseNameID && strings.IndexByte(strpool.String(strpool.ID(id)), '.') < 0
		case rhino.Token_THIS:
			return cur.parentID == 0 && cur.baseNameID == gnThisStrID
		case rhino.Token_SUPER:
			return cur.parentID == 0 && cur.baseNameID == gnSuperStrID
		case rhino.Token_GETPROP:
			if cur.parentID == 0 {
				return false
			}
			id := node.GetStringID()
			if id != cur.baseNameID || strings.IndexByte(strpool.String(strpool.ID(id)), '.') >= 0 {
				return false
			}
			cur = g.nameByID(cur.parentID)
			node = node.GetFirstChild()
		default:
			return false
		}
	}
}

// HasProps reports whether this Name has any child property Names.
func (n *GlobalNamespace_Name) HasProps() bool {
	return n.firstPropID != 0
}

// Props iterates child property Names present at the start of iteration.
func (n *GlobalNamespace_Name) Props() iter.Seq[*GlobalNamespace_Name] {
	return func(yield func(*GlobalNamespace_Name) bool) {
		endID := n.lastPropID
		if endID == 0 {
			return
		}
		g := n.outer
		for curID := n.firstPropID; curID != 0; {
			p := g.nameByID(curID)
			nextID := p.nextPropID
			if !yield(p) {
				return
			}
			if curID == endID {
				return
			}
			curID = nextID
		}
	}
}

// AppendPropsTo appends all child property Names to *dst.
func (n *GlobalNamespace_Name) AppendPropsTo(dst *[]*GlobalNamespace_Name) {
	if n.firstPropID == 0 {
		return
	}
	g := n.outer
	for curID := n.firstPropID; curID != 0; {
		p := g.nameByID(curID)
		*dst = append(*dst, p)
		curID = p.nextPropID
	}
}

// PropsSlice returns a newly allocated slice of child property Names (or nil
// if there are no child properties).
func (n *GlobalNamespace_Name) PropsSlice() []*GlobalNamespace_Name {
	if n.firstPropID == 0 {
		return nil
	}
	var out []*GlobalNamespace_Name
	n.AppendPropsTo(&out)
	return out
}

// Java: GlobalNamespace.Name.getBaseName()
func (n *GlobalNamespace_Name) GetBaseName() string {
	return strpool.String(strpool.ID(n.baseNameID))
}

// Java: GlobalNamespace.Name.inExterns()
func (n *GlobalNamespace_Name) InExterns() bool {
	return n.GetBooleanProperty(GlobalNamespace_NameProp_IS_EXTERN)
}

// Java: GlobalNamespace.Name.subclassingGetCount()
func (n *GlobalNamespace_Name) SubclassingGetCount() int {
	return int(n.subclassingGets)
}

// Java: GlobalNamespace.Name.getName()
func (n *GlobalNamespace_Name) GetName() string {
	return n.GetFullName()
}

// Java: GlobalNamespace.Name.getFullName()
func (n *GlobalNamespace_Name) GetFullName() string {
	if n.parentID == 0 {
		return strpool.String(strpool.ID(n.baseNameID))
	}
	var segsBuf [32]string
	segs := segsBuf[:0]
	totalLen := -1
	g := n.outer
	for cur := n; cur != nil; {
		s := strpool.String(strpool.ID(cur.baseNameID))
		segs = append(segs, s)
		totalLen += 1 + len(s)
		if cur.parentID == 0 {
			break
		}
		cur = g.nameByID(cur.parentID)
	}
	var sb strings.Builder
	sb.Grow(totalLen)
	for i := len(segs) - 1; i >= 0; i-- {
		if i < len(segs)-1 {
			sb.WriteByte('.')
		}
		sb.WriteString(segs[i])
	}
	return sb.String()
}

// Java: GlobalNamespace.Name.usesHasOwnProperty()
func (n *GlobalNamespace_Name) UsesHasOwnProperty() bool {
	return n.GetBooleanProperty(GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY)
}

// GetDeclaration returns the declaration Ref (a rhino.StaticRef), or nil.
//
// Java: GlobalNamespace.Name.getDeclaration()
func (n *GlobalNamespace_Name) GetDeclaration() rhino.StaticRef {
	if n.declarationID == 0 {
		return nil
	}
	return n.outer.refByID(n.declarationID)
}

// getDeclarationRef is getDeclaration with its Java return type Ref.
//
// Go-only: language-adapter — see GetDeclaration.
func (n *GlobalNamespace_Name) getDeclarationRef() *GlobalNamespace_Ref {
	return n.outer.refByID(n.declarationID)
}

// Java: GlobalNamespace.Name.getInitialization()
func (n *GlobalNamespace_Name) GetInitialization() *GlobalNamespace_Ref {
	return n.outer.refByID(n.initializationID)
}

// Java: GlobalNamespace.Name.isFunction()
func (n *GlobalNamespace_Name) IsFunction() bool {
	return n.GetBooleanProperty(GlobalNamespace_NameProp_FUNCTION)
}

// Java: GlobalNamespace.Name.isClass()
func (n *GlobalNamespace_Name) IsClass() bool {
	return n.GetBooleanProperty(GlobalNamespace_NameProp_CLASS)
}

// Java: GlobalNamespace.Name.isObjectLiteral()
func (n *GlobalNamespace_Name) IsObjectLiteral() bool {
	return n.GetBooleanProperty(GlobalNamespace_NameProp_OBJECTLIT)
}

// Java: GlobalNamespace.Name.getAliasingGets()
func (n *GlobalNamespace_Name) GetAliasingGets() int {
	return int(n.aliasingGets)
}

// Java: GlobalNamespace.Name.getSubclassingGets()
func (n *GlobalNamespace_Name) GetSubclassingGets() int {
	return int(n.subclassingGets)
}

// Java: GlobalNamespace.Name.getLocalSets()
func (n *GlobalNamespace_Name) GetLocalSets() int {
	return int(n.localSets)
}

// Java: GlobalNamespace.Name.getGlobalSets()
func (n *GlobalNamespace_Name) GetGlobalSets() int {
	return int(n.globalSets)
}

// Java: GlobalNamespace.Name.getTotalSets()
func (n *GlobalNamespace_Name) GetTotalSets() int {
	return int(n.globalSets + n.localSets)
}

// Java: GlobalNamespace.Name.getCallGets()
func (n *GlobalNamespace_Name) GetCallGets() int {
	return int(n.callGets)
}

// Java: GlobalNamespace.Name.getTotalGets()
func (n *GlobalNamespace_Name) GetTotalGets() int {
	return int(n.totalGets)
}

// Java: GlobalNamespace.Name.getDeleteProps()
func (n *GlobalNamespace_Name) GetDeleteProps() int {
	return int(n.deleteProps)
}

// Java: GlobalNamespace.Name.getParent()
func (n *GlobalNamespace_Name) GetParent() *GlobalNamespace_Name {
	return n.outer.nameByID(n.parentID)
}

// GetDeepestCommonAncestorChunk returns the JSChunk that is an ancestor (or
// equal to) all the chunks in which this name is referenced, inclusive of
// the actual chunks in the refs.
//
// Most of the time this is identical to the declaration's chunk, but it's
// still valid to conditionally reference names from later child chunks in
// an earlier chunk.
//
// Java: GlobalNamespace.Name.getDeepestCommonAncestorChunk(JSChunkGraph)
func (n *GlobalNamespace_Name) GetDeepestCommonAncestorChunk(chunkGraph *JSChunkGraph) *JSChunk {
	commonAncestor := n.getDeclarationRef().GetChunk()
	for r := range n.Refs() {
		rChunk := r.GetChunk()
		commonAncestor = chunkGraph.GetDeepestCommonDependencyInclusive(rChunk, commonAncestor)
	}
	return commonAncestor
}

// Java: GlobalNamespace.Name.getScope()
func (n *GlobalNamespace_Name) GetScope() rhino.StaticScope {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: GlobalNamespace.Name.setBooleanProperty(NameProp)
func (n *GlobalNamespace_Name) SetBooleanProperty(property GlobalNamespace_NameProp) {
	n.propertyBitSet = n.propertyBitSet | property.Bit()
}

// Java: GlobalNamespace.Name.getBooleanProperty(NameProp)
func (n *GlobalNamespace_Name) GetBooleanProperty(property GlobalNamespace_NameProp) bool {
	return (n.propertyBitSet & property.Bit()) != 0
}

func (n *GlobalNamespace_Name) multiRefsSlot() *globalNamespaceMultiRefs {
	return &n.outer.multiRefs[n.ref0&^gnMultiRefsBit]
}

func (n *GlobalNamespace_Name) compactMultiRefs(mr *globalNamespaceMultiRefs) {
	if mr.tombstones == 0 {
		return
	}
	g := n.outer
	w := 0
	for _, refID := range mr.refIDs {
		if refID != 0 {
			mr.refIDs[w] = refID
			if mr.refsIndex != nil {
				r := g.refByID(refID)
				if !r.node.IsNil() {
					mr.refsIndex[r.node] = w
				}
			}
			w++
		}
	}
	clear(mr.refIDs[w:])
	mr.refIDs = mr.refIDs[:w]
	mr.tombstones = 0
}

func (n *GlobalNamespace_Name) ensureRefsIndex(mr *globalNamespaceMultiRefs) map[rhino.Node]int {
	if mr.refsIndex != nil {
		return mr.refsIndex
	}
	if len(mr.refIDs)-int(mr.tombstones) > 24 {
		n.compactMultiRefs(mr)
		g := n.outer
		idxMap := make(map[rhino.Node]int, len(mr.refIDs)*2)
		for i, refID := range mr.refIDs {
			if refID != 0 {
				r := g.refByID(refID)
				if !r.node.IsNil() {
					idxMap[r.node] = i
				}
			}
		}
		mr.refsIndex = idxMap
		return idxMap
	}
	return nil
}

// Java: GlobalNamespace.Name.addRef(JSChunk, Scope, Node, Ref.Type)
func (n *GlobalNamespace_Name) AddRef(chunk *JSChunk, scope *Scope, node rhino.Node, type_ GlobalNamespace_Ref_Type) {
	n.CheckNoExistingRefsForNode(node)
	refID, ref := n.createNewRefWithID(chunk, scope, node, type_)
	n.putNewRefID(node, refID)
	n.updateStateForAddedRef(refID, ref)
}

// Java: GlobalNamespace.Name.checkNoExistingRefsForNode(Node)
func (n *GlobalNamespace_Name) CheckNoExistingRefsForNode(node rhino.Node) {
	if n.ref0 == 0 {
		return
	}
	g := n.outer
	if n.ref0&gnMultiRefsBit == 0 {
		r0 := g.refByID(n.ref0)
		jutil.CheckStateMsg1(r0.node != node, "Ref already exists for node: %s", r0)
		if n.ref1 != 0 {
			r1 := g.refByID(n.ref1)
			jutil.CheckStateMsg1(r1.node != node, "Ref already exists for node: %s", r1)
		}
		return
	}
	mr := n.multiRefsSlot()
	if idxMap := n.ensureRefsIndex(mr); idxMap != nil {
		if idx, ok := idxMap[node]; ok {
			jutil.CheckStateMsg1(false, "Ref already exists for node: %s", g.refByID(mr.refIDs[idx]))
		}
		return
	}
	for _, refID := range mr.refIDs {
		if refID != 0 {
			r := g.refByID(refID)
			jutil.CheckStateMsg1(r.node != node, "Ref already exists for node: %s", r)
		}
	}
}

func (n *GlobalNamespace_Name) createNewRefWithID(chunk *JSChunk, scope *Scope, node rhino.Node, type_ GlobalNamespace_Ref_Type) (uint32, *GlobalNamespace_Ref) {
	if !n.outer.forProcessDefines {
		jutil.CheckNotNull(scope)
	}
	rhino.CheckNotNull(node) // may be null later, but not on creation
	return n.outer.allocRef(chunk, scope, node, type_)
}

// Java: GlobalNamespace.Name.createNewRef(JSChunk, Scope, Node, Ref.Type)
func (n *GlobalNamespace_Name) CreateNewRef(chunk *JSChunk, scope *Scope, node rhino.Node, type_ GlobalNamespace_Ref_Type) *GlobalNamespace_Ref {
	_, ref := n.createNewRefWithID(chunk, scope, node, type_)
	return ref
}

func (n *GlobalNamespace_Name) putNewRefID(node rhino.Node, refID uint32) {
	if n.ref0 == 0 {
		n.ref0 = refID
		return
	}
	g := n.outer
	if n.ref0&gnMultiRefsBit == 0 {
		n.propertyBitSet |= gnNamePropRefsMapState
		if n.ref1 == 0 {
			n.ref1 = refID
			return
		}
		ids := g.allocRefSlice3Of4(n.ref0, n.ref1, refID)
		mrIdx := uint32(len(g.multiRefs))
		g.multiRefs = append(g.multiRefs, globalNamespaceMultiRefs{refIDs: ids})
		n.ref0 = gnMultiRefsBit | mrIdx
		n.ref1 = 0
		return
	}
	mr := n.multiRefsSlot()
	oldLen := len(mr.refIDs)
	if oldLen == cap(mr.refIDs) && int(mr.tombstones)*2 >= oldLen {
		n.compactMultiRefs(mr)
		oldLen = len(mr.refIDs)
	}
	if oldLen < cap(mr.refIDs) {
		mr.refIDs = mr.refIDs[:oldLen+1]
		mr.refIDs[oldLen] = refID
	} else {
		newCap := max(cap(mr.refIDs)*2, 8)
		next := make([]uint32, oldLen+1, newCap)
		copy(next, mr.refIDs)
		next[oldLen] = refID
		mr.refIDs = next
	}
	if mr.refsIndex != nil && !node.IsNil() {
		mr.refsIndex[node] = oldLen
	}
}

func (n *GlobalNamespace_Name) putRefID(node rhino.Node, refID uint32) {
	if n.ref0 == 0 {
		n.ref0 = refID
		return
	}
	g := n.outer
	if n.ref0&gnMultiRefsBit == 0 {
		if g.refByID(n.ref0).node == node {
			n.ref0 = refID
			return
		}
		if n.ref1 != 0 && g.refByID(n.ref1).node == node {
			n.ref1 = refID
			return
		}
		n.putNewRefID(node, refID)
		return
	}
	mr := n.multiRefsSlot()
	if idxMap := n.ensureRefsIndex(mr); idxMap != nil {
		if idx, ok := idxMap[node]; ok {
			mr.refIDs[idx] = refID
			return
		}
	} else {
		for i, id := range mr.refIDs {
			if id != 0 && g.refByID(id).node == node {
				mr.refIDs[i] = refID
				return
			}
		}
	}
	n.putNewRefID(node, refID)
}

// Java: GlobalNamespace.Name.putRef(Node, Ref)
func (n *GlobalNamespace_Name) PutRef(node rhino.Node, ref *GlobalNamespace_Ref) {
	n.putRefID(node, n.outer.findOrCopyRefID(ref))
}

// Java: GlobalNamespace.Name.addSingleRefForTesting(Node, Ref.Type)
func (n *GlobalNamespace_Name) AddSingleRefForTesting(node rhino.Node, type_ GlobalNamespace_Ref_Type) *GlobalNamespace_Ref {
	refID, ref := n.outer.allocRef( /* chunk= */ nil /* scope= */, nil /* node= */, node, type_)
	n.putRefID(node, refID)
	n.updateStateForAddedRef(refID, ref)
	return ref
}

// AddAliasingGetClonedFromDeclaration adds an ALIASING_GET Ref for the given
// Node using the same Ref properties as the declaration Ref, which must
// exist.
//
// Only for use by CollapseProperties.
//
// newRefNode is a newly added AST node that refers to this Name and appears
// in the same module and scope as the Ref that declares this Name.
//
// Java: GlobalNamespace.Name.addAliasingGetClonedFromDeclaration(Node)
func (n *GlobalNamespace_Name) AddAliasingGetClonedFromDeclaration(newRefNode rhino.Node) {
	// TODO(bradfordcsmith): It would be good to add checks that the scope is correct.
	declRef := jutil.CheckNotNull(n.getDeclarationRef())
	n.AddRef(declRef.chunk, declRef.scope, newRefNode, GlobalNamespace_Ref_Type_ALIASING_GET)
}

// UpdateStateForAddedRef updates counters and JSDocInfo recorded for the
// name to include a newly added Ref.
//
// Must be called exactly once when a new Ref is added.
//
// ref is a Ref that has just been added for this Name.
//
// Java: GlobalNamespace.Name.updateStateForAddedRef(Ref)
func (n *GlobalNamespace_Name) UpdateStateForAddedRef(ref *GlobalNamespace_Ref) {
	n.updateStateForAddedRef(n.outer.findOrCopyRefID(ref), ref)
}

func (n *GlobalNamespace_Name) updateStateForAddedRef(refID uint32, ref *GlobalNamespace_Ref) {
	switch ref.type_ {
	case GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL, GlobalNamespace_Ref_Type_SET_FROM_GLOBAL:
		if n.declarationID == 0 {
			n.declarationID = refID
		}
		if n.initializationID == 0 && !ref.IsUninitializedDeclaration() {
			// Record the reference where the first value is actually assigned.
			// Do not include the automatically-assigned `undefined` value case.
			// e.g. `var name;`
			n.initializationID = refID
		}
		if n.firstDeclDocID == 0 {
			// JSDocInfo from the first SET_FROM_GLOBAL will be assumed to be canonical
			// Note that this will not change if the first declaration is later removed
			// by optimizations.
			n.firstDeclDocID = n.outer.storeDoc(n.GetDocInfoForDeclaration(ref))
		}
		n.globalSets++
		if ref.type_ == GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL {
			n.aliasingGets++
			n.totalGets++
		}
	case GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL, GlobalNamespace_Ref_Type_SET_FROM_LOCAL:
		n.localSets++
		var info *rhino.JSDocInfo
		if !ref.GetNode().IsNil() {
			info = NodeUtil_GetBestJSDocInfo(ref.GetNode())
		}
		if info != nil && info.IsNoCollapse() {
			n.localSetsWithNoCollapse++
		}
		if ref.type_ == GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL {
			n.aliasingGets++
			n.totalGets++
		}
	case GlobalNamespace_Ref_Type_PROTOTYPE_GET, GlobalNamespace_Ref_Type_DIRECT_GET:
		node := ref.GetNode()
		if n.firstQnameDocID == 0 &&
			GlobalNamespace_IsQnameDeclarationWithoutAssignment(node) {
			// /** @type {sometype} */
			// some.qname.ref;
			n.firstQnameDocID = n.outer.storeDoc(node.GetJSDocInfo())
		}
		n.totalGets++
	case GlobalNamespace_Ref_Type_ALIASING_GET:
		n.aliasingGets++
		n.totalGets++
	case GlobalNamespace_Ref_Type_CALL_GET:
		n.callGets++
		n.totalGets++
	case GlobalNamespace_Ref_Type_DELETE_PROP:
		n.deleteProps++
	case GlobalNamespace_Ref_Type_SUBCLASSING_GET:
		n.subclassingGets++
		n.totalGets++
	}
}

// UpdateRefNode is the only safe way to update the Node belonging to a Ref
// once it is added to a Name.
//
// This is a specialized method that exists only for use by
// CollapseProperties.
//
// ref is the reference to update - it must belong to this name; newNode the
// new value for the ref's node.
//
// Java: GlobalNamespace.Name.updateRefNode(Ref, Node)
func (n *GlobalNamespace_Name) UpdateRefNode(ref *GlobalNamespace_Ref, newNode rhino.Node) {
	jutil.CheckArgumentMsg1(ref.node != newNode, "redundant update to Ref node: %s", ref)

	// Once a Ref's node is set to null, it shouldn't ever be set to anything else.
	// TODO(bradfordcsmith): Document here what it means when we set the node to null.
	//     Seems to be a way to keep name.getDeclaration() returning the original declaration
	//     Ref even though its node is no longer in the AST.
	oldNode := ref.GetNode()
	jutil.CheckStateMsg1(!oldNode.IsNil(), "Ref's node is already null: %s", ref)
	g := n.outer
	if n.ref0&gnMultiRefsBit != 0 {
		n.ensureRefsIndex(n.multiRefsSlot())
	}
	ref.node = newNode

	if n.ref0 == 0 {
		if n.propertyBitSet&gnNamePropRefsMapState == 0 {
			n.ref0 = g.findOrCopyRefID(ref)
		} else if !newNode.IsNil() {
			n.putNewRefID(newNode, g.findOrCopyRefID(ref))
		}
		return
	}
	if n.ref0&gnMultiRefsBit == 0 {
		if n.propertyBitSet&gnNamePropRefsMapState == 0 {
			// No update needed, since refsForNode is a singleton.
			jutil.CheckState(g.refByID(n.ref0) == ref)
			return
		}
		slot := -1
		var refID uint32
		if r0 := g.refByID(n.ref0); r0.node == oldNode || r0 == ref {
			slot = 0
			refID = n.ref0
		} else if n.ref1 != 0 {
			if r1 := g.refByID(n.ref1); r1.node == oldNode || r1 == ref {
				slot = 1
				refID = n.ref1
			}
		}
		if !newNode.IsNil() {
			if slot != 0 {
				if r0 := g.refByID(n.ref0); r0.node == newNode {
					jutil.CheckArgumentMsg1(false, "refs already exist: %s", r0)
				}
			}
			if slot != 1 && n.ref1 != 0 {
				if r1 := g.refByID(n.ref1); r1.node == newNode {
					jutil.CheckArgumentMsg1(false, "refs already exist: %s", r1)
				}
			}
		}
		if slot == 0 {
			n.ref0 = n.ref1
			n.ref1 = 0
		} else if slot == 1 {
			n.ref1 = 0
		}
		if !newNode.IsNil() {
			if refID == 0 {
				refID = g.findOrCopyRefID(ref)
			}
			n.putNewRefID(newNode, refID)
		}
		return
	}

	mr := n.multiRefsSlot()
	idx := -1
	var refID uint32
	if mr.refsIndex != nil {
		if i, ok := mr.refsIndex[oldNode]; ok {
			idx = i
			refID = mr.refIDs[i]
		}
	} else {
		for i, id := range mr.refIDs {
			if id != 0 {
				r := g.refByID(id)
				if r.node == oldNode || r == ref {
					idx = i
					refID = id
					break
				}
			}
		}
	}
	if !newNode.IsNil() {
		if mr.refsIndex != nil {
			if existingIdx, ok := mr.refsIndex[newNode]; ok && existingIdx != idx {
				jutil.CheckArgumentMsg1(false, "refs already exist: %s", g.refByID(mr.refIDs[existingIdx]))
			}
		} else {
			for i, id := range mr.refIDs {
				if i != idx && id != 0 {
					r := g.refByID(id)
					if r.node == newNode {
						jutil.CheckArgumentMsg1(false, "refs already exist: %s", r)
					}
				}
			}
		}
	}
	if idx >= 0 {
		mr.refIDs[idx] = 0
		mr.tombstones++
		if mr.refsIndex != nil {
			delete(mr.refsIndex, oldNode)
		}
		if !newNode.IsNil() {
			n.putNewRefID(newNode, refID)
		} else if len(mr.refIDs) == int(mr.tombstones) {
			mr.refIDs = mr.refIDs[:0]
			mr.tombstones = 0
			n.ref0 = 0
			n.ref1 = 0
		}
	} else if !newNode.IsNil() {
		if refID == 0 {
			refID = g.findOrCopyRefID(ref)
		}
		n.putNewRefID(newNode, refID)
	}
}

// RemoveRef removes the given Ref, which must belong to this Name.
//
// NOTE: if this is a twin ref, i.e. both a get and a set of this Name, this
// removes both the get and the set.
//
// Java: GlobalNamespace.Name.removeRef(Ref)
func (n *GlobalNamespace_Name) RemoveRef(ref *GlobalNamespace_Ref) {
	if n.ref0 == 0 {
		if n.propertyBitSet&gnNamePropRefsMapState == 0 {
			panic(jutil.NewNullPointerExceptionMsg("removeRef(" + jutil.ToJavaString(ref) + "): unknown ref"))
		}
		jutil.CheckStateMsg1(false, "removeRef(%s): unknown ref", ref)
	}
	g := n.outer
	if n.ref0&gnMultiRefsBit == 0 {
		if n.propertyBitSet&gnNamePropRefsMapState == 0 {
			jutil.CheckStateMsg1(g.refByID(n.ref0) == ref, "removeRef(%s): unknown ref", ref)
			n.ref0 = 0
		} else {
			refNode := ref.GetNode()
			found := g.refByID(n.ref0).node == refNode || (n.ref1 != 0 && g.refByID(n.ref1).node == refNode)
			jutil.CheckStateMsg1(found, "removeRef(%s): unknown ref", ref)
			if !refNode.IsNil() {
				n.RemoveRefFromNodeMap(ref)
			}
		}
	} else {
		mr := n.multiRefsSlot()
		refNode := ref.GetNode()
		found := false
		if idxMap := n.ensureRefsIndex(mr); idxMap != nil {
			_, found = idxMap[refNode]
		} else {
			for _, id := range mr.refIDs {
				if id != 0 && g.refByID(id).node == refNode {
					found = true
					break
				}
			}
		}
		jutil.CheckStateMsg1(found, "removeRef(%s): unknown ref", ref)
		if !refNode.IsNil() {
			n.RemoveRefFromNodeMap(ref)
		}
	}
	n.RemoveRefAndUpdateState(ref)
}

// RemoveRefAndUpdateState updates counts, declaration, and JSDoc to reflect
// removal of the given Ref.
//
// Java: GlobalNamespace.Name.removeRefAndUpdateState(Ref)
func (n *GlobalNamespace_Name) RemoveRefAndUpdateState(ref *GlobalNamespace_Ref) {
	g := n.outer
	if n.declarationID != 0 && g.refByID(n.declarationID) == ref {
		n.declarationID = 0
		if n.ref0 != 0 {
			if n.ref0&gnMultiRefsBit == 0 {
				if g.refByID(n.ref0).type_ == GlobalNamespace_Ref_Type_SET_FROM_GLOBAL {
					n.declarationID = n.ref0
				} else if n.ref1 != 0 && g.refByID(n.ref1).type_ == GlobalNamespace_Ref_Type_SET_FROM_GLOBAL {
					n.declarationID = n.ref1
				}
			} else {
				mr := n.multiRefsSlot()
				for _, id := range mr.refIDs {
					if id != 0 && g.refByID(id).type_ == GlobalNamespace_Ref_Type_SET_FROM_GLOBAL {
						n.declarationID = id
						break
					}
				}
			}
		}
	}

	var info *rhino.JSDocInfo
	switch ref.type_ {
	case GlobalNamespace_Ref_Type_SET_FROM_GLOBAL:
		n.globalSets--
	case GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL:
		n.aliasingGets--
		n.totalGets--
		n.globalSets--
	case GlobalNamespace_Ref_Type_SET_FROM_LOCAL, GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL:
		n.localSets--
		info = nil
		if !ref.GetNode().IsNil() {
			info = NodeUtil_GetBestJSDocInfo(ref.GetNode())
		}
		if info != nil && info.IsNoCollapse() {
			n.localSetsWithNoCollapse--
		}
		if ref.type_ == GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL {
			n.aliasingGets--
			n.totalGets--
		}
	case GlobalNamespace_Ref_Type_PROTOTYPE_GET, GlobalNamespace_Ref_Type_DIRECT_GET:
		n.totalGets--
	case GlobalNamespace_Ref_Type_ALIASING_GET:
		n.aliasingGets--
		n.totalGets--
	case GlobalNamespace_Ref_Type_CALL_GET:
		n.callGets--
		n.totalGets--
	case GlobalNamespace_Ref_Type_DELETE_PROP:
		n.deleteProps--
	case GlobalNamespace_Ref_Type_SUBCLASSING_GET:
		n.subclassingGets--
		n.totalGets--
		// Leaving off default: allows compile-time enforcement that all values are covered
	}
}

// Java: GlobalNamespace.Name.removeRefFromNodeMap(Ref)
func (n *GlobalNamespace_Name) RemoveRefFromNodeMap(ref *GlobalNamespace_Ref) {
	refNode := rhino.CheckNotNullMsg1(ref.GetNode(), "%s", ref)
	if n.ref0 == 0 {
		if n.propertyBitSet&gnNamePropRefsMapState == 0 {
			panic(jutil.NewNullPointerExceptionMsg("Missing ref when trying to remove it: " + jutil.ToJavaString(ref)))
		}
		jutil.CheckStateMsg2(false, "Unexpected Refs for Node: %s: when removing Ref: %s", (*GlobalNamespace_Ref)(nil), ref)
	}
	g := n.outer
	if n.ref0&gnMultiRefsBit == 0 {
		if n.propertyBitSet&gnNamePropRefsMapState == 0 {
			n.ref0 = 0
			return
		}
		if r0 := g.refByID(n.ref0); r0.node == refNode {
			jutil.CheckStateMsg2(r0 == ref, "Unexpected Refs for Node: %s: when removing Ref: %s", r0, ref)
			n.ref0 = n.ref1
			n.ref1 = 0
			return
		}
		if n.ref1 != 0 {
			if r1 := g.refByID(n.ref1); r1.node == refNode {
				jutil.CheckStateMsg2(r1 == ref, "Unexpected Refs for Node: %s: when removing Ref: %s", r1, ref)
				n.ref1 = 0
				return
			}
		}
		jutil.CheckStateMsg2(false, "Unexpected Refs for Node: %s: when removing Ref: %s", (*GlobalNamespace_Ref)(nil), ref)
		return
	}
	mr := n.multiRefsSlot()
	idx := -1
	var refsForNode *GlobalNamespace_Ref
	if idxMap := n.ensureRefsIndex(mr); idxMap != nil {
		if i, ok := idxMap[refNode]; ok {
			idx = i
			refsForNode = g.refByID(mr.refIDs[i])
		}
	} else {
		for i, id := range mr.refIDs {
			if id != 0 {
				r := g.refByID(id)
				if r.node == refNode {
					refsForNode = r
					idx = i
					break
				}
			}
		}
	}
	jutil.CheckStateMsg2(
		refsForNode == ref,
		"Unexpected Refs for Node: %s: when removing Ref: %s",
		refsForNode,
		ref)
	mr.refIDs[idx] = 0
	mr.tombstones++
	if mr.refsIndex != nil {
		delete(mr.refsIndex, refNode)
	}
	if len(mr.refIDs) == int(mr.tombstones) {
		mr.refIDs = mr.refIDs[:0]
		mr.tombstones = 0
		n.ref0 = 0
		n.ref1 = 0
	}
}

// Refs iterates the Refs of this name in order with zero heap allocations.
func (n *GlobalNamespace_Name) Refs() iter.Seq[*GlobalNamespace_Ref] {
	return func(yield func(*GlobalNamespace_Ref) bool) {
		if n.ref0 == 0 {
			return
		}
		g := n.outer
		if n.ref0&gnMultiRefsBit == 0 {
			if !yield(g.refByID(n.ref0)) {
				return
			}
			if n.ref1 != 0 {
				yield(g.refByID(n.ref1))
			}
			return
		}
		mr := n.multiRefsSlot()
		if mr.tombstones > 0 {
			n.compactMultiRefs(mr)
		}
		for _, id := range mr.refIDs {
			if id != 0 && !yield(g.refByID(id)) {
				return
			}
		}
	}
}

// snapshotRefs returns a stable snapshot slice of this Name's Refs, using buf
// when len <= len(buf) so callers that mutate refs during iteration avoid heap
// allocation in the common case.
func (n *GlobalNamespace_Name) snapshotRefs(buf *[8]*GlobalNamespace_Ref) []*GlobalNamespace_Ref {
	if n.ref0 == 0 {
		return nil
	}
	g := n.outer
	if n.ref0&gnMultiRefsBit == 0 {
		buf[0] = g.refByID(n.ref0)
		if n.ref1 == 0 {
			return buf[:1]
		}
		buf[1] = g.refByID(n.ref1)
		return buf[:2]
	}
	mr := n.multiRefsSlot()
	if mr.tombstones > 0 {
		n.compactMultiRefs(mr)
	}
	cnt := len(mr.refIDs)
	var out []*GlobalNamespace_Ref
	if cnt <= len(buf) {
		out = buf[:cnt]
	} else {
		out = make([]*GlobalNamespace_Ref, cnt)
	}
	for i, id := range mr.refIDs {
		out[i] = g.refByID(id)
	}
	return out
}

// GetRefs returns the Refs of this name.
//
// Go-only: language-adapter — Java returns a Collection (an immutable set
// or the live values() view of the map); this returns a snapshot slice in
// the same order. No caller modifies the Refs while iterating them.
//
// Java: GlobalNamespace.Name.getRefs()
func (n *GlobalNamespace_Name) GetRefs() []*GlobalNamespace_Ref {
	if n.ref0 == 0 {
		return nil
	}
	g := n.outer
	if n.ref0&gnMultiRefsBit == 0 {
		if n.ref1 == 0 {
			return []*GlobalNamespace_Ref{g.refByID(n.ref0)}
		}
		return []*GlobalNamespace_Ref{g.refByID(n.ref0), g.refByID(n.ref1)}
	}
	mr := n.multiRefsSlot()
	if mr.tombstones > 0 {
		n.compactMultiRefs(mr)
	}
	out := make([]*GlobalNamespace_Ref, len(mr.refIDs))
	for i, id := range mr.refIDs {
		out[i] = g.refByID(id)
	}
	return out
}

// GetRefForNode gets the Ref for this name that belongs to the given node.
//
// Returns null if there are no Refs corresponding to the node.
//
// Java: GlobalNamespace.Name.getRefForNode(Node)
func (n *GlobalNamespace_Name) GetRefForNode(node rhino.Node) *GlobalNamespace_Ref {
	rhino.CheckNotNull(node)
	if n.ref0 == 0 {
		return nil
	}
	g := n.outer
	if n.ref0&gnMultiRefsBit == 0 {
		if r0 := g.refByID(n.ref0); r0.node == node {
			return r0
		}
		if n.ref1 != 0 {
			if r1 := g.refByID(n.ref1); r1.node == node {
				return r1
			}
		}
		return nil
	}
	mr := n.multiRefsSlot()
	if idxMap := n.ensureRefsIndex(mr); idxMap != nil {
		if idx, ok := idxMap[node]; ok {
			return g.refByID(mr.refIDs[idx])
		}
		return nil
	}
	for _, id := range mr.refIDs {
		if id != 0 {
			r := g.refByID(id)
			if r.node == node {
				return r
			}
		}
	}
	return nil
}

// Java: GlobalNamespace.Name.getFirstRef()
func (n *GlobalNamespace_Name) GetFirstRef() *GlobalNamespace_Ref {
	if n.ref0 == 0 {
		if n.propertyBitSet&gnNamePropRefsMapState == 0 {
			panic(jutil.NewNullPointerExceptionMsg("no first Ref to get"))
		}
		panic(jutil.NewNoSuchElementException())
	}
	g := n.outer
	if n.ref0&gnMultiRefsBit == 0 {
		return g.refByID(n.ref0)
	}
	mr := n.multiRefsSlot()
	if mr.tombstones > 0 {
		n.compactMultiRefs(mr)
	}
	if len(mr.refIDs) > 0 {
		return g.refByID(mr.refIDs[0])
	}
	panic(jutil.NewNoSuchElementException())
}

func (n *GlobalNamespace_Name) refCount() int {
	if n.ref0 == 0 {
		return 0
	}
	if n.ref0&gnMultiRefsBit == 0 {
		if n.ref1 == 0 {
			return 1
		}
		return 2
	}
	mr := n.multiRefsSlot()
	return len(mr.refIDs) - int(mr.tombstones)
}

// Java: GlobalNamespace.Name.canEliminate()
func (n *GlobalNamespace_Name) CanEliminate() bool {
	if !n.CanCollapseUnannotatedChildNames() || n.totalGets > 0 {
		return false
	}

	if n.HasProps() {
		for p := range n.Props() {
			if !p.CanCollapse() {
				return false
			}
		}
	}
	return true
}

// Java: GlobalNamespace.Name.isSimpleStubDeclaration()
func (n *GlobalNamespace_Name) IsSimpleStubDeclaration() bool {
	if n.refCount() == 1 {
		ref := n.GetFirstRef()
		if ref.node.GetParent().IsExprResult() {
			return true
		}
	}
	return false
}

// Java: GlobalNamespace.Name.isCollapsingExplicitlyDenied()
func (n *GlobalNamespace_Name) IsCollapsingExplicitlyDenied() bool {
	docInfo := n.GetJSDocInfo()
	return docInfo != nil && docInfo.IsNoCollapse()
}

// -------------------------------------------------------------------------

// GlobalNamespace_Ref is a global name reference. Contains references to the
// relevant parse tree node and its ancestors that may be affected.
type GlobalNamespace_Ref struct {
	chunk *JSChunk

	// The scope in which the reference is resolved. Note that for
	// ALIASING_GETS like "var x = ns;" this scope may not be the correct
	// hoist scope of the aliasing VAR.
	scope *Scope

	// Not final because CollapseProperties needs to update the namespace in-place.
	node  rhino.Node
	type_ GlobalNamespace_Ref_Type
}

// GlobalNamespace_Ref_Type is the enum GlobalNamespace.Ref.Type.
//
// Note: we are more aggressive about collapsing @enum and @constructor
// declarations than implied here, see Name#canCollapse. Non-private for
// testing.
type GlobalNamespace_Ref_Type int32

const (
	// Set in the scope in which a name is declared, either the global scope
	// or a module scope: `a.b.c = 0;` or `goog.module('mod'); exports.Foo =
	// class {};`
	GlobalNamespace_Ref_Type_SET_FROM_GLOBAL GlobalNamespace_Ref_Type = iota

	// Set in a local scope: function f() { a.b.c = 0; }
	GlobalNamespace_Ref_Type_SET_FROM_LOCAL

	// Combined get and set in the scope in which a name is declared, either
	// the global scope or a module scope: `const c = a.b.c = 0;` or
	// `goog.module('mod'); exports.Foo = class {};`
	GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL

	// Combined get and set in a local scope: function f() { return a.b.c = 0; }
	GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL

	// Get a name's prototype: a.b.c.prototype
	GlobalNamespace_Ref_Type_PROTOTYPE_GET

	// Includes all uses that prevent a name's properties from being
	// collapsed: var x = a.b.c f(a.b.c) new Foo(a.b.c)
	GlobalNamespace_Ref_Type_ALIASING_GET

	// Includes all uses that prevent a name from being completely
	// eliminated: goog.inherits(anotherName, a.b.c) new a.b.c() x instanceof
	// a.b.c void a.b.c if (a.b.c) {}
	GlobalNamespace_Ref_Type_DIRECT_GET

	// Calling a name: a.b.c(); Prevents a name from being collapsed if never
	// set.
	GlobalNamespace_Ref_Type_CALL_GET

	// Deletion of a property: delete a.b.c; Prevents a name from being
	// collapsed at all.
	GlobalNamespace_Ref_Type_DELETE_PROP

	// ES6 subclassing ref: class extends A {}
	GlobalNamespace_Ref_Type_SUBCLASSING_GET
)

var globalNamespace_Ref_Type_names = [...]string{
	"SET_FROM_GLOBAL",
	"SET_FROM_LOCAL",
	"GET_AND_SET_FROM_GLOBAL",
	"GET_AND_SET_FROM_LOCAL",
	"PROTOTYPE_GET",
	"ALIASING_GET",
	"DIRECT_GET",
	"CALL_GET",
	"DELETE_PROP",
	"SUBCLASSING_GET",
}

// Java: GlobalNamespace.Ref.Type.name()
func (t GlobalNamespace_Ref_Type) String() string { return globalNamespace_Ref_Type_names[t] }

// Java: GlobalNamespace.Ref.Type.ordinal()
func (t GlobalNamespace_Ref_Type) Ordinal() int { return int(t) }

// NewGlobalNamespace_Ref creates a Ref
//
// No parameter checking is done here, because we allow nulls for several
// fields in Refs created just for testing. However, all Refs for real use
// must be created by methods on the Name class, which does do argument
// checking.
//
// Java: GlobalNamespace.Ref(JSChunk, Scope, Node, Type)
func NewGlobalNamespace_Ref(chunk *JSChunk, scope *Scope, node rhino.Node, type_ GlobalNamespace_Ref_Type) *GlobalNamespace_Ref {
	return &GlobalNamespace_Ref{chunk: chunk, node: node, type_: type_, scope: scope}
}

// Java: GlobalNamespace.Ref.getChunk()
func (r *GlobalNamespace_Ref) GetChunk() *JSChunk {
	return r.chunk
}

// Java: GlobalNamespace.Ref.getNode()
func (r *GlobalNamespace_Ref) GetNode() rhino.Node {
	return r.node
}

// Java: GlobalNamespace.Ref.getSourceFile()
func (r *GlobalNamespace_Ref) GetSourceFile() rhino.StaticSourceFile {
	if !r.node.IsNil() {
		return r.node.GetStaticSourceFile()
	}
	return nil
}

// Java: GlobalNamespace.Ref.getSymbol()
func (r *GlobalNamespace_Ref) GetSymbol() rhino.StaticSlot {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: GlobalNamespace.Ref.isDeleteProp()
func (r *GlobalNamespace_Ref) IsDeleteProp() bool {
	return r.type_ == GlobalNamespace_Ref_Type_DELETE_PROP
}

// Java: GlobalNamespace.Ref.isSubclassingGet()
func (r *GlobalNamespace_Ref) IsSubclassingGet() bool {
	return r.type_ == GlobalNamespace_Ref_Type_SUBCLASSING_GET
}

// IsTwin reports whether this is a "twin" ref, i.e. a ref that is both a
// get and a set.
//
// Example: `a.b` from `x = a.b = 0;`
//
// Java: GlobalNamespace.Ref.isTwin()
func (r *GlobalNamespace_Ref) IsTwin() bool {
	switch r.type_ {
	case GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL, GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL:
		return true
	default:
		return false
	}
}

// Java: GlobalNamespace.Ref.isGet()
func (r *GlobalNamespace_Ref) IsGet() bool {
	switch r.type_ {
	case GlobalNamespace_Ref_Type_DIRECT_GET,
		GlobalNamespace_Ref_Type_ALIASING_GET,
		GlobalNamespace_Ref_Type_SUBCLASSING_GET,
		GlobalNamespace_Ref_Type_CALL_GET,
		GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL,
		GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL,
		GlobalNamespace_Ref_Type_PROTOTYPE_GET:
		return true
	default:
		return false
	}
}

// Java: GlobalNamespace.Ref.isAliasingGet()
func (r *GlobalNamespace_Ref) IsAliasingGet() bool {
	switch r.type_ {
	case GlobalNamespace_Ref_Type_ALIASING_GET,
		GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL,
		GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL:
		return true
	default:
		return false
	}
}

// Java: GlobalNamespace.Ref.isSet()
func (r *GlobalNamespace_Ref) IsSet() bool {
	switch r.type_ {
	case GlobalNamespace_Ref_Type_SET_FROM_GLOBAL,
		GlobalNamespace_Ref_Type_SET_FROM_LOCAL,
		GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL,
		GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL:
		return true
	default:
		return false
	}
}

// Java: GlobalNamespace.Ref.isSetFromGlobal()
func (r *GlobalNamespace_Ref) IsSetFromGlobal() bool {
	return r.type_ == GlobalNamespace_Ref_Type_SET_FROM_GLOBAL || r.type_ == GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL
}

// IsUninitializedDeclaration: true for `var name;` and similar cases.
//
// Java: GlobalNamespace.Ref.isUninitializedDeclaration()
func (r *GlobalNamespace_Ref) IsUninitializedDeclaration() bool {
	return r.node.IsName() && NodeUtil_IsNameDeclaration(r.node.GetParent()) && !r.node.HasChildren()
}

// Java: GlobalNamespace.Ref.toString()
func (r *GlobalNamespace_Ref) String() string {
	// MoreObjects.toStringHelper(this).omitNullValues()
	var sb strings.Builder
	sb.WriteString("Ref{type=")
	sb.WriteString(r.type_.String())
	if !r.node.IsNil() {
		sb.WriteString(", node=")
		sb.WriteString(r.node.String())
	}
	if r.scope != nil {
		sb.WriteString(", scope=")
		sb.WriteString(r.scope.String())
	}
	sb.WriteByte('}')
	return sb.String()
}
