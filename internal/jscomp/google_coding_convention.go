// Port of com/google/javascript/jscomp/GoogleCodingConvention.java @ closure-compiler v20261021

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// GoogleCodingConvention describes the Google-specific JavaScript coding
// conventions. Within Google, variable names are semantically significant.
type GoogleCodingConvention struct {
	CodingConventions_ProxyBase
}

const googleCodingConvention_serialVersionUID int64 = 1

const GoogleCodingConvention_OPTIONAL_ARG_PREFIX = "opt_"

const GoogleCodingConvention_VAR_ARGS_NAME = "var_args"

var GoogleCodingConvention_ENUM_KEY_PATTERN = jutil.Pattern_Compile("[A-Z0-9][A-Z0-9_]*")

var GoogleCodingConvention_PACKAGE_WITH_TEST_DIR = jutil.Pattern_Compile("^(.*)/(?:test|tests|testing)/(?:[^/]+)$")

var GoogleCodingConvention_GENFILES_DIR = jutil.Pattern_Compile("-out/.*/(bin|genfiles)/(.*)$")

// NewGoogleCodingConvention: by default, decorate the
// ClosureCodingConvention.
//
// Java: GoogleCodingConvention()
func NewGoogleCodingConvention() *GoogleCodingConvention {
	return NewGoogleCodingConventionWithConvention(NewClosureCodingConvention())
}

// NewGoogleCodingConventionWithConvention decorates a wrapped
// CodingConvention.
//
// Java: GoogleCodingConvention(CodingConvention)
func NewGoogleCodingConventionWithConvention(convention CodingConvention) *GoogleCodingConvention {
	c := &GoogleCodingConvention{}
	c.self = c
	c.InitCodingConventions_ProxyBase(convention)
	return c
}

// IsConstant enforces the Google const name convention, that the first
// character after the last $ must be an upper-case letter and all
// subsequent letters must be upper case. The name must be at least 2
// characters long.
//
// Examples:
//
//	aaa          Not constant - lower-case letters in the name
//	A            Not constant - too short
//	goog$A       Constant - letters after the $ are upper-case.
//	AA17         Constant - digits can appear after the first letter
//	goog$7A      Not constant - first character after the $ must be
//	             upper case.
//	$A           Constant - doesn't have to be anything in front of the $
//
// Java: GoogleCodingConvention.isConstant(String)
func (c *GoogleCodingConvention) IsConstant(name string) bool {
	if len(name) <= 1 {
		return false
	}

	lastDollar := -1
	ascii := true
	for i := 0; i < len(name); i++ {
		b := name[i]
		if b >= 0x80 {
			ascii = false
			break
		}
		if b == '$' {
			lastDollar = i
		}
	}
	if ascii {
		start := 0
		if lastDollar >= 0 {
			start = lastDollar + 1
		}
		return googleCodingConvention_isConstantKeyASCII(name, start, len(name))
	}

	if jutil.StringLength(name) <= 1 {
		return false
	}

	// In compiled code, '$' is often a namespace delimiter. To allow inlining
	// of namespaced constants, we strip off any namespaces here.
	pos := jutil.StringLastIndexOf(name, "$")
	start := 0
	if pos >= 0 {
		start = pos + 1
	}

	return GoogleCodingConvention_IsConstantKeyWithRange(name, start, jutil.StringLength(name))
}

// Java: GoogleCodingConvention.isConstantKey(String)
func (c *GoogleCodingConvention) IsConstantKey(name string) bool {
	if len(name) == 0 {
		return false
	}
	ascii := true
	for i := 0; i < len(name); i++ {
		if name[i] >= 0x80 {
			ascii = false
			break
		}
	}
	if ascii {
		return googleCodingConvention_isConstantKeyASCII(name, 0, len(name))
	}
	return GoogleCodingConvention_IsConstantKeyWithRange(name, 0, jutil.StringLength(name))
}

func googleCodingConvention_isConstantKeyASCII(name string, start int, end int) bool {
	if start >= end {
		return false
	}
	first := name[start]
	if first < 'A' || first > 'Z' {
		return false
	}
	for i := start + 1; i < end; i++ {
		ch := name[i]
		if ch >= 'a' && ch <= 'z' {
			return false
		}
	}
	return true
}

