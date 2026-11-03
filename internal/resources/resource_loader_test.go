package resources

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"testing"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

var resourceLoaderVectors = []struct {
	path   string
	exists bool
	result string
}{
	{JsDir + "polyfills.txt", true, "len=11721 sha256=8986c52293e50e1b416449301929f7ef97d3dab2c21702ee6c0b03f75ed08491"},
	{JsDir + "transpilation_libs.txt", true, "len=3965 sha256=bb817d2962ab4d5e29948cc553a61aced001b7260136cc2ad51b76707867b22e"},
	{JsDir + "es6_runtime.js", true, "len=1234 sha256=9cde552822c9a16038349459630de758cbff0ddd64e38c05e14c905e3b1e614b"},
	{JsDir + "es6/util/arrayiterator.js", true, "len=1395 sha256=a2505a52eb1bae870a5975927ffdfa357c11de977612c9dceb6360e0d5b5c26a"},
	{JsDir + "nonexistent.js", false, "java.lang.RuntimeException: No such resource: " + JsDir + "nonexistent.js"},
	{"/" + JsDir + "polyfills.txt", true, "len=11721 sha256=8986c52293e50e1b416449301929f7ef97d3dab2c21702ee6c0b03f75ed08491"},
	{JsDir + "../js/polyfills.txt", false, "java.lang.RuntimeException: No such resource: " + JsDir + "../js/polyfills.txt"},
	{"com/google/javascript/jscomp/js//polyfills.txt", false, "java.lang.RuntimeException: No such resource: com/google/javascript/jscomp/js//polyfills.txt"},
	{"./" + JsDir + "polyfills.txt", false, "java.lang.RuntimeException: No such resource: ./" + JsDir + "polyfills.txt"},
	{JsDir, true, "len=0 sha256=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	{"com/google/javascript/jscomp/js", true, "len=0 sha256=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
}

func TestResourceLoaderMatchesJava(t *testing.T) {
	for _, v := range resourceLoaderVectors {
		var got string
		jutil.Try(func() {
			s := LoadText(v.path)
			got = "len=" + strconv.Itoa(jutil.StringLength(s)) + fmt.Sprintf(" sha256=%x", sha256.Sum256([]byte(s)))
		}, func(e jutil.Throwable) { got = e.Error() })
		if got != v.result {
			t.Errorf("LoadText(%q) = %s, want %s", v.path, got, v.result)
		}
		if e := Exists(v.path); e != v.exists {
			t.Errorf("Exists(%q) = %v, want %v", v.path, e, v.exists)
		}
		_, err := LoadBytes(v.path)
		if (err == nil) != v.exists {
			t.Errorf("LoadBytes(%q) err = %v, exists = %v", v.path, err, v.exists)
		}
		if !v.exists && !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("LoadBytes(%q) err = %v, want errors.Is(fs.ErrNotExist)", v.path, err)
		}
	}
}

func TestEmbeddedResourcesAreValidUTF8(t *testing.T) {
	n := 0
	err := fs.WalkDir(classpath, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		n++
		data, err := fs.ReadFile(classpath, path)
		if err != nil {
			return err
		}
		if path == "externs.zip" || path == "runtime_libs.typedast" {
			return nil // a binary resource
		}
		if !utf8.Valid(data) {
			t.Errorf("%s is not valid UTF-8", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if n != 177 {
		t.Errorf("%d embedded resources, want 177", n)
	}
}

func TestEmbeddedRuntimeLibsTypedAst(t *testing.T) {
	data, err := LoadBytes("runtime_libs.typedast")
	if err != nil {
		t.Fatalf("runtime_libs.typedast: %v", err)
	}
	sum := sha256.Sum256(data)
	if got := hex.EncodeToString(sum[:]); len(data) != 3540605 ||
		got != "894791ed6d0069bbfba7a26410bfe4f2600a3bfe21f989967e10f8c42429b968" {
		t.Errorf("runtime_libs.typedast: %d bytes, sha256 %s", len(data), got)
	}
}
