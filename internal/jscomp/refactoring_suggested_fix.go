// Port of com/google/javascript/refactoring/SuggestedFix.java @ closure-compiler v20261015

/*
 * Copyright 2014 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SuggestedFix is the object representing the fixes to apply to the source
// code to create the refactoring CL. To create a class, use the
// SuggestedFix_Builder class and helper functions.
//
// Java: public final class SuggestedFix
type SuggestedFix struct {
	matchedNodeInfo *SuggestedFix_MatchedNodeInfo
	// Multimap of filename to a modification to that file.
	replacements *jutil.SetMultimap[string, *CodeReplacement]

	// An optional description of the fix, to distinguish between the various
	// possible fixes for errors that have multiple fixes.
	description *string

	// Alternative fixes for the same problem. The fix itself is always the
	// first entry in this list. If you cannot ask the developer which fix is
	// appropriate, apply the first fix instead of any alternatives.
	alternatives []*SuggestedFix
}

// SuggestedFix_ImportType is the kind of import a fix adds.
//
// Java: static enum SuggestedFix.ImportType
type SuggestedFix_ImportType int32

const (
	SuggestedFix_ImportType_REQUIRE SuggestedFix_ImportType = iota
	SuggestedFix_ImportType_REQUIRE_TYPE
)

var suggestedFix_ImportType_names = [...]string{"REQUIRE", "REQUIRE_TYPE"}

// SuggestedFix_ImportType_values is ImportType.values().
func SuggestedFix_ImportType_values() []SuggestedFix_ImportType {
	return []SuggestedFix_ImportType{SuggestedFix_ImportType_REQUIRE, SuggestedFix_ImportType_REQUIRE_TYPE}
}

// Name is Enum.name().
func (t SuggestedFix_ImportType) Name() string { return suggestedFix_ImportType_names[t] }

// String is Enum.toString().
func (t SuggestedFix_ImportType) String() string { return t.Name() }

// Ordinal is Enum.ordinal().
func (t SuggestedFix_ImportType) Ordinal() int { return int(t) }

// Java: SuggestedFix(MatchedNodeInfo, SetMultimap<String, CodeReplacement>, String, ImmutableList<SuggestedFix>)
func newSuggestedFix(
	matchedNodeInfo *SuggestedFix_MatchedNodeInfo,
	replacements *jutil.SetMultimap[string, *CodeReplacement],
	description *string,
	alternatives []*SuggestedFix) *SuggestedFix {
	f := &SuggestedFix{}
	f.matchedNodeInfo = matchedNodeInfo
	f.replacements = replacements
	f.description = description
	all := make([]*SuggestedFix, 0, 1+len(alternatives))
	all = append(all, f)
	all = append(all, alternatives...)
	f.alternatives = all
	return f
}

// GetMatchedNodeInfo returns information about the original JS Compiler
// Node that caused this SuggestedFix to be constructed.
//
// Java: SuggestedFix.getMatchedNodeInfo()
func (f *SuggestedFix) GetMatchedNodeInfo() *SuggestedFix_MatchedNodeInfo {
	return f.matchedNodeInfo
}

// GetReplacements returns a multimap from filename to all the replacements
// that should be applied for this given fix.
//
// Java: SuggestedFix.getReplacements()
func (f *SuggestedFix) GetReplacements() *jutil.SetMultimap[string, *CodeReplacement] {
	return f.replacements
}

// GetDescription may return nil.
//
// Java: SuggestedFix.getDescription()
func (f *SuggestedFix) GetDescription() *string {
	return f.description
}

// GetAlternatives gets all possible fixes for this problem, including this
// fix.
//
// Java: SuggestedFix.getAlternatives()
func (f *SuggestedFix) GetAlternatives() []*SuggestedFix {
	return f.alternatives
}

// GetNonDefaultAlternatives gets all alternative fixes, excluding this fix.
//
// Java: SuggestedFix.getNonDefaultAlternatives()
func (f *SuggestedFix) GetNonDefaultAlternatives() []*SuggestedFix {
	return f.alternatives[1:len(f.alternatives)]
}

// Java: SuggestedFix.isNoOp()
func (f *SuggestedFix) IsNoOp() bool {
	return f.replacements.IsEmpty()
}

// Java: SuggestedFix.toString()
func (f *SuggestedFix) String() string {
	if f.IsNoOp() {
		return "<no-op SuggestedFix>"
	}
	var sb strings.Builder
	for key, values := range f.replacements.AsMap().All() {
		sb.WriteString("Replacements for file: ")
		sb.WriteString(key)
		sb.WriteString("\n")
		// Joiner.on("\n\n").appendTo(sb, entry.getValue())
		first := true
		for replacement := range values.All() {
			if !first {
				sb.WriteString("\n\n")
			}
			first = false
			sb.WriteString(replacement.String())
		}
	}
	return sb.String()
}

// SuggestedFix_Builder is the builder class for SuggestedFix that contains
// helper functions to manipulate JS nodes.
//
// Java: public static final class SuggestedFix.Builder
type SuggestedFix_Builder struct {
	matchedNodeInfo *SuggestedFix_MatchedNodeInfo
	// Go-only: language-adapter — ImmutableSetMultimap.Builder: keys and
	// values in first-insertion order, values deduplicated by equals; Build
	// copies it, as ImmutableSetMultimap.Builder.build() does.
	replacements *jutil.SetMultimap[string, *CodeReplacement]
	alternatives []*SuggestedFix
	description  *string
}

// Java: SuggestedFix.Builder()
func NewSuggestedFix_Builder() *SuggestedFix_Builder {
	b := &SuggestedFix_Builder{}
	b.matchedNodeInfo = nil
	b.replacements = jutil.NewSetMultimap[string, *CodeReplacement]()
	b.alternatives = []*SuggestedFix{}
	b.description = nil
	return b
}

// AttachMatchedNodeInfo sets the node on this SuggestedFix that caused this
// SuggestedFix to be built in the first place.
//
// Java: SuggestedFix.Builder.attachMatchedNodeInfo(Node, AbstractCompiler)
func (b *SuggestedFix_Builder) AttachMatchedNodeInfo(node *rhino.Node, compiler AbstractCompiler) *SuggestedFix_Builder {
	b.matchedNodeInfo =
		SuggestedFix_MatchedNodeInfo_Create(
			node, RefactoringUtils_IsInClosurizedFile(node, NewNodeMetadata(compiler)))
	return b
}

// Java: SuggestedFix.Builder.addAlternative(SuggestedFix)
func (b *SuggestedFix_Builder) AddAlternative(alternative *SuggestedFix) *SuggestedFix_Builder {
	jutil.CheckStateMsg(
		len(alternative.GetNonDefaultAlternatives()) == 0,
		"Alternative SuggestedFix must have no alternatives of their own.")
	b.alternatives = append(b.alternatives, alternative)
	return b
}

// ReplaceText replaces text starting at the given node position.
//
// Java: SuggestedFix.Builder.replaceText(Node, int, String)
func (b *SuggestedFix_Builder) ReplaceText(node *rhino.Node, length int, newContent string) *SuggestedFix_Builder {
	startPosition := node.GetSourceOffset()
	b.replacements.Put(
		*node.GetSourceFileName(), CodeReplacement_Create(startPosition, length, newContent))
	return b
}

// AddChildToFront inserts a new node as the first child of the provided
// node.
//
// Java: SuggestedFix.Builder.addChildToFront(Node, String)
func (b *SuggestedFix_Builder) AddChildToFront(parentNode *rhino.Node, content string) *SuggestedFix_Builder {
	jutil.CheckStateMsg(
		parentNode.IsBlock(), "addChildToFront is only supported for BLOCK statements.")
	startPosition := parentNode.GetSourceOffset() + 1
	b.replacements.Put(
		*parentNode.GetSourceFileName(), CodeReplacement_Create(startPosition, 0, "\n"+content))
	return b
}

// InsertAfter inserts the text after the given node.
//
// Java: SuggestedFix.Builder.insertAfter(Node, String)
func (b *SuggestedFix_Builder) InsertAfter(node *rhino.Node, text string) *SuggestedFix_Builder {
	position := node.GetSourceOffset() + node.GetLength()
	b.replacements.Put(*node.GetSourceFileName(), CodeReplacement_Create(position, 0, text))
	return b
}

// InsertBeforeNode inserts a new node before the provided node.
//
// Java: SuggestedFix.Builder.insertBefore(Node, Node, AbstractCompiler)
func (b *SuggestedFix_Builder) InsertBeforeNode(nodeToInsertBefore *rhino.Node, n *rhino.Node, compiler AbstractCompiler) *SuggestedFix_Builder {
	return b.InsertBeforeNodeWithSortKey(nodeToInsertBefore, n, compiler, "")
}

// Java: SuggestedFix.Builder.insertBefore(Node, Node, AbstractCompiler, String)
func (b *SuggestedFix_Builder) InsertBeforeNodeWithSortKey(
	nodeToInsertBefore *rhino.Node, n *rhino.Node, compiler AbstractCompiler, sortKey string) *SuggestedFix_Builder {
	return b.InsertBeforeWithSortKey(nodeToInsertBefore, b.GenerateCode(compiler, n), sortKey)
}

// InsertBefore inserts a string before the provided node. This is useful
// for inserting comments into a file since the JS Compiler doesn't
// currently support printing comments.
//
// Java: SuggestedFix.Builder.insertBefore(Node, String)
func (b *SuggestedFix_Builder) InsertBefore(nodeToInsertBefore *rhino.Node, content string) *SuggestedFix_Builder {
	return b.InsertBeforeWithSortKey(nodeToInsertBefore, content, "")
}

// Java: SuggestedFix.Builder.insertBefore(Node, String, String)
func (b *SuggestedFix_Builder) InsertBeforeWithSortKey(nodeToInsertBefore *rhino.Node, content string, sortKey string) *SuggestedFix_Builder {
	startPosition := SuggestedFix_GetStartPositionForNodeConsideringComments(nodeToInsertBefore)
	jutil.CheckNotNullMsg1(nodeToInsertBefore.GetSourceFileName(),
		"No source file name for node: %s", nodeToInsertBefore)
	b.replacements.Put(
		*nodeToInsertBefore.GetSourceFileName(),
		CodeReplacement_CreateWithSortKey(startPosition, 0, content, sortKey))
	return b
}

// Delete deletes a node and its contents from the source file. If the node
// is a child of a block or top level statement, this will also delete the
// whitespace before the node.
//
// Java: SuggestedFix.Builder.delete(Node)
func (b *SuggestedFix_Builder) Delete(n *rhino.Node) *SuggestedFix_Builder {
	return b.DeleteWithWhitespaceBefore(n, true)
}

// DeleteWithWhitespaceBefore deletes a node and its contents from the
// source file.
//
// Java: SuggestedFix.Builder.delete(Node, boolean)
func (b *SuggestedFix_Builder) DeleteWithWhitespaceBefore(n *rhino.Node, deleteWhitespaceBefore bool) *SuggestedFix_Builder {
	startPosition := SuggestedFix_GetStartPositionForNodeConsideringComments(n)
	startOffsetWithoutComments := n.GetSourceOffset()
	length := (startOffsetWithoutComments - startPosition) + n.GetLength()

	if n.GetNext() != nil &&
		NodeUtil_GetBestJSDocInfo(n.GetNext()) == nil &&
		n.GetNext().GetNonJSDocComment() == nil {
		length = n.GetNext().GetSourceOffset() - startPosition
	}

	// Variable declarations and string keys require special handling since
	// the node doesn't contain enough if it has a child. The NAME node in a
	// var/let/const declaration doesn't include its child in its length, and
	// the code needs to know how to delete the commas. The same is true for
	// string keys in object literals and object destructuring patterns.
	// TODO(mknichel): Move this logic and the start position logic to a
	// helper function so that it can be reused in other methods.
	if (n.IsName() && NodeUtil_IsNameDeclaration(n.GetParent())) || n.IsStringKey() {
		if n.GetNext() != nil {
			length = SuggestedFix_GetStartPositionForNodeConsideringComments(n.GetNext()) - startPosition
		} else if n.HasChildren() {
			child := n.GetFirstChild()
			length = (child.GetSourceOffset() + child.GetLength()) - startPosition
		}
		if n.GetParent().GetLastChild() == n && n != n.GetParent().GetFirstChild() {
			previousSibling := n.GetPrevious()
			if previousSibling.HasChildren() {
				child := previousSibling.GetFirstChild()
				startPositionDiff := startPosition - (child.GetSourceOffset() + child.GetLength())
				startPosition -= startPositionDiff
				length += startPositionDiff
			} else {
				startPositionDiff :=
					startPosition - (previousSibling.GetSourceOffset() + previousSibling.GetLength())
				startPosition -= startPositionDiff
				length += startPositionDiff
			}
		}
	}

	parent := n.GetParent()
	if deleteWhitespaceBefore &&
		parent != nil &&
		(parent.IsScript() || parent.IsBlock()) {
		previousSibling := n.GetPrevious()
		if previousSibling != nil {
			previousSiblingEndPosition :=
				previousSibling.GetSourceOffset() + previousSibling.GetLength()
			length += (startPosition - previousSiblingEndPosition)
			startPosition = previousSiblingEndPosition
		}
	}
	b.replacements.Put(*n.GetSourceFileName(), CodeReplacement_Create(startPosition, length, ""))
	return b
}

// DeleteWithoutRemovingWhitespaceBefore deletes a node and its contents
// from the source file.
//
// Java: SuggestedFix.Builder.deleteWithoutRemovingWhitespaceBefore(Node)
func (b *SuggestedFix_Builder) DeleteWithoutRemovingWhitespaceBefore(n *rhino.Node) *SuggestedFix_Builder {
	return b.DeleteWithWhitespaceBefore(n, false)
}

// DeleteWithoutRemovingWhitespace deletes a node without touching any
// surrounding whitespace.
//
// Java: SuggestedFix.Builder.deleteWithoutRemovingWhitespace(Node)
func (b *SuggestedFix_Builder) DeleteWithoutRemovingWhitespace(n *rhino.Node) *SuggestedFix_Builder {
	b.replacements.Put(
		*n.GetSourceFileName(), CodeReplacement_Create(n.GetSourceOffset(), n.GetLength(), ""))
	return b
}

// Rename renames a given node to the provided name.
//
// n is the node to rename; name is the new name for the node.
//
// Java: SuggestedFix.Builder.rename(Node, String)
func (b *SuggestedFix_Builder) Rename(n *rhino.Node, name string) *SuggestedFix_Builder {
	return b.RenameWithReplaceNameSubtree(n, name, false)
}

// RenameWithReplaceNameSubtree renames a given node to the provided name.
//
// n is the node to rename; name is the new name for the node;
// replaceNameSubtree is true to replace the entire name subtree below the
// node. The default is to replace just the last property in the node with
// the new name. For instance, if replaceNameSubtree is false, then
// `this.foo()` will be renamed to `this.bar()`. However, if it is true, it
// will be renamed to `bar()`.
//
// Java: SuggestedFix.Builder.rename(Node, String, boolean)
func (b *SuggestedFix_Builder) RenameWithReplaceNameSubtree(n *rhino.Node, name string, replaceNameSubtree bool) *SuggestedFix_Builder {
	var range_ *rhino.Node
	switch n.GetToken() {
	case rhino.Token_CALL, rhino.Token_TAGGED_TEMPLATELIT:
		return b.RenameWithReplaceNameSubtree(n.GetFirstChild(), name, replaceNameSubtree)
	case rhino.Token_GETPROP:
		if replaceNameSubtree {
			range_ = SuggestedFix_Builder_SubtreeRangeOfIdentifier(n)
		} else {
			range_ = n
		}
	case rhino.Token_STRINGLIT:
		jutil.CheckStateMsg1(n.GetParent().IsGetProp(), "%s", n)
		range_ = n
	case rhino.Token_STRING_KEY, rhino.Token_NAME:
		range_ = n
	default:
		panic(jutil.NewUnsupportedOperationExceptionMsg(
			"Rename is not implemented for this node type: " + n.String()))
	}

	b.replacements.Put(
		*range_.GetSourceFileName(),
		CodeReplacement_Create(range_.GetSourceOffset(), range_.GetLength(), name))
	return b
}

// ReplaceRange replaces a range of nodes with the given content.
//
// Java: SuggestedFix.Builder.replaceRange(Node, Node, String)
func (b *SuggestedFix_Builder) ReplaceRange(first *rhino.Node, last *rhino.Node, newContent string) *SuggestedFix_Builder {
	jutil.CheckState(first.GetParent() == last.GetParent())

	start := SuggestedFix_GetStartPositionForNodeConsideringComments(first)
	if start == 0 {
		// if there are file-level comments at the top of the file, we do not
		// wish to remove them
		start = first.GetSourceOffset()
	}
	end := last.GetSourceOffset() + last.GetLength()
	length := end - start
	b.replacements.Put(
		*first.GetSourceFileName(), CodeReplacement_Create(start, length, newContent))
	return b
}

// Replace replaces the provided node with new node in the source file.
//
// Java: SuggestedFix.Builder.replace(Node, Node, AbstractCompiler)
func (b *SuggestedFix_Builder) Replace(original *rhino.Node, newNode *rhino.Node, compiler AbstractCompiler) *SuggestedFix_Builder {
	parent := original.GetParent()
	// EXPR_RESULT nodes will contain the trailing semicolons, but the child
	// node will not. Replace the EXPR_RESULT node to ensure that the
	// semicolons are correct in the final output.
	if parent != nil && parent.IsExprResult() {
		original = parent
	}
	// TODO(mknichel): Move this logic to CodePrinter.
	newCode := b.GenerateCode(compiler, newNode)
	// The generated code may contain a trailing newline but that is never
	// wanted.
	if strings.HasSuffix(newCode, "\n") {
		newCode = newCode[:len(newCode)-1]
	}

	// Most replacements don't need the semicolon in the new generated code -
	// however, some statements that are blocks or expressions will need the
	// semicolon.
	needsSemicolon :=
		parent != nil &&
			(parent.IsExprResult() ||
				parent.IsBlock() ||
				parent.IsScript() ||
				parent.IsModuleBody())
	if strings.HasSuffix(newCode, ";") && !needsSemicolon {
		newCode = newCode[:len(newCode)-1]
	}

	// If the replacement has lower precedence then we may need to add
	// parentheses.
	if parent != nil && rhino.IR_MayBeExpression(parent) {
		replacement := newNode
		for (replacement.IsBlock() || replacement.IsScript() || replacement.IsModuleBody()) &&
			replacement.HasOneChild() {
			replacement = replacement.GetOnlyChild()
		}
		if replacement.IsExprResult() {
			replacement = replacement.GetOnlyChild()
		}
		if rhino.IR_MayBeExpression(replacement) {
			outer := NodeUtil_Precedence(parent.GetToken())
			inner := NodeUtil_Precedence(original.GetToken())
			newInner := NodeUtil_Precedence(replacement.GetToken())
			if newInner < NodeUtil_Precedence(rhino.Token_CALL) && newInner <= outer && inner >= outer {
				newCode = "(" + newCode + ")"
			}
		}
	}

	range_ := original
	if original.IsGetProp() {
		range_ = SuggestedFix_Builder_SubtreeRangeOfIdentifier(original)
	}

	b.replacements.Put(
		*range_.GetSourceFileName(),
		CodeReplacement_Create(range_.GetSourceOffset(), range_.GetLength(), newCode))
	return b
}

// AddCast adds a cast of the given type to the provided node.
//
// Java: SuggestedFix.Builder.addCast(Node, AbstractCompiler, String)
func (b *SuggestedFix_Builder) AddCast(n *rhino.Node, compiler AbstractCompiler, type_ string) *SuggestedFix_Builder {
	// TODO(mknichel): Figure out the best way to output the typecast.
	b.replacements.Put(
		*n.GetSourceFileName(),
		CodeReplacement_Create(
			n.GetSourceOffset(),
			n.GetLength(),
			"/** @type {"+type_+"} */ ("+b.GenerateCode(compiler, n)+")"))
	return b
}

