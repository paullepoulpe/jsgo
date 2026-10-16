// Port of com/google/javascript/jscomp/ExternExportsPass.java @ closure-compiler v20260415

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ExternExportsPass creates an externs file containing all exported symbols
// and properties for later consumption.
type ExternExportsPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	// The exports found.
	exports *jutil.ArrayList[ExternExportsPass_Export]

	// A map of all assigns to their parent nodes.
	definitionMap *jutil.LinkedHashMap[string, *rhino.Node]

	// The parent compiler.
	compiler AbstractCompiler

	// The AST root which holds the externs generated.
	externsRoot *rhino.Node

	// A mapping of internal paths to exported paths.
	mappedPaths *jutil.LinkedHashMap[string, string]

	// A list of exported paths.
	alreadyExportedPaths *jutil.LinkedHashSet[string]

	// A list of function names used to export symbols.
	exportSymbolFunctionNames *jutil.LinkedHashSet[string]

	// A list of function names used to export properties.
	exportPropertyFunctionNames *jutil.LinkedHashSet[string]
}

var _ CompilerPass = (*ExternExportsPass)(nil)
var _ NodeTraversal_Callback = (*ExternExportsPass)(nil)

// Q_NAME_JOINER is Joiner.on('.').
func externExportsPass_qNameJoin(a string, b string) string {
	return a + "." + b
}

// ExternExportsPass_Export is ExternExportsPass.Export (abstract, inner).
type ExternExportsPass_Export interface {
	// Java: ExternExportsPass.Export.generateExterns()
	GenerateExterns()

	// GetExportedPath returns the path exported by this export.
	//
	// Java: ExternExportsPass.Export.getExportedPath()
	GetExportedPath() string

	GetExternExportsPass_ExportBase() *ExternExportsPass_ExportBase
}

// ExternExportsPass_ExportBase holds the fields and base methods of Export.
type ExternExportsPass_ExportBase struct {
	self  ExternExportsPass_Export
	outer *ExternExportsPass

	symbolName string
	value      *rhino.Node
}

// Java: ExternExportsPass.Export(String, Node)
func (e *ExternExportsPass_ExportBase) InitExternExportsPass_ExportBase(
	self ExternExportsPass_Export, outer *ExternExportsPass, symbolName string, value *rhino.Node) {
	e.self = self
	e.outer = outer
	e.symbolName = symbolName
	e.value = jutil.CheckNotNull(value)
}

// GetExternExportsPass_ExportBase gives access to the base fields.
func (e *ExternExportsPass_ExportBase) GetExternExportsPass_ExportBase() *ExternExportsPass_ExportBase {
	return e
}

// GenerateExterns generates the externs representation of this export and
// appends it to the externsRoot AST.
//
// Java: ExternExportsPass.Export.generateExterns()
func (e *ExternExportsPass_ExportBase) GenerateExterns() {
	e.AppendExtern(e.self.GetExportedPath(), e.GetValue())
}

