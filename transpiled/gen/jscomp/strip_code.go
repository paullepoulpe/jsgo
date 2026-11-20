/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/StripCode.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_strip_code_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_strip_code_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_strip_code_3 j2gort.Ext_java_util_function_Function
var lambdaSingleton_strip_code_4 j2gort.Ext_java_util_function_IntFunction
var lambdaSingleton_strip_code_5 j2gort.Ext_java_util_function_Function
var lambdaSingleton_strip_code_6 j2gort.Ext_java_util_function_IntFunction
var lambdaSingleton_strip_code_7 j2gort.Ext_java_util_function_Function
var lambdaSingleton_strip_code_8 j2gort.Ext_java_util_function_IntFunction
var lambdaSingleton_strip_code_9 j2gort.Ext_java_util_function_Function
var lambdaSingleton_strip_code_10 j2gort.Ext_java_util_function_Supplier

func init() {
	lambdaSingleton_strip_code_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return StripCode_ToStreamWithCollapsedVersions(j2gort.CastString(mrp0_)) })
	lambdaSingleton_strip_code_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return StripCode_ToStreamWithCollapsedVersions(j2gort.CastString(mrp0_)) })
	lambdaSingleton_strip_code_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return StripCode_ToStreamWithCollapsedVersions(j2gort.CastString(mrp0_)) })
	lambdaSingleton_strip_code_4 = j2gort.NewLambda_Ext_java_util_function_IntFunction(func(mrp0_ int32) any { return j2gort.SliceToAny(j2gort.NewStringArray(mrp0_)) })
	lambdaSingleton_strip_code_5 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_1 any) any {
		s := j2gort.CastString(lp_1)
		_ = s
		return j2gort.Ext_java_util_stream_Stream_Of_1_slice_any([]any{j2gort.StrValString(s) + ".", j2gort.StrValString(s) + "$"})
	})
	lambdaSingleton_strip_code_6 = j2gort.NewLambda_Ext_java_util_function_IntFunction(func(mrp0_ int32) any { return j2gort.SliceToAny(j2gort.NewStringArray(mrp0_)) })
	lambdaSingleton_strip_code_7 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_5 any) any {
		s := j2gort.CastString(lp_5)
		_ = s
		return j2gort.BoxString(javalang.StringToLowerCaseASCII(s))
	})
	lambdaSingleton_strip_code_8 = j2gort.NewLambda_Ext_java_util_function_IntFunction(func(mrp0_ int32) any { return j2gort.SliceToAny(j2gort.NewStringArray(mrp0_)) })
	lambdaSingleton_strip_code_9 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_9 any) any {
		s := j2gort.CastString(lp_9)
		_ = s
		return j2gort.BoxString(javalang.StringToLowerCaseASCII(s))
	})
	lambdaSingleton_strip_code_10 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "removing function call"
	})
}

// Java class: com.google.javascript.jscomp.StripCode
type StripCode struct {
	compiler                       AbstractCompiler                                  `gson:"compiler"`
	stripNameSuffixes              j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"stripNameSuffixes"`
	stripNamePrefixes              j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"stripNamePrefixes"`
	varsToRemove                   *pkg_base.LinkedIdentityHashSet[string]           `gson:"varsToRemove"`
	stripTypesList                 []string                                          `gson:"stripTypesList"`
	stripTypePrefixesList          []string                                          `gson:"stripTypePrefixesList"`
	stripNamePrefixesLowerCaseList []string                                          `gson:"stripNamePrefixesLowerCaseList"`
	stripNameSuffixesLowerCaseList []string                                          `gson:"stripNameSuffixesLowerCaseList"`
	idHash_                        atomic.Int32
}

func (this_ *StripCode) BindSelf_(s any) {
	_ = s
}
func (this_ *StripCode) IsJavaNil_() bool { return this_ == nil }
func (*StripCode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.StripCode", "StripCode"
}
func (this_ *StripCode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*StripCode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.StripCode", "com.google.javascript.jscomp.StripCode", func(v_ any) bool {
		return j2gort.InstanceOf[*StripCode](v_)
	})
}
func (*StripCode) Implements_jscomp_CompilerPass() {}

var static_StripCode_STRIP_TYPE_INHERIT_ERROR *DiagnosticType

func StaticGet_StripCode_STRIP_TYPE_INHERIT_ERROR() *DiagnosticType {
	clinit_StripCode()
	return static_StripCode_STRIP_TYPE_INHERIT_ERROR
}
func StaticSet_StripCode_STRIP_TYPE_INHERIT_ERROR(v *DiagnosticType) {
	clinit_StripCode()
	static_StripCode_STRIP_TYPE_INHERIT_ERROR = v
}

var static_StripCode_STRIP_ASSIGNMENT_ERROR *DiagnosticType

