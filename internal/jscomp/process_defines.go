// Port of com/google/javascript/jscomp/ProcessDefines.java @ closure-compiler v20261021

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) where it reads
// shapes: the default value of goog.define('name', value) (its second
// argument, processDefines_googDefineDefault, where Java reads child 2 of
// the call), the operands of the operators a define value may use, and the
// CLOSURE_DEFINES declaration and assignments (shared by the Go-only walk
// and Visit). The goog.define checks stay on rhino: they pass arguments
// that may be missing to Java's verifyNotNull and verifyIsLast.

// ProcessDefines processes variables annotated as `@define`. A define is a
// special constant that may be overridden by later files and manipulated by
// the compiler, much like C preprocessor `#define`s.
type ProcessDefines struct {
	compiler                   AbstractCompiler
	registry                   *rhino.JSTypeRegistry
	replacementValuesFromFlags OrderedStringMap[rhino.Node]
	mode                       ProcessDefines_Mode
	namespaceSupplier          func() *GlobalNamespace
	recognizeClosureDefines    bool
	enableZonesDefineName      string
	hasEnableZonesDefineName   bool
	zoneInputPattern           *jutil.Pattern
	unknownDefinesToIgnore     map[string]struct{}

	knownDefineJsdocs    map[*rhino.JSDocInfo]struct{}
	knownGoogDefineCalls map[rhino.Node]struct{}
	defineByDefineName   map[string]*ProcessDefines_Define
	defines              []*ProcessDefines_Define
	// from var CLOSURE_DEFINES = {
	replacementValuesFromClosureDefines map[string]rhino.Node
	closureDefineOrder                  []string
	validDefineValueExpressions         map[rhino.Node]struct{}
	hasZoneInput                        bool

	namespace *GlobalNamespace
}

var _ CompilerPass = (*ProcessDefines)(nil)

// Defines in this set will not be flagged with "unknown define" warnings.
// There are flags that always set these defines, even when they might not be
// in the binary.
var processDefinesKnownDefineList = [...]string{
	"COMPILED",
	"goog.DEBUG",
	"$jscomp.ASSUME_ES5",
	"$jscomp.ASSUME_ES6",
	"$jscomp.ASSUME_ES2020",
	"$jscomp.ISOLATE_POLYFILLS",
	"$jscomp.INSTRUMENT_ASYNC_CONTEXT",
}

var processDefines_KNOWN_DEFINES = map[string]struct{}{
	"COMPILED":                         {},
	"goog.DEBUG":                       {},
	"$jscomp.ASSUME_ES5":               {},
	"$jscomp.ASSUME_ES6":               {},
	"$jscomp.ASSUME_ES2020":            {},
	"$jscomp.ISOLATE_POLYFILLS":        {},
	"$jscomp.INSTRUMENT_ASYNC_CONTEXT": {},
}

var processDefines_GOOG_DEFINE = rhino.IR_Getprop(rhino.IR_Name("goog"), "define")

// Warnings

// Java: ProcessDefines.UNKNOWN_DEFINE_WARNING
var ProcessDefines_UNKNOWN_DEFINE_WARNING = DiagnosticType_Warning("JSC_UNKNOWN_DEFINE_WARNING", "unknown @define variable {0}")

// Errors

// Java: ProcessDefines.INVALID_DEFINE_NAME_ERROR
var ProcessDefines_INVALID_DEFINE_NAME_ERROR = DiagnosticType_Error(
	"JSC_INVALID_DEFINE_NAME_ERROR", "\"{0}\" is not a valid JS identifier name")

// Java: ProcessDefines.MISSING_DEFINE_ANNOTATION
var ProcessDefines_MISSING_DEFINE_ANNOTATION = DiagnosticType_Error("JSC_INVALID_MISSING_DEFINE_ANNOTATION", "Missing @define annotation")

// Java: ProcessDefines.INVALID_DEFINE_TYPE
var ProcessDefines_INVALID_DEFINE_TYPE = DiagnosticType_Error("JSC_INVALID_DEFINE_TYPE", "@define tag only permits primitive types")

