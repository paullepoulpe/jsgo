// Port of com/google/javascript/jscomp/serialization/SerializeTypesToPointers.java @ closure-compiler v20261021

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SerializeTypesToPointers grabs an integer TypePool pointer for each JSType
// on the AST and logs information about the pointers.
type SerializeTypesToPointers struct {
	compiler                  AbstractCompiler
	jstypeReconserializer     *JSTypeReconserializer
	propertiesReferencedInAst jutil.BitSet
	typePointersByJstype      map[rhino.JSType]int
	typePool                  *TypePool
	lastPropID                uint32
	gathered                  bool
}

var serializeTypesToPointers_GSON = jutil.NewGson()

// Java: SerializeTypesToPointers.create(AbstractCompiler, StringPool.Builder, SerializationOptions)
func SerializeTypesToPointers_Create(
	compiler AbstractCompiler,
	stringPoolBuilder *StringPool_Builder,
	serializationOptions *SerializationOptions) *SerializeTypesToPointers {
	invalidatingTypes :=
		NewInvalidatingTypes_Builder(compiler.GetTypeRegistry()).
			AddAllTypeMismatches(compiler.GetTypeMismatches()).
			Build()

	s := &SerializeTypesToPointers{
		compiler: compiler,
	}
	jsTypeReconserializer :=
		JSTypeReconserializer_Create(
			compiler.GetTypeRegistry(),
			invalidatingTypes,
			stringPoolBuilder,
			nil,
			serializationOptions)
	jsTypeReconserializer.allowedProperties = &s.propertiesReferencedInAst
	s.jstypeReconserializer = jsTypeReconserializer
	return s
}

// Java: SerializeTypesToPointers.gatherTypesOnAst(Node)
func (s *SerializeTypesToPointers) GatherTypesOnAst(root rhino.Node) {
	jutil.CheckStateMsg(!s.gathered, "Cannot call process() twice")
	s.gathered = true
	s.gatherPropertiesAndTypesWalk(root)

	// these types are only used when debug logging is enabled, but we always serialize them as not
	// to have a different TypePool with and without debug logging.
	for _, mismatch := range s.compiler.GetTypeMismatches() {
		s.jstypeReconserializer.SerializeType(mismatch.Found())
		s.jstypeReconserializer.SerializeType(mismatch.Required())
	}

	s.typePool = s.jstypeReconserializer.GenerateTypePool()

	s.logSerializationDebugInfo(s.jstypeReconserializer)
}

func (s *SerializeTypesToPointers) gatherTypesForColors(root rhino.Node) {
	jutil.CheckStateMsg(!s.gathered, "Cannot call process() twice")
	s.gathered = true
	s.gatherPropertiesAndTypesWalk(root)

	// Always record mismatch types so pointer assignments and debug logs are
	// identical to GatherTypesOnAst without allocating a TypePool/StringPool.
	for _, mismatch := range s.compiler.GetTypeMismatches() {
		s.jstypeReconserializer.SerializeType(mismatch.Found())
		s.jstypeReconserializer.SerializeType(mismatch.Required())
	}

	s.jstypeReconserializer.prepareForColorPool()

	s.logSerializationDebugInfo(s.jstypeReconserializer)
}

func (s *SerializeTypesToPointers) gatherPropertiesAndTypesWalk(n rhino.Node) {
	if n.IsScript() && NodeUtil_IsFromTypeSummary(n) {
		return
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		s.gatherPropertiesAndTypesWalk(c)
	}
	switch n.GetToken() {
	case rhino.Token_GETPROP, // "name" from (someObject.name)
		rhino.Token_OPTCHAIN_GETPROP: // "name" from (someObject?.name)
		propID := n.GetStringID()
		if propID != s.lastPropID {
			s.lastPropID = propID
			s.propertiesReferencedInAst.Set(int(propID))
		}
	case rhino.Token_STRING_KEY, // "name" from obj = {name: 0}
		rhino.Token_MEMBER_FUNCTION_DEF, // "name" from class C { name() {} }
		rhino.Token_MEMBER_FIELD_DEF,    // "name" from class C { name = 0; }
		rhino.Token_GETTER_DEF,          // "name" from class C { get name() {} }
		rhino.Token_SETTER_DEF:          // "name" from class C { set name(n) {} }
		if !n.IsQuotedStringKey() {
			propID := n.GetStringID()
			if propID != s.lastPropID {
				s.lastPropID = propID
				s.propertiesReferencedInAst.Set(int(propID))
			}
		}
	case rhino.Token_ROOT:
		// the ROOT node is given the 'global this' type for use during typechecking, but this is
		// not needed for optimizations. (The 'global this' type will still be serialized if
		// referenced in actual code though.)
		return
	case rhino.Token_NAME:
		if tb := n.GetJSTypeBase(); tb != nil && !s.jstypeReconserializer.hasTypePointer(tb) {
			s.jstypeReconserializer.SerializeType(tb.Self())
		}
		if shadow := n.GetClosureUnawareShadow(); !shadow.IsNil() {
			// Shadow roots are structured as
			// ROOT -> SCRIPT -> EXPR_RESULT -> FUNCTION
			s.gatherTypesWalk(shadow.GetFirstFirstChild().GetFirstChild())
		}
		return
	}
	if tb := n.GetJSTypeBase(); tb != nil && !s.jstypeReconserializer.hasTypePointer(tb) {
		s.jstypeReconserializer.SerializeType(tb.Self())
	}
}

