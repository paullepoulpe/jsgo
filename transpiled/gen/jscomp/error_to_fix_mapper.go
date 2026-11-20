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
// Port of com/google/javascript/refactoring/ErrorToFixMapper.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.refactoring.ErrorToFixMapper
type ErrorToFixMapper struct {
	compiler              AbstractCompiler                                   `gson:"compiler"`
	metadataForEachScript j2gort.OrderedMap[pkg_rhino.Node, *ScriptMetadata] `gson:"metadataForEachScript"`
	idHash_               atomic.Int32
}

func (this_ *ErrorToFixMapper) BindSelf_(s any) {
	_ = s
}
func (this_ *ErrorToFixMapper) IsJavaNil_() bool { return this_ == nil }
func (*ErrorToFixMapper) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.ErrorToFixMapper", "ErrorToFixMapper"
}
func (this_ *ErrorToFixMapper) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ErrorToFixMapper) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.ErrorToFixMapper", "com.google.javascript.refactoring.ErrorToFixMapper", func(v_ any) bool {
		return j2gort.InstanceOf[*ErrorToFixMapper](v_)
	})
}

var static_ErrorToFixMapper_DID_YOU_MEAN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_ErrorToFixMapper_DID_YOU_MEAN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_ErrorToFixMapper()
	return static_ErrorToFixMapper_DID_YOU_MEAN
}
func StaticSet_ErrorToFixMapper_DID_YOU_MEAN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_ErrorToFixMapper()
	static_ErrorToFixMapper_DID_YOU_MEAN = v
}

var static_ErrorToFixMapper_EARLY_REF *j2gort.Ext_java_util_regex_Pattern

func StaticGet_ErrorToFixMapper_EARLY_REF() *j2gort.Ext_java_util_regex_Pattern {
	clinit_ErrorToFixMapper()
	return static_ErrorToFixMapper_EARLY_REF
}
func StaticSet_ErrorToFixMapper_EARLY_REF(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_ErrorToFixMapper()
	static_ErrorToFixMapper_EARLY_REF = v
}

var static_ErrorToFixMapper_MISSING_REQUIRE *j2gort.Ext_java_util_regex_Pattern

func StaticGet_ErrorToFixMapper_MISSING_REQUIRE() *j2gort.Ext_java_util_regex_Pattern {
	clinit_ErrorToFixMapper()
	return static_ErrorToFixMapper_MISSING_REQUIRE
}
func StaticSet_ErrorToFixMapper_MISSING_REQUIRE(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_ErrorToFixMapper()
	static_ErrorToFixMapper_MISSING_REQUIRE = v
}

var static_ErrorToFixMapper_FULLY_QUALIFIED_NAME *j2gort.Ext_java_util_regex_Pattern

func StaticGet_ErrorToFixMapper_FULLY_QUALIFIED_NAME() *j2gort.Ext_java_util_regex_Pattern {
	clinit_ErrorToFixMapper()
	return static_ErrorToFixMapper_FULLY_QUALIFIED_NAME
}
func StaticSet_ErrorToFixMapper_FULLY_QUALIFIED_NAME(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_ErrorToFixMapper()
	static_ErrorToFixMapper_FULLY_QUALIFIED_NAME = v
}

var static_ErrorToFixMapper_USE_SHORT_NAME *j2gort.Ext_java_util_regex_Pattern

func StaticGet_ErrorToFixMapper_USE_SHORT_NAME() *j2gort.Ext_java_util_regex_Pattern {
	clinit_ErrorToFixMapper()
	return static_ErrorToFixMapper_USE_SHORT_NAME
}
func StaticSet_ErrorToFixMapper_USE_SHORT_NAME(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_ErrorToFixMapper()
	static_ErrorToFixMapper_USE_SHORT_NAME = v
}

var clinitDone_ErrorToFixMapper uint32
var clinitState_ErrorToFixMapper j2gort.ClinitState

