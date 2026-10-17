// Port of com/google/javascript/jscomp/AliasStrings.java @ closure-compiler v20260415

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AliasStrings is a compiler pass for aliasing strings. String declarations
// contribute to garbage collection, which becomes a problem in large
// applications. Strings that should be aliased occur many times in the code,
// or occur on codepaths that get executed frequently.
//
// 2017/09/17 Notes: - Turning on this pass usually hurts code size after
// gzip. - It was originally written to deal with performance problems on
// some older browser VMs. - However, projects that make heavy use of
// jslayout may need to enable this pass even for modern browsers, because
// jslayout generates so many duplicate strings.
type AliasStrings struct {
	compiler AbstractCompiler

	chunkGraph *JSChunkGraph

	outputStringUsage bool

	stringInfoMap *jutil.TreeMap[string, *AliasStrings_StringInfo]

	usedHashedAliases *jutil.LinkedHashSet[string]

	aliasStringsMode CompilerOptions_AliasStringsMode

	// Map from chunk to the node in that chunk that should parent any string
	// variable declarations that have to be moved into that chunk
	chunkVarParentMap *jutil.LinkedHashMap[*JSChunk, *rhino.Node]

	// package private. This value is AND-ed with the hash function to allow
	// unit tests to reduce the range of hash values to test collision cases
	unitTestHashReductionMask int32
}

var aliasStrings_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.AliasStrings")

// Prefix for variable names for the aliased strings
const AliasStrings_STRING_ALIAS_PREFIX = "$$S_"

// Alias strings longer than 100 characters if
// aliasStringsMode=AliasStringsMode.LARGE
const AliasStrings_ALIAS_LARGE_STRINGS_LENGTH = 100

// NewAliasStrings creates an instance.
//
// compiler is the compiler; chunkGraph is the chunk graph, or null if there
// are no chunks; outputStringUsage outputs all strings and the number of
// times they were used in the application to the server log;
// aliasStringsMode is the alias strings policy set to either ALL or LARGE.
//
// Java: AliasStrings(AbstractCompiler, JSChunkGraph, boolean, AliasStringsMode)
func NewAliasStrings(
	compiler AbstractCompiler,
	chunkGraph *JSChunkGraph,
	outputStringUsage bool,
	aliasStringsMode CompilerOptions_AliasStringsMode) *AliasStrings {
	a := &AliasStrings{
		stringInfoMap:             jutil.NewTreeMap[string, *AliasStrings_StringInfo](jutil.StringCompareTo),
		usedHashedAliases:         jutil.NewLinkedHashSet[string](),
		chunkVarParentMap:         jutil.NewLinkedHashMap[*JSChunk, *rhino.Node](),
		unitTestHashReductionMask: ^0,
	}
	a.compiler = compiler
	a.chunkGraph = chunkGraph
	a.outputStringUsage = outputStringUsage
	jutil.CheckState(aliasStringsMode != CompilerOptions_AliasStringsMode_NONE)
	a.aliasStringsMode = aliasStringsMode
	return a
}

// Java: AliasStrings.process(Node, Node)
func (a *AliasStrings) Process(externs *rhino.Node, root *rhino.Node) {
	aliasStrings_logger.Fine("Aliasing common strings")

	// Traverse the tree and collect strings
	NodeTraversal_Traverse(a.compiler, root, a)

	// 1st edit pass: replace some strings with aliases
	a.ReplaceStringsWithAliases()

	// 2nd edit pass: add variable declarations for aliased strings.
	a.AddAliasDeclarationNodes()

	if a.outputStringUsage {
		a.OutputStringUsage()
	}
}

// Java: AliasStrings.shouldTraverse(NodeTraversal, Node, Node)
func (a *AliasStrings) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_TEMPLATELIT,
		rhino.Token_TAGGED_TEMPLATELIT,
		rhino.Token_TEMPLATELIT_SUB: // technically redundant, since it must be a child of the others
		// TODO(bradfordcsmith): Consider replacing long and/or frequently occurring substrings
		// within template literals with template substitutions.
		return false
	case rhino.Token_CALL:
		return !ReplaceMessagesConstants_IsProtectedMessage(n)
	default:
		return true
	}
}

