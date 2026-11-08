// Port of com/google/javascript/jscomp/serialization/ConvertTypesToColors.java @ closure-compiler v20261104

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

func (r *convertTypesToColors_removeTypes) visitNode(n rhino.Node) {
	_, _, jsdoc := n.TakeTypesAndJSDocForColorConversion()
	if jsdoc != nil {
		n.SetJSDocInfo(jsDocSerializer_convertJSDocInfoWithBuilder(jsdoc, &r.jsdocBuilder))
	}
}

func (r *convertTypesToColors_removeTypes) walk(n rhino.Node) {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		r.walk(c)
	}
	r.visitNode(n)
}

// Java: private static class ConvertTypesToColors.RemoveTypesAndApplyColors extends RemoveTypes
type ConvertTypesToColors_RemoveTypesAndApplyColors struct {
	offsetToColor         []*colors.Color
	jstypeReconserializer *JSTypeReconserializer
	jsdocBuilder          rhino.JSDocInfo_Builder
	// detached: walking a subtree no longer attached to the AST
	// (convertTypesToColors_forDetachedRoots, Go-only).
	detached bool
}

func (r *ConvertTypesToColors_RemoveTypesAndApplyColors) visitNode(n rhino.Node) {
	oldType, hadBeforeCast, jsdoc := n.TakeTypesAndJSDocForColorConversion()
	if jsdoc != nil {
		n.SetJSDocInfo(jsDocSerializer_convertJSDocInfoWithBuilder(jsdoc, &r.jsdocBuilder))
	}

	appliedColor := false
	if oldType != nil && r.jstypeReconserializer != nil {
		if pointer, ok := r.jstypeReconserializer.typePointers[oldType]; ok {
			if c := r.offsetToColor[pointer]; c != nil {
				n.SetColor(c)
				appliedColor = true
			}
		}
	}

	// Go-only: a detached node's types were never gathered (Java does not
	// visit detached nodes at all), so it gets no color and must not be
	// marked as having its color from a cast: Node.setColorFromTypeCast
	// requires a color.
	if hadBeforeCast && (!r.detached || appliedColor) {
		// used by FunctionInjector and InlineVariables when inlining as a hint that a node has a
		// more specific color
		n.SetColorFromTypeCast()
	}
}

// walkDetached is walk for a subtree that is no longer attached to the AST.
func (r *ConvertTypesToColors_RemoveTypesAndApplyColors) walkDetached(n rhino.Node) {
	r.detached = true
	r.walk(n)
	r.detached = false
}

func (r *ConvertTypesToColors_RemoveTypesAndApplyColors) walk(n rhino.Node) {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		r.walk(c)
	}
	r.visitNode(n)
}

func collectScriptsForColorPass(externs, root, externsAndJsRoot rhino.Node, skipTypeSummaries bool) ([]rhino.Node, bool) {
	if externsAndJsRoot.IsNil() || !externsAndJsRoot.IsRoot() || !externsAndJsRoot.HasTwoChildren() {
		return nil, false
	}
	if externsAndJsRoot.GetFirstChild() != externs || externsAndJsRoot.GetSecondChild() != root {
		return nil, false
	}
	extCount, ok1 := countRootScripts(externs)
	rootCount, ok2 := countRootScripts(root)
	if !ok1 || !ok2 {
		return nil, false
	}
	total := extCount + rootCount
	if total == 0 {
		return nil, true
	}
	scripts := make([]rhino.Node, 0, total)
	for _, parent := range [2]rhino.Node{externs, root} {
		for s := parent.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
			if skipTypeSummaries && NodeUtil_IsFromTypeSummary(s) {
				continue
			}
			scripts = append(scripts, s)
		}
	}
	return scripts, true
}

type colorTypeGatherWorker struct {
	workerIdx  int
	seen       map[*rhino.JSTypeBase]uint32
	gen        uint32
	types      []*rhino.JSTypeBase
	props      *jutil.BitSet
	lastPropID uint32
	baseSeen   map[*rhino.JSTypeBase]int32
}

