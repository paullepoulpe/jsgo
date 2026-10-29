// Port of com/google/javascript/jscomp/InlineAndCollapseProperties.java @ closure-compiler v20261021

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CollapseDeclarationOfNameAndDescendants collapses definitions of the
// collapsible properties of a global name. Recurs on subnames that also
// represent JavaScript objects with collapsible properties.
//
// n is a node representing a global name; alias the flattened name for n.
//
// Java: InlineAndCollapseProperties.CollapseProperties.collapseDeclarationOfNameAndDescendants(Name, String, Set<Name>)
func (c *InlineAndCollapseProperties_CollapseProperties) CollapseDeclarationOfNameAndDescendants(
	n *GlobalNamespace_Name, alias string) {
	childNameInlinability := n.CanCollapseOrInlineChildNames()
	var canCollapseChildNames bool
	if !childNameInlinability.CanCollapse() {
		if c.outer.decisionsLog.IsLogging() {
			c.LogDecisionForNameSupplier(
				n,
				func() string {
					return "child name inlinability: " + childNameInlinability.String() + ": will not collapse child names"
				})
		}
		canCollapseChildNames = false
	} else if n.isEscaped() {
		c.LogDecisionForName(n, "escapes: will not collapse child names")
		canCollapseChildNames = false
	} else {
		canCollapseChildNames = true
	}

	// Handle this name first so that nested object literals get unrolled.
	if c.CanCollapse(n) {
		c.LogDecisionForName(n, "collapsing")
		c.UpdateGlobalNameDeclaration(n, alias, canCollapseChildNames)
	}

	if !n.HasProps() || n.isEscaped() {
		return
	}
	c.LogDecisionForName(n, "collapsing descendants")
	isLogging := c.outer.decisionsLog.IsLogging()
	for p := range n.Props() {
		if !isLogging && !p.HasProps() && !c.CanCollapse(p) {
			continue
		}
		var propAlias string
		if c.CanCollapse(p) || (p.HasProps() && !p.isEscaped()) {
			propAlias = c.AppendPropForAlias(alias, p.GetBaseName())
		}
		c.CollapseDeclarationOfNameAndDescendants(p, propAlias)
	}
}

// UpdateTwinnedDeclaration updates the initial assignment to a collapsible
// property at global scope by adding a VAR stub and collapsing the
// property. e.g. c = a.b = 1; => var a$b; c = a$b = 1; This specifically
// handles "twinned" assignments, which are those where the assignment is
// also used as a reference and which need special handling.
//
// alias is the flattened property name (e.g. "a$b"); refName the name for
// the reference being updated; ref an object containing information about
// the assignment getting updated.
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateTwinnedDeclaration(String, Name, Ref)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateTwinnedDeclaration(alias string, refName *GlobalNamespace_Name, ref *GlobalNamespace_Ref) {
	jutil.CheckStateMsg1(ref.IsTwin(), "%s", ref)
	// Don't handle declarations of an already flat name, just qualified names.
	if !ref.GetNode().IsGetProp() {
		return
	}
	rvalue := ref.GetNode().GetNext()
	parent := ref.GetNode().GetParent()
	grandparent := parent.GetParent()

	if !rvalue.IsNil() && rvalue.IsFunction() {
		c.CheckForReceiverAffectedByCollapse(rvalue, refName.GetJSDocInfo(), refName)
	}

	// Create the new alias node.
	nameNode := NodeUtil_NewNameWithOriginalName(c.outer.compiler, alias, grandparent.GetFirstChild(), refName.GetFullName())
	NodeUtil_CopyNameAnnotations(ref.GetNode(), nameNode)

	// BEFORE:
	// ... (x.y = 3);
	//
	// AFTER:
	// var x$y;
	// ... (x$y = 3);

	current := grandparent
	currentParent := grandparent.GetParent()
	for ; !currentParent.IsScript() && !currentParent.IsBlock(); current, currentParent = currentParent, currentParent.GetParent() {
	}

	// Create a stub variable declaration right
	// before the current statement.
	stubVar := rhino.IR_Var(nameNode.CloneTree()).SrcrefIfMissing(nameNode)
	stubVar.InsertBefore(current)

	ref.GetNode().ReplaceWith(nameNode)
	c.outer.compiler.ReportChangeToEnclosingScope(nameNode)
}

