/*
 * Copyright 2026 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/J2clSourceUtils.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.J2clSourceUtils
type J2clSourceUtils struct {
	idHash_ atomic.Int32
}

func (this_ *J2clSourceUtils) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clSourceUtils) IsJavaNil_() bool { return this_ == nil }
func (*J2clSourceUtils) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.J2clSourceUtils", "J2clSourceUtils" }
func (this_ *J2clSourceUtils) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*J2clSourceUtils) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clSourceUtils", "com.google.javascript.jscomp.J2clSourceUtils", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clSourceUtils](v_)
	})
}

func (this_ *J2clSourceUtils) initFields_J2clSourceUtils() {
	_ = this_
}

// Java constructor: J2clSourceUtils()
func NewJ2clSourceUtils() *J2clSourceUtils {
	this_ := &J2clSourceUtils{}
	InitJ2clSourceUtils(this_)
	return this_
}

func InitJ2clSourceUtils(this_ *J2clSourceUtils) {
	this_.initFields_J2clSourceUtils()
}

// Java static method: J2clSourceUtils.isJ2clSource(java.lang.String sourceFileName)
func J2clSourceUtils_IsJ2clSource_1_string(sourceFileName string) bool {
	_ = sourceFileName
	return ((sourceFileName != j2gort.NullString) && javalang.StringEndsWith(j2gort.CheckNotNullString(sourceFileName), ".java.js"))
}

// Java static method: J2clSourceUtils.isJ2clSource(com.google.javascript.rhino.Node scriptNode)
func J2clSourceUtils_IsJ2clSource_1_Node(scriptNode pkg_rhino.Node) bool {
	_ = scriptNode
	if !(scriptNode.IsScript()) { j2gort.CheckArgument(false) }
	return J2clSourceUtils_IsJ2clSource_1_string(scriptNode.GetSourceFileName())
}

// Java static method: J2clSourceUtils.isInsideJ2clSource(com.google.javascript.rhino.Node node)
func J2clSourceUtils_IsInsideJ2clSource(node pkg_rhino.Node) bool {
	_ = node
	var scriptNode pkg_rhino.Node
	_ = scriptNode
	if node.IsScript() {
		scriptNode = node
	} else {
		scriptNode = NodeUtil_GetEnclosingScript(node)
	}
	if ((scriptNode == 0)) {
		return false
	}
	return J2clSourceUtils_IsJ2clSource_1_Node(scriptNode)
}

// Java static method: J2clSourceUtils.isJ2clHeaderSource(java.lang.String sourceFileName)
func J2clSourceUtils_IsJ2clHeaderSource(sourceFileName string) bool {
	_ = sourceFileName
	return (((sourceFileName != j2gort.NullString) && javalang.StringEndsWith(j2gort.CheckNotNullString(sourceFileName), ".java.js")) && !(javalang.StringEndsWith(j2gort.CheckNotNullString(sourceFileName), ".impl.java.js")))
}

