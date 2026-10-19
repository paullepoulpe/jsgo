// Port of com/google/javascript/jscomp/parsing/ParserConfiguration.java @ closure-compiler v20261015

/*
 * Copyright 2020 The Closure Compiler Authors.
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

package parsing

import (
	_ "embed"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: language-adapter — classpath resource via go:embed. This is
// com/google/javascript/jscomp/parsing/ParserConfig.properties, extracted
// byte for byte from the pinned closure-compiler.jar.
//
//go:embed ParserConfig.properties
var parserConfiguration_properties string

// Go-only: language-adapter — ResourceBundle.getBundle loads (and caches) the
// PropertyResourceBundle once; the mutex is the JDK's own synchronization of
// that cache.
var (
	parserConfiguration_bundleOnce sync.Once
	parserConfiguration_bundle     *jutil.LinkedHashMap[string, string]
)

// ParserConfiguration reads the parser's ParserConfig.properties resource.
//
// Java: final class ParserConfiguration
type ParserConfiguration struct{}

// ParserConfiguration_GetString is
// ResourceBundle.getBundle("com.google.javascript.jscomp.parsing.ParserConfig").getString(name).
//
// Java: ParserConfiguration.getString(String)
func ParserConfiguration_GetString(name string) string {
	parserConfiguration_bundleOnce.Do(func() {
		parserConfiguration_bundle = jutil.PropertiesLoad(parserConfiguration_properties)
	})
	value, ok := parserConfiguration_bundle.Lookup(name)
	if !ok {
		panic(jutil.NewMissingResourceException(
			"Can't find resource for bundle java.util.PropertyResourceBundle, key "+name,
			"java.util.PropertyResourceBundle", name))
	}
	return value
}

// Java: ParserConfiguration()
func NewParserConfiguration() *ParserConfiguration { return &ParserConfiguration{} }
