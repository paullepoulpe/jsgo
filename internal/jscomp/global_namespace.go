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
	"strconv"
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
	gnNamePropRefsMapState  int32 = 1 << 29
	gnNamePropIndexedGlobal int32 = 1 << 30
)

var (
	gnPrototypeStrID = uint32(strpool.Intern("prototype"))
	gnThisStrID      = uint32(strpool.Intern("this"))
	gnSuperStrID     = uint32(strpool.Intern("super"))
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
	return NewGlobalNamespaceWithLogAndRoots( /* decisionsLog= */ nil, compiler, rhino.NilNode, root)
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
	return NewGlobalNamespaceWithLogAndRoots(nil, compiler, externsRoot, root)
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
	g.refCount = idx + 1
	pageIdx := idx >> gnRefPageShift
	if int(pageIdx) >= len(g.refPages) {
		g.refPages = append(g.refPages, new([gnRefPageLen]GlobalNamespace_Ref))
	}
	r := &g.refPages[pageIdx][idx&gnRefPageMask]
	*r = GlobalNamespace_Ref{chunk: chunk, scope: scope, node: node, type_: type_}
	return idx + 1, r
}

func (g *GlobalNamespace) findOrCopyRefID(ref *GlobalNamespace_Ref) uint32 {
	if ref == nil {
		return 0
	}
	for idx := g.refCount; idx > 0; idx-- {
		i := idx - 1
		if &g.refPages[i>>gnRefPageShift][i&gnRefPageMask] == ref {
			return idx
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

// GlobalNamespace_AstChange is the interface GlobalNamespace.AstChange.
type GlobalNamespace_AstChange interface {
	// Java: GlobalNamespace.AstChange.node()
	Node() rhino.Node

	// Java: GlobalNamespace.AstChange.chunk()
	Chunk() *JSChunk

	// Java: GlobalNamespace.AstChange.scope()
	Scope() *Scope
}

// GlobalNamespace_SimpleAstChange is the record SimpleAstChange(Node node,
// JSChunk chunk, Scope scope).
//
// Go-only: language-adapter — a comparable value type, so that two equal
// records are equal keys of a set of AstChange, as in Java.
type GlobalNamespace_SimpleAstChange struct {
	node  rhino.Node
	chunk *JSChunk
	scope *Scope
}

// Java: GlobalNamespace.SimpleAstChange(Node, JSChunk, Scope)
func NewGlobalNamespace_SimpleAstChange(node rhino.Node, chunk *JSChunk, scope *Scope) GlobalNamespace_SimpleAstChange {
	return GlobalNamespace_SimpleAstChange{node: node, chunk: chunk, scope: scope}
}

// Java: GlobalNamespace.SimpleAstChange.node()
func (c GlobalNamespace_SimpleAstChange) Node() rhino.Node { return c.node }

// Java: GlobalNamespace.SimpleAstChange.chunk()
func (c GlobalNamespace_SimpleAstChange) Chunk() *JSChunk { return c.chunk }

// Java: GlobalNamespace.SimpleAstChange.scope()
func (c GlobalNamespace_SimpleAstChange) Scope() *Scope { return c.scope }

// GlobalNamespace_RefBasedAstChange is a simple data class that contains the
// information necessary to inspect a node for changes to the global
// namespace: the record RefBasedAstChange(Ref ref, Node node).
//
// Go-only: language-adapter — a comparable value type, so that two equal
// records are equal keys of a set of AstChange, as in Java.
type GlobalNamespace_RefBasedAstChange struct {
	ref  *GlobalNamespace_Ref
	node rhino.Node
}

// Java: GlobalNamespace.RefBasedAstChange(Ref, Node)
func NewGlobalNamespace_RefBasedAstChange(ref *GlobalNamespace_Ref, node rhino.Node) GlobalNamespace_RefBasedAstChange {
	return GlobalNamespace_RefBasedAstChange{ref: ref, node: node}
}

// Java: GlobalNamespace.RefBasedAstChange.ref()
func (c GlobalNamespace_RefBasedAstChange) Ref() *GlobalNamespace_Ref { return c.ref }

// Java: GlobalNamespace.RefBasedAstChange.node()
func (c GlobalNamespace_RefBasedAstChange) Node() rhino.Node { return c.node }

// Java: GlobalNamespace.RefBasedAstChange.chunk()
func (c GlobalNamespace_RefBasedAstChange) Chunk() *JSChunk {
	return c.ref.GetChunk()
}

// Java: GlobalNamespace.RefBasedAstChange.scope()
func (c GlobalNamespace_RefBasedAstChange) Scope() *Scope {
	return c.ref.scope
}

// ScanNewNodes: if the client adds new nodes to the AST, scan these new
// nodes to see if they've added any references to the global namespace.
//
// newNodes are the new nodes to check.
//
// Java: GlobalNamespace.scanNewNodes(Set<AstChange>)
func (g *GlobalNamespace) ScanNewNodes(newNodes *jutil.LinkedHashSet[GlobalNamespace_AstChange]) {
	builder := NewGlobalNamespace_BuildGlobalNamespace(g)

	for info := range newNodes.All() {
		if !info.Node().IsQualifiedName() && !NodeUtil_MayBeObjectLitKey(info.Node()) {
			continue
		}
		g.ScanFromNode(builder, info.Scope(), info.Node(), info.Chunk())
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

// -------------------------------------------------------------------------

type globalNamespace_objLitCacheEntry struct {
	prefix string
	root   rhino.Node
	scope  *Scope
	state  uint8 // 0: invalid, 1: LValue, 2: Object.defineProperties
}

// GlobalNamespace_BuildGlobalNamespace builds a tree representation of the
// global namespace. Omits prototypes.
type GlobalNamespace_BuildGlobalNamespace struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer *GlobalNamespace

	curModuleRoot rhino.Node
	curMetadata   *ModuleMetadataMap_ModuleMetadata

	lastName       string
	lastMeta       *ModuleMetadataMap_ModuleMetadata
	lastNode       *GlobalNamespace_Name
	lastParentName string
	lastParentMeta *ModuleMetadataMap_ModuleMetadata
	lastParentNode *GlobalNamespace_Name

	objLitCache     map[rhino.Node]globalNamespace_objLitCacheEntry
	confirmExisting bool
}

// Java: GlobalNamespace.BuildGlobalNamespace()
func NewGlobalNamespace_BuildGlobalNamespace(outer *GlobalNamespace) *GlobalNamespace_BuildGlobalNamespace {
	return &GlobalNamespace_BuildGlobalNamespace{outer: outer, curModuleRoot: rhino.NilNode, curMetadata: nil, confirmExisting: true}
}

func globalNamespace_topNameOfQName(n rhino.Node) string {
	if n.IsNil() {
		return ""
	}
	left := n
	for left.IsGetProp() {
		left = left.GetFirstChild()
	}
	switch left.GetToken() {
	case rhino.Token_NAME:
		return GlobalNamespace_GetTopVarName(left.GetString())
	case rhino.Token_THIS:
		return "this"
	case rhino.Token_SUPER:
		return "super"
	default:
		return ""
	}
}

func (b *GlobalNamespace_BuildGlobalNamespace) putObjLitCache(objLit rhino.Node, prefix string, root rhino.Node, scope *Scope, state uint8) (string, rhino.Node, *Scope, uint8) {
	if b.objLitCache == nil {
		b.objLitCache = make(map[rhino.Node]globalNamespace_objLitCacheEntry, 16)
	}
	b.objLitCache[objLit] = globalNamespace_objLitCacheEntry{
		prefix: prefix,
		root:   root,
		scope:  scope,
		state:  state,
	}
	return prefix, root, scope, state
}

func globalNamespace_hasPrototype(s string) bool {
	return strings.HasSuffix(s, ".prototype") || strings.Contains(s, ".prototype.")
}

func globalNamespace_nodeHasPrototype(n rhino.Node) bool {
	for cur := n; !cur.IsNil() && cur.IsGetProp(); cur = cur.GetFirstChild() {
		if cur.GetString() == "prototype" {
			return true
		}
	}
	return false
}

func (b *GlobalNamespace_BuildGlobalNamespace) resolveObjectLit(objLit rhino.Node, scope *Scope, t *NodeTraversal) (string, rhino.Node, *Scope, uint8) {
	if e, ok := b.objLitCache[objLit]; ok {
		return e.prefix, e.root, e.scope, e.state
	}

	g := b.outer
	owner := NodeUtil_GetBestLValue(objLit)
	if owner.IsNil() {
		objLitParent := objLit.GetParent()
		if NodeUtil_IsObjectDefinePropertiesDefinition(objLitParent) {
			receiverNode := objLitParent.GetSecondChild()
			if globalNamespace_nodeHasPrototype(receiverNode) {
				return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
			}
			topName := globalNamespace_topNameOfQName(receiverNode)
			if topName == "" {
				return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
			}
			jutil.CheckStateMsg1(objLit == receiverNode.GetNext(), "%s", objLitParent)
			root := b.getRootNodeOfTopNameWithTraversal(topName, scope, t)
			if !b.IsTopLevelScopeRoot(root) {
				return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
			}
			qname, ok := receiverNode.GetQualifiedNameString()
			if !ok || globalNamespace_hasPrototype(qname) {
				return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
			}
			if scope == nil && !g.forProcessDefines {
				scope = t.GetScope()
			}
			return b.putObjLitCache(objLit, qname, root, scope, 2)
		}
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
	}

	ownerParent := owner.GetParent()
	if ownerParent.IsNil() {
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
	}
	if ownerParent.IsClassMembers() && !owner.IsComputedProp() {
		if !owner.IsStaticMember() || owner.GetString() == "prototype" {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
		}
		nameNode := NodeUtil_GetNameNode(ownerParent.GetParent())
		if globalNamespace_nodeHasPrototype(nameNode) {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
		}
		topName := globalNamespace_topNameOfQName(nameNode)
		if topName == "" {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
		}
		root := b.getRootNodeOfTopNameWithTraversal(topName, scope, t)
		if !b.IsTopLevelScopeRoot(root) {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
		}
		className, ok := nameNode.GetQualifiedNameString()
		if !ok {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
		}
		ownerName := className + "." + owner.GetString()
		if globalNamespace_hasPrototype(ownerName) {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
		}
		if scope == nil && !g.forProcessDefines {
			scope = t.GetScope()
		}
		return b.putObjLitCache(objLit, ownerName, root, scope, 1)
	}
	if NodeUtil_MayBeObjectLitKey(owner) {
		if !ownerParent.IsObjectLit() {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
		}
		parentPrefix, root, resolvedScope, parentState := b.resolveObjectLit(ownerParent, scope, t)
		if parentState != 1 {
			return b.putObjLitCache(objLit, "", rhino.NilNode, resolvedScope, 0)
		}
		key := NodeUtil_GetObjectOrClassLitKeyName(owner)
		if key == "prototype" || !rhino.TokenStream_IsJSIdentifier(key) {
			return b.putObjLitCache(objLit, "", rhino.NilNode, resolvedScope, 0)
		}
		return b.putObjLitCache(objLit, parentPrefix+"."+key, root, resolvedScope, 1)
	}

	if owner.IsName() && t != nil && b.curModuleRoot.IsNil() && !t.currentHoistScopeRoot.IsNil() {
		switch ownerParent.GetToken() {
		case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
		}
	}
	if globalNamespace_nodeHasPrototype(owner) {
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
	}
	topName := globalNamespace_topNameOfQName(owner)
	if topName == "" {
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
	}
	root := b.getRootNodeOfTopNameWithTraversal(topName, scope, t)
	if !b.IsTopLevelScopeRoot(root) {
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
	}
	ownerName, ok := owner.GetQualifiedNameString()
	if !ok || globalNamespace_hasPrototype(ownerName) {
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, 0)
	}
	if scope == nil && !g.forProcessDefines {
		scope = t.GetScope()
	}
	return b.putObjLitCache(objLit, ownerName, root, scope, 1)
}

func (b *GlobalNamespace_BuildGlobalNamespace) getRootNodeOfTopName(topName string, s *Scope) rhino.Node {
	return b.outer.getRootNodeOfTopName(topName, s)
}

func (b *GlobalNamespace_BuildGlobalNamespace) getRootNodeOfTopNameWithTraversal(topName string, s *Scope, t *NodeTraversal) rhino.Node {
	return b.getRootNodeOfTopNameByID(0, topName, s, t)
}

func (b *GlobalNamespace_BuildGlobalNamespace) getRootNodeOfTopNameByID(topID uint32, topName string, s *Scope, t *NodeTraversal) rhino.Node {
	if s != nil || t == nil {
		return b.getRootNodeOfTopName(topName, s)
	}
	g := b.outer
	if b.curModuleRoot.IsNil() && len(t.scopes) > 0 {
		globalScope := t.scopes[0].untypedScope
		if globalScope == nil {
			globalScope = t.getScopeSlow(0)
		}
		if globalScope.GetVar(topName) == nil {
			if !g.hasProvidedNames {
				return rhino.NilNode
			}
			if topID == 0 {
				topID = uint32(strpool.Lookup(topName))
			}
			if !g.isProvidedRootID(topID) {
				return rhino.NilNode
			}
		}
		if len(t.scopes) == 1 {
			return g.globalRoot
		}
	}
	v := t.LookupVar(topName)
	if v == nil {
		if g.hasProvidedNames {
			if topID == 0 {
				topID = uint32(strpool.Lookup(topName))
			}
			if g.isProvidedRootID(topID) {
				return g.globalRoot
			}
		}
		return rhino.NilNode
	}
	if v.IsLocal() {
		return v.GetScopeRoot()
	}
	return g.globalRoot
}

// ShouldTraverse collects the references in pre-order.
//
// Java: GlobalNamespace.BuildGlobalNamespace.shouldTraverse(NodeTraversal, Node, Node)
func (b *GlobalNamespace_BuildGlobalNamespace) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	g := b.outer
	if n.IsScript() && g.shouldTraverseScript != nil && !g.shouldTraverseScript(n) {
		return false
	}
	if g.HasExternsRoot() && n.IsScript() {
		// When checking type-summary files, we want to consider them like normal code
		// for some things (like alias inlining) but like externs for other things.
		g.sourceKind = GlobalNamespace_SourceKind_FromScriptNode(n)
	} else if n == g.root {
		g.sourceKind = GlobalNamespace_SourceKind_CODE
	}
	if g.forProcessDefines && g.processDefinesCheck && g.sourceKind == GlobalNamespace_SourceKind_CODE {
		if jsdoc := n.GetJSDocInfo(); (jsdoc != nil && jsdoc.IsDefine()) ||
			(n.IsCall() && n.GetFirstChild().MatchesQualifiedNameNode(processDefines_GOOG_DEFINE)) ||
			(n.IsName() && n.GetString() == "CLOSURE_DEFINES") {
			if !t.currentScript.IsNil() {
				if g.defineCandidateScripts == nil {
					g.defineCandidateScripts = make(map[rhino.Node]struct{}, 8)
				}
				g.defineCandidateScripts[t.currentScript] = struct{}{}
			} else {
				g.hasNonScriptDefineCand = true
			}
		}
	}
	if n.IsModuleBody() || NodeUtil_IsBundledGoogModuleScopeRoot(n) {
		b.SetupModuleMetadata(n)
	} else if n.IsScript() || NodeUtil_IsBundledGoogModuleCall(n) {
		b.curModuleRoot = rhino.NilNode
		b.curMetadata = nil
	}

	b.collect(nil, nil, t, n)

	return true
}

// SetupModuleMetadata initializes the ModuleMetadata for a goog;.module or
// ES module.
//
// moduleRoot is either a MODULE_BODY or a goog.loadModule BLOCK.
//
// Java: GlobalNamespace.BuildGlobalNamespace.setupModuleMetadata(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) SetupModuleMetadata(moduleRoot rhino.Node) {
	moduleMap := b.outer.compiler.GetModuleMap()
	if moduleMap == nil {
		return
	}
	b.curModuleRoot = moduleRoot

	b.curMetadata = jutil.CheckNotNull(
		ModuleImportResolver_GetModuleFromScopeRoot(moduleMap, b.outer.compiler, moduleRoot).
			Metadata())
	if b.curMetadata.IsGoogModule() {
		b.GetOrCreateName("exports", b.curMetadata)
	}
}

// Java: GlobalNamespace.BuildGlobalNamespace.collect(Scope, JSChunk, Node)
func (b *GlobalNamespace_BuildGlobalNamespace) Collect(scope *Scope, chunk *JSChunk, n rhino.Node) {
	b.collect(scope, chunk, nil, n)
}

func (b *GlobalNamespace_BuildGlobalNamespace) collect(scope *Scope, chunk *JSChunk, t *NodeTraversal, n rhino.Node) {
	g := b.outer
	parent := n.GetParent()

	var nameStr string
	var resolvedName *GlobalNamespace_Name
	var root rhino.Node
	hasName := false
	isSet := false
	type_ := GlobalNamespace_NameProp_OTHER_OBJECT

	switch n.GetToken() {
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF:
		if parent.IsObjectLit() {
			prefix, r, s, state := b.resolveObjectLit(parent, scope, t)
			if state != 1 {
				return
			}
			key := n.GetString()
			if key == "prototype" || !rhino.TokenStream_IsJSIdentifier(key) {
				return
			}
			nameStr = prefix + "." + key
			root = r
			scope = s
			hasName = true
		} else if parent.IsClassMembers() {
			if !n.IsStaticMember() || n.IsComputedProp() {
				return // Within a class, only static members define global names.
			}
			nameNode := NodeUtil_GetNameNode(parent.GetParent())
			topName := globalNamespace_topNameOfQName(nameNode)
			if topName == "" {
				return
			}
			root = b.getRootNodeOfTopNameWithTraversal(topName, scope, t)
			if !b.IsTopLevelScopeRoot(root) {
				return
			}
			className, ok := nameNode.GetQualifiedNameString()
			if !ok {
				return
			}
			nameStr = className + "." + n.GetString()
			hasName = true
		} else {
			nameStr, hasName = nodeUtil_getBestLValueNameString(n)
			if !hasName {
				return
			}
			root = b.getRootNodeOfTopNameWithTraversal(GlobalNamespace_GetTopVarName(nameStr), scope, t)
			if !b.IsTopLevelScopeRoot(root) {
				return
			}
		}
		isSet = true

		if n.IsMemberFunctionDef() {
			type_ = GlobalNamespace_NameProp_FUNCTION
		} else if n.IsMemberFieldDef() {
			type_ = GlobalNamespace_NameProp_OTHER_OBJECT
		} else {
			type_ = GlobalNamespace_NameProp_GET_SET
		}
	case rhino.Token_STRING_KEY:
		if parent.IsObjectLit() {
			prefix, r, s, state := b.resolveObjectLit(parent, scope, t)
			if state == 0 {
				return
			}
			key := n.GetString()
			if key == "prototype" {
				return
			}
			if state == 1 {
				if !rhino.TokenStream_IsJSIdentifier(key) {
					return
				}
				nameStr = prefix + "." + key
				type_ = b.GetValueType(n.GetOnlyChild())
			} else {
				nameStr = prefix + "." + key
				type_ = GlobalNamespace_NameProp_GET_SET
			}
			root = r
			scope = s
			hasName = true
			isSet = true
		} else if parent.IsObjectPattern() {
			if g.forProcessDefines && (!g.processDefinesCheck || (t != nil && !t.currentHoistScopeRoot.IsNil() && b.curModuleRoot.IsNil())) {
				return
			}
			patternParent := parent.GetParent()
			if !patternParent.IsAssign() && !patternParent.IsDestructuringLhs() {
				return
			}
			rhs := patternParent.GetSecondChild()
			topName := globalNamespace_topNameOfQName(rhs)
			if topName == "" {
				return
			}
			root = b.getRootNodeOfTopNameWithTraversal(topName, scope, t)
			if !b.IsTopLevelScopeRoot(root) {
				return
			}
			qname, ok := rhs.GetQualifiedNameString()
			if !ok {
				return
			}
			nameStr = qname + "." + n.GetString()
			type_ = b.GetValueType(n.GetFirstChild())
			hasName = true
			// not a set
		} else {
			// else not a reference we should record
			return
		}
	case rhino.Token_NAME, rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		// "a.b?.c" is not a reference to the global name "a.b.c" for the
		// purposes of GlobalNamespace, but we include OPTCHAIN_GETPROP to
		// detect optional chain hasOwnProperty guards.
		pTok := parent.GetToken()
		if pTok == rhino.Token_GETPROP {
			// This name is nested in a getprop. Return and only create a Ref for the outermost
			// getprop in the chain.
			return
		}
		if !g.forProcessDefines && (pTok == rhino.Token_CALL || pTok == rhino.Token_OPTCHAIN_CALL) &&
			n.IsFirstChildOf(parent) && b.IsObjectHasOwnPropertyCall(parent) {
			qname, _ := n.GetFirstChild().GetQualifiedNameString()
			globalName := b.GetOrCreateName(qname, b.curMetadata)
			globalName.SetBooleanProperty(GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY)
		}
		if n.GetToken() == rhino.Token_OPTCHAIN_GETPROP {
			return
		}
		// TODO(b/127505242): CAST parents may indicate a set.
		// This may be a variable get or set.
		switch pTok {
		case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
			isSet = true
			rvalue := n.GetFirstChild()
			if rvalue.IsNil() {
				type_ = GlobalNamespace_NameProp_OTHER_OBJECT
			} else {
				type_ = b.GetValueType(rvalue)
			}
		case rhino.Token_ASSIGN:
			if parent.GetFirstChild() == n {
				isSet = true
				type_ = b.GetValueType(n.GetNext())
			}
		case rhino.Token_FUNCTION:
			grandparent := parent.GetParent()
			if grandparent.IsNil() || NodeUtil_IsFunctionExpression(parent) {
				return
			}
			isSet = true
			type_ = GlobalNamespace_NameProp_FUNCTION
		case rhino.Token_CATCH, rhino.Token_INC, rhino.Token_DEC:
			isSet = true
			type_ = GlobalNamespace_NameProp_OTHER_OBJECT
		case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
			if parent.GetFirstChild() == n {
				isSet = true
				type_ = GlobalNamespace_NameProp_OTHER_OBJECT
			}
		case rhino.Token_CLASS:
			// The first child is the class name, and the second child is the superclass name.
			if parent.GetFirstChild() == n {
				isSet = true
				type_ = GlobalNamespace_NameProp_CLASS
			}
		case rhino.Token_STRING_KEY,
			rhino.Token_ARRAY_PATTERN,
			rhino.Token_DEFAULT_VALUE,
			rhino.Token_COMPUTED_PROP,
			rhino.Token_ITER_REST,
			rhino.Token_OBJECT_REST:
			// This may be a set.
			// TODO(b/120303257): this should extend to qnames too, but doing
			// so causes invalid output. Covered in CollapsePropertiesTest
			if n.IsName() && NodeUtil_IsLhsByDestructuring(n) {
				isSet = true
				type_ = GlobalNamespace_NameProp_OTHER_OBJECT
			}
		case rhino.Token_ITER_SPREAD, rhino.Token_OBJECT_SPREAD, rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL:
			// isSet = false, type = OTHER.
		default:
			if NodeUtil_IsAssignmentOp(parent) && parent.GetFirstChild() == n {
				isSet = true
				type_ = GlobalNamespace_NameProp_OTHER_OBJECT
			}
		}
		if g.forProcessDefines && !isSet && (!g.processDefinesCheck || (t != nil && !t.currentHoistScopeRoot.IsNil() && b.curModuleRoot.IsNil())) {
			return
		}
		if n.GetToken() == rhino.Token_NAME {
			topID := n.GetStringID()
			if topID <= 1 || pTok == rhino.Token_PARAM_LIST {
				return
			}
			if t != nil && b.curModuleRoot.IsNil() && !t.currentHoistScopeRoot.IsNil() {
				switch pTok {
				case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_FUNCTION, rhino.Token_CATCH:
					return
				case rhino.Token_CLASS:
					if parent.GetFirstChild() == n {
						return
					}
				}
			}
			topName := strpool.String(strpool.ID(topID))
			if strings.IndexByte(topName, '.') < 0 {
				root = b.getRootNodeOfTopNameByID(topID, topName, scope, t)
				if !b.IsTopLevelScopeRoot(root) || b.skipsLeafID(topID) {
					return
				}
				var nameMetadata *ModuleMetadataMap_ModuleMetadata
				if root != g.globalRoot {
					nameMetadata = b.curMetadata
				}
				resolvedName = b.getOrCreateRootByID(topID, nameMetadata)
			} else {
				nameStr = topName
				root = b.getRootNodeOfTopNameWithTraversal(GlobalNamespace_GetTopVarName(nameStr), scope, t)
				if !b.IsTopLevelScopeRoot(root) {
					return
				}
				hasName = true
			}
		} else {
			var segBuf [32]uint32
			var segExtra []uint32
			segLen := 0
			protoCut := 0
			hasDotInSeg := false
			var protoTarget rhino.Node
			left := n
			for left.IsGetProp() {
				segID := left.GetStringID()
				if segID == gnPrototypeStrID {
					protoTarget = left.GetFirstChild()
					protoCut = segLen + 1
				} else if !hasDotInSeg && strings.IndexByte(strpool.String(strpool.ID(segID)), '.') >= 0 {
					hasDotInSeg = true
				}
				if segLen < len(segBuf) {
					segBuf[segLen] = segID
				} else {
					if segExtra == nil {
						segExtra = make([]uint32, 0, 64)
						segExtra = append(segExtra, segBuf[:]...)
					}
					segExtra = append(segExtra, segID)
				}
				segLen++
				left = left.GetFirstChild()
			}
			if !left.IsName() {
				return
			}
			topID := left.GetStringID()
			if topID <= 1 {
				return
			}
			topName := strpool.String(strpool.ID(topID))
			if !hasDotInSeg && strings.IndexByte(topName, '.') >= 0 {
				hasDotInSeg = true
			}
			if !hasDotInSeg {
				root = b.getRootNodeOfTopNameByID(topID, topName, scope, t)
				if !b.IsTopLevelScopeRoot(root) {
					return
				}
				segs := segBuf[:segLen]
				if segExtra != nil {
					segs = segExtra
				}
				var nameMetadata *ModuleMetadataMap_ModuleMetadata
				if root != g.globalRoot {
					nameMetadata = b.curMetadata
				}
				if protoTarget.IsNil() {
					if b.skipsLeafID(segs[0]) {
						return
					}
					resolvedName = b.resolveChainByID(topID, segs, nameMetadata)
				} else {
					leafID := topID
					if protoCut < len(segs) {
						leafID = segs[protoCut]
					}
					if b.skipsLeafID(leafID) {
						return
					}
					if !g.forProcessDefines || scope != nil || t == nil {
						if scope == nil {
							scope = t.GetScope()
						}
						if chunk == nil && t != nil {
							chunk = t.GetChunk()
						}
					}
					b.AddOrConfirmRef(b.resolveChainByID(topID, segs[protoCut:], nameMetadata), protoTarget, GlobalNamespace_Ref_Type_PROTOTYPE_GET, scope, chunk)
					return
				}
			} else {
				root = b.getRootNodeOfTopNameWithTraversal(GlobalNamespace_GetTopVarName(topName), scope, t)
				if !b.IsTopLevelScopeRoot(root) {
					return
				}
				if protoTarget.IsNil() {
					var ok bool
					nameStr, ok = n.GetQualifiedNameString()
					if !ok {
						return
					}
					hasName = true
				} else {
					prefix, ok := protoTarget.GetQualifiedNameString()
					if !ok || b.skipsName(prefix) {
						return
					}
					var nameMetadata *ModuleMetadataMap_ModuleMetadata
					if root != g.globalRoot {
						nameMetadata = b.curMetadata
					}
					if !g.forProcessDefines || scope != nil || t == nil {
						if scope == nil {
							scope = t.GetScope()
						}
						if chunk == nil && t != nil {
							chunk = t.GetChunk()
						}
					}
					b.AddOrConfirmRef(b.GetOrCreateName(prefix, nameMetadata), protoTarget, GlobalNamespace_Ref_Type_PROTOTYPE_GET, scope, chunk)
					return
				}
			}
		}
	case rhino.Token_CALL:
		if parent.IsExprResult() &&
			GlobalNamespace_GOOG_PROVIDE.Matches(n.GetFirstChild()) &&
			n.GetSecondChild().IsStringLit() {
			// goog.provide goes through a different code path than regular sets because it can
			// create multiple names, e.g. `goog.provide('a.b.c');` creates the global names
			// a, a.b, and a.b.c. Other sets only create a single global name.
			b.CreateNamesFromProvide(n.GetSecondChild().GetString())
			return
		}
		return
	default:
		return
	}

	if resolvedName != nil {
		if g.forProcessDefines && scope == nil && t != nil {
			if isSet {
				hoistIdx, hoistRoot := t.closestScopeHoistRoot()
				if hoistIdx == 0 || (root != g.globalRoot && hoistRoot == b.curModuleRoot) {
					b.handleResolvedSetFromGlobal(nil, nil, n, resolvedName, type_)
				} else {
					b.handleResolvedSetFromLocal(nil, nil, n, resolvedName)
				}
			} else {
				b.handleResolvedGet(nil, nil, n, resolvedName)
			}
			return
		}
		if scope == nil {
			scope = t.GetScope()
		}
		if chunk == nil && t != nil {
			chunk = t.GetChunk()
		}
		if isSet {
			hoistScope := scope.GetClosestHoistScope()
			if hoistScope.IsGlobal() ||
				(root != g.globalRoot && hoistScope.GetRootNode() == b.curModuleRoot) {
				b.handleResolvedSetFromGlobal(chunk, scope, n, resolvedName, type_)
			} else {
				b.handleResolvedSetFromLocal(chunk, scope, n, resolvedName)
			}
		} else {
			b.handleResolvedGet(chunk, scope, n, resolvedName)
		}
		return
	}

	var nameMetadata *ModuleMetadataMap_ModuleMetadata
	if root != g.globalRoot {
		nameMetadata = b.curMetadata
	}
	if g.forProcessDefines && scope == nil && t != nil {
		if isSet {
			hoistIdx, hoistRoot := t.closestScopeHoistRoot()
			if hoistIdx == 0 || (root != g.globalRoot && hoistRoot == b.curModuleRoot) {
				b.HandleSetFromGlobal(nil, nil, n, nameStr, type_, nameMetadata)
			} else {
				b.HandleSetFromLocal(nil, nil, n, nameStr, nameMetadata)
			}
		} else {
			b.HandleGet(nil, nil, n, nameStr, nameMetadata)
		}
		return
	}
	if scope == nil {
		scope = t.GetScope()
	}
	if chunk == nil && t != nil {
		chunk = t.GetChunk()
	}
	if isSet {
		// Use the closest hoist scope to select handleSetFromGlobal or handleSetFromLocal
		// because they use the term 'global' in an ES5, pre-block-scoping sense.
		hoistScope := scope.GetClosestHoistScope()
		// Consider a set to be 'global' if it is in the hoist scope in which the name is defined.
		// For example, a global name set in a module scope is a 'local' set, but a module-level
		// name set in a module scope is a 'global' set.
		if hoistScope.IsGlobal() ||
			(root != g.globalRoot && hoistScope.GetRootNode() == b.curModuleRoot) {
			b.HandleSetFromGlobal(chunk, scope, n, nameStr, type_, nameMetadata)
		} else {
			b.HandleSetFromLocal(chunk, scope, n, nameStr, nameMetadata)
		}
	} else {
		b.HandleGet(chunk, scope, n, nameStr, nameMetadata)
	}
}

func (b *GlobalNamespace_BuildGlobalNamespace) analyzeObjLitStringKey(stringKeyNode rhino.Node) (string, GlobalNamespace_NameProp, bool) {
	if nameString, ok := nodeUtil_getBestLValueNameString(stringKeyNode); ok {
		// `parent.qname = { myPropName: myValue }`;
		// `NodeUtil.getBestLValueName()` finds the name being assigned for this case.
		return nameString, b.GetValueType(stringKeyNode.GetOnlyChild()), true
	}
	// maybe we have a case like
	// `Object.defineProperties(parentName, { myPropName: { get: ..., set: ..., ... })`
	objLitNode := stringKeyNode.GetParent()
	jutil.CheckArgumentMsg1(objLitNode.IsObjectLit(), "%s", objLitNode)
	objLitParentNode := objLitNode.GetParent()
	if NodeUtil_IsObjectDefinePropertiesDefinition(objLitParentNode) {
		receiverNode := objLitParentNode.GetSecondChild()
		if qname, ok := receiverNode.GetQualifiedNameString(); ok {
			jutil.CheckStateMsg1(objLitNode == receiverNode.GetNext(), "%s", objLitParentNode)
			defined := qname + "." + stringKeyNode.GetString()
			return defined, GlobalNamespace_NameProp_GET_SET, true
		}
	}
	return "", GlobalNamespace_NameProp_OTHER_OBJECT, false
}

// Java: GlobalNamespace.BuildGlobalNamespace.createObjLitStringKeyAnalysis(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) CreateObjLitStringKeyAnalysis(stringKeyNode rhino.Node) *GlobalNamespace_ObjLitStringKeyAnalysis {
	nameString, nameType, ok := b.analyzeObjLitStringKey(stringKeyNode)
	if !ok {
		return GlobalNamespace_ObjLitStringKeyAnalysis_ForNonReference()
	}
	return NewGlobalNamespace_ObjLitStringKeyAnalysis(&nameString, nameType)
}