// AppendExtern appends the exported function and all paths necessary for
// the path to be declared. For example, for a property "a.b.c", the
// initializers for paths "a", "a.b" will be appended (if they have not
// already) and a.b.c will be initialized with the exported version of the
// function:
//
//	var a = {};
//	a.b = {};
//	a.b.c = function(x,y) { }
//
// Java: ExternExportsPass.Export.appendExtern(String, Node)
func (e *ExternExportsPass_ExportBase) AppendExtern(path string, valueToExport *rhino.Node) {
	pathPrefixes := externExportsPass_computePathPrefixes(path)

	for i := 0; i < len(pathPrefixes); i++ {
		pathPrefix := pathPrefixes[i]

		// The complete path (the last path prefix) must be emitted and
		// it gets initialized to the externed version of the value.
		isCompletePathPrefix := i == len(pathPrefixes)-1

		skipPathPrefix :=
			strings.HasSuffix(pathPrefix, ".prototype") ||
				(e.outer.alreadyExportedPaths.Contains(pathPrefix) && !isCompletePathPrefix)
		if skipPathPrefix {
			continue
		}

		exportedValueDefinesNewType := false

		if valueToExport != nil {
			jsdoc := NodeUtil_GetBestJSDocInfo(valueToExport)
			if valueToExport.IsClass() || (jsdoc != nil && jsdoc.ContainsTypeDefinition()) {
				exportedValueDefinesNewType = true
			}
		}

		// Namespaces get initialized to {}, functions to externed versions of
		// their value, and if we can't figure out where the value came from we
		// initialize it to {}.
		//
		// Since externs are always exported in sorted order, we know that if we
		// export a.b = function() {} and later a.b.c = function then a.b will
		// always be in alreadyExportedPaths when we emit a.b.c and thus we will
		// never overwrite the function exported for a.b with a namespace.
		var initializer *rhino.Node
		var jsdoc *rhino.JSDocInfo
		if isCompletePathPrefix && valueToExport != nil {
			if valueToExport.IsFunction() {
				initializer = e.CreateExternFunction(valueToExport)
			} else if valueToExport.IsClass() {
				initializer = e.CreateExternFunctionForEs6Class(valueToExport)
			} else {
				jutil.CheckState(valueToExport.IsObjectLit())
				initializer = e.CreateExternObjectLit(valueToExport)
			}
		} else if !isCompletePathPrefix && exportedValueDefinesNewType {
			jsdoc = e.BuildNamespaceJSDoc()
			initializer = e.CreateExternObjectLit(rhino.IR_Objectlit())
			// Don't add the empty jsdoc here
			initializer.SetJSDocInfo(nil)
		} else {
			initializer = rhino.IR_Empty()
		}

		e.AppendPathDefinition(pathPrefix, initializer, jsdoc)
	}
}

// Java: ExternExportsPass.Export.appendPathDefinition(String, Node, JSDocInfo)
func (e *ExternExportsPass_ExportBase) AppendPathDefinition(
	path string, initializer *rhino.Node, jsdoc *rhino.JSDocInfo) {
	var pathDefinition *rhino.Node

	if strings.Contains(path, ".") {
		qualifiedPath := NodeUtil_NewQName(e.outer.compiler, path)
		if initializer.IsEmpty() {
			pathDefinition = NodeUtil_NewExpr(qualifiedPath)
		} else {
			pathDefinition = NodeUtil_NewExpr(rhino.IR_Assign(qualifiedPath, initializer))
		}
	} else {
		if initializer.IsEmpty() {
			pathDefinition = rhino.IR_Var(rhino.IR_Name(path))
		} else {
			pathDefinition = NodeUtil_NewVarNode(path, initializer)
		}
	}

	if jsdoc != nil {
		if pathDefinition.IsExprResult() {
			pathDefinition.GetFirstChild().SetJSDocInfo(jsdoc)
		} else {
			jutil.CheckState(pathDefinition.IsVar())
			pathDefinition.SetJSDocInfo(jsdoc)
		}
	}

	e.outer.externsRoot.AddChildToBack(pathDefinition)

	e.outer.alreadyExportedPaths.Add(path)
}

// CreateExternFunction: given a function to export, create the empty
// function that will be put in the externs file. This extern function should
// have the same type as the original function and the same parameter name
// but no function body.
//
// We create a warning here if the the function to export is missing
// parameter or return types.
//
// Java: ExternExportsPass.Export.createExternFunction(Node)
func (e *ExternExportsPass_ExportBase) CreateExternFunction(exportedFunction *rhino.Node) *rhino.Node {
	paramList := e.CreateExternsParamListFromOriginalFunction(exportedFunction)
	externFunction := rhino.IR_Function(rhino.IR_Name(""), paramList, rhino.IR_Block())

	externFunction.SetJSType(exportedFunction.GetJSType())

	return externFunction
}

