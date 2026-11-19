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
// Port of com/google/javascript/refactoring/NodeMetadata.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.refactoring.NodeMetadata
type NodeMetadata struct {
	compiler AbstractCompiler `gson:"compiler"`
	traversal *NodeTraversal `gson:"traversal"`
	idHash_ atomic.Int32
}

func (this_ *NodeMetadata) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeMetadata) IsJavaNil_() bool { return this_ == nil }
func (*NodeMetadata) JavaClassInfo_() (string, string) { return "com.google.javascript.refactoring.NodeMetadata", "NodeMetadata" }
func (this_ *NodeMetadata) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeMetadata) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.NodeMetadata", "com.google.javascript.refactoring.NodeMetadata", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeMetadata](v_)
	})
}

func (this_ *NodeMetadata) initFields_NodeMetadata() {
	_ = this_
}

// Java constructor: NodeMetadata(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewNodeMetadata_1(compiler AbstractCompiler) *NodeMetadata {
	this_ := &NodeMetadata{}
	InitNodeMetadata_1(this_, compiler)
	return this_
}

func InitNodeMetadata_1(this_ *NodeMetadata, compiler AbstractCompiler) {
	_ = compiler
	InitNodeMetadata_2(this_, compiler, nil)
}

// Java constructor: NodeMetadata(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable NodeTraversal traversal)
func NewNodeMetadata_2(compiler AbstractCompiler, traversal *NodeTraversal) *NodeMetadata {
	this_ := &NodeMetadata{}
	InitNodeMetadata_2(this_, compiler, traversal)
	return this_
}

func InitNodeMetadata_2(this_ *NodeMetadata, compiler AbstractCompiler, traversal *NodeTraversal) {
	_ = compiler
	_ = traversal
	this_.initFields_NodeMetadata()
	this_.compiler = compiler
	this_.traversal = traversal
}

// Java static method: NodeMetadata.fromTraversal(com.google.javascript.jscomp.NodeTraversal traversal)
func NodeMetadata_FromTraversal(traversal *NodeTraversal) *NodeMetadata {
	_ = traversal
	return NewNodeMetadata_2(traversal.GetCompiler(), traversal)
}

// Java method: NodeMetadata.getCompiler()
func (this_ *NodeMetadata) GetCompiler() AbstractCompiler {
	return this_.compiler
}

// Java method: NodeMetadata.getTraversal()
func (this_ *NodeMetadata) GetTraversal() *NodeTraversal {
	return this_.traversal
}

