// Port of com/google/javascript/jscomp/ConvertChunksToESModules.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ConvertChunksToESModules finds all references to global symbols in a
// different output chunk and adds ES Module imports and exports for them.
type ConvertChunksToESModules struct {
	compiler               AbstractCompiler
	exportChunks           []*JSChunk
	exportsByChunk         map[*JSChunk]*chunkStringSet
	importChunks           []*JSChunk
	importsByChunk         map[*JSChunk]*chunkImportMap
	dynamicImportCallbacks []rhino.Node
}

type chunkImportType uint8

const (
	chunkImportStatic chunkImportType = iota
	chunkImportDynamic
)

type chunkStringSet struct {
	names []string
	seen  map[string]struct{}
}

func (s *chunkStringSet) add(name string) {
	if s.seen == nil {
		s.seen = make(map[string]struct{})
	}
	if _, ok := s.seen[name]; !ok {
		s.seen[name] = struct{}{}
		s.names = append(s.names, name)
	}
}

type chunkImportMap struct {
	chunks  []*JSChunk
	byChunk map[*JSChunk]*chunkStringSet
}

func (m *chunkImportMap) ensure(chunk *JSChunk) *chunkStringSet {
	if s, ok := m.byChunk[chunk]; ok {
		return s
	}
	if m.byChunk == nil {
		m.byChunk = make(map[*JSChunk]*chunkStringSet)
	}
	s := &chunkStringSet{}
	m.byChunk[chunk] = s
	m.chunks = append(m.chunks, chunk)
	return s
}

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

// NewConvertChunksToESModules creates a new ConvertChunksToESModules compiler pass.
//
// Java: ConvertChunksToESModules(AbstractCompiler)
func NewConvertChunksToESModules(compiler AbstractCompiler) *ConvertChunksToESModules {
	return &ConvertChunksToESModules{
		compiler:       compiler,
		exportsByChunk: make(map[*JSChunk]*chunkStringSet),
		importsByChunk: make(map[*JSChunk]*chunkImportMap),
	}
}

func (c *ConvertChunksToESModules) ensureExports(chunk *JSChunk) *chunkStringSet {
	if s, ok := c.exportsByChunk[chunk]; ok {
		return s
	}
	s := &chunkStringSet{}
	c.exportsByChunk[chunk] = s
	c.exportChunks = append(c.exportChunks, chunk)
	return s
}

func (c *ConvertChunksToESModules) ensureImports(chunk *JSChunk) *chunkImportMap {
	if m, ok := c.importsByChunk[chunk]; ok {
		return m
	}
	m := &chunkImportMap{}
	c.importsByChunk[chunk] = m
	c.importChunks = append(c.importChunks, chunk)
	return m
}

// Java: ConvertChunksToESModules.process(Node, Node)
func (c *ConvertChunksToESModules) Process(externs rhino.Node, root rhino.Node) {
	// Find global names that are used in more than one chunk. Those that
	// are have to have import and export statements added.
	NodeTraversal_Traverse(c.compiler, root, &findCrossChunkReferences{outer: c})

	// Force every output chunk to parse as an ES Module. If a chunk has no imports and
	// no exports, add an empty export list to generate an empty export statement:
	// example: export {};
	for _, chunk := range c.compiler.GetChunkGraph().GetAllChunks() {
		_, hasExports := c.exportsByChunk[chunk]
		_, hasImports := c.importsByChunk[chunk]
		if !hasExports && !hasImports && len(chunk.GetInputs()) > 0 {
			c.ensureExports(chunk)
		}
	}

	c.convertChunkSourcesToModules()
	c.addExportStatements()
	c.addImportStatements()
	c.rewriteDynamicImportCallbacks()

	NodeUtil_AddFeatureToAllScripts(root, parser.FeatureSet_Feature_MODULES, c.compiler)
}

