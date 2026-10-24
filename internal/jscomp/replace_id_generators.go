// Port of com/google/javascript/jscomp/ReplaceIdGenerators.java @ closure-compiler v20261021

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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// ReplaceIdGenerators replaces calls to id generators with ids.
//
// Use this to get unique and short ids.
//
// Go-only: language-adapter — Java uses null as a map key for a generator
// whose name is unknown (JSDoc @idGenerator on a node other than an
// assignment to a qualified name, a name declaration or a function, or a
// call whose target has no qualified name). Such keys are the string "null"
// here: no qualified name can be "null" (it is a keyword), so lookups match
// Java's, and IdMappingUtil serializes Java's null key as "null" too.
type ReplaceIdGenerators struct {
	compiler                      AbstractCompiler
	templateLiteralsAreTranspiled bool
	nameGenerators                *jutil.LinkedHashMap[string, ReplaceIdGenerators_NameSupplier]
	consistNameMap                *jutil.LinkedHashMap[string, *jutil.LinkedHashMap[string, string]]

	idGeneratorMaps *jutil.LinkedHashMap[string, *jutil.LinkedHashMap[string, string]]
	previousMap     *jutil.LinkedHashMap[string, *jutil.BiMap[string, string]]

	generatePseudoNames bool
	xidHashFunction     Xid_HashFunction
}

// Java: ReplaceIdGenerators.NON_GLOBAL_ID_GENERATOR_CALL
var ReplaceIdGenerators_NON_GLOBAL_ID_GENERATOR_CALL = DiagnosticType_Error(
	"JSC_NON_GLOBAL_ID_GENERATOR_CALL", "Id generator call must be in the global scope")

// Java: ReplaceIdGenerators.CONDITIONAL_ID_GENERATOR_CALL
var ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL = DiagnosticType_Error(
	"JSC_CONDITIONAL_ID_GENERATOR_CALL", "Id generator call must be unconditional")

// Java: ReplaceIdGenerators.MISSING_NAME_MAP_FOR_GENERATOR
var ReplaceIdGenerators_MISSING_NAME_MAP_FOR_GENERATOR = DiagnosticType_Warning(
	"JSC_MISSING_NAME_MAP_FOR_GENERATOR",
	"The mapped id generator, does not have a renaming map supplied.")

// Java: ReplaceIdGenerators.INVALID_GENERATOR_PARAMETER
var ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER = DiagnosticType_Warning(
	"JSC_INVALID_GENERATOR_PARAMETER", "An id generator must be called with a literal.")

// Java: ReplaceIdGenerators.INVALID_TEMPLATE_LITERAL_PARAMETER
var ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER = DiagnosticType_Warning(
	"JSC_INVALID_GENERATOR_PARAMETER",
	"An id generator must be called with a template literals with no invalid escape"+
		" sequences.")

// Java: ReplaceIdGenerators.SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN
var ReplaceIdGenerators_SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN = DiagnosticType_Error(
	"JSC_SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN",
	"Object literal shorthand functions is not allowed in the "+
		"arguments of an id generator")

// Java: ReplaceIdGenerators.COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN
var ReplaceIdGenerators_COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN = DiagnosticType_Error(
	"JSC_COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN",
	"Object literal computed property name is not allowed in the "+
		"arguments of an id generator")

// replaceIdGenerators_nullableName is the map key for a nullable Java String
// name (see the type comment).
//
// Go-only: language-adapter — Java null map keys.
func replaceIdGenerators_nullableName(name *string) string {
	if name == nil {
		return "null"
	}
	return *name
}

