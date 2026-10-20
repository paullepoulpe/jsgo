// Port of com/google/javascript/jscomp/InlineObjectLiterals.java @ closure-compiler v20261015

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InlineObjectLiterals: using the infrastructure provided by
// ReferenceCollector, identify variables that are only ever assigned to
// object literals and that are never used in entirety, and expand the
// objects into individual variables.
//
// Based on the InlineVariables pass
type InlineObjectLiterals struct {
	compiler AbstractCompiler

	safeNameIdSupplier func() string
}

const InlineObjectLiterals_VAR_PREFIX = "JSCompiler_object_inline_"
const InlineObjectLiterals_STRING_KEY_IDENTIFIER = "string_key"

// Java: InlineObjectLiterals(AbstractCompiler, Supplier<String>)
func NewInlineObjectLiterals(compiler AbstractCompiler, safeNameIdSupplier func() string) *InlineObjectLiterals {
	return &InlineObjectLiterals{compiler: compiler, safeNameIdSupplier: safeNameIdSupplier}
}

var inlineObjectLiterals_isLocalVar = func(v *Var) bool { return v.IsLocal() }

// Java: InlineObjectLiterals.process(Node, Node)
func (o *InlineObjectLiterals) Process(externs *rhino.Node, root *rhino.Node) {
	callback := NewReferenceCollectorWithVarFilter(
		o.compiler,
		NewInlineObjectLiterals_InliningBehavior(o),
		NewSyntacticScopeCreator(o.compiler),
		inlineObjectLiterals_isLocalVar)
	callback.Process(externs, root)
}

// InlineObjectLiterals_InliningBehavior builds up information about nodes in
// each scope. When exiting the scope, inspects all variables in that scope,
// and inlines any that we can.
type InlineObjectLiterals_InliningBehavior struct {
	outer *InlineObjectLiterals

	// A list of variables that should not be inlined, because their
	// reference information is out of sync with the state of the AST.
	staleVars *jutil.LinkedHashSet[*Var]

	validProperties jutil.LinkedHashSet[string]
}

// Java: InlineObjectLiterals.InliningBehavior()
func NewInlineObjectLiterals_InliningBehavior(outer *InlineObjectLiterals) *InlineObjectLiterals_InliningBehavior {
	return &InlineObjectLiterals_InliningBehavior{outer: outer, staleVars: jutil.NewLinkedHashSet[*Var]()}
}

// Java: InlineObjectLiterals.InliningBehavior.afterExitScope(NodeTraversal, ReferenceMap)
func (b *InlineObjectLiterals_InliningBehavior) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	if t.InGlobalScope() {
		return
	}
	if len(t.scopes) > 0 && t.scopes[len(t.scopes)-1].scope == nil {
		return
	}
	scope := t.GetScope()
	if scope.GetVarCount() == 0 {
		return
	}
	for v := range scope.GetVarIterable() {
		if b.IsVarInlineForbidden(v) {
			continue
		}

		referenceInfo := referenceMap.GetReferences(v)
		if referenceInfo == nil {
			continue
		}

		if b.IsInlinableObject(referenceInfo.references) {
			// Skiplist the object itself, as well as any other values
			// that it refers to, since they will have been moved around.
			b.staleVars.Add(v)

			init := referenceInfo.GetInitializingReference()

			// Split up the object into individual variables if the object
			// is never referenced directly in full.
			b.SplitObject(v, init, referenceInfo)
		}
	}
}

// RecordStaleVarReferencesInTree: if there are any variable references in
// the given node tree, skiplist them to prevent the pass from trying to
// inline the variable. Any code modifications will have potentially made the
// ReferenceCollection invalid.
//
// Java: InlineObjectLiterals.InliningBehavior.recordStaleVarReferencesInTree(Node, Scope)
func (b *InlineObjectLiterals_InliningBehavior) RecordStaleVarReferencesInTree(root *rhino.Node, scope *Scope) {
	if root.IsName() {
		b.staleVars.Add(scope.GetVar(root.GetString()))
	}
	if !root.IsFunction() {
		for c := root.GetFirstChild(); c != nil; c = c.GetNext() {
			b.RecordStaleVarReferencesInTree(c, scope)
		}
	}
}

// IsVarInlineForbidden reports whether the given variable is forbidden from
// being inlined.
//
// Java: InlineObjectLiterals.InliningBehavior.isVarInlineForbidden(Var)
func (b *InlineObjectLiterals_InliningBehavior) IsVarInlineForbidden(var_ *Var) bool {
	// A variable may not be inlined if:
	// 1) The variable is defined in the externs
	// 2) The variable is exported,
	// 3) Don't inline the special RENAME_PROPERTY_FUNCTION_NAME
	// 4) A reference to the variable has been inlined. We're downstream
	//    of the mechanism that creates variable references, so we don't
	//    have a good way to update the reference. Just punt on it.

	// Additionally, exclude global variables for now.

	return var_.IsGlobal() ||
		var_.IsExtern() ||
		b.outer.compiler.GetCodingConvention().IsExportedWithLocal(var_.GetName(), true /* local= */) ||
		b.outer.compiler.GetCodingConvention().IsPropertyRenameFunction(var_.GetNameNode()) ||
		b.staleVars.Contains(var_)
}