// RemoveCast removes a cast from the given node.
//
// Java: SuggestedFix.Builder.removeCast(Node, AbstractCompiler)
func (b *SuggestedFix_Builder) RemoveCast(n *rhino.Node, compiler AbstractCompiler) *SuggestedFix_Builder {
	jutil.CheckArgument(n.IsCast())
	jsDoc := n.GetJSDocInfo()
	b.replacements.Put(
		*n.GetSourceFileName(),
		CodeReplacement_Create(
			jsDoc.GetOriginalCommentPosition(),
			n.GetFirstChild().GetSourceOffset()-jsDoc.GetOriginalCommentPosition(),
			""))
	b.replacements.Put(
		*n.GetSourceFileName(),
		CodeReplacement_Create(n.GetSourceOffset()+n.GetLength()-1 /* length= */, 1, ""))
	return b
}

// AddOrReplaceJsDoc adds or replaces the JS Doc for the given node.
//
// Java: SuggestedFix.Builder.addOrReplaceJsDoc(Node, String)
func (b *SuggestedFix_Builder) AddOrReplaceJsDoc(n *rhino.Node, newJsDoc string) *SuggestedFix_Builder {
	offset := n.GetSourceOffset()
	length := 0

	if n.IsGetProp() {
		offset = SuggestedFix_Builder_SubtreeRangeOfIdentifier(n).GetSourceOffset()
	}

	jsDoc := NodeUtil_GetBestJSDocInfo(n)
	if jsDoc != nil {
		offset = jsDoc.GetOriginalCommentPosition()
		length = jutil.StringLength(*jsDoc.GetOriginalCommentString())
	}

	b.replacements.Put(*n.GetSourceFileName(), CodeReplacement_Create(offset, length, newJsDoc))
	return b
}

