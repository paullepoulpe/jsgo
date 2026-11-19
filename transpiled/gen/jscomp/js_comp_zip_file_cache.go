/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/JSCompZipFileCache.java (generated mechanically by j2go).

package jscomp

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_js_comp_zip_file_cache_1 j2gort.Ext_com_google_common_cache_RemovalListener

func init() {
	lambdaSingleton_js_comp_zip_file_cache_1 = j2gort.NewLambda_Ext_com_google_common_cache_RemovalListener(func(notification *j2gort.Ext_com_google_common_cache_RemovalNotification) {
		_ = notification
		func() {
			defer func() {
				if rec_1 := recover(); rec_1 != nil {
					jex_2 := javalang.RecoverException(rec_1)
					_ = jex_2
					if j2gort.InstanceOf[javaio.IIOException](jex_2) {
						e := j2gort.AsIface[javaio.IIOException](jex_2)
						_ = e
						panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
					} else {
						panic(rec_1)
					}
				}
			}()
			j2gort.CastPtr[JSCompZipFileCache_CachedZipFile](notification.GetValue()).priv_CachedZipFile_maybeClose_3551()
		}()
	})
}

// Java class: com.google.javascript.jscomp.JSCompZipFileCache
type JSCompZipFileCache struct {
	idHash_ atomic.Int32
}

func (this_ *JSCompZipFileCache) BindSelf_(s any) {
	_ = s
}
func (this_ *JSCompZipFileCache) IsJavaNil_() bool { return this_ == nil }
func (*JSCompZipFileCache) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JSCompZipFileCache", "JSCompZipFileCache"
}
func (this_ *JSCompZipFileCache) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSCompZipFileCache) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JSCompZipFileCache", "com.google.javascript.jscomp.JSCompZipFileCache", func(v_ any) bool {
		return j2gort.InstanceOf[*JSCompZipFileCache](v_)
	})
}

var static_JSCompZipFileCache_ZIP_CACHE_SIZE int32

func StaticGet_JSCompZipFileCache_ZIP_CACHE_SIZE() int32 {
	clinit_JSCompZipFileCache()
	return static_JSCompZipFileCache_ZIP_CACHE_SIZE
}
func StaticSet_JSCompZipFileCache_ZIP_CACHE_SIZE(v int32) {
	clinit_JSCompZipFileCache()
	static_JSCompZipFileCache_ZIP_CACHE_SIZE = v
}

var static_JSCompZipFileCache_zipFileCache j2gort.Ext_com_google_common_cache_LoadingCache

func StaticGet_JSCompZipFileCache_ZipFileCache() j2gort.Ext_com_google_common_cache_LoadingCache {
	clinit_JSCompZipFileCache()
	return static_JSCompZipFileCache_zipFileCache
}
func StaticSet_JSCompZipFileCache_ZipFileCache(v j2gort.Ext_com_google_common_cache_LoadingCache) {
	clinit_JSCompZipFileCache()
	static_JSCompZipFileCache_zipFileCache = v
}

var clinitDone_JSCompZipFileCache uint32
var clinitState_JSCompZipFileCache j2gort.ClinitState

func clinit_JSCompZipFileCache() {
	if atomic.LoadUint32(&clinitDone_JSCompZipFileCache) != 0 {
		return
	}
	clinitSlow_JSCompZipFileCache()
}

func clinitSlow_JSCompZipFileCache() {
	if !j2gort.EnterClinit(&clinitState_JSCompZipFileCache, &clinitDone_JSCompZipFileCache) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSCompZipFileCache, &clinitDone_JSCompZipFileCache)
	static_JSCompZipFileCache_ZIP_CACHE_SIZE = j2gort.ParseInt(j2gort.Ext_java_lang_System_GetProperty_2("jscomp.zipfile.cachesize", "1000"), 10)
	static_JSCompZipFileCache_zipFileCache = j2gort.Ext_com_google_common_cache_CacheBuilder_NewBuilder().MaximumSize(int64(StaticGet_JSCompZipFileCache_ZIP_CACHE_SIZE())).RemovalListener(lambdaSingleton_js_comp_zip_file_cache_1).Build_1(NewJSCompZipFileCache__Anon_55())
}

func (this_ *JSCompZipFileCache) initFields_JSCompZipFileCache() {
	_ = this_
}

// Java constructor: JSCompZipFileCache()
func NewJSCompZipFileCache() *JSCompZipFileCache {
	clinit_JSCompZipFileCache()
	this_ := &JSCompZipFileCache{}
	InitJSCompZipFileCache(this_)
	return this_
}