// CreateExternsParamListFromOriginalFunction creates a PARAM_LIST to store
// in the AST we'll use to generate externs for a function with the given
// type.
//
// If the NODE defining the original function is available, it would be
// better to use createExternsParamListFromOriginalFunction(), because that
// one will keep the parameter names the same instead of generating arbitrary
// parameter names.
//
// exportedFunction is the FUNCTION Node of the original function.
//
// Java: ExternExportsPass.Export.createExternsParamListFromOriginalFunction(Node)
func (e *ExternExportsPass_ExportBase) CreateExternsParamListFromOriginalFunction(exportedFunction *rhino.Node) *rhino.Node {
	originalParamList := NodeUtil_GetFunctionParameters(exportedFunction)
	// First get all of the original positional parameter list names we can.
	// Place empty stings in the positions where we'll need to generate names.
	originalParamNames := []string{}
	for originalParam := originalParamList.GetFirstChild(); originalParam != nil; originalParam = originalParam.GetNext() {
		// We'll get an empty string for a destructuring pattern.
		// Also if originalParamList came from a FunctionType instead of an
		// actual FUNCTION node, then all of the NAME nodes in it will have
		// empty strings, so we'll end up generating names for all of them.
		originalParamNames = append(originalParamNames, e.GetOriginalNameForParam(originalParam))
	}
	return e.CreateExternsParamListFromOriginalParamList(originalParamNames)
}

// CreateExternsParamListFromFunctionType creates a PARAM_LIST to store in
// the AST we'll use to generate externs for a function with the given type.
//
// If the NODE defining the original function is available, it would be
// better to use createExternsParamListFromOriginalFunction(), because that
// one will keep the parameter names the same instead of generating arbitrary
// parameter names.
//
// functionType is the JSType read from the FUNCTION (or possibly CLASS)
// node.
//
// Java: ExternExportsPass.Export.createExternsParamListFromFunctionType(JSType)
func (e *ExternExportsPass_ExportBase) CreateExternsParamListFromFunctionType(functionType rhino.JSType) *rhino.Node {
	// Place empty stings in the positions where we'll need to generate names.
	emptyParamNames := make([]string, len(functionType.AssertFunctionType().GetParameters()))

	return e.CreateExternsParamListFromOriginalParamList(emptyParamNames)
}

// CreateExternsParamListFromOriginalParamList creates a PARAM_LIST to store
// in the AST we'll use to generate externs for a function.
// originalParamNames are names for the parameters, possibly synthetic.
//
// Java: ExternExportsPass.Export.createExternsParamListFromOriginalParamList(List<String>)
func (e *ExternExportsPass_ExportBase) CreateExternsParamListFromOriginalParamList(originalParamNames []string) *rhino.Node {
	paramList := rhino.IR_ParamList()
	var nameGenerator NameGenerator = NewDefaultNameGeneratorWithReserved(
		jutil.NewLinkedHashSet(originalParamNames...),
		"",
		/* reservedCharacters= */ jutil.NewLinkedHashSet[uint16]())
	for _, originalParamName := range originalParamNames {
		externParamName := originalParamName
		if originalParamName == "" {
			externParamName = nameGenerator.GenerateNextName()
		}
		paramList.AddChildToBack(rhino.IR_Name(externParamName))
	}
	return paramList
}

// GetOriginalNameForParam: paramNode is expected to be a node in a
// PARAM_LIST. Returns the original name of the parameter, if possible,
// otherwise an empty string.
//
// Java: ExternExportsPass.Export.getOriginalNameForParam(Node)
func (e *ExternExportsPass_ExportBase) GetOriginalNameForParam(paramNode *rhino.Node) string {
	var nameOrPatternNode *rhino.Node
	if paramNode.IsRest() {
		// get name or pattern from `...nameOrPattern`
		nameOrPatternNode = paramNode.GetOnlyChild()
	} else if paramNode.IsDefaultValue() {
		// get name or pattern from `nameOrPattern = defaultValue`
		nameOrPatternNode = paramNode.GetFirstChild()
	} else {
		nameOrPatternNode = paramNode
	}
	if nameOrPatternNode.IsName() {
		originalName := nameOrPatternNode.GetOriginalName()
		if originalName != nil {
			return *originalName
		}
		return nameOrPatternNode.GetString()
	} else {
		jutil.CheckStateMsg1(nameOrPatternNode.IsDestructuringPattern(), "%s", nameOrPatternNode)
		return ""
	}
}