func clinit_ErrorToFixMapper() {
	if atomic.LoadUint32(&clinitDone_ErrorToFixMapper) != 0 {
		return
	}
	clinitSlow_ErrorToFixMapper()
}

func clinitSlow_ErrorToFixMapper() {
	if !j2gort.EnterClinit(&clinitState_ErrorToFixMapper, &clinitDone_ErrorToFixMapper) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ErrorToFixMapper, &clinitDone_ErrorToFixMapper)
	static_ErrorToFixMapper_DID_YOU_MEAN = j2gort.Ext_java_util_regex_Pattern_Compile_1(".*Did you mean (.*)\\?")
	static_ErrorToFixMapper_EARLY_REF = j2gort.Ext_java_util_regex_Pattern_Compile_1("Variable referenced before declaration: (.*)")
	static_ErrorToFixMapper_MISSING_REQUIRE = j2gort.Ext_java_util_regex_Pattern_Compile_2("'([^']+)' references a .*", int32(32))
	static_ErrorToFixMapper_FULLY_QUALIFIED_NAME = j2gort.Ext_java_util_regex_Pattern_Compile_1("Reference to fully qualified import name '([^']+)'.*")
	static_ErrorToFixMapper_USE_SHORT_NAME = j2gort.Ext_java_util_regex_Pattern_Compile_1(".*Please use the short name '(.*)' instead.")
}

func (this_ *ErrorToFixMapper) initFields_ErrorToFixMapper() {
	_ = this_
	this_.metadataForEachScript = j2gort.NewOrderedMap[pkg_rhino.Node, *ScriptMetadata](

	// Java constructor: ErrorToFixMapper(com.google.javascript.jscomp.AbstractCompiler compiler)
	)
}

func NewErrorToFixMapper(compiler AbstractCompiler) *ErrorToFixMapper {
	clinit_ErrorToFixMapper()
	this_ := &ErrorToFixMapper{}
	InitErrorToFixMapper(this_, compiler)
	return this_
}

func InitErrorToFixMapper(this_ *ErrorToFixMapper, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ErrorToFixMapper()
	this_.compiler = compiler
}

// Java method: ErrorToFixMapper.getFixesForJsError(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) GetFixesForJsError(error_ *JSError) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = error_
	var fix *SuggestedFix = this_.GetFixForJsError(error_)
	_ = fix
	if fix != nil {
		return j2gort.ImmListOfPtr(fix)
	}
	switch j2gort.CheckNotNullString(error_.Type().key) {
	case "JSC_IMPLICITLY_NONNULL_JSDOC", "JSC_IMPLICITLY_NULLABLE_JSDOC", "JSC_MISSING_NULLABILITY_MODIFIER_JSDOC", "JSC_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC":
		return this_.priv_ErrorToFixMapper_getFixesForImplicitNullabilityErrors_7944(error_)
	default:
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	panic("unreachable switch expression")
}

