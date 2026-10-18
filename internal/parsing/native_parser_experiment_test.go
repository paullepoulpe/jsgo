package parsing

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type nativeExperimentReporter struct{ diagnostics []string }

func (r *nativeExperimentReporter) Warning(message, file string, line, column int) {
	r.diagnostics = append(r.diagnostics, fmt.Sprintf("warning:%s:%d:%d:%s", file, line, column, message))
}
func (r *nativeExperimentReporter) Error(message, file string, line, column int) {
	r.diagnostics = append(r.diagnostics, fmt.Sprintf("error:%s:%d:%d:%s", file, line, column, message))
}

func nativeExperimentCompare(t *testing.T, source string, config *Config, requireNative bool) {
	t.Helper()
	file := rhino.NewSimpleSourceFile("native.js", rhino.StaticSourceFile_SourceKind_STRONG)
	baselineReporter, experimentReporter := &nativeExperimentReporter{}, &nativeExperimentReporter{}
	baseline := parserRunnerParseUncached(file, source, config, baselineReporter)
	_, accepted := tryNativeParser(file, source, config)
	if accepted != requireNative {
		t.Fatalf("native acceptance %v, want %v for %q", accepted, requireNative, source)
	}
	experiment := parserRunnerParseWithNativeFallback(file, source, config, experimentReporter)
	if !reflect.DeepEqual(experimentReporter.diagnostics, baselineReporter.diagnostics) {
		t.Fatalf("diagnostics changed for %q: %v, want %v", source, experimentReporter.diagnostics, baselineReporter.diagnostics)
	}
	if baseline.Ast == nil || experiment.Ast == nil {
		if baseline.Ast != experiment.Ast {
			t.Fatalf("AST availability changed for %q", source)
		}
		return
	}
	// This includes node properties, literal values, UTF-16 line/column/length,
	// parenthesis/call flags and features, beyond equivalent printed JS alone.
	if got, want := experiment.Ast.ToStringTree(), baseline.Ast.ToStringTree(); got != want {
		t.Fatalf("AST changed for %q\n%s\nwant:\n%s", source, got, want)
	}
	if experiment.Features.String() != baseline.Features.String() ||
		!reflect.DeepEqual(experiment.Comments, baseline.Comments) ||
		!reflect.DeepEqual(experiment.TopLevelStatementRanges, baseline.TopLevelStatementRanges) ||
		!reflect.DeepEqual(experiment.SourceMapURL, baseline.SourceMapURL) {
		t.Fatalf("parse metadata changed for %q", source)
	}
}

func TestNativeParserExpressionSubset(t *testing.T) {
	cases := []string{
		"", " \n\t ", ";;", "var x;", "var x=1,y=2; x=x+y*3;",
		"'use strict'; var x=(1+2)*3;",
		"var data={name:'😀', flag:true, nested:[1,2,{value:null}], 'quoted key':0.125};",
		"var x=-a + +b * ~c; var y=typeof x; void y;",
		"var x=a?b=c:d?e:f;",
		"var x=a.b[2](1, b.c); obj.member=call(x);",
		"(a).b(c)[d]; (a)=1; a?b=c:d=e; a=b=c;",
		"var data={1:true, 'é':null}; var \\u0061=1; \\u0061;",
		"var o={1e309:0,1e-9999:1,0.0:2,2.:3}; var zero=-0;",
		"(obj.method)(); (obj['method'])(); 'custom';'use strict';x;",
		"var o={__proto__:1,__proto__:2};",
		"var x=a||b&&c|d^e&f==g<h<<i+j*k;",
		"var x=1e3 + .5 + 2.; var π='é'; π;",
		"var s='😀';\r\nvar t = {\n key : s\n};\u2028t.key;\n",
	}
	for _, mode := range []Config_LanguageMode{Config_LanguageMode_ECMASCRIPT3, Config_LanguageMode_ECMASCRIPT5, Config_LanguageMode_ES_NEXT} {
		for _, strict := range []Config_StrictMode{Config_StrictMode_STRICT, Config_StrictMode_SLOPPY} {
			config := ParserRunner_CreateConfig(mode, nil, strict)
			for _, source := range cases {
				nativeExperimentCompare(t, source, config, true)
			}
		}
	}
}