func InitJSCompZipFileCache(this_ *JSCompZipFileCache) {
	this_.initFields_JSCompZipFileCache()
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java static method: JSCompZipFileCache.getEntryStream(java.lang.String zipName, java.lang.String entryName)
func JSCompZipFileCache_GetEntryStream(zipName string, entryName string) j2gort.Ext_java_io_InputStream {
	clinit_JSCompZipFileCache()
	_ = zipName
	_ = entryName
	return j2gort.CastPtr[JSCompZipFileCache_CachedZipFile](StaticGet_JSCompZipFileCache_ZipFileCache().GetUnchecked(j2gort.BoxString(zipName))).GetEntryStream(entryName)
}

// Java class: com.google.javascript.jscomp.JSCompZipFileCache.CachedZipFile
type JSCompZipFileCache_CachedZipFile struct {
	path         j2gort.Ext_java_nio_file_Path                `gson:"path"`
	zipFile      *j2gort.Ext_java_util_zip_ZipFile            `gson:"zipFile"`
	lastModified *j2gort.Ext_java_nio_file_attribute_FileTime `gson:"lastModified"`
	_syncMu      sync.Mutex                                   `gson:"-"`
	idHash_      atomic.Int32
}

func (this_ *JSCompZipFileCache_CachedZipFile) BindSelf_(s any) {
	_ = s
}
func (this_ *JSCompZipFileCache_CachedZipFile) IsJavaNil_() bool { return this_ == nil }
func (*JSCompZipFileCache_CachedZipFile) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JSCompZipFileCache$CachedZipFile", "CachedZipFile"
}
func (this_ *JSCompZipFileCache_CachedZipFile) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSCompZipFileCache_CachedZipFile) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JSCompZipFileCache.CachedZipFile", "com.google.javascript.jscomp.JSCompZipFileCache$CachedZipFile", func(v_ any) bool {
		return j2gort.InstanceOf[*JSCompZipFileCache_CachedZipFile](v_)
	})
}

func (this_ *JSCompZipFileCache_CachedZipFile) initFields_JSCompZipFileCache_CachedZipFile() {
	_ = this_
}

// Java constructor: CachedZipFile(java.lang.String zipName)
func NewJSCompZipFileCache_CachedZipFile(zipName string) *JSCompZipFileCache_CachedZipFile {
	this_ := &JSCompZipFileCache_CachedZipFile{}
	InitJSCompZipFileCache_CachedZipFile(this_, zipName)
	return this_
}

func InitJSCompZipFileCache_CachedZipFile(this_ *JSCompZipFileCache_CachedZipFile, zipName string) {
	_ = zipName
	this_.initFields_JSCompZipFileCache_CachedZipFile()
	this_.path = j2gort.Ext_java_nio_file_Path_Of_2(zipName, []string{})
}

// Java method: CachedZipFile.getEntryStream(java.lang.String entryName)
func (this_ *JSCompZipFileCache_CachedZipFile) GetEntryStream(entryName string) j2gort.Ext_java_io_InputStream {
	_ = entryName
	this_.priv_CachedZipFile_refreshIfNeeded_3550()
	var entry *j2gort.Ext_java_util_zip_ZipEntry = j2gort.CheckNotNullPtr(this_.zipFile.GetEntry(entryName))
	_ = entry
	return this_.zipFile.GetInputStream(entry)
}

// Java method: CachedZipFile.refreshIfNeeded()
func (this_ *JSCompZipFileCache_CachedZipFile) priv_CachedZipFile_refreshIfNeeded_3550() {
	var newLastModified *j2gort.Ext_java_nio_file_attribute_FileTime = j2gort.Ext_java_nio_file_Files_GetLastModifiedTime(this_.path, []*j2gort.Ext_java_nio_file_LinkOption{})
	_ = newLastModified
	if newLastModified.Equals(j2gort.PtrToAny(this_.lastModified)) {
		return
	}
	{
		if this_ == nil {
			j2gort.ThrowNPE()
		}
		this_._syncMu.Lock()
		defer this_._syncMu.Unlock()
		if newLastModified.Equals(j2gort.PtrToAny(this_.lastModified)) {
			return
		}
		this_.priv_CachedZipFile_maybeClose_3551()
		this_.zipFile = j2gort.NewExt_java_util_zip_ZipFile_1_Ext_java_io_File(this_.path.ToFile())
		this_.lastModified = newLastModified
	}
}

// Java method: CachedZipFile.maybeClose()
func (this_ *JSCompZipFileCache_CachedZipFile) priv_CachedZipFile_maybeClose_3551() {
	if this_.zipFile != nil {
		this_.zipFile.Close_0()
	}
}

// Java class:
type JSCompZipFileCache__Anon_55 struct {
	j2gort.Ext_com_google_common_cache_CacheLoaderBase
	idHash_ atomic.Int32
}

func (this_ *JSCompZipFileCache__Anon_55) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_common_cache_CacheLoaderBase.BindExtSelf_(s)
}
func (this_ *JSCompZipFileCache__Anon_55) IsJavaNil_() bool { return this_ == nil }
func (*JSCompZipFileCache__Anon_55) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JSCompZipFileCache$1", ""
}
func (this_ *JSCompZipFileCache__Anon_55) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSCompZipFileCache__Anon_55) Implements_j2gort_Ext_com_google_common_cache_CacheLoader() {}

func (this_ *JSCompZipFileCache__Anon_55) initFields_JSCompZipFileCache__Anon_55() {
	_ = this_
}

func NewJSCompZipFileCache__Anon_55() *JSCompZipFileCache__Anon_55 {
	this_ := &JSCompZipFileCache__Anon_55{}
	this_.BindSelf_(this_)
	j2gort.InitExt_com_google_common_cache_CacheLoader(&this_.Ext_com_google_common_cache_CacheLoaderBase)
	this_.initFields_JSCompZipFileCache__Anon_55()
	return this_
}

// Java method: .load(java.lang.String key)
func (this_ *JSCompZipFileCache__Anon_55) Load(arg_3 any) any {
	key := j2gort.CastString(arg_3)
	_ = key
	return NewJSCompZipFileCache_CachedZipFile(key)
}