// Java: ProcessDefines.INVALID_DEFINE_VALUE
var ProcessDefines_INVALID_DEFINE_VALUE = DiagnosticType_Error(
	"JSC_INVALID_DEFINE_VALUE", "invalid initialization value for @define {0}")

// Java: ProcessDefines.INVALID_DEFINE_LOCATION
var ProcessDefines_INVALID_DEFINE_LOCATION = DiagnosticType_Error(
	"JSC_INVALID_DEFINE_LOCATION",
	"@define must be initialized on a static qualified name in global or module scope")

// Java: ProcessDefines.NON_CONST_DEFINE
var ProcessDefines_NON_CONST_DEFINE = DiagnosticType_Error("JSC_NON_CONST_DEFINE", "@define {0} has already been set at {1}.")

// Java: ProcessDefines.CLOSURE_DEFINES_ERROR
var ProcessDefines_CLOSURE_DEFINES_ERROR = DiagnosticType_Error("JSC_CLOSURE_DEFINES_ERROR", "Invalid CLOSURE_DEFINES definition")

// Java: ProcessDefines.CLOSURE_DEFINES_MULTIPLE
var ProcessDefines_CLOSURE_DEFINES_MULTIPLE = DiagnosticType_Error(
	"JSC_CLOSURE_DEFINES_MULTIPLE",
	"Multiple CLOSURE_DEFINES definitions for {0}. First occurrence: {1}")

// Java: ProcessDefines.NON_GLOBAL_CLOSURE_DEFINES_ERROR
var ProcessDefines_NON_GLOBAL_CLOSURE_DEFINES_ERROR = DiagnosticType_Error(
	"JSC_NON_GLOBAL_CLOSURE_DEFINES_ERROR",
	"CLOSURE_DEFINES definition must be in top-level global scope")

// Java: ProcessDefines.DEFINE_CALL_WITHOUT_ASSIGNMENT
var ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT = DiagnosticType_Error(
	"JSC_DEFINE_CALL_WITHOUT_ASSIGNMENT",
	"The result of a goog.define call must be assigned as an isolated statement.")

// Java: ProcessDefines.ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT
var ProcessDefines_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT = DiagnosticType_Error(
	"JSC_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT",
	"ZoneJS is incompatible with language level ES2017 or higher (See go/ngissue/31730)\n"+
		"Please set `--language_out=ECMASCRIPT_2016` (or older) in your flags.")

// Java: ProcessDefines.DEFINE_WITHOUT_GOOG_DEFINE
var ProcessDefines_DEFINE_WITHOUT_GOOG_DEFINE = DiagnosticType_Disabled(
	"JSC_DEFINE_WITHOUT_GOOG_DEFINE",
	"The const declaration for @define ''{0}'' must be initialized with a call to"+
		" goog.define().")

// ProcessDefinesOptions configures a ProcessDefines pass.
type ProcessDefinesOptions struct {
	Replacements            OrderedStringMap[rhino.Node]
	Mode                    ProcessDefines_Mode
	NamespaceSupplier       func() *GlobalNamespace
	RecognizeClosureDefines bool
	EnableZonesDefineName   *string
	ZoneInputPattern        *jutil.Pattern
	UnknownDefinesToIgnore  []string
}

