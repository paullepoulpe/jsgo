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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
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
	compiler              AbstractCompiler
	stripNameSuffixesSet  map[string]struct{}
	stripNameSuffixesList []string
	stripNamePrefixesSet  map[string]struct{}
	stripNamePrefixesList []string
	varsToRemove          map[uint32]struct{}

	stripTypesSet                  map[string]struct{}
	stripTypesList                 []string
	stripTypePrefixesList          []string
	stripNamePrefixesLowerCaseList []string
	stripNameSuffixesLowerCaseList []string
}

var _ CompilerPass = (*StripCode)(nil)

// Java: StripCode.STRIP_TYPE_INHERIT_ERROR
var StripCode_STRIP_TYPE_INHERIT_ERROR = DiagnosticType_Error(
	"JSC_STRIP_TYPE_INHERIT_ERROR", "Non-strip type {0} cannot inherit from strip type {1}")

// Java: StripCode.STRIP_ASSIGNMENT_ERROR
var StripCode_STRIP_ASSIGNMENT_ERROR = DiagnosticType_Error("JSC_STRIP_ASSIGNMENT_ERROR", "Unable to strip assignment to {0}")

// StripCode_ToStreamWithCollapsedVersions returns a slice containing `s`
// and, if `s` contains a ".", also all the possible collapsed forms of `s`.
// (e.g. for `'a.b.c'` we generate `'a$b.c'` and `'a$b$c'`)
//
// StripCode now runs after `CollapseProperties`, so it needs to look for
// both the original and collapsed versions of qualified names.
//
// Java: StripCode.toStreamWithCollapsedVersions(String)
func StripCode_ToStreamWithCollapsedVersions(s string) []string {
	if strings.IndexByte(s, '.') == -1 {
		return []string{s}
	}
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

func stripCode_expandCollapsed(src *jutil.LinkedHashSet[string], extra string) ([]string, map[string]struct{}) {
	n := src.Size()
	if extra != "" {
		n++
	}
	if n == 0 {
		return nil, nil
	}
	set := make(map[string]struct{}, n*2)
	list := make([]string, 0, n*2)
	addCollapsed := func(item string) {
		for _, form := range StripCode_ToStreamWithCollapsedVersions(item) {
			if _, exists := set[form]; !exists {
				set[form] = struct{}{}
				list = append(list, form)
			}
		}
	}
	for item := range src.All() {
		addCollapsed(item)
	}
	if extra != "" {
		addCollapsed(extra)
	}
	return list, set
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
	s := &StripCode{compiler: compiler}

	s.stripNameSuffixesList, s.stripNameSuffixesSet = stripCode_expandCollapsed(stripNameSuffixes, "")
	s.stripNamePrefixesList, s.stripNamePrefixesSet = stripCode_expandCollapsed(stripNamePrefixes, "")

	extraType := ""
	if enableTweakStripping {
		extraType = "goog.tweak"
	}
	s.stripTypesList, s.stripTypesSet = stripCode_expandCollapsed(stripTypes, extraType)

	// We want to strip types that are defined on a type that is being stripped, otherwise the
	// resulting code will be invalid. So, we'll also check for prefixes that indicate such child
	// names.
	// TODO(johnlenz): I'm not sure what the original intent of "type prefix" stripping was.
	// Verify that we can always assume a complete namespace and simplify this logic.
	if len(s.stripTypesList) > 0 {
		s.stripTypePrefixesList = make([]string, 0, 2*len(s.stripTypesList))
		for _, t := range s.stripTypesList {
			// look for both non-collapsed and collapsed child names
			s.stripTypePrefixesList = append(s.stripTypePrefixesList, t+".", t+"$")
		}
	}

	// Precalculate the lowercase versions of the string to avoid repeated
	// lowercase conversions.
	if len(s.stripNamePrefixesList) > 0 {
		s.stripNamePrefixesLowerCaseList = make([]string, len(s.stripNamePrefixesList))
		for i, p := range s.stripNamePrefixesList {
			s.stripNamePrefixesLowerCaseList[i] = jutil.StringToLowerCaseRoot(p)
		}
	}

	if len(s.stripNameSuffixesList) > 0 {
		s.stripNameSuffixesLowerCaseList = make([]string, len(s.stripNameSuffixesList))
		for i, p := range s.stripNameSuffixesList {
			s.stripNameSuffixesLowerCaseList[i] = jutil.StringToLowerCaseRoot(p)
		}
	}
	return s
}

// Java: StripCode.process(Node, Node)
func (s *StripCode) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(s.compiler.GetLifeCycleStage().IsNormalized())
	decisionsLog := s.compiler.CreateOrReopenIndexedLog("StripCode", "decisions.log")
	defer decisionsLog.Close()
	if decisionsLog.IsLogging() {
		decisionsLog.LogString(s.configRecord())
	}
	decisionsLog.LogString("\n=== decisions ===\n")
	NodeTraversal_Traverse(s.compiler, root, NewStripCode_Strip(s, decisionsLog))
}

