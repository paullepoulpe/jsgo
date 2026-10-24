// Port of com/google/javascript/jscomp/DevirtualizeMethods.java @ closure-compiler v20261015

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

package jscomp

import (
	"sync"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DevirtualizeMethods rewrites prototype and static methods as global, free
// functions that take the receiver as their first argument.
//
// This transformation simplifies the call graph so smart name removal, cross
// chunk code motion and other passes can do more.
//
// To work effectively, this pass depends on
// com.google.javascript.jscomp.disambiguate.DisambiguateProperties running
// first to do a lot of heavy-lifting. It assumes that different methods
// will have unique names which in general isn't true for source JavaScript.
//
// This pass should only be used in production code if property and
// variable renaming are turned on. Resulting code may also benefit from
// `--collapse_anonymous_functions` and `--collapse_variable_declarations`
//
// This pass only rewrites functions that are part of a type's prototype or
// are statics on a ctor/interface function. A host of other preconditions
// must also be met. Functions that access the "arguments" variable
// arguments object are not eligible for this optimization, for example.
//
// For example:
//
//	A.prototype.accumulate = function(value) {
//	  this.total += value; return this.total
//	}
//	var total = a.accumulate(2)
//
// will be rewritten as:
//
//	var accumulate = function(self, value) {
//	  self.total += value; return self.total
//	}
//	var total = accumulate(a, 2)
//
// A similar transformation occurs for:
//
//	/** @constructor *\/
//	function A() { }
//
//	A.accumulate = function(value) {
//	  this.total += value; return this.total
//	}
//	var total = a.accumulate(2)
type DevirtualizeMethods struct {
	compiler     AbstractCompiler
	refMap       *OptimizeCalls_ReferenceMap
	defSites     []rhino.Node
	defFunctions []rhino.Node
	callSites    []rhino.Node
}

// Java: DevirtualizeMethods(AbstractCompiler)
func NewDevirtualizeMethods(compiler AbstractCompiler) *DevirtualizeMethods {
	return &DevirtualizeMethods{compiler: compiler}
}

// Java: DevirtualizeMethods.process(Node, Node, ReferenceMap)
func (d *DevirtualizeMethods) Process(externs rhino.Node, root rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
	jutil.CheckStateMsg(d.refMap == nil, "`process` should only be called once.")

	d.refMap = refMap
	defer func() {
		clear(d.defSites)
		d.defSites = d.defSites[:0]
		clear(d.defFunctions)
		d.defFunctions = d.defFunctions[:0]
		clear(d.callSites)
		d.callSites = d.callSites[:0]
	}()

	for key, value := range refMap.GetPropReferences() {
		d.ProcessReferenceList(key, value.Slice())
	}
}

// Java: DevirtualizeMethods.processReferenceList(String, List<Node>)
func (d *DevirtualizeMethods) ProcessReferenceList(name string, sites []rhino.Node) {
	defSites := d.defSites[:0]
	defFunctions := d.defFunctions[:0]
	callSites := d.callSites[:0]
	for _, s := range sites {
		prevLen := len(defFunctions)
		defFunctions = optimizeCalls_appendDefinitionFunctionNodesFor(defFunctions, s)
		if len(defFunctions) > prevLen {
			defSites = append(defSites, s)
		} else {
			// If a site has no associated functions, it must be a call site.
			callSites = append(callSites, s)
		}
	}
	d.defSites = defSites
	d.defFunctions = defFunctions
	d.callSites = callSites

	if len(defFunctions) == 0 {
		return // We can't devirtualize without a definition.
	}
	if len(callSites) == 0 {
		return // No need to devirtualize without calls.
	}

	// Use the first definition to ensure that all invocations are after the definition
	// (temporally). It's  possible that if class `A` is defined after class `B` there are calls to
	// `A::foo` before `B::foo` is defined.
	canonicalDefinitionSite := defSites[0]

	// Check that this method is safe for devirtualization. These are in (estimated) increasing
	// order of cost.
	if !d.allSitesAreEligibleDefinitionSites(name, defSites) {
		return
	} else if !d.allFunctionsAreEligibleDefinitionFunctions(defFunctions) {
		return
	} else if !d.allCallSitesAreEligible(callSites, canonicalDefinitionSite) {
		return
	} else if !d.AllDefinitionsEquivalent(defFunctions) {
		// Remember that this is only valid because we already checked the scoping of the definitions.
		return
	}

	devirtualizedName := DevirtualizeMethods_RewrittenMethodNameOf(name)
	bestJSDoc := NodeUtil_GetBestJSDocInfo(canonicalDefinitionSite)
	isConstantName := bestJSDoc != nil && bestJSDoc.IsConstant()

	for _, callSite := range callSites {
		d.RewriteCall(callSite, devirtualizedName, isConstantName)
	}
	// We only have to rewrite one definition. We've checked they're all identical so any of them
	// can replace the others. The un-rewritten ones will be dead-code eliminated.
	d.RewriteDefinition(canonicalDefinitionSite, devirtualizedName, isConstantName)
}

// allSitesAreEligibleDefinitionSites is
// `functionsBySite.keySet().stream().allMatch((s) -> isEligibleDefinitionSite(name, s))`.
//
// Go-only: language-adapter — a Java stream.
func (d *DevirtualizeMethods) allSitesAreEligibleDefinitionSites(name string, defSites []rhino.Node) bool {
	for _, s := range defSites {
		if !d.IsEligibleDefinitionSite(name, s) {
			return false
		}
	}
	return true
}

// allFunctionsAreEligibleDefinitionFunctions is
// `functionsBySite.values().stream().allMatch(this::isEligibleDefinitionFunction)`.
//
// Go-only: language-adapter — a Java stream.
func (d *DevirtualizeMethods) allFunctionsAreEligibleDefinitionFunctions(defFunctions []rhino.Node) bool {
	for _, f := range defFunctions {
		if !d.IsEligibleDefinitionFunction(f) {
			return false
		}
	}
	return true
}

// allCallSitesAreEligible is
// `callSites.stream().allMatch((c) -> isEligibleCallSite(c, canonicalDefinitionSite))`.
//
// Go-only: language-adapter — a Java stream.
func (d *DevirtualizeMethods) allCallSitesAreEligible(callSites []rhino.Node, canonicalDefinitionSite rhino.Node) bool {
	for _, c := range callSites {
		if !d.IsEligibleCallSite(c, canonicalDefinitionSite) {
			return false
		}
	}
	return true
}

// Java: DevirtualizeMethods.isPrototypeOrStaticMethodDefinition(Node)
func (d *DevirtualizeMethods) IsPrototypeOrStaticMethodDefinition(node rhino.Node) bool {
	parent := node.GetParent()
	grandparent := node.GetGrandparent()
	if parent.IsNil() || grandparent.IsNil() {
		return false
	}

	switch node.GetToken() {
	case rhino.Token_MEMBER_FUNCTION_DEF:
		if NodeUtil_IsEs6ConstructorMemberFunctionDef(node) {
			return false // Constructors aren't methods.
		}

		if parent.IsObjectLit() {
			// Case: `Foo.prototype = {
			//          bar() { },
			//        }`
			return DevirtualizeMethods_IsPrototypeObjectLit(parent)
		}

		return true // Case: `class Foo { bar() { } }`
	case rhino.Token_GETPROP:
		// Case: `Foo.prototype.bar = function() { };
		if !node.IsFirstChildOf(parent) ||
			!NodeUtil_IsExprAssign(grandparent) ||
			!parent.GetLastChild().IsFunction() {
			return false
		}

		if NodeUtil_IsPrototypeProperty(node) {
			return true
		}

		if d.IsDefinitelyCtorOrInterface(node.GetFirstChild()) {
			return true
		}

		return false
	case rhino.Token_STRING_KEY:
		// Case: `Foo.prototype = {
		//          bar: function() { },
		//        }`
		jutil.CheckArgumentMsg1(parent.IsObjectLit(), "%s", parent)

		return DevirtualizeMethods_IsPrototypeObjectLit(parent) && node.GetFirstChild().IsFunction()
	default:
		return false
	}
}

// DevirtualizeMethods_IsPrototypeObjectLit returns whether objectLit is the
// value assigned in `Foo.prototype = {...}`.
//
// Methods on other object literals are not eligible: plain objects are
// routinely enumerated, copied, or mixed in, so removing a method from one is
// unsafe without the escape analysis that InlineAndCollapseProperties
// performs.
//
// Java: DevirtualizeMethods.isPrototypeObjectLit(Node)
func DevirtualizeMethods_IsPrototypeObjectLit(objectLit rhino.Node) bool {
	jutil.CheckArgumentMsg1(objectLit.IsObjectLit(), "%s", objectLit)
	parent := objectLit.GetParent()
	return objectLit.IsSecondChildOf(parent) &&
		NodeUtil_IsPrototypeAssignment(parent.GetFirstChild())
}

// Java: DevirtualizeMethods.isDefinitelyCtorOrInterface(Node)
func (d *DevirtualizeMethods) IsDefinitelyCtorOrInterface(receiver rhino.Node) bool {
	qname := receiver.GetQualifiedName()
	if qname == nil {
		return false
	}

	// It's safe to rely on the global scope because normalization has made all names unique. No
	// local will shadow a global name, and we're ok with missing some statics.
	var_ := d.refMap.GetGlobalScope().GetVar(*qname)
	if var_ == nil {
		return false
	}

	if var_.IsClass() {
		return true
	}

	jsdoc := var_.GetJSDocInfo()
	if jsdoc == nil {
		return false
	} else if jsdoc.IsConstructorOrInterface() {
		return true // Case: `@constructor`, `@interface`, `@record`.
	}

	return false
}

// IsEligibleDefinitionSite determines if a method definition site is
// eligible for rewrite as a global.
//
// In order to be eligible for rewrite, the definition site must:
//
//   - Not be exported
//   - Be for a prototype method
//
// Java: DevirtualizeMethods.isEligibleDefinitionSite(String, Node)
func (d *DevirtualizeMethods) IsEligibleDefinitionSite(name string, definitionSite rhino.Node) bool {
	switch definitionSite.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_STRING_KEY, rhino.Token_MEMBER_FIELD_DEF:
	default:
		// No other node types are supported.
		panic(jutil.NewIllegalArgumentExceptionMsg(definitionSite.String()))
	}

	// Exporting a method prevents rewrite.
	codingConvention := d.compiler.GetCodingConvention()
	if codingConvention.IsExportedWithLocal(name /* local= */, false) {
		return false
	}

	if !d.IsPrototypeOrStaticMethodDefinition(definitionSite) {
		return false
	}

	return true
}

