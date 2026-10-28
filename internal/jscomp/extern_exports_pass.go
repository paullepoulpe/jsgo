// Port of com/google/javascript/jscomp/ExternExportsPass.java @ closure-compiler v20261021

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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ExternExportsPass creates an externs file containing all exported symbols
// and properties for later consumption.
type ExternExportsPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	// The exports found.
	exports []externExport

	// A map of all assigns to their parent nodes.
	definitionMap *jutil.LinkedHashMap[string, rhino.Node]

	// The parent compiler.
	compiler AbstractCompiler

	// The AST root which holds the externs generated.
	externsRoot rhino.Node

	// A mapping of internal paths to exported paths.
	mappedPaths *jutil.LinkedHashMap[string, string]

	// A set of exported paths.
	alreadyExportedPaths map[string]struct{}

	// Function name used to export symbols (in addition to google_exportSymbol).
	exportSymbolFunction string

	// Function name used to export properties (in addition to google_exportProperty).
	exportPropertyFunction string
}

var _ CompilerPass = (*ExternExportsPass)(nil)
var _ NodeTraversal_Callback = (*ExternExportsPass)(nil)

type exportKind uint8

const (
	exportKindSymbol exportKind = iota
	exportKindProperty
)

// externExport collapses ExternExportsPass.Export, SymbolExport, and PropertyExport.
type externExport struct {
	kind       exportKind
	exportPath string
	symbolName string
	value      rhino.Node
}

// Q_NAME_JOINER is Joiner.on('.').
func externExportsPass_qNameJoin(a string, b string) string {
	return a + "." + b
}

func (p *ExternExportsPass) addSymbolExport(symbolName string, value rhino.Node) {
	rhino.CheckNotNull(value)
	if qualifiedName := value.GetQualifiedName(); qualifiedName != nil {
		p.mappedPaths.Put(*qualifiedName, symbolName)
	}
	p.exports = append(p.exports, externExport{
		kind:       exportKindSymbol,
		symbolName: symbolName,
		value:      value,
	})
}

func (p *ExternExportsPass) addPropertyExport(exportPath string, symbolName string, value rhino.Node) {
	rhino.CheckNotNull(value)
	p.exports = append(p.exports, externExport{
		kind:       exportKindProperty,
		exportPath: exportPath,
		symbolName: symbolName,
		value:      value,
	})
}

func (p *ExternExportsPass) exportedPath(e externExport) string {
	if e.kind == exportKindSymbol {
		return e.symbolName
	}
	// Find the longest path that has been mapped (if any).
	prefixes := externExportsPass_computePathPrefixes(e.exportPath)
	for i := len(prefixes) - 1; i >= 0; i-- {
		currentPath := prefixes[i]
		jutil.CheckState(len(currentPath) > 0)

		mappedPath, ok := p.mappedPaths.Lookup(currentPath)
		if !ok {
			continue
		}

		// Append the remaining path segments, including a leading separator.
		mappedPath += e.exportPath[len(currentPath):]
		return externExportsPass_qNameJoin(mappedPath, e.symbolName)
	}
	return externExportsPass_qNameJoin(e.exportPath, e.symbolName)
}

func (p *ExternExportsPass) exportValue(e externExport) rhino.Node {
	qualifiedName := e.value.GetQualifiedName()
	if qualifiedName == nil {
		return rhino.NilNode
	}
	definition := p.definitionMap.Get(*qualifiedName)
	if definition.IsNil() {
		return rhino.NilNode
	}
	if definition.IsFunction() || definition.IsClass() || definition.IsObjectLit() {
		return definition
	}
	return rhino.NilNode
}

