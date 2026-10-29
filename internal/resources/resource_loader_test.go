package resources

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"strconv"
	"testing"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Expected values were produced by the pinned jar's ResourceLoader on JDK 21
// (loadTextResource: UTF-16 length and SHA-256 of the UTF-8 text, or the
// exception; resourceExists). Closure has no ResourceLoaderTest.
var resourceLoaderVectors = []struct {
	clazz, path string
	exists      bool
	result      string
}{
	{"com.google.javascript.jscomp.RemoveUnusedCode", "js/polyfills.txt", true, "len=11379 sha256=4b501797a9eca44c23281dfdb1aef1e8b7af88e48f2adb0b4025542009dc69e7"},
	{"com.google.javascript.jscomp.js.RuntimeJsLibManager", "transpilation_libs.txt", true, "len=3934 sha256=0857894f26d26542f0db7dc9090861133d5ce80e0a716a60fc74d8f6339e0d37"},
	{"com.google.javascript.jscomp.Compiler", "js/es6_runtime.js", true, "len=1234 sha256=9cde552822c9a16038349459630de758cbff0ddd64e38c05e14c905e3b1e614b"},
	{"com.google.javascript.jscomp.Compiler", "js/es6/util/arrayiterator.js", true, "len=1395 sha256=a2505a52eb1bae870a5975927ffdfa357c11de977612c9dceb6360e0d5b5c26a"},
	{"com.google.javascript.jscomp.Compiler", "js/nonexistent.js", false, "java.lang.RuntimeException: No such resource: js/nonexistent.js"},
	{"com.google.javascript.jscomp.Compiler", "/com/google/javascript/jscomp/js/polyfills.txt", true, "len=11379 sha256=4b501797a9eca44c23281dfdb1aef1e8b7af88e48f2adb0b4025542009dc69e7"},
	{"com.google.javascript.jscomp.Compiler", "js/../js/polyfills.txt", false, "java.lang.RuntimeException: No such resource: js/../js/polyfills.txt"},
	{"com.google.javascript.jscomp.Compiler", "js//polyfills.txt", false, "java.lang.RuntimeException: No such resource: js//polyfills.txt"},
	{"com.google.javascript.jscomp.Compiler", "./js/polyfills.txt", false, "java.lang.RuntimeException: No such resource: ./js/polyfills.txt"},
	{"com.google.javascript.jscomp.Compiler", "js/", true, "len=0 sha256=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	{"com.google.javascript.jscomp.Compiler", "js", true, "len=0 sha256=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	{"com.google.javascript.jscomp.Compiler$CodeBuilder", "js/polyfills.txt", true, "len=11379 sha256=4b501797a9eca44c23281dfdb1aef1e8b7af88e48f2adb0b4025542009dc69e7"},
	{"com.google.javascript.jscomp.Compiler", "", true, "len=0 sha256=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
}

func TestResourceLoaderMatchesJava(t *testing.T) {
	for _, v := range resourceLoaderVectors {
		var got string
		jutil.Try(func() {
			s := ResourceLoader_LoadTextResource(v.clazz, v.path)
			got = "len=" + strconv.Itoa(jutil.StringLength(s)) + fmt.Sprintf(" sha256=%x", sha256.Sum256([]byte(s)))
		}, func(e jutil.Throwable) { got = e.Error() })
		if got != v.result {
			t.Errorf("loadTextResource(%s, %q) = %s, Java %s", v.clazz, v.path, got, v.result)
		}
		if e := ResourceLoader_ResourceExists(v.clazz, v.path); e != v.exists {
			t.Errorf("resourceExists(%s, %q) = %v, Java %v", v.clazz, v.path, e, v.exists)
		}
	}
}

// The embedded text resources must be valid UTF-8: on such input Java's
// decoding (InputStreamReader with UTF_8) is the identity, which is what
// resourceLoader_charStreamsToString relies on. (externs.zip and
// runtime_libs.typedast are binary resources, read as bytes only.)
func TestEmbeddedResourcesAreValidUTF8(t *testing.T) {
	n := 0
	err := fs.WalkDir(resourceLoader_classpath, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		n++
		data, err := fs.ReadFile(resourceLoader_classpath, path)
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
	// unzip -Z1 closure-compiler.jar 'com/google/javascript/jscomp/js/*' minus directories and .class files,
	// plus externs.zip (runtime_libs.typedast is embedded on its own).
	if n != 170 {
		t.Errorf("%d embedded resources, want 170", n)
	}
}

// runtime_libs.typedast is the pinned jar's entry, byte for byte
// (unzip -p closure-compiler.jar runtime_libs.typedast | sha256sum).
func TestEmbeddedRuntimeLibsTypedAst(t *testing.T) {
	data := ResourceLoader_GetResourceAsStream("com.google.javascript.jscomp.Compiler", "/runtime_libs.typedast")
	if data == nil {
		t.Fatal("/runtime_libs.typedast is not embedded")
	}
	sum := sha256.Sum256(data)
	if got := hex.EncodeToString(sum[:]); len(data) != 3530744 ||
		got != "feb246c5d2a4074923caa36814fba7913094789f5fe3aa119a528271e4dad601" {
		t.Errorf("runtime_libs.typedast: %d bytes, sha256 %s", len(data), got)
	}
}