// IsEligibleDefinitionFunction determines if a method definition function
// is eligible for rewrite as a global function.
//
// In order to be eligible for rewrite, the definition function must:
//
//   - Be instantiated exactly once
//   - Be the only possible implementation at a given site
//   - Not refer to its `arguments`; no implicit varags
//   - Not be an arrow function
//
// Java: DevirtualizeMethods.isEligibleDefinitionFunction(Node)
func (d *DevirtualizeMethods) IsEligibleDefinitionFunction(definitionFunction rhino.Node) bool {
	jutil.CheckArgumentMsg1(definitionFunction.IsFunction(), "%s", definitionFunction)

	if definitionFunction.IsArrowFunction() {
		return false
	}

	for ancestor := definitionFunction.GetParent(); !ancestor.IsNil(); ancestor = ancestor.GetParent() {
		// The definition must be made exactly once. (i.e. not in a loop, conditional, or function)
		if DevirtualizeMethods_IsScopingOrBranchingConstruct(ancestor) {
			return false
		}

		// TODO(nickreid): Support this so long as the definition doesn't reference the name.
		// We can't allow this in general because references to the local name:
		//  - won't be rewritten correctly
		//  - won't be the same across multiple definitions, even if they are node-wise identical.
		if ancestor.IsClass() && DevirtualizeMethods_LocalNameIsDeclaredByClass(ancestor) {
			return false
		}
	}

	if NodeUtil_Has(definitionFunction, (rhino.Node).IsSuper, func(rhino.Node) bool { return true }) {
		// TODO(b/120452418): Remove this when we have a rewrite for `super`. We punted initially due
		// to complexity.
		return false
	}

	if NodeUtil_DoesFunctionReferenceOwnArgumentsObject(definitionFunction) {
		// Functions that access "arguments" are not eligible since rewriting changes the structure of
		// the function params.
		return false
	}

	return true
}

