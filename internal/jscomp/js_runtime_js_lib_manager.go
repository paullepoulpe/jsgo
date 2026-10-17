// Port of com/google/javascript/jscomp/js/RuntimeJsLibManager.java @ closure-compiler v20261015

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

package jscomp

import (
	"strings"

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
	nodeForCodeInsertion func() *rhino.Node
	injectedLibs         *jutil.LinkedHashSet[string]
	internedFields       *jutil.LinkedHashMap[string, *RuntimeJsLibManager_InternalField]

	lastInjectedLibrary *rhino.Node
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

// Ordinal is Enum.ordinal().
func (m RuntimeJsLibManager_RuntimeLibraryMode) Ordinal() int { return int(m) }

// RuntimeJsLibManager_ResourceProvider loads /js resources into AST
// format: it returns an AST for the given resource name, with a synthetic
// source file with the given path.
//
// resourceName is for example, es6/set; path is for example,
// /my/filesystem/jscomp/js/es6/set.
//
// Go-only: language-adapter — a functional interface, so a Go func type
// (Java: Node parse(String resourceName, String path)).
type RuntimeJsLibManager_ResourceProvider func(resourceName string, path string) *rhino.Node

// Java: RuntimeJsLibManager.RUNTIME_LIB_DIR
const RuntimeJsLibManager_RUNTIME_LIB_DIR = "src/com/google/javascript/jscomp/js/"

// runtimeJsLibManager_FieldsTable_INSTANCE holds the names of all possible
// $jscomp fields & the file in which they're defined. For example,
// "inherits" -> "es6/inherits". (This doesn't include nested properties -
// so no "$jscomp.a.b.c")
//
// Go-only: language-adapter — Java's FieldsTable holder class loads
// INSTANCE lazily, on first use; so does this variable.
var runtimeJsLibManager_FieldsTable_INSTANCE *jutil.LinkedHashMap[string, string]

func runtimeJsLibManager_FieldsTable() *jutil.LinkedHashMap[string, string] {
	if runtimeJsLibManager_FieldsTable_INSTANCE == nil {
		runtimeJsLibManager_FieldsTable_INSTANCE = runtimeJsLibManager_FieldsTable_loadFields()
	}
	return runtimeJsLibManager_FieldsTable_INSTANCE
}

// Java: RuntimeJsLibManager.FieldsTable.loadFields()
func runtimeJsLibManager_FieldsTable_loadFields() *jutil.LinkedHashMap[string, string] {
	transpilationLibsTxt := resources.ResourceLoader_LoadTextResource(
		"com.google.javascript.jscomp.js.RuntimeJsLibManager", "transpilation_libs.txt")

	allFields := jutil.NewImmutableMapBuilder[string, string]()
	// Splitter.on('\n').omitEmptyStrings()
	for _, line := range strings.Split(transpilationLibsTxt, "\n") {
		if line == "" {
			continue
		}
		// Splitter.on(',').limit(2).splitToList(line)
		tokens := strings.SplitN(line, ",", 2)
		if len(tokens) != 2 {
			// checkState(boolean, Object errorMessage)
			panic(jutil.NewIllegalStateExceptionMsg(jutil.CollectionToString(tokens)))
		}
		fieldName := tokens[0]
		resourcePath := tokens[1]
		allFields.Put(fieldName, resourcePath)
	}
	return allFields.BuildOrThrow()
}

// Java: RuntimeJsLibManager(ChangeTracker, RuntimeLibraryMode, ResourceProvider, Supplier<Node>)
func newRuntimeJsLibManager(
	changeTracker *ChangeTracker,
	mode RuntimeJsLibManager_RuntimeLibraryMode,
	resourceProvider RuntimeJsLibManager_ResourceProvider,
	nodeForCodeInsertion func() *rhino.Node) *RuntimeJsLibManager {
	return &RuntimeJsLibManager{
		mode:                 mode,
		changeTracker:        changeTracker,
		resourceProvider:     resourceProvider,
		nodeForCodeInsertion: nodeForCodeInsertion,
		injectedLibs:         jutil.NewLinkedHashSet[string](),
		internedFields:       jutil.NewLinkedHashMap[string, *RuntimeJsLibManager_InternalField](),
	}
}

// Java: RuntimeJsLibManager.create(RuntimeLibraryMode, ResourceProvider, ChangeTracker, Supplier<Node>)
func RuntimeJsLibManager_Create(
	mode RuntimeJsLibManager_RuntimeLibraryMode,
	resourceProvider RuntimeJsLibManager_ResourceProvider,
	changeTracker *ChangeTracker,
	nodeForCodeInsertion func() *rhino.Node) *RuntimeJsLibManager {

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
	for f := range m.internedFields.Values() {
		if f.isInjected() {
			out = append(out, f)
		}
	}
	return out
}

// Java: RuntimeJsLibManager.createField(String)
func (m *RuntimeJsLibManager) createField(fieldName string) *RuntimeJsLibManager_InternalField {
	if fieldName == "" {
		// checkArgument(boolean, Object errorMessage)
		panic(jutil.NewIllegalArgumentExceptionMsg(fieldName))
	}

	parts := strings.Split(fieldName, ".")
	jutil.CheckArgumentMsg1(len(parts) >= 2, "Field name must start with $jscomp., found %s", fieldName)
	jutil.CheckArgumentMsg1(
		parts[0] == "$jscomp", "Field name must start with $jscomp, found %s", fieldName)
	propName := parts[1]

	resourceName, ok := runtimeJsLibManager_FieldsTable().Lookup(propName)
	if !ok {
		panic(jutil.NewNullPointerExceptionMsg(jutil.LenientFormat("Cannot find definition of %s", fieldName)))
	}

	var compiledName string
	switch m.mode {
	case RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES:
		compiledName = strings.ReplaceAll(fieldName, ".", "_")
	case RuntimeJsLibManager_RuntimeLibraryMode_NO_OP, RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY,
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS, RuntimeJsLibManager_RuntimeLibraryMode_INJECT:
		compiledName = fieldName
	}

	return newRuntimeJsLibManager_InternalField(m, resourceName, compiledName, fieldName)
}

// RuntimeJsLibManager_InternalField is a $jscomp.* field (could be a
// method, class, or any arbitrary value) that exists in some runtime
// library under the js/ directory.
type RuntimeJsLibManager_InternalField struct {
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
func newRuntimeJsLibManager_InternalField(
	outer *RuntimeJsLibManager, resourceName string, qualifiedName string, uncompiledName string) *RuntimeJsLibManager_InternalField {
	return &RuntimeJsLibManager_InternalField{
		outer:          outer,
		resourceName:   resourceName,
		qualifiedName:  qualifiedName,
		uncompiledName: uncompiledName,
	}
}

// Java: RuntimeJsLibManager.InternalField.isInjected()
func (f *RuntimeJsLibManager_InternalField) isInjected() bool {
	return f.injected
}

// Java: RuntimeJsLibManager.InternalField.markInjected()
func (f *RuntimeJsLibManager_InternalField) markInjected() {
	f.injected = true
}

// Java: RuntimeJsLibManager.InternalField.assertInjected()
func (f *RuntimeJsLibManager_InternalField) AssertInjected() RuntimeJsLibManager_InjectedJsLibField {
	switch f.outer.mode {
	case RuntimeJsLibManager_RuntimeLibraryMode_NO_OP, RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY:
	case RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS, RuntimeJsLibManager_RuntimeLibraryMode_INJECT,
		RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES:
		jutil.CheckStateMsg1(f.injected, "Field %s is not injected", f.uncompiledName)
	}
	return f
}

// Java: RuntimeJsLibManager.InternalField.matches(Node)
func (f *RuntimeJsLibManager_InternalField) Matches(node *rhino.Node) bool {
	return node.MatchesQualifiedName(f.qualifiedName)
}

// Don't override .equals/.hashCode and just use reference equality: we
// already intern Fields by qualifiedName.

// Java: RuntimeJsLibManager.InternalField.toString()
func (f *RuntimeJsLibManager_InternalField) String() string {
	return jutil.LenientFormat("Field<%s, %s>", f.uncompiledName, f.resourceName)
}

// Java: RuntimeJsLibManager.InternalField.qualifiedName()
func (f *RuntimeJsLibManager_InternalField) QualifiedName() string {
	return f.qualifiedName
}

// Java: RuntimeJsLibManager.InternalField.uncompiledName()
func (f *RuntimeJsLibManager_InternalField) UncompiledName() string {
	return f.uncompiledName
}

// Java: RuntimeJsLibManager.InternalField.resourceName()
func (f *RuntimeJsLibManager_InternalField) ResourceName() string {
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
func (m *RuntimeJsLibManager) getJsLibFieldInternal(fieldName string) *RuntimeJsLibManager_InternalField {
	return m.internedFields.ComputeIfAbsent(fieldName, m.createField)
}

// RuntimeJsLibManager_JsLibField is a $jscomp.* field (could be a method,
// class, or any arbitrary value) that exists in some runtime library under
// the js/ directory.
type RuntimeJsLibManager_JsLibField interface {
	// Java: RuntimeJsLibManager.JsLibField.matches(Node)
	Matches(node *rhino.Node) bool

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
		externRoot := m.nodeForCodeInsertion()
		if !externRoot.IsFromExterns() {
			// checkState(boolean, Object errorMessage)
			panic(jutil.NewIllegalStateExceptionMsg(externRoot.String()))
		}
		declaration := rhino.IR_Var(rhino.IR_Name(field.QualifiedName()))
		declaration.GetFirstChild().PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		externRoot.AddChildToBack(declaration)
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
func (m *RuntimeJsLibManager) EnsureLibraryInjected(resourceName string, force bool) *rhino.Node {
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

	path := strings.Join([]string{RuntimeJsLibManager_RUNTIME_LIB_DIR, resourceName, ".js"}, "")
	ast := m.resourceProvider(resourceName, path)

	m.injectLibraryDependencies(ast, force) // may also modify 'ast'
	return m.injectLibrary(ast)
}

// Java: RuntimeJsLibManager.injectLibrary(Node)
func (m *RuntimeJsLibManager) injectLibrary(ast *rhino.Node) *rhino.Node {
	if !ast.HasChildren() {
		// Require-only libraries may be empty at this point. Nothing to do
		// here.
		return m.lastInjectedLibrary
	}

	for child := ast.GetFirstChild(); child != nil; child = child.GetNext() {
		m.changeTracker.MarkNewScopesChanged(child)
	}
	endOfLib := ast.GetLastChild()
	firstChild := ast.RemoveChildren()

	// Insert the code immediately after the last-inserted runtime library,
	// if any.
	parent := m.nodeForCodeInsertion()
	if m.lastInjectedLibrary == nil {
		parent.AddChildrenToFront(firstChild)
	} else {
		parent.AddChildrenAfter(firstChild, m.lastInjectedLibrary)
	}
	m.lastInjectedLibrary = endOfLib

	m.changeTracker.ReportChangeToEnclosingScope(parent)
	return endOfLib
}

// Java: RuntimeJsLibManager.injectLibraryDependencies(Node, boolean)
func (m *RuntimeJsLibManager) injectLibraryDependencies(ast *rhino.Node, force bool) {
	// Look for string literals of the form 'require foo bar'
	// As we process each one, remove it from its parent.
	for node := ast.GetFirstChild(); node != nil && node.IsExprResult() && node.GetFirstChild().IsStringLit(); node = ast.GetFirstChild() {
		directive := node.GetFirstChild().GetString()
		// Splitter.on(' ').limit(2).splitToList(directive)
		words := strings.SplitN(directive, " ", 2)
		switch words[0] {
		// 'use strict' is ignored (and deleted).
		case "use":
		// 'require lib'; pulls in the named library before this one.
		case "require":
			if len(words) < 2 {
				panic(jutil.NewIndexOutOfBoundsException(jutil.LenientFormat("Index %s out of bounds for length %s", 1, len(words))))
			}
			m.EnsureLibraryInjected(words[1], force)
		default:
			panic(jutil.NewIllegalStateExceptionMsg("Bad directive: " + directive))
		}
		node.Detach()
	}
}

// Java: RuntimeJsLibManager.setLastInjectedLibrary(Node)
func (m *RuntimeJsLibManager) SetLastInjectedLibrary(lastInjectedLibrary *rhino.Node) {
	m.lastInjectedLibrary = lastInjectedLibrary
}

// Java: RuntimeJsLibManager.getLastInjectedLibrary()
func (m *RuntimeJsLibManager) GetLastInjectedLibrary() *rhino.Node {
	return m.lastInjectedLibrary
}
