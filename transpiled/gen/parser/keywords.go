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
// Port of com/google/javascript/jscomp/parsing/parser/Keywords.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.parsing.parser.Keywords
type Keywords struct {
	j2gort.EnumBase
	value string
	type_ TokenType
}

func (this_ *Keywords) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *Keywords) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*Keywords) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*Keywords) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *Keywords) initFields_Keywords() {
	_ = this_
	this_.value = j2gort.NullString
}

func initEnumKeywords(this_ *Keywords, enumName_ string, enumOrd_ int32, value string, type_ TokenType) *Keywords {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = value
	_ = type_
	this_.initFields_Keywords()
	this_.value = value
	this_.type_ = type_
	return this_
}

var enum_Keywords_BREAK *Keywords

func StaticGet_Keywords_BREAK() *Keywords {
	clinit_Keywords()
	return enum_Keywords_BREAK
}

var enum_Keywords_CASE *Keywords

func StaticGet_Keywords_CASE() *Keywords {
	clinit_Keywords()
	return enum_Keywords_CASE
}

var enum_Keywords_CATCH *Keywords

func StaticGet_Keywords_CATCH() *Keywords {
	clinit_Keywords()
	return enum_Keywords_CATCH
}

var enum_Keywords_CONTINUE *Keywords

func StaticGet_Keywords_CONTINUE() *Keywords {
	clinit_Keywords()
	return enum_Keywords_CONTINUE
}

var enum_Keywords_DEBUGGER *Keywords

func StaticGet_Keywords_DEBUGGER() *Keywords {
	clinit_Keywords()
	return enum_Keywords_DEBUGGER
}

var enum_Keywords_DEFAULT *Keywords

func StaticGet_Keywords_DEFAULT() *Keywords {
	clinit_Keywords()
	return enum_Keywords_DEFAULT
}

var enum_Keywords_DELETE *Keywords

func StaticGet_Keywords_DELETE() *Keywords {
	clinit_Keywords()
	return enum_Keywords_DELETE
}

var enum_Keywords_DO *Keywords

func StaticGet_Keywords_DO() *Keywords {
	clinit_Keywords()
	return enum_Keywords_DO
}

var enum_Keywords_ELSE *Keywords

func StaticGet_Keywords_ELSE() *Keywords {
	clinit_Keywords()
	return enum_Keywords_ELSE
}

var enum_Keywords_FINALLY *Keywords

func StaticGet_Keywords_FINALLY() *Keywords {
	clinit_Keywords()
	return enum_Keywords_FINALLY
}

var enum_Keywords_FOR *Keywords

func StaticGet_Keywords_FOR() *Keywords {
	clinit_Keywords()
	return enum_Keywords_FOR
}

var enum_Keywords_FUNCTION *Keywords

func StaticGet_Keywords_FUNCTION() *Keywords {
	clinit_Keywords()
	return enum_Keywords_FUNCTION
}

var enum_Keywords_IF *Keywords

func StaticGet_Keywords_IF() *Keywords {
	clinit_Keywords()
	return enum_Keywords_IF
}

var enum_Keywords_IN *Keywords

func StaticGet_Keywords_IN() *Keywords {
	clinit_Keywords()
	return enum_Keywords_IN
}

var enum_Keywords_INSTANCEOF *Keywords

func StaticGet_Keywords_INSTANCEOF() *Keywords {
	clinit_Keywords()
	return enum_Keywords_INSTANCEOF
}

var enum_Keywords_NEW *Keywords

func StaticGet_Keywords_NEW() *Keywords {
	clinit_Keywords()
	return enum_Keywords_NEW
}

var enum_Keywords_RETURN *Keywords

func StaticGet_Keywords_RETURN() *Keywords {
	clinit_Keywords()
	return enum_Keywords_RETURN
}

var enum_Keywords_SWITCH *Keywords

func StaticGet_Keywords_SWITCH() *Keywords {
	clinit_Keywords()
	return enum_Keywords_SWITCH
}

var enum_Keywords_THIS *Keywords

func StaticGet_Keywords_THIS() *Keywords {
	clinit_Keywords()
	return enum_Keywords_THIS
}

var enum_Keywords_THROW *Keywords

func StaticGet_Keywords_THROW() *Keywords {
	clinit_Keywords()
	return enum_Keywords_THROW
}

var enum_Keywords_TRY *Keywords

