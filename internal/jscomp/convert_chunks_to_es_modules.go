// Port of com/google/javascript/jscomp/ConvertChunksToESModules.java @ closure-compiler v20261015

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ConvertChunksToESModules finds all references to global symbols in a
// different output chunk and add ES Module imports and exports for them.
//
//	// chunk1.js
//	var a = 1;
//	function b() { return a }
//
//	// chunk2.js
//	console.log(a);
//
// becomes
//
//	// chunk1.js
//	var a = 1;
//	var b = function b() { return a };
//	export {a};
//
//	// chunk2.js
//	import {a} from './chunk1.js';
//	console.log(a);
//
// This allows splitting code into es modules that depend on each other's
// symbols, without using a global namespace or polluting the global scope.
type ConvertChunksToESModules struct {
	compiler               AbstractCompiler
	crossChunkExports      *jutil.LinkedHashMap[*JSChunk, *jutil.LinkedHashSet[string]]
	crossChunkImports      *jutil.LinkedHashMap[*JSChunk, *jutil.LinkedHashMap[*JSChunk, *jutil.LinkedHashSet[string]]]
	dynamicImportCallbacks *jutil.ArrayList[rhino.Node]
}

// ConvertChunksToESModules_ImportType is the ImportType enum.
type ConvertChunksToESModules_ImportType int32

const (
	ConvertChunksToESModules_ImportType_STATIC ConvertChunksToESModules_ImportType = iota
	ConvertChunksToESModules_ImportType_DYNAMIC
)

var convertChunksToESModules_ImportType_names = [...]string{"STATIC", "DYNAMIC"}

// ConvertChunksToESModules_ImportType_values is ImportType.values().
func ConvertChunksToESModules_ImportType_values() []ConvertChunksToESModules_ImportType {
	return []ConvertChunksToESModules_ImportType{
		ConvertChunksToESModules_ImportType_STATIC, ConvertChunksToESModules_ImportType_DYNAMIC,
	}
}

// Name is Enum.name().
func (t ConvertChunksToESModules_ImportType) Name() string {
	return convertChunksToESModules_ImportType_names[t]
}

// String is Enum.toString().
func (t ConvertChunksToESModules_ImportType) String() string { return t.Name() }

// Ordinal is Enum.ordinal().
func (t ConvertChunksToESModules_ImportType) Ordinal() int { return int(t) }

const ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN = "jscomp$DynamicImportCallback"

var ConvertChunksToESModules_ASSIGNMENT_TO_IMPORT = DiagnosticType_Error(
	"JSC_IMPORT_ASSIGN",
	"Imported symbol \"{0}\" in chunk \"{1}\" cannot be assigned (defined in \"{2}\")")

var ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH = DiagnosticType_Error(
	"JSC_UNABLE_TO_COMPUTE_RELATIVE_PATH",
	"Unable to compute relative import path from \"{0}\" to \"{1}\"")

var ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK = DiagnosticType_Error(
	"JSC_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK",
	"Dynamic import callback encountered wih an invalid format.{0}")

// NewConvertChunksToESModules is the constructor for the
// ConvertChunksToESModules compiler pass.
//
// compiler is the JSCompiler, for reporting code changes.
//
// Java: ConvertChunksToESModules(AbstractCompiler)
func NewConvertChunksToESModules(compiler AbstractCompiler) *ConvertChunksToESModules {
	c := &ConvertChunksToESModules{}
	c.crossChunkExports = jutil.NewLinkedHashMap[*JSChunk, *jutil.LinkedHashSet[string]]()
	c.crossChunkImports = jutil.NewLinkedHashMap[*JSChunk, *jutil.LinkedHashMap[*JSChunk, *jutil.LinkedHashSet[string]]]()
	c.dynamicImportCallbacks = jutil.NewArrayList[rhino.Node]()
	c.compiler = compiler
	return c
}

