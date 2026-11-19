/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/disambiguate/PropertyClustering.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.disambiguate.PropertyClustering
type PropertyClustering struct {
	name                   string                                             `gson:"name"`
	useSites               j2gort.OrderedMap[pkg_rhino.Node, *ColorGraphNode] `gson:"useSites"`
	clusters               *pkg_graph.StandardUnionFind[*ColorGraphNode]      `gson:"clusters"`
	originalNameClusterRep *ColorGraphNode                                    `gson:"originalNameClusterRep"`
	lastInvalidation       Invalidation                                       `gson:"lastInvalidation"`
	idHash_                atomic.Int32
}

func (this_ *PropertyClustering) BindSelf_(s any) {
	_ = s
}
func (this_ *PropertyClustering) IsJavaNil_() bool { return this_ == nil }
func (*PropertyClustering) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.PropertyClustering", "PropertyClustering"
}
func (this_ *PropertyClustering) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PropertyClustering) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.PropertyClustering", "com.google.javascript.jscomp.disambiguate.PropertyClustering", func(v_ any) bool {
		return j2gort.InstanceOf[*PropertyClustering](v_)
	})
}

func (this_ *PropertyClustering) initFields_PropertyClustering() {
	_ = this_
	this_.name = j2gort.NullString
	this_.useSites = j2gort.NewOrderedMap[pkg_rhino.Node, *ColorGraphNode]()
	this_.clusters = pkg_graph.NewStandardUnionFind[*ColorGraphNode]()
}

// Java constructor: PropertyClustering(java.lang.String name)
func NewPropertyClustering(name string) *PropertyClustering {
	this_ := &PropertyClustering{}
	InitPropertyClustering(this_, name)
	return this_
}

func InitPropertyClustering(this_ *PropertyClustering, name string) {
	_ = name
	this_.initFields_PropertyClustering()
	this_.name = j2gort.CheckNotNullString(name)
}

// Java method: PropertyClustering.getName()
func (this_ *PropertyClustering) GetName() string {
	return this_.name
}

// Java method: PropertyClustering.putUseSite(com.google.javascript.rhino.Node site, com.google.javascript.jscomp.disambiguate.ColorGraphNode owner)
func (this_ *PropertyClustering) PutUseSite(site pkg_rhino.Node, owner *ColorGraphNode) {
	_ = site
	_ = owner
	this_.useSites.Put(site, owner)
}

// Java method: PropertyClustering.forEachUseSite(java.util.function.BiConsumer<com.google.javascript.rhino.Node,com.google.javascript.jscomp.disambiguate.ColorGraphNode> consumer)
func (this_ *PropertyClustering) ForEachUseSite(consumer j2gort.Ext_java_util_function_BiConsumer) {
	_ = consumer
	this_.useSites.ForEachBiConsumer(consumer)
}

// Java method: PropertyClustering.getUseSites()
func (this_ *PropertyClustering) GetUseSites() *j2gort.OrderedMap[pkg_rhino.Node, *ColorGraphNode] {
	return &this_.useSites

}

// Java method: PropertyClustering.getClusters()
func (this_ *PropertyClustering) GetClusters() *pkg_graph.StandardUnionFind[*ColorGraphNode] {
	return this_.clusters
}

// Java method: PropertyClustering.getOriginalNameClusterRep()
func (this_ *PropertyClustering) GetOriginalNameClusterRep() *ColorGraphNode {
	if !(!(this_.IsInvalidated())) {
		j2gort.CheckState(false)
	}
	if this_.originalNameClusterRep == nil {
		return nil
	} else {
		return this_.clusters.Find(this_.originalNameClusterRep)
	}
}

// Java method: PropertyClustering.isInvalidated()
func (this_ *PropertyClustering) IsInvalidated() bool {
	return this_.lastInvalidation != nil && !this_.lastInvalidation.IsJavaNil_()
}

// Java method: PropertyClustering.invalidate(com.google.javascript.jscomp.disambiguate.Invalidation invalidation)
func (this_ *PropertyClustering) Invalidate(invalidation Invalidation) {
	_ = invalidation
	this_.lastInvalidation = j2gort.CheckNotNullIface(invalidation)
}

// Java method: PropertyClustering.getLastInvalidation()
func (this_ *PropertyClustering) GetLastInvalidation() Invalidation {
	if !(this_.IsInvalidated()) {
		j2gort.CheckState(false)
	}
	return this_.lastInvalidation
}

// Java method: PropertyClustering.registerOriginalNameType(com.google.javascript.jscomp.disambiguate.ColorGraphNode type)
func (this_ *PropertyClustering) RegisterOriginalNameType(type_ *ColorGraphNode) {
	_ = type_
	if !(!(this_.IsInvalidated())) {
		j2gort.CheckState(false)
	}
	if this_.originalNameClusterRep == nil {
		this_.originalNameClusterRep = type_
	}
	this_.clusters.Union(this_.originalNameClusterRep, type_)
}

// Java method: PropertyClustering.toString()
func (this_ *PropertyClustering) String() string {
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).Add_2_string_any("name", j2gort.BoxString(this_.name)).String()
}
