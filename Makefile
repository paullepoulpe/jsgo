SHELL := /bin/bash
CLOSURE_VERSION := v20260415
CLOSURE_JAR_SHA256 := 20313094c6e9cac8c685931e0b3fd67f82e1534c30448cdff404ee350460e391
CLOSURE_JAR_URL := https://repo1.maven.org/maven2/com/google/javascript/closure-compiler/$(CLOSURE_VERSION)/closure-compiler-$(CLOSURE_VERSION).jar
CLOSURE_JAR := .cache/closure-compiler.jar
CLOSURE_EXTERNS_DIR := .cache/closure-externs
CLOSURE_FLAGS := --checks_only --jscomp_error=checkTypes --warning_level=VERBOSE --language_in=ECMASCRIPT_NEXT
CLOSURE_FLAGS_STRICT := $(CLOSURE_FLAGS) --jscomp_warning=strictCheckTypes

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
	# Extract Closure's upstream externs (es3.js, es5.js, es6.js, browser/*.js)
	# from the nested externs.zip inside the jar so the Go externs loader can
	# read them from $(CLOSURE_EXTERNS_DIR). Mirrors how Closure's
	# AbstractCommandLineRunner.getBuiltinExterns reads /externs.zip from the
	# classpath.
	rm -rf $(CLOSURE_EXTERNS_DIR)
	mkdir -p $(CLOSURE_EXTERNS_DIR)
	unzip -p $@ externs.zip > $(CLOSURE_EXTERNS_DIR)/externs.zip
	unzip -o -q $(CLOSURE_EXTERNS_DIR)/externs.zip -d $(CLOSURE_EXTERNS_DIR)
	rm $(CLOSURE_EXTERNS_DIR)/externs.zip

oracle: $(CLOSURE_JAR)
	bash scripts/gen-oracles.sh

clean:
	rm -rf bin/ .cache/
