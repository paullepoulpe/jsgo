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
// Port of com/google/javascript/refactoring/SuggestedFix.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_suggested_fix_1 j2gort.Ext_java_util_function_Supplier

func init() {
	lambdaSingleton_suggested_fix_1 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError"))
	})
}

// Java class: com.google.javascript.refactoring.SuggestedFix
type SuggestedFix struct {
	matchedNodeInfo *SuggestedFix_MatchedNodeInfo                    `gson:"matchedNodeInfo"`
	replacements    j2gort.Ext_com_google_common_collect_SetMultimap `gson:"replacements"`
	description     string                                           `gson:"description"`
	alternatives    []*SuggestedFix                                  `gson:"alternatives"`
	idHash_         atomic.Int32
}

func (this_ *SuggestedFix) BindSelf_(s any) {
	_ = s
}
func (this_ *SuggestedFix) IsJavaNil_() bool { return this_ == nil }
func (*SuggestedFix) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.SuggestedFix", "SuggestedFix"
}
func (this_ *SuggestedFix) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SuggestedFix) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.SuggestedFix", "com.google.javascript.refactoring.SuggestedFix", func(v_ any) bool {
		return j2gort.InstanceOf[*SuggestedFix](v_)
	})
}

func (this_ *SuggestedFix) initFields_SuggestedFix() {
	_ = this_
	this_.description = j2gort.NullString
}

// Java constructor: SuggestedFix(com.google.javascript.refactoring.SuggestedFix.MatchedNodeInfo matchedNodeInfo, com.google.common.collect.SetMultimap<java.lang.String,com.google.javascript.refactoring.CodeReplacement> replacements, java.lang.@org.jspecify.annotations.Nullable String description, com.google.common.collect.ImmutableList<com.google.javascript.refactoring.SuggestedFix> alternatives)
func NewSuggestedFix(matchedNodeInfo *SuggestedFix_MatchedNodeInfo, replacements j2gort.Ext_com_google_common_collect_SetMultimap, description string, alternatives j2gort.Ext_com_google_common_collect_ImmutableList) *SuggestedFix {
	this_ := &SuggestedFix{}
	InitSuggestedFix(this_, matchedNodeInfo, replacements, description, alternatives)
	return this_
}

func InitSuggestedFix(this_ *SuggestedFix, matchedNodeInfo *SuggestedFix_MatchedNodeInfo, replacements j2gort.Ext_com_google_common_collect_SetMultimap, description string, alternatives j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = matchedNodeInfo
	_ = replacements
	_ = description
	_ = alternatives
	this_.initFields_SuggestedFix()
	this_.matchedNodeInfo = matchedNodeInfo
	this_.replacements = replacements
	this_.description = description
	this_.alternatives = j2gort.ImmListToSliceNullable[*SuggestedFix]((&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}).AddSelf_(j2gort.PtrToAny(this_)).AddAllSelf_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](alternatives)).BuildImmList_())
}

// Java method: SuggestedFix.getMatchedNodeInfo()
func (this_ *SuggestedFix) GetMatchedNodeInfo() *SuggestedFix_MatchedNodeInfo {
	return this_.matchedNodeInfo
}

// Java method: SuggestedFix.getReplacements()
func (this_ *SuggestedFix) GetReplacements() j2gort.Ext_com_google_common_collect_SetMultimap {
	return this_.replacements
}

// Java method: SuggestedFix.getDescription()
func (this_ *SuggestedFix) GetDescription() string {
	return this_.description
}

// Java method: SuggestedFix.getAlternatives()
func (this_ *SuggestedFix) GetAlternatives() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListFromTypedSlice(this_.alternatives)
}

// Java method: SuggestedFix.getNonDefaultAlternatives()
func (this_ *SuggestedFix) GetNonDefaultAlternatives() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListFromTypedSlice(j2gort.ImmSliceSubList(this_.alternatives, int32(1), int32(len(this_.alternatives))))
}

// Java method: SuggestedFix.isNoOp()
func (this_ *SuggestedFix) IsNoOp() bool {
	return this_.replacements.IsEmpty()
}

// Java method: SuggestedFix.toString()
func (this_ *SuggestedFix) String() string {
	if this_.IsNoOp() {
		return "<no-op SuggestedFix>"
	}
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	{
		cur_1 := j2gort.NewMapEntryCursor(j2gort.AsIface[j2gort.Ext_java_util_Map](this_.replacements.AsMap()))
		for cur_1.NextEntry() {
			j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string("Replacements for file: ")).Append_1_string(cur_1.KeyStr())).Append_1_string("\n")
			j2gort.Ext_com_google_common_base_Joiner_On_1_string("\n\n").AppendTo_2_Ext_java_lang_StringBuilder_Ext_java_lang_Iterable(sb, j2gort.AsIface[j2gort.Ext_java_util_Collection](cur_1.GetValue()))
		}
	}
	return sb.String()
}

// Java static method: SuggestedFix.getStartPositionForNodeConsideringComments(com.google.javascript.rhino.Node node)
func SuggestedFix_GetStartPositionForNodeConsideringComments(node pkg_rhino.Node) int32 {
	_ = node
	var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(node)
	_ = jsdoc
	var associatedNonJSDocComment *pkg_rhino.NonJSDocComment = node.GetNonJSDocComment()
	_ = associatedNonJSDocComment
	var start int32 = node.GetSourceOffset()
	_ = start
	if jsdoc != nil {
		start = jsdoc.GetOriginalCommentPosition()
	}
	if associatedNonJSDocComment != nil {
		start = min(start, associatedNonJSDocComment.GetStartPosition().GetOffset())
	}
	return start
}

// Java enum: com.google.javascript.refactoring.SuggestedFix.ImportType
type SuggestedFix_ImportType uint8

const (
	enum_SuggestedFix_ImportType_REQUIRE      SuggestedFix_ImportType = 1
	enum_SuggestedFix_ImportType_REQUIRE_TYPE SuggestedFix_ImportType = 2
)

var enumBaseTable_SuggestedFix_ImportType = [...]j2gort.EnumBase{
	{Name_: "REQUIRE", Ordinal_: 0},
	{Name_: "REQUIRE_TYPE", Ordinal_: 1},
}
var enumValues_SuggestedFix_ImportType = []SuggestedFix_ImportType{1, 2}

func StaticGet_SuggestedFix_ImportType_REQUIRE() SuggestedFix_ImportType {
	return enum_SuggestedFix_ImportType_REQUIRE
}
func StaticGet_SuggestedFix_ImportType_REQUIRE_TYPE() SuggestedFix_ImportType {
	return enum_SuggestedFix_ImportType_REQUIRE_TYPE
}

