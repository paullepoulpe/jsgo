/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/TagNameToType.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.TagNameToType
type TagNameToType struct {
	idHash_ atomic.Int32
}

func (this_ *TagNameToType) BindSelf_(s any) {
	_ = s
}
func (this_ *TagNameToType) IsJavaNil_() bool { return this_ == nil }
func (*TagNameToType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TagNameToType", "TagNameToType"
}
func (this_ *TagNameToType) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TagNameToType) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TagNameToType", "com.google.javascript.jscomp.TagNameToType", func(v_ any) bool {
		return j2gort.InstanceOf[*TagNameToType](v_)
	})
}

func (this_ *TagNameToType) initFields_TagNameToType() {
	_ = this_
}

// Java constructor: TagNameToType()
func NewTagNameToType() *TagNameToType {
	this_ := &TagNameToType{}
	InitTagNameToType(this_)
	return this_
}

func InitTagNameToType(this_ *TagNameToType) {
	this_.initFields_TagNameToType()
}

// Java static method: TagNameToType.getMap()
func TagNameToType_GetMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.NewExt_com_google_common_collect_ImmutableMap_Builder().Put_2("a", "HTMLAnchorElement")).Put_2("area", "HTMLAreaElement")).Put_2("audio", "HTMLAudioElement")).Put_2("base", "HTMLBaseElement")).Put_2("body", "HTMLBodyElement")).Put_2("br", "HTMLBRElement")).Put_2("button", "HTMLButtonElement")).Put_2("canvas", "HTMLCanvasElement")).Put_2("caption", "HTMLTableCaptionElement")).Put_2("col", "HTMLTableColElement")).Put_2("content", "HTMLContentElement")).Put_2("data", "HTMLDataElement")).Put_2("datalist", "HTMLDataListElement")).Put_2("del", "HTMLModElement")).Put_2("dialog", "HTMLDialogElement")).Put_2("dir", "HTMLDirectoryElement")).Put_2("div", "HTMLDivElement")).Put_2("dl", "HTMLDListElement")).Put_2("embed", "HTMLEmbedElement")).Put_2("fieldset", "HTMLFieldSetElement")).Put_2("font", "HTMLFontElement")).Put_2("form", "HTMLFormElement")).Put_2("frame", "HTMLFrameElement")).Put_2("frameset", "HTMLFrameSetElement")).Put_2("h1", "HTMLHeadingElement")).Put_2("head", "HTMLHeadElement")).Put_2("hr", "HTMLHRElement")).Put_2("html", "HTMLHtmlElement")).Put_2("iframe", "HTMLIFrameElement")).Put_2("img", "HTMLImageElement")).Put_2("input", "HTMLInputElement")).Put_2("keygen", "HTMLKeygenElement")).Put_2("label", "HTMLLabelElement")).Put_2("legend", "HTMLLegendElement")).Put_2("li", "HTMLLIElement")).Put_2("link", "HTMLLinkElement")).Put_2("map", "HTMLMapElement")).Put_2("marquee", "HTMLMarqueeElement")).Put_2("menu", "HTMLMenuElement")).Put_2("menuitem", "HTMLMenuItemElement")).Put_2("meta", "HTMLMetaElement")).Put_2("meter", "HTMLMeterElement")).Put_2("object", "HTMLObjectElement")).Put_2("ol", "HTMLOListElement")).Put_2("optgroup", "HTMLOptGroupElement")).Put_2("option", "HTMLOptionElement")).Put_2("output", "HTMLOutputElement")).Put_2("p", "HTMLParagraphElement")).Put_2("param", "HTMLParamElement")).Put_2("pre", "HTMLPreElement")).Put_2("progress", "HTMLProgressElement")).Put_2("q", "HTMLQuoteElement")).Put_2("script", "HTMLScriptElement")).Put_2("select", "HTMLSelectElement")).Put_2("shadow", "HTMLShadowElement")).Put_2("source", "HTMLSourceElement")).Put_2("span", "HTMLSpanElement")).Put_2("style", "HTMLStyleElement")).Put_2("table", "HTMLTableElement")).Put_2("tbody", "HTMLTableSectionElement")).Put_2("tfoot", "HTMLTableSectionElement")).Put_2("template", "HTMLTemplateElement")).Put_2("td", "HTMLTableSectionElement")).Put_2("textarea", "HTMLTextAreaElement")).Put_2("th", "HTMLTableCellElement")).Put_2("thead", "HTMLTableCellElement")).Put_2("time", "HTMLTimeElement")).Put_2("title", "HTMLTitleElement")).Put_2("tr", "HTMLTableRowElement")).Put_2("track", "HTMLTrackElement")).Put_2("ul", "HTMLUListElement")).Put_2("video", "HTMLVideoElement")).BuildOrThrow())
}