// UpdateGlobalNameDeclaration updates the first initialization (a.k.a
// "declaration") of a global name. This involves flattening the global name
// (if it's not just a global variable name already), collapsing object
// literal keys into global variables, declaring stub global variables for
// properties added later in a local scope.
//
// It may seem odd that this function also takes care of declaring stubs for
// direct children. The ultimate goal of this function is to eliminate the
// global name entirely (when possible), so that "middlemen" namespaces
// disappear, and to do that we need to make sure that all the direct
// children will be collapsed as well.
//
// n is an object representing a global name (e.g. "a", "a.b.c"); alias the
// flattened name for n (e.g. "a", "a$b$c"); canCollapseChildNames whether
// it's possible to collapse children of this name. (This is mostly passed
// for convenience; it's equivalent to n.canCollapseChildNames()).
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclaration(Name, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclaration(n *GlobalNamespace_Name, alias string, canCollapseChildNames bool) {
	decl := n.getDeclarationRef()
	if decl == nil {
		// Some names do not have declarations, because they
		// are only defined in local scopes.
		c.LogDecisionForName(n, "no global declaration found")
		return
	}

	isSimpleName := n.IsSimpleName()
	declNode := decl.GetNode()
	switch declNode.GetParent().GetToken() {
	case rhino.Token_ASSIGN:
		c.LogDeclarationAction(n, declNode, "updating assignment")
		c.UpdateGlobalNameDeclarationAtAssignNode(n, alias, canCollapseChildNames)
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		c.LogDeclarationAction(n, declNode, "updating variable declaration")
		c.UpdateGlobalNameDeclarationAtVariableNode(n, canCollapseChildNames)
	case rhino.Token_FUNCTION:
		c.LogDeclarationAction(n, declNode, "updating function declaration")
		c.UpdateGlobalNameDeclarationAtFunctionNode(n, canCollapseChildNames)
	case rhino.Token_CLASS:
		c.LogDeclarationAction(n, declNode, "updating class declaration")
		c.UpdateGlobalNameDeclarationAtClassNode(n, canCollapseChildNames)
	case rhino.Token_CLASS_MEMBERS:
		c.LogDeclarationAction(n, declNode, "updating static member declaration")
		c.UpdateGlobalNameDeclarationAtStaticMemberNode(n, alias, canCollapseChildNames)
	default:
		c.LogDeclarationAction(
			n, declNode, "not updating an unsupported type of declaration node")
	}

	if !isSimpleName {
		commonAncestor := n.GetDeepestCommonAncestorChunk(c.outer.compiler.GetChunkGraph())
		if commonAncestor != decl.GetChunk() {
			c.MoveDeclarationIntoHigherChunk(decl, commonAncestor)
		}
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.logDeclarationAction(Name, Node, String)
func (c *InlineAndCollapseProperties_CollapseProperties) LogDeclarationAction(name *GlobalNamespace_Name, declarationNode rhino.Node, message string) {
	if c.outer.decisionsLog.IsLogging() {
		c.LogDecisionForNameSupplier(name, func() string { return declarationNode.String() + ": " + message })
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.moveDeclarationIntoHigherChunk(Ref, JSChunk)
func (c *InlineAndCollapseProperties_CollapseProperties) MoveDeclarationIntoHigherChunk(decl *GlobalNamespace_Ref, commonAncestor *JSChunk) {
	// Var hoisting means that early references are generally ok. E.g. `a$b; var a$b;` is valid.
	// However, if `var a$b;` were in a different, conditionally loaded chunk, then it might not
	// be parsed before `a$b;`. So we need to avoid a ReferenceError when referencing `a$b`.
	newNode := decl.GetNode()
	varNode := newNode.GetParent()
	initializer := newNode.GetFirstChild()
	jutil.CheckStateMsg1(varNode.IsVar(), "Unexpected declaration %s", varNode)

	newRoot := commonAncestor.GetFirst().GetAstRoot(c.outer.compiler)
	if initializer.IsNil() {
		// Already a stub declaration - "var a$b;". Move it to the ancestor chunk.
		newRoot.AddChildToFront(varNode.Detach())
	} else {
		// Define `var a$b;` in the ancestor chunk.
		newVarNode := rhino.IR_Var(newNode.CloneNode()).Srcref(newNode)
		newRoot.AddChildToFront(newVarNode)
		// Replace `var a$b = () => {};` with `a$b = () = {};`. Having multiple `var` declarations
		// of the same name is valid JavaScript but violates Normalize constraints.
		newAssign := rhino.IR_ExprResult(rhino.IR_Assign(newNode.Detach(), initializer.Detach())).
			SrcrefTreeIfMissing(newNode)
		varNode.ReplaceWith(newAssign)
	}
	c.outer.compiler.ReportChangeToChangeScope(newRoot)
}

// UpdateGlobalNameDeclarationAtAssignNode updates the first initialization
// (a.k.a "declaration") of a global name that occurs at an ASSIGN node. See
// comment for updateGlobalNameDeclaration.
//
// n is an object representing a global name (e.g. "a", "a.b.c"); alias the
// flattened name for n (e.g. "a", "a$b$c").
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtAssignNode(Name, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtAssignNode(
	n *GlobalNamespace_Name, alias string, canCollapseChildNames bool) {
	// NOTE: It's important that we don't add additional nodes
	// (e.g. a var node before the exprstmt) because the exprstmt might be
	// the child of an if statement that's not inside a block).

	// All qualified names - even for variables that are initially declared as LETS and CONSTS -
	// are being declared as VAR statements, but this is not incorrect because
	// we are only collapsing for global names.
	ref := n.getDeclarationRef()
	refNode := ref.GetNode()
	rvalue := refNode.GetNext()
	if ref.IsTwin() {
		c.UpdateTwinnedDeclaration(alias, n, ref)
		return
	}
	assignNode := refNode.GetParent()
	grandparent := assignNode.GetParent()
	varParent := grandparent.GetParent()
	isObjLit := rvalue.IsObjectLit()
	insertedVarNode := false
	var varNode rhino.Node

	if isObjLit && c.CanEliminate(n, rvalue) {
		// Eliminate the object literal altogether.
		if grandparent.CanResetInPlace() {
			grandparent.MutateAttachedToFreshVarInPlace(rhino.NilNode, nil)
			varNode = grandparent
		} else {
			varNode = rhino.NewNode(rhino.Token_VAR)
			grandparent.ReplaceWith(varNode)
		}
		n.UpdateRefNode(ref, rhino.NilNode)
		insertedVarNode = true
		c.outer.compiler.ReportChangeToEnclosingScope(varNode)
	} else if !n.IsSimpleName() {
		// Create a VAR node to declare the name.
		if rvalue.IsFunction() {
			c.CheckForReceiverAffectedByCollapse(rvalue, n.GetJSDocInfo(), n)
		}

		c.outer.compiler.ReportChangeToEnclosingScope(rvalue)
		rvalue.Detach()

		isConstName := c.ShouldAddConstantName(refNode, alias)
		info := NodeUtil_GetBestJSDocInfo(assignNode)
		var nameNode rhino.Node
		if assignNode.CanResetInPlace() && grandparent.CanResetInPlace() {
			assignNode.Detach()
			nameNode = assignNode.ResetDetachedToNameWithOriginalNameInPlace(alias, n.GetFullName(), grandparent, isConstName)
			nameNode.AddChildToFront(rvalue)
			grandparent.MutateAttachedToFreshVarInPlace(nameNode, info)
			varNode = grandparent
		} else {
			nameNode = NodeUtil_NewNameWithOriginalName(c.outer.compiler, alias, grandparent, n.GetFullName())
			nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, isConstName)
			varNode = rhino.NewNode(rhino.Token_VAR)
			if info != nil {
				varNode.SetJSDocInfo(info)
			}
			varNode.AddChildToBack(nameNode)
			nameNode.AddChildToFront(rvalue)
			grandparent.ReplaceWith(varNode)
		}

		// Update the node ancestry stored in the reference.
		n.UpdateRefNode(ref, nameNode)
		insertedVarNode = true
		c.outer.compiler.ReportChangeToEnclosingScope(varNode)
	} else if canCollapseChildNames {
		varNode = rhino.NewNode(rhino.Token_VAR)
	}

	if canCollapseChildNames {
		if isObjLit {
			c.DeclareVariablesForObjLitValues(n, alias, rvalue, varNode, varNode.GetPrevious())
		}

		c.AddStubsForUndeclaredProperties(n, alias, varParent, varNode)
	}

	if insertedVarNode {
		if !varNode.HasChildren() {
			varNode.Detach()
		}
	}
}

// CheckForReceiverAffectedByCollapse warns about any references to "this" in
// the given FUNCTION. The function is getting collapsed, so the references
// will change.
//
// Java: InlineAndCollapseProperties.CollapseProperties.checkForReceiverAffectedByCollapse(Node, JSDocInfo, Name)
func (c *InlineAndCollapseProperties_CollapseProperties) CheckForReceiverAffectedByCollapse(function rhino.Node, docInfo *rhino.JSDocInfo, name *GlobalNamespace_Name) {
	jutil.CheckState(function.IsFunction())

	if docInfo != nil {
		// Don't rely on type inference for this check.

		if docInfo.IsConstructorOrInterface() {
			return // Ctors and interfaces need to be able to reference `this`
		}
		if docInfo.HasThisType() {
			/*
			 * Use `@this` as a signal that the reference is intentional.
			 *
			 * <p>TODO(b/156823102): This signal also silences the check on all transpiled static
			 * methods.
			 */
			return
		}
	}

	// Use the NodeUtil method so we don't forget to update this logic.
	if NodeUtil_ReferencesOwnReceiver(function) {
		c.outer.compiler.Report(JSError_MakeWithNode(function, InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE, name.GetFullName()))
	}
}

// UpdateGlobalNameDeclarationAtVariableNode updates the first
// initialization (a.k.a "declaration") of a global name that occurs at a VAR
// node. See comment for updateGlobalNameDeclaration.
//
// n is an object representing a global name (e.g. "a").
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtVariableNode(Name, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtVariableNode(n *GlobalNamespace_Name, canCollapseChildNames bool) {
	if !canCollapseChildNames {
		c.LogDecisionForName(n, "cannot collapse child names: skipping")
		return
	}

	ref := n.getDeclarationRef()
	name := ref.GetNode().GetString()
	rvalue := ref.GetNode().GetFirstChild()
	variableNode := ref.GetNode().GetParent()
	grandparent := variableNode.GetParent()

	isObjLit := rvalue.IsObjectLit()

	if isObjLit {
		c.DeclareVariablesForObjLitValues(n, name, rvalue, variableNode, variableNode.GetPrevious())
	}

	c.AddStubsForUndeclaredProperties(n, name, grandparent, variableNode)

	if isObjLit && c.CanEliminate(n, rvalue) {
		ref.GetNode().Detach()
		c.outer.compiler.ReportChangeToEnclosingScope(variableNode)
		if !variableNode.HasChildren() {
			variableNode.Detach()
		}

		// Clear out the object reference, since we've eliminated it from the
		// parse tree.
		n.UpdateRefNode(ref, rhino.NilNode)
	}
}

// UpdateGlobalNameDeclarationAtFunctionNode updates the first
// initialization (a.k.a "declaration") of a global name that occurs at a
// FUNCTION node. See comment for updateGlobalNameDeclaration.
//
// n is an object representing a global name (e.g. "a").
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtFunctionNode(Name, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtFunctionNode(n *GlobalNamespace_Name, canCollapseChildNames bool) {
	if !canCollapseChildNames || !c.CanCollapse(n) {
		return
	}

	ref := n.getDeclarationRef()
	fnName := ref.GetNode().GetString()
	c.AddStubsForUndeclaredProperties(
		n, fnName, ref.GetNode().GetAncestor(2), ref.GetNode().GetParent())
}

// UpdateGlobalNameDeclarationAtClassNode updates the first initialization
// (a.k.a "declaration") of a global name that occurs at a CLASS node. See
// comment for updateGlobalNameDeclaration.
//
// n is an object representing a global name (e.g. "a").
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtClassNode(Name, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtClassNode(n *GlobalNamespace_Name, canCollapseChildNames bool) {
	if !canCollapseChildNames || !c.CanCollapse(n) {
		return
	}

	ref := n.getDeclarationRef()
	className := ref.GetNode().GetString()
	c.AddStubsForUndeclaredProperties(
		n, className, ref.GetNode().GetAncestor(2), ref.GetNode().GetParent())
}

// UpdateGlobalNameDeclarationAtStaticMemberNode updates the first
// initialization (a.k.a "declaration") of a global name that occurs in a
// static MEMBER_FUNCTION_DEF or MEMBER_FIELD_DEF in a class. See comment for
// updateGlobalNameDeclaration.
//
// n is a static MEMBER_FUNCTION_DEF or MEMBER_FIELD_DEF in a class assigned
// to a global name (e.g. `a.b`); alias the new flattened name for `n` (e.g.
// "a$b"); canCollapseChildNames whether properties of `n` are also
// collapsible, meaning that any properties only assigned locally need stub
// declarations.
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtStaticMemberNode(Name, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtStaticMemberNode(
	n *GlobalNamespace_Name, alias string, canCollapseChildNames bool) {

	declaration := n.getDeclarationRef()
	classNode := declaration.GetNode().GetGrandparent()
	jutil.CheckStateMsg1(classNode.IsClass(), "%s", classNode)
	enclosingStatement := NodeUtil_GetEnclosingStatement(classNode)

	if canCollapseChildNames {
		c.AddStubsForUndeclaredProperties(n, alias, enclosingStatement.GetParent(), classNode)
	}

	// detach `static m() {}` from `class Foo { static m() {} }`
	member := declaration.GetNode().Detach()
	var varDecl rhino.Node
	if member.HasChildren() {
		memberVal := member.GetOnlyChild()
		c.CheckForReceiverAffectedByCollapse(memberVal, member.GetJSDocInfo(), n)
		memberVal = memberVal.Detach()

		// add a var declaration, creating `var Foo$m = function() {}; class Foo {}`
		varDecl = rhino.IR_VarWithValue(NodeUtil_NewNameWithSrcref(c.outer.compiler, alias, member), memberVal).Srcref(member)
	} else {
		// Member field defs may not have initializers.
		varDecl = rhino.IR_Var(NodeUtil_NewNameWithSrcref(c.outer.compiler, alias, member)).Srcref(member)
	}
	varDecl.InsertBefore(enclosingStatement)
	// We would lose optimization-relevant jsdoc tags here because they are stored on the class
	// member node, not the function node. Copy them over to the new declaration statement so
	// later passes can make use of them.
	varDecl.SetJSDocInfo(member.GetJSDocInfo())
	c.outer.compiler.ReportChangeToEnclosingScope(varDecl)

	// collapsing this name's properties requires updating this Ref
	n.UpdateRefNode(declaration, varDecl.GetFirstChild())
}

// DeclareVariablesForObjLitValues declares global variables to serve as
// aliases for the values in an object literal, optionally removing all of
// the object literal's keys and values.
//
// alias is the object literal's flattened name (e.g. "a$b$c"); objlit the
// OBJLIT node; varNode the VAR node to which new global variables should be
// added as children; nameToAddAfter the child of varNode after which new
// variables should be added (may be null).
//
// Java: InlineAndCollapseProperties.CollapseProperties.declareVariablesForObjLitValues(Name, String, Node, Node, Node)
func (c *InlineAndCollapseProperties_CollapseProperties) DeclareVariablesForObjLitValues(
	objlitName *GlobalNamespace_Name, alias string, objlit rhino.Node, varNode rhino.Node, nameToAddAfter rhino.Node) {
	arbitraryNameCounter := 0
	discardKeys := !objlitName.ShouldKeepKeys()

	for member := range ast.ReadAhead(ast.MustObjectLit(objlit).FirstMember()) {
		key := member.Node()
		value := key.GetFirstChild()

		// A computed property, or a get or a set can not be rewritten as a VAR. We don't know what
		// properties will be generated by a spread.
		switch key.GetToken() {
		case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_COMPUTED_PROP, rhino.Token_OBJECT_SPREAD:
			continue
		case rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF:
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of OBJECTLIT: " + key.ToStringTree()))
		}

		// We generate arbitrary names for keys that aren't valid JavaScript
		// identifiers, since those keys are never referenced. (If they were,
		// this object literal's child names wouldn't be collapsible.) The only
		// reason that we don't eliminate them entirely is the off chance that
		// their values are expressions that have side effects.
		isJsIdentifier := !key.IsNumber() && rhino.TokenStream_IsJSIdentifier(key.GetString())
		var propName string
		if isJsIdentifier {
			propName = key.GetString()
		} else {
			arbitraryNameCounter++
			propName = strconv.Itoa(arbitraryNameCounter)
		}

		// If the name cannot be collapsed, skip it.
		p := objlitName.LookupChild(propName)
		if p != nil && !c.CanCollapse(p) {
			continue
		}

		propAlias := c.AppendPropForAlias(alias, propName)
		var refNode rhino.Node
		isConstantName := c.ShouldAddConstantName(key, propAlias)
		keyDoc := key.GetJSDocInfo()
		if discardKeys {
			key.Detach()
			value.Detach()
			// Don't report a change here because the objlit has already been removed from the tree.
		} else {
			// Substitute a reference for the value.
			refNode = rhino.IR_Name(propAlias)
			if isConstantName {
				refNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
			}

			value.ReplaceWith(refNode)
			c.outer.compiler.ReportChangeToEnclosingScope(refNode)
		}

		// Declare the collapsed name as a variable with the original value.
		nameNode := rhino.IR_Name(propAlias)
		nameNode.AddChildToFront(value)
		if isConstantName {
			nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}
		var newVar rhino.Node
		if discardKeys && key.CanResetInPlace() && (p == nil || p.getDeclarationRef() == nil || p.getDeclarationRef().GetNode() == key) {
			nameNode.SrcrefTreeIfMissing(key)
			key.MutateAttachedStmtToVarFromSrcref(nameNode, key, nil)
			newVar = key
		} else {
			newVar = rhino.IR_Var(nameNode).SrcrefTreeIfMissing(key)
		}
		if !nameToAddAfter.IsNil() {
			newVar.InsertAfter(nameToAddAfter)
		} else {
			newVar.InsertBefore(varNode)
		}
		c.outer.compiler.ReportChangeToEnclosingScope(newVar)
		nameToAddAfter = newVar

		// Update the global name's node ancestry if it hasn't already been
		// done. (Duplicate keys in an object literal can bring us here twice
		// for the same global name.)
		if isJsIdentifier && p != nil {
			if !discardKeys {
				p.AddAliasingGetClonedFromDeclaration(refNode)
			}

			p.UpdateRefNode(p.getDeclarationRef(), nameNode)

			if value.IsFunction() {
				c.CheckForReceiverAffectedByCollapse(value, keyDoc, p)
			}
		}
	}
}

// AddStubsForUndeclaredProperties adds global variable "stubs" for any
// properties of a global name that are only set in a local scope or read
// but never set.
//
// n is an object representing a global name (e.g. "a", "a.b.c"); alias the
// flattened name of the object whose properties we are adding stubs for
// (e.g. "a$b$c"); parent the node to which new global variables should be
// added as children; addAfter the child of after which new variables should
// be added.
//
// Java: InlineAndCollapseProperties.CollapseProperties.addStubsForUndeclaredProperties(Name, String, Node, Node)
func (c *InlineAndCollapseProperties_CollapseProperties) AddStubsForUndeclaredProperties(
	n *GlobalNamespace_Name, alias string, parent rhino.Node, addAfter rhino.Node) {
	jutil.CheckStateMsg1(n.CanCollapseUnannotatedChildNames(), "%s", n)
	jutil.CheckArgumentMsg1(NodeUtil_IsStatementBlock(parent), "%s", parent)
	rhino.CheckNotNull(addAfter)
	if !n.HasProps() {
		return
	}
	for p := range n.Props() {
		if !p.NeedsToBeStubbed() {
			continue
		}

		propAlias := c.AppendPropForAlias(alias, p.GetBaseName())
		nameNode := rhino.IR_Name(propAlias)
		newVar := rhino.IR_Var(nameNode).SrcrefTreeIfMissing(addAfter)
		newVar.InsertAfter(addAfter)

		// Determine if this is a constant var by checking the first
		// reference to it. Don't check the declaration, as it might be null.
		constPropNode := p.GetFirstRef().GetNode()
		isConstantName :=
			// Don't call shouldAddConstantName because it expects to see a declaration node for
			// a name, and this name is, by definition, undeclared.
			constPropNode.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) ||
				c.outer.compiler.GetCodingConvention().IsConstant(propAlias)
		if isConstantName {
			nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}

		c.outer.compiler.ReportChangeToEnclosingScope(newVar)
		addAfter = newVar
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.appendPropForAlias(String, String)
func (c *InlineAndCollapseProperties_CollapseProperties) AppendPropForAlias(root string, prop string) string {
	if strings.IndexByte(prop, '$') != -1 {
		// Encode '$' in a property as '$0'. Because '0' cannot be the
		// start of an identifier, this will never conflict with our
		// encoding from '.' -> '$'.
		prop = strings.ReplaceAll(prop, "$", "$0")
	}
	result := root + "$" + prop
	id := 1
	for c.namespace.ContainsName(result) {
		result = root + "$" + prop + "$" + strconv.Itoa(id)
		id++
	}
	return result
}
