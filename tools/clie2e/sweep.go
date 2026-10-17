package main

// The sweep suite: every file of a corpus list, compiled alone in a few
// configurations by both compilers, each result classified as identical,
// diff, jsgo panic, jsgo hang (timeout) or Java error.
//
// Results are kept per job under <work>/{java,go}/res/<config>/<file>/
// (stdout, stderr, exit, millis), so a rerun can keep the Java side
// (-reuse-java) and redo only jsgo after a fix. The summary goes to
// stdout and to <work>/sweep.tsv (one line per job: config, file, class,
// Java ms, jsgo ms, signature).

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	oraclecorpus "github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

var (
	timeout   = flag.Duration("timeout", 3*time.Minute, "sweep: jsgo time limit per job (then SIGQUIT, for its goroutine dump)")
	reuseJava = flag.Bool("reuse-java", false, "sweep: keep the Java results already in the work directory")
	locale    = flag.String("lang", "C.UTF-8", "sweep: LANG for both compilers")
)

// sweepConfigs are the sweep's columns: flags added before --js.
var sweepConfigs = []config{
	{name: "simple", args: []string{"--compilation_level", "SIMPLE"}},
	{name: "advanced", args: []string{"--compilation_level", "ADVANCED"}},
	{name: "checks", args: []string{"--checks_only", "--jscomp_warning=checkTypes"}},
}

// Result classes, in report order.
const (
	clIdentical = "identical"
	clDiff      = "diff"
	clPanic     = "jsgo-panic"
	clHang      = "jsgo-hang"
	clJavaError = "java-error"
	clMissing   = "java-missing"
)

var classes = []string{clIdentical, clDiff, clPanic, clHang, clJavaError, clMissing}

type sweepJob struct {
	conf, file string
	args       []string
}

func (j sweepJob) res() string { return filepath.Join("res", j.conf, j.file) }

func sweepEnv() []string {
	return append(env(), "LANG="+*locale)
}

