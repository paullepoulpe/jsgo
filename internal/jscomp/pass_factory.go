// Port of com/google/javascript/jscomp/PassFactory.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/AutoValue_PassFactory.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/CompilerPass.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/PassNames.java @ closure-compiler v20261021

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

/*
 * Copyright 2006 The Closure Compiler Authors.
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

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PassFactory creates JSCompiler passes and holds pass metadata (name,
// preconditions, and whether the pass runs in a fixed-point loop).
type PassFactory struct {
	name                string
	condition           func(options *CompilerOptions) bool
	preconditionCheck   func(options *CompilerOptions) *PassFactory_PreconditionResult
	runInFixedPointLoop bool
	internalFactory     func(compiler AbstractCompiler) CompilerPass
}

func passFactoryAlwaysEnabled(*CompilerOptions) bool { return true }

func passFactoryAlwaysValid(*CompilerOptions) *PassFactory_PreconditionResult {
	return PassFactory_PreconditionResult_SUCCESS
}

func makePassFactory(name string, factory func(AbstractCompiler) CompilerPass) *PassFactory {
	return &PassFactory{
		name:              name,
		condition:         passFactoryAlwaysEnabled,
		preconditionCheck: passFactoryAlwaysValid,
		internalFactory:   factory,
	}
}

func makeLoopPassFactory(name string, factory func(AbstractCompiler) CompilerPass) *PassFactory {
	return &PassFactory{
		name:                name,
		condition:           passFactoryAlwaysEnabled,
		preconditionCheck:   passFactoryAlwaysValid,
		runInFixedPointLoop: true,
		internalFactory:     factory,
	}
}

func makeCheckedPassFactory(
	name string,
	precond func(*CompilerOptions) *PassFactory_PreconditionResult,
	loop bool,
	factory func(AbstractCompiler) CompilerPass,
) *PassFactory {
	return &PassFactory{
		name:                name,
		condition:           passFactoryAlwaysEnabled,
		preconditionCheck:   precond,
		runInFixedPointLoop: loop,
		internalFactory:     factory,
	}
}

// GetName returns the name of the pass as it will appear in logs.
func (f *PassFactory) GetName() string { return f.name }

func (f *PassFactory) GetCondition() func(options *CompilerOptions) bool { return f.condition }

// PassFactory_PreconditionResult is the result of a pass precondition check.
type PassFactory_PreconditionResult struct {
	success bool
	message *string
}

func NewPassFactory_PreconditionResult(success bool, message *string) *PassFactory_PreconditionResult {
	return &PassFactory_PreconditionResult{success: success, message: message}
}

var PassFactory_PreconditionResult_SUCCESS = NewPassFactory_PreconditionResult(true, nil)

func (r *PassFactory_PreconditionResult) Success() bool    { return r.success }
func (r *PassFactory_PreconditionResult) Message() *string { return r.message }

// GetPreconditionCheck runs before this pass runs to validate preconditions.
func (f *PassFactory) GetPreconditionCheck() func(options *CompilerOptions) *PassFactory_PreconditionResult {
	return f.preconditionCheck
}

// IsRunInFixedPointLoop reports whether this factory runs inside a PhaseOptimizer loop.
func (f *PassFactory) IsRunInFixedPointLoop() bool { return f.runInFixedPointLoop }

// GetInternalFactory returns the underlying pass constructor function.
func (f *PassFactory) GetInternalFactory() func(compiler AbstractCompiler) CompilerPass {
	return f.internalFactory
}

func (f *PassFactory) String() string {
	return "PassFactory{" + f.name + "}"
}

func (f *PassFactory) Equals(o any) bool {
	that, ok := o.(*PassFactory)
	if !ok || that == nil {
		return false
	}
	if that == f {
		return true
	}
	return f.name == that.name &&
		f.runInFixedPointLoop == that.runInFixedPointLoop &&
		jutil.ObjectsEquals(f.condition, that.condition) &&
		jutil.ObjectsEquals(f.preconditionCheck, that.preconditionCheck) &&
		jutil.ObjectsEquals(f.internalFactory, that.internalFactory)
}

func (f *PassFactory) ToBuilder() *PassFactory_Builder {
	return &PassFactory_Builder{
		name:                f.name,
		condition:           f.condition,
		preconditionCheck:   f.preconditionCheck,
		runInFixedPointLoop: f.runInFixedPointLoop,
		internalFactory:     f.internalFactory,
	}
}

// PassFactory_Builder builds a PassFactory.
type PassFactory_Builder struct {
	name                string
	condition           func(options *CompilerOptions) bool
	preconditionCheck   func(options *CompilerOptions) *PassFactory_PreconditionResult
	runInFixedPointLoop bool
	internalFactory     func(compiler AbstractCompiler) CompilerPass
}

func PassFactory_Builder_() *PassFactory_Builder {
	return &PassFactory_Builder{
		condition:         passFactoryAlwaysEnabled,
		preconditionCheck: passFactoryAlwaysValid,
	}
}

func (b *PassFactory_Builder) SetName(name string) *PassFactory_Builder {
	b.name = name
	return b
}

func (b *PassFactory_Builder) SetCondition(condition func(options *CompilerOptions) bool) *PassFactory_Builder {
	if condition == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null condition"))
	}
	b.condition = condition
	return b
}

func (b *PassFactory_Builder) SetPreconditionCheck(preconditionCheck func(options *CompilerOptions) *PassFactory_PreconditionResult) *PassFactory_Builder {
	if preconditionCheck == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null preconditionCheck"))
	}
	b.preconditionCheck = preconditionCheck
	return b
}

func (b *PassFactory_Builder) SetRunInFixedPointLoop(runInFixedPointLoop bool) *PassFactory_Builder {
	b.runInFixedPointLoop = runInFixedPointLoop
	return b
}

func (b *PassFactory_Builder) SetInternalFactory(internalFactory func(compiler AbstractCompiler) CompilerPass) *PassFactory_Builder {
	if internalFactory == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null internalFactory"))
	}
	b.internalFactory = internalFactory
	return b
}

func (b *PassFactory_Builder) Build() *PassFactory {
	if b.condition == nil || b.preconditionCheck == nil || b.internalFactory == nil || b.name == "" {
		b.failBuild()
	}
	return &PassFactory{
		name:                b.name,
		condition:           b.condition,
		preconditionCheck:   b.preconditionCheck,
		runInFixedPointLoop: b.runInFixedPointLoop,
		internalFactory:     b.internalFactory,
	}
}

func (b *PassFactory_Builder) failBuild() {
	if b.condition == nil || b.preconditionCheck == nil || b.internalFactory == nil {
		panic(jutil.NewIllegalStateExceptionMsg("Missing required PassFactory properties"))
	}
	jutil.CheckState(b.name != "")
}

func (f *PassFactory) ValidatePreconditions(options *CompilerOptions) {
	canRun := f.preconditionCheck(options)
	if !canRun.success {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			jutil.StringFormat("Precondition for pass %s failed: %s", f.name, canRun.message)))
	}
}

// PassFactory_CreateEmptyPass creates a no-op pass that can only run once.
// Used to break up loops and mark pipeline boundaries.
func PassFactory_CreateEmptyPass(name string) *PassFactory {
	if name == "" {
		jutil.CheckState(false)
	}
	return makePassFactory(name, passFactoryEmptyInternal)
}

type passFactory_emptyPass struct{}

var passFactory_emptyPassInstance = &passFactory_emptyPass{}

func passFactoryEmptyInternal(AbstractCompiler) CompilerPass { return passFactory_emptyPassInstance }

func (*passFactory_emptyPass) Process(externs rhino.Node, root rhino.Node) {}

// Create creates a new compiler pass to be run.
func (f *PassFactory) Create(compiler AbstractCompiler) CompilerPass {
	return f.internalFactory(compiler)
}

// CompilerPass is the interface for classes that can compile JS.
//
// Class has single function "process", which is passed the root node of the
// parsed JS tree, as well as the root node of the external JS tree (used to
// provide a public API and prevent renaming of system functions).
//
// Use this class to support testing with BaseCompilerTest
type CompilerPass interface {
	// Process the JS with root node root. Can modify the contents of each
	// Node tree.
	//
	// externs is the top of external JS tree; root is the top of JS tree.
	//
	// Java: CompilerPass.process(Node, Node)
	Process(externs rhino.Node, root rhino.Node)
}

// CompilerPassFunc lets a Go closure stand for a Java lambda or method
// reference implementing CompilerPass, e.g. `(externs, root) -> ...`.
//
// Go-only: language-adapter — functional interface implemented by a lambda.
type CompilerPassFunc func(externs rhino.Node, root rhino.Node)

// Process calls f(externs, root).
func (f CompilerPassFunc) Process(externs rhino.Node, root rhino.Node) { f(externs, root) }

// PassNames: if the name of a pass is used in more than one place in the
// source, it's good to create a symbolic name here.
type PassNames struct{}

const (
	PassNames_AFTER_EARLY_OPTIMIZATION_LOOP                = "afterEarlyOptimizationLoop"
	PassNames_AFTER_MAIN_OPTIMIZATIONS                     = "afterMainOptimizations"
	PassNames_AFTER_STANDARD_CHECKS                        = "afterStandardChecks"
	PassNames_AMBIGUATE_PROPERTIES                         = "ambiguateProperties"
	PassNames_ANALYZER_CHECKS                              = "analyzerChecks"
	PassNames_ANGULAR_PASS                                 = "angularPass"
	PassNames_BEFORE_EARLY_OPTIMIZATION_LOOP               = "beforeEarlyOptimizationLoop"
	PassNames_BEFORE_EARLY_OPTIMIZATIONS_TRANSPILATION     = "beforeEarlyOptimizationsTranspilation"
	PassNames_BEFORE_STANDARD_OPTIMIZATIONS                = "beforeStandardOptimizations"
	PassNames_OBFUSCATION_PASS_MARKER                      = "obfuscationPassMarker"
	PassNames_BEFORE_MAIN_OPTIMIZATIONS                    = "beforeMainOptimizations"
	PassNames_BEFORE_TYPE_CHECKING                         = "beforeTypeChecking"
	PassNames_BEFORE_SERIALIZATION                         = "beforeSerialization"
	PassNames_BEFORE_VARIABLE_RENAMING                     = "beforeVariableRenaming"
	PassNames_BEFORE_EXTRACT_PROTOTYPE_MEMBER_DECLARATIONS = "beforeExtractPrototypeMemberDeclarations"
	PassNames_CHECK_CONFORMANCE                            = "checkConformance"
	PassNames_CHECK_REG_EXP                                = "checkRegExp"
	PassNames_CHECK_TYPES                                  = "checkTypes"
	PassNames_CHECK_VARIABLE_REFERENCES                    = "checkVariableReferences"
	PassNames_CHECK_VARS                                   = "checkVars"
	PassNames_COALESCE_VARIABLE_NAMES                      = "coalesceVariableNames"
	PassNames_COLLAPSE_ANONYMOUS_FUNCTIONS                 = "collapseAnonymousFunctions"
	PassNames_COLLAPSE_OBJECT_LITERALS                     = "collapseObjectLiterals"
	PassNames_COLLAPSE_PROPERTIES                          = "collapseProperties"
	PassNames_COLLAPSE_VARIABLE_DECLARATIONS               = "collapseVariableDeclarations"
	PassNames_CONVERT_TO_DOTTED_PROPERTIES                 = "convertToDottedProperties"
	PassNames_CREATE_MODULE_MAP                            = "createModuleMap"
	PassNames_CROSS_CHUNK_CODE_MOTION                      = "crossChunkCodeMotion"
	PassNames_CROSS_CHUNK_METHOD_MOTION                    = "crossChunkMethodMotion"
	PassNames_DEAD_ASSIGNMENT_ELIMINATION                  = "deadAssignmentsElimination"
	PassNames_DECLARED_GLOBAL_EXTERNS_ON_WINDOW            = "declaredGlobalExternsOnWindow"
	PassNames_DESERIALIZE_COMPILER_STATE                   = "deserializeCompilerState"
	PassNames_DEVIRTUALIZE_METHODS                         = "devirtualizeMethods"
	PassNames_DISAMBIGUATE_PROPERTIES                      = "disambiguateProperties"
	PassNames_ES6_NORMALIZE_CLASSES                        = "es6NormalizeClasses"
	PassNames_EXPORT_TEST_FUNCTIONS                        = "exportTestFunctions"
	PassNames_EXTERN_EXPORTS                               = "externExports"
	PassNames_EXTRACT_PROTOTYPE_MEMBER_DECLARATIONS        = "extractPrototypeMemberDeclarations"
	PassNames_FLOW_SENSITIVE_INLINE_VARIABLES              = "flowSensitiveInlineVariables"
	PassNames_GATHER_GETTERS_AND_SETTERS                   = "gatherGettersAndSetters"
	PassNames_GATHER_MODULE_METADATA                       = "gatherModuleMetadata"
	PassNames_GATHER_RAW_EXPORTS                           = "gatherRawExports"
	PassNames_GENERATE_EXPORTS                             = "generateExports"
	PassNames_INFER_CONSTS                                 = "inferConsts"
	PassNames_INFER_TYPES                                  = "inferTypes"
	PassNames_INLINE_FUNCTIONS                             = "inlineFunctions"
	PassNames_INLINE_PROPERTIES                            = "inlineProperties"
	PassNames_INLINE_TYPE_ALIASES                          = "inlineTypeAliases"
	PassNames_INLINE_VARIABLES                             = "inlineVariables"
	PassNames_LATE_PEEPHOLE_OPTIMIZATIONS                  = "latePeepholeOptimizations"
	PassNames_LINT_CHECKS                                  = "lintChecks"
	PassNames_MARK_UNNORMALIZED                            = "markUnnormalized"
	PassNames_NORMALIZE                                    = "normalize"
	PassNames_OPTIMIZATIONS_HALFWAY_POINT                  = "optimizationsHalfwayPoint"
	PassNames_OPTIMIZE_CALLS                               = "optimizeCalls"
	PassNames_PARSE_INPUTS                                 = "parseInputs"
	PassNames_PEEPHOLE_OPTIMIZATIONS                       = "peepholeOptimizations"
	PassNames_POST_NORMALIZE_PEEPHOLE                      = "postNormalizePeephole"
	PassNames_REWRITE_COMMON_JS_MODULES                    = "rewriteCommonJsModules"
	PassNames_REWRITE_SCRIPTS_TO_ES6_MODULES               = "rewriteScriptsToEs6Modules"
	PassNames_REMOVE_UNUSED_CODE                           = "removeUnusedCode"
	PassNames_REPLACE_ID_GENERATORS                        = "replaceIdGenerators"
	PassNames_REPLACE_MESSAGES                             = "replaceMessages"
	PassNames_RESOLVE_TYPES                                = "resolveTypes"
	PassNames_REWRITE_FUNCTION_EXPRESSIONS                 = "rewriteFunctionExpressions"
	PassNames_RENAME_PROPERTIES                            = "renameProperties"
	PassNames_BEFORE_RENAME_PROPERTIES                     = "beforeRenameProperties"
	PassNames_STRIP_SIDE_EFFECT_PROTECTION                 = "stripSideEffectProtection"
	PassNames_WIZ_PASS                                     = "wizPass"
)

// Go-only: measurement hook for the per-pass profile (cmd/jsgo,
// JSGO_PASS_LOG). It has no Java counterpart and never changes what the
// compiler does: when DebugPassLogHook is nil (the default) nothing below
// runs.

// DebugPassEvent describes the start or the end of one pass run (a
// PhaseOptimizer pass) or of one Compiler tracer section (parseInputs,
// toSource, ...).
type DebugPassEvent struct {
	Name   string
	Tracer bool // a Compiler.newTracer section rather than a PhaseOptimizer pass
	Start  bool
	// Loop is the 1-based ordinal of the fixed-point loop the pass runs in
	// (0 outside loops), Iter the 1-based iteration of that loop.
	Loop, Iter int
	// End events only: change reports made by the pass, and the number of
	// distinct change scopes (functions/scripts) it marked changed.
	Reports, ChangedScopes int
}

// DebugPassLogHook, when set, receives every DebugPassEvent.
var DebugPassLogHook func(ev DebugPassEvent)

const debugPassTimelineKey = "__debugPassLog__"

type debugPassChangeCounter struct{ n int }

func (h *debugPassChangeCounter) ReportChange() { h.n++ }

// debugPassStart reports the start of a PhaseOptimizer pass and returns the
// function that reports its end.
func (p *PhaseOptimizer) debugPassStart(name string) func() {
	ev := DebugPassEvent{Name: name, Start: true}
	if p.inLoop {
		ev.Loop, ev.Iter = p.debugLoopOrdinal, p.debugLoopIter
	}
	DebugPassLogHook(ev)
	counter := &debugPassChangeCounter{}
	p.changeTracker.AddChangeHandler(counter)
	p.changeTracker.changeTimeline.Mark(debugPassTimelineKey)
	return func() {
		p.changeTracker.RemoveChangeHandler(counter)
		scopes := p.changeTracker.changeTimeline.GetSince(debugPassTimelineKey)
		p.changeTracker.changeTimeline.Unmark(debugPassTimelineKey)
		ev.Start = false
		ev.Reports = counter.n
		ev.ChangedScopes = len(scopes)
		DebugPassLogHook(ev)
	}
}
