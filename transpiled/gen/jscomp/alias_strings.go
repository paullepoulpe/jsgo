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
// Port of com/google/javascript/jscomp/AliasStrings.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.AliasStrings
type AliasStrings struct {
	compiler                  AbstractCompiler                            `gson:"compiler"`
	chunkGraph                *JSChunkGraph                               `gson:"chunkGraph"`
	stringInfoMap             j2gort.Ext_java_util_SortedMap              `gson:"stringInfoMap"`
	usedHashedAliases         j2gort.OrderedSet[string]                   `gson:"usedHashedAliases"`
	chunkVarParentMap         j2gort.OrderedMap[*JSChunk, pkg_rhino.Node] `gson:"chunkVarParentMap"`
	unitTestHashReductionMask int32                                       `gson:"unitTestHashReductionMask"`
	idHash_                   atomic.Int32
	outputStringUsage         bool                             `gson:"outputStringUsage"`
	aliasStringsMode          CompilerOptions_AliasStringsMode `gson:"aliasStringsMode"`
}

func (this_ *AliasStrings) BindSelf_(s any) {
	_ = s
}
func (this_ *AliasStrings) IsJavaNil_() bool { return this_ == nil }
func (*AliasStrings) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AliasStrings", "AliasStrings"
}
func (this_ *AliasStrings) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AliasStrings) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AliasStrings", "com.google.javascript.jscomp.AliasStrings", func(v_ any) bool {
		return j2gort.InstanceOf[*AliasStrings](v_)
	})
}
func (*AliasStrings) Implements_jscomp_CompilerPass()           {}
func (*AliasStrings) Implements_jscomp_NodeTraversal_Callback() {}

var static_AliasStrings_logger *j2gort.Ext_java_util_logging_Logger

func StaticGet_AliasStrings_Logger() *j2gort.Ext_java_util_logging_Logger {
	clinit_AliasStrings()
	return static_AliasStrings_logger
}
func StaticSet_AliasStrings_Logger(v *j2gort.Ext_java_util_logging_Logger) {
	clinit_AliasStrings()
	static_AliasStrings_logger = v
}
func StaticGet_AliasStrings_STRING_ALIAS_PREFIX() string       { return "$$S_" }
func StaticGet_AliasStrings_ALIAS_LARGE_STRINGS_LENGTH() int32 { return int32(100) }

var clinitDone_AliasStrings uint32
var clinitState_AliasStrings j2gort.ClinitState

func clinit_AliasStrings() {
	if atomic.LoadUint32(&clinitDone_AliasStrings) != 0 {
		return
	}
	clinitSlow_AliasStrings()
}

func clinitSlow_AliasStrings() {
	if !j2gort.EnterClinit(&clinitState_AliasStrings, &clinitDone_AliasStrings) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AliasStrings, &clinitDone_AliasStrings)
	static_AliasStrings_logger = j2gort.Ext_java_util_logging_Logger_GetLogger_1(j2gort.CheckNotNullPtr(j2gort.ClassLiteral("com.google.javascript.jscomp.AliasStrings", "AliasStrings")).GetName())
}

func (this_ *AliasStrings) initFields_AliasStrings() {
	_ = this_
	this_.stringInfoMap = j2gort.NewExt_java_util_TreeMap()
	this_.usedHashedAliases = j2gort.NewOrderedSet[string]()
	this_.chunkVarParentMap = j2gort.NewOrderedMap[*JSChunk, pkg_rhino.Node]()
	this_.unitTestHashReductionMask = ^(int32(0))
}

// Java constructor: AliasStrings(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.JSChunkGraph chunkGraph, boolean outputStringUsage, com.google.javascript.jscomp.CompilerOptions.AliasStringsMode aliasStringsMode)
func NewAliasStrings(compiler AbstractCompiler, chunkGraph *JSChunkGraph, outputStringUsage bool, aliasStringsMode CompilerOptions_AliasStringsMode) *AliasStrings {
	clinit_AliasStrings()
	this_ := &AliasStrings{}
	InitAliasStrings(this_, compiler, chunkGraph, outputStringUsage, aliasStringsMode)
	return this_
}

