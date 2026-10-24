// Port of com/google/javascript/jscomp/GenerateExports.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// GenerateExports generates goog.exportSymbol/goog.exportProperty for the
// @export annotation.
type GenerateExports struct {
	compiler AbstractCompiler

	exportSymbolFunction *string

	exportPropertyFunction *string

	allowNonGlobalExports bool

	exportedVariables *jutil.LinkedHashSet[string]
}

var _ CompilerPass = (*GenerateExports)(nil)

const generateExports_PROTOTYPE_PROPERTY = "prototype"

// Java: GenerateExports.MISSING_EXPORT_CONVENTION
var GenerateExports_MISSING_EXPORT_CONVENTION = DiagnosticType_Error(
	"JSC_MISSING_EXPORT_CONVENTION",
	"@export cannot be used without defining a exportProperty and exportSymbol function in"+
		" the coding convention")

// Java: GenerateExports.MISSING_GOOG_FOR_EXPORT
var GenerateExports_MISSING_GOOG_FOR_EXPORT = DiagnosticType_Error(
	"JSC_MISSING_EXPORT_SYMBOL_DEFINITION",
	"@export cannot be used without including closure/base.js or other definition of"+
		" {0} and {1}")

// NewGenerateExports creates a new generate exports compiler pass.
//
// The exportSymbolFunction and exportPropertyFunction are optional, but the
// compiler will report an error if they are not passed *and* this pass
// finds @export annotations. They are allowed to be optional if no code uses
// the @export annotation.
//
// compiler is the JS compiler; exportSymbolFunction the function used for
// exporting symbols; exportPropertyFunction the function used for exporting
// property names.
//
// Java: GenerateExports(AbstractCompiler, boolean, String, String)
func NewGenerateExports(
	compiler AbstractCompiler,
	allowNonGlobalExports bool,
	exportSymbolFunction *string,
	exportPropertyFunction *string) *GenerateExports {
	jutil.CheckNotNull(compiler)

	g := &GenerateExports{exportedVariables: jutil.NewLinkedHashSet[string]()}
	g.compiler = compiler
	g.allowNonGlobalExports = allowNonGlobalExports
	g.exportSymbolFunction = exportSymbolFunction
	g.exportPropertyFunction = exportPropertyFunction
	return g
}

// Java: GenerateExports.getExportedVariableNames()
func (g *GenerateExports) GetExportedVariableNames() *jutil.LinkedHashSet[string] {
	return g.exportedVariables
}

// Java: GenerateExports.process(Node, Node)
func (g *GenerateExports) Process(externs *rhino.Node, root *rhino.Node) {
	findExportableNodes := NewFindExportableNodes(g.compiler, g.allowNonGlobalExports)
	NodeTraversal_Traverse(g.compiler, root, findExportableNodes)
	exports := findExportableNodes.GetExports()
	es6Exports := findExportableNodes.GetEs6ClassExports()
	localExports := findExportableNodes.GetLocalExports()

	for export := range localExports.All() {
		g.AddExtern(export)
	}

	if !exports.IsEmpty() || !es6Exports.IsEmpty() {
		if !g.ValidateExportMethodsIncluded(exports, es6Exports, root) {
			return
		}
	}

	for k, v := range es6Exports.All() {
		g.AddExportForEs6Method(k, v)
	}

	for export, context := range exports.All() {
		g.AddExportMethod(exports, export, context)
	}
}

// ValidateExportMethodsIncluded validates that a) the user has configured
// methods to call to export symbols and b) those methods are included in
// this binary. Returns whether to continue trying to export methods.
//
// Java: GenerateExports.validateExportMethodsIncluded(Map<String, Node>, Map<Node, String>, Node)
func (g *GenerateExports) ValidateExportMethodsIncluded(
	exports *jutil.LinkedHashMap[string, *rhino.Node], es6Exports *jutil.LinkedHashMap[*rhino.Node, string], root *rhino.Node) bool {

	// Pick an arbitrary @export to report the warning on.
	var errorLocation *rhino.Node
	if !exports.IsEmpty() {
		for v := range exports.Values() {
			errorLocation = v
			break
		}
	} else {
		for k := range es6Exports.Keys() {
			errorLocation = k
			break
		}
	}

	if g.exportSymbolFunction == nil || g.exportPropertyFunction == nil {
		g.compiler.Report(JSError_MakeWithNode(errorLocation, GenerateExports_MISSING_EXPORT_CONVENTION))
		// don't try to rewrite @export since there's nothing we can rewrite
		// them to.
		return false
	}

	if !g.IncludesExportMethods(root.GetParent()) {
		g.compiler.Report(
			JSError_MakeWithNode(
				errorLocation,
				GenerateExports_MISSING_GOOG_FOR_EXPORT,
				*g.exportSymbolFunction,
				*g.exportPropertyFunction))
		g.DeclareExportMethodsInExterns(errorLocation)
		return true
	}
	return true
}

