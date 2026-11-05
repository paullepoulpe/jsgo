// Port of com/google/javascript/jscomp/ProcessDefines.java @ closure-compiler v20261028

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
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// Go-only: this pass uses the typed AST view (internal/ast) where it reads
// shapes: the default value of goog.define('name', value) (its second
// argument, googDefineDefault, where Java reads child 2 of the call), the
// operands of the operators a define value may use, and the CLOSURE_DEFINES
// declaration and assignments (shared by the Go-only walk and Visit). The
// goog.define checks stay on rhino: they pass arguments that may be missing
// to Java's verifyNotNull and verifyIsLast.

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
	enableZonesDefineName      base.OptString
	zoneInputPattern           *jutil.Pattern
	unknownDefinesToIgnore     map[string]struct{}

	knownDefineJsdocs    map[*rhino.JSDocInfo]struct{}
	knownGoogDefineCalls map[rhino.Node]struct{}
	defineByDefineName   map[string]*ProcessDefines_Define
	defines              []*ProcessDefines_Define
	// from var CLOSURE_DEFINES = {
	replacementValuesFromClosureDefines OrderedStringMap[rhino.Node]
	validDefineValueExpressions         map[rhino.Node]struct{}
	hasZoneInput                        bool

	namespace *GlobalNamespace
}

var _ CompilerPass = (*ProcessDefines)(nil)

// Defines in this set will not be flagged with "unknown define" warnings.
// There are flags that always set these defines, even when they might not be
// in the binary.
var knownProcessDefinesList = [...]string{
	"COMPILED",
	"goog.DEBUG",
	"$jscomp.ASSUME_ES5",
	"$jscomp.ASSUME_ES6",
	"$jscomp.ASSUME_ES2020",
	"$jscomp.ISOLATE_POLYFILLS",
	"$jscomp.INSTRUMENT_ASYNC_CONTEXT",
}

var knownProcessDefines = func() map[string]struct{} {
	m := make(map[string]struct{}, len(knownProcessDefinesList))
	for _, d := range knownProcessDefinesList {
		m[d] = struct{}{}
	}
	return m
}()

var processDefines_GOOG_DEFINE = rhino.QualifiedName_Of("goog.define")

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
	EnableZonesDefineName   base.OptString
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
		enableZonesDefineName:   opts.EnableZonesDefineName,
		zoneInputPattern:        opts.ZoneInputPattern,
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
	processDefinesModeNames    = [...]string{"CHECK", "OPTIMIZE", "CHECK_AND_OPTIMIZE"}
	processDefinesModeCheck    = [...]bool{true, false, true}
	processDefinesModeOptimize = [...]bool{false, true, true}
)

// Java: ProcessDefines.Mode.name()
func (m ProcessDefines_Mode) String() string { return processDefinesModeNames[m] }

// Check is the Mode.check field.
func (m ProcessDefines_Mode) Check() bool { return processDefinesModeCheck[m] }