// appendExtern appends the exported function and all paths necessary for
// the path to be declared.
func (p *ExternExportsPass) appendExtern(path string, valueToExport rhino.Node) {
	pathPrefixes := externExportsPass_computePathPrefixes(path)

	for i, pathPrefix := range pathPrefixes {
		// The complete path (the last path prefix) must be emitted and
		// it gets initialized to the externed version of the value.
		isCompletePathPrefix := i == len(pathPrefixes)-1

		_, alreadyExported := p.alreadyExportedPaths[pathPrefix]
		if strings.HasSuffix(pathPrefix, ".prototype") || (alreadyExported && !isCompletePathPrefix) {
			continue
		}

		exportedValueDefinesNewType := false
		if !valueToExport.IsNil() {
			jsdoc := NodeUtil_GetBestJSDocInfo(valueToExport)
			if valueToExport.IsClass() || (jsdoc != nil && jsdoc.ContainsTypeDefinition()) {
				exportedValueDefinesNewType = true
			}
		}

		var initializer rhino.Node
		var jsdoc *rhino.JSDocInfo
		if isCompletePathPrefix && !valueToExport.IsNil() {
			if valueToExport.IsFunction() {
				initializer = p.createExternFunction(valueToExport)
			} else if valueToExport.IsClass() {
				initializer = p.createExternFunctionForEs6Class(valueToExport)
			} else {
				jutil.CheckState(valueToExport.IsObjectLit())
				initializer = p.createExternObjectLit(valueToExport)
			}
		} else if !isCompletePathPrefix && exportedValueDefinesNewType {
			jsdoc = buildNamespaceJSDoc()
			initializer = p.createExternObjectLit(rhino.IR_Objectlit())
			// Don't add the empty jsdoc here
			initializer.SetJSDocInfo(nil)
		} else {
			initializer = rhino.IR_Empty()
		}

		p.appendPathDefinition(pathPrefix, initializer, jsdoc)
	}
}

func (p *ExternExportsPass) appendPathDefinition(path string, initializer rhino.Node, jsdoc *rhino.JSDocInfo) {
	var pathDefinition rhino.Node

	if strings.Contains(path, ".") {
		qualifiedPath := NodeUtil_NewQName(p.compiler, path)
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
		if exprStmt, ok := ast.AsExprStmt(pathDefinition); ok {
			exprStmt.Expr().Node().SetJSDocInfo(jsdoc)
		} else {
			jutil.CheckState(pathDefinition.IsVar())
			pathDefinition.SetJSDocInfo(jsdoc)
		}
	}

	p.externsRoot.AddChildToBack(pathDefinition)
	p.alreadyExportedPaths[path] = struct{}{}
}

func (p *ExternExportsPass) createExternFunction(exportedFunction rhino.Node) rhino.Node {
	paramList := createExternsParamListFromOriginalFunction(exportedFunction)
	externFunction := rhino.IR_Function(rhino.IR_Name(""), paramList, rhino.IR_Block())
	externFunction.SetJSType(exportedFunction.GetJSType())
	return externFunction
}

func createExternsParamListFromOriginalFunction(exportedFunction rhino.Node) rhino.Node {
	params := ast.MustFunction(exportedFunction).Params()
	originalParamNames := make([]string, 0, params.NumParams())
	for originalParam := range params.Params() {
		originalParamNames = append(originalParamNames, getOriginalNameForParam(originalParam.Node()))
	}
	return createExternsParamListFromOriginalParamList(originalParamNames)
}

func createExternsParamListFromFunctionType(functionType rhino.JSType) rhino.Node {
	emptyParamNames := make([]string, len(functionType.AssertFunctionType().GetParameters()))
	return createExternsParamListFromOriginalParamList(emptyParamNames)
}

