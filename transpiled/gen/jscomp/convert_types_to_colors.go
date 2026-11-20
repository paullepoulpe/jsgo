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
// Port of com/google/javascript/jscomp/serialization/ConvertTypesToColors.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

// Java class: com.google.javascript.jscomp.serialization.ConvertTypesToColors
type ConvertTypesToColors struct {
	compiler             AbstractCompiler      `gson:"compiler"`
	serializationOptions *SerializationOptions `gson:"serializationOptions"`
	idHash_              atomic.Int32
}

func (this_ *ConvertTypesToColors) BindSelf_(s any) {
	_ = s
}
func (this_ *ConvertTypesToColors) IsJavaNil_() bool { return this_ == nil }
func (*ConvertTypesToColors) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.ConvertTypesToColors", "ConvertTypesToColors"
}
func (this_ *ConvertTypesToColors) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConvertTypesToColors) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.ConvertTypesToColors", "com.google.javascript.jscomp.serialization.ConvertTypesToColors", func(v_ any) bool {
		return j2gort.InstanceOf[*ConvertTypesToColors](v_)
	})
}
func (*ConvertTypesToColors) Implements_jscomp_CompilerPass() {}

func (this_ *ConvertTypesToColors) initFields_ConvertTypesToColors() {
	_ = this_
}

// Java constructor: ConvertTypesToColors(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.serialization.SerializationOptions serializationOptions)
func NewConvertTypesToColors(compiler AbstractCompiler, serializationOptions *SerializationOptions) *ConvertTypesToColors {
	this_ := &ConvertTypesToColors{}
	InitConvertTypesToColors(this_, compiler, serializationOptions)
	return this_
}

func InitConvertTypesToColors(this_ *ConvertTypesToColors, compiler AbstractCompiler, serializationOptions *SerializationOptions) {
	_ = compiler
	_ = serializationOptions
	this_.initFields_ConvertTypesToColors()
	this_.compiler = compiler
	this_.serializationOptions = serializationOptions
}

// Java method: ConvertTypesToColors.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ConvertTypesToColors) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer nodestore.ReleaseCurrentStoreJSTypes()

	if this_.compiler.GetLifeCycleStage().HasColorAndSimplifiedJSDoc() {
		return
	}
	if !(!(this_.compiler.GetLifeCycleStage().IsNormalized())) {
		j2gort.CheckStateStr(false, "Not expected to run after normalization")
	}
	var externsAndJsRoot pkg_rhino.Node = root.GetParent()
	_ = externsAndJsRoot
	if !(this_.compiler.HasTypeCheckingRun()) {
		NodeTraversal_Traverse(this_.compiler, externsAndJsRoot, Upcast_NodeTraversal_Callback(NewConvertTypesToColors_RemoveTypes()))
		this_.compiler.ClearJSTypeRegistry()
		this_.compiler.InitRuntimeLibraryTypedAsts(j2gort.Ext_com_google_common_base_Optional_Absent())
		return
	}
	var stringPoolBuilder *StringPool_Builder = StringPool_S_Builder()
	_ = stringPoolBuilder
	var serializeJstypes *SerializeTypesToPointers = SerializeTypesToPointers_Create(this_.compiler, stringPoolBuilder, this_.serializationOptions)
	_ = serializeJstypes
	serializeJstypes.GatherTypesOnAst(externsAndJsRoot)
	var typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool = serializeJstypes.GetTypePool()
	_ = typePool
	var stringPool *StringPool = stringPoolBuilder.Build()
	_ = stringPool
	var colorPoolBuilder *ColorPool_Builder = ColorPool_S_Builder()
	_ = colorPoolBuilder
	this_.compiler.InitRuntimeLibraryTypedAsts(j2gort.Ext_com_google_common_base_Optional_Of(j2gort.PtrToAny(colorPoolBuilder)))
	var callback *ConvertTypesToColors_RemoveTypesAndApplyColors = NewConvertTypesToColors_RemoveTypesAndApplyColors(colorPoolBuilder.AddShard(typePool, stringPool), serializeJstypes)
	_ = callback
	var colorPool *ColorPool = colorPoolBuilder.Build()
	_ = colorPool
	NodeTraversal_Traverse(this_.compiler, externsAndJsRoot, callback)
	this_.compiler.ClearJSTypeRegistry()
	this_.compiler.SetColorRegistry(colorPool.GetRegistry())
}

