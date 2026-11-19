/*
 * Copyright 2025 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/NestedCompilerRunner.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.NestedCompilerRunner
type NestedCompilerRunner struct {
	shadowCompiler   *Compiler        `gson:"shadowCompiler"`
	inputs           []*CompilerInput `gson:"inputs"`
	baseChunk        *JSChunk         `gson:"baseChunk"`
	original         AbstractCompiler `gson:"original"`
	shadowOptions    *CompilerOptions `gson:"shadowOptions"`
	shadowPassConfig PassConfig       `gson:"shadowPassConfig"`
	idHash_          atomic.Int32
	mode             NestedCompilerRunner_Mode `gson:"mode"`
}

func (this_ *NestedCompilerRunner) BindSelf_(s any) {
	_ = s
}
func (this_ *NestedCompilerRunner) IsJavaNil_() bool { return this_ == nil }
func (*NestedCompilerRunner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NestedCompilerRunner", "NestedCompilerRunner"
}
func (this_ *NestedCompilerRunner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NestedCompilerRunner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NestedCompilerRunner", "com.google.javascript.jscomp.NestedCompilerRunner", func(v_ any) bool {
		return j2gort.InstanceOf[*NestedCompilerRunner](v_)
	})
}

func (this_ *NestedCompilerRunner) initFields_NestedCompilerRunner() {
	_ = this_
	this_.shadowCompiler = NewCompiler()
	this_.inputs = []*CompilerInput{}
	this_.baseChunk = NewJSChunk("shadow_base")
}

// Java constructor: NestedCompilerRunner(com.google.javascript.jscomp.AbstractCompiler original, com.google.javascript.jscomp.CompilerOptions shadowOptions, com.google.javascript.jscomp.NestedCompilerRunner.Mode mode, com.google.javascript.jscomp.PassConfig shadowPassConfig)
func NewNestedCompilerRunner(original AbstractCompiler, shadowOptions *CompilerOptions, mode NestedCompilerRunner_Mode, shadowPassConfig PassConfig) *NestedCompilerRunner {
	this_ := &NestedCompilerRunner{}
	InitNestedCompilerRunner(this_, original, shadowOptions, mode, shadowPassConfig)
	return this_
}

func InitNestedCompilerRunner(this_ *NestedCompilerRunner, original AbstractCompiler, shadowOptions *CompilerOptions, mode NestedCompilerRunner_Mode, shadowPassConfig PassConfig) {
	_ = original
	_ = shadowOptions
	_ = mode
	_ = shadowPassConfig
	this_.initFields_NestedCompilerRunner()
	this_.original = original
	this_.shadowOptions = shadowOptions
	this_.mode = mode
	this_.shadowPassConfig = shadowPassConfig
}

// Java static method: NestedCompilerRunner.create(com.google.javascript.jscomp.AbstractCompiler original, com.google.javascript.jscomp.CompilerOptions shadowOptions, com.google.javascript.jscomp.NestedCompilerRunner.Mode mode, com.google.javascript.jscomp.PassConfig shadowPassConfig)
func NestedCompilerRunner_Create(original AbstractCompiler, shadowOptions *CompilerOptions, mode NestedCompilerRunner_Mode, shadowPassConfig PassConfig) *NestedCompilerRunner {
	_ = original
	_ = shadowOptions
	_ = mode
	_ = shadowPassConfig
	return NewNestedCompilerRunner(original, shadowOptions, mode, shadowPassConfig)
}

// Java method: NestedCompilerRunner.addScript(com.google.javascript.rhino.Node shadowScript, java.lang.String uniqueShadowId)
func (this_ *NestedCompilerRunner) AddScript(shadowScript pkg_rhino.Node, uniqueShadowId string) *NestedCompilerRunner {
	_ = shadowScript
	_ = uniqueShadowId
	if !(shadowScript.GetIsInClosureUnawareSubtree()) {
		j2gort.CheckArgumentFmt(false, "Expected closureUnaware script, found %s", j2gort.BoxNode(shadowScript))
	}
	var shadowFile *SourceFile = SourceFile_FromCode_3(uniqueShadowId, "", pkg_rhino.StaticGet_StaticSourceFile_SourceKind_STRONG())
	_ = shadowFile
	var val_1 CompilerInput
	InitCompilerInput_2_SourceFile_InputId(&val_1, shadowFile, pkg_rhino.NewInputId(uniqueShadowId))
	var shadowInput *CompilerInput = &val_1
	_ = shadowInput
	shadowInput.InitShadowAst(shadowScript)
	shadowScript.SetInputId(shadowInput.GetInputId())
	shadowScript.SetStaticSourceFile(pkg_rhino.ToIface_StaticSourceFile(shadowFile))
	this_.inputs = append(this_.inputs, shadowInput)
	return this_
}

// Java method: NestedCompilerRunner.compile()
func (this_ *NestedCompilerRunner) Compile() {
	this_.priv_NestedCompilerRunner_initializeCompiler_3919()
	this_.priv_NestedCompilerRunner_runCompilation_3921()
	this_.priv_NestedCompilerRunner_checkInvariants_3923()
}

// Java method: NestedCompilerRunner.initializeCompiler()
func (this_ *NestedCompilerRunner) priv_NestedCompilerRunner_initializeCompiler_3919() {
	var externNames j2gort.Ext_com_google_common_collect_ImmutableSet = NodeUtil_CollectExternVariableNames(this_.original, this_.original.GetRoot().GetFirstChild())
	_ = externNames
	this_.shadowCompiler.AddExportedNames(j2gort.AsIface[j2gort.Ext_java_util_Set](externNames))
	this_.shadowCompiler.SetExternProperties(this_.original.GetExternProperties())
	this_.shadowCompiler.SetDebugMessage("<shadow AST compilation>")
	var chunks j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_NestedCompilerRunner_createChunks_3920()
	_ = chunks
	this_.shadowCompiler.InitChunks(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()), j2gort.AsIface[j2gort.Ext_java_util_List](chunks), this_.shadowOptions)
	this_.shadowCompiler.SetPassConfig(this_.shadowPassConfig)
	this_.shadowCompiler.ParseForCompilation()
}

// Java method: NestedCompilerRunner.createChunks()
func (this_ *NestedCompilerRunner) priv_NestedCompilerRunner_createChunks_3920() j2gort.Ext_com_google_common_collect_ImmutableList {
	this_.baseChunk.Add_1_CompilerInput(NewCompilerInput_1(SourceFile_FromCode_3("synthetic_base", "", pkg_rhino.StaticGet_StaticSourceFile_SourceKind_STRONG())))
	var chunks *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}).AddSelf_(j2gort.PtrToAny(this_.baseChunk))
	_ = chunks
	for _, input := range this_.inputs {
		_ = input
		var chunk *JSChunk = NewJSChunk(("shadow_chunk_" + j2gort.StrValString(input.GetSourceFile().GetName())))
		_ = chunk
		chunk.Add_1_CompilerInput(input)
		chunk.AddDependency(this_.baseChunk)
		chunks.AddVoid_(chunk)
	}

	return chunks.BuildImmList_()
}

// Java method: NestedCompilerRunner.runCompilation()
func (this_ *NestedCompilerRunner) priv_NestedCompilerRunner_runCompilation_3921() {
	switch (this_.mode).Ordinal() {
	case int32(0) /* Mode.TRANSPILE_AND_OPTIMIZE */ :
		{
			clinit_VarCheck()
			var val_4 VarCheck
			InitVarCheck_1(&val_4, ToIface_AbstractCompiler(this_.shadowCompiler))
			(&val_4).Process_2(this_.shadowCompiler.GetExternsRoot(), this_.shadowCompiler.GetJsRoot())
		}
		this_.priv_NestedCompilerRunner_delegateDiagnosticsToOriginal_3922()
		if this_.shadowCompiler.HasHaltingErrors() {
			return
		}
		this_.shadowCompiler.Stage2Passes(StaticGet_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS())
		this_.priv_NestedCompilerRunner_delegateDiagnosticsToOriginal_3922()
		if this_.shadowCompiler.HasHaltingErrors() {
			return
		}
		this_.shadowCompiler.Stage3Passes()
		this_.priv_NestedCompilerRunner_delegateDiagnosticsToOriginal_3922()
	case int32(1) /* Mode.TRANSPILE_ONLY */ :
		this_.shadowCompiler.TranspileAndDontCheck()
		this_.priv_NestedCompilerRunner_delegateDiagnosticsToOriginal_3922()
	}
}