func StaticGet_Keywords_TRY() *Keywords {
	clinit_Keywords()
	return enum_Keywords_TRY
}

var enum_Keywords_TYPEOF *Keywords

func StaticGet_Keywords_TYPEOF() *Keywords {
	clinit_Keywords()
	return enum_Keywords_TYPEOF
}

var enum_Keywords_VAR *Keywords

func StaticGet_Keywords_VAR() *Keywords {
	clinit_Keywords()
	return enum_Keywords_VAR
}

var enum_Keywords_VOID *Keywords

func StaticGet_Keywords_VOID() *Keywords {
	clinit_Keywords()
	return enum_Keywords_VOID
}

var enum_Keywords_WHILE *Keywords

func StaticGet_Keywords_WHILE() *Keywords {
	clinit_Keywords()
	return enum_Keywords_WHILE
}

var enum_Keywords_WITH *Keywords

func StaticGet_Keywords_WITH() *Keywords {
	clinit_Keywords()
	return enum_Keywords_WITH
}

var enum_Keywords_CLASS *Keywords

func StaticGet_Keywords_CLASS() *Keywords {
	clinit_Keywords()
	return enum_Keywords_CLASS
}

var enum_Keywords_CONST *Keywords

func StaticGet_Keywords_CONST() *Keywords {
	clinit_Keywords()
	return enum_Keywords_CONST
}

var enum_Keywords_ENUM *Keywords

func StaticGet_Keywords_ENUM() *Keywords {
	clinit_Keywords()
	return enum_Keywords_ENUM
}

var enum_Keywords_EXPORT *Keywords

func StaticGet_Keywords_EXPORT() *Keywords {
	clinit_Keywords()
	return enum_Keywords_EXPORT
}

var enum_Keywords_EXTENDS *Keywords

func StaticGet_Keywords_EXTENDS() *Keywords {
	clinit_Keywords()
	return enum_Keywords_EXTENDS
}

var enum_Keywords_IMPORT *Keywords

func StaticGet_Keywords_IMPORT() *Keywords {
	clinit_Keywords()
	return enum_Keywords_IMPORT
}

var enum_Keywords_SUPER *Keywords

func StaticGet_Keywords_SUPER() *Keywords {
	clinit_Keywords()
	return enum_Keywords_SUPER
}

var enum_Keywords_IMPLEMENTS *Keywords

func StaticGet_Keywords_IMPLEMENTS() *Keywords {
	clinit_Keywords()
	return enum_Keywords_IMPLEMENTS
}

var enum_Keywords_INTERFACE *Keywords

func StaticGet_Keywords_INTERFACE() *Keywords {
	clinit_Keywords()
	return enum_Keywords_INTERFACE
}

var enum_Keywords_LET *Keywords

func StaticGet_Keywords_LET() *Keywords {
	clinit_Keywords()
	return enum_Keywords_LET
}

var enum_Keywords_PACKAGE *Keywords

func StaticGet_Keywords_PACKAGE() *Keywords {
	clinit_Keywords()
	return enum_Keywords_PACKAGE
}

var enum_Keywords_PRIVATE *Keywords

func StaticGet_Keywords_PRIVATE() *Keywords {
	clinit_Keywords()
	return enum_Keywords_PRIVATE
}

var enum_Keywords_PROTECTED *Keywords

func StaticGet_Keywords_PROTECTED() *Keywords {
	clinit_Keywords()
	return enum_Keywords_PROTECTED
}

var enum_Keywords_PUBLIC *Keywords

func StaticGet_Keywords_PUBLIC() *Keywords {
	clinit_Keywords()
	return enum_Keywords_PUBLIC
}

var enum_Keywords_STATIC *Keywords

func StaticGet_Keywords_STATIC() *Keywords {
	clinit_Keywords()
	return enum_Keywords_STATIC
}

var enum_Keywords_YIELD *Keywords

func StaticGet_Keywords_YIELD() *Keywords {
	clinit_Keywords()
	return enum_Keywords_YIELD
}

var enum_Keywords_NULL *Keywords

func StaticGet_Keywords_NULL() *Keywords {
	clinit_Keywords()
	return enum_Keywords_NULL
}

var enum_Keywords_TRUE *Keywords

func StaticGet_Keywords_TRUE() *Keywords {
	clinit_Keywords()
	return enum_Keywords_TRUE
}

var enum_Keywords_FALSE *Keywords