func StaticGet_StripCode_STRIP_ASSIGNMENT_ERROR() *DiagnosticType {
	clinit_StripCode()
	return static_StripCode_STRIP_ASSIGNMENT_ERROR
}
func StaticSet_StripCode_STRIP_ASSIGNMENT_ERROR(v *DiagnosticType) {
	clinit_StripCode()
	static_StripCode_STRIP_ASSIGNMENT_ERROR = v
}

var clinitDone_StripCode uint32
var clinitState_StripCode j2gort.ClinitState

func clinit_StripCode() {
	if atomic.LoadUint32(&clinitDone_StripCode) != 0 {
		return
	}
	clinitSlow_StripCode()
}

func clinitSlow_StripCode() {
	if !j2gort.EnterClinit(&clinitState_StripCode, &clinitDone_StripCode) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_StripCode, &clinitDone_StripCode)
	static_StripCode_STRIP_TYPE_INHERIT_ERROR = DiagnosticType_Error("JSC_STRIP_TYPE_INHERIT_ERROR", "Non-strip type {0} cannot inherit from strip type {1}")
	static_StripCode_STRIP_ASSIGNMENT_ERROR = DiagnosticType_Error("JSC_STRIP_ASSIGNMENT_ERROR", "Unable to strip assignment to {0}")
}

func (this_ *StripCode) initFields_StripCode() {
	_ = this_
	this_.varsToRemove = pkg_base.NewLinkedIdentityHashSet[string]()
}

// Java constructor: StripCode(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.common.collect.ImmutableSet<java.lang.String> stripTypes, com.google.common.collect.ImmutableSet<java.lang.String> stripNameSuffixes, com.google.common.collect.ImmutableSet<java.lang.String> stripNamePrefixes, boolean enableTweakStripping)
func NewStripCode(compiler AbstractCompiler, stripTypes j2gort.Ext_com_google_common_collect_ImmutableSet, stripNameSuffixes j2gort.Ext_com_google_common_collect_ImmutableSet, stripNamePrefixes j2gort.Ext_com_google_common_collect_ImmutableSet, enableTweakStripping bool) *StripCode {
	clinit_StripCode()
	this_ := &StripCode{}
	InitStripCode(this_, compiler, stripTypes, stripNameSuffixes, stripNamePrefixes, enableTweakStripping)
	return this_
}

func InitStripCode(this_ *StripCode, compiler AbstractCompiler, stripTypes j2gort.Ext_com_google_common_collect_ImmutableSet, stripNameSuffixes j2gort.Ext_com_google_common_collect_ImmutableSet, stripNamePrefixes j2gort.Ext_com_google_common_collect_ImmutableSet, enableTweakStripping bool) {
	_ = compiler
	_ = stripTypes
	_ = stripNameSuffixes
	_ = stripNamePrefixes
	_ = enableTweakStripping
	this_.initFields_StripCode()
	this_.compiler = compiler
	this_.stripNameSuffixes = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](stripNameSuffixes.Stream().FlatMap(lambdaSingleton_strip_code_1).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet()))
	this_.stripNamePrefixes = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](stripNamePrefixes.Stream().FlatMap(lambdaSingleton_strip_code_2).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet()))
	var stripTypesStream j2gort.Ext_java_util_stream_Stream = stripTypes.Stream()
	_ = stripTypesStream
	if enableTweakStripping {
		stripTypesStream = j2gort.Ext_java_util_stream_Stream_Concat(stripTypesStream, j2gort.Ext_java_util_stream_Stream_Of_1_any("goog.tweak"))
	}
	var stripTypesAdjusted j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](stripTypesStream.FlatMap(lambdaSingleton_strip_code_3).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet()))
	_ = stripTypesAdjusted
	this_.stripTypesList = j2gort.ColToNewStringArray(stripTypesAdjusted)
	this_.stripTypePrefixesList = func() []string {
		arr_2 := stripTypesAdjusted.Stream().FlatMap(lambdaSingleton_strip_code_5).ToArray_1(lambdaSingleton_strip_code_4)
		if arr_2 == nil {
			return nil
		}
		out_3 := make([]string, len(arr_2))
		for i_, el_4 := range arr_2 {
			out_3[i_] = j2gort.CastString(el_4)
		}
		return out_3
	}()
	this_.stripNamePrefixesLowerCaseList = func() []string {
		arr_6 := this_.stripNamePrefixes.Stream().Map(lambdaSingleton_strip_code_7).ToArray_1(lambdaSingleton_strip_code_6)
		if arr_6 == nil {
			return nil
		}
		out_7 := make([]string, len(arr_6))
		for i_, el_8 := range arr_6 {
			out_7[i_] = j2gort.CastString(el_8)
		}
		return out_7
	}()
	this_.stripNameSuffixesLowerCaseList = func() []string {
		arr_10 := this_.stripNameSuffixes.Stream().Map(lambdaSingleton_strip_code_9).ToArray_1(lambdaSingleton_strip_code_8)
		if arr_10 == nil {
			return nil
		}
		out_11 := make([]string, len(arr_10))
		for i_, el_12 := range arr_10 {
			out_11[i_] = j2gort.CastString(el_12)
		}
		return out_11
	}()
}

