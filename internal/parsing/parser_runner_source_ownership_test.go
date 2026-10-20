package parsing

import (
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type retainedSourceCommentRecorder struct{}

func (retainedSourceCommentRecorder) RecordComment(parser.Comment_Type, *parser.SourceRange, string) {
}

func retainedSourceConfig(mode Config_JsDocParsing) *Config {
	return ParserRunner_CreateConfigWithModes(Config_LanguageMode_ES_NEXT, mode,
		Config_RunMode_KEEP_GOING, nil, true, Config_StrictMode_STRICT)
}

func retainedSourceParse(code string, config *Config) *ParserRunner_ParseResult {
	return ParserRunner_Parse(rhino.NewSimpleSourceFile("published.js", rhino.StaticSourceFile_SourceKind_STRONG),
		code, config, rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE)
}

func findRetainedSourceComment(n *rhino.Node) *rhino.NonJSDocComment {
	if n == nil {
		return nil
	}
	if comment := n.GetNonJSDocComment(); comment != nil {
		return comment
	}
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		if comment := findRetainedSourceComment(child); comment != nil {
			return comment
		}
	}
	return n.GetTrailingNonJSDocComment()
}

func TestParserRunnerPublishedSourceFilesRemainUsable(t *testing.T) {
	for _, tc := range []struct {
		name       string
		code       string
		mode       Config_JsDocParsing
		statements int
		comments   int
	}{
		{"range", "const café = '😀';\n", Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE, 1, 0},
		{"comment", "// café 😀\n", Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE, 0, 1},
		{"astComment", "// café 😀\nconst value = '😀';\n", Config_JsDocParsing_INCLUDE_ALL_COMMENTS, 1, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := retainedSourceConfig(tc.mode)
			result := retainedSourceParse(tc.code, config)
			if result.Ast == nil || len(result.TopLevelStatementRanges) != tc.statements || len(result.Comments) != tc.comments {
				t.Fatal("fixture did not publish the expected source views")
			}
			var position *parser.SourcePosition
			switch tc.name {
			case "range":
				position = result.TopLevelStatementRanges[0].Start
			case "comment":
				position = result.Comments[0].Location.Start
			case "astComment":
				comment := findRetainedSourceComment(result.Ast)
				if comment == nil {
					t.Fatal("fixture did not retain an AST comment")
				}
				position = comment.GetStartPosition()
			}
			if position == nil || position.Source == nil {
				t.Fatal("published position lost its source file")
			}
			file := position.Source
			originalPosition := *position
			wantUTF16 := jutil.StringToUTF16(tc.code)
			check := func() {
				t.Helper()
				reporter := NewParserRunner_Es6ErrorReporter(rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE, true)
				p := parser.NewParser(ParserRunner_NewParserConfig(config), reporter, file)
				program := p.ParseProgram()
				if program == nil || len(program.SourceElements) != tc.statements || len(p.GetComments()) != tc.comments {
					t.Fatal("public parser could not reparse the published immutable source file")
				}
				scanner := parser.NewScanner(reporter, retainedSourceCommentRecorder{}, file, 0)
				for {
					token := scanner.NextToken().GetTokenBase()
					if token.Type == parser.TokenType_END_OF_FILE {
						if token.Location.Start.Offset != len(wantUTF16) {
							t.Fatal("public scanner lost the published file's UTF-16 offsets")
						}
						break
					}
				}
				if file.Name != "published.js" || file.Contents != tc.code || !slices.Equal(file.ContentsUTF16, wantUTF16) || *position != originalPosition {
					t.Fatal("published source contents or position changed")
				}
			}
			check()
			// A shorter parse fits the previous buffer capacity and would
			// overwrite it if the published source still belonged to the arena.
			for range 3 {
				retainedSourceParse("var q=0;", retainedSourceConfig(Config_JsDocParsing_TYPES_ONLY))
			}
			check()
		})
	}
}

func TestParserRunnerUnpublishedSourceStillReusesUTF16Arena(t *testing.T) {
	arena := parser.AcquireParseArena()
	seed := parser.NewSourceFileInArena("seed.js", strings.Repeat(" ", 256), arena)
	backing := &seed.ContentsUTF16[0]
	arena.Release()
	result := retainedSourceParse("const reusable=1;", retainedSourceConfig(Config_JsDocParsing_TYPES_ONLY))
	if result.Ast == nil || len(result.Comments) != 0 || len(result.TopLevelStatementRanges) != 0 {
		t.Fatal("ordinary parse unexpectedly published source views")
	}
	reused := parser.AcquireParseArena()
	defer reused.Release()
	probe := parser.NewSourceFileInArena("probe.js", "x", reused)
	if reused != arena || &probe.ContentsUTF16[0] != backing {
		t.Fatal("ordinary unpublished parse stopped recycling its UTF-16 backing storage")
	}
}