// Java: ReplaceIdGenerators(AbstractCompiler, boolean, Map<String, RenamingMap>, boolean, String, Xid.HashFunction)
func NewReplaceIdGenerators(
	compiler AbstractCompiler,
	templateLiteralsAreTranspiled bool,
	idGens *jutil.LinkedHashMap[string, RenamingMap],
	generatePseudoNames bool,
	previousMapSerialized *string,
	xidHashFunction Xid_HashFunction) *ReplaceIdGenerators {
	r := &ReplaceIdGenerators{}
	r.compiler = compiler
	r.templateLiteralsAreTranspiled = templateLiteralsAreTranspiled
	r.generatePseudoNames = generatePseudoNames
	r.xidHashFunction = xidHashFunction
	r.nameGenerators = jutil.NewLinkedHashMap[string, ReplaceIdGenerators_NameSupplier]()
	r.idGeneratorMaps = jutil.NewLinkedHashMap[string, *jutil.LinkedHashMap[string, string]]()
	r.consistNameMap = jutil.NewLinkedHashMap[string, *jutil.LinkedHashMap[string, string]]()

	previousMap := IdMappingUtil_ParseSerializedIdMappings(previousMapSerialized)
	r.previousMap = previousMap

	if idGens != nil {
		for name, map_ := range idGens.All() {
			if renamingToken, ok := map_.(RenamingToken); ok {
				switch renamingToken {
				case RenamingToken_DISABLE:
					r.nameGenerators.Put(name, nil)
					continue
					// don't put an entry in idGeneratorsMap
				case RenamingToken_INCONSISTENT:
					r.nameGenerators.Put(
						name, r.CreateNameSupplier(ReplaceIdGenerators_RenameStrategy_INCONSISTENT, previousMap.Get(name)))
				case RenamingToken_STABLE:
					r.nameGenerators.Put(
						name, r.CreateNameSupplier(ReplaceIdGenerators_RenameStrategy_STABLE, previousMap.Get(name)))
				}
			} else {
				r.nameGenerators.Put(name, ReplaceIdGenerators_CreateNameSupplier(ReplaceIdGenerators_RenameStrategy_MAPPED, map_))
			}
			r.idGeneratorMaps.Put(name, jutil.NewLinkedHashMap[string, string]())
		}
	}
	return r
}

// ReplaceIdGenerators_RenameStrategy is the Java enum
// ReplaceIdGenerators.RenameStrategy.
type ReplaceIdGenerators_RenameStrategy int32

const (
	ReplaceIdGenerators_RenameStrategy_CONSISTENT ReplaceIdGenerators_RenameStrategy = iota
	ReplaceIdGenerators_RenameStrategy_INCONSISTENT
	ReplaceIdGenerators_RenameStrategy_MAPPED
	ReplaceIdGenerators_RenameStrategy_STABLE
	ReplaceIdGenerators_RenameStrategy_XID
)

var replaceIdGenerators_RenameStrategy_names = [...]string{"CONSISTENT", "INCONSISTENT", "MAPPED", "STABLE", "XID"}

// ReplaceIdGenerators_RenameStrategy_values is RenameStrategy.values().
func ReplaceIdGenerators_RenameStrategy_values() []ReplaceIdGenerators_RenameStrategy {
	return []ReplaceIdGenerators_RenameStrategy{
		ReplaceIdGenerators_RenameStrategy_CONSISTENT,
		ReplaceIdGenerators_RenameStrategy_INCONSISTENT,
		ReplaceIdGenerators_RenameStrategy_MAPPED,
		ReplaceIdGenerators_RenameStrategy_STABLE,
		ReplaceIdGenerators_RenameStrategy_XID,
	}
}

// Name is Enum.name().
func (s ReplaceIdGenerators_RenameStrategy) Name() string {
	return replaceIdGenerators_RenameStrategy_names[s]
}

// String is Enum.toString().
func (s ReplaceIdGenerators_RenameStrategy) String() string { return s.Name() }

// Ordinal is Enum.ordinal().
func (s ReplaceIdGenerators_RenameStrategy) Ordinal() int { return int(s) }

// ReplaceIdGenerators_NameSupplier is the private interface NameSupplier.
type ReplaceIdGenerators_NameSupplier interface {
	// Java: NameSupplier.getName(String, String)
	GetName(id string, name string) string

	// Java: NameSupplier.getRenameStrategy()
	GetRenameStrategy() ReplaceIdGenerators_RenameStrategy
}

// ReplaceIdGenerators_ObfuscatedNameSupplier is the private class
// ObfuscatedNameSupplier.
type ReplaceIdGenerators_ObfuscatedNameSupplier struct {
	generator        NameGenerator
	previousMappings *jutil.BiMap[string, string]
	renameStrategy   ReplaceIdGenerators_RenameStrategy
}

