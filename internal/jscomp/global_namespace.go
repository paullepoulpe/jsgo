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

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

const (
	// Keep page sizes <= 8 KB so Go allocates single-page spans instead of
	// multi-page heap spans on 35M-node inputs (see jutil.Slab).
	gnNamePageShift = 6
	gnNamePageLen   = 1 << gnNamePageShift // 64 * 96B = 6144B
	gnNamePageMask  = gnNamePageLen - 1

	gnRefPageShift = 8
	gnRefPageLen   = 1 << gnRefPageShift // 256 * 24B = 6144B
	gnRefPageMask  = gnRefPageLen - 1

	gnMultiRefsBit uint32 = 1 << 31

	// Internal bits in Name.propertyBitSet (NameProp uses bits 0..13).
	gnNameEscaped           int32 = 1 << 28
	gnNamePropRefsMapState  int32 = 1 << 29
	gnNamePropIndexedGlobal int32 = 1 << 30
)

var (
	gnPrototypeStrID            = uint32(strpool.Intern("prototype"))
	gnThisStrID                 = uint32(strpool.Intern("this"))
	gnSuperStrID                = uint32(strpool.Intern("super"))
	gnClosureDefinesStrID       = uint32(strpool.Intern("CLOSURE_DEFINES"))
	gnHasOwnPropertyStrID       = uint32(strpool.Intern("hasOwnProperty"))
	gnPropertyIsEnumerableStrID = uint32(strpool.Intern("propertyIsEnumerable"))
)

type gnChildSlot struct {
	key    uint64 // (uint64(parentID) << 32) | uint64(baseNameID); 0 = empty
	nameID uint32 // 1-based Name ID
}

type globalNamespaceMultiRefs struct {
	refIDs     []uint32
	refsIndex  map[rhino.Node]int
	tombstones int32
}

// GlobalNamespace builds a namespace of all qualified names whose root is in
// the global scope or a module, plus an index of all references to those
// global names.
//
// This class tracks assignments to qualified names (e.g. `a.b.c`), because
// we often want to treat them as if they were global variables (e.g.
// CollapseProperties). However, when a qualified name begins an optional
// chain, we will not consider the optional chain to be part of the qualified
// name (e.g. `a.b?.c` is the qualified name `a.b` with an optional reference
// to property `c`.) We will record such optional chains as ALIASING_GET
// references to the non-optional qualified name part.
//
// When used as a StaticScope this class acts like a single parentless global
// scope. The module references are currently only accessible by
// getNameFromModule(ModuleMetadata, String), as many use cases only care
// about global names. (This may change as module rewriting is moved later in
// compilation). Module tracking also only occurs when
// com.google.javascript.jscomp.modules.ModuleMapCreator has run.
//
// The namespace can be updated as the AST is changed. Removing names or
// references should be done by the methods on Name. Adding new names should
// be done with scanNewNodes.
//
// Go-only: language-adapter — Java narrows StaticScope's getSlot/getOwnSlot
// and StaticSlot's getDeclaration to Name/Ref. They keep the rhino
// interfaces' Go signatures (never returning a typed nil); callers that need
// the Java type assert with the comma-ok form, as for scopes (see
// abstract_scope.go).
type GlobalNamespace struct {
	compiler                      AbstractCompiler
	enableImplicitlyAliasedValues bool
	root                          rhino.Node
	externsRoot                   rhino.Node
	globalRoot                    rhino.Node
	spreadSiblingCache            map[rhino.Node]bool
	sourceKind                    GlobalNamespace_SourceKind
	generated                     bool

	// Records decisions made by this class.
	//
	// Since this class is a utility used by others, the creating class may
	// provide the log file.
	decisionsLog diagnostic.LogFile

	// Global namespace tree (1-based Name IDs into namePages).
	globalNameIDs []uint32

	// Flat noscan open-addressed hash table mapping (parentID, baseNameID) -> 1-based Name ID.
	// parentID == 0 holds global roots; parentID != 0 holds child properties of both
	// global and module-scoped names (since parentID is already unique per module root).
	childTable         []gnChildSlot
	childCount         uint32
	childShift         uint8
	indexedGlobalCount uint32
	hasProvidedNames   bool

	// Maps (ModuleMetadata, root baseNameID) to 1-based root Name IDs in that module.
	nameMapByModule map[globalNamespaceModuleKey]uint32

	// Limits traversal to scripts matching the given predicate.
	shouldTraverseScript func(rhino.Node) bool

	// Go-only: performance — Names and Refs live in flat uint32-indexed pages
	// with minimal pointer footprints so GC mark work stays negligible.
	namePages []*[gnNamePageLen]GlobalNamespace_Name
	nameCount uint32
	refPages  []*[gnRefPageLen]GlobalNamespace_Ref
	refCount  uint32
	refIDSlab []uint32
	multiRefs []globalNamespaceMultiRefs
	jsDocs    []*rhino.JSDocInfo

	forProcessDefines      bool
	processDefinesCheck    bool
	hasNonScriptDefineCand bool
	defineCandidateScripts map[rhino.Node]struct{}

	// Go-only: performance — for ProcessDefines, the strpool IDs of leaf
	// segments that can belong to a define or (in check mode) a global constant
	// whose value a define may use; the namespace records no other names. nil
	// when the namespace serves any other pass. See collectDefineCandidates.
	defineCandidateLeafIDs map[uint32]struct{}
}

