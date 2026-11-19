/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/ParsingUtil.java (generated mechanically by j2go).

package parsing

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.ParsingUtil
type ParsingUtil struct {
	idHash_ atomic.Int32
}

func (this_ *ParsingUtil) BindSelf_(s any) {
	_ = s
}
func (this_ *ParsingUtil) IsJavaNil_() bool { return this_ == nil }
func (*ParsingUtil) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.ParsingUtil", "ParsingUtil"
}
func (this_ *ParsingUtil) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ParsingUtil) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.ParsingUtil", "com.google.javascript.jscomp.parsing.ParsingUtil", func(v_ any) bool {
		return j2gort.InstanceOf[*ParsingUtil](v_)
	})
}

func (this_ *ParsingUtil) initFields_ParsingUtil() {
	_ = this_
}

// Java constructor: ParsingUtil()
func NewParsingUtil() *ParsingUtil {
	this_ := &ParsingUtil{}
	InitParsingUtil(this_)
	return this_
}

func InitParsingUtil(this_ *ParsingUtil) {
	this_.initFields_ParsingUtil()
}

// Java static method: ParsingUtil.getParamOrPatternNames(com.google.javascript.rhino.Node n, java.util.function.Consumer<com.google.javascript.rhino.Node> cb)
func ParsingUtil_GetParamOrPatternNames(n pkg_rhino.Node, cb j2gort.Ext_java_util_function_Consumer) {
	_ = n
	_ = cb
	j2gort.CheckNotNullNode(n)
	switch n.TokenSwitchOrd_() {
	case int32(88) /* Token.EMPTY */, int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */ :
	case int32(32) /* Token.NAME */ :
		cb.Accept(j2gort.BoxNode(n))
	case int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */, int32(100) /* Token.STRING_KEY */, int32(134) /* Token.DEFAULT_VALUE */ :
		ParsingUtil_GetParamOrPatternNames(n.GetFirstChild(), cb)
	case int32(129) /* Token.COMPUTED_PROP */ :
		ParsingUtil_GetParamOrPatternNames(n.GetSecondChild(), cb)
	case int32(103) /* Token.OBJECT_PATTERN */, int32(102) /* Token.ARRAY_PATTERN */, int32(49) /* Token.PARAM_LIST */ :
		{
			var c pkg_rhino.Node = n.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				ParsingUtil_GetParamOrPatternNames(c, cb)
			}
		}
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected parameter structure")))
	}
}