// Java method: ErrorToFixMapper.getFixForJsError(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) GetFixForJsError(error_ *JSError) *SuggestedFix {
	_ = error_
	switch j2gort.CheckNotNullString(error_.Type().key) {
	case "JSC_REDECLARED_VARIABLE":
		return this_.priv_ErrorToFixMapper_getFixForRedeclaration_7941(error_)
	case "JSC_REFERENCE_BEFORE_DECLARE":
		return this_.priv_ErrorToFixMapper_getFixForEarlyReference_7942(error_)
	case "JSC_MISSING_SEMICOLON":
		return this_.priv_ErrorToFixMapper_getFixForMissingSemicolon_7946(error_)
	case "JSC_REQUIRES_NOT_SORTED":
		return this_.priv_ErrorToFixMapper_getFixForUnsortedRequires_7952(error_)
	case "JSC_PROVIDES_NOT_SORTED":
		return this_.priv_ErrorToFixMapper_getFixForUnsortedProvides_7953(error_)
	case "JSC_DEBUGGER_STATEMENT_PRESENT":
		return NewSuggestedFix_Builder().AttachMatchedNodeInfo(error_.Node(), this_.compiler).SetDescription("Remove debugger statement").Delete_1(error_.Node()).Build()
	case "JSC_USELESS_EMPTY_STATEMENT":
		return this_.priv_ErrorToFixMapper_removeEmptyStatement_7945(error_)
	case "JSC_INEXISTENT_PROPERTY_WITH_SUGGESTION", "JSC_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION":
		return this_.priv_ErrorToFixMapper_getFixForInexistentProperty_7949(error_)
	case "JSC_MISSING_CALL_TO_SUPER":
		return this_.priv_ErrorToFixMapper_getFixForMissingSuper_7947(error_)
	case "JSC_INVALID_SUPER_CALL_WITH_SUGGESTION":
		return this_.priv_ErrorToFixMapper_getFixForInvalidSuper_7948(error_)
	case "JSC_MISSING_REQUIRE", "JSC_MISSING_REQUIRE_IN_PROVIDES_FILE":
		return this_.priv_ErrorToFixMapper_getFixForMissingRequire_7950(error_, StaticGet_SuggestedFix_ImportType_REQUIRE())
	case "JSC_MISSING_REQUIRE_TYPE", "JSC_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE":
		return this_.priv_ErrorToFixMapper_getFixForMissingRequire_7950(error_, StaticGet_SuggestedFix_ImportType_REQUIRE_TYPE())
	case "JSC_EXTRA_REQUIRE_WARNING":
		return this_.priv_ErrorToFixMapper_getFixForExtraRequire_7951(error_)
	case "JSC_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME", "JSC_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME":
		return this_.priv_ErrorToFixMapper_getFixForReferenceToShortImportByLongName_7943(error_)
	case "JSC_REDUNDANT_NULLABILITY_MODIFIER_JSDOC":
		return this_.priv_ErrorToFixMapper_getFixForRedundantNullabilityModifierJsDoc_7954(error_)
	case "JSC_MISSING_CONST_ON_CONSTANT_CASE":
		return this_.priv_ErrorToFixMapper_getFixForConstantCaseErrors_7955(error_)
	default:
		return nil
	}
	panic("unreachable switch expression")
}

// Java method: ErrorToFixMapper.getFixForRedeclaration(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForRedeclaration_7941(error_ *JSError) *SuggestedFix {
	_ = error_
	var name pkg_rhino.Node = error_.Node()
	_ = name
	if !(name.IsName()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(name))
	}
	var parent pkg_rhino.Node = name.GetParent()
	_ = parent
	if !(NodeUtil_IsNameDeclaration(parent)) {
		return nil
	}
	var fix *SuggestedFix_Builder = NewSuggestedFix_Builder().AttachMatchedNodeInfo(name, this_.compiler)
	_ = fix
	if !(name.HasChildren()) {
		var nodeToDelete pkg_rhino.Node
		_ = nodeToDelete
		if parent.HasOneChild() {
			nodeToDelete = parent
		} else {
			nodeToDelete = error_.Node()
		}
		return fix.SetDescription("Remove redundant declaration").Delete_1(nodeToDelete).Build()
	}
	fix.SetDescription("Convert redundant declaration to assignment")
	var assign pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(name.CloneNode_0(), name.GetFirstChild().CloneTree_0()))
	_ = assign
	if parent.HasOneChild() {
		return fix.Replace(parent, assign, this_.compiler).Build()
	}
	var childrenOfAddedVarStatement []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = childrenOfAddedVarStatement
	{
		var n pkg_rhino.Node = parent.GetFirstChild()
		_ = n
		for ; n != 0; n = n.GetNext() {
			if n == name {
				break
			}
			childrenOfAddedVarStatement = append(childrenOfAddedVarStatement, n)
		}
	}
	if len(childrenOfAddedVarStatement) != 0 {
		var var_ pkg_rhino.Node = pkg_rhino.NewNode_1(parent.GetToken())
		_ = var_
		for _, n_2 := range childrenOfAddedVarStatement {
			_ = n_2
			var_.AddChildToBack(n_2.CloneTree_0())
		}

		fix.InsertBefore_4(parent, var_, this_.compiler, "1")
	}
	if name.GetNext() != 0 {
		for _, n_3 := range childrenOfAddedVarStatement {
			_ = n_3
			fix.Delete_1(n_3)
		}

		fix.Delete_1(name)
		fix.InsertBefore_4(parent, assign, this_.compiler, "2")
	} else {
		fix.Replace(parent, assign, this_.compiler)
	}
	return fix.Build()
}