// CreateExternFunctionForEs6Class: given a class to export, create the empty
// function that will be put in the externs file.
//
// This extern function should have the same type as the original function
// and the same parameter name but no function body.
//
// TODO(b/123352214): It would be nice if we could put ES6 classes in the
// generated externs, but we'd have to fix some things first.
//
// Java: ExternExportsPass.Export.createExternFunctionForEs6Class(Node)
func (e *ExternExportsPass_ExportBase) CreateExternFunctionForEs6Class(exportedClass *rhino.Node) *rhino.Node {
	constructorMethodDefinition :=
		NodeUtil_GetEs6ClassConstructorMemberFunctionDef(exportedClass)
	if constructorMethodDefinition == nil {
		// no constructor for the class, so just create an empty function with
		// parameters to match the parameters indicated in the JSType, which
		// should have inherited parameters from the superclass, if any.
		classJSType := exportedClass.GetJSType()
		paramList := e.CreateExternsParamListFromFunctionType(classJSType)
		externFunction := rhino.IR_Function(rhino.IR_Name(""), paramList, rhino.IR_Block())
		externFunction.SetJSType(classJSType)
		return externFunction
	} else {
		// The JSType on the constructor function definition is the same as the
		// JSType on the whole class, so we can just pretend that the function is
		// an ES5 constructor function.
		return e.CreateExternFunction(constructorMethodDefinition.GetOnlyChild())
	}
}

// Java: ExternExportsPass.Export.buildEmptyJSDoc()
func (e *ExternExportsPass_ExportBase) BuildEmptyJSDoc() *rhino.JSDocInfo {
	// TODO(johnlenz): share the JSDocInfo here rather than building
	// a new one each time.
	return rhino.JSDocInfo_Builder_().BuildWithAlways(true)
}

// Java: ExternExportsPass.Export.buildNamespaceJSDoc()
func (e *ExternExportsPass_ExportBase) BuildNamespaceJSDoc() *rhino.JSDocInfo {
	builder := rhino.JSDocInfo_Builder_()
	builder.RecordConstancy()
	builder.RecordSuppressions(jutil.NewLinkedHashSet("const", "duplicate"))
	return builder.Build()
}

// CreateExternObjectLit: given an object literal to export, create an
// object lit with all its string properties. We don't care what the values
// of those properties are because they are not checked.
//
// Java: ExternExportsPass.Export.createExternObjectLit(Node)
func (e *ExternExportsPass_ExportBase) CreateExternObjectLit(exportedObjectLit *rhino.Node) *rhino.Node {
	lit := rhino.IR_Objectlit()
	lit.SetJSType(exportedObjectLit.GetJSType())

	// This is an indirect way of telling the typed code generator
	// "print the type of this"
	lit.SetJSDocInfo(e.BuildEmptyJSDoc())

	index := 1
	for child := exportedObjectLit.GetFirstChild(); child != nil; child = child.GetNext() {
		// TODO(dimvar): handle getters or setters?
		if child.IsStringKey() {
			lit.AddChildToBack(
				rhino.IR_Propdef(
					rhino.IR_StringKey(child.GetString()),
					rhino.IR_Number(float64(index))))
			index++
		}
	}
	return lit
}

// GetValue: if the given value is a qualified name which refers a function
// or object literal, the node is returned. Otherwise, null is returned.
//
// Java: ExternExportsPass.Export.getValue()
func (e *ExternExportsPass_ExportBase) GetValue() *rhino.Node {
	qualifiedName := e.value.GetQualifiedName()

	if qualifiedName == nil {
		// We expect to see
		// goog.exportSymbol('exportedName', some.path);
		// goog.exportProperty(some.path, 'exportedName', some.path.prop);
		//
		// In either case `value` will be the last argument, which we expect to
		// be a qualified name If it isn't we won't include any type information
		// in the output externs. It would be very strange to use a literal value
		// as the final argument, since it wouldn't then be accessible by any
		// non-exported name.
		return nil
	}

	definition := e.outer.definitionMap.Get(*qualifiedName)
	if definition == nil {
		// Couldn't find any assignment to the qualified name
		return nil
	}

	if definition.IsFunction() || definition.IsClass() || definition.IsObjectLit() {
		// We can generate good type information for all of these cases.
		return definition
	}

	// value was something unusual, so we won't return any node from which to
	// get type information.
	return nil
}