// Java: ObfuscatedNameSupplier(RenameStrategy, BiMap<String, String>)
func NewReplaceIdGenerators_ObfuscatedNameSupplier(
	renameStrategy ReplaceIdGenerators_RenameStrategy, previousMappings *jutil.BiMap[string, string]) *ReplaceIdGenerators_ObfuscatedNameSupplier {
	s := &ReplaceIdGenerators_ObfuscatedNameSupplier{}
	s.previousMappings = previousMappings.Inverse()
	// Go-only: language-adapter — Java passes the live keySet() view;
	// previousMappings is an immutable BiMap, so a copy of its keys is the
	// same set.
	reservedNames := jutil.NewLinkedHashSet[string]()
	for k := range previousMappings.Keys() {
		reservedNames.Add(k)
	}
	s.generator = NewDefaultNameGeneratorWithReserved(reservedNames, "", jutil.NewLinkedHashSet[uint16]())
	s.renameStrategy = renameStrategy
	return s
}

// Java: ObfuscatedNameSupplier.getName(String, String)
func (s *ReplaceIdGenerators_ObfuscatedNameSupplier) GetName(id string, name string) string {
	newName, ok := s.previousMappings.Lookup(id)
	if !ok {
		newName = s.generator.GenerateNextName()
	}
	return newName
}

// Java: ObfuscatedNameSupplier.getRenameStrategy()
func (s *ReplaceIdGenerators_ObfuscatedNameSupplier) GetRenameStrategy() ReplaceIdGenerators_RenameStrategy {
	return s.renameStrategy
}

// ReplaceIdGenerators_PseudoNameSupplier is the private class
// PseudoNameSupplier.
type ReplaceIdGenerators_PseudoNameSupplier struct {
	counter        int32
	renameStrategy ReplaceIdGenerators_RenameStrategy
}

// Java: PseudoNameSupplier(RenameStrategy)
func NewReplaceIdGenerators_PseudoNameSupplier(renameStrategy ReplaceIdGenerators_RenameStrategy) *ReplaceIdGenerators_PseudoNameSupplier {
	s := &ReplaceIdGenerators_PseudoNameSupplier{counter: 0}
	s.renameStrategy = renameStrategy
	return s
}

// Java: PseudoNameSupplier.getName(String, String)
func (s *ReplaceIdGenerators_PseudoNameSupplier) GetName(id string, name string) string {
	if s.renameStrategy == ReplaceIdGenerators_RenameStrategy_INCONSISTENT {
		c := s.counter
		s.counter++
		return name + "$" + strconv.Itoa(int(c))
	}
	return name + "$0"
}

// Java: PseudoNameSupplier.getRenameStrategy()
func (s *ReplaceIdGenerators_PseudoNameSupplier) GetRenameStrategy() ReplaceIdGenerators_RenameStrategy {
	return s.renameStrategy
}

// ReplaceIdGenerators_StableNameSupplier is the private class
// StableNameSupplier.
type ReplaceIdGenerators_StableNameSupplier struct{}

// Java: StableNameSupplier()
func NewReplaceIdGenerators_StableNameSupplier() *ReplaceIdGenerators_StableNameSupplier {
	return &ReplaceIdGenerators_StableNameSupplier{}
}

// Java: StableNameSupplier.getName(String, String)
func (s *ReplaceIdGenerators_StableNameSupplier) GetName(id string, name string) string {
	return sourcemap.Base64_Base64EncodeInt(jutil.StringHashCode(name))
}

// Java: StableNameSupplier.getRenameStrategy()
func (s *ReplaceIdGenerators_StableNameSupplier) GetRenameStrategy() ReplaceIdGenerators_RenameStrategy {
	return ReplaceIdGenerators_RenameStrategy_STABLE
}

// ReplaceIdGenerators_XidNameSupplier is the private class XidNameSupplier.
type ReplaceIdGenerators_XidNameSupplier struct {
	xid *Xid
}

// Java: XidNameSupplier(Xid.HashFunction)
func NewReplaceIdGenerators_XidNameSupplier(hashFunction Xid_HashFunction) *ReplaceIdGenerators_XidNameSupplier {
	s := &ReplaceIdGenerators_XidNameSupplier{}
	if hashFunction == nil {
		s.xid = NewXid()
	} else {
		s.xid = NewXidWithHasher(hashFunction)
	}
	return s
}

// Java: XidNameSupplier.getName(String, String)
func (s *ReplaceIdGenerators_XidNameSupplier) GetName(id string, name string) string {
	return s.xid.Get(name)
}

// Java: XidNameSupplier.getRenameStrategy()
func (s *ReplaceIdGenerators_XidNameSupplier) GetRenameStrategy() ReplaceIdGenerators_RenameStrategy {
	return ReplaceIdGenerators_RenameStrategy_XID
}

