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

// ReplaceIdGenerators replaces calls to id generators with unique, short ids.
//
// Java uses null as a map key for a generator whose name is unknown (JSDoc
// @idGenerator on a node other than an assignment to a qualified name, a name
// declaration or a function, or a call whose target has no qualified name).
// Such keys are represented by the string "null" here: no qualified name can be
// "null" (it is a keyword), and IdMappingUtil serializes Java's null key as
// "null" as well.
type ReplaceIdGenerators struct {
	compiler                      AbstractCompiler
	templateLiteralsAreTranspiled bool
	nameGenerators                map[string]idNameSupplier
	consistNameMap                map[string]map[string]string

	generatorOrder  []string
	idGeneratorMaps map[string]*idGenMapping
	previousMap     *jutil.LinkedHashMap[string, *jutil.BiMap[string, string]]

	generatePseudoNames bool
	xidHashFunction     Xid_HashFunction
}

type idGenMapping struct {
	keys   []string
	values map[string]string
}

func (m *idGenMapping) put(key, value string) {
	if m.values == nil {
		m.values = make(map[string]string)
	}
	if _, exists := m.values[key]; !exists {
		m.keys = append(m.keys, key)
	}
	m.values[key] = value
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

var replaceIdGeneratorsCreateTemplateTagFirstArg = rhino.QualifiedName_Of("$jscomp.createTemplateTagFirstArg")

// Java: ReplaceIdGenerators(AbstractCompiler, boolean, Map<String, RenamingMap>, boolean, String, Xid.HashFunction)
func NewReplaceIdGenerators(
	compiler AbstractCompiler,
	templateLiteralsAreTranspiled bool,
	idGens *jutil.LinkedHashMap[string, RenamingMap],
	generatePseudoNames bool,
	previousMapSerialized *string,
	xidHashFunction Xid_HashFunction) *ReplaceIdGenerators {
	r := &ReplaceIdGenerators{
		compiler:                      compiler,
		templateLiteralsAreTranspiled: templateLiteralsAreTranspiled,
		generatePseudoNames:           generatePseudoNames,
		xidHashFunction:               xidHashFunction,
		nameGenerators:                make(map[string]idNameSupplier),
		idGeneratorMaps:               make(map[string]*idGenMapping),
	}

	if previousMapSerialized != nil && *previousMapSerialized != "" {
		r.previousMap = IdMappingUtil_ParseSerializedIdMappings(previousMapSerialized)
	}

	if idGens != nil {
		for name, map_ := range idGens.All() {
			if renamingToken, ok := map_.(RenamingToken); ok {
				switch renamingToken {
				case RenamingToken_DISABLE:
					r.nameGenerators[name] = nil
					continue // don't put an entry in idGeneratorMaps
				case RenamingToken_INCONSISTENT:
					r.nameGenerators[name] = r.createNameSupplier(idRenameInconsistent, r.prevMappingsFor(name))
				case RenamingToken_STABLE:
					r.nameGenerators[name] = r.createNameSupplier(idRenameStable, r.prevMappingsFor(name))
				}
			} else {
				r.nameGenerators[name] = &mappedIdNameSupplier{map_: map_}
			}
			r.registerGeneratorMap(name)
		}
	}
	return r
}

func (r *ReplaceIdGenerators) prevMappingsFor(name string) *jutil.BiMap[string, string] {
	if r.previousMap == nil {
		return nil
	}
	return r.previousMap.Get(name)
}

func (r *ReplaceIdGenerators) registerGeneratorMap(name string) {
	if _, exists := r.idGeneratorMaps[name]; !exists {
		r.generatorOrder = append(r.generatorOrder, name)
	}
	r.idGeneratorMaps[name] = &idGenMapping{}
}

type idRenameStrategy uint8

const (
	idRenameConsistent idRenameStrategy = iota
	idRenameInconsistent
	idRenameMapped
	idRenameStable
	idRenameXid
)

type idNameSupplier interface {
	getName(id, name string) string
	renameStrategy() idRenameStrategy
}

type obfuscatedIdNameSupplier struct {
	generator        NameGenerator
	previousMappings *jutil.BiMap[string, string]
	strategy         idRenameStrategy
}

func newObfuscatedIdNameSupplier(strategy idRenameStrategy, previousMappings *jutil.BiMap[string, string]) *obfuscatedIdNameSupplier {
	reservedNames := jutil.NewLinkedHashSet[string]()
	var inv *jutil.BiMap[string, string]
	if previousMappings != nil {
		inv = previousMappings.Inverse()
		for k := range previousMappings.Keys() {
			reservedNames.Add(k)
		}
	}
	return &obfuscatedIdNameSupplier{
		generator:        NewDefaultNameGeneratorWithReserved(reservedNames, "", nil),
		previousMappings: inv,
		strategy:         strategy,
	}
}

func (s *obfuscatedIdNameSupplier) getName(id, name string) string {
	if s.previousMappings != nil {
		if newName, ok := s.previousMappings.Lookup(id); ok {
			return newName
		}
	}
	return s.generator.GenerateNextName()
}

func (s *obfuscatedIdNameSupplier) renameStrategy() idRenameStrategy { return s.strategy }

type pseudoIdNameSupplier struct {
	counter  int32
	strategy idRenameStrategy
}

func (s *pseudoIdNameSupplier) getName(id, name string) string {
	if s.strategy == idRenameInconsistent {
		c := s.counter
		s.counter++
		return name + "$" + strconv.Itoa(int(c))
	}
	return name + "$0"
}

func (s *pseudoIdNameSupplier) renameStrategy() idRenameStrategy { return s.strategy }

type stableIdNameSupplier struct{}

func (stableIdNameSupplier) getName(id, name string) string {
	return sourcemap.Base64_Base64EncodeInt(jutil.StringHashCode(name))
}

func (stableIdNameSupplier) renameStrategy() idRenameStrategy { return idRenameStable }

type xidNameSupplier struct {
	xid *Xid
}

func newXidNameSupplier(hashFunction Xid_HashFunction) *xidNameSupplier {
	if hashFunction == nil {
		return &xidNameSupplier{xid: NewXid()}
	}
	return &xidNameSupplier{xid: NewXidWithHasher(hashFunction)}
}

func (s *xidNameSupplier) getName(id, name string) string   { return s.xid.Get(name) }
func (s *xidNameSupplier) renameStrategy() idRenameStrategy { return idRenameXid }

type mappedIdNameSupplier struct {
	map_ RenamingMap
}

func (s *mappedIdNameSupplier) getName(id, name string) string {
	return *s.map_.Get(name)
}

func (s *mappedIdNameSupplier) renameStrategy() idRenameStrategy { return idRenameMapped }

func (r *ReplaceIdGenerators) createNameSupplier(
	strategy idRenameStrategy, previousMappings *jutil.BiMap[string, string]) idNameSupplier {
	switch {
	case strategy == idRenameStable:
		return stableIdNameSupplier{}
	case strategy == idRenameXid:
		return newXidNameSupplier(r.xidHashFunction)
	case r.generatePseudoNames:
		return &pseudoIdNameSupplier{strategy: strategy}
	default:
		return newObfuscatedIdNameSupplier(strategy, previousMappings)
	}
}

type gatherIdGenerators struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	r *ReplaceIdGenerators
}

