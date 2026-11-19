/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/NonJSDocComment.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
)

// Java class: com.google.javascript.rhino.NonJSDocComment
type NonJSDocComment struct {
	contents string `gson:"contents"`
	startPosition pkg_parser.SourcePosition `gson:"startPosition"`
	endPosition pkg_parser.SourcePosition `gson:"endPosition"`
	idHash_ atomic.Int32
	endsAsLineComment bool `gson:"endsAsLineComment"`
	isInline bool `gson:"isInline"`
}

func (this_ *NonJSDocComment) BindSelf_(s any) {
	_ = s
}
func (this_ *NonJSDocComment) IsJavaNil_() bool { return this_ == nil }
func (*NonJSDocComment) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.NonJSDocComment", "NonJSDocComment" }
func (this_ *NonJSDocComment) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NonJSDocComment) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.NonJSDocComment", "com.google.javascript.rhino.NonJSDocComment", func(v_ any) bool {
		return j2gort.InstanceOf[*NonJSDocComment](v_)
	})
}

func (this_ *NonJSDocComment) initFields_NonJSDocComment() {
	_ = this_
	this_.contents = j2gort.NullString
}

// Java constructor: NonJSDocComment(com.google.javascript.jscomp.parsing.parser.util.SourcePosition start, com.google.javascript.jscomp.parsing.parser.util.SourcePosition end, java.lang.String s)
func NewNonJSDocComment(start pkg_parser.SourcePosition, end pkg_parser.SourcePosition, s string) *NonJSDocComment {
	this_ := &NonJSDocComment{}
	InitNonJSDocComment(this_, start, end, s)
	return this_
}

func InitNonJSDocComment(this_ *NonJSDocComment, start pkg_parser.SourcePosition, end pkg_parser.SourcePosition, s string) {
	_ = start
	_ = end
	_ = s
	this_.initFields_NonJSDocComment()
	this_.startPosition = start
	this_.endPosition = end
	this_.contents = s
	this_.endsAsLineComment = false
	this_.isInline = false
}

// Java method: NonJSDocComment.getCommentString()
func (this_ *NonJSDocComment) GetCommentString() string {
	if ((this_.contents == j2gort.NullString)) {
		return ""
	}
	return this_.contents
}

// Java method: NonJSDocComment.getStartPosition()
func (this_ *NonJSDocComment) GetStartPosition() pkg_parser.SourcePosition {
	return this_.startPosition
}

// Java method: NonJSDocComment.getEndPosition()
func (this_ *NonJSDocComment) GetEndPosition() pkg_parser.SourcePosition {
	return this_.endPosition
}

// Java method: NonJSDocComment.setEndsAsLineComment(boolean endsAsLineComment)
func (this_ *NonJSDocComment) SetEndsAsLineComment(endsAsLineComment bool) {
	_ = endsAsLineComment
	this_.endsAsLineComment = endsAsLineComment
}

// Java method: NonJSDocComment.setIsInline(boolean isInline)
func (this_ *NonJSDocComment) SetIsInline(isInline bool) {
	_ = isInline
	this_.isInline = isInline
}

// Java method: NonJSDocComment.isEndingAsLineComment()
func (this_ *NonJSDocComment) IsEndingAsLineComment() bool {
	return this_.endsAsLineComment
}

// Java method: NonJSDocComment.isInline()
func (this_ *NonJSDocComment) IsInline() bool {
	return this_.isInline
}

// Java method: NonJSDocComment.toString()
func (this_ *NonJSDocComment) String() string {
	return ("NonJSDocComment : " + j2gort.StrValString(this_.GetCommentString()))
}

