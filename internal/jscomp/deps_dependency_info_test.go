// Port of test/com/google/javascript/jscomp/deps/DependencyInfoTest.java @ closure-compiler v20261015

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestDependencyInfo_WriteAddDependency(t *testing.T) {
	sb := jutil.NewStringBuilder()
	flags := jutil.NewLinkedHashMap[string, string]()
	flags.Put("module", "goog")
	flags.Put("lang", "es6")
	DependencyInfo_Util_WriteAddDependency(
		sb,
		SimpleDependencyInfo_Builder_("../some/relative/path.js", "/unused/absolute/path.js").
			SetProvides([]string{"provided.symbol", "other.provide"}).
			SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("required.symbol"), DependencyInfo_Require_GoogRequireSymbol("other.require")).
			SetLoadFlags(flags).
			Build())
	jsaEq(t, sb.String(),
		"goog.addDependency('../some/relative/path.js', ['provided.symbol', 'other.provide'], "+
			"['required.symbol', 'other.require'], {'module': 'goog', 'lang': 'es6'});\n")
}

func TestDependencyInfo_WriteAddDependency_emptyArguments(t *testing.T) {
	sb := jutil.NewStringBuilder()
	DependencyInfo_Util_WriteAddDependency(sb, SimpleDependencyInfo_Builder_("path.js", "unused.js").Build())
	jsaEq(t, sb.String(), "goog.addDependency('path.js', [], []);\n")
}
