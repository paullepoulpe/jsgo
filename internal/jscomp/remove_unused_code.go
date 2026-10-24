// Port of com/google/javascript/jscomp/RemoveUnusedCode.java @ closure-compiler v20261015

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

package jscomp

import (
	"reflect"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type propRemovableList struct {
	head *RemoveUnusedCode_RemovableBase
	tail *RemoveUnusedCode_RemovableBase
}

// rucVarHash hashes the identity of a variable: RemoveUnusedCode re-creates
// the scopes of a function, so one variable has several Var objects (Java:
// Var.equals compares the name and the scope root).
func rucVarHash(scopeRoot rhino.Node, nameID uint32) uint32 {
	h := uint32(31*int32(31+nameID) + scopeRoot.IdentityHashCode())
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	return h
}

// removeUnusedCodeState is the state of one run of the pass; Process drops
// it when the run ends.
type removeUnusedCodeState struct {
	worklist     []RemoveUnusedCode_Continuation
	worklistHead int
	worklistSize int

	// Java: varInfoMap. varBuckets is a hash table of the VarInfos keyed by
	// (scope root, name), chained through RealVarInfo.hashNext; varInfos
	// holds them in insertion order.
	//
	// Go-only: performance — a Go map with that key allocated several times
	// more (up to 60 MB in the bundle benchmark).
	varBuckets []*RemoveUnusedCode_RealVarInfo
	varInfos   []*RemoveUnusedCode_RealVarInfo

	pinnedPropertyNames map[string]struct{}
	propRemovables      map[string]propRemovableList
	propRemovableKeys   []string

	allFunctionParamScopes []*Scope

	extraContSlab []RemoveUnusedCode_Continuation

	// Go-only: performance — one run's removables and var infos are created
	// together and die together with this state (see jutil.Slab).
	removables  jutil.Slab[RemoveUnusedCode_RemovableBase]
	varInfoSlab jutil.Slab[RemoveUnusedCode_RealVarInfo]
}

// RemoveUnusedCode is garbage collection for variable and function
// definitions. Basically performs a mark-and-sweep type algorithm over the
// JavaScript parse tree.
//
// For each scope: (1) Scan the variable/function declarations at that scope.
// (2) Traverse the scope for references, marking all referenced variables.
// Unlike other compiler passes, this is a pre-order traversal, not a
// post-order traversal. (3) If the traversal encounters an assign without
// other side-effects, create a continuation. Continue the continuation iff
// the assigned variable is referenced. (4) When the traversal completes,
// remove all unreferenced variables.
//
// If it makes it easier, you can think of the continuations of the traversal
// as a reference graph. Each continuation represents a set of edges, where
// the source node is a known variable, and the destination nodes are lazily
// evaluated when the continuation is executed.
//
// This algorithm is similar to the algorithm used by SmartNameRemoval.
// SmartNameRemoval maintains an explicit graph of dependencies between
// global symbols. However, SmartNameRemoval cannot handle non-trivial edges
// in the reference graph ("A is referenced iff both B and C are
// referenced"), or local variables. SmartNameRemoval is also substantially
// more complicated because it tries to handle namespaces (which is largely
// unnecessary in the presence of CollapseProperties.
//
// This pass also uses a more complex analysis of assignments, where an
// assignment to a variable or a property of that variable does not
// necessarily count as a reference to that variable, unless we can prove
// that it modifies external state. This is similar to
// FlowSensitiveInlineVariables, except that it works for variables used
// across scopes.
//
// Multiple datastructures are used to accumulate nodes, some of which are
// later removed. Since some nodes encompass a subtree of nodes, the removal
// can sometimes pre-remove other nodes which are also referenced in these
// datastructures for later removal. Attempting double-removal violates
// scope change notification constraints so there is a desire to excise
// already-removed subtree nodes from these datastructures. But not all of
// the datastructures are conducive to flexible removal and the ones that are
// conducive don't necessarily track all flavors of nodes. So instead of
// updating datastructures on the fly a pre-check is performed to skip
// already-removed nodes right before the moment an attempt to remove them
// would otherwise be made.
type RemoveUnusedCode struct {
	compiler    AbstractCompiler
	astAnalyzer *AstAnalyzer

	codingConvention CodingConvention

	removeLocalVars bool
	removeGlobals   bool

	preserveFunctionExpressionNames bool

	externProperties    *jutil.LinkedHashSet[string]
	accessorSummary     *AccessorSummary
	lastAccessorSummary *AccessorSummary
	hasGetterAccessors  bool
	hasSetterAccessors  bool
	getterLenMask       [256]uint64
	getterEndMask       [256]uint64
	setterLenMask       [256]uint64
	setterEndMask       [256]uint64

	// Single value to use for all vars for which we cannot remove anything at
	// all.
	canonicalUnremovableVarInfo    *RemoveUnusedCode_RealVarInfo
	canonicalUnremovableVarInfoVal RemoveUnusedCode_RealVarInfo

	// Stores the names of all "leaf" properties that are polyfilled, to avoid
	// unnecessary qualified name matching and searches for all the other
	// properties. This includes global names such as "Promise" and "Map",
	// static methods on global names such as "Array.from" and "Math.fround",
	// and instance properties such as "String.prototype.repeat" and
	// "Promise.prototype.finally".
	polyfills       *jutil.HashMultimap[string, RemoveUnusedCode_PolyfillInfo]
	hasPolyfills    bool
	polyfillLenMask [256]uint64
	polyfillEndMask [256]uint64

	guardedUsages         *jutil.LinkedHashSet[rhino.Node]
	guardedUsagesComputed bool
	root                  rhino.Node

	polyfillsFromTable *PolyfillUsageFinder_Polyfills

	scopeCreator *SyntacticScopeCreator

	removeUnusedPrototypeProperties               bool
	removeUnusedThisProperties                    bool
	removeUnusedObjectDefinePropertiesDefinitions bool
	removeUnusedPolyfills                         bool
	assumeGettersArePure                          bool

	// Allocated & cleaned up by process()
	removalLog     diagnostic.LogFile
	unremovableLog diagnostic.LogFile

	removeUnusedCodeState
}

// Properties that are implicitly used as part of the JS language.
var RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES = []string{"length", "toString", "valueOf", "constructor", "prototype"}

func isImplicitlyUsedProperty(name string) bool {
	switch name {
	case "length", "toString", "valueOf", "constructor", "prototype":
		return true
	default:
		return false
	}
}

// Java: RemoveUnusedCode(Builder)
func NewRemoveUnusedCode(builder *RemoveUnusedCode_Builder) *RemoveUnusedCode {
	sc := NewSyntacticScopeCreator(builder.compiler)
	r := &RemoveUnusedCode{
		compiler:                        builder.compiler,
		astAnalyzer:                     builder.compiler.GetAstAnalyzer(),
		codingConvention:                builder.compiler.GetCodingConvention(),
		scopeCreator:                    sc,
		removeLocalVars:                 builder.removeLocalVars,
		removeGlobals:                   builder.removeGlobals,
		preserveFunctionExpressionNames: builder.preserveFunctionExpressionNames,
		removeUnusedPrototypeProperties: builder.removeUnusedPrototypeProperties,
		removeUnusedThisProperties:      builder.removeUnusedThisProperties,
		removeUnusedObjectDefinePropertiesDefinitions: builder.removeUnusedObjectDefinePropertiesDefinitions,
		removeUnusedPolyfills:                         builder.removeUnusedPolyfills,
		polyfillsFromTable:                            defaultPolyfillsTable(),
		assumeGettersArePure:                          builder.assumeGettersArePure,
	}

	// All Vars that are completely unremovable will share this VarInfo instance.
	r.canonicalUnremovableVarInfoVal.outer = r
	r.canonicalUnremovableVarInfoVal.flags = rucVarFlagCanonical
	r.canonicalUnremovableVarInfo = &r.canonicalUnremovableVarInfoVal
	return r
}

// RemoveUnusedCode_Builder is the nested class Builder.
type RemoveUnusedCode_Builder struct {
	compiler AbstractCompiler

	removeLocalVars                               bool
	removeGlobals                                 bool
	preserveFunctionExpressionNames               bool
	removeUnusedPrototypeProperties               bool
	removeUnusedThisProperties                    bool
	removeUnusedObjectDefinePropertiesDefinitions bool
	removeUnusedPolyfills                         bool
	assumeGettersArePure                          bool
}

// Java: RemoveUnusedCode.Builder(AbstractCompiler)
func NewRemoveUnusedCode_Builder(compiler AbstractCompiler) *RemoveUnusedCode_Builder {
	return &RemoveUnusedCode_Builder{compiler: compiler}
}

// Java: RemoveUnusedCode.Builder.removeLocalVars(boolean)
func (b *RemoveUnusedCode_Builder) RemoveLocalVars(value bool) *RemoveUnusedCode_Builder {
	b.removeLocalVars = value
	return b
}

// Java: RemoveUnusedCode.Builder.removeGlobals(boolean)
func (b *RemoveUnusedCode_Builder) RemoveGlobals(value bool) *RemoveUnusedCode_Builder {
	b.removeGlobals = value
	return b
}

// Java: RemoveUnusedCode.Builder.preserveFunctionExpressionNames(boolean)
func (b *RemoveUnusedCode_Builder) PreserveFunctionExpressionNames(value bool) *RemoveUnusedCode_Builder {
	b.preserveFunctionExpressionNames = value
	return b
}

// Java: RemoveUnusedCode.Builder.removeUnusedPrototypeProperties(boolean)
func (b *RemoveUnusedCode_Builder) RemoveUnusedPrototypeProperties(value bool) *RemoveUnusedCode_Builder {
	b.removeUnusedPrototypeProperties = value
	return b
}

// Java: RemoveUnusedCode.Builder.removeUnusedThisProperties(boolean)
func (b *RemoveUnusedCode_Builder) RemoveUnusedThisProperties(value bool) *RemoveUnusedCode_Builder {
	b.removeUnusedThisProperties = value
	return b
}

// Java: RemoveUnusedCode.Builder.removeUnusedObjectDefinePropertiesDefinitions(boolean)
func (b *RemoveUnusedCode_Builder) RemoveUnusedObjectDefinePropertiesDefinitions(value bool) *RemoveUnusedCode_Builder {
	b.removeUnusedObjectDefinePropertiesDefinitions = value
	return b
}

// Java: RemoveUnusedCode.Builder.removeUnusedPolyfills(boolean)
func (b *RemoveUnusedCode_Builder) RemoveUnusedPolyfills(value bool) *RemoveUnusedCode_Builder {
	b.removeUnusedPolyfills = value
	return b
}

// Java: RemoveUnusedCode.Builder.assumeGettersArePure(boolean)
func (b *RemoveUnusedCode_Builder) AssumeGettersArePure(value bool) *RemoveUnusedCode_Builder {
	b.assumeGettersArePure = value
	return b
}

// Java: RemoveUnusedCode.Builder.build()
func (b *RemoveUnusedCode_Builder) Build() *RemoveUnusedCode {
	return NewRemoveUnusedCode(b)
}

// RemoveUnusedCode_RemovalLogRecord supplies the string needed for an entry
// in the removal log.
type RemoveUnusedCode_RemovalLogRecord struct {
	kind                 string
	nameSupplier         func() string
	functionNameSupplier func() string
}

// Get returns a log entry string.
//
// Each entry is one tab-separated line of the form:
//
//	KIND NAME [FUNCTION_NAME]
//
// See specific methods below for details.
//
// Java: RemoveUnusedCode.RemovalLogRecord.get()
func (r *RemoveUnusedCode_RemovalLogRecord) Get() string {
	return strings.Join([]string{r.kind, r.nameSupplier(), r.functionNameSupplier()}, "\t")
}

// Java: RemoveUnusedCode.RemovalLogRecord(String, Supplier<String>, Supplier<String>)
func NewRemoveUnusedCode_RemovalLogRecordWithFunctionName(
	kind string, nameSupplier func() string, functionNameSupplier func() string) *RemoveUnusedCode_RemovalLogRecord {
	return &RemoveUnusedCode_RemovalLogRecord{
		kind:                 kind, // checkNotNull: a Go string is never null.
		nameSupplier:         jutil.CheckNotNull(nameSupplier),
		functionNameSupplier: jutil.CheckNotNull(functionNameSupplier),
	}
}

// Java: RemoveUnusedCode.RemovalLogRecord(String, Supplier<String>)
func NewRemoveUnusedCode_RemovalLogRecord(kind string, nameSupplier func() string) *RemoveUnusedCode_RemovalLogRecord {
	// No function name
	return NewRemoveUnusedCode_RemovalLogRecordWithFunctionName(kind, nameSupplier, func() string { return "" })
}

// Java: RemoveUnusedCode.RemovalLogRecord.forProperty(String)
func RemoveUnusedCode_RemovalLogRecord_ForProperty(propName string) *RemoveUnusedCode_RemovalLogRecord {
	return NewRemoveUnusedCode_RemovalLogRecord("prop", func() string { return propName })
}

// Java: RemoveUnusedCode.RemovalLogRecord.forVar(Var)
func RemoveUnusedCode_RemovalLogRecord_ForVar(var_ *Var) *RemoveUnusedCode_RemovalLogRecord {
	return NewRemoveUnusedCode_RemovalLogRecord("var", var_.GetName)
}

// Java: RemoveUnusedCode.RemovalLogRecord.forPolyfill(PolyfillInfo)
func RemoveUnusedCode_RemovalLogRecord_ForPolyfill(polyfillInfo RemoveUnusedCode_PolyfillInfo) *RemoveUnusedCode_RemovalLogRecord {
	return NewRemoveUnusedCode_RemovalLogRecord("poly", polyfillInfo.GetName)
}

// RemoveUnusedCode_RemovalLogRecord_ForNamedArg records removal of a named
// function parameter.
//
// nameNode is the parameter's NAME node; argList the function's PARAM_LIST
// node.
//
// Java: RemoveUnusedCode.RemovalLogRecord.forNamedArg(Node, Node)
func RemoveUnusedCode_RemovalLogRecord_ForNamedArg(nameNode rhino.Node, argList rhino.Node) *RemoveUnusedCode_RemovalLogRecord {
	return NewRemoveUnusedCode_RemovalLogRecordWithFunctionName(
		"arg", nameNode.GetString, RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier(argList))
}

// RemoveUnusedCode_RemovalLogRecord_ForDestructuringArg records removal of a
// destructuring function parameter.
//
// argList is the function's PARAM_LIST node.
//
// Java: RemoveUnusedCode.RemovalLogRecord.forDestructuringArg(Node)
func RemoveUnusedCode_RemovalLogRecord_ForDestructuringArg(argList rhino.Node) *RemoveUnusedCode_RemovalLogRecord {
	return NewRemoveUnusedCode_RemovalLogRecordWithFunctionName(
		"arg", func() string { return "<pattern>" }, RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier(argList))
}

// RemoveUnusedCode_RemovalLogRecord_ForMarkingNamedArg records that a named
// parameter is marked as unused for possible removal by OptimizeParameters.
//
// nameNode is the parameter's NAME node; argList the function's PARAM_LIST
// node.
//
// Java: RemoveUnusedCode.RemovalLogRecord.forMarkingNamedArg(Node, Node)
func RemoveUnusedCode_RemovalLogRecord_ForMarkingNamedArg(nameNode rhino.Node, argList rhino.Node) *RemoveUnusedCode_RemovalLogRecord {
	return NewRemoveUnusedCode_RemovalLogRecordWithFunctionName(
		"argmark", nameNode.GetString, RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier(argList))
}

// RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier returns
// a supplier for the FUNCTION_NAME field of an argument removal log entry.
//
// If no good name can be found, then "<anonymous>" will be supplied.
//
// argList is the function's PARAM_LIST node.
//
// Java: RemoveUnusedCode.RemovalLogRecord.getLoggableFunctionNameSupplier(Node)
func RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier(argList rhino.Node) func() string {
	return func() string {
		functionName := NodeUtil_GetNearestFunctionName(rhino.CheckNotNull(argList).GetParent())
		if functionName == nil {
			return "<anonymous>"
		}
		return *functionName
	}
}

// Process traverses the root, removing all unused variables. Multiple
// traversals may occur to ensure all unused variables are removed.
//
// Java: RemoveUnusedCode.process(Node, Node)
func (r *RemoveUnusedCode) Process(externs rhino.Node, root rhino.Node) {
	if !r.compiler.GetLifeCycleStage().IsNormalized() {
		jutil.CheckState(false)
	}
	r.externProperties = r.compiler.GetExternProperties()
	r.accessorSummary = r.compiler.GetAccessorSummary()

	defer func() {
		r.removeUnusedCodeState = removeUnusedCodeState{}
	}()

	removalLogFile :=
		r.compiler.CreateOrReopenIndexedLog(reflect.TypeFor[RemoveUnusedCode](), "removals.log")
	defer removalLogFile.Close()
	keepLogFile :=
		r.compiler.CreateOrReopenIndexedLog(reflect.TypeFor[RemoveUnusedCode](), "unremovable.log")
	defer keepLogFile.Close()
	defer func() {
		r.removalLog = nil
		r.unremovableLog = nil
	}()
	r.removalLog = removalLogFile // avoid passing the log file through a bunch of methods
	r.unremovableLog = keepLogFile
	r.TraverseAndRemoveUnusedReferences(root)
}

func (r *RemoveUnusedCode) enqueueContinuation(c RemoveUnusedCode_Continuation) {
	if r.worklistSize == len(r.worklist) {
		r.growWorklist()
	}
	idx := r.worklistHead + r.worklistSize
	if idx >= len(r.worklist) {
		idx -= len(r.worklist)
	}
	r.worklist[idx] = c
	r.worklistSize++
}

func (r *RemoveUnusedCode) growWorklist() {
	n := len(r.worklist) * 2
	if n < 256 {
		n = 256
	}
	nb := make([]RemoveUnusedCode_Continuation, n)
	if r.worklistSize > 0 {
		firstPart := len(r.worklist) - r.worklistHead
		if firstPart >= r.worklistSize {
			copy(nb, r.worklist[r.worklistHead:r.worklistHead+r.worklistSize])
		} else {
			copy(nb, r.worklist[r.worklistHead:])
			copy(nb[firstPart:], r.worklist[:r.worklistSize-firstPart])
		}
	}
	r.worklist = nb
	r.worklistHead = 0
}

func (r *RemoveUnusedCode) popContinuation() RemoveUnusedCode_Continuation {
	c := r.worklist[r.worklistHead]
	r.worklist[r.worklistHead] = RemoveUnusedCode_Continuation{}
	r.worklistHead++
	if r.worklistHead == len(r.worklist) {
		r.worklistHead = 0
	}
	r.worklistSize--
	return c
}

// TraverseAndRemoveUnusedReferences traverses a node recursively. Call this
// once per pass.
//
// Java: RemoveUnusedCode.traverseAndRemoveUnusedReferences(Node)
func (r *RemoveUnusedCode) TraverseAndRemoveUnusedReferences(root rhino.Node) {
	if r.pinnedPropertyNames == nil {
		r.pinnedPropertyNames = make(map[string]struct{}, 64)
		r.propRemovables = make(map[string]propRemovableList, 64)
	}
	// Create scope from parent of root node, which also has externs as a child, so we'll
	// have extern definitions in scope.
	scope := r.scopeCreator.CreateScopeOfScope(root.GetParent(), nil)
	if r.varBuckets == nil {
		r.varBuckets = make([]*RemoveUnusedCode_RealVarInfo, 1024)
	}
	if !scope.HasSlot(NodeUtil_JSC_PROPERTY_NAME_FN) {
		// TODO(b/70730762): Passes that add references to this should ensure it is declared.
		// NOTE: null input makes this an extern var.
		scope.Declare(
			NodeUtil_JSC_PROPERTY_NAME_FN /* no declaration node */, rhino.Node{} /* no input */, nil)
	}

	r.root = root
	r.guardedUsagesComputed = false

	r.enqueueContinuation(NewRemoveUnusedCode_Continuation(r, root, scope))
	for r.worklistSize > 0 {
		continuation := r.popContinuation() // worklist.remove()
		continuation.Apply(r)
	}

	r.RemoveUnreferencedVarsAndPolyfills()
	r.RemoveIndependentlyRemovableProperties()
	for _, fparamScope := range r.allFunctionParamScopes {
		r.RemoveUnreferencedFunctionArgs(fparamScope)
	}
}

// Java: RemoveUnusedCode.storePolyfill(PolyfillUsage)
func (r *RemoveUnusedCode) StorePolyfill(polyfillUsage *PolyfillUsageFinder_PolyfillUsage) {
	if r.guardedUsages == nil {
		r.guardedUsages = jutil.NewLinkedHashSet[rhino.Node]()
	}
	r.guardedUsages.Add(polyfillUsage.Node())
}

// Java: RemoveUnusedCode.removeIndependentlyRemovableProperties()
func (r *RemoveUnusedCode) RemoveIndependentlyRemovableProperties() {
	for _, propName := range r.propRemovableKeys {
		list, ok := r.propRemovables[propName]
		if !ok || list.head == nil {
			continue
		}
		if r.removalLog.IsLogging() {
			r.removalLog.LogSupplier(RemoveUnusedCode_RemovalLogRecord_ForProperty(propName).Get)
		}
		for cur := list.head; cur != nil; {
			next := cur.propNext
			cur.propNext = nil
			cur.flags &^= remFlagInPropList
			cur.Remove(r)
			cur = next
		}
	}
}

func hasBlockScopedDeclarations(n rhino.Node) bool {
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if hasBlockScopedDeclarationsChild(child) {
			return true
		}
	}
	return false
}

