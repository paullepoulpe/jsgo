// Port of com/google/javascript/jscomp/StripCode.java @ closure-compiler v20261021

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
	"reflect"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// StripCode is a pass for stripping a list of provided JavaScript object
// types.
//
// The stripping strategy is as follows:
//
//   - Provide: 1) a list of types that should be stripped, and 2) a list of
//     suffixes of field/variable names that should be stripped.
//   - Remove declarations of variables that are initialized using static
//     methods of strip types (e.g. var x = goog.debug.Logger.getLogger(...);).
//   - Remove all references to variables that are stripped.
//   - Remove all object literal keys with strip names.
//   - Remove all assignments to 1) field names that are strip names and 2)
//     qualified names that begin with strip types.
//   - Remove all statements containing calls to static methods of strip
//     types.
type StripCode struct {
	compiler          AbstractCompiler
	stripNameSuffixes *jutil.LinkedHashSet[string]
	stripNamePrefixes *jutil.LinkedHashSet[string]
	varsToRemove      map[string]bool

	stripTypesList                 []string
	stripTypePrefixesList          []string
	stripNamePrefixesLowerCaseList []string
	stripNameSuffixesLowerCaseList []string
}

// Java: StripCode.STRIP_TYPE_INHERIT_ERROR
var StripCode_STRIP_TYPE_INHERIT_ERROR = DiagnosticType_Error(
	"JSC_STRIP_TYPE_INHERIT_ERROR", "Non-strip type {0} cannot inherit from strip type {1}")

// Java: StripCode.STRIP_ASSIGNMENT_ERROR
var StripCode_STRIP_ASSIGNMENT_ERROR = DiagnosticType_Error("JSC_STRIP_ASSIGNMENT_ERROR", "Unable to strip assignment to {0}")

// StripCode_ToStreamWithCollapsedVersions returns a stream containing `s`
// and, if `s` contains a ".", also all the possible collapsed forms of `s`.
// (e.g. for `'a.b.c'` we generate `'a$b.c'` and `'a$b$c'`)
//
// StripCode now runs after `CollapseProperties`, so it needs to look for
// both the original and collapsed versions of qualified names.
//
// Java: StripCode.toStreamWithCollapsedVersions(String)
func StripCode_ToStreamWithCollapsedVersions(s string) []string {
	possibleForms := []string{s}
	// Go-only: language-adapter — '.' and '$' are ASCII, so replacing them in
	// the WTF-8 bytes is the same as replacing them in Java's char array.
	chars := []byte(s)
	for i := 0; i < len(chars); i++ {
		if chars[i] == '.' {
			chars[i] = '$'
			possibleForms = append(possibleForms, string(chars))
		}
	}
	return possibleForms
}

// NewStripCode creates an instance.
//
// Java: StripCode(AbstractCompiler, ImmutableSet<String>, ImmutableSet<String>, ImmutableSet<String>, boolean)
func NewStripCode(
	compiler AbstractCompiler,
	stripTypes *jutil.LinkedHashSet[string],
	stripNameSuffixes *jutil.LinkedHashSet[string],
	stripNamePrefixes *jutil.LinkedHashSet[string],
	enableTweakStripping bool) *StripCode {
	s := &StripCode{varsToRemove: make(map[string]bool)}

	s.compiler = compiler

	// stripNameSuffixes.stream().flatMap(StripCode::toStreamWithCollapsedVersions)
	//     .collect(toImmutableSet())
	s.stripNameSuffixes = jutil.NewLinkedHashSet[string]()
	for n := range stripNameSuffixes.All() {
		s.stripNameSuffixes.AddAll(StripCode_ToStreamWithCollapsedVersions(n)...)
	}
	s.stripNamePrefixes = jutil.NewLinkedHashSet[string]()
	for n := range stripNamePrefixes.All() {
		s.stripNamePrefixes.AddAll(StripCode_ToStreamWithCollapsedVersions(n)...)
	}

	stripTypesStream := stripTypes.ToSlice()
	// Add "tweak" class stripping if requested
	if enableTweakStripping {
		stripTypesStream = append(stripTypesStream, "goog.tweak")
	}
	stripTypesAdjusted := jutil.NewLinkedHashSet[string]()
	for _, t := range stripTypesStream {
		stripTypesAdjusted.AddAll(StripCode_ToStreamWithCollapsedVersions(t)...)
	}

	// Iteration overhead was a high cost in this pass. Using a native array
	// is trivial and avoid those costs.
	s.stripTypesList = stripTypesAdjusted.ToSlice()

	// We want to strip types that are defined on a type that is being stripped, otherwise the
	// resulting code will be invalid. So, we'll also check for prefixes that indicate such child
	// names.
	// TODO(johnlenz): I'm not sure what the original intent of "type prefix" stripping was.
	// Verify that we can always assume a complete namespace and simplify this logic.
	s.stripTypePrefixesList = make([]string, 0, 2*stripTypesAdjusted.Size())
	for t := range stripTypesAdjusted.All() {
		// look for both non-collapsed and collapsed child names
		s.stripTypePrefixesList = append(s.stripTypePrefixesList, t+".", t+"$")
	}

	// Precalculate the lowercase versions of the string to avoid repeated
	// lowercase conversions.
	s.stripNamePrefixesLowerCaseList = make([]string, 0, s.stripNamePrefixes.Size())
	for p := range s.stripNamePrefixes.All() {
		s.stripNamePrefixesLowerCaseList = append(s.stripNamePrefixesLowerCaseList, jutil.StringToLowerCaseRoot(p))
	}

	s.stripNameSuffixesLowerCaseList = make([]string, 0, s.stripNameSuffixes.Size())
	for p := range s.stripNameSuffixes.All() {
		s.stripNameSuffixesLowerCaseList = append(s.stripNameSuffixesLowerCaseList, jutil.StringToLowerCaseRoot(p))
	}
	return s
}