// ExternExportsPass_SymbolExport is a symbol export.
type ExternExportsPass_SymbolExport struct {
	ExternExportsPass_ExportBase
}

// Java: ExternExportsPass.SymbolExport(String, Node)
func (p *ExternExportsPass) NewSymbolExport(symbolName string, value *rhino.Node) *ExternExportsPass_SymbolExport {
	e := &ExternExportsPass_SymbolExport{}
	e.InitExternExportsPass_ExportBase(e, p, symbolName, value)

	qualifiedName := value.GetQualifiedName()

	if qualifiedName != nil {
		p.mappedPaths.Put(*qualifiedName, symbolName)
	}
	return e
}

// Java: ExternExportsPass.SymbolExport.getExportedPath()
func (e *ExternExportsPass_SymbolExport) GetExportedPath() string {
	return e.symbolName
}

// ExternExportsPass_PropertyExport is a property export.
type ExternExportsPass_PropertyExport struct {
	ExternExportsPass_ExportBase
	exportPath string
}

// Java: ExternExportsPass.PropertyExport(String, String, Node)
func (p *ExternExportsPass) NewPropertyExport(exportPath string, symbolName string, value *rhino.Node) *ExternExportsPass_PropertyExport {
	e := &ExternExportsPass_PropertyExport{}
	e.InitExternExportsPass_ExportBase(e, p, symbolName, value)

	e.exportPath = exportPath
	return e
}

// Java: ExternExportsPass.PropertyExport.getExportedPath()
func (e *ExternExportsPass_PropertyExport) GetExportedPath() string {
	// Find the longest path that has been mapped (if any).
	prefixes := externExportsPass_computePathPrefixes(e.exportPath)
	for i := len(prefixes) - 1; i >= 0; i-- {
		currentPath := prefixes[i]
		jutil.CheckState(len(currentPath) > 0)

		// If this path is mapped, return the mapped path plus any remaining
		// pieces.
		mappedPath, ok := e.outer.mappedPaths.Lookup(currentPath)
		if !ok {
			continue
		}

		// Append the remaining path segments, including a leading separator.
		mappedPath += e.exportPath[len(currentPath):]
		return externExportsPass_qNameJoin(mappedPath, e.symbolName)
	}

	return externExportsPass_qNameJoin(e.exportPath, e.symbolName)
}

// externExportsPass_computePathPrefixes computes a list of the path prefixes
// constructed from the components of the path.
//
// E.g., if the path is:
//
//	"a.b.c"
//
// then then path prefixes will be
//
//	["a","a.b","a.b.c"]:
//
// Java: ExternExportsPass.computePathPrefixes(String)
func externExportsPass_computePathPrefixes(path string) []string {
	pieces := strings.Split(path, ".")
	var pathPrefixes []string

	partial := pieces[0] // There will always be at least 1.
	pathPrefixes = append(pathPrefixes, partial)
	for i := 1; i < len(pieces); i++ {
		partial = externExportsPass_qNameJoin(partial, pieces[i])
		pathPrefixes = append(pathPrefixes, partial)
	}

	return pathPrefixes
}

// NewExternExportsPass creates an instance.
//
// Java: ExternExportsPass(AbstractCompiler)
func NewExternExportsPass(compiler AbstractCompiler) *ExternExportsPass {
	p := &ExternExportsPass{}
	p.exports = jutil.NewArrayList[ExternExportsPass_Export]()
	p.compiler = compiler
	p.definitionMap = jutil.NewLinkedHashMap[string, *rhino.Node]()
	p.externsRoot = rhino.IR_Script()
	p.alreadyExportedPaths = jutil.NewLinkedHashSet[string]()
	p.mappedPaths = jutil.NewLinkedHashMap[string, string]()

	p.InitExportMethods()
	return p
}