func hasBlockScopedDeclarationsChild(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CATCH, rhino.Token_IMPORT:
		return true
	case rhino.Token_EXPORT:
		first := n.GetFirstChild()
		return !first.IsNil() && hasBlockScopedDeclarationsChild(first)
	case rhino.Token_FUNCTION:
		return !NodeUtil_IsFunctionExpression(n) && n.GetFirstChild().GetString() != ""
	case rhino.Token_CLASS:
		return !NodeUtil_IsClassExpression(n) && n.GetFirstChild().GetString() != ""
	case rhino.Token_BLOCK:
		if scanVarsBlockCreatesScope(n) {
			return false
		}
		return hasBlockScopedDeclarations(n)
	case rhino.Token_ROOT, rhino.Token_DO, rhino.Token_WHILE, rhino.Token_WITH, rhino.Token_IF,
		rhino.Token_LABEL, rhino.Token_TRY, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
		return hasBlockScopedDeclarations(n)
	default:
		return false
	}
}

func (r *RemoveUnusedCode) traverseNameReference(n rhino.Node, scope *Scope) {
	// The only cases that should reach this point are parameter declarations and references
	// to names. The name node does not have children in these cases.
	if n.HasChildren() {
		jutil.CheckState(false)
	}
	parent := n.GetParent()
	// the parameter declaration is not a read of the name
	if pt := parent.GetToken(); pt != rhino.Token_PARAM_LIST {
		// var|let|const name;
		// function name() {}
		// class name() {}
		// are handled at a higher level.
		if pt == rhino.Token_VAR || pt == rhino.Token_LET || pt == rhino.Token_CONST ||
			((pt == rhino.Token_FUNCTION || pt == rhino.Token_CLASS) && parent.GetFirstChild() == n) {
			jutil.CheckState(false)
		}
		varInfo := r.TraverseNameNode(n, scope)
		if varInfo.IsRemovable() {
			var reason func() string
			if r.unremovableLog.IsLogging() {
				reason = func() string { return "reference found: " + n.GetLocation() }
			}
			varInfo.SetIsExplicitlyNotRemovable(reason)
		}
	}
}

// TraverseNode traverses everything in the current scope and marks
// variables that are referenced.
//
// During traversal, we identify subtrees that will only be referenced if
// their enclosing variables are referenced. Instead of traversing those
// subtrees, we create a continuation for them, and traverse them lazily.
//
// Java: RemoveUnusedCode.traverseNode(Node, Scope)
func (r *RemoveUnusedCode) TraverseNode(n rhino.Node, scope *Scope) {
	if n.GetFirstChild().IsNil() {
		if n.GetToken() == rhino.Token_NAME {
			r.traverseNameReference(n, scope)
		}
		return
	}
	type_ := n.GetToken()
	switch type_ {
	case rhino.Token_CATCH:
		r.TraverseCatch(n, scope)
	case rhino.Token_FUNCTION:
		{
			// If this function is a removable var, then create a continuation
			// for it instead of traversing immediately.
			if NodeUtil_IsFunctionDeclaration(n) {
				varInfo := r.TraverseNameNode(n.GetFirstChild(), scope)
				functionDeclaration :=
					NewRemoveUnusedCode_RemovableBuilder(r).
						AddContinuation(NewRemoveUnusedCode_Continuation(r, n, scope)).
						BuildFunctionDeclaration(n)
				varInfo.AddRemovable(functionDeclaration)
				if n.GetParent().IsExport() {
					varInfo.SetIsExplicitlyNotRemovable(func() string { return "exported class" })
				}
			} else {
				r.TraverseFunction(n, scope)
			}
		}
	case rhino.Token_ASSIGN:
		r.TraverseAssign(n, scope)
	case rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD:
		r.TraverseCompoundAssign(n, scope)
	case rhino.Token_INC, rhino.Token_DEC:
		r.TraverseIncrementOrDecrementOp(n, scope)
	case rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL:
		r.TraverseCall(n, scope)
	case rhino.Token_SWITCH_BODY:
		if hasBlockScopedDeclarations(n) {
			r.TraverseChildren(n, r.scopeCreator.CreateScopeOfScope(n, scope))
		} else {
			r.TraverseChildren(n, scope)
		}
	case rhino.Token_BLOCK:
		// This case if for if there are let and const variables in block scopes.
		// Otherwise other variables will be hoisted up into the global scope and already be
		// handled.
		if scanVarsBlockCreatesScope(n) &&
			(NodeUtil_IsFunctionBlock(n) || NodeUtil_IsClassStaticBlock(n) || hasBlockScopedDeclarations(n)) {
			r.TraverseChildren(n, r.scopeCreator.CreateScopeOfScope(n, scope))
		} else {
			r.TraverseChildren(n, scope)
		}
	case rhino.Token_MODULE_BODY:
		r.TraverseChildren(n, r.scopeCreator.CreateScopeOfScope(n, scope))
	case rhino.Token_CLASS:
		r.TraverseClass(n, scope)
	case rhino.Token_CLASS_MEMBERS:
		r.TraverseClassMembers(n, scope)
	case rhino.Token_ARRAY_PATTERN, rhino.Token_PARAM_LIST:
		r.TraverseIndirectAssignmentList(n, scope)
	case rhino.Token_OBJECT_PATTERN:
		r.TraverseObjectPattern(n, scope)
	case rhino.Token_OBJECTLIT:
		r.TraverseObjectLiteral(n, scope)
	case rhino.Token_FOR:
		r.TraverseVanillaFor(n, scope)
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		r.TraverseEnhancedFor(n, scope)
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		// for-loop cases are handled by custom traversal methods.
		if !NodeUtil_IsStatement(n) {
			jutil.CheckState(false)
		}
		r.TraverseDeclarationStatement(n, scope)
	case rhino.Token_INSTANCEOF:
		r.TraverseInstanceof(n, scope)
	case rhino.Token_NAME:
		r.traverseNameReference(n, scope)
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		r.TraverseNormalOrOptChainGetProp(n, scope)
	default:
		r.TraverseChildren(n, scope)
	}
}

// Java: RemoveUnusedCode.traverseInstanceof(Node, Scope)
func (r *RemoveUnusedCode) TraverseInstanceof(instanceofNode rhino.Node, scope *Scope) {
	if !instanceofNode.IsInstanceOf() {
		jutil.CheckArgumentMsg1(false, "%s", instanceofNode)
	}
	lhs := instanceofNode.GetFirstChild()
	rhs := lhs.GetNext()
	r.TraverseNode(lhs, scope)
	if rhs.IsName() {
		varInfo := r.TraverseNameNode(rhs, scope)
		builder := NewRemoveUnusedCode_RemovableBuilder(r)
		varInfo.AddRemovable(builder.BuildInstanceofName(instanceofNode))
	} else {
		r.TraverseNode(rhs, scope)
	}
}

func (r *RemoveUnusedCode) addPolyfillMask(name string) {
	r.hasPolyfills = true
	if len(name) > 0 {
		shift := uint(len(name))
		if shift >= 64 {
			shift = 63
		}
		mask := uint64(1) << shift
		r.polyfillLenMask[name[0]] |= mask
		r.polyfillEndMask[name[len(name)-1]] |= mask
	} else {
		r.polyfillLenMask[0] |= 1
		r.polyfillEndMask[0] |= 1
	}
}

func (r *RemoveUnusedCode) mightHavePolyfill(name string) bool {
	if !r.hasPolyfills {
		return false
	}
	if len(name) == 0 {
		return (r.polyfillLenMask[0] & 1) != 0
	}
	shift := uint(len(name))
	if shift >= 64 {
		shift = 63
	}
	mask := uint64(1) << shift
	return (r.polyfillLenMask[name[0]]&mask) != 0 && (r.polyfillEndMask[name[len(name)-1]]&mask) != 0
}

// TraverseNormalOrOptChainGetProp traverses `expr.prop` or `expr?.prop`.
//
// Note that this method is called only for RHS nodes. Property references
// that are being assigned to are handled by the logic traversing their
// parent (e.g. ASSIGN) node.
//
// The primary purpose of this method is to make sure the property reference
// is correctly recorded.
//
// Java: RemoveUnusedCode.traverseNormalOrOptChainGetProp(Node, Scope)
func (r *RemoveUnusedCode) TraverseNormalOrOptChainGetProp(getProp rhino.Node, scope *Scope) {
	if !NodeUtil_IsNormalOrOptChainGetProp(getProp) {
		jutil.CheckStateMsg1(false, "%s", getProp)
	}
	objectNode := getProp.GetFirstChild()
	propertyName := getProp.GetString()

	if r.mightHavePolyfill(propertyName) && r.polyfills != nil && !r.polyfills.IsEmpty() && r.polyfills.ContainsKey(propertyName) {
		for info := range r.polyfills.Get(propertyName).All() {
			if info.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable {
				info.ConsiderPossibleReference(getProp)
			}
		}
	}

	if NodeUtil_IsExpressionResultUsed(getProp) ||
		r.ConsiderForAccessorSideEffects(getProp, AccessorSummary_PropertyAccessKind_GETTER_ONLY) {
		// must record as reference to the property and continue traversal.
		r.MarkPropertyNameAsPinned(propertyName)
		r.TraverseNode(objectNode, scope)
	} else if objectNode.IsThis() {
		// This is probably the declaration of a class field in a constructor.
		// /** @private {number} */
		// this.propName;
		// We don't want to consider this a real usage that should prevent removal.
		builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsThisDotPropertyReference(true)
		r.ConsiderForIndependentRemoval(builder.BuildUnusedReadReference(getProp, getProp))
	} else if RemoveUnusedCode_IsDotPrototype(objectNode) {
		// (objExpression).prototype.propName;
		builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsPrototypeDotPropertyReference(true)
		objExpression := objectNode.GetFirstChild()
		if objExpression.IsName() {
			// name.prototype.propName;
			varInfo := r.TraverseNameNode(objExpression, scope)
			varInfo.AddRemovable(builder.BuildUnusedReadReference(getProp, getProp))
		} else {
			// (objExpression).prototype.propName;
			if r.MayHaveSideEffects(objExpression) {
				r.TraverseNode(objExpression, scope)
			} else {
				builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, objExpression, scope))
			}
			r.ConsiderForIndependentRemoval(builder.BuildUnusedReadReference(getProp, getProp))
		}
	} else {
		// TODO(bradfordcsmith): add removal of `varName.propName;`
		r.MarkPropertyNameAsPinned(propertyName)
		r.TraverseNode(objectNode, scope)
	}
}

// TraverseIncrementOrDecrementOp.
//
// TODO(b/137380742): Combine with `traverseCompoundAssign`.
//
// Java: RemoveUnusedCode.traverseIncrementOrDecrementOp(Node, Scope)
func (r *RemoveUnusedCode) TraverseIncrementOrDecrementOp(incOrDecOp rhino.Node, scope *Scope) {
	jutil.CheckArgumentMsg1(incOrDecOp.IsInc() || incOrDecOp.IsDec(), "%s", incOrDecOp)
	arg := incOrDecOp.GetOnlyChild()
	if NodeUtil_IsExpressionResultUsed(incOrDecOp) {
		// If expression result is used, then this expression is definitely not removable.
		r.TraverseNode(arg, scope)
	} else if arg.IsGetProp() {
		getPropObj := arg.GetFirstChild()

		if r.ConsiderForAccessorSideEffects(arg, AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER) {
			r.TraverseNode(getPropObj, scope) // Don't re-traverse the GETPROP as a read.
		} else if getPropObj.IsThis() {
			// this.propName++
			builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsThisDotPropertyReference(true)
			r.ConsiderForIndependentRemoval(builder.BuildIncOrDepOp(incOrDecOp, arg, rhino.Node{}))
		} else if RemoveUnusedCode_IsDotPrototype(getPropObj) {
			// someExpression.prototype.propName++
			exprObj := getPropObj.GetFirstChild()
			builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsPrototypeDotPropertyReference(true)
			if exprObj.IsName() {
				// varName.prototype.propName++
				varInfo := r.TraverseNameNode(exprObj, scope)
				varInfo.AddRemovable(builder.BuildIncOrDepOp(incOrDecOp, arg, rhino.Node{}))
			} else {
				// (someExpression).prototype.propName++
				var toPreserve rhino.Node
				if r.MayHaveSideEffects(exprObj) {
					toPreserve = exprObj
					r.TraverseNode(exprObj, scope)
				} else {
					builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, exprObj, scope))
				}
				r.ConsiderForIndependentRemoval(builder.BuildIncOrDepOp(incOrDecOp, arg, toPreserve))
			}
		} else {
			// someExpression.propName++ is not removable except in the cases covered above
			r.TraverseNode(arg, scope)
		}
	} else {
		// TODO(bradfordcsmith): varName++ should be removable if varName is otherwise unused
		r.TraverseNode(arg, scope)
	}
}

// TraverseCompoundAssign.
//
// TODO(b/137380742): Combine with `traverseIncrementOrDecrement`.
//
// Java: RemoveUnusedCode.traverseCompoundAssign(Node, Scope)
func (r *RemoveUnusedCode) TraverseCompoundAssign(compoundAssignNode rhino.Node, scope *Scope) {
	// We'll allow removal of compound assignment to a `this` property as long as the result of the
	// assignment is unused.
	// e.g. `this.prop += 3;`

	// NOTE: Some history here, as there were questions about "why is 'this' special".  The "remove
	// unused properties" is not a general property removal algorithm.  It only removes unreferenced
	// properties that are part of class definitions.  "SomeClass.prop += 3" and
	// "SomeClass.prototype.prop += 3" could so be candidates but they aren't considered here.

	targetNode := compoundAssignNode.GetFirstChild()
	valueNode := compoundAssignNode.GetLastChild()
	if targetNode.IsGetProp() {
		if r.ConsiderForAccessorSideEffects(targetNode, AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER) {
			r.TraverseNode(targetNode.GetFirstChild(), scope) // Don't re-traverse the GETPROP as a read.
			r.TraverseNode(valueNode, scope)
		} else if targetNode.GetFirstChild().IsThis() &&
			!NodeUtil_IsExpressionResultUsed(compoundAssignNode) {
			builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsThisDotPropertyReference(true)
			r.TraverseRemovableAssignValue(valueNode, builder, scope)
			r.ConsiderForIndependentRemoval(
				builder.BuildNamedPropertyAssign(compoundAssignNode, targetNode))
		} else {
			r.TraverseNode(targetNode, scope)
			r.TraverseNode(valueNode, scope)
		}
	} else {
		r.TraverseNode(targetNode, scope)
		r.TraverseNode(valueNode, scope)
	}
}

// Java: RemoveUnusedCode.traverseNameNode(Node, Scope)
func (r *RemoveUnusedCode) TraverseNameNode(n rhino.Node, scope *Scope) *RemoveUnusedCode_RealVarInfo {
	if r.hasPolyfills {
		name := n.GetString()
		if r.mightHavePolyfill(name) && r.polyfills != nil && !r.polyfills.IsEmpty() && r.polyfills.ContainsKey(name) {
			for info := range r.polyfills.Get(name).All() {
				if info.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable {
					info.ConsiderPossibleReference(n)
				}
			}
		}
	}
	v := scope.GetVarByID(n.GetStringID())
	if v == nil {
		jutil.CheckNotNullMsg1(v, "%s", n)
	}
	return r.TraverseVar(v)
}

func isSpecialCallCandidate(callee rhino.Node) bool {
	switch callee.GetToken() {
	case rhino.Token_GETPROP:
		switch callee.GetString() {
		case "inherits", "mixin", "defineProperties", "objectProperty",
			"addSingletonGetter", "polyfill", "patch", "polyfillTypedArrayMethod", "weakUsage":
			return true
		}
		return false
	case rhino.Token_NAME:
		name := callee.GetString()
		if len(name) < 6 {
			return false
		}
		switch name[len(name)-1] {
		case 'y':
			return name == NodeUtil_JSC_PROPERTY_NAME_FN || name == "goog$reflect$objectProperty"
		case 'l':
			return name == "$jscomp$polyfill"
		case 'h':
			return name == "$jscomp$patch"
		case 'd':
			return name == "$jscomp$polyfillTypedArrayMethod"
		case 'e':
			return name == "goog$weakUsage"
		case 'r':
			return name == "goog$addSingletonGetter"
		case 's':
			return strings.HasSuffix(name, "$inherits")
		case 'n':
			return strings.HasSuffix(name, "$mixin")
		}
		return false
	default:
		return false
	}
}