// ChangeJsDocType changes the JS Doc Type of the given node.
//
// Java: SuggestedFix.Builder.changeJsDocType(Node, AbstractCompiler, String)
func (b *SuggestedFix_Builder) ChangeJsDocType(n *rhino.Node, compiler AbstractCompiler, type_ string) *SuggestedFix_Builder {
	typeNode := parsing.JsDocInfoParser_ParseTypeString(type_)
	jutil.CheckNotNullMsg1(typeNode, "Invalid type: %s", type_)
	typeExpr := rhino.NewJSTypeExpression(typeNode, "jsflume")
	newJsType := typeExpr.Evaluate(nil, compiler.GetTypeRegistry())
	if newJsType == nil {
		panic(jutil.NewRuntimeException("JS Compiler does not recognize type: " + type_))
	}

	// TODO(mknichel): Use the JSDocInfoParser to find the end of the type
	// declaration. This would also handle multiple lines, and record types
	// (which contain '{')

	// Only "@type" allows type names without "{}"
	b.ReplaceTypePattern(n, type_, jutil.Pattern_Compile(
		"@(type) *\\{?[^@\\s}]+\\}?"))

	// Text following other annotations may be a comment, not a type.
	b.ReplaceTypePattern(n, type_, jutil.Pattern_Compile(
		"@(export|package|private|protected|public|const|return) *\\{[^}]+\\}"))

	return b
}

