package resources

import (
	"crypto/sha256"
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
	{"com.google.javascript.jscomp.RemoveUnusedCode", "js/polyfills.txt", true, "len=11051 sha256=47155f596ff9f7d21a21c9597157f760ca06617fd2d22c253db8bca387c030ca"},
	{"com.google.javascript.jscomp.js.RuntimeJsLibManager", "transpilation_libs.txt", true, "len=3561 sha256=3305d19c3dd4661f76d2ad3c2fc46d953e0a1b06c789c9d2423246ddf6884be6"},
	{"com.google.javascript.jscomp.Compiler", "js/es6_runtime.js", true, "len=1178 sha256=69949d204bd59fa5e5396e3a937d2b63232f6a30fdb62191fdfcf632385d7786"},
	{"com.google.javascript.jscomp.Compiler", "js/es6/util/arrayiterator.js", true, "len=1368 sha256=d2817f96df566b92d34946a6da2d67fdab62cbb5eb0270b78e3e3566ab77435c"},
	{"com.google.javascript.jscomp.Compiler", "js/nonexistent.js", false, "java.lang.RuntimeException: No such resource: js/nonexistent.js"},
	{"com.google.javascript.jscomp.Compiler", "/com/google/javascript/jscomp/js/polyfills.txt", true, "len=11051 sha256=47155f596ff9f7d21a21c9597157f760ca06617fd2d22c253db8bca387c030ca"},
	{"com.google.javascript.jscomp.Compiler", "js/../js/polyfills.txt", false, "java.lang.RuntimeException: No such resource: js/../js/polyfills.txt"},
	{"com.google.javascript.jscomp.Compiler", "js//polyfills.txt", false, "java.lang.RuntimeException: No such resource: js//polyfills.txt"},
	{"com.google.javascript.jscomp.Compiler", "./js/polyfills.txt", false, "java.lang.RuntimeException: No such resource: ./js/polyfills.txt"},
	{"com.google.javascript.jscomp.Compiler", "js/", true, "len=0 sha256=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	{"com.google.javascript.jscomp.Compiler", "js", true, "len=0 sha256=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	{"com.google.javascript.jscomp.Compiler$CodeBuilder", "js/polyfills.txt", true, "len=11051 sha256=47155f596ff9f7d21a21c9597157f760ca06617fd2d22c253db8bca387c030ca"},
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

// The embedded resources must be valid UTF-8: on such input Java's decoding
// (InputStreamReader with UTF_8) is the identity, which is what
// resourceLoader_charStreamsToString relies on.
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
		if path == "externs.zip" {
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
	// plus externs.zip.
	if n != 158 {
		t.Errorf("%d embedded resources, want 158", n)
	}
}

func TestLoadPropertiesMapMatchesJava(t *testing.T) {
	// PropertiesParser_Parse of the jar's js/polyfills.txt is in the oracle;
	// loadPropertiesMap must read the same text and give the same map.
	text := ResourceLoader_LoadTextResource("com.google.javascript.jscomp.RemoveUnusedCode", "js/polyfills.txt")
	for _, r := range readPropertiesOracle(t) {
		if r.input != text {
			continue
		}
		var sb []byte
		for k, v := range ResourceLoader_LoadPropertiesMap("com.google.javascript.jscomp.RemoveUnusedCode", "js/polyfills.txt").All() {
			sb = append(sb, escape(k+"="+v)...)
			sb = append(sb, '\n')
		}
		if string(sb) != r.want {
			t.Errorf("loadPropertiesMap(polyfills.txt):\n--- java\n%s--- go\n%s", r.want, sb)
		}
		return
	}
	t.Fatal("polyfills.txt not in the oracle")
}
