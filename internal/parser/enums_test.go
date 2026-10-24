package parser

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
)

// testdata/enums.txt is generated with the pinned Java compiler by
// tools/oracle/java/ParserEnumDump.java; its lines are:
//
//	enum <JavaClass> <ordinal> <name> <toString or <Exception>>
//	featureset <FIELD> <version()> <toString()>
//	keyword <name> <type> strict=<isStrictKeyword(type)>
//
// It checks every enum table (ordinals, names, Java toString overrides) and
// the contents of every static FeatureSet constant, which the parse-tree
// oracle cannot reach.

type enumValue interface {
	Name() string
	String() string
}

func enumValues[E enumValue](values []E) []enumValue {
	out := make([]enumValue, len(values))
	for i, v := range values {
		out[i] = v
	}
	return out
}

var javaEnums = map[string][]enumValue{
	"TokenType":                             enumValues(TokenType_values()),
	"Keywords":                              enumValues(Keywords_values()),
	"FeatureSet$Feature":                    enumValues(FeatureSet_Feature_values()),
	"FeatureSet$LangVersion":                enumValues(FeatureSet_LangVersion_values()),
	"TemplateLiteralToken$ErrorLevel":       enumValues(TemplateLiteralToken_ErrorLevel_values()),
	"Parser$FunctionFlavor":                 enumValues(Parser_FunctionFlavor_values()),
	"Parser$Config$Mode":                    enumValues(Parser_Config_Mode_values()),
	"Parser$Expression":                     enumValues(Parser_Expression_values()),
	"Parser$PatternKind":                    enumValues(Parser_PatternKind_values()),
	"ParseTreeType":                         enumValues(ParseTreeType_values()),
	"Comment$Type":                          enumValues(Comment_Type_values()),
	"FunctionDeclarationTree$Kind":          enumValues(FunctionDeclarationTree_Kind_values()),
	"UpdateExpressionTree$OperatorPosition": enumValues(UpdateExpressionTree_OperatorPosition_values()),
}

var javaFeatureSets = map[string]*FeatureSet{
	"BARE_MINIMUM":   FeatureSet_BARE_MINIMUM,
	"ES3":            FeatureSet_ES3,
	"ES5":            FeatureSet_ES5,
	"ES2015_MODULES": FeatureSet_ES2015_MODULES,
	"ES2015":         FeatureSet_ES2015,
	"ES2016_MODULES": FeatureSet_ES2016_MODULES,
	"ES2016":         FeatureSet_ES2016,
	"ES2017_MODULES": FeatureSet_ES2017_MODULES,
	"ES2017":         FeatureSet_ES2017,
	"ES2018_MODULES": FeatureSet_ES2018_MODULES,
	"ES2018":         FeatureSet_ES2018,
	"ES2019_MODULES": FeatureSet_ES2019_MODULES,
	"ES2019":         FeatureSet_ES2019,
	"ES2020_MODULES": FeatureSet_ES2020_MODULES,
	"ES2020":         FeatureSet_ES2020,
	"ES2021_MODULES": FeatureSet_ES2021_MODULES,
	"ES2021":         FeatureSet_ES2021,
	"ES2022_MODULES": FeatureSet_ES2022_MODULES,
	"ES2022":         FeatureSet_ES2022,
	"ES_NEXT":        FeatureSet_ES_NEXT,
	"ES_UNSTABLE":    FeatureSet_ES_UNSTABLE,
	"ES_UNSUPPORTED": FeatureSet_ES_UNSUPPORTED,
	"BROWSER_2019":   FeatureSet_BROWSER_2019,
	"BROWSER_2020":   FeatureSet_BROWSER_2020,
	"BROWSER_2021":   FeatureSet_BROWSER_2021,
	"BROWSER_2022":   FeatureSet_BROWSER_2022,
	"BROWSER_2023":   FeatureSet_BROWSER_2023,
	"BROWSER_2024":   FeatureSet_BROWSER_2024,
	"BROWSER_2025":   FeatureSet_BROWSER_2025,
	"BROWSER_2026":   FeatureSet_BROWSER_2026,
	"ALL":            FeatureSet_ALL,
}

// javaToString renders v.String() like the generator renders toString(),
// mapping a Go nil dereference to Java's NullPointerException.
func javaToString(v enumValue) (s string) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(runtime.Error); ok && strings.Contains(re.Error(), "nil pointer dereference") {
				s = "<NullPointerException>"
				return
			}
			s = fmt.Sprintf("<panic %v>", r)
		}
	}()
	return v.String()
}

func TestEnumsAndFeatureSetsMatchJava(t *testing.T) {
	f, err := os.Open("testdata/enums.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	seenEnum := map[string]int{}
	seenFS := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		switch kind, rest, _ := strings.Cut(line, " "); kind {
		case "enum":
			var class, name string
			var ordinal int
			fmt.Sscanf(rest, "%s %d %s", &class, &ordinal, &name)
			str := rest[len(fmt.Sprintf("%s %d %s ", class, ordinal, name)):]
			values, ok := javaEnums[class]
			if !ok {
				t.Errorf("no Go enum for Java %s", class)
				continue
			}
			seenEnum[class]++
			if ordinal >= len(values) {
				t.Errorf("%s: Go has no ordinal %d (%s)", class, ordinal, name)
				continue
			}
			if got := values[ordinal].Name(); got != name {
				t.Errorf("%s ordinal %d: Go name %s, Java %s", class, ordinal, got, name)
			}
			if got := javaToString(values[ordinal]); got != str {
				t.Errorf("%s.%s.toString(): Go %q, Java %q", class, name, got, str)
			}
		case "featureset":
			name, rest, _ := strings.Cut(rest, " ")
			version, str, _ := strings.Cut(rest, " ")
			fs, ok := javaFeatureSets[name]
			if !ok {
				t.Errorf("no Go FeatureSet_%s", name)
				continue
			}
			seenFS++
			if got := fs.Version(); got != version {
				t.Errorf("FeatureSet_%s.Version() = %s, Java %s", name, got, version)
			}
			if got := fs.String(); got != str {
				t.Errorf("FeatureSet_%s = %s\nJava %s", name, got, str)
			}
		case "keyword":
			var name, typ, strict string
			fmt.Sscanf(rest, "%s %s %s", &name, &typ, &strict)
			found := false
			for _, k := range Keywords_values() {
				if k.Name() == name {
					found = true
					if k.Type().Name() != typ {
						t.Errorf("Keywords.%s.type = %s, Java %s", name, k.Type().Name(), typ)
					}
					if got := fmt.Sprintf("strict=%v", Keywords_IsStrictKeyword(k.Type())); got != strict {
						t.Errorf("isStrictKeyword(%s) %s, Java %s", typ, got, strict)
					}
					if Keywords_Get(k.Value()) == nil || *Keywords_Get(k.Value()) != k {
						t.Errorf("Keywords_Get(%q) does not find %s", k.Value(), name)
					}
				}
			}
			if !found {
				t.Errorf("no Go keyword %s", name)
			}
		}
	}
	for class, values := range javaEnums {
		if seenEnum[class] != len(values) {
			t.Errorf("%s: Java has %d constants, Go %d", class, seenEnum[class], len(values))
		}
	}
	if seenFS != len(javaFeatureSets) {
		t.Errorf("Java has %d FeatureSet constants, Go table %d", seenFS, len(javaFeatureSets))
	}
}