func SuggestedFix_ImportType_Values() []SuggestedFix_ImportType {
	return j2gort.CloneSlice(enumValues_SuggestedFix_ImportType)
}

func SuggestedFix_ImportType_ValueOf(name string) SuggestedFix_ImportType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_SuggestedFix_ImportType {
		if enumBaseTable_SuggestedFix_ImportType[i_].Name_ == name {
			return SuggestedFix_ImportType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.refactoring.SuggestedFix.ImportType." + name))
}

func (this_ SuggestedFix_ImportType) IsJavaNil_() bool                              { return this_ == 0 }
func (SuggestedFix_ImportType) HasIdentityEquals_()                                 {}
func (SuggestedFix_ImportType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (SuggestedFix_ImportType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ SuggestedFix_ImportType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (SuggestedFix_ImportType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.SuggestedFix$ImportType", "ImportType"
}
func (this_ SuggestedFix_ImportType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_SuggestedFix_ImportType[this_-1]
}
func (this_ SuggestedFix_ImportType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ SuggestedFix_ImportType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ SuggestedFix_ImportType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_SuggestedFix_ImportType[this_-1].Name_
}
func (this_ SuggestedFix_ImportType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_SuggestedFix_ImportType[this_-1].Name_
}
func (this_ SuggestedFix_ImportType) Equals(other_ any) bool {
	o_, ok_ := other_.(SuggestedFix_ImportType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ SuggestedFix_ImportType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_SuggestedFix_ImportType[this_-1].HashCode()
}
func (this_ SuggestedFix_ImportType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_SuggestedFix_ImportType[this_-1].CompareTo(other_)
}
func (this_ SuggestedFix_ImportType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_SuggestedFix_ImportType[this_-1].Name_)
}
func init() {
	vals_ := []any{SuggestedFix_ImportType(1), SuggestedFix_ImportType(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.refactoring.SuggestedFix.ImportType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.refactoring.SuggestedFix$ImportType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.SuggestedFix.ImportType", "com.google.javascript.refactoring.SuggestedFix$ImportType", func(v_ any) bool {
		return j2gort.InstanceOf[SuggestedFix_ImportType](v_)
	})
}

// Java class: com.google.javascript.refactoring.SuggestedFix.Builder
type SuggestedFix_Builder struct {
	matchedNodeInfo *SuggestedFix_MatchedNodeInfo                                      `gson:"matchedNodeInfo"`
	replacements    *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap_Builder `gson:"replacements"`
	alternatives    *j2gort.Ext_com_google_common_collect_ImmutableList_Builder        `gson:"alternatives"`
	description     string                                                             `gson:"description"`
	idHash_         atomic.Int32
}

func (this_ *SuggestedFix_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *SuggestedFix_Builder) IsJavaNil_() bool { return this_ == nil }
func (*SuggestedFix_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.SuggestedFix$Builder", "Builder"
}
func (this_ *SuggestedFix_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SuggestedFix_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.SuggestedFix.Builder", "com.google.javascript.refactoring.SuggestedFix$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*SuggestedFix_Builder](v_)
	})
}

func (this_ *SuggestedFix_Builder) initFields_SuggestedFix_Builder() {
	_ = this_

	this_.replacements = j2gort.Ext_com_google_common_collect_ImmutableSetMultimap_S_Builder()
	this_.alternatives = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	this_.description = j2gort.NullString
}

// Java constructor: Builder()
func NewSuggestedFix_Builder() *SuggestedFix_Builder {
	this_ := &SuggestedFix_Builder{}
	InitSuggestedFix_Builder(this_)
	return this_
}

func InitSuggestedFix_Builder(this_ *SuggestedFix_Builder) {
	this_.initFields_SuggestedFix_Builder()
}

// Java method: Builder.attachMatchedNodeInfo(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *SuggestedFix_Builder) AttachMatchedNodeInfo(node pkg_rhino.Node, compiler AbstractCompiler) *SuggestedFix_Builder {
	_ = node
	_ = compiler
	this_.matchedNodeInfo = SuggestedFix_MatchedNodeInfo_Create(node, RefactoringUtils_IsInClosurizedFile(node, NewNodeMetadata_1(compiler)))
	return this_
}

// Java method: Builder.addAlternative(com.google.javascript.refactoring.SuggestedFix alternative)
func (this_ *SuggestedFix_Builder) AddAlternative(alternative *SuggestedFix) *SuggestedFix_Builder {
	_ = alternative
	if !(alternative.GetNonDefaultAlternatives().IsEmpty()) {
		j2gort.CheckStateStr(false, "Alternative SuggestedFix must have no alternatives of their own.")
	}
	this_.alternatives.AddVoid_(j2gort.PtrToAny(alternative))
	return this_
}

// Java method: Builder.replaceText(com.google.javascript.rhino.Node node, int length, java.lang.String newContent)
func (this_ *SuggestedFix_Builder) ReplaceText(node pkg_rhino.Node, length int32, newContent string) *SuggestedFix_Builder {
	_ = node
	_ = length
	_ = newContent
	var startPosition int32 = node.GetSourceOffset()
	_ = startPosition
	j2gort.MultimapBuilderPutStr(this_.replacements, node.GetSourceFileName(), CodeReplacement_Create_3(startPosition, length, newContent))
	return this_
}

// Java method: Builder.addChildToFront(com.google.javascript.rhino.Node parentNode, java.lang.String content)
func (this_ *SuggestedFix_Builder) AddChildToFront(parentNode pkg_rhino.Node, content string) *SuggestedFix_Builder {
	_ = parentNode
	_ = content
	if !(parentNode.IsBlock()) {
		j2gort.CheckStateStr(false, "addChildToFront is only supported for BLOCK statements.")
	}
	var startPosition int32 = (parentNode.GetSourceOffset() + int32(1))
	_ = startPosition
	j2gort.MultimapBuilderPutStr(this_.replacements, parentNode.GetSourceFileName(), CodeReplacement_Create_3(startPosition, int32(0), ("\n"+j2gort.StrValString(content))))
	return this_
}

// Java method: Builder.insertAfter(com.google.javascript.rhino.Node node, java.lang.String text)
func (this_ *SuggestedFix_Builder) InsertAfter(node pkg_rhino.Node, text string) *SuggestedFix_Builder {
	_ = node
	_ = text
	var position int32 = (node.GetSourceOffset() + node.GetLength())
	_ = position
	j2gort.MultimapBuilderPutStr(this_.replacements, node.GetSourceFileName(), CodeReplacement_Create_3(position, int32(0), text))
	return this_
}

