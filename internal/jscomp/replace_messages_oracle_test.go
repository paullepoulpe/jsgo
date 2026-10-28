package jscomp

// Go-only: test-helper — the differential oracle for ReplaceMessages with a
// bundle built from each input: the Go mirror of
// tools/oracle/java/ReplaceMessagesDump.java, whose header documents the
// bundle, the variants (full, full-strict, protect, complete) and the output
// format.
//
//   - TestReplaceMessagesDumpGolden compares against the checked-in Java
//     output testdata/replace_messages/golden.txt.gz for the inputs listed in
//     testdata/replace_messages/files.txt (sources in
//     testdata/replace_messages/src; always runs). Regenerate with
//       java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:<build> \
//         com.google.javascript.jscomp.ReplaceMessagesDump --externs=tools/oracle/typecheck-externs.js \
//         --full-every=2 internal/jscomp/testdata/replace_messages/golden.txt.gz \
//         internal/jscomp/testdata/replace_messages/src \
//         --list=internal/jscomp/testdata/replace_messages/files.txt
//   - TestReplaceMessagesDumpOracle compares any other run when
//     REPLACE_MESSAGES_ORACLE (the Java <out.gz>), REPLACE_MESSAGES_ROOT and
//     REPLACE_MESSAGES_FILES (the list given with --list) are set, with
//     REPLACE_MESSAGES_SAMPLE and REPLACE_MESSAGES_FULL_EVERY matching the
//     Java run's --sample and --full-every (REPLACE_MESSAGES_FULL=1 for
//     --full); REPLACE_MESSAGES_DIFF_DIR, if set, receives <rel>.go.txt and
//     <rel>.java.txt for every mismatching file. Use absolute paths.

import (
	"iter"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

var replaceMessagesDump_VARIANTS = []string{"full", "full-strict", "protect", "complete"}

// replaceMessagesDump_collector is ReplaceMessagesDump.Collector: it records
// every message the visitor finds, in visiting order.
type replaceMessagesDump_collector struct {
	JsMessageVisitor
	found []*JsMessage
}

func newReplaceMessagesDump_collector(compiler AbstractCompiler) *replaceMessagesDump_collector {
	c := &replaceMessagesDump_collector{}
	c.InitJsMessageVisitor(
		compiler,
		nil,
		func(definition JsMessageDefinition) { c.found = append(c.found, definition.GetMessage()) },
		func(definition IcuTemplateDefinition) { c.found = append(c.found, definition.GetMessage()) },
		nil,
	)
	return c
}

// replaceMessagesDump_bundle is ReplaceMessagesDump.Bundle.
type replaceMessagesDump_bundle struct {
	messages *jutil.LinkedHashMap[string, *JsMessage]
}

func (b *replaceMessagesDump_bundle) GetMessage(id string) *JsMessage {
	return b.messages.Get(id)
}

func (b *replaceMessagesDump_bundle) GetAllMessages() iter.Seq[*JsMessage] {
	return b.messages.Values()
}

func (b *replaceMessagesDump_bundle) IdGenerator() JsMessage_IdGenerator {
	return nil
}

// replaceMessagesDump_buildBundle is ReplaceMessagesDump.buildBundle.
func replaceMessagesDump_buildBundle(found []*JsMessage) *replaceMessagesDump_bundle {
	byId := jutil.NewLinkedHashMap[string, *JsMessage]()
	for _, m := range found {
		byId.PutIfAbsent(m.GetId(), m)
	}
	bundle := &replaceMessagesDump_bundle{messages: jutil.NewLinkedHashMap[string, *JsMessage]()}
	i := 0
	for m := range byId.Values() {
		skip := i%4 == 3
		i++
		if skip {
			continue
		}
		b := NewJsMessage_Builder().SetKey(m.GetKey()).SetId(m.GetId())
		parts := m.GetParts()
		for j := len(parts) - 1; j >= 0; j-- {
			p := parts[j]
			if p.IsPlaceholder() {
				b.AppendPart(p)
			} else {
				b.AppendStringPart("~" + p.GetString())
			}
		}
		bundle.messages.Put(m.GetId(), b.Build())
	}
	return bundle
}

// replaceMessagesDump_dumpFile is ReplaceMessagesDump.dumpFile.
func replaceMessagesDump_dumpFile(sb *strings.Builder, externs, rel, src string, fullThis bool) {
	typeCheckDump_line(sb, "## file "+rel)
	externsFile := SourceFile_FromCode("externs.js", externs)
	inputFile := SourceFile_FromCode(rel, src)
	var bundle *replaceMessagesDump_bundle
	{
		s := newSerializationDump_setup(externsFile, inputFile)
		var thrown jutil.Throwable
		jutil.Try(func() { s.compiler.Parse() }, func(e jutil.IRuntimeException) { thrown = e })
		if thrown != nil {
			typeCheckDump_line(sb, "# parse exception "+typeCheckDump_className(thrown)+": "+
				jutil.ToJavaString(jutil.GetMessageOrNull(thrown)))
			return
		}
		if s.errors.GetErrorCount() == 0 {
			typeCheckDump_line(sb, "# parse ok")
		} else {
			typeCheckDump_line(sb, "# parse errors")
			for _, d := range s.errors.lines {
				typeCheckDump_line(sb, "diag "+d)
			}
			return
		}
		jutil.Try(func() {
			collector := newReplaceMessagesDump_collector(s.compiler)
			collector.Process(s.compiler.GetExternsRoot(), s.compiler.GetJsRoot())
			bundle = replaceMessagesDump_buildBundle(collector.found)
		}, func(e jutil.IRuntimeException) { thrown = e })
		if thrown != nil {
			typeCheckDump_line(sb, "exception bundle "+typeCheckDump_className(thrown)+": "+
				jutil.ToJavaString(jutil.GetMessageOrNull(thrown)))
			return
		}
	}
	var bundleLines []string
	for m := range bundle.messages.Values() {
		bundleLines = append(bundleLines, m.GetId())
		for _, p := range m.GetParts() {
			if p.IsPlaceholder() {
				bundleLines = append(bundleLines, "  p "+p.GetCanonicalPlaceholderName())
			} else {
				bundleLines = append(bundleLines, "  s "+p.GetString())
			}
		}
	}
	serializationDump_section(sb, "bundle", bundleLines, fullThis)
	for _, variant := range replaceMessagesDump_VARIANTS {
		typeCheckDump_line(sb, "## variant "+variant)
		replaceMessagesDump_dumpVariant(sb, variant, externsFile, inputFile, bundle, fullThis)
	}
}

// replaceMessagesDump_dumpVariant is ReplaceMessagesDump.dumpVariant.
func replaceMessagesDump_dumpVariant(
	sb *strings.Builder, variant string, externsFile, inputFile *SourceFile,
	bundle *replaceMessagesDump_bundle, fullThis bool) {
	s := newSerializationDump_setup(externsFile, inputFile)
	stage := "checks"
	jutil.Try(func() {
		s.compiler.Parse()
		externsRoot := s.compiler.GetExternsRoot()
		jsRoot := s.compiler.GetJsRoot()
		serializationDump_typeCheck(s.compiler, s.options)
		NewRemoveCastNodes(s.compiler).Process(externsRoot, jsRoot)
		NewConvertTypesToColors(
			s.compiler, SerializationOptions_Builder_().SetIncludeDebugInfo(false).Build()).
			Process(externsRoot, jsRoot)
		s.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC)
		before := len(s.errors.lines)
		typeCheckDump_line(sb, "checks "+strconv.Itoa(before)+" "+
			serializationDump_sha16(serializationDump_prefixed(s.errors.lines)))

		stage = "pass"
		switch variant {
		case "full":
			NewReplaceMessages(s.compiler, bundle, false).GetFullReplacementPass().Process(externsRoot, jsRoot)
		case "full-strict":
			NewReplaceMessages(s.compiler, bundle, true).GetFullReplacementPass().Process(externsRoot, jsRoot)
		case "protect":
			NewReplaceMessages(s.compiler, bundle, false).GetMsgProtectionPass().Process(externsRoot, jsRoot)
		case "complete":
			NewReplaceMessages(s.compiler, bundle, false).GetMsgProtectionPass().Process(externsRoot, jsRoot)
			NewReplaceMessages(s.compiler, bundle, false).GetReplacementCompletionPass().Process(externsRoot, jsRoot)
		default:
			panic(jutil.NewIllegalArgumentExceptionMsg(variant))
		}
		serializationDump_section(sb, "diags", serializationDump_prefixed(s.errors.lines[before:]), fullThis)
		serializationDump_section(sb, "tree", renamingDump_dumpTree(nil, jsRoot, 0), fullThis)
		serializationDump_section(sb, "externs", renamingDump_dumpTree(nil, externsRoot, 0), fullThis)
	}, func(e jutil.IRuntimeException) {
		typeCheckDump_line(sb, "exception "+stage+" "+typeCheckDump_className(e)+": "+
			jutil.ToJavaString(jutil.GetMessageOrNull(e)))
	})
}

