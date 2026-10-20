// Port of com/google/javascript/jscomp/serialization/SerializeTypesToPointers.java @ closure-compiler v20261015

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"reflect"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SerializeTypesToPointers grabs an integer TypePool pointer for each JSType
// on the AST and logs information about the pointers.
type SerializeTypesToPointers struct {
	compiler                  AbstractCompiler
	jstypeReconserializer     *JSTypeReconserializer
	propertiesReferencedInAst *jutil.LinkedHashSet[string]
	typePointersByJstype      *base.LinkedIdentityHashMap[rhino.JSType, int]
	typePool                  *TypePool
	lastProp                  string
	hasLastProp               bool
	lastType                  rhino.JSType
}

var serializeTypesToPointers_GSON = jutil.NewGson()

// Java: SerializeTypesToPointers(AbstractCompiler, JSTypeReconserializer, LinkedHashSet<String>)
func newSerializeTypesToPointers(
	compiler AbstractCompiler,
	jstypeReconserializer *JSTypeReconserializer,
	propertiesReferencedInAst *jutil.LinkedHashSet[string]) *SerializeTypesToPointers {
	return &SerializeTypesToPointers{
		typePointersByJstype:      base.NewLinkedIdentityHashMap[rhino.JSType, int](),
		compiler:                  compiler,
		jstypeReconserializer:     jstypeReconserializer,
		propertiesReferencedInAst: propertiesReferencedInAst,
	}
}

// Java: SerializeTypesToPointers.create(AbstractCompiler, StringPool.Builder, SerializationOptions)
func SerializeTypesToPointers_Create(
	compiler AbstractCompiler,
	stringPoolBuilder *StringPool_Builder,
	serializationOptions *SerializationOptions) *SerializeTypesToPointers {
	invalidatingTypes :=
		NewInvalidatingTypes_Builder(compiler.GetTypeRegistry()).
			AddAllTypeMismatches(compiler.GetTypeMismatches()).
			Build()

	// this set requires access to the externs and src ASTs, so can't be populated yet.
	propertiesReferencedInAst := jutil.NewLinkedHashSet[string]()
	jsTypeReconserializer :=
		JSTypeReconserializer_Create(
			compiler.GetTypeRegistry(),
			invalidatingTypes,
			stringPoolBuilder,
			propertiesReferencedInAst.Contains,
			serializationOptions)
	return newSerializeTypesToPointers(compiler, jsTypeReconserializer, propertiesReferencedInAst)
}

// Java: SerializeTypesToPointers.gatherTypesOnAst(Node)
func (s *SerializeTypesToPointers) GatherTypesOnAst(root *rhino.Node) {
	jutil.CheckStateMsg(s.typePool == nil, "Cannot call process() twice")
	s.gatherPropertiesWalk(root)
	s.gatherTypesWalk(root)

	// these types are only used when debug logging is enabled, but we always serialize them as not
	// to have a different TypePool with and without debug logging.
	for _, mismatch := range s.compiler.GetTypeMismatches() {
		s.jstypeReconserializer.SerializeType(mismatch.Found())
		s.jstypeReconserializer.SerializeType(mismatch.Required())
	}

	s.typePool = s.jstypeReconserializer.GenerateTypePool()

	s.logSerializationDebugInfo(s.jstypeReconserializer, s.typePool)
}

func (s *SerializeTypesToPointers) gatherPropertiesWalk(n *rhino.Node) {
	if n.IsScript() && NodeUtil_IsFromTypeSummary(n) {
		return
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		s.gatherPropertiesWalk(c)
	}
	switch n.GetToken() {
	case rhino.Token_GETPROP, // "name" from (someObject.name)
		rhino.Token_OPTCHAIN_GETPROP: // "name" from (someObject?.name)
		prop := n.GetString()
		if !s.hasLastProp || prop != s.lastProp {
			s.lastProp = prop
			s.hasLastProp = true
			s.propertiesReferencedInAst.Add(prop)
		}
	case rhino.Token_STRING_KEY, // "name" from obj = {name: 0}
		rhino.Token_MEMBER_FUNCTION_DEF, // "name" from class C { name() {} }
		rhino.Token_MEMBER_FIELD_DEF,    // "name" from class C { name = 0; }
		rhino.Token_GETTER_DEF,          // "name" from class C { get name() {} }
		rhino.Token_SETTER_DEF:          // "name" from class C { set name(n) {} }
		if !n.IsQuotedStringKey() {
			prop := n.GetString()
			if !s.hasLastProp || prop != s.lastProp {
				s.lastProp = prop
				s.hasLastProp = true
				s.propertiesReferencedInAst.Add(prop)
			}
		}
	default:
	}
}