// Java method: Builder.insertBefore(com.google.javascript.rhino.Node nodeToInsertBefore, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *SuggestedFix_Builder) InsertBefore_3_Node_Node_AbstractCompiler(nodeToInsertBefore pkg_rhino.Node, n pkg_rhino.Node, compiler AbstractCompiler) *SuggestedFix_Builder {
	_ = nodeToInsertBefore
	_ = n
	_ = compiler
	return this_.InsertBefore_4(nodeToInsertBefore, n, compiler, "")
}

// Java method: Builder.insertBefore(com.google.javascript.rhino.Node nodeToInsertBefore, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String sortKey)
func (this_ *SuggestedFix_Builder) InsertBefore_4(nodeToInsertBefore pkg_rhino.Node, n pkg_rhino.Node, compiler AbstractCompiler, sortKey string) *SuggestedFix_Builder {
	_ = nodeToInsertBefore
	_ = n
	_ = compiler
	_ = sortKey
	return this_.priv_Builder_insertBefore_7974(nodeToInsertBefore, this_.GenerateCode(compiler, n), sortKey)
}

// Java method: Builder.insertBefore(com.google.javascript.rhino.Node nodeToInsertBefore, java.lang.String content)
func (this_ *SuggestedFix_Builder) InsertBefore_2(nodeToInsertBefore pkg_rhino.Node, content string) *SuggestedFix_Builder {
	_ = nodeToInsertBefore
	_ = content
	return this_.priv_Builder_insertBefore_7974(nodeToInsertBefore, content, "")
}

// Java method: Builder.insertBefore(com.google.javascript.rhino.Node nodeToInsertBefore, java.lang.String content, java.lang.String sortKey)
func (this_ *SuggestedFix_Builder) priv_Builder_insertBefore_7974(nodeToInsertBefore pkg_rhino.Node, content string, sortKey string) *SuggestedFix_Builder {
	_ = nodeToInsertBefore
	_ = content
	_ = sortKey
	var startPosition int32 = SuggestedFix_GetStartPositionForNodeConsideringComments(nodeToInsertBefore)
	_ = startPosition
	j2gort.CheckNotNullString(nodeToInsertBefore.GetSourceFileName())
	j2gort.MultimapBuilderPutStr(this_.replacements, nodeToInsertBefore.GetSourceFileName(), CodeReplacement_Create_4(startPosition, int32(0), content, sortKey))
	return this_
}

// Java method: Builder.delete(com.google.javascript.rhino.Node n)
func (this_ *SuggestedFix_Builder) Delete_1(n pkg_rhino.Node) *SuggestedFix_Builder {
	_ = n
	return this_.priv_Builder_delete_7976(n, true)
}