var globalNamespace_globalRoot = rhino.IR_Root()

// Java: GlobalNamespace.GOOG_PROVIDE
var GlobalNamespace_GOOG_PROVIDE = rhino.QualifiedName_Of("goog.provide")

// GlobalNamespace_SourceKind is the enum GlobalNamespace.SourceKind.
type GlobalNamespace_SourceKind int32

const (
	GlobalNamespace_SourceKind_EXTERN GlobalNamespace_SourceKind = iota
	GlobalNamespace_SourceKind_TYPE_SUMMARY
	GlobalNamespace_SourceKind_CODE
)

var globalNamespace_SourceKind_names = [...]string{"EXTERN", "TYPE_SUMMARY", "CODE"}

// Java: GlobalNamespace.SourceKind.name()
func (k GlobalNamespace_SourceKind) String() string { return globalNamespace_SourceKind_names[k] }

// Java: GlobalNamespace.SourceKind.ordinal()
func (k GlobalNamespace_SourceKind) Ordinal() int { return int(k) }

// Java: GlobalNamespace.SourceKind.fromScriptNode(Node)
func GlobalNamespace_SourceKind_FromScriptNode(n rhino.Node) GlobalNamespace_SourceKind {
	if !n.IsFromExterns() {
		return GlobalNamespace_SourceKind_CODE
	} else if NodeUtil_IsFromTypeSummary(n) {
		return GlobalNamespace_SourceKind_TYPE_SUMMARY
	} else {
		return GlobalNamespace_SourceKind_EXTERN
	}
}

// NewGlobalNamespaceWithLog creates an instance that may emit warnings when
// building the namespace.
//
// compiler is the AbstractCompiler, for reporting code changes; root the
// root of the rest of the code to build a namespace for.
//
// Java: GlobalNamespace(LogFile, AbstractCompiler, Node)
func NewGlobalNamespaceWithLog(decisionsLog diagnostic.LogFile, compiler AbstractCompiler, root rhino.Node) *GlobalNamespace {
	return NewGlobalNamespaceWithLogAndRoots(decisionsLog, compiler, rhino.NilNode, root)
}

// NewGlobalNamespace creates an instance that may emit warnings when
// building the namespace.
//
// compiler is the AbstractCompiler, for reporting code changes; root the
// root of the rest of the code to build a namespace for.
//
// Java: GlobalNamespace(AbstractCompiler, Node)
func NewGlobalNamespace(compiler AbstractCompiler, root rhino.Node) *GlobalNamespace {
	return NewGlobalNamespaceWithLogAndRoots( /* decisionsLog= */ diagnostic.LogFile{}, compiler, rhino.NilNode, root)
}

// NewGlobalNamespaceWithRoots creates an instance that may emit warnings
// when building the namespace.
//
// compiler is the AbstractCompiler, for reporting code changes; externsRoot
// the root of the externs to build a namespace for. If this is null, externs
// and properties defined on extern types will not be included in the global
// namespace. If non-null, it allows user-defined function on extern types to
// be included in the global namespace. E.g. String.foo. root is the root of
// the rest of the code to build a namespace for.
//
// Java: GlobalNamespace(AbstractCompiler, Node, Node)
func NewGlobalNamespaceWithRoots(compiler AbstractCompiler, externsRoot rhino.Node, root rhino.Node) *GlobalNamespace {
	return NewGlobalNamespaceWithLogAndRoots(diagnostic.LogFile{}, compiler, externsRoot, root)
}

// NewGlobalNamespaceWithLogAndRoots creates an instance that may emit
// warnings when building the namespace.
//
// decisionsLog is where to log decisions made by this instance; compiler the
// AbstractCompiler, for reporting code changes; externsRoot the root of the
// externs to build a namespace for. If this is null, externs and properties
// defined on extern types will not be included in the global namespace. If
// non-null, it allows user-defined function on extern types to be included
// in the global namespace. E.g. String.foo. root is the root of the rest of
// the code to build a namespace for.
//
// Java: GlobalNamespace(LogFile, AbstractCompiler, Node, Node)
func NewGlobalNamespaceWithLogAndRoots(
	decisionsLog diagnostic.LogFile,
	compiler AbstractCompiler,
	externsRoot rhino.Node,
	root rhino.Node) *GlobalNamespace {
	g := &GlobalNamespace{
		globalRoot: globalNamespace_globalRoot,
		generated:  false,
	}
	g.decisionsLog = decisionsLog
	g.compiler = compiler
	g.externsRoot = externsRoot
	g.root = root
	g.enableImplicitlyAliasedValues = !compiler.GetOptions().GetAssumeStaticInheritanceIsNotUsed()
	return g
}