// ReplaceTypePattern: the pattern supplied here should have one matching
// group, the annotation with associated the type expression, the entire
// pattern should match the annotation and the type expression to be
// replaced.
//
// Java: SuggestedFix.Builder.replaceTypePattern(Node, String, Pattern)
func (b *SuggestedFix_Builder) ReplaceTypePattern(n *rhino.Node, type_ string, pattern *jutil.Pattern) {
	info := NodeUtil_GetBestJSDocInfo(n)
	jutil.CheckNotNullMsg1(info, "Node %s does not have JS Doc associated with it.", n)
	originalComment := info.GetOriginalCommentString()
	originalPosition := info.GetOriginalCommentPosition()
	if originalComment != nil {
		m := pattern.Matcher(*originalComment)
		for m.Find() {
			b.replacements.Put(
				*n.GetSourceFileName(),
				CodeReplacement_Create(
					originalPosition+m.Start(),
					m.End()-m.Start(),
					"@"+*m.Group(1)+" {"+type_+"}"))
		}
	}
}

// InsertArguments inserts arguments into an existing function call.
//
// Java: SuggestedFix.Builder.insertArguments(Node, int, String...)
func (b *SuggestedFix_Builder) InsertArguments(n *rhino.Node, position int, args ...string) *SuggestedFix_Builder {
	jutil.CheckArgumentMsg(n.IsCall(), "insertArguments is only applicable to function call nodes.")
	var startPosition int
	argument := n.GetSecondChild()
	i := 0
	for argument != nil && i < position {
		argument = argument.GetNext()
		i++
	}
	if argument == nil {
		jutil.CheckArgumentMsg(
			position == i, "The specified position must be less than the number of arguments.")
		startPosition = n.GetSourceOffset() + n.GetLength() - 1
	} else {
		startPosition = SuggestedFix_GetStartPositionForNodeConsideringComments(argument)
	}

	newContent := strings.Join(args, ", ")
	if argument != nil {
		newContent += ", "
	} else if i > 0 {
		newContent = ", " + newContent
	}
	b.replacements.Put(*n.GetSourceFileName(), CodeReplacement_Create(startPosition, 0, newContent))

	return b
}