// Java method: Builder.delete(com.google.javascript.rhino.Node n, boolean deleteWhitespaceBefore)
func (this_ *SuggestedFix_Builder) priv_Builder_delete_7976(n pkg_rhino.Node, deleteWhitespaceBefore bool) *SuggestedFix_Builder {
	_ = n
	_ = deleteWhitespaceBefore
	var startPosition int32 = SuggestedFix_GetStartPositionForNodeConsideringComments(n)
	_ = startPosition
	var startOffsetWithoutComments int32 = n.GetSourceOffset()
	_ = startOffsetWithoutComments
	var length int32 = ((startOffsetWithoutComments - startPosition) + n.GetLength())
	_ = length
	if ((n.GetNext() != 0) && (NodeUtil_GetBestJSDocInfo(n.GetNext()) == nil)) && (n.GetNext().GetNonJSDocComment() == nil) {
		length = (n.GetNext().GetSourceOffset() - startPosition)
	}
	if (n.IsName() && NodeUtil_IsNameDeclaration(n.GetParent())) || n.IsStringKey() {
		if n.GetNext() != 0 {
			length = (SuggestedFix_GetStartPositionForNodeConsideringComments(n.GetNext()) - startPosition)
		} else if n.HasChildren() {
			var child pkg_rhino.Node = n.GetFirstChild()
			_ = child
			length = ((child.GetSourceOffset() + child.GetLength()) - startPosition)
		}
		if (n.GetParent().GetLastChild() == n) && (n != n.GetParent().GetFirstChild()) {
			var previousSibling pkg_rhino.Node = n.GetPrevious()
			_ = previousSibling
			if previousSibling.HasChildren() {
				var child_2 pkg_rhino.Node = previousSibling.GetFirstChild()
				_ = child_2
				var startPositionDiff int32 = (startPosition - (child_2.GetSourceOffset() + child_2.GetLength()))
				_ = startPositionDiff
				startPosition = (startPosition - startPositionDiff)
				length = (length + startPositionDiff)
			} else {
				var startPositionDiff_2 int32 = (startPosition - (previousSibling.GetSourceOffset() + previousSibling.GetLength()))
				_ = startPositionDiff_2
				startPosition = (startPosition - startPositionDiff_2)
				length = (length + startPositionDiff_2)
			}
		}
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if (deleteWhitespaceBefore && (parent != 0)) && (parent.IsScript() || parent.IsBlock()) {
		var previousSibling_2 pkg_rhino.Node = n.GetPrevious()
		_ = previousSibling_2
		if previousSibling_2 != 0 {
			var previousSiblingEndPosition int32 = (previousSibling_2.GetSourceOffset() + previousSibling_2.GetLength())
			_ = previousSiblingEndPosition
			length = (length + (startPosition - previousSiblingEndPosition))
			startPosition = previousSiblingEndPosition
		}
	}
	j2gort.MultimapBuilderPutStr(this_.replacements, n.GetSourceFileName(), CodeReplacement_Create_3(startPosition, length, ""))
	return this_
}

// Java method: Builder.deleteWithoutRemovingWhitespaceBefore(com.google.javascript.rhino.Node n)
func (this_ *SuggestedFix_Builder) DeleteWithoutRemovingWhitespaceBefore(n pkg_rhino.Node) *SuggestedFix_Builder {
	_ = n
	return this_.priv_Builder_delete_7976(n, false)
}

// Java method: Builder.deleteWithoutRemovingWhitespace(com.google.javascript.rhino.Node n)
func (this_ *SuggestedFix_Builder) DeleteWithoutRemovingWhitespace(n pkg_rhino.Node) *SuggestedFix_Builder {
	_ = n
	j2gort.MultimapBuilderPutStr(this_.replacements, n.GetSourceFileName(), CodeReplacement_Create_3(n.GetSourceOffset(), n.GetLength(), ""))
	return this_
}

// Java method: Builder.rename(com.google.javascript.rhino.Node n, java.lang.String name)
func (this_ *SuggestedFix_Builder) Rename_2(n pkg_rhino.Node, name string) *SuggestedFix_Builder {
	_ = n
	_ = name
	return this_.Rename_3(n, name, false)
}

// Java method: Builder.rename(com.google.javascript.rhino.Node n, java.lang.String name, boolean replaceNameSubtree)
func (this_ *SuggestedFix_Builder) Rename_3(n pkg_rhino.Node, name string, replaceNameSubtree bool) *SuggestedFix_Builder {
	_ = n
	_ = name
	_ = replaceNameSubtree
	var range_ pkg_rhino.Node = pkg_rhino.Node(0)
	_ = range_
	switch n.TokenSwitchOrd_() {
	case int32(28) /* Token.CALL */, int32(130) /* Token.TAGGED_TEMPLATELIT */ :
		return this_.Rename_3(n.GetFirstChild(), name, replaceNameSubtree)
	case int32(26) /* Token.GETPROP */ :
		if replaceNameSubtree {
			range_ = SuggestedFix_Builder_SubtreeRangeOfIdentifier(n)
		} else {
			range_ = n
		}
	case int32(35) /* Token.STRINGLIT */ :
		if !(n.GetParent().IsGetProp()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
		}
		range_ = n
	case int32(100) /* Token.STRING_KEY */, int32(32) /* Token.NAME */ :
		range_ = n
	default:
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Rename is not implemented for this node type: "+j2gort.StrVal(n))))
	}
	j2gort.MultimapBuilderPutStr(this_.replacements, range_.GetSourceFileName(), CodeReplacement_Create_3(range_.GetSourceOffset(), range_.GetLength(), name))
	return this_
}

// Java method: Builder.replaceRange(com.google.javascript.rhino.Node first, com.google.javascript.rhino.Node last, java.lang.String newContent)
func (this_ *SuggestedFix_Builder) ReplaceRange(first pkg_rhino.Node, last pkg_rhino.Node, newContent string) *SuggestedFix_Builder {
	_ = first
	_ = last
	_ = newContent
	if !(first.GetParent() == last.GetParent()) {
		j2gort.CheckState(false)
	}
	var start int32 = SuggestedFix_GetStartPositionForNodeConsideringComments(first)
	_ = start
	if start == int32(0) {
		start = first.GetSourceOffset()
	}
	var end int32 = (last.GetSourceOffset() + last.GetLength())
	_ = end
	var length int32 = (end - start)
	_ = length
	j2gort.MultimapBuilderPutStr(this_.replacements, first.GetSourceFileName(), CodeReplacement_Create_3(start, length, newContent))
	return this_
}

// Java method: Builder.replace(com.google.javascript.rhino.Node original, com.google.javascript.rhino.Node newNode, com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *SuggestedFix_Builder) Replace(original pkg_rhino.Node, newNode pkg_rhino.Node, compiler AbstractCompiler) *SuggestedFix_Builder {
	_ = original
	_ = newNode
	_ = compiler
	var parent pkg_rhino.Node = original.GetParent()
	_ = parent
	if (parent != 0) && parent.IsExprResult() {
		original = parent
	}
	var newCode string = this_.GenerateCode(compiler, newNode)
	_ = newCode
	if javalang.StringEndsWith(j2gort.CheckNotNullString(newCode), "\n") {
		newCode = j2gort.JStrSub(newCode, int32(0), (j2gort.JStrLen(newCode) - int32(1)))
	}
	var needsSemicolon bool = ((parent != 0) && (((parent.IsExprResult() || parent.IsBlock()) || parent.IsScript()) || parent.IsModuleBody()))
	_ = needsSemicolon
	if javalang.StringEndsWith(j2gort.CheckNotNullString(newCode), ";") && !(needsSemicolon) {
		newCode = j2gort.JStrSub(newCode, int32(0), (j2gort.JStrLen(newCode) - int32(1)))
	}
	if (parent != 0) && pkg_rhino.IR_MayBeExpression(parent) {
		var replacement pkg_rhino.Node = newNode
		_ = replacement
		for ((replacement.IsBlock() || replacement.IsScript()) || replacement.IsModuleBody()) && replacement.HasOneChild() {
			replacement = replacement.GetOnlyChild()
		}
		if replacement.IsExprResult() {
			replacement = replacement.GetOnlyChild()
		}
		if pkg_rhino.IR_MayBeExpression(replacement) {
			var outer int32 = NodeUtil_Precedence(parent.GetToken())
			_ = outer
			var inner int32 = NodeUtil_Precedence(original.GetToken())
			_ = inner
			var newInner int32 = NodeUtil_Precedence(replacement.GetToken())
			_ = newInner
			if ((newInner < NodeUtil_Precedence(pkg_rhino.StaticGet_Token_CALL())) && (newInner <= outer)) && (inner >= outer) {
				newCode = ("(" + j2gort.StrValString(newCode) + ")")
			}
		}
	}
	var range_ pkg_rhino.Node = original
	_ = range_
	if original.IsGetProp() {
		range_ = SuggestedFix_Builder_SubtreeRangeOfIdentifier(original)
	}
	j2gort.MultimapBuilderPutStr(this_.replacements, range_.GetSourceFileName(), CodeReplacement_Create_3(range_.GetSourceOffset(), range_.GetLength(), newCode))
	return this_
}

// Java method: Builder.addCast(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String type)
func (this_ *SuggestedFix_Builder) AddCast(n pkg_rhino.Node, compiler AbstractCompiler, type_ string) *SuggestedFix_Builder {
	_ = n
	_ = compiler
	_ = type_
	j2gort.MultimapBuilderPutStr(this_.replacements, n.GetSourceFileName(), CodeReplacement_Create_3(n.GetSourceOffset(), n.GetLength(), ("/** @type {"+j2gort.StrValString(type_)+"} */ ("+j2gort.StrValString(this_.GenerateCode(compiler, n))+")")))
	return this_
}

