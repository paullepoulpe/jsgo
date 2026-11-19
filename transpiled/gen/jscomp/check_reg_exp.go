/*
 * Copyright 2010 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CheckRegExp.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_regex "github.com/paullepoulpe/jsgo/transpiled/gen/regex"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CheckRegExp
type CheckRegExp struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler                           AbstractCompiler `gson:"compiler"`
	reportErrors                       bool             `gson:"reportErrors"`
	globalRegExpPropertiesUsed         bool             `gson:"globalRegExpPropertiesUsed"`
	assumeAllGlobalRegexpUsagesVisible bool             `gson:"assumeAllGlobalRegexpUsagesVisible"`
}

func (this_ *CheckRegExp) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckRegExp) IsJavaNil_() bool { return this_ == nil }
func (*CheckRegExp) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckRegExp", "CheckRegExp"
}
func (this_ *CheckRegExp) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckRegExp) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckRegExp", "com.google.javascript.jscomp.CheckRegExp", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckRegExp](v_)
	})
}
func (*CheckRegExp) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckRegExp) Implements_jscomp_CompilerPass()                            {}
func (*CheckRegExp) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckRegExp_REGEXP_REFERENCE *DiagnosticType

func StaticGet_CheckRegExp_REGEXP_REFERENCE() *DiagnosticType {
	clinit_CheckRegExp()
	return static_CheckRegExp_REGEXP_REFERENCE
}
func StaticSet_CheckRegExp_REGEXP_REFERENCE(v *DiagnosticType) {
	clinit_CheckRegExp()
	static_CheckRegExp_REGEXP_REFERENCE = v
}

var static_CheckRegExp_MALFORMED_REGEXP *DiagnosticType

func StaticGet_CheckRegExp_MALFORMED_REGEXP() *DiagnosticType {
	clinit_CheckRegExp()
	return static_CheckRegExp_MALFORMED_REGEXP
}
func StaticSet_CheckRegExp_MALFORMED_REGEXP(v *DiagnosticType) {
	clinit_CheckRegExp()
	static_CheckRegExp_MALFORMED_REGEXP = v
}

var static_CheckRegExp_REGEXP_PROPERTY_SKIPLIST j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_CheckRegExp_REGEXP_PROPERTY_SKIPLIST() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CheckRegExp()
	return static_CheckRegExp_REGEXP_PROPERTY_SKIPLIST
}
func StaticSet_CheckRegExp_REGEXP_PROPERTY_SKIPLIST(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_CheckRegExp()
	static_CheckRegExp_REGEXP_PROPERTY_SKIPLIST = v
}

var clinitDone_CheckRegExp uint32
var clinitState_CheckRegExp j2gort.ClinitState

func clinit_CheckRegExp() {
	if atomic.LoadUint32(&clinitDone_CheckRegExp) != 0 {
		return
	}
	clinitSlow_CheckRegExp()
}

func clinitSlow_CheckRegExp() {
	if !j2gort.EnterClinit(&clinitState_CheckRegExp, &clinitDone_CheckRegExp) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckRegExp, &clinitDone_CheckRegExp)
	static_CheckRegExp_REGEXP_REFERENCE = DiagnosticType_Warning("JSC_REGEXP_REFERENCE", "References to the global RegExp object prevents optimization of regular expressions.")
	static_CheckRegExp_MALFORMED_REGEXP = DiagnosticType_Warning("JSC_MALFORMED_REGEXP", "Malformed Regular Expression: {0}")
	static_CheckRegExp_REGEXP_PROPERTY_SKIPLIST = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("$1", "$2", "$3", "$4", "$5", "$6", []any{"$7", "$8", "$9", "$_", "$input", "input", "lastMatch", "lastParen", "leftContext", "rightContext", "global", "ignoreCase", "lastIndex", "multiline", "source"})
}

func (this_ *CheckRegExp) initFields_CheckRegExp() {
	_ = this_

}

// Java constructor: CheckRegExp(com.google.javascript.jscomp.AbstractCompiler compiler, boolean assumeAllGlobalRegexpUsagesVisible, boolean reportErrors)
func NewCheckRegExp(compiler AbstractCompiler, assumeAllGlobalRegexpUsagesVisible bool, reportErrors bool) *CheckRegExp {
	clinit_CheckRegExp()
	this_ := &CheckRegExp{}
	InitCheckRegExp(this_, compiler, assumeAllGlobalRegexpUsagesVisible, reportErrors)
	return this_
}

func InitCheckRegExp(this_ *CheckRegExp, compiler AbstractCompiler, assumeAllGlobalRegexpUsagesVisible bool, reportErrors bool) {
	_ = compiler
	_ = assumeAllGlobalRegexpUsagesVisible
	_ = reportErrors
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckRegExp()
	this_.compiler = compiler
	this_.assumeAllGlobalRegexpUsagesVisible = assumeAllGlobalRegexpUsagesVisible
	this_.reportErrors = reportErrors
}

// Java method: CheckRegExp.isGlobalRegExpPropertiesUsed()
func (this_ *CheckRegExp) IsGlobalRegExpPropertiesUsed() bool {
	if this_.assumeAllGlobalRegexpUsagesVisible {
		return this_.globalRegExpPropertiesUsed
	}
	return true
}

// Java method: CheckRegExp.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckRegExp) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckRegExp.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckRegExp) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if NodeUtil_IsReferenceName(n) {
		var name string = n.GetString()
		_ = name
		if (j2gort.CheckNotNullString(name) == "RegExp") && (j2gort.CastPtr[Var](t.GetScope().GetVar(name)) == nil) {
			var parentType pkg_rhino.Token = parent.GetToken()
			_ = parentType
			var first bool = (n == parent.GetFirstChild())
			_ = first
			if !((((((((((parentType == pkg_rhino.StaticGet_Token_NEW()) && first) || ((parentType == pkg_rhino.StaticGet_Token_CALL()) && first)) || ((parentType == pkg_rhino.StaticGet_Token_INSTANCEOF()) && !(first))) || (parentType == pkg_rhino.StaticGet_Token_EQ())) || (parentType == pkg_rhino.StaticGet_Token_NE())) || (parentType == pkg_rhino.StaticGet_Token_SHEQ())) || (parentType == pkg_rhino.StaticGet_Token_SHNE())) || (parentType == pkg_rhino.StaticGet_Token_CASE())) || (((parentType == pkg_rhino.StaticGet_Token_GETPROP()) && first) && !(j2gort.SetContainsStr(StaticGet_CheckRegExp_REGEXP_PROPERTY_SKIPLIST(), parent.GetString())))) {
				if this_.reportErrors {
					t.Report_3(n, StaticGet_CheckRegExp_REGEXP_REFERENCE(), []string{})
				}
				this_.globalRegExpPropertiesUsed = true
			}
		}
	} else if this_.reportErrors && n.IsRegExp() {
		var pattern string = n.GetFirstChild().GetString()
		_ = pattern
		var flags string
		_ = flags
		if n.HasTwoChildren() {
			flags = n.GetLastChild().GetString()
		} else {
			flags = ""
		}
		func() {
			defer func() {
				if rec_1 := recover(); rec_1 != nil {
					jex_2 := javalang.RecoverException(rec_1)
					_ = jex_2
					if j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_2) || j2gort.InstanceOf[javalang.IIndexOutOfBoundsException](jex_2) {
						ex := jex_2
						_ = ex
						t.Report_3(n, StaticGet_CheckRegExp_MALFORMED_REGEXP(), []string{j2gort.ThrowableGetMessage(ex)})
					} else {
						panic(rec_1)
					}
				}
			}()
			pkg_regex.RegExpTree_ParseRegExp(pattern, flags)
		}()
	}
}
func (*CheckRegExp) AlwaysShouldTraverse_()      {}
func (*CheckRegExp) CanReleaseSyntacticScopes_() {}
