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
// Port of com/google/javascript/jscomp/PhaseOptimizer.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.PhaseOptimizer
type PhaseOptimizer struct {
	compiler                               AbstractCompiler                                     `gson:"compiler"`
	changeTracker                          *ChangeTracker                                       `gson:"changeTracker"`
	tracker                                *PerformanceTracker                                  `gson:"tracker"`
	passes                                 []CompilerPass                                       `gson:"passes"`
	validityCheck                          PassFactory                                          `gson:"validityCheck"`
	currentPass                            *PhaseOptimizer_NamedPass                            `gson:"currentPass"`
	lastRuns                               *j2gort.OrderedMap[*PhaseOptimizer_NamedPass, int32] `gson:"lastRuns"`
	changeVerifier                         *ChangeVerifier                                      `gson:"changeVerifier"`
	lastChange                             int32                                                `gson:"lastChange"`
	jsRoot                                 pkg_rhino.Node                                       `gson:"jsRoot"`
	optimizationLoopMaxIterations          int32                                                `gson:"optimizationLoopMaxIterations"`
	idHash_                                atomic.Int32
	inLoop                                 bool `gson:"inLoop"`
	printAstHashcodes                      bool `gson:"printAstHashcodes"`
	useSizeHeuristicToStopOptimizationLoop bool `gson:"useSizeHeuristicToStopOptimizationLoop"`
}

func (this_ *PhaseOptimizer) BindSelf_(s any) {
	_ = s
}
func (this_ *PhaseOptimizer) IsJavaNil_() bool { return this_ == nil }
func (*PhaseOptimizer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PhaseOptimizer", "PhaseOptimizer"
}
func (this_ *PhaseOptimizer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PhaseOptimizer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PhaseOptimizer", "com.google.javascript.jscomp.PhaseOptimizer", func(v_ any) bool {
		return j2gort.InstanceOf[*PhaseOptimizer](v_)
	})
}
func (*PhaseOptimizer) Implements_jscomp_CompilerPass() {}

var static_PhaseOptimizer_logger *j2gort.Ext_java_util_logging_Logger

func StaticGet_PhaseOptimizer_Logger() *j2gort.Ext_java_util_logging_Logger {
	clinit_PhaseOptimizer()
	return static_PhaseOptimizer_logger
}
func StaticSet_PhaseOptimizer_Logger(v *j2gort.Ext_java_util_logging_Logger) {
	clinit_PhaseOptimizer()
	static_PhaseOptimizer_logger = v
}
func StaticGet_PhaseOptimizer_START_TIME() int32 { return int32(0) }

var static_PhaseOptimizer_CODE_REMOVING_PASSES j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_PhaseOptimizer_CODE_REMOVING_PASSES() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_PhaseOptimizer()
	return static_PhaseOptimizer_CODE_REMOVING_PASSES
}
func StaticSet_PhaseOptimizer_CODE_REMOVING_PASSES(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_PhaseOptimizer()
	static_PhaseOptimizer_CODE_REMOVING_PASSES = v
}
func StaticGet_PhaseOptimizer_MAX_LOOPS() int32 { return int32(100) }
func StaticGet_PhaseOptimizer_OPTIMIZE_LOOP_ERROR() string {
	return "Fixed point loop exceeded the maximum number of iterations."
}

var clinitDone_PhaseOptimizer uint32
var clinitState_PhaseOptimizer j2gort.ClinitState

func clinit_PhaseOptimizer() {
	if atomic.LoadUint32(&clinitDone_PhaseOptimizer) != 0 {
		return
	}
	clinitSlow_PhaseOptimizer()
}

func clinitSlow_PhaseOptimizer() {
	if !j2gort.EnterClinit(&clinitState_PhaseOptimizer, &clinitDone_PhaseOptimizer) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PhaseOptimizer, &clinitDone_PhaseOptimizer)
	static_PhaseOptimizer_logger = j2gort.Ext_java_util_logging_Logger_GetLogger_1(j2gort.CheckNotNullPtr(j2gort.ClassLiteral("com.google.javascript.jscomp.PhaseOptimizer", "PhaseOptimizer")).GetName())
	static_PhaseOptimizer_CODE_REMOVING_PASSES = j2gort.ImmListOfStr1("peepholeOptimizations")
}