// CreateNamesFromProvide declares all subnamespaces from
// `goog.provide('some.long.namespace')` globally.
//
// Java: GlobalNamespace.BuildGlobalNamespace.createNamesFromProvide(String)
func (b *GlobalNamespace_BuildGlobalNamespace) CreateNamesFromProvide(namespace string) {
	b.outer.hasProvidedNames = true
	var name *GlobalNamespace_Name
	dot := 0

	for dot >= 0 {
		if dot+1 < len(namespace) {
			if idx := strings.IndexByte(namespace[dot+1:], '.'); idx >= 0 {
				dot = dot + 1 + idx
			} else {
				dot = -1
			}
		} else {
			dot = -1
		}
		var subNamespace string
		if dot < 0 {
			subNamespace = namespace
		} else {
			subNamespace = namespace[:dot]
		}
		jutil.CheckState(subNamespace != "")
		// The top-level name resolves the roots of qualified names
		// (getRootNodeOfTopName), so it is kept even for ProcessDefines.
		if name != nil && b.skipsName(subNamespace) {
			continue
		}
		name = b.GetOrCreateName(subNamespace, nil)
		name.SetBooleanProperty(GlobalNamespace_NameProp_IS_PROVIDED)
	}

	if b.skipsName(namespace) {
		return
	}
	newName := b.GetOrCreateName(namespace, nil)
	newName.SetBooleanProperty(GlobalNamespace_NameProp_IS_PROVIDED)
	clear(b.objLitCache)
}