// convertChunkSourcesToModules moves all code in a chunk into the first
// input and marks it as an ESModule. At this point in the compilation, all
// input files should be scripts.
func (c *ConvertChunksToESModules) convertChunkSourcesToModules() {
	for _, chunk := range c.compiler.GetChunkGraph().GetAllChunks() {
		inputs := chunk.GetInputs()
		if len(inputs) == 0 {
			continue
		}

		var firstInput *CompilerInput
		for _, input := range inputs {
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
				combinedFeatureSet := firstInputScriptFeatures.Union(scriptFeatures)
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

func (c *ConvertChunksToESModules) addExportStatements() {
	for _, chunk := range c.exportChunks {
		exports := c.exportsByChunk[chunk]
		firstInput := chunk.GetFirst()
		moduleBody := firstInput.GetAstRoot(c.compiler).GetFirstChild()
		jutil.CheckState(!moduleBody.IsNil() && moduleBody.IsModuleBody())
		exportSpecs := rhino.NewNode(rhino.Token_EXPORT_SPECS)
		for _, name := range exports.names {
			exportSpec := rhino.NewNode(rhino.Token_EXPORT_SPEC)
			exportSpec.AddChildToFront(rhino.IR_Name(name))
			exportSpec.AddChildToFront(rhino.IR_Name(name))
			exportSpec.PutIntProp(rhino.Node_IS_SHORTHAND_PROPERTY, 1)
			exportSpecs.AddChildToBack(exportSpec)
		}
		importsByChunk := c.importsByChunk[chunk]

		// Force the chunk to parse as a module by adding an empty export spec when no actual
		// static imports or exports exist
		if exportSpecs.HasChildren() || importsByChunk == nil || len(importsByChunk.chunks) == 0 {
			export := rhino.IR_Export(exportSpecs).SrcrefTree(moduleBody)
			moduleBody.AddChildToBack(export)
			c.compiler.ReportChangeToEnclosingScope(moduleBody)
		}
	}
}

func chunkFileName(chunk *JSChunk) string {
	return chunk.GetName() + ".js"
}

func isValidChunkPathString(s string) bool {
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b == 0 {
			return false
		}
		if b == 0xED && i+1 < len(s) && s[i+1] >= 0xA0 && s[i+1] <= 0xBF {
			return false
		}
	}
	return true
}

func relativeChunkImportPath(fromUriPath, toUriPath string) (string, bool) {
	if !isValidChunkPathString(fromUriPath) || !isValidChunkPathString(toUriPath) {
		return "", false
	}
	fromPath := jutil.Path_Of(fromUriPath)
	toPath := jutil.Path_Of(toUriPath)
	fromFolder := fromPath.GetParent()
	if fromFolder == nil {
		if toPath.GetParent() == nil {
			return "./" + toUriPath, true
		}
		if strings.HasPrefix(toUriPath, ".") ||
			strings.HasPrefix(toPath.String(), "/") ||
			strings.HasPrefix(toPath.String(), "\\") {
			return toUriPath, true
		}
		return "", false
	}
	if fromFolder.IsAbsolute() != toPath.IsAbsolute() {
		return "", false
	}
	if !fromFolder.Equals(toPath) && fromFolder.String() != "" {
		baseStr := fromFolder.Normalize().String()
		if baseStr == ".." || strings.HasPrefix(baseStr, "../") {
			targetStr := toPath.Normalize().String()
			remBase, remTarget := baseStr, targetStr
			for remBase == ".." || strings.HasPrefix(remBase, "../") {
				if remTarget == ".." {
					remTarget = ""
				} else if after, ok := strings.CutPrefix(remTarget, "../"); ok {
					remTarget = after
				} else {
					return "", false
				}
				if remBase == ".." {
					remBase = ""
				} else {
					remBase = remBase[3:]
				}
			}
		}
	}
	calculatedPath := fromFolder.Relativize(toPath).String()
	if strings.HasPrefix(calculatedPath, ".") || strings.HasPrefix(calculatedPath, "/") {
		return calculatedPath, true
	}
	return "./" + calculatedPath, true
}

func (c *ConvertChunksToESModules) addImportStatements() {
	var importStatements []rhino.Node
	for _, importingChunk := range c.importChunks {
		chunkImports := c.importsByChunk[importingChunk]
		importStatements = importStatements[:0]
		firstInput := importingChunk.GetFirst()
		moduleBody := firstInput.GetAstRoot(c.compiler).GetFirstChild()
		jutil.CheckState(!moduleBody.IsNil() && moduleBody.IsModuleBody())
		importingName := chunkFileName(importingChunk)

		// For each distinct chunk where a referenced symbol is defined, create an import statement
		// referencing the names.
		for _, exportingChunk := range chunkImports.chunks {
			importsForChunk := chunkImports.byChunk[exportingChunk]
			importSpecs := rhino.NewNode(rhino.Token_IMPORT_SPECS)
			for _, name := range importsForChunk.names {
				importSpec := rhino.NewNode(rhino.Token_IMPORT_SPEC)
				importSpec.AddChildToFront(rhino.IR_Name(name))
				importSpec.AddChildToFront(rhino.IR_Name(name))
				importSpec.PutIntProp(rhino.Node_IS_SHORTHAND_PROPERTY, 1)
				importSpecs.AddChildToBack(importSpec)
			}
			importStatement := rhino.NewNode(rhino.Token_IMPORT)
			exportingName := chunkFileName(exportingChunk)
			importPath, ok := relativeChunkImportPath(importingName, exportingName)
			if !ok {
				importPath = exportingName
				c.compiler.Report(
					JSError_MakeWithNode(
						moduleBody,
						ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH,
						importingName,
						exportingName))
			}
			importStatement.AddChildToFront(rhino.IR_String(importPath))
			if importSpecs.HasChildren() {
				importStatement.AddChildToFront(importSpecs)
			} else {
				// Empty import of a dependent chunk for side effects: import './chunk.js'
				importStatement.AddChildToFront(rhino.IR_Empty())
			}
			importStatement.AddChildToFront(rhino.IR_Empty())

			importStatement.SrcrefTree(moduleBody)
			importStatements = append(importStatements, importStatement)
		}
		for i := len(importStatements) - 1; i >= 0; i-- {
			moduleBody.AddChildToFront(importStatements[i])
		}
		c.compiler.ReportChangeToEnclosingScope(moduleBody)
	}
}

// ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace finds and
// returns the module namespace name node in a dynamic import callback function.
//
// Java: ConvertChunksToESModules.getDynamicImportCallbackModuleNamespace(AbstractCompiler, Node)
func ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(
	compiler AbstractCompiler, call rhino.Node) rhino.Node {
	_, isCall := ast.AsCall(call)
	jutil.CheckState(isCall)
	callbackFnNode := NodeUtil_GetArgumentForCallOrNew(call, 0)
	callbackFn, isFn := ast.AsFunction(callbackFnNode)
	if !isFn || callbackFn.Params().Node().HasChildren() {
		compiler.Report(
			JSError_MakeWithNode(
				call,
				ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK,
				" Unable to find valid callback function."))
		return rhino.NilNode
	}
	callbackBody := callbackFn.Body().Node()

	// The callback body should have a single statement that returns a name.
	// Support both standard and arrow function semantics
	if _, ok := ast.AsName(callbackBody); ok {
		return callbackBody
	}
	if block, ok := ast.AsBlock(callbackBody); ok && block.Node().HasOneChild() {
		if ret, isRet := ast.AsReturn(block.Node().GetFirstChild()); isRet {
			if val, hasVal := ret.Value(); hasVal {
				if _, isName := ast.AsName(val.Node()); isName {
					return val.Node()
				}
			}
		}
	}
	compiler.Report(
		JSError_MakeWithNode(
			call,
			ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK,
			" Unable to find valid namespace reference."))
	return rhino.NilNode
}

// rewriteDynamicImportCallbacks removes the special external function call injected
// by RewriteDynamicImports and updates the original module namespace reference to be
// a property of the chunk namespace:
//
//	import('./chunk0.js').then(jscomp$DynamicImportCallback(() => module$input0));
//	-> import('./chunk0.js').then(($) => $.module$input0);
func (c *ConvertChunksToESModules) rewriteDynamicImportCallbacks() {
	astFactory := c.compiler.CreateAstFactory()
	for _, dynamicImportCallback := range c.dynamicImportCallbacks {
		jutil.CheckState(dynamicImportCallback.IsCall())
		moduleNamespace := ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(c.compiler, dynamicImportCallback)
		if moduleNamespace.IsNil() {
			continue
		}
		callbackFn := NodeUtil_GetArgumentForCallOrNew(dynamicImportCallback, 0)
		callbackParamList := NodeUtil_GetFunctionParameters(callbackFn)
		importNamespaceParam := astFactory.CreateNameWithUnknownType("$").Srcref(moduleNamespace)
		callbackParamList.AddChildToFront(importNamespaceParam)
		c.compiler.ReportChangeToEnclosingScope(importNamespaceParam)

		namespaceGetprop := astFactory.CreateGetPropWithUnknownType(
			astFactory.CreateNameWithUnknownType("$").Srcref(moduleNamespace),
			moduleNamespace.GetString())

		moduleNamespace.ReplaceWith(namespaceGetprop)
		c.compiler.ReportChangeToEnclosingScope(namespaceGetprop)
		innerCallback := callbackFn.Detach()
		dynamicImportCallback.ReplaceWith(innerCallback)
		c.compiler.ReportChangeToEnclosingScope(innerCallback)
	}
}

type findCrossChunkReferences struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer *ConvertChunksToESModules
}

func (f *findCrossChunkReferences) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		f.visitScript(t, n)
	case rhino.Token_CALL:
		return f.outer.visitCallAndTraverse(t, n)
	case rhino.Token_NAME:
		f.outer.visitName(t, n, chunkImportStatic)
	}
	return true
}