// Java static method: StripCode.toStreamWithCollapsedVersions(java.lang.String s)
func StripCode_ToStreamWithCollapsedVersions(s string) j2gort.Ext_java_util_stream_Stream {
	clinit_StripCode()
	_ = s
	var possibleForms j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = possibleForms
	possibleForms.Add_1(j2gort.BoxString(s))
	var chars []uint16 = javalang.StringToUTF16(s)
	_ = chars
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(chars)); i = i + 1 {
			if int32(chars[i]) == int32(uint16(46)) {
				{
					a_13 := chars
					i_14 := i
					a_13[i_14] = uint16(36)
				}
				possibleForms.Add_1(javalang.StringFromUTF16(chars))
			}
		}
	}
	return possibleForms.Stream()
}

// Java method: StripCode.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *StripCode) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	func() {
		var decisionsLog pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenIndexedLog(j2gort.GetObjectClass(this_), "decisions.log", []string{})
		_ = decisionsLog
		defer j2gort.CloseResource(decisionsLog)
		if decisionsLog.IsLogging() {
			decisionsLog.Log_1_Ext_java_util_function_Supplier(NewStripCode_StripCodeConfigRecord(this_))
		}
		decisionsLog.Log_1_string("\n=== decisions ===\n")
		NodeTraversal_Traverse(this_.compiler, root, NewStripCode_Strip(this_, decisionsLog))
	}()
}

// Java class: com.google.javascript.jscomp.StripCode.StripCodeConfigRecord
type StripCode_StripCodeConfigRecord struct {
	outer_  *StripCode
	idHash_ atomic.Int32
}

func (this_ *StripCode_StripCodeConfigRecord) BindSelf_(s any) {
	_ = s
}
func (this_ *StripCode_StripCodeConfigRecord) IsJavaNil_() bool { return this_ == nil }
func (*StripCode_StripCodeConfigRecord) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.StripCode$StripCodeConfigRecord", "StripCodeConfigRecord"
}
func (this_ *StripCode_StripCodeConfigRecord) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*StripCode_StripCodeConfigRecord) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.StripCode.StripCodeConfigRecord", "com.google.javascript.jscomp.StripCode$StripCodeConfigRecord", func(v_ any) bool {
		return j2gort.InstanceOf[*StripCode_StripCodeConfigRecord](v_)
	})
}
func (this_ *StripCode_StripCodeConfigRecord) SetOuter_(o *StripCode)                       { this_.outer_ = o }
func (this_ *StripCode_StripCodeConfigRecord) GetOuter_() *StripCode                        { return this_.outer_ }
func (*StripCode_StripCodeConfigRecord) Implements_j2gort_Ext_java_util_function_Supplier() {}

func (this_ *StripCode_StripCodeConfigRecord) initFields_StripCode_StripCodeConfigRecord() {
	_ = this_
}

// Java constructor: StripCodeConfigRecord()
func NewStripCode_StripCodeConfigRecord(outer_ *StripCode) *StripCode_StripCodeConfigRecord {
	this_ := &StripCode_StripCodeConfigRecord{}
	this_.outer_ = outer_
	InitStripCode_StripCodeConfigRecord(this_)
	return this_
}

func InitStripCode_StripCodeConfigRecord(this_ *StripCode_StripCodeConfigRecord) {
	this_.initFields_StripCode_StripCodeConfigRecord()
}

// Java method: StripCodeConfigRecord.get()
func (this_ *StripCode_StripCodeConfigRecord) Get() any {
	var builder j2gort.LocalStringBuilder

	builder.AppendStr("=== stripNameSuffixes ===\n")
	{
		cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.outer_.stripNameSuffixes))
		for cur_15.HasNext() {
			stripNameSuffix := cur_15.NextStr()
			_ = stripNameSuffix
			builder.AppendStr(stripNameSuffix).AppendStr("\n")
		}
	}
	builder.AppendStr("\n")
	builder.AppendStr("=== stripNamePrefixes ===\n")
	{
		cur_16 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.outer_.stripNamePrefixes))
		for cur_16.HasNext() {
			stripNamePrefix := cur_16.NextStr()
			_ = stripNamePrefix
			builder.AppendStr(stripNamePrefix).AppendStr("\n")
		}
	}
	builder.AppendStr("\n")
	builder.AppendStr("=== stripTypesList ===\n")
	for _, el_17 := range j2gort.CheckNotNullSlice(this_.outer_.stripTypesList) {
		stripType := el_17
		_ = stripType
		builder.AppendStr(stripType).AppendStr("\n")
	}
	builder.AppendStr("\n")
	builder.AppendStr("=== stripTypePrefixesList ===\n")
	for _, el_18 := range j2gort.CheckNotNullSlice(this_.outer_.stripTypePrefixesList) {
		stripNamePrefix_2 := el_18
		_ = stripNamePrefix_2
		builder.AppendStr(stripNamePrefix_2).AppendStr("\n")
	}
	builder.AppendStr("\n")
	return j2gort.BoxString(builder.String())
}