// IsTopLevelScopeRoot reports whether the given name root represents a
// global or module-level name.
//
// This method will return false for functions and blocks and true for
// module bodies and the globalRoot. The one exception is if the function or
// block is from a goog.loadModule argument, as those functions/blocks are
// treated as module roots.
//
// Java: GlobalNamespace.BuildGlobalNamespace.isTopLevelScopeRoot(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) IsTopLevelScopeRoot(root rhino.Node) bool {
	if root.IsNil() {
		return false
	} else if root == b.outer.globalRoot {
		return true
	} else if root == b.curModuleRoot {
		return true
	}
	// Given
	//   goog.loadModule(function(exports) {
	// pretend that assignments to `exports` or `exports.x = ...` are scoped to the function body,
	// although `exports` is really in the enclosing function parameter scope.
	return !b.curModuleRoot.IsNil() && b.curModuleRoot.IsBlock() && root == b.curModuleRoot.GetParent()
}

func (b *GlobalNamespace_BuildGlobalNamespace) getNameForObjectPatternKey(stringKey rhino.Node) (string, bool) {
	parent := stringKey.GetParent()
	jutil.CheckState(parent.IsObjectPattern())

	patternParent := parent.GetParent()
	if patternParent.IsAssign() || patternParent.IsDestructuringLhs() {
		// this is a top-level string key. we find the name.
		rhs := patternParent.GetSecondChild()
		if rhs.IsNil() {
			// The rhs is null for patterns in parameter lists, enhanced for loops, and catch exprs
			return "", false
		}
		if qname, ok := rhs.GetQualifiedNameString(); ok {
			return qname + "." + stringKey.GetString(), true
		}
		return "", false
	}
	// skip this step for nested patterns for now
	return "", false
}

// GetNameForObjectPatternKey gets the fully qualified name corresponding to
// an object pattern key, as long as it is not in a nested pattern and is
// destructuring an qualified name.
//
// stringKey is a child of an OBJECT_PATTERN node. Returns the global name,
// or null if n doesn't correspond to the key of an object literal that can
// be named.
//
// Java: GlobalNamespace.BuildGlobalNamespace.getNameForObjectPatternKey(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) GetNameForObjectPatternKey(stringKey rhino.Node) *string {
	if s, ok := b.getNameForObjectPatternKey(stringKey); ok {
		return &s
	}
	return nil
}

