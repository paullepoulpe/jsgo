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
// Port of com/google/debugging/sourcemap/SourceMapGenerator.java (generated mechanically by j2go).

package sourcemap

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.debugging.sourcemap.SourceMapGenerator
type SourceMapGenerator interface {
	Implements_sourcemap_SourceMapGenerator()
	IsJavaNil_() bool
	AppendTo(j2gort.Ext_java_lang_Appendable, string)
	AppendIndexMapTo(j2gort.Ext_java_lang_Appendable, string, j2gort.Ext_java_util_List)
	Reset()
	AddMapping(string, string, *FilePosition, *FilePosition, *FilePosition)
	AddSourcesContent(string, string)
	SetWrapperPrefix(string)
	SetStartingPosition(int32, int32)
	Validate(bool)
}
func ToIface_SourceMapGenerator[T interface {
	SourceMapGenerator
	comparable
}](v T) SourceMapGenerator {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapGenerator", "com.google.debugging.sourcemap.SourceMapGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[SourceMapGenerator](v_)
	})
}