func (g *gatherIdGenerators) Visit(unused *NodeTraversal, n rhino.Node, parent rhino.Node) {
	r := g.r
	doc := n.GetJSDocInfo()
	if doc == nil || !doc.IsAnyIdGenerator() {
		return
	}

	key := "null"
	if n.IsAssign() {
		if qname, ok := n.GetFirstChild().GetQualifiedNameString(); ok {
			key = qname
		}
	} else if NodeUtil_IsNameDeclaration(n) {
		key = n.GetFirstChild().GetString()
	} else if n.IsFunction() {
		s := n.GetFirstChild().GetString()
		if s == "" {
			return
		}
		key = s
	}

	if _, exists := r.nameGenerators[key]; exists {
		// This generator is already registered from our constructor; don't override it.
		return
	}
	if doc.IsConsistentIdGenerator() {
		if r.consistNameMap == nil {
			r.consistNameMap = make(map[string]map[string]string)
		}
		r.consistNameMap[key] = make(map[string]string)
		r.nameGenerators[key] = r.createNameSupplier(idRenameConsistent, r.prevMappingsFor(key))
	} else if doc.IsStableIdGenerator() {
		r.nameGenerators[key] = r.createNameSupplier(idRenameStable, r.prevMappingsFor(key))
	} else if doc.IsXidGenerator() {
		r.nameGenerators[key] = r.createNameSupplier(idRenameXid, r.prevMappingsFor(key))
	} else if doc.IsIdGenerator() {
		r.nameGenerators[key] = r.createNameSupplier(idRenameInconsistent, r.prevMappingsFor(key))
	} else if doc.IsMappedIdGenerator() {
		supplier := r.nameGenerators[key]
		if supplier == nil || supplier.renameStrategy() != idRenameMapped {
			r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_MISSING_NAME_MAP_FOR_GENERATOR))
			return
		}
	} else {
		panic(jutil.NewIllegalStateExceptionMsg("unexpected"))
	}
	r.registerGeneratorMap(key)
}

