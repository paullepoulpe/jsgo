package main

// The projects suite: multi-file, project-level builds (Closure Library
// apps, ES modules and CommonJS packages, chunked output, output options)
// compiled by both compilers, each job classified as in the sweep.
//
// A project is a JSON file of tools/clie2e/projects (see projectFile). It
// names a workspace layout (symlinks into the cache or the repository, and
// small generated files), the inputs (files and globs, expanded here in
// sorted order, without the files of tools/oracle/gpl-excluded-corpus.txt)
// and its jobs: argument vectors over those inputs. A project file can
// also import the end-to-end benchmark's tasks (tools/bench/tasks), one job
// each.
//
// Every job runs in a fresh workspace, <work>/ws/<project>/<job>: first
// Java (CliE2E, one JVM per job, its working directory the workspace),
// then jsgo, in the same directory, so that both see the same relative
// and absolute paths. After each run the files the compiler wrote there
// (every regular file that is not part of the layout) are moved to
// <work>/{java,go}/res/<project>/<job>/files/, next to stdout, stderr,
// exit and millis. Compared byte for byte: stdout, stderr (without Java
// stack frames), the exit status and every output file.
//
// Summary on stdout and in <work>/projects.tsv: project, job, class, Java
// ms, jsgo ms, and what differs.
//
// The Java side needs -javabuild built by tools/clie2e/build-java.sh: CliE2E
// and CorpusExclusions, and a working XtbMessageBundle (the pinned jar's own
// cannot load, so --translations_file fails there with NoClassDefFoundError).

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
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
	projectsDir = flag.String("projects-dir", "tools/clie2e/projects", "projects: directory of the project files")
	projectSet  = flag.String("set", "dev", "projects: which projects to run: dev (development corpus and fixtures), holdout (held-out corpus), all")
	onlyProj    = flag.String("projects", "", "projects: comma-separated project or project/job names to run (default: all of -set)")
	cacheDir    = flag.String("cache", "/home/user/jsgo/.cache", "projects: the cache directory (cache: links)")
)

// projectFile is the format of a project file.
type projectFile struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Set is "dev" (the default: development corpus, bench inputs and
	// fixtures) or "holdout" (the held-out corpus, for final checks only).
	Set string `json:"set"`
	// Links are symlinks of the workspace: workspace path → target, where
	// the target starts with "cache:" (the -cache directory) or "repo:"
	// (the repository root). Parent directories are created.
	Links map[string]string `json:"links"`
	// Files are generated files of the workspace: path → content.
	Files map[string]string `json:"files"`
	// Zips are generated zip files of the workspace (for --jszip): path →
	// the workspace directory whose files it holds, in sorted order, with
	// a fixed modification time.
	Zips map[string]string `json:"zips"`
	// Inputs are the default inputs of the jobs.
	Inputs []inputSpec `json:"inputs"`
	// Flags start every job's argument vector.
	Flags []string `json:"flags"`
	// BenchTasks imports tasks of tools/bench/tasks, one job each, named
	// after the task: its flags, inputs and output arguments.
	BenchTasks []string     `json:"bench_tasks"`
	Jobs       []projectJob `json:"jobs"`
}

// inputSpec is one input or glob of inputs, as in tools/bench/tasks: a
// glob (Python glob.glob with recursive=True semantics) is expanded in
// sorted order, without the paths matching an exclude pattern (on the
// whole path or its base name). Each file becomes `<flag> <file>` (flag
// default --js). Group names the inputs for {inputs:<group>} and
// {count:<group>}.
type inputSpec struct {
	File    string   `json:"file"`
	Glob    string   `json:"glob"`
	Exclude []string `json:"exclude"`
	Flag    string   `json:"flag"`
	Group   string   `json:"group"`
}

