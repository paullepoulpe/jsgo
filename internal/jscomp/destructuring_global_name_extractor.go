// Port of com/google/javascript/jscomp/DestructuringGlobalNameExtractor.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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

// DestructuringGlobalNameExtractor is a helper for changing the value of an
// lvalue in a destructuring pattern. Intended for use by
// InlineAndCollapseProperties.CollapseProperties and
// InlineAndCollapseProperties.AggressiveInlineAliases only. This class makes
// some assumptions that don't generally hold in order to preserve
// GlobalNamespace validity and avoid creating temporary variables.
type DestructuringGlobalNameExtractor struct{}

// DestructuringGlobalNameExtractor_ReassignDestructringLvalue, given an
// lvalue in a destructuring pattern, and a detached subtree, rewrites the
// AST to assign the lvalue to the subtree instead of its previous value,
// while preserving the rest of the destructuring pattern
//
// For example: given stringKey: 'y' and newName: 'new.name', where 'y' is
// contained in the declaration `const {x, y, z} = original;`, this method
// will produce `const {x} = original; const y = new.name; const {z} =
// original;`
//
// This method only handles a limited subset of destructuring patterns and
// is intended only for CollapseProperties and AggressiveInlineAliases.
// Preconditions are:
//
//   - the pattern must be the lhs of an assignment or declaration (e.g. not
//     a nested pattern)
//   - the original rvalue for the pattern must be an effectively constant
//     qualified name. it is safe to evaluate the rvalue multiple times.
//
// stringKey is the STRING_KEY node representing the ref to rewrite, e.g.
// `y`; newName what the lvalue in the STRING_KEY should now be assigned to,
// e.g. foo.bar; newNodes optionally a set to add new AstChanges to, if you
// need to keep the GlobalNamespace up-to-date. Otherwise null. ref is the
// Ref corresponding to the `stringKey`.
//
// Java: DestructuringGlobalNameExtractor.reassignDestructringLvalue(Node, Node, Set<AstChange>, Ref, AbstractCompiler)
func DestructuringGlobalNameExtractor_ReassignDestructringLvalue(
	stringKey rhino.Node,
	newName rhino.Node,
	newNodes *astChangeSet,
	ref *GlobalNamespace_Ref,
	compiler AbstractCompiler) {
	pattern := stringKey.GetParent()
	assignmentType := pattern.GetParent()
	jutil.CheckStateMsg1(assignmentType.IsAssign() || assignmentType.IsDestructuringLhs(), "%s", assignmentType)
	originalRvalue := pattern.GetNext()                // e.g. `original`
	jutil.CheckState(originalRvalue.IsQualifiedName()) // don't handle rvalues with side effects

	// Create a new assignment using the provided qualified name, e.g. `const y = new.name;`
	var lvalueToReassign rhino.Node
	if stringKey.GetOnlyChild().IsDefaultValue() {
		lvalueToReassign = stringKey.GetOnlyChild().GetFirstChild()
	} else {
		lvalueToReassign = stringKey.GetOnlyChild()
	}
	if newNodes != nil {
		newNodes.add(NewGlobalNamespace_RefBasedAstChange(ref, newName))
	}
	rvalue := DestructuringGlobalNameExtractor_MakeNewRvalueForDestructuringKey(stringKey, newName, newNodes, ref)

	// Add that new assignment to the AST after the original destructuring pattern
	if stringKey.GetPrevious().IsNil() && !stringKey.GetNext().IsNil() {
		// Remove the original item completely if the given string key doesn't have any preceding
		// string keys, /and/ it has succeeding string keys.
		DestructuringGlobalNameExtractor_ReplaceDestructuringAssignment(pattern, lvalueToReassign.Detach(), rvalue)
	} else {
		DestructuringGlobalNameExtractor_AddAfter(pattern, lvalueToReassign.Detach(), rvalue)
	}

	// Remove any lvalues in the pattern that are after the given stringKey, and put them in
	// a second destructuring pattern. This is necessary in case they depend on side effects from
	// assigning `lvalueToReassign`. e.g. create `const {z} = original;`
	if !stringKey.GetNext().IsNil() {
		newPattern := DestructuringGlobalNameExtractor_CreateNewObjectPatternFromSuccessiveKeys(stringKey).Srcref(pattern)

		// reuse the original rvalue if we don't need it for earlier keys, otherwise make a copy.
		var newRvalue rhino.Node
		if stringKey.GetPrevious().IsNil() {
			newRvalue = originalRvalue.Detach()
		} else {
			newRvalue = originalRvalue.CloneTree()
			if newNodes != nil {
				newNodes.add(NewGlobalNamespace_RefBasedAstChange(ref, newRvalue))
			}
		}
		DestructuringGlobalNameExtractor_AddAfter(lvalueToReassign, newPattern, newRvalue)
	}

	stringKey.Detach()
	compiler.ReportChangeToEnclosingScope(lvalueToReassign)
}

