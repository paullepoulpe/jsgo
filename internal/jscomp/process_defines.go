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
	"github.com/paullepoulpe/jsgo/internal/base"
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
	replacementValuesFromFlags map[string]rhino.Node
	flagReplacementOrder       []string
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

// newProcessDefines creates a pass that overrides define constants.
//
// Java: ProcessDefines(ProcessDefines.Builder)
func newProcessDefines(builder *ProcessDefines_Builder) *ProcessDefines {
	p := &ProcessDefines{
		mode:                     builder.mode,
		compiler:                 builder.compiler,
		namespaceSupplier:        builder.namespaceSupplier,
		recognizeClosureDefines:  builder.recognizeClosureDefines,
		enableZonesDefineName:    builder.enableZonesDefineName,
		hasEnableZonesDefineName: builder.hasEnableZonesDefineName,
		zoneInputPattern:         builder.zoneInputPattern,
	}
	if p.mode.Check() {
		p.registry = p.compiler.GetTypeRegistry()
	}
	if len(builder.replacementOrder) > 0 {
		p.replacementValuesFromFlags = make(map[string]rhino.Node, len(builder.replacementOrder))
		p.flagReplacementOrder = make([]string, len(builder.replacementOrder))
		copy(p.flagReplacementOrder, builder.replacementOrder)
		for _, k := range builder.replacementOrder {
			p.replacementValuesFromFlags[k] = rhino.CheckNotNull(builder.replacementValues[k])
		}
	}
	if len(builder.unknownDefinesToIgnore) > 0 {
		p.unknownDefinesToIgnore = make(map[string]struct{}, len(builder.unknownDefinesToIgnore))
		for _, d := range builder.unknownDefinesToIgnore {
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

// ProcessDefines_Builder is the builder for ProcessDefines.
type ProcessDefines_Builder struct {
	compiler                 AbstractCompiler
	replacementValues        map[string]rhino.Node
	replacementOrder         []string
	mode                     ProcessDefines_Mode
	namespaceSupplier        func() *GlobalNamespace
	recognizeClosureDefines  bool
	enableZonesDefineName    string
	hasEnableZonesDefineName bool
	zoneInputPattern         *jutil.Pattern
	unknownDefinesToIgnore   []string
}

// Java: ProcessDefines.Builder(AbstractCompiler)
func NewProcessDefines_Builder(compiler AbstractCompiler) *ProcessDefines_Builder {
	return &ProcessDefines_Builder{
		compiler:                compiler,
		recognizeClosureDefines: true,
	}
}

// PutReplacement sets a single define replacement value, preserving
// first-insertion order across overrides.
func (b *ProcessDefines_Builder) PutReplacement(name string, value rhino.Node) *ProcessDefines_Builder {
	if b.replacementValues == nil {
		b.replacementValues = make(map[string]rhino.Node)
	}
	if _, exists := b.replacementValues[name]; !exists {
		b.replacementOrder = append(b.replacementOrder, name)
	}
	b.replacementValues[name] = value
	return b
}

// Java: ProcessDefines.Builder.putReplacements(Map<String, Node>)
func (b *ProcessDefines_Builder) PutReplacements(replacementValues jutil.MapSource[string, rhino.Node]) *ProcessDefines_Builder {
	for k, v := range replacementValues.All() {
		b.PutReplacement(k, v)
	}
	return b
}

// Java: ProcessDefines.Builder.setMode(ProcessDefines.Mode)
func (b *ProcessDefines_Builder) SetMode(x ProcessDefines_Mode) *ProcessDefines_Builder {
	b.mode = x
	return b
}

// InjectNamespace injects a pre-computed global namespace, so that the same
// namespace can be re-used for multiple check passes. Accepts a supplier
// because the namespace may not exist at pass-creation time.
//
// Java: ProcessDefines.Builder.injectNamespace(Supplier<GlobalNamespace>)
func (b *ProcessDefines_Builder) InjectNamespace(namespaceSupplier func() *GlobalNamespace) *ProcessDefines_Builder {
	b.namespaceSupplier = namespaceSupplier
	return b
}

// Java: ProcessDefines.Builder.setUnknownDefinesToIgnore(ImmutableList<String>)
func (b *ProcessDefines_Builder) SetUnknownDefinesToIgnore(unknownDefinesToIgnore []string) *ProcessDefines_Builder {
	b.unknownDefinesToIgnore = unknownDefinesToIgnore
	return b
}

// Java: ProcessDefines.Builder.setRecognizeClosureDefines(boolean)
func (b *ProcessDefines_Builder) SetRecognizeClosureDefines(recognizeClosureDefines bool) *ProcessDefines_Builder {
	b.recognizeClosureDefines = recognizeClosureDefines
	return b
}

// Java: ProcessDefines.Builder.setEnableZonesDefineName(String)
func (b *ProcessDefines_Builder) SetEnableZonesDefineName(enableZonesDefineName *string) *ProcessDefines_Builder {
	if enableZonesDefineName != nil {
		b.enableZonesDefineName = *enableZonesDefineName
		b.hasEnableZonesDefineName = true
	} else {
		b.enableZonesDefineName = ""
		b.hasEnableZonesDefineName = false
	}
	return b
}

// Java: ProcessDefines.Builder.setZoneInputPattern(Pattern)
func (b *ProcessDefines_Builder) SetZoneInputPattern(zoneInputPattern *jutil.Pattern) *ProcessDefines_Builder {
	b.zoneInputPattern = zoneInputPattern
	return b
}

// Java: ProcessDefines.Builder.build()
func (b *ProcessDefines_Builder) Build() *ProcessDefines {
	return newProcessDefines(b)
}

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

	if len(p.flagReplacementOrder) > 0 || len(p.closureDefineOrder) > 0 {
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
		for _, unknownDefine := range p.flagReplacementOrder {
			if !isUsed(unknownDefine) {
				p.compiler.Report(JSError_Make(ProcessDefines_UNKNOWN_DEFINE_WARNING, unknownDefine))
			}
		}
		for _, unknownDefine := range p.closureDefineOrder {
			if _, fromFlags := p.replacementValuesFromFlags[unknownDefine]; !fromFlags && !isUsed(unknownDefine) {
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
	if replacementFromFlags, ok := p.replacementValuesFromFlags[defineName]; ok && !replacementFromFlags.IsNil() {
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

// CollectDefines finds all defines, and creates a Define data structure for
// each one.
//
// Java: ProcessDefines.collectDefines(Node)
func (p *ProcessDefines) CollectDefines(root rhino.Node) {
	if p.recognizeClosureDefines && !root.IsNil() {
		c := ProcessDefines_ClosureDefinesCollector{outer: p}
		c.walk(root)
	}
	for name := range p.namespace.GetAllSymbols() {
		declaration := p.SelectDefineDeclaration(name)
		if declaration == nil {
			continue
		}

		totalSets := name.GetTotalSets()
		valueParent := processDefines_getValueParentForDefine(declaration)
		var value rhino.Node
		if !valueParent.IsNil() {
			value = valueParent.GetLastChild()
		}

		var defineName string
		if p.IsGoogDefineCall(value) && p.VerifyGoogDefine(value) {
			// VerifyGoogDefine checked that the first argument is a string.
			nameArg, _ := ast.MustCall(value).FirstArg()
			defineName = ast.MustStringLit(nameArg.Node()).Value()
		} else {
			defineName = name.GetFullName()
		}

		if existingDefine, exists := p.defineByDefineName[defineName]; exists {
			declaration = existingDefine.declaration
			totalSets += existingDefine.name.GetTotalSets()
		} else {
			def := NewProcessDefines_Define(defineName, name, declaration, valueParent, value)
			if p.defineByDefineName == nil {
				p.defineByDefineName = make(map[string]*ProcessDefines_Define)
			}
			p.defineByDefineName[defineName] = def
			p.defines = append(p.defines, def)
		}

		/*
		 * We have to report this here because otherwise we don't remember
		 * which names have the same define name. It's not worth it tracking a
		 * set of names, because it makes the rest of the pass more complex.
		 */
		if totalSets > 1 {
			for ref := range name.Refs() {
				if ref.IsSet() && ref != declaration {
					p.compiler.Report(
						JSError_MakeWithNode(
							ref.GetNode(),
							ProcessDefines_NON_CONST_DEFINE,
							defineName,
							declaration.GetNode().GetLocation()))
				}
			}
		}
	}
	defineNames := make([]string, 0, len(processDefinesKnownDefineList)+len(p.defines))
	defineNames = append(defineNames, processDefinesKnownDefineList[:]...)
	for _, def := range p.defines {
		if _, known := processDefines_KNOWN_DEFINES[def.defineName]; !known {
			defineNames = append(defineNames, def.defineName)
		}
	}
	p.compiler.SetDefineNames(defineNames)
}

// Java: ProcessDefines.selectDefineDeclaration(GlobalNamespace.Name)
func (p *ProcessDefines) SelectDefineDeclaration(name *GlobalNamespace_Name) *GlobalNamespace_Ref {
	if name.GetGlobalSets() == 0 {
		return nil
	}
	for ref := range name.Refs() {
		// Make sure we don't select a local set as the declaration.
		if !ref.IsSetFromGlobal() {
			continue
		}

		refNode := ref.GetNode()
		if !refNode.IsQualifiedName() {
			continue
		}

		jsdoc := NodeUtil_GetBestJSDocInfo(refNode)
		if jsdoc == nil || !jsdoc.IsDefine() {
			continue
		}

		if p.knownDefineJsdocs == nil {
			p.knownDefineJsdocs = make(map[*rhino.JSDocInfo]struct{})
		}
		p.knownDefineJsdocs[jsdoc] = struct{}{}
		return ref
	}

	return nil
}

// Java: ProcessDefines.getValueParentForDefine(GlobalNamespace.Ref)
func processDefines_getValueParentForDefine(declaration *GlobalNamespace_Ref) rhino.Node {
	// Note: this may be a NAME, a GETPROP, or even STRING_KEY or GETTER_DEF.
	// We only care about the first two, in which case the parent should be
	// either VAR/CONST or ASSIGN. We could accept STRING_KEY (i.e. `@define`
	// on a property in an object literal), but there's no reason to add
	// another new way to do the same thing.
	declarationNode := declaration.GetNode()
	declarationParent := declarationNode.GetParent()

	if declarationParent.IsVar() || declarationParent.IsConst() {
		// Simple case of `var` or `const`. There's no reason to support `let`
		// here, and we don't explicitly check that it's not `let` anywhere
		// else.
		jutil.CheckStateMsg1(declarationNode.IsName(), "%s", declarationNode)
		return declarationNode
	} else if assign, ok := ast.AsAssign(declarationParent); ok && assign.Target().Is(declarationNode) {
		// Assignment. Must either assign to a qualified name, or else be a
		// different ref than the declaration to not emit an error (we don't
		// allow assignment before it's declared).
		return declarationParent
	}
	return rhino.NilNode
}

// Java: ProcessDefines.collectValidDefineValueExpressions()
func (p *ProcessDefines) CollectValidDefineValueExpressions() {
	if !p.mode.Check() || len(p.defines) == 0 {
		return
	}

	if p.validDefineValueExpressions == nil {
		p.validDefineValueExpressions = make(map[rhino.Node]struct{})
	}

	defineNamesSet := make(map[*GlobalNamespace_Name]struct{}, len(p.defines))
	// All defines are implicitly valid in the values of other defines.
	for _, define := range p.defines {
		defineNamesSet[define.name] = struct{}{}
		for r := range define.name.Refs() {
			if !r.IsSet() {
				p.validDefineValueExpressions[r.GetNode()] = struct{}{}
			}
		}
	}

	hasUnknownDefineValue := false
	for _, define := range p.defines {
		if !define.valueParent.IsNil() && !define.valueParent.IsFromExterns() && p.IsValidDefineValue(define.value) == base.Tri_UNKNOWN {
			hasUnknownDefineValue = true
			break
		}
	}
	if !hasUnknownDefineValue {
		return
	}

	var namesToCheck []*GlobalNamespace_Name
	for name := range p.namespace.GetAllSymbols() {
		if !processDefines_isGlobalConst(name) {
			continue
		}
		if _, isDef := defineNamesSet[name]; isDef {
			continue
		}
		declValue := processDefines_getConstantDeclValue(name.getDeclarationRef().GetNode())
		if declValue.IsNil() {
			continue
		}
		switch p.IsValidDefineValue(declValue) {
		case base.Tri_TRUE:
			for ref := range name.Refs() {
				p.validDefineValueExpressions[ref.GetNode()] = struct{}{}
			}
		case base.Tri_UNKNOWN:
			namesToCheck = append(namesToCheck, name)
		default:
		}
	}

	// Do a breadth-first search of all const names to find those defined in
	// terms of valid values.
	for len(namesToCheck) > 0 {
		var namesToCheckAgain []*GlobalNamespace_Name

		for _, name := range namesToCheck {
			declValue := processDefines_getConstantDeclValue(name.getDeclarationRef().GetNode())
			switch p.IsValidDefineValue(declValue) {
			case base.Tri_TRUE:
				for ref := range name.Refs() {
					p.validDefineValueExpressions[ref.GetNode()] = struct{}{}
				}
			case base.Tri_UNKNOWN:
				namesToCheckAgain = append(namesToCheckAgain, name)
			default:
			}
		}

		if len(namesToCheckAgain) == len(namesToCheck) {
			break
		}
		namesToCheck = namesToCheckAgain
	}
}

// Java: ProcessDefines.validateDefineDeclarations()
func (p *ProcessDefines) ValidateDefineDeclarations() {
	if !p.mode.Check() {
		return
	}

	for _, define := range p.defines {
		declarationNode := define.declaration.GetNode()

		_, known := processDefines_KNOWN_DEFINES[define.defineName]
		if !p.IsGoogDefineCall(define.value) &&
			!known &&
			!strings.HasPrefix(define.defineName, "$jscomp.") {
			errorNode := define.value
			if errorNode.IsNil() {
				errorNode = define.valueParent
				if errorNode.IsNil() {
					errorNode = rhino.CheckNotNull(declarationNode)
				}
			}
			p.compiler.Report(
				JSError_MakeWithNode(
					errorNode,
					ProcessDefines_DEFINE_WITHOUT_GOOG_DEFINE,
					define.defineName,
					"go/js-practices/locals#defines")) // Example help link
		}

		if !p.HasValidValue(define) {
			errorNode := define.value
			if errorNode.IsNil() {
				errorNode = define.valueParent
				if errorNode.IsNil() {
					errorNode = rhino.CheckNotNull(declarationNode)
				}
			}
			p.compiler.Report(
				JSError_MakeWithNode(
					errorNode,
					ProcessDefines_INVALID_DEFINE_VALUE,
					define.defineName))
		}

		/*
		 * Process defines should not depend on check types being enabled, so
		 * we look for the JSDoc instead of the inferred type.
		 */
		jsdoc := NodeUtil_GetBestJSDocInfo(declarationNode)
		if jsdoc == nil || !p.IsValidDefineType(jsdoc.GetType()) {
			p.compiler.Report(JSError_MakeWithNode(declarationNode, ProcessDefines_INVALID_DEFINE_TYPE))
		}
	}
}

// ReportInvalidDefineLocations checks for misplaced @define and goog.define
// calls.
//
// Java: ProcessDefines.reportInvalidDefineLocations(Node)
func (p *ProcessDefines) ReportInvalidDefineLocations(root rhino.Node) {
	if !p.mode.Check() {
		return
	}

	/*
	 * This has to be done using a traversal because the global namespace
	 * doesn't record symbols which only appear in local scopes.
	 *
	 * <p>We don't check the externs because they can't contain local vars.
	 */
	cb := func(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
		jsdoc := n.GetJSDocInfo()
		if jsdoc != nil && jsdoc.IsDefine() {
			if _, known := p.knownDefineJsdocs[jsdoc]; !known {
				if p.knownDefineJsdocs == nil {
					p.knownDefineJsdocs = make(map[*rhino.JSDocInfo]struct{})
				}
				p.knownDefineJsdocs[jsdoc] = struct{}{}
				p.compiler.Report(JSError_MakeWithNode(n, ProcessDefines_INVALID_DEFINE_LOCATION))
			}
		}

		if p.IsGoogDefineCall(n) {
			if _, known := p.knownGoogDefineCalls[n]; !known {
				p.VerifyGoogDefine(n)
			}
		}

		if n.MatchesName("CLOSURE_DEFINES") &&
			(NodeUtil_IsNameDeclaration(parent) ||
				(parent.IsGetElem() && parent.GetParent().IsAssign())) &&
			!NodeUtil_GetEnclosingScopeRoot(n).IsRoot() {
			p.compiler.Report(JSError_MakeWithNode(n, ProcessDefines_NON_GLOBAL_CLOSURE_DEFINES_ERROR))
		}
	}
	if p.namespace != nil && p.namespace.forProcessDefines && p.namespace.generated && !root.IsNil() && root.IsRoot() && !p.namespace.hasNonScriptDefineCand {
		if len(p.namespace.defineCandidateScripts) == 0 {
			return
		}
		b := NodeTraversal_Builder_().SetCompiler(p.compiler).SetCallbackFunc(cb)
		for s := root.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
			if _, ok := p.namespace.defineCandidateScripts[s]; ok {
				b.Traverse(s)
			}
		}
		return
	}
	NodeTraversal_Builder_().
		SetCompiler(p.compiler).
		SetCallbackFunc(cb).
		Traverse(root)
}

// ProcessDefines_ClosureDefinesCollector is
// ProcessDefines.ClosureDefinesCollector (an inner class).
type ProcessDefines_ClosureDefinesCollector struct {
	outer *ProcessDefines
}

func (c *ProcessDefines_ClosureDefinesCollector) walk(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_EXPR_RESULT:
		for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
			c.walk(ch)
		}
	case rhino.Token_SCRIPT:
		for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
			c.walk(ch)
		}
		p := c.outer
		if p.zoneInputPattern != nil {
			source := n.GetSourceFileName()
			if source != nil && p.zoneInputPattern.Matches(*source) {
				p.hasZoneInput = true
			}
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
			c.walk(ch)
		}
		c.outer.handleClosureDefinesDeclaration(ast.MustNameDeclaration(n))
	case rhino.Token_ASSIGN:
		for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
			c.walk(ch)
		}
		c.outer.handleClosureDefinesAssign(ast.MustAssign(n))
	}
}

// handleClosureDefinesDeclaration records the values of
// `var CLOSURE_DEFINES = {'name': value, ...}`.
//
// Go-only: shared by the walk above and Visit.
func (p *ProcessDefines) handleClosureDefinesDeclaration(decl ast.NameDeclaration) {
	first, _ := decl.FirstDeclarator()
	name, ok := first.AsName()
	if !ok || name.Ident() != "CLOSURE_DEFINES" {
		return
	}
	value, ok := name.Init()
	if !ok {
		return
	}
	obj, ok := value.AsObjectLit()
	if !ok {
		return
	}
	for member := range obj.Members() {
		// The first child is the value of a string key; HandleClosureDefinesValue
		// reports the other kinds of members.
		p.HandleClosureDefinesValue(member.Node(), member.Node().GetFirstChild(), decl.Node())
	}
}

// handleClosureDefinesAssign records the value of
// `CLOSURE_DEFINES['name'] = value`.
//
// Go-only: shared by the walk above and Visit.
func (p *ProcessDefines) handleClosureDefinesAssign(assign ast.Assign) {
	getElem, ok := assign.Target().AsGetElem()
	if ok && getElem.Object().Node().MatchesName("CLOSURE_DEFINES") {
		p.HandleClosureDefinesValue(getElem.Index().Node(), assign.Value().Node(), assign.Node())
	}
}

// Java: ProcessDefines.ClosureDefinesCollector.shouldTraverse(NodeTraversal, Node, Node)
func (c *ProcessDefines_ClosureDefinesCollector) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// In particular, don't traverse into modules or functions - only script
	// top level scopes.
	return n.IsRoot() ||
		n.IsScript() ||
		n.IsExprResult() ||
		n.IsAssign() ||
		NodeUtil_IsNameDeclaration(n)
}

// Java: ProcessDefines.ClosureDefinesCollector.visit(NodeTraversal, Node, Node)
func (c *ProcessDefines_ClosureDefinesCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	p := c.outer
	if n.IsScript() && p.zoneInputPattern != nil {
		source := n.GetSourceFileName()
		if source != nil && p.zoneInputPattern.Matches(*source) {
			p.hasZoneInput = true
		}
	}

	if NodeUtil_IsNameDeclaration(n) {
		// var CLOSURE_DEFINES = {...};
		p.handleClosureDefinesDeclaration(ast.MustNameDeclaration(n))
	} else if assign, ok := ast.AsAssign(n); ok {
		// CLOSURE_DEFINES['...'] = ...;
		p.handleClosureDefinesAssign(assign)
	}
}

// Java: ProcessDefines.handleClosureDefinesValue(Node, Node, Node)
func (p *ProcessDefines) HandleClosureDefinesValue(stringNode rhino.Node, valueNode rhino.Node, errorNode rhino.Node) {
	if (stringNode.IsStringKey() || stringNode.IsStringLit()) &&
		processDefines_isValidClosureDefinesValue(valueNode) {
		key := stringNode.GetString()
		if prev, exists := p.replacementValuesFromClosureDefines[key]; exists {
			p.compiler.Report(
				JSError_MakeWithNode(
					errorNode,
					ProcessDefines_CLOSURE_DEFINES_MULTIPLE,
					key,
					prev.GetLocation()))
		} else {
			if p.replacementValuesFromClosureDefines == nil {
				p.replacementValuesFromClosureDefines = make(map[string]rhino.Node)
			}
			p.closureDefineOrder = append(p.closureDefineOrder, key)
		}
		p.replacementValuesFromClosureDefines[key] = valueNode
	} else if p.mode.Check() {
		p.compiler.Report(JSError_MakeWithNode(errorNode, ProcessDefines_CLOSURE_DEFINES_ERROR))
	}
}

// Java: ProcessDefines.isValidClosureDefinesValue(Node)
func processDefines_isValidClosureDefinesValue(val rhino.Node) bool {
	// Values allowed in 'var CLOSURE_DEFINES = {'
	// Must be a subset of the values allowed for <val> in
	// /** @define {...} */ var DEF = <val>
	switch val.GetToken() {
	case rhino.Token_STRINGLIT, rhino.Token_NUMBER, rhino.Token_TRUE, rhino.Token_FALSE:
		return true
	case rhino.Token_NEG:
		return ast.MustUnary(val).Operand().Token() == rhino.Token_NUMBER
	default:
		return false
	}
}

// Java: ProcessDefines.hasValidValue(ProcessDefines.Define)
func (p *ProcessDefines) HasValidValue(define *ProcessDefines_Define) bool {
	if define.valueParent.IsNil() {
		return false
	} else if define.valueParent.IsFromExterns() {
		return true
	} else {
		return p.IsValidDefineValue(define.value).ToBoolean(false)
	}
}

// Java: ProcessDefines.isGlobalConst(GlobalNamespace.Name)
func processDefines_isGlobalConst(name *GlobalNamespace_Name) bool {
	if name.GetTotalSets() != 1 {
		return false
	}
	declaration, _ := name.GetDeclaration().(*GlobalNamespace_Ref)
	return declaration != nil && declaration.IsSetFromGlobal()
}

// IsValidDefineValue determines whether the given value may be assigned to a
// define. val is the value being assigned.
//
// Java: ProcessDefines.isValidDefineValue(Node)
func (p *ProcessDefines) IsValidDefineValue(val rhino.Node) base.Tri {
	if val.IsNil() {
		return base.Tri_FALSE
	}

	switch val.GetToken() {
	case rhino.Token_STRINGLIT, rhino.Token_NUMBER, rhino.Token_TRUE, rhino.Token_FALSE:
		return base.Tri_TRUE
		// Binary operators are only valid if both children are valid.
	case rhino.Token_AND,
		rhino.Token_OR,
		rhino.Token_COALESCE,
		rhino.Token_ADD,
		rhino.Token_BITAND,
		rhino.Token_BITOR,
		rhino.Token_BITXOR,
		rhino.Token_DIV,
		rhino.Token_EQ,
		rhino.Token_EXPONENT,
		rhino.Token_GE,
		rhino.Token_GT,
		rhino.Token_LE,
		rhino.Token_LSH,
		rhino.Token_LT,
		rhino.Token_MOD,
		rhino.Token_MUL,
		rhino.Token_NE,
		rhino.Token_RSH,
		rhino.Token_SHEQ,
		rhino.Token_SHNE,
		rhino.Token_SUB,
		rhino.Token_URSH:
		op := ast.MustBinary(val)
		return p.IsValidDefineValue(op.Left().Node()).And(p.IsValidDefineValue(op.Right().Node()))
	case rhino.Token_HOOK:
		hook := ast.MustHook(val)
		return p.IsValidDefineValue(hook.Cond().Node()).
			And(p.IsValidDefineValue(hook.Then().Node())).
			And(p.IsValidDefineValue(hook.Else().Node()))
		// Unary operators are valid if the child is valid.
		//
		// Go-only: Java lists BITNOT with the binary operators, where its first
		// and last child are both its operand; x.and(x) is x, so it belongs here.
	case rhino.Token_NOT, rhino.Token_NEG, rhino.Token_POS, rhino.Token_BITNOT:
		return p.IsValidDefineValue(ast.MustUnary(val).Operand().Node())
		// Names are valid if and only if they are defines themselves.
	case rhino.Token_NAME, rhino.Token_GETPROP:
		if val.IsQualifiedName() {
			if _, ok := p.validDefineValueExpressions[val]; ok {
				return base.Tri_TRUE
			}
			return base.Tri_UNKNOWN
		}
		// Allow goog.define('XYZ', <val>) calls if and only if <val> is valid.
	case rhino.Token_CALL:
		if !p.IsGoogDefineCall(val) {
			return base.Tri_FALSE
		}
		defaultValue, ok := processDefines_googDefineDefault(ast.MustCall(val))
		if !ok {
			// goog.define call with wrong arg count. Warn elsewhere and treat
			// this call as valid.
			return base.Tri_TRUE
		}
		return p.IsValidDefineValue(defaultValue.Node())
	default:
	}

	return base.Tri_FALSE
}

// processDefines_getConstantDeclValue checks whether the NAME node is
// inside either a CONST or a @const VAR. Returns the RHS node if so,
// otherwise returns NilNode.
//
// Java: ProcessDefines.getConstantDeclValue(Node)
func processDefines_getConstantDeclValue(name rhino.Node) rhino.Node {
	parent := name.GetParent()
	if parent.IsNil() {
		return rhino.NilNode
	}
	if nameNode, ok := ast.AsName(name); ok {
		if parent.IsConst() {
			init, _ := nameNode.Init()
			return init.Node()
		}
		if !parent.IsVar() {
			return rhino.NilNode
		}
		jsdoc := NodeUtil_GetBestJSDocInfo(name)
		if jsdoc != nil && jsdoc.IsConstant() {
			init, _ := nameNode.Init()
			return init.Node()
		}
		return rhino.NilNode
	}
	if assign, ok := ast.AsAssign(parent); ok && name.IsGetProp() && assign.Target().Is(name) {
		jsdoc := NodeUtil_GetBestJSDocInfo(name)
		if jsdoc != nil && jsdoc.IsConstant() {
			return assign.Value().Node()
		}
	}
	return rhino.NilNode
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

// Java: ProcessDefines.isGoogDefineCall(Node)
func (p *ProcessDefines) IsGoogDefineCall(node rhino.Node) bool {
	if !p.recognizeClosureDefines {
		return false
	}
	call, ok := ast.AsCall(node)
	return ok && call.Callee().Node().MatchesQualifiedNameNode(processDefines_GOOG_DEFINE)
}

// VerifyGoogDefine verifies that a goog.define method call has exactly two
// arguments, with the first a string literal whose contents is a valid JS
// qualified name. Reports a compile error if it doesn't. Returns whether the
// argument checked out okay.
//
// Java: ProcessDefines.verifyGoogDefine(Node)
func (p *ProcessDefines) VerifyGoogDefine(callNode rhino.Node) bool {
	if p.knownGoogDefineCalls == nil {
		p.knownGoogDefineCalls = make(map[rhino.Node]struct{})
	}
	p.knownGoogDefineCalls[callNode] = struct{}{}

	parent := callNode.GetParent()
	methodName := callNode.GetFirstChild()
	args := callNode.GetSecondChild()

	// Calls to goog.define must be in the global hoist scope after module
	// rewriting
	if !NodeUtil_GetEnclosingFunction(callNode).IsNil() {
		p.compiler.Report(JSError_MakeWithNode(methodName.GetParent(), ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR))
		return false
	}

	// It is an error for goog.define to show up anywhere except immediately
	// after =.
	if parent.IsAssign() && parent.GetParent().IsExprResult() {
		parent = parent.GetParent()
	} else if parent.IsName() && NodeUtil_IsNameDeclaration(parent.GetParent()) {
		parent = parent.GetParent()
	} else {
		p.compiler.Report(JSError_MakeWithNode(methodName.GetParent(), ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT))
		return false
	}

	// Verify first arg
	arg := args
	if !p.VerifyNotNull(methodName, arg) || !p.VerifyOfType(methodName, arg, rhino.Token_STRINGLIT) {
		return false
	}

	// Verify second arg
	arg = arg.GetNext()
	if !args.IsFromExterns() &&
		(!p.VerifyNotNull(methodName, arg) || !p.VerifyIsLast(methodName, arg)) {
		return false
	}

	name := args.GetString()
	if !NodeUtil_IsValidQualifiedNameFeatureSet(
		p.compiler.GetOptions().GetLanguageIn().ToFeatureSet(), name) {
		p.compiler.Report(JSError_MakeWithNode(args, ProcessDefines_INVALID_DEFINE_NAME_ERROR, name))
		return false
	}

	infoNode := parent
	if exprStmt, ok := ast.AsExprStmt(parent); ok {
		infoNode = exprStmt.Expr().Node()
	}
	info := infoNode.GetJSDocInfo()
	if info == nil || !info.IsDefine() {
		p.compiler.Report(JSError_MakeWithNode(parent, ProcessDefines_MISSING_DEFINE_ANNOTATION))
		return false
	}
	return true
}

func processDefines_qualifiedNameOrNull(n rhino.Node) string {
	if qname, ok := n.GetQualifiedNameString(); ok {
		return qname
	}
	return "null"
}

// VerifyNotNull returns whether the argument checked out okay.
//
// Java: ProcessDefines.verifyNotNull(Node, Node)
func (p *ProcessDefines) VerifyNotNull(methodName rhino.Node, arg rhino.Node) bool {
	if arg.IsNil() {
		p.compiler.Report(
			JSError_MakeWithNode(
				methodName,
				ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR,
				processDefines_qualifiedNameOrNull(methodName)))
		return false
	}
	return true
}

// VerifyIsLast returns whether the argument checked out okay.
//
// Java: ProcessDefines.verifyIsLast(Node, Node)
func (p *ProcessDefines) VerifyIsLast(methodName rhino.Node, arg rhino.Node) bool {
	if !arg.GetNext().IsNil() {
		p.compiler.Report(
			JSError_MakeWithNode(
				methodName,
				ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR,
				processDefines_qualifiedNameOrNull(methodName)))
		return false
	}
	return true
}

// VerifyOfType returns whether the argument checked out okay.
//
// Java: ProcessDefines.verifyOfType(Node, Node, Token)
func (p *ProcessDefines) VerifyOfType(methodName rhino.Node, arg rhino.Node, desiredType rhino.Token) bool {
	if arg.GetToken() != desiredType {
		p.compiler.Report(
			JSError_MakeWithNode(
				methodName,
				ClosurePrimitiveErrors_INVALID_ARGUMENT_ERROR,
				processDefines_qualifiedNameOrNull(methodName)))
		return false
	}
	return true
}
