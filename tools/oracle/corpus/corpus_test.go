package corpus

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestKey(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"mathjs@15.2.0/package/lib/esm/a/b.js", "lib/esm/a/b.js"},
		{"@scope/pkg@1.0.0/package/x.js", "x.js"},
		{"some/dir/x.js", "some/dir/x.js"},
		{"pkg@1.0.0/package", "pkg@1.0.0/package"},
	} {
		if got := Key(tc.in); got != tc.want {
			t.Errorf("Key(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestMatch(t *testing.T) {
	ex := []string{"mathjs@15.2.0/package/lib/esm/sparse/csPost.js"}
	for _, tc := range []struct {
		path string
		want bool
	}{
		{"mathjs@15.2.0/package/lib/esm/sparse/csPost.js", true},
		{"/home/u/.cache/corpus-holdout/mathjs@15.2.0/package/lib/esm/sparse/csPost.js", true},
		{"package/lib/esm/sparse/csPost.js", true},
		{"lib/esm/sparse/csPost.js", true},
		{`lib\esm\sparse\csPost.js`, filepath.Separator == '\\'},
		{"mathjs@15.2.0/package/lib/esm/sparse/xcsPost.js", false},
		{"mathjs@15.2.0/package/lib/esm/sparse/csPost.js.map", false},
		{"mathjs@15.2.0/package/lib/cjs/sparse/csPost.js", false},
	} {
		if got := Match(ex, tc.path); got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestParseList(t *testing.T) {
	got := ParseList([]byte("# comment\n\n a@1/package/x.js \n#b\nc.js\n"))
	if want := []string{"a@1/package/x.js", "c.js"}; !slices.Equal(got, want) {
		t.Errorf("ParseList = %q, want %q", got, want)
	}
}

// TestRepositoryList checks that the list is found from a package directory
// and that ReadList and Filter drop its files.
func TestRepositoryList(t *testing.T) {
	ex, err := Excluded()
	if err != nil {
		t.Fatal(err)
	}
	if len(ex) == 0 {
		t.Fatal("empty list")
	}
	list := filepath.Join(t.TempDir(), "files.txt")
	if err := os.WriteFile(list, []byte("a@1.0.0/package/x.js\n\n"+ex[0]+"\nb.js\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := ReadList(list)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"a@1.0.0/package/x.js", "b.js"}; !slices.Equal(got, want) {
		t.Errorf("ReadList = %q, want %q", got, want)
	}
	if got := Filter([]string{"/r/" + ex[len(ex)-1], "c.js"}); !slices.Equal(got, []string{"c.js"}) {
		t.Errorf("Filter = %q", got)
	}
	if !IsExcluded(ex[0]) || IsExcluded("c.js") {
		t.Error("IsExcluded")
	}
}