// GetValueType gets the type of a value or simple expression.
//
// n is an r-value in an assignment or variable declaration (not null).
//
// Java: GlobalNamespace.BuildGlobalNamespace.getValueType(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) GetValueType(n rhino.Node) GlobalNamespace_NameProp {
	switch n.GetToken() {
	case rhino.Token_CLASS:
		return GlobalNamespace_NameProp_CLASS
	case rhino.Token_OBJECTLIT:
		return GlobalNamespace_NameProp_OBJECTLIT
	case rhino.Token_FUNCTION:
		return GlobalNamespace_NameProp_FUNCTION
	case rhino.Token_OR:
		// Recurse on the second value. If the first value were an object
		// literal or function, then the OR would be meaningless and the
		// second value would be dead code. Assume that if the second value
		// is an object literal or function, then the first value will also
		// evaluate to one when it doesn't evaluate to false.
		return b.GetValueType(n.GetLastChild())
	case rhino.Token_HOOK:
		// The same line of reasoning used for the OR case applies here.
		second := n.GetSecondChild()
		t := b.GetValueType(second)
		if t != GlobalNamespace_NameProp_OTHER_OBJECT {
			return t
		}
		third := second.GetNext()
		return b.GetValueType(third)
	default:
	}
	return GlobalNamespace_NameProp_OTHER_OBJECT
}

// HandleSetFromGlobal updates our representation of the global namespace to
// reflect an assignment to a global name in any scope where variables are
// hoisted to the global scope (i.e. the global scope in an ES5 sense).
//
// scope is the current scope; n the node currently being visited; name the
// global name (e.g. "a" or "a.b.c.d"); type the type of the value that the
// name is being assigned.
//
// Java: GlobalNamespace.BuildGlobalNamespace.handleSetFromGlobal(JSChunk, Scope, Node, String, NameProp, ModuleMetadata)
func (b *GlobalNamespace_BuildGlobalNamespace) HandleSetFromGlobal(
	chunk *JSChunk, scope *Scope, n rhino.Node, name string, type_ GlobalNamespace_NameProp, metadata *ModuleMetadataMap_ModuleMetadata) {
	if b.MaybeHandlePrototypePrefix(chunk, scope, n, name, metadata) || b.skipsName(name) {
		return
	}
	b.handleResolvedSetFromGlobal(chunk, scope, n, b.GetOrCreateName(name, metadata), type_)
}

func (b *GlobalNamespace_BuildGlobalNamespace) handleResolvedSetFromGlobal(
	chunk *JSChunk, scope *Scope, n rhino.Node, nameObj *GlobalNamespace_Name, type_ GlobalNamespace_NameProp) {
	if !nameObj.IsGetOrSetDefinition() {
		// Don't change the type of a getter or setter. This is because given:
		//   var a = {set b(item) {}}; a.b = class {};
		// `a.b = class {};` does not change the runtime value of a.b, and we do not want to change
		// the 'type' of a.b to Type.CLASS.
		// TODO(lharker): for non-setter cases, do we really want to just treat the last set of
		// a name as canonical? e.g. what if a name is first set to a class, then an object literal?
		nameObj.SetNameType(type_)
	}
	if n.GetBooleanProp(rhino.Node_MODULE_EXPORT) {
		nameObj.SetBooleanProperty(GlobalNamespace_NameProp_IS_MODULE_PROP)
	}

	if b.IsNestedAssign(n.GetParent()) {
		// This assignment is both a set and a get that creates an alias.
		refType := GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL
		b.AddOrConfirmRef(nameObj, n, refType, scope, chunk)
	} else {
		b.AddOrConfirmRef(nameObj, n, GlobalNamespace_Ref_Type_SET_FROM_GLOBAL, scope, chunk)
		nameObj.SetDeclaredTypeKind(b.GetDeclaredTypeKind(n))
	}
}

// GetDeclaredTypeKind determines whether a set operation is a constructor or
// enumeration or interface declaration. The set operation may either be an
// assignment to a name, a variable declaration, or an object literal key
// mapping.
//
// n is the node that represents the name being set.
//
// Java: GlobalNamespace.BuildGlobalNamespace.getDeclaredTypeKind(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) GetDeclaredTypeKind(n rhino.Node) GlobalNamespace_NameProp {
	valueNode := NodeUtil_GetRValueOfLValue(n)
	var kind GlobalNamespace_NameProp
	if valueNode.IsNil() {
		kind = GlobalNamespace_NameProp_NOT_A_TYPE
	} else if valueNode.IsClass() {
		// Always treat classes as having a declared type. (Transpiled classes are annotated
		// @constructor)
		kind = GlobalNamespace_NameProp_CONSTRUCTOR_TYPE
	} else {
		info := NodeUtil_GetBestJSDocInfo(n)
		// Heed the annotations only if they're sensibly used.
		if info == nil {
			kind = GlobalNamespace_NameProp_NOT_A_TYPE
		} else if info.IsConstructor() && valueNode.IsFunction() {
			kind = GlobalNamespace_NameProp_CONSTRUCTOR_TYPE
		} else if info.IsInterface() && valueNode.IsFunction() {
			kind = GlobalNamespace_NameProp_INTERFACE_TYPE
		} else if info.HasEnumParameterType() && valueNode.IsObjectLit() {
			kind = GlobalNamespace_NameProp_ENUM_TYPE
		} else {
			kind = GlobalNamespace_NameProp_NOT_A_TYPE
		}
	}
	return kind
}

// HandleSetFromLocal updates our representation of the global namespace to
// reflect an assignment to a global name in a local scope.
//
// scope is the current scope; n the node currently being visited; name the
// global name (e.g. "a" or "a.b.c.d").
//
// Java: GlobalNamespace.BuildGlobalNamespace.handleSetFromLocal(JSChunk, Scope, Node, String, ModuleMetadata)
func (b *GlobalNamespace_BuildGlobalNamespace) HandleSetFromLocal(
	chunk *JSChunk, scope *Scope, n rhino.Node, name string, metadata *ModuleMetadataMap_ModuleMetadata) {
	if b.MaybeHandlePrototypePrefix(chunk, scope, n, name, metadata) || b.skipsName(name) {
		return
	}
	b.handleResolvedSetFromLocal(chunk, scope, n, b.GetOrCreateName(name, metadata))
}

func (b *GlobalNamespace_BuildGlobalNamespace) handleResolvedSetFromLocal(
	chunk *JSChunk, scope *Scope, n rhino.Node, nameObj *GlobalNamespace_Name) {
	if n.GetBooleanProp(rhino.Node_MODULE_EXPORT) {
		nameObj.SetBooleanProperty(GlobalNamespace_NameProp_IS_MODULE_PROP)
	}

	if b.IsNestedAssign(n.GetParent()) {
		// This assignment is both a set and a get that creates an alias.
		b.AddOrConfirmRef(nameObj, n, GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL, scope, chunk)
	} else {
		b.AddOrConfirmRef(nameObj, n, GlobalNamespace_Ref_Type_SET_FROM_LOCAL, scope, chunk)
	}
}

// HandleGet updates our representation of the global namespace to reflect a
// read of a global name.
//
// scope is the current scope; n the node currently being visited; name the
// global name (e.g. "a" or "a.b.c.d").
//
// Java: GlobalNamespace.BuildGlobalNamespace.handleGet(JSChunk, Scope, Node, String, ModuleMetadata)
func (b *GlobalNamespace_BuildGlobalNamespace) HandleGet(chunk *JSChunk, scope *Scope, n rhino.Node, name string, metadata *ModuleMetadataMap_ModuleMetadata) {
	if b.MaybeHandlePrototypePrefix(chunk, scope, n, name, metadata) || b.skipsName(name) {
		return
	}
	type_ := b.DetermineRefTypeForGet(n, n, name)

	b.AddOrConfirmRef(b.GetOrCreateName(name, metadata), n, type_, scope, chunk)
}

func (b *GlobalNamespace_BuildGlobalNamespace) handleResolvedGet(chunk *JSChunk, scope *Scope, n rhino.Node, nameObj *GlobalNamespace_Name) {
	type_ := b.determineRefTypeForGetNode(n, n, nameObj)
	b.AddOrConfirmRef(nameObj, n, type_, scope, chunk)
}

// DetermineRefTypeForGet determines the Ref.Type for referenceNode by
// inspecting parents and recusively ascending as necessary, where n
// represents.
//
// n is the node currently being visited; referenceNode the node hosting the
// name; name the global name (e.g. "a" or "a.b.c.d").
//
// Java: GlobalNamespace.BuildGlobalNamespace.determineRefTypeForGet(Node, Node, String)
func (b *GlobalNamespace_BuildGlobalNamespace) DetermineRefTypeForGet(n rhino.Node, referenceNode rhino.Node, name string) GlobalNamespace_Ref_Type {
	return b.determineRefTypeForGetImpl(n, referenceNode, name, nil)
}

func (b *GlobalNamespace_BuildGlobalNamespace) determineRefTypeForGetNode(n rhino.Node, referenceNode rhino.Node, nameObj *GlobalNamespace_Name) GlobalNamespace_Ref_Type {
	return b.determineRefTypeForGetImpl(n, referenceNode, "", nameObj)
}

func (b *GlobalNamespace_BuildGlobalNamespace) determineRefTypeForGetImpl(
	n rhino.Node, referenceNode rhino.Node, name string, nameObj *GlobalNamespace_Name) GlobalNamespace_Ref_Type {
	var type_ GlobalNamespace_Ref_Type
	parent := n.GetParent()
	switch parent.GetToken() {
	case rhino.Token_EXPR_RESULT,
		rhino.Token_IF,
		rhino.Token_WHILE,
		rhino.Token_FOR,
		rhino.Token_INSTANCEOF,
		rhino.Token_TYPEOF,
		rhino.Token_VOID,
		rhino.Token_NOT,
		rhino.Token_BITNOT,
		rhino.Token_POS,
		rhino.Token_NEG,
		rhino.Token_SHEQ,
		rhino.Token_EQ,
		rhino.Token_SHNE,
		rhino.Token_NE,
		rhino.Token_LT,
		rhino.Token_LE,
		rhino.Token_GT,
		rhino.Token_GE,
		rhino.Token_ADD,
		rhino.Token_SUB,
		rhino.Token_MUL,
		rhino.Token_DIV,
		rhino.Token_MOD,
		rhino.Token_EXPONENT,
		rhino.Token_BITAND,
		rhino.Token_BITOR,
		rhino.Token_BITXOR,
		rhino.Token_LSH,
		rhino.Token_RSH,
		rhino.Token_URSH:
		type_ = GlobalNamespace_Ref_Type_DIRECT_GET
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		if n == parent.GetFirstChild() {
			// It is a call target
			type_ = GlobalNamespace_Ref_Type_CALL_GET
		} else if b.IsClassDefiningCall(parent) {
			type_ = GlobalNamespace_Ref_Type_DIRECT_GET
		} else {
			type_ = GlobalNamespace_Ref_Type_ALIASING_GET
		}
	case rhino.Token_NEW:
		if n == parent.GetFirstChild() {
			type_ = GlobalNamespace_Ref_Type_DIRECT_GET
		} else {
			type_ = GlobalNamespace_Ref_Type_ALIASING_GET
		}
	case rhino.Token_CAST, rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		// This node is x or y in (x||y), (x&&y), or (x??y). We only know that an
		// alias is not getting created for this name if the result is used
		// in a boolean context or assigned to the same name
		// (e.g. var a = a || {}).
		type_ = b.determineRefTypeForGetImpl(parent, referenceNode, name, nameObj)
	case rhino.Token_NAME:
		// Only LET, CONST, VAR declarations have NAME nodes
		// with children.
		// Of particular interest is "var n = n || {}"
		sameName := false
		if n != referenceNode {
			if nameObj != nil {
				sameName = nameObj.parentID == 0 && nameObj.baseNameID == parent.GetStringID()
			} else {
				sameName = name == parent.GetString()
			}
		}
		if sameName {
			type_ = GlobalNamespace_Ref_Type_DIRECT_GET
		} else {
			type_ = GlobalNamespace_Ref_Type_ALIASING_GET
		}
	case rhino.Token_COMMA, rhino.Token_HOOK:
		if n != parent.GetFirstChild() {
			// This node is y or z in (x?y:z) or (x,y). We only know that an alias is
			// not getting created for this name if the result is assigned to
			// the same name (e.g. var a = a ? a : {}).
			type_ = b.determineRefTypeForGetImpl(parent, referenceNode, name, nameObj)
		} else {
			type_ = GlobalNamespace_Ref_Type_DIRECT_GET
		}
	case rhino.Token_DELPROP:
		type_ = GlobalNamespace_Ref_Type_DELETE_PROP
	case rhino.Token_CLASS:
		// This node is the superclass in an extends clause.
		type_ = GlobalNamespace_Ref_Type_SUBCLASSING_GET
	case rhino.Token_DESTRUCTURING_LHS, rhino.Token_ASSIGN:
		lhs := n.GetPrevious()
		if lhs.IsNil() {
			// TODO(b/127505242): CAST confused the "is this a get or set?"
			// logic and "handleGet" should not have been called.
			type_ = GlobalNamespace_Ref_Type_ALIASING_GET
			break
		}
		for lhs.IsCast() {
			// Case: `/** @type {!Foo} */ (x) = ...`; or multiple casts like `(cast(cast(x)) =`
			lhs = lhs.GetOnlyChild()
		}

		// This is a recursive ascent check if this an assignment
		// to itself.  This handles cases like: "a.b = a.b || {}"
		if n != referenceNode {
			if nameObj != nil {
				if nameObj.MatchesNode(lhs) {
					return GlobalNamespace_Ref_Type_DIRECT_GET
				}
			} else if lhs.MatchesQualifiedName(name) {
				return GlobalNamespace_Ref_Type_DIRECT_GET
			}
		}

		switch lhs.GetToken() {
		case rhino.Token_NAME,
			rhino.Token_GETPROP,
			rhino.Token_GETELEM,
			// The rhs of an assign or a name declaration is escaped if it's assigned to a name
			// directly ...
			rhino.Token_ARRAY_PATTERN,
			rhino.Token_OBJECT_PATTERN:
			// ... or referenced through numeric/object keys.
			type_ = GlobalNamespace_Ref_Type_ALIASING_GET
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected previous sibling of " + n.GetToken().String() + ": " + jutil.ToJavaString(n.GetPrevious())))
		}
	case rhino.Token_OBJECT_PATTERN, // Handle STRING_KEYS in object patterns.
		rhino.Token_ITER_SPREAD,
		rhino.Token_OBJECT_SPREAD,
		rhino.Token_RETURN,
		rhino.Token_THROW:
		// NOTE: There are likely more cases where we should be returning
		// DIRECT_GET.
		type_ = GlobalNamespace_Ref_Type_ALIASING_GET
	default:
		// NOTE: There are likely more cases where we should be returning
		// DIRECT_GET.
		type_ = GlobalNamespace_Ref_Type_ALIASING_GET
	}
	return type_
}

