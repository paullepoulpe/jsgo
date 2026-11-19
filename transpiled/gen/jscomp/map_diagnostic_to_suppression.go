/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/MapDiagnosticToSuppression.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.MapDiagnosticToSuppression
type MapDiagnosticToSuppression struct {
	idHash_ atomic.Int32
	output  DiagnosticToSuppressionMapper_OutputFormat `gson:"output"`
	help    bool                                       `gson:"help"`
}

func (this_ *MapDiagnosticToSuppression) BindSelf_(s any) {
	_ = s
}
func (this_ *MapDiagnosticToSuppression) IsJavaNil_() bool { return this_ == nil }
func (*MapDiagnosticToSuppression) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MapDiagnosticToSuppression", "MapDiagnosticToSuppression"
}
func (this_ *MapDiagnosticToSuppression) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MapDiagnosticToSuppression) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MapDiagnosticToSuppression", "com.google.javascript.jscomp.MapDiagnosticToSuppression", func(v_ any) bool {
		return j2gort.InstanceOf[*MapDiagnosticToSuppression](v_)
	})
}

func (this_ *MapDiagnosticToSuppression) initFields_MapDiagnosticToSuppression() {
	_ = this_
	this_.output = StaticGet_DiagnosticToSuppressionMapper_OutputFormat_JSON()

}

// Java constructor: MapDiagnosticToSuppression()
func NewMapDiagnosticToSuppression() *MapDiagnosticToSuppression {
	this_ := &MapDiagnosticToSuppression{}
	InitMapDiagnosticToSuppression(this_)
	return this_
}

func InitMapDiagnosticToSuppression(this_ *MapDiagnosticToSuppression) {
	this_.initFields_MapDiagnosticToSuppression()
}

// Java static method: MapDiagnosticToSuppression.main(java.lang.String[] args)
func MapDiagnosticToSuppression_Main(args []string) {
	_ = args
	{
		var val_1 MapDiagnosticToSuppression
		InitMapDiagnosticToSuppression(&val_1)
		(&val_1).priv_MapDiagnosticToSuppression_printSuppressions_3858(args)
	}
}

// Java static method: MapDiagnosticToSuppression.printHelp()
func MapDiagnosticToSuppression_PrintHelp() {
	j2gort.Ext_java_lang_System_F_Out.Println_1_string("Usage: MapDiagnosticToSuppression --output={json,md}")
}

// Java method: MapDiagnosticToSuppression.printSuppressions(java.lang.String[] args)
func (this_ *MapDiagnosticToSuppression) priv_MapDiagnosticToSuppression_printSuppressions_3858(args []string) {
	_ = args
	var parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser = j2gort.NewExt_org_kohsuke_args4j_CmdLineParser_1(this_)
	_ = parser
	var didRet_2 bool
	func() {
		defer func() {
			if rec_3 := recover(); rec_3 != nil {
				jex_4 := javalang.RecoverException(rec_3)
				_ = jex_4
				if j2gort.InstanceOf[j2gort.Ext_org_kohsuke_args4j_CmdLineException](jex_4) {
					ex := j2gort.AsIface[j2gort.Ext_org_kohsuke_args4j_CmdLineException](jex_4)
					_ = ex
					j2gort.Ext_java_lang_System_F_Err.Println_1_string(("Could not parse arguments: " + j2gort.StrVal(ex)))
					MapDiagnosticToSuppression_PrintHelp()
					didRet_2 = true
					return
				} else {
					panic(rec_3)
				}
			}
		}()
		parser.ParseArgument_1_slice_string(args)
	}()
	if didRet_2 {
		return
	}
	if this_.help {
		MapDiagnosticToSuppression_PrintHelp()
		return
	}
	{
		var val_5 DiagnosticToSuppressionMapper
		InitDiagnosticToSuppressionMapper(&val_5, j2gort.AsIface[j2gort.Ext_java_util_Set](pkg_parsing.ParserRunner_GetSuppressionNames()), DiagnosticGroups_GetRegisteredGroups())
		(&val_5).PrintSuppressions(this_.output)
	}
}

func (this_ *MapDiagnosticToSuppression) RegisterArgs4jOptions(p *j2gort.Ext_org_kohsuke_args4j_CmdLineParser) {
	_ = this_
	_ = p
	{
		opt_0_ := j2gort.NewArgs4jNamedOptionDef("--output", nil, "Sets the desired output format.\nOptions: md, jsonDefault: json", "", false, false, false, false)
		set_0_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.DiagnosticToSuppressionMapper.OutputFormat", "OutputFormat"), func(val_ any) {
			this_.output = func() DiagnosticToSuppressionMapper_OutputFormat {
				_ = DiagnosticToSuppressionMapper_OutputFormat_Values()
				return j2gort.Cast[DiagnosticToSuppressionMapper_OutputFormat](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.DiagnosticToSuppressionMapper.OutputFormat", val_))
			}()
		}, func() any {
			return this_.output
		})
		var iifeArg_1 j2gort.
			Ext_org_kohsuke_args4j_spi_OptionHandler
		{

			_ = DiagnosticToSuppressionMapper_OutputFormat_Values()
			iifeArg_1 = j2gort.NewArgs4jBuiltinHandler(p, opt_0_, set_0_, "enum", "com.google.javascript.jscomp.DiagnosticToSuppressionMapper.OutputFormat")
		}
		j2gort.Args4jRegisterHandler(p, opt_0_, iifeArg_1, false)
	}
	{
		opt_1_ := j2gort.NewArgs4jNamedOptionDef("--help", nil, "Usage", "", false, false, false, false)
		set_1_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.help = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.help)
		})
		j2gort.Args4jRegisterHandler(p, opt_1_, j2gort.NewArgs4jBuiltinHandler(p, opt_1_, set_1_, "boolean", ""), false)
	}
}