func (s *SerializeTypesToPointers) gatherTypesWalk(n *rhino.Node) {
	if n.IsScript() && NodeUtil_IsFromTypeSummary(n) {
		return
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		s.gatherTypesWalk(c)
	}
	if n.IsRoot() {
		// the ROOT node is given the 'global this' type for use during typechecking, but this is
		// not needed for optimizations. (The 'global this' type will still be serialized if
		// referenced in actual code though.)
		return
	}
	if type_ := n.GetJSType(); type_ != nil && type_ != s.lastType {
		s.lastType = type_
		if _, ok := s.typePointersByJstype.Lookup(type_); !ok {
			s.typePointersByJstype.Put(type_, s.jstypeReconserializer.SerializeType(type_))
		}
	}
	if n.IsName() {
		if shadow := n.GetClosureUnawareShadow(); shadow != nil {
			// Shadow roots are structured as
			// ROOT -> SCRIPT -> EXPR_RESULT -> FUNCTION
			s.gatherTypesWalk(shadow.GetFirstFirstChild().GetFirstChild())
		}
	}
}

// serializeTypesToPointers_PropertySearchCallback finds all unquoted
// property names referenced outside @typeSummary files
//
// Only property names found in this traversal will be serialized onto a
// color's list of "own properties". Properties referenced only inside
// `@typeSummary` files may be excluded because they don't matter for
// optimization of this library's srcs.
//
// Java: private final class SerializeTypesToPointers.PropertySearchCallback implements NodeTraversal.Callback
type serializeTypesToPointers_PropertySearchCallback struct {
	outer *SerializeTypesToPointers
}

// Java: SerializeTypesToPointers.PropertySearchCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *serializeTypesToPointers_PropertySearchCallback) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return !n.IsScript() || !NodeUtil_IsFromTypeSummary(n)
}

// Java: SerializeTypesToPointers.PropertySearchCallback.visit(NodeTraversal, Node, Node)
func (c *serializeTypesToPointers_PropertySearchCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETPROP, // "name" from (someObject.name)
		rhino.Token_OPTCHAIN_GETPROP: // "name" from (someObject?.name)
		c.outer.propertiesReferencedInAst.Add(n.GetString())
	case rhino.Token_STRING_KEY, // "name" from obj = {name: 0}
		rhino.Token_MEMBER_FUNCTION_DEF, // "name" from class C { name() {} }
		rhino.Token_MEMBER_FIELD_DEF,    // "name" from class C { name = 0; }
		rhino.Token_GETTER_DEF,          // "name" from class C { get name() {} }
		rhino.Token_SETTER_DEF:          // "name" from class C { set name(n) {} }
		if !n.IsQuotedStringKey() {
			c.outer.propertiesReferencedInAst.Add(n.GetString())
		}
	default:
	}
}

// Java: private final class SerializeTypesToPointers.TypeSearchCallback implements NodeTraversal.Callback
type serializeTypesToPointers_TypeSearchCallback struct {
	outer *SerializeTypesToPointers
}

// Java: SerializeTypesToPointers.TypeSearchCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *serializeTypesToPointers_TypeSearchCallback) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return !n.IsScript() || !NodeUtil_IsFromTypeSummary(n)
}

// Java: SerializeTypesToPointers.TypeSearchCallback.visit(NodeTraversal, Node, Node)
func (c *serializeTypesToPointers_TypeSearchCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsRoot() {
		// the ROOT node is given the 'global this' type for use during typechecking, but this is
		// not needed for optimizations. (The 'global this' type will still be serialized if
		// referenced in actual code though.)
		return
	}
	type_ := n.GetJSType()
	if type_ != nil {
		c.outer.typePointersByJstype.ComputeIfAbsent(type_, c.outer.jstypeReconserializer.SerializeType)
	}
	shadow := n.GetClosureUnawareShadow()
	if shadow != nil {
		// Shadow roots are structured as
		// ROOT -> SCRIPT -> EXPR_RESULT -> FUNCTION
		NodeTraversal_Traverse(
			c.outer.compiler, shadow.GetFirstFirstChild().GetFirstChild(), &serializeTypesToPointers_TypeSearchCallback{outer: c.outer})
	}
}