// Java class: com.google.javascript.jscomp.StripCode.Strip
type StripCode_Strip struct {
	outer_       *StripCode
	decisionsLog pkg_diagnostic.LogFile `gson:"decisionsLog"`
	idHash_      atomic.Int32
}

func (this_ *StripCode_Strip) BindSelf_(s any) {
	_ = s
}
func (this_ *StripCode_Strip) IsJavaNil_() bool { return this_ == nil }
func (*StripCode_Strip) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.StripCode$Strip", "Strip"
}
func (this_ *StripCode_Strip) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*StripCode_Strip) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.StripCode.Strip", "com.google.javascript.jscomp.StripCode$Strip", func(v_ any) bool {
		return j2gort.InstanceOf[*StripCode_Strip](v_)
	})
}
func (this_ *StripCode_Strip) SetOuter_(o *StripCode)              { this_.outer_ = o }
func (this_ *StripCode_Strip) GetOuter_() *StripCode               { return this_.outer_ }
func (*StripCode_Strip) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *StripCode_Strip) initFields_StripCode_Strip() {
	_ = this_
}

// Java constructor: Strip(com.google.javascript.jscomp.diagnostic.LogFile decisionsLog)
func NewStripCode_Strip(outer_ *StripCode, decisionsLog pkg_diagnostic.LogFile) *StripCode_Strip {
	this_ := &StripCode_Strip{}
	this_.outer_ = outer_
	InitStripCode_Strip(this_, decisionsLog)
	return this_
}

func InitStripCode_Strip(this_ *StripCode_Strip, decisionsLog pkg_diagnostic.LogFile) {
	_ = decisionsLog
	this_.initFields_StripCode_Strip()
	this_.decisionsLog = decisionsLog
}

// Java method: Strip.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *StripCode_Strip) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(28) /* Token.CALL */, int32(23) /* Token.NEW */ :
		if this_.IsMethodOrCtorCallThatTriggersRemoval(t, n, parent) {
			this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(lambdaSingleton_strip_code_10)
			this_.ReplaceHighestNestedCallWithNull(t, n, parent)
			return false
		} else {
			return true
		}
	default:
		return true
	}
	panic("unreachable")
}

// Java method: Strip.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *StripCode_Strip) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(84) /* Token.VAR */, int32(97) /* Token.CONST */, int32(111) /* Token.LET */ :
		this_.RemoveVarDeclarationsByNameOrRvalue(t, n, parent)
	case int32(32) /* Token.NAME */ :
		this_.MaybeRemoveReferenceToRemovedVariable(t, n, parent)
	case int32(51) /* Token.ASSIGN */, int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(58) /* Token.ASSIGN_ADD */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */ :
		this_.MaybeEliminateAssignmentByLvalueName(t, n, parent)
	case int32(47) /* Token.OBJECTLIT */ :
		this_.EliminateKeysWithStripNamesFromObjLit(n)
	case int32(93) /* Token.EXPR_RESULT */ :
		this_.MaybeEliminateExpressionByName(n)
	case int32(105) /* Token.CLASS */ :
		this_.MaybeEliminateClassByNameOrExtends(t, n, parent)
	default:
	}
}

// Java method: Strip.removeVarDeclarationsByNameOrRvalue(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *StripCode_Strip) RemoveVarDeclarationsByNameOrRvalue(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var next pkg_rhino.Node = pkg_rhino.Node(0)
	_ = next
	{
		var nameNode pkg_rhino.Node = n.GetFirstChild()
		_ = nameNode
		for ; nameNode != 0; nameNode = next {
			next = nameNode.GetNext()
			if nameNode.IsDestructuringLhs() {
				continue
			}
			if !(nameNode.IsName()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(nameNode))
			}
			var name string = nameNode.GetString()
			_ = name
			var lastDollarSign int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(36)))
			_ = lastDollarSign
			var possibleStripName string
			_ = possibleStripName
			if lastDollarSign != int32(-1) {
				possibleStripName = j2gort.JStrSubFrom(name, (lastDollarSign + int32(1)))
			} else {
				possibleStripName = name
			}
			if (this_.IsStripName(possibleStripName) || this_.QualifiedNameBeginsWithStripType_1_Node(nameNode)) || this_.IsCallWhoseReturnValueShouldBeStripped(nameNode.GetFirstChild()) {
				this_.outer_.varsToRemove.Add(name)
				if j2gort.JStrIndexOf(name, "$") >= 0 {
					if this_.decisionsLog.IsLogging() {
						this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
							return j2gort.BoxString((j2gort.StrValString(name) + ": initialize with null (" + j2gort.StrValString(possibleStripName) + ")"))
						}))
					}
					if nameNode.HasChildren() {
						this_.ReplaceWithNull(nameNode.GetOnlyChild())
					} else {
						nameNode.AddChildToFront(pkg_rhino.IR_NullNode().Srcref(nameNode))
					}
					t.ReportCodeChange_0()
				} else {
					if this_.decisionsLog.IsLogging() {
						this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
							return j2gort.StrValString(name) + ": removing declaration"
						}))
					}
					nameNode.Detach()
					NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(nameNode, this_.outer_.compiler)
				}
			}
		}
	}
	if !(n.HasChildren()) {
		this_.ReplaceWithEmpty(n, parent)
		t.ReportCodeChange_0()
	}
}