// Java: ConvertChunksToESModules.process(Node, Node)
func (c *ConvertChunksToESModules) Process(externs rhino.Node, root rhino.Node) {
	// Find global names that are used in more than one chunk. Those that
	// are have to have import and export statements added.
	NodeTraversal_Traverse(c.compiler, root, newConvertChunksToESModules_FindCrossChunkReferences(c))

	// Force every output chunk to parse as an ES Module. If a chunk has no imports and
	// no exports, add an empty export list to generate an empty export statement:
	// example: export {};
	for _, chunk := range c.compiler.GetChunkGraph().GetAllChunks() {
		if !c.crossChunkExports.ContainsKey(chunk) &&
			!c.crossChunkImports.ContainsKey(chunk) &&
			!(len(chunk.GetInputs()) == 0) {
			c.crossChunkExports.Put(chunk, jutil.NewLinkedHashSet[string]())
		}
	}

	c.ConvertChunkSourcesToModules()
	c.AddExportStatements()
	c.AddImportStatements()
	c.RewriteDynamicImportCallbacks()

	NodeUtil_AddFeatureToAllScripts(root, parser.FeatureSet_Feature_MODULES, c.compiler)
}

// ConvertChunkSourcesToModules moves all code in a chunk into the first
// input and mark it as an ESModule. At this point in the compilation, all
// input files should be scripts.
//
// Java: ConvertChunksToESModules.convertChunkSourcesToModules()
func (c *ConvertChunksToESModules) ConvertChunkSourcesToModules() {
	for _, chunk := range c.compiler.GetChunkGraph().GetAllChunks() {
		if len(chunk.GetInputs()) == 0 {
			continue
		}

		var firstInput *CompilerInput = nil
		for _, input := range chunk.GetInputs() {
			astRoot := input.GetAstRoot(c.compiler)
			scriptFeatures := NodeUtil_GetFeatureSetOfScript(astRoot)
			jutil.CheckState(!scriptFeatures.Contains(parser.FeatureSet_ES2015_MODULES))
			if firstInput == nil {
				firstInput = input
				scriptFeatures = scriptFeatures.Union(parser.FeatureSet_ES2015_MODULES)
				astRoot.PutProp(rhino.Node_FEATURE_SET, scriptFeatures)
				moduleBody := rhino.NewNode(rhino.Token_MODULE_BODY)
				moduleBody.Srcref(astRoot)
				moduleBody.AddChildrenToFront(astRoot.RemoveChildren())
				astRoot.AddChildToFront(moduleBody)
				c.compiler.ReportChangeToEnclosingScope(moduleBody)
			} else {
				firstInputAstRoot := firstInput.GetAstRoot(c.compiler)
				firstInputScriptFeatures := NodeUtil_GetFeatureSetOfScript(firstInputAstRoot)
				combinedFeatureSet :=
					firstInputScriptFeatures.Union(NodeUtil_GetFeatureSetOfScript(astRoot))
				astRoot.PutProp(rhino.Node_FEATURE_SET, combinedFeatureSet)
				moduleBody := firstInputAstRoot.GetFirstChild()
				jutil.CheckState(!moduleBody.IsNil() && moduleBody.IsModuleBody())
				moduleBody.AddChildrenToBack(astRoot.RemoveChildren())
				c.compiler.ReportChangeToEnclosingScope(firstInputAstRoot)
				c.compiler.ReportChangeToChangeScope(astRoot)
			}
		}
	}
}

// AddExportStatements adds export statements to chunks
//
// Java: ConvertChunksToESModules.addExportStatements()
func (c *ConvertChunksToESModules) AddExportStatements() {
	for jsModuleExportsKey, jsModuleExportsValue := range c.crossChunkExports.All() {
		firstInput := jsModuleExportsKey.GetFirst()
		moduleBody := firstInput.GetAstRoot(c.compiler).GetFirstChild()
		jutil.CheckState(!moduleBody.IsNil() && moduleBody.IsModuleBody())
		exportSpecs := rhino.NewNode(rhino.Token_EXPORT_SPECS)
		for name := range jsModuleExportsValue.All() {
			exportSpec := rhino.NewNode(rhino.Token_EXPORT_SPEC)
			exportSpec.AddChildToFront(rhino.IR_Name(name))
			exportSpec.AddChildToFront(rhino.IR_Name(name))
			exportSpec.PutIntProp(rhino.Node_IS_SHORTHAND_PROPERTY, 1)
			exportSpecs.AddChildToBack(exportSpec)
		}
		importsByChunk := c.crossChunkImports.Get(jsModuleExportsKey)

		// Force the chunk to parse as a module by adding an empty export spec when no actual
		// static imports or exports exist
		if exportSpecs.HasChildren() || importsByChunk == nil || importsByChunk.IsEmpty() {
			export := rhino.IR_Export(exportSpecs).SrcrefTree(moduleBody)
			moduleBody.AddChildToBack(export)
			c.compiler.ReportChangeToEnclosingScope(moduleBody)
		}
	}
}

