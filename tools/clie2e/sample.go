package main

// Sampling for the sweep's flag configurations: a stratified random sample
// of a corpus list (-sample), and each configuration's share of it
// (-per-config).

import (
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// packageOf is a corpus path's stratum: its npm package and version, the
// first component of the path ("@scope/name@version" for scoped packages).
func packageOf(path string) string {
	parts := strings.SplitN(path, "/", 3)
	if strings.HasPrefix(parts[0], "@") && len(parts) > 1 {
		return parts[0] + "/" + parts[1]
	}
	return parts[0]
}

// stratifiedSample draws n files of the list, stratified by package. Each
// package's share is proportional to the square root of its size, so that
// the small packages are represented, except Closure Library's, which gets
// twice that: it is the code that --define, the Closure primitives and the
// type-based flags are about. A package smaller than its share is taken
// whole, and the rest is shared again. Inside a package the files are drawn
// at random. With maxSize > 0, files larger than maxSize bytes (under root)
// are left out first. The sample keeps the list's order.
func stratifiedSample(files []string, n int, seed int64, root string, maxSize int64) []string {
	if maxSize > 0 {
		var kept []string
		for _, f := range files {
			if st, err := os.Stat(filepath.Join(root, f)); err == nil && st.Size() <= maxSize {
				kept = append(kept, f)
			}
		}
		files = kept
	}
	if n >= len(files) {
		return files
	}
	byPkg := map[string][]string{}
	var pkgs []string
	for _, f := range files {
		p := packageOf(f)
		if byPkg[p] == nil {
			pkgs = append(pkgs, p)
		}
		byPkg[p] = append(byPkg[p], f)
	}
	sort.Strings(pkgs)
	weight := func(p string) float64 {
		w := math.Sqrt(float64(len(byPkg[p])))
		if strings.HasPrefix(p, "google-closure-library@") {
			w *= 2
		}
		return w
	}
	quota := map[string]int{}
	open := append([]string{}, pkgs...)
	left := n
	for left > 0 && len(open) > 0 {
		total := 0.0
		for _, p := range open {
			total += weight(p)
		}
		var still []string
		given := 0
		for _, p := range open {
			q := int(math.Floor(float64(left) * weight(p) / total))
			if room := len(byPkg[p]) - quota[p]; q >= room {
				q = room
			} else {
				still = append(still, p)
			}
			quota[p] += q
			given += q
		}
		if given == 0 {
			// Rounding left a few: one more each to the largest packages
			// that still have room.
			sort.SliceStable(still, func(a, b int) bool { return weight(still[a]) > weight(still[b]) })
			for _, p := range still {
				if given == left {
					break
				}
				quota[p]++
				given++
			}
			if given == 0 {
				break
			}
		}
		left -= given
		open = still
	}
	r := rand.New(rand.NewSource(seed))
	chosen := map[string]bool{}
	for _, p := range pkgs {
		fs := append([]string{}, byPkg[p]...)
		r.Shuffle(len(fs), func(i, j int) { fs[i], fs[j] = fs[j], fs[i] })
		for _, f := range fs[:quota[p]] {
			chosen[f] = true
		}
	}
	var out []string
	for _, f := range files {
		if chosen[f] {
			out = append(out, f)
		}
	}
	return out
}

// spread returns the files in an order where each package's files are
// evenly spaced (a random permutation inside each package), so that any
// run of consecutive files is itself about stratified.
func spread(files []string, r *rand.Rand) []string {
	byPkg := map[string][]string{}
	for _, f := range files {
		byPkg[packageOf(f)] = append(byPkg[packageOf(f)], f)
	}
	type keyed struct {
		f string
		k float64
	}
	var ks []keyed
	pkgs := make([]string, 0, len(byPkg))
	for p := range byPkg {
		pkgs = append(pkgs, p)
	}
	sort.Strings(pkgs)
	for _, p := range pkgs {
		fs := byPkg[p]
		r.Shuffle(len(fs), func(i, j int) { fs[i], fs[j] = fs[j], fs[i] })
		off := r.Float64()
		for i, f := range fs {
			ks = append(ks, keyed{f, (float64(i) + off) / float64(len(fs))})
		}
	}
	sort.SliceStable(ks, func(a, b int) bool { return ks[a].k < ks[b].k })
	out := make([]string, len(ks))
	for i, k := range ks {
		out[i] = k.f
	}
	return out
}

// assignFiles gives each configuration perConfig files (or its own count),
// dealing the files out in passes so that every file is used about as often
// as any other: pass p deals a fresh spread order, a run of consecutive
// files to each configuration in turn. A configuration's files keep the
// order they were dealt in (which decides its variants).
func assignFiles(files []string, confs []config, perConfig int, seed int64) map[string][]string {
	r := rand.New(rand.NewSource(seed + 1))
	var deck []string
	next := 0
	draw := func() string {
		if next == len(deck) {
			deck = spread(files, r)
			next = 0
		}
		next++
		return deck[next-1]
	}
	out := map[string][]string{}
	for _, c := range confs {
		want := perConfig
		if c.files > 0 {
			want = c.files
		}
		if want > len(files) {
			want = len(files)
		}
		seen := map[string]bool{}
		var mine []string
		for len(mine) < want {
			f := draw()
			if !seen[f] {
				seen[f] = true
				mine = append(mine, f)
			}
		}
		out[c.name] = mine
	}
	return out
}
