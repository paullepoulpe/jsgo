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
// Port of com/google/javascript/jscomp/PerformanceTrackerCodeSizeEstimator.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.PerformanceTrackerCodeSizeEstimator
type PerformanceTrackerCodeSizeEstimator struct {
	CodeConsumerBase
	output      *j2gort.Ext_java_io_ByteArrayOutputStream `gson:"output"`
	stream      j2gort.Ext_java_util_zip_GZIPOutputStream `gson:"stream"`
	size        int32                                     `gson:"size"`
	lastChar    uint16                                    `gson:"lastChar"`
	trackGzSize bool                                      `gson:"trackGzSize"`
}

func (this_ *PerformanceTrackerCodeSizeEstimator) BindSelf_(s any) {
	_ = s
	this_.CodeConsumerBase.BindSelf_(s)
}
func (this_ *PerformanceTrackerCodeSizeEstimator) IsJavaNil_() bool { return this_ == nil }
func (*PerformanceTrackerCodeSizeEstimator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PerformanceTrackerCodeSizeEstimator", "PerformanceTrackerCodeSizeEstimator"
}
func (this_ *PerformanceTrackerCodeSizeEstimator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CodeConsumerBase.IdentityHash_()
}
func (*PerformanceTrackerCodeSizeEstimator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PerformanceTrackerCodeSizeEstimator", "com.google.javascript.jscomp.PerformanceTrackerCodeSizeEstimator", func(v_ any) bool {
		return j2gort.InstanceOf[*PerformanceTrackerCodeSizeEstimator](v_)
	})
}
func (*PerformanceTrackerCodeSizeEstimator) Implements_jscomp_CodeConsumer() {}

func (this_ *PerformanceTrackerCodeSizeEstimator) initFields_PerformanceTrackerCodeSizeEstimator() {
	_ = this_

}

// Java constructor: PerformanceTrackerCodeSizeEstimator(boolean trackGzSize)
func NewPerformanceTrackerCodeSizeEstimator(trackGzSize bool) *PerformanceTrackerCodeSizeEstimator {
	this_ := &PerformanceTrackerCodeSizeEstimator{}
	this_.BindSelf_(this_)
	InitPerformanceTrackerCodeSizeEstimator(this_, trackGzSize)
	return this_
}

func InitPerformanceTrackerCodeSizeEstimator(this_ *PerformanceTrackerCodeSizeEstimator, trackGzSize bool) {
	_ = trackGzSize
	InitCodeConsumer(&this_.CodeConsumerBase)
	this_.initFields_PerformanceTrackerCodeSizeEstimator()
	this_.trackGzSize = trackGzSize
	if trackGzSize {
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
			this_.output = j2gort.NewExt_java_io_ByteArrayOutputStream()
			this_.stream = j2gort.NewExt_java_util_zip_GZIPOutputStream_1(j2gort.AsIface[j2gort.Ext_java_io_OutputStream](this_.output))
		}()
	} else {
		this_.output = nil
		this_.stream = nil
	}
}

// Java static method: PerformanceTrackerCodeSizeEstimator.estimate(com.google.javascript.rhino.Node jsRoot, boolean trackGzSize)
func PerformanceTrackerCodeSizeEstimator_Estimate(jsRoot pkg_rhino.Node, trackGzSize bool) *PerformanceTrackerCodeSizeEstimator {
	_ = jsRoot
	_ = trackGzSize
	var estimator *PerformanceTrackerCodeSizeEstimator = NewPerformanceTrackerCodeSizeEstimator(trackGzSize)
	_ = estimator
	CodeGenerator_ForCostEstimation(estimator).Add_1_Node(jsRoot)
	return estimator
}

// Java method: PerformanceTrackerCodeSizeEstimator.append(java.lang.String str)
func (this_ *PerformanceTrackerCodeSizeEstimator) Append(str string) {
	_ = str
	var len_ int32 = j2gort.JStrLen(str)
	_ = len_
	if len_ > int32(0) {
		this_.size = (this_.size + len_)
		this_.lastChar = j2gort.JStrCharAt(str, (len_ - int32(1)))
		if this_.trackGzSize {
			func() {
				defer func() {
					if rec_3 := recover(); rec_3 != nil {
						jex_4 := javalang.RecoverException(rec_3)
						_ = jex_4
						if j2gort.InstanceOf[javaio.IIOException](jex_4) {
							e := j2gort.AsIface[javaio.IIOException](jex_4)
							_ = e
							panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
						} else {
							panic(rec_3)
						}
					}
				}()
				this_.stream.Write_1_slice_int8(func() []int8 {
					_ = j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8
					return j2gort.JStrGetBytes(str)
				}())
			}()
		}
	}
}

// Java method: PerformanceTrackerCodeSizeEstimator.getLastChar()
func (this_ *PerformanceTrackerCodeSizeEstimator) GetLastChar() uint16 {
	return this_.lastChar
}

// Java method: PerformanceTrackerCodeSizeEstimator.getCodeSize()
func (this_ *PerformanceTrackerCodeSizeEstimator) GetCodeSize() int32 {
	return this_.size
}

// Java method: PerformanceTrackerCodeSizeEstimator.getZippedCodeSize()
func (this_ *PerformanceTrackerCodeSizeEstimator) GetZippedCodeSize() int32 {
	if this_.trackGzSize {
		var didRet_5 bool
		var retVal_6 int32
		func() {
			defer func() {
				if rec_7 := recover(); rec_7 != nil {
					jex_8 := javalang.RecoverException(rec_7)
					_ = jex_8
					if j2gort.InstanceOf[javaio.IIOException](jex_8) {
						e := j2gort.AsIface[javaio.IIOException](jex_8)
						_ = e
						panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
					} else {
						panic(rec_7)
					}
				}
			}()
			this_.stream.Finish()
			this_.stream.Close_0()
			retVal_6 = this_.output.Size()
			didRet_5 = true
			return
		}()
		if didRet_5 {
			return retVal_6
		}
	} else {
		return int32(0)
	}
	panic("unreachable")
}
