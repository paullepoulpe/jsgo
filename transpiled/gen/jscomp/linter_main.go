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
// Port of com/google/javascript/jscomp/LinterMain.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.LinterMain
type LinterMain struct {
	files          j2gort.Ext_java_util_List `gson:"files"`
	idHash_        atomic.Int32
	fix            bool                        `gson:"fix"`
	resolutionMode ModuleLoader_ResolutionMode `gson:"resolutionMode"`
}

func (this_ *LinterMain) BindSelf_(s any) {
	_ = s
}
func (this_ *LinterMain) IsJavaNil_() bool { return this_ == nil }
func (*LinterMain) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LinterMain", "LinterMain"
}
func (this_ *LinterMain) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinterMain) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LinterMain", "com.google.javascript.jscomp.LinterMain", func(v_ any) bool {
		return j2gort.InstanceOf[*LinterMain](v_)
	})
}

func (this_ *LinterMain) initFields_LinterMain() {
	_ = this_

	this_.resolutionMode = StaticGet_ModuleLoader_ResolutionMode_BROWSER()
	this_.files = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: LinterMain()
func NewLinterMain() *LinterMain {
	this_ := &LinterMain{}
	InitLinterMain(this_)
	return this_
}

func InitLinterMain(this_ *LinterMain) {
	this_.initFields_LinterMain()
}

// Java static method: LinterMain.main(java.lang.String[] args)
func LinterMain_Main(args []string) {
	_ = args
	{
		var val_1 LinterMain
		InitLinterMain(&val_1)
		(&val_1).priv_LinterMain_run_3791(args)
	}
}

// Java method: LinterMain.run(java.lang.String[] args)
func (this_ *LinterMain) priv_LinterMain_run_3791(args []string) {
	_ = args
	var parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser = j2gort.NewExt_org_kohsuke_args4j_CmdLineParser_1(this_)
	_ = parser
	parser.ParseArgument_1_slice_string(args)
	var linter *Linter = Linter_S_Builder().WithModuleResolutionMode(this_.resolutionMode).Build()
	_ = linter
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.files))
		for cur_2.HasNext() {
			filename := cur_2.NextStr()
			_ = filename
			if this_.fix {
				linter.FixRepeatedly_1(filename)
			} else {
				linter.Lint_1(filename)
			}
		}
	}
}

func (this_ *LinterMain) RegisterArgs4jOptions(p *j2gort.Ext_org_kohsuke_args4j_CmdLineParser) {
	_ = this_
	_ = p
	{
		opt_0_ := j2gort.NewArgs4jNamedOptionDef("--fix", nil, "Fix lint warnings automatically", "", false, false, false, false)
		set_0_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.fix = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.fix)
		})
		j2gort.Args4jRegisterHandler(p, opt_0_, j2gort.NewArgs4jBuiltinHandler(p, opt_0_, set_0_, "boolean", ""), false)
	}
	{
		opt_1_ := j2gort.NewArgs4jNamedOptionDef("--resolution_mode", nil, "Expected path resolution mode for resolving module specifiers.", "", false, false, false, false)
		set_1_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode", "ResolutionMode"), func(val_ any) {
			this_.resolutionMode = func() ModuleLoader_ResolutionMode {
				_ = ModuleLoader_ResolutionMode_Values()
				return j2gort.Cast[ModuleLoader_ResolutionMode](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode", val_))
			}()
		}, func() any {
			return this_.resolutionMode
		})
		var iifeArg_1 j2gort.
			Ext_org_kohsuke_args4j_spi_OptionHandler
		{

			_ = ModuleLoader_ResolutionMode_Values()
			iifeArg_1 = j2gort.NewArgs4jBuiltinHandler(p, opt_1_, set_1_, "enum", "com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode")
		}
		j2gort.Args4jRegisterHandler(p, opt_1_, iifeArg_1, false)
	}
	{
		opt_2_ := j2gort.NewArgs4jArgumentDef("", "", false, false, true, 0)
		set_2_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.files == nil {
				this_.files = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.files, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_2_, j2gort.NewArgs4jBuiltinHandler(p, opt_2_, set_2_, "string", ""), true)
	}
}