func (this_ *PhaseOptimizer) initFields_PhaseOptimizer() {
	_ = this_

}

// Java constructor: PhaseOptimizer(com.google.javascript.jscomp.AbstractCompiler comp, com.google.javascript.jscomp.PerformanceTracker tracker)
func NewPhaseOptimizer(comp AbstractCompiler, tracker *PerformanceTracker) *PhaseOptimizer {
	clinit_PhaseOptimizer()
	this_ := &PhaseOptimizer{}
	InitPhaseOptimizer(this_, comp, tracker)
	return this_
}

func InitPhaseOptimizer(this_ *PhaseOptimizer, comp AbstractCompiler, tracker *PerformanceTracker) {
	_ = comp
	_ = tracker
	this_.initFields_PhaseOptimizer()
	this_.compiler = comp
	this_.jsRoot = comp.GetJsRoot()
	this_.tracker = tracker
	this_.changeTracker = comp.GetChangeTracker()
	this_.passes = []CompilerPass{}
	this_.inLoop = false
	this_.lastChange = int32(0)
	this_.useSizeHeuristicToStopOptimizationLoop = comp.GetOptions().ShouldUseSizeHeuristicToStopOptimizationLoop()
	var maxIterations int32 = comp.GetOptions().GetMaxOptimizationLoopIterations()
	_ = maxIterations
	if (maxIterations > int32(0)) && (maxIterations <= int32(100)) {
		this_.optimizationLoopMaxIterations = maxIterations
	} else {
		this_.optimizationLoopMaxIterations = int32(100)
	}
}

// Java method: PhaseOptimizer.consume(java.util.List<com.google.javascript.jscomp.PassFactory> factories)
func (this_ *PhaseOptimizer) Consume(factories j2gort.Ext_java_util_List) {
	_ = factories
	var val_1 PhaseOptimizer_Loop
	val_1.outer_ = this_
	InitPhaseOptimizer_Loop(&val_1)
	var currentLoop *PhaseOptimizer_Loop = &val_1
	_ = currentLoop
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(factories))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			factory := j2gort.AsIface[PassFactory](rawElem_3)
			_ = factory
			if factory.IsRunInFixedPointLoop() {
				currentLoop.AddLoopedPass(factory)
			} else {
				if currentLoop.IsPopulated() {
					this_.passes = append(this_.passes, Upcast_CompilerPass(currentLoop))
					currentLoop = NewPhaseOptimizer_Loop(this_)
				}
				this_.AddOneTimePass(factory)
			}
		}
	}
	if currentLoop.IsPopulated() {
		this_.passes = append(this_.passes, Upcast_CompilerPass(currentLoop))
	}
}

// Java method: PhaseOptimizer.addOneTimePass(com.google.javascript.jscomp.PassFactory factory)
func (this_ *PhaseOptimizer) AddOneTimePass(factory PassFactory) {
	_ = factory
	this_.passes = append(this_.passes, NewPhaseOptimizer_NamedPass(this_, factory))
}

// Java method: PhaseOptimizer.addFixedPointLoop()
func (this_ *PhaseOptimizer) AddFixedPointLoop() *PhaseOptimizer_Loop {
	var val_4 PhaseOptimizer_Loop
	val_4.outer_ = this_
	InitPhaseOptimizer_Loop(&val_4)
	var loop *PhaseOptimizer_Loop = &val_4
	_ = loop
	this_.passes = append(this_.passes, Upcast_CompilerPass(loop))
	return loop
}

// Java method: PhaseOptimizer.setValidityCheck(com.google.javascript.jscomp.PassFactory validityCheck)
func (this_ *PhaseOptimizer) SetValidityCheck(validityCheck PassFactory) {
	_ = validityCheck
	this_.validityCheck = validityCheck
	this_.changeVerifier = NewChangeVerifier(this_.compiler).Snapshot(this_.jsRoot)
}

// Java method: PhaseOptimizer.setPrintAstHashcodes(boolean printAstHashcodes)
func (this_ *PhaseOptimizer) SetPrintAstHashcodes(printAstHashcodes bool) {
	_ = printAstHashcodes
	this_.printAstHashcodes = printAstHashcodes
}

