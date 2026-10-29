// Port of com/google/javascript/jscomp/DevirtualizeMethods.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
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
func (d *DevirtualizeMethods) ProcessWithReferenceMap(externs rhino.Node, root rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
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

	for id, value := range refMap.GetPropReferencesByID() {
		d.processReferenceListByID(id, value)
	}
}

// Java: DevirtualizeMethods.processReferenceList(String, List<Node>)
func (d *DevirtualizeMethods) processReferenceListByID(nameID uint32, sites []rhino.Node) {
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

	if len(defFunctions) == 0 || len(callSites) == 0 {
		return
	}

	// Use the first definition to ensure that all invocations are after the definition
	// (temporally). It's possible that if class `A` is defined after class `B` there are calls to
	// `A::foo` before `B::foo` is defined.
	canonicalDefinitionSite := defSites[0]

	// Check that this method is safe for devirtualization in increasing order of cost.
	name := rhino.RhinoStringPool_StringFromID(nameID)
	isExported := d.compiler.GetCodingConvention().IsExportedWithLocal(name /* local= */, false)
	for _, s := range defSites {
		if !d.isEligibleDefinitionSite(s, isExported) {
			return
		}
	}
	for _, f := range defFunctions {
		if !d.isEligibleDefinitionFunction(ast.MustFunction(f)) {
			return
		}
	}
	var chunkGraph *JSChunkGraph
	var defChunk *JSChunk
	haveDefChunk := false
	for _, c := range callSites {
		call, ok := ast.AsCall(c.GetParent())
		if !ok || !call.Callee().Is(c) {
			// Accessing the property in any way besides CALL has issues:
			//   - tear-off: invocations can't be tracked
			//   - as constructor: unsupported rewrite
			//   - as tagged template string: unsupported rewrite
			//   - as optional chaining call: original invocation target may be null/undefined and
			//     rewriting creates an unconditional call to a well-defined global function
			return
		}
		if !haveDefChunk {
			chunkGraph = d.compiler.GetChunkGraph()
			defChunk = d.chunkForNode(canonicalDefinitionSite)
			haveDefChunk = true
		}
		if !d.isEligibleCallSiteChunk(c, chunkGraph, defChunk) {
			return
		}
	}
	if !d.allDefinitionsEquivalent(defFunctions) {
		// Remember that this is only valid because we already checked the scoping of the definitions.
		return
	}

	devirtualizedName := "JSCompiler_StaticMethods_" + name
	bestJSDoc := NodeUtil_GetBestJSDocInfo(canonicalDefinitionSite)
	isConstantName := bestJSDoc != nil && bestJSDoc.IsConstant()

	for _, callSite := range callSites {
		d.rewriteCall(callSite, devirtualizedName, isConstantName)
	}
	// We only have to rewrite one definition. We've checked they're all identical so any of them
	// can replace the others. The un-rewritten ones will be dead-code eliminated.
	d.rewriteDefinition(canonicalDefinitionSite, devirtualizedName, isConstantName)
}

// Java: DevirtualizeMethods.isPrototypeOrStaticMethodDefinition(Node)
func (d *DevirtualizeMethods) isPrototypeOrStaticMethodDefinition(node rhino.Node) bool {
	parent := node.GetParent()
	grandparent := node.GetGrandparent()
	if parent.IsNil() || grandparent.IsNil() {
		return false
	}

	switch node.GetToken() {
	case rhino.Token_MEMBER_FUNCTION_DEF:
		if ast.MustMemberFunctionDef(node).IsEs6Constructor() {
			return false // Constructors aren't methods.
		}
		if objLit, ok := ast.AsObjectLit(parent); ok {
			// Case: `Foo.prototype = { bar() { } }`
			return devirtualizeMethods_isPrototypeObjectLit(objLit)
		}
		return true // Case: `class Foo { bar() { } }`

	case rhino.Token_GETPROP:
		// Case: `Foo.prototype.bar = function() { };`
		assign, ok := ast.AsAssign(parent)
		if !ok || !assign.Target().Is(node) ||
			!NodeUtil_IsExprAssign(grandparent) ||
			!assign.Value().Node().IsFunction() {
			return false
		}
		if NodeUtil_IsPrototypeProperty(node) {
			return true
		}
		return d.isDefinitelyCtorOrInterface(ast.MustGetProp(node).Object().Node())

	case rhino.Token_STRING_KEY:
		// Case: `Foo.prototype = { bar: function() { } }`
		objLit := ast.MustObjectLit(parent)
		return devirtualizeMethods_isPrototypeObjectLit(objLit) &&
			ast.MustStringKey(node).Value().Node().IsFunction()

	default:
		return false
	}
}