// configRecord formats the StripCode configuration summary for decisions.log
// (flattened from Java inner class StripCode.StripCodeConfigRecord).
//
// Java: StripCodeConfigRecord.get()
func (s *StripCode) configRecord() string {
	var builder strings.Builder
	builder.WriteString("=== stripNameSuffixes ===\n")
	for _, stripNameSuffix := range s.stripNameSuffixesList {
		builder.WriteString(stripNameSuffix)
		builder.WriteByte('\n')
	}
	builder.WriteByte('\n')
	builder.WriteString("=== stripNamePrefixes ===\n")
	for _, stripNamePrefix := range s.stripNamePrefixesList {
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

func stripCode_qualifiedNameOrNull(n rhino.Node) string {
	if qname, ok := n.GetQualifiedNameString(); ok {
		return qname
	}
	return "null"
}

var stripCodeVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_VAR)
	s.Add(rhino.Token_CONST)
	s.Add(rhino.Token_LET)
	s.Add(rhino.Token_NAME)
	s.Add(rhino.Token_ASSIGN)
	s.Add(rhino.Token_ASSIGN_BITOR)
	s.Add(rhino.Token_ASSIGN_BITXOR)
	s.Add(rhino.Token_ASSIGN_BITAND)
	s.Add(rhino.Token_ASSIGN_LSH)
	s.Add(rhino.Token_ASSIGN_RSH)
	s.Add(rhino.Token_ASSIGN_URSH)
	s.Add(rhino.Token_ASSIGN_ADD)
	s.Add(rhino.Token_ASSIGN_SUB)
	s.Add(rhino.Token_ASSIGN_MUL)
	s.Add(rhino.Token_ASSIGN_DIV)
	s.Add(rhino.Token_ASSIGN_MOD)
	s.Add(rhino.Token_OBJECTLIT)
	s.Add(rhino.Token_EXPR_RESULT)
	s.Add(rhino.Token_CLASS)
	return s
}()

var _ nodeTraversal_visitTokens = (*StripCode_Strip)(nil)

// StripCode_Strip is a callback that strips debug code from a JavaScript
// parse tree.
type StripCode_Strip struct {
	this_        *StripCode
	decisionsLog diagnostic.LogFile
}

func (*StripCode_Strip) visitTokens() *ast.TokenSet {
	return &stripCodeVisitTokens
}