// Java method: PhaseOptimizer.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *PhaseOptimizer) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	for _, pass := range this_.passes {
		_ = pass
		if Platform_IsThreadInterrupted() {
			panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", j2gort.NewExt_java_lang_InterruptedException())))
		}
		pass.Process_2(externs, root)
		if this_.priv_PhaseOptimizer_hasHaltingErrors_4354() {
			return
		}
	}

}

// Java method: PhaseOptimizer.maybePrintAstHashcodes(java.lang.String passName, com.google.javascript.rhino.Node root)
func (this_ *PhaseOptimizer) priv_PhaseOptimizer_maybePrintAstHashcodes_4352(passName string, root pkg_rhino.Node) {
	_ = passName
	_ = root
	if this_.printAstHashcodes {
		var hashCodeMsg string = ("AST hashCode after " + j2gort.StrValString(passName) + ": " + j2gort.StrValInt32(javalang.StringHashCode(j2gort.CheckNotNullString(this_.compiler.ToSource_1_Node(root)))))
		_ = hashCodeMsg
		StaticGet_PhaseOptimizer_Logger().Info_1_string(hashCodeMsg)
	}
}

// Java method: PhaseOptimizer.maybeRunValidityCheck(java.lang.String passName, com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *PhaseOptimizer) priv_PhaseOptimizer_maybeRunValidityCheck_4353(passName string, externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = passName
	_ = externs
	_ = root
	if this_.validityCheck == nil || this_.validityCheck.IsJavaNil_() {
		return
	}
	func() {
		defer func() {
			if rec_7 := recover(); rec_7 != nil {
				jex_8 := javalang.RecoverException(rec_7)
				_ = jex_8
				if j2gort.InstanceOf[javalang.IRuntimeException](jex_8) {
					e := j2gort.AsIface[javalang.IRuntimeException](jex_8)
					_ = e
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Validity checks failed for pass: "+j2gort.StrValString(passName), e)))
				} else {
					panic(rec_7)
				}
			}
		}()
		this_.validityCheck.GetPassFactoryBase().Create(this_.compiler).Process_2(externs, root)
		this_.changeVerifier.CheckRecordedChanges_2(passName, this_.jsRoot)
	}()
}

// Java method: PhaseOptimizer.hasHaltingErrors()
func (this_ *PhaseOptimizer) priv_PhaseOptimizer_hasHaltingErrors_4354() bool {
	return this_.compiler.HasHaltingErrors()
}

// Java method: PhaseOptimizer.hasScopeChanged(com.google.javascript.rhino.Node n)
func (this_ *PhaseOptimizer) HasScopeChanged(n pkg_rhino.Node) bool {
	_ = n
	if !(this_.inLoop) {
		return true
	}
	var timeOfLastRun int32
iifeSwitch_1:
	switch {
	default:
		if v_, ok_ := this_.lastRuns.GetOk(this_.currentPass); ok_ {
			timeOfLastRun = v_
			break iifeSwitch_1
		}
		j2gort.ThrowNPE()
		timeOfLastRun = 0
		break iifeSwitch_1
	}

	_ = timeOfLastRun
	return ((timeOfLastRun == int32(0)) || (n.GetChangeTime() > timeOfLastRun))
}

// Java enum: com.google.javascript.jscomp.PhaseOptimizer.State
type PhaseOptimizer_State uint8

const (
	enum_PhaseOptimizer_State_RUN_PASSES_NOT_RUN_IN_PREV_ITER          PhaseOptimizer_State = 1
	enum_PhaseOptimizer_State_RUN_PASSES_THAT_CHANGED_STH_IN_PREV_ITER PhaseOptimizer_State = 2
)

var enumBaseTable_PhaseOptimizer_State = [...]j2gort.EnumBase{
	{Name_: "RUN_PASSES_NOT_RUN_IN_PREV_ITER", Ordinal_: 0},
	{Name_: "RUN_PASSES_THAT_CHANGED_STH_IN_PREV_ITER", Ordinal_: 1},
}
var enumValues_PhaseOptimizer_State = []PhaseOptimizer_State{1, 2}