// devirtualizeMethods_isPrototypeObjectLit reports whether objectLit is the
// value assigned in `Foo.prototype = {...}`.
//
// Java: DevirtualizeMethods.isPrototypeObjectLit(Node)
func devirtualizeMethods_isPrototypeObjectLit(objectLit ast.ObjectLit) bool {
	assign, ok := ast.AsAssign(objectLit.Node().GetParent())
	return ok &&
		assign.Value().Is(objectLit.Node()) &&
		NodeUtil_IsPrototypeAssignment(assign.Target().Node())
}

// Java: DevirtualizeMethods.isDefinitelyCtorOrInterface(Node)
func (d *DevirtualizeMethods) isDefinitelyCtorOrInterface(receiver rhino.Node) bool {
	var var_ *Var
	if receiver.IsName() {
		id := receiver.GetStringID()
		if id <= rhino.RhinoStringPool_EmptyID {
			return false
		}
		var_ = d.refMap.GetGlobalScope().GetVarByID(id)
	} else {
		qname, ok := receiver.GetQualifiedNameString()
		if !ok {
			return false
		}
		// It's safe to rely on the global scope because normalization has made all names unique. No
		// local will shadow a global name, and we're ok with missing some statics.
		if id := rhino.RhinoStringPool_LookupID(qname); id != 0 {
			var_ = d.refMap.GetGlobalScope().GetVarByID(id)
		}
	}
	if var_ == nil {
		return false
	}
	if var_.IsClass() {
		return true
	}
	jsdoc := var_.GetJSDocInfo()
	return jsdoc != nil && jsdoc.IsConstructorOrInterface()
}

// Java: DevirtualizeMethods.isEligibleDefinitionSite(String, Node)
func (d *DevirtualizeMethods) isEligibleDefinitionSite(definitionSite rhino.Node, isExported bool) bool {
	switch definitionSite.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_STRING_KEY, rhino.Token_MEMBER_FIELD_DEF:
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg(definitionSite.String()))
	}

	if isExported {
		return false
	}
	return d.isPrototypeOrStaticMethodDefinition(definitionSite)
}

// Java: DevirtualizeMethods.isEligibleDefinitionFunction(Node)
func (d *DevirtualizeMethods) isEligibleDefinitionFunction(definitionFunction ast.Function) bool {
	if definitionFunction.IsArrow() {
		return false
	}

	fnNode := definitionFunction.Node()
	for ancestor := fnNode.GetParent(); !ancestor.IsNil(); ancestor = ancestor.GetParent() {
		// The definition must be made exactly once (not in a loop, conditional, or function).
		if devirtualizeMethods_isScopingOrBranchingConstruct(ancestor) {
			return false
		}
		if cls, ok := ast.AsClass(ancestor); ok && devirtualizeMethods_localNameIsDeclaredByClass(cls) {
			return false
		}
	}

	if NodeUtil_Has(fnNode, (rhino.Node).IsSuper, func(rhino.Node) bool { return true }) {
		// TODO(b/120452418): Remove this when we have a rewrite for `super`.
		return false
	}

	if NodeUtil_DoesFunctionReferenceOwnArgumentsObject(fnNode) {
		return false
	}

	return true
}

func (d *DevirtualizeMethods) isEligibleCallSiteChunk(access rhino.Node, chunkGraph *JSChunkGraph, definitionChunk *JSChunk) bool {
	callChunk := d.chunkForNode(access)
	if definitionChunk == callChunk {
		return true
	}
	if callChunk == nil {
		return false
	}
	return chunkGraph.DependsOn(callChunk, definitionChunk)
}

// Java: DevirtualizeMethods.allDefinitionsEquivalent(Collection<Node>)
func (d *DevirtualizeMethods) allDefinitionsEquivalent(definitions []rhino.Node) bool {
	if len(definitions) == 0 {
		return true
	}
	definition := definitions[0]
	jutil.CheckArgumentMsg1(definition.IsFunction(), "%s", definition)

	for _, def := range definitions[1:] {
		if !d.compiler.AreNodesEqualForInlining(def, definition) {
			return false
		}
	}
	return true
}

