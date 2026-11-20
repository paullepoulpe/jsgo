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
// Port of com/google/javascript/jscomp/PrintStreamErrorReportGenerator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.PrintStreamErrorReportGenerator
type PrintStreamErrorReportGenerator struct {
	formatter          MessageFormatter                `gson:"formatter"`
	stream             *j2gort.Ext_java_io_PrintStream `gson:"stream"`
	summaryDetailLevel int32                           `gson:"summaryDetailLevel"`
	idHash_            atomic.Int32
}

func (this_ *PrintStreamErrorReportGenerator) BindSelf_(s any) {
	_ = s
}
func (this_ *PrintStreamErrorReportGenerator) IsJavaNil_() bool { return this_ == nil }
func (*PrintStreamErrorReportGenerator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PrintStreamErrorReportGenerator", "PrintStreamErrorReportGenerator"
}
func (this_ *PrintStreamErrorReportGenerator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PrintStreamErrorReportGenerator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PrintStreamErrorReportGenerator", "com.google.javascript.jscomp.PrintStreamErrorReportGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[*PrintStreamErrorReportGenerator](v_)
	})
}
func (*PrintStreamErrorReportGenerator) Implements_jscomp_SortingErrorManager_ErrorReportGenerator() {
}

func (this_ *PrintStreamErrorReportGenerator) initFields_PrintStreamErrorReportGenerator() {
	_ = this_
}

// Java constructor: PrintStreamErrorReportGenerator(com.google.javascript.jscomp.MessageFormatter formatter, java.io.PrintStream stream, int summaryDetailLevel)
func NewPrintStreamErrorReportGenerator(formatter MessageFormatter, stream *j2gort.Ext_java_io_PrintStream, summaryDetailLevel int32) *PrintStreamErrorReportGenerator {
	this_ := &PrintStreamErrorReportGenerator{}
	InitPrintStreamErrorReportGenerator(this_, formatter, stream, summaryDetailLevel)
	return this_
}

func InitPrintStreamErrorReportGenerator(this_ *PrintStreamErrorReportGenerator, formatter MessageFormatter, stream *j2gort.Ext_java_io_PrintStream, summaryDetailLevel int32) {
	_ = formatter
	_ = stream
	_ = summaryDetailLevel
	this_.initFields_PrintStreamErrorReportGenerator()
	this_.formatter = formatter
	this_.stream = stream
	this_.summaryDetailLevel = summaryDetailLevel
}

// Java method: PrintStreamErrorReportGenerator.generateReport(com.google.javascript.jscomp.SortingErrorManager manager)
func (this_ *PrintStreamErrorReportGenerator) GenerateReport(manager SortingErrorManager) {
	_ = manager
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(manager.GetSortingErrorManagerBase().GetSortedDiagnostics()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			e := j2gort.CastPtr[SortingErrorManager_ErrorWithLevel](rawElem_2)
			_ = e
			this_.priv_PrintStreamErrorReportGenerator_println_4451(e.level, e.error_)
		}
	}
	this_.priv_PrintStreamErrorReportGenerator_printSummary_4452(Upcast_ErrorManager(manager))
}

// Java method: PrintStreamErrorReportGenerator.println(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.JSError error)
func (this_ *PrintStreamErrorReportGenerator) priv_PrintStreamErrorReportGenerator_println_4451(level *CheckLevel, error_ *JSError) {
	_ = level
	_ = error_
	this_.stream.Println_1_string(error_.Format(level, this_.formatter))
}

// Java method: PrintStreamErrorReportGenerator.printSummary(com.google.javascript.jscomp.ErrorManager manager)
func (this_ *PrintStreamErrorReportGenerator) priv_PrintStreamErrorReportGenerator_printSummary_4452(manager ErrorManager) {
	_ = manager
	if ((this_.summaryDetailLevel >= int32(3)) || ((this_.summaryDetailLevel >= int32(1)) && ((manager.GetErrorCount() + manager.GetWarningCount()) > int32(0)))) || ((this_.summaryDetailLevel >= int32(2)) && (manager.GetTypedPercent() > float64(0.0))) {
		if manager.GetTypedPercent() > float64(0.0) {
			this_.stream.Printf_2("%d error(s), %d warning(s), %.1f%% typed%n", []any{j2gort.ValueOfInt(manager.GetErrorCount()), j2gort.ValueOfInt(manager.GetWarningCount()), j2gort.ValueOfDouble(manager.GetTypedPercent())})
		} else {
			this_.stream.Printf_2("%d error(s), %d warning(s)%n", []any{j2gort.ValueOfInt(manager.GetErrorCount()), j2gort.ValueOfInt(manager.GetWarningCount())})
		}
	}
}
