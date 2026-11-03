package sourcemap_test

// Differential test against the pinned Java implementation: replays the
// calls recorded by tools/oracle/java/SourceMapDiff.java (seeded random
// generator scripts, merges, index maps, mutated maps, hand-written
// consumer/parser edge cases, Base64/VLQ/escaping vectors) on the Go port
// and compares every result with Java's. See SourceMapDiff.java for the
// format.

import (
	"bufio"
	"compress/gzip"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

func diffEsc(s string) string {
	var sb strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			sb.WriteString("\\\\")
		case c == '\t':
			sb.WriteString("\\t")
		case c == '\n':
			sb.WriteString("\\n")
		case c >= 0x20 && c <= 0x7e:
			sb.WriteByte(byte(c))
		default:
			sb.WriteString(jutil.StringFormat("\\u%04x", int32(c)))
		}
	}
	return sb.String()
}

func diffUnescape(t *testing.T, s string) sourcemap.OptionalString {
	if s == "\\0" {
		return sourcemap.OptionalString{}
	}
	var units []uint16
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\\' {
			units = append(units, uint16(c))
			continue
		}
		i++
		switch s[i] {
		case '\\':
			units = append(units, '\\')
		case 't':
			units = append(units, '\t')
		case 'n':
			units = append(units, '\n')
		case 'u':
			v, err := strconv.ParseUint(s[i+1:i+5], 16, 16)
			if err != nil {
				t.Fatalf("bad escape in %q", s)
			}
			units = append(units, uint16(v))
			i += 4
		default:
			t.Fatalf("bad escape in %q", s)
		}
	}
	return sourcemap.SomeString(jutil.StringFromUTF16(units))
}

// diffExc is SourceMapDiff.exc on a recovered panic. JVM exceptions
// without a Java-defined message are reported by kind (Go runtime errors
// are their Go counterparts).
func diffExc(r any) string {
	runtimeKind := func(msg string) string {
		switch {
		case strings.Contains(msg, "nil pointer dereference"):
			return "!NPE"
		case strings.Contains(msg, "index out of range"), strings.Contains(msg, "slice bounds out of range"):
			return "!INDEX"
		}
		return "!GO " + msg
	}
	switch e := r.(type) {
	case *sourcemap.SourceMapParseException:
		return "!SMPE " + e.GetMessage()
	case jutil.INullPointerException:
		return "!NPE"
	case jutil.IIndexOutOfBoundsException:
		return "!INDEX"
	case jutil.Throwable:
		if strings.HasPrefix(e.Error(), "java.lang.RuntimeException: runtime error: ") {
			return runtimeKind(e.Error())
		}
		return "!" + e.Error()
	case runtime.Error:
		return runtimeKind(e.Error())
	}
	return "!GO " + jutil.ToJavaString(r)
}

func diffSimpleName(o any) string {
	switch o.(type) {
	case string:
		return "String"
	case int:
		return "Integer"
	case bool:
		return "Boolean"
	case *jutil.JsonObject:
		return "JsonObject"
	case *jutil.JsonArray:
		return "JsonArray"
	case *jutil.JsonPrimitive:
		return "JsonPrimitive"
	case *jutil.JsonNull:
		return "JsonNull"
	}
	return "?"
}

func diffDescribe(o any) string {
	if jutil.IsNil(o) {
		return "null"
	}
	return diffSimpleName(o) + ":" + jutil.ToJavaString(o)
}

func diffOptStr(s sourcemap.OptionalString) string {
	if !s.Valid {
		return "<null>"
	}
	return s.Value
}

func diffPairStr(s string, ok bool) string {
	if !ok {
		return "<null>"
	}
	return s
}

