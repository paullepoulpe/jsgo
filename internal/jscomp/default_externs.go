// Port of com/google/javascript/jscomp/DefaultExterns.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// DefaultExterns contains information on default externs files. Provide
// DefaultExterns#prepareExterns to filter externs source files.
type DefaultExterns struct {
	// Go-only: language-adapter — distinct identity (Java gives every `new`
	// one; a zero-size struct would not).
	_ byte
}

// Java: DefaultExterns()
func NewDefaultExterns() *DefaultExterns {
	return &DefaultExterns{}
}

// Core language externs. When the environment is CUSTOM, only these externs
// will be included.
var defaultExterns_BUILTIN_LANG_EXTERNS = []string{"es3.js", "es5.js", "es6.js", "es6_typed_array.js", "es6_collections.js"}

// Ordered browser externs. Externs not included in this list are added
// last.
var defaultExterns_BROWSER_EXTERN_DEP_ORDER = []string{
	"intl.js",
	"w3c_event.js",
	"w3c_event3.js",
	"gecko_event.js",
	"ie_event.js",
	"webkit_event.js",
	"w3c_device_sensor_event.js",
	"w3c_dom1.js",
	"w3c_dom2.js",
	"w3c_dom3.js",
	"w3c_dom4.js",
	"gecko_dom.js",
	"ie_dom.js",
	"webkit_dom.js",
	"w3c_css.js",
	"gecko_css.js",
	"ie_css.js",
	"webkit_css.js",
	"w3c_touch_event.js",
	"webstorage.js",
	"whatwg_console.js",
	"nonstandard_console.js",
	// w3d_rtc.js must be evaluated before nonstandard_rtc.js
	"w3c_rtc.js",
}

// DefaultExterns_PrepareExterns filters and orders the passed externs for
// the specified environment.
//
// env is the environment being used; externs is the flat filename to
// source externs map. Must be mutable and will be modified. Returns the
// ordered list of externs. (The remaining externs come in the map's order;
// the only caller, AbstractCommandLineRunner, passes a LinkedHashMap.)
//
// Java: DefaultExterns.prepareExterns(CompilerOptions.Environment, Map<String, SourceFile>)
func DefaultExterns_PrepareExterns(env CompilerOptions_Environment,
	externs *jutil.LinkedHashMap[string, *SourceFile]) []*SourceFile {
	var out []*SourceFile

	for _, key := range defaultExterns_BUILTIN_LANG_EXTERNS {
		jutil.CheckStateMsg1(externs.ContainsKey(key), "Externs must contain builtin: %s", key)
		out = append(out, externs.Remove(key))
	}

	if env == CompilerOptions_Environment_BROWSER {
		for _, key := range defaultExterns_BROWSER_EXTERN_DEP_ORDER {
			jutil.CheckStateMsg2(externs.ContainsKey(key),
				"Externs must contain builtin for env %s: %s", env, key)
			out = append(out, externs.Remove(key))
		}

		for value := range externs.Values() {
			out = append(out, value)
		}
	}

	if out == nil {
		out = []*SourceFile{}
	}
	return out
}
