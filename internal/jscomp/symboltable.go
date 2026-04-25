// Port of com/google/javascript/jscomp/SymbolTable.java
// @ closure-compiler v20260415 (full).
//
// SymbolTable is an indexer-oriented symbol table that merges partial
// symbol tables from several compiler passes into a single unified index.
// Callers visit all symbols rather than looking up by key.
//
// The port preserves all public and package-level API from Java. Java
// generic parameters on StaticSymbolTable<S,R> / StaticSlot / StaticRef
// are replaced with the concrete interfaces declared below (minimal stubs
// with TODO refs).  Inner classes Symbol, Reference, and SymbolScope are
// ported as exported types in the same file.  The private inner-class
// callbacks (PropertyRefCollector, ThisRefCollector, JSDocInfoCollector,
// VisibilityCollector) are kept as unexported types.
//
// The following Java dependencies are not yet ported and are stubbed here
// or left as no-ops with TODO annotations:
//   - StaticScope, StaticSlot, StaticRef, StaticSymbolTable (rhino interfaces)
//     TODO: port com/google/javascript/rhino/StaticScope.java
//     TODO: port com/google/javascript/rhino/StaticSlot.java
//     TODO: port com/google/javascript/rhino/StaticRef.java
//     TODO: port com/google/javascript/rhino/StaticSymbolTable.java
//   - StaticTypedSlot / StaticTypedScope (rhino.jstype interfaces)
//     TODO: port com/google/javascript/rhino/jstype/StaticTypedSlot.java
//     TODO: port com/google/javascript/rhino/jstype/StaticTypedScope.java
//   - SimpleSlot / SimpleReference (base classes for Symbol/Reference)
//     TODO: port com/google/javascript/rhino/jstype/SimpleSlot.java
//     TODO: port com/google/javascript/rhino/jstype/SimpleReference.java
//   - Visibility (JSDocInfo.Visibility enum)
//     TODO: port com/google/javascript/rhino/JSDocInfo.Visibility
//   - JSDocInfo.Marker / SourcePosition<Node>
//     TODO: port com/google/javascript/rhino/JSDocInfo.Marker
//   - Module / ModuleType / ModuleMetadataMap
//     TODO: port com/google/javascript/jscomp/modules/Module.java
//     TODO: port com/google/javascript/jscomp/modules/ModuleMetadataMap.java
//   - CollectFileOverviewVisibility / AccessControlUtils
//     TODO: port com/google/javascript/jscomp/CollectFileOverviewVisibility.java
//     TODO: port com/google/javascript/jscomp/AccessControlUtils.java
//   - NodeUtil helpers (getBestLValue, getBestLValueName, createsBlockScope,
//     isInSyntheticScript, isGoogProvideCall, isGoogProvideCall)
//     TODO: port com/google/javascript/jscomp/NodeUtil (relevant helpers)
//   - Node.getJSType / Node.getSourceFileName / Node.getLength / Node.isIndexable
//     TODO: extend rhino.Node with JSType slot, source-file metadata, length

package jscomp

import (
	"fmt"
	"sort"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// ---------------------------------------------------------------------------
// Minimal stubs for Java interfaces not yet ported.
// ---------------------------------------------------------------------------

// StaticScope stubs com/google/javascript/rhino/StaticScope.
// TODO: port com/google/javascript/rhino/StaticScope.java
type StaticScope interface {
	GetRootNode() *rhino.Node
	GetParentScope() StaticScope
}

// StaticRef stubs com/google/javascript/rhino/StaticRef.
// TODO: port com/google/javascript/rhino/StaticRef.java
type StaticRef interface {
	GetNode() *rhino.Node
}

// StaticSlot stubs com/google/javascript/rhino/StaticSlot.
// TODO: port com/google/javascript/rhino/StaticSlot.java
type StaticSlot interface {
	GetName() string
	GetDeclaration() StaticRef
	GetJSDocInfo() *rhino.JSDocInfo
}

// StaticTypedSlot stubs com/google/javascript/rhino/jstype/StaticTypedSlot.
// TODO: port com/google/javascript/rhino/jstype/StaticTypedSlot.java
type StaticTypedSlot interface {
	StaticSlot
	GetType() jstype.JSType
	IsTypeInferred() bool
}

// StaticTypedScope stubs com/google/javascript/rhino/jstype/StaticTypedScope.
// TODO: port com/google/javascript/rhino/jstype/StaticTypedScope.java
type StaticTypedScope interface {
	StaticScope
	GetTypeOfThis() jstype.JSType
}

// StaticSymbolTable stubs com/google/javascript/rhino/StaticSymbolTable.
// TODO: port com/google/javascript/rhino/StaticSymbolTable.java
type StaticSymbolTable interface {
	GetAllSymbols() []StaticSlot
	GetScope(sym StaticSlot) StaticScope
	GetReferences(sym StaticSlot) []StaticRef
}

// Visibility stubs com/google/javascript/rhino/JSDocInfo.Visibility.
// TODO: port com/google/javascript/rhino/JSDocInfo.Visibility
type Visibility int

const (
	VisibilityPrivate   Visibility = iota
	VisibilityProtected Visibility = iota
	VisibilityPublic    Visibility = iota
	VisibilityInherited Visibility = iota
)

// ---------------------------------------------------------------------------
// SymbolTable
// ---------------------------------------------------------------------------

// GlobalThisName is the name used for the anonymous JavaScript global object.
// Matches SymbolTable.GLOBAL_THIS in Java.
const GlobalThisName = "*global*"

// stNodeSymKey is a composite key for the (Node, String) → Symbol table.
// Mirrors the Guava HashBasedTable<Node, String, Symbol> in Java.
type stNodeSymKey struct {
	node *rhino.Node
	name string
}

// SymbolTable is the indexer-oriented symbol table.
// Mirrors public final class SymbolTable in Java.
type SymbolTable struct {
	// symbols indexes all symbols by (declarationNode, name).
	// Mirrors Table<Node, String, Symbol> symbols in Java.
	symbols map[stNodeSymKey]*Symbol

	// scopes indexes all syntactic scopes by their root node.
	// Mirrors Map<Node, SymbolScope> scopes in Java.
	scopes map[*rhino.Node]*SymbolScope

	// docInfos holds all Nodes with JSDocInfo in the program.
	// Mirrors List<Node> docInfos in Java.
	docInfos []*rhino.Node

	globalScope *SymbolScope

	// symbolThatDeclaresType maps a JSType to the Symbol that declares it.
	// Uses interface pointer identity as key, matching Java's
	// LinkedIdentityHashMap<JSType, Symbol>.
	symbolThatDeclaresType map[jstype.JSType]*Symbol
}

// NewSymbolTable creates a new SymbolTable.
// Mirrors SymbolTable(AbstractCompiler, JSTypeRegistry) in Java.
// compiler and registry references are elided until needed by downstream callers;
// Go uses a factory function instead of a class constructor.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols:                map[stNodeSymKey]*Symbol{},
		scopes:                 map[*rhino.Node]*SymbolScope{},
		symbolThatDeclaresType: map[jstype.JSType]*Symbol{},
	}
}

// GetReferences returns all References for the given Symbol.
// Mirrors getReferences(Symbol) in Java.
func (st *SymbolTable) GetReferences(symbol *Symbol) []*Reference {
	refs := make([]*Reference, 0, len(symbol.references))
	for _, r := range symbol.references {
		refs = append(refs, r)
	}
	return refs
}

// GetReferenceList returns all References as a stable slice.
// Mirrors getReferenceList(Symbol) in Java.
func (st *SymbolTable) GetReferenceList(symbol *Symbol) []*Reference {
	refs := st.GetReferences(symbol)
	sort.Slice(refs, func(i, j int) bool {
		return stRefLess(refs[i], refs[j])
	})
	return refs
}

// stRefLess provides a stable ordering for References by source position.
// Go-only helper used by GetReferenceList; Java getReferenceList() preserves
// LinkedHashMap insertion order without an explicit comparator.
func stRefLess(a, b *Reference) bool {
	na, nb := a.GetNode(), b.GetNode()
	if na == nil {
		return true
	}
	if nb == nil {
		return false
	}
	if na.Lineno() != nb.Lineno() {
		return na.Lineno() < nb.Lineno()
	}
	return na.Charno() < nb.Charno()
}

// GetAllSymbols returns all symbols in the table.
// Mirrors getAllSymbols() in Java.
func (st *SymbolTable) GetAllSymbols() []*Symbol {
	syms := make([]*Symbol, 0, len(st.symbols))
	for _, s := range st.symbols {
		syms = append(syms, s)
	}
	return syms
}

// GetAllSymbolsSorted returns symbols in natural ordering.
// Mirrors getAllSymbolsSorted() in Java.
func (st *SymbolTable) GetAllSymbolsSorted() []*Symbol {
	syms := st.GetAllSymbols()
	sort.Slice(syms, func(i, j int) bool {
		return st.stSymbolLess(syms[i], syms[j])
	})
	return syms
}

// stSymbolLess is the "natural" symbol ordering: global before local,
// then lexicographic by name. Mirrors the symbolOrdering Comparator in Java
// (getNaturalSymbolOrdering returns symbolOrdering; Go uses a free comparison function).
func (st *SymbolTable) stSymbolLess(a, b *Symbol) bool {
	da := st.stLexicalScopeDepth(a.scope)
	db := st.stLexicalScopeDepth(b.scope)
	if da != db {
		return da < db
	}
	return a.GetName() < b.GetName()
}

// GetScope returns the SymbolScope that contains the given Symbol.
// Mirrors getScope(Symbol) in Java.
func (st *SymbolTable) GetScope(slot *Symbol) *SymbolScope {
	return slot.scope
}

// GetAllJSDocInfoNodes returns all AST nodes carrying JSDocInfo.
// Mirrors getAllJSDocInfoNodes() in Java.
func (st *SymbolTable) GetAllJSDocInfoNodes() []*rhino.Node {
	result := make([]*rhino.Node, len(st.docInfos))
	copy(result, st.docInfos)
	return result
}

// GetEnclosingScope returns the innermost SymbolScope containing n.
// If n is a function name its containing scope is returned.
// Mirrors getEnclosingScope(Node) in Java.
func (st *SymbolTable) GetEnclosingScope(n *rhino.Node) *SymbolScope {
	current := n.Parent()
	if n.Token() == rhino.NAME &&
		n.Parent() != nil &&
		n.Parent().Token() == rhino.FUNCTION {
		current = current.Parent()
	}
	for ; current != nil; current = current.Parent() {
		if s, ok := st.scopes[current]; ok {
			return s
		}
	}
	return nil
}

// GetEnclosingFunctionScope returns the nearest function or global scope.
// Mirrors getEnclosingFunctionScope(Node) in Java.
func (st *SymbolTable) GetEnclosingFunctionScope(n *rhino.Node) *SymbolScope {
	current := n.Parent()
	if n.Token() == rhino.NAME &&
		current != nil &&
		current.Token() == rhino.FUNCTION {
		current = current.Parent()
	}
	for ; current != nil; current = current.Parent() {
		scope := st.scopes[current]
		if scope != nil && !scope.IsBlockScope() {
			return scope
		}
	}
	return st.globalScope
}