// ReplaceIdGenerators_MappedNameSupplier is the private class
// MappedNameSupplier.
type ReplaceIdGenerators_MappedNameSupplier struct {
	map_ RenamingMap
}

// Java: MappedNameSupplier(RenamingMap)
func NewReplaceIdGenerators_MappedNameSupplier(map_ RenamingMap) *ReplaceIdGenerators_MappedNameSupplier {
	return &ReplaceIdGenerators_MappedNameSupplier{map_: map_}
}

// Java: MappedNameSupplier.getName(String, String)
func (s *ReplaceIdGenerators_MappedNameSupplier) GetName(id string, name string) string {
	// Go-only: language-adapter — a null replacement makes every caller of
	// getObfuscatedName throw in Java (IR.string(null) and
	// Node.setString(null) reject it); dereferencing it panics here.
	return *s.map_.Get(name)
}

// Java: MappedNameSupplier.getRenameStrategy()
func (s *ReplaceIdGenerators_MappedNameSupplier) GetRenameStrategy() ReplaceIdGenerators_RenameStrategy {
	return ReplaceIdGenerators_RenameStrategy_MAPPED
}

// Java: ReplaceIdGenerators.createNameSupplier(RenameStrategy, BiMap<String, String>)
func (r *ReplaceIdGenerators) CreateNameSupplier(
	renameStrategy ReplaceIdGenerators_RenameStrategy, previousMappings *jutil.BiMap[string, string]) ReplaceIdGenerators_NameSupplier {
	if previousMappings == nil {
		previousMappings = jutil.ImmutableBiMap_CopyOf(jutil.NewLinkedHashMap[string, string]())
	}
	if renameStrategy == ReplaceIdGenerators_RenameStrategy_STABLE {
		return NewReplaceIdGenerators_StableNameSupplier()
	} else if renameStrategy == ReplaceIdGenerators_RenameStrategy_XID {
		return NewReplaceIdGenerators_XidNameSupplier(r.xidHashFunction)
	} else if r.generatePseudoNames {
		return NewReplaceIdGenerators_PseudoNameSupplier(renameStrategy)
	} else {
		return NewReplaceIdGenerators_ObfuscatedNameSupplier(renameStrategy, previousMappings)
	}
}

// Java: ReplaceIdGenerators.createNameSupplier(RenameStrategy, RenamingMap)
func ReplaceIdGenerators_CreateNameSupplier(
	renameStrategy ReplaceIdGenerators_RenameStrategy, mappings RenamingMap) ReplaceIdGenerators_NameSupplier {
	jutil.CheckState(renameStrategy == ReplaceIdGenerators_RenameStrategy_MAPPED)
	return NewReplaceIdGenerators_MappedNameSupplier(mappings)
}

// Java: ReplaceIdGenerators.CREATE_TEMPLATE_TAG_FIRST_ARG
var ReplaceIdGenerators_CREATE_TEMPLATE_TAG_FIRST_ARG = rhino.QualifiedName_Of("$jscomp.createTemplateTagFirstArg")

// ReplaceIdGenerators_GatherGenerators is the inner class GatherGenerators.
type ReplaceIdGenerators_GatherGenerators struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	this_ *ReplaceIdGenerators // the enclosing instance
}

// Java: GatherGenerators()
func NewReplaceIdGenerators_GatherGenerators(outer *ReplaceIdGenerators) *ReplaceIdGenerators_GatherGenerators {
	return &ReplaceIdGenerators_GatherGenerators{this_: outer}
}