func runSweep() int {
	if *work == "" {
		fmt.Fprintln(os.Stderr, "-work is required")
		return 2
	}
	wd, _ := filepath.Abs(*work)
	javaDir, goDir := filepath.Join(wd, "java"), filepath.Join(wd, "go")
	if !*reuseJava {
		os.RemoveAll(javaDir)
	}
	os.RemoveAll(goDir)
	for _, d := range []string{javaDir, goDir} {
		os.MkdirAll(d, 0o777)
		os.Remove(filepath.Join(d, "in"))
		os.Symlink(*corpus, filepath.Join(d, "in"))
	}
	bin := buildJsgo(wd)

	files, err := oraclecorpus.ReadList(*list)
	if err != nil {
		panic(err)
	}
	if *n > 0 && *n < len(files) {
		files = files[:*n]
	}
	selected := map[string]bool{}
	for _, c := range strings.Split(*onlyConf, ",") {
		if c != "" {
			selected[c] = true
		}
	}
	var order []string
	var confs []config
	for _, c := range sweepConfigs {
		if len(selected) == 0 || selected[c.name] {
			confs = append(confs, c)
			order = append(order, c.name)
		}
	}
	// File-major order, so that both sides reach the same jobs at about
	// the same time.
	var jobs []sweepJob
	for _, f := range files {
		for _, c := range confs {
			args := append(append([]string{}, c.args...), "--js", "in/"+f)
			jobs = append(jobs, sweepJob{c.name, f, args})
		}
	}

	var wg sync.WaitGroup
	// Java side: shards of one JVM each (CliE2E).
	excludedList, err := filepath.Abs(oraclecorpus.ListPath)
	if err != nil {
		panic(err)
	}
	var javaJobs []sweepJob
	for _, j := range jobs {
		if _, err := os.Stat(filepath.Join(javaDir, j.res(), "exit")); err == nil && *reuseJava {
			continue
		}
		javaJobs = append(javaJobs, j)
	}
	for s := 0; s < *shards; s++ {
		var buf bytes.Buffer
		for i := s; i < len(javaJobs); i += *shards {
			buf.WriteString(javaJobs[i].res() + "\t" + strings.Join(javaJobs[i].args, "\t") + "\n")
		}
		if buf.Len() == 0 {
			continue
		}
		wg.Add(1)
		go func(input []byte) {
			defer wg.Done()
			cmd := exec.Command("nice", "-n", "10", "java", "-Djsgo.excludedCorpus="+excludedList, "-cp", *jar+":"+*javaBld, "CliE2E")
			cmd.Dir = javaDir
			cmd.Env = sweepEnv()
			cmd.Stdin = bytes.NewReader(input)
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, "java shard:", err)
			}
		}(buf.Bytes())
	}
	// Go side.
	sem := make(chan bool, *jobsN)
	var done, total = 0, len(jobs)
	var mu sync.Mutex
	start := time.Now()
	for _, j := range jobs {
		wg.Add(1)
		sem <- true
		go func(j sweepJob) {
			defer func() { <-sem; wg.Done() }()
			runSweepJsgo(bin, goDir, j)
			mu.Lock()
			done++
			if done%100 == 0 {
				fmt.Fprintf(os.Stderr, "jsgo: %d/%d jobs (%s)\n", done, total, time.Since(start).Round(time.Second))
			}
			mu.Unlock()
		}(j)
	}
	wg.Wait()

	// Classify.
	counts := map[string]map[string]int{}
	var lines []string
	for _, j := range jobs {
		class, sig := classify(javaDir, goDir, j)
		if counts[j.conf] == nil {
			counts[j.conf] = map[string]int{}
		}
		counts[j.conf][class]++
		lines = append(lines, strings.Join([]string{
			j.conf, j.file, class,
			readTrim(filepath.Join(javaDir, j.res(), "millis")),
			readTrim(filepath.Join(goDir, j.res(), "millis")),
			sig}, "\t"))
	}
	os.WriteFile(filepath.Join(wd, "sweep.tsv"), []byte(strings.Join(lines, "\n")+"\n"), 0o666)

	fmt.Printf("%-10s", "config")
	for _, c := range classes {
		fmt.Printf(" %12s", c)
	}
	fmt.Println()
	for _, name := range order {
		fmt.Printf("%-10s", name)
		for _, c := range classes {
			fmt.Printf(" %12d", counts[name][c])
		}
		fmt.Println()
	}
	// Panic and hang signatures, most frequent first.
	sigs := map[string]int{}
	example := map[string]string{}
	for _, l := range lines {
		f := strings.Split(l, "\t")
		if f[2] == clPanic || f[2] == clHang || f[2] == clJavaError {
			k := f[2] + "\t" + f[5]
			sigs[k]++
			if example[k] == "" {
				example[k] = f[0] + " " + f[1]
			}
		}
	}
	keys := make([]string, 0, len(sigs))
	for k := range sigs {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(a, b int) bool {
		if sigs[keys[a]] != sigs[keys[b]] {
			return sigs[keys[a]] > sigs[keys[b]]
		}
		return keys[a] < keys[b]
	})
	for _, k := range keys {
		fmt.Printf("%5d  %s\n       e.g. %s\n", sigs[k], k, example[k])
	}
	return 0
}

