// Port of com/google/javascript/jscomp/ChromeCodingConvention.java @ closure-compiler v20261021

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

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino"

// ChromeCodingConvention is the coding convention used by the Chrome team to
// compile Chrome's JS.
type ChromeCodingConvention struct {
	CodingConventions_ProxyBase
}

var chromeAddSingletonGetterQ = rhino.QualifiedName_Of("cr.addSingletonGetter")

func NewChromeCodingConvention() *ChromeCodingConvention {
	return NewChromeCodingConventionWithConvention(CodingConventions_GetDefault())
}

func NewChromeCodingConventionWithConvention(wrapped CodingConvention) *ChromeCodingConvention {
	c := &ChromeCodingConvention{}
	c.CodingConventions_ProxyBase = newProxyBase(c, wrapped)
	return c
}

func (c *ChromeCodingConvention) GetSingletonGetterClassName(callNode rhino.Node) *string {
	callArg := callNode.GetFirstChild()
	if !chromeAddSingletonGetterQ.Matches(callArg) || !callNode.HasTwoChildren() {
		return c.CodingConventions_ProxyBase.GetSingletonGetterClassName(callNode)
	}
	return callArg.GetNext().GetQualifiedName()
}

func (c *ChromeCodingConvention) ApplySingletonGetter(classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType) {
	defSite := classType.Constructor().GetSource()
	classType.DeclareConstructorProperty("getInstance", getterType, defSite)
	classType.DeclareConstructorProperty("instance_", classType.Instance(), defSite)
}
