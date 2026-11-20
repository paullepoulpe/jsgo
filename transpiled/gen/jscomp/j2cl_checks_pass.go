/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/J2clChecksPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.J2clChecksPass
type J2clChecksPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *J2clChecksPass) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clChecksPass) IsJavaNil_() bool { return this_ == nil }
func (*J2clChecksPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clChecksPass", "J2clChecksPass"
}
func (this_ *J2clChecksPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*J2clChecksPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clChecksPass", "com.google.javascript.jscomp.J2clChecksPass", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clChecksPass](v_)
	})
}
func (*J2clChecksPass) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*J2clChecksPass) Implements_jscomp_CompilerPass()                            {}
func (*J2clChecksPass) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_J2clChecksPass_J2CL_REFERENCE_EQUALITY *DiagnosticType

func StaticGet_J2clChecksPass_J2CL_REFERENCE_EQUALITY() *DiagnosticType {
	clinit_J2clChecksPass()
	return static_J2clChecksPass_J2CL_REFERENCE_EQUALITY
}
func StaticSet_J2clChecksPass_J2CL_REFERENCE_EQUALITY(v *DiagnosticType) {
	clinit_J2clChecksPass()
	static_J2clChecksPass_J2CL_REFERENCE_EQUALITY = v
}
func StaticGet_J2clChecksPass_INTEGER_TYPE_NAME() string { return "java.lang.Integer" }
func StaticGet_J2clChecksPass_INTEGER_FILE_NAME() string { return "java/lang/Integer.impl.java.js" }
func StaticGet_J2clChecksPass_FLOAT_TYPE_NAME() string   { return "java.lang.Float" }
func StaticGet_J2clChecksPass_FLOAT_FILE_NAME() string   { return "java/lang/Float.impl.java.js" }
func StaticGet_J2clChecksPass_LONG_TYPE_NAME() string    { return "goog.math.Long" }
func StaticGet_J2clChecksPass_LONG_FILE_NAME() string    { return "javascript/closure/math/long.js" }

var static_J2clChecksPass_REFERENCE_EQUALITY_TYPE_PATTERNS j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_J2clChecksPass_REFERENCE_EQUALITY_TYPE_PATTERNS() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_J2clChecksPass()
	return static_J2clChecksPass_REFERENCE_EQUALITY_TYPE_PATTERNS
}
func StaticSet_J2clChecksPass_REFERENCE_EQUALITY_TYPE_PATTERNS(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_J2clChecksPass()
	static_J2clChecksPass_REFERENCE_EQUALITY_TYPE_PATTERNS = v
}

var clinitDone_J2clChecksPass uint32
var clinitState_J2clChecksPass j2gort.ClinitState

func clinit_J2clChecksPass() {
	if atomic.LoadUint32(&clinitDone_J2clChecksPass) != 0 {
		return
	}
	clinitSlow_J2clChecksPass()
}

func clinitSlow_J2clChecksPass() {
	if !j2gort.EnterClinit(&clinitState_J2clChecksPass, &clinitDone_J2clChecksPass) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_J2clChecksPass, &clinitDone_J2clChecksPass)
	static_J2clChecksPass_J2CL_REFERENCE_EQUALITY = DiagnosticType_Warning("JSC_J2CL_REFERENCE_EQUALITY", "Reference equality may not be used with the specified type: {0}")
	static_J2clChecksPass_REFERENCE_EQUALITY_TYPE_PATTERNS = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_6("java.lang.Integer", "java/lang/Integer.impl.java.js", "java.lang.Float", "java/lang/Float.impl.java.js", "goog.math.Long", "javascript/closure/math/long.js")
}

func (this_ *J2clChecksPass) initFields_J2clChecksPass() {
	_ = this_
}

// Java constructor: J2clChecksPass(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewJ2clChecksPass(compiler AbstractCompiler) *J2clChecksPass {
	clinit_J2clChecksPass()
	this_ := &J2clChecksPass{}
	InitJ2clChecksPass(this_, compiler)
	return this_
}

func InitJ2clChecksPass(this_ *J2clChecksPass, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_J2clChecksPass()
	this_.compiler = compiler
}