// Java: SerializeTypesToPointers.getTypePointersByJstype()
func (s *SerializeTypesToPointers) GetTypePointersByJstype() *base.LinkedIdentityHashMap[rhino.JSType, int] {
	return s.typePointersByJstype
}

// Java: SerializeTypesToPointers.getTypePool()
func (s *SerializeTypesToPointers) GetTypePool() *TypePool {
	return s.typePool
}

// Java: SerializeTypesToPointers.logSerializationDebugInfo(JSTypeReconserializer, TypePool)
func (s *SerializeTypesToPointers) logSerializationDebugInfo(serializer *JSTypeReconserializer, typePool *TypePool) {
	// Log information about how the JSTypes correspond to the colors. This may be useful later on
	// in optimizations.
	func() {
		log := s.compiler.CreateOrReopenLog(reflect.TypeFor[SerializeTypesToPointers](), "object_uuids.log")
		defer log.Close()
		if log.IsLogging() {
			allSerializedTypes :=
				serializer.GetColorIdToJSTypeMapForDebugging()
			// Stream json writing here rather than building up the entire json representation at once
			// because the latter used to cause OOMs.
			log.LogJsonProducer(newSerializeTypesToPointers_StreamObjectUuidsJson(allSerializedTypes))
		}
	}()

	// Log type mismatches, which contribute to the definition of an "invalidating" type
	func() {
		log := s.compiler.CreateOrReopenLog(reflect.TypeFor[SerializeTypesToPointers](), "mismatches.log")
		defer log.Close()
		log.LogSupplier(
			func() string {
				return serializeTypesToPointers_GSON.ToJson(s.logTypeMismatches(s.compiler.GetTypeMismatches(), serializer, typePool))
			})
	}()
}

// logTypeMismatches serializes a type not necessarily attached to an AST
// node.
//
// Not part of the main API for this callback. For use when serializing
// additional types for debug logging.
//
// Java: SerializeTypesToPointers.logTypeMismatches(Iterable<TypeMismatch>, JSTypeReconserializer, TypePool)
func (s *SerializeTypesToPointers) logTypeMismatches(
	typeMismatches []*TypeMismatch, serializer *JSTypeReconserializer, typePool *TypePool) []*serializeTypesToPointers_TypeMismatchJson {
	result := []*serializeTypesToPointers_TypeMismatchJson{}
	for _, mismatch := range typeMismatches {
		result = append(result, serializeTypesToPointers_TypeMismatchJson_Create(mismatch, serializer, typePool))
	}
	// collect(toImmutableSortedSet(naturalOrder())): sorted by compareTo,
	// keeping the first of the elements that compare equal.
	slices.SortStableFunc(result, func(a, b *serializeTypesToPointers_TypeMismatchJson) int { return a.CompareTo(b) })
	return slices.CompactFunc(result, func(a, b *serializeTypesToPointers_TypeMismatchJson) bool { return a.CompareTo(b) == 0 })
}

// serializeTypesToPointers_StreamObjectUuidsJson writes a JSON object whose
// keys are color ids and values are arrays of JSType strings.
//
// Example: `{0: ['*', '?', 'None'], 10f34lksdf: ['SomeProtoCtor']}`
//
// Java: private static class SerializeTypesToPointers.StreamObjectUuidsJson implements LogFile.StreamedJsonProducer
type serializeTypesToPointers_StreamObjectUuidsJson struct {
	allSerializedTypes *jutil.ListMultimap[string, rhino.JSType]
}

// Java: SerializeTypesToPointers.StreamObjectUuidsJson(ImmutableMap<String, Collection<JSType>>)
func newSerializeTypesToPointers_StreamObjectUuidsJson(allSerializedTypes *jutil.ListMultimap[string, rhino.JSType]) *serializeTypesToPointers_StreamObjectUuidsJson {
	return &serializeTypesToPointers_StreamObjectUuidsJson{allSerializedTypes: allSerializedTypes}
}