// NewProcessDefines creates a pass that overrides define constants.
func NewProcessDefines(compiler AbstractCompiler, opts ProcessDefinesOptions) *ProcessDefines {
	p := &ProcessDefines{
		mode:                    opts.Mode,
		compiler:                compiler,
		namespaceSupplier:       opts.NamespaceSupplier,
		recognizeClosureDefines: opts.RecognizeClosureDefines,
		zoneInputPattern:        opts.ZoneInputPattern,
	}
	if opts.EnableZonesDefineName != nil {
		p.enableZonesDefineName = *opts.EnableZonesDefineName
		p.hasEnableZonesDefineName = true
	}
	if p.mode.Check() {
		p.registry = p.compiler.GetTypeRegistry()
	}
	for _, e := range opts.Replacements.Entries() {
		rhino.CheckNotNull(e.Value)
	}
	p.replacementValuesFromFlags = opts.Replacements
	if len(opts.UnknownDefinesToIgnore) > 0 {
		p.unknownDefinesToIgnore = make(map[string]struct{}, len(opts.UnknownDefinesToIgnore))
		for _, d := range opts.UnknownDefinesToIgnore {
			p.unknownDefinesToIgnore[d] = struct{}{}
		}
	}
	return p
}

// ProcessDefines_Mode is ProcessDefines.Mode.
type ProcessDefines_Mode int32

const (
	ProcessDefines_Mode_CHECK ProcessDefines_Mode = iota
	ProcessDefines_Mode_OPTIMIZE
	ProcessDefines_Mode_CHECK_AND_OPTIMIZE
)

var (
	processDefines_Mode_names    = [...]string{"CHECK", "OPTIMIZE", "CHECK_AND_OPTIMIZE"}
	processDefines_Mode_check    = [...]bool{true, false, true}
	processDefines_Mode_optimize = [...]bool{false, true, true}
)

// Java: ProcessDefines.Mode.values()
func ProcessDefines_Mode_values() []ProcessDefines_Mode {
	return []ProcessDefines_Mode{ProcessDefines_Mode_CHECK, ProcessDefines_Mode_OPTIMIZE, ProcessDefines_Mode_CHECK_AND_OPTIMIZE}
}

// Java: ProcessDefines.Mode.name()
func (m ProcessDefines_Mode) String() string { return processDefines_Mode_names[m] }

// Check is the Mode.check field.
func (m ProcessDefines_Mode) Check() bool { return processDefines_Mode_check[m] }

// Optimize is the Mode.optimize field.
func (m ProcessDefines_Mode) Optimize() bool { return processDefines_Mode_optimize[m] }

// Java: ProcessDefines.process(Node, Node)
func (p *ProcessDefines) Process(externs rhino.Node, root rhino.Node) {
	if p.namespaceSupplier != nil {
		p.namespace = p.namespaceSupplier()
	}
	ownNamespace := p.namespace == nil
	p.InitNamespace(externs, root)
	if ownNamespace && p.namespace != nil {
		defer func() {
			p.namespace = nil
		}()
	}
	p.CollectDefines(root)
	p.ReportInvalidDefineLocations(root)
	p.CollectValidDefineValueExpressions()
	p.ValidateDefineDeclarations()
	p.OverrideDefines()
	p.ValidateDefines()
}

// Java: ProcessDefines.initNamespace(Node, Node)
func (p *ProcessDefines) InitNamespace(externs rhino.Node, root rhino.Node) {
	if p.namespace == nil && p.namespaceSupplier != nil {
		p.namespace = p.namespaceSupplier()
	}
	if p.namespace == nil {
		p.namespace = NewGlobalNamespaceWithRoots(p.compiler, externs, root)
		p.namespace.forProcessDefines = true
		p.namespace.processDefinesCheck = p.mode.Check()
	}
}