// Java: GatherGenerators.visit(NodeTraversal, Node, Node)
func (g *ReplaceIdGenerators_GatherGenerators) Visit(unused *NodeTraversal, n rhino.Node, parent rhino.Node) {
	r := g.this_
	doc := n.GetJSDocInfo()
	if doc == nil || !doc.IsAnyIdGenerator() {
		return
	}

	var name *string = nil
	if n.IsAssign() {
		name = n.GetFirstChild().GetQualifiedName()
	} else if NodeUtil_IsNameDeclaration(n) {
		s := n.GetFirstChild().GetString()
		name = &s
	} else if n.IsFunction() {
		s := n.GetFirstChild().GetString()
		name = &s
		if *name == "" {
			return
		}
	}
	key := replaceIdGenerators_nullableName(name)
	if r.nameGenerators.ContainsKey(key) {
		// This generator is already registered from our constructor.
		// Don't override it.
		return
	}
	if doc.IsConsistentIdGenerator() {
		r.consistNameMap.Put(key, jutil.NewLinkedHashMap[string, string]())
		r.nameGenerators.Put(
			key, r.CreateNameSupplier(ReplaceIdGenerators_RenameStrategy_CONSISTENT, r.previousMap.Get(key)))
	} else if doc.IsStableIdGenerator() {
		r.nameGenerators.Put(key, r.CreateNameSupplier(ReplaceIdGenerators_RenameStrategy_STABLE, r.previousMap.Get(key)))
	} else if doc.IsXidGenerator() {
		r.nameGenerators.Put(key, r.CreateNameSupplier(ReplaceIdGenerators_RenameStrategy_XID, r.previousMap.Get(key)))
	} else if doc.IsIdGenerator() {
		r.nameGenerators.Put(
			key, r.CreateNameSupplier(ReplaceIdGenerators_RenameStrategy_INCONSISTENT, r.previousMap.Get(key)))
	} else if doc.IsMappedIdGenerator() {
		supplier := r.nameGenerators.Get(key)
		if supplier == nil || supplier.GetRenameStrategy() != ReplaceIdGenerators_RenameStrategy_MAPPED {
			r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_MISSING_NAME_MAP_FOR_GENERATOR))
			// skip registering the name in the list of Generators if there no
			// mapping.
			return
		}
	} else {
		panic(jutil.NewIllegalStateExceptionMsg("unexpected"))
	}
	r.idGeneratorMaps.Put(key, jutil.NewLinkedHashMap[string, string]())
}

// Java: ReplaceIdGenerators.process(Node, Node)
func (r *ReplaceIdGenerators) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, NewReplaceIdGenerators_GatherGenerators(r))
	if !r.nameGenerators.IsEmpty() {
		NodeTraversal_Traverse(r.compiler, root, NewReplaceIdGenerators_ReplaceGenerators(r))
	}
}

// ReplaceIdGenerators_ReplaceGenerators is the inner class
// ReplaceGenerators.
type ReplaceIdGenerators_ReplaceGenerators struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	this_            *ReplaceIdGenerators // the enclosing instance
	hasNullGenerator bool
	lastPropSet      map[string]struct{}
}

// Java: ReplaceGenerators()
func NewReplaceIdGenerators_ReplaceGenerators(outer *ReplaceIdGenerators) *ReplaceIdGenerators_ReplaceGenerators {
	g := &ReplaceIdGenerators_ReplaceGenerators{this_: outer}
	if outer != nil && outer.nameGenerators != nil && !outer.nameGenerators.IsEmpty() {
		g.lastPropSet = make(map[string]struct{}, outer.nameGenerators.Size())
		for key := range outer.nameGenerators.Keys() {
			if key == "null" {
				g.hasNullGenerator = true
			}
			if lastDot := strings.LastIndexByte(key, '.'); lastDot >= 0 {
				g.lastPropSet[key[lastDot+1:]] = struct{}{}
			} else {
				g.lastPropSet[key] = struct{}{}
			}
		}
	}
	return g
}

// Java: ReplaceGenerators.visit(NodeTraversal, Node, Node)
func (g *ReplaceIdGenerators_ReplaceGenerators) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	r := g.this_
	if !n.IsCall() && !n.IsTaggedTemplateLit() {
		return
	}

	qname := NodeUtil_GetCallTargetResolvingIndirectCalls(n)
	var callName string
	if qname.IsName() {
		callName = qname.GetString()
		if callName == "" {
			callName = "null"
		}
	} else if qname.IsGetProp() {
		if _, ok := g.lastPropSet[qname.GetString()]; !ok && !g.hasNullGenerator {
			return
		}
		if s, ok := qname.GetQualifiedNameString(); ok {
			callName = s
		} else {
			callName = "null"
		}
	} else {
		if s, ok := qname.GetQualifiedNameString(); ok {
			callName = s
		} else {
			callName = "null"
		}
	}
	nameGenerator := r.nameGenerators.Get(callName)
	if nameGenerator == nil {
		return
	}

	if !t.InGlobalHoistScope() &&
		nameGenerator.GetRenameStrategy() == ReplaceIdGenerators_RenameStrategy_INCONSISTENT {
		// Warn about calls not in the global scope.
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_NON_GLOBAL_ID_GENERATOR_CALL))
		return
	}

	if nameGenerator.GetRenameStrategy() == ReplaceIdGenerators_RenameStrategy_INCONSISTENT {
		for ancestor := range n.GetAncestors().All() {
			if NodeUtil_IsControlStructure(ancestor) {
				// Warn about conditional calls.
				r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL))
				return
			}
		}
	}

	if n.IsCall() {
		g.MaybeReplaceCall(t, n, callName, nameGenerator)
	} else {
		g.MaybeReplaceTaggedTemplateLit(t, n, callName, nameGenerator)
	}
}

