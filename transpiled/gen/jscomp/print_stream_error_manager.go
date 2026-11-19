/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PrintStreamErrorManager.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.PrintStreamErrorManager
type PrintStreamErrorManager struct {
	BasicErrorManagerBase
	formatter MessageFormatter `gson:"formatter"`
	stream *j2gort.Ext_java_io_PrintStream `gson:"stream"`
	summaryDetailLevel int32 `gson:"summaryDetailLevel"`
}

func (this_ *PrintStreamErrorManager) BindSelf_(s any) {
	_ = s
	this_.BasicErrorManagerBase.BindSelf_(s)
}
func (this_ *PrintStreamErrorManager) IsJavaNil_() bool { return this_ == nil }
func (*PrintStreamErrorManager) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.PrintStreamErrorManager", "PrintStreamErrorManager" }
func (this_ *PrintStreamErrorManager) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.BasicErrorManagerBase.IdentityHash_()
}
func (*PrintStreamErrorManager) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PrintStreamErrorManager", "com.google.javascript.jscomp.PrintStreamErrorManager", func(v_ any) bool {
		return j2gort.InstanceOf[*PrintStreamErrorManager](v_)
	})
}
func (*PrintStreamErrorManager) Implements_jscomp_BasicErrorManager() {}
func (*PrintStreamErrorManager) Implements_jscomp_SortingErrorManager() {}
func (*PrintStreamErrorManager) Implements_jscomp_ErrorManager() {}
func (*PrintStreamErrorManager) Implements_jscomp_ErrorHandler() {}

func (this_ *PrintStreamErrorManager) initFields_PrintStreamErrorManager() {
	_ = this_
	this_.summaryDetailLevel = int32(1)
}

// Java constructor: PrintStreamErrorManager(com.google.javascript.jscomp.MessageFormatter formatter, java.io.PrintStream stream)
func NewPrintStreamErrorManager_2(formatter MessageFormatter, stream *j2gort.Ext_java_io_PrintStream) *PrintStreamErrorManager {
	this_ := &PrintStreamErrorManager{}
	this_.BindSelf_(this_)
	InitPrintStreamErrorManager_2(this_, formatter, stream)
	return this_
}

func InitPrintStreamErrorManager_2(this_ *PrintStreamErrorManager, formatter MessageFormatter, stream *j2gort.Ext_java_io_PrintStream) {
	_ = formatter
	_ = stream
	InitBasicErrorManager(&this_.BasicErrorManagerBase)
	this_.initFields_PrintStreamErrorManager()
	this_.formatter = formatter
	this_.stream = stream
}

// Java constructor: PrintStreamErrorManager(java.io.PrintStream stream)
func NewPrintStreamErrorManager_1(stream *j2gort.Ext_java_io_PrintStream) *PrintStreamErrorManager {
	this_ := &PrintStreamErrorManager{}
	this_.BindSelf_(this_)
	InitPrintStreamErrorManager_1(this_, stream)
	return this_
}

func InitPrintStreamErrorManager_1(this_ *PrintStreamErrorManager, stream *j2gort.Ext_java_io_PrintStream) {
	_ = stream
	InitPrintStreamErrorManager_2(this_, StaticGet_ErrorFormat_SOURCELESS().ToFormatter(nil, false), stream)
}

// Java method: PrintStreamErrorManager.println(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.JSError error)
func (this_ *PrintStreamErrorManager) Println(level *CheckLevel, error_ *JSError) {
	_ = level
	_ = error_
	this_.stream.Println_1_string(error_.Format(level, this_.formatter))
}

// Java method: PrintStreamErrorManager.setSummaryDetailLevel(int summaryDetailLevel)
func (this_ *PrintStreamErrorManager) SetSummaryDetailLevel(summaryDetailLevel int32) {
	_ = summaryDetailLevel
	this_.summaryDetailLevel = summaryDetailLevel
}

// Java method: PrintStreamErrorManager.printSummary()
func (this_ *PrintStreamErrorManager) PrintSummary() {
	if ((((this_.summaryDetailLevel >= int32(3)) || (((this_.summaryDetailLevel >= int32(1)) && ((this_.GetErrorCount() + this_.GetWarningCount()) > int32(0))))) || (((this_.summaryDetailLevel >= int32(2)) && (this_.GetTypedPercent() > float64(0.0)))))) {
		if ((this_.GetTypedPercent() > float64(0.0))) {
			this_.stream.Printf_2("%d error(s), %d warning(s), %.1f%% typed%n", []any{j2gort.ValueOfInt(this_.GetErrorCount()), j2gort.ValueOfInt(this_.GetWarningCount()), j2gort.ValueOfDouble(this_.GetTypedPercent())})
		} else {
			this_.stream.Printf_2("%d error(s), %d warning(s)%n", []any{j2gort.ValueOfInt(this_.GetErrorCount()), j2gort.ValueOfInt(this_.GetWarningCount())})
		}
	}
}