// Java: ExternExportsPass.initExportMethods()
func (p *ExternExportsPass) InitExportMethods() {
	convention := p.compiler.GetCodingConvention()
	// Go-only: language-adapter — ImmutableSet.of rejects a null element
	// (NullPointerException), as CheckNotNull does.
	p.exportSymbolFunctionNames =
		jutil.NewLinkedHashSet(
			*jutil.CheckNotNull(convention.GetExportSymbolFunction()), // goog.exportSymbol(name, value)
			"google_exportSymbol") // used within Google

	p.exportPropertyFunctionNames =
		jutil.NewLinkedHashSet(
			*jutil.CheckNotNull(convention.GetExportPropertyFunction()), // goog.exportProperty(owner, name, value)
			"google_exportProperty") // used within Google
}

// Java: ExternExportsPass.process(Node, Node)
func (p *ExternExportsPass) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(p.compiler, root, p)

	// Sort by path length to ensure that the longer
	// paths (which may depend on the shorter ones)
	// come later.
	sorted := jutil.NewTreeSet[ExternExportsPass_Export](func(a, b ExternExportsPass_Export) int {
		return jutil.StringCompareTo(a.GetExportedPath(), b.GetExportedPath())
	})

	for e := range p.exports.All() {
		sorted.Add(e)
	}

	for export := range sorted.All() {
		export.GenerateExterns()
	}

	p.SetGeneratedExternsOnCompiler()
}

// Java: ExternExportsPass.setGeneratedExternsOnCompiler()
func (p *ExternExportsPass) SetGeneratedExternsOnCompiler() {
	builder := NewCodePrinter_Builder(p.externsRoot).
		SetPrettyPrint(true).
		SetOutputTypes(true).
		SetTypeRegistry(p.compiler.GetTypeRegistry())

	p.compiler.SetExternExports(
		"/**\n" +
			" * @fileoverview Generated externs.\n" +
			" * @externs\n" +
			" */\n" +
			builder.Build())
}

// Java: ExternExportsPass.visit(NodeTraversal, Node, Node)
func (p *ExternExportsPass) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	p.LookForQnameDefinition(n)
	p.LookForAtExportOnThisDotProperty(t, n)
	p.LookForSymbolExportCall(n)
	p.LookForPropertyExportCall(n)
}

// Java: ExternExportsPass.lookForQnameDefinition(Node)
func (p *ExternExportsPass) LookForQnameDefinition(n *rhino.Node) {
	// TODO(b/123725559): There are lots of cases where this could fail to find
	//     the right definition or be fooled by there being multiple
	//     definitions.
	if n.IsClass() {
		if NodeUtil_IsClassDeclaration(n) {
			// class Foo {...}
			p.definitionMap.Put(n.GetFirstChild().GetString(), n)
		}
	} else if n.IsFunction() {
		if NodeUtil_IsFunctionDeclaration(n) {
			// function foo() {...}
			p.definitionMap.Put(n.GetFirstChild().GetString(), n)
		}
	} else if n.IsAssign() {
		// TODO(b/123718645): Add support for destructuring assignments
		lhs := n.GetFirstChild()
		if lhs.IsQualifiedName() && (!lhs.HasChildren() || !lhs.GetFirstChild().IsThis()) {
			// qualified.name = value;
			//   but not
			// this.prop = value;
			p.definitionMap.Put(*lhs.GetQualifiedName(), n.GetLastChild())
		}
	} else if n.IsName() {
		// TODO(b/123718645): Add support for destructuring declarations
		parent := jutil.CheckNotNullMsg1(n.GetParent(), "%s", n)
		if NodeUtil_IsNameDeclaration(parent) {
			value := n.GetFirstChild()
			if value != nil {
				// const foo = value;
				p.definitionMap.Put(n.GetString(), value)
			}
		}
	} else if n.IsMemberFunctionDef() {
		// Try to find a fully qualified name for the method
		lvalueName := NodeUtil_GetBestLValueName(n)
		if lvalueName != nil {
			// Store the function as the value
			p.definitionMap.Put(*lvalueName, n.GetOnlyChild())
		}
	}
	// TODO(b/123725422): Getters and setters?
}