// Java: RemoveUnusedCode.traverseCall(Node, Scope)
func (r *RemoveUnusedCode) TraverseCall(callNode rhino.Node, scope *Scope) {
	callee := callNode.GetFirstChild()
	if !isSpecialCallCandidate(callee) {
		r.TraverseChildren(callNode, scope)
		return
	}

	if r.codingConvention.IsPropertyRenameFunction(callee) {
		propertyNameNode := callee.GetNext()
		if !propertyNameNode.IsNil() && propertyNameNode.IsStringLit() {
			r.MarkPropertyNameAsPinned(propertyNameNode.GetString())
		}
		r.TraverseChildren(callNode, scope)
	} else if NodeUtil_IsObjectDefinePropertiesDefinition(callNode) {
		// TODO(bradfordcsmith): Should also handle Object.create() and Object.defineProperty().
		r.TraverseObjectDefinePropertiesCall(callNode, scope)
	} else if r.removeUnusedPolyfills && RemoveUnusedCode_IsJscompPolyfill(callee) {
		firstArg := callee.GetNext()
		polyfillName := firstArg.GetString()
		if strings.HasSuffix(callee.GetString(), "polyfillTypedArrayMethod") {
			polyfillName = "TypedArray.prototype." + polyfillName
		}
		info := r.CreatePolyfillInfo(callNode, scope, polyfillName)
		if r.polyfills == nil {
			r.polyfills = jutil.HashMultimap_Create[string, RemoveUnusedCode_PolyfillInfo]()
		}
		key := info.GetRemoveUnusedCode_PolyfillInfoBase().key
		r.addPolyfillMask(key)
		r.polyfills.Put(key, info)
		// Only traverse the callee (to mark it as used).  The arguments may be traversed later.
		r.TraverseNode(callNode.GetFirstChild(), scope)
	} else if NodeUtil_IsGoogWeakUsageCall(callNode) &&
		callNode.HasTwoChildren() &&
		callNode.GetSecondChild().IsName() {
		// goog.weakUsage() should have exactly one argument, and it should be either a name or a
		// qualified name (this is checked in ProcessClosurePrimitives.java).
		//
		// If it is a qualified name, then we do not attempt to remove it at this time (see condition
		// above). We rely on CollapseProperties to turn qualified names into simple names where
		// possible.

		// Mark this call as removable if the var is not otherwise referenced.
		varInfo := r.TraverseNameNode(callNode.GetSecondChild(), scope)
		builder := NewRemoveUnusedCode_RemovableBuilder(r)
		varInfo.AddRemovable(builder.BuildWeakUsageCall(callNode))

		// We need to traverse the goog.weakUsage function itself (to mark it as used, in case our
		// usage of it is not removed).
		r.TraverseNode(callNode.GetFirstChild(), scope)
	} else {
		parent := callNode.GetParent()
		var classVarName *string
		classDefiningCall := false

		// A call that is a statement unto itself or the left side of a comma expression might be
		// a call to a known method for doing class setup
		// e.g. $jscomp.inherits(Class, BaseClass) or goog.addSingletonGetter(Class)
		// Such methods never have meaningful return values, so we won't look for them in other
		// contexts
		if parent.IsExprResult() || (parent.IsComma() && parent.GetFirstChild() == callNode) {
			subclassRelationship :=
				r.codingConvention.GetClassesDefinedByCall(callNode)
			if subclassRelationship != nil {
				// e.g. goog.inherits(DerivedClass, BaseClass);
				// NOTE: DerivedClass and BaseClass must be QNames. Otherwise getClassesDefinedByCall()
				// will return null.
				subclassName := subclassRelationship.subclassName
				classVarName = &subclassName
				classDefiningCall = true
			} else {
				// Look for calls to addSingletonGetter calls.
				classVarName = r.codingConvention.GetSingletonGetterClassName(callNode)
			}
		}

		var classVar *Var
		if classVarName != nil && NodeUtil_IsValidSimpleName(*classVarName) {
			classVar = scope.GetVar(*classVarName)
			if classVar == nil {
				jutil.CheckNotNullMsg1(classVar, "%s", *classVarName)
			}
		}

		if classVar == nil || !classVar.IsGlobal() {
			// The call we are traversing does not modify a class definition,
			// or the class is not specified with a simple variable name,
			// or the variable name is not global.
			// TODO(bradfordcsmith): It would be more correct to check whether the class name
			// references a known constructor and expand to allow QNames.
			r.TraverseChildren(callNode, scope)
		} else {
			builder := NewRemoveUnusedCode_RemovableBuilder(r)
			for child := callNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
				builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, child, scope))
			}
			r.TraverseVar(classVar).
				AddRemovable(builder.BuildClassSetupCallWithClassDefiningCall(callNode, classDefiningCall))
		}
	}
}

// RemoveUnusedCode_IsJscompPolyfill checks whether this is a recognizable
// call to $jscomp.polyfill.
//
// Java: RemoveUnusedCode.isJscompPolyfill(Node)
func RemoveUnusedCode_IsJscompPolyfill(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		{
			// Need to work correctly after CollapseProperties.
			name := n.GetString()
			return (name == "$jscomp$polyfill" ||
				name == "$jscomp$patch" ||
				name == "$jscomp$polyfillTypedArrayMethod") &&
				n.GetNext().IsStringLit()
		}
	case rhino.Token_GETPROP:
		{
			// Need to work correctly without CollapseProperties.
			propertyName := n.GetString()
			return (propertyName == "polyfill" ||
				propertyName == "patch" ||
				propertyName == "polyfillTypedArrayMethod") &&
				n.GetFirstChild().IsName() &&
				n.GetFirstChild().GetString() == "$jscomp" &&
				n.GetNext().IsStringLit()
		}
	default:
		return false
	}
}

// TraverseObjectDefinePropertiesCall traverses
// `Object.defineProperties(someObject, propertyDefinitions);`.
//
// Java: RemoveUnusedCode.traverseObjectDefinePropertiesCall(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectDefinePropertiesCall(callNode rhino.Node, scope *Scope) {
	// First child is Object.defineProperties or some equivalent of it.
	callee := callNode.GetFirstChild()
	targetObject := callNode.GetSecondChild()
	propertyDefinitions := targetObject.GetNext()

	if (targetObject.IsName() || r.IsNameDotPrototype(targetObject)) &&
		!NodeUtil_IsExpressionResultUsed(callNode) {
		// NOTE: Object.defineProperties() returns its first argument, so if its return value is used
		// that counts as a use of the targetObject.
		var nameNode rhino.Node
		if targetObject.IsName() {
			nameNode = targetObject
		} else {
			nameNode = targetObject.GetFirstChild()
		}
		varInfo := r.TraverseNameNode(nameNode, scope)
		builder := NewRemoveUnusedCode_RemovableBuilder(r)
		// TODO(bradfordcsmith): Is it really necessary to traverse the callee
		// (aka. Object.defineProperties)?
		builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, callee, scope))
		if r.MayHaveSideEffects(propertyDefinitions) {
			r.TraverseNode(propertyDefinitions, scope)
		} else {
			builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, propertyDefinitions, scope))
		}
		varInfo.AddRemovable(builder.BuildClassSetupCall(callNode))
	} else {
		// TODO(bradfordcsmith): Is it really necessary to traverse the callee
		// (aka. Object.defineProperties)?
		r.TraverseNode(callee, scope)
		r.TraverseNode(targetObject, scope)
		r.TraverseNode(propertyDefinitions, scope)
	}
}

// TraverseObjectDefinePropertiesLiteral traverses the object literal passed
// as the second argument to `Object.defineProperties()`.
//
// Java: RemoveUnusedCode.traverseObjectDefinePropertiesLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectDefinePropertiesLiteral(propertyDefinitions rhino.Node, scope *Scope) {
	for property := propertyDefinitions.GetFirstChild(); !property.IsNil(); property = property.GetNext() {
		if property.IsQuotedStringKey() {
			// Quoted property name counts as a reference to the property and protects it from removal.
			r.MarkPropertyNameAsPinned(property.GetString())
			r.TraverseNode(property.GetOnlyChild(), scope)
		} else if property.IsStringKey() {
			definition := property.GetOnlyChild()
			if r.MayHaveSideEffects(definition) {
				r.TraverseNode(definition, scope)
			} else {
				r.ConsiderForIndependentRemoval(
					NewRemoveUnusedCode_RemovableBuilder(r).
						AddContinuation(NewRemoveUnusedCode_Continuation(r, definition, scope)).
						BuildObjectDefinePropertiesDefinition(property))
			}
		} else {
			// TODO(bradfordcsmith): Maybe report error for anything other than a computed property,
			// since getters, setters, and methods don't make much sense in this context.
			r.TraverseNode(property, scope)
		}
	}
}

// Java: RemoveUnusedCode.getVarForNameNode(Node, Scope)
func (r *RemoveUnusedCode) GetVarForNameNode(nameNode rhino.Node, scope *Scope) *Var {
	v := scope.GetVarByNode(nameNode)
	if v == nil {
		jutil.CheckNotNullMsg1(v, "%s", nameNode)
	}
	return v
}

// Java: RemoveUnusedCode.traverseObjectLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectLiteral(objectLiteral rhino.Node, scope *Scope) {
	if !objectLiteral.IsObjectLit() {
		jutil.CheckArgumentMsg1(false, "%s", objectLiteral)
	}
	// Is this an object literal that is assigned directly to a 'prototype' property?
	if r.IsAssignmentToPrototype(objectLiteral.GetParent()) {
		r.TraversePrototypeLiteral(objectLiteral, scope)
	} else if r.IsObjectDefinePropertiesSecondArgument(objectLiteral) {
		// TODO(bradfordcsmith): Consider restricting special handling of the properties literal to
		// cases where the target object is a known class, prototype, or this.
		r.TraverseObjectDefinePropertiesLiteral(objectLiteral, scope)
	} else {
		r.TraverseNonPrototypeObjectLiteral(objectLiteral, scope)
	}
}

// Java: RemoveUnusedCode.isObjectDefinePropertiesSecondArgument(Node)
func (r *RemoveUnusedCode) IsObjectDefinePropertiesSecondArgument(n rhino.Node) bool {
	parent := n.GetParent()
	return parent.IsCall() && parent.GetLastChild() == n && NodeUtil_IsObjectDefinePropertiesDefinition(parent)
}

// Java: RemoveUnusedCode.traverseNonPrototypeObjectLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraverseNonPrototypeObjectLiteral(objectLiteral rhino.Node, scope *Scope) {
	for propertyNode := objectLiteral.GetFirstChild(); !propertyNode.IsNil(); propertyNode = propertyNode.GetNext() {
		if propertyNode.IsStringKey() {
			// A property name in an object literal counts as a reference,
			// because of some reflection patterns.
			// Note that we are intentionally treating both quoted and unquoted keys as
			// references.
			r.MarkPropertyNameAsPinned(propertyNode.GetString())
			r.TraverseNode(propertyNode.GetFirstChild(), scope)
		} else {
			r.TraverseNode(propertyNode, scope)
		}
	}
}

// Java: RemoveUnusedCode.traversePrototypeLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraversePrototypeLiteral(objectLiteral rhino.Node, scope *Scope) {
	for propertyNode := objectLiteral.GetFirstChild(); !propertyNode.IsNil(); propertyNode = propertyNode.GetNext() {
		if propertyNode.IsComputedProp() || propertyNode.IsQuotedStringKey() {
			r.TraverseChildren(propertyNode, scope)
		} else {
			valueNode := propertyNode.GetOnlyChild()
			if r.MayHaveSideEffects(valueNode) {
				// TODO(bradfordcsmith): Ideally we should preserve the side-effect without keeping the
				// property itself alive.
				r.TraverseNode(valueNode, scope)
			} else {
				// If we've come this far, we already know we're keeping the prototype literal itself,
				// but we may be able to remove unreferenced properties in it.
				r.ConsiderForIndependentRemoval(
					NewRemoveUnusedCode_RemovableBuilder(r).
						AddContinuation(NewRemoveUnusedCode_Continuation(r, valueNode, scope)).
						BuildClassOrPrototypeNamedProperty(propertyNode))
			}
		}
	}
}

// Java: RemoveUnusedCode.isAssignmentToPrototype(Node)
func (r *RemoveUnusedCode) IsAssignmentToPrototype(n rhino.Node) bool {
	return n.IsAssign() && RemoveUnusedCode_IsDotPrototype(n.GetFirstChild())
}

// RemoveUnusedCode_IsDotPrototype: true for `someExpression.prototype`.
//
// Java: RemoveUnusedCode.isDotPrototype(Node)
func RemoveUnusedCode_IsDotPrototype(n rhino.Node) bool {
	return NodeUtil_IsNormalOrOptChainGetProp(n) && n.GetString() == "prototype"
}

// Java: RemoveUnusedCode.traverseCatch(Node, Scope)
func (r *RemoveUnusedCode) TraverseCatch(catchNode rhino.Node, scope *Scope) {
	exceptionNameNode := catchNode.GetFirstChild()
	block := exceptionNameNode.GetNext()
	if exceptionNameNode.IsName() {
		// exceptionNameNode can be an empty node if not using a binding in 2019.
		exceptionVarInfo := r.TraverseNameNode(exceptionNameNode, scope)
		exceptionVarInfo.SetIsExplicitlyNotRemovable(func() string { return "catch variable" })
	} else {
		r.TraverseNode(exceptionNameNode, scope)
	}
	r.TraverseNode(block, scope)
}

// Java: RemoveUnusedCode.traverseEnhancedFor(Node, Scope)
func (r *RemoveUnusedCode) TraverseEnhancedFor(enhancedFor rhino.Node, scope *Scope) {
	forScope := scope
	if hasBlockScopedDeclarations(enhancedFor) {
		forScope = r.scopeCreator.CreateScopeOfScope(enhancedFor, scope)
	}
	// for (iterationTarget in|of collection) body;
	iterationTarget := enhancedFor.GetFirstChild()
	collection := iterationTarget.GetNext()
	body := collection.GetNext()
	if iterationTarget.IsName() {
		// using previously-declared loop variable. e.g.
		// `for (varName of collection) {}`
		varInfo := r.TraverseNameNode(iterationTarget, forScope)
		varInfo.SetIsExplicitlyNotRemovable(func() string { return "for-of or for-in loop variable" })
	} else if NodeUtil_IsNameDeclaration(iterationTarget) {
		// loop has const/var/let declaration
		declNode := iterationTarget.GetOnlyChild()
		if declNode.IsDestructuringLhs() {
			// e.g.
			// `for (const [a, b] of pairList) {}`
			// destructuring is handled at a lower level
			// Note that destructuring assignments are always considered to set an unknown value
			// equivalent to what we set for the var name case above and below.
			// It isn't necessary to set the variable names as not removable, though, because the
			// thing that isn't removable is the destructuring pattern itself, which we never remove.
			// TODO(bradfordcsmith): The need to explain all the above shows this should be reworked.
			r.TraverseNode(declNode, forScope)
		} else {
			// e.g.
			// `for (const varName of collection) {}`
			jutil.CheckState(declNode.IsName())
			jutil.CheckState(!declNode.HasChildren())
			// We can never remove the loop variable of a for-in or for-of loop, because it's
			// essential to loop syntax.
			varInfo := r.TraverseNameNode(declNode, forScope)
			varInfo.SetIsExplicitlyNotRemovable(func() string { return "for-of or for-in loop variable" })
		}
	} else {
		// using some general LHS value e.g.
		// `for ([a, b] of collection) {}` destructuring with existing vars
		// `for (a.x of collection) {}` using a property as the loop var
		// TODO(bradfordcsmith): This should be considered a write if it's a property reference.
		r.TraverseNode(iterationTarget, forScope)
	}
	r.TraverseNode(collection, forScope)
	r.TraverseNode(body, forScope)
}

// Java: RemoveUnusedCode.traverseVanillaFor(Node, Scope)
func (r *RemoveUnusedCode) TraverseVanillaFor(forNode rhino.Node, scope *Scope) {
	forScope := scope
	if hasBlockScopedDeclarations(forNode) {
		forScope = r.scopeCreator.CreateScopeOfScope(forNode, scope)
	}
	initialization := forNode.GetFirstChild()
	condition := initialization.GetNext()
	update := condition.GetNext()
	block := update.GetNext()
	if NodeUtil_IsNameDeclaration(initialization) {
		r.TraverseVanillaForNameDeclarations(initialization, forScope)
	} else {
		r.TraverseNode(initialization, forScope)
	}
	r.TraverseNode(condition, forScope)
	r.TraverseNode(update, forScope)
	r.TraverseNode(block, forScope)
}

// Java: RemoveUnusedCode.traverseVanillaForNameDeclarations(Node, Scope)
func (r *RemoveUnusedCode) TraverseVanillaForNameDeclarations(nameDeclaration rhino.Node, scope *Scope) {
	for child := nameDeclaration.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if !child.IsName() {
			// TODO(bradfordcsmith): Customize handling of destructuring
			r.TraverseNode(child, scope)
		} else {
			nameNode := child
			valueNode := child.GetFirstChild()
			varInfo := r.TraverseNameNode(nameNode, scope)
			if valueNode.IsNil() {
				varInfo.AddRemovable(NewRemoveUnusedCode_RemovableBuilder(r).BuildVanillaForNameDeclaration(nameNode))
			} else if r.MayHaveSideEffects(valueNode) {
				// TODO(bradfordcsmith): Actually allow for removing the variable while keeping the
				// valueNode for its side-effects.
				varInfo.SetIsExplicitlyNotRemovable(
					func() string { return "for-loop variable initialization has side-effects" })
				r.TraverseNode(valueNode, scope)
			} else {
				vanillaForNameDeclaration :=
					NewRemoveUnusedCode_RemovableBuilder(r).
						AddContinuation(NewRemoveUnusedCode_Continuation(r, valueNode, scope)).
						BuildVanillaForNameDeclaration(nameNode)
				varInfo.AddRemovable(vanillaForNameDeclaration)
			}
		}
	}
}

// Java: RemoveUnusedCode.traverseDeclarationStatement(Node, Scope)
func (r *RemoveUnusedCode) TraverseDeclarationStatement(declarationStatement rhino.Node, scope *Scope) {
	// Normalization should ensure that declaration statements always have just one child.
	nameNode := declarationStatement.GetOnlyChild()
	if !nameNode.IsName() {
		// Destructuring declarations are handled elsewhere.
		r.TraverseNode(nameNode, scope)
	} else {
		valueNode := nameNode.GetFirstChild()
		varInfo := r.TraverseNameNode(nameNode, scope)
		builder := NewRemoveUnusedCode_RemovableBuilder(r)
		if valueNode.IsNil() {
			varInfo.AddRemovable(builder.BuildNameDeclarationStatement(declarationStatement))
		} else {
			if r.MayHaveSideEffects(valueNode) {
				r.TraverseNode(valueNode, scope)
			} else {
				builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, valueNode, scope))
			}
			removable :=
				builder.BuildNameDeclarationStatement(declarationStatement)
			varInfo.AddRemovable(removable)
		}
	}
}