// Java: ConvertChunksToESModules.getChunkName(JSChunk)
func ConvertChunksToESModules_GetChunkName(chunk *JSChunk) string {
	return chunk.GetName() + ".js"
}

// AddImportStatements adds import statements to chunks
//
// Java: ConvertChunksToESModules.addImportStatements()
func (c *ConvertChunksToESModules) AddImportStatements() {
	for importingChunk, chunkImportsValue := range c.crossChunkImports.All() {
		importStatements := jutil.NewArrayList[rhino.Node]()
		firstInput := importingChunk.GetFirst()
		moduleBody := firstInput.GetAstRoot(c.compiler).GetFirstChild()
		jutil.CheckState(!moduleBody.IsNil() && moduleBody.IsModuleBody())

		// For each distinct chunk where a referenced symbol is defined, create an import statement
		// referencing the names.
		for exportingChunk, importsByChunkValue := range chunkImportsValue.All() {
			importSpecs := rhino.NewNode(rhino.Token_IMPORT_SPECS)
			for name := range importsByChunkValue.All() {
				importSpec := rhino.NewNode(rhino.Token_IMPORT_SPEC)
				importSpec.AddChildToFront(rhino.IR_Name(name))
				importSpec.AddChildToFront(rhino.IR_Name(name))
				importSpec.PutIntProp(rhino.Node_IS_SHORTHAND_PROPERTY, 1)
				importSpecs.AddChildToBack(importSpec)
			}
			importStatement := rhino.NewNode(rhino.Token_IMPORT)
			importPath := ConvertChunksToESModules_GetChunkName(exportingChunk)
			jutil.Try(func() {
				importPath =
					ModuleLoader_RelativePathFrom(
						ConvertChunksToESModules_GetChunkName(importingChunk), ConvertChunksToESModules_GetChunkName(exportingChunk))
			}, func(e jutil.IIllegalArgumentException) {
				c.compiler.Report(
					JSError_MakeWithNode(
						moduleBody,
						ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH,
						ConvertChunksToESModules_GetChunkName(importingChunk),
						ConvertChunksToESModules_GetChunkName(exportingChunk)))
			})
			importStatement.AddChildToFront(rhino.IR_String(importPath))
			if importSpecs.HasChildren() {
				importStatement.AddChildToFront(importSpecs)
			} else {
				// Empty import of a dependent chunk for side effects
				// import './chunk.js'
				importStatement.AddChildToFront(rhino.IR_Empty())
			}
			importStatement.AddChildToFront(rhino.IR_Empty())

			importStatement.SrcrefTree(moduleBody)
			importStatements.AddAt(0, importStatement)
		}
		for importStatement := range importStatements.All() {
			moduleBody.AddChildToFront(importStatement)
		}
		c.compiler.ReportChangeToEnclosingScope(moduleBody)
	}
}

// ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace finds and
// returns the module namespace name node in a dynamic import callback
// function
//
// Java: ConvertChunksToESModules.getDynamicImportCallbackModuleNamespace(AbstractCompiler, Node)
func ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(
	compiler AbstractCompiler, call rhino.Node) rhino.Node {
	jutil.CheckState(call.IsCall())
	callbackFn := NodeUtil_GetArgumentForCallOrNew(call, 0)
	if callbackFn.IsNil() ||
		!callbackFn.IsFunction() ||
		NodeUtil_GetFunctionParameters(callbackFn).HasChildren() {
		compiler.Report(
			JSError_MakeWithNode(
				call,
				ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK,
				" Unable to find valid callback function."))
		return rhino.NilNode
	}
	callbackBody := NodeUtil_GetFunctionBody(callbackFn)

	// The callback body should have a single statement that returns a name.
	// Support both standard and arrow function semantics
	if callbackBody.IsName() {
		return callbackBody
	} else if callbackBody.IsBlock() &&
		callbackBody.HasOneChild() &&
		callbackBody.GetFirstChild().IsReturn() &&
		callbackBody.GetFirstChild().HasOneChild() &&
		callbackBody.GetFirstFirstChild().IsName() {
		return callbackBody.GetFirstFirstChild()
	}
	compiler.Report(
		JSError_MakeWithNode(
			call,
			ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK,
			" Unable to find valid namespace reference."))
	return rhino.NilNode
}