func (f *findCrossChunkReferences) visitScript(t *NodeTraversal, script rhino.Node) {
	jutil.CheckState(script.IsScript())
	chunk := t.GetChunk()
	chunkDependencies := chunk.GetDependencies()

	f.outer.ensureExports(chunk)

	// Ensure every chunk dependency is explicitly listed with an import.
	// Dependent chunks may have side effects even if there isn't an explicit name reference.
	if len(chunkDependencies) > 0 {
		namesToImportByChunk := f.outer.ensureImports(chunk)
		for _, dependency := range chunkDependencies {
			namesToImportByChunk.ensure(dependency)
		}
	}
}

// ConvertChunksToESModules_IsDynamicImportCallback tests if a node is a
// .then callback as inserted by the RewriteDynamicImports pass.
//
// Java: ConvertChunksToESModules.isDynamicImportCallback(Node)
func ConvertChunksToESModules_IsDynamicImportCallback(call rhino.Node) bool {
	return call.IsCall() && NodeUtil_IsCallTo(call, ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN)
}

func (c *ConvertChunksToESModules) visitCallAndTraverse(t *NodeTraversal, call rhino.Node) bool {
	jutil.CheckState(call.IsCall())
	if !ConvertChunksToESModules_IsDynamicImportCallback(call) {
		return true
	}

	moduleNamespace := ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(c.compiler, call)
	if moduleNamespace.IsNil() {
		return true
	}
	if c.visitName(t, moduleNamespace, chunkImportDynamic) {
		c.dynamicImportCallbacks = append(c.dynamicImportCallbacks, call)
	} else {
		c.compiler.Report(
			JSError_MakeWithNode(
				call,
				ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK,
				" Unable to find valid namespace reference."))
	}
	return false
}

func (c *ConvertChunksToESModules) visitName(t *NodeTraversal, nameNode rhino.Node, importType chunkImportType) bool {
	jutil.CheckState(nameNode.IsName())
	name := nameNode.GetString()
	if name == "" {
		return false
	}

	v := t.GetScope().GetVar(name)
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
				name,
				chunkFileName(referencingChunk),
				chunkFileName(definingChunk))
		}

		// Mark the chunk where the name is declared as needing an export for this name
		c.ensureExports(definingChunk).add(name)

		// Add an import for this name to this chunk from the source chunk
		namesToImportByChunk := c.ensureImports(referencingChunk)
		if importType == chunkImportStatic {
			namesToImportByChunk.ensure(definingChunk).add(name)
		}
	}
	return true
}