func InitAliasStrings(this_ *AliasStrings, compiler AbstractCompiler, chunkGraph *JSChunkGraph, outputStringUsage bool, aliasStringsMode CompilerOptions_AliasStringsMode) {
	_ = compiler
	_ = chunkGraph
	_ = outputStringUsage
	_ = aliasStringsMode
	this_.initFields_AliasStrings()
	this_.compiler = compiler
	this_.chunkGraph = chunkGraph
	this_.outputStringUsage = outputStringUsage
	if !(aliasStringsMode != StaticGet_CompilerOptions_AliasStringsMode_NONE()) {
		j2gort.CheckState(false)
	}
	this_.aliasStringsMode = aliasStringsMode
}

// Java method: AliasStrings.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *AliasStrings) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	StaticGet_AliasStrings_Logger().Fine_1_string("Aliasing common strings")
	NodeTraversal_Traverse(this_.compiler, root, this_)
	this_.priv_AliasStrings_replaceStringsWithAliases_545()
	this_.priv_AliasStrings_addAliasDeclarationNodes_546()
	if this_.outputStringUsage {
		this_.priv_AliasStrings_outputStringUsage_549()
	}
}

// Java method: AliasStrings.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *AliasStrings) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(131) /* Token.TEMPLATELIT */, int32(130) /* Token.TAGGED_TEMPLATELIT */, int32(132) /* Token.TEMPLATELIT_SUB */ :
		return false
	case int32(28) /* Token.CALL */ :
		return !(ReplaceMessagesConstants_IsProtectedMessage(n))
	default:
		return true
	}
	panic("unreachable switch expression")
}

// Java method: AliasStrings.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *AliasStrings) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsStringLit() && !(parent.IsRegExp()) {
		var str string = n.GetString()
		_ = str
		if "undefined" == str {
			return
		}
		if (this_.aliasStringsMode == StaticGet_CompilerOptions_AliasStringsMode_LARGE()) && (j2gort.JStrLen(str) <= int32(100)) {
			return
		}
		var occurrence pkg_rhino.Node = n
		_ = occurrence
		var info *AliasStrings_StringInfo = this_.priv_AliasStrings_getOrCreateStringInfo_544(str)
		_ = info
		info.occurrences = append(info.occurrences, occurrence)
		var chunk *JSChunk = t.GetChunk()
		_ = chunk
		if int32(len(info.occurrences)) != int32(1) {
			if ((chunk != nil) && (info.chunkToContainDecl != nil)) && (chunk != info.chunkToContainDecl) {
				chunk = this_.chunkGraph.GetDeepestCommonDependencyInclusive_2(chunk, info.chunkToContainDecl)
			} else {
				return
			}
		}
		var varParent pkg_rhino.Node
		{
			bound_1 := this_.compiler
			varParent = j2gort.OrderedMapComputeIfAbsent(&this_.chunkVarParentMap, chunk, pkg_rhino.Node(0), func(mrp0_ *JSChunk) pkg_rhino.Node {
				return bound_1.GetNodeForCodeInsertion(mrp0_)
			})
		}

		_ = varParent
		{
			r_2 := info
			r_2.chunkToContainDecl = chunk
		}
		{
			r_3 := info
			r_3.parentForNewVarDecl = varParent
		}
		{
			r_4 := info
			r_4.siblingToInsertVarDeclBefore = varParent.GetFirstChild()
		}
	}
}

// Java method: AliasStrings.getOrCreateStringInfo(java.lang.String string)
func (this_ *AliasStrings) priv_AliasStrings_getOrCreateStringInfo_544(string_ string) *AliasStrings_StringInfo {
	_ = string_
	var info *AliasStrings_StringInfo = j2gort.CastPtr[AliasStrings_StringInfo](j2gort.MapGetStr(this_.stringInfoMap, string_))
	_ = info
	if info == nil {
		info = NewAliasStrings_StringInfo(this_, this_.stringInfoMap.Size())
		j2gort.MapPutStr(this_.stringInfoMap, string_, j2gort.PtrToAny(info))
	}
	return info
}

// Java method: AliasStrings.replaceStringsWithAliases()
func (this_ *AliasStrings) priv_AliasStrings_replaceStringsWithAliases_545() {
	{
		cur_5 := j2gort.NewMapEntryCursor(this_.stringInfoMap)
		for cur_5.NextEntry() {
			var literal string = cur_5.KeyStr()
			_ = literal
			var info *AliasStrings_StringInfo = j2gort.CastPtr[AliasStrings_StringInfo](cur_5.GetValue())
			_ = info
			if this_.priv_AliasStrings_shouldReplaceWithAlias_547(literal, info) {
				for _, node := range info.occurrences {
					_ = node
					this_.priv_AliasStrings_replaceStringWithAliasName_548(node, info.GetVariableName(literal), info)
				}

			}
		}
	}
}