// projectJob is one compiler invocation. Its argument vector is the
// project's flags then its args, where
//
//	{inputs}          is replaced by the inputs, as `--js <file>` pairs;
//	{inputs:<group>}  by the inputs of one group;
//	{count:<group>}   (inside an argument) by the number of files in it.
//
// Without an {inputs...} argument, the inputs are appended at the end.
// Inputs, when set, replace the project's inputs. Stdin names a workspace
// file (generated or linked) to use as the compiler's standard input.
//
// Steps, instead of Args, makes a job of several invocations run one after
// the other in the same workspace (a multistage build: each step's args
// are expanded like Args). Each step's stdout, stderr and exit status are
// compared; the output files are compared once all steps have run, except
// the Intermediate ones, which both sides must write but whose content is
// not compared: the compiler state that --filename_to_save_to writes is
// Java serialization in Java and the port's own encoding in jsgo
// (Compiler.saveState), gzipped by different deflaters.
type projectJob struct {
	Name         string      `json:"name"`
	Args         []string    `json:"args"`
	Steps        [][]string  `json:"steps"`
	Intermediate []string    `json:"intermediate"`
	Inputs       []inputSpec `json:"inputs"`
	Stdin        string      `json:"stdin"`
}

// benchTask is the part of a tools/bench/tasks file the suite reads.
type benchTask struct {
	Name       string      `json:"name"`
	Flags      []string    `json:"flags"`
	Inputs     []inputSpec `json:"inputs"`
	OutputArgs []string    `json:"output_args"`
	PerFile    any         `json:"per_file"`
}

// projJob is a job ready to run.
type projJob struct {
	project      *projectFile
	name         string
	steps        [][]string // argument vectors
	intermediate map[string]bool
	stdin        string
}

func (j *projJob) id() string { return j.project.Name + "/" + j.name }

func runProjects() int {
	if *work == "" {
		fmt.Fprintln(os.Stderr, "-work is required")
		return 2
	}
	wd, _ := filepath.Abs(*work)
	repo, _ := os.Getwd()
	timeoutSet := false
	flag.Visit(func(f *flag.Flag) { timeoutSet = timeoutSet || f.Name == "timeout" })
	if !timeoutSet {
		*timeout = 30 * time.Minute // a whole project, not one file
	}
	projects, err := loadProjects(repo)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	selected := map[string]bool{}
	for _, s := range strings.Split(*onlyProj, ",") {
		if s != "" {
			selected[s] = true
		}
	}
	var jobs []*projJob
	for _, p := range projects {
		set := p.Set
		if set == "" {
			set = "dev"
		}
		if *projectSet != "all" && set != *projectSet && !selected[p.Name] {
			continue
		}
		for _, pj := range p.Jobs {
			if len(selected) > 0 && !selected[p.Name] && !selected[p.Name+"/"+pj.Name] {
				continue
			}
			j, err := expandJob(repo, p, pj)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s/%s: %v\n", p.Name, pj.Name, err)
				return 2
			}
			jobs = append(jobs, j)
		}
	}
	if len(jobs) == 0 {
		fmt.Fprintln(os.Stderr, "no jobs selected")
		return 2
	}
	javaDir, goDir := filepath.Join(wd, "java"), filepath.Join(wd, "go")
	if !*classifyOnly {
		bin := ""
		if !*skipGo {
			bin = buildJsgo(wd)
		}
		runProjectJobs(repo, wd, javaDir, goDir, bin, jobs)
	}
	return summarizeProjects(wd, javaDir, goDir, jobs)
}