// IsEligibleCallSite determines if a method call is eligible for rewrite as
// a global function.
//
// In order to be eligible for rewrite, the call must:
//
//   - Property is never accessed outside a function call context.
//
// Java: DevirtualizeMethods.isEligibleCallSite(Node, Node)
func (d *DevirtualizeMethods) IsEligibleCallSite(access rhino.Node, definitionSite rhino.Node) bool {
	invocation := access.GetParent()
	if !NodeUtil_IsInvocationTarget(access) || !invocation.IsCall() {
		// TODO(nickreid): Use the same definition of "a call" as
		// `OptimizeCalls::ReferenceMap::isCallTarget`.
		//
		// Accessing the property in any way besides CALL has issues:
		//   - tear-off: invocations can't be tracked
		//   - as constructor: unsupported rewrite
		//   - as tagged template string: unsupported rewrite
		//   - as optional chaining call: original invocation target may be null/undefined and
		// rewriting creates an unconditional call to a well-defined global function
		return false
	}

	// We can't rewrite functions called in chunks that do not depend on the defining chunk.
	// This is due to a subtle execution order change introduced by rewriting. Example:
	//
	//     `x.foo().bar()` => `JSCompiler_StaticMethods_bar(x.foo())`
	//
	// Note how `JSCompiler_StaticMethods_bar` will be resolved before `x.foo()` is executed. In
	// the case that `x.foo()` defines `JSCompiler_StaticMethods_bar` (e.g. by dynamically loading
	// the defining chunk) this change in ordering will cause a `ReferenceError`. No error would
	// be thrown by the original code because `bar` would be resolved later.
	//
	// We choose to use chunk ordering to avoid this issue because:
	//   - The other eligibility checks for devirtualization prevent any other dangerous cases
	//     that JSCompiler supports.
	//   - Rewriting all call-sites in a way that preserves exact ordering (e.g. using
	//     `ExpressionDecomposer`) has a significant code-size impact (circa 2018-11-19).

	chunkGraph := d.compiler.GetChunkGraph()
	definitionChunk := d.ChunkForNode(definitionSite)
	callChunk := d.ChunkForNode(access)
	if definitionChunk == callChunk {
		// Do nothing.
	} else if callChunk == nil {
		return false
	} else if !chunkGraph.DependsOn(callChunk, definitionChunk) {
		return false
	}

	return true
}