// Java polymorphic class interface: com.google.javascript.jscomp.serialization.ConvertTypesToColors.RemoveTypes
type ConvertTypesToColors_RemoveTypes interface {
	Implements_jscomp_ConvertTypesToColors_RemoveTypes()
	IsJavaNil_() bool
	GetConvertTypesToColors_RemoveTypesBase() *ConvertTypesToColors_RemoveTypesBase
	NodeTraversal_AbstractPostOrderCallback
}

func ToIface_ConvertTypesToColors_RemoveTypes[T interface {
	ConvertTypesToColors_RemoveTypes
	comparable
}](v T) ConvertTypesToColors_RemoveTypes {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ConvertTypesToColors_RemoveTypes[T ConvertTypesToColors_RemoveTypes](v T) ConvertTypesToColors_RemoveTypes {
	var vi ConvertTypesToColors_RemoveTypes = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.serialization.ConvertTypesToColors.RemoveTypes
type ConvertTypesToColors_RemoveTypesBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}

func (this_ *ConvertTypesToColors_RemoveTypesBase) BindSelf_(s any) {
	_ = s
}
func (this_ *ConvertTypesToColors_RemoveTypesBase) IsJavaNil_() bool { return this_ == nil }
func (*ConvertTypesToColors_RemoveTypesBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.ConvertTypesToColors$RemoveTypes", "RemoveTypes"
}
func (this_ *ConvertTypesToColors_RemoveTypesBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.ConvertTypesToColors.RemoveTypes", "com.google.javascript.jscomp.serialization.ConvertTypesToColors$RemoveTypes", func(v_ any) bool {
		return j2gort.InstanceOf[ConvertTypesToColors_RemoveTypes](v_)
	})
}
func (b *ConvertTypesToColors_RemoveTypesBase) GetConvertTypesToColors_RemoveTypesBase() *ConvertTypesToColors_RemoveTypesBase {
	return b
}
func (*ConvertTypesToColors_RemoveTypesBase) Implements_jscomp_ConvertTypesToColors_RemoveTypes() {}
func (*ConvertTypesToColors_RemoveTypesBase) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ConvertTypesToColors_RemoveTypesBase) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ConvertTypesToColors_RemoveTypesBase) initFields_ConvertTypesToColors_RemoveTypes() {
	_ = this_
}

// Java constructor: RemoveTypes()
func NewConvertTypesToColors_RemoveTypes() ConvertTypesToColors_RemoveTypes {
	this_ := &ConvertTypesToColors_RemoveTypesBase{}
	InitConvertTypesToColors_RemoveTypes(this_)
	return this_
}

func InitConvertTypesToColors_RemoveTypes(this_ *ConvertTypesToColors_RemoveTypesBase) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ConvertTypesToColors_RemoveTypes()
}

// Java method: RemoveTypes.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ConvertTypesToColors_RemoveTypesBase) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	n.SetJSType(nil)
	n.SetJSTypeBeforeCast(nil)
	n.SetDeclaredTypeExpression(pkg_rhino.Node(0))
	n.SetTypedefTypeProp(nil)
	var jsdoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = jsdoc
	if jsdoc != nil {
		n.SetJSDocInfo(JSDocSerializer_ConvertJSDocInfoForOptimizations(jsdoc))
	}
}

// Java class: com.google.javascript.jscomp.serialization.ConvertTypesToColors.RemoveTypesAndApplyColors
type ConvertTypesToColors_RemoveTypesAndApplyColors struct {
	ConvertTypesToColors_RemoveTypesBase
	colorPoolShard       *ColorPool_ShardView      `gson:"colorPoolShard"`
	typePointersByJstype *SerializeTypesToPointers `gson:"typePointersByJstype"`
}