// Java method: ErrorToFixMapper.getFixForEarlyReference(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForEarlyReference_7942(error_ *JSError) *SuggestedFix {
	_ = error_
	var m *j2gort.Ext_java_util_regex_Matcher = StaticGet_ErrorToFixMapper_EARLY_REF().Matcher(j2gort.BoxString(error_.Description()))
	_ = m
	if m.Matches() {
		var name string = m.Group_1_int32(int32(1))
		_ = name
		var stmt pkg_rhino.Node = NodeUtil_GetEnclosingStatement(error_.Node())
		_ = stmt
		return NewSuggestedFix_Builder().AttachMatchedNodeInfo(error_.Node(), this_.compiler).SetDescription("Insert var declaration statement").InsertBefore_2(stmt, ("var " + j2gort.StrValString(name) + ";\n")).Build()
	}
	return nil
}

// Java method: ErrorToFixMapper.getFixForReferenceToShortImportByLongName(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForReferenceToShortImportByLongName_7943(error_ *JSError) *SuggestedFix {
	_ = error_
	var fix *SuggestedFix_Builder = NewSuggestedFix_Builder().AttachMatchedNodeInfo(error_.Node(), this_.compiler)
	_ = fix
	var val_3 NodeMetadata
	InitNodeMetadata_1(&val_3, this_.compiler)
	var metadata *NodeMetadata = &val_3
	_ = metadata
	var val_4 Match
	InitMatch(&val_4, error_.Node(), metadata)
	var match *Match = &val_4
	_ = match
	var script pkg_rhino.Node = this_.compiler.GetScriptNode(error_.SourceName())
	_ = script
	var scriptMetadata *ScriptMetadata = this_.priv_ErrorToFixMapper_getMetadataForScript_7956(script)
	_ = scriptMetadata
	var fullNameMatcher *j2gort.Ext_java_util_regex_Matcher = StaticGet_ErrorToFixMapper_FULLY_QUALIFIED_NAME().Matcher(j2gort.BoxString(error_.Description()))
	_ = fullNameMatcher
	j2gort.CheckStateMsg(fullNameMatcher.Matches(), j2gort.BoxString(error_.Description()))
	var fullName string = fullNameMatcher.Group_1_int32(int32(1))
	_ = fullName
	var shortNameMatcher *j2gort.Ext_java_util_regex_Matcher = StaticGet_ErrorToFixMapper_USE_SHORT_NAME().Matcher(j2gort.BoxString(error_.Description()))
	_ = shortNameMatcher
	var shortName string = j2gort.NullString
	_ = shortName
	if shortNameMatcher.Matches() {
		shortName = shortNameMatcher.Group_1_int32(int32(1))
	} else {
		fix.AddGoogRequire(match, fullName, scriptMetadata)
		shortName = scriptMetadata.GetAlias(fullName)
	}
	var oldName string
	_ = oldName
	if error_.Node().IsQualifiedName() {
		oldName = error_.Node().GetQualifiedName()
	} else {
		oldName = error_.Node().GetString()
	}
	return fix.Replace(error_.Node(), NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, javalang.StringReplace(oldName, j2gort.CheckNotNullString(fullName), j2gort.CheckNotNullString(shortName))), this_.compiler).Build()
}