func TestReplaceMessagesDumpGolden(t *testing.T) {
	dir := filepath.Join("testdata", "replace_messages")
	files := typeCheckDump_readList(t, filepath.Join(dir, "files.txt"))
	oracleDump_compare(t, replaceMessagesDump_dumpFile, typeCheckDump_externs(t), filepath.Join(dir, "src"),
		filepath.Join(dir, "golden.txt.gz"), files, 1, 2, false, os.Getenv("REPLACE_MESSAGES_DIFF_DIR"))
}

func TestReplaceMessagesDumpOracle(t *testing.T) {
	oracle := os.Getenv("REPLACE_MESSAGES_ORACLE")
	root := os.Getenv("REPLACE_MESSAGES_ROOT")
	list := os.Getenv("REPLACE_MESSAGES_FILES")
	if oracle == "" || root == "" || list == "" {
		t.Skip("set REPLACE_MESSAGES_ORACLE, REPLACE_MESSAGES_ROOT and REPLACE_MESSAGES_FILES to run the corpus oracle")
	}
	atoi := func(name string, def int) int {
		s := os.Getenv(name)
		if s == "" {
			return def
		}
		n, err := strconv.Atoi(s)
		if err != nil || n < 0 {
			t.Fatalf("%s=%q: want a non-negative integer", name, s)
		}
		return n
	}
	sample := atoi("REPLACE_MESSAGES_SAMPLE", 1)
	if sample < 1 {
		sample = 1
	}
	files := typeCheckDump_readList(t, list)
	oracleDump_compare(t, replaceMessagesDump_dumpFile, typeCheckDump_externs(t), root, oracle, files, sample,
		atoi("REPLACE_MESSAGES_FULL_EVERY", 0), os.Getenv("REPLACE_MESSAGES_FULL") == "1",
		os.Getenv("REPLACE_MESSAGES_DIFF_DIR"))
}