// Java: Strip(LogFile)
func NewStripCode_Strip(outer *StripCode, decisionsLog diagnostic.LogFile) *StripCode_Strip {
	return &StripCode_Strip{
		this_:        outer,
		decisionsLog: decisionsLog,
	}
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
			if st.decisionsLog.IsLogging() {
				st.decisionsLog.LogString("removing function call")
			}
			st.ReplaceHighestNestedCallWithNull(t, n, parent)
			return false
		}
		return true
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
	decl := ast.MustNameDeclaration(n)
	for declarator := range ast.ReadAhead(decl.FirstDeclarator()) {
		nameDecl, ok := declarator.AsName()
		if !ok {
			// Destructuring LHS — StripCode does not optimize destructuring patterns.
			continue
		}
		nameNode := nameDecl.Node()
		name := nameDecl.Ident()
		// If this variable represents a collapsed property, it's the original property name we're
		// supposed to be matching against. (ASCII '$' never occurs inside multi-byte UTF-8/WTF-8.)
		lastDollarSign := strings.LastIndexByte(name, '$')
		possibleStripName := name
		if lastDollarSign != -1 {
			possibleStripName = name[lastDollarSign+1:]
		}
		initExpr, hasInit := nameDecl.Init()
		if st.IsStripName(possibleStripName) ||
			st.QualifiedNameBeginsWithStripType(nameNode) ||
			(hasInit && st.IsCallWhoseReturnValueShouldBeStripped(initExpr.Node())) {
			// Remove the NAME.
			if s.varsToRemove == nil {
				s.varsToRemove = make(map[uint32]struct{})
			}
			s.varsToRemove[nameNode.GetStringID()] = struct{}{}
			if lastDollarSign != -1 {
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
				if st.decisionsLog.IsLogging() {
					st.decisionsLog.LogString(name + ": initialize with null (" + possibleStripName + ")")
				}
				if hasInit {
					st.ReplaceWithNull(initExpr.Node())
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
				if st.decisionsLog.IsLogging() {
					st.decisionsLog.LogString(name + ": removing declaration")
				}
				nameDecl.Detach()
				NodeUtil_MarkFunctionsDeleted(nameNode, s.compiler)
			}
		}
	}
	if !decl.HasDeclarators() {
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
			if st.decisionsLog.IsLogging() {
				st.decisionsLog.LogString(n.GetString() + ": removing getelem/getprop/call chain")
			}
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
			assign := ast.MustAssignment(parent)
			if assign.Target().Is(n) {
				grandparent := parent.GetParent()
				if st.decisionsLog.IsLogging() {
					st.decisionsLog.LogString(stripCode_qualifiedNameOrNull(n) + ": removing assignment to stripped var")
				}
				if grandparent.IsExprResult() {
					// Remove the assignment.
					greatGrandparent := grandparent.GetParent()
					st.ReplaceWithEmpty(grandparent, greatGrandparent)
					t.ReportCodeChange()
				} else {
					// Substitute the r-value for the assignment.
					rvalue := assign.Value().Node()
					rvalue.Detach()
					parent.ReplaceWith(rvalue)
					t.ReportCodeChange()
				}
			} else {
				// The var reference is the r-value. Replace it with null.
				if st.decisionsLog.IsLogging() {
					st.decisionsLog.LogString(stripCode_qualifiedNameOrNull(n) + ": replacing rhs reference with null")
				}
				st.ReplaceWithNull(n)
				t.ReportCodeChange()
			}
		}
	case rhino.Token_NEW, rhino.Token_CALL:
		if !n.IsFirstChildOf(parent) && st.IsReferenceToRemovedVar(n) {
			// NOTE: the callee is handled when we visit the CALL or NEW node
			if st.decisionsLog.IsLogging() {
				st.decisionsLog.LogString(stripCode_qualifiedNameOrNull(n) + ": replacing parameter reference with null")
			}
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
			if st.decisionsLog.IsLogging() {
				st.decisionsLog.LogString(stripCode_qualifiedNameOrNull(n) + ": replacing reference in comma expr with null")
			}
			st.ReplaceWithNull(n)
			t.ReportCodeChange()
		}
	default:
		if st.IsReferenceToRemovedVar(n) {
			if st.decisionsLog.IsLogging() {
				st.decisionsLog.LogString(stripCode_qualifiedNameOrNull(n) + ": replacing reference with null")
			}
			st.ReplaceWithNull(n)
			t.ReportCodeChange()
		}
	}
}

