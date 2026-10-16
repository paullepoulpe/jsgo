// Port of com/google/javascript/jscomp/modules/Module.java @ closure-compiler v20260415
// Port of com/google/javascript/jscomp/modules/AutoBuilder_Module_Builder.java @ closure-compiler v20260415

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Module holds information for modules, particularly ES modules, that is
// useful for rewriting. The primary pieces of information are what
// variables are exported (transitive or local), and what names are
// imported.
//
// path is the path of this module, nil if this module is a nested
// goog.loadModule.
//
// namespace is the map of exported identifiers to originating binding. Note
// that the keys are different than boundNames(). Here the keys are the
// exported name, not the local name. Examples:
//   - Locally defined exports (e.g. export let x;, let x; export {x as v};)
//     creates an entry with the exported name for the local module's export
//     definition.
//   - export default function foo() {}; creates an entry with the name
//     "default" for the local module's default export definition.
//   - export {x as v} from 'mod'; creates an entry with the name "v" for the
//     export definition x from 'mod'.
//   - import statements make no entries on their own. If imported values
//     are exported with export {}; then an entry is created like export {}
//     from.
//   - exports.foo = bar; creates an entry with the name "foo" for the
//     expression on the right-hand side. This is not bound to a local name.
//
// boundNames is the map of local identifiers to originating binding. This
// includes all names bound by import and exported names which originate in
// this module. Used for rewriting in later stages of the compiler. ES
// modules may have names bound by both imports and exports. Closure modules
// only have names bound by imports, as it is impossible to create a new
// local identifier in an export. Examples:
//   - import {x as v} from 'mod'; creates an entry with the name "v" for the
//     export definition x from 'mod'.
//   - import * as ns from 'mod'; creates an entry with the name "ns" with a
//     binding containing all of mod's bindings.
//   - export default function foo() {} creates an entry with the name "foo"
//     for the local module's export definition.
//   - export {x as v} from 'mod'; does not create any entry in this module.
//   - const C = goog.require('mod.C') creates an entry with the name "C" for
//     the binding containing the default export of 'mod.C'
//
// localNameToLocalExport is the map of local identifier name to local
// export definition.
//
// closureNamespace is the specific Closure namespace this module
// represents, if any. This can be from goog.provide, goog.module, or
// goog.module.declareNamespace. nil otherwise.
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
		m.metadata == other.metadata &&
		m.path == other.path &&
		jutil.ObjectsEquals(m.namespace, other.namespace) &&
		jutil.ObjectsEquals(m.boundNames, other.boundNames) &&
		jutil.ObjectsEquals(m.localNameToLocalExport, other.localNameToLocalExport) &&
		modules_nullableStringEquals(m.closureNamespace, other.closureNamespace)
}

// HashCode is the record's hashCode.
//
// Go-only: language-adapter — identity hashes are jutil.IdentityHashCode.
func (m *Module) HashCode() int32 {
	h := int32(0)
	h = 31*h + jutil.IdentityHashCode(m.metadata)
	h = 31*h + jutil.IdentityHashCode(m.path)
	h = 31*h + jutil.ObjectHashCode(m.namespace)
	h = 31*h + jutil.ObjectHashCode(m.boundNames)
	h = 31*h + jutil.ObjectHashCode(m.localNameToLocalExport)
	h = 31*h + modules_nullableStringHash(m.closureNamespace)
	return h
}

// String is the record's toString.
func (m *Module) String() string {
	return "Module[metadata=" + m.metadata.String() +
		", path=" + jutil.ToJavaString(m.path) +
		", namespace=" + m.namespace.String() +
		", boundNames=" + m.boundNames.String() +
		", localNameToLocalExport=" + m.localNameToLocalExport.String() +
		", closureNamespace=" + modules_nullableString(m.closureNamespace) + "]"
}

// Module_Builder_ creates a new builder.
//
// Java: Module.builder()
func Module_Builder_() *Module_Builder {
	return newAutoBuilder_Module_Builder()
}