func (g *GlobalNamespace) nameByID(id uint32) *GlobalNamespace_Name {
	if id == 0 {
		return nil
	}
	idx := id - 1
	return &g.namePages[idx>>gnNamePageShift][idx&gnNamePageMask]
}

func (g *GlobalNamespace) allocName(baseNameID uint32, parentID uint32, sourceKind GlobalNamespace_SourceKind) *GlobalNamespace_Name {
	idx := g.nameCount
	g.nameCount = idx + 1
	pageIdx := idx >> gnNamePageShift
	if int(pageIdx) >= len(g.namePages) {
		g.namePages = append(g.namePages, new([gnNamePageLen]GlobalNamespace_Name))
	}
	id := idx + 1
	n := &g.namePages[pageIdx][idx&gnNamePageMask]
	n.outer = g
	n.selfID = id
	n.baseNameID = baseNameID
	n.parentID = parentID
	n.SetBooleanProperty(GlobalNamespace_NameProp_NOT_A_TYPE)
	n.SetBooleanProperty(GlobalNamespace_NameProp_OTHER_OBJECT)
	if sourceKind == GlobalNamespace_SourceKind_EXTERN {
		n.SetBooleanProperty(GlobalNamespace_NameProp_IS_EXTERN)
	}
	return n
}

func (g *GlobalNamespace) initChildTable(capHint int) {
	minSlots := max((capHint*10)/7+1, 64)
	log2 := uint8(6)
	for (1 << log2) < minSlots {
		log2++
	}
	g.childTable = make([]gnChildSlot, 1<<log2)
	g.childShift = 64 - log2
}

func (g *GlobalNamespace) lookupChildID(parentID, baseNameID uint32) uint32 {
	tbl := g.childTable
	if len(tbl) == 0 || baseNameID == 0 {
		return 0
	}
	key := (uint64(parentID) << 32) | uint64(baseNameID)
	mask := uint32(len(tbl) - 1)
	slot := uint32((key * 0x9e3779b97f4a7c15) >> g.childShift)
	for {
		e := tbl[slot]
		if e.key == key {
			return e.nameID
		}
		if e.key == 0 {
			return 0
		}
		slot = (slot + 1) & mask
	}
}

func (g *GlobalNamespace) putChildID(parentID, baseNameID, nameID uint32) {
	if int(g.childCount+1)*10 > len(g.childTable)*7 {
		g.growChildTable()
	}
	key := (uint64(parentID) << 32) | uint64(baseNameID)
	tbl := g.childTable
	mask := uint32(len(tbl) - 1)
	slot := uint32((key * 0x9e3779b97f4a7c15) >> g.childShift)
	for {
		e := tbl[slot]
		if e.key == 0 {
			tbl[slot] = gnChildSlot{key: key, nameID: nameID}
			g.childCount++
			return
		}
		if e.key == key {
			tbl[slot].nameID = nameID
			return
		}
		slot = (slot + 1) & mask
	}
}

func (g *GlobalNamespace) growChildTable() {
	old := g.childTable
	var newLog2 uint8 = 6
	if len(old) > 0 {
		newLog2 = (64 - g.childShift) + 1
	}
	next := make([]gnChildSlot, 1<<newLog2)
	newShift := 64 - newLog2
	mask := uint32(len(next) - 1)
	for _, e := range old {
		if e.key == 0 {
			continue
		}
		slot := uint32((e.key * 0x9e3779b97f4a7c15) >> newShift)
		for next[slot].key != 0 {
			slot = (slot + 1) & mask
		}
		next[slot] = e
	}
	g.childTable = next
	g.childShift = newShift
}

func (g *GlobalNamespace) allocRefSlice3Of4(id0, id1, id2 uint32) []uint32 {
	if len(g.refIDSlab) < 4 {
		g.refIDSlab = make([]uint32, 1024)
	}
	s := g.refIDSlab[:3:4]
	g.refIDSlab = g.refIDSlab[4:]
	s[0] = id0
	s[1] = id1
	s[2] = id2
	return s
}

func (g *GlobalNamespace) refByID(id uint32) *GlobalNamespace_Ref {
	if id == 0 {
		return nil
	}
	idx := id - 1
	return &g.refPages[idx>>gnRefPageShift][idx&gnRefPageMask]
}

func (g *GlobalNamespace) allocRef(chunk *JSChunk, scope *Scope, node rhino.Node, type_ GlobalNamespace_Ref_Type) (uint32, *GlobalNamespace_Ref) {
	idx := g.refCount
	id := idx + 1
	g.refCount = id
	pageIdx := idx >> gnRefPageShift
	if int(pageIdx) >= len(g.refPages) {
		g.refPages = append(g.refPages, new([gnRefPageLen]GlobalNamespace_Ref))
	}
	r := &g.refPages[pageIdx][idx&gnRefPageMask]
	*r = GlobalNamespace_Ref{chunk: chunk, scope: scope, node: node, type_: type_, pageRefID: id}
	return id, r
}

