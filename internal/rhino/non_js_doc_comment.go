// Port of com/google/javascript/rhino/NonJSDocComment.java @ closure-compiler v20261021

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

package rhino

import "github.com/paullepoulpe/jsgo/internal/parser"

// NonJSDocComment is a minimal class holding information about a nonJSDoc
// comment's source location and contents.
type NonJSDocComment struct {
	startPosition     parser.LinePosition
	endPosition       parser.LinePosition
	contents          string
	endsAsLineComment bool
	isInline          bool
}

// Java: NonJSDocComment(SourcePosition, SourcePosition, String)
func NewNonJSDocComment(start parser.LinePosition, end parser.LinePosition, s *string) *NonJSDocComment {
	var contents string
	if s != nil {
		contents = *s
	}
	return &NonJSDocComment{
		startPosition: start,
		endPosition:   end,
		contents:      contents,
	}
}

// Java: NonJSDocComment.getCommentString()
func (c *NonJSDocComment) GetCommentString() string {
	return c.contents
}

// Java: NonJSDocComment.getStartPosition()
func (c *NonJSDocComment) GetStartPosition() parser.LinePosition {
	return c.startPosition
}

// Java: NonJSDocComment.getEndPosition()
func (c *NonJSDocComment) GetEndPosition() parser.LinePosition {
	return c.endPosition
}

// Java: NonJSDocComment.setEndsAsLineComment(boolean)
func (c *NonJSDocComment) SetEndsAsLineComment(endsAsLineComment bool) {
	c.endsAsLineComment = endsAsLineComment
}

// Java: NonJSDocComment.setIsInline(boolean)
func (c *NonJSDocComment) SetIsInline(isInline bool) {
	c.isInline = isInline
}

// Java: NonJSDocComment.isEndingAsLineComment()
func (c *NonJSDocComment) IsEndingAsLineComment() bool {
	return c.endsAsLineComment
}

// Java: NonJSDocComment.isInline()
func (c *NonJSDocComment) IsInline() bool {
	return c.isInline
}

// Java: NonJSDocComment.toString()
func (c *NonJSDocComment) String() string {
	return "NonJSDocComment : " + c.GetCommentString()
}
