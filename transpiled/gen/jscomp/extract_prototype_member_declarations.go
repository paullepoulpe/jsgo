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
// Port of com/google/javascript/jscomp/ExtractPrototypeMemberDeclarations.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_extract_prototype_member_declarations_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_extract_prototype_member_declarations_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_12 any) any {
		k := j2gort.CastPtr[JSChunk](lp_12)
		_ = k
		return NewExtractPrototypeMemberDeclarations_ExtractionInstanceInfo()
	})
}

// Java class: com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations
type ExtractPrototypeMemberDeclarations struct {
	compiler AbstractCompiler                            `gson:"compiler"`
	pattern  *ExtractPrototypeMemberDeclarations_Pattern `gson:"pattern"`
	idHash_  atomic.Int32
}

func (this_ *ExtractPrototypeMemberDeclarations) BindSelf_(s any) {
	_ = s
}
func (this_ *ExtractPrototypeMemberDeclarations) IsJavaNil_() bool { return this_ == nil }
func (*ExtractPrototypeMemberDeclarations) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations", "ExtractPrototypeMemberDeclarations"
}
func (this_ *ExtractPrototypeMemberDeclarations) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ExtractPrototypeMemberDeclarations) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations", "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations", func(v_ any) bool {
		return j2gort.InstanceOf[*ExtractPrototypeMemberDeclarations](v_)
	})
}
func (*ExtractPrototypeMemberDeclarations) Implements_jscomp_CompilerPass() {}

func StaticGet_ExtractPrototypeMemberDeclarations_PROTOTYPE_ALIAS() string {
	return "JSCompiler_prototypeAlias"
}
func (this_ *ExtractPrototypeMemberDeclarations) initFields_ExtractPrototypeMemberDeclarations() {
	_ = this_
}

// Java constructor: ExtractPrototypeMemberDeclarations(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.Pattern pattern)
func NewExtractPrototypeMemberDeclarations(compiler AbstractCompiler, pattern *ExtractPrototypeMemberDeclarations_Pattern) *ExtractPrototypeMemberDeclarations {
	this_ := &ExtractPrototypeMemberDeclarations{}
	InitExtractPrototypeMemberDeclarations(this_, compiler, pattern)
	return this_
}

func InitExtractPrototypeMemberDeclarations(this_ *ExtractPrototypeMemberDeclarations, compiler AbstractCompiler, pattern *ExtractPrototypeMemberDeclarations_Pattern) {
	_ = compiler
	_ = pattern
	this_.initFields_ExtractPrototypeMemberDeclarations()
	this_.compiler = compiler
	this_.pattern = pattern
}

// Java method: ExtractPrototypeMemberDeclarations.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ExtractPrototypeMemberDeclarations) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var extractionInfo *ExtractPrototypeMemberDeclarations_GatherExtractionInfo = NewExtractPrototypeMemberDeclarations_GatherExtractionInfo(this_)
	_ = extractionInfo
	NodeTraversal_Traverse(this_.compiler, root, extractionInfo)
	this_.priv_ExtractPrototypeMemberDeclarations_maybeDoExtraction_2860(extractionInfo)
}