func createExternsParamListFromOriginalParamList(originalParamNames []string) rhino.Node {
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

func getOriginalNameForParam(paramNode rhino.Node) string {
	nameOrPatternNode := paramNode
	if rest, ok := ast.AsIterRest(paramNode); ok {
		nameOrPatternNode = rest.Target().Node()
	} else if defVal, ok := ast.AsDefaultValue(paramNode); ok {
		nameOrPatternNode = defVal.Target().Node()
	}
	if name, ok := ast.AsName(nameOrPatternNode); ok {
		if originalName := nameOrPatternNode.GetOriginalName(); originalName != nil {
			return *originalName
		}
		return name.Ident()
	}
	jutil.CheckStateMsg1(nameOrPatternNode.IsDestructuringPattern(), "%s", nameOrPatternNode)
	return ""
}

func (p *ExternExportsPass) createExternFunctionForEs6Class(exportedClass rhino.Node) rhino.Node {
	if ctorDef, ok := ast.MustClass(exportedClass).Constructor(); ok {
		return p.createExternFunction(ctorDef.Function().Node())
	}
	classJSType := exportedClass.GetJSType()
	paramList := createExternsParamListFromFunctionType(classJSType)
	externFunction := rhino.IR_Function(rhino.IR_Name(""), paramList, rhino.IR_Block())
	externFunction.SetJSType(classJSType)
	return externFunction
}

func buildEmptyJSDoc() *rhino.JSDocInfo {
	return rhino.JSDocInfo_Builder_().BuildWithAlways(true)
}

func buildNamespaceJSDoc() *rhino.JSDocInfo {
	builder := rhino.JSDocInfo_Builder_()
	builder.RecordConstancy()
	builder.RecordSuppressions(jutil.NewLinkedHashSet("const", "duplicate"))
	return builder.Build()
}

func (p *ExternExportsPass) createExternObjectLit(exportedObjectLit rhino.Node) rhino.Node {
	lit := rhino.IR_Objectlit()
	lit.SetJSType(exportedObjectLit.GetJSType())
	lit.SetJSDocInfo(buildEmptyJSDoc())

	index := 1
	for member := range ast.MustObjectLit(exportedObjectLit).Members() {
		if sk, ok := member.AsStringKey(); ok {
			lit.AddChildToBack(
				rhino.IR_Propdef(
					rhino.IR_StringKey(sk.Key()),
					rhino.IR_Number(float64(index))))
			index++
		}
	}
	return lit
}

// externExportsPass_computePathPrefixes computes a list of the path prefixes
// constructed from the components of the path (e.g. "a.b.c" -> ["a","a.b","a.b.c"]).
func externExportsPass_computePathPrefixes(path string) []string {
	pieces := strings.Split(path, ".")
	pathPrefixes := make([]string, 0, len(pieces))

	partial := pieces[0]
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
	p := &ExternExportsPass{
		compiler:             compiler,
		definitionMap:        jutil.NewLinkedHashMap[string, rhino.Node](),
		externsRoot:          rhino.IR_Script(),
		alreadyExportedPaths: make(map[string]struct{}),
		mappedPaths:          jutil.NewLinkedHashMap[string, string](),
	}

	p.InitExportMethods()
	return p
}

// Java: ExternExportsPass.initExportMethods()
func (p *ExternExportsPass) InitExportMethods() {
	convention := p.compiler.GetCodingConvention()
	p.exportSymbolFunction = *jutil.CheckNotNull(convention.GetExportSymbolFunction())     // goog.exportSymbol(name, value)
	p.exportPropertyFunction = *jutil.CheckNotNull(convention.GetExportPropertyFunction()) // goog.exportProperty(owner, name, value)
}

// Java: ExternExportsPass.process(Node, Node)
func (p *ExternExportsPass) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(p.compiler, root, p)

	type resolvedExport struct {
		path   string
		export externExport
	}
	resolved := make([]resolvedExport, len(p.exports))
	for i, e := range p.exports {
		resolved[i] = resolvedExport{path: p.exportedPath(e), export: e}
	}
	slices.SortStableFunc(resolved, func(a, b resolvedExport) int {
		return jutil.StringCompareTo(a.path, b.path)
	})

	for i, item := range resolved {
		if i > 0 && item.path == resolved[i-1].path {
			continue
		}
		p.appendExtern(item.path, p.exportValue(item.export))
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
func (p *ExternExportsPass) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	p.LookForQnameDefinition(n)
	p.LookForAtExportOnThisDotProperty(t, n)
	p.LookForSymbolExportCall(n)
	p.LookForPropertyExportCall(n)
}

// Java: ExternExportsPass.lookForQnameDefinition(Node)
func (p *ExternExportsPass) LookForQnameDefinition(n rhino.Node) {
	if cls, ok := ast.AsClass(n); ok {
		if NodeUtil_IsClassDeclaration(n) {
			p.definitionMap.Put(cls.NameNode().GetString(), n)
		}
	} else if fn, ok := ast.AsFunction(n); ok {
		if NodeUtil_IsFunctionDeclaration(n) {
			p.definitionMap.Put(fn.NameNode().Ident(), n)
		}
	} else if assign, ok := ast.AsAssign(n); ok {
		lhs := assign.Target().Node()
		if lhs.IsQualifiedName() && (!lhs.HasChildren() || !lhs.GetFirstChild().IsThis()) {
			p.definitionMap.Put(*lhs.GetQualifiedName(), assign.Value().Node())
		}
	} else if name, ok := ast.AsName(n); ok {
		parent := rhino.CheckNotNullMsg1(n.GetParent(), "%s", n)
		if NodeUtil_IsNameDeclaration(parent) {
			if value, ok := name.Init(); ok {
				p.definitionMap.Put(name.Ident(), value.Node())
			}
		}
	} else if methodDef, ok := ast.AsMemberFunctionDef(n); ok {
		if lvalueName := NodeUtil_GetBestLValueName(n); lvalueName != nil {
			p.definitionMap.Put(*lvalueName, methodDef.Function().Node())
		}
	}
}

// Java: ExternExportsPass.lookForSymbolExportCall(Node)
func (p *ExternExportsPass) LookForSymbolExportCall(n rhino.Node) {
	call, ok := p.asCallToOneOf(n, p.exportSymbolFunction, "google_exportSymbol")
	if !ok || call.NumArgs() != 2 {
		return
	}
	callOrNew := call.AsCallOrNew()
	nameArg, _ := callOrNew.ArgAt(0)
	valueArg, _ := callOrNew.ArgAt(1)
	nameLit, ok := nameArg.AsStringLit()
	if !ok {
		return
	}
	p.addSymbolExport(nameLit.Value(), valueArg.Node())
}

// Java: ExternExportsPass.lookForPropertyExportCall(Node)
func (p *ExternExportsPass) LookForPropertyExportCall(n rhino.Node) {
	call, ok := p.asCallToOneOf(n, p.exportPropertyFunction, "google_exportProperty")
	if !ok || call.NumArgs() != 3 {
		return
	}
	callOrNew := call.AsCallOrNew()
	objectArg, _ := callOrNew.ArgAt(0)
	nameArg, _ := callOrNew.ArgAt(1)
	valueArg, _ := callOrNew.ArgAt(2)
	if !objectArg.Node().IsQualifiedName() {
		return
	}
	nameLit, ok := nameArg.AsStringLit()
	if !ok {
		return
	}
	p.addPropertyExport(*objectArg.Node().GetQualifiedName(), nameLit.Value(), valueArg.Node())
}

func (p *ExternExportsPass) asCallToOneOf(n rhino.Node, fn1, fn2 string) (ast.Call, bool) {
	call, ok := ast.AsCall(n)
	if !ok {
		return ast.Call{}, false
	}
	callee := call.Callee().Node()
	if !callee.IsQualifiedName() {
		return ast.Call{}, false
	}
	qname := *callee.GetQualifiedName()
	if qname != fn1 && qname != fn2 {
		return ast.Call{}, false
	}
	return call, true
}

// Java: ExternExportsPass.lookForAtExportOnThisDotProperty(NodeTraversal, Node)
func (p *ExternExportsPass) LookForAtExportOnThisDotProperty(t *NodeTraversal, thisDotPropName rhino.Node) {
	gp, ok := ast.AsGetProp(thisDotPropName)
	if !ok || !gp.Object().Node().IsThis() {
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

	var classNode rhino.Node
	if NodeUtil_IsEs6Constructor(constructorNode) {
		classNode = NodeUtil_GetEnclosingClass(constructorNode)
	} else {
		classNode = constructorNode
	}
	className := NodeUtil_GetName(classNode)
	propertyName := gp.Property()
	prototypeName := jutil.ToJavaString(className) + ".prototype"
	propertyNameNode := NodeUtil_NewQName(p.compiler, "this."+propertyName)

	p.addPropertyExport(prototypeName, propertyName, propertyNameNode)
}