// DeleteArgument deletes an argument from an existing function call,
// including any JS doc that precedes it. WARNING: If jsdoc erroneously
// follows the argument, it will not be removed as the parser considers the
// comment to belong to the next argument.
//
// Java: SuggestedFix.Builder.deleteArgument(Node, int)
func (b *SuggestedFix_Builder) DeleteArgument(n *rhino.Node, position int) *SuggestedFix_Builder {
	jutil.CheckArgumentMsg(
		n.IsCall() || n.IsNew(), "deleteArgument is only applicable to function call nodes.")

	// A CALL node's first child is the name of the function being called,
	// and subsequent children are the arguments being passed to that
	// function.
	numArguments := n.GetChildCount() - 1
	jutil.CheckStateMsg(
		numArguments > 0, "deleteArgument() cannot be used on a function call with no arguments")
	jutil.CheckArgumentMsg(
		position >= 0 && position < numArguments,
		"The specified position must be less than the number of arguments.")
	argument := n.GetSecondChild()

	// Points at the first position in the code we will remove.
	startOfArgumentToRemove := -1
	// Points one past the last position in the code we will remove.
	endOfArgumentToRemove := -1
	i := 0
	for argument != nil {
		// If we are removing the first argument, we remove from the start of
		// it (including any jsdoc). Otherwise, we remove from the end of the
		// previous argument (to remove the comma and any whitespace).

		// If we are removing the first argument and it's not the only
		// argument, we remove to the beginning of the next argument (to
		// remove the comma and any whitespace). Otherwise we remove to the
		// end of the argument.
		if i < position {
			startOfArgumentToRemove = argument.GetSourceOffset() + argument.GetLength()
		} else if i == position {
			if position == 0 {
				startOfArgumentToRemove = SuggestedFix_GetStartPositionForNodeConsideringComments(argument)
			}
			endOfArgumentToRemove = argument.GetSourceOffset() + argument.GetLength()
		} else if i > position {
			if position == 0 {
				endOfArgumentToRemove = argument.GetSourceOffset()
			}
			// We have all the information we need to remove the argument,
			// break early.
			break
		}

		argument = argument.GetNext()
		i++
	}

	// Remove the argument by replacing it with an empty string.
	lengthOfArgumentToRemove := endOfArgumentToRemove - startOfArgumentToRemove
	b.replacements.Put(
		*n.GetSourceFileName(),
		CodeReplacement_Create(startOfArgumentToRemove, lengthOfArgumentToRemove, ""))
	return b
}