// Java method: ExtractPrototypeMemberDeclarations.maybeDoExtraction(com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.GatherExtractionInfo info)
func (this_ *ExtractPrototypeMemberDeclarations) priv_ExtractPrototypeMemberDeclarations_maybeDoExtraction_2860(info *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) {
	_ = info
	if ((this_.pattern == StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE()) || (this_.pattern == StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP())) && !(info.priv_GatherExtractionInfo_shouldExtractGlobal_2863()) {
		return
	}
	if this_.pattern == StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP() {
		var injectionPoint pkg_rhino.Node = this_.compiler.GetNodeForCodeInsertion(nil)
		_ = injectionPoint
		var var_ pkg_rhino.Node = NodeUtil_NewVarNode_2_string_Node("JSCompiler_prototypeAlias", pkg_rhino.Node(0)).SrcrefTreeIfMissing(injectionPoint)
		_ = var_
		injectionPoint.AddChildToFront(var_)
		this_.compiler.ReportChangeToEnclosingScope(var_)
	}
	{
		cur_1 := info.instancesByChunk.IterEntries()
		for cur_1.NextEntry() {
			var alias string = "JSCompiler_prototypeAlias"
			_ = alias
			if this_.pattern == StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP() {
				if info.priv_GatherExtractionInfo_shouldExtractChunk_2864(cur_1.Key()) {
					var injectionPoint_2 pkg_rhino.Node = this_.compiler.GetNodeForCodeInsertion(cur_1.Key())
					_ = injectionPoint_2
					alias = ("JSCompiler_prototypeAlias" + j2gort.StrValInt32(cur_1.Key().GetIndex()))
					var var__2 pkg_rhino.Node = NodeUtil_NewVarNode_2_string_Node(alias, pkg_rhino.Node(0)).SrcrefTreeIfMissing(injectionPoint_2)
					_ = var__2
					injectionPoint_2.AddChildToFront(var__2)
					this_.compiler.ReportChangeToEnclosingScope(var__2)
				} else {
					continue
				}
			}
			for _, instance := range cur_1.Value().instances {
				_ = instance
				this_.priv_ExtractPrototypeMemberDeclarations_extractInstance_2861(instance, alias)
			}

		}
	}
}

// Java method: ExtractPrototypeMemberDeclarations.extractInstance(com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.ExtractionInstance instance, java.lang.String alias)
func (this_ *ExtractPrototypeMemberDeclarations) priv_ExtractPrototypeMemberDeclarations_extractInstance_2861(instance *ExtractPrototypeMemberDeclarations_ExtractionInstance, alias string) {
	_ = instance
	_ = alias
	var first *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration = instance.declarations[int32(0)]
	_ = first
	var className string = first.qualifiedClassName
	_ = className
	if (this_.pattern == StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP()) || (this_.pattern == StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP()) {
		var classNameNode pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, className)
		_ = classNameNode
		classNameNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), first.constant)
		var stmt pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(pkg_rhino.IR_Name(alias), pkg_rhino.IR_Getprop_2(classNameNode, "prototype"))).SrcrefTreeIfMissing(first.node)
		_ = stmt
		stmt.InsertBefore(first.node)
		this_.compiler.ReportChangeToEnclosingScope(stmt)
	} else if this_.pattern == StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE() {
		var block pkg_rhino.Node = pkg_rhino.IR_Block_0()
		_ = block
		var func_ pkg_rhino.Node = pkg_rhino.IR_Function(pkg_rhino.IR_Name(""), pkg_rhino.IR_ParamList([]pkg_rhino.Node{pkg_rhino.IR_Name(alias)}), block)
		_ = func_
		var call pkg_rhino.Node = pkg_rhino.IR_Call(func_, []pkg_rhino.Node{NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.compiler, (j2gort.StrValString(className) + ".prototype"), instance.parent, (j2gort.StrValString(className) + ".prototype"))})
		_ = call
		call.PutIntProp(pkg_rhino.StaticGet_Node_FREE_CALL(), int32(1))
		var stmt_2 pkg_rhino.Node = pkg_rhino.IR_ExprResult(call)
		_ = stmt_2
		stmt_2.SrcrefTreeIfMissing(first.node)
		stmt_2.InsertBefore(first.node)
		this_.compiler.ReportChangeToEnclosingScope(stmt_2)
		for _, declar := range instance.declarations {
			_ = declar
			this_.compiler.ReportChangeToEnclosingScope(declar.node)
			block.AddChildToBack(declar.node.Detach())
		}

	}
	for _, declar_2 := range instance.declarations {
		_ = declar_2
		this_.priv_ExtractPrototypeMemberDeclarations_replacePrototypeMemberDeclaration_2862(declar_2, alias)
	}

}

