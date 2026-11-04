package jutil

import (
	"sync"
	"testing"
)

func TestIsValidPattern(t *testing.T) {
	cases := []struct {
		pattern string
		want    bool
	}{
		{"", true},
		{`.*\.closure\.js`, true},
		{`^foo/(bar|baz)[0-9]+\.js$`, true},
		{`\Q([unclosed\E+`, true},
		{`\Qunclosed_quote([`, true},
		{`[\Q]-\E]`, true},
		{`\Eabc`, true},
		{`[]a]`, true},
		{`[^]a]`, true},
		{`[\d\w\s\v\h]`, true},
		{`(foo)\1`, true},
		{`\bword\b\z`, true},
		{`a{2}{3}`, true},
		{`a{2000}`, true},
		{`\`, false},
		{`\q`, false},
		{`\M`, false},
		{`[abc`, false},
		{`[]`, false},
		{`[^]`, false},
		{`(abc`, false},
		{`abc)`, false},
		{`*abc`, false},
		{`[z-a]`, false},
	}
	for _, tc := range cases {
		if got := IsValidPattern(tc.pattern); got != tc.want {
			t.Errorf("IsValidPattern(%q) = %v, want %v", tc.pattern, got, tc.want)
		}
	}
}

func TestPatternCompileCacheAccounting(t *testing.T) {
	const uniqueRe = `^round96_unique_cache_test_[0-9]+$`
	patternCache.Delete(uniqueRe)

	before := patternCacheSize.Load()
	var wg sync.WaitGroup
	start := make(chan struct{})
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_ = Pattern_Compile(uniqueRe)
		}()
	}
	close(start)
	wg.Wait()

	after := patternCacheSize.Load()
	if before < patternCacheLimit {
		if after != before+1 {
			t.Fatalf("patternCacheSize after 16 concurrent compiles of same regex = %d, want %d", after, before+1)
		}
	}
	// Repeat hit does not increment patternCacheSize.
	_ = Pattern_Compile(uniqueRe)
	if got := patternCacheSize.Load(); got != after {
		t.Fatalf("patternCacheSize after cache hit = %d, want %d", got, after)
	}

	// At patternCacheLimit, compiling an uncached regex neither stores nor grows patternCacheSize.
	saved := patternCacheSize.Swap(patternCacheLimit)
	defer patternCacheSize.Store(saved)

	const overflowRe = `^round96_overflow_uncached_[a-z]+$`
	patternCache.Delete(overflowRe)
	_ = Pattern_Compile(overflowRe)
	if got := patternCacheSize.Load(); got != patternCacheLimit {
		t.Fatalf("patternCacheSize at limit = %d, want %d", got, patternCacheLimit)
	}
	if _, ok := patternCache.Load(overflowRe); ok {
		t.Fatalf("expected %q not to be cached when patternCacheSize == patternCacheLimit", overflowRe)
	}
}