// Java: StripCode.process(Node, Node)
func (s *StripCode) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(s.compiler.GetLifeCycleStage().IsNormalized())
	decisionsLog := s.compiler.CreateOrReopenIndexedLog(reflect.TypeFor[StripCode](), "decisions.log")
	defer decisionsLog.Close()
	decisionsLog.LogSupplier(NewStripCode_StripCodeConfigRecord(s).Get)
	decisionsLog.LogString("\n=== decisions ===\n")
	NodeTraversal_Traverse(s.compiler, root, NewStripCode_Strip(s, decisionsLog))
}

// -------------------------------------------------------------------------

// StripCode_StripCodeConfigRecord is the inner class StripCodeConfigRecord
// (a Supplier<String>).
type StripCode_StripCodeConfigRecord struct {
	this_ *StripCode // the enclosing instance
}

// Java: StripCodeConfigRecord()
func NewStripCode_StripCodeConfigRecord(outer *StripCode) *StripCode_StripCodeConfigRecord {
	return &StripCode_StripCodeConfigRecord{this_: outer}
}

// Java: StripCodeConfigRecord.get()
func (r *StripCode_StripCodeConfigRecord) Get() string {
	s := r.this_
	var builder strings.Builder
	builder.WriteString("=== stripNameSuffixes ===\n")
	for stripNameSuffix := range s.stripNameSuffixes.All() {
		builder.WriteString(stripNameSuffix)
		builder.WriteByte('\n')
	}
	builder.WriteByte('\n')
	builder.WriteString("=== stripNamePrefixes ===\n")
	for stripNamePrefix := range s.stripNamePrefixes.All() {
		builder.WriteString(stripNamePrefix)
		builder.WriteByte('\n')
	}
	builder.WriteByte('\n')
	builder.WriteString("=== stripTypesList ===\n")
	for _, stripType := range s.stripTypesList {
		builder.WriteString(stripType)
		builder.WriteByte('\n')
	}
	builder.WriteByte('\n')
	builder.WriteString("=== stripTypePrefixesList ===\n")
	for _, stripNamePrefix := range s.stripTypePrefixesList {
		builder.WriteString(stripNamePrefix)
		builder.WriteByte('\n')
	}
	builder.WriteByte('\n')
	return builder.String()
}

// StripCode_Strip is a callback that strips debug code from a JavaScript
// parse tree.
type StripCode_Strip struct {
	this_ *StripCode // the enclosing instance

	decisionsLog diagnostic.LogFile
}

// Java: Strip(LogFile)
func NewStripCode_Strip(outer *StripCode, decisionsLog diagnostic.LogFile) *StripCode_Strip {
	st := &StripCode_Strip{this_: outer}
	st.decisionsLog = decisionsLog
	return st
}