// IncludesExportMethods: in order to use @export, the user must have
// included the defintiion of the property and symbol export methods in their
// binary.
//
// In real code these methods are typically goog.exportProperty and
// goog.exportSymbol.
//
// Java: GenerateExports.includesExportMethods(Node)
func (g *GenerateExports) IncludesExportMethods(root *rhino.Node) bool {
	globalScope := NewSyntacticScopeCreator(g.compiler).CreateScopeOfScope(root /* parent= */, nil)
	exportPropertyRoot := NodeUtil_GetRootOfQualifiedNameString(*g.exportPropertyFunction)
	exportSymbolRoot := NodeUtil_GetRootOfQualifiedNameString(*g.exportSymbolFunction)
	return globalScope.HasSlot(exportPropertyRoot) && globalScope.HasSlot(exportSymbolRoot)
}

// DeclareExportMethodsInExterns adds a synthesized externs declaration to
// prevent crashes later on in VarCheck.
//
// Java: GenerateExports.declareExportMethodsInExterns(Node)
func (g *GenerateExports) DeclareExportMethodsInExterns(srcref *rhino.Node) {

	exportPropertyRoot := NodeUtil_GetRootOfQualifiedNameString(*g.exportPropertyFunction)
	exportSymbolRoot := NodeUtil_GetRootOfQualifiedNameString(*g.exportSymbolFunction)

	if exportPropertyRoot == exportSymbolRoot {
		g.DeclareSyntheticExternsVar(exportPropertyRoot, srcref)
	} else {
		g.DeclareSyntheticExternsVar(exportPropertyRoot, srcref)
		g.DeclareSyntheticExternsVar(exportSymbolRoot, srcref)
	}
}

// Java: GenerateExports.declareSyntheticExternsVar(String, Node)
func (g *GenerateExports) DeclareSyntheticExternsVar(name string, srcref *rhino.Node) {
	syntheticInput := g.compiler.GetSynthesizedExternsInput()
	syntheticVarRoot := syntheticInput.GetAstRoot(g.compiler)

	varDeclaration := rhino.IR_Var(rhino.IR_Name(name)).SrcrefTree(srcref)
	varDeclaration.SetStaticSourceFile(g.compiler.GetSynthesizedExternsInput().GetSourceFile())
	syntheticVarRoot.AddChildToBack(varDeclaration)
	g.compiler.ReportChangeToEnclosingScope(varDeclaration)
}

// Java: GenerateExports.addExtern(String)
func (g *GenerateExports) AddExtern(export string) {
	objectPrototype := NodeUtil_NewQName(g.compiler, "Object.prototype")
	objCtor := g.compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_OBJECT_FUNCTION_TYPE)
	objectPrototype.GetFirstChild().SetJSType(objCtor)
	propstmt := rhino.IR_ExprResult(rhino.IR_Getprop(objectPrototype, export))
	propstmt.SrcrefTree(g.GetSynthesizedExternsRoot())
	propstmt.SetOriginalName(&export)
	g.GetSynthesizedExternsRoot().AddChildToBack(propstmt)
	g.compiler.ReportChangeToEnclosingScope(propstmt)
}

// Java: GenerateExports.recordExportSymbol(String)
func (g *GenerateExports) RecordExportSymbol(qname string) {
	dot := strings.IndexByte(qname, '.')
	if dot == -1 {
		g.exportedVariables.Add(qname)
	} else {
		g.exportedVariables.Add(qname[:dot])
	}
}

// Java: GenerateExports.addExportForEs6Method(Node, String)
func (g *GenerateExports) AddExportForEs6Method(memberFunction *rhino.Node, ownerName string) {
	// We always export ES6 member methods as properties.
	jutil.CheckArgumentMsg1(memberFunction.IsMemberFunctionDef(), "%s", memberFunction)
	jutil.CheckArgumentMsg1(ownerName != "", "%s", ownerName)
	fullExport := ownerName + "." + memberFunction.GetString()
	g.AddExportPropertyCall(&ownerName, memberFunction, fullExport, memberFunction.GetString())
}