// Java: ExternExportsPass.lookForSymbolExportCall(Node)
func (p *ExternExportsPass) LookForSymbolExportCall(n *rhino.Node) {
	if !p.IsCallToOneOf(n, p.exportSymbolFunctionNames) {
		return // not a call to goog.exportSymbol()
	}
	// TODO(b/123725716): We should report errors for malformed calls instead
	// of just ignoring them. Ensure that we only check valid calls with the 2
	// arguments (plus the GETPROP node itself).
	if !n.HasXChildren(3) {
		return
	}

	thisNode := n.GetFirstChild()
	nameArg := thisNode.GetNext()
	valueArg := nameArg.GetNext()

	// Confirm the arguments are the expected types. If they are not,
	// then we have an export that we cannot statically identify.
	if !nameArg.IsStringLit() {
		return
	}

	// Add the export to the list.
	p.exports.Add(p.NewSymbolExport(nameArg.GetString(), valueArg))
}

// Java: ExternExportsPass.lookForPropertyExportCall(Node)
func (p *ExternExportsPass) LookForPropertyExportCall(n *rhino.Node) {
	if !p.IsCallToOneOf(n, p.exportPropertyFunctionNames) {
		return // not a call to goog.exportProperty()
	}
	// TODO(b/123725716): We should report errors for malformed calls instead
	// of just ignoring them. Ensure that we only check valid calls with the 3
	// arguments (plus the GETPROP node itself).
	if !n.HasXChildren(4) {
		return
	}

	thisNode := n.GetFirstChild()
	objectArg := thisNode.GetNext()
	nameArg := objectArg.GetNext()
	valueArg := nameArg.GetNext()

	// Confirm the arguments are the expected types. If they are not,
	// then we have an export that we cannot statically identify.
	if !objectArg.IsQualifiedName() {
		return
	}

	if !nameArg.IsStringLit() {
		return
	}

	// Add the export to the list.
	p.exports.Add(
		p.NewPropertyExport(*objectArg.GetQualifiedName(),
			nameArg.GetString(),
			valueArg))
}

// Java: ExternExportsPass.isCallToOneOf(Node, ImmutableSet<String>)
func (p *ExternExportsPass) IsCallToOneOf(n *rhino.Node, functionQnames *jutil.LinkedHashSet[string]) bool {
	if !n.IsCall() {
		return false
	} else {
		callee := n.GetFirstChild()
		return callee.IsQualifiedName() && functionQnames.Contains(*callee.GetQualifiedName())
	}
}

// Java: ExternExportsPass.lookForAtExportOnThisDotProperty(NodeTraversal, Node)
func (p *ExternExportsPass) LookForAtExportOnThisDotProperty(t *NodeTraversal, thisDotPropName *rhino.Node) {
	if !thisDotPropName.IsGetProp() || !thisDotPropName.GetFirstChild().IsThis() {
		return // not this.propName
	}
	jsdoc := NodeUtil_GetBestJSDocInfo(thisDotPropName)
	if jsdoc == nil || !jsdoc.IsExport() {
		return // no @export on this.propName
	}

	constructorNode := t.GetEnclosingFunction()
	if !NodeUtil_IsConstructor(constructorNode) {
		return // @export on this.propName only works within a constructor
	}

	var classNode *rhino.Node
	if NodeUtil_IsEs6Constructor(constructorNode) {
		classNode = NodeUtil_GetEnclosingClass(constructorNode)
	} else {
		classNode = constructorNode
	}
	className := NodeUtil_GetName(classNode)
	propertyName := thisDotPropName.GetString()
	prototypeName := jutil.ToJavaString(className) + ".prototype"
	propertyNameNode := NodeUtil_NewQName(p.compiler, "this."+propertyName)

	// Add the export to the list.
	p.exports.Add(p.NewPropertyExport(prototypeName, propertyName, propertyNameNode))
}