func (g *GlobalNamespace) findOrCopyRefID(ref *GlobalNamespace_Ref) uint32 {
	if ref == nil {
		return 0
	}
	if id := ref.pageRefID; id != 0 && id <= g.refCount {
		i := id - 1
		if &g.refPages[i>>gnRefPageShift][i&gnRefPageMask] == ref {
			return id
		}
	}
	id, _ := g.allocRef(ref.chunk, ref.scope, ref.node, ref.type_)
	return id
}

func (g *GlobalNamespace) storeDoc(doc *rhino.JSDocInfo) uint32 {
	if doc == nil {
		return 0
	}
	g.jsDocs = append(g.jsDocs, doc)
	return uint32(len(g.jsDocs))
}

func (g *GlobalNamespace) docByID(id uint32) *rhino.JSDocInfo {
	if id == 0 {
		return nil
	}
	return g.jsDocs[id-1]
}

// LookupName returns the Name for fullName if present, without interning
// fullName when absent.
func (g *GlobalNamespace) LookupName(fullName string) *GlobalNamespace_Name {
	if g.indexedGlobalCount == 0 || fullName == "" {
		return nil
	}
	var curID uint32
	rem := fullName
	for {
		dot := strings.IndexByte(rem, '.')
		var seg string
		if dot < 0 {
			seg = rem
		} else {
			seg = rem[:dot]
		}
		segID := uint32(strpool.Lookup(seg))
		if segID == 0 {
			return nil
		}
		curID = g.lookupChildID(curID, segID)
		if curID == 0 {
			return nil
		}
		if dot < 0 {
			n := g.nameByID(curID)
			if n.propertyBitSet&gnNamePropIndexedGlobal != 0 {
				return n
			}
			return nil
		}
		rem = rem[dot+1:]
	}
}

// ContainsName reports whether fullName is indexed in the global namespace.
func (g *GlobalNamespace) ContainsName(fullName string) bool {
	return g.LookupName(fullName) != nil
}

// lookupGlobalNode resolves an AST qualified-name node directly against the
// interned child table without allocating an intermediate qualified-name
// string. ok is false iff qname is not a valid qualified name.
func (g *GlobalNamespace) lookupGlobalNode(qname rhino.Node) (*GlobalNamespace_Name, bool) {
	if qname.IsNil() {
		return nil, false
	}
	g.EnsureGenerated()
	switch qname.GetToken() {
	case rhino.Token_NAME:
		id := qname.GetStringID()
		if id <= rhino.RhinoStringPool_EmptyID {
			return nil, false
		}
		if g.indexedGlobalCount == 0 {
			return nil, true
		}
		if s := strpool.String(strpool.ID(id)); strings.IndexByte(s, '.') >= 0 {
			return g.LookupName(s), true
		}
		if nameID := g.lookupChildID(0, id); nameID != 0 {
			if n := g.nameByID(nameID); n.propertyBitSet&gnNamePropIndexedGlobal != 0 {
				return n, true
			}
		}
		return nil, true
	case rhino.Token_THIS:
		if g.indexedGlobalCount != 0 {
			if nameID := g.lookupChildID(0, gnThisStrID); nameID != 0 {
				if n := g.nameByID(nameID); n.propertyBitSet&gnNamePropIndexedGlobal != 0 {
					return n, true
				}
			}
		}
		return nil, true
	case rhino.Token_SUPER:
		if g.indexedGlobalCount != 0 {
			if nameID := g.lookupChildID(0, gnSuperStrID); nameID != 0 {
				if n := g.nameByID(nameID); n.propertyBitSet&gnNamePropIndexedGlobal != 0 {
					return n, true
				}
			}
		}
		return nil, true
	case rhino.Token_GETPROP:
		var segBuf [8]uint32
		segs := segBuf[:0]
		hasDot := false
		cur := qname
		for cur.IsGetProp() {
			segID := cur.GetStringID()
			if !hasDot && (segID <= rhino.RhinoStringPool_EmptyID || strings.IndexByte(strpool.String(strpool.ID(segID)), '.') >= 0) {
				hasDot = true
			}
			segs = append(segs, segID)
			cur = cur.GetFirstChild()
		}
		var topID uint32
		switch cur.GetToken() {
		case rhino.Token_NAME:
			topID = cur.GetStringID()
			if topID <= rhino.RhinoStringPool_EmptyID {
				return nil, false
			}
			if !hasDot && strings.IndexByte(strpool.String(strpool.ID(topID)), '.') >= 0 {
				hasDot = true
			}
		case rhino.Token_THIS:
			topID = gnThisStrID
		case rhino.Token_SUPER:
			topID = gnSuperStrID
		default:
			return nil, false
		}
		if g.indexedGlobalCount == 0 {
			return nil, true
		}
		if hasDot {
			fullName, ok := qname.GetQualifiedNameString()
			if !ok {
				return nil, false
			}
			return g.LookupName(fullName), true
		}
		curNameID := g.lookupChildID(0, topID)
		if curNameID == 0 {
			return nil, true
		}
		for i := len(segs) - 1; i >= 0; i-- {
			curNameID = g.lookupChildID(curNameID, segs[i])
			if curNameID == 0 {
				return nil, true
			}
		}
		if n := g.nameByID(curNameID); n.propertyBitSet&gnNamePropIndexedGlobal != 0 {
			return n, true
		}
		return nil, true
	default:
		return nil, false
	}
}