// RewriteDynamicImportCallbacks: the RewriteDynamicImports pass adds direct
// references to the original input module namespace and wraps the callback
// in a special external function call so that this pass can recognize
// them.
//
// Example:
//
// import('./chunk0.js').then(jscomp$DynamicImportCallback(() => module$input0));
//
// We need to remove the special external function call and update the
// original module namespace reference to be a property of the chunk
// namespace.
//
// import('./chunk0.js').then(($) => $.module$input0);
//
// Java: ConvertChunksToESModules.rewriteDynamicImportCallbacks()
func (c *ConvertChunksToESModules) RewriteDynamicImportCallbacks() {
	astFactory := c.compiler.CreateAstFactory()
	for dynamicImportCallback := range c.dynamicImportCallbacks.All() {
		jutil.CheckState(dynamicImportCallback.IsCall())
		moduleNamespace :=
			ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(c.compiler, dynamicImportCallback)
		if moduleNamespace.IsNil() {
			continue
		}
		callbackFn := NodeUtil_GetArgumentForCallOrNew(dynamicImportCallback, 0)
		callbackParamList := NodeUtil_GetFunctionParameters(callbackFn)
		importNamespaceParam := astFactory.CreateNameWithUnknownType("$").Srcref(moduleNamespace)
		callbackParamList.AddChildToFront(importNamespaceParam)
		c.compiler.ReportChangeToEnclosingScope(importNamespaceParam)

		namespaceGetprop :=
			astFactory.CreateGetPropWithUnknownType(
				astFactory.CreateNameWithUnknownType("$").Srcref(moduleNamespace),
				moduleNamespace.GetString())

		moduleNamespace.ReplaceWith(namespaceGetprop)
		c.compiler.ReportChangeToEnclosingScope(namespaceGetprop)
		innerCallback := NodeUtil_GetArgumentForCallOrNew(dynamicImportCallback, 0).Detach()
		dynamicImportCallback.ReplaceWith(innerCallback)
		c.compiler.ReportChangeToEnclosingScope(innerCallback)
	}
}

// ConvertChunksToESModules_FindCrossChunkReferences finds names in a chunk
// that are defined in a different chunk.
type ConvertChunksToESModules_FindCrossChunkReferences struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	outer *ConvertChunksToESModules
}

// Java: ConvertChunksToESModules.FindCrossChunkReferences()
func newConvertChunksToESModules_FindCrossChunkReferences(outer *ConvertChunksToESModules) *ConvertChunksToESModules_FindCrossChunkReferences {
	return &ConvertChunksToESModules_FindCrossChunkReferences{outer: outer}
}

// Java: ConvertChunksToESModules.FindCrossChunkReferences.shouldTraverse(NodeTraversal, Node, Node)
func (f *ConvertChunksToESModules_FindCrossChunkReferences) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		f.VisitScript(t, n)
		return true
	} else if n.IsCall() {
		return f.outer.VisitCallAndTraverse(t, n)
	} else if n.IsName() {
		f.outer.VisitName(t, n, ConvertChunksToESModules_ImportType_STATIC)
		return true
	}
	return true
}

// Java: ConvertChunksToESModules.FindCrossChunkReferences.visitScript(NodeTraversal, Node)
func (f *ConvertChunksToESModules_FindCrossChunkReferences) VisitScript(t *NodeTraversal, script rhino.Node) {
	jutil.CheckState(script.IsScript())
	chunk := t.GetChunk()
	chunkDependencies := chunk.GetDependencies()

	f.outer.crossChunkExports.PutIfAbsent(chunk, jutil.NewLinkedHashSet[string]())

	// Ensure every chunk dependency is explicitly listed with an import
	// Dependent chunks may have side effects even if there isn't an explicit name reference
	if !(len(chunkDependencies) == 0) {
		namesToImportByChunk :=
			f.outer.crossChunkImports.ComputeIfAbsent(chunk, func(k *JSChunk) *jutil.LinkedHashMap[*JSChunk, *jutil.LinkedHashSet[string]] {
				return jutil.NewLinkedHashMap[*JSChunk, *jutil.LinkedHashSet[string]]()
			})
		for _, dependency := range chunkDependencies {
			namesToImportByChunk.ComputeIfAbsent(dependency, func(k *JSChunk) *jutil.LinkedHashSet[string] {
				return jutil.NewLinkedHashSet[string]()
			})
		}
	}
}