// Java: ReplaceIdGenerators.process(Node, Node)
func (r *ReplaceIdGenerators) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, &gatherIdGenerators{r: r})
	if len(r.nameGenerators) > 0 {
		NodeTraversal_Traverse(r.compiler, root, newReplaceIdGeneratorsPass(r))
	}
}

type replaceIdGeneratorsPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	r                *ReplaceIdGenerators
	hasNullGenerator bool
	lastPropSet      map[string]struct{}
}

func newReplaceIdGeneratorsPass(r *ReplaceIdGenerators) *replaceIdGeneratorsPass {
	g := &replaceIdGeneratorsPass{r: r}
	if len(r.nameGenerators) > 0 {
		g.lastPropSet = make(map[string]struct{}, len(r.nameGenerators))
		for key := range r.nameGenerators { // order-independent: populates unordered lookup set and boolean flag
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

func (g *replaceIdGeneratorsPass) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	r := g.r
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
	} else if s, ok := qname.GetQualifiedNameString(); ok {
		callName = s
	} else {
		callName = "null"
	}

	nameGenerator := r.nameGenerators[callName]
	if nameGenerator == nil {
		return
	}

	if nameGenerator.renameStrategy() == idRenameInconsistent {
		if !t.InGlobalHoistScope() {
			r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_NON_GLOBAL_ID_GENERATOR_CALL))
			return
		}
		for ancestor := n.GetParent(); !ancestor.IsNil(); ancestor = ancestor.GetParent() {
			if NodeUtil_IsControlStructure(ancestor) {
				r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL))
				return
			}
		}
	}

	if n.IsCall() {
		g.maybeReplaceCall(t, n, callName, nameGenerator)
	} else {
		g.maybeReplaceTaggedTemplateLit(t, n, callName, nameGenerator)
	}
}

func (g *replaceIdGeneratorsPass) maybeReplaceTaggedTemplateLit(
	t *NodeTraversal, n rhino.Node, callName string, nameGenerator idNameSupplier) {
	r := g.r
	arg := n.GetLastChild()

	if arg.IsNil() || !arg.IsTemplateLit() {
		panic(jutil.NewIllegalStateException())
	}
	if arg.HasOneChild() {
		cooked := arg.GetFirstChild().GetCookedString()
		if cooked == nil {
			r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER))
			return
		}
		rename := g.getObfuscatedName(arg, callName, nameGenerator, *cooked)
		n.ReplaceWith(rhino.IR_String(rename))
		t.ReportCodeChange()
		return
	}

	// Alternating sequence of template literals and expressions: preserve the call, obfuscate literals.
	newTemplateLit := rhino.IR_TemplateLiteral()
	for child := arg.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if child.IsTemplateLitString() {
			cooked := child.GetCookedString()
			if cooked == nil {
				r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER))
				return
			}
			rename := g.getObfuscatedName(child, callName, nameGenerator, *cooked)
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

