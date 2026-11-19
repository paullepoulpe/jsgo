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
// Port of com/google/javascript/refactoring/JsFlumeCallback.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.refactoring.JsFlumeCallback
type JsFlumeCallback struct {
	scanner            Scanner                             `gson:"scanner"`
	includeFilePattern *j2gort.Ext_java_util_regex_Pattern `gson:"includeFilePattern"`
	matches            j2gort.Ext_java_util_List           `gson:"matches"`
	fixes              j2gort.Ext_java_util_List           `gson:"fixes"`
	idHash_            atomic.Int32
}

func (this_ *JsFlumeCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *JsFlumeCallback) IsJavaNil_() bool { return this_ == nil }
func (*JsFlumeCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.JsFlumeCallback", "JsFlumeCallback"
}
func (this_ *JsFlumeCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsFlumeCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.JsFlumeCallback", "com.google.javascript.refactoring.JsFlumeCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*JsFlumeCallback](v_)
	})
}
func (*JsFlumeCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *JsFlumeCallback) initFields_JsFlumeCallback() {
	_ = this_
	this_.matches = j2gort.NewExt_java_util_ArrayList()
	this_.fixes = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: JsFlumeCallback(com.google.javascript.refactoring.Scanner scanner, java.util.regex.@org.jspecify.annotations.Nullable Pattern includeFilePattern)
func NewJsFlumeCallback(scanner Scanner, includeFilePattern *j2gort.Ext_java_util_regex_Pattern) *JsFlumeCallback {
	this_ := &JsFlumeCallback{}
	InitJsFlumeCallback(this_, scanner, includeFilePattern)
	return this_
}

func InitJsFlumeCallback(this_ *JsFlumeCallback, scanner Scanner, includeFilePattern *j2gort.Ext_java_util_regex_Pattern) {
	_ = scanner
	_ = includeFilePattern
	this_.initFields_JsFlumeCallback()
	this_.scanner = scanner
	this_.includeFilePattern = includeFilePattern
}

// Java method: JsFlumeCallback.getMatches()
func (this_ *JsFlumeCallback) GetMatches() j2gort.Ext_java_util_List {
	return this_.matches
}

// Java method: JsFlumeCallback.getFixes()
func (this_ *JsFlumeCallback) GetFixes() j2gort.Ext_java_util_List {
	return this_.fixes
}

// Java method: JsFlumeCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *JsFlumeCallback) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if n.IsFromExterns() {
		return false
	}
	if n.IsScript() {
		var filename string = n.GetSourceFileName()
		_ = filename
		if ((this_.includeFilePattern != nil) && !(j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(this_.includeFilePattern.Pattern()))) && !(j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(filename)) {
			return this_.includeFilePattern.Matcher(j2gort.BoxString(filename)).Find_0()
		}
	}
	return true
}

// Java method: JsFlumeCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *JsFlumeCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var metadata *NodeMetadata = NodeMetadata_FromTraversal(t)
	_ = metadata
	if this_.scanner.Matches(n, metadata) {
		var val_1 Match
		InitMatch(&val_1, n, metadata)
		var match *Match = &val_1
		_ = match
		this_.fixes.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.scanner.ProcessMatch(match)))
		this_.matches.Add_1(match)
	}
}