func (w *colorTypeGatherWorker) recordType(tb *rhino.JSTypeBase) {
	if _, axiomatic := w.baseSeen[tb]; axiomatic {
		return
	}
	if w.seen[tb] == w.gen {
		return
	}
	if w.seen == nil {
		w.seen = make(map[*rhino.JSTypeBase]uint32, 256)
	}
	w.seen[tb] = w.gen
	w.types = append(w.types, tb)
}

func (w *colorTypeGatherWorker) walkPropsAndTypes(n rhino.Node) {
	if n.IsScript() && NodeUtil_IsFromTypeSummary(n) {
		return
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		w.walkPropsAndTypes(c)
	}
	switch n.GetToken() {
	case rhino.Token_GETPROP,
		rhino.Token_OPTCHAIN_GETPROP:
		propID := n.GetStringID()
		if propID != w.lastPropID {
			w.lastPropID = propID
			w.props.Set(int(propID))
		}
	case rhino.Token_STRING_KEY,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF:
		if !n.IsQuotedStringKey() {
			propID := n.GetStringID()
			if propID != w.lastPropID {
				w.lastPropID = propID
				w.props.Set(int(propID))
			}
		}
	case rhino.Token_ROOT:
		return
	case rhino.Token_NAME:
		if tb := n.GetJSTypeBase(); tb != nil {
			w.recordType(tb)
		}
		if shadow := n.GetClosureUnawareShadow(); !shadow.IsNil() {
			w.walkShadowTypes(shadow.GetFirstFirstChild().GetFirstChild())
		}
		return
	}
	if tb := n.GetJSTypeBase(); tb != nil {
		w.recordType(tb)
	}
}

func (w *colorTypeGatherWorker) walkShadowTypes(n rhino.Node) {
	if n.IsScript() && NodeUtil_IsFromTypeSummary(n) {
		return
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		w.walkShadowTypes(c)
	}
	if n.IsRoot() {
		return
	}
	if tb := n.GetJSTypeBase(); tb != nil {
		w.recordType(tb)
	}
	if n.IsName() {
		if shadow := n.GetClosureUnawareShadow(); !shadow.IsNil() {
			w.walkShadowTypes(shadow.GetFirstFirstChild().GetFirstChild())
		}
	}
}

func convertTypesToColors_gatherTypesParallel(
	compiler AbstractCompiler,
	serializeJstypes *SerializeTypesToPointers,
	externs, root, externsAndJsRoot rhino.Node,
) bool {
	if serializeJstypes.gathered {
		return false
	}
	scripts, ok := collectScriptsForColorPass(externs, root, externsAndJsRoot, true)
	if !ok {
		return false
	}
	numWorkers := parallelThreadCount(compiler, len(scripts))
	if numWorkers <= 1 {
		return false
	}

	recon := serializeJstypes.jstypeReconserializer
	baseSeen := recon.typePointers
	perScriptTypes := make([][]*rhino.JSTypeBase, len(scripts))
	workerProps := make([]jutil.BitSet, numWorkers)

	func() {
		if compiler != nil {
			prewarmCompilerForParallelReads(compiler)
			restore := enableParallelTypeRegistryReads(compiler)
			defer restore()
		}
		runParallelWorkersIndexed(
			len(scripts),
			numWorkers,
			func(workerIdx int) colorTypeGatherWorker {
				return colorTypeGatherWorker{
					workerIdx: workerIdx,
					props:     &workerProps[workerIdx],
					baseSeen:  baseSeen,
				}
			},
			func(w *colorTypeGatherWorker, idx int) {
				w.gen = uint32(idx + 1)
				w.types = nil
				w.walkPropsAndTypes(scripts[idx])
				perScriptTypes[idx] = w.types
				w.types = nil
			},
			nil,
		)
	}()

	serializeJstypes.gathered = true
	for i := range workerProps {
		serializeJstypes.propertiesReferencedInAst.Or(&workerProps[i])
	}
	for _, scriptTypes := range perScriptTypes {
		for _, tb := range scriptTypes {
			if !recon.hasTypePointer(tb) {
				recon.SerializeType(tb.Self())
			}
		}
	}
	for _, mismatch := range serializeJstypes.compiler.GetTypeMismatches() {
		recon.SerializeType(mismatch.Found)
		recon.SerializeType(mismatch.Required)
	}
	recon.prepareForColorPool()
	serializeJstypes.logSerializationDebugInfo(recon)
	return true
}

