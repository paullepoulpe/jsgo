// Port of com/google/javascript/jscomp/js/RuntimeJsLibManager.java @ closure-compiler v20261028
// Port of com/google/javascript/jscomp/InjectRuntimeLibraries.java @ closure-compiler v20261028

/*
 * Copyright 2025 The Closure Compiler Authors.
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
	"strings"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/resources"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RuntimeJsLibManager injects runtime libraries from the jscomp/js
// directory into an AST.
//
// Supports injecting libraries either based on the library path name, or
// by the specific `$jscomp.*` field/class name.
type RuntimeJsLibManager struct {
	changeTracker        *ChangeTracker
	mode                 RuntimeJsLibManager_RuntimeLibraryMode
	resourceProvider     RuntimeJsLibManager_ResourceProvider
	nodeForCodeInsertion func() rhino.Node
	injectedLibs         OrderedSet[string]
	internedFields       OrderedStringMap[*runtimeJsLibField]

	lastInjectedLibrary rhino.Node
}

// RuntimeJsLibManager_RuntimeLibraryMode is what to do with the runtime
// libraries under the js/ directory: the compiler can ignore them
// completely; can validate any attempted library usage is correct but not
// modify the AST; or both validate & add to the AST.
type RuntimeJsLibManager_RuntimeLibraryMode int32

const (
	RuntimeJsLibManager_RuntimeLibraryMode_INJECT RuntimeJsLibManager_RuntimeLibraryMode = iota
	// mode where the compiler is building a TypedAST for the runtime
	// libraries & can't use them, or where the compiler is in
	// transpile-only mode and running per-file.
	RuntimeJsLibManager_RuntimeLibraryMode_NO_OP
	// for testing - compiler records that it was asked to inject a runtime
	// library, but does not modify the AST
	RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY
	// for testing - a stricter version of RECORD_ONLY. The compiler records
	// that it was asked to inject a runtime library, and will also validate
	// injection in JsLibField::assertInjected, but does not modify the AST.
	RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS
	// for @closureUnaware nested compilation. When asked to inject a field,
	// the compiler adds an @extern definition of the field name to prevent
	// misoptimization, and uses a different identifier from the original
	// uncompiled field name. The RuntimeJsLibManager is not responsible for
	// actually defining the external field name definitions; the code using
	// EXTERN_FIELD_NAMES mode must enure that everything field from
	// getExternedFields() is actually loaded at runtime.
	RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES
)

var runtimeJsLibManager_RuntimeLibraryMode_names = [...]string{
	"INJECT", "NO_OP", "RECORD_ONLY", "RECORD_AND_VALIDATE_FIELDS", "EXTERN_FIELD_NAMES",
}

// RuntimeJsLibManager_RuntimeLibraryMode_values is RuntimeLibraryMode.values().
func RuntimeJsLibManager_RuntimeLibraryMode_values() []RuntimeJsLibManager_RuntimeLibraryMode {
	return []RuntimeJsLibManager_RuntimeLibraryMode{
		RuntimeJsLibManager_RuntimeLibraryMode_INJECT,
		RuntimeJsLibManager_RuntimeLibraryMode_NO_OP,
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY,
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS,
		RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES,
	}
}

// Name is Enum.name().
func (m RuntimeJsLibManager_RuntimeLibraryMode) Name() string {
	return runtimeJsLibManager_RuntimeLibraryMode_names[m]
}

// String is Enum.toString().
func (m RuntimeJsLibManager_RuntimeLibraryMode) String() string { return m.Name() }

// RuntimeJsLibManager_ResourceProvider loads /js resources into AST
// format: it returns an AST for the given resource name, with a synthetic
// source file with the given path.
//
// resourceName is for example, es6/set; path is for example,
// /my/filesystem/jscomp/js/es6/set.
//
// Go-only: language-adapter — a functional interface, so a Go func type
// (Java: Node parse(String resourceName, String path)).
type RuntimeJsLibManager_ResourceProvider func(resourceName string, path string) rhino.Node

// Java: RuntimeJsLibManager.RUNTIME_LIB_DIR
const RuntimeJsLibManager_RUNTIME_LIB_DIR = "src/com/google/javascript/jscomp/js/"

type runtimeLibDirHolder struct {
	dir string
}

var runtimeLibDirOverride atomic.Pointer[runtimeLibDirHolder]

// SetRuntimeLibDir overrides the synthetic source-path directory prefix used
// for injected Closure runtime libraries (defaulting to
// RuntimeJsLibManager_RUNTIME_LIB_DIR when dir == ""). If non-empty and missing
// a trailing '/', one is appended automatically.
func SetRuntimeLibDir(dir string) {
	if dir == "" || dir == RuntimeJsLibManager_RUNTIME_LIB_DIR {
		runtimeLibDirOverride.Store(nil)
		return
	}
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	runtimeLibDirOverride.Store(&runtimeLibDirHolder{dir: dir})
}

// GetRuntimeLibDir returns the active runtime library directory prefix.
func GetRuntimeLibDir() string {
	if p := runtimeLibDirOverride.Load(); p != nil {
		return p.dir
	}
	return RuntimeJsLibManager_RUNTIME_LIB_DIR
}

type runtimeLibFieldsSnapshot struct {
	epoch  uint64
	fields map[string]string
}

var (
	runtimeLibFieldsMu    sync.Mutex
	runtimeLibFieldsCache atomic.Pointer[runtimeLibFieldsSnapshot]
)

// runtimeLibFieldsTable returns the read-only lookup table mapping
// $jscomp field names to the runtime library path in which they're defined
// (e.g. "inherits" -> "es6/util/inherits"). Nested properties such as
// "$jscomp.a.b.c" are keyed by their top-level property ("a").
// Automatically reloads when resources.RegisterOverride mutates resources.
func runtimeLibFieldsTable() map[string]string {
	epoch := resources.OverridesEpoch()
	if snap := runtimeLibFieldsCache.Load(); snap != nil && snap.epoch == epoch {
		return snap.fields
	}
	runtimeLibFieldsMu.Lock()
	defer runtimeLibFieldsMu.Unlock()
	epoch = resources.OverridesEpoch()
	if snap := runtimeLibFieldsCache.Load(); snap != nil && snap.epoch == epoch {
		return snap.fields
	}
	fields := loadRuntimeLibFields()
	runtimeLibFieldsCache.Store(&runtimeLibFieldsSnapshot{epoch: epoch, fields: fields})
	return fields
}

// Java: RuntimeJsLibManager.FieldsTable.loadFields()
func loadRuntimeLibFields() map[string]string {
	transpilationLibsTxt := resources.LoadText(resources.JsDir + "transpilation_libs.txt")
	allFields := make(map[string]string, strings.Count(transpilationLibsTxt, "\n")+1)
	// Splitter.on('\n').omitEmptyStrings()
	for rem := transpilationLibsTxt; rem != ""; {
		var line string
		line, rem, _ = strings.Cut(rem, "\n")
		if line == "" {
			continue
		}
		// Splitter.on(',').limit(2).splitToList(line)
		fieldName, resourcePath, ok := strings.Cut(line, ",")
		if !ok {
			panic(jutil.NewIllegalStateExceptionMsg("[" + line + "]"))
		}
		if prev, dup := allFields[fieldName]; dup {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Multiple entries with same key: " + fieldName + "=" + resourcePath + " and " + fieldName + "=" + prev))
		}
		allFields[fieldName] = resourcePath
	}
	return allFields
}

// Java: RuntimeJsLibManager(ChangeTracker, RuntimeLibraryMode, ResourceProvider, Supplier<Node>)
func newRuntimeJsLibManager(
	changeTracker *ChangeTracker,
	mode RuntimeJsLibManager_RuntimeLibraryMode,
	resourceProvider RuntimeJsLibManager_ResourceProvider,
	nodeForCodeInsertion func() rhino.Node) *RuntimeJsLibManager {
	return &RuntimeJsLibManager{
		mode:                 mode,
		changeTracker:        changeTracker,
		resourceProvider:     resourceProvider,
		nodeForCodeInsertion: nodeForCodeInsertion,
	}
}

// Java: RuntimeJsLibManager.create(RuntimeLibraryMode, ResourceProvider, ChangeTracker, Supplier<Node>)
func RuntimeJsLibManager_Create(
	mode RuntimeJsLibManager_RuntimeLibraryMode,
	resourceProvider RuntimeJsLibManager_ResourceProvider,
	changeTracker *ChangeTracker,
	nodeForCodeInsertion func() rhino.Node) *RuntimeJsLibManager {

	return newRuntimeJsLibManager(changeTracker, mode, resourceProvider, nodeForCodeInsertion)
}

// RecordLibraryInjected records that the given library has already been
// injected into the compilation.
//
// This doesn't actually modify the AST - it can be used if e.g. restoring
// compilation state after serialization/deserialization.
//
// resourceName is the name of the library. For example, if "base" is is
// specified, then we record js/base.js
//
// Java: RuntimeJsLibManager.recordLibraryInjected(String)
func (m *RuntimeJsLibManager) RecordLibraryInjected(resourceName string) {
	m.injectedLibs.Add(resourceName)
}

// HasInjectedLibrary returns true if the given library has already been
// injected into the compilation.
//
// Java: RuntimeJsLibManager.hasInjectedLibrary(String)
func (m *RuntimeJsLibManager) HasInjectedLibrary(resourceName string) bool {
	return m.injectedLibs.Contains(resourceName)
}

// GetInjectedLibraries returns a list of all library paths previously
// injected, via one of the other methods on this class.
//
// Java: RuntimeJsLibManager.getInjectedLibraries()
func (m *RuntimeJsLibManager) GetInjectedLibraries() []string {
	// Go-only: language-adapter — ImmutableList.copyOf is never null, and a
	// nil slice reads as null to SerializationOptions.Builder (the
	// typed-AST output of a compilation that injects no library threw
	// NullPointerException("Null runtimeLibraries")).
	if libs := m.injectedLibs.ToSlice(); libs != nil {
		return libs
	}
	return []string{}
}

// GetExternedFields returns a list of all fields previously injected, via
// one of the other methods on this class.
//
// Java: RuntimeJsLibManager.getExternedFields()
func (m *RuntimeJsLibManager) GetExternedFields() []RuntimeJsLibManager_ExternedField {
	if m.mode != RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES {
		return []RuntimeJsLibManager_ExternedField{}
	}
	out := []RuntimeJsLibManager_ExternedField{}
	for _, e := range m.internedFields.Entries() {
		if e.Value.isInjected() {
			out = append(out, e.Value)
		}
	}
	return out
}

// Java: RuntimeJsLibManager.createField(String)
func (m *RuntimeJsLibManager) createField(fieldName string) *runtimeJsLibField {
	if fieldName == "" {
		panic(jutil.NewIllegalArgumentExceptionMsg(fieldName))
	}

	prefix, rest, hasDot := strings.Cut(fieldName, ".")
	if !hasDot {
		panic(jutil.NewIllegalArgumentExceptionMsg("Field name must start with $jscomp., found " + fieldName))
	}
	if prefix != "$jscomp" {
		panic(jutil.NewIllegalArgumentExceptionMsg("Field name must start with $jscomp, found " + fieldName))
	}
	propName, _, _ := strings.Cut(rest, ".")

	resourceName, ok := runtimeLibFieldsTable()[propName]
	if !ok {
		panic(jutil.NewNullPointerExceptionMsg("Cannot find definition of " + fieldName))
	}

	var compiledName string
	switch m.mode {
	case RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES:
		compiledName = strings.ReplaceAll(fieldName, ".", "_")
	case RuntimeJsLibManager_RuntimeLibraryMode_NO_OP, RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY,
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS, RuntimeJsLibManager_RuntimeLibraryMode_INJECT:
		compiledName = fieldName
	}

	return newRuntimeJsLibField(m, resourceName, compiledName, fieldName)
}

// runtimeJsLibField is a $jscomp.* field (could be a method, class, or any
// arbitrary value) that exists in some runtime library under the js/
// directory.
type runtimeJsLibField struct {
	outer *RuntimeJsLibManager

	// The file containing this field, e.g. "es6/generator".
	resourceName string

	// The name by which compiler passes should refer to this field, e.g.
	// `$jscomp.inherits` or `$jscomp_inherits`.
	qualifiedName string

	// The original fully qualified name of this field, e.g.
	// `$jscomp.inherits`.
	uncompiledName string

	injected bool
}

// Java: RuntimeJsLibManager.InternalField(String, String, String)
func newRuntimeJsLibField(
	outer *RuntimeJsLibManager, resourceName string, qualifiedName string, uncompiledName string) *runtimeJsLibField {
	return &runtimeJsLibField{
		outer:          outer,
		resourceName:   resourceName,
		qualifiedName:  qualifiedName,
		uncompiledName: uncompiledName,
	}
}

// Java: RuntimeJsLibManager.InternalField.isInjected()
func (f *runtimeJsLibField) isInjected() bool {
	return f.injected
}

// Java: RuntimeJsLibManager.InternalField.markInjected()
func (f *runtimeJsLibField) markInjected() {
	f.injected = true
}

// Java: RuntimeJsLibManager.InternalField.assertInjected()
func (f *runtimeJsLibField) AssertInjected() RuntimeJsLibManager_InjectedJsLibField {
	switch f.outer.mode {
	case RuntimeJsLibManager_RuntimeLibraryMode_NO_OP, RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY:
	case RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS, RuntimeJsLibManager_RuntimeLibraryMode_INJECT,
		RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES:
		if !f.injected {
			panic(jutil.NewIllegalStateExceptionMsg("Field " + f.uncompiledName + " is not injected"))
		}
	}
	return f
}

// Java: RuntimeJsLibManager.InternalField.matches(Node)
func (f *runtimeJsLibField) Matches(node rhino.Node) bool {
	return node.MatchesQualifiedName(f.qualifiedName)
}

// Don't override .equals/.hashCode and just use reference equality: we
// already intern Fields by qualifiedName.

// Java: RuntimeJsLibManager.InternalField.toString()
func (f *runtimeJsLibField) String() string {
	return "Field<" + f.uncompiledName + ", " + f.resourceName + ">"
}

// Java: RuntimeJsLibManager.InternalField.qualifiedName()
func (f *runtimeJsLibField) QualifiedName() string {
	return f.qualifiedName
}

// Java: RuntimeJsLibManager.InternalField.uncompiledName()
func (f *runtimeJsLibField) UncompiledName() string {
	return f.uncompiledName
}

// Java: RuntimeJsLibManager.InternalField.resourceName()
func (f *runtimeJsLibField) ResourceName() string {
	return f.resourceName
}

// GetJsLibField returns a field definition for the given name, without
// actually asserting or requiring its definition to be injected.
//
// Throws an IllegalArgumentException if passed a nonexistent field name.
//
// Java: RuntimeJsLibManager.getJsLibField(String)
func (m *RuntimeJsLibManager) GetJsLibField(fieldName string) RuntimeJsLibManager_JsLibField {
	// wrapper around getJsLibFieldInternal that returns the broader
	// JsLibField interface
	return m.getJsLibFieldInternal(fieldName)
}

// Java: RuntimeJsLibManager.getJsLibFieldInternal(String)
func (m *RuntimeJsLibManager) getJsLibFieldInternal(fieldName string) *runtimeJsLibField {
	if f, ok := m.internedFields.Lookup(fieldName); ok {
		return f
	}
	f := m.createField(fieldName)
	m.internedFields.Put(fieldName, f)
	return f
}

// RuntimeJsLibManager_JsLibField is a $jscomp.* field (could be a method,
// class, or any arbitrary value) that exists in some runtime library under
// the js/ directory.
type RuntimeJsLibManager_JsLibField interface {
	// Java: RuntimeJsLibManager.JsLibField.matches(Node)
	Matches(node rhino.Node) bool

	// Java: RuntimeJsLibManager.JsLibField.resourceName()
	ResourceName() string

	// Java: RuntimeJsLibManager.JsLibField.assertInjected()
	AssertInjected() RuntimeJsLibManager_InjectedJsLibField
}

// RuntimeJsLibManager_InjectedJsLibField is a JsLibField that is statically
// guaranteed to be available at runtime, i.e. that has been injected into
// the AST.
//
// Only InjectedJsLibField instances provide direct access to the field
// name - this is to prevent code from creating new AST references to a
// field that's not actually injected (yet).
type RuntimeJsLibManager_InjectedJsLibField interface {
	RuntimeJsLibManager_JsLibField

	// QualifiedName returns the name of the field as it is available to the
	// compiler.
	//
	// Java: RuntimeJsLibManager.InjectedJsLibField.qualifiedName()
	QualifiedName() string
}

// RuntimeJsLibManager_ExternedField is an InjectedJsLibField that was
// injected in RuntimeLibraryMode#EXTERN_FIELD_NAMES mode, and so has both
// a regular "qualified name" by which transpilation passes should refer to
// it - this is the externed name - but also provides the uncompiled, raw
// name of the field.
//
// Only ExternedField instances provide access to the uncompiled name. Most
// use sites should use the compiled/qualified name, which may or may not
// equal the uncompiled name. So to try and minimize confusion, we only
// expose the uncompiled name when this is actually an external field
// injection.
type RuntimeJsLibManager_ExternedField interface {
	RuntimeJsLibManager_InjectedJsLibField

	// UncompiledName returns the uncompiled name of the field, as it is
	// seen in the raw runtime lib.
	//
	// Java: RuntimeJsLibManager.ExternedField.uncompiledName()
	UncompiledName() string
}

// InjectLibForField injects the runtime library that defines the given
// $jscomp.* field.
//
// Java: RuntimeJsLibManager.injectLibForField(String)
func (m *RuntimeJsLibManager) InjectLibForField(fieldName string) {
	field := m.getJsLibFieldInternal(fieldName)
	if field.isInjected() {
		return // already done.
	}
	field.markInjected()
	if m.mode == RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES {
		externRoot := rhino.CheckNotNull(m.nodeForCodeInsertion())
		if !externRoot.IsFromExterns() {
			panic(jutil.NewIllegalStateExceptionMsg(externRoot.String()))
		}
		name := ast.NewName(field.QualifiedName())
		name.SetConstantName(true)
		externRoot.AddChildToBack(rhino.IR_Var(name.Node()))
	}
	m.EnsureLibraryInjected(field.resourceName /* force= */, false)
}