// AllDefinitionsEquivalent, given a set of method definitions, verifies
// they are the same.
//
// Java: DevirtualizeMethods.allDefinitionsEquivalent(Collection<Node>)
func (d *DevirtualizeMethods) AllDefinitionsEquivalent(definitions []rhino.Node) bool {
	if len(definitions) == 0 {
		return true
	}

	definition := definitions[0]
	jutil.CheckArgumentMsg1(definition.IsFunction(), "%s", definition)

	for _, def := range definitions {
		if !d.compiler.AreNodesEqualForInlining(def, definition) {
			return false
		}
	}
	return true
}

// RewriteCall rewrites object method call sites as calls to global
// functions that take "this" as their first argument.
//
// Before: o.foo(a, b, c)
//
// After: foo(o, a, b, c)
//
// Java: DevirtualizeMethods.rewriteCall(Node, String, boolean)
func (d *DevirtualizeMethods) RewriteCall(getprop rhino.Node, newMethodName string, isConstantName bool) {
	jutil.CheckArgumentMsg1(getprop.IsGetProp(), "%s", getprop)
	call := getprop.GetParent()
	jutil.CheckArgumentMsg1(call.IsCall(), "%s", call)
	receiver := getprop.GetFirstChild()

	// This rewriting does not exactly preserve order of operations; the newly inserted static
	// method name will be resolved before `receiver` is evaluated. This is known to be safe due
	// to the eligibility checks earlier in the pass.
	//
	// We choose not to do a full-fidelity rewriting (e.g. using `ExpressionDecomposer`) because
	// doing so means extracting `receiver` into a new variable at each call-site. This  has a
	// significant code-size impact (circa 2018-11-19).

	if !getprop.MutateCallTargetGetPropToNameInPlace(call, receiver, newMethodName, isConstantName) {
		receiver.Detach()
		getprop.ReplaceWith(receiver)
		newReceiver := rhino.IR_Name(newMethodName).CopyTypeFrom(getprop).Srcref(getprop)
		if isConstantName {
			newReceiver.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}
		call.AddChildToFront(newReceiver)
	}

	if receiver.IsSuper() {
		// Case: `super.foo(a, b)` => `foo(this, a, b)`
		receiver.SetToken(rhino.Token_THIS)
	}

	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	d.compiler.ReportChangeToEnclosingScope(call)
}