// lookupGlobalChildByID resolves an indexed global child property of parent by
// interned baseNameID without string concatenation.
func (g *GlobalNamespace) lookupGlobalChildByID(parent *GlobalNamespace_Name, baseNameID uint32) *GlobalNamespace_Name {
	if parent == nil || parent.firstPropID == 0 || baseNameID == 0 {
		return nil
	}
	if s := strpool.String(strpool.ID(baseNameID)); strings.IndexByte(s, '.') >= 0 {
		return g.LookupName(parent.GetFullName() + "." + s)
	}
	childID := g.lookupChildID(parent.selfID, baseNameID)
	if childID == 0 {
		return nil
	}
	if child := g.nameByID(childID); child.propertyBitSet&gnNamePropIndexedGlobal != 0 {
		return child
	}
	return nil
}

// Java: GlobalNamespace.setShouldTraverseScriptPredicate(Predicate<Node>)
func (g *GlobalNamespace) SetShouldTraverseScriptPredicate(shouldTraverseScript func(rhino.Node) bool) {
	g.shouldTraverseScript = shouldTraverseScript
}

// Java: GlobalNamespace.hasExternsRoot()
func (g *GlobalNamespace) HasExternsRoot() bool {
	return !g.externsRoot.IsNil()
}

// Java: GlobalNamespace.getRootNode()
func (g *GlobalNamespace) GetRootNode() rhino.Node {
	return g.root.GetParent()
}

// GetRootNodeOfName returns the root node of the scope in which the root of
// a qualified name is declared, or null.
//
// name is a variable name (e.g. "a"); s the scope in which the name is
// referenced. Returns the root node of the scope in which this is defined,
// or null if this is undeclared.
//
// Java: GlobalNamespace.getRootNode(String, Scope)
func (g *GlobalNamespace) GetRootNodeOfName(name string, s *Scope) rhino.Node {
	return g.getRootNodeOfTopName(GlobalNamespace_GetTopVarName(name), s)
}

func (g *GlobalNamespace) isProvidedRootID(topID uint32) bool {
	if !g.hasProvidedNames || topID == 0 {
		return false
	}
	id := g.lookupChildID(0, topID)
	return id != 0 && g.nameByID(id).GetBooleanProperty(GlobalNamespace_NameProp_IS_PROVIDED)
}

func (g *GlobalNamespace) getRootNodeOfTopName(name string, s *Scope) rhino.Node {
	v := s.GetVar(name)
	if v == nil {
		if g.hasProvidedNames && g.isProvidedRootID(uint32(strpool.Lookup(name))) {
			return g.globalRoot
		}
		return rhino.NilNode
	}
	if v.IsLocal() {
		return v.GetScopeRoot()
	}
	return g.globalRoot
}

// Java: GlobalNamespace.getParentScope()
func (g *GlobalNamespace) GetParentScope() rhino.StaticScope {
	return nil
}

// GetSlot returns the Name (a rhino.StaticSlot), or nil.
//
// Java: GlobalNamespace.getSlot(String)
func (g *GlobalNamespace) GetSlot(name string) rhino.StaticSlot {
	return g.GetOwnSlot(name)
}

// GetOwnSlot returns the Name (a rhino.StaticSlot), or nil.
//
// Java: GlobalNamespace.getOwnSlot(String)
func (g *GlobalNamespace) GetOwnSlot(name string) rhino.StaticSlot {
	g.EnsureGenerated()
	if n := g.LookupName(name); n != nil {
		return n
	}
	return nil
}

// GetTopmostScopeOfEventualDeclaration is StaticScope's default method,
// which GlobalNamespace inherits.
//
// Java: StaticScope.getTopmostScopeOfEventualDeclaration(String) (default method)
func (g *GlobalNamespace) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(g, name)
}

// Java: GlobalNamespace.getReferences(Name)
func (g *GlobalNamespace) GetReferences(slot *GlobalNamespace_Name) iter.Seq[*GlobalNamespace_Ref] {
	g.EnsureGenerated()
	return slot.Refs()
}

// Java: GlobalNamespace.getScope(Name)
func (g *GlobalNamespace) GetScope(slot *GlobalNamespace_Name) rhino.StaticScope {
	return g
}

// AllNames iterates all global names in insertion order without materializing
// a string-keyed LinkedHashMap.
func (g *GlobalNamespace) AllNames() iter.Seq[*GlobalNamespace_Name] {
	g.EnsureGenerated()
	return func(yield func(*GlobalNamespace_Name) bool) {
		limit := g.nameCount
		for idx := uint32(0); idx < limit; idx++ {
			n := &g.namePages[idx>>gnNamePageShift][idx&gnNamePageMask]
			if n.propertyBitSet&gnNamePropIndexedGlobal != 0 {
				if !yield(n) {
					return
				}
			}
		}
	}
}