// DestructuringGlobalNameExtractor_AddAfter adds the new assign or name
// declaration after the original assign or name declaration.
//
// Java: DestructuringGlobalNameExtractor.addAfter(Node, Node, Node)
func DestructuringGlobalNameExtractor_AddAfter(originalLvalue rhino.Node, newLvalue rhino.Node, newRvalue rhino.Node) {
	parent := originalLvalue.GetParent()
	if parent.IsAssign() {
		// create `(<originalLvalue = ...>, <newLvalue = newRvalue>)`
		newAssign := rhino.IR_Assign(newLvalue, newRvalue).Srcref(parent)
		newComma := rhino.NewNodeWithChild(rhino.Token_COMMA, newAssign)
		parent.ReplaceWith(newComma)
		newComma.AddChildToFront(parent)
		return
	}
	// This must have been in a var/let/const.
	if newLvalue.IsDestructuringPattern() {
		newLvalue = rhino.NewNodeWith2Children(rhino.Token_DESTRUCTURING_LHS, newLvalue, newRvalue).Srcref(parent)
	} else {
		newLvalue.AddChildToBack(newRvalue)
	}
	var declaration rhino.Node
	if parent.IsDestructuringLhs() {
		declaration = originalLvalue.GetGrandparent()
	} else {
		declaration = parent
	}
	jutil.CheckStateMsg1(NodeUtil_IsNameDeclaration(declaration), "%s", declaration)
	if NodeUtil_IsStatementParent(declaration.GetParent()) {
		// `const {} = originalRvalue; const newLvalue = newRvalue;`
		// create an entirely new statement
		newDeclaration := rhino.NewNode(declaration.GetToken()).Srcref(declaration)
		newDeclaration.AddChildToBack(newLvalue)
		newDeclaration.InsertAfter(declaration)
	} else {
		// `const {} = originalRvalue, newLvalue = newRvalue;`
		// The Normalize pass tries to ensure name declarations are always in statement blocks, but
		// currently has not implemented normalization for `for (let x = 0; ...`
		// so we can't add a new statement
		declaration.AddChildToBack(newLvalue)
	}
}

// DestructuringGlobalNameExtractor_ReplaceDestructuringAssignment replaces
// the given assignment or declaration with the new lvalue/rvalue.
//
// pattern is a destructuring pattern in an ASSIGN or VAR/LET/CONST.
//
// Java: DestructuringGlobalNameExtractor.replaceDestructuringAssignment(Node, Node, Node)
func DestructuringGlobalNameExtractor_ReplaceDestructuringAssignment(pattern rhino.Node, newLvalue rhino.Node, newRvalue rhino.Node) {
	parent := pattern.GetParent()
	if parent.IsAssign() {
		newAssign := rhino.IR_Assign(newLvalue, newRvalue).Srcref(parent)
		parent.ReplaceWith(newAssign)
	} else if newLvalue.IsName() {
		jutil.CheckState(parent.IsDestructuringLhs())
		parent.ReplaceWith(newLvalue)
		newLvalue.AddChildToBack(newRvalue)
	} else {
		pattern.GetNext().Detach()
		pattern.Detach()
		parent.AddChildToBack(newLvalue)
		parent.AddChildToBack(newRvalue)
	}
}

// DestructuringGlobalNameExtractor_MakeNewRvalueForDestructuringKey makes a
// default value expression from the rvalue, or otherwise just returns it
//
// e.g. for `const {x = defaultValue} = y;`, and the new rvalue `rvalue`,
// returns `void 0 === rvalue ? defaultValue : rvalue`
//
// Java: DestructuringGlobalNameExtractor.makeNewRvalueForDestructuringKey(Node, Node, Set<AstChange>, Ref)
func DestructuringGlobalNameExtractor_MakeNewRvalueForDestructuringKey(
	stringKey rhino.Node, rvalue rhino.Node, newNodes *astChangeSet, ref *GlobalNamespace_Ref) rhino.Node {
	if stringKey.GetOnlyChild().IsDefaultValue() {
		defaultValue := stringKey.GetFirstChild().GetSecondChild().Detach()
		// Assume `rvalue` has no side effects since it's a qname, and we can create multiple
		// references to it. This ignores getters/setters.
		rvalueForSheq := rvalue.CloneTree()
		if newNodes != nil {
			newNodes.add(NewGlobalNamespace_RefBasedAstChange(ref, rvalueForSheq))
		}
		// `void 0 === rvalue ? defaultValue : rvalue`
		rvalue = rhino.IR_Hook(rhino.IR_Sheq(NodeUtil_NewUndefinedNode(rvalue), rvalueForSheq), defaultValue, rvalue).
			SrcrefTree(defaultValue)
	}
	return rvalue
}

// DestructuringGlobalNameExtractor_CreateNewObjectPatternFromSuccessiveKeys
// removes any keys after the given key, and adds them in order to a new
// object pattern.
//
// Java: DestructuringGlobalNameExtractor.createNewObjectPatternFromSuccessiveKeys(Node)
func DestructuringGlobalNameExtractor_CreateNewObjectPatternFromSuccessiveKeys(stringKey rhino.Node) rhino.Node {
	newPattern := stringKey.GetParent().CloneNode() // copies the JSType
	for next := stringKey.GetNext(); !next.IsNil(); {
		newKey := next
		next = newKey.GetNext()
		newPattern.AddChildToBack(newKey.Detach())
	}
	return newPattern
}

// Java: DestructuringGlobalNameExtractor()
func NewDestructuringGlobalNameExtractor() *DestructuringGlobalNameExtractor {
	return &DestructuringGlobalNameExtractor{}
}