// IsInlinableObject counts the number of direct (full) references to an
// object. Specifically, we check for references of the following type:
//
//	x;
//	x.fn();
//
// Java: InlineObjectLiterals.InliningBehavior.isInlinableObject(List<Reference>)
func (b *InlineObjectLiterals_InliningBehavior) IsInlinableObject(refs *jutil.ArrayList[*Reference]) bool {
	ret := false
	validProperties := &b.validProperties
	if !validProperties.IsEmpty() {
		validProperties.Clear()
	}
	refSlice := refs.Slice()
	for _, ref := range refSlice {
		name := ref.GetNode()
		parent := ref.GetParent()
		grandparent := ref.GetGrandparent()

		// Ignore most indirect references, like x.y (but not x.y(),
		// since the function referenced by y might reference 'this').
		//
		if NodeUtil_IsNormalOrOptChainGetProp(parent) {
			jutil.CheckState(parent.GetFirstChild() == name)
			// A call target may be using the object as a 'this' value.
			if NodeUtil_IsInvocationTarget(parent) {
				return false
			}

			// Deleting a property has different semantics from deleting
			// a variable, so deleted properties should not be inlined.
			if grandparent.IsDelProp() {
				return false
			}

			// NOTE(nicksantos): This pass's object-splitting algorithm has
			// a blind spot. It assumes that if a property isn't defined on an
			// object, then the value is undefined. This is not true, because
			// Object.prototype can have arbitrary properties on it.
			//
			// We short-circuit this problem by bailing out if we see a reference
			// to a property that isn't defined on the object literal. This
			// isn't a perfect algorithm, but it should catch most cases.
			propName := parent.GetString()
			if !validProperties.Contains(propName) {
				if NodeUtil_IsNameDeclOrSimpleAssignLhs(parent, grandparent) {
					validProperties.Add(propName)
				} else {
					return false
				}
			}
			continue
		}

		// Only rewrite VAR declarations or simple assignment statements
		if !b.IsVarOrAssignExprLhs(name) {
			return false
		}

		// Don't try to handle rewriting VAR/CONST/LET declarations inside for loops.
		// Currently, normalization moves var declarations out of for loop initializers anyway.
		// let/const are more difficult. Declaring each property outside the for loop puts them
		// in an incorrect scope. Declaring them in the loop would initialize them multiple times.
		if NodeUtil_IsNameDeclaration(parent) && NodeUtil_IsAnyFor(grandparent) {
			return false
		}

		val := ref.GetAssignedValue()
		if val == nil {
			// A var with no assignment.
			continue
		}

		// We're looking for object literal assignments only.
		if !val.IsObjectLit() {
			return false
		}

		// Make sure that the value is not self-referential. IOW,
		// disallow things like x = {b: x.a}.
		//
		// TODO(dimvar): Only exclude unorderable self-referential
		// assignments. i.e. x = {a: x.b, b: x.a} is not orderable,
		// but x = {a: 1, b: x.a} is.
		//
		// Also, ES5 getters/setters aren't handled by this pass.
		for child := val.GetFirstChild(); child != nil; child = child.GetNext() {
			switch child.GetToken() {
			case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_COMPUTED_PROP, rhino.Token_OBJECT_SPREAD:
				// ES5 get/set not supported.
				// Don't inline computed property names
				// Spread can overwrite any preceding prop if there are matching keys.
				// TODO(b/126567617): Allow inlining props declared after the SPREAD.
				return false
			case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_STRING_KEY:
			default:
				panic(jutil.NewIllegalStateExceptionMsg(
					"Unexpected child of OBJECTLIT: " + child.ToStringTree()))
			}

			validProperties.Add(child.GetString())

			childVal := child.GetFirstChild()
			// Check if childVal is the parent of any of the passed in
			// references, as that is how self-referential assignments
			// will happen.
			for _, t := range refSlice {
				refNode := t.GetParent()
				for !NodeUtil_IsStatementBlock(refNode) {
					if refNode == childVal {
						// There's a self-referential assignment
						return false
					}
					refNode = refNode.GetParent()
				}
			}
		}

		// We have found an acceptable object literal assignment. As
		// long as there are no other assignments that mess things up,
		// we can inline.
		ret = true
	}
	return ret
}