// ReplaceHighestNestedCallWithNull uses a loop to get up out of any
// nested calls. For example, if we have just detected that we need to remove
// the a.b() call in a.b().c().d(), we'll have to remove all of the calls,
// and it will take a few iterations through this loop to get up to d().
//
// Java: Strip.replaceHighestNestedCallWithNull(NodeTraversal, Node, Node)
func (st *StripCode_Strip) ReplaceHighestNestedCallWithNull(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	ancestor := parent
	ancestorChild := node
	for {
		ancestorParent := ancestor.GetParent()
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
		if assign, ok := ast.AsAssign(ancestor); ok {
			ancestor.ReplaceWith(assign.Value().Node().Detach())
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
	lvalue := ast.MustAssignment(n).Target().Node()
	if st.NameIncludesFieldNameToStrip(lvalue) || st.QualifiedNameBeginsWithStripType(lvalue) {
		// Limit to EXPR_RESULT because it is not
		// safe to eliminate assignment in complex expressions,
		// e.g. in ((x = 7) + 8)
		if parent.IsExprResult() {
			if st.decisionsLog.IsLogging() {
				st.decisionsLog.LogString(lvalue.GetString() + ": removing assignment statement")
			}
			grandparent := parent.GetParent()
			// the assignment may already have been removed when visiting either the lhs
			// or the rhs.
			if !grandparent.IsNil() {
				st.ReplaceWithEmpty(parent, grandparent)
				s.compiler.ReportChangeToEnclosingScope(grandparent)
			}
		} else {
			t.Report(n, StripCode_STRIP_ASSIGNMENT_ERROR, stripCode_qualifiedNameOrNull(lvalue))
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
	exprStmt, ok := ast.AsExprStmt(n)
	jutil.CheckArgumentMsg1(ok, "%s", n)
	parent := n.GetParent()
	if parent.IsNil() {
		// This EXPR_RESULT was already removed when one of its child nodes was visited.
		return
	}
	expression := exprStmt.Expr().Node()
	if st.NameIncludesFieldNameToStrip(expression) ||
		st.QualifiedNameBeginsWithStripType(expression) {
		if st.decisionsLog.IsLogging() {
			st.decisionsLog.LogString(expression.GetString() + ": removing property declaration statement")
		}
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
	for member := range ast.ReadAhead(ast.MustObjectLit(n).FirstMember()) {
		if key, ok := member.AsPlainKeyMember(); ok && st.IsStripName(key.Key()) {
			key.Detach()
			NodeUtil_MarkFunctionsDeleted(key.Node(), s.compiler)
			s.compiler.ReportChangeToEnclosingScope(n)
		}
	}
}

// MaybeEliminateClassByNameOrExtends removes a class definition if the name
// is a strip type. Warns if a non-strippable class is extending a
// strippable type.
//
// Java: Strip.maybeEliminateClassByNameOrExtends(NodeTraversal, Node, Node)
func (st *StripCode_Strip) MaybeEliminateClassByNameOrExtends(t *NodeTraversal, classNode rhino.Node, parent rhino.Node) {
	cls := ast.MustClass(classNode)
	nameNode := NodeUtil_GetNameNode(classNode)
	className := "<anonymous>"
	// Replace class with null if it is a strip type
	if qname, ok := nameNode.GetQualifiedNameString(); ok {
		className = qname
		if st.qualifiedNameBeginsWithStripType(className, true) {
			if st.decisionsLog.IsLogging() {
				st.decisionsLog.LogString(className + ": removing class")
			}
			if NodeUtil_IsStatementParent(parent) {
				st.ReplaceWithEmpty(classNode, parent)
			} else {
				st.ReplaceWithNull(classNode)
			}
			t.ReportCodeChange()
			return
		}
	}

	// If the class is not a strip type, the superclass also cannot be a strip type
	if superclass, ok := cls.Extends(); ok {
		if superclassName, ok := superclass.Node().GetQualifiedNameString(); ok &&
			st.qualifiedNameBeginsWithStripType(superclassName, true) {
			t.Report(classNode, StripCode_STRIP_TYPE_INHERIT_ERROR, className, superclassName)
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
	name, ok := n.GetQualifiedNameString()
	return st.qualifiedNameBeginsWithStripType(name, ok)
}

func (st *StripCode_Strip) qualifiedNameBeginsWithStripType(name string, hasName bool) bool {
	s := st.this_
	if hasName {
		if _, ok := s.stripTypesSet[name]; ok {
			if st.decisionsLog.IsLogging() {
				st.LogStripName(name, "equals strip type")
			}
			return true
		}
		for _, prefix := range s.stripTypePrefixesList {
			if strings.HasPrefix(name, prefix) {
				if st.decisionsLog.IsLogging() {
					st.LogStripName(name, "starts with strip type prefix")
				}
				return true
			}
		}
		if st.decisionsLog.IsLogging() {
			st.LogNotAStripName(name, "does not begin with a strip type")
		}
		return false
	}
	if st.decisionsLog.IsLogging() {
		st.LogNotAStripName("null", "does not begin with a strip type")
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
	_, ok := st.this_.varsToRemove[n.GetStringID()]
	return ok
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

	if fnName, ok := ast.AsName(function); ok && st.IsStripName(fnName.Ident()) {
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
	for cur := n; !cur.IsNil(); {
		if prop, ok := ast.AsGetProp(cur); ok {
			if st.IsStripName(prop.Property()) {
				return true
			}
			cur = prop.Object().Node()
			continue
		}
		if name, ok := ast.AsName(cur); ok {
			nameString := name.Ident()
			// CollapseProperties may have turned "a.b.c" into "a$b$c",
			// so split that up and match its parts. (ASCII '$' never occurs inside multi-byte UTF-8/WTF-8.)
			end := strings.IndexByte(nameString, '$')
			if end == -1 {
				return false
			}
			for {
				if st.IsStripName(nameString[:end]) {
					return true
				}
				nameString = nameString[end+1:]
				end = strings.IndexByte(nameString, '$')
				if end == -1 {
					return st.IsStripName(nameString)
				}
			}
		}
		return false
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
		if st.qualifiedNameBeginsWithStripType(classes.subclassName, true) {
			st.LogStripName(classes.subclassName, "class defining call")
			return true
		}

		// report an error if a non-strip type inherits from a
		// strip type.
		if st.qualifiedNameBeginsWithStripType(classes.superclassName, true) {
			t.Report(
				callNode, StripCode_STRIP_TYPE_INHERIT_ERROR, classes.subclassName, classes.superclassName)
		}
	}

	return false
}

func stripCode_startsWithUpperCase(name string) bool {
	c := name[0]
	if c < 0x80 {
		return c >= 'A' && c <= 'Z'
	}
	return jutil.CharacterIsUpperCase(jutil.StringCharAt(name, 0))
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
	if _, ok := s.stripNameSuffixesSet[name]; ok {
		st.LogStripName(name, "matches a suffix")
		return true
	}

	if _, ok := s.stripNamePrefixesSet[name]; ok {
		st.LogStripName(name, "matches a prefix")
		return true
	}

	if name == "" || stripCode_startsWithUpperCase(name) {
		st.LogNotAStripName(name, "empty or starts with uppercase")
		return false
	}

	if len(s.stripNamePrefixesLowerCaseList) > 0 || len(s.stripNameSuffixesLowerCaseList) > 0 {
		lcName := jutil.StringToLowerCaseRoot(name)
		for _, stripName := range s.stripNamePrefixesLowerCaseList {
			if strings.HasPrefix(lcName, stripName) {
				if st.decisionsLog.IsLogging() {
					st.LogStripName(name, "matches lc prefix: "+stripName)
				}
				return true
			}
		}

		for _, stripName := range s.stripNameSuffixesLowerCaseList {
			if strings.HasSuffix(lcName, stripName) {
				if st.decisionsLog.IsLogging() {
					st.LogStripName(name, "matches lc suffix: "+stripName)
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
		st.decisionsLog.LogString(name + "\tnot a strip name: " + reason)
	}
}

// Java: Strip.logStripName(String, String)
func (st *StripCode_Strip) LogStripName(name string, reason string) {
	if st.decisionsLog.IsLogging() {
		st.decisionsLog.LogString(name + "\tstrip name: " + reason)
	}
}

// ReplaceWithNull replaces a node with a NULL node. This is useful where a
// value is expected.
//
// Java: Strip.replaceWithNull(Node)
func (st *StripCode_Strip) ReplaceWithNull(n rhino.Node) {
	if st.decisionsLog.IsLogging() {
		st.decisionsLog.LogString("replace with null: " + n.GetLocation())
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
		st.decisionsLog.LogString("replace with empty: " + n.GetLocation())
	}
	NodeUtil_RemoveChild(parent, n)
	NodeUtil_MarkFunctionsDeleted(n, st.this_.compiler)
}