func (g *replaceIdGeneratorsPass) maybeReplaceCall(
	t *NodeTraversal, n rhino.Node, callName string, nameGenerator idNameSupplier) {
	r := g.r
	arg := n.GetSecondChild()
	if arg.IsNil() {
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER))
	} else if arg.IsStringLit() {
		rename := g.getObfuscatedName(arg, callName, nameGenerator, arg.GetString())
		n.ReplaceWith(rhino.IR_String(rename))
		t.ReportCodeChange()
	} else if arg.IsTemplateLit() && arg.HasOneChild() {
		cooked := arg.GetFirstChild().GetCookedString()
		if cooked == nil {
			r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER))
		} else {
			rename := g.getObfuscatedName(arg, callName, nameGenerator, *cooked)
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

			rename := g.getObfuscatedName(key, callName, nameGenerator, key.GetString())
			key.SetString(rename)
			// Prevent standard renaming by marking the key as quoted.
			key.PutBooleanProp(rhino.Node_QUOTED_PROP, true)
		}
		arg.Detach()
		n.ReplaceWith(arg)
		t.ReportCodeChange()
	} else if r.templateLiteralsAreTranspiled && arg.IsName() {
		ttlVar := findTtlVar(t, arg)
		if ttlVar.IsNil() {
			r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER))
			return
		}
		ttlFunctionCall := ttlVar.GetFirstFirstChild()
		paramsArrayLiteral := ttlFunctionCall.GetSecondChild()
		if paramsArrayLiteral.IsNil() || !paramsArrayLiteral.IsArrayLit() {
			panic(jutil.NewIllegalStateExceptionMsg("bad transpiled structure"))
		}
		if paramsArrayLiteral.GetChildCount() == 1 {
			originalName := paramsArrayLiteral.GetFirstChild().GetString()
			rename := g.getObfuscatedName(arg, callName, nameGenerator, originalName)
			n.ReplaceWith(rhino.IR_String(rename))
			t.ReportCodeChange()
			t.ReportCodeChangeWithNode(ttlVar)
			ttlVar.Detach()
		} else {
			for param := paramsArrayLiteral.GetFirstChild(); !param.IsNil(); param = param.GetNext() {
				rename := g.getObfuscatedName(arg, callName, nameGenerator, param.GetString())
				param.SetString(rename)
			}
			t.ReportCodeChangeWithNode(paramsArrayLiteral)
		}
	} else {
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER))
	}
}

func findTtlVar(t *NodeTraversal, arg rhino.Node) rhino.Node {
	ttlVarName := t.GetScope().GetVar(arg.GetString()).GetNode()
	parent := ttlVarName.GetParent()
	// Only interested in top-level var definitions created by the transpiler.
	if !parent.IsVar() || !parent.GetParent().IsScript() {
		return rhino.NilNode
	}
	jutil.CheckStateMsg(
		parent.HasOneChild(),
		"The AST is normalized so there should only be one name per var")
	callNode := ttlVarName.GetFirstChild()
	if !callNode.IsCall() {
		return rhino.NilNode
	}
	callExpr := callNode.GetFirstChild()
	if callExpr.MatchesName("$jscomp$createTemplateTagFirstArg") ||
		replaceIdGeneratorsCreateTemplateTagFirstArg.Matches(callExpr) {
		return parent
	}
	return rhino.NilNode
}

func (g *replaceIdGeneratorsPass) getObfuscatedName(
	id rhino.Node, callName string, nameGenerator idNameSupplier, name string) string {
	r := g.r
	strategy := nameGenerator.renameStrategy()
	instanceId := getIdForGeneratorNode(strategy != idRenameInconsistent, id, name)
	var rename string
	if strategy == idRenameConsistent {
		entry := r.consistNameMap[callName]
		var ok bool
		rename, ok = entry[instanceId]
		if !ok {
			rename = nameGenerator.getName(instanceId, name)
			entry[instanceId] = rename
		}
	} else {
		rename = nameGenerator.getName(instanceId, name)
	}
	r.idGeneratorMaps[callName].put(rename, instanceId)
	return rename
}

// GetSerializedIdMappings returns the serialized map of generators and their
// ids and their replacements.
//
// Java: ReplaceIdGenerators.getSerializedIdMappings()
func (r *ReplaceIdGenerators) GetSerializedIdMappings() string {
	var sb strings.Builder
	for _, key := range r.generatorOrder {
		m := r.idGeneratorMaps[key]
		if m == nil || len(m.keys) == 0 {
			continue
		}
		sb.WriteByte('[')
		sb.WriteString(key)
		sb.WriteString("]\n\n")
		for _, rKey := range m.keys {
			sb.WriteString(rKey)
			sb.WriteByte(':')
			sb.WriteString(m.values[rKey])
			sb.WriteByte('\n')
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func getIdForGeneratorNode(consistent bool, n rhino.Node, name string) string {
	if consistent {
		return name
	}
	return jutil.ToJavaString(n.GetSourceFileName()) + ":" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
}