// Java method: ExtractPrototypeMemberDeclarations.replacePrototypeMemberDeclaration(com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.PrototypeMemberDeclaration declar, java.lang.String alias)
func (this_ *ExtractPrototypeMemberDeclarations) priv_ExtractPrototypeMemberDeclarations_replacePrototypeMemberDeclaration_2862(declar *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration, alias string) {
	_ = declar
	_ = alias
	var assignment pkg_rhino.Node = declar.node.GetFirstChild()
	_ = assignment
	var lhs pkg_rhino.Node = assignment.GetFirstChild()
	_ = lhs
	var name pkg_rhino.Node = NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.compiler, (j2gort.StrValString(alias) + "." + j2gort.StrValString(declar.memberName)), declar.node, declar.memberName)
	_ = name
	var accessNode pkg_rhino.Node = declar.lhs.GetFirstFirstChild()
	_ = accessNode
	var originalName string = accessNode.GetOriginalName()
	_ = originalName
	var className string
	_ = className
	if originalName != j2gort.NullString {
		className = originalName
	} else {
		className = "?"
	}
	name.GetFirstChild().SrcrefTree(lhs)
	name.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), lhs.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME()))
	name.GetFirstChild().SetOriginalName((j2gort.StrValString(className) + ".prototype"))
	lhs.ReplaceWith(name)
	this_.compiler.ReportChangeToEnclosingScope(name)
}

// Java enum: com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.Pattern
type ExtractPrototypeMemberDeclarations_Pattern struct {
	j2gort.EnumBase
	globalOverhead        int32
	perExtractionOverhead int32
	perMemberOverhead     int32
}

func (this_ *ExtractPrototypeMemberDeclarations_Pattern) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *ExtractPrototypeMemberDeclarations_Pattern) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*ExtractPrototypeMemberDeclarations_Pattern) Implements_j2gort_Ext_java_lang_Enum() {}
func (*ExtractPrototypeMemberDeclarations_Pattern) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ *ExtractPrototypeMemberDeclarations_Pattern) initFields_ExtractPrototypeMemberDeclarations_Pattern() {
	_ = this_
}

func initEnumExtractPrototypeMemberDeclarations_Pattern(this_ *ExtractPrototypeMemberDeclarations_Pattern, enumName_ string, enumOrd_ int32, globalOverHead int32, perExtractionOverhead int32, perMemberOverhead int32) *ExtractPrototypeMemberDeclarations_Pattern {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = globalOverHead
	_ = perExtractionOverhead
	_ = perMemberOverhead
	this_.initFields_ExtractPrototypeMemberDeclarations_Pattern()
	this_.globalOverhead = globalOverHead
	this_.perExtractionOverhead = perExtractionOverhead
	this_.perMemberOverhead = perMemberOverhead
	return this_
}

var enum_ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP *ExtractPrototypeMemberDeclarations_Pattern

func StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP() *ExtractPrototypeMemberDeclarations_Pattern {
	clinit_ExtractPrototypeMemberDeclarations_Pattern()
	return enum_ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP
}

var enum_ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP *ExtractPrototypeMemberDeclarations_Pattern

func StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP() *ExtractPrototypeMemberDeclarations_Pattern {
	clinit_ExtractPrototypeMemberDeclarations_Pattern()
	return enum_ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP
}

var enum_ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE *ExtractPrototypeMemberDeclarations_Pattern

func StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE() *ExtractPrototypeMemberDeclarations_Pattern {
	clinit_ExtractPrototypeMemberDeclarations_Pattern()
	return enum_ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE
}

var clinitDone_ExtractPrototypeMemberDeclarations_Pattern uint32
var clinitState_ExtractPrototypeMemberDeclarations_Pattern j2gort.ClinitState

func clinit_ExtractPrototypeMemberDeclarations_Pattern() {
	if atomic.LoadUint32(&clinitDone_ExtractPrototypeMemberDeclarations_Pattern) != 0 {
		return
	}
	clinitSlow_ExtractPrototypeMemberDeclarations_Pattern()
}