// Java: ReplaceGenerators.maybeReplaceTaggedTemplateLit(NodeTraversal, Node, String, NameSupplier)
func (g *ReplaceIdGenerators_ReplaceGenerators) MaybeReplaceTaggedTemplateLit(
	t *NodeTraversal, n rhino.Node, callName string, nameGenerator ReplaceIdGenerators_NameSupplier) {
	r := g.this_
	arg := n.GetLastChild()

	if arg.IsNil() || !arg.IsTemplateLit() {
		panic(jutil.NewIllegalStateException())
	} else if arg.HasOneChild() {
		cooked := arg.GetFirstChild().GetCookedString()
		if cooked == nil {
			// We don't allow strings with odd escape sequences... We could but it doesn't seem
			// necessary
			r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER))
			return
		}

		rename := g.GetObfuscatedName(arg, callName, nameGenerator, *cooked)
		n.ReplaceWith(rhino.IR_String(rename))
		t.ReportCodeChange()
	} else {
		// There is an alternating sequence of template literals and expressions, in this case we
		// need to preserve the function call but obfuscate the literals.
		newTemplateLit := rhino.IR_TemplateLiteral()
		for child := arg.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if child.IsTemplateLitString() {
				cooked := child.GetCookedString()
				if cooked == nil {
					r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER))
					return
				}
				rename := g.GetObfuscatedName(child, callName, nameGenerator, *cooked)
				newTemplateLit.AddChildToBack(
					rhino.IR_TemplateLiteralString(&rename, rename).SrcrefIfMissing(child))
			} else {
				newTemplateLit.AddChildToBack(
					rhino.IR_TemplateLiteralSubstitution(child.RemoveFirstChild()).SrcrefIfMissing(child))
			}
		}
		arg.ReplaceWith(newTemplateLit)
		t.ReportCodeChangeWithNode(newTemplateLit)
	}
}

// Java: ReplaceGenerators.maybeReplaceCall(NodeTraversal, Node, String, NameSupplier)
func (g *ReplaceIdGenerators_ReplaceGenerators) MaybeReplaceCall(
	t *NodeTraversal, n rhino.Node, callName string, nameGenerator ReplaceIdGenerators_NameSupplier) {
	r := g.this_
	arg := n.GetSecondChild()
	if arg.IsNil() {
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER))
	} else if arg.IsStringLit() {
		rename := g.GetObfuscatedName(arg, callName, nameGenerator, arg.GetString())
		n.ReplaceWith(rhino.IR_String(rename))
		t.ReportCodeChange()
	} else if arg.IsTemplateLit() && arg.HasOneChild() {
		cooked := arg.GetFirstChild().GetCookedString()
		if cooked == nil {
			r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER))
		} else {
			rename := g.GetObfuscatedName(arg, callName, nameGenerator, *cooked)
			n.ReplaceWith(rhino.IR_String(rename))
			t.ReportCodeChange()
		}
	} else if arg.IsObjectLit() {
		for key := arg.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
			if key.IsMemberFunctionDef() {
				r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN))
				return
			}
			if key.IsComputedProp() {
				r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN))
				return
			}

			rename := g.GetObfuscatedName(key, callName, nameGenerator, key.GetString())
			key.SetString(rename)
			// Prevent standard renaming by marking the key as quoted.
			key.PutBooleanProp(rhino.Node_QUOTED_PROP, true)
		}
		arg.Detach()
		n.ReplaceWith(arg)
		t.ReportCodeChange()
	} else if r.templateLiteralsAreTranspiled && arg.IsName() {
		// This might be a transpiled template literal.  Find the definition.
		// The pass always injects it at the current script root, typically at the top, but not
		// always.
		ttlVar := g.FindTtlVar(t, arg)
		if ttlVar.IsNil() {
			// This could be some other call to the ttl function not using ttl syntax.  These are
			// weird but not illegal.
			r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER))
			return
		}
		ttlFunctionCall := ttlVar.GetFirstFirstChild()
		paramsArrayLiteral := ttlFunctionCall.GetSecondChild()
		if paramsArrayLiteral.IsNil() || !paramsArrayLiteral.IsArrayLit() {
			panic(jutil.NewIllegalStateExceptionMsg("bad transpiled structure"))
		}
		// 2 cases, if there is exactly 1 element we can just replace everything
		// otherwise we need to rewrite the array in place.
		if paramsArrayLiteral.GetChildCount() == 1 {
			originalName := paramsArrayLiteral.GetFirstChild().GetString()
			rename := g.GetObfuscatedName(arg, callName, nameGenerator, originalName)
			n.ReplaceWith(rhino.IR_String(rename))
			t.ReportCodeChange()
			t.ReportCodeChangeWithNode(ttlVar)
			ttlVar.Detach()
		} else {
			// We have parameters, so we need to rewrite the TTL array in place and leave the function
			// call.
			for param := paramsArrayLiteral.GetFirstChild(); !param.IsNil(); param = param.GetNext() {
				originalName := param.GetString()
				rename := g.GetObfuscatedName(arg, callName, nameGenerator, originalName)
				param.SetString(rename)
			}
			t.ReportCodeChangeWithNode(paramsArrayLiteral)
		}
	} else {
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER))
	}
}