// Java method: ErrorToFixMapper.getFixesForImplicitNullabilityErrors(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixesForImplicitNullabilityErrors_7944(error_ *JSError) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = error_
	if error_.Node().GetSourceFileName() == j2gort.NullString {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var qmark *SuggestedFix = NewSuggestedFix_Builder().AttachMatchedNodeInfo(error_.Node(), this_.compiler).InsertBefore_2(error_.Node(), "?").SetDescription("Make nullability explicit").Build()
	_ = qmark
	var bang *SuggestedFix = NewSuggestedFix_Builder().AttachMatchedNodeInfo(error_.Node(), this_.compiler).InsertBefore_2(error_.Node(), "!").SetDescription("Make type non-nullable").Build()
	_ = bang
	switch j2gort.CheckNotNullString(error_.Type().key) {
	case "JSC_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC":
		return j2gort.ImmListOfPtr(qmark)
	case "JSC_MISSING_NULLABILITY_MODIFIER_JSDOC":
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.PtrToAny(bang), j2gort.PtrToAny(qmark))
	case "JSC_IMPLICITLY_NULLABLE_JSDOC":
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.PtrToAny(qmark), j2gort.PtrToAny(bang))
	case "JSC_IMPLICITLY_NONNULL_JSDOC":
		return j2gort.ImmListOfPtr(bang)
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unexpected JSError Type: "+j2gort.StrValString(error_.Type().key))))
	}
	panic("unreachable switch expression")
}

// Java method: ErrorToFixMapper.removeEmptyStatement(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_removeEmptyStatement_7945(error_ *JSError) *SuggestedFix {
	_ = error_
	return NewSuggestedFix_Builder().AttachMatchedNodeInfo(error_.Node(), this_.compiler).SetDescription("Remove empty statement").DeleteWithoutRemovingWhitespace(error_.Node()).Build()
}

// Java method: ErrorToFixMapper.getFixForMissingSemicolon(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForMissingSemicolon_7946(error_ *JSError) *SuggestedFix {
	_ = error_
	return NewSuggestedFix_Builder().AttachMatchedNodeInfo(error_.Node(), this_.compiler).InsertAfter(error_.Node(), ";").Build()
}

// Java method: ErrorToFixMapper.getFixForMissingSuper(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForMissingSuper_7947(error_ *JSError) *SuggestedFix {
	_ = error_
	var constructorFunction pkg_rhino.Node = error_.Node().GetFirstChild()
	_ = constructorFunction
	var body pkg_rhino.Node = NodeUtil_GetFunctionBody(constructorFunction)
	_ = body
	return NewSuggestedFix_Builder().AttachMatchedNodeInfo(error_.Node(), this_.compiler).AddChildToFront(body, "super();").Build()
}

// Java method: ErrorToFixMapper.getFixForInvalidSuper(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForInvalidSuper_7948(error_ *JSError) *SuggestedFix {
	_ = error_
	var m *j2gort.Ext_java_util_regex_Matcher = StaticGet_ErrorToFixMapper_DID_YOU_MEAN().Matcher(j2gort.BoxString(error_.Description()))
	_ = m
	if m.Matches() {
		var superDotSuggestion string = j2gort.CheckNotNullString(m.Group_1_int32(int32(1)))
		_ = superDotSuggestion
		return NewSuggestedFix_Builder().AttachMatchedNodeInfo(error_.Node(), this_.compiler).SetDescription(("Call '" + j2gort.StrValString(superDotSuggestion) + "' instead")).Replace(error_.Node(), NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, superDotSuggestion), this_.compiler).Build()
	}
	return nil
}

// Java method: ErrorToFixMapper.getFixForInexistentProperty(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForInexistentProperty_7949(error_ *JSError) *SuggestedFix {
	_ = error_
	var m *j2gort.Ext_java_util_regex_Matcher = StaticGet_ErrorToFixMapper_DID_YOU_MEAN().Matcher(j2gort.BoxString(error_.Description()))
	_ = m
	if m.Matches() {
		var suggestedPropName string = m.Group_1_int32(int32(1))
		_ = suggestedPropName
		return NewSuggestedFix_Builder().AttachMatchedNodeInfo(error_.Node(), this_.compiler).SetDescription(("Change property name to '" + j2gort.StrValString(suggestedPropName) + "'")).Rename_2(error_.Node(), suggestedPropName).Build()
	}
	return nil
}