// GoogleCodingConvention_IsConstantKeyWithRange returns true iff the
// substring of name from [start, end) is a constant key.
//
// This method takes start and end indices, rather than checking an entire
// string, to avoid the need for unnecessary .substring(start, end) calls
//
// Java: GoogleCodingConvention.isConstantKey(String, int, int)
func GoogleCodingConvention_IsConstantKeyWithRange(name string, start int, end int) bool {
	if start >= end || !jutil.CharacterIsUpperCase(jutil.StringCharAt(name, start)) {
		return false
	}
	for i := start + 1; i < end; i++ {
		// check this instead of Character.isUpperCase - they have different results for some
		// characters like symbols. symbols are allowed unless at the start index.
		ch := jutil.StringCharAt(name, i)
		if jutil.CharacterToUpperCase(ch) != ch {
			return false
		}
	}
	return true
}

// IsValidEnumKey enforces Google's convention about enum key names. They
// must match the regular expression [A-Z0-9][A-Z0-9_]*.
//
// Examples:
//   - A
//   - 213
//   - FOO_BAR
//
// Java: GoogleCodingConvention.isValidEnumKey(String)
func (c *GoogleCodingConvention) IsValidEnumKey(key string) bool {
	if len(key) == 0 {
		return false
	}
	first := key[0]
	if !((first >= 'A' && first <= 'Z') || (first >= '0' && first <= '9')) {
		return false
	}
	for i := 1; i < len(key); i++ {
		ch := key[i]
		if !((ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return false
		}
	}
	return true
}

// IsOptionalParameter: in Google code, parameter names beginning with opt_
// are treated as optional arguments.
//
// Java: GoogleCodingConvention.isOptionalParameter(Node)
func (c *GoogleCodingConvention) IsOptionalParameter(parameter *rhino.Node) bool {
	return c.CodingConventions_ProxyBase.IsOptionalParameter(parameter) ||
		(parameter.IsName() && strings.HasPrefix(parameter.GetString(), GoogleCodingConvention_OPTIONAL_ARG_PREFIX))
}

// Java: GoogleCodingConvention.isVarArgsParameter(Node)
func (c *GoogleCodingConvention) IsVarArgsParameter(parameter *rhino.Node) bool {
	return c.CodingConventions_ProxyBase.IsVarArgsParameter(parameter) ||
		(parameter.IsName() && GoogleCodingConvention_VAR_ARGS_NAME == parameter.GetString())
}

// IsExportedWithLocal: in Google code, any global name starting with an
// underscore is considered exported.
//
// Java: GoogleCodingConvention.isExported(String, boolean)
func (c *GoogleCodingConvention) IsExportedWithLocal(name string, local bool) bool {
	return c.CodingConventions_ProxyBase.IsExportedWithLocal(name, local) || (!local && strings.HasPrefix(name, "_"))
}

// Java: GoogleCodingConvention.isClassFactoryCall(Node)
func (c *GoogleCodingConvention) IsClassFactoryCall(callNode *rhino.Node) bool {
	callTarget := callNode.GetFirstChild()
	// Go-only: performance — in Java, `callTarget.isName() && callTarget.toString().equals("Polymer")`
	// always evaluates to false because `Node.toString()` starts with `"NAME "`. Calling `callTarget.String()`
	// here formats the entire Node and its JSType into a new string on every CALL node during CheckJSDoc.
	return c.CodingConventions_ProxyBase.IsClassFactoryCall(callNode) ||
		(callTarget.IsName() && false)
}

// GetPackageName: in Google code, the package name of a source file is its
// file path. Exceptions: if a source file's parent directory is "test",
// "tests", or "testing", that directory is stripped from the package name.
// If a file is generated, strip the "genfiles" prefix to try to match the
// package of the generating file.
//
// Java: GoogleCodingConvention.getPackageName(StaticSourceFile)
func (c *GoogleCodingConvention) GetPackageName(source rhino.StaticSourceFile) *string {
	name := source.GetName()
	genfilesMatcher := GoogleCodingConvention_GENFILES_DIR.Matcher(name)
	if genfilesMatcher.Find() {
		name = *genfilesMatcher.Group(2)
	}

	m := GoogleCodingConvention_PACKAGE_WITH_TEST_DIR.Matcher(name)
	if m.Find() {
		return m.Group(1)
	} else {
		lastSlash := jutil.StringLastIndexOf(name, "/")
		var result string
		if lastSlash == -1 {
			result = ""
		} else {
			result = jutil.StringSubstring(name, 0, lastSlash)
		}
		return &result
	}
}