func StaticGet_Keywords_FALSE() *Keywords {
	clinit_Keywords()
	return enum_Keywords_FALSE
}

var static_Keywords_KEYWORDS_BY_NAME j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_Keywords_KEYWORDS_BY_NAME() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_Keywords()
	return static_Keywords_KEYWORDS_BY_NAME
}
func StaticSet_Keywords_KEYWORDS_BY_NAME(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_Keywords()
	static_Keywords_KEYWORDS_BY_NAME = v
}

var static_Keywords_KEYWORDS_BY_TYPE j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_Keywords_KEYWORDS_BY_TYPE() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_Keywords()
	return static_Keywords_KEYWORDS_BY_TYPE
}
func StaticSet_Keywords_KEYWORDS_BY_TYPE(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_Keywords()
	static_Keywords_KEYWORDS_BY_TYPE = v
}

var clinitDone_Keywords uint32
var clinitState_Keywords j2gort.ClinitState

func clinit_Keywords() {
	if atomic.LoadUint32(&clinitDone_Keywords) != 0 {
		return
	}
	clinitSlow_Keywords()
}

func clinitSlow_Keywords() {
	if !j2gort.EnterClinit(&clinitState_Keywords, &clinitDone_Keywords) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Keywords, &clinitDone_Keywords)
	enum_Keywords_BREAK = &Keywords{}
	initEnumKeywords(enum_Keywords_BREAK, "BREAK", int32(0), "break", StaticGet_TokenType_BREAK())
	enum_Keywords_CASE = &Keywords{}
	initEnumKeywords(enum_Keywords_CASE, "CASE", int32(1), "case", StaticGet_TokenType_CASE())
	enum_Keywords_CATCH = &Keywords{}
	initEnumKeywords(enum_Keywords_CATCH, "CATCH", int32(2), "catch", StaticGet_TokenType_CATCH())
	enum_Keywords_CONTINUE = &Keywords{}
	initEnumKeywords(enum_Keywords_CONTINUE, "CONTINUE", int32(3), "continue", StaticGet_TokenType_CONTINUE())
	enum_Keywords_DEBUGGER = &Keywords{}
	initEnumKeywords(enum_Keywords_DEBUGGER, "DEBUGGER", int32(4), "debugger", StaticGet_TokenType_DEBUGGER())
	enum_Keywords_DEFAULT = &Keywords{}
	initEnumKeywords(enum_Keywords_DEFAULT, "DEFAULT", int32(5), "default", StaticGet_TokenType_DEFAULT())
	enum_Keywords_DELETE = &Keywords{}
	initEnumKeywords(enum_Keywords_DELETE, "DELETE", int32(6), "delete", StaticGet_TokenType_DELETE())
	enum_Keywords_DO = &Keywords{}
	initEnumKeywords(enum_Keywords_DO, "DO", int32(7), "do", StaticGet_TokenType_DO())
	enum_Keywords_ELSE = &Keywords{}
	initEnumKeywords(enum_Keywords_ELSE, "ELSE", int32(8), "else", StaticGet_TokenType_ELSE())
	enum_Keywords_FINALLY = &Keywords{}
	initEnumKeywords(enum_Keywords_FINALLY, "FINALLY", int32(9), "finally", StaticGet_TokenType_FINALLY())
	enum_Keywords_FOR = &Keywords{}
	initEnumKeywords(enum_Keywords_FOR, "FOR", int32(10), "for", StaticGet_TokenType_FOR())
	enum_Keywords_FUNCTION = &Keywords{}
	initEnumKeywords(enum_Keywords_FUNCTION, "FUNCTION", int32(11), "function", StaticGet_TokenType_FUNCTION())
	enum_Keywords_IF = &Keywords{}
	initEnumKeywords(enum_Keywords_IF, "IF", int32(12), "if", StaticGet_TokenType_IF())
	enum_Keywords_IN = &Keywords{}
	initEnumKeywords(enum_Keywords_IN, "IN", int32(13), "in", StaticGet_TokenType_IN())
	enum_Keywords_INSTANCEOF = &Keywords{}
	initEnumKeywords(enum_Keywords_INSTANCEOF, "INSTANCEOF", int32(14), "instanceof", StaticGet_TokenType_INSTANCEOF())
	enum_Keywords_NEW = &Keywords{}
	initEnumKeywords(enum_Keywords_NEW, "NEW", int32(15), "new", StaticGet_TokenType_NEW())
	enum_Keywords_RETURN = &Keywords{}
	initEnumKeywords(enum_Keywords_RETURN, "RETURN", int32(16), "return", StaticGet_TokenType_RETURN())
	enum_Keywords_SWITCH = &Keywords{}
	initEnumKeywords(enum_Keywords_SWITCH, "SWITCH", int32(17), "switch", StaticGet_TokenType_SWITCH())
	enum_Keywords_THIS = &Keywords{}
	initEnumKeywords(enum_Keywords_THIS, "THIS", int32(18), "this", StaticGet_TokenType_THIS())
	enum_Keywords_THROW = &Keywords{}
	initEnumKeywords(enum_Keywords_THROW, "THROW", int32(19), "throw", StaticGet_TokenType_THROW())
	enum_Keywords_TRY = &Keywords{}
	initEnumKeywords(enum_Keywords_TRY, "TRY", int32(20), "try", StaticGet_TokenType_TRY())
	enum_Keywords_TYPEOF = &Keywords{}
	initEnumKeywords(enum_Keywords_TYPEOF, "TYPEOF", int32(21), "typeof", StaticGet_TokenType_TYPEOF())
	enum_Keywords_VAR = &Keywords{}
	initEnumKeywords(enum_Keywords_VAR, "VAR", int32(22), "var", StaticGet_TokenType_VAR())
	enum_Keywords_VOID = &Keywords{}
	initEnumKeywords(enum_Keywords_VOID, "VOID", int32(23), "void", StaticGet_TokenType_VOID())
	enum_Keywords_WHILE = &Keywords{}
	initEnumKeywords(enum_Keywords_WHILE, "WHILE", int32(24), "while", StaticGet_TokenType_WHILE())
	enum_Keywords_WITH = &Keywords{}
	initEnumKeywords(enum_Keywords_WITH, "WITH", int32(25), "with", StaticGet_TokenType_WITH())
	enum_Keywords_CLASS = &Keywords{}
	initEnumKeywords(enum_Keywords_CLASS, "CLASS", int32(26), "class", StaticGet_TokenType_CLASS())
	enum_Keywords_CONST = &Keywords{}
	initEnumKeywords(enum_Keywords_CONST, "CONST", int32(27), "const", StaticGet_TokenType_CONST())
	enum_Keywords_ENUM = &Keywords{}
	initEnumKeywords(enum_Keywords_ENUM, "ENUM", int32(28), "enum", StaticGet_TokenType_ENUM())
	enum_Keywords_EXPORT = &Keywords{}
	initEnumKeywords(enum_Keywords_EXPORT, "EXPORT", int32(29), "export", StaticGet_TokenType_EXPORT())
	enum_Keywords_EXTENDS = &Keywords{}
	initEnumKeywords(enum_Keywords_EXTENDS, "EXTENDS", int32(30), "extends", StaticGet_TokenType_EXTENDS())
	enum_Keywords_IMPORT = &Keywords{}
	initEnumKeywords(enum_Keywords_IMPORT, "IMPORT", int32(31), "import", StaticGet_TokenType_IMPORT())
	enum_Keywords_SUPER = &Keywords{}
	initEnumKeywords(enum_Keywords_SUPER, "SUPER", int32(32), "super", StaticGet_TokenType_SUPER())
	enum_Keywords_IMPLEMENTS = &Keywords{}
	initEnumKeywords(enum_Keywords_IMPLEMENTS, "IMPLEMENTS", int32(33), "implements", StaticGet_TokenType_IMPLEMENTS())
	enum_Keywords_INTERFACE = &Keywords{}
	initEnumKeywords(enum_Keywords_INTERFACE, "INTERFACE", int32(34), "interface", StaticGet_TokenType_INTERFACE())
	enum_Keywords_LET = &Keywords{}
	initEnumKeywords(enum_Keywords_LET, "LET", int32(35), "let", StaticGet_TokenType_LET())
	enum_Keywords_PACKAGE = &Keywords{}
	initEnumKeywords(enum_Keywords_PACKAGE, "PACKAGE", int32(36), "package", StaticGet_TokenType_PACKAGE())
	enum_Keywords_PRIVATE = &Keywords{}
	initEnumKeywords(enum_Keywords_PRIVATE, "PRIVATE", int32(37), "private", StaticGet_TokenType_PRIVATE())
	enum_Keywords_PROTECTED = &Keywords{}
	initEnumKeywords(enum_Keywords_PROTECTED, "PROTECTED", int32(38), "protected", StaticGet_TokenType_PROTECTED())
	enum_Keywords_PUBLIC = &Keywords{}
	initEnumKeywords(enum_Keywords_PUBLIC, "PUBLIC", int32(39), "public", StaticGet_TokenType_PUBLIC())
	enum_Keywords_STATIC = &Keywords{}
	initEnumKeywords(enum_Keywords_STATIC, "STATIC", int32(40), "static", StaticGet_TokenType_STATIC())
	enum_Keywords_YIELD = &Keywords{}
	initEnumKeywords(enum_Keywords_YIELD, "YIELD", int32(41), "yield", StaticGet_TokenType_YIELD())
	enum_Keywords_NULL = &Keywords{}
	initEnumKeywords(enum_Keywords_NULL, "NULL", int32(42), "null", StaticGet_TokenType_NULL())
	enum_Keywords_TRUE = &Keywords{}
	initEnumKeywords(enum_Keywords_TRUE, "TRUE", int32(43), "true", StaticGet_TokenType_TRUE())
	enum_Keywords_FALSE = &Keywords{}
	initEnumKeywords(enum_Keywords_FALSE, "FALSE", int32(44), "false", StaticGet_TokenType_FALSE())
	enumValues_Keywords = []*Keywords{enum_Keywords_BREAK, enum_Keywords_CASE, enum_Keywords_CATCH, enum_Keywords_CONTINUE, enum_Keywords_DEBUGGER, enum_Keywords_DEFAULT, enum_Keywords_DELETE, enum_Keywords_DO, enum_Keywords_ELSE, enum_Keywords_FINALLY, enum_Keywords_FOR, enum_Keywords_FUNCTION, enum_Keywords_IF, enum_Keywords_IN, enum_Keywords_INSTANCEOF, enum_Keywords_NEW, enum_Keywords_RETURN, enum_Keywords_SWITCH, enum_Keywords_THIS, enum_Keywords_THROW, enum_Keywords_TRY, enum_Keywords_TYPEOF, enum_Keywords_VAR, enum_Keywords_VOID, enum_Keywords_WHILE, enum_Keywords_WITH, enum_Keywords_CLASS, enum_Keywords_CONST, enum_Keywords_ENUM, enum_Keywords_EXPORT, enum_Keywords_EXTENDS, enum_Keywords_IMPORT, enum_Keywords_SUPER, enum_Keywords_IMPLEMENTS, enum_Keywords_INTERFACE, enum_Keywords_LET, enum_Keywords_PACKAGE, enum_Keywords_PRIVATE, enum_Keywords_PROTECTED, enum_Keywords_PUBLIC, enum_Keywords_STATIC, enum_Keywords_YIELD, enum_Keywords_NULL, enum_Keywords_TRUE, enum_Keywords_FALSE}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.Keywords", []any{enum_Keywords_BREAK, enum_Keywords_CASE, enum_Keywords_CATCH, enum_Keywords_CONTINUE, enum_Keywords_DEBUGGER, enum_Keywords_DEFAULT, enum_Keywords_DELETE, enum_Keywords_DO, enum_Keywords_ELSE, enum_Keywords_FINALLY, enum_Keywords_FOR, enum_Keywords_FUNCTION, enum_Keywords_IF, enum_Keywords_IN, enum_Keywords_INSTANCEOF, enum_Keywords_NEW, enum_Keywords_RETURN, enum_Keywords_SWITCH, enum_Keywords_THIS, enum_Keywords_THROW, enum_Keywords_TRY, enum_Keywords_TYPEOF, enum_Keywords_VAR, enum_Keywords_VOID, enum_Keywords_WHILE, enum_Keywords_WITH, enum_Keywords_CLASS, enum_Keywords_CONST, enum_Keywords_ENUM, enum_Keywords_EXPORT, enum_Keywords_EXTENDS, enum_Keywords_IMPORT, enum_Keywords_SUPER, enum_Keywords_IMPLEMENTS, enum_Keywords_INTERFACE, enum_Keywords_LET, enum_Keywords_PACKAGE, enum_Keywords_PRIVATE, enum_Keywords_PROTECTED, enum_Keywords_PUBLIC, enum_Keywords_STATIC, enum_Keywords_YIELD, enum_Keywords_NULL, enum_Keywords_TRUE, enum_Keywords_FALSE})
	var keywordsByName j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = keywordsByName
	var keywordsByType *j2gort.Ext_java_util_EnumMap = j2gort.NewExt_java_util_EnumMap_1_Class(j2gort.ClassLiteral("com.google.javascript.jscomp.parsing.parser.TokenType", "TokenType"))
	_ = keywordsByType
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_java_util_EnumSet_AllOf(j2gort.ClassLiteral("com.google.javascript.jscomp.parsing.parser.Keywords", "Keywords"))))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			kw := j2gort.CastPtr[Keywords](rawElem_2)
			_ = kw
			keywordsByName.Put_2(j2gort.BoxString(kw.value), j2gort.PtrToAny(kw))
			keywordsByType.Put_2(j2gort.BoxNullableVal(kw.type_), j2gort.PtrToAny(kw))
		}
	}
	StaticSet_Keywords_KEYWORDS_BY_NAME(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](keywordsByName.BuildOrThrow()))
	StaticSet_Keywords_KEYWORDS_BY_TYPE(j2gort.Ext_com_google_common_collect_Maps_ImmutableEnumMap(keywordsByType))
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.parsing.parser.Keywords", "com.google.javascript.jscomp.parsing.parser.Keywords", func() {
		clinit_Keywords()
	})
}