// Java: RemoveUnusedCode.traverseAssign(Node, Scope)
func (r *RemoveUnusedCode) TraverseAssign(assignNode rhino.Node, scope *Scope) {
	if !NodeUtil_IsAssignmentOp(assignNode) {
		jutil.CheckState(false)
	}

	lhs := assignNode.GetFirstChild()
	valueNode := assignNode.GetLastChild()
	if lhs.IsName() {
		// varName = something
		varInfo := r.TraverseNameNode(lhs, scope)
		builder := NewRemoveUnusedCode_RemovableBuilder(r)
		r.TraverseRemovableAssignValue(valueNode, builder, scope)
		varInfo.AddRemovable(builder.BuildVariableAssign(assignNode, varInfo))
	} else if lhs.IsGetElem() {
		getElemObj := lhs.GetFirstChild()
		getElemKey := lhs.GetLastChild()
		var varNameNode rhino.Node
		if getElemObj.IsName() {
			varNameNode = getElemObj
		} else if r.IsNameDotPrototype(getElemObj) {
			varNameNode = getElemObj.GetFirstChild()
		} else {
			varNameNode = rhino.Node{}
		}

		if !varNameNode.IsNil() {
			// varName[someExpression] = someValue
			// OR
			// varName.prototype[someExpression] = someValue
			varInfo := r.TraverseNameNode(varNameNode, scope)
			builder := NewRemoveUnusedCode_RemovableBuilder(r)
			if r.MayHaveSideEffects(getElemKey) {
				r.TraverseNode(getElemKey, scope)
			} else {
				builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, getElemKey, scope))
			}
			r.TraverseRemovableAssignValue(valueNode, builder, scope)
			varInfo.AddRemovable(builder.BuildComputedPropertyAssign(assignNode, getElemKey, varInfo))
		} else {
			r.TraverseNode(getElemObj, scope)
			r.TraverseNode(getElemKey, scope)
			r.TraverseNode(valueNode, scope)
		}
	} else if lhs.IsGetProp() {
		getPropLhs := lhs.GetFirstChild()
		// Assignments `Foo.prototype.bar = function() {`
		isDotPrototypeLhs := RemoveUnusedCode_IsDotPrototype(getPropLhs)
		isPrototypeMethodDef := isDotPrototypeLhs && valueNode.IsFunction()

		if !isPrototypeMethodDef &&
			r.ConsiderForAccessorSideEffects(lhs, AccessorSummary_PropertyAccessKind_SETTER_ONLY) {
			// And the possible side-effects mean we can't do any removal. We don't use the
			// `AstAnalyzer` because we only want to consider side-effect from the assignment, not the
			// entire l-value subtree.
			// Assume prototype method assignments never trigger setters, matching ES class semantics
			r.TraverseNode(getPropLhs, scope) // Don't re-traverse the GETPROP as a read.
			r.TraverseNode(valueNode, scope)
		} else if getPropLhs.IsName() {
			// varName.propertyName = someValue
			varInfo := r.TraverseNameNode(getPropLhs, scope)
			builder := NewRemoveUnusedCode_RemovableBuilder(r)
			r.TraverseRemovableAssignValue(valueNode, builder, scope)
			varInfo.AddRemovable(builder.BuildNamedPropertyAssignWithVarInfo(assignNode, lhs, varInfo))
		} else if isDotPrototypeLhs {
			// objExpression.prototype.propertyName = someValue
			objExpression := getPropLhs.GetFirstChild()
			builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsPrototypeDotPropertyReference(true)
			r.TraverseRemovableAssignValue(valueNode, builder, scope)
			if objExpression.IsName() {
				// varName.prototype.propertyName = someValue
				varInfo := r.TraverseNameNode(getPropLhs.GetFirstChild(), scope)
				varInfo.AddRemovable(builder.BuildNamedPropertyAssignWithVarInfo(assignNode, lhs, varInfo))
			} else {
				// (someExpression).prototype.propertyName = someValue
				if r.MayHaveSideEffects(objExpression) {
					r.TraverseNode(objExpression, scope)
				} else {
					builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, objExpression, scope))
				}
				r.ConsiderForIndependentRemoval(
					builder.BuildAnonymousPrototypeNamedPropertyAssign(assignNode, lhs.GetString()))
			}
		} else if getPropLhs.IsThis() {
			// this.propertyName = someValue
			builder := NewRemoveUnusedCode_RemovableBuilder(r).SetIsThisDotPropertyReference(true)
			r.TraverseRemovableAssignValue(valueNode, builder, scope)
			r.ConsiderForIndependentRemoval(builder.BuildNamedPropertyAssign(assignNode, lhs))
		} else {
			r.TraverseNode(lhs, scope)
			r.TraverseNode(valueNode, scope)
		}
	} else {
		// no other cases are removable
		r.TraverseNode(lhs, scope)
		r.TraverseNode(valueNode, scope)
	}
}

// Java: RemoveUnusedCode.traverseRemovableAssignValue(Node, RemovableBuilder, Scope)
func (r *RemoveUnusedCode) TraverseRemovableAssignValue(valueNode rhino.Node, builder *RemoveUnusedCode_RemovableBuilder, scope *Scope) {
	if r.MayHaveSideEffects(valueNode) || NodeUtil_IsExpressionResultUsed(valueNode.GetParent()) {
		r.TraverseNode(valueNode, scope)
	} else {
		builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, valueNode, scope))
	}
}

// Java: RemoveUnusedCode.isNameDotPrototype(Node)
func (r *RemoveUnusedCode) IsNameDotPrototype(n rhino.Node) bool {
	return n.IsGetProp() && n.GetFirstChild().IsName() && n.GetString() == "prototype"
}

// Java: RemoveUnusedCode.traverseObjectPattern(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectPattern(pattern rhino.Node, scope *Scope) {
	if !pattern.IsObjectPattern() {
		jutil.CheckStateMsg1(false, "%s", pattern)
	}

	for elem := pattern.GetFirstChild(); !elem.IsNil(); elem = elem.GetNext() {
		switch elem.GetToken() {
		case rhino.Token_COMPUTED_PROP:
			r.TraverseIndirectAssignment(elem, elem.GetSecondChild(), scope)
		case rhino.Token_STRING_KEY:
			if !elem.IsQuotedStringKey() {
				r.MarkPropertyNameAsPinned(elem.GetString())
			}
			r.TraverseIndirectAssignment(elem, elem.GetOnlyChild(), scope)
		case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
			// Recall that the rest target can be any l-value expression
			r.TraverseIndirectAssignment(elem, elem.GetOnlyChild(), scope)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of " + pattern.GetToken().String() + ": " + elem.ToStringTree()))
		}
	}
}

// Java: RemoveUnusedCode.traverseIndirectAssignmentList(Node, Scope)
func (r *RemoveUnusedCode) TraverseIndirectAssignmentList(list rhino.Node, scope *Scope) {
	if !list.IsArrayPattern() && !list.IsParamList() {
		jutil.CheckStateMsg1(false, "%s", list)
	}

	for elem := list.GetFirstChild(); !elem.IsNil(); elem = elem.GetNext() {
		switch elem.GetToken() {
		case rhino.Token_EMPTY:
		case rhino.Token_ARRAY_PATTERN, rhino.Token_DEFAULT_VALUE, rhino.Token_GETELEM, rhino.Token_GETPROP, rhino.Token_NAME, rhino.Token_OBJECT_PATTERN:
			r.TraverseIndirectAssignment(elem, elem, scope)
		case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
			r.TraverseIndirectAssignment(elem, elem.GetOnlyChild(), scope)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of " + list.GetToken().String() + ": " + elem.ToStringTree()))
		}
	}
}

// TraverseIndirectAssignment traverses an AST structure representing an
// assignment operation for which the target and value are far apart.
//
// Examples include destructurings and function parameters.
//
// root is the root of the assignment subtree; target the l-value expression
// being assigned to.
//
// Java: RemoveUnusedCode.traverseIndirectAssignment(Node, Node, Scope)
func (r *RemoveUnusedCode) TraverseIndirectAssignment(root rhino.Node, target rhino.Node, scope *Scope) {
	rootParent := root.GetParent()
	if !rootParent.IsDestructuringPattern() && !rootParent.IsParamList() {
		jutil.CheckArgumentMsg1(false, "%s", rootParent)
	}

	// Flatten out the case where the target is a default value. We always have to consider it.
	if target.IsDefaultValue() {
		target = target.GetFirstChild()
	}

	if target.IsGetProp() {
		r.ConsiderForAccessorSideEffects(target, AccessorSummary_PropertyAccessKind_SETTER_ONLY)
	}

	builder :=
		NewRemoveUnusedCode_RemovableBuilder(r).AddContinuation(NewRemoveUnusedCode_Continuation(r, root, scope))

	if r.MayHaveSideEffects(root) {
		// If anywhere in the assignment subtree has side-effects, it means that even if the target is
		// removable the subtree is not.
		r.TraverseNode(root, scope)
		// TODO(bradfordcsmith): Preserve side effects without preventing removal of variables and
		// properties. We could probably do this by subbing in an empty object pattern.
	} else if target.IsName() {
		varInfo := r.TraverseNameNode(target, scope)
		varInfo.AddRemovable(builder.BuildIndirectAssign(root, target))
	} else if r.IsNameDotPrototype(target) || RemoveUnusedCode_IsThisDotProperty(target) {
		r.ConsiderForIndependentRemoval(builder.BuildIndirectAssign(root, target))
	} else {
		// TODO(bradfordcsmith): Handle property assignments also
		// e.g. `({a: foo.bar, b: foo.baz}) = {a: 1, b: 2}`
		r.TraverseNode(root, scope)
	}
}

// Java: RemoveUnusedCode.traverseChildren(Node, Scope)
func (r *RemoveUnusedCode) TraverseChildren(n rhino.Node, scope *Scope) {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.GetToken() == rhino.Token_NAME {
				r.traverseNameReference(c, scope)
			}
			continue
		}
		r.TraverseNode(c, scope)
	}
}

// TraverseClass handles a class that is not the RHS child of an assignment
// or a variable declaration initializer.
//
// Java: RemoveUnusedCode.traverseClass(Node, Scope)
func (r *RemoveUnusedCode) TraverseClass(classNode rhino.Node, scope *Scope) {
	jutil.CheckArgument(classNode.IsClass())
	if NodeUtil_IsClassDeclaration(classNode) {
		r.TraverseClassDeclaration(classNode, scope)
	} else {
		r.TraverseClassExpression(classNode, scope)
	}
}

// Java: RemoveUnusedCode.traverseClassDeclaration(Node, Scope)
func (r *RemoveUnusedCode) TraverseClassDeclaration(classNode rhino.Node, scope *Scope) {
	jutil.CheckArgument(classNode.IsClass())
	classNameNode := classNode.GetFirstChild()
	baseClassExpression := classNameNode.GetNext()
	classBodyNode := baseClassExpression.GetNext()
	// Class declarations never declare the class name in the inner class scope (only class
	// expressions do), so creating an inner scope here would always produce an empty scope.
	classScope := scope

	varInfo := r.TraverseNameNode(classNameNode, scope)
	if classNode.GetParent().IsExport() {
		// Cannot remove an exported class.
		varInfo.SetIsExplicitlyNotRemovable(func() string { return "exported class" })
		r.TraverseNode(baseClassExpression, scope)
		// Use traverseChildren() here, because we should not consider any properties on the exported
		// class to be removable.
		r.TraverseChildren(classBodyNode, classScope)
	} else if r.MayHaveSideEffects(baseClassExpression) {
		// TODO(bradfordcsmith): implement removal without losing side-effects for this case
		varInfo.SetIsExplicitlyNotRemovable(func() string { return "base class expression has side-effects" })
		r.TraverseNode(baseClassExpression, scope)
		r.TraverseClassMembers(classBodyNode, classScope)
	} else if r.MayHaveSideEffects(classBodyNode) {
		varInfo.SetIsExplicitlyNotRemovable(func() string { return "class body has side-effects" })
		r.TraverseNode(baseClassExpression, scope)
		r.TraverseClassMembers(classBodyNode, classScope)
	} else {
		builder :=
			NewRemoveUnusedCode_RemovableBuilder(r).
				AddContinuation(NewRemoveUnusedCode_Continuation(r, baseClassExpression, classScope)).
				AddContinuation(NewRemoveUnusedCode_Continuation(r, classBodyNode, classScope))
		varInfo.AddRemovable(builder.BuildClassDeclaration(classNode))
	}
}

// Java: RemoveUnusedCode.traverseClassExpression(Node, Scope)
func (r *RemoveUnusedCode) TraverseClassExpression(classNode rhino.Node, scope *Scope) {
	jutil.CheckArgument(classNode.IsClass())
	classNameNode := classNode.GetFirstChild()
	baseClassExpression := classNameNode.GetNext()
	classBodyNode := baseClassExpression.GetNext()
	classScope := scope

	if classNameNode.IsName() {
		classScope = r.scopeCreator.CreateScopeOfScope(classNode, scope)
		// We may be able to remove the name node if nothing ends up referring to it.
		varInfo := r.TraverseNameNode(classNameNode, classScope)
		// The class is non-local, because it is accessible by unknown code outside
		// of the scope where InnerName is defined.
		// e.g. `use(class InnerName {})`
		varInfo.SetHasNonLocalOrNonLiteralValue()
		varInfo.AddRemovable(NewRemoveUnusedCode_RemovableBuilder(r).BuildNamedClassExpression(classNode))
	}
	// If we're traversing the class expression, we've already decided we cannot remove it.
	r.TraverseNode(baseClassExpression, scope)
	r.TraverseClassMembers(classBodyNode, classScope)
}

// Java: RemoveUnusedCode.traverseClassMembers(Node, Scope)
func (r *RemoveUnusedCode) TraverseClassMembers(node rhino.Node, scope *Scope) {
	jutil.CheckArgumentMsg1(node.IsClassMembers(), "%s", node)
	if !r.removeUnusedPrototypeProperties {
		r.TraverseChildren(node, scope)
		return
	}

	for member := node.GetFirstChild(); !member.IsNil(); member = member.GetNext() {
		switch member.GetToken() {
		case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
			// If we get as far as traversing the members of a class, we've already decided that
			// we cannot remove the class itself, so just consider individual members for removal.
			r.ConsiderForIndependentRemoval(
				NewRemoveUnusedCode_RemovableBuilder(r).
					AddContinuation(NewRemoveUnusedCode_Continuation(r, member, scope)).
					BuildClassOrPrototypeNamedProperty(member))
		case rhino.Token_MEMBER_FIELD_DEF:
			// TODO(bradfordcsmith): currently if the RHS of a field has side effects, we do not
			// remove any part of the field. The proper behavior of class C { x = alert(); }
			// would be to remove x, leaving class C { constructor() { alert(); } }
			// but currently we aren't removing anything.
			if !member.HasChildren() || !r.MayHaveSideEffects(member.GetFirstChild()) {
				r.ConsiderForIndependentRemoval(
					NewRemoveUnusedCode_RemovableBuilder(r).
						AddContinuation(NewRemoveUnusedCode_Continuation(r, member, scope)).
						BuildClassOrPrototypeNamedProperty(member))
			}
			// TODO: b/354704593 - remove the entire class when it is not referenced
			if member.HasChildren() {
				r.TraverseChildren(member, scope)
			}
		case rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF:
			r.TraverseChildren(member, scope)
		case rhino.Token_BLOCK:
			r.TraverseChildren(member, r.scopeCreator.CreateScopeOfScope(member, scope))
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of CLASS_MEMBERS: " + member.ToStringTree()))
		}
	}
}

// TraverseFunction traverses a function
//
// ES6 scopes of a function include the parameter scope and the body scope of
// the function.
//
// Note that CATCH blocks also create a new scope, but only for the catch
// variable. Declarations within the block actually belong to the enclosing
// scope. Because we don't remove catch variables, there's no need to treat
// CATCH blocks differently like we do functions.
//
// Java: RemoveUnusedCode.traverseFunction(Node, Scope)
func (r *RemoveUnusedCode) TraverseFunction(function rhino.Node, parentScope *Scope) {
	nameNode := function.GetFirstChild()
	var paramlist, body rhino.Node
	if !nameNode.IsNil() {
		paramlist = nameNode.GetNext()
		if !paramlist.IsNil() {
			body = paramlist.GetNext()
		}
	}
	if body.IsNil() || !body.GetNext().IsNil() || !function.IsFunction() {
		jutil.CheckStateMsg1(function.HasXChildren(3), "%s", function)
		jutil.CheckStateMsg1(function.IsFunction(), "%s", function)
	}
	if !body.IsBlock() {
		jutil.CheckStateMsg1(false, "%s", body)
	}

	// Checking the parameters
	fparamScope := r.scopeCreator.CreateScopeOfScope(function, parentScope)

	// Checking the function body
	fbodyScope := r.scopeCreator.CreateScopeOfScope(body, fparamScope)

	if nameNode.GetString() != "" {
		// var x = function funcName() {};
		// make sure funcName gets into the varInfoMap so it will be considered for removal.
		varInfo := r.TraverseNameNode(nameNode, fparamScope)
		if NodeUtil_IsExpressionResultUsed(function) {
			// var f = function g() {};
			// The f is an alias for g, so g escapes from the scope where it is defined.
			varInfo.SetHasNonLocalOrNonLiteralValue()
		}
	}

	if !paramlist.GetFirstChild().IsNil() {
		r.TraverseIndirectAssignmentList(paramlist, fparamScope)
		if r.removeGlobals {
			r.allFunctionParamScopes = append(r.allFunctionParamScopes, fparamScope)
		}
	}
	r.TraverseChildren(body, fbodyScope)
}

// Java: RemoveUnusedCode.canRemoveParameters(Node)
func (r *RemoveUnusedCode) CanRemoveParameters(parameterList rhino.Node) bool {
	if !parameterList.IsParamList() {
		jutil.CheckState(false)
	}
	function := parameterList.GetParent()
	return r.removeGlobals && !NodeUtil_IsGetOrSetKey(function.GetParent())
}

// RemoveUnreferencedFunctionArgs removes unreferenced arguments from a
// function declaration and when possible the function's callSites.
//
// fparamScope is the function parameter scope.
//
// Java: RemoveUnusedCode.removeUnreferencedFunctionArgs(Scope)
func (r *RemoveUnusedCode) RemoveUnreferencedFunctionArgs(fparamScope *Scope) {
	// Notice that removing unreferenced function args breaks
	// Function.prototype.length. In advanced mode, we don't really care
	// about this: we consider "length" the equivalent of reflecting on
	// the function's lexical source.
	//
	// Rather than create a new option for this, we assume that if the user
	// is removing globals, then it's OK to remove unused function args.
	//
	// See http://blickly.github.io/closure-compiler-issues/#253
	if !r.removeGlobals {
		return
	}

	function := fparamScope.GetRootNode()
	if !function.IsFunction() {
		jutil.CheckState(false)
	}
	if NodeUtil_IsGetOrSetKey(function.GetParent()) {
		// The parameters object literal setters can not be removed.
		return
	}

	argList := NodeUtil_GetFunctionParameters(function)
	// Strip as many unreferenced args off the end of the function declaration as possible.
	r.MaybeRemoveUnusedTrailingParameters(argList, fparamScope)

	// Mark any remaining unused parameters are unused to OptimizeParameters can try to remove
	// them.
	r.MarkUnusedParameters(argList, fparamScope)
}

func (r *RemoveUnusedCode) isPropertyPinned(propertyName string) bool {
	if isImplicitlyUsedProperty(propertyName) {
		return true
	}
	if _, ok := r.pinnedPropertyNames[propertyName]; ok {
		return true
	}
	if r.externProperties != nil && r.externProperties.Contains(propertyName) {
		r.pinnedPropertyNames[propertyName] = struct{}{}
		return true
	}
	return false
}

// Java: RemoveUnusedCode.markPropertyNameAsPinned(String)
func (r *RemoveUnusedCode) MarkPropertyNameAsPinned(propertyName string) {
	if r.isPropertyPinned(propertyName) {
		return
	}
	r.pinnedPropertyNames[propertyName] = struct{}{}
	// Continue traversal of all of the property name's values and no longer consider them for
	// removal.
	if list, ok := r.propRemovables[propertyName]; ok {
		delete(r.propRemovables, propertyName)
		for cur := list.head; cur != nil; {
			next := cur.propNext
			cur.propNext = nil
			cur.flags &^= remFlagInPropList
			cur.ApplyContinuations(r)
			cur = next
		}
	}
}

