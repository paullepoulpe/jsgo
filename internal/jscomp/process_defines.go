// Port of com/google/javascript/jscomp/ProcessDefines.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ProcessDefines processes variables annotated as `@define`. A define is a
// special constant that may be overridden by later files and manipulated by
// the compiler, much like C preprocessor `#define`s.
type ProcessDefines struct {
	compiler                   AbstractCompiler
	registry                   *rhino.JSTypeRegistry
	replacementValuesFromFlags *jutil.LinkedHashMap[string, rhino.Node]
	mode                       ProcessDefines_Mode
	namespaceSupplier          func() *GlobalNamespace
	recognizeClosureDefines    bool
	enableZonesDefineName      *string
	zoneInputPattern           *jutil.Pattern
	unknownDefinesToIgnore     *jutil.LinkedHashSet[string]

	knownDefineJsdocs    *jutil.LinkedHashSet[*rhino.JSDocInfo]
	knownGoogDefineCalls *jutil.LinkedHashSet[rhino.Node]
	defineByDefineName   *jutil.LinkedHashMap[string, *ProcessDefines_Define]
	// from var CLOSURE_DEFINES = {
	replacementValuesFromClosureDefines *jutil.LinkedHashMap[string, rhino.Node]
	validDefineValueExpressions         *jutil.LinkedHashSet[rhino.Node]
	hasZoneInput                        bool

	namespace *GlobalNamespace
}

var _ CompilerPass = (*ProcessDefines)(nil)

// Defines in this set will not be flagged with "unknown define" warnings.
// There are flags that always set these defines, even when they might not be
// in the binary.
var processDefines_KNOWN_DEFINES = jutil.NewLinkedHashSet(
	"COMPILED",
	"goog.DEBUG",
	"$jscomp.ASSUME_ES5",
	"$jscomp.ASSUME_ES6",
	"$jscomp.ASSUME_ES2020",
	"$jscomp.ISOLATE_POLYFILLS",
	"$jscomp.INSTRUMENT_ASYNC_CONTEXT")

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
		knownDefineJsdocs:                   jutil.NewLinkedHashSet[*rhino.JSDocInfo](),
		knownGoogDefineCalls:                jutil.NewLinkedHashSet[rhino.Node](),
		defineByDefineName:                  jutil.NewLinkedHashMap[string, *ProcessDefines_Define](),
		replacementValuesFromClosureDefines: jutil.NewLinkedHashMap[string, rhino.Node](),
		validDefineValueExpressions:         jutil.NewLinkedHashSet[rhino.Node](),
	}
	p.mode = builder.mode
	p.compiler = builder.compiler
	if p.mode.Check() {
		p.registry = p.compiler.GetTypeRegistry()
	}
	p.replacementValuesFromFlags = closureModuleProcessor_immutableMapCopyOf(builder.replacementValues)
	p.namespaceSupplier = builder.namespaceSupplier
	p.recognizeClosureDefines = builder.recognizeClosureDefines
	p.enableZonesDefineName = builder.enableZonesDefineName
	p.zoneInputPattern = builder.zoneInputPattern
	p.unknownDefinesToIgnore = jutil.NewLinkedHashSet(builder.unknownDefinesToIgnore...)
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
	compiler                AbstractCompiler
	replacementValues       *jutil.LinkedHashMap[string, rhino.Node]
	mode                    ProcessDefines_Mode
	namespaceSupplier       func() *GlobalNamespace
	recognizeClosureDefines bool
	enableZonesDefineName   *string
	zoneInputPattern        *jutil.Pattern
	unknownDefinesToIgnore  []string
}

// Java: ProcessDefines.Builder(AbstractCompiler)
func NewProcessDefines_Builder(compiler AbstractCompiler) *ProcessDefines_Builder {
	return &ProcessDefines_Builder{
		compiler:                compiler,
		replacementValues:       jutil.NewLinkedHashMap[string, rhino.Node](),
		recognizeClosureDefines: true,
		unknownDefinesToIgnore:  []string{},
	}
}

