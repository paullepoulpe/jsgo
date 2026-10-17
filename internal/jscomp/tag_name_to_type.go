// Port of com/google/javascript/jscomp/TagNameToType.java @ closure-compiler v20261015

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TagNameToType contains a mapping from HTML Element tag name to the
// javascript type of that element at runtime.
type TagNameToType struct{}

// Java: TagNameToType.getMap()
func TagNameToType_GetMap() *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	m.Put("a", "HTMLAnchorElement")
	m.Put("area", "HTMLAreaElement")
	m.Put("audio", "HTMLAudioElement")
	m.Put("base", "HTMLBaseElement")
	m.Put("body", "HTMLBodyElement")
	m.Put("br", "HTMLBRElement")
	m.Put("button", "HTMLButtonElement")
	m.Put("canvas", "HTMLCanvasElement")
	m.Put("caption", "HTMLTableCaptionElement")
	m.Put("col", "HTMLTableColElement")
	m.Put("content", "HTMLContentElement")
	m.Put("data", "HTMLDataElement")
	m.Put("datalist", "HTMLDataListElement")
	m.Put("del", "HTMLModElement")
	m.Put("dialog", "HTMLDialogElement")
	m.Put("dir", "HTMLDirectoryElement")
	m.Put("div", "HTMLDivElement")
	m.Put("dl", "HTMLDListElement")
	m.Put("embed", "HTMLEmbedElement")
	m.Put("fieldset", "HTMLFieldSetElement")
	m.Put("font", "HTMLFontElement")
	m.Put("form", "HTMLFormElement")
	m.Put("frame", "HTMLFrameElement")
	m.Put("frameset", "HTMLFrameSetElement")
	m.Put("h1", "HTMLHeadingElement")
	m.Put("head", "HTMLHeadElement")
	m.Put("hr", "HTMLHRElement")
	m.Put("html", "HTMLHtmlElement")
	m.Put("iframe", "HTMLIFrameElement")
	m.Put("img", "HTMLImageElement")
	m.Put("input", "HTMLInputElement")
	m.Put("keygen", "HTMLKeygenElement")
	m.Put("label", "HTMLLabelElement")
	m.Put("legend", "HTMLLegendElement")
	m.Put("li", "HTMLLIElement")
	m.Put("link", "HTMLLinkElement")
	m.Put("map", "HTMLMapElement")
	m.Put("marquee", "HTMLMarqueeElement")
	m.Put("menu", "HTMLMenuElement")
	m.Put("menuitem", "HTMLMenuItemElement")
	m.Put("meta", "HTMLMetaElement")
	m.Put("meter", "HTMLMeterElement")
	m.Put("object", "HTMLObjectElement")
	m.Put("ol", "HTMLOListElement")
	m.Put("optgroup", "HTMLOptGroupElement")
	m.Put("option", "HTMLOptionElement")
	m.Put("output", "HTMLOutputElement")
	m.Put("p", "HTMLParagraphElement")
	m.Put("param", "HTMLParamElement")
	m.Put("pre", "HTMLPreElement")
	m.Put("progress", "HTMLProgressElement")
	m.Put("q", "HTMLQuoteElement")
	m.Put("script", "HTMLScriptElement")
	m.Put("select", "HTMLSelectElement")
	m.Put("shadow", "HTMLShadowElement")
	m.Put("source", "HTMLSourceElement")
	m.Put("span", "HTMLSpanElement")
	m.Put("style", "HTMLStyleElement")
	m.Put("table", "HTMLTableElement")
	m.Put("tbody", "HTMLTableSectionElement")
	m.Put("tfoot", "HTMLTableSectionElement")
	m.Put("template", "HTMLTemplateElement")
	m.Put("td", "HTMLTableSectionElement")
	m.Put("textarea", "HTMLTextAreaElement")
	m.Put("th", "HTMLTableCellElement")
	m.Put("thead", "HTMLTableCellElement")
	m.Put("time", "HTMLTimeElement")
	m.Put("title", "HTMLTitleElement")
	m.Put("tr", "HTMLTableRowElement")
	m.Put("track", "HTMLTrackElement")
	m.Put("ul", "HTMLUListElement")
	m.Put("video", "HTMLVideoElement")
	return m
}