// Java method: J2clChecksPass.visit(com.google.javascript.jscomp.NodeTraversal unused, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *J2clChecksPass) Visit_3(unused *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = unused
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(40) /* Token.SHEQ */, int32(4) /* Token.EQ */, int32(41) /* Token.SHNE */, int32(5) /* Token.NE */ :
	default:
		return
	}
	var firstJsType pkg_rhino.JSType = n.GetFirstChild().GetJSType()
	_ = firstJsType
	var lastJsType pkg_rhino.JSType = n.GetLastChild().GetJSType()
	_ = lastJsType
	if J2clChecksPass_IsNullType(firstJsType) || J2clChecksPass_IsNullType(lastJsType) {
		return
	}
	var firstSource string = J2clChecksPass_GetRestrictedJ2clSourceName(firstJsType)
	_ = firstSource
	var lastSource string = J2clChecksPass_GetRestrictedJ2clSourceName(lastJsType)
	_ = lastSource
	if (firstSource == j2gort.NullString) && (lastSource == j2gort.NullString) {
		return
	}
	this_.priv_J2clChecksPass_reportIfMatches_3466(n, firstSource, lastSource, "java.lang.Integer", "java/lang/Integer.impl.java.js")
	this_.priv_J2clChecksPass_reportIfMatches_3466(n, firstSource, lastSource, "java.lang.Float", "java/lang/Float.impl.java.js")
	this_.priv_J2clChecksPass_reportIfMatches_3466(n, firstSource, lastSource, "goog.math.Long", "javascript/closure/math/long.js")
}

// Java method: J2clChecksPass.reportIfMatches(com.google.javascript.rhino.Node n, java.lang.@org.jspecify.annotations.Nullable String firstSource, java.lang.@org.jspecify.annotations.Nullable String lastSource, java.lang.String typeName, java.lang.String fileName)
func (this_ *J2clChecksPass) priv_J2clChecksPass_reportIfMatches_3466(n pkg_rhino.Node, firstSource string, lastSource string, typeName string, fileName string) {
	_ = n
	_ = firstSource
	_ = lastSource
	_ = typeName
	_ = fileName
	if ((firstSource != j2gort.NullString) && javalang.StringEndsWith(j2gort.CheckNotNullString(firstSource), j2gort.CheckNotNullString(fileName))) || ((lastSource != j2gort.NullString) && javalang.StringEndsWith(j2gort.CheckNotNullString(lastSource), j2gort.CheckNotNullString(fileName))) {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_J2clChecksPass_J2CL_REFERENCE_EQUALITY(), []string{typeName}))
	}
}

// Java static method: J2clChecksPass.isNullType(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType jsType)
func J2clChecksPass_IsNullType(jsType pkg_rhino.JSType) bool {
	clinit_J2clChecksPass()
	_ = jsType
	if jsType == nil || jsType.IsJavaNil_() {
		return false
	}
	base_jsType_ := jsType.GetJSTypeBase()
	return (base_jsType_.IsNullType() || base_jsType_.IsVoidType())
}

// Java static method: J2clChecksPass.matchesAnyTargetPattern(java.lang.@org.jspecify.annotations.Nullable String sourceName)
func J2clChecksPass_MatchesAnyTargetPattern(sourceName string) bool {
	clinit_J2clChecksPass()
	_ = sourceName
	return ((sourceName != j2gort.NullString) && ((javalang.StringEndsWith(j2gort.CheckNotNullString(sourceName), "java/lang/Integer.impl.java.js") || javalang.StringEndsWith(j2gort.CheckNotNullString(sourceName), "java/lang/Float.impl.java.js")) || javalang.StringEndsWith(j2gort.CheckNotNullString(sourceName), "javascript/closure/math/long.js")))
}