// AddOrConfirmRef: if there is already a Ref for the given name & node,
// confirm it matches what we would create. Otherwise add a new one.
//
// Java: GlobalNamespace.BuildGlobalNamespace.addOrConfirmRef(Name, Node, Ref.Type, Scope, JSChunk)
func (b *GlobalNamespace_BuildGlobalNamespace) AddOrConfirmRef(
	nameObj *GlobalNamespace_Name, node rhino.Node, refType GlobalNamespace_Ref_Type, scope *Scope, chunk *JSChunk) {
	var existingRef *GlobalNamespace_Ref
	if b.confirmExisting {
		existingRef = nameObj.GetRefForNode(node)
	}
	if existingRef == nil {
		refID, ref := nameObj.createNewRefWithID(chunk, scope, node, refType)
		nameObj.putNewRefID(node, refID)
		nameObj.updateStateForAddedRef(refID, ref)
	} else {
		// module and scope are dependent on Node, so not much point in checking them
		existingRefType := existingRef.type_
		jutil.CheckStateMsg2(
			existingRefType == refType,
			"existing ref type: %s expected: %s",
			existingRefType,
			refType)
	}
}

// Java: GlobalNamespace.BuildGlobalNamespace.isClassDefiningCall(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) IsClassDefiningCall(callNode rhino.Node) bool {
	convention := b.outer.compiler.GetCodingConvention()
	// Look for goog.inherits and J2CL mixin calls
	classes := convention.GetClassesDefinedByCall(callNode)
	if classes != nil {
		return true
	}

	// Look for calls to goog.addSingletonGetter calls.
	className := convention.GetSingletonGetterClassName(callNode)
	return className != nil
}

// IsObjectHasOwnPropertyCall detects calls of the form
// a.b.hasOwnProperty(c); that prevent property collapsing on a.b
//
// Java: GlobalNamespace.BuildGlobalNamespace.isObjectHasOwnPropertyCall(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) IsObjectHasOwnPropertyCall(callNode rhino.Node) bool {
	jutil.CheckArgumentMsg1(callNode.IsCall() || callNode.IsOptChainCall(), "%s", callNode)
	if !callNode.HasTwoChildren() {
		return false
	}
	callee := callNode.GetFirstChild()
	if !callee.IsGetProp() && !callee.IsOptChainGetProp() {
		return false
	}
	receiver := callee.GetFirstChild()
	propName := callee.GetString()
	return ("hasOwnProperty" == propName || "propertyIsEnumerable" == propName) &&
		receiver.IsQualifiedName()
}

// MaybeHandlePrototypePrefix updates our representation of the global
// namespace to reflect a read of a global name's longest prefix before the
// "prototype" property if the name includes the "prototype" property. Does
// nothing otherwise.
//
// scope is the current scope; n the node currently being visited; name the
// global name (e.g. "a" or "a.b.c.d"). Returns whether the name was handled.
//
// Java: GlobalNamespace.BuildGlobalNamespace.maybeHandlePrototypePrefix(JSChunk, Scope, Node, String, ModuleMetadata)
func (b *GlobalNamespace_BuildGlobalNamespace) MaybeHandlePrototypePrefix(
	chunk *JSChunk, scope *Scope, n rhino.Node, name string, metadata *ModuleMetadataMap_ModuleMetadata) bool {
	// We use a string-based approach instead of inspecting the parse tree
	// to avoid complexities with object literals, possibly nested, beneath
	// assignments.

	var numLevelsToRemove int
	var prefix string
	if strings.HasSuffix(name, ".prototype") {
		numLevelsToRemove = 1
		prefix = name[:len(name)-10]
	} else {
		i := strings.Index(name, ".prototype.")
		if i == -1 {
			return false
		}
		prefix = name[:i]
		numLevelsToRemove = 2
		for j := i + 11; j < len(name); j++ {
			if name[j] == '.' {
				numLevelsToRemove++
			}
		}
	}

	if NodeUtil_MayBeObjectLitKey(n) {
		// Object literal keys have no prefix that's referenced directly per
		// key, so we're done.
		return true
	}

	if b.skipsName(prefix) {
		return true
	}
	for i := 0; i < numLevelsToRemove; i++ {
		n = n.GetFirstChild()
	}

	b.AddOrConfirmRef(b.GetOrCreateName(prefix, metadata), n, GlobalNamespace_Ref_Type_PROTOTYPE_GET, scope, chunk)
	return true
}

// IsNestedAssign determines whether an assignment is nested (i.e. whether
// its return value is used).
//
// parent is the parent of the current traversal node (not null). Returns
// whether it appears that the return value of the assignment is used.
//
// Java: GlobalNamespace.BuildGlobalNamespace.isNestedAssign(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) IsNestedAssign(parent rhino.Node) bool {
	return parent.IsAssign() && !parent.GetParent().IsExprResult()
}

func (b *GlobalNamespace_BuildGlobalNamespace) getOrCreateRootByID(baseNameID uint32, metadata *ModuleMetadataMap_ModuleMetadata) *GlobalNamespace_Name {
	g := b.outer
	if metadata == nil {
		if id := g.lookupChildID(0, baseNameID); id != 0 {
			return g.nameByID(id)
		}
		node := g.allocName(baseNameID, 0, g.sourceKind)
		node.propertyBitSet |= gnNamePropIndexedGlobal
		g.indexedGlobalCount++
		g.globalNameIDs = append(g.globalNameIDs, node.selfID)
		g.putChildID(0, baseNameID, node.selfID)
		return node
	}
	if g.nameMapByModule != nil {
		if id := g.nameMapByModule[globalNamespaceModuleKey{module: metadata, nameID: baseNameID}]; id != 0 {
			return g.nameByID(id)
		}
	} else {
		g.nameMapByModule = make(map[globalNamespaceModuleKey]uint32)
	}
	node := g.allocName(baseNameID, 0, g.sourceKind)
	g.nameMapByModule[globalNamespaceModuleKey{module: metadata, nameID: baseNameID}] = node.selfID
	return node
}

func (b *GlobalNamespace_BuildGlobalNamespace) getOrCreateChildByID(parent *GlobalNamespace_Name, baseNameID uint32, metadata *ModuleMetadataMap_ModuleMetadata) *GlobalNamespace_Name {
	g := b.outer
	if id := g.lookupChildID(parent.selfID, baseNameID); id != 0 {
		return g.nameByID(id)
	}
	node := parent.addPropertyWithBaseNameID(baseNameID, g.sourceKind)
	if metadata == nil {
		node.propertyBitSet |= gnNamePropIndexedGlobal
		g.indexedGlobalCount++
	}
	g.putChildID(parent.selfID, baseNameID, node.selfID)
	return node
}

func (b *GlobalNamespace_BuildGlobalNamespace) resolveChainByID(topID uint32, revSegs []uint32, metadata *ModuleMetadataMap_ModuleMetadata) *GlobalNamespace_Name {
	cur := b.getOrCreateRootByID(topID, metadata)
	for i := len(revSegs) - 1; i >= 0; i-- {
		cur = b.getOrCreateChildByID(cur, revSegs[i], metadata)
	}
	return cur
}

// GetOrCreateName gets a Name instance for a global name. Creates it if
// necessary, as well as instances for any of its prefixes that are not yet
// defined.
//
// name is a global name (e.g. "a", "a.b.c.d"). Returns the Name instance
// for name.
//
// Java: GlobalNamespace.BuildGlobalNamespace.getOrCreateName(String, ModuleMetadata)
func (b *GlobalNamespace_BuildGlobalNamespace) GetOrCreateName(name string, metadata *ModuleMetadataMap_ModuleMetadata) *GlobalNamespace_Name {
	if b.lastNode != nil && name == b.lastName && metadata == b.lastMeta {
		return b.lastNode
	}
	if b.lastParentNode != nil && name == b.lastParentName && metadata == b.lastParentMeta {
		return b.lastParentNode
	}
	var node *GlobalNamespace_Name
	i := strings.LastIndexByte(name, '.')
	if i >= 0 {
		parentName := name[:i]
		baseSeg := name[i+1:]
		parent := b.GetOrCreateName(parentName, metadata)
		b.lastParentName = parentName
		b.lastParentMeta = metadata
		b.lastParentNode = parent
		node = b.getOrCreateChildByID(parent, uint32(strpool.Intern(baseSeg)), metadata)
	} else {
		node = b.getOrCreateRootByID(uint32(strpool.Intern(name)), metadata)
	}
	b.lastName = name
	b.lastMeta = metadata
	b.lastNode = node
	return node
}

// -------------------------------------------------------------------------

// Java: GlobalNamespace.createNameForTesting(String)
func (g *GlobalNamespace) CreateNameForTesting(name string) *GlobalNamespace_Name {
	return NewGlobalNamespace_Name(g, name, nil, GlobalNamespace_SourceKind_CODE)
}

// GlobalNamespace_Inlinability is how much to inline a Name.
//
// The `Inlinability#INLINE_BUT_KEEP_DECLARATION_*` cass are really an
// indicator that something 'unsafe' is happening in order to not break
// CollapseProperties as badly. Sadly Inlinability #INLINE_COMPLETELY may
// *also* be unsafe.
type GlobalNamespace_Inlinability int32

const (
	GlobalNamespace_Inlinability_INLINE_COMPLETELY GlobalNamespace_Inlinability = iota
	GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION
	GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY
	GlobalNamespace_Inlinability_DO_NOT_INLINE
)

var globalNamespace_Inlinability_names = [...]string{
	"INLINE_COMPLETELY",
	"INLINE_BUT_KEEP_DECLARATION",
	"INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY",
	"DO_NOT_INLINE",
}

// The enum fields, indexed by ordinal.
var (
	globalNamespace_Inlinability_shouldInlineUsages      = [...]bool{true, true, true, false}
	globalNamespace_Inlinability_shouldRemoveDeclaration = [...]bool{true, false, false, false}
	globalNamespace_Inlinability_canCollapse             = [...]bool{true, true, true, false}
)

