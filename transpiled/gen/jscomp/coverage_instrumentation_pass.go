/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/instrumentation/CoverageInstrumentationPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass
type CoverageInstrumentationPass struct {
	compiler                           AbstractCompiler                                    `gson:"compiler"`
	instrumentationData                j2gort.OrderedMap[string, *FileInstrumentationData] `gson:"instrumentationData"`
	reach                              CoverageInstrumentationPass_CoverageReach           `gson:"reach"`
	instrumentOption                   CompilerOptions_InstrumentOption                    `gson:"instrumentOption"`
	productionInstrumentationArrayName string                                              `gson:"productionInstrumentationArrayName"`
	idHash_                            atomic.Int32
}

func (this_ *CoverageInstrumentationPass) BindSelf_(s any) {
	_ = s
}
func (this_ *CoverageInstrumentationPass) IsJavaNil_() bool { return this_ == nil }
func (*CoverageInstrumentationPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass", "CoverageInstrumentationPass"
}
func (this_ *CoverageInstrumentationPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CoverageInstrumentationPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass", "com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass", func(v_ any) bool {
		return j2gort.InstanceOf[*CoverageInstrumentationPass](v_)
	})
}
func (*CoverageInstrumentationPass) Implements_jscomp_CompilerPass() {}

func StaticGet_CoverageInstrumentationPass_JS_INSTRUMENTATION_OBJECT_NAME() string { return "__jscov" }
func (this_ *CoverageInstrumentationPass) initFields_CoverageInstrumentationPass() {
	_ = this_
	this_.productionInstrumentationArrayName = j2gort.NullString
}

// Java constructor: CoverageInstrumentationPass(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass.CoverageReach reach, com.google.javascript.jscomp.CompilerOptions.InstrumentOption instrumentOption, java.lang.String productionInstrumentationArrayName)
func NewCoverageInstrumentationPass_4(compiler AbstractCompiler, reach CoverageInstrumentationPass_CoverageReach, instrumentOption CompilerOptions_InstrumentOption, productionInstrumentationArrayName string) *CoverageInstrumentationPass {
	this_ := &CoverageInstrumentationPass{}
	InitCoverageInstrumentationPass_4(this_, compiler, reach, instrumentOption, productionInstrumentationArrayName)
	return this_
}

func InitCoverageInstrumentationPass_4(this_ *CoverageInstrumentationPass, compiler AbstractCompiler, reach CoverageInstrumentationPass_CoverageReach, instrumentOption CompilerOptions_InstrumentOption, productionInstrumentationArrayName string) {
	_ = compiler
	_ = reach
	_ = instrumentOption
	_ = productionInstrumentationArrayName
	this_.initFields_CoverageInstrumentationPass()
	this_.compiler = compiler
	this_.reach = reach
	this_.instrumentOption = instrumentOption
	this_.productionInstrumentationArrayName = productionInstrumentationArrayName
	this_.instrumentationData = j2gort.NewOrderedMap[string, *FileInstrumentationData]()
}

// Java constructor: CoverageInstrumentationPass(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass.CoverageReach reach)
func NewCoverageInstrumentationPass_2(compiler AbstractCompiler, reach CoverageInstrumentationPass_CoverageReach) *CoverageInstrumentationPass {
	this_ := &CoverageInstrumentationPass{}
	InitCoverageInstrumentationPass_2(this_, compiler, reach)
	return this_
}

func InitCoverageInstrumentationPass_2(this_ *CoverageInstrumentationPass, compiler AbstractCompiler, reach CoverageInstrumentationPass_CoverageReach) {
	_ = compiler
	_ = reach
	InitCoverageInstrumentationPass_4(this_, compiler, reach, StaticGet_CompilerOptions_InstrumentOption_LINE_ONLY(), "")
}

// Java method: CoverageInstrumentationPass.addHeaderCode(com.google.javascript.rhino.Node script)
func (this_ *CoverageInstrumentationPass) priv_CoverageInstrumentationPass_addHeaderCode_6504(script pkg_rhino.Node) {
	_ = script
	script.AddChildToFront(this_.priv_CoverageInstrumentationPass_createConditionalObjectDecl_6507("__jscov", script))
	script.AddChildToFront(this_.compiler.ParseSyntheticCode("coverage_instrumentation_header", "(function(self) { if (!self.window) { self.window = self; self.window.top = self; }})(typeof self !== 'undefined' ? self : globalThis)").RemoveFirstChild().SrcrefTreeIfMissing(script))
}

