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
// Port of com/google/javascript/jscomp/instrumentation/reporter/ProductionInstrumentationReporter.java (generated mechanically by j2go).

package reporter

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.instrumentation.reporter.ProductionInstrumentationReporter
type ProductionInstrumentationReporter struct {
	instrumentationMappingLocation  string `gson:"instrumentationMappingLocation"`
	instrumentationReportsDirectory string `gson:"instrumentationReportsDirectory"`
	finalResultOutput               string `gson:"finalResultOutput"`
	idHash_                         atomic.Int32
}

func (this_ *ProductionInstrumentationReporter) BindSelf_(s any) {
	_ = s
}
func (this_ *ProductionInstrumentationReporter) IsJavaNil_() bool { return this_ == nil }
func (*ProductionInstrumentationReporter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.reporter.ProductionInstrumentationReporter", "ProductionInstrumentationReporter"
}
func (this_ *ProductionInstrumentationReporter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProductionInstrumentationReporter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.reporter.ProductionInstrumentationReporter", "com.google.javascript.jscomp.instrumentation.reporter.ProductionInstrumentationReporter", func(v_ any) bool {
		return j2gort.InstanceOf[*ProductionInstrumentationReporter](v_)
	})
}

func (this_ *ProductionInstrumentationReporter) initFields_ProductionInstrumentationReporter() {
	_ = this_

}

// Java constructor: ProductionInstrumentationReporter()
func NewProductionInstrumentationReporter() *ProductionInstrumentationReporter {
	this_ := &ProductionInstrumentationReporter{}
	InitProductionInstrumentationReporter(this_)
	return this_
}

func InitProductionInstrumentationReporter(this_ *ProductionInstrumentationReporter) {
	this_.initFields_ProductionInstrumentationReporter()
}

// Java static method: ProductionInstrumentationReporter.main(java.lang.String[] args)
func ProductionInstrumentationReporter_Main(args []string) {
	_ = args
	func() {
		defer func() {
			if rec_1 := recover(); rec_1 != nil {
				jex_2 := javalang.RecoverException(rec_1)
				_ = jex_2
				if j2gort.InstanceOf[javaio.IIOException](jex_2) {
					e := j2gort.AsIface[javaio.IIOException](jex_2)
					_ = e
					j2gort.Ext_java_lang_System_F_Err.Println_1_string(j2gort.ThrowableGetMessage(e))
				} else {
					panic(rec_1)
				}
			}
		}()
		{
			var val_3 ProductionInstrumentationReporter
			InitProductionInstrumentationReporter(&val_3)
			(&val_3).priv_ProductionInstrumentationReporter_doMain_6533(args)
		}
	}()
}

// Java static method: ProductionInstrumentationReporter.readFile(java.lang.String filePath)
func ProductionInstrumentationReporter_ReadFile(filePath string) string {
	_ = filePath
	return j2gort.Ext_com_google_common_io_CharStreams_ToString(j2gort.AsIface[j2gort.Ext_java_lang_Readable](j2gort.Ext_java_nio_file_Files_NewBufferedReader_2(j2gort.Ext_java_nio_file_Path_Of_2(filePath, []string{}), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)))
}

// Java method: ProductionInstrumentationReporter.getAllExecutionResults(java.io.File folder)
func (this_ *ProductionInstrumentationReporter) priv_ProductionInstrumentationReporter_getAllExecutionResults_6531(folder *j2gort.Ext_java_io_File) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = folder
	var result []j2gort.Ext_java_util_Map = []j2gort.Ext_java_util_Map{}
	_ = result
	for _, el_4 := range j2gort.CheckNotNullSlice(folder.ListFiles_0()) {
		fileEntry := el_4
		_ = fileEntry
		var executionResult string = ProductionInstrumentationReporter_ReadFile(fileEntry.GetAbsolutePath())
		_ = executionResult
		var type_ any = NewProductionInstrumentationReporter__Anon_121(this_).GetType()
		_ = type_
		var executedInstrumentationData j2gort.Ext_java_util_Map = j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.NewExt_com_google_gson_Gson().FromJson_2_string_any(executionResult, type_))
		_ = executedInstrumentationData
		result = append(result, j2gort.CheckNotNullElem(executedInstrumentationData))
	}
	return j2gort.ImmListFromTypedSlice(result)
}

// Java method: ProductionInstrumentationReporter.createFile(java.lang.String fileName, java.lang.String fileContents)
func (this_ *ProductionInstrumentationReporter) priv_ProductionInstrumentationReporter_createFile_6532(fileName string, fileContents string) {
	_ = fileName
	_ = fileContents
	var fold *j2gort.Ext_java_io_File = j2gort.NewExt_java_io_File_1_string(fileName)
	_ = fold
	fold.Delete()
	var myObj *j2gort.Ext_java_io_File = j2gort.NewExt_java_io_File_1_string(fileName)
	_ = myObj
	myObj.CreateNewFile()
	func() {
		var myWriter j2gort.Ext_java_io_Writer = j2gort.AsIface[j2gort.Ext_java_io_Writer](j2gort.Ext_java_nio_file_Files_NewBufferedWriter_3(j2gort.Ext_java_nio_file_Path_Of_2(fileName, []string{}), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8, []j2gort.Ext_java_nio_file_OpenOption{}))
		_ = myWriter
		defer j2gort.CloseResource(myWriter)
		myWriter.Write_1_string(fileContents)
	}()
}