// Java: InlineObjectLiterals.InliningBehavior.isVarOrAssignExprLhs(Node)
func (b *InlineObjectLiterals_InliningBehavior) IsVarOrAssignExprLhs(n *rhino.Node) bool {
	parent := n.GetParent()
	return NodeUtil_IsNameDeclaration(parent) ||
		(parent.IsAssign() &&
			parent.GetFirstChild() == n &&
			parent.GetParent().IsExprResult())
}

// ComputeVarList computes a list of ever-referenced keys in the object being
// inlined, and returns a mapping of key name -> generated variable name.
//
// Java: InlineObjectLiterals.InliningBehavior.computeVarList(ReferenceCollection)
func (b *InlineObjectLiterals_InliningBehavior) ComputeVarList(referenceInfo *ReferenceCollection) *jutil.LinkedHashMap[string, string] {
	varmap := jutil.NewLinkedHashMap[string, string]()

	for _, ref := range referenceInfo.references.Slice() {
		if ref.IsLvalue() || ref.IsInitializingDeclaration() {
			val := ref.GetAssignedValue()
			if val != nil {
				jutil.CheckStateMsg1(val.IsObjectLit(), "%s", val)
				for child := val.GetFirstChild(); child != nil; child = child.GetNext() {
					varname := child.GetString()
					if varmap.ContainsKey(varname) {
						continue
					}
					var_ := varname
					if !rhino.TokenStream_IsJSIdentifier(varname) {
						var_ = InlineObjectLiterals_STRING_KEY_IDENTIFIER
					}
					var_ = InlineObjectLiterals_VAR_PREFIX + var_ + "_" + b.outer.safeNameIdSupplier()
					varmap.Put(varname, var_)
				}
			}
		} else if NodeUtil_IsNameDeclaration(ref.GetParent()) {
			// This is the var. There is no value.
		} else {
			getprop := ref.GetParent()
			jutil.CheckStateMsg1(NodeUtil_IsNormalOrOptChainGetProp(getprop), "%s", getprop)

			// The key being looked up in the original map.
			varname := getprop.GetString()
			if varmap.ContainsKey(varname) {
				continue
			}

			var_ := InlineObjectLiterals_VAR_PREFIX + varname + "_" + b.outer.safeNameIdSupplier()
			varmap.Put(varname, var_)
		}
	}

	return varmap
}

// FillInitialValues populates a map of key names -> initial assigned values.
// The object literal these are being pulled from is invalidated as a
// result.
//
// Java: InlineObjectLiterals.InliningBehavior.fillInitialValues(Reference, Map<String, Node>)
func (b *InlineObjectLiterals_InliningBehavior) FillInitialValues(init *Reference, initvals *jutil.LinkedHashMap[string, *rhino.Node]) {
	object := init.GetAssignedValue()
	jutil.CheckStateMsg1(object.IsObjectLit(), "%s", object)
	for key := object.GetFirstChild(); key != nil; key = key.GetNext() {
		initvals.Put(key.GetString(), key.RemoveFirstChild())
	}
}

// ReplaceAssignmentExpression replaces an assignment like x = {...} with
// t1=a,t2=b,t3=c,true. Note that the resulting expression will always
// evaluate to true, as would the x = {...} expression.
//
// Java: InlineObjectLiterals.InliningBehavior.replaceAssignmentExpression(Var, Reference, Map<String, String>)
func (b *InlineObjectLiterals_InliningBehavior) ReplaceAssignmentExpression(v *Var, ref *Reference, varmap *jutil.LinkedHashMap[string, string]) {
	// Compute all of the assignments necessary
	var nodes []*rhino.Node
	val := ref.GetAssignedValue()
	b.RecordStaleVarReferencesInTree(val, v.GetScope().(*Scope))
	jutil.CheckStateMsg1(val.IsObjectLit(), "%s", val)
	all := jutil.NewLinkedHashSet(varmap.KeySlice()...)
	for key := val.GetFirstChild(); key != nil; key = key.GetNext() {
		var_ := key.GetString()
		value := key.RemoveFirstChild()
		// TODO(user): Copy type information.
		nodes = append(nodes, rhino.IR_Assign(rhino.IR_Name(varmap.Get(var_)), value))
		all.Remove(var_)
	}

	// TODO(user): Better source information.
	for var_ := range all.All() {
		nodes = append(nodes, rhino.IR_Assign(rhino.IR_Name(varmap.Get(var_)), NodeUtil_NewUndefinedNode(nil)))
	}

	var replacement *rhino.Node
	if len(nodes) == 0 {
		replacement = rhino.IR_TrueNode()
	} else {
		// All assignments evaluate to true, so make sure that the
		// expr statement evaluates to true in case it matters.
		nodes = append(nodes, rhino.IR_TrueNode())

		// Join these using COMMA.  A COMMA node must have 2 children, so we
		// create a tree. In the tree the first child be the COMMA to match
		// the parser, otherwise tree equality tests fail.
		// Lists.reverse(nodes): the original order is not used again.
		slices.Reverse(nodes)
		replacement = rhino.NewNode(rhino.Token_COMMA)
		cur := replacement
		var i int
		for i = 0; i < len(nodes)-2; i++ {
			cur.AddChildToFront(nodes[i])
			t := rhino.NewNode(rhino.Token_COMMA)
			cur.AddChildToFront(t)
			cur = t
		}
		cur.AddChildToFront(nodes[i])
		cur.AddChildToFront(nodes[i+1])
	}

	replace := ref.GetParent()
	replacement.SrcrefTreeIfMissing(replace)

	if NodeUtil_IsNameDeclaration(replace) {
		replace.ReplaceWith(NodeUtil_NewExpr(replacement))
	} else {
		replace.ReplaceWith(replacement)
	}
}