// AddExportMethod emits a call to either goog.exportProperty or
// goog.exportSymbol.
//
// Attempts to optimize by creating a property export instead of a symbol
// export, because property exports are significantly simpler/faster.
//
// export is the fully qualified name of the object we want to export;
// context the node on which the @export annotation was found.
//
// Java: GenerateExports.addExportMethod(Map<String, Node>, String, Node)
func (g *GenerateExports) AddExportMethod(exports *jutil.LinkedHashMap[string, *rhino.Node], export string, context *rhino.Node) {
	// We can export as a property if any of the following conditions holds:
	// a) ES6 class members, which the above `addExportForEs6Method` handles
	// b) this is a property on a name which is also being exported
	// c) this is a prototype property
	var methodOwnerName *string            // the object this method is on, null for exported names.
	isEs5StylePrototypeAssignment := false // If this is a prototype property
	var propertyName string

	child := context.GetFirstChild()
	if child != nil && child.IsGetProp() { // e.g. `/** @export */ a.prototype.b = obj;`
		node := context.GetFirstChild()                // e.g. get `a.prototype.b`
		ownerNode := node.GetFirstChild()              // e.g. get `a.prototype`
		methodOwnerName = ownerNode.GetQualifiedName() // e.g. get the string "a.prototype"
		if ownerNode.IsGetProp() && ownerNode.GetString() == generateExports_PROTOTYPE_PROPERTY {
			// e.g. true if ownerNode is `a.prototype`
			// false if this export were `/** @export */ a.b = obj;` instead
			isEs5StylePrototypeAssignment = true
		}
		propertyName = node.GetString()
	}

	useExportSymbol := true
	if isEs5StylePrototypeAssignment {
		useExportSymbol = false
	} else if methodOwnerName != nil && exports.ContainsKey(*methodOwnerName) {
		useExportSymbol = false
	}

	if useExportSymbol {
		g.AddExportSymbolCall(export, context)
	} else {
		g.AddExportPropertyCall(methodOwnerName, context, export, propertyName)
	}
}

// Java: GenerateExports.addExportPropertyCall(String, Node, String, String)
func (g *GenerateExports) AddExportPropertyCall(
	methodOwnerName *string, context *rhino.Node, export string, propertyName string) {
	// JS output: exportProperty(object, publicName, symbol);
	jutil.CheckNotNull(methodOwnerName)
	call :=
		rhino.IR_Call(
			NodeUtil_NewQNameWithOriginalName(
				g.compiler, *g.exportPropertyFunction,
				context, *g.exportPropertyFunction),
			NodeUtil_NewQNameWithOriginalName(
				g.compiler, *methodOwnerName,
				context, *g.exportPropertyFunction),
			rhino.IR_String(propertyName),
			NodeUtil_NewQNameWithOriginalName(
				g.compiler, export,
				context, *g.exportPropertyFunction))

	expression := rhino.IR_ExprResult(call).SrcrefTreeIfMissing(context)

	g.AddStatement(context, expression)
}

// Java: GenerateExports.addExportSymbolCall(String, Node)
func (g *GenerateExports) AddExportSymbolCall(export string, context *rhino.Node) {
	// JS output: exportSymbol(publicPath, object);
	g.RecordExportSymbol(export)

	call :=
		rhino.IR_Call(
			NodeUtil_NewQNameWithOriginalName(
				g.compiler, *g.exportSymbolFunction,
				context, export),
			rhino.IR_String(export),
			NodeUtil_NewQNameWithOriginalName(
				g.compiler, export,
				context, export))

	expression := rhino.IR_ExprResult(call).SrcrefTreeIfMissing(context)

	g.AddStatement(context, expression)
}

// Java: GenerateExports.addStatement(Node, Node)
func (g *GenerateExports) AddStatement(context *rhino.Node, stmt *rhino.Node) {
	convention := g.compiler.GetCodingConvention()

	n := context
	exprRoot := n
	for !NodeUtil_IsStatementBlock(exprRoot.GetParent()) {
		exprRoot = exprRoot.GetParent()
	}

	// It's important that any class-building calls (goog.inherits)
	// come right after the class definition, so move the export after that.
	for {
		next := exprRoot.GetNext()
		if next != nil &&
			NodeUtil_IsExprCall(next) &&
			convention.GetClassesDefinedByCall(next.GetFirstChild()) != nil {
			exprRoot = next
		} else {
			break
		}
	}

	stmt.InsertAfter(exprRoot)
	g.compiler.ReportChangeToEnclosingScope(stmt)
}

// GetSynthesizedExternsRoot lazily creates a "new" externs root for
// undeclared variables.
//
// Java: GenerateExports.getSynthesizedExternsRoot()
func (g *GenerateExports) GetSynthesizedExternsRoot() *rhino.Node {
	return g.compiler.GetSynthesizedExternsInput().GetAstRoot(g.compiler)
}
