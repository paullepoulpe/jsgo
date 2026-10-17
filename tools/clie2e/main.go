// Command clie2e is the end-to-end oracle for the jsgo binary: it runs
// `java -jar closure-compiler.jar` (in one JVM, tools/oracle/java/CliE2E)
// and the jsgo binary with identical argument vectors and compares stdout,
// stderr, the exit status and every output file byte for byte.
//
// Four suites:
//
//	fixtures  testdata/fixtures, whose Java outputs are recorded in
//	          testdata/oracle and testdata/oracle-strict (CliBatch format:
//	          stdout and stderr combined, then ---EXIT--- and the status).
//	          No JVM is needed.
//	matrix    corpus files × compilation levels, --language_out variants,
//	          --create_source_map, --formatting PRETTY_PRINT and
//	          --warning_level, plus --help, --version and bad-flag jobs.
//	sweep     every file of a corpus list compiled alone with SIMPLE,
//	          ADVANCED and --checks_only --jscomp_warning=checkTypes
//	          (output on stdout), each job classified as identical, diff,
//	          jsgo panic, jsgo hang (-timeout) or Java error; see sweep.go.
//	          -configs @flags runs the flag configurations instead
//	          (sweep_flags.go), usually on a stratified -sample with
//	          -per-config files each (sample.go).
//	projects  multi-file, project-level builds defined as data in
//	          tools/clie2e/projects/*.json (the tools/bench tasks, Closure
//	          Library apps, ES modules and CommonJS packages of the corpora,
//	          chunked output, output options), every output file compared
//	          too; classified like the sweep; see projects.go.
//
// Both sides run in the C/POSIX locale without JAVA_TOOL_OPTIONS (which
// makes the JVM print a banner on stderr). In stderr, Java stack-frame
// lines ("\tat ...", "\t... n more") are dropped before comparing: Go has
// no Java frames (jutil.Throwable_PrintStackTrace).
//
// Usage:
//
//	go run ./tools/clie2e -suite fixtures [-jsgo BIN] [-n N]
//	go run ./tools/clie2e -suite matrix -work DIR [-n N] [-seed S] [-shards K] [-jobs J]
//	go run ./tools/clie2e -suite sweep -work DIR -corpus ROOT -list LIST [-configs C,...|@flags] [-sample N] [-per-config K] [-max-size B] [-timeout D] [-reuse-java] [-shards K] [-jobs J]
//	go run ./tools/clie2e -suite projects -work DIR -javabuild DIR [-set dev|holdout|all] [-projects P[/JOB],...] [-timeout D] [-reuse-java|-skip-java|-skip-go|-classify-only] [-jobs J] [-v]
//
// A job whose jsgo run reached a not-yet-ported class (stderr contains
// "pending port:") is counted as PENDING, not as a diff.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	oraclecorpus "github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

var (
	suite    = flag.String("suite", "matrix", "fixtures, matrix, sweep or projects")
	jsgoBin  = flag.String("jsgo", "", "jsgo binary (default: build ./cmd/jsgo into the work dir)")
	work     = flag.String("work", "", "work directory (matrix suite)")
	n        = flag.Int("n", 0, "number of inputs to use (0: all fixtures / 20 corpus files)")
	seed     = flag.Int64("seed", 1, "corpus sample seed")
	shards   = flag.Int("shards", 2, "number of JVMs")
	jobsN    = flag.Int("jobs", 3, "concurrent jsgo processes")
	javaBld  = flag.String("javabuild", "/home/user/jsgo/.cache/oracle-build", "directory holding the compiled CliE2E (tools/oracle/build.sh)")
	jar      = flag.String("jar", "/home/user/jsgo/.cache/closure-compiler.jar", "pinned closure-compiler.jar")
	corpus   = flag.String("corpus", "/home/user/jsgo/.cache/corpus", "corpus root")
	list     = flag.String("list", "/home/user/jsgo/.cache/corpus-files.txt", "corpus file list")
	keep     = flag.Bool("v", false, "print the first lines of each diff")
	onlyConf = flag.String("configs", "", "comma-separated config names to run (default: all)")
)