func (r *RemoveUnusedCode) isExportedByConvention(name string, local bool) bool {
	if len(name) == 0 || (name[0] != '_' && name[0] != '$') {
		return false
	}
	return r.codingConvention.IsExportedWithLocal(name, local)
}

// Java: RemoveUnusedCode.considerForIndependentRemoval(Removable)
func (r *RemoveUnusedCode) ConsiderForIndependentRemoval(removable *RemoveUnusedCode_RemovableBase) {
	if removable.IsNamedProperty() {
		propertyName := removable.GetPropertyName()

		if r.isPropertyPinned(propertyName) ||
			r.isExportedByConvention(propertyName /* local= */, false) {
			// Referenced or exported, so not removable.
			removable.ApplyContinuations(r)
		} else if r.IsIndependentlyRemovable(removable) {
			// Store for possible removal later.
			if removable.flags&remFlagInPropList == 0 {
				removable.flags |= remFlagInPropList
				removable.propNext = nil
				list, exists := r.propRemovables[propertyName]
				if !exists {
					r.propRemovableKeys = append(r.propRemovableKeys, propertyName)
					r.propRemovables[propertyName] = propRemovableList{head: removable, tail: removable}
				} else {
					list.tail.propNext = removable
					list.tail = removable
					r.propRemovables[propertyName] = list
				}
			}
		} else {
			removable.ApplyContinuations(r)
			// This assignment counts as a reference, since we won't be removing it.
			// This is necessary in order to preserve getters and setters for the property.
			r.MarkPropertyNameAsPinned(propertyName)
		}
	} else {
		removable.ApplyContinuations(r)
	}
}

func (r *RemoveUnusedCode) ensureAccessorMasks(summary *AccessorSummary) {
	r.lastAccessorSummary = summary
	r.hasGetterAccessors = false
	r.hasSetterAccessors = false
	r.getterLenMask = [256]uint64{}
	r.getterEndMask = [256]uint64{}
	r.setterLenMask = [256]uint64{}
	r.setterEndMask = [256]uint64{}
	if summary.assumeAlwaysGetterAndSetter {
		r.hasGetterAccessors = true
		r.hasSetterAccessors = true
		return
	}
	for name, kind := range summary.accessors.All() {
		var shift uint
		var b0, bEnd byte
		if len(name) > 0 {
			shift = uint(len(name))
			if shift >= 64 {
				shift = 63
			}
			b0 = name[0]
			bEnd = name[len(name)-1]
		}
		mask := uint64(1) << shift
		if kind.HasGetter() {
			r.hasGetterAccessors = true
			r.getterLenMask[b0] |= mask
			r.getterEndMask[bEnd] |= mask
		}
		if kind.HasSetter() {
			r.hasSetterAccessors = true
			r.setterLenMask[b0] |= mask
			r.setterEndMask[bEnd] |= mask
		}
	}
}

// ConsiderForAccessorSideEffects returns whether or not accessor
// side-effect are a possibility.
//
// Java: RemoveUnusedCode.considerForAccessorSideEffects(Node, PropertyAccessKind)
func (r *RemoveUnusedCode) ConsiderForAccessorSideEffects(getprop rhino.Node, usage AccessorSummary_PropertyAccessKind) bool {
	// Other node types may make sense in the future.
	if !NodeUtil_IsNormalOrOptChainGetProp(getprop) {
		jutil.CheckStateMsg1(false, "%s", getprop)
	}
	if r.assumeGettersArePure && !usage.HasSetter() {
		return false
	}
	summary := r.accessorSummary
	if summary == nil {
		summary = r.compiler.GetAccessorSummary()
		r.accessorSummary = summary
	}
	if summary != r.lastAccessorSummary {
		r.ensureAccessorMasks(summary)
	}
	propName := getprop.GetString()
	if !summary.assumeAlwaysGetterAndSetter {
		var shift uint
		var b0, bEnd byte
		if len(propName) > 0 {
			shift = uint(len(propName))
			if shift >= 64 {
				shift = 63
			}
			b0 = propName[0]
			bEnd = propName[len(propName)-1]
		}
		mask := uint64(1) << shift
		mightMatch := false
		if usage.HasSetter() && r.hasSetterAccessors &&
			(r.setterLenMask[b0]&mask) != 0 && (r.setterEndMask[bEnd]&mask) != 0 {
			mightMatch = true
		} else if !r.assumeGettersArePure && usage.HasGetter() && r.hasGetterAccessors &&
			(r.getterLenMask[b0]&mask) != 0 && (r.getterEndMask[bEnd]&mask) != 0 {
			mightMatch = true
		}
		if !mightMatch {
			return false
		}
	}
	recorded := summary.GetKind(propName)
	if (recorded.HasGetter() && usage.HasGetter() && !r.assumeGettersArePure) ||
		(recorded.HasSetter() && usage.HasSetter()) {
		r.MarkPropertyNameAsPinned(propName)
		return true
	}

	return false
}

// Java: RemoveUnusedCode.isIndependentlyRemovable(Removable)
func (r *RemoveUnusedCode) IsIndependentlyRemovable(removable *RemoveUnusedCode_RemovableBase) bool {
	if removable.IsPrototypeProperty() {
		// `foo.prototype.prop = something;`
		// `class C { prop() {} }`
		return r.removeUnusedPrototypeProperties
	} else if removable.IsObjectDefinePropertiesDefinition() {
		// `Object.defineProperties({ prop: {...}});`
		return r.removeUnusedObjectDefinePropertiesDefinitions
	} else if removable.IsThisDotPropertyReference() {
		// `this.prop = something;`
		return r.removeUnusedThisProperties
	} else if removable.IsStaticProperty() {
		// `class Foo { static prop() {} }`
		// `Foo.otherStaticProp = value;`
		// TODO(b/139319709): removeUnusedThisProperties has ended up covering more than it was
		// originally intended to cover for arbitrary reasons.
		return r.removeUnusedThisProperties
	} else {
		return false
	}
}

// MarkUnusedParameters marks any remaining unused parameters as being unused
// so it can be used elsewhere.
//
// paramList is the list of function's parameters.
//
// Java: RemoveUnusedCode.markUnusedParameters(Node, Scope)
func (r *RemoveUnusedCode) MarkUnusedParameters(paramList rhino.Node, fparamScope *Scope) {
	if !paramList.IsParamList() {
		jutil.CheckArgumentMsg1(false, "%s", paramList)
	}

	for param := paramList.GetFirstChild(); !param.IsNil(); param = param.GetNext() {
		if param.IsUnusedParameter() {
			// already marked
			continue
		}

		paramNameNode := RemoveUnusedCode_NameOfParam(param)
		if paramNameNode.IsNil() {
			// destructuring pattern parameters don't have a name that applies to the whole parameter
			// TODO(bradfordcsmith): We could mark this if we determined that all vars created by
			// the pattern are unused.
			continue
		}

		varInfo := r.TraverseNameNode(paramNameNode, fparamScope)
		if varInfo.IsRemovable() {
			param.SetUnusedParameter(true)
			r.compiler.ReportChangeToEnclosingScope(paramList)
			if r.removalLog.IsLogging() {
				r.removalLog.LogSupplier(RemoveUnusedCode_RemovalLogRecord_ForMarkingNamedArg(paramNameNode, paramList).Get)
			}
		}
	}
}

// MaybeRemoveUnusedTrailingParameters strips as many unreferenced args off
// the end of the function declaration as possible. We start from the end of
// the function declaration because removing parameters from the middle of
// the param list could mess up the interpretation of parameters being sent
// over by any function calls.
//
// argList is the list of function's arguments.
//
// Java: RemoveUnusedCode.maybeRemoveUnusedTrailingParameters(Node, Scope)
func (r *RemoveUnusedCode) MaybeRemoveUnusedTrailingParameters(argList rhino.Node, fparamScope *Scope) {
	if !argList.IsParamList() {
		jutil.CheckArgumentMsg1(false, "%s", argList)
	}
	for {
		lastArg := argList.GetLastChild()
		if lastArg.IsNil() {
			break
		}
		argNode := lastArg
		if lastArg.IsDefaultValue() {
			argNode = lastArg.GetFirstChild()
			if r.MayHaveSideEffects(lastArg.GetLastChild()) {
				break
			}
		}

		if argNode.IsRest() {
			argNode = argNode.GetFirstChild()
		}

		if argNode.IsDestructuringPattern() {
			if argNode.HasChildren() {
				// TODO(johnlenz): handle the case where there are no assignments.
				break
			} else {
				// Remove empty destructuring patterns and their associated object literal assignment
				// if it exists and if the right hand side does not have side effects. Note, a
				// destructuring pattern with a "leftover" property key as in {a:{}} is not considered
				// empty in this case!
				NodeUtil_DeleteNode(lastArg, r.compiler)
				if r.removalLog.IsLogging() {
					r.removalLog.LogSupplier(RemoveUnusedCode_RemovalLogRecord_ForDestructuringArg(argList).Get)
				}
				continue
			}
		}

		varInfo := r.GetVarInfo(r.GetVarForNameNode(argNode, fparamScope))
		if varInfo.IsRemovable() {
			NodeUtil_DeleteNode(lastArg, r.compiler)
			if r.removalLog.IsLogging() {
				r.removalLog.LogSupplier(RemoveUnusedCode_RemovalLogRecord_ForNamedArg(argNode, argList).Get)
			}
		} else {
			break
		}
	}
}

// TraverseVar handles a variable reference seen during traversal and
// returns a VarInfo object appropriate for the given Var.
//
// This is a wrapper for getVarInfo that handles additional logic needed when
// we're getting the VarInfo during traversal.
//
// Java: RemoveUnusedCode.traverseVar(Var)
func (r *RemoveUnusedCode) TraverseVar(var_ *Var) *RemoveUnusedCode_RealVarInfo {
	if var_ == nil {
		panic(jutil.NewNullPointerException())
	}
	if r.removeLocalVars && var_.IsArguments() {
		// If we are considering removing local variables, that includes parameters.
		// If `arguments` is used in a function we must consider all parameters to be referenced.
		varScope, _ := var_.GetScope().(*Scope)
		functionScope := varScope.GetClosestHoistScope()
		paramList := NodeUtil_GetFunctionParameters(functionScope.GetRootNode())
		for param := paramList.GetFirstChild(); !param.IsNil(); param = param.GetNext() {
			lValue := RemoveUnusedCode_NameOfParam(param)
			if lValue.IsNil() {
				continue
			}

			r.GetVarInfo(r.GetVarForNameNode(lValue, functionScope)).
				SetIsExplicitlyNotRemovable(func() string { return "parameter in function using arguments" })
		}
		// `arguments` is never removable.
		return r.canonicalUnremovableVarInfo
	} else {
		return r.GetVarInfo(var_)
	}
}

// RemoveUnusedCode_NameOfParam returns the NAME node associated with a
// function parameter (the child of a PARAM_LIST), or null if there is no
// single name.
//
// Java: RemoveUnusedCode.nameOfParam(Node)
func RemoveUnusedCode_NameOfParam(param rhino.Node) rhino.Node {
	switch param.GetToken() {
	case rhino.Token_NAME:
		return param
	case rhino.Token_DEFAULT_VALUE:
		return RemoveUnusedCode_NameOfParam(param.GetFirstChild())
	case rhino.Token_ITER_REST:
		return RemoveUnusedCode_NameOfParam(param.GetOnlyChild())
	case rhino.Token_ARRAY_PATTERN, rhino.Token_OBJECT_PATTERN:
		return rhino.Node{}
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected child of PARAM_LIST: " + param.ToStringTree()))
	}
}

// GetVarInfo gets the right VarInfo object to use for the given Var.
//
// This method is responsible for managing the entries in varInfoMap.
//
// Note: Several Vars may share the same VarInfo when they should be treated
// the same way.
//
// Java: RemoveUnusedCode.getVarInfo(Var)
func (r *RemoveUnusedCode) GetVarInfo(var_ *Var) *RemoveUnusedCode_RealVarInfo {
	if var_ == nil {
		panic(jutil.NewNullPointerException())
	}
	isGlobal := var_.IsGlobal()
	if var_.IsExtern() {
		if r.unremovableLog.IsLogging() {
			name := var_.GetName()
			r.unremovableLog.LogSupplier(func() string { return name + ": extern" })
		}
		return r.canonicalUnremovableVarInfo
	} else if name := var_.GetName(); r.isExportedByConvention(name /* local= */, !isGlobal) {
		if r.unremovableLog.IsLogging() {
			r.unremovableLog.LogSupplier(func() string { return name + ": exported by convention" })
		}
		return r.canonicalUnremovableVarInfo
	} else if var_.IsArguments() {
		// No point in logging that we cannot remove "arguments"
		return r.canonicalUnremovableVarInfo
	}
	scopeRoot := var_.GetScopeRoot()
	h := rucVarHash(scopeRoot, var_.nameID)
	for v := r.varBuckets[h&uint32(len(r.varBuckets)-1)]; v != nil; v = v.hashNext {
		if v.hash == h && v.var_.nameID == var_.nameID && v.var_.GetScopeRoot() == scopeRoot {
			return v
		}
	}
	varInfo := NewRemoveUnusedCode_RealVarInfo(r, var_)
	varInfo.hash = h
	if var_.GetParentNode().IsParamList() {
		varInfo.SetHasNonLocalOrNonLiteralValue()
	}
	// Cannot use canonicalUnremovableVarInfo for the 2 non-removable cases below, because each
	// varInfo needs to track what value is assigned to it for the purpose of correctly allowing
	// or preventing removal of properties set on it.
	if !r.removeGlobals && isGlobal {
		varInfo.SetIsExplicitlyNotRemovable(func() string { return "not removing globals" })
	} else if !r.removeLocalVars && !isGlobal {
		varInfo.SetIsExplicitlyNotRemovable(func() string { return "not removing locals" })
	}
	r.putVarInfo(varInfo)
	return varInfo
}

func (r *RemoveUnusedCode) putVarInfo(v *RemoveUnusedCode_RealVarInfo) {
	r.varInfos = appendDoubling(r.varInfos, v)
	if len(r.varInfos)*4 > len(r.varBuckets)*3 {
		buckets := make([]*RemoveUnusedCode_RealVarInfo, 2*len(r.varBuckets))
		mask := uint32(len(buckets) - 1)
		for _, info := range r.varInfos {
			info.hashNext = buckets[info.hash&mask]
			buckets[info.hash&mask] = info
		}
		r.varBuckets = buckets
		return
	}
	i := v.hash & uint32(len(r.varBuckets)-1)
	v.hashNext = r.varBuckets[i]
	r.varBuckets[i] = v
}

// RemoveUnreferencedVarsAndPolyfills removes any vars in the scope that were
// not referenced. Removes any assignments to those variables as well.
//
// Java: RemoveUnusedCode.removeUnreferencedVarsAndPolyfills()
func (r *RemoveUnusedCode) RemoveUnreferencedVarsAndPolyfills() {
	for _, varInfo := range r.varInfos {
		if !varInfo.IsRemovable() {
			continue
		}
		var_ := varInfo.var_

		if r.removalLog.IsLogging() {
			r.removalLog.LogSupplier(RemoveUnusedCode_RemovalLogRecord_ForVar(var_).Get)
		}
		// Regardless of what happens to the original declaration,
		// we need to remove all assigns, because they may contain references
		// to other unreferenced variables.
		varInfo.RemoveAllRemovables()

		nameNode := var_.GetNameNode()
		toRemove := nameNode.GetParent()
		if toRemove.IsNil() || RemoveUnusedCode_AlreadyRemoved(toRemove) {
			// assignedVarInfo.removeAllRemovables () already removed it
		} else if NodeUtil_IsFunctionExpression(toRemove) {
			// TODO(bradfordcsmith): Add a Removable for this case.
			if !r.preserveFunctionExpressionNames {
				fnNameNode := toRemove.GetFirstChild()
				r.compiler.ReportChangeToEnclosingScope(fnNameNode)
				fnNameNode.SetString("")
			}
		} else {
			// Removables are not created for theses cases.
			// function foo(unused1 = someSideEffectingValue, ...unused2) {}
			// removeUnreferencedFunctionArgs() is responsible for removing these.
			// TODO(bradfordcsmith): handle parameter declarations with removables
			jutil.CheckStateMsg1(
				toRemove.IsParamList() ||
					(toRemove.GetParent().IsParamList() &&
						(toRemove.IsDefaultValue() || toRemove.IsRest())),
				"unremoved code: %s",
				toRemove)
		}
	}

	// Iterator<PolyfillInfo> iter = polyfills.values().iterator(); ... iter.remove();
	//
	// Go-only: language-adapter — the removable polyfills are collected in
	// iteration order and removed from the multimap after the iteration
	// (the JavaHashMap iterators are fail-fast); removing a polyfill's code
	// does not touch the multimap, so the order of effects is unchanged.
	if r.polyfills != nil && !r.polyfills.IsEmpty() {
		var removed []RemoveUnusedCode_PolyfillInfo
		for polyfill := range r.polyfills.Values() {
			if polyfill.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable {
				if r.removalLog.IsLogging() {
					r.removalLog.LogSupplier(RemoveUnusedCode_RemovalLogRecord_ForPolyfill(polyfill).Get)
				}
				polyfill.GetRemoveUnusedCode_PolyfillInfoBase().removable.Remove(r)
				removed = append(removed, polyfill)
			}
		}
		for _, polyfill := range removed {
			r.polyfills.Remove(polyfill.GetRemoveUnusedCode_PolyfillInfoBase().key, polyfill)
		}
		if len(removed) > 0 {
			r.hasPolyfills = false
			r.polyfillLenMask = [256]uint64{}
			r.polyfillEndMask = [256]uint64{}
			for key := range r.polyfills.KeySet() {
				r.addPolyfillMask(key)
			}
		}
	}
}

// RemoveUnusedCode_Continuation: our progress in a traversal can be
// expressed completely as the current node and scope. The continuation lets
// us save that information so that we can continue the traversal later.
type RemoveUnusedCode_Continuation struct {
	node  rhino.Node
	scope *Scope
}

// Java: RemoveUnusedCode.Continuation(Node, Scope)
func NewRemoveUnusedCode_Continuation(_ *RemoveUnusedCode, node rhino.Node, scope *Scope) RemoveUnusedCode_Continuation {
	return RemoveUnusedCode_Continuation{node: node, scope: scope}
}

// Java: RemoveUnusedCode.Continuation.apply()
func (c RemoveUnusedCode_Continuation) Apply(outer *RemoveUnusedCode) {
	if c.node.IsFunction() {
		// Calling traverseNode here would create infinite recursion for a function declaration
		outer.TraverseFunction(c.node, c.scope)
	} else {
		outer.TraverseNode(c.node, c.scope)
	}
}

// RemoveUnusedCode_Removable represents a portion of the AST that can be
// removed.
type RemoveUnusedCode_Removable = *RemoveUnusedCode_RemovableBase

const (
	remKindUnusedReadRef uint8 = iota + 1
	remKindInstanceofName
	remKindWeakUsageCall
	remKindIncOrDecOp
	remKindIndirectAssign
	remKindPolyfill
	remKindClassDecl
	remKindNamedClassExpr
	remKindClassOrProtoProp
	remKindObjectDefinePropsDef
	remKindFunctionDecl
	remKindNameDeclStmt
	remKindAssignVar
	remKindAssignNamedProp
	remKindAssignComputedProp
	remKindAnonProtoPropAssign
	remKindClassSetupCall
	remKindVanillaForNameDecl
)