// rewriteCall rewrites `o.foo(a, b, c)` into `foo(o, a, b, c)`.
//
// Java: DevirtualizeMethods.rewriteCall(Node, String, boolean)
func (d *DevirtualizeMethods) rewriteCall(getprop rhino.Node, newMethodName string, isConstantName bool) {
	gp := ast.MustGetProp(getprop)
	call := ast.MustCall(getprop.GetParent())
	callNode := call.Node()
	receiver := gp.Object().Node()

	// Note: MutateCallTargetGetPropToNameInPlace mutates getprop's token in place;
	// do not use gp after this call.
	if !getprop.MutateCallTargetGetPropToNameInPlace(callNode, receiver, newMethodName, isConstantName) {
		receiver.Detach()
		getprop.ReplaceWith(receiver)
		newCallee := ast.NewName(newMethodName)
		newCallee.Node().CopyTypeFrom(getprop).Srcref(getprop)
		if isConstantName {
			newCallee.SetConstantName(true)
		}
		callNode.AddChildToFront(newCallee.Node())
	}

	if receiver.IsSuper() {
		// Case: `super.foo(a, b)` => `foo(this, a, b)`
		receiver.SetToken(rhino.Token_THIS)
	}

	call.SetFreeCall(true)
	d.compiler.ReportChangeToEnclosingScope(callNode)
}

// rewriteDefinition rewrites `a.prototype.b = function(a, b) {...}` into
// `var b = function(self, a, b) {...}`.
//
// Java: DevirtualizeMethods.rewriteDefinition(Node, String, boolean)
func (d *DevirtualizeMethods) rewriteDefinition(
	definitionSite rhino.Node, newMethodName string, isConstantName bool) {
	if gp, ok := ast.AsGetProp(definitionSite); ok {
		assignNode := definitionSite.GetParent()
		statement := NodeUtil_GetEnclosingStatement(definitionSite)
		recv := gp.Object().Node()
		if assign, isAssign := ast.AsAssign(assignNode); isAssign &&
			statement == assignNode.GetParent() && statement.IsExprResult() &&
			!recv.IsNil() && recv.IsQualifiedName() &&
			definitionSite.CanResetInPlace() && statement.CanResetInPlace() && recv.CanResetInPlace() {
			originalJSDoc := NodeUtil_GetBestJSDocInfo(definitionSite)
			fn := ast.MustFunction(assign.Value().Node())

			// Note: in-place arena token mutations follow; use raw rhino.Node handles.
			statement.DetachChildren()
			assignNode.DetachChildren()
			definitionSite.DetachChildren()
			recv.DetachChildren()

			methodNameID := rhino.RhinoStringPool_InternID(newMethodName)
			selfName := newMethodName + "$self"
			selfNameID := rhino.RhinoStringPool_InternID(selfName)

			definitionSite.ResetDetachedGetPropToNameID(methodNameID, isConstantName)
			definitionSite.AddChildToFront(fn.Node())
			statement.MutateAttachedStmtToVarFromSrcref(definitionSite, definitionSite, originalJSDoc)

			paramList := fn.Params().Node()
			paramList.AddChildToFront(recv.ResetDetachedToFreshNameID(selfNameID).SrcrefIfMissing(fn.Node()))
			d.compiler.ReportChangeToEnclosingScope(paramList)

			d.replaceReferencesToThisID(paramList, selfName, selfNameID)
			d.replaceReferencesToThisID(fn.Body().Node(), selfName, selfNameID)

			d.compiler.ReportChangeToEnclosingScope(statement.GetParent())
			d.compiler.ReportChangeToEnclosingScope(statement)
			return
		}
	}

	var fn ast.Function
	var subtreeToRemove rhino.Node
	switch definitionSite.GetToken() {
	case rhino.Token_GETPROP:
		fn = ast.MustFunction(ast.MustAssign(definitionSite.GetParent()).Value().Node())
		subtreeToRemove = NodeUtil_GetEnclosingStatement(definitionSite)
	case rhino.Token_STRING_KEY:
		fn = ast.MustFunction(ast.MustStringKey(definitionSite).Value().Node())
		subtreeToRemove = definitionSite
	case rhino.Token_MEMBER_FUNCTION_DEF:
		fn = ast.MustMemberFunctionDef(definitionSite).Function()
		subtreeToRemove = definitionSite
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg(definitionSite.String()))
	}

	// Define a new variable before the original statement.
	statement := NodeUtil_GetEnclosingStatement(definitionSite)
	newName := ast.NewName(newMethodName)
	newName.Node().SrcrefIfMissing(definitionSite)
	if isConstantName {
		newName.SetConstantName(true)
	}
	newVarNode := rhino.IR_Var(newName.Node()).SrcrefIfMissing(definitionSite)
	newVarNode.InsertBefore(statement)
	newVarNode.SetJSDocInfo(NodeUtil_GetBestJSDocInfo(definitionSite))

	// Attach the function to the new variable.
	fn.Detach()
	newName.SetInit(fn.AsExpr())

	// Create the `this` param.
	selfName := newMethodName + "$self"
	selfNameID := rhino.RhinoStringPool_InternID(selfName)
	paramList := fn.Params()
	selfParam := ast.NewName(selfName)
	selfParam.Node().SrcrefIfMissing(fn.Node())
	paramList.PrependParam(selfParam.AsParam())
	d.compiler.ReportChangeToEnclosingScope(paramList.Node())

	// Eliminate `this`.
	d.replaceReferencesToThisID(paramList.Node(), selfName, selfNameID) // In default param values.
	d.replaceReferencesToThisID(fn.Body().Node(), selfName, selfNameID) // In function body.

	// Clean up dangling AST.
	NodeUtil_DeleteNode(subtreeToRemove, d.compiler)
	d.compiler.ReportChangeToEnclosingScope(newVarNode)
}