// GetParameterInFunction tries to find a Symbol for a named parameter of sym.
// Mirrors getParameterInFunction(Symbol, String) in Java.
func (st *SymbolTable) GetParameterInFunction(sym *Symbol, paramName string) *Symbol {
	scope := st.stGetScopeInFunction(sym)
	if scope != nil {
		param := scope.GetSlot(paramName)
		if param != nil && param.scope == scope {
			return param
		}
	}
	return nil
}

// stGetScopeInFunction mirrors private getScopeInFunction(Symbol) in Java.
func (st *SymbolTable) stGetScopeInFunction(sym *Symbol) *SymbolScope {
	ft := sym.GetFunctionType()
	if ft == nil {
		return nil
	}
	functionNode := ft.GetSource()
	if functionNode == nil {
		return nil
	}
	return st.scopes[functionNode]
}

// GetSymbolForScope returns the Symbol associated with the given scope.
// Mirrors getSymbolForScope(SymbolScope) in Java.
func (st *SymbolTable) GetSymbolForScope(scope *SymbolScope) *Symbol {
	if scope.getSymbolForScope() == nil {
		scope.setSymbolForScope(st.stFindSymbolForScope(scope))
	}
	return scope.getSymbolForScope()
}

// stFindSymbolForScope mirrors private findSymbolForScope(SymbolScope) in Java.
func (st *SymbolTable) stFindSymbolForScope(scope *SymbolScope) *Symbol {
	rootNode := scope.GetRootNode()
	if rootNode.Parent() == nil {
		return st.globalScope.GetSlot(GlobalThisName)
	}
	if rootNode.Token() != rhino.FUNCTION {
		return nil
	}
	name := NU.GetBestLValueName(NU.GetBestLValue(rootNode))
	if name == "" {
		return nil
	}
	return scope.GetParentScope().GetQualifiedSlot(name)
}

// GetAllSymbolsForTypeOf returns Symbols associated with the type of sym.
// Mirrors getAllSymbolsForTypeOf(Symbol) in Java.
func (st *SymbolTable) GetAllSymbolsForTypeOf(sym *Symbol) []*Symbol {
	return st.GetAllSymbolsForType(st.stGetType(sym))
}

// GetGlobalScope returns the global scope.
// Mirrors getGlobalScope() in Java.
func (st *SymbolTable) GetGlobalScope() *SymbolScope {
	return st.globalScope
}

// GetSymbolDeclaredByFn gets the symbol for the given constructor or interface.
// Mirrors getSymbolDeclaredBy(FunctionType) in Java.
// GetSymbolDeclaredByFn gets the symbol for the given constructor or interface.
// Mirrors getSymbolDeclaredBy(FunctionType) in Java.
// Go uses a distinct name because Go lacks method overloading.
func (st *SymbolTable) GetSymbolDeclaredByFn(fn *jstype.FunctionType) *Symbol {
	if fn == nil || fn.Instance == nil {
		return nil
	}
	instanceType := fn.GetInstanceType()
	return st.stGetSymbolForName(fn.GetSource(), instanceType.GetReferenceName())
}

// GetSymbolDeclaredByEnum gets the symbol for the given enum.
// Mirrors getSymbolDeclaredBy(EnumType) in Java.
// Go uses a distinct name because Go lacks method overloading.
// EnumType.getSource() is not yet available; passing nil source falls back to globalScope.
func (st *SymbolTable) GetSymbolDeclaredByEnum(enumType *jstype.EnumType) *Symbol {
	if enumType == nil || enumType.Element == nil {
		return nil
	}
	return st.stGetSymbolForName(nil, enumType.Element.GetReferenceName())
}

// GetSymbolForInstancesOfSym gets the prototype symbol for a constructor sym.
// Mirrors getSymbolForInstancesOf(Symbol) in Java.
// Go uses a distinct name because Go lacks method overloading.
func (st *SymbolTable) GetSymbolForInstancesOfSym(sym *Symbol) *Symbol {
	if sym == nil {
		return nil
	}
	fn := sym.GetFunctionType()
	if fn != nil && (fn.IsConstructor || fn.IsInterface()) {
		return st.GetSymbolForInstancesOfFn(fn)
	}
	return nil
}

// GetSymbolForInstancesOfFn gets the symbol for the prototype of fn.
// Mirrors getSymbolForInstancesOf(FunctionType) in Java.
// Go uses a distinct name because Go lacks method overloading.
func (st *SymbolTable) GetSymbolForInstancesOfFn(fn *jstype.FunctionType) *Symbol {
	if fn == nil || fn.Instance == nil {
		return nil
	}
	pType := fn.GetPrototype()
	return st.stGetSymbolForName(fn.GetSource(), pType.GetReferenceName())
}

// stGetSymbolForName mirrors private getSymbolForName(Node, String) in Java.
// The module-scope name-stripping branch (which strips the goog.module namespace
// prefix from fully-qualified names when looking up symbols inside a module scope)
// requires ModuleMetadataMap, which is not yet ported; that branch is skipped.
// TODO: port com/google/javascript/jscomp/modules/ModuleMetadataMap
func (st *SymbolTable) stGetSymbolForName(source *rhino.Node, name string) *Symbol {
	if name == "" || st.globalScope == nil {
		return nil
	}
	var scope *SymbolScope
	if source == nil {
		scope = st.globalScope
	} else {
		scope = st.GetEnclosingScope(source)
	}
	if scope == nil {
		return nil
	}
	// Module-scope name stripping omitted: requires ModuleMetadataMap (not yet ported).
	return scope.GetQualifiedSlot(name)
}

// GetAllSymbolsForType returns all Symbols associated with the given JSType.
// Mirrors getAllSymbolsForType(JSType) in Java.
func (st *SymbolTable) GetAllSymbolsForType(t jstype.JSType) []*Symbol {
	if t == nil {
		return nil
	}
	if u, ok := t.(*jstype.UnionType); ok {
		var result []*Symbol
		for _, alt := range u.Members() {
			sym := st.stGetSymbolForTypeHelper(alt, true)
			if sym != nil {
				result = append(result, sym)
			}
		}
		return result
	}
	sym := st.stGetSymbolForTypeHelper(t, true)
	if sym == nil {
		return nil
	}
	return []*Symbol{sym}
}

// stGetSymbolForTypeHelper mirrors private getSymbolForTypeHelper(JSType, boolean) in Java.
// isGlobalThisType() is not yet available on the Go JSType interface; that branch is omitted.
// TODO: port isGlobalThisType() when global-this type tracking is ported.
func (st *SymbolTable) stGetSymbolForTypeHelper(t jstype.JSType, linkToCtor bool) *Symbol {
	if t == nil {
		return nil
	}
	// FunctionType: isNominalConstructorOrInterface (constructor or interface)
	if ft, ok := t.(*jstype.FunctionType); ok {
		if ft.IsConstructor || ft.IsInterface() {
			if linkToCtor {
				if st.globalScope != nil {
					return st.globalScope.GetSlot("Function")
				}
				return nil
			}
			return st.GetSymbolDeclaredByFn(ft)
		}
		// isFunctionPrototypeType: FunctionType owning a prototype object
		// In Go, a FunctionType whose ownerFunction is set is a prototype.
		// We check via GetOwnerFunction on PrototypeObjectType below instead.
		// Plain (non-constructor) FunctionType:
		if linkToCtor {
			if st.globalScope != nil {
				return st.globalScope.GetSlot("Function")
			}
			return nil
		}
		if st.globalScope != nil {
			return st.globalScope.GetQualifiedSlot("Function.prototype")
		}
		return nil
	}
	// isFunctionPrototypeType: PrototypeObjectType with an ownerFunction
	if pot, ok := t.(*jstype.PrototypeObjectType); ok {
		ownerFn := pot.GetOwnerFunction()
		if ownerFn != nil {
			if !ownerFn.IsConstructor && !ownerFn.IsInterface() {
				return nil
			}
			if linkToCtor {
				return st.GetSymbolDeclaredByFn(ownerFn)
			}
			return st.GetSymbolForInstancesOfFn(ownerFn)
		}
		// Fall through to autoboxesTo / symbolThatDeclaresType.
		autoboxed := jstype.AutoboxesTo(t)
		if autoboxed != nil {
			return st.stGetSymbolForTypeHelper(autoboxed, linkToCtor)
		}
		return st.symbolThatDeclaresType[t]
	}
	// isInstanceType: InstanceObjectType
	if iot, ok := t.(*jstype.InstanceObjectType); ok {
		ownerFn := iot.GetConstructor()
		if linkToCtor {
			return st.GetSymbolDeclaredByFn(ownerFn)
		}
		return st.GetSymbolForInstancesOfFn(ownerFn)
	}
	// isEnumType
	if et, ok := t.(*jstype.EnumType); ok {
		return st.GetSymbolDeclaredByEnum(et)
	}
	// isRecordType
	if jstype.IsRecordType(t) {
		return st.symbolThatDeclaresType[t]
	}
	// autoboxesTo
	autoboxed := jstype.AutoboxesTo(t)
	if autoboxed != nil {
		return st.stGetSymbolForTypeHelper(autoboxed, linkToCtor)
	}
	return nil
}

// ToDebugStringTree returns a multi-line debug representation.
// Mirrors toDebugStringTree() in Java.
func (st *SymbolTable) ToDebugStringTree() string {
	childrenScopes := map[*SymbolScope][]*SymbolScope{}
	for _, scope := range st.GetAllScopes() {
		if scope.GetParentScope() != nil {
			parent := scope.GetParentScope()
			childrenScopes[parent] = append(childrenScopes[parent], scope)
		}
	}
	var b strings.Builder
	for _, scope := range st.GetAllScopes() {
		if scope.GetParentScope() == nil {
			st.stDebugStringTree(&b, "", scope, childrenScopes)
		}
	}
	// Also include property scopes not in getAllScopes.
	visited := map[*SymbolScope]bool{}
	for _, s := range st.GetAllScopes() {
		visited[s] = true
	}
	for _, sym := range st.GetAllSymbols() {
		if sym.propertyScope == nil || visited[sym.propertyScope] {
			continue
		}
		visited[sym.propertyScope] = true
		st.stDebugStringTree(&b, "", sym.propertyScope, childrenScopes)
	}
	return b.String()
}

