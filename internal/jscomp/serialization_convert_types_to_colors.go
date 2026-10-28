// Port of com/google/javascript/jscomp/serialization/ConvertTypesToColors.java @ closure-compiler v20261021

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ConvertTypesToColors is a pass to convert JSType objects from TypeChecking
// that are attached to the AST into Color objects whose sole use is to enable
// running optimizations and delete all other references to JSTypes.
//
// This pass is also responsible for logging debug information that needs to
// know about both JSType objects and their corresponding colors, and for
// telling the compiler to prepare for reading runtime libraries from
// precompiled TypedASTs instead of source.
//
// Java: public final class ConvertTypesToColors implements CompilerPass
type ConvertTypesToColors struct {
	compiler             AbstractCompiler
	serializationOptions *SerializationOptions
}

// Java: ConvertTypesToColors(AbstractCompiler, SerializationOptions)
func NewConvertTypesToColors(
	compiler AbstractCompiler, serializationOptions *SerializationOptions) *ConvertTypesToColors {
	c := &ConvertTypesToColors{}
	c.compiler = compiler
	c.serializationOptions = serializationOptions
	return c
}

type convertTypesToColors_removeTypes struct {
	jsdocBuilder rhino.JSDocInfo_Builder
}

func (r *convertTypesToColors_removeTypes) walk(n rhino.Node) {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		r.walk(c)
	}
	_, _, jsdoc := n.TakeTypesAndJSDocForColorConversion()
	if jsdoc != nil {
		n.SetJSDocInfo(jsDocSerializer_convertJSDocInfoWithBuilder(jsdoc, &r.jsdocBuilder))
	}
}

// Java: private static class ConvertTypesToColors.RemoveTypesAndApplyColors extends RemoveTypes
type ConvertTypesToColors_RemoveTypesAndApplyColors struct {
	offsetToColor         []*colors.Color
	jstypeReconserializer *JSTypeReconserializer
	jsdocBuilder          rhino.JSDocInfo_Builder
}

func (r *ConvertTypesToColors_RemoveTypesAndApplyColors) visitNode(n rhino.Node) {
	oldType, hadBeforeCast, jsdoc := n.TakeTypesAndJSDocForColorConversion()
	if jsdoc != nil {
		n.SetJSDocInfo(jsDocSerializer_convertJSDocInfoWithBuilder(jsdoc, &r.jsdocBuilder))
	}

	if oldType != nil && r.jstypeReconserializer != nil {
		if pointer, ok := r.jstypeReconserializer.typePointers[oldType]; ok {
			if c := r.offsetToColor[pointer]; c != nil {
				n.SetColor(c)
			}
		}
	}

	if hadBeforeCast {
		// used by FunctionInjector and InlineVariables when inlining as a hint that a node has a
		// more specific color
		n.SetColorFromTypeCast()
	}
}

func (r *ConvertTypesToColors_RemoveTypesAndApplyColors) walk(n rhino.Node) {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		r.walk(c)
	}
	r.visitNode(n)
}