// runSweepJsgo runs one job in dir with the time limit, writing
// <res>/stdout, stderr, exit, millis (and timeout when it hit the limit).
func runSweepJsgo(bin, dir string, j sweepJob) {
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, j.args...)
	cmd.Dir = dir
	cmd.Env = sweepEnv()
	// SIGQUIT makes the Go runtime dump every goroutine on stderr: the
	// hang's location. SIGKILL follows if it does not exit.
	cmd.Cancel = func() error { return cmd.Process.Signal(syscall.SIGQUIT) }
	cmd.WaitDelay = 20 * time.Second
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	t0 := time.Now()
	err := cmd.Run()
	ms := time.Since(t0).Milliseconds()
	status := 0
	if ee, ok := err.(*exec.ExitError); ok {
		status = ee.ExitCode()
	} else if err != nil {
		status = -1
	}
	res := filepath.Join(dir, j.res())
	os.MkdirAll(res, 0o777)
	os.WriteFile(filepath.Join(res, "stdout"), stdout.Bytes(), 0o666)
	os.WriteFile(filepath.Join(res, "stderr"), stderr.Bytes(), 0o666)
	os.WriteFile(filepath.Join(res, "exit"), []byte(strconv.Itoa(status)+"\n"), 0o666)
	os.WriteFile(filepath.Join(res, "millis"), []byte(strconv.FormatInt(ms, 10)+"\n"), 0o666)
	if ctx.Err() != nil {
		os.WriteFile(filepath.Join(res, "timeout"), nil, 0o666)
	}
}

func readTrim(p string) string {
	b, _ := os.ReadFile(p)
	return strings.TrimSpace(string(b))
}

var (
	// A Java exception escaping main, or one CommandLineRunner reports.
	javaExceptionRe = regexp.MustCompile(`(?m)^(Exception in thread "main" |INTERNAL COMPILER ERROR|java\.lang\.\w+(Error|Exception)\b)`)
	goPanicRe       = regexp.MustCompile(`(?m)^(panic: |fatal error: )`)
	goFrameRe       = regexp.MustCompile(`(?m)^\s*(?:at )?(github\.com/paullepoulpe/jsgo/internal/[^\s(]+)`)
)

// classify returns the job's class and, for crashes and hangs, a short
// signature: the exception line and the first frame in the port.
func classify(javaDir, goDir string, j sweepJob) (string, string) {
	read := func(dir, name string) []byte {
		b, _ := os.ReadFile(filepath.Join(dir, j.res(), name))
		return b
	}
	javaExit, err := os.ReadFile(filepath.Join(javaDir, j.res(), "exit"))
	if err != nil {
		return clMissing, ""
	}
	if _, err := os.Stat(filepath.Join(goDir, j.res(), "timeout")); err == nil {
		return clHang, hangSignature(read(goDir, "stderr"))
	}
	javaErr, goErr := stripFrames(read(javaDir, "stderr")), stripFrames(read(goDir, "stderr"))
	if bytes.Equal(javaExit, read(goDir, "exit")) &&
		bytes.Equal(read(javaDir, "stdout"), read(goDir, "stdout")) &&
		bytes.Equal(javaErr, goErr) {
		return clIdentical, ""
	}
	if javaExceptionRe.Match(read(javaDir, "stderr")) {
		return clJavaError, firstMatchLine(javaExceptionRe, read(javaDir, "stderr"))
	}
	goStderr := read(goDir, "stderr")
	if goPanicRe.Match(goStderr) || javaExceptionRe.Match(goStderr) {
		sig := firstMatchLine(goPanicRe, goStderr)
		if sig == "" {
			sig = firstMatchLine(javaExceptionRe, goStderr)
		}
		if m := goFrameRe.FindSubmatch(goStderr); m != nil {
			sig += " @ " + string(m[1])
		}
		return clPanic, sig
	}
	return clDiff, ""
}

func firstMatchLine(re *regexp.Regexp, b []byte) string {
	loc := re.FindIndex(b)
	if loc == nil {
		return ""
	}
	line := b[loc[0]:]
	if i := bytes.IndexByte(line, '\n'); i >= 0 {
		line = line[:i]
	}
	s := string(line)
	if len(s) > 200 {
		s = s[:200]
	}
	return s
}

// hangSignature is the innermost port frame of the main goroutine in the
// SIGQUIT dump.
func hangSignature(stderr []byte) string {
	i := bytes.Index(stderr, []byte("goroutine 1 ["))
	if i < 0 {
		return "no goroutine dump"
	}
	if m := goFrameRe.FindSubmatch(stderr[i:]); m != nil {
		return string(m[1])
	}
	return "no port frame"
}