// SuggestedFix_Builder_CreateImportNode: alias may be nil.
//
// Java: SuggestedFix.Builder.createImportNode(ImportType, String, String)
func SuggestedFix_Builder_CreateImportNode(
	importType SuggestedFix_ImportType, alias *string, namespace string) *rhino.Node {
	var requireFlavor string
	switch importType {
	case SuggestedFix_ImportType_REQUIRE:
		requireFlavor = "require"
	case SuggestedFix_ImportType_REQUIRE_TYPE:
		requireFlavor = "requireType"
	}

	callNode := rhino.IR_Call(rhino.IR_Getprop(rhino.IR_Name("goog"), requireFlavor), rhino.IR_String(namespace))

	if alias != nil {
		return rhino.IR_ConstNode(rhino.IR_Name(*alias), callNode)
	} else {
		return rhino.IR_ExprResult(callNode)
	}
}

// Java: SuggestedFix.Builder.addGoogRequire(Match, String, ScriptMetadata)
func (b *SuggestedFix_Builder) AddGoogRequire(m *Match, namespace string, scriptMetadata *ScriptMetadata) *SuggestedFix_Builder {
	return b.AddImport(m, namespace, SuggestedFix_ImportType_REQUIRE, scriptMetadata)
}

// Java: SuggestedFix.Builder.addGoogRequireType(Match, String, ScriptMetadata)
func (b *SuggestedFix_Builder) AddGoogRequireType(m *Match, namespace string, scriptMetadata *ScriptMetadata) *SuggestedFix_Builder {
	return b.AddImport(m, namespace, SuggestedFix_ImportType_REQUIRE_TYPE, scriptMetadata)
}

// AddImport adds a goog.require/requireType for the given namespace if it
// does not already exist.
//
// Java: SuggestedFix.Builder.addImport(Match, String, ImportType, ScriptMetadata)
func (b *SuggestedFix_Builder) AddImport(
	m *Match, namespace string, importType SuggestedFix_ImportType, scriptMetadata *ScriptMetadata) *SuggestedFix_Builder {
	var alias *string
	if scriptMetadata.SupportsRequireAliases() {
		existingAlias := scriptMetadata.GetAlias(namespace)
		if existingAlias != nil {
			/*
			 * Each fix muct be independently valid, so go through the steps of adding a require even
			 * if one may already exist or have been added by another fix.
			 */
			alias = existingAlias
		} else if strings.IndexByte(namespace, '.') == -1 {
			/*
			 * For unqualified names, the exisiting references will still be valid so long as we keep
			 * the same name for the alias.
			 */
			alias = &namespace
		} else {
			// stream(RequireAliasGenerator.over(namespace))
			//     .filter((a) -> !scriptMetadata.usesName(a))
			//     .findFirst()
			//     .orElseThrow(AssertionError::new)
			for a := range RequireAliasGenerator_Over(namespace).All() {
				if !scriptMetadata.UsesName(a) {
					alias = &a
					break
				}
			}
			if alias == nil {
				panic(jutil.NewAssertionErrorNoMessage())
			}
		}
		scriptMetadata.AddAlias(namespace, *alias)
	} else {
		alias = nil
	}

	metadata := m.GetMetadata()
	existingNode := SuggestedFix_Builder_FindGoogRequireNode(m.GetNode(), metadata, namespace)

	if existingNode != nil {
		// TODO(b/139953612): Destructured goog.requires are not supported.

		// Add an alias to a naked require if allowed in this file.
		if existingNode.IsExprResult() && alias != nil {
			var newNode *rhino.Node
			// Replace goog.forwardDeclare with the appropriate alternative
			if NodeUtil_IsCallTo(existingNode.GetFirstChild(), "goog.forwardDeclare") {
				newNode = SuggestedFix_Builder_CreateImportNode(importType, alias, namespace)
			} else {
				newNode = rhino.IR_ConstNode(rhino.IR_Name(*alias), existingNode.GetFirstChild().CloneTree())
			}
			b.Replace(existingNode, newNode, m.GetMetadata().GetCompiler())
			scriptMetadata.AddAlias(namespace, *alias)
		}

		return b
	}

	// Find the right goog.require node to insert this after.
	script := scriptMetadata.GetScript()
	if script.GetFirstChild().IsModuleBody() {
		script = script.GetFirstChild()
	}

	var lastModuleOrProvideNode *rhino.Node
	var lastGoogRequireNode *rhino.Node
	var nodeToInsertBefore *rhino.Node
	child := script.GetFirstChild()
	for child != nil {
		if Matchers_GoogModule().Matches(child, metadata) {
			lastModuleOrProvideNode = child
		}
		if NodeUtil_IsExprCall(child) {
			// TODO(mknichel): Replace this logic with a function argument
			// Matcher when it exists.
			grandchild := child.GetFirstChild()
			if Matchers_GoogModuleOrProvide().Matches(grandchild, metadata) {
				lastModuleOrProvideNode = grandchild
			} else if Matchers_GoogRequirelike().Matches(grandchild, metadata) {
				lastGoogRequireNode = grandchild
				if grandchild.GetLastChild().IsStringLit() &&
					jutil.StringCompareTo(namespace, grandchild.GetLastChild().GetString()) < 0 {
					nodeToInsertBefore = child
					break
				}
			}
		} else if NodeUtil_IsNameDeclaration(child) &&
			child.GetFirstFirstChild() != nil &&
			Matchers_GoogRequirelike().Matches(child.GetFirstFirstChild(), metadata) {
			lastGoogRequireNode = child.GetFirstFirstChild()
			requireName := child.GetFirstChild().GetString()
			originalName := child.GetFirstChild().GetOriginalName()
			if originalName != nil {
				requireName = *originalName
			}
			if jutil.StringCompareTo(*alias, requireName) < 0 {
				nodeToInsertBefore = child
				break
			}
		}
		child = child.GetNext()
	}

	newImportNode := SuggestedFix_Builder_CreateImportNode(importType, alias, namespace)
	if nodeToInsertBefore == nil {
		// The file has goog.provide or goog.require nodes but they come before
		// the new goog.require node alphabetically.
		if lastModuleOrProvideNode != nil || lastGoogRequireNode != nil {
			var nodeToInsertAfter *rhino.Node
			if lastGoogRequireNode != nil {
				nodeToInsertAfter = lastGoogRequireNode
			} else {
				nodeToInsertAfter = lastModuleOrProvideNode
			}
			startPosition :=
				nodeToInsertAfter.GetSourceOffset() + nodeToInsertAfter.GetLength() + 2
			b.replacements.Put(
				*nodeToInsertAfter.GetSourceFileName(),
				CodeReplacement_CreateWithSortKey(
					startPosition,
					0,
					b.GenerateCode(m.GetMetadata().GetCompiler(), newImportNode),
					namespace))
			return b
		} else {
			// The file has no goog.provide or goog.require nodes.
			if script.HasChildren() {
				nodeToInsertBefore = script.GetFirstChild()
			} else {
				b.replacements.Put(
					*script.GetSourceFileName(),
					CodeReplacement_CreateWithSortKey(
						0, 0, b.GenerateCode(m.GetMetadata().GetCompiler(), newImportNode), namespace))
				return b
			}
		}
	}

	return b.InsertBeforeNodeWithSortKey(
		nodeToInsertBefore, newImportNode, m.GetMetadata().GetCompiler(), namespace)
}