func StaticGet_PhaseOptimizer_State_RUN_PASSES_NOT_RUN_IN_PREV_ITER() PhaseOptimizer_State {
	return enum_PhaseOptimizer_State_RUN_PASSES_NOT_RUN_IN_PREV_ITER
}
func StaticGet_PhaseOptimizer_State_RUN_PASSES_THAT_CHANGED_STH_IN_PREV_ITER() PhaseOptimizer_State {
	return enum_PhaseOptimizer_State_RUN_PASSES_THAT_CHANGED_STH_IN_PREV_ITER
}

func PhaseOptimizer_State_Values() []PhaseOptimizer_State {
	return j2gort.CloneSlice(enumValues_PhaseOptimizer_State)
}

func PhaseOptimizer_State_ValueOf(name string) PhaseOptimizer_State {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_PhaseOptimizer_State {
		if enumBaseTable_PhaseOptimizer_State[i_].Name_ == name {
			return PhaseOptimizer_State(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.PhaseOptimizer.State." + name))
}

func (this_ PhaseOptimizer_State) IsJavaNil_() bool                              { return this_ == 0 }
func (PhaseOptimizer_State) HasIdentityEquals_()                                 {}
func (PhaseOptimizer_State) Implements_j2gort_Ext_java_lang_Enum()               {}
func (PhaseOptimizer_State) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ PhaseOptimizer_State) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (PhaseOptimizer_State) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PhaseOptimizer$State", "State"
}
func (this_ PhaseOptimizer_State) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_PhaseOptimizer_State[this_-1]
}
func (this_ PhaseOptimizer_State) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ PhaseOptimizer_State) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ PhaseOptimizer_State) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_PhaseOptimizer_State[this_-1].Name_
}
func (this_ PhaseOptimizer_State) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_PhaseOptimizer_State[this_-1].Name_
}
func (this_ PhaseOptimizer_State) Equals(other_ any) bool {
	o_, ok_ := other_.(PhaseOptimizer_State)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ PhaseOptimizer_State) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_PhaseOptimizer_State[this_-1].HashCode()
}
func (this_ PhaseOptimizer_State) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_PhaseOptimizer_State[this_-1].CompareTo(other_)
}
func (this_ PhaseOptimizer_State) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_PhaseOptimizer_State[this_-1].Name_)
}
func init() {
	vals_ := []any{PhaseOptimizer_State(1), PhaseOptimizer_State(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.PhaseOptimizer.State", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.PhaseOptimizer$State", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PhaseOptimizer.State", "com.google.javascript.jscomp.PhaseOptimizer$State", func(v_ any) bool {
		return j2gort.InstanceOf[PhaseOptimizer_State](v_)
	})
}

// Java class: com.google.javascript.jscomp.PhaseOptimizer.NamedPass
type PhaseOptimizer_NamedPass struct {
	outer_  *PhaseOptimizer
	name    string      `gson:"name"`
	factory PassFactory `gson:"factory"`
	tracer  *Tracer     `gson:"tracer"`
	idHash_ atomic.Int32
}

func (this_ *PhaseOptimizer_NamedPass) BindSelf_(s any) {
	_ = s
}
func (this_ *PhaseOptimizer_NamedPass) IsJavaNil_() bool { return this_ == nil }
func (*PhaseOptimizer_NamedPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PhaseOptimizer$NamedPass", "NamedPass"
}
func (this_ *PhaseOptimizer_NamedPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PhaseOptimizer_NamedPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PhaseOptimizer.NamedPass", "com.google.javascript.jscomp.PhaseOptimizer$NamedPass", func(v_ any) bool {
		return j2gort.InstanceOf[*PhaseOptimizer_NamedPass](v_)
	})
}
func (this_ *PhaseOptimizer_NamedPass) SetOuter_(o *PhaseOptimizer) { this_.outer_ = o }
func (this_ *PhaseOptimizer_NamedPass) GetOuter_() *PhaseOptimizer  { return this_.outer_ }
func (*PhaseOptimizer_NamedPass) Implements_jscomp_CompilerPass()   {}

func (this_ *PhaseOptimizer_NamedPass) initFields_PhaseOptimizer_NamedPass() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: NamedPass(com.google.javascript.jscomp.PassFactory factory)
func NewPhaseOptimizer_NamedPass(outer_ *PhaseOptimizer, factory PassFactory) *PhaseOptimizer_NamedPass {
	this_ := &PhaseOptimizer_NamedPass{}
	this_.outer_ = outer_
	InitPhaseOptimizer_NamedPass(this_, factory)
	return this_
}