// Java: Strip.shouldTraverse(NodeTraversal, Node, Node)
func (st *StripCode_Strip) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Here we check for cases where we're going to remove a large chunk of code,
	// and so should not traverse into it to avoid both wasting time and causing
	// problems with logic duplication.
	switch n.GetToken() {
	case rhino.Token_CALL, rhino.Token_NEW:
		// If we're removing the whole call / new
		if st.IsMethodOrCtorCallThatTriggersRemoval(t, n, parent) {
			st.decisionsLog.LogSupplier(func() string { return "removing function call" })
			st.ReplaceHighestNestedCallWithNull(t, n, parent)
			return false
		} else {
			return true
		}
	default:
		return true
	}
}

// Java: Strip.visit(NodeTraversal, Node, Node)
func (st *StripCode_Strip) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		st.RemoveVarDeclarationsByNameOrRvalue(t, n, parent)
	case rhino.Token_NAME:
		st.MaybeRemoveReferenceToRemovedVariable(t, n, parent)
	case rhino.Token_ASSIGN,
		rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD:
		st.MaybeEliminateAssignmentByLvalueName(t, n, parent)
	case rhino.Token_OBJECTLIT:
		st.EliminateKeysWithStripNamesFromObjLit(n)
	case rhino.Token_EXPR_RESULT:
		st.MaybeEliminateExpressionByName(n)
	case rhino.Token_CLASS:
		st.MaybeEliminateClassByNameOrExtends(t, n, parent)
	default:
	}
}

// RemoveVarDeclarationsByNameOrRvalue removes declarations of any variables
// whose names are strip names or whose r-values are static method calls on
// strip types. Builds a set of removed variables so that all references to
// them can be removed.
//
// t is the traversal, n a VAR, CONST, or LET node, parent n's parent.
//
// Java: Strip.removeVarDeclarationsByNameOrRvalue(NodeTraversal, Node, Node)
func (st *StripCode_Strip) RemoveVarDeclarationsByNameOrRvalue(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	s := st.this_
	var next rhino.Node = rhino.NilNode
	for nameNode := n.GetFirstChild(); !nameNode.IsNil(); nameNode = next {
		next = nameNode.GetNext()
		if nameNode.IsDestructuringLhs() {
			continue
		}
		jutil.CheckStateMsg1(nameNode.IsName(), "%s", nameNode)
		name := nameNode.GetString()
		// If this variable represents a collapsed property, it's the original property name we're
		// supposed to be matching against.
		lastDollarSign := jutil.StringLastIndexOf(name, "$")
		possibleStripName := name
		if lastDollarSign != -1 {
			possibleStripName = jutil.StringSubstringFrom(name, lastDollarSign+1)
		}
		if st.IsStripName(possibleStripName) ||
			st.QualifiedNameBeginsWithStripType(nameNode) ||
			st.IsCallWhoseReturnValueShouldBeStripped(nameNode.GetFirstChild()) {
			// Remove the NAME.
			s.varsToRemove[name] = true
			if strings.Contains(name, "$") {
				// We need to be careful with this code pattern that appears after
				// collapsing properties.
				// ```javascript
				// /** @constructor */
				// var a$b$C = function() {
				//   this.nonStrippedName = a$b$C$strippedByNameOrValue;
				// };
				// var a$b$C$strippedByNameOrValue = strippedFunction();
				// ```
				// Note that the declaration of `a$b$C$nonStrippedByNameOrValue` will be visited
				// **after** its first use, so it's too late to go back and remove the
				// reference (without restructuring this class quite a bit). Instead,
				// we'll preserve here the behavior you would get before collapsing and
				// just replace the rhs with `null`.
				st.decisionsLog.LogSupplier(func() string { return name + ": initialize with null (" + possibleStripName + ")" })
				if nameNode.HasChildren() {
					st.ReplaceWithNull(nameNode.GetOnlyChild())
				} else {
					// `var my$name = null;` is a bit easier to optimize away than
					// `var my$name;`, because we can clearly see that it is initialized, so we have
					// a value to inline in later passes.
					nameNode.AddChildToFront(rhino.IR_NullNode().Srcref(nameNode))
				}
				t.ReportCodeChange()
			} else {
				// Assume that the declaration comes before any reference.
				// We will remove the references when we see them later.
				st.decisionsLog.LogSupplier(func() string { return name + ": removing declaration" })
				nameNode.Detach()
				NodeUtil_MarkFunctionsDeleted(nameNode, s.compiler)
			}
		}
	}
	if !n.HasChildren() {
		// Must also remove the VAR.
		st.ReplaceWithEmpty(n, parent)
		t.ReportCodeChange()
	}
}

