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
// Port of com/google/javascript/jscomp/disambiguate/ClusterPropagator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.disambiguate.ClusterPropagator
type ClusterPropagator struct {
	idHash_ atomic.Int32
}

func (this_ *ClusterPropagator) BindSelf_(s any) {
	_ = s
}
func (this_ *ClusterPropagator) IsJavaNil_() bool { return this_ == nil }
func (*ClusterPropagator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.ClusterPropagator", "ClusterPropagator"
}
func (this_ *ClusterPropagator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClusterPropagator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.ClusterPropagator", "com.google.javascript.jscomp.disambiguate.ClusterPropagator", func(v_ any) bool {
		return j2gort.InstanceOf[*ClusterPropagator](v_)
	})
}
func (*ClusterPropagator) Implements_graph_FixedPointGraphTraversal_EdgeCallback() {}

func (this_ *ClusterPropagator) initFields_ClusterPropagator() {
	_ = this_
}

// Java constructor: ClusterPropagator()
func NewClusterPropagator() *ClusterPropagator {
	this_ := &ClusterPropagator{}
	InitClusterPropagator(this_)
	return this_
}

func InitClusterPropagator(this_ *ClusterPropagator) {
	this_.initFields_ClusterPropagator()
}

// Java method: ClusterPropagator.traverseEdge(com.google.javascript.jscomp.disambiguate.ColorGraphNode src, java.lang.Object unused, com.google.javascript.jscomp.disambiguate.ColorGraphNode dest)
func (this_ *ClusterPropagator) TraverseEdge(src *ColorGraphNode, unused any, dest *ColorGraphNode) bool {
	_ = src
	_ = unused
	_ = dest
	if src.GetColor() == pkg_colors.StaticGet_StandardColors_TOP_OBJECT() {
		return false
	}
	var startDestPropCount int32 = dest.GetAssociatedProps().Size()
	_ = startDestPropCount
	{
		cur_1 := src.GetAssociatedProps().IterKeys()
		for cur_1.HasNext() {

			prop := cur_1.
				Next()
			_ = prop
			if prop.IsInvalidated() {
				continue
			}
			dest.GetAssociatedProps().PutIfAbsent(prop, StaticGet_ColorGraphNode_PropAssociation_SUPERTYPE())
			prop.GetClusters().Union(src, dest)
		}
	}
	return (startDestPropCount < dest.GetAssociatedProps().Size())
}
