// Port of com/google/javascript/jscomp/FindExportableNodes.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// FindExportableNodes records all of the symbols and properties that should
// be exported.
//
// Currently applies to: - function foo() {} - var foo = function() {} -
// foo.bar = function() {} - var FOO = ...; - foo.BAR = ...;
//
// FOO = BAR = 5; and var FOO = BAR = 5; are not supported because the
// annotation is ambiguous to whether it applies to all the variables or only
// the first one.
type FindExportableNodes struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler

	// The set of node with @export annotations and their associated fully
	// qualified names
	exports *jutil.LinkedHashMap[string, rhino.Node]

	// The set of class member function nodes with @export annotations and
	// the fully qualified name of the owning class (not including the member
	// fn name).
	es6ClassExports *jutil.LinkedHashMap[rhino.Node, string]

	// The set of property names associated with @export annotations that do
	// not have an associated fully qualified name.
	localExports *jutil.LinkedHashSet[string]

	allowLocalExports bool
}

var FindExportableNodes_NON_GLOBAL_ERROR = DiagnosticType_Error(
	"JSC_NON_GLOBAL_ERROR",
	"@export only allowed on symbols/properties defined in the global scope.\n"+
		"Convert to a global definition or enable --export_local_property_definitions/"+
		" options.setExportLocalPropertyDefinitions(true)")

var FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED = DiagnosticType_Error(
	"JSC_EXPORT_ANNOTATION_NOT_ALLOWED", "@export is not supported on this expression.")

// Java: FindExportableNodes(AbstractCompiler, boolean)
func NewFindExportableNodes(compiler AbstractCompiler, allowLocalExports bool) *FindExportableNodes {
	f := &FindExportableNodes{}
	f.exports = jutil.NewLinkedHashMap[string, rhino.Node]()
	f.es6ClassExports = jutil.NewLinkedHashMap[rhino.Node, string]()
	f.localExports = jutil.NewLinkedHashSet[string]()
	f.compiler = compiler
	f.allowLocalExports = allowLocalExports
	return f
}

