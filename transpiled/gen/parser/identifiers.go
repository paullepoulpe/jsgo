/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/Identifiers.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.Identifiers
type Identifiers struct {
	idHash_ atomic.Int32
}

func (this_ *Identifiers) BindSelf_(s any) {
	_ = s
}
func (this_ *Identifiers) IsJavaNil_() bool { return this_ == nil }
func (*Identifiers) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.Identifiers", "Identifiers"
}
func (this_ *Identifiers) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Identifiers) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.Identifiers", "com.google.javascript.jscomp.parsing.parser.Identifiers", func(v_ any) bool {
		return j2gort.InstanceOf[*Identifiers](v_)
	})
}

func (this_ *Identifiers) initFields_Identifiers() {
	_ = this_
}

// Java constructor: Identifiers()
func NewIdentifiers() *Identifiers {
	this_ := &Identifiers{}
	InitIdentifiers(this_)
	return this_
}

func InitIdentifiers(this_ *Identifiers) {
	this_.initFields_Identifiers()
}

// Java static method: Identifiers.isIdentifierStart(char ch)
func Identifiers_IsIdentifierStart(ch uint16) bool {
	_ = ch
	if int32(ch) <= int32(127) {
		{
			var ba_9 bool
			{
				var ba_5 bool
				{

					ba_1 := (int32(ch) >= int32(uint16(65)))
					bb_2 := (int32(ch) <= int32(uint16(90)))
					ba_5 = ba_1 && bb_2
				}
				var bb_6 bool
				{

					ba_3 := (int32(ch) >= int32(uint16(97)))
					bb_4 := (int32(ch) <= int32(uint16(122)))
					bb_6 = ba_3 && bb_4
				}
				ba_9 = ba_5 || bb_6
			}
			var bb_10 bool
			{

				ba_7 := (int32(ch) == int32(uint16(95)))
				bb_8 := (int32(ch) == int32(uint16(36)))
				bb_10 = ba_7 || bb_8
			}
			return ba_9 || bb_10
		}
	}
	return j2gort.BoxChar_IsLetter_1_uint16(ch)
}

// Java static method: Identifiers.isIdentifierPart(char ch)
func Identifiers_IsIdentifierPart(ch uint16) bool {
	_ = ch
	if int32(ch) <= int32(127) {
		{
			var ba_23 bool
			{
				var ba_19 bool
				{
					var ba_15 bool
					{

						ba_11 := (int32(ch) >= int32(uint16(65)))
						bb_12 := (int32(ch) <= int32(uint16(90)))
						ba_15 = ba_11 && bb_12
					}
					var bb_16 bool
					{

						ba_13 := (int32(ch) >= int32(uint16(97)))
						bb_14 := (int32(ch) <= int32(uint16(122)))
						bb_16 = ba_13 && bb_14
					}
					ba_19 = ba_15 || bb_16
				}
				var bb_20 bool
				{

					ba_17 := (int32(ch) >= int32(uint16(48)))
					bb_18 := (int32(ch) <= int32(uint16(57)))
					bb_20 = ba_17 && bb_18
				}
				ba_23 = ba_19 || bb_20
			}
			var bb_24 bool
			{

				ba_21 := (int32(ch) == int32(uint16(95)))
				bb_22 := (int32(ch) == int32(uint16(36)))
				bb_24 = ba_21 || bb_22
			}
			return ba_23 || bb_24
		}
	}
	return (Identifiers_IsIdentifierStart(ch) || j2gort.BoxChar_IsDigit_1_uint16(ch))
}