// ToBuilder returns this module in builder form.
//
// Java: Module.toBuilder()
func (m *Module) ToBuilder() *Module_Builder {
	return newAutoBuilder_Module_BuilderFrom(m)
}

// Module_Builder is the builder for Module (AutoBuilder_Module_Builder).
type Module_Builder struct {
	metadata               *ModuleMetadataMap_ModuleMetadata
	path                   *ModuleLoader_ModulePath
	namespace              *jutil.LinkedHashMap[string, *Binding]
	boundNames             *jutil.LinkedHashMap[string, *Binding]
	localNameToLocalExport *jutil.LinkedHashMap[string, *Export]
	closureNamespace       *string
}

// Java: AutoBuilder_Module_Builder()
func newAutoBuilder_Module_Builder() *Module_Builder {
	return &Module_Builder{}
}

// Java: AutoBuilder_Module_Builder(Module)
func newAutoBuilder_Module_BuilderFrom(source *Module) *Module_Builder {
	return &Module_Builder{
		metadata:               source.Metadata(),
		path:                   source.Path(),
		namespace:              source.Namespace(),
		boundNames:             source.BoundNames(),
		localNameToLocalExport: source.LocalNameToLocalExport(),
		closureNamespace:       source.ClosureNamespace(),
	}
}

// Java: AutoBuilder_Module_Builder.metadata(ModuleMetadataMap.ModuleMetadata)
func (b *Module_Builder) Metadata(metadata *ModuleMetadataMap_ModuleMetadata) *Module_Builder {
	if metadata == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null metadata"))
	}
	b.metadata = metadata
	return b
}

// Path takes a possibly-nil value.
//
// Java: AutoBuilder_Module_Builder.path(ModuleLoader.ModulePath)
func (b *Module_Builder) Path(path *ModuleLoader_ModulePath) *Module_Builder {
	b.path = path
	return b
}

// Java: AutoBuilder_Module_Builder.namespace(ImmutableMap<String, Binding>)
func (b *Module_Builder) Namespace(namespace *jutil.LinkedHashMap[string, *Binding]) *Module_Builder {
	if namespace == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null namespace"))
	}
	b.namespace = namespace
	return b
}

// Java: AutoBuilder_Module_Builder.boundNames(ImmutableMap<String, Binding>)
func (b *Module_Builder) BoundNames(boundNames *jutil.LinkedHashMap[string, *Binding]) *Module_Builder {
	if boundNames == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null boundNames"))
	}
	b.boundNames = boundNames
	return b
}

// Java: AutoBuilder_Module_Builder.localNameToLocalExport(ImmutableMap<String, Export>)
func (b *Module_Builder) LocalNameToLocalExport(localNameToLocalExport *jutil.LinkedHashMap[string, *Export]) *Module_Builder {
	if localNameToLocalExport == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null localNameToLocalExport"))
	}
	b.localNameToLocalExport = localNameToLocalExport
	return b
}

// ClosureNamespace takes a possibly-nil value.
//
// Java: AutoBuilder_Module_Builder.closureNamespace(String)
func (b *Module_Builder) ClosureNamespace(closureNamespace *string) *Module_Builder {
	b.closureNamespace = closureNamespace
	return b
}

// Java: AutoBuilder_Module_Builder.build()
func (b *Module_Builder) Build() *Module {
	if b.metadata == nil ||
		b.namespace == nil ||
		b.boundNames == nil ||
		b.localNameToLocalExport == nil {
		missing := ""
		if b.metadata == nil {
			missing += " metadata"
		}
		if b.namespace == nil {
			missing += " namespace"
		}
		if b.boundNames == nil {
			missing += " boundNames"
		}
		if b.localNameToLocalExport == nil {
			missing += " localNameToLocalExport"
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing))
	}
	return NewModule(
		b.metadata,
		b.path,
		b.namespace,
		b.boundNames,
		b.localNameToLocalExport,
		b.closureNamespace)
}