func (d *DevirtualizeMethods) replaceThisNode(thisNode rhino.Node, name string, nameID uint32) {
	if thisNode.MutateThisToNameIDInPlace(nameID) {
		d.compiler.ReportChangeToEnclosingScope(thisNode)
	} else {
		newName := rhino.IR_Name(name).Srcref(thisNode).CopyTypeFrom(thisNode)
		thisNode.ReplaceWith(newName)
		d.compiler.ReportChangeToEnclosingScope(newName)
	}
}

func (d *DevirtualizeMethods) replaceReferencesToThisID(node rhino.Node, name string, nameID uint32) {
	if (node.IsFunction() && !node.IsArrowFunction()) ||
		node.IsMemberFieldDef() ||
		NodeUtil_IsClassStaticBlock(node) {
		// Functions (besides arrows), member field defs, and static blocks create a new binding for
		// `this`.
		return
	}

	if cf, ok := ast.AsComputedFieldDef(node); ok {
		// Computed field keys are evaluated in the enclosing scope, but field initializers create a
		// new binding for `this`.
		key := cf.Key().Node()
		if !key.IsNil() {
			if key.IsThis() {
				d.replaceThisNode(key, name, nameID)
			} else {
				d.replaceReferencesToThisID(key, name, nameID)
			}
		}
		return
	}

	for child := node.GetFirstChild(); !child.IsNil(); {
		next := child.GetNext()
		if child.IsThis() {
			d.replaceThisNode(child, name, nameID)
		} else {
			d.replaceReferencesToThisID(child, name, nameID)
		}
		child = next
	}
}

// Java: DevirtualizeMethods.chunkForNode(Node)
func (d *DevirtualizeMethods) chunkForNode(node rhino.Node) *JSChunk {
	script := NodeUtil_GetEnclosingScript(node)
	input := d.compiler.GetInput(script.GetInputId())
	return input.GetChunk()
}

// Java: DevirtualizeMethods.isScopingOrBranchingConstruct(Node)
func devirtualizeMethods_isScopingOrBranchingConstruct(node rhino.Node) bool {
	return NodeUtil_IsControlStructure(node) || // Branching.
		node.IsAnd() || // Branching.
		node.IsOr() || // Branching.
		node.IsFunction() || // Branching & scoping.
		node.IsBlock() // Scoping.
}

// Java: DevirtualizeMethods.localNameIsDeclaredByClass(Node)
func devirtualizeMethods_localNameIsDeclaredByClass(cls ast.Class) bool {
	if _, hasName := cls.Name(); !hasName {
		return false
	}
	return !NodeUtil_IsStatement(cls.Node())
}