// RemoveGoogRequire removes a goog.require for the given namespace to the
// file if it already exists.
//
// Java: SuggestedFix.Builder.removeGoogRequire(Match, String)
func (b *SuggestedFix_Builder) RemoveGoogRequire(m *Match, namespace string) *SuggestedFix_Builder {
	googRequireNode := SuggestedFix_Builder_FindGoogRequireNode(m.GetNode(), m.GetMetadata(), namespace)
	if googRequireNode != nil {
		return b.DeleteWithoutRemovingWhitespaceBefore(googRequireNode)
	}
	return b
}

// SuggestedFix_Builder_FindGoogRequireNode finds the goog.require node for
// the given namespace (or nil if there isn't one). If there is more than
// one:
//
//   - If there is at least one standalone goog.require, this will return the
//     first standalone goog.require.
//   - If not, this will return the first goog.require.
//
// Java: SuggestedFix.Builder.findGoogRequireNode(Node, NodeMetadata, String)
func SuggestedFix_Builder_FindGoogRequireNode(
	n *rhino.Node, metadata *NodeMetadata, namespace string) *rhino.Node {
	script := metadata.GetCompiler().GetScriptNode(*n.GetSourceFileName())
	if script.GetFirstChild().IsModuleBody() {
		script = script.GetFirstChild()
	}

	for child := script.GetFirstChild(); child != nil; child = child.GetNext() {
		if NodeUtil_IsExprCall(child) &&
			Matchers_GoogRequirelikeWithNamespace(namespace).Matches(child.GetFirstChild(), metadata) {
			return child
		}
	}

	for child := script.GetFirstChild(); child != nil; child = child.GetNext() {
		if NodeUtil_IsNameDeclaration(child) &&
			// TODO(b/139953612): respect destructured goog.requires
			!child.GetFirstChild().IsDestructuringLhs() &&
			child.GetFirstChild().GetLastChild() != nil &&
			Matchers_GoogRequirelikeWithNamespace(namespace).
				Matches(child.GetFirstChild().GetLastChild(), metadata) {
			return child
		}
	}
	return nil
}

// Java: SuggestedFix.Builder.subtreeRangeOfIdentifier(Node)
func SuggestedFix_Builder_SubtreeRangeOfIdentifier(n *rhino.Node) *rhino.Node {
	jutil.CheckStateMsg(n.IsGetProp(), "Support other identifier nodes")

	leftmost := n
	for leftmost.HasChildren() {
		leftmost = leftmost.GetFirstChild()
	}

	result := rhino.IR_Empty()
	result.SetStaticSourceFile(n.GetStaticSourceFile())
	result.SetLinenoCharno(leftmost.GetLineno(), leftmost.GetCharno())
	result.SetLength(n.GetLength() + (n.GetSourceOffset() - leftmost.GetSourceOffset()))
	return result
}

