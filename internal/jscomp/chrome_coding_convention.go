// Port of com/google/javascript/jscomp/ChromeCodingConvention.java @ closure-compiler v20261015

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

// Java: ChromeCodingConvention()
func NewChromeCodingConvention() *ChromeCodingConvention {
	return NewChromeCodingConventionWithConvention(CodingConventions_GetDefault())
}

// Java: ChromeCodingConvention(CodingConvention)
func NewChromeCodingConventionWithConvention(wrapped CodingConvention) *ChromeCodingConvention {
	c := &ChromeCodingConvention{}
	c.self = c
	c.InitCodingConventions_ProxyBase(wrapped)
	return c
}

var ChromeCodingConvention_CR_ADD_SINGLETON_GETTER = rhino.QualifiedName_Of("cr.addSingletonGetter")

// Java: ChromeCodingConvention.getSingletonGetterClassName(Node)
func (c *ChromeCodingConvention) GetSingletonGetterClassName(callNode *rhino.Node) *string {
	callArg := callNode.GetFirstChild()
	if !ChromeCodingConvention_CR_ADD_SINGLETON_GETTER.Matches(callArg) || !callNode.HasTwoChildren() {
		return c.CodingConventions_ProxyBase.GetSingletonGetterClassName(callNode)
	}
	return callArg.GetNext().GetQualifiedName()
}

// Java: ChromeCodingConvention.applySingletonGetter(NominalTypeBuilder, FunctionType)
func (c *ChromeCodingConvention) ApplySingletonGetter(classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType) {
	defSite := classType.Constructor().GetSource()
	classType.DeclareConstructorProperty("getInstance", getterType, defSite)
	classType.DeclareConstructorProperty("instance_", classType.Instance(), defSite)
}
