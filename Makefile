# See PLAN.md for the roadmap and PORTING.md for the translation rules.
.PHONY: all build vet test oracle-tools

all: build vet test

build:
	go build ./...

# Vet everything under internal/ (the whole port); check that
# internal/jscomp/pending_gen.go is current and duplicates no real
# declaration (without the pinned jar, only the latter); and check the
# licence notices (tools/licenses; see PORTING.md, top).
vet:
	go vet ./internal/...
	go run ./tools/pending -check
	go run ./tools/licenses -check

test:
	go test ./...

# Pinned Closure Compiler jar + sources, and the Java oracle tools (needs a JDK).
oracle-tools:
	bash tools/oracle/build.sh
