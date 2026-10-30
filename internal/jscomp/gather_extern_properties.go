// Port of com/google/javascript/jscomp/GatherExternProperties.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// ExternPropSet stores interned extern property names as both a dense BitSet
// indexed by strpool.ID (for O(1) membership checks without string hashing) and
// a deduplicated insertion-ordered slice of IDs (for deterministic iteration).
type ExternPropSet struct {
	bits jutil.BitSet
	ids  []strpool.ID
}

// NewExternPropSetFromIDs builds an ExternPropSet from interned string IDs,
// preserving first-seen order and ignoring zero (unset) IDs.
func NewExternPropSetFromIDs(ids []strpool.ID) *ExternPropSet {
	s := &ExternPropSet{
		ids: make([]strpool.ID, 0, len(ids)),
	}
	for _, id := range ids {
		s.addID(id)
	}
	return s
}

// NewExternPropSetFromStrings interns each property name in names and builds an
// ExternPropSet preserving first-seen iteration order.
func NewExternPropSetFromStrings(names iter.Seq[string]) *ExternPropSet {
	s := &ExternPropSet{}
	if names != nil {
		for name := range names {
			s.addID(strpool.Intern(name))
		}
	}
	return s
}

func (s *ExternPropSet) addID(id strpool.ID) {
	if id == 0 || s.bits.Get(int(id)) {
		return
	}
	s.bits.Set(int(id))
	s.ids = append(s.ids, id)
}

func (s *ExternPropSet) clone() *ExternPropSet {
	if s == nil {
		return &ExternPropSet{}
	}
	out := &ExternPropSet{
		ids: slices.Clone(s.ids),
	}
	out.bits.Or(&s.bits)
	return out
}

// HasID reports whether id is in the extern property set.
func (s *ExternPropSet) HasID(id strpool.ID) bool {
	return s != nil && id != 0 && s.bits.Get(int(id))
}

// Has reports whether name is in the extern property set without interning
// unknown candidate strings.
func (s *ExternPropSet) Has(name string) bool {
	if s == nil {
		return false
	}
	id := strpool.Lookup(name)
	return id != 0 && s.bits.Get(int(id))
}

// Bits returns the underlying BitSet indexed by uint32(strpool.ID), or nil if s is nil.
func (s *ExternPropSet) Bits() *jutil.BitSet {
	if s == nil {
		return nil
	}
	return &s.bits
}

// IDs returns the deduplicated extern property IDs in insertion order.
func (s *ExternPropSet) IDs() []strpool.ID {
	if s == nil {
		return nil
	}
	return s.ids
}

// Len returns the number of unique extern property names in s.
func (s *ExternPropSet) Len() int {
	if s == nil {
		return 0
	}
	return len(s.ids)
}

// AllStrings yields each extern property name in insertion order.
func (s *ExternPropSet) AllStrings() iter.Seq[string] {
	return func(yield func(string) bool) {
		if s == nil {
			return
		}
		for _, id := range s.ids {
			if !yield(strpool.String(id)) {
				return
			}
		}
	}
}

// ToLinkedHashSet materializes a LinkedHashSet[string] in insertion order, or
// returns nil if s is nil.
func (s *ExternPropSet) ToLinkedHashSet() *jutil.LinkedHashSet[string] {
	if s == nil {
		return nil
	}
	out := jutil.NewLinkedHashSetWithCapacity[string](len(s.ids))
	for _, id := range s.ids {
		out.Add(strpool.String(id))
	}
	return out
}

// Equals reports whether s and other contain the exact same set of extern properties.
func (s *ExternPropSet) Equals(other *ExternPropSet) bool {
	if s == other {
		return true
	}
	if s == nil || other == nil {
		return false
	}
	return len(s.ids) == len(other.ids) && s.bits.Equals(&other.bits)
}

// GatherExternProperties gathers property names defined in externs.
//
// The collection of these property names could easily happen during type
// checking. However, when exporting local property definitions, the externs
// may be modified after type checking, and we want to collect the new names
// as well.
//
// Mode.CHECK is set when checksOnly=true and only collects the externs in
// the type expressions in JSDoc. Mode.OPTIMIZE is set when checksOnly=false
// and only collects the externs in AST properties.
//
// To optimize build time, we collect only type expression externs in
// getChecks() and only AST property externs in getOptimizations() instead of
// collecting all in both.
type GatherExternProperties struct {
	compiler AbstractCompiler

	externProperties *ExternPropSet
	published        bool

	mode GatherExternProperties_Mode
}

// GatherExternProperties_Mode is the Java enum GatherExternProperties.Mode.
type GatherExternProperties_Mode int32

const (
	GatherExternProperties_Mode_CHECK GatherExternProperties_Mode = iota
	GatherExternProperties_Mode_OPTIMIZE
	GatherExternProperties_Mode_CHECK_AND_OPTIMIZE
)

var gatherExternProperties_Mode_names = [...]string{"CHECK", "OPTIMIZE", "CHECK_AND_OPTIMIZE"}