const (
	remFlagHasPropertyName uint8 = 1 << iota
	remFlagProtoDotPropRef
	remFlagThisDotPropRef
	remFlagContsApplied
	remFlagRemoved
	remFlagInPropList
	remFlagUninitVarDecl
	remFlagSubBit // isPatch for Polyfill, classDefiningCall for ClassSetupCall
)

// RemoveUnusedCode_RemovableBase holds all Removable state in a single 64-byte
// cache line (8 words, 6 pointers).
type RemoveUnusedCode_RemovableBase struct {
	propNext *RemoveUnusedCode_RemovableBase

	cont0 RemoveUnusedCode_Continuation

	// Primary AST node for the removable construct.
	targetNode rhino.Node

	// Auxiliary AST node: root for IndirectAssign, toPreserve for IncOrDecOp.
	auxNode rhino.Node

	// Associated VarInfo for Assign.
	varInfo *RemoveUnusedCode_RealVarInfo

	propertyNameID uint32
	extraContIdx   uint32
	numConts       uint16
	remKind        uint8
	flags          uint8
}

func (b *RemoveUnusedCode_RemovableBase) initFromBuilder(remKind uint8, targetNode rhino.Node, builder *RemoveUnusedCode_RemovableBuilder) {
	b.cont0 = builder.cont0
	b.targetNode = targetNode
	b.propertyNameID = builder.propertyNameID
	b.extraContIdx = builder.extraContIdx
	b.numConts = builder.numConts
	b.remKind = remKind
	b.flags = builder.flags
}

// GetRemoveUnusedCode_RemovableBase gives access to the base.
func (b *RemoveUnusedCode_RemovableBase) GetRemoveUnusedCode_RemovableBase() *RemoveUnusedCode_RemovableBase {
	return b
}

// Java: RemoveUnusedCode.Removable.getPropertyName()
func (b *RemoveUnusedCode_RemovableBase) GetPropertyName() string {
	if b.remKind == remKindIndirectAssign {
		jutil.CheckStateMsg1(b.targetNode.IsGetProp(), "%s", b.targetNode)
		return b.targetNode.GetString()
	}
	if b.flags&remFlagHasPropertyName == 0 {
		panic(jutil.NewNullPointerException())
	}
	return rhino.RhinoStringPool_StringFromID(b.propertyNameID)
}

// Remove removes the associated nodes from the AST, unless they've already
// been removed.
//
// Java: RemoveUnusedCode.Removable.remove(AbstractCompiler)
func (b *RemoveUnusedCode_RemovableBase) Remove(r *RemoveUnusedCode) {
	if b.flags&remFlagRemoved == 0 {
		b.flags |= remFlagRemoved
		b.RemoveInternal(r)
	}
}

// Java: RemoveUnusedCode.Removable.applyContinuations()
func (b *RemoveUnusedCode_RemovableBase) ApplyContinuations(r *RemoveUnusedCode) {
	if b.flags&remFlagContsApplied == 0 {
		b.flags |= remFlagContsApplied
		// Enqueue the continuation for processing.
		// Don't invoke the continuation immediately, because that can lead to concurrent
		// modification of data structures.
		if b.numConts >= 1 {
			r.enqueueContinuation(b.cont0)
			b.cont0 = RemoveUnusedCode_Continuation{}
		}
		if b.numConts >= 2 {
			end := b.extraContIdx + uint32(b.numConts-1)
			extras := r.extraContSlab[b.extraContIdx:end]
			for _, c := range extras {
				r.enqueueContinuation(c)
			}
			clear(extras)
		}
		b.numConts = 0
	}
}

// IsVariableAssignment: true if this object represents assignment to a
// variable.
//
// Java: RemoveUnusedCode.Removable.isVariableAssignment()
func (b *RemoveUnusedCode_RemovableBase) IsVariableAssignment() bool {
	switch b.remKind {
	case remKindAssignVar, remKindClassDecl, remKindFunctionDecl, remKindNameDeclStmt, remKindVanillaForNameDecl:
		return true
	case remKindIndirectAssign:
		return b.targetNode.IsName()
	default:
		return false
	}
}

// IsNamedProperty: true if this object represents a named property, either
// assignment or declaration.
//
// Java: RemoveUnusedCode.Removable.isNamedProperty()
func (b *RemoveUnusedCode_RemovableBase) IsNamedProperty() bool {
	if b.remKind == remKindIndirectAssign {
		return b.targetNode.IsGetProp()
	}
	return b.flags&remFlagHasPropertyName != 0
}

// IsNamedPropertyAssignment: true if this object represents assignment to a
// named property.
//
// This does not include class or object literal member declarations.
//
// Java: RemoveUnusedCode.Removable.isNamedPropertyAssignment()
func (b *RemoveUnusedCode_RemovableBase) IsNamedPropertyAssignment() bool {
	switch b.remKind {
	case remKindAssignNamedProp:
		return true
	case remKindIndirectAssign:
		return b.targetNode.IsGetProp()
	default:
		return false
	}
}

// Java: RemoveUnusedCode.Removable.isAssignedValueLocal()
func (b *RemoveUnusedCode_RemovableBase) IsAssignedValueLocal() bool {
	switch b.remKind {
	case remKindClassDecl, remKindFunctionDecl:
		return true
	case remKindNameDeclStmt:
		nameNode := b.targetNode.GetOnlyChild()
		initialValueNode := nameNode.GetFirstChild()
		if initialValueNode.IsNil() {
			return true
		}
		valueNode := RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(nameNode, initialValueNode)
		return NodeUtil_EvaluatesToLocalValue(valueNode)
	case remKindVanillaForNameDecl:
		initialValueNode := b.targetNode.GetFirstChild()
		if initialValueNode.IsNil() {
			return true
		}
		valueNode := RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(b.targetNode, initialValueNode)
		return NodeUtil_EvaluatesToLocalValue(valueNode)
	case remKindAssignVar, remKindAssignNamedProp, remKindAssignComputedProp:
		return !b.GetLocalAssignedValue().IsNil()
	default:
		return false
	}
}

// GetLocalAssignedValue returns the Node representing the local value that
// is being assigned or `null` if the value is non-local or cannot be
// determined.
//
// Java: RemoveUnusedCode.Removable.getLocalAssignedValue()
func (b *RemoveUnusedCode_RemovableBase) GetLocalAssignedValue() rhino.Node {
	switch b.remKind {
	case remKindClassDecl, remKindFunctionDecl:
		return b.targetNode
	case remKindNameDeclStmt:
		nameNode := b.targetNode.GetOnlyChild()
		initialValueNode := nameNode.GetFirstChild()
		if initialValueNode.IsNil() {
			return rhino.Node{}
		}
		valueNode := RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(nameNode, initialValueNode)
		if NodeUtil_EvaluatesToLocalValue(valueNode) {
			return valueNode
		}
		return rhino.Node{}
	case remKindVanillaForNameDecl:
		initialValueNode := b.targetNode.GetFirstChild()
		if initialValueNode.IsNil() {
			return rhino.Node{}
		}
		valueNode := RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(b.targetNode, initialValueNode)
		if NodeUtil_EvaluatesToLocalValue(valueNode) {
			return valueNode
		}
		return rhino.Node{}
	case remKindAssignVar, remKindAssignNamedProp, remKindAssignComputedProp:
		if NodeUtil_IsExpressionResultUsed(b.targetNode) {
			return rhino.Node{}
		}
		valueNode := RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(
			b.targetNode.GetFirstChild(), b.targetNode.GetLastChild())
		if NodeUtil_EvaluatesToLocalValue(valueNode) {
			return valueNode
		}
		return rhino.Node{}
	default:
		return rhino.Node{}
	}
}

// IsPrototypeAssignment: is this a direct assignment to `varName.prototype`?
//
// Java: RemoveUnusedCode.Removable.isPrototypeAssignment()
func (b *RemoveUnusedCode_RemovableBase) IsPrototypeAssignment() bool {
	if !b.IsNamedPropertyAssignment() {
		return false
	}
	return b.GetPropertyName() == "prototype"
}

// IsPrototypeDotPropertyReference: is this an assignment to a property on a
// prototype object?
//
// Java: RemoveUnusedCode.Removable.isPrototypeDotPropertyReference()
func (b *RemoveUnusedCode_RemovableBase) IsPrototypeDotPropertyReference() bool {
	return b.flags&remFlagProtoDotPropRef != 0
}

// Java: RemoveUnusedCode.Removable.isClassOrPrototypeNamedProperty()
func (b *RemoveUnusedCode_RemovableBase) IsClassOrPrototypeNamedProperty() bool {
	return b.remKind == remKindClassOrProtoProp && !b.targetNode.IsStaticMember()
}

// Java: RemoveUnusedCode.Removable.isPrototypeProperty()
func (b *RemoveUnusedCode_RemovableBase) IsPrototypeProperty() bool {
	if b.remKind == remKindAnonProtoPropAssign {
		return true
	}
	return b.IsPrototypeDotPropertyReference() || b.IsClassOrPrototypeNamedProperty()
}

// Java: RemoveUnusedCode.Removable.isThisDotPropertyReference()
func (b *RemoveUnusedCode_RemovableBase) IsThisDotPropertyReference() bool {
	if b.remKind == remKindIndirectAssign {
		return RemoveUnusedCode_IsThisDotProperty(b.targetNode)
	}
	return b.flags&remFlagThisDotPropRef != 0
}

// Java: RemoveUnusedCode.Removable.isObjectDefinePropertiesDefinition()
func (b *RemoveUnusedCode_RemovableBase) IsObjectDefinePropertiesDefinition() bool {
	return b.remKind == remKindObjectDefinePropsDef
}

// IsStaticProperty.
//
// TODO(b/134610338): Combine this method with `isPrototypeProperty`.
//
// Java: RemoveUnusedCode.Removable.isStaticProperty()
func (b *RemoveUnusedCode_RemovableBase) IsStaticProperty() bool {
	switch b.remKind {
	case remKindClassOrProtoProp:
		return b.targetNode.IsStaticMember()
	case remKindAssignNamedProp:
		if b.varInfo != nil && b.varInfo.HasFunctionOrClassLiteralValue() {
			return b.targetNode.GetFirstFirstChild().IsName()
		}
		return false
	default:
		return false
	}
}

// PreventsRemovalOfVariableWithNonLocalValueOrPrototype: would a nonlocal
// or nonliteral value prevent removal of a variable associated with this
// Removable?
//
// Java: RemoveUnusedCode.Removable.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
func (b *RemoveUnusedCode_RemovableBase) PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool {
	switch b.remKind {
	case remKindInstanceofName, remKindClassSetupCall, remKindAssignNamedProp, remKindAssignComputedProp:
		return true
	case remKindIndirectAssign:
		if b.targetNode.IsGetProp() {
			getPropLhs := b.targetNode.GetFirstChild()
			return getPropLhs.IsName() || (getPropLhs.IsGetProp() && getPropLhs.GetFirstChild().IsName() && getPropLhs.GetString() == "prototype")
		}
		return false
	default:
		return false
	}
}

func (b *RemoveUnusedCode_RemovableBase) RemoveInternal(r *RemoveUnusedCode) {
	compiler := r.compiler
	switch b.remKind {
	case remKindUnusedReadRef:
		if !RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			if RemoveUnusedCode_IsThisDotProperty(b.targetNode) {
				r.RemoveExpressionCompletely(b.targetNode)
			} else {
				if !RemoveUnusedCode_IsDotPrototypeDotProperty(b.targetNode) {
					jutil.CheckStateMsg1(false, "%s", b.targetNode)
				}
				// objExpression.prototype.propertyName
				objExpression := b.targetNode.GetFirstFirstChild()
				if r.MayHaveSideEffects(objExpression) {
					r.ReplaceNodeWith(b.targetNode, objExpression.Detach())
				} else {
					r.RemoveExpressionCompletely(b.targetNode)
				}
			}
		}
	case remKindInstanceofName:
		if !RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			lhs := b.targetNode.GetFirstChild()
			falseNode := rhino.IR_FalseNode().Srcref(b.targetNode)
			if r.MayHaveSideEffects(lhs) {
				r.ReplaceNodeWith(b.targetNode, rhino.IR_Comma(lhs.Detach(), falseNode).Srcref(b.targetNode))
			} else {
				r.ReplaceNodeWith(b.targetNode, falseNode)
			}
		}
	case remKindWeakUsageCall:
		if !RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			undefinedNode := NodeUtil_NewUndefinedNode(b.targetNode)
			r.ReplaceNodeWith(b.targetNode, undefinedNode)
		}
	case remKindIncOrDecOp:
		if RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			return
		}
		arg := b.targetNode.GetOnlyChild()
		jutil.CheckStateMsg1(arg.IsGetProp(), "%s", arg)
		if b.auxNode.IsNil() {
			r.RemoveExpressionCompletely(b.targetNode)
		} else {
			r.ReplaceNodeWith(b.targetNode, b.auxNode.Detach())
		}
	case remKindIndirectAssign:
		if !RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			b.removeIndirectAssignRoot(r)
		}
	case remKindPolyfill, remKindClassDecl, remKindClassOrProtoProp, remKindObjectDefinePropsDef, remKindFunctionDecl:
		NodeUtil_DeleteNode(b.targetNode, compiler)
	case remKindNamedClassExpr:
		if !RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			nameNode := b.targetNode.GetFirstChild()
			if !nameNode.IsEmpty() {
				nameNode.ReplaceWith(rhino.IR_Empty().Srcref(nameNode))
				compiler.ReportChangeToEnclosingScope(b.targetNode)
			}
		}
	case remKindNameDeclStmt:
		nameNode := b.targetNode.GetOnlyChild()
		valueNode := nameNode.GetFirstChild()
		if !valueNode.IsNil() && r.MayHaveSideEffects(valueNode) {
			compiler.ReportChangeToEnclosingScope(b.targetNode)
			valueNode.Detach()
			b.targetNode.ReplaceWith(rhino.IR_ExprResult(valueNode).Srcref(valueNode))
		} else {
			NodeUtil_DeleteNode(b.targetNode, compiler)
		}
	case remKindAssignVar, remKindAssignNamedProp, remKindAssignComputedProp:
		if RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			return
		}
		parent := b.targetNode.GetParent()
		compiler.ReportChangeToEnclosingScope(parent)
		lhs := b.targetNode.GetFirstChild()
		rhs := b.targetNode.GetSecondChild()
		mustPreserveRhs :=
			r.MayHaveSideEffects(rhs) || NodeUtil_IsExpressionResultUsed(b.targetNode)
		mustPreserveGetElmExpr := lhs.IsGetElem() && r.MayHaveSideEffects(lhs.GetLastChild())

		if mustPreserveRhs && mustPreserveGetElmExpr {
			replacement := rhino.IR_Comma(lhs.GetLastChild().Detach(), rhs.Detach()).Srcref(b.targetNode)
			r.ReplaceNodeWith(b.targetNode, replacement)
		} else if mustPreserveGetElmExpr {
			r.ReplaceNodeWith(b.targetNode, lhs.GetLastChild().Detach())
		} else if mustPreserveRhs {
			r.ReplaceNodeWith(b.targetNode, rhs.Detach())
		} else {
			r.RemoveExpressionCompletely(b.targetNode)
		}
	case remKindAnonProtoPropAssign:
		if RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			return
		}
		parent := b.targetNode.GetParent()
		compiler.ReportChangeToEnclosingScope(parent)
		lhs := b.targetNode.GetFirstChild()
		rhs := b.targetNode.GetLastChild()

		jutil.CheckStateMsg1(lhs.IsGetProp(), "%s", lhs)
		objDotPrototype := lhs.GetFirstChild()
		jutil.CheckStateMsg1(objDotPrototype.IsGetProp(), "%s", objDotPrototype)
		objExpression := objDotPrototype.GetFirstChild()
		jutil.CheckStateMsg1(objDotPrototype.GetString() == "prototype", "%s", objDotPrototype)

		mustPreserveRhs :=
			r.MayHaveSideEffects(rhs) || NodeUtil_IsExpressionResultUsed(b.targetNode)
		mustPreserveObjExpression := r.MayHaveSideEffects(objExpression)

		if mustPreserveRhs && mustPreserveObjExpression {
			replacement := rhino.IR_Comma(objExpression.Detach(), rhs.Detach()).Srcref(b.targetNode)
			r.ReplaceNodeWith(b.targetNode, replacement)
		} else if mustPreserveObjExpression {
			r.ReplaceNodeWith(b.targetNode, objExpression.Detach())
		} else if mustPreserveRhs {
			r.ReplaceNodeWith(b.targetNode, rhs.Detach())
		} else {
			r.RemoveExpressionCompletely(b.targetNode)
		}
	case remKindClassSetupCall:
		parent := b.targetNode.GetParent()
		classDefiningCall := b.flags&remFlagSubBit != 0

		var replacement rhino.Node
		b.targetNode.RemoveFirstChild()
		for arg := b.targetNode.GetLastChild(); !arg.IsNil(); arg = b.targetNode.GetLastChild() {
			arg.Detach()
			if !classDefiningCall && r.MayHaveSideEffects(arg) {
				if replacement.IsNil() {
					replacement = arg
				} else {
					replacement = rhino.IR_Comma(arg, replacement).Srcref(b.targetNode)
				}
			} else {
				NodeUtil_MarkFunctionsDeleted(arg, compiler)
			}
		}

		if !replacement.IsNil() {
			r.ReplaceNodeWith(b.targetNode, replacement)
		} else if parent.IsExprResult() {
			NodeUtil_DeleteNode(parent, compiler)
		} else {
			if parent.IsComma() || parent.IsAnd() || parent.IsOr() {
				if parent.GetFirstChild() == b.targetNode {
					rhs := rhino.CheckNotNull(b.targetNode.GetNext())
					compiler.ReportChangeToEnclosingScope(parent)
					parent.ReplaceWith(rhs.Detach())
				} else {
					lhs := parent.GetFirstChild()
					compiler.ReportChangeToEnclosingScope(parent)
					parent.ReplaceWith(lhs.Detach())
				}
			} else {
				compiler.ReportChangeToEnclosingScope(parent)
				b.targetNode.ReplaceWith(rhino.IR_Number(0))
			}
		}
	case remKindVanillaForNameDecl:
		declaration := rhino.CheckNotNull(b.targetNode.GetParent())
		compiler.ReportChangeToEnclosingScope(declaration)
		if b.targetNode.GetPrevious().IsNil() && b.targetNode.GetNext().IsNil() {
			declaration.ReplaceWith(rhino.IR_Empty().Srcref(declaration))
		} else {
			b.targetNode.Detach()
		}
		NodeUtil_MarkFunctionsDeleted(b.targetNode, compiler)
	}
}