// RewriteDefinition rewrites method definitions as global functions that
// take "this" as their first argument.
//
// Before: a.prototype.b = function(a, b, c) {...}
//
// After: var b = function(self, a, b, c) {...}
//
// Java: DevirtualizeMethods.rewriteDefinition(Node, String, boolean)
func (d *DevirtualizeMethods) RewriteDefinition(
	definitionSite rhino.Node, newMethodName string, isConstantName bool) {
	if definitionSite.GetToken() == rhino.Token_GETPROP {
		assign := definitionSite.GetParent()
		statement := NodeUtil_GetEnclosingStatement(definitionSite)
		recv := definitionSite.GetFirstChild()
		if !assign.IsNil() && statement == assign.GetParent() && statement.IsExprResult() &&
			!recv.IsNil() && recv.IsQualifiedName() &&
			definitionSite.CanResetInPlace() && statement.CanResetInPlace() && recv.CanResetInPlace() {
			originalJSDoc := NodeUtil_GetBestJSDocInfo(definitionSite)
			function := assign.GetLastChild()

			statement.DetachChildren()
			assign.DetachChildren()
			definitionSite.DetachChildren()
			recv.DetachChildren()

			methodNameID := rhino.RhinoStringPool_InternID(newMethodName)
			selfName := allocPooledStringConcat2(newMethodName, "$self")
			selfNameID := rhino.RhinoStringPool_InternID(selfName)

			definitionSite.ResetDetachedGetPropToNameID(methodNameID, isConstantName)
			definitionSite.AddChildToFront(function)
			statement.MutateAttachedStmtToVarFromSrcref(definitionSite, definitionSite, originalJSDoc)

			paramList := function.GetSecondChild()
			paramList.AddChildToFront(recv.ResetDetachedToFreshNameID(selfNameID).SrcrefIfMissing(function))
			d.compiler.ReportChangeToEnclosingScope(paramList)

			d.replaceReferencesToThisID(function.GetSecondChild(), selfName, selfNameID)
			d.replaceReferencesToThisID(function.GetLastChild(), selfName, selfNameID)

			d.compiler.ReportChangeToEnclosingScope(statement.GetParent())
			d.compiler.ReportChangeToEnclosingScope(statement)
			return
		}
	}

	var function rhino.Node
	var subtreeToRemove rhino.Node
	var nameSource rhino.Node

	switch definitionSite.GetToken() {
	case rhino.Token_GETPROP:
		function = definitionSite.GetParent().GetLastChild()
		nameSource = definitionSite
		subtreeToRemove = NodeUtil_GetEnclosingStatement(definitionSite)
	case rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF:
		function = definitionSite.GetLastChild()
		nameSource = definitionSite
		subtreeToRemove = definitionSite
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg(definitionSite.String()))
	}

	// Define a new variable after the original declaration.
	statement := NodeUtil_GetEnclosingStatement(definitionSite)
	newNameNode := rhino.IR_Name(newMethodName).SrcrefIfMissing(nameSource)
	newVarNode := rhino.IR_Var(newNameNode).SrcrefIfMissing(nameSource)
	newVarNode.InsertBefore(statement)
	if isConstantName {
		newNameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}

	// Copy the JSDocInfo, if any, from the original declaration
	originalJSDoc := NodeUtil_GetBestJSDocInfo(definitionSite)
	newVarNode.SetJSDocInfo(originalJSDoc)

	// Attach the function to the new variable.
	function.Detach()
	newNameNode.AddChildToFront(function)

	// Create the `this` param.
	selfName := allocPooledStringConcat2(newMethodName, "$self")
	selfNameID := rhino.RhinoStringPool_InternID(selfName)
	paramList := function.GetSecondChild()
	paramList.AddChildToFront(rhino.IR_Name(selfName).SrcrefIfMissing(function))
	d.compiler.ReportChangeToEnclosingScope(paramList)

	// Eliminate `this`.
	d.replaceReferencesToThisID(function.GetSecondChild(), selfName, selfNameID) // In default param values.
	d.replaceReferencesToThisID(function.GetLastChild(), selfName, selfNameID)   // In function body.

	// Clean up dangling AST.
	NodeUtil_DeleteNode(subtreeToRemove, d.compiler)
	d.compiler.ReportChangeToEnclosingScope(newVarNode)
}

// ReplaceReferencesToThis replaces references to "this" with references to
// name. Do not traverse function boundaries.
//
// Java: DevirtualizeMethods.replaceReferencesToThis(Node, String)
func (d *DevirtualizeMethods) ReplaceReferencesToThis(node rhino.Node, name string) {
	d.replaceReferencesToThisID(node, name, rhino.RhinoStringPool_InternID(name))
}