// loadProjects reads every project file, in name order.
func loadProjects(repo string) ([]*projectFile, error) {
	dir := filepath.Join(repo, *projectsDir)
	names, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	var out []*projectFile
	seen := map[string]bool{}
	for _, f := range names {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		var p projectFile
		dec := json.NewDecoder(bytes.NewReader(b))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&p); err != nil {
			return nil, fmt.Errorf("%s: %v", f, err)
		}
		if p.Name == "" {
			p.Name = strings.TrimSuffix(filepath.Base(f), ".json")
		}
		if seen[p.Name] {
			return nil, fmt.Errorf("%s: duplicate project %s", f, p.Name)
		}
		seen[p.Name] = true
		for _, t := range p.BenchTasks {
			tb, err := os.ReadFile(filepath.Join(repo, "tools/bench/tasks", t+".json"))
			if err != nil {
				return nil, fmt.Errorf("%s: %v", f, err)
			}
			var bt benchTask
			if err := json.Unmarshal(tb, &bt); err != nil {
				return nil, fmt.Errorf("%s: bench task %s: %v", f, t, err)
			}
			if bt.PerFile != nil {
				return nil, fmt.Errorf("%s: bench task %s is per-file", f, t)
			}
			args := append(append(append([]string{}, bt.Flags...), "{inputs}"), bt.OutputArgs...)
			p.Jobs = append(p.Jobs, projectJob{Name: bt.Name, Args: args, Inputs: bt.Inputs})
		}
		names := map[string]bool{}
		for _, j := range p.Jobs {
			if j.Name == "" || names[j.Name] || strings.ContainsAny(j.Name, "/\t") {
				return nil, fmt.Errorf("%s: bad or duplicate job name %q", f, j.Name)
			}
			names[j.Name] = true
		}
		out = append(out, &p)
	}
	return out, nil
}

// inputGroup is the expanded inputs of one group.
type inputGroup struct {
	args  []string // --js <file> pairs
	count int
}

// expandJob builds the job's argument vectors, expanding its inputs.
func expandJob(repo string, p *projectFile, pj projectJob) (*projJob, error) {
	specs := pj.Inputs
	if specs == nil {
		specs = p.Inputs
	}
	// The globs are expanded against the links' targets.
	groups := map[string]*inputGroup{}
	var all []string
	for _, s := range specs {
		files, err := expandSpec(repo, p, s)
		if err != nil {
			return nil, err
		}
		flagName := s.Flag
		if flagName == "" {
			flagName = "--js"
		}
		g := groups[s.Group]
		if g == nil {
			g = &inputGroup{}
			groups[s.Group] = g
		}
		for _, f := range files {
			g.args = append(g.args, flagName, f)
			g.count++
			all = append(all, flagName, f)
		}
	}
	stepArgs := pj.Steps
	if stepArgs == nil {
		stepArgs = [][]string{pj.Args}
	} else if pj.Args != nil {
		return nil, fmt.Errorf("both args and steps")
	}
	j := &projJob{project: p, name: pj.Name, stdin: pj.Stdin, intermediate: map[string]bool{}}
	for _, f := range pj.Intermediate {
		j.intermediate[filepath.FromSlash(f)] = true
	}
	for _, sa := range stepArgs {
		args, err := expandArgs(p, sa, groups, all)
		if err != nil {
			return nil, err
		}
		j.steps = append(j.steps, args)
	}
	return j, nil
}

// expandArgs builds one argument vector of a job (see projectJob).
func expandArgs(p *projectFile, jobArgs []string, groups map[string]*inputGroup, all []string) ([]string, error) {
	countRe := regexp.MustCompile(`\{count:([^}]*)\}`)
	var args []string
	usedInputs := false
	for _, a := range append(append([]string{}, p.Flags...), jobArgs...) {
		switch {
		case a == "{inputs}":
			args = append(args, all...)
			usedInputs = true
		case strings.HasPrefix(a, "{inputs:") && strings.HasSuffix(a, "}"):
			g := groups[a[len("{inputs:"):len(a)-1]]
			if g == nil {
				return nil, fmt.Errorf("no input group in %s", a)
			}
			args = append(args, g.args...)
			usedInputs = true
		default:
			var bad error
			a = countRe.ReplaceAllStringFunc(a, func(m string) string {
				g := groups[countRe.FindStringSubmatch(m)[1]]
				if g == nil {
					bad = fmt.Errorf("no input group in %s", m)
					return m
				}
				return strconv.Itoa(g.count)
			})
			if bad != nil {
				return nil, bad
			}
			args = append(args, a)
		}
	}
	if !usedInputs {
		args = append(args, all...)
	}
	// tools/oracle/gpl-excluded-corpus.txt: no job may name one of those
	// files, and a glob left for the compiler to expand must not reach
	// into a corpus.
	for i, a := range args {
		if oraclecorpus.IsExcluded(a) {
			return nil, fmt.Errorf("names an excluded corpus file: %s", a)
		}
		// CommandLineRunner expands "*" in --js values (findJsFiles).
		jsValue := strings.HasPrefix(a, "--js=") || (i > 0 && args[i-1] == "--js")
		if jsValue && strings.Contains(a, "*") {
			for _, t := range p.Links {
				if strings.HasPrefix(t, "cache:corpus") {
					return nil, fmt.Errorf("compiler-expanded glob %q in a project linked to a corpus", a)
				}
			}
		}
	}
	return args, nil
}