func TestNativeParserFallsBackBeforePublishing(t *testing.T) {
	config := ParserRunner_CreateConfig(Config_LanguageMode_ES_NEXT, nil, Config_StrictMode_STRICT)
	for _, source := range []string{
		"var x=1", "function f(){return 1;}", "{foo:1};", "let x=1;", "const {a}=b;",
		"var x=[1,,2];", "var x=[1,];", "var x={a:1,};", "f(1,);", "var x={a};",
		"var x=0x10;", "var x=01;", "var x=1_000;", "var x=1n;", "var x='a\\n';",
		"var x=/re/;", "var x=`template`;", "var x=a?.b;", "a++;", "var x=(a,b);",
		"var eval=1;", "var arguments=1;", "var x=#private;", "var x=(a?b:c)=d;",
		"var x=1; // late comment", "var x=1; /** @fileoverview Late */",
		"var x=1; /** @license License */", "var x=1;\n//# sourceMappingURL=test.map",
		"/** @closureUnaware */ var x=1;", "'use strict'; var x='\\1'; var y=;",
		"var x='unterminated", "var x=1; function (", "goog.module('example');",
		"var \\u0063lass=1;", "obj.default;", "var x=1; /* unterminated",
		"<!-- html comment\nvar x=1;", "#!/bin/node\nvar x=1;",
		strings.Repeat("(", 150) + "x" + strings.Repeat(")", 150) + ";",
	} {
		nativeExperimentCompare(t, source, config, false)
		keepGoing := config.ToBuilder().SetRunMode(Config_RunMode_KEEP_GOING).Build()
		nativeExperimentCompare(t, source, keepGoing, false)
	}
	withComments := config.ToBuilder().SetJsDocParsingMode(Config_JsDocParsing_INCLUDE_ALL_COMMENTS).Build()
	nativeExperimentCompare(t, "var x=1;", withComments, false)
}

func TestNativeParserGeneratedExpressions(t *testing.T) {
	rng := rand.New(rand.NewSource(651017))
	var expression func(int) string
	expression = func(depth int) string {
		if depth == 0 {
			return []string{"a", "3", "null", "true", "'😀'"}[rng.Intn(5)]
		}
		switch rng.Intn(8) {
		case 0:
			return "(" + expression(depth-1) + ")"
		case 1:
			return "f(" + expression(depth-1) + ", " + expression(depth-1) + ")"
		case 2:
			return "[" + expression(depth-1) + ", " + expression(depth-1) + "]"
		case 3:
			return "{a:" + expression(depth-1) + ", b:" + expression(depth-1) + "}"
		case 4:
			return "(" + expression(depth-1) + ").member"
		case 5:
			return "!(" + expression(depth-1) + ")"
		case 6:
			return "(" + expression(depth-1) + "?" + expression(depth-1) + ":" + expression(depth-1) + ")"
		default:
			return "(" + expression(depth-1) + []string{"+", "*", "||", "===", "<<"}[rng.Intn(5)] + expression(depth-1) + ")"
		}
	}
	config := ParserRunner_CreateConfig(Config_LanguageMode_ES_NEXT, nil, Config_StrictMode_STRICT)
	for range 200 {
		nativeExperimentCompare(t, "var result="+expression(3)+";", config, true)
	}
}

var nativeExperimentBenchmarkResult *ParserRunner_ParseResult

func BenchmarkNativeParserDataProgram(b *testing.B) {
	var source strings.Builder
	for i := 0; i < 128; i++ {
		fmt.Fprintf(&source, "var entry%d={id:%d, name:'item', enabled:true, values:[1,2,3], computed:base.member+%d*2};\n", i, i, i)
	}
	text := source.String()
	config := ParserRunner_CreateConfig(Config_LanguageMode_ES_NEXT, nil, Config_StrictMode_STRICT)
	file := rhino.NewSimpleSourceFile("bench.js", rhino.StaticSourceFile_SourceKind_STRONG)
	for _, native := range []bool{false, true} {
		b.Run(fmt.Sprintf("native=%v", native), func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(text)))
			for i := 0; i < b.N; i++ {
				if native {
					result, ok := tryNativeParser(file, text, config)
					if !ok {
						b.Fatal("benchmark left supported grammar")
					}
					nativeExperimentBenchmarkResult = result
				} else {
					nativeExperimentBenchmarkResult = parserRunnerParseUncached(file, text, config, rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE)
				}
			}
		})
	}
}