var enumValues_Keywords []*Keywords

func Keywords_Values() []*Keywords {
	clinit_Keywords()
	return j2gort.CloneSlice(enumValues_Keywords)
}

func Keywords_ValueOf(name string) *Keywords {
	clinit_Keywords()
	for _, c := range enumValues_Keywords {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.Keywords." + name))
}

// Java method: Keywords.toString()
func (this_ *Keywords) String() string {
	return this_.value
}

// Java static method: Keywords.isKeyword(java.lang.String value)
func Keywords_IsKeyword_1_string(value string) bool {
	clinit_Keywords()
	_ = value
	return (Keywords_Get_1_string(value) != nil)
}

// Java static method: Keywords.isKeyword(com.google.javascript.jscomp.parsing.parser.TokenType token)
func Keywords_IsKeyword_1_TokenType(token TokenType) bool {
	clinit_Keywords()
	_ = token
	return (Keywords_Get_1_TokenType(token) != nil)
}

// Java static method: Keywords.isStrictKeyword(com.google.javascript.jscomp.parsing.parser.TokenType token)
func Keywords_IsStrictKeyword(token TokenType) bool {
	clinit_Keywords()
	_ = token
	switch (token).Ordinal() {
	case int32(36) /* TokenType.IMPLEMENTS */, int32(37) /* TokenType.INTERFACE */, int32(38) /* TokenType.LET */, int32(39) /* TokenType.PACKAGE */, int32(40) /* TokenType.PRIVATE */, int32(41) /* TokenType.PROTECTED */, int32(42) /* TokenType.PUBLIC */, int32(43) /* TokenType.STATIC */, int32(44) /* TokenType.YIELD */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: Keywords.getTokenType(java.lang.String value)
func Keywords_GetTokenType(value string) TokenType {
	clinit_Keywords()
	_ = value
	return j2gort.CastPtr[Keywords](j2gort.MapGetStr(StaticGet_Keywords_KEYWORDS_BY_NAME(), value)).type_
}

// Java static method: Keywords.get(java.lang.String value)
func Keywords_Get_1_string(value string) *Keywords {
	clinit_Keywords()
	_ = value
	return j2gort.CastPtr[Keywords](j2gort.MapGetStr(StaticGet_Keywords_KEYWORDS_BY_NAME(), value))
}

// Java static method: Keywords.get(com.google.javascript.jscomp.parsing.parser.TokenType token)
func Keywords_Get_1_TokenType(token TokenType) *Keywords {
	clinit_Keywords()
	_ = token
	return j2gort.CastPtr[Keywords](StaticGet_Keywords_KEYWORDS_BY_TYPE().Get(j2gort.BoxNullableVal(token)))
}

func (this_ *Keywords) IsJavaNil_() bool { return this_ == nil }
func (*Keywords) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.Keywords", "Keywords"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.Keywords", "com.google.javascript.jscomp.parsing.parser.Keywords", func(v_ any) bool {
		return j2gort.InstanceOf[*Keywords](v_)
	})
}
