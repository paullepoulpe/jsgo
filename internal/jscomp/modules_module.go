// Port of com/google/javascript/jscomp/modules/Module.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/modules/AutoBuilder_Module_Builder.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"strings"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Module holds information for modules, particularly ES modules, that is
// useful for rewriting. The primary pieces of information are what
// variables are exported (transitive or local), and what names are
// imported.
//
// Go-only: language-adapter — the AutoBuilder-generated
// AutoBuilder_Module_Builder is folded into Module_Builder.
type Module struct {
	metadata               *ModuleMetadataMap_ModuleMetadata
	path                   *ModuleLoader_ModulePath
	namespace              *jutil.LinkedHashMap[string, *Binding]
	boundNames             *jutil.LinkedHashMap[string, *Binding]
	localNameToLocalExport *jutil.LinkedHashMap[string, *Export]
	closureNamespace       *string
	hash                   int32
}

// Java: Module(ModuleMetadata, ModulePath, ImmutableMap<String, Binding>, ImmutableMap<String, Binding>, ImmutableMap<String, Export>, String)
func NewModule(
	metadata *ModuleMetadataMap_ModuleMetadata,
	path *ModuleLoader_ModulePath,
	namespace *jutil.LinkedHashMap[string, *Binding],
	boundNames *jutil.LinkedHashMap[string, *Binding],
	localNameToLocalExport *jutil.LinkedHashMap[string, *Export],
	closureNamespace *string) *Module {
	jutil.CheckNotNullMsg(metadata, "metadata")
	jutil.CheckNotNullMsg(namespace, "namespace")
	jutil.CheckNotNullMsg(boundNames, "boundNames")
	jutil.CheckNotNullMsg(localNameToLocalExport, "localNameToLocalExport")
	return &Module{
		metadata:               metadata,
		path:                   path,
		namespace:              namespace,
		boundNames:             boundNames,
		localNameToLocalExport: localNameToLocalExport,
		closureNamespace:       closureNamespace,
	}
}

// Java: Module.metadata()
func (m *Module) Metadata() *ModuleMetadataMap_ModuleMetadata { return m.metadata }

// Path may return nil.
//
// Java: Module.path()
func (m *Module) Path() *ModuleLoader_ModulePath { return m.path }

// Java: Module.namespace()
func (m *Module) Namespace() *jutil.LinkedHashMap[string, *Binding] { return m.namespace }

// Java: Module.boundNames()
func (m *Module) BoundNames() *jutil.LinkedHashMap[string, *Binding] { return m.boundNames }

// Java: Module.localNameToLocalExport()
func (m *Module) LocalNameToLocalExport() *jutil.LinkedHashMap[string, *Export] {
	return m.localNameToLocalExport
}

// ClosureNamespace may return nil.
//
// Java: Module.closureNamespace()
func (m *Module) ClosureNamespace() *string { return m.closureNamespace }

// Equals is the record's componentwise equals (ModuleMetadata and
// ModulePath compare by identity, the maps by their entries).
func (m *Module) Equals(o any) bool {
	other, ok := o.(*Module)
	return ok && other != nil &&
		(m == other || (m.metadata == other.metadata &&
			m.path == other.path &&
			modules_bindingMapEquals(m.namespace, other.namespace) &&
			modules_bindingMapEquals(m.boundNames, other.boundNames) &&
			modules_exportMapEquals(m.localNameToLocalExport, other.localNameToLocalExport) &&
			modules_nullableStringEquals(m.closureNamespace, other.closureNamespace)))
}

// HashCode is the record's hashCode.
//
// Go-only: language-adapter — identity hashes use IdentityHashCodeOfPointer.
func (m *Module) HashCode() int32 {
	if h := atomic.LoadInt32(&m.hash); h != 0 {
		return h
	}
	return m.computeHashCode()
}

//go:noinline
func (m *Module) computeHashCode() int32 {
	h := int32(0)
	h = 31*h + jutil.IdentityHashCodeOfPointer(m.metadata)
	h = 31*h + jutil.IdentityHashCodeOfPointer(m.path)
	h = 31*h + modules_bindingMapHash(m.namespace)
	h = 31*h + modules_bindingMapHash(m.boundNames)
	h = 31*h + modules_exportMapHash(m.localNameToLocalExport)
	h = 31*h + modules_nullableStringHash(m.closureNamespace)
	atomic.StoreInt32(&m.hash, h)
	return h
}

func modules_bindingMapEquals(a, b *jutil.LinkedHashMap[string, *Binding]) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil || a.Size() != b.Size() {
		return false
	}
	for k, v := range a.All() {
		ov, ok := b.Lookup(k)
		if !ok || !(v == ov || (v != nil && v.Equals(ov))) {
			return false
		}
	}
	return true
}