func InitPhaseOptimizer_NamedPass(this_ *PhaseOptimizer_NamedPass, factory PassFactory) {
	_ = factory
	this_.initFields_PhaseOptimizer_NamedPass()
	this_.name = factory.GetName()
	this_.factory = factory
}

// Java method: NamedPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *PhaseOptimizer_NamedPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	StaticGet_PhaseOptimizer_Logger().Fine_1_string(("Running pass " + j2gort.StrValString(this_.name)))
	this_.factory.GetPassFactoryBase().ValidatePreconditions(this_.outer_.compiler.GetOptions())
	if this_.outer_.validityCheck != nil && !this_.outer_.validityCheck.IsJavaNil_() {
		this_.outer_.changeVerifier = NewChangeVerifier(this_.outer_.compiler).Snapshot(this_.outer_.jsRoot)
	}
	if this_.outer_.tracker != nil {
		this_.outer_.tracker.RecordPassStart(this_.name, !(this_.factory.IsRunInFixedPointLoop()))
	}
	this_.tracer = NewTracer_2("Compiler", this_.name)
	this_.outer_.compiler.BeforePass(this_.name)
	this_.factory.GetPassFactoryBase().Create(this_.outer_.compiler).Process_2(externs, root)
	this_.outer_.compiler.AfterPass(this_.name)
	func() {
		defer func() {
			if rec_9 := recover(); rec_9 != nil {
				jex_10 := javalang.RecoverException(rec_9)
				_ = jex_10
				if j2gort.InstanceOf[javalang.IIllegalStateException](jex_10) {
					e := j2gort.AsIface[javalang.IIllegalStateException](jex_10)
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Validity check failed for "+j2gort.StrValString(this_.name), e)))
				} else {
					panic(rec_9)
				}
			}
		}()
		var traceRuntime int64 = this_.tracer.Stop_0()
		_ = traceRuntime
		if this_.outer_.tracker != nil {
			this_.outer_.tracker.RecordPassStop(this_.name, traceRuntime)
		}
		this_.outer_.priv_PhaseOptimizer_maybePrintAstHashcodes_4352(this_.name, root)
		this_.outer_.priv_PhaseOptimizer_maybeRunValidityCheck_4353(this_.name, externs, root)
	}()
}

// Java method: NamedPass.toString()
func (this_ *PhaseOptimizer_NamedPass) String() string {
	return ("pass: " + j2gort.StrValString(this_.name))
}

// Java class: com.google.javascript.jscomp.PhaseOptimizer.ScopedChangeHandler
type PhaseOptimizer_ScopedChangeHandler struct {
	outer_              *PhaseOptimizer
	lastCodeChangeQuery int32 `gson:"lastCodeChangeQuery"`
	idHash_             atomic.Int32
}

func (this_ *PhaseOptimizer_ScopedChangeHandler) BindSelf_(s any) {
	_ = s
}
func (this_ *PhaseOptimizer_ScopedChangeHandler) IsJavaNil_() bool { return this_ == nil }
func (*PhaseOptimizer_ScopedChangeHandler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PhaseOptimizer$ScopedChangeHandler", "ScopedChangeHandler"
}
func (this_ *PhaseOptimizer_ScopedChangeHandler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PhaseOptimizer_ScopedChangeHandler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PhaseOptimizer.ScopedChangeHandler", "com.google.javascript.jscomp.PhaseOptimizer$ScopedChangeHandler", func(v_ any) bool {
		return j2gort.InstanceOf[*PhaseOptimizer_ScopedChangeHandler](v_)
	})
}
func (this_ *PhaseOptimizer_ScopedChangeHandler) SetOuter_(o *PhaseOptimizer)    { this_.outer_ = o }
func (this_ *PhaseOptimizer_ScopedChangeHandler) GetOuter_() *PhaseOptimizer     { return this_.outer_ }
func (*PhaseOptimizer_ScopedChangeHandler) Implements_jscomp_CodeChangeHandler() {}

func (this_ *PhaseOptimizer_ScopedChangeHandler) initFields_PhaseOptimizer_ScopedChangeHandler() {
	_ = this_
}