// Java method: Builder.removeCast(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *SuggestedFix_Builder) RemoveCast(n pkg_rhino.Node, compiler AbstractCompiler) *SuggestedFix_Builder {
	_ = n
	_ = compiler
	if !(n.IsCast()) {
		j2gort.CheckArgument(false)
	}
	var jsDoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = jsDoc
	j2gort.MultimapBuilderPutStr(this_.replacements, n.GetSourceFileName(), CodeReplacement_Create_3(jsDoc.GetOriginalCommentPosition(), (n.GetFirstChild().GetSourceOffset()-jsDoc.GetOriginalCommentPosition()), ""))
	j2gort.MultimapBuilderPutStr(this_.replacements, n.GetSourceFileName(), CodeReplacement_Create_3(((n.GetSourceOffset()+n.GetLength())-int32(1)), int32(1), ""))
	return this_
}

// Java method: Builder.addOrReplaceJsDoc(com.google.javascript.rhino.Node n, java.lang.String newJsDoc)
func (this_ *SuggestedFix_Builder) AddOrReplaceJsDoc(n pkg_rhino.Node, newJsDoc string) *SuggestedFix_Builder {
	_ = n
	_ = newJsDoc
	var offset int32 = n.GetSourceOffset()
	_ = offset
	var length int32 = int32(0)
	_ = length
	if n.IsGetProp() {
		offset = SuggestedFix_Builder_SubtreeRangeOfIdentifier(n).GetSourceOffset()
	}
	var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
	_ = jsDoc
	if jsDoc != nil {
		offset = jsDoc.GetOriginalCommentPosition()
		length = j2gort.JStrLen(jsDoc.GetOriginalCommentString())
	}
	j2gort.MultimapBuilderPutStr(this_.replacements, n.GetSourceFileName(), CodeReplacement_Create_3(offset, length, newJsDoc))
	return this_
}

// Java method: Builder.changeJsDocType(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String type)
func (this_ *SuggestedFix_Builder) ChangeJsDocType(n pkg_rhino.Node, compiler AbstractCompiler, type_ string) *SuggestedFix_Builder {
	_ = n
	_ = compiler
	_ = type_
	var typeNode pkg_rhino.Node = pkg_parsing.JsDocInfoParser_ParseTypeString(type_)
	_ = typeNode
	j2gort.CheckNotNullNode(typeNode)
	var typeExpr pkg_rhino.JSTypeExpression = pkg_rhino.NewJSTypeExpression(typeNode, "jsflume")
	_ = typeExpr
	var newJsType pkg_rhino.JSType = typeExpr.Evaluate(nil, compiler.GetTypeRegistry())
	_ = newJsType
	if newJsType == nil || newJsType.IsJavaNil_() {
		panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "JS Compiler does not recognize type: "+j2gort.StrValString(type_))))
	}
	this_.priv_Builder_replaceTypePattern_7987(n, type_, j2gort.Ext_java_util_regex_Pattern_Compile_1("@(type) *\\{?[^@\\s}]+\\}?"))
	this_.priv_Builder_replaceTypePattern_7987(n, type_, j2gort.Ext_java_util_regex_Pattern_Compile_1("@(export|package|private|protected|public|const|return) *\\{[^}]+\\}"))
	return this_
}

// Java method: Builder.replaceTypePattern(com.google.javascript.rhino.Node n, java.lang.String type, java.util.regex.Pattern pattern)
func (this_ *SuggestedFix_Builder) priv_Builder_replaceTypePattern_7987(n pkg_rhino.Node, type_ string, pattern *j2gort.Ext_java_util_regex_Pattern) {
	_ = n
	_ = type_
	_ = pattern
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
	_ = info
	j2gort.CheckNotNullPtr(info)
	var originalComment string = info.GetOriginalCommentString()
	_ = originalComment
	var originalPosition int32 = info.GetOriginalCommentPosition()
	_ = originalPosition
	if originalComment != j2gort.NullString {
		var m *j2gort.Ext_java_util_regex_Matcher = pattern.Matcher(j2gort.BoxString(originalComment))
		_ = m
		for m.Find_0() {
			j2gort.MultimapBuilderPutStr(this_.replacements, n.GetSourceFileName(), CodeReplacement_Create_3((originalPosition+m.Start_0()), (m.End_0()-m.Start_0()), ("@"+j2gort.StrValString(m.Group_1_int32(int32(1)))+" {"+j2gort.StrValString(type_)+"}")))
		}
	}
}

// Java method: Builder.insertArguments(com.google.javascript.rhino.Node n, int position, java.lang.String[] args)
func (this_ *SuggestedFix_Builder) InsertArguments(n pkg_rhino.Node, position int32, args []string) *SuggestedFix_Builder {
	_ = n
	_ = position
	_ = args
	if !(n.IsCall()) {
		j2gort.CheckArgumentStr(false, "insertArguments is only applicable to function call nodes.")
	}
	var startPosition int32 = int32(0)
	_ = startPosition
	var argument pkg_rhino.Node = n.GetSecondChild()
	_ = argument
	var i int32 = int32(0)
	_ = i
	for (argument != 0) && (i < position) {
		argument = argument.GetNext()
		i = i + 1
	}
	if argument == 0 {
		if !(position == i) {
			j2gort.CheckArgumentStr(false, "The specified position must be less than the number of arguments.")
		}
		startPosition = ((n.GetSourceOffset() + n.GetLength()) - int32(1))
	} else {
		startPosition = SuggestedFix_GetStartPositionForNodeConsideringComments(argument)
	}
	var newContent string = j2gort.Ext_com_google_common_base_Joiner_On_1_string(", ").Join_1_slice_any(func() []any {
		arr_2 := args
		if arr_2 == nil {
			return nil
		}
		out_3 := make([]any, len(arr_2))
		for i_, el_4 := range arr_2 {
			out_3[i_] = j2gort.BoxString(el_4)
		}
		return out_3
	}())
	_ = newContent
	if argument != 0 {
		newContent = (j2gort.StrValString(newContent) + ", ")
	} else if i > int32(0) {
		newContent = (", " + j2gort.StrValString(newContent))
	}
	j2gort.MultimapBuilderPutStr(this_.replacements, n.GetSourceFileName(), CodeReplacement_Create_3(startPosition, int32(0), newContent))
	return this_
}