// Java method: NestedCompilerRunner.delegateDiagnosticsToOriginal()
func (this_ *NestedCompilerRunner) priv_NestedCompilerRunner_delegateDiagnosticsToOriginal_3922() {
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.shadowCompiler.GetErrors()))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			error_ := j2gort.CastPtr[JSError](rawElem_6)
			_ = error_
			this_.original.Report_1(error_)
		}
	}
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.shadowCompiler.GetWarnings()))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			error__2 := j2gort.CastPtr[JSError](rawElem_8)
			_ = error__2
			this_.original.Report_1(error__2)
		}
	}
}

// Java method: NestedCompilerRunner.checkInvariants()
func (this_ *NestedCompilerRunner) priv_NestedCompilerRunner_checkInvariants_3923() {
	if !(this_.baseChunk.GetInputs().Size() == int32(1)) {
		j2gort.CheckStateFmt(false, "Expected exactly 1 base chunk input, got %s", this_.baseChunk.GetInputs())
	}
	var baseInputNode pkg_rhino.Node = j2gort.CastPtr[CompilerInput](this_.baseChunk.GetInputs().Get(int32(0))).GetAstRoot(ToIface_AbstractCompiler(this_.shadowCompiler))
	_ = baseInputNode
	if !(!(baseInputNode.HasChildren())) {
		j2gort.CheckStateFmt(false, "Expected synthetic base input chunk to have zero children, but found %s", j2gort.BoxString(baseInputNode.ToStringTree()))
	}
}