// config is one column of the matrix: flags added before --js.
type config struct {
	name    string
	args    []string
	srcMap  bool
	perFile bool // false: the job does not take corpus inputs
	optIn   bool // sweep: only when -configs names it (or its group)
	// sweep only:
	group string // -configs @group selects every config of the group
	// variants, when set, are rotated over the config's files: the i-th
	// file of the config gets args, then variants[i%len(variants)].
	variants [][]string
	// files, when smaller, overrides -per-config for this config (a config
	// whose flags fail before any input is read needs few files).
	files int
}

var configs = []config{
	{name: "whitespace", args: []string{"--compilation_level", "WHITESPACE_ONLY"}, perFile: true},
	{name: "whitespace-pretty", args: []string{"--compilation_level", "WHITESPACE_ONLY", "--formatting", "PRETTY_PRINT"}, perFile: true},
	{name: "whitespace-es5", args: []string{"--compilation_level", "WHITESPACE_ONLY", "--language_out", "ECMASCRIPT5"}, perFile: true},
	{name: "simple", args: []string{"--compilation_level", "SIMPLE"}, perFile: true},
	{name: "simple-srcmap", args: []string{"--compilation_level", "SIMPLE"}, srcMap: true, perFile: true},
	{name: "simple-es2015", args: []string{"--compilation_level", "SIMPLE", "--language_out", "ECMASCRIPT_2015"}, perFile: true},
	{name: "simple-es5", args: []string{"--compilation_level", "SIMPLE", "--language_out", "ECMASCRIPT5"}, perFile: true},
	{name: "simple-verbose", args: []string{"--compilation_level", "SIMPLE", "--warning_level", "VERBOSE"}, perFile: true},
	{name: "simple-quiet", args: []string{"--compilation_level", "SIMPLE", "--warning_level", "QUIET"}, perFile: true},
	{name: "advanced", args: []string{"--compilation_level", "ADVANCED"}, perFile: true},
	{name: "advanced-srcmap-pretty", args: []string{"--compilation_level", "ADVANCED", "--formatting", "PRETTY_PRINT"}, srcMap: true, perFile: true},
}

// flagJobs are jobs without corpus inputs.
var flagJobs = [][]string{
	{"--help"},
	{"--help_markdown"},
	{"--version"},
	{"--bogus"},
	{"--compilation_level=FOO"},
	{"--warning_level=LOUD"},
	{"--language_in=ES1", "--js", "in/none.js"},
	{"--js"},
	{"--charset=bad name"},
	{"--charset=windows-1252"},
	{"--define="},
	{"--output_wrapper=x"},
	{"--chunk", "x"},
	{"--js", "in/does-not-exist.js"},
	{"--isolation_mode", "IIFE", "--output_wrapper", "%output%"},
	{"--compilation_level", "BUNDLE"},
	{"--renaming=false", "--compilation_level", "ADVANCED"},
	{"--browser_featureset_year", "2020", "--language_out", "ECMASCRIPT_NEXT"},
	{"--flagfile", "in/does-not-exist.flags"},
	{"--jscomp_error", "nosuchgroup", "--js", "in/none.js"},
}

type job struct {
	id   string
	conf string
	args []string
}

func main() {
	flag.Parse()
	switch *suite {
	case "fixtures":
		os.Exit(runFixtures())
	case "matrix":
		os.Exit(runMatrix())
	case "sweep":
		os.Exit(runSweep())
	case "projects":
		os.Exit(runProjects())
	}
	fmt.Fprintln(os.Stderr, "unknown -suite", *suite)
	os.Exit(2)
}