// Java method: Strip.maybeRemoveReferenceToRemovedVariable(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *StripCode_Strip) MaybeRemoveReferenceToRemovedVariable(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(84) /* Token.VAR */, int32(97) /* Token.CONST */, int32(111) /* Token.LET */ :
	case int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
		if (parent.GetFirstChild() == n) && this_.IsReferenceToRemovedVar(n) {
			if this_.decisionsLog.IsLogging() {
				this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
					return j2gort.StrValString(n.GetString()) + ": removing getelem/getprop/call chain"
				}))
			}
			this_.ReplaceHighestNestedCallWithNull(t, parent, parent.GetParent())
		}
	case int32(51) /* Token.ASSIGN */, int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(58) /* Token.ASSIGN_ADD */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */ :
		if this_.IsReferenceToRemovedVar(n) {
			if parent.GetFirstChild() == n {
				var grandparent pkg_rhino.Node = parent.GetParent()
				_ = grandparent
				if this_.decisionsLog.IsLogging() {
					this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
						return j2gort.StrValString(n.GetQualifiedName()) + ": removing assignment to stripped var"
					}))
				}
				if grandparent.IsExprResult() {
					var greatGrandparent pkg_rhino.Node = grandparent.GetParent()
					_ = greatGrandparent
					this_.ReplaceWithEmpty(grandparent, greatGrandparent)
					t.ReportCodeChange_0()
				} else {
					var rvalue pkg_rhino.Node = n.GetNext()
					_ = rvalue
					rvalue.Detach()
					parent.ReplaceWith(rvalue)
					t.ReportCodeChange_0()
				}
			} else {
				if this_.decisionsLog.IsLogging() {
					this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
						return j2gort.StrValString(n.GetQualifiedName()) + ": replacing rhs reference with null"
					}))
				}
				this_.ReplaceWithNull(n)
				t.ReportCodeChange_0()
			}
		}
	case int32(23) /* Token.NEW */, int32(28) /* Token.CALL */ :
		if !(n.IsFirstChildOf(parent)) && this_.IsReferenceToRemovedVar(n) {
			if this_.decisionsLog.IsLogging() {
				this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
					return j2gort.StrValString(n.GetQualifiedName()) + ": replacing parameter reference with null"
				}))
			}
			this_.ReplaceWithNull(n)
			t.ReportCodeChange_0()
		}
	case int32(50) /* Token.COMMA */ :
		var grandparent_2 pkg_rhino.Node = parent.GetParent()
		_ = grandparent_2
		var isLastChild bool = (parent.GetLastChild() == n)
		_ = isLastChild
		var parentIsCallee bool = ((grandparent_2.IsCall() || grandparent_2.IsNew()) && parent.IsFirstChildOf(grandparent_2))
		_ = parentIsCallee
		var isSafeToRemove bool = (!(isLastChild) || !(parentIsCallee))
		_ = isSafeToRemove
		if isSafeToRemove && this_.IsReferenceToRemovedVar(n) {
			if this_.decisionsLog.IsLogging() {
				this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
					return j2gort.StrValString(n.GetQualifiedName()) + ": replacing reference in comma expr with null"
				}))
			}
			this_.ReplaceWithNull(n)
			t.ReportCodeChange_0()
		}
	default:
		if this_.IsReferenceToRemovedVar(n) {
			if this_.decisionsLog.IsLogging() {
				this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
					return j2gort.StrValString(n.GetQualifiedName()) + ": replacing reference with null"
				}))
			}
			this_.ReplaceWithNull(n)
			t.ReportCodeChange_0()
		}
	}
}

// Java method: Strip.replaceHighestNestedCallWithNull(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *StripCode_Strip) ReplaceHighestNestedCallWithNull(t *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = node
	_ = parent
	var ancestor pkg_rhino.Node = parent
	_ = ancestor
	var ancestorChild pkg_rhino.Node = node
	_ = ancestorChild
	var ancestorParent pkg_rhino.Node = pkg_rhino.Node(0)
	_ = ancestorParent
	for {
		ancestorParent = ancestor.GetParent()
		if ancestorParent == 0 {
			return
		}
		if ancestor.GetFirstChild() != ancestorChild {
			this_.ReplaceWithNull(ancestorChild)
			break
		}
		if ancestor.IsExprResult() {
			this_.ReplaceWithEmpty(ancestor, ancestorParent)
			break
		}
		if ancestor.IsAssign() {
			ancestor.ReplaceWith(ancestor.GetLastChild().Detach())
			break
		}
		if !(NodeUtil_IsNormalGet(ancestor)) && !(ancestor.IsCall()) {
			this_.ReplaceWithNull(ancestorChild)
			break
		}
		ancestorChild = ancestor
		ancestor = ancestorParent
	}
	t.ReportCodeChange_0()
}