// Java enum: com.google.javascript.jscomp.NestedCompilerRunner.Mode
type NestedCompilerRunner_Mode uint8

const (
	enum_NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE NestedCompilerRunner_Mode = 1
	enum_NestedCompilerRunner_Mode_TRANSPILE_ONLY         NestedCompilerRunner_Mode = 2
)

var enumBaseTable_NestedCompilerRunner_Mode = [...]j2gort.EnumBase{
	{Name_: "TRANSPILE_AND_OPTIMIZE", Ordinal_: 0},
	{Name_: "TRANSPILE_ONLY", Ordinal_: 1},
}
var enumValues_NestedCompilerRunner_Mode = []NestedCompilerRunner_Mode{1, 2}

func StaticGet_NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE() NestedCompilerRunner_Mode {
	return enum_NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE
}
func StaticGet_NestedCompilerRunner_Mode_TRANSPILE_ONLY() NestedCompilerRunner_Mode {
	return enum_NestedCompilerRunner_Mode_TRANSPILE_ONLY
}

func NestedCompilerRunner_Mode_Values() []NestedCompilerRunner_Mode {
	return j2gort.CloneSlice(enumValues_NestedCompilerRunner_Mode)
}

func NestedCompilerRunner_Mode_ValueOf(name string) NestedCompilerRunner_Mode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_NestedCompilerRunner_Mode {
		if enumBaseTable_NestedCompilerRunner_Mode[i_].Name_ == name {
			return NestedCompilerRunner_Mode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.NestedCompilerRunner.Mode." + name))
}

func (this_ NestedCompilerRunner_Mode) IsJavaNil_() bool                              { return this_ == 0 }
func (NestedCompilerRunner_Mode) HasIdentityEquals_()                                 {}
func (NestedCompilerRunner_Mode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (NestedCompilerRunner_Mode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ NestedCompilerRunner_Mode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (NestedCompilerRunner_Mode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NestedCompilerRunner$Mode", "Mode"
}
func (this_ NestedCompilerRunner_Mode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_NestedCompilerRunner_Mode[this_-1]
}
func (this_ NestedCompilerRunner_Mode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ NestedCompilerRunner_Mode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ NestedCompilerRunner_Mode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_NestedCompilerRunner_Mode[this_-1].Name_
}
func (this_ NestedCompilerRunner_Mode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_NestedCompilerRunner_Mode[this_-1].Name_
}
func (this_ NestedCompilerRunner_Mode) Equals(other_ any) bool {
	o_, ok_ := other_.(NestedCompilerRunner_Mode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ NestedCompilerRunner_Mode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_NestedCompilerRunner_Mode[this_-1].HashCode()
}
func (this_ NestedCompilerRunner_Mode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_NestedCompilerRunner_Mode[this_-1].CompareTo(other_)
}
func (this_ NestedCompilerRunner_Mode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_NestedCompilerRunner_Mode[this_-1].Name_)
}
func init() {
	vals_ := []any{NestedCompilerRunner_Mode(1), NestedCompilerRunner_Mode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.NestedCompilerRunner.Mode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.NestedCompilerRunner$Mode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NestedCompilerRunner.Mode", "com.google.javascript.jscomp.NestedCompilerRunner$Mode", func(v_ any) bool {
		return j2gort.InstanceOf[NestedCompilerRunner_Mode](v_)
	})
}