// env is the process environment for both sides: no JAVA_TOOL_OPTIONS and
// the C/POSIX locale.
func env() []string {
	var e []string
	for _, kv := range os.Environ() {
		k := kv[:strings.IndexByte(kv, '=')]
		switch k {
		case "JAVA_TOOL_OPTIONS", "LANG", "LANGUAGE":
			continue
		}
		if strings.HasPrefix(k, "LC_") {
			continue
		}
		e = append(e, kv)
	}
	return e
}

func buildJsgo(dir string) string {
	if *jsgoBin != "" {
		abs, _ := filepath.Abs(*jsgoBin)
		return abs
	}
	bin := filepath.Join(dir, "jsgo")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/jsgo")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "building jsgo:", err)
		os.Exit(2)
	}
	return bin
}

// runJsgo runs one job in dir, writing <res>/stdout, stderr and exit.
func runJsgo(bin, dir, res string, args []string) {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	cmd.Env = env()
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if in, err := os.Open(filepath.Join(dir, res, "stdin")); err == nil {
		cmd.Stdin = in
		defer in.Close()
	}
	err := cmd.Run()
	status := 0
	if ee, ok := err.(*exec.ExitError); ok {
		status = ee.ExitCode()
	} else if err != nil {
		status = -1
	}
	os.MkdirAll(filepath.Join(dir, res), 0o777)
	os.WriteFile(filepath.Join(dir, res, "stdout"), stdout.Bytes(), 0o666)
	os.WriteFile(filepath.Join(dir, res, "stderr"), stderr.Bytes(), 0o666)
	os.WriteFile(filepath.Join(dir, res, "exit"), []byte(strconv.Itoa(status)+"\n"), 0o666)
}

// stripFrames drops Java stack-frame lines.
func stripFrames(b []byte) []byte {
	var out bytes.Buffer
	for _, line := range bytes.SplitAfter(b, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("\tat ")) || bytes.HasPrefix(line, []byte("\t... ")) {
			continue
		}
		out.Write(line)
	}
	return out.Bytes()
}

func isPending(stderr []byte) bool {
	return bytes.Contains(stderr, []byte("pending port:"))
}

type tally struct {
	same, pending, diff int
	diffs               []string
}

func report(t map[string]*tally, order []string) {
	totalSame, total := 0, 0
	for _, name := range order {
		c := t[name]
		if c == nil {
			continue
		}
		all := c.same + c.pending + c.diff
		total += all
		totalSame += c.same
		fmt.Printf("%-24s %5d/%-5d identical  %5d pending  %5d diff\n", name, c.same, all, c.pending, c.diff)
		for i, d := range c.diffs {
			if i >= 5 {
				fmt.Printf("    ... %d more\n", len(c.diffs)-5)
				break
			}
			fmt.Printf("    %s\n", d)
		}
	}
	fmt.Printf("%-24s %5d/%-5d identical\n", "TOTAL", totalSame, total)
}