// EnsureLibraryInjected: the subdir js/ contains libraries of code that we
// inject at compile-time only if requested by this function.
//
// Notice that these libraries will almost always create global symbols.
//
// resourceName is the name of the library. For example, if "base" is is
// specified, then we load js/base.js. force injects the library even if
// compiler options say not to.
//
// Returns the last node of the most-recently-injected runtime library. If
// new code was injected, this will be the last expression node of the
// library. If the caller needs to add additional code, they should add it
// as the next sibling of this node. If no runtime libraries have been
// injected, then null is returned.
//
// Java: RuntimeJsLibManager.ensureLibraryInjected(String, boolean)
func (m *RuntimeJsLibManager) EnsureLibraryInjected(resourceName string, force bool) rhino.Node {
	if !force {
		switch m.mode {
		case RuntimeJsLibManager_RuntimeLibraryMode_NO_OP:
			return m.lastInjectedLibrary
		case RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY, RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS,
			RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES:
			m.RecordLibraryInjected(resourceName)
			return m.lastInjectedLibrary
		case RuntimeJsLibManager_RuntimeLibraryMode_INJECT: // Keep going.
		}
	}

	if m.injectedLibs.Contains(resourceName) {
		return m.lastInjectedLibrary
	}
	m.RecordLibraryInjected(resourceName)

	path := GetRuntimeLibDir() + resourceName + ".js"
	script := ast.MustScript(rhino.CheckNotNull(m.resourceProvider(resourceName, path)))

	m.injectLibraryDependencies(script, force) // may also modify 'script'
	return m.injectLibrary(script)
}