// Java method: ErrorToFixMapper.getFixForMissingRequire(com.google.javascript.jscomp.JSError error, com.google.javascript.refactoring.SuggestedFix.ImportType importType)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForMissingRequire_7950(error_ *JSError, importType SuggestedFix_ImportType) *SuggestedFix {
	_ = error_
	_ = importType
	var regexMatcher *j2gort.Ext_java_util_regex_Matcher = StaticGet_ErrorToFixMapper_MISSING_REQUIRE().Matcher(j2gort.BoxString(error_.Description()))
	_ = regexMatcher
	j2gort.CheckStateFmt(regexMatcher.Matches(), "Unexpected error description: %s", j2gort.BoxString(error_.Description()))
	var namespaceToRequire string = regexMatcher.Group_1_int32(int32(1))
	_ = namespaceToRequire
	var qName string
	_ = qName
	if error_.Node().IsQualifiedName() {
		qName = error_.Node().GetQualifiedName()
	} else {
		qName = error_.Node().GetString()
	}
	j2gort.CheckStateFmt(javalang.StringStartsWith(j2gort.CheckNotNullString(qName), j2gort.CheckNotNullString(namespaceToRequire)), "Expected error location %s to start with namespace <%s>", j2gort.BoxNode(error_.Node()), j2gort.BoxString(namespaceToRequire))
	var script pkg_rhino.Node = this_.compiler.GetScriptNode(error_.SourceName())
	_ = script
	var scriptMetadata *ScriptMetadata = this_.priv_ErrorToFixMapper_getMetadataForScript_7956(script)
	_ = scriptMetadata
	var val_5 NodeMetadata
	InitNodeMetadata_1(&val_5, this_.compiler)
	var metadata *NodeMetadata = &val_5
	_ = metadata
	var val_6 Match
	InitMatch(&val_6, error_.Node(), metadata)
	var match *Match = &val_6
	_ = match
	var fix *SuggestedFix_Builder = NewSuggestedFix_Builder().AttachMatchedNodeInfo(error_.Node(), this_.compiler).AddImport(match, namespaceToRequire, importType, scriptMetadata)
	_ = fix
	var alias string = scriptMetadata.GetAlias(namespaceToRequire)
	_ = alias
	if alias != j2gort.NullString {
		fix.Replace(error_.Node(), NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, javalang.StringReplace(qName, j2gort.CheckNotNullString(namespaceToRequire), j2gort.CheckNotNullString(alias))), this_.compiler)
	}
	return fix.Build()
}

// Java method: ErrorToFixMapper.getFixForExtraRequire(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForExtraRequire_7951(error_ *JSError) *SuggestedFix {
	_ = error_
	var node pkg_rhino.Node = error_.Node()
	_ = node
	var fix *SuggestedFix_Builder = NewSuggestedFix_Builder().AttachMatchedNodeInfo(node, this_.compiler)
	_ = fix
	var destructuring bool = (NodeUtil_GetEnclosingType(node, pkg_rhino.StaticGet_Token_OBJECT_PATTERN()) != 0)
	_ = destructuring
	if destructuring {
		fix.SetDescription("Delete unused symbol")
		if node.IsStringKey() {
			fix.Delete_1(node)
		} else {
			if !(node.GetParent().IsStringKey()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(node.GetParent()))
			}
			fix.Delete_1(node.GetParent())
		}
	} else {
		fix.SetDescription("Delete extra require")
		fix.DeleteWithoutRemovingWhitespaceBefore(NodeUtil_GetEnclosingStatement(node))
	}
	return fix.Build()
}