// Java method: AliasStrings.addAliasDeclarationNodes()
func (this_ *AliasStrings) priv_AliasStrings_addAliasDeclarationNodes_546() {
	{
		cur_7 := j2gort.NewMapEntryCursor(this_.stringInfoMap)
		for cur_7.NextEntry() {
			var info *AliasStrings_StringInfo = j2gort.CastPtr[AliasStrings_StringInfo](cur_7.GetValue())
			_ = info
			if !(info.isAliased) {
				continue
			}
			var alias string = info.GetVariableName(cur_7.KeyStr())
			_ = alias
			var var_ pkg_rhino.Node = pkg_rhino.IR_Var_2(pkg_rhino.IR_Name(alias), pkg_rhino.IR_String(cur_7.KeyStr()))
			_ = var_
			var firstUse pkg_rhino.Node = info.occurrences[int32(0)]
			_ = firstUse
			var_.SrcrefTree(firstUse)
			if info.siblingToInsertVarDeclBefore == 0 {
				info.parentForNewVarDecl.AddChildToFront(var_)
			} else {
				var_.InsertBefore(info.siblingToInsertVarDeclBefore)
			}
			this_.compiler.ReportChangeToEnclosingScope(var_)
		}
	}
}

// Java method: AliasStrings.shouldReplaceWithAlias(java.lang.String str, com.google.javascript.jscomp.AliasStrings.StringInfo info)
func (this_ *AliasStrings) priv_AliasStrings_shouldReplaceWithAlias_547(str string, info *AliasStrings_StringInfo) bool {
	_ = str
	_ = info
	if this_.aliasStringsMode == StaticGet_CompilerOptions_AliasStringsMode_ALL_AGGRESSIVE() {
		return true
	}
	var sizeOfLiteral int32 = (int32(2) + j2gort.JStrLen(str))
	_ = sizeOfLiteral
	var count int32 = int32(len(info.occurrences))
	_ = count
	var sizeOfStrings int32 = (count * sizeOfLiteral)
	_ = sizeOfStrings
	var sizeOfVariable int32 = int32(3)
	_ = sizeOfVariable
	var sizeOfAliases int32 = (((int32(6) + sizeOfVariable) + sizeOfLiteral) + (count * sizeOfVariable))
	_ = sizeOfAliases
	return (sizeOfAliases < sizeOfStrings)
}

// Java method: AliasStrings.replaceStringWithAliasName(com.google.javascript.rhino.Node n, java.lang.String name, com.google.javascript.jscomp.AliasStrings.StringInfo info)
func (this_ *AliasStrings) priv_AliasStrings_replaceStringWithAliasName_548(n pkg_rhino.Node, name string, info *AliasStrings_StringInfo) {
	_ = n
	_ = name
	_ = info
	var nameNode pkg_rhino.Node = pkg_rhino.IR_Name(name)
	_ = nameNode
	n.ReplaceWith(nameNode)
	{
		r_8 := info
		r_8.isAliased = true
	}
	this_.compiler.ReportChangeToEnclosingScope(nameNode)
}

// Java method: AliasStrings.outputStringUsage()
func (this_ *AliasStrings) priv_AliasStrings_outputStringUsage_549() {
	var sb j2gort.LocalStringBuilder
	sb.InitStr("Strings used more than once:\n")

	{
		cur_9 := j2gort.NewMapEntryCursor(this_.stringInfoMap)
		for cur_9.NextEntry() {
			var info *AliasStrings_StringInfo = j2gort.CastPtr[AliasStrings_StringInfo](cur_9.GetValue())
			_ = info
			var count int32 = int32(len(info.occurrences))
			_ = count
			if count > int32(1) {
				sb.AppendInt(count)
				sb.AppendStr(": ")
				sb.AppendStr(cur_9.KeyStr())
				sb.AppendChar(uint16(10))
			}
		}
	}
	StaticGet_AliasStrings_Logger().Fine_1_string(sb.String())
}

// Java class: com.google.javascript.jscomp.AliasStrings.StringInfo
type AliasStrings_StringInfo struct {
	outer_                       *AliasStrings
	occurrences                  []pkg_rhino.Node `gson:"occurrences"`
	chunkToContainDecl           *JSChunk         `gson:"chunkToContainDecl"`
	aliasName                    string           `gson:"aliasName"`
	id                           int32            `gson:"id"`
	parentForNewVarDecl          pkg_rhino.Node   `gson:"parentForNewVarDecl"`
	siblingToInsertVarDeclBefore pkg_rhino.Node   `gson:"siblingToInsertVarDeclBefore"`
	idHash_                      atomic.Int32
	isAliased                    bool `gson:"isAliased"`
}

