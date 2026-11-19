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
// Port of com/google/javascript/refactoring/Match.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.refactoring.Match
type Match struct {
	metadata *NodeMetadata `gson:"metadata"`
	node pkg_rhino.Node `gson:"node"`
	idHash_ atomic.Int32
}

func (this_ *Match) BindSelf_(s any) {
	_ = s
}
func (this_ *Match) IsJavaNil_() bool { return this_ == nil }
func (*Match) JavaClassInfo_() (string, string) { return "com.google.javascript.refactoring.Match", "Match" }
func (this_ *Match) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Match) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.Match", "com.google.javascript.refactoring.Match", func(v_ any) bool {
		return j2gort.InstanceOf[*Match](v_)
	})
}

func (this_ *Match) initFields_Match() {
	_ = this_
}

// Java constructor: Match(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func NewMatch(node pkg_rhino.Node, metadata *NodeMetadata) *Match {
	this_ := &Match{}
	InitMatch(this_, node, metadata)
	return this_
}

func InitMatch(this_ *Match, node pkg_rhino.Node, metadata *NodeMetadata) {
	_ = node
	_ = metadata
	this_.initFields_Match()
	this_.node = node
	this_.metadata = metadata
}

// Java method: Match.getNode()
func (this_ *Match) GetNode() pkg_rhino.Node {
	return this_.node
}

// Java method: Match.getMetadata()
func (this_ *Match) GetMetadata() *NodeMetadata {
	return this_.metadata
}