// Java: GlobalNamespace.getAllSymbols()
func (g *GlobalNamespace) GetAllSymbols() iter.Seq[*GlobalNamespace_Name] {
	return g.AllNames()
}

// Java: GlobalNamespace.ensureGenerated()
func (g *GlobalNamespace) EnsureGenerated() {
	if !g.generated {
		g.ProcessNamespace()
	}
}

// NameForest iterates the top-level roots of the global name forest without
// allocating an ArrayList.
func (g *GlobalNamespace) NameForest() iter.Seq[*GlobalNamespace_Name] {
	g.EnsureGenerated()
	return func(yield func(*GlobalNamespace_Name) bool) {
		for _, id := range g.globalNameIDs {
			if !yield(g.nameByID(id)) {
				return
			}
		}
	}
}

// NameForestSlice returns a freshly allocated slice of the top-level roots of
// the global name forest.
func (g *GlobalNamespace) NameForestSlice() []*GlobalNamespace_Name {
	g.EnsureGenerated()
	out := make([]*GlobalNamespace_Name, len(g.globalNameIDs))
	for i, id := range g.globalNameIDs {
		out[i] = g.nameByID(id)
	}
	return out
}

// GetNameForest gets a list of the roots of the forest of the global names,
// where the roots are the top-level names.
//
// Java: GlobalNamespace.getNameForest()
func (g *GlobalNamespace) GetNameForest() []*GlobalNamespace_Name {
	return g.NameForestSlice()
}

// GetNameIndex gets an index of all the global names, indexed by full
// qualified name (as in "a", "a.b.c", etc.).
//
// Java: GlobalNamespace.getNameIndex()
func (g *GlobalNamespace) GetNameIndex() *jutil.LinkedHashMap[string, *GlobalNamespace_Name] {
	g.EnsureGenerated()
	m := jutil.NewLinkedHashMapWithCapacity[string, *GlobalNamespace_Name](int(g.indexedGlobalCount))
	for n := range g.AllNames() {
		m.PutNew(n.GetFullName(), n)
	}
	return m
}

// GlobalNamespace_AstChange records a node added or modified in the AST so
// GlobalNamespace can rescan it.
type GlobalNamespace_AstChange struct {
	chunk *JSChunk
	scope *Scope
	node  rhino.Node
}

// Java: GlobalNamespace.SimpleAstChange(Node, JSChunk, Scope)
func NewGlobalNamespace_SimpleAstChange(node rhino.Node, chunk *JSChunk, scope *Scope) GlobalNamespace_AstChange {
	return GlobalNamespace_AstChange{chunk: chunk, scope: scope, node: node}
}

// Java: GlobalNamespace.RefBasedAstChange(Ref, Node)
func NewGlobalNamespace_RefBasedAstChange(ref *GlobalNamespace_Ref, node rhino.Node) GlobalNamespace_AstChange {
	return GlobalNamespace_AstChange{chunk: ref.chunk, scope: ref.scope, node: node}
}

// Java: GlobalNamespace.AstChange.node()
func (c GlobalNamespace_AstChange) Node() rhino.Node { return c.node }

// Java: GlobalNamespace.AstChange.chunk()
func (c GlobalNamespace_AstChange) Chunk() *JSChunk { return c.chunk }

// Java: GlobalNamespace.AstChange.scope()
func (c GlobalNamespace_AstChange) Scope() *Scope { return c.scope }

type astChangeSet struct {
	items []GlobalNamespace_AstChange
	seen  map[GlobalNamespace_AstChange]struct{}
}

func (s *astChangeSet) reset() {
	clear(s.items)
	s.items = s.items[:0]
	if len(s.seen) > 0 {
		clear(s.seen)
	}
}

func (s *astChangeSet) add(c GlobalNamespace_AstChange) {
	if s.seen != nil && len(s.items) > 8 {
		if _, ok := s.seen[c]; ok {
			return
		}
		s.seen[c] = struct{}{}
		s.items = append(s.items, c)
		return
	}
	for i := range s.items {
		if s.items[i] == c {
			return
		}
	}
	if len(s.items) == 8 {
		if s.seen == nil {
			s.seen = make(map[GlobalNamespace_AstChange]struct{}, 16)
		}
		for _, prev := range s.items {
			s.seen[prev] = struct{}{}
		}
		s.seen[c] = struct{}{}
	}
	s.items = append(s.items, c)
}

// ScanNewNodes: if the client adds new nodes to the AST, scan these new
// nodes to see if they've added any references to the global namespace.
//
// newNodes are the new nodes to check.
//
// Java: GlobalNamespace.scanNewNodes(Set<AstChange>)
func (g *GlobalNamespace) ScanNewNodes(newNodes *astChangeSet) {
	if newNodes == nil || len(newNodes.items) == 0 {
		return
	}
	builder := NewGlobalNamespace_BuildGlobalNamespace(g)

	for i := range newNodes.items {
		info := newNodes.items[i]
		if !info.node.IsQualifiedName() && !NodeUtil_MayBeObjectLitKey(info.node) {
			continue
		}
		g.ScanFromNode(builder, info.scope, info.node, info.chunk)
	}
}

