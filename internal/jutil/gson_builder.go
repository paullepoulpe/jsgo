// Port of com/google/gson/GsonBuilder.java (Gson 2.9.1), with
// TypeAdapters.newTypeHierarchyFactory of
// com/google/gson/internal/bind/TypeAdapters.java

/*
 * Copyright (C) 2008 Google Inc. (GsonBuilder.java)
 * Copyright (C) 2011 Google Inc. (TypeAdapters.java)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jutil

// Go-only: language-adapter — com.google.gson.GsonBuilder (Gson 2.9.1, as
// shaded into the pinned closure-compiler.jar): the settings that
// Gson.toJson uses (gson_to_json.go).

// gson_typeHierarchyFactory is the factory that
// TypeAdapters.newTypeHierarchyFactory(Class, TypeAdapter) returns: it
// applies typeAdapter to every value that is an instance of baseType.
type gson_typeHierarchyFactory struct {
	baseType    func(value any) bool
	typeAdapter TypeAdapter
}

// GsonBuilder is com.google.gson.GsonBuilder (the settings toJson uses).
type GsonBuilder struct {
	hierarchyFactories []*gson_typeHierarchyFactory
	serializeNulls     bool
	prettyPrinting     bool
	escapeHtmlChars    bool
}

// NewGsonBuilder is new GsonBuilder(), whose configuration is new Gson()'s.
//
// Java: GsonBuilder()
func NewGsonBuilder() *GsonBuilder {
	return &GsonBuilder{escapeHtmlChars: true}
}

// SerializeNulls makes Gson serialize null fields and map values.
//
// Java: GsonBuilder.serializeNulls()
func (b *GsonBuilder) SerializeNulls() *GsonBuilder {
	b.serializeNulls = true
	return b
}

// DisableHtmlEscaping makes Gson pass <, >, &, = and ' through as is.
//
// Java: GsonBuilder.disableHtmlEscaping()
func (b *GsonBuilder) DisableHtmlEscaping() *GsonBuilder {
	b.escapeHtmlChars = false
	return b
}

// SetPrettyPrinting makes Gson indent its output by two spaces per level.
//
// Java: GsonBuilder.setPrettyPrinting()
func (b *GsonBuilder) SetPrettyPrinting() *GsonBuilder {
	b.prettyPrinting = true
	return b
}

// RegisterTypeHierarchyAdapter registers typeAdapter for baseType and all
// its subtypes. baseType stands for the Class<?> argument: it reports
// whether a value is an instance of the class (`value instanceof Base`).
//
// Java: GsonBuilder.registerTypeHierarchyAdapter(Class<?>, Object)
func (b *GsonBuilder) RegisterTypeHierarchyAdapter(baseType func(value any) bool, typeAdapter TypeAdapter) *GsonBuilder {
	CheckArgument(typeAdapter != nil)
	b.hierarchyFactories = append(b.hierarchyFactories,
		&gson_typeHierarchyFactory{baseType: baseType, typeAdapter: typeAdapter})
	return b
}

// Create creates a Gson instance with the builder's configuration. Type
// hierarchy adapters registered later take precedence.
//
// Java: GsonBuilder.create()
func (b *GsonBuilder) Create() *Gson {
	factories := make([]*gson_typeHierarchyFactory, 0, len(b.hierarchyFactories))
	for i := len(b.hierarchyFactories) - 1; i >= 0; i-- { // Collections.reverse
		factories = append(factories, b.hierarchyFactories[i])
	}
	return &Gson{
		lenient:            false,
		serializeNulls:     b.serializeNulls,
		htmlSafe:           b.escapeHtmlChars,
		prettyPrinting:     b.prettyPrinting,
		factoriesToBeAdded: factories,
	}
}