// Java: AliasStrings.visit(NodeTraversal, Node, Node)
func (a *AliasStrings) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsStringLit() && !parent.IsRegExp() {
		str := n.GetString()

		// "undefined" is special-cased, since it needs to be used when JS code
		// is unloading and therefore variable references aren't available.
		// This is because of a bug in Firefox.
		if "undefined" == str {
			return
		}

		if a.aliasStringsMode == CompilerOptions_AliasStringsMode_LARGE &&
			jutil.StringLength(str) <= AliasStrings_ALIAS_LARGE_STRINGS_LENGTH {
			return
		}

		occurrence := n
		info := a.GetOrCreateStringInfo(str)

		info.occurrences.Add(occurrence)

		// The current chunk.
		chunk := t.GetChunk()
		if info.occurrences.Size() != 1 {
			// Check whether the current chunk depends on the chunk containing
			// the declaration.
			if chunk != nil && info.chunkToContainDecl != nil && chunk != info.chunkToContainDecl {
				// We need to declare this string in the deepest chunk in the
				// chunk dependency graph that both of these chunks depend on.
				chunk = a.chunkGraph.GetDeepestCommonDependencyInclusive(chunk, info.chunkToContainDecl)
			} else {
				// use the previously saved insertion location.
				return
			}
		}
		varParent := a.chunkVarParentMap.ComputeIfAbsent(chunk, a.compiler.GetNodeForCodeInsertion)
		info.chunkToContainDecl = chunk
		info.parentForNewVarDecl = varParent
		info.siblingToInsertVarDeclBefore = varParent.GetFirstChild()
	}
}

// GetOrCreateStringInfo looks up the StringInfo object for a JavaScript
// string. Creates it if necessary.
//
// Java: AliasStrings.getOrCreateStringInfo(String)
func (a *AliasStrings) GetOrCreateStringInfo(string_ string) *AliasStrings_StringInfo {
	info := a.stringInfoMap.Get(string_)
	if info == nil {
		info = NewAliasStrings_StringInfo(a, a.stringInfoMap.Size())
		a.stringInfoMap.Put(string_, info)
	}
	return info
}

// ReplaceStringsWithAliases replaces strings with references to alias
// variables.
//
// Java: AliasStrings.replaceStringsWithAliases()
func (a *AliasStrings) ReplaceStringsWithAliases() {
	for literal, info := range a.stringInfoMap.All() {
		if a.ShouldReplaceWithAlias(literal, info) {
			for _, node := range info.occurrences.Slice() {
				a.ReplaceStringWithAliasName(node, info.GetVariableName(literal), info)
			}
		}
	}
}

// AddAliasDeclarationNodes creates a var declaration for each aliased
// string. Var declarations are inserted as close to the first use of the
// string as possible.
//
// Java: AliasStrings.addAliasDeclarationNodes()
func (a *AliasStrings) AddAliasDeclarationNodes() {
	for key, info := range a.stringInfoMap.All() {
		if !info.isAliased {
			continue
		}
		alias := info.GetVariableName(key)
		var_ := rhino.IR_VarWithValue(rhino.IR_Name(alias), rhino.IR_String(key))
		firstUse := info.occurrences.Get(0)
		var_.SrcrefTree(firstUse)
		if info.siblingToInsertVarDeclBefore == nil {
			info.parentForNewVarDecl.AddChildToFront(var_)
		} else {
			var_.InsertBefore(info.siblingToInsertVarDeclBefore)
		}
		a.compiler.ReportChangeToEnclosingScope(var_)
	}
}

// ShouldReplaceWithAlias dictates the policy for replacing a string with an
// alias.
//
// str is the string literal; info is the accumulated information about a
// string.
//
// Java: AliasStrings.shouldReplaceWithAlias(String, StringInfo)
func (a *AliasStrings) ShouldReplaceWithAlias(str string, info *AliasStrings_StringInfo) bool {
	// Always alias strings if the mode is ALL_AGGRESSIVE.
	if a.aliasStringsMode == CompilerOptions_AliasStringsMode_ALL_AGGRESSIVE {
		return true
	}

	// Optimize for code size.  Are aliases smaller than strings?
	//
	// This logic optimizes for the size of uncompressed code, but it tends to
	// get good results for the size of the gzipped code too.
	//
	// gzip actually prefers that strings are not aliased - it compresses N
	// string literals better than 1 string literal and N+1 short variable
	// names, provided each string is within 32k of the previous copy.  We
	// follow the uncompressed logic as insurance against there being multiple
	// strings more than 32k apart.

	sizeOfLiteral := 2 + jutil.StringLength(str)
	count := info.occurrences.Size()
	sizeOfStrings := count * sizeOfLiteral
	sizeOfVariable := 3
	//  '6' comes from: 'var =;' in var XXX="...";
	sizeOfAliases := 6 +
		sizeOfVariable +
		sizeOfLiteral + // declaration
		count*sizeOfVariable // + uses

	return sizeOfAliases < sizeOfStrings
}