// Java method: ErrorToFixMapper.getFixForUnsortedRequires(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForUnsortedRequires_7952(error_ *JSError) *SuggestedFix {
	_ = error_
	var script pkg_rhino.Node = NodeUtil_GetEnclosingScript(error_.Node())
	_ = script
	clinit_CheckRequiresSorted()
	var val_7 CheckRequiresSorted
	InitCheckRequiresSorted(&val_7, StaticGet_CheckRequiresSorted_Mode_COLLECT_ONLY())
	var callback *CheckRequiresSorted = &val_7
	_ = callback
	NodeTraversal_Traverse(this_.compiler, script, callback)
	if !(callback.NeedsFix()) {
		return nil
	}
	return NewSuggestedFix_Builder().AttachMatchedNodeInfo(callback.GetFirstNode(), this_.compiler).ReplaceRange(callback.GetFirstNode(), callback.GetLastNode(), callback.GetReplacement()).Build()
}

// Java method: ErrorToFixMapper.getFixForUnsortedProvides(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForUnsortedProvides_7953(error_ *JSError) *SuggestedFix {
	_ = error_
	var script pkg_rhino.Node = NodeUtil_GetEnclosingScript(error_.Node())
	_ = script
	clinit_CheckProvidesSorted()
	var val_8 CheckProvidesSorted
	InitCheckProvidesSorted(&val_8, StaticGet_CheckProvidesSorted_Mode_COLLECT_ONLY())
	var callback *CheckProvidesSorted = &val_8
	_ = callback
	NodeTraversal_Traverse(this_.compiler, script, callback)
	if !(callback.NeedsFix()) {
		return nil
	}
	return NewSuggestedFix_Builder().AttachMatchedNodeInfo(callback.GetFirstNode(), this_.compiler).ReplaceRange(callback.GetFirstNode(), callback.GetLastNode(), callback.GetReplacement()).Build()
}

// Java method: ErrorToFixMapper.getFixForRedundantNullabilityModifierJsDoc(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForRedundantNullabilityModifierJsDoc_7954(error_ *JSError) *SuggestedFix {
	_ = error_
	return NewSuggestedFix_Builder().AttachMatchedNodeInfo(error_.Node(), this_.compiler).ReplaceText(error_.Node(), int32(1), "").Build()
}

// Java method: ErrorToFixMapper.getFixForConstantCaseErrors(com.google.javascript.jscomp.JSError error)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getFixForConstantCaseErrors_7955(error_ *JSError) *SuggestedFix {
	_ = error_
	var n pkg_rhino.Node = error_.Node()
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if !(n.IsName()) || !(NodeUtil_IsNameDeclaration(parent)) {
		return nil
	}
	if parent.IsLet() {
		return NewSuggestedFix_Builder().SetDescription("Make explicitly constant").AttachMatchedNodeInfo(parent, this_.compiler).ReplaceText(parent, int32(3), "const").Build()
	}
	if !(parent.IsVar()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
	}
	var info *pkg_rhino.JSDocInfo = parent.GetJSDocInfo()
	_ = info
	if info == nil {
		return NewSuggestedFix_Builder().SetDescription("Make explicitly constant").AttachMatchedNodeInfo(parent, this_.compiler).AddOrReplaceJsDoc(parent, "/** @const */\n").Build()
	}
	return nil
}

// Java method: ErrorToFixMapper.getMetadataForScript(com.google.javascript.rhino.Node script)
func (this_ *ErrorToFixMapper) priv_ErrorToFixMapper_getMetadataForScript_7956(script pkg_rhino.Node) *ScriptMetadata {
	_ = script
	if !(script.IsScript()) {
		j2gort.CheckArgument(false)
	}
	return j2gort.OrderedMapComputeIfAbsent(&this_.metadataForEachScript, script, nil, func(lp_13 pkg_rhino.Node) *ScriptMetadata {
		s := lp_13
		_ = s
		return ScriptMetadata_Create_2_Node_AbstractCompiler(s, this_.compiler)
	})
}