// Java method: Builder.deleteArgument(com.google.javascript.rhino.Node n, int position)
func (this_ *SuggestedFix_Builder) DeleteArgument(n pkg_rhino.Node, position int32) *SuggestedFix_Builder {
	_ = n
	_ = position
	if !(n.IsCall() || n.IsNew()) {
		j2gort.CheckArgumentStr(false, "deleteArgument is only applicable to function call nodes.")
	}
	var numArguments int32 = (n.GetChildCount() - int32(1))
	_ = numArguments
	if !(numArguments > int32(0)) {
		j2gort.CheckStateStr(false, "deleteArgument() cannot be used on a function call with no arguments")
	}
	if !((position >= int32(0)) && (position < numArguments)) {
		j2gort.CheckArgumentStr(false, "The specified position must be less than the number of arguments.")
	}
	var argument pkg_rhino.Node = n.GetSecondChild()
	_ = argument
	var startOfArgumentToRemove int32 = int32(-1)
	_ = startOfArgumentToRemove
	var endOfArgumentToRemove int32 = int32(-1)
	_ = endOfArgumentToRemove
	var i int32 = int32(0)
	_ = i
	for argument != 0 {
		if i < position {
			startOfArgumentToRemove = (argument.GetSourceOffset() + argument.GetLength())
		} else if i == position {
			if position == int32(0) {
				startOfArgumentToRemove = SuggestedFix_GetStartPositionForNodeConsideringComments(argument)
			}
			endOfArgumentToRemove = (argument.GetSourceOffset() + argument.GetLength())
		} else if i > position {
			if position == int32(0) {
				endOfArgumentToRemove = argument.GetSourceOffset()
			}
			break
		}
		argument = argument.GetNext()
		i = i + 1
	}
	var lengthOfArgumentToRemove int32 = (endOfArgumentToRemove - startOfArgumentToRemove)
	_ = lengthOfArgumentToRemove
	j2gort.MultimapBuilderPutStr(this_.replacements, n.GetSourceFileName(), CodeReplacement_Create_3(startOfArgumentToRemove, lengthOfArgumentToRemove, ""))
	return this_
}

// Java static method: Builder.createImportNode(com.google.javascript.refactoring.SuggestedFix.ImportType importType, java.lang.@org.jspecify.annotations.Nullable String alias, java.lang.String namespace)
func SuggestedFix_Builder_CreateImportNode(importType SuggestedFix_ImportType, alias string, namespace string) pkg_rhino.Node {
	_ = importType
	_ = alias
	_ = namespace
	var requireFlavor string = j2gort.NullString
	_ = requireFlavor
	switch (importType).Ordinal() {
	case int32(0) /* ImportType.REQUIRE */ :
		requireFlavor = "require"
	case int32(1) /* ImportType.REQUIRE_TYPE */ :
		requireFlavor = "requireType"
	}
	var callNode pkg_rhino.Node = pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), requireFlavor), []pkg_rhino.Node{pkg_rhino.IR_String(namespace)})
	_ = callNode
	if alias != j2gort.NullString {
		return pkg_rhino.IR_ConstNode(pkg_rhino.IR_Name(alias), callNode)
	} else {
		return pkg_rhino.IR_ExprResult(callNode)
	}
}

// Java method: Builder.addGoogRequire(com.google.javascript.refactoring.Match m, java.lang.String namespace, com.google.javascript.refactoring.ScriptMetadata scriptMetadata)
func (this_ *SuggestedFix_Builder) AddGoogRequire(m *Match, namespace string, scriptMetadata *ScriptMetadata) *SuggestedFix_Builder {
	_ = m
	_ = namespace
	_ = scriptMetadata
	return this_.AddImport(m, namespace, StaticGet_SuggestedFix_ImportType_REQUIRE(), scriptMetadata)
}

// Java method: Builder.addGoogRequireType(com.google.javascript.refactoring.Match m, java.lang.String namespace, com.google.javascript.refactoring.ScriptMetadata scriptMetadata)
func (this_ *SuggestedFix_Builder) AddGoogRequireType(m *Match, namespace string, scriptMetadata *ScriptMetadata) *SuggestedFix_Builder {
	_ = m
	_ = namespace
	_ = scriptMetadata
	return this_.AddImport(m, namespace, StaticGet_SuggestedFix_ImportType_REQUIRE_TYPE(), scriptMetadata)
}

