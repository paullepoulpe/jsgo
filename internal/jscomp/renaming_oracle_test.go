package jscomp

// Go-only: test-helper — the differential oracle for the property renaming
// passes (DisambiguateProperties, AmbiguateProperties, RenameProperties,
// RescopeGlobalSymbols): the Go mirror of tools/oracle/java/RenamingDump.java,
// whose header documents the variants and the output format.
//
//   - TestRenamingDumpGolden compares against the checked-in Java output
//     testdata/renaming/golden.txt.gz for the inputs listed in
//     testdata/renaming/files.txt (sources in testdata/renaming/src; always
//     runs). Regenerate with
//       java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:<build> \
//         com.google.javascript.jscomp.RenamingDump --externs=tools/oracle/typecheck-externs.js \
//         --full-every=2 internal/jscomp/testdata/renaming/golden.txt.gz \
//         internal/jscomp/testdata/renaming/src --list=internal/jscomp/testdata/renaming/files.txt
//   - TestRenamingDumpOracle compares any other run when RENAMING_ORACLE (the
//     Java <out.gz>), RENAMING_ROOT and RENAMING_FILES (the list given with
//     --list) are set, with RENAMING_SAMPLE and RENAMING_FULL_EVERY matching
//     the Java run's --sample and --full-every (RENAMING_FULL=1 for --full);
//     RENAMING_DIFF_DIR, if set, receives <rel>.go.txt and <rel>.java.txt for
//     every mismatching file. Use absolute paths.

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var renamingDump_VARIANTS = []string{
	"disambiguate", "ambiguate", "rename", "rename-pseudo", "rescope", "rescope-extern",
}

// renamingDump_dumpFile is RenamingDump.dumpFile.
func renamingDump_dumpFile(sb *strings.Builder, externs, rel, src string, fullThis bool) {
	typeCheckDump_line(sb, "## file "+rel)
	externsFile := SourceFile_FromCode("externs.js", externs)
	inputFile := SourceFile_FromCode(rel, src)
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
	}
	for _, variant := range renamingDump_VARIANTS {
		typeCheckDump_line(sb, "## variant "+variant)
		renamingDump_dumpVariant(sb, variant, externsFile, inputFile, fullThis)
	}
}

// renamingDump_charSet passes a nullable option set on as Java does (a null
// Set is a nil interface).
func renamingDump_charSet(s *jutil.LinkedHashSet[uint16]) jutil.Set[uint16] {
	if s == nil {
		return nil
	}
	return s
}

// renamingDump_dumpVariant is RenamingDump.dumpVariant.
func renamingDump_dumpVariant(sb *strings.Builder, variant string, externsFile, inputFile *SourceFile, fullThis bool) {
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
		Normalize_CreateNormalizeForOptimizations(s.compiler).Process(externsRoot, jsRoot)
		GatherGetterAndSetterProperties_Update(s.compiler, externsRoot, jsRoot)
		NewGatherExternProperties(s.compiler, GatherExternProperties_Mode_OPTIMIZE).
			Process(externsRoot, jsRoot)
		before := len(s.errors.lines)
		typeCheckDump_line(sb, "checks "+strconv.Itoa(before)+" "+
			serializationDump_sha16(serializationDump_prefixed(s.errors.lines)))

		stage = "pass"
		options := s.options
		var renaming []string
		hasMap := false
		switch variant {
		case "disambiguate":
			NewDisambiguateProperties(s.compiler, options.GetPropertiesThatMustDisambiguate()).
				Process(externsRoot, jsRoot)
		case "ambiguate":
			// RenamingDump sets the renamingMap field as makePassForTesting does.
			ap := AmbiguateProperties_MakePassForTesting(
				s.compiler,
				renamingDump_charSet(options.GetPropertyReservedNamingFirstChars()),
				renamingDump_charSet(options.GetPropertyReservedNamingNonFirstChars()),
				s.compiler.GetExternProperties())
			ap.Process(externsRoot, jsRoot)
			hasMap = true
			for k, v := range ap.GetRenamingMap().All() {
				renaming = append(renaming, k+" "+v)
			}
		case "rename", "rename-pseudo":
			rprop := NewRenameProperties(
				s.compiler,
				variant == "rename-pseudo",
				nil,
				renamingDump_charSet(options.GetPropertyReservedNamingFirstChars()),
				renamingDump_charSet(options.GetPropertyReservedNamingNonFirstChars()),
				options.GetNameGenerator())
			rprop.Process(externsRoot, jsRoot)
			hasMap = true
			for k, v := range rprop.GetPropertyMap().GetOriginalNameToNewNameMap().All() {
				renaming = append(renaming, k+" "+v)
			}
		case "rescope":
			NewRescopeGlobalSymbolsWithAddExtern(s.compiler, "_", false, true).Process(externsRoot, jsRoot)
		case "rescope-extern":
			NewRescopeGlobalSymbolsWithAddExtern(s.compiler, "_", true, false).Process(externsRoot, jsRoot)
		default:
			panic(jutil.NewIllegalArgumentExceptionMsg(variant))
		}
		serializationDump_section(sb, "diags", serializationDump_prefixed(s.errors.lines[before:]), fullThis)
		serializationDump_section(sb, "tree", renamingDump_dumpTree(nil, jsRoot, 0), fullThis)
		serializationDump_section(sb, "externs", renamingDump_dumpTree(nil, externsRoot, 0), fullThis)
		if hasMap {
			serializationDump_section(sb, "map", renaming, fullThis)
		}
	}, func(e jutil.IRuntimeException) {
		typeCheckDump_line(sb, "exception "+stage+" "+typeCheckDump_className(e)+": "+
			jutil.ToJavaString(jutil.GetMessageOrNull(e)))
	})
}

// renamingDump_dumpTree is RenamingDump.dumpTree.
func renamingDump_dumpTree(out []string, n *rhino.Node, level int) []string {
	out = append(out, strings.Repeat("  ", level)+n.ToStringWithFlags(true, true, false))
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		out = renamingDump_dumpTree(out, c, level+1)
	}
	return out
}

func TestRenamingDumpGolden(t *testing.T) {
	dir := filepath.Join("testdata", "renaming")
	files := typeCheckDump_readList(t, filepath.Join(dir, "files.txt"))
	oracleDump_compare(t, renamingDump_dumpFile, typeCheckDump_externs(t), filepath.Join(dir, "src"),
		filepath.Join(dir, "golden.txt.gz"), files, 1, 2, false, os.Getenv("RENAMING_DIFF_DIR"))
}

func TestRenamingDumpOracle(t *testing.T) {
	oracle := os.Getenv("RENAMING_ORACLE")
	root := os.Getenv("RENAMING_ROOT")
	list := os.Getenv("RENAMING_FILES")
	if oracle == "" || root == "" || list == "" {
		t.Skip("set RENAMING_ORACLE, RENAMING_ROOT and RENAMING_FILES to run the corpus oracle")
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
	sample := atoi("RENAMING_SAMPLE", 1)
	if sample < 1 {
		sample = 1
	}
	files := typeCheckDump_readList(t, list)
	oracleDump_compare(t, renamingDump_dumpFile, typeCheckDump_externs(t), root, oracle, files, sample,
		atoi("RENAMING_FULL_EVERY", 0), os.Getenv("RENAMING_FULL") == "1", os.Getenv("RENAMING_DIFF_DIR"))
}