// Java: GlobalNamespace.scanFromNode(BuildGlobalNamespace, Scope, Node, JSChunk)
func (g *GlobalNamespace) ScanFromNode(builder *GlobalNamespace_BuildGlobalNamespace, scope *Scope, n rhino.Node, chunk *JSChunk) {
	// Check affected parent nodes first.
	parent := n.GetParent()
	if (n.IsName() || n.IsGetProp()) && parent.IsGetProp() {
		// e.g. when replacing "my.alias.prop" with "foo.bar.prop"
		// we want also want to visit "foo.bar.prop", since that's a new global qname we are now
		// referencing.
		g.ScanFromNode(builder, scope, n.GetParent(), chunk)
	} else if !n.GetPrevious().IsNil() && n.GetPrevious().IsObjectPattern() {
		// e.g. if we change `const {x} = bar` to `const {x} = foo`, add a new reference to `foo.x`
		// attached to the STRING_KEY `x`
		pattern := n.GetPrevious()
		for key := pattern.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
			if key.IsStringKey() {
				g.ScanFromNode(builder, scope, key, chunk)
			}
		}
	}
	builder.Collect(scope, chunk, n)
}

// ProcessNamespace builds the namespace lazily.
//
// Go-only: language-adapter — the private Java method process() is
// ProcessNamespace, so that GlobalNamespace is not mistaken for a
// CompilerPass (process(Node, Node)).
//
// Java: GlobalNamespace.process()
func (g *GlobalNamespace) ProcessNamespace() {
	isInitial := !g.generated && g.childCount == 0 && g.nameMapByModule == nil
	if g.forProcessDefines && isInitial {
		g.collectDefineCandidates()
	}
	if len(g.childTable) == 0 {
		capHint := g.estimateNameCapacity()
		if g.defineCandidateLeafIDs != nil {
			capHint = 2 * len(g.defineCandidateLeafIDs)
		}
		if capHint > 64 {
			g.initChildTable(capHint)
		}
	}
	builder := NewGlobalNamespace_BuildGlobalNamespace(g)
	builder.confirmExisting = !isInitial
	traversal := NodeTraversal_Builder_().SetCompiler(g.compiler).SetCallback(builder)

	if g.HasExternsRoot() {
		traversal.TraverseRoots(g.externsRoot, g.root)
	} else {
		traversal.Traverse(g.root)
	}

	g.generated = true
}

// collectDefineCandidates fills defineCandidateLeafIDs: the leaf segment ID of
// every NAME or GETPROP whose best JSDoc is @define, plus, in check mode,
// every one whose best JSDoc is constant or that a `const` declares.
//
// ProcessDefines reads a name only when its declaration carries @define
// (selectDefineDeclaration) or when it is a global constant with a `const`
// or @const declaration (isGlobalConst + getConstantDeclValue), and then
// only that name's own sets and gets. The names collected here are a
// superset of those (local declarations and reads are not told apart), so
// skipping every other name leaves ProcessDefines with the same defines,
// the same constants and the same refs in the same order.
//
// Go-only: performance.
func (g *GlobalNamespace) collectDefineCandidates() {
	g.defineCandidateLeafIDs = make(map[uint32]struct{}, 1024)
	check := g.processDefinesCheck
	var walk func(n rhino.Node)
	walk = func(n rhino.Node) {
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			walk(c)
		}
		switch n.GetToken() {
		case rhino.Token_NAME:
			if check && !n.GetParent().IsNil() && n.GetParent().IsConst() {
				g.defineCandidateLeafIDs[n.GetStringID()] = struct{}{}
				return
			}
		case rhino.Token_GETPROP:
		default:
			return
		}
		info := NodeUtil_GetBestJSDocInfo(n)
		if info == nil || !(info.IsDefine() || (check && info.IsConstant())) {
			return
		}
		if n.GetToken() == rhino.Token_NAME || n.IsQualifiedName() {
			g.defineCandidateLeafIDs[n.GetStringID()] = struct{}{}
		}
	}
	if g.HasExternsRoot() {
		walk(g.externsRoot)
	}
	walk(g.root)
}

func (b *GlobalNamespace_BuildGlobalNamespace) skipsLeafID(leafID uint32) bool {
	if b.outer.defineCandidateLeafIDs == nil {
		return false
	}
	if leafID == 0 {
		return true
	}
	_, ok := b.outer.defineCandidateLeafIDs[leafID]
	return !ok
}

// skipsName reports whether the namespace, built for ProcessDefines, leaves
// out the given name (see collectDefineCandidates).
func (b *GlobalNamespace_BuildGlobalNamespace) skipsName(name string) bool {
	if b.outer.defineCandidateLeafIDs == nil {
		return false
	}
	lastSeg := name
	if i := strings.LastIndexByte(name, '.'); i >= 0 {
		lastSeg = name[i+1:]
	}
	return b.skipsLeafID(uint32(strpool.Lookup(lastSeg)))
}