// Java: RuntimeJsLibManager.injectLibrary(Node)
func (m *RuntimeJsLibManager) injectLibrary(script ast.Script) rhino.Node {
	lastItem, ok := script.LastItem()
	if !ok {
		// Require-only libraries may be empty at this point. Nothing to do
		// here.
		return m.lastInjectedLibrary
	}

	for item, more := script.FirstItem(); more; item, more = item.Next() {
		m.changeTracker.MarkNewScopesChanged(item.Node())
	}
	endOfLib := lastItem.Node()
	firstChild := script.Node().RemoveChildren()

	// Insert the code immediately after the last-inserted runtime library,
	// if any.
	parent := rhino.CheckNotNull(m.nodeForCodeInsertion())
	if m.lastInjectedLibrary.IsNil() {
		parent.AddChildrenToFront(firstChild)
	} else {
		parent.AddChildrenAfter(firstChild, m.lastInjectedLibrary)
	}
	m.lastInjectedLibrary = endOfLib

	m.changeTracker.ReportChangeToEnclosingScope(parent)
	return endOfLib
}

// Java: RuntimeJsLibManager.injectLibraryDependencies(Node, boolean)
func (m *RuntimeJsLibManager) injectLibraryDependencies(script ast.Script, force bool) {
	// Look for string literals of the form 'require foo bar'
	// As we process each one, remove it from its parent.
	for {
		item, ok := script.FirstItem()
		if !ok {
			break
		}
		exprStmt, ok := item.AsExprStmt()
		if !ok {
			break
		}
		strLit, ok := exprStmt.Expr().AsStringLit()
		if !ok {
			break
		}
		directive := strLit.Value()
		kind, arg, hasArg := strings.Cut(directive, " ")
		switch kind {
		// 'use strict' is ignored (and deleted).
		case "use":
		// 'require lib'; pulls in the named library before this one.
		case "require":
			if !hasArg {
				panic(jutil.NewIndexOutOfBoundsException("Index 1 out of bounds for length 1"))
			}
			m.EnsureLibraryInjected(arg, force)
		default:
			panic(jutil.NewIllegalStateExceptionMsg("Bad directive: " + directive))
		}
		item.Detach()
	}
}