func (b *RemoveUnusedCode_RemovableBase) removeIndirectAssignRoot(r *RemoveUnusedCode) {
	root := b.auxNode
	rootParent := root.GetParent()

	switch rootParent.GetToken() {
	case rhino.Token_ARRAY_PATTERN:
		r.ReplaceNodeWith(root, rhino.IR_Empty().Srcref(root))
		for maybeEmpty := rootParent.GetLastChild(); !maybeEmpty.IsNil() && maybeEmpty.IsEmpty(); maybeEmpty = rootParent.GetLastChild() {
			maybeEmpty.Detach()
		}
		r.compiler.ReportChangeToEnclosingScope(rootParent)
	case rhino.Token_PARAM_LIST:
		if !root.IsDefaultValue() {
			return
		}
		r.compiler.ReportChangeToEnclosingScope(rootParent)
		name := root.GetFirstChild()
		jutil.CheckState(name.IsName())
		if root == rootParent.GetLastChild() &&
			r.removeGlobals &&
			r.CanRemoveParameters(rootParent) {
			root.Detach()
		} else {
			root.ReplaceWith(name.Detach())
		}
		NodeUtil_MarkFunctionsDeleted(root, r.compiler)
	case rhino.Token_OBJECT_PATTERN:
		NodeUtil_DeleteNode(root, r.compiler)
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected parent of indirect assignment: " + rootParent.ToStringTree()))
	}
}

func (b *RemoveUnusedCode_RemovableBase) String() string {
	switch b.remKind {
	case remKindUnusedReadRef:
		return "UnusedReadReference:" + b.targetNode.String()
	case remKindInstanceofName:
		return "InstanceofName:" + b.targetNode.String()
	case remKindWeakUsageCall:
		return "weakUsageCall:" + b.targetNode.String()
	case remKindIncOrDecOp:
		return "IncOrDecOp:" + b.targetNode.String()
	case remKindPolyfill:
		if b.flags&remFlagSubBit != 0 {
			return "Patch:" + b.targetNode.String()
		}
		return "Polyfill:" + b.targetNode.String()
	case remKindClassDecl:
		return "ClassDeclaration:" + b.targetNode.String()
	case remKindNamedClassExpr:
		return "NamedClassExpression:" + b.targetNode.String()
	case remKindClassOrProtoProp:
		return "ClassOrPrototypeNamedProperty:" + b.targetNode.String()
	case remKindFunctionDecl:
		return "FunctionDeclaration:" + b.targetNode.String()
	case remKindNameDeclStmt:
		return "NameDeclStmt:" + b.targetNode.String()
	case remKindAssignVar, remKindAssignNamedProp, remKindAssignComputedProp:
		return "Assign:" + b.targetNode.String()
	case remKindAnonProtoPropAssign:
		return "AnonymousPrototypeNamedPropertyAssign:" + b.targetNode.String()
	case remKindClassSetupCall:
		return "ClassSetupCall:" + b.targetNode.String()
	default:
		return "Removable:" + b.targetNode.String()
	}
}

// RemoveUnusedCode_RemovableBuilder is the inner class RemovableBuilder.
type RemoveUnusedCode_RemovableBuilder struct {
	outer          *RemoveUnusedCode
	cont0          RemoveUnusedCode_Continuation
	propertyNameID uint32
	extraContIdx   uint32
	numConts       uint16
	flags          uint8
}

// Java: RemoveUnusedCode.RemovableBuilder()
func NewRemoveUnusedCode_RemovableBuilder(outer *RemoveUnusedCode) *RemoveUnusedCode_RemovableBuilder {
	return &RemoveUnusedCode_RemovableBuilder{outer: outer}
}

// Java: RemoveUnusedCode.RemovableBuilder.addContinuation(Continuation)
func (b *RemoveUnusedCode_RemovableBuilder) AddContinuation(continuation RemoveUnusedCode_Continuation) *RemoveUnusedCode_RemovableBuilder {
	switch b.numConts {
	case 0:
		b.cont0 = continuation
	case 1:
		b.extraContIdx = uint32(len(b.outer.extraContSlab))
		b.outer.extraContSlab = append(b.outer.extraContSlab, continuation)
	default:
		b.outer.extraContSlab = append(b.outer.extraContSlab, continuation)
	}
	b.numConts++
	return b
}

// Java: RemoveUnusedCode.RemovableBuilder.setIsPrototypeDotPropertyReference(boolean)
func (b *RemoveUnusedCode_RemovableBuilder) SetIsPrototypeDotPropertyReference(value bool) *RemoveUnusedCode_RemovableBuilder {
	if value {
		b.flags |= remFlagProtoDotPropRef
	} else {
		b.flags &^= remFlagProtoDotPropRef
	}
	return b
}

// Java: RemoveUnusedCode.RemovableBuilder.setIsThisDotPropertyReference(boolean)
func (b *RemoveUnusedCode_RemovableBuilder) SetIsThisDotPropertyReference(value bool) *RemoveUnusedCode_RemovableBuilder {
	if value {
		b.flags |= remFlagThisDotPropRef
	} else {
		b.flags &^= remFlagThisDotPropRef
	}
	return b
}

func (b *RemoveUnusedCode_RemovableBuilder) setPropertyName(name string) {
	b.propertyNameID = uint32(rhino.RhinoStringPool_InternID(name))
	b.flags |= remFlagHasPropertyName
}

func (b *RemoveUnusedCode_RemovableBuilder) setPropertyNameNode(n rhino.Node) {
	id := n.GetStringID()
	if id == 0 {
		s := n.GetString()
		if s != "" {
			id = uint32(rhino.RhinoStringPool_InternID(s))
		}
	}
	b.propertyNameID = id
	b.flags |= remFlagHasPropertyName
}