// estimateNameCapacity guesses the number of names from the number of
// top-level statements in the code (externs hold few qualified names).
//
// Go-only: performance — growing the name index to ~150k names (the
// typed-full benchmark) allocates about twice its final size.
func (g *GlobalNamespace) estimateNameCapacity() int {
	if g.root.IsNil() {
		return 0
	}
	stmtCount := 0
	for s := g.root.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
		if g.shouldTraverseScript != nil && !g.shouldTraverseScript(s) {
			continue
		}
		for c := s.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if c.IsModuleBody() {
				stmtCount += c.GetChildCount()
			} else {
				stmtCount++
			}
		}
	}
	return stmtCount * 5
}

// GlobalNamespace_GetTopVarName gets the top variable name from a possibly
// namespaced name.
//
// name is a variable or qualified property name (e.g. "a" or "a.b.c.d").
// Returns the top variable name (e.g. "a").
//
// Java: GlobalNamespace.getTopVarName(String)
func GlobalNamespace_GetTopVarName(name string) string {
	firstDotIndex := strings.IndexByte(name, '.')
	if firstDotIndex == -1 {
		return name
	}
	return name[:firstDotIndex]
}

type globalNamespaceModuleKey struct {
	module *ModuleMetadataMap_ModuleMetadata
	nameID uint32
}

// Java: GlobalNamespace.getNameFromModule(ModuleMetadata, String)
func (g *GlobalNamespace) GetNameFromModule(moduleMetadata *ModuleMetadataMap_ModuleMetadata, name string) *GlobalNamespace_Name {
	jutil.CheckNotNull(moduleMetadata)
	g.EnsureGenerated()
	if g.nameMapByModule == nil || name == "" {
		return nil
	}
	dot := strings.IndexByte(name, '.')
	var rootSeg string
	if dot < 0 {
		rootSeg = name
	} else {
		rootSeg = name[:dot]
	}
	rootStrID := uint32(strpool.Lookup(rootSeg))
	if rootStrID == 0 {
		return nil
	}
	curID := g.nameMapByModule[globalNamespaceModuleKey{module: moduleMetadata, nameID: rootStrID}]
	if curID == 0 {
		return nil
	}
	for dot >= 0 {
		rem := name[dot+1:]
		nextDot := strings.IndexByte(rem, '.')
		var seg string
		if nextDot < 0 {
			seg = rem
			dot = -1
		} else {
			seg = rem[:nextDot]
			dot = dot + 1 + nextDot
		}
		segID := uint32(strpool.Lookup(seg))
		if segID == 0 {
			return nil
		}
		curID = g.lookupChildID(curID, segID)
		if curID == 0 {
			return nil
		}
	}
	return g.nameByID(curID)
}

// DeclarationHasFollowingObjectSpreadSibling returns whether a declaration
// node, inside an object-literal, has a following OBJECT_SPREAD sibling.
//
// This check is implemented using a cache because otherwise it has aggregate
// `O(n^2)` performance in terms of the size of an OBJECTLIT. If each
// declaration checked each sibling independently, each sibling would be
// checked up-to once for each of it's preceeding siblings.
//
// Java: GlobalNamespace.declarationHasFollowingObjectSpreadSibling(Node)
func (g *GlobalNamespace) DeclarationHasFollowingObjectSpreadSibling(declaration rhino.Node) bool {
	objLit := declaration.GetParent()
	jutil.CheckStateMsg1(objLit.IsObjectLit(), "%s", declaration)

	if g.spreadSiblingCache != nil {
		if cached, ok := g.spreadSiblingCache[objLit]; ok && !cached {
			return false
		}
		if cached, ok := g.spreadSiblingCache[declaration]; ok {
			return cached
		}
	} else {
		g.spreadSiblingCache = make(map[rhino.Node]bool, 64)
	}

	hasSpread := false
	for sibling := objLit.GetLastChild(); !sibling.IsNil(); sibling = sibling.GetPrevious() {
		if sibling.IsSpread() {
			hasSpread = true
			break
		}
	}
	if !hasSpread {
		g.spreadSiblingCache[objLit] = false
		return false
	}
	g.spreadSiblingCache[objLit] = true

	/*
	 * Iterate backward over all children of the object-literal, filling in the cache.
	 *
	 * <p>We iterate the entire literal because we expect to eventually need the result for each of
	 * them. Additionally, it makes the loop conditions simpler.
	 *
	 * <p>We use a loop rather than recursion to minimize stack depth. Large object-literals were
	 * the reason caching was added.
	 */
	toCache := false
	for sibling := objLit.GetLastChild(); !sibling.IsNil(); sibling = sibling.GetPrevious() {
		if sibling.IsSpread() {
			toCache = true
		}
		g.spreadSiblingCache[sibling] = toCache
	}

	return g.spreadSiblingCache[declaration]
}
