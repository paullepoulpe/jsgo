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

func (n *GlobalNamespace_Name) markEscaped() {
	n.propertyBitSet |= gnNameEscaped
}

func (n *GlobalNamespace_Name) isEscaped() bool {
	return n.propertyBitSet&gnNameEscaped != 0
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