// Java method: Strip.maybeEliminateAssignmentByLvalueName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *StripCode_Strip) MaybeEliminateAssignmentByLvalueName(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var lvalue pkg_rhino.Node = n.GetFirstChild()
	_ = lvalue
	if this_.NameIncludesFieldNameToStrip(lvalue) || this_.QualifiedNameBeginsWithStripType_1_Node(lvalue) {
		if parent.IsExprResult() {
			if this_.decisionsLog.IsLogging() {
				this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
					return j2gort.StrValString(lvalue.GetString()) + ": removing assignment statement"
				}))
			}
			var grandparent pkg_rhino.Node = parent.GetParent()
			_ = grandparent
			if grandparent != 0 {
				this_.ReplaceWithEmpty(parent, grandparent)
				this_.outer_.compiler.ReportChangeToEnclosingScope(grandparent)
			}
		} else {
			t.Report_3(n, StaticGet_StripCode_STRIP_ASSIGNMENT_ERROR(), []string{lvalue.GetQualifiedName()})
		}
	}
}

// Java method: Strip.maybeEliminateExpressionByName(com.google.javascript.rhino.Node n)
func (this_ *StripCode_Strip) MaybeEliminateExpressionByName(n pkg_rhino.Node) {
	_ = n
	if !(n.IsExprResult()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent == 0 {
		return
	}
	var expression pkg_rhino.Node = n.GetFirstChild()
	_ = expression
	if this_.NameIncludesFieldNameToStrip(expression) || this_.QualifiedNameBeginsWithStripType_1_Node(expression) {
		if this_.decisionsLog.IsLogging() {
			this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
				return j2gort.StrValString(expression.GetString()) + ": removing property declaration statement"
			}))
		}
		this_.ReplaceWithEmpty(n, parent)
		this_.outer_.compiler.ReportChangeToEnclosingScope(parent)
	}
}

// Java method: Strip.eliminateKeysWithStripNamesFromObjLit(com.google.javascript.rhino.Node n)
func (this_ *StripCode_Strip) EliminateKeysWithStripNamesFromObjLit(n pkg_rhino.Node) {
	_ = n
	var key pkg_rhino.Node = n.GetFirstChild()
	_ = key
	for key != 0 {
		switch key.TokenSwitchOrd_() {
		case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(100) /* Token.STRING_KEY */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
			if this_.IsStripName(key.GetString()) {
				var next pkg_rhino.Node = key.GetNext()
				_ = next
				key.Detach()
				NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(key, this_.outer_.compiler)
				key = next
				this_.outer_.compiler.ReportChangeToEnclosingScope(n)
				break
			}
			fallthrough
		default:
			key = key.GetNext()
		}
	}
}

// Java method: Strip.maybeEliminateClassByNameOrExtends(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node classNode, com.google.javascript.rhino.Node parent)
func (this_ *StripCode_Strip) MaybeEliminateClassByNameOrExtends(t *NodeTraversal, classNode pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = classNode
	_ = parent
	var nameNode pkg_rhino.Node = NodeUtil_GetNameNode(classNode)
	_ = nameNode
	var className string = j2gort.NullString
	_ = className
	if (nameNode != 0) && nameNode.IsQualifiedName() {
		className = nameNode.GetQualifiedName()
		if this_.QualifiedNameBeginsWithStripType_1_string(className) {
			if this_.decisionsLog.IsLogging() {
				this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
					return j2gort.StrValString(className) + ": removing class"
				}))
			}
			if NodeUtil_IsStatementParent(parent) {
				this_.ReplaceWithEmpty(classNode, parent)
			} else {
				this_.ReplaceWithNull(classNode)
			}
			t.ReportCodeChange_0()
			return
		}
	} else {
		className = "<anonymous>"
	}
	var superclassNode pkg_rhino.Node = classNode.GetSecondChild()
	_ = superclassNode
	if (superclassNode != 0) && superclassNode.IsQualifiedName() {
		var superclassName string = superclassNode.GetQualifiedName()
		_ = superclassName
		if this_.QualifiedNameBeginsWithStripType_1_string(superclassName) {
			t.Report_3(classNode, StaticGet_StripCode_STRIP_TYPE_INHERIT_ERROR(), []string{className, superclassName})
		}
	}
}