// Java: RuntimeJsLibManager.setLastInjectedLibrary(Node)
func (m *RuntimeJsLibManager) SetLastInjectedLibrary(lastInjectedLibrary rhino.Node) {
	m.lastInjectedLibrary = lastInjectedLibrary
}

// Java: RuntimeJsLibManager.getLastInjectedLibrary()
func (m *RuntimeJsLibManager) GetLastInjectedLibrary() rhino.Node {
	return m.lastInjectedLibrary
}

// InjectRuntimeLibraries adds runtime libraries to the beginning of the
// AST. Any libraries explicitly requested via the
// CompilerOptions#forceLibraryInjection field.
//
// TODO(b/120486392): merge this pass with InjectTranspilationRuntimeLibraries.
type InjectRuntimeLibraries struct {
	runtimeLibs            *RuntimeJsLibManager
	forceInjectedLibraries []string
}

// Java: InjectRuntimeLibraries(AbstractCompiler, ImmutableSet<String>)
func NewInjectRuntimeLibraries(compiler AbstractCompiler, forceInjectedLibraries jutil.Set[string]) *InjectRuntimeLibraries {
	if forceInjectedLibraries == nil {
		panic(jutil.NewNullPointerException())
	}
	var libs []string
	if n := forceInjectedLibraries.Size(); n > 0 {
		libs = make([]string, 0, n)
		for lib := range forceInjectedLibraries.All() {
			libs = append(libs, lib)
		}
	}
	return newInjectRuntimeLibrariesFromSlice(compiler, libs)
}

func newInjectRuntimeLibrariesFromSlice(compiler AbstractCompiler, forceInjectedLibraries []string) *InjectRuntimeLibraries {
	return &InjectRuntimeLibraries{
		runtimeLibs:            compiler.GetRuntimeJsLibManager(),
		forceInjectedLibraries: forceInjectedLibraries,
	}
}

// Java: InjectRuntimeLibraries.process(Node, Node)
func (i *InjectRuntimeLibraries) Process(externs rhino.Node, root rhino.Node) {
	i.InjectLibraries()
}

// Java: InjectRuntimeLibraries.injectLibraries()
func (i *InjectRuntimeLibraries) InjectLibraries() {
	for _, forced := range i.forceInjectedLibraries {
		i.runtimeLibs.EnsureLibraryInjected(forced, true)
	}
}