// resolveLink returns the target of a link value.
func resolveLink(repo, target string) (string, error) {
	switch {
	case strings.HasPrefix(target, "cache:"):
		return filepath.Join(*cacheDir, strings.TrimPrefix(target, "cache:")), nil
	case strings.HasPrefix(target, "repo:"):
		return filepath.Join(repo, strings.TrimPrefix(target, "repo:")), nil
	}
	return "", fmt.Errorf("link target %q: want cache: or repo:", target)
}

// fsPath maps a workspace path to the file system through the project's
// links (the longest link that is a prefix of it).
func fsPath(repo string, p *projectFile, wsPath string) (string, error) {
	best := ""
	for l := range p.Links {
		if (wsPath == l || strings.HasPrefix(wsPath, l+"/")) && len(l) > len(best) {
			best = l
		}
	}
	if best == "" {
		return "", fmt.Errorf("%s is not under a link", wsPath)
	}
	t, err := resolveLink(repo, p.Links[best])
	if err != nil {
		return "", err
	}
	return filepath.Join(t, strings.TrimPrefix(wsPath, best)), nil
}

// expandSpec returns the workspace paths of one input spec, without the
// excluded corpus files.
func expandSpec(repo string, p *projectFile, s inputSpec) ([]string, error) {
	if s.File != "" {
		if _, ok := p.Files[s.File]; ok {
			return []string{s.File}, nil
		}
		if _, ok := p.Zips[s.File]; ok {
			return []string{s.File}, nil
		}
		fp, err := fsPath(repo, p, s.File)
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(fp); err != nil {
			return nil, fmt.Errorf("missing input %s (%s)", s.File, fp)
		}
		return oraclecorpus.Filter([]string{s.File}), nil
	}
	var files []string
	if err := globWalk(repo, p, "", strings.Split(s.Glob, "/"), &files); err != nil {
		return nil, err
	}
	sort.Strings(files)
	var out []string
	for _, f := range files {
		excluded := false
		for _, e := range s.Exclude {
			if fnmatch(e, f) || fnmatch(e, path.Base(f)) {
				excluded = true
			}
		}
		if !excluded {
			out = append(out, f)
		}
	}
	out = oraclecorpus.Filter(out)
	if len(out) == 0 {
		return nil, fmt.Errorf("input glob matched nothing: %s", s.Glob)
	}
	return out, nil
}