// Java: Mode.check
var gatherExternProperties_Mode_check = [...]bool{true, false, true}

// Java: Mode.optimize
var gatherExternProperties_Mode_optimize = [...]bool{false, true, true}

// Name is Enum.name().
func (m GatherExternProperties_Mode) Name() string { return gatherExternProperties_Mode_names[m] }

// String is Enum.toString().
func (m GatherExternProperties_Mode) String() string { return m.Name() }

// Ordinal is Enum.ordinal().
func (m GatherExternProperties_Mode) Ordinal() int { return int(m) }

// Java: GatherExternProperties(AbstractCompiler, Mode)
func NewGatherExternProperties(compiler AbstractCompiler, mode GatherExternProperties_Mode) *GatherExternProperties {
	return &GatherExternProperties{
		compiler:         compiler,
		externProperties: compiler.GetExternPropSet().clone(),
		mode:             mode,
	}
}

// Java: GatherExternProperties.process(Node, Node)
func (g *GatherExternProperties) Process(externs rhino.Node, root rhino.Node) {
	if g.published {
		g.externProperties = g.externProperties.clone()
		g.published = false
	}
	NodeTraversal_Traverse(g.compiler, externs, g)
	g.published = true
	g.compiler.SetExternPropSet(g.externProperties)
}

// Java: GatherExternProperties.shouldTraverse(NodeTraversal, Node, Node)
func (g *GatherExternProperties) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return !n.IsScript() || !NodeUtil_IsFromTypeSummary(n)
}

// Java: GatherExternProperties.visit(NodeTraversal, Node, Node)
func (g *GatherExternProperties) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if gatherExternProperties_Mode_optimize[g.mode] {
		switch n.GetToken() {
		case rhino.Token_GETPROP:
			// Gathers "name" from (someObject.name).
			g.externProperties.addID(strpool.ID(n.GetStringID()))
		case rhino.Token_STRING_KEY:
			if parent.IsObjectLit() {
				g.externProperties.addID(strpool.ID(n.GetStringID()))
			}
		case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
			g.externProperties.addID(strpool.ID(n.GetStringID()))
		default:
		}
	}

	jsDocInfo := n.GetJSDocInfo()
	if jsDocInfo != nil && gatherExternProperties_Mode_check[g.mode] {
		g.GatherPropertiesFromJSDocInfo(jsDocInfo)
	}
}

// Java: GatherExternProperties.gatherPropertiesFromJSDocInfo(JSDocInfo)
func (g *GatherExternProperties) GatherPropertiesFromJSDocInfo(jsDocInfo *rhino.JSDocInfo) {
	for it := jsDocInfo.TypeExprIter(); it.Next(); {
		g.GatherPropertiesFromJsTypeExpressionNode(it.Node())
	}
}

// Java: GatherExternProperties.gatherPropertiesFromJsTypeExpressionNode(Node)
func (g *GatherExternProperties) GatherPropertiesFromJsTypeExpressionNode(jsTypeExpressionNode rhino.Node) {
	switch jsTypeExpressionNode.GetToken() {
	case rhino.Token_LB:
		g.GatherPropertiesFromJsDocRecordType(jsTypeExpressionNode)
	default:
		for child := jsTypeExpressionNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			g.GatherPropertiesFromJsTypeExpressionNode(child)
		}
	}
}

// Java: GatherExternProperties.gatherPropertiesFromJsDocRecordType(Node)
func (g *GatherExternProperties) GatherPropertiesFromJsDocRecordType(jsDocRecordNode rhino.Node) {
	jutil.CheckStateMsg1(jsDocRecordNode.GetToken() == rhino.Token_LB, "%s", jsDocRecordNode)
	for fieldNode := jsDocRecordNode.GetFirstChild(); !fieldNode.IsNil(); fieldNode = fieldNode.GetNext() {
		var fieldNameNode rhino.Node
		var fieldTypeNode rhino.Node
		if fieldNode.GetToken() == rhino.Token_COLON {
			fieldNameNode = fieldNode.GetFirstChild()
			fieldTypeNode = fieldNameNode.GetNext()
		} else {
			fieldNameNode = fieldNode
			fieldTypeNode = rhino.NilNode
		}
		jutil.CheckStateMsg1(fieldNameNode.IsStringKey(), "%s", fieldNameNode)
		id := strpool.ID(fieldNameNode.GetStringID())
		fieldName := strpool.String(id)
		// TODO(bradfordcsmith): The JSDoc parser should do this.
		if len(fieldName) > 0 && (fieldName[0] == '\'' || fieldName[0] == '"') {
			fieldName = jutil.StringSubstring(fieldName, 1, jutil.StringLength(fieldName)-1)
			id = strpool.Intern(fieldName)
		}
		g.externProperties.addID(id)
		if !fieldTypeNode.IsNil() {
			g.GatherPropertiesFromJsTypeExpressionNode(fieldTypeNode)
		}
	}
}
