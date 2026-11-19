/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/SourceMapInput.java (generated mechanically by j2go).

package jscomp

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_sourcemap "github.com/paullepoulpe/jsgo/transpiled/gen/sourcemap"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.SourceMapInput
type SourceMapInput struct {
	sourceFile      *SourceFile                        `gson:"sourceFile"`
	parsedSourceMap *pkg_sourcemap.SourceMapConsumerV3 `gson:"parsedSourceMap"`
	_syncMu         sync.Mutex                         `gson:"-"`
	idHash_         atomic.Int32
	cached          bool `gson:"cached"`
}

func (this_ *SourceMapInput) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapInput) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapInput) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceMapInput", "SourceMapInput"
}
func (this_ *SourceMapInput) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapInput) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceMapInput", "com.google.javascript.jscomp.SourceMapInput", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapInput](v_)
	})
}

var static_SourceMapInput_SOURCEMAP_RESOLVE_FAILED *DiagnosticType

func StaticGet_SourceMapInput_SOURCEMAP_RESOLVE_FAILED() *DiagnosticType {
	clinit_SourceMapInput()
	return static_SourceMapInput_SOURCEMAP_RESOLVE_FAILED
}
func StaticSet_SourceMapInput_SOURCEMAP_RESOLVE_FAILED(v *DiagnosticType) {
	clinit_SourceMapInput()
	static_SourceMapInput_SOURCEMAP_RESOLVE_FAILED = v
}

var static_SourceMapInput_SOURCEMAP_PARSE_FAILED *DiagnosticType

func StaticGet_SourceMapInput_SOURCEMAP_PARSE_FAILED() *DiagnosticType {
	clinit_SourceMapInput()
	return static_SourceMapInput_SOURCEMAP_PARSE_FAILED
}
func StaticSet_SourceMapInput_SOURCEMAP_PARSE_FAILED(v *DiagnosticType) {
	clinit_SourceMapInput()
	static_SourceMapInput_SOURCEMAP_PARSE_FAILED = v
}

var clinitDone_SourceMapInput uint32
var clinitState_SourceMapInput j2gort.ClinitState

func clinit_SourceMapInput() {
	if atomic.LoadUint32(&clinitDone_SourceMapInput) != 0 {
		return
	}
	clinitSlow_SourceMapInput()
}

func clinitSlow_SourceMapInput() {
	if !j2gort.EnterClinit(&clinitState_SourceMapInput, &clinitDone_SourceMapInput) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SourceMapInput, &clinitDone_SourceMapInput)
	static_SourceMapInput_SOURCEMAP_RESOLVE_FAILED = DiagnosticType_Warning("SOURCEMAP_RESOLVE_FAILED", "Failed to resolve sourcemap at {0}: {1}")
	static_SourceMapInput_SOURCEMAP_PARSE_FAILED = DiagnosticType_Warning("SOURCEMAP_PARSE_FAILED", "Failed to parse malformed sourcemap in {0}: {1}")
}

func (this_ *SourceMapInput) initFields_SourceMapInput() {
	_ = this_

}

// Java constructor: SourceMapInput(com.google.javascript.jscomp.SourceFile sourceFile)
func NewSourceMapInput(sourceFile *SourceFile) *SourceMapInput {
	clinit_SourceMapInput()
	this_ := &SourceMapInput{}
	InitSourceMapInput(this_, sourceFile)
	return this_
}

func InitSourceMapInput(this_ *SourceMapInput, sourceFile *SourceFile) {
	_ = sourceFile
	this_.initFields_SourceMapInput()
	this_.sourceFile = sourceFile
}

// Java method: SourceMapInput.getSourceMap(com.google.javascript.jscomp.ErrorManager errorManager)
func (this_ *SourceMapInput) GetSourceMap(errorManager ErrorManager) *pkg_sourcemap.SourceMapConsumerV3 {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	this_._syncMu.Lock()
	defer this_._syncMu.Unlock()
	_ = errorManager
	if !(this_.cached) {
		this_.cached = true
		var sourceMapPath string = this_.sourceFile.GetName()
		_ = sourceMapPath
		func() {
			defer func() {
				if rec_1 := recover(); rec_1 != nil {
					jex_2 := javalang.RecoverException(rec_1)
					_ = jex_2
					if j2gort.InstanceOf[javaio.IIOException](jex_2) {
						e := j2gort.AsIface[javaio.IIOException](jex_2)
						_ = e
						var error_ *JSError = JSError_Make_2(StaticGet_SourceMapInput_SOURCEMAP_RESOLVE_FAILED(), []string{sourceMapPath, j2gort.ThrowableGetMessage(e)})
						_ = error_
						errorManager.Report_2(error_.DefaultLevel(), error_)
					} else if j2gort.InstanceOf[*pkg_sourcemap.SourceMapParseException](jex_2) {
						e_2 := j2gort.AsIface[*pkg_sourcemap.SourceMapParseException](jex_2)
						_ = e_2
						var error__2 *JSError = JSError_Make_2(StaticGet_SourceMapInput_SOURCEMAP_PARSE_FAILED(), []string{sourceMapPath, j2gort.ThrowableGetMessage(e_2)})
						_ = error__2
						errorManager.Report_2(error__2.DefaultLevel(), error__2)
					} else {
						panic(rec_1)
					}
				}
			}()
			var sourceMapContents string = this_.sourceFile.GetCode()
			_ = sourceMapContents
			var consumer *pkg_sourcemap.SourceMapConsumerV3 = pkg_sourcemap.NewSourceMapConsumerV3()
			_ = consumer
			consumer.Parse_1(sourceMapContents)
			if ((((((consumer.GetFile() != j2gort.NullString) && !(j2gort.CheckNotNullString(consumer.GetFile()) == "")) && (consumer.GetOriginalSources().Size() == int32(1))) && consumer.GetOriginalSources().Contains(j2gort.BoxString(consumer.GetFile()))) && (consumer.GetOriginalSourcesContent().Size() == int32(1))) && consumer.GetOriginalNames().IsEmpty()) && (consumer.GetLineCount() == int32(-1)) {
				consumer = nil
			}
			this_.parsedSourceMap = consumer
		}()
	}
	return this_.parsedSourceMap
}

// Java method: SourceMapInput.getOriginalPath()
func (this_ *SourceMapInput) GetOriginalPath() string {
	return this_.sourceFile.GetName()
}

// Java method: SourceMapInput.getRawSourceMapContents()
func (this_ *SourceMapInput) GetRawSourceMapContents() string {
	var didRet_3 bool
	var retVal_4 string
	func() {
		defer func() {
			if rec_5 := recover(); rec_5 != nil {
				jex_6 := javalang.RecoverException(rec_5)
				_ = jex_6
				if j2gort.InstanceOf[javaio.IIOException](jex_6) {
					e := j2gort.AsIface[javaio.IIOException](jex_6)
					_ = e
					retVal_4 = j2gort.NullString
					didRet_3 = true
					return
				} else {
					panic(rec_5)
				}
			}
		}()
		retVal_4 = this_.sourceFile.GetCode()
		didRet_3 = true
		return
	}()
	if didRet_3 {
		return retVal_4
	}
	panic("unreachable")
}