func convertTypesToColors_removeTypesAst(
	compiler AbstractCompiler,
	externs, root, externsAndJsRoot rhino.Node,
	remover *convertTypesToColors_removeTypes,
) {
	if scripts, ok := collectScriptsForColorPass(externs, root, externsAndJsRoot, false); ok {
		if numWorkers := parallelThreadCount(compiler, len(scripts)); numWorkers > 1 {
			runParallelWorkersIndexed(
				len(scripts),
				numWorkers,
				func(_ int) convertTypesToColors_removeTypes {
					return convertTypesToColors_removeTypes{}
				},
				func(w *convertTypesToColors_removeTypes, idx int) {
					w.walk(scripts[idx])
				},
				nil,
			)
			remover.visitNode(externs)
			remover.visitNode(root)
			remover.visitNode(externsAndJsRoot)
			return
		}
	}
	remover.walk(externsAndJsRoot)
}

func convertTypesToColors_applyColorsAst(
	compiler AbstractCompiler,
	externs, root, externsAndJsRoot rhino.Node,
	callback *ConvertTypesToColors_RemoveTypesAndApplyColors,
) {
	if scripts, ok := collectScriptsForColorPass(externs, root, externsAndJsRoot, false); ok {
		if numWorkers := parallelThreadCount(compiler, len(scripts)); numWorkers > 1 {
			runParallelWorkersIndexed(
				len(scripts),
				numWorkers,
				func(_ int) ConvertTypesToColors_RemoveTypesAndApplyColors {
					return ConvertTypesToColors_RemoveTypesAndApplyColors{
						offsetToColor:         callback.offsetToColor,
						jstypeReconserializer: callback.jstypeReconserializer,
					}
				},
				func(w *ConvertTypesToColors_RemoveTypesAndApplyColors, idx int) {
					w.walk(scripts[idx])
				},
				nil,
			)
			callback.visitNode(externs)
			callback.visitNode(root)
			callback.visitNode(externsAndJsRoot)
			return
		}
	}
	callback.walk(externsAndJsRoot)
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
		convertTypesToColors_removeTypesAst(c.compiler, externs, root, externsAndJsRoot, &remover)
		convertTypesToColors_forDetachedRoots(c.compiler, externsAndJsRoot, remover.walk)
		c.compiler.ClearJSTypeRegistry()
		c.compiler.InitRuntimeLibraryTypedAsts(nil) // Optional.absent()
		return
	}

	serializeJstypes := SerializeTypesToPointers_Create(c.compiler, nil, c.serializationOptions)
	if !convertTypesToColors_gatherTypesParallel(c.compiler, serializeJstypes, externs, root, externsAndJsRoot) {
		serializeJstypes.gatherTypesForColors(externsAndJsRoot)
	}

	colorPoolBuilder := NewColorPoolBuilder()
	c.compiler.InitRuntimeLibraryTypedAsts(colorPoolBuilder)
	colorPoolBuilder.addJSTypeShard(serializeJstypes.jstypeReconserializer)

	colorPool := colorPoolBuilder.Build()

	var callback ConvertTypesToColors_RemoveTypesAndApplyColors
	callback.jstypeReconserializer = serializeJstypes.jstypeReconserializer
	callback.offsetToColor = serializeJstypes.jstypeReconserializer.offsetToColor

	convertTypesToColors_applyColorsAst(c.compiler, externs, root, externsAndJsRoot, &callback)
	convertTypesToColors_forDetachedRoots(c.compiler, externsAndJsRoot, callback.walkDetached)

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
		visitBindings := func(bindings *OrderedStringMap[*Binding]) {
			for _, entry := range bindings.Entries() {
				b := entry.Value
				visit(b.SourceNode())
				visitExport(b.OriginatingExport())
			}
		}
		visitModule := func(m *Module) {
			visit(m.Metadata().RootNode())
			visitBindings(m.NamespaceMap())
			visitBindings(m.BoundNamesMap())
			for _, entry := range m.LocalExportsMap().Entries() {
				visitExport(entry.Value)
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