func diffOptStrs(a []sourcemap.OptionalString) string {
	if a == nil {
		return "<nullarray>"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, s := range a {
		if i > 0 {
			sb.WriteString("|")
		}
		sb.WriteString(diffOptStr(s))
	}
	sb.WriteString("]")
	return sb.String()
}

func mustStr(s sourcemap.OptionalString) string {
	if !s.Valid {
		panic(jutil.NewNullPointerException())
	}
	return s.Value
}

func diffOm(c *sourcemap.SourceMapConsumerV3, lineNumber, column int) string {
	m, hasIdentifier := sourcemap.LookupMappingForDiff(c, lineNumber, column)
	if m == nil {
		return "-"
	}
	id := ""
	if hasIdentifier {
		id = "=" + m.GetIdentifier()
	}
	return m.GetOriginalFile() + ":" + strconv.Itoa(m.GetLineNumber()) + ":" +
		strconv.Itoa(m.GetColumnPosition()) + id + ":" + strconv.Itoa(int(m.GetPrecision()))
}

type diffState struct {
	buffers []string
	g       *sourcemap.SourceMapGeneratorV3
	c       *sourcemap.SourceMapConsumerV3
	out     []string
}

func (d *diffState) result(s string) { d.out = append(d.out, s) }

func (d *diffState) resolve(arg sourcemap.OptionalString) string {
	s := mustStr(arg)
	if strings.HasPrefix(s, "@") {
		n, _ := strconv.Atoi(s[1:])
		return d.buffers[n]
	}
	return s
}

func (d *diffState) emit(op []sourcemap.OptionalString) {
	defer func() {
		if r := recover(); r != nil {
			d.result(diffExc(r))
		}
	}()
	d.run(op)
}

func atoi(s sourcemap.OptionalString) int {
	n, err := strconv.Atoi(mustStr(s))
	if err != nil {
		panic(err)
	}
	return n
}

type diffMergeAction struct{}

func (diffMergeAction) Merge(key string, cur, nw any) any {
	return "merged(" + key + "," + diffDescribe(cur) + "," + diffDescribe(nw) + ")"
}

func (d *diffState) run(op []sourcemap.OptionalString) {
	switch mustStr(op[0]) {
	case "case":
		d.buffers = nil
	case "gnew":
		d.g = sourcemap.NewSourceMapGeneratorV3()
	case "gfactory":
		format := sourcemap.SourceMapFormat_V3
		if mustStr(op[1]) == "DEFAULT" {
			format = sourcemap.SourceMapFormat_DEFAULT
		}
		d.g = sourcemap.SourceMapGeneratorFactory_GetInstance(format)
	case "greset":
		d.g.Reset()
	case "gvalidate":
		d.g.Validate(mustStr(op[1]) == "true")
	case "gsourceroot":
		d.g.SetSourceRoot(mustStr(op[1]))
	case "gprefix":
		d.g.SetWrapperPrefix(mustStr(op[1]))
	case "gstart":
		d.g.SetStartingPosition(atoi(op[1]), atoi(op[2]))
	case "gmap":
		origPos := sourcemap.NewFilePosition(atoi(op[3]), atoi(op[4]))
		startPos := sourcemap.NewFilePosition(atoi(op[5]), atoi(op[6]))
		endPos := sourcemap.NewFilePosition(atoi(op[7]), atoi(op[8]))
		if op[1].Valid {
			d.g.AddMappingPositions(op[1].Value, op[2], origPos, startPos, endPos)
		}
	case "gcontent":
		d.g.AddSourcesContent(mustStr(op[1]), op[2])
	case "gext":
		var value any
		switch mustStr(op[2]) {
		case "str":
			value = mustStr(op[3])
		case "int":
			value = atoi(op[3])
		case "bool":
			value = strings.EqualFold(mustStr(op[3]), "true")
		case "json":
			if e := jutil.NewGson().FromJsonElement(mustStr(op[3])); e != nil {
				value = e
			}
		}
		d.g.AddExtension(mustStr(op[1]), value)
	case "gextrm":
		d.g.RemoveExtension(mustStr(op[1]))
	case "gexthas":
		d.result(jutil.ToJavaString(d.g.HasExtension(mustStr(op[1]))))
	case "gextget":
		d.result(diffDescribe(d.g.GetExtension(mustStr(op[1]))))
	case "gappend":
		var sb strings.Builder
		defer func() {
			d.buffers = append(d.buffers, sb.String())
			d.result(sb.String())
		}()
		d.g.AppendToOptional(&sb, op[1])
	case "gindex":
		var sections []*sourcemap.SourceMapSection
		n := atoi(op[2])
		for i := 0; i < n; i++ {
			type_ := mustStr(op[3+4*i])
			value := d.resolve(op[4+4*i])
			line := atoi(op[5+4*i])
			column := atoi(op[6+4*i])
			if type_ == "URL" {
				sections = append(sections, sourcemap.SourceMapSection_ForURL(value, line, column))
			} else {
				sections = append(sections, sourcemap.SourceMapSection_ForMap(value, line, column))
			}
		}
		var sb strings.Builder
		defer func() {
			d.buffers = append(d.buffers, sb.String())
			d.result(sb.String())
		}()
		d.g.AppendIndexMapTo(&sb, mustStr(op[1]), sections)
	case "gmerge":
		if err := d.g.MergeMapSection(atoi(op[1]), atoi(op[2]), d.resolve(op[3])); err != nil {
			d.result(diffExc(err))
			return
		}
		d.result("ok")
	case "gmergex":
		if err := d.g.MergeMapSectionWithMergeAction(atoi(op[1]), atoi(op[2]), d.resolve(op[3]), diffMergeAction{}); err != nil {
			d.result(diffExc(err))
			return
		}
		d.result("ok")
	case "cparse":
		d.c = sourcemap.NewSourceMapConsumerV3()
		if err := d.c.Parse(d.resolve(op[1])); err != nil {
			d.result(diffExc(err))
			return
		}
		d.result("ok")
	case "cfactory":
		maps := make(map[string]string, len(op)-2)
		for _, kv := range op[2:] {
			kvs := mustStr(kv)
			eq := strings.IndexByte(kvs, '=')
			v := kvs[eq+1:]
			maps[kvs[:eq]] = d.resolve(sourcemap.SomeString(v))
		}
		m, err := sourcemap.ParseSourceMapWithSupplier(d.resolve(op[1]), mapSupplierFunc(func(url string) (string, bool) {
			v, ok := maps[url]
			if !ok {
				return "", false
			}
			if v == "!io" {
				panic(jutil.NewIOException("cannot read " + url))
			}
			return v, true
		}))
		if err != nil {
			d.result(diffExc(err))
			return
		}
		d.c = m
		d.result("ok")
	case "cfactory0":
		m, err := sourcemap.ParseSourceMap(d.resolve(op[1]))
		if err != nil {
			d.result(diffExc(err))
			return
		}
		d.c = m
		d.result("ok")
	case "gbig":
		var sb strings.Builder
		sb.WriteString("{\"version\":3,\"mappings\":\"")
		sb.WriteString(mustStr(op[3]))
		sb.WriteString("\",\"sources\":[")
		for i := 0; i < atoi(op[1]); i++ {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString("\"s" + strconv.Itoa(i) + "\"")
		}
		sb.WriteString("],\"names\":[")
		for i := 0; i < atoi(op[2]); i++ {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString("\"n" + strconv.Itoa(i) + "\"")
		}
		sb.WriteString("]}")
		d.buffers = append(d.buffers, sb.String())
		d.result(strconv.Itoa(jutil.StringLength(sb.String())))
	case "clookup":
		for i := 1; i+1 < len(op); i += 2 {
			d.result(diffOm(d.c, atoi(op[i]), atoi(op[i+1])))
		}
	case "creverse":
		var sb strings.Builder
		for _, m := range d.c.GetReverseMapping(mustStr(op[1]), atoi(op[2]), 0) {
			sb.WriteString(" " + strconv.Itoa(m.GetLineNumber()) + "," + strconv.Itoa(m.GetColumnPosition()))
		}
		d.result(sb.String())
	case "cvisitcount":
		n := 0
		var sb strings.Builder
		d.c.VisitMappings(entryVisitorFunc(func(sourceName, symbolName sourcemap.OptionalString, src, start, end sourcemap.FilePosition) {
			if n < 6 {
				sb.WriteString(diffOptStr(sourceName) + "|" + diffOptStr(symbolName) + " ")
			}
			n++
		}))
		d.result(strconv.Itoa(n) + " " + sb.String())
	case "cdump":
		d.dumpConsumer()
	case "oparse":
		obj, err := sourcemap.ParseSourceMapObject(d.resolve(op[1]))
		if err != nil {
			d.result(diffExc(err))
			return
		}
		d.dumpObject(obj)
	case "venc":
		var sb strings.Builder
		sourcemap.Base64VLQ_Encode(&sb, int32(atoi(op[1])))
		d.result(sb.String())
	case "vdec":
		it := &base64VLQTest_CharIteratorImpl{}
		it.set(mustStr(op[1]))
		var sb strings.Builder
		defer func() { d.result(sb.String()) }()
		for it.HasNext() {
			sb.WriteString(jutil.ToJavaString(sourcemap.Base64VLQ_Decode(it)))
			sb.WriteString(" ")
		}
	case "b64":
		d.result(sourcemap.Base64_Base64EncodeInt(int32(atoi(op[1]))))
	case "b64from":
		d.result(jutil.ToJavaString(sourcemap.Base64_FromBase64(jutil.StringToUTF16(mustStr(op[1]))[0])))
	case "esc":
		d.result(sourcemap.Util_EscapeString(mustStr(op[1])))
	default:
		panic("unknown op " + mustStr(op[0]))
	}
}

func (d *diffState) guarded(body func()) {
	defer func() {
		if r := recover(); r != nil {
			d.result(diffExc(r))
		}
	}()
	body()
}

func (d *diffState) dumpObject(o *sourcemap.SourceMapObject) {
	d.result("version=" + strconv.Itoa(o.GetVersion()) + " lineCount=" + strconv.Itoa(o.GetLineCount()))
	d.result("file=" + diffPairStr(o.GetFile()))
	d.result("mappings=" + diffPairStr(o.GetMappings()))
	d.result("sourceRoot=" + diffPairStr(o.GetSourceRoot()))
	d.result("sources=" + diffOptStrs(o.GetSources()))
	d.result("sourcesContent=" + diffOptStrs(o.GetSourcesContent()))
	d.result("names=" + diffOptStrs(o.GetNames()))
	if o.GetSections() == nil {
		d.result("sections=null")
	} else {
		for _, s := range o.GetSections() {
			d.result("section " + s.GetSectionType().String() + " " + strconv.Itoa(s.GetLine()) + " " +
				strconv.Itoa(s.GetColumn()) + " " + diffPairStr(s.GetSectionValue()))
		}
	}
	for _, ext := range o.GetExtensions() {
		d.result("ext " + ext.Key + "=" + diffDescribe(ext.Value))
	}
}

func (d *diffState) dumpConsumer() {
	c := d.c
	d.guarded(func() {
		d.result("file=" + diffPairStr(c.GetFile()) + " lineCount=" + strconv.Itoa(c.GetLineCount()) +
			" sourceRoot=" + diffPairStr(c.GetSourceRoot()))
	})
	d.guarded(func() { d.result("sources=" + diffOptStrs(c.OriginalSources())) })
	d.guarded(func() { d.result("sourcesContent=" + diffOptStrs(c.OriginalSourcesContent())) })
	d.guarded(func() { d.result("names=" + diffOptStrs(c.GetOriginalNames())) })
	d.guarded(func() {
		for _, ext := range c.GetExtensions() {
			d.result("ext " + ext.Key + "=" + diffDescribe(ext.Value))
		}
	})
	max0, max1, max2 := 0, 0, 0 // line, column, source line
	d.guarded(func() {
		c.VisitMappings(entryVisitorFunc(func(sourceName, symbolName sourcemap.OptionalString, src, start, end sourcemap.FilePosition) {
			d.result("visit " + diffOptStr(sourceName) + "|" + diffOptStr(symbolName) + "|" +
				strconv.Itoa(src.GetLine()) + "," + strconv.Itoa(src.GetColumn()) + "|" +
				strconv.Itoa(start.GetLine()) + "," + strconv.Itoa(start.GetColumn()) + "|" +
				strconv.Itoa(end.GetLine()) + "," + strconv.Itoa(end.GetColumn()))
			max0 = max(max0, max(start.GetLine(), end.GetLine()))
			max1 = max(max1, max(start.GetColumn(), end.GetColumn()))
			max2 = max(max2, src.GetLine())
		}))
	})
	maxLine := min(jutil.IntAdd(max(c.GetLineCount(), max0), 2), 300)
	maxCol := min(jutil.IntAdd(max1, 3), 300)
	for line := 1; line <= maxLine; line++ {
		l := line
		d.guarded(func() {
			var sb strings.Builder
			sb.WriteString("L" + strconv.Itoa(l) + ":")
			var prev sourcemap.OptionalString
			runStart := 0
			for col := 0; col <= jutil.IntAdd(maxCol, 1); col++ {
				var cur sourcemap.OptionalString
				if col <= maxCol {
					cur = sourcemap.SomeString(diffOm(c, l, col))
				}
				if prev.Valid && prev != cur {
					sb.WriteString(" " + strconv.Itoa(runStart) + "-" + strconv.Itoa(col-1) + "=" + prev.Value)
					runStart = col
				}
				if !prev.Valid {
					runStart = col
				}
				prev = cur
			}
			d.result(sb.String())
		})
	}
	probes := [][2]int{
		{0, 1}, {-1, 1}, {1, -5}, {1, math.MaxInt32}, {1, math.MinInt32},
		{math.MaxInt32, 1}, {math.MinInt32, 1}, {2, math.MinInt32},
	}
	for _, p := range probes {
		d.guarded(func() {
			d.result("P" + strconv.Itoa(p[0]) + "," + strconv.Itoa(p[1]) + "=" + diffOm(c, p[0], p[1]))
		})
	}
	d.guarded(func() {
		var files []string
		seen := make(map[string]bool)
		addFile := func(s string) {
			if !seen[s] {
				seen[s] = true
				files = append(files, s)
			}
		}
		for _, s := range c.OriginalSources() {
			if s.Valid {
				addFile(s.Value)
			}
		}
		addFile("not-a-source.js")
		for _, f := range files {
			for line := -1; line <= min(jutil.IntAdd(max2, 1), 200); line++ {
				rs := c.GetReverseMapping(f, line, 0)
				if len(rs) > 0 {
					var sb strings.Builder
					sb.WriteString("R " + f + " " + strconv.Itoa(line) + ":")
					for _, m := range rs {
						sb.WriteString(" " + strconv.Itoa(m.GetLineNumber()) + "," + strconv.Itoa(m.GetColumnPosition()))
					}
					d.result(sb.String())
				}
			}
		}
	})
}

func TestDifferentialAgainstJava(t *testing.T) {
	f, err := os.Open("testdata/differential.txt.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	type call struct {
		line int
		op   []string
		want []string
	}
	var calls []*call
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 1<<24), 1<<24)
	n := 0
	for sc.Scan() {
		n++
		text := sc.Text()
		fields := strings.Split(text, "\t")
		switch fields[0] {
		case ">":
			calls = append(calls, &call{line: n, op: fields[1:]})
		case "<":
			calls[len(calls)-1].want = append(calls[len(calls)-1].want, fields[1])
		default:
			t.Fatalf("line %d: %q", n, text)
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	if len(calls) < 10000 {
		t.Fatalf("only %d calls", len(calls))
	}

	d := &diffState{g: sourcemap.NewSourceMapGeneratorV3(), c: sourcemap.NewSourceMapConsumerV3()}
	diffs, outputs := 0, 0
	skipFastCdump := false
	for _, c := range calls {
		op := make([]sourcemap.OptionalString, len(c.op))
		for i, a := range c.op {
			op[i] = diffUnescape(t, a)
		}
		switch mustStr(op[0]) {
		case "ofast", "proto":
			continue
		case "cfast":
			if len(c.want) == 1 && c.want[0] == "ok" {
				cand := sourcemap.NewSourceMapConsumerV3()
				if cand.TryParse(d.resolve(op[1])) == nil {
					d.c = cand
					skipFastCdump = false
				} else {
					skipFastCdump = true
				}
			}
			continue
		case "case", "cparse":
			skipFastCdump = false
		case "cdump":
			if skipFastCdump {
				continue
			}
		}
		d.out = nil
		d.emit(op)
		got := make([]string, len(d.out))
		for i, s := range d.out {
			got[i] = diffEsc(s)
		}
		outputs += len(c.want)
		if strings.Join(got, "\n") != strings.Join(c.want, "\n") {
			diffs++
			if diffs <= 5 {
				for i := 0; i < len(got) || i < len(c.want); i++ {
					var g, w string
					if i < len(got) {
						g = got[i]
					}
					if i < len(c.want) {
						w = c.want[i]
					}
					if g != w {
						t.Errorf("line %d %s (result %d):\n java: %.600s\n   go: %.600s", c.line, strings.Join(c.op, " "), i, w, g)
						break
					}
				}
			}
		}
	}
	t.Logf("%d calls, %d results compared, %d calls differ", len(calls), outputs, diffs)
	if diffs > 0 {
		t.Fatalf("%d of %d calls differ from Java", diffs, len(calls))
	}
}
