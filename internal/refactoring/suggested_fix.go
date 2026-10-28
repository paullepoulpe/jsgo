// Port of com/google/javascript/refactoring/SuggestedFix.java @ closure-compiler v20261021

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

package refactoring

import (
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SuggestedFix represents the fixes to apply to source code to create a
// refactoring CL.
type SuggestedFix struct {
	matchedNodeInfo *SuggestedFix_MatchedNodeInfo
	// Multimap of filename to a modification to that file.
	replacements *jutil.SetMultimap[string, *CodeReplacement]
	description  *string
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
	f := &SuggestedFix{
		matchedNodeInfo: matchedNodeInfo,
		replacements:    replacements,
		description:     description,
	}
	all := make([]*SuggestedFix, 0, 1+len(alternatives))
	all = append(all, f)
	all = append(all, alternatives...)
	f.alternatives = all
	return f
}

// Java: SuggestedFix.getMatchedNodeInfo()
func (f *SuggestedFix) GetMatchedNodeInfo() *SuggestedFix_MatchedNodeInfo {
	return f.matchedNodeInfo
}

// Java: SuggestedFix.getReplacements()
func (f *SuggestedFix) GetReplacements() *jutil.SetMultimap[string, *CodeReplacement] {
	return f.replacements
}

// Java: SuggestedFix.getDescription()
func (f *SuggestedFix) GetDescription() *string {
	return f.description
}

// Java: SuggestedFix.getAlternatives()
func (f *SuggestedFix) GetAlternatives() []*SuggestedFix {
	return f.alternatives
}

// Java: SuggestedFix.getNonDefaultAlternatives()
func (f *SuggestedFix) GetNonDefaultAlternatives() []*SuggestedFix {
	return f.alternatives[1:]
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

// SuggestedFix_Builder builds SuggestedFix instances and provides AST edit
// helpers.
type SuggestedFix_Builder struct {
	matchedNodeInfo *SuggestedFix_MatchedNodeInfo
	replacements    *jutil.SetMultimap[string, *CodeReplacement]
	alternatives    []*SuggestedFix
	description     *string
}

// Java: SuggestedFix.Builder()
func NewSuggestedFix_Builder() *SuggestedFix_Builder {
	return &SuggestedFix_Builder{
		replacements: jutil.NewSetMultimap[string, *CodeReplacement](),
	}
}

// Java: SuggestedFix.Builder.attachMatchedNodeInfo(Node, AbstractCompiler)
func (b *SuggestedFix_Builder) AttachMatchedNodeInfo(node rhino.Node, compiler jscomp.AbstractCompiler) *SuggestedFix_Builder {
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

// Java: SuggestedFix.Builder.replaceText(Node, int, String)
func (b *SuggestedFix_Builder) ReplaceText(node rhino.Node, length int, newContent string) *SuggestedFix_Builder {
	b.replacements.Put(
		*node.GetSourceFileName(), CodeReplacement_Create(node.GetSourceOffset(), length, newContent))
	return b
}

// Java: SuggestedFix.Builder.addChildToFront(Node, String)
func (b *SuggestedFix_Builder) AddChildToFront(parentNode rhino.Node, content string) *SuggestedFix_Builder {
	jutil.CheckStateMsg(
		parentNode.IsBlock(), "addChildToFront is only supported for BLOCK statements.")
	startPosition := parentNode.GetSourceOffset() + 1
	b.replacements.Put(
		*parentNode.GetSourceFileName(), CodeReplacement_Create(startPosition, 0, "\n"+content))
	return b
}

// Java: SuggestedFix.Builder.insertAfter(Node, String)
func (b *SuggestedFix_Builder) InsertAfter(node rhino.Node, text string) *SuggestedFix_Builder {
	position := node.GetSourceOffset() + node.GetLength()
	b.replacements.Put(*node.GetSourceFileName(), CodeReplacement_Create(position, 0, text))
	return b
}

// Java: SuggestedFix.Builder.insertBefore(Node, Node, AbstractCompiler)
func (b *SuggestedFix_Builder) InsertBeforeNode(nodeToInsertBefore rhino.Node, n rhino.Node, compiler jscomp.AbstractCompiler) *SuggestedFix_Builder {
	return b.InsertBeforeNodeWithSortKey(nodeToInsertBefore, n, compiler, "")
}

// Java: SuggestedFix.Builder.insertBefore(Node, Node, AbstractCompiler, String)
func (b *SuggestedFix_Builder) InsertBeforeNodeWithSortKey(
	nodeToInsertBefore rhino.Node, n rhino.Node, compiler jscomp.AbstractCompiler, sortKey string) *SuggestedFix_Builder {
	return b.InsertBeforeWithSortKey(nodeToInsertBefore, b.GenerateCode(compiler, n), sortKey)
}

// Java: SuggestedFix.Builder.insertBefore(Node, String)
func (b *SuggestedFix_Builder) InsertBefore(nodeToInsertBefore rhino.Node, content string) *SuggestedFix_Builder {
	return b.InsertBeforeWithSortKey(nodeToInsertBefore, content, "")
}

// Java: SuggestedFix.Builder.insertBefore(Node, String, String)
func (b *SuggestedFix_Builder) InsertBeforeWithSortKey(nodeToInsertBefore rhino.Node, content string, sortKey string) *SuggestedFix_Builder {
	startPosition := SuggestedFix_GetStartPositionForNodeConsideringComments(nodeToInsertBefore)
	jutil.CheckNotNullMsg1(nodeToInsertBefore.GetSourceFileName(),
		"No source file name for node: %s", nodeToInsertBefore)
	b.replacements.Put(
		*nodeToInsertBefore.GetSourceFileName(),
		CodeReplacement_CreateWithSortKey(startPosition, 0, content, sortKey))
	return b
}

// Java: SuggestedFix.Builder.delete(Node)
func (b *SuggestedFix_Builder) Delete(n rhino.Node) *SuggestedFix_Builder {
	return b.DeleteWithWhitespaceBefore(n, true)
}

// Java: SuggestedFix.Builder.delete(Node, boolean)
func (b *SuggestedFix_Builder) DeleteWithWhitespaceBefore(n rhino.Node, deleteWhitespaceBefore bool) *SuggestedFix_Builder {
	startPosition := SuggestedFix_GetStartPositionForNodeConsideringComments(n)
	startOffsetWithoutComments := n.GetSourceOffset()
	length := (startOffsetWithoutComments - startPosition) + n.GetLength()

	next := n.GetNext()
	if !next.IsNil() &&
		jscomp.NodeUtil_GetBestJSDocInfo(next) == nil &&
		next.GetNonJSDocComment() == nil {
		length = next.GetSourceOffset() - startPosition
	}

	parent := n.GetParent()
	_, isNameInDecl := ast.AsNameDeclaration(parent)
	isNameInDecl = isNameInDecl && n.IsName()
	if isNameInDecl || n.IsStringKey() {
		if !next.IsNil() {
			length = SuggestedFix_GetStartPositionForNodeConsideringComments(next) - startPosition
		} else if child := n.GetFirstChild(); !child.IsNil() {
			length = (child.GetSourceOffset() + child.GetLength()) - startPosition
		}
		if parent.GetLastChild() == n && n != parent.GetFirstChild() {
			previousSibling := n.GetPrevious()
			prevEndNode := previousSibling
			if child := previousSibling.GetFirstChild(); !child.IsNil() {
				prevEndNode = child
			}
			startPositionDiff := startPosition - (prevEndNode.GetSourceOffset() + prevEndNode.GetLength())
			startPosition -= startPositionDiff
			length += startPositionDiff
		}
	}

	if deleteWhitespaceBefore && !parent.IsNil() && (parent.IsScript() || parent.IsBlock()) {
		if previousSibling := n.GetPrevious(); !previousSibling.IsNil() {
			previousSiblingEndPosition := previousSibling.GetSourceOffset() + previousSibling.GetLength()
			length += startPosition - previousSiblingEndPosition
			startPosition = previousSiblingEndPosition
		}
	}
	b.replacements.Put(*n.GetSourceFileName(), CodeReplacement_Create(startPosition, length, ""))
	return b
}

// Java: SuggestedFix.Builder.deleteWithoutRemovingWhitespaceBefore(Node)
func (b *SuggestedFix_Builder) DeleteWithoutRemovingWhitespaceBefore(n rhino.Node) *SuggestedFix_Builder {
	return b.DeleteWithWhitespaceBefore(n, false)
}

// Java: SuggestedFix.Builder.deleteWithoutRemovingWhitespace(Node)
func (b *SuggestedFix_Builder) DeleteWithoutRemovingWhitespace(n rhino.Node) *SuggestedFix_Builder {
	b.replacements.Put(
		*n.GetSourceFileName(), CodeReplacement_Create(n.GetSourceOffset(), n.GetLength(), ""))
	return b
}

// Java: SuggestedFix.Builder.rename(Node, String)
func (b *SuggestedFix_Builder) Rename(n rhino.Node, name string) *SuggestedFix_Builder {
	return b.RenameWithReplaceNameSubtree(n, name, false)
}

// Java: SuggestedFix.Builder.rename(Node, String, boolean)
func (b *SuggestedFix_Builder) RenameWithReplaceNameSubtree(n rhino.Node, name string, replaceNameSubtree bool) *SuggestedFix_Builder {
	var range_ rhino.Node
	switch n.GetToken() {
	case rhino.Token_CALL:
		return b.RenameWithReplaceNameSubtree(ast.MustCall(n).Callee().Node(), name, replaceNameSubtree)
	case rhino.Token_TAGGED_TEMPLATELIT:
		return b.RenameWithReplaceNameSubtree(ast.MustTaggedTemplateLit(n).Tag().Node(), name, replaceNameSubtree)
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

// Java: SuggestedFix.Builder.replaceRange(Node, Node, String)
func (b *SuggestedFix_Builder) ReplaceRange(first rhino.Node, last rhino.Node, newContent string) *SuggestedFix_Builder {
	jutil.CheckState(first.GetParent() == last.GetParent())

	start := SuggestedFix_GetStartPositionForNodeConsideringComments(first)
	if start == 0 {
		start = first.GetSourceOffset()
	}
	end := last.GetSourceOffset() + last.GetLength()
	b.replacements.Put(
		*first.GetSourceFileName(), CodeReplacement_Create(start, end-start, newContent))
	return b
}

// Java: SuggestedFix.Builder.replace(Node, Node, AbstractCompiler)
func (b *SuggestedFix_Builder) Replace(original rhino.Node, newNode rhino.Node, compiler jscomp.AbstractCompiler) *SuggestedFix_Builder {
	parent := original.GetParent()
	if !parent.IsNil() && parent.IsExprResult() {
		original = parent
	}
	newCode := strings.TrimSuffix(b.GenerateCode(compiler, newNode), "\n")

	needsSemicolon :=
		!parent.IsNil() &&
			(parent.IsExprResult() ||
				parent.IsBlock() ||
				parent.IsScript() ||
				parent.IsModuleBody())
	if !needsSemicolon {
		newCode = strings.TrimSuffix(newCode, ";")
	}

	if !parent.IsNil() && rhino.IR_MayBeExpression(parent) {
		replacement := newNode
		for (replacement.IsBlock() || replacement.IsScript() || replacement.IsModuleBody()) &&
			replacement.HasOneChild() {
			replacement = replacement.GetOnlyChild()
		}
		if exprStmt, ok := ast.AsExprStmt(replacement); ok {
			replacement = exprStmt.Expr().Node()
		}
		if rhino.IR_MayBeExpression(replacement) {
			outer := jscomp.NodeUtil_Precedence(parent.GetToken())
			inner := jscomp.NodeUtil_Precedence(original.GetToken())
			newInner := jscomp.NodeUtil_Precedence(replacement.GetToken())
			if newInner < jscomp.NodeUtil_Precedence(rhino.Token_CALL) && newInner <= outer && inner >= outer {
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

// Java: SuggestedFix.Builder.addCast(Node, AbstractCompiler, String)
func (b *SuggestedFix_Builder) AddCast(n rhino.Node, compiler jscomp.AbstractCompiler, type_ string) *SuggestedFix_Builder {
	b.replacements.Put(
		*n.GetSourceFileName(),
		CodeReplacement_Create(
			n.GetSourceOffset(),
			n.GetLength(),
			"/** @type {"+type_+"} */ ("+b.GenerateCode(compiler, n)+")"))
	return b
}

// Java: SuggestedFix.Builder.removeCast(Node, AbstractCompiler)
func (b *SuggestedFix_Builder) RemoveCast(n rhino.Node, compiler jscomp.AbstractCompiler) *SuggestedFix_Builder {
	cast, ok := ast.AsCast(n)
	jutil.CheckArgument(ok)
	jsDoc := n.GetJSDocInfo()
	b.replacements.Put(
		*n.GetSourceFileName(),
		CodeReplacement_Create(
			jsDoc.GetOriginalCommentPosition(),
			cast.Expr().Node().GetSourceOffset()-jsDoc.GetOriginalCommentPosition(),
			""))
	b.replacements.Put(
		*n.GetSourceFileName(),
		CodeReplacement_Create(n.GetSourceOffset()+n.GetLength()-1, 1, ""))
	return b
}

// Java: SuggestedFix.Builder.addOrReplaceJsDoc(Node, String)
func (b *SuggestedFix_Builder) AddOrReplaceJsDoc(n rhino.Node, newJsDoc string) *SuggestedFix_Builder {
	offset := n.GetSourceOffset()
	length := 0
	if n.IsGetProp() {
		offset = SuggestedFix_Builder_SubtreeRangeOfIdentifier(n).GetSourceOffset()
	}
	if jsDoc := jscomp.NodeUtil_GetBestJSDocInfo(n); jsDoc != nil {
		offset = jsDoc.GetOriginalCommentPosition()
		length = jutil.StringLength(*jsDoc.GetOriginalCommentString())
	}
	b.replacements.Put(*n.GetSourceFileName(), CodeReplacement_Create(offset, length, newJsDoc))
	return b
}

// Java: SuggestedFix.Builder.changeJsDocType(Node, AbstractCompiler, String)
func (b *SuggestedFix_Builder) ChangeJsDocType(n rhino.Node, compiler jscomp.AbstractCompiler, type_ string) *SuggestedFix_Builder {
	typeNode := parsing.JsDocInfoParser_ParseTypeString(type_)
	rhino.CheckNotNullMsg1(typeNode, "Invalid type: %s", type_)
	typeExpr := rhino.NewJSTypeExpression(typeNode, "jsflume")
	newJsType := typeExpr.Evaluate(nil, compiler.GetTypeRegistry())
	if newJsType == nil {
		panic(jutil.NewRuntimeException("JS Compiler does not recognize type: " + type_))
	}

	b.ReplaceTypePattern(n, type_, jutil.Pattern_Compile("@(type) *\\{?[^@\\s}]+\\}?"))
	b.ReplaceTypePattern(n, type_, jutil.Pattern_Compile("@(export|package|private|protected|public|const|return) *\\{[^}]+\\}"))
	return b
}

// Java: SuggestedFix.Builder.replaceTypePattern(Node, String, Pattern)
func (b *SuggestedFix_Builder) ReplaceTypePattern(n rhino.Node, type_ string, pattern *jutil.Pattern) {
	info := jscomp.NodeUtil_GetBestJSDocInfo(n)
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

// Java: SuggestedFix.Builder.insertArguments(Node, int, String...)
func (b *SuggestedFix_Builder) InsertArguments(n rhino.Node, position int, args ...string) *SuggestedFix_Builder {
	call, ok := ast.AsCall(n)
	jutil.CheckArgumentMsg(ok, "insertArguments is only applicable to function call nodes.")
	numArgs := call.NumArgs()
	jutil.CheckArgumentMsg(
		position >= 0 && position <= numArgs,
		"The specified position must be less than the number of arguments.")

	var startPosition int
	var argument rhino.Node
	if position < numArgs {
		arg, _ := call.AsCallOrNew().ArgAt(position)
		argument = arg.Node()
		startPosition = SuggestedFix_GetStartPositionForNodeConsideringComments(argument)
	} else {
		startPosition = n.GetSourceOffset() + n.GetLength() - 1
	}

	newContent := strings.Join(args, ", ")
	if !argument.IsNil() {
		newContent += ", "
	} else if position > 0 {
		newContent = ", " + newContent
	}
	b.replacements.Put(*n.GetSourceFileName(), CodeReplacement_Create(startPosition, 0, newContent))
	return b
}

// Java: SuggestedFix.Builder.deleteArgument(Node, int)
func (b *SuggestedFix_Builder) DeleteArgument(n rhino.Node, position int) *SuggestedFix_Builder {
	jutil.CheckArgumentMsg(
		n.IsCall() || n.IsNew(), "deleteArgument is only applicable to function call nodes.")
	callOrNew := ast.MustCallOrNew(n)
	numArguments := n.GetChildCount() - 1
	jutil.CheckStateMsg(
		numArguments > 0, "deleteArgument() cannot be used on a function call with no arguments")
	jutil.CheckArgumentMsg(
		position >= 0 && position < numArguments,
		"The specified position must be less than the number of arguments.")

	targetArg, _ := callOrNew.ArgAt(position)
	targetNode := targetArg.Node()

	var startOfArgumentToRemove, endOfArgumentToRemove int
	if position == 0 {
		startOfArgumentToRemove = SuggestedFix_GetStartPositionForNodeConsideringComments(targetNode)
		if nextArg, ok := callOrNew.ArgAt(1); ok {
			endOfArgumentToRemove = nextArg.Node().GetSourceOffset()
		} else {
			endOfArgumentToRemove = targetNode.GetSourceOffset() + targetNode.GetLength()
		}
	} else {
		prevArg, _ := callOrNew.ArgAt(position - 1)
		prevNode := prevArg.Node()
		startOfArgumentToRemove = prevNode.GetSourceOffset() + prevNode.GetLength()
		endOfArgumentToRemove = targetNode.GetSourceOffset() + targetNode.GetLength()
	}

	b.replacements.Put(
		*n.GetSourceFileName(),
		CodeReplacement_Create(startOfArgumentToRemove, endOfArgumentToRemove-startOfArgumentToRemove, ""))
	return b
}

// Java: SuggestedFix.Builder.createImportNode(ImportType, String, String)
func SuggestedFix_Builder_CreateImportNode(
	importType SuggestedFix_ImportType, alias *string, namespace string) rhino.Node {
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
	}
	return rhino.IR_ExprResult(callNode)
}

// Java: SuggestedFix.Builder.addGoogRequire(Match, String, ScriptMetadata)
func (b *SuggestedFix_Builder) AddGoogRequire(m *Match, namespace string, scriptMetadata *ScriptMetadata) *SuggestedFix_Builder {
	return b.AddImport(m, namespace, SuggestedFix_ImportType_REQUIRE, scriptMetadata)
}

// Java: SuggestedFix.Builder.addGoogRequireType(Match, String, ScriptMetadata)
func (b *SuggestedFix_Builder) AddGoogRequireType(m *Match, namespace string, scriptMetadata *ScriptMetadata) *SuggestedFix_Builder {
	return b.AddImport(m, namespace, SuggestedFix_ImportType_REQUIRE_TYPE, scriptMetadata)
}

// Java: SuggestedFix.Builder.addImport(Match, String, ImportType, ScriptMetadata)
func (b *SuggestedFix_Builder) AddImport(
	m *Match, namespace string, importType SuggestedFix_ImportType, scriptMetadata *ScriptMetadata) *SuggestedFix_Builder {
	var alias *string
	if scriptMetadata.SupportsRequireAliases() {
		if existingAlias := scriptMetadata.GetAlias(namespace); existingAlias != nil {
			alias = existingAlias
		} else if strings.IndexByte(namespace, '.') == -1 {
			alias = &namespace
		} else {
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
	}

	metadata := m.GetMetadata()
	existingNode := SuggestedFix_Builder_FindGoogRequireNode(m.GetNode(), metadata, namespace)

	if !existingNode.IsNil() {
		if exprStmt, ok := ast.AsExprStmt(existingNode); ok && alias != nil {
			expr := exprStmt.Expr().Node()
			var newNode rhino.Node
			if jscomp.NodeUtil_IsCallTo(expr, "goog.forwardDeclare") {
				newNode = SuggestedFix_Builder_CreateImportNode(importType, alias, namespace)
			} else {
				newNode = rhino.IR_ConstNode(rhino.IR_Name(*alias), expr.CloneTree())
			}
			b.Replace(existingNode, newNode, metadata.GetCompiler())
			scriptMetadata.AddAlias(namespace, *alias)
		}
		return b
	}

	script := scriptMetadata.GetScript()
	if script.GetFirstChild().IsModuleBody() {
		script = script.GetFirstChild()
	}

	var lastModuleOrProvideNode rhino.Node
	var lastGoogRequireNode rhino.Node
	var nodeToInsertBefore rhino.Node
	modMatcher := Matchers_GoogModule()
	modOrProvideMatcher := Matchers_GoogModuleOrProvide()
	reqLikeMatcher := Matchers_GoogRequirelike()

	for child := script.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if modMatcher.Matches(child, metadata) {
			lastModuleOrProvideNode = child
		}
		if exprStmt, ok := ast.AsExprStmt(child); ok {
			if call, isCall := ast.AsCall(exprStmt.Expr().Node()); isCall {
				callNode := call.Node()
				if modOrProvideMatcher.Matches(callNode, metadata) {
					lastModuleOrProvideNode = callNode
				} else if reqLikeMatcher.Matches(callNode, metadata) {
					lastGoogRequireNode = callNode
					if lastArg, hasArg := call.LastArg(); hasArg {
						if strLit, isStr := lastArg.AsStringLit(); isStr &&
							jutil.StringCompareTo(namespace, strLit.Value()) < 0 {
							nodeToInsertBefore = child
							break
						}
					}
				}
			}
		} else if decl, ok := ast.AsNameDeclaration(child); ok {
			if firstDecl, hasDecl := decl.FirstDeclarator(); hasDecl {
				if declName, isName := ast.AsName(firstDecl.Node()); isName {
					if initExpr, hasInit := declName.Init(); hasInit &&
						reqLikeMatcher.Matches(initExpr.Node(), metadata) {
						lastGoogRequireNode = initExpr.Node()
						requireName := declName.Ident()
						if originalName := declName.Node().GetOriginalName(); originalName != nil {
							requireName = *originalName
						}
						if jutil.StringCompareTo(*alias, requireName) < 0 {
							nodeToInsertBefore = child
							break
						}
					}
				}
			}
		}
	}

	newImportNode := SuggestedFix_Builder_CreateImportNode(importType, alias, namespace)
	if nodeToInsertBefore.IsNil() {
		if !lastModuleOrProvideNode.IsNil() || !lastGoogRequireNode.IsNil() {
			nodeToInsertAfter := lastModuleOrProvideNode
			if !lastGoogRequireNode.IsNil() {
				nodeToInsertAfter = lastGoogRequireNode
			}
			startPosition := nodeToInsertAfter.GetSourceOffset() + nodeToInsertAfter.GetLength() + 2
			b.replacements.Put(
				*nodeToInsertAfter.GetSourceFileName(),
				CodeReplacement_CreateWithSortKey(
					startPosition,
					0,
					b.GenerateCode(metadata.GetCompiler(), newImportNode),
					namespace))
			return b
		}
		if script.HasChildren() {
			nodeToInsertBefore = script.GetFirstChild()
		} else {
			b.replacements.Put(
				*script.GetSourceFileName(),
				CodeReplacement_CreateWithSortKey(
					0, 0, b.GenerateCode(metadata.GetCompiler(), newImportNode), namespace))
			return b
		}
	}

	return b.InsertBeforeNodeWithSortKey(
		nodeToInsertBefore, newImportNode, metadata.GetCompiler(), namespace)
}

// Java: SuggestedFix.Builder.removeGoogRequire(Match, String)
func (b *SuggestedFix_Builder) RemoveGoogRequire(m *Match, namespace string) *SuggestedFix_Builder {
	googRequireNode := SuggestedFix_Builder_FindGoogRequireNode(m.GetNode(), m.GetMetadata(), namespace)
	if !googRequireNode.IsNil() {
		return b.DeleteWithoutRemovingWhitespaceBefore(googRequireNode)
	}
	return b
}

// Java: SuggestedFix.Builder.findGoogRequireNode(Node, NodeMetadata, String)
func SuggestedFix_Builder_FindGoogRequireNode(
	n rhino.Node, metadata *NodeMetadata, namespace string) rhino.Node {
	script := metadata.GetCompiler().GetScriptNode(*n.GetSourceFileName())
	if script.GetFirstChild().IsModuleBody() {
		script = script.GetFirstChild()
	}

	nsMatcher := Matchers_GoogRequirelikeWithNamespace(namespace)
	for child := script.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if exprStmt, ok := ast.AsExprStmt(child); ok {
			expr := exprStmt.Expr().Node()
			if expr.IsCall() && nsMatcher.Matches(expr, metadata) {
				return child
			}
		}
	}

	for child := script.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if decl, ok := ast.AsNameDeclaration(child); ok {
			if firstDecl, hasDecl := decl.FirstDeclarator(); hasDecl {
				if declName, isName := ast.AsName(firstDecl.Node()); isName {
					if initExpr, hasInit := declName.Init(); hasInit &&
						nsMatcher.Matches(initExpr.Node(), metadata) {
						return child
					}
				}
			}
		}
	}
	return rhino.NilNode
}

// Java: SuggestedFix.Builder.subtreeRangeOfIdentifier(Node)
func SuggestedFix_Builder_SubtreeRangeOfIdentifier(n rhino.Node) rhino.Node {
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
func (b *SuggestedFix_Builder) GenerateCode(compiler jscomp.AbstractCompiler, node rhino.Node) string {
	if node.IsBlock() {
		node.SetToken(rhino.Token_SCRIPT)
	}
	compilerOptions := jscomp.NewCompilerOptions()
	compilerOptions.SetPreferSingleQuotes(true)
	compilerOptions.SetUseOriginalNamesInOutput(true)
	compilerOptions.SetPreserveNonJSDocComments(true)
	compilerOptions.SetTrustedStrings(true)
	return jscomp.NewCodePrinter_Builder(node).
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
		b.matchedNodeInfo, jutil.SetMultimapCopyOf(b.replacements), b.description, slices.Clone(b.alternatives))
}

// SuggestedFix_MatchedNodeInfo holds source location metadata for the matched node.
type SuggestedFix_MatchedNodeInfo struct {
	sourceFilename   string
	lineno           int
	charno           int
	inClosurizedFile bool
}

// Java: SuggestedFix.MatchedNodeInfo(String, int, int, boolean)
func NewSuggestedFix_MatchedNodeInfo(
	sourceFilename *string, lineno int, charno int, inClosurizedFile bool) *SuggestedFix_MatchedNodeInfo {
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

func (i *SuggestedFix_MatchedNodeInfo) SourceFilename() string    { return i.sourceFilename }
func (i *SuggestedFix_MatchedNodeInfo) Lineno() int               { return i.lineno }
func (i *SuggestedFix_MatchedNodeInfo) Charno() int               { return i.charno }
func (i *SuggestedFix_MatchedNodeInfo) InClosurizedFile() bool    { return i.inClosurizedFile }
func (i *SuggestedFix_MatchedNodeInfo) GetSourceFilename() string { return i.sourceFilename }
func (i *SuggestedFix_MatchedNodeInfo) GetLineno() int            { return i.lineno }
func (i *SuggestedFix_MatchedNodeInfo) GetCharno() int            { return i.charno }
func (i *SuggestedFix_MatchedNodeInfo) IsInClosurizedFile() bool  { return i.inClosurizedFile }

// Java: SuggestedFix.MatchedNodeInfo.create(Node, boolean)
func SuggestedFix_MatchedNodeInfo_Create(node rhino.Node, closurized bool) *SuggestedFix_MatchedNodeInfo {
	return NewSuggestedFix_MatchedNodeInfo(
		jscomp.NodeUtil_GetSourceName(node), node.GetLineno(), node.GetCharno(), closurized)
}

// Java: SuggestedFix.getStartPositionForNodeConsideringComments(Node)
func SuggestedFix_GetStartPositionForNodeConsideringComments(node rhino.Node) int {
	start := node.GetSourceOffset()
	if jsdoc := jscomp.NodeUtil_GetBestJSDocInfo(node); jsdoc != nil {
		start = jsdoc.GetOriginalCommentPosition()
	}
	if associatedNonJSDocComment := node.GetNonJSDocComment(); associatedNonJSDocComment != nil {
		start = min(start, associatedNonJSDocComment.GetStartPosition().GetOffset())
	}
	return start
}