func (d *DevirtualizeMethods) replaceReferencesToThisID(node rhino.Node, name string, nameID uint32) {
	if (node.IsFunction() && !node.IsArrowFunction()) ||
		node.IsMemberFieldDef() ||
		NodeUtil_IsClassStaticBlock(node) {
		// Functions (besides arrows), member field defs, and static blocks create a new binding for
		// `this`.
		return
	}

	if node.IsComputedFieldDef() {
		// Computed field keys are evaluated in the enclosing scope, but field initializers create a
		// new binding for `this`.
		key := node.GetFirstChild()
		if !key.IsNil() {
			if key.IsThis() {
				if key.MutateThisToNameIDInPlace(nameID) {
					d.compiler.ReportChangeToEnclosingScope(key)
				} else {
					newName := rhino.IR_Name(name).Srcref(key).CopyTypeFrom(key)
					key.ReplaceWith(newName)
					d.compiler.ReportChangeToEnclosingScope(newName)
				}
			} else {
				d.replaceReferencesToThisID(key, name, nameID)
			}
		}
		return
	}

	for child := node.GetFirstChild(); !child.IsNil(); {
		next := child.GetNext()
		if child.IsThis() {
			if child.MutateThisToNameIDInPlace(nameID) {
				d.compiler.ReportChangeToEnclosingScope(child)
			} else {
				newName := rhino.IR_Name(name).Srcref(child).CopyTypeFrom(child)
				child.ReplaceWith(newName)
				d.compiler.ReportChangeToEnclosingScope(newName)
			}
		} else {
			d.replaceReferencesToThisID(child, name, nameID)
		}
		child = next
	}
}

// Java: DevirtualizeMethods.chunkForNode(Node)
func (d *DevirtualizeMethods) ChunkForNode(node rhino.Node) *JSChunk {
	script := NodeUtil_GetEnclosingScript(node)
	input := d.compiler.GetInput(script.GetInputId())
	return input.GetChunk()
}

type devirtStrSlab struct {
	buf []byte
}

var devirtStrSlabPool = sync.Pool{
	New: func() any { return &devirtStrSlab{} },
}

func allocPooledStringConcat2(prefix, suffix string) string {
	n := len(prefix) + len(suffix)
	if n == 0 {
		return ""
	}
	if n > 2048 {
		return prefix + suffix
	}
	s := devirtStrSlabPool.Get().(*devirtStrSlab)
	if len(s.buf) < n {
		s.buf = make([]byte, 32760) // <= maxSmallSize - MallocHeaderSize noscan
	}
	b := s.buf[:n:n]
	s.buf = s.buf[n:]
	devirtStrSlabPool.Put(s)
	copy(b, prefix)
	copy(b[len(prefix):], suffix)
	return unsafe.String(&b[0], n)
}

// Java: DevirtualizeMethods.rewrittenMethodNameOf(String)
func DevirtualizeMethods_RewrittenMethodNameOf(originalMethodName string) string {
	return allocPooledStringConcat2("JSCompiler_StaticMethods_", originalMethodName)
}

// DevirtualizeMethods_IsScopingOrBranchingConstruct returns true iff a node
// may change the variable bindings of its subtree or cause that subtree to
// be executed not exactly once.
//
// This method does not include CLASS because CLASS does not always create a
// new binding and it is important for the success of this optimization to
// consider class methods.
//
// See localNameIsDeclaredByClass().
//
// Java: DevirtualizeMethods.isScopingOrBranchingConstruct(Node)
func DevirtualizeMethods_IsScopingOrBranchingConstruct(node rhino.Node) bool {
	return NodeUtil_IsControlStructure(node) || // Branching.
		node.IsAnd() || // Branching.
		node.IsOr() || // Branching.
		node.IsFunction() || // Branching & scoping.
		node.IsBlock() // Scoping.
}

// DevirtualizeMethods_LocalNameIsDeclaredByClass returns true iff a CLASS
// subtree declares a name local to the class body.
//
// Example:
//
//	const Foo = class Bar {
//	  qux() { return Bar; }
//	}
//
// Java: DevirtualizeMethods.localNameIsDeclaredByClass(Node)
func DevirtualizeMethods_LocalNameIsDeclaredByClass(clazz rhino.Node) bool {
	jutil.CheckArgumentMsg1(clazz.IsClass(), "%s", clazz)

	if clazz.GetFirstChild().IsEmpty() {
		return false // There must be a name.
	} else if NodeUtil_IsStatement(clazz) {
		return false // The name must be local.
	}

	return true
}
