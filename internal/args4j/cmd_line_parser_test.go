package args4j

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type customDelimBean struct {
	opt  string
	flag bool
	defs *jutil.OrderedMap[string, *string]
}

func (b *customDelimBean) Args4jFields() []AnnotatedField {
	return []AnnotatedField{
		{Option: &Option{Name: "-opt", Aliases: []string{"--option"}}, Field: String(&b.opt)},
		{Option: &Option{Name: "-f"}, Field: Bool(&b.flag)},
		{Option: &Option{Name: "-D"}, Field: OrderedStringMap(&b.defs)},
	}
}

func TestCustomOptionValueDelimiterColon(t *testing.T) {
	t.Run("simple key:val", func(t *testing.T) {
		bean := &customDelimBean{}
		props := ParserProperties_Defaults().WithOptionValueDelimiter(":")
		p := NewCmdLineParserWithProperties(bean, props)
		if err := p.TryParseArgument("-opt:val"); err != nil {
			t.Fatalf("unexpected parse error: %v", err)
		}
		if bean.opt != "val" {
			t.Fatalf("bean.opt = %q, want %q", bean.opt, "val")
		}
	})

	t.Run("value containing equals -opt:a=b", func(t *testing.T) {
		bean := &customDelimBean{}
		props := ParserProperties_Defaults().WithOptionValueDelimiter(":")
		p := NewCmdLineParserWithProperties(bean, props)
		if err := p.TryParseArgument("-opt:a=b"); err != nil {
			t.Fatalf("unexpected parse error: %v", err)
		}
		if bean.opt != "a=b" {
			t.Fatalf("bean.opt = %q, want %q", bean.opt, "a=b")
		}
	})

	t.Run("fallback to equals -opt=val", func(t *testing.T) {
		bean := &customDelimBean{}
		props := ParserProperties_Defaults().WithOptionValueDelimiter(":")
		p := NewCmdLineParserWithProperties(bean, props)
		if err := p.TryParseArgument("-opt=val"); err != nil {
			t.Fatalf("unexpected parse error: %v", err)
		}
		if bean.opt != "val" {
			t.Fatalf("bean.opt = %q, want %q", bean.opt, "val")
		}
	})

	t.Run("map option -D:key=value", func(t *testing.T) {
		bean := &customDelimBean{}
		props := ParserProperties_Defaults().WithOptionValueDelimiter(":")
		p := NewCmdLineParserWithProperties(bean, props)
		if err := p.TryParseArgument("-D:foo=bar"); err != nil {
			t.Fatalf("unexpected parse error: %v", err)
		}
		if bean.defs == nil {
			t.Fatal("bean.defs is nil")
		}
		got := bean.defs.Get("foo")
		if got == nil || *got != "bar" {
			t.Fatalf("bean.defs[\"foo\"] = %v, want %q", got, "bar")
		}
	})

	t.Run("multi-character delimiter", func(t *testing.T) {
		bean := &customDelimBean{}
		props := ParserProperties_Defaults().WithOptionValueDelimiter(":=")
		p := NewCmdLineParserWithProperties(bean, props)
		if err := p.TryParseArgument("--option:=x=1:2"); err != nil {
			t.Fatalf("unexpected parse error: %v", err)
		}
		if bean.opt != "x=1:2" {
			t.Fatalf("bean.opt = %q, want %q", bean.opt, "x=1:2")
		}
	})

	t.Run("empty delimiter still parses space-separated and equals flags", func(t *testing.T) {
		bean := &customDelimBean{}
		props := ParserProperties_Defaults().WithOptionValueDelimiter("")
		p := NewCmdLineParserWithProperties(bean, props)
		if err := p.TryParseArgument("-f", "-opt", "hello"); err != nil {
			t.Fatalf("unexpected parse error: %v", err)
		}
		if !bean.flag || bean.opt != "hello" {
			t.Fatalf("got flag=%v opt=%q, want true/hello", bean.flag, bean.opt)
		}
	})
}