// Java: GlobalNamespace.Inlinability.name()
func (i GlobalNamespace_Inlinability) String() string { return globalNamespace_Inlinability_names[i] }

// Java: GlobalNamespace.Inlinability.ordinal()
func (i GlobalNamespace_Inlinability) Ordinal() int { return int(i) }

// Java: GlobalNamespace.Inlinability.shouldInlineUsages()
func (i GlobalNamespace_Inlinability) ShouldInlineUsages() bool {
	return globalNamespace_Inlinability_shouldInlineUsages[i]
}

// Java: GlobalNamespace.Inlinability.shouldRemoveDeclaration()
func (i GlobalNamespace_Inlinability) ShouldRemoveDeclaration() bool {
	return globalNamespace_Inlinability_shouldRemoveDeclaration[i]
}

// Java: GlobalNamespace.Inlinability.canCollapse()
func (i GlobalNamespace_Inlinability) CanCollapse() bool {
	return globalNamespace_Inlinability_canCollapse[i]
}

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

// Java: GlobalNamespace.Name.castRefsForNodeMap()
func (n *GlobalNamespace_Name) CastRefsForNodeMap() *jutil.LinkedHashMap[rhino.Node, *GlobalNamespace_Ref] {
	return nil
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

// CalculateInlinability returns whether to treat this alias as completely
// inlineable or to keep the aliasing assignment
//
// This method used to only return true/false, but now returns an enum in
// order to track more information about "unsafely" inlineable names.
//
// CollapseProperties will flatten `@constructor` properties even if they are
// potentially accessed by a reference other than their fully qualified name,
// which breaks those other refs. To avoid breakages AggressiveInlineAliases
// must unsafely inline constructor properties that alias another global
// name. Existing code depends on this behavior, and it's not easily
// determinable where these dependencies are.
//
// However, AggressiveInlineAliases must not also remove the initializtion of
// an alias if it is not safely inlineable. (i.e. if
// Inlinability#shouldRemoveDeclaration()). It's possible that a third name
// aliases the alias - we might later inline the third name (as an alias of
// the original alias) and don't want to set the third name to null.
//
// Java: GlobalNamespace.Name.calculateInlinability()
func (n *GlobalNamespace_Name) CalculateInlinability() GlobalNamespace_Inlinability {
	// Only simple aliases with direct usage are inlineable.
	// Exactly 2 global sets could be OK, if the first is the declaration with no initializer.
	if n.InExterns() || !n.HasOneRealGlobalSet() || n.localSets != 0 {
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	// TODO(lharker): consider separating canCollapseOrInline() into this method, since it
	// duplicates some logic here
	inlinability := n.CanCollapseOrInline()
	if !inlinability.ShouldInlineUsages() {
		// if you can't even inline the usages, do nothing.
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}
	var initialChunk *JSChunk
	if initRef := n.GetInitialization(); initRef != nil {
		initialChunk = initRef.GetChunk()
	} else {
		initialChunk = n.getDeclarationRef().GetChunk()
	}

	// Only allow inlining of simple references.
	for ref := range n.Refs() {
		switch ref.type_ {
		case GlobalNamespace_Ref_Type_SET_FROM_GLOBAL, GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL:
			// Expect one global set
			jutil.CheckStateMsg1(
				ref.IsUninitializedDeclaration() || ref.GetChunk() == initialChunk, "%s", ref)
			continue
		case GlobalNamespace_Ref_Type_SET_FROM_LOCAL, GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL:
			panic(jutil.NewIllegalStateException())
		case GlobalNamespace_Ref_Type_ALIASING_GET,
			GlobalNamespace_Ref_Type_DIRECT_GET,
			GlobalNamespace_Ref_Type_PROTOTYPE_GET,
			GlobalNamespace_Ref_Type_CALL_GET,
			GlobalNamespace_Ref_Type_SUBCLASSING_GET:
			// This name has a reference in a different chunk that is not guaranteed to be loaded
			// before the alias is initialized. In an ideal world, we might just back off inlining
			// this alias entirely. Unfortunately, that causes new runtime errors in practice,
			// because property collapsing is intentionally unsafe and will collapse properties even
			// when it will break code, and so it's safer to still inling usages of this alias
			// except for the specific invalid cross-chunk dependency(s).
			refChunk := ref.GetChunk()
			if inlinability != GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY &&
				refChunk != initialChunk &&
				!n.outer.compiler.GetChunkGraph().DependsOn(refChunk, initialChunk) {
				inlinability = GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY
			}
			continue
		case GlobalNamespace_Ref_Type_DELETE_PROP:
			return GlobalNamespace_Inlinability_DO_NOT_INLINE
		}
		panic(jutil.NewAssertionErrorNoMessage())
	}
	return inlinability
}

// Java: GlobalNamespace.Name.hasOneRealGlobalSet()
func (n *GlobalNamespace_Name) HasOneRealGlobalSet() bool {
	return n.globalSets == 1 || (n.globalSets == 2 && n.getDeclarationRef().IsUninitializedDeclaration())
}

// Java: GlobalNamespace.Name.canCollapse()
func (n *GlobalNamespace_Name) CanCollapse() bool {
	return n.CanCollapseOrInline().CanCollapse()
}

// CanCollapseOrInline determines whether it's safe to collapse properties on
// an objects
//
// For legacy reasons, both CollapseProperties and AggressiveInlineAliases
// share the same logic when deciding whether to inline properties or to
// collapse them.
//
// The main reasons we cannot inline/collapse properties of a name name are:
//
//	a) it is set multiple times or set once in a local scope,
//	b) one or more of the above conditions in canCollapseOrInlineChildNames
//	   applies to the namespace it's on,
//	c) it's annotated at-nocollapse
//	d) it's in the externs,
//	e) or it's a known getter or setter, not a regular property
//	f) it's an ES6 class static method that references `super` or the internal class name
//
// We ignore conditions (a) and (b) on at-constructor and at-enum names in
// CollapseProperties.
//
// In AggressiveInlineAliases we want to do some partial backoff if (a) and
// (b) are false for at-constructor or at-enum names, which is why we return
// an enum value instead of a boolean.
//
// Java: GlobalNamespace.Name.canCollapseOrInline()
func (n *GlobalNamespace_Name) CanCollapseOrInline() GlobalNamespace_Inlinability {
	if n.InExterns() {
		// condition (d)
		n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "declared in externs")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}
	if n.IsGetOrSetDefinition() {
		// condition (e)
		n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "getter / setter")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}
	if n.IsCollapsingExplicitlyDenied() {
		// condition (c)
		n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "@nocollapse")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.ReferencesSuperOrInnerClassName() {
		// condition (f)
		n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "references super or inner class name")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.IsToStringValueOfInObjectLiteral() {
		n.LogDecision(
			GlobalNamespace_Inlinability_DO_NOT_INLINE,
			"references explicit definition of toString/valueOf functions used implicitly in the JS"+
				" language")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.deleteProps > 0 {
		// If we inline or collapse, then the delete operation will be incorrect.
		n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "delete operator is used on this property")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if declRef := n.getDeclarationRef(); declRef != nil {
		declaration := declRef.GetNode()
		if declaration.GetParent().IsObjectLit() {
			if n.outer.DeclarationHasFollowingObjectSpreadSibling(declaration) {
				// Case: `var x = {a: 0, ...b, c: 2}` where declaration is `a` but not `c`.
				// Following spreads may overwrite the declaration.
				n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "obj lit property followed by spread")
				return GlobalNamespace_Inlinability_DO_NOT_INLINE
			}
			// We may be in a deeply nested object literal like, `{a: {b: {c: 1}}}`, so find the
			// outermost object literal node and check its ancestor expressions to determine
			// whether it is used conditionally.
			child := n.outer.GetOutermostObjectLit(declaration)
			for p := child.GetParent(); !p.IsNil() && !p.GetParent().IsNil() && !NodeUtil_IsStatement(p); child, p = p, p.GetParent() {
				if (p.IsOr() && p.GetFirstChild() != child) ||
					(p.IsHook() && p.GetFirstChild() != child) ||
					(p.IsAnd() && p.GetFirstChild() != child) {
					// Case: `var x = y || {a: b}` or `var x = cond ? y : {a: b}` or `var x = cond && {a:
					// b}`.
					n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "conditional definition")
					return GlobalNamespace_Inlinability_DO_NOT_INLINE
				}
			}
		}
	}

	// condition (a)
	isUnchangedThroughFullName :=
		(n.globalSets > 0 || n.localSets > 0) && n.localSetsWithNoCollapse == 0 && n.deleteProps == 0
	// additional information about condition (b)
	parentInlinability := GlobalNamespace_Inlinability_INLINE_COMPLETELY
	if parent := n.GetParent(); parent != nil {
		parentInlinability = parent.CanCollapseOrInlineChildNames()
	}

	// if condition (a) or condition (b) is not true, but this is a declared name, we may need
	// to allow inlining usages of a variable but keep the declaration.
	switch parentInlinability {
	case GlobalNamespace_Inlinability_INLINE_COMPLETELY:
		if isUnchangedThroughFullName {
			n.LogDecision(
				GlobalNamespace_Inlinability_INLINE_COMPLETELY, "parent inlineable: unchanged through full name")
			return GlobalNamespace_Inlinability_INLINE_COMPLETELY
		} else {
			// maybe inline usages of this name, but only if a declared type. non-declared-types
			// just
			// back off and don't inline at all
			unsafeInlinablility := n.GetUnsafeInlinabilityBasedOnDeclaredType()
			n.LogDecision(unsafeInlinablility, "parent inlineable: changed through full name")
			return unsafeInlinablility
		}
	case GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION,
		GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY:
		// this is definitely not safe to completely inline/collapse of its parent
		// if it's a declared type, we should still partially inline it and completely collapse it
		// if not a declared type we should partially inline it iff the other conditions hold
		if n.IsDeclaredType() {
			unsafeInlinability := n.GetUnsafeInlinabilityBasedOnDeclaredType()
			n.LogDecision(unsafeInlinability, "parent unsafely inlineable & is declared type")
			return unsafeInlinability
		} else if isUnchangedThroughFullName {
			n.LogDecision(
				parentInlinability, "parent unsafely inlineable & unchanged through full name")
			return parentInlinability
		} else {
			// Not a declared type. We may still 'partially' inline it because it must be a property
			// on an @enum or @constructor, but only if it actually matches conditions (a) and (b)
			n.LogDecision(
				GlobalNamespace_Inlinability_DO_NOT_INLINE,
				"parent unsafely inlineable & changed through full name")
			return GlobalNamespace_Inlinability_DO_NOT_INLINE
		}
	case GlobalNamespace_Inlinability_DO_NOT_INLINE:
		// If the parent is unsafely to collapse/inline, we will still inline it if it's on
		// a declaredType (i.e. @constructor or @enum), but we propagate the information that
		// the parent is unsafe. If this is not a declared type, return DO_NOT_INLINE.
		unsafeInlinability := n.GetUnsafeInlinabilityBasedOnDeclaredType()
		n.LogDecision(unsafeInlinability, "parent cannot be inlined")
		return unsafeInlinability
	}
	panic(jutil.NewIllegalStateExceptionMsg("unknown enum value " + parentInlinability.String()))
}

// Java: GlobalNamespace.Name.logDecision(Inlinability, String)
func (n *GlobalNamespace_Name) LogDecision(inlinability GlobalNamespace_Inlinability, reason string) {
	if n.outer.decisionsLog != nil && n.outer.decisionsLog.IsLogging() {
		n.outer.decisionsLog.LogFormat("%s: %s: %s", n.GetFullName(), inlinability.String(), reason)
	}
}

// Java: GlobalNamespace.Name.getUnsafeInlinabilityBasedOnDeclaredType()
func (n *GlobalNamespace_Name) GetUnsafeInlinabilityBasedOnDeclaredType() GlobalNamespace_Inlinability {
	if n.GetBooleanProperty(GlobalNamespace_NameProp_CONSTRUCTOR_TYPE) {
		return GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION
	} else if n.GetBooleanProperty(GlobalNamespace_NameProp_INTERFACE_TYPE) {
		return GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION
	} else if n.GetBooleanProperty(GlobalNamespace_NameProp_ENUM_TYPE) {
		return GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION
	} else if n.GetBooleanProperty(GlobalNamespace_NameProp_NOT_A_TYPE) {
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}
	panic(jutil.NewIllegalStateExceptionMsg("name missing declaredType value: " + n.String()))
}