// runFixtures compares jsgo with the recorded oracles of testdata.
func runFixtures() int {
	root, _ := os.Getwd()
	tmp, err := os.MkdirTemp("", "clie2e-fixtures")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)
	bin := buildJsgo(tmp)
	type fx struct{ js, oracle, strict string }
	var fixtures []fx
	filepath.WalkDir(filepath.Join(root, "testdata/fixtures"), func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(p, ".js") {
			rel, _ := filepath.Rel(root, p)
			base := strings.TrimSuffix(strings.TrimPrefix(rel, "testdata/fixtures/"), ".js") + ".txt"
			fixtures = append(fixtures, fx{rel, filepath.Join(root, "testdata/oracle", base), filepath.Join(root, "testdata/oracle-strict", base)})
		}
		return nil
	})
	sort.Slice(fixtures, func(i, j int) bool { return fixtures[i].js < fixtures[j].js })
	if *n > 0 && *n < len(fixtures) {
		fixtures = fixtures[:*n]
	}
	base := []string{"--checks_only", "--jscomp_error=checkTypes", "--warning_level=VERBOSE", "--language_in=ECMASCRIPT_NEXT"}
	t := map[string]*tally{"oracle": {}, "oracle-strict": {}}
	var mu sync.Mutex
	sem := make(chan bool, *jobsN)
	var wg sync.WaitGroup
	for _, f := range fixtures {
		for _, kind := range []string{"oracle", "oracle-strict"} {
			args := append([]string{}, base...)
			want := f.oracle
			if kind == "oracle-strict" {
				args = append(args, "--jscomp_warning=strictCheckTypes")
				want = f.strict
			}
			args = append(args, "--js", f.js)
			wg.Add(1)
			sem <- true
			go func(kind, want string, args []string) {
				defer func() { <-sem; wg.Done() }()
				// CliBatch's format: stdout and stderr through one UTF-8
				// PrintStream, then the exit status.
				cmd := exec.Command(bin, args...)
				cmd.Dir = root
				cmd.Env = append(env(), "LC_ALL=C.UTF-8")
				var combined bytes.Buffer
				cmd.Stdout, cmd.Stderr = &combined, &combined
				err := cmd.Run()
				status := 0
				if ee, ok := err.(*exec.ExitError); ok {
					status = ee.ExitCode()
				}
				got := append(combined.Bytes(), []byte("---EXIT---\n"+strconv.Itoa(status)+"\n")...)
				wantB, _ := os.ReadFile(want)
				mu.Lock()
				defer mu.Unlock()
				c := t[kind]
				switch {
				case bytes.Equal(stripFrames(got), stripFrames(wantB)):
					c.same++
				case isPending(got):
					c.pending++
				default:
					c.diff++
					c.diffs = append(c.diffs, want)
				}
			}(kind, want, args)
		}
	}
	wg.Wait()
	report(t, []string{"oracle", "oracle-strict"})
	return 0
}

// runMatrix runs the corpus matrix on both sides.
func runMatrix() int {
	if *work == "" {
		fmt.Fprintln(os.Stderr, "-work is required")
		return 2
	}
	wd, _ := filepath.Abs(*work)
	javaDir, goDir := filepath.Join(wd, "java"), filepath.Join(wd, "go")
	for _, d := range []string{javaDir, goDir} {
		os.RemoveAll(d)
		os.MkdirAll(d, 0o777)
		os.Symlink(*corpus, filepath.Join(d, "in"))
	}
	bin := buildJsgo(wd)

	// Sample corpus files, without those of
	// tools/oracle/gpl-excluded-corpus.txt.
	files, err := oraclecorpus.ReadList(*list)
	if err != nil {
		panic(err)
	}
	r := rand.New(rand.NewSource(*seed))
	r.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })
	count := *n
	if count == 0 {
		count = 20
	}
	if count > len(files) {
		count = len(files)
	}
	files = files[:count]

	selected := map[string]bool{}
	for _, c := range strings.Split(*onlyConf, ",") {
		if c != "" {
			selected[c] = true
		}
	}
	var jobs []job
	order := []string{}
	for _, c := range configs {
		if len(selected) > 0 && !selected[c.name] {
			continue
		}
		order = append(order, c.name)
		for i, f := range files {
			id := fmt.Sprintf("%s/%04d", c.name, i)
			args := append(append([]string{}, c.args...), "--js", "in/"+f, "--js_output_file", "out/"+id+"/out.js")
			if c.srcMap {
				args = append(args, "--create_source_map", "out/"+id+"/out.js.map")
			}
			jobs = append(jobs, job{id, c.name, args})
		}
	}
	if len(selected) == 0 || selected["flags"] {
		order = append(order, "flags")
		for i, a := range flagJobs {
			jobs = append(jobs, job{fmt.Sprintf("flags/%04d", i), "flags", a})
		}
	}
	for _, j := range jobs {
		os.MkdirAll(filepath.Join(javaDir, "out", j.id), 0o777)
		os.MkdirAll(filepath.Join(goDir, "out", j.id), 0o777)
	}

	// Java side: shards of one JVM each. The JVMs run in the work
	// directory, so CliE2E's CorpusExclusions gets the list's path.
	excludedList, err := filepath.Abs(oraclecorpus.ListPath)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for s := 0; s < *shards; s++ {
		var buf bytes.Buffer
		for i := s; i < len(jobs); i += *shards {
			buf.WriteString("res/" + jobs[i].id + "\t" + strings.Join(jobs[i].args, "\t") + "\n")
		}
		wg.Add(1)
		go func(input []byte) {
			defer wg.Done()
			cmd := exec.Command("nice", "-n", "10", "java", "-Djsgo.excludedCorpus="+excludedList, "-cp", *jar+":"+*javaBld, "CliE2E")
			cmd.Dir = javaDir
			cmd.Env = env()
			cmd.Stdin = bytes.NewReader(input)
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, "java shard:", err)
			}
		}(buf.Bytes())
	}
	// Go side.
	sem := make(chan bool, *jobsN)
	for _, j := range jobs {
		wg.Add(1)
		sem <- true
		go func(j job) {
			defer func() { <-sem; wg.Done() }()
			runJsgo(bin, goDir, "res/"+j.id, j.args)
		}(j)
	}
	wg.Wait()

	t := map[string]*tally{}
	for _, j := range jobs {
		c := t[j.conf]
		if c == nil {
			c = &tally{}
			t[j.conf] = c
		}
		why := compareJob(javaDir, goDir, j)
		switch {
		case why == "":
			c.same++
		case why == "pending":
			c.pending++
		default:
			c.diff++
			c.diffs = append(c.diffs, j.id+": "+why+" ("+strings.Join(j.args, " ")+")")
		}
	}
	report(t, order)
	return 0
}