// Java method: CoverageInstrumentationPass.addProductionHeaderCode(com.google.javascript.rhino.Node script, java.lang.String arrayName)
func (this_ *CoverageInstrumentationPass) priv_CoverageInstrumentationPass_addProductionHeaderCode_6505(script pkg_rhino.Node, arrayName string) {
	_ = script
	_ = arrayName
	var arrayLit pkg_rhino.Node = pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{})
	_ = arrayLit
	var name pkg_rhino.Node = pkg_rhino.IR_Name(arrayName)
	_ = name
	var varNode pkg_rhino.Node = pkg_rhino.IR_Var_2(name, arrayLit)
	_ = varNode
	script.AddChildToFront(varNode.SrcrefTreeIfMissing(script))
}

// Java method: CoverageInstrumentationPass.checkIfArrayNameExternDeclared(com.google.javascript.rhino.Node externsNode, java.lang.String arrayName)
func (this_ *CoverageInstrumentationPass) priv_CoverageInstrumentationPass_checkIfArrayNameExternDeclared_6506(externsNode pkg_rhino.Node, arrayName string) {
	_ = externsNode
	_ = arrayName
	if !(j2gort.SetContainsStr(NodeUtil_CollectExternVariableNames(this_.compiler, externsNode), arrayName)) {
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "The array name passed to --production_instrumentation_array_name was not declared as an extern. This will result in undefined behaviour")))
	}
}

// Java method: CoverageInstrumentationPass.process(com.google.javascript.rhino.Node externsNode, com.google.javascript.rhino.Node rootNode)
func (this_ *CoverageInstrumentationPass) Process_2(externsNode pkg_rhino.Node, rootNode pkg_rhino.Node) {
	_ = externsNode
	_ = rootNode
	if rootNode.HasChildren() {
		if this_.instrumentOption == StaticGet_CompilerOptions_InstrumentOption_BRANCH_ONLY() {
			NodeTraversal_Traverse(this_.compiler, rootNode, NewBranchCoverageInstrumentationCallback(this_.compiler, &this_.instrumentationData))
		} else if this_.instrumentOption == StaticGet_CompilerOptions_InstrumentOption_PRODUCTION() {
			var val_1 ProductionCoverageInstrumentationCallback
			InitProductionCoverageInstrumentationCallback(&val_1, this_.compiler, this_.productionInstrumentationArrayName)
			var productionCoverageInstrumentationCallback *ProductionCoverageInstrumentationCallback = &val_1
			_ = productionCoverageInstrumentationCallback
			NodeTraversal_Traverse(this_.compiler, rootNode, productionCoverageInstrumentationCallback)
			this_.priv_CoverageInstrumentationPass_checkIfArrayNameExternDeclared_6506(externsNode, this_.productionInstrumentationArrayName)
			var instrumentationMapping *VariableMap = productionCoverageInstrumentationCallback.GetInstrumentationMapping()
			_ = instrumentationMapping
			this_.compiler.SetInstrumentationMapping(instrumentationMapping)
		} else {
			NodeTraversal_Traverse(this_.compiler, rootNode, NewCoverageInstrumentationCallback(&this_.instrumentationData, this_.reach))
		}
		var firstScript pkg_rhino.Node = rootNode.GetFirstChild()
		_ = firstScript
		if !(firstScript.IsScript()) {
			j2gort.CheckState(false)
		}
		if firstScript.HasChildren() && firstScript.GetFirstChild().IsModuleBody() {
			firstScript = firstScript.GetFirstChild()
		}
		if this_.instrumentOption == StaticGet_CompilerOptions_InstrumentOption_PRODUCTION() {
			this_.priv_CoverageInstrumentationPass_addProductionHeaderCode_6505(firstScript, this_.productionInstrumentationArrayName)
		} else {
			this_.priv_CoverageInstrumentationPass_addHeaderCode_6504(firstScript)
		}
	}
}

// Java method: CoverageInstrumentationPass.createConditionalObjectDecl(java.lang.String name, com.google.javascript.rhino.Node srcref)
func (this_ *CoverageInstrumentationPass) priv_CoverageInstrumentationPass_createConditionalObjectDecl_6507(name string, srcref pkg_rhino.Node) pkg_rhino.Node {
	_ = name
	_ = srcref
	var jscovData pkg_rhino.Node = pkg_rhino.Node(0)
	_ = jscovData
	switch (this_.instrumentOption).Ordinal() {
	case int32(2) /* InstrumentOption.BRANCH_ONLY */ :
		jscovData = pkg_rhino.IR_Objectlit([]pkg_rhino.Node{pkg_rhino.IR_QuotedStringKey("fileNames", pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{})), pkg_rhino.IR_QuotedStringKey("branchPresent", pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{})), pkg_rhino.IR_QuotedStringKey("branchesInLine", pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{})), pkg_rhino.IR_QuotedStringKey("branchesTaken", pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{}))})
	case int32(1) /* InstrumentOption.LINE_ONLY */ :
		jscovData = pkg_rhino.IR_Objectlit([]pkg_rhino.Node{pkg_rhino.IR_QuotedStringKey("fileNames", pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{})), pkg_rhino.IR_QuotedStringKey("instrumentedLines", pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{})), pkg_rhino.IR_QuotedStringKey("executedLines", pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{}))})
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Unexpected option: "+j2gort.StrVal(this_.instrumentOption))))
	}
	var var_ pkg_rhino.Node = pkg_rhino.IR_Var_2(pkg_rhino.IR_Name(name), pkg_rhino.IR_Or(pkg_rhino.IR_Getelem(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("window"), "top"), pkg_rhino.IR_String("__jscov")), pkg_rhino.IR_Assign(pkg_rhino.IR_Getelem(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("window"), "top"), pkg_rhino.IR_String("__jscov")), jscovData)))
	_ = var_
	return var_.SrcrefTreeIfMissing(srcref)
}