// ConvertChunksToESModules_IsDynamicImportCallback tests if a node is a
// .then callback as inserted by the RewriteDynamicImports pass Only finds
// callbacks when wrapped in the specially named extern function injected
// when the chunk output type is ES_MODULES
//
// Java: ConvertChunksToESModules.isDynamicImportCallback(Node)
func ConvertChunksToESModules_IsDynamicImportCallback(call rhino.Node) bool {
	if !call.IsCall() {
		return false
	}
	if NodeUtil_IsCallTo(call, ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN) {
		return true
	}
	return false
}

// VisitCallAndTraverse tests if a node is a .then callback as inserted by
// the RewriteDynamicImports pass
//
// Java: ConvertChunksToESModules.visitCallAndTraverse(NodeTraversal, Node)
func (c *ConvertChunksToESModules) VisitCallAndTraverse(t *NodeTraversal, call rhino.Node) bool {
	jutil.CheckState(call.IsCall())
	if !ConvertChunksToESModules_IsDynamicImportCallback(call) {
		return true
	}

	moduleNamespace := ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(c.compiler, call)
	if moduleNamespace.IsNil() {
		return true
	}
	isValidModuleNamespace := c.VisitName(t, moduleNamespace, ConvertChunksToESModules_ImportType_DYNAMIC)
	if isValidModuleNamespace {
		c.dynamicImportCallbacks.Add(call)
	} else {
		c.compiler.Report(
			JSError_MakeWithNode(
				call,
				ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK,
				" Unable to find valid namespace reference."))
	}
	return false
}

// Java: ConvertChunksToESModules.visitName(NodeTraversal, Node, ImportType)
func (c *ConvertChunksToESModules) VisitName(t *NodeTraversal, nameNode rhino.Node, importType ConvertChunksToESModules_ImportType) bool {
	jutil.CheckState(nameNode.IsName())
	name := nameNode.GetString()

	if name == "" {
		return false
	}

	s := t.GetScope()
	v := s.GetVar(name)
	if v == nil || !v.IsGlobal() {
		return false
	}
	input := v.GetInput()
	if input == nil {
		return false
	}
	definingChunk := input.GetChunk()
	referencingChunk := t.GetChunk()

	if definingChunk != referencingChunk {
		if NodeUtil_IsLhsOfAssign(nameNode) {
			t.Report(
				nameNode,
				ConvertChunksToESModules_ASSIGNMENT_TO_IMPORT,
				nameNode.GetString(),
				ConvertChunksToESModules_GetChunkName(referencingChunk),
				ConvertChunksToESModules_GetChunkName(definingChunk))
		}

		// Mark the chunk where the name is declared as needing an export for this name
		namesToExport :=
			c.crossChunkExports.ComputeIfAbsent(definingChunk, func(k *JSChunk) *jutil.LinkedHashSet[string] {
				return jutil.NewLinkedHashSet[string]()
			})
		namesToExport.Add(name)

		// Add an import for this name to this chunk from the source chunk
		namesToImportByChunk :=
			c.crossChunkImports.ComputeIfAbsent(referencingChunk, func(k *JSChunk) *jutil.LinkedHashMap[*JSChunk, *jutil.LinkedHashSet[string]] {
				return jutil.NewLinkedHashMap[*JSChunk, *jutil.LinkedHashSet[string]]()
			})
		if importType == ConvertChunksToESModules_ImportType_STATIC {
			importsForChunk :=
				namesToImportByChunk.ComputeIfAbsent(
					definingChunk, func(k *JSChunk) *jutil.LinkedHashSet[string] {
						return jutil.NewLinkedHashSet[string]()
					})
			importsForChunk.Add(name)
		}
	}
	return true
}