// Java: ProcessDefines.Builder.putReplacements(Map<String, Node>)
func (b *ProcessDefines_Builder) PutReplacements(replacementValues jutil.MapSource[string, rhino.Node]) *ProcessDefines_Builder {
	for k, v := range replacementValues.All() {
		b.replacementValues.Put(k, v)
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
	b.enableZonesDefineName = enableZonesDefineName
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
			p.namespace.Release()
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
	if p.mode.Optimize() {
		for define := range p.defineByDefineName.Values() {
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
					globalDefineName :=
						NodeUtil_NewNameWithSrcref(
							p.compiler, ProcessDefines_GetGlobalDefineAlias(define.defineName), define.valueParent)
					var defineLhs rhino.Node
					if define.valueParent.IsAssign() {
						defineLhs = define.valueParent.GetFirstChild().CloneTree()
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
	}

	if p.mode.Optimize() {
		// Sets.difference(Sets.union(flags, closureDefines),
		//     Sets.union(Sets.union(KNOWN_DEFINES, defineByDefineName), ignore))
		var unusedReplacements []string
		isUsed := func(k string) bool {
			return processDefines_KNOWN_DEFINES.Contains(k) ||
				p.defineByDefineName.ContainsKey(k) ||
				p.unknownDefinesToIgnore.Contains(k)
		}
		for k := range p.replacementValuesFromFlags.Keys() {
			if !isUsed(k) {
				unusedReplacements = append(unusedReplacements, k)
			}
		}
		for k := range p.replacementValuesFromClosureDefines.Keys() {
			if !p.replacementValuesFromFlags.ContainsKey(k) && !isUsed(k) {
				unusedReplacements = append(unusedReplacements, k)
			}
		}

		for _, unknownDefine := range unusedReplacements {
			p.compiler.Report(JSError_Make(ProcessDefines_UNKNOWN_DEFINE_WARNING, unknownDefine))
		}
	}
}

// Java: ProcessDefines.getGlobalDefineAlias(String)
func ProcessDefines_GetGlobalDefineAlias(defineName string) string {
	// Known defines are already globally accessible and may not be present in
	// code. Therefore we'll reference them directly.
	// goog.LOCALE is a special case as it's not processed as a define, it is
	// instead late-substituted.
	if processDefines_KNOWN_DEFINES.Contains(defineName) || defineName == "goog.LOCALE" {
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
	replacement := p.GetReplacementForDefineName(&define.defineName)
	if !replacement.IsNil() {
		return replacement
	}

	if p.IsGoogDefineCall(define.value) && define.value.GetChildCount() == 3 {
		// Return the second argument of goog.define('name', false);
		return define.value.GetChildAtIndex(2)
	}
	return rhino.Node{}
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
func (p *ProcessDefines) GetReplacementForDefineName(defineName *string) rhino.Node {
	// Go-only: language-adapter — a null key finds nothing in either map.
	if defineName == nil {
		return rhino.Node{}
	}
	replacementFromFlags := p.replacementValuesFromFlags.Get(*defineName)
	if !replacementFromFlags.IsNil() {
		return replacementFromFlags
	}

	replacementFromClosureDefines := p.replacementValuesFromClosureDefines.Get(*defineName)
	if !replacementFromClosureDefines.IsNil() {
		return replacementFromClosureDefines
	}

	return rhino.Node{}
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
			nameNode := value.GetSecondChild()
			defineName = nameNode.GetString()
		} else {
			defineName = name.GetFullName()
		}
		existingDefine :=
			p.defineByDefineName.PutIfAbsent(
				defineName, NewProcessDefines_Define(defineName, name, declaration, valueParent, value))

		if existingDefine != nil {
			declaration = existingDefine.declaration
			totalSets += existingDefine.name.GetTotalSets()
		}

		/*
		 * We have to report this here because otherwise we don't remember
		 * which names have the same define name. It's not worth it tracking a
		 * set of names, because it makes the rest of the pass more complex.
		 */
		if totalSets > 1 {
			for _, ref := range name.GetRefs() {
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
	defineNames := jutil.NewLinkedHashSet(processDefines_KNOWN_DEFINES.ToSlice()...)
	for k := range p.defineByDefineName.Keys() {
		defineNames.Add(k)
	}
	p.compiler.SetDefineNames(defineNames.ToSlice())
}

// Java: ProcessDefines.selectDefineDeclaration(GlobalNamespace.Name)
func (p *ProcessDefines) SelectDefineDeclaration(name *GlobalNamespace_Name) *GlobalNamespace_Ref {
	if name.GetGlobalSets() == 0 {
		return nil
	}
	for _, ref := range name.GetRefs() {
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

		p.knownDefineJsdocs.Add(jsdoc)
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
	} else if declarationParent.IsAssign() && declarationNode.IsFirstChildOf(declarationParent) {
		// Assignment. Must either assign to a qualified name, or else be a
		// different ref than the declaration to not emit an error (we don't
		// allow assignment before it's declared).
		return declarationParent
	}
	return rhino.Node{}
}

// Java: ProcessDefines.collectValidDefineValueExpressions()
func (p *ProcessDefines) CollectValidDefineValueExpressions() {
	if !p.mode.Check() || p.defineByDefineName.IsEmpty() {
		return
	}

	defineNamesSet := make(map[*GlobalNamespace_Name]struct{}, p.defineByDefineName.Size())
	// All defines are implicitly valid in the values of other defines.
	for define := range p.defineByDefineName.Values() {
		defineNamesSet[define.name] = struct{}{}
		for _, r := range define.name.GetRefs() {
			if !r.IsSet() {
				p.validDefineValueExpressions.Add(r.GetNode())
			}
		}
	}

	hasUnknownDefineValue := false
	for define := range p.defineByDefineName.Values() {
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
			for _, ref := range name.GetRefs() {
				p.validDefineValueExpressions.Add(ref.GetNode())
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
				for _, ref := range name.GetRefs() {
					p.validDefineValueExpressions.Add(ref.GetNode())
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

	for define := range p.defineByDefineName.Values() {
		declarationNode := define.declaration.GetNode()

		if !p.IsGoogDefineCall(define.value) &&
			!processDefines_KNOWN_DEFINES.Contains(define.defineName) &&
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
		if jsdoc != nil && jsdoc.IsDefine() && p.knownDefineJsdocs.Add(jsdoc) {
			p.compiler.Report(JSError_MakeWithNode(n, ProcessDefines_INVALID_DEFINE_LOCATION))
		}

		if p.IsGoogDefineCall(n) && p.knownGoogDefineCalls.Add(n) {
			p.VerifyGoogDefine(n)
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
		if first := n.GetFirstChild(); !first.IsNil() && first.MatchesName("CLOSURE_DEFINES") {
			valueNode := first.GetFirstChild()
			if !valueNode.IsNil() && valueNode.IsObjectLit() {
				p := c.outer
				for ch := valueNode.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
					p.HandleClosureDefinesValue(ch, ch.GetFirstChild(), n)
				}
			}
		}
	case rhino.Token_ASSIGN:
		for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
			c.walk(ch)
		}
		lhs := n.GetFirstChild()
		if lhs.IsGetElem() && lhs.GetFirstChild().MatchesName("CLOSURE_DEFINES") {
			c.outer.HandleClosureDefinesValue(lhs.GetSecondChild(), n.GetSecondChild(), n)
		}
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

	if NodeUtil_IsNameDeclaration(n) && n.GetFirstChild().MatchesName("CLOSURE_DEFINES") {
		// var CLOSURE_DEFINES = {...};
		valueNode := n.GetFirstFirstChild()
		if !valueNode.IsNil() && valueNode.IsObjectLit() {
			for c := valueNode.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
				p.HandleClosureDefinesValue(c, c.GetFirstChild(), n)
			}
		}
	} else if n.IsAssign() {
		// CLOSURE_DEFINES['...'] = ...;
		lhs := n.GetFirstChild()
		if lhs.IsGetElem() && lhs.GetFirstChild().MatchesName("CLOSURE_DEFINES") {
			p.HandleClosureDefinesValue(lhs.GetSecondChild(), n.GetSecondChild(), n)
		}
	}
}

// Java: ProcessDefines.handleClosureDefinesValue(Node, Node, Node)
func (p *ProcessDefines) HandleClosureDefinesValue(stringNode rhino.Node, valueNode rhino.Node, errorNode rhino.Node) {
	if (stringNode.IsStringKey() || stringNode.IsStringLit()) &&
		processDefines_isValidClosureDefinesValue(valueNode) {
		key := stringNode.GetString()
		if p.replacementValuesFromClosureDefines.ContainsKey(key) {
			p.compiler.Report(
				JSError_MakeWithNode(
					errorNode,
					ProcessDefines_CLOSURE_DEFINES_MULTIPLE,
					key,
					p.replacementValuesFromClosureDefines.Get(key).GetLocation()))
		}
		p.replacementValuesFromClosureDefines.Put(key, valueNode)
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
		return val.GetFirstChild().IsNumber()
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
		rhino.Token_BITNOT,
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
		return p.IsValidDefineValue(val.GetFirstChild()).And(p.IsValidDefineValue(val.GetLastChild()))
	case rhino.Token_HOOK:
		return p.IsValidDefineValue(val.GetFirstChild()).
			And(p.IsValidDefineValue(val.GetSecondChild())).
			And(p.IsValidDefineValue(val.GetLastChild()))
		// Unary operators are valid if the child is valid.
	case rhino.Token_NOT, rhino.Token_NEG, rhino.Token_POS:
		return p.IsValidDefineValue(val.GetFirstChild())
		// Names are valid if and only if they are defines themselves.
	case rhino.Token_NAME, rhino.Token_GETPROP:
		if val.IsQualifiedName() {
			if p.validDefineValueExpressions.Contains(val) {
				return base.Tri_TRUE
			}
			return base.Tri_UNKNOWN
		}
		// Allow goog.define('XYZ', <val>) calls if and only if <val> is valid.
	case rhino.Token_CALL:
		if !p.IsGoogDefineCall(val) {
			return base.Tri_FALSE
		}
		if !val.HasXChildren(3) {
			// goog.define call with wrong arg count. Warn elsewhere and treat
			// this call as valid.
			return base.Tri_TRUE
		}
		return p.IsValidDefineValue(val.GetChildAtIndex(2))
	default:
	}

	return base.Tri_FALSE
}

// processDefines_getConstantDeclValue checks whether the NAME node is
// inside either a CONST or a @const VAR. Returns the RHS node if so,
// otherwise returns null.
//
// Java: ProcessDefines.getConstantDeclValue(Node)
func processDefines_getConstantDeclValue(name rhino.Node) rhino.Node {
	parent := name.GetParent()
	if parent.IsNil() {
		return rhino.Node{}
	}
	if name.IsName() {
		if parent.IsConst() {
			return name.GetFirstChild()
		} else if !parent.IsVar() {
			return rhino.Node{}
		}
		jsdoc := NodeUtil_GetBestJSDocInfo(name)
		if jsdoc != nil && jsdoc.IsConstant() {
			return name.GetFirstChild()
		}
		return rhino.Node{}
	} else if name.IsGetProp() && parent.IsAssign() {
		jsdoc := NodeUtil_GetBestJSDocInfo(name)
		if jsdoc != nil && jsdoc.IsConstant() {
			return name.GetNext()
		}
		return rhino.Node{}
	}
	return rhino.Node{}
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

	// The connonical set ref with an `@define` or `goog.define`.
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

	if node.IsNil() || !node.IsCall() {
		return false
	}
	return node.GetFirstChild().MatchesQualifiedNameNode(processDefines_GOOG_DEFINE)
}

// VerifyGoogDefine verifies that a goog.define method call has exactly two
// arguments, with the first a string literal whose contents is a valid JS
// qualified name. Reports a compile error if it doesn't. Returns whether the
// argument checked out okay.
//
// Java: ProcessDefines.verifyGoogDefine(Node)
func (p *ProcessDefines) VerifyGoogDefine(callNode rhino.Node) bool {
	p.knownGoogDefineCalls.Add(callNode)

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
	if parent.IsExprResult() {
		infoNode = parent.GetFirstChild()
	}
	info := infoNode.GetJSDocInfo()
	if info == nil || !info.IsDefine() {
		p.compiler.Report(JSError_MakeWithNode(parent, ProcessDefines_MISSING_DEFINE_ANNOTATION))
		return false
	}
	return true
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
				jutil.ToJavaString(methodName.GetQualifiedName())))
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
				jutil.ToJavaString(methodName.GetQualifiedName())))
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
				jutil.ToJavaString(methodName.GetQualifiedName())))
		return false
	}
	return true
}
