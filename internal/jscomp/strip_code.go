// Port of com/google/javascript/jscomp/StripCode.java @ closure-compiler v20261028

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
	"iter"
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
	compiler                       AbstractCompiler
	decisionsLog                   diagnostic.LogFile
	stripNameSuffixes              OrderedSet[string]
	stripNamePrefixes              OrderedSet[string]
	stripTypes                     OrderedSet[string]
	varsToRemove                   map[uint32]struct{}
	stripTypePrefixesList          []string
	stripNamePrefixesLowerCaseList []string
	stripNameSuffixesLowerCaseList []string
}

var (
	_ CompilerPass              = (*StripCode)(nil)
	_ NodeTraversal_Callback    = (*StripCode)(nil)
	_ nodeTraversal_visitTokens = (*StripCode)(nil)
)

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

func expandCollapsedInto(dst *OrderedSet[string], src jutil.Set[string], extra string) {
	var capHint int
	var seq iter.Seq[string]
	if src != nil && !src.IsEmpty() {
		capHint = src.Size()
		seq = src.All()
	}
	if extra != "" {
		capHint++
	}
	if capHint == 0 {
		return
	}
	*dst = *NewOrderedSetWithCapacity[string](capHint * 2)
	addCollapsed := func(item string) {
		for _, form := range StripCode_ToStreamWithCollapsedVersions(item) {
			dst.Add(form)
		}
	}
	if seq != nil {
		for item := range seq {
			addCollapsed(item)
		}
	}
	if extra != "" {
		addCollapsed(extra)
	}
}

func (s *StripCode) initDerivedLists() {
	// We want to strip types that are defined on a type that is being stripped, otherwise the
	// resulting code will be invalid. So, we'll also check for prefixes that indicate such child
	// names.
	// TODO(johnlenz): I'm not sure what the original intent of "type prefix" stripping was.
	// Verify that we can always assume a complete namespace and simplify this logic.
	typeEntries := s.stripTypes.Entries()
	if len(typeEntries) > 0 {
		s.stripTypePrefixesList = make([]string, 0, 2*len(typeEntries))
		for _, e := range typeEntries {
			// look for both non-collapsed and collapsed child names
			s.stripTypePrefixesList = append(s.stripTypePrefixesList, e.Key+".", e.Key+"$")
		}
	}

	// Precalculate the lowercase versions of the string to avoid repeated
	// lowercase conversions.
	prefixEntries := s.stripNamePrefixes.Entries()
	if len(prefixEntries) > 0 {
		s.stripNamePrefixesLowerCaseList = make([]string, len(prefixEntries))
		for i, e := range prefixEntries {
			s.stripNamePrefixesLowerCaseList[i] = jutil.StringToLowerCaseRoot(e.Key)
		}
	}

	suffixEntries := s.stripNameSuffixes.Entries()
	if len(suffixEntries) > 0 {
		s.stripNameSuffixesLowerCaseList = make([]string, len(suffixEntries))
		for i, e := range suffixEntries {
			s.stripNameSuffixesLowerCaseList[i] = jutil.StringToLowerCaseRoot(e.Key)
		}
	}
}

// NewStripCode creates an instance.
//
// Java: StripCode(AbstractCompiler, ImmutableSet<String>, ImmutableSet<String>, ImmutableSet<String>, boolean)
func NewStripCode(
	compiler AbstractCompiler,
	stripTypes jutil.Set[string],
	stripNameSuffixes jutil.Set[string],
	stripNamePrefixes jutil.Set[string],
	enableTweakStripping bool) *StripCode {
	s := &StripCode{compiler: compiler}
	expandCollapsedInto(&s.stripNameSuffixes, stripNameSuffixes, "")
	expandCollapsedInto(&s.stripNamePrefixes, stripNamePrefixes, "")
	extraType := ""
	if enableTweakStripping {
		extraType = "goog.tweak"
	}
	expandCollapsedInto(&s.stripTypes, stripTypes, extraType)
	s.initDerivedLists()
	return s
}