// MaybeRemoveReferenceToRemovedVariable removes a reference if it is a
// reference to a removed variable.
//
// t is the traversal, n a NAME node, parent n's parent.
//
// Java: Strip.maybeRemoveReferenceToRemovedVariable(NodeTraversal, Node, Node)
func (st *StripCode_Strip) MaybeRemoveReferenceToRemovedVariable(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch parent.GetToken() {
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		// This is a variable declaration, not a reference.
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		// GETPROP
		//   NAME
		//   STRING (property name)
		// GETELEM
		//   NAME
		//   NUMBER|STRING|NAME|...
		if parent.GetFirstChild() == n && st.IsReferenceToRemovedVar(n) {
			st.decisionsLog.LogSupplier(func() string { return n.GetString() + ": removing getelem/getprop/call chain" })
			st.ReplaceHighestNestedCallWithNull(t, parent, parent.GetParent())
		}
	case rhino.Token_ASSIGN,
		rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD:
		if st.IsReferenceToRemovedVar(n) {
			if parent.GetFirstChild() == n {
				grandparent := parent.GetParent()
				st.decisionsLog.LogSupplier(
					func() string {
						return jutil.ToJavaString(n.GetQualifiedName()) + ": removing assignment to stripped var"
					})
				if grandparent.IsExprResult() {
					// Remove the assignment.
					greatGrandparent := grandparent.GetParent()
					st.ReplaceWithEmpty(grandparent, greatGrandparent)
					t.ReportCodeChange()
				} else {
					// Substitute the r-value for the assignment.
					rvalue := n.GetNext()
					rvalue.Detach()
					parent.ReplaceWith(rvalue)
					t.ReportCodeChange()
				}
			} else {
				// The var reference is the r-value. Replace it with null.
				st.decisionsLog.LogSupplier(
					func() string { return jutil.ToJavaString(n.GetQualifiedName()) + ": replacing rhs reference with null" })
				st.ReplaceWithNull(n)
				t.ReportCodeChange()
			}
		}
	case rhino.Token_NEW, rhino.Token_CALL:
		if !n.IsFirstChildOf(parent) && st.IsReferenceToRemovedVar(n) {
			// NOTE: the callee is handled when we visit the CALL or NEW node
			st.decisionsLog.LogSupplier(
				func() string {
					return jutil.ToJavaString(n.GetQualifiedName()) + ": replacing parameter reference with null"
				})
			st.ReplaceWithNull(n)
			t.ReportCodeChange()
		}
	case rhino.Token_COMMA:
		grandparent := parent.GetParent()
		// The last child in a comma expression is its result, so we need to be careful replacing
		// it with null. We don't want to replace the entire comma expression with null because
		// there could be other elements in it with side-effects.
		isLastChild := parent.GetLastChild() == n
		// The only parent where replacing with null is an issue is likely where the comma
		// expression is the first child (callee) of a CALL or NEW node.
		parentIsCallee :=
			(grandparent.IsCall() || grandparent.IsNew()) && parent.IsFirstChildOf(grandparent)
		isSafeToRemove := !isLastChild || !parentIsCallee
		if isSafeToRemove && st.IsReferenceToRemovedVar(n) {
			st.decisionsLog.LogSupplier(
				func() string {
					return jutil.ToJavaString(n.GetQualifiedName()) + ": replacing reference in comma expr with null"
				})
			st.ReplaceWithNull(n)
			t.ReportCodeChange()
		}
	default:
		if st.IsReferenceToRemovedVar(n) {
			st.decisionsLog.LogSupplier(
				func() string { return jutil.ToJavaString(n.GetQualifiedName()) + ": replacing reference with null" })
			st.ReplaceWithNull(n)
			t.ReportCodeChange()
		}
	}
}