// ReferencesSuperOrInnerClassName examines ES6 class members for some syntax
// that blocks collapsing
//
// Specifically, this looks for super references and references to inner
// class names. These are unique to ES6 static class members so we don't
// need more general handling.
//
// TODO(b/122665204): also return false on `super` in an object lit method
//
// Java: GlobalNamespace.Name.referencesSuperOrInnerClassName()
func (n *GlobalNamespace_Name) ReferencesSuperOrInnerClassName() bool {
	ref := n.getDeclarationRef()
	if ref == nil {
		return false
	}
	member := ref.GetNode()
	if member.IsNil() || !(member.IsStaticMember() && member.GetParent().IsClassMembers()) {
		return false
	}
	// Get either the function body or the member field initializer (which may be null if there is
	// no initializer).
	var body rhino.Node
	if member.IsMemberFieldDef() {
		body = member.GetFirstChild()
	} else {
		body = NodeUtil_GetFunctionBody(member.GetFirstChild())
	}
	if !body.IsNil() && NodeUtil_ReferencesSuper(body) {
		return true
	}

	classNode := member.GetGrandparent()
	if NodeUtil_IsClassDeclaration(classNode) {
		return false // e.g. class C {}
	}

	innerNameNode := classNode.GetFirstChild()
	return !innerNameNode.IsEmpty() && // e.g. const C = class {};
		NodeUtil_IsNameReferenced(member, innerNameNode.GetString())
}

// Java: GlobalNamespace.Name.isSetInLoop()
func (n *GlobalNamespace_Name) IsSetInLoop() bool {
	ref := n.getDeclarationRef()
	if ref != nil {
		node := ref.GetNode()
		if !node.IsNil() {
			return NodeUtil_IsWithinLoop(node)
		}
	}
	return false
}

// Java: GlobalNamespace.Name.isGetOrSetDefinition()
func (n *GlobalNamespace_Name) IsGetOrSetDefinition() bool {
	return n.GetBooleanProperty(GlobalNamespace_NameProp_GET_SET)
}

// Java: GlobalNamespace.Name.canCollapseUnannotatedChildNames()
func (n *GlobalNamespace_Name) CanCollapseUnannotatedChildNames() bool {
	return n.CanCollapseOrInlineChildNames().CanCollapse()
}

// IsToStringValueOfInObjectLiteral: toString/valueOf are implicitly used as
// part of the JS language should not be collapsed.
//
// Java: GlobalNamespace.Name.isToStringValueOfInObjectLiteral()
func (n *GlobalNamespace_Name) IsToStringValueOfInObjectLiteral() bool {
	parent := n.GetParent()
	baseName := n.GetBaseName()
	return n.IsFunction() &&
		parent != nil &&
		parent.IsObjectLiteral() &&
		(baseName == "toString" || baseName == "valueOf")
}

// CanCollapseOrInlineChildNames returns whether to assume that child
// properties of this name are collapsible/inlineable
//
// For legacy reasons, both CollapseProperties and AggressiveInlineAliases
// share the same logic when deciding whether to inline properties or to
// collapse them.
//
// The main reasons we cannot inline/collapse properties of a name name are:
//
//	a) it is set multiple times
//	b) its properties might not be referred to by their full qname but on a different object
//	c) one or more of the above conditions applies to a parent name
//	d) it's annotated @nocollapse
//	e) it's in the externs.
//	f) it's assigned a value that supports being aliased (e.g. ctors, so their static
//	   properties can be accessed via `this`)
//
// However, in some cases for properties of `@constructor` or `@enum` names,
// we ignore some of these conditions in order to more aggressively collapse
// `@constructor`s used in goog.provide namespace chains.
//
// Java: GlobalNamespace.Name.canCollapseOrInlineChildNames()
func (n *GlobalNamespace_Name) CanCollapseOrInlineChildNames() GlobalNamespace_Inlinability {
	// condition (a) and (b)
	if n.GetBooleanProperty(GlobalNamespace_NameProp_OTHER_OBJECT) {
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "NameProp.OTHER_OBJECT")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	} else if n.IsGetOrSetDefinition() {
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "getter/setter")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	} else if n.globalSets != 1 {
		if n.outer.decisionsLog != nil && n.outer.decisionsLog.IsLogging() {
			n.LogChildNamesDecisionSupplier(
				GlobalNamespace_Inlinability_DO_NOT_INLINE, func() string { return "set " + strconv.Itoa(int(n.globalSets)) + " times globally" })
		}
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	} else if n.localSets != 0 {
		if n.outer.decisionsLog != nil && n.outer.decisionsLog.IsLogging() {
			n.LogChildNamesDecisionSupplier(
				GlobalNamespace_Inlinability_DO_NOT_INLINE, func() string { return "set " + strconv.Itoa(int(n.localSets)) + " times locally" })
		}
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	} else if n.deleteProps != 0 {
		if n.outer.decisionsLog != nil && n.outer.decisionsLog.IsLogging() {
			n.LogChildNamesDecisionSupplier(
				GlobalNamespace_Inlinability_DO_NOT_INLINE,
				func() string { return "properties are deleted " + strconv.Itoa(int(n.deleteProps)) + " times" })
		}
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	// Don't try to collapse if the one global set is a twin reference.
	// We could theoretically handle this case in CollapseProperties, but
	// it's probably not worth the effort.
	declRef := jutil.CheckNotNull(n.getDeclarationRef())
	if declRef.IsTwin() {
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "twinned declaration")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.IsCollapsingExplicitlyDenied() {
		// condition (d)
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "@nocollapse")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.IsSetInLoop() {
		// condition (a)
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "set in a loop")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.UsesHasOwnProperty() {
		// condition (b)
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "hasOwnProperty() call exists")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.ValueImplicitlySupportsAliasing() {
		// condition (f)
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "value implicitly supports aliasing")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	// If this is a key of an aliased object literal, then it will be aliased
	// later. So we won't be able to collapse its properties.
	// condition (b)
	parent := n.GetParent()
	if parent != nil && parent.ShouldKeepKeys() {
		unsafeInlinability := n.GetUnsafeInlinabilityBasedOnDeclaredType()
		n.LogChildNamesDecision(unsafeInlinability, "parent.shouldKeepKeys()")
		return unsafeInlinability
	}

	// If this is aliased, then its properties can't be collapsed either. but we may do so anyway
	// if it's a declared type.
	// condition (b)
	if n.aliasingGets > 0 {
		unsafeInlinability := n.GetUnsafeInlinabilityBasedOnDeclaredType()
		if n.outer.decisionsLog != nil && n.outer.decisionsLog.IsLogging() {
			n.LogChildNamesDecisionSupplier(
				unsafeInlinability, func() string { return strconv.Itoa(int(n.aliasingGets)) + " aliasing gets exist" })
		}
		return unsafeInlinability
	}

	if parent == nil {
		// this is completely safe to inline! yay
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_INLINE_COMPLETELY, "no reason not to inline")
		return GlobalNamespace_Inlinability_INLINE_COMPLETELY
	}

	// Cases are:
	//  - parent is safe to completely inline. then same for this name
	//  - parent is unsafe but should still be inlined. then same for this name
	//  - parent is unsafe, should not be inlined at all. then return either DO_NOT_INLINE,
	//    or maybe unsafely inline if this is a ctor property
	parentInlinability := parent.CanCollapseOrInlineChildNames()
	if parentInlinability == GlobalNamespace_Inlinability_DO_NOT_INLINE {
		// the parent name is used in a way making this unsafe to inline, but we might want to
		// inline usages of this name
		unsafeInlinability := n.GetUnsafeInlinabilityBasedOnDeclaredType()
		n.LogChildNamesDecision(unsafeInlinability, "parent is not inlineable")
		return unsafeInlinability
	}
	n.LogChildNamesDecision(parentInlinability, "inherited from parent")
	return parentInlinability
}

// Java: GlobalNamespace.Name.logChildNamesDecision(Inlinability, String)
func (n *GlobalNamespace_Name) LogChildNamesDecision(inlinability GlobalNamespace_Inlinability, reason string) {
	if n.outer.decisionsLog != nil && n.outer.decisionsLog.IsLogging() {
		n.outer.decisionsLog.LogFormat("%s: children: %s: %s", n.GetFullName(), inlinability.String(), reason)
	}
}

// Java: GlobalNamespace.Name.logChildNamesDecision(Inlinability, Supplier<String>)
func (n *GlobalNamespace_Name) LogChildNamesDecisionSupplier(inlinability GlobalNamespace_Inlinability, reasonSupplier func() string) {
	if n.outer.decisionsLog != nil && n.outer.decisionsLog.IsLogging() {
		n.outer.decisionsLog.LogSupplier(
			func() string {
				return n.GetFullName() + ": children: " + inlinability.String() + ": " + reasonSupplier()
			})
	}
}

// Java: GlobalNamespace.Name.valueImplicitlySupportsAliasing()
func (n *GlobalNamespace_Name) ValueImplicitlySupportsAliasing() bool {
	if !n.outer.enableImplicitlyAliasedValues {
		return false
	}

	if n.GetBooleanProperty(GlobalNamespace_NameProp_CLASS) {
		// Properties on classes may be referenced via `this` in static methods.
		return true
	}
	if n.GetBooleanProperty(GlobalNamespace_NameProp_FUNCTION) {
		// We want ES5 ctors/interfaces to behave consistently with ES6 because:
		// - transpilation should not change behaviour
		// - updating code shouldn't be hindered by behaviour changes
		jsdoc := n.GetJSDocInfo()
		return jsdoc != nil && jsdoc.IsConstructorOrInterface()
	}
	return false
}

// ShouldKeepKeys reports whether this is an object literal that needs to
// keep its keys.
//
// Java: GlobalNamespace.Name.shouldKeepKeys()
func (n *GlobalNamespace_Name) ShouldKeepKeys() bool {
	return n.IsObjectLiteral() && (n.aliasingGets > 0 || n.IsCollapsingExplicitlyDenied())
}

// Java: GlobalNamespace.Name.needsToBeStubbed()
func (n *GlobalNamespace_Name) NeedsToBeStubbed() bool {
	return n.globalSets == 0 &&
		n.localSets > 0 &&
		n.localSetsWithNoCollapse == 0 &&
		!n.IsCollapsingExplicitlyDenied()
}

// Java: GlobalNamespace.Name.setDeclaredTypeKind(NameProp)
func (n *GlobalNamespace_Name) SetDeclaredTypeKind(declaredType GlobalNamespace_NameProp) {
	jutil.CheckArgumentMsg1(
		(declaredType.Bit()&GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK) != 0,
		"Unexpected NameProp for declaredType %s",
		declaredType)
	n.propertyBitSet =
		(n.propertyBitSet &^ GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK) | declaredType.Bit()
	if declaredType != GlobalNamespace_NameProp_NOT_A_TYPE {
		g := n.outer
		for ancestorID := n.parentID; ancestorID != 0; {
			ancestor := g.nameByID(ancestorID)
			ancestor.SetBooleanProperty(GlobalNamespace_NameProp_IS_DECLARED)
			ancestorID = ancestor.parentID
		}
	}
}

// Java: GlobalNamespace.Name.setNameType(NameProp)
func (n *GlobalNamespace_Name) SetNameType(type_ GlobalNamespace_NameProp) {
	jutil.CheckArgumentMsg1(
		(type_.Bit()&GlobalNamespace_NameProp_NAME_KIND_MASK) != 0, "Unexpected NameProp for nameType %s", type_)
	n.propertyBitSet = (n.propertyBitSet &^ GlobalNamespace_NameProp_NAME_KIND_MASK) | type_.Bit()
}

// Java: GlobalNamespace.Name.isDeclaredType()
func (n *GlobalNamespace_Name) IsDeclaredType() bool {
	jutil.CheckStateMsg1(
		(n.propertyBitSet&GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK) != 0, "%s", n.propertyBitSet)
	return !n.GetBooleanProperty(GlobalNamespace_NameProp_NOT_A_TYPE)
}

// Java: GlobalNamespace.Name.isConstructor()
func (n *GlobalNamespace_Name) IsConstructor() bool {
	declNode := n.getDeclarationRef().node
	rvalueNode := NodeUtil_GetRValueOfLValue(declNode)
	jsdoc := NodeUtil_GetBestJSDocInfo(declNode)
	return !rvalueNode.IsNil() &&
		rvalueNode.IsFunction() &&
		jsdoc != nil &&
		jsdoc.IsConstructor()
}