// Java method: Builder.addImport(com.google.javascript.refactoring.Match m, java.lang.String namespace, com.google.javascript.refactoring.SuggestedFix.ImportType importType, com.google.javascript.refactoring.ScriptMetadata scriptMetadata)
func (this_ *SuggestedFix_Builder) AddImport(m *Match, namespace string, importType SuggestedFix_ImportType, scriptMetadata *ScriptMetadata) *SuggestedFix_Builder {
	_ = m
	_ = namespace
	_ = importType
	_ = scriptMetadata
	var alias string = j2gort.NullString
	_ = alias
	if scriptMetadata.SupportsRequireAliases() {
		var existingAlias string = scriptMetadata.GetAlias(namespace)
		_ = existingAlias
		if existingAlias != j2gort.NullString {
			alias = existingAlias
		} else if j2gort.JStrIndexOfCp(namespace, int32(uint16(46))) == int32(-1) {
			alias = namespace
		} else {
			alias = j2gort.CastString(j2gort.Ext_com_google_common_collect_Streams_Stream_1_Ext_java_lang_Iterable(RequireAliasGenerator_Over(namespace)).Filter(j2gort.NewLambda_Ext_java_util_function_Predicate(func(lp_5 any) bool {
				a := j2gort.CastString(lp_5)
				_ = a
				return !(scriptMetadata.UsesName(a))
			})).FindFirst().OrElseThrow_1(lambdaSingleton_suggested_fix_1))
		}
		scriptMetadata.AddAlias(namespace, alias)
	} else {
		alias = j2gort.NullString
	}
	var metadata *NodeMetadata = m.GetMetadata()
	_ = metadata
	var existingNode pkg_rhino.Node = SuggestedFix_Builder_FindGoogRequireNode(m.GetNode(), metadata, namespace)
	_ = existingNode
	if existingNode != 0 {
		if existingNode.IsExprResult() && (alias != j2gort.NullString) {
			var newNode pkg_rhino.Node = pkg_rhino.Node(0)
			_ = newNode
			if NodeUtil_IsCallTo_2_Node_string(existingNode.GetFirstChild(), "goog.forwardDeclare") {
				newNode = SuggestedFix_Builder_CreateImportNode(importType, alias, namespace)
			} else {
				newNode = pkg_rhino.IR_ConstNode(pkg_rhino.IR_Name(alias), existingNode.GetFirstChild().CloneTree_0())
			}
			this_.Replace(existingNode, newNode, m.GetMetadata().GetCompiler())
			scriptMetadata.AddAlias(namespace, alias)
		}
		return this_
	}
	var script pkg_rhino.Node = scriptMetadata.GetScript()
	_ = script
	if script.GetFirstChild().IsModuleBody() {
		script = script.GetFirstChild()
	}
	var lastModuleOrProvideNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = lastModuleOrProvideNode
	var lastGoogRequireNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = lastGoogRequireNode
	var nodeToInsertBefore pkg_rhino.Node = pkg_rhino.Node(0)
	_ = nodeToInsertBefore
	var child pkg_rhino.Node = script.GetFirstChild()
	_ = child
	for child != 0 {
		if Matchers_GoogModule().Matches(child, metadata) {
			lastModuleOrProvideNode = child
		}
		if NodeUtil_IsExprCall(child) {
			var grandchild pkg_rhino.Node = child.GetFirstChild()
			_ = grandchild
			if Matchers_GoogModuleOrProvide().Matches(grandchild, metadata) {
				lastModuleOrProvideNode = grandchild
			} else if Matchers_GoogRequirelike_0().Matches(grandchild, metadata) {
				lastGoogRequireNode = grandchild
				if grandchild.GetLastChild().IsStringLit() && (j2gort.JStrCompareTo(namespace, grandchild.GetLastChild().GetString()) < int32(0)) {
					nodeToInsertBefore = child
					break
				}
			}
		} else if (NodeUtil_IsNameDeclaration(child) && (child.GetFirstFirstChild() != 0)) && Matchers_GoogRequirelike_0().Matches(child.GetFirstFirstChild(), metadata) {
			lastGoogRequireNode = child.GetFirstFirstChild()
			var requireName string = child.GetFirstChild().GetString()
			_ = requireName
			var originalName string = child.GetFirstChild().GetOriginalName()
			_ = originalName
			if originalName != j2gort.NullString {
				requireName = originalName
			}
			if j2gort.JStrCompareTo(alias, requireName) < int32(0) {
				nodeToInsertBefore = child
				break
			}
		}
		child = child.GetNext()
	}
	var newImportNode pkg_rhino.Node = SuggestedFix_Builder_CreateImportNode(importType, alias, namespace)
	_ = newImportNode
	if nodeToInsertBefore == 0 {
		if (lastModuleOrProvideNode != 0) || (lastGoogRequireNode != 0) {
			var nodeToInsertAfter pkg_rhino.Node
			_ = nodeToInsertAfter
			if lastGoogRequireNode != 0 {
				nodeToInsertAfter = lastGoogRequireNode
			} else {
				nodeToInsertAfter = lastModuleOrProvideNode
			}
			var startPosition int32 = ((nodeToInsertAfter.GetSourceOffset() + nodeToInsertAfter.GetLength()) + int32(2))
			_ = startPosition
			j2gort.MultimapBuilderPutStr(this_.replacements, nodeToInsertAfter.GetSourceFileName(), CodeReplacement_Create_4(startPosition, int32(0), this_.GenerateCode(m.GetMetadata().GetCompiler(), newImportNode), namespace))
			return this_
		} else {
			if script.HasChildren() {
				nodeToInsertBefore = script.GetFirstChild()
			} else {
				j2gort.MultimapBuilderPutStr(this_.replacements, script.GetSourceFileName(), CodeReplacement_Create_4(int32(0), int32(0), this_.GenerateCode(m.GetMetadata().GetCompiler(), newImportNode), namespace))
				return this_
			}
		}
	}
	return this_.InsertBefore_4(nodeToInsertBefore, newImportNode, m.GetMetadata().GetCompiler(), namespace)
}

// Java method: Builder.removeGoogRequire(com.google.javascript.refactoring.Match m, java.lang.String namespace)
func (this_ *SuggestedFix_Builder) RemoveGoogRequire(m *Match, namespace string) *SuggestedFix_Builder {
	_ = m
	_ = namespace
	var googRequireNode pkg_rhino.Node = SuggestedFix_Builder_FindGoogRequireNode(m.GetNode(), m.GetMetadata(), namespace)
	_ = googRequireNode
	if googRequireNode != 0 {
		return this_.DeleteWithoutRemovingWhitespaceBefore(googRequireNode)
	}
	return this_
}

// Java static method: Builder.findGoogRequireNode(com.google.javascript.rhino.Node n, com.google.javascript.refactoring.NodeMetadata metadata, java.lang.String namespace)
func SuggestedFix_Builder_FindGoogRequireNode(n pkg_rhino.Node, metadata *NodeMetadata, namespace string) pkg_rhino.Node {
	_ = n
	_ = metadata
	_ = namespace
	var script pkg_rhino.Node = metadata.GetCompiler().GetScriptNode(n.GetSourceFileName())
	_ = script
	if script.GetFirstChild().IsModuleBody() {
		script = script.GetFirstChild()
	}
	{
		var child pkg_rhino.Node = script.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if NodeUtil_IsExprCall(child) && Matchers_GoogRequirelike_1(namespace).Matches(child.GetFirstChild(), metadata) {
				return child
			}
		}
	}
	{
		var child_2 pkg_rhino.Node = script.GetFirstChild()
		_ = child_2
		for ; child_2 != 0; child_2 = child_2.GetNext() {
			if ((NodeUtil_IsNameDeclaration(child_2) && !(child_2.GetFirstChild().IsDestructuringLhs())) && (child_2.GetFirstChild().GetLastChild() != 0)) && Matchers_GoogRequirelike_1(namespace).Matches(child_2.GetFirstChild().GetLastChild(), metadata) {
				return child_2
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java static method: Builder.subtreeRangeOfIdentifier(com.google.javascript.rhino.Node n)
func SuggestedFix_Builder_SubtreeRangeOfIdentifier(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsGetProp()) {
		j2gort.CheckStateStr(false, "Support other identifier nodes")
	}
	var leftmost pkg_rhino.Node = n
	_ = leftmost
	for leftmost.HasChildren() {
		leftmost = leftmost.GetFirstChild()
	}
	var result pkg_rhino.Node = pkg_rhino.IR_Empty()
	_ = result
	result.SetStaticSourceFile(n.GetStaticSourceFile())
	result.SetLinenoCharno(leftmost.GetLineno(), leftmost.GetCharno())
	result.SetLength((n.GetLength() + (n.GetSourceOffset() - leftmost.GetSourceOffset())))
	return result
}

// Java method: Builder.generateCode(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node node)
func (this_ *SuggestedFix_Builder) GenerateCode(compiler AbstractCompiler, node pkg_rhino.Node) string {
	_ = compiler
	_ = node
	if node.IsBlock() {
		node.SetToken(pkg_rhino.StaticGet_Token_SCRIPT())
	}
	clinit_CompilerOptions()
	var val_6 CompilerOptions
	InitCompilerOptions(&val_6)
	var compilerOptions *CompilerOptions = &val_6
	_ = compilerOptions
	compilerOptions.SetPreferSingleQuotes(true)
	compilerOptions.SetUseOriginalNamesInOutput(true)
	compilerOptions.SetPreserveNonJSDocComments(true)
	compilerOptions.SetTrustedStrings(true)
	return NewCodePrinter_Builder(node).SetCompilerOptions(compilerOptions).SetTypeRegistry(compiler.GetTypeRegistry()).SetPrettyPrint(true).SetLineBreak(true).SetOutputTypes(true).Build()
}

// Java method: Builder.setDescription(java.lang.String description)
func (this_ *SuggestedFix_Builder) SetDescription(description string) *SuggestedFix_Builder {
	_ = description
	this_.description = description
	return this_
}

// Java method: Builder.build()
func (this_ *SuggestedFix_Builder) Build() *SuggestedFix {
	return NewSuggestedFix(this_.matchedNodeInfo, j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableSetMultimap](this_.replacements.Build())), this_.description, this_.alternatives.BuildImmList_())
}

