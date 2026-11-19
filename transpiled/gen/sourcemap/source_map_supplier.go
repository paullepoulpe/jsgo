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
// Port of com/google/debugging/sourcemap/SourceMapSupplier.java (generated mechanically by j2go).

package sourcemap

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.debugging.sourcemap.SourceMapSupplier
type SourceMapSupplier interface {
	Implements_sourcemap_SourceMapSupplier()
	IsJavaNil_() bool
	GetSourceMap(string) string
}
func ToIface_SourceMapSupplier[T interface {
	SourceMapSupplier
	comparable
}](v T) SourceMapSupplier {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapSupplier", "com.google.debugging.sourcemap.SourceMapSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[SourceMapSupplier](v_)
	})
}
