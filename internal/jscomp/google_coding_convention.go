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
	CodingConvention
}

const (
	googleOptionalArgPrefix = "opt_"
	googleVarArgsName       = "var_args"
)

var googleVarArgsID = rhino.RhinoStringPool_InternID(googleVarArgsName)

func NewGoogleCodingConvention() *GoogleCodingConvention {
	return NewGoogleCodingConventionWithConvention(NewClosureCodingConvention())
}

func NewGoogleCodingConventionWithConvention(convention CodingConvention) *GoogleCodingConvention {
	return &GoogleCodingConvention{CodingConvention: convention}
}

func (c *GoogleCodingConvention) IsExported(name string) bool {
	return CodingConvention_IsExported(c, name)
}

// IsConstant enforces the Google const name convention: the first character
// after the last $ must be an upper-case letter, all subsequent letters must be
// upper case, and the name must be at least 2 characters long.
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
		return isConstantKeyASCII(name, lastDollar+1, len(name))
	}

	utf16Len := jutil.StringLength(name)
	if utf16Len <= 1 {
		return false
	}
	pos := jutil.StringLastIndexOf(name, "$")
	return isConstantKeyUTF16(name, pos+1, utf16Len)
}

func (c *GoogleCodingConvention) IsConstantKey(name string) bool {
	if len(name) == 0 {
		return false
	}
	for i := 0; i < len(name); i++ {
		if name[i] >= 0x80 {
			return isConstantKeyUTF16(name, 0, jutil.StringLength(name))
		}
	}
	return isConstantKeyASCII(name, 0, len(name))
}

func isConstantKeyASCII(name string, start, end int) bool {
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

func isConstantKeyUTF16(name string, start, end int) bool {
	if start >= end || !jutil.CharacterIsUpperCase(jutil.StringCharAt(name, start)) {
		return false
	}
	for i := start + 1; i < end; i++ {
		ch := jutil.StringCharAt(name, i)
		if jutil.CharacterToUpperCase(ch) != ch {
			return false
		}
	}
	return true
}

// IsValidEnumKey enforces [A-Z0-9][A-Z0-9_]*.
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

func (c *GoogleCodingConvention) IsOptionalParameter(parameter rhino.Node) bool {
	return c.CodingConvention.IsOptionalParameter(parameter) ||
		(parameter.IsName() && strings.HasPrefix(parameter.GetString(), googleOptionalArgPrefix))
}

func (c *GoogleCodingConvention) IsVarArgsParameter(parameter rhino.Node) bool {
	return c.CodingConvention.IsVarArgsParameter(parameter) ||
		(parameter.IsName() && parameter.GetStringID() == googleVarArgsID)
}

func (c *GoogleCodingConvention) IsExportedWithLocal(name string, local bool) bool {
	return c.CodingConvention.IsExportedWithLocal(name, local) || (!local && strings.HasPrefix(name, "_"))
}

func (c *GoogleCodingConvention) IsClassFactoryCall(callNode rhino.Node) bool {
	return c.CodingConvention.IsClassFactoryCall(callNode)
}

// GetPackageName returns the package directory path of source, stripping any
// Bazel "-out/.../(bin|genfiles)/" prefix and trailing "/(test|tests|testing)"
// parent directory.
func (c *GoogleCodingConvention) GetPackageName(source rhino.StaticSourceFile) *string {
	name := source.GetName()
	if !strings.ContainsAny(name, "\r\n") {
		if outIdx := strings.Index(name, "-out/"); outIdx >= 0 {
			after := name[outIdx+len("-out/"):]
			binIdx := strings.LastIndex(after, "/bin/")
			genIdx := strings.LastIndex(after, "/genfiles/")
			switch {
			case binIdx >= 0 && binIdx > genIdx:
				name = after[binIdx+len("/bin/"):]
			case genIdx >= 0:
				name = after[genIdx+len("/genfiles/"):]
			}
		}
	}

	lastSlash := strings.LastIndexByte(name, '/')
	if lastSlash < 0 {
		return &emptyPackageName
	}
	pkg := name[:lastSlash]
	if lastSlash+1 < len(name) && !strings.ContainsAny(name, "\r\n") {
		if prevSlash := strings.LastIndexByte(pkg, '/'); prevSlash >= 0 {
			switch pkg[prevSlash+1:] {
			case "test", "tests", "testing":
				pkg = pkg[:prevSlash]
			}
		}
	}
	return &pkg
}