// Java constructor: ScopedChangeHandler()
func NewPhaseOptimizer_ScopedChangeHandler(outer_ *PhaseOptimizer) *PhaseOptimizer_ScopedChangeHandler {
	this_ := &PhaseOptimizer_ScopedChangeHandler{}
	this_.outer_ = outer_
	InitPhaseOptimizer_ScopedChangeHandler(this_)
	return this_
}

func InitPhaseOptimizer_ScopedChangeHandler(this_ *PhaseOptimizer_ScopedChangeHandler) {
	this_.initFields_PhaseOptimizer_ScopedChangeHandler()
	this_.lastCodeChangeQuery = this_.outer_.changeTracker.GetChangeStamp()
}

// Java method: ScopedChangeHandler.reportChange()
func (this_ *PhaseOptimizer_ScopedChangeHandler) ReportChange() {
	this_.outer_.lastChange = this_.outer_.changeTracker.GetChangeStamp()
}

// Java method: ScopedChangeHandler.hasCodeChangedSinceLastCall()
func (this_ *PhaseOptimizer_ScopedChangeHandler) priv_ScopedChangeHandler_hasCodeChangedSinceLastCall_4357() bool {
	var result bool = (this_.outer_.lastChange > this_.lastCodeChangeQuery)
	_ = result
	this_.lastCodeChangeQuery = this_.outer_.changeTracker.GetChangeStamp()
	this_.outer_.changeTracker.IncrementChangeStamp()
	return result
}

// Java class: com.google.javascript.jscomp.PhaseOptimizer.Loop
type PhaseOptimizer_Loop struct {
	outer_                          *PhaseOptimizer
	myPasses                        []*PhaseOptimizer_NamedPass         `gson:"myPasses"`
	myNames                         j2gort.OrderedSet[string]           `gson:"myNames"`
	scopeHandler                    *PhaseOptimizer_ScopedChangeHandler `gson:"scopeHandler"`
	howmanyIterationsUnderThreshold int32                               `gson:"howmanyIterationsUnderThreshold"`
	idHash_                         atomic.Int32
	isCodeRemovalLoop               bool `gson:"isCodeRemovalLoop"`
}

func (this_ *PhaseOptimizer_Loop) BindSelf_(s any) {
	_ = s
}
func (this_ *PhaseOptimizer_Loop) IsJavaNil_() bool { return this_ == nil }
func (*PhaseOptimizer_Loop) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PhaseOptimizer$Loop", "Loop"
}
func (this_ *PhaseOptimizer_Loop) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PhaseOptimizer_Loop) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PhaseOptimizer.Loop", "com.google.javascript.jscomp.PhaseOptimizer$Loop", func(v_ any) bool {
		return j2gort.InstanceOf[*PhaseOptimizer_Loop](v_)
	})
}
func (this_ *PhaseOptimizer_Loop) SetOuter_(o *PhaseOptimizer) { this_.outer_ = o }
func (this_ *PhaseOptimizer_Loop) GetOuter_() *PhaseOptimizer  { return this_.outer_ }
func (*PhaseOptimizer_Loop) Implements_jscomp_CompilerPass()   {}

func (this_ *PhaseOptimizer_Loop) initFields_PhaseOptimizer_Loop() {
	_ = this_
	this_.myPasses = []*PhaseOptimizer_NamedPass{}
	this_.myNames = j2gort.NewOrderedSet[string]()

}

// Java constructor: Loop()
func NewPhaseOptimizer_Loop(outer_ *PhaseOptimizer) *PhaseOptimizer_Loop {
	this_ := &PhaseOptimizer_Loop{}
	this_.outer_ = outer_
	InitPhaseOptimizer_Loop(this_)
	return this_
}

func InitPhaseOptimizer_Loop(this_ *PhaseOptimizer_Loop) {
	this_.initFields_PhaseOptimizer_Loop()
}

// Java method: Loop.addLoopedPass(com.google.javascript.jscomp.PassFactory factory)
func (this_ *PhaseOptimizer_Loop) AddLoopedPass(factory PassFactory) {
	_ = factory
	var name string = factory.GetName()
	_ = name
	if !(!(this_.myNames.Contains(name))) {
		j2gort.CheckArgumentFmt(false, "Already a pass with name '%s' in this loop", j2gort.BoxString(name))
	}
	this_.myNames.Add(name)
	this_.myPasses = append(this_.myPasses, NewPhaseOptimizer_NamedPass(this_.outer_, factory))
}