// globWalk matches the pattern segments under the workspace directory
// prefix, like Python's glob.glob(recursive=True), which bench.py uses: "**"
// matches zero or more directories, symlinks are followed, and a wildcard
// does not match a name starting with "." unless the segment does.
func globWalk(repo string, p *projectFile, prefix string, segs []string, out *[]string) error {
	join := func(name string) string {
		if prefix == "" {
			return name
		}
		return prefix + "/" + name
	}
	if len(segs) == 0 {
		return nil
	}
	seg, rest := segs[0], segs[1:]
	if !strings.ContainsAny(seg, "*?[") {
		next := join(seg)
		if len(rest) == 0 {
			if fp, err := fsPath(repo, p, next); err == nil {
				if _, err := os.Stat(fp); err == nil {
					*out = append(*out, next)
				}
			}
			return nil
		}
		return globWalk(repo, p, next, rest, out)
	}
	dir := prefix
	fp, err := fsPath(repo, p, dir)
	if err != nil {
		// The workspace root holds only links: list them.
		if dir != "" {
			return err
		}
	}
	var names []string
	if dir == "" {
		seen := map[string]bool{}
		for l := range p.Links {
			n := strings.SplitN(l, "/", 2)[0]
			if !seen[n] {
				seen[n] = true
				names = append(names, n)
			}
		}
	} else {
		ents, err := os.ReadDir(fp)
		if err != nil {
			return nil
		}
		for _, e := range ents {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	isDir := func(ws string) bool {
		fp, err := fsPath(repo, p, ws)
		if err != nil {
			return true // a directory of the layout above the links
		}
		st, err := os.Stat(fp)
		return err == nil && st.IsDir()
	}
	if seg == "**" {
		// Zero directories.
		if err := globWalk(repo, p, prefix, rest, out); err != nil {
			return err
		}
		for _, n := range names {
			if strings.HasPrefix(n, ".") || !isDir(join(n)) {
				continue
			}
			// One or more: recurse with "**" still in front.
			if err := globWalk(repo, p, join(n), segs, out); err != nil {
				return err
			}
		}
		return nil
	}
	for _, n := range names {
		if strings.HasPrefix(n, ".") && !strings.HasPrefix(seg, ".") {
			continue
		}
		if !fnmatch(seg, n) {
			continue
		}
		if len(rest) == 0 {
			*out = append(*out, join(n))
		} else if isDir(join(n)) {
			if err := globWalk(repo, p, join(n), rest, out); err != nil {
				return err
			}
		}
	}
	return nil
}

// fnmatch is Python's fnmatch.fnmatchcase, which bench.py uses: "*"
// matches any string, "/" included, "?" one character, "[seq]" and
// "[!seq]" a character class.
func fnmatch(pattern, name string) bool {
	var re strings.Builder
	re.WriteString("^(?s:")
	for i := 0; i < len(pattern); i++ {
		switch c := pattern[i]; c {
		case '*':
			re.WriteString(".*")
		case '?':
			re.WriteString(".")
		case '[':
			j := i + 1
			if j < len(pattern) && pattern[j] == '!' {
				j++
			}
			if j < len(pattern) && pattern[j] == ']' {
				j++
			}
			for j < len(pattern) && pattern[j] != ']' {
				j++
			}
			if j >= len(pattern) {
				re.WriteString(`\[`)
				continue
			}
			class := pattern[i+1 : j]
			if strings.HasPrefix(class, "!") {
				class = "^" + class[1:]
			} else if strings.HasPrefix(class, "^") {
				class = `\` + class
			}
			re.WriteString("[" + strings.ReplaceAll(class, `\`, `\\`) + "]")
			i = j
		default:
			re.WriteString(regexp.QuoteMeta(string(c)))
		}
	}
	re.WriteString(")$")
	m, err := regexp.MatchString(re.String(), name)
	return err == nil && m
}

// makeWorkspace (re)creates the job's workspace: its links, its generated
// files and an empty out/ directory for the outputs (the compiler creates
// the directories of chunk outputs, but not those of --js_output_file and
// the other output files).
func makeWorkspace(repo, ws string, p *projectFile) error {
	if err := os.RemoveAll(ws); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(ws, "out"), 0o777); err != nil {
		return err
	}
	for l, t := range p.Links {
		target, err := resolveLink(repo, t)
		if err != nil {
			return err
		}
		if _, err := os.Stat(target); err != nil {
			return fmt.Errorf("link %s: %v", l, err)
		}
		dst := filepath.Join(ws, filepath.FromSlash(l))
		os.MkdirAll(filepath.Dir(dst), 0o777)
		if err := os.Symlink(target, dst); err != nil {
			return err
		}
	}
	for f, content := range p.Files {
		dst := filepath.Join(ws, filepath.FromSlash(f))
		os.MkdirAll(filepath.Dir(dst), 0o777)
		if err := os.WriteFile(dst, []byte(content), 0o666); err != nil {
			return err
		}
	}
	for z, dir := range p.Zips {
		if err := writeZip(filepath.Join(ws, filepath.FromSlash(z)), filepath.Join(ws, filepath.FromSlash(dir))); err != nil {
			return fmt.Errorf("zip %s: %v", z, err)
		}
	}
	return nil
}

// writeZip writes the regular files under dir (following symlinks) to a
// zip file, sorted by path, with a fixed time.
func writeZip(dst, dir string) error {
	var names []string
	var walk func(rel string) error
	walk = func(rel string) error {
		ents, err := os.ReadDir(filepath.Join(dir, rel))
		if err != nil {
			return err
		}
		for _, e := range ents {
			r := path.Join(rel, e.Name())
			st, err := os.Stat(filepath.Join(dir, r))
			if err != nil {
				return err
			}
			if st.IsDir() {
				if err := walk(r); err != nil {
					return err
				}
			} else {
				names = append(names, r)
			}
		}
		return nil
	}
	if err := walk(""); err != nil {
		return err
	}
	sort.Strings(names)
	os.MkdirAll(filepath.Dir(dst), 0o777)
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	for _, n := range names {
		b, err := os.ReadFile(filepath.Join(dir, n))
		if err != nil {
			return err
		}
		w, err := zw.CreateHeader(&zip.FileHeader{Name: n, Method: zip.Deflate,
			Modified: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)})
		if err != nil {
			return err
		}
		if _, err := w.Write(b); err != nil {
			return err
		}
	}
	return zw.Close()
}

// collectOutputs moves the files the compiler wrote in the workspace to
// res/files.
func collectOutputs(ws, res string, p *projectFile) {
	layout := map[string]bool{}
	for f := range p.Files {
		layout[filepath.FromSlash(f)] = true
	}
	for f := range p.Zips {
		layout[filepath.FromSlash(f)] = true
	}
	dst := filepath.Join(res, "files")
	os.RemoveAll(dst)
	filepath.WalkDir(ws, func(pth string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		rel, _ := filepath.Rel(ws, pth)
		if layout[rel] {
			return nil
		}
		to := filepath.Join(dst, rel)
		os.MkdirAll(filepath.Dir(to), 0o777)
		os.Rename(pth, to)
		return nil
	})
}

func (j *projJob) res(side string) string {
	return filepath.Join(side, "res", j.project.Name, j.name)
}

// stepRes is where step i's stdout, stderr, exit and millis go: the job's
// result directory, or step<i> in it for a job of several steps.
func (j *projJob) stepRes(res string, i int) string {
	if len(j.steps) == 1 {
		return res
	}
	return filepath.Join(res, "step"+strconv.Itoa(i))
}

// runProjectJobs runs the jobs, -jobs at a time, each Java then jsgo.
func runProjectJobs(repo, wd, javaDir, goDir, bin string, jobs []*projJob) {
	excludedList, err := filepath.Abs(oraclecorpus.ListPath)
	if err != nil {
		panic(err)
	}
	var mu sync.Mutex
	sem := make(chan bool, *jobsN)
	var wg sync.WaitGroup
	start := time.Now()
	done := 0
	for _, j := range jobs {
		wg.Add(1)
		sem <- true
		go func(j *projJob) {
			defer func() { <-sem; wg.Done() }()
			ws := filepath.Join(wd, "ws", j.project.Name, j.name)
			javaRes, goRes := j.res(javaDir), j.res(goDir)
			_, haveJava := os.Stat(filepath.Join(j.stepRes(javaRes, len(j.steps)-1), "exit"))
			if !*skipJava && !(*reuseJava && haveJava == nil) {
				if err := runProjectJava(repo, ws, javaRes, excludedList, j); err != nil {
					fmt.Fprintf(os.Stderr, "%s: java: %v\n", j.id(), err)
				}
			}
			if !*skipGo {
				if err := runProjectGo(repo, ws, goRes, bin, j); err != nil {
					fmt.Fprintf(os.Stderr, "%s: jsgo: %v\n", j.id(), err)
				}
			}
			os.RemoveAll(ws)
			mu.Lock()
			done++
			fmt.Fprintf(os.Stderr, "projects: %d/%d %s (%s)\n", done, len(jobs), j.id(), time.Since(start).Round(time.Second))
			mu.Unlock()
		}(j)
	}
	wg.Wait()
}

func prepareRes(res string) error {
	if err := os.RemoveAll(res); err != nil {
		return err
	}
	return os.MkdirAll(res, 0o777)
}

func runProjectJava(repo, ws, res, excludedList string, j *projJob) error {
	if err := makeWorkspace(repo, ws, j.project); err != nil {
		return err
	}
	if err := prepareRes(res); err != nil {
		return err
	}
	// One CliE2E line per step, all in one JVM.
	var lines strings.Builder
	for i, args := range j.steps {
		stepRes, _ := filepath.Abs(j.stepRes(res, i))
		os.MkdirAll(stepRes, 0o777)
		if j.stdin != "" {
			b, err := os.ReadFile(filepath.Join(ws, j.stdin))
			if err != nil {
				return err
			}
			os.WriteFile(filepath.Join(stepRes, "stdin"), b, 0o666)
		}
		for _, a := range args {
			if strings.ContainsAny(a, "\t\n") {
				return fmt.Errorf("argument with a tab or newline: %q", a)
			}
		}
		lines.WriteString(stepRes + "\t" + strings.Join(args, "\t") + "\n")
	}
	// -javabuild goes first: tools/clie2e/build-java.sh puts there the
	// XtbMessageBundle that the jar cannot load.
	cmd := exec.Command("nice", "-n", "10", "java", "-Xmx"+*javaHeap, "-Djsgo.excludedCorpus="+excludedList, "-cp", *javaBld+":"+*jar, "CliE2E")
	cmd.Dir = ws
	cmd.Env = sweepEnv()
	cmd.Stdin = strings.NewReader(lines.String())
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	err := cmd.Run()
	if out.Len() > 0 {
		os.WriteFile(filepath.Join(res, "jvm.log"), out.Bytes(), 0o666)
	}
	collectOutputs(ws, res, j.project)
	if err != nil {
		return fmt.Errorf("%v: %s", err, lastLines(out.Bytes(), 5))
	}
	return nil
}

func runProjectGo(repo, ws, res, bin string, j *projJob) error {
	if err := makeWorkspace(repo, ws, j.project); err != nil {
		return err
	}
	if err := prepareRes(res); err != nil {
		return err
	}
	for i, args := range j.steps {
		hung, err := runProjectGoStep(ws, j.stepRes(res, i), bin, args, j.stdin)
		if err != nil {
			return err
		}
		if hung {
			break
		}
	}
	collectOutputs(ws, res, j.project)
	return nil
}

// runProjectGoStep runs one invocation and reports whether it hit the time
// limit.
func runProjectGoStep(ws, res, bin string, args []string, stdin string) (bool, error) {
	os.MkdirAll(res, 0o777)
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "nice", append([]string{"-n", "10", bin}, args...)...)
	cmd.Dir = ws
	cmd.Env = sweepEnv()
	cmd.Cancel = func() error { return cmd.Process.Signal(syscall.SIGQUIT) }
	cmd.WaitDelay = 20 * time.Second
	if stdin != "" {
		in, err := os.Open(filepath.Join(ws, stdin))
		if err != nil {
			return false, err
		}
		defer in.Close()
		cmd.Stdin = in
	}
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
	os.WriteFile(filepath.Join(res, "stdout"), stdout.Bytes(), 0o666)
	os.WriteFile(filepath.Join(res, "stderr"), stderr.Bytes(), 0o666)
	os.WriteFile(filepath.Join(res, "exit"), []byte(strconv.Itoa(status)+"\n"), 0o666)
	os.WriteFile(filepath.Join(res, "millis"), []byte(strconv.FormatInt(ms, 10)+"\n"), 0o666)
	if ctx.Err() != nil {
		os.WriteFile(filepath.Join(res, "timeout"), nil, 0o666)
		return true, nil
	}
	return false, nil
}

func lastLines(b []byte, n int) string {
	lines := strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, " | ")
}

// classifyProject returns the job's class and a detail: what differs, or
// the crash signature.
func classifyProject(javaDir, goDir string, j *projJob) (string, string) {
	jr, gr := j.res(javaDir), j.res(goDir)
	read := func(dir, name string) []byte {
		b, _ := os.ReadFile(filepath.Join(dir, name))
		return b
	}
	last := len(j.steps) - 1
	if _, err := os.Stat(filepath.Join(j.stepRes(jr, last), "exit")); err != nil {
		return clMissing, ""
	}
	var why []string
	var javaStderr, goStderr []byte
	for i := range j.steps {
		js, gs := j.stepRes(jr, i), j.stepRes(gr, i)
		if _, err := os.Stat(filepath.Join(gs, "timeout")); err == nil {
			return clHang, hangSignature(read(gs, "stderr"))
		}
		prefix := ""
		if len(j.steps) > 1 {
			prefix = "step" + strconv.Itoa(i) + "/"
		}
		for _, f := range []string{"exit", "stdout", "stderr"} {
			a, b := read(js, f), read(gs, f)
			if f == "stderr" {
				javaStderr, goStderr = append(javaStderr, a...), append(goStderr, b...)
				a, b = stripFrames(a), stripFrames(b)
			}
			if !bytes.Equal(a, b) {
				why = append(why, prefix+f+firstDiff(a, b))
			}
		}
	}
	files := map[string]bool{}
	for _, dir := range []string{jr, gr} {
		root := filepath.Join(dir, "files")
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
		a, errA := os.ReadFile(filepath.Join(jr, "files", f))
		b, errB := os.ReadFile(filepath.Join(gr, "files", f))
		switch {
		case errA != nil:
			why = append(why, f+" (only jsgo)")
		case errB != nil:
			why = append(why, f+" (only Java)")
		case j.intermediate[f]:
			// Both wrote it; its content is not compared (see projectJob).
		case !bytes.Equal(a, b):
			why = append(why, f+firstDiff(a, b))
		}
	}
	if len(why) == 0 {
		return clIdentical, fmt.Sprintf("%d files", len(names))
	}
	if javaExceptionRe.Match(javaStderr) {
		return clJavaError, firstMatchLine(javaExceptionRe, javaStderr)
	}
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
	return clDiff, strings.Join(why, "; ")
}

// millis is the summed wall time of the job's steps on one side.
func (j *projJob) millis(side string) string {
	total := 0
	for i := range j.steps {
		ms, err := strconv.Atoi(readTrim(filepath.Join(j.stepRes(j.res(side), i), "millis")))
		if err != nil {
			return ""
		}
		total += ms
	}
	return strconv.Itoa(total)
}

func summarizeProjects(wd, javaDir, goDir string, jobs []*projJob) int {
	counts := map[string]int{}
	var lines []string
	fmt.Printf("%-44s %-12s %9s %9s  %s\n", "job", "class", "java ms", "jsgo ms", "detail")
	for _, j := range jobs {
		class, detail := classifyProject(javaDir, goDir, j)
		counts[class]++
		jms, gms := j.millis(javaDir), j.millis(goDir)
		lines = append(lines, strings.Join([]string{j.project.Name, j.name, class, jms, gms, detail}, "\t"))
		if len(detail) > 600 {
			detail = detail[:600] + "..."
		}
		fmt.Printf("%-44s %-12s %9s %9s  %s\n", j.id(), class, jms, gms, detail)
	}
	os.WriteFile(filepath.Join(wd, "projects.tsv"), []byte(strings.Join(lines, "\n")+"\n"), 0o666)
	fmt.Printf("TOTAL %d jobs:", len(jobs))
	for _, c := range classes {
		fmt.Printf(" %s %d", c, counts[c])
	}
	fmt.Println()
	return 0
}