// Java static method: J2clChecksPass.getRestrictedJ2clSourceName(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType jsType)
func J2clChecksPass_GetRestrictedJ2clSourceName(jsType pkg_rhino.JSType) string {
	clinit_J2clChecksPass()
	_ = jsType
	if jsType == nil || jsType.IsJavaNil_() {
		return j2gort.NullString
	}
	base_jsType_ := jsType.GetJSTypeBase()
	var union *pkg_rhino.UnionType
	_ = union
	var ok_1 bool
	union, ok_1 = j2gort.TryCast[*pkg_rhino.UnionType](jsType)
	if !ok_1 {
		if base_jsType_.IsNullType() || base_jsType_.IsVoidType() {
			return j2gort.NullString
		}
		var src string = J2clChecksPass_GetObjectSourceName(base_jsType_.ToMaybeObjectType())
		_ = src
		if J2clChecksPass_MatchesAnyTargetPattern(src) {
			return src
		} else {
			return j2gort.NullString
		}
	}
	var n int32 = union.GetAlternatesCount()
	_ = n
	var singleNonNull pkg_rhino.JSType = nil
	_ = singleNonNull
	var nonNullCount int32 = int32(0)
	_ = nonNullCount
	var hasNestedUnion bool = false
	_ = hasNestedUnion
	var allNonNullAreObjects bool = true
	_ = allNonNullAreObjects
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var alt pkg_rhino.JSType = union.GetAlternateAt(i)
			_ = alt
			if (alt.GetJSTypeBase().IsNullType() || alt.GetJSTypeBase().IsVoidType()) || alt.GetJSTypeBase().IsNoType() {
				continue
			}
			if alt.GetJSTypeBase().IsUnionType() {
				hasNestedUnion = true
				break
			}
			nonNullCount = nonNullCount + 1
			if nonNullCount == int32(1) {
				singleNonNull = alt
			}
			if j2gort.IfaceEquals[pkg_rhino.ObjectType](alt.GetJSTypeBase().ToMaybeObjectType(), nil) {
				allNonNullAreObjects = false
			}
		}
	}
	if hasNestedUnion {
		var src_2 string = J2clChecksPass_GetObjectSourceName(union.RestrictByNotNullOrUndefined().GetJSTypeBase().ToMaybeObjectType())
		_ = src_2
		if J2clChecksPass_MatchesAnyTargetPattern(src_2) {
			return src_2
		} else {
			return j2gort.NullString
		}
	}
	if nonNullCount == int32(0) {
		return j2gort.NullString
	}
	if nonNullCount == int32(1) {
		var src_3 string = J2clChecksPass_GetObjectSourceName(singleNonNull.GetJSTypeBase().ToMaybeObjectType())
		_ = src_3
		if J2clChecksPass_MatchesAnyTargetPattern(src_3) {
			return src_3
		} else {
			return j2gort.NullString
		}
	}
	if !(allNonNullAreObjects) {
		return j2gort.NullString
	}
	var hasJ2clCandidate bool = false
	_ = hasJ2clCandidate
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < n; i_2 = i_2 + 1 {
			var alt_2 pkg_rhino.JSType = union.GetAlternateAt(i_2)
			_ = alt_2
			if ((!(alt_2.GetJSTypeBase().IsNullType()) && !(alt_2.GetJSTypeBase().IsVoidType())) && !(alt_2.GetJSTypeBase().IsNoType())) && J2clChecksPass_MatchesAnyTargetPattern(J2clChecksPass_GetObjectSourceName(alt_2.GetJSTypeBase().ToMaybeObjectType())) {
				hasJ2clCandidate = true
				break
			}
		}
	}
	if !(hasJ2clCandidate) {
		return j2gort.NullString
	}
	var src_4 string = J2clChecksPass_GetObjectSourceName(union.RestrictByNotNullOrUndefined().GetJSTypeBase().ToMaybeObjectType())
	_ = src_4
	if J2clChecksPass_MatchesAnyTargetPattern(src_4) {
		return src_4
	} else {
		return j2gort.NullString
	}
}

// Java static method: J2clChecksPass.getObjectSourceName(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType objType)
func J2clChecksPass_GetObjectSourceName(objType pkg_rhino.ObjectType) string {
	clinit_J2clChecksPass()
	_ = objType
	if objType == nil || objType.IsJavaNil_() {
		return j2gort.NullString
	}
	var constructor *pkg_rhino.FunctionType = objType.GetConstructor()
	_ = constructor
	if constructor == nil {
		return j2gort.NullString
	}
	return NodeUtil_GetSourceName(constructor.GetSource())
}

// Java method: J2clChecksPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *J2clChecksPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}
func (*J2clChecksPass) AlwaysShouldTraverse_()      {}
func (*J2clChecksPass) CanReleaseSyntacticScopes_() {}
