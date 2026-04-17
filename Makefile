SHELL := /bin/bash
CLOSURE_VERSION := v20260415
CLOSURE_JAR_SHA256 := 20313094c6e9cac8c685931e0b3fd67f82e1534c30448cdff404ee350460e391
CLOSURE_JAR_URL := https://repo1.maven.org/maven2/com/google/javascript/closure-compiler/$(CLOSURE_VERSION)/closure-compiler-$(CLOSURE_VERSION).jar
CLOSURE_JAR := .cache/closure-compiler.jar
CLOSURE_FLAGS := --checks_only --jscomp_error=checkTypes --warning_level=VERBOSE --language_in=ECMASCRIPT_NEXT

.PHONY: all build test ci-guards oracle clean

all: build test ci-guards

build:
	go build -o bin/jsgo ./cmd/jsgo

test:
	go test ./...

ci-guards:
	bash scripts/ci-guards.sh

$(CLOSURE_JAR):
	mkdir -p .cache
	curl -sSL -o $@ $(CLOSURE_JAR_URL)
	echo "$(CLOSURE_JAR_SHA256)  $@" | sha256sum -c -

oracle: $(CLOSURE_JAR)
	@for f in $$(find testdata/fixtures -name '*.js' | sort); do \
	    rel=$${f#testdata/fixtures/}; \
	    out=testdata/oracle/$${rel%.js}.txt; \
	    mkdir -p $$(dirname $$out) || exit 1; \
	    tmp=$$(mktemp); \
	    java -jar $(CLOSURE_JAR) $(CLOSURE_FLAGS) --js $$f > $$tmp 2>&1; ec=$$?; \
	    { cat $$tmp; echo "---EXIT---"; echo $$ec; } > $$out; \
	    rm -f $$tmp; \
	    echo "oracle: $$out"; \
	done

clean:
	rm -rf bin/ .cache/