func clinitSlow_ExtractPrototypeMemberDeclarations_Pattern() {
	if !j2gort.EnterClinit(&clinitState_ExtractPrototypeMemberDeclarations_Pattern, &clinitDone_ExtractPrototypeMemberDeclarations_Pattern) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ExtractPrototypeMemberDeclarations_Pattern, &clinitDone_ExtractPrototypeMemberDeclarations_Pattern)
	enum_ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP = &ExtractPrototypeMemberDeclarations_Pattern{}
	initEnumExtractPrototypeMemberDeclarations_Pattern(enum_ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP, "USE_GLOBAL_TEMP", int32(0), j2gort.JStrLen("var t;"), j2gort.JStrLen("t=y.prototype;"), (j2gort.JStrLen("t.y=") - j2gort.JStrLen("x[p].y=")))
	enum_ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP = &ExtractPrototypeMemberDeclarations_Pattern{}
	initEnumExtractPrototypeMemberDeclarations_Pattern(enum_ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP, "USE_CHUNK_TEMP", int32(1), j2gort.JStrLen("var t;"), j2gort.JStrLen("t=y.prototype;"), (j2gort.JStrLen("t.y=") - j2gort.JStrLen("x[p].y=")))
	enum_ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE = &ExtractPrototypeMemberDeclarations_Pattern{}
	initEnumExtractPrototypeMemberDeclarations_Pattern(enum_ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE, "USE_IIFE", int32(2), int32(0), j2gort.JStrLen("(function(t){})(y.prototype);"), (j2gort.JStrLen("t.y=") - j2gort.JStrLen("x.prototype.y=")))
	enumValues_ExtractPrototypeMemberDeclarations_Pattern = []*ExtractPrototypeMemberDeclarations_Pattern{enum_ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP, enum_ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP, enum_ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.Pattern", []any{enum_ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP, enum_ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP, enum_ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations$Pattern", []any{enum_ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP, enum_ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP, enum_ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.Pattern", "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations$Pattern", func() {
		clinit_ExtractPrototypeMemberDeclarations_Pattern()
	})
}

var enumValues_ExtractPrototypeMemberDeclarations_Pattern []*ExtractPrototypeMemberDeclarations_Pattern

func ExtractPrototypeMemberDeclarations_Pattern_Values() []*ExtractPrototypeMemberDeclarations_Pattern {
	clinit_ExtractPrototypeMemberDeclarations_Pattern()
	return j2gort.CloneSlice(enumValues_ExtractPrototypeMemberDeclarations_Pattern)
}

func ExtractPrototypeMemberDeclarations_Pattern_ValueOf(name string) *ExtractPrototypeMemberDeclarations_Pattern {
	clinit_ExtractPrototypeMemberDeclarations_Pattern()
	for _, c := range enumValues_ExtractPrototypeMemberDeclarations_Pattern {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.Pattern." + name))
}

func (this_ *ExtractPrototypeMemberDeclarations_Pattern) IsJavaNil_() bool { return this_ == nil }
func (*ExtractPrototypeMemberDeclarations_Pattern) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations$Pattern", "Pattern"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.Pattern", "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations$Pattern", func(v_ any) bool {
		return j2gort.InstanceOf[*ExtractPrototypeMemberDeclarations_Pattern](v_)
	})
}

// Java class: com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.ExtractionInstanceInfo
type ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo struct {
	instances  []*ExtractPrototypeMemberDeclarations_ExtractionInstance `gson:"instances"`
	totalDelta int32                                                    `gson:"totalDelta"`
	idHash_    atomic.Int32
}

func (this_ *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo) IsJavaNil_() bool {
	return this_ == nil
}
func (*ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations$ExtractionInstanceInfo", "ExtractionInstanceInfo"
}
func (this_ *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.ExtractionInstanceInfo", "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations$ExtractionInstanceInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo](v_)
	})
}

func (this_ *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo) initFields_ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo() {
	_ = this_
	this_.instances = []*ExtractPrototypeMemberDeclarations_ExtractionInstance{

		// Java constructor: ExtractionInstanceInfo()
	}

}

func NewExtractPrototypeMemberDeclarations_ExtractionInstanceInfo() *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo {
	this_ := &ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo{}
	InitExtractPrototypeMemberDeclarations_ExtractionInstanceInfo(this_)
	return this_
}

func InitExtractPrototypeMemberDeclarations_ExtractionInstanceInfo(this_ *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo) {
	this_.initFields_ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo()
}

// Java class: com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.GatherExtractionInfo
type ExtractPrototypeMemberDeclarations_GatherExtractionInfo struct {
	NodeTraversal_AbstractShallowCallbackBase
	outer_           *ExtractPrototypeMemberDeclarations
	instancesByChunk j2gort.OrderedMap[*JSChunk, *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo] `gson:"instancesByChunk"`
}