// SplitObject splits up the object literal into individual variables, and
// updates all uses.
//
// Java: InlineObjectLiterals.InliningBehavior.splitObject(Var, Reference, ReferenceCollection)
func (b *InlineObjectLiterals_InliningBehavior) SplitObject(v *Var, init *Reference, referenceInfo *ReferenceCollection) {
	// First figure out the FULL set of possible keys, so that they
	// can all be properly set as necessary.
	varmap := b.ComputeVarList(referenceInfo)

	initvals := jutil.NewLinkedHashMap[string, *rhino.Node]()
	// Figure out the top-level of the var assign node. If it's a plain
	// ASSIGN, then there's an EXPR_STATEMENT above it, if it's a
	// VAR then it should be directly replaced.
	var vnode *rhino.Node
	defined := referenceInfo.IsWellDefined() && NodeUtil_IsNameDeclaration(init.GetParent())
	if defined {
		vnode = init.GetParent()
		b.FillInitialValues(init, initvals)
	} else if v.GetParentNode().IsLet() || v.GetParentNode().IsConst() {
		// Find the beginning of the current scope.
		scopeRoot := v.GetScope().(*Scope).GetRootNode()
		// Assuming scopeRoot is a BLOCK, then we want to insert at the top of the block, before the
		// first statement.
		vnode = scopeRoot.GetFirstChild()
		// Some scope-creating nodes might not be BLOCK nodes
		if !NodeUtil_IsStatement(vnode) && NodeUtil_IsStatement(scopeRoot) {
			vnode = scopeRoot
		} else if scopeRoot.IsSwitchBody() {
			vnode = scopeRoot.GetParent()
		}
	} else {
		// Find the beginning of the function body / script.
		vnode = v.GetScope().(*Scope).GetClosestHoistScope().GetRootNode().GetFirstChild()
	}

	jutil.CheckStateMsg1(NodeUtil_IsStatement(vnode), "%s", vnode)

	for key, value := range varmap.All() {
		val := initvals.Get(key)
		newVarNode := NodeUtil_NewVarNode(value, val)
		if val == nil {
			// is this right?
			newVarNode.SrcrefTreeIfMissing(vnode)
		} else {
			b.RecordStaleVarReferencesInTree(val, v.GetScope().(*Scope))
		}
		newVarNode.InsertBefore(vnode)
		b.outer.compiler.ReportChangeToEnclosingScope(vnode)
	}

	if defined {
		b.outer.compiler.ReportChangeToEnclosingScope(vnode.GetParent())
		vnode.Detach()
	}

	for _, ref := range referenceInfo.references.Slice() {

		// The init/decl have already been converted.
		if defined && ref == init {
			continue
		}
		b.outer.compiler.ReportChangeToEnclosingScope(ref.GetNode())

		if ref.IsLvalue() {
			// Assignments have to be handled specially, since they
			// expand out into multiple assignments.
			b.ReplaceAssignmentExpression(v, ref, varmap)
		} else if NodeUtil_IsNameDeclaration(ref.GetParent()) {
			// The old variable declaration. It didn't have a
			// value. Remove it entirely as it should now be unused.
			ref.GetParent().Detach()
		} else {
			// Make sure that the reference is a GETPROP as we expect it to be.
			getprop := ref.GetParent()
			jutil.CheckStateMsg1(NodeUtil_IsNormalOrOptChainGetProp(getprop), "%s", getprop)

			// The key being looked up in the original map.
			var_ := getprop.GetString()

			// If the variable hasn't already been declared, add an empty
			// declaration near all the other declarations.
			jutil.CheckState(varmap.ContainsKey(var_))

			// Replace the GETPROP node with a NAME.
			replacement := rhino.IR_Name(varmap.Get(var_))
			replacement.SrcrefIfMissing(getprop)
			ref.GetParent().ReplaceWith(replacement)
		}
	}
}