func modules_bindingMapHash(m *jutil.LinkedHashMap[string, *Binding]) int32 {
	if m == nil || m.IsEmpty() {
		return 0
	}
	var sum int32
	for k, v := range m.All() {
		var vh int32
		if v != nil {
			vh = v.HashCode()
		}
		sum += jutil.StringHashCode(k) ^ vh
	}
	return sum
}

func modules_exportMapEquals(a, b *jutil.LinkedHashMap[string, *Export]) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil || a.Size() != b.Size() {
		return false
	}
	for k, v := range a.All() {
		ov, ok := b.Lookup(k)
		if !ok || !v.equalsExport(ov) {
			return false
		}
	}
	return true
}

func modules_exportMapHash(m *jutil.LinkedHashMap[string, *Export]) int32 {
	if m == nil || m.IsEmpty() {
		return 0
	}
	var sum int32
	for k, v := range m.All() {
		var vh int32
		if v != nil {
			vh = v.HashCode()
		}
		sum += jutil.StringHashCode(k) ^ vh
	}
	return sum
}

// String is the record's toString.
func (m *Module) String() string {
	path := "null"
	if m.path != nil {
		path = m.path.String()
	}
	return "Module[metadata=" + m.metadata.String() +
		", path=" + path +
		", namespace=" + m.namespace.String() +
		", boundNames=" + m.boundNames.String() +
		", localNameToLocalExport=" + m.localNameToLocalExport.String() +
		", closureNamespace=" + modules_nullableString(m.closureNamespace) + "]"
}

// Module_Builder_ creates a new builder.
//
// Java: Module.builder()
func Module_Builder_() *Module_Builder {
	return &Module_Builder{}
}

// ToBuilder returns this module in builder form.
//
// Java: Module.toBuilder()
func (m *Module) ToBuilder() *Module_Builder {
	b := &Module_Builder{m: *m}
	b.m.hash = 0
	return b
}

// Module_Builder is the builder for Module (AutoBuilder_Module_Builder).
type Module_Builder struct {
	m     Module
	built bool
}

// Java: AutoBuilder_Module_Builder.metadata(ModuleMetadataMap.ModuleMetadata)
func (b *Module_Builder) Metadata(metadata *ModuleMetadataMap_ModuleMetadata) *Module_Builder {
	if metadata == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null metadata"))
	}
	b.m.metadata = metadata
	return b
}

// Path takes a possibly-nil value.
//
// Java: AutoBuilder_Module_Builder.path(ModuleLoader.ModulePath)
func (b *Module_Builder) Path(path *ModuleLoader_ModulePath) *Module_Builder {
	b.m.path = path
	return b
}

// Java: AutoBuilder_Module_Builder.namespace(ImmutableMap<String, Binding>)
func (b *Module_Builder) Namespace(namespace *jutil.LinkedHashMap[string, *Binding]) *Module_Builder {
	if namespace == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null namespace"))
	}
	b.m.namespace = namespace
	return b
}

// Java: AutoBuilder_Module_Builder.boundNames(ImmutableMap<String, Binding>)
func (b *Module_Builder) BoundNames(boundNames *jutil.LinkedHashMap[string, *Binding]) *Module_Builder {
	if boundNames == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null boundNames"))
	}
	b.m.boundNames = boundNames
	return b
}

// Java: AutoBuilder_Module_Builder.localNameToLocalExport(ImmutableMap<String, Export>)
func (b *Module_Builder) LocalNameToLocalExport(localNameToLocalExport *jutil.LinkedHashMap[string, *Export]) *Module_Builder {
	if localNameToLocalExport == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null localNameToLocalExport"))
	}
	b.m.localNameToLocalExport = localNameToLocalExport
	return b
}

// ClosureNamespace takes a possibly-nil value.
//
// Java: AutoBuilder_Module_Builder.closureNamespace(String)
func (b *Module_Builder) ClosureNamespace(closureNamespace *string) *Module_Builder {
	b.m.closureNamespace = closureNamespace
	return b
}

// Java: AutoBuilder_Module_Builder.build()
func (b *Module_Builder) Build() *Module {
	if b.m.metadata == nil ||
		b.m.namespace == nil ||
		b.m.boundNames == nil ||
		b.m.localNameToLocalExport == nil {
		var sb strings.Builder
		sb.WriteString("Missing required properties:")
		if b.m.metadata == nil {
			sb.WriteString(" metadata")
		}
		if b.m.namespace == nil {
			sb.WriteString(" namespace")
		}
		if b.m.boundNames == nil {
			sb.WriteString(" boundNames")
		}
		if b.m.localNameToLocalExport == nil {
			sb.WriteString(" localNameToLocalExport")
		}
		panic(jutil.NewIllegalStateExceptionMsg(sb.String()))
	}
	if !b.built {
		b.built = true
		return &b.m
	}
	cp := b.m
	cp.hash = 0
	return &cp
}