// Java: ConvertTypesToColors.process(Node, Node)
func (c *ConvertTypesToColors) Process(externs rhino.Node, root rhino.Node) {
	if c.compiler.GetLifeCycleStage().HasColorAndSimplifiedJSDoc() {
		// Pass is a no-op if we already have optimization colors or have finished the checks phase
		return
	}
	// NOTE(lharker): this pass could probably safely run after normalization, except that the
	// LifeCycleStage enum assumes that normalization implies "colors + simplified JSDoc"
	jutil.CheckStateMsg(
		!c.compiler.GetLifeCycleStage().IsNormalized(), "Not expected to run after normalization")

	externsAndJsRoot := root.GetParent()

	if !c.compiler.HasTypeCheckingRun() {
		var remover convertTypesToColors_removeTypes
		remover.walk(externsAndJsRoot)
		convertTypesToColors_forDetachedRoots(c.compiler, externsAndJsRoot, remover.walk)
		c.compiler.ClearJSTypeRegistry()
		c.compiler.InitRuntimeLibraryTypedAsts(nil) // Optional.absent()
		return
	}

	serializeJstypes := SerializeTypesToPointers_Create(c.compiler, nil, c.serializationOptions)
	serializeJstypes.gatherTypesForColors(externsAndJsRoot)

	colorPoolBuilder := ColorPool_Builder_()
	c.compiler.InitRuntimeLibraryTypedAsts(colorPoolBuilder)
	colorPoolBuilder.addJSTypeShard(serializeJstypes.jstypeReconserializer)

	colorPool := colorPoolBuilder.Build()

	var callback ConvertTypesToColors_RemoveTypesAndApplyColors
	callback.jstypeReconserializer = serializeJstypes.jstypeReconserializer
	callback.offsetToColor = serializeJstypes.jstypeReconserializer.offsetToColor

	callback.walk(externsAndJsRoot)
	convertTypesToColors_forDetachedRoots(c.compiler, externsAndJsRoot, callback.walk)

	c.compiler.ClearJSTypeRegistry()
	c.compiler.SetColorRegistry(colorPool.GetRegistry())
}

// convertTypesToColors_forDetachedRoots calls walk on the root of every
// subtree that is no longer attached to externsAndJsRoot but is still
// referenced by the compiler: the change-tracker timeline (function and
// script nodes changed by the checks and never reported deleted), the
// module metadata and module map (goog.loadModule calls, module bodies
// and export/import nodes the module rewriting detached) and the reported
// diagnostics. Java keeps those nodes too, and with them their
// JSTypes; here the types are converted like the live nodes' so that the
// JSType graph (every type points to the registry) can be collected.
//
// Go-only.
func convertTypesToColors_forDetachedRoots(compiler AbstractCompiler, externsAndJsRoot rhino.Node, walk func(rhino.Node)) {
	seen := map[rhino.Node]struct{}{}
	visit := func(n rhino.Node) {
		if n.IsNil() {
			return
		}
		root := n
		for p := root.GetParent(); !p.IsNil(); p = p.GetParent() {
			root = p
		}
		if root == externsAndJsRoot {
			return
		}
		if _, ok := seen[root]; ok {
			return
		}
		seen[root] = struct{}{}
		walk(root)
	}
	compiler.GetChangeTracker().changeTimeline.Values(func(n rhino.Node) bool {
		visit(n)
		return true
	})
	if mm := compiler.GetModuleMetadataMap(); mm != nil {
		for m := range mm.GetAllModuleMetadata().All() {
			visit(m.RootNode())
		}
	}
	if mm := compiler.GetModuleMap(); mm != nil {
		visitExport := func(e *Export) {
			if e != nil {
				visit(e.ExportNode())
				visit(e.NameNode())
			}
		}
		visitBindings := func(bindings *jutil.LinkedHashMap[string, *Binding]) {
			for b := range bindings.Values() {
				visit(b.SourceNode())
				visitExport(b.OriginatingExport())
			}
		}
		visitModule := func(m *Module) {
			visit(m.Metadata().RootNode())
			visitBindings(m.Namespace())
			visitBindings(m.BoundNames())
			for e := range m.LocalNameToLocalExport().Values() {
				visitExport(e)
			}
		}
		for m := range mm.GetModulesByPath().Values() {
			visitModule(m)
		}
		for m := range mm.GetModulesByClosureNamespace().Values() {
			visitModule(m)
		}
	}
	// Diagnostics keep their node for deferred formatting; the checks report
	// on nodes they then remove (goog.require calls, casts, ...).
	em := compiler.GetErrorManager()
	for _, e := range em.GetErrors() {
		visit(e.Node())
	}
	for _, e := range em.GetWarnings() {
		visit(e.Node())
	}
}