// Java method: Strip.isCallWhoseReturnValueShouldBeStripped(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node n)
func (this_ *StripCode_Strip) IsCallWhoseReturnValueShouldBeStripped(n pkg_rhino.Node) bool {
	_ = n
	if ((n == 0) || (!(n.IsCall()) && !(n.IsNew()))) || !(n.HasChildren()) {
		return false
	}
	var function pkg_rhino.Node = NodeUtil_GetCallTargetResolvingIndirectCalls(n)
	_ = function
	return (this_.QualifiedNameBeginsWithStripType_1_Node(function) || this_.NameIncludesFieldNameToStrip(function))
}

// Java method: Strip.qualifiedNameBeginsWithStripType(com.google.javascript.rhino.Node n)
func (this_ *StripCode_Strip) QualifiedNameBeginsWithStripType_1_Node(n pkg_rhino.Node) bool {
	_ = n
	var name string = n.GetQualifiedName()
	_ = name
	return this_.QualifiedNameBeginsWithStripType_1_string(name)
}

// Java method: Strip.qualifiedNameBeginsWithStripType(java.lang.String name)
func (this_ *StripCode_Strip) QualifiedNameBeginsWithStripType_1_string(name string) bool {
	_ = name
	if name != j2gort.NullString {
		for _, el_19 := range j2gort.CheckNotNullSlice(this_.outer_.stripTypesList) {
			type_ := el_19
			_ = type_
			if j2gort.CheckNotNullString(name) == type_ {
				this_.priv_Strip_logStripName_5225(name, "equals strip type")
				return true
			}
		}
		for _, el_20 := range j2gort.CheckNotNullSlice(this_.outer_.stripTypePrefixesList) {
			type__2 := el_20
			_ = type__2
			if javalang.StringStartsWith(j2gort.CheckNotNullString(name), j2gort.CheckNotNullString(type__2)) {
				this_.priv_Strip_logStripName_5225(name, "starts with strip type prefix")
				return true
			}
		}
	}
	this_.priv_Strip_logNotAStripName_5224(name, "does not begin with a strip type")
	return false
}

// Java method: Strip.isReferenceToRemovedVar(com.google.javascript.rhino.Node n)
func (this_ *StripCode_Strip) IsReferenceToRemovedVar(n pkg_rhino.Node) bool {
	_ = n
	return this_.outer_.varsToRemove.Contains(n.GetString())
}

// Java method: Strip.isMethodOrCtorCallThatTriggersRemoval(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *StripCode_Strip) IsMethodOrCtorCallThatTriggersRemoval(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	var function pkg_rhino.Node = NodeUtil_GetCallTargetResolvingIndirectCalls(n)
	_ = function
	if (function == 0) || !(function.IsQualifiedName()) {
		return false
	}
	if (parent != 0) && parent.IsName() {
		var grandparent pkg_rhino.Node = parent.GetParent()
		_ = grandparent
		if (grandparent != 0) && NodeUtil_IsNameDeclaration(grandparent) {
			return false
		}
	}
	if function.IsName() && this_.IsStripName(function.GetString()) {
		return true
	}
	var callee pkg_rhino.Node = function.GetFirstChild()
	_ = callee
	return (((this_.NameIncludesFieldNameToStrip(callee) || this_.NameIncludesFieldNameToStrip(function)) || this_.QualifiedNameBeginsWithStripType_1_Node(function)) || this_.priv_Strip_actsOnStripType_5222(t, n))
}

// Java method: Strip.nameIncludesFieldNameToStrip(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node n)
func (this_ *StripCode_Strip) NameIncludesFieldNameToStrip(n pkg_rhino.Node) bool {
	_ = n
	if n == 0 {
		return false
	} else if n.IsGetProp() {
		return (this_.IsStripName(n.GetString()) || this_.NameIncludesFieldNameToStrip(n.GetFirstChild()))
	} else if n.IsName() {
		var nameString string = n.GetString()
		_ = nameString
		var end int32 = j2gort.JStrIndexOfCp(nameString, int32(uint16(36)))
		_ = end
		if end == int32(-1) {
			return false
		}
		var start int32 = int32(0)
		_ = start
		for {
			if this_.IsStripName(j2gort.JStrSub(nameString, start, end)) {
				return true
			}
			if end == j2gort.JStrLen(nameString) {
				break
			}
			var newStart int32 = (end + int32(1))
			_ = newStart
			var newEnd int32 = j2gort.JStrIndexOfCpFrom(nameString, int32(uint16(36)), newStart)
			_ = newEnd
			start = newStart
			if newEnd == int32(-1) {
				end = j2gort.JStrLen(nameString)
			} else {
				end = newEnd
			}
		}
	}
	return false
}