func (this_ *ConvertTypesToColors_RemoveTypesAndApplyColors) BindSelf_(s any) {
	_ = s
}
func (this_ *ConvertTypesToColors_RemoveTypesAndApplyColors) IsJavaNil_() bool { return this_ == nil }
func (*ConvertTypesToColors_RemoveTypesAndApplyColors) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.ConvertTypesToColors$RemoveTypesAndApplyColors", "RemoveTypesAndApplyColors"
}
func (this_ *ConvertTypesToColors_RemoveTypesAndApplyColors) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ConvertTypesToColors_RemoveTypesBase.IdentityHash_()
}
func (*ConvertTypesToColors_RemoveTypesAndApplyColors) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.ConvertTypesToColors.RemoveTypesAndApplyColors", "com.google.javascript.jscomp.serialization.ConvertTypesToColors$RemoveTypesAndApplyColors", func(v_ any) bool {
		return j2gort.InstanceOf[*ConvertTypesToColors_RemoveTypesAndApplyColors](v_)
	})
}
func (*ConvertTypesToColors_RemoveTypesAndApplyColors) Implements_jscomp_ConvertTypesToColors_RemoveTypes() {
}
func (*ConvertTypesToColors_RemoveTypesAndApplyColors) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ConvertTypesToColors_RemoveTypesAndApplyColors) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ConvertTypesToColors_RemoveTypesAndApplyColors) initFields_ConvertTypesToColors_RemoveTypesAndApplyColors() {
	_ = this_
}

// Java constructor: RemoveTypesAndApplyColors(com.google.javascript.jscomp.serialization.ColorPool.ShardView colorPoolShard, com.google.javascript.jscomp.serialization.SerializeTypesToPointers typePointersByJstype)
func NewConvertTypesToColors_RemoveTypesAndApplyColors(colorPoolShard *ColorPool_ShardView, typePointersByJstype *SerializeTypesToPointers) *ConvertTypesToColors_RemoveTypesAndApplyColors {
	this_ := &ConvertTypesToColors_RemoveTypesAndApplyColors{}
	InitConvertTypesToColors_RemoveTypesAndApplyColors(this_, colorPoolShard, typePointersByJstype)
	return this_
}

func InitConvertTypesToColors_RemoveTypesAndApplyColors(this_ *ConvertTypesToColors_RemoveTypesAndApplyColors, colorPoolShard *ColorPool_ShardView, typePointersByJstype *SerializeTypesToPointers) {
	_ = colorPoolShard
	_ = typePointersByJstype
	InitConvertTypesToColors_RemoveTypes(&this_.ConvertTypesToColors_RemoveTypesBase)
	this_.initFields_ConvertTypesToColors_RemoveTypesAndApplyColors()
	this_.colorPoolShard = colorPoolShard
	this_.typePointersByJstype = typePointersByJstype
}

// Java method: RemoveTypesAndApplyColors.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ConvertTypesToColors_RemoveTypesAndApplyColors) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var oldType pkg_rhino.JSType = n.GetJSType()
	_ = oldType
	var oldTypeBeforeCast pkg_rhino.JSType = n.GetJSTypeBeforeCast()
	_ = oldTypeBeforeCast
	this_.ConvertTypesToColors_RemoveTypesBase.Visit_3(t, n, parent)
	if oldType != nil && !oldType.IsJavaNil_() {
		var pointer int32 = this_.typePointersByJstype.GetTypePointer(oldType)
		_ = pointer
		if pointer >= int32(0) {
			n.SetColor(this_.colorPoolShard.GetColor(pointer))
		}
	}
	if oldTypeBeforeCast != nil && !oldTypeBeforeCast.IsJavaNil_() {
		n.SetColorFromTypeCast()
	}
}
func (*ConvertTypesToColors_RemoveTypesBase) CanReleaseSyntacticScopes_() {
}
func (*ConvertTypesToColors_RemoveTypesAndApplyColors) CanReleaseSyntacticScopes_() {
}
func (*ConvertTypesToColors_RemoveTypesAndApplyColors) AlwaysShouldTraverse_() {}