// ReplaceHighestNestedCallWithNull uses a while loop to get up out of any
// nested calls. For example, if we have just detected that we need to remove
// the a.b() call in a.b().c().d(), we'll have to remove all of the calls,
// and it will take a few iterations through this loop to get up to d().
//
// Java: Strip.replaceHighestNestedCallWithNull(NodeTraversal, Node, Node)
func (st *StripCode_Strip) ReplaceHighestNestedCallWithNull(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	ancestor := parent
	ancestorChild := node
	var ancestorParent rhino.Node
	for {
		ancestorParent = ancestor.GetParent()
		if ancestorParent.IsNil() {
			return // the call was already removed from the AST
		}

		if ancestor.GetFirstChild() != ancestorChild {
			st.ReplaceWithNull(ancestorChild)
			break
		}
		if ancestor.IsExprResult() {
			// Remove the entire expression statement.
			st.ReplaceWithEmpty(ancestor, ancestorParent)
			break
		}
		if ancestor.IsAssign() {
			ancestor.ReplaceWith(ancestor.GetLastChild().Detach())
			break
		}
		if !NodeUtil_IsNormalGet(ancestor) && !ancestor.IsCall() {
			st.ReplaceWithNull(ancestorChild)
			break
		}

		// Is not executed on the last iteration so can't be used for change reporting.
		ancestorChild = ancestor
		ancestor = ancestorParent
	}
	t.ReportCodeChange()
}

// MaybeEliminateAssignmentByLvalueName eliminates an assignment if the
// l-value is:
//
//   - A field name that's a strip name
//   - A qualified name that begins with a strip type
//
// t is the traversal, n an ASSIGN node, parent n's parent.
//
// Java: Strip.maybeEliminateAssignmentByLvalueName(NodeTraversal, Node, Node)
func (st *StripCode_Strip) MaybeEliminateAssignmentByLvalueName(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	s := st.this_
	// ASSIGN
	//   l-value
	//   r-value
	lvalue := n.GetFirstChild()
	if st.NameIncludesFieldNameToStrip(lvalue) || st.QualifiedNameBeginsWithStripType(lvalue) {

		// Limit to EXPR_RESULT because it is not
		// safe to eliminate assignment in complex expressions,
		// e.g. in ((x = 7) + 8)
		if parent.IsExprResult() {
			st.decisionsLog.LogSupplier(func() string { return lvalue.GetString() + ": removing assignment statement" })
			grandparent := parent.GetParent()
			// the assignment may already have been removed when visiting either the lhs
			// or the rhs.
			if !grandparent.IsNil() {
				st.ReplaceWithEmpty(parent, grandparent)
				s.compiler.ReportChangeToEnclosingScope(grandparent)
			}
		} else {
			t.Report(n, StripCode_STRIP_ASSIGNMENT_ERROR, jutil.ToJavaString(lvalue.GetQualifiedName()))
		}
	}
}

// MaybeEliminateExpressionByName eliminates an expression if it refers to:
//
//   - A field name that's a strip name
//   - A qualified name that begins with a strip type
//
// This gets rid of construct like: a.prototype.logger; (used instead of
// a.prototype.logger = null;) This expression is not an assignment and so
// will not be caught by maybeEliminateAssignmentByLvalueName.
//
// n is an EXPR_RESULT node.
//
// Java: Strip.maybeEliminateExpressionByName(Node)
func (st *StripCode_Strip) MaybeEliminateExpressionByName(n rhino.Node) {
	s := st.this_
	// EXPR_RESULT
	//   expression
	jutil.CheckArgumentMsg1(n.IsExprResult(), "%s", n)
	parent := n.GetParent()
	if parent.IsNil() {
		// This EXPR_RESULT was already removed when one of its child nodes was visited.
		return
	}
	expression := n.GetFirstChild()
	if st.NameIncludesFieldNameToStrip(expression) ||
		st.QualifiedNameBeginsWithStripType(expression) {
		st.decisionsLog.LogSupplier(
			func() string { return expression.GetString() + ": removing property declaration statement" })
		st.ReplaceWithEmpty(n, parent)
		s.compiler.ReportChangeToEnclosingScope(parent)
	}
}

// EliminateKeysWithStripNamesFromObjLit eliminates any object literal keys
// in an object literal declaration that have strip names.
//
// n is an OBJLIT node.
//
// Java: Strip.eliminateKeysWithStripNamesFromObjLit(Node)
func (st *StripCode_Strip) EliminateKeysWithStripNamesFromObjLit(n rhino.Node) {
	s := st.this_
	// OBJLIT
	//   key1
	//     value1
	//   key2
	//   ...
	key := n.GetFirstChild()
	for !key.IsNil() {
		switch key.GetToken() {
		case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF:
			if st.IsStripName(key.GetString()) {
				next := key.GetNext()
				key.Detach()
				NodeUtil_MarkFunctionsDeleted(key, s.compiler)
				key = next
				s.compiler.ReportChangeToEnclosingScope(n)
				break
			}
			// fall through
			key = key.GetNext()
		default:
			key = key.GetNext()
		}
	}
}

