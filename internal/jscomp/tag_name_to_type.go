// Port of com/google/javascript/jscomp/TagNameToType.java @ closure-compiler v20261021

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

// TagNameToType contains a mapping from HTML Element tag name to the
// javascript type of that element at runtime.
type TagNameToType struct{}

var tagNameToTypeMap = map[string]string{
	"a":        "HTMLAnchorElement",
	"area":     "HTMLAreaElement",
	"audio":    "HTMLAudioElement",
	"base":     "HTMLBaseElement",
	"body":     "HTMLBodyElement",
	"br":       "HTMLBRElement",
	"button":   "HTMLButtonElement",
	"canvas":   "HTMLCanvasElement",
	"caption":  "HTMLTableCaptionElement",
	"col":      "HTMLTableColElement",
	"content":  "HTMLContentElement",
	"data":     "HTMLDataElement",
	"datalist": "HTMLDataListElement",
	"del":      "HTMLModElement",
	"dialog":   "HTMLDialogElement",
	"dir":      "HTMLDirectoryElement",
	"div":      "HTMLDivElement",
	"dl":       "HTMLDListElement",
	"embed":    "HTMLEmbedElement",
	"fieldset": "HTMLFieldSetElement",
	"font":     "HTMLFontElement",
	"form":     "HTMLFormElement",
	"frame":    "HTMLFrameElement",
	"frameset": "HTMLFrameSetElement",
	"h1":       "HTMLHeadingElement",
	"head":     "HTMLHeadElement",
	"hr":       "HTMLHRElement",
	"html":     "HTMLHtmlElement",
	"iframe":   "HTMLIFrameElement",
	"img":      "HTMLImageElement",
	"input":    "HTMLInputElement",
	"keygen":   "HTMLKeygenElement",
	"label":    "HTMLLabelElement",
	"legend":   "HTMLLegendElement",
	"li":       "HTMLLIElement",
	"link":     "HTMLLinkElement",
	"map":      "HTMLMapElement",
	"marquee":  "HTMLMarqueeElement",
	"menu":     "HTMLMenuElement",
	"menuitem": "HTMLMenuItemElement",
	"meta":     "HTMLMetaElement",
	"meter":    "HTMLMeterElement",
	"object":   "HTMLObjectElement",
	"ol":       "HTMLOListElement",
	"optgroup": "HTMLOptGroupElement",
	"option":   "HTMLOptionElement",
	"output":   "HTMLOutputElement",
	"p":        "HTMLParagraphElement",
	"param":    "HTMLParamElement",
	"pre":      "HTMLPreElement",
	"progress": "HTMLProgressElement",
	"q":        "HTMLQuoteElement",
	"script":   "HTMLScriptElement",
	"select":   "HTMLSelectElement",
	"shadow":   "HTMLShadowElement",
	"source":   "HTMLSourceElement",
	"span":     "HTMLSpanElement",
	"style":    "HTMLStyleElement",
	"table":    "HTMLTableElement",
	"tbody":    "HTMLTableSectionElement",
	"tfoot":    "HTMLTableSectionElement",
	"template": "HTMLTemplateElement",
	"td":       "HTMLTableSectionElement",
	"textarea": "HTMLTextAreaElement",
	"th":       "HTMLTableCellElement",
	"thead":    "HTMLTableCellElement",
	"time":     "HTMLTimeElement",
	"title":    "HTMLTitleElement",
	"tr":       "HTMLTableRowElement",
	"track":    "HTMLTrackElement",
	"ul":       "HTMLUListElement",
	"video":    "HTMLVideoElement",
}

// Java: TagNameToType.getMap()
func TagNameToType_GetMap() map[string]string {
	return tagNameToTypeMap
}
