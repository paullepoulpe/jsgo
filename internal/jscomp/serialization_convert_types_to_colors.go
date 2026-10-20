// Port of com/google/javascript/jscomp/serialization/ConvertTypesToColors.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/base"
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

// ConvertTypesToColors_RemoveTypes is the interface of RemoveTypes and its
// subclass RemoveTypesAndApplyColors.
//
// Java: private static class ConvertTypesToColors.RemoveTypes extends AbstractPostOrderCallback
type ConvertTypesToColors_RemoveTypes interface {
	NodeTraversal_AbstractPostOrderCallback
	GetConvertTypesToColors_RemoveTypesBase() *ConvertTypesToColors_RemoveTypesBase
}

// ConvertTypesToColors_RemoveTypesBase holds RemoveTypes' method bodies.
type ConvertTypesToColors_RemoveTypesBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}

// Java: ConvertTypesToColors.RemoveTypes()
func NewConvertTypesToColors_RemoveTypes() *ConvertTypesToColors_RemoveTypesBase {
	return &ConvertTypesToColors_RemoveTypesBase{}
}

// GetConvertTypesToColors_RemoveTypesBase gives field access through the
// interface.
func (r *ConvertTypesToColors_RemoveTypesBase) GetConvertTypesToColors_RemoveTypesBase() *ConvertTypesToColors_RemoveTypesBase {
	return r
}

// Java: ConvertTypesToColors.RemoveTypes.visit(NodeTraversal, Node, Node)
func (r *ConvertTypesToColors_RemoveTypesBase) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	convertTypesToColors_removeTypesFromNode(n)
}

func convertTypesToColors_removeTypesFromNode(n *rhino.Node) {
	_, _, jsdoc := n.TakeTypesAndJSDocForColorConversion()
	if jsdoc != nil {
		n.SetJSDocInfo(JSDocSerializer_ConvertJSDocInfoForOptimizations(jsdoc))
	}
}

func convertTypesToColors_removeTypesWalk(n *rhino.Node) {
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		convertTypesToColors_removeTypesWalk(c)
	}
	convertTypesToColors_removeTypesFromNode(n)
}

// Java: private static class ConvertTypesToColors.RemoveTypesAndApplyColors extends RemoveTypes
type ConvertTypesToColors_RemoveTypesAndApplyColors struct {
	ConvertTypesToColors_RemoveTypesBase
	colorPoolShard       *ColorPool_ShardView
	typePointersByJstype *base.LinkedIdentityHashMap[rhino.JSType, int]
	colorsByPointer      []*colors.Color
	lastType             rhino.JSType
	lastColor            *colors.Color
	hasLast              bool
}

// Java: ConvertTypesToColors.RemoveTypesAndApplyColors(ColorPool.ShardView, LinkedIdentityHashMap<JSType, Integer>)
func NewConvertTypesToColors_RemoveTypesAndApplyColors(
	colorPoolShard *ColorPool_ShardView,
	typePointersByJstype *base.LinkedIdentityHashMap[rhino.JSType, int]) *ConvertTypesToColors_RemoveTypesAndApplyColors {
	r := &ConvertTypesToColors_RemoveTypesAndApplyColors{}
	r.colorPoolShard = colorPoolShard
	r.typePointersByJstype = typePointersByJstype
	return r
}

func (r *ConvertTypesToColors_RemoveTypesAndApplyColors) initColors() {
	if r.colorPoolShard != nil && r.colorPoolShard.colorPool != nil {
		axiomaticCount := len(TypePointers_OFFSET_TO_AXIOMATIC_COLOR)
		trimmedIds := r.colorPoolShard.trimmedOffsetToId
		idToColor := r.colorPoolShard.colorPool.idToColor
		colorsByPointer := make([]*colors.Color, axiomaticCount+len(trimmedIds))
		copy(colorsByPointer, TypePointers_OFFSET_TO_AXIOMATIC_COLOR)
		for i, id := range trimmedIds {
			colorsByPointer[axiomaticCount+i] = idToColor.Get(id)
		}
		r.colorsByPointer = colorsByPointer
	}
}

// Java: ConvertTypesToColors.RemoveTypesAndApplyColors.visit(NodeTraversal, Node, Node)
func (r *ConvertTypesToColors_RemoveTypesAndApplyColors) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	r.visitNode(n)
}

func (r *ConvertTypesToColors_RemoveTypesAndApplyColors) visitNode(n *rhino.Node) {
	oldType, hadBeforeCast, jsdoc := n.TakeTypesAndJSDocForColorConversion()
	if jsdoc != nil {
		n.SetJSDocInfo(JSDocSerializer_ConvertJSDocInfoForOptimizations(jsdoc))
	}

	if oldType != nil {
		var c *colors.Color
		if r.hasLast && oldType == r.lastType {
			c = r.lastColor
		} else {
			if pointer, ok := r.typePointersByJstype.Lookup(oldType); ok {
				if r.colorsByPointer != nil {
					c = r.colorsByPointer[pointer]
				} else {
					c = r.colorPoolShard.GetColor(pointer)
				}
			}
			r.lastType = oldType
			r.lastColor = c
			r.hasLast = true
		}
		if c != nil {
			n.SetColor(c)
		}
	}

	if hadBeforeCast {
		// used by FunctionInjector and InlineVariables when inlining as a hint that a node has a
		// more specific color
		n.SetColorFromTypeCast()
	}
}

func (r *ConvertTypesToColors_RemoveTypesAndApplyColors) walk(n *rhino.Node) {
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		r.walk(c)
	}
	r.visitNode(n)
}

// Java: ConvertTypesToColors.process(Node, Node)
func (c *ConvertTypesToColors) Process(externs *rhino.Node, root *rhino.Node) {
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
		convertTypesToColors_removeTypesWalk(externsAndJsRoot)
		c.compiler.ClearJSTypeRegistry()
		c.compiler.InitRuntimeLibraryTypedAsts(nil) // Optional.absent()
		return
	}

	stringPoolBuilder := StringPool_Builder_()
	serializeJstypes := SerializeTypesToPointers_Create(c.compiler, stringPoolBuilder, c.serializationOptions)
	serializeJstypes.GatherTypesOnAst(externsAndJsRoot)

	typePool := serializeJstypes.GetTypePool()
	stringPool := stringPoolBuilder.Build()

	colorPoolBuilder := ColorPool_Builder_()
	c.compiler.InitRuntimeLibraryTypedAsts(colorPoolBuilder)

	var callback ConvertTypesToColors_RemoveTypesAndApplyColors
	callback.colorPoolShard = colorPoolBuilder.AddShard(typePool, stringPool)
	callback.typePointersByJstype = serializeJstypes.GetTypePointersByJstype()

	colorPool := colorPoolBuilder.Build()
	callback.initColors()
	callback.walk(externsAndJsRoot)

	c.compiler.ClearJSTypeRegistry()
	c.compiler.SetColorRegistry(colorPool.GetRegistry())
}