// Java: StripCode.process(Node, Node)
func (s *StripCode) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(s.compiler.GetLifeCycleStage().IsNormalized())
	decisionsLog := s.compiler.CreateOrReopenIndexedLog("StripCode", "decisions.log")
	s.decisionsLog = decisionsLog
	defer func() {
		s.decisionsLog = diagnostic.LogFile{}
		decisionsLog.Close()
	}()
	if decisionsLog.IsLogging() {
		decisionsLog.LogString(s.configRecord())
	}
	decisionsLog.LogString("\n=== decisions ===\n")
	NodeTraversal_Traverse(s.compiler, root, s)
}

// configRecord formats the StripCode configuration summary for decisions.log
// (flattened from Java inner class StripCode.StripCodeConfigRecord).
//
// Java: StripCodeConfigRecord.get()
func (s *StripCode) configRecord() string {
	var builder strings.Builder
	builder.WriteString("=== stripNameSuffixes ===\n")
	for _, e := range s.stripNameSuffixes.Entries() {
		builder.WriteString(e.Key)
		builder.WriteByte('\n')
	}
	builder.WriteByte('\n')
	builder.WriteString("=== stripNamePrefixes ===\n")
	for _, e := range s.stripNamePrefixes.Entries() {
		builder.WriteString(e.Key)
		builder.WriteByte('\n')
	}
	builder.WriteByte('\n')
	builder.WriteString("=== stripTypesList ===\n")
	for _, e := range s.stripTypes.Entries() {
		builder.WriteString(e.Key)
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

func qualifiedNameOrNull(n rhino.Node) string {
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
	s.Add(rhino.Token_ASSIGN_EXPONENT)
	s.Add(rhino.Token_ASSIGN_OR)
	s.Add(rhino.Token_ASSIGN_AND)
	s.Add(rhino.Token_ASSIGN_COALESCE)
	s.Add(rhino.Token_INC)
	s.Add(rhino.Token_DEC)
	s.Add(rhino.Token_OBJECTLIT)
	s.Add(rhino.Token_EXPR_RESULT)
	s.Add(rhino.Token_CLASS)
	return s
}()

func (*StripCode) visitTokens() *ast.TokenSet {
	return &stripCodeVisitTokens
}

// Java: Strip.shouldTraverse(NodeTraversal, Node, Node)
func (s *StripCode) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Here we check for cases where we're going to remove a large chunk of code,
	// and so should not traverse into it to avoid both wasting time and causing
	// problems with logic duplication.
	switch n.GetToken() {
	case rhino.Token_CALL, rhino.Token_NEW:
		// If we're removing the whole call / new
		if s.IsMethodOrCtorCallThatTriggersRemoval(t, n, parent) {
			if s.decisionsLog.IsLogging() {
				s.decisionsLog.LogString("removing function call")
			}
			s.ReplaceHighestNestedCallWithNull(t, n, parent)
			return false
		}
		return true
	default:
		return true
	}
}

// Java: Strip.visit(NodeTraversal, Node, Node)
func (s *StripCode) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		s.RemoveVarDeclarationsByNameOrRvalue(t, n, parent)
	case rhino.Token_NAME:
		s.MaybeRemoveReferenceToRemovedVariable(t, n, parent)
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
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_ASSIGN_OR,
		rhino.Token_ASSIGN_AND,
		rhino.Token_ASSIGN_COALESCE,
		rhino.Token_INC,
		rhino.Token_DEC:
		s.MaybeEliminateAssignmentByLvalueName(t, n, parent)
	case rhino.Token_OBJECTLIT:
		s.EliminateKeysWithStripNamesFromObjLit(n)
	case rhino.Token_EXPR_RESULT:
		s.MaybeEliminateExpressionByName(n)
	case rhino.Token_CLASS:
		s.MaybeEliminateClassByNameOrExtends(t, n, parent)
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
func (s *StripCode) RemoveVarDeclarationsByNameOrRvalue(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
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
		if s.IsStripName(possibleStripName) ||
			s.QualifiedNameBeginsWithStripType(nameNode) ||
			(hasInit && s.IsCallWhoseReturnValueShouldBeStripped(initExpr.Node())) {
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
				if s.decisionsLog.IsLogging() {
					s.decisionsLog.LogString(name + ": initialize with null (" + possibleStripName + ")")
				}
				if hasInit {
					s.ReplaceWithNull(initExpr.Node())
				} else {
					// `var my$name = null;` is a bit easier to optimize away than
					// `var my$name;`, because we can clearly see that it is initialized, so we have
					// a value to inline in later passes.
					nameDecl.SetInit(ast.MustExpr(rhino.IR_NullNode().Srcref(nameNode)))
				}
				t.ReportCodeChange()
			} else {
				// Assume that the declaration comes before any reference.
				// We will remove the references when we see them later.
				if s.decisionsLog.IsLogging() {
					s.decisionsLog.LogString(name + ": removing declaration")
				}
				nameDecl.Detach()
				NodeUtil_MarkFunctionsDeleted(nameNode, s.compiler)
			}
		}
	}
	if !decl.HasDeclarators() && !parent.IsNil() {
		// Must also remove the VAR.
		s.ReplaceWithEmpty(n, parent)
		t.ReportCodeChange()
	}
}