// compareJob returns "" when both sides agree, "pending" when jsgo reached
// a pending port, and what differs otherwise.
func compareJob(javaDir, goDir string, j job) string {
	read := func(dir, name string) []byte {
		b, _ := os.ReadFile(filepath.Join(dir, "res", j.id, name))
		return b
	}
	if isPending(read(goDir, "stderr")) {
		return "pending"
	}
	var why []string
	for _, f := range []string{"exit", "stdout", "stderr"} {
		a, b := read(javaDir, f), read(goDir, f)
		if f == "stderr" {
			a, b = stripFrames(a), stripFrames(b)
		}
		if !bytes.Equal(a, b) {
			why = append(why, f+firstDiff(a, b))
		}
	}
	// Output files.
	files := map[string]bool{}
	for _, dir := range []string{javaDir, goDir} {
		root := filepath.Join(dir, "out", j.id)
		filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				rel, _ := filepath.Rel(root, p)
				files[rel] = true
			}
			return nil
		})
	}
	names := make([]string, 0, len(files))
	for f := range files {
		names = append(names, f)
	}
	sort.Strings(names)
	for _, f := range names {
		a, errA := os.ReadFile(filepath.Join(javaDir, "out", j.id, f))
		b, errB := os.ReadFile(filepath.Join(goDir, "out", j.id, f))
		if (errA == nil) != (errB == nil) || !bytes.Equal(a, b) {
			why = append(why, "out/"+f+firstDiff(a, b))
		}
	}
	return strings.Join(why, "; ")
}

// firstDiff describes the first differing line when -v is set.
func firstDiff(a, b []byte) string {
	if !*keep {
		return ""
	}
	sa := bufio.NewScanner(bytes.NewReader(a))
	sb := bufio.NewScanner(bytes.NewReader(b))
	sa.Buffer(nil, 1<<26)
	sb.Buffer(nil, 1<<26)
	for line := 1; ; line++ {
		okA, okB := sa.Scan(), sb.Scan()
		if !okA && !okB {
			return ""
		}
		la, lb := sa.Text(), sb.Text()
		if la != lb || okA != okB {
			cut := func(s string) string {
				if len(s) > 120 {
					return s[:120] + "..."
				}
				return s
			}
			return fmt.Sprintf(" [line %d: java %q / go %q]", line, cut(la), cut(lb))
		}
	}
}