func (s *SerializeTypesToPointers) gatherTypesWalk(n rhino.Node) {
	if n.IsScript() && NodeUtil_IsFromTypeSummary(n) {
		return
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		s.gatherTypesWalk(c)
	}
	if n.IsRoot() {
		return
	}
	if tb := n.GetJSTypeBase(); tb != nil && !s.jstypeReconserializer.hasTypePointer(tb) {
		s.jstypeReconserializer.SerializeType(tb.Self())
	}
	if n.IsName() {
		if shadow := n.GetClosureUnawareShadow(); !shadow.IsNil() {
			s.gatherTypesWalk(shadow.GetFirstFirstChild().GetFirstChild())
		}
	}
}

// Java: SerializeTypesToPointers.getTypePointersByJstype()
func (s *SerializeTypesToPointers) GetTypePointersByJstype() map[rhino.JSType]int {
	if s.typePointersByJstype == nil {
		s.typePointersByJstype = make(map[rhino.JSType]int, len(s.jstypeReconserializer.recordsByPointer))
		for _, record := range s.jstypeReconserializer.recordsByPointer {
			for _, jstype := range record.jstypes {
				s.typePointersByJstype[jstype] = record.pointer
			}
		}
	}
	return s.typePointersByJstype
}

// Java: SerializeTypesToPointers.getTypePool()
func (s *SerializeTypesToPointers) GetTypePool() *TypePool {
	return s.typePool
}

// Java: SerializeTypesToPointers.logSerializationDebugInfo(JSTypeReconserializer, TypePool)
func (s *SerializeTypesToPointers) logSerializationDebugInfo(serializer *JSTypeReconserializer) {
	// Log information about how the JSTypes correspond to the colors. This may be useful later on
	// in optimizations.
	func() {
		log := s.compiler.CreateOrReopenLog("SerializeTypesToPointers", "object_uuids.log")
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
		log := s.compiler.CreateOrReopenLog("SerializeTypesToPointers", "mismatches.log")
		defer log.Close()
		if log.IsLogging() {
			log.LogSupplier(
				func() string {
					return serializeTypesToPointers_GSON.ToJson(s.logTypeMismatches(s.compiler.GetTypeMismatches(), serializer))
				})
		}
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
	typeMismatches []*TypeMismatch, serializer *JSTypeReconserializer) []*serializeTypesToPointers_TypeMismatchJson {
	result := make([]*serializeTypesToPointers_TypeMismatchJson, 0, len(typeMismatches))
	for _, mismatch := range typeMismatches {
		result = append(result, serializeTypesToPointers_TypeMismatchJson_Create(mismatch, serializer))
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
	allSerializedTypes *jutil.LinkedHashMap[string, []rhino.JSType]
}

// Java: SerializeTypesToPointers.StreamObjectUuidsJson(ImmutableMap<String, Collection<JSType>>)
func newSerializeTypesToPointers_StreamObjectUuidsJson(allSerializedTypes *jutil.LinkedHashMap[string, []rhino.JSType]) *serializeTypesToPointers_StreamObjectUuidsJson {
	return &serializeTypesToPointers_StreamObjectUuidsJson{allSerializedTypes: allSerializedTypes}
}

// Java: SerializeTypesToPointers.StreamObjectUuidsJson.writeJson(JsonWriter)
func (p *serializeTypesToPointers_StreamObjectUuidsJson) WriteJson(jsonWriter *jutil.JsonWriter) {
	jsonWriter.BeginObject()
	// allSerializedTypes.asMap().entrySet()
	for key, types := range p.allSerializedTypes.All() {
		jsonWriter.Name(key) // color id
		jsonWriter.BeginArray()
		// sort JSTypes by string representation
		var jstypes []string
		for _, jstype := range types {
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
	Location        string `gson:"location"`
	FoundColorId    string `gson:"foundColorId"`
	RequiredColorId string `gson:"requiredColorId"`
}

// Java: SerializeTypesToPointers.TypeMismatchJson(TypeMismatch, ColorId, ColorId)
func newSerializeTypesToPointers_TypeMismatchJson(x *TypeMismatch, found colors.ColorId, required colors.ColorId) *serializeTypesToPointers_TypeMismatchJson {
	return &serializeTypesToPointers_TypeMismatchJson{
		Location:        x.Location().GetLocation(),
		FoundColorId:    found.String(),
		RequiredColorId: required.String(),
	}
}

// Java: SerializeTypesToPointers.TypeMismatchJson.create(TypeMismatch, JSTypeReconserializer, TypePool)
func serializeTypesToPointers_TypeMismatchJson_Create(
	x *TypeMismatch, serializer *JSTypeReconserializer) *serializeTypesToPointers_TypeMismatchJson {
	foundRecord := serializer.RecordType(x.Found())
	requiredRecord := serializer.RecordType(x.Required())
	return newSerializeTypesToPointers_TypeMismatchJson(x, foundRecord.colorId, requiredRecord.colorId)
}

// Java: SerializeTypesToPointers.TypeMismatchJson.compareTo(TypeMismatchJson)
func (j *serializeTypesToPointers_TypeMismatchJson) CompareTo(x *serializeTypesToPointers_TypeMismatchJson) int {
	if c := jutil.StringCompareTo(j.FoundColorId, x.FoundColorId); c != 0 {
		return c
	}
	if c := jutil.StringCompareTo(j.RequiredColorId, x.RequiredColorId); c != 0 {
		return c
	}
	return jutil.StringCompareTo(j.Location, x.Location)
}
