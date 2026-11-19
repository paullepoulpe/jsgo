/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/resources/ResourceLoader.java (generated mechanically by j2go).

package resources

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.resources.ResourceLoader
type ResourceLoader struct {
	idHash_ atomic.Int32
}

func (this_ *ResourceLoader) BindSelf_(s any) {
	_ = s
}
func (this_ *ResourceLoader) IsJavaNil_() bool { return this_ == nil }
func (*ResourceLoader) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.resources.ResourceLoader", "ResourceLoader"
}
func (this_ *ResourceLoader) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ResourceLoader) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.resources.ResourceLoader", "com.google.javascript.jscomp.resources.ResourceLoader", func(v_ any) bool {
		return j2gort.InstanceOf[*ResourceLoader](v_)
	})
}

func (this_ *ResourceLoader) initFields_ResourceLoader() {
	_ = this_
}

// Java constructor: ResourceLoader()
func NewResourceLoader() *ResourceLoader {
	this_ := &ResourceLoader{}
	InitResourceLoader(this_)
	return this_
}

func InitResourceLoader(this_ *ResourceLoader) {
	this_.initFields_ResourceLoader()
}

// Java static method: ResourceLoader.loadTextResource(java.lang.Class<?> clazz, java.lang.String path)
func ResourceLoader_LoadTextResource(clazz *javalang.Class, path string) string {
	_ = clazz
	_ = path
	var didRet_1 bool
	var retVal_2 string
	func() {
		defer func() {
			if rec_3 := recover(); rec_3 != nil {
				jex_4 := javalang.RecoverException(rec_3)
				_ = jex_4
				if j2gort.InstanceOf[javalang.INullPointerException](jex_4) {
					e := j2gort.AsIface[javalang.INullPointerException](jex_4)
					_ = e
					if !(ResourceLoader_ResourceExists(clazz, path)) {
						panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "No such resource: "+j2gort.StrValString(path))))
					}
					panic(j2gort.CheckNotNull(e))
				} else if j2gort.InstanceOf[javaio.IIOException](jex_4) {
					e_2 := j2gort.AsIface[javaio.IIOException](jex_4)
					_ = e_2
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e_2)))
				} else {
					panic(rec_3)
				}
			}
		}()
		var iifeArg_1 j2gort.
			Ext_java_io_InputStream
	iifeSwitch_1:
		switch {
		default:

			if b_, ok_ := j2gort.ClassResourceBytes(clazz, path); ok_ {
				s_ := &j2gort.Ext_java_io_InputStreamBase{Core_: b_}
				s_.Self_ = s_
				iifeArg_1 = s_
				break iifeSwitch_1
			}
			iifeArg_1 = nil
			break iifeSwitch_1
		}
		retVal_2 = j2gort.Ext_com_google_common_io_CharStreams_ToString(j2gort.NewExt_java_io_InputStreamReader_2_Ext_java_io_InputStream_Ext_java_nio_charset_Charset(iifeArg_1, j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8))
		didRet_1 = true
		return
	}()
	if didRet_1 {
		return retVal_2
	}
	panic("unreachable")
}

// Java static method: ResourceLoader.loadPropertiesMap(java.lang.Class<?> clazz, java.lang.String resourceName)
func ResourceLoader_LoadPropertiesMap(clazz *javalang.Class, resourceName string) j2gort.Ext_com_google_common_collect_ImmutableMap {
	_ = clazz
	_ = resourceName
	return PropertiesParser_Parse(ResourceLoader_LoadTextResource(clazz, resourceName))
}

// Java static method: ResourceLoader.resourceExists(java.lang.Class<?> clazz, java.lang.String path)
func ResourceLoader_ResourceExists(clazz *javalang.Class, path string) bool {
	_ = clazz
	_ = path
	return (func() *j2gort.Ext_java_net_URL {
		if j2gort.ClassResourceExists(clazz, path) {
			return &j2gort.Ext_java_net_URL{}
		}
		return nil
	}() != nil)
}