// Java method: Loop.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *PhaseOptimizer_Loop) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(!(this_.outer_.inLoop)) {
		j2gort.CheckStateStr(false, "Nested loops are forbidden")
	}
	this_.outer_.inLoop = true
	this_.isCodeRemovalLoop = this_.priv_Loop_isCodeRemovalLoop_4361()
	this_.scopeHandler = NewPhaseOptimizer_ScopedChangeHandler(this_.outer_)
	this_.outer_.changeTracker.AddChangeHandler(ToIface_CodeChangeHandler(this_.scopeHandler))
	this_.outer_.lastRuns = &j2gort.OrderedMap[*PhaseOptimizer_NamedPass, int32]{}
	for _, pass := range this_.myPasses {
		_ = pass
		this_.outer_.lastRuns.Put(pass, int32(0))
	}

	var madeChanges j2gort.LocalOrderedSet[*PhaseOptimizer_NamedPass] = j2gort.NewLocalOrderedSet[*PhaseOptimizer_NamedPass]()
	_ = madeChanges
	var runInPrevIter j2gort.LocalOrderedSet[*PhaseOptimizer_NamedPass] = j2gort.NewLocalOrderedSet[*PhaseOptimizer_NamedPass]()
	_ = runInPrevIter
	var didNotMakeChanges j2gort.LocalOrderedSet[*PhaseOptimizer_NamedPass] = j2gort.NewLocalOrderedSet[*PhaseOptimizer_NamedPass]()
	_ = didNotMakeChanges
	var state PhaseOptimizer_State = StaticGet_PhaseOptimizer_State_RUN_PASSES_NOT_RUN_IN_PREV_ITER()
	_ = state
	var lastIterMadeChanges bool = false
	_ = lastIterMadeChanges
	var count int32 = int32(1)
	_ = count
	var astSize int32 = NodeUtil_CountAstSize(root)
	_ = astSize
	var previousAstSize int32 = astSize
	_ = previousAstSize
	defer func() {
		this_.outer_.inLoop = false
		this_.outer_.changeTracker.RemoveChangeHandler(ToIface_CodeChangeHandler(this_.scopeHandler))
	}()
	{
		for {
			if (count > this_.outer_.optimizationLoopMaxIterations) && this_.isCodeRemovalLoop {
				return
			}
			if count > int32(100) {
				this_.outer_.compiler.ThrowInternalError("Fixed point loop exceeded the maximum number of iterations.", nil)
			}
			count = count + 1
			lastIterMadeChanges = false
			for _, pass_2 := range this_.myPasses {
				_ = pass_2
				if (((state == StaticGet_PhaseOptimizer_State_RUN_PASSES_NOT_RUN_IN_PREV_ITER()) && !(runInPrevIter.Contains(pass_2))) && !(didNotMakeChanges.Contains(pass_2))) || ((state == StaticGet_PhaseOptimizer_State_RUN_PASSES_THAT_CHANGED_STH_IN_PREV_ITER()) && madeChanges.Contains(pass_2)) {
					this_.outer_.changeTracker.IncrementChangeStamp()
					this_.outer_.currentPass = pass_2
					pass_2.Process_2(externs, root)
					runInPrevIter.Add(pass_2)
					this_.outer_.lastRuns.Put(pass_2, this_.outer_.changeTracker.GetChangeStamp())
					if this_.outer_.priv_PhaseOptimizer_hasHaltingErrors_4354() {
						return
					} else if this_.scopeHandler.priv_ScopedChangeHandler_hasCodeChangedSinceLastCall_4357() {
						madeChanges.Add(pass_2)
						didNotMakeChanges.Clear()
						lastIterMadeChanges = true
					} else {
						madeChanges.Remove(pass_2)
						didNotMakeChanges.Add(pass_2)
					}
				} else {
					runInPrevIter.Remove(pass_2)
				}
			}

			previousAstSize = astSize
			astSize = NodeUtil_CountAstSize(root)
			if state == StaticGet_PhaseOptimizer_State_RUN_PASSES_NOT_RUN_IN_PREV_ITER() {
				if lastIterMadeChanges && this_.priv_Loop_isAstSufficientlyChanging_4359(previousAstSize, astSize) {
					state = StaticGet_PhaseOptimizer_State_RUN_PASSES_THAT_CHANGED_STH_IN_PREV_ITER()
				} else {
					return
				}
			} else {
				if !(state == StaticGet_PhaseOptimizer_State_RUN_PASSES_THAT_CHANGED_STH_IN_PREV_ITER()) {
					j2gort.CheckState(false)
				}
				if !(lastIterMadeChanges) || !(this_.priv_Loop_isAstSufficientlyChanging_4359(previousAstSize, astSize)) {
					state = StaticGet_PhaseOptimizer_State_RUN_PASSES_NOT_RUN_IN_PREV_ITER()
				}
			}
		}
	}
}

