// Port of com/google/javascript/jscomp/CssRenamingMap.java @ closure-compiler v20261015

/*
 * Copyright 2009 The Closure Compiler Authors.
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

// CssRenamingMap is the interface used by ReplaceCssNames to substitute
// CSS class names.
type CssRenamingMap interface {
	RenamingMap

	// Java: CssRenamingMap.getStyle()
	GetStyle() CssRenamingMap_Style
}

// CssRenamingMap_Style is the kind of renaming map.
type CssRenamingMap_Style int32

const (
	CssRenamingMap_Style_BY_WHOLE CssRenamingMap_Style = iota
	CssRenamingMap_Style_BY_PART
)

var cssRenamingMap_Style_names = [...]string{"BY_WHOLE", "BY_PART"}

// CssRenamingMap_Style_values is Style.values().
func CssRenamingMap_Style_values() []CssRenamingMap_Style {
	return []CssRenamingMap_Style{CssRenamingMap_Style_BY_WHOLE, CssRenamingMap_Style_BY_PART}
}

// Name is Enum.name().
func (s CssRenamingMap_Style) Name() string { return cssRenamingMap_Style_names[s] }

// String is Enum.toString().
func (s CssRenamingMap_Style) String() string { return s.Name() }

// Ordinal is Enum.ordinal().
func (s CssRenamingMap_Style) Ordinal() int { return int(s) }

// CssRenamingMap_Style_ValueOf is Style.valueOf(String).
func CssRenamingMap_Style_ValueOf(name string) CssRenamingMap_Style {
	for i, n := range cssRenamingMap_Style_names {
		if n == name {
			return CssRenamingMap_Style(i)
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(
		"No enum constant com.google.javascript.jscomp.CssRenamingMap.Style." + name))
}

// CssRenamingMap_ByPart is the ByPart renaming map: an abstract class whose
// subclasses embed CssRenamingMap_ByPartBase and implement Get.
type CssRenamingMap_ByPart interface {
	CssRenamingMap
	GetCssRenamingMap_ByPartBase() *CssRenamingMap_ByPartBase
}

// CssRenamingMap_ByPartBase holds ByPart's base implementation.
type CssRenamingMap_ByPartBase struct{}

// GetCssRenamingMap_ByPartBase gives access to the base through the
// interface (and marks ByPart subclasses for `instanceof`).
func (b *CssRenamingMap_ByPartBase) GetCssRenamingMap_ByPartBase() *CssRenamingMap_ByPartBase {
	return b
}

// Java: CssRenamingMap.ByPart.getStyle()
func (b *CssRenamingMap_ByPartBase) GetStyle() CssRenamingMap_Style {
	return CssRenamingMap_Style_BY_PART
}

// CssRenamingMap_ByWhole is the ByWhole renaming map: an abstract class
// whose subclasses embed CssRenamingMap_ByWholeBase and implement Get.
type CssRenamingMap_ByWhole interface {
	CssRenamingMap
	GetCssRenamingMap_ByWholeBase() *CssRenamingMap_ByWholeBase
}

// CssRenamingMap_ByWholeBase holds ByWhole's base implementation.
type CssRenamingMap_ByWholeBase struct{}

// GetCssRenamingMap_ByWholeBase gives access to the base through the
// interface (and marks ByWhole subclasses for `instanceof`).
func (b *CssRenamingMap_ByWholeBase) GetCssRenamingMap_ByWholeBase() *CssRenamingMap_ByWholeBase {
	return b
}

// Java: CssRenamingMap.ByWhole.getStyle()
func (b *CssRenamingMap_ByWholeBase) GetStyle() CssRenamingMap_Style {
	return CssRenamingMap_Style_BY_WHOLE
}