// MaybeEliminateClassByNameOrExtends removes a class definition if the name
// is a strip type. Warns if a non-strippable class is extending a
// strippable type.
//
// Java: Strip.maybeEliminateClassByNameOrExtends(NodeTraversal, Node, Node)
func (st *StripCode_Strip) MaybeEliminateClassByNameOrExtends(t *NodeTraversal, classNode rhino.Node, parent rhino.Node) {
	nameNode := NodeUtil_GetNameNode(classNode)
	var className string
	// Replace class with null if it is a strip type
	if !nameNode.IsNil() && nameNode.IsQualifiedName() {
		className = *nameNode.GetQualifiedName()
		if st.QualifiedNameBeginsWithStripTypeName(&className) {
			st.decisionsLog.LogSupplier(func() string { return className + ": removing class" })
			if NodeUtil_IsStatementParent(parent) {
				st.ReplaceWithEmpty(classNode, parent)
			} else {
				st.ReplaceWithNull(classNode)
			}
			t.ReportCodeChange()
			return
		}
	} else {
		className = "<anonymous>"
	}

	// If the class is not a strip type, the superclass also cannot be a strip type
	superclassNode := classNode.GetSecondChild()
	if !superclassNode.IsNil() && superclassNode.IsQualifiedName() {
		superclassName := superclassNode.GetQualifiedName()
		if st.QualifiedNameBeginsWithStripTypeName(superclassName) {
			t.Report(classNode, StripCode_STRIP_TYPE_INHERIT_ERROR, className, *superclassName)
		}
	}
}

// IsCallWhoseReturnValueShouldBeStripped gets whether a node is a CALL node
// whose return value should be stripped. A call's return value should be
// stripped if the function getting called is a static method in a class
// that gets stripped. For example, if "goog.debug.Logger" is a strip name,
// then this function returns true for a call such as
// "goog.debug.Logger.getLogger(...)". It may also simply be a function that
// is getting stripped. For example, if "getLogger" is a strip name, but not
// "goog.debug.Logger", this will still return true.
//
// n is a node (typically a CALL node). Returns whether the call's return
// value should be stripped.
//
// Java: Strip.isCallWhoseReturnValueShouldBeStripped(Node)
func (st *StripCode_Strip) IsCallWhoseReturnValueShouldBeStripped(n rhino.Node) bool {
	if n.IsNil() || (!n.IsCall() && !n.IsNew()) || !n.HasChildren() {
		return false
	}

	function := NodeUtil_GetCallTargetResolvingIndirectCalls(n)
	return st.QualifiedNameBeginsWithStripType(function) || st.NameIncludesFieldNameToStrip(function)
}

// QualifiedNameBeginsWithStripType gets whether a qualified name begins
// with a strip name. The names "goog.debug", "goog.debug.Logger", and
// "goog.debug.Logger.Level" are examples of strip names that would result in
// this function returning true for a node representing the name
// "goog.debug.Logger.Level".
//
// n is a node (typically a NAME or GETPROP node). Returns whether the name
// begins with a strip name.
//
// Java: Strip.qualifiedNameBeginsWithStripType(Node)
func (st *StripCode_Strip) QualifiedNameBeginsWithStripType(n rhino.Node) bool {
	name := n.GetQualifiedName()
	return st.QualifiedNameBeginsWithStripTypeName(name)
}

// QualifiedNameBeginsWithStripTypeName gets whether a qualified name begins
// with a strip name. The names "goog.debug", "goog.debug.Logger", and
// "goog.debug.Logger.Level" are examples of strip names that would result in
// this function returning true for a node representing the name
// "goog.debug.Logger.Level".
//
// name is a qualified class name (nil is Java null). Returns whether the name
// begins with a strip name.
//
// Java: Strip.qualifiedNameBeginsWithStripType(String)
func (st *StripCode_Strip) QualifiedNameBeginsWithStripTypeName(name *string) bool {
	s := st.this_
	if name != nil {
		for _, type_ := range s.stripTypesList {
			if *name == type_ {
				if st.decisionsLog.IsLogging() {
					st.LogStripName(jutil.ToJavaString(name), "equals strip type")
				}
				return true
			}
		}
		for _, type_ := range s.stripTypePrefixesList {
			if jutil.StringStartsWith(*name, type_) {
				if st.decisionsLog.IsLogging() {
					st.LogStripName(jutil.ToJavaString(name), "starts with strip type prefix")
				}
				return true
			}
		}
	}
	if st.decisionsLog.IsLogging() {
		st.LogNotAStripName(jutil.ToJavaString(name), "does not begin with a strip type")
	}
	return false
}