// IsNamespaceObjectLit determines whether this name is a prefix of at least
// one class or enum name. Because classes and enums are always collapsed,
// the namespace will have different properties in compiled code than in
// uncompiled code.
//
// For example, if foo.bar.DomHelper is a class, then foo and foo.bar are
// considered namespaces.
//
// Java: GlobalNamespace.Name.isNamespaceObjectLit()
func (n *GlobalNamespace_Name) IsNamespaceObjectLit() bool {
	return n.GetBooleanProperty(GlobalNamespace_NameProp_IS_DECLARED) && n.IsObjectLiteral()
}

// IsSimpleName determines whether this is a simple name (as opposed to a
// qualified name).
//
// Java: GlobalNamespace.Name.isSimpleName()
func (n *GlobalNamespace_Name) IsSimpleName() bool {
	return n.parentID == 0
}

// Java: GlobalNamespace.Name.getTypeDebugString()
func (n *GlobalNamespace_Name) GetTypeDebugString() string {
	if n.IsClass() {
		return "CLASS"
	} else if n.IsFunction() {
		return "FUNCTION"
	} else if n.IsObjectLiteral() {
		return "OBJECTLIT"
	} else if n.IsGetOrSetDefinition() {
		return "GET_SET"
	} else if n.GetBooleanProperty(GlobalNamespace_NameProp_OTHER_OBJECT) {
		return "OTHER"
	}
	panic(jutil.NewAssertionError("Missing NameProp for name kind in " + strconv.Itoa(int(n.propertyBitSet))))
}

// Java: GlobalNamespace.Name.toString()
func (n *GlobalNamespace_Name) String() string {
	return n.GetFullName() +
		" (" +
		n.GetTypeDebugString() +
		"): " +
		strings.Join([]string{
			"globalSets=" + strconv.Itoa(int(n.globalSets)),
			"localSets=" + strconv.Itoa(int(n.localSets)),
			"totalGets=" + strconv.Itoa(int(n.totalGets)),
			"aliasingGets=" + strconv.Itoa(int(n.aliasingGets)),
			"callGets=" + strconv.Itoa(int(n.callGets)),
			"subclassingGets=" + strconv.Itoa(int(n.subclassingGets)),
		}, ", ")
}

// Java: GlobalNamespace.Name.getJSDocInfo()
func (n *GlobalNamespace_Name) GetJSDocInfo() *rhino.JSDocInfo {
	// e.g.
	// /** @type {string} */ X.numProp;     // could be a declaration, but...
	// /** @type {number} */ X.numProp = 3; // assignment wins
	if n.firstDeclDocID != 0 {
		return n.outer.docByID(n.firstDeclDocID)
	}
	return n.outer.docByID(n.firstQnameDocID)
}

// GetDocInfoForDeclaration tries to get the doc info for a given
// declaration ref.
//
// Java: GlobalNamespace.Name.getDocInfoForDeclaration(Ref)
func (n *GlobalNamespace_Name) GetDocInfoForDeclaration(ref *GlobalNamespace_Ref) *rhino.JSDocInfo {
	if !ref.node.IsNil() {
		refParent := ref.node.GetParent()
		if refParent.IsNil() {
			// May happen when inlineAliases removes refs from the AST.
			return nil
		}
		switch refParent.GetToken() {
		case rhino.Token_FUNCTION, rhino.Token_ASSIGN, rhino.Token_CLASS:
			return refParent.GetJSDocInfo()
		case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
			if ref.node == refParent.GetFirstChild() {
				return refParent.GetJSDocInfo()
			}
			return ref.node.GetJSDocInfo()
		case rhino.Token_OBJECTLIT, rhino.Token_CLASS_MEMBERS:
			return ref.node.GetJSDocInfo()
		default:
		}
	}

	return nil
}

// Java: GlobalNamespace.Name.isModuleExport()
func (n *GlobalNamespace_Name) IsModuleExport() bool {
	return n.GetBooleanProperty(GlobalNamespace_NameProp_IS_MODULE_PROP)
}

// GetOutermostObjectLit, given something like the `'c'` STRING_KEY node in
// `x = {a: {b: {c: 0}}};`, returns the Node for the outermost object
// literal.
//
// objLitKey must be the child of an OBJECT_LIT node. Returns the first
// ancestor that is an OBJECT_LIT and whose grandparent is not an
// OBJECT_LIT.
//
// Java: GlobalNamespace.getOutermostObjectLit(Node)
func (g *GlobalNamespace) GetOutermostObjectLit(objLitKey rhino.Node) rhino.Node {
	outermostObjectLit := objLitKey.GetParent()
	jutil.CheckStateMsg1(outermostObjectLit.IsObjectLit(), "%s", outermostObjectLit)
	for {
		objLitGrandparent := outermostObjectLit.GetGrandparent()
		if !objLitGrandparent.IsNil() && objLitGrandparent.IsObjectLit() {
			outermostObjectLit = objLitGrandparent
		} else {
			return outermostObjectLit
		}
	}
}

// GlobalNamespace_IsQnameDeclarationWithoutAssignment: true if the given
// Node is the GETPROP in a statement like `some.q.name;`
//
// Such do-nothing statements often have JSDoc on them and are intended to
// declare the qualified name.
//
// node is any Node, or even null.
//
// Java: GlobalNamespace.isQnameDeclarationWithoutAssignment(Node)
func GlobalNamespace_IsQnameDeclarationWithoutAssignment(node rhino.Node) bool {
	return !node.IsNil() && node.IsGetProp() && node.GetParent().IsExprResult()
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

// GlobalNamespace_NameProp is an enum of boolean properties of a Name, to be
// stored in one bit set per name.
//
// Non-private for @AutoValue. Code outside this class shouldn't need this
// enum though.
type GlobalNamespace_NameProp int32

const (
	// Increment the "index" passed to the constructor by 1 for each enum member
	GlobalNamespace_NameProp_IS_DECLARED              GlobalNamespace_NameProp = iota // index 0
	GlobalNamespace_NameProp_IS_MODULE_PROP                                           // index 1
	GlobalNamespace_NameProp_IS_PROVIDED                                              // index 2: If this name was in any goog.provide() calls.
	GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY                                 // index 3
	GlobalNamespace_NameProp_IS_EXTERN                                                // index 4: corresponds to SourceKind.EXTERN

	// Mutually exclusive properties indicating what kind of Closure type this Name is, if any.
	// Corresponds to DECLARED_TYPE_KIND_MASK.
	GlobalNamespace_NameProp_CONSTRUCTOR_TYPE // index 5: a non-interface class {} or a `/** @constructor */ function`
	GlobalNamespace_NameProp_INTERFACE_TYPE   // index 6: a `/** @interface */` or `/** @record */`
	GlobalNamespace_NameProp_ENUM_TYPE        // index 7: an `/** @enum */`
	GlobalNamespace_NameProp_NOT_A_TYPE       // index 8

	// Mutually exclusive properties indicating what kind of JavaScript entity this Name is.
	// Corresponds to NAME_KIND_MASK.
	GlobalNamespace_NameProp_CLASS        // index 9: class C {}
	GlobalNamespace_NameProp_OBJECTLIT    // index 10: var x = {};
	GlobalNamespace_NameProp_FUNCTION     // index 11: function f() {}
	GlobalNamespace_NameProp_GET_SET      // index 12: a getter, setter, or both; e.g. `obj.b` in `const obj = {set b(x) {}};`
	GlobalNamespace_NameProp_OTHER_OBJECT // index 13: anything else, including `var x = 1;`, var x = new Something();`, etc.
)

var globalNamespace_NameProp_names = [...]string{
	"IS_DECLARED",
	"IS_MODULE_PROP",
	"IS_PROVIDED",
	"IS_USED_HAS_OWN_PROPERTY",
	"IS_EXTERN",
	"CONSTRUCTOR_TYPE",
	"INTERFACE_TYPE",
	"ENUM_TYPE",
	"NOT_A_TYPE",
	"CLASS",
	"OBJECTLIT",
	"FUNCTION",
	"GET_SET",
	"OTHER_OBJECT",
}

// globalNamespace_NameProp_index is the enum constructor argument (the
// index of the property's bit), indexed by ordinal.
var globalNamespace_NameProp_index = [...]int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

// Java: GlobalNamespace.NameProp.DECLARED_TYPE_KIND_MASK
var GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK = GlobalNamespace_NameProp_CONSTRUCTOR_TYPE.Bit() |
	GlobalNamespace_NameProp_INTERFACE_TYPE.Bit() |
	GlobalNamespace_NameProp_ENUM_TYPE.Bit() |
	GlobalNamespace_NameProp_NOT_A_TYPE.Bit()

// Java: GlobalNamespace.NameProp.NAME_KIND_MASK
var GlobalNamespace_NameProp_NAME_KIND_MASK = GlobalNamespace_NameProp_CLASS.Bit() |
	GlobalNamespace_NameProp_OBJECTLIT.Bit() |
	GlobalNamespace_NameProp_FUNCTION.Bit() |
	GlobalNamespace_NameProp_GET_SET.Bit() |
	GlobalNamespace_NameProp_OTHER_OBJECT.Bit()

// Java: GlobalNamespace.NameProp.name()
func (p GlobalNamespace_NameProp) String() string { return globalNamespace_NameProp_names[p] }

// Java: GlobalNamespace.NameProp.ordinal()
func (p GlobalNamespace_NameProp) Ordinal() int { return int(p) }

// Bit is the enum field `bit` (some power of 2): `1 << index`.
func (p GlobalNamespace_NameProp) Bit() int32 {
	return jutil.Shl32(1, int64(globalNamespace_NameProp_index[p]))
}

// GlobalNamespace_ObjLitStringKeyAnalysis is the record
// ObjLitStringKeyAnalysis(@Nullable String nameString, @Nullable NameProp
// nameType).
type GlobalNamespace_ObjLitStringKeyAnalysis struct {
	nameString *string
	nameType   GlobalNamespace_NameProp
}

// Java: GlobalNamespace.ObjLitStringKeyAnalysis(String, NameProp)
func NewGlobalNamespace_ObjLitStringKeyAnalysis(nameString *string, nameType GlobalNamespace_NameProp) *GlobalNamespace_ObjLitStringKeyAnalysis {
	return &GlobalNamespace_ObjLitStringKeyAnalysis{nameString: nameString, nameType: nameType}
}

// Java: GlobalNamespace.ObjLitStringKeyAnalysis.nameString()
func (a *GlobalNamespace_ObjLitStringKeyAnalysis) NameString() *string { return a.nameString }

// Java: GlobalNamespace.ObjLitStringKeyAnalysis.nameType()
func (a *GlobalNamespace_ObjLitStringKeyAnalysis) NameType() GlobalNamespace_NameProp {
	return a.nameType
}

// GlobalNamespace_ObjLitStringKeyAnalysis_ForObjectDefineProperty: the
// object literal key is used to define a property.
// `Object.defineProperty(parent.qname, { strKeyName: value, { get: ..., } })`
//
// Java: GlobalNamespace.ObjLitStringKeyAnalysis.forObjectDefineProperty(String)
func GlobalNamespace_ObjLitStringKeyAnalysis_ForObjectDefineProperty(nameString string) *GlobalNamespace_ObjLitStringKeyAnalysis {
	// Technically the definition may not have a getter or setter, but we'll just
	// always pretend it does, because we cannot inline and collapse properties defined this
	// way.
	return NewGlobalNamespace_ObjLitStringKeyAnalysis(&nameString, GlobalNamespace_NameProp_GET_SET)
}

// GlobalNamespace_ObjLitStringKeyAnalysis_ForObjLitAssignment: the object
// literal key represents `parent.qname = { strKeyName: value }`.
//
// Java: GlobalNamespace.ObjLitStringKeyAnalysis.forObjLitAssignment(String, NameProp)
func GlobalNamespace_ObjLitStringKeyAnalysis_ForObjLitAssignment(nameString string, nameType GlobalNamespace_NameProp) *GlobalNamespace_ObjLitStringKeyAnalysis {
	return NewGlobalNamespace_ObjLitStringKeyAnalysis(&nameString, nameType)
}

// GlobalNamespace_ObjLitStringKeyAnalysis_ForNonReference: the object
// literal key does not represent a qualified name assignment.
//
// Java: GlobalNamespace.ObjLitStringKeyAnalysis.forNonReference()
func GlobalNamespace_ObjLitStringKeyAnalysis_ForNonReference() *GlobalNamespace_ObjLitStringKeyAnalysis {
	return NewGlobalNamespace_ObjLitStringKeyAnalysis( /* nameString= */ nil, GlobalNamespace_NameProp_OTHER_OBJECT)
}