// Optimize is the Mode.optimize field.
func (m ProcessDefines_Mode) Optimize() bool { return processDefinesModeOptimize[m] }

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
				if !defineLhs.IsName() && !defineLhs.IsQualifiedName() {
					panic(jutil.NewIllegalStateExceptionMsg(defineLhs.String()))
				}
				globalDefineVar :=
					rhino.IR_VarWithValue(globalDefineName, defineLhs).SrcrefTreeIfMissing(define.valueParent)
				globalDefineVar.InsertAfter(define.valueParent.GetParent())
				p.compiler.ReportChangeToEnclosingScope(define.valueParent.GetParent())
			}
		}
	}

	if !p.replacementValuesFromFlags.IsEmpty() || !p.replacementValuesFromClosureDefines.IsEmpty() {
		// Sets.difference(Sets.union(flags, closureDefines),
		//     Sets.union(Sets.union(KNOWN_DEFINES, defineByDefineName), ignore))
		isUsed := func(k string) bool {
			if _, known := knownProcessDefines[k]; known {
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
		for _, e := range p.replacementValuesFromClosureDefines.Entries() {
			if !p.replacementValuesFromFlags.ContainsKey(e.Key) && !isUsed(e.Key) {
				p.compiler.Report(JSError_Make(ProcessDefines_UNKNOWN_DEFINE_WARNING, e.Key))
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
	if _, known := knownProcessDefines[defineName]; known || defineName == "goog.LOCALE" {
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
		if defaultValue, ok := googDefineDefault(ast.MustCall(define.value)); ok {
			return defaultValue.Node()
		}
	}
	return rhino.NilNode
}

// googDefineDefault returns the default value of goog.define('name', value),
// its second argument; ok is false if the call does not have exactly two.
//
// Go-only: Java's getChildCount() == 3 and getChildAtIndex(2) (the callee is
// the first child).
func googDefineDefault(call ast.Call) (ast.Arg, bool) {
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
	if replacementFromClosureDefines, ok := p.replacementValuesFromClosureDefines.Lookup(defineName); ok && !replacementFromClosureDefines.IsNil() {
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
	defineName, ok := p.enableZonesDefineName.Get()
	if !ok {
		return enableZonesDefault
	}
	zoneEnabled := p.GetReplacementForDefineName(defineName)
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

type gnDefineScript struct {
	node   rhino.Node
	isCode bool
}

const (
	gnScriptTraverse uint8 = 1 << iota
	gnScriptDefineCand
	gnScriptFuncCand
)

func mayBeGlobalSetInFunc(n, parent rhino.Node) bool {
	if parent.IsNil() || parent.IsGetProp() {
		return false
	}
	switch ast.RoleIn(n, parent) {
	case ast.RoleAssignTarget,
		ast.RoleAssignOpTarget,
		ast.RoleIncDecTarget,
		ast.RoleForInTarget,
		ast.RoleForOfTarget,
		ast.RoleForAwaitOfTarget:
		return true
	case ast.RoleStringKeyValue,
		ast.RoleArrayPatternElements,
		ast.RoleDefaultValueTarget,
		ast.RoleDefaultValueValue,
		ast.RoleComputedPropKey,
		ast.RoleComputedPropValue,
		ast.RoleIterRestTarget,
		ast.RoleObjectRestTarget:
		return n.IsName() && NodeUtil_IsLhsByDestructuring(n)
	}
	return false
}

// prepareDefineNamespace populates defineCandidateLeafIDs (and, in check mode,
// defineCandidateScripts) and identifies the subset of SCRIPT roots that can
// contribute symbols or references to ProcessDefines.
//
// ProcessDefines inspects only global names whose global-scope declaration is
// annotated @define (SelectDefineDeclaration) or, when a @define initializer
// references another identifier, global constants declared at global hoist
// scope (includeGlobalConsts). Neither kind of declaration can reside inside a
// FUNCTION body or MODULE_BODY, so Pass 1 scans only top-level AST nodes.
// Pass 2 then scans scripts in parallel to skip any script that contains none
// of the candidate leaf IDs.
//
// Go-only: performance.
func (g *GlobalNamespace) prepareDefineNamespace() (skipTraversal bool, traverseScripts map[rhino.Node]struct{}) {
	preSeeded := g.includeGlobalConsts && len(g.defineCandidateLeafIDs) > 0
	if !preSeeded {
		g.defineCandidateLeafIDs = make(map[uint32]struct{}, 64)
	}
	needCheckCand := g.processDefinesCheck && !g.includeGlobalConsts
	scripts, ok := g.collectDefineScripts()
	if !ok {
		if !preSeeded {
			var stackBuf [64]rhino.Node
			var idBuf [64]uint32
			stack := stackBuf[:0]
			ids := idBuf[:0]
			if g.HasExternsRoot() {
				ids, stack = collectScriptDefineLeafIDs(g.externsRoot, g.includeGlobalConsts, ids, stack)
			}
			ids, _ = collectScriptDefineLeafIDs(g.root, g.includeGlobalConsts, ids, stack)
			for _, id := range ids {
				g.addDefineCandidateLeafID(id)
			}
		}
		return !needCheckCand && len(g.defineCandidateLeafIDs) == 0, nil
	}
	if len(scripts) == 0 {
		return true, nil
	}

	numWorkers := parallelThreadCount(g.compiler, len(scripts))
	if !preSeeded {
		if numWorkers > 1 && len(scripts) >= 8 {
			perWorker := make([][]uint32, numWorkers)
			if p := spawnIndexedWorkers(len(scripts), numWorkers, func(wIdx int, claim func() (int, bool)) {
				var ids []uint32
				var stack []rhino.Node
				for {
					idx, ok := claim()
					if !ok {
						break
					}
					ids, stack = collectScriptDefineLeafIDs(scripts[idx].node, g.includeGlobalConsts, ids, stack)
				}
				perWorker[wIdx] = ids
			}); p != nil {
				panic(p)
			}
			for _, ids := range perWorker {
				for _, id := range ids {
					g.addDefineCandidateLeafID(id)
				}
			}
		} else {
			var stackBuf [64]rhino.Node
			var idBuf [64]uint32
			stack := stackBuf[:0]
			ids := idBuf[:0]
			for i := range scripts {
				ids, stack = collectScriptDefineLeafIDs(scripts[i].node, g.includeGlobalConsts, ids[:0], stack)
				for _, id := range ids {
					g.addDefineCandidateLeafID(id)
				}
			}
		}
	}

	hasLeafCandidates := len(g.defineCandidateLeafIDs) > 0
	if !needCheckCand && !hasLeafCandidates {
		return true, nil
	}

	var bloom [256]uint64
	if hasLeafCandidates {
		// order-independent: bitwise OR of per-ID bloom filter bits is commutative.
		for id := range g.defineCandidateLeafIDs {
			bit := (id * 0x9e3779b1) >> 18
			bloom[bit>>6] |= uint64(1) << (bit & 63)
		}
	}
	matchesLeafID := func(id uint32, checkDots bool) bool {
		if id <= strpool.EmptyID {
			return false
		}
		bit := (id * 0x9e3779b1) >> 18
		if bloom[bit>>6]&(uint64(1)<<(bit&63)) != 0 {
			if _, ok := g.defineCandidateLeafIDs[id]; ok {
				return true
			}
		}
		if checkDots {
			if s := strpool.String(strpool.ID(id)); strings.IndexByte(s, '.') >= 0 {
				last := s[strings.LastIndexByte(s, '.')+1:]
				if lid := uint32(strpool.Lookup(last)); lid > strpool.EmptyID {
					lbit := (lid * 0x9e3779b1) >> 18
					if bloom[lbit>>6]&(uint64(1)<<(lbit&63)) != 0 {
						_, ok := g.defineCandidateLeafIDs[lid]
						return ok
					}
				}
			}
		}
		return false
	}
	provideMatchesLeafIDs := func(ns string) bool {
		for {
			seg, rest, hasMore := strings.Cut(ns, ".")
			if lid := uint32(strpool.Lookup(seg)); lid > strpool.EmptyID {
				bit := (lid * 0x9e3779b1) >> 18
				if bloom[bit>>6]&(uint64(1)<<(bit&63)) != 0 {
					if _, ok := g.defineCandidateLeafIDs[lid]; ok {
						return true
					}
				}
			}
			if !hasMore {
				return false
			}
			ns = rest
		}
	}
	const (
		popFuncSentinel   rhino.Node = 0
		popModSentinel    rhino.Node = ^rhino.Node(0)
		enterFuncSentinel rhino.Node = ^rhino.Node(1)
	)
	scanScript := func(s gnDefineScript, stack []rhino.Node) (uint8, []rhino.Node) {
		checkCand := needCheckCand && s.isCode
		if !hasLeafCandidates && !checkCand {
			return 0, stack
		}
		needsTraverse := false
		hasDefineCand := false
		hasFuncCand := false
		funcDepth := 0
		modDepth := 0
		stack = append(stack[:0], s.node)
		for len(stack) > 0 {
			n := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if n == enterFuncSentinel {
				funcDepth++
				continue
			}
			if n == popFuncSentinel {
				funcDepth--
				continue
			}
			if n == popModSentinel {
				modDepth--
				continue
			}
			tok := n.GetToken()
			inLocalFunc := funcDepth > 0 && modDepth == 0
			if checkCand && !hasDefineCand {
				if n.HasJSDocInfo() {
					if jsdoc := n.GetJSDocInfo(); jsdoc != nil && jsdoc.IsDefine() {
						hasDefineCand = true
					}
				}
				if !hasDefineCand {
					if tok == rhino.Token_CALL && processDefines_GOOG_DEFINE.Matches(ast.MustCall(n).Callee().Node()) {
						hasDefineCand = true
					} else if tok == rhino.Token_NAME && n.GetStringID() == gnClosureDefinesStrID {
						hasDefineCand = true
					}
				}
			}
			if hasLeafCandidates && (!needsTraverse || (inLocalFunc && !hasFuncCand)) {
				matched := false
				switch tok {
				case rhino.Token_NAME, rhino.Token_GETPROP:
					parent := n.GetParent()
					if !parent.IsGetProp() && (!inLocalFunc || mayBeGlobalSetInFunc(n, parent)) {
						matched = matchesLeafID(n.GetStringID(), true)
					}
				case rhino.Token_STRING_KEY,
					rhino.Token_GETTER_DEF,
					rhino.Token_SETTER_DEF,
					rhino.Token_MEMBER_FUNCTION_DEF,
					rhino.Token_MEMBER_FIELD_DEF:
					matched = matchesLeafID(n.GetStringID(), false)
				case rhino.Token_CALL:
					parent := n.GetParent()
					if parent.IsExprResult() {
						call := ast.MustCall(n)
						if GlobalNamespace_GOOG_PROVIDE.Matches(call.Callee().Node()) {
							if arg, ok := call.FirstArg(); ok {
								if str, ok := arg.AsStringLit(); ok {
									if !parent.GetParent().IsScript() || provideMatchesLeafIDs(str.Value()) {
										matched = true
									}
								}
							}
						}
					}
				}
				if matched {
					needsTraverse = true
					if inLocalFunc {
						hasFuncCand = true
					}
				}
			}
			if (!hasLeafCandidates || (needsTraverse && hasFuncCand)) && (!checkCand || hasDefineCand) {
				break
			}
			if tok == rhino.Token_MODULE_BODY {
				modDepth++
				stack = append(stack, popModSentinel)
				for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
					stack = append(stack, c)
				}
				continue
			}
			if tok == rhino.Token_FUNCTION {
				isBundledMod := NodeUtil_IsBundledGoogModuleCall(n.GetParent())
				fnName := n.GetFirstChild()
				if isBundledMod {
					modDepth++
					stack = append(stack, popModSentinel)
					for c := fnName; !c.IsNil(); c = c.GetNext() {
						stack = append(stack, c)
					}
				} else {
					hasBodyChildren := false
					for c := fnName.GetNext(); !c.IsNil(); c = c.GetNext() {
						if !hasBodyChildren {
							hasBodyChildren = true
							stack = append(stack, popFuncSentinel)
						}
						stack = append(stack, c)
					}
					if hasBodyChildren {
						stack = append(stack, enterFuncSentinel)
					}
					if !fnName.IsNil() {
						stack = append(stack, fnName)
					}
				}
				continue
			}
			for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
				stack = append(stack, c)
			}
		}
		var f uint8
		if needsTraverse {
			f |= gnScriptTraverse
		}
		if hasDefineCand {
			f |= gnScriptDefineCand
		}
		if hasFuncCand {
			f |= gnScriptFuncCand
		}
		return f, stack
	}

	flags := make([]uint8, len(scripts))
	if numWorkers > 1 && len(scripts) >= 8 {
		if p := spawnIndexedWorkers(len(scripts), numWorkers, func(_ int, claim func() (int, bool)) {
			var stack []rhino.Node
			for {
				idx, ok := claim()
				if !ok {
					return
				}
				flags[idx], stack = scanScript(scripts[idx], stack)
			}
		}); p != nil {
			panic(p)
		}
	} else {
		var stackBuf [128]rhino.Node
		stack := stackBuf[:0]
		for i := range scripts {
			flags[i], stack = scanScript(scripts[i], stack)
		}
	}

	g.defineCandScriptsPopulated = true
	traverseCount := 0
	for i, f := range flags {
		if f&gnScriptTraverse != 0 {
			traverseCount++
		}
		if f&gnScriptFuncCand != 0 {
			g.hasFuncDefineCand = true
		}
		if f&gnScriptDefineCand != 0 {
			if g.defineCandidateScripts == nil {
				g.defineCandidateScripts = make(map[rhino.Node]struct{}, 8)
			}
			g.defineCandidateScripts[scripts[i].node] = struct{}{}
		}
	}
	if traverseCount == 0 {
		return true, nil
	}
	if traverseCount == len(scripts) && g.shouldTraverseScript == nil {
		return false, nil
	}
	traverseScripts = make(map[rhino.Node]struct{}, traverseCount)
	for i, f := range flags {
		if f&gnScriptTraverse != 0 {
			traverseScripts[scripts[i].node] = struct{}{}
		}
	}
	return false, traverseScripts
}

func (g *GlobalNamespace) collectDefineScripts() ([]gnDefineScript, bool) {
	rootCount, ok := countRootScripts(g.root)
	if !ok {
		return nil, false
	}
	extCount := 0
	if g.HasExternsRoot() {
		var okExt bool
		extCount, okExt = countRootScripts(g.externsRoot)
		if !okExt || g.externsRoot.GetParent().IsNil() || g.externsRoot.GetParent() != g.root.GetParent() {
			return nil, false
		}
	}
	if extCount+rootCount == 0 {
		return nil, true
	}
	scripts := make([]gnDefineScript, 0, extCount+rootCount)
	pred := g.shouldTraverseScript
	if g.HasExternsRoot() {
		for s := g.externsRoot.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
			if pred != nil && !pred(s) {
				continue
			}
			scripts = append(scripts, gnDefineScript{node: s, isCode: !s.IsFromExterns()})
		}
	}
	for s := g.root.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
		if pred != nil && !pred(s) {
			continue
		}
		scripts = append(scripts, gnDefineScript{node: s, isCode: !g.HasExternsRoot() || !s.IsFromExterns()})
	}
	return scripts, true
}

func (g *GlobalNamespace) addDefineCandidateLeafID(id uint32) {
	if id <= strpool.EmptyID {
		return
	}
	g.defineCandidateLeafIDs[id] = struct{}{}
	if s := strpool.String(strpool.ID(id)); strings.IndexByte(s, '.') >= 0 {
		if last := s[strings.LastIndexByte(s, '.')+1:]; last != "" {
			g.defineCandidateLeafIDs[uint32(strpool.Intern(last))] = struct{}{}
		}
	}
}

func collectScriptDefineLeafIDs(root rhino.Node, includeConsts bool, out []uint32, stack []rhino.Node) ([]uint32, []rhino.Node) {
	if root.IsNil() {
		return out, stack
	}
	stack = append(stack[:0], root)
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch n.GetToken() {
		case rhino.Token_MODULE_BODY:
			continue
		case rhino.Token_FUNCTION:
			if !NodeUtil_IsFunctionExpression(n) {
				if name := n.GetFirstChild(); name.IsName() {
					stack = append(stack, name)
				}
			}
			continue
		case rhino.Token_NAME:
			for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
				stack = append(stack, c)
			}
			parent := n.GetParent()
			if parent.IsGetProp() {
				continue
			}
			if includeConsts && !parent.IsNil() && parent.IsConst() {
				out = append(out, n.GetStringID())
				continue
			}
			if info := NodeUtil_GetBestJSDocInfo(n); info != nil && (info.IsDefine() || (includeConsts && info.IsConstant())) {
				out = append(out, n.GetStringID())
			}
			continue
		case rhino.Token_GETPROP:
			for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
				stack = append(stack, c)
			}
			if n.GetParent().IsGetProp() {
				continue
			}
			if info := NodeUtil_GetBestJSDocInfo(n); info != nil && (info.IsDefine() || (includeConsts && info.IsConstant())) && n.IsQualifiedName() {
				out = append(out, n.GetStringID())
			}
			continue
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			stack = append(stack, c)
		}
	}
	return out, stack
}
