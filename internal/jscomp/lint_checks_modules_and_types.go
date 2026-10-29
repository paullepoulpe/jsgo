// Port of com/google/javascript/jscomp/lint/CheckPrototypeProperties.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/lint/CheckEs6Modules.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/lint/CheckNoMutatedEs6Exports.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/lint/CheckPrimitiveAsObject.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/lint/CheckGoogModuleTypeScriptName.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/lint/CheckArrayWithGoogObject.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/lint/CheckConstantCaseNames.java @ closure-compiler v20261021

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

/*
 * Copyright 2016 The Closure Compiler Authors.
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
 * Copyright 2021 The Closure Compiler Authors.
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
 * Copyright 2015 The Closure Compiler Authors.
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
 * Copyright 2019 The Closure Compiler Authors.
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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ---------------------------------------------------------------------------
// CheckPrototypeProperties
// ---------------------------------------------------------------------------

// CheckPrototypeProperties checks when a mutable property is assigned to a
// prototype. This is generally undesirable because it can lead to the
// following unexpected situation.
//
//	/** @constructor * /
//	function MyClass() {}
//	MyClass.prototype.prop = [];
//	x = new MyClass;
//	y = new MyClass;
//	x.prop.push(1);
//	console.log(y.prop) // [1]
type CheckPrototypeProperties struct {
	compiler AbstractCompiler
}

var CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER = DiagnosticType_Disabled(
	"JSC_ILLEGAL_PROTOTYPE_MEMBER",
	"Prototype property {0} should be a primitive, not an Array or Object.")

// NewCheckPrototypeProperties is new CheckPrototypeProperties(compiler).
//
// Java: CheckPrototypeProperties(AbstractCompiler)
func NewCheckPrototypeProperties(compiler AbstractCompiler) *CheckPrototypeProperties {
	return &CheckPrototypeProperties{compiler: compiler}
}

func (*CheckPrototypeProperties) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewCheckPrototypeProperties(c)
}

// Java: CheckPrototypeProperties.process(Node, Node)
func (c *CheckPrototypeProperties) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

var checkPrototypePropertiesVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_EXPR_RESULT)
	return s
}()

var (
	_ nodeTraversal_visitTokens   = (*CheckPrototypeProperties)(nil)
	_ nodeTraversal_noVisitLeaves = (*CheckPrototypeProperties)(nil)
)

func (*CheckPrototypeProperties) noVisitLeaves() {}
func (*CheckPrototypeProperties) visitTokens() *ast.TokenSet {
	return &checkPrototypePropertiesVisitTokens
}

// Java: CheckPrototypeProperties.visit(NodeTraversal, Node, Node)
func (c *CheckPrototypeProperties) Visit(unused *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !NodeUtil_IsPrototypePropertyDeclaration(n) {
		return
	}
	exprStmt, ok := ast.AsExprStmt(n)
	if !ok {
		return
	}
	assign, ok := exprStmt.Expr().AsAssign()
	if !ok {
		return
	}
	rhs := assign.Value().Node()
	if rhs.IsArrayLit() || rhs.IsObjectLit() {
		jsDoc := NodeUtil_GetBestJSDocInfo(rhs)
		if jsDoc != nil && jsDoc.HasEnumParameterType() {
			// Don't report for @enum's on the prototype. Sometimes this is necessary, for example,
			// to expose the enum values to an Angular template.
			return
		}
		// Note: Java calls assign.getFirstChild().getLastChild().getString(), which on a
		// single-child GETPROP (MyClass.prototype.prop) reads the receiver's string ("prototype").
		propName := assign.Target().Node().GetLastChild().GetString()
		c.compiler.Report(JSError_MakeWithNode(assign.Node(), CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER, propName))
	}
}

// Java: CheckPrototypeProperties.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckPrototypeProperties) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

// ---------------------------------------------------------------------------
// CheckEs6Modules
// ---------------------------------------------------------------------------

// CheckEs6Modules: miscellaneous checks for style in ES6 modules.
type CheckEs6Modules struct {
	compiler         AbstractCompiler
	importSpecifiers map[uint32]rhino.Node
}

var CheckEs6Modules_DUPLICATE_IMPORT = DiagnosticType_Warning(
	"JSC_DUPLICATE_IMPORT", "The module \"{0}\" has already been imported at {1}, {2}.")

var CheckEs6Modules_NO_DEFAULT_EXPORT = DiagnosticType_Warning(
	"JSC_DEFAULT_EXPORT",
	"Do not use the default export. There is no way to force consistent naming when "+
		"imported.")

// NewCheckEs6Modules is new CheckEs6Modules(compiler).
//
// Java: CheckEs6Modules(AbstractCompiler)
func NewCheckEs6Modules(compiler AbstractCompiler) *CheckEs6Modules {
	return &CheckEs6Modules{
		compiler: compiler,
	}
}

func (*CheckEs6Modules) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewCheckEs6Modules(c)
}

var checkEs6ModulesVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_MODULE_BODY)
	return s
}()

var (
	_ nodeTraversal_visitTokens   = (*CheckEs6Modules)(nil)
	_ nodeTraversal_noVisitLeaves = (*CheckEs6Modules)(nil)
)

func (*CheckEs6Modules) noVisitLeaves()             {}
func (*CheckEs6Modules) visitTokens() *ast.TokenSet { return &checkEs6ModulesVisitTokens }

// Java: CheckEs6Modules.visit(NodeTraversal, Node, Node)
func (c *CheckEs6Modules) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if _, ok := ast.AsModuleBody(n); ok {
		clear(c.importSpecifiers)
	}
}

// Java: CheckEs6Modules.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckEs6Modules) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_MODULE_BODY:
		return true
	case rhino.Token_SCRIPT:
		return n.GetBooleanProp(rhino.Node_ES6_MODULE)
	case rhino.Token_IMPORT:
		c.VisitImport(t, n)
		return false
	case rhino.Token_EXPORT:
		c.VisitExport(t, n)
		return false
	default:
		return false
	}
}

// Java: CheckEs6Modules.visitImport(NodeTraversal, Node)
func (c *CheckEs6Modules) VisitImport(t *NodeTraversal, importNode rhino.Node) {
	src := ast.MustImport(importNode).Source()
	specifierID := src.Node().GetStringID()
	if duplicateImport, ok := c.importSpecifiers[specifierID]; ok {
		t.Report(
			importNode,
			CheckEs6Modules_DUPLICATE_IMPORT,
			src.Value(),
			strconv.Itoa(duplicateImport.GetLineno()),
			strconv.Itoa(duplicateImport.GetCharno()))
	} else {
		if c.importSpecifiers == nil {
			c.importSpecifiers = make(map[uint32]rhino.Node)
		}
		c.importSpecifiers[specifierID] = importNode
	}
}

// Java: CheckEs6Modules.visitExport(NodeTraversal, Node)
func (c *CheckEs6Modules) VisitExport(t *NodeTraversal, export rhino.Node) {
	if ast.MustExport(export).IsDefault() {
		t.Report(export, CheckEs6Modules_NO_DEFAULT_EXPORT)
	}
}

// Java: CheckEs6Modules.process(Node, Node)
func (c *CheckEs6Modules) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// ---------------------------------------------------------------------------
// CheckPrimitiveAsObject
// ---------------------------------------------------------------------------

// CheckPrimitiveAsObject checks for explicit creation of the object
// equivalents of primitive types (e.g. Boolean instead of boolean) and their
// use in type declarations.
//
// Using these is confusing and gives no benefit. For example, the result of
// typeof (new Boolean(true)) is "object". and the result of (new
// Boolean(false)) ? "true" : "false" is "true".
type CheckPrimitiveAsObject struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

var CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT = DiagnosticType_Warning("JSC_PRIMITIVE_OBJECT", "Explicit creation of a {0} object.")

var CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION = DiagnosticType_Warning(
	"JSC_PRIMITIVE_OBJECT_DECLARATION",
	"Declaration of {0} object instead of primitive type.")

func checkPrimitiveAsObject_isPrimitiveObjectConstructor(name string) bool {
	switch name {
	case "Boolean", "Number", "String":
		return true
	default:
		return false
	}
}

// NewCheckPrimitiveAsObject is new CheckPrimitiveAsObject(compiler).
//
// Java: CheckPrimitiveAsObject(AbstractCompiler)
func NewCheckPrimitiveAsObject(compiler AbstractCompiler) *CheckPrimitiveAsObject {
	return &CheckPrimitiveAsObject{compiler: compiler}
}

func (*CheckPrimitiveAsObject) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewCheckPrimitiveAsObject(c)
}

// Java: CheckPrimitiveAsObject.process(Node, Node)
func (c *CheckPrimitiveAsObject) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckPrimitiveAsObject.visit(NodeTraversal, Node, Node)
func (c *CheckPrimitiveAsObject) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsNew() {
		c.CheckForPrimitiveObjectConstructor(t, n)
	}
	if n.GetJSDocInfo() != nil {
		c.CheckForPrimitiveObjectDeclaration(t, n)
	}
}

// Java: CheckPrimitiveAsObject.checkForPrimitiveObjectDeclaration(NodeTraversal, Node)
func (c *CheckPrimitiveAsObject) CheckForPrimitiveObjectDeclaration(t *NodeTraversal, n rhino.Node) {
	jsDocInfo := n.GetJSDocInfo()
	if jsDocInfo != nil {
		for it := jsDocInfo.TypeExprIter(); it.Next(); {
			c.CheckTypeNodeForPrimitiveObjectDeclaration(t, it.Node())
		}
	}
}

// Java: CheckPrimitiveAsObject.checkTypeNodeForPrimitiveObjectDeclaration(NodeTraversal, Node)
func (c *CheckPrimitiveAsObject) CheckTypeNodeForPrimitiveObjectDeclaration(t *NodeTraversal, typeRoot rhino.Node) {
	NodeUtil_VisitPreOrder(
		typeRoot,
		func(node rhino.Node) {
			if strLit, ok := ast.AsStringLit(node); ok {
				typeName := strLit.Value()
				if checkPrimitiveAsObject_isPrimitiveObjectConstructor(typeName) {
					t.Report(node, CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION, typeName)
				}
			}
		})
}

// Java: CheckPrimitiveAsObject.checkForPrimitiveObjectConstructor(NodeTraversal, Node)
func (c *CheckPrimitiveAsObject) CheckForPrimitiveObjectConstructor(t *NodeTraversal, n rhino.Node) {
	if newExpr, ok := ast.AsNew(n); ok {
		if calleeName, ok := newExpr.Callee().AsName(); ok {
			constructorName := calleeName.Ident()
			if checkPrimitiveAsObject_isPrimitiveObjectConstructor(constructorName) {
				t.Report(n, CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT, constructorName)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// CheckGoogModuleTypeScriptName
// ---------------------------------------------------------------------------

// CheckGoogModuleTypeScriptName checks that goog.module statement matches the
// generated TypeScript module namespace, which is based on the file path.
type CheckGoogModuleTypeScriptName struct {
	compiler AbstractCompiler

	finished bool
}

var CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE = DiagnosticType_Disabled(
	"JSC_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE",
	"goog.module namespace does not match the future TypeScript namespace, which is generated"+
		" from the file path."+
		" The correct namespace is: \"{0}\"")

// MOE::begin_strip
var checkGoogModuleTypeScriptName_allowedDirectories = [...]string{
	"google3/gws/",
	"google3/java/com/google/gws/",
	"google3/javascript/search/",
}

// MOE::end_strip

// NewCheckGoogModuleTypeScriptName is new CheckGoogModuleTypeScriptName(compiler).
//
// Java: CheckGoogModuleTypeScriptName(AbstractCompiler)
func NewCheckGoogModuleTypeScriptName(compiler AbstractCompiler) *CheckGoogModuleTypeScriptName {
	return &CheckGoogModuleTypeScriptName{compiler: compiler}
}

// Java: CheckGoogModuleTypeScriptName.process(Node, Node)
func (c *CheckGoogModuleTypeScriptName) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckGoogModuleTypeScriptName.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckGoogModuleTypeScriptName) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Traverse top-level statements until a `goog.module` statement is found.
	return !c.finished &&
		(parent.IsNil() || parent.IsRoot() || parent.IsScript() || parent.IsModuleBody())
}

var checkGoogModuleTypeScriptNameVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_EXPR_RESULT)
	return s
}()

var (
	_ nodeTraversal_visitTokens   = (*CheckGoogModuleTypeScriptName)(nil)
	_ nodeTraversal_noVisitLeaves = (*CheckGoogModuleTypeScriptName)(nil)
)

func (*CheckGoogModuleTypeScriptName) noVisitLeaves() {}
func (*CheckGoogModuleTypeScriptName) visitTokens() *ast.TokenSet {
	return &checkGoogModuleTypeScriptNameVisitTokens
}

// Java: CheckGoogModuleTypeScriptName.visit(NodeTraversal, Node, Node)
func (c *CheckGoogModuleTypeScriptName) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if NodeUtil_IsGoogModuleCall(n) {
		c.CheckGoogModuleNamespace(t, n)
		c.finished = true
	}
}

// Java: CheckGoogModuleTypeScriptName.checkGoogModuleNamespace(NodeTraversal, Node)
func (c *CheckGoogModuleTypeScriptName) CheckGoogModuleNamespace(t *NodeTraversal, n rhino.Node) {
	originalNamespace := n.GetFirstChild().GetSecondChild().GetString()
	sourceNamePtr := t.GetSourceName()
	// Go-only: language-adapter — Node.getString() is a Go string, which
	// cannot be null, so `originalNamespace == null` is always false.
	if sourceNamePtr == nil {
		return
	}
	sourceName := *sourceNamePtr
	// MOE::begin_strip
	google3Index := strings.Index(sourceName, "google3")
	if google3Index == -1 {
		sourceName = "google3/" + sourceName
	} else if google3Index != 0 {
		sourceName = sourceName[google3Index:]
	}
	// MOE::end_strip
	replaced := strings.ReplaceAll(sourceName, "/", ".")
	var replacementNamespace string
	if nBytes := len(replaced); nBytes >= len(".js") &&
		replaced[nBytes-3] < 0x80 && replaced[nBytes-2] < 0x80 && replaced[nBytes-1] < 0x80 {
		replacementNamespace = replaced[:nBytes-len(".js")]
	} else {
		replacementNamespace = jutil.StringSubstring(
			replaced, 0, jutil.StringLength(sourceName)-len(".js"))
	}
	if originalNamespace != replacementNamespace {
		for _, allowedDirectory := range checkGoogModuleTypeScriptName_allowedDirectories {
			if strings.HasPrefix(sourceName, allowedDirectory) {
				t.Report(n, CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE, replacementNamespace)
				return
			}
		}
	}
}

// ---------------------------------------------------------------------------
// CheckArrayWithGoogObject
// ---------------------------------------------------------------------------

// CheckArrayWithGoogObject lints against passing arrays to goog.object
// methods with the intention of iterating over them as though with a for-in
// loop, which is discouraged with arrays.
type CheckArrayWithGoogObject struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

func checkArrayWithGoogObject_isGoogObjectMethod(name string) bool {
	switch name {
	case "goog.object.forEach",
		"goog.object.filter",
		"goog.object.map",
		"goog.object.some",
		"goog.object.every",
		"goog.object.getCount",
		"goog.object.getAnyKey",
		"goog.object.getAnyValue",
		"goog.object.contains",
		"goog.object.getValues",
		"goog.object.getKeys",
		"goog.object.findKey",
		"goog.object.findValue",
		"goog.object.isEmpty",
		"goog.object.clear",
		"goog.object.remove",
		"goog.object.equals",
		"goog.object.clone",
		"goog.object.transpose":
		return true
	default:
		return false
	}
}

var CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT = DiagnosticType_Warning(
	"JSC_ARRAY_PASSED_TO_GOOG_OBJECT",
	"{0} expects an object, not an array. Did you mean to use goog.array?")

// NewCheckArrayWithGoogObject is new CheckArrayWithGoogObject(compiler).
//
// Java: CheckArrayWithGoogObject(AbstractCompiler)
func NewCheckArrayWithGoogObject(compiler AbstractCompiler) *CheckArrayWithGoogObject {
	return &CheckArrayWithGoogObject{compiler: compiler}
}

func (*CheckArrayWithGoogObject) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewCheckArrayWithGoogObject(c)
}

func (c *CheckArrayWithGoogObject) googObjectArrayCallName(n rhino.Node) (string, bool) {
	call, ok := ast.AsCall(n)
	if !ok {
		return "", false
	}
	callee := call.Callee().Node()
	if !callee.IsQualifiedName() {
		return "", false
	}
	name := *callee.GetQualifiedName()
	if !checkArrayWithGoogObject_isGoogObjectMethod(name) {
		return "", false
	}
	firstArg, hasArg := call.FirstArg()
	if !hasArg {
		return "", false
	}
	type_ := firstArg.Node().GetJSType()
	if type_ == nil || !c.ContainsArray(type_) {
		return "", false
	}
	return name, true
}

// Java: CheckArrayWithGoogObject.isGoogObjectIterationOverArray(Node)
func (c *CheckArrayWithGoogObject) IsGoogObjectIterationOverArray(n rhino.Node) bool {
	_, ok := c.googObjectArrayCallName(n)
	return ok
}

// Java: CheckArrayWithGoogObject.containsArray(JSType)
func (c *CheckArrayWithGoogObject) ContainsArray(type_ rhino.JSType) bool {
	// Check if type is itself an array
	if type_.IsArrayType() {
		return true
	}
	templatizedType := type_.ToMaybeTemplatizedType()
	if templatizedType != nil && templatizedType.GetReferencedType().IsArrayType() {
		return true
	}

	// Check if this is a union that contains an array
	if type_.IsUnionType() {
		arrayType := c.compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)
		for _, alternate := range type_.ToMaybeUnionType().GetAlternates() {
			if alternate.IsSubtypeOf(arrayType) {
				return true
			}
		}
	}
	return false
}

var checkArrayWithGoogObjectVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_CALL)
	return s
}()

var (
	_ nodeTraversal_visitTokens   = (*CheckArrayWithGoogObject)(nil)
	_ nodeTraversal_noVisitLeaves = (*CheckArrayWithGoogObject)(nil)
)

func (*CheckArrayWithGoogObject) noVisitLeaves() {}
func (*CheckArrayWithGoogObject) visitTokens() *ast.TokenSet {
	return &checkArrayWithGoogObjectVisitTokens
}

// Java: CheckArrayWithGoogObject.visit(NodeTraversal, Node, Node)
func (c *CheckArrayWithGoogObject) Visit(unused *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if name, ok := c.googObjectArrayCallName(n); ok {
		c.compiler.Report(
			JSError_MakeWithNode(n, CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT, name))
	}
}

// Java: CheckArrayWithGoogObject.process(Node, Node)
func (c *CheckArrayWithGoogObject) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// ---------------------------------------------------------------------------
// CheckConstantCaseNames
// ---------------------------------------------------------------------------

// CheckConstantCaseNames looks for module-level variable declarations that
// use CONSTANT_CASE, according to the Google style guide, and verifies that
// they are also annotated @const or are in a const clause.
//
// This pass could be extended to check CONSTANT_CASE properties in the
// future.
//
// Non-module-level variables should always use camel case according to the
// Google style guide. In order to not confuse users, this pass does not warn
// that they should be @const. (A more correct lint check could warn that
// non-module-locals should not be constant case.)
type CheckConstantCaseNames struct {
	compiler   AbstractCompiler
	convention CodingConvention
	// Maps CONSTANT_CASE module-level name IDs to their initializing NAME node,
	// accompanied by invalidNamesOrder to preserve declaration order.
	invalidNamesPerModule map[uint32]rhino.Node
	invalidNamesOrder     []uint32
	// Subset of variables in `invalidNamesPerModule` that are mutated post-declaration.
	reassignedNames map[uint32]bool
}

var CheckConstantCaseNames_MISSING_CONST_PROPERTY = DiagnosticType_Disabled(
	"JSC_MISSING_CONST_ON_CONSTANT_CASE",
	"CONSTANT_CASE name \"{0}\" is constant-by-convention, so must be explicitly"+
		" `const` or @const")

var CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME = DiagnosticType_Disabled(
	"JSC_REASSIGNED_CONSTANT_CASE_NAME",
	"CONSTANT_CASE name \"{0}\" is constant-by-convention but is reassigned. "+
		"Use camelCase instead.")

// NewCheckConstantCaseNames is new CheckConstantCaseNames(compiler).
//
// Java: CheckConstantCaseNames(AbstractCompiler)
func NewCheckConstantCaseNames(compiler AbstractCompiler) *CheckConstantCaseNames {
	return &CheckConstantCaseNames{
		compiler:   compiler,
		convention: compiler.GetCodingConvention(),
	}
}

func (*CheckConstantCaseNames) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewCheckConstantCaseNames(c)
}

// Java: CheckConstantCaseNames.process(Node, Node)
func (c *CheckConstantCaseNames) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckConstantCaseNames.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckConstantCaseNames) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Only need to warn for module-level names, so don't visit other files.
	if script, ok := ast.AsScript(n); ok {
		first, hasFirst := script.FirstItem()
		if !hasFirst {
			return false
		}
		_, isModule := first.AsModuleBody()
		return isModule
	}
	return true
}

var checkConstantCaseNamesVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_MODULE_BODY)
	s.Add(rhino.Token_VAR)
	s.Add(rhino.Token_LET)
	s.Add(rhino.Token_NAME)
	return s
}()

var (
	_ nodeTraversal_visitTokens    = (*CheckConstantCaseNames)(nil)
	_ nodeTraversal_onlyNameLeaves = (*CheckConstantCaseNames)(nil)
)

func (*CheckConstantCaseNames) onlyNameLeaves()            {}
func (*CheckConstantCaseNames) visitTokens() *ast.TokenSet { return &checkConstantCaseNamesVisitTokens }

// Java: CheckConstantCaseNames.visit(NodeTraversal, Node, Node)
func (c *CheckConstantCaseNames) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if _, ok := ast.AsModuleBody(n); ok {
		c.ReportWarningsAndClear()
		return
	}
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET:
		// Skip CONST as it automatically meets the criteria and only look for module-level vars.
		if !t.InModuleScope() {
			return
		}
		info := n.GetJSDocInfo()
		if info != nil && info.HasConstAnnotation() {
			break
		}
		NodeUtil_VisitLhsNodesInNode(
			n,
			func(name rhino.Node) {
				ident := name.GetString()
				if c.convention.IsConstant(ident) {
					id := name.GetStringID()
					if c.invalidNamesPerModule == nil {
						c.invalidNamesPerModule = make(map[uint32]rhino.Node)
					}
					if _, exists := c.invalidNamesPerModule[id]; !exists {
						c.invalidNamesOrder = append(c.invalidNamesOrder, id)
					}
					c.invalidNamesPerModule[id] = name
				}
			})
	case rhino.Token_NAME:
		if len(c.invalidNamesPerModule) == 0 {
			return
		}
		id := n.GetStringID()
		if _, exists := c.invalidNamesPerModule[id]; !exists {
			return
		}
		if !NodeUtil_IsLValue(n) {
			return
		}
		// Verify this name is referring to the actual module-level var and not a local shadow.
		v := t.GetScope().GetVarByID(id)
		if v.GetScopeRoot().IsModuleBody() {
			if c.reassignedNames == nil {
				c.reassignedNames = make(map[uint32]bool)
			}
			c.reassignedNames[id] = true
		}
	default:
	}
}

// Java: CheckConstantCaseNames.reportWarningsAndClear()
func (c *CheckConstantCaseNames) ReportWarningsAndClear() {
	for _, id := range c.invalidNamesOrder {
		nameNode := c.invalidNamesPerModule[id]
		name := rhino.RhinoStringPool_StringFromID(id)
		if c.reassignedNames[id] {
			c.compiler.Report(JSError_MakeWithNode(nameNode, CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME, name))
		} else {
			c.compiler.Report(JSError_MakeWithNode(nameNode, CheckConstantCaseNames_MISSING_CONST_PROPERTY, name))
		}
	}
	clear(c.invalidNamesPerModule)
	c.invalidNamesOrder = c.invalidNamesOrder[:0]
	clear(c.reassignedNames)
}

// ---------------------------------------------------------------------------
// CheckNoMutatedEs6Exports
// ---------------------------------------------------------------------------

// CheckNoMutatedEs6Exports checks that exports of ES6 modules are not mutated
// outside of module initialization.
type CheckNoMutatedEs6Exports struct {
	compiler           AbstractCompiler
	mutatedNames       map[uint32][]rhino.Node
	mutatedOrder       []uint32
	exportedLocalNames map[uint32]struct{}
}

var CheckNoMutatedEs6Exports_MUTATED_EXPORT = DiagnosticType_Warning(
	"JSC_MUTATED_EXPORT",
	"The name \"{0}\" is exported and should not be mutated outside of module "+
		"initialization. Mutable exports are generally difficult to reason about. You "+
		"can work around this by exporting getter/setter functions, or an object with "+
		"mutable properties instead.")

func NewCheckNoMutatedEs6Exports(compiler AbstractCompiler) *CheckNoMutatedEs6Exports {
	return &CheckNoMutatedEs6Exports{
		compiler: compiler,
	}
}

func (*CheckNoMutatedEs6Exports) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewCheckNoMutatedEs6Exports(c)
}

var checkNoMutatedEs6ExportsVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_MODULE_BODY)
	return s
}()

var (
	_ nodeTraversal_visitTokens   = (*CheckNoMutatedEs6Exports)(nil)
	_ nodeTraversal_noVisitLeaves = (*CheckNoMutatedEs6Exports)(nil)
)

func (*CheckNoMutatedEs6Exports) noVisitLeaves() {}
func (*CheckNoMutatedEs6Exports) visitTokens() *ast.TokenSet {
	return &checkNoMutatedEs6ExportsVisitTokens
}

func (c *CheckNoMutatedEs6Exports) CheckNoMutations() {
	for _, mutatedID := range c.mutatedOrder {
		if _, ok := c.exportedLocalNames[mutatedID]; !ok {
			continue
		}
		mutatedExport := rhino.RhinoStringPool_StringFromID(mutatedID)
		for _, mutation := range c.mutatedNames[mutatedID] {
			c.compiler.Report(JSError_MakeWithNode(mutation, CheckNoMutatedEs6Exports_MUTATED_EXPORT, mutatedExport))
		}
	}

	clear(c.mutatedNames)
	c.mutatedOrder = c.mutatedOrder[:0]
	clear(c.exportedLocalNames)
}

func (c *CheckNoMutatedEs6Exports) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsModuleBody() {
		c.CheckNoMutations()
	}
}

func (c *CheckNoMutatedEs6Exports) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		return n.GetBooleanProp(rhino.Node_ES6_MODULE)
	case rhino.Token_EXPORT:
		c.VisitExport(n)
	case rhino.Token_NAME:
		c.VisitName(t, n)
	}
	return true
}

func (c *CheckNoMutatedEs6Exports) addExportedLocalName(id uint32) {
	if c.exportedLocalNames == nil {
		c.exportedLocalNames = make(map[uint32]struct{})
	}
	c.exportedLocalNames[id] = struct{}{}
}

func (c *CheckNoMutatedEs6Exports) VisitExport(export rhino.Node) {
	if !export.HasOneChild() {
		return
	}
	first := export.GetFirstChild()
	if first.IsExportSpecs() {
		for exportSpec := first.GetFirstChild(); !exportSpec.IsNil(); exportSpec = exportSpec.GetNext() {
			jutil.CheckState(exportSpec.HasTwoChildren())
			c.addExportedLocalName(exportSpec.GetFirstChild().GetStringID())
		}
		return
	}
	if export.GetBooleanProp(rhino.Node_EXPORT_ALL_FROM) {
		return
	}
	if NodeUtil_IsNameDeclaration(first) {
		NodeUtil_VisitLhsNodesInNode(first, func(lhs rhino.Node) {
			jutil.CheckState(lhs.IsName())
			c.addExportedLocalName(lhs.GetStringID())
		})
	} else if export.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
		if !first.IsClass() && !first.IsFunction() {
			return
		}
		nameNode := first.GetFirstChild()
		if !nameNode.IsEmpty() && nameNode.GetStringID() > rhino.RhinoStringPool_EmptyID {
			c.addExportedLocalName(nameNode.GetStringID())
		}
	} else {
		jutil.CheckState(first.IsClass() || first.IsFunction())
		c.addExportedLocalName(first.GetFirstChild().GetStringID())
	}
}

func (c *CheckNoMutatedEs6Exports) VisitName(t *NodeTraversal, name rhino.Node) {
	if !NodeUtil_IsLValue(name) {
		return
	}
	scope := t.GetScope()
	if scope.GetClosestHoistScope().IsModuleScope() {
		return
	}
	id := name.GetStringID()
	var_ := scope.GetVarByID(id)
	if var_ == nil || !var_.scope.IsModuleScope() {
		return
	}
	if c.mutatedNames == nil {
		c.mutatedNames = make(map[uint32][]rhino.Node)
	}
	nodes, exists := c.mutatedNames[id]
	if !exists {
		c.mutatedOrder = append(c.mutatedOrder, id)
	}
	if !slices.Contains(nodes, name) {
		c.mutatedNames[id] = append(nodes, name)
	}
}

func (c *CheckNoMutatedEs6Exports) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}