// ReplaceStringWithAliasName replaces a string literal with a reference to
// the string's alias variable.
//
// Java: AliasStrings.replaceStringWithAliasName(Node, String, StringInfo)
func (a *AliasStrings) ReplaceStringWithAliasName(n *rhino.Node, name string, info *AliasStrings_StringInfo) {
	nameNode := rhino.IR_Name(name)
	n.ReplaceWith(nameNode)
	info.isAliased = true
	a.compiler.ReportChangeToEnclosingScope(nameNode)
}

// OutputStringUsage outputs a log of all strings used more than once in the
// code.
//
// Java: AliasStrings.outputStringUsage()
func (a *AliasStrings) OutputStringUsage() {
	sb := jutil.NewStringBuilder()
	sb.Append("Strings used more than once:\n")
	for key, info := range a.stringInfoMap.All() {
		count := info.occurrences.Size()
		if count > 1 {
			sb.AppendInt(int64(count))
			sb.Append(": ")
			sb.Append(key)
			sb.AppendChar('\n')
		}
	}
	// TODO(user): Make this save to file OR output to the application
	aliasStrings_logger.Fine(sb.String())
}

// -------------------------------------------------------------------------

// AliasStrings_StringInfo is a class that holds information about a
// JavaScript string that might become aliased.
type AliasStrings_StringInfo struct {
	outer *AliasStrings

	id int

	isAliased bool // set to 'true' when reference to alias created

	occurrences *jutil.ArrayList[*rhino.Node]

	chunkToContainDecl           *JSChunk
	parentForNewVarDecl          *rhino.Node
	siblingToInsertVarDeclBefore *rhino.Node

	aliasName *string
}

// Java: AliasStrings.StringInfo(int)
func NewAliasStrings_StringInfo(outer *AliasStrings, id int) *AliasStrings_StringInfo {
	s := &AliasStrings_StringInfo{outer: outer, occurrences: jutil.NewArrayList[*rhino.Node]()}
	s.id = id
	s.isAliased = false
	return s
}

// GetVariableName returns the JS variable name to be substituted for this
// string.
//
// Java: AliasStrings.StringInfo.getVariableName(String)
func (s *AliasStrings_StringInfo) GetVariableName(stringLiteral string) string {
	if s.aliasName == nil {
		aliasName := s.EncodeStringAsIdentifier(AliasStrings_STRING_ALIAS_PREFIX, stringLiteral)
		s.aliasName = &aliasName
	}
	return *s.aliasName
}

// EncodeStringAsIdentifier returns a legal identifier that uniquely
// characterizes string 's'.
//
// We want the identifier to be a function of the string value because that
// makes the identifiers stable as the program is changed.
//
// The digits of a good hash function would be adequate, but for short
// strings the following algorithm is easier to work with for unit tests.
//
// ASCII alphanumerics are mapped to themselves. Other characters are mapped
// to $XXX or $XXX_ where XXX is a variable number of hex digits. The
// underscore is inserted as necessary to avoid ambiguity when the character
// following is a hex digit. E.g. '\n1' maps to '$a_1', distinguished by the
// underscore from '¡' which maps to '$a1'.
//
// If the string is short enough, this is sufficient. Longer strings are
// truncated after encoding an initial prefix and appended with a hash value.
//
// Java: AliasStrings.StringInfo.encodeStringAsIdentifier(String, String)
func (s *AliasStrings_StringInfo) EncodeStringAsIdentifier(prefix string, str string) string {
	// Limit to avoid generating very long identifiers
	const maxLimit = 20
	length := jutil.StringLength(str)
	limit := min(length, maxLimit)

	sb := jutil.NewStringBuilder()
	sb.Append(prefix)
	protectHex := false

	for i := 0; i < limit; i++ {
		ch := jutil.StringCharAt(str, i)

		if protectHex {
			if (ch >= '0' && ch <= '9') ||
				(ch >= 'a' && ch <= 'f') { // toHexString generate lowercase
				sb.AppendChar('_')
			}
			protectHex = false
		}

		if (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
			sb.AppendChar(ch)
		} else {
			sb.AppendChar('$')
			sb.Append(jutil.IntegerToHexString(int32(ch)))
			protectHex = true
		}
	}

	if length == limit {
		return sb.String()
	}

	// The identifier is not unique because we omitted part, so add a
	// checksum as a hashcode.
	hash := jutil.StringHashCode(str) & s.outer.unitTestHashReductionMask
	sb.AppendChar('_')
	sb.Append(jutil.IntegerToHexString(hash))
	encoded := sb.String()
	if !s.outer.usedHashedAliases.Add(encoded) {
		// A collision has been detected (which is very rare). Use the sequence
		// id to break the tie. This means that the name is no longer invariant
		// across source code changes and recompilations.
		encoded += "_" + jutil.ToJavaString(s.id)
	}
	return encoded
}