// Java: ProcessDefines.overrideDefines()
func (p *ProcessDefines) OverrideDefines() {
	if !p.mode.Optimize() {
		return
	}

	for _, define := range p.defines {
		if define.valueParent.IsNil() {
			continue
		}

		inputValue := p.GetReplacementForDefine(define)
		if inputValue.IsNil() || inputValue == define.value {
			continue
		}

		changed :=
			define.value.IsNil() ||
				inputValue.GetToken() != define.value.GetToken() ||
				!inputValue.IsEquivalentTo(define.value)
		if changed {
			if define.value.IsNil() {
				define.valueParent.AddChildToBack(inputValue.CloneTree())
			} else {
				define.value.ReplaceWith(inputValue.CloneTree())
			}

			p.compiler.ReportChangeToEnclosingScope(define.valueParent)
		}

		// When requested, create a globally accessible alias for all define
		// values. This will provide a hook for passes like
		// J2clUtilGetDefineRewriterPass to read them in out-of-scope
		// contexts.
		if J2clSourceFileChecker_ShouldRunJ2clPasses(p.compiler) {
			alias := ProcessDefines_GetGlobalDefineAlias(define.defineName)
			if alias != define.defineName {
				// If we had:
				//  var x = goog.define('y', ...);
				// we'll add an additional statement:
				//  var goog$defines$y = x;
				globalDefineName := NodeUtil_NewNameWithSrcref(p.compiler, alias, define.valueParent)
				var defineLhs rhino.Node
				if assign, ok := ast.AsAssign(define.valueParent); ok {
					defineLhs = assign.Target().Node().CloneTree()
				} else {
					defineLhs = define.valueParent.CloneNode()
				}
				jutil.CheckStateMsg1(defineLhs.IsName() || defineLhs.IsQualifiedName(), "%s", defineLhs)
				globalDefineVar :=
					rhino.IR_VarWithValue(globalDefineName, defineLhs).SrcrefTreeIfMissing(define.valueParent)
				globalDefineVar.InsertAfter(define.valueParent.GetParent())
				p.compiler.ReportChangeToEnclosingScope(define.valueParent.GetParent())
			}
		}
	}

	if !p.replacementValuesFromFlags.IsEmpty() || len(p.closureDefineOrder) > 0 {
		// Sets.difference(Sets.union(flags, closureDefines),
		//     Sets.union(Sets.union(KNOWN_DEFINES, defineByDefineName), ignore))
		isUsed := func(k string) bool {
			if _, known := processDefines_KNOWN_DEFINES[k]; known {
				return true
			}
			if _, ok := p.defineByDefineName[k]; ok {
				return true
			}
			_, ignored := p.unknownDefinesToIgnore[k]
			return ignored
		}
		for _, e := range p.replacementValuesFromFlags.Entries() {
			if !isUsed(e.Key) {
				p.compiler.Report(JSError_Make(ProcessDefines_UNKNOWN_DEFINE_WARNING, e.Key))
			}
		}
		for _, unknownDefine := range p.closureDefineOrder {
			if !p.replacementValuesFromFlags.ContainsKey(unknownDefine) && !isUsed(unknownDefine) {
				p.compiler.Report(JSError_Make(ProcessDefines_UNKNOWN_DEFINE_WARNING, unknownDefine))
			}
		}
	}
}

// Java: ProcessDefines.getGlobalDefineAlias(String)
func ProcessDefines_GetGlobalDefineAlias(defineName string) string {
	// Known defines are already globally accessible and may not be present in
	// code. Therefore we'll reference them directly.
	// goog.LOCALE is a special case as it's not processed as a define, it is
	// instead late-substituted.
	if _, known := processDefines_KNOWN_DEFINES[defineName]; known || defineName == "goog.LOCALE" {
		return defineName
	}
	return "jscomp$defines$" + strings.ReplaceAll(defineName, ".", "$")
}

// GetReplacementForDefine returns the replacement value for a @define, if
// any.
//
//  1. First checks the flags/compiler options `--define=FOO=1`
//  2. If nothing was found, check for values in a "var CLOSURE_DEFINES =
//     {'FOO': 1}` definition
//  3. If nothing was found, and this is defined via a goog.define call,
//     replace the call with the default value.
//
// Java: ProcessDefines.getReplacementForDefine(ProcessDefines.Define)
func (p *ProcessDefines) GetReplacementForDefine(define *ProcessDefines_Define) rhino.Node {
	replacement := p.GetReplacementForDefineName(define.defineName)
	if !replacement.IsNil() {
		return replacement
	}

	if p.IsGoogDefineCall(define.value) {
		if defaultValue, ok := processDefines_googDefineDefault(ast.MustCall(define.value)); ok {
			return defaultValue.Node()
		}
	}
	return rhino.NilNode
}