// IsReferenceToRemovedVar determines whether a NAME node represents a
// reference to a variable that has been removed.
//
// n is a NAME node. Returns whether the variable was removed.
//
// Java: Strip.isReferenceToRemovedVar(Node)
func (st *StripCode_Strip) IsReferenceToRemovedVar(n rhino.Node) bool {
	return st.this_.varsToRemove[n.GetString()]
}

// IsMethodOrCtorCallThatTriggersRemoval gets whether a CALL node triggers
// statement removal, based on the name of the object whose method is being
// called, or the name of the method. Checks whether the name begins with a
// strip type, includes a field name that's a strip name, or belongs to the
// set of global class-defining functions (e.g. goog.inherits).
//
// t is the traversal, n a CALL or NEW node. Returns whether the node triggers
// statement removal.
//
// Java: Strip.isMethodOrCtorCallThatTriggersRemoval(NodeTraversal, Node, Node)
func (st *StripCode_Strip) IsMethodOrCtorCallThatTriggersRemoval(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// CALL/NEW
	//   GETPROP (function)         <-- we're interested in this, the function
	//     GETPROP (callee object)  <-- or the object on which it is called
	//       ...
	//       STRING (field name)
	//     STRING (method name)
	//   ... (arguments)

	function := NodeUtil_GetCallTargetResolvingIndirectCalls(n)
	if function.IsNil() || !function.IsQualifiedName() {
		return false
	}

	if !parent.IsNil() && parent.IsName() {
		grandparent := parent.GetParent()
		if !grandparent.IsNil() && NodeUtil_IsNameDeclaration(grandparent) {
			// The call's return value is being used to initialize a newly
			// declared variable. We should leave the call intact for now.
			// That way, when the traversal reaches the variable declaration,
			// we'll recognize that the variable and all references to it need
			// to be eliminated.
			return false
		}
	}

	if function.IsName() && st.IsStripName(function.GetString()) {
		return true
	}
	callee := function.GetFirstChild()
	return st.NameIncludesFieldNameToStrip(callee) ||
		st.NameIncludesFieldNameToStrip(function) ||
		st.QualifiedNameBeginsWithStripType(function) ||
		st.ActsOnStripType(t, n)
}

// NameIncludesFieldNameToStrip returns whether a name includes a field name
// that should be stripped. E.g., "foo.stripMe.bar", "(foo.bar).stripMe",
// etc.
//
// Java: Strip.nameIncludesFieldNameToStrip(Node)
func (st *StripCode_Strip) NameIncludesFieldNameToStrip(n rhino.Node) bool {
	if n.IsNil() {
		return false
	} else if n.IsGetProp() {
		return st.IsStripName(n.GetString()) || st.NameIncludesFieldNameToStrip(n.GetFirstChild())
	} else if n.IsName() {
		nameString := n.GetString()
		// CollapseProperties may have turned "a.b.c" into "a$b$c",
		// so split that up and match its parts.
		end := jutil.StringIndexOfChar(nameString, '$')
		if end == -1 {
			return false
		}
		// iterate over all the $-deliminated parts of the name, where each part is from [start,end)
		start := 0
		for {
			if st.IsStripName(jutil.StringSubstring(nameString, start, end)) {
				return true
			}
			if end == jutil.StringLength(nameString) {
				break
			}
			newStart := end + 1
			newEnd := jutil.StringIndexOfCharFrom(nameString, '$', newStart)
			start = newStart
			if newEnd == -1 {
				end = jutil.StringLength(nameString)
			} else {
				end = newEnd
			}
		}
	}
	return false
}