// MaybeRemoveReferenceToRemovedVariable removes a reference if it is a
// reference to a removed variable.
//
// t is the traversal, n a NAME node, parent n's parent.
//
// Java: Strip.maybeRemoveReferenceToRemovedVariable(NodeTraversal, Node, Node)
func (s *StripCode) MaybeRemoveReferenceToRemovedVariable(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsNil() {
		return
	}
	parent = n.GetParent()
	if parent.IsNil() {
		return
	}
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
		if ast.MustNormalGet(parent).Object().Is(n) && s.IsReferenceToRemovedVar(n) {
			if s.decisionsLog.IsLogging() {
				s.decisionsLog.LogString(n.GetString() + ": removing getelem/getprop/call chain")
			}
			s.ReplaceHighestNestedCallWithNull(t, parent, parent.GetParent())
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
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_ASSIGN_OR,
		rhino.Token_ASSIGN_AND,
		rhino.Token_ASSIGN_COALESCE:
		if s.IsReferenceToRemovedVar(n) {
			assign := ast.MustAssignment(parent)
			if assign.Target().Is(n) {
				grandparent := parent.GetParent()
				if grandparent.IsNil() {
					return
				}
				if s.decisionsLog.IsLogging() {
					s.decisionsLog.LogString(qualifiedNameOrNull(n) + ": removing assignment to stripped var")
				}
				if grandparent.IsExprResult() {
					// Remove the assignment.
					if greatGrandparent := grandparent.GetParent(); !greatGrandparent.IsNil() {
						s.ReplaceWithEmpty(grandparent, greatGrandparent)
						t.ReportCodeChange()
					}
				} else {
					// Substitute the r-value for the assignment.
					rvalue := assign.Value().Node()
					if !rvalue.IsNil() {
						rvalue.Detach()
						parent.ReplaceWith(rvalue)
						t.ReportCodeChange()
					}
				}
			} else {
				// The var reference is the r-value. Replace it with null.
				if s.decisionsLog.IsLogging() {
					s.decisionsLog.LogString(qualifiedNameOrNull(n) + ": replacing rhs reference with null")
				}
				s.ReplaceWithNull(n)
				t.ReportCodeChange()
			}
		}
	case rhino.Token_INC, rhino.Token_DEC:
		if s.IsReferenceToRemovedVar(n) {
			grandparent := parent.GetParent()
			if grandparent.IsNil() {
				return
			}
			if s.decisionsLog.IsLogging() {
				s.decisionsLog.LogString(qualifiedNameOrNull(n) + ": removing assignment to stripped var")
			}
			if grandparent.IsExprResult() {
				greatGrandparent := grandparent.GetParent()
				if !greatGrandparent.IsNil() {
					s.ReplaceWithEmpty(grandparent, greatGrandparent)
					t.ReportCodeChange()
				}
			} else {
				s.ReplaceWithNull(parent)
				t.ReportCodeChange()
			}
		}
	case rhino.Token_NEW, rhino.Token_CALL:
		if !ast.MustCallOrNew(parent).Callee().Is(n) && s.IsReferenceToRemovedVar(n) {
			// NOTE: the callee is handled when we visit the CALL or NEW node
			if s.decisionsLog.IsLogging() {
				s.decisionsLog.LogString(qualifiedNameOrNull(n) + ": replacing parameter reference with null")
			}
			s.ReplaceWithNull(n)
			t.ReportCodeChange()
		}
	case rhino.Token_COMMA:
		grandparent := parent.GetParent()
		// The last child in a comma expression is its result, so we need to be careful replacing
		// it with null. We don't want to replace the entire comma expression with null because
		// there could be other elements in it with side-effects.
		bin, isBin := ast.AsBinary(parent)
		isLastChild := isBin && bin.Right().Is(n)
		// The only parent where replacing with null is an issue is likely where the comma
		// expression is the first child (callee) of a CALL or NEW node.
		parentIsCallee := (grandparent.IsCall() || grandparent.IsNew()) && grandparent.GetFirstChild() == parent
		isSafeToRemove := !isLastChild || !parentIsCallee
		if isSafeToRemove && s.IsReferenceToRemovedVar(n) {
			if s.decisionsLog.IsLogging() {
				s.decisionsLog.LogString(qualifiedNameOrNull(n) + ": replacing reference in comma expr with null")
			}
			s.ReplaceWithNull(n)
			t.ReportCodeChange()
		}
	default:
		if s.IsReferenceToRemovedVar(n) {
			if s.decisionsLog.IsLogging() {
				s.decisionsLog.LogString(qualifiedNameOrNull(n) + ": replacing reference with null")
			}
			s.ReplaceWithNull(n)
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
func (s *StripCode) ReplaceHighestNestedCallWithNull(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	if node.IsNil() || parent.IsNil() {
		return
	}
	ancestor := parent
	ancestorChild := node
	for {
		ancestorParent := ancestor.GetParent()
		if ancestorParent.IsNil() {
			return // the call was already removed from the AST
		}

		if exprStmt, ok := ast.AsExprStmt(ancestor); ok {
			if exprStmt.Expr().Is(ancestorChild) {
				// Remove the entire expression statement.
				s.ReplaceWithEmpty(ancestor, ancestorParent)
			} else {
				s.ReplaceWithNull(ancestorChild)
			}
			break
		}
		if assign, ok := ast.AsAssignment(ancestor); ok {
			if assign.Target().Is(ancestorChild) && !assign.Value().Node().IsNil() {
				ancestor.ReplaceWith(assign.Value().Detach().Node())
			} else {
				s.ReplaceWithNull(ancestorChild)
			}
			break
		}
		if incDec, ok := ast.AsIncDec(ancestor); ok {
			if incDec.Target().Is(ancestorChild) {
				if ancestorParent.IsExprResult() {
					if ggp := ancestorParent.GetParent(); !ggp.IsNil() {
						s.ReplaceWithEmpty(ancestorParent, ggp)
					}
				} else {
					s.ReplaceWithNull(ancestor)
				}
			} else {
				s.ReplaceWithNull(ancestorChild)
			}
			break
		}
		if get, ok := ast.AsNormalGet(ancestor); ok {
			if !get.Object().Is(ancestorChild) {
				s.ReplaceWithNull(ancestorChild)
				break
			}
		} else if call, ok := ast.AsCall(ancestor); ok {
			if !call.Callee().Is(ancestorChild) {
				s.ReplaceWithNull(ancestorChild)
				break
			}
		} else {
			s.ReplaceWithNull(ancestorChild)
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
// t is the traversal, n an assignment or inc/dec node, parent n's parent.
//
// Java: Strip.maybeEliminateAssignmentByLvalueName(NodeTraversal, Node, Node)
func (s *StripCode) MaybeEliminateAssignmentByLvalueName(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsNil() || parent.IsNil() || n.GetParent().IsNil() {
		return
	}
	// ASSIGN* / INC / DEC
	//   l-value
	//   [r-value]
	var lvalue rhino.Node
	if incDec, ok := ast.AsIncDec(n); ok {
		lvalue = incDec.Target().Node()
	} else if assign, ok := ast.AsAssignment(n); ok {
		lvalue = assign.Target().Node()
	}
	if lvalue.IsNil() {
		return
	}
	if s.NameIncludesFieldNameToStrip(lvalue) || s.QualifiedNameBeginsWithStripType(lvalue) {
		// Limit to EXPR_RESULT because it is not
		// safe to eliminate assignment in complex expressions,
		// e.g. in ((x = 7) + 8)
		if parent.IsExprResult() {
			if s.decisionsLog.IsLogging() {
				s.decisionsLog.LogString(lvalue.GetString() + ": removing assignment statement")
			}
			grandparent := parent.GetParent()
			// the assignment may already have been removed when visiting either the lhs
			// or the rhs.
			if !grandparent.IsNil() {
				s.ReplaceWithEmpty(parent, grandparent)
				s.compiler.ReportChangeToEnclosingScope(grandparent)
			}
		} else {
			t.Report(n, StripCode_STRIP_ASSIGNMENT_ERROR, qualifiedNameOrNull(lvalue))
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
func (s *StripCode) MaybeEliminateExpressionByName(n rhino.Node) {
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
	if expression.IsNil() {
		return
	}
	if s.NameIncludesFieldNameToStrip(expression) ||
		s.QualifiedNameBeginsWithStripType(expression) {
		if s.decisionsLog.IsLogging() {
			s.decisionsLog.LogString(expression.GetString() + ": removing property declaration statement")
		}
		s.ReplaceWithEmpty(n, parent)
		s.compiler.ReportChangeToEnclosingScope(parent)
	}
}

// EliminateKeysWithStripNamesFromObjLit eliminates any object literal keys
// in an object literal declaration that have strip names.
//
// n is an OBJLIT node.
//
// Java: Strip.eliminateKeysWithStripNamesFromObjLit(Node)
func (s *StripCode) EliminateKeysWithStripNamesFromObjLit(n rhino.Node) {
	// OBJLIT
	//   key1
	//     value1
	//   key2
	//   ...
	for member := range ast.ReadAhead(ast.MustObjectLit(n).FirstMember()) {
		if key, ok := member.AsPlainKeyMember(); ok && s.IsStripName(key.Key()) {
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
func (s *StripCode) MaybeEliminateClassByNameOrExtends(t *NodeTraversal, classNode rhino.Node, parent rhino.Node) {
	cls := ast.MustClass(classNode)
	nameNode := NodeUtil_GetNameNode(classNode)
	className := "<anonymous>"
	// Replace class with null if it is a strip type
	if qname, ok := nameNode.GetQualifiedNameString(); ok {
		className = qname
		if s.qualifiedNameBeginsWithStripType(className, true) {
			if s.decisionsLog.IsLogging() {
				s.decisionsLog.LogString(className + ": removing class")
			}
			if NodeUtil_IsStatementParent(parent) {
				s.ReplaceWithEmpty(classNode, parent)
			} else {
				s.ReplaceWithNull(classNode)
			}
			t.ReportCodeChange()
			return
		}
	}

	// If the class is not a strip type, the superclass also cannot be a strip type
	if superclass, ok := cls.Extends(); ok {
		if superclassName, ok := superclass.Node().GetQualifiedNameString(); ok &&
			s.qualifiedNameBeginsWithStripType(superclassName, true) {
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
func (s *StripCode) IsCallWhoseReturnValueShouldBeStripped(n rhino.Node) bool {
	if (!n.IsCall() && !n.IsNew()) || !n.HasChildren() {
		return false
	}

	function := NodeUtil_GetCallTargetResolvingIndirectCalls(n)
	return s.QualifiedNameBeginsWithStripType(function) || s.NameIncludesFieldNameToStrip(function)
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
func (s *StripCode) QualifiedNameBeginsWithStripType(n rhino.Node) bool {
	name, ok := n.GetQualifiedNameString()
	return s.qualifiedNameBeginsWithStripType(name, ok)
}

func (s *StripCode) qualifiedNameBeginsWithStripType(name string, hasName bool) bool {
	if hasName {
		if s.stripTypes.Has(name) {
			if s.decisionsLog.IsLogging() {
				s.LogStripName(name, "equals strip type")
			}
			return true
		}
		for _, prefix := range s.stripTypePrefixesList {
			if strings.HasPrefix(name, prefix) {
				if s.decisionsLog.IsLogging() {
					s.LogStripName(name, "starts with strip type prefix")
				}
				return true
			}
		}
		if s.decisionsLog.IsLogging() {
			s.LogNotAStripName(name, "does not begin with a strip type")
		}
		return false
	}
	if s.decisionsLog.IsLogging() {
		s.LogNotAStripName("null", "does not begin with a strip type")
	}
	return false
}

// IsReferenceToRemovedVar determines whether a NAME node represents a
// reference to a variable that has been removed.
//
// n is a NAME node. Returns whether the variable was removed.
//
// Java: Strip.isReferenceToRemovedVar(Node)
func (s *StripCode) IsReferenceToRemovedVar(n rhino.Node) bool {
	if _, ok := ast.AsName(n); !ok {
		return false
	}
	_, ok := s.varsToRemove[n.GetStringID()]
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
func (s *StripCode) IsMethodOrCtorCallThatTriggersRemoval(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// CALL/NEW
	//   GETPROP (function)         <-- we're interested in this, the function
	//     GETPROP (callee object)  <-- or the object on which it is called
	//       ...
	//       STRING (field name)
	//     STRING (method name)
	//   ... (arguments)

	if (!n.IsCall() && !n.IsNew()) || !n.HasChildren() {
		return false
	}
	function := NodeUtil_GetCallTargetResolvingIndirectCalls(n)
	if function.IsNil() || !function.IsQualifiedName() {
		return false
	}

	if parentName, isName := ast.AsName(parent); isName {
		if init, hasInit := parentName.Init(); hasInit && init.Is(n) {
			if _, isDecl := ast.AsNameDeclaration(parent.GetParent()); isDecl {
				// The call's return value is being used to initialize a newly
				// declared variable. We should leave the call intact for now.
				// That way, when the traversal reaches the variable declaration,
				// we'll recognize that the variable and all references to it need
				// to be eliminated.
				return false
			}
		}
	}

	if fnName, ok := ast.AsName(function); ok && s.IsStripName(fnName.Ident()) {
		return true
	}
	var callee rhino.Node
	if gp, ok := ast.AsGetProp(function); ok {
		callee = gp.Object().Node()
	}
	return s.NameIncludesFieldNameToStrip(callee) ||
		s.NameIncludesFieldNameToStrip(function) ||
		s.QualifiedNameBeginsWithStripType(function) ||
		s.ActsOnStripType(t, n)
}

// NameIncludesFieldNameToStrip returns whether a name includes a field name
// that should be stripped. E.g., "foo.stripMe.bar", "(foo.bar).stripMe",
// etc.
//
// Java: Strip.nameIncludesFieldNameToStrip(Node)
func (s *StripCode) NameIncludesFieldNameToStrip(n rhino.Node) bool {
	for cur := n; !cur.IsNil(); {
		if prop, ok := ast.AsGetProp(cur); ok {
			if s.IsStripName(prop.Property()) {
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
				if s.IsStripName(nameString[:end]) {
					return true
				}
				nameString = nameString[end+1:]
				end = strings.IndexByte(nameString, '$')
				if end == -1 {
					return s.IsStripName(nameString)
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
func (s *StripCode) ActsOnStripType(t *NodeTraversal, callNode rhino.Node) bool {
	classes := s.compiler.GetCodingConvention().GetClassesDefinedByCall(callNode)
	if classes != nil {
		// It's okay to strip a type that inherits from a non-stripped type
		// e.g. goog.inherits(goog.debug.Logger, Object)
		if s.qualifiedNameBeginsWithStripType(classes.subclassName, true) {
			s.LogStripName(classes.subclassName, "class defining call")
			return true
		}

		// report an error if a non-strip type inherits from a
		// strip type.
		if s.qualifiedNameBeginsWithStripType(classes.superclassName, true) {
			t.Report(
				callNode, StripCode_STRIP_TYPE_INHERIT_ERROR, classes.subclassName, classes.superclassName)
		}
	}

	return false
}

func startsWithUpperCase(name string) bool {
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
func (s *StripCode) IsStripName(name string) bool {
	if s.stripNameSuffixes.Has(name) {
		s.LogStripName(name, "matches a suffix")
		return true
	}

	if s.stripNamePrefixes.Has(name) {
		s.LogStripName(name, "matches a prefix")
		return true
	}

	if name == "" || startsWithUpperCase(name) {
		s.LogNotAStripName(name, "empty or starts with uppercase")
		return false
	}

	if len(s.stripNamePrefixesLowerCaseList) > 0 || len(s.stripNameSuffixesLowerCaseList) > 0 {
		lcName := jutil.StringToLowerCaseRoot(name)
		for _, stripName := range s.stripNamePrefixesLowerCaseList {
			if strings.HasPrefix(lcName, stripName) {
				if s.decisionsLog.IsLogging() {
					s.LogStripName(name, "matches lc prefix: "+stripName)
				}
				return true
			}
		}

		for _, stripName := range s.stripNameSuffixesLowerCaseList {
			if strings.HasSuffix(lcName, stripName) {
				if s.decisionsLog.IsLogging() {
					s.LogStripName(name, "matches lc suffix: "+stripName)
				}
				return true
			}
		}
	}

	s.LogNotAStripName(name, "no matches")
	return false
}

// Java: Strip.logNotAStripName(String, String)
func (s *StripCode) LogNotAStripName(name string, reason string) {
	if s.decisionsLog.IsLogging() {
		s.decisionsLog.LogString(name + "\tnot a strip name: " + reason)
	}
}

// Java: Strip.logStripName(String, String)
func (s *StripCode) LogStripName(name string, reason string) {
	if s.decisionsLog.IsLogging() {
		s.decisionsLog.LogString(name + "\tstrip name: " + reason)
	}
}

// ReplaceWithNull replaces a node with a NULL node. This is useful where a
// value is expected.
//
// Java: Strip.replaceWithNull(Node)
func (s *StripCode) ReplaceWithNull(n rhino.Node) {
	if s.decisionsLog.IsLogging() {
		s.decisionsLog.LogString("replace with null: " + n.GetLocation())
	}
	n.ReplaceWith(rhino.IR_NullNode())
	NodeUtil_MarkFunctionsDeleted(n, s.compiler)
}

// ReplaceWithEmpty replaces a node with an EMPTY node. This is useful where
// a statement is expected.
//
// n is a node, parent n's parent.
//
// Java: Strip.replaceWithEmpty(Node, Node)
func (s *StripCode) ReplaceWithEmpty(n rhino.Node, parent rhino.Node) {
	if s.decisionsLog.IsLogging() {
		s.decisionsLog.LogString("replace with empty: " + n.GetLocation())
	}
	NodeUtil_RemoveChild(parent, n)
	NodeUtil_MarkFunctionsDeleted(n, s.compiler)
}