// processDefines_googDefineDefault returns the default value of
// goog.define('name', value), its second argument; ok is false if the call
// does not have exactly two.
//
// Go-only: Java's getChildCount() == 3 and getChildAtIndex(2) (the callee is
// the first child).
func processDefines_googDefineDefault(call ast.Call) (ast.Arg, bool) {
	if call.NumArgs() != 2 {
		return ast.Arg{}, false
	}
	return call.LastArg()
}

// GetReplacementForDefineName returns the replacement value for a @define,
// if any.
//
//  1. First checks the flags/compiler options `--define=FOO=1`
//  2. If nothing was found, check for values in a "var CLOSURE_DEFINES =
//     {'FOO': 1}` definition
//  3. If still not found, returns `null`.
//
// Java: ProcessDefines.getReplacementForDefineName(String)
func (p *ProcessDefines) GetReplacementForDefineName(defineName string) rhino.Node {
	if replacementFromFlags, ok := p.replacementValuesFromFlags.Lookup(defineName); ok && !replacementFromFlags.IsNil() {
		return replacementFromFlags
	}
	if replacementFromClosureDefines, ok := p.replacementValuesFromClosureDefines[defineName]; ok && !replacementFromClosureDefines.IsNil() {
		return replacementFromClosureDefines
	}
	return rhino.NilNode
}

// IsValidDefineType: only defines of literal number, string, or boolean are
// supported.
//
// Java: ProcessDefines.isValidDefineType(JSTypeExpression)
func (p *ProcessDefines) IsValidDefineType(expression *rhino.JSTypeExpression) bool {
	type_ := p.registry.EvaluateTypeExpressionInGlobalScope(expression)
	return !type_.IsUnknownType() && type_.IsSubtypeOf(p.registry.GetNativeType(rhino.JSTypeNative_NUMBER_STRING_BOOLEAN))
}

// Java: ProcessDefines.validateDefines()
func (p *ProcessDefines) ValidateDefines() {
	// Validate that a Zone-enabled app does not include native async/await,
	// which is incompatible.
	willOutputAsyncFunctions :=
		p.compiler.GetOptions().GetOutputFeatureSet().ContainsFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	if p.IsZoneEnabled() && p.hasZoneInput && willOutputAsyncFunctions {
		p.compiler.Report(JSError_Make(ProcessDefines_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT))
	}
}

// Java: ProcessDefines.isZoneEnabled()
func (p *ProcessDefines) IsZoneEnabled() bool {
	enableZonesDefault := true
	if !p.hasEnableZonesDefineName {
		return enableZonesDefault
	}
	zoneEnabled := p.GetReplacementForDefineName(p.enableZonesDefineName)
	if zoneEnabled.IsNil() {
		return enableZonesDefault
	}

	return NodeUtil_GetBooleanValue(zoneEnabled).ToBoolean(enableZonesDefault)
}

// ProcessDefines_Define is ProcessDefines.Define.
type ProcessDefines_Define struct {
	defineName string
	name       *GlobalNamespace_Name

	// The canonical set ref with an `@define` or `goog.define`.
	//
	// This may not be the same as `name.getDeclaration()`.
	declaration *GlobalNamespace_Ref

	valueParent rhino.Node
	value       rhino.Node
}

// Java: ProcessDefines.Define(String, GlobalNamespace.Name, GlobalNamespace.Ref, Node, Node)
func NewProcessDefines_Define(
	defineName string,
	name *GlobalNamespace_Name,
	declaration *GlobalNamespace_Ref,
	valueParent rhino.Node,
	value rhino.Node) *ProcessDefines_Define {
	jutil.CheckState(valueParent.IsNil() || value.IsNil() || value.GetParent() == valueParent)
	jutil.CheckState(declaration.IsSet())

	return &ProcessDefines_Define{
		defineName:  defineName,
		name:        name,
		declaration: declaration,
		valueParent: valueParent,
		value:       value,
	}
}