// ActsOnStripType determines whether the given node helps to define a strip
// type. For example, goog.inherits(stripType, Object) would be such a call.
//
// Also reports an error if a non-strip type inherits from a strip type.
//
// t is the current traversal, callNode the CALL node.
//
// Java: Strip.actsOnStripType(NodeTraversal, Node)
func (st *StripCode_Strip) ActsOnStripType(t *NodeTraversal, callNode rhino.Node) bool {
	s := st.this_
	classes := s.compiler.GetCodingConvention().GetClassesDefinedByCall(callNode)
	if classes != nil {
		// It's okay to strip a type that inherits from a non-stripped type
		// e.g. goog.inherits(goog.debug.Logger, Object)
		subclassName := classes.subclassName
		if st.QualifiedNameBeginsWithStripTypeName(&subclassName) {
			st.LogStripName(classes.subclassName, "class defining call")
			return true
		}

		// report an error if a non-strip type inherits from a
		// strip type.
		superclassName := classes.superclassName
		if st.QualifiedNameBeginsWithStripTypeName(&superclassName) {
			t.Report(
				callNode, StripCode_STRIP_TYPE_INHERIT_ERROR, classes.subclassName, classes.superclassName)
		}
	}

	return false
}

// IsStripName gets whether a JavaScript identifier is the name of a
// variable or property that should be stripped.
//
// name is a JavaScript identifier. Returns whether name is a name that
// triggers removal.
//
// Java: Strip.isStripName(String)
func (st *StripCode_Strip) IsStripName(name string) bool {
	s := st.this_
	if !s.stripNameSuffixes.IsEmpty() && s.stripNameSuffixes.Contains(name) {
		st.LogStripName(name, "matches a suffix")
		return true
	}

	if !s.stripNamePrefixes.IsEmpty() && s.stripNamePrefixes.Contains(name) {
		st.LogStripName(name, "matches a prefix")
		return true
	}

	if name == "" || jutil.CharacterIsUpperCase(jutil.StringCharAt(name, 0)) {
		st.LogNotAStripName(name, "empty or starts with uppercase")
		return false
	}

	if len(s.stripNamePrefixesLowerCaseList) > 0 || len(s.stripNameSuffixesLowerCaseList) > 0 {
		lcName := jutil.StringToLowerCaseRoot(name)
		for _, stripName := range s.stripNamePrefixesLowerCaseList {
			if jutil.StringStartsWith(lcName, stripName) {
				if st.decisionsLog.IsLogging() {
					st.LogStripNameSupplier(name, func() string { return "matches lc prefix: " + stripName })
				}
				return true
			}
		}

		for _, stripName := range s.stripNameSuffixesLowerCaseList {
			if jutil.StringEndsWith(lcName, stripName) {
				if st.decisionsLog.IsLogging() {
					st.LogStripNameSupplier(name, func() string { return "matches lc suffix: " + stripName })
				}
				return true
			}
		}
	}

	st.LogNotAStripName(name, "no matches")
	return false
}

// Java: Strip.logNotAStripName(String, String)
func (st *StripCode_Strip) LogNotAStripName(name string, reason string) {
	if st.decisionsLog.IsLogging() {
		st.decisionsLog.LogSupplier(func() string { return name + "\tnot a strip name: " + reason })
	}
}

// Java: Strip.logStripName(String, String)
func (st *StripCode_Strip) LogStripName(name string, reason string) {
	if st.decisionsLog.IsLogging() {
		st.decisionsLog.LogSupplier(func() string { return name + "\tstrip name: " + reason })
	}
}

// Java: Strip.logStripName(String, Supplier<String>)
func (st *StripCode_Strip) LogStripNameSupplier(name string, reasonSupplier func() string) {
	if st.decisionsLog.IsLogging() {
		st.decisionsLog.LogSupplier(func() string { return name + "\tstrip name: " + reasonSupplier() })
	}
}

// ReplaceWithNull replaces a node with a NULL node. This is useful where a
// value is expected.
//
// Java: Strip.replaceWithNull(Node)
func (st *StripCode_Strip) ReplaceWithNull(n rhino.Node) {
	if st.decisionsLog.IsLogging() {
		st.decisionsLog.LogSupplier(func() string { return "replace with null: " + n.GetLocation() })
	}
	n.ReplaceWith(rhino.IR_NullNode())
	NodeUtil_MarkFunctionsDeleted(n, st.this_.compiler)
}

// ReplaceWithEmpty replaces a node with an EMPTY node. This is useful where
// a statement is expected.
//
// n is a node, parent n's parent.
//
// Java: Strip.replaceWithEmpty(Node, Node)
func (st *StripCode_Strip) ReplaceWithEmpty(n rhino.Node, parent rhino.Node) {
	if st.decisionsLog.IsLogging() {
		st.decisionsLog.LogSupplier(func() string { return "replace with empty: " + n.GetLocation() })
	}
	NodeUtil_RemoveChild(parent, n)
	NodeUtil_MarkFunctionsDeleted(n, st.this_.compiler)
}