func (this_ *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) IsJavaNil_() bool {
	return this_ == nil
}
func (*ExtractPrototypeMemberDeclarations_GatherExtractionInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations$GatherExtractionInfo", "GatherExtractionInfo"
}
func (this_ *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractShallowCallbackBase.IdentityHash_()
}
func (*ExtractPrototypeMemberDeclarations_GatherExtractionInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.GatherExtractionInfo", "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations$GatherExtractionInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*ExtractPrototypeMemberDeclarations_GatherExtractionInfo](v_)
	})
}
func (this_ *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) SetOuter_(o *ExtractPrototypeMemberDeclarations) {
	this_.outer_ = o
}
func (this_ *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) GetOuter_() *ExtractPrototypeMemberDeclarations {
	return this_.outer_
}
func (*ExtractPrototypeMemberDeclarations_GatherExtractionInfo) Implements_jscomp_NodeTraversal_AbstractShallowCallback() {
}
func (*ExtractPrototypeMemberDeclarations_GatherExtractionInfo) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) initFields_ExtractPrototypeMemberDeclarations_GatherExtractionInfo() {
	_ = this_
	this_.instancesByChunk = j2gort.NewOrderedMap[*JSChunk, *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo]()
}

// Java constructor: GatherExtractionInfo()
func NewExtractPrototypeMemberDeclarations_GatherExtractionInfo(outer_ *ExtractPrototypeMemberDeclarations) *ExtractPrototypeMemberDeclarations_GatherExtractionInfo {
	this_ := &ExtractPrototypeMemberDeclarations_GatherExtractionInfo{}
	this_.outer_ = outer_
	InitExtractPrototypeMemberDeclarations_GatherExtractionInfo(this_)
	return this_
}

func InitExtractPrototypeMemberDeclarations_GatherExtractionInfo(this_ *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) {
	InitNodeTraversal_AbstractShallowCallback(&this_.NodeTraversal_AbstractShallowCallbackBase)
	this_.initFields_ExtractPrototypeMemberDeclarations_GatherExtractionInfo()
}

// Java method: GatherExtractionInfo.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsScript()) && !(n.IsBlock()) {
		return
	}
	{
		var cur pkg_rhino.Node = n.GetFirstChild()
		_ = cur
		for ; cur != 0; cur = cur.GetNext() {
			var prototypeMember *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration = ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_ExtractDeclaration(cur)
			_ = prototypeMember
			if prototypeMember == nil {
				continue
			}
			var instance *ExtractPrototypeMemberDeclarations_ExtractionInstance = NewExtractPrototypeMemberDeclarations_ExtractionInstance(this_.outer_, prototypeMember, n)
			_ = instance
			cur = j2gort.SliceGetLast(instance.declarations).node
			if instance.IsFavorable() {
				j2gort.OrderedMapComputeIfAbsentFn(&this_.instancesByChunk, t.GetChunk(), nil, lambdaSingleton_extract_prototype_member_declarations_1)

				var instanceInfo *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo = this_.instancesByChunk.Get(t.GetChunk())
				_ = instanceInfo
				instanceInfo.instances = append(instanceInfo.instances, instance)
				r_13 := instanceInfo
				r_13.totalDelta = (r_13.totalDelta + instance.delta)
			}
		}
	}
}

// Java method: GatherExtractionInfo.shouldExtractGlobal()
func (this_ *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) priv_GatherExtractionInfo_shouldExtractGlobal_2863() bool {
	var allModulesDelta int32 = int32(0)
	_ = allModulesDelta
	{
		cur_14 := this_.instancesByChunk.IterValues()
		for cur_14.HasNext() {

			instanceInfo := cur_14.
				Next()
			_ = instanceInfo
			allModulesDelta = (allModulesDelta + instanceInfo.totalDelta)
		}
	}
	return ((allModulesDelta + this_.outer_.pattern.globalOverhead) < int32(0))
}