// Java: SuggestedFix.Builder.generateCode(AbstractCompiler, Node)
func (b *SuggestedFix_Builder) GenerateCode(compiler AbstractCompiler, node *rhino.Node) string {
	// TODO(mknichel): Fix all the formatting problems with this code.
	// How does this play with goog.scope?
	if node.IsBlock() {
		// Avoid printing the {}'s
		node.SetToken(rhino.Token_SCRIPT)
	}
	compilerOptions := NewCompilerOptions()
	compilerOptions.SetPreferSingleQuotes(true)
	compilerOptions.SetUseOriginalNamesInOutput(true)
	compilerOptions.SetPreserveNonJSDocComments(true)
	// We're refactoring existing code, so no need to escape values inside
	// strings.
	compilerOptions.SetTrustedStrings(true)
	return NewCodePrinter_Builder(node).
		SetCompilerOptions(compilerOptions).
		SetTypeRegistry(compiler.GetTypeRegistry()).
		SetPrettyPrint(true).
		SetLineBreak(true).
		SetOutputTypes(true).
		Build()
}

// Java: SuggestedFix.Builder.setDescription(String)
func (b *SuggestedFix_Builder) SetDescription(description string) *SuggestedFix_Builder {
	b.description = &description
	return b
}

// Java: SuggestedFix.Builder.build()
func (b *SuggestedFix_Builder) Build() *SuggestedFix {
	return newSuggestedFix(
		b.matchedNodeInfo, jutil.SetMultimapCopyOf(b.replacements), b.description, append([]*SuggestedFix{}, b.alternatives...))
}

// SuggestedFix_MatchedNodeInfo is information about the node that was
// matched for the suggested fix. This information can be used later on when
// processing the SuggestedFix.
//
// NOTE: Since this class can be retained for a long time when running
// refactorings over large blobs of code, it's important that it does not
// contain any memory intensive objects in order to keep memory to a
// reasonable amount.
//
// Java: public record SuggestedFix.MatchedNodeInfo(String sourceFilename, int lineno, int charno, boolean inClosurizedFile)
type SuggestedFix_MatchedNodeInfo struct {
	sourceFilename   string
	lineno           int
	charno           int
	inClosurizedFile bool
}

// Java: SuggestedFix.MatchedNodeInfo(String, int, int, boolean)
func NewSuggestedFix_MatchedNodeInfo(
	sourceFilename *string, lineno int, charno int, inClosurizedFile bool) *SuggestedFix_MatchedNodeInfo {
	// requireNonNull(sourceFilename, "sourceFilename")
	if sourceFilename == nil {
		panic(jutil.NewNullPointerExceptionMsg("sourceFilename"))
	}
	return &SuggestedFix_MatchedNodeInfo{
		sourceFilename:   *sourceFilename,
		lineno:           lineno,
		charno:           charno,
		inClosurizedFile: inClosurizedFile,
	}
}

// SourceFilename is the record accessor.
//
// Java: SuggestedFix.MatchedNodeInfo.sourceFilename()
func (i *SuggestedFix_MatchedNodeInfo) SourceFilename() string {
	return i.sourceFilename
}

// Lineno is the record accessor.
//
// Java: SuggestedFix.MatchedNodeInfo.lineno()
func (i *SuggestedFix_MatchedNodeInfo) Lineno() int {
	return i.lineno
}

// Charno is the record accessor.
//
// Java: SuggestedFix.MatchedNodeInfo.charno()
func (i *SuggestedFix_MatchedNodeInfo) Charno() int {
	return i.charno
}

// InClosurizedFile is the record accessor.
//
// Java: SuggestedFix.MatchedNodeInfo.inClosurizedFile()
func (i *SuggestedFix_MatchedNodeInfo) InClosurizedFile() bool {
	return i.inClosurizedFile
}

// Java: SuggestedFix.MatchedNodeInfo.getSourceFilename()
func (i *SuggestedFix_MatchedNodeInfo) GetSourceFilename() string {
	return i.SourceFilename()
}

// Java: SuggestedFix.MatchedNodeInfo.getLineno()
func (i *SuggestedFix_MatchedNodeInfo) GetLineno() int {
	return i.Lineno()
}

// Java: SuggestedFix.MatchedNodeInfo.getCharno()
func (i *SuggestedFix_MatchedNodeInfo) GetCharno() int {
	return i.Charno()
}

// Java: SuggestedFix.MatchedNodeInfo.isInClosurizedFile()
func (i *SuggestedFix_MatchedNodeInfo) IsInClosurizedFile() bool {
	return i.InClosurizedFile()
}

// Java: SuggestedFix.MatchedNodeInfo.create(Node, boolean)
func SuggestedFix_MatchedNodeInfo_Create(node *rhino.Node, closurized bool) *SuggestedFix_MatchedNodeInfo {
	return NewSuggestedFix_MatchedNodeInfo(
		NodeUtil_GetSourceName(node), node.GetLineno(), node.GetCharno(), closurized)
}

// SuggestedFix_GetStartPositionForNodeConsideringComments is a helper
// function to return the source offset of this node considering that JSDoc
// comments, non-JDDoc comments, or both may or may not be attached.
//
// Java: SuggestedFix.getStartPositionForNodeConsideringComments(Node)
func SuggestedFix_GetStartPositionForNodeConsideringComments(node *rhino.Node) int {
	jsdoc := NodeUtil_GetBestJSDocInfo(node)
	associatedNonJSDocComment := node.GetNonJSDocComment()
	start := node.GetSourceOffset()
	if jsdoc != nil {
		start = jsdoc.GetOriginalCommentPosition()
	}
	if associatedNonJSDocComment != nil {
		start = min(start, associatedNonJSDocComment.GetStartPosition().GetOffset())
	}
	return start
}