func (this_ *AliasStrings_StringInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *AliasStrings_StringInfo) IsJavaNil_() bool { return this_ == nil }
func (*AliasStrings_StringInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AliasStrings$StringInfo", "StringInfo"
}
func (this_ *AliasStrings_StringInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AliasStrings_StringInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AliasStrings.StringInfo", "com.google.javascript.jscomp.AliasStrings$StringInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*AliasStrings_StringInfo](v_)
	})
}
func (this_ *AliasStrings_StringInfo) SetOuter_(o *AliasStrings) { this_.outer_ = o }
func (this_ *AliasStrings_StringInfo) GetOuter_() *AliasStrings  { return this_.outer_ }

func (this_ *AliasStrings_StringInfo) initFields_AliasStrings_StringInfo() {
	_ = this_
	this_.occurrences = []pkg_rhino.Node{}
	this_.aliasName = j2gort.NullString
}

// Java constructor: StringInfo(int id)
func NewAliasStrings_StringInfo(outer_ *AliasStrings, id int32) *AliasStrings_StringInfo {
	this_ := &AliasStrings_StringInfo{}
	this_.outer_ = outer_
	InitAliasStrings_StringInfo(this_, id)
	return this_
}

func InitAliasStrings_StringInfo(this_ *AliasStrings_StringInfo, id int32) {
	_ = id
	this_.initFields_AliasStrings_StringInfo()
	this_.id = id
	this_.isAliased = false
}

// Java method: StringInfo.getVariableName(java.lang.String stringLiteral)
func (this_ *AliasStrings_StringInfo) GetVariableName(stringLiteral string) string {
	_ = stringLiteral
	if this_.aliasName == j2gort.NullString {
		this_.aliasName = this_.EncodeStringAsIdentifier("$$S_", stringLiteral)
	}
	return this_.aliasName
}

// Java method: StringInfo.encodeStringAsIdentifier(java.lang.String prefix, java.lang.String s)
func (this_ *AliasStrings_StringInfo) EncodeStringAsIdentifier(prefix string, s string) string {
	_ = prefix
	_ = s
	var maxLimit int32 = int32(20)
	_ = maxLimit
	var length int32 = j2gort.JStrLen(s)
	_ = length
	var limit int32 = min(length, maxLimit)
	_ = limit
	var sb j2gort.LocalStringBuilder

	sb.AppendStr(prefix)
	var protectHex bool = false
	_ = protectHex
	{
		var i int32 = int32(0)
		_ = i
		for ; i < limit; i = i + 1 {
			var ch uint16 = j2gort.JStrCharAt(s, i)
			_ = ch
			if protectHex {
				if ((int32(ch) >= int32(uint16(48))) && (int32(ch) <= int32(uint16(57)))) || ((int32(ch) >= int32(uint16(97))) && (int32(ch) <= int32(uint16(102)))) {
					sb.AppendChar(uint16(95))
				}
				protectHex = false
			}
			if (((int32(ch) >= int32(uint16(48))) && (int32(ch) <= int32(uint16(57)))) || ((int32(ch) >= int32(uint16(65))) && (int32(ch) <= int32(uint16(90))))) || ((int32(ch) >= int32(uint16(97))) && (int32(ch) <= int32(uint16(122)))) {
				sb.AppendChar(ch)
			} else {
				sb.AppendChar(uint16(36))
				sb.AppendStr(javalang.IntegerToHexString(int32(ch)))
				protectHex = true
			}
		}
	}
	if length == limit {
		return sb.String()
	}
	var hash int32 = (javalang.StringHashCode(j2gort.CheckNotNullString(s)) & this_.outer_.unitTestHashReductionMask)
	_ = hash
	sb.AppendChar(uint16(95))
	sb.AppendStr(javalang.IntegerToHexString(hash))
	var encoded string = sb.String()
	_ = encoded
	if !(this_.outer_.usedHashedAliases.Add(encoded)) {
		encoded = (j2gort.StrValString(encoded) + "_" + j2gort.StrValInt32(this_.id))
	}
	return encoded
}
func (*AliasStrings) CanReleaseSyntacticScopes_() {}