// Java method: GatherExtractionInfo.shouldExtractChunk(com.google.javascript.jscomp.JSChunk chunk)
func (this_ *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) priv_GatherExtractionInfo_shouldExtractChunk_2864(chunk *JSChunk) bool {
	_ = chunk
	var instanceInfo *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo = this_.instancesByChunk.Get(chunk)
	_ = instanceInfo
	if instanceInfo == nil {
		return false
	}
	return ((instanceInfo.totalDelta + this_.outer_.pattern.globalOverhead) < int32(0))
}

// Java class: com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.ExtractionInstance
type ExtractPrototypeMemberDeclarations_ExtractionInstance struct {
	outer_       *ExtractPrototypeMemberDeclarations
	declarations []*ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration `gson:"declarations"`
	delta        int32                                                            `gson:"delta"`
	parent       pkg_rhino.Node                                                   `gson:"parent"`
	idHash_      atomic.Int32
}

func (this_ *ExtractPrototypeMemberDeclarations_ExtractionInstance) BindSelf_(s any) {
	_ = s
}
func (this_ *ExtractPrototypeMemberDeclarations_ExtractionInstance) IsJavaNil_() bool {
	return this_ == nil
}
func (*ExtractPrototypeMemberDeclarations_ExtractionInstance) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations$ExtractionInstance", "ExtractionInstance"
}
func (this_ *ExtractPrototypeMemberDeclarations_ExtractionInstance) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ExtractPrototypeMemberDeclarations_ExtractionInstance) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.ExtractionInstance", "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations$ExtractionInstance", func(v_ any) bool {
		return j2gort.InstanceOf[*ExtractPrototypeMemberDeclarations_ExtractionInstance](v_)
	})
}
func (this_ *ExtractPrototypeMemberDeclarations_ExtractionInstance) SetOuter_(o *ExtractPrototypeMemberDeclarations) {
	this_.outer_ = o
}
func (this_ *ExtractPrototypeMemberDeclarations_ExtractionInstance) GetOuter_() *ExtractPrototypeMemberDeclarations {
	return this_.outer_
}

func (this_ *ExtractPrototypeMemberDeclarations_ExtractionInstance) initFields_ExtractPrototypeMemberDeclarations_ExtractionInstance() {
	_ = this_
	this_.declarations = []*ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration{

		// Java constructor: ExtractionInstance(com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.PrototypeMemberDeclaration head, com.google.javascript.rhino.Node parent)
	}

}

func NewExtractPrototypeMemberDeclarations_ExtractionInstance(outer_ *ExtractPrototypeMemberDeclarations, head *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration, parent pkg_rhino.Node) *ExtractPrototypeMemberDeclarations_ExtractionInstance {
	this_ := &ExtractPrototypeMemberDeclarations_ExtractionInstance{}
	this_.outer_ = outer_
	InitExtractPrototypeMemberDeclarations_ExtractionInstance(this_, head, parent)
	return this_
}

func InitExtractPrototypeMemberDeclarations_ExtractionInstance(this_ *ExtractPrototypeMemberDeclarations_ExtractionInstance, head *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration, parent pkg_rhino.Node) {
	_ = head
	_ = parent
	this_.initFields_ExtractPrototypeMemberDeclarations_ExtractionInstance()
	this_.parent = parent
	this_.declarations = append(this_.declarations, head)
	this_.delta = (this_.outer_.pattern.perExtractionOverhead + this_.outer_.pattern.perMemberOverhead)
	{
		var cur pkg_rhino.Node = head.node.GetNext()
		_ = cur
		for ; cur != 0; cur = cur.GetNext() {
			if cur.IsFunction() {
				continue
			}
			var prototypeMember *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration = ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_ExtractDeclaration(cur)
			_ = prototypeMember
			if (prototypeMember == nil) || !(head.priv_PrototypeMemberDeclaration_isSameClass_2866(prototypeMember)) {
				break
			}
			this_.declarations = append(this_.declarations, prototypeMember)
			this_.delta = (this_.delta + this_.outer_.pattern.perMemberOverhead)
		}
	}
}

// Java method: ExtractionInstance.isFavorable()
func (this_ *ExtractPrototypeMemberDeclarations_ExtractionInstance) IsFavorable() bool {
	return (this_.delta <= int32(0))
}