// Java record: com.google.javascript.refactoring.SuggestedFix.MatchedNodeInfo
type SuggestedFix_MatchedNodeInfo struct {
	sourceFilename   string `gson:"sourceFilename"`
	lineno           int32  `gson:"lineno"`
	charno           int32  `gson:"charno"`
	inClosurizedFile bool   `gson:"inClosurizedFile"`
}

func (this_ *SuggestedFix_MatchedNodeInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *SuggestedFix_MatchedNodeInfo) IsJavaNil_() bool { return this_ == nil }
func (*SuggestedFix_MatchedNodeInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.SuggestedFix$MatchedNodeInfo", "MatchedNodeInfo"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.SuggestedFix.MatchedNodeInfo", "com.google.javascript.refactoring.SuggestedFix$MatchedNodeInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*SuggestedFix_MatchedNodeInfo](v_)
	})
}
func (*SuggestedFix_MatchedNodeInfo) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *SuggestedFix_MatchedNodeInfo) initFields_SuggestedFix_MatchedNodeInfo() {
	_ = this_
}

// Java constructor: MatchedNodeInfo(java.lang.String sourceFilename, int lineno, int charno, boolean inClosurizedFile)
func NewSuggestedFix_MatchedNodeInfo(sourceFilename string, lineno int32, charno int32, inClosurizedFile bool) *SuggestedFix_MatchedNodeInfo {
	this_ := &SuggestedFix_MatchedNodeInfo{}
	InitSuggestedFix_MatchedNodeInfo(this_, sourceFilename, lineno, charno, inClosurizedFile)
	return this_
}

func InitSuggestedFix_MatchedNodeInfo(this_ *SuggestedFix_MatchedNodeInfo, sourceFilename string, lineno int32, charno int32, inClosurizedFile bool) {
	_ = sourceFilename
	_ = lineno
	_ = charno
	_ = inClosurizedFile
	this_.initFields_SuggestedFix_MatchedNodeInfo()
	j2gort.CheckNotNullString(sourceFilename)
	this_.sourceFilename = sourceFilename
	this_.lineno = lineno
	this_.charno = charno
	this_.inClosurizedFile = inClosurizedFile
}

// Java method: MatchedNodeInfo.getSourceFilename()
func (this_ *SuggestedFix_MatchedNodeInfo) GetSourceFilename() string {
	return this_.SourceFilename()
}

// Java method: MatchedNodeInfo.getLineno()
func (this_ *SuggestedFix_MatchedNodeInfo) GetLineno() int32 {
	return this_.Lineno()
}

// Java method: MatchedNodeInfo.getCharno()
func (this_ *SuggestedFix_MatchedNodeInfo) GetCharno() int32 {
	return this_.Charno()
}

// Java method: MatchedNodeInfo.isInClosurizedFile()
func (this_ *SuggestedFix_MatchedNodeInfo) IsInClosurizedFile() bool {
	return this_.InClosurizedFile()
}

// Java static method: MatchedNodeInfo.create(com.google.javascript.rhino.Node node, boolean closurized)
func SuggestedFix_MatchedNodeInfo_Create(node pkg_rhino.Node, closurized bool) *SuggestedFix_MatchedNodeInfo {
	_ = node
	_ = closurized
	return NewSuggestedFix_MatchedNodeInfo(NodeUtil_GetSourceName(node), node.GetLineno(), node.GetCharno(), closurized)
}

func (this_ *SuggestedFix_MatchedNodeInfo) SourceFilename() string { return this_.sourceFilename }
func (this_ *SuggestedFix_MatchedNodeInfo) Lineno() int32          { return this_.lineno }
func (this_ *SuggestedFix_MatchedNodeInfo) Charno() int32          { return this_.charno }
func (this_ *SuggestedFix_MatchedNodeInfo) InClosurizedFile() bool { return this_.inClosurizedFile }
func (this_ *SuggestedFix_MatchedNodeInfo) Equals(other_ any) bool {
	o_, ok_ := other_.(*SuggestedFix_MatchedNodeInfo)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.sourceFilename != o_.sourceFilename {
		return false
	}
	if this_.lineno != o_.lineno {
		return false
	}
	if this_.charno != o_.charno {
		return false
	}
	if this_.inClosurizedFile != o_.inClosurizedFile {
		return false
	}
	return true
}
func (this_ *SuggestedFix_MatchedNodeInfo) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_7 := this_.sourceFilename
		h_ = 31 * h_
		if hs_7 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_7)
		}
	}
	h_ = 31*h_ + (this_.lineno)
	h_ = 31*h_ + (this_.charno)
	if this_.inClosurizedFile {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	return h_
}
func (this_ *SuggestedFix_MatchedNodeInfo) String() string {
	_ = fmt.Sprint
	return "MatchedNodeInfo[" + "sourceFilename=" + j2gort.StrVal(j2gort.BoxString(this_.sourceFilename)) + ", " + "lineno=" + j2gort.StrVal(j2gort.ValueOfInt(this_.lineno)) + ", " + "charno=" + j2gort.StrVal(j2gort.ValueOfInt(this_.charno)) + ", " + "inClosurizedFile=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.inClosurizedFile)) + "]"
}
func (this_ *SuggestedFix) Slice_GetAlternatives() []*SuggestedFix {
	return this_.alternatives
}