// Java method: Strip.actsOnStripType(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node callNode)
func (this_ *StripCode_Strip) priv_Strip_actsOnStripType_5222(t *NodeTraversal, callNode pkg_rhino.Node) bool {
	_ = t
	_ = callNode
	var classes *CodingConvention_SubclassRelationship = this_.outer_.compiler.GetCodingConvention().GetClassesDefinedByCall(callNode)
	_ = classes
	if classes != nil {
		if this_.QualifiedNameBeginsWithStripType_1_string(classes.subclassName) {
			this_.priv_Strip_logStripName_5225(classes.subclassName, "class defining call")
			return true
		}
		if this_.QualifiedNameBeginsWithStripType_1_string(classes.superclassName) {
			t.Report_3(callNode, StaticGet_StripCode_STRIP_TYPE_INHERIT_ERROR(), []string{classes.subclassName, classes.superclassName})
		}
	}
	return false
}

// Java method: Strip.isStripName(java.lang.String name)
func (this_ *StripCode_Strip) IsStripName(name string) bool {
	_ = name
	if j2gort.SetContainsStr(this_.outer_.stripNameSuffixes, name) {
		this_.priv_Strip_logStripName_5225(name, "matches a suffix")
		return true
	}
	if j2gort.SetContainsStr(this_.outer_.stripNamePrefixes, name) {
		this_.priv_Strip_logStripName_5225(name, "matches a prefix")
		return true
	}
	if (j2gort.CheckNotNullString(name) == "") || j2gort.BoxChar_IsUpperCase_1_uint16(j2gort.JStrCharAt(name, int32(0))) {
		this_.priv_Strip_logNotAStripName_5224(name, "empty or starts with uppercase")
		return false
	}
	var lcName string = javalang.StringToLowerCaseASCII(name)
	_ = lcName
	for _, el_21 := range j2gort.CheckNotNullSlice(this_.outer_.stripNamePrefixesLowerCaseList) {
		stripName := el_21
		_ = stripName
		if javalang.StringStartsWith(j2gort.CheckNotNullString(lcName), j2gort.CheckNotNullString(stripName)) {
			if this_.decisionsLog.IsLogging() {
				this_.priv_Strip_logStripName_5226(name, j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
					return "matches lc prefix: " + j2gort.StrValString(stripName)
				}))
			}
			return true
		}
	}
	for _, el_22 := range j2gort.CheckNotNullSlice(this_.outer_.stripNameSuffixesLowerCaseList) {
		stripName_2 := el_22
		_ = stripName_2
		if javalang.StringEndsWith(j2gort.CheckNotNullString(lcName), j2gort.CheckNotNullString(stripName_2)) {
			if this_.decisionsLog.IsLogging() {
				this_.priv_Strip_logStripName_5226(name, j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
					return "matches lc suffix: " + j2gort.StrValString(stripName_2)
				}))
			}
			return true
		}
	}
	this_.priv_Strip_logNotAStripName_5224(name, "no matches")
	return false
}

// Java method: Strip.logNotAStripName(java.lang.String name, java.lang.String reason)
func (this_ *StripCode_Strip) priv_Strip_logNotAStripName_5224(name string, reason string) {
	_ = name
	_ = reason
	if this_.decisionsLog.IsLogging() {
		this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
			return j2gort.BoxString((j2gort.StrValString(name) + "\tnot a strip name: " + j2gort.StrValString(reason)))
		}))
	}
}

// Java method: Strip.logStripName(java.lang.String name, java.lang.String reason)
func (this_ *StripCode_Strip) priv_Strip_logStripName_5225(name string, reason string) {
	_ = name
	_ = reason
	if this_.decisionsLog.IsLogging() {
		this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
			return j2gort.BoxString((j2gort.StrValString(name) + "\tstrip name: " + j2gort.StrValString(reason)))
		}))
	}
}

// Java method: Strip.logStripName(java.lang.String name, java.util.function.Supplier<java.lang.String> reasonSupplier)
func (this_ *StripCode_Strip) priv_Strip_logStripName_5226(name string, reasonSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = name
	_ = reasonSupplier
	if this_.decisionsLog.IsLogging() {
		this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
			return j2gort.BoxString((j2gort.StrValString(name) + "\tstrip name: " + j2gort.StrValString(j2gort.CastString(reasonSupplier.Get()))))
		}))
	}
}

// Java method: Strip.replaceWithNull(com.google.javascript.rhino.Node n)
func (this_ *StripCode_Strip) ReplaceWithNull(n pkg_rhino.Node) {
	_ = n
	if this_.decisionsLog.IsLogging() {
		this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
			return "replace with null: " + j2gort.StrValString(n.GetLocation())
		}))
	}
	n.ReplaceWith(pkg_rhino.IR_NullNode())
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(n, this_.outer_.compiler)
}

// Java method: Strip.replaceWithEmpty(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *StripCode_Strip) ReplaceWithEmpty(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	if this_.decisionsLog.IsLogging() {
		this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
			return "replace with empty: " + j2gort.StrValString(n.GetLocation())
		}))
	}
	NodeUtil_RemoveChild(parent, n)
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(n, this_.outer_.compiler)
}
func (*StripCode_Strip) CanReleaseSyntacticScopes_() {}