// Java: RemoveUnusedCode.RemovableBuilder.buildIndirectAssign(Node, Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildIndirectAssign(root rhino.Node, targetNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_IndirectAssign(b.outer, b, root, targetNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildPolyfill(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildPolyfill(polyfillNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_Polyfill(b.outer, b, polyfillNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildClassDeclaration(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildClassDeclaration(classNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_ClassDeclaration(b.outer, b, classNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildNamedClassExpression(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildNamedClassExpression(classNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_NamedClassExpression(b.outer, b, classNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildClassOrPrototypeNamedProperty(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildClassOrPrototypeNamedProperty(propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	if !(propertyNode.IsMemberFunctionDef() ||
		propertyNode.IsMemberFieldDef() ||
		NodeUtil_IsGetOrSetKey(propertyNode) ||
		(propertyNode.IsStringKey() && !propertyNode.IsQuotedStringKey())) {
		jutil.CheckArgumentMsg1(false, "%s", propertyNode)
	}
	b.setPropertyNameNode(propertyNode)
	return NewRemoveUnusedCode_ClassOrPrototypeNamedProperty(b.outer, b, propertyNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildObjectDefinePropertiesDefinition(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildObjectDefinePropertiesDefinition(propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	b.setPropertyNameNode(propertyNode)
	return NewRemoveUnusedCode_ObjectDefinePropertiesDefinition(b.outer, b, propertyNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildFunctionDeclaration(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildFunctionDeclaration(functionNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_FunctionDeclaration(b.outer, b, functionNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildNameDeclarationStatement(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildNameDeclarationStatement(declarationStatement rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_NameDeclarationStatement(b.outer, b, declarationStatement)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildNamedPropertyAssign(Node, Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildNamedPropertyAssign(assignNode rhino.Node, propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return b.BuildNamedPropertyAssignWithVarInfo(assignNode, propertyNode, nil)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildNamedPropertyAssign(Node, Node, VarInfo)
func (b *RemoveUnusedCode_RemovableBuilder) BuildNamedPropertyAssignWithVarInfo(assignNode rhino.Node, propertyNode rhino.Node, varInfo *RemoveUnusedCode_RealVarInfo) *RemoveUnusedCode_RemovableBase {
	b.setPropertyNameNode(propertyNode)
	return NewRemoveUnusedCode_Assign(b.outer, b, assignNode, RemoveUnusedCode_Kind_NAMED_PROPERTY, propertyNode, varInfo)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildComputedPropertyAssign(Node, Node, VarInfo)
func (b *RemoveUnusedCode_RemovableBuilder) BuildComputedPropertyAssign(assignNode rhino.Node, propertyNode rhino.Node, varInfo *RemoveUnusedCode_RealVarInfo) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_Assign(b.outer, b, assignNode, RemoveUnusedCode_Kind_COMPUTED_PROPERTY, propertyNode, varInfo)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildVariableAssign(Node, VarInfo)
func (b *RemoveUnusedCode_RemovableBuilder) BuildVariableAssign(assignNode rhino.Node, varInfo *RemoveUnusedCode_RealVarInfo) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_Assign(b.outer, b, assignNode, RemoveUnusedCode_Kind_VARIABLE /* propertyNode= */, rhino.Node{}, varInfo)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildClassSetupCall(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildClassSetupCall(callNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return b.BuildClassSetupCallWithClassDefiningCall(callNode /* classDefiningCall= */, false)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildClassSetupCall(Node, boolean)
func (b *RemoveUnusedCode_RemovableBuilder) BuildClassSetupCallWithClassDefiningCall(callNode rhino.Node, classDefiningCall bool) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_ClassSetupCall(b.outer, b, callNode, classDefiningCall)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildVanillaForNameDeclaration(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildVanillaForNameDeclaration(nameNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_VanillaForNameDeclaration(b.outer, b, nameNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildAnonymousPrototypeNamedPropertyAssign(Node, String)
func (b *RemoveUnusedCode_RemovableBuilder) BuildAnonymousPrototypeNamedPropertyAssign(
	assignNode rhino.Node, propertyName string) *RemoveUnusedCode_RemovableBase {
	b.setPropertyName(propertyName)
	return NewRemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign(b.outer, b, assignNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildIncOrDepOp(Node, Node, Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildIncOrDepOp(incOrDecOp rhino.Node, propertyNode rhino.Node, toPreseve rhino.Node) *RemoveUnusedCode_RemovableBase {
	b.setPropertyNameNode(propertyNode)
	return NewRemoveUnusedCode_IncOrDecOp(b.outer, b, incOrDecOp, toPreseve)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildUnusedReadReference(Node, Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildUnusedReadReference(referenceNode rhino.Node, propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	b.setPropertyNameNode(propertyNode)
	return NewRemoveUnusedCode_UnusedReadReference(b.outer, b, referenceNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildInstanceofName(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildInstanceofName(instanceofNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_InstanceofName(b.outer, b, instanceofNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildWeakUsageCall(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildWeakUsageCall(weakUsageCall rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_WeakUsageCall(b.outer, b, weakUsageCall)
}

// Java: RemoveUnusedCode.UnusedReadReference(RemovableBuilder, Node)
func NewRemoveUnusedCode_UnusedReadReference(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, referenceNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	if !(RemoveUnusedCode_IsThisDotProperty(referenceNode) || RemoveUnusedCode_IsDotPrototypeDotProperty(referenceNode)) {
		jutil.CheckStateMsg1(false, "%s", referenceNode)
	}
	u := outer.removables.New()
	u.initFromBuilder(remKindUnusedReadRef, referenceNode, builder)
	return u
}

// Java: RemoveUnusedCode.InstanceofName(RemovableBuilder, Node)
func NewRemoveUnusedCode_InstanceofName(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, instanceofNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	if !instanceofNode.IsInstanceOf() {
		jutil.CheckArgumentMsg1(false, "%s", instanceofNode)
	}
	i := outer.removables.New()
	i.initFromBuilder(remKindInstanceofName, instanceofNode, builder)
	return i
}

// Java: RemoveUnusedCode.WeakUsageCall(RemovableBuilder, Node)
func NewRemoveUnusedCode_WeakUsageCall(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, weakCallNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	jutil.CheckArgumentMsg1(weakCallNode.IsCall(), "%s", weakCallNode)
	w := outer.removables.New()
	w.initFromBuilder(remKindWeakUsageCall, weakCallNode, builder)
	return w
}

// Java: RemoveUnusedCode.IncOrDecOp(RemovableBuilder, Node, Node)
func NewRemoveUnusedCode_IncOrDecOp(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, incOrDecNode rhino.Node, toPreserve rhino.Node) *RemoveUnusedCode_RemovableBase {
	jutil.CheckArgumentMsg1(incOrDecNode.IsInc() || incOrDecNode.IsDec(), "%s", incOrDecNode)
	arg := incOrDecNode.GetOnlyChild()
	jutil.CheckStateMsg1(RemoveUnusedCode_IsThisDotProperty(arg) || RemoveUnusedCode_IsDotPrototypeDotProperty(arg), "%s", arg)
	o := outer.removables.New()
	o.initFromBuilder(remKindIncOrDecOp, incOrDecNode, builder)
	o.auxNode = toPreserve
	return o
}

// RemoveUnusedCode_IsThisDotProperty: true for `this.propertyName`.
//
// Java: RemoveUnusedCode.isThisDotProperty(Node)
func RemoveUnusedCode_IsThisDotProperty(n rhino.Node) bool {
	return NodeUtil_IsNormalOrOptChainGetProp(n) && n.GetFirstChild().IsThis()
}

// RemoveUnusedCode_IsDotPrototypeDotProperty: true for
// `(something).prototype.propertyName`.
//
// Java: RemoveUnusedCode.isDotPrototypeDotProperty(Node)
func RemoveUnusedCode_IsDotPrototypeDotProperty(n rhino.Node) bool {
	return NodeUtil_IsNormalOrOptChainGetProp(n) && RemoveUnusedCode_IsDotPrototype(n.GetFirstChild())
}

// Java: RemoveUnusedCode.IndirectAssign(RemovableBuilder, Node, Node)
func NewRemoveUnusedCode_IndirectAssign(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, root rhino.Node, targetNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	rootParent := root.GetParent()
	if !(rootParent.IsDestructuringPattern() || rootParent.IsParamList()) {
		jutil.CheckStateMsg1(false, "%s", rootParent)
	}
	if !(targetNode.IsName() || targetNode.IsGetProp()) {
		jutil.CheckStateMsg1(false, "%s", targetNode)
	}
	a := outer.removables.New()
	a.initFromBuilder(remKindIndirectAssign, targetNode, builder)
	a.auxNode = root
	return a
}

// Java: RemoveUnusedCode.Polyfill(RemovableBuilder, Node)
func NewRemoveUnusedCode_Polyfill(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, polyfillNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	p := outer.removables.New()
	p.initFromBuilder(remKindPolyfill, polyfillNode, builder)
	if strings.Contains(polyfillNode.GetFirstFirstChild().GetString(), "patch") {
		p.flags |= remFlagSubBit
	}
	return p
}

// Java: RemoveUnusedCode.ClassDeclaration(RemovableBuilder, Node)
func NewRemoveUnusedCode_ClassDeclaration(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, classDeclarationNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	c := outer.removables.New()
	c.initFromBuilder(remKindClassDecl, classDeclarationNode, builder)
	return c
}

// Java: RemoveUnusedCode.NamedClassExpression(RemovableBuilder, Node)
func NewRemoveUnusedCode_NamedClassExpression(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, classNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	c := outer.removables.New()
	c.initFromBuilder(remKindNamedClassExpr, classNode, builder)
	return c
}

// Java: RemoveUnusedCode.ClassOrPrototypeNamedProperty(RemovableBuilder, Node)
func NewRemoveUnusedCode_ClassOrPrototypeNamedProperty(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	p := outer.removables.New()
	p.initFromBuilder(remKindClassOrProtoProp, propertyNode, builder)
	return p
}

// Java: RemoveUnusedCode.ObjectDefinePropertiesDefinition(RemovableBuilder, Node)
func NewRemoveUnusedCode_ObjectDefinePropertiesDefinition(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	d := outer.removables.New()
	d.initFromBuilder(remKindObjectDefinePropsDef, propertyNode, builder)
	return d
}

// Java: RemoveUnusedCode.FunctionDeclaration(RemovableBuilder, Node)
func NewRemoveUnusedCode_FunctionDeclaration(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, functionDeclarationNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	f := outer.removables.New()
	f.initFromBuilder(remKindFunctionDecl, functionDeclarationNode, builder)
	return f
}

// Java: RemoveUnusedCode.NameDeclarationStatement(RemovableBuilder, Node)
func NewRemoveUnusedCode_NameDeclarationStatement(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, declarationStatement rhino.Node) *RemoveUnusedCode_RemovableBase {
	if !NodeUtil_IsNameDeclaration(declarationStatement) {
		jutil.CheckArgumentMsg1(false, "%s", declarationStatement)
	}
	s := outer.removables.New()
	s.initFromBuilder(remKindNameDeclStmt, declarationStatement, builder)
	if declarationStatement.GetOnlyChild().GetFirstChild().IsNil() {
		s.flags |= remFlagUninitVarDecl
	}
	return s
}

// RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode: targetNode is the
// node to which a value is being assigned; valueNode the value being
// assigned. If `valueNode` has the form `qualifiedName || defaultValue` and
// `qualifiedName` matches `targetNode`, return `defaultValue`. Otherwise
// return `valueNode`.
//
// Java: RemoveUnusedCode.maybeUnwrapQnameOrDefaultValueNode(Node, Node)
func RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(targetNode rhino.Node, valueNode rhino.Node) rhino.Node {
	if valueNode.IsOr() && targetNode.IsQualifiedName() {
		lhsOfOr := rhino.CheckNotNull(valueNode.GetFirstChild())
		if lhsOfOr.MatchesQualifiedNameNode(targetNode) {
			return valueNode.GetLastChild()
		}
	}
	return valueNode
}

// RemoveUnusedCode_Kind is the nested enum Kind.
type RemoveUnusedCode_Kind int32

const (
	// X = something;
	RemoveUnusedCode_Kind_VARIABLE RemoveUnusedCode_Kind = iota
	// X.propertyName = something;
	// X.prototype.propertyName = something;
	RemoveUnusedCode_Kind_NAMED_PROPERTY
	// X[expression] = something;
	// X.prototype[expression] = something;
	RemoveUnusedCode_Kind_COMPUTED_PROPERTY
)

var removeUnusedCode_Kind_names = [...]string{"VARIABLE", "NAMED_PROPERTY", "COMPUTED_PROPERTY"}

// Java: RemoveUnusedCode.Kind.values()
func RemoveUnusedCode_Kind_values() []RemoveUnusedCode_Kind {
	return []RemoveUnusedCode_Kind{
		RemoveUnusedCode_Kind_VARIABLE, RemoveUnusedCode_Kind_NAMED_PROPERTY, RemoveUnusedCode_Kind_COMPUTED_PROPERTY,
	}
}

// Java: RemoveUnusedCode.Kind.name()
func (k RemoveUnusedCode_Kind) String() string { return removeUnusedCode_Kind_names[k] }

// Java: RemoveUnusedCode.Assign(RemovableBuilder, Node, Kind, Node, VarInfo)
func NewRemoveUnusedCode_Assign(
	outer *RemoveUnusedCode,
	builder *RemoveUnusedCode_RemovableBuilder,
	assignNode rhino.Node,
	kind RemoveUnusedCode_Kind,
	propertyNode rhino.Node,
	varInfo *RemoveUnusedCode_RealVarInfo) *RemoveUnusedCode_RemovableBase {
	if !NodeUtil_IsAssignmentOp(assignNode) {
		jutil.CheckArgumentMsg1(false, "%s", assignNode)
	}
	var rk uint8
	if kind == RemoveUnusedCode_Kind_VARIABLE {
		if !propertyNode.IsNil() {
			jutil.CheckArgumentMsg1(
				false,
				"got property node for simple variable assignment: %s",
				propertyNode)
		}
		if varInfo == nil {
			jutil.CheckArgumentMsg1(false, "missing VarInfo for variable assignment: %s", propertyNode)
		}
		rk = remKindAssignVar
	} else {
		if propertyNode.IsNil() {
			jutil.CheckArgumentMsg(false, "missing property node")
		}
		if kind == RemoveUnusedCode_Kind_NAMED_PROPERTY {
			if !propertyNode.IsGetProp() {
				jutil.CheckArgumentMsg1(
					false, "property name is not a GETPROP: %s", propertyNode)
			}
			rk = remKindAssignNamedProp
		} else {
			rk = remKindAssignComputedProp
		}
	}
	a := outer.removables.New()
	a.initFromBuilder(rk, assignNode, builder)
	a.varInfo = varInfo
	return a
}

// Java: RemoveUnusedCode.mayHaveSideEffects(Node)
func (r *RemoveUnusedCode) MayHaveSideEffects(node rhino.Node) bool {
	switch node.GetToken() {
	case rhino.Token_NAME, rhino.Token_THIS, rhino.Token_NUMBER, rhino.Token_STRINGLIT,
		rhino.Token_TRUE, rhino.Token_FALSE, rhino.Token_NULL, rhino.Token_EMPTY,
		rhino.Token_FUNCTION, rhino.Token_BIGINT:
		return false
	case rhino.Token_VOID:
		return r.MayHaveSideEffects(node.GetFirstChild())
	case rhino.Token_CALL:
		if jsDocInfo := node.GetJSDocInfo(); jsDocInfo != nil && jsDocInfo.IsPureOrBreakMyCode() {
			return false
		}
	}
	return r.astAnalyzer.MayHaveSideEffects(node)
}

// Java: RemoveUnusedCode.AnonymousPrototypeNamedPropertyAssign(RemovableBuilder, Node)
func NewRemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, assignNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	if builder.flags&remFlagHasPropertyName == 0 {
		panic(jutil.NewNullPointerException())
	}
	jutil.CheckArgumentMsg1(assignNode.IsAssign(), "%s", assignNode)
	a := outer.removables.New()
	a.initFromBuilder(remKindAnonProtoPropAssign, assignNode, builder)
	return a
}

// Java: RemoveUnusedCode.ClassSetupCall(RemovableBuilder, Node, boolean)
func NewRemoveUnusedCode_ClassSetupCall(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, callNode rhino.Node, classDefiningCall bool) *RemoveUnusedCode_RemovableBase {
	c := outer.removables.New()
	c.initFromBuilder(remKindClassSetupCall, callNode, builder)
	if classDefiningCall {
		c.flags |= remFlagSubBit
	}
	return c
}

// Java: RemoveUnusedCode.alreadyRemoved(Node)
func RemoveUnusedCode_AlreadyRemoved(n rhino.Node) bool {
	for parent := n.GetParent(); !parent.IsNil(); parent = parent.GetParent() {
		if parent.IsRoot() {
			return false
		}
	}
	return true
}

// RemoveUnusedCode_VarInfo tracks whether a variable is removable or not,
// including tracking the Removable objects associated with it.
type RemoveUnusedCode_VarInfo interface {
	AddRemovable(removable *RemoveUnusedCode_RemovableBase)
	IsRemovable() bool
	SetIsExplicitlyNotRemovable(reasonSupplier func() string)
	SetHasNonLocalOrNonLiteralValue()
	HasFunctionOrClassLiteralValue() bool
	RemoveAllRemovables()
}

const (
	rucVarFlagEntirelyRemovable uint8 = 1 << iota
	rucVarFlagHasNonLocalOrNonLit
	rucVarFlagHasFuncOrClassLit
	rucVarFlagRequiresLocalLit
	rucVarFlagCanonical
)

// RemoveUnusedCode_RealVarInfo tracks the removable code and other state
// related to variables we may be able to remove.
type RemoveUnusedCode_RealVarInfo struct {
	outer    *RemoveUnusedCode
	var_     *Var
	remHead  *RemoveUnusedCode_RemovableBase
	remTail  *RemoveUnusedCode_RemovableBase
	hashNext *RemoveUnusedCode_RealVarInfo // see removeUnusedCodeState.varBuckets

	hash  uint32
	flags uint8
}

// Java: RemoveUnusedCode.RealVarInfo(String)
func NewRemoveUnusedCode_RealVarInfo(outer *RemoveUnusedCode, var_ *Var) *RemoveUnusedCode_RealVarInfo {
	v := outer.varInfoSlab.New()
	*v = RemoveUnusedCode_RealVarInfo{outer: outer, var_: var_, flags: rucVarFlagEntirelyRemovable}
	return v
}

// Java: RemoveUnusedCode.RealVarInfo.addRemovable(Removable)
func (v *RemoveUnusedCode_RealVarInfo) AddRemovable(removable *RemoveUnusedCode_RemovableBase) {
	if v.flags&rucVarFlagCanonical != 0 {
		v.outer.ConsiderForIndependentRemoval(removable)
		return
	}
	if removable.IsVariableAssignment() {
		if localValue := removable.GetLocalAssignedValue(); !localValue.IsNil() {
			if localValue.IsFunction() || localValue.IsClass() {
				v.flags |= rucVarFlagHasFuncOrClassLit
			}
		} else if removable.flags&remFlagUninitVarDecl == 0 {
			v.flags |= rucVarFlagHasNonLocalOrNonLit
		}
	} else if removable.IsPrototypeAssignment() && !removable.IsAssignedValueLocal() {
		v.flags |= rucVarFlagHasNonLocalOrNonLit
	}
	if removable.PreventsRemovalOfVariableWithNonLocalValueOrPrototype() {
		v.flags |= rucVarFlagRequiresLocalLit
	}
	if v.flags&(rucVarFlagHasNonLocalOrNonLit|rucVarFlagRequiresLocalLit) == (rucVarFlagHasNonLocalOrNonLit | rucVarFlagRequiresLocalLit) {
		v.SetIsExplicitlyNotRemovable(
			func() string { return "hasNonLocalOrNonLiteralValue && requiresLocalLiteralValueForRemoval" })
	}

	if v.flags&rucVarFlagEntirelyRemovable != 0 {
		removable.propNext = nil
		if v.remHead == nil {
			v.remHead = removable
			v.remTail = removable
		} else {
			v.remTail.propNext = removable
			v.remTail = removable
		}
	} else {
		v.outer.ConsiderForIndependentRemoval(removable)
	}
}

// Java: RemoveUnusedCode.RealVarInfo.isRemovable()
func (v *RemoveUnusedCode_RealVarInfo) IsRemovable() bool {
	return v.flags&rucVarFlagEntirelyRemovable != 0
}

// Java: RemoveUnusedCode.RealVarInfo.setIsExplicitlyNotRemovable(Supplier<String>)
func (v *RemoveUnusedCode_RealVarInfo) SetIsExplicitlyNotRemovable(reasonSupplier func() string) {
	if v.flags&rucVarFlagEntirelyRemovable != 0 {
		v.flags &^= rucVarFlagEntirelyRemovable
		if v.outer.unremovableLog.IsLogging() {
			v.outer.unremovableLog.LogString(v.var_.GetName() + ": " + reasonSupplier())
		}
		cur := v.remHead
		v.remHead = nil
		v.remTail = nil
		for cur != nil {
			next := cur.propNext
			cur.propNext = nil
			v.outer.ConsiderForIndependentRemoval(cur)
			cur = next
		}
	}
}

// Java: RemoveUnusedCode.RealVarInfo.setHasNonLocalOrNonLiteralValue()
func (v *RemoveUnusedCode_RealVarInfo) SetHasNonLocalOrNonLiteralValue() {
	v.flags |= rucVarFlagHasNonLocalOrNonLit
}

// Java: RemoveUnusedCode.RealVarInfo.hasFunctionOrClassLiteralValue()
func (v *RemoveUnusedCode_RealVarInfo) HasFunctionOrClassLiteralValue() bool {
	return v.flags&rucVarFlagHasFuncOrClassLit != 0
}

// Java: RemoveUnusedCode.RealVarInfo.removeAllRemovables()
func (v *RemoveUnusedCode_RealVarInfo) RemoveAllRemovables() {
	jutil.CheckState(v.flags&rucVarFlagEntirelyRemovable != 0)
	cur := v.remHead
	v.remHead = nil
	v.remTail = nil
	for cur != nil {
		next := cur.propNext
		cur.propNext = nil
		cur.Remove(v.outer)
		cur = next
	}
}

// CreatePolyfillInfo makes a new PolyfillInfo, including the correct
// Removable. Parses the name to determine whether this is a global, static,
// or prototype polyfill.
//
// Java: RemoveUnusedCode.createPolyfillInfo(Node, Scope, String)
func (r *RemoveUnusedCode) CreatePolyfillInfo(call rhino.Node, scope *Scope, name string) RemoveUnusedCode_PolyfillInfo {
	jutil.CheckState(call.GetParent().IsExprResult())
	if !r.guardedUsagesComputed {
		r.guardedUsagesComputed = true
		NewPolyfillUsageFinder(r.compiler, r.polyfillsFromTable).
			TraverseOnlyGuarded(r.root, r.StorePolyfill)
	}
	// Make the removable and polyfill info.  Add continuations for all arguments.
	builder := NewRemoveUnusedCode_RemovableBuilder(r)
	for n := call.GetFirstChild().GetNext(); !n.IsNil(); n = n.GetNext() {
		builder.AddContinuation(NewRemoveUnusedCode_Continuation(r, n, scope))
	}
	removable := builder.BuildPolyfill(call.GetParent())
	lastDot := jutil.StringLastIndexOf(name, ".")
	if lastDot < 0 {
		return NewRemoveUnusedCode_GlobalPolyfillInfo(r, removable, name)
	}
	owner := jutil.StringSubstring(name, 0, lastDot)
	prop := jutil.StringSubstringFrom(name, lastDot+1)
	if strings.HasSuffix(owner, RemoveUnusedCode_DOT_PROTOTYPE) {
		owner = jutil.StringSubstring(owner, 0, jutil.StringLength(owner)-jutil.StringLength(RemoveUnusedCode_DOT_PROTOTYPE))
		return NewRemoveUnusedCode_PrototypePropertyPolyfillInfo(r, removable, prop, owner)
	}
	return NewRemoveUnusedCode_StaticPropertyPolyfillInfo(r, removable, prop, owner)
}

const RemoveUnusedCode_DOT_PROTOTYPE = ".prototype"

// RemoveUnusedCode_PolyfillInfo stores information about definitions and
// usages of polyfills.
type RemoveUnusedCode_PolyfillInfo interface {
	// Java: RemoveUnusedCode.PolyfillInfo.considerPossibleReference(Node)
	ConsiderPossibleReference(n rhino.Node)

	// Java: RemoveUnusedCode.PolyfillInfo.getName()
	GetName() string

	// ConsiderPossibleReferenceInternal is the template method to check the
	// node.
	//
	// Java: RemoveUnusedCode.PolyfillInfo.considerPossibleReferenceInternal(Node)
	ConsiderPossibleReferenceInternal(n rhino.Node)

	GetRemoveUnusedCode_PolyfillInfoBase() *RemoveUnusedCode_PolyfillInfoBase
}

// RemoveUnusedCode_PolyfillInfoBase holds PolyfillInfo's fields and base
// implementations.
type RemoveUnusedCode_PolyfillInfoBase struct {
	outer *RemoveUnusedCode
	self  RemoveUnusedCode_PolyfillInfo

	// The Polyfill instance corresponding to the polyfill's definition.
	removable *RemoveUnusedCode_RemovableBase

	// The rightmost component of the polyfill's qualified name (does not
	// contain a dot).
	key string

	// Whether the polyfill is unreferenced and this can be removed safely.
	isRemovable bool
}

// InitRemoveUnusedCode_PolyfillInfoBase is the PolyfillInfo(Polyfill,
// String) constructor; self is the polyfill info being constructed.
//
// Java: RemoveUnusedCode.PolyfillInfo(Polyfill, String)
func (b *RemoveUnusedCode_PolyfillInfoBase) InitRemoveUnusedCode_PolyfillInfoBase(self RemoveUnusedCode_PolyfillInfo, outer *RemoveUnusedCode, removable *RemoveUnusedCode_RemovableBase, key string) {
	b.self = self
	b.outer = outer
	b.isRemovable = true
	b.removable = removable
	b.key = key
}

// GetRemoveUnusedCode_PolyfillInfoBase gives access to the base through the
// interface.
func (b *RemoveUnusedCode_PolyfillInfoBase) GetRemoveUnusedCode_PolyfillInfoBase() *RemoveUnusedCode_PolyfillInfoBase {
	return b
}

// ConsiderPossibleReference accepts a NAME or GETPROP node whose (property)
// string matches key and checks whether the node should be considered as a
// possible reference to this polyfill. If so, mark the polyfill as
// referenced and therefore not removable.
//
// Java: RemoveUnusedCode.PolyfillInfo.considerPossibleReference(Node)
func (b *RemoveUnusedCode_PolyfillInfoBase) ConsiderPossibleReference(n rhino.Node) {
	if b.isRemovable && (b.removable.flags&remFlagSubBit != 0 || b.outer.guardedUsages == nil || !b.outer.guardedUsages.Contains(n)) {
		b.self.ConsiderPossibleReferenceInternal(n)
		if !b.isRemovable {
			b.removable.ApplyContinuations(b.outer)
		}
	}
}

// Java: RemoveUnusedCode.PolyfillInfo.getName()
func (b *RemoveUnusedCode_PolyfillInfoBase) GetName() string {
	return b.key
}

// RemoveUnusedCode_GlobalPolyfillInfo is the inner class GlobalPolyfillInfo.
type RemoveUnusedCode_GlobalPolyfillInfo struct {
	RemoveUnusedCode_PolyfillInfoBase
}

// Java: RemoveUnusedCode.GlobalPolyfillInfo(Polyfill, String)
func NewRemoveUnusedCode_GlobalPolyfillInfo(outer *RemoveUnusedCode, removable *RemoveUnusedCode_RemovableBase, name string) *RemoveUnusedCode_GlobalPolyfillInfo {
	g := &RemoveUnusedCode_GlobalPolyfillInfo{}
	g.InitRemoveUnusedCode_PolyfillInfoBase(g, outer, removable, name)
	return g
}

// Java: RemoveUnusedCode.GlobalPolyfillInfo.considerPossibleReferenceInternal(Node)
func (g *RemoveUnusedCode_GlobalPolyfillInfo) ConsiderPossibleReferenceInternal(possiblyReferencingNode rhino.Node) {
	if possiblyReferencingNode.IsName() {
		g.isRemovable = false
	} else if NodeUtil_IsNormalOrOptChainGetProp(possiblyReferencingNode) {
		g.isRemovable = false
	}
}

// RemoveUnusedCode_StaticPropertyPolyfillInfo is the inner class
// StaticPropertyPolyfillInfo.
type RemoveUnusedCode_StaticPropertyPolyfillInfo struct {
	RemoveUnusedCode_PolyfillInfoBase
	// Name of the owning type, used only for debugging.
	polyfillOwnerName string
}

// Java: RemoveUnusedCode.StaticPropertyPolyfillInfo(Polyfill, String, String)
func NewRemoveUnusedCode_StaticPropertyPolyfillInfo(outer *RemoveUnusedCode, removable *RemoveUnusedCode_RemovableBase, key string, ownerName string) *RemoveUnusedCode_StaticPropertyPolyfillInfo {
	s := &RemoveUnusedCode_StaticPropertyPolyfillInfo{}
	s.InitRemoveUnusedCode_PolyfillInfoBase(s, outer, removable, key)
	s.polyfillOwnerName = ownerName
	return s
}

// Java: RemoveUnusedCode.StaticPropertyPolyfillInfo.getName()
func (s *RemoveUnusedCode_StaticPropertyPolyfillInfo) GetName() string {
	return s.polyfillOwnerName + "." + s.key
}

// Java: RemoveUnusedCode.StaticPropertyPolyfillInfo.considerPossibleReferenceInternal(Node)
func (s *RemoveUnusedCode_StaticPropertyPolyfillInfo) ConsiderPossibleReferenceInternal(possiblyReferencingNode rhino.Node) {
	if NodeUtil_IsNormalOrOptChainGetProp(possiblyReferencingNode) {
		s.isRemovable = false
	}
}

// RemoveUnusedCode_PrototypePropertyPolyfillInfo is the inner class
// PrototypePropertyPolyfillInfo.
type RemoveUnusedCode_PrototypePropertyPolyfillInfo struct {
	RemoveUnusedCode_PolyfillInfoBase
	// Name of the owning type, used only for debugging.
	polyfillOwnerName string
}

// Java: RemoveUnusedCode.PrototypePropertyPolyfillInfo(Polyfill, String, String)
func NewRemoveUnusedCode_PrototypePropertyPolyfillInfo(outer *RemoveUnusedCode, removable *RemoveUnusedCode_RemovableBase, key string, polyfillOwnerName string) *RemoveUnusedCode_PrototypePropertyPolyfillInfo {
	p := &RemoveUnusedCode_PrototypePropertyPolyfillInfo{}
	p.InitRemoveUnusedCode_PolyfillInfoBase(p, outer, removable, key)
	p.polyfillOwnerName = polyfillOwnerName
	return p
}

// Java: RemoveUnusedCode.PrototypePropertyPolyfillInfo.getName()
func (p *RemoveUnusedCode_PrototypePropertyPolyfillInfo) GetName() string {
	return p.polyfillOwnerName + ".prototype." + p.key
}

// Java: RemoveUnusedCode.PrototypePropertyPolyfillInfo.considerPossibleReferenceInternal(Node)
func (p *RemoveUnusedCode_PrototypePropertyPolyfillInfo) ConsiderPossibleReferenceInternal(possiblyReferencingNode rhino.Node) {
	if NodeUtil_IsNormalOrOptChainGetProp(possiblyReferencingNode) {
		p.isRemovable = false
	}
}

// Java: RemoveUnusedCode.VanillaForNameDeclaration(RemovableBuilder, Node)
func NewRemoveUnusedCode_VanillaForNameDeclaration(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, nameNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	d := outer.removables.New()
	d.initFromBuilder(remKindVanillaForNameDecl, nameNode, builder)
	if nameNode.GetFirstChild().IsNil() {
		d.flags |= remFlagUninitVarDecl
	}
	return d
}

// Java: RemoveUnusedCode.removeExpressionCompletely(Node)
func (r *RemoveUnusedCode) RemoveExpressionCompletely(expression rhino.Node) {
	jutil.CheckStateMsg1(!NodeUtil_IsExpressionResultUsed(expression), "%s", expression)
	parent := expression.GetParent()
	if parent.IsExprResult() {
		NodeUtil_DeleteNode(parent, r.compiler)
	} else if parent.IsComma() {
		// Expression is probably the first child of the comma,
		// but it could be the second if the entire comma expression value is unused.
		otherChild := expression.GetNext()
		if otherChild.IsNil() {
			otherChild = expression.GetPrevious()
		}
		r.ReplaceNodeWith(parent, otherChild.Detach())
	} else {
		// value isn't needed, but we need to keep the AST valid.
		r.ReplaceNodeWith(expression, rhino.IR_Number(0).Srcref(expression))
	}
}

// Java: RemoveUnusedCode.replaceNodeWith(Node, Node)
func (r *RemoveUnusedCode) ReplaceNodeWith(n rhino.Node, replacement rhino.Node) {
	r.compiler.ReportChangeToEnclosingScope(n)
	n.ReplaceWith(replacement)
	NodeUtil_MarkFunctionsDeleted(n, r.compiler)
}