// Java method: Loop.isAstSufficientlyChanging(int oldAstSize, int newAstSize)
func (this_ *PhaseOptimizer_Loop) priv_Loop_isAstSufficientlyChanging_4359(oldAstSize int32, newAstSize int32) bool {
	_ = oldAstSize
	_ = newAstSize
	if this_.outer_.useSizeHeuristicToStopOptimizationLoop && this_.isCodeRemovalLoop {
		var percentChange float32 = (float32(int32(100)) * (float32(func() int32 {
			ab_15 := (newAstSize - oldAstSize)
			if ab_15 < 0 {
				return -ab_15
			}
			return ab_15
		}()) / float32(oldAstSize)))
		_ = percentChange
		if float64(percentChange) < float64(0.05) {
			this_.howmanyIterationsUnderThreshold = this_.howmanyIterationsUnderThreshold + 1
		} else {
			this_.howmanyIterationsUnderThreshold = int32(0)
		}
		return (this_.howmanyIterationsUnderThreshold < int32(2))
	}
	return true
}

// Java method: Loop.isPopulated()
func (this_ *PhaseOptimizer_Loop) IsPopulated() bool {
	return len(this_.myPasses) != 0
}

// Java method: Loop.isCodeRemovalLoop()
func (this_ *PhaseOptimizer_Loop) priv_Loop_isCodeRemovalLoop_4361() bool {
	for _, pass := range this_.myPasses {
		_ = pass
		if j2gort.SetContainsStr(StaticGet_PhaseOptimizer_CODE_REMOVING_PASSES(), pass.name) {
			return true
		}
	}

	return false
}

// Java class: com.google.javascript.jscomp.PhaseOptimizer.ProgressRange
type PhaseOptimizer_ProgressRange struct {
	initialValue float64 `gson:"initialValue"`
	maxValue     float64 `gson:"maxValue"`
	idHash_      atomic.Int32
}

func (this_ *PhaseOptimizer_ProgressRange) BindSelf_(s any) {
	_ = s
}
func (this_ *PhaseOptimizer_ProgressRange) IsJavaNil_() bool { return this_ == nil }
func (*PhaseOptimizer_ProgressRange) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PhaseOptimizer$ProgressRange", "ProgressRange"
}
func (this_ *PhaseOptimizer_ProgressRange) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PhaseOptimizer_ProgressRange) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PhaseOptimizer.ProgressRange", "com.google.javascript.jscomp.PhaseOptimizer$ProgressRange", func(v_ any) bool {
		return j2gort.InstanceOf[*PhaseOptimizer_ProgressRange](v_)
	})
}

func (this_ *PhaseOptimizer_ProgressRange) initFields_PhaseOptimizer_ProgressRange() {
	_ = this_
}

// Java constructor: ProgressRange(double initialValue, double maxValue)
func NewPhaseOptimizer_ProgressRange(initialValue float64, maxValue float64) *PhaseOptimizer_ProgressRange {
	this_ := &PhaseOptimizer_ProgressRange{}
	InitPhaseOptimizer_ProgressRange(this_, initialValue, maxValue)
	return this_
}

func InitPhaseOptimizer_ProgressRange(this_ *PhaseOptimizer_ProgressRange, initialValue float64, maxValue float64) {
	_ = initialValue
	_ = maxValue
	this_.initFields_PhaseOptimizer_ProgressRange()
	this_.initialValue = initialValue
	this_.maxValue = maxValue
}