// Java: SerializeTypesToPointers.StreamObjectUuidsJson.writeJson(JsonWriter)
func (p *serializeTypesToPointers_StreamObjectUuidsJson) WriteJson(jsonWriter *jutil.JsonWriter) {
	jsonWriter.BeginObject()
	// allSerializedTypes.asMap().entrySet()
	for key := range p.allSerializedTypes.KeySet() {
		jsonWriter.Name(key) // color id
		jsonWriter.BeginArray()
		// sort JSTypes by string representation
		var jstypes []string
		for _, jstype := range p.allSerializedTypes.Get(key) {
			jstypes = append(jstypes, jstype.String())
		}
		slices.SortStableFunc(jstypes, jutil.StringCompareTo)
		for _, typeName := range jstypes { // all corresponding JSTypes
			jsonWriter.ValueString(&typeName)
		}
		jsonWriter.EndArray()
	}
	jsonWriter.EndObject()
}

// Java: private static final class SerializeTypesToPointers.TypeMismatchJson implements Comparable<TypeMismatchJson>
type serializeTypesToPointers_TypeMismatchJson struct {
	location        string
	foundColorId    string
	requiredColorId string
}

// Java: SerializeTypesToPointers.TypeMismatchJson(TypeMismatch, ColorId, ColorId)
func newSerializeTypesToPointers_TypeMismatchJson(x *TypeMismatch, found colors.ColorId, required colors.ColorId) *serializeTypesToPointers_TypeMismatchJson {
	return &serializeTypesToPointers_TypeMismatchJson{
		location:        x.Location().GetLocation(),
		foundColorId:    found.String(),
		requiredColorId: required.String(),
	}
}

// Java: SerializeTypesToPointers.TypeMismatchJson.create(TypeMismatch, JSTypeReconserializer, TypePool)
func serializeTypesToPointers_TypeMismatchJson_Create(
	x *TypeMismatch, serializer *JSTypeReconserializer, typePool *TypePool) *serializeTypesToPointers_TypeMismatchJson {
	foundPointer := serializer.SerializeType(x.Found())
	requiredPointer := serializer.SerializeType(x.Required())

	return newSerializeTypesToPointers_TypeMismatchJson(
		x, serializeTypesToPointers_TypeMismatchJson_typePointerToId(foundPointer, typePool),
		serializeTypesToPointers_TypeMismatchJson_typePointerToId(requiredPointer, typePool))
}

// serializeTypesToPointers_TypeMismatchJson_typePointerToId returns the
// unique ID of this pointer if in the type pool.
//
// The given type may not be in the type pool because the type pool was
// generated based on all types reachable from the AST, while a TypeMismatch
// may contain a type in dead code no longer reachable from the AST.
//
// Java: SerializeTypesToPointers.TypeMismatchJson.typePointerToId(int, TypePool)
func serializeTypesToPointers_TypeMismatchJson_typePointerToId(poolOffset int, typePool *TypePool) colors.ColorId {
	if TypePointers_IsAxiomatic(poolOffset) {
		return TypePointers_OFFSET_TO_AXIOMATIC_COLOR[poolOffset].GetId()
	}

	typeProto := typePool.GetTypeList()[TypePointers_TrimOffset(poolOffset)]
	switch typeProto.GetKindCase() {
	case TypeProto_KindCase_UNION:
		ids := jutil.NewLinkedHashSet[colors.ColorId]()
		for _, pointer := range typeProto.GetUnion().GetUnionMemberList() {
			ids.Add(serializeTypesToPointers_TypeMismatchJson_typePointerToId(int(pointer), typePool))
		}
		return colors.ColorId_Union(ids)
	case TypeProto_KindCase_OBJECT:
		return colors.ColorId_FromBytes(typeProto.GetObject().GetUuid())
	case TypeProto_KindCase_KIND_NOT_SET:
	}
	panic(jutil.NewAssertionError("Unrecognized TypeProto " + typeProto.String()))
}

// Java: SerializeTypesToPointers.TypeMismatchJson.compareTo(TypeMismatchJson)
func (j *serializeTypesToPointers_TypeMismatchJson) CompareTo(x *serializeTypesToPointers_TypeMismatchJson) int {
	if c := jutil.StringCompareTo(j.foundColorId, x.foundColorId); c != 0 {
		return c
	}
	if c := jutil.StringCompareTo(j.requiredColorId, x.requiredColorId); c != 0 {
		return c
	}
	return jutil.StringCompareTo(j.location, x.location)
}