// Java class: com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.PrototypeMemberDeclaration
type ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration struct {
	memberName         string         `gson:"memberName"`
	qualifiedClassName string         `gson:"qualifiedClassName"`
	node               pkg_rhino.Node `gson:"node"`
	lhs                pkg_rhino.Node `gson:"lhs"`
	idHash_            atomic.Int32
	constant           bool `gson:"constant"`
}

func (this_ *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration) BindSelf_(s any) {
	_ = s
}
func (this_ *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration) IsJavaNil_() bool {
	return this_ == nil
}
func (*ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations$PrototypeMemberDeclaration", "PrototypeMemberDeclaration"
}
func (this_ *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.PrototypeMemberDeclaration", "com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations$PrototypeMemberDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration](v_)
	})
}

func (this_ *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration) initFields_ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration() {
	_ = this_
	this_.memberName = j2gort.NullString
	this_.qualifiedClassName = j2gort.NullString
}

// Java constructor: PrototypeMemberDeclaration(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node node)
func NewExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration(lhs pkg_rhino.Node, node pkg_rhino.Node) *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration {
	this_ := &ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration{}
	InitExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration(this_, lhs, node)
	return this_
}

func InitExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration(this_ *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration, lhs pkg_rhino.Node, node pkg_rhino.Node) {
	_ = lhs
	_ = node
	this_.initFields_ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration()
	if !(NodeUtil_IsExprAssign(node)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
	}
	this_.lhs = lhs
	this_.memberName = NodeUtil_GetPrototypePropertyName(lhs)
	this_.node = node
	var classNode pkg_rhino.Node = ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_GetPrototypeClassName(lhs)
	_ = classNode
	this_.qualifiedClassName = classNode.GetQualifiedName()
	this_.constant = classNode.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME())
}

// Java method: PrototypeMemberDeclaration.isSameClass(com.google.javascript.jscomp.ExtractPrototypeMemberDeclarations.PrototypeMemberDeclaration other)
func (this_ *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration) priv_PrototypeMemberDeclaration_isSameClass_2866(other *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration) bool {
	_ = other
	{

		sr_16 := this_.qualifiedClassName
		se_17 := other.qualifiedClassName
		return j2gort.CheckNotNullString(sr_16) == se_17
	}

}

// Java static method: PrototypeMemberDeclaration.getPrototypeClassName(com.google.javascript.rhino.Node qName)
func ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_GetPrototypeClassName(qName pkg_rhino.Node) pkg_rhino.Node {
	_ = qName
	var cur pkg_rhino.Node = qName
	_ = cur
	for cur.IsGetProp() {
		if j2gort.CheckNotNullString(cur.GetString()) == "prototype" {
			return cur.GetFirstChild()
		} else {
			cur = cur.GetFirstChild()
		}
	}
	return pkg_rhino.Node(0)
}

// Java static method: PrototypeMemberDeclaration.isPrototypePropertyDeclaration(com.google.javascript.rhino.Node n)
func ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_IsPrototypePropertyDeclaration(n pkg_rhino.Node) bool {
	_ = n
	if !(NodeUtil_IsExprAssign(n)) {
		return false
	}
	var lvalue pkg_rhino.Node = n.GetFirstFirstChild()
	_ = lvalue
	if lvalue.IsGetProp() {
		var cur pkg_rhino.Node = lvalue.GetFirstChild()
		_ = cur
		for cur.IsGetProp() {
			if j2gort.CheckNotNullString(cur.GetString()) == "prototype" {
				return cur.IsQualifiedName()
			}
			cur = cur.GetFirstChild()
		}
	}
	return false
}

// Java static method: PrototypeMemberDeclaration.extractDeclaration(com.google.javascript.rhino.Node n)
func ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_ExtractDeclaration(n pkg_rhino.Node) *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration {
	_ = n
	if !(ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_IsPrototypePropertyDeclaration(n)) {
		return nil
	}
	var lhs pkg_rhino.Node = n.GetFirstFirstChild()
	_ = lhs
	return NewExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration(lhs, n)
}
func (*ExtractPrototypeMemberDeclarations_GatherExtractionInfo) CanReleaseSyntacticScopes_() {}