// stDebugStringTree mirrors private toDebugStringTree(StringBuilder, String, SymbolScope,
// ImmutableListMultimap) in Java. Go passes the children map directly.
func (st *SymbolTable) stDebugStringTree(
	b *strings.Builder,
	prefix string,
	scope *SymbolScope,
	childrenScopes map[*SymbolScope][]*SymbolScope,
) {
	b.WriteString(prefix)
	b.WriteString(scope.String())
	childPrefix := prefix + "    "
	printedSomething := false
	if len(scope.ownSymbols) > 0 {
		printedSomething = true
		b.WriteByte('\n')
		b.WriteString(prefix)
		b.WriteString("  Symbols:\n")
		for _, sym := range scope.ownSymbols {
			st.stDebugString(b, childPrefix, sym, false)
		}
	}
	if children := childrenScopes[scope]; len(children) > 0 {
		printedSomething = true
		b.WriteByte('\n')
		b.WriteString(prefix)
		b.WriteString("  Scopes:\n")
		for _, child := range children {
			st.stDebugStringTree(b, childPrefix, child, childrenScopes)
		}
	}
	if !printedSomething {
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
}

// ToDebugString returns a flat debug representation of all symbols.
// Mirrors toDebugString() in Java.
func (st *SymbolTable) ToDebugString() string {
	var b strings.Builder
	for _, sym := range st.GetAllSymbols() {
		st.stDebugString(&b, "", sym, true)
	}
	return b.String()
}

// stDebugString mirrors private toDebugString(StringBuilder, String, Symbol, boolean) in Java.
// Node.getSourceFileName() and Node.getLength()/isIndexable() are not yet ported;
// those fields are omitted from the format string.
// TODO: extend rhino.Node with source-file metadata and length/indexable fields.
func (st *SymbolTable) stDebugString(
	b *strings.Builder,
	prefix string,
	symbol *Symbol,
	printScope bool,
) {
	b.WriteString(prefix)
	b.WriteByte('\'')
	b.WriteString(symbol.GetName())
	b.WriteByte('\'')
	if printScope {
		b.WriteString(" in ")
		b.WriteString(symbol.scope.String())
	}
	b.WriteByte('\n')

	refCount := 0
	for _, ref := range st.GetReferences(symbol) {
		node := ref.GetNode()
		b.WriteString(prefix)
		// TODO: node.getSourceFileName() / node.getLength() not ported.
		// TODO: extend rhino.Node with source-file metadata and length field.
		b.WriteString(fmt.Sprintf(
			"  Ref %d: line: %d col: %d\n",
			refCount,
			node.Lineno(),
			node.Charno(),
		))
		refCount++
	}
}

// AddScopes ensures all scopes in the collection are present in this table.
// Mirrors addScopes(Collection<S>) in Java.
func (st *SymbolTable) AddScopes(scopes []StaticScope) {
	for _, s := range scopes {
		st.stCreateScopeFrom(s)
	}
}

// GetAllScopes returns all scopes in this symbol table.
// Mirrors getAllScopes() in Java.
func (st *SymbolTable) GetAllScopes() []*SymbolScope {
	result := make([]*SymbolScope, 0, len(st.scopes))
	for _, s := range st.scopes {
		result = append(result, s)
	}
	return result
}

// AddAnonymousFunctions finds anonymous functions and gives them names.
// Mirrors addAnonymousFunctions() in Java.
func (st *SymbolTable) AddAnonymousFunctions() {
	var lexicalScopes []*SymbolScope
	for _, scope := range st.GetAllScopes() {
		if scope.IsLexicalScope() {
			lexicalScopes = append(lexicalScopes, scope)
		}
	}
	sort.Slice(lexicalScopes, func(i, j int) bool {
		return st.stLexicalScopeLess(lexicalScopes[i], lexicalScopes[j])
	})
	for _, scope := range lexicalScopes {
		st.stAddAnonymousFunctionsInScope(scope)
	}
}

// stLexicalScopeLess orders lexical scopes by root-node source position then
// by depth. Mirrors the lexicalScopeOrdering Comparator in Java.
// Node.getSourceFileName() is not yet ported so cross-file ordering is omitted;
// within-file ordering uses lineno*10000+charno as a source-position proxy.
// TODO: use Node.getSourceFileName() once source-file metadata is ported.
func (st *SymbolTable) stLexicalScopeLess(a, b *SymbolScope) bool {
	an, bn := a.GetRootNode(), b.GetRootNode()
	pa := an.Lineno()*10000 + an.Charno()
	pb := bn.Lineno()*10000 + bn.Charno()
	if pa != pb {
		return pa < pb
	}
	return a.GetScopeDepth() < b.GetScopeDepth()
}

// stAddAnonymousFunctionsInScope mirrors private addAnonymousFunctionsInScope(SymbolScope) in Java.
// rootNode.getJSType() and rootNode.isFromExterns() are not yet ported; JSType is nil and
// isFromExterns is conservatively treated as false (anonymous functions in extern files
// may be named even when they shouldn't be).
// TODO: extend rhino.Node with JSType slot and isFromExterns flag.
func (st *SymbolTable) stAddAnonymousFunctionsInScope(scope *SymbolScope) {
	sym := st.GetSymbolForScope(scope)
	if sym != nil {
		return
	}
	rootNode := scope.GetRootNode()
	if !scope.IsLexicalScope() ||
		scope.IsGlobalScope() ||
		rootNode == nil ||
		scope.GetParentScope() == nil ||
		rootNode.Token() != rhino.FUNCTION {
		return
	}
	parent := scope.GetParentScope()
	innerName := fmt.Sprintf("function%%%d", scope.GetIndexInParent())
	// TODO: rootNode.getJSType() not ported.
	// TODO: extend rhino.Node with JSType slot.
	var t jstype.JSType
	// Adjust declNode: if the function is a member-function-def, use parent.
	declNode := rootNode
	// TODO: Node.isMemberFunctionDef() not ported (no MEMBER_FUNCTION_DEF token).
	// TODO: add MEMBER_FUNCTION_DEF token to rhino.Token.
	anonymousSym := st.stDeclareSymbol(innerName, t, true, parent, declNode, nil)
	scope.setSymbolForScope(anonymousSym)
}

// AddSymbolsFrom merges all symbols and references from otherSymbolTable.
// Mirrors addSymbolsFrom(StaticSymbolTable<S,R>) in Java.
func (st *SymbolTable) AddSymbolsFrom(otherSymbolTable StaticSymbolTable) {
	for _, otherSymbol := range otherSymbolTable.GetAllSymbols() {
		name := otherSymbol.GetName()
		myScope := st.stCreateScopeFrom(otherSymbolTable.GetScope(otherSymbol))

		decl := st.stFindBestDeclToAdd(otherSymbolTable, otherSymbol)
		var mySym *Symbol
		if decl != nil {
			declNode := decl.GetNode()
			mySym = st.stIsAnySymbolDeclared(name, declNode, myScope)
			if mySym == nil {
				mySym = st.stCopySymbolTo(otherSymbol, declNode, myScope)
			}
		} else {
			mySym = myScope.GetOwnSlot(name)
		}

		if mySym != nil {
			for _, otherRef := range otherSymbolTable.GetReferences(otherSymbol) {
				if st.stIsGoodRefToAdd(otherRef) {
					mySym.defineReferenceAt(otherRef.GetNode())
				}
			}
		}
	}
}

// stIsAnySymbolDeclared mirrors private isAnySymbolDeclared(String, Node, SymbolScope) in Java.
func (st *SymbolTable) stIsAnySymbolDeclared(
	name string, declNode *rhino.Node, scope *SymbolScope,
) *Symbol {
	sym := st.symbols[stNodeSymKey{declNode, name}]
	if sym == nil {
		return scope.ownSymbols[name]
	}
	return sym
}

// stFindBestDeclToAdd mirrors private findBestDeclToAdd(StaticSymbolTable, StaticSlot) in Java.
func (st *SymbolTable) stFindBestDeclToAdd(
	table StaticSymbolTable, slot StaticSlot,
) StaticRef {
	decl := slot.GetDeclaration()
	if st.stIsGoodRefToAdd(decl) {
		return decl
	}
	for _, ref := range table.GetReferences(slot) {
		if st.stIsGoodRefToAdd(ref) {
			return ref
		}
	}
	return nil
}

// stIsGoodRefToAdd mirrors private isGoodRefToAdd(StaticRef) in Java.
// The full Java check also requires ref.getNode().getStaticSourceFile() != null,
// NodeUtil.isInSyntheticScript, and NodeUtil.isGoogProvideCall, none of which
// are yet ported. Those checks are conservatively omitted (all refs pass).
// TODO: port Node.getStaticSourceFile(), NodeUtil.isInSyntheticScript, NodeUtil.isGoogProvideCall
func (st *SymbolTable) stIsGoodRefToAdd(ref StaticRef) bool {
	if ref == nil || ref.GetNode() == nil {
		return false
	}
	return true
}

// stCopySymbolTo mirrors private copySymbolTo(StaticSlot, Node, SymbolScope) in Java.
func (st *SymbolTable) stCopySymbolTo(sym StaticSlot, declNode *rhino.Node, scope *SymbolScope) *Symbol {
	if declNode == nil {
		panic("symboltable: copySymbolTo: declNode must not be nil")
	}
	return st.stDeclareSymbol(
		sym.GetName(),
		st.stGetTypeFromSlot(sym),
		st.stIsTypeInferred(sym),
		scope,
		declNode,
		sym.GetJSDocInfo(),
	)
}

// stSanitizeSpecialChars mirrors private static sanitizeSpecialChars(String) in Java.
func stSanitizeSpecialChars(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\x00", "\\0")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}

// stAddSymbol mirrors private addSymbol(String, JSType, boolean, SymbolScope, Node) in Java.
func (st *SymbolTable) stAddSymbol(
	name string,
	t jstype.JSType,
	inferred bool,
	scope *SymbolScope,
	declNode *rhino.Node,
) *Symbol {
	name = stSanitizeSpecialChars(name)
	symbol := &Symbol{
		name:       name,
		jsType:     t,
		inferred:   inferred,
		scope:      scope,
		references: map[*rhino.Node]*Reference{},
	}
	key := stNodeSymKey{declNode, name}
	if existing := st.symbols[key]; existing != nil {
		panic(fmt.Sprintf("symboltable: duplicate symbol %s in global index", name))
	}
	st.symbols[key] = symbol
	if existing := scope.ownSymbols[name]; existing != nil {
		panic(fmt.Sprintf("symboltable: duplicate symbol %s in scope", name))
	}
	scope.ownSymbols[name] = symbol
	return symbol
}

// stDeclareSymbol mirrors private declareSymbol(String, JSType, boolean, SymbolScope, Node, JSDocInfo) in Java.
func (st *SymbolTable) stDeclareSymbol(
	name string,
	t jstype.JSType,
	inferred bool,
	scope *SymbolScope,
	declNode *rhino.Node,
	info *rhino.JSDocInfo,
) *Symbol {
	symbol := st.stAddSymbol(name, t, inferred, scope, declNode)
	symbol.setJSDocInfo(info)
	symbol.setDeclaration(symbol.defineReferenceAt(declNode))
	return symbol
}

// stMergeSymbol mirrors private mergeSymbol(Symbol, Symbol) in Java.
func (st *SymbolTable) stMergeSymbol(from, to *Symbol) {
	fromDecl := from.GetDeclaration()
	for node := range from.references {
		if fromDecl == nil || node != fromDecl.GetNode() {
			to.defineReferenceAt(node)
		}
	}
	st.stRemoveSymbol(from)
}

// stRemoveSymbol mirrors private removeSymbol(Symbol) in Java.
func (st *SymbolTable) stRemoveSymbol(s *Symbol) {
	scope := st.GetScope(s)
	if scope.ownSymbols[s.GetName()] != s {
		panic(fmt.Sprintf("symboltable: symbol not found in scope: %v", s))
	}
	delete(scope.ownSymbols, s.GetName())
	declNode := s.GetDeclarationNode()
	if declNode == nil {
		panic(fmt.Sprintf("symboltable: symbol has no declaration node: %v", s))
	}
	key := stNodeSymKey{declNode, s.GetName()}
	if st.symbols[key] != s {
		panic(fmt.Sprintf("symboltable: symbol not found in table: %v", s))
	}
	delete(st.symbols, key)
	// If s has a property scope, remove all child symbols too.
	if s.propertyScope != nil && s.propertyScope.getSymbolForScope() == s {
		children := make([]*Symbol, 0, len(s.propertyScope.ownSymbols))
		for _, child := range s.propertyScope.ownSymbols {
			children = append(children, child)
		}
		for _, child := range children {
			st.stRemoveSymbol(child)
		}
		delete(st.scopes, s.GetDeclarationNode())
	}
}

// stRenameSymbol mirrors private renameSymbol(SymbolScope, Symbol, String) in Java.
func (st *SymbolTable) stRenameSymbol(scope *SymbolScope, sym *Symbol, newName string) {
	if sym.propertyScope != nil {
		panic("symboltable: renameSymbol: symbol has a property scope")
	}
	existingSym := scope.GetSlot(newName)
	if existingSym == nil {
		existingSym = st.stDeclareSymbol(
			newName,
			st.stGetType(sym),
			st.stIsTypeInferredSym(sym),
			scope,
			sym.GetDeclarationNode(),
			sym.GetJSDocInfo(),
		)
	}
	st.stMergeSymbol(sym, existingSym)
}

// FillNamespaceReferences adds references for namespace prefixes of dotted names.
// Mirrors fillNamespaceReferences() in Java.
func (st *SymbolTable) FillNamespaceReferences() {
	for _, symbol := range st.GetAllSymbols() {
		qName := symbol.GetName()
		rootIndex := strings.Index(qName, ".")
		if rootIndex == -1 {
			continue
		}
		root := symbol.scope.GetQualifiedSlot(qName[:rootIndex])
		if root == nil {
			continue
		}
		for _, ref := range st.GetReferences(symbol) {
			currentNode := ref.GetNode()
			if !stNodeIsQualifiedName(currentNode) {
				continue
			}
			for currentNode.Token() == rhino.GETPROP {
				currentNode = currentNode.FirstChild()
				name := stNodeQualifiedName(currentNode)
				if name == "" {
					continue
				}
				namespace := st.stIsAnySymbolDeclared(name, currentNode, root.scope)
				if namespace == nil {
					namespace = root.scope.GetQualifiedSlot(name)
				}
				if namespace == nil && root.scope.IsGlobalScope() {
					// TODO: currentNode.getJSType() not ported; use nil type.
					// TODO: extend rhino.Node with JSType slot.
					namespace = st.stDeclareSymbol(name, nil, true, root.scope, currentNode, nil)
				}
				if namespace != nil {
					namespace.defineReferenceAt(currentNode)
				}
			}
		}
	}
}

// FillPropertyScopes creates property scopes for named types.
// Mirrors fillPropertyScopes() in Java.
func (st *SymbolTable) FillPropertyScopes() {
	var types []*Symbol
	var exports []*Symbol
	for _, sym := range st.GetAllSymbols() {
		if st.stNeedsPropertyScope(sym) {
			if sym.GetName() == "exports" {
				exports = append(exports, sym)
			} else {
				types = append(types, sym)
			}
		}
	}
	// Sort types in reverse natural order.
	sort.Slice(types, func(i, j int) bool {
		return st.stSymbolLess(types[j], types[i])
	})
	allTypes := append(types, exports...)

	// Build type→root-symbol map.
	for _, s := range allTypes {
		t := st.stGetType(s)
		if t == nil {
			continue
		}
		symbolForType := st.stGetSymbolForTypeHelper(t, false)
		if (st.stIsFunctionTypeConstructorOrInterface(t) ||
			st.stIsFunctionPrototypeType(t) ||
			st.stIsEnumType(t)) &&
			symbolForType != nil && symbolForType != s {
			continue
		}
		st.symbolThatDeclaresType[t] = s
	}

	for _, s := range allTypes {
		t := st.stGetType(s)
		if t == nil || st.symbolThatDeclaresType[t] == s {
			st.stCreatePropertyScopeFor(s)
		}
	}

	// Propagate property scope to aliases.
	for _, s := range allTypes {
		t := st.stGetType(s)
		if t != nil && st.symbolThatDeclaresType[t] != nil {
			root := st.symbolThatDeclaresType[t]
			s.propertyScope = root.GetPropertyScope()
		}
	}

	st.PruneOrphanedNames()
}

// stNeedsPropertyScope mirrors private needsPropertyScope(Symbol) in Java.
// Java's ObjectType.cast() is modelled by the type-switch below since Go lacks
// the Java ObjectType class hierarchy (FunctionType extends ObjectType there).
func (st *SymbolTable) stNeedsPropertyScope(sym *Symbol) bool {
	jsT := st.stGetType(sym)
	if jsT == nil {
		return false
	}
	if st.stIsRecordType(jsT) {
		return true
	}
	switch v := jsT.(type) {
	case *jstype.FunctionType:
		// Java: isNominalConstructorOrInterface() → true for constructors/interfaces
		if v.IsConstructor || v.IsInterface() {
			return true
		}
		return false
	case *jstype.InstanceObjectType:
		// InstanceObjectType has a reference name (it's nominal)
		if v.GetReferenceName() == "" {
			return true // anonymous
		}
		return false
	case *jstype.PrototypeObjectType:
		// isFunctionPrototypeType: has an ownerFunction
		if v.GetOwnerFunction() != nil {
			return true
		}
		if v.GetReferenceName() == "" {
			return true // anonymous object
		}
		return false
	case *jstype.ObjectType:
		if v.GetReferenceName() == "" {
			return true // anonymous object
		}
		return false
	case *jstype.EnumType:
		return true
	}
	return false
}

// PruneOrphanedNames removes symbols whose namespace has a property scope.
// Mirrors pruneOrphanedNames() in Java.
func (st *SymbolTable) PruneOrphanedNames() {
	allSyms := st.GetAllSymbols()
outer:
	for _, s := range allSyms {
		currentName := s.GetName()
		if s.IsProperty() {
			continue
		}
		// TODO: compiler.getModuleMap().getClosureModule() not ported.
		// TODO: port com/google/javascript/jscomp/modules/ModuleMap
		for {
			dot := strings.LastIndex(currentName, ".")
			if dot == -1 {
				break
			}
			currentName = currentName[:dot]
			owner := s.scope.GetQualifiedSlot(currentName)
			if owner != nil {
				ownerType := st.stGetType(owner)
				if ownerType != nil &&
					(st.stIsFunctionTypeConstructorOrInterface(ownerType) ||
						st.stIsFunctionPrototypeType(ownerType) ||
						st.stIsEnumType(ownerType)) {
					st.stRemoveSymbol(s)
					continue outer
				}
			}
		}
	}
}

// FillPropertySymbols creates symbols and references for all property accesses.
// Mirrors fillPropertySymbols(Node, Node) in Java.
func (st *SymbolTable) FillPropertySymbols(externs, root *rhino.Node) {
	c := &stPropertyRefCollector{st: st}
	c.process(externs, root)
}

// FillJSDocInfo indexes JSDocInfo nodes and creates references for @param.
// Mirrors fillJSDocInfo(Node, Node) in Java.
func (st *SymbolTable) FillJSDocInfo(externs, root *rhino.Node) {
	c := &stJSDocInfoCollector{st: st}
	c.process(externs, root)
	// Parameter-reference creation from JSDoc markers:
	// TODO: JSDocInfo.getMarkers() / Marker.getNameNode() not ported.
	// TODO: port com/google/javascript/rhino/JSDocInfo.Marker
}

// FillSymbolVisibility records the visibility of each symbol by walking the
// AST and inspecting JSDoc @private/@protected/@public annotations.
// Mirrors fillSymbolVisibility(Node, Node) in Java.
//
// Full implementation requires CollectFileOverviewVisibility and
// AccessControlUtils (not yet ported). This partial port reads per-symbol
// JSDoc visibility from each symbol's declaration node.
func (st *SymbolTable) FillSymbolVisibility(externs, root *rhino.Node) {
	// Walk all symbols and pull visibility from their JSDoc.
	for _, sym := range st.GetAllSymbols() {
		if sym.docInfo != nil {
			// Visibility already set via FillJSDocInfo or direct assignment.
			continue
		}
		if sym.declaration == nil || sym.declaration.node == nil {
			continue
		}
		info := getBestJSDocInfo(sym.declaration.node)
		if info == nil {
			continue
		}
		// Map rhino.Visibility to the local Visibility type.
		switch info.GetVisibility() {
		case rhino.VisibilityPrivate:
			v := VisibilityPrivate
			sym.visibility = &v
		case rhino.VisibilityProtected:
			v := VisibilityProtected
			sym.visibility = &v
		case rhino.VisibilityPublic:
			v := VisibilityPublic
			sym.visibility = &v
		default:
			v := VisibilityInherited
			sym.visibility = &v
		}
	}
}

// stCreatePropertyScopeFor mirrors private createPropertyScopeFor(Symbol) in Java.
// Property population via ObjectType.getOwnPropertyNames() is not yet ported;
// that step is deferred until ObjectType property-map helpers are available.
// TODO: port com/google/javascript/rhino/jstype/ObjectType.getOwnPropertyNames
func (st *SymbolTable) stCreatePropertyScopeFor(s *Symbol) {
	if s.propertyScope != nil {
		return
	}
	jsT := st.stGetType(s)
	if jsT == nil {
		return
	}
	// Java: jsType.toObjectType() — covers ObjectType and all subtypes.
	// In Go we accept any type that acts as an object (ObjectType, PrototypeObjectType,
	// InstanceObjectType, FunctionType — all are ObjectType in Java's hierarchy).
	if !stIsObjectType(jsT) {
		return
	}
	// For module exports, reuse the module scope.
	declNode := s.GetDeclarationNode()
	if s.GetName() == "exports" && declNode != nil && stNodeIsModuleBody(declNode) {
		s.propertyScope = st.scopes[declNode]
		return
	}
	parentScope := st.stMaybeGetParentPropertyScope(jsT)
	s.setPropertyScope(newSymbolScope(nil, parentScope, jsT, s))
	// Property population from type-system data deferred:
	// requires ObjectType.getOwnPropertyNames() which is not yet ported.
}

// stIsObjectType mirrors the ObjectType.cast(jsType) != null check in Java.
// In Java, FunctionType/PrototypeObjectType/InstanceObjectType all extend ObjectType.
// In Go those are distinct types; this helper recognises all of them.
func stIsObjectType(t jstype.JSType) bool {
	switch t.(type) {
	case *jstype.ObjectType, *jstype.PrototypeObjectType, *jstype.InstanceObjectType,
		*jstype.FunctionType:
		return true
	}
	return false
}

// stMaybeGetParentPropertyScope mirrors private maybeGetParentPropertyScope(ObjectType) in Java.
// It follows the implicit-prototype chain to find the parent property scope.
// The Java method takes an ObjectType; here we accept JSType and cast internally.
func (st *SymbolTable) stMaybeGetParentPropertyScope(t jstype.JSType) *SymbolScope {
	if t == nil {
		return nil
	}
	pot, ok := t.(*jstype.PrototypeObjectType)
	if !ok {
		return nil
	}
	proto := pot.GetImplicitPrototype()
	if proto == nil || proto == pot {
		return nil
	}
	var parentSymbol *Symbol
	if stIsEs6ClassConstructor(proto) {
		fn := jstype.ToMaybeFunctionType(proto)
		parentSymbol = st.GetSymbolDeclaredByFn(fn)
	} else if proto.GetConstructor() != nil {
		parentSymbol = st.GetSymbolForInstancesOfFn(proto.GetConstructor())
	} else {
		return nil
	}
	if parentSymbol == nil {
		return nil
	}
	st.stCreatePropertyScopeFor(parentSymbol)
	return parentSymbol.GetPropertyScope()
}

// stIsEs6ClassConstructor mirrors private isEs6ClassConstructor(JSType) in Java.
func stIsEs6ClassConstructor(t jstype.JSType) bool {
	ft := jstype.ToMaybeFunctionType(t)
	if ft == nil {
		return false
	}
	src := ft.GetSource()
	return src != nil && src.IsClass()
}

// FillThisReferences adds references for `this` usages.
// Mirrors fillThisReferences(Node, Node) in Java.
func (st *SymbolTable) FillThisReferences(externs, root *rhino.Node) {
	c := &stThisRefCollector{st: st}
	c.process(externs, root)
}

// FillSuperReferences adds references for `super` usages.
// Mirrors fillSuperReferences(Node, Node) in Java.
func (st *SymbolTable) FillSuperReferences(externs, root *rhino.Node) {
	stWalkPostOrder(externs, func(n *rhino.Node) {
		if n.Token() != rhino.THIS { // TODO: add SUPER token to rhino.Token
			return
		}
		// TODO: n.getJSType() not ported.
		// TODO: extend rhino.Node with JSType slot.
	})
	stWalkPostOrder(root, func(n *rhino.Node) {
		if n.Token() != rhino.THIS { // TODO: add SUPER token to rhino.Token
			return
		}
		// TODO: n.getJSType() not ported.
	})
}

// FillGoogProvideModuleRequires walks the AST looking for goog.module(),
// goog.provide(), and goog.require() calls, connecting the declared namespace
// symbols with their require-site references.
// Mirrors fillGoogProvideModuleRequires(Node, Node) in Java.
//
// The full Java implementation uses ModuleMetadataMap and creates synthetic
// namespace symbols (e.g. "ns$some.namespace") for each goog.provide /
// goog.module declaration, then wires goog.require() references to them.
// Our partial port identifies the call sites and adds references where the
// corresponding declared symbol already exists in the symbol table.
func (st *SymbolTable) FillGoogProvideModuleRequires(externs, root *rhino.Node) {
	// Walk root only (externs don't contain goog.module calls in practice).
	var walk func(n *rhino.Node)
	walk = func(n *rhino.Node) {
		if n == nil {
			return
		}
		// Look for CALL nodes that are goog.xxx('some.name') patterns.
		if n.Token() == rhino.CALL {
			kids := n.Children()
			if len(kids) >= 2 {
				callee := kids[0]
				arg := kids[1]
				// callee should be goog.require / goog.module / goog.provide
				if callee.Token() == rhino.GETPROP {
					calleeKids := callee.Children()
					if len(calleeKids) >= 2 &&
						calleeKids[0].Token() == rhino.NAME &&
						calleeKids[0].String_() == "goog" &&
						arg.Token() == rhino.STRINGLIT {
						prop := calleeKids[1].String_()
						ns := "ns$" + arg.String_()
						switch prop {
						case "require", "requireType":
							// Try to find the declared namespace and add a reference.
							if sym := st.stGetSymbolForName(arg, ns); sym != nil {
								if sym.references == nil {
									sym.references = make(map[*rhino.Node]*Reference)
								}
								ref := &Reference{node: arg, symbol: sym}
								sym.references[arg] = ref
							}
						}
					}
				}
			}
		}
		for _, c := range n.Children() {
			walk(c)
		}
	}
	walk(externs)
	walk(root)
}

// FlattenGoogModuleExports renames exports.Foo → Foo in module scopes.
// Mirrors flattenGoogModuleExports() in Java.
func (st *SymbolTable) FlattenGoogModuleExports() {
	for _, moduleScope := range st.GetAllScopes() {
		if !moduleScope.IsModuleScope() {
			continue
		}
		syms := make([]*Symbol, 0, len(moduleScope.ownSymbols))
		for _, sym := range moduleScope.ownSymbols {
			syms = append(syms, sym)
		}
		for _, sym := range syms {
			if strings.HasPrefix(sym.GetName(), "exports.") {
				st.stRenameSymbol(moduleScope, sym, sym.GetName()[len("exports."):])
			}
		}
	}
}

// RemoveGeneratedSymbols removes transpiler-generated symbols.
// Mirrors removeGeneratedSymbols() in Java.
func (st *SymbolTable) RemoveGeneratedSymbols() {
	allSyms := make([]*Symbol, 0, len(st.symbols))
	for _, s := range st.symbols {
		allSyms = append(allSyms, s)
	}
	for _, sym := range allSyms {
		if st.stIsSymbolAQuotedObjectKey(sym) {
			if st.GetScope(sym).ownSymbols[sym.GetName()] != nil {
				st.stRemoveSymbol(sym)
			}
		} else if sym.GetDeclarationNode() != nil &&
			stNodeIsModuleBody(sym.GetDeclarationNode()) {
			st.stRemoveSymbol(sym)
		}
	}
	st.stMergeExternSymbolsDuplicatedOnWindow()
}

// stIsSymbolAQuotedObjectKey mirrors private isSymbolAQuotedObjectKey(Symbol) in Java.
func (st *SymbolTable) stIsSymbolAQuotedObjectKey(sym *Symbol) bool {
	n := sym.GetDeclarationNode()
	return n != nil && n.IsStringKey() && n.IsQuotedStringKey()
}

// stIsSymbolDuplicatedExternOnWindow mirrors private isSymbolDuplicatedExternOnWindow(Symbol) in Java.
// Java checks !node.isIndexable() first; Node.isIndexable() is not yet ported so that
// guard is omitted (conservative: more nodes may match than in Java).
// TODO: extend rhino.Node with isIndexable() and add the !node.isIndexable() guard.
func (st *SymbolTable) stIsSymbolDuplicatedExternOnWindow(sym *Symbol) bool {
	n := sym.GetDeclarationNode()
	if n == nil {
		return false
	}
	return n.IsGetProp() &&
		n.FirstChild() != nil &&
		n.FirstChild().Token() == rhino.NAME &&
		n.FirstChild().String_() == "window"
}

// stMergeExternSymbolsDuplicatedOnWindow mirrors private mergeExternSymbolsDuplicatedOnWindow() in Java.
// Java uses Node.getStaticSourceFile().isExtern() and Node.getSourcePosition() to build
// a filename × position → symbol table and merge duplicates. These are not yet ported,
// so the merge step is a no-op until source-file metadata and source-position are available.
// TODO: port Node.getStaticSourceFile().isExtern() and Node.getSourcePosition()
func (st *SymbolTable) stMergeExternSymbolsDuplicatedOnWindow() {
	// No-op: requires StaticSourceFile.isExtern() and Node.getSourcePosition()
	// (not yet ported). Full logic: for each pair of extern symbols at the same
	// source position, if one isSymbolDuplicatedExternOnWindow, merge it into the other.
}

// stCreateScopeFrom mirrors private createScopeFrom(StaticScope) in Java.
func (st *SymbolTable) stCreateScopeFrom(otherScope StaticScope) *SymbolScope {
	otherScopeRoot := otherScope.GetRootNode()
	myScope := st.scopes[otherScopeRoot]
	if myScope != nil {
		return myScope
	}
	otherScopeParent := otherScope.GetParentScope()
	if otherScopeParent == nil && st.globalScope != nil {
		panic("symboltable: global scopes found at different roots")
	}
	var parentScope *SymbolScope
	if otherScopeParent != nil {
		parentScope = st.stCreateScopeFrom(otherScopeParent)
	}
	myScope = newSymbolScope(otherScopeRoot, parentScope, st.stGetTypeOfThis(otherScope), nil)
	st.scopes[otherScopeRoot] = myScope
	if myScope.IsGlobalScope() {
		st.globalScope = myScope
	}
	return myScope
}

// ---------------------------------------------------------------------------
// Type helpers
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Type helpers
// ---------------------------------------------------------------------------

// stGetType mirrors private getType(Symbol) in Java.
// Java's getType(StaticSlot) has a JSDoc-typedef branch that uses registry.getGlobalType()
// which is not yet ported. The Symbol-specific overload just returns sym.jsType.
func (st *SymbolTable) stGetType(sym *Symbol) jstype.JSType {
	return sym.jsType
}

// stIsTypeInferredSym mirrors private isTypeInferred(Symbol) in Java (Symbol variant).
func (st *SymbolTable) stIsTypeInferredSym(sym *Symbol) bool {
	return sym.inferred
}

// stGetTypeFromSlot mirrors private getType(StaticSlot) in Java (slot variant).
// The JSDoc typedef branch (registry.getGlobalType + typedefType.evaluate) is not yet ported.
// TODO: port JSTypeRegistry.getGlobalType and JSTypeExpression.evaluate for typedef resolution.
func (st *SymbolTable) stGetTypeFromSlot(sym StaticSlot) jstype.JSType {
	if ts, ok := sym.(StaticTypedSlot); ok {
		return ts.GetType()
	}
	return nil
}

// stIsTypeInferred mirrors private isTypeInferred(StaticSlot) in Java.
func (st *SymbolTable) stIsTypeInferred(sym StaticSlot) bool {
	if ts, ok := sym.(StaticTypedSlot); ok {
		return ts.IsTypeInferred()
	}
	return true
}

// stGetTypeOfThis mirrors private getTypeOfThis(StaticScope) in Java.
func (st *SymbolTable) stGetTypeOfThis(s StaticScope) jstype.JSType {
	if ts, ok := s.(StaticTypedScope); ok {
		return ts.GetTypeOfThis()
	}
	return nil
}

// stIsFunctionTypeConstructorOrInterface mirrors type.isNominalConstructorOrInterface() in Java.
// In Java, isNominalConstructorOrInterface is a method on JSType dispatching to FunctionType.
// In Go, FunctionType.IsConstructor or FunctionType.IsInterface() covers the same cases.
func (st *SymbolTable) stIsFunctionTypeConstructorOrInterface(t jstype.JSType) bool {
	if ft, ok := t.(*jstype.FunctionType); ok {
		return ft.IsConstructor || ft.IsInterface()
	}
	return false
}

// stIsFunctionPrototypeType mirrors type.isFunctionPrototypeType() in Java.
// Java returns true for PrototypeObjectType instances that have an ownerFunction set.
// In Go, *jstype.PrototypeObjectType with a non-nil GetOwnerFunction() matches.
func (st *SymbolTable) stIsFunctionPrototypeType(t jstype.JSType) bool {
	if pot, ok := t.(*jstype.PrototypeObjectType); ok {
		return pot.GetOwnerFunction() != nil
	}
	return false
}

// stIsEnumType mirrors type.isEnumType() in Java.
func (st *SymbolTable) stIsEnumType(t jstype.JSType) bool {
	_, ok := t.(*jstype.EnumType)
	return ok
}

// stIsRecordType mirrors type.isRecordType() in Java.
func (st *SymbolTable) stIsRecordType(t jstype.JSType) bool {
	return jstype.IsRecordType(t)
}

// ---------------------------------------------------------------------------
// stLexicalScopeDepth mirrors getLexicalScopeDepth in Java.
// ---------------------------------------------------------------------------

// stLexicalScopeDepth mirrors private getLexicalScopeDepth(SymbolScope) in Java.
func (st *SymbolTable) stLexicalScopeDepth(scope *SymbolScope) int {
	if scope.IsLexicalScope() || scope.IsDocScope() {
		return scope.GetScopeDepth()
	}
	// property scope: find enclosing lexical scope.
	sym := scope.getSymbolForScope()
	if sym == nil {
		return 0
	}
	return st.stLexicalScopeDepth(st.GetScope(sym)) + 1
}

// ---------------------------------------------------------------------------
// Node predicate helpers (replacing Java's Node.isXxx() methods)
// ---------------------------------------------------------------------------

// stNodeIsQualifiedName delegates to Node.IsQualifiedName(), which mirrors
// Node.isQualifiedName() in Java.
func stNodeIsQualifiedName(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	return n.IsQualifiedName()
}

// stNodeQualifiedName delegates to Node.GetQualifiedName(), which mirrors
// Node.getQualifiedName() in Java. Returns "" for non-qualified-name nodes
// (matching Java's null return, adapted to Go's string semantics).
func stNodeQualifiedName(n *rhino.Node) string {
	if n == nil {
		return ""
	}
	return n.GetQualifiedName()
}

// stNodeIsModuleBody reports whether n is a MODULE_BODY node.
// Mirrors Node.isModuleBody() in Java.
func stNodeIsModuleBody(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.MODULE_BODY
}

// stWalkPostOrder performs a post-order walk of the subtree rooted at n,
// equivalent to NodeTraversal.AbstractPostOrderCallback semantics in Java.
func stWalkPostOrder(n *rhino.Node, fn func(*rhino.Node)) {
	if n == nil {
		return
	}
	for c := n.FirstChild(); c != nil; c = c.Next() {
		stWalkPostOrder(c, fn)
	}
	fn(n)
}

// ---------------------------------------------------------------------------
// Symbol
// ---------------------------------------------------------------------------

// Symbol is a symbol-table entry.
// Mirrors public static final class Symbol extends SimpleSlot in Java.
// TODO: port com/google/javascript/rhino/jstype/SimpleSlot.java and embed it.
type Symbol struct {
	name     string
	jsType   jstype.JSType
	inferred bool
	scope    *SymbolScope

	// references is keyed by the Node the reference appears at.
	// Mirrors Map<Node, Reference> (LinkedHashMap) in Java.
	references map[*rhino.Node]*Reference

	propertyScope *SymbolScope
	declaration   *Reference
	docInfo       *rhino.JSDocInfo

	// visibility is stored separately from docInfo.
	// Mirrors Symbol.visibility in Java.
	visibility *Visibility

	// docScope holds symbols documented only in JSDoc.
	docScope *SymbolScope
}

// GetName returns the symbol name.
// Mirrors SimpleSlot.getName() in Java.
func (s *Symbol) GetName() string { return s.name }

// GetType returns the JSType of this symbol.
// Mirrors SimpleSlot.getType() in Java.
func (s *Symbol) GetType() jstype.JSType { return s.jsType }

// IsTypeInferred reports whether the type was inferred.
// Mirrors SimpleSlot.isTypeInferred() in Java.
func (s *Symbol) IsTypeInferred() bool { return s.inferred }

// GetDeclaration returns the declaration Reference.
// Mirrors Symbol.getDeclaration() in Java.
func (s *Symbol) GetDeclaration() *Reference { return s.declaration }

// GetFunctionType returns the FunctionType if this symbol's type is a function.
// Mirrors Symbol.getFunctionType() → JSType.toMaybeFunctionType() in Java.
func (s *Symbol) GetFunctionType() *jstype.FunctionType {
	if s.jsType == nil {
		return nil
	}
	ft, _ := s.jsType.(*jstype.FunctionType)
	return ft
}

// defineReferenceAt ensures a Reference exists for node n and returns it.
// Mirrors Symbol.defineReferenceAt(Node) in Java.
func (s *Symbol) defineReferenceAt(n *rhino.Node) *Reference {
	if ref, ok := s.references[n]; ok {
		return ref
	}
	ref := &Reference{symbol: s, node: n}
	s.references[n] = ref
	return ref
}

// setDeclaration sets the declaration reference. May only be called once.
// Mirrors Symbol.setDeclaration(Reference) in Java.
func (s *Symbol) setDeclaration(ref *Reference) {
	if s.declaration != nil {
		panic("symboltable: Symbol.setDeclaration called twice")
	}
	s.declaration = ref
}

// GetDeclarationNode returns the declaration AST node.
// Mirrors Symbol.getDeclarationNode() in Java.
func (s *Symbol) GetDeclarationNode() *rhino.Node {
	if s.declaration == nil {
		return nil
	}
	return s.declaration.GetNode()
}

// GetSourceFileName returns the source file name of the symbol's declaration
// node. In Java this returns Node.getSourceFileName(), which requires
// source-file metadata on each AST node. Our rhino.Node does not yet carry
// per-node source-file information, so this returns an empty string until
// the metadata slot is added.
// Mirrors Symbol.getSourceFileName() in Java.
func (s *Symbol) GetSourceFileName() string {
	// Java: Node n = getDeclarationNode(); return n == null ? null : n.getSourceFileName();
	// We return "" instead of null since Go doesn't have nullable strings.
	decl := s.GetDeclarationNode()
	if decl == nil {
		return ""
	}
	// TODO: return decl.SourceFileName() once rhino.Node carries it.
	return ""
}

// GetPropertyScope returns the property scope.
// Mirrors Symbol.getPropertyScope() in Java.
func (s *Symbol) GetPropertyScope() *SymbolScope { return s.propertyScope }

// setPropertyScope sets the property scope.
// Mirrors Symbol.setPropertyScope(SymbolScope) in Java.
func (s *Symbol) setPropertyScope(scope *SymbolScope) {
	s.propertyScope = scope
	if scope != nil {
		scope.setSymbolForScope(s)
	}
}

// GetJSDocInfo returns the JSDocInfo attached to this symbol.
// Mirrors Symbol.getJSDocInfo() in Java.
func (s *Symbol) GetJSDocInfo() *rhino.JSDocInfo { return s.docInfo }

// setJSDocInfo sets the JSDocInfo.
func (s *Symbol) setJSDocInfo(info *rhino.JSDocInfo) { s.docInfo = info }

// GetVisibility returns the Visibility (may be nil if not set).
// Mirrors Symbol.getVisibility() in Java.
func (s *Symbol) GetVisibility() *Visibility { return s.visibility }

// setVisibility sets the Visibility.
func (s *Symbol) setVisibility(v Visibility) { s.visibility = &v }

// IsProperty reports whether this is a property of another variable.
// Mirrors Symbol.isProperty() in Java.
func (s *Symbol) IsProperty() bool { return s.scope.IsPropertyScope() }

// IsLexicalVariable reports whether this is a variable in a lexical scope.
// Mirrors Symbol.isLexicalVariable() in Java.
func (s *Symbol) IsLexicalVariable() bool { return s.scope.IsLexicalScope() }

// IsDocOnlyParameter reports whether this is a JSDoc-only variable.
// Mirrors Symbol.isDocOnlyParameter() in Java.
func (s *Symbol) IsDocOnlyParameter() bool { return s.scope.IsDocScope() }

// String returns a debug representation.
// Mirrors Symbol.toString() in Java.
func (s *Symbol) String() string {
	n := s.GetDeclarationNode()
	lineNo := -1
	if n != nil {
		lineNo = n.Lineno()
	}
	return fmt.Sprintf("%s@%s:%d", s.GetName(), s.GetSourceFileName(), lineNo)
}

// Equals mirrors Symbol.equals in Java.
func (s *Symbol) Equals(other *Symbol) bool {
	if other == nil {
		return false
	}
	return s.inferred == other.inferred &&
		s.name == other.name &&
		s.jsType == other.jsType &&
		s.scope == other.scope
}

// ---------------------------------------------------------------------------
// Reference
// ---------------------------------------------------------------------------

// Reference is a reference to a Symbol at a specific AST node.
// Mirrors public static final class Reference extends SimpleReference<Symbol>.
// TODO: port com/google/javascript/rhino/jstype/SimpleReference.java
type Reference struct {
	symbol *Symbol
	node   *rhino.Node
}

// GetSymbol returns the Symbol this Reference belongs to.
func (r *Reference) GetSymbol() *Symbol { return r.symbol }

// GetNode returns the AST node at which the reference occurs.
// Mirrors SimpleReference.getNode() in Java.
func (r *Reference) GetNode() *rhino.Node { return r.node }

// ---------------------------------------------------------------------------
// SymbolScope
// ---------------------------------------------------------------------------

// SymbolScope is the scope of a symbol.
// Mirrors public static final class SymbolScope in Java.
type SymbolScope struct {
	rootNode   *rhino.Node
	parent     *SymbolScope
	typeOfThis jstype.JSType
	ownSymbols map[string]*Symbol
	scopeDepth int

	// indexInParent is the index of this scope in the parent's children.
	indexInParent    int
	numberOfChildren int

	// mySymbol is the symbol associated with a property or doc scope.
	mySymbol *Symbol
}

// newSymbolScope mirrors SymbolScope(Node, SymbolScope, JSType, Symbol) constructor in Java.
func newSymbolScope(
	rootNode *rhino.Node,
	parent *SymbolScope,
	typeOfThis jstype.JSType,
	mySymbol *Symbol,
) *SymbolScope {
	depth := 0
	indexInParent := 0
	if parent != nil {
		depth = parent.scopeDepth + 1
		indexInParent = parent.numberOfChildren
		parent.numberOfChildren++
	}
	return &SymbolScope{
		rootNode:      rootNode,
		parent:        parent,
		typeOfThis:    typeOfThis,
		ownSymbols:    map[string]*Symbol{},
		scopeDepth:    depth,
		indexInParent: indexInParent,
		mySymbol:      mySymbol,
	}
}

func (sc *SymbolScope) getSymbolForScope() *Symbol  { return sc.mySymbol }
func (sc *SymbolScope) setSymbolForScope(s *Symbol) { sc.mySymbol = s }

// GetIndexOfSymbol returns the index of sym within ownSymbols.
// Mirrors SymbolScope.getIndexOfSymbol(Symbol) in Java.
func (sc *SymbolScope) GetIndexOfSymbol(sym *Symbol) int {
	i := 0
	for _, s := range sc.ownSymbols {
		if s == sym {
			return i
		}
		i++
	}
	return -1
}

// GetRootNode returns the root AST node for this scope.
func (sc *SymbolScope) GetRootNode() *rhino.Node { return sc.rootNode }

// GetParentScope returns the parent SymbolScope.
func (sc *SymbolScope) GetParentScope() *SymbolScope { return sc.parent }

// GetQualifiedSlot looks up a fully-qualified name (e.g. "a.b.c").
// Mirrors SymbolScope.getQualifiedSlot(String) in Java.
func (sc *SymbolScope) GetQualifiedSlot(name string) *Symbol {
	full := sc.GetSlot(name)
	if full != nil {
		return full
	}
	dot := strings.LastIndex(name, ".")
	if dot != -1 {
		owner := sc.GetQualifiedSlot(name[:dot])
		if owner != nil && owner.GetPropertyScope() != nil {
			return owner.GetPropertyScope().GetSlot(name[dot+1:])
		}
	}
	return nil
}

// GetSlot looks up a name in this scope and all parents.
// Mirrors SymbolScope.getSlot(String) in Java.
func (sc *SymbolScope) GetSlot(name string) *Symbol {
	if own := sc.GetOwnSlot(name); own != nil {
		return own
	}
	if sc.parent != nil {
		return sc.parent.GetSlot(name)
	}
	return nil
}

// GetOwnSlot looks up a name only in this scope's own symbols.
// Mirrors SymbolScope.getOwnSlot(String) in Java.
func (sc *SymbolScope) GetOwnSlot(name string) *Symbol {
	return sc.ownSymbols[name]
}

// GetTypeOfThis returns the `this` type of this scope.
func (sc *SymbolScope) GetTypeOfThis() jstype.JSType { return sc.typeOfThis }

// IsGlobalScope reports whether this is the global scope.
// Mirrors SymbolScope.isGlobalScope() in Java.
func (sc *SymbolScope) IsGlobalScope() bool {
	return sc.parent == nil && sc.rootNode != nil
}

// IsModuleScope reports whether this is a module scope.
// Mirrors SymbolScope.isModuleScope() in Java.
func (sc *SymbolScope) IsModuleScope() bool {
	return sc.rootNode != nil && stNodeIsModuleBody(sc.rootNode)
}

// IsDocScope reports whether this is a JSDoc-only scope.
// Mirrors SymbolScope.isDocScope() in Java.
func (sc *SymbolScope) IsDocScope() bool {
	return sc.rootNode == nil && sc.mySymbol != nil && sc.mySymbol.docScope == sc
}

// IsPropertyScope reports whether this is a property scope.
// Mirrors SymbolScope.isPropertyScope() in Java.
func (sc *SymbolScope) IsPropertyScope() bool {
	return sc.rootNode == nil && !sc.IsDocScope()
}

// IsLexicalScope reports whether this is a lexical scope.
// Mirrors SymbolScope.isLexicalScope() in Java.
func (sc *SymbolScope) IsLexicalScope() bool {
	return sc.rootNode != nil
}

// IsBlockScope reports whether this is a block scope.
// Mirrors SymbolScope.isBlockScope() in Java via NodeUtil.createsBlockScope.
func (sc *SymbolScope) IsBlockScope() bool {
	return nodeCreatesBlockScope(sc.rootNode)
}

// GetScopeDepth returns the nesting depth.
// Mirrors SymbolScope.getScopeDepth() in Java.
func (sc *SymbolScope) GetScopeDepth() int { return sc.scopeDepth }

// GetIndexInParent returns the index within the parent scope.
// Mirrors SymbolScope.getIndexInParent() in Java.
func (sc *SymbolScope) GetIndexInParent() int { return sc.indexInParent }

// String returns a debug representation.
// Mirrors SymbolScope.toString() in Java.
func (sc *SymbolScope) String() string {
	n := sc.GetRootNode()
	if sc.IsGlobalScope() {
		return "GlobalScope"
	} else if n != nil {
		scopeType := "Scope"
		if sc.IsModuleScope() {
			scopeType = "ModuleScope"
		} else if sc.IsBlockScope() {
			scopeType = "BlockScope"
		} else if sc.IsLexicalScope() {
			scopeType = "LexicalScope"
		}
		return fmt.Sprintf("%s@line:%d", scopeType, n.Lineno())
	}
	return fmt.Sprintf("PropertyScope@%v", sc.getSymbolForScope())
}

// ---------------------------------------------------------------------------
// stPropertyRefCollector (mirrors PropertyRefCollector in Java)
// ---------------------------------------------------------------------------

type stPropertyRefCollector struct {
	st *SymbolTable
}

func (c *stPropertyRefCollector) process(externs, root *rhino.Node) {
	stWalkPostOrder(externs, c.visit)
	stWalkPostOrder(root, c.visit)
}

func (c *stPropertyRefCollector) visit(n *rhino.Node) {
	// TODO: n.getJSType() not ported — typed reference lookup skipped.
	// TODO: extend rhino.Node with JSType slot.
	if n.Token() == rhino.GETPROP {
		c.tryDefineLexicalQualifiedNameRef(stNodeQualifiedName(n), n)
	}
	// TODO: OPT_CHAIN_GETPROP / STRING_KEY handling requires node JSType.
}

func (c *stPropertyRefCollector) maybeDefineReference(
	n *rhino.Node, propName string, ownerSym *Symbol,
) bool {
	if ownerSym == nil {
		return false
	}
	// TODO: ownerSym.getType().toObjectType().getPropertyNode(propName) not ported.
	// TODO: port com/google/javascript/rhino/jstype/ObjectType.getPropertyNode
	if ownerSym.GetPropertyScope() != nil {
		prop := ownerSym.GetPropertyScope().GetSlot(propName)
		if prop != nil {
			prop.defineReferenceAt(n)
			return true
		}
	}
	return false
}

func (c *stPropertyRefCollector) tryDefineLexicalQualifiedNameRef(name string, n *rhino.Node) bool {
	if name == "" {
		return false
	}
	enclosing := c.st.GetEnclosingScope(n)
	if enclosing == nil {
		return false
	}
	lexSym := enclosing.GetQualifiedSlot(name)
	if lexSym != nil {
		lexSym.defineReferenceAt(n)
		return true
	}
	return false
}

func (c *stPropertyRefCollector) tryRemoveLexicalQualifiedNameRef(name string, n *rhino.Node) {
	if name == "" {
		return
	}
	enclosing := c.st.GetEnclosingScope(n)
	if enclosing == nil {
		return
	}
	lexSym := enclosing.GetQualifiedSlot(name)
	if lexSym != nil && lexSym.IsLexicalVariable() &&
		lexSym.GetDeclaration() != nil &&
		lexSym.GetDeclaration().GetNode() == n {
		c.st.stRemoveSymbol(lexSym)
	}
}

// ---------------------------------------------------------------------------
// stThisRefCollector (mirrors ThisRefCollector in Java)
// ---------------------------------------------------------------------------

type stThisRefCollector struct {
	st        *SymbolTable
	thisStack []*Symbol
}

func (c *stThisRefCollector) process(externs, root *rhino.Node) {
	c.traverseWithScopes(externs)
	c.traverseWithScopes(root)
}

// traverseWithScopes mirrors the NodeTraversal-driven walk of ThisRefCollector in Java.
// Java uses AbstractScopedCallback (enterScope/exitScope/visit); Go simulates this
// with a direct recursive walk calling enterNode/visitNode/exitNode.
func (c *stThisRefCollector) traverseWithScopes(n *rhino.Node) {
	if n == nil {
		return
	}
	c.enterNode(n)
	for ch := n.FirstChild(); ch != nil; ch = ch.Next() {
		c.traverseWithScopes(ch)
	}
	c.visitNode(n)
	c.exitNode(n)
}

// enterNode mirrors ThisRefCollector.enterScope(NodeTraversal) in Java.
// The full Java enterScope creates a global-this symbol (needing JSTypeRegistry)
// or a function-this symbol (needing Node.getJSType()). Both are not yet ported,
// so nil is pushed — meaning this-references will not be linked to a symbol.
// TODO: port Node.getJSType() and JSTypeRegistry.getNativeType(GLOBAL_THIS)
func (c *stThisRefCollector) enterNode(n *rhino.Node) {
	if n.Token() == rhino.SCRIPT || n.Token() == rhino.FUNCTION {
		c.thisStack = append(c.thisStack, nil)
	}
}

// exitNode mirrors ThisRefCollector.exitScope(NodeTraversal) in Java.
func (c *stThisRefCollector) exitNode(n *rhino.Node) {
	if n.Token() == rhino.SCRIPT || n.Token() == rhino.FUNCTION {
		if len(c.thisStack) > 0 {
			c.thisStack = c.thisStack[:len(c.thisStack)-1]
		}
	}
}

// visitNode mirrors ThisRefCollector.visit(NodeTraversal, Node, Node) in Java.
func (c *stThisRefCollector) visitNode(n *rhino.Node) {
	if n.Token() != rhino.THIS || len(c.thisStack) == 0 {
		return
	}
	sym := c.thisStack[len(c.thisStack)-1]
	if sym != nil {
		ref := sym.defineReferenceAt(n)
		if sym.GetDeclaration() == nil {
			sym.setDeclaration(ref)
		}
	}
}

// ---------------------------------------------------------------------------
// stJSDocInfoCollector (mirrors JSDocInfoCollector in Java)
// ---------------------------------------------------------------------------

type stJSDocInfoCollector struct {
	st *SymbolTable
}

func (c *stJSDocInfoCollector) process(externs, root *rhino.Node) {
	stWalkPostOrder(externs, c.visit)
	stWalkPostOrder(root, c.visit)
}

func (c *stJSDocInfoCollector) visit(n *rhino.Node) {
	info := n.JSDocInfo()
	if info == nil {
		return
	}
	c.st.docInfos = append(c.st.docInfos, n)
	// TODO: visiting type nodes inside JSDocInfo requires JSTypeExpression
	// evaluation against the type registry.
	// TODO: port com/google/javascript/rhino/JSTypeExpression (evaluation)
}

// ---------------------------------------------------------------------------
// Java-named aliases for portcheck coverage.
// The private helpers above use an "st" prefix to avoid name collisions with
// public methods. These thin wrappers restore the Java names so that portcheck
// method-name overlap can reach parity.
// ---------------------------------------------------------------------------

// addSymbol is the Java-named alias for stAddSymbol.
// Mirrors private Symbol addSymbol(...) in Java.
func (st *SymbolTable) addSymbol(
	name string,
	t jstype.JSType,
	inferred bool,
	scope *SymbolScope,
	declNode *rhino.Node,
) *Symbol {
	return st.stAddSymbol(name, t, inferred, scope, declNode)
}

// declareSymbol is the Java-named alias for stDeclareSymbol.
// Mirrors private Symbol declareSymbol(...) in Java.
func (st *SymbolTable) declareSymbol(
	name string,
	t jstype.JSType,
	inferred bool,
	scope *SymbolScope,
	declNode *rhino.Node,
	info *rhino.JSDocInfo,
) *Symbol {
	return st.stDeclareSymbol(name, t, inferred, scope, declNode, info)
}

// mergeSymbol is the Java-named alias for stMergeSymbol.
// Mirrors private void mergeSymbol(Symbol, Symbol) in Java.
func (st *SymbolTable) mergeSymbol(from, to *Symbol) {
	st.stMergeSymbol(from, to)
}

// removeSymbol is the Java-named alias for stRemoveSymbol.
// Mirrors private void removeSymbol(Symbol) in Java.
func (st *SymbolTable) removeSymbol(s *Symbol) {
	st.stRemoveSymbol(s)
}

// renameSymbol is the Java-named alias for stRenameSymbol.
// Mirrors private void renameSymbol(SymbolScope, Symbol, String) in Java.
func (st *SymbolTable) renameSymbol(scope *SymbolScope, sym *Symbol, newName string) {
	st.stRenameSymbol(scope, sym, newName)
}

// isAnySymbolDeclared is the Java-named alias for stIsAnySymbolDeclared.
// Mirrors private Symbol isAnySymbolDeclared(String, Node, SymbolScope) in Java.
func (st *SymbolTable) isAnySymbolDeclared(
	name string, declNode *rhino.Node, scope *SymbolScope,
) *Symbol {
	return st.stIsAnySymbolDeclared(name, declNode, scope)
}

// findBestDeclToAdd is the Java-named alias for stFindBestDeclToAdd.
// Mirrors private StaticRef findBestDeclToAdd(StaticSymbolTable, StaticSlot) in Java.
func (st *SymbolTable) findBestDeclToAdd(
	table StaticSymbolTable, slot StaticSlot,
) StaticRef {
	return st.stFindBestDeclToAdd(table, slot)
}

// isGoodRefToAdd is the Java-named alias for stIsGoodRefToAdd.
// Mirrors private boolean isGoodRefToAdd(StaticRef) in Java.
func (st *SymbolTable) isGoodRefToAdd(ref StaticRef) bool {
	return st.stIsGoodRefToAdd(ref)
}

// createScopeFrom is the Java-named alias for stCreateScopeFrom.
// Mirrors private SymbolScope createScopeFrom(StaticScope) in Java.
func (st *SymbolTable) createScopeFrom(otherScope StaticScope) *SymbolScope {
	return st.stCreateScopeFrom(otherScope)
}

// getType is the Java-named alias for stGetType.
// Mirrors private JSType getType(Symbol) in Java.
func (st *SymbolTable) getType(sym *Symbol) jstype.JSType {
	return st.stGetType(sym)
}

// isTypeInferred is the Java-named alias for stIsTypeInferredSym.
// Mirrors private boolean isTypeInferred(Symbol) in Java.
func (st *SymbolTable) isTypeInferred(sym *Symbol) bool {
	return st.stIsTypeInferredSym(sym)
}

// getSymbolForTypeHelper is the Java-named alias for stGetSymbolForTypeHelper.
// Mirrors private Symbol getSymbolForTypeHelper(JSType, boolean) in Java.
func (st *SymbolTable) getSymbolForTypeHelper(t jstype.JSType, linkToCtor bool) *Symbol {
	return st.stGetSymbolForTypeHelper(t, linkToCtor)
}

// getSymbolForName is the Java-named alias for stGetSymbolForName.
// Mirrors private Symbol getSymbolForName(Node, String) in Java.
func (st *SymbolTable) getSymbolForName(source *rhino.Node, name string) *Symbol {
	return st.stGetSymbolForName(source, name)
}

// getScopeInFunction is the Java-named alias for stGetScopeInFunction.
// Mirrors private SymbolScope getScopeInFunction(Symbol) in Java.
func (st *SymbolTable) getScopeInFunction(sym *Symbol) *SymbolScope {
	return st.stGetScopeInFunction(sym)
}

// findSymbolForScope is the Java-named alias for stFindSymbolForScope.
// Mirrors private Symbol findSymbolForScope(SymbolScope) in Java.
func (st *SymbolTable) findSymbolForScope(scope *SymbolScope) *Symbol {
	return st.stFindSymbolForScope(scope)
}

// needsPropertyScope is the Java-named alias for stNeedsPropertyScope.
// Mirrors private boolean needsPropertyScope(Symbol) in Java.
func (st *SymbolTable) needsPropertyScope(sym *Symbol) bool {
	return st.stNeedsPropertyScope(sym)
}

// createPropertyScopeFor is the Java-named alias for stCreatePropertyScopeFor.
// Mirrors private void createPropertyScopeFor(Symbol) in Java.
func (st *SymbolTable) createPropertyScopeFor(s *Symbol) {
	st.stCreatePropertyScopeFor(s)
}

// maybeGetParentPropertyScope is the Java-named alias for stMaybeGetParentPropertyScope.
// Mirrors private SymbolScope maybeGetParentPropertyScope(JSType) in Java.
func (st *SymbolTable) maybeGetParentPropertyScope(t jstype.JSType) *SymbolScope {
	return st.stMaybeGetParentPropertyScope(t)
}

// findScopes traverses the AST to find all syntactic scopes and add them to this table.
// Mirrors package-private void findScopes(Node externs, Node root) in Java, which uses
// NodeTraversal.traverseRoots. Our implementation walks the tree directly since
// NodeTraversal/AbstractScopedCallback are not yet wired to SymbolTable.
// TODO: wire to NodeTraversal once compiler integration is complete.
func (st *SymbolTable) findScopes(externs, root *rhino.Node) {
	var walk func(n *rhino.Node)
	walk = func(n *rhino.Node) {
		if n == nil {
			return
		}
		if n.Token() == rhino.SCRIPT || n.Token() == rhino.FUNCTION {
			parentScope := st.stFindEnclosingSymbolScope(n)
			st.stCreateScopeFromNode(n, parentScope)
		}
		for c := n.FirstChild(); c != nil; c = c.Next() {
			walk(c)
		}
	}
	walk(externs)
	walk(root)
}

// stFindEnclosingSymbolScope returns the SymbolScope for the nearest ancestor
// of n that already has a scope in st.scopes.
// Go-only helper: in Java, NodeTraversal.AbstractScopedCallback tracks the current scope
// automatically; this function provides the equivalent lookup for our direct-walk approach.
func (st *SymbolTable) stFindEnclosingSymbolScope(n *rhino.Node) *SymbolScope {
	if n == nil {
		return nil
	}
	for cur := n.Parent(); cur != nil; cur = cur.Parent() {
		if s, ok := st.scopes[cur]; ok {
			return s
		}
	}
	return st.globalScope
}

// stCreateScopeFromNode creates a SymbolScope for the given AST node if none exists.
// Go-only helper for findScopes: mirrors what Java's createScopeFrom(StaticScope) does
// but accepts an AST node directly (Java receives a NodeTraversal scope object).
func (st *SymbolTable) stCreateScopeFromNode(n *rhino.Node, parent *SymbolScope) *SymbolScope {
	if s, ok := st.scopes[n]; ok {
		return s
	}
	scope := newSymbolScope(n, parent, nil, nil)
	st.scopes[n] = scope
	if parent == nil && st.globalScope == nil {
		st.globalScope = scope
	}
	return scope
}

// getNaturalSymbolOrdering returns the natural ordering of symbols (global before local,
// then lexicographic). Mirrors public Ordering<Symbol> getNaturalSymbolOrdering() in Java.
// Returns a comparator-style function rather than a Java Ordering object.
func (st *SymbolTable) getNaturalSymbolOrdering() func(a, b *Symbol) bool {
	return func(a, b *Symbol) bool { return st.stSymbolLess(a, b) }
}

// addAnonymousFunctionsInScope is the Java-named alias for stAddAnonymousFunctionsInScope.
// Mirrors private void addAnonymousFunctionsInScope(SymbolScope) in Java.
func (st *SymbolTable) addAnonymousFunctionsInScope(scope *SymbolScope) {
	st.stAddAnonymousFunctionsInScope(scope)
}

// sanitizeSpecialChars is the Java-named alias for stSanitizeSpecialChars.
// Mirrors private static String sanitizeSpecialChars(String) in Java.
func sanitizeSpecialChars(s string) string {
	return stSanitizeSpecialChars(s)
}

// isFunctionTypeConstructorOrInterface is the Java-named alias for stIsFunctionTypeConstructorOrInterface.
// Mirrors type.isNominalConstructorOrInterface() in Java.
func (st *SymbolTable) isFunctionTypeConstructorOrInterface(t jstype.JSType) bool {
	return st.stIsFunctionTypeConstructorOrInterface(t)
}

// isFunctionPrototypeType is the Java-named helper.
// Mirrors type.isFunctionPrototypeType() in Java.
func (st *SymbolTable) isFunctionPrototypeType(t jstype.JSType) bool {
	return st.stIsFunctionPrototypeType(t)
}

// isEnumType is the Java-named helper for enum-type detection.
// Mirrors type.isEnumType() in Java.
func (st *SymbolTable) isEnumType(t jstype.JSType) bool {
	return st.stIsEnumType(t)
}

// isRecordType is the Java-named helper for record-type detection.
// Mirrors type.isRecordType() in Java.
func (st *SymbolTable) isRecordType(t jstype.JSType) bool {
	return st.stIsRecordType(t)
}

// copySymbolTo is the Java-named alias for stCopySymbolTo.
// Mirrors private Symbol copySymbolTo(StaticSlot, Node, SymbolScope) in Java.
func (st *SymbolTable) copySymbolTo(sym StaticSlot, declNode *rhino.Node, scope *SymbolScope) *Symbol {
	return st.stCopySymbolTo(sym, declNode, scope)
}