// Java: ReplaceGenerators.findTtlVar(NodeTraversal, Node)
func (g *ReplaceIdGenerators_ReplaceGenerators) FindTtlVar(t *NodeTraversal, arg rhino.Node) rhino.Node {
	name := arg.GetString()
	ttlVarName := t.GetScope().GetVar(name).GetNode()
	parent := ttlVarName.GetParent()
	// We are only interested in top level var definitions since that is what the transpiler
	// creates
	if !parent.IsVar() || !parent.GetParent().IsScript() {
		return rhino.NilNode
	}
	// Because the AST is normalized there is only one child
	jutil.CheckStateMsg(
		ttlVarName.GetParent().HasOneChild(),
		"The AST is normalized so there should only be one name per var")
	callNode := ttlVarName.GetFirstChild()
	if !callNode.IsCall() {
		return rhino.NilNode
	}
	callExpr := callNode.GetFirstChild()
	if callExpr.MatchesName("$jscomp$createTemplateTagFirstArg") ||
		// The dot case is about our unit tests mostly.
		ReplaceIdGenerators_CREATE_TEMPLATE_TAG_FIRST_ARG.Matches(callExpr) {
		return ttlVarName.GetParent()
	}

	return rhino.NilNode
}

// Java: ReplaceGenerators.getObfuscatedName(Node, String, NameSupplier, String)
func (g *ReplaceIdGenerators_ReplaceGenerators) GetObfuscatedName(
	id rhino.Node, callName string, nameGenerator ReplaceIdGenerators_NameSupplier, name string) string {
	r := g.this_
	var rename string
	idGeneratorMap := r.idGeneratorMaps.Get(callName)
	instanceId := ReplaceIdGenerators_GetIdForGeneratorNode(
		nameGenerator.GetRenameStrategy() != ReplaceIdGenerators_RenameStrategy_INCONSISTENT, id, name)
	if nameGenerator.GetRenameStrategy() == ReplaceIdGenerators_RenameStrategy_CONSISTENT {
		entry := r.consistNameMap.Get(callName)
		var ok bool
		rename, ok = entry.Lookup(instanceId)
		if !ok {
			rename = nameGenerator.GetName(instanceId, name)
			entry.Put(instanceId, rename)
		}
	} else {
		rename = nameGenerator.GetName(instanceId, name)
	}
	idGeneratorMap.Put(rename, instanceId)
	return rename
}

// GetSerializedIdMappings returns the serialize map of generators and their
// ids and their replacements.
//
// Java: ReplaceIdGenerators.getSerializedIdMappings()
func (r *ReplaceIdGenerators) GetSerializedIdMappings() string {
	return IdMappingUtil_GenerateSerializedIdMappings(r.idGeneratorMaps)
}

// Java: ReplaceIdGenerators.getIdForGeneratorNode(boolean, Node, String)
func ReplaceIdGenerators_GetIdForGeneratorNode(consistent bool, n rhino.Node, name string) string {
	if consistent {
		return name
	} else {
		return jutil.ToJavaString(n.GetSourceFileName()) + ":" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
	}
}