// Java method: ProductionInstrumentationReporter.doMain(java.lang.String[] args)
func (this_ *ProductionInstrumentationReporter) priv_ProductionInstrumentationReporter_doMain_6533(args []string) {
	_ = args
	this_.priv_ProductionInstrumentationReporter_parseCmdLineArguments_6534(args)
	var mapping j2gort.Ext_java_util_Map = ReportDecoder_ParseMappingFromFile(this_.instrumentationMappingLocation)
	_ = mapping
	var folder *j2gort.Ext_java_io_File = j2gort.NewExt_java_io_File_1_string(this_.instrumentationReportsDirectory)
	_ = folder
	var listOfExecutionResults j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_ProductionInstrumentationReporter_getAllExecutionResults_6531(folder)
	_ = listOfExecutionResults
	this_.priv_ProductionInstrumentationReporter_createFile_6532(this_.finalResultOutput, ProfilingReport_CreateTabSeparatedProfilingReport(mapping, listOfExecutionResults))
}

// Java method: ProductionInstrumentationReporter.parseCmdLineArguments(java.lang.String[] args)
func (this_ *ProductionInstrumentationReporter) priv_ProductionInstrumentationReporter_parseCmdLineArguments_6534(args []string) {
	_ = args
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(args)); i = i + 1 {
			var arg string = args[i]
			_ = arg
			var name string = arg
			_ = name
			var val string = j2gort.NullString
			_ = val
			var eq int32 = j2gort.JStrIndexOfCp(arg, int32(uint16(61)))
			_ = eq
			if eq >= int32(0) {
				name = j2gort.JStrSub(arg, int32(0), eq)
				val = j2gort.JStrSubFrom(arg, (eq + int32(1)))
			} else if (i + int32(1)) < int32(len(args)) {
				val = args[func() int32 { old_5 := i; nw_6 := old_5 + 1; i = nw_6; return nw_6 }()]
			}
			if val == j2gort.NullString {
				j2gort.Ext_java_lang_System_F_Err.Println_1_string(("Option \"" + j2gort.StrValString(name) + "\" takes an operand"))
				return
			}
			switch j2gort.CheckNotNullString(name) {
			case "--mapping_file":
				this_.instrumentationMappingLocation = val
			case "--reports_directory":
				this_.instrumentationReportsDirectory = val
			case "--result_output":
				this_.finalResultOutput = val
			default:
				j2gort.Ext_java_lang_System_F_Err.Println_1_string(("\"" + j2gort.StrValString(arg) + "\" is not a valid option"))
				return
			}
		}
	}
}

// Java class:
type ProductionInstrumentationReporter__Anon_121 struct {
	j2gort.Ext_com_google_gson_reflect_TypeTokenBase
	outer_ *ProductionInstrumentationReporter
}

func (this_ *ProductionInstrumentationReporter__Anon_121) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_gson_reflect_TypeTokenBase.BindExtSelf_(s)
	this_.Ext_com_google_gson_reflect_TypeTokenBase.Core_ = "java.util.Map<java.lang.String,java.lang.Long>"
}
func (this_ *ProductionInstrumentationReporter__Anon_121) IsJavaNil_() bool { return this_ == nil }
func (*ProductionInstrumentationReporter__Anon_121) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.reporter.ProductionInstrumentationReporter$1", ""
}
func (this_ *ProductionInstrumentationReporter__Anon_121) SetOuter_(o *ProductionInstrumentationReporter) {
	this_.outer_ = o
}
func (this_ *ProductionInstrumentationReporter__Anon_121) GetOuter_() *ProductionInstrumentationReporter {
	return this_.outer_
}
func (*ProductionInstrumentationReporter__Anon_121) Implements_j2gort_Ext_com_google_gson_reflect_TypeToken() {
}

func (this_ *ProductionInstrumentationReporter__Anon_121) initFields_ProductionInstrumentationReporter__Anon_121() {
	_ = this_
}

func NewProductionInstrumentationReporter__Anon_121(outer_ *ProductionInstrumentationReporter) *ProductionInstrumentationReporter__Anon_121 {
	this_ := &ProductionInstrumentationReporter__Anon_121{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	j2gort.InitExt_com_google_gson_reflect_TypeToken(&this_.Ext_com_google_gson_reflect_TypeTokenBase)
	this_.initFields_ProductionInstrumentationReporter__Anon_121()
	return this_
}