// Java enum: com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass.CoverageReach
type CoverageInstrumentationPass_CoverageReach uint8

const (
	enum_CoverageInstrumentationPass_CoverageReach_ALL         CoverageInstrumentationPass_CoverageReach = 1
	enum_CoverageInstrumentationPass_CoverageReach_CONDITIONAL CoverageInstrumentationPass_CoverageReach = 2
)

var enumBaseTable_CoverageInstrumentationPass_CoverageReach = [...]j2gort.EnumBase{
	{Name_: "ALL", Ordinal_: 0},
	{Name_: "CONDITIONAL", Ordinal_: 1},
}
var enumValues_CoverageInstrumentationPass_CoverageReach = []CoverageInstrumentationPass_CoverageReach{1, 2}

func StaticGet_CoverageInstrumentationPass_CoverageReach_ALL() CoverageInstrumentationPass_CoverageReach {
	return enum_CoverageInstrumentationPass_CoverageReach_ALL
}
func StaticGet_CoverageInstrumentationPass_CoverageReach_CONDITIONAL() CoverageInstrumentationPass_CoverageReach {
	return enum_CoverageInstrumentationPass_CoverageReach_CONDITIONAL
}

func CoverageInstrumentationPass_CoverageReach_Values() []CoverageInstrumentationPass_CoverageReach {
	return j2gort.CloneSlice(enumValues_CoverageInstrumentationPass_CoverageReach)
}

func CoverageInstrumentationPass_CoverageReach_ValueOf(name string) CoverageInstrumentationPass_CoverageReach {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CoverageInstrumentationPass_CoverageReach {
		if enumBaseTable_CoverageInstrumentationPass_CoverageReach[i_].Name_ == name {
			return CoverageInstrumentationPass_CoverageReach(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass.CoverageReach." + name))
}

func (this_ CoverageInstrumentationPass_CoverageReach) IsJavaNil_() bool                { return this_ == 0 }
func (CoverageInstrumentationPass_CoverageReach) HasIdentityEquals_()                   {}
func (CoverageInstrumentationPass_CoverageReach) Implements_j2gort_Ext_java_lang_Enum() {}
func (CoverageInstrumentationPass_CoverageReach) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ CoverageInstrumentationPass_CoverageReach) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CoverageInstrumentationPass_CoverageReach) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass$CoverageReach", "CoverageReach"
}
func (this_ CoverageInstrumentationPass_CoverageReach) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CoverageInstrumentationPass_CoverageReach[this_-1]
}
func (this_ CoverageInstrumentationPass_CoverageReach) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CoverageInstrumentationPass_CoverageReach) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CoverageInstrumentationPass_CoverageReach) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CoverageInstrumentationPass_CoverageReach[this_-1].Name_
}
func (this_ CoverageInstrumentationPass_CoverageReach) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CoverageInstrumentationPass_CoverageReach[this_-1].Name_
}
func (this_ CoverageInstrumentationPass_CoverageReach) Equals(other_ any) bool {
	o_, ok_ := other_.(CoverageInstrumentationPass_CoverageReach)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CoverageInstrumentationPass_CoverageReach) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CoverageInstrumentationPass_CoverageReach[this_-1].HashCode()
}
func (this_ CoverageInstrumentationPass_CoverageReach) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CoverageInstrumentationPass_CoverageReach[this_-1].CompareTo(other_)
}
func (this_ CoverageInstrumentationPass_CoverageReach) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CoverageInstrumentationPass_CoverageReach[this_-1].Name_)
}
func init() {
	vals_ := []any{CoverageInstrumentationPass_CoverageReach(1), CoverageInstrumentationPass_CoverageReach(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass.CoverageReach", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass$CoverageReach", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass.CoverageReach", "com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass$CoverageReach", func(v_ any) bool {
		return j2gort.InstanceOf[CoverageInstrumentationPass_CoverageReach](v_)
	})
}