// Java: FindExportableNodes.visit(NodeTraversal, Node, Node)
func (f *FindExportableNodes) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	docInfo := n.GetJSDocInfo()
	if docInfo == nil || !docInfo.IsExport() {
		return
	}

	if parent.IsAssign() && (n.IsFunction() || n.IsClass()) {
		parentInfo := parent.GetJSDocInfo()
		if parentInfo != nil && parentInfo.IsExport() {
			// ScopedAliases produces export annotations on both the function/class
			// node and assign node, we only want to visit the assign node.
			return
		}
	}

	mode := findExportableNodes_Mode_null
	var export *string = nil
	var context rhino.Node = rhino.NilNode

	switch n.GetToken() {
	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		if parent.IsScript() {
			export = NodeUtil_GetName(n)
			context = n
			mode = FindExportableNodes_Mode_EXPORT
		}

	case rhino.Token_ASSIGN:
		grandparent := parent.GetParent()
		child := n.GetFirstChild()
		if parent.IsExprResult() && !n.GetLastChild().IsAssign() {
			if !grandparent.IsNil() && grandparent.IsScript() && child.IsQualifiedName() {
				export = child.GetQualifiedName()
				context = n
				mode = FindExportableNodes_Mode_EXPORT
			} else if f.allowLocalExports && child.IsGetProp() {
				s := child.GetString()
				export = &s
				mode = FindExportableNodes_Mode_EXTERN
			}
		}

	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		if parent.IsScript() {
			if n.GetFirstChild().HasChildren() && !n.GetFirstFirstChild().IsAssign() {
				s := n.GetFirstChild().GetString()
				export = &s
				context = n
				mode = FindExportableNodes_Mode_EXPORT
			}
		}

	case rhino.Token_GETPROP:
		{
			// TODO(b/63582201): currently, the IF body is executed even for code in the global
			// scope, e.g., when the pass looks at code transpiled from ES6 that uses getters.
			// This means that the top-level code we want to rewrite works by accident, only when
			// allowLocalExports happens to be true.
			if f.allowLocalExports && parent.IsExprResult() {
				s := n.GetString()
				export = &s
				mode = FindExportableNodes_Mode_EXTERN
			}
		}

	case rhino.Token_MEMBER_FIELD_DEF:
		if n.IsStaticMember() {
			classNode := parent.GetParent()
			className := NodeUtil_GetName(classNode)
			if className == nil {
				t.Report(n, FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED)
				return
			}
			s := *className + "." + n.GetString()
			export = &s
			mode = FindExportableNodes_Mode_EXPORT
			context = n
		} else if f.allowLocalExports {
			s := n.GetString()
			export = &s
			mode = FindExportableNodes_Mode_EXTERN
		}

	case rhino.Token_MEMBER_FUNCTION_DEF:
		if parent.GetParent().IsClass() {
			methodOwnerName :=
				NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(parent.GetParent()))
			if methodOwnerName == nil {
				t.Report(n, FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED)
				return
			}
			suffix := ".prototype"
			if n.IsStaticMember() {
				suffix = ""
			}
			f.es6ClassExports.Put(n, *methodOwnerName+suffix)
			return
		}
		fallthrough

	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		if f.allowLocalExports {
			s := n.GetString()
			export = &s
			mode = FindExportableNodes_Mode_EXTERN
		}
	default:
	}

	if export != nil {
		if mode == FindExportableNodes_Mode_EXPORT {
			rhino.CheckNotNull(context)
			f.exports.Put(*export, context)
		} else {
			jutil.CheckState(context.IsNil())
			jutil.CheckState(mode == FindExportableNodes_Mode_EXTERN)
			jutil.CheckState(!(*export == ""))
			f.localExports.Add(*export)
		}
	} else if !(n.IsGetProp() && parent.IsExprResult()) {
		// Silently ignore exports of the form:
		// /** @export */ Foo.prototype.myprop;
		// They are currently used on interfaces and records and have no effect.
		// If we cleanup the code base in the future, we can warn again.

		// Don't produce extra warnings for functions values of object literals
		if !n.IsFunction() || !NodeUtil_MayBeObjectLitKey(parent) {
			if f.allowLocalExports {
				f.compiler.Report(JSError_MakeWithNode(n, FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED))
			} else {
				f.compiler.Report(JSError_MakeWithNode(n, FindExportableNodes_NON_GLOBAL_ERROR))
			}
		}
	}
}

// Java: FindExportableNodes.getExports()
func (f *FindExportableNodes) GetExports() *jutil.LinkedHashMap[string, rhino.Node] {
	return f.exports
}

// Java: FindExportableNodes.getLocalExports()
func (f *FindExportableNodes) GetLocalExports() *jutil.LinkedHashSet[string] {
	return f.localExports
}

// Java: FindExportableNodes.getEs6ClassExports()
func (f *FindExportableNodes) GetEs6ClassExports() *jutil.LinkedHashMap[rhino.Node, string] {
	return f.es6ClassExports
}

// FindExportableNodes_Mode is the Mode enum.
type FindExportableNodes_Mode int32

const (
	FindExportableNodes_Mode_EXPORT FindExportableNodes_Mode = iota
	FindExportableNodes_Mode_EXTERN
)

// findExportableNodes_Mode_null stands for a null Mode variable.
//
// Go-only: language-adapter — Java null for an enum-typed local.
const findExportableNodes_Mode_null FindExportableNodes_Mode = -1

var findExportableNodes_Mode_names = [...]string{"EXPORT", "EXTERN"}

// FindExportableNodes_Mode_values is Mode.values().
func FindExportableNodes_Mode_values() []FindExportableNodes_Mode {
	return []FindExportableNodes_Mode{FindExportableNodes_Mode_EXPORT, FindExportableNodes_Mode_EXTERN}
}

// Name is Enum.name().
func (m FindExportableNodes_Mode) Name() string { return findExportableNodes_Mode_names[m] }

// String is Enum.toString().
func (m FindExportableNodes_Mode) String() string { return m.Name() }

// Ordinal is Enum.ordinal().
func (m FindExportableNodes_Mode) Ordinal() int { return int(m) }
